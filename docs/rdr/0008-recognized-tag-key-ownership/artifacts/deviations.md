# Deviations — RDR 0008 Ownership of the recognized-outcome tag key name

Pre-seeded by the `0002-0009` cluster gate, iteration 3 (2026-08-24 —
`docs/rdr/cluster-reconcile/0002-0009/iter-3/reconcile-report.md`). Stage 8
opens this file; running each named check IS the entry's disposition. An entry
escalates only if its check contradicts a contract.

---

## D1 — Scenario 3 assumes a view-capturing guard seam that RDR 0007's fence forbids

- **Type**: TEST-FIXTURE
- **Status**: OPEN (deferred at cluster gate; DEFER-TO-IMPLEMENTATION)
- **Source**: `iter-3/pairwise-0007-0008.md` F1 (iter-2 critique C-12, never
  in a reconcile findings table → NET-NEW at this gate).
- **Conditions carried**: (a) `0008:2390-2413` is an MVV scenario, outside
  every ```normative fence (nearest fences close at 1635/1647); (b) the check
  below is a compile/test; (c) 0008's contract (reserve the key; enforce at
  the kernel input boundary) is unchanged — only the test's capture point
  moves.
- **Contradiction**: scenario 3 expects "the captured guard-side view
  satisfies `Lookup("recognized") == …` and `Len()` equals the key count";
  0007 (fenced, `0007:1282-1290`): "`Evaluate(atom, value)` … It never sees
  the view."
- **Check** (Stage 8, Phase 1 test authoring): write scenario 3 against
  0007's per-atom seam, capturing the `value` handed to the guard atom over
  `recognized` and asserting it equals `(in.Recognized, ProvenanceRecognized)`
  while a `Match` on the same key selects the row in the same resolve. If the
  same-view property cannot be asserted at that seam, escalate as SPEC-DEFECT
  (0008 scenario 3 / A2 vs 0007 §D4 seam).

## D2 — A10 evidence rests on the pre-§D7 `[accessors.<id>]` `{Mode, Path}` layout

- **Type**: TEST-FIXTURE
- **Status**: OPEN (deferred at cluster gate; DEFER-TO-IMPLEMENTATION)
- **Source**: `iter-3/pairwise-0008-0002.md` N2; critique N-8.
- **Conditions carried**: (a) `0008:882-918` is assumption evidence, unfenced;
  (b) grep below; (c) A10's conclusion (no accessor-side site can mint the
  reserved key) survives via block 4 — no meaning change.
- **Check** (Stage 8): after 0002 re-locks under §JD-17, `grep -n
  'recognized' <0002 layout clause + rdr-fixture.toml>` shows no
  `[read.<id>]`/`[write.<id>]`/`[gate.<id>]` `keys` entry may name
  `recognized` (a writer key must be owned; `recognized` is kernel-supplied).
  If the re-locked layout admits it, escalate as SPEC-DEFECT against 0002.
- Also unfenced and stale after §D5: `0008:1772` "`<clear>` sentinel is
  syntactically un-declarable" — it is refused by load category; citation
  repair at 0008's next touch.
