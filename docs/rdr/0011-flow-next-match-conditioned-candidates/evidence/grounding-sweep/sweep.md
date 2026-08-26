Model: claude-fable-5

| # | anchor | verdict | found |
|---|--------|---------|-------|
| 1 | `internal/cli/flow_next.go::excluded` | CONFIRMED | L276-294: `probe.Match = nil`, `probe.Escape = nil`, `Recognized: row.Outcome`, `return result.Refusal.Kind == resolve.KindNoMatch`; doc comment carries both quoted phrases (line-wrapped). |
| 2 | `internal/cli/flow_next.go::summarize` | CONFIRMED | L168-207: `for _, atom := range row.Atoms` with no Block filter, `view[atom.Key]` presence test, absent keys appended to `c.Unresolved`; `row.RequiresOwned` absent keys and, when `!gatesRan`, `row.Gate` ids appended likewise. |
| 3 | `internal/cli/flow_next.go::assembledView` | CONFIRMED | L232: `func assembledView(owned []resolve.Tag, observed []resolveTag) map[string]string`; owned written first, then observed. |
| 4 | `internal/cli/flow_next.go::runFlowNext` | CONFIRMED | L115 `if len(row.Escape) != 0 { continue }`; L157 `payload.Outcomes = slices.Clone(req.model.Outcomes)`; L149-156 comment rejects filtering the alphabet by surviving candidate rows, "Recorded as DEV-4"; L97 `invokedReaders(req.model, "")`. |
| 5 | `internal/table/model.go::Row.KernelRow` | CONFIRMED | L232-246: `a.Block == BlockMatch` -> `match`, else -> `guard`. |
| 6 | `internal/table/model.go::Row` | CONFIRMED | L170-171 `// Atoms is the unified predicate set ...`; type doc L151-153: atoms travel as ONE unified set with authored block retained, Match/guard split applied at kernel handoff. |
| 7 | `internal/resolve/resolve.go::TagSet.matches` | CONFIRMED | L174-177: `if !ok \|\| tv.conflicted \|\| tv.value != w.Value { return false }`. |
| 8 | `internal/resolve/resolve.go::Resolve` | CONFIRMED | Doc step 2 L418-419 "candidate rows are the non-escape rows for that outcome whose match pattern holds"; L472 `if !view.matches(row.Match)`; L483-493 `switch len(selected)` case 1 plan / case 0 KindNoMatch / default KindAmbiguousMatch. |
| 9 | `internal/resolve/resolve.go::gate` | CONFIRMED | L518-535: `if verdict == GuardFalse { continue }` prunes; `missingOwned` -> `KindOwnedStateUnavailable` returned before the undecidable/`KindGuardUnevaluable` branch (L561). |
| 10 | `internal/resolve/resolve.go::KindNoMatch` et al. | CONFIRMED | L52 "no_match", L54 "ambiguous_match", L57 "owned_state_unavailable", L60 "guard_unevaluable". |
| 11 | `internal/cli/flow.go::registerSelectionFlags` | CONFIRMED | L92 `func registerSelectionFlags(cmd *cobra.Command)`. |
| 12 | `internal/cli/flow_exec.go::invokedReaders` | CONFIRMED | L86 `func invokedReaders(m *table.Model, outcome string) []string`. |
| 13 | `internal/cli/flow_next_0005_test.go` / fixtures | CONFIRMED | `TestReq40And74_NextPayloadCarriesTheAlphabetAndCandidateSummaries` (L30, seeds `status=draft` L32) and `TestReq42And46_OnlyReportedCandidatesGatesRunAndResultsRideTheCandidate` (L197, seeds `status=draft` L199); fixtures `flowMVVModel` (L77) and `flowGatedNextModel` (L245) with rules `gated-reported` (L305) and `gated-excluded` (L317), each `[rule.match.status] eq = "draft"`. |
| 14 | `models/rdr.toml` | REFUTED | `grep -c '^\[\[rule\]\]'` = **22** `[[rule]]` tables (lines 109..350), not 21. `[rule.match.stage]` and `[rule.match.recognized]` each appear 22 times, so every rule carries both (that half holds). |
| 15 | `0005:C1` | CONFIRMED | All four quoted phrases present verbatim in the projected C1 text. |
| 16 | `0005:§decision-rationale` | CONFIRMED | Premortem paragraph: "`next` omits enough condition detail to constrain a skill". |
| 17 | `0005:§briefly-rejected` | CONFIRMED | "cannot carry legal outcome alphabets, conditional" / "summaries, or typed tag sets" (line-wrapped). |
| 18 | Peer elements exist | CONFIRMED | `0005:A3`, `0005:A7`, `0005:D-identity`, `0005:D-selection-predicate`, `0007:C8`, `0010:C4`, `0010:A11` all project non-empty. |
| 19 | `0010:C4` | CONFIRMED | Ends "`flow next`, `flow read-state`, and `flow set-state` are unchanged by this RDR." |
| 20 | `0010:A11` / `0010:§minimum-viable-validation` | CONFIRMED | A11 text: "re-verified against cli/0011's candidate predicate once 0011 locks" (Status Pending); MVV step 5 contains "`flow next --model m.toml --as=json` with no `--artifact` → exit 0, every candidate with empty `required` (A11)". |
| 21 | qmuntal-stateless README / statemachine.go / states.go | CONFIRMED | README L134 quoted sentence verbatim; statemachine.go L175 `func (sm *StateMachine) PermittedTriggersCtx`; states.go L201-207 appends `key` only when `len(tb.UnmetGuardConditions(ctx, unmet[:0], args...)) == 0`. |
| 22 | transitions core.py | CONFIRMED | L952-961 `get_triggers` docstring "Collects all triggers FROM certain states", body is the quoted list comprehension with no condition evaluation; L897-911 `_can_trigger` iterates `self.get_triggers(state)` and evaluates transitions. |
| 23 | xstate CHANGELOG / State.ts | CONFIRMED | CHANGELOG L1525, L3076 "Removed `MachineSnapshot['nextEvents']`"; State.ts L305 `const machineSnapshotCan = function can(` calling `this.machine.getTransitionData(this, event)` at L315. |
| 24 | scxmlcc user-manual.md | CONFIRMED | L134-135 under `##### cond`: "The transition is only executed if the condition evaluates to true". |
