Model: claude-opus-5[1m]

# CoVe pre-lock lens — RDR 0026, iteration 2 (DELTA pass)

Leaf run. Scope: only the seven anchors the iteration-1 resolving pass edited. The six
iteration-1 findings are settled and are NOT re-reviewed. Nothing outside the delta was
opened.

Spans read (source):
- `internal/accessor/executor.go:340-420` (`(*Executor).Write` write arm + read-back arm)
- `internal/accessor/executor.go:185-236` (`invokeRead` + `Gate` deadline-first arms)
- `internal/cli/flow_exec.go:355-475` (`accessorFailure`, `withStderrTail`,
  `accessorFailureOf`, `detailMayHaveApplied`)
- `internal/cli/cmdbind/cmdbind.go:232-242` (pipe creation + read-end defers),
  `:275-279`, `:345-362` (`readBounded`), `:383-392` (`wrap`), `:54-61` (caps)

Record read via projector only: `0026:C1`, `0026:A8`, `0026:A4`, `0026:A1`, `0026:S1`,
`0026:S2`, `0026:S3`, `0026:S6`, `0026:S7`, `0026:S9`, `0026:F6`, `0026:MVV`,
`0026:D-selection-predicate`, `0026:JC1`, `§conditional-mini-checks`,
`§existing-infrastructure-audit`, `§decision-rationale`, `§illustrative-code`,
`§technical-design`, `§consequences`, `§key-discoveries`.

---

## Delta item 1 — `0026:C1` `refusal:` two-leg split — **VERIFIED**

C1 now reads: "The DIRECT write invocation refuses `execution_failure`, and it is the leg
that needs the gain … The READ-BACK leg … re-classed `read_back_incomplete` with `applied`
already set … which `accessorFailureOf`'s `ClassReadBackIncomplete` arm already renders
with `detailMayHaveApplied`."

Every clause holds against source:

- Direct write arm, `internal/accessor/executor.go:367-370`:
  `if err != nil { return WriteResult{Refusal: refusalOf(def, timeout, ClassExecutionFailure, err)} }`
  — comment "The command failed before mutating: this one did NOT occur." No
  `r.applied` assignment on that arm. Class is `ClassExecutionFailure`. The "sets no
  applied sense today" claim is exact.
- Read-back arm, `executor.go:398-406`:
  `if raw.class != "" { r := refusalOf(def, readTimeout, ClassReadBackIncomplete, raw.err); r.applied = true; r.Keys = compared; r.Expected = planned; return WriteResult{Refusal: r} }`
  — class IS `ClassReadBackIncomplete`, `applied` IS already true. Exact.
- CLI rendering, `internal/cli/flow_exec.go:418-424`:
  `case accessor.ClassReadBackIncomplete: ce := envErr(codeReadBackIncomplete, id, "…did not complete"); ce.Detail = detailMayHaveApplied; return ce`
  — the read-back arm DOES already set the applied-sense Detail. Exact.
- `ClassExecutionFailure` arm, `flow_exec.go:410-412`:
  `return envErr(codeAccessorFailed, id, "the accessor \`"+id+"\` could not be executed")`
  — no phase check, no `Detail`. C1's "carries neither a phase check nor a `Detail`" is
  exact as a statement about THAT ARM (and C1 no longer claims the user sees a bare
  message — see item 2).

No residual flattening anywhere in the record: a sweep of `§technical-design`,
`§consequences`, `§key-discoveries` for `applied` / `read.back` / `both legs` found no
surviving text that asserts a uniform class across the two legs. `§consequences` line
"a command WRITE that applied its effect and exited 0 but left a held stdout is refused
`execution_failure`" is about the DIRECT leg (the child exited 0 and the write ran) and is
consistent, not a leftover.

## Delta item 2 — `0026:A8` narrowed to one site, "bare" corrected, Verified→Pending — **VERIFIED**

- One-site claim: correct. `executor.go:401` already sets `r.applied = true` on the
  read-back arm; `:368-370` sets none. The `errors.As` match therefore lands on exactly
  the `err != nil` arm. A8's Evidence now says "the re-key is one site (the `err != nil`
  write arm), not two."
- "bare" correction: correct. `flow_exec.go:369-371`
  `func accessorFailure(refusal accessor.Refusal, at phase) *clierr.CLIError { return withStderrTail(accessorFailureOf(refusal, at), refusal.Detail) }`
  and `withStderrTail:383-393` assigns `ce.Detail = tail` when the CLIError has none. A8
  now states "The held-pipe REASON still reaches the user one frame up … so what the gain
  adds is the applied-sense sentence ahead of that tail, not the diagnostic itself." That
  is what the code does.
- Line-number cites inside A8 spot-checked: `model.go:314` `applied bool`, `:323`
  `Applied()`; the six `r.applied = true` sites `:363, :375, :393, :401, :409, :472`;
  `detailMayHaveApplied` at `flow_exec.go:469`. All present.
- One imprecision, NOT raised as a finding: A8's Evidence writes "`ClassExecutionFailure`
  at `:410-412`" and, three clauses earlier, "`:410-414` (`ClassReadBackIncomplete`)".
  Against source, `ClassExecutionFailure` is at `:410-412` and `ClassReadBackIncomplete`
  is at `:418-424`. The earlier `:410-414` cite is off by one arm. This is a stale
  line-range in a *supporting* cite, not a claim: the substance (both arms exist,
  `ReadBackIncomplete` sets `detailMayHaveApplied`, `ExecutionFailure` does not) is
  correct and independently stated correctly later in the same Evidence block. Line
  numbers are not stable anchors and iteration 1's edits are not what introduced it (the
  cite predates the split). Recorded, not a finding.
- Status flip is right: A8's scope changed materially (two sites → one) and C1 `refusal:`
  and S1 both depend on it. Pending is the honest state; Stage 6 closes it.

## Delta item 3 — `0026:A4` scoped to "three DIRECT binding paths", left Verified — **VERIFIED**

The count "three" is right. The direct binding paths with a textual deadline test are:
- `executor.go:198-200` (`invokeRead`): "deadline makes both classes true, and `timeout`
  is the one reported" → `return readOutcome{class: ClassTimeout}` at `:200`, ahead of
  the `err` arm at `:203`.
- `executor.go:231-232` (`Gate`): `ClassTimeout` refusal at `:232`, ahead of the
  `ClassExecutionFailure` err arm at `:235`.
- `executor.go:354-359` (`Write`): `appliedDeadline := errors.Is(applyCtx.Err(), context.DeadlineExceeded)`
  at `:356`, consumed at `:358-364`, textually ahead of `if err != nil` at `:367`.

There is no fourth direct binding call: `grep -n 'context.WithTimeout' internal/ --exclude-tests`
returns exactly `executor.go:192, :224, :354` — three, matching. The read-back is not a
fourth `WithTimeout`; it reuses `ctx` and calls `e.invokeRead(ctx, reader, …)` at `:394`,
whose own `:192` bound is one of the three. So "three direct paths + one transitive leg"
is the correct decomposition, not four.

The transitive claim is exact: `executor.go:395` is `if raw.class == ClassTimeout`, a
class `invokeRead` derived at `:200` from ITS deadline check. A4's new text says exactly
this, including that the read-back returns the WRITE's refusal (`WriteResult{Refusal: …}`,
`applied = true`) rather than laundering it into a read result — confirmed at `:395-400`.

Leaving A4 **Verified** is defensible. A4's assertion is "the executor's deadline-first
classification turns the new refusal into `timeout` whenever the declared timeout elapsed,
on every path, so `spawn` needs no timeout class of its own." That conclusion was never
disturbed — only the stated MECHANISM was over-generalised, and the edit narrowed the
mechanism sentence without weakening the conclusion. The "If wrong" consequence (a
bound-expired drain misreported as `execution_failure`) is still ruled out on all four
paths. Contrast with A8, whose *scope* changed (two sites → one) and which correctly went
Pending. The asymmetry is principled.

## Delta item 4 — `0026:S1` gain/regression split — **VERIFIED**

S1 now asserts, on the DIRECT leg, `execution_failure` + "what proves the re-key landed"
because `ClassExecutionFailure` "carries no phase check and sets no `Detail` today"; and
on the READ-BACK leg, `read_back_incomplete` as "a REGRESSION assertion, not a proof of
the gain: that arm already sets `detailMayHaveApplied` and the executor already sets
`applied`."

Both halves match the code (item 1 cites). The split is also oracle-correct: the direct
row FAILS today (no `Applied()` key on that arm → no applied sense rendered), so it
discriminates the gain; the read-back row PASSES today, so it can only fail if the re-key
broke a correct path — which is precisely what a regression assertion is for. S1's own
text says this. No overstatement survives: S1 no longer claims "write covering both its
invocation and its read-back" as uniform work.

## Delta item 5 — new Infra-Audit row "Read-end close on every path" — **VERIFIED**

Row claims, and source:
- "spawn already defers both closes … set up at pipe creation" — `cmdbind.go:236`
  `defer func() { _ = outR.Close() }()` immediately after `os.Pipe()` at `:232`, and
  `:242` `defer func() { _ = errR.Close() }()` after `:237`. Both present, both at pipe
  creation, both unconditional.
- "closes at FUNCTION RETURN, which is after the join either way" — the join is inside
  `spawn`; deferred closes run at `spawn`'s return, strictly later. Correct.
- "`os.File.Close` returns `ErrClosed` on the second call rather than panicking" —
  verified empirically, not just by doc reading: a scratch program doing
  `r, w, _ := os.Pipe(); r.Close(); r.Close()` prints
  `first: <nil> second: close |0: file already closed isErrClosed: true`. Confirmed on
  this darwin/Go toolchain.

**No contradiction with C1 `precedence:`.** C1 says "After the join every read end is
closed (deadline to unblock, close to release the fd)". That is a statement about the
end-state at the join's completion, and the defers satisfy it — the row says exactly that
("is ALREADY satisfied by the defers"). The row does not license removing the defers, and
the Illustrative Code's `closeReads(outR, errR)` is explicitly demoted to "shape, not a
new obligation", consistent with the Illustrative Code block's own disclaimer
("`waitBounded`/`deadlineIn`/`closeReads`/`heldPipes`/`errHeldPipe` are names for the
shape, not the contract"). Decision `Reuse` and Spec Impact "Nothing in this change moves
the fds' release" are both true. The silence iteration-1 finding 5 named is closed and no
new silence replaces it.

## Delta item 6 — new `#### Conditional Mini-Checks` section (five tables) — **VERIFIED**

Placement: `§conditional-mini-checks` occupies lines 182-251, inside Technical Design,
after `§normative-contracts` (165-180) and before `§load-bearing-decisions` (253-288). As
specified.

**`authority` table** — five rows, all correct:
- "Refusal class … the EXECUTOR; the binding never classifies `timeout`" — matches C1
  `class:` and source (`spawn` returns `error`, never a class; `wrap:383-392` builds an
  `*accessor.ExecError`, no class field).
- **"Applied sense" row** — "direct arm (`err != nil`, gains it) and read-back arm
  (already sets it)"; "the two legs carry DIFFERENT classes"; canonical
  "`Refusal.applied`; `Applied()` is the only reader key (A8)". This agrees exactly with
  the corrected C1 and A8. This was the cell most at risk of carrying the old two-site
  version and it does not.
- "Held-vs-EOF terminal condition … none — `readBounded` today swallows it" — correct,
  `cmdbind.go:352-357` returns `b` on any `err`.
- "Stderr tail into `Detail` … `cmdbind.go::wrap` … `withStderrTail` composes into one
  slot, applied-sense first" — correct: `wrap:383-392` mints `ExecError{Detail, Err}`;
  `withStderrTail:383-393` is the single-slot composer and its doc comment states the
  applied-sense-first ordering verbatim.
- "Drain bound … `DrainGrace` is an OFFSET, not a rival bound … `WaitDelay`; no second
  knob" — matches C1 `bound:` and `cmdbind.go:63-66`.

**`oracle` table** — six rows, each pairing a scenario with a named failing control. Cross-
checked against `0026:MVV` rows 3/5/6 and `0026:S2`/`S3`/`S9`: every quoted expectation
(elapsed ≤ 3 s + `timeout` + empty `Detail`; `execution_failure` within `2·WaitDelay` with
`Err` naming "stdout, stderr"/"exited 0" and envelope not parsed; the deadline-of-`now`
control; 65536-vs-1 MiB; cap boundary; the byte-every-10-ms escaped writer) appears
verbatim in the cited element. No invented control.

**`fidelity` table** — six rows. Cap figures check out against source: `StdoutCap = 1 << 20`
(`cmdbind.go:57`), `StderrTailCap = 4 << 10` (`:61`), `tail` keeps the LAST
`StderrTailCap` bytes (`:364-368`), overflow detected at `StdoutCap+1` (`:275`, `:326`).
The 65536-prefix and concurrent-drain rows match `S7` and C1 `whole:`.

**`disposition` table** — eight rows. The two rows this pass had to scrutinise:
- "Held pipe on the DIRECT write | `execution_failure` + `Applied()` true | held-pipe
  reason | no read-back | LOUD once the re-key lands (A8)" — correct, and the
  "once the re-key lands" qualifier is the honest tense (it is NOT loud today).
- "Held pipe on the write's READ-BACK | `read_back_incomplete` + `applied` already true |
  read-back's own error | no confirmation | **LOUD today**" — correct: class from
  `executor.go:401`, `applied` from `:401`, and "LOUD today" from
  `flow_exec.go:423` `ce.Detail = detailMayHaveApplied`. Both rows agree with the
  corrected C1 and A8.
Other rows spot-checked: the `timeout`-with-no-`Detail` row matches F6 and 0025:C4; the
stdin row matches A7 and `spawn`'s non-`ExitError` arm; the `os.ErrNoDeadline` row matches
F5/S5 and the pipe-creation pollability check; the escapee row matches C1 `residue:`/F3.

**`trace` table** — ten steps. Steps map onto the MVV's numbered fixture steps and each
cites an in-record assertion plus a witness. Step 5's "500.1–508.2 ms, 40/40" and step
10's "2.013 s / 2.082 s" match A6 as quoted in `§key-discoveries` and the MVV's measured
basis. Step 8 ("read ends closed | C1 `precedence:`; Infra-Audit close-ownership row |
defers already close them; no double-close hazard") is consistent with the new Infra-Audit
row — it does not re-assert a new close obligation.

**The "No CONTRADICTION row" claim — TESTED, holds.** The section says: "No CONTRADICTION
row. Steps 6 and 9 are the two that a naive reading collides: the grace makes the final
read happen, and the class still flips to `timeout` when the deadline elapsed — so the
held-pipe REASON is computed and then not rendered. That is F6, recorded as
under-explained rather than contradictory."

I probed for a contradiction three ways and found none:
1. Steps 6/9 do not contradict: computing the held-pipe reason and then not rendering it
   is a lossy-but-consistent outcome, and it is already an owned, named failure mode —
   `0026:F6` states it directly ("a held pipe that also outran the deadline reports
   `timeout` with no `Err`/`Detail` … so the held-pipe reason is visible only on
   `execution_failure`"), and `0026:S6` tests it. Owned + tested = not a contradiction.
2. Steps 3 and 8 do not collide: step 3's "no arm returns between Wait and the join" (C1
   `stdin:`) and step 8's close-at-return are compatible because the defers fire strictly
   after the join, which the Infra-Audit row now states.
3. Steps 2 and 5/10 do not collide: step 2's deadline-first classification is the
   executor's, step 5/10's bound is the CLI mechanism's — C1 `bound:` explicitly reconciles
   them ("the deadline leg's unused allowance absorbs" the drain-join overrun), and
   `D-naming` / the Sibling-path check keep the two decisions separate.
The claim is not merely asserted; it survives the probe.

## Delta item 7 — Decision Rationale sibling-path search string — **VERIFIED**

New text: "searched `internal/` for another deadline/timer decision (`SetReadDeadline`,
`time.After`, `time.NewTimer`, `WaitDelay`, `context.WithTimeout`): none besides the
`cmd.WaitDelay` assignment inside `cmdbind.go::spawn` exists … The `context.WithTimeout`
sites in `executor.go` are not siblings: each bounds the DECLARED timeout for an
invocation, which is the author's policy, where this bound is the CLI's own post-exit
mechanism cost — a different decision, as D-naming states."

Re-ran the widened search: `grep -rn 'context.WithTimeout' internal/ --include='*.go'`
excluding tests returns exactly `executor.go:192`, `:224`, `:354` — the three declared-
timeout bounds (item 3). `SetReadDeadline`, `time.After`, `time.NewTimer` remain zero in
non-test `internal/`. `WaitDelay` remains `cmdbind.go:63,66,210,224` only.

The stated REASON is sound and is not special pleading: each `WithTimeout` bounds a
model-declared `timeout` (the accessor author's policy, read from the definition via
`reader.timeout()` / `def`), whereas `WaitDelay`/`DrainGrace` are CLI-owned constants that
are explicitly "not model-declarable" (C1 `bound:`) — a different owner and a different
input class, which is exactly the distinction `D-naming` draws. The conclusion the search
supports ("no parallel signal is invented") is unchanged and still correct.

---

## Findings

**No net-new findings.**

All seven edited anchors verify against source and are internally consistent with the
untouched parts of the record. No surviving text asserts the superseded two-site version;
no mini-check cell is factually wrong; the new Infra-Audit row closes iteration-1
finding 5 without contradicting C1 `precedence:`; the "No CONTRADICTION row" claim
survives an independent probe; and the two status decisions (A8 → Pending, A4 stays
Verified) are each defensible on the principle that scope changes demote and mechanism-
wording narrowings do not.

One sub-claim cite imprecision is **recorded, not raised**: A8's Evidence writes
`:410-414` for `ClassReadBackIncomplete` where source has that arm at `flow_exec.go:418-424`
(`:410-412` is the `ClassExecutionFailure` arm, which A8 also cites correctly). Line
ranges are not stable anchors, the substance is stated correctly elsewhere in the same
block, and the imprecision predates the iteration-1 edits — so it is out of delta scope
and does not meet the bar for a finding.
