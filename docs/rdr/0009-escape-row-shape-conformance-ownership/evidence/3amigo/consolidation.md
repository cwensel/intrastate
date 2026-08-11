Model: claude-opus-5[1m]
Lens: 3amigo — consolidation (iteration 1)

Cross-persona overlap over `persona-1-pm.md`, `persona-2-implementer.md`,
`persona-3-qa.md`. Persona IDs are per-file; the `T-N` rows below are the merged
origin ledger the resolve pass dispositions against.

## Hotspot passages (2+ personas)

| # | Passage (grepable anchor) | Personas | Type |
| --- | --- | --- | --- |
| T-1 | Normative Contracts: "the error MUST remain diagnostic in that case by stating the breach count alongside the identities" | IMP-2, QA-4 | unspecified carrier + internal contradiction with A8 |
| T-2 | Normative Contracts: "Both the error type and the sentinel MUST be EXPORTED" / Load-Bearing Decisions: "the exact type/sentinel spelling is sharpened at Pre-Lock" | IMP-1, QA-3, QA-5 | deferred-to-a-stage-that-has-run; blocks every assertion |
| T-3 | Testing Strategy sc.11 + Prerequisites: "(binds the verb that first calls `Resolve`)" / A2 "Obligation, not a gap" | PM-4, IMP-4, QA-5 | scope silence: unexecutable scenario listed flush with executable ones; no `Code`/`Hint`/wire-key literal |
| T-4 | Normative Contracts typed-error block vs Failure Modes vs sc.9: `errors.Is`/`As` vs "must traverse the aggregate" vs `Unwrap() []error` | IMP-5, QA-3 | traversal contract indeterminate (depth, per-element type, per-element `errors.Is`) |
| T-5 | Consequences: "no write-bearing escape plan can be emitted by any producer path" vs Capability Dependencies "Deferred" | PM-1, PM-2 (+ QA-5 on the same deferral seam) | overclaim vs what Phases 1–2 ship; named user protected only by deferred halves |

## Single-persona findings carried into the ledger

| # | Passage | Origin | Type |
| --- | --- | --- | --- |
| T-6 | Testing Strategy sc.6: "Mutation check — revert the entry check and re-run. Expected: At least one test fails" | QA-1 | vacuous-by-construction; no mutant, no surviving oracle |
| T-7 | Testing Strategy sc.5: "Every assertion still passes and still discriminates" | QA-2 | no oracle for "discriminates"; baseline stale-by-design |
| T-8 | Technical Design: "called by `resolve.go::Resolve` at entry before `assemble`" + Precedence pin | IMP-3 | is the placement load-bearing; must the frozen REQ-1 doc list be amended |
| T-9 | Testing Strategy: "Done means: every producer path is closed…" | PM-3 | success stated only as codebase properties, never as the user outcome |
| T-10 | Consequences: dormant-row behavior change + Phase 3 `writes = []` strictness | PM-5 | author friction never weighed; rests on A5, which expires when RDR 0002 ships |
| T-11 | sc.3: "resolves exactly as the A3 baseline" | QA-2 (below-top-5) | empty-vs-nil never executed in any A3 run; `len==0` cannot discriminate the two states the scenario exists to separate |
| T-12 | Load-Bearing Decisions → Precedence: "precedes every modeled disposition" | QA-3 (below-top-5) | scenario'd 1-of-5 kinds against the MVV's exhaustive-over-kinds house style |
| T-13 | Technical Design: "callers MUST check the error before reading the disposition" | QA-3 (below-top-5) | the novel `Refused() == false` trap has no scenario |
| T-14 | sc.7: "Verdicts identical to `Resolve`'s entry check" | QA-3 (below-top-5) | no comparison relation named for two `errors.Join` values |

## Signal check

**Healthy.** Three distinct lists, no reused findings, no generic advice. Each
persona found what its lens is built to find: PM found outcome/ownership
overclaim, Implementer found unnamed identifiers and data shape, QA found
unassertable expected values. The convergence at T-1/T-2/T-3/T-4 is genuine
overlap on four passages, not three copies of one critique.

**No false claim about the code was found by any persona.** Both code-reading
personas independently verified the RDR's grounding at HEAD (`d82f1b5`) and
reported every citation accurate — the 16 `escapeRow` call sites, both
`escape.Writes` overrides, all-`nil` `Resolve` returns, `Cause` `json:"-"`,
the exit-code mapping, `config.Load` carrying no `Hint`. The findings are
**under-specification, not misstatement**, which is the load-bearing signal:
this RDR's evidence is solid and its *forms* are unnamed.

## Dispatcher-side grounding (run alongside the personas)

Checked the frozen package-boundary guards none of the three personas read,
because Phase 1 adds exported symbols and two imports to a package whose surface
is frozen by test:

- `resolve_test.go::TestReq25_NoExportedSymbolImpliesOrchestrationOrPersistence`
  bans exact names `Run/Execute/Apply/Persist/Save/Commit/Start/Orchestrate`;
  `TestReq36_…` bans `Encode/Decode/Marshal/Unmarshal/Import/Export/Parse/
  Serialize/Deserialize/Invert`; `TestReq37_…` bans `Hash/Checksum/
  Canonicalize/Fingerprint/Digest`. All are **exact-match** comparisons, so
  `CheckValid` / `Validate` and an `Err…` sentinel collide with none of them.
- Import guards: `TestReq26_…` (no third-party), `TestReq37_…` (no
  `encoding/json`, hashes), `TestReq28_…` + `mvv_test.go` leg 3 (no CLI
  package). `errors` and `fmt` are forbidden by none.

**Verdict: Phase 1's new surface clears every frozen boundary guard.** Recorded
so the naming decision (T-2) is made against a checked constraint set rather
than an assumed-open one. Note the corollary for T-3: the kernel may not import
`encoding/json` or any CLI package, so the wire-carrier work is necessarily
CLI-side — the kernel cannot serialize its own breach.
