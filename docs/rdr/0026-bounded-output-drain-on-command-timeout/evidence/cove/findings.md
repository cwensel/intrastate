Model: claude-opus-5[1m]

# CoVe pre-lock lens — RDR 0026 (bounded output drain on command timeout)

Leaf run. Edges JSON already CONFIRMED all 45 `source-anchor` symbols textually exist;
this pass verifies BEHAVIOUR and sweeps the unanchored prose.

Spans widened past the RDR's scoped ranges, and why:
- `internal/accessor/executor.go:300-480` (whole `(*Executor).Write` body) — A8 and C1
  `refusal:` claim a class and an applied sense for the write path AND its read-back;
  neither is judgeable from the cited line numbers alone.
- `internal/cli/flow_exec.go:355-470` (`accessorFailure` + `withStderrTail` +
  `accessorFailureOf`) — A8 cites only `accessorFailureOf`; the Detail claim is decided
  one frame up, in `accessorFailure`.
- `internal/cli/cmdbind/cmdbind.go:556-737` (the three `spawn` callers' full bodies) —
  A1's "no fast path parses a held stream" needs the post-`err` code, not just the
  `if err != nil` line.
- `internal/accessor/model.go:300-330`, `internal/cli/cmdbind/procgroup_other.go` —
  applied-sense field and the non-Unix arm.

## Step 0 — Grounding sweep

Note: Decision Rationale carries `Ground-sweep: clean (14 anchors, 66 sub-facts
CONFIRMED)`. Noted; swept independently anyway. My sweep agrees on the anchor-existence
and most behaviour, and diverges on two behavioural claims (G6, G7 below).

| # | Claim (element) | Verdict | Cite / what I found |
| --- | --- | --- | --- |
| G1 | `spawn` holds the package's only `exec.CommandContext`, only `cmd.Wait()`, only `drains.Wait()`; three callers only (A1) | CONFIRMED | `internal/cli/cmdbind/cmdbind.go::spawn` at `:158`; `exec.CommandContext` `:205`; `cmd.Wait()` `:282`; `drains.Wait()` `:302`. `grep -n 'spawn('` returns exactly `:566` (`Reader.Read`), `:657` (`Gate.Gate`), `:722` (`(*Writer).Apply`) plus the definition. |
| G2 | Each caller does `inv, err := spawn(...)` then unconditionally `if err != nil { return }` BEFORE touching `inv` (A1) | CONFIRMED | `cmdbind.go:566-570` (`return nil, nil, err`), `:657-661` (`return "", "", err`), `:722-726` (`return err`). In all three the first `inv` field read (`len(inv.stdout)` / `inv.signaled`) is strictly after that guard. |
| G3 | Between `cmd.Wait()` and the join the only statement is `reapGroup(cmd.Process)`; no branch, no return (A1, C1 `stdin:`) | CONFIRMED | `cmdbind.go:282` `werr := cmd.Wait()` → `:297` `reapGroup(cmd.Process)` → `:302` `drains.Wait()`. Only comments intervene. Every `werr` switch arm is at `:304+`. |
| G4 | `readBounded` swallows the error and returns bytes on ANY error, so a deadline is indistinguishable from EOF today (Illustrative Code note; Infra Audit row) | CONFIRMED | `cmdbind.go::readBounded` — `b, err := io.ReadAll(io.LimitReader(r, int64(limit))); if err != nil { return b }`. No terminal condition reported out. Also confirms the "plumbing is part of this change" caveat is honest. |
| G5 | `WaitDelay = 500` at `:66`, applied at `:224`; its doc already reads "bounds the stdin write and the post-kill pipe drain" (A5, D-naming, Sibling-path check) | CONFIRMED | `cmdbind.go:63-66` const + doc; `:224` `cmd.WaitDelay = WaitDelay * time.Millisecond`. |
| G6 | A4: "on all four paths the `DeadlineExceeded` test is the first statement after the binding call and textually precedes the `if err != nil` arm" — including the **read-back** path | **REFUTED (as to the read-back leg)** | Read `executor.go:392-406`. The read-back leg does NOT test `applyCtx`/`readCtx` deadline itself; it tests `raw.class == ClassTimeout` — a value already computed inside `invokeRead`. The deadline-first rule therefore holds *transitively*, not textually, on that leg. More consequentially: the arm at `:398` maps `raw.class != ""` (i.e. a held-pipe `ClassExecutionFailure` from the read-back reader) to **`ClassReadBackIncomplete`**, not to `execution_failure`. See G7. |
| G7 | C1 `refusal:`: "a held pipe on the write path (the write invocation, **or its read-back after the write ran**) … class stays `execution_failure` and `Applied()` is true" | **REFUTED for the read-back leg** | `executor.go:398-406`: `if raw.class != "" { r := refusalOf(def, readTimeout, ClassReadBackIncomplete, raw.err); r.applied = true; … }`. A held pipe inside the read-back reader arrives as `raw.class == ClassExecutionFailure`, is re-labelled **`ClassReadBackIncomplete`**, and `applied` is *already* set true there. So on that leg the class does NOT stay `execution_failure`, and the `Applied()` derivation A8 proposes is not needed. |
| G8 | A8: `executor.go::Write` sets `r.applied = true` at six refusal sites by in-package field assignment | CONFIRMED | `executor.go:363`, `:375`, `:393`, `:401`, `:409`, `:472` — six `r.applied = true`. `applied bool` at `model.go:314`; `Applied()` at `model.go:323`. |
| G9 | A8: the `err != nil` arm of `Write` (`:368`) "today sets none" (no applied sense) | CONFIRMED | `executor.go:368-371`: `if err != nil { return WriteResult{Refusal: refusalOf(def, timeout, ClassExecutionFailure, err)} }` — comment "The command failed before mutating: this one did NOT occur." No `applied`. This is the site the RDR proposes to `errors.As` at, and it is reachable. |
| G10 | A8: `Applied()` "has zero non-test callers repo-wide" | CONFIRMED | `grep -rn 'Applied()' internal/ cmd/` returns only `model.go:323` (the definition) and `*_test.go` sites. No production caller. |
| G11 | A8 / C1 `refusal:`: `flow_exec.go::accessorFailureOf`'s `ClassExecutionFailure` arm "carries neither a phase check nor a `Detail` at all, so a held-pipe write would today render as a bare 'could not be executed'" | **PARTLY REFUTED** | The arm itself is bare — `flow_exec.go:409-411`: `return envErr(codeAccessorFailed, id, "the accessor \`"+id+"\` could not be executed")`, no phase, no `Detail`. But it is never reached bare: its caller `accessorFailure` (`flow_exec.go:370-372`) wraps it — `return withStderrTail(accessorFailureOf(refusal, at), refusal.Detail)` — and `withStderrTail` (`:378-390`) *does* land `refusal.Detail` in the one slot. So a held-pipe `execution_failure` would render the held-pipe reason + stderr tail today; what it would NOT render is the applied-but-unverified sentence. A8's "bare" wording overstates the gap; the *substance* of the gap (no `Applied()` key) stands. |
| G12 | C1 `class:` / F6: a `timeout` refusal carries no `Err`/`Detail` | CONFIRMED | `executor.go::refusalOf:74-87` — `Detail` is populated only via `errors.As(err, &ee)`. Every `ClassTimeout` mint passes `nil` for `err` (`:199`, `:231`, `:361`, `:394`). `Refusal` has no exported `Err` field (`model.go:308`: "set only where a binding returned an *ExecError"). |
| G13 | Sibling-path check: "searched `internal/` for another deadline/timer decision (`SetReadDeadline`, `time.After`, `time.NewTimer`, `WaitDelay`): none besides the `cmd.WaitDelay` assignment inside `cmdbind.go::spawn`" | CONFIRMED, with a caveat | `grep -rn 'SetReadDeadline\|time.After\|time.NewTimer\|WaitDelay' internal/ --exclude '*_test.go'` returns only `cmdbind.go:63,66,210,224`. Zero `SetReadDeadline`, zero `time.After`, zero `time.NewTimer` in non-test code. Caveat: the RDR's own search string omitted `context.WithTimeout`, which has three non-test sites (`executor.go:192`, `:224`, `:354`) — those are the *declared-timeout* decisions the RDR relies on, so no parallel signal is invented, but the check as written did not look for them. |
| G14 | The inverse the RDR did not check: does a sibling/adjacent path in this codebase ALREADY impose a bound on a post-exit drain, or already classify a held/unreadable output stream? | **searched, none exists** | The only bound family in non-test `internal/` is (a) `cmd.WaitDelay` at `cmdbind.go:224` and (b) `context.WithTimeout` on the *declared* timeout at `executor.go:192/224/354`. There is no other drain, no other `sync.WaitGroup` join over readers (`grep -rn 'sync.WaitGroup' internal/ --exclude '*_test.go'` → `cmdbind.go:271` only), and no held-pipe/unreadable-stream classification anywhere. C1's new rule has no existing sibling. |
| G15 | A7: `cmd.Stdin` is a `bytes.Reader` (Cmd-owned pipe), so `ErrWaitDelay` reaches the non-`ExitError` `default:` arm after `reapGroup` and the join | CONFIRMED | `cmdbind.go:207` `cmd.Stdin = bytes.NewReader(stdin)`; `:224` WaitDelay; `errors.As(werr, &ee)` at `:311`; `default: return inv, wrap(inv.stderr, werr)` at `:320-322`, which is after `:297` and `:302`. |
| G16 | Approach: "`spawn` hands `*os.File` ends to `cmd.Stdout`/`cmd.Stderr` precisely so the group can be released before the drains are joined (76c121b), which is why `Cmd`'s bound never reached them" | CONFIRMED | `cmdbind.go:226-262` — `os.Pipe()` pair, `cmd.Stdout = outW`, `cmd.Stderr = errW`, with the in-code comment stating exactly that rationale ("Owning them lets the group be released FIRST and the drains then run to a true EOF"). The "true EOF" promise the Problem Statement calls out is verbatim in that comment. |
| G17 | D-the-drains-stay-concurrent: "`spawn` already comments the verbose-tool deadlock this avoids" | CONFIRMED | `cmdbind.go:264-282` comment block: "BOTH channels drain CONCURRENTLY. Draining stdout to completion before touching stderr deadlocks a merely VERBOSE tool…". |
| G18 | C1 `precedence:` "After the join every read end is closed (deadline to unblock, close to release the fd)" — is there already a close? | CONFIRMED (already present, via defer) | `cmdbind.go:233` `defer func() { _ = outR.Close() }()` and `:239` `defer func() { _ = errR.Close() }()`. The reads are already closed on every return path today; C1's explicit `closeReads` is a re-statement, and the RDR does not note that the defers already exist (so the illustrative `closeReads(outR, errR)` would be a *double* close unless the defers are removed). |
| G19 | C1 `stdin:` "no arm of `spawn` returns between Wait and the join" | CONFIRMED | See G3. |
| G20 | Non-Unix arm exists and refuses pre-spawn (context for C1's platform reach) | CONFIRMED | `internal/cli/cmdbind/procgroup_other.go` — `procGroupSupported = false`, `killGroup` returns `errNoProcGroup`; `cmdbind.go:178-181` refuses on `Unsupported()` before spawn. C1 adds no platform clause; 0025:C4's refuse-list still governs. |

## Step 1+2 — Twelve verification questions, answered independently

**Q1 (source).** Does `readBounded` today report EOF-vs-deadline to its caller?
**A.** No. `internal/cli/cmdbind/cmdbind.go::readBounded` is
`b, err := io.ReadAll(io.LimitReader(r, int64(limit))); if err != nil { return b }` —
`[]byte` return only, error discarded. `heldPipes()` has nothing to read from today.
The RDR states this (Illustrative Code, Infra Audit row); confirmed accurate.

**Q2 (source).** Are the two read ends already closed on every return path from `spawn`,
such that C1's "close to release the fd" is new work?
**A.** They are already closed. `cmdbind.go:233` and `:239` install
`defer func() { _ = outR.Close() }()` / `defer func() { _ = errR.Close() }()`
immediately after each `os.Pipe()`. C1's `precedence:` clause and the illustrative
`closeReads(outR, errR)` describe an explicit close that duplicates the existing defers.
Not wrong, but the RDR nowhere says the defers exist or that they must be removed.

**Q3 (source).** On the write path's **read-back**, what class does a held-pipe
`execution_failure` from the reader actually produce?
**A.** `ClassReadBackIncomplete`, not `execution_failure`. `internal/accessor/executor.go:398-406`:
`if raw.class != "" { r := refusalOf(def, readTimeout, ClassReadBackIncomplete, raw.err); r.applied = true; r.Keys = compared; r.Expected = planned; return WriteResult{Refusal: r} }`.
C1 `refusal:` says the class "stays `execution_failure`" for "the write invocation, or its
read-back after the write ran". For the read-back leg that is false against the code.

**Q4 (source).** On that read-back leg, does `Applied()` already come out true without any
change from this RDR?
**A.** Yes. `executor.go:401` — `r.applied = true` is already on that arm. So the A8 gain
(deriving `Applied()` from a typed held-pipe error) is needed only on the **direct write
invocation** arm at `executor.go:368`, not on the read-back leg.

**Q5 (source).** Is `flow_exec.go`'s `ClassExecutionFailure` arm really rendered with no
`Detail`?
**A.** No — the arm builds none, but its caller adds one.
`flow_exec.go::accessorFailure:370-372` is `return withStderrTail(accessorFailureOf(refusal, at), refusal.Detail)`,
and `withStderrTail:378-390` assigns `ce.Detail = tail` when the CLIError has none. All
three call sites go through it: `flow_exec.go:285` (read), `:690` (gate),
`flow_state.go:372` (write). So a held-pipe `execution_failure` would render the held-pipe
reason today; only the applied-but-unverified *sentence* is missing.

**Q6 (source).** Does any sibling or adjacent path in `internal/` already bound a
post-exit output drain or classify a held stream?
**A.** No. `grep -rn 'SetReadDeadline\|time.After\|time.NewTimer\|WaitDelay\|sync.WaitGroup' internal/`
excluding tests yields only `cmdbind.go:63,66,210,224,271`. The three
`context.WithTimeout` sites (`executor.go:192,224,354`) bound the *declared* timeout, a
different decision. No held-pipe/unreadable-stream classifier exists anywhere. C1's rule
has no existing sibling — recorded as "searched, none exists".

**Q7 (source).** Is the deadline-first classification textually first on *all four* paths,
as A4 claims?
**A.** On three, yes: `executor.go:199` (`invokeRead`), `:231` (`Gate`), `:356-359`
(`Write`, via `appliedDeadline` computed immediately after `binding.Apply` and consumed at
`:359` before the `err != nil` arm at `:368`). On the fourth — the read-back inside
`Write` — the test at `:392` is `if raw.class == ClassTimeout`, a value already derived
inside `invokeRead`; there is no textual `DeadlineExceeded` test on that leg. The
deadline-first *effect* holds transitively; A4's textual claim does not.

**Q8 (source).** Would a fast path in any `spawn` caller parse a held stream before the
error is checked?
**A.** No. `cmdbind.go:566-570` (`Reader.Read`), `:657-661` (`Gate.Gate`), `:722-726`
(`(*Writer).Apply`) each `return` on `err != nil` before any `inv.stdout` / `inv.exitCode`
/ `inv.signaled` read. Confirmed by reading each caller's full body, not just the guard.

**Q9 (RDR-internal).** Does the RDR account for the fact that `spawn` already closes both
read ends via `defer`, when C1 mandates an explicit close after the join?
**A.** RDR is silent. `§technical-design` step 3 says "closes both read ends";
Illustrative Code shows `closeReads(outR, errR)`; the Existing Infrastructure Audit has no
row for pipe-close ownership. Nothing states whether the defers stay (double close, benign
but sloppy) or are removed.

**Q10 (RDR-internal).** Does the RDR name what happens if the drain goroutine returns
between the timer firing and `deadlineIn` being called — i.e. the benign race where the
join times out but the drain completes microseconds later?
**A.** RDR is effectively silent on the *ordering*, though the outcome is covered:
C1 `precedence:` says "each drain's FINAL read then reports the pipe's state" and S9 tests
the late-drain-no-holder case. But the illustrative shape calls `deadlineIn` on a pipe
whose reader may already have returned; `SetReadDeadline` on an fd whose read completed is
harmless, so this is a documentation gap, not a defect. Recorded, not raised as a finding.

**Q11 (RDR-internal).** Does C1 or any assumption state the class for a held pipe on a
**gate** invocation, and how the CLI renders it?
**A.** Partly. C1 `refusal:` says "On read and gate invocations `Applied()` is false as
today", and C1 `class:` gives the general timeout/execution_failure rule. Confirmed against
`executor.go::Gate:231-236` (deadline-first, then `ClassExecutionFailure` with the err) and
`flow_exec.go:690` (`accessorFailure(..., phaseGate)`). Consistent; no gap.

**Q12 (RDR-internal).** Does the RDR state the total bound for a **command write**, which
is two invocations?
**A.** Yes, explicitly. C1 `bound:`: "Per invocation, the COMMITTED return bound is
`timeout + 2·WaitDelay`; a command write is two invocations (write, then read-back), each
under its own bound." Consequences repeats the per-invocation form. No silence.

## Step 3 — Findings

1. **`0026:C1`** — class (a)/(c). C1 `refusal:` asserts that a held pipe on "the write
   invocation, **or its read-back after the write ran**" keeps class `execution_failure`;
   on the read-back leg the shipped code re-labels it **`ClassReadBackIncomplete`**.
   Cite: `internal/accessor/executor.go:398-406` —
   `if raw.class != "" { r := refusalOf(def, readTimeout, ClassReadBackIncomplete, raw.err); … }`.
   Either C1 must scope its class claim to the direct write invocation and name
   `read_back_incomplete` for the read-back leg, or it must state that it is *changing*
   that mapping (which would be a 0004:C13/C14 amendment this record does not declare).

2. **`0026:A8`** — class (a). A8's proposed `Applied()` derivation is unnecessary on the
   read-back leg: `executor.go:401` **already** sets `r.applied = true` on the
   `raw.class != ""` arm. The gain is real only at the direct-write arm
   (`executor.go:368-371`, which sets no applied sense — confirmed). A8's scope, and
   Testing Strategy S1's "write covering both its invocation and its read-back", overstate
   the work.

3. **`0026:A8`** — class (a). A8 states the CLI "would today render as a bare 'could not
   be executed'" because the `ClassExecutionFailure` arm "carries … no `Detail` at all".
   Refuted at the frame above: `internal/cli/flow_exec.go::accessorFailure:370-372` wraps
   every `accessorFailureOf` result in `withStderrTail(..., refusal.Detail)` (`:378-390`),
   which lands the Detail. The real gap is narrower and should be stated as such: the
   held-pipe reason DOES render; the applied-but-unverified sentence does not.

4. **`0026:A4`** — class (a). A4 claims "on all four paths the `DeadlineExceeded` test is
   the first statement after the binding call and textually precedes the `if err != nil`
   arm". True at `executor.go:199`, `:231`, `:356-359`; **not** true on the read-back leg,
   where `executor.go:392` tests `raw.class == ClassTimeout` — a class already derived
   inside `invokeRead` — with no `DeadlineExceeded` test of its own. The conclusion (the
   binding needs no timeout class) survives transitively; the stated method does not.

5. **`0026:C1`** — class (d), silence. C1 `precedence:` mandates "After the join every read
   end is closed (deadline to unblock, close to release the fd)" and the Illustrative Code
   adds `closeReads(outR, errR)`, but the record nowhere notes that `spawn` **already**
   closes both read ends via `defer` at `internal/cli/cmdbind/cmdbind.go:233` and `:239`.
   Silent on whether those defers stay (making the close double) or are removed. The
   Existing Infrastructure Audit, which enumerates every reused surface, has no row for
   read-end close ownership.

6. **`0026:D-selection-predicate`** — class (b), the inverse check, recorded as clean.
   Searched for an existing sibling that already bounds a post-exit drain or classifies a
   held output stream: `grep -rn 'SetReadDeadline\|time.After\|time.NewTimer\|WaitDelay\|sync.WaitGroup' internal/`
   (non-test) returns only `cmdbind.go:63,66,210,224,271`; the three
   `context.WithTimeout` sites (`executor.go:192,224,354`) bound the declared timeout, a
   different decision. **No sibling exists** — the RDR's Sibling-path check conclusion is
   correct, though its search string omitted `context.WithTimeout`.

### Not findings (checked, clean)
- A1's single-spawn-site and no-fast-path claims: fully CONFIRMED against all three
  callers' bodies (`cmdbind.go:566`, `:657`, `:722`).
- A7's stdin chain: CONFIRMED end-to-end (`:207`, `:224`, `:311`, `:320`).
- F6 / C1 `class:` "a `timeout` carries no Detail": CONFIRMED at
  `executor.go::refusalOf:82-85` plus every `ClassTimeout` mint passing `nil`.
- `readBounded`'s error-swallowing and the "plumbing is part of this change" caveat:
  CONFIRMED, honestly stated.
- Q12's per-invocation write bound: explicitly stated in C1 `bound:`; no silence.

---

## Dispositions (resolve half, iteration 1)

Grounding gate: all six re-checked against `main` by the resolving context before any
edit. No finding was dismissed as ungrounded.

| # | Anchor | Class | Disposition | Edit |
| --- | --- | --- | --- | --- |
| 1 | 0026:C1 | (a)/(c) | **fixed** | C1 `refusal:` no longer flattens the write's two legs: the direct invocation refuses `execution_failure` and needs the gain; the read-back refuses `read_back_incomplete` with `applied` already set. Confirmed at `executor.go`'s `raw.class != ""` arm and `flow_exec.go:423`. |
| 2 | 0026:A8 | (a) | **fixed** | A8 narrowed to ONE site (the `err != nil` write arm). A8's own Evidence already listed `:401` among the six sites that set `applied` — the record carried the refuting fact and still drew a two-site conclusion. S1 split into a gain assertion (direct) and a regression assertion (read-back). |
| 3 | 0026:A8 | (a) | **fixed** | "renders as a bare 'could not be executed'" corrected: `flow_exec.go::accessorFailure` wraps every mapping in `withStderrTail(…, refusal.Detail)`, so the held-pipe reason DOES render. The gain is the applied-sense sentence ahead of the tail, not the diagnostic. |
| 4 | 0026:A4 | (a) | **fixed** | Method claim scoped to the three DIRECT binding paths. The read-back leg tests `raw.class == ClassTimeout` — a class derived one frame down — so there is no second `DeadlineExceeded` test. Conclusion survives transitively; the wording no longer claims a uniform mechanism. A4 stays **Verified** (the conclusion was never disturbed). |
| 5 | 0026:C1 | (d) silence | **fixed** | New Infra-Audit row on read-end close ownership. Both read ends are ALREADY closed by defer at pipe creation; the Illustrative Code's `closeReads` is shape, not a new obligation. Double-close verified empirically to return `os.ErrClosed`, not panic. |
| 6 | 0026:D-selection-predicate | (b) | **fixed** (clean check, recorded gap) | The inverse check itself came back CLEAN — no sibling bounds a post-exit drain. But the recorded search string omitted `context.WithTimeout`; added, with the reason those sites are not siblings (they bound the declared timeout, the author's policy, not the CLI's mechanism cost). |

**Flag-as-you-go.** `A8` flipped **Verified → Pending** — its scope changed materially
(two sites → one) and C1 `refusal:` and S1 both depend on it. Stage 6 closes it.
`A4` stays Verified: only the method wording was narrowed, and the deadline-first
conclusion holds transitively.

**Net-new scope:** none. All six trace to the origin ledger (this pass's findings).

**Mini-checks fired:** all five — `authority`, `oracle`, `fidelity`, `disposition`,
`trace`. Tables written into Technical Design. The desk trace produced no
CONTRADICTION row; steps 6/9 are F6, already recorded as under-explained.
