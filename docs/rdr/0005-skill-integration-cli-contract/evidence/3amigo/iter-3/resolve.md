# 3amigo Delta Resolve - Iteration 3

Grounding sources:

- Code: `internal/cli/respond::OK`, `internal/cli/respond::Fail`, and
  `internal/cli/clierr::CLIError` remain the user-facing success/failure gateway
  for any accessor refusal surfaced by the CLI.
- Resources: `.rdr/resources.md` names `docs/cli-output-contract.md` and
  `CLAUDE.md` as authoritative output contracts; no direct output path is added.
- Peer RDRs: RDR 0001 keeps the resolver kernel free of artifact discovery and
  accessor execution; RDR 0004 requires read, gate, and write accessors to have
  distinct result shapes and says CLI mapping belongs to RDR 0005.

Dispositions:

1. fixed - O5 gate accessor invocation surface - A3/A4, Technical Design,
   Normative Contracts, Implementation Plan, and Validation now state that
   `flow next` / `flow resolve` may invoke declared gate accessors from explicit
   `--artifact role=path` bindings before calling the pure resolver kernel;
   `flow read-state` remains read-accessor-only, and `flow set-state` remains
   write-accessor/read-back-only.

Needs verification:

- A4 flipped to Pending. Stage 6 must verify the new verb-boundary accessor
  contract against RDR 0001 and RDR 0004 before lock.

Tiebreakers:

- None.
