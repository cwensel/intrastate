# Persona 3 — QA / Tester

- **QA-001** — I cannot write a pass/fail schema test for exact field names because the RDR says the field layout was finalized but does not make the fixture layout normative. Trigger: `Minimum Viable Validation` and `Validation / Testing Strategy` say to promote the spike fixtures, while `Illustrative Code` says the shape is not load-bearing.
- **QA-002** — I cannot write the explicit-clear test without guessing whether the accepted input is rule-level `clear` or a write-block sentinel. Trigger: `Normative Contracts` conflicts with `evidence/spikes/rdr-fixture.toml`.
- **QA-003** — I cannot write the ambiguous-selection test from the named positive fixtures alone because their four rows do not intentionally overlap. Trigger: `Minimum Viable Validation` says "one ambiguous tag-set is refused" but does not say to use a negative fixture variant.

