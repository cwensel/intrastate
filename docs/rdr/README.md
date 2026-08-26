# intrastate — Recommendation Decisioning Records

Project-scoped RDRs for `intrastate`. Draft new RDRs from the shared
`TEMPLATE.md` in the RDR engine (`$RDR_HOME/TEMPLATE.md`); `/rdr-seed`
materializes a copy automatically. Rationale + the full stage flow live in the
engine README — this file is only the per-project index.

## Index

| ID | Title | Status | Priority |
| --- | --- | --- | --- |
| 0001 | Resolution kernel | Implemented | — |
| [0002](0002-transition-table-as-reviewable-data.md) | Transition table as reviewable data | Implemented | High |
| [0003](0003-guard-predicate-exhaustiveness.md) | Guard predicate exhaustiveness | Implemented | High |
| [0004](0004-accessor-execution-safety-model.md) | Accessor execution safety model | Implemented | High |
| [0005](0005-skill-integration-cli-contract.md) | Skill integration CLI contract | Implemented | High |
| 0006 | Graph lint authority and guarantees | Implemented | High |
| [0007](0007-guard-predicate-totality.md) | Guard predicate totality over an incomplete evaluation view | Implemented | Medium |
| [0008](0008-recognized-tag-key-ownership.md) | Ownership of the recognized-outcome tag key name | Implemented | Medium |
| [0009](0009-escape-row-shape-conformance-ownership.md) | Ownership of escape-row shape conformance | Implemented | Medium |

## Implementing

Read **[BUILD-ORDER.md](BUILD-ORDER.md)** before `/rdr-implement` on any
`0002-0009` member. Each RDR is implemented to completion in one run; the run
order is:

```
0007 → 0002 → 0003 → 0006 → 0004 → 0005 → 0008 → 0009
```

The first two are forced: **RDR 0007 defines the `internal/resolve::Row` that
RDR 0002 normalizes to**, so 0007 goes first even though 0002 produces the wire
format — run order is not the lock order the cluster gate reasoned about, and
the two invert at the head of the graph. 0002 must then precede 0003, because
JDR 0001 §D13 moved the set-value encoding declaration to 0002 and that is what
broke the 0007↔0003 cycle. Neither open joint decision (§JD-5, §JD-18) blocks
any run.

## Status legend

- **Draft** — during the planning/research phase
- **Final** — locked, ready for or during implementation
- **Implemented** — implementation complete
- **Reverted** — implemented then undone (document why)
- **Abandoned** — RDR not implemented
- **Superseded** — replaced by another RDR
- **Demoted** — judged not RDR-shaped; refiled as a plain issue (carry `Demoted [→ <issue link>]`)
