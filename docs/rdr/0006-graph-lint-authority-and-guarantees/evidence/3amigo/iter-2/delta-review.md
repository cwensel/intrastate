# 3amigo Delta Review - Iteration 2

Scope: original 3amigo ledger entries only.

Result: no open entries remain.

- Command and CI authority boundary is now explicit: the authoritative surface
  is `intrastate lint`, CI runs that production command shape, and hooks,
  wrappers, or resolver-local flags may call the same engine but cannot define
  different acceptance rules.
- Stable finding-code taxonomy is now explicit: the RDR pins
  `graph-lint-failed` as the aggregate error and names each mandatory blocking
  invariant code, including `graph-unprovable-coverage`.
- Failure envelope shape is now explicit: multiple blocking findings are carried
  under one aggregate `CLIError`, JSON mode keeps findings machine-readable, and
  deterministic finding order is testable by identity.
- Normalized graph interface minimum is now explicit: lint consumes a normalized
  graph value with model identity, rows, declarations, finite-domain metadata,
  terminals/escapes, and normalized references, not sparse TOML or Cobra state.

Net-new scope:

- None.
