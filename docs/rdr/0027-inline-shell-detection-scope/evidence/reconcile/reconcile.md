Model: claude-opus-5

# Stage 6 — Reconcile Spikes & Assumptions — RDR 0027

## Preflight (Stage 5 completeness)

| Check | Call | Result |
| --- | --- | --- |
| Row complete | `--outcome lens` | `emit.next = /rdr-reconcile`; `emit.row = none` (rule `lens-large-row-complete`) |
| Critique finished | `--outcome critique` | `none` — second pass taken and diffed (rule `critique-large-diffed`); `critique_models=differ`, no single-model fallback |
| Repeatability owed? | `--outcome repeatability` → chain | `resolve:determinacy` → `none` (rule `determinacy-na`); the `Determinacy:` line reads `n/a`, judged and written at Stage 5 |

No caveat to carry. Profile `large`; accretion floor `none` (rule `floor-below-two`), so the recorded Profile is at or above the floor.

## Open set (four sources)

1. **Pre-Lock needs-verification list** — none pasted at invocation; reconstructed from the lens evidence. The `critique/iter-2/delta.md` closing note flags A7 as "the same class of gap as A5's durability caveat... noting for Stage 6", and `§mini-check-trace`'s closing line names two OPEN rows that "resolve at Stage 6". Both reduce to A5/A7, already in source 2.
2. **Pending/Unverified assumptions** — A5, A6, A7 (`ca_pending_ids` at entry).
3. **Unrun spikes** — `spikes_unrun=[]` in the tracked form, but all three Pending assumptions carried a `Method: Spike` with the spike named-but-unrun in prose. Those three are the real source-3 items; run in this pass.
4. **Exactness-word delta** — `rdr lint 0027 | grep prose:exactness` returns nothing. Read by eye, the round-introduced quantifiers ("every wrapper", "the only forms", "nothing before argv[i]") each already carry an Evidence Record from A3/A4 at Stage 4.

**Absorption audit** — lens dirs present: `grounding`, `3amigo` (+`iter-2`), `critique` (+`iter-2`), `spikes`, `research`. `3amigo/iter-2/findings.md`: "None. The delta is clean" (PASS, zero findings). `critique/iter-2/delta.md`: all four delta-scope checks CLEAR, "No new defects." No unabsorbed residue survives either round; the only carried-forward items are the three assumptions above, which the rounds explicitly deferred here by name.

## Dispositions

| Item | Source | Disposition | Evidence pointer / plan |
| --- | --- | --- | --- |
| A5 — stdin-fed-interpreter hazard owned by the charted `stdin` successor | 1, 2, 3 | **DOWNGRADED** | `evidence/spikes/a5-stdin-write-path-withholding-point.md`. Withholding premise falsified as present-tense fact: `cmd.Stdin` assigned unconditionally at `internal/cli/cmdbind/cmdbind.go:207`, sole assignment, no branch/config/field read before it; refusal ladder `:163-202` never touches stdin; all three bindings send the same envelope (`:567`, `:658`, `:723`); no `Stdin` field on `table.Accessor` (`internal/table/model.go:145-174`). **Not a refutation of C1** — C1 names the stdin forms admitted by the predicate's own definition and calls the routing "an ownership claim, not a closure" (§approach), so the falsification bounds a durability claim only. Survivable: adding the withholding point is a one-branch, one-field local change at that call site. Named plan: this record ships NO code for the axis; S4 asserts the admitted forms stay green, which is the whole obligation. |
| A6 — a category description is surfaceable without provoking a refusal | 1, 2, 3 | **VERIFIED** | `evidence/spikes/a6-description-surface.md`. Landing site exists by shipped precedent: `internal/cli/lint.go::lintExtendedDesc` (registered `lint.go:58` via `withExtendedHelp`) already renders a closed taxonomy's descriptive text live from `graphlint.BlockingCodes()`/`AdvisoryCodes()` with zero defect present; `internal/cli/help_all.go` keeps that stream off the respond gateway by construction; `internal/cli/docs.go::runDocs` mirrors it into the committed `docs/cli-reference.md`. `docs/cli-output-contract.md` names `intrastate lint --help-all` as the live finding taxonomy's home and governs only the refusal envelope's wire payload, which this surface does not touch. |
| A7 — the new false-refusal class is rare enough to accept as a stated cost | 1, 2, 3 | **VERIFIED** (scope-limited) | `evidence/spikes/a7-false-refusal-base-rate.md`. 264 unique, independently-authored shell-free argv bindings (Docker `CMD`/`ENTRYPOINT` exec-form, compose `command:`) from ~20 unrelated OSS projects, run against C1's predicate implemented verbatim; interpreter map checked against `internal/table/load.go::shellInterpreters`. 23 refusals, **23/23 true positives**, every one firing at `i=0, j=1`. Only 2/264 (0.8%) carry a listed basename at any non-argv0 position, neither a false refusal. **0/264 exhibit the A7 shape.** Sampling caveat recorded in the RDR: container entrypoints are a convention-heavy genre, so this bounds the rate rather than certifying a population rate; A7's If-wrong stays the re-open condition. |

All three dispositions are written into the RDR's Critical Assumptions section, not only here. Post-edit: `ca_verified=6`, `ca_pending=1` (`0027:A5`), `ca_off_vocabulary=0`.

## Amendment sweep (§amendment-sweep)

A6's flip to Verified made C1's `promise:` clause — written CONDITIONAL on A6 being Pending — stale. Every dependent site was re-read and updated in the same pass:

| Site | Change |
| --- | --- |
| `0027:C1` `promise:` | Conditional rewritten as settled: names the `--help-all` surface, keeps "predicate line alone holds until Phase 3 lands it" as the honest present-tense fallback (S7 is that assertion) |
| `Determinacy:` line | "what remains open (A5, A6, A7)" → records all three settled at Stage 6 |
| §capability-dependencies | Added the description-surface row as **Available**; A5 row now states the falsification and that the record depends on none of it |
| §phase-3-promise-wording | "Where it lands is the implementer's call (A6)" → the settled surface + the pattern to follow |
| §testing-strategy framing | "7 depends on A6 resolving" → "7 reads the description from the `--help-all` surface A6 settled" |
| `0027:S7` | "read from wherever Phase 3 lands it" → reads the `--help-all` body with no model loaded and no defect present |
| `0027:D-selection-predicate` | "open under A7... judged against whichever way A7 settles" → judged against A7's settled reading |
| §consequences | Position-freedom bullet: cost ACCEPTED and named, "no wrapper table" economy NOT re-priced; residual risk is the sampling caveat, not an open design question |
| §key-discoveries | `i < j` bound now carries the measured 0-of-264 result |
| §mini-check-trace 3c + closing | Both OPEN rows CLOSED (the closing line had promised "Both resolve at Stage 6"); neither moves to CONTRADICTION |

Subject-token sweep on the amended clause: `--help-all`, `promise:`, and `A6` greps return only the sites above. No clause added or removed, so no sibling-agreement re-read owed beyond C1's own block, which was re-read whole.

## Hard rules

- **Refutation → BLOCKER.** None fired. A5's falsification is the near miss and was tested against the rule directly: C1 must RELY on the refuted claim for a BLOCKER. It does not — C1 names the stdin forms as admitted by the predicate's own definition, and §approach already called the successor routing "an ownership claim, not a closure." The claim that died was a durability claim about a future record, not a predicate C1 asserts. No route-back, so no §punt-ledger row is owed and no §strong-consult was triggered.
- **No MVV-critical item deferred past lock.** The MVV's five steps rest on the predicate (A1–A4, all Verified) and on S7's description assertion (A6, now Verified). A5 is the only carried Pending and the MVV consumes nothing from it — S4 asserts the stdin forms stay GREEN, which is what the predicate delivers with or without the successor.

## Completeness check

- `_Draft placeholder._` — 0 occurrences. `this is a seed skeleton` — 0 occurrences. No hollow body section.
- `## References` — fully authored (peer RDR elements, source symbols, prior-art paths, trackers, Stage 2 evidence); no template brackets.
- `rdr lint 0027` — `PASS blocking=0 resolution=0`. The 5 surviving `placeholder:survived` findings and the `gate:inline` finding are all inside §finalization-gate (lines 472-592): unanswered Stage 7 gate prompts (`gate_written=false`), which Stage 7 owns and Stage 6 does not touch.

## Verdict

**RECONCILED** — all four sources exhausted, every item terminal (2 VERIFIED, 1 DOWNGRADED), no BLOCKER. Ready for Finalize.
