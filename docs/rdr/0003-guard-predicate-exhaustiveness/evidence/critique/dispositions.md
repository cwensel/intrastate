# Critique Dispositions

- **fixed** - Scoped coverage product; origin: critique item 1; sections
  touched: `Critical Assumptions`, `Approach`, `Technical Design`, `Normative
  Contracts`, `Load-Bearing Decisions`, `Risks and Mitigations`, `Minimum Viable
  Validation`, `Phase 2`, `Testing Strategy`, `Performance Expectations`.
  Grounding: RDR 0002 supplies normalized candidate rows with source identity;
  RDR 0006 owns graph-lint grouping and blocking exhaustiveness/overlap
  findings; this RDR owns predicate-domain semantics.
- **fixed** - Set-valued finite-domain proof limit; origin: critique item 2;
  sections touched: `Critical Assumptions`, `Technical Design`, `Normative
  Contracts`, `Load-Bearing Decisions`, `Risks and Mitigations`, `Failure
  Modes`, `Phase 2`, `Testing Strategy`, `Performance Expectations`. Grounding:
  A2 already allowed an implementation-equivalent bitset; the fix makes the
  deterministic representation and refusal/downgrade behavior explicit.
- **fixed** - Predicate error ownership handoff; origin: critique item 3;
  sections touched: `Critical Assumptions`, `Load-Bearing Decisions`, `Testing
  Strategy`. Grounding: `internal/cli/clierr/clierr.go::CLIError` and
  `internal/cli/respond/respond.go::Fail` provide the current envelope; RDR
  0006 owns graph-lint finding codes, and RDR 0005 owns CLI command/envelope
  mapping.

Needs verification:

- Implementation MVV must include one row-group coverage gap and one row-group
  overlap that are only visible in the multi-dimensional product.
- Implementation MVV must include a set-valued `contains` proof using a
  deterministic symbolic or bitset-equivalent representation.
- Implementation MVV must include refusal/downgrade behavior for a finite
  scoped product too large to prove deterministically.
- RDR 0005 and RDR 0006 must preserve the predicate semantic-kind ownership
  chain when mapping predicate failures to lint findings and CLI envelopes.

Tiebreakers: None.
