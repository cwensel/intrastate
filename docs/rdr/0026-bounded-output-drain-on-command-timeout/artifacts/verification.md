# Verification — RDR 0026 bounded-output-drain-on-command-timeout

Stage 8 launch verification. Phase 3a records `FAIL-N` entries; Phase 3b
records `ADV-N` entries. The two phases are written independently and
neither reads the other's findings before recording its own.

## Phase 3b — adversarial review

Anchor: Trade-offs → Failure Modes (F1..F6), plus `0026:C1` `precedence:`
and `bound:` where a Failure Mode's promise is stated as a number.

Three failure modes, each with the test that catches it. **All three tests
currently FAIL against the implementation at `ee72e35`** — they are real
defects, not regression traps.

Tests added:

- `internal/cli/cmdbind/adversarial_0026_test.go` (in-package; ADV-1, ADV-2)
- `internal/cli/cmdbind/adversarial_bound_0026_test.go` (external; ADV-3)

ADV-1 and ADV-2 are in-package on purpose. Both defects live inside
`readBounded`'s terminal-condition report and inside `spawn`'s fold of the
two reports. Driving `readBounded` against a real `os.Pipe` with a real held
write end reproduces the exact state `spawn` reaches — a parent that has
armed `now + DrainGrace` on its own read end while a writer it cannot reach
holds the other — and names which of the three ordered conditions the drain
reported. Nothing is mocked: the pipes are real, the deadline is the real
one the parent arms, and the loop under test is the shipped one.

---

### ADV-1 — a held stderr past the read bound is laundered into overflow, and the refusal is WITHHELD

**Failure mode**: F1. "an escaped grandchild holds a pipe past the bound —
the invocation refuses ... with `Err` naming the pipe(s) and the bound."
C1 `residue:` states the same guarantee from the other side: "Leakage of a
process, or of THAT process's output, is admitted; a withheld refusal is
not."

**The defect**. `spawn` drains stderr with `readBounded(errR, StdoutCap)` —
a limit of exactly `StdoutCap`, where stdout uses `StdoutCap+1`.
`readBounded` applies the overflow rank whenever `len(kept) >= limit`,
whatever the terminal condition. So a verbose child whose stderr reaches
1 MiB and is then held by a `setsid(2)` escapee reports `drainOverflow`
rather than `drainHeld`.

`spawn` then folds the two reports. Its overflow arm consults only the
STDOUT report (`outReport == drainOverflow || len(inv.stdout) > StdoutCap`)
— there is no stderr overflow refusal in the contract, because stderr is
tailed rather than capped — and `heldPipes` counts only `drainHeld`. The
held condition is therefore consulted by NEITHER arm: `spawn` falls through
to `return inv, nil` and the invocation succeeds with no error at all. If
stdout was whole, the caller parses a perfectly good answer while stderr sat
held by an unreachable writer — the refusal F1 owes is never delivered.

C1 `precedence:` scopes the overflow rank deliberately: "Overflow outranks
held on the same drain: a drain past `StdoutCap` REFUSES the existing
overflow error whatever its terminal condition, since the bytes were read
and the pipe's state is then not what the invocation turns on." The rank
exists BECAUSE overflow refuses. On stderr it does not refuse, so ranking a
held stderr as overflow **discards** the condition instead of outranking it.

**Reachability**: a tool emitting >1 MiB of stderr (a build tool, a linter)
that also leaves a `setsid` helper holding stderr. Uncommon, but the failure
is silent — it produces an accepted, parsed answer.

**Test**: `TestAdv1_AHeldStderrPastTheReadBoundStillReportsHeldAndStillRefuses`
(`internal/cli/cmdbind/adversarial_0026_test.go`).

**Currently fails**: YES.

```
a stderr pipe HELD by an unreachable writer past the drain bound reported
overflow (past the read bound), not held.
```

---

### ADV-2 — the liveness ceiling reports HELD on a paced tail that reaches a real EOF, with no writer at all

**Failure mode**: F4. "bytes the direct child wrote can be lost at the
bound; the refusal still fires, so the loss is bounded to a refused
invocation, never a parsed one." F4 scopes the loss to a pipe an escapee
still holds. This drain's pipe has NO holder: the writer closes, EOF
arrives, and every byte is deliverable.

**The defect**. `readBounded`'s liveness ceiling (deviation D4) is checked
at the TOP of the loop, before the read, and expires `2·WaitDelay` after the
drain entered the grace. It is armed once and progress never re-arms it. A
writer whose per-read idle gaps are all UNDER `DrainGrace` — so the grace
itself never fires, and the re-arm mechanism is working exactly as specified
— but whose tail takes longer than `2·WaitDelay` to finish is cut off
mid-tail and reported HELD.

That is the failure C1 `precedence:` forbids by name, twice: "the grace
bounds the IDLE GAP between reads and never the size or **total duration**
of the tail", and "A single absolute deadline would instead bound the whole
remaining tail against the clock, which turns a slow-arriving tail into a
**false held-pipe report on a pipe with no writer at all**". The ceiling IS
that single absolute deadline; only its constant is larger.

Deviation D4 asserts the ceiling "never binds on a drain that reaches EOF"
and is "deliberately slack — 1 s against a measured 0.605 s worst case".
That claim does not hold. 0.605 s is the one shape A9 happened to measure
(1 MiB across 32 chunks); nothing in the contract bounds a tail's total
duration, and the measured case already sits at ~1.5x of the ceiling. A
loaded host, a slower writer cadence, or more chunks reaches it. The test's
witness — 40 chunks x 4 KiB with 40 ms gaps — is a smaller payload than
A9's and still breaches it.

The pacing is deliberately inside the grace on every gap (40 ms against a
50 ms `DrainGrace`), so the failure cannot be blamed on a gap the grace
legitimately refuses: the ONLY thing that ends this drain early is the
total-duration ceiling.

**Test**: `TestAdv2_APacedTailReachingARealEOFIsNeverReportedHeld`
(`internal/cli/cmdbind/adversarial_0026_test.go`).

**Currently fails**: YES.

```
a paced tail with NO holder — every idle gap 40ms, strictly under the 50 ms
grace, and the writer closing for a real EOF — reported held
(os.ErrDeadlineExceeded on the final read) with 131072 of 163840 bytes
recovered.
```

20% of a complete, deliverable answer lost, and the invocation refused.

---

### ADV-3 — the ceiling STACKS on the parent's join timer, so the trickler overruns the committed return bound

**Failure mode**: F1 and F2. F1 promises the trickling escapee is REFUSED,
and the existing trickle oracle asserts exactly that. What no oracle asserts
is WHEN. C1 `bound:` publishes the number as a commitment: "Per invocation,
the COMMITTED return bound is timeout + 2·WaitDelay, asserted with a
scheduling tolerance of +100 ms ... a return past the total but inside the
tolerance is a pass, and **past the tolerance is a contract violation**."
F2's consequence is what makes it load-bearing: a caller sizing its own
watchdog from `timeout` plus the published bound is killed mid-invocation.

**The defect**. The ceiling is `2·WaitDelay` measured from the DRAIN's entry
into the grace, not carved out of the parent's allowance. The parent has
already spent its full `WaitDelay` on the join timer before it arms the
grace at all, so on a trickler the two costs STACK:

```
WaitDelay (join timer) + DrainGrace (first grace window) + 2·WaitDelay (ceiling)
  = 500 + 50 + 1000 = ~1550 ms
```

against a committed drain allowance of `2·WaitDelay` = 1000 ms plus a
100 ms tolerance. Measured end to end through the real `Reader.Read`:
**1.825 s against a 1.1 s budget** — 725 ms past the tolerance.

Deviation D4 states the ceiling is "taken from the committed
per-invocation total C1 `bound:` already publishes (`2·WaitDelay`)". It is
taken from that total's VALUE but not from its BUDGET: measuring it from the
drain's own entry into the grace is what makes it stack rather than fit
inside. The trickler is the one shape the ceiling exists for, so this is the
ceiling's own motivating case overrunning the number it was added
underneath.

The existing trickle oracle
(`TestReq6_ATricklingEscapeeIsRefusedThoughItsDrainIsNeverEmpty`) does not
catch this: it asserts only `elapsed < 10s`, the declared deadline, and
never the committed total. The package already ships `escBoundFor`, the
helper that computes the committed total from the package constant; this
test uses it.

**Test**: `TestAdv3_ATricklingEscapeeStillReturnsInsideTheCommittedTotal`
(`internal/cli/cmdbind/adversarial_bound_0026_test.go`).

**Currently fails**: YES.

```
the trickling escapee refused after 1.825457625s, past the committed total
of 1.1s (2*WaitDelay = 1s, plus the record's stated 100ms scheduling
tolerance).
```

---

### Phase 3b verdict

**BLOCK.** Three added tests, three currently failing, three distinct
defects. ADV-2 and ADV-3 are both consequences of the liveness ceiling
(deviation D4) and are the evidence D4's "needs author decision" status was
opened for: as implemented, the ceiling violates C1 `precedence:`'s
prohibition on bounding a tail's total duration (ADV-2) and C1 `bound:`'s
committed return (ADV-3). ADV-1 is independent of D4 and is a plain fold
defect in `spawn` plus a limit asymmetry in `readBounded`.

Areas probed and found SOUND (no test added):

- **The C1 `precedence:` (b) happens-before edge.** `drains.Wait()` runs in
  its own goroutine feeding a `close(joined)`, and both the bounded and the
  non-pollable arms block on `<-joined` before `inv` is built. The chain is
  a valid happens-before edge; the join remains a real join on every path.
- **Goroutine liveness on the ceiling path.** `readBounded` returns from the
  ceiling branch like any other, so each drain goroutine's
  `defer drains.Done()` always fires. No leak, including on the ceiling.
- **fd lifetime.** No read end outlives `spawn`; the two creation-time
  defers are the only closes and nothing double-closes.
- **The cap boundary.** A stdout held at exactly `StdoutCap` reports held,
  not overflow (`limit` is `StdoutCap+1`), and a clean EOF with no deadline
  ever armed reports whole. The rewrite did not shift the boundary — the
  asymmetry is on stderr's limit alone, which is ADV-1.
- **`timeout` carrying a Detail it shouldn't.** The executor's deadline-first
  arm returns `readOutcome{class: ClassTimeout}` with no `err`, so the
  held-pipe reason cannot leak onto a `timeout` refusal (F6 as accepted).
- **The write arm's applied sense.** `errors.As(err, &hp)` matches the
  held-pipe sibling specifically and not `*ExecError` generally, so no
  pre-mutation failure shape is flipped to "may have been applied".

---

## Phase 3a — CoVe (chain-of-verification)

Recorded independently of the Phase 3b section above; neither read the
other before writing.

The method is the inverse of the implementation's: for each REQ under probe, name an input that would make a
CORRECT implementation visibly violate it, then run that input against the
build. Inputs were derived from the record (`0026:C1`'s seven labelled
lines, `0026:MVV` rows 1-6, S1-S10) WITHOUT reading the phase-1 test files,
so a shared blind spot between the tests and the implementation is
detectable.

Probed against `ee72e35`. Probes were driven through the real
`cmdbind.Reader.Read` and the real `accessor.Executor.Write`/`Read`/`Gate`
wherever the input class allowed, and against `readBounded`/`heldPipes`/
`invocation.detail` directly where the input had to be shaped at the pipe.

### FAIL-1 — a drain making progress inside the grace is refused as held once the tail outlives 2·WaitDelay

**Contract violated.** `0026:C1` `precedence:`:

> The deadline is re-armed at `now + DrainGrace` BEFORE EACH read of the
> final drain, so the grace bounds the IDLE GAP between reads and never the
> size or total duration of the tail.

and

> A re-implementation that shifts the cap boundary is a contract violation,
> not a refactor: a drain making progress keeps its deadline ahead of it and
> runs to EOF, and only a gap longer than the grace ends it as held.

and `0026:C1` `whole:` — "a drain that reaches EOF — before the timer, or on
its final read after it — is whole".

**Cause.** `internal/cli/cmdbind/cmdbind.go::readBounded` carries a second,
absolute bound on top of the per-read grace — a "liveness ceiling" of
`2·WaitDelay` measured from the drain's entry into the grace:

```go
if graced {
        if time.Now().After(ceiling) {
                if len(kept) >= limit {
                        return kept, drainOverflow
                }
                return kept, drainHeld
        }
        _ = r.SetReadDeadline(time.Now().Add(DrainGrace * time.Millisecond))
}
```

with `ceiling = time.Now().Add(2 * WaitDelay * time.Millisecond)`. The
discriminating variable is therefore the TOTAL DURATION of the tail, which
is exactly what the clause forbids and what MVV row 6(b) exists to catch —
row 6(b)'s own shape (32 chunks × 20 ms ≈ 640 ms) simply sits under the
1000 ms ceiling and so does not reach it. The ceiling is a rival bound that
appears nowhere in the record; `bound:` states `WaitDelay` is "the ONE
bound" and that `DrainGrace` "is NOT a second bound".

### Failing input (end-to-end, through `Reader.Read`)

A child that spawns a `SysProcAttr{Setsid: true}` grandchild inheriting
stdout/stderr and exits 0 at once. The grandchild writes the envelope in 60
chunks separated by 25 ms idle gaps — every gap HALF of `DrainGrace` (50 ms)
— then closes its stdio and exits, so the drain reaches a real EOF. Total
tail ≈ 1.5 s, past the 1000 ms ceiling.

```
elapsed=1.802071583s values=[] err=the child's stdout, stderr stayed held past the drain bound (exited 0)
```

**Observed**: `*accessor.HeldPipeError{Held: [stdout stderr], ExitStatus:
"exited 0"}` — the whole envelope is refused and never parsed.
**Expected per C1**: `drainWhole`, the envelope returned, no refusal.

**Control (proves the ceiling and not the grace is the cause)** — the same
fixture with 20 chunks (≈ 500 ms tail, UNDER the ceiling), gaps unchanged at
25 ms:

```
elapsed=771.130917ms values=[{k v false}] unreadable=[] err=<nil>
```

### Failing input (unit, at `readBounded`)

A pipe whose writer sends 40 chunks separated by 40 ms (still inside the
50 ms grace) and then CLOSES — no holder at all, the premortem P-2 shape:

```
bytes=38912 want=40960 report=1 (whole=0 held=1 overflow=2)
```

**Observed**: truncated at 38912 of 40960 bytes and reported `drainHeld`.
**Expected per C1**: all 40960 bytes, `drainWhole`. This is the same failure
signature MVV row 6(b) records for the rejected single-absolute-deadline
mechanism (there: 15.6% truncation, terminal `os.ErrDeadlineExceeded`), and
it reaches a pipe with NO WRITER AT ALL — the case C1 names as the reason
the re-arm exists.

MVV row 6(b) as literally written (32 chunks × 20 ms) PASSES:
`row6b bytes=1048576 want=1048576 report=0`. The row is therefore not a
witness for this defect; only a tail longer than the ceiling is.

### FAIL-2 — the committed per-invocation return bound is exceeded by a full WaitDelay when the drain ends at the liveness ceiling

**Contract violated.** `0026:C1` `bound:`:

> Per invocation, the COMMITTED return bound is timeout + 2·WaitDelay,
> asserted with a scheduling tolerance of +100 ms (the premortem's P-8
> figure) ... a return past the total but inside the tolerance is a pass,
> and past the tolerance is a contract violation.

**Cause.** The same liveness ceiling. The drain leg is entered at
`WaitDelay` after `Wait` returns (the join timer), and the ceiling then
allows a further `2·WaitDelay` measured from the drain's entry into the
grace. The drain leg can therefore run `3·WaitDelay`, not the `WaitDelay`
the total's arithmetic assumes, so a fixture that also spends the full
deadline leg (S8's competing-legs shape) returns at
`timeout + 3·WaitDelay`.

### Failing input

A child that spawns a `Setsid` grandchild inheriting stdout/stderr which
emits one byte to each pipe every 10 ms FOREVER — S9's companion ("an
escaped writer emitting a byte every 10 ms forever is REFUSED", "never
blocked, never EOF") — while the child itself sleeps past the declared
`timeout` of 2 s, so the deadline leg spends its own allowance rather than
donating it.

```
elapsed=3.56179625s total=3s tolerance=100ms values=[] err=the child's stdout, stderr stayed held past the drain bound (killed by signal)
```

**Observed**: 3.5617 s, 3.5623 s, 3.5553 s across 3/3 runs — an overshoot of
455-463 ms past `total + tolerance`, i.e. essentially one whole extra
`WaitDelay`.
**Expected per C1 `bound:`**: ≤ 3.0 s + 100 ms.

The refusal itself is correct — the trickle IS refused, so S9's companion
holds — and this entry is only about the elapsed time. Both FAIL entries
share one root cause: the ceiling is a second bound stacked OUTSIDE the one
the contract publishes, rather than a mechanism inside it.

### Probed and found conformant

Every item below was driven with an input chosen to break it, and did not
break.

- `precedence:` cap boundary — exactly `StdoutCap` reports `drainWhole` and
  is a value; `StdoutCap+1` reports `drainOverflow` and is a failure, at
  both `readBounded` and `spawn`'s overflow predicate. The `readBounded`
  rewrite did not shift the boundary (S3).
- `precedence:` overflow > held on the SAME drain — a drain past the cap
  whose writer never closes reports `drainOverflow`, not `drainHeld`.
- `precedence:` the fold's fixed order and its input — `heldPipes` takes the
  two REPORTED conditions and derives nothing from byte counts; returns
  `[stdout stderr]` in that order for both-held, the bare name for either
  alone, and nothing when a drain reports overflow rather than held.
- `precedence:` happens-before edge (b) — the parent must not read the
  output slices until every drain goroutine returns. Every probe above,
  including the paced and burst shapes that write to the slices while the
  parent is joining, ran clean under `-race`.
- `precedence:` MVV row 6(a), the burst shape — 60 KiB buffered, the parent
  arms the grace, the writer then delivers another 60 KiB and closes: 122880
  bytes recovered, `drainWhole`. A deadline of literally `now` would have
  short-circuited to n=0.
- `refusal:` `Detail` composition ORDER — the held-pipe reason LEADS and the
  stderr tail FOLLOWS (the contract fence, not S10's inverse wording):
  observed `"the child's stdout, stderr stayed held past the 500 ms drain
  bound (exited 0); close or redirect the helper's inherited stdio\nsome
  stderr noise\n"`. Probed with an EMPTY tail (no stray separator), with a
  tail whose bytes MIMIC the reason prefix (the reason still leads), and
  with the F5 non-pollable line present (reason, then F5 line, then tail).
- `refusal:` the reason's four required parts — the pipe(s), the bound
  ("500 ms drain bound"), the child's exit status, the remediation.
- `refusal:` `*accessor.HeldPipeError` declared in `internal/accessor`,
  exported, and matchable with `errors.As` THROUGH the `*ExecError` carrier.
- `refusal:` applied sense on the DIRECT write leg — a held-pipe error
  yields `execution_failure` with `Applied()` TRUE.
- `refusal:` the Phase-4 NEGATIVE CONTROL, the input the record names as
  the inverse defect: three other shapes reaching the SAME shared
  `err != nil` arm — a plain `*ExecError` with no wrapped held error, a bare
  non-`ExecError`, and an `*ExecError` wrapping a non-held error — each
  yields `Applied()` FALSE. The match is on the held-pipe sibling and not on
  `*ExecError` generally.
- `refusal:` `Applied()` is FALSE on read and gate even when the binding
  fails with a held-pipe error.
- `class:` deadline-first — `invokeRead` checks `bounded.Err()` for
  `context.DeadlineExceeded` BEFORE the binding error, and the `timeout`
  arm returns `readOutcome{class: ClassTimeout}` with a nil error, so
  `refusalOf` copies no `Detail`: a `timeout` refusal carries no Err/Detail
  (F6, S6).
- `bound:` MVV rows 1-4, the deadline shape — the `Setsid` escapee holds
  both pipes and the child sleeps past a 2 s deadline: the refusal is
  DELIVERED (never withheld), ctx reports `DeadlineExceeded`, elapsed
  2.607 s, inside `timeout + 2·WaitDelay`.
- `bound:` MVV row 5, the success shape — child writes a whole envelope and
  exits 0 while the escapee holds both pipes: refused in 857 ms (inside
  `2·WaitDelay` + tolerance), `Err` naming `stdout, stderr` and `exited 0`,
  and the envelope NOT parsed (`values=[]`).
- S10 single-pipe form — stdout reaches EOF cleanly while the escapee holds
  stderr alone: refused, `Held=[stderr]` as the bare name, and the stderr
  tail the grace exists to preserve survives into `Detail`.
- `whole:` a paced child with no grandchild is unaffected — the join timer
  starts only after `Wait` returns, so a slow-writing child that exits
  normally never reaches the grace path at all (1.89 s, whole envelope
  parsed).

### Phase 3a verdict

BLOCK — two reproduced divergences, FAIL-1 and FAIL-2, both rooted in the
`readBounded` liveness ceiling, which is a second absolute bound the record
does not admit (`bound:` names `WaitDelay` as "the ONE bound" and states
`DrainGrace` "is NOT a second bound"). Everything else probed was
conformant, including the Phase-4 negative control and the `-race`
happens-before edge.

Read AFTER writing, for the record: Phase 3b reached FAIL-1 and FAIL-2
independently as ADV-2 and ADV-3, from the Failure Modes rather than from
C1's labelled lines and with different witnesses. Two independent derivations
converging on the same root cause raises confidence that the ceiling — not
the witness shapes — is the defect. Phase 3b additionally found ADV-1 (a
held stderr laundered into `drainOverflow` by the `StdoutCap` read limit,
withholding the refusal entirely); this phase observed the same asymmetric
limit but wrongly reasoned it conformant under "overflow outranks held",
missing that the rank presupposes overflow REFUSES, which it does not on
stderr. ADV-1 stands.
