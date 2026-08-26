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
