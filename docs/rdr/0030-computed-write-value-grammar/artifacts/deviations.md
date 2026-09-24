# Deviations — RDR 0030 Counting and stepping a tag without writing out every value by hand

Phase 2 implementation deviations. Append-only.

---

## D1 — REQ-3's literal negative control wrote a value outside its own domain

- Type: TEST-FIXTURE
- Where: `internal/table/step_grammar_0030_test.go::TestReq3_0030_AnUnrepresentableIntWidthIsRefusedForAStep`
- Evidence: the "max - min == MaxInt ending at -1" arm declares `{MinInt..-1}` and its
  "same declaration written literally" twin wrote `x = 0`, which the pre-existing
  `conformDomain` refuses (`0 is above max -1`) on the base before any step exists. The
  twin could never load, so the arm could not discriminate the step's refusal.
- Change: the twin writes the declared `min` (`tc.minV`), a member of every arm's domain.
  The refusal assertion is unchanged.
- Status: mechanical translation

## D2 — REQ-52's `#`-domain negative control dropped a member the fixture still guards on

- Type: TEST-FIXTURE
- Where: `internal/cli/lint_mvv_0030_test.go::TestReq52_0030_EveryStepGrammarAndKindDefectIsANamedLoadRefusal`, case "# in the stepped domain"
- Evidence: replacing the ladder's `tier` domain with `["small", "mid", "la#rge"]` removes
  `large`, which `ladder-step.toml`'s `escalate-cap` names in `guard.all tier eq`, so the
  unstepped control refused `malformed_predicate_atom` (`"large" is outside the declared
  domain`) rather than loading.
- Change: both the refusal and its control use `["small", "mid", "large", "la#rge"]`; the
  `#` member still sits in the stepped domain, so the refusal is the step-scoped `#` guard.
- Status: mechanical translation

## D3 — The relocated seam's constructor and the shared comparison are exported from `internal/resolve`

- Type: IMPL-DECISION
- Where: `internal/resolve/evaluator.go` (`Evaluator`, `NewEvaluator`, `CompareValue`, `IntWidth`)
- Evidence: IP Phase 2 moves `Evaluator` to `internal/resolve` and extracts the
  render-then-evaluate core there; JDR 0004 JD-3 has the loader construct over its own
  declarations; REQ-14 asserts an exported constructor returning the concrete evaluator and
  REQ-40 asserts one exported comparison all three callers reach. `IntWidth` is named by C1.
  `guard.Evaluator` is a type alias and `guard.NewEvaluator` delegates, so 0012's pins keep
  their names. The canonicalizing set renderer stays unexported inside `CompareValue`.
- Status: mechanical translation

## D4 — An `exists` atom on a stepped tag is decided from presence during admission

- Type: IMPL-DECISION
- Where: `internal/table/step.go::admits`
- Evidence: C1 evaluates every positive atom "as the runtime evaluator would"; `exists` never
  reaches the value seam (it answers unevaluable there) and the kernel decides it from key
  presence. A cell is a present value, so `exists = true` admits and `exists = false`
  excludes, as the runtime would, instead of excluding every cell via the seam's
  unevaluable answer.
- Status: mechanical translation

## D5 — Per-row lint attribution changes on models that step nothing (Q1 reading (a))

- Type: SPEC-DEFECT
- Where: `internal/graphlint/{analysis,coverage,groups}.go` per-row finding sites;
  `internal/graphlint/engine.go::rowIdentity`/`authoredID`
- Evidence: Q1 resolved to reading (a) (req-list.md Q1 ASSUMPTION): every row carrying a
  suffix publishes `rule#suffix…` in `Rule`/`Element`, `in`-minted rows included. That
  contradicts G/C1's "this record changes no behaviour on a model that steps nothing" for
  the attribution field alone: an `in`-expanded row's per-row finding now reads
  `fan#small`, not `fan`. `groupHasOverlap` joins on the recovered authored id.
- Status: mechanical translation

## D6 — 0012's zero-value check keyed the Evaluator's named type to `internal/guard`

- Type: DEPENDENCY-LIMIT
- Where: `internal/guard/seam_carrier_0012_test.go::r12IsEvaluator`, `::r12AllowedZeroSite`
  (drives `TestReq48_0012_NoZeroValueEvaluatorSurvivesOutsideTheOneAllowListedFunction` and
  `guard_adv_0012_test.go::TestAdv1_0012_TheZeroValueCheckCatchesEveryZeroValueForm`)
- Evidence: after the relocation (REQ-40) `guard.Evaluator` is an alias of
  `resolve.Evaluator`, so `types.Unalias` resolves every site to a named type whose package is
  `internal/resolve`; the check matched only `internal/guard` and stopped finding
  `atomAdmitsValue`'s still-present `guard.Evaluator{}`. This is the "zero-value allow-list keyed
  on the `internal/guard/` prefix" the req-list ASSUMPTION names, migrated on the precedent of
  0012 deviations D3.
- Change: the named type is matched at either home, and `NewEvaluator`'s own body is allowed in
  `internal/resolve/` as in `internal/guard/`. The ONE allow-listed function and every plant in
  the adversarial suite are unchanged and still caught.
- Status: mechanical translation

## D7 — 0001's no-canonical-serialization pin scanned the whole `internal/resolve` package

- Type: SPEC-DEFECT
- Where: `internal/resolve/resolve_test.go::TestReq37_KernelIntroducesNoHashOrCanonicalSerialization`
- Evidence: REQ-40 moves the `Evaluator` (which decodes JDR 0001 §D13's canonical JSON set
  literal) and the canonicalizing set renderer into `internal/resolve`, so the package now imports
  `encoding/json`, which the 0001 pin bans package-wide. 0030's `Overrides` field names no 0001
  clause. JDR 0004 JD-1 (binds 0030) settles the scope: the kernel's fences govern its RESOLUTION
  PATH, not the package as a namespace, and option (c), a fourth package to keep `resolve`
  literally untouched, was rejected. The serialization is §D13's, already decided at the kernel
  seam, not one RDR 0001 introduces.
- Change: the import half of the pin scans the package minus the relocated seam file
  (`evaluator.go`); the exported-symbol half is unchanged and still passes.
- Status: mechanical translation (derived choice, JDR 0004 JD-1)

## D8 — graph-lint's private canonical set renderer is retired

- Type: IMPL-DECISION
- Where: `internal/graphlint/reach.go` (`renderSetLiteral` removed)
- Evidence: `atomAdmitsValue`'s fall-through now renders through `resolve.CompareValue`'s
  canonicalizing renderer (REQ-40, "rather than minting a fourth copy"), leaving the private copy
  unused (`golangci-lint` `unused`). The rendered form is identical (sorted, compacted, `[]` for
  empty).
- Status: mechanical translation

## D9 — The step's int walk was sized by the declared width, not the admitted cells

- Type: IMPL-GAP
- Where: `internal/table/step.go::stepCells`, `::resolveStep` (verification FAIL-1, ADV-1)
- Evidence: C1 admits a representable width "and merely expensive … the record takes no
  load-time row ceiling" (REQ-6); A7 bounds the rows by the admitted cells. `stepCells`
  pre-sized `make([]string, 0, width)`, so `{0..MaxInt-1}` panicked (`makeslice: cap out of
  range`) before any atom was read, and a walk over the full width never returns.
- Change: `stepCells` yields candidates lazily and, for `int`, narrows them first by the
  positive atoms: the ordered bounds (`lt`/`lte`/`gt`/`gte`) clip `min..max`, and the smallest
  `eq`/`in` member set pins the candidates. It drops only members some positive atom answers
  false at, and every candidate still goes through `admits`, so the admitted set, the domain
  order and C2's first-failing-cell are unchanged. No ceiling is taken; a width no atom narrows
  is still walked in full.
- Status: mechanical translation

## D10 — ADV-1's fixture bound no reader to its owned tags

- Type: TEST-FIXTURE
- Where: `internal/table/adversarial_0030_test.go::advStepModel`
- Evidence: once the panic was gone the child load refused
  `malformed_accessor_binding: owned tag attempt is served by 0 readers`, a pre-existing rule
  the panic had masked. ADV-3 shares the helper and refuses earlier, on the step.
- Change: the helper declares the `read.state` / `write.state` pair over `status` and
  `attempt`, as `internal/graphlint/adversarial_0030_test.go::adv2Model` does. The assertions
  are unchanged.
- Status: mechanical translation

## D11 — The non-representable width refusal did not name the admitted kinds

- Type: IMPL-GAP
- Where: `internal/table/step.go::parseStep` (verification ADV-3)
- Evidence: F1 has the width arm refuse "the same way" as the kind arm, and S5 expects each
  refusal's detail to name "the admitted form or kinds". The width detail named neither.
- Change: the width detail keeps its bound and cause and appends the kind refusal's text (one
  constant, shared by both arms).
- Status: mechanical translation

## D12 — `graph-coverage-closed-by-escape` named an expanded escape row by its bare id

- Type: IMPL-GAP
- Where: `internal/graphlint/coverage.go::bareEscapeFor` (verification FAIL-2)
- Evidence: REQ-36/REQ-37 under Q1 reading (a): every surface that names a ROW publishes
  `rule#suffix…`. The finding's `Rule` and message name the closing ROW, not a `firstRuleID`
  group site, yet read `esc` for the row `esc#retry`.
- Change: `bareEscapeFor` collects `rowIdentity(row)`. The guard-side lint (`internal/guard/lint.go`)
  keeps bare ids on every row, as Phase 2 left all its per-row findings; REQ-37 scopes graph-lint.
- Status: mechanical translation

## D13 — A step over a subsumed match `in` lints a coverage gap on every claimed cell

- Type: SPEC-DEFECT
- Where: `internal/table/normalize.go` step expansion (REQ-13 row shape) meeting
  `internal/guard/product.go::Dimensions`/`::Product` (the scoped product); verification ADV-2;
  test `internal/graphlint/adversarial_0030_test.go::TestAdv2_0030_SubsumedInLadderLintsNoCoverageGap`
- Evidence: C1 (REQ-13) mandates each row of a step over a subsumed `in` carry BOTH the
  rewritten `match eq = <cell>` and `guard.all eq = <cell>`. The match `eq` puts each row in
  its own group; `0006:C7` has the scoped product range over the group's guard dimensions'
  DECLARED domains, with match keys contributing "no assignment to either side", so `attempt`
  spans `{0..9}` in every per-cell group and 9 of 10 read uncovered: 5 blocking
  `graph-coverage-gap` on a ladder whose every cell is claimed (F4 names the gap as the
  UNCLAIMED cell's signal). The hand-unrolled `match eq` ladder lints zero gaps; a literal
  ladder carrying both atoms per row lints the same 5 gaps (probed on the built CLI), so the
  expansion matches its literal twin (REQ-25) and the defect is the mandated shape itself.
  Neither fix is available without re-deciding a clause: dropping the `guard.all eq` where an
  `in` was subsumed contradicts REQ-13's "exactly two atoms"; narrowing a guard dimension by
  the group's own match `eq` contradicts `0006:C7`.
- Change (b9c340f): `internal/table/normalize.go::expand` marks a step point pinned when a
  match `eq` on its key scopes the row (the subsumed-`in` rewrite, or an authored/inherited
  match `eq`); it always keeps the rewritten `match eq = <cell>` and appends
  `guard.all eq = <cell>` only when not pinned. `TestReq13_0030` re-cut; ADV-2 green unchanged.
- Status: needs author decision → RESOLVED (3d strong consult: option (a) widened — emit the expansion's `guard.all eq = <cell>` only when no match `eq` pins the stepped tag (a subsumed-`in` rewrite, or an authored/inherited match `eq` on that key); REQ-13 "exactly two" and REQ-16 read as "the per-cell atom lands in the block the author scoped the tag with"; cite JDR 0001 §D6 (docs/jdr/0001-resolve-kernel-seam.md:260-292), 0030 lines 651-657 and C1 "exactly as `expand`'s existing `case expanding:` arm rewrites an `in`" (normalize.go:917-925 mints match `eq`, no guard atom); (b) rejected — re-decides 0003:C13/0006:C7 and breaks 0030 "changes no behaviour on a model that steps nothing"; (c) rejected — contradicts F4)

## D14 — A step could land on a declared `<clear>` domain member

- Type: IMPL-GAP
- Where: `internal/table/step.go::resolveStep` (triage row 3)
- Evidence: C2 has the stepped literal conform "exactly as an authored literal write does", and
  0002:C11 refuses an authored `<clear>` write as `reserved_tag_value`. The declaration never
  inspects domain members for the sentinel, so `domain = ["small", "<clear>"]` stepping `+1` off
  `small` loaded a row writing `<clear>`, which resolve and flowbind read as removal.
- Change: after `stepValue`, a result equal to `ClearSentinel` refuses under `reserved_tag_value`
  (the literal write's category, 0002:C11's existing refusal at a new site, not one REQ-10's
  minted set names), detail naming the rule, the tag and the cell. The source-cell arm is
  unreachable (`eq`/`in` naming the sentinel already refuse).
  Test: `internal/table/step_expand_0030_test.go::TestReq27_0030_AStepLandingOnTheClearSentinelIsRefusedAsTheLiteralIs`.
- Status: mechanical translation
