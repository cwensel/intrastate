# REQ list — RDR 0026 bounded-output-drain-on-command-timeout

Phase 0 spec audit. Source: `docs/rdr/0026-bounded-output-drain-on-command-timeout.md`
(1074 lines; C=1, MVV=1, A=13, F=6, S=10, D=4, G=5, JC=1, ALT=3, BR=4).

Element ids are carried where a REQ derives from a labelled element, so a later
stage traces the REQ back to its contract. Quotes for `0026:C1` and `0026:MVV`
are exact bytes copied from the projector (`rdr inspect --select <id>`, verified
byte-identical to the record at lines 200-210 / 748-794); quotes outside those
fences are read from the record. Only line-wrapping is changed — no wording.

`0026:C1` is the record's ONE contract, carrying seven labelled lines
(`precedence:`, `refusal:`, `class:`, `stdin:`, `bound:`, `whole:`, `residue:`).
Each line carries several independent obligations, so C1 alone yields most of
the REQ set. The Conditional Mini-Check tables, Load-Bearing Decisions, Failure
Modes, Implementation Plan phases and Testing Strategy carry further testable
prose outside the fence — per PHASE 0, the projection narrows the read, it does
not replace it.

Section abbreviations:

- `NC` = Proposed Solution / Technical Design / Normative Contracts (fenced C1)
- `AP` = Proposed Solution / Approach
- `TD` = Proposed Solution / Technical Design (unfenced prose)
- `MC` = Technical Design / Conditional Mini-Checks (authority / oracle /
  fidelity / disposition / trace tables)
- `LBD` = Technical Design / Load-Bearing Decisions
- `IC` = Technical Design / Illustrative Code
- `EIA` = Implementation Plan / Existing Infrastructure Audit
- `FM` = Trade-offs / Failure Modes
- `CON` = Trade-offs / Consequences
- `MVV` = Implementation Plan / Minimum Viable Validation
- `IP` = Implementation Plan (Prerequisites, Phases 1-5)
- `TS` = Validation / Testing Strategy (S1-S10)
- `XC` = Finalization Gate / Cross-Cutting Concerns (`0026:G-cross-cutting`)
- `DEV` = artifacts/deviations.md (pre-seeded 7.1 entries D1, D2)

**Standing note on ownership.** This record amends `0025:C4` at exactly one
point (bound-vs-whole precedence) and restates `0025:F4`'s residue class. It
does NOT re-decide the `execution_failure` sub-reason carrier — that is
JDR 0003 §D1, which C1 `refusal:` cites as an instance. It does NOT edit 0025.
Classification remains the executor's, unchanged. Anything not named below is
out of scope for this REQ set.

---

## REQ-MVV — Minimum Viable Validation (`0026:MVV`)

- [REQ-MVV] "Fixture **FX-deadline-escape**, deadline shape, run through the real `cmdbind.Reader.Read` under `--allow-commands` with a declared `timeout = \"2s\"`" — (0026:MVV, §minimum-viable-validation)
- [REQ-MVV.1] "The child (a Go helper) spawns a grandchild with `SysProcAttr{Setsid: true}` that inherits stdout and stderr and sleeps 60 s; the child itself writes nothing and sleeps past the deadline." — (0026:MVV, §minimum-viable-validation)
- [REQ-MVV.2] "`flow resolve` invokes the reader; the context deadline expires at 2 s; `Cancel` kills the group; `Wait` returns; `reapGroup` signals the group; the drains are joined under the bound." — (0026:MVV, §minimum-viable-validation)
- [REQ-MVV.3] "Expected end-state: the read returns within `2 s + 2·500 ms` (Normative: elapsed ≤ 3 s, asserted with the C1 `bound:` tolerance), the executor classifies `timeout`, and the refusal's `Detail` is empty as 0025:C4 states for `timeout`." — (0026:MVV, §minimum-viable-validation)
- [REQ-MVV.3a] "Asserted at the `flow resolve` boundary, not only at the binding: the defect was reported there (\"still blocked after 15 s\"), so the row also asserts that `flow resolve` itself terminates and emits the class in its envelope. A binding-level pass with a hung CLI would not close the original report." — (0026:MVV, §minimum-viable-validation)
- [REQ-MVV.4] "The grandchild is still alive (asserted — the admitted residue) and the test kills it." — (0026:MVV, §minimum-viable-validation)
- [REQ-MVV.5] "Success shape, same fixture: the child writes a whole envelope and exits 0 while the grandchild holds the pipes — expected `execution_failure` within `2·500 ms` of the exit, `Err` naming \"stdout, stderr\" and \"exited 0\", and the envelope NOT parsed (no values returned)." — (0026:MVV, §minimum-viable-validation)
- [REQ-MVV.6] "Late-but-unheld control (Normative), TWO shapes — one guards each rejected mechanism, and shape (a) alone is not sufficient" — (0026:MVV, §minimum-viable-validation)
- [REQ-MVV.6a] "**Burst**: a child that writes 1 MiB to stdout in one `Write` and exits with no grandchild, run with the drain goroutine artificially delayed past the timer, returns the whole value — the final read under `DrainGrace` drains buffered bytes and sees EOF; no refusal. This shape guards C1 `precedence:` against a deadline of literally `now`, which short-circuits before the syscall and returns n=0 (A2, measured on both OSes)." — (0026:MVV, §minimum-viable-validation)
- [REQ-MVV.6b] "**Paced — the mechanism discriminator**: the same 1 MiB written in 32 chunks separated by idle gaps of 0.4·`DrainGrace` (20 ms), total tail duration ~12·`DrainGrace`, then exit with no grandchild. Every chunk is recovered and the drain ends on EOF, because the deadline is re-armed before each read and so bounds only the idle gap. An implementation that instead sets ONE absolute `now + DrainGrace` for the whole final drain truncates here — measured at 15.6% of the payload (163840 of 1048576 bytes, 5 reads, ~51 ms), terminal `os.ErrDeadlineExceeded` — and FAILS this row." — (0026:MVV, §minimum-viable-validation)
- [REQ-MVV.6c] "Shape (a) does NOT distinguish re-armed from single-absolute and must never be the only witness" — (0026:MVV, §minimum-viable-validation) — NEGATIVE REQ: row 6(a) may not be shipped as the sole witness for the re-arm mechanism.

---

## C1 `precedence:` — the drain precedence (`0026:C1`)

- [REQ-1] "precedence: a bounded wait outranks a whole read." — (0026:C1 `precedence:`, NC)
- [REQ-2] "Order in `spawn` is fixed: Wait (direct child reaped) → reapGroup (one kill(2) to -pgid, non-blocking) → drain join under ONE timer of WaitDelay covering BOTH read drains." — (0026:C1 `precedence:`, NC)
- [REQ-3] "When the timer fires the parent sets a read deadline of now + DrainGrace (50 ms) on its own read ends" — (0026:C1 `precedence:`, NC)
- [REQ-4] "each drain's FINAL read then reports the pipe's state, not the clock: bytes still buffered are delivered, EOF (every writer gone) is whole output as today, and a deadline error is a HELD pipe." — (0026:C1 `precedence:`, NC)
- [REQ-5] "The predicate is the drain's REPORTED terminal condition, not the byte count: a drain that ends on `os.ErrDeadlineExceeded` reports held however many bytes it delivered during the grace, and one that ends on EOF reports whole." — (0026:C1 `precedence:`, NC)
- [REQ-6] "\"Empty, with a writer the group signal could not reach\" describes the common case; it is not the test, because a writer trickling bytes slower than the grace is exactly the case that must refuse (S9's companion) and is never empty." — (0026:C1 `precedence:`, NC)
- [REQ-7] "Overflow outranks held on the same drain: a drain past `StdoutCap` refuses the existing overflow error whatever its terminal condition, since the bytes were read and the pipe's state is then not what the invocation turns on." — (0026:C1 `precedence:`, NC)
- [REQ-8] "The held fold takes the two drains' REPORTED conditions as its input — it derives nothing from the byte counts or from closed-over state — and returns the held pipe names in the fixed order stdout, stderr, so `len(held) != 0` is the consulted predicate and the single-pipe form is the bare name (S10, MVV row 5)." — (0026:C1 `precedence:`, NC)
- [REQ-9] "The grace is what MAKES that final read happen: Go checks the deadline BEFORE attempting the syscall (`internal/poll/fd_unix.go::(*FD).Read` calls `prepareRead` ahead of `syscall.Read`), so a deadline of NOW short-circuits with n=0 and reports a held pipe without ever looking at the pipe." — (0026:C1 `precedence:`, NC) — NEGATIVE REQ: a deadline of `now` is forbidden as the mechanism.
- [REQ-10] "The deadline is re-armed at `now + DrainGrace` BEFORE EACH read of the final drain, so the grace bounds the IDLE GAP between reads and never the size or total duration of the tail." — (0026:C1 `precedence:`, NC)
- [REQ-11] "That re-arm has NO site today and this clause creates one: `readBounded` (`cmdbind.go:352`) is `io.ReadAll(io.LimitReader(r, limit))` followed by `io.Copy(io.Discard, r)`, so `io.ReadAll` owns the loop and there is no per-read point to re-arm. Implementing this clause replaces that body with an explicit `for { SetReadDeadline(now+DrainGrace); Read }` loop — a rewrite of the function, not the signature change alone." — (0026:C1 `precedence:`, NC)
- [REQ-12] "The terminal condition is carried as ONE three-valued report per drain, not as a pair of booleans and not as a bare `error`: the three states are mutually exclusive and totally ordered here (overflow > held > whole), and a shape that cannot express that order — two independent booleans, or an `error` in which overflow is not an error — forces the caller to re-derive the precedence this clause fixes." — (0026:C1 `precedence:`, NC)
- [REQ-13] "`readBounded` returns that report alongside its bytes and the drain records it for `heldPipes()` to fold; the identifier and the Go spelling are the implementer's, the three-valued-and-ordered shape is not." — (0026:C1 `precedence:`, NC) — NEGATIVE REQ: no test may pin the identifier or Go spelling of the report type or of `heldPipes`.
- [REQ-14] "The rewrite must also preserve three behaviors the existing suite pins: the read stays bounded at `StdoutCap+1` so exactly-at-cap is a value and cap+1 is a failure (the overflow check at `cmdbind.go:326` reads `len(inv.stdout) > StdoutCap`), the discard-remainder exit must itself carry the deadline, and the pre-timer path must stay deadline-free." — (0026:C1 `precedence:`, NC)
- [REQ-15] "A re-implementation that shifts the cap boundary is a contract violation, not a refactor" — (0026:C1 `precedence:`, NC)
- [REQ-16] "a drain making progress keeps its deadline ahead of it and runs to EOF, and only a gap longer than the grace ends it as held." — (0026:C1 `precedence:`, NC)
- [REQ-17] "A single absolute deadline would instead bound the whole remaining tail against the clock, which turns a slow-arriving tail into a false held-pipe report on a pipe with no writer at all" — (0026:C1 `precedence:`, NC) — NEGATIVE REQ: a single absolute deadline for the whole final drain is forbidden.
- [REQ-18] "refuse on \"no EOF after a final bounded read\", not on elapsed time" — (0026:C1 `precedence:` quoting P-2, NC) — NEGATIVE REQ: no elapsed-time refusal predicate.
- [REQ-19] "Only a held pipe refuses." — (0026:C1 `precedence:`, NC)
- [REQ-20] "the parent's deadline-set on the read ends is the mark — `os.File`'s own poller is internally synchronised, so `SetReadDeadline` from the joining goroutine against a `Read` in flight is the sanctioned cross-goroutine call (A2 measured exactly this) and no separate `graceOn` flag is introduced" — (0026:C1 `precedence:` (a), NC) — NEGATIVE REQ: no `graceOn` flag.
- [REQ-21] "the parent still MUST NOT read `inv.stdout`/`inv.stderr` until every drain goroutine has returned, so the join remains a real join — the bound governs how long the parent waits before setting the deadline, never whether it waits for the goroutines to finish." — (0026:C1 `precedence:` (b), NC)
- [REQ-22] "A drain that ends on the deadline still runs `drains.Done()`, so the join completes; a design in which the parent proceeds past a live drain is forbidden here by name." — (0026:C1 `precedence:`, NC) — NEGATIVE REQ.
- [REQ-23] "The MVV cannot catch a violation of (b) — its fixture's child is silent, so the drain is blocked in its FIRST read and no partial write to the slice is in flight — so S3 and S7 run under `-race` and that is the check." — (0026:C1 `precedence:`, NC); the named check is "S3 + S7 under `-race -count=25` — is a Phase-1 exit condition, not a follow-up" — (0026:G-cross-cutting, XC)
- [REQ-24] "No read end outlives `spawn` (deadline to unblock, close to release the fd). The existing defers at pipe creation already satisfy this — they fire at function return, which is after the join but also after the `werr` switch and the overflow check, and that is the guarantee: no fd leaks, not that the fd is released at the join point. An explicit close before them would double-close, so none is added and the Illustrative Code draws none" — (0026:C1 `precedence:`, NC) — NEGATIVE REQ: no explicit close before the existing defers.

## C1 `refusal:` — the refusal carrier and applied sense (`0026:C1`)

- [REQ-25] "refusal:  a held pipe refuses through `*accessor.ExecError`, this record's instance of JDR 0003 §D1" — (0026:C1 `refusal:`, NC)
- [REQ-26] "`Detail` carries the held-pipe reason — the pipe(s) held (\"stdout\"/\"stderr\"/both), the bound, the direct child's exit status (exited N / killed by signal) and the remediation (\"close or redirect the helper's inherited stdio\") — composed AHEAD of the stderr tail collected up to the bound." — (0026:C1 `refusal:`, NC)
- [REQ-27] "That exit status is composed from fields `spawn` already populates from `werr` — `inv.exited` and `inv.signaled` beside `inv.exitCode` — so no new field is added to carry it" — (0026:C1 `refusal:`, NC) — NEGATIVE REQ: no new invocation field for exit status.
- [REQ-28] "`Err` wraps a typed held-pipe error carrying the pipes held and the exit status, which the executor matches with `errors.As` at the refusal site." — (0026:C1 `refusal:`, NC)
- [REQ-29] "The type is declared in `internal/accessor`, NOT in `cmdbind`: `cmdbind` imports `accessor` (`cmdbind.go:44`) and `accessor` imports no `cli` package (`go list -deps ./internal/accessor` returns `internal/resolve`, `internal/table` only), so a concrete `cmdbind` type named at an `errors.As` site inside `executor.go` is an import cycle that does not compile." — (0026:C1 `refusal:`, NC)
- [REQ-30] "Being matched from `executor.go` and populated from `cmdbind`, both outside its own package, the type is EXPORTED and its held-pipes and exit-status fields with it — an unexported type would be unnameable at the `errors.As` site" — (0026:C1 `refusal:`, NC)
- [REQ-31] "Applied sense: a held pipe on the write path after the direct child was reaped is a post-run refusal — `Applied()` is true, carried per JDR 0003 §D1 through the executor's write arm." — (0026:C1 `refusal:`, NC)
- [REQ-32] "The write's TWO legs already differ in class and this clause does not flatten them. The DIRECT write invocation refuses `execution_failure`, and it is the leg that needs the gain: `executor.go::Write`'s `err != nil` arm sets no applied sense today, and `flow_exec.go::accessorFailureOf`'s `ClassExecutionFailure` arm carries neither a phase check nor a `Detail`, so a held-pipe write slips past the \"may have been applied\" rendering (JDR 0003 §D1 rule 3 lands that gain here)." — (0026:C1 `refusal:`, NC)
- [REQ-33] "The READ-BACK leg after the write ran is already whole and gains nothing: a held pipe there returns the write's own refusal re-classed `read_back_incomplete` with `applied` already set (`executor.go`'s `raw.class != \"\"` arm), which `accessorFailureOf`'s `ClassReadBackIncomplete` arm already renders with `detailMayHaveApplied`. So one site changes, not two" — (0026:C1 `refusal:`, NC) — NEGATIVE REQ: the read-back arm must NOT be re-keyed.
- [REQ-34] "On read and gate invocations `Applied()` is false as today." — (0026:C1 `refusal:`, NC)
- [REQ-35] "The returned invocation carries no stdout; stdout read on a held drain is never parsed on ANY path — read envelope, gate verdict, write, write read-back (JDR 0003 §D1)." — (0026:C1 `refusal:`, NC)
- [REQ-36] "\"Carries no stdout\" is a CALLER obligation, not a structural guarantee, and is stated as one so a later caller cannot inherit it by accident: `spawn` builds `inv` before the `werr` switch (`cmdbind.go:304`) and both error arms return that populated value (`:322`, `:326`), so the field is readable on every error path. The obligation is that every caller checks `err` before touching `inv` — true of today's three callers by inspection (A1), and NOT enforced by any type or signature." — (0026:C1 `refusal:`, NC)
- [REQ-37] "Implementation carries this as a comment at the `inv` construction site naming the two error returns; making it structural (an error-only return, or zeroing `inv` on the refusal arms) is a `readBounded`-adjacent change this record does not take and 0028's command-reader read-back leg should revisit" — (0026:C1 `refusal:`, NC) — NEGATIVE REQ: do NOT make it structural (no error-only return, no zeroing `inv`).

## C1 `class:` — classification stays the executor's (`0026:C1`)

- [REQ-38] "class:    the class is the executor's, as today: ctx `DeadlineExceeded` ⇒ `timeout` (deadline-first rule; per 0025:C4 a `timeout` refusal carries no Err/Detail, so the held-pipe reason survives only on `execution_failure`), else the ExecError ⇒ `execution_failure`." — (0026:C1 `class:`, NC)
- [REQ-39] "The binding never classifies `timeout` itself" — (0026:C1 `class:`, NC) — NEGATIVE REQ.

## C1 `stdin:` — the stdin leg is unchanged (`0026:C1`)

- [REQ-40] "stdin:    unchanged and named: the child's stdin is a Cmd-owned pipe (`cmd.Stdin` is a `bytes.Reader`), bounded by Cmd's own WaitDelay; a non-reading child or escapee holding its read end ends in `exec.ErrWaitDelay` from Wait, which `spawn`'s non-ExitError arm refuses as `execution_failure` AFTER reapGroup and the join — no arm of `spawn` returns between Wait and the join" — (0026:C1 `stdin:`, NC)

## C1 `bound:` — the one bound and the committed total (`0026:C1`)

- [REQ-41] "bound:    WaitDelay (500 ms) is the ONE bound — the same constant, by design, for every \"the CLI waits past the child's exit\" case: Cmd's stdin write and post-Cancel wait, and the owned drain join; no second knob, not model-declarable." — (0026:C1 `bound:`, NC) — NEGATIVE REQ: no second bound, not model-declarable.
- [REQ-42] "DrainGrace (50 ms) is NOT a second bound: it answers a different question (how far ahead the read deadline must sit for the kernel to deliver already-buffered bytes), is never model-declarable, and cannot move the refuse/accept boundary — it is a mechanism constant inside the one bound, named here because every shipped constant of this family is named in its contract (0025:C4 does the same for StdoutCap and StderrTailCap)." — (0026:C1 `bound:`, NC)
- [REQ-43] "Cancel sends one signal and returns, so Cmd's timer starts at the deadline." — (0026:C1 `bound:`, NC)
- [REQ-44] "Per invocation, the COMMITTED return bound is timeout + 2·WaitDelay, asserted with a scheduling tolerance of +100 ms (the premortem's P-8 figure) — stated here once so the scenarios cite it rather than each carrying its own: a return past the total but inside the tolerance is a pass, and past the tolerance is a contract violation." — (0026:C1 `bound:`, NC)
- [REQ-45] "The bound is PER INVOCATION, and a command write runs up to THREE of them through the same reader, not two: the pre-write baseline read when `protectedKeys` returns non-empty (`executor.go:320`), the write itself (`:355`), and the read-back (`:397`). A write journey against an escaping helper therefore costs up to `3·(timeout + 2·WaitDelay)` end to end" — (0026:C1 `bound:`, NC)
- [REQ-46] "the baseline leg's refusal is DISCARDED (`:326` sets `baselineUnread` and continues), so that leg spends a full bound and yields no signal ... This record does not change that discard (it is 0004:C13's deliberate \"unverifiable, not unconstrained\" arm); it names the cost so the published total is not read as the journey's." — (0026:C1 `bound:`, NC) — NEGATIVE REQ: the baseline-leg discard is NOT changed.
- [REQ-47] "The two legs are NOT separately bounded and the contract does not claim they are: the drain join runs one full WaitDelay plus timer-fire and goroutine-scheduling latency (measured 500.1-508.2 ms in 40/40 samples, A6), which the deadline leg's unused allowance absorbs." — (0026:C1 `bound:`, NC) — NEGATIVE REQ: no per-leg assertion may be minted.
- [REQ-48] "a child that RESISTS termination spends leg (a)'s allowance and is the case where the total tightens — S8 tests it" — (0026:C1 `bound:`, NC)

## C1 `whole:` — the unchanged whole-output path and its antecedent (`0026:C1`)

- [REQ-49] "whole:    a drain that reaches EOF — before the timer, or on its final read after it — is whole, as today; the 1 MiB stdout cap and 4 KiB stderr tail are unchanged. A child (and group) that closes its pipes observes no change." — (0026:C1 `whole:`, NC)
- [REQ-50] "The completeness guarantee (A3) carries an antecedent, stated here because nothing structural enforces it: the two drains run CONCURRENTLY and are never stalled." — (0026:C1 `whole:`, NC)
- [REQ-51] "A child whose payload exceeds the ~64 KiB pipe buffer blocks in write(2) until drained, so a serialized or stalled drain converts a large-answer success into a deadline kill — measured: exactly one pipe buffer (65536 bytes) survives, as a correct prefix, where a prompt drain returns the whole 1 MiB." — (0026:C1 `whole:`, NC)
- [REQ-52] "this clause is what keeps a later refactor from silently voiding A3, since the prohibition is carried by the contract and S7, not by a build error" — (0026:C1 `whole:`, NC) — NEGATIVE REQ: the drains may not be serialized.

## C1 `residue:` — the restated 0025:F4 residue (`0026:C1`)

- [REQ-53] "residue:  0025:F4 restated — a process outside the child's process group may outlive the refusal; it keeps write ends whose reader is gone, so its next write fails with EPIPE (SIGPIPE under the default disposition), and whatever it wrote is never read." — (0026:C1 `residue:`, NC)
- [REQ-54] "Leakage of a process, or of THAT process's output, is admitted; a withheld refusal is not." — (0026:C1 `residue:`, NC)
- [REQ-55] "Diagnosis: the refusal's Detail names the held pipe(s); the leaked process is visible to `ps` until its next write" — (0026:C1 `residue:`, NC)

---

## Approach / Technical Design prose (outside the fence)

- [REQ-56] "the order stays as it is (`cmd.Wait()`, then `reapGroup`, then the drain join), and the join becomes bounded by the existing `cmdbind.go::WaitDelay` (500 ms), not by a second constant." — (AP)
- [REQ-57] "Output read on a bound-expired drain is never parsed on any path." — (AP)
- [REQ-58] "the refusal is never withheld, and the invocation returns within `timeout + 2·WaitDelay`." — (AP)
- [REQ-59] "`reapGroup` as today — release first, so the common backgrounded-grandchild case closes its ends before any timer matters." — (TD)
- [REQ-60] "On EOF within the bound: unchanged — whole output, the 1 MiB cap and 4 KiB tail apply as today." — (TD)
- [REQ-61] "only if a drain reported the deadline — returns `wrap(detail, errHeldPipe)` — `Detail` leading with the held pipe(s) and the child's exit status ahead of the stderr tail, `Err` the typed held-pipe error (JDR 0003 §D1)." — (TD)
- [REQ-62] "Nothing read from a held stdout is handed to `Reader.parse`, `Gate.Gate`'s verdict mapping, or the writer's read-back." — (TD)
- [REQ-63] "Pollability is checked once at pipe creation (a far-future deadline; `os.ErrNoDeadline` marks the pipe non-pollable) so the fallback is decided before any child exists (F5)." — (TD)
- [REQ-64] "Classification is unchanged and lives where it lives today: the executor reads the bounded context's `DeadlineExceeded` first (⇒ `timeout`), then the binding's error (⇒ `execution_failure`). `spawn` therefore never has to know whether the bound expired because of a deadline or because a success-path child leaked a writer." — (TD)

## Illustrative Code — shape-only, with two binding negatives

- [REQ-65] "Illustrative: `heldPipes`/`errHeldPipe` are names for the shape, not the contract; the timings in the MVV are Normative." — (IC) — NEGATIVE REQ.
- [REQ-66] "The block deliberately draws NO `graceOn` and NO `closeReads` call ... the existing defers already close both read ends, so an explicit close would double-close" — (IC) — NEGATIVE REQ.
- [REQ-67] "the re-arm lives INSIDE that read loop, not in the parent. A one-shot `SetReadDeadline` from the parent before `drains.Wait()` would bound the whole remaining tail against the clock" — (IC) — NEGATIVE REQ. See ASSUMPTION-1 for the reconciliation with REQ-3.

## Existing Infrastructure Audit

- [REQ-68] "C1 lands once; no per-binding wait. Its signature is unchanged by this record: it returns the `invocation` by VALUE (not a pointer), which is what makes C1 `refusal:`'s \"the field is readable on every error path\" hold without a nil check" — (EIA, `spawn` row) — NEGATIVE REQ: `spawn`'s signature is unchanged.
- [REQ-69] "`exited`/`signaled` already carry C1 `refusal:`'s \"exited N / killed by signal\"; F5's pollability result is recorded here as one more field, which is how it reaches the `Detail` line S5 asserts on" — (EIA, `invocation` row)
- [REQ-70] "C1 `bound:` — no second bound" — (EIA, drain-bound-constant row; Decision: Reuse `cmdbind.go::WaitDelay`)

## Load-Bearing Decisions

- [REQ-71] "the bound is `cmdbind.go::WaitDelay`, reused; rejected: a new `DrainBound` constant (a parallel knob for one decision), and a model-declarable drain bound" — (LBD `0026:D-naming`) — NEGATIVE REQ.
- [REQ-72] "`DrainGrace` (50 ms) ships beside `WaitDelay` ... It is not model-declarable, moves no refuse/accept boundary, and invents no second signal ... it is an offset on the single existing timer's expiry." — (LBD `0026:D-the-drain-grace-is-a-mechanism-constant-not-the-rejected-knob`)
- [REQ-73] "50 ms is chosen against the measured scheduling tail (A5 saw drain-goroutine latency reach ~28 ms under 4x-core saturation), so the grace clears it with margin while costing 10% of one bound." — (LBD `0026:D-the-drain-grace-...`)
- [REQ-74] "the tail reaches the user on the `execution_failure` leg (MVV row 5, S10) and NOT when the declared timeout also elapsed" — (LBD `0026:D-the-drain-grace-...`)
- [REQ-75] "**The drains stay concurrent** — A3's no-loss guarantee holds only while they do ... Recorded here, and pinned by S7, so the prohibition is load-bearing rather than decorative." — (LBD `0026:D-the-drains-stay-concurrent`)
- [REQ-76] "when EOF and the bound both qualify: EOF wins if it arrives first (whole output, as today); the bound wins otherwise and the invocation refuses." — (LBD `0026:D-selection-predicate`)
- [REQ-77] "When `timeout` and `execution_failure` both qualify: `timeout` (the executor's existing deadline-first rule, not a new tie-break)." — (LBD `0026:D-selection-predicate`)
- [REQ-78] "When overflow and held both qualify on one drain ... overflow wins, and the existing \"stdout exceeded the 1 MiB bound\" refusal after the join (`cmdbind.go:326-329`) stands unchanged." — (LBD `0026:D-selection-predicate`)
- [REQ-79] "When the `werr` non-ExitError arm (`exec.ErrWaitDelay`, C1 `stdin:`) and a drain condition both qualify, the drain conditions are selected FIRST — overflow, then held — and the non-ExitError arm refuses only if neither fired." — (LBD `0026:D-selection-predicate`)

## Conditional Mini-Checks — `authority` census

- [REQ-80] "Refusal class (`timeout` vs `execution_failure`)" canonical: "the EXECUTOR; the binding never classifies `timeout` (C1 `class:`)" — (MC `authority`)
- [REQ-81] "Applied sense on a write" canonical: "`Refusal.applied`; `Applied()` is the only reader key (A8)" — (MC `authority`)
- [REQ-82] "Held-vs-EOF-vs-overflow terminal condition" canonical: "the drain's reported condition; NOT the clock, NOT the byte count (C1 `precedence:`); overflow outranks held (D-selection-predicate), which is why the report is one ordered value and not two booleans" — (MC `authority`)
- [REQ-83] "Stderr tail into `Detail`" canonical: "`refusal.Detail`, ordered per 0025:C4"; writer `cmdbind.go::wrap`, "`withStderrTail` composes into one slot", "one slot, applied-sense first" — (MC `authority`)
- [REQ-84] "Held-pipe error type" canonical: "`internal/accessor` — `cmdbind` imports `accessor`, never the reverse (C1 `refusal:`)"; "one declaration, one match site" — (MC `authority`)
- [REQ-85] "Drain bound" canonical: "`WaitDelay`; no second knob (D-naming)"; "`DrainGrace` is an OFFSET, not a rival bound" — (MC `authority`)

## Conditional Mini-Checks — `fidelity` invariants

- [REQ-86] "Whole read, child closes pipes" → "byte-for-byte unchanged from today (C1 `whole:`, S3)"; lossy exemption "none" — (MC `fidelity`)
- [REQ-87] "stdout at the 1 MiB cap" → "exactly the cap is a value; one byte more is a detected overflow"; exemption "the discarded remainder — bounded on what is KEPT, as today" — (MC `fidelity`)
- [REQ-88] "stderr tail" → "LAST 4 KiB preserved" — (MC `fidelity`)
- [REQ-89] "Drain past the ~64 KiB pipe buffer, prompt" → "whole 1 MiB recovered (S7)"; exemption "none while the drains stay concurrent" — (MC `fidelity`)
- [REQ-90] "Final drain after the timer, writer gone" → "whole remaining tail, however many reads it takes — the grace bounds the idle gap between reads, not the tail (C1 `precedence:`, A9)"; "none; a tail that stalls longer than the grace is held, not truncated" — (MC `fidelity`)
- [REQ-91] "Drain past the pipe buffer, STALLED" → "exactly one pipe buffer (65536 bytes) as a **correct prefix**" — (MC `fidelity`)
- [REQ-92] "Held pipe" → "stdout is never parsed on ANY path"; exemption "the whole stream — refused, not truncated-then-read" — (MC `fidelity`)

## Conditional Mini-Checks — `disposition` (input class → outcome)

- [REQ-93] "Child closes pipes, exits 0" → "success", "none", "envelope parsed", "loud (the value)" — (MC `disposition`)
- [REQ-94] "Held pipe, deadline NOT elapsed" → "`execution_failure`", "the held-pipe error (accessor-side, C1 `refusal:`) in `Err`, reason + tail in `Detail`", "no stdout", "LOUD" — (MC `disposition`)
- [REQ-95] "Held pipe, deadline ALSO elapsed" → "`timeout`", "none — 0025:C4 gives `timeout` no Err/Detail", "no stdout", "loud in class, **quiet in reason** (F6, S6)" — (MC `disposition`)
- [REQ-96] "Held pipe on the DIRECT write" → "`execution_failure` + `Applied()` true", "held-pipe reason", "no read-back", "LOUD once the re-key lands (A8)" — (MC `disposition`)
- [REQ-97] "Held pipe on the write's READ-BACK" → "`read_back_incomplete` + `applied` already true", "read-back's own error", "no confirmation", "LOUD today" — (MC `disposition`)
- [REQ-98] "Non-reading child, stdin past the buffer" → "`execution_failure`", "`exec.ErrWaitDelay` via `spawn`'s non-ExitError arm", "none", "loud (S4, A7)" — (MC `disposition`)
- [REQ-99] "Non-pollable pipe (`os.ErrNoDeadline`)" → "unbounded join — the old hang", "recorded on the invocation at the check", "a `Detail` line naming the host condition", "**loud at the check, by design** — as a `Detail` line, not a log line; `cmdbind` has no logger (F5, S5)" — (MC `disposition`)
- [REQ-100] "Escapee outside the group" → "not reaped", "EPIPE on its next write", "none", "**SILENT — admitted residue** (C1 `residue:`, F3)" — (MC `disposition`)

## Conditional Mini-Checks — `oracle` (named failing controls)

- [REQ-101] MVV row 5 negative control: "a build that parses held stdout returns values and fails the not-parsed assertion" — (MC `oracle`)
- [REQ-102] MVV row 6 negative control: "**this row IS the control, and it takes BOTH shapes** — 6(a) fails under a deadline of `now` (A2) but passes under either grace mechanism, so it cannot discriminate them; 6(b) is what fails the single-absolute variant" — (MC `oracle`)
- [REQ-103] S10 negative control: "stdout-held-alone: the same refusal with the other pipe named, and no tail" — (MC `oracle`)
- [REQ-104] S3 negative control: "one byte past the cap must still refuse overflow" — (MC `oracle`)
- [REQ-105] S9 negative control: "companion: escaped writer emitting a byte every 10 ms is REFUSED" — (MC `oracle`)
- [REQ-106] S2 negative control: "S7's stalled drain recovers exactly 65536 bytes where a prompt drain returns 1 MiB" — (MC `oracle`)

## Failure Modes

- [REQ-107] "**F5** ... It runs in `spawn` at pipe creation: a probe `SetReadDeadline` on each read end, immediately cleared, whose error is the signal. `cmdbind` has no logger and gains none (its imports carry no `log`/`slog`), so \"logged\" is not a log line: the condition is carried on the invocation and surfaces as a `Detail` line on any refusal that invocation produces, which is what a test asserts on. On that host the join stays unbounded" — (FM `0026:F5`) — NEGATIVE REQ: `cmdbind` gains no logger.
- [REQ-108] "**F6** Accepted trade-off, not a gap: a held pipe that also outran the deadline reports `timeout` with no `Err`/`Detail` (0025:C4's rule for the class), so the held-pipe reason is visible only on `execution_failure`." — (FM `0026:F6`) — NEGATIVE REQ: 0025:C4's `detail:` rule is NOT amended here.
- [REQ-109] "**F4** Silent risk: if A3 is wrong, bytes the direct child wrote can be lost at the bound; the refusal still fires, so the loss is bounded to a refused invocation, never a parsed one." — (FM `0026:F4`)

## Implementation Plan — phases

- [REQ-110] "A `Setsid: true` Go test helper exists in `internal/cli/cmdbind` so the escape is real on both OSes (Background: a shell fixture does not escape on macOS)" — (IP Prerequisites)
- [REQ-111] "Declare the held-pipe error type in `internal/accessor` beside `ExecError` ... Bound the drain join in `spawn` by `WaitDelay`; on expiry the parent re-arms a per-read deadline of `now + DrainGrace` on its own read ends (mechanism from the A2 spike), joins the drains, and refuses through `wrap` with the held pipe(s) named. `readBounded` gains a reported terminal condition ... No other line of `spawn` moves." — (IP Phase 1) — NEGATIVE REQ: no other line of `spawn` moves.
- [REQ-112] "Re-point `cmdbind.go`'s C4 comments (\"necessary and sufficient\", \"run to a true EOF\", `WaitDelay`'s doc) at 0026:C1; 0025 itself is not edited." — (IP Phase 2) — NEGATIVE REQ: 0025 is not edited.
- [REQ-113] "Add FX-deadline-escape (both shapes) beside the existing grandchild-leak and cancel oracles; re-run the in-group grandchild oracle under the bound (A5) so the prior point-fixes stay proven." — (IP Phase 3)
- [REQ-114] "match the held-pipe error in `executor.go::Write`'s `err != nil` arm ... and set the applied sense there; then key `flow_exec.go::accessorFailureOf`'s `ClassExecutionFailure` arm (`:410-412`) on `Applied()`." — (IP Phase 4)
- [REQ-115] "**Not descopable.** ... Phases 1-3 without Phase 4 are therefore not a shippable subset; if Phase 4 must slip, the bound ships with the write path excluded (reads and gates only) rather than with an unsensed write." — (IP Phase 4)
- [REQ-116] "**Negative control (required, not optional).** The `errors.As` match lands on the SHARED `err != nil` arm at `:368-371`, which every other pre-mutation failure shape also reaches ... Only the held-pipe error may set the applied sense there." — (IP Phase 4) — NEGATIVE REQ: no broadened `*accessor.ExecError` match.
- [REQ-117] "S1 gains a negative row: a binding whose command cannot start refuses `execution_failure` with `Applied()` FALSE and no `detailMayHaveApplied` in the envelope." — (IP Phase 4)
- [REQ-118] "SURVEY, then state. The survey runs FIRST ... run the existing reader/gate/write fixtures and the repo's own declared bindings under the new refusal, and record whether any binding that passed before now refuses. Acceptance for this half is the recorded result, empty or not ... A NON-empty result is a Stage-6 route-back, not a doc note" — (IP Phase 5)
- [REQ-119] "Then state the restated residue in `docs/cli-output-contract.md`, beside the command-entry `timeout` text authors already read — the \"a process may leak and its output is never read\" consequence, and the remediation the refusal names (close or redirect the helper's inherited stdio). Acceptance: that doc names the held-pipe refusal and the residue" — (IP Phase 5)

## Testing Strategy (S1-S10)

- [REQ-120] "**Scenario**: Escaped grandchild holds the pipes, on the read, gate and write verbs (write covering both its invocation and its read-back) ... **Expected**: each refuses within `timeout + 2·WaitDelay`; no envelope parsed, no gate verdict derived, no write confirmed on held output; on BOTH write legs `Applied()` is true AND the CLI renders the \"may have been applied\" text" — (TS `0026:S1`)
- [REQ-121] "The rendered half asserts on the envelope's `detail` field carrying `detailMayHaveApplied` (`flow_exec.go:423`'s existing constant, per docs/cli-output-contract.md's error envelope) — the same observable the read-back leg already renders, so the two legs are asserted the same way and neither needs a new golden." — (TS `0026:S1`) — NEGATIVE REQ: no new golden for the rendered half.
- [REQ-122] "The READ-BACK leg refuses `read_back_incomplete` and is a REGRESSION assertion, not a proof of the gain" — (TS `0026:S1`)
- [REQ-123] "The existing in-group grandchild-leak and cancel oracles re-run under the bounded join, N times under `-race` and CPU load (A5). **Expected**: unchanged results, elapsed far below the bound" — (TS `0026:S2`)
- [REQ-124] "Ordinary paths — child closes its pipes; output at the 1 MiB stdout cap and 4 KiB stderr tail boundaries. **Expected**: byte-for-byte unchanged from today (C1 `whole:`). \"Today\" is the existing `cmdbind` suite passing unmodified — no new golden is minted for this row ... Assert both boundary cases explicitly here rather than relying on the suite's incidental coverage; relaxing an existing expectation to make the rewritten loop pass is the failure this row names." — (TS `0026:S3`) — NEGATIVE REQ: no existing expectation may be relaxed.
- [REQ-125] "Non-reading child with a stdin payload past the pipe buffer (A7). **Expected**: `exec.ErrWaitDelay` from `Wait` reaches `spawn`'s non-`ExitError` arm and refuses `execution_failure` after `reapGroup` and the join." — (TS `0026:S4`)
- [REQ-126] "Non-pollable pipe — the creation-time pollability check returns `os.ErrNoDeadline` (F5), simulated at the check by injecting a read end that refuses a deadline (`os.NewFile` over a dup'd pipe fd, the shape F5 confirms produces the error). **Expected**: the condition is recorded at the check, not discovered at the join, and is observable as the F5 `Detail` line on that invocation's refusal ... The fallback is the documented unbounded path." — (TS `0026:S5`)
- [REQ-127] "Held stdout with the deadline also elapsed (F6). **Expected**: classified `timeout` with empty `Detail` per 0025:C4" — (TS `0026:S6`)
- [REQ-128] "payload past the ~64 KiB pipe buffer (1 MiB), drain prompt vs artificially stalled — the stall is ONE seam, shared by this row, S9 and MVV row 6 so they cannot build two: a test-only delay hook on the drain goroutine's start ... one unexported package-level hook set by the test and nil in production (not a build tag, not a sleep in production code). S7 is a standing regression trap, so the seam outlives the fixture" — (TS `0026:S7`) — NEGATIVE REQ: one seam only; not a build tag; no production sleep.
- [REQ-129] "**Expected**: prompt drain returns the whole 1 MiB; a stalled drain leaves the child blocked in `write(2)`, killed by the group signal, with exactly one pipe buffer (65536 bytes) recovered as a correct prefix." — (TS `0026:S7`)
- [REQ-130] "a child that RESISTS termination (ignores SIGTERM, survives to its `WaitDelay`) while a grandchild holds the pipes ... **Expected**: the invocation still returns within the C1 `bound:` total plus its stated scheduling tolerance" — (TS `0026:S8`)
- [REQ-131] "a LATE drain with no holder — the drain goroutine stalled past the timer with bytes buffered and every writer gone (premortem P-2). **Expected**: accepted, whole output, no refusal — the final read under `DrainGrace` delivers the buffered bytes and then sees EOF. Companion: an escaped writer emitting a byte every 10 ms forever is REFUSED (never blocked, never EOF)." — (TS `0026:S9`)
- [REQ-132] "stdout reaches EOF cleanly while STDERR alone is held by the escapee ... **Expected**: refused; `Detail` names `stderr` alone — the single-pipe form is the bare pipe name, and the both-pipes form is the two names comma-separated in the fixed order stdout, stderr (MVV row 5's `\"stdout, stderr\"`) — and the stderr tail collected up to the bound survives into `Detail` ahead of the held-pipe reason." — (TS `0026:S10`) — see QUESTION-1 on the `Detail` ordering.
- [REQ-133] "Coverage goal: the FX-deadline-escape helper is exercised on darwin and linux in CI, so the escape is real on both (Background)." — (TS)

## Cross-Cutting Concerns (`0026:G-cross-cutting`)

- [REQ-134] "A11 is the open record, and its named check — S3 + S7 under `-race -count=25` — is a Phase-1 exit condition, not a follow-up." — (XC, Concurrency model)
- [REQ-135] "Unix-only process-group syscalls stay behind the build tags established by intrastate#1drn (`procgroup_unix.go`); Windows continues to take 0025:C4's platform refusal. A12 confirms the drain-start stall seam this RDR adds needs **no third build tag** — it is portable Go in `cmdbind`." — (XC, Build tool compatibility) — NEGATIVE REQ.
- [REQ-136] "the existing 1 MiB stdout cap and 4 KiB stderr tail are unchanged and remain the only bound on retained output ... Output read on a bound-expired drain is never parsed on any path, so a partial buffer is never retained as a value." — (XC, Memory management)
- [REQ-137] "No opt-out flag ships and none is proposed — a per-binding \"accept partial output\" knob is ALT1 by another name, rejected for ALT1's reason; the escape hatch is redirection at the binding's command, which costs no contract." — (XC, Incremental adoption; also CON) — NEGATIVE REQ.

## Pre-seeded deviations (7.1 cluster-reconcile; binding on REQ wording)

- [REQ-138] (D1, TEST-FIXTURE, OPEN) "A unit test on `accessorFailureOf`: (write phase, `ClassExecutionFailure`, `Applied()` false) → no `detailMayHaveApplied`; (read phase, refusal wrapping the typed held-pipe error) → `Applied()` false. The `errors.As` for the held-pipe error lives in `Executor.Write`'s `err != nil` arm, never in the shared `refusalOf` — `0026:A8`'s \"either placement compiles\" is the trap." — (DEV D1) — NEGATIVE REQ: the `errors.As` may NOT live in `refusalOf`.
- [REQ-139] (D2, IMPL-DECISION, decided at JDR 0003 §D4 (b)) "a distinct exit-3 CLI code for the applied write refusal (spelling non-normative, e.g. `flow-write-failed-applied`), keyed on `Applied()` in `flow_exec.go::accessorFailureOf`'s write arm; `detail` text unchanged; not-applied refusals keep `flow-accessor-failed`. JDR 0001 §D10's table gains the code by citation repair." — (DEV D2)
- [REQ-140] (D2 check) "S1's rendered half asserts the new code AND the prose on the direct write leg; D1's negative witness asserts `flow-accessor-failed` with no prose on the `Applied()`-false write refusal." — (DEV D2) — refines REQ-120/REQ-121: the rendered assertion is the new exit-3 code PLUS `detailMayHaveApplied`, not the prose alone.

---

## ASSUMPTIONS

- **ASSUMPTION-1 (REQ-3 / REQ-10 / REQ-67 — who arms the deadline).** C1
  `precedence:` says "the parent sets a read deadline of now + DrainGrace ... on
  its own read ends" and Phase 1 repeats it, while the Illustrative Code says
  "the re-arm is HERE, per read, never in the parent" and draws
  `drains.Wait() // bounded by the drains' own per-read deadlines; the parent
  sets none`. Read as one mechanism in two moments, not two rivals: the PARENT
  sets the FIRST deadline at timer-fire (that set is also C1 `precedence:` (a)'s
  cross-goroutine happens-before mark, which only the parent can supply), and the
  DRAIN then re-arms `now + DrainGrace` before each subsequent read of its final
  drain. This is the only reading under which REQ-14's "the pre-timer path must
  stay deadline-free" holds (a drain arming its own deadline from the first read
  would deadline the pre-timer path) AND REQ-10's per-read re-arm holds. It is
  also the only reading under which MVV row 6(b) discriminates. Implementers
  therefore need a per-drain signal that the timer fired; per REQ-20 that signal
  is the deadline itself (the read returning `os.ErrDeadlineExceeded`), NOT a
  `graceOn` flag.
- **ASSUMPTION-2 (REQ-7 / REQ-78 — where overflow is detected).** Overflow
  outranks held, and the existing overflow refusal at `cmdbind.go:326` (which
  reads `len(inv.stdout) > StdoutCap` AFTER the join) "stands unchanged". So the
  three-valued report's `overflow` state and the existing length check coexist:
  the report carries the ordering the fold consumes, and the shipped refusal text
  still comes from the post-join length check. No second overflow refusal is
  minted.
- **ASSUMPTION-3 (REQ-63 / REQ-99 / REQ-126 — the F5 fallback shape).** On a
  non-pollable pipe the join "stays unbounded", i.e. the bounded-join code path is
  skipped entirely for that invocation and `drains.Wait()` behaves as today; the
  pollability result is recorded on the `invocation` (EIA) and rendered as an
  extra `Detail` line on whatever refusal that invocation later produces. A
  non-pollable pipe does not itself refuse.
- **ASSUMPTION-4 (REQ-41 / REQ-42 — constant placement).** `DrainGrace` ships
  as an exported package-level constant in `internal/cli/cmdbind` beside
  `WaitDelay`/`StdoutCap`/`StderrTailCap`, at 50 ms. C1 names the value, not the
  identifier's export status; the sibling constants there are exported, and no
  cross-package reader is named, so unexported is also defensible — but the
  contract-naming convention it cites (0025:C4 names `StdoutCap` and
  `StderrTailCap`, both exported) makes the exported form the single defensible
  reading.
- **ASSUMPTION-5 (REQ-44 — how the tolerance is asserted).** The +100 ms
  scheduling tolerance applies to every timing assertion that cites the C1
  `bound:` total (MVV row 3, S1, S8), not only to MVV row 3; C1 states it "once
  so the scenarios cite it rather than each carrying its own". MVV row 5's
  "within `2·500 ms` of the exit" is a different anchor (the child's exit, not
  the invocation start) and takes the same tolerance.
- **ASSUMPTION-6 (REQ-128 — the stall seam's shape).** The single test-only
  drain-start delay hook is an unexported package-level `var` in
  `internal/cli/cmdbind`, nil in production, set and restored by the tests that
  need it (S7, S9, MVV row 6a). A12 confirms it needs no build tag. Exactly one
  such hook ships; S9 and MVV row 6 reuse S7's.
- **ASSUMPTION-7 (REQ-26 / REQ-83 — Detail composition).** The held-pipe reason
  and the stderr tail share one `Detail` slot, composed through the existing
  `cmdbind.go::wrap` / `withStderrTail` path with the applied-sense text first
  (MC `authority`, stderr-tail row). No new `Detail` slot or field is added.
- **ASSUMPTION-8 (REQ-139 — the new CLI code).** The exit-3 code for the applied
  write refusal is added to the `codeAccessor*` family in `internal/cli` beside
  `codeAccessorFailed`; D2 states the spelling is non-normative, so any
  `flow-write-failed-applied`-shaped id satisfies it, and no test may pin the
  literal beyond what the envelope golden requires.

## QUESTIONS

- **QUESTION-1 (REQ-26 vs REQ-132 — the order of the held-pipe reason and the
  stderr tail inside `Detail`).** C1 `refusal:` is explicit: the held-pipe reason
  is "composed AHEAD of the stderr tail collected up to the bound", and TD REQ-61
  agrees ("`Detail` leading with the held pipe(s) and the child's exit status
  ahead of the stderr tail"). S10 states the inverse: "the stderr tail collected
  up to the bound survives into `Detail` **ahead of the held-pipe reason**". The
  two readings produce materially different assertion text on the one row that
  exists to pin the single-pipe `Detail` form, and neither predecessor resolves
  it (0025:C4 orders `Detail` generally but not this pair; 0004 does not reach
  it). **Provisional reading taken in REQ-26/REQ-61/REQ-83, pending the
  orchestrator's answer**: the CONTRACT wins — the held-pipe reason leads, the
  stderr tail follows — because C1 is the normative fence, TD corroborates it,
  and S10 is a scenario whose load-bearing claim is that the tail SURVIVES into
  `Detail` (which holds under either order), not that it leads. If the answer is
  the other way, REQ-26, REQ-61 and REQ-83 change wording.
