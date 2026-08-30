# Charted to successor — repeatability lens, RDR 0025

- **Module layout for the command binding (D2, G2)** — the three runs produced
  three different package/file layouts (`internal/cli/cmdbind` sibling package;
  `internal/accessor/command_binding.go`; unnamed) and three different helper
  decompositions. Out of scope: RDR 0025 fixes the wire, the carrier grammar,
  the load-time vocabulary and the seam it implements — not where the code
  lives. C1 now says so explicitly rather than leaving the silence to be read
  as an unstated constraint. Successor: none needed; this is implementation
  latitude by design, and Stage 8 picks the layout.
