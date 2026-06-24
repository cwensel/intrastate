# 3amigo Consolidation

Highest-priority overlapping passages:

1. **Command and CI authority boundary** — PM finding 1, PM finding 3, QA
   finding 5, and implementer finding 1 all hit the same gap: the RDR selects
   `intrastate lint` plus CI, but leaves the exact command authority,
   package/engine reuse boundary, and wrapper drift rules too loose.

2. **Stable finding-code taxonomy** — PM finding 2, implementer finding 3, QA
   finding 1, and QA finding 3 overlap: the RDR requires stable codes but does
   not enumerate the mandatory invariant code names, including the
   inability-to-prove case.

3. **Failure envelope shape for multiple findings** — implementer finding 5 and
   QA findings 2/4 overlap with PM finding 4: the RDR requires rich per-finding
   diagnostics but does not say how a single `CLIError` carries multiple
   blocking findings, how success treats advisories, or how finding order is
   made testable.

4. **Normalized graph interface minimum** — implementer finding 2 is concrete
   and in-scope. It is less cross-persona than the first three, but it gates
   implementation because the lint engine consumes a peer-owned API.
