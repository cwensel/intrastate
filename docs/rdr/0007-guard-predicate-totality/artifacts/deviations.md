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

---

## D7 — The block boundary needed a fail-closed rule the RDR states only for operators

- **Type**: SPEC-UNDER
- **Status**: CLOSED by Phase 3c. The evidence base resolves it; no author
  decision is owed.
- **The gap.** `0007:C3` states the drift rule for the OPERATOR boundary and
  makes it fail closed in both directions ("an existence atom carrying a
  foreign TOKEN is a value atom to the kernel"; "A foreign LITERAL … is
  UNEVALUABLE"). `0007:C6` fences the row verdict over `all` and `unless`
  only. Neither clause states what the kernel does with an atom in a THIRD
  block — and JDR 0001 §D12 made a third block representable in `Row.Guard`
  by putting `BlockMatch` on the same `Block` type ("one type, no separate
  slice on `Row`"). The zero value `Block("")` and any future token are
  equally representable.
- **What the implementation did before the fix.** `evaluateAtoms` routed
  `BlockUnless` explicitly and folded everything else into `all_result`, so
  the block boundary failed OPEN where the operator boundary fails closed.
  Phase 3a FAIL-1 and Phase 3b ADV-1/ADV-2 are the two consequences: a
  `BlockMatch` atom over an absent key refused the NON-escapable
  `guard_unevaluable` where `0007:F4` requires the escapable `no_match`; and
  a decided-FALSE atom in an unrecognized block pruned its row with its
  owned-state obligation (D8) while a modeled escape rescued `Escaped:true`
  — kata `xg7p`'s masking probe, this RDR's own originating defect,
  reproduced on the post-RDR kernel.
- **Evidence** (`{RDR_RESOURCES}` → design docs; the RDR's own fences):
  1. `0007:C6` (`docs/rdr/0007-guard-predicate-totality.md:1447-1449`) —
     "The row verdict is `all_result ∧ ¬(unless_conj)`, where `unless_conj`
     is the conjunction of the `unless` block's atoms." An atom in a third
     block is an operand of neither term. The formula is the fence, and it
     is exhaustive over what an operand can be.
  2. `0007:C5` supplies the shape for "contributes nothing": an omitted
     `unless` block "contributes no operand, not a neutral one". The same
     construction applies one level down, at the atom.
  3. `0007:C11` ratifies pruning "conditional on the domain rule: pruning is
     safe exactly because GuardFalse can only arise from decided atoms" — a
     condition stated over GUARD atoms. An atom in an unfenced block is not
     one, so folding it into `all_result` voids the condition D8 rests on.
  4. `0007:C3`'s posture generalizes: drift at a kernel boundary fails
     CLOSED. The operator boundary fails closed in both directions; a block
     boundary that failed open would be the one boundary in the kernel where
     an unrecognized token can DECIDE a row.
  5. JDR 0001 §D12 (`docs/jdr/0001-resolve-kernel-seam.md:640-652`) — "0007's
     per-atom payload may therefore carry `match`, which no refusal names —
     **harmless, since match atoms are never unevaluable**." §D12's own
     harmlessness claim only holds while a match atom contributes no
     operand; folding it in made it drive refusals, which is what §D12
     asserts cannot happen.
- **Resolution.** `evaluateAtoms` replaces the `default:` arm with explicit
  `BlockAll` and `BlockUnless` cases plus a fail-closed fallback
  (`isGuardBlock`) that contributes no operand to either conjunction, emits
  no payload entry, and never consults the seam. Q1 reading (a) is
  unchanged: `BlockMatch` remains a constant on the `Block` type and `Row`'s
  match pattern is untouched.
- **Why SPEC-UNDER and not SPEC-DEFECT.** No clause is wrong. The RDR states
  the fail-closed posture for one boundary and the verdict formula for the
  other, and the composition of the two — what happens to an atom the
  formula does not name — is unstated. §D12 opened the case after 0007
  locked.
- **New public surface**: none. `isGuardBlock` is unexported.
- **Escalation**: none.

---

## D8 — The payload sort key was total over atoms but not over rows

- **Type**: IMPL-DECISION
- **Status**: mechanical translation — the RDR states the rule; the
  implementation under-applied it. Recorded because the correction changes
  a comparator's contract.
- **The gap.** `0007:C8` fences the ordering as "the tuple `(RuleID,
  SourceLocator, key, block, operator token, literal)`, compared field by
  field in that order" and states "**The sort key MUST be TOTAL over payload
  entries.**" The implementation split that one tuple across two
  comparators: `compareUndecidedAtoms` over the atom half, and
  `compareUndecidedRows` over `(RuleID, SourceLocator)` alone. Since
  `slices.SortFunc` is not stable, two undecidable rows tying on identity
  retained `Table.Rows` position — reproduced by Phase 3a FAIL-2 on both the
  candidate set and the escape set.
- **Evidence** (the RDR's own fence,
  `docs/rdr/0007-guard-predicate-totality.md`, `0007:C8`): the clause names
  the failure mode it exists to close — "Two entries would then tie … leaving
  the order to an unstable tie-break — the atom-order dependence this clause
  exists to forbid, **one level below the row order ADV-3 already freezes**."
  It fences totality over ENTRIES, and REQ-55's "two entries equal on all six
  name the same atom" makes row identity a COMPONENT of the entry key rather
  than a partition boundary. Whether RDR 0002 guarantees `(RuleID,
  SourceLocator)` uniqueness is not stated in 0007, and `0007:C8` places the
  obligation on the kernel unconditionally.
- **Resolution.** `compareUndecidedRows` breaks an identity tie on the atom
  half, entry by entry over each row's already-sorted atom list, then by
  atom count. Two rows equal under it carry the same identity and the same
  entries, so their relative order is unobservable in the emitted payload —
  which is what "a function of the input tuple" requires. `Reason` is
  compared alongside each atom for the same reason: it is part of the emitted
  entry even though it is not one of REQ-55's six ordering fields.
- **New public surface**: none.
- **Escalation**: none.

---

## D9 — A duplicated tag key resolves to UNEVALUABLE, not to a positional value

- **Type**: SPEC-UNDER
- **Status**: CLOSED by Phase 3c. The evidence base resolves the choice; no
  author decision is owed.
- **The gap.** `0007:C8` requires the disposition to be "a function of the
  input tuple (RDR 0001 REQ-1), **never of atom or row order**". It closes
  atom order and row order. It does not reach the third ordering the kernel
  is exposed to — the order of the caller's TAG slices. `resolve.go::assemble`
  resolved a key repeated WITHIN one provenance by "last tag wins", and A13
  leaves caller-supplied `Observed` unconstrained, so two orderings of ONE
  duplicated observed key assembled different values. Phase 3b ADV-3
  observed `[gate=open, gate=closed]` → `Escaped:true` plan routing around
  an absent owned key, and `[gate=closed, gate=open]` →
  `owned_state_unavailable`: same tag multiset, opposite disposition. Under
  D8 that is the difference between a pruned row and a survivor, and
  therefore between the masking plan and the honest refusal.
- **Scope ruling.** `assemble` is RDR 0001 code, but `0007:C8` is 0007's own
  contract and its guarantee is exactly what a positional duplicate-key
  resolution breaks. Shipping as-is ships a `0007:C8` that does not hold.
  Fixed at the kernel boundary only: RDR 0001's merge semantics for the
  NON-duplicate case, and its REQ-3 replay determinism, are byte-for-byte
  unchanged (verified: identical repeats are not conflicts, cross-provenance
  precedence is untouched, and a permuted non-duplicate input assembles an
  identical view).
- **The two candidate rules, and why unevaluable wins.** `0007:F6`'s own
  suggested resolutions are "reject it as a malformed input tuple, or make
  the collision itself unevaluable."
  1. **Rejection is refused.** RDR 0001 reserves the Go error return for
     programmer mistakes: `Resolve`'s contract is "A modeled refusal travels
     the `Result` value with a nil error; the error return is reserved for
     programmer mistakes, not for modeled refusals." A duplicated
     caller-supplied observed tag is caller INPUT that A13 explicitly leaves
     unconstrained, not a programmer mistake, so routing it to the error
     channel would put unconstrained caller input on the panic-adjacent
     path. A new refusal kind is equally refused — REQ-7 pins the kind set at
     exactly five, and `0007:C7` already rejects reopening that taxonomy for
     a strictly better-motivated case (combined reporting).
  2. **Unevaluable is the RDR's own posture, stated four times.**
     `0007:C2` — a value-comparing atom that cannot be decided "MUST evaluate
     to unevaluable — never to false and never to true". The Research
     Findings' in-repo house rule, RDR 0004 A8: "Indeterminate MUST be a
     refusal-class result, not a false allow and not a false deny." The
     external anchors carry the same shape: SCXML §5.9.1 folds unevaluable
     into false only with a mandated observable error on a second channel,
     and the kernel has one channel; SQL:2003 strict routines make the HOST
     refuse to invoke the routine on an argument outside its domain rather
     than let the routine invent an answer — which is precisely "do not hand
     the seam one of the colliding values".
  3. `0007:C8`'s reason set supplies the exact word: `uncomparable` is "the
     key was present and its value was not compared to a verdict". A
     conflicted key is present and its value was not compared. No third
     reason is minted, and the closed two-member set holds.
- **Resolution.** `assemble` marks a key repeated within one provenance with
  DIFFERING values as conflicted. Presence stays provenance-blind and
  non-positional (`0007:C4`), so an existence atom still decides from
  presence alone — a collision is about the VALUE, not about whether the key
  arrived. A value-comparing atom over a conflicted key returns
  `GuardUnevaluable` / `ReasonUncomparable` and the seam is not consulted.
  Both orderings then agree, and neither can mask absent owned state: the
  row stays a SURVIVOR (`0007:C10`), so its `RequiresOwned` obligation is
  still raised.
- **Two narrowings, both deliberate.** A repeat carrying the SAME value is
  not a conflict — either resolution is the same value, so the result is
  already a function of the tuple. A key crossing PROVENANCES is resolved by
  the owned > observed > recognized precedence, which is a property of the
  tuple and not of slice order, and is left untouched.
- **New public surface**: none. `taggedValue.conflicted` and
  `TagSet.conflicting` / `TagSet.merge` are unexported; `TagSet.Lookup`'s
  exported three-value shape, fenced by `0007:C1`'s presence test, is
  unchanged.
- **Why SPEC-UNDER and not SPEC-DEFECT.** `0007:C8` is not wrong; it
  enumerates the orderings it had in view (atom, row) and A13 admits a third
  the clause never names. Per the BUILD-ORDER standing rule this is
  fix-now, not rdr-seed.
- **Escalation**: none.
