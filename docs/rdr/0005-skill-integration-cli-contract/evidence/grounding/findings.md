# Grounding Findings

No REFUTED, NOT-FOUND, or new-rule-with-existing-sibling findings.

Confirmed codebase claims:

- `internal/cli/respond/respond.go::Success`
- `internal/cli/respond/respond.go::OK`
- `internal/cli/respond/respond.go::Fail`
- `internal/cli/respond/respond.go::ValidateMode`
- `internal/cli/clierr/clierr.go::CLIError`
- `internal/cli/clierr/clierr.go::ErrorCode`
- `internal/cli/clierr/clierr.go::ExitCodeFor`
- `internal/cli/root.go::NewRootCmd`
- `internal/cli/root.go::ExecuteAndEmit`
- `internal/cli/version.go::newVersionCmd`

Confirmed documentation claims:

- `docs/cli-output-contract.md` defines `--as text|json`, text stdout/stderr
  behavior, JSON stdout terminal records, JSON advisory stderr records, and
  points exit-code mapping at `clierr.ExitCodeFor`.

Confirmed peer-RDR claims:

- RDR 0001 defines the stateless resolver kernel and delegates CLI behavior to
  RDR 0005.
- RDR 0002 defines the transition-model/outcome-alphabet seam consumed by
  `flow next` and `flow resolve`.
- RDR 0003 owns symbolic guard predicate semantics outside the CLI.
- RDR 0004 owns accessor execution and read-back verification.
- RDR 0006 owns graph lint authority outside this contract.

Confirmed sibling prior-art claims:

- `../state-machines/attic/RESOLVER-DESIGN.md` names `resolve`, `next`,
  `read-state`, and `set-state`, keeps the resolver stateless, and makes
  `set-state` persist an already-decided result with read-back verification.
- `../state-machines/attic/RESOLVER-CLI.md` names the same four-operation thin
  CLI and rejects a runtime/driver shape.
- `../state-machines/MODEL-transition.md` names `next(state-tags)`,
  `resolve(state-tags)`, `read-state(location)`, and
  `set-state(location, next-tags)`, with `resolve` pure and location-free.

Inverse search:

- Searched `internal/` and `cmd/` for resolver, transition, table loader,
  guard, predicate, accessor, graph lint, `flow`, recognized outcome, outcome
  alphabet, `read-state`, and `set-state` implementation symbols. No existing
  resolver CLI surface, table loader, accessor executor, graph lint command, or
  sibling discriminator that already decides the verb set was found.
