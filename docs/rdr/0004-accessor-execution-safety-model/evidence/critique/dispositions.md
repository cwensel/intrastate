# Critique Dispositions

- fixed - origin: write verification blesses collateral mutation - sections
  touched: Critical Assumptions A7, Technical Design, Normative Contracts,
  Round-Trip / Inverse Invariants, Implementation Plan Phase 3, Testing
  Strategy, Validation Scenario 3.
- dismissed-with-cite - origin: timeout refusal is mistaken for rollback -
  sections touched: none. The RDR already declines undo/transaction semantics in
  Alternatives and Day 2 Operations, and the critique did not ground a new
  implementable contract beyond retry discipline.
- dismissed-with-cite - origin: typed capability names become a thin wrapper
  over arbitrary integration code - sections touched: none. The RDR already
  rejects raw shell-out, ambient artifact discovery, and host callbacks in
  Normative Contracts and Alternatives; additional binding-interface shape
  belongs to implementation design unless a concrete unsafe metadata field is
  introduced.

Needs verification:

- A7 is newly Pending. Verify by extending the spike and MVV tests so a write
  accessor that changes an observed or recognized tag while writing the expected
  owned tag reports `read_back_mismatch`.

Tiebreakers: None.
