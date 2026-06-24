# Persona 2 - Implementer

- Technical Design / Normative Contracts: "timeout policy" is named as required
  model data, but the implementer would ask what makes timeout metadata valid.
  The RDR should say that a missing or non-positive timeout is a validation
  failure before execution.
- Technical Design / Round-Trip Invariants: "re-read through the same accessor
  boundary" is not precise enough for implementation. Does the write path invoke
  an independent read accessor, reuse the write accessor's artifact role, or
  inspect ambient state? The text should pin the verification boundary.
- Technical Design / Validation: duplicate or multiply-bound accessors are
  mentioned in tests and Load-Bearing Decisions, but not in the normative
  contract. The validator needs a pass/fail rule for missing and multiply-bound
  accessor identities.
