# Grounding Findings

No REFUTED, NOT-FOUND, or new-rule-with-existing-sibling findings.

Confirmed codebase claims:

- `internal/cli/clierr/clierr.go::CLIError`
- `internal/cli/clierr/clierr.go::ExitCodeFor`
- `internal/cli/respond/respond.go::ValidateMode`
- `internal/cli/respond/respond.go::OK`
- `internal/cli/respond/respond.go::Fail`
- `internal/cli/root.go::ExecuteAndEmit`

Inverse search:

- Searched `internal/` and `cmd/` for resolver, transition, guard, tag,
  predicate, recognized outcome, and accessor implementation symbols. No
  existing resolver kernel, transition table, guard evaluator, or accessor
  executor implementation found.
