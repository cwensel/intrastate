# Post-Mortem: RDR-0028 Write a planned value into a line of a text file without a wrapper script

Opened at a Stage-7.1 route-back, before implementation. Only the Escaped-Defect
Ledger is populated; the remaining post-mortem sections are authored after
implementation.

## Escaped-Defect Ledger

One row per finding in the implementation's `<art>/triage.md`, enriched with
the two fields triage cannot assign. Route-back rows (rdr-common
§punt-ledger) append here at the moment a stage reopens a completed stage —
odc columns n/a until implementation triage.

| Finding | odc-type | odc-trigger | Expected-catching stage | Precursor lens | Escape distance |
| --- | --- | --- | --- | --- | --- |
| C1.6 admits `{tag.<key>}` as a whole argv element at any position with no value rule, so a binding lints green under 0027's inline-shell promise while the executed argv is whatever the caller binds (`["{tag.a}","{tag.b}","{tag.c}"]` + `--tag a=sh --tag b=-c`). `cmdbind.go::substitute` already refuses a `-`-prefixed `{artifact}`; C1.6 carried no such rule and the Cross-Cutting gate answered "adds no new source of caller-supplied value". JC1 compared 0027 on the shared function's clauses and never on the promise C1.6 undercuts. Decided at JDR 0003 §D2 (a): argv0 exclusion + `-` refusal; re-entry at Stage 3, STAGE-SCOPED, re-verify A7, A9. | n/a | n/a | `2-propose` | 3amigo, critique (ran clean; both reviewed C1.6 as a placeholder-vocabulary extension, not as a caller-bound argv word against 0027's promise) | 5 stages (caught at 7.1-cluster-reconcile/whole-set critique; owed at 2-propose, where JC1 disposed 0027 as "no shared decision, no bridge surface") |
