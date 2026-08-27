# Deviations — RDR 0010 Phase 1 (tests first)

Each entry records a place where the record's literal text could not be
executed as written, what was written instead, and why the substitute pins
the same obligation. Nothing here weakens a REQ; each is a fixture-shape or
surface-availability finding a later phase must carry forward.

---

## D1 — the class-omitted control cannot be a byte-for-byte `class` strip

**REQs**: REQ-2, REQ-83 (MVV step 6), REQ-94 (SC-4), REQ-103 (IP Phase 2).

**The record says**: "the same model with `class` omitted → lint
`graph-dangling-edge` naming the absent `[initial]` (0006:C18 unchanged)"
and "the class-omitted control" among SC-4's fixtures.

**What is executable**: stripping `class = "decision-table"` from the MVV
table leaves a model whose ordinary rules carry no write block. Under the
state-machine class `0002:C4`'s write-block arm is LIVE — that is precisely
what REQ-16 conditions — so the stripped document refuses at LOAD with
`malformed rule shape` and never reaches lint. It cannot exhibit the finding
step 6 names.

**Written instead**: the class-omitted control is authored as a model of the
same SHAPE (same outcome, same observed dimension, guard-atom
discrimination, no `[initial]`) carrying the owned-tag scaffolding a state
machine needs to satisfy `0002:C4` — one owned tag, its reader and writer,
and a write block on the ordinary rule. Verified against `main`: it loads
and lints to exactly one `graph-dangling-edge` at `element = model`.

**Why it pins the same obligation**: the clause under test is that a
state-machine model with an EMPTY `[initial]` still traverses nothing and
takes `0006:C18`'s missing-root finding (REQ-53) — that the class-keyed seed
AUGMENTS the `len(Initial)` test rather than replacing it. The absent
`[initial]` is what the control varies; the write-block scaffolding is
inert to that clause. `TestReq69_ClassReadersNeverRederiveFromTheOwnedSet`
takes the complementary case (a zero-owned state machine with no root,
authored as an escape row) so the `len(owned) == 0` inference is closed in
both directions.

The same reasoning gives `smZeroOwned` (SC-1's counterpart control, REQ-5 /
REQ-89) its shape: a zero-owned model's ORDINARY rule cannot carry a write
block, so its one rule is an ESCAPE row — the only rule shape a zero-owned
state machine can author while `0002:C4` binds it.

**For Phase 2**: no implementation consequence. If the implementer prefers a
different control shape, the invariant to preserve is "state-machine class +
zero `[initial]` → `graph-dangling-edge` at `element = model`".

---

## D2 — there is no root `dump` CLI verb

**REQs**: REQ-82 (MVV step 5).

**The record says**: "`intrastate dump --model m.toml` renders the `emit`
column for every row with no `[dump]` declared".

**What is executable**: `internal/cli/root.go::NewRootCmd` registers exactly
three commands — `version`, `lint`, and `flow`. No `dump` verb exists, and
this RDR's Implementation Plan does not add one (Phase 4 names a fixture and
two docs). RDR 0002 fences the dump as the `internal/table` surface
(`0002:C19`, `internal/table/dump.go::Dump`), which is what the shipped
`internal/table/dump_test.go` asserts against.

**Written instead**: `mvvStep5DumpAndNext` takes the dump half through
`table.Dump(m)` over the MVV model loaded with no `[dump]` declared,
asserting the `emit` cell renders on every row and that `DumpOrder` carries
`emit`. The `flow next` half is taken through the CLI unchanged.

**Why it pins the same obligation**: the clause is about the DEFAULT column
set including `emit` and the cell rendering for every row — a property of
`dumpColumns` and `column()`, not of any command wrapper. Routing through a
CLI verb that does not exist would test nothing.

**For Phase 2/3**: if a `dump` verb is later added, this assertion should
move to it. Nothing in REQ-82 requires one.

---

## D3 — `internal/table/helpers_test.go::clonedRowSliceFields` must gain `Emit`

**REQs**: REQ-22 (`Row.Emit []EmitValue`), REQ-87/REQ-88 (licensed diffs).

**The situation**: `helpers_test.go` declares `clonedRowSliceFields` and
`assertCloneRowsCoversEverySliceField`, a completeness guard whose own
doc-comment says "a slice field added to `Row` fails loudly rather than
silently reopening the aliasing hole". Adding `Row.Emit` trips it, reddening
the shipped `TestReq44_RenderedKindDoesNotFeedBackIntoTheNormalizedValue`
(RDR 0002) for a reason unrelated to that test's clause.

**Assessment**: this is the guard performing its stated function, and the
sanctioned response is to EXTEND the field list and the detaching helper —
not to relax either. It is nonetheless a diff to a checked-in expectation
that is NOT one of REQ-87's four licensed shapes, so it is recorded here
rather than made silently.

**Status at Phase 1**: NOT made. The list is left as shipped so the red gate
is not muddied. `TestReq44_RenderedKindDoesNotFeedBackIntoTheNormalizedValue`
will be red against a Phase-2 tree that adds `Row.Emit` until the list gains
`"Emit"` and `cloneRows` detaches it.

**For Phase 2**: append `"Emit"` to `clonedRowSliceFields` and add
`r.Emit = slices.Clone(r.Emit)` to `cloneRows`. Record it as a fifth diff
shape under the Done clause, or note it as a deviation there. Do NOT weaken
`assertCloneRowsCoversEverySliceField`.

---

## D4 — `TestReq80_UnprovableCoverageCarriesAReasonFromTheClosedSet` must gain the fourth reason

**REQs**: REQ-56, REQ-57.

**The situation**: `internal/graphlint/findings_0006_test.go::TestReq80_…`
asserts the `reason` set is EXACTLY `{dimension-not-finite,
row-can-refuse, tag-not-single-valued}`. REQ-56 appends
`no-participating-dimension` to that closed, append-only set, and C5 states
the REQ-80 test "must stay exhaustive and green".

**Status at Phase 1**: NOT edited. RDR 0010's own
`TestReq56_TheZeroDimensionArmCarriesTheFourthReasonValue` asserts the
four-member set and is RED; the shipped REQ-80 test will go red the moment
Phase 2 makes the append.

**For Phase 2**: extend REQ-80's `want` list with the fourth member. That is
the append 0006 licensed structurally (its set is declared "closed,
**append-only**"), not a weakening — the test stays an EXACT-set assertion.

---

## D5 — `table.IsDecisionTable` is the helper REQ-8 licenses

**REQs**: REQ-7, REQ-8, REQ-69.

**The record says**: "A reader MAY spell the zero value through a small
helper rather than repeating the empty-string comparison, but the field is
the storage".

**What the tests assume**: an exported `table.IsDecisionTable(*Model) bool`.
Four sites read the class (`normalizeRule`, `reach`, `checkDanglingEdge`'s
root arm, `checkCoverage`'s zero-dim arm) and two of them live in
`internal/graphlint`, so the helper must be exported for those readers to
use it at all. The tests assert BOTH halves: the field is the storage
(`TestReq8_TheFieldIsTheStorageForTheClass` mutates `Model.Class` and
requires the helper's answer to follow) and the zero value reads as
`state-machine`.

**For Phase 2**: the name is not fixed by the record. If a different
spelling is chosen, the two table-package tests and the graphlint
`TestReq52_TheLoadedModelCarriesTheClassForLintToRead` need the new name;
nothing else does.

---

## D6 — REQ-100/REQ-101 (SC-6 regression sweep) has no dedicated test

**REQs**: REQ-86, REQ-100, REQ-101.

SC-6 is "every checked-in model and fixture under `make check`" with the
oracle "changed only by licensed diffs". That is the SUITE itself plus a
diff review, not a test a test file can host — writing a Go test that
re-runs `make check` would recurse. REQ-101 is explicitly a negative REQ
("do not write a two-toolchain test") and is discharged by writing none.

Coverage is taken indirectly:
`TestReq71_EveryCheckedInDumpCarryingFixtureNamesEmit` pins the one
non-silent widening across the checked-in fixtures, and
`TestReq85_TheCheckedInNavigatorModelIsUntouched` plus
`mvvStep6NegativeControls`' `models/rdr.toml` arm pin the checked-in model.
The licensed-diff rule itself (REQ-87/REQ-88) is a review obligation on the
Phase-2 diff, recorded in `coverage.md` as such.

---

## D7 — RDR 0011's diff guard fires on this RDR's Phase-2/3 edits

**REQs**: REQ-88 (the scope gate), REQ-102/REQ-104 (IP Phases 1 and 3).

**The situation**: `internal/cli/flow_next_0011_test.go::TestReq120_TheProductionDiffIsFlowNextPlusOneTermInFlowExec`
diffs the working tree against the branch point with `main` and fails when
any production file outside `flow_next.go` / the one `flow_exec.go` term has
changed. Verified by adding this RDR's two field edits — `Model.Class`,
`Row.Emit`, `resolvePayload.Emit` — as temporary stubs: the guard fires
naming `internal/cli/flow_resolve.go` and `internal/table/model.go`.

**Assessment**: this is RDR 0011's scope guard doing its job against a
DIFFERENT RDR's diff. It is not a defect in either RDR, and it is NOT one of
REQ-87's licensed diff shapes — it is a guard whose premise ("the only
in-flight change is 0011's") stops holding once 0010 lands on the same
branch.

**Status at Phase 1**: NOT edited. No test file of this RDR touches it.

**For Phase 2/3**: the guard must learn that 0010's edits are licensed, or
0010 must land on a branch 0011's guard is not evaluated against. Do NOT
weaken the guard to a permissive check — extending its allow-list with the
three files this RDR's Technical Design names (`internal/table/model.go`,
`internal/table/source.go`, `internal/table/load.go`,
`internal/table/normalize.go`, `internal/table/dump.go`,
`internal/graphlint/reach.go`, `internal/graphlint/analysis.go`,
`internal/graphlint/coverage.go`, `internal/graphlint/taxonomy.go`,
`internal/cli/flow_resolve.go`) is the shape that preserves what it checks.
Escalate as a scope question if the allow-list edit is judged out of bounds.

---

## D8 — a shipped 0011 test hard-codes the `flow resolve` payload

**REQs**: REQ-87 shape (ii) — "an added `\"emit\":` payload member".

**The situation**: `internal/cli/flow_demand_0011_test.go::TestReq35And118_ResolveOverTheCheckedInModelAndShippedFixturesIsUnchanged`
compares `flow resolve`'s whole JSON envelope against four hard-coded
expectations. Adding `emit` to `resolvePayload` reddens all four. Verified
with the temporary stub.

**Assessment**: squarely inside REQ-87's licensed shape (ii). Each of the
four expectations gains `"emit":{}` immediately after `"gates":[]` and
changes in no other way — which is also an independent confirmation of C4's
position clause.

**For Phase 3**: make the four edits, verify each changed line differs ONLY
by the added `"emit":{}` member, and do not touch the surrounding
assertions. This is the same class of edit as the 103 `[dump]` fixtures
(ASSUMPTION-8) and should be verified the same way — against the licensed
shape, not re-derived.

---

# Phase 2 (implementation)

## D9 — RDR 0011's diff guard learns 0010's licensed shapes (resolves D7)

**REQs**: REQ-87, REQ-88, REQ-102/REQ-104 (IP Phases 1 and 3).

**Type**: DEPENDENCY-LIMIT. **Status**: mechanical translation.

D7 recorded that `internal/cli/flow_next_0011_test.go::TestReq120_…` fires
on this RDR's production edits. Resolved as D7 itself prescribed: the
guard's allow-list is EXTENDED — never relaxed — with the ten production
files RDR 0010's Technical Design names, plus a second, separately named
predicate for the test-side surface (`_0010_test.go`,
`internal/table/testdata/`, and the three checked-in expectations 0010's
Done clause licenses by shape: `helpers_test.go`, `dump_test.go`,
`findings_0006_test.go`).

The guard's premise — "the only in-flight change is 0011's" — is what
stopped holding; its CHECK is unchanged. A file outside both lists still
fails it, and `internal/resolve` is in neither, so the kernel-untouched
claim the guard exists for is still enforced. Evidence: D7's own prescribed
shape, and the guard's pre-existing `_0011_test.go` exemption, which is the
same "this RDR's own suite" carve-out taken for the sibling RDR.

---

## D10 — the promoted-fixture content pin needs the licensed `emit` append

**REQs**: REQ-71, REQ-87 shape (iii), ASSUMPTION-8.

**Type**: TEST-FIXTURE. **Status**: mechanical translation.

**Found at Phase 2, not recorded at Phase 1.** Beyond the 100 `[dump]`
lists themselves, `internal/table/roundtrip_test.go::TestReq118_PromotedFixtureSetIsNotNarrowed`
leg 2 compares every promoted fixture BYTE FOR BYTE against its approved
iter-2 spike copy under `docs/rdr/0002-…/evidence/spikes/`. Appending
`"emit"` to a `[dump]` order list reddens 100 of them.

The spike copies are RDR EVIDENCE and are never amended. Leg 2 already
carries exactly one recorded normalization for exactly this situation — the
`kind = "string"` → `kind = "scalar"` rename (0002 deviations.md D1) —
applied to the SPIKE side so the promoted fixture stays compared byte for
byte. The `emit` append is added the same way and nowhere else: the
comparison is still byte-for-byte equality, and any other edit to a
promoted fixture still fails.

Grounded in the RDR's own text: REQ-71 states the dump vocabulary is "the
one non-silent widening", A4 censused the 103 fixtures, and REQ-87 shape
(iii) licenses precisely "an added `\"emit\"` member inside a `[dump]`
`order = [ … ]` list". This is that diff reaching one more expectation than
Phase 1 enumerated, not a new class of change.

---

## D11 — `TestReq6`'s class-omitted half must be the zero-owned document

**REQs**: REQ-6, REQ-13, REQ-89.

**Type**: TEST-FIXTURE. **Status**: mechanical translation.

`TestReq6_ClassIsDeclaredNeverInferredFromTheOwnedSet` pairs a DECLARED
decision table with a class-OMITTED model and closes with a fixture
invariant: both must carry the same (empty) owned set, since that is what
makes the pair discriminate an implementation inferring the class from
`len(owned) == 0`. Phase 1 used `dtClassOmitted` for the omitted half — but
D1 gives that fixture an owned tag and its accessors by necessity, so the
invariant it asserts is false of it and the test failed on its own guard.

The pair is repointed at `smZeroOwned` with its `class` line stripped —
the same document `TestReq5`'s "class omitted" subtest already loads, and
genuinely zero-owned. That STRENGTHENS the test: the shared empty owned set
is now real rather than merely asserted, so the discrimination REQ-6 asks
for actually holds. No assertion was relaxed; the closing invariant is
unchanged and still fatal.

---

## D12 — the zero-dimension arm is per-GROUP, and a group's name carries its scoping match atoms

**REQs**: REQ-55, REQ-58, REQ-59, REQ-60, REQ-95.

**Type**: TEST-FIXTURE. **Status**: mechanical translation.

Three Phase-1 assertions encoded fixture assumptions about the group
partition that the shipped partition does not hold:

1. `TestReq55` expected `f.Element == "dt/decide"`. The group's canonical
   name is `guard.Selection.String()`, which carries the match atoms that
   SCOPE the context — for `dtMatchOnly0010` that is `dt/decide a.eq=x`,
   the very `[rule.match.a]` atom whose use instead of a guard atom IS the
   defect being reported. `Element: g.Context.String()` is the spelling
   every other group-scoped finding in `coverage.go` already uses, and
   `dt/decide` names a DIFFERENT group in the same model
   (`dtMatchOnlyBareEscape0010` carries both). The expectation is
   repointed at the real name; the assertion is otherwise unchanged.

2. `TestReq58`/`TestReq60` used `requireOneCode` over
   `dtMatchOnlyBareEscape0010`. That fixture carries TWO zero-dimension
   groups — the bare escape row authors no match atom, so it scopes its own
   `dt/decide` beside `dt/decide a.eq=x` — and C5 binds the arm per GROUP
   ("a decision-table **group** whose scoped product has zero participating
   dimensions"), so two findings is the contract, not a fall-through. Both
   are now asserted BY ELEMENT and BY REASON, which is strictly more than
   the count was checking.

**The clause under test is untouched and still fails without the early
`return`.** Verified: the report over the bare-escape fixture carries
`graph-unprovable-coverage` and nothing else — no `graph-coverage-gap`, no
`graph-coverage-closed-by-escape`, no `graph-product-too-large`. REQ-58's
count check is replaced by an EXACT code-set equality, so a fall-through
that emitted the closure advisory still reddens it. Evidence: the shipped
partition (`internal/guard/product.go::selectionOf`, which builds the
context from a row's match atoms) and `0006:C7`/`0003:C13`, which the RDR
cites in the same terms ("match keys select a row's group").
