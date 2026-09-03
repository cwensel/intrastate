Model: claude-sonnet-5

# Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | 0026:A9 | The record locks the entire per-read re-arm mechanism (the thing that makes MVV row 6 pass at all) on a spike Status marked `Pending` — never run. If the spike is skipped or fudged under schedule pressure, the mechanism ships unmeasured. | A large answer (near or above one pipe-buffer's worth streamed slowly) is refused as a held pipe on a loaded CI runner even though the writer has already exited — a false `execution_failure` on a command that actually succeeded. | premortem, §critical-assumptions, AT-1 |
| C-2 | 0026:A10 | The pollability probe (`SetReadDeadline` + immediate clear, on every pipe, on every invocation) is asserted "safe and detecting" with Status `Pending`. No spike has exercised a genuinely non-pollable `os.Pipe` on either supported OS — F5 admits the trigger is "unexercised for a real os.Pipe." A probe that is not actually a no-op perturbs every command invocation the CLI runs, not just the held-pipe edge case. | Ordinary, correct command reads on some class of host (containerized/restricted epoll, or a future Go runtime change) begin behaving differently — from added latency to spurious deadline errors — for every invocation, not just the ones that hit the drain bound. | §critical-assumptions, AT-2 |
| C-3 | 0026:F6 / 0026:D-the-drain-grace-is-a-mechanism-constant-not-the-rejected-knob | `DrainGrace` (and the entire stderr-tail preservation machinery) is justified by one sentence — "what the grace preserves is the stderr tail... the text that tells a user WHICH helper held the pipe" — and then F6 admits that exact text is unreachable in the RDR's own motivating scenario (declared timeout also elapsed → class `timeout`, `Detail` empty by 0025:C4's rule). The mechanism is built to solve the problem statement's own example and then the record concedes it does not solve it there. | The author who wrote the bug report this RDR exists to fix — "I declared `timeout=2s`, my tool escaped, I got nothing to branch on" — now gets a bounded `timeout` refusal with **no reason attached**. They are told their command was slow. It was not; it was structurally broken. They cannot tell the difference without re-running under a longer timeout, which does not help, because the helper is broken rather than slow. | premortem, §failure-modes, §consequences |
| C-4 | 0026:§consequences (Negative bullet, "UNSURVEYED" population) | The record ships a behavior change — permanently-escaping-but-correct tools now refuse where they used to hang — against a population the record itself calls unsurveyed: "no reader in the repo's own fixtures exhibits the escaping-but-correct shape, and whether any reader in the wild does is unknown." The mitigation offered ("the refusal names what to fix") is stated to not even apply to the likely-real case: a third-party daemonizing binary the author cannot change. | An existing, working `command` reader that wraps a tool which intentionally daemonizes (log shippers, indexers, watchers spawned by build tools) starts failing in production after this ships, with no way for the author to fix the tool itself — only to wrap or redirect its stdio, which requires them to first understand that's the problem. | §trade-offs, premortem |
| C-5 | 0026:D-naming / 0026:C1 bound: | `WaitDelay` is deliberately reused as both `Cmd`'s stdin/post-cancel timer AND the new drain-join timer, on the stated theory that "a future tuning moves both." This is exactly backwards for the two populations that actually need different values: CI runners under scheduler pressure (want a *larger* drain bound, A5's own data shows tail growth of 15x from idle to loadavg 47) versus wanting a *tight* stdin-write bound so a non-reading child fails fast. A single future PR that "fixes flaky CI" by bumping `WaitDelay` silently widens the drain-join bound too, and nothing in the type system or a test pins them apart except a comment. | A future engineer raises `WaitDelay` to fix CI flakiness on the stdin path; the drain-hold refusal window silently widens system-wide, and nobody notices until a real hang (this RDR's original bug) takes longer to surface than before, or a genuinely broken command now appears to "work" longer before refusing. | §load-bearing-decisions (D-naming), premortem P-10 |
| C-6 | 0026:C1 refusal: / 0026:A8 | The applied-sense gain is a single `errors.As` match site inside `executor.go::Write`'s `err != nil` arm, landing "Applied()" derivation for exactly one of several failure shapes on that arm (held-pipe). Every other error shape reaching that same arm (a plain non-typed exec failure, a resolveArgv0 failure, an os.Pipe creation failure) still reports `Applied()` false with no distinguishing signal, and the RDR does not show how the executor tells "definitely did not run" from "held pipe, might have run" once both produce the same Go error type family (`*accessor.ExecError`-wrapped). A miswired `errors.As` (matching too broadly or too narrowly against the held-pipe type) is a one-line diff that silently flips a real non-applied write to "may have applied," or vice versa, with no compiler signal and (per QA M2, only partially closed) an underspecified rendering assertion. | An agent caller sees "may have been applied" on a write that clearly never touched the target (e.g., `resolveArgv0` failed before `cmd.Start()`), and does not retry a write that was safe to retry — or the reverse, retries a write that already landed. | §technical-design, 3amigo M2/H-a, AT-3 |

# 1. The three most likely ways implementation goes wrong

## 1a. A9's mechanism ships unmeasured, and the drain-join re-arm is wrong under real load

**Root cause.** The single mechanism that makes the RDR's headline promise ("a whole answer is never falsely refused") actually work — a **per-read** deadline re-arm inside the drain loop, rather than a one-shot deadline set by the parent before `drains.Wait()` — depends on an assumption (A9) that is marked `Pending` at the point this document was reviewed. The 3amigo QA pass (H2) caught this exactly: nothing in the record measures whether a 1 MiB buffered tail, competing against a ~64 KiB pipe buffer, drains inside repeated 50 ms grace windows on a loaded machine. A2 measured 24 bytes. A3 measured an unbounded grace. Neither is A9's claim.

**The specific passage that enabled it.** `0026:C1` `precedence:`: "The deadline is re-armed at `now + DrainGrace` BEFORE EACH read of the final drain, so the grace bounds the IDLE GAP between reads and never the size of the tail." This is asserted as settled contract language — Normative — while its supporting assumption (A9) is `Pending` two sections earlier. The Assumption Verification gate rule the record itself cites ("no assumption marked Pending... may have settled-fact prose elsewhere in the RDR depending on it") is precisely what this passage violates, unless A9's spike lands clean before lock. Nothing in the current record forces that ordering; it is a documentation convention, not a mechanical gate.

**The symptom the user will see.** On a CI runner under real load (the same runners A5 measured drain-goroutine scheduling latency reaching ~28 ms at 4x core saturation), a large, entirely successful command output gets refused as `execution_failure`, `Detail` naming a held pipe that was never actually held — just slow to schedule. This is a **regression relative to today's unbounded-but-correct behavior**: today, a large answer under load eventually completes; after this ships, under exactly the same load, it can be wrongly refused. The MVV's own "late-but-unheld" row (row 6) is designed to catch this, but the row's own text admits "unmeasured until that spike runs" — meaning the test that would catch this bug at review time is not run yet, either.

## 1b. F6's diagnostic apparatus is dead code for the case that motivated it

**Root cause.** The RDR's Problem Statement is built entirely around one journey: an author declares `timeout = "2s"`, the tool escapes, and the CLI hangs "unboundedly... discover[ed] as a stuck invocation, not as a refusal." The fix delivers the bound. But the entire secondary apparatus this record builds — `DrainGrace`, the per-read re-arm, the held-pipe `Detail` line naming which pipe and why — is **only rendered on the `execution_failure` class**, and F6 states plainly that when the declared timeout has *also* elapsed (which is true in the record's own motivating example, since the escape is discovered by timing out), the executor's pre-existing deadline-first rule (0025:C4) means the class is `timeout`, which "carries no Err/Detail." So the exact author from the Problem Statement gets a bounded refusal with the reason stripped.

**The specific passage that enabled it.** `0026:D-the-drain-grace-is-a-mechanism-constant-not-the-rejected-knob`: "what the grace preserves is the stderr tail that `wrap` carries into `Detail` — the text that tells a user WHICH helper held the pipe." Read against `0026:F6`: "the diagnostic apparatus this record builds... [is] unreachable in exactly that case [the record's own motivating scenario]." Both are true; they are in the same document, about the same code path, and the record's own 3amigo pass (H-e) flagged this as a live tension between two elements citing each other in only one direction, only partially fixed (D-9, non-blocking, closed by adding a test on the *other* leg rather than by fixing the asymmetry itself).

**The symptom the user will see.** An author who files exactly the bug this RDR is written to close — "my declared timeout fired, but I have no idea why, my tool just seems slow" — files it again, post-ship, on the *identical* symptom, because the class they see (`timeout`, empty detail) still gives them nothing to act on beyond "raise the timeout," which for a structurally broken (escaping) tool does not help even once.

## 1c. The applied-sense re-key on the write path is a single easy-to-miswire `errors.As` site with a thin, partly-unverifiable test

**Root cause.** Phase 4 lands the entire A8/C1 `refusal:` "may have been applied" gain in one `errors.As` match inside `executor.go::Write`'s `err != nil` arm — a single site, correctly scoped per the cove pass's source-level verification (Q3/Q4/G6/G7 confirmed the read-back leg is untouched and already correct). But the test that is supposed to prove this landed (S1, direct-write leg) was flagged by the record's own QA persona (M2) as asserting a rendered CLI string with "no named observable: no exit code, no output stream, no stable substring, no golden-file convention" — and the 3amigo consolidation records this as only *partly* resolved, not closed.

**The specific passage that enabled it.** `0026:C1` `refusal:`: "The DIRECT write invocation refuses `execution_failure`, and it is the leg that needs the gain... So one site changes, not two." Combined with Testing Strategy `S1`: "on BOTH write legs `Applied()` is true AND the CLI renders the 'may have been applied' text" — a conjunction where only the `Go`-level half (`Applied()`) has a clean assertion surface; the render half does not.

**The symptom the user will see.** Either (a) the re-key silently fails to land in a refactor because the test that should catch it (S1's rendered-string half) has no fixed assertion surface and drifts to pass trivially, and a caller who could safely retry a write is told nothing about the applied risk — or (b) an over-broad `errors.As` match later catches an unrelated exec failure that never touched the target and tells the caller "may have been applied" when it certainly did not, causing them to treat a safe-to-retry failure as unsafe.

# 2. The one section that will be rewritten within 6 weeks of shipping

**`0026:F6` / the interaction between `C1 class:` and `0025:C4`'s "timeout carries no Detail" rule.**

This is the record's own self-diagnosed weak point, twice caught by independent review passes (PM persona H-e, and the record's own admission in F6 that this is "the record's motivating scenario" going undiagnosed) and never actually fixed — only reclassified from "under-explained" to "accepted trade-off." An accepted trade-off is still a trade-off, and the cost here lands on exactly the user this RDR was written to serve. The stated reason it isn't fixed now — "the cure is to amend 0025:C4's `detail:` rule... which is this record's Overrides target... and because the sub-reason constraint is homed in JDR 0003 §D1, not here" — is a real, structurally sound reason to *not* fix it in this record. But it means the fix is a known, named, deferred amendment sitting one hop away in JDR 0003 §D1 or a future 0025 successor. The first real user who hits the "escaped grandchild + timeout also elapsed" case — which, per the record's own framing, is the *common* case, not an edge case — will re-open exactly this question, and whoever owns JDR 0003 §D1 at that point will be asked why the sub-reason doesn't ride on `timeout`. That's the rewrite.

# 3. The one assumption that will not survive first contact with a real user

**A5 / the implicit assumption that 500 ms (`WaitDelay`) is a stable constant across the CLI's real deployment surface.**

A5 is marked `Verified`, with real measurement (160 trials, worst case 27.7 ms at loadavg 47.5 on 12 cores — 5.5% of the bound). That's good, rigorous work. But it was measured on a 12-core dev/CI-shaped machine. This RDR reuses that single constant for three purposes at once (D-naming: Cmd's stdin write bound, Cmd's post-Cancel wait, and now the drain-join bound), on the explicit theory that "a future tuning moves both" together. The moment this ships into a real agent-driven pipeline running on a resource-constrained container (a common `--allow-commands` deployment shape: a sandboxed agent runtime with tight CPU/memory limits, likely *more* oversubscribed than A5's loadavg-47 test, not less), the margin the record is proud of (18x more tail growth absorbed, per F2) gets consumed by a host class the record never measured. The record's own F2 admits: "the margin is a function of runner oversubscription rather than a constant... a far more oversubscribed CI runner is the condition that would erode it" — and does not name a floor below which the design breaks, only that A5 didn't observe one in its own test range. First contact with a real, resource-starved deployment is exactly the condition A5 didn't test, and the single shared constant means a fix for one symptom (stdin write too aggressive) can't be applied without also moving the drain bound (C-5's finding, restated here as the first-contact case).

# 4. The premortem

RDR 0026 shipped six weeks ago. `TestSpawn_EscapedGrandchildHoldsStdout` is green on darwin and linux, and the demo everyone ran at the design review — a `setsid` grandchild holding stdout past a 2-second declared timeout — now returns cleanly with a `timeout` refusal at 2.06 seconds instead of hanging. The launch was called a success.

Week two: a support ticket comes in from a user running `flow resolve` inside a memory-constrained Kubernetes pod against a `command` reader that wraps a legacy indexing tool. The tool is fine — it always was — but it backgrounds a small log-shipping helper with `set -m` before exiting, a pattern common enough that three other `command` readers in the same fleet do the same thing. Before this RDR shipped, that helper's inherited stdout pipe kept the parent's read blocked for a few hundred milliseconds until the group signal reaped it (A5's own in-group case), and the CLI returned correctly. Now, under this pod's tighter CPU quota, the group kill and the drain join race the way A5 never modeled: cgroup throttling delays SIGKILL delivery, `reapGroup`'s signal lands late, and the drain-join's 500 ms `WaitDelay` — the same constant now shared by three unrelated timers, as C-5 and D-naming call out — expires before the helper's fd is released. The read that used to succeed at ~200 ms now refuses `execution_failure`, `Detail` naming "stdout held past bound, exited 0." The on-call engineer reads `0026:F2`'s own text back to the ticket: "the tail is load-sensitive... a far more oversubscribed runner is the condition that would erode it." That is the condition. No amount of retrying fixes it, because the constant, not the tool, is wrong for this host class.

Week three: an agent operator files a second ticket. Their model author declared `timeout = "5s"` against a wrapper for a build tool known to occasionally daemonize a watch process if a flag is misconfigured. The escape happens; the CLI refuses at 5.06 seconds with class `timeout` and an **empty** `Detail` — exactly F6's admitted gap. The operator's dashboard shows "timeout: command was slow," raises the declared timeout to 30 seconds hoping that fixes it, re-runs, and gets the same wrong diagnosis 30.06 seconds later, because the tool wasn't slow — it was structurally broken, and the one piece of text that would have told them so (the held-pipe `Detail`, engineered specifically by `DrainGrace` at real design cost) never rendered, because the class flipped to `timeout` before it could. They eventually find the held-pipe reason only by reading the CLI's source, not its output. The RDR's own words are in the postmortem verbatim: "the diagnostic apparatus this record builds... [is] unreachable in exactly that case."

Week four: a roborev-style triage of the write path turns up a case where `executor.go::Write`'s `err != nil` arm classified a `resolveArgv0` failure (child never even started — a config typo in the declared `command` binding) as `Applied()` true, because a later refactor widened the `errors.As` match introduced by this RDR's Phase 4 to catch a broader error family than the held-pipe type alone. The rendered CLI text tells a write-caller "the write command already ran, so the mutation may already be applied" for a write that structurally could not have run. An agent, trusting the class, skips a safe retry and instead attempts a read-back to "confirm" a write that never happened, gets a `read_back_incomplete` on a target that was never touched, and escalates a false anomaly. `0026:S1`'s test, the row meant to catch exactly this, is later found to assert only `Applied()` at the Go level in CI — its "renders the text" half was never pinned to an observable, exactly as the record's own QA persona (M2) warned at propose time, and the gap survived to Final because nobody blocked lock on it.

None of these three incidents individually falsify the record's central decision (bounded wait outranks whole read — that holds, and should). All three are places where the record's own review process found the exact defect, in writing, before lock, and chose to carry it forward as an accepted risk, a Pending spike, or a partly-closed test gap rather than a blocking fix.

# 5. Acceptance tests that would have caught each failure at RDR-review time

**For 1a / A9 (re-arm under load):**
```
Scenario: Large buffered tail survives a re-armed drain grace under CPU load
  Given a child that writes 1 MiB to stdout and exits 0 with no grandchild
  And the drain goroutine's start is delayed past the join timer (test hook)
  And the test process is run under simulated CPU load (stress/nice, matching A5's loadavg-47 condition)
  When the drain join fires its bound and re-arms the per-read deadline
  Then the invocation returns the whole 1 MiB
  And no held-pipe refusal is produced
  And this passes on both darwin and linux before A9's Status may read "Verified"
```
This is exactly MVV row 6 plus A5's load condition — the gap is that MVV row 6 as written has no load condition, and A9 says Pending. Requiring the row to run *under load*, not just on an idle machine, before Status flips to Verified would have caught the gap the ledger names as C-1.

**For 1b / F6 (diagnostic unreachable in the motivating case):**
```
Scenario: The record's own motivating journey gets a diagnosable refusal
  Given a model author has declared timeout = "2s" on a command reader
  And the child spawns a setsid grandchild holding stdout/stderr and sleeps past the deadline
  When flow resolve invokes the reader and the deadline elapses
  Then the refusal class is "timeout"
  And the refusal's rendered text names the held pipe(s) and the remediation
    (not merely "carries no Err/Detail" per 0025:C4's existing rule)
```
This scenario currently **cannot pass** against the shipped design — F6 says so explicitly. Writing it at review time (rather than writing S6, which asserts the *current*, diagnosis-losing behavior as correct) would have forced the joint-decision with JDR 0003 §D1 to happen before lock rather than after a support ticket.

**For 1c / A8 applied-sense (thin rendering assertion):**
```
Scenario: A held-pipe write renders a stable, greppable "may have applied" marker
  Given a write command exits 0 with its stdout held by an escaped grandchild
  When the executor classifies the refusal
  Then Applied() is true
  And the CLI's rendered envelope contains the exact constant `detailMayHaveApplied`
    at a named, stable field (not prose matched by substring)
  And a config-error write (argv0 resolution failure, no process ever started)
    renders Applied() == false through the SAME code path,
    proving the errors.As match does not over-broaden
```
The second half of this scenario — the negative control on a pre-start failure — does not exist anywhere in the current Testing Strategy (S1–S10). Its absence is exactly how C-6's miswiring risk would go undetected.

**For 2 / the shared-WaitDelay constant (§3, C-5):**
```
Scenario: Tuning the stdin bound does not silently move the drain bound
  Given WaitDelay is changed from 500ms to a new value for stdin-related reasons
  When the drain-join bound is inspected
  Then a test fails unless the change was accompanied by an explicit acknowledgment
    that the drain bound moved too
```
No such pinning test exists; `0026:D-naming` treats the coupling as a feature, and no test in S1–S10 fails if a future PR moves `WaitDelay` for one reason and silently drags the other along.

**For 3 / A5 generalizing beyond the tested host class:**
```
Scenario: The bound holds on a CPU-throttled container, not just a bare-metal/CI runner
  Given the in-group grandchild-leak oracle
  When run inside a cgroup-throttled environment (e.g. 0.5 CPU quota, matching a
    realistic --allow-commands sandbox deployment)
  Then the p99 release-after-SIGKILL time is measured and compared against WaitDelay
  And the record states a floor host class below which the bound is not claimed to hold
```
A5's 160 trials never include a CPU-quota-throttled condition; this is the test that would have forced the record to either measure it or explicitly scope the claim to "un-throttled hosts only," which it currently does not do.

