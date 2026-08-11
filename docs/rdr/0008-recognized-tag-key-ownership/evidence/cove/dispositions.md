Model: claude-opus-5[1m]

# cove lens — resolve pass (iteration 1)

Origin ledger = `findings.md` F-1…F-9. Every entry exits once.
Grounding gate applied per finding: code on `main`, `{RDR_RESOURCES}`
(`docs/cli-output-contract.md`, `CLAUDE.md`), and the RDR's own decided text.

| # | Type | Disposition | Section touched |
| --- | --- | --- | --- |
| F-1 | (c) | **dismissed-with-cite**, partial fix | Decision Rationale (joint-check) |
| F-2 | (b) | **fixed** | Decision Rationale (sibling-path check) |
| F-3 | (a) | **fixed** | A6 evidence |
| F-4 | (a) | **fixed** | Key Discoveries; References |
| F-5 | (d) | **dismissed-with-cite**, partial fix | Phase 2; Capability Dependencies |
| F-6 | (c) | **fixed** | A2 evidence |
| F-7 | (c) | **fixed** | Normative block 2; Identity; scenario 4 |
| F-8 | (d) | **fixed** | New normative block 5; A7; scenario 9; Approach 5–6 |
| F-9 | (d) | **fixed** | Normative block 1 |

## The two dismissals

**F-1 — "Joint-check: clear (7 peers) is false; a Final peer's Verified
assumption is hostage to this Draft's open locus."** Dismissed on grounding
source 1 (the peer text itself). RDR 0007 A13's conclusion is *"No constraint
exists anywhere"* — an accepted-exposure verdict built from a five-peer sweep
(0005, 0004, 0001, 0008), each item showing an **absence** of constraint. The
0008 sentence is the last corroborating item, cited to show the only
input-side obligation in the set is about a key *name* rather than provenance
origin. That reading is stable across all three Enforcement-locus candidates,
since (a)/(b)/(c) all constrain the same name. If block 4 changed or vanished,
A13's verdict would be unchanged or reinforced — the dependency direction the
finding asserts is inverted. Corroborating: 0007 lists this RDR under "Related
but distinct seams, **deliberately not folded in**" (`0007…md:158-163`) and
disclaims 0008 ownership of anything it depends on (`:811-815`).

*Partial fix kept*: the inbound reference was genuinely undisclosed, and A13's
"exactly one input-side obligation" claim is falsifiable by a *future* peer
adding a second. Both are now recorded on the joint-check line as a
cluster-reconcile watch item — disclosure, not a constraint on this RDR.

**F-5 — "RDR 0006 holds Final blocking-lint authority and is never
mentioned; the locus is unreconciled."** Dismissed on grounding source 1: 0006
draws the boundary itself, and 0002 already occupies the layer in question.
0006 "receives a normalized graph value, not Cobra command state and **not
sparse TOML**" (`0006…md:244-246`); takes "declared tags with provenance" as
already-valid *input* (`:251`); and leaves "table source and normalization …
in RDR 0002" (`:458-460`, also `:62-66`, `:781-785`). Its "MUST NOT define
different acceptance rules" clause governs *CLI surfaces re-implementing graph
acceptance* (hooks, aliases, resolver-local flags), not a normalizer's
data-level load validation. Decisive peer practice: RDR 0002 (Final) homes
`unknown tag` / `unknown context` / `unknown accessor` — declaration-name
validation — in its own data-level category block with **zero** references to
0006, and states "malformed TOML or schema-invalid rules are load/lint
failures before the resolver sees a table" (`0002…md:200-201`). So this RDR's
check sits in an established pre-acceptance layer, not in 0006's scope. The
richer-payload concern (0006 requires stable code + model identity + severity
+ span) applies to 0006 findings, which these are not.

*Partial fix kept*: the boundary was never *stated*. Phase 2 now names it with
cites, so a reader need not rediscover it.

## F-8 — the one finding that changed the design

Real and reachable, confirmed by reading source rather than the RDR:
`::missingOwned` iterates `row.RequiresOwned` and tests
`view.has(key, ProvenanceOwned)`; `::TagSet.has` compares
`tv.provenance == prov`. A key present under `ProvenanceRecognized` therefore
fails the owned-only test, and the kernel refuses `owned_state_unavailable`
naming the reserved key — the confusing-refusal shape this RDR exists to
eliminate, one field over. RDR 0002 says nothing about `RequiresOwned` (zero
hits), so no peer covered it.

**Tiebreaker collapsed rather than escalated.** The apparent fork was
"add the clause and collide with RDR 0007's declared ownership of
`RequiresOwned`, or omit it and leave the hole." The evidence dissolves it:
0007 owns what the field *means* (post-guard write-dependency keys; the
owned-only vs provenance-blind test distinction). This RDR owns the *name*
`recognized`. A name reservation is not a semantics change, and 0007's own
Metadata prescribes the composition rule — "both peers cite rather than
restate." So block 5 reserves the name, cites 0007 for meaning, and restates
none of it. Booked as A7 (`Pending`) with the ownership half routed to Stage
7.1 cluster-reconcile, since 0007's concurrence is a peer question rather than
a source-readable fact.

## Needs (re)verification (Stage 6 closes these)

- **A7 (new, `Pending`)** — two halves. *Reachability*: read at Pre-Lock from
  `::missingOwned` / `::TagSet.has`; to be pinned by test, not by reading
  (scenario 9). *Ownership*: whether a name reservation reaches RDR 0007's
  contract — reconcile at Stage 7.1.
- **No previously-Verified assumption was invalidated.** A2's fix was an
  arithmetic correction (five read statements across four sites, not four); its
  conclusion — one tag-vocabulary site — was independently re-confirmed at
  `resolve.go:151-156` and its `Verified` stamp stands.
- **Watch item, not an assumption**: RDR 0007 A13's negative-existential
  ("the set holds exactly one input-side producer obligation") is falsifiable
  by a future peer, not by this RDR. Carried to cluster-reconcile.

## Charted to successor

None. Every finding either landed in scope or was dismissed with a cite; no
net-new scope was absorbed.
