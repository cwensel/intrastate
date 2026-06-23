# 3amigo Dispositions

- **fixed** - Operator and kind contract; origin: consolidation item 1;
  section touched: `Technical Design`. Added the operator/kind/literal matrix
  required by the existing normative contract.
- **fixed** - Set containment versus existing spike evidence; origin:
  consolidation item 2; sections touched: `Phase 3: Target-Flow Fixture`,
  `Validation`, `Performance Expectations`. Clarified that Resolve evidence
  covered the target-flow subset and that implementation MVV must add
  `contains` before accepting the full operator vocabulary.
- **fixed** - `unless` pass/fail semantics; origin: consolidation item 3;
  sections touched: `Technical Design`, `Testing Strategy`. Stated the
  conjunctive-subtraction model and the disabled-row test case.
- **fixed** - Validation criteria for lint/refusal behavior; origin:
  consolidation item 4; section touched: `Testing Strategy`. Clarified
  source-order independence and multiple-match refusal behavior.

Needs verification:

- Implementation MVV must include at least one `contains` predicate over a
  declared set-valued tag before the full operator vocabulary is accepted.
- Implementation MVV must assert an `unless` disabled-row case where all
  positive predicates match and the row is disabled solely by the full
  conjunctive `unless` block.

Tiebreakers: None.
