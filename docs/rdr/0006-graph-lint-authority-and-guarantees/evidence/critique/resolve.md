# Critique Resolve

- **fixed** - origin: C1 command authority drifts away from the resolver CLI
  namespace; section touched: `Proposed Solution`, `Technical Design`,
  `Normative Contracts`, `Implementation Plan`, `Validation`. Grounding: RDR
  0005 owns runtime `flow next` / `flow resolve` and leaves graph lint
  authority to RDR 0006; this draft now pins root `intrastate lint` as the
  canonical design-time command and requires any alias or resolver-local helper
  to share the same request builder and graph-lint engine.
- **fixed** - origin: C2 machine-readable findings are promised but not made
  implementable; section touched: `Technical Design`, `Normative Contracts`,
  `Load-Bearing Decisions`, `Implementation Plan`. Grounding:
  `internal/cli/clierr/clierr.go::CLIError` currently has only Code, Message,
  Param, Detail, Hint, Group, and Cause; `internal/cli/respond/respond.go::Fail`
  emits that envelope. The draft now requires a respond/clierr-owned optional
  typed findings field serialized in JSON, not a verb-local wrapper or
  text-only Detail string.
- **fixed** - origin: C3 CI authority is named but not enforceable from the
  repo's actual gate; section touched: `Critical Assumptions`,
  `Implementation Plan`, `Validation`, `Finalization Gate`. Grounding:
  `.github/workflows/ci.yml` and `Makefile` exist, but no graph-lint command or
  checked-in model corpus exists yet. A5 is no longer marked Verified by an MVV
  promise; it is Pending with a Spike plan to wire `make check` or the GitHub
  workflow to the production `intrastate lint` command once implementation
  artifacts exist.

Needs verification:

- A5 is Pending. Verify during Stage 6 after implementation artifacts exist:
  capture the repository gate running the production `intrastate lint` command
  over the checked-in transition model or lint fixture corpus, with one legal
  model passing and illegal fixtures for every blocking invariant class failing.

Tiebreakers:

- None.
