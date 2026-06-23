# Persona 2 - Implementer

1. **What is the operator/kind compatibility matrix?**
   - Passage: `Normative Contracts`, `Load-Bearing Decisions`.
   - Clarification request: the RDR requires each operator to declare accepted
     tag value kinds but does not actually declare the matrix. Before coding, an
     implementer needs the allowed kind for `eq`, `in`, `lt/lte/gt/gte`,
     `exists`, and set containment.

2. **How should `unless` be represented internally for coverage math?**
   - Passage: `Technical Design`, `A2`.
   - Clarification request: the draft says `unless` is conjunctive and subtracted
     from the row's accepted assignments. It should state that a row's accepted
     domain is `all` intersection minus the single conjunctive `unless`
     intersection, not per-atom negation or multiple separate exclusions.

3. **Does the MVV have to extend the existing spike fixture?**
   - Passage: `A1`, `Minimum Viable Validation`, `Testing Strategy`,
     `Performance Expectations`.
   - Clarification request: the spike output lists five operators and no
     `contains`, but the validation text says representative rows include set
     containment. The implementation plan should distinguish past Resolve
     evidence from the production MVV that still has to cover `contains`.

4. **Which errors are parse errors versus lint errors?**
   - Passage: `Failure Modes`, `Testing Strategy`.
   - Clarification request: unknown operator, unsupported operator/tag-kind
     pair, and literal parse mismatch are pre-resolution validation failures;
     non-exhaustive domains and overlaps are lint findings. The draft should not
     require the implementer to infer this split.
