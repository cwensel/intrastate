# Persona 3 - QA / Tester

- Minimum Viable Validation: cannot write a complete gate-result test from the
  MVV summary because it names gate-indeterminate but not gate-denied. The
  detailed Testing Strategy includes deny, so the MVV should match it.
- Normative Contracts / Testing Strategy: cannot write a timeout validation test
  unless the spec says whether zero, negative, or missing timeouts are invalid.
- Normative Contracts / Testing Strategy: cannot write a read-back boundary test
  unless the spec states that read-back verification uses the same caller-supplied
  artifact role from the write binding and must not discover another artifact
  from ambient process state.
- Validation Scenario 1: cannot write a duplicate-binding validator test with a
  normative source unless the RDR explicitly says duplicate or multiply-bound
  accessor identities are validation failures.
