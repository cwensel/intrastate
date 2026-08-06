Model: gpt-5

# 3-amigo — Implementer

1. **IMP-1 — Critical: the locator-agreement implementation boundary is
   unresolved.** Trigger: `Critical Assumptions / A4`, `Technical Design`, and
   `Implementation Plan / Prerequisites`. `SourceLocator` is opaque to `plan`,
   while an unnamed normalized-table boundary must validate that it designates
   the same model/rule. **Blocks:** package/API placement and whether RDR 0008
   or RDR 0002 supplies the validator. This re-raises COVE-3/A4.
2. **IMP-2 — High: ownership of the canonical identity type is unclear.**
   Trigger: `Proposed Solution / Approach` and `Capability Dependencies`. The
   normalizer supplies RDR 0002 identity while this RDR introduces
   `SelectedRule`, without saying whether that is the normalized canonical
   type or a resolver projection. **Blocks:** package placement, imports,
   constructor ownership, and avoiding duplicate structs.
3. **IMP-3 — High: invariant-error taxonomy is unspecified.** Trigger:
   `Technical Design`, `Normative Contracts`, and Validation scenario 4.
   "Programmer/invariant error" selects the Go error channel but names no
   sentinel, typed error, wrapping contract, or stable classification.
   **Blocks:** the new `plan` error implementation and assertions stronger than
   `err != nil`.
4. **IMP-4 — High: copy-isolation validation describes the wrong slice
   mutation.** Trigger: `Minimum Viable Validation` and Validation scenario 3.
   Replacing `edge.Action.Writes` changes only the slice header and can pass
   while the result still aliases the old backing array. **Blocks:** a test
   that actually proves defensive copying.
5. **IMP-5 — Medium: the MVV does not require a non-empty expansion suffix.**
   Trigger: `Technical Design`, `Load-Bearing Decisions / Identity`, and the
   MVV. Both fixtures could use the permitted empty suffix, allowing an
   implementation that drops suffixes to pass. **Blocks:** fixture selection
   that proves the complete identity handoff.

Lower-severity overflow: 2 — nil-versus-empty `Writes` representation; migration
expectations for existing exported `Edge` literals.
