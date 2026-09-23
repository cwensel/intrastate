# deviations — 0012 declared-kind carrier at the guard value seam

## D1 — a nil-mapping evaluator cannot byte-compare, so the match site does it

Type: SPEC-DEFECT
Status: mechanical translation (derived choice: the match path's `eq`/`in` byte comparison is made inside `internal/graphlint/reach.go::atomAdmitsValue`; every other operator still goes through its zero-value `guard.Evaluator{}`)

The record asks for three things that cannot all hold at the same time:

- `0012:C1` (REQ-7, REQ-49): a nil-mapping evaluator answers
  GuardUnevaluable for every `eq`/`in` atom and never falls back to
  raw-string comparison.
- `0012:C4` (REQ-35): `atomAdmitsValue` "constructs `guard.Evaluator{}` with
  no mapping and compares raw bytes. That is required, not tolerated", and
  it keeps its `main` shape.
- `0012:S6` (REQ-50): no committed model's lint verdict flips.

`atomAdmitsValue` reads the verdict as `verdict != resolve.GuardFalse`. Once
C1 lands, its zero-value evaluator answers every match `eq`/`in`
unevaluable, and that reads as "admissible". Every owned match edge then
becomes traversable, and on the pre-change golden REQ-50 flipped 15
committed models. Each lost its `graph-dead-end` and
`graph-terminal-escape` findings (for example
`internal/table/testdata/merge-distinct.toml`, node
`{iter=0, stage=propose, status=Draft}`). `0012:S4` calls this folding
"harmless" on the match path, but that holds only while the zero-value
evaluator still byte-compares, and C1 removes that behaviour.

The record states what it wants in several places. C4 says the site
"compares raw bytes" because the kernel byte-compares at
`resolve.go::TagSet.matches`, and that typing the site "would create the
divergence this record closes". S6 says no verdict flips. The derived
choice meets all of them. `atomAdmitsValue` now compares `eq`
(`strings.Join(literal) == held`) and `in` (`slices.Contains`) byte for
byte, which is exactly what the zero-value seam used to answer, including
an empty `in` list, which reads as admissible. It keeps its name, its home,
its `bool` return and its `guard.Evaluator{}` construction for every other
operator, and it reads no declarations. REQ-35, REQ-48 (the allow-list site
is still a zero-value construction) and REQ-50 all pass. No surface was
added.

## D2 — the published conformance fixture's key names

Type: IMPL-DECISION
Status: mechanical translation

`0012:C3` fixes `ContractKinds()` as a function that returns a fresh
`map[string]string` with one key per kind token. It does not name the keys.
The implementation publishes `phase`→`enum`, `label`→`scalar`,
`count`→`int`, `flag`→`bool` and `labels`→`set`. The keying follows REQ-29:

- legacy `eq` cases → `phase` (enum)
- legacy `in` cases → `label` (scalar)
- `gte` cases → `count`
- `contains` cases → `labels`
- new int, bool and set legs → `count`, `flag` and `labels`
- the absent-key leg → `undeclared`
- the unknown-token leg → `measure`→`decimal`, built over a caller-local
  copy so the published fixture stays exactly the five tokens (the
  req-list ASSUMPTION)

The suite also carries the REQ-26 meta-check itself. Recorded because
`ContractKinds` is a cross-RDR surface that later implementers key their
own cases against.

## D3 — predecessor tests construct their seam over declarations

Type: IMPL-DECISION
Status: mechanical translation

Under C1, the RDR 0003, 0005, 0009 and 0011 tests that handed the kernel a
zero-value evaluator went red. Their `eq`/`in` atoms now answer
unevaluable. The record names this migration: IP Phase 1 says "internal/guard's own test
migration", and §Consequences counts "eighteen zero-value construction
sites" in the test surface. Each site now constructs the way a production
site does. No assertion was changed:

- `guard_fixtures_0003_test.go::resolveWith` takes the model its table came
  from and builds `guard.NewEvaluator(guard.DeclaredKinds(m))`. Its 34
  callers pass the model they already hold.
- Kernel-only tables with no model (`guard_evaluator_0003_test.go`
  REQ-35/36/37, `guard_narrowing_0003_test.go` REQ-41/103, and the direct
  `Evaluate` checks in `guard_grammar_0003_test.go` REQ-1/12 and
  `guard_scope_0003_test.go` REQ-119/120) construct over a literal mapping
  that declares each key the atom names.
- The `internal/cli` escape-shape tests call `guardSeam(nil)`. Their tables
  carry no guard atoms. `flow_next_0011_test.go` passes `probeRow` the seam
  `guardSeam(model)` builds.

The migrations the record itself prescribes landed as written: the A2
narrowing of `TestReq34`'s field check, the one-line adapter at the
`guard_evaluator_0003_test.go` suite call, and the re-typed `var fn` pin
plus the `conformingContractSeam` constructor in `guard_mvv_test.go`.
`conformingContractSeam` compares `eq`/`in` under its mapping's kind, with
an `strconv.Atoi` parse so that the suite's overflow leg is unevaluable.
The file's own `parseInt` wraps on overflow.

## D4 — the S4 zero-value check matched the spelling, not the type

Type: IMPL-GAP
Status: mechanical translation (fixed 6c113a5)

`0012:S4` (REQ-48) asks for "a typed check, not a text grep", one that
fails on any composite literal, `var`, `new` or struct field "of type
`guard.Evaluator`". The Phase 2 check parsed files with `go/parser` and
matched the selector `<alias>.Evaluator`. Phase 3a (FAIL-48) and 3b (ADV-1)
each planted forms that passed it: an alias-typed, dot-imported or elided
composite literal, a named result, an array `var`, and a `make`-allocated
element.

`internal/guard/seam_carrier_0012_test.go::r12ZeroValueSites` now
type-checks every main-module package's non-test files with `go/types`,
importing dependencies from `go list -deps -export` data (no new module
dependency). It matches `guard.Evaluator` by its type object through any
alias. A zero-value `var` includes a named result, and an array of the
type counts as holding it, as does `new`/`make([]…)` over it. The allow-list
is unchanged: `NewEvaluator`'s body and `reach.go::atomAdmitsValue`, by
name. The check sees the files `go list` selects for the host platform, so
a file behind another GOOS/GOARCH constraint is not type-checked. The one
such non-test file today, `internal/cli/cmdbind/procgroup_other.go`, does
not reference `internal/guard`. `TestAdv1_0012_…` carries every Phase 3 form as a
plant.

## D5 — the S4 check missed array literals that zero-fill an element

Type: IMPL-GAP
Status: mechanical translation (fixed 31ba7bb)

Triage row 2 (roborev 8247). D4 already reads "holds `guard.Evaluator`" as
including an array of it, but the composite-literal arm of
`internal/guard/seam_carrier_0012_test.go::r12FileZeroValueSites` matched
only a literal whose own type is `guard.Evaluator`. An array literal with
fewer elements than its length (`[1]guard.Evaluator{}[0]`,
`[2]guard.Evaluator{guard.NewEvaluator(nil)}[1]`, `[...]T{2: v}`) zero-fills
the omitted indices with no nested literal for the walk to find. The arm now
also records "array literal with omitted elements" when the literal's
underlying type is an array whose element holds an Evaluator and
`len(Elts) < Len()`. Go rejects duplicate constant indices, so that
inequality is exactly "some index omitted". `TestAdv1_0012_…` carries both
array-literal forms as plants.
