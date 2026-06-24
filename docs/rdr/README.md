# intrastate — Recommendation Decisioning Records

Project-scoped RDRs for `intrastate`. Draft new RDRs from the shared
`TEMPLATE.md` in the RDR engine (`$RDR_HOME/TEMPLATE.md`); `/rdr-seed`
materializes a copy automatically. Rationale + the full stage flow live in the
engine README — this file is only the per-project index.

## Index

| ID | Title | Status | Priority |
| --- | --- | --- | --- |
| 0001 | Resolution kernel | Final | — |
| 0002 | Transition table as reviewable data | Draft [revised from Final 2026-06-24; re-verify A2 A3 — cluster 0001-0006 found missing table-modeled escape-edge semantics] | — |
| 0003 | Guard predicate exhaustiveness | Final | — |
| 0004 | Accessor execution safety model | Final | — |
| 0005 | Skill integration CLI contract | Draft [revised from Final 2026-06-24; re-verify A4 — cluster 0001-0006 found read-state/gate accessor result-shape gap] | — |
| 0006 | Graph lint authority and guarantees | Final | — |

## Status legend

- **Draft** — during the planning/research phase
- **Final** — locked, ready for or during implementation
- **Implemented** — implementation complete
- **Reverted** — implemented then undone (document why)
- **Abandoned** — RDR not implemented
- **Superseded** — replaced by another RDR
- **Demoted** — judged not RDR-shaped; refiled as a plain issue (carry `Demoted [→ <issue link>]`)
