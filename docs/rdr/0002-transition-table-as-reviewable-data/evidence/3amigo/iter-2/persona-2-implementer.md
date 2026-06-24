# 3amigo Iteration 2 — Persona 2 Implementer

- **IMP2-001 should malformed TOML, unknown fields, and missing outcomes get stable data-level categories?** — First-hour implementation question: the parser/validator can encounter malformed TOML before schema validation, unknown root/table fields while decoding, and missing `outcomes`; the RDR names those as failures/tests but omits them from the stable-category contract. Passage: Normative Contracts / validation failures, Failure Modes, Testing Strategy scenario 3.

