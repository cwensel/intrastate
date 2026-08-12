# intrastate — Recommendation Decisioning Records

Project-scoped RDRs for `intrastate`. Draft new RDRs from the shared
`TEMPLATE.md` in the RDR engine (`$RDR_HOME/TEMPLATE.md`); `/rdr-seed`
materializes a copy automatically. Rationale + the full stage flow live in the
engine README — this file is only the per-project index.

## Index

| ID | Title | Status | Priority |
| --- | --- | --- | --- |
| 0001 | Resolution kernel | Implemented | — |
| 0002 | Transition table as reviewable data | Draft [revised from Final 2026-08-12; re-verify A2, A7] | High |
| 0003 | Guard predicate exhaustiveness | Draft [revised from Final 2026-08-12; re-verify A5] | — |
| 0004 | Accessor execution safety model | Draft [revised from Final 2026-08-12; re-verify none] | — |
| 0005 | Skill integration CLI contract | Final [joint decision → JDR 0001 §JD-8, §JD-9] | High |
| 0006 | Graph lint authority and guarantees | Final [joint decision → JDR 0001 §JD-4] | — |
| [0007](0007-guard-predicate-totality.md) | Guard predicate totality over an incomplete evaluation view | Draft [revised from Final 2026-08-12; re-verify A3, A10, A15] | Medium |
| [0008](0008-recognized-tag-key-ownership.md) | Ownership of the recognized-outcome tag key name | Draft [revised from Final 2026-08-11; re-verify A6, A11] | Medium |
| [0009](0009-escape-row-shape-conformance-ownership.md) | Ownership of escape-row shape conformance | Final [joint decision → JDR 0001 §JD-5] | Medium |

## Status legend

- **Draft** — during the planning/research phase
- **Final** — locked, ready for or during implementation
- **Implemented** — implementation complete
- **Reverted** — implemented then undone (document why)
- **Abandoned** — RDR not implemented
- **Superseded** — replaced by another RDR
- **Demoted** — judged not RDR-shaped; refiled as a plain issue (carry `Demoted [→ <issue link>]`)
