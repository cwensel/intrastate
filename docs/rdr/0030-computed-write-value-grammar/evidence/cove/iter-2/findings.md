Model: claude-fable-5-1

# CoVe iter-2 — delta-scoped verification of the iter-1 fixes (RDR 0030)

Scope: 0030:C1, C2, C3, A1, A3, A4, A5, A9, A10, A11, D-identity, S1, S4,
§mini-check-tables, the audit's `valueSatisfies` row, Phase 2, Failure
Modes, Risks and Mitigations, Consequences, MVV. Source = working tree on
`main`. §finalization-gate ignored.

## Verification of each fix

| # | Fix | Verdict | Source cite |
| --- | --- | --- | --- |
| 1 | C1 subsumption rule | DEFECT (ordering claim SOUND; "retained as `in`" REFUTED — F1; inherited-context `in` unaddressed — F9) | `internal/table/normalize.go::compareAtoms` orders (Key, Block, Operator, Literal) by `strings.Compare`, so `"all"` < `"match"`; the intersection claim holds (composition ordered after the `in` would admit `{m} ∩ bounds` per member = the same set). But `internal/table/model.go::KernelRow` renders EVERY match-block atom as an equality `resolve.Tag{Key, Value: seamValue(a.Literal, …)}`, and `seamValue` on a non-set key returns `members[0]`. A retained `match in` is not an equality the kernel can test (`0002:C13`). `case expanding:` and the suffix are otherwise undisturbed; A10 is an honest Pending but inherits F1 in its per-row atom assertion. |
| 2 | C1 both-carriers clause | DEFECT (obligation SOUND; collision motive REFUTED — F5) | `internal/graphlint/engine.go::Fingerprint` hashes `row.Atoms` THEN `row.NextTags`; each expanded row carries a distinct `guard.all eq <cell>` atom, so fingerprints cannot collide even with `Writes`-only. Census: `graph_document.go:219-220` (both), `flow_next.go:397-409` (both via `KernelRow()`), `guard/lint.go:180-183` (both), `dump.go:125-127` (both), `reach.go::writesOf` (union). Omitted but harmless: `groups.go:347 checkIdempotentWrites` reads `Writes`; `model.go:390-391 KernelRow` copies both. No third carrier: `Row` has `NextTags`, `Writes`, and `RequiresOwned` (keys only, `0002:C14`). A11's "no third carrier" is honest; its Fingerprint oracle is non-discriminating. |
| 3 | C1 always-append suffix | NEEDS_DECISION (totality SOUND; peer iff contradicted — F6) | `normalize.go:791 suffixed := expanding && len(members) > 1`. Distinct cells give distinct suffixes, so `0002:C13` totality holds. But `0002:C13` (Implemented) also fences "A suffix is non-empty exactly when the rule produced more than one row", which a one-cell step breaks; C1 names the code flag, not the peer sentence. D-identity defers to C1 and does not restate it — no double statement. |
| 4 | C1 width refusal | SOUND (gloss REFUTED — F7, low) | `internal/guard/declaration.go::intWidth`: `minV > maxV` → false; `span := maxV - minV; span < 0 \|\| span == math.MaxInt` → false. `load.go::tagDecl` already refuses `min > max` ("min exceeds max"), so the RDR's two-condition restatement is the same rule on any loadable declaration. The gloss "lint reports as having no domain" is wrong: `declaration.go:127-135` returns `cardinalityCeiling, true` on `!ok`, with a comment saying "reporting no domain at all would instead misreport". |
| 5 | C1 zero-cell refusal | SOUND | `load.go::conform` `case "lt","lte","gt","gte": return conformKind(decl, members[0])` — no `conformDomain`; so `lt = 0` over `min=0` loads today and the refusal is new. `0003:C6`'s disposition (unsatisfiable atom → load error, not a silent never-matching row) is correctly cited. Category unstated in C1 — F10 (low). |
| 6 | C1 enum duplicate-domain + `[initial]` split | DEFECT (dup-domain SOUND; `[initial]` SOUND; predicate path REFUTED — F3) | `load.go::tagDecl` (848-908) checks `src.Domain` only for kind agreement — no duplicate or emptiness check. `load.go:1688-1690` `[initial]` → `valueMembers` error → `CatMalformedInitialDeclaration`. `normalize.go::loader.atom` (`valueMembers` error → `badAtom` → `CatMalformedPredicateAtom`): a predicate-literal inline table refuses under `malformed_predicate_atom`, not `malformed_tag_declaration` as S4, MVV item 4 and the disposition table assert. Write path `renderWrites` → `CatMalformedTagDeclaration` (A9 correct). |
| 7 | C2 first-failing-cell order | SOUND | C1's admitted cells are enumerated in domain order (`{min..max}`, authored `domain`); C2 names the first in that order; Failure Modes F3 repeats it verbatim. MVV 4's `retry`/`attempt`/cell 9/value 10 and `tier`/`large` are the last cells and are consistent. |
| 8 | Audit `valueSatisfies` row + Phase 2 | DEFECT (all four source claims SOUND; extraction target unreachable — F4) | `guard/product.go::valueSatisfies` returns `resolve.GuardResult`; `graphlint/reach.go::atomAdmitsValue` = same core + `return verdict != resolve.GuardFalse`, comment "the false-green direction"; `guard/assignment.go::renderSet` marshals as authored; `reach.go::renderSetLiteral` applies `canonicalValues` (`slices.Compact(slices.Sorted(...))`) and substitutes `[]` for nil. But both cores call `guard.Evaluator{}.Evaluate` — `Evaluator` is declared in `internal/guard/grammar.go:97`, package `guard`, which imports `table` (`go list`). `resolve` imports nothing intra-repo and `table` imports only `resolve`; neither can name `guard.Evaluator`. Phase 2's "needs only `resolve.GuardAtom` and a `[]string` literal" omits the evaluator. |
| 9 | S1/A5 finding-set projection | SOUND (with caveat) | `internal/cli/clierr/clierr.go::Finding` has `Code`, `Key`, `Dimension`, `Class`, `Reason`, `Rule`, `Span`, `Element`, `Fingerprint` — all real. Caveat: `graph-overlap` and `graph-redundant-row` carry no Key/Dimension, so they project to `(code, "", "", class, "")` — presence only; `graph-coverage-gap` carries `Dimension` but its cell count is in `Message` only. Not vacuous on the MVV fixture (zero overlap/gap asserted separately; owned-before-write and idempotent-write project on key+dimension), but set comparison loses multiplicity. Tuple stated inconsistently — F8. |
| 10 | A5 scoping + MVV 5b `graph-overlap` | DEFECT — F2 | `guard/product.go::acceptedIn` subtracts the row's `unless` block (`accepted.subtract(excluded)`), and C1 retains the authored `unless` on every row; the over-admitted row's accepted set is EMPTY. `groups.go::emitOverlaps` `continue`s on `left.Intersect(other).Len() == 0`, so `graph-overlap` cannot fire; `checkRedundantRows` `properSubset(∅, sibling)` → `graph-redundant-row` (advisory) fires instead. |
| 11 | Five mini-check tables | DEFECT (internal — F8; oracle/disposition rows rest on F2) | `trace` has no CONTRADICTION row per se: steps 1-7 each cite the clause in force and the witness agrees. Step 8 and `oracle` row 1-2 state the projection as `(code, key, dimension)` multiset where S1/A5/MVV say `(code, key, dimension, class, reason)` as sets. `oracle` row 1-2's negative control and `disposition`'s `unless` row assert `graph-overlap` (F2). `authority` table's first row assumes the shim is reachable from the loader (F4). |
| 12 | Residue of the OLD C1 | SOUND (residue sweep clean) with new gaps F9, F10 | Scratch sweep of the whole record for "unenumerable", "compose(s) with", "already refuses", "`Writes` only": no residue outside the sentences that name the old rule to reject it. New terms: "subsumes" defined in C1 before use in A1/A10. Two new silences: inherited-context `match in` on the stepped tag (F9); category of the zero-cell and `#`-member refusals (F10). |

## Findings

### F1 — Subsumed `match in` "retained as an authored atom" is unmatchable at the kernel
- anchor: 0030:C1 (also 0030:A10 evidence)
- class: contradiction
- claim: "The subsumed `in` is still retained as an authored atom on every row it admits; only its expansion is taken over." A10's MVV assertion: "each row carrying both the authored `in` and the generated `eq`".
- evidence: `internal/table/model.go::KernelRow` — `if a.Block == BlockMatch { match = append(match, resolve.Tag{Key: a.Key, Value: seamValue(a.Literal, r.isSet(a.Key))}) }`; `::seamValue` non-set arm: "A non-set value is one member by construction: load refuses a multi-member literal … at every authoring site … return members[0]". `0002:C13`: "In each expanded row the `in` atom becomes an `eq` atom on the chosen member … so every match-block atom the kernel receives is an equality the `Match` pattern can test."
- why it matters: a stepped rule with `match.<tag> in = [1, 3]` would reach the kernel as `Tag{tag, "1"}` on every expanded row — the cell-3 row can never match. The subsumed `in` must be rewritten per row to `match eq = <cell>` (what `expand`'s `case expanding:` does today), not retained; A10's oracle must assert that shape.

### F2 — Scenario 5b cannot produce `graph-overlap`; the over-admitted row is dead
- anchor: 0030:A5 | 0030:MVV item 5 | 0030:S1 | §mini-check-tables (`oracle` row 1-2, `disposition` `unless` row) | §risks-and-mitigations
- class: refuted
- claim: an `unless` on the stepped tag over-admits an interior cell another row claims; "that case lints `graph-overlap`" (A5), "loads, and lints `graph-overlap`" (MVV 5), "surfaces … as blocking `graph-overlap` at lint" (Risks).
- evidence: C1 retains the authored atoms, `unless` included. `internal/guard/product.go::acceptedIn`: `if atom.Block == table.BlockUnless { unlessTerms = append(...) }` … `accepted = accepted.subtract(excluded)`. For a row carrying `guard.all eq 3` and `unless eq 3` the accepted set is empty and projectable (`assignment.go::subtract` returns `newSet(s.dims)`). `internal/graphlint/groups.go::emitOverlaps`: `if !otherOK || left.Intersect(other).Len() == 0 { continue }`. `::checkRedundantRows`: `properSubset(mine, sibling)` is true for an empty set against any non-empty sibling → `CodeRedundantRow` (advisory).
- why it matters: the MVV scenario, A5's scope bound, the `oracle` negative control and the `disposition` row all assert a blocking finding that the lint cannot emit; the real consequence of an over-admitted interior cell is an advisory `graph-redundant-row` (or nothing, if no sibling), and the real cost of "`unless` not consulted" is C2's bound refusal on an over-admitted CAP cell — which the Risk already names. The negative control for MVV rows 1-2 is therefore not a control.

### F3 — A predicate-literal `{ step = 1 }` refuses under `malformed_predicate_atom`, not `malformed_tag_declaration`
- anchor: 0030:S4 | 0030:MVV item 4 | §mini-check-tables (`disposition`, `oracle` row 5)
- class: refuted
- claim: "a predicate literal `{ step = 1 }` refuses `malformed_tag_declaration`" (MVV 4); S4 expects "the predicate one `malformed_tag_declaration`".
- evidence: `internal/table/normalize.go::loader.atom` — `badAtom := func(detail string) (Atom, error) { return Atom{}, fail(CatMalformedPredicateAtom, where+": "+detail) }` … `members, err := valueMembers(raw); if err != nil { return badAtom(err.Error()) }`. `valueMembers`' default arm returns "value %v is not a tag value" for an inline table. Only the write path (`::renderWrites`) files that error under `CatMalformedTagDeclaration`.
- why it matters: S4 says "asserting the code, not merely 'refuses', is what pins the clause" — the pinned code is wrong, so the MVV would fail against a correct implementation or force a category change the record does not intend. C1's `[initial]` paragraph names only the `[initial]` category; the predicate category should be stated there as `malformed_predicate_atom` (`0003:C6`'s site rule) and S4/MVV 4/disposition corrected.

### F4 — The extracted shim in `internal/resolve` cannot reach `guard.Evaluator`
- anchor: 0030:§phase-2-expand-into-cells | §existing-infrastructure-audit (`valueSatisfies` row) | §mini-check-tables (`authority` row 1)
- class: not-found
- claim: "Extract the two-valued core of `valueSatisfies` — render the atom's literal, call the evaluator, return its three-valued verdict — to `internal/resolve` … it needs only `resolve.GuardAtom` and a `[]string` literal, so no `table` type is named and the cycle … does not arise."
- evidence: `internal/guard/product.go::valueSatisfies` — `seam := Evaluator{}` … `seam.Evaluate(...)`; `internal/guard/grammar.go:1 package guard`, `:97 type Evaluator struct{}`. `go list -f '{{.Imports}}'`: `resolve` imports no intra-repo package; `table` imports `resolve` only; `guard` imports `resolve` and `table`. A4 fixes `guard.Evaluator` as "the sole implementation of `resolve.GuardEvaluator`".
- why it matters: the core's third dependency is the evaluator, and it lives one package ABOVE `table`. As written, the loader's admits filter cannot call "the same comparison the runtime does" without either moving `Evaluator` into `resolve` (grammar.go imports only `resolve` + stdlib, so it is movable — but that is a larger change A4 and the audit do not name) or injecting a `resolve.GuardEvaluator` into the loader (an API change the record does not name). Phase 2 must pick one and say so.

### F5 — "Colliding their fingerprints" is not what a `Writes`-only stepped literal does
- anchor: 0030:C1 (both-carriers paragraph) | 0030:A11
- class: refuted
- claim: "A stepped literal written to `Writes` only would leave every expanded row of one rule sharing one `NextTags`, colliding their fingerprints"; A11 If-wrong: "expanded rows of one stepped rule collide on `Fingerprint`"; A11 oracle: "assert distinct `Fingerprint` values across a stepped rule's expanded rows".
- evidence: `internal/graphlint/engine.go::Fingerprint` writes every `row.Atoms` entry (`Key|Block|Operator|Literal;`) before `#` and the `NextTags`; C1 gives each expanded row its own `guard.all eq = <cell>` atom, so the atom prefix already differs per row.
- why it matters: the obligation (both carriers) still stands on the second half — "publishing the unstepped value as the successor" — but A11's Fingerprint assertion passes with the defect present and so verifies nothing; keep the `graph` per-row `next` assertion and drop or replace the Fingerprint one.

### F6 — Always-append contradicts `0002:C13`'s "non-empty exactly when more than one row"
- anchor: 0030:C1 (always-append clause) | 0030:D-identity
- class: contradiction
- claim: "The suffix element is appended at EVERY admitted cell, including when exactly one is admitted … The identity tuple stays total over the product (`0002:C13`) either way."
- evidence: `0002:C13` (record status Implemented): "A suffix is non-empty exactly when the rule produced more than one row." `normalize.go:791`: `suffixed := expanding && len(members) > 1`.
- why it matters: totality is preserved, but C1 breaks a peer's fenced iff while citing that peer for support. Per the joint-decision rule this is either a narrowing C1 must name ("`0002:C13`'s iff is scoped to `in` expansion; a step point always contributes an element") or a JDR. As written a reader of 0002 alone is told something 0030 makes false.

### F7 — `intWidth` failure is reported by lint as a saturating ceiling, not "no domain"
- anchor: 0030:C1 (width clause) | 0030:F1
- class: refuted (gloss only — the rule itself matches)
- claim: "a declaration lint reports as having no domain must not be one the loader enumerates" (C1); "the declarations lint already reports as having no domain" (F1).
- evidence: `internal/guard/declaration.go:127-135` — `width, ok := intWidth(*d.Min, *d.Max); if !ok { // … Reporting the ceiling makes Cardinality saturate into the too-large refusal … reporting no domain at all would instead misreport … return cardinalityCeiling, true }`. Only `::IntDomain` returns `nil`.
- why it matters: low; the refusal rule C1 restates is exactly `intWidth`'s, but the justification names the opposite of what `Cardinality` does. Reword to "a declaration lint refuses as over-large (`graph-product-too-large`) must not be one the loader enumerates".

### F8 — The MVV projection is stated two ways
- anchor: 0030:§mini-check-tables (`oracle` row 1-2; `trace` step 8) vs 0030:S1 / 0030:A5 / 0030:MVV item 2
- class: contradiction (internal)
- claim: tables: "a different `(code, key, dimension)` multiset"; "`(code, key, dimension)` equal". S1/A5/MVV: "PROJECTION `(code, key, dimension, class, reason)` … compared as sets".
- evidence: the two spellings in the record itself (`/tmp/full0030.md` lines 658, 704 vs 1218, 1245).
- why it matters: low; two tuples and set-vs-multiset in one record for one assertion is the "two clauses saying one thing differently" this pass was asked to catch. Pick the S1 form in the tables.

### F9 — Inherited-context `match in` on the stepped tag is not covered by the subsumption rule
- anchor: 0030:C1 (subsumption paragraph)
- class: silence
- claim: the admits filter conjoins "every positive atom the rule AUTHORS on that tag"; the step point "SUBSUMES a match `in` on that tag".
- evidence: `0002:C13`: "Every `in` atom in a rule's match blocks — local `match` and inherited context `match`, on `recognized` or on any other declared tag — expands"; `normalize.go::expand` receives the merged `predicates` and does not distinguish local from inherited.
- why it matters: an inherited context's `match.<tag> in` would still be a second choice point on the stepped tag unless "authors" is read to include inherited atoms; if it is not subsumed, the off-diagonal unreachable rows C1 argues against reappear. One sentence ("authored locally or inherited from a context") closes it.

### F10 — Zero-cell and `#`-member refusals carry no category in C1
- anchor: 0030:C1 (last sentences of the expansion half) | §mini-check-tables (`disposition`)
- class: silence
- claim: "A rule admitting zero cells is refused at load. A step write over a domain containing a member with the suffix separator `#` is refused at load." — no category; the `disposition` table alone assigns `malformed_tag_declaration`.
- evidence: `internal/table/category.go` — the category is wire vocabulary (`Categories()`), and `0003:C6` files its unsatisfiable-atom analogue under `malformed predicate atom`.
- why it matters: low; the contract leaves the emitted category to a non-normative table, and the table's choice differs from the disposition it cites. State the category in C1.
