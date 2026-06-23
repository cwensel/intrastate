# Grounding Findings

No REFUTED, NOT-FOUND, or new-rule-with-existing-sibling findings.

Confirmed codebase claims:

- `internal/cli/clierr/clierr.go::CLIError`
- `internal/cli/respond/respond.go::Fail`
- `internal/cli/config/config.go::Load`

Confirmed peer-RDR claims:

- RDR 0002 defines tag declarations with `owned`, `observed`, and `recognized`
  provenance.
- RDR 0002 defines positive `all` and negative `unless` guard lists and requires
  normalized rows to retain source rule ids/source locators.
- RDR 0006 consumes the normalized candidate-row graph and finite-domain
  predicate semantics for graph lint.

Confirmed spike claims:

- `docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes/guard-fixture.toml`
  uses `eq`, `in`, `lt`, `gte`, and `exists` across the representative guard
  rows.
- `docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes/output.txt`
  records the successful fixture check.

Inverse search:

- Searched `internal/` and `cmd/` for guard, predicate, transition, resolver,
  recognized outcome, tag, and state implementation symbols. No existing guard
  predicate evaluator, predicate identity rule, or transition resolver
  implementation found.
