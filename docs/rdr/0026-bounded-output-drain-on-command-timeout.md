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
  - **Population caveat**: what these 160 trials measure is scheduler and pipe-teardown latency for a Go helper that dies promptly to SIGKILL, on a 12-core machine with no cgroup quota. They do NOT measure process-EXIT cost, which is the quantity the bound actually has to cover: a large-heap runtime unmapping its arena, a helper blocked in uninterruptible sleep on NFS or FUSE (premortem P-7), or any process in a CPU-quota-throttled container, where wall-clock latency is not a function of loadavg at all. Oversubscription (loadavg 47.5 on 12 cores) is the only stress dimension exercised, and F2 bounds the residual risk with this same 27.7 ms figure scaled only by that dimension. The 5.5%-of-bound headroom is therefore a headline for the measured population, not for the population that will fail.
  - **If wrong**: correct tools whose helpers take longer to die are refused as held-pipe failures — a regression of the prior point-fix, seen as intermittent `execution_failure` on a working reader. The likely first report is CI-only or container-only: works on the author's machine, fails in a throttled pod, class says "could not be executed."
- **A6 [The invocation returns within `timeout + 2·WaitDelay` in the escaped-grandchild case — the TOTAL is the claim: `Cmd.Wait` returns well inside one `WaitDelay` of the context deadline (its own timer bounds the stdin-copy goroutine and the post-kill wait) and donates the remainder to the drain join, which spends one full `WaitDelay` plus a scheduling epsilon. The two legs are not separately bounded]**
  - **Status**: Verified, with the per-leg claim narrowed
  - **Method**: Spike
  - **Evidence**: `evidence/spikes/a3-a6-drain-bound/`, fixture **FX-deadline-escape** (Go helper, `SysProcAttr{Setsid: true}` grandchild inheriting both pipes), both shapes, n=10 each, timeout 1.5 s against `WaitDelay = 500` ms (`cmdbind.go:66`), budget 2.5 s. **The escape is real on both OSes** — darwin child pid/pgid 37707/37707 vs grandchild 37709/37709; linux 828/828 vs 834/834; "grandchild alive after `kill(-pgid, SIGKILL)`" true on both, so the drain genuinely sees no EOF. **The total bound holds with wide margin**: worst total 2.013 s (darwin) / 2.082 s (linux) against 2.5 s — 487 ms and 418 ms of headroom. Leg (a), deadline→`Wait`, holds comfortably at 1–79 ms, two orders of magnitude inside one `WaitDelay`.
  - **Narrows C1 `bound:`**: leg (b), `Wait`→join, measured **500.1–508.2 ms in all 40 samples — never under 500 ms**. This is structural, not noise: the escaped grandchild holds the pipes so the join always waits the full timer, then pays timer-fire latency plus goroutine scheduling to observe the deadline and unwind. So "the drain join is bounded by one more `WaitDelay`" is strictly false; the true per-leg bound is one `WaitDelay` plus a small scheduling epsilon (worst overshoot 8.2 ms linux, 3.0 ms darwin). The TOTAL claim survives unharmed because leg (a) consumes almost none of its allowance and absorbs the overshoot many times over. Scope caveat: these figures are for a child that dies promptly to the group kill; a child that RESISTS termination would push leg (a) toward its full `WaitDelay`, and only then does the budget tighten — that case was not exercised. **The total is therefore Normative on an unexercised composite.** Leg (b)'s measured floor is at or above the whole per-leg allowance in 40/40 samples, so the total holds only while leg (a) underspends; the one case where it does not underspend is the resisting child, and that is exactly the case with no measurement. S8 is the test and it is UNRUN — until it runs, C1 `bound:`'s total is verified for the prompt-death shape only, and a resisting child plus an escaped grandchild is the shape most likely to breach it in production first. S8 runs before Phase 1 closes, not at the end; if it breaches, the fix is the stated scheduling tolerance or the constant, not a relaxed assertion.
  - **If wrong**: C1's stated bound is false; a caller sizing its own watchdog from `timeout` still waits longer than promised.
- **A7 [The child's stdin is already bounded and its expiry already refuses: `cmd.Stdin` is a `bytes.Reader`, so `Cmd` owns that pipe and its copy goroutine under `cmd.WaitDelay`; a child (or escapee) that never reads it ends in `exec.ErrWaitDelay` from `Wait`, which reaches `spawn`'s non-`ExitError` arm — after `reapGroup` and the join — and refuses `execution_failure`]**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `cmdbind.go:207` assigns `cmd.Stdin = bytes.NewReader(stdin)` (non-`*os.File`) and `:224` sets `cmd.WaitDelay = WaitDelay * time.Millisecond` (`WaitDelay = 500`, `:66`). Go's own text covers the stdin goroutine without the gloss: `os/exec/exec.go` WaitDelay doc — "those pipes are closed in order to unblock any goroutines currently blocked on Read or **Write** calls" — and the mechanism is in code, not only doc: `exec.go::childStdin:542` hands an `*os.File` stdin straight to the child with no goroutine, but the `bytes.Reader` branch (`:545-561`) appends the write end to `c.parentIOPipes` and the `io.Copy` to `c.goroutine`, which `awaitGoroutines` force-closes at `exec.go:999` before returning `ErrWaitDelay` (`:1002`). `ErrWaitDelay` is not an `*exec.ExitError`, so `errors.As` at `cmdbind.go:311` fails and control reaches the `default:` arm at `:320-322` — after `reapGroup` (`:297`) and after the join (`:302`) — refusing `execution_failure`. The third pipe does not re-open the hole.
  - **If wrong**: a write whose payload exceeds the pipe buffer hangs on a non-reading child — the third inherited pipe re-opens the hole this RDR closes for the other two.
- **A9 [The per-read deadline re-arm bounds the idle gap, not the tail: re-arming `now + DrainGrace` before each read of the final drain lets a 1 MiB buffered tail drain to EOF through a ~64 KiB pipe buffer without a false held report, on darwin and linux]**
  - **Status**: Verified, with the witness corrected
  - **Method**: Spike
  - **Evidence**: `evidence/spikes/a9-a10-regrace/` (`fixture.go`, `a9_regrace.go`, `output.txt`), darwin/arm64 go1.26.6 and linux/arm64 go1.26.8 under `golang:1.26`, drain started 200 ms late, `DrainGrace = 50` ms, completeness by `bytes.Equal` against an independently reconstructed `ENVELOPE-BEGIN:`/`:ENVELOPE-END` framed payload so a truncation anywhere would show. **Confirmed on both OSes**: with the deadline re-armed before each read, a 1 MiB paced tail recovered 1048576/1048576 bytes ending on EOF in 12/12 trials, 0/12 deadline errors, 32 reads. The decisive pair — **max per-read gap 26.8 ms (darwin) / 31.6 ms (linux), under the 50 ms grace, while total elapsed ran 612 ms / 663 ms = 12.2x / 13.3x the grace** — is the claim itself: the grace bounded the gap and not the tail. No flake in any configuration, either OS, across 48+ trials per OS.
  - **Corrects its own witness (scope, not a refutation)**: the shape this assumption originally named — a 1 MiB BUFFERED tail — produced a NULL result and does not discriminate the mechanisms. A burst tail drains in 1.6 ms (darwin) / 9.5 ms (linux) across 32 non-waiting reads, so a single absolute deadline recovers it whole too; both arms pass. Duration, not size, is the operative variable: the spike added a paced child (32 chunks, 20 ms idle gaps = 0.4x grace each, summing to ~12x grace), under which the single absolute deadline failed 0/12 complete, truncating at exactly 163840 of 1048576 bytes (15.6%) after 5 reads at ~51 ms with terminal `os.ErrDeadlineExceeded` — a false held report on a live writer, identical on both OSes. Even a 1 KiB paced payload loses 62%. MVV row 6 was amended to carry both shapes because the burst shape alone would pass a wrongly-implemented single-absolute drain; C1 `precedence:` now reads "size or total duration".
  - **If wrong**: MVV row 6 is non-deterministic — a slow-arriving answer is refused as a held pipe on a loaded machine, which is a false refusal on a pipe with no writer at all.
  - **Known sensitivity (recorded, not blocking)**: the measured per-read gap (27-32 ms) sits under 2x `DrainGrace` (50 ms), so a heavily loaded host or a slower writer cadence could push a LEGITIMATE gap past the grace and produce exactly the false held report the re-arm exists to prevent. This is a tuning question about the VALUE of `DrainGrace` — it argues for a comfortable grace, never for the single-absolute mechanism, which fails the same case far worse. Revisit the constant if S7/S9 flake under CI load.

- **A10 [The pollability probe is safe and detecting: a `SetReadDeadline` on each read end at pipe creation, immediately cleared, changes no subsequent read behavior on a pollable pipe, and returns `os.ErrNoDeadline` on a non-pollable one]**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `evidence/spikes/a9-a10-regrace/` (`a10_probe.go`, `output.txt`), darwin/arm64 go1.26.6 and linux/arm64 go1.26.8, identical on both OSes on every point. The check is NEW — a repo-wide sweep of `internal/` for `SetReadDeadline`/`SetDeadline`/`SetWriteDeadline`/`ErrNoDeadline` returns ZERO hits, production or test. **DETECTING**: both non-pollable shapes return `"file type does not support deadline"` with `errors.Is(err, os.ErrNoDeadline)` true — a regular file, and `os.NewFile` over a `syscall.Dup` of a pipe read end; an ordinary `os.Pipe` read end returns nil on both arm and clear. The dup'd-pipe shape is the load-bearing one: the fd is a genuine PIPE, so the probe detects poller REGISTRATION rather than file kind, which is the property F5 needs. Both non-pollable files still read correctly after the failed probe (21 bytes / 15 bytes), so a failed probe costs nothing but the information it returns. **SAFE**: the strong test, not the round-trip — after probe+clear, a blocking `Read` on an EMPTY pipe with a LIVE writer stayed parked 12/12 at 300 ms (a leaked 10 ms probe deadline would have fired 30x sooner), then woke on the write with correct bytes 12/12, 0 errors; worst wake latency 56 µs (darwin) / 799 µs (linux) against an unprobed baseline of 45 µs / 332 µs. Full-payload round-trip `bytes.Equal` 12/12 in every cell at 1 KiB / 64 KiB / 1 MiB, with and without the probe, both OSes.
  - **If wrong**: a probe on the hot path perturbs ordinary reads — the byte-for-byte-unchanged claim (C1 `whole:`, S3) fails for every command, not just the held case.

- **A11 [The bounded join preserves the happens-before edge the unconditional `drains.Wait()` supplies today: a drain goroutine ended by the grace deadline still runs `drains.Done()`, so the parent's join returns only after both goroutines have finished writing the `stdout`/`stderr` slices, and no read of `inv.stdout`/`inv.stderr` races a live drain]**
  - **Status**: Pending — DOWNGRADED at Stage 6, survivable into implementation
  - **Method**: MVV Test
  - **Evidence**: added by the critique pass (D-8). Today `drains.Wait()` (`cmdbind.go:302`) is the ONLY edge between the goroutine writes at `:275-279` and the parent read at `:304`; the bounded join replaces the unconditional wait, so the edge must be re-established rather than inherited. A2 measured a parent `SetReadDeadline` against a blocked `Read` on an `os.File`, whose poller is internally synchronised — that covers the unblock, not the slice handoff. The MVV fixture cannot detect a violation: its child is silent, so each drain is blocked in its FIRST read and no partial slice write is ever in flight. Verify by running S3 and S7 (a 1 MiB payload, so the drain is mid-write when the timer could fire) under `-race`.
  - **Stage 6 disposition — DOWNGRADED, not deferred by choice**: this assumption is UNVERIFIABLE before implementation by construction, and that is why it is not run here. Its named check is S3 + S7 under `-race`; S7 requires the drain-start stall seam, which A12 confirms does NOT exist in `cmdbind` today and which THIS record adds as new production code. There is no bounded join to race until Phase 1 writes one, so no spike can stand in for the check — the race is a property of the code this RDR authors, not of the platform (which is what a spike can answer, and what A2/A9 already did). NOT MVV-critical in the hard-rule sense: the MVV provably cannot exercise this edge (silent child, drain blocked in its first read), which is precisely why the assumption names S3/S7 instead of an MVV row. **Named plan, binding**: Phase 1 closes only when S3 and S7 pass under `-race -count=25`; a `-race` report on the slice handoff is a Phase-1 blocker, not a follow-up. **Survivable because** the failure is loud and bounded — `-race` flags the handoff, or the parent proceeds past a live drain — and the repair is local (the join stays a real join; C1 `precedence:` (b) already forbids the failing design by name).
  - **If wrong**: `-race` flags the handoff intermittently, or a drain's final bytes are lost or read torn — and in the worst shape the parent proceeds past a live drain that never returns, which is the original hang wearing a timer that reads as correct.
- **A12 [The drain-start stall seam S7/S9/MVV row 6 share is a NEW production-code addition, and an unexported package-level hook nil in production is enough to build it: `cmdbind` has no existing test-injection surface to hang it on, and adding one perturbs no ordinary invocation]**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: added by the repeatability pass (G-6); confirmed at Stage 6 against the tree. S7 previously asserted the seam was "injected through the package's existing test surface"; there is no such surface. The package's non-test files hold exactly five package-level vars, and only one is an injection seam: `goos` (`cmdbind.go:97`), whose own doc scopes it to making the refusal testable on Unix CI and which is read only at `:104` (`Unsupported()`) and `:180` — unrelated to drain timing. The others are data or sentinels: `AllowlistedVars` (`:87`), `procGroupPlatforms` (`procgroup_unix.go:14`, `procgroup_other.go:14`), `errNoProcGroup` (`procgroup_other.go:27`). `export_test.go` exposes only `SetGOOSForTest` (`:15`), `ProcGroupPlatforms` (`:26`) and `ProcGroupSupported`/`UnsupportedOn` (`:35`) — no timing hook. Today's timing tests drive real subprocesses, confirming the gap: `adversarial_0025_test.go:67,96,134` (`exec sleep 30`), `:264,320` (`( sleep 20 ) &`), `invocation_0025_test.go:1131`, backed by `fixtures_0025_test.go:181,190`. **The shape is viable**: the two drain goroutines are inline closures launched at `cmdbind.go:273` and `:277`, each opening `defer drains.Done()`, so an unexported package-level `var drainStartHook func()` checked immediately after that defer inserts at `:274`/`:278` with no signature change, no build tag and no production sleep — and it sequences correctly ahead of `cmd.Wait()` (`:281`), `reapGroup` (`:295`) and `drains.Wait()` (`:300`). The only `//go:build` lines in the package are the two-file unix/other process-group split (`procgroup_unix.go:1`, `procgroup_other.go:1`); `procgroup_unix.go:16-20` states that only the mechanism is tag-split, so a third tag would be the novelty this record forbids.
  - **If wrong**: three scenarios that share one seam need three fixtures, or the seam lands as a build tag / production sleep — the two shapes S7 rules out by name.
- **A8 [The executor's write arm can derive `Applied()` from the binding's typed held-pipe error: `internal/accessor/executor.go::Write` has, or can gain, a refusal site that matches the accessor-side held-pipe error with `errors.As` (C1 `refusal:` — a `cmdbind` type cannot be named here; the import runs the other way) after the direct child was reaped and sets the unexported applied sense there; the CLI's "may have been applied" rendering must GAIN its `Applied()` key on the DIRECT write arm, since the shipped mapping keys on class and phase alone and that one case slips past it — the read-back leg already carries both the applied sense and the rendering (JDR 0003 §D1 rule 3, which assigns that gain to this record)]**
  - **Status**: Verified
  - **Method**: Source Search
  - **Stage 6 confirmation**: every site below was re-checked against the tree and all hold, with no line drift. The two decisive facts are confirmed: `flow_exec.go`'s `ClassExecutionFailure` arm (`:410-412`) is a bare `return envErr(codeAccessorFailed, id, "the accessor `"+id+"` could not be executed")` — no `at == phaseWrite` check, no `ce.Detail` assignment, no `ce` local at all — so the DIRECT write arm genuinely lacks the applied sense; and `Applied()` has ZERO non-test callers (23 `.Applied()` hits across `internal/`, every one in a `_test.go`), so the CLI derives the applied sense from class+phase and never from the accessor. `Write` sets `r.applied = true` at exactly six sites, no seventh (`executor.go:363,375,400,409,423,472`). One refinement to the plan below: the `errors.As` match need not be written inline in the `err != nil` arm — `refusalOf` (`executor.go:74`) is the shared refusal constructor and already runs `errors.As(err, &ee)` for `*ExecError` at `:83`, and every refusal site funnels through it, so the held-pipe match has an existing home. The `err != nil` arm itself (`:367-370`) is a plain guarded block with `err` in scope and `errors` already imported, so either placement compiles with no structural change. Import direction confirmed clean: `go list -deps ./internal/accessor` returns only `internal/resolve` and `internal/table` as intrastate-internal deps and no `internal/cli` package, so a typed-error read at the CLI boundary introduces no cycle.
  - **Evidence**: the setter is reachable — `internal/accessor/model.go:314` holds the unexported `applied bool` whose doc reserves it for the write path, `:324` the `Applied()` accessor; `executor.go::Write` already sets `r.applied = true` at six refusal sites (`:363`, `:375`, `:400`, `:409`, `:423`, `:472`) by direct in-package field assignment, so the proposed `errors.As` match takes the `if err != nil` arm at `:367-371` with no structural change. The rendering is the gain, not a found fact: `internal/cli/flow_exec.go::accessorFailureOf` (`:395`) keys `detailMayHaveApplied` on class plus phase at `:399-405` (`ClassTimeout` + `at == phaseWrite`) and `:419-424` (`ClassReadBackIncomplete`) and never calls `Applied()` — which has zero non-test callers repo-wide. Decisively, `ClassExecutionFailure` at `:410-412` carries no phase check and sets no `Detail` of its own, so a held-pipe write on the DIRECT invocation renders today with no applied sense. The held-pipe REASON still reaches the user one frame up — `flow_exec.go::accessorFailure` wraps every mapping in `withStderrTail(…, refusal.Detail)`, which lands the tail in the one `Detail` slot — so what the gain adds is the applied-sense sentence ahead of that tail, not the diagnostic itself. The read-back leg is NOT part of the gain: `ClassReadBackIncomplete` at `:419-424` already sets `detailMayHaveApplied`, and `executor.go`'s read-back arm already sets `applied`, so the re-key is one site (the `err != nil` write arm), not two. JDR 0003 §D1 rule 3 records exactly this ("the shipped mapping keys on class and phase, which this case slips past") and lands the `Applied()` derivation here, so the gain is a decided constraint this record implements, not an open question. Its cost is carried in C1 `refusal:` and Testing Strategy S1.
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

Determinacy: fired — C1 `precedence:` (step ordering in `spawn`, and the total order
overflow > held > whole the drain's terminal-condition report must carry).

**C1**

```normative
precedence: a bounded wait outranks a whole read. Order in `spawn` is fixed: Wait (direct child reaped) → reapGroup (one kill(2) to -pgid, non-blocking) → drain join under ONE timer of WaitDelay covering BOTH read drains. When the timer fires the parent sets a read deadline of now + DrainGrace (50 ms) on its own read ends; each drain's FINAL read then reports the pipe's state, not the clock: bytes still buffered are delivered, EOF (every writer gone) is whole output as today, and a deadline error is a HELD pipe. The predicate is the drain's REPORTED terminal condition, not the byte count: a drain that ends on `os.ErrDeadlineExceeded` reports held however many bytes it delivered during the grace, and one that ends on EOF reports whole. "Empty, with a writer the group signal could not reach" describes the common case; it is not the test, because a writer trickling bytes slower than the grace is exactly the case that must refuse (S9's companion) and is never empty. Overflow outranks held on the same drain: a drain past `StdoutCap` refuses the existing overflow error whatever its terminal condition, since the bytes were read and the pipe's state is then not what the invocation turns on. The held fold takes the two drains' REPORTED conditions as its input — it derives nothing from the byte counts or from closed-over state — and returns the held pipe names in the fixed order stdout, stderr, so `len(held) != 0` is the consulted predicate and the single-pipe form is the bare name (S10, MVV row 5). The grace is what MAKES that final read happen: Go checks the deadline BEFORE attempting the syscall (`internal/poll/fd_unix.go::(*FD).Read` calls `prepareRead` ahead of `syscall.Read`), so a deadline of NOW short-circuits with n=0 and reports a held pipe without ever looking at the pipe. The deadline is re-armed at `now + DrainGrace` BEFORE EACH read of the final drain, so the grace bounds the IDLE GAP between reads and never the size or total duration of the tail. That re-arm has NO site today and this clause creates one: `readBounded` (`cmdbind.go:352`) is `io.ReadAll(io.LimitReader(r, limit))` followed by `io.Copy(io.Discard, r)`, so `io.ReadAll` owns the loop and there is no per-read point to re-arm. Implementing this clause replaces that body with an explicit `for { SetReadDeadline(now+DrainGrace); Read }` loop — a rewrite of the function, not the signature change alone. The terminal condition is carried as ONE three-valued report per drain, not as a pair of booleans and not as a bare `error`: the three states are mutually exclusive and totally ordered here (overflow > held > whole), and a shape that cannot express that order — two independent booleans, or an `error` in which overflow is not an error — forces the caller to re-derive the precedence this clause fixes. `readBounded` returns that report alongside its bytes and the drain records it for `heldPipes()` to fold; the identifier and the Go spelling are the implementer's, the three-valued-and-ordered shape is not. The rewrite must also preserve three behaviors the existing suite pins: the read stays bounded at `StdoutCap+1` so exactly-at-cap is a value and cap+1 is a failure (the overflow check at `cmdbind.go:326` reads `len(inv.stdout) > StdoutCap`), the discard-remainder exit must itself carry the deadline, and the pre-timer path must stay deadline-free. S3 asserts that preservation; A9 measures the loop. A re-implementation that shifts the cap boundary is a contract violation, not a refactor: a drain making progress keeps its deadline ahead of it and runs to EOF, and only a gap longer than the grace ends it as held. A single absolute deadline would instead bound the whole remaining tail against the clock, which turns a slow-arriving tail into a false held-pipe report on a pipe with no writer at all — the failure MVV row 6(b) exists to catch, measured at 15.6% truncation on both OSes (A9); a large tail already BUFFERED does not trip it, which is why row 6(b)'s paced shape, not row 6(a)'s burst, is the discriminating witness — the elapsed-time refusal the premortem rejected (P-2: refuse on "no EOF after a final bounded read", not on elapsed time), and, one layer down, the guessed absence 0004:D-absent-vs-unreadable forbids. Verified on darwin and linux (A2). Only a held pipe refuses. The handoff is SYNCHRONISED and the edge is named, because removing the unconditional join removes today's only one: `drains.Wait()` (`cmdbind.go:302`) is currently the sole happens-before edge between the drain goroutines' writes to the `stdout`/`stderr` slices (`:275-279`) and the parent's read of them (`:304`), and a bounded join that can return while a drain is still live voids it. Two edges are therefore required, not one: (a) the parent's deadline-set on the read ends is the mark — `os.File`'s own poller is internally synchronised, so `SetReadDeadline` from the joining goroutine against a `Read` in flight is the sanctioned cross-goroutine call (A2 measured exactly this) and no separate `graceOn` flag is introduced; (b) the parent still MUST NOT read `inv.stdout`/`inv.stderr` until every drain goroutine has returned, so the join remains a real join — the bound governs how long the parent waits before setting the deadline, never whether it waits for the goroutines to finish. A drain that ends on the deadline still runs `drains.Done()`, so the join completes; a design in which the parent proceeds past a live drain is forbidden here by name. The MVV cannot catch a violation of (b) — its fixture's child is silent, so the drain is blocked in its FIRST read and no partial write to the slice is in flight — so S3 and S7 run under `-race` and that is the check. No read end outlives `spawn` (deadline to unblock, close to release the fd). The existing defers at pipe creation already satisfy this — they fire at function return, which is after the join but also after the `werr` switch and the overflow check, and that is the guarantee: no fd leaks, not that the fd is released at the join point. An explicit close before them would double-close, so none is added and the Illustrative Code draws none
refusal:  a held pipe refuses through `*accessor.ExecError`, this record's instance of JDR 0003 §D1: `Detail` carries the held-pipe reason — the pipe(s) held ("stdout"/"stderr"/both), the bound, the direct child's exit status (exited N / killed by signal) and the remediation ("close or redirect the helper's inherited stdio") — composed AHEAD of the stderr tail collected up to the bound. That exit status is composed from fields `spawn` already populates from `werr` — `inv.exited` and `inv.signaled` beside `inv.exitCode` — so no new field is added to carry it; `Err` wraps a typed held-pipe error carrying the pipes held and the exit status, which the executor matches with `errors.As` at the refusal site. The type is declared in `internal/accessor`, NOT in `cmdbind`: `cmdbind` imports `accessor` (`cmdbind.go:44`) and `accessor` imports no `cli` package (`go list -deps ./internal/accessor` returns `internal/resolve`, `internal/table` only), so a concrete `cmdbind` type named at an `errors.As` site inside `executor.go` is an import cycle that does not compile. It follows the carrier already in use: `accessor.ExecError` (`model.go:423`), which `executor.go:82` already matches with `errors.As`; the held-pipe error is its sibling on the same side of the edge, and `cmdbind` populates it as it already populates `*accessor.ExecError`. Being matched from `executor.go` and populated from `cmdbind`, both outside its own package, the type is EXPORTED and its held-pipes and exit-status fields with it — an unexported type would be unnameable at the `errors.As` site, which is the same constraint that puts the declaration in `accessor` rather than `cmdbind`. Applied sense: a held pipe on the write path after the direct child was reaped is a post-run refusal — `Applied()` is true, carried per JDR 0003 §D1 through the executor's write arm. The write's TWO legs already differ in class and this clause does not flatten them. The DIRECT write invocation refuses `execution_failure`, and it is the leg that needs the gain: `executor.go::Write`'s `err != nil` arm sets no applied sense today, and `flow_exec.go::accessorFailureOf`'s `ClassExecutionFailure` arm carries neither a phase check nor a `Detail`, so a held-pipe write slips past the "may have been applied" rendering (JDR 0003 §D1 rule 3 lands that gain here). The READ-BACK leg after the write ran is already whole and gains nothing: a held pipe there returns the write's own refusal re-classed `read_back_incomplete` with `applied` already set (`executor.go`'s `raw.class != ""` arm), which `accessorFailureOf`'s `ClassReadBackIncomplete` arm already renders with `detailMayHaveApplied`. So one site changes, not two — the read-back leg is named here because a reader who assumes the class is uniform across both legs would re-key an arm that is correct. On read and gate invocations `Applied()` is false as today. The returned invocation carries no stdout; stdout read on a held drain is never parsed on ANY path — read envelope, gate verdict, write, write read-back (JDR 0003 §D1). "Carries no stdout" is a CALLER obligation, not a structural guarantee, and is stated as one so a later caller cannot inherit it by accident: `spawn` builds `inv` before the `werr` switch (`cmdbind.go:304`) and both error arms return that populated value (`:322`, `:326`), so the field is readable on every error path. The obligation is that every caller checks `err` before touching `inv` — true of today's three callers by inspection (A1), and NOT enforced by any type or signature. A new caller that reads `inv.stdout` or `inv.exitCode` ahead of `err` derives a verdict from a stream whose live writer is unreachable, which is the silent-wrong-answer class 0025 exists to prevent. Implementation carries this as a comment at the `inv` construction site naming the two error returns; making it structural (an error-only return, or zeroing `inv` on the refusal arms) is a `readBounded`-adjacent change this record does not take and 0028's command-reader read-back leg should revisit
class:    the class is the executor's, as today: ctx `DeadlineExceeded` ⇒ `timeout` (deadline-first rule; per 0025:C4 a `timeout` refusal carries no Err/Detail, so the held-pipe reason survives only on `execution_failure`), else the ExecError ⇒ `execution_failure`. The binding never classifies `timeout` itself
stdin:    unchanged and named: the child's stdin is a Cmd-owned pipe (`cmd.Stdin` is a `bytes.Reader`), bounded by Cmd's own WaitDelay; a non-reading child or escapee holding its read end ends in `exec.ErrWaitDelay` from Wait, which `spawn`'s non-ExitError arm refuses as `execution_failure` AFTER reapGroup and the join — no arm of `spawn` returns between Wait and the join
bound:    WaitDelay (500 ms) is the ONE bound — the same constant, by design, for every "the CLI waits past the child's exit" case: Cmd's stdin write and post-Cancel wait, and the owned drain join; no second knob, not model-declarable. DrainGrace (50 ms) is NOT a second bound: it answers a different question (how far ahead the read deadline must sit for the kernel to deliver already-buffered bytes), is never model-declarable, and cannot move the refuse/accept boundary — it is a mechanism constant inside the one bound, named here because every shipped constant of this family is named in its contract (0025:C4 does the same for StdoutCap and StderrTailCap). Cancel sends one signal and returns, so Cmd's timer starts at the deadline. Per invocation, the COMMITTED return bound is timeout + 2·WaitDelay, asserted with a scheduling tolerance of +100 ms (the premortem's P-8 figure) — stated here once so the scenarios cite it rather than each carrying its own: a return past the total but inside the tolerance is a pass, and past the tolerance is a contract violation. The bound is PER INVOCATION, and a command write runs up to THREE of them through the same reader, not two: the pre-write baseline read when `protectedKeys` returns non-empty (`executor.go:320`), the write itself (`:355`), and the read-back (`:397`). A write journey against an escaping helper therefore costs up to `3·(timeout + 2·WaitDelay)` end to end, and a caller sizing a watchdog from a single invocation's bound is killed mid-write. Worse, the baseline leg's refusal is DISCARDED (`:326` sets `baselineUnread` and continues), so that leg spends a full bound and yields no signal — it is the one leg whose held-pipe refusal the user never sees. This record does not change that discard (it is 0004:C13's deliberate "unverifiable, not unconstrained" arm); it names the cost so the published total is not read as the journey's. The two legs are NOT separately bounded and the contract does not claim they are: the drain join runs one full WaitDelay plus timer-fire and goroutine-scheduling latency (measured 500.1-508.2 ms in 40/40 samples, A6), which the deadline leg's unused allowance absorbs. That absorption is why the total holds, so it is stated rather than left as a coincidence between two independent decisions (0025:C6's rule); a child that RESISTS termination spends leg (a)'s allowance and is the case where the total tightens — S8 tests it
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
| Held-vs-EOF-vs-overflow terminal condition | each drain goroutine, reporting out ONE three-valued report (C1 `precedence:`) | `heldPipes()` at the join, which takes both reports as input | one join in `spawn` | none — `readBounded` returns only `[]byte` today (`cmdbind.go:352-364`) and reports nothing | the drain's reported condition; NOT the clock, NOT the byte count (C1 `precedence:`); overflow outranks held (D-selection-predicate), which is why the report is one ordered value and not two booleans |
| Stderr tail into `Detail` | `cmdbind.go::wrap` | `withStderrTail` composes into one slot | one slot, applied-sense first | held-pipe reason and tail share the slot | `refusal.Detail`, ordered per 0025:C4 |
| Held-pipe error type | declared in `internal/accessor` beside `ExecError` | `executor.go::Write` via `errors.As`; `cmdbind` populates it | one declaration, one match site | `accessor.ExecError` is the existing sibling | `internal/accessor` — `cmdbind` imports `accessor`, never the reverse (C1 `refusal:`) |
| Drain bound | `cmdbind.go::WaitDelay` | `spawn`'s join timer, Cmd's own timer | one constant, two timers | `DrainGrace` is an OFFSET, not a rival bound | `WaitDelay`; no second knob (D-naming) |

**`oracle` — test-discriminability.** The rows that could pass by absence-of-error
carry their named failing control.

| MVV / scenario | Fails if X is wrong because Y | Negative / failing control |
| --- | --- | --- |
| MVV row 3 (deadline shape) | asserts elapsed ≤ 3 s under the C1 `bound:` tolerance AND class `timeout` AND empty `Detail`; a withheld refusal never returns | today's code hangs unboundedly — the bug itself is the control |
| MVV row 5 (success shape) | asserts `execution_failure` within `2·WaitDelay` under the same C1 `bound:` tolerance, `Err` naming "stdout, stderr" and "exited 0", envelope NOT parsed | a build that parses held stdout returns values and fails the not-parsed assertion |
| MVV row 6 (late-but-unheld) | 6(a) burst: a deadline of `now` short-circuits before the syscall and loses the buffered bytes. 6(b) paced: a SINGLE absolute `now + DrainGrace` bounds the whole tail against the clock and truncates it | **this row IS the control, and it takes BOTH shapes** — 6(a) fails under a deadline of `now` (A2) but passes under either grace mechanism, so it cannot discriminate them; 6(b) is what fails the single-absolute variant (A9: 0/12 complete, truncated at 15.6%, vs 12/12 whole re-armed, both OSes) |
| S10 (stderr held alone) | the both-pipes case passes without ever exercising the single-pipe `Detail` form or the tail the grace exists to preserve | stdout-held-alone: the same refusal with the other pipe named, and no tail |
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
| Final drain after the timer, writer gone | whole remaining tail, however many reads it takes — the grace bounds the idle gap between reads, not the tail (C1 `precedence:`, A9) | none; a tail that stalls longer than the grace is held, not truncated |
| Drain past the pipe buffer, STALLED | exactly one pipe buffer (65536 bytes) as a **correct prefix** | the remainder — admitted, and the reason S7 pins the concurrency (A3, D-the-drains-stay-concurrent) |
| Held pipe | stdout is never parsed on ANY path | the whole stream — refused, not truncated-then-read |

**`disposition` — input class → outcome.**

| Input class | Exit / class | Event / error | Artifact minted | Silent vs loud |
| --- | --- | --- | --- | --- |
| Child closes pipes, exits 0 | success | none | envelope parsed | loud (the value) |
| Held pipe, deadline NOT elapsed | `execution_failure` | the held-pipe error (accessor-side, C1 `refusal:`) in `Err`, reason + tail in `Detail` | no stdout | LOUD |
| Held pipe, deadline ALSO elapsed | `timeout` | none — 0025:C4 gives `timeout` no Err/Detail | no stdout | loud in class, **quiet in reason** (F6, S6) |
| Held pipe on the DIRECT write | `execution_failure` + `Applied()` true | held-pipe reason | no read-back | LOUD once the re-key lands (A8) |
| Held pipe on the write's READ-BACK | `read_back_incomplete` + `applied` already true | read-back's own error | no confirmation | LOUD today |
| Non-reading child, stdin past the buffer | `execution_failure` | `exec.ErrWaitDelay` via `spawn`'s non-ExitError arm | none | loud (S4, A7) |
| Non-pollable pipe (`os.ErrNoDeadline`) | unbounded join — the old hang | recorded on the invocation at the check | a `Detail` line naming the host condition | **loud at the check, by design** — as a `Detail` line, not a log line; `cmdbind` has no logger (F5, S5) |
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
| 7. final read reports pipe state | C1 `precedence:` (held = the drain's reported terminal condition) | deadline error ⇒ held, whatever bytes the grace delivered |
| 8. read ends closed | C1 `precedence:`; Infra-Audit close-ownership row | defers already close them; no double-close hazard |
| 9. classify | C1 `class:`; F6 | deadline elapsed ⇒ `timeout`, empty `Detail` (MVV row 3) |
| 10. return | C1 `bound:` total | ≤ 3 s asserted; worst measured 2.013 s / 2.082 s (A6) |

No CONTRADICTION row. Steps 6 and 9 are the two that a naive reading collides:
the grace makes the final read happen, and the class still flips to `timeout`
when the deadline elapsed — so the held-pipe REASON is computed and then not
rendered. That is F6, recorded as an accepted trade-off rather than a contradiction.

#### Load-Bearing Decisions

- **Naming** — the bound is `cmdbind.go::WaitDelay`, reused; rejected: a new `DrainBound`
  constant (a parallel knob for one decision), and a model-declarable drain bound (the
  author declared `timeout`; the drain bound is the CLI's own mechanism cost, not policy).
  The reuse is deliberate, not incidental: Cmd's stdin/post-Cancel timer and the drain
  timer are two timers sharing one value because they answer one question ("how long past
  the child's exit does the CLI wait"), so a future tuning moves both — accepted, with the
  coupling's cost named rather than left implicit. The two populations want OPPOSITE moves:
  a CPU-throttled container wants a LARGER drain bound (A5's population caveat), while a
  non-reading child with a stdin payload wants a TIGHT stdin bound (A7, S4). So "a future
  tuning moves both" is not a neutral consequence — raising the constant to stop CI
  flakiness on the stdin path silently widens the held-pipe refusal window system-wide, and
  a genuinely broken command then appears to work longer before refusing. Two things follow.
  The record's own escape ("widening it is a one-line follow-up, not a contract change") is
  FALSE once C1 `bound:` publishes `timeout + 2·WaitDelay`: the constant is in the contract,
  so moving it is a contract change requiring an amendment, and it must be treated as one.
  And S4 and S8 are the pair that pins the coupling — S4 fails if the constant grows past
  what the stdin path tolerates, S8 fails if it shrinks past what a resisting child needs —
  so a future PR moving it for one population's reason breaks the other's row rather than
  passing silently. That pairing is the enforcement; the shared value is not.
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
  into `Detail` — the text that tells a user WHICH helper held the pipe. That benefit is
  real but partial, and F6 is where the limit is recorded: the tail reaches the user on the
  `execution_failure` leg (MVV row 5, S10) and NOT when the declared timeout also elapsed,
  which is the record's motivating scenario. The grace is justified on the leg that renders
  it, not on all of them. 50 ms is chosen
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
  deadline-first rule, not a new tie-break). When overflow and held both qualify on one
  drain — the bytes passed `StdoutCap` and the discard drain (`cmdbind.go:352-364`, whose
  `io.Copy(io.Discard, r)` is the third exit and today discards its error) then ended on the
  deadline: overflow wins, and the existing "stdout exceeded the 1 MiB bound" refusal after
  the join (`cmdbind.go:326-329`) stands unchanged. Both are `execution_failure`, so the
  class cannot separate them; overflow is chosen because the bytes were read, which makes
  the pipe's terminal state not what the invocation turns on. When the `werr`
  non-ExitError arm (`exec.ErrWaitDelay`, C1 `stdin:`) and a drain condition both qualify,
  the drain conditions are selected FIRST — overflow, then held — and the non-ExitError
  arm refuses only if neither fired. A held pipe names the helper that held it, which
  `ErrWaitDelay` does not, so the more diagnostic reason wins; both are
  `execution_failure`, so this orders the reason, not the class.

#### Illustrative Code

Shape only — the mechanism inside step 3 is the A2 spike's.

```go
// in each drain goroutine — the re-arm is HERE, per read, never in the parent:
//   for { r.SetReadDeadline(time.Now().Add(DrainGrace)); n, err := r.Read(buf); … }
// so the grace bounds the IDLE GAP between reads and never the size or total
// duration of the tail (A9: gaps 27-32 ms under a 50 ms grace while the tail
// ran 12-13x the grace, and still reached EOF whole).

werr := cmd.Wait()          // direct child reaped (≤ WaitDelay past ctx deadline)
reapGroup(cmd.Process)      // release the group FIRST, as today
drains.Wait()               // bounded by the drains' own per-read deadlines; the parent
                            // sets none. Each drain returns EOF, held, or overflow.
if held := heldPipes(outEnd, errEnd); len(held) != 0 {
    return inv, wrap(tail(stderr), errHeldPipe(held, exitStatus))
}
// EOF (before or after the timer): whole output, unchanged from here on
```

Illustrative: `heldPipes`/`errHeldPipe` are names for the shape, not the contract; the
timings in the MVV are Normative. The block deliberately draws NO `graceOn` and NO
`closeReads` call: C1 `precedence:` introduces no separate `graceOn` flag (the per-read
deadline in the loop above is the whole mechanism), and the existing defers already close
both read ends, so an explicit close would double-close (Infra-Audit close-ownership row).
An earlier draft drew both as live parent-level calls, and a reconstruction that took the
drawing at face value rebuilt the one-shot parent deadline — precisely the mechanism MVV
row 6(b) exists to fail, and the reason that row needs its paced shape: row 6(a)'s burst
tail would have let that reconstruction pass (A9).

Two shipped details the shape depends on, both C1 `precedence:` obligations. First,
`readBounded` (`cmdbind.go:352-364`) returns `[]byte` and NO error at all, so today a
deadline error is indistinguishable from a clean EOF at the join. The held-vs-EOF-vs-overflow
distinction is therefore not a fact the current code exposes — each drain must report its
terminal condition out, as C1 `precedence:` now pins, for `heldPipes()` to have anything to
read. Second, the re-arm lives INSIDE that read loop, not in the parent. A one-shot
`SetReadDeadline` from the parent before `drains.Wait()` would bound the whole remaining
tail against the clock — the false held report on a writerless pipe that C1 `precedence:`
names and MVV row 6(b) catches (measured: 15.6% recovered, terminal
`os.ErrDeadlineExceeded`, on both OSes — A9). Both are part of this change, not
pre-existing capabilities.

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| One spawn site for every command path | `internal/cli/cmdbind/cmdbind.go::spawn` | none known (A1 confirms) | Extend | C1 lands once; no per-binding wait. Its signature is unchanged by this record: it returns the `invocation` by VALUE (not a pointer), which is what makes C1 `refusal:`'s "the field is readable on every error path" hold without a nil check |
| Per-invocation condition carrier | `internal/cli/cmdbind/cmdbind.go::invocation` — ships `stdout`, `stderr`, `exitCode`, `exited`, `signaled` | none | Extend | `exited`/`signaled` already carry C1 `refusal:`'s "exited N / killed by signal"; F5's pollability result is recorded here as one more field, which is how it reaches the `Detail` line S5 asserts on |
| Drain bound constant | `internal/cli/cmdbind/cmdbind.go::WaitDelay` | today applied only to `Cmd`'s own goroutines | Reuse | C1 `bound:` — no second bound |
| Read-deadline grace | none — no grace/slack constant exists in this repo or the corpus | n/a | New | C1 `bound:` names `DrainGrace` (50 ms); a mechanism constant, not the knob D-naming rejected |
| Group release before join | `internal/cli/cmdbind/cmdbind.go::reapGroup` / `procgroup_unix.go::killGroup` | cannot reach a `setsid(2)` escapee | Reuse | C1 `residue:` names the escapee as the admitted leak |
| Bounded, overflow-detecting read | `internal/cli/cmdbind/cmdbind.go::readBounded` | blocks in `Read` until EOF or an ended fd, and SWALLOWS the error (`:352-364` returns `[]byte` and no error at all), so a deadline is today indistinguishable from EOF | Extend | the join ends the read from outside; `readBounded` must additionally report its terminal condition so `heldPipes()` can tell held from whole |
| Refusal channel with stderr tail | `*accessor.ExecError` via `cmdbind.go::wrap` (0025:C4 `detail:`) | none | Reuse | the held-pipe reason rides in `Err`, tail in `Detail` |
| Deadline-first classification | `internal/accessor/executor.go::invokeRead` (+ Gate / Apply) | A4 checks the write and read-back paths | Reuse | C1 `class:` — no binding-side `timeout` |
| Stdin bound | `cmd.Stdin` = `bytes.Reader` + `cmd.WaitDelay` (Cmd-owned pipe) | A7 confirms `ErrWaitDelay` reaches the refusing arm | Reuse | C1 `stdin:` names it; no change |
| Escape fixture | 0025's `evidence/spikes/a1-deadline/` S2 (grandchild in group) | never escaped the group; macOS has no `setsid` binary | Extend | FX-deadline-escape adds a `Setsid: true` Go helper |
| Read-end close on every path | `spawn` already defers both closes (`cmdbind.go::spawn`, the `defer func() { _ = outR.Close() }()` / `errR` pair set up at pipe creation) | closes at FUNCTION RETURN, which is after the join either way | Reuse | C1 `precedence:` "after the join every read end is closed" is ALREADY satisfied by the defers, and the Illustrative Code draws no explicit close. The defers stay — an explicit close before them would double-close, and `os.File.Close` returns `ErrClosed` on the second call rather than panicking, so drawing one would cost a swallowed error, not a crash. Nothing in this change moves the fds' release |

### Decision Rationale

Questions-Options-Criteria — the deciding rows are **correctness fit** and **prior-art
alignment**; **blast radius** breaks the tie against handing the pipes back to `Cmd`.

| Criterion | A. Bounded join, refuse (chosen) | B. Bounded join, accept partial | C. Widen kill reach | D. Pipes back to `Cmd` |
| --- | --- | --- | --- | --- |
| Correctness fit (refusal always delivered; output never silently unproven) | both: refusal within `timeout + 2·WaitDelay`; nothing unproven is parsed | refusal delivered; an envelope from a stream with a live unreachable writer is parsed anyway | whole read only if the escapee is found; a double-fork re-escapes | refusal delivered by `ErrWaitDelay`; loses release-group-first, so a correct backgrounding tool is refused |
| Prior-art alignment | matches Go `awaitGoroutines` / `ErrWaitDelay` and both peer CLIs that met this failure | no precedent found for parsing past a bound | no portable precedent; Linux-only primitives | is the Go precedent verbatim, but for `Cmd`-owned pipes only |
| Reversibility | join site plus a rewritten `readBounded` and a new accessor error type; revert is the old `drains.Wait()` and dropping both | same | new syscall surface per platform | re-opens the ADV-2 / FAIL-2 ordering 76c121b fixed |
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

**Variant considered and also rejected — restrict ALT1 to an EXITED child.** A3 does
observe byte-for-byte completeness (1 KiB and 63 KiB, drain reading strictly after the child
was gone), so the "never observed" con above is not unqualified: it holds for a child blocked
in `write(2)`, not for one that exited. The observable distinction is available at the join —
`spawn` knows `werr` and `inv.exited` before the drain verdict (`cmdbind.go:304-321`) — so a
variant could parse when the direct child EXITED cleanly and refuse only when it did not.
Rejected anyway, and on its own grounds rather than the parent alternative's: the direct
child exiting says nothing about the GRANDCHILD still holding the write end, which is the
whole failure shape — an escapee can write into the pipe after the direct child is reaped,
so "child exited" does not make the tail complete, it only makes it complete SO FAR. Parsing
on that basis is the same silent class one condition later. This variant is named because
without it the rejection reads as ignoring A3, and it is the variant a reader will propose
first.

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
  rather than accepted — visibly, with the pipe named. That population is UNSURVEYED: no
  reader in the repo's own fixtures exhibits the escaping-but-correct shape, and whether any
  reader in the wild does is unknown. A5's margin measurement does not speak to it — that
  covers the in-group helper that merely misses the bound, not a permanent escapee. For an
  author whose helper is a third-party binary that daemonizes, the mitigation "the refusal
  names what to fix" does not apply, because they cannot change it; their path is to wrap or
  redirect the helper's stdio, which is what Phase 5 documents. Because the population is
  unknown rather than known-empty, Phase 5 SURVEYS it before it documents: run the existing
  reader/gate/write fixtures plus the repo's own declared bindings under the new refusal and
  record whether any previously-passing binding now refuses. A non-empty result is a
  Stage-6 route-back, not a doc note. No opt-out flag ships and none is proposed — a
  per-binding "accept partial output" knob is ALT1 by another name and was rejected for
  ALT1's reason (a partial answer parsed as whole is the silent-wrong-answer class); the
  escape hatch for an author who cannot change a third-party binary is redirection at the
  binding's command, which costs no contract.
- Negative: the leaked process residue is now explicitly admitted with its output; nothing
  in the CLI reaps a `setsid(2)` escapee.
- Negative: when the declared timeout also elapsed — the common case for an escaping helper —
  the refusal is `timeout` with no reason attached, so the author is told the command was slow
  when it was structurally broken (F6). The bound is delivered; the diagnostic is not.
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
- **F5** Silent risk: a pipe whose poller registration failed cannot be deadlined. The
  pollability check is NEW and this record adds it — nothing in `internal/` calls
  `SetReadDeadline`/`SetDeadline` or names `os.ErrNoDeadline` today. It runs in `spawn` at
  pipe creation: a probe `SetReadDeadline` on each read end, immediately cleared, whose
  error is the signal. `cmdbind` has no logger and gains none (its imports carry no
  `log`/`slog`), so "logged" is not a log line: the condition is carried on the invocation
  and surfaces as a `Detail` line on any refusal that invocation produces, which is what a
  test asserts on. On that host the join stays unbounded — the old hang, confined to a
  non-pollable host and recorded at the check rather than discovered at the join.
  Scope, from the A2 spike: this trigger is **unexercised for a real `os.Pipe`** — every
  read end accepted a deadline (nil) on darwin and linux, and no naturally non-pollable
  pipe was constructible. The error value itself is confirmed on the two shapes that do
  produce it (a regular file; `os.NewFile` over a dup'd pipe fd, adopted without the
  original's poller registration), so S5 simulates the condition at the check rather than
  reproducing it. The residual risk is a host whose pipes are non-pollable — not observed
  on either supported OS.
- **F6** Accepted trade-off, not a gap: a held pipe that also outran the deadline reports
  `timeout` with no `Err`/`Detail` (0025:C4's rule for the class), so the held-pipe reason
  is visible only on `execution_failure`. This is the record's own motivating scenario — an
  author declaring `timeout = "2s"` against a helper that escapes — so the diagnostic
  apparatus this record builds, and the stderr tail `D-the-drain-grace` cites to justify
  `DrainGrace`, are unreachable in exactly that case. What the author gets there is a bounded
  `timeout` instead of a hang, which is the outcome this record owes; what they do not get is
  the reason, and raising the declared timeout does not help, because the helper is broken
  rather than slow. It is accepted rather than fixed because the cure is to amend 0025:C4's
  `detail:` rule ("`Detail` is empty on every refusal carrying no error — timeout, read-back,
  gate-off"), which is Implemented, is this record's Overrides target, and governs two other
  classes; and because the sub-reason constraint is homed in JDR 0003 §D1, not here (JC1). A
  record that declared it in both places would have to be demoted at 7.1 to undo it. The
  coupling rides to 7.1 with the rest of the 0025 relationship.
  **The consequence, stated because acceptance without it is a diagnostic that never
  ships:** the held-pipe reason is unreachable on the record's own PRIMARY journey, so the
  `DrainGrace` mechanism — a new constant, and this record's explicit divergence from Go and
  both peer CLIs — is justified in C1 and in `D-the-drain-grace-is-a-mechanism-constant-not-the-rejected-knob`
  by a benefit that journey does not deliver. The mechanism is still owed: it is what makes
  the FINAL read report the pipe's state rather than the clock, which is the refuse/accept
  predicate itself, and it delivers the reason on the `execution_failure` journeys (a helper
  that answers fast and escapes; S1's write and gate rows). But the justification is narrowed
  here to that, not to the timeout journey. **Delivery is scheduled, not left to the
  relationship:** the sub-reason carrier is homed in JDR 0003 §D1 and 0026 does not amend
  0025:C4; the amendment is owed by whichever record next opens 0025:C4's `detail:` rule, and
  0026 registers the demand at 7.1 via JC1 rather than assuming a peer will notice. If 7.1
  closes with no record owning it, this is a Stage-6 route-back for 0025, not a residue to
  re-accept silently — that re-acceptance is how the diagnostic goes missing for a third
  record in a row.

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
   ≤ 3 s, asserted with the C1 `bound:` tolerance), the executor classifies `timeout`, and
   the refusal's `Detail` is empty as 0025:C4 states for `timeout`. Asserted at the
   `flow resolve` boundary, not only at the binding: the defect was reported there ("still
   blocked after 15 s"), so the row also asserts that `flow resolve` itself terminates and
   emits the class in its envelope. A binding-level pass with a hung CLI would not close
   the original report.
4. The grandchild is still alive (asserted — the admitted residue) and the test kills it.
5. Success shape, same fixture: the child writes a whole envelope and exits 0 while the
   grandchild holds the pipes — expected `execution_failure` within `2·500 ms` of the exit,
   `Err` naming "stdout, stderr" and "exited 0", and the envelope NOT parsed (no values
   returned).
6. Late-but-unheld control (Normative), TWO shapes — one guards each rejected mechanism,
   and shape (a) alone is not sufficient:
   (a) **Burst**: a child that writes 1 MiB to stdout in one `Write` and exits with no
   grandchild, run with the drain goroutine artificially delayed past the timer, returns
   the whole value — the final read under `DrainGrace` drains buffered bytes and sees EOF;
   no refusal. This shape guards C1 `precedence:` against a deadline of literally `now`,
   which short-circuits before the syscall and returns n=0 (A2, measured on both OSes).
   (b) **Paced — the mechanism discriminator**: the same 1 MiB written in 32 chunks
   separated by idle gaps of 0.4·`DrainGrace` (20 ms), total tail duration ~12·`DrainGrace`,
   then exit with no grandchild. Every chunk is recovered and the drain ends on EOF, because
   the deadline is re-armed before each read and so bounds only the idle gap. An
   implementation that instead sets ONE absolute `now + DrainGrace` for the whole final
   drain truncates here — measured at 15.6% of the payload (163840 of 1048576 bytes, 5
   reads, ~51 ms), terminal `os.ErrDeadlineExceeded` — and FAILS this row.
   Shape (a) does NOT distinguish re-armed from single-absolute and must never be the only
   witness: a 1 MiB tail already resident in the pipe drains in 1.6 ms (darwin) / 9.5 ms
   (linux) across 32 non-waiting reads, so both mechanisms pass it (A9,
   `evidence/spikes/a9-a10-regrace/`). Duration, not size, is the discriminating variable —
   even a 1 KiB paced payload loses 62% under the single absolute deadline.

Measured basis (A6, `evidence/spikes/a3-a6-drain-bound/`): worst end-to-end 2.013 s
(darwin) / 2.082 s (linux) against the 2.5 s budget — 487 ms and 418 ms of headroom. The
drain-join leg alone measured 500.1-508.2 ms in 40/40 samples, so the assertion above is
on the TOTAL, which is what C1 commits, and never on a per-leg figure the evidence
contradicts.

### Phase 1: Bounded join

Declare the held-pipe error type in `internal/accessor` beside `ExecError` (C1 `refusal:` —
the import direction forbids its declaration in `cmdbind`). Bound the drain join in `spawn`
by `WaitDelay`; on expiry the parent re-arms a per-read deadline of `now + DrainGrace` on
its own read ends (mechanism from the A2 spike), joins the drains, and refuses through
`wrap` with the held pipe(s) named. `readBounded` gains a reported terminal condition —
today it returns `[]byte` and NO error (`cmdbind.go:352-364`), so the held/EOF/overflow
distinction the join now turns on does not exist in its signature. No other line of `spawn`
moves.

### Phase 2: Contract re-citation

Re-point `cmdbind.go`'s C4 comments ("necessary and sufficient", "run to a true EOF",
`WaitDelay`'s doc) at 0026:C1; 0025 itself is not edited.

### Phase 3: Fixtures and oracles

Add FX-deadline-escape (both shapes) beside the existing grandchild-leak and cancel oracles;
re-run the in-group grandchild oracle under the bound (A5) so the prior point-fixes stay
proven.

### Phase 4: Executor and CLI applied sense

Land the A8 gain, which is outside `cmdbind` and so outside every phase above: match the
held-pipe error in `executor.go::Write`'s `err != nil` arm (`:368-371`, which today returns
`refusalOf(def, timeout, ClassExecutionFailure, err)` and sets no applied sense, unlike the
six sites at `:363`, `:375`, `:400`, `:409`, `:423`, `:472`) and set the applied sense there;
then key `flow_exec.go::accessorFailureOf`'s `ClassExecutionFailure` arm (`:410-412`) on
`Applied()`. That arm is the odd one out today — `ClassTimeout` (`:399-405`) already carries
a phase check plus `detailMayHaveApplied`, and `ClassReadBackIncomplete` (`:419-424`) already
sets it. `Applied()` (`model.go:324`) has zero non-test callers repo-wide, so this phase
creates its first. S1's direct-write leg is the row that proves this phase landed.

**Not descopable.** This phase is late and in a different package from every phase above,
which makes it the one most likely to be deferred when the bound lands and the pressure
comes off — and it is the ONLY user-visible write-safety gain in the record. Deferring it
ships the hang fix together with a new failure: a command write that applied its effect and
exited 0 refuses `execution_failure` with no applied sense, and an agent's documented branch
on that class is retry, so the mutation is applied twice. Phases 1-3 without Phase 4 are
therefore not a shippable subset; if Phase 4 must slip, the bound ships with the write path
excluded (reads and gates only) rather than with an unsensed write. That is the partial-ship
disposition, stated so it is a decision rather than an omission.

**Negative control (required, not optional).** The `errors.As` match lands on the SHARED
`err != nil` arm at `:368-371`, which every other pre-mutation failure shape also reaches —
`resolveArgv0` failure, `os.Pipe` creation failure, a plain exec failure from a mistyped
binding. Only the held-pipe error may set the applied sense there. A match broadened during
implementation (matching `*accessor.ExecError` generally rather than the held-pipe sibling)
flips a write that structurally could not have run to "may have been applied", and the agent
skips a safe retry — the exact inverse of the defect this phase fixes, and invisible to
S1, which only exercises the positive case. S1 gains a negative row: a binding whose command
cannot start refuses `execution_failure` with `Applied()` FALSE and no `detailMayHaveApplied`
in the envelope.

### Phase 5: Surface

SURVEY, then state. The survey runs FIRST because the population it measures is the one
`§consequences` records as UNSURVEYED: run the existing reader/gate/write fixtures and the
repo's own declared bindings under the new refusal, and record whether any binding that
passed before now refuses. Acceptance for this half is the recorded result, empty or not —
"no previously-passing binding refuses" is the expected finding and must be written down
rather than assumed. A NON-empty result is a Stage-6 route-back, not a doc note: it means the
escaping-but-correct shape exists in practice and the trade-off `§consequences` accepts was
priced against a population that turned out to be non-empty.

Then state the restated residue in `docs/cli-output-contract.md`, beside the command-entry
`timeout` text authors already read — the "a process may leak and its output is never read"
consequence, and the remediation the refusal names (close or redirect the helper's inherited
stdio). Acceptance: that doc names the held-pipe refusal and the residue; a reviewer checks
it as they would a scenario. This phase is the record's only user-facing artifact
(Profile: `user-facing yes`) and `§consequences`' claims that the refusal is "visible, with
the pipe named" and the residue "explicitly admitted" hold only once it ships.

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
   `Applied()` is true AND the CLI renders the "may have been applied" text — observed
   on both legs, but GAINED on only one: the direct leg is the site this record changes
   and the read-back leg already behaves this way today (C1 `refusal:` "one site changes,
   not two"). A reading that takes this row as licence to re-key both arms re-keys an arm
   that is already correct. The rendered
   half asserts on the envelope's `detail` field carrying `detailMayHaveApplied`
   (`flow_exec.go:423`'s existing constant, per docs/cli-output-contract.md's error
   envelope) — the same observable the read-back leg already renders, so the two legs are
   asserted the same way and neither needs a new golden. The two
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
   **Expected**: byte-for-byte unchanged from today (C1 `whole:`). "Today" is the existing
   `cmdbind` suite passing unmodified — no new golden is minted for this row; the baseline
   is the current expectations, and a diff in them is the failure. This row is load-bearing
   BECAUSE C1 `precedence:` rewrites `readBounded`'s body into an explicit read loop: the cap
   boundary (exactly `StdoutCap` is a value, `StdoutCap+1` is the overflow failure) and the
   discard-remainder exit are re-implemented, not merely re-called. Assert both boundary cases
   explicitly here rather than relying on the suite's incidental coverage; relaxing an existing
   expectation to make the rewritten loop pass is the failure this row names.

4. **Scenario**: Non-reading child with a stdin payload past the pipe buffer (A7).
   **Expected**: `exec.ErrWaitDelay` from `Wait` reaches `spawn`'s non-`ExitError` arm
   and refuses `execution_failure` after `reapGroup` and the join.
5. **Scenario**: Non-pollable pipe — the creation-time pollability check returns
   `os.ErrNoDeadline` (F5), simulated at the check by injecting a read end that refuses a
   deadline (`os.NewFile` over a dup'd pipe fd, the shape F5 confirms produces the error).
   **Expected**: the condition is recorded at the check, not discovered at the join, and is
   observable as the F5 `Detail` line on that invocation's refusal — that line, naming the
   host condition, is the assertion. The fallback is the documented unbounded path.
6. **Scenario**: Held stdout with the deadline also elapsed (F6).
   **Expected**: classified `timeout` with empty `Detail` per 0025:C4 — asserted so the
   class flip is a tested consequence rather than a surprise.

7. **Scenario**: payload past the ~64 KiB pipe buffer (1 MiB), drain prompt vs
   artificially stalled — the stall is ONE seam, shared by this row, S9 and MVV row 6 so
   they cannot build two: a test-only delay hook on the drain goroutine's start. The
   package has NO existing test-injection surface to hang it on — `cmdbind`'s only
   pre-existing test var is `goos` (platform gating), and today's timing tests drive real
   subprocesses instead — so this seam is NEW and this record adds it: one unexported
   package-level hook set by the test and nil in production (not a build tag, not a sleep
   in production code). S7 is a standing regression trap, so the seam outlives the fixture — the antecedent C1 `whole:` states (A3).
   **Expected**: prompt drain returns the whole 1 MiB; a stalled drain leaves the child
   blocked in `write(2)`, killed by the group signal, with exactly one pipe buffer
   (65536 bytes) recovered as a correct prefix. The test exists so a later refactor that
   serializes or delays the drains fails here rather than silently voiding A3.
8. **Scenario**: a child that RESISTS termination (ignores SIGTERM, survives to its
   `WaitDelay`) while a grandchild holds the pipes — the case A6 did not exercise, and
   the one where leg (a) spends its allowance instead of donating it to leg (b).
   **Expected**: the invocation still returns within the C1 `bound:` total plus its stated
   scheduling tolerance, so the total bound is asserted where the two legs genuinely compete
   rather than only where one underspends (0025:C6).
9. **Scenario**: a LATE drain with no holder — the drain goroutine stalled past the
   timer with bytes buffered and every writer gone (premortem P-2).
   **Expected**: accepted, whole output, no refusal — the final read under `DrainGrace`
   delivers the buffered bytes and then sees EOF. Companion: an escaped writer emitting
   a byte every 10 ms forever is REFUSED (never blocked, never EOF). Together these pin
   "refuse on no-EOF-after-a-final-bounded-read, not on elapsed time"; with a deadline of
   `now` the first of them fails, which is what makes `DrainGrace` load-bearing.

10. **Scenario**: stdout reaches EOF cleanly while STDERR alone is held by the escapee —
   the single-pipe form of C1 `refusal:`, and the case `D-the-drain-grace` cites as the
   grace's purpose ("what the grace preserves is the stderr tail").
   **Expected**: refused; `Detail` names `stderr` alone — the single-pipe form is the bare
   pipe name, and the both-pipes form is the two names comma-separated in the fixed order
   stdout, stderr (MVV row 5's `"stdout, stderr"`) — and the stderr tail collected up to the
   bound survives into `Detail` ahead of the held-pipe reason. Without this row the grace's
   load-bearing justification has no test.

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

**Stage 6 (reconcile) closed this finding — see the disposition
note at the end of this section.** The gate's original text is
kept below as the record of what was open at Stage 5.

No `Docs Only` records. Three remain `Pending`, each with a
plan, and **all three fail the status-consistency clause as
written** — this is the gate's finding, not a clean pass:

- **A9** (per-read re-arm bounds the idle gap) — C1
  `precedence:` states the re-arm as settled normative text and
  MVV row 6 is marked Normative on top of it, while no spike has
  measured a large tail against a re-armed grace. Plan: extend
  `evidence/spikes/a3-a6-drain-bound` per A9, run under CPU
  contention, before Phase 1 closes. Until it runs, C1
  `precedence:`'s re-arm sentence and MVV row 6 are contingent.
- **A10** (pollability probe is safe and detecting) — the probe
  runs on EVERY invocation's hot path and F5 concedes its trigger
  is unexercised for a real `os.Pipe`. C1 `whole:`'s
  byte-for-byte claim and S3 depend on the probe being a no-op.
  Plan: A10's spike, plus S5, before Phase 1 closes.
- **A8** (executor derives `Applied()` from the typed error) —
  C1 `refusal:`'s applied-sense sentence depends on it, and the
  gain sits in Phase 4. Method is Source Search and the sites are
  already named; the residual is the CLI rendering gain, not a
  found fact. Plan: verify at Phase 4 entry.

Disposition: the three are survivable as `Pending` **into
implementation** because each names a specific pre-phase
verification and each "If wrong" is a stated, bounded regression
rather than a redesign. They are NOT survivable into lock as
settled prose, so the clauses that lean on them are marked
contingent above and in C1. Stage 6 closes them; a lock taken
with A9 or A10 still open must demote the dependent sentences
first.

**Stage 6 disposition (2026-09-03) — the contingency is
discharged.** All five then-`Pending` records are now terminal
and the status-consistency clause passes:

- **A9 — VERIFIED by spike** (`evidence/spikes/a9-a10-regrace/`,
  both OSes). C1 `precedence:`'s re-arm sentence and MVV row 6
  are no longer contingent. The spike also CORRECTED the
  assumption's witness: a 1 MiB *buffered* tail does not
  discriminate the re-armed mechanism from a single absolute
  deadline (both recover it whole), so MVV row 6 was amended to
  carry a paced shape 6(b) that does — 0/12 vs 12/12, truncating
  at 15.6% on both OSes. C1 now reads "size or total duration".
  Duration, not size, is the operative variable.
- **A10 — VERIFIED by spike** (same dir). The probe is detecting
  (`os.ErrNoDeadline` on a genuine dup'd pipe fd, so it detects
  poller registration, not file kind) and a no-op on the hot
  path (parked-read 12/12 at 300 ms with a 10 ms probe window;
  round-trip identical in every cell). C1 `whole:`'s
  byte-for-byte claim and S3 are no longer contingent on it.
- **A8 — VERIFIED by source search.** Every cited site holds
  with no drift; the two decisive facts (the
  `ClassExecutionFailure` arm carries no phase check and no
  `Detail`; `Applied()` has zero non-test callers) confirm the
  gain is real and one-sited.
- **A12 — VERIFIED by source search.** The seam is new, the
  hook shape is viable at `cmdbind.go:273`/`:277`, and no third
  build tag is needed.
- **A11 — DOWNGRADED, survivable into implementation.**
  Unverifiable before implementation by construction: its check
  is S3 + S7 under `-race`, and S7 needs the stall seam A12
  confirms this record ADDS. Not MVV-critical — the MVV
  provably cannot exercise the edge. Bound by a named plan:
  Phase 1 closes only on S3 + S7 passing `-race -count=25`, and
  a race report there is a Phase-1 blocker.

No assumption now carries settled-fact prose it does not
support. A11 is the one `Pending` record, and no clause states
its claim as settled: C1 `precedence:` (b) states the
requirement and names `-race` as the check rather than
asserting the edge holds.

No `Verified` stamp is self-referential. A5's and A6's cited
`path::Symbol` anchors resolve on `main` (`cmdbind.go::WaitDelay`
`:66`, applied `:224`; `drains.Wait()` `:302`); A6 carries its
own narrowing of the per-leg claim and is stamped "Verified, with
the per-leg claim narrowed", which is consistent.

### Scope Verification

[Gate key: scope]

[Confirm the Minimum Viable Validation is in scope
and will be executed during implementation, not
deferred. State the specific test or proof.]

### Cross-Cutting Concerns

[Gate key: cross-cutting]

- **Concurrency model** — this RDR's own seam, and the only concern it authors
  policy for. The bounded join replaces `drains.Wait()` as the sole
  happens-before edge between the drain goroutines' slice writes and the parent's
  read of `inv.stdout`/`inv.stderr`; C1 `precedence:` (b) requires the join stay a
  real join (a drain ended by the grace deadline still runs `drains.Done()`) and
  forbids by name the design that lets the parent proceed past a live drain. The
  edge is not asserted as settled: A11 is the open record, and its named check —
  S3 + S7 under `-race -count=25` — is a Phase-1 exit condition, not a follow-up.
  The drains stay concurrent (`D-the-drains-stay-concurrent`): serializing them
  would deadlock a child that fills one pipe while the parent reads the other.
- **Build tool compatibility** — Unix-only process-group syscalls stay behind the
  build tags established by intrastate#1drn (`procgroup_unix.go`); Windows
  continues to take 0025:C4's platform refusal. A12 confirms the drain-start stall
  seam this RDR adds needs **no third build tag** — it is portable Go in
  `cmdbind`. Every spike backing C1 (A2, A9, A10) ran on both darwin and linux, so
  the bound is not a single-platform measurement.
- **Memory management** — the existing 1 MiB stdout cap and 4 KiB stderr tail are
  unchanged and remain the only bound on retained output. C1 `precedence:` rewrites
  `readBounded`'s body into an explicit read loop, so those boundaries are
  re-implemented rather than re-called; S3 asserts both cap boundary cases
  (exactly `StdoutCap`, and `StdoutCap+1` as overflow) explicitly, and relaxing an
  existing expectation to make the rewritten loop pass is that row's named failure.
  Output read on a bound-expired drain is never parsed on any path, so a partial
  buffer is never retained as a value.
- **Incremental adoption** — the narrowing of 0025:C4 ("true EOF" → "whole within
  the bound") can newly refuse a binding that answers correctly but leaves a
  `setsid` helper holding stdout. That population is **unsurveyed**, not
  known-empty, so Phase 5 SURVEYS before it documents: the existing
  reader/gate/write fixtures plus the repo's own declared bindings run under the
  new refusal, and a previously-passing binding that now refuses is a Stage-6
  route-back, not a doc note. No opt-out flag ships and none is proposed — a
  per-binding "accept partial output" knob is ALT1 by another name, rejected for
  ALT1's reason; the escape hatch is redirection at the binding's command, which
  costs no contract.

Two more are owned elsewhere and cited, not restated:

- **Error/refusal surface** — the `execution_failure` sub-reason carrier (`Detail`
  for the reason, `Err` executor-facing, `Applied()` for the applied sense) is
  **JDR 0003 §D1**. C1 `refusal:` is this record's instance of that policy;
  JC1 records the fire, symmetric with cli/0028:JC1's arm-2 fire on the same
  literal. This RDR does not re-decide the carrier.
- **Deadline triple and platform refusal** — 0025:C4, which this RDR declares in
  Overrides and amends only at the point C4 left undecided (bound-vs-whole
  precedence). The residue class stays 0025:F4's, restated honestly rather than
  redefined.

Not applicable, and deliberately not N/A-bulleted above: versioning, licensing,
deployment model, IDE compatibility, secret/credential lifecycle, character
encoding. This RDR claims no byte-identical output, content-addressed identity or
replay-stable hash, so the canonical-form checklist does not apply — S3's
"byte-for-byte unchanged" is a *regression* assertion against today's existing
expectations on the ordinary paths, not a determinism claim about a newly minted
artifact.

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
- `evidence/spikes/a9-a10-regrace/` (A9 the per-read re-arm vs a single absolute deadline, paced and burst shapes; A10 the pollability probe's detection and hot-path safety; darwin and linux)
- JDR 0003 §D1 (`docs/jdr/0003-accessor-binding-seam.md`) — the `execution_failure` sub-reason carrier; C1's `refusal:` line is this record's instance
