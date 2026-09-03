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
- **Profile**: foundational — accretion floor (3 prior point-fixes at the C4 seam, Seam Lineage below); the contract axis alone reads mid: one locked contract amended (0025:C4's execution-safety triple) on a user-facing surface (whether a `timeout` refusal is delivered at all).
- **Priority**: High
- **Related Issues**: intrastate#936e (defect tracker; stays open until Stage 8); intrastate#v0hb (closed, RDR 0025's tracker); closed point-fixes at the seam: intrastate#1drn, intrastate#x5vq, intrastate#878e
- **Predecessors**: 0025-command-invoking-accessor-bindings, 0004-accessor-execution-safety-model
- **Overrides**: 0025:C4 (the "necessary and sufficient" deadline triple — to be amended to a bounded drain with a stated precedence) and 0025:F4's residue class (leakage of a process, never loss of the refusal)
- **Seam Lineage**: `internal/cli/cmdbind::reapGroup` / `internal/cli/cmdbind::drains.Wait` (`area:internal-cli`) — 3 prior closed point-fixes; trail: cd9a09d + intrastate#1drn (process-group syscall build tags), d22cd7a + intrastate#x5vq (grandchild-leak oracle), b4b7431 + intrastate#878e (cancel oracle); plus 76c121b, the manual `os.Pipe` ownership change inside 0025's own implementation (ADV-2 / FAIL-2), which is the change that opened this gap. The count is to be confirmed, not re-derived, at Resolve.

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
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: to be taken at Resolve — the callers of `spawn` (`cmdbind.go::Reader.Read`, `cmdbind.go::Gate.Gate`, `cmdbind.go::(*Writer).Apply`) and a sweep for any other `drains`/`cmd.Wait` site under `internal/cli/cmdbind`.
  - **If wrong**: a path outside `spawn` keeps the unbounded join and the hang survives on that verb only — visible as the same stuck invocation, on gates or writes. The same search confirms every caller reads `err` before `inv` (a gate fast path on `inv.exitCode` would parse a held stream) and that no arm of `spawn` returns between `Wait` and the join.
- **A2 [The joining goroutine can END a pending blocking `Read` on an `os.Pipe` read end from another goroutine — `SetReadDeadline` or `Close` on the read end — on darwin and linux, and the drain goroutine then returns promptly with the bytes it had read]**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: Stage-2 read only: `os/file_unix.go::newFile` marks `kindPipe` files `pollable`, `os/file.go` documents `SetReadDeadline` on pollable files, and `internal/poll/fd_poll_runtime.go::(*pollDesc).evict` ("unblocking any I/O running on fd") runs from `(*FD).Close`. Whether the drain actually returns, how fast, and with which error, is an order/reachability claim: spike `{SPIKE_DIR}a2-unblock-read/` on both OSes, holding the write end in a `setsid(2)` grandchild.
  - **If wrong**: the join's timer fires but the drain goroutine never returns; if `spawn` then returns anyway it leaks the goroutine and the fd per invocation, and if it waits the hang is unchanged. The spike also records what `SetReadDeadline` returns on a pipe whose poller registration failed (`os.ErrNoDeadline`), which is F5's trigger.
- **A3 [After `Wait` + `reapGroup`, a drain still blocked on an EMPTY pipe past the bound has already consumed every byte the exited direct child wrote — pipe FIFO delivers outstanding data before EOF — so the only writer left is a process outside the group, and refusing is a policy choice, not loss of the direct child's answer]**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: Stage-2 read only: Kerrisk, *TLPI* §44.3 ("the reader sees end-of-file (once it has read any outstanding data in the pipe)"). The spike must show the parent's buffer holding the child's complete envelope at the moment the bound expires, on both OSes, with a child that writes its envelope and then exits leaving an escaped writer.
  - **If wrong**: C1's residue wording ("its output is never read") understates the loss — the direct child's own bytes can be lost at the bound — and the Failure Modes must say so.
- **A4 [The executor's deadline-first classification (`internal/accessor/executor.go::invokeRead`, and its Gate and Apply counterparts) turns the new refusal into `timeout` whenever the declared timeout elapsed, on every path, so `spawn` needs no timeout class of its own and the binding's error is read only when the deadline did not expire]**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: `invokeRead` ("Deadline first: … `timeout` is the one reported") and `(*Executor).Gate` were opened at Stage 2; the write path's `appliedDeadline` check and the read-back's own executor call are to be opened at Resolve, since a read-back that times out inside a WRITE must still report as the write's refusal.
  - **If wrong**: a bound-expired drain after the deadline is reported as `execution_failure` with a held-pipe reason instead of `timeout`, so a caller branching on class is misled.
- **A5 [500 ms (`cmdbind.go::WaitDelay`) is long enough that a NON-escaping backgrounded grandchild SIGKILLed by `reapGroup` releases its pipe ends within the bound on CI, so the ordinary "wrapper script backgrounds a helper" success case (intrastate#x5vq's oracle) pays nothing and is never refused]**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: the existing grandchild-leak oracle in `internal/cli/cmdbind` re-run under the bounded join, N times under `-race` and CPU load, asserting the value is returned and elapsed stays far below the bound.
  - **If wrong**: correct tools whose helpers take longer to die are refused as held-pipe failures — a regression of the prior point-fix, seen as intermittent `execution_failure` on a working reader.
- **A6 [The invocation returns within `timeout + 2·WaitDelay` in the escaped-grandchild case: `Cmd.Wait` returns within one `WaitDelay` of the context deadline (its own timer bounds the stdin-copy goroutine and the post-kill wait) and the drain join is bounded by one more]**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: fixture **FX-deadline-escape** — a Go helper spawning a `SysProcAttr{Setsid: true}` grandchild that inherits both pipes and sleeps 60 s, in two shapes (child also sleeps past the deadline; child exits 0 with a whole envelope), timed end to end on darwin and linux; the spike also confirms the escape is real on both (macOS ships no `setsid` binary, so a shell fixture does not escape — Background).
  - **If wrong**: C1's stated bound is false; a caller sizing its own watchdog from `timeout` still waits longer than promised.
- **A7 [The child's stdin is already bounded and its expiry already refuses: `cmd.Stdin` is a `bytes.Reader`, so `Cmd` owns that pipe and its copy goroutine under `cmd.WaitDelay`; a child (or escapee) that never reads it ends in `exec.ErrWaitDelay` from `Wait`, which reaches `spawn`'s non-`ExitError` arm — after `reapGroup` and the join — and refuses `execution_failure`]**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: `cmdbind.go::spawn` (`cmd.Stdin`, `cmd.WaitDelay`, the `default:` arm of the `werr` switch) and `os/exec/exec.go` `WaitDelay` doc (it bounds the copying goroutines, stdin's included — "bounds the stdin write" is 0025:C4's gloss, not Go's text); 0025's A1 spike S4 (non-reading child, 1 MiB stdin) is the in-group case.
  - **If wrong**: a write whose payload exceeds the pipe buffer hangs on a non-reading child — the third inherited pipe re-opens the hole this RDR closes for the other two.
- **A8 [The executor's write arm can derive `Applied()` from the binding's typed held-pipe error: `internal/accessor/executor.go::Write` has, or can gain, a refusal site that matches `cmdbind.HeldPipeError` with `errors.As` after the direct child was reaped and sets the unexported applied sense there, and the CLI's "may have been applied" rendering keys on `Applied()` rather than on class and phase alone (JDR 0003 §D1)]**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: to be taken at Resolve — `internal/accessor/executor.go::Write` (its refusal sites and where the unexported applied sense is set), the `Applied()` accessor and its callers, and the CLI's "may have been applied" rendering site; read against JDR 0003 §D1.
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
3. When the timer fires: the parent sets `SetReadDeadline(now)` on `outR`/`errR` (a
   pollable `os.Pipe` file wakes a pending `Read`; a read with bytes buffered still returns
   them, so a merely LATE drain finishes whole and only an empty pipe with a live writer
   reports `os.ErrDeadlineExceeded`), joins the drains (now returning), closes both read
   ends, and — only if a drain reported the deadline — returns
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
precedence: a bounded wait outranks a whole read. Order in `spawn` is fixed: Wait (direct child reaped) → reapGroup (one kill(2) to -pgid, non-blocking) → drain join under ONE timer of WaitDelay covering BOTH read drains. When the timer fires the parent sets a read deadline of now on its own read ends; each drain's FINAL read then reports the pipe's state, not the clock: bytes still buffered are delivered, EOF (every writer gone) is whole output as today, and a deadline error is a HELD pipe — empty, with a writer the group signal could not reach. Only a held pipe refuses. After the join every read end is closed (deadline to unblock, close to release the fd)
refusal:  a held pipe refuses through `*accessor.ExecError`, this record's instance of JDR 0003 §D1: `Detail` carries the held-pipe reason — the pipe(s) held ("stdout"/"stderr"/both), the bound, the direct child's exit status (exited N / killed by signal) and the remediation ("close or redirect the helper's inherited stdio") — composed AHEAD of the stderr tail collected up to the bound; `Err` wraps the typed held-pipe error `cmdbind.HeldPipeError` (pipes held, exit status), which the executor matches with `errors.As` at the refusal site. Applied sense: a held pipe on the write path (the write invocation, or its read-back after the write ran) after the direct child was reaped is a post-run refusal — class stays `execution_failure` and `Applied()` is true, carried per JDR 0003 §D1 through the executor's write arm (which refusal site sets it, and whether the CLI's "may have been applied" rendering already keys on `Applied()` rather than on class and phase, is A8); on read and gate invocations `Applied()` is false as today. The returned invocation carries no stdout; stdout read on a held drain is never parsed on ANY path — read envelope, gate verdict, write, write read-back (JDR 0003 §D1)
class:    the class is the executor's, as today: ctx `DeadlineExceeded` ⇒ `timeout` (deadline-first rule; per 0025:C4 a `timeout` refusal carries no Err/Detail, so the held-pipe reason survives only on `execution_failure`), else the ExecError ⇒ `execution_failure`. The binding never classifies `timeout` itself
stdin:    unchanged and named: the child's stdin is a Cmd-owned pipe (`cmd.Stdin` is a `bytes.Reader`), bounded by Cmd's own WaitDelay; a non-reading child or escapee holding its read end ends in `exec.ErrWaitDelay` from Wait, which `spawn`'s non-ExitError arm refuses as `execution_failure` AFTER reapGroup and the join — no arm of `spawn` returns between Wait and the join
bound:    WaitDelay (500 ms) is the ONE bound — the same constant, by design, for every "the CLI waits past the child's exit" case: Cmd's stdin write and post-Cancel wait, and the owned drain join; no second knob, not model-declarable. Cancel sends one signal and returns, so Cmd's timer starts at the deadline. Per invocation, the return bound is timeout + 2·WaitDelay; a command write is two invocations (write, then read-back), each under its own bound
whole:    a drain that reaches EOF — before the timer, or on its final read after it — is whole, as today; the 1 MiB stdout cap and 4 KiB stderr tail are unchanged. A child (and group) that closes its pipes observes no change
residue:  0025:F4 restated — a process outside the child's process group may outlive the refusal; it keeps write ends whose reader is gone, so its next write fails with EPIPE (SIGPIPE under the default disposition), and whatever it wrote is never read. Leakage of a process, or of THAT process's output, is admitted; a withheld refusal is not. Diagnosis: the refusal's Detail names the held pipe(s); the leaked process is visible to `ps` until its next write
```

#### Load-Bearing Decisions

- **Naming** — the bound is `cmdbind.go::WaitDelay`, reused; rejected: a new `DrainBound`
  constant (a parallel knob for one decision), and a model-declarable drain bound (the
  author declared `timeout`; the drain bound is the CLI's own mechanism cost, not policy).
  The reuse is deliberate, not incidental: Cmd's stdin/post-Cancel timer and the drain
  timer are two timers sharing one value because they answer one question ("how long past
  the child's exit does the CLI wait"), so a future tuning moves both — accepted.
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
    deadlineNow(outR, errR)     // wakes an EMPTY-pipe read; buffered bytes still drain
    drains.Wait()               // now returns; each drain reports EOF or deadline
}
closeReads(outR, errR)          // release the fds on every path
if held := heldPipes(); len(held) != 0 {
    return inv, wrap(tail(stderr), errHeldPipe(held, exitStatus))
}
// EOF (before or after the timer): whole output, unchanged from here on
```

Illustrative: `waitBounded`/`deadlineNow`/`closeReads`/`heldPipes`/`errHeldPipe` are
names for the shape, not the contract; the timings in the MVV are Normative.

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| One spawn site for every command path | `internal/cli/cmdbind/cmdbind.go::spawn` | none known (A1 confirms) | Extend | C1 lands once; no per-binding wait |
| Drain bound constant | `internal/cli/cmdbind/cmdbind.go::WaitDelay` | today applied only to `Cmd`'s own goroutines | Reuse | C1 `bound:` — no second constant |
| Group release before join | `internal/cli/cmdbind/cmdbind.go::reapGroup` / `procgroup_unix.go::killGroup` | cannot reach a `setsid(2)` escapee | Reuse | C1 `residue:` names the escapee as the admitted leak |
| Bounded, overflow-detecting read | `internal/cli/cmdbind/cmdbind.go::readBounded` | blocks in `Read` until EOF or an ended fd | Reuse | the join ends the read from outside; `readBounded` itself is unchanged |
| Refusal channel with stderr tail | `*accessor.ExecError` via `cmdbind.go::wrap` (0025:C4 `detail:`) | none | Reuse | the held-pipe reason rides in `Err`, tail in `Detail` |
| Deadline-first classification | `internal/accessor/executor.go::invokeRead` (+ Gate / Apply) | A4 checks the write and read-back paths | Reuse | C1 `class:` — no binding-side `timeout` |
| Stdin bound | `cmd.Stdin` = `bytes.Reader` + `cmd.WaitDelay` (Cmd-owned pipe) | A7 confirms `ErrWaitDelay` reaches the refusing arm | Reuse | C1 `stdin:` names it; no change |
| Escape fixture | 0025's `evidence/spikes/a1-deadline/` S2 (grandchild in group) | never escaped the group; macOS has no `setsid` binary | Extend | FX-deadline-escape adds a `Setsid: true` Go helper |

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
`WaitDelay`): none besides the `cmd.WaitDelay` assignment inside `cmdbind.go::spawn` exists, so no
parallel signal is invented.

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
- **Documented** — `os.Pipe` ends are pollable on Unix (`os/file_unix.go::newFile`), so a
  pending `Read` can be ended from another goroutine by deadline or close; that it does so
  promptly with the bytes read is **Assumed** (A2, spike).
- **Documented** — the executor classifies `timeout` from its own bounded context before it
  reads the binding's error (`executor.go::invokeRead` "Deadline first"), so the binding
  needs no timeout class (A4 checks the write and read-back paths).
- **Assumed** — the total bound `timeout + 2·WaitDelay` (A6) and the "500 ms is enough for
  a killed in-group grandchild" premise (A5) are order claims: spike and oracle, not quotes.
- **Assumed** — a blocked drain past the bound has consumed the direct child's whole
  output (A3): TLPI states the FIFO rule; the implementation must be shown to observe it.

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
- **F3** Silent risk: a process outside the group outlives the refusal and keeps running
  (the restated 0025:F4 residue); this RDR admits it and does not reap it.
- **F4** Silent risk: if A3 is wrong, bytes the direct child wrote can be lost at the bound;
  the refusal still fires, so the loss is bounded to a refused invocation, never a parsed one.
- **F5** Silent risk: a pipe whose poller registration failed (`os.ErrNoDeadline` at the
  pollability check) cannot be deadlined; on that host the join stays unbounded — the old
  hang, confined to a non-pollable host and logged at the check rather than discovered.
- **F6** Visible but under-explained: a held pipe that also outran the deadline reports
  `timeout` with no `Err`/`Detail` (0025:C4's rule for the class), so the held-pipe reason
  is visible only on `execution_failure`; the class flips with the declared `timeout` for the
  same helper, as it does today for any slow command near its deadline.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified (A2 and A6 by spike on darwin and linux before any code)
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
   the whole value — the deadline read drains buffered bytes and sees EOF; no refusal.

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
   gate verdict derived, no write confirmed on held output; on the write path
   `Applied()` is true (C1 `refusal:`).
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
