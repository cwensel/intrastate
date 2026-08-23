# Post-Mortem: RDR 0003 Guard Predicate Exhaustiveness

Started at a Stage 6 route-back (rdr-common §punt-ledger), before
implementation. Only the Escaped-Defect Ledger is populated; the remaining
post-mortem sections are authored after implementation.

## Escaped-Defect Ledger

One row per finding in the implementation's `<art>/triage.md`, enriched with
the two fields triage cannot assign. Route-back rows (rdr-common
§punt-ledger) append here at the moment a stage reopens a completed stage —
odc columns n/a until implementation triage.

| Finding | odc-type | odc-trigger | Expected-catching stage | Precursor lens | Escape distance |
| --- | --- | --- | --- | --- | --- |
| Stage 6 (iteration 5) reached lock with A16 booked as a non-blocking peer request while the clause that made it blocking was already in the draft. The repeatability iteration-3 pass introduced the single-value-operator projection clause — `eq`/`in`/comparisons over a dimension not declared `single_valued` do not project at all and take the blocking outcome — and correctly raised A16's consequence on the assumption record itself. But the three surfaces the gate reads were not swept with it: the Authorability paragraph still scoped A16 to "MVV Scenario 5, and Scenario 4's set-kind rejection", the Prerequisites checklist still read "Not lock-blocking", and the aggregate assessment still applied a test ("requires a *different contract here* rather than a field elsewhere") that A16 passes while still emptying the MVV of its only green verdict. The projection clause removes the MVV's positive proof outcome — Scenario 2's passing partition and Scenario 8's green negative control both become unauthorable — so the second Stage 6 hard rule fires. | n/a | n/a | `5-prelock` (the repeatability pass that authored the projection clause owed the consequence sweep across the MVV authorability and gate surfaces, not only the A16 record) | `repeatability` iter-3 (raised A16's consequence on the record; did not sweep the MVV/gate surfaces) | 1 stage (caught at 6-reconcile; owed at 5-prelock) |
