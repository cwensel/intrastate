# 3amigo Resolve

- **fixed** - origin: command and CI authority boundary; section touched:
  `Technical Design`, `Normative Contracts`, `Validation`. Grounding: RDR 0005
  leaves graph lint authority to RDR 0006 while requiring CLI verbs to use
  `respond`; this RDR now names `intrastate lint`, the internal graph-lint
  package boundary, CI use of the production command, and wrapper drift rules.
- **fixed** - origin: stable finding-code taxonomy; section touched:
  `Technical Design`, `Normative Contracts`, `Minimum Viable Validation`.
  Grounding: RDR 0006 already owns the invariant taxonomy and `clierr.CLIError`
  supports stable codes; the draft now pins the aggregate error code and each
  mandatory blocking finding code.
- **fixed** - origin: failure envelope shape for multiple findings; section
  touched: `Technical Design`, `Normative Contracts`, `Validation`. Grounding:
  `internal/cli/clierr/clierr.go::CLIError` is append-only with `omitempty`
  optional fields, and `internal/cli/respond/respond.go::Fail` is the gateway;
  the draft now requires one aggregate `graph-lint-failed` error with
  machine-readable finding records in JSON mode and `GroupUserEnv` exit
  behavior.
- **fixed** - origin: normalized graph interface minimum; section touched:
  `Technical Design`, `Implementation Plan`. Grounding: RDR 0002 defines
  normalized candidate rows with source identity, predicates, writes, and tag
  declarations; RDR 0003 defines finite-domain predicate semantics and
  provenance checks. The draft now lists the minimum normalized graph value
  lint consumes.

Needs verification:

- None. The edits pin implementation/test contracts using already-verified
  peer RDRs and the existing CLI error gateway. The new deterministic ordering
  and finding-code claims are covered by the updated MVV scenarios.

Tiebreakers:

- None.
