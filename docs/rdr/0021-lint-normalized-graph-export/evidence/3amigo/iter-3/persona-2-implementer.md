Model: claude-sonnet-5
Iteration: 3 (pass 2 — delta-scoped)

# Persona 2 — Implementer (delta-scoped second pass)

Question asked throughout: if I started coding this Monday morning, what would I have to ask in the first hour — with focus on gaps the pass-1 rewrite may have opened in `C1`–`C5` and `D-*`.

## Scope note

Stayed within `C1`–`C5`, `D-identity`, `D-wire-byte-format`, `D-naming`, `D-selection-predicate`, plus the `source-anchor` edges and the Pre-Lock Mini-Checks tables (`0021:§pre-lock-mini-checks`) that those clauses cite as authority/oracle. Widened briefly into `internal/graphlint/analysis.go` and `taxonomy.go` beyond the `reach.go` anchor already named for `C4`/`A8`, because `C4`'s own text makes a claim about a second function (`checkNodeCeiling`) that isn't named in `A8`'s anchors — the claim sent me there.

## Findings

### 1. `0021:C4` — "distinct guard-product bound" misidentifies its own comparison target

`C4`'s last sentence: "The ceiling this refusal names is `reach`'s node-count completeness, NOT `analysis.go`'s `checkNodeCeiling`, a distinct guard-product bound."

Grounding against `internal/graphlint/analysis.go` and `internal/graphlint/taxonomy.go`:

- `analysis.go::newAnalysis` (line 46) calls `nodes, complete := reach(m)` — the exact same private `reach()` the export's completeness surface (`A8`) reads.
- `analysis.go::checkNodeCeiling` (line 177) fires on `a.complete` being false — i.e., on the SAME `reach()` completeness the export needs, using the SAME `nodeCeiling` constant (`taxonomy.go:136`, value 4096, exposed as `graphlint.NodeCeiling()` at `taxonomy.go:152`).
- The codebase does have a second, genuinely distinct bound — `guard.Bound()`, exposed as `graphlint.ProductBound()` (`taxonomy.go:147`) — but `checkNodeCeiling` never calls it. `ProductBound`/`guard.Bound()` is the guard-product bound; `checkNodeCeiling` is not.

So `checkNodeCeiling` IS the node-count ceiling check (same mechanism, same constant as the export needs) — it is not "a distinct guard-product bound." The sentence has swapped which function carries which bound. This is exactly backwards from what an implementer needs to know: the correct fact — the one that matters for wiring `graph-export-too-large` correctly — is that lint's `checkNodeCeiling` and the export's new refusal are THE SAME bound observed through two call sites (which is consistent with, and actually strengthens, `C4`'s own neutrality claim: same traversal, same ceiling, two consumers). The RDR's attempt to pre-empt a reviewer question ("is this a new bound?") instead points the implementer at a nonexistent distinction and a wrong function pairing.

Practical effect Monday morning: an implementer checking "do I need to also thread `guard.Bound()`/`ProductBound()` into the export's refusal path" will read this sentence, go looking for how `checkNodeCeiling` relates to the guard-product bound, find that it doesn't, and lose time reconciling the RDR's claim against the tree before concluding (correctly) that the export refusal needs only the node ceiling, not the product bound. The blocking decision is which bound(s) the new `graph-export-too-large` refusal must check and whether it needs any guard-product input at all — the sentence as written answers this wrong.

Severity: medium. It doesn't block writing correct code (the surrounding sentence and `A8` still make node-count completeness the right target), but it is a factually incorrect grounding statement inside a normative clause that will cost real implementer time and could mislead a reviewer verifying `C4` against the tree into thinking there's a second bound to reconcile.

Blocks: the wording of `C4`'s refusal-authority sentence (whether the new completeness surface needs to be reconciled with, or is independent of, `guard.Bound()`/`ProductBound()`).

## Checked and clean (no finding)

- `C2`'s row field list (`identity, source, kind, outcome, atoms, next, writes, requires_owned, gate, escape, emit`) matches `table.dumpColumns` (`internal/table/dump.go:17`) verbatim, same order — no over- or under-specification introduced by the rewrite.
- `C2`'s `reach.nodes[].values` shape matches `reach.go::Node.Values map[string][]string` (already dispositioned in pass 1 as fixed; re-confirmed clean, not re-reported as a finding).
- `C4`'s neutrality mechanism-independence claim holds: `analysis.go::newAnalysis` and the export path both would reach the same private `reach()`; nothing in the rewrite routes the export through a second traversal.
- `C1`'s arm codes (`flag-mutually-exclusive`, `flag-invalid-value`, `flag-required`, `model-unreadable`, `model-invalid`) match `internal/cli/lint.go::runLint` verbatim, and the narrowed "mirrors verbatim" scope (arm set + code spellings, not help text) is consistent with `lint.go`'s own per-verb `Long` help text pattern — the `--emit` arm's own `flag-invalid-value naming emit` is fully specified as a sixth, verb-owned arm alongside the five mirrored ones; nothing is left implicit.
- `C5`'s mode-coexistence text agrees with `respond.go`: `--as` defaults to `ModeText` (`respond.go:99-109`), matching C5's "under `--as=text` (the default)"; `respond.FlagName = "as"` (`respond.go:46`) confirms `--as` is the persistent root flag C5 and D-selection-predicate both reference.
- The Determinacy line added to Normative Contracts is a gate-status citation of `C2`+`C3`, not new normative content — it introduces no requirement `C2`/`C3` don't already state, so it does not function as a sixth de facto clause.
- `guard.AssignmentCount` (`internal/guard/declaration.go:80`) exists with the signature `C2` and `A8`'s neighbor assumptions rely on — confirmed real, not re-reported (already verified clean in pass 1).
