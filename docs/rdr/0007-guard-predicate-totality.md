# Recommendation 0007: Guard predicate totality over an incomplete evaluation view

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-08-09
- **Status**: Draft
  <!--
  - `Demoted` is the terminal status for an RDR judged
    *not RDR-shaped* — the decision was never a real
    design fork, so it leaves the RDR lifecycle and is
    refiled as a plain issue. Carry the destination on the
    live value: `Demoted [→ <issue link>]`, and record the
    same link under **Related Issues**. A `Demoted` RDR runs
    no further stages. (Distinct from the 08.1 *demotion*
    below, which is a `Final → Draft` flip that keeps the
    RDR in the lifecycle — that flip never writes
    `Status: Demoted`; see the disambiguation note there.)
  - A Draft demoted from Final by the 08.1 cluster gate
    carries a qualifier on the live value:
    `Draft [revised from Final YYYY-MM-DD; re-verify A2,A4
    — <one-line reason>]`. It is still a `Draft` for every
    binary Draft/Final gate; only Stage 4 (scoped
    re-verify) and Stage 8 (re-lock) parse the qualifier.
    The Stage 8 flip to `Final` overwrites the whole value,
    so the qualifier self-clears at re-lock — no separate
    cleanup. This 08.1 "demotion" is a *verb* describing the
    Final→Draft flip; it is **not** the `Demoted` status
    above (which exits the lifecycle to an issue) — do not
    conflate the two. (`Reverted` above is the unrelated
    terminal "implementation rolled back" status — also do
    not conflate.)
  -->
- **Type**: Architecture
- **Profile**: foundational — provisional; one contract
  (the domain of guard-predicate evaluation over an
  incomplete view), but it is implemented by RDR 0003's
  evaluator and consumed by RDR 0001's refusal taxonomy, so
  it spans two already-Final RDRs. Resolve overwrites from
  the verified count.
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
  is floored at FOUNDATIONAL regardless of the contract axis.
  Matrix: rdr/stages/README.md. Seed estimates from the design
  shape; Resolve overwrites from the verified count; Stage 8
  Gate locks it at Draft → Final. Never skip lenses off a
  Draft Profile until Resolve has run. -->
- **Priority**: Medium
- **Related Issues**: kata `xg7p` — "RDR 0001:
  RequiresOwned conflates guard-input state with post-guard
  transition state" (`src:roborev`, `severity:medium`,
  `area:internal-resolve`).
- **Predecessors**: 0001-resolution-kernel,
  0003-guard-predicate-exhaustiveness
- **Overrides**: None superseded. This RDR *narrows* RDR
  0001's documented `Row.RequiresOwned` meaning (to
  post-guard write-dependency keys) and *adds* an
  evaluation-domain obligation to RDR 0003's guard
  evaluator; it is the single normative home of the domain
  rule, which both peers cite rather than restate.
- **Seam Lineage**: `area:internal-resolve` (guard-verdict
  semantics at the RDR 0001 ↔ 0003 boundary) — no prior
  accretion.

## Problem Statement

Someone authoring a transition table — and the guard
predicates its rows depend on — needs to know what happens
when a guard asks about state the runtime view does not
carry. They want a refusal they can act on: told plainly
that the artifact state needed to decide was missing, not
handed a plan that quietly routed around it. They discover
they need this the first time a table with a modeled escape
produces a plan in a situation where the honest answer was
"I could not tell." At that point the table looks like it
worked, and nothing in the output distinguishes "this row
does not apply" from "I could not read what I needed to
find out."

Internally, the resolution kernel must fix the **domain of
guard-predicate evaluation**. When a guard atom references
a tag absent from the assembled evaluation view, the seam
today may answer with a two-valued verdict — treating
absence as non-satisfaction — or refuse to answer at all.
The two answers are not interchangeable: a `GuardFalse`
prunes the row before it can raise an owned-state
obligation, and the resulting `no_match` is escapable per
RDR 0002, while `owned_state_unavailable` is not. So the
choice of domain decides whether missing artifact state can
be masked behind an escapable refusal class.

Nothing normative states which answer is required.
`Row.RequiresOwned` is documented as naming "owned tag keys
the row's evaluation needs," which reads as covering both
the state a guard reads to be decidable and the state a
transition's writes depend on once the guard holds. Neither
RDR 0001 nor RDR 0003 declares which of those two roles a
given key plays, so a table author and a guard evaluator can
each hold a defensible but incompatible reading. This RDR
must state, normatively and in one place, whether the guard
predicate is a total function over the view or a partial
function whose domain is the set of tags it references —
and which side of the RDR 0001 / RDR 0003 boundary owns
that rule.

## Context

### Background

Raised by a roborev finding against the RDR 0001
implementation (kata `xg7p`, `src:roborev`), then grounded
by live probes against the current `internal/.../resolve.go`
during seed triage. The probes established the boundary
precisely; they are recorded here as *observed current
behavior*, not as a design position:

- An overlap key present as *observed* with a matching value
  yields `owned_state_unavailable` with `MissingOwned:[gate]`.
- An overlap key absent entirely yields `no_match` — the row
  genuinely does not match, per the `req-list.md`
  ASSUMPTION's own exclusion clause.
- A FALSE guard plus an absent `RequiresOwned` key plus a
  modeled `no_match` escape yields a plan via the escape,
  `Escaped:true`. This is the path where missing artifact
  state is masked.

No input produces an unambiguously wrong disposition, which
is why this is under-specification rather than a defect. The
exposure is an unconstrained seam, and the seam is exactly
where an author-written predicate lands.

Constraints inherited from prior rounds. RDR 0001's
Phase 3c deviation **D8** (guard-FALSE prunes first; a
pruned row contributes neither candidacy nor an owned-state
obligation) is pinned by frozen tests ADV-1 and ADV-1b,
which two independent verifiers proved. D8 supersedes D4's
per-row ordering half and is recorded as an explicitly
requested author decision. Deviation **D5** scopes
`RequiresOwned` and guard evaluation to candidate rows only,
so that an unrelated row's key cannot poison a legal
resolution. RDR 0002 L292 closes escape lists to `no_match`
and `ambiguous_match`, which is the asymmetry that makes
verdict ordering load-bearing in both directions.

Seed triage corrected one premise carried in the original
report: it assumed RDR 0003 was unimplemented. RDR 0003 is
**Final**. Its normative blocks — grammar, typing,
`all`/`unless`, exhaustiveness, row-group scoping,
set-domain proof, diagnostics, provenance lint — contain no
rule for what an atom evaluates to when its tag is absent
from the runtime view; its only absent-state rule (owned
read-before-write) is a lint rejection, not a runtime
verdict. The fork is live.

Sequencing. `triage.md` records that D8 should be ratified
before RDR 0003 is implemented, because 0003's
guard-evaluator contract depends on how the fork resolves.
This makes 0007 the sequencing-critical member of its seed
batch.

Scope note carried from triage: the originating kata also
listed a secondary item — ratifying deviations D8 and D1.
D8 is this fork's downstream consequence and is settled by
whichever branch is chosen. D1 (escape scope) is
independent of guard semantics and already resolved by
evidence rather than open design; it needs an author
signature, not an RDR, and is not in scope here.

### Technical Environment

Go; the resolution kernel under `internal/` (`resolve.go`,
`Row.RequiresOwned`, the guard-verdict enum). The contract
sits at the boundary between two Final RDRs:

- **RDR 0001 — Resolution kernel** (`Implemented`). Owns
  the closed refusal taxonomy and the selection rule: "the
  only successful selection is exactly one matching edge
  after guard evaluation. Zero, multiple, unavailable, or
  unevaluable candidates are refusals unless the table
  contains a modeled escape edge that itself matches exactly
  once." The clause "after guard evaluation" is what
  authorizes pruning. RDR 0001 is the consumer that maps a
  guard verdict onto its five refusal kinds.
- **RDR 0003 — Guard predicate exhaustiveness** (`Final`,
  not yet implemented). Owns predicate semantics; its
  assumption A4 already claims the "unevaluable guard"
  semantic kind.
- **RDR 0002 — Transition table as reviewable data**
  (`Final`). Supplies the escapable-class constraint at
  L292.

Related but distinct seams, deliberately not folded in
here: katas `z53t` (recognized-outcome tag key binding, now
RDR 0008) and `tgzh` (escape-row shape conformance, now RDR
0009) are RDR 0001 ↔ 0002 questions about row shape and key
naming; this is an RDR 0001 ↔ 0003 question about
guard-verdict semantics.

## Research Findings

### Investigation

Prior art was read before enumerating approaches; the full
query ledger and accepted/rejected citations are cached at
`docs/rdr/0007-guard-predicate-totality/evidence/research/propose-prior-art.md`.
⚠ no prior-art coverage in the arc corpora
(`StateMachineRes` ×3 queries, `StateMachineLit` ×1) for
the undecidable-guard problem class; external coverage came
from one bounded web read of the W3C SCXML Recommendation.
SCXML §5.9.1 Conditional Expressions is the decisive
external anchor: "If a conditional expression cannot be
evaluated as a boolean value ('true' or 'false') or if its
evaluation causes an error, the SCXML Processor MUST treat
the expression as if it evaluated to 'false' and MUST place
the error 'error.execution' in the internal event queue."
The only major standard that folds an unevaluable condition
into `false` does so only in tandem with a *mandatory
observable error* on a second channel; silent
absence-as-false has no prior-art support. Intrastate's
kernel has no second channel — the verdict is the only
observable — so an honest transplant surfaces the error in
the verdict itself.

In-repo, the sibling-path check found the discriminator
already exists rather than needing invention: the kernel's
`internal/resolve/resolve.go::GuardResult` is already
three-valued (`GuardFalse`, `GuardTrue`,
`GuardUnevaluable`), `resolve.go::KindGuardUnevaluable` is
already a closed refusal kind, and
`resolve.go::evaluateGuard` already answers the nil-seam
case — the same class of undecidability — with
`GuardUnevaluable`, never `GuardFalse`. RDR 0003's
normative operator vocabulary ("equality, membership,
bounded integer comparison, existence, and set
containment") already includes **existence**, which gives
authors an explicit, first-class way to test absence —
value-comparing operators do not need to double as implicit
existence tests.

### Key Discoveries

- **Documented** — SCXML §5.9.1 (quoted above): unevaluable
  → `false` is legitimate *only* paired with a mandated
  observable error signal; the spec never permits silent
  masking. Source: W3C SCXML Recommendation, fetched
  2026-08-11; quote cached in the research file.
- **Documented** — The kernel side of a three-valued
  contract is already built and frozen:
  `resolve.go::GuardResult` (three values),
  `resolve.go::KindGuardUnevaluable`, and the `gate`
  ordering (owned-state before undecidable-guard among
  survivors; GuardFalse prunes) are pinned by the ADV test
  suite. Choosing the partial-function branch requires no
  kernel behavior change.
- **Documented** — RDR 0003's operator vocabulary contains
  a dedicated existence operator (Normative Contracts,
  "The initial operator vocabulary MUST be closed and
  typed: equality, membership, bounded integer comparison,
  existence, and set containment"), so absence remains
  expressible by authors under a partial-function rule.
- **Documented** — RDR 0002's escape constraint ("An
  `escape` list MUST contain only resolver failure classes
  that RDR 0001 allows the table to model: `no_match` and
  `ambiguous_match`") makes `guard_unevaluable`
  structurally non-escapable, so routing absence to the
  unevaluable verdict closes the masking path by
  construction.
- **Assumed** — Strong-Kleene combination (definite results
  from present tags dominate; unevaluability propagates
  otherwise) never lets absence contribute falsity; needs
  the small truth-table derivation at Resolve (A2).

### Critical Assumptions

- **A1 RDR 0003's closed operator vocabulary needs no new
  operator for this rule: existence is the only operator
  whose semantics inspect absence, and every other operator
  can be classified value-comparing.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: RDR 0003 Normative Contracts, operator
    vocabulary block ("equality, membership, bounded
    integer comparison, existence, and set containment");
    verify each operator's declared value-kind acceptance
    classifies it cleanly into exactly one of the two
    classes.
  - **If wrong**: An operator that is neither total-over-
    presence nor value-comparing (e.g. a future default-
    valued comparison) has no domain rule, and the seam is
    unconstrained again for that operator.
- **A2 Under strong-Kleene combination of atom verdicts,
  no combination of present-tag verdicts and absent-tag
  unevaluability yields `GuardFalse` where absence
  contributed the falsity.**
  - **Status**: Pending
  - **Method**: Derivation
  - **Evidence**: Truth-table over `all` (conjunction) and
    `unless` (conjunctive exclusion, negated) for
    {T, F, U}; to be shown inline at Resolve.
  - **If wrong**: The masking path reopens through verdict
    combination even though single atoms are honest, and
    the contract's core guarantee fails.
- **A3 The kernel requires no code change: the
  `GuardUnevaluable` → `guard_unevaluable` mapping, the
  gate ordering, escape-candidate gating, AND the
  resolution-level aggregation rule (any surviving
  candidate with an unevaluable guard blocks selection — a
  decided-true sibling row is never selected past it)
  already implement the consumer side of this contract.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: `internal/resolve/resolve.go::gate`
    (undecidable survivors return a blocking refusal
    before selection), `resolve.go::evaluateGuard`,
    `resolve.go::escapeOrRefuse`; ADV-1/ADV-1b in
    `internal/resolve/adversarial_test.go` freeze D8.
    Premortem P-2/P-3 name the collapse this must rule
    out: a two-valued selection branch quietly treating
    unevaluable as pruned.
  - **If wrong**: This RDR's implementation touches frozen
    kernel behavior and must re-open RDR 0001 deviations
    rather than ship as spec-plus-fixtures; sibling-row
    selection past an unevaluable candidate would be the
    masking path reopened without any escape row.
- **A4 Prior art supports "definite results from present
  data dominate; unknown propagates otherwise" as the
  standard three-valued treatment of missing data
  (SQL/Kleene K3).**
  - **Status**: Pending
  - **Method**: Prior Art
  - **Evidence**: To cite at Resolve from a standard SQL
    reference or Kleene-logic text (corpus pass found no
    coverage; claim is currently from the model prior and
    is deliberately not load-bearing for the choice).
  - **If wrong**: The combination rule loses its external
    precedent and must stand on the A2 derivation alone —
    the choice survives, the citation does not.
- **A5 Authors who need "row applies when tag X is absent"
  — including the default-value pattern "absent OR equals
  v" — can express it with existence atoms, using the
  two-row pattern (one row guarded on absence via `unless`
  existence, one on the value) where a single row cannot
  carry the disjunction, without ambiguity under exact-one
  selection and without driving authors to stamp sentinel
  values upstream.**
  - **Status**: Pending
  - **Method**: Design Decision
  - **Evidence**: The decision: absence tests MUST be
    explicit existence atoms; implicit absence-as-false on
    value operators is rejected (Alternative 1). Verify at
    Resolve that RDR 0003's `all`/`unless` + existence
    composition covers the pattern, that the two rows
    cannot both match (absence and presence are mutually
    exclusive), and record the blessed pattern for the
    authoring docs. Premortem P-6 names the erosion risk:
    sentinel-stamping upstream defeats partiality wholesale.
  - **If wrong**: Authors cannot express a legitimate
    absence-conditional row and either demand a grammar
    extension from RDR 0003 or mask absence behind
    sentinels at the authoring layer.
- **A6 The assembled view is faithful for guard purposes:
  a tag key absent from the view means the artifact state
  is genuinely absent or unreadable — never dropped by
  assembly (merge precedence, provenance override, partial
  accessor read) — and the seam's treatment of
  "read failed" vs "genuinely absent" is stated rather
  than left to conflation.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: `internal/resolve/resolve.go::assemble`
    (merge and provenance precedence); RDR 0004's accessor
    contract for what an accessor must surface, may omit,
    and how a failed read is reported. Premortem P-5/P-9/
    P-11: the sharpest consumer is a total negative
    existence atom, which decides TRUE from absence — a
    dropped-but-real tag turns it into a plan against
    state that exists.
  - **If wrong**: Two horns — spurious non-escapable
    `guard_unevaluable` refusal storms on state that
    exists (value operators), or silent wrong plans
    (negative existence); operators cannot tell retryable
    read failure from genuine absence.
- **A7 RDR 0003's rule that a guard predicate is "a
  symbolic atom over a declared tag" gives review-time
  coverage for every guard-referenced key: a typo'd or
  undeclared key fails parse/lint at table load, so
  narrowing `RequiresOwned` does not leave guard inputs
  undeclared anywhere.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: RDR 0003 Normative Contracts, first
    block ("A guard predicate MUST be a symbolic atom over
    a declared tag…") plus the typed-operator rejection
    block; verify "declared tag" binds atoms to a checkable
    tag vocabulary, not merely to syntax. Premortem P-4/P-8
    names the failure this must answer.
  - **If wrong**: A misspelled key is forever absent —
    permanently unevaluable or permanently shadowed by
    FALSE-domination — invisible at review time and
    undiagnosable from the refusal; a vocabulary lint must
    then be added to RDR 0003's obligations before this
    RDR's rule is safe.

## Proposed Solution

### Approach

**The guard predicate is a partial function whose domain is
the set of tags it references; outside that domain the
evaluator answers "unevaluable," never "false."** Concretely:
value-comparing operators (equality, membership, bounded
integer comparison, set containment) require their
referenced tag key to be present in the assembled view — an
atom over an absent key is *unevaluable*. The existence
operator is the sole total operator: presence/absence is
exactly what it decides, so absence yields a decided
verdict. Atom verdicts combine under strong-Kleene
three-valued logic: a guard is `GuardFalse` or `GuardTrue`
only when the present tags alone decide it; any dependence
on an unevaluable atom makes the guard `GuardUnevaluable`.
The kernel (RDR 0001, already implemented) maps
`GuardUnevaluable` to `guard_unevaluable`, which RDR 0002
excludes from escape lists — so missing artifact state can
never be masked behind an escapable refusal class, which is
the user outcome this RDR exists to secure.

Ownership: **this RDR is the single normative home of the
domain rule.** RDR 0003's evaluator *implements* it (its A4
already claims the "unevaluable guard" semantic kind); RDR
0001's kernel *consumes* it unchanged. `Row.RequiresOwned`
is narrowed to its post-guard role: the owned tag keys the
row's transition writes depend on once the guard holds.
Guard decidability is carried by the guard's own
referenced-tag set under the domain rule — resolving kata
`xg7p`'s conflation without a kernel change. D8
(guard-FALSE prunes first) is ratified as safe *because of*
this rule: pruning only ever fires on falsity decided from
present tags.

### Technical Design

The contract constrains one seam — `GuardEvaluator.Evaluate`
(`internal/resolve/resolve.go::GuardEvaluator`) — and is
implemented by RDR 0003's evaluator when that RDR is built.
No kernel data flow changes: `assemble` builds the view,
`gate` prunes decided-FALSE rows, reports missing owned
state among survivors, then reports undecidable guards;
`escapeOrRefuse` applies the same gate to escape candidates.
This RDR adds (a) the normative domain rule below, (b)
doc-contract narrowing on `Row.RequiresOwned`, and (c) a
conformance fixture set the 0003 evaluator must pass.

#### Normative Contracts

```normative
Value-comparing guard operators (equality, membership,
bounded integer comparison, set containment) are PARTIAL
over the assembled evaluation view: an atom whose
referenced tag key is absent from the view MUST evaluate to
unevaluable — never to false and never to true.
```

```normative
The existence operator is the sole TOTAL operator: it MUST
decide true or false from the presence or absence of its
referenced tag key alone. Absence tests MUST be expressed
as explicit existence atoms; no value-comparing operator
may act as an implicit existence test.
```

```normative
Atom verdicts combine under strong-Kleene three-valued
logic across `all` and `unless`. A guard verdict MUST be
GuardTrue or GuardFalse only when the present tags alone
decide it; if the combined verdict depends on any
unevaluable atom, the evaluator MUST return
GuardUnevaluable. Absence MUST NOT contribute truth or
falsity to any value-comparing atom.
```

```normative
The kernel maps GuardUnevaluable to the refusal kind
`guard_unevaluable`, which is not a modelable escape class
(RDR 0002). Missing artifact state MUST NOT be maskable
behind an escapable refusal class.
```

```normative
Row.RequiresOwned names the owned tag keys the row's
post-guard transition (writes and clears) depends on. Guard
decidability is not RequiresOwned's job: a guard's input
coverage is enforced by the domain rule above. Listing a
guard-read owned key in RequiresOwned remains legal and
yields the more precise owned_state_unavailable diagnosis
among surviving rows.
```

```normative
Aggregation is resolution-level, not merely per-guard: if
any surviving candidate row's guard is GuardUnevaluable,
the resolution MUST refuse guard_unevaluable — a
decided-GuardTrue sibling row MUST NOT be selected while
an unevaluable candidate exists. Only decided-GuardFalse
rows are pruned from consideration.
```

```normative
Deviation D8 (guard-FALSE prunes first; a pruned row
contributes neither candidacy nor an owned-state
obligation) is ratified as normative, conditional on the
domain rule above: pruning is safe exactly because
GuardFalse can only arise from tags present in the view.
```

#### Load-Bearing Decisions

- **Selection / predicate** — when a guard's atoms mix
  decided and unevaluable verdicts, strong-Kleene
  combination selects the guard verdict: a definite
  `GuardFalse` derivable from present tags alone prunes
  (preserving D8); otherwise any unevaluable atom forces
  `GuardUnevaluable`. Why: refusing rows whose falsity is
  already decided would refuse legal resolutions the ADV
  suite proves must succeed.
- **Naming** — the verdict keeps the shipped names
  `GuardUnevaluable` / `guard_unevaluable`
  (`internal/resolve/resolve.go`); rejected: a new
  `guard_input_missing` kind (would open RDR 0001's closed
  five-kind taxonomy for a distinction the guard text
  already carries).

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Three-valued `GuardResult` + `guard_unevaluable` refusal kind | Predecessor (RDR 0001, implemented) | Available | Consumer side of the contract exists; no kernel change |
| Guard evaluator honoring the domain rule | RDR 0003 implementation (future) | Deferred | Conformance fixtures from this RDR bind it |
| Existence operator in the guard grammar | Predecessor (RDR 0003, Final) | Introduced (spec'd, unimplemented) | Carries all explicit absence tests (A5) |
| Escape-class closure (`no_match`, `ambiguous_match` only) | Predecessor (RDR 0002, Final) | Available | Makes `guard_unevaluable` non-escapable by construction |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Undecidability verdict | `internal/resolve/resolve.go::GuardUnevaluable` | Verdict carries no missing-tag payload | Reuse | Diagnostics naming the absent tag stay with RDR 0003's semantic kinds (Phase 3) |
| Refusal mapping + ordering | `internal/resolve/resolve.go::gate` | None | Reuse | Frozen by ADV suite; A3 verifies no change needed |
| Owned-state precheck | `internal/resolve/resolve.go::missingOwned` | Owned-provenance keys only | Reuse | RequiresOwned narrowed in doc contract only |

### Decision Rationale

The accretion-gate framing (no accretion recorded, but the
discipline applies): the undecided contract is **"when a
guard atom references a tag absent from the assembled view,
what verdict must the evaluator return — and which RDR owns
that rule?"** Approaches were enumerated as answers to that
question.

Scored matrix (criteria weighted by the Problem Statement's
user outcome — an actionable refusal, never a masked one):

| Criterion | A: Total, absence⇒false | B: Partial, absence⇒unevaluable (chosen) | C: Total over declared reads (precheck) |
| --- | --- | --- | --- |
| Correctness fit (masking closed?) | ✗ — absence⇒false prunes, `no_match` escapable: the exact masking path | ✓ — absence routes to `guard_unevaluable`, non-escapable per RDR 0002 | ✓ if declarations complete; ✗ silently when a guard-read tag is undeclared (drift) |
| Prior-art alignment | ✗ — SCXML pairs false with a mandatory error event; silent false has no support | ✓ — honest transplant of SCXML §5.9.1 into a one-channel kernel; matches K3 treatment of missing data (A4) | ~ — no prior system found declaring guard reads separately from guard text |
| Reversibility | ~ — flipping later reclassifies shipped table behavior | ✓ — spec + fixtures now; evaluator unbuilt, kernel unchanged | ✗ — adds a declared-reads surface to Row that must be maintained forever |
| Blast radius | Kernel unchanged, but D8 becomes unsound (prune on absence) | Kernel unchanged; D8 ratified sound; 0003 evaluator constrained | Kernel `gate` must grow a second precheck; ADV ordering re-opened |
| Cost | Low now, high at first masked escape | Low — the discriminator already exists end to end | Medium — new Row field or overloaded RequiresOwned (re-creates xg7p's conflation) |

The deciding rows are correctness fit and blast radius: B
is the only approach that closes the masking path without
touching frozen kernel behavior, and the sibling-path check
shows its discriminator already shipped
(`GuardResult`/`KindGuardUnevaluable`), so B reuses an
existing signal where A and C would bend or duplicate one.
A is rejected because it *is* the masking path. C is
rejected because declared reads duplicate what the guard
text already names (drift = silent unsoundness), extend the
exact `RequiresOwned` conflation this RDR exists to remove,
and re-open the ADV-proven gate ordering.

Premortem: hardened — critic verdict PASS, no switch
forced; interior semantics survived every attack, and the
findings cluster at the unwritten seams. Folded: P-2/P-3
(sibling-row masking through aggregation → new normative
aggregation clause; A3 scope extended to `gate`'s
any-undecidable-blocks behavior; named MVV tests), P-7
(refusal-kind accounting → answered: `guard_unevaluable`
is already one of the closed five kinds; the Naming
decision records the rejected sixth kind), P-4/P-8 (typo'd
or undeclared guard key invisible after `RequiresOwned`
narrowing → A7 added: RDR 0003's "symbolic atom over a
declared tag" rule is the review-time replacement; Failure
Modes updated), P-6 (absent-or-equals inexpressible in one
row → A5 widened to the two-row pattern and the
sentinel-stamping erosion risk), P-5/P-9/P-11/P-1
(fallible view assembly, sharpest through total negative
existence; read-failed vs absent conflation → A6 widened;
Risks and Failure Modes updated), P-10/P-12 (cross-team
evaluator drift → Phase 2 hardened into a normative golden
conformance vector suite).
Ledger: `docs/rdr/0007-guard-predicate-totality/evidence/propose-premortem/critic.md`.

Joint-check: clear (7 peers)

## Alternatives Considered

### Alternative 1: Total closed-world guard (absence ⇒ atom false)

**Description**: Keep the guard a total boolean function
over the view: an atom whose tag is absent simply does not
hold. This is SQL's WHERE-clause behavior (unknown rows are
dropped) and the superficially simplest evaluator.

**Pros**:

- Simplest evaluator; no three-valued combination logic.
- Matches the "row does not apply" intuition for many
  absent-tag cases.

**Cons**:

- Is exactly the masking path: absence ⇒ `GuardFalse` ⇒
  prune (D8) ⇒ obligation vanishes ⇒ `no_match` ⇒ escapable
  ⇒ plan, with nothing distinguishing "did not apply" from
  "could not tell" — the Problem Statement's failure.
- Makes every value operator an implicit existence test,
  making RDR 0003's dedicated existence operator redundant
  and author intent inexpressible.
- SCXML, the closest standard, only permits the false-fold
  alongside a mandatory observable error; the kernel has no
  second channel, so this transplant would be dishonest.

**Reason for rejection**: fails the deciding correctness
criterion — it converts missing artifact state into an
escapable refusal class by construction.

### Alternative 2: Totality via declared guard reads (view-completeness precheck)

**Description**: Keep guards two-valued but require every
guard-read tag to be declared (in `RequiresOwned` or a new
`Row.Reads` field). The kernel prechecks the view before
guard evaluation and refuses when a declared read is
absent, so the evaluator never sees an incomplete domain.

**Pros**:

- Guards stay boolean; no three-valued logic anywhere.
- Missing-input diagnosis can name the absent key directly
  from the declaration.

**Cons**:

- The declaration duplicates the referenced-tag set the
  guard text already names; drift between them is silent
  unsoundness — an undeclared read reverts to the
  unconstrained seam.
- Extends the `RequiresOwned` conflation (kata `xg7p`)
  rather than removing it, or adds a permanent new Row
  surface (RDR 0002's reviewable-data shape would also
  need reopening).
- Precheck ordering vs guard-FALSE pruning re-opens
  ADV-1's poisoning problem: a decided-FALSE row's declared
  reads could block a legal resolution — the exact defect
  D8 exists to prevent — so `gate` and its frozen tests
  would need rework.

**Reason for rejection**: highest blast radius (kernel gate
rework + Row shape + RDR 0002 reopen) to reach a guarantee
B achieves with a spec and fixtures.

### Briefly Rejected

- **Lint-only totality (statically reject guards over
  possibly-absent tags)**: observed-tag presence is a
  runtime fact; a static rule either rejects nearly every
  legitimate observed-tag guard or proves nothing — and
  RDR 0003's lint already rejects the one statically
  decidable case (owned read-before-write).
- **Evaluator-defined (leave it to RDR 0003's
  implementation)**: perpetuates the unconstrained seam
  this RDR exists to close; the implementer would decide a
  cross-RDR contract silently.
- **New refusal kind `guard_input_missing`**: reopens RDR
  0001's closed five-kind taxonomy for a distinction the
  refusal's guard text already carries.

## Trade-offs

### Consequences

- Positive: missing artifact state is never maskable behind
  an escape; the user always learns "I could not tell" as a
  distinct, non-escapable refusal.
- Positive: no kernel code change; D8/D5 and the ADV suite
  stand as shipped; kata `xg7p`'s conflation resolves by
  doc-contract narrowing.
- Positive: absence remains expressible — authors state it
  intentionally via existence atoms instead of relying on
  implicit closed-world reads.
- Negative: tables whose authors expected closed-world
  reads ("x == 1 just fails when x is missing") now refuse
  where they previously escaped or pruned; those guards
  must be rewritten with explicit existence atoms.
- Negative: `guard_unevaluable` is non-escapable by
  design, so an artifact in genuinely degraded state has no
  modeled-escape hatch through a guard that cannot read it;
  recovery is fixing the artifact state or the accessor,
  not the table.

### Risks and Mitigations

- **Risk**: RDR 0003's evaluator is implemented later by a
  session that drifts from this rule (e.g. folds absence
  into false for one operator).
  **Mitigation**: The conformance fixture set (Phase 2) is
  normative and lives in the kernel's test tree; 0003's
  implementation must run it, and the Prerequisites bind
  0003's implement stage to this RDR being Final.
- **Risk**: Strong-Kleene combination has a subtle case
  where absence leaks into a decided verdict (e.g. `unless`
  negation).
  **Mitigation**: A2's truth-table derivation at Resolve;
  MVV scenario exercises the mixed-verdict combinations.
- **Risk**: View assembly drops state the artifact
  actually has — producing either spurious non-escapable
  `guard_unevaluable` refusal storms (value operators) or,
  worse, silent plans through a total negative-existence
  atom deciding TRUE on a dropped-but-real tag.
  **Mitigation**: A6 verifies assembly faithfulness and
  the read-failed vs absent distinction against `assemble`
  and RDR 0004's accessor contract; the recovery path for
  an unevaluable refusal is operational (fix the artifact
  state or the accessor), never a table escape, and the
  authoring docs must say so.
- **Risk**: Authors defeat partiality by stamping sentinel
  values upstream so keys are never absent, making value
  comparisons decide on placeholders.
  **Mitigation**: A5's blessed two-row pattern gives the
  legitimate outlet; authoring guidance flags
  sentinel-stamping as the anti-pattern (premortem P-6).

### Failure Modes

- **Visible break**: a previously escaping table now
  refuses `guard_unevaluable`. Diagnosis: the refusal
  carries the guard text (`Refusal.Guard`) and the
  contributing rows; the named tags point at the absent
  state. Recovery: make the artifact state readable, or
  rewrite the guard with an explicit existence atom if
  absence was intended.
- **Silent failure guarded against**: absence folded into
  `GuardFalse` by a drifting evaluator would re-open
  masking with no visible symptom; the conformance fixtures
  are the tripwire (they fail loudly on any false-fold).
- **Diagnostic gap**: `GuardResult` carries no missing-tag
  payload, so "which tag was absent" is inferable only from
  the guard text; if that proves insufficient in practice,
  the enrichment belongs to RDR 0003's semantic-kind
  diagnostics (Phase 3 hands it off), not to a kernel enum
  change.
- **Review-time blind spot guarded against**: a typo'd or
  undeclared guard key would be permanently absent —
  forever unevaluable, or forever shadowed when
  FALSE-domination decides the guard from its other atoms
  (premortem P-4/P-8). The catch is RDR 0003's
  declared-tag rule at parse/lint (A7): the defect fails
  table load, never reaching a production refusal.
- **Conflated recovery signal**: `guard_unevaluable`
  cannot itself distinguish transient read failure
  (retryable) from genuine absence (not retryable); until
  A6 settles what the accessor contract reports, operators
  should treat the refusal as "inspect the artifact and
  accessor," not "retry" (premortem P-11).

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified
- [ ] This RDR Final **before** RDR 0003's implementation
      begins (triage sequencing: D8 ratification precedes
      the evaluator build)

### Minimum Viable Validation

Kernel-level tests (ADV-style, with a stub evaluator
conforming to the domain rule) proving the masking probe
inverts: a guard referencing an absent tag, plus an absent
`RequiresOwned` key, plus a modeled `no_match` escape,
yields `guard_unevaluable` with `Escaped:false` — while the
same table with the tag present and decided FALSE still
prunes and escapes (D8 preserved). Two named scenarios from
the premortem are in scope: *unevaluable-not-no_match* (one
row, guard over an absent key → `guard_unevaluable`, never
`no_match`) and *unevaluable-blocks-true-sibling* (row A
unevaluable beside row B decided true → refusal, no plan).
Mixed-verdict combination cases (decided-FALSE beside
unevaluable, in `all` and in `unless`) assert the
strong-Kleene selection.

### Phase 1: Contract doc alignment

Narrow `Row.RequiresOwned`'s doc contract to post-guard
write-dependency keys and state the domain rule on
`GuardEvaluator` — doc comments in
`internal/resolve/resolve.go` only; no behavior change.

### Phase 2: Normative golden conformance vectors

Encode the domain rule and the MVV scenarios as a named,
versioned golden vector suite in the kernel's test tree —
operators × present/absent × strong-Kleene verdict pairs ×
`unless` atoms × row aggregation — that any
`GuardEvaluator` implementation (RDR 0003's build) must
pass verbatim: the enforcement point for cross-team
evaluator drift (premortem P-10: unevaluable-producing
inputs are rare in dev data, so a diverging evaluator
otherwise ships green).

### Phase 3: Diagnostics and authoring handoff

Record with RDR 0003's surfaces: the "name the absent tag"
diagnostic obligation (its A4 semantic kinds), confirmation
that the declared-tag rule is a checkable vocabulary lint
(A7), and the authoring guidance — blessed two-row
absence pattern, sentinel-stamping anti-pattern (A5) — so
unevaluable refusals become fully actionable without a
kernel change.

## Validation

### Testing Strategy

[Required — never omit. Test scenarios and coverage goals — what to test and
what constitutes "done." For non-functional concerns
(performance, security): state measurement strategy,
not estimates.]

1. **Scenario**: [Description]
   **Expected**: [Result]

## Finalization Gate

> Complete each item with a written response before
> marking this RDR as **Final**. Written responses
> prevent rubber-stamping and produce a review record.
>
> First run the mechanical pre-sweep
> (`prompts/gate/tooling-pass.md`): TEMPLATE section
> coverage, Method-label vocabulary, `Source Search`
> self-reference, `Docs Only` on load-bearing claims. It
> catches what the review rounds disturbed; resolve any
> BLOCK before the written responses below.

### Contradiction Check

[State any conflicts between Research Findings and
the Proposed Solution. If none exist, state
"No contradictions found between research findings,
design principles, and proposed solution."]

### Assumption Verification

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

[Confirm the Minimum Viable Validation is in scope
and will be executed during implementation, not
deferred. State the specific test or proof.]

### Cross-Cutting Concerns

[List only concerns that apply to this RDR. For each,
state either how this RDR addresses it, or which peer
RDR owns the project-wide policy this RDR conforms
to. Omit (rather than N/A-bullet) anything that does
not apply.]

### Proportionality

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
Profile, correct the field and do not lock until the
missing lenses have run. Also confirm form: value +
one clause naming the contract(s); strip any
matrix/provenance prose left from the template or Seed
(it belongs in the template comment, not the instance).]

## References

- W3C SCXML Recommendation, §5.9.1 Conditional Expressions
  (error-in-cond rule; quoted in
  `docs/rdr/0007-guard-predicate-totality/evidence/research/propose-prior-art.md`)
- RDR 0001 — Resolution kernel (refusal taxonomy; selection
  rule; deviations D5/D8 in
  `docs/rdr/0001-resolution-kernel/artifacts/deviations.md`)
- RDR 0002 — Transition table as reviewable data (escape
  class closure)
- RDR 0003 — Guard predicate exhaustiveness (operator
  vocabulary; A4 "unevaluable guard" semantic kind)
- `internal/resolve/resolve.go` (`GuardResult`,
  `GuardEvaluator`, `gate`, `evaluateGuard`,
  `missingOwned`, `escapeOrRefuse`)
- `internal/resolve/adversarial_test.go` (ADV-1, ADV-1b)
- kata `xg7p` (originating finding)
