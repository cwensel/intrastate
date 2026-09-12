# Charted to successor — 3amigo, 0021

Real findings, out of scope for this RDR. Recorded durably here and
dismissed from this lens's loop per the resolve prompt's
charted-to-successor disposition. Never absorbed into the current draft.

Charted: `0021:MVV` / `0021:S6` (PM #11) — the diagram-review outcome is
validated only as DOT well-formedness plus node/edge id-set equality,
never as legibility or usability for a human reviewer. Out of scope: DOT
STYLING is explicitly non-normative in this RDR (Load-Bearing Decisions,
"Wire / byte format"; F3 scopes the normative fixture to the node/edge
SET and abstraction-marker placement), so a legibility criterion cannot
be asserted here without reversing that decision. Suggested successor: a
styling/layout RDR or a kata, if reviewer legibility ever becomes a
stated acceptance criterion rather than a downstream consumer's concern.

Charted: `0021:§decision-rationale` (PM #12) — the CI-diffability
outcome has no scenario that diffs one model's export across two
commits; determinism (C3, Scenario 1) and the golden pin (Scenario 2)
are proxies for it. Out of scope: a true cross-commit diff test needs
two builds of the tool, which C3's DELIBERATE NARROWING to (model,
build) places outside this RDR's promise — the cross-build golden named
in C3 is the acknowledged seam. Suggested successor: a CI-integration
kata that pins an export across a version bump, where the two-build
harness legitimately lives.
