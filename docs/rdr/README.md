# intrastate — Recommendation Decisioning Records

Project-scoped RDRs for `intrastate`. Draft new RDRs from the shared
`TEMPLATE.md` in the RDR engine (`$RDR_HOME/TEMPLATE.md`); `/rdr-seed`
materializes a copy automatically. Rationale + the full stage flow live in the
engine README — this file is only the per-project index.

## Index

| ID | Title | Status | Priority |
| --- | --- | --- | --- |
| 0001 | Resolution kernel | Implemented | — |
| [0002](0002-transition-table-as-reviewable-data.md) | Transition table as reviewable data | Final [joint decision → JDR 0001 §JD-15, §JD-16, §JD-17] | High |
| [0003](0003-guard-predicate-exhaustiveness.md) | Guard predicate exhaustiveness | Final [joint decision → JDR 0001 §JD-16, §JD-18] | High |
| 0004 | Accessor execution safety model | Final [joint decision → JDR 0001 §JD-15, §JD-17] | High |
| 0005 | Skill integration CLI contract | Final [joint decision → JDR 0001 §JD-8, §JD-9] | High |
| 0006 | Graph lint authority and guarantees | Final [joint decision → JDR 0001 §JD-17] | High |
| [0007](0007-guard-predicate-totality.md) | Guard predicate totality over an incomplete evaluation view | Final [joint decision → JDR 0001 §JD-8, §JD-18] | Medium |
| [0008](0008-recognized-tag-key-ownership.md) | Ownership of the recognized-outcome tag key name | Draft [revised from Final 2026-08-23; re-verify A9, A4 — cluster 0002-0009 iter-2] | Medium |
| [0009](0009-escape-row-shape-conformance-ownership.md) | Ownership of escape-row shape conformance | Final [joint decision → JDR 0001 §JD-5, §JD-8, §JD-15] | Medium |

## Status legend

- **Draft** — during the planning/research phase
- **Final** — locked, ready for or during implementation
- **Implemented** — implementation complete
- **Reverted** — implemented then undone (document why)
- **Abandoned** — RDR not implemented
- **Superseded** — replaced by another RDR
- **Demoted** — judged not RDR-shaped; refiled as a plain issue (carry `Demoted [→ <issue link>]`)
