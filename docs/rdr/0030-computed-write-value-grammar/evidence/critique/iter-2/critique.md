Model: claude-fable-5

# Critique — delta re-run (iter-2)

Delta re-run over the passages the first critique pass touched, plus the
additions the rewrite made. Scope: `0030:A2`, `A4`, `A5`, `A7`, `A8`,
`A10`, `A12`, `A14`, `A15`, `C1`, `C2`, `C3`, `D-identity`,
`D-selection-predicate`, `S1`, `S9`, `S11`, `MVV` (item 2),
`§consequences`, `§decision-rationale`, `§performance-expectations`,
`§problem-statement`, `§proportionality`, `§illustrative-code`,
`§finalization-gate`. The record is read as it stands; no change-history
commentary. Claims were checked against `internal/graphlint/{groups,
coverage,analysis,taxonomy}.go`, `internal/table/{normalize,category}.go`,
`internal/guard/product.go`, `internal/cli/clierr/clierr.go`, and the A7
spike under `evidence/spikes/a7-row-growth.md`.

## Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
| --- | --- | --- | --- | --- |
| E-1 | `0030:C3` (lint paragraph), `0030:A15` | Suffixing per-row `Rule` breaks an in-package JOIN that reads it back by bare id: `coverage.go::groupHasOverlap` matches the EMITTED `graph-overlap` finding's `Rule`/`Element` against `ruleIDsOf(g)` (bare `row.RuleID`s) to decide whether the `ambiguous_match` coverage arm is reachable. An overlap whose left row is an expanded row publishes `retry#3`, never matches, and the arm is "treated as vacuously closed". A15's verification plan (per-row sites hold the Row; `firstRuleID` sites disjoint) passes while this ships. | `lint` reports the blocking `graph-overlap` on `retry#3` but then SKIPS the `ambiguous_match` coverage arm for that group — no `graph-coverage-gap` demanding the rescue row the overlap makes reachable. A false green on exactly the arm the overlap unlocks; the model runs into `ambiguous_match` at runtime with no escape row. | Rewrite (C3 lint paragraph + A15 are new) |
| E-2 | `0030:C3` (ROW/RULE split), `0030:S11` | The split is stated for `Rule` only, but two per-row findings name a SECOND row in `Element`: `graph-overlap` (`Element: right.RuleID`, `groups.go:209`) and `graph-redundant-row` (`Element: rows[j].RuleID`, `groups.go:320`). Where the partner is an expanded row, `Element` publishes the bare `retry` — a ROW named by the bare id, the inversion C3 says it closes, on the same finding whose `Rule` it fixed. S11 asserts `rule` only, so nothing pins it either way. | A `graph-overlap` between a literal row and `retry#3` reads `rule: "other", element: "retry"` — the reviewer is pointed at five rows, not one, on the surface C3 promises names the cell. Two findings on one finding record disagree about whether a row is suffixed. | Rewrite (new lint paragraph is silent on `Element`) |
| E-3 | `0030:C3` (fourth hazard), cross-checked against `0030:§performance-expectations` and `0030:A7` | The new lint-cost hazard says the cliff is "reported as `complete=false` rather than as a refusal" and that "the author's only in-band signal is a degraded verdict". The A7 spike records the opposite for the range the hazard names ("a few hundred cells"): every lint run that completed returned `exit 0` with EMPTY findings, and the 1,000-row run was cancelled past 600 s with no verdict at all. `complete=false` is the node-ceiling (4096) path and `graph-product-too-large` the product-bound (2048) path — neither trips between ~200 and 2048 cells, which is precisely the quartic band. | The authoring guide written from C3 tells the author to look for `complete=false`; what the author actually gets is a `lint` that never returns. No in-band signal exists in the band the hazard describes, so the guide's remedy pointer ("narrow the rule's atoms") is attached to a symptom the author will never see. | Rewrite (hazard 4 is new) |
| E-4 | `0030:C1` ("the one signature change this clause implies") vs Phase 2's split | C1 counts ONE signature change: widening the `renderWrites` return / `expand` `writes` pair. Phase 2 places "admitted-cell evaluation" in `renderWrites`, but `renderWrites(rule *sourceRule, id string, isEscape bool)` has no predicate set — the merged predicates are built in `normalizeRule` (`normalize.go:436-442`) and handed only to `expand`. Admitting cells over "the CONJUNCTION of every positive atom the rule AUTHORS on that tag" (C1) therefore needs a second widening (predicates into `renderWrites`) or a third site. Implementable, but not as counted. | An implementer reading C1's "one signature change" and Phase 2's placement cannot satisfy both; either `renderWrites` grows a parameter C1 does not license, or the admits step lands in `normalizeRule`, a site Phase 2 says it does not add. | Rewrite (Phase 2 split is new; C1's count did not follow it) |
| E-5 | `0030:S1` | S1's own worked example of the excluded `rule` field reads "(`retry` vs `retry-2`)". Under C3's lint paragraph the stepped side publishes `retry#2`, not `retry`; S11 says as much ("S1 passes on an implementation that emits `retry` on all five rows and this scenario is the only thing that fails it"). S1 is left describing the pre-C3 shape as the expected difference. | An implementer authoring S1 from its text expects bare `retry` on the stepped side and treats `retry#2` as a fixture drift; S1 and S11 disagree on what the stepped ladder's `rule` field is. | Rewrite (C3/S11 moved; S1's example did not) |

## Justification

### E-1 — `groupHasOverlap` joins per-row `Rule` by bare id

`internal/graphlint/coverage.go:519-530`:

```go
func (a *analysis) groupHasOverlap(g guard.Group) bool {
	ids := ruleIDsOf(g)
	for _, f := range a.findings {
		if f.Code != CodeOverlap || f.Class != "" {
			continue
		}
		if slices.Contains(ids, f.Rule) && slices.Contains(ids, f.Element) {
			return true
		}
	}
	return false
}
```

`ruleIDsOf` (`coverage.go:575-584`) collects bare `row.RuleID`. The
overlap finding is emitted per row at `groups.go:205-209` with
`Rule: rows[i].RuleID, Element: right.RuleID`. C3 now says a per-row
finding on an expanded row publishes `rule#cell` in `Rule`. Once it does,
an overlap whose left row is expanded fails `slices.Contains(ids, f.Rule)`
and `emitCoverageArms` (`coverage.go:391-396`) skips the `ambiguous_match`
arm: `if class == ambiguousMatch && !overlapping { continue }`. The
comment above it states the arm is "treated as vacuously closed" — so the
rescue-row obligation the overlap makes reachable is silently dropped.

This is a group-level consumer reading per-row `Rule` back, which is
exactly the category A15 claims is undisturbed ("without disturbing the
group-level findings"). A15's verification plan checks that per-row sites
hold the `Row` and that `firstRuleID` sites are disjoint — both true — and
never checks for a read-back join, so the plan verifies green over a live
defect. The fix is small (join on `Span`, which stays `model:rule`, or on
a bare id derived from the suffixed form) but C3/A15 must name it; as
written, an implementation satisfying S11 to the letter regresses the
`ambiguous_match` arm.

The same join exists outside the repo: `clierr.go:364` documents `Rule` as
"the source rule/context id", and `docs/cli-output-contract.md:200` lists
`rule` as "graph-lint attribution". Any consumer resolving `rule` against
authored `[[rule]] id`s gets no match on `retry#3`. C3 asserts "no new
vocabulary member" — true — but the field's documented meaning changes
from source-rule id to normalized row identity, and the record should say
`span` remains the bare join key rather than leave that to the reader.

### E-2 — `Element` also names a row

`groups.go:207-209` (`graph-overlap`) and `groups.go:318-320`
(`graph-redundant-row`) set `Element` to the partner row's bare
`RuleID`. C3's split rule is "every surface that names a ROW publishes the
suffixed identity"; `Element` on these two findings names a row. The lint
paragraph fixes `Rule` and is silent on `Element`, so a finding can carry
`rule: "retry#3"` and `element: "retry"` for two rows of the same
mechanism, or `rule: "other", element: "retry"` when the expanded row is
on the right. S11 pins `rule` and the group-level bare id; it does not
pin `element`, so either behaviour passes. The `Message` strings at the
same sites also interpolate the bare ids ("rows %q and %q"), which C3 does
not govern but a reader will see side-by-side with the suffixed `rule`.

Note E-1 and E-2 interact: if `Element` is ALSO suffixed to satisfy the
split, `groupHasOverlap`'s second `slices.Contains(ids, f.Element)` fails
too; if it is not, the split is incoherent. The record has to choose and
state the join.

### E-3 — the fourth hazard misdescribes the in-band signal

C3: "pushes `lint` past the point it can finish, reported as
`complete=false` rather than as a refusal (A7, §Performance
Expectations)" and "the author's only in-band signal is a degraded
verdict that does not name the step as its cause."

`evidence/spikes/a7-row-growth.md` §(c): "Every lint run that completed
returned exit 0 with an empty findings list … Lint does not refuse these
models; it simply gets slower, quartically, until it is unusable past a
few hundred rows." Table row for 1,000 rows: lint ">600s (cancelled)".
The spike then separates the two bound mechanisms — product bound 2048
(`guard/product.go:420`, emits `graph-product-too-large`) and node ceiling
4096 (`reach.go:131` sets `complete=false`) — and says "Not on this
fixture shape" for both.

So in the band C3 names ("a few hundred cells": product ≤ 2048, rows
≤ 4096) neither bound trips and no degraded verdict is emitted; the
process simply does not return. `§performance-expectations` is careful
here ("a model that loads instantly and that `lint` cannot finish
analysing" — no claim of a finding), so the contradiction is between the
new C3 hazard text and its own cited evidence. The authoring-guide
obligation C3 mints ("states the shape of that cliff") is grounded on a
symptom the author will not observe; the guide should say "lint does not
return" and name the fix, not point at `complete=false`.

### E-4 — Phase 2's placement needs a widening C1 does not count

`normalize.go:436-444`: `predicates` (match atoms merged with
`guardAll`/`guardUnless`) is assembled in `normalizeRule`, and
`renderWrites(rule, id, isEscape)` is called immediately after with no
access to it. `normalize.go:580`: `func (l *loader) renderWrites(rule
*sourceRule, id string, isEscape bool) ([]TagValue, []string, error)`.
C1's admits filter is "the CONJUNCTION of every positive atom the rule
AUTHORS on that tag … whether the `in` is authored locally or inherited
from a context, since `expand` receives the merged predicate set". Phase
2 puts admitted-cell evaluation in `renderWrites`. That is only possible
if `renderWrites` receives the merged predicate set — a parameter
widening beyond the return/`writes` pair C1 counts as "the one signature
change this clause implies". The split is implementable (the data is one
line above the call) and does not contradict "one choice-point kind in
`expand`" — row minting stays there — but C1's count is now wrong by one,
and an implementer will have to pick between C1 and Phase 2.

### E-5 — S1's example still describes the pre-split shape

`0030:S1`: "The row-naming fields differ by construction and are
excluded: `rule` (`retry` vs `retry-2`)". C3 and S11 now say the stepped
side publishes `retry#2`. The exclusion itself is still correct (the
values differ either way), but the parenthetical is the one place S1
tells the implementer what to expect on the stepped side, and it names
the shape S11 explicitly calls the failing implementation. Low severity;
one-token fix.

## Per-claim verdicts (1-8)

1. **C3 lint on the ROW side / A15 / S11** — DOES NOT HOLD as written.
   The split is coherent for `Rule` on the emitting side and does not
   contradict the "no new vocabulary member" claim (no field or code is
   added), A8's narrowing (the `resolve.Plan` seam and graph edge `rule`
   are untouched), or `0010:C3`'s `rowByID` (a flow-verb join, not lint).
   It fails on the READ-BACK side inside graphlint (E-1) and on `Element`
   (E-2). A15's Evidence confirms correctly that no non-test
   `Row.Identity()` caller exists and that every per-row site holds the
   `Row`; its verification plan omits the one join that breaks.
2. **C1 calls `graph-coverage-gap` BLOCKING** — HOLDS.
   `taxonomy.go:85-89` lists `CodeCoverageGap` in `blockingCodes`.
   `§consequences` (add-a-tier bullet, drift bullet), A5, S5b, and F4 all
   agree or are tier-silent; none says advisory. The exclude-vs-admit
   argument still works: exclude fails loud at LINT with a named cell and
   a claiming remedy; admit fails at LOAD with a misattributed bound
   message. One wording nit inside the same sentence: "the two directions
   fail differently at LOAD" — the exclude direction fails at lint, as the
   next clause itself says. Not a ledger row.
3. **Four hazards in C3's authoring-guide list** — HOLDS on count
   (domain order; `unless`; `enum` terminal exclusion; lint-cost cliff).
   The fourth hazard's stated symptom is wrong — E-3.
4. **§consequences pins the add-a-tier symptom** — HOLDS. No load
   refusal (the `in` still admits only named cells), no silent widening,
   blocking `graph-coverage-gap` on the new member; consistent with C1's
   tier and A5.
5. **A5 records the `checkRedundantRows` precondition** — HOLDS.
   `groups.go:300-315` confirms same group, `Kind()` match, both sides
   `Projectable()`, `properSubset`. S5b restates the precondition and
   defers the zero-overlap half to S5c; S5c supplies the projectable
   control. A5, S5b, S5c agree. (Observed, out of this pass's scope: MVV
   item 5 still states 5b and 5c as one fixture.)
6. **Phase 2 split** — HOLDS in kind, not in count. Placing cell
   resolution in `renderWrites` and row minting in `expand` matches the
   real shapes (`expand` is package-level with no `*loader`), and does not
   contradict "one choice-point kind in `expand`". It requires
   `renderWrites` to receive the merged predicates, which C1's "one
   signature change" does not admit — E-4.
7. **C2's `Failure` field census** — HOLDS. `category.go:201-239`:
   `Category`, `Detail`, `Offending`, `Remedy`, `Rule`, `Line`, nothing
   else. `Offending`/`Remedy`/`Rule` populated at two sites
   (`load.go:252-254`, `:264`); `RuleKernelOwned`/`RuleAuthorMustRename`
   are the only `Rule` values; `flow_input.go:259` gates the rename hint
   on it. The census and its consequence are correct.
8. **S9 (d) arm; round-trip bullet test-only** — HOLDS. S9 (d) exercises
   a step point sorted against an `in` on another key and agrees with
   D-identity (`retry#<attempt>#<mode>`, key order, no tie possible since
   subsumption removes the same-key match `in` and guard-block `in` is not
   a choice point). §consequences states the equivalence oracle is
   test-only and names the manual mitigation.

## Summary

Five rows. Two are substantive and both sit on the passage the rewrite
added (C3's lint paragraph + A15): an in-package bare-id join the
suffixing breaks (E-1, false green on the `ambiguous_match` arm) and an
unstated `Element` half of the ROW/RULE split (E-2). One factual
misdescription in the new fourth hazard (E-3). Two low-severity
consistency slips (E-4, E-5). Claims 2, 4, 5, 7, 8 hold outright; 3 and 6
hold on structure with a defect in detail; 1 does not hold.
