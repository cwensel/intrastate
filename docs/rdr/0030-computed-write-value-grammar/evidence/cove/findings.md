Model: claude-fable-5-1 + claude-opus-5 (dual-model cove, fanned out)

# CoVe pre-lock lens — cli/0030 — consolidated origin ledger

Two independent passes: `findings-a.md` (fable) and `findings-b.md` (opus),
neither seeing the other. Step 0 grounding agreed on every CONFIRMED row.
Rows below are the union, deduped by element id; `both` marks a convergent
hotspot (independent agreement, not a re-dump).

| # | anchor | class | finding | pass |
| --- | --- | --- | --- | --- |
| L1 | 0030:A3 | not-found | A3's Evidence cites `internal/table/normalize.go::normalizeAtom`; no such symbol. The arm is `func (l *loader) atom(...)` at `normalize.go:43`, `#` ban at `:120-128`. The claim A3 makes is TRUE at the three real sites; only the symbol name is wrong. | both (A-F3, B-F1) |
| L2 | 0030:C3 | silence | C3 says surfaces read expanded rows "through `Row.Writes` and `Row.Atoms` alone". `Row.NextTags` is populated independently (`model.go:293-295`, `0002:C15`) from the same `writes` slice (`normalize.go:836-837`) and is read WITHOUT `Writes` by `graphlint/engine.go:147` (`Fingerprint`), `cli/graph_document.go:219`, `cli/flow_next.go:397`, `guard/lint.go:183`, `table/dump.go:125`. If the stepped literal lands in `Writes` only — which is what C1 specifies — every expanded row keeps identical `NextTags`: `Fingerprint` collides across all N rows, `graph`'s `next` publishes the unstepped value, `reach.go::successor` traverses wrong. | B-F2 |
| L3 | 0030:C1 | contradiction | C1 says admits are evaluated "against the atoms as already chosen in the expanded row — a match `in` on the stepped tag has expanded to its `eq` member first". `compareAtoms` (`normalize.go:266-277`) orders by key then `string(a.Block)`, and `BlockAll="all" < BlockMatch="match"` (`internal/resolve/guard.go:16,24`). A step point keyed as its emitted `guard.all` atom sorts BEFORE the match `in` on the same key, so the `eq` member is NOT yet chosen. Admitting against the authored `in` set then expanding it mints \|members\|² rows, \|members\|²−\|members\| of them carrying `match eq = a` beside `guard.all eq = b` — dead rows. A1's "slot into `compareAtoms`' sort" and C1's composition claim cannot both hold. | A-F1 |
| L4 | 0030:S1 | contradiction | S1/A5 assert "identical finding sets as sets over `findings[]` (`model` normalized)". Findings carry `Rule`, `Span`, `Fingerprint` (`graphlint/groups.go:355-366`; `clierr.go:366,367,386`), and `Fingerprint` hashes every atom (`engine.go:131-145`). Literal row `retry-2` carries `{attempt all eq 2}`; the step row `retry#2` carries `{attempt all eq 2, attempt all lt 5}` (C1 retains authored atoms) under a different rule id and locator. The assertion is unsatisfiable unless both sets are empty — but S1 also says the comparison spans the advisory tier. | A-F2 |
| L5 | 0030:C1 | refuted | Technical Design says a zero-admitted-cell rule "is the unsatisfiable-atom case `0003:C6` already refuses at load". It does not: `load.go:1834-1838` ordered-operator arm returns `conformKind` only, never `conformDomain` ("kind-checked but not domain-checked", `0002:C17`); no `unsatisf*` refusal exists in the loader. `[rule.guard.all.attempt] lt = 0` over `{0..9}` loads today as a never-matching row. The zero-cell refusal is NEW. | A-F5 |
| L6 | 0030:C1 | refuted | C1: "An `int` whose declared `min..max` spans the full integer width is refused too: its domain is unenumerable". `guard/declaration.go:190-199` refuses `span < 0 \|\| span == math.MaxInt` — a WIDER set than full-width (e.g. `{0..MaxInt}`), and the reason is width-arithmetic representability, not enumerability. By Performance Expectations' own ~5.8 KB/row, `{0..2^40}` is "enumerable" under C1 yet needs ~6 PB, so enumerability is not the separating property. `guard` is unreachable from `table` (`go list -deps`), so the loader must re-state the rule — and the record re-states a different one. Two authorities would disagree on `{0..MaxInt}`. | A-F7 |
| L7 | audit `valueSatisfies` row / Phase 2 | refuted | The audit calls `graphlint/reach.go::atomAdmitsValue` a copy of `guard/product.go::valueSatisfies`, and Phase 2 "repoints" both at one extracted shim. `reach.go:502-514` is `valueSatisfies` PLUS an admit-on-UNEVALUABLE collapse (`return verdict != resolve.GuardFalse`, documented "the false-green direction") — the OPPOSITE of C1's exclude-on-unevaluable. `guard/product.go:591-597` `acceptedIn` takes a third direction (unprojectable set). Three callers, three third-arm policies; the extraction is behaviour-preserving only if each caller keeps its own collapse at the call site. "C1's 'as the runtime evaluator would' then holds structurally" overstates: what holds is the two-valued core, never the third arm — the arm C1 spends a paragraph on. | both (A-F6, B-F3) |
| L8 | audit `valueSatisfies` row | refuted | The audit says "the set renderer twice (`guard/assignment.go::renderSet`, `graphlint/reach.go::renderSetLiteral`)". They are different functions: `renderSet` (`assignment.go:353`) marshals members as authored; `renderSetLiteral` (`reach.go:519`) first applies `canonicalValues` = `slices.Compact(slices.Sorted(...))` plus nil→`[]`. Folding them silently changes one of the two `in`-literal comparisons. Phase 2 must name which the loader takes. | B-F4 |
| L9 | audit `valueSatisfies` row | refuted | The audit's "duplicated three ways … a test-local copy in `internal/resolve`" is wrong on the third: `guard_mvv_test.go:486-510` `conformingContractSeam` is a second EVALUATOR (parses ints itself, never renders or calls `guard.Evaluator`); `guard_fixtures_test.go:268-278` `d13Set` is only a renderer. The shim exists twice. The third "copy" must stay independent to be a contract stand-in — it cannot be repointed. | A-F4 |
| L10 | 0030:C1 | silence | C1's enum step is position-indexed ("the domain member `n` positions from the cell") over a domain `tagDecl` never checks for duplicates or emptiness: `load.go:848-905` has only `len(src.Domain) > 0 && src.Kind != "enum"` and the unsorted carry `Domain: src.Domain`; `firstDuplicate` (`normalize.go:200`) is called only on atom members (`:113`, `:139`). `domain = ["a","b","a"]` is authorable on `main` — "n positions from cell `a`" has two answers, and the suffix `retry#a` is ambiguous across two cells (the exact collision `#` was banned to prevent). C1 names a `#`-member refusal it needed A3 to discover; this is the same class and is unstated. | B-F5 |
| L11 | 0030:C1, 0030:S4 | refuted | C1: "`[initial]` values and predicate literals keep the literal-only grammar" — read as "refuses the same way". `load.go:1688-1690` files the `[initial]` refusal under `CatMalformedInitialDeclaration`, not `CatMalformedTagDeclaration` (the write path, `normalize.go:615,625,638`). S4 asserts only "both refuse" with no code, so the MVV passes whichever fires. | B-F6 |
| L12 | 0030:D-identity | contradiction | D-identity states the single-cell suffix reversal ("a stepped row is never the authored row, so it always carries the cell — the opposite of the single-member `in` rule") as a Load-Bearing Decision parenthetical. `normalize.go:131-132` is `suffixed := expanding && len(members) > 1` — ONE boolean shared by every choice point, so this is a real code fork: `suffixed` stops being a function of `len(members)` and becomes a function of the choice-point KIND, varying the identity-tuple totality argument `0002:C13` rests on. C1's expansion paragraph is silent on the single-cell case. It belongs in C1's fence where a peer can cite it. | B-F7 |
| L13 | 0030:A4 | refuted (narrow) | A4's Evidence says "`internal/table/load.go::conform` kind-checks their bound". Correct as to the parse — `conformKind` IS `strconv.Atoi` for `int` — but `conform` has two arms and the ordered one (`:1834-1838`) skips `conformDomain` entirely. Live consequence the record does not state: `lt = 500` on `min=0,max=9` loads today, so a step rule may author an ordered atom whose bound is outside the domain (harmless for admission — the conjunction is per member — but "the rule's own positive atoms" can be vacuous in a way `conform` never reports). Narrow the Evidence to `conformKind`. | B-F8 |
| L14 | 0030:C1 (`unless` para), §risks-and-mitigations | silence | C1: an over-admitted cell "surfaces as the bound refusal with the `guard.all` complement as the named remedy"; Risks repeats it. The bound refusal fires only when `cell + n` leaves the domain. An over-admitted INTERIOR cell (`attempt` `0..9`, `guard.all.attempt lt = 5`, `guard.unless.attempt eq = 3`, step 1) yields a row for cell 3 writing 4 — conforms, no refusal. Its `guard.all eq = 3` atom then OVERLAPS the row the author wrote to claim cell 3: `graph-overlap` at lint (blocking finding, not a refusal, and not the one C1 names). A5 ("no overlap the hand-unrolled table did not have") is Pending and its MVV fixture pair carries no `unless` on a stepped tag. | B-F9 |
| L15 | 0030:C2 | silence | C2's detail "names the rule, the tag, the cell, and the stepped value"; F3 says a too-large step "refuses at every admitted cell under C2, naming the first". The record never states the order "first" is taken in (domain order ascending for `int`, authored order for `enum`, or the row's suffix sort). S3's single-failing-cell expectation is consistent with either reading; S5's oversize-step fixture is not. | A-F8 |
| L16 | 0030:§finalization-gate | silence | The gate section carries TEMPLATE.md's bracketed guidance verbatim (1053-1172); `recs lint` reports `gate:inline` + five `placeholder:survived`; the `artifact` edge at 1056 still points at the unexpanded `{ARTIFACT_DIR}/gate.md`. Noted for completeness — this is pre-lock state the Finalization Gate resolves, NOT a design defect, and NOT this stage's obligation. | B-F10 |

## Convergence note

L1 and L7 were found independently by both models — real agreement, a hotspot.
Each pass also found blocking material the other missed (L3 sort ordering,
L2 `NextTags`), which is the dual-model draw paying for itself.

## Dispositions (iteration 1)

| # | anchor | disposition | section touched |
| --- | --- | --- | --- |
| L1 | 0030:A3 | fixed — symbol re-pointed to `normalize.go::loader.atom` | Critical Assumptions |
| L2 | 0030:C3 | fixed — C1 gains a both-carriers clause (`Writes` + `NextTags`, `0002:C15`); C3, A1, A9 and the audit row amended; A11 booked Pending | Normative Contracts, CAs, audit |
| L3 | 0030:C1 | fixed — step point SUBSUMES a match `in` on the stepped tag rather than composing; A1's sort fact stated; A10 booked Pending | Normative Contracts, A1 |
| L4 | 0030:S1 | fixed — finding-set comparison projected to `(code, key, dimension, class, reason)`; A5, MVV and Consequences aligned | Testing Strategy, A5, MVV |
| L5 | 0030:C1 | fixed — zero-cell refusal restated as NEW, taking `0003:C6`'s disposition rather than its check | Technical Design, Failure Modes, Key Discoveries |
| L6 | 0030:C1 | fixed — width refusal restated as `intWidth`'s representability rule, not enumerability | Normative Contracts, Failure Modes |
| L7 | audit / Phase 2 | fixed — extraction scoped to the two-valued core; each of three callers keeps its own third-arm policy at the call site | audit, Phase 2 |
| L8 | audit | fixed — two renderers, not two copies; loader takes the canonicalizing form | audit |
| L9 | audit | fixed — shim exists twice; the `resolve` test item is a second evaluator, not a copy | audit |
| L10 | 0030:C1 | fixed — duplicate enum domain member refused, same class as the `#` guard | Normative Contracts |
| L11 | 0030:C1, S4 | fixed — `[initial]` keeps `malformed_initial_declaration` | Normative Contracts, Testing Strategy, MVV |
| L12 | 0030:D-identity | fixed — always-append hoisted into C1's fence; D-identity defers | Normative Contracts, Load-Bearing Decisions |
| L13 | 0030:A4 | fixed — Evidence narrowed to `conformKind`; the un-range-checked bound stated | Critical Assumptions |
| L14 | 0030:C1, Risks | fixed — interior over-admission separated from the bound refusal | Technical Design, Risks |
| L15 | 0030:C2 | fixed — first failing cell defined in the tag's own domain order | Normative Contracts, Failure Modes |
| L16 | §finalization-gate | dismissed-with-cite — Stage 7's obligation, not a design defect; `recs lint` reports it as the expected pre-lock state | none |

Mini-checks fired: authority, oracle, fidelity, disposition, trace (all five;
tables written into Technical Design §Mini-check tables). No CONTRADICTION row
survived.

Loop: `lens-loop` answered `rerun` (substantial rewrite of C1) → iteration 2.

## Dispositions (iteration 2 — delta re-run, `iter-2/findings.md`)

The delta pass verified the 12 iteration-1 fixes against source: 8 sound,
4 defective. All 10 of its findings dispositioned:

| # | anchor | disposition | section touched |
| --- | --- | --- | --- |
| F1 | 0030:C1, A10 | fixed — a subsumed `in` is REWRITTEN per row to `match eq = <cell>`, never retained; `KernelRow`/`seamValue` would otherwise truncate it to its first member on every row | Normative Contracts, A10 |
| F2 | 0030:A5, MVV, Risks, tables | fixed — an over-admitted interior cell yields a DEAD row (retained `unless` empties its accepted set), lint advisory `graph-redundant-row`, never blocking `graph-overlap`; A5's claim un-scoped as a result | Technical Design, A5, Risks, MVV, Testing Strategy, tables |
| F3 | 0030:S4, MVV, tables | fixed — a predicate literal refuses `malformed_predicate_atom`; all three paths' categories now named in C1 | Normative Contracts, Testing Strategy, MVV, tables |
| F4 | Phase 2, audit, tables | fixed — `guard/grammar.go::Evaluator` moves to `internal/resolve` with the shim (it imports only stdlib + `resolve` and names no `table` type); A12 booked Pending | Phase 2, audit, tables |
| F5 | 0030:C1, A11 | fixed — the collision motive was wrong (`Fingerprint` hashes `Atoms` first); restated as a wrong-successor defect and A11's oracle repointed at `graph`/`dump` `next` | Normative Contracts, A11 |
| F6 | 0030:C1, D-identity | fixed — C1 names the narrowing of `0002:C13`'s suffix iff and why a step point has no one-edge twin; resolved as a scoped narrowing, not a JDR (C13's own reason does not reach this case) | Normative Contracts |
| F7 | 0030:C1, F1 | fixed — gloss corrected: lint saturates `Cardinality` to its ceiling, it does not report "no domain" | Normative Contracts, Failure Modes |
| F8 | tables vs S1 | fixed — one projection spelling, `(code, key, dimension, class, reason)` as sets, everywhere | Mini-check tables |
| F9 | 0030:C1 | fixed — subsumption covers an inherited-context `in`, not only a locally authored one | Normative Contracts |
| F10 | 0030:C1 | fixed — every refusal this clause mints is on the write-block path, `malformed_tag_declaration`, stated in C1 | Normative Contracts |

Net-new anchors at iteration 2 (A10, A11, A5, F1, MVV, mini-check tables,
Phase 2) are all elements iteration 1's own fixes created or rewrote — the
delta re-run checking its predecessor's work, not scope expansion. Nothing
was charted to a successor and nothing was absorbed that the ledger did not
originate.

Loop: `lens-loop` answered `converged`.
