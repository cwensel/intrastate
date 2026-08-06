Model: gpt-5.6-sol

# Critique resolution — RDR 0008

- **fixed** — origin: `CRIT-1`; sections touched: `Critical Assumptions / A5`,
  `Approach`, `Technical Design`, `Normative Contracts`, `Existing
  Infrastructure Audit`, `Risks and Mitigations`, `Failure Modes`,
  `Implementation Plan`, and `Validation`. Current
  `internal/resolver/resolver.go::{Input,Edge}` confirms raw row construction;
  RDR 0007 confirms the planned normalized semantic table value. The draft now
  requires a validated production table boundary that makes independently
  composed identity/action rows invalid or unrepresentable, with ordinary and
  escape tests.
- **fixed** — origin: `CRIT-2`; sections touched: `Critical Assumptions / A4`,
  `Approach`, `Technical Design`, `Normative Contracts`, `Capability
  Dependencies`, `Prerequisites`, and `Validation`. RDR 0002 requires the
  locator to identify at least model id and rule id but does not define a string
  grammar. The draft no longer locks an opaque string; it requires the RDR 0002
  boundary to supply a typed locator and a named constructor/validator. A4
  remains Pending until Stage 6 names and verifies that boundary.
- **fixed** — origin: `CRIT-3`; sections touched: `Critical Assumptions / A3`,
  `Technical Design`, `Normative Contracts`, `Capability Dependencies`,
  `Prerequisites`, `Phase 3`, and `Validation`. RDR 0005 `Technical Design`
  lists matched identity, but its normative `flow resolve` paragraph requires
  only a next tag-set or refusal. A3 is now Pending; RDR 0008 cannot lock until
  RDR 0005 normatively requires and tests direct `SelectedRule` and `Action`
  projection in JSON and text modes.

Needs verification:

- **A3** — align RDR 0005 `Normative Contracts` and production text/JSON
  acceptance with direct selected-rule/action projection.
- **A4** — name and verify the RDR 0002 typed locator constructor/validator,
  zero-value behavior, and model/rule mismatch tests.
- **A5** — name and verify the validated normalized-table production type and
  constructor/API tests for forged ordinary and escape rows.
- Confirm RDR 0001's resolver implementation is present on the implementation
  baseline before Stage 8.

Tiebreakers: none.

## Fresh-context fallback resolution

- **fixed** — fallback `FALLBACK-1`, origin `CRIT-1`; sections touched:
  `Critical Assumptions / A6`, `Approach`, `Technical Design`, `Normative
  Contracts`, `Capability Dependencies`, `Decision Rationale`, `Trade-offs`,
  `Implementation Plan`, and `Validation`. RDR 0002 forbids write/clear data on
  an escape while RDR 0001 and `resolver.Resolve` return an exactly-one escape
  through `plan`; A6 now blocks lock until one production end-to-end meaning is
  chosen.
- **fixed** — fallback `FALLBACK-3`, origins `CRIT-1` and `CRIT-2`; sections
  touched: `Critical Assumptions / A2`, `Approach`, `Technical Design`,
  `Normative Contracts`, `Existing Infrastructure Audit`, `Trade-offs`,
  `Minimum Viable Validation`, and `Validation`. The typed locator and
  snapshot-owned table invalidate A2's string-only, post-plan spike evidence;
  A2 is Pending and the new plan requires external-package ownership tests,
  pre/post-resolution mutation, defensive inspection, and a race test.
- **dismissed-with-cite** — fallback `FALLBACK-2`, origin `CRIT-2`; section
  touched: `Technical Design` and `Load-Bearing Decisions / Identity`. RDR 0007
  owns revision derivation and its `Existing Infrastructure Audit / Replay
  revision claim` charts association enforcement to a successor. The RDR now
  states that `SelectedRule` is logical identity within a revision-bound input,
  not standalone historical event identity; `TableRevision` is not folded into
  this carrier.
- **fixed** — fallback `FALLBACK-3` error-boundary residue, origin `CRIT-1`;
  sections touched: `Technical Design`, `Normative Contracts`, `Prerequisites`,
  `Phase 3`, and Validation scenario 5. User-authored invalidity is a
  load/config failure mapped to stable `CLIError`; only an impossible invalid
  value inside a validated table uses the internal programmer error path.

Needs verification:

- **A2** — rerun ownership evidence against the concrete typed locator and
  opaque normalized-table API.
- **A3** — align RDR 0005 selected-rule/action projection and load/config error
  mapping.
- **A4** — name and verify the typed locator constructor/validator.
- **A5** — name and verify the opaque, snapshot-owned normalized-table API.
- **A6** — reconcile modeled-escape disposition/action semantics end to end.

Tiebreakers: none.
