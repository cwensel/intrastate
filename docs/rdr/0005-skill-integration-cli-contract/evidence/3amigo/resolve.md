# 3amigo Resolve

Grounding sources:

- Code: `internal/cli/respond::Success`, `internal/cli/respond::OK`,
  `internal/cli/respond::Fail`, `internal/cli/respond::ValidateMode`,
  `internal/cli/clierr::CLIError`, and `internal/cli/clierr::ExitCodeFor`
  confirm the existing success envelope and stable error-code surface.
- Resources: `.rdr/resources.md` names `docs/cli-output-contract.md`,
  `CLAUDE.md`, `respond`, and `clierr` as authoritative contracts for CLI I/O.
- Peer RDRs: RDR 0001 owns exact-one resolver selection and refusal classes;
  RDR 0002 owns recognized outcomes, normalized row identity, source locators,
  tag declarations, and accessor references; RDR 0003 owns guard facts; RDR
  0004 owns accessor execution and read-back verification.

Dispositions:

1. fixed — O1 next payload contract — Technical Design, Normative Contracts,
   and Validation now define minimum `next` JSON fields, candidate summaries,
   and same-content text rendering.

2. fixed — O2 request input grammar — Technical Design and Illustrative Code
   now pin MVP `--flow`, repeated `--tag name=value`, `--artifact role=path`,
   and deferred structured-input scope.

3. fixed — O3 stable refusal code list — Failure Modes now names the minimum
   `flow-*` code strings and exit-code groups.

4. fixed — O4 set-state semantics — Technical Design, Normative Contracts,
   Illustrative Code, and Validation now distinguish context `--tag` values
   from planned owned-tag `--write` mutations and read-back-confirmed values.

Needs verification:

- A6 added as Pending. Stage 6/Reconcile must run the MVV through the production
  Cobra path and check pinned request grammar, success payload fields, and
  `flow-*` code spellings before lock.

Tiebreakers:

- None.
