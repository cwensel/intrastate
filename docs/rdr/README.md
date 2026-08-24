# intrastate — Recommendation Decisioning Records

Project-scoped RDRs for `intrastate`. Draft new RDRs from the shared
`TEMPLATE.md` in the RDR engine (`$RDR_HOME/TEMPLATE.md`); `/rdr-seed`
materializes a copy automatically. Rationale + the full stage flow live in the
engine README — this file is only the per-project index.

## Index

| ID | Title | Status | Priority |
| --- | --- | --- | --- |
| 0001 | Resolution kernel | Implemented | — |
| [0002](0002-transition-table-as-reviewable-data.md) | Transition table as reviewable data | Final [joint decision → JDR 0001 §JD-9, §JD-19, §JD-20, §JD-21, §JD-22] | High |
| [0003](0003-guard-predicate-exhaustiveness.md) | Guard predicate exhaustiveness | Final [joint decision → JDR 0001 §JD-18, §JD-22] | High |
| [0004](0004-accessor-execution-safety-model.md) | Accessor execution safety model | Final [joint decision → JDR 0001 §JD-19, §JD-20, §JD-22] | High |
| 0005 | Skill integration CLI contract | Final [joint decision → JDR 0001 §JD-8, §JD-9, §JD-19, §JD-20] | High |
| 0006 | Graph lint authority and guarantees | Final | High |
| [0007](0007-guard-predicate-totality.md) | Guard predicate totality over an incomplete evaluation view | Final [joint decision → JDR 0001 §JD-8, §JD-18, §JD-21, §JD-22] | Medium |
| [0008](0008-recognized-tag-key-ownership.md) | Ownership of the recognized-outcome tag key name | Final [joint decision → JDR 0001 §JD-5, §JD-8, §JD-9] | Medium |
| [0009](0009-escape-row-shape-conformance-ownership.md) | Ownership of escape-row shape conformance | Final [joint decision → JDR 0001 §JD-5, §JD-8] | Medium |

## Status legend

- **Draft** — during the planning/research phase
- **Final** — locked, ready for or during implementation
- **Implemented** — implementation complete
- **Reverted** — implemented then undone (document why)
- **Abandoned** — RDR not implemented
- **Superseded** — replaced by another RDR
- **Demoted** — judged not RDR-shaped; refiled as a plain issue (carry `Demoted [→ <issue link>]`)
