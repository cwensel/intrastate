Model: claude-opus-5[1m] (dispatcher; persona passes stamped in their own files)

# 3amigo iter-2 (delta) — consolidation — cli/0012

Delta re-run after iteration 1's substantial rewrite, scoped to the
still-open anchors. Three isolated personas, no cross-persona
visibility.

## Hotspot

| id | personas | what converged |
| --- | --- | --- |
| `0012:C3` | 2 (IMPL-D1, QA-D1) | the rewrite's own call-site census was wrong by one, and the missed site drives a ZERO-VALUE evaluator through the suite — directly contradicting the scenario the rewrite had just added |

One hotspot, down from seven in iteration 1.

## Findings (8) and dispositions

### High
- **QA-D1** `0012:C3`,`0012:S5` — `internal/guard/guard_evaluator_0003_test.go:67`
  drives `resolve.TestGuardEvaluatorContract(t, ev)` over
  `var ev guard.Evaluator`. The suite asserts GuardTrue for `eq`/`in`
  on exactly the zero value scenario 5 requires to answer
  GuardUnevaluable. Verified directly. **fixed** — C3 now says the
  site is RE-POINTED (not that a new caller is added), names all three
  call sites, and S5 states the conflict and its resolution.

### Medium
- **IMPL-D1** `0012:C3` — same census error, reached independently.
  **fixed** with QA-D1.
- **IMPL-D2** `0012:C3` — the widened signature breaks
  `guard_mvv_test.go:222`'s `var fn func(*testing.T, resolve.GuardEvaluator)`,
  a compile-time pin whose doc marks name AND signature as normative
  BOUNDARY surface for `0007:REQ-71`. The rewrite changed it silently.
  **fixed** — C3 now amends `0007:REQ-71` explicitly, preserving the
  importable-surface guarantee and recording the deviation the way A2
  does for `0003:C9`.
- **PM-D1** `0012:§consequences` — the Negative bullet closed "A4
  bounds the blast", but A4 is now Pending AND only ever measured the
  authored side, while the bullet's own `"07" eq 7` example is a
  held-side flip (A6's unmeasured population). **fixed** — the two
  legs are now bounded separately and the silent leg is named
  unmeasured.
- **QA-D2/D3/D4** — three dangling cross-references from the scenario
  renumber. **fixed**: MVV step 1's "S8 and S8" → "S8 and S9";
  Phase 4's "Lands with S8's regression" → S7 (S8 is the unchanged CLI
  control, not a regression). The third (S5→S6 in C5/F1/authority) had
  already been corrected in the same pass that introduced it.

### Low
- **IMPL-D3** `0012:C4`,`0012:§phase-3-producer-alignment` — the
  `guardSeam` arity change hits four `internal/cli` test call sites.
  **fixed** — Phase 3 now names all four with file:line.
- **IMPL-D4** `0012:C1` — `DeclaredKinds` takes `*table.Model` while
  the package's `CanRefuse` takes a decls map. **fixed** — C1 states
  why the model form is chosen.
- **PM-D2** `0012:A6` — the silent prune→match leg appeared in no
  PM-owned section. **fixed** — the MVV acceptance paragraph now names
  it and points at A6.

Also found and fixed beyond the personas' reports: Phase 1's
literal count was eleven composite `Evaluator{}` forms but missed
seven `var ev guard.Evaluator` declarations in four files (two never
named), so the phase as scoped could not reach a green `make check`.
Now eighteen sites, enumerated.

## Verification quality

The implementer pass checked 15 claims the rewrite introduced against
`main`; 13 held exactly, 2 did not (IMPL-D1, IMPL-D2) and are the two
fixed above. The QA pass audited all 24 S-number citations after the
renumber. Both are grounded, not speculative.

## Resolved by the rewrite (not re-filed)

All six iteration-1 QA findings, all eight iteration-1 implementer
findings, and three of four PM findings. C5's narrowing is confirmed
sound from the product view: the owned-ingress fix never depended on
C5, so narrowing it does not weaken the delivered outcome.
