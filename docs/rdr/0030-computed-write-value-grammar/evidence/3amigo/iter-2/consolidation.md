Model: claude-opus-5[1m]

# 3amigo consolidation — cli/0030, iteration 2 (delta)

Delta re-run after iteration 1's resolve, ordered by `rdr-loop.toml`
(`fix=substantial` → `rerun`: a rewrite can open gaps). Three fresh
isolated persona contexts, each scoped to the passages iteration 1
rewrote, each stamped `claude-fable-5-1`. The re-run earned its keep: it
caught four overclaims the iteration-1 fixes introduced.

## Hotspots (>=2 isolated personas)

| id | personas | what converged |
| --- | --- | --- |
| `0030:C2` | PM, Impl, QA (3) | the "structured `rule`/`key`/`literal` fields" claim — invented; `Failure` has no such fields |
| `0030:C1` | Impl, QA, PM (3) | `IntDomain` cannot move (takes `table.TagDecl`); step-spec carrier unstated |
| `0030:C3` | all three (3) | guide scope; whether an envelope field addition would break no-new-vocabulary |
| `0030:A13` / `0030:A14` | all three (3) | new elements: A13 miscited its guard row and had no fixture |
| `0030:S1` | PM, QA (2) | the both-tiers-in-one-run oracle is unsatisfiable against shipped `lint.go` |
| `0030:A12`, `0030:D-identity`, `0030:MVV`, `0030:§prerequisites` | 2 each | move scope, candidate interleaving, count ambiguity, 0012 sequencing |

## Findings introduced BY iteration 1's fixes (the re-run's real yield)

- **`0030:C2` structured fields — high, 3 personas.** Iteration 1 wrote
  "the detail's structured `rule`, `key` and `literal` fields are
  populated on both arms". Grounded REFUTED:
  `internal/table/category.go::Failure` carries `Category`, `Detail`,
  `Offending`, `Rule` and nothing else; `fail(cat, detail string)` takes a
  string. Worse, `Failure.Rule` is a direction identifier
  (`RuleKernelOwned`/`RuleAuthorMustRename`) that
  `internal/cli/flow_input.go` gates a rename hint on, so populating it
  with a rule id would publish that hint on the wire; and minting new
  fields would contradict C3's no-new-vocabulary clause. FIXED: retracted
  to a `Detail` string that names the cell first, with `Category` as the
  structured half of the assertion. S3 amended to match.
- **`0030:S1` both-tiers oracle — high, 2 personas.** Iteration 1's
  vacuous-pass guard demanded one blocking-tier AND one advisory-tier
  finding from one `lint --as json` run. Grounded REFUTED:
  `internal/cli/lint.go` exits through `respond.Fail` with
  `report.Blocking()` alone on any blocking finding, and builds the
  advisory list only on the success path — never both. FIXED: split into
  a clean pair (advisory tier, reused by every other scenario) and a
  dedicated blocking pair reused by nothing, so the deliberate defect
  cannot perturb MVV item 7's rule count.
- **`0030:C1` / `0030:A12` `IntDomain` move — high, 1 persona.** Iteration
  1 wrote that `intWidth` and `IntDomain` both move to `internal/resolve`.
  Grounded PARTIAL: `intWidth(minV, maxV int)` is type-free and moves;
  `IntDomain(d table.TagDecl)` cannot, since `resolve` must not import
  `table`. FIXED: only the type-free core moves; `IntDomain` stays as the
  decl-unpacking wrapper it already is. The same cycle bites
  `valueSatisfies`, now stated.
- **`0030:§consequences` remedy — medium, 1 persona.** Iteration 1
  rewrote C2 to say a literal row cannot clear the refusal, but left
  Consequences bullet 2 offering it as an alternative. FIXED (this is
  §amendment-sweep catching its own miss).

## Findings on pre-existing text

- **`0030:S5b` zero-overlap is vacuous — medium.** The `unless`-carrying
  fixture cannot exercise the overlap oracle where
  `groups.go::decidableAccepted` falls back to the scoped-product path (an
  `unless` row returns undecided and is skipped). Nuance the persona
  overstated: a row that projects whole through
  `guard.AcceptedAssignments` IS decided, `unless` included. FIXED: the
  projectability precondition is now explicit, and A5's overlap control
  moved to a new 5c whose exclusion is a positive atom pair.
- **`0030:S2` unclaimed-cell arm has no instance — medium.** S1 asserts
  zero `graph-coverage-gap` on the pair, so no unclaimed cell exists
  there. FIXED: moved to a new 2b on the F4 cap-cell shape.
- **`0030:A13` had no fixture — medium.** Method said MVV Test; no
  scenario declared an overflowing-step fixture (S5's item is an
  unrepresentable WIDTH, a declaration defect). FIXED: added to S5, with
  the expected detail pinned to "cell and bound", not "is not an int".
- **`0030:S5` negative width — medium.** `load.go::tagDecl` already
  refuses `min exceeds max` today, so that arm is not a new refusal and
  has no available control. FIXED: scoped to the
  `max - min == math.MaxInt` arm and stated as a restatement.
- **`0030:C1` step-spec carrier — medium.** `renderWrites` returns
  `([]TagValue, []string, error)` and `expand` takes `writes []TagValue`;
  neither carries a spec, and `RequiresOwned` derives from
  `maps.Keys(assignments)`, so a spec tracked beside the map would drop
  the stepped key from the owned-state gate. FIXED: the stepped key goes
  in `assignments`; the signature widening is named.
- **`0030:D-identity` interleaving — medium.** Order-by-key fixed the
  order among step points but not against `in` points on other keys.
  FIXED: the same key comparison places the step point in the whole
  candidate list; no tie is possible.
- **Lows, all FIXED**: S6 cited `internal/cli/tier_assertions_0029_test.go`
  (the file is under `internal/graphlint/`); S8 named an internal Go API
  as its oracle against the section's CLI-level rule (now `dump`'s `atoms`
  column); S3c asserted admitted rows on a load that refuses (now two
  fixtures); MVV item 7's approximate figures (now a strict inequality
  derived from the fixtures); A13's miscited guard row; C1's `Itoa`
  attributed to `IntDomain` rather than each caller; the enum residual
  added to C3's guide hazards as the third.

## Judged and holding (iteration 1 fixes the re-run confirmed)

C1's step-scoped `#`/duplicate guards; C1's clear-list refusal against
A14 (`renderWrites`' write loop precedes the clear loop, silent overwrite
confirmed); C2's domain-order-first and checked-add-before-render against
A13 (`conform` runs `conformKind`'s `Atoi` first — the ordering is
required); C3's row-vs-rule identity split; D-identity's order-by-key
against C1's one-point-per-key; §prerequisites' no-wait-on-0012 (0012 is
Draft, `Evaluator` is `struct{}`, `grammar.go` imports only stdlib +
`resolve`); phase-2's test-reference count (exactly 25 across 7 files);
S4, S9, S10 all judged writable with no finding.

No CONTRADICTION row and no net-new scope: every entry above traces to
the iteration-1 ledger or to a passage iteration 1 rewrote.
