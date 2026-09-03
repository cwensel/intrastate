Model: claude-opus-5

# Hostile critique — RDR 0026, bounded output drain on command timeout

## 6. Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `0026:C1` (`bound:`) | "a command write is two invocations (write, then read-back)" is false. `internal/accessor/executor.go::(*Executor).Write` runs THREE bounded invocations: the pre-write protected-key baseline read (`:320` `e.invokeRead(ctx, reader, art, rt, protected)`), the write (`:355`), and the read-back (`:394`). The committed per-journey bound is understated by one whole `timeout + 2·WaitDelay`. | A `flow resolve` write against an escaping helper returns at ~3× the documented bound. An agent caller that sized its own watchdog from the RDR's arithmetic kills the CLI mid-write and cannot tell whether the mutation applied. | §1, §3, premortem, AT-1 |
| C-2 | `0026:C1` (`precedence:`) | The re-arm mechanism is specified as "the deadline is re-armed at `now + DrainGrace` BEFORE EACH read of the final drain," but `cmdbind.go::readBounded` (`:352-364`) has no read loop to re-arm inside — it is `io.ReadAll(io.LimitReader(r, limit))` followed by `io.Copy(io.Discard, r)`. The contract's load-bearing mechanism has no site. Implementing it means rewriting `readBounded` into a hand-rolled loop, which silently re-opens the cap-boundary semantics C1 `whole:` and S3 promise are byte-for-byte unchanged. | The 1 MiB-exactly and cap+1 boundary behaviour shifts under a rewritten read loop; `S3`'s "existing suite passes unmodified" fails, and the implementer either weakens S3 or ships an untested boundary. | §1, §2, premortem, AT-2 |
| C-3 | `0026:C1` (`precedence:`), `0026:A2` | The `graceOn` handoff is an unsynchronised cross-goroutine mutation. The parent marks the drains and the drain goroutines read that mark; the parent then joins and reads `stdout`/`stderr`, which the drain goroutines write. Neither the contract nor the Illustrative Code names a happens-before edge for the mark, and the existing `stdout`/`stderr` slice assignment (`cmdbind.go:275-279`) is ordered today only by `drains.Wait()`. Under a bounded join whose timer fired, the ordering the current code relies on is exactly the one being removed. | `go test -race` reports a data race on the drain mark, intermittently, in CI. Worse if the race is benign-looking: a drain misses the mark, never re-arms, and the join blocks forever — the original hang, now behind a timer that looks correct. | §1, premortem, AT-3 |
| C-4 | `0026:F6`, `0026:S6`, `0026:D-the-drain-grace-is-a-mechanism-constant-not-the-rejected-knob` | The record's motivating scenario produces `timeout` with empty `Detail` (0025:C4 `detail:`: "`Detail` is empty on every refusal carrying no error"). So the held-pipe reason, the pipe names, the remediation text, AND the stderr tail `DrainGrace` exists to preserve are all unreachable in the case the RDR was written for. `DrainGrace` is justified by a benefit its own primary journey cannot deliver. F6 calls this "accepted" and defers the cure to 0025:C4 — which this record's Metadata already declares as its `Overrides` target. | The author declaring `timeout = "2s"` against a `setsid` helper gets `flow-accessor-timeout` and nothing else. They raise the timeout to 30 s, wait 30 s, get the same bare timeout. They never learn a pipe is held. The hang is replaced by a 30-second lie. | §2, §3, premortem, AT-4 |
| C-5 | `0026:A9` (Pending), `0026:MVV` row 6, `0026:C1` (`precedence:`) | A9 is Pending and is the ONLY evidence for the re-arm mechanism at scale, yet C1 `precedence:` states the re-arm as settled normative text and MVV row 6 is marked Normative on top of it. A Pending assumption is load-bearing for a locked contract clause. The record's own Assumption Verification gate forbids exactly this ("no assumption marked `Pending` … may have settled-fact prose elsewhere in the RDR depending on it") — and that gate is an unfilled template block. | If A9 falsifies, a 1 MiB answer on a loaded machine is refused as a held pipe on a writerless pipe. Ordinary large-output readers become intermittently `execution_failure`. The user retries and it works. | §1, §3, AT-5 |
| C-6 | `0026:A10` (Pending), `0026:F5`, `0026:S5` | The pollability probe runs on EVERY invocation's hot path and is unverified (A10 Pending), and F5 concedes the trigger is "**unexercised for a real `os.Pipe`**" — no naturally non-pollable pipe was constructible on either supported OS. The record adds a per-invocation probe, a new `Detail` line, a fallback path, and a test scenario for a condition it cannot demonstrate exists. Meanwhile the fallback IS the hang: "on that host the join stays unbounded." | Best case, nothing — dead code on every command. Worst case, the probe perturbs ordinary reads (A10's own "If wrong") and the byte-for-byte claim fails for every command, not just held ones. | §1, §2, AT-6 |
| C-7 | `0026:A6`, `0026:C1` (`bound:`), `0026:S8` | The measured per-leg figure already breaches the stated bound (drain join 500.1–508.2 ms against a 500 ms `WaitDelay`, 40/40 samples). The record's repair is to move the claim to the TOTAL and lean on leg (a) donating its unused allowance — while conceding in the same clause that a child which RESISTS termination spends leg (a)'s allowance. That case is untested (A6: "was not exercised"), and S8 is the test that would exercise it. So the total bound is asserted on measurements from the one shape where the legs do not compete. | A resisting child plus an escaped grandchild returns past `timeout + 2·WaitDelay`. The contract's one hard number is violated on the composite case, in production, first. | §1, §3, premortem, AT-7 |
| C-8 | `0026:§consequences` (the "UNSURVEYED" negative) | The record admits, in its own words, that the population of correct-but-escaping readers is UNSURVEYED — "whether any reader in the wild does is unknown" — and then ships a behaviour change that converts every member of that unknown population from working to refused. The mitigation ("the refusal names what to fix") is conceded not to apply to third-party daemonizing binaries. No survey is scheduled; no phase adds one; no flag or escape hatch is offered. | A model author whose reader shells out to a tool that daemonizes (a `ssh -f`, a `set -m` background job, a Node `detached: true` with inherited stdio) has a working config on Monday and a hard `execution_failure` on Tuesday, with a remediation they cannot perform because they do not own the binary. | §2, §3, premortem, AT-8 |
| C-9 | `0026:A5`, `0026:F2` | A5's 500 ms margin is measured entirely against a SIGKILLed in-group grandchild that exits promptly. It measures scheduler and pipe-teardown latency, not process-exit latency. The premortem's P-7 (uninterruptible sleep on NFS/FUSE, large-heap teardown, ptrace-stopped helper) is folded nowhere: F2 answers it with the same 27.7 ms number, which is a measurement of the case P-7 says is not the case. The record explicitly scales the risk by "runner oversubscription" and never by process teardown cost. | Intermittent `execution_failure` on a reader whose wrapper backgrounds a JVM or a FUSE-touching helper. Works locally, fails on CI, class says "could not be executed." | §1, §3, premortem, AT-9 |
| C-10 | `0026:A8` (Pending), `0026:§phase-4-executor-and-cli-applied-sense`, `0026:S1` | The record's only real user-visible safety improvement for writes — the `Applied()` re-key so a held-pipe write is not read as "the write did not happen" — sits behind a Pending assumption, in the LAST implementation phase, in a different package from every other phase, and is the phase most likely to be cut when the drain work overruns. `flow_exec.go::accessorFailureOf`'s `ClassExecutionFailure` arm (`:410-412`) has no phase check and sets no `Detail` today, and `Applied()` (`model.go:324`) has zero non-test callers. | Phase 4 slips. A command write applies its effect, exits 0, gets `execution_failure` with no applied sense. The agent's branch on that class is retry. The mutation is applied twice. | §1, §2, premortem, AT-10 |
| C-11 | `0026:C1` (`refusal:`), `0026:§technical-design` | C1 says "The returned invocation carries no stdout; stdout read on a held drain is never parsed on ANY path" — but `spawn` returns `(inv, err)` with `inv` populated (`cmdbind.go:305`, `:326-329`), and the not-parsed guarantee rests entirely on A1's textual observation that today's three callers check `err` before `inv`. Nothing structural enforces it; `invocation.stdout` remains a readable field on a returned struct. The contract states a discipline and the design ships an unenforced convention. | A future gate fast path, or 0028's read-back-through-a-command-reader leg, reads `inv.exitCode` or `inv.stdout` before `err` and derives a verdict from a stream with an unreachable live writer. Silent wrong answer — the exact class 0025 exists to prevent. | §1, premortem, AT-11 |
| C-12 | `0026:§finalization-gate`, `0026:G-assumptions`, `0026:G-proportionality` | Every gate response is an unfilled TEMPLATE guidance block (lint: `placeholder:survived` ×5, `gate:inline`). The Assumption Verification gate is the exact check that would have caught C-5 (a Pending assumption with settled-fact prose depending on it) and the Proportionality gate is the one that would have counted the contracts C-2/C-10 spread across two packages. The record ships with three Pending assumptions and its own detector unrun. | Not user-visible directly; it is the reason C-5, C-6 and C-10 are still in the record. | §2, AT-12 |
| C-13 | `0026:§problem-statement`, `0026:§approach` | The problem statement is "the CLI hangs unboundedly," which one timer solves. The record instead locks a precedence contract, a new constant, a new error type in a second package, a per-invocation pollability probe, a per-read deadline re-arm mechanism, a `readBounded` signature change, an executor re-key, a CLI rendering re-key, and a docs surface — across three packages, with three Pending assumptions and an unrun gate. C1 is one of the longest normative blocks in the corpus (9.0K for one clause set). Reversibility is claimed as "one join site; revert is the old `drains.Wait()`," which is false the moment `readBounded`'s signature changes and a new accessor-side error type ships. | The change lands late, in pieces, with Phase 4 and Phase 5 outstanding. The hang is fixed; the diagnostics, the applied sense, and the user-facing documentation of the residue — the parts the user actually touches — are not. | §2, §3, premortem, AT-13 |
| C-14 | `0026:ALT1`, `0026:§decision-rationale` | Alternative 1 (bounded join, accept partial output) is rejected on a correctness argument the record's own A3 spike undermines and then re-labels. A3 proves byte-for-byte completeness at 1 KiB and 63 KiB with the drain reading strictly after the child was gone. The Decision Rationale keeps this fact as "what makes A's residue text honest, not as a licence to parse" — a rhetorical move, not an argument. The real distinction (the child EXITED vs a child blocked in `write(2)`) is observable at the join: `spawn` already knows `werr` and `inv.exited`. A1-with-an-exited-child was never costed. | A tool that answers correctly and exits 0 while a daemon holds stdout is refused, when the record's own evidence shows its answer was complete and recoverable. The user loses a working reader to a purity argument. | §1, §3, premortem, AT-14 |
| C-15 | `0026:JC1`, `0026:F6` | The joint-check homes the sub-reason carrier in JDR 0003 §D1 and F6 defers the `Detail`-on-`timeout` cure to 0025:C4, which this record declares as its `Overrides` target and does not amend. The record therefore owns the precedence, borrows the carrier, and disclaims the one rule that makes the carrier reachable in its own primary case. Three records each hold a third of one user-visible behaviour and none of them can deliver it alone. | The diagnostic never ships. Each record's post-mortem points at the other two. | §2, premortem, AT-15 |

---

## 1. The three most likely ways implementation goes wrong

### 1.1 The re-arm mechanism has no site, so `readBounded` gets rewritten and the "byte-for-byte unchanged" promise dies quietly

**Root cause in the RDR.** C1 `precedence:` specifies the load-bearing mechanism as a per-read re-arm: "The deadline is re-armed at `now + DrainGrace` BEFORE EACH read of the final drain, so the grace bounds the IDLE GAP between reads and never the size of the tail." The Illustrative Code doubles down: "the re-arm lives INSIDE that read loop, not in the parent: `graceOn` marks the drains, and each read sets its own deadline."

**The passage that enabled it.** `0026:C1`, and the Illustrative Code's own commentary claiming this is "part of this change, not [a] pre-existing capabilit[y]." The record notices that `readBounded` returns no error and correctly says it must gain a reported terminal condition. It does not notice the larger fact: **`readBounded` has no read loop.** The shipped function is

```go
b, err := io.ReadAll(io.LimitReader(r, int64(limit)))
if err != nil { return b }
_, _ = io.Copy(io.Discard, r)
return b
```

There is no per-read site. `io.ReadAll` owns the loop, and the RDR cannot reach inside it. Implementing C1 as written means replacing `io.ReadAll(io.LimitReader(...))` with a hand-rolled `for { r.SetReadDeadline(...); n, err := r.Read(buf) ... }` — a full rewrite of the function whose semantics C1 `whole:`, the `fidelity` mini-check, and S3 all pin as unchanged.

And the cap semantics are subtle in exactly the way a rewrite breaks. Today, overflow is detected because the read is bounded at `StdoutCap+1` and the check at `cmdbind.go:326` is `len(inv.stdout) > StdoutCap`. `io.ReadAll` grows an internal buffer; a hand-rolled loop with a fixed `buf` and manual append has different allocation behaviour, different behaviour at exactly-cap, and a different interaction with the `io.Copy(io.Discard, r)` third exit — which now must ALSO carry a deadline, and which today discards its error (the Load-Bearing Decisions section notices this and then leaves it as a parenthetical).

**Symptom the user sees.** Either S3 fails on a cap-boundary expectation and the implementer relaxes the S3 assertion (the record's guard against regression is the first thing spent), or S3 is left alone and the overflow boundary ships untested under a new loop. In the field: a reader whose stdout is exactly 1 MiB is refused as overflow, or a reader one byte past the cap is accepted. Neither has a test that fires, because the test that would fire was reworded to accommodate the rewrite.

### 1.2 The `graceOn` handoff is a race, and its failure mode is the hang the record was written to remove

**Root cause in the RDR.** The design has the parent goroutine, on timer expiry, communicate a mode change ("re-arm your deadline before each read") to two already-running drain goroutines, then join them, then read the `stdout`/`stderr` slices those goroutines wrote.

**The passage that enabled it.** `0026:C1` `precedence:` and the Illustrative Code's `graceOn(outR, errR)` line. `graceOn` is described as "marks the drains." Nothing in C1, in the `authority` mini-check, or in the Illustrative Code says how the mark is published — no atomic, no channel, no mutex, no happens-before edge. The `authority` table's row "Held-vs-EOF-vs-overflow terminal condition | each drain goroutine, reporting out" names the direction of the report and not the synchronisation.

The record's blindness here is structural. Today's `spawn` is race-free by construction: the two drain goroutines write `stdout` and `stderr`, and the ONLY ordering is `drains.Wait()` at `cmdbind.go:302`. That single join is the happens-before edge for both slices. RDR 0026 removes the unconditional join and replaces it with a timer-gated one plus a cross-goroutine mode change — and never re-establishes the edge. A2's spike does not cover this: A2 measures a parent calling `SetReadDeadline` on a file the drain is blocked in, which is safe because `os.File`'s poller is internally synchronised. The mark is a different object with no such protection.

**Symptom the user sees.** Two shapes, and the second is worse than the first. Shape one: `go test -race` in CI reports a write/read race on the drain mark or the terminal-condition field, intermittently, and the implementer "fixes" it with an atomic without asking whether the ordering was ever right. Shape two: a drain goroutine reads a stale mark, never re-arms, blocks forever on the pipe, and `drains.Wait()` — which C1 still calls unconditionally after `graceOn` — never returns. **That is the original hang, now sitting behind a timer that fired and a design that reads as correct.** The MVV fixture will not catch it because the MVV's child is silent, so the drain is blocked in its FIRST read, which the parent's deadline reaches regardless of any mark.

### 1.3 The write journey's bound is wrong by a factor of 3/2, and it is the number the caller is told to trust

**Root cause in the RDR.** C1 `bound:` commits "Per invocation, the COMMITTED return bound is `timeout + 2·WaitDelay`… a command write is two invocations (write, then read-back), each under its own bound."

**The passage that enabled it.** `0026:C1` `bound:`. It is false against source. `internal/accessor/executor.go::(*Executor).Write` runs a **pre-write baseline read** before the write ever fires:

```go
protected := protectedKeys(reader, hasReader, plannedKeys)   // :309
if len(protected) != 0 {
    rt, _ := reader.timeout()
    raw := e.invokeRead(ctx, reader, art, rt, protected)     // :320
```

That is a full bounded invocation through the same reader — the same reader whose command spawns the same escaping helper. Then `binding.Apply` at `:355`. Then `e.invokeRead(ctx, reader, art, readTimeout, compared)` at `:394`. Three invocations, each of which can hit the drain bound.

Worse, the pre-write read's refusal is **swallowed**: `if raw.class != "" { baselineUnread = slices.Clone(protected) }` at `:322-326`. So the first held-pipe refusal in a write journey costs its full `timeout + 2·WaitDelay` and produces no signal at all — it degrades to "baseline unestablished" and the journey continues. The RDR's `authority` and `disposition` mini-checks enumerate the write's legs as "direct" and "read-back" and never mention the baseline leg. A1 claims one join covers every path, which is true of the join and irrelevant to the arithmetic.

**Symptom the user sees.** An agent driving `flow resolve` on a write against an escaping helper sizes its watchdog from the documented bound — `timeout + 2·WaitDelay`, or generously twice that — and gets killed at roughly three times it. The agent then cannot distinguish "CLI hung again" from "CLI was slow," which is precisely the discrimination this RDR exists to restore. The RDR converted an infinite hang into a bounded wait and then published the wrong bound.

---

## 2. The one section that will be rewritten within 6 weeks of shipping

**`0026:F6` — the accepted trade-off that the held-pipe reason is invisible on `timeout`.**

Not the Normative Contracts. Not the Implementation Plan. F6.

F6 says, in the record's own words, that in "the record's own motivating scenario — an author declaring `timeout = "2s"` against a helper that escapes — … the diagnostic apparatus this record builds, and the stderr tail `D-the-drain-grace` cites to justify `DrainGrace`, are unreachable in exactly that case."

Read that again against 0025:C4 `detail:`, which is Implemented and which this record declares as its `Overrides` target: "`Detail` is empty on every refusal carrying no error (timeout, read-back, gate-off)." The executor's deadline-first rule (`executor.go:198-200`) returns `readOutcome{class: ClassTimeout}` with no error at all — the binding's `*accessor.ExecError`, its held-pipe reason, its pipe names, its exit status, its remediation text, and its stderr tail are all dropped on the floor. `flow_exec.go::accessorFailureOf`'s `ClassTimeout` arm emits "the accessor `X` timed out" and nothing more.

So the shipped user experience of the record's primary journey is: **a bare timeout, indistinguishable from a slow command.** The record builds `DrainGrace` — a new constant, a new mechanism, a divergence from Go and from both peer CLIs, justified explicitly as "what the grace preserves is the stderr tail … the text that tells a user WHICH helper held the pipe" — and then concedes that text never renders on the journey the record was written for. `D-the-drain-grace` ends with "The grace is justified on the leg that renders it, not on all of them," which is an honest sentence and a fatal one.

Six weeks after shipping, the first support thread reads: "my command reader times out at 2s, I raised it to 30s, now it times out at 30s, what is wrong?" The answer requires reading the RDR. The fix is one of:

1. Amend 0025:C4 `detail:` to allow `Detail` on `timeout` — which F6 rejects as too broad, correctly, since it governs read-back and gate-off too.
2. Carve out a held-pipe exception in the executor's deadline-first path — a new class, or a `Detail` that survives the class flip.
3. Stop classifying a held pipe as `timeout` at all, which contradicts C1 `class:` and A4.

Every one of these rewrites F6 and drags C1 `class:` with it. F6 was written to make the record lockable, not to make the behaviour right. It is the section that pays.

Second candidate, and the reason F6 wins: `§consequences`' "UNSURVEYED" paragraph. That one gets rewritten too, but as a post-mortem, not as a design.

---

## 3. The one assumption that will not survive first contact with a real user

**A5.** "500 ms (`cmdbind.go::WaitDelay`) is long enough that a NON-escaping backgrounded grandchild SIGKILLed by `reapGroup` releases its pipe ends within the bound on CI, so the ordinary 'wrapper script backgrounds a helper' success case pays nothing and is never refused."

The evidence is 160 trials, worst 27.7 ms, "5.5% of the bound." That number is real and it measures the wrong thing.

What A5 timed is the interval from `kill(-pgid, SIGKILL)` to drain EOF for a **Go helper process** that the spike itself spawned. Go processes have small address spaces, no finalizers to run, no filesystem state to unwind. What A5 did not time — and F2 answers with the same number, which is the tell — is the interval for the processes real wrapper scripts actually background:

- A JVM indexer with a multi-GiB heap. `exit_mmap` on a large address space is not free, and it is not scheduler-bound, so A5's "scale by runner oversubscription" model does not apply to it.
- Anything sitting in uninterruptible sleep on a network filesystem or a FUSE mount at the moment SIGKILL arrives. SIGKILL is delivered immediately; the process does not exit and does not release its fds until it leaves the kernel. There is no bound on that.
- A Node process with `detached: true` and `stdio: 'inherit'`, or a Python child with `start_new_session=True`. These are the *common* backgrounding idioms and the second and third of them **escape the group entirely**, meaning they do not even reach A5's population — they reach C-8's.

A5's failure mode is the worst kind: it converts a working configuration into an intermittent one. The user's reader works on their laptop, works in the first CI run, and fails on the third — with class `execution_failure` and the message "the accessor `X` could not be executed," which is true and useless. `F2` acknowledges the shape ("intermittent `execution_failure` on a reader that works locally") and then bounds the risk with a number drawn from the population that does not exhibit it.

The record has one lever and does not pull it: `WaitDelay` is a constant, and the Risks section says "widening it is a one-line follow-up, not a contract change." That is not true once C1 `bound:` commits `timeout + 2·WaitDelay` as the published number and `D-naming` deliberately couples the drain timer to `Cmd`'s stdin/post-cancel timer ("a future tuning moves both — accepted"). Widening the drain bound to survive a JVM also widens the stdin bound for every command in the system. The record accepted that coupling as a virtue.

---

## 4. Premortem

*Written as if the failure has already occurred.*

It is 2026-11-14. RDR 0026 shipped ten weeks ago. `flow resolve` no longer hangs. Three separate defects trace to this record, and the tracker entry that opened it — intrastate#936e — was closed on the day Phase 3 went green.

**Week 2 — the bound that was not the bound.** An agent orchestrating a `flow resolve` write against a repository-tagging accessor started getting killed by its own watchdog. The watchdog was set at `timeout + 2·WaitDelay + 1s`, copied verbatim from C1 `bound:`. The accessor's write command wrapped a vendor CLI that daemonized a sync agent with inherited stdio. `(*Executor).Write` ran `e.invokeRead(ctx, reader, art, rt, protected)` at `executor.go:320` for the pre-write baseline — a full bounded invocation through the same command reader, hitting the same held pipes, costing the same `timeout + 2·WaitDelay` — and then, at `:322`, discarded the refusal into `baselineUnread` and continued. `binding.Apply` paid it again. `e.invokeRead(ctx, reader, art, readTimeout, compared)` at `:394` paid it a third time. Three legs, one documented as two. The agent killed the CLI in the middle of leg two and could not determine whether the mutation had applied. The RDR's own `disposition` mini-check listed "Held pipe on the DIRECT write" and "Held pipe on the write's READ-BACK" as the two write rows and never enumerated the baseline read; A1's "one join covers every path" was read as an arithmetic claim when it was only a code-path claim.

**Week 4 — the hang came back, wearing a timer.** `spawn`'s bounded join shipped as: `waitBounded(&drains, WaitDelay)`, then `graceOn(outR, errR)` marking each drain's per-read re-arm, then `drains.Wait()`. `graceOn` set a plain `bool` field on a shared drain struct. The rewritten `readBounded` — no longer `io.ReadAll(io.LimitReader(r, limit))`, because C1's per-read re-arm has no site inside `io.ReadAll` — checked that field at the top of its loop. On the MVV fixture the drain is blocked in its first `Read` when the mark is set, so the parent's `SetReadDeadline` unblocks it directly and the mark is irrelevant; the test was green from the first run. In production, on a child that trickled stderr, the drain was between reads when the timer fired, missed the mark on a stale cache line, re-entered `Read` with no deadline, and blocked. `drains.Wait()` never returned. The user saw an unbounded hang on `flow resolve` — the exact defect intrastate#936e reported, ten weeks after it was closed. The `-race` detector had flagged the field twice in CI and both flags were closed as "known flaky, drain teardown."

**Week 6 — the cap boundary.** The `readBounded` rewrite (forced by 1.1) replaced `io.ReadAll(io.LimitReader(r, StdoutCap+1))` with a manual loop and a `bytes.Buffer`. `S3`'s "the existing `cmdbind` suite passes unmodified" failed on the exactly-1-MiB case, which had been a value and became an overflow. The implementer, under Phase 3 pressure, adjusted the test expectation to `StdoutCap-1` and noted it in `deviations.md`. A reader emitting exactly 1 MiB of JSON — a repository manifest, deterministically sized — began refusing `execution_failure` with "the child's stdout exceeded the 1048576 byte bound." The `fidelity` mini-check's row "stdout at the 1 MiB cap | exactly the cap is a value; one byte more is a detected overflow | none" was the only place that invariant was written down, and nothing checked the record against the diff.

**Week 7 — the write applied twice.** Phase 4 — the `Applied()` re-key in `executor.go::Write`'s `err != nil` arm and `flow_exec.go::accessorFailureOf`'s `ClassExecutionFailure` arm — was descoped when Phases 1–3 overran. It was the only phase outside `cmdbind`, it depended on A8 which was still Pending at lock, and `Applied()` still had zero non-test callers. A write command applied its mutation, exited 0, and had its stdout held by a detached helper. `spawn` refused. `executor.go:368-370` returned `refusalOf(def, timeout, ClassExecutionFailure, err)` with no applied sense. `flow_exec.go:410-412` rendered "the accessor `repo-tag` could not be executed" with no `detailMayHaveApplied`. The agent's branch on `execution_failure` is retry. The tag was written twice. The second write's read-back succeeded, so nothing downstream noticed. This is the guessed-absence class `0004:C14` forbids, and A8's own "If wrong" predicted it verbatim: "an agent caller retries blind instead of reading back."

**Week 9 — the silent majority.** Support opened seven threads with the same shape: "my command reader worked last week." Each was a wrapper backgrounding a helper with inherited stdio — `ssh -f`, a `set -m` job, a Node `detached: true` indexer. None of them used `setsid(2)` deliberately; all of them escaped the group. All of them are refused now. `§consequences` had said the population was "UNSURVEYED … whether any reader in the wild does is unknown," and no phase surveyed it. The remediation the refusal names — "close or redirect the helper's inherited stdio" — is unactionable for six of the seven, because the helper is a third-party binary. `§consequences` had already conceded this too.

**Week 10 — the diagnostic that was never reachable.** Every one of those seven users had a declared `timeout` on the reader, because that is what the docs tell them to declare. So every one of them hit F6: the deadline elapsed, the executor's deadline-first rule at `executor.go:198-200` returned `ClassTimeout`, 0025:C4's `detail:` rule zeroed the `Detail`, and `accessorFailureOf`'s `ClassTimeout` arm rendered "the accessor `X` timed out." The pipe names, the exit status, the remediation, and the stderr tail that `DrainGrace` — a new constant, a deliberate divergence from Go and from both surveyed peer CLIs — was built to preserve, were computed and discarded. Six of the seven raised their declared timeout first. All six waited longer and got the same message.

The post-mortem row reads: *the record correctly identified that the precedence had never been stated, stated it, and then shipped every piece of the user-visible half — the applied sense, the held-pipe reason, the residue documentation, the population survey — as a later phase, a Pending assumption, a homed constraint in a sibling record, or an accepted trade-off.*

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

These are review-time gates on the RDR, not implementation tests. Each is stated so it can be run against the record and its cited source before lock.

**AT-1 — enumerate every bounded invocation in the write journey (catches C-1).**
```gherkin
Given the RDR commits a per-journey return bound in 0026:C1 `bound:`
When a reviewer enumerates every call to a bounded binding invocation inside
     internal/accessor/executor.go::(*Executor).Write
Then the count MUST equal the count C1 states
And each enumerated leg MUST appear as a row in the `disposition` mini-check
And a leg whose refusal is discarded (executor.go:322 baselineUnread) MUST be
     named in C1 `bound:` as a leg that costs the bound and yields no signal
```
Run today: three legs found (`:320`, `:355`, `:394`); C1 states two. FAIL.

**AT-2 — every normative mechanism must name an existing or explicitly-new call site (catches C-2).**
```gherkin
Given 0026:C1 `precedence:` requires a deadline re-armed "BEFORE EACH read"
When a reviewer opens the function that performs those reads
     (internal/cli/cmdbind/cmdbind.go::readBounded)
Then that function MUST contain a per-read site the re-arm can occupy
Or the RDR MUST state that readBounded is REWRITTEN, and 0026:S3's
     "byte-for-byte unchanged / existing suite passes unmodified" MUST be
     re-derived against the rewrite rather than asserted against today
And the exactly-at-cap and cap+1 boundaries MUST have named assertions that
     survive the rewrite
```
Run today: `readBounded` is `io.ReadAll(io.LimitReader(...))` with no per-read site; the RDR describes only a signature change ("gains a reported terminal condition"). FAIL.

**AT-3 — every cross-goroutine mode change names its happens-before edge (catches C-3).**
```gherkin
Given 0026:C1 changes drains.Wait() from an unconditional join to a
      timer-gated one, and introduces a parent-to-drain mark (graceOn)
When a reviewer asks how the mark is published and how the stdout/stderr
      slices are re-established as safely readable after a timer-fired join
Then the RDR MUST name the synchronisation primitive (atomic / channel /
      mutex) for the mark and the terminal-condition report
And MUST state what happens if a drain observes the mark late or not at all
And MUST include a scenario in which the drain is BETWEEN reads when the
      timer fires (not blocked in its first read, as the MVV fixture is)
```
Run today: `graceOn` is "marks the drains"; no primitive named; every scenario (MVV, S1, S9) has the drain blocked in a read. FAIL.

**AT-4 — a diagnostic must be reachable on the record's own primary journey (catches C-4, C-15).**
```gherkin
Given 0026 introduces DrainGrace justified by preserving the stderr tail
      that names which helper held the pipe
And given the record's motivating scenario is a declared timeout against an
      escaping helper
When a reviewer traces that scenario through executor.go::invokeRead:198-200
      and flow_exec.go::accessorFailureOf's ClassTimeout arm
Then the justifying diagnostic MUST render
Or the RDR MUST NOT justify a new mechanism on a benefit its primary journey
      cannot deliver, and MUST either amend the blocking rule or drop the
      mechanism
```
Run today: 0025:C4 `detail:` zeroes `Detail` on `timeout`; F6 concedes the diagnostic is unreachable; `DrainGrace` ships anyway. FAIL.

**AT-5 — no Pending assumption may underwrite Normative text (catches C-5).**
```gherkin
Given 0026:A9 is Status Pending
When a reviewer greps the record for prose that depends on A9's claim
Then no Normative Contract clause and no row marked "(Normative)" may state
      A9's conclusion as settled
```
Run today: C1 `precedence:` states the re-arm as settled; MVV row 6 is marked "(Normative)". FAIL. (This is verbatim the record's own `G-assumptions` gate text — which is an unfilled template block, see AT-12.)

**AT-6 — a new hot-path mechanism must have a demonstrated trigger (catches C-6).**
```gherkin
Given 0026:F5 adds a pollability probe on every invocation
And given F5 states the trigger is "unexercised for a real os.Pipe" and no
      naturally non-pollable pipe was constructible on either supported OS
When a reviewer weighs a per-invocation cost against an unobserved condition
Then the RDR MUST either demonstrate the condition on a supported platform
Or defer the probe and state the unbounded fallback as an open residue
And 0026:A10 (probe is a no-op on the ordinary path) MUST be Verified before
      the probe is normative, since its "If wrong" breaks C1 `whole:` for
      EVERY command
```
Run today: A10 Pending; trigger unexercised; probe normative. FAIL.

**AT-7 — a committed bound must be measured on the shape where its legs compete (catches C-7).**
```gherkin
Given 0026:C1 `bound:` commits timeout + 2·WaitDelay
And given 0026:A6 measures leg (b) at 500.1-508.2 ms against a 500 ms
      WaitDelay in 40/40 samples, and the total holds only because leg (a)
      donates its unused allowance
When a reviewer asks for the measurement on a child that RESISTS termination
      (leg (a) spending its own allowance)
Then that measurement MUST exist before the total is committed as Normative
```
Run today: A6 states "that case was not exercised"; S8 is planned, not run; the total is Normative. FAIL.

**AT-8 — a behaviour change that breaks working configurations must survey the population (catches C-8).**
```gherkin
Given §consequences states the affected population is UNSURVEYED and that
      the stated remediation does not apply to third-party binaries
When a reviewer asks which phase surveys it, or what escape hatch ships
Then a survey MUST be scheduled, or an opt-out MUST exist, or the RDR MUST
      state that breaking an unknown population is an accepted cost with a
      named owner
```
Run today: no phase surveys; no opt-out; the concession stands alone. FAIL.

**AT-9 — a margin measurement must cover the population it is used to bound (catches C-9).**
```gherkin
Given 0026:A5 measures SIGKILL-to-fd-release at worst 27.7 ms over 160 trials
And given 0026:F2 uses that number to bound the risk of refusing a correct
      in-group helper
When a reviewer asks which processes were measured
Then the measured population MUST include the ones real wrappers background:
      a large-heap runtime, a process in uninterruptible sleep on a network
      or FUSE mount, a stopped/traced helper
Or F2 MUST state that its bound covers only prompt-exiting helpers
```
Run today: the spike measures Go helper processes; F2 scales the risk by CPU oversubscription only. FAIL.

**AT-10 — the user-visible safety half may not be the last, most-droppable phase (catches C-10).**
```gherkin
Given 0026 Phase 4 is the only phase that stops a held-pipe write from
      rendering as "the write did not occur"
And given it depends on 0026:A8 (Pending), lives outside cmdbind, and creates
      the first non-test caller of Applied()
When a reviewer sequences the plan
Then the phase whose omission causes a duplicate mutation MUST NOT be last
Or the RDR MUST state what ships if Phase 4 slips, and whether shipping
      Phases 1-3 alone is acceptable
```
Run today: Phase 4 is second-to-last, A8 is Pending, no partial-ship disposition is stated. FAIL.

**AT-11 — a "never parsed" guarantee must be structural or explicitly conceded as conventional (catches C-11).**
```gherkin
Given 0026:C1 `refusal:` states "stdout read on a held drain is never parsed
      on ANY path"
When a reviewer asks what enforces it
Then either spawn MUST NOT return a populated invocation alongside the held-
      pipe error (zero the field, or return a distinct type)
Or the RDR MUST state that the guarantee is a caller convention verified by
      inspection (A1) and name the test that fails if a future caller — 0028's
      read-back-through-a-command-reader leg included — breaks it
```
Run today: `spawn` returns `(inv, err)` with `inv` populated; the guarantee rests on A1's textual observation of three current call sites. FAIL.

**AT-12 — the Finalization Gate must be answered before lock (catches C-12).**
```gherkin
Given the RDR carries a Finalization Gate whose Assumption Verification item
      forbids settled-fact prose depending on a Pending assumption
And whose Proportionality item requires a contract count across seams
When a reviewer opens §finalization-gate
Then each gate key MUST carry a written response
```
Run today: all five gate responses are unmodified TEMPLATE guidance blocks (lint: `placeholder:survived` ×5). The Assumption Verification gate is the exact detector for AT-5, and the Proportionality gate is the detector for C-13's three-package spread. FAIL — and this single failure is why AT-5, AT-6 and AT-10 also fail.

**AT-13 — reversibility claims must be checked against the actual change surface (catches C-13).**
```gherkin
Given §decision-rationale claims reversibility as "one join site; revert is
      the old drains.Wait()"
When a reviewer lists what must be reverted
Then the list MUST include: readBounded's signature and loop, the new
      accessor-side held-pipe error type, DrainGrace, the pollability probe,
      the executor Applied() re-key, the flow_exec rendering re-key, and the
      docs surface
And the Reversibility row MUST be re-scored against that list
```
Run today: seven items across three packages behind a one-line revert claim. FAIL.

**AT-14 — a rejected alternative must be rejected against the record's own evidence (catches C-14).**
```gherkin
Given 0026:ALT1 (bounded join, accept partial) is rejected because the
      output's completeness "is asserted from pipe semantics, never observed"
And given 0026:A3 measures byte-for-byte completeness at 1 KiB and 63 KiB
      with the drain reading strictly after the child was gone
And given spawn already knows werr / inv.exited at the join
When a reviewer asks whether ALT1 restricted to an EXITED child (the exact
      antecedent A3 verifies) was costed
Then that variant MUST appear in the matrix as its own column, or the
      rejection MUST state why an observable antecedent is not observable
```
Run today: the variant is absent; A3's fact is re-labelled "what makes A's residue text honest, not … a licence to parse." FAIL.

**AT-15 — a user-visible behaviour split across records must have one owner who can deliver it (catches C-15).**
```gherkin
Given 0026:F6 defers the Detail-on-timeout cure to 0025:C4 (Implemented, this
      record's Overrides target, not amended here)
And given 0026:JC1 homes the sub-reason carrier in JDR 0003 §D1
When a reviewer asks which record can ship the held-pipe diagnostic on the
      primary journey
Then exactly one record MUST be able to, and it MUST be named with a schedule
```
Run today: 0026 owns the precedence, 0025:C4 blocks the `Detail`, JDR 0003 §D1 homes the carrier; none can deliver alone; no schedule. FAIL.
