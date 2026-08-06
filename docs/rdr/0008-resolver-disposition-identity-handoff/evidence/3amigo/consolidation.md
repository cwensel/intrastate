Model: gpt-5

# 3-amigo consolidation — RDR 0008

## Origin ledger

1. **AMIGO-1 — Critical: locator-agreement boundary** (`PM-1`, `IMP-1`,
   `QA-1`). `Critical Assumptions / A4`, `Technical Design`, `Normative
   Contracts`, and `Implementation Plan / Prerequisites` do not yet make the
   locator/model-rule agreement implementable or testable. This is an
   acknowledged re-raise of COVE-3.
2. **AMIGO-2 — High: copy-isolation validation wording** (`IMP-4`, `QA-2`).
   The MVV and Validation scenario 3 must mutate a `Writes` element in the
   original backing array, not merely replace the slice.
3. **AMIGO-3 — High: expansion-suffix positive coverage** (`IMP-5`, `QA-3`).
   At least one successful fixture must carry a non-empty suffix.
4. **AMIGO-4 — High: replay acceptance ownership** (`PM-3`, `QA-5`). The draft
   promises replay consumption but neither names an owning test nor charts the
   downstream obligation.

Consolidation overflow: 1 — invariant-error classification (`IMP-3`, `QA-4`).
