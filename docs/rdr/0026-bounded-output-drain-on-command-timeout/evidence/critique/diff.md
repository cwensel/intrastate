Model: claude-opus-5
Passes diffed: critique.md (claude-opus-5) x critique-modelB.md (claude-sonnet-5)

# Cross-model critique diff — RDR 0026

Pass A = claude-opus-5 (15 ledger rows). Pass B = claude-sonnet-5 (6 ledger rows).
IDs are per-file and do not correspond across files; every row below was matched by
RDR passage anchor first, then by underlying defect.

---

## Merged findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Raised by | Confidence |
|----|-------------|--------------|-------------------|-----------|------------|
| D-1 | `0026:C1` (`precedence:`) | The load-bearing mechanism — "the deadline is re-armed at `now + DrainGrace` BEFORE EACH read of the final drain" — has no site. `cmdbind.go:352-364` `readBounded` is `io.ReadAll(io.LimitReader(r, limit))` + `io.Copy(io.Discard, r)`; there is no per-read loop the re-arm can occupy. Implementing C1 as written is a full rewrite of the function whose semantics C1 `whole:`, the `fidelity` mini-check and S3 all pin as unchanged. | A cap-boundary shift under the rewritten loop: a reader emitting exactly 1 MiB is refused as overflow, or cap+1 is accepted. S3's "existing suite passes unmodified" fails and the implementer relaxes the assertion instead of the code. | A | high (grounded: `internal/cli/cmdbind/cmdbind.go:352-364`) |
| D-2 | `0026:C1` (`bound:`) | "a command write is two invocations (write, then read-back)" is false. `(*Executor).Write` runs THREE bounded invocations through the same reader: pre-write baseline read `executor.go:320`, `binding.Apply` `:355`, read-back `:397`. Worse, the baseline leg's refusal is discarded at `:322-326` (`baselineUnread = slices.Clone(protected)`), so it costs a full bound and yields no signal. The `disposition` mini-check enumerates only the direct and read-back legs. | An agent sizing its watchdog from the published `timeout + 2·WaitDelay` is killed at roughly 3x it, mid-write, and cannot tell whether the mutation applied — the exact discrimination this RDR exists to restore. | A | high (grounded: `internal/accessor/executor.go:320`, `:326`, `:355`, `:397`) |
| D-3 | `0026:A9` (Pending), `0026:MVV` row 6, `0026:C1` (`precedence:`) | A9 is Pending and is the only evidence for the per-read re-arm at scale (A2 measured 24 bytes; A3 measured an unbounded grace). C1 `precedence:` states the re-arm as settled normative text and MVV row 6 is marked "(Normative)" on top of it. The record's own `G-assumptions` gate forbids exactly this — and that gate is an unfilled template block. | A large answer on a loaded runner is refused as a held pipe on a pipe with no writer at all. Ordinary large-output readers become intermittently `execution_failure`; retry works. A regression relative to today's unbounded-but-correct behaviour. | both | high (grounded: A9 Status `Pending`; MVV row 6 text carries "(Normative)"; `G-assumptions` is verbatim template) |
| D-4 | `0026:F6`, `0026:D-the-drain-grace-is-a-mechanism-constant-not-the-rejected-knob`, `0026:§consequences` | `DrainGrace` — a new constant and an explicit divergence from Go and both peer CLIs — is justified solely by preserving the stderr tail "that tells a user WHICH helper held the pipe," and F6 concedes that text is unreachable in the record's own motivating scenario (declared timeout also elapses ⇒ class `timeout` ⇒ `Detail` empty per 0025:C4). The mechanism is justified by a benefit its primary journey cannot deliver. | The author who filed the originating bug gets `flow-accessor-timeout` and nothing else. They raise the timeout to 30 s, wait 30 s, get the same message. The hang is replaced by a longer lie. | both | high (grounded: F6 text; `flow_exec.go:399-408` `ClassTimeout` arm emits only "the accessor `X` timed out") |
| D-5 | `0026:§consequences` (UNSURVEYED negative) | The record ships a behaviour change converting an unknown population of correct-but-escaping readers from working to refused, states the population is UNSURVEYED ("whether any reader in the wild does is unknown"), and concedes the mitigation does not apply to third-party daemonizing binaries. No phase surveys it; no opt-out or escape hatch ships. | A reader wrapping a tool that daemonizes (`ssh -f`, `set -m`, Node `detached: true`) works Monday, hard-fails Tuesday, with a remediation the author cannot perform because they do not own the binary. | both | high (grounded: `§consequences` third bullet; no phase in `§implementation-plan` surveys) |
| D-6 | `0026:A6`, `0026:C1` (`bound:`), `0026:S8` | The measured per-leg figure already breaches the stated bound (drain join 500.1-508.2 ms against a 500 ms `WaitDelay`, 40/40 samples). The repair moves the claim to the TOTAL and leans on leg (a) donating its unused allowance — while conceding in the same clause that a child which RESISTS termination spends leg (a)'s allowance. A6: "that case was not exercised." S8 is planned, not run. The total is Normative. | A resisting child plus an escaped grandchild returns past `timeout + 2·WaitDelay`. The contract's one hard number is violated on the composite case, in production, first. | A | high (grounded: A6 Evidence + Narrows clause; S8 is unrun) |
| D-7 | `0026:A10` (Pending), `0026:F5`, `0026:S5` | The pollability probe runs on EVERY invocation's hot path and is unverified. F5 concedes the trigger is "unexercised for a real `os.Pipe`" — no naturally non-pollable pipe was constructible on either supported OS. A10's own "If wrong" is that the probe perturbs ordinary reads and the byte-for-byte claim fails for every command. The fallback IS the hang ("on that host the join stays unbounded"). | Best case, dead code on every command. Worst case, added latency or spurious deadline errors on every invocation on some host class, not just held-pipe ones. | both | high (grounded: A10 Status `Pending`, Method Spike; F5 "unexercised for a real `os.Pipe`") |
| D-8 | `0026:C1` (`precedence:`), `0026:A2` | The `graceOn` handoff is an unsynchronised cross-goroutine mutation. Nothing in C1, the `authority` mini-check, or the Illustrative Code names a primitive or a happens-before edge for the mark. Today `spawn` is race-free by construction: `drains.Wait()` at `cmdbind.go:302` is the only edge for the `stdout`/`stderr` slices assigned at `:275-279` and read at `:304`. RDR 0026 removes the unconditional join and never re-establishes the edge. A2 measured a parent `SetReadDeadline` on an `os.File` (internally synchronised poller) — a different object. | `-race` flags the drain mark intermittently in CI; or worse, a drain misses a stale mark, never re-arms, blocks forever, and `drains.Wait()` never returns — the original hang, now behind a timer that reads as correct. The MVV fixture cannot catch it (its child is silent, so the drain is blocked in its FIRST read). | A | high (grounded: `cmdbind.go:275-279`, `:302`, `:304`; A2 scope) |
| D-9 | `0026:A5`, `0026:F2` | A5's 500 ms margin is measured on Go helper processes on a 12-core dev/CI machine that exit promptly to SIGKILL. It measures scheduler and pipe-teardown latency, not process-exit latency, and never a CPU-quota-throttled container. F2 bounds the risk with the same 27.7 ms number and scales it only by "runner oversubscription," never by teardown cost or cgroup throttling. Premortem P-7 (uninterruptible sleep on NFS/FUSE, large-heap teardown) is folded nowhere. | Intermittent `execution_failure` on a reader whose wrapper backgrounds a JVM or a FUSE-touching helper, or that runs in a memory/CPU-constrained pod. Works locally, fails on CI, class says "could not be executed." | both | high (grounded: A5 Evidence names 12 cores, loadavg 16.2/47.5, Go helper; F2 text) |
| D-10 | `0026:A8` (Pending), `0026:§phase-4-executor-and-cli-applied-sense`, `0026:S1` | The only user-visible write-safety gain — the `Applied()` re-key so a held-pipe write is not read as "the write did not happen" — sits behind a Pending assumption, in the second-to-last phase, in a different package from every other phase. `flow_exec.go:410-412`'s `ClassExecutionFailure` arm has no phase check and sets no `Detail`; `Applied()` (`model.go:323`) has zero non-test callers. No partial-ship disposition is stated. | Phase 4 slips. A command write applies its effect, exits 0, gets `execution_failure` with no applied sense. The agent's branch on that class is retry. The mutation is applied twice. | A | high (grounded: A8 Status `Pending`; `flow_exec.go:410-412`; `.Applied()` callers are all in `adversarial_0004_test.go`) |
| D-11 | `0026:C1` (`refusal:`), `0026:A8` | The `errors.As` match that derives the applied sense lands on the shared `err != nil` arm at `executor.go:367-371`, which every other pre-mutation failure shape also reaches (`resolveArgv0` failure, `os.Pipe` creation failure, plain exec failure). The record does not state a negative control: nothing pins that an over-broadened match cannot flip a definitely-not-applied write to "may have applied." No scenario in S1-S10 exercises a pre-start failure through that arm. | An agent is told "the write command already ran, so the mutation may already be applied" for a write that structurally could not have run (a config typo in the declared binding), and skips a safe retry — or the reverse. | B | medium (grounded: the shared arm at `executor.go:367-371` is real; the miswiring is a future-refactor risk, not a present record defect) |
| D-12 | `0026:C1` (`refusal:`), `0026:§technical-design` | C1 states "The returned invocation carries no stdout; stdout read on a held drain is never parsed on ANY path," but `spawn` returns `(inv, err)` with `inv` fully populated on every error arm (`cmdbind.go:304`, `:322`, `:326-328`). The guarantee rests entirely on A1's textual observation that today's three callers check `err` before `inv`. Nothing structural enforces it. | A future gate fast path, or 0028's read-back-through-a-command-reader leg, reads `inv.exitCode` or `inv.stdout` before `err` and derives a verdict from a stream with an unreachable live writer — the silent-wrong-answer class 0025 exists to prevent. | A | high (grounded: `cmdbind.go:304` builds `inv` before the switch; `:322` and `:326-328` both `return inv, wrap(...)`) |
| D-13 | `0026:§finalization-gate`, `0026:G-assumptions`, `0026:G-proportionality` | Every gate response is an unfilled TEMPLATE guidance block. The Assumption Verification gate is verbatim the detector for D-3 ("no assumption marked `Pending` or `Unverified` may have settled-fact prose elsewhere in the RDR depending on it"); the Proportionality gate is the detector for the multi-package spread in D-15. The record carries three Pending assumptions with its own detector unrun. | Not directly user-visible; it is the mechanical reason D-3, D-7 and D-10 are still in the record at lock. | A | high (grounded: `G-assumptions` and `G-proportionality` bodies are bracketed template text, not responses) |
| D-14 | `0026:D-naming`, `0026:C1` (`bound:`) | `WaitDelay` is deliberately shared by three timers (Cmd's stdin write, Cmd's post-Cancel wait, the new drain join) on the theory "a future tuning moves both — accepted." The two populations want opposite moves: a throttled container wants a LARGER drain bound; a non-reading child wants a TIGHT stdin bound. Nothing but a comment pins them apart, and no scenario in S1-S10 fails if a future PR moves one for the other's reason. The record's own lever ("widening it is a one-line follow-up, not a contract change") is false once C1 `bound:` publishes `timeout + 2·WaitDelay`. | An engineer raises `WaitDelay` to fix CI flakiness on the stdin path and silently widens the held-pipe refusal window system-wide; a genuinely broken command appears to work longer before refusing. | B | high (grounded: `D-naming` text; `C1 bound:` publishes the constant; no S-row pins the coupling) |
| D-15 | `0026:§problem-statement`, `0026:§approach`, `0026:§decision-rationale` | The stated problem is "the CLI hangs unboundedly," which one timer solves. The record locks a precedence contract, a new constant, a new accessor-side error type, a per-invocation pollability probe, a per-read re-arm, a `readBounded` signature change, an executor re-key, a CLI rendering re-key and a docs surface, across three packages with three Pending assumptions. Reversibility is scored "one join site; revert is the old `drains.Wait()`" — false the moment `readBounded`'s signature changes and a new error type ships. | The change lands late and in pieces, with Phases 4 and 5 outstanding. The hang is fixed; the diagnostics, the applied sense and the documented residue — the parts the user touches — are not. | A | medium (Reversibility row is grounded and understated; the proportionality charge is a judgement call, and the record's own Proportionality gate was never run — see D-13) |
| D-16 | `0026:ALT1`, `0026:§decision-rationale` | ALT1 is rejected because the output's completeness "is asserted from pipe semantics, never observed" — but A3 OBSERVES byte-for-byte completeness at 1 KiB and 63 KiB with the drain reading strictly after the child was gone. The Rationale re-labels A3's fact as "what makes A's residue text honest, not as a licence to parse," which is a rhetorical move, not an argument. The real distinction (child EXITED vs child blocked in `write(2)`) is observable at the join — `spawn` already knows `werr` and `inv.exited` at `cmdbind.go:304-321`. ALT1-restricted-to-an-exited-child was never costed as a column. | A tool that answers correctly and exits 0 while a daemon holds stdout is refused, when the record's own evidence shows the answer was complete. The user loses a working reader to a purity argument. | A | medium (the A3/ALT1 tension and the unavailable-variant gap are grounded; whether the variant should have been costed is a design judgement, and A3's own antecedent clause partly defends the rejection) |
| D-17 | `0026:JC1`, `0026:F6` | The user-visible diagnostic is split three ways: 0026 owns the precedence, 0025:C4 (Implemented, this record's `Overrides` target, not amended here) blocks the `Detail`, JDR 0003 §D1 homes the sub-reason carrier. None can deliver the held-pipe diagnostic on the primary journey alone, and no schedule names which one will. | The diagnostic never ships. Each record's post-mortem points at the other two. | A | medium (grounded structurally: F6 defers to 0025:C4, JC1 homes the carrier in JDR 0003 §D1, no schedule appears; but F6's stated reason for deferring is the correct never-declare-in-both discipline, so the defect is the missing schedule, not the homing) |

---

## 1. AGREED — raised by both passes

Five defects. These are the highest-confidence findings; two models with no shared context
landed on the same passages.

| Merged | Pass A row | Pass B row | Anchor match |
|--------|-----------|-----------|--------------|
| D-3 | C-5 (`0026:A9` Pending under Normative text) | C-1 (`0026:A9` Pending, mechanism unmeasured) | Same anchor, `0026:A9`. A frames it as a gate violation (Pending assumption underwriting Normative prose, with `G-assumptions` unrun); B frames it as a schedule risk (the spike gets skipped under pressure). Same defect, two framings; A's is the sharper because it names the record's own unrun detector. |
| D-4 | C-4 (`0026:F6`, `0026:S6`, `0026:D-the-drain-grace…`) | C-3 (`0026:F6` / `0026:D-the-drain-grace…`) | Same anchors. Both independently selected F6 as the section rewritten within 6 weeks (§2 of both passes). Strongest agreement in the diff. |
| D-5 | C-8 (`0026:§consequences` UNSURVEYED) | C-4 (`0026:§consequences` Negative bullet) | Same anchor, near-identical wording of the defect. Both note the mitigation is conceded not to apply to third-party binaries; both note no phase surveys. |
| D-7 | C-6 (`0026:A10` Pending, `0026:F5`, `0026:S5`) | C-2 (`0026:A10` Pending) | Same anchor, `0026:A10`. A adds the F5 "unexercised trigger" and the "fallback IS the hang" observation; B adds the containerized/restricted-epoll host class. Complementary, not conflicting. |
| D-9 | C-9 (`0026:A5`, `0026:F2` — wrong population measured) | §3 (A5 does not generalize beyond the tested host class) | DIFFERENT anchors describing the same underlying defect: A5's margin is measured on a population that does not exhibit the failure. A names the population as process-teardown cost (JVM heap, uninterruptible sleep on NFS/FUSE); B names it as host class (cgroup-throttled container, more oversubscribed than loadavg 47). Both independently chose A5 as the assumption that will not survive first contact (§3 of both passes). The two named populations are disjoint and both real, which strengthens the finding rather than splitting it. |

Note on framing agreement beyond the ledger: both passes independently chose `0026:F6` for
their §2 and `0026:A5` for their §3. Two of the three §-level judgements converge exactly.

---

## 2. PASS-A-ONLY — raised only by the opus pass

Ten rows. All ten survive as real defects; none is judged artifact or overreach, though
three are scored `medium` in the merged ledger for the reasons given.

| Pass A row | Merged | Judgement |
|-----------|--------|-----------|
| C-1 (`0026:C1` `bound:` understates the write legs) | D-2 | REAL, missed by B. Grounded below. B never opened `executor.go::Write` and so never counted the legs. This is the single largest factual error in the record that only one pass found. |
| C-2 (`0026:C1` `precedence:` re-arm has no site) | D-1 | REAL, missed by B. B quotes the SAME `precedence:` sentence in its §1a and treats the mechanism as merely unmeasured under load — it does not notice the mechanism has nowhere to live. Grounded below. |
| C-3 (`graceOn` is an unsynchronised handoff) | D-8 | REAL, missed by B. Grounded: `cmdbind.go:302` is today's only happens-before edge for the slices assigned at `:275-279` and read at `:304`; C1 removes the unconditional join and names no replacement primitive. |
| C-7 (`0026:A6` per-leg figure breaches the bound; S8 unrun) | D-6 | REAL, missed by B. A6's own text says the per-leg claim is "strictly false" and that the resisting-child case "was not exercised," while the total remains Normative. |
| C-10 (Phase 4 applied-sense is last and Pending) | D-10 | REAL, missed by B as a *sequencing* defect. B's C-6 touches the same code (the `errors.As` site) but as a miswiring risk, not a descope risk. Kept as separate rows D-10 and D-11: they fail differently and have different fixes. |
| C-11 (`refusal:` "never parsed" is an unenforced convention) | D-12 | REAL, missed by B. Grounded below: `spawn` returns a populated `inv` on every error arm. |
| C-12 (Finalization Gate is five unfilled template blocks) | D-13 | REAL, missed by B. Verified directly: `G-assumptions` and `G-proportionality` bodies are bracketed template guidance, not responses. This is the mechanical cause of D-3, D-7 and D-10 surviving to lock. |
| C-13 (proportionality / reversibility mis-scored) | D-15 | REAL but partly judgement. The Reversibility row's "one join site; revert is the old `drains.Wait()`" is grounded-false (a `readBounded` rewrite and a new `accessor` error type are in scope). The broader "this should be split" charge is a Proportionality-gate call the gate was never run to make — which is why it is scored `medium`, not dropped. |
| C-14 (`0026:ALT1` rejected against the record's own A3) | D-16 | REAL but partly overreach. The tension is genuine and the un-costed exited-child variant is a real gap. However, A3 carries an explicit antecedent clause ("a child blocked in `write(2)` never exited"), and C1 `whole:` restates it — so the record does distinguish the two cases, it simply never turns that distinction into an ALT column. Scored `medium`: the missing column is the defect, not the rejection itself. |
| C-15 (three records each own a third of one behaviour) | D-17 | REAL but partly overreach on the framing. F6's stated reason for not amending 0025:C4 here is the correct never-declare-in-both discipline (a record declaring it in both places "would have to be demoted at 7.1"), which A acknowledges and then charges as a defect anyway. The surviving defect is narrower and real: no schedule and no named owner for the delivery. Scored `medium` and retained on that narrower ground. |

No Pass A row is excluded from the merged ledger as artifact.

---

## 3. PASS-B-ONLY — raised only by the sonnet pass

Two rows. One survives; one is partly refuted and survives only in reduced form.

| Pass B row | Merged | Judgement |
|-----------|--------|-----------|
| C-5 (`0026:D-naming` — shared `WaitDelay` across three timers) | D-14 | REAL, missed by A as a ledger row. A touches the coupling in its §3 (as the reason the "widening `WaitDelay` is a one-line follow-up" lever is unavailable) but never ledgers it as its own defect. B's framing is better: the two populations want OPPOSITE moves, and no S-row pins them apart. Verified: `D-naming` records the reuse as deliberate ("a future tuning moves both — accepted"), and nothing in S1-S10 fails if a future PR moves the constant for one reason. B found a defect A only glanced at. |
| C-6 (`0026:C1` `refusal:` / `0026:A8` — `errors.As` miswiring, thin S1 assertion) | D-11 (reduced) | PARTLY REFUTED, retained in reduced form. B's SECOND half — "S1 asserts a rendered CLI string with no named observable, no stable substring, no golden-file convention" — is REFUTED against the record: S1 explicitly names the observable ("The rendered half asserts on the envelope's `detail` field carrying `detailMayHaveApplied` (`flow_exec.go:423`'s existing constant, per docs/cli-output-contract.md's error envelope) — … so the two legs are asserted the same way and neither needs a new golden"). B is citing the 3amigo QA M2 finding as if it were still open; the record shows it CLOSED, with a named field, a named constant and a named source line. B's premortem week-four narrative and its AT-3 both rest on this refuted half. The FIRST half survives: the `errors.As` match does land on the shared `err != nil` arm at `executor.go:367-371` that pre-start failures also reach, and no scenario in S1-S10 supplies the negative control B's AT-3 asks for. Retained as D-11 at `medium`, on the negative-control gap alone. |

No Pass B row is fully excluded, but B's S1-observable charge is named here as the one
substantive artifact in either pass.

---

## 4. CONTRADICTIONS — adjudicated against record and source

Four. This is the bucket the lens exists to produce.

### CONTRA-1 — Does the per-read re-arm have an implementation site?

- **Pass A (C-2, §1.1):** the mechanism has NO site. `readBounded` is `io.ReadAll(io.LimitReader(r, limit))`; there is no read loop to re-arm inside, so C1 as written forces a rewrite of the function S3 pins as unchanged.
- **Pass B (C-1, §1a):** treats the same `precedence:` sentence as an existing, implementable mechanism whose only problem is that A9 has not MEASURED it under load. B's AT-1 says "When the drain join fires its bound and re-arms the per-read deadline" — presupposing a per-read site exists.

These cannot both be true: either the re-arm has a site and only needs measuring, or it has no site and the RDR is specifying a rewrite it does not acknowledge.

**PASS A IS RIGHT.** Grounded at `/Users/cwensel/sandbox/newcoinc/intrastate/internal/cli/cmdbind/cmdbind.go:352-364`:

```go
func readBounded(r io.Reader, limit int) []byte {
	b, err := io.ReadAll(io.LimitReader(r, int64(limit)))
	if err != nil {
		return b
	}
	_, _ = io.Copy(io.Discard, r)
	return b
}
```

`io.ReadAll` owns the loop. There is no per-read site the RDR can reach into. `0026:C1`
`precedence:` requires "The deadline is re-armed at `now + DrainGrace` BEFORE EACH read of
the final drain," and `0026:S3` requires "byte-for-byte unchanged from today … the existing
`cmdbind` suite passing unmodified — no new golden is minted for this row." Both cannot
hold: implementing the re-arm requires replacing `io.ReadAll(io.LimitReader(...))` with a
hand-rolled `for { SetReadDeadline; Read }` loop, which changes the exactly-at-cap and
cap+1 detection path (`cmdbind.go:326`, `len(inv.stdout) > StdoutCap`, fed by the `StdoutCap+1`
bound at `:275`) and the `io.Copy(io.Discard, r)` third exit, which must now also carry a
deadline. B's finding is a strict subset of A's and understates the severity: measuring a
mechanism that has no site cannot succeed. Merged as D-1 (A's charge) with D-3 (B's and A's
shared A9 charge) kept separate, because they fail at different stages.

### CONTRA-2 — How many bounded invocations does a command write run?

- **Pass A (C-1, §1.3, premortem week 2):** THREE — `executor.go:320` (pre-write baseline read), `:355` (`binding.Apply`), `:397` (read-back) — and the baseline leg's refusal is discarded at `:322-326`.
- **Pass B:** silent, but its C-6 and §1c both describe the write path in terms of `executor.go::Write`'s `err != nil` arm and "the read-back leg," reproducing the record's own two-leg framing without challenge. B's ledger accepts the two-leg model implicitly.

**PASS A IS RIGHT.** Grounded at `/Users/cwensel/sandbox/newcoinc/intrastate/internal/accessor/executor.go`:

- `:320` — `raw := e.invokeRead(ctx, reader, art, rt, protected)`, guarded by `if len(protected) != 0`, a full bounded invocation through the SAME reader (`reader, hasReader := e.Registry.readerFor(def.Accessor.Role)`) whose command spawns the same escaping helper.
- `:326` — `baselineUnread = slices.Clone(protected)` inside `if raw.class != ""`. The refusal is discarded; the journey continues. So the first held-pipe refusal in a write journey costs a full `timeout + 2·WaitDelay` and produces no signal at all.
- `:355` — `err := binding.Apply(applyCtx, art, slices.Clone(planned))`.
- `:397` — `raw := e.invokeRead(ctx, reader, art, readTimeout, compared)`.

`0026:C1` `bound:` states "a command write is two invocations (write, then read-back), each
under its own bound." That is false against source whenever `protectedKeys` returns
non-empty. Pass A's line numbers are accurate to within the two lines the `slices.Clone`
sits on (`:326`, A said `:322-326`) and its read-back cite is `:394` against an actual
`:397` — cosmetic drift, the claim itself confirmed. A found a factual error in a Normative
clause that B did not look for. Merged as D-2, the highest-severity row.

### CONTRA-3 — Is S1's rendering assertion pinned to an observable?

- **Pass A:** does not raise it; A's C-10 treats Phase 4's risk as descope, not as test thinness.
- **Pass B (C-6, §1c, AT-3):** S1's rendered half "was flagged by the record's own QA persona (M2) as asserting a rendered CLI string with 'no named observable: no exit code, no output stream, no stable substring, no golden-file convention'" and "the 3amigo consolidation records this as only *partly* resolved, not closed."

**PASS B IS REFUTED.** Grounded against `0026:S1`, which reads:

> "The rendered half asserts on the envelope's `detail` field carrying `detailMayHaveApplied`
> (`flow_exec.go:423`'s existing constant, per docs/cli-output-contract.md's error envelope)
> — the same observable the read-back leg already renders, so the two legs are asserted the
> same way and neither needs a new golden."

A named field (`detail`), a named constant (`detailMayHaveApplied`), a named source line,
and a named contract document. The constant is confirmed to exist at
`/Users/cwensel/sandbox/newcoinc/intrastate/internal/cli/flow_exec.go:469` and is already
set at `:404` and `:423`. B is quoting a QA finding the record CLOSED and reporting it as
open. B's premortem week-four incident ("S1's test is later found to assert only `Applied()`
at the Go level in CI — its 'renders the text' half was never pinned to an observable") is
built on the refuted half and does not stand. B's other half — the missing negative control
on a pre-start failure — is unrefuted and survives as D-11.

### CONTRA-4 — Is the write-path applied-sense gain one site or a broad match?

- **Pass A (C-10):** the gain is a single re-key in the LAST-but-one phase, most likely to be cut; the risk is that it never lands.
- **Pass B (C-6):** the gain is a single `errors.As` match that is "a one-line diff that silently flips a real non-applied write to 'may have applied,'" landing applied-sense "for exactly one of several failure shapes on that arm," and the record "does not show how the executor tells 'definitely did not run' from 'held pipe, might have run.'"

Not a direct contradiction of fact, but opposite risk directions from the same site — A says
it will not land, B says it will land too broadly — and B's "does not show how" is a claim
about the record that can be checked.

**BOTH SURVIVE; B'S "DOES NOT SHOW HOW" IS PARTLY REFUTED.** `0026:C1` `refusal:` DOES show
how: "`Err` wraps a typed held-pipe error carrying the pipes held and the exit status, which
the executor matches with `errors.As` at the refusal site," with the type "declared in
`internal/accessor`, NOT in `cmdbind`" for a stated import-cycle reason. The discrimination
mechanism is specified — a distinct type. What is NOT specified, and what B's AT-3 correctly
asks for, is a negative control proving the match does not over-broaden. Grounded: the match
lands on `executor.go:367-371`'s `if err != nil` arm, which is shared by `resolveArgv0`
failures and `os.Pipe` creation failures, and no scenario in S1-S10 exercises a pre-start
failure through it. A's descope risk (D-10) and B's over-broadening risk (D-11) are distinct
failures of the same site with different fixes, and both are merged as separate rows.

---

## 5. Ground-check verdicts on the sharpest claims

### Pass A: C1's `bound:` clause understates the number of write legs — **CONFIRMED**

Anchor: `0026:C1` `bound:` ("a command write is two invocations (write, then read-back), each
under its own bound").

Against `/Users/cwensel/sandbox/newcoinc/intrastate/internal/accessor/executor.go`:
- baseline read at `:320` (`raw := e.invokeRead(ctx, reader, art, rt, protected)`) — A said ~`:320`. Exact.
- apply at `:355` (`err := binding.Apply(applyCtx, art, slices.Clone(planned))`) — A said `:355`. Exact.
- read-back at `:397` (`raw := e.invokeRead(ctx, reader, art, readTimeout, compared)`) — A said ~`:394`. Off by three; claim confirmed.
- baseline refusal discarded at `:326` (`baselineUnread = slices.Clone(protected)` inside `if raw.class != ""`) — A said ~`:322`, citing the `if` at `:322` and the assignment through `:326`. Confirmed.

All three call sites route through the same `reader` obtained at `:302`
(`e.Registry.readerFor(def.Accessor.Role)`), so all three hit the same command binding and
the same `spawn` join per A1. The `bound:` clause's leg count is wrong by one whenever
`protectedKeys` returns non-empty, and the record's `disposition` mini-check enumerates only
the direct and read-back rows. **CONFIRMED.**

### Pass A: the per-read deadline re-arm has no named site because `readBounded` is `io.ReadAll` — **CONFIRMED**

Anchor: `0026:C1` `precedence:` ("The deadline is re-armed at `now + DrainGrace` BEFORE EACH
read of the final drain").

Against `/Users/cwensel/sandbox/newcoinc/intrastate/internal/cli/cmdbind/cmdbind.go:352-364`:
the entire body of `readBounded` is `io.ReadAll(io.LimitReader(r, int64(limit)))`, an early
return on error, and `io.Copy(io.Discard, r)`. No `for`, no `Read`, no `SetReadDeadline`
anywhere in the function — and `grep` over the package confirms `SetReadDeadline` appears
nowhere in `internal/` at all (which `0026:A10`'s Evidence and `0026:F5` both independently
state: "nothing in `internal/` calls `SetReadDeadline`/`SetDeadline` … today (verified
repo-wide)"). The record acknowledges only a signature change to `readBounded` (it "must gain
a reported terminal condition") and never states the function is rewritten, while `0026:S3`
simultaneously requires the existing suite to pass unmodified. **CONFIRMED.**

### Pass B: C1's normative language rests on A9, which is Pending and unmeasured under load — **CONFIRMED**

Anchor: `0026:A9`, `0026:C1` `precedence:`, `0026:MVV` row 6.

- `0026:A9` Status is **Pending**, Method **Spike**. Its Evidence says in the record's own words: "A2 measured 24 bytes and A3 measured a delayed drain under an UNBOUNDED grace, so no spike has run a large tail against a re-armed 50 ms deadline." The spike is described in the future tense ("Extend the `a3-a6-drain-bound` spike").
- `0026:C1` `precedence:` states the re-arm as settled normative text, and names the alternative it rejects as "the failure MVV row 6 exists to catch, and one no spike measured."
- `0026:MVV` row 6 is labelled "Late-but-unheld control (**Normative**)" and calls itself "the MVV's guard on C1 `precedence:`, not a formality."
- `0026:G-assumptions` — the gate that forbids exactly this ("no assumption marked `Pending` or `Unverified` may have settled-fact prose elsewhere in the RDR depending on it") — is an unfilled bracketed template block, never answered.

Every element of B's claim checks out, and A's C-5 states the same thing with the gate cite
attached. The "unmeasured under load" half is additionally supported by A9's own "If wrong"
("a large answer is refused as a held pipe on a loaded machine"). **CONFIRMED.**

### Additional ground-check: Pass A's C-11 ("never parsed" is unenforced) — **CONFIRMED**

Anchor: `0026:C1` `refusal:` ("The returned invocation carries no stdout; stdout read on a
held drain is never parsed on ANY path").

`/Users/cwensel/sandbox/newcoinc/intrastate/internal/cli/cmdbind/cmdbind.go:304` builds
`inv := invocation{stdout: stdout, stderr: tail(stderr)}` BEFORE the `werr` switch, and both
error returns — `:322` (`return inv, wrap(inv.stderr, werr)`) and `:326-328` (`return inv,
wrap(inv.stderr, errors.New("the child's stdout exceeded …"))`) — return that populated
`inv` alongside the error. The field is readable on every error path. `0026:A1`'s Evidence
confirms the guarantee is textual, not structural: it enumerates the three callers at
`:566`, `:657`, `:722` and observes each is "followed immediately and unconditionally by
`if err != nil { return … }`." That is a convention verified by inspection of three current
call sites, exactly as A charges. **CONFIRMED.**

### Additional ground-check: Pass A's C-10 (`Applied()` has zero non-test callers) — **CONFIRMED**

`grep -rn "\.Applied()" internal/` returns only `internal/accessor/adversarial_0004_test.go`
(lines 294, 444, 632, 716, 816). The accessor is declared at
`/Users/cwensel/sandbox/newcoinc/intrastate/internal/accessor/model.go:323` (A cited `:324`;
off by one). `/Users/cwensel/sandbox/newcoinc/intrastate/internal/cli/flow_exec.go:410-412`
is the `ClassExecutionFailure` arm and carries no phase check and sets no `Detail` — it
returns `envErr(codeAccessorFailed, id, "the accessor `"+id+"` could not be executed")` and
nothing else. `0026:A8`, which owns this gain, is Status **Pending**. **CONFIRMED.**

### Additional ground-check: Pass B's S1 "no named observable" — **REFUTED**

See CONTRA-3. `0026:S1` names the `detail` envelope field, the `detailMayHaveApplied`
constant, `flow_exec.go:423`, and `docs/cli-output-contract.md`. The constant exists at
`/Users/cwensel/sandbox/newcoinc/intrastate/internal/cli/flow_exec.go:469`. B reports a
closed QA finding as open. **REFUTED.**
