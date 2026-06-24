# Reconcile

| item | source | disposition | evidence pointer or plan |
| --- | --- | --- | --- |
| A5 CI can run `intrastate lint` as the blocking graph-acceptance authority | 1, 2, 4 | DOWNGRADED | RDR A5 is now `Status: Pending`, `Method: MVV Test`. Validation scenario 7 must capture `make check` or `.github/workflows/ci.yml` invoking the production `intrastate lint` command over the checked-in transition model or fixture corpus. Source-search support: `Makefile::check`, `Makefile::build`, `.github/workflows/ci.yml::jobs`, `internal/cli/root.go::NewRootCmd`. |
| Command and CLI authority boundary | 1, 4 | VERIFIED | Absorbed from 3amigo and critique. The current RDR pins root `intrastate lint` as canonical and requires aliases or resolver-local helpers to use the same request builder and graph-lint engine. Evidence: `Technical Design`, `Normative Contracts`, `Validation` scenario 6. |
| Stable finding-code taxonomy and deterministic finding order | 1, 4 | VERIFIED | Absorbed from 3amigo. The current RDR names aggregate `graph-lint-failed`, every mandatory blocking invariant code, and deterministic finding order by identity. Evidence: `Technical Design`, `Normative Contracts`, `Minimum Viable Validation`. |
| Machine-readable multi-finding failure envelope | 1, 4 | VERIFIED | Absorbed from 3amigo and critique. The current RDR requires a `clierr`/`respond`-owned optional typed `findings` field for JSON mode, not a verb-local wrapper or text-only `Detail`. Evidence: `Technical Design`, `Normative Contracts`, `Validation` scenarios 2 and 4. |
| Normalized graph input minimum | 1 | VERIFIED | Absorbed from 3amigo. The current RDR lists the normalized graph value lint consumes: model identity, normalized candidate rows, declared tags/domains, recognized outcomes, terminal and escape data, and accessor/context references. Evidence: `Technical Design`. |

## Absorption Audit

- Grounding: no residue. The grounding pass reported no refuted or missing source claims.
- 3amigo: no residue. Iteration 1 findings were absorbed; iteration 2 found no open original entries and no net-new scope.
- Critique: C1 and C2 absorbed. C3 survived until this reconcile pass because A5 was still stamped `Verified` while the critique required an implementation-time gate proof. This pass downgrades A5 to a pending MVV test and names the exact validation plan.

## Completeness

- No named spike remains unrun before lock. A5's production-gate proof cannot run until the lint command and checked-in model or fixture corpus exist, so it is carried as an implementation-time MVV test rather than a pre-lock spike.
- No `_Draft placeholder._` body section or `this is a seed skeleton` header remains.

## Verdict

RECONCILED. All disturbed items now have a terminal Stage 6 disposition, and no blocker requires returning to Propose, Refine, or Resolve.
