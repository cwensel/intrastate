# Consolidation

Origin ledger for 3amigo resolve:

1. **O1 next payload contract** — PM, implementer, and QA all flagged that
   `flow next` says "legal-outcome alphabet plus conditional summaries" without
   enough payload shape to implement or test.

2. **O2 request input grammar** — PM, implementer, and QA all flagged that model
   selection, repeated tags, structured input files, and artifact role bindings
   are described broadly but not pinned enough for the MVP.

3. **O3 stable refusal code list** — PM, implementer, and QA all flagged that
   the RDR requires stable codes while leaving canonical resolver/accessor code
   strings and group mapping unspecified.

4. **O4 set-state semantics** — implementer and QA both flagged that `set-state`
   examples and prose do not distinguish target owned tags, observed context,
   resolver-planned writes, and read-back expectations.

Highest-priority rewrites: Normative Contracts, Technical Design, Illustrative
Code, Failure Modes, and Validation / Testing Strategy.
