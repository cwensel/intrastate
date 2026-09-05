# Post-Mortem: RDR-0019 Owned-state initialization semantics

Ledger-only. Opened at Stage 6 reconcile to record one escaped defect: a
normative contract citing a resolver the implementing package cannot call.
The remaining post-mortem sections are authored later.

## Escaped-Defect Ledger

One row per finding, enriched with the two fields triage cannot assign. Rows
whose disposition is a record-level contract correction carry `n/a` odc
columns — there is no implementation defect to classify.

| Finding | odc-type | odc-trigger | Expected-catching stage | Precursor lens | Escape distance |
| --- | --- | --- | --- | --- | --- |
| `0019:C1`'s read-side carrier gate and `0019:A8` both named `internal/accessor/model.go::Registry.readerFor` as the resolver by which the verb learns a needed role's reader before invocation. That method is UNEXPORTED (lowercase, `model.go:249`) and its only non-test caller is `internal/accessor/executor.go:343`, inside the declaring package — so `internal/cli`, the verb's package, cannot call it. The contract named a check the verb could not perform. Repaired in place at reconcile after a strong consult returned PASS: the gate re-derives the reader over the exported `::Registry.Definitions` field, filtering `Identity.Capability == accessor.CapRead && Accessor.Role == role`, FIRST MATCH IN REGISTRY ORDER — which the consult established is exactly `readerFor`'s semantics (bare first-match, no `::bound` check, no normalization; loader stores `Role: *a.Role` verbatim). First-match had to be made normative rather than left to the implementer because read-side role uniqueness is not enforced today (no arm in `internal/accessor/validate.go` or the loader's `::accessorTable`; RDR 0016, which would make role→reader a function, is Draft), so a collect-and-refuse-on-ambiguity filter would be WIDER than the symbol it replaces. | n/a | n/a | `4-resolve` | **A7 is the precursor, and it is an exact-shape precedent, not a near miss** — A7 found this identical defect one package over (`flowbind::commandBacked` unexported, contract naming an uncallable discriminator) and C1 already carries its repair (type-switch on the exported `::Definition.Binding`) plus the sentence "Exporting `commandBacked` is not required and is not authorized here". The read-side SELECTION step was never subjected to the reachability question A7 had already answered for the write-side DETECTION step. Neither the critique lens (which opened A8 and corrected its pointer/value spelling) nor repeatability re-asked reachability. | 2 stages (owed at `4-resolve`, where A7's own reachability finding was verified; caught at `6-reconcile`) |

**Notes.**

- Root cause worth naming: a reachability finding was repaired at its
  discovery SITE rather than generalized into a check. A7 established that
  "this contract names a symbol the verb's package cannot reach" is a live
  defect class in this codebase's `internal/` layering, and C1 absorbed the
  specific repair — but nothing re-ran that question over the record's other
  cited symbols. The A8 citation was authored AFTER A7 was verified, into a
  record that already contained A7's warning.
- Cheap generalization for a future stage: for each `path::Symbol` anchor a
  normative contract cites as a check the implementation performs, confirm the
  symbol is exported or in the implementing package. This is a grep against
  the anchor list the grounding sweep already builds, not new analysis.
- The A6 companion assumption verified clean at the same gate with no route-back
  (`evidence/spikes/a6-emptiness-carrier.md`) and is not a ledger row.
