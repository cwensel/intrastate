# 3amigo Iteration 2 — Persona 3 QA / Tester

- **QA2-001 negative tests lack a normative category oracle for three promised failures** — I can write malformed TOML, unknown-schema-field, and missing-root-outcome-alphabet fixtures from the RDR text, but before the fix I cannot assert their required stable data-level categories from the normative contract. Passage: Normative Contracts / validation failures, Failure Modes, Testing Strategy scenario 3.

