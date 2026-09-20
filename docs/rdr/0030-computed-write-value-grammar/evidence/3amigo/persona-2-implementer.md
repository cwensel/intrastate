Model: claude-fable-5-1

# Persona 2 — Implementer

Question held: if I started coding this Monday, what would I ask in the
first hour? Starting set: 0030:C1, 0030:C2, 0030:C3, 0030:D-identity,
0030:D-wire-byte-format, 0030:D-naming, 0030:D-selection-predicate, plus
the `source-anchor` edges. Widened where noted; each widening names what
sent me. Source consulted: `internal/table/normalize.go` (`expand`,
`renderWrites`, `compareAtoms`, `mergeAtoms`), `internal/table/load.go`
(`valueMembers`, `conform`, `conformDomain`, `loadInitial`),
`internal/table/model.go`, `internal/table/source.go`,
`internal/guard/declaration.go`, `internal/guard/grammar.go`,
`internal/guard/product.go`, `internal/graphlint/reach.go`,
`internal/graphlint/groups.go`, `internal/resolve/resolve.go`.

No high finding. Every item below is a clarification request; none
refutes the design.

- **0030:C1** — medium — C1 states "the stepped tag carries exactly two
  atoms per row: the rewritten `match eq = <cell>` … and the expansion's
  own `guard.all eq = <cell>`", then in the next sentence retains "every
  OTHER authored atom on the tag — the `guard.all` bounds". The
  §mini-check-tables `trace` row 4 witness shows THREE atoms on the tag
  (`{attempt match eq <cell>, attempt all lt 5, attempt all eq <cell>}`).
  "Exactly two" reads as a countable invariant an implementer would assert
  (`len(atomsOnTag) == 2`) and it is false whenever a bound is retained.
  Ask: reword to "exactly two atoms FROM this mechanism" (or "at most one
  rewritten match `eq` plus one emitted `guard.all eq`"). — blocks: the
  per-row atom emission in Phase 2 and any S-test asserting the stepped
  tag's atom set.

- **0030:C2** — medium — the bound detail "names the rule, the tag, the
  cell, and the stepped value". For an `enum` step that leaves the domain
  there IS no stepped value — no member sits `n` positions from the cell.
  Widened to 0030:S3 (sent by that phrase): S3's enum sibling quietly
  names only `tier` and `large`, no value. Ask: what does the enum detail
  print in the value slot — the position offset (`large +1`), the bare
  `n`, or nothing? — blocks: the C2 detail format, and the enum half of
  S3's / MVV row 4's assertion.

- **0030:D-identity** — medium — "the cell appended as one more suffix
  element in the choice-point sort order" and "two step points in one
  rule order their suffix elements by the choice-point sort, as `in` atoms
  do". A step point is not an `Atom`, so `compareAtoms`' tuple (key,
  block, operator, literal) does not apply to it as written; widened to
  0030:A1, which says only that it "must slot into" that sort. After
  subsumption at most one suffix-emitting point exists per key, so the
  only consequential order is BY KEY — but that reasoning is mine, not the
  record's. Ask: state the tuple (e.g. `(key, BlockAll, "eq", ∅)`) or
  state "by key, total because subsumption leaves one suffix-emitting
  point per key". — blocks: suffix order for a rule with two step points
  or a step point beside an `in` on another key; hence `compareRows`,
  the dump order, and every `rule#a#b` identity in the fixtures.

- **0030:§prerequisites** — medium — widened from 0030:A4 / 0030:A12
  (sent by "0012's typed evaluator, `0012:C2`") and the record's
  Predecessors edge to 0012. RDR 0012 is `Status: Draft`; today's
  `internal/guard/grammar.go::Evaluator` is `struct{}` with string `eq`,
  and `0012:C2` would make `eq`/`in` typed against a kind mapping. The
  Prerequisites name only 0002 as Implemented. Ask: does 0030 land on
  today's untyped evaluator (so the move to `resolve` precedes 0012, and
  0012's edits then target `resolve`), or does it wait? The answer fixes
  the shim's signature — whether the loader's admits call passes the
  declared kind. — blocks: Phase 2's start and the extracted shim's
  signature.

- **0030:C2** — low — "computed without overflow at every admitted cell,
  whatever the step's magnitude" plus §phase-2's "the bound refusal
  through `conform`". `conformDomain`'s `int` arm calls `strconv.Atoi`
  on the rendered literal; an exactly-computed stepped value beyond
  int64 renders a decimal `Atoi` rejects, so the detail reads "is not an
  int" rather than "is above max". Ask: either pre-check (`|n| >= width`
  implies the first domain-order cell fails, so refuse before rendering)
  or say the arithmetic is checked and the detail form for that arm. —
  blocks: the arithmetic representation (checked add vs `math/big`) and
  S3's message assertion.

- **0030:C1** — low — `intWidth`'s rule is RESTATED in `table` (a copy)
  because `guard` is unreachable from `table`, while the very same cycle
  is solved for `Evaluator` by a MOVE into `resolve`. The `{min..max}`
  enumeration and the decimal cell rendering are copied too, and must
  byte-match `guard/assignment.go`'s `strconv.Itoa` cells or the emitted
  `guard.all eq = <cell>` never satisfies at lint. Ask: one sentence on
  why `intWidth`/`IntDomain` are copied rather than moved down (or moved
  to `resolve` beside the shim), so the implementer does not re-decide. —
  blocks: whether to duplicate or relocate `intWidth`/`IntDomain`.

- **0030:C1** — low — silence: a rule that both steps a key and lists it
  in `clear`. Today `renderWrites` lets the clear list overwrite
  `assignments[key]` with `<clear>` silently (no refusal exists for
  write+clear on one key). With a step, the choice point would either
  expand N rows all writing `<clear>` or the step would vanish. Ask:
  refuse the pair under `malformed_tag_declaration`, or define the
  overwrite. — blocks: whether the step spec is derived before or after
  the clear-list pass in `renderWrites`.

- **0030:§phase-2-expand-into-cells** — low — "relocates unchanged" and
  §decision-rationale's "two callers repointed" understate the move: the
  non-test repoints are `graphlint/reach.go::atomAdmitsValue`,
  `guard/product.go::valueSatisfies` and `cli/flow_resolve.go::guardSeam`,
  plus roughly twenty `guard.Evaluator` references across
  `internal/guard/*_test.go` (including a field-count reflection test in
  `guard_evaluator_0003_test.go`). Ask: confirm the test files move or are
  repointed in the same phase and that the reflection test's package
  assumption is accepted scope. — blocks: the Phase 2 scope estimate,
  nothing in the design.

Not findings, recorded so the next reader need not re-check: `[initial]`
and predicate literals already refuse a TOML inline table via
`valueMembers`' default arm under their own categories (MVV row 5 pins
existing behaviour); `sourceRule.Write` is `*map[string]any`, so
`{ step = n }` decodes without a struct change; `flow set-state` exists
(`internal/cli/flow_input.go`); `docs/model-authoring.md` has a write-block
section; `checkIdempotentWrites` compares a match `eq` literal to the
write string and cannot fire on a non-zero step.
