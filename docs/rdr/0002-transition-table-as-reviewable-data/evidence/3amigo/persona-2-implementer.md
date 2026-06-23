# Persona 2 — Implementer

- **IMP-001** — Which exact TOML keys must production structs decode: the spike's root `outcomes`, `[model]`, `[tags.<tag>]`, `[accessors.<id>]`, `[context.<id>]`, `[[rule]]`, `[rule.write]`, rule-level `clear`, and `[dump]`, or are those still illustrative? Trigger: `Load-Bearing Decisions / Wire` says exact field names are finalized during Resolve, but the `Illustrative Code` block says it is not a locked schema.
- **IMP-002** — Is explicit clearing inside `[rule.write]` or a rule-level `clear` list? Trigger: `Normative Contracts` says the clear sentinel is in the write block; the verified RDR fixture uses `clear = ["prelock_lens"]` at rule level and the spike normalizer turns that into `<clear>`.
- **IMP-003** — What source-location contract must normalization preserve? Trigger: `Normative Contracts` requires a source span, while the spike rows preserve model id, rule id, and the `source` string but not byte line/column spans.

