# intrastate — Recommendation Decisioning Records

Project-scoped RDRs for `intrastate`. Draft new RDRs from the shared
`TEMPLATE.md` in the RDR engine (`$RDR_HOME/TEMPLATE.md`); `/rdr-seed`
materializes a copy automatically. Rationale + the full stage flow live in the
engine README — this file is only the per-project index.

## Index

| ID | Title | Status | Priority |
| --- | --- | --- | --- |
| 0001 | Resolution kernel | Implemented | — |
| [0002](0002-transition-table-as-reviewable-data.md) | Transition table as reviewable data | Final [all joint decisions answered] | High |
| [0003](0003-guard-predicate-exhaustiveness.md) | Guard predicate exhaustiveness | Final [joint decision → JDR 0001 §JD-18] | High |
| [0004](0004-accessor-execution-safety-model.md) | Accessor execution safety model | Final [all joint decisions answered] | High |
| [0005](0005-skill-integration-cli-contract.md) | Skill integration CLI contract | Final [all joint decisions answered] | High |
| 0006 | Graph lint authority and guarantees | Final | High |
| [0007](0007-guard-predicate-totality.md) | Guard predicate totality over an incomplete evaluation view | Final [joint decision → JDR 0001 §JD-18] | Medium |
| [0008](0008-recognized-tag-key-ownership.md) | Ownership of the recognized-outcome tag key name | Final [joint decision → JDR 0001 §JD-5] | Medium |
| [0009](0009-escape-row-shape-conformance-ownership.md) | Ownership of escape-row shape conformance | Final [joint decision → JDR 0001 §JD-5] | Medium |

## Implementing

Read **[BUILD-ORDER.md](BUILD-ORDER.md)** before `/rdr-implement` on any
`0002-0009` member. Build order is **not** the lock order the cluster gate
reasoned about, and the two run in opposite directions at the head of the
graph: **RDR 0007 Phase 1 defines the `internal/resolve::Row` that RDR 0002
Phase 2 normalizes to**, so 0007 goes first even though 0002 produces the wire
format. Neither open joint decision (§JD-5, §JD-18) blocks the build.

## Status legend

- **Draft** — during the planning/research phase
- **Final** — locked, ready for or during implementation
- **Implemented** — implementation complete
- **Reverted** — implemented then undone (document why)
- **Abandoned** — RDR not implemented
- **Superseded** — replaced by another RDR
- **Demoted** — judged not RDR-shaped; refiled as a plain issue (carry `Demoted [→ <issue link>]`)
