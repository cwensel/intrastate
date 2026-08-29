# Recommendation 0018: Dual-breach precondition precedence and error carrier

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
- **Profile**: foundational — provisional: one contract (the
  dual-breach error's identity and classifiability) closing
  JDR 0001 §JD-5 across the RDR 0008 × RDR 0009 seam.
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
- **Priority**: Medium
- **Related Issues**: kata `intrastate#6cek` (1557); roborev job
  6281; the open joint decision JD-5 in JDR 0001
- **Predecessors**: 0008-recognized-tag-key-ownership,
  0009-escape-row-shape-conformance-ownership
- **Seam Lineage**: no prior accretion

## Problem Statement

A programmatic producer handing `Resolve` an `Input` that breaches
both entry preconditions — RDR 0009's escape-shape breach
(`Table.CheckValid()`) and RDR 0008's reserved-key breach
(`CheckInput`) — needs to classify the single `error` it gets back.
Today that value's identity and classifiability are unowned: no test
pins which precondition reports, the current order at HEAD
(`internal/resolve/resolve.go`, 0009 before 0008) is documented as a
cheapness preference and explicitly "not an observable contract"
(0009 REQ-22; its req-list ASSUMPTION forbids asserting on placement),
and reordering the checks would silently change the returned error
category. JDR 0001 §JD-5 holds the decision open — "Either order is
defensible — pick one and pin it... The silence is the defect, not the
choice" — with a sharpening adding the second half: 0008 has no
sentinel, while 0009's is `errors.Is`-classifiable.

The decision is one fork, not two. First, which precondition is the
reporting one — or whether both report through one `errors.Join`
aggregate so precedence stops mattering. Second, given that choice,
the error carrier the losing/joined side presents: does 0008's breach
gain an exported package-level sentinel plus typed error mirroring
0009's `ErrEscapeShapeBreach`/`*EscapeShapeBreachError`, or does it
stay deliberately unclassifiable so that classifiability itself is the
observable marker of which precondition fired? The halves are
inseparable: pinning an order with no sentinel on the winner is
untestable beyond message text, and giving both sides sentinels makes
the order either free (join) or newly load-bearing (fixed). The
dual-breach pinning test lands as part of this RDR.

## Critical Assumptions

- **A1 At HEAD, a dual-breach `Input` observably returns 0009's
  escape-shape aggregate — pinning `CheckValid`-before-`CheckInput`
  is a pin, not a behavior change.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: pending — an emission/order claim is never
    quote-confirmed; the MVV's dual-breach scenario is the proof,
    and it MUST first run against pre-change HEAD so the "pin, not a
    change" half is itself observed, not inferred (0009 Phase 3c
    probed a hand-built kernel table and observed 0009 reporting
    first, but deliberately did not pin it).
  - **If wrong**: the precedence contract C1 demands a reorder of
    `Resolve`'s entry checks — a caller-visible category change, not
    a pin — and A3's no-behavior-change framing falls with it.
- **A2 Wrapping the three unexported reserved-key channel errors
  with an exported sentinel breaks no existing test or caller: no
  assertion depends on the current errors being unwrapped plain
  values or on their exact message text.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: pending — sweep `internal/resolve/*_test.go` and
    all `CheckInput` callers for message-text or error-equality
    assertions on the reserved-key breach (the 0008 suites are
    external-package tests, so they cannot reference the unexported
    vars directly, but message substrings remain possible).
  - **If wrong**: Phase 1 additionally rewrites the affected
    assertions to `errors.Is` form — scope grows by those tests, and
    the "additive only" consequence needs restating.
- **A3 Pinning 0009-first contradicts no Final fence: `0008:C4`
  licenses either error on a doubly-breaching input and explicitly
  does not foreclose a fixed order, and 0009's precedence fences
  order the shape breach only against modeled dispositions, which
  0008's Go-error breach is not.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: pending — re-verify against `0008:C4` ("JD-5 may
    narrow this to a fixed order; nothing here forecloses that"),
    `0009:C3` ("precedes every modeled disposition"), and `0009:C4`
    (carrier fences unchanged by this RDR).
  - **If wrong**: the choice needs a 07.1 SPEC-DEFECT route against
    the contradicted Final fence instead of landing here — the
    Proposed Solution reopens.
- **A4 No consumer outside `internal/cli` branches on the kernel's
  error return today, and `internal/cli` branches only via
  `errors.Is(err, resolve.ErrEscapeShapeBreach)` — so adding
  classifiability to the reserved-key breach changes no shipped
  branch.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: pending — sweep non-test callers of
    `resolve.Resolve` / `resolve.CheckInput` for error inspection
    beyond the `internal/cli/flow_resolve.go::kernelResolveFailure`
    discrimination.
  - **If wrong**: an existing negative-classification branch
    ("unclassifiable ⇒ reserved key") silently inverts when the
    sentinel lands; that caller must be rewritten in Phase 1.

## Proposed Solution

### Approach

Answer both halves of the fork with the smallest surface that makes
each half structurally testable. **Precedence**: pin the order the
code already has — on an `Input` breaching both preconditions,
`Resolve` returns 0009's escape-shape aggregate; `Table.CheckValid()`
reporting before `CheckInput` becomes an observable contract owned
here (it was "a cheapness preference, not an observable contract"
under 0009 REQ-22 — this RDR supersedes that framing by citation, not
amendment). **Carrier**: the reserved-key breach becomes
`errors.Is`-classifiable by one exported package-level sentinel,
`ErrReservedTagKey`, which the three existing unexported channel
errors wrap. No typed error and no structured payload: the reserved
key is a compile-time constant and `0008:C4` fixes first-breach
single-error reporting, so there is no payload a caller needs — the
0009-style typed mirror is not reproduced. The dual-breach pinning
test then asserts both halves structurally — `errors.Is(err,
ErrEscapeShapeBreach)` true and `errors.Is(err, ErrReservedTagKey)`
false — with no message-text assertion, and JDR 0001 §JD-5 closes at
its home citing this RDR as the implementation carrier.

### Technical Design

Two surfaces move, both inside `internal/resolve`:

- `internal/resolve/precondition.go` — gains the exported sentinel;
  the three channel errors (`errReservedOwnedTag`,
  `errReservedObservedTag`, `errReservedRequiresOwned`) become
  wrappers that unwrap to it. Channel identity stays diagnostic
  prose, exactly as today.
- `internal/resolve/resolve.go::Resolve` — the call order
  (`in.Table.CheckValid()` then `CheckInput(in)`) is unchanged; its
  comment block stops disclaiming observability and cites this RDR's
  precedence contract instead.

The CLI mapping (`internal/cli/flow_resolve.go::kernelResolveFailure`)
is deliberately untouched: a reserved-key breach keeps riding the
generic internal-error branch. Whether it gains a dedicated CLI code
is a code-table question owned by JDR 0001 §JD-8/§D10 and RDR 0005's
re-entry, not by this RDR — the sentinel makes that later
discrimination possible without deciding it here.

#### Normative Contracts

**C1**

```normative
Dual-breach precedence. When an Input breaches both entry
preconditions — RDR 0009's escape-row shape rule
(Table.CheckValid()) and RDR 0008's reserved-key rule
(CheckInput) — Resolve MUST return the escape-shape breach
error: errors.Is(err, ErrEscapeShapeBreach) reports true and
errors.Is(err, ErrReservedTagKey) reports false, with the zero
Result. This makes the CheckValid-before-CheckInput evaluation
order at Resolve entry an observable contract (superseding, by
citation, 0009 REQ-22's "cheapness preference, not an observable
contract" framing); reordering the entry checks is a breaking
change to this contract, not an implementation detail. The
returned error is 0009's aggregate exactly as `0009:C4` fences
it — verbatim, flat, every element a *EscapeShapeBreachError —
and this RDR adds nothing to it. A dual-breach pinning test MUST
assert both errors.Is verdicts above; it MUST NOT assert on
message text.
```

**C2**

```normative
Reserved-key breach carrier. The kernel MUST export a
package-level sentinel, ErrReservedTagKey, and every error
CheckInput returns MUST report errors.Is(err, ErrReservedTagKey)
== true — all three breach channels (owned tag, observed tag,
Row.RequiresOwned) classify to the one sentinel. The breach
stays a single non-aggregate error with first-breach-wins
reporting and implementation-latitude scan order, exactly as
`0008:C4` fences it; this RDR adds classifiability only. No
typed error and no structured payload accompany the sentinel:
channel identity and the offending row remain diagnostic prose,
and a test or caller needing them MUST NOT parse message text —
that need would be a design change routed through a successor
RDR, not latitude here.
```

#### Load-Bearing Decisions

- **Identity** — classification is per-category, not per-channel:
  one sentinel names "reserved-key breach"; the three channels stay
  distinguishable only in prose. Callers branch on the category
  (mirrors the authored path, where one `reserved_tag_key` load
  category covers every channel).
- **Naming** — `ErrReservedTagKey`, aligning the Go sentinel with
  the existing category literal `reserved_tag_key`
  (`internal/table/category.go::CatReservedTagKey`) and the
  kernel-aligned code-name rule of JDR 0001 §D10. Rejected:
  `ErrReservedKey` (drops the tag vocabulary), `ErrReservedInput`
  (names the argument, not the rule).
- **Selection / predicate** — when both preconditions qualify to
  report, the escape-shape breach wins. Why this winner: a shape
  breach is the table malformed *as a value*, independent of the
  input tuple (`0009:C3`'s own rationale — the table does travel
  inside `Input`, but `CheckValid` reads only `in.Table` while
  `CheckInput` reads the call's tag sequences, so the asymmetry is
  scope-of-read, not a phantom partition); it matches the order at
  HEAD (a pin, not a change — A1); and it is the order under which a
  dual-breach caller keeps the richer report (0009's per-row
  structured aggregate vs 0008's single prose line).

#### Illustrative Code

```go
// Illustrative — shape only.
// precondition.go
var ErrReservedTagKey = errors.New(
    "resolve: input supplies the reserved tag key " + recognizedTagKey)

var errReservedOwnedTag = fmt.Errorf(
    "%w as an owned tag; the recognized outcome enters only "+
        "through Input.Recognized", ErrReservedTagKey)

// dual-breach pinning test (external package)
_, err := resolve.Resolve(dualBreachInput)
// errors.Is(err, resolve.ErrEscapeShapeBreach) == true
// errors.Is(err, resolve.ErrReservedTagKey) == false
```

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Reserved-key breach detection | `internal/resolve/precondition.go::CheckInput` + three unexported channel errors | Unexported ⇒ unclassifiable outside the package | Extend (wrap with exported sentinel) | C2 |
| Entry-order evaluation | `internal/resolve/resolve.go::Resolve` (CheckValid at 0, CheckInput at 0.5) | Order documented as non-observable | Reuse (pin as observable) | C1 |
| Escape-shape classification | `resolve.ErrEscapeShapeBreach` / `*EscapeShapeBreachError` (`0009:C4`) | None | Reuse unchanged | C1 cites, adds nothing |
| CLI kernel-error mapping | `internal/cli/flow_resolve.go::kernelResolveFailure` | Non-0009 errors collapse to generic `flow-accessor-failed` | Reuse unchanged (code-table change rides JDR 0001 §D10 / RDR 0005 re-entry) | none here |

### Decision Rationale

The QOC matrix (Alternatives Considered carries the losing analyses)
decided on three rows: **fence compatibility** eliminated the join
(`0009:C4`'s verbatim/flat-aggregate fence cannot host a
cross-precondition join without reopening a Final record);
**structural testability of the losing side** eliminated the
no-sentinel option (a marker-by-absence degrades to "not 0009" the
moment a third kernel error class exists, and its test is
message-text-only — the untestability JD-5's sharpening names); and
**behavior preservation + report richness** picked 0009-first over
0008-first (a pin of HEAD rather than a silent category change, and
the dual-breach caller keeps per-row structured diagnostics). The
carrier stops at a sentinel because the guidance fork for Go error
surfaces — caller needs matching + static message ⇒ top-level `var`
(uber-go-guide §Error Types) — lands there, and 0009's heavier typed
surface was justified by a structural payload (RowRef findings) that
0008's first-breach-single-error fence deliberately does not owe.
Sibling-path check: the "which precondition fired" discriminator
already exists once at
`internal/cli/flow_resolve.go::kernelResolveFailure`
(`errors.Is(err, resolve.ErrEscapeShapeBreach)`) — this RDR extends
that existing signal family rather than inventing a parallel one, and
the sentinel's name reuses the authored path's existing category
identity `internal/table/category.go::CatReservedTagKey`
(`"reserved_tag_key"`) rather than minting new vocabulary.

Premortem: hardened (hardened)
Ground-sweep: clean (14 anchors)
Joint-check: clear (12 peers)

Premortem mitigations folded (critic ledger P-1..P-12,
`evidence/propose-premortem/critic.md`): the pre-change-HEAD
characterization run (P-1/P-12 → A1, MVV step 5), byte-preserved
messages under a `%w` chain (P-3/P-5 → Phase 1, Risk 1, A2's
message-assertion sweep), the non-aggregate carrier assertion (P-4 →
MVV step 4), white-box per-channel wrap proofs (P-7 → MVV step 4),
the JD-5-closure prerequisite so the pin is never solely
self-declared (P-2 → Prerequisites, Phase 3), the negative-assertion
blind-spot acknowledgment (P-6 → Failure Modes), the scope-of-read
clarification (P-8 → D-Selection), and the honest locator-cost
wording on the R4 rejection (P-10 → Briefly Rejected). P-9's owner
is named in Technical Design (JDR 0001 §D10 / RDR 0005 re-entry);
P-11's demanded quotation already sits verbatim in Alternative 1 and
was CONFIRMED by the ground sweep. Ground-sweep anchor 14
was refuted cosmetically (a dual-breach *construction* exists in
test; the missing thing is the ordering assertion) and corrected
inline — the choice was not reopened. Joint-check context: the open
peers are the 12 Draft roster siblings; the two grep hits
(`CheckInput` in 0012, inside its *rejected* Alternative 3;
`reserved_tag_key` in 0024, as an append-precedent citation) are
context, not decisions at this locus. RDRs 0008/0009 are
`Implemented`, outside the open-peer set; their coupling on this
question is already homed at JDR 0001 §JD-5, which this RDR closes —
cite, don't restate.

## Alternatives Considered

Scored QOC matrix. Question: *on a dual-breach `Input`, which
precondition's error does `Resolve` return, and what carrier does the
reserved-key breach travel?* Options: **O1** pin 0009-first + export
`ErrReservedTagKey` (chosen); **O2** pin 0008-first + same sentinel;
**O3** `errors.Join` both so order stops mattering; **O4** pin
0009-first, no sentinel (classifiability-as-marker). One clause per
cell; ✓ favorable, ✗ disqualifying, ~ neutral.

| Criterion | O1 (chosen) | O2 0008-first | O3 join both | O4 no sentinel |
| --- | --- | --- | --- | --- |
| Correctness fit (both halves structurally testable) | ✓ two `errors.Is` assertions | ✓ same carrier | ✓ both classifiable | ✗ losing side testable only via message text |
| Final-fence compatibility | ✓ `0008:C4` invites the narrowing | ✓ likewise | ✗ breaks `0009:C4` verbatim/flat-aggregate fence | ✓ no fence touched |
| Prior-art alignment | ✓ sentinel per uber-go-guide §Error Types | ✓ same | ~ k8s `field.ErrorList` precedent is user-input validation, not programmer-mistake preconditions | ✗ negative classification has no precedent |
| Behavior preservation | ✓ pins HEAD (A1) | ✗ silently swaps dual-breach category | ✗ changes return shape for dual breach | ✓ pins HEAD |
| Dual-breach report richness | ✓ per-row structured aggregate wins | ✗ single prose line wins | ✓ both present | ✓ aggregate wins |
| Reversibility / blast radius | ✓ additive sentinel + one test | ~ additive + reorder churn | ✗ requires reopening Final 0009 | ✓ smallest diff, but locks in fragility |
| Future third error class | ✓ each class self-names | ✓ same | ~ aggregate grows | ✗ "unclassifiable ⇒ 0008" inverts |

### Alternative 1: `errors.Join` both breaches (O3)

**Description**: `Resolve` evaluates both preconditions
unconditionally and returns `errors.Join(shapeErr, inputErr)` when
both fire, so precedence stops existing as a question and both
breaches classify via `errors.Is`. This is the Kubernetes-validation
shape (`field.ErrorList` appended across checks, returned as one
invalid error — `langref/kubebuilder` cronjob webhook).

**Pros**:

- Dissolves the precedence half entirely; no order is load-bearing.
- A dual-breach producer sees both mistakes in one round trip.

**Cons**:

- `0009:C4` fences the escape-shape return: `Resolve` MUST return
  `CheckValid`'s error VERBATIM, and the aggregate's
  `Unwrap() []error` MUST be exactly one level deep with *every*
  element a `*EscapeShapeBreachError`. Wrapping it in an outer join
  (or admitting a foreign element) breaks both clauses — the CLI's
  flat traversal (`internal/cli/flow_resolve.go::escapeShapeBreaches`)
  depends on them — so O3 cannot land without reopening Final 0009.
- The prior art misfits the channel: aggregation serves user-facing
  declarative validation; this is the programmer-mistake path, where
  `0008:C4` already rules a producer "is not owed an exhaustive
  list".

**Reason for rejection**: disqualified on the fence-compatibility
row — it is the only option whose cost is a Final-record reopen; its
whole benefit (both mistakes in one trip) is worth less than that on
a programmer-mistake channel.

### Alternative 2: Reserved-key-first (O2)

**Description**: same carrier work as O1, but pin `CheckInput` to
report before `Table.CheckValid()` — the dual-breach test asserts the
inverse pair of `errors.Is` verdicts, and the entry checks reorder.

**Pros**:

- Equally pinnable and structurally testable.
- Arguably checks the cheaper predicate first (three sequence scans
  vs a whole-row-set scan — both trivial).

**Cons**:

- Silently changes the de facto dual-breach category at HEAD (A1) —
  the exact silent-reorder hazard the Problem Statement exists to
  close would be exercised once, deliberately, to no benefit.
- The dual-breach caller trades 0009's per-row structured aggregate
  for a single prose line — strictly poorer diagnostics.
- No semantic rationale beats `0009:C3`'s: a shape breach is the
  table malformed as a value, independent of the tuple; occupancy of
  one call's `Input` is the narrower fact.

**Reason for rejection**: loses the behavior-preservation and
report-richness rows and wins none.

### Alternative 3: Classifiability as the marker (O4)

**Description**: pin 0009-first but export nothing new — the fact
that a dual-breach error satisfies `errors.Is(err,
ErrEscapeShapeBreach)` while a pure reserved-key breach satisfies no
sentinel *is* the observable marker of which precondition fired.

**Pros**:

- Zero new exported surface; nothing to name or lock.

**Cons**:

- The reserved-key side of the pinning test can only assert "not
  classifiable as 0009" — which is also true of any future kernel
  error — or fall back to message text, the exact untestability
  JD-5's sharpening flags ("0008 has no sentinel").
- Negative classification is already under pressure at HEAD:
  `kernelResolveFailure`'s fall-through comment ("A blanket recode
  would mislabel RDR 0008's reserved-key breach") shows a consumer
  reasoning about 0008 by absence today.

**Reason for rejection**: answers the precedence half while leaving
the carrier half open — the JD-5 sharpening explicitly requires
naming the error wrapping, and a marker-by-absence stops marking the
moment a third error class arrives.

### Briefly Rejected

- **Full 0009-style mirror (exported typed error + structured
  payload for the reserved-key breach)**: the key is a compile-time
  constant and `0008:C4` fixes first-breach single-error reporting;
  the breach *site* (which sequence, which row) stays diagnostic
  prose — that locator cost is accepted deliberately (a producer
  with a large generated input diagnoses from the message), and a
  future structured-locator need is a successor-RDR design change
  (C2 says so), not grounds to lock a typed surface with no present
  consumer.
- **Three exported per-channel sentinels**: no caller needs channel
  granularity; it triples the locked surface and invites branching on
  a distinction `0008:C4` keeps as prose.

## Context

### Background

JDR 0001 §JD-5 sits in the JDR's Interface record unclosed (JD-6/
JD-7 around it are closed); RDR 0009 has since moved to
`Implemented` with no Status qualifier, so the JDR entry is the only
live pointer holding the question open. The gap is logged verbatim
in the 0009
implementation's `coverage.md` §Q-E ("JD-5's precedence ordering stays
unconstrained — no test decides it") and Phase 3c probed a dual-breach
input against a hand-built kernel table, observing 0009 reporting
first but deliberately not pinning it. Independently re-raised by the
roborev range review (job 6281); tracked as kata `intrastate#6cek`.
Exposure is to programmatic producers only: the authored TOML path
cannot construct a dual breach (`internal/table/normalize.go` rejects
any escape rule bearing a write block). Triage checked and rejected
the collapse hypothesis — 0009's fenced MUSTs (REQ-19/21) order the
shape breach only against *modeled dispositions*, and 0008's
reserved-key breach is explicitly not one, so both halves are
genuinely open and pinning inside 0009 alone would unilaterally close
a decision JDR 0001 reserves to the joint seam. Home: an amendment
closing JDR 0001 §JD-5, with this RDR carrying the implementation.

### Technical Environment

Go module `github.com/cwensel/intrastate`. Surfaces:
`internal/resolve/resolve.go` (`Table.CheckValid()` / `CheckInput`
call order), 0009's exported `ErrEscapeShapeBreach` /
`*EscapeShapeBreachError`, `internal/table/normalize.go` (escape
write-block rejection). Governing records: JDR 0001 §JD-5, RDR 0008
(A11, reserved-key normative block), RDR 0009 (C3, REQ-19, REQ-21,
REQ-22).

## Research Findings

### Investigation

Read before enumerating (Stage 2 pass; queries, rejected branches,
and accepted citations in
`evidence/research/propose-prior-art.md`): the governing fences via
the projector (`0008:C4`, `0009:C3`, `0009:C4`, JDR 0001 §JD-5 and
§D10), the seam at HEAD (`internal/resolve/resolve.go::Resolve`,
`internal/resolve/precondition.go::CheckInput`,
`internal/cli/flow_resolve.go::kernelResolveFailure`), 0009's
implementation artifacts (req-list REQ-19/21/22 and the REQ-22
ASSUMPTION), and a bounded external pass. ⚠ no prior-art coverage for
dual-precondition precedence ordering in the corpora (DevRef and
StateMachineRes queries returned off-class hits only) — the
precedence choice rests on repo-internal anchors; the carrier choice
rests on uber-go-guide §Error Types and the k8s aggregate precedent
recorded above.

### Key Discoveries

- **Documented** — `0008:C4` closes with: "apply this predicate at
  entry and report its breach; if the table also breaches RDR 0009's
  shape rule, **either error is conforming** … JD-5 may narrow this
  to a fixed order; nothing here forecloses that." ⇒ pinning an order
  here contradicts no 0008 fence; the writer of the licensed
  latitude is 0008 itself, and JD-5 is its named narrowing venue.
- **Documented** — `0009:C3` fences precedence only against *modeled
  dispositions* ("it precedes every modeled disposition (a malformed
  table is malformed as a value, independent of the input tuple)";
  REQ-21 carries the same bound via `0009:D-selection-predicate`);
  0008's breach travels the Go-error path and is not one. ⇒ 0009's fences do not decide the dual-breach
  order on their own; this RDR's C1 is a new contract, not a
  restatement (this is triage's anti-collapse finding, kept because
  JD-5's sharpening can be misread the other way).
- **Documented** — `0009:C4` fences the escape-shape return as
  verbatim and flat (one level, all elements
  `*EscapeShapeBreachError`), and
  `internal/cli/flow_resolve.go::escapeShapeBreaches` traverses
  exactly that shape. ⇒ any cross-precondition join is fence-breaking
  (kills O3), and C1 must add nothing to the aggregate.
- **Documented** — `internal/resolve/resolve.go::Resolve` calls
  `in.Table.CheckValid()` then `CheckInput(in)`, with the comment "a
  cheapness preference, not an observable contract"; the three
  reserved-key errors in `internal/resolve/precondition.go` are
  unexported plain `errors.New` values. ⇒ the pin matches HEAD and
  the carrier gap is real.
- **Documented** — uber-go-guide §Error Types: caller needs matching
  + static message ⇒ "top-level `var` with errors.New"; the typed-
  error row is for dynamic/structured needs. ⇒ sentinel-only is the
  convention-aligned carrier; 0009's typed surface was payload-
  driven (`0009:C4` RowRef) and is not precedent for payload-free
  0008.
- **Verified** — a dual-breach test DOES exist at HEAD
  (`internal/resolve/reserved_key_0008_test.go::TestReq44And45And90_ReservedKeyBreachIsNotSkippedByACoincidentShapeBreach`
  builds the combined breach), but it deliberately asserts only that
  *some* breach reports — never which — per `0008:C4`'s
  not-skipped clause; no test asserts the ordering, matching
  coverage.md §Q-E ("no test asserts the relative order", scoped to
  the 0009 suite). ⇒ the MVV's pinning test extends this existing
  fixture pattern rather than inventing one, and the gap is the
  assertion, not the construction.
- **Assumed** — dual-breach emission at HEAD is 0009-first (observed
  by 0009 Phase 3c, but an order claim is never quote-confirmed) —
  A1, Method: MVV Test.

## Trade-offs

### Consequences

- Positive: a programmatic producer can finally branch on the
  reserved-key breach (`errors.Is`) instead of parsing prose, and the
  dual-breach category becomes test-pinned — reordering the entry
  checks now fails a test instead of silently swapping what callers
  see.
- Positive: JDR 0001 §JD-5 gets its closing answer with zero change
  to either Final peer's fenced text — 0008's licensed latitude
  narrows exactly the way its own fence anticipated.
- Negative: the entry-check order is now load-bearing; a future
  precondition added to `Resolve` must state its own order against
  both existing checks at introduction (C1 covers only this pair).
- Negative: `ErrReservedTagKey` is permanently locked exported
  surface on the kernel package; a later structured-payload need is a
  successor-RDR design change, not latitude (C2 says so explicitly).
- Neutral: a dual-breach producer discovers its two mistakes
  sequentially (fix shape breach, then see reserved-key breach) —
  acceptable on the programmer-mistake channel where `0008:C4`
  already owes no exhaustive list.

### Risks and Mitigations

- **Risk**: wrapping the three channel errors changes their rendered
  message text, breaking an unnoticed message-substring assertion or
  log scraper.
  **Mitigation**: A2 sweeps for message/equality assertions before
  implementation; C2 keeps messages diagnostic-only, so any such
  assertion found is rewritten to `errors.Is` form as part of
  Phase 1.
- **Risk**: a consumer already branches by negative classification
  ("not `ErrEscapeShapeBreach` ⇒ reserved key"), which the new
  sentinel silently strengthens or an added third error class later
  weakens.
  **Mitigation**: A4 sweeps all `resolve.Resolve`/`CheckInput`
  callers; the only known discriminator
  (`kernelResolveFailure`) treats non-0009 generically and stays
  conforming.
- **Risk**: JD-5's closure is recorded at the JDR home in a way that
  restates rather than cites, drifting from this RDR's C1.
  **Mitigation**: the closure entry cites `0018:C1`/`0018:C2` —
  cite-don't-restate is the JDR's own interface-record convention.

### Failure Modes

- **Visible**: a dual-breach `Resolve` call returns an error that
  classifies as `ErrReservedTagKey` — the pinning test fails; the
  diagnosis is an entry-check reorder in `Resolve`, and the recovery
  is restoring `CheckValid` before `CheckInput` (or a successor RDR
  re-deciding C1).
- **Visible**: a reserved-key breach error stops classifying
  (`errors.Is` false) — a channel error lost its wrap; the
  single-breach classification tests name the broken channel.
- **Silent (guarded)**: before this RDR, a reorder swapped the
  dual-breach category with no test noise — that is the defect being
  closed; after it, the same edit is loud by construction.
- **Silent (bounded, accepted)**: the dual-breach pin's
  reserved-key half is a *negative* assertion
  (`errors.Is == false`), which cannot by itself distinguish "shape
  check won" from "reserved-key check deleted" — the suite as a
  whole closes that hole (deleting `CheckInput` turns the
  single-breach channel tests red even while the dual-breach test
  stays green), and the MVV's flip demonstration shows the pin is
  sensitive to the reorder it guards.
- **Diagnosis path**: `errors.Is` against the two exported sentinels
  partitions every kernel entry-precondition error; anything
  classifying as neither is by construction some other kernel error,
  and message text remains the human-facing diagnostic only.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified (A1 rides the MVV itself;
      A2/A4 are pre-implementation source sweeps; A3 is a peer-fence
      re-read)
- [ ] The JD-5 closure entry (Phase 3) is drafted alongside the
      code, so no implementation ships calling the order "pinned"
      while the joint decision's home still reads open — `0008:C4`
      routes the narrowing authority there, not here

### Minimum Viable Validation

1. In an external-package kernel test, build a dual-breach `Input`
   programmatically: one escape row carrying a write (`0009` breach)
   and an owned tag keyed `recognized` (`0008` breach).
2. Call `resolve.Resolve`; assert the error is non-nil and the
   `Result` is zero.
3. Assert `errors.Is(err, resolve.ErrEscapeShapeBreach) == true` and
   `errors.Is(err, resolve.ErrReservedTagKey) == false` — the
   precedence pin (C1), with no message-text assertion.
4. Rebuild the `Input` with only the reserved-key breach (each of
   the three channels in turn); assert
   `errors.Is(err, resolve.ErrReservedTagKey) == true` and that the
   error is a single non-aggregate value (no `Unwrap() []error`) —
   the carrier (C2). Add three white-box one-liners asserting
   `errors.Is` directly on each unexported channel error, so "each
   channel wraps" is a checked fact, not a routed guess (a
   channel-routing mistake in the black-box input otherwise fakes
   coverage).
5. End-state: both tests green at HEAD order — with step 3 run once
   against pre-change HEAD before the carrier lands (A1's
   observation); flipping the two entry checks in `Resolve` makes
   step 3 fail — demonstrated once during MVV, then reverted.

### Phase 1: Carrier

Export `ErrReservedTagKey` in `internal/resolve/precondition.go` and
make the three channel errors wrap it via a `%w` chain (C2),
preserving each rendered message byte-for-byte where the wrap
allows — a message that must change is recorded, not slipped;
rewrite any assertion A2 uncovered to `errors.Is` form.

### Phase 2: Precedence pin

Land the dual-breach pinning test (C1) and replace the
`resolve.go::Resolve` comment's "not an observable contract"
disclaimer with a citation of `0018:C1`.

### Phase 3: Close the joint decision at its home

Record JD-5's closure in JDR 0001 (Open → closed, citing
`0018:C1`/`0018:C2` — cite, don't restate); RDR 0009's Status-line
qualifier clearing then follows that record's own rule at its next
gate, outside this RDR.

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

- [Requirements/standards with section numbers]
- [Dependency docs, source paths reviewed]
- [Dependency repos searched (clone + code search)]
- [Related issues, articles, discussions]
