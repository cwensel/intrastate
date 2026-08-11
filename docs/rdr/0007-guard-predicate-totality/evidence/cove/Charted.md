# Charted to successor — cove iteration 1

Charted: **F-8 `recognized`-key provenance flip.** An owned or observed tag
keyed `"recognized"` overwrites the recognized entry's value *and* provenance in
`resolve.go::assemble`; `recognizedTagKey` is the bare string `"recognized"`
(resolve.go L110) with no reservation or collision check. Because
`missingOwned` keys on `view.has(key, ProvenanceOwned)`, a row listing
`"recognized"` in `RequiresOwned` could be satisfied by such a shadowing tag.

Why out-of-scope for 0007: this is a *key-binding / reservation* question, not a
guard-evaluation-domain question. It does not refute A6a (presence stays
monotone; only value and provenance change), and it changes no verdict under
this RDR's domain rule — `exists` decides on presence, which shadowing cannot
flip. RDR 0007 already declares recognized-tag key binding as a deliberately
separate seam ("katas `z53t` (recognized-outcome tag key binding, now RDR 0008)
… deliberately not folded in here").

Suggested successor: **RDR 0008** (recognized outcome tag key binding) — whether
`recognized` is a reserved key, and what happens when an owned or observed tag
collides with it.
