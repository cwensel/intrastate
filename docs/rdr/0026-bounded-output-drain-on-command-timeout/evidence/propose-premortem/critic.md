Model: claude-fable-5

# Premortem critic: bounded output drain on command timeout (propose stage)

Worked from the brief alone. No repository, record, or peer file was read.

## Findings ledger

| ID | passage/claim | Failure mode | Symptom user sees | Origin |
|----|---------------|--------------|-------------------|--------|
| P-1 | "one bounded join covers every path"; residue "the invocation returns within timeout + 2·WaitDelay" | The fix bounds the two READ drains (stdout, stderr). The third inherited pipe, the child's STDIN (the write binding's payload carrier), has a writer goroutine on the parent side whose join is not named anywhere in the approach. An escaped grandchild that inherits the stdin read end and never reads it leaves the parent's payload write blocked on a full pipe with a live reader (no EPIPE). If stdin is a manual `os.Pipe` like stdout/stderr, nothing bounds it; if it is an `io.Reader` handed to `Cmd`, `Cmd.Wait` returns `exec.ErrWaitDelay` on a zero exit and the spawn's handling of that non-ExitError is untraced. Same class as S4: fixed for the subset the fixture exercised. | `flow resolve` hangs unboundedly on a `write` whose payload exceeds the pipe buffer (64 KiB Linux, 16-64 KiB darwin) when the wrapper backgrounds a detached helper and exits without reading stdin. No refusal. | Prospective hindsight / Refutation targets (S4) |
| P-2 | Claim 3: "a drain still blocked with an empty pipe past the bound has already consumed every byte" | The bound as written fires on ELAPSED time, not on BLOCKED state. A drain goroutine that is merely late (descheduled, GC pause, CPU-starved CI runner) with bytes still in the pipe and NO holder is refused with a Detail that names a held pipe that nothing held. Claim 3's premise ("blocked with an empty pipe") is a property the join must enforce, not a fact about time. | Intermittent `execution_failure` "stdout held open after exit" on ordinary commands under CI load; retried, it passes. Flaky pipeline, wrong diagnosis in Detail. | Obstacle negation (C3) |
| P-3 | R1 rejection: "a write's read-back through such a reader would confirm a write on unproven output" | The WRITE command itself applied its side effect and exited 0; only its stdout is held. The refusal is classed `execution_failure`, which to an agent caller reads "the write did not happen". The caller's branch on that class is retry. Duplicate writes. The refusal must carry the child's exit status so the caller can distinguish "did not run" from "ran, output unproven". | Duplicate commits / duplicate rows / double-posted webhooks after an agent retries a write the CLI refused. | Obstacle negation (R1) |
| P-4 | Claim 4: deadline-first rule classifies `timeout` "whenever the declared timeout elapsed" | The held-pipe join adds up to WaitDelay AFTER the child exited. A child that exits 0 at t with an escaped holder is classed `timeout` when t + WaitDelay >= declared timeout and `execution_failure` otherwise. The class of the SAME defect flips with the declared timeout. If Detail differs by class, the user chases a phantom slow command. Also unverified: whether the executor tests `ctx.Err()` or `errors.Is(err, DeadlineExceeded)` on the binding's error; if the latter and the new typed error does not wrap the context error, timeouts with a held pipe become `execution_failure`. | A 50 ms command with a detached helper is reported as `timeout` at timeout=500ms and `execution_failure` at timeout=5s. The author raises the timeout, the class changes, the cause is still hidden. | Obstacle negation (C4) / Refutation targets (S3) |
| P-5 | "Output read on a bound-expired drain is never parsed, on any path" | Stated as a property; the four bindings each decide (inv, err) precedence separately. If spawn returns a populated invocation struct alongside the typed error and any binding (or the gate's allow/deny fast path) inspects `inv.exitCode == 0` or `inv.stdout` before `err`, the held-pipe output IS parsed. Separately, any early return in spawn on a non-nil `cmd.Wait()` error (ExitError, ErrWaitDelay, ctx error) that sits ABOVE `reapGroup` and the join skips both: the group is not killed and the drains leak. | A gate allows on output from a stream with an unreachable writer; or a non-zero-exit child with a backgrounded helper leaks a goroutine and two fds per invocation until the CLI exits. | Refutation targets (S1) |
| P-6 | Claim 2: the joiner can end a pending `Read` on an `os.Pipe` read end | True only for pollable fds. If `os.Pipe` falls back to a blocking fd (nonblock or poller registration failed), `SetReadDeadline` returns `os.ErrNoDeadline`, `Close` does not unblock the `Read`, and a join that waits for the goroutine after the deadline blocks forever: the old hang returns on exactly the hosts least able to explain it. If the join does NOT wait for the goroutine after the deadline, the goroutine keeps appending to the stdout buffer while the caller inspects it (data race) and the fd is never closed (fd leak per refusal). | Rare hosts hang as before; or `go test -race` reports a write/read race on the invocation's stdout buffer; or fd exhaustion after many refusals. | Obstacle negation (C2) |
| P-7 | Claim 5: 500 ms suffices for a SIGKILLed non-escaping grandchild to release its pipe ends | A group member in uninterruptible sleep (NFS, FUSE, a kernel driver) or a large-heap runtime (JVM, node with a big RSS) can take longer than 500 ms between SIGKILL and fd release. The correct answer is refused. Separately, the "ordinary" escape is far more common than `setsid(2)`: `set -m` in a shell script, Node `detached: true` with `stdio: 'inherit'`, Python `start_new_session=True`, `ssh -f` ControlMaster. These are refused every time. The approach is right to refuse them, but without a remediation hint in Detail (redirect the helper's stdio) the author has no path. | "stdout held open" refusal on a wrapper that backgrounds a helper under `set -m`; author has no idea what to change. | Obstacle negation (C5) |
| P-8 | Claim 6: total bound `timeout + 2·WaitDelay` | Three unstated cardinalities: (a) two drains bounded SEQUENTIALLY give 3·WaitDelay; (b) if `reapGroup` polls `kill(-pgid, 0)` until ESRCH, that poll has no stated bound (a D-state member holds it); (c) `Cmd`'s WaitDelay timer starts only after `Cancel` RETURNS, so a Cancel that does more than send a signal delays it. And the user-visible write journey is two invocations (write, read-back), so its bound is 2×. | Return at timeout + 1.5 s instead of + 1 s; or the write journey takes twice the documented bound. | Refutation targets (S2) / Obstacle negation (C6) |
| P-9 | Claim 7: Go precedent and "two comparable Go CLIs" establish "bound wins, parent closes, caller sees an error"; "no precedent for preserving a whole-read promise" | `Cmd.Output()` under `ErrWaitDelay` returns BOTH the bytes read and the error; Go's precedent leaves the accept/refuse choice to the caller. At least one widely used Go consumer that hit this exact failure (git/ssh/gpg helpers holding the go command's pipes) is known to tolerate `ErrWaitDelay` on a zero exit with output. The two CLIs are unnamed. If the resolve stage cannot quote the handling code, this is borrowed authority for a disposition the precedent may not contain. | None directly; the risk is a rejected R1 being re-argued at review with the same citation inverted. | Refutation targets (S3) / Obstacle negation (C7) |
| P-10 | Claim 8: reusing WaitDelay "makes the predecessor's wording true" | Two independent timers now share one name and one value. `cmd.WaitDelay` never touched manual pipes (that was the defect); the sentence becomes true only of the NEW drain timer. A later change that sets `cmd.WaitDelay = 0` to disable the stdin bound, or tunes it for slow stdin consumers, silently changes the drain bound too. | A future tuning for slow stdin readers lengthens or unbounds the held-pipe refusal. | Obstacle negation (C8) |
| P-11 | R3 rejection cites "3.0 s refusal vs 221 ms" | Measured "before the manual-pipe change", under whatever WaitDelay was then. With WaitDelay = 500 ms the R3 cost is 500 ms, not 3.0 s. The ordering argument still stands; the number does not. | None; a reviewer catches the stale figure and doubts the rest. | Refutation targets (S3) |
| P-12 | Residue: the escaped process "retain[s] the pipes it inherited" | Once the parent closes its read end, the escaped process's next write gets EPIPE / SIGPIPE and, under default disposition, dies. The residue is not "retains the pipes"; it is "keeps broken write ends and is killed on its next write unless it ignores SIGPIPE". Also: "closing / deadlining" must be BOTH (deadline to unblock, then close) or fds leak per refusal (R5's own objection). | A deliberately launched background service dies on its first log line after the CLI returns; the author blames the CLI. | Obstacle negation (R5) |
| P-13 | "carrying the stderr tail" vs "output read on a bound-expired drain is never parsed" | The stderr tail IS read from a possibly bound-expired drain and surfaced. Not a defect, but the two sentences contradict as written; the rule should be "stdout is never parsed; the stderr tail collected up to the bound is carried as Detail". | Reviewer confusion; an implementer drops the tail to satisfy the stronger sentence, and the refusal loses its only diagnostic. | Obstacle negation |
| P-14 | R2 rejection: "no portable Unix primitive exists" | The seam already forks by build tag; a Linux-only subreaper is consistent with that pattern. The rejection's true reason is that it would be a fourth mechanism patch leaving precedence unstated, which stands; the portability reason does not. | None; reviewer notes the inconsistency. | Obstacle negation (R2) |

## 1. Prospective hindsight: the failure narrative

RDR 0026 shipped. The fixture that carried it (a child that `setsid`s a grandchild holding stdout and stderr, then exits 0) returned a refusal in 550 ms on darwin and linux, `TestSpawn_EscapedGrandchildHoldsStdout` was green, and the residue clause was rewritten to promise the refusal is never withheld.

Nineteen days later a pipeline hung again. The model author's `write` binding fed a 210 KiB payload to a wrapper script that started a detached indexer (`node --detached`, stdio inherited) and exited 0 without reading stdin. The parent's stdin writer blocked at 64 KiB on a full pipe whose read end the indexer held; the write never returned EPIPE because a reader existed. The two output drains hit their bound, ended cleanly, and then the spawn function joined its stdin writer, which had no bound. The agent driving `flow resolve` waited forever, exactly the symptom 0026 existed to remove. The post-mortem row read: "owed at propose; the approach named two of the three inherited pipes; the fixture held stdout and stderr; nobody asked what held stdin." Class S4, second occurrence at the same seam.

In parallel, a slower failure had been accumulating. On the shared CI runner, 0.4% of ordinary `read` invocations were refused `execution_failure` with Detail "stdout held open after exit". Nothing held them. The drain goroutine had been descheduled under load with 40 KiB still in the pipe; the join's timer fired on elapsed time and the deadline ended a read that would have reached EOF a few milliseconds later. Retries passed. The team added a retry loop at the executor, which masked the next real defect.

And once: a `write` command applied its change, exited 0, and was refused `execution_failure` because a detached helper held stdout. The agent's branch on `execution_failure` was retry. The change was applied twice.

## 2. Obstacle negation

Separator between items is `---`.

### Claim 1: every invocation passes through one spawn function

Negation: one function does not mean one join. The stdin carrier is a third pipe; the write binding's read-back is a second spawn under either its own or the write's context. If the read-back shares the write's context, it starts with an exhausted deadline after a held-pipe write and is classed `timeout` for a hold unrelated to time. Scenario: write refuses at t = timeout + 1 s; read-back begins with ctx already done; `Cmd.Start` under a done context returns the ctx error before the process runs; the journey reports two `timeout`s for a 50 ms command.

---

### Claim 2: the joiner can end a pending Read from another goroutine

Negation: only while the fd is registered with the runtime poller. `os.Pipe` falls back to blocking mode when `SetNonblock` or poller registration fails; then `SetReadDeadline` returns `os.ErrNoDeadline` and `Close` does not unblock the `Read`. The join then either waits forever (hang returns) or abandons the goroutine (race on the stdout buffer, fd never closed). Scenario: a container runtime with a restricted epoll/kqueue budget; the CLI hangs on that host alone and the refusal that "is never withheld" is withheld.

---

### Claim 3: a drain blocked past the bound has consumed every byte

Negation: the claim is about BLOCKED; the mechanism is about ELAPSED. Pipe FIFO guarantees delivery before EOF only if the reader keeps reading. A late reader with bytes pending is refused as held. Scenario above (P-2). The inversion is exact: the approach refuses in a case where the whole read was achievable within a few milliseconds more, so "bounded wait outranks whole read" is being applied where no conflict existed.

---

### Claim 4: deadline-first classification needs no timeout class in the binding

Negation (two ways). First, the flip: exit at t, hold detected at t + WaitDelay; class depends on whether that crosses the declared timeout. Second, the mechanism: if the executor classifies by `errors.Is(bindingErr, context.DeadlineExceeded)` rather than by `ctx.Err()`, a typed exec error that does not wrap the context error hides the timeout. The brief states the rule as a property of the executor; S1 says trace it. Scenario: child killed at the deadline, grandchild holds stdout, binding returns `ErrPipeHeld{stdout}` with no wrapped ctx error; executor sees no DeadlineExceeded and classes `execution_failure`; the author's timeout branch never fires.

---

### Claim 5: 500 ms suffices for a SIGKILLed group member to release its pipes

Negation: SIGKILL is delivered at once; fd release happens at exit, and exit waits for the process to leave the kernel. D-state on a network filesystem, a FUSE mount, or a driver ioctl; teardown of a multi-GiB address space; a stopped-and-traced helper. Each pushes release past 500 ms on a loaded runner. Scenario: a wrapper backgrounds a `java` indexer with inherited stdio on an NFS-backed CI; the group kill lands; the JVM spends 700 ms in `exit_mmap`; the join refuses a correct answer as held. The prior point-fix's success case regresses to a refusal on that host.

---

### Claim 6: the total bound is timeout + 2·WaitDelay

Negation: (a) two drains joined one after the other with one timer each is 3·WaitDelay; (b) a `reapGroup` that polls for group death has its own unstated bound; (c) `Cmd`'s WaitDelay timer starts after `Cancel` returns; (d) `Process.Wait` after SIGKILL still waits for the child to actually exit, unbounded for a D-state direct child; (e) the write journey is two invocations. Scenario: both pipes held, sequential joins, `reapGroup` polls a D-state member for 2 s: return at timeout + 3.5 s against a documented + 1 s.

---

### Claim 7: precedent establishes "bound wins, parent closes, caller sees an error"

Negation: Go's precedent establishes "bound wins, parent closes, caller sees an error AND the bytes". `Cmd.Output()` returns both. The disposition of whether to use the bytes is the caller's, and at least one prominent Go consumer of subprocess output that met this exact git/ssh/gpg-helper failure is known to treat `ErrWaitDelay` on a zero exit as success with output. The two CLIs the brief cites are unnamed. Until the resolve stage quotes their handling code, this is S3-class borrowed authority: a precedent cited for a disposition it may not contain. The precedent supports the BOUND; it does not by itself support REFUSE over ACCEPT.

---

### Claim 8: reusing WaitDelay makes the predecessor's wording true

Negation: it makes the wording true of a different timer. `cmd.WaitDelay` still does not touch manual pipes; the drain bound is a second timer that happens to read the same constant. The wording should say "the same 500 ms value bounds Cmd's stdin copy and the spawn's drain join, as two timers". Coupling risk as in P-10.

---

### R1: accept partial output on the non-deadline path

Attack on the rejection: the rejection is sound for stdout parsing, but the sentence about read-back confirming a write conflates two invocations. The write's own side effect is not "unproven" when the child exited 0; the OUTPUT is. Refusing the write invocation with a class that means "did not run" invites a retry of a completed write. The rejection survives; the refusal's Detail and the class semantics must carry `exited=true, exitCode=0, stdout held` so the caller can branch correctly.

---

### R2: widen the kill reach

Attack: "no portable primitive" is inconsistent with a seam already forked by build tag; the surviving reason is "fourth mechanism patch, precedence still unstated". Fine as rejected; fix the reason.

---

### R3: hand the pipes back to Cmd

Attack: the measurement is stale (P-11). The ordering argument (release the group BEFORE the drain bound starts) stands and is the real reason R3 loses. Note that the approach re-implements `Cmd.WaitDelay`'s pipe-closing behaviour by hand with a different timer start; the RDR should say so plainly so a reviewer does not ask why `Cmd` is not used.

---

### R4: status quo plus documentation

Stands.

---

### R5: abandon the drains with a timer

Attack on the approach from R5's own objection: if the join uses `SetReadDeadline` only and never `Close`s, or abandons the goroutine after the deadline without waiting for it, the approach has R5's leak (fds, goroutines) plus a data race. The approach must be two-phase: deadline, wait for goroutine return, close.

---

### R6: pty

Stands.

## 3. Consumer artifacts

Each names the test or journey that would have caught the finding at review time.

- P-1: `TestSpawn_StdinHeldByEscapedGrandchild_ReturnsWithinBound` in `internal/cli/cmdbind`: write binding, payload 4× the pipe buffer, child `setsid`s a grandchild that inherits stdin and sleeps, child exits 0 without reading; assert return ≤ timeout + 2·WaitDelay + 100 ms, typed error names `stdin`, `goleak.VerifyNone`. Companion: `TestSpawn_WaitReturnsErrWaitDelay_StillReapsAndJoins` if stdin is Cmd-owned.
- P-2: `TestSpawn_LateDrainNoHolder_Accepts`: test hook stalls the stdout drain goroutine for 2·WaitDelay after the child exits with a 40 KiB envelope and no holder; assert acceptance and whole output. Companion `TestSpawn_ContinuousEscapedWriter_Refuses`: escaped grandchild writes a byte every 10 ms forever; assert refusal (the drain is never blocked yet never reaches EOF). The rule under test: refuse on "no EOF after a final bounded read", not on elapsed time.
- P-3: `TestWriter_HeldStdoutAfterExitZero_DetailReportsExit`: assert Detail contains the exit status and the phrase "command exited 0; stdout held open by a process outside its group". User journey "agent retries a refused write": the journey document must state what the agent should do on this refusal.
- P-4: `TestExecutor_HeldPipeClassification` table: declared timeout ∈ {200 ms, 1 s, 5 s} × child exits 0 at 50 ms with an escaped holder; assert Detail is byte-identical across rows and class follows the documented rule. Plus `TestExecutor_DeadlineFirst_TypedErrorWrapsContext`: child killed at deadline with holder; assert class `timeout`.
- P-5: four tests, one per binding, sharing a fixture that prints a VALID envelope, `setsid`s a holder, exits 0: `TestReader_ValidEnvelopeThenHeldPipe_Refuses`, `TestGate_...`, `TestWriter_...`, `TestReadBack_...`. Plus `TestSpawn_NonZeroExitWithHolder_ReapsAndJoins` (exit 3) and `TestSpawn_SignaledWithHolder_ReapsAndJoins` (SIGTERM'd child), both under `goleak`.
- P-6: `TestSpawn_PipesArePollable`: assert `SetReadDeadline(time.Time{})` returns nil on both pipe read ends right after `os.Pipe`; spawn refuses pre-start with a typed error on `os.ErrNoDeadline`. Race test: `go test -race` on the held-pipe fixture. `TestSpawn_HeldPipeRefusal_ClosesFds`: 200 held-pipe invocations; open-fd count of the test process stable.
- P-7: `TestSpawn_SetMBackgroundHelper_RefusedWithRemediation`: wrapper `set -m; helper & ; exit 0`; assert Detail includes the remediation sentence. CI note: run the non-escaping helper success case under `stress`/`nice` load in the flake sweep; record the p99 of release-after-SIGKILL to justify 500 ms.
- P-8: `TestSpawn_BothPipesHeld_ReturnsWithinBound`: both drains held, ctx deadline hit; assert ≤ timeout + 2·WaitDelay + 100 ms. `TestReapGroup_DoesNotBlock`: reapGroup returns within 50 ms with a member that ignores nothing but is SIGSTOPped. `TestCancel_IsSignalOnly`. Journey doc states the write journey bound is 2×.
- P-9: resolve-stage evidence file quoting the exact `ErrWaitDelay` handling lines of each named CLI, with the branch (surface vs tolerate) marked.
- P-10: a single `const drainBound = waitDelay` with a doc comment naming both timers; `TestSpawn_DrainBoundEqualsWaitDelay` pins the equality so a later tuning is a visible diff.
- P-11: replace the figure with a measurement under the current WaitDelay or drop the number.
- P-12: `TestSpawn_EscapedWriterAfterRefusal_GetsEPIPE`: after refusal, the escaped grandchild attempts a write; assert EPIPE/SIGPIPE. Residue clause restated accordingly.
- P-13: wording fix; no test.
- P-14: wording fix; no test.

## 4. Refutation targets (seed classes found here)

- S1 (order claim stated as property): P-5. "Never parsed on any path" and "the order stays Wait → reapGroup → join" are both properties; the early-return structure on a non-nil `Wait` error and each binding's (inv, err) precedence must be traced, not asserted. Owed at propose: the approach should name the (inv, err) contract of spawn ("on a typed error the invocation struct is nil") so the property is structural.
- S2 (cardinality premise borrowed, never checked against the producer): P-1 and P-8. "One spawn, one join, one bound" borrows "one" three times; the producer has three pipes, two drains, and a write journey of two invocations.
- S3 (borrowed authority): P-4 ("existing deadline-first rule"), P-9 (unnamed CLIs; Go precedent returns the bytes), P-11 (stale measurement).
- S4 (answered for the reachable subset): P-1 is the exact class, at the exact seam, for the second time. The fixture holds stdout and stderr; stdin is the unreached sibling. Also P-2 and P-7 are the "reachable subset" of the bound argument: the fixture's holder is a sleeping escaped process; a late drain and a slow-dying group member are the shapes it does not exercise.

## Disposition

The precedence decision (bounded wait outranks whole read) survives every negation; nothing here forces a switch to a rejected alternative. The approach does not survive as written: P-1 (stdin carrier unbounded) reproduces the RDR's own symptom, and P-2 (elapsed vs blocked) refuses correct answers under load. Both fold into the approach without changing the decision. P-3, P-4, P-5, P-6, P-8, P-12 are contract sentences and tests the propose text must add so the refine and resolve stages have something to check.
