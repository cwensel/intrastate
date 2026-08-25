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

---

## D3 — D1's Stage-8 checks run; all four pass

- **Type**: TEST-FIXTURE
- **Status**: CLOSED by Phase 2. Running each named check IS the entry's
  disposition, and none contradicts a contract.
- **Check results** (worktree `worktree-rdr-0007`, commit `d63950c`):
  - `grep -c '§D1[023]' docs/rdr/0007-guard-predicate-totality.md` → 1 each
    for `§D10`, `§D12`, `§D13`. ≥1 each: **pass**.
  - `BlockMatch` exists in the Phase 1 type: `internal/resolve/guard.go:24`
    declares `BlockMatch Block = "match"` as a third constant on the SAME
    `Block` type, with no separate slice on `Row`. **pass**.
  - The `contains` contract-test leg is written against 0002's JSON-array
    form: `guardcontract.go`'s `in` and `contains` cases carry
    `["alpha","beta"]` / `["alpha"]` — sorted, duplicate-free, compact —
    and `guard_fixtures_test.go::d13Set` renders the same form for the
    kernel tests. **pass**.
  - `0007:2182` names 0002 as the declarer: the citation repair stands in
    this file and in `req-list.md`'s standing correction 2; the RDR
    markdown is not amended (we never amend RDRs). **pass**.
- **The §D12 fence change was NOT taken silently.** It landed as the scoped
  normative re-entry D1 describes: one type, one constant, no widening of
  the K3 verdict formula `0007:C6` fences. Per Q1 reading (a),
  `BlockMatch` atoms are not evaluated by the guard pipeline —
  `evaluateAtoms` folds any non-`unless` block into the conjunctive
  reading, and `TestReq78` holds `TagSet.matches` unchanged.

---

## D4 — D2's six Phase-1 test-author decisions all confirmed

- **Type**: TEST-FIXTURE / NAMING
- **Status**: CLOSED. The implementation compiles against the Phase 1
  tests unchanged; no test file was edited in Phase 2. Item by item:
  1. **`GuardAtom`** — adopted verbatim. No rename.
  2. **`Row.Guard` keeps its name, changes its type** — adopted:
     `Guard []GuardAtom` replaces the field in place.
  3. **`BlockMatch = "match"`** — adopted. No 0007 clause reads the block's
     bytes, so the assumption stays a one-line change if §D12's home later
     spells it otherwise.
  4. **`TestGuardEvaluatorContract` inspected, not failed** — the exported
     harness satisfies both halves as written: 4 operators (`eq`, `gte`,
     `in`, `contains`) over 12 cases, two values against one literal on
     both `eq` and `gte`, and two unparseable-value legs
     (`gte`/`"many"`, `contains`/`"alpha"`) whose obligation is
     `unevaluable, never false`.
  5. **Frozen-suite migration** — confirmed sound. The full frozen RDR 0001
     suite passes against the implementation, `TestFixup1d` (RE-DECIDED)
     and `TestReq33` (payload plumbing) included.
  6. **Escape-row payload (Q2)** — confirmed. `escapeOrRefuse` delegates to
     the same `gate`, which populates `Undecided`, so an unevaluable escape
     row's payload rides along. No parallel implementation was written.

---

## D5 — A foreign `OpExists` literal reports `uncomparable` even when the key is ABSENT

- **Type**: IMPL-DECISION
- **Status**: mechanical translation — the RDR states the rule
  unconditionally; recorded because it affects how a reader reconciles two
  clauses.
- **The apparent tension.** `0007:C8` (REQ-45) glosses the reason set as
  `absent` = "the key was not in the view" and `uncomparable` = "the key
  was PRESENT and its value was not compared to a verdict". Read alone,
  that gloss says an atom over an absent key is always `absent`. But
  `0007:C3` (REQ-25) states the foreign-literal rule with no presence
  condition at all.
- **Evidence** (`{RDR_RESOURCES}` → design docs; the RDR's own normative
  fence, read at `docs/rdr/0007-guard-predicate-totality.md:1371-1381`):
  *"A foreign LITERAL on an `OpExists` atom — any value that is neither
  `LiteralTrue` nor `LiteralFalse`, the empty literal included — is
  UNEVALUABLE at the kernel, reason `uncomparable`; the kernel MUST NOT
  decide such an atom from presence."* The clause is unconditional, and
  the phrase "MUST NOT decide such an atom from presence" is precisely a
  direction to stop consulting presence for these atoms — so presence
  cannot select the reason either. Phase 1 asserts the same reading
  independently: `TestReq25_ForeignLiteralOnAnExistsAtomIsKernelUncomparable`
  runs each foreign literal over BOTH `reviews` (present) and
  `iterations` (absent) and requires `uncomparable` in both.
- **Resolution.** `evaluateAtom` returns `ReasonUncomparable` for a foreign
  `OpExists` literal regardless of presence. `0007:C8`'s gloss is a
  description of the two reasons' ordinary provenance, not a fenced
  predicate; `0007:C3`'s rule is the fenced one and is more specific. The
  substantive reading is consistent: what failed is the atom's own
  literal, not the view, so `absent` would name the wrong fact — exactly
  the argument REQ-15 makes for the nil-seam case ("its key is present, so
  `absent` would be a lie"), applied to the other direction.
- **Downstream note.** RDR 0005's renderer reads `Reason` as an opaque
  closed-set value, so this affects diagnosis wording only, never a
  verdict: every foreign-literal atom is UNEVALUABLE under both readings.
- **Escalation**: none. The evidence base resolves it. **Status:
  mechanical translation.**

---

## D6 — No unnamed public surface was added

- **Type**: IMPL-DECISION
- **Status**: informational; recorded so the ADDITIVE-IS-NOT-EXEMPT gate
  has a written disposition rather than silence.
- Every exported symbol this build adds is named by a REQ in
  `req-list.md`: `Block` / `BlockAll` / `BlockUnless` (REQ-6),
  `BlockMatch` (REQ-6 as widened by §D12), `GuardAtom` and its four fields
  (REQ-1, REQ-5), `OpExists` / `LiteralTrue` / `LiteralFalse` (REQ-22,
  REQ-23), `Reason` / `ReasonAbsent` / `ReasonUncomparable` / `Reasons()`
  (REQ-49), `UndecidedRow` / `UndecidedAtom` / `Refusal.Undecided`
  (REQ-48), and the narrowed `GuardEvaluator.Evaluate(GuardAtom, string)`
  (REQ-10). `TestGuardEvaluatorContract` is the exception the launch
  prompt already grants, named by REQ-71 as normative surface.
- Everything else the build adds is unexported: `evaluateAtoms`,
  `evaluateAtom`, `appendAtom`, `kleeneAnd`, `kleeneNot`, `boolResult`,
  `compareUndecidedAtoms`, `compareUndecidedRows`. The retired
  `evaluateGuard` and `gate`'s `slices.MinFunc` call are removed, not
  re-typed (REQ-52).
- No **SPEC-UNDER** entry is owed.
