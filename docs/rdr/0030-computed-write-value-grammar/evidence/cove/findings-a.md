Model: claude-fable-5-1

# CoVe pre-lock lens — RDR 0030 (findings-a)

Record read whole via the projector (metadata, problem statement, critical
assumptions, approach, technical design, audit, rationale, alternatives,
context, research, trade-offs, plan, validation, gate, references).

## Step 0 — grounding ledger

Scope: `recs inspect --json --filter edges,elements 0030` yields 60
`source-anchor` edges, ALL `resolved: true` — the projector worklist is
empty. Decision Rationale carries `Ground-sweep: clean (33 anchors)`, so
27 anchors and the surrounding prose were added after propose; those, plus
every claim that mints no edge, were swept by hand against the working
tree (`main`).

| # | Claim (element) | Verdict | Cite |
| --- | --- | --- | --- |
| 1 | `expand`'s discriminant is `Operator == "in" && Block == BlockMatch`; `case expanding:` mints a fresh `eq` atom; members land in `Row.Suffix` (A1) | CONFIRMED | `internal/table/normalize.go:784` `expanding := c.atom.Operator == "in" && c.atom.Block == BlockMatch`; `:804-812` appends `Atom{Operator: "eq", Literal: []string{member}}`; `:795-797` `next.Suffix = append(...)` |
| 2 | `renderWrites` runs BEFORE `expand`; the tail clones one write set onto every row (A1) | CONFIRMED | `normalize.go:471` `return expand(base, predicates, outcomeAtom, writes)`; `:830-837` `rows[i].Writes = cloneTagValues(writes)` |
| 3 | `expand`'s seed copies `Gate`, `RequiresOwned`, `Escape`, `Emit`; one `SourceLocator` shared (audit row, rationale) | CONFIRMED | `normalize.go:764-775` |
| 4 | `compareAtoms` sort is (key, block, operator, literal) (A1, D-identity) | CONFIRMED | `normalize.go:266-277` |
| 5 | `tagDecl` carries `Domain: src.Domain` unsorted and performs no `#` check (A2, A3) | CONFIRMED | `internal/table/load.go:896-905`; no `suffixSep`/`'#'` reference in `tagDecl` |
| 6 | No sort over `TagDecl.Domain` in table/graph_document/graphlint; the one domain sort is `EmitDecl.Domain` (A2) | CONFIRMED | sort sweep: `load.go:377 slices.Sort(members)` inside the emit-decl builder (comment `:372` "value-preserving, not order-preserving"); `graph_document.go:187 tag.Domain = slices.Clone(decl.Domain)` |
| 7 | The `#` ban has exactly three sites: `normalizeAtom` `case "in":` gated `BlockMatch`, `normalizeRules`, `loadOutcomes` (A3) | NOT-FOUND (symbol) / CONFIRMED (fact) | three `suffixSep` sites: `normalize.go:119-127` (inside `func (l *loader) atom`, line 43 — there is no `normalizeAtom` symbol in `internal/table`), `normalize.go:297` (`normalizeRules`), `load.go:163` (`loadOutcomes`) |
| 8 | `Identity` joins model id with `.` and each suffix element with `#` (ground-sweep note) | CONFIRMED | `internal/table/model.go:337-347` |
| 9 | `Evaluator` is `struct{}`; `lt/lte/gt/gte` `strconv.Atoi` both sides → `compare`; `eq` string equality; `in` `slices.Contains` over a parsed set; unknown operator → `GuardUnevaluable` (A4) | CONFIRMED | `internal/guard/grammar.go:97,104-152,198` |
| 10 | `Evaluator` is the sole non-test `GuardEvaluator` implementation; `guardcontract.go` is a harness (A4) | CONFIRMED | non-test `Evaluate(atom resolve.GuardAtom` only at `grammar.go:104`; `internal/resolve/guardcontract.go:14 func TestGuardEvaluatorContract(t *testing.T, seam GuardEvaluator)` |
| 11 | `operatorKinds` confines `lt/lte/gt/gte` to `int`; claims to mirror `guard.Accepts` with no cross-package test (A4, rationale) | CONFIRMED | `model.go:110-119`; comment `:100-104` "`guard.Accepts` is the authority and its REQ-4 test pins the matrix" |
| 12 | `conform` kind-checks a comparison bound and domain-checks `eq`/`in` members (A4, C2) | CONFIRMED | `load.go:1826-1850` `case "lt","lte","gt","gte": … return conformKind(decl, members[0])` |
| 13 | "A rule with zero admitted cells … is the unsatisfiable-atom case `0003:C6` already refuses at load" (Technical Design, C1) | REFUTED (partial) | the loader refuses a domain literal outside the domain (`conform` → `conformDomain`) but a BOUND atom is never domain-checked (`load.go:1840-1844`, comment "kind-checked but not domain-checked (`0002:C17`)"); `rg unsatisf internal/table` → no loader site. `lt = 0` on `{0..9}` loads today. See F5 |
| 14 | `valueMembers` admits string/bool/int64/float64/[]any, refuses the rest as "is not a tag value" (problem statement) | CONFIRMED | `load.go:1748-1786` |
| 15 | `renderWrites` prefixes `rule <id> write <key>: ` and files kind/domain conformance under `malformed_tag_declaration` via `conform(decl, "eq", members)` (S3, D-naming) | CONFIRMED | `normalize.go:614-626` |
| 16 | `Row.KernelRow` builds `resolve.Row{RuleID: r.RuleID, …}` dropping `Suffix`; flow verbs publish that bare id; `rowByID` matches `row.RuleID == ruleID` (A8) | CONFIRMED | `model.go:383-384`; `flow_next.go:337 Rule: row.RuleID`; `flow_resolve.go:306 Rule: plan.RuleID`, `:408-410` |
| 17 | dump `case "identity": return r.Identity()`; graph `graphRowDoc.Identity`; edge `Rule` field (A8, C3) | CONFIRMED | `dump.go:114-115`; `graph_document.go:87,137,214` |
| 18 | `flow_next.go` carries the "answered a question it does not own … disagreed with the kernel silently" comment (rationale) | CONFIRMED | `internal/cli/flow_next.go:436-437` |
| 19 | `docs/cli-output-contract.md` says "a candidate is a rule" (A8) | CONFIRMED (near-quote) | `docs/cli-output-contract.md:620` "A candidate is a rule whose match and guard both HOLD or are UNDECIDED" |
| 20 | `go list -deps ./internal/table` names neither `guard` nor `graphlint`; `internal/resolve` imports nothing intra-repo (A7, audit) | CONFIRMED | `go list -deps ./internal/table` → resolve, table only; `go list -f '{{join .Imports "\n"}}' ./internal/resolve` → no intra-repo import |
| 21 | No row-count gate in `internal/table` (A7) | CONFIRMED | `rg 'too many|ceiling|len\(rows\) >' internal/table` → none (the only "truncated" text is the dump-column vocabulary at `load.go:1606`) |
| 22 | `valueSatisfies` is a ~13-line render-then-`Evaluate` shim; `atomAdmitsValue` and `renderSet`/`renderSetLiteral` duplicate it; plus "a test-local copy in `internal/resolve`" (audit, rationale) | CONFIRMED ×2 / REFUTED for the third | `guard/product.go:541-556`; `graphlint/reach.go:502-515`; `guard/assignment.go:353`, `reach.go:519`. In `internal/resolve` tests there is no render-then-`Evaluate` copy: `guard_mvv_test.go:486-510 conformingContractSeam` is a second EVALUATOR, `guard_fixtures_test.go:268-278 d13Set` a set renderer. See F4 |
| 23 | `atomAdmitsValue` collapses the seam's third verdict (A4/C1 third-arm context) | CONFIRMED, direction ADMIT | `reach.go:514 return verdict != resolve.GuardFalse` — the opposite of C1's exclude. See F6 |
| 24 | `Denotation`/`acceptedIn` yield the unprojectable set on an undecidable dimension (audit) | CONFIRMED | `product.go:591-597` `if CanRefuse(m.Tags, row) { return AssignmentSet{} }` |
| 25 | `intWidth` "already reports that span as unusable"; `IntDomain` enumerates `min..max` (F1, audit) | CONFIRMED with a set mismatch | `declaration.go:190-199` refuses `span < 0 || span == math.MaxInt` — wider than "spans the full integer width". See F7 |
| 26 | `AssignmentCount` is lint-side and saturates at `cardinalityCeiling` (A7) | CONFIRMED | `declaration.go:80-107` |
| 27 | Finding codes `CodeOverlap`, `CodeCoverageGap`, `CodeOwnedBeforeWrite` (blocking), `CodeIdempotentWrite` (advisory) exist (A5) | CONFIRMED | `graphlint/taxonomy.go:28,29,33,55,88-93,112` |
| 28 | `checkIdempotentWrites` compares against match-block `eq` only (audit) | CONFIRMED | `graphlint/groups.go:342-354` `atom.Block != table.BlockMatch \|\| atom.Operator != "eq"` |
| 29 | `flow_resolve` payload carries `next`, `writes`, `clear`; `--plan-only` exists (A6, S2) | CONFIRMED | `flow_resolve.go:89-91,126,163` |
| 30 | "state cell driven by rewriting the bound artifact file, the pattern `flow_exit_0005_test.go` uses" (A6) | CONFIRMED with a nuance | that file seeds a fresh artifact per state (`seedArtifact(t, model, "status=draft")` at `:81,143,183,291`, `artifactBinding` at `flow_fixtures_0005_test.go:372`); it binds distinct files rather than rewriting one — same capability, not literally "rewriting" |
| 31 | `TestReq92_ANonFirstExpandedRowResolvesToTheAuthoredEmitBlock`, `TestReq34_TheRuleIDJoinIsSoundOverAnExpandingRule` exist and assert the authored id (rationale, S7) | CONFIRMED | `internal/cli/decision_table_0010_test.go:859,888` |
| 32 | `<topic>_mvv_<rdr>_test.go` precedent (Testing Strategy) | CONFIRMED | 12 files, e.g. `internal/cli/lint_mvv_0006_test.go`, `graph_mvv_0021_test.go` |
| 33 | `rdr-fixture.toml` `[tags.iter]` int `min = 0, max = 9` and enum domains (Testing Strategy) | CONFIRMED | `internal/table/testdata/rdr-fixture.toml:40-44,19-57` |
| 34 | `docs/model-authoring.md` §Discriminate with guard atoms quote "contributes **no dimension**"; §Guard and match atoms exists (Technical Design) | CONFIRMED | `docs/model-authoring.md:1077,1091-1094`; `:193` |
| 35 | `strconv.Atoi("+1")` is `1` (D-wire) | CONFIRMED | `go run` → `1 <nil>` |
| 36 | `engine.go` canonicalizes `NextTags`; `reach.go::successor` clones `row.Writes`; `Fingerprint` hashes atoms (problem statement, audit) | CONFIRMED | `engine.go:147 canonicalTags(row.NextTags)`, `:131-145`; `reach.go:304,326` |
| 37 | `internal/table::Categories()`, `guard/lint.go::writtenKeys`, `accessor/executor.go`, `flow_state.go` read `Writes` (rationale ten-site list) | CONFIRMED | `category.go:93`; `guard/lint.go:177`; `flow_state.go:418,599-618` |
| 38 | Peer ids `0002:C4/C7/C11/C13/C17/C19/C22/C24`, `0003:C3/C6`, `0005:C1`, `0007:C1/C2`, `0010:C3`, `0012:C2/A5`, `0013:C1/C3`, `0021:C2`, `0022:A1`, `0023:C2`, `0029:C4` exist; `0002:C17` scopes the `#` ban to `in` members "since only match blocks expand"; `0003:C6` says "rather than a silently-never-matching row"; `0021:C2` "AUTHORED members"; `0010:C3` "first-match `rowByID` is correct as it stands"; `0007:C2` partiality is over the assembled view | CONFIRMED | `recs inspect --select` on each |
| 39 | Evidence artifacts `evidence/spikes/a7-row-growth.md`, `evidence/research/ground-sweep.md`, `evidence/propose-premortem/critic.md`, `evidence/research/prior-art.md` exist | CONFIRMED | `ls` of the evidence tree |

Spans widened past the record's anchors (and why):
- `internal/resolve/guard.go:16-24` (`BlockAll = "all"`, `BlockMatch = "match"`, `BlockUnless = "unless"`) — never cited by the record; needed to test C1's "a match `in` on the stepped tag has expanded to its `eq` member first" against the sort the record says the step point slots into (F1).
- `internal/cli/clierr/clierr.go:366-386` (`Finding.Rule/Span/Fingerprint`) and `graphlint/engine.go::Fingerprint` — never cited; needed to test S1's "identical finding sets … normalizing `model`" (F2).
- `internal/resolve/guard_mvv_test.go`, `guard_fixtures_test.go` — to locate the claimed "test-local copy" of the shim (F4).
- `internal/table/load.go::conform` bound arm + `rg unsatisf` — to test "already refuses at load" for the zero-cell case (F5).

## Step 1 — questions

1. (source) In `expand`, does the existing choice-point sort place a `guard.all` atom BEFORE a `match` atom on the same key, so that a step point slotted "into `compareAtoms`' sort" is processed before the match `in` it is supposed to compose with?
2. (source) Does today's loader refuse a bound atom that admits no domain member (e.g. `[rule.guard.all.attempt] lt = 0` on `min = 0, max = 9`), as the record's "0003:C6 already refuses at load" implies?
3. (source) Does the symbol `internal/table/normalize.go::normalizeAtom` (A3, references) exist?
4. (source) Is there a test-local render-then-`Evaluate` copy of `valueSatisfies` in `internal/resolve`, making the shim "duplicated three ways"?
5. (source) Do lint findings carry row-specific fields that differ between `retry-2` (literal ladder) and `retry#2` (step ladder), making "identical finding sets" impossible unless the sets are empty?
6. (source) Which direction does the existing `graphlint::atomAdmitsValue` take on an UNEVALUABLE verdict — admit or exclude?
7. (source) Does `guard/declaration.go::intWidth` refuse exactly "an `int` whose `min..max` spans the full integer width" (C1), or a wider set?
8. (source) Is `internal/guard` (and `graphlint`) genuinely unreachable from `internal/table`, and does `internal/resolve` import nothing intra-repo, so the proposed extraction target creates no cycle?
9. (source) Does the graph export emit `domain` in authored order, and under what condition is it present at all?
10. (RDR) Does the record define the step point's sort tuple — its (key, block, operator, literal) or otherwise — relative to `in` atoms on other keys and the outcome atom, so two authorings yield one suffix order?
11. (RDR) Does the record say how the MVV's finding-set comparison treats `rule`, `span` and `fingerprint`, which name rows?
12. (RDR) Is C1's full-width refusal ("its domain is unenumerable") consistent with A7's "no load-time ceiling" and Performance Expectations' 5.8 KB/row figure — i.e. does the record say why `{0..2^40}` is enumerable but `{MinInt..MaxInt}` is not?

## Step 2 — answers

1. YES. `internal/resolve/guard.go:16` `BlockAll Block = "all"`, `:24` `BlockMatch Block = "match"`; `normalize.go::compareAtoms` compares `string(a.Block)` after the key, and `expand` (`:751-753`) sorts candidates with it and walks them in that order. `"all" < "match"`, so on the same key a guard-block choice point precedes the match `in`. If the step point sorts as the `guard.all eq` atom it emits (the only tuple the record hints at), it runs BEFORE the match `in` on that key has expanded — the `eq` member is not "already chosen". Evaluating admits against the authored `in` (set membership) and then letting the `in` expand mints |members|×|members| rows, the off-diagonal ones carrying `match eq = m1` beside `guard.all eq = m2` on one key — dead rows.
2. NO. `internal/table/load.go:1826-1850` `conform`: `case "lt", "lte", "gt", "gte": … return conformKind(decl, members[0])` — a comparison bound is "kind-checked but not domain-checked (`0002:C17`)". `rg -n 'unsatisf' internal/table/*.go` hits only a comment in `edit.go`. So `lt = 0` over `{0..9}` loads today as a never-matching row; the loader's `0003:C6` refusal covers domain LITERALS (`eq`/`in` outside `domain`/`min..max` via `conformDomain`), not bounds.
3. NO. `rg -n normalizeAtom internal/table/` → nothing. The `#` ban the record attributes to it is in `func (l *loader) atom(decl TagDecl, key, operator string, raw any, b Block, owner string)` at `normalize.go:43`, `case "in":` arm `:119-127`, gated `if b == BlockMatch`.
4. NO. `internal/resolve/guard_mvv_test.go:486-510` `conformingContractSeam.Evaluate` is a second typed EVALUATOR (parses ints itself, never calls `guard.Evaluator`); `guard_fixtures_test.go:268-278` `d13Set` is a set renderer; `valueSeam`/`eqSeam` (`:90-135`) are programmable stubs. The render-then-`Evaluate` shim exists exactly twice (`guard/product.go:541`, `graphlint/reach.go:502`), the set renderer twice in non-test code plus `d13Set` in a test.
5. YES. `graphlint/groups.go:355-366` emits `Rule: row.RuleID, Span: row.SourceLocator, … Fingerprint: Fingerprint(row)`; `clierr.go:366 Rule`, `:367 Span`, `:386 Fingerprint` are envelope fields; `engine.go:131-145 Fingerprint` hashes EVERY atom (key|block|operator|literal). A literal row `retry-2` carries atoms `{attempt all eq 2}`; the step row `retry#2` carries `{attempt all eq 2, attempt all lt 5}` (C1: "The authored atoms are retained") under rule `retry` at a different locator. Any finding either ladder emits differs on `rule`, `span` and `fingerprint`.
6. ADMIT. `internal/graphlint/reach.go:514` `return verdict != resolve.GuardFalse`, comment `:499-501` "An UNEVALUABLE verdict is taken as satisfiable: the relation over-approximates". `guard/product.go:591-597` `acceptedIn` takes a third direction (unprojectable set). C1's EXCLUDE is a third policy for the same arm.
7. WIDER. `internal/guard/declaration.go:190-199`: `span := maxV - minV; if span < 0 || span == math.MaxInt { return 0, false }`. That refuses `{MinInt..MaxInt}` (wrap) AND e.g. `{-1..MaxInt-1}` or `{0..MaxInt}` (`span == MaxInt`, the inclusive `+1` would wrap) — declarations that do not "span the full integer width". `IntDomain` (`:241-256`) returns nil for all of them.
8. YES. `go list -deps ./internal/table` → `internal/resolve`, `internal/table` only; `internal/resolve` has no intra-repo import; `internal/guard` imports `resolve`+`table`; `internal/graphlint` imports `clierr`, `guard`, `resolve`, `table`. A shim in `resolve` taking `resolve.GuardAtom` + `[]string` is reachable from all three without a cycle.
9. Authored order, conditionally. `internal/cli/graph_document.go:185-188`: `if _, finite := guard.AssignmentCount(decl); finite && len(decl.Domain) > 0 { tag.Domain = slices.Clone(decl.Domain) }` — unsorted clone; present only when `AssignmentCount` reports finite (always true for an `enum` with a non-empty domain that `agrees`).
10. RDR is silent on the tuple. D-identity says "in the choice-point sort order" and "Two step points in one rule order their suffix elements by the choice-point sort, as `in` atoms do"; A1 says the step point is "a third kind whose suffix element must slot into `compareAtoms`' sort". Neither names the (key, block, operator, literal) the step point sorts under, and C1's composition text presumes an order (`in` first) the natural tuple does not give (answer 1).
11. RDR is silent beyond `model`. S1: "identical finding sets as sets over `findings[]` (`model` normalized)"; A5 Evidence: "finding sets compared as sets over `findings[]` (normalizing `model`)". Consequences admits "their row identities differ (suffix vs authored id) … the MVV compares plans and findings, not dumps" — but findings carry the row identity too (answer 5).
12. Not consistent as written. C1: "An `int` whose declared `min..max` spans the full integer width is refused too: its domain is unenumerable, so the expansion has no cell set to compute over." A7/Performance Expectations: no ceiling; "~5.8 KB per row"; "Were a ceiling ever wanted it would be memory-motivated and first bite near 1,000,000 rows (~6 GB)". By the record's own numbers a `{0..2^40}` int is "enumerable" under C1 yet needs ~6 PB; the refusal line C1 draws is the width-arithmetic overflow `intWidth` guards (answer 7), not enumerability — the record gives the wrong reason for a rule it copies from `intWidth` with a narrower set.

## Step 3 — findings

### F1 — Step point cannot both "slot into compareAtoms' sort" and see the match `in` "already chosen"
- anchor: 0030:C1
- class: contradiction
- claim: C1 "evaluated … against the atoms as already chosen in the expanded row — a match `in` on the stepped tag has expanded to its `eq` member first, so the two compose to exactly its members"; A1 "a step point is a third kind whose suffix element must slot into `internal/table/normalize.go::compareAtoms`' sort"; D-identity "the cell appended as one more suffix element in the choice-point sort order".
- evidence: `internal/resolve/guard.go:16` `BlockAll Block = "all"`, `:24` `BlockMatch Block = "match"`; `internal/table/normalize.go:266-277` `compareAtoms` orders by `string(a.Block)` after the key; `expand` `:751-753` walks candidates in that order. A step point keyed as its emitted `guard.all` atom sorts BEFORE the match `in` on the same key, so the `eq` member is not yet chosen when cells are admitted; admitting against the authored `in` set and then expanding it mints |members|² rows, |members|²−|members| of them contradictory (`match eq = a` beside `guard.all eq = b`).
- why it matters: the row set C1 promises ("exactly its members") depends on an ordering the record never states and the cited sort does not give; an implementer following A1 literally ships dead rows and a wrong `identity` suffix order. The record must fix the step point's sort tuple (or defer admission until after the product) — a C1/D-identity edit, not a re-decision of the approach.

### F2 — S1/A5 "identical finding sets" is unsatisfiable unless both sets are empty
- anchor: 0030:S1
- class: contradiction
- claim: S1 "identical finding sets as sets over `findings[]` (`model` normalized), spanning both the blocking codes and the advisory `graph-idempotent-write`"; A5 Evidence "finding sets compared as sets over `findings[]` (normalizing `model`)".
- evidence: `internal/graphlint/groups.go:355-366` emits `Rule: row.RuleID, Span: row.SourceLocator, Fingerprint: Fingerprint(row)`; `internal/cli/clierr/clierr.go:366,367,386` publish them; `internal/graphlint/engine.go:131-145` `Fingerprint` hashes every atom, and C1 says "The authored atoms are retained" (step row atoms = `{eq cell, lt 5}` vs literal row `{eq 2}`). Rule ids differ by construction (`retry-2` vs `retry`).
- why it matters: as written the MVV either passes only when both ladders lint clean (then "spanning … the advisory tier" is vacuous) or fails on every non-empty finding. S1 needs to say which: assert empty finding sets, or normalize `rule`/`span`/`fingerprint` (and say the comparison is by `code`+`key`+`dimension`).

### F3 — `internal/table/normalize.go::normalizeAtom` does not exist
- anchor: 0030:A3
- class: not-found
- claim: A3 Evidence "`internal/table/normalize.go::normalizeAtom`'s `case "in":` arm, gated `if b == BlockMatch`"; A4 "`internal/table/normalize.go` refuses an unknown operator" (fine); References list.
- evidence: `rg -n normalizeAtom internal/table/` → no hits. The site is `internal/table/normalize.go:43` `func (l *loader) atom(decl TagDecl, key, operator string, raw any, b Block, owner string)`, `#` ban at `:119-127`. The projector marked the anchor resolved; it is not.
- why it matters: the fact A3 rests on holds (three `suffixSep` sites, `tagDecl` has none), but the cited symbol will not resolve for the implementer or the gate's "each cited `path::Symbol` resolves on `main`" check.

### F4 — The shim is duplicated twice, not "three ways"; the resolve test-local item is a second evaluator
- anchor: none (§Existing Infrastructure Audit, `valueSatisfies` row; Decision Rationale "where does the admits comparison live?")
- class: refuted
- claim: "It is already duplicated three ways (`guard/product.go::valueSatisfies`, `graphlint/reach.go::atomAdmitsValue`, a test-local copy in `internal/resolve`)".
- evidence: `internal/resolve/guard_mvv_test.go:486-510` `conformingContractSeam.Evaluate` parses ints itself and never renders a set literal or calls `guard.Evaluator` — a stand-in EVALUATOR, not the render-then-`Evaluate` shim; `guard_fixtures_test.go:268-278` `d13Set` is only a renderer. The shim exists at `guard/product.go:541` and `graphlint/reach.go:502` only.
- why it matters: low — the extraction decision survives on two copies plus two renderers; the count in the rationale is wrong and the third "copy" cannot be repointed at the extracted shim (it must stay independent to be a contract stand-in).

### F5 — The zero-cell refusal is a NEW load refusal, not one "0003:C6 already refuses"
- anchor: 0030:C1
- class: refuted
- claim: Technical Design "A rule with zero admitted cells is refused: it is the unsatisfiable-atom case `0003:C6` already refuses at load rather than yielding a never-matching row"; Key Discoveries "a zero-cell step rule takes the same disposition".
- evidence: `internal/table/load.go:1840-1844` `case "lt", "lte", "gt", "gte": … return conformKind(decl, members[0])` — "an ordered BOUND, which is kind-checked but not domain-checked (`0002:C17`)"; no `unsatisf*` refusal anywhere in the loader. `[rule.guard.all.attempt] lt = 0` over `{0..9}` loads today; only `eq`/`in` literals outside the domain refuse (`conformDomain`).
- why it matters: a step rule whose zero cells come from a BOUND atom (or from two bounds that conjoin to nothing) will be refused where the same atoms without a step write load — a behaviour the record presents as pre-existing. State it as a new refusal that follows C6's disposition, and add the bound-atom shape to S5 so the MVV covers it.

### F6 — The undecided-arm direction already has two shipped siblings, in two other directions
- anchor: 0030:C1
- class: sibling-exists
- claim: C1 "an atom answering neither true nor false at a member does NOT admit that member … a soundness fence"; audit row cites only `Denotation`/`acceptedIn`'s unprojectable-set direction.
- evidence: `internal/graphlint/reach.go:514` `return verdict != resolve.GuardFalse` (comment `:499-501`: UNEVALUABLE "is taken as satisfiable: the relation over-approximates"); `internal/guard/product.go:591-597` returns `AssignmentSet{}` (unprojectable). C1's EXCLUDE is a third policy for the same seam verdict, and the extracted shim (Phase 2) carries only render+`Evaluate` — each caller keeps its own collapse, so "as the runtime evaluator would … then holds structurally" covers the verdict, not its disposition.
- why it matters: the record chooses the direction deliberately and the arm is unreachable for the admitted kinds, so no re-decision is needed; but the rationale should name `atomAdmitsValue`'s opposite choice so the three policies are a recorded per-consumer split rather than an apparent inconsistency.

### F7 — C1's full-width refusal copies `intWidth` with a narrower set and the wrong reason
- anchor: 0030:C1 (also 0030:F1)
- class: sibling-exists
- claim: C1 "An `int` whose declared `min..max` spans the full integer width is refused too: its domain is unenumerable"; F1 "`internal/guard/declaration.go::intWidth` already reports that span as unusable".
- evidence: `internal/guard/declaration.go:190-199` `span := maxV - minV; if span < 0 || span == math.MaxInt { return 0, false }` — refuses the wrap AND `span == MaxInt` (e.g. `{0..MaxInt}`), a set larger than "the full integer width"; `IntDomain` `:241-256` returns nil for the same set. `guard` is unreachable from `table` (`go list -deps`), so the loader must re-state the rule; the record re-states a different one. And by Performance Expectations' own ~5.8 KB/row, "enumerable" is not the property that separates `{0..2^40}` (admitted) from `{MinInt..MaxInt}` (refused).
- why it matters: two authorities on "which int declaration has a readable width" would disagree on `{0..MaxInt}` — lint says no domain, the loader would enumerate. Define the loader's refusal as "width not representable (`intWidth`'s rule)", and let A7's no-ceiling stand as the explicit acceptance of huge-but-representable widths.

### F8 — C2 does not say which cell a multi-cell bound failure names
- anchor: 0030:C2
- class: silence
- claim: C2 "whose detail names the rule, the tag, the cell, and the stepped value"; Failure Modes F3 "A step whose magnitude exceeds the domain width refuses at every admitted cell under C2, naming the first".
- evidence: RDR is silent on the order "first" is taken in (domain order ascending for `int`, authored order for `enum`, or the row's suffix sort); C2 names one cell and F3 names "the first" without an order. S3's expectation (`cell 9, value 10`) is consistent with either reading for a single failing cell only.
- why it matters: the MVV asserts the refusal text; without a stated order the assertion and the implementation can each pick a different "first" cell for the oversize-step fixture in S5.

## Verdict basis

F1 and F2 are record-internal contradictions that would misdirect an implementer (wrong row set; an unsatisfiable MVV assertion); F3, F5 and F7 are codebase claims that do not hold as stated. None reopens the chosen approach: each is a C1/C2/S1 edit plus a symbol fix.
