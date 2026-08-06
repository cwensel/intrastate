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
