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

---

# Phase 3c fixup — new entries

## D10 — `acceptedIn` is expressed over the group's DECIDABLE sub-product

- **Type**: SPEC-UNDER
- **Status**: mechanical translation
- **What**: the RDR states the scoped product as one product over every
  participating dimension, and states that an atom lint cannot project
  takes the blocking inability-to-prove outcome "for that dimension". It
  does not say what a row constraining only PROVABLE dimensions denotes
  when some OTHER dimension of its group is unprovable. Read as one
  indivisible product, every row in such a group becomes undecidable and
  the surviving overlap check (REQ-44) has nothing left to run on.
- **Evidence**: REQ-44 scopes the surviving check to "the group's
  **decidable** rows", which presupposes a group can hold both kinds at
  once. The MVV Scenario-2 walkthrough (`0003:1408`) declines overlap
  because "step 4 produced none **for the two ordinary rows**" — the rows
  that CARRY the unprojectable atoms — rather than because the group holds
  one. REQ-58 assigns the outcome "for that dimension". REQ-95 requires
  every decidable defect in one pass.
- **Resolution**: unprojectability is scoped to the dimension that carries
  it. `acceptedIn` projects onto the group's decidable sub-product — the
  participating keys with a finite declared domain over which no row
  carries an unprojectable atom. A row carrying an atom over an
  undecidable dimension still yields the unprojectable set: dropping the
  atom would widen the row past what it denotes. Where every dimension is
  decidable the sub-product IS the scoped product, so the provable path is
  byte-identical to before, and COVERAGE is still compared against the
  full product — never against the restriction, which would be a green
  claim over a proper subset of the product.
- **New public surface**: none. `decidableKeys` / `decidableProduct` /
  `productOver` are unexported.

## D11 — The published bound applies per dimension, not only per product

- **Type**: SPEC-UNDER
- **Status**: mechanical translation
- **What**: REQ-86 – REQ-92 state the bound over the SCOPED PRODUCT's
  cardinality. They do not say what a single dimension whose own
  assignment count exceeds the bound denotes. Read strictly per-product,
  `Denotation` was free to enumerate a 2^22 value dimension before any
  bound was consulted — on the refusal path itself (ADV-2).
- **Evidence**: the bound's own rationale is "the largest product this
  implementation's enumerating proof representation completes over within
  its budget", so a dimension carrying more assignments than that can
  appear in NO product this implementation enumerates. REQ-108 forbids
  silently capping enumeration; the record's Risks mitigation forbids
  naive powerset enumeration outright and requires the refusal be
  reachable without it. A15 claims the cardinality bound PREDICTS
  provability, which a per-product-only reading refutes.
- **Resolution**: `valueAssignments` declines a dimension whose own
  assignment count exceeds `Bound()`. `lintGroup` tests the bound before
  the projection scan and after every dimension is known finite — the
  order REQ-93 states — so a finite over-large product still refuses
  `graph-product-too-large` carrying `(computed size, bound)`, and a
  product carrying an unprovable dimension still reports no computed size
  (REQ-94). No verdict changes; only the cost of reaching it.
- **New public surface**: none.

## D12 — The exponential must saturate, not wrap

- **Type**: SPEC-DEFECT (implementation), mechanically resolvable
- **Status**: mechanical translation
- **What**: `domainSize` computed a `set` dimension as `1 << |elements|`
  and `spread` computed an unmarked dimension as `1 << |domain|`, neither
  guarded. In Go a shift at or past the integer width yields ZERO and a
  shift of 63 yields a negative. A 64-element element universe therefore
  reported a value dimension of size 0 and a whole-product cardinality of
  0 — comfortably UNDER the published bound, fully "provable", and GREEN.
- **Evidence**: `Cardinality` already saturates at `cardinalityCeiling`
  on exactly this reasoning, stated in its own comment: "a wrapped
  negative would read as under the bound and certify the very product the
  clause refuses". The rule was stated for the product and not applied to
  the dimensions composing it. REQ-86 requires the bound be ENFORCED,
  which a wrapped comparison does not do.
- **Resolution**: `spread` saturates at the same ceiling and floors a
  negative domain at zero. Found while fixing ADV-2; not in the Phase 3b
  findings.
- **New public surface**: none.

## D13 — ADV-4: an inverted int bound carries no readable domain

- **Type**: SPEC-UNDER
- **Status**: mechanical translation
- **What**: Phase 3b recorded ADV-4 as LATENT, because `table.Load`
  refuses `min > max` before a declaration reaches this package. But
  `AssignmentCount` is exported and takes a `table.TagDecl` directly, and
  `agrees`' doc comment states the agreement rule is read "on this RDR's
  own surface, so a declaration built in memory is judged by it too" —
  which it was not: `agrees` never compared the endpoints.
- **Evidence**: REQ-18 fixes `{min..max}` as inclusive at both endpoints,
  "so `{0..3}` has cardinality 4", which presumes an ordered pair; `{3..0}`
  names no value. REQ-15's family makes a declaration that carries no
  readable finite domain take the blocking inability-to-prove outcome
  rather than yielding a count.
- **Resolution**: `agrees` requires `*Min <= *Max`, so an inverted bound
  carries no readable domain — the same answer the loader gives, now
  reachable on this RDR's own surface. `spread`'s floor (D12)
  independently removes the `1 << -2` panic. **Weighed and fixed rather
  than left latent**: the defence-in-depth cost is one comparison, the
  exposure is an exported function, and the doc comment already claimed
  the behaviour. The tripwire `TestAdv_RecordedIntBoundInversion` is
  unchanged and still green — it asserts the loader premise, not this.
- **New public surface**: none.

## D14 — ADV-1: the non-empty rule is the LITERAL's, not the held value's

- **Type**: SPEC-DEFECT (implementation), mechanically resolvable
- **Status**: mechanical translation
- **What**: ADV-1's two surfaces disagreed about whether a held `[]` is a
  value at all. REQ-89 makes a `set` dimension `2^|element universe|` — "a
  set-valued tag holds any SUBSET" — so the scoped product enumerates the
  empty subset, and `Conforms` admits `caps=[]` as a conforming view. But
  `Evaluator.Evaluate`'s `contains` arm parsed the HELD value with
  `parseSetLiteral`, whose non-empty requirement is the operator/kind
  matrix's published shape for the AUTHORED LITERAL (REQ-4). A held `[]`
  was therefore unevaluable, and the kernel refused `guard_unevaluable` on
  a conforming view lint had proved — the false exhaustiveness claim
  REQ-67 forbids.
- **Evidence**: REQ-57 states `contains` denotes "the assignments whose
  held set **contains** every listed element". Containment is TOTAL over
  sets: `[]` does not contain `x`, so the verdict is FALSE, not
  undecidable. Nothing in the record extends the literal's non-empty shape
  to a held value, and REQ-13's declaration model carries no
  non-emptiness marker a `set` tag could be declared with.
- **Resolution**: `parseHeldSet` decodes a held §D13 canonical JSON array
  with no cardinality rule; `parseSetLiteral` keeps the non-empty rule and
  is used for the authored literal only. `Denotation` additionally
  propagates the seam's three-valued verdict instead of collapsing
  UNEVALUABLE into false, so any future unevaluable assignment makes the
  atom unprojectable rather than crediting its complement.
- **Alternative rejected**: dropping the empty subset from the product —
  named as a candidate in the Phase 3b note — would make a `set`
  dimension `2^n - 1` and contradict REQ-89, whose own test
  (`TestReq89_…`) pins the powerset. It would also assert non-emptiness of
  a `set` tag, a declaration property the model does not carry.
- **New public surface**: none.

## D15 — ADV-1's second assertion is stale against its own premise

- **Type**: TEST-FIXTURE
- **Status**: needs author decision — **and the work continued**
- **What**: `TestAdv1_GreenClaimCoversAViewTheRuntimeRefuses` closes with a
  second assertion guarded as "so a fix that merely stops the kernel
  refusing does not satisfy this test: either the empty subset leaves the
  product, or the group must not be green". After D14 the group IS green,
  the empty subset IS in the product, and the two surfaces AGREE — the
  runtime resolves `caps=[]` to `adv1-not` and `caps=["x"]` to `adv1-has`,
  a genuine complete partition. The assertion still fires, on a message
  whose own premise ("the runtime finds its guard UNEVALUABLE over a held
  `[]`") is no longer true.
- **Why this is not a weakened test**: the obligation the assertion states
  is "Lint MUST NOT credit a row with an assignment the runtime cannot
  decide." That obligation is now SATISFIED — the runtime can decide it.
  The author foresaw two fixes and guarded against a third that would have
  weakened the runtime veto (REQ-67). D14 is not that third fix: it
  removes a SPURIOUS refusal the RDR's own `contains` semantics never
  called for (REQ-57), rather than weakening a genuine veto. Both of the
  author's own candidates are ruled out by the enforcement surface —
  dropping the empty subset contradicts REQ-89's pinned powerset, and
  withholding the group's claim would require declaring the dimension
  unprovable, which REQ-58 reserves for atoms that cannot project.
- **What was NOT done**: the assertion was not edited, relaxed, or
  deleted. It stands exactly as authored, and the guard suite therefore
  reports one failure.
- **Standing in its place**: `TestFixup_EmptyHeldSetIsDecidedRatherThanRefused`
  pins the underlying obligation directly — every conforming view a green
  group covers, the runtime decides; the empty subset stays in the
  product — and was verified to FAIL against the pre-fix evaluator.
- **The decision the author owns**: whether to retire the stale assertion
  now that its premise is gone, or to reject D14's reading of REQ-57 and
  require a different resolution. The implementation follows the evidence
  as recorded above.
