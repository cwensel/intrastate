# Grounding Findings

No REFUTED, NOT-FOUND, or new-rule-with-existing-sibling findings.

Confirmed codebase claims:

- `internal/cli/clierr/clierr.go::CLIError`
- `internal/cli/clierr/clierr.go::ExitCodeFor`
- `internal/cli/clierr/clierr.go::GroupUserEnv`
- `internal/cli/respond/respond.go::OK`
- `internal/cli/respond/respond.go::Fail`
- `internal/cli/respond/respond.go::ValidateMode`
- `internal/cli/root.go::NewRootCmd`
- `internal/cli/root.go::ExecuteAndEmit`
- `Makefile::check`
- `.github/workflows/ci.yml::jobs`

Confirmed documentation claims:

- `docs/cli-output-contract.md` defines the text/json output envelope and
  routes exit-code meaning to `clierr.ExitCodeFor`.

Confirmed peer-RDR claims:

- RDR 0002 defines sparse TOML transition data, normalized candidate rows,
  source rule ids/source locators, deterministic dumps, tag declarations, and
  explicit writes/clears for lint and resolver use.
- RDR 0003 defines symbolic predicate atoms, finite-domain coverage/overlap
  proof semantics, source identity for predicate diagnostics, and
  owned/observed/recognized provenance checks.
- RDR 0004 keeps accessor execution outside predicate evaluation and defines
  owned-tag write/read-back safety at the accessor boundary.
- RDR 0005 owns the resolver CLI surface and output-envelope mapping while
  leaving static graph lint authority to RDR 0006.

Confirmed sibling prior-art claims:

- `../state-machines/MODEL-transition.md` names design-time lint failures for
  overlapping predicates, non-exhaustive predicates, and owned-tag
  read-before-write.
- `../state-machines/attic/RESOLVER-DESIGN.md` names a design-time CI lint over
  table data with no dangling edges, no dead ends except declared terminals,
  determinism, and closed guard coverage.
- `../state-machines/repos/README.md` records the small-graph conclusion:
  declarative table plus short CLI/lint instead of adopting a formal model
  checker as the runtime authority.

Inverse search:

- Searched `internal/`, `cmd/`, `Makefile`, `.github/`, and the CLI output
  contract for lint, transition, resolver, guard, predicate, candidate row,
  owned/recognized outcome, and state implementation symbols. No existing
  transition graph model, resolver, guard predicate evaluator, graph lint
  command, or sibling discriminator already makes the graph-lint authority
  decision. Existing repository gates can host the future command, but neither
  `make check` nor CI currently invokes `intrastate lint` because the command
  does not exist yet.
