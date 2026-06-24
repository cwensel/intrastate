# Cluster Reconcile Report 0001-0006

Date: 2026-06-24

## Membership

The cluster is valid for Stage 7.1:

| RDR | Status before pass | Implemented? | Relatedness |
| --- | --- | --- | --- |
| 0001 | Final | No | Foundational resolver kernel; named predecessor/peer for the set. |
| 0002 | Final | No | Predecessor depends on 0001; owns table contract consumed by 0001/0003/0006. |
| 0003 | Final | No | Predecessor depends on 0001; owns guard semantics consumed by 0002/0006. |
| 0004 | Final | No | Depends on 0001/0002/0003; owns accessor execution consumed by 0005. |
| 0005 | Final | No | Depends on 0001/0002/0003 and maps 0004 accessor outcomes to CLI. |
| 0006 | Final | No | Depends on 0001/0002/0003 and consumes normalized rows/predicates. |

## Pairwise Table

| Pair | Finding type | Severity | Disposition | RDR returns to Draft | Target stage | Scope | Why |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 0001-0002 | gap | blocks-impl | SPEC-DEFECT | 0002 | 4 re-resolve | STAGE-SCOPED | 0001 requires explicit table-modeled escape edges; 0002 owns the table contract but does not define the syntax or semantics. |
| 0001-0002 | contradiction | risks-impl | SPEC-DEFECT | 0002 | 4 re-resolve | STAGE-SCOPED | RDR 0002 blurs malformed table rules into typed refusals, disturbing the resolver/table validation boundary. |
| 0002-0003 | gap | risks-impl | Covered by 0002 SPEC-DEFECT | 0002 | 4 re-resolve | STAGE-SCOPED | Recognized outcome alphabet versus recognized tag declaration invariant should be fixed with 0002's reopened table contract. |
| 0002-0003 | duplication | cosmetic | NO CONFLICT | - | - | - | `all`/`unless` duplication is benign if 0002 remains the container owner and 0003 remains the semantic owner. |
| 0004-0005 | gap | risks-impl | SPEC-DEFECT | 0005 | 4 re-resolve | STAGE-SCOPED | RDR 0005 allows `read-state` to invoke gate accessors but still promises a tag-set response; RDR 0004 gates return allow/deny/indeterminate. |
| 0004-0005 | round-trip | cosmetic | NO CONFLICT | - | - | - | RDR 0005's CLI-visible set/read invariant correctly delegates protected non-owned tag checks to RDR 0004. |
| 0005-0006 | gap | risks-impl | NO CONFLICT | - | - | - | RDR 0005's "no new envelope fields" claim is scoped to resolver/accessor failures; RDR 0006 owns the graph-lint findings exception. |
| 0005-0006 | no conflict | cosmetic | NO CONFLICT | - | - | - | Runtime `flow` verbs and root `intrastate lint` authority are separated correctly. |

## Verdict

NOT RECONCILED.

The cluster may not proceed to implementation until the demoted RDRs re-enter,
repair, and re-lock:

- RDR 0002 is now Draft with a cluster re-entry note. Re-enter at Stage 4
  (`$rdr-resolve 0002`) with STAGE-SCOPED re-verification of A2 and A3, then
  continue through reconcile/finalize.
- RDR 0005 is now Draft with a cluster re-entry note. Re-enter at Stage 4
  (`$rdr-resolve 0005`) with STAGE-SCOPED re-verification of A4, then continue
  through reconcile/finalize.

RDRs 0001, 0003, 0004, and 0006 remain Final.

## Review Gate

- Cluster membership is correct: all six members were Final and unimplemented,
  and the relatedness is explicit through predecessors, peer assumptions, and
  shared resolver/table/guard/accessor/CLI/lint concerns.
- No SPEC-DEFECT was fixed by editing a Final design body. RDR 0002 and RDR
  0005 were flipped to Draft and annotated with re-entry context.
- The less foundational RDR was demoted in each blocking pair: RDR 0002 yields
  to RDR 0001's resolver contract; RDR 0005 yields to RDR 0004's accessor
  result contract.
- Scope is the narrowest that covers the blast radius: STAGE-SCOPED Stage 4
  re-entry, because the chosen approaches survive but specific assumptions and
  peer evidence must be re-verified before re-lock.

