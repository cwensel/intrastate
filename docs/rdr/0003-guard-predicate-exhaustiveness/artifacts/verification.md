# Verification — RDR 0003 Guard Predicate Exhaustiveness

Phase 3a CoVe artifact. Findings below were produced by driving the shipped
implementation (`internal/table` loader + `internal/guard`) with probe models
authored against the REQ list, WITHOUT reading the Phase 1 test files. Each
entry names the REQ it violates, the exact failing input, and observed vs
spec-required behaviour.

---

## FAIL-1 — a can-refuse row is credited with an accepted-assignment set

- **Violates**: REQ-43 (`0003:C17`), and REQ-44's scoping of the surviving
  overlap check.
- **Quoted text**: "**A can-refuse row contributes no assignments to the
  coverage union** … a row that can refuse has **no decidable
  accepted-assignment set to contribute**, so lint MUST NOT credit it with its
  `all`-intersection unsubtracted"; and "Overlap findings among the group's
  **decidable** rows are unaffected and still MUST be emitted."

### Failing input

One scoped row group. `opt` is declared **optional** (no `required` marker —
the conservative default REQ-20 fixes); `req` is always-present. Both are
single-valued enums, so every dimension is finite and projectable.

```toml
[tags.opt]
provenance = "observed"
kind = "enum"
domain = ["x", "y"]
single_valued = true
# no `required` -> optional

[tags.req]
provenance = "observed"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[[rule]]
id = "dec"
[rule.match.recognized]
eq = "go"
[rule.guard.all.req]
eq = "a"
[rule.write]
mark = "x"

[[rule]]
id = "refuser"
[rule.match.recognized]
eq = "go"
[rule.guard.all.opt]
eq = "x"
[rule.write]
mark = "y"
```

`refuser` carries a **value atom over a key declared optional**, so it "can
refuse" under REQ-70, and `guard.CanRefuse` correctly reports `true`.

### Observed behaviour

```
GROUP "p/go" dims=[opt req] presdims=[opt]
  cardinality=8 finite=true productProjectable=true productLen=8
  row dec        canRefuse=false accepted.projectable=true  accepted.len=4
  row refuser    canRefuse=true  accepted.projectable=true  accepted.len=2
  coverageUnion projectable=true len=5
```

and from `guard.Lint`:

```
GROUP "p/go" rules=[dec refuser] verdict=withheld green=false provable=false
   finding code=graph-unprovable-coverage rules=[refuser] dim="opt" blocking=true atom={opt all eq [x]}
   finding code=graph-overlap rules=[dec refuser] class=""
```

Two spec-visible symptoms of one root cause:

1. `guard.AcceptedAssignments` / `guard.CoverageUnion` return a **projectable**
   set for `refuser` (2 assignments) and a union of 5 that **includes them**.
   REQ-43 requires a can-refuse row to have no decidable accepted-assignment
   set and to contribute nothing to the union.
2. `graph-overlap` is emitted for the pair `(dec, refuser)`. REQ-44 scopes the
   surviving overlap check to "the group's **decidable** rows"; `refuser` is
   not one, so this pair must never be intersected. Adding a second decidable
   row reproduces it at scale — a three-row group emits `(d1,d2)`, `(d1,refuser)`
   and `(d2,refuser)` where only `(d1,d2)` is a real finding.

### Spec-required behaviour

`refuser` yields the unprojectable set; `CoverageUnion` returns 4 (the
decidable rows only, or the unprojectable set); and `graph-overlap` is emitted
for decidable pairs only — none in the two-row model above.

### Mechanism

`product.go::acceptedIn` returns the unprojectable set only when the group's
product is unprojectable or when some atom's `Denotation` is unprojectable. It
never consults `CanRefuse`. When every dimension is finite and single-valued —
which is exactly the narrowing's own interesting case — a can-refuse row's
atoms all project, so the row is credited. `lint.go::coverageUnion` skips on
`!accepted.Projectable()` and its comment claims "A can-refuse or unprojectable
row contributes NO assignments", but the can-refuse half of that comment is
never reached. `AcceptedAssignments`' doc comment ("a can-refuse row and a row
carrying an unprojectable atom both return the unprojectable set") is likewise
not what the code does.

`Lint`'s group verdict is unaffected — it withholds on `CanRefuse` before
reading the union, so **no false green is reachable through this path**. The
violation is in the exported denotation/union API and in the emitted overlap
set, both of which the REQs constrain directly.

---

## Probes that PASSED (no finding)

Recorded so the density of this pass is legible.

- Finite-domain exhaustiveness: complete partition green; one-sided gap; gap
  and overlap visible only in the 2-D product both reported from one run with
  a concrete witness (REQ-60, REQ-95, REQ-99, REQ-129).
- Single-valued marker: with/without the marker on otherwise identical rows
  returns different verdicts (`exhaustive` vs `withheld`) — REQ-136, REQ-26,
  REQ-58.
- Operator/kind matrix, cell by cell, incl. `eq`/`in` refused on `set`,
  `contains` refused on `enum`, `lt` refused on `enum`, `exists` on all five
  (REQ-4, REQ-3, REQ-6).
- Possibly-absent-key narrowing: withheld with the row and atom named, one
  finding per refusing row, in `all` and in `unless` alike; the always-present
  negative control certifies green (REQ-42, REQ-70, REQ-73, REQ-96, REQ-139,
  REQ-140, REQ-MVV 4 and 5).
- `exists` presence dimension: `{present,absent}` pair closes coverage;
  `exists=true` alone gaps at `{absent}`; vacuous over an always-present key
  and reported rather than rejected (REQ-52, REQ-55, REQ-56, REQ-5).
- Set containment over a declared element universe: powerset assignment space,
  `contains`/`unless contains` partition, gap, optional-key withholding
  (REQ-17, REQ-89 `set` row, REQ-124).
- Bounded integer edges: `lt`/`gte` partition, `lte`/`gt` partition, the
  off-by-one gap at the boundary, the `lte`/`gte` overlap at the boundary,
  inclusive endpoints (`{0..3}` → 4), negative ranges, unbounded (`min` only)
  → unprovable (REQ-18, REQ-19, REQ-61, REQ-62).
- Published bound: 2048, reported beside the computed size; the REQ-131 shape
  pair (two 64-value dims vs twelve booleans, both 4096) receives the same
  refusal, and both prove under the bound (REQ-86 – REQ-92, REQ-108, REQ-132).
- No computed size for a product carrying an unprovable dimension, and one
  blocking finding per unprovable dimension (REQ-93, REQ-94, REQ-107).
- Mixed `all`/`unless`: the block is subtracted as one conjunctive set, not
  per-atom negation (REQ-39, REQ-41).
- Escape rows: bare escape closes only its declared class and is reported as
  `graph-coverage-closed-by-escape`; per-class overlap populations, a row in
  several classes placed in each, one finding per shared class, no finding for
  rows sharing no class, and no escape-vs-ordinary overlap (REQ-75 – REQ-82,
  REQ-109).
- Identity tuple: `all`/`unless` over one key distinguishable; literal
  distinguishes; rule id distinguishes; delimiter-stuffing and `["a","b"]` vs
  `["ab"]` stay distinct; semantic equality drops block and source (REQ-8,
  REQ-9, REQ-10, REQ-138).
- Source-order independence, including reordered `in` literal elements and
  reordered rows (REQ-83, REQ-84, REQ-137).
- Declaration errors at the loader as `malformed_tag_declaration`: `min`/`max`
  on an `enum`, `elements` on any non-`set`, `single_valued` on a `set` or a
  `scalar`, `domain` on `bool`/`int`/`set`, an unknown kind token (REQ-28,
  REQ-29, REQ-31, REQ-32).
- Atom rejections as `malformed_predicate_atom` / `unknown_tag`: unknown
  operator, literal parse mismatch on `int`/`bool`, non-boolean `exists`
  literal, non-array `in`, multi-member `eq`, array comparison bound, empty and
  duplicate-bearing set literals, and **literal-outside-declared-domain** for
  both `enum` and `int` (REQ-2, REQ-30, REQ-112).
- Evaluator seam: value semantics over a present value only; a present value
  the operator cannot parse is unevaluable rather than false; `exists` and any
  unknown operator return unevaluable (REQ-34, REQ-35).
- Participation: a key only one row constrains still enters the product; a key
  every row constrains identically still bounds it; match keys never enter; the
  same key is a match key in one group and a guard dimension in another
  (REQ-48 – REQ-51).
- Provenance: `graph-owned-before-write` for a matched owned tag with no
  reachable write, silent once written, and never raised for an observed tag
  (REQ-100, REQ-102).

## Read but not a FAIL

- `single_valued = true` on an `enum` or `int` declaring **no** finite domain
  is accepted. REQ-134 names "a single-valued marker on a `set` kind or on a
  kind carrying no finite domain", but REQ-29's agreement table is per-KIND and
  enumerates exactly `set` and `scalar`, which the loader enforces. The
  mandatory `[tags.recognized]` declaration is itself an `enum` with
  `single_valued` and no `domain`, so the per-kind reading is the only one the
  model admits. Spec ambiguity, not a deviation.
- A scoped row group carrying no guard atom is reported with no verdict and no
  green claim. The RDR is silent; `deviations.md` D7 records the choice.

# Phase 3b — Adversarial Review (RDR 0003)

Independent adversarial pass. Three failure modes named, each anchored in
the record's `Trade-offs` sections, each with a test that FAILS against the
current implementation. Tests live in
`internal/guard/guard_adversarial_0003_test.go`. No existing test was
edited or weakened. Defects are recorded, not fixed — Phase 3c owns the
fix.

Validated with `go test ./...`: the three adversarial tests are the only
failures in the repository; every pre-existing test still passes.

## ADV-1 — A green claim covers a conforming view the runtime refuses

- **Failure mode**: False exhaustiveness claim. Lint certifies a scoped row
  group `exhaustive` / green, and its coverage union contains a conforming
  view on which the kernel refuses `guard_unevaluable`.
- **RDR anchor**: `Trade-offs / Failure Modes` — "Silent failure would be a
  false exhaustiveness claim; the recovery path is to keep every exactness
  claim tied to A2 and the MVV fixture." Also narrows against REQ-67 ("An
  exhaustiveness claim MUST NOT be stronger than the runtime it describes").
- **Test**: `TestAdv1_GreenClaimCoversAViewTheRuntimeRefuses`
- **Currently**: **FAILS**
- **Mechanism**: A `set` tag's value dimension is the powerset of its
  element universe, so the scoped product enumerates the EMPTY subset as an
  assignment, and `guard.Conforms` admits `caps = []` as a conforming view.
  But `grammar.go::parseSetLiteral` refuses an empty JSON array — the
  published literal shape for both set-shaped operators is a NON-EMPTY typed
  element set — so `Evaluator.Evaluate` returns `GuardUnevaluable` for a
  held `[]`. `product.go::valueSatisfies` compares that verdict against
  `GuardTrue` only, collapsing UNEVALUABLE and FALSE into one answer. The
  `unless contains` complement row is therefore credited with the
  empty-subset assignment, closing coverage and certifying green, while the
  kernel refuses `guard_unevaluable` on that same view.
- **Minimal reproduction**: one always-present `set` tag over a
  ONE-element universe, partitioned by `contains = ["x"]` in `all` and the
  same atom in `unless`. Lint: `verdict=exhaustive green=true`, no findings.
  Runtime on `caps=[]`: `refusal:guard_unevaluable`.
- **Note for Phase 3c**: the two surfaces disagree about whether an empty
  held set is a value at all. Either the empty subset leaves the product
  (and a `set` tag is asserted non-empty), or `valueSatisfies` must
  propagate UNEVALUABLE as unprojectable rather than as false. The test
  asserts both halves — the runtime agreement and the lint-side
  acceptance — so a fix that only silences one surface will not pass it.

## ADV-2 — The over-large refusal enumerates the product it refuses

- **Failure mode**: Naive powerset enumeration on the refusal path. The
  `graph-product-too-large` refusal is emitted, but only after the full
  powerset of the element universe has been materialized, so the published
  bound provides no protection.
- **RDR anchor**: `Trade-offs / Risks and Mitigations` — "**Risk**:
  Set-valued domains are implemented by naive powerset enumeration.
  **Mitigation**: Require symbolic or bitset-equivalent proof and explicit
  refusal/downgrade for finite products that are too large to prove."
- **Test**: `TestAdv2_OverLargeSetProductIsRefusedWithoutEnumeratingIt`
- **Currently**: **FAILS**
- **Mechanism**: `Cardinality` and `Product` both decline correctly and in
  microseconds. `Denotation` consults no bound at all: it calls
  `valueAssignments`, which calls `subsets(d.Elements)` and builds all 2^n
  canonical JSON arrays. `lint.go::unprovableReason` invokes `Denotation`
  for every guard atom on every participating dimension, so the enumeration
  runs before the bound is ever compared — on the refusal path itself.
- **Measured** (this implementation, `Bound() == 2048`): n=16 → 32,768
  assignments in ~160ms; n=20 → 524,288 in ~2.3s; n=22 → 2,097,152 in
  ~10s. n=22 is a factor of 1024 above the bound, materialized in memory,
  from an authored model of twenty-two element names and one `contains`
  atom. Scaling is 2^n.
- **Assertion shape**: the test compares lint cost for two universes four
  elements apart (14 and 18), both far past the bound and both refused. A
  bounded refusal decides each from declaration arithmetic and costs about
  the same; naive enumeration costs 16x. Tolerance is 8x with a 50ms
  absolute floor; the observed ratio is a stable 16-19x across runs.
- **Note for Phase 3c**: `product.go::Denotation` needs the same guard
  `Product` already applies. This also contradicts the documented rationale
  for `bound` — "the largest product this implementation's enumerating
  proof representation completes over within its budget" — and weakens A15,
  which claims the cardinality bound predicts provability.

## ADV-3 — An unprovable dimension silences an unrelated, decidable overlap

- **Failure mode**: Overlapping candidate rows — a visible lint failure the
  record enumerates — go unreported. One opaque `scalar` key on an
  unrelated row suppresses the overlap check for the entire scoped row
  group.
- **RDR anchor**: `Trade-offs / Failure Modes` — "Visible failures should be
  typed load or lint failures: … non-exhaustive finite domain, **overlapping
  candidate rows**, guard dimension not provable because it lacks a finite
  domain …". The implementation states the same obligation in
  `lint.go::lintGroup`: "overlap among the group's decidable rows is
  unaffected and still MUST be [reported]".
- **Test**: `TestAdv3_UnprovableDimensionSilencesAnUnrelatedOverlap`
- **Currently**: **FAILS**
- **Mechanism**: `Product` returns the unprojectable set as soon as ANY
  participating dimension lacks a finite domain. `acceptedIn` short-circuits
  on an unprojectable product and returns the unprojectable set for EVERY
  row in the group — including rows that constrain only fully provable
  dimensions. `pairwiseOverlaps` then skips all of them. Unprojectability is
  being scoped to the group rather than to the row that carries the
  offending atom.
- **Reproduction**: one group with `adv3-dup-a` and `adv3-dup-b` carrying
  byte-identical `p eq true` guards (`p` is `bool`, `single_valued`,
  `required`) plus `adv3-opaque` guarding on a `scalar` key. Lint emits
  exactly one finding — `graph-unprovable-coverage` on dimension `s` — and
  ZERO `graph-overlap`. The kernel refuses `ambiguous_match` on
  `{p:true, s:z}`.
- **Distinct from existing coverage**: `TestReq44_WithheldGroupEmitsNoCoverageGapButStillEmitsOverlap`
  exercises the OTHER withholding path — a can-refuse row over an optional
  key — where the product stays projectable and overlap still computes
  correctly. This is the unprovable-dimension path, which no existing test
  reaches.
- **Consequence**: the author is told a dimension is unprovable and is NOT
  told two of their rows collide, which is the defect they can actually fix.
  The findings are correct but incomplete, and the record requires every
  decidable defect in one pass (REQ-95).

## ADV-4 — (recorded, latent) Inverted int bounds on the exported surface

- **Failure mode**: `guard.AssignmentCount` reports a NEGATIVE cardinality
  as a valid finite count, and panics with "negative shift amount" for the
  unmarked (co-occurring) variant.
- **RDR anchor**: `Trade-offs / Failure Modes` — the "literal parse failure"
  / declaration-rejection family, and the declaration model's own clause
  that the kind/field agreement rule is read "on this RDR's own surface, so
  a declaration built in memory is judged by it too"
  (`declaration.go::agrees`).
- **Status**: **LATENT — not asserted as a failure.** `table.Load` rejects
  `min > max` (`load.go:213`), so no declaration reaching the guard package
  through the loader can trigger it. `guard.AssignmentCount` is exported and
  takes a `table.TagDecl` directly, so the defect is reachable by any
  in-package or future caller that builds a declaration in memory — which
  the existing suite already does in several places.
- **Observed**: `AssignmentCount({Kind:"int", Min:3, Max:0, SingleValued:true})`
  returns `(-2, true)`; the same declaration without `SingleValued` panics
  in `spread` on `1 << -2`.
- **Test**: `TestAdv_RecordedIntBoundInversion` — asserts only the
  reachability premise (that the loader still refuses the declaration), so
  it PASSES today and turns red the day the loader stops holding the line.
  It is not a vacuous adversarial test: it is a tripwire on the invariant
  that keeps ADV-4 latent.
- **Note for Phase 3c**: `agrees` does not check `*Min <= *Max`, and
  `domainSize` does not floor the computed range. Whether to harden this is
  a judgement call — it is defence in depth behind a loader check that
  currently holds.

---

# Phase 3c — fixup outcomes

How each defect was closed. Every fix is the minimum change that makes the
governing clause true; no existing assertion was edited, relaxed, or
deleted. New deviations D10–D15 in `deviations.md`.

## FAIL-1 — CLOSED

- **Root cause**: `acceptedIn` consulted only PROJECTION — the group's
  product and each atom's `Denotation` — and never `CanRefuse`. The
  withholding decision and the projection decision are independent: when
  every dimension is finite and single-valued, which is the narrowing's own
  interesting case, a can-refuse row's atoms all project.
- **Fix**: `acceptedIn` returns the unprojectable set for a row that can
  refuse. One condition, at the head of the function. Both doc comments
  (`AcceptedAssignments`, `coverageUnion`) already described this.
- **Effect**: the row contributes nothing to the union (REQ-43) and is never
  paired by `pairwiseOverlaps` (REQ-44). The withholding finding still
  fires, so the fix cannot be mistaken for dropping the row from lint.
- **Regression test**: `TestFail1_CanRefuseRowContributesNoAcceptedAssignments`
  — verified to fail on all three symptoms before the fix.

## ADV-1 — CLOSED (one assertion outstanding; see D15)

- **Root cause**: NOT where Phase 3b placed it. The defect is in the
  EVALUATOR, not in lint's projection. `Evaluate`'s `contains` arm parsed
  the HELD value with `parseSetLiteral`, whose non-empty requirement is the
  operator/kind matrix's published shape for the AUTHORED LITERAL (REQ-4).
  A held `[]` was refused as unevaluable, so the kernel refused
  `guard_unevaluable` on a conforming view lint's own product enumerates.
- **What settles it**: REQ-57 — `contains` denotes "the assignments whose
  held set CONTAINS every listed element". Containment is TOTAL over sets;
  `[]` does not contain `x`, so the verdict is FALSE, not undecidable. The
  literal's shape rule never governed the held value.
- **Fix**: `parseHeldSet` decodes a held §D13 array with no cardinality
  rule; `parseSetLiteral` keeps the non-empty rule for the authored literal.
  `Denotation` additionally propagates the seam's three-valued verdict
  rather than collapsing UNEVALUABLE into false — reading a three-valued
  seam two-valued is the general shape of the defect, and the `contains` arm
  was the reachable instance.
- **Both Phase 3b candidates were weighed and rejected**: dropping the empty
  subset makes a `set` dimension `2^n-1`, contradicting REQ-89, whose own
  test pins the powerset; propagating UNEVALUABLE as unprojectable ALONE
  (without fixing the evaluator) was implemented, measured, and reverted —
  it broke four Phase-1 tests (REQ-84, REQ-129, REQ-57, and the MVV
  gap-and-overlap scenario) by making every `contains` dimension
  unprojectable, since `[]` is a member of every powerset.
- **Verified outcome**: the runtime now resolves `caps=[]` to `adv1-not` and
  `caps=["x"]` to `adv1-has` — a genuine complete partition. The group's
  green claim is TRUE: every conforming view it covers, the runtime decides.
- **Regression test**: `TestFixup_EmptyHeldSetIsDecidedRatherThanRefused` —
  verified to fail against the pre-fix evaluator.
- **Outstanding**: `TestAdv1_…`'s second assertion still fires, on a message
  whose own premise is now false. It was NOT edited. See D15 — the only
  entry needing an author decision.

## ADV-2 — CLOSED

- **Root cause**: the bound was compared AFTER the work rather than before
  it. `Denotation` consulted no bound, calling `valueAssignments` → the full
  powerset; `unprovableReason` reached `Denotation` for every guard atom, so
  the enumeration ran on the refusal path itself.
- **Fix**, three parts, one cause:
  1. `valueAssignments` declines a dimension whose own assignment count
     exceeds `Bound()` — a dimension larger than the whole product bound can
     appear in no product this implementation enumerates (D11).
  2. `lintGroup` tests the bound BEFORE the projection scan and after every
     dimension is known finite, the order REQ-93 states. `graph-product-too-large`
     still carries `(computed size, bound)`; a product with an unprovable
     dimension still reports no computed size (REQ-94).
  3. `spread` saturates at the same ceiling `Cardinality` saturates at (D12).
- **Effect**: guard suite wall-clock 1.6s → 0.4s. No verdict changed.
- **Regression test**: `TestFixup_HugeSetUniverseIsRefusedRatherThanWrappingUnderTheBound`
  — covers the D12 wrap at |universe| ∈ {40, 63, 64, 65, 96}; verified to
  fail (OOM/hang) against the unguarded shift.

## ADV-3 — CLOSED

- **Root cause**: unprojectability was scoped to the GROUP rather than to
  the dimension carrying it. `Product` returned the unprojectable set as
  soon as any dimension lacked a finite domain, and `acceptedIn`
  short-circuited on that, so one opaque `scalar` key made every row in the
  group undecidable.
- **Fix**: `acceptedIn` projects onto the group's DECIDABLE sub-product (D10).
  A row carrying an atom over an undecidable dimension still yields the
  unprojectable set. Where every dimension is decidable the sub-product IS
  the scoped product, so the provable path is unchanged; coverage is still
  compared against the FULL product.
- **Grounding**: the MVV Scenario-2 walkthrough (`0003:1408`) declines
  overlap because "step 4 produced none FOR THE TWO ORDINARY ROWS" — the
  rows carrying the unprojectable atoms — which is the per-row reading.
- **Regression test**: the Phase 3b `TestAdv3_…`, now passing on both its
  halves (the emitted finding AND the accepted-assignment surface).

## ADV-4 — FIXED, not left latent

- **Weighed**: the Phase 3b note called hardening "a judgement call — it is
  defence in depth behind a loader check that currently holds". Fixed,
  because the cost is one comparison, `AssignmentCount` is EXPORTED and
  takes a `table.TagDecl` directly, and `agrees`' own doc comment already
  claims the rule is read "on this RDR's own surface, so a declaration built
  in memory is judged by it too" — a claim the code did not honour.
- **Fix**: `agrees` requires `*Min <= *Max` (D13); `spread`'s floor (D12)
  independently removes the `1 << -2` panic. REQ-18 fixes both endpoints
  inclusive, which presumes an ordered pair — `{3..0}` names no value, so it
  carries no readable finite domain.
- **Tripwire**: `TestAdv_RecordedIntBoundInversion` is UNCHANGED and still
  green. It asserts the loader premise, which still holds.
- **Regression test**: `TestFixup_InvertedIntBoundCarriesNoReadableDomain` —
  asserts no count is reported, none is negative, and the ordered `{0..3}`
  sibling still counts 4.

## Suite state

`go build ./...` clean; `golangci-lint run` 0 issues; `go test -race`
clean. `internal/resolve`, `internal/table`, `internal/cli` green. In
`internal/guard`: 148 tests pass — the 142 Phase-1 tests, the four
Phase-3c regressions, `TestAdv2_…`, `TestAdv3_…`, and the `TestAdv_Recorded…`
tripwire. One assertion fails: `TestAdv1_…`'s second half, unedited, per
D15. No pre-existing test was edited, weakened, or broken.
