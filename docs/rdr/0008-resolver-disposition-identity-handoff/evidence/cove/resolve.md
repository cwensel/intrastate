Model: gpt-5

# COVE resolve — RDR 0008

- **fixed** — origin: COVE-1 mainline grounding mismatch; sections touched: `Technical Environment`, `Key Discoveries`, `Implementation Plan / Prerequisites`. The draft now distinguishes the integration-branch evidence from `main` and makes the RDR 0001 implementation an explicit baseline prerequisite.
- **fixed** — origin: COVE-2 consumer-test ownership contradiction; sections touched: `Implementation Plan / Phase 3`, `Validation / Testing Strategy`. RDR 0008 now exposes the handoff; RDR 0005 owns the eventual production CLI acceptance fixture.
- **fixed** — origin: COVE-3 identity/provenance consistency silence; sections touched: `Critical Assumptions`, `Technical Design`, `Normative Contracts`, `Failure Modes`, `Implementation Plan / Prerequisites`. Normalization must reject locator/model-rule disagreement before resolution.

Needs verification:

- A4 is Pending. Stage 6 must verify the new normalized identity/provenance consistency requirement against RDR 0002 and ensure its implementation boundary can reject disagreement before `Resolve`.
- Confirm the RDR 0001 resolver implementation is present on the implementation baseline before Stage 8 work begins.

Tiebreakers: none.
