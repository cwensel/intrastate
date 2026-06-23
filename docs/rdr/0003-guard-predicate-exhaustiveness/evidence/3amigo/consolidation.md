# 3amigo Consolidation

Highest-priority rewrite hotspots:

1. **Operator and kind contract** appeared under all three personas. The RDR
   requires an operator/kind compatibility matrix but does not state it.

2. **Set containment versus existing spike evidence** appeared under all three
   personas. The RDR normatively includes set containment while the Resolve
   fixture only exercised `eq`, `in`, `lt`, `gte`, and `exists`.

3. **`unless` pass/fail semantics** appeared under Implementer and QA. The RDR
   should make clear that the whole conjunctive `unless` block is subtracted
   from the row's accepted assignments and must be tested with a disabled-row
   case.

4. **Validation criteria for lint/refusal behavior** appeared under PM and QA.
   The MVV should distinguish parse/type failures from lint failures and assert
   ambiguity/no-match behavior without source-order fallback.
