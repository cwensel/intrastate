# 3amigo Dispositions

- **fixed** — HOT-001 / PM-001, IMP-001, QA-001: added a normative source-schema contract for the Resolve spike field layout and made `Wire / byte format` name the canonical fields and fixtures. Section touched: Proposed Solution / Normative Contracts; Load-Bearing Decisions.
- **fixed** — HOT-002 / IMP-002, QA-002: aligned clear semantics with the verified fixture by making `clear` a rule-level list normalized into `<clear>` writes. Section touched: Technical Design; Normative Contracts.
- **fixed** — HOT-003 / PM-002, IMP-003: replaced unsupported byte-level "source span" language with a source locator contract requiring at least model id and rule id, leaving byte line/column as optional diagnostic detail. Section touched: Technical Design; Normative Contracts; Validation.
- **fixed** — HOT-004 / QA-003: specified that ambiguous selection is tested with a deliberately overlapping malformed negative fixture variant, not only the two positive fixtures. Section touched: Minimum Viable Validation; Testing Strategy.

Needs verification: None. The fixes either pin the already-verified spike layout or narrow unsupported wording; they introduce no new assumption beyond A1/A3/A6.

Tiebreakers: None.
