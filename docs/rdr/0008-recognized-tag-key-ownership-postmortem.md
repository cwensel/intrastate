# Post-Mortem: RDR-0008 Ownership of the recognized-outcome tag key name

Ledger-only. Opened at Stage 6 reconcile (re-entry) to record a route-back
that had escaped without a row. The remaining post-mortem sections are
authored after implementation.

## Escaped-Defect Ledger

| Finding | odc-type | odc-trigger | Expected-catching stage | Precursor lens | Escape distance |
| --- | --- | --- | --- | --- | --- |
| A6/A11 stamped `Verified` on a false negative-existential about peer RDR 0009 ("the string `Input` appears zero times"; it appears 5x in the Final text), plus two bare peer line citations that went stale when the peer re-locked | n/a until implementation triage | n/a until implementation triage | 4-resolve | none (critique, 3amigo, repeatability, cove all ran clean) | 3 stages (4-resolve -> 7-finalize -> 7.1-cluster-reconcile) |
| Failure Modes claimed "no user data is reclassified as a programmer mistake" on the strength of A10, whose sweep covers only the accessor half of the data channel; the caller-supplied `--tag` half is user data and is reclassified | n/a until implementation triage | n/a until implementation triage | 4-resolve | none (A8's "zero non-test callers at HEAD" sweep ran clean and was correct at HEAD) | 3 stages (4-resolve -> 7-finalize -> 7.1-cluster-reconcile) |
| A9/A4 stamped `Verified` against RDR 0002's normalizer spike as it stood 2026-08-11 (no outcome field; three fixtures to rename); 0002's Stage 4 rebuilt the spike with outcome lifting, renamed its fixtures, and re-locked with fenced text that refuses guard-position `recognized` atoms and makes load fail-fast — block 2's predicate-position clause and scenario 5's two-failure count contradicted Final 0002 | n/a until implementation triage | n/a until implementation triage | 4-resolve | none (the first re-entry's Stage 4/6/7 all ran after 0002's Stage 4 rebuild and re-read neither the peer's spike nor its Draft text) | 3 stages (4-resolve -> 7-finalize -> 7.1-cluster-reconcile) |

**Notes.**

- Both rows share one root cause: an assumption verified against a *moving*
  peer or a *future* caller was stamped with a point-in-time fact and no
  re-verification trigger. The first read Draft 0009 and never re-read the
  Final; the second swept a package with zero non-test callers and generalized
  past the callers RDR 0005 has not written yet.
- The first row's citations were bare peer line numbers, which is what made the
  staleness invisible. Stage 6 has since converted every live peer-RDR citation
  in this RDR to a durable anchor (section heading or assumption ID) per
  rdr-common's anchor doctrine.
- The third row is the same root cause a second time, against a different peer artifact: the first re-entry (2026-08-21) re-verified A6/A11 against Final 0009 but carried A9/A4 forward as `Verified` while 0002 — then `Draft [revised from Final]` and actively rebuilding its spike — was moving. A scoped re-entry that carries assumptions forward needs the same "peer is moving" check on every carried-forward assumption whose Evidence names a peer artifact, not only on the listed IDs.
- Neither escape was caught by a pre-lock lens. All four required lenses ran and
  returned clean on these passages, so the miss is not lens coverage but the
  absence of a "peer is Final and may move" check at Stage 4.
