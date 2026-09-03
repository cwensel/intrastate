# Recommendation 0026: Deliver the command timeout even when a detached grandchild holds the output pipes

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-08-31
- **Status**: Draft
- **Type**: Bug Fix
- **Profile**: foundational — accretion floor; C1 fixes the drain precedence at `cmdbind::spawn` (bounded wait outranks whole read, held-pipe refusal, the one bound plus DrainGrace); user-facing yes; locks cross-rdr
- **Priority**: High
- **Related Issues**: intrastate#936e (defect tracker; stays open until Stage 8); intrastate#v0hb (closed, RDR 0025's tracker); closed point-fixes at the seam: intrastate#1drn, intrastate#x5vq, intrastate#878e
- **Predecessors**: 0025-command-invoking-accessor-bindings, 0004-accessor-execution-safety-model
- **Overrides**: 0025:C4 (the "necessary and sufficient" deadline triple — to be amended to a bounded drain with a stated precedence) and 0025:F4's residue class (leakage of a process, never loss of the refusal)
- **Seam Lineage**: `internal/cli/cmdbind::reapGroup` / `internal/cli/cmdbind::drains.Wait` (`area:internal-cli`) — 3 prior closed point-fixes; trail: cd9a09d + intrastate#1drn (process-group syscall build tags), d22cd7a + intrastate#x5vq (grandchild-leak oracle), b4b7431 + intrastate#878e (cancel oracle); plus 76c121b, the manual `os.Pipe` ownership change inside 0025's own implementation (ADV-2 / FAIL-2), which is the change that opened this gap. Count confirmed at Resolve against `git log -- internal/cli/cmdbind/`: exactly the three closed point-fixes named, each with its tracker, plus 76c121b.

## Problem Statement

A model author who declares a `command` reader with a `timeout` (RDR 0025) expects that
when the child misbehaves, `intrastate` returns the `timeout` / `execution_failure`
refusal within a bound they can reason about, and their agent caller moves on. Today a
child that spawns a grandchild which leaves the process group via `setsid(2)` while
holding the inherited stdout/stderr write ends makes the CLI hang **unboundedly**: no
refusal is ever delivered, so a caller (an agent driving `flow resolve` in a pipeline)
waits forever with nothing to branch on. They discover it as a stuck invocation, not as
a refusal.

System-internally, 0025:C4 states a deadline triple (context deadline, `WaitDelay`,
group-scoped `SIGKILL` via `reapGroup`) and calls it "necessary and sufficient"; its
implementation (`internal/cli/cmdbind/cmdbind.go::spawn`, 76c121b) then promises, in the
comment that justifies the manual pipes, that the drains "run to a true EOF" with nothing
truncated — a promise C4's own text never makes. The drain join
(`drains.Wait()`) is unbounded, and the pipes are manually owned (`os.Pipe`, chosen in
76c121b so `Cmd.Wait` cannot truncate an in-flight read), which puts them outside
`Cmd.WaitDelay`'s reach. An escaped grandchild therefore survives `reapGroup`'s
group-scoped kill and keeps the write ends open; nothing bounds the join. 0025:F4's
residue admits only that a *process* may outlive the refusal — never that the
*refusal* is withheld — so this is a liveness hole in C4, not an accepted residue.

The decision this RDR owns: what C4 promises when "drains reach a true EOF, nothing
truncated" and "no unbounded wait" conflict — i.e. the precedence between a bounded
wait and a whole read, and the restated residue class that follows from it. Every
candidate repair changes what C4 says about output completeness, so it is a
contract-level amendment, not a patch.

## Critical Assumptions

- **A1 [Every command-binding invocation — read, gate, write, and a write's read-back through a command reader — passes through `internal/cli/cmdbind/cmdbind.go::spawn`, so one bounded drain join covers every path and no binding keeps a private wait]**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/cli/cmdbind/cmdbind.go::spawn` holds the package's only `exec.CommandContext` (`cmdbind.go:205`), only `cmd.Wait()` (`:282`) and only `drains.Wait()` (`:302`) — a sweep of the package for `exec.Command`/`.Wait()`/`drains` returns no other non-test site, so no binding keeps a private wait. Its three callers are `cmdbind.go::Reader.Read` (`:566`), `cmdbind.go::Gate.Gate` (`:657`) and `cmdbind.go::(*Writer).Apply` (`:722`), each of the form `inv, err := spawn(...)` followed immediately and unconditionally by `if err != nil { return … }` (`:569`, `:660`, `:725`) — ahead of every `inv.stdout` / `inv.exitCode` touch, so no caller has a fast path that would parse a held stream. Between `cmd.Wait()` (`:282`) and the join (`:302`) the only statement is `reapGroup(cmd.Process)` (`:297`): no branch, no return; every `werr`-switch arm returns at `:304` or later, after the join.
  - **If wrong**: a path outside `spawn` keeps the unbounded join and the hang survives on that verb only — visible as the same stuck invocation, on gates or writes. The same search confirms every caller reads `err` before `inv` (a gate fast path on `inv.exitCode` would parse a held stream) and that no arm of `spawn` returns between `Wait` and the join.
- **A2 [The joining goroutine can END a pending blocking `Read` on an `os.Pipe` read end from another goroutine — `SetReadDeadline` or `Close` on the read end — on darwin and linux, and the drain goroutine then returns promptly with the bytes it had read]**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `evidence/spikes/a2-unblock-read/` (`a2_unblock_read.go` S1-S6, `a2_buffered_bytes.go` T1-T3), run on darwin/arm64 go1.26.6 natively and linux/arm64 go1.26.8 under `golang:1.26`, with a `SysProcAttr{Setsid: true}` grandchild holding the write end so no EOF is possible. The reader is confirmed genuinely parked (still blocked at 300 ms, writer alive) before each measurement. `SetReadDeadline` returns nil and ends the blocked `Read` in **101.5 µs** (darwin) / **33.4 µs** (linux) with `&fs.PathError{Op:"read", Err: os.ErrDeadlineExceeded}` — `errors.Is(err, os.ErrDeadlineExceeded)` true. A drain loop that accumulates before inspecting `err` keeps what it already read (24/24 bytes; loop ended 48.8 µs / 34.4 µs after the deadline). The `Close` variant also unblocks (232 µs / 104 µs, `os.ErrClosed`) but makes unread buffered bytes unrecoverable, which is why the deadline — not `Close` — is the mechanism. Results were identical on both OSes on every point.
  - **Refines C1**: the spike falsified "a read deadline of **now**" as the way to keep buffered bytes. A past deadline is checked BEFORE the read is attempted, so it short-circuits with `n=0` and skips bytes still in the pipe (they are deferred, not lost — clearing the deadline returns all of them). Delivering buffered bytes requires a SHORT FUTURE deadline: at `now+grace` the first `Read` returns the whole buffered payload with `err == nil`, and only the next read, on an empty pipe with a live writer, reports the deadline. C1's `precedence:` line carries this as the drain discipline.
  - **If wrong**: the join's timer fires but the drain goroutine never returns; if `spawn` then returns anyway it leaks the goroutine and the fd per invocation, and if it waits the hang is unchanged. The spike also records what `SetReadDeadline` returns on a pipe whose poller registration failed (`os.ErrNoDeadline`), which is F5's trigger.
- **A3 [After `Wait` + `reapGroup`, a drain still blocked on an EMPTY pipe past the bound has already consumed every byte the exited direct child wrote — pipe FIFO delivers outstanding data before EOF — so the only writer left is a process outside the group, and refusing is a policy choice, not loss of the direct child's answer]**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `evidence/spikes/a3-a6-drain-bound/` (`fixture.go`, `spike.go`), darwin/arm64 go1.26.6 and linux/arm64 go1.26.8, identical on both. The adversarial ordering is what makes it evidence: the stdout drain reads nothing until AFTER `cmd.Wait()` returned and the group SIGKILL was sent, so every recovered byte was read strictly after the writing child was gone. Completeness is asserted with `bytes.Equal` against an independently reconstructed envelope, not by length. Delayed-drain result: **complete at 1 KiB and 63 KiB**, child exit 0 — the direct child's whole answer survives the bound with only an escaped writer left, confirming TLPI §44.3's rule holds in the implementation.
  - **Antecedent (scope, not a refutation)**: past the ~64 KiB pipe buffer with the drain delayed, the child **cannot exit** — it blocks in `write(2)`, is still blocked at the deadline, and dies to the group signal; exactly one pipe buffer (65536 bytes) survives, as a correct prefix. A3 speaks of "every byte the EXITED direct child wrote": a child blocked in write never exited and never finished writing, so no complete answer existed for FIFO to deliver and nothing successfully written was lost. The guarantee is therefore conditional on the drain keeping up — and the shipped concurrent-drain design is what keeps that antecedent satisfiable: with the drain prompt, the same 1 MiB completes cleanly on both OSes.
  - **If wrong**: C1's residue wording ("its output is never read") understates the loss — the direct child's own bytes can be lost at the bound — and the Failure Modes must say so.
- **A4 [The executor's deadline-first classification (`internal/accessor/executor.go::invokeRead`, and its Gate and Apply counterparts) turns the new refusal into `timeout` whenever the declared timeout elapsed, on every path, so `spawn` needs no timeout class of its own and the binding's error is read only when the deadline did not expire]**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: on the three DIRECT binding paths the `DeadlineExceeded` test is the first statement after the binding call and textually precedes the `if err != nil` arm — read `internal/accessor/executor.go::invokeRead:199` ("Deadline first: … `timeout` is the one reported"), gate `(*Executor).Gate:231`, write `(*Executor).Write:356` (`appliedDeadline := errors.Is(applyCtx.Err(), context.DeadlineExceeded)`, consumed at `:359` ahead of the `err != nil` arm at `:368`). The fourth path, the read-back leg, reaches the same conclusion by a DIFFERENT mechanism, stated here rather than folded into the sentence above: `executor.go:392` tests `raw.class == ClassTimeout`, a class `invokeRead` already derived from ITS deadline check, so there is no second `DeadlineExceeded` test to precede anything. Deadline-first still holds transitively — the precedence was decided one frame down — and a read-back that times out returns the WRITE's refusal (`WriteResult{Refusal: …}`, `:396`) with class `ClassTimeout` and `applied = true`, not laundered into a read result, with the generic `raw.class != ""` arm following at `:398`. No path reads the binding's error before the context deadline.
  - **If wrong**: a bound-expired drain after the deadline is reported as `execution_failure` with a held-pipe reason instead of `timeout`, so a caller branching on class is misled.
- **A5 [500 ms (`cmdbind.go::WaitDelay`) is long enough that a NON-escaping backgrounded grandchild SIGKILLed by `reapGroup` releases its pipe ends within the bound on CI, so the ordinary "wrapper script backgrounds a helper" success case (intrastate#x5vq's oracle) pays nothing and is never refused]**
  - **Status**: Verified
  - **Method**: MVV Test
  - **Evidence**: `WaitDelay = 500` ms confirmed at `internal/cli/cmdbind/cmdbind.go:66` (applied `:224`). The oracles — `adversarial_0025_test.go::TestAdvBackgroundGrandchildDoesNotSurviveASuccessfulRead:313` and `::TestAdvBackgroundGrandchildMustNotStallASuccessfulRead:258` (intrastate#x5vq), `::TestAdvKilledChildOnParentCancelIsNotAnAnswer:131` (intrastate#878e) — pass 75/75 under `-race -count=25` with no race reports, run concurrently with CPU load. Because `spawn` joins unbounded today (`cmdbind.go:302`), the oracle proves termination but not margin, so the margin was measured directly: `evidence/spikes/a5-ingroup-margin/` times `kill(-pgid, SIGKILL)` → drain EOF over 160 trials (40 idle; 60 at loadavg 16.2; 60 at loadavg 47.5 on 12 cores). Median 0.37/1.19/0.89 ms, p99 0.46/19.9/15.3 ms, worst observation anywhere 27.7 ms = 5.5% of the bound. Zero of 160 trials reached or approached 500 ms; every trial returned the correct envelope.
  - **If wrong**: correct tools whose helpers take longer to die are refused as held-pipe failures — a regression of the prior point-fix, seen as intermittent `execution_failure` on a working reader.
- **A6 [The invocation returns within `timeout + 2·WaitDelay` in the escaped-grandchild case — the TOTAL is the claim: `Cmd.Wait` returns well inside one `WaitDelay` of the context deadline (its own timer bounds the stdin-copy goroutine and the post-kill wait) and donates the remainder to the drain join, which spends one full `WaitDelay` plus a scheduling epsilon. The two legs are not separately bounded]**
  - **Status**: Verified, with the per-leg claim narrowed
  - **Method**: Spike
  - **Evidence**: `evidence/spikes/a3-a6-drain-bound/`, fixture **FX-deadline-escape** (Go helper, `SysProcAttr{Setsid: true}` grandchild inheriting both pipes), both shapes, n=10 each, timeout 1.5 s against `WaitDelay = 500` ms (`cmdbind.go:66`), budget 2.5 s. **The escape is real on both OSes** — darwin child pid/pgid 37707/37707 vs grandchild 37709/37709; linux 828/828 vs 834/834; "grandchild alive after `kill(-pgid, SIGKILL)`" true on both, so the drain genuinely sees no EOF. **The total bound holds with wide margin**: worst total 2.013 s (darwin) / 2.082 s (linux) against 2.5 s — 487 ms and 418 ms of headroom. Leg (a), deadline→`Wait`, holds comfortably at 1–79 ms, two orders of magnitude inside one `WaitDelay`.
  - **Narrows C1 `bound:`**: leg (b), `Wait`→join, measured **500.1–508.2 ms in all 40 samples — never under 500 ms**. This is structural, not noise: the escaped grandchild holds the pipes so the join always waits the full timer, then pays timer-fire latency plus goroutine scheduling to observe the deadline and unwind. So "the drain join is bounded by one more `WaitDelay`" is strictly false; the true per-leg bound is one `WaitDelay` plus a small scheduling epsilon (worst overshoot 8.2 ms linux, 3.0 ms darwin). The TOTAL claim survives unharmed because leg (a) consumes almost none of its allowance and absorbs the overshoot many times over. Scope caveat: these figures are for a child that dies promptly to the group kill; a child that RESISTS termination would push leg (a) toward its full `WaitDelay`, and only then does the budget tighten — that case was not exercised.
  - **If wrong**: C1's stated bound is false; a caller sizing its own watchdog from `timeout` still waits longer than promised.
- **A7 [The child's stdin is already bounded and its expiry already refuses: `cmd.Stdin` is a `bytes.Reader`, so `Cmd` owns that pipe and its copy goroutine under `cmd.WaitDelay`; a child (or escapee) that never reads it ends in `exec.ErrWaitDelay` from `Wait`, which reaches `spawn`'s non-`ExitError` arm — after `reapGroup` and the join — and refuses `execution_failure`]**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `cmdbind.go:207` assigns `cmd.Stdin = bytes.NewReader(stdin)` (non-`*os.File`) and `:224` sets `cmd.WaitDelay = WaitDelay * time.Millisecond` (`WaitDelay = 500`, `:66`). Go's own text covers the stdin goroutine without the gloss: `os/exec/exec.go` WaitDelay doc — "those pipes are closed in order to unblock any goroutines currently blocked on Read or **Write** calls" — and the mechanism is in code, not only doc: `exec.go::childStdin:542` hands an `*os.File` stdin straight to the child with no goroutine, but the `bytes.Reader` branch (`:545-561`) appends the write end to `c.parentIOPipes` and the `io.Copy` to `c.goroutine`, which `awaitGoroutines` force-closes at `exec.go:999` before returning `ErrWaitDelay` (`:1002`). `ErrWaitDelay` is not an `*exec.ExitError`, so `errors.As` at `cmdbind.go:311` fails and control reaches the `default:` arm at `:320-322` — after `reapGroup` (`:297`) and after the join (`:302`) — refusing `execution_failure`. The third pipe does not re-open the hole.
  - **If wrong**: a write whose payload exceeds the pipe buffer hangs on a non-reading child — the third inherited pipe re-opens the hole this RDR closes for the other two.
- **A8 [The executor's write arm can derive `Applied()` from the binding's typed held-pipe error: `internal/accessor/executor.go::Write` has, or can gain, a refusal site that matches `cmdbind.HeldPipeError` with `errors.As` after the direct child was reaped and sets the unexported applied sense there; the CLI's "may have been applied" rendering must GAIN its `Applied()` key on the DIRECT write arm, since the shipped mapping keys on class and phase alone and that one case slips past it — the read-back leg already carries both the applied sense and the rendering (JDR 0003 §D1 rule 3, which assigns that gain to this record)]**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: the setter is reachable — `internal/accessor/model.go:314` holds the unexported `applied bool` whose doc reserves it for the write path, `:323` the `Applied()` accessor; `executor.go::Write` already sets `r.applied = true` at six refusal sites (`:363`, `:375`, `:393`, `:401`, `:409`, `:472`) by direct in-package field assignment, so the proposed `errors.As` match takes the `if err != nil` arm at `:368` with no structural change. The rendering is the gain, not a found fact: `internal/cli/flow_exec.go::accessorFailureOf` (`:395`) keys `detailMayHaveApplied` (`:469`) on class plus phase at `:398-404` (`ClassTimeout` + `at == phaseWrite`) and `:410-414` (`ClassReadBackIncomplete`) and never calls `Applied()` — which has zero non-test callers repo-wide. Decisively, `ClassExecutionFailure` at `:410-412` carries no phase check and sets no `Detail` of its own, so a held-pipe write on the DIRECT invocation renders today with no applied sense. The held-pipe REASON still reaches the user one frame up — `flow_exec.go::accessorFailure` wraps every mapping in `withStderrTail(…, refusal.Detail)`, which lands the tail in the one `Detail` slot — so what the gain adds is the applied-sense sentence ahead of that tail, not the diagnostic itself. The read-back leg is NOT part of the gain: `ClassReadBackIncomplete` at `:410-423` already sets `detailMayHaveApplied`, and `executor.go`'s read-back arm already sets `applied`, so the re-key is one site (the `err != nil` write arm), not two. JDR 0003 §D1 rule 3 records exactly this ("the shipped mapping keys on class and phase, which this case slips past") and lands the `Applied()` derivation here, so the gain is a decided constraint this record implements, not an open question. Its cost is carried in C1 `refusal:` and Testing Strategy S1.
  - **If wrong**: a held pipe on a write is rendered as a write that did not occur — the guessed absence 0004:C14 forbids — and an agent caller retries blind instead of reading back.

## Proposed Solution

### Approach

Decide once what 0025:C4 left undecided — **when a bounded wait and a whole read
conflict, the bound wins** — and put the bound where C4 already says it is. The undecided
contract the three prior point-fixes danced around: *what does the command binding promise
when the only remaining writer on its output pipes is a process its termination cannot
reach?* Each fix made the reachable writer set larger (group signal on the success path,
build-tagged syscalls, a cancel oracle) without ever stating what happens at the edge of
that set, so the edge stayed an unbounded wait.

The answer this RDR locks: in `internal/cli/cmdbind/cmdbind.go::spawn` — the one function
every command binding (read, gate, write, read-back) invokes — the order stays as it is
(`cmd.Wait()`, then `reapGroup`, then the drain join), and the join becomes bounded by the
existing `cmdbind.go::WaitDelay` (500 ms), not by a second constant. A drain still open at
the bound is ended by the parent closing (or deadlining) its own read end, and the
invocation **refuses** through the existing `*accessor.ExecError` channel, naming the pipe(s)
held past the bound and carrying the stderr tail as today. The executor's deadline-first rule
already classifies that refusal as `timeout` when the declared timeout had elapsed; otherwise
it is `execution_failure`. Output read on a bound-expired drain is never parsed on any path.
The residue class of 0025:F4 is restated honestly: a process outside the child's group may
outlive the refusal AND keep the pipes it inherited, and its output is never read — but the
refusal is never withheld, and the invocation returns within `timeout + 2·WaitDelay`.

This follows Go's own disposition for the same conflict, quoted rather than paraphrased:
`os/exec/exec.go::(*Cmd).awaitGoroutines` — "If c.WaitDelay elapses before the goroutines
complete, awaitGoroutines forcibly closes their pipes and returns ErrWaitDelay" ⇒ the parent
closes its own ends and a success-status child with an unclosed pipe is REPORTED as an
error, never as a silent success. Go leaves accept-or-refuse to the caller (`Output` returns
the bytes beside `ErrWaitDelay`); the refusal is this RDR's choice, and it is the one both
peer CLIs that met this failure made — roborev maps `exec.ErrWaitDelay` to its context
error, beads sets the bound for exactly the inherited-pipe grandchild (Investigation). `spawn` hands `*os.File` ends to `cmd.Stdout`/`cmd.Stderr`
precisely so the group can be released before the drains are joined (76c121b), which is why
`Cmd`'s bound never reached them; this RDR applies the same rule to the owned drains.

### Technical Design

`spawn`'s post-`Wait` sequence today: `reapGroup(cmd.Process)` (group SIGKILL; releases every
writer the group signal can reach) → `drains.Wait()` (unbounded join of the two
`readBounded` goroutines). The design changes only the join:

1. `reapGroup` as today — release first, so the common backgrounded-grandchild case closes its
   ends before any timer matters.
2. Join the drains with a `WaitDelay` bound. On EOF within the bound: unchanged — whole
   output, the 1 MiB cap and 4 KiB tail apply as today.
3. When the timer fires: the parent sets `SetReadDeadline(now + DrainGrace)` on
   `outR`/`errR` (a pollable `os.Pipe` file wakes a pending `Read`; a read with bytes
   buffered still returns them, so a merely LATE drain finishes whole and only an empty
   pipe with a live writer reports `os.ErrDeadlineExceeded`), joins the drains (now
   returning), closes both read ends, and — only if a drain reported the deadline — returns
   `wrap(detail, errHeldPipe)` — `Detail` leading with the held pipe(s) and the child's exit
   status ahead of the stderr tail, `Err` the typed held-pipe error (JDR 0003 §D1).
   Nothing read from a held stdout is handed to `Reader.parse`, `Gate.Gate`'s verdict
   mapping, or the writer's read-back. Pollability is checked once at pipe creation (a
   far-future deadline; `os.ErrNoDeadline` marks the pipe non-pollable) so the fallback is
   decided before any child exists (F5).
4. Classification is unchanged and lives where it lives today: the executor reads the bounded
   context's `DeadlineExceeded` first (⇒ `timeout`), then the binding's error
   (⇒ `execution_failure`). `spawn` therefore never has to know whether the bound expired
   because of a deadline or because a success-path child leaked a writer.

State read cited above: `readBounded`'s pending `Read` is the state; its writer is the
escaped process holding the pipe's write end — set at the grandchild's `setsid(2)` + inherit,
cleared only by that process closing or exiting, which nothing in this CLI can force. That is
exactly why the bound must be on the parent's side.

#### Normative Contracts

The drain precedence, amending 0025:C4's `deadline:` line (0025 is not edited; this
clause is the successor text and `cmdbind.go`'s C4 comments re-cite it).

**C1**

```normative
precedence: a bounded wait outranks a whole read. Order in `spawn` is fixed: Wait (direct child reaped) → reapGroup (one kill(2) to -pgid, non-blocking) → drain join under ONE timer of WaitDelay covering BOTH read drains. When the timer fires the parent sets a read deadline of now + DrainGrace (50 ms) on its own read ends; each drain's FINAL read then reports the pipe's state, not the clock: bytes still buffered are delivered, EOF (every writer gone) is whole output as today, and a deadline error is a HELD pipe — empty, with a writer the group signal could not reach. The grace is what MAKES that final read happen: Go checks the deadline BEFORE attempting the syscall (`internal/poll/fd_unix.go::(*FD).Read` calls `prepareRead` ahead of `syscall.Read`), so a deadline of NOW short-circuits with n=0 and reports a held pipe without ever looking at the pipe — the elapsed-time refusal the premortem rejected (P-2: refuse on "no EOF after a final bounded read", not on elapsed time), and, one layer down, the guessed absence 0004:D-absent-vs-unreadable forbids. Verified on darwin and linux (A2). Only a held pipe refuses. After the join every read end is closed (deadline to unblock, close to release the fd)
refusal:  a held pipe refuses through `*accessor.ExecError`, this record's instance of JDR 0003 §D1: `Detail` carries the held-pipe reason — the pipe(s) held ("stdout"/"stderr"/both), the bound, the direct child's exit status (exited N / killed by signal) and the remediation ("close or redirect the helper's inherited stdio") — composed AHEAD of the stderr tail collected up to the bound; `Err` wraps the typed held-pipe error `cmdbind.HeldPipeError` (pipes held, exit status), which the executor matches with `errors.As` at the refusal site. Applied sense: a held pipe on the write path after the direct child was reaped is a post-run refusal — `Applied()` is true, carried per JDR 0003 §D1 through the executor's write arm. The write's TWO legs already differ in class and this clause does not flatten them. The DIRECT write invocation refuses `execution_failure`, and it is the leg that needs the gain: `executor.go::Write`'s `err != nil` arm sets no applied sense today, and `flow_exec.go::accessorFailureOf`'s `ClassExecutionFailure` arm carries neither a phase check nor a `Detail`, so a held-pipe write slips past the "may have been applied" rendering (JDR 0003 §D1 rule 3 lands that gain here). The READ-BACK leg after the write ran is already whole and gains nothing: a held pipe there returns the write's own refusal re-classed `read_back_incomplete` with `applied` already set (`executor.go`'s `raw.class != ""` arm), which `accessorFailureOf`'s `ClassReadBackIncomplete` arm already renders with `detailMayHaveApplied`. So one site changes, not two — the read-back leg is named here because a reader who assumes the class is uniform across both legs would re-key an arm that is correct. On read and gate invocations `Applied()` is false as today. The returned invocation carries no stdout; stdout read on a held drain is never parsed on ANY path — read envelope, gate verdict, write, write read-back (JDR 0003 §D1)
class:    the class is the executor's, as today: ctx `DeadlineExceeded` ⇒ `timeout` (deadline-first rule; per 0025:C4 a `timeout` refusal carries no Err/Detail, so the held-pipe reason survives only on `execution_failure`), else the ExecError ⇒ `execution_failure`. The binding never classifies `timeout` itself
stdin:    unchanged and named: the child's stdin is a Cmd-owned pipe (`cmd.Stdin` is a `bytes.Reader`), bounded by Cmd's own WaitDelay; a non-reading child or escapee holding its read end ends in `exec.ErrWaitDelay` from Wait, which `spawn`'s non-ExitError arm refuses as `execution_failure` AFTER reapGroup and the join — no arm of `spawn` returns between Wait and the join
bound:    WaitDelay (500 ms) is the ONE bound — the same constant, by design, for every "the CLI waits past the child's exit" case: Cmd's stdin write and post-Cancel wait, and the owned drain join; no second knob, not model-declarable. DrainGrace (50 ms) is NOT a second bound: it answers a different question (how far ahead the read deadline must sit for the kernel to deliver already-buffered bytes), is never model-declarable, and cannot move the refuse/accept boundary — it is a mechanism constant inside the one bound, named here because every shipped constant of this family is named in its contract (0025:C4 does the same for StdoutCap and StderrTailCap). Cancel sends one signal and returns, so Cmd's timer starts at the deadline. Per invocation, the COMMITTED return bound is timeout + 2·WaitDelay; a command write is two invocations (write, then read-back), each under its own bound. The two legs are NOT separately bounded and the contract does not claim they are: the drain join runs one full WaitDelay plus timer-fire and goroutine-scheduling latency (measured 500.1-508.2 ms in 40/40 samples, A6), which the deadline leg's unused allowance absorbs. That absorption is why the total holds, so it is stated rather than left as a coincidence between two independent decisions (0025:C6's rule); a child that RESISTS termination spends leg (a)'s allowance and is the case where the total tightens — S8 tests it
whole:    a drain that reaches EOF — before the timer, or on its final read after it — is whole, as today; the 1 MiB stdout cap and 4 KiB stderr tail are unchanged. A child (and group) that closes its pipes observes no change. The completeness guarantee (A3) carries an antecedent, stated here because nothing structural enforces it: the two drains run CONCURRENTLY and are never stalled. A child whose payload exceeds the ~64 KiB pipe buffer blocks in write(2) until drained, so a serialized or stalled drain converts a large-answer success into a deadline kill — measured: exactly one pipe buffer (65536 bytes) survives, as a correct prefix, where a prompt drain returns the whole 1 MiB. `spawn`'s concurrent-drain comment already reasons this way for the verbose-tool deadlock it cites 0025:F4 for; this clause is what keeps a later refactor from silently voiding A3, since the prohibition is carried by the contract and S7, not by a build error
residue:  0025:F4 restated — a process outside the child's process group may outlive the refusal; it keeps write ends whose reader is gone, so its next write fails with EPIPE (SIGPIPE under the default disposition), and whatever it wrote is never read. Leakage of a process, or of THAT process's output, is admitted; a withheld refusal is not. Diagnosis: the refusal's Detail names the held pipe(s); the leaked process is visible to `ps` until its next write
```

#### Conditional Mini-Checks

Five structural cues fire on this draft. Each table is the decision, recorded here.

**`authority` — source-authority census.** Cue: a fallback path (F5), a derived
applied-sense, and two arms classifying the same refusal.

| Input / decision | Writer | Readers | Call sites | Sibling arms | Canonical |
| --- | --- | --- | --- | --- | --- |
| Refusal class (`timeout` vs `execution_failure`) | the executor | CLI rendering, agent caller | `executor.go::invokeRead`, `Gate`, `Write` | `spawn` returns an error, never a class | the EXECUTOR; the binding never classifies `timeout` (C1 `class:`) |
| Applied sense on a write | `executor.go::Write`, by direct field assignment | `flow_exec.go::accessorFailureOf` | direct arm (`err != nil`, gains it) and read-back arm (already sets it) | the two legs carry DIFFERENT classes | `Refusal.applied`; `Applied()` is the only reader key (A8) |
| Held-vs-EOF terminal condition | each drain goroutine, reporting out | `heldPipes()` at the join | one join in `spawn` | none — `readBounded` today swallows it | the drain's reported condition; NOT the clock (C1 `precedence:`) |
| Stderr tail into `Detail` | `cmdbind.go::wrap` | `withStderrTail` composes into one slot | one slot, applied-sense first | held-pipe reason and tail share the slot | `refusal.Detail`, ordered per 0025:C4 |
| Drain bound | `cmdbind.go::WaitDelay` | `spawn`'s join timer, Cmd's own timer | one constant, two timers | `DrainGrace` is an OFFSET, not a rival bound | `WaitDelay`; no second knob (D-naming) |

**`oracle` — test-discriminability.** The rows that could pass by absence-of-error
carry their named failing control.

| MVV / scenario | Fails if X is wrong because Y | Negative / failing control |
| --- | --- | --- |
| MVV row 3 (deadline shape) | asserts elapsed ≤ 3 s AND class `timeout` AND empty `Detail`; a withheld refusal never returns | today's code hangs unboundedly — the bug itself is the control |
| MVV row 5 (success shape) | asserts `execution_failure` within `2·WaitDelay`, `Err` naming "stdout, stderr" and "exited 0", envelope NOT parsed | a build that parses held stdout returns values and fails the not-parsed assertion |
| MVV row 6 (late-but-unheld) | a deadline of `now` short-circuits before the syscall and loses the buffered bytes | **this row IS the control** — it fails under the rejected mechanism (A2) |
| S2 (existing oracles re-run) | passes by absence-of-change, so it is paired with S7's stalled arm | S7's stalled drain recovers exactly 65536 bytes where a prompt drain returns 1 MiB |
| S3 (ordinary paths byte-for-byte) | absence-of-change oracle; discriminated by the cap boundaries | one byte past the cap must still refuse overflow |
| S9 (late drain, no holder) | accepted-whole vs refused is the whole point | companion: escaped writer emitting a byte every 10 ms is REFUSED |

**`fidelity` — round-trip / fidelity.** Byte-equality invariants this record commits.

| Operation | Invariant | Lossy exemption |
| --- | --- | --- |
| Whole read, child closes pipes | byte-for-byte unchanged from today (C1 `whole:`, S3) | none |
| stdout at the 1 MiB cap | exactly the cap is a value; one byte more is a detected overflow | the discarded remainder — bounded on what is KEPT, as today |
| stderr tail | LAST 4 KiB preserved | the earlier bytes, as today |
| Drain past the ~64 KiB pipe buffer, prompt | whole 1 MiB recovered (S7) | none while the drains stay concurrent |
| Drain past the pipe buffer, STALLED | exactly one pipe buffer (65536 bytes) as a **correct prefix** | the remainder — admitted, and the reason S7 pins the concurrency (A3, D-the-drains-stay-concurrent) |
| Held pipe | stdout is never parsed on ANY path | the whole stream — refused, not truncated-then-read |

**`disposition` — input class → outcome.**

| Input class | Exit / class | Event / error | Artifact minted | Silent vs loud |
| --- | --- | --- | --- | --- |
| Child closes pipes, exits 0 | success | none | envelope parsed | loud (the value) |
| Held pipe, deadline NOT elapsed | `execution_failure` | `HeldPipeError` in `Err`, reason + tail in `Detail` | no stdout | LOUD |
| Held pipe, deadline ALSO elapsed | `timeout` | none — 0025:C4 gives `timeout` no Err/Detail | no stdout | loud in class, **quiet in reason** (F6, S6) |
| Held pipe on the DIRECT write | `execution_failure` + `Applied()` true | held-pipe reason | no read-back | LOUD once the re-key lands (A8) |
| Held pipe on the write's READ-BACK | `read_back_incomplete` + `applied` already true | read-back's own error | no confirmation | LOUD today |
| Non-reading child, stdin past the buffer | `execution_failure` | `exec.ErrWaitDelay` via `spawn`'s non-ExitError arm | none | loud (S4, A7) |
| Non-pollable pipe (`os.ErrNoDeadline`) | unbounded join — the old hang | logged AT THE CHECK | none | **loud at the check, by design** (F5, S5) |
| Escapee outside the group | not reaped | EPIPE on its next write | none | **SILENT — admitted residue** (C1 `residue:`, F3) |

**`trace` — desk trace of the MVV.** Assertions in force at each step, with a witness.

| Step | Assertions in force | Witness |
| --- | --- | --- |
| 1. child spawns `Setsid` grandchild, sleeps past deadline | C1 `residue:` (escapee unreachable); A1 (one spawn site) | grandchild alive after refusal, killed by the test |
| 2. ctx deadline expires at 2 s | C1 `class:` deadline-first; A4 | `timeout` wins over any binding error |
| 3. `Cancel` kills the group; `Wait` returns | C1 `stdin:` (no arm returns between Wait and the join) | `werr` held, not returned |
| 4. `reapGroup` signals the group | C1 `residue:`; 0025:C4 | escapee survives — that is the case under test |
| 5. join under ONE timer of `WaitDelay` | C1 `bound:` (one bound, both drains) | drain-join leg measured 500.1–508.2 ms, 40/40 (A6) |
| 6. timer fires → deadline of now + `DrainGrace` | C1 `precedence:`; D-the-drain-grace | final read RUNS; A2 confirms a deadline of `now` would not |
| 7. final read reports pipe state | C1 `precedence:` (held = empty + live writer) | empty + deadline error ⇒ held |
| 8. read ends closed | C1 `precedence:`; Infra-Audit close-ownership row | defers already close them; no double-close hazard |
| 9. classify | C1 `class:`; F6 | deadline elapsed ⇒ `timeout`, empty `Detail` (MVV row 3) |
| 10. return | C1 `bound:` total | ≤ 3 s asserted; worst measured 2.013 s / 2.082 s (A6) |

No CONTRADICTION row. Steps 6 and 9 are the two that a naive reading collides:
the grace makes the final read happen, and the class still flips to `timeout`
when the deadline elapsed — so the held-pipe REASON is computed and then not
rendered. That is F6, recorded as under-explained rather than contradictory.

#### Load-Bearing Decisions

- **Naming** — the bound is `cmdbind.go::WaitDelay`, reused; rejected: a new `DrainBound`
  constant (a parallel knob for one decision), and a model-declarable drain bound (the
  author declared `timeout`; the drain bound is the CLI's own mechanism cost, not policy).
  The reuse is deliberate, not incidental: Cmd's stdin/post-Cancel timer and the drain
  timer are two timers sharing one value because they answer one question ("how long past
  the child's exit does the CLI wait"), so a future tuning moves both — accepted.
- **The drain grace is a mechanism constant, not the rejected knob** — `DrainGrace`
  (50 ms) ships beside `WaitDelay`, and the distinction is the criterion the row above
  states: `DrainBound` was refused for answering the SAME question with a second value;
  `DrainGrace` answers a different one — how far ahead the read deadline must sit for the
  kernel to deliver bytes already buffered. It is not model-declarable, moves no
  refuse/accept boundary, and invents no second signal (the Sibling-path check's test):
  it is an offset on the single existing timer's expiry. It is named in C1 `bound:`
  because a shipped constant of this family is always named in its contract (0025:C4).
  **This is a deliberate divergence from prior art, recorded as one**: Go and both peer
  CLIs hard-close and accept the loss — `os/exec/exec.go::awaitGoroutines` calls
  `closeDescriptors` and discards the in-flight read, and no peer in the surveyed set
  uses a future deadline to drain (the only deadline-unblock idiom found anywhere is
  `SetReadDeadline(time.Now())`, which discards the tail on purpose). We diverge because
  the bytes at stake are the refusal's own diagnostic: on a held pipe stdout is refused
  and never parsed, so what the grace preserves is the stderr tail that `wrap` carries
  into `Detail` — the text that tells a user WHICH helper held the pipe. 50 ms is chosen
  against the measured scheduling tail (A5 saw drain-goroutine latency reach ~28 ms under
  4x-core saturation), so the grace clears it with margin while costing 10% of one bound.
- **The drains stay concurrent** — A3's no-loss guarantee holds only while they do, and
  nothing in the type system prevents serializing them. `spawn` already comments the
  verbose-tool deadlock this avoids (a child filling the stderr pipe buffer while the
  parent reads stdout); the new fact is that the same coupling scopes A3, since a child
  past the pipe buffer cannot exit until drained. Recorded here, and pinned by S7, so the
  prohibition is load-bearing rather than decorative.
- **Selection / predicate** — when EOF and the bound both qualify: EOF wins if it arrives
  first (whole output, as today); the bound wins otherwise and the invocation refuses.
  When `timeout` and `execution_failure` both qualify: `timeout` (the executor's existing
  deadline-first rule, not a new tie-break).

#### Illustrative Code

Shape only — the mechanism inside step 3 is the A2 spike's.

```go
werr := cmd.Wait()          // direct child reaped (≤ WaitDelay past ctx deadline)
reapGroup(cmd.Process)      // release the group FIRST, as today
if !waitBounded(&drains, WaitDelay) {
    deadlineIn(outR, errR, DrainGrace)  // a FUTURE deadline: the final read still runs,
                                        // so buffered bytes drain and only an empty pipe
                                        // with a live writer reports the deadline
    drains.Wait()                       // now returns; each drain reports EOF or deadline
}
closeReads(outR, errR)          // release the fds on every path
if held := heldPipes(); len(held) != 0 {
    return inv, wrap(tail(stderr), errHeldPipe(held, exitStatus))
}
// EOF (before or after the timer): whole output, unchanged from here on
```

Illustrative: `waitBounded`/`deadlineIn`/`closeReads`/`heldPipes`/`errHeldPipe` are
names for the shape, not the contract; the timings in the MVV are Normative.

One shipped detail the shape hides, and C1 `precedence:` depends on: `readBounded`
(`cmdbind.go:352-356`) returns its bytes on ANY error, so today a deadline error is
indistinguishable from a clean EOF at the join. The held-vs-EOF distinction is therefore
not a fact the current code exposes — each drain must report its terminal condition out
(EOF or `os.ErrDeadlineExceeded`) for `heldPipes()` to have anything to read. That
plumbing is part of this change, not a pre-existing capability.

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| One spawn site for every command path | `internal/cli/cmdbind/cmdbind.go::spawn` | none known (A1 confirms) | Extend | C1 lands once; no per-binding wait |
| Drain bound constant | `internal/cli/cmdbind/cmdbind.go::WaitDelay` | today applied only to `Cmd`'s own goroutines | Reuse | C1 `bound:` — no second bound |
| Read-deadline grace | none — no grace/slack constant exists in this repo or the corpus | n/a | New | C1 `bound:` names `DrainGrace` (50 ms); a mechanism constant, not the knob D-naming rejected |
| Group release before join | `internal/cli/cmdbind/cmdbind.go::reapGroup` / `procgroup_unix.go::killGroup` | cannot reach a `setsid(2)` escapee | Reuse | C1 `residue:` names the escapee as the admitted leak |
| Bounded, overflow-detecting read | `internal/cli/cmdbind/cmdbind.go::readBounded` | blocks in `Read` until EOF or an ended fd, and SWALLOWS the error (`:352-356` returns `b` on any `err`), so a deadline is today indistinguishable from EOF | Extend | the join ends the read from outside; `readBounded` must additionally report its terminal condition so `heldPipes()` can tell held from whole |
| Refusal channel with stderr tail | `*accessor.ExecError` via `cmdbind.go::wrap` (0025:C4 `detail:`) | none | Reuse | the held-pipe reason rides in `Err`, tail in `Detail` |
| Deadline-first classification | `internal/accessor/executor.go::invokeRead` (+ Gate / Apply) | A4 checks the write and read-back paths | Reuse | C1 `class:` — no binding-side `timeout` |
| Stdin bound | `cmd.Stdin` = `bytes.Reader` + `cmd.WaitDelay` (Cmd-owned pipe) | A7 confirms `ErrWaitDelay` reaches the refusing arm | Reuse | C1 `stdin:` names it; no change |
| Escape fixture | 0025's `evidence/spikes/a1-deadline/` S2 (grandchild in group) | never escaped the group; macOS has no `setsid` binary | Extend | FX-deadline-escape adds a `Setsid: true` Go helper |
| Read-end close on every path | `spawn` already defers both closes (`cmdbind.go::spawn`, the `defer func() { _ = outR.Close() }()` / `errR` pair set up at pipe creation) | closes at FUNCTION RETURN, which is after the join either way | Reuse | C1 `precedence:` "after the join every read end is closed" is ALREADY satisfied by the defers; the Illustrative Code's `closeReads` is shape, not a new obligation. The defers stay — an explicit close before them would double-close, and `os.File.Close` returns `ErrClosed` on the second call rather than panicking, so the shape's cost is a swallowed error, not a crash. Nothing in this change moves the fds' release |

### Decision Rationale

Questions-Options-Criteria — the deciding rows are **correctness fit** and **prior-art
alignment**; **blast radius** breaks the tie against handing the pipes back to `Cmd`.

| Criterion | A. Bounded join, refuse (chosen) | B. Bounded join, accept partial | C. Widen kill reach | D. Pipes back to `Cmd` |
| --- | --- | --- | --- | --- |
| Correctness fit (refusal always delivered; output never silently unproven) | both: refusal within `timeout + 2·WaitDelay`; nothing unproven is parsed | refusal delivered; an envelope from a stream with a live unreachable writer is parsed anyway | whole read only if the escapee is found; a double-fork re-escapes | refusal delivered by `ErrWaitDelay`; loses release-group-first, so a correct backgrounding tool is refused |
| Prior-art alignment | matches Go `awaitGoroutines` / `ErrWaitDelay` and both peer CLIs that met this failure | no precedent found for parsing past a bound | no portable precedent; Linux-only primitives | is the Go precedent verbatim, but for `Cmd`-owned pipes only |
| Reversibility | one join site; revert is the old `drains.Wait()` | same | new syscall surface per platform | re-opens the ADV-2 / FAIL-2 ordering 76c121b fixed |
| Blast radius | `spawn` post-`Wait` block + one contract clause | same + parse-path semantics on three verbs | build tags per OS, proc-table walks | every command path's timing; intrastate#x5vq's case regresses |
| Cost | small: timer + read deadline + oracle | small, but the residue text becomes "may parse unproven output" | large; unbounded on macOS | medium; the 221 ms success case becomes a `WaitDelay` refusal |
| Portability | every Unix with pollable pipes | same | Linux only (subreaper / cgroup) | same as A |

Rationale: A is the only column that is green on both deciding rows. B is rejected on
correctness — it accepts an envelope whose stream still had a writer the CLI could not
reach, diverging from Go's disposition and from 0025's visible-over-silent stance, and a
write's read-back through such a reader would confirm a write on unproven output; the FIFO
fact that makes B *tempting* (A3) is kept as what makes A's residue text honest, not as a
licence to parse. C is rejected first as a fourth mechanism patch that still leaves the
precedence unstated, and second on portability. D is rejected on blast radius: `Cmd` signals the group
only from `Cancel`, so release-before-drain is lost and the common backgrounded-grandchild
success case pays the whole bound and is refused — the regression 76c121b measured.

Sibling-path check: the bound reuses `internal/cli/cmdbind/cmdbind.go::WaitDelay`, whose doc
already reads "bounds the stdin write and the post-kill pipe drain" — searched `internal/`
for another deadline/timer decision (`SetReadDeadline`, `time.After`, `time.NewTimer`,
`WaitDelay`, `context.WithTimeout`): none besides the `cmd.WaitDelay` assignment inside
`cmdbind.go::spawn` exists, so no parallel signal is invented. The `context.WithTimeout`
sites in `executor.go` are not siblings: each bounds the DECLARED timeout for an
invocation, which is the author's policy, where this bound is the CLI's own post-exit
mechanism cost — a different decision, as D-naming states.

Premortem: hardened (hardened) — critic ledger P-1…P-14 in `evidence/propose-premortem/critic.md`; folded: stdin named in C1 (P-1, A7), held-state not elapsed-time refusal via read deadline (P-2), exit status + remediation in Err (P-3, P-7), one timer for both drains and per-invocation bound (P-8), Go precedent restated as caller-decides with the peers named (P-9), shared-timer naming made deliberate (P-10), stale figure dropped (P-11), EPIPE residue (P-12), stderr-tail wording (P-13), Alt 2 reason reordered (P-14); F5/F6 added (P-4, P-6)
Ground-sweep: clean (14 anchors, 66 sub-facts CONFIRMED; two cosmetic wording notes folded) — a fresh-context checker read every cited `path::Symbol`, the Go stdlib anchors, 0025:C4/F4/A1, commit 76c121b and the two peer-CLI anchors without the justifying prose
Joint-check: fired → 0028 (home: JDR 0003 §D1) — symmetric to cli/0028:C3's arm-2 fire on `execution_failure`, homed as a constraint: the sub-reason rides `Detail`, `Err` is executor-facing, the applied sense is `Applied()` (JDR 0003 §D1, cited not restated); C1's `refusal:` line aligned. Recorded at the batch propose of 2026-08-31 after this record's own check ran clear: all three arms run. Arm 1 (`index --anchor-intersect`, repo-resolved): `overlaps[]` empty. Arm 2 (`index --literal-intersect`): `overlaps[]` empty. Arm 3 (absence, manual): this proposal withdraws no refusal — it converts a withheld answer into one — so no refusal token is absent; the wording it stops honouring ("necessary-and-sufficient", "true EOF", `WaitDelay`) was grepped across every Draft and Implemented peer at depth 1: only 0025 carries "necessary-and-sufficient" / `WaitDelay`, and it is this record's declared Overrides target (0025:C4, Implemented — never edited; the coupling rides to 7.1); no Draft peer names any. Unproposed roster seeds 0027 (C5 deny-list authority surface) and 0028 (in-process `edit` write carrier) were compared on problem statement: neither owns or binds the drain precedence; 0028's read-back through a command reader inherits C1 transitively without locking anything about it

## Alternatives Considered

### Alternative 1: Bounded join that accepts partial output

**Description**: Bound the join as in the chosen approach, but on the non-deadline path hand
whatever was read to `Reader.parse` / the gate verdict mapping; refuse only when the deadline
expired. Rests on the FIFO fact (A3) that a blocked drain has already consumed the direct
child's bytes.

**Pros**:

- A leaky-but-correct tool (one that backgrounds a `setsid` daemon and answers) keeps working.
- Smallest visible behaviour change for authors.

**Cons**:

- Parses an envelope from a stream that still had a writer the CLI could not reach; the
  completeness of the answer is asserted from pipe semantics, never observed.
- No prior art parses past a bound: Go returns `ErrWaitDelay` on a success-status child.
- A write's read-back through such a reader confirms the write on unproven output — the
  silent class 0025:F5/F6 exist to keep visible.

**Reason for rejection**: correctness row of the matrix — the residue would have to read
"output of unproven completeness may be parsed", which is a wider silent class than the
process leak it replaces.

### Alternative 2: Widen the kill reach so the whole read always terminates

**Description**: keep the unbounded join and make the writer set reachable — session-level
kill by walking the process table, `PR_SET_CHILD_SUBREAPER` so orphans re-parent to the CLI,
or a cgroup per invocation.

**Pros**:

- Preserves C4's "true EOF, nothing truncated" wording literally.

**Cons**:

- No portable primitive: subreaper and cgroups are Linux-only; macOS has neither; a
  process-table walk races a double-forking escapee.

**Reason for rejection**: it is a fourth mechanism at the seam that still never states the
precedence — the Seam Lineage is exactly the record of that pattern; portability is
secondary (the seam already forks by build tag, but 0025:C4's platform clause admits every
Unix with process groups, so a Linux-only widening cannot be the contract).

### Alternative 3: Hand the pipes back to `exec.Cmd`

**Description**: use a non-`*os.File` `cmd.Stdout`/`cmd.Stderr` so `Cmd`'s own copy
goroutines run under its `WaitDelay`, and take `ErrWaitDelay` as the refusal.

**Pros**:

- The Go precedent verbatim; no owned join at all.

**Cons**:

- `Cmd` signals the group only from `Cancel` (context end), so on the success path the group
  is not released before the copy goroutines are awaited: a non-escaping backgrounded
  grandchild costs the full `WaitDelay` and the correct answer is refused as `ErrWaitDelay`
  — the ordering regression 76c121b's manual pipe ownership was introduced to avoid.
- Re-opens the ordering that the manual `os.Pipe` ownership was introduced to fix.

**Reason for rejection**: blast-radius row — regresses intrastate#x5vq's case on every
wrapper-script tool.

### Briefly Rejected

- **Status quo plus documentation of the hang as residue**: a withheld refusal gives an agent caller nothing to branch on; 0025:F4 never admitted it.
- **Abandon the drain goroutines from the executor on a timer**: leaks a goroutine and two fds per invocation and still never learns whether output was whole.
- **Run the child under a pty**: one line discipline to hang up, but it changes the byte stream the envelope parser reads (line editing, CR translation) and merges stdout with stderr.
- **A separate, model-declarable drain bound**: the author declared `timeout`; the drain bound is mechanism cost, and a second knob invites the parallel-signal defect the sibling-path check exists to catch.

## Context

### Background

Found by the roborev triage of RDR 0025's implementation (`src:roborev`, batch
`rdr-0025`) and reproduced at HEAD through the real `cmdbind.Reader.Read`: a Go helper
spawning a `SysProcAttr{Setsid: true}` grandchild that inherits the pipes and sleeps
60 s, driven with a declared `timeout = "2s"`, was still blocked after 15 s — deadline
(2 s) and `WaitDelay` (500 ms) long expired.

Reproduction caveat worth carrying: macOS ships no `setsid` *binary*, so a shell fixture
using `setsid sh -c …` does **not** escape the process group — `reapGroup` kills it, the
drains close, and the read returns bounded in ~200 ms. A first probe of that shape
wrongly showed "no hang". The escape must be made with a real `setsid(2)` call.

Impact: a CLI hang with no refusal, reachable only under the `--allow-commands` opt-in
(0025), by any child whose descendants detach from the group while holding the pipes.
Constraint: the manual `os.Pipe` ownership was introduced deliberately (76c121b, fixing
the sequential-drain deadlock ADV-2 / FAIL-2) so that `Cmd.Wait` cannot truncate an
in-flight read; a repair that hands the pipes back to `Cmd` re-opens that deadlock.
Closing one C4 liveness hole opened another, which is why the precedence must be
decided once rather than patched a third time (see Seam Lineage).

The triage comment on intrastate#936e names the fork the successor RDR must weigh
(timed join accepting partial output; deadline-capable pipe reads; session-level
reaping) and the correlated clauses (F4's residue restated as honest leakage; an
`FX-deadline` scenario in which a `setsid(2)` grandchild holding both pipes must still
refuse within `timeout + bound`). None of that is decided here.

### Technical Environment

Go, `os/exec` with `SysProcAttr` process-group handling behind build tags
(`internal/cli/cmdbind`, Unix-only syscalls tagged after intrastate#1drn; Windows takes
C4's platform refusal). `internal/cli/cmdbind/cmdbind.go`: `Reader.Read`, the
manually-owned `os.Pipe` drains, `drains.Wait()`, `reapGroup` (`syscall.Kill(-pid,
SIGKILL)`). Output caps in play: the 1 MiB stdout cap and 4 KiB stderr tail that C4's
"partial output" wording would have to be read against. Governing records: 0004
(accessor execution safety model), 0025:C4 / 0025:F4, and 0025's implementation
artifacts (`deviations.md`, `verification.md`), which do not record the `os.Pipe`
ownership change against C4.

## Research Findings

### Investigation

Prior art was read before any approach was named (`evidence/research/prior-art.md`). The
instance the code is built on, Go `os/exec`, already decided this conflict for the pipes it
owns: `Cmd.WaitDelay` "bounds the time spent waiting on … a child process that exits but
leaves its I/O pipes unclosed", `(*Cmd).awaitGoroutines` "forcibly closes their pipes and
returns ErrWaitDelay", and `ErrWaitDelay` is returned "if the process exits with a successful
status code but its output pipes are not closed" ⇒ bound wins, the parent closes, success
with an open pipe is an error. The two comparable Go CLIs that met this exact failure
(the peer checkout set was swept for `WaitDelay`; gh-cli, helm, goreleaser and kubebuilder
set none) both took that disposition unchanged: roborev maps `exec.ErrWaitDelay` to the
context error; beads sets a 10 s `WaitDelay` because "a grandchild … that inherited the
output pipes would otherwise keep Wait blocked indefinitely after the kill". The class
frame is Kerrisk *TLPI* §44.3 — EOF is a property of the whole writer set, and "the reader
sees end-of-file (once it has read any outstanding data in the pipe)". Code paths read:
`cmdbind.go::spawn` (post-`Wait` order, owned `os.Pipe` ends, `drains.Wait()`),
`cmdbind.go::reapGroup`, `procgroup_unix.go::killGroup`, `cmdbind.go::readBounded`,
`executor.go::invokeRead` and `(*Executor).Gate` (deadline-first classification); Go
`os/file_unix.go::newFile` (`kindPipe` ⇒ pollable) and `internal/poll` (`Close` evicts
pending I/O). Constraint carried from 76c121b: release the group BEFORE joining the drains,
which is what keeps the manual pipe ownership.

### Key Discoveries

- **Documented** — Go's own precedence for this conflict is "bound wins; parent closes;
  success with an open pipe is `ErrWaitDelay`" (`os/exec/exec.go` `WaitDelay` /
  `awaitGoroutines` docs, quoted in `evidence/research/prior-art.md`).
- **Documented** — `Cmd.WaitDelay` reaches only pipes `Cmd` creates; `spawn`'s `*os.File`
  ends are neither copied nor closed by `Wait`, so the owned join is outside it today.
- **Measured** — `os.Pipe` ends are pollable on Unix (`os/file_unix.go::newFile`), and a
  pending `Read` is ended from another goroutine in 33-102 µs by deadline (A2). The spike
  also settled WHICH deadline: Go checks it before attempting the syscall, so only a
  FUTURE deadline runs the final read that delivers buffered bytes — a deadline of `now`
  skips them. `Close` unblocks too but makes them unrecoverable, which is why the
  mechanism is the deadline.
- **Documented** — the executor classifies `timeout` from its own bounded context before it
  reads the binding's error (`executor.go::invokeRead` "Deadline first"), so the binding
  needs no timeout class (A4 checks the write and read-back paths).
- **Measured** — the total bound `timeout + 2·WaitDelay` (A6) holds with 418-487 ms of
  headroom, but the drain-join leg alone runs 500.1-508.2 ms against a 500 ms `WaitDelay`
  in every sample: the total is the claim, the legs are not separately bounded. The
  "500 ms is enough for a killed in-group grandchild" premise (A5) is confirmed with a
  worst case of 27.7 ms across 160 trials — 5.5% of the bound.
- **Measured** — a blocked drain past the bound has consumed the direct child's whole
  output (A3): confirmed byte-for-byte at 1 KiB and 63 KiB with the drain reading only
  after the child was gone. The rule holds for an EXITED child; past the ~64 KiB pipe
  buffer a stalled drain leaves the child blocked in `write(2)` and unable to exit, which
  is why C1 `whole:` now states the concurrent-drain antecedent.

## Trade-offs

### Consequences

- Positive: every command invocation returns within `timeout + 2·WaitDelay`; an agent caller
  always has a refusal class to branch on.
- Positive: the ordinary paths (child closes its pipes; in-group backgrounded grandchild) are
  byte-for-byte unchanged — the bound is reached only when the group signal could not reach
  a writer.
- Negative: C4's "true EOF, nothing truncated" is narrowed to "whole within the bound"; a
  tool that answers correctly but leaves a `setsid` daemon holding its stdout is now refused
  rather than accepted — visibly, with the pipe named.
- Negative: the leaked process residue is now explicitly admitted with its output; nothing
  in the CLI reaps a `setsid(2)` escapee.
- Negative: a command WRITE that applied its effect and exited 0 but left a held stdout is
  refused `execution_failure`; the applied sense that lets a caller
  choose read-back over blind retry is `Applied()` (JDR 0003 §D1; the executor-side
  derivation is A8) — the same class 0025 already accepts for a write killed
  at its deadline.

### Risks and Mitigations

- **Risk**: the join's timer fires but the drain goroutine does not return (A2 wrong on one
  OS). **Mitigation**: the A2 spike runs on darwin and linux before the join is written; the
  design never returns from `spawn` with a drain still running.
- **Risk**: 500 ms is too short under CI load and the in-group grandchild case regresses
  (A5). **Mitigation**: the existing grandchild-leak oracle re-run N times under `-race`;
  the bound is a constant, so widening it is a one-line follow-up, not a contract change.
- **Risk**: a success-path child whose helper legitimately holds the pipe briefly past
  500 ms is refused. **Mitigation**: the refusal names the held pipe and the bound, so the
  author sees exactly what to fix (close or redirect the helper's inherited fds).

### Failure Modes

- **F1** Visible: an escaped grandchild holds a pipe past the bound — the invocation refuses
  `timeout` (deadline elapsed) or `execution_failure` (deadline not elapsed) with `Err`
  naming the pipe(s) and the bound; diagnosis: the refusal text plus a lingering process
  visible to `ps` holding the CLI's former pipe.
- **F2** Visible: a correct tool that backgrounds a `setsid` daemon inheriting stdout is
  refused where it used to hang; recovery: the tool closes or redirects the daemon's fds.
  The adjacent risk — an IN-group helper refused because it missed the bound — is bounded
  by measurement, not assumption: A5 timed 160 trials at a worst 27.7 ms against 500 ms.
  The tail is load-sensitive (idle→loadavg 47 grew the max ~15x while the median grew ~3x),
  so the margin is a function of runner oversubscription rather than a constant; at that
  scaling the bound absorbs roughly another 18x of tail growth. A far more oversubscribed
  CI runner is the condition that would erode it, and the symptom would be intermittent
  `execution_failure` on a reader that works locally.
- **F3** Silent risk: a process outside the group outlives the refusal and keeps running
  (the restated 0025:F4 residue); this RDR admits it and does not reap it.
- **F4** Silent risk: if A3 is wrong, bytes the direct child wrote can be lost at the bound;
  the refusal still fires, so the loss is bounded to a refused invocation, never a parsed one.
- **F5** Silent risk: a pipe whose poller registration failed (`os.ErrNoDeadline` at the
  pollability check) cannot be deadlined; on that host the join stays unbounded — the old
  hang, confined to a non-pollable host and logged at the check rather than discovered.
  Scope, from the A2 spike: this trigger is **unexercised for a real `os.Pipe`** — every
  read end accepted a deadline (nil) on darwin and linux, and no naturally non-pollable
  pipe was constructible. The error value itself is confirmed on the two shapes that do
  produce it (a regular file; `os.NewFile` over a dup'd pipe fd, adopted without the
  original's poller registration), so S5 simulates the condition at the check rather than
  reproducing it. The residual risk is a host whose pipes are non-pollable — not observed
  on either supported OS.
- **F6** Visible but under-explained: a held pipe that also outran the deadline reports
  `timeout` with no `Err`/`Detail` (0025:C4's rule for the class), so the held-pipe reason
  is visible only on `execution_failure`; the class flips with the declared `timeout` for the
  same helper, as it does today for any slow command near its deadline.

## Implementation Plan

### Prerequisites

- [x] All Critical Assumptions verified (A2 and A6 by spike on darwin and linux before any code — `evidence/spikes/a2-unblock-read/`, `a3-a6-drain-bound/`, `a5-ingroup-margin/`)
- [ ] A `Setsid: true` Go test helper exists in `internal/cli/cmdbind` so the escape is real on both OSes (Background: a shell fixture does not escape on macOS)

### Minimum Viable Validation

Fixture **FX-deadline-escape**, deadline shape, run through the real `cmdbind.Reader.Read`
under `--allow-commands` with a declared `timeout = "2s"`:

1. The child (a Go helper) spawns a grandchild with `SysProcAttr{Setsid: true}` that
   inherits stdout and stderr and sleeps 60 s; the child itself writes nothing and sleeps
   past the deadline.
2. `flow resolve` invokes the reader; the context deadline expires at 2 s; `Cancel` kills
   the group; `Wait` returns; `reapGroup` signals the group; the drains are joined under the
   bound.
3. Expected end-state: the read returns within `2 s + 2·500 ms` (Normative: elapsed
   ≤ 3 s, asserted with margin), the executor classifies `timeout`, and the refusal's
   `Detail` is empty as 0025:C4 states for `timeout`.
4. The grandchild is still alive (asserted — the admitted residue) and the test kills it.
5. Success shape, same fixture: the child writes a whole envelope and exits 0 while the
   grandchild holds the pipes — expected `execution_failure` within `2·500 ms` of the exit,
   `Err` naming "stdout, stderr" and "exited 0", and the envelope NOT parsed (no values
   returned).
6. Late-but-unheld control (Normative): a child that writes 1 MiB to stdout and exits with
   no grandchild, run with the drain goroutine artificially delayed past the timer, returns
   the whole value — the final read under `DrainGrace` drains buffered bytes and sees EOF;
   no refusal. This row is what distinguishes the two candidate mechanisms: with a deadline
   of `now` the read short-circuits before the syscall and this control FAILS (A2, measured
   on both OSes), so it is the MVV's guard on C1 `precedence:`, not a formality.

Measured basis (A6, `evidence/spikes/a3-a6-drain-bound/`): worst end-to-end 2.013 s
(darwin) / 2.082 s (linux) against the 2.5 s budget — 487 ms and 418 ms of headroom. The
drain-join leg alone measured 500.1-508.2 ms in 40/40 samples, so the assertion above is
on the TOTAL, which is what C1 commits, and never on a per-leg figure the evidence
contradicts.

### Phase 1: Bounded join

Bound the drain join in `spawn` by `WaitDelay`; on expiry the parent ends its own read ends
(mechanism from the A2 spike), joins the drains, and refuses through `wrap` with the held
pipe(s) named. No other line of `spawn` moves.

### Phase 2: Contract re-citation

Re-point `cmdbind.go`'s C4 comments ("necessary and sufficient", "run to a true EOF",
`WaitDelay`'s doc) at 0026:C1; 0025 itself is not edited.

### Phase 3: Fixtures and oracles

Add FX-deadline-escape (both shapes) beside the existing grandchild-leak and cancel oracles;
re-run the in-group grandchild oracle under the bound (A5) so the prior point-fixes stay
proven.

### Phase 4: Surface

State the restated residue where authors read about `timeout` for command entries, so the
"a process may leak and its output is never read" consequence is documented, not discovered.

## Validation

### Testing Strategy

Done means: every command path refuses within its bound, and no held stream is ever
parsed. The MVV fixture is the deadline case; these are the rest of the matrix, all in
`internal/cli/cmdbind` unless noted.

1. **Scenario**: Escaped grandchild holds the pipes, on the read, gate and write verbs
   (write covering both its invocation and its read-back) — A1's claim that one join
   covers every path.
   **Expected**: each refuses within `timeout + 2·WaitDelay`; no envelope parsed, no
   gate verdict derived, no write confirmed on held output; on BOTH write legs
   `Applied()` is true AND the CLI renders the "may have been applied" text. The two
   legs assert it against different arms, and the scenario keeps them apart: the DIRECT
   invocation refuses `execution_failure` and is what proves the re-key landed —
   `flow_exec.go::accessorFailureOf`'s `ClassExecutionFailure` arm carries no phase
   check and sets no `Detail` today (A8). The READ-BACK leg refuses
   `read_back_incomplete` and is a REGRESSION assertion, not a proof of the gain: that
   arm already sets `detailMayHaveApplied` and the executor already sets `applied`, so
   the row fails only if the re-key broke a path that was correct (C1 `refusal:`).
2. **Scenario**: The existing in-group grandchild-leak and cancel oracles re-run under
   the bounded join, N times under `-race` and CPU load (A5).
   **Expected**: unchanged results, elapsed far below the bound — the prior point-fixes
   stay proven and no correct tool is newly refused.
3. **Scenario**: Ordinary paths — child closes its pipes; output at the 1 MiB stdout cap
   and 4 KiB stderr tail boundaries.
   **Expected**: byte-for-byte unchanged from today (C1 `whole:`).
4. **Scenario**: Non-reading child with a stdin payload past the pipe buffer (A7).
   **Expected**: `exec.ErrWaitDelay` from `Wait` reaches `spawn`'s non-`ExitError` arm
   and refuses `execution_failure` after `reapGroup` and the join.
5. **Scenario**: Non-pollable pipe — the creation-time pollability check returns
   `os.ErrNoDeadline` (F5), simulated at the check.
   **Expected**: the condition is logged at the check, not discovered at the join; the
   fallback is the documented unbounded path, and the log names the host condition.
6. **Scenario**: Held stdout with the deadline also elapsed (F6).
   **Expected**: classified `timeout` with empty `Detail` per 0025:C4 — asserted so the
   class flip is a tested consequence rather than a surprise.

7. **Scenario**: payload past the ~64 KiB pipe buffer (1 MiB), drain prompt vs
   artificially stalled — the antecedent C1 `whole:` states (A3).
   **Expected**: prompt drain returns the whole 1 MiB; a stalled drain leaves the child
   blocked in `write(2)`, killed by the group signal, with exactly one pipe buffer
   (65536 bytes) recovered as a correct prefix. The test exists so a later refactor that
   serializes or delays the drains fails here rather than silently voiding A3.
8. **Scenario**: a child that RESISTS termination (ignores SIGTERM, survives to its
   `WaitDelay`) while a grandchild holds the pipes — the case A6 did not exercise, and
   the one where leg (a) spends its allowance instead of donating it to leg (b).
   **Expected**: the invocation still returns within `timeout + 2·WaitDelay + 100 ms`
   (the premortem's P-8 tolerance), so the total bound is asserted where the two legs
   genuinely compete rather than only where one underspends (0025:C6).
9. **Scenario**: a LATE drain with no holder — the drain goroutine stalled past the
   timer with bytes buffered and every writer gone (premortem P-2).
   **Expected**: accepted, whole output, no refusal — the final read under `DrainGrace`
   delivers the buffered bytes and then sees EOF. Companion: an escaped writer emitting
   a byte every 10 ms forever is REFUSED (never blocked, never EOF). Together these pin
   "refuse on no-EOF-after-a-final-bounded-read, not on elapsed time"; with a deadline of
   `now` the first of them fails, which is what makes `DrainGrace` load-bearing.

Coverage goal: the FX-deadline-escape helper is exercised on darwin and linux in CI, so
the escape is real on both (Background).

## Finalization Gate

> Complete each item with a written response in
> `{ARTIFACT_DIR}/gate.md` before marking this RDR as
> **Final**. Written responses prevent rubber-stamping
> and produce a review record.
>
> First run the mechanical pre-sweep
> (`prompts/gate/tooling-pass.md`): TEMPLATE section
> coverage, Method-label vocabulary, `Source Search`
> self-reference, `Docs Only` on load-bearing claims. It
> catches what the review rounds disturbed; resolve any
> BLOCK before the written responses.
>
> At lock, replace Contradiction Check, Assumption
> Verification, Scope Verification and Proportionality
> with the one-line pointer to gate.md — those four
> judge THIS record at THIS lock and no peer cites
> them. **Cross-Cutting Concerns stays here**, below
> the pointer: it names the project-wide policy other
> RDRs conform to, so it must stay projected and
> citable as `cli/NNNN:G-cross-cutting`. Cite it that
> way, not by section name.

### Contradiction Check

[Gate key: contradiction — a gate response is cited as
`cli/NNNN:G-<key>`, so the key is a stable id and is
not derived from this heading, which may be reworded.]

[State any conflicts between Research Findings and
the Proposed Solution. If none exist, state
"No contradictions found between research findings,
design principles, and proposed solution."]

### Assumption Verification

[Gate key: assumptions]

[Confirm every Critical Assumption Evidence Record
is internally consistent: Status, Method, and
Evidence agree, and "If wrong" is non-empty. List
any record whose Method is `Docs Only` (these block
lock unless paired with a Spike or Source Search
plan) and any that remain `Pending` or `Unverified`
with a plan to verify before implementation begins.
Confirm no `Verified` stamp is self-referential or
proves only an adjacent claim, and that each cited
`path::Symbol` resolves on `main`. **Status
consistency:** no assumption marked `Pending` or
`Unverified` may have settled-fact prose elsewhere in
the RDR depending on it.]

### Scope Verification

[Gate key: scope]

[Confirm the Minimum Viable Validation is in scope
and will be executed during implementation, not
deferred. State the specific test or proof.]

### Cross-Cutting Concerns

[Gate key: cross-cutting]

[Retained at lock — this sub-section stays in the RDR
when the other gate responses move to gate.md, because
peer RDRs cite it as `cli/NNNN:G-cross-cutting` and an
element that is not projected cannot be cited.]

[List only concerns that apply to this RDR. For each,
state either how this RDR addresses it, or which peer
RDR owns the project-wide policy this RDR conforms
to. Omit (rather than N/A-bullet) anything that does
not apply.]

Candidate concerns (include only those that apply):
versioning · build tool compatibility · licensing ·
deployment model · IDE compatibility · incremental
adoption · secret/credential lifecycle · memory
management · concurrency model · character encoding ·
canonical-form / determinism (see note below).

If this RDR claims byte-identical output,
content-addressed identity, or replay-stable hashes,
also confirm: hash function + library, pre-image
byte layout, primitive encodings, map iteration order,
whitespace policy, case folding, empty/null/absent
distinguishability, and a version marker for future
evolution.

### Proportionality

[Gate key: proportionality]

[Is the document right-sized for the change? Flag
any sections that should be trimmed before locking.
The split test is **contract count, not word count**:
confirm this RDR is the sole author of at most one
independent load-bearing contract (per the Normative
Contracts split signal). If it owns more than one
seam, flag it for splitting rather than locking the
seams together.

Re-validate the **Profile** Metadata field against the
contracts you just counted: confirm the value Resolve
wrote still matches (one contract + no user-facing
surface → `small`; etc. per the applicability matrix).
If the lenses that actually ran disagree with the
Profile (e.g. Profile says `small` but the change locks
a contract that warranted `mid`+ lenses, or the lenses
were skipped on a wrong `small`), correct the field and
do not lock until the missing lenses have run. This is
the latch's backstop — a wrong Profile cannot route
past the lens battery undetected. A `Transient`-marked
contract with a named deleting sibling and schedule is a
recorded lifespan disposition, not an under-sized
Profile — do not count it when re-deriving. Also confirm form:
value + one clause naming the contract(s); strip any
matrix/provenance prose left from the template or Seed
(it belongs in the template comment, not the instance).]

## References

- 0025:C4 (deadline triple, `detail:` channel, platform clause) and 0025:F4 (the residue this RDR restates); 0025:A1's spike `evidence/spikes/a1-deadline/` (in-group grandchild case S2)
- Go `os/exec/exec.go` — `Cmd.WaitDelay`, `ErrWaitDelay`, `(*Cmd).awaitGoroutines`; `os/file_unix.go::newFile`, `os/file.go` `SetReadDeadline`; `internal/poll/fd_poll_runtime.go::(*pollDesc).evict`
- Kerrisk, *The Linux Programming Interface*, §44.3 (pipe EOF requires every write descriptor closed; outstanding data delivered first)
- Peer CLIs swept for `WaitDelay`: roborev `internal/agent/process.go::configureSubprocess`; beads `internal/storage/dolt/store.go::cliExecWaitDelay`
- `internal/cli/cmdbind/cmdbind.go::spawn`, `::reapGroup`, `::readBounded`, `::WaitDelay`; `internal/cli/cmdbind/procgroup_unix.go::killGroup`; `internal/accessor/executor.go::invokeRead`
- Commits at the seam: 76c121b (manual `os.Pipe` ownership, concurrent drains, reap on success), cd9a09d, d22cd7a, b4b7431; intrastate#936e (this defect), intrastate#1drn, intrastate#x5vq, intrastate#878e
- `evidence/research/prior-art.md` (queries, accepted quotes, rejected branches)
- JDR 0003 §D1 (`docs/jdr/0003-accessor-binding-seam.md`) — the `execution_failure` sub-reason carrier; C1's `refusal:` line is this record's instance
