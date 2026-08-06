# intrastate — Recommendation Decisioning Records

Project-scoped RDRs for `intrastate`. Draft new RDRs from the shared
`TEMPLATE.md` in the RDR engine (`$RDR_HOME/TEMPLATE.md`); `/rdr-seed`
materializes a copy automatically. Rationale + the full stage flow live in the
engine README — this file is only the per-project index.

## Index

| ID | Title | Status | Priority |
| --- | --- | --- | --- |
| [0001](0001-resolution-kernel.md) | Resolution kernel | Final | — |
| [0002](0002-transition-table-as-reviewable-data.md) | Transition table as reviewable data | Final | High |
| [0003](0003-guard-predicate-exhaustiveness.md) | Guard predicate exhaustiveness | Final | — |
| [0004](0004-accessor-execution-safety-model.md) | Accessor execution safety model | Final | — |
| [0005](0005-skill-integration-cli-contract.md) | Skill integration CLI contract | Final | High |
| [0006](0006-graph-lint-authority-and-guarantees.md) | Graph lint authority and guarantees | Final | — |
| [0007](0007-normalized-table-revision-binding.md) | Normalized-table revision binding | Draft | Medium |
| [0008](0008-resolver-disposition-identity-handoff.md) | Resolver disposition identity handoff | Draft | Medium |

## Status legend

- **Draft** — during the planning/research phase
- **Final** — locked, ready for or during implementation
- **Implemented** — implementation complete
- **Reverted** — implemented then undone (document why)
- **Abandoned** — RDR not implemented
- **Superseded** — replaced by another RDR
- **Demoted** — judged not RDR-shaped; refiled as a plain issue (carry `Demoted [→ <issue link>]`)
