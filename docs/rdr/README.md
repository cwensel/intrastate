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
| [0010](0010-stateless-decision-tables.md) | Owned state is optional: stateless decision tables are first-class | Implemented | High |
| [0011](0011-flow-next-match-conditioned-candidates.md) | flow next selects by match; --all enumerates the alphabet | Implemented | High |
| [0012](0012-declared-kind-carrier-at-the-guard-seam.md) | Declared-kind carrier at the guard value seam | Draft | Medium |
| [0013](0013-ungated-proof-completion-observation-seam.md) | Ungated proof-completion observation seam | Draft | Low |
| [0014](0014-repository-gate-verification-oracle.md) | Repository gate verification oracle | Draft | Medium |
| [0015](0015-dead-end-quantifier-over-merged-nodes.md) | Dead-end quantifier over merged nodes | Draft | Medium |
| [0016](0016-artifact-role-reader-cardinality.md) | Reader cardinality of the artifact-role read binding | Draft | Medium |
| [0017](0017-per-finding-code-identity.md) | Per-finding code identity on mixed-disposition rows | Draft | Medium |

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

## Landing

An RDR's build-time **diff oracles** — tests asserting `git diff` against the
branch point, the checkable form of an "X is untouched" claim — are deleted in
the landing commit. They are discharged proof obligations over one branch's
changeset (`0011:S6`: a claim about the DIFF "is discharged at review by
reading the changed test files"), not standing invariants: post-merge there is
no branch point, so they are unrunnable by construction. Deleting one does not
retract the claim; the proof is in git history. A property that must hold going
forward is written as a durable tree/runtime assertion instead — see
`internal/cli/mvv_0010_test.go::TestReq85_TheCheckedInNavigatorModelIsUntouched`.

## Status legend

- **Draft** — during the planning/research phase
- **Final** — locked, ready for or during implementation
- **Implemented** — implementation complete
- **Reverted** — implemented then undone (document why)
- **Abandoned** — RDR not implemented
- **Superseded** — replaced by another RDR
- **Demoted** — judged not RDR-shaped; refiled as a plain issue (carry `Demoted [→ <issue link>]`)
