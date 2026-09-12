# 3amigo consolidation — 0021, iteration 3 (pass 2, delta-scoped)

Model: claude-opus-5 (consolidator; persona passes all `claude-sonnet-5`)
Date: 2026-09-12
Iteration: 3 (the pass-1 loop verdict was `rerun`; `iter-2/` holds only a
post-fix lint — pass 2's lens work was never written, so this iteration owes it)

Three isolated persona passes, no cross-persona visibility. Consolidation is
mechanical: hotspots are `rdr anchors --record 0021` over each persona file,
counted, `>=2`. Overlap marks a hotspot PASSAGE, not a validated finding.

## Scope

Delta-scoped to what the pass-1 rewrite touched: `C1`, `C2`, `C4`, `C5`,
`S1`, `S5`, `S6`, `S8`, `A7` (flipped Verified), `A8` (new), the added
`Determinacy:` line, and `MVV`. Pass-1 findings were NOT re-raised; each
persona brief carried the settled/charted list.

## Isolation check

No persona file cites another persona's output (grepped). Anchor sets are
near-disjoint along ownership. Isolation held.

## Hotspots (cited by >=2 isolated personas)

| id | Personas | What converged |
| --- | --- | --- |
| `0021:C4` | PM, implementer | PM cleared it (ceiling arm has S7 coverage); implementer found the refusal-authority sentence factually backwards. One real defect. |
| `0021:C2` | PM, QA | Both cleared it — PM on outcome coherence, QA on spellings/orders being pinned by S2's golden. No defect. |
| `0021:S8` | PM, QA | Both cleared it. No defect. |

Hotspot ≠ defect: only `C4` carries one.

## Merged ledger

`P1`=PM, `P2`=implementer, `P3`=QA.

| # | id | Persona | Finding | Severity |
| --- | --- | --- | --- | --- |
| 1 | `0021:S6` | P3 | The rewritten oracle says each hostile label must "DECODE back to its source tag value" but names no decoder, grammar, or tool — and the `§pre-lock-mini-checks` `oracle`/`trace` rows for step 3 still assert the old exit-0-plus-set-match bar S6 itself now calls insufficient. Scenario and its own mini-check backing contradict each other. **HOTSPOT passage: no.** | high |
| 2 | `0021:C4` | P2 | The closing sentence names `analysis.go::checkNodeCeiling` "a distinct guard-product bound". It is neither distinct nor guard-product: it fires on the same `reach()` completeness using the same `nodeCeiling` constant. **HOTSPOT.** | medium |

PM: zero anchored findings (valid result). Widened to `S4`/`S7`/`S8` to
confirm C4's new ceiling arm has acceptance coverage — it does (`S7`), so
not a gap.

## Dispositions

Grounded against code on `main`, `{RDR_RESOURCES}`, and the RDR's own
decided text (incl. `rulings.md` F3) before any edit.

| # | id | Disposition | Basis |
| --- | --- | --- | --- |
| 1 | `0021:S6` | **fixed** | CONFIRMED the record names no decoder anywhere (grepped `decode`/`unescap`/`-Tplain`/`graphviz`: only JSON round-trip hits). The mechanism does exist — A6's spike records a complete escaping ORDER (backslash → quote → literal-newline → UTF-8 passthrough) and its Evidence already says that order "carries into implementation"; the `fidelity` mini-check row cites it too. Pinning a NEW normative DOT escaping grammar was rejected as a re-raise: ruling F3 scoped the normative fixture to the node/edge SET and marker placement, "not styling", and called the separator bug styling-side. So S6's inverse is now stated as inverting A6's recorded order, asserted over node/edge IDENTIFIERS only (the set F3 made normative), with styling explicitly left out of the oracle's scope. Both stale mini-check rows (`oracle` row 3, `trace` row 3) swept to match (§amendment-sweep). |
| 2 | `0021:C4` | **fixed** | CONFIRMED from source: `nodeCeiling = 4096` (`taxonomy.go:136`), published by `NodeCeiling()` (`taxonomy.go:152`); `reach.go:131` compares `len(nodes) > nodeCeiling`; `analysis.go::checkNodeCeiling` (`analysis.go:177`, via `engine.go:71`) fires on the same `reach()` completeness and the same constant. The genuine guard-product bound is `graphlint.ProductBound()` → `guard.Bound()` (`taxonomy.go:147`, `guard/product.go:424`), which `checkNodeCeiling` never calls. Sentence rewritten: same bound at two call sites (which strengthens C4's neutrality claim rather than excepting it), and `ProductBound()` named as the bound the refusal does NOT involve. |

## Needs (re)verification — carried to Stage 6

- **A8 (Pending, unchanged)** — the exported completeness surface. Untouched
  by this pass; the C4 edit corrected a grounding statement about which bound
  is at issue, and did not alter what A8 claims.
- **No new assumption raised.** Neither fix adds a load-bearing claim: #2
  replaced a false grounding statement with one verified from source this
  pass; #1 cites an escaping order already recorded and already carried by
  A6's Evidence and the `fidelity` row.

## Loop verdict

See the close packet: `rdr-loop.toml --outcome lens-loop`, tags recorded there.
