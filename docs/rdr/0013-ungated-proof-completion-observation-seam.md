# Recommendation 0013: Ungated proof-completion observation seam

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-08-28
- **Status**: Draft
  <!--
  - `Deferred` is the parked-with-a-revisit-trigger status for a
    Draft that cannot proceed because **no acceptable mechanism
    exists yet** — every in-our-control path is ruled out and the
    one that would work is outside our control. It is a *pause in
    the lifecycle*, not an exit from it: the RDR stays intact and
    re-enters at the stage it stopped when the trigger fires.
    Carry the condition on the live value:
    `Deferred [revisit when <condition>]`, and say in the same
    field what was ruled out and why (Alternatives Considered
    carries the long form). Distinct from `Abandoned`, which is
    terminal — an Abandoned RDR is closed, owes a post-mortem, and
    never re-enters. A Deferred RDR owes **no** post-mortem
    (nothing was implemented), and its `Priority` records what the
    fix is *worth*, not what is scheduled. Do not defer merely to
    park work that is possible but unfunded — that is a Priority,
    not a Status.
  - `Demoted` is the terminal status for an RDR judged
    *not RDR-shaped* — the decision was never a real
    design fork, so it leaves the RDR lifecycle and is
    refiled as a plain issue. Carry the destination on the
    live value: `Demoted [→ <issue link>]`, and record the
    same link under **Related Issues**. A `Demoted` RDR runs
    no further stages. (Distinct from the 07.1 *demotion*
    below, which is a `Final → Draft` flip that keeps the
    RDR in the lifecycle — that flip never writes
    `Status: Demoted`; see the disambiguation note there.)
  - A Draft demoted from Final by the 07.1 cluster gate
    carries a qualifier on the live value:
    `Draft [revised from Final YYYY-MM-DD; re-verify A2,A4
    — <one-line reason>]`. It is still a `Draft` for every
    binary Draft/Final gate; only Stage 4 (scoped
    re-verify) and Stage 7 (re-lock) parse the qualifier.
    The Stage 7 flip to `Final` overwrites the whole value,
    so the qualifier self-clears at re-lock — no separate
    cleanup. This 07.1 "demotion" is a *verb* describing the
    Final→Draft flip; it is **not** the `Demoted` status
    above (which exits the lifecycle to an issue) — do not
    conflate the two. (`Reverted` above is the unrelated
    terminal "implementation rolled back" status — also do
    not conflate.)
  - A Final tolerated at the 07.1 gate under a JOINT-DECISION
    carries `Final [joint decision → <home §-anchor>: <the
    open question>]`. It is still a `Final` for every binary
    gate. The qualifier is an **open obligation, not a
    coherence claim**: it says the named question is
    unanswered here, not that this RDR agrees with the answer.
    So it does not self-clear. When the home answers, this RDR
    owes a scoped check of that answer against its own
    normative fences before it re-locks or implements —
    consistent → drop the qualifier and record the clearing;
    contradicts fenced text → a 07.1 SPEC-DEFECT. A re-lock
    that comes first carries the qualifier forward unchanged;
    it is never silently dropped.
  -->
- **Type**: Technical Debt
- **Profile**: mid — provisional: one contract (the
  proof-completion observation seam, or the recorded decision to
  keep the bound's roles fused), either arm amending locked
  RDR 0003 surfaces (A15, REQ-133, `ProofCompletes`).
  <!-- Do not paste the matrix below into the field; it is the
  Stage 5 routing latch, provisional on `Draft`, made
  authoritative by Resolve.
  Sized by BLAST RADIUS — the MAX of two axes, not
  contract count or word count.
  (1) contract axis: small = one contract, no user-facing
  surface (skips Stage 5); mid = one contract + user-facing
  surface OR locks a contract; large = locks an enum/hash/
  format/grammar/destructive-op; foundational = cross-RDR
  producer / spans modules.
  (2) accretion axis (HARD floor): if `Seam Lineage` below
  carries ≥2 closed prior point-fixes at this locus, Profile
  is floored at FOUNDATIONAL regardless of the contract axis
  — a seam with prior point-fixes is never small/mid (it
  spans the prior RDRs/patches = the matrix's cross-RDR
  trigger). The only escape is a written accretion disposition
  in the Seam Lineage field. This floor is what stops a
  "one contract → mid" sizing from under-gating an accreting
  seam.
  Matrix: rdr/stages/README.md. Seed estimates from the design
  shape; Resolve overwrites from the verified count; Stage 8
  Gate locks it at Draft → Final. Never skip lenses off a
  Draft Profile until Resolve has run. -->
- **Priority**: Low
- **Related Issues**: kata `intrastate#dh40` (1531); katas
  `s1r1` (shipped tautology removal), `9yeq` (removed wall-clock
  ratio), `ge67` (sequencing constraint)
- **Predecessors**: 0003-guard-predicate-exhaustiveness
- **Seam Lineage**: no prior accretion

## Problem Statement

A maintainer locking RDR 0003 needs its Pending assumption `0003:A15`
— a single published cardinality bound predicts provability — to be
discharged by evidence that could refute it. Today it cannot be: every
exported enumeration seam (`internal/guard/product.go::Product`,
`internal/guard/assignment.go::valueAssignments`) consults `Bound()`
before enumerating, so REQ-133's completion oracle can only observe
agreement the bound itself defines. Demonstrated: lowering the bound
2048→64 — an over-conservative value that should read as a
misprediction — still passes REQ-133. A bound that does not actually
predict provability would ship as if verified.

The decision: `Bound()` currently fuses three roles — the published
constant (REQ-86's model-independence carrier), the enumeration gate,
and the completion oracle — and the record declares them one quantity
deliberately (`internal/guard/product.go` doc, deviation D6). Either
**split** the roles by naming a second, bound-independent observation
surface (an ungated enumerator entry point or an injectable enumerator,
with a declared *structural* stopping condition), or **keep them
fused** and downgrade A15 to a definitional design constraint, changing
its Method from `MVV Test` and restating REQ-133 as a consistency
check. Two adjudicated exclusions bound the fork: the observation must
not be derived from the bound under test (A15's own Plan), and it must
not be a wall-clock, memory, or subprocess-cap oracle (REQ-86/ADV-2;
kata `9yeq` already shipped the removal of a wall-clock ratio). Either
arm rewrites A15's Status/Method, the `ProofCompletes` doc contract,
and REQ-133's assertion shape — the choice cannot be deferred to the
test.

## Critical Assumptions

- **A1 A test-scoped export (`export_test.go` in `internal/guard`)
  reaches the unexported ungated entries from the external
  `package guard_test` files without growing the production API.**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: needed — a compile-and-run spike: an
    `internal/guard/export_test.go` exposing an unexported symbol to
    `guard_mvv_0003_test.go`, plus `go build ./...` showing the symbol
    absent from the non-test build. Precedent exists in a sibling
    package (`internal/table/export_test.go`), so the spike confirms
    reach from THIS package's external tests rather than the mechanism
    itself.
  - **If wrong**: the observation seam needs a production export after
    all, which reopens the public-surface question dh40 was routed out
    over, and the choice falls back to Alternative 1.
- **A2 Structural exhaustion is observable at test-feasible cost for the
  over-bound fixtures: the `2*Bound()`-cardinality shape pair and the
  2^18-subset value dimension enumerate to exhaustion in ordinary test
  time.**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: needed — timed run of the ungated enumeration over
    `shapePairSource(_, true)` and `adv2BigSetSource(18)` fixtures
    (an emission/termination claim; never quote-confirmable from source
    reading).
  - **If wrong**: the over-bound side of the differential stays
    unobserved and REQ-133's restatement narrows to the admitted region
    plus a smaller over-bound fixture chosen at implementation.
- **A3 After limit-parameterization, `Bound()` is the only value any
  production call site passes as the enumeration limit, so REQ-86's
  published-constant posture (not per-model, per-group, or per-run
  configurable) is preserved.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: needed — a call-site sweep over `internal/guard`
    showing every non-test caller of the parameterized entries passes
    `Bound()`.
  - **If wrong**: the seam silently becomes the injectable gate
    Alternative 2 rejects, and two runs of the same model could apply
    different limits.
- **A4 The definitional disposition of the bound's value is consistent
  with RDR 0003's own record: REQ-86 mandates *publishing* the bound and
  verdict determinism, and deviation D6 already records the value 2048
  as an implementation choice — neither clause states the value is an
  empirical prediction subject to refutation by over-conservatism.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: needed — `0003:C24` (REQ-86's clause) read against
    this RDR's C3, plus 0003's `artifacts/deviations.md` §D6; the
    reading must also dispose of the tension between REQ-86's
    "same verdict on every conforming implementation" sentence and the
    `internal/guard/product.go::bound` doc's "another conforming
    implementation MAY publish a different bound".
  - **If wrong**: the bound's value is adequacy-bearing after all, an
    independent resource observation is owed, and the REQ-86/ADV-2
    exclusions leave no admissible oracle — the fork reopens toward
    Alternative 3's pure consistency check.
- **A5 Peer bounded-verification practice treats the analysis bound as a
  declared scope rather than a refutable prediction (the Alloy-scope /
  bounded-model-checking stance the class prior art suggests).**
  - **Status**: Pending
  - **Method**: Prior Art
  - **Evidence**: needed — the resolving corpora had no coverage for the
    instance question at Propose (see
    `evidence/research/propose-prior-art.md`); the quoted TraceFix p9
    passage ("outside those bounds the guarantee is incomplete") is the
    only anchored support, so the broader peer-tool claim is demoted
    here rather than leaned on.
  - **If wrong**: C3's definitional framing loses its external support
    and rests solely on 0003's D6 — the choice survives (D6 suffices)
    but the rationale must drop the peer-practice sentence.

## Proposed Solution

### Approach

Split the observation from the gate by making the enumeration limit an
explicit caller-supplied parameter at the two existing gate sites — the
product assembly in `internal/guard/product.go::Product` and the
per-dimension precheck in
`internal/guard/assignment.go::valueAssignments` — with every production
call site passing `Bound()` and the ungated entries reaching tests only
through a test-scoped export (`export_test.go`). REQ-133's completion
half is restated as a **structural differential**: the ungated
enumeration, stopped by a cap derived from declaration arithmetic
(`internal/guard/product.go::Cardinality`, which consults no bound) and
never from `Bound()`, must exhaust and yield exactly the arithmetic
count — on **both** sides of the published bound. The half of 0003:A15
that this cannot observe — whether 2048 is *adequate* as a resource
budget — is closed as definitional, adopting the keep-fused arm's
insight: the bound's value is a published analysis scope (0003 deviation
D6), so an over-conservative value is a smaller scope, not a
misprediction. The seed's 2048→64 demonstration therefore stops being a
blind spot by adjudication, not by detection: what becomes detectable is
an enumerator that fails to exhaust an admitted product, yields a count
disagreeing with `Cardinality`, or behaves differently on
equal-cardinality shapes — the refutable content A15 actually carries
for an enumerating representation.

### Technical Design

Three facets of one seam (the proof-completion observation contract):
the mechanical seam (C1), the restated oracle (C2), and the recorded
disposition of the bound's value (C3). Data flow: production verdicts
are unchanged — `Product`/`Denotation` still decline from declaration
arithmetic before enumerating (`ADV-2` posture intact); only the test
binary gains a second, bound-independent route to the same enumerator
(`internal/guard/product.go::cross`, which is already ungated — the gate
sits in its callers).

#### Normative Contracts

These three blocks are facets of one contract — the proof-completion
observation seam — not three independent seams (the split signal counts
one).

**C1**

```normative
The enumeration limit is an explicit caller-supplied parameter at the
two gate sites (the product assembly reached from `Product`, and
`valueAssignments`' per-dimension precheck). Every production call site
passes `Bound()` and no other value; the published constant remains the
single production limit — not per-model, per-group, or per-run
configurable (REQ-86 posture). The ungated entries stay unexported;
tests reach them only through a test-scoped export
(`internal/guard/export_test.go`), so the conforming production surface
is unchanged and the ungated route is absent from the non-test build.
```

**C2**

```normative
REQ-133's completion observation is structural exhaustion, never a
resource oracle: the ungated enumeration runs under a stopping cap
derived from declaration arithmetic (`Cardinality` /
`AssignmentCount`), never from `Bound()`, and the observation is
(exhausted, distinct-assignment count). Assertions:
(a) soundness — every product the bound ADMITS exhausts with count ==
    `Cardinality(m, g)`;
(b) shape-independence — equal-cardinality shape pairs (few wide vs
    many narrow dimensions) yield identical (exhausted, count) on BOTH
    sides of the published bound;
(c) dimension side — a value dimension whose domain exceeds the bound
    (the ADV-2 powerset fixture) exhausts ungated to exactly its
    arithmetic `domainSize`, while production `Denotation` still
    refuses it.
An over-bound product exhausting structurally is headroom, not a
refutation (see C3). Wall-clock appears at most as test-harness hang
protection, never as the oracle (REQ-86 / ADV-2 exclusions).
```

**C3**

```normative
The published bound's VALUE is a definitional analysis scope
(0003 deviation D6), not an empirical prediction: `ProofCompletes`
remains defined as `card <= Bound()`, its doc restated to say the
budget IS the bound by definition rather than "the largest product
that completes within its budget". 0003:A15's discharge claim narrows
to what C2 observes — enumerator soundness and cardinality as the sole
completion discriminator; the budget-adequacy reading is closed as
definitional. The disposition is recorded in this RDR and in 0003's
landing artifacts; 0003's body is never amended.
```

#### Load-Bearing Decisions

- **Naming** — the observation is "structural exhaustion", and the seam
  is the "ungated proof-completion observation seam"; rejected names:
  "completion oracle" alone (ambiguous with the gate it replaces as
  observer) and "test bound" (a limit parameter is not a second
  published bound).
- **Selection / predicate** — "completes" means: the ungated
  enumeration terminates having yielded exactly
  `Cardinality(m, g)` distinct assignments before the caller's
  arithmetic-derived cap trips; chosen over wall-clock return (excluded
  by REQ-86/ADV-2 and kata `9yeq`) and over `Projectable()` (which
  reads the gate under test back as the observation — the defect this
  RDR removes).

#### Illustrative Code

Shape only — not load-bearing; signatures sharpen at Resolve/Pre-Lock.

```go
// production (unchanged verdicts): callers pass the published bound
values, ok := valueAssignments(decl, Bound())

// export_test.go — compiled only into the guard test binary
var ProductWithin = productWithin // (m, g, cap) -> (AssignmentSet, exhausted bool)

// REQ-133 test, completion half (structural differential)
card, _ := guard.Cardinality(m, g)
set, exhausted := guard.ProductWithin(m, g, card+1)
// assert exhausted && set.Len() == card, on both sides of Bound()
```

### Existing Infrastructure Audit

Sibling-path check: the test-scoped-export decision already exists in
an adjacent path — `internal/table/export_test.go` — so C1 reuses that
convention rather than inventing one; and no third limit signal is
invented — the parameter threads through the two `Bound()`-comparison
sites that already exist (`internal/guard/product.go::Product`,
`internal/guard/assignment.go::valueAssignments`).

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Ungated enumeration | `internal/guard/product.go::cross` | none — already ungated; gate sits in callers | Reuse | none |
| Bound-independent count | `internal/guard/product.go::Cardinality`, `internal/guard/declaration.go::AssignmentCount` | shares `domainSize` with the enumerator's ranges (partial common-mode; see Failure Modes) | Reuse | cap derivation in C2 |
| Gated production entries | `::Product`, `::valueAssignments`, `::Denotation` | limit is hard-wired to `Bound()` today | Extend (limit parameter; production passes `Bound()`) | C1 |
| Completion test harness | `guard_mvv_0003_test.go::TestReq133_TheBoundsVerdictAgreesWithWhetherTheProofCompletes`, `completesWithin` | observes `Projectable()`, which reads the gate back | Replace assertion shape (C2); wall-clock demoted to hang protection | C2 |

### Decision Rationale

Key factors: (1) the record's own writer/consequence chain — the
observation surfaces (`Product`, `Denotation`) consult `Bound()` before
enumerating (`internal/guard/product.go::Product`'s
`card > Bound()` arm; `valueAssignments`' `n > Bound()` precheck)
⇒ any oracle routed through them echoes the gate, which is the residual
`s1r1` left; (2) `cross` and `Cardinality` already provide an ungated
enumerator and a bound-free count ⇒ the differential needs a seam, not
new machinery; (3) REQ-133 is a Testing Strategy obligation, not a
runtime contract ⇒ a test-scoped seam suffices and the production
surface stays fixed; (4) under the REQ-86/ADV-2 exclusions (no
wall-clock, memory, or subprocess-cap oracle) budget-adequacy is
unobservable ⇒ the only honest dispositions are "observe the structural
part" and "close the resource part as definitional", and the chosen
approach does both. Alternative 3 (keep fused) was the serious rival:
respectable on prior art (a bounded verifier's guarantee is scoped, not
predicted — TraceFix p9), and its definitional insight is adopted as
C3; it was rejected as the *whole* answer because it leaves the
over-bound side of the shape pair and the refusal-path enumerator
observable only as the gate's own echo. Alternatives 1 and 2 fail on
surface growth and on REQ-86's non-configurability posture
respectively.

Premortem: hardened (paragraph) — one Failure Modes line folded
(representation swap reopens C2); see
`evidence/propose-premortem/premortem.md`.

Ground-sweep: clean (19 anchors) — 18 confirmed verbatim at source; 1
cosmetic miss fixed inline (test-scoped-export precedent DOES exist:
`internal/table/export_test.go` — it supports the choice; A1 and the
Existing Infrastructure Audit now cite it).

Joint-check: clear (12 peers) — context: 0021/0022 cite
`guard.AssignmentCount` and 0021 cites `product.go::Groups` as read
dependencies on the closed 0003 surface (no peer modifies them, and
this RDR does not either); 0012 modifies `grammar.go::Evaluator` /
`valueSatisfies`, disjoint from this seam's entries.

Joint-check: fired → 0030 (home: `cli/0013 §Normative Contracts` C1).
Context: 0030 (Draft) cites `internal/guard/declaration.go::AssignmentCount`
and rests its A7 on the enumeration limit being this record's
caller-supplied lint scope, not a load-time bound; neither record modifies
the symbol. Recorded symmetrically by 0030's propose, 2026-09-19.

## Alternatives Considered

### Alternative 1: Exported production observation API

**Description**: Export an ungated enumerator from `internal/guard`
(e.g. a public `EnumerateProduct(m, g, cap)`), making the observation
surface part of the conforming package API.

**Pros**:

- No test-scoped shim; any future caller (a diagnostic verb, a
  benchmark) can observe enumeration directly.
- The SPEC-UNDER gap is closed by naming the surface in a spec rather
  than in test scaffolding.

**Cons**:

- Grows the conforming surface for a spec-integrity-only need (the
  record classifies this risk as no user-facing runtime defect).
- A production-reachable ungated enumerator is an invitation to run
  exactly the enumeration ADV-2 exists to forbid on the refusal path.
- REQ-133 is a Testing Strategy obligation; nothing at runtime consumes
  the observation.

**Reason for rejection**: the observation has exactly one consumer (the
REQ-133 test), and a public ungated entry re-creates the ADV-2 hazard
as API.

### Alternative 2: Injectable enumeration gate

**Description**: Make the gate swappable — a package-level limit
variable, a functional option, or a config knob — and have tests inject
a structural stop in place of `Bound()`.

**Pros**:

- No new entry points; the existing surfaces become observable in
  place.

**Cons**:

- Collides with the published-constant posture REQ-86 anchors and the
  `internal/guard/product.go::bound` doc states ("not per-model,
  per-group, or configurable per run"): a mutable gate is a run-time
  knob even when only tests are meant to touch it.
- A swapped gate observes the system in a configuration production
  never runs; the parameter form keeps the production path byte-for-
  byte the code under test.

**Reason for rejection**: it converts the single published constant
into ambient mutable state to solve a problem a call-site parameter
solves without one.

### Alternative 3: Keep the roles fused (no-mechanism disposition)

**Description**: Add no seam. Record that `Bound()` is definitionally
the budget (0003 deviation D6), downgrade 0003:A15 to a definitional
design constraint (Method off `MVV Test`), and restate REQ-133 as a
consistency check confined to what needs no new surface: shape-pair
verdict agreement plus `Product(m, g).Len() == Cardinality(m, g)` in
the admitted region.

**Pros**:

- Zero code surface; the smallest honest response to an unobservable
  claim.
- Prior-art aligned: a bounded verifier's guarantee is scoped by its
  bound — "outside those bounds the guarantee is incomplete"
  (TraceFix p9) — so treating the bound as definition, not prediction,
  is the peer stance.
- The admitted-region soundness check is already bound-independent
  (`Cardinality` consults no bound).

**Cons**:

- The over-bound side of REQ-133 stays vacuous: `Product` returns the
  unprojectable set *because of the gate*, so "does not complete" over
  the bound is the gate's echo — the exact residual this RDR was opened
  on.
- Shape-independence — the one refutable claim A15 makes for an
  enumerating representation — goes untested precisely where shapes
  diverge most (past the bound).
- The refusal-path enumerator (`valueAssignments` on a powerset
  dimension) is never observed at all.

**Reason for rejection**: right about the bound's *value*, wrong about
the enumerator — its definitional insight is adopted (C3) while its
refusal to observe the observable is not.

### Briefly Rejected

- **Dual independent proof representation** (a second bitset/symbolic
  engine as differential N-version oracle): over-scoped by an order of
  magnitude for a Low-priority spec-integrity risk.
- **Wall-clock / memory / subprocess-cap oracles**: already adjudicated
  out — REQ-86 forbids a bound discovered by exhausting a resource,
  ADV-2's comment records why a clock cannot tell a bounded
  implementation from fast hardware, and kata `9yeq` shipped the
  removal of the wall-clock ratio.

## Context

### Background

Raised by roborev job 6143 during kata `s1r1`'s refine loop and routed
out as out-of-scope (a test-only kata cannot add a production export
seam); tracked as kata `intrastate#dh40`. `s1r1` already removed the
outright tautology — REQ-133 formerly compared `card <= Bound()`
against `ProofCompletes`, which *is* `card <= Bound()`; it now observes
completion via `Product` and `Denotation` under a deadline and does
catch a `Product`-side A15 divergence. This RDR closes the residual:
the oracle's remaining bound-gated circularity. The split arm is a new
public surface RDR 0003 does not name — its own Phase-2 rule classifies
that as SPEC-UNDER requiring author decision. Spec-integrity risk only;
no user-facing runtime defect. A decision here constrains but does not
decide kata `ge67`'s finding-set question.

### Technical Environment

Go module `github.com/cwensel/intrastate`. Surfaces:
`internal/guard/product.go::Product`, `::Bound`, `::ProofCompletes`,
`::Denotation`, `internal/guard/assignment.go::valueAssignments`.
Governing record: RDR 0003
(A15, REQ-86, REQ-133, ADV-2, deviation D6).

## Research Findings

### Investigation

Source read on `main`: `internal/guard/product.go` (`Product`'s
`card > Bound()` decline before enumeration; `ProofCompletes` defined
as `card <= Bound()`; `bound = 2048` with its "not configurable per
run" doc; `Cardinality` consulting only `cardinalityCeiling`, never
`Bound()`; `cross` ungated), `internal/guard/assignment.go`
(`valueAssignments`' `n > Bound()` precheck — the second gate site),
`internal/guard/declaration.go` (`AssignmentCount` wrapping
`domainSize` — the arithmetic the cap derives from),
`guard_mvv_0003_test.go` (REQ-133's current observation:
`Projectable()` under a 10s deadline via `productCompletes` /
`denotationCompletes`), `guard_adversarial_0003_test.go` (ADV-2's
structural oracle and its recorded wall-clock rejection). Governing
record: `0003:A15` (Pending, Method `MVV Test`), REQ-86 (`0003:C24`),
0003 `artifacts/deviations.md` §D6. Prior-art pass and query ledger:
`evidence/research/propose-prior-art.md`. One observed tension for
Resolve (carried by A4): REQ-86's "same verdict on every conforming
implementation" sentence vs the `bound` doc's "another conforming
implementation MAY publish a different bound".
⚠ no prior-art coverage for the instance question (peer state-machine
tools' too-large disposition) in the resolving corpora; the
peer-practice claim is demoted to A5 rather than leaned on.

### Key Discoveries

- **Documented** — the circularity is total on the over-bound side:
  both observation routes (`Product`, and `Denotation` via
  `valueAssignments`) consult `Bound()` before enumerating, so the
  current REQ-133 assertion compares the gate to itself there; the
  seed's 2048→64 demonstration follows statically (prediction and
  observation move in lockstep).
- **Documented** — `Cardinality` is already bound-free (it saturates at
  `cardinalityCeiling`, a distinct constant), and `cross` is already
  ungated: the differential's two halves exist; only the seam between
  them is missing.
- **Documented** — the fusion is deliberate and recorded twice:
  `ProofCompletes`' doc ("stating them as one quantity … the claim")
  and deviation D6 ("the bound doubles as the proof representation's
  budget"), so the split must not silently un-record that choice — C3
  restates it as definition instead.
- **Documented** — prior art supports both arms' halves: proof cost may
  not track state-space size for non-enumerating representations
  (TraceFix p7), and a bounded guarantee is scoped rather than
  predicted (TraceFix p9); an oracle knows its answer from outside the
  computation under test (developer-testing p171).
- **Assumed** — test-feasible exhaustion cost over-bound (A2) and the
  export_test reach (A1) need spikes; the peer-practice stance (A5)
  needs a prior-art pass with actual coverage.

## Trade-offs

### Consequences

- Positive: 0003:A15 becomes dischargeable by evidence that could
  refute it — an enumerator that skips, duplicates, or diverges on a
  shape now fails a test instead of shipping as verified.
- Positive: production behavior, verdicts, and the ADV-2 refusal
  posture are byte-for-byte unchanged; the seam exists only in the test
  binary.
- Negative: over-conservatism of the bound's value is *permanently*
  unobservable by test — that is now recorded policy (C3), and a future
  reader who wants the 2048→64 case to fail must reopen C3, not add a
  clock.
- Negative: two entries grow a limit parameter, and the REQ-133 test's
  meaning shifts — a reader of 0003 alone will find the test asserting
  something its Testing Strategy text does not say verbatim (the
  disposition lives here and in 0003's landing artifacts).

### Risks and Mitigations

- **Risk**: the differential is partially common-mode — the cap and the
  enumerator's ranges both derive from `domainSize`, so an arithmetic
  bug could fool both sides at once.
  **Mitigation**: the enumeration count is a real loop-iteration fact
  over a distinct code path (`subsets`/`spreadValues`/`cross` vs
  multiplication in `Cardinality`), and the set's dedup makes `Len()`
  count *distinct* rendered assignments — the representative regression
  classes (gate reordering, materialize-before-check, skipped or
  collapsed points) diverge the two; the residue is named in A4's orbit
  and accepted.
- **Risk**: the test-scoped export pattern spreads and becomes a
  backdoor convention for reaching internals.
  **Mitigation**: C1 confines it to the two named entries; anything
  further is a deviation under the launch prompt's Phase-2 rule.

### Failure Modes

- Visible: the C2 differential fails — the message carries
  (exhausted, count, `Cardinality`, side of the bound), so the
  developer reads which half diverged; recovery is fixing the
  enumerator or the arithmetic, never widening the cap.
- Visible: a production call site passing a non-`Bound()` limit fails
  A3's call-site sweep at Resolve/review.
- Silent (accepted, recorded): a common-mode `domainSize` bug that
  shifts arithmetic and enumeration identically; and any genuine
  resource inadequacy of 2048 on some machine — both are outside the
  REQ-86/ADV-2-admissible observation space, and C3 records the second
  as definitional. Diagnosis path for suspicion: run
  `BenchmarkAdv2RefusalCost` by hand (its comment already assigns it
  the cost-signal duty).
- Spec event, not a test event: swapping the proof representation
  (e.g. to a bitset/symbolic engine) reopens C2 — the observation is a
  claim about the *enumerating* representation, and a new
  representation publishes its own bound under REQ-86; the C2 test must
  not be read as covering a representation it never observes.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified (A1/A2 spikes, A3 sweep, A4
      peer-RDR read, A5 prior-art pass)

### Minimum Viable Validation

1. Build the four `shapePairSource` fixtures (wide/narrow ×
   over/under, per D6 at `2*Bound()` and `Bound()/2`) and the
   `adv2BigSetSource(18)` fixture.
2. For each shape-pair fixture, call the test-scoped
   `ProductWithin(m, g, Cardinality(m, g)+1)`.
3. Assert exhausted with `Len() == Cardinality(m, g)` for all four —
   including the two the bound refuses — and identical
   (exhausted, count) within each equal-cardinality pair.
4. For the set fixture, call the test-scoped ungated
   `valueAssignments` variant under an arithmetic-derived cap; assert
   it exhausts to exactly `domainSize` values while production
   `Denotation` still returns the unprojectable set (ADV-2 intact).
5. Mutation checks: remove `Product`'s bound gate → ADV-2 still fails;
   make `cross` skip one point → step 3 fails; double one point →
   the dedup'd `Len()` still catches it via the exhaustion cap count.

End-state: `go test ./internal/guard` green; REQ-133's completion half
no longer routes through `Bound()`; `go build ./...` contains no
ungated entry.

### Phase 1: Limit parameterization

Thread an explicit limit parameter through the two gate sites;
every production call site passes `Bound()`; verdicts byte-identical.

### Phase 2: Test-scoped observation seam

Add `internal/guard/export_test.go` exposing the ungated
product-assembly and value-enumeration entries (with structural caps)
to `package guard_test`.

### Phase 3: Oracle restatement

Rewrite `TestReq133`'s completion half to the C2 structural
differential; demote the wall-clock deadline to hang protection or
delete it; restate the `ProofCompletes`/`bound` doc comments to the C3
definitional wording.

### Phase 4: Disposition recording

Record the A15 disposition (narrowed discharge claim, definitional
budget) in 0003's landing artifacts per its deviations pattern; 0003's
body is not amended.

## Validation

### Testing Strategy

[Required — never omit. Test scenarios and coverage goals — what to test and
what constitutes "done." For non-functional concerns
(performance, security): state measurement strategy,
not estimates.]

1. **Scenario**: [Description]
   **Expected**: [Result]

### Performance Expectations

[Conditional — omit (don't N/A-bullet) this section unless
comparing alternatives on empirical performance grounds.
Do not include effort estimates or speculative
throughput targets. Rough performance metrics are
appropriate only when comparing alternatives — note
empirical data or obvious gains that support the
chosen approach over a rejected one.]

## Finalization Gate

> Complete each item with a written response in
> `{ARTIFACT_DIR}/gate.md` before marking this RDR as
> **Final**. Written responses prevent rubber-stamping
> and produce a review record.
>
> First run the mechanical pre-sweep
> (`prompts/gate/tooling-pass.md`): TEMPLATE section
> coverage, Method-label vocabulary, `Source Search`
> self-reference, `Docs Only` on load-bearing claims. It
> catches what the review rounds disturbed; resolve any
> BLOCK before the written responses.
>
> At lock, replace Contradiction Check, Assumption
> Verification, Scope Verification and Proportionality
> with the one-line pointer to gate.md — those four
> judge THIS record at THIS lock and no peer cites
> them. **Cross-Cutting Concerns stays here**, below
> the pointer: it names the project-wide policy other
> RDRs conform to, so it must stay projected and
> citable as `cli/NNNN:G-cross-cutting`. Cite it that
> way, not by section name.

### Contradiction Check

[Gate key: contradiction — a gate response is cited as
`cli/NNNN:G-<key>`, so the key is a stable id and is
not derived from this heading, which may be reworded.]

[State any conflicts between Research Findings and
the Proposed Solution. If none exist, state
"No contradictions found between research findings,
design principles, and proposed solution."]

### Assumption Verification

[Gate key: assumptions]

[Confirm every Critical Assumption Evidence Record
is internally consistent: Status, Method, and
Evidence agree, and "If wrong" is non-empty. List
any record whose Method is `Docs Only` (these block
lock unless paired with a Spike or Source Search
plan) and any that remain `Pending` or `Unverified`
with a plan to verify before implementation begins.
Confirm no `Verified` stamp is self-referential or
proves only an adjacent claim, and that each cited
`path::Symbol` resolves on `main`. **Status
consistency:** no assumption marked `Pending` or
`Unverified` may have settled-fact prose elsewhere in
the RDR depending on it.]

### Scope Verification

[Gate key: scope]

[Confirm the Minimum Viable Validation is in scope
and will be executed during implementation, not
deferred. State the specific test or proof.]

### Cross-Cutting Concerns

[Gate key: cross-cutting]

[Retained at lock — this sub-section stays in the RDR
when the other gate responses move to gate.md, because
peer RDRs cite it as `cli/NNNN:G-cross-cutting` and an
element that is not projected cannot be cited.]

[List only concerns that apply to this RDR. For each,
state either how this RDR addresses it, or which peer
RDR owns the project-wide policy this RDR conforms
to. Omit (rather than N/A-bullet) anything that does
not apply.]

Candidate concerns (include only those that apply):
versioning · build tool compatibility · licensing ·
deployment model · IDE compatibility · incremental
adoption · secret/credential lifecycle · memory
management · concurrency model · character encoding ·
canonical-form / determinism (see note below).

If this RDR claims byte-identical output,
content-addressed identity, or replay-stable hashes,
also confirm: hash function + library, pre-image
byte layout, primitive encodings, map iteration order,
whitespace policy, case folding, empty/null/absent
distinguishability, and a version marker for future
evolution.

### Proportionality

[Gate key: proportionality]

[Is the document right-sized for the change? Flag
any sections that should be trimmed before locking.
The split test is **contract count, not word count**:
confirm this RDR is the sole author of at most one
independent load-bearing contract (per the Normative
Contracts split signal). If it owns more than one
seam, flag it for splitting rather than locking the
seams together.

Re-validate the **Profile** Metadata field against the
contracts you just counted: confirm the value Resolve
wrote still matches (one contract + no user-facing
surface → `small`; etc. per the applicability matrix).
If the lenses that actually ran disagree with the
Profile (e.g. Profile says `small` but the change locks
a contract that warranted `mid`+ lenses, or the lenses
were skipped on a wrong `small`), correct the field and
do not lock until the missing lenses have run. This is
the latch's backstop — a wrong Profile cannot route
past the lens battery undetected. A `Transient`-marked
contract with a named deleting sibling and schedule is a
recorded lifespan disposition, not an under-sized
Profile — do not count it when re-deriving. Also confirm form:
value + one clause naming the contract(s); strip any
matrix/provenance prose left from the template or Seed
(it belongs in the template comment, not the instance).]

## References

- RDR 0003 — `0003:A15`, `0003:C24` (REQ-86); artifacts:
  `docs/rdr/0003-guard-predicate-exhaustiveness/artifacts/req-list.md`
  (REQ-133), `artifacts/deviations.md` §D6
- Source reviewed: `internal/guard/product.go`,
  `internal/guard/assignment.go`, `internal/guard/declaration.go`,
  `internal/guard/guard_mvv_0003_test.go`,
  `internal/guard/guard_adversarial_0003_test.go`
- Prior-art query ledger and citations:
  `docs/rdr/0013-ungated-proof-completion-observation-seam/evidence/research/propose-prior-art.md`
- Related issues: kata `intrastate#dh40` (this RDR's tracker); katas
  `s1r1`, `9yeq`, `ge67`; roborev job 6143
