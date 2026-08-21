# Post-Mortem: RDR 0007 Guard predicate totality over an incomplete evaluation view

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
| Stage 6 reached with two of four required `foundational` lenses never run against the re-entered design. After the JDR 0001 demotion the draft was rewritten (1666 insertions / 2137 deletions on a 1927-line file), dissolving the very seam `critique` and `repeatability` had reviewed — critique's central finding C-2 targeted assumption A10, which no longer exists, and all three repeatability runs reconstructed the superseded `Evaluate(guard string, view TagSet)` seam. Both lenses' iteration-1 evidence therefore attests a document that no longer exists, and the Stage 5 lens row was treated as discharged by it. | n/a | n/a | `5-prelock` | none (cove and 3amigo both ran iter-2 and ran clean; the two stale lenses did not re-run) | 1 stage (caught at 6-reconcile; owed at 5-prelock) |
