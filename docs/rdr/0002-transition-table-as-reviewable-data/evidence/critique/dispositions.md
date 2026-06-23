# Critique Dispositions

- **fixed** — C-001 unsupported model versions can be accepted because `version` is present but inert: added a normative `[model].version = 1` compatibility gate and carried unsupported-version refusal into MVV/testing. Section touched: Technical Design; Normative Contracts; Load-Bearing Decisions; Minimum Viable Validation; Testing Strategy.
- **fixed** — C-002 expanded-table determinism is under-specified: defined the expanded-table value and deterministic ordering over row identity, predicates, and writes; clarified that the spike SHA is evidence only, while production tests assert the normative value. Section touched: Normative Contracts; Round-Trip / Inverse Invariants; Testing Strategy; Performance Expectations.
- **fixed** — C-003 stable CLI parse/lint failures can collapse into generic errors: required stable data-level validation categories before CLI mapping, with minimum categories for unknown tag/context, non-owned writes, unknown accessors, unsupported version, and ambiguous overlap. Section touched: Technical Design; Normative Contracts; Testing Strategy.

Needs verification: None. The fixes pin design decisions and route deterministic dump behavior into the existing MVV/test plan; they introduce no new external-behavior assumption beyond the fixtures and validation tests already owned by implementation.

Tiebreakers: None.
