# Deviations — RDR 0003 Guard predicate exhaustiveness

Pre-seeded by the `0002-0009` cluster gate, iteration 4 (2026-08-24 —
`docs/rdr/cluster-reconcile/0002-0009/iter-4/reconcile-report.md`). Stage 8
opens this file; running each named check IS the entry's disposition. An entry
escalates only if its check contradicts a contract.

**Run order: 3 of 8** — consumes 0007's atom shape and 0002's normalized model,
and cites 0002's §D13 encoding. Owns the tag declaration model that 0006 reads
in run 4. Instantiates `resolve.TestGuardEvaluatorContract` (0007 Phase 3). See
[`../../BUILD-ORDER.md`](../../BUILD-ORDER.md).

---

## D1 — Empty / omitted `unless` identity is stated by 0007, not here

- **Type**: TEST-FIXTURE
- **Status**: OPEN (deferred at cluster gate; DEFER-TO-IMPLEMENTATION)
- **Source**: `iter-4/pairwise-0007-0003.md` F1; critique C-16 (open since
  iteration 2).
- **Conditions carried**: (a) 0003's fence (`0003:1143-1170`) is scoped to
  decided atoms and is *silent* on a group with no `unless`; the identity
  (`unless = ∅` subtracts nothing) is fenced in 0007 (`0007:1395-1409`),
  and 0003's prose at `741-747` is unfenced; (b) check below; (c) filling
  the silence with 0007's identity changes no 0003 clause.
- **Check** (Stage 8): 0003 MVV Scenario 2 run over a group with no `unless`
  block yields the full product (no subtraction); 0007 row 10 is the kernel-
  side twin. If Scenario 2 cannot pass without a fenced change, escalate.

## D2 — Stale peer claims (citation repair)

- **Type**: TEST-FIXTURE
- **Status**: OPEN (citation repair pending 0003's next touch)
- Sites: `0003:2093-2095` MVV Scenario 4 "a predicate semantic kind that RDR
  0006 can map to a lint finding" vs this RDR's own fence `0003:920-925` and
  JDR 0001 §D7(iii) ("RDR 0006 mints nothing"); `0003:583-586` A20 "RDR 0007
  states no atom carrier" vs `0007:1522-1536` `UndecidedAtom`; `0003:1530`
  "RDR 0002 | Pending" vs `0003:1539`; 0006 as lint-finding enveloper at
  `0003:1481, 1364, 1354, 1567, 1545` vs `0005:447`.
- **Check**: `grep -n '0006 can map' docs/rdr/0003-*.md` → 0.

## D3 — §JD-22 answered: cite §D13 for the set-value encoding

- **Type**: TEST-FIXTURE
- **Status**: OPEN (citation repair pending 0003's next touch)
- **Source**: JDR 0001 §D13 (`d937eec`, settled `5c2b96b`, both 2026-08-24) —
  answered *after* the `0002-0009` iteration-4 gate recorded §JD-22 as
  unanswered. Status qualifier and README row corrected 2026-08-24 (§JD-22
  dropped; §JD-18 retained).
- **The answer**: `Tag.Value` stays `string`; a set crosses the kernel seam as
  its canonical JSON array — members sorted, duplicate-free, compact encoding —
  and **RDR 0002 declares it**. Read-back equality is byte equality; no third
  encoding exists.
- **Scoped answer-vs-fences check: CONSISTENT.** §D13 governs *carriage* of a
  set tag value across `Tag.Value`; this RDR owns the set **literal** spelling
  (A13, `0003:1241-1248`) and the **element universe** (A9) — distinct things,
  and §D13 says so ("0003 … declares spelling and universe only"; "Encoding is
  carriage, not declaration semantics"). No fenced clause here contradicts it:
  A13's unordered/duplicate-free literal and §D13's sorted/duplicate-free
  carriage share one normal form, and `grep -n 'Tag\.Value'` over this RDR
  returns no fenced tag-value byte claim. The silence §D13 fills is real
  silence, so §D13 *narrows* what 0007 had assigned here rather than
  contradicting a clause.
- **Effect on peers**: 0007's `contains` contract-test leg (`0007:2159-2170`,
  "Blocked on one RDR 0003 declaration") and Testing Strategy scenario 8's
  present-key half are unblocked — the encoding is 0002's, and it is stated.
- **Check** (Stage 8): A13's Evidence and the References list cite JDR 0001
  §D13, noting the tag-value encoding is RDR 0002's while the literal spelling
  stays A13's. `grep -n '§D13' docs/rdr/0003-*.md` → ≥1. If the citation cannot
  be added without a fenced change, escalate.

---

# Stage 8 implementation — dispositions and new entries

## D1 — DISPOSED

- **Check run**: REQ-129's D1 leg asserts the MVV RDR slice's
  `stage.eq=propose` group carries no `unless` atom and that its coverage
  union equals its scoped product. `TestReq129_…` passes.
- **Outcome**: an omitted `unless` block subtracts nothing, exactly as
  RDR 0007's fenced identity states. No 0003 clause needed changing, so
  condition (c) holds and the entry closes.
- **Status**: CLOSED — check passed.

## D2 — OPEN (unchanged)

- **Check run**: `grep -c '0006 can map' docs/rdr/0003-*.md` → **1**, and
  the check expects 0.
- **Outcome**: the stale citation is still in the record. The repair is an
  edit to a locked RDR, which this stage does not make (project rule: RDRs
  are never amended). Nothing in the implementation depends on it — the
  code routes both rejections to RDR 0002's load categories and mints no
  RDR 0006 code for either, which is what the fence requires.
- **Status**: OPEN — citation repair still pending, not implementation-
  blocking.

## D3 — OPEN (unchanged)

- **Check run**: `grep -c '§D13' docs/rdr/0003-*.md` → **0**, and the
  check expects ≥1.
- **Outcome**: the §D13 citation is still absent from the record, and
  adding it is an RDR edit this stage does not make. The implementation
  follows §D13 regardless: a set literal and a set-valued held value both
  cross as the canonical JSON array (`assignment.go::renderSet`,
  `grammar.go::parseSetLiteral`), matching RDR 0002's `seamValue`.
- **Status**: OPEN — citation repair still pending, not implementation-
  blocking.

---

## D4 — The operator/kind matrix had no enforcement point

- **Type**: DEPENDENCY-LIMIT
- **Status**: mechanical translation
- **What**: REQ-6, REQ-112 and REQ-134 require an unsupported
  operator/tag-kind pair to refuse at load as `malformed_predicate_atom`.
  RDR 0002's loader enforced the operator SET and the literal SHAPE but
  not the matrix, so `contains` on an `enum` and `gte` on an `enum` loaded
  clean.
- **Evidence**: `0002:C16` ("This RDR does not mint operators and does not
  widen the set") and `0002:C22` already land RDR 0003's rejection rules
  in that category, and JDR 0001 §D7(iii) routes them there. The loader is
  the only place a rule's atoms are validated, and `guard` imports `table`
  rather than the reverse, so the mirror sits beside the existing
  `declaredKinds` / `operators` mirrors with `guard.Accepts` named as the
  authority.
- **Resolution**: `internal/table/normalize.go::atom` gains the gate;
  `internal/table/model.go` gains the unexported `operatorKinds` mirror.
  Scoped to GUARD blocks: the RDR states the matrix of the guard atom
  shape and separates match keys from guard dimensions deliberately
  (`0003:C13`), and a match block keeps its own `eq`/`in` restriction.

## D5 — `contains` accepted a repeated element and an empty literal

- **Type**: DEPENDENCY-LIMIT
- **Status**: mechanical translation
- **What**: REQ-83 requires "a repeated element MUST be rejected at parse
  rather than silently collapsed" of every set literal — the right-hand
  side of `in` AND of `contains`. The loader rejected a repeated `in`
  member but let a repeated `contains` member through to the canonicalizing
  dedup, which is exactly the silent collapse the clause names. REQ-4
  likewise publishes both set-shaped literals as NON-EMPTY.
- **Evidence**: `0003:C22` names both operators in one clause; the
  loader's own `in` arm already cites it and calls `firstDuplicate`.
- **Resolution**: `internal/table/normalize.go::atom` applies the same
  duplicate check to `contains`, and both arms reject an empty literal.

## D6 — The published bound is an implementation choice REQ-131 constrains

- **Type**: IMPL-DECISION
- **Status**: derived choice — recorded because it affects interpretation
- **What**: the RDR requires the implementation to PUBLISH a bound but
  fixes no value. REQ-131 then builds its equal-cardinality shape pair
  relative to B: `2^n ≥ 2B` for the narrow shape and `side² ≥ 2B` for the
  wide one, repeated at `B/2`. The two agree only when `2B` and `B/2` are
  both EVEN powers of two, so not every integer is a bound this RDR's own
  Scenario 3 can be built around.
- **Resolution**: `guard.Bound()` publishes **2048** — the pair builds at
  4096 (over, both shapes refuse) and 1024 (under, both prove). The bound
  doubles as the proof representation's budget, which is what makes
  `ProofCompletes` and the bound's verdict one quantity rather than two
  that happen to track (REQ-133 / A15).

## D7 — A scoped row group with no guard dimension carries no claim

- **Type**: SPEC-UNDER
- **Status**: derived choice
- **What**: the RDR does not say what a group whose rows carry NO guard
  atom is. Read literally, its scoped product is the empty product — one
  assignment, which any row accepts — so such a group would certify GREEN.
- **Evidence**: REQ-63 scopes the claim to "every scoped row group whose
  participating dimensions are all finitely declared"; a group with no
  participating dimension has nothing for a claim to range over, and
  REQ-82 requires a reader be able to tell a group PROVED over its
  declared domains from one that was not.
- **Resolution**: such a group is reported with no verdict and no green
  claim. It is not blocked either — there is no dimension to name in a
  blocking finding — so it is neither proved nor refused, which is the one
  reading consistent with both clauses. Reachable via
  `optionalButAlwaysWrittenSource`'s writer row (REQ-101).

## D8 — Phase 1 fixture helpers could not reach the clauses under test

- **Type**: TEST-FIXTURE
- **Status**: mechanical translation
- **What**: three Phase 1 fixtures asserted contracts they could not
  reach.
  - `loadCategory` built a rule with no `[rule.write]`, so RDR 0002's
    rule-SHAPE check refused every case as `malformed_rule_shape` before
    any atom was parsed. An empty write block restores the atom defect as
    the reported one and constrains nothing this RDR owns.
  - `resolveWith` supplied the view as OBSERVED tags only, so a row whose
    `RequiresOwned` names a key refused `owned_state_unavailable` before
    its guard was reached, and it hardcoded `Recognized: "go"`, which the
    MVV slices' own alphabets do not contain. It now supplies the one
    evaluation view under both provenances (value-identical: `assemble`
    resolves a cross-provenance key by precedence, never as a conflict)
    and resolves under the table's declared alphabet.
  - `twoPopulationSource`'s ordinary rows closed coverage on their own, so
    REQ-79's "the escape rows were excluded from coverage" held whether
    the escape rows were counted or dropped. A third domain value only the
    ESCAPE rows reach restores the discrimination.
- **Note**: no assertion was relaxed in any of the three; each fixture now
  reaches the clause its test quotes.

## D9 — RDR 0002's FAIL-1 fixture authored a matrix-illegal atom

- **Type**: TEST-FIXTURE
- **Status**: mechanical translation
- **What**: `TestFail1_MultiMemberLiteralIsNeverTruncatedAtTheKernelSeam`
  asserted `contains` carriage over a `scalar`-kind tag, which D4's gate
  now refuses. The test predates the gate; its contract — a set-valued
  guard literal crosses the seam WHOLE, canonical, and distinguishable —
  is about carriage, not about which kind carries it.
- **Resolution**: the fixture gains a declared `set` tag and its three
  `contains` arms move onto it. The `in` arm stays on the `scalar`, so the
  test still shows carriage is decided by the OPERATOR rather than by the
  declared kind — the point the fixture's own comment makes. No assertion
  changed.
