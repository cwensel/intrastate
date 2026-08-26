# Deviations — RDR 0008 Ownership of the recognized-outcome tag key name

Pre-seeded by the `0002-0009` cluster gate, iteration 3 (2026-08-24 —
`docs/rdr/cluster-reconcile/0002-0009/iter-3/reconcile-report.md`). Stage 8
opens this file; running each named check IS the entry's disposition. An entry
escalates only if its check contradicts a contract.

**Run order: 7 of 8** — §JD-5's ordering resolves in run 8, when the second of
this RDR's and 0009's `Resolve`-entry checks is written. See
[`../../BUILD-ORDER.md`](../../BUILD-ORDER.md).

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

## D3 — JD-8/JD-9 answered: cite §D10 and §D8

- **Type**: TEST-FIXTURE
- **Status**: OPEN (citation repair pending 0008's next touch)
- **Source**: JDR 0001 §D10/§D8 (`d937eec`, settled `5c2b96b`, both
  2026-08-24) — answered *after* the `0002-0009` iteration-4 gate. Status
  qualifier and README row corrected 2026-08-24. **§JD-5 remains open**
  (precondition precedence vs 0009) and stays in the qualifier.
- **The answers**: §D10 (§JD-8) — one `Code` per caller-branchable failure,
  kernel-aligned names, one `omitempty` `findings` field, exit 3 = environment
  could not be consulted; `reserved_tag_key` gets its constant here. §D8
  (§JD-9) — `--tag` enters as `Observed` and never satisfies an owned-state
  dependency; `--tag` on an owned key **or on `recognized`** is a
  `GroupUserEnv` refusal at the CLI. That is this RDR's reserved-key check
  given a site and an error group.
- **Check** (Stage 8): `grep -n '§D10\|§D8\b' docs/rdr/0008-*.md` → ≥1 each;
  the reserved-key refusal names `GroupUserEnv` and the CLI site. Escalate only
  if the CLI-site placement contradicts this RDR's `Resolve`-entry check —
  note that ordering question is §JD-5's, still open, and must not be
  pre-decided by this citation.

---

## D5 — Near-miss advisory has no settled delivery channel (Phase 0 Q1)

- **Type**: SPEC-UNDER
- **Status**: mechanical translation (unattended run: recorded, proceeding
  under the Phase 0 auditor's reading; not escalated to the author)
- **Source**: Phase 0 spec audit QUESTIONS Q1, `req-list.md`.
- **Gap**: §D10 places the near-miss advisory in `Finding.Hint`, but
  `Findings` rides `CLIError` (the failure envelope) while a near-miss by
  definition loads *clean*. No contract names the clean-load carrier.
- **Resolution proceeded under**: table-side advisory list returned
  separately from the error; CLI delivery deferred to RDR 0005/0006, which
  own the CLI output contract. This keeps 0008's surface at the table/kernel
  boundary the RDR actually claims and adds no new CLI public surface.

## D6 — Block 3's payload requirement over `unknown tag` is vacuous (Phase 0 Q2)

- **Type**: TEST-FIXTURE
- **Status**: mechanical translation (unattended run: recorded, proceeding
  under the Phase 0 auditor's reading)
- **Source**: Phase 0 spec audit QUESTIONS Q2, `req-list.md`.
- **Gap**: block 3 requires the three-field payload on `unknown tag`
  failures over the reserved key, but implemented RDR 0002 trips
  `CatMalformedModelDeclaration` ("no [tags.recognized] declaration")
  before any `unknown tag` can fire on the reserved key.
- **Resolution proceeded under**: the payload binds `reserved_tag_key`
  failures only. The `unknown tag` arm is unreachable given 0002 as
  implemented, so asserting it would test a path the predecessor forecloses.

---

## Stage 8 Phase 1 (test authoring) dispositions — 2026-08-26

### D1 — DISCHARGED, no SPEC-DEFECT

D1's check was run. Scenario 3 is written against RDR 0007's per-atom seam in
`internal/resolve/reserved_key_sameview_0008_test.go`
(`TestReq76And77And78_GuardSeamAndMatcherReadTheSameRecognizedBinding`), using
`valueSeam` — the per-atom stand-in with a call recorder already present at
`internal/resolve/guard_fixtures_test.go:64`. That is a NEW seam alongside
`fixtureGuards`, not a change to it, which is REQ-76's own requirement.

The test captures the `value` handed to the guard atom over `recognized` and
asserts it equals `in.Recognized`, while a `Match` on the same key selects the
row in the same resolve. Provenance is read as `(in.Recognized,
ProvenanceRecognized)` on the package-internal path over the identical `Input`,
because `0007:C1` fences `Evaluate` from the view.

**The same-view property IS assertable at that seam** — the test passes against
the shipped kernel. D1's escalation branch ("if the same-view property cannot
be asserted at that seam, escalate as SPEC-DEFECT") therefore does not fire.

### D5 — applied as written

The near-miss advisory is asserted at the `internal/table` boundary as a list
returned separately from the error, via a net-new
`LoadWithAdvisories(src, sourceID) (*Model, []Advisory, error)`. `Load`'s
existing signature is untouched, so no caller changes. No CLI surface is added;
REQ-102 pins only that §D10's `Finding.Hint` carrier already exists.

### D6 — applied as written

REQ-26's `unknown tag` arm is not asserted. Confirmed empirically against HEAD:
`neg/neg-no-recognized-decl.toml` refuses `malformed_model_declaration` at
`internal/table/load.go:169-171`, before any `unknown tag` failure can carry the
reserved key. The three-field payload binds `reserved_tag_key` failures only,
pinned by `TestReq15And16_TheMissingDeclarationLowerBoundStaysOutsideReservedTagKey`.

### D3 — citation half carried forward

The record-side half of D3 (`grep -n '§D10\|§D8\b' docs/rdr/0008-*.md` must hit
a body clause, not only the Status line) is a documentation edit, not a test.
It is unaddressed by this Phase and is flagged as the sole REQ-side orphan in
`coverage.md`.

### New surface the tests pin (implementation contract)

Phase 1 chose the names the tests bind, within the latitude REQ-34/64/70 grant:

- `resolve.CheckInput(in Input) error` — the one exported predicate
  (`0008:C4`). Not named `Final` (REQ-40/70). `Resolve` must call this exact
  function at entry, not a second independently-written check (REQ-68, pinned
  by a two-site agreement test that fails when the predicate exists but
  `Resolve` does not call it).
- `table.Failure` gains `Offending`, `Remedy`, `Rule string` — block 3's three
  distinct machine-readable fields.
- `table.Advisory{Authored, Reserved, Rule string}` plus
  `table.LoadWithAdvisories` — the near-miss advisory's separate channel.
- Rule identifier literals, asserted byte-for-byte:
  `reserved-tag-key/kernel-owned`, `reserved-tag-key/author-must-rename`,
  `reserved-tag-key/near-miss`.
- A doc comment on `internal/resolve/resolve.go::recognizedTagKey` citing
  RDR 0008 (Phase 1's IP deliverable, REQ-9).

`internal/resolve/export_0008_test.go` is a `package resolve` test file adding
`ResolveBypassingPreconditionForTest`, `AssembledBindingForTest`, and
`AssembledLenForTest`. It exists only in the test binary and adds no production
surface; it is what keeps blocks 5 and 6's "documented residual" executable once
the entry precondition makes those paths unreachable through `Resolve`.

---

## Stage 8 Phase 2 (implementation) dispositions — 2026-08-26

All 23 net-new tests are green with no new deviation of contract type. Phase 1's
§"New surface the tests pin" was followed exactly: `resolve.CheckInput` (the one
exported `Input` predicate, called by `Resolve` at entry), `Failure.Offending`
/ `.Remedy` / `.Rule`, and `table.Advisory` / `table.LoadWithAdvisories`. No
name was renamed, no field reshaped, and no new public surface beyond those
three items was added.

D1, D5 and D6 were honoured as Phase 1 settled them; none reopened.

### D7 — the advisory scan reads the source bytes, not the loaded model

- **Type**: IMPL-DECISION
- **Status**: mechanical translation
- **Site**: `internal/table/advisory.go::nearMissAdvisories`
- **The latitude**: REQ-64 leaves the carrier's Go type and field names to the
  implementer, and REQ-63 pins only that the advisory channel is *separate*
  from the validation-failure list and never alters the verdict.
- **The choice**: `nearMissAdvisories` re-decodes the source bytes with a
  permissive `map[string]any` probe over `[tags]` alone rather than reading
  `Model.Tags`. `Load` is fail-fast (`0002:C3`), so a document that refuses
  returns a nil model and would carry no key set to advise on; reading the
  source keeps the two channels genuinely independent, which is what
  `TestReq64_TheAdvisoryChannelIsSeparateFromTheFailureChannel` asserts by
  loading a *failing* fixture through `LoadWithAdvisories`.
- **Recorded because** it affects future interpretation: an implementer who
  later moves the scan onto the normalized model would couple the advisory to
  the load verdict and silently break REQ-63's separateness on refusing
  documents, even though the fixtures in the current suite would stay green.
- Advisory order is sorted by authored spelling. Nothing in the RDR fixes an
  order; the sort exists only so the channel is reproducible for a golden
  consumer, and no test asserts it.

### D8 — the near-miss fold is `strings.EqualFold` over the trimmed key

- **Type**: IMPL-DECISION
- **Status**: mechanical translation
- **Site**: `internal/table/advisory.go::isNearMiss`
- Implements REQ-59/REQ-60's disjunction as
  `key != "recognized" && strings.EqualFold(strings.TrimSpace(key), "recognized")`.
  A single expression covers folding-alone, trimming-alone, and both-together,
  which is exactly the disjunctive trigger REQ-60 requires and what
  ASSUMPTION-7 already anticipated (Go's `EqualFold` is simple folding).
  Recorded only so the disjunction's single-expression form is legible as
  deliberate rather than as an accidental conjunction.

### Still open, unchanged by this Phase

- **D2** — its check is `TestReq106_NoAccessorKeysEntryMayNameTheReservedKey`,
  green against the predecessor at `internal/table/load.go`. No escalation.
- **D3** — the record-citation half (`grep -n '§D10\|§D8\b'` must hit a body
  clause of `docs/rdr/0008-*.md`, not only the Status line) remains owed. It is
  a documentation edit on the record, and this Phase writes implementation
  only. Carried forward as the sole open item; also flagged in `coverage.md`.

**No `needs author decision` entry was recorded in this Phase.**

---

## Stage 8 Phase 3c (fixup) dispositions — 2026-08-26

Phase 3a returned PASS with no `FAIL-N`. Phase 3b returned BLOCK with three
added tests (ADV-1/2/3). The two verifiers converged on the same three seams
and disagreed on whether they violate *this* RDR. This Phase adjudicates.

### ADV-3 / 3a's S1 — FIXED (no deviation)

Both verifiers agreed independently, so no adjudication was needed. The scan at
`internal/table/load.go` ranged a Go map, so a model breaching **both** naming
directions returned one of two **contradictory** remedies at random (measured
161/39 over 200 loads). Fixed by scanning `slices.Sorted(maps.Keys(decls))`.

REQ-82 (TS-5) licenses an unspecified choice between two **same-direction**
declarations; it does not license a coin-flip between opposite edits, and
REQ-27/33 make the rule identifier a token a golden test asserts byte-for-byte
and a consumer uses for remediation lookup. Sorting satisfies both readings and
narrows no latitude. `TestAdv0008_DoublyBreachingModelReportsOneStableDirection`
is green across repeated runs.

### D9 — the reserved-key payload and near-miss advisory stop at the table boundary

- **Type**: SPEC-UNDER
- **Status**: mechanical translation (unattended run: adjudicated against the
  record; not escalated to the author). Supersedes nothing; extends D5.
- **Site**: `internal/cli/lint.go` (`model-invalid`, zero findings) and
  `internal/cli/flow_input.go::loadFailure` (`Param`/`Hint` unpopulated).
- **The seams** (3b's ADV-2 and ADV-1, 3a's S3 and S2). Both confirmed real by
  both verifiers: the three-field payload and the near-miss advisory are
  complete and correct on `table.Failure` / `table.LoadWithAdvisories`, and
  neither reaches a user-observable CLI surface.

**Adjudicated as belonging to RDR 0005/0006, not this RDR.** The evidence:

1. **The RDR disclaims the surface in its own Approach** (item 3, unfenced but
   scope-setting): this RDR "guarantees the failure *exists and carries its
   guidance* at load/lint; which user-facing command surfaces it, and under
   which exit code, is RDR 0005's mapping decision and is **not settled
   here**." The same sentence appears for the advisory in the AP item 3 note.
2. **The normative fence binds the data level, not the wire.** `0008:C3`:
   the three fields MUST be carried "**at the data level**", and "the guidance
   travels in the failure data, not the renderer." `table.Failure` carries all
   three; verified by Phase 3a's T3–T6/G8 and re-asserted by this Phase.
3. **REQ-93 is explicit**: "The consumer is a test stub, **not RDR 0005's
   exit-code map**: this RDR asserts the payload contract, and 0005 owns
   whatever mapping it later adds."
4. **D5 already settled the advisory half** unattended at Phase 0: "a
   table-side advisory list returned separately from the error; CLI delivery
   left to 0005/0006." Reversing that here would relitigate a settled entry on
   no new evidence — 3b's finding is the *predicted consequence* of D5, not a
   discovery that contradicts it.
5. **The carriers are other RDRs' surfaces, already Implemented.** BUILD-ORDER
   puts 0006 at run 4 and 0005 at run 6, both ahead of 0008 at run 7, and both
   are `Status: Implemented`. `internal/cli/flow_input.go` is headed "RDR 0005";
   `internal/cli/lint.go`'s bare `model-invalid` block traces to RDR 0006's own
   skeleton commit `d0194d9` and was never touched by 0008.
6. **The gap is already forbidden by its owner's contract.** RDR 0006's REQ-91:
   blocking findings "MUST remain machine-readable in JSON mode through an
   append-only typed `findings` field owned by `clierr`, not through a
   verb-local wrapper or a **text-only `Detail` string**." `lint.go` does
   exactly the latter for *every* load category — 0002's twenty-odd categories
   included, not just `reserved_tag_key`. It is one pre-existing defect on
   0006's surface with a whole-category blast radius, not a 0008 omission.
7. **`Finding.Rule` is 0006's field, not 0008's.** `clierr.go` documents it as
   "the source rule/context id when the normalized model can provide one";
   3b's test repurposed it for 0008's rule identifier. Binding 0008's token to
   a field another RDR defines differently would itself need a joint decision.

**Why not fix it anyway.** Routing the payload and advisory through
`clierr.Finding` means adding CLI public surface — a new advisory carrier on
the *success* path (a near-miss loads clean, so `CLIError.Findings` cannot hold
it; §D10's own reconciliation note puts the success payload on 0006's
`data.findings`) — that **0008's Normative Contracts do not name**. That is the
"ADDITIVE IS NOT EXEMPT" `SPEC-UNDER` violation precisely: a plausible,
useful-looking extension outside the record's fenced surface, taken
unilaterally by the run that noticed it. Q1 recorded this trap at Phase 0 and
chose reading (b) for exactly this reason.

**Weighed against**: REQ-58's "normative, not deferred". The reading applied is
that "not deferred" scopes the *decision* — the advisory's trigger, values, and
separateness are settled now rather than left to a later RDR — not the delivery
surface the same document twice disclaims. REQ-63's ban on the advisory altering
the load/lint verdict independently rules out the failure envelope as its
carrier, which leaves no carrier this RDR owns.

**Disposition.** No CLI surface added. 3b's two tests are **retained, not
deleted and not weakened**: each now asserts unconditionally the half RDR 0008
owns (the three payload fields on `table.Failure`; the advisory's two spellings
and rule identifier, plus the clean verdict) and then `t.Skip`s the wire half
with the owning RDR named in the skip reason. Both retained assertions were
mutation-checked — flipping `Failure.Remedy` fails the test — so neither is
vacuous.

**Handoff.** The delivery gap is real and reaches no author until 0005/0006
close it. It is recorded here as the named successor obligation rather than
silently absorbed. Worth a follow-up kata against RDR 0006's `lint.go` carrier,
whose REQ-91 the current `Detail`-string collapse already violates for every
load category.
