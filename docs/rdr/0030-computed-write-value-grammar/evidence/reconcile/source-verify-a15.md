Model: claude-sonnet-5

# Source-search verification: RDR 0030 A15

Claim: a per-row graph-lint finding can carry the expanded row's suffixed
identity in `Rule`/`Element` without disturbing group-level findings, any
finding code, or any in-package consumer that reads those slots back.

## 1. PRODUCERS (per-row `Rule:` sites) — VERIFIED

Every site holds the `table.Row` value itself at construction (loop
variable `row`, or `rows[i]`/`rows[j]` from a cloned `[]table.Row`), so
`row.RuleID`/`row.Suffix` are reachable without threading new state:

- `internal/graphlint/groups.go:80` — `checkOwnedBeforeMatch`, `for _, row
  := range rows` (`rows := slices.Clone(g.Rows)`).
- `internal/graphlint/groups.go:207` — `emitOverlaps`, `for i := range
  rows`, `Rule: rows[i].RuleID`, `Element: right.RuleID` (also a `table.Row`
  loop var from `rows[i+1:]`).
- `internal/graphlint/groups.go:318` — `checkRedundantRows`, `Rule:
  rows[i].RuleID`, `Element: rows[j].RuleID`.
- `internal/graphlint/groups.go:363` — `checkIdempotentWrites`, `for _, row
  := range rows`.
- `internal/graphlint/groups.go:399` — `checkVacuousAtoms`, `for _, row :=
  range rows`.
- `internal/graphlint/analysis.go:213` — `checkSingleValuedState`, `for _,
  row := range a.model.Rows`.
- `internal/graphlint/analysis.go:567` — `checkUnreachableRules`, `for _,
  row := range rows`.
- `internal/graphlint/coverage.go:330` — `emitWithholdings`, `for _, row :=
  range rows`.

All eight sites (RDR's three files, groups.go alone has five) hold `Row`,
not a bare id string. Load-bearing half confirmed.

## 2. GROUP-LEVEL SITES — VERIFIED, DISJOINT

- `coverage.go:91` and `coverage.go:140` — `Rule: firstRuleID(g)`, over
  `guard.Group`, no `Row` in scope.
- `coverage.go:426` — `gap.Rule = firstRuleID(g)`.
- `coverage.go:439` — `Rule: closedBy`, where `closedBy` is set from
  `bareEscapeFor(g, class)` (coverage.go:557), which returns a bare
  `row.RuleID` string picked by sort — still bare, still group-scoped
  (names the distinguished escape row, not the row under test).

None of these four sites overlaps the eight per-row sites in (1); they
operate over `guard.Group`, not a per-row loop variable, and read
`ruleIDsOf(g)`/`firstRuleID(g)`, both bare. Disjoint confirmed.

## 3. READ-BACKS — PARTIAL (one unenumerated read-back found; see §7)

- `coverage.go:519` `groupHasOverlap(g guard.Group) bool`: `ids :=
  ruleIDsOf(g)` (bare authored ids); loops `a.findings` for `f.Code ==
  CodeOverlap && f.Class == ""`; test is `slices.Contains(ids, f.Rule) &&
  slices.Contains(ids, f.Element)`. This is a two-slot join of an emitted
  `graph-overlap` finding's `Rule`+`Element` (currently bare `RuleID`,
  written at groups.go:207) against the group's bare rule-id set, gating
  the `ambiguous_match` coverage arm at coverage.go:391-393. Exactly as
  the RDR describes. If `Rule`/`Element` become suffixed, this join breaks
  unless it recovers the bare id first (truncate at first `#`) — this is
  the site C3 must retarget.
- `engine.go:98-111` `identityKey(f clierr.Finding) string`: reads
  `f.Rule`/`f.Element` to build the sort key `sortFindings` (engine.go:88)
  orders on. NOT named in A15's evidence text or C3's clause. Analysis:
  benign — `TestReq88_RuleIDsSortBeforeGraphElementIDs`
  (findings_0006_test.go:901) tests only the RULE-namespace-before-
  ELEMENT-namespace bucketing (`f.Rule != ""` vs `f.Rule == "" &&
  f.Element != ""`), not literal string content; a suffixed non-empty
  `Rule` stays in the same namespace bucket, so REQ-88/114 ordering is
  unaffected. This read-back does not need to move; it is a genuine
  read-back A15's text omitted, but not a defect — see §7.
- `export.go:174` `compareEdges` reads `a.Rule`/`b.Rule` — FALSE POSITIVE
  on naming only: this is `Edge.Rule` (export.go:16-20), a distinct
  struct unrelated to `clierr.Finding`, part of the graph-edge vocabulary
  C3 explicitly holds fixed. Not a read-back of the finding payload.
- Test-only read-backs of `.Rule`/`.Element` are widespread (adversarial_,
  class_0010_, coverage_0006_, findings_0006_test.go, etc.), consistent
  with A15's own carve-out ("the package calls `Row.Identity()` nowhere
  outside tests" — tests are the expected exception, not silently
  omitted). `findings_0006_test.go:115`
  (`TestReq127_FindingWithoutARuleIDCarriesASpanOrElementID`) asserts
  `slices.Contains(authored, got.Rule)` against a bare-authored-id list,
  but its fixture rule (`advance`) is not `in`-expanded, so it is not
  exercised by A15's change and would not fail.

## 4. RECOVERY — VERIFIED, with one precision correction to A15's text

- `internal/table/model.go:337-347` `Row.Identity()` renders
  `ModelID + "." + RuleID + (suffixSep + s for s in Suffix)`, where
  `suffixSep = "#"` (model.go:25). Note: `Identity()` prepends
  `ModelID + "."` — the suffixed form C3 actually specifies for the
  finding's `Rule` slot is `rule#cell` (RDR line ~1032: "publishes the
  suffixed identity `rule#cell`"), i.e. `RuleID` + `#`-joined `Suffix`,
  built directly from those two fields — NOT `row.Identity()` verbatim
  (which would also carry the `ModelID.` prefix). This matches A15's own
  evidence ("the package calls `Row.Identity()` nowhere outside tests" —
  confirmed true both today and under the proposed change, since the
  producer sites would construct `rule#cell` directly, not call
  `Identity()`).
- Authored rule ids cannot contain `#`: `internal/table/normalize.go:297`
  inside `(*loader).normalizeRules`, checked immediately after reading
  `id := *rule.ID`: `if strings.Contains(id, suffixSep) { return
  fail(CatMalformedRuleID, ...) }`. Ban confirmed at load.
- `Suffix` members (the `in`-expansion cell labels) are also banned from
  containing `#`: `internal/table/load.go:163-165`
  (`case strings.Contains(o, suffixSep): fail(...)`), applied at the
  point candidate members are validated, before `normalize.go:800`
  (`next.Suffix = append(next.Suffix, member)`) ever appends one.
- Given both bans, `RuleID#s1#s2#...` truncated at the FIRST `#` always
  recovers the bare `RuleID`, for any number of suffix elements: neither
  the head (`RuleID`) nor any tail element can itself contain `#`.
  Recovery is sound.

## 5. `ruleIDsOf` — VERIFIED, needs no change

`coverage.go:574-584`: `for _, row := range g.Rows { ... row.RuleID ...
}` — reads `RuleID` directly, never `Identity()`. Produces bare authored
ids today and would continue to under A15's change, since `Row.RuleID`
itself is untouched by the proposal (only the finding payload's `Rule`/
`Element` slots change). No modification needed.

## 6. `Row.Identity()` callers outside tests — VERIFIED: none

`grep -rn "\.Identity()" internal/graphlint/*.go` returns exactly two
hits, both in `_test.go` files: `authority_0006_test.go:319`
(`row.Identity() == ""`) and `class_0010_test.go:776`
(`before.Rows[i].Identity()`). No production caller in
`internal/graphlint`.

## 7. Other in-package read-backs beyond `groupHasOverlap`

One unenumerated read-back found: `engine.go::identityKey`
(engine.go:98-111), called from `sortFindings` (engine.go:88), reads
both `f.Rule` and `f.Element`. Assessed BENIGN, not blocking: it only
uses non-emptiness of `f.Rule` to pick the sort-order namespace
(`namespace := "0"` if `Rule != ""`, else falls to `Element`, then
`Span`), then uses the raw string as a secondary sort key within that
namespace. A suffixed `Rule` value:
  - stays non-empty, so it stays in namespace "0" — REQ-88's
    rule-before-element guarantee (`TestReq88_RuleIDsSortBeforeGraphElementIDs`)
    is preserved.
  - changes the byte-order of findings sharing one (model, code) bucket
    only among themselves (e.g. `retry#1` sorts before `retry#2` — still
    a valid total order, just a different one than the unsuffixed form
    would give). No test pins the exact string content of this
    tie-break field, only that it forms a stable total order
    (`TestReq86And114_FindingsAreOrderedByTheFindingIdentityTuple`,
    findings_0006_test.go:730, which sorts and only checks monotonicity
    of `identityKey`, not its literal values).

This read-back was not named in A15's evidence text or in C3's normative
clause. It does not need to move (unlike `groupHasOverlap`, which must
retarget to a recovered bare id), but A15's claim that the producer
census — three files, N sites — is "the half that misses" the
consumer census is itself proven true here: a second consumer
(`identityKey`) exists beyond the one A15 explicitly names
(`groupHasOverlap`). It happens to be harmless, but its existence was
not enumerated, which is exactly the gap-shape A15 warns about. Recommend
C3's text name `identityKey` explicitly (as "safe, verified") rather than
leave it undiscovered.

## Verdict rationale

All eight sub-claims individually verified/refuted as detailed above.
The overall A15 claim is essentially sound — the one join that must move
(`groupHasOverlap`) is correctly identified and its fix is well-specified
(recover via truncate-at-first-`#`, sound per §4). The one gap is that
A15's own text under-enumerates: it says "every READ-BACK" but only
names `groupHasOverlap`; `identityKey` is a second one, undiscovered by
A15's own search, that happens not to break anything. Since A15 exists
specifically to catch "any in-package consumer that reads those slots
back," and its own census stopped one short, this is a PARTIAL on the
completeness of A15's enumeration (not a defect in the underlying code
or in C3's fix), worth a one-line addition to A15/C3 before lock.
