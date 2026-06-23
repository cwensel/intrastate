# Persona 1 - Product Manager

Question: does this RDR actually deliver the user outcome?

Findings:

1. `Proposed Solution / Approach` and `Technical Design` blur the user outcome of replay-safe resolution by saying the kernel applies accessors to artifacts, while other passages say the kernel consumes an owned-state snapshot and returns values. The outcome is unclear because "same input produces same output" depends on whether accessor reads are inside or outside the resolver boundary.

2. `Validation / Testing Strategy` validates no-match, ambiguous-match, and unavailable/unevaluable refusals, but the problem statement and normative contract also require refusing an unmodeled recognized outcome. The user outcome includes "unmodeled matches must be refused", so the test matrix should name that pass/fail case explicitly.
