# Deviations — RDR 0002 Transition table as reviewable data

Pre-seeded by the `0002-0009` cluster gate, iteration 4 (2026-08-24 —
`docs/rdr/cluster-reconcile/0002-0009/iter-4/reconcile-report.md`). Stage 8
opens this file; running each named check IS the entry's disposition. An entry
escalates only if its check contradicts a contract.

**Run order: this RDR runs SECOND (2 of 8), after 0007** — the atom-shaped
`Row` and its constants do not exist in `internal/resolve` until 0007's reshape
lands (`0002:1651`, `1662`, `1772`, `2263`, `2278`). It must also precede 0003:
D5's §D13 clause is the encoding 0003 cites, and moving that declaration here
is what broke the 0007↔0003 cycle. See
[`../../BUILD-ORDER.md`](../../BUILD-ORDER.md).

---

## D1 — Canonical fixtures vs RDR 0003's declaration vocabulary

- **Type**: TEST-FIXTURE
- **Status**: OPEN (deferred at cluster gate; DEFER-TO-IMPLEMENTATION)
- **Source**: `iter-4/pairwise-0002-0003.md` F1; `iter-4/critique-set.md` Q-1.
- **Conditions carried**: (a) the offending text is fixture/evidence/
  illustrative — `evidence/spikes/iter-2/rdr-fixture.toml:65` (`kind =
  "string"`), Illustrative Code `0002:1698`, A13 evidence `0002:340`,
  Round-Trip prose `0002:1551`, and the fixture's `single_valued` /
  assignment shape on `iter` / `cluster_ready`; 0002's own fence defers
  the kind vocabulary to RDR 0003 (`0003:791-793`: exactly five tokens,
  `scalar` not `string`) and the assignment table to `0003:948-952`,
  `1273-1290`; (b) checks below; (c) no clause's meaning changes — the
  fences already agree, the artifacts lag them.
- **Checks** (Stage 8):
  1. `grep -rn 'kind = "string"' docs/rdr/0002-*/evidence/spikes/iter-2/`
     → 0 after the fixtures are rewritten to 0003's five tokens; the loader
     refuses an unknown kind token as `malformed tag declaration`.
  2. Run RDR 0006's lint (or, before it exists, 0003's assignment table by
     hand) over the promoted fixtures: `iter` / `cluster_ready` declared
     `single_valued` must satisfy `0003:948-952`. If a fixture cannot be
     made to conform without a fenced change in 0002 or 0003, escalate as
     SPEC-DEFECT.

## D2 — Empty write block on an escape row has no fixture

- **Type**: TEST-FIXTURE
- **Status**: OPEN (deferred at cluster gate; DEFER-TO-IMPLEMENTATION)
- **Source**: `iter-4/pairwise-0009-0002.md` F1 (JD-11 row 3 residual).
- **Conditions carried**: (a) both fences agree (`0002:833-836`,
  `0009:843-846`: `writes = []` is a presence-keyed rejection); the gap
  is in evidence — the promoted spike checks `len(rule.Write) > 0`
  (`iter-2/main.go:546`) and `gen-cases.py:60-62` populates
  `neg-escape-with-write`; (b) check below; (c) no meaning change.
- **Check** (Stage 8): mint `neg-escape-with-empty-write` (`write = []` on
  an escape row) and assert `malformed escape declaration`; the loader must
  key on key presence, not length. RDR 0009 Scenario 8 is the second oracle.

## D3 — Write-block value kind/domain conformance names no load category

- **Type**: TEST-FIXTURE
- **Status**: OPEN (deferred at cluster gate; DEFER-TO-IMPLEMENTATION)
- **Source**: `iter-4/pairwise-0002-0006.md` F1.
- **Conditions carried**: (a) unfenced on both sides — 0002's arity rule
  claims the check "at minimum" but lists no category; 0006 checks it in
  invariants 1 and 5; (b) check below; (c) additive — no clause changes.
- **Check** (Stage 8): one Scenario-3 fixture writing a literal outside the
  declared domain (and one of the wrong kind) → a named 0002 load category
  (`malformed tag declaration` / literal-outside-domain per §D7(iii)).
  If the category set must widen in fenced text, escalate.

## D4 — Stale peer-status sentences inside fences (citation repair)

- **Type**: TEST-FIXTURE
- **Status**: OPEN (citation repair pending 0002's next touch)
- Sites: `0002:1053-1058` ("RDR 0004 does not yet carry that clause … not
  end-to-end until RDR 0004 lands §D5") — 0004 re-locked 2026-08-24 with the
  clause fenced at `0004:357-364`; `0002:1229-1233` ("RDR 0003 is `Draft
  [revised from Final 2026-08-24]`") — 0003 re-locked the same day. Both are
  self-scoped status claims; no obligation changes. Also `0002:2073-2077`
  ("the reshape has no owner") vs 0007 Phase 1 (`0007:2117-2126`), and
  `0002:2297-2298` Perf Expectations "suffix … an alphabet member" vs the
  fenced sequence-suffix rule. Artifact of record: `docs/rdr/README.md`.
- **Check**: `grep -n 'does not yet carry\|has no owner' docs/rdr/0002-*.md`
  → 0, and the `0003 is Draft` sentence at `0002:1229` gone.

## D5 — JD-9/19/20/21/22 answered: land §D13, cite §D8/§D9/§D11/§D12

- **Type**: TEST-FIXTURE
- **Status**: OPEN (one landing + citations, pending 0002's next touch)
- **Source**: JDR 0001 §D8–§D13 (`d937eec`, settled `5c2b96b`, both
  2026-08-24) — answered *after* the `0002-0009` iteration-4 gate recorded
  JD-9/19/20/21/22 as unanswered. Status qualifier and README row corrected
  2026-08-24. No joint decision is now open against this RDR.
- **The landing (§D13, §JD-22)**: this RDR **declares** the set-value
  encoding — `Tag.Value` stays `string`; a set crosses as its canonical JSON
  array, members sorted, duplicate-free, compact. Read-back equality is byte
  equality; it is the same form the CLI accepts in `--write` and emits in
  `writes` (§D11); no third encoding exists. §D13 is explicit: "Lands in
  **0002** — one normative clause". This is a *landing*, not a citation, and
  it is the only one in this sweep.
- **The citations**: §D8 (§JD-9) `resolve`/`next` assemble `Input.Owned`,
  `--tag` stays `Observed`; §D9 (§JD-19) gates run after exact-one selection,
  `deny` is a refusal, `set-state` never gates; §D11 (§JD-20) `--clear <key>`,
  JSON-array set writes, `--write k=<clear>` refused; §D12 (§JD-21) `Block`
  gains `BlockMatch`, added in 0007 Phase 1 — this RDR's "tell an implementer
  to widen" sentence (`0002:923-938`) is satisfied by that, not by a change
  here.
- **Check** (Stage 8): the set-value encoding clause exists in a `normative`
  fence citing §D13; `grep -n '§D1[123]\|§D[89]\b' docs/rdr/0002-*.md` → ≥1 per
  answered JD. If the §D13 landing cannot be written without contradicting an
  existing fence, escalate — that is a contract conflict, not a citation.

## D6 — REQ-124's `crypto/sha256` scan fires against its own source

- **Type**: TEST-FIXTURE
- **Status**: mechanical translation (Phase 2; fixture repaired, assertion
  preserved)
- **Site**: `internal/table/roundtrip_test.go::TestReq124_NoGoldenHashOfRenderedText`.
- **The defect**: the test walks every `.go` file in the package and fails
  any file containing the literal `"crypto/sha256"`. Its sibling check —
  the one for the spike SHA — carries a self-exemption
  (`e.Name() != "roundtrip_test.go"`); the hashing check does not. The
  literal `"crypto/sha256"` is written in `roundtrip_test.go` itself, as
  the needle of that very check, so the file matches its own scan. The
  assertion is therefore RED against every possible implementation,
  including a correct one: no code change can make it pass. It is a
  Phase-1 authoring slip, not a spec obligation the implementation failed.
- **Grounding**: `0002:TS` scenario 2 states the obligation as "**This SHA
  MUST NOT be asserted as a golden hash by any implementation test** …
  Assert over the normalized value" (REQ-124). The obligation binds
  *implementation* tests hashing *output*; the scan's needle in the test's
  own source is neither. Nothing in the record asks the scanner to exempt
  itself from the SHA check but not from the hashing check — the asymmetry
  is unmotivated, and the sibling exemption two lines above shows the
  intended shape.
- **Resolution**: apply the SAME self-exemption the SHA check already
  carries. The assertion's discriminating power is untouched: any other
  file in the package — implementation or test — that imports
  `crypto/sha256` or `crypto/md5` still fails. No implementation file
  imports either.

## D7 — REQ-137's third leg omits an owned key its own scenario needs

- **Type**: TEST-FIXTURE
- **Status**: mechanical translation (Phase 2; fixture repaired, assertion
  preserved)
- **Site**: `internal/table/mvv_test.go::TestReq137_ResolveOverTheNormalizedFixtureRows`,
  subtest "an unevaluable sibling refuses guard_unevaluable".
- **The defect**: the leg's `owned` map carries `status`, `stage`,
  `profile`, and `iter`, and deliberately omits `cluster_ready` so
  `continue-prelock-cluster`'s `[rule.guard.all.cluster_ready] eq = true`
  is unevaluable. But that rule also *writes* `prelock_lens`, so its
  `RequiresOwned` is `[prelock_lens, stage]` (REQ-77: "the sorted,
  duplicate-free set of tag keys named by the rule's write block and clear
  list"), and `prelock_lens` is absent from the same map. The kernel's
  gate reports owned state before an undecidable guard (`0007:C7`,
  source-verified at `internal/resolve/resolve.go::gate`), so the row
  refuses `owned_state_unavailable` and the guard verdict is never
  reached. The leg asserts `guard_unevaluable`.
- **Grounding**: the two clauses are jointly satisfiable and neither is
  wrong. `0002:TS` scenario 4 expects "the unevaluable-sibling tag-set
  refuses `guard_unevaluable` even though a decidable sibling and a
  `no_match` escape row exist" (REQ-137), and `0002:C14` derives
  `RequiresOwned` from the write-plus-clear key set with no guard-read
  keys added (REQ-77, REQ-79) — which the shipped normalizer does, and
  which `TestReq77` asserts as `[prelock_lens, stage]` on this very row.
  The kernel's precedence is `0007:C7`'s and is "pinned to shipped
  behavior". So the contradiction is not between contracts: it is that
  this leg's tag-set does not satisfy the *precondition* its own comment
  states — the guard must be the row's only undecided input.
  Leg 1 does not hit this because `cluster_ready = "false"` decides that
  guard FALSE, pruning the row and its owned obligation with it.
- **Resolution**: add `"prelock_lens": "critique"` to the leg's owned map.
  The scenario is unchanged and every discriminating element survives:
  `cluster_ready` stays absent, the cluster row still survives the prune
  as unevaluable, `continue-prelock` is still decidable, and the
  `no_match` escape row still exists. Verified: the refusal is
  `guard_unevaluable`, and its payload names exactly
  `{cluster_ready, all, eq, true, absent}` on
  `continue-prelock-cluster` — the atom the leg exists to witness.

---

## Phase 2 dispositions for the pre-seeded entries

Running each entry's named check IS its disposition. None contradicted a
contract, so none escalated.

- **D1** — DISCHARGED at the promoted set. Check 1:
  `grep -rn 'kind = "string"' internal/table/testdata/` → **0**; the
  promoted fixtures carry RDR 0003's five tokens, and the loader refuses
  an unknown kind token as `malformed_tag_declaration`
  (`load.go::tagDecl`, asserted by `TestReq64`). Check 2 (`iter` /
  `cluster_ready` conform to `0003:948-952`): `iter` is
  `kind = "int", min = 0, max = 9, required = true` and `cluster_ready` is
  `kind = "bool"` — neither declares `single_valued`, so the assignment
  table is satisfied without a fenced change on either side. The *spike*
  directory retains 71 hits and is deliberately untouched: it is Stage-4/6
  evidence of what was reviewed, and this build is read-only over it. The
  residual is re-running `gen-cases.py` to refresh that evidence; it moves
  no clause and no assertion.
- **D2** — DISCHARGED. `neg-escape-with-empty-write.toml` is promoted and
  refuses `malformed_escape_declaration`. The loader keys on key PRESENCE,
  not length: `sourceRule.Write` is `*map[string]any`, so `write = []`
  and a populated block refuse alike (`normalize.go::normalizeRule`).
- **D3** — DISCHARGED, no widening. Both scenario-3 fixtures
  (`neg-write-value-outside-domain`, `neg-write-value-wrong-kind`) refuse
  `malformed_tag_declaration`, the category D3's check names, via
  `renderWrites`'s `conform` call. No fenced category had to widen, so the
  escalation arm did not fire.
- **D4** — NO ACTION for this build; it is a citation repair on the record,
  and RDRs are never amended. Its substance was already discharged in
  `req-list.md`'s standing correction 2: the 0007 reshape HAS landed, so
  every REQ marked *(was blocked)* was satisfiable and is now green.
- **D5** — DISCHARGED as a LANDING. The §D13 set-value encoding is
  implemented at the kernel seam (`model.go::seamValue`): a set crosses as
  its canonical JSON array, members sorted, duplicate-free, compact,
  matching the form RDR 0007 already wrote its `contains` leg against
  rather than redefining it. Read-back equality is byte equality over that
  array. `TestReq93` asserts all four legs. It contradicted no existing
  fence, so the escalation arm did not fire.

---

## Phase 3c deviations

Opened while fixing the Phase 3a (`FAIL-N`) and Phase 3b (`ADV-N`) findings.
Resolutions and the tests that guard them are in
[`verification.md`](verification.md) § *Phase 3c — Fixup resolutions*.

## D8 — A guard atom's literal carriage is keyed on the OPERATOR, which no fence states

- **Type**: SPEC-UNDER
- **Status**: resolved from evidence (no author decision needed)
- **Site**: `internal/table/model.go::setValuedLiteral`, `::seamValue`,
  `Row.KernelRow`. Filed against FAIL-1.
- **The gap.** REQ-93 (the §D13 landing) says "a **set** crosses as its
  canonical JSON array", covering "any set-valued atom literal handed to
  the guard seam" — but it never says what makes an *atom literal* set-
  valued. The shipped build read it as the tag's declared `kind`, which
  truncated a multi-member `in`/`contains` literal on any non-`set` tag to
  `members[0]`. Nothing in RDR 0002's fences distinguishes the two
  readings.
- **Evidence the resolution rests on.**
  1. `internal/resolve/guardcontract.go` — the cross-RDR contract RDR 0007
     shipped, which RDR 0003's evaluator is instantiated against — fixes
     the byte form per OPERATOR and mentions no declared kind:
     `{"in", ["alpha","beta"], "alpha"}` and
     `{"contains", ["alpha"], ["alpha","beta"]}`. BUILD-ORDER is explicit
     that 0007 wrote these against §D13 first and "0002's later run must
     **match**, not redefine".
  2. RDR 0003 `0003:948-960` fences which operators take a set:
     "**`eq`, `in`, and the integer comparisons are single-value
     operators** … an atom denotes the assignments in which the tag's
     single held value **is the literal** (`eq`), is a member of the
     literal set (`in`) … `contains` is the operator for a tag whose …".
     So `in` and `contains` are exactly the two whose right-hand side is a
     member set, whatever the tag's kind.
  3. `0002:C17` (REQ-89) already makes conformance "per operator, not per
     literal", so keying carriage on the operator is the same axis the
     record uses for the neighbouring obligation.
- **Resolution.** The tag VALUE side (`Tag.Value` on match tags,
  `NextTags`, `Writes`) stays keyed on the declared kind — that is the
  question §D13 answers, and a one-member set must stay `["a"]` or RDR
  0004's read-back could not tell it from the scalar `a`. The guard atom
  LITERAL side is additionally keyed on the operator. The two rules answer
  different questions and neither displaces the other; no fence had to
  widen, so the escalation arm did not fire.
- **Additive surface**: none. `setValuedLiteral` is unexported and no
  exported identifier, signature, or category was added.
- **Residue for the record.** RDR 0002 is §D13's landing document, so the
  operator half of the carriage rule arguably belongs in its normative
  clause rather than only in this package's comments. RDRs are never
  amended, so it is noted here for whoever next opens the record.

## D9 — Arity is not stated as part of "well-formed for the declared kind"

- **Type**: SPEC-UNDER
- **Status**: resolved from evidence (no author decision needed)
- **Sites**: `normalize.go::atom` (the `eq` arm), `::renderWrites`,
  `load.go::loadInitial`. Filed against FAIL-1's refusal half.
- **The gap.** REQ-30 requires every `[initial]` value to be "well-formed
  for that tag's declared kind and domain", and `0002:C4` says "for a
  `set` kind the array literal is the whole new set" — but neither says
  what a multi-member array MEANS on a kind that is not `set`, and
  `0002:C17` lists `eq` among the operators that "take members" (plural)
  without fencing its arity. The shipped `conform` checked each member
  individually, so a multi-member literal passed every stated check and
  was then silently truncated at the seam.
- **Evidence the resolution rests on.** RDR 0003's single-value-operator
  fence (`0003:948-960`, quoted in D8) gives `eq` and the comparisons a
  denotation over "a tag holding exactly one value", so a multi-member
  `eq` literal denotes nothing; `0002:C22` (REQ-64) gives this package
  "the two load categories that carry RDR 0003's rejection rules", of
  which `malformed predicate atom` is one. `0002:C3` (REQ-6) supplies the
  governing principle for the write and `[initial]` sites: a malformed
  authoring is "a stable refusal rather than a silent no-op".
- **Resolution.** A multi-member literal on `eq` refuses
  `malformed_predicate_atom`; a multi-member value on a non-`set` kind
  refuses `malformed_tag_declaration` at the write site and
  `malformed_initial_declaration` at `[initial]` — each the category that
  site already carries for a kind/domain defect, so no category was
  minted. A one-element array stays admitted, since `0002:C13` fixes
  `eq = "x"` and the one-member spelling as one spelling of one edge.
- **Additive surface**: none. No category was added to `Categories()`; the
  closed set of 25 is unchanged.
- **Why a refusal and not only a seam encoding.** Encoding a multi-member
  `eq` as an array would invent a literal form no consumer parses: the
  shipped contract drives `eq` with a bare string, so an array would reach
  RDR 0003's evaluator as an unparseable literal and answer
  `GuardUnevaluable` — trading a silent truncation for a silent deadlock.

## D10 — REQ-57's member sort is exercised on `contains`, applied here to `in`

- **Type**: IMPL-DECISION
- **Status**: mechanical translation
- **Site**: `internal/table/normalize.go::atom`. Filed against ADV-3.
- **The reading.** REQ-57 (`0002:C10`) states "Members sort
  byte-lexicographically so two authored orderings of one set are one
  literal" with no operator qualifier, and REQ-56 makes the obligation
  bind "a set-valued literal". An `in` literal is a member set — RDR 0003
  fences it as "a member of the literal **set**" — so the sort binds it.
  The shipped build excluded `in` from the sort; `TestReq57` happens to
  exercise `contains`, so nothing observed the exclusion.
- **Why it is safe.** The row set is unchanged: a match-block `in` still
  yields one row per member, and rows re-sort by the identity tuple
  afterwards (REQ-98), so only the enumeration order inside `expand`
  moves. The `slices.Compact` that accompanies the sort can drop nothing,
  because a repeated `in` element now refuses at load (ADV-3) before the
  sort runs — the two changes are ordered deliberately.
- **Why it was needed.** Phase 3b's `TestAdv1_GuardBlockInMustNotExpand`
  asserts the surviving guard `unless` literal is `[mid small]` from an
  authored `["small", "mid"]`. Without the sort the assertion could only
  be satisfied by weakening it, which `§scope-discipline` forbids.
- **Additive surface**: none.
