# Deviations — RDR 0007 Guard predicate totality over an incomplete evaluation view

Opened 2026-08-24 to record the joint decisions JDR 0001 answered after the
`0002-0009` cluster gate, iteration 4
(`docs/rdr/cluster-reconcile/0002-0009/iter-4/reconcile-report.md`). Stage 8
opens this file; running each named check IS the entry's disposition. An entry
escalates only if its check contradicts a contract.

**Run order: this RDR runs FIRST (1 of 8)** — it defines the
`internal/resolve::Row` that RDR 0002 normalizes to, so it precedes 0002
despite 0002 owning the wire format. §D12 (`BlockMatch`) lands here. Phase 3's
`contains` leg cites §D13's JSON-array form directly — 0002 writes that clause
in run 2 and must match it, not redefine it. See
[`../../BUILD-ORDER.md`](../../BUILD-ORDER.md).

---

## D1 — JD-8/JD-21/JD-22 answered: two landings and one citation

- **Type**: TEST-FIXTURE
- **Status**: OPEN (pending 0007's next touch)
- **Source**: JDR 0001 §D10/§D12/§D13 (`d937eec`, settled `5c2b96b`, both
  2026-08-24) — answered *after* the iteration-4 gate recorded them
  unanswered. Status qualifier and README row corrected 2026-08-24.
  **§JD-18 remains open** (conforming-view enforcer, sibling 0003) and stays
  in the qualifier.
- **§D12 (§JD-21) — a landing in this RDR.** One Go type: `Block` gains
  `BlockMatch`, **added in this RDR's Phase 1**. The fence at `0007:1264-1266`
  ("exactly two constants") is the clause §D12 overrides; Phase 1's reshape is
  the named site. This is a normative change to a Final RDR's fenced
  cardinality — it cannot be done as a citation. Escalate at Stage 8 and take
  the scoped re-entry rather than widening the type silently.
- **§D13 (§JD-22) — unblocks this RDR.** The set-value encoding is RDR 0002's
  canonical sorted, duplicate-free JSON array in `Tag.Value`. This RDR's
  `contains` contract-test leg (`0007:2159-2170`, "Blocked on one RDR 0003
  declaration") and Testing Strategy scenario 8's present-key half were
  blocked on exactly this and are now writable. Note §D13 assigns the
  declaration to **0002**, not to 0003 as `0007:2182` states — that sentence
  is a citation repair.
- **§D10 (§JD-8) — citation.** This RDR's A19 is **accepted**: the envelope
  carries one `omitempty` `findings` field, and §D4's routing is reversed.
- **Check** (Stage 8): `grep -n '§D1[023]' docs/rdr/0007-*.md` → ≥1 each;
  `BlockMatch` exists in the Phase 1 type; the `contains` contract-test leg is
  written against 0002's JSON-array form; `0007:2182` names 0002 as the
  declarer. Escalate the §D12 fence change as above.

---

## D2 — Phase 1 test-author decisions taken unattended

- **Type**: TEST-FIXTURE / NAMING
- **Status**: OPEN (for Phase 2 to confirm or correct in one line each)
- **Source**: Stage 8 Phase 1 (spec tests, red). No human in the loop; each
  item below would normally be a question and is recorded instead.

1. **Atom type name `GuardAtom`.** `0007:C1` names the atom's four FIELDS
   normatively but not the type. `0007:347` and A26 (`0007:1111`) both name
   `GuardAtom` as the proposed exported surface, and A26 clears it against
   RDR 0001's frozen boundary-symbol tests. The tests are written against
   `resolve.GuardAtom`. If Phase 2 spells it otherwise, this is a
   mechanical rename across the four new test files.

2. **`Row.Guard` keeps its NAME, changes its TYPE.** The Technical Design
   says "`Row.Guard` becomes an atom slice" (`0007:1244`), so the field is
   `Guard []GuardAtom`, not a renamed field beside a deleted one.
   `TestReq69_TheGuardTextFieldIsActuallyRemoved` asserts the type changed
   rather than that the name went away.

3. **`BlockMatch = "match"`** — the byte value §D12 does not spell.
   Recorded as an ASSUMPTION in `req-list.md`; asserted verbatim by
   `TestReq6_BlockIsAnExportedNamedStringTypeWithThreeConstants`. No 0007
   clause reads the `match` block's bytes, so a later respelling is a
   one-line change. Q1 reading (a) is honoured throughout: `Row`'s match
   pattern is unchanged, and `TestReq78` asserts that as a NEGATIVE
   contract.

4. **`TestGuardEvaluatorContract` is inspected, not failed.** REQ-71's
   behavioural obligation ("unparseable value → unevaluable, never false")
   would ideally be checked by running the contract test against a
   non-conforming seam and asserting it FAILS. Go's `testing` propagates a
   subtest's failure with no supported suppression, so a deliberately
   failing run would fail the file. The test instead drives the contract
   test with a CONFORMING seam wrapped in a recorder and asserts on the
   product it exercised (≥2 operators, ≥2 values against one literal, ≥1
   unparseable value), plus a satisfiability leg. Verified to reject both
   a hollow and a one-case contract test.

5. **Frozen-suite migration performed here, not deferred.** REQ-68 assigns
   fixture migration to Phase 1. Four frozen files break at the COMPILE
   surface and were adapted mechanically; `fixtures_test.go::namedGuard`
   renders each former guard string as one atom over the PRESENT key
   `status` so the seam decides it (never an absence). Fixup-1d is the one
   RE-DECIDED case (REQ-17). See `coverage.md` for the per-file table. No
   behavioural assertion outside 0007's surface was weakened.

6. **Escape-row payload (Q2) asserted, not just the kind.** Per the
   `disposition` mini-check and REQ-40's delegation reading,
   `TestReq41_EscapeSetScopingBothLegs` asserts the escape set's payload
   contents. Testing Strategy row 14 pins only the kind; this is the
   stronger of the two readings and the one `req-list.md` Q2 proceeds on.

- **Check** (Phase 2): each item is confirmed by the implementation
  compiling against the tests unchanged, or corrected by a one-line edit
  in the named place.
