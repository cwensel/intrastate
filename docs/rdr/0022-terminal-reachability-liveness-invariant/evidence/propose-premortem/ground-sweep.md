Model: claude-fable-5

# Factored grounding sweep — per-anchor verdicts

Checked read-only against the current checkout; peer-RDR anchors read via the
projector (`rdr inspect --select`), never the record files.

## Code anchors

| # | Anchor | Verdict | Evidence |
| --- | --- | --- | --- |
| 1 | `internal/graphlint/reach.go::reach` — returns `(nodes, complete bool)`; incomplete above ceiling | CONFIRMED | Signature `func reach(m *table.Model) (nodes []Node, complete bool)`; loop sets `complete = false` and breaks when `len(nodes) > nodeCeiling`. |
| 2 | `reach.go::successor` — source clone with writes applied, clears removed | CONFIRMED | `out := src.clone()`; clear sentinel deletes the key; owned writes replace via `heldValues`. |
| 3 | `reach.go::subsumes` — same key set required; value-set coverage | CONFIRMED | Rejects on `len(n.Values) != len(want.Values)` or a missing key; every wanted value must be contained in the held set. |
| 4 | `reach.go::indexOf` — first node in slice order that subsumes | CONFIRMED | Linear scan `for i, n := range nodes`, returns first `subsumes(n, want)`, else -1. |
| 5 | `reach.go::matchSatisfiable` — only owned match atoms; guards never consulted | CONFIRMED | Skips `a.Block != table.BlockMatch` and non-owned provenance; doc comment: "no guard atom is consulted at all". |
| 6 | `reach.go::ownedAtomSatisfiable` — existential over held values | CONFIRMED | "whether SOME value the node holds for the atom's key satisfies it"; loop returns true on any admitting value. |
| 7 | `reach.go::heldValues` — opaque value when no finite domain, gated on `guard.AssignmentCount` | CONFIRMED | `if _, finite := guard.AssignmentCount(m.Tags[key]); !finite { return []string{OpaqueValue} }`. |
| 8 | `reach.go::key` — canonical, injective node identity, escaped fields | CONFIRMED | Sorted keys, `escapeField`/`escapeJoin` per field; doc comment states the rendering is injective. |
| 9 | `internal/graphlint/analysis.go::splitNode` — splits on given keys only, others pass through | CONFIRMED | Only keys in `keys` with >1 held value fan out; clone carries every other key unchanged. |
| 10 | `analysis.go::terminalKeys` — sorted owned keys in terminal predicates | CONFIRMED | Collects owned-provenance atom keys from `a.model.Terminal`, dedupes, `slices.Sort`. |
| 11 | `analysis.go::satisfiesSomeTerminal` — universal over per-tag value sets | CONFIRMED | Via `nodeMeetsAll`: every held value for each atom key must pass `atomAdmitsValue`; doc says "UNIVERSAL over the node's per-tag value sets". |
| 12 | `analysis.go::hasOutgoingOrdinaryRow` — excludes escape rows; doc books accepted false-negative residual citing D12 | CONFIRMED | Skips `row.Kind() == table.KindEscape`; comment: "the accepted false-NEGATIVE the record books against invariant 2 … See D12." |
| 13 | `analysis.go::checkDeadEnd` — runs on split nodes; early return via `table.IsDecisionTable` | CONFIRMED | `if table.IsDecisionTable(a.model) { return }`; then iterates `splitNode(n, keys)`. |
| 14 | `analysis.go::checkNodeCeiling` — incomplete traversal emits blocking graph-product-too-large, element "traversal" | CONFIRMED | Emits `CodeProductTooLarge` with `Element: elementTraversal` when `!a.complete`; that code is in `blockingCodes`. |
| 15 | `analysis.go::newAnalysis` — runs reach once, stores nodes + complete | CONFIRMED | `nodes, complete := reach(m)` once; both stored on the `analysis` struct. |
| 16 | `internal/graphlint/taxonomy.go::blockingCodes` — ends with CodeProductTooLarge; contains CodeTerminalEscape; `nodeCeiling = 4096`; ReasonNoParticipatingDimension appended (RDR 0010) | CONFIRMED | All four properties present; the reason's comment cites `0010:C5` and names it "the append" to the closed set. |
| 17 | `internal/graphlint/engine.go::sortFindings` — sorts by identity key (model, code, namespace, id, fingerprint...) | CONFIRMED | `slices.SortFunc` over `identityKey`, which joins `Model, Code, namespace, id, Fingerprint`, then tie-break fields. |
| 18 | `internal/cli/lint.go` — enumerates `graphlint.BlockingCodes()` dynamically | CONFIRMED | `lintExtendedDesc` loops over `BlockingCodes()` and `AdvisoryCodes()` to build the --help-all body; no hardcoded list. |
| 19 | `internal/cli/help_all_test.go` — iterates BlockingCodes() and AdvisoryCodes() | CONFIRMED | Lines 170 and 175: `for _, code := range graphlint.BlockingCodes()` / `...AdvisoryCodes()`. |
| 20 | `models/rdr.toml` — `terminal = ["landed", "abandoned"]`; rows matching/emitting "revise" (cycles) | CONFIRMED | Line 13 declares exactly that terminal list; multiple rules match `recognized eq "revise"` (e.g. `seed-revise` writes `stage = "seeded"` back, a self-cycle). |
| 21 | `docs/cli-reference.md` and `docs/model-authoring.md` — each contains a table or list of graph-* finding codes | REFUTED (half holds) | cli-reference.md: CONFIRMED — lines ~731-748 list all ten blocking and four advisory codes. model-authoring.md: no table or list of finding codes exists; its 24 graph-* mentions (incl. graph-dead-end at lines 413, 673) are all inline in prose sentences or inside bullets about other topics. |

## Peer-RDR anchors (projector)

| # | Anchor | Verdict | Evidence |
| --- | --- | --- | --- |
| 22 | 0006:C3 — "at least" the listed blocking invariant classes; no stability clause inside C3 | CONFIRMED | "MUST check at least these blocking invariant classes: dangling edge, dead end, determinism/overlap, guard exhaustiveness/gap, single-valued state, owned-set-before-match, and declared terminal/escape handling." Nothing else in C3. |
| 23 | 0006:C13 — stable code, model identity, severity, message, source rule/context id or span when available | CONFIRMED | Verbatim match, plus atom-attribution and escape-population clauses. |
| 24 | 0006:C14 — aggregate CLIError code `graph-lint-failed`; findings field append-only, typed, owned by clierr | CONFIRMED | "one aggregate `CLIError` with code `graph-lint-failed`" and "append-only typed `findings` field owned by `clierr`". |
| 25 | 0006:C15 — order by model id, invariant code, rule/context id or element id, then fingerprint; rule ids before element ids | CONFIRMED | Exact tuple stated; "a source rule/context id sorts before any graph element id." |
| 26 | 0006:C16 — every defect decidable in one pass, never first-failure | CONFIRMED | "MUST report every defect it can decide in one pass over a row group, not the first it encounters." |
| 27 | 0006:C17 — advisory tier closed at four named codes; advisory never changes success disposition | CONFIRMED | Names the four codes; "Advisory findings MUST NOT change the success disposition." |
| 28 | 0006:C18 — no initial owned state rejected with a blocking finding | CONFIRMED | "MUST be rejected with a blocking finding; lint MUST NOT treat an absent root as an empty reachable set." |
| 29 | 0006:ALT4 — external model checker rejected; native lint over the consumed model | CONFIRMED | Rejection: "the authoritative gate should be native lint over the model intrastate actually consumes." |
| 30 | 0006:D-reachability-relation — escape rows are self-loop edges (no write block, no clear list); excluded from invariant 2's outgoing-row test | CONFIRMED | "Escape rows are edges too, and they are self-loops"; "carries neither a write block nor a clear list"; "Escape self-loops are excluded from invariant 2's 'at least one outgoing row' test." |
| 31 | 0006:D-soundness-direction-is-per-invariant-not-global — universal checks must not read merged nodes; dead end's remedy is splitting on terminal-participating keys | CONFIRMED | "Universal checks must not" read merged nodes; dead end's remedy "is to split the node on terminal-participating keys first, recovering exactness." |
| 32 | 0010:C5 — machine-only invariants (incl. dead end) vacuous over decision tables, must not emit | CONFIRMED | "dead end (2), always-present-owned (5), owned-set-before-match (6), and single-valued state are vacuous by construction over this class and MUST NOT emit" (missing-root arm of invariant 1 likewise). |
| 33 | 0017:C1 — findings[].code names its own disposition, from a per-finding vocabulary its producer declares | CONFIRMED | "code is the finding's OWN discriminator … drawn from a per-finding vocabulary its producer declares." |
| 34 | 0024:§approach — new emit checks land in the load pipeline, not graphlint; emit values stay uninterpreted | CONFIRMED | "whose checks land in the **load pipeline** (not graphlint)"; "emit values stay uninterpreted and byte-compared (`0010:C3`)"; "No graphlint code changes." |

## Totals

33 CONFIRMED, 1 REFUTED (anchor 21 — half holds: cli-reference.md yes,
model-authoring.md has no table/list of graph-* codes, only inline prose
mentions), 0 NOT-FOUND.
