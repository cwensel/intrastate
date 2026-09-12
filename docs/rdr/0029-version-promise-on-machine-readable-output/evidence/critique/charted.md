# Charted — critique lens, cli/0029

Findings real but out of scope for this RDR's lock. Recorded here so the
disposition is durable rather than a silent drop.

- **D-19 (`0029:A1`) — a Draft amends an Implemented peer's contract and
  rewrites three of its tests as a prerequisite.** Out of scope for this
  lock because the record already routes it and the project's own precedent
  says the shape is in-bounds: `0029:§risks-and-mitigations` names the risk
  and sends it to 7.1 cluster reconcile over {0006, 0029}, and A1's dated
  Ruling (2026-09-11) cites `JDR 0001 §D10` rule 3 as precedent for a ruled
  citation-repair. Successor: the existing 7.1 cluster reconcile over
  {0006, 0029} — no new record needed. Pass A asked for "a separate record
  with its own gate"; that would duplicate the 7.1 pass already owed.

- **Presentation-only contract fold (from `0029:§proportionality`).**
  `rdr-write --outcome profile` emits `stopped:split-signal` because it
  counts four labelled `**Cn**` blocks (`contracts_durable=2+`). The split
  test is seams, not labels, and all four clauses bind one seam — the
  `--as=json` terminal envelope, via `respond::Success` and
  `clierr::EmitJSON` and no third surface. The remedy its `surface` line
  names — fold the clauses under one `**C1**` — changes no normative word
  and is a restructure of Normative Contracts, so it is not taken inside a
  lens pass. Successor: Stage 6 reconcile, or Stage 7's gate, whichever
  reaches the section first. The gate response records the reasoning.
