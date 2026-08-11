# Recommendation 0007: Guard predicate totality over an incomplete evaluation view

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-08-09
- **Status**: Draft
- **Type**: Architecture
- **Profile**: foundational — one contract, the domain of
  guard-predicate evaluation over an incomplete view;
  cross-RDR producer, implemented by RDR 0003's evaluator
  and consumed by RDR 0001's kernel.
- **Priority**: Medium
- **Related Issues**: kata `xg7p` — "RDR 0001:
  RequiresOwned conflates guard-input state with post-guard
  transition state" (`src:roborev`, `severity:medium`,
  `area:internal-resolve`).
- **Predecessors**: 0001-resolution-kernel,
  0003-guard-predicate-exhaustiveness
- **Overrides**: None superseded. This RDR *fixes* the
  meaning of `Row.RequiresOwned` (to post-guard
  write-dependency keys) — a field no RDR ever defined,
  carrying only a Go doc comment — and *adds* an
  evaluation-domain obligation to RDR 0003's guard
  evaluator; it is the single normative home of the domain
  rule, which both peers cite rather than restate. No peer
  RDR text is reopened.
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

Nothing normative states which answer is required. The only
statement of `Row.RequiresOwned`'s meaning anywhere is a Go
doc comment — `internal/resolve/resolve.go`'s "RequiresOwned
names owned tag keys the row's evaluation needs" — which
reads as covering both the state a guard reads to be
decidable and the state a transition's writes depend on once
the guard holds. No RDR defines the field at all: RDR 0001
introduced it as an implementation surface without fixing
its meaning, and RDR 0003 does not mention it. So a table
author and a guard evaluator can each hold a defensible but
incompatible reading. This RDR
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

RDR 0003 is **Final** and unimplemented. Its normative
blocks — grammar, typing, `all`/`unless`, exhaustiveness,
row-group scoping, set-domain proof, diagnostics, provenance
lint — contain no rule for what an atom evaluates to when
its tag is absent from the runtime view; its only
absent-state rule (owned read-before-write) is a lint
rejection, not a runtime verdict. The fork is live.

Sequencing. D8 must be ratified before RDR 0003 is
implemented, because 0003's guard-evaluator contract depends
on how this fork resolves — which makes 0007 the
sequencing-critical member of its seed batch. D8 is this
fork's downstream consequence and is settled by whichever
branch is chosen.

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

**What this citation does and does not authorize.** SCXML's
*disposition* is to keep running: the machine folds the
condition to false, raises `error.execution`, and continues.
This RDR halts resolution instead. So the citation supports
the narrow claim it is used for — that an undecidable
condition MUST NOT silently become `false`, and that the
error must be observable — and it supports nothing else. It
does not underwrite the resolution-level veto (a table-wide
halt is strictly stronger than SCXML's continue), and it
does not underwrite omitting the absent tag from the
refusal — SCXML mandates a payload naming what failed, which
is the diagnostic gap this RDR records under Failure Modes
and hands to Phase 3. On both of those points the citation
cuts *against* the current design, and they stand on the
in-repo house rule (A8) and the accepted-cost argument in
Decision Rationale, not on SCXML.

In-repo, the sibling-path check found both the *principle*
and the *discriminator* already present rather than needing
invention. The principle is settled at the adjacent seam:
RDR 0004 requires a gate accessor's `indeterminate` to be "a
refusal-class result, not a false allow and not a false
deny" (A8) — the same refusal to collapse an undecided third
value, one seam over, carried into RDR 0005 and shipped as
`flow-gate-indeterminate`. This RDR therefore applies a
house rule rather than transplanting a foreign one; SCXML
and K3 corroborate a choice the codebase has already made
elsewhere. The discriminator likewise exists: the kernel's
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

- **Documented** — SCXML §5.9.1 (quoted under
  *Investigation*): unevaluable → `false` is legitimate
  *only* paired with a mandated observable error signal;
  the spec never permits silent masking. Source: W3C SCXML
  Recommendation, fetched 2026-08-11; quote cached in the
  research file.
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
  operator for this rule: `exists` is the only operator
  whose declared semantics name absence, and the remaining
  four are value-comparing. RDR 0003 is *silent* on
  absent-operand behavior for those four, so this RDR fixes
  that rule rather than restating one — consistent with,
  but not stated by, RDR 0003.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: RDR 0003 Normative Contracts —
    "The initial operator vocabulary MUST be closed and
    typed: equality, membership, bounded integer
    comparison, existence, and set containment. Unknown
    operators MUST be rejected during parse or lint before
    resolution." Vocabulary is closed with no runtime
    extension point. RDR 0003's operator/kind matrix marks
    `exists` alone as absence-inspecting ("Tests presence
    or absence, not value equality"; accepted kind
    "optional scalar or optional set-valued tag"), while
    `eq`/`in`/`lt,lte,gt,gte`/`contains` are stated in
    value terms ("compares a tag value to one typed
    literal"; "checks a scalar tag against a typed literal
    set"; "Narrows a bounded integer domain by
    comparison"; "Narrows the set-valued domain to
    assignments containing every listed element"). A sweep
    of RDR 0003 for `absent|absence|optional|default`
    returns only the `exists` clauses — no absent-operand
    or default-value rule exists for the other four. This
    RDR's Normative Contracts supply it, and pin
    `contains`: an absent set-valued tag is *unevaluable*,
    not the empty set.
  - **If wrong**: An operator that is neither total-over-
    presence nor value-comparing (e.g. a future default-
    valued comparison) has no domain rule, and the seam is
    unconstrained again for that operator.
- **A2 Under strong-Kleene combination of atom verdicts,
  no combination of present-tag verdicts and absent-tag
  unevaluability yields `GuardFalse` where absence
  contributed the falsity.**
  - **Status**: Verified
  - **Method**: Derivation
  - **Evidence**: Strong-Kleene conjunction is `min` under
    the truth order `F < U < T`; negation is
    `¬T=F, ¬F=T, ¬U=U`. The row-verdict *shape*
    `all_result ∧ ¬(unless_conj)` is RDR 0003's — "Positive
    `all` predicates are conjunctive requirements. Negative
    `unless` predicates are also conjunctive within the
    excluded predicate set: if all `unless` predicates
    hold, the candidate row is disabled", and "`unless` is
    not per-atom negation." Note that RDR 0003 states this
    set-theoretically, over *complete* assignments in a
    finite declared domain ("a row's accepted assignments
    are the intersection of all positive `all` atom domains
    minus the single conjunctive assignment set matched by
    the row's full `unless` block"), where no third value
    can arise. Lifting that shape to a three-valued runtime
    is **this RDR's extension**, not a restatement of 0003 —
    which is precisely the gap A1 records 0003 as leaving
    open. The shape is inherited; the K3 semantics below are
    fixed here. Conjunction table:

    | ∧ | T | F | U |
    |---|---|---|---|
    | **T** | T | F | U |
    | **F** | F | F | F |
    | **U** | U | F | U |

    `GuardFalse` arises in exactly two ways, and neither
    admits absence as the source of falsity:
    (i) `all_result = F` requires some atom `aᵢ = F` — a
    *present* tag whose value comparison failed. An
    absent-key value atom is `U`, and `U` never yields `F`
    in a conjunction (the `U` row/column contains no `F`
    except where an `F` operand is already present).
    (ii) `¬(unless_conj) = F` requires `unless_conj = T`,
    which requires *every* `uⱼ = T`; a `U` atom cannot
    contribute `T`, so an absent key blocks this branch
    rather than driving it.
    The `F ∧ U = F` cell is the case that looks like a leak
    and is not: the falsity is *witnessed* by the
    present-tag atom alone, and holds under every
    substitution of `T`/`F` for the unevaluable atom — so
    pruning discards no row that could have survived. This
    is precisely what preserves D8. Absence reaches a
    decided verdict only through `exists`, the operator
    whose whole job is deciding presence. Cross-checked
    against the K3 tables in A4's citation (PostgreSQL 16
    §9.1: `FALSE AND NULL = FALSE`), which is the *strong*
    Kleene variant — weak/Bochvar would give `U` here.
  - **If wrong**: The masking path reopens through verdict
    combination even though single atoms are honest, and
    the contract's core guarantee fails.
- **A3 The kernel requires no behavior change (Phase 1
  edits doc comments only): the `GuardUnevaluable` →
  `guard_unevaluable` mapping, the gate ordering,
  escape-candidate gating, AND the resolution-level
  aggregation rule (any surviving candidate with an
  unevaluable guard blocks selection — a decided-true
  sibling row is never selected past it) already implement
  the consumer side of this contract.**
  - **Status**: Verified — after correcting the escape
    half of the aggregation clause to match frozen kernel
    behavior (see Evidence).
  - **Method**: Source Search
  - **Evidence**: Mapping — `resolve.go::GuardUnevaluable`
    → `resolve.go::KindGuardUnevaluable`, realized in
    `resolve.go::gate`; the five-kind set is pinned closed
    by `resolve_test.go::TestReq7_RefusalKindSetIsExactlyTheFiveNamedKinds`.
    Ordering — `gate` calls `missingOwned(survivors, view)`
    (not `rows`), so GuardFalse prunes before the
    owned-state sweep, which precedes the undecidable
    loop; D8's prune-first half is frozen by
    `adversarial_test.go::TestAdv1_GuardFalseRowsRequiresOwnedMustNotPoisonAnExactOneMatch`
    and `::TestAdv1b_AllGuardsFalseIsNoMatchNotOwnedStateUnavailable`
    (both fail if `rows` replaced `survivors`).
    Resolution-level aggregation — `gate`'s
    `len(undecidable) > 0` return overwrites the named
    return `selected` with `nil`, discarding accumulated
    GuardTrue rows, so a decided-true sibling is never
    selected past an unevaluable candidate; confirmed by
    spike Probe B
    (`evidence/spikes/aggregation-probe.md`). A decidable
    escape row cannot mask an unevaluable candidate:
    `Resolve` returns on `gate(candidates,…)`'s block
    before `escapeOrRefuse` is reachable.
    **Correction found at Resolve**: the clause's other
    half — "an unevaluable escape row MUST NOT convert a
    candidate-set refusal into `guard_unevaluable`" — was
    refuted. `resolve.go::escapeOrRefuse` returns the
    escape set's blocking refusal in place of the
    candidate-set refusal (`if blocked != nil { return
    refuse(in, *blocked) }`), and this is frozen
    deliberately by
    `adversarial_test.go::TestAdv2_EscapeEdgeMustNotBypassTheGuardSeam`
    subtest "guard UNEVALUABLE must not rescue",
    `fixup_test.go::TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue`,
    and `fixup_test.go::TestFixupGateIsUniformAcrossOrdinaryAndEscapeCandidates`.
    Spike Probe A reproduced it against the shipped kernel.
    The Normative Contracts aggregation block was corrected
    to the shipped behavior rather than the kernel changed;
    Phase 1 therefore remains doc-comments-only.
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
  - **Status**: Verified
  - **Method**: Prior Art
  - **Evidence**: PostgreSQL 16 Documentation, Ch. 9
    "Functions and Operators", §9.1 "Logical Operators",
    p. 231 — "SQL uses a three-valued logic system with
    true, false, and null, which represents 'unknown'",
    followed by the full AND/OR/NOT truth tables. The two
    load-bearing rows are explicit: `FALSE AND NULL =
    FALSE` and `TRUE OR NULL = TRUE` — definite results
    from present data dominate, unknown propagates
    otherwise. The `FALSE AND NULL = FALSE` row is the
    discriminating case identifying this as *strong*
    Kleene (K3, conjunction as `min` under `F < U < T`)
    rather than weak/Bochvar, where `U` is infectious. The
    spec's commutativity note ("you can switch the left
    and right operands without affecting the result")
    additionally rules out a short-circuit reading, so the
    dominance is semantic rather than evaluation-order
    dependent — which matters for an evaluator free to
    visit atoms in any order. Quote cached in
    `evidence/research/resolve-citations.md` (C3).
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
  - **Status**: Pending — the runtime half is settled; the
    **load-time** half is refuted pending A12. See *Carried
    constraint*.
  - **Method**: Design Decision (runtime) + Source Search
    (load-time, owed)
  - **Evidence**: The decision: absence tests MUST be
    explicit existence atoms; implicit absence-as-false on
    value operators is rejected (Alternative 1). Verified
    against RDR 0003: (a) `unless` accepts any operator —
    the placement rule is "Positive guard atoms MUST live
    in `all`; negative guard atoms MUST live in `unless`",
    a placement constraint only, and RDR 0003's own spike
    fixture (`0003-…/evidence/spikes/guard-fixture.toml`)
    shows `exists` in `all` and a non-`eq` operator in
    `unless`, so `unless … exists = true` is grammar-legal;
    (b) no intra-guard disjunction exists — RDR 0003 A3
    states "Inline `not` and nested boolean expressions
    are rejected", and its Alternative 3 (boolean
    expression grammar) was rejected with row factoring
    accepted as the cost — so the two-row pattern is the
    only route; (c) the two rows are provably disjoint: no
    assignment has X both absent and present, and row-group
    scoping "MUST be scoped to a normalized row group …
    as one product" narrows what is compared without
    merging distinct assignments, while "`unless` … does
    not create source-order priority" keeps disjointness
    independent of ordering. Exact-one selection is
    therefore satisfied. Blessed pattern recorded for the
    Phase 3 authoring handoff. Premortem P-6 names the
    erosion risk: sentinel-stamping upstream defeats
    partiality wholesale.
  - **Carried constraint**: the sub-clauses above establish
    disjointness at *runtime*, over assembled views, which is
    the wrong forum for the binding risk. The author meets
    **table load** first, and there the pattern is not
    obviously legal. RDR 0003 computes "overlap is any
    non-empty `row_i accepted assignments intersect row_j
    accepted assignments`" over a *finite domain product*
    built from declared domains, and no tag declaration in
    the tree carries an absent member
    (`0003-…/evidence/spikes/guard-fixture.toml`:
    `cluster_eligible domain = [true, false]`,
    `prelock_iterations min = 0 / max = 3`). If absence is
    not an element of that product, the absence-row and the
    value-row may project onto the same assignment set.
    **Overlap is not downgradeable**: RDR 0002 retains
    "ambiguous overlap" among the stable data-level
    validation categories a table MUST be rejected on, unlike
    the *exhaustiveness* claim RDR 0003 permits lint to
    refuse or downgrade. An earlier draft of this constraint
    scoped the question to exhaustiveness and therefore
    under-read the risk. Whether RDR 0003's row-group scoping
    ("MUST be scoped to a normalized row group … as one
    product") narrows the two rows out of a shared product,
    or whether absence needs standing in the domain, is A12.
  - **If wrong**: Authors cannot express a legitimate
    absence-conditional row and either demand a grammar
    extension from RDR 0003 or mask absence behind
    sentinels at the authoring layer — the anti-pattern whose
    only sanctioned alternative *is* this pattern, so the
    Risks mitigation for sentinel-stamping fails with it.
- **A6a Kernel-side assembly is faithful: no tag key
  present in the input tuple is ever dropped from the
  assembled view by merge precedence or provenance
  override. Presence is monotone in the input; shadowing
  changes a value, never a key's presence.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/resolve/resolve.go::assemble`
    — every write is an unconditional map assignment into
    `view.tags` keyed by tag key; there is no `delete`, no
    filter, no conditional skip, no key rewriting, and no
    error return. The single conditional
    (`if in.Recognized != ""`) gates an *insert*, never a
    removal. The view has exactly four
    readers — `TagSet.Lookup`, `TagSet.has`, `TagSet.matches`
    (the candidate-selection path), and `TagSet.Len` — and
    none removes a key: the first two read the map directly
    with no filtering, `matches` is a match-pattern test
    over `want`, and `Len` is a count. Owned-over-observed
    precedence overwrites a `taggedValue` in place, so the
    key survives with `ProvenanceOwned` — value-level
    shadowing only, which A6a does not claim about and
    which cannot manufacture absence. Since absence is what
    `exists` decides on, shadowing cannot flip a negative
    existence atom.
  - **If wrong**: Spurious non-escapable
    `guard_unevaluable` refusal storms on state that
    exists (value operators), or silent wrong plans
    (negative existence) — premortem P-5/P-9.
- **A6b OPEN: whether a *failed* accessor read is
  distinguishable from *genuine* absence at the
  accessor→kernel boundary is not stated by any current
  contract.**
  - **Status**: Pending — open obligation on RDR 0004 (or
    wherever the accessor executor lands); not resolvable
    inside this RDR.
  - **Method**: Source Search
  - **Evidence**: RDR 0004's **read-accessor** clause is
    the whole of what governs this: "A read accessor MUST
    return typed tag values or a typed refusal. It MUST NOT
    mutate authoritative artifacts." This is a disjunction
    with **no completeness requirement** on the "typed tag
    values" branch: nothing says the returned values are the
    complete tag set, and nothing requires a partially-read
    artifact to take the refusal branch. RDR 0004's named
    failure taxonomy ("artifact unavailable, timeout,
    execution failure, gate denied, gate indeterminate,
    write attempted for a non-owned tag, read-back
    mismatch") covers whole-artifact failure but not a
    corrupt or partially-parsed artifact where some tags
    read and others did not; its own definition of "silent
    failure" is scoped exclusively to the *write* path,
    guarded by read-back, with no read-side analogue. RDR
    0004 *does* rule on third-value collapse for the **gate**
    accessor (see A8) — but that clause governs a reported
    execution outcome, not read completeness, so it does not
    close this gap.
    A sweep of `docs/rdr/` for read-failure-vs-absence
    language returns hits only inside this assumption's own
    text. Because `Input.Owned []Tag` carries no error
    channel, a truncated snapshot and a genuinely-absent
    tag are byte-identical at the kernel boundary — so
    A6a's faithfulness is faithfulness *to the tuple*, and
    does not rescue A6b.
  - **Partial mitigation in kernel**:
    `resolve.go::missingOwned` raises
    `owned_state_unavailable` for keys a row declares in
    `RequiresOwned`, converting absence into a refusal
    rather than a decision — but only for declared keys. A
    guard-read key with no `RequiresOwned` entry gets no
    such protection, which is exactly the surface this
    RDR's narrowing leaves to the domain rule.
  - **If wrong / while open**: `guard_unevaluable` cannot
    distinguish transient read failure (retryable) from
    genuine absence (not retryable), and a dropped-but-real
    tag makes a negative existence atom decide TRUE,
    yielding a plan against state that exists. Until the
    read-completeness seam is stated, operators must treat
    an unevaluable refusal as "inspect the artifact and
    accessor," never "retry" (see Failure Modes).
- **A7 RDR 0003's rule that a guard predicate is "a
  symbolic atom over a declared tag" gives review-time
  coverage for every guard-referenced key: a typo'd or
  undeclared key fails parse/lint at table load, so
  narrowing `RequiresOwned` does not leave guard inputs
  undeclared anywhere.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: "Declared tag" binds to a checkable
    vocabulary, not to syntax. RDR 0003 Normative
    Contracts: "A guard predicate MUST be a symbolic atom
    over a declared tag, not a host language callback and
    not a free-form expression string." Its Technical
    Design makes resolution an obligation — "The tag name
    must resolve to a declared tag" — and names the
    declaration's payload: "Each tag declaration supplies
    the value kind and, when lint must prove
    exhaustiveness, the finite domain." A vocabulary
    supplying value kind and finite domain is an
    enumerated set, not a syntactic category. RDR 0003
    Failure Modes lists `unknown tag` among typed
    parse/lint failures, and its Testing Strategy Scenario
    4 requires each such failure be "rejected before
    resolution" — load time, not production. The
    vocabulary itself lives in RDR 0002 Normative
    Contracts: "The model MUST declare every tag it
    matches or writes, including each tag's provenance:
    owned, observed, or recognized", with `unknown tag`
    among the stable data-level validation categories that
    "MUST" be retained. Crucially that obligation is keyed
    on *matching*, independent of ownership — so narrowing
    `RequiresOwned` cannot leave a guard-referenced key
    undeclared. Premortem P-4/P-8 names the failure this
    answers.
  - **Anchor note**: the normative MUST-reject for an
    undeclared key is RDR 0002's validation-category block;
    RDR 0003 states the binding ("must resolve to a
    declared tag") and the diagnostic, but its
    `unknown tag` rejection sits in Failure Modes /
    Testing Strategy prose rather than inside a
    ```normative``` block. Cite RDR 0002 for the hard
    obligation.
  - **If wrong**: A misspelled key is forever absent —
    permanently unevaluable or permanently shadowed by
    FALSE-domination — invisible at review time and
    undiagnosable from the refusal; a vocabulary lint must
    then be added to RDR 0003's obligations before this
    RDR's rule is safe.
- **A8 The no-collapse principle this RDR applies is
  already settled in-repo at the adjacent accessor seam, so
  this RDR conforms to a house rule rather than importing
  one: RDR 0004 requires an undecided third value to be a
  refusal class, never folded into either boolean.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: RDR 0004 Normative Contracts — "A gate
    accessor MUST return allow, deny, or indeterminate.
    Indeterminate MUST be a refusal-class result, not a
    false allow and not a false deny." Carried forward
    normatively by RDR 0005 ("MUST NOT … coerce gate allow,
    deny, or indeterminate results into tag values") and
    shipped as the CLI code `flow-gate-indeterminate`. The
    seams are distinct and this one is not pre-empted: 0004's
    gate accessor *reports* indeterminate as an execution
    outcome, whereas this RDR *derives* unevaluability from
    an absent key in the assembled view — 0004 states no
    absent-key rule. The principle is the same and the
    phrasing is near-identical ("not a false allow and not a
    false deny" / "never to false and never to true"), which
    is why the chosen approach is the consistent one for
    this codebase.
  - **If wrong**: The rule stands on external prior art
    (SCXML, K3) alone and loses its in-repo consistency
    argument; the choice survives, the house-rule framing
    does not.
- **A9 The empty-atom-block identities this RDR fixes do not
  contradict RDR 0002's normalized row shape: an empty
  `unless` block is representable-and-absent rather than a
  vacuously-true conjunction, and RDR 0002's normalization
  does not synthesize an empty `unless` block that would
  disable rows under `all ∧ ¬(unless_conj)`.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: Needed. RDR 0003 states the row-verdict
    shape but not the empty-block identity, and this RDR now
    fixes it (Normative Contracts). Verify against RDR 0002's
    normalization contract whether an omitted `unless` block
    normalizes to absent or to an empty collection, and
    confirm no lint/exhaustiveness rule in RDR 0003 depends
    on the opposite reading.
  - **If wrong**: Either the identity belongs to RDR 0002's
    normalization (not here), or a table with an omitted
    `unless` block disables every row carrying one — a
    silent whole-table failure.
- **A10 The atom structure this RDR's domain rule quantifies
  over (operator, referenced tag key, literal, `all`/`unless`
  placement) is recoverable by the evaluator from what it is
  handed at the `GuardEvaluator.Evaluate` seam — i.e. some
  total, lossless mapping from the authored guard to
  `Row.Guard string` exists or can be specified by RDR 0003
  without reopening RDR 0001's `Row` shape.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: Needed. Established so far: `Row.Guard` is
    an opaque `string` and `Evaluate(guard string, view
    TagSet)` is the whole seam
    (`internal/resolve/resolve.go::GuardEvaluator`); the
    kernel attaches no meaning to the text
    (`evaluateGuard` only special-cases `""`). Guards are
    authored *structurally* — RDR 0003's spike fixture uses
    `[rule.guard.all.profile]`, `[rule.guard.unless.
    prelock_iterations]` — and RDR 0002's normalization
    "MUST combine both into one candidate-row" predicate
    set. No RDR states what that combination yields at the
    `Row.Guard` boundary, and the only in-tree consumer
    (`fixtures_test.go::fixtureGuards`) treats it as an
    opaque *name* keyed into `decided map[string]bool`.
    Verify against RDR 0002's normalization contract and
    RDR 0003's grammar whether the mapping is specified
    anywhere; if not, it is an obligation handed to RDR 0003
    (Phase 3), not a defect fixed here.
  - **If wrong**: Every atom-level normative clause in this
    RDR is unimplementable as written and the Phase 2
    vectors for scenarios 4/6/7/8 cannot be encoded at all —
    the domain rule would bind a seam that cannot see what
    the rule is about. Raised by the 3amigo lens
    (IMP-1/QA-1) as the hinge finding.
- **A11 Deciding guard presence provenance-blind conflicts
  with no shipped kernel behavior: no code path requires a
  guard atom to see only owned tags.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: The view has exactly four readers.
    `TagSet.Lookup` is provenance-blind (returns `ok` on any
    key present, `resolve.go` L113–119) and is the only
    *exported* accessor — the one an evaluator outside the
    package can use at all. `TagSet.has` takes a provenance
    and is unexported, with a single caller,
    `missingOwned`, which is scoped to `RequiresOwned` and
    the `owned_state_unavailable` diagnosis — a different
    question from guard decidability. `TagSet.matches`
    (candidate selection) compares key and value with no
    provenance test, and `Len` is a count. So the
    provenance-blind reading is both the only one an
    external evaluator can implement and the one already
    used everywhere except the deliberately owned-scoped
    precheck. RDR 0003's fixture declares
    `[tags.cluster_eligible] provenance = "observed"` and
    guards it with `exists = true`, confirming the spec'd
    intent is not owned-only.
  - **If wrong**: `exists` over an observed or recognized
    tag decides differently than over an owned one, RDR
    0003's own fixture breaks, and absence tests become
    silently provenance-dependent.
- **A12 The blessed two-row absence pattern (A5) survives
  RDR 0003's load-time OVERLAP check, not merely runtime
  evaluation: either row-group scoping keeps the
  absence-guarded row and the value-guarded row out of one
  compared product, or absence has standing in the declared
  domain product — otherwise the pattern is rejected at table
  load under RDR 0002's mandatory `ambiguous overlap`
  category.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: Needed. Established so far: RDR 0003
    Technical Design defines overlap as "any non-empty
    `row_i accepted assignments intersect row_j accepted
    assignments`" over a product of declared finite domains,
    and normatively requires coverage/overlap checks "MUST be
    scoped to a normalized row group … as one product". RDR
    0002 lists "ambiguous overlap" among validation
    categories that "MUST retain stable data-level
    categories" — a hard reject, not a downgradeable claim
    (contrast RDR 0003's exhaustiveness permission ceiling,
    which explicitly allows lint to "refuse or downgrade").
    No tag declaration in the tree admits an absent element.
    The asymmetry is textual: every "refuse or downgrade"
    escape in RDR 0003 attaches to the *exhaustiveness*
    claim, and none of them mentions overlap — so the weaker
    claim is the only one with a relief valve, while the hard
    reject runs over a product that cannot represent absence.
    RDR 0003's own fixture sits on the fault line: the
    `foundational-to-cove` rule guards `cluster_eligible`
    with `exists = true` over a tag whose declared domain is
    `[true, false]`. Verify whether an existence atom over a
    dimension with no absent member makes a row's
    accepted-assignment set empty, total, or undefined under
    0003's proof, and whether the two rows share a normalized
    row group at all.
  - **If wrong**: A5's blessed pattern is unusable — the
    domain rule leaves authors with no sanctioned way to
    express an absence-conditional row, and the only
    remaining outlets are the sentinel-stamping anti-pattern
    or a grammar extension from RDR 0003. The domain rule
    itself survives; its authoring story does not.
- **A13 Provenance-blind guard presence is safe against
  caller-supplied input, or the exposure is accepted
  knowingly: an observed tag supplied by the caller can turn
  a `guard_unevaluable` into a decided guard verdict, which
  no shipped test contends.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: Needed. Established so far: `Input.Observed`
    is caller-supplied (`resolve.go::Input`; `ProvenanceObserved`
    is documented as "a non-owned context tag supplied by the"
    caller) and `assemble` writes it into `view.tags`
    unconditionally, so presence under the provenance-blind
    rule is caller-reachable. The kernel already refuses this
    substitution on the *owned-state* path — frozen by
    `adversarial_test.go::TestAdv4_ObservedTagsMustNotShadowTheOwnedSnapshot`
    and `::TestAdv5_ObservedTagCannotSatisfyAnOwnedStateRequirement`
    — but neither test covers the guard path, and before this
    RDR nothing stated that guards read the view
    provenance-blind. A11 verifies the rule *conforms* to
    shipped `TagSet.Lookup` behavior; it does not establish
    that conformance is safe. Verify whether the
    accessor/input boundary (RDR 0004 / RDR 0005) constrains
    what a caller may supply as observed, and record the
    residue here if it does not.
  - **If wrong**: an operator can unstick a
    `guard_unevaluable` refusal by supplying the missing tag
    on the command line, producing a plan whose `Writes` are
    computed from caller input rather than artifact state —
    the ADV-4/ADV-5 defect reintroduced one seam over, with
    no test contending it.

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
`escapeOrRefuse` applies the same gate to the escape rows
as a separate row set (per D5 scoping — see the aggregation
contract below).
This RDR adds (a) the normative domain rule below, (b)
doc-contract narrowing on `Row.RequiresOwned`, and (c) a
conformance fixture set the 0003 evaluator must pass.

#### Normative Contracts

```normative
SCOPE OF THE ATOM VOCABULARY. The clauses below quantify
over guard ATOMS — operator, referenced tag key, literal —
and over the `all`/`unless` blocks that combine them. That
structure is RDR 0003's; the kernel does not model it. At
the kernel seam a guard is the opaque `Row.Guard string`
handed to `GuardEvaluator.Evaluate`, and the kernel
attaches no meaning to it (`resolve.go::GuardEvaluator`:
"The kernel calls it; it does not implement operator
semantics").

This RDR therefore binds the EVALUATOR, not the kernel: the
implementation of `Evaluate` MUST recover the atom
structure of the guard it is given and apply these rules to
it. How guard structure survives the trip from the
authored table (RDR 0003's spike fixture authors it as
`[rule.guard.all.<tag>]` / `[rule.guard.unless.<tag>]`)
into `Row.Guard` — a parsed representation, or a key into
one the evaluator holds — is RDR 0003's to state, and it
MUST be stated before that evaluator is built. This RDR
requires only that the mapping be total and lossless for
operator, referenced tag key, literal, and block placement:
an evaluator that cannot recover which tag keys an atom
references cannot implement the domain rule at all.
```

```normative
Value-comparing guard operators (equality, membership,
bounded integer comparison, set containment) are PARTIAL
over the assembled evaluation view: an atom whose
referenced tag key is absent from the view MUST evaluate to
unevaluable — never to false and never to true.

This rule is fixed HERE. RDR 0003 declares the closed
operator vocabulary and names absence only for `exists`; it
states no absent-operand behavior for the other four. In
particular, an absent set-valued tag MUST be treated as
unevaluable under set containment, NOT as the empty set —
absence is not an empty value.
```

```normative
The existence operator is the sole TOTAL operator: it MUST
decide true or false from the presence or absence of its
referenced tag key alone. Absence tests MUST be expressed
as explicit existence atoms; no value-comparing operator
may act as an implicit existence test.
```

```normative
PRESENCE IS PROVENANCE-BLIND. "Present in the assembled
view" means the tag key is in the view under ANY
provenance — owned, observed, or recognized. A guard atom
MUST NOT be decided differently according to how a tag
reached the view, and `exists` MUST decide TRUE for an
observed- or recognized-supplied key exactly as for an
owned one. (RDR 0003's own guard fixture puts `exists` on
`cluster_eligible`, declared `provenance = "observed"`, so
an owned-only reading would break a spec'd example.)

This is deliberately NOT the kernel's owned-state
predicate: `resolve.go::missingOwned` tests
`view.has(key, ProvenanceOwned)` because
`owned_state_unavailable` is about the owned snapshot
specifically. Guard decidability is the broader question
and takes the broader test — `TagSet.Lookup`'s
provenance-blind `ok`. The two predicates answer different
questions and MUST NOT be conflated.
```

```normative
The domain rule governs guard ATOMS. A row carrying no
guard has no atoms and is therefore outside the rule's
scope: it is decided TRUE without consulting the evaluation
view, as the kernel already does
(`resolve.go::evaluateGuard`, empty guard yields GuardTrue
before the seam is reached). An unguarded row reads no tags,
so it cannot mask missing artifact state.

Empty atom blocks take the conventional identities: an empty
`all` block is TRUE (empty conjunction), and an empty
`unless` block MUST be treated as absent — NOT as a
vacuously-true conjunction, which under
`all ∧ ¬(unless_conj)` would disable every row carrying one.
```

**Open ownership question (A9, Pending) — not part of the
normative block above.** If RDR 0002's normalization already
fixes the empty-`unless` identity, this clause is 0002's and
is struck here rather than restated; if it does not, the
clause stands as written. A normative MUST cannot carry its
own strike-out condition, so the question is recorded here
and A9 blocks lock until it is answered either way. The
identity itself is not in doubt — only which RDR owns it.

```normative
Atom verdicts combine under strong-Kleene three-valued
logic across `all` and `unless`. A guard verdict MUST be
GuardTrue or GuardFalse only when the present tags alone
decide it; if the combined verdict has any *unresolved*
dependence on an unevaluable atom — one the present tags do
not already decide — the evaluator MUST return
GuardUnevaluable. Absence MUST NOT contribute truth or
falsity to any value-comparing atom.

A present-tag atom decided FALSE therefore still yields
GuardFalse for the whole `all` block even beside an
unevaluable atom (strong-Kleene `F ∧ U = F`): the falsity
is witnessed by present state and holds under every
resolution of the unevaluable atom. This is what makes D8
sound, and it is conformance to strong Kleene, not a
deviation from it.
```

```normative
The kernel maps GuardUnevaluable to the refusal kind
`guard_unevaluable`, which is not a modelable escape class
(RDR 0002). Missing artifact state MUST NOT be maskable
behind an escapable refusal class.

Among surviving rows, absent owned state is reported BEFORE
an undecidable guard: a survivor set carrying both MUST
refuse `owned_state_unavailable`. This is the shipped
ordering (`resolve.go::gate` runs `missingOwned` before the
undecidable loop) and is pinned here because no test
currently contends the two within one survivor set.

The ordering is pinned to shipped behavior, NOT to its
original rationale, which this RDR's narrowing invalidates.
D8 justified the precedence as "absent owned state is the
more precise diagnosis and is frequently the reason the seam
could not decide the predicate"
(`0001-…/artifacts/deviations.md`, echoed in
`resolve.go::gate`'s doc comment) — true only while
`RequiresOwned` named guard *inputs*. Once it names
post-guard *write* dependencies, a missing RequiresOwned key
is no longer a plausible cause of the guard's
undecidability, so "more precise" no longer follows: the two
refusals now diagnose independent problems. The precedence
is retained anyway because it is frozen kernel behavior and
reversing it is a kernel change this RDR declines to make
(A3) — but the honest cost is that a row failing both ways
surfaces the write-dependency problem first and the
decidability problem only on the next run. Phase 2's pinning
vector MUST assert the ordering as behavior, and the
authoring docs MUST NOT repeat the superseded rationale.

The refusal MUST identify the rows it came from:
`Refusal.Rows` carries every undecidable row.
`Refusal.Guard` is single-valued — the kernel selects the
lowest row by `(RuleID, SourceLocator)` — so it is a
diagnostic convenience, NOT the discriminator: any contract
or test distinguishing *which* row went unevaluable MUST
assert on `Rows`.
```

```normative
Row.RequiresOwned names the owned tag keys the row's
post-guard transition depends on — the keys its `Writes`
require, which per RDR 0002 includes an authored clear
(normalization "renders a `<clear>` write", so the kernel
`Row` carries no separate clear list). Guard decidability is
not RequiresOwned's job: a guard's input coverage is
enforced by the domain rule above. Listing a guard-read
owned key in RequiresOwned remains legal and yields the
more precise owned_state_unavailable diagnosis among
surviving rows.
```

```normative
Aggregation is resolution-level, not merely per-guard: if
any surviving candidate row's guard is GuardUnevaluable,
the resolution MUST refuse guard_unevaluable — a
decided-GuardTrue sibling row MUST NOT be selected while
an unevaluable candidate exists. Only decided-GuardFalse
rows are pruned from consideration.

Aggregation is scoped to one row set at a time (RDR 0001
deviation D5): the candidate rows aggregate among
themselves, and the escape rows are gated as their own
set. An unevaluable CANDIDATE row MUST NOT be masked by a
decidable escape row. An unevaluable ESCAPE row yields
guard_unevaluable in place of the candidate-set refusal:
the kernel MUST NOT claim the escape failed to rescue when
it could not decide the escape at all.

This veto is intended, and it is the cost the RDR accepts:
one unreadable row refuses a table whose other rows decide
cleanly. It does not contradict D5's anti-poisoning
rationale, which scopes evaluation to *matching* candidates
so an UNRELATED row's obligation cannot poison a legal
resolution. An unevaluable row that matched the outcome and
the tag pattern is not unrelated — it is a live edge whose
applicability is unknown, and selecting a sibling past it
would assert exactly the "I could tell" the Problem
Statement forbids. Narrowing the veto would require
deciding that an undecided edge is a non-edge, which is
absence-as-false at resolution scope.
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
| Refusal mapping + ordering | `internal/resolve/resolve.go::gate` | Owned-before-unevaluable precedence is unfrozen (no contending test) | Reuse | ADV-1/1b freeze D8; A3 verifies no change needed; Phase 2 pins the precedence |
| Owned-state precheck | `internal/resolve/resolve.go::missingOwned` | Owned-provenance keys only; covers declared keys, not guard-read ones | Reuse | RequiresOwned narrowed in doc contract only |
| Verdict-only stub evaluator for the MVV | `internal/resolve/fixtures_test.go::fixtureGuards` | `Evaluate(guard string, _ resolve.TagSet)` **discards the view**; verdicts come from a `decided map[string]bool` keyed on guard text | Reuse (rows 1–3 only) | Sufficient where a verdict is handed in; cannot exercise presence/absence, so it backs MVV rows 1–3 and no more |
| View-reading evaluator for the domain-rule scenarios | None — no `GuardEvaluator` implementation reads the `TagSet` | Scenarios 4/6/7/8 assert operator × presence/absence behavior the stub cannot express | Build | Phase 2 must build it alongside the vectors; it is not free reuse |
| Golden conformance vector harness | None found under `internal/resolve/` (no `testdata/` repo-wide) | — | Build | Phase 2 introduces it; with the evaluator above, this RDR's new test-tree surface |

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
| Operability (recovery when a guard cannot read state) | ~ — masked, so nothing to recover *from*; the failure is silent | ✗ — worst of the three: `guard_unevaluable` is non-escapable, so there is no table-level hatch; recovery is fixing the artifact or accessor, and while A6b is open a transient read failure is indistinguishable from real absence | ~ — precheck refusal names the absent key directly, but is equally non-escapable |
| Blast radius on the authoring surface | None — closed-world reads keep working | ✗ — absence-conditional rows need the two-row pattern, whose load-time legality is open (A12); the project's own fixtures become migration candidates (Consequences) | ~ — authors must maintain a declared-reads list |
| Cost | Low now, high at first masked escape | Low — the discriminator already exists end to end | Medium — new Row field or overloaded RequiresOwned (re-creates xg7p's conflation) |

The deciding rows are correctness fit and blast radius: B
is the only approach that closes the masking path without
touching frozen kernel behavior, and the sibling-path check
shows its discriminator already shipped
(`GuardResult`/`KindGuardUnevaluable`), so B reuses an
existing signal where A and C would bend or duplicate one.

B loses the operability and authoring-blast-radius rows
outright, and is chosen anyway. That trade is accepted, not
overlooked: the alternative to "no recovery hatch" is
precisely the escapable masking A rejects, so an escape
hatch here would reintroduce the defect the RDR exists to
close. The cost is bounded to tables that read state they
cannot see, and the honest refusal is what makes the
condition visible enough to fix.

The narrower variant the operability row invites — degrade
an unevaluable row to an escapable class when a decided
sibling exists — was considered and rejected under
*Briefly Rejected* below: it is absence-as-false at
resolution scope, since it decides that an undecided edge is
a non-edge. What the matrix cannot settle is how *often* the
veto fires on tables whose rows are unrelated in practice;
that is an empirical question with no installed base to
answer it (Consequences), so it is recorded as accepted risk
rather than scored away.
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

- Is exactly the masking path the Problem Statement
  describes, with nothing distinguishing "did not apply"
  from "could not tell".
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
- **Narrowed veto (an unevaluable row degrades to an
  escapable class when a decided-TRUE sibling exists)**:
  buys back the operability the aggregation veto costs, and
  is rejected because it *is* absence-as-false at resolution
  scope — selecting the sibling asserts the unevaluable row
  does not apply, which is exactly the "I could tell" the
  Problem Statement forbids, and routes it through an
  escapable class, which is the masking path Alternative 1
  was rejected for. The variant that scopes the veto to rows
  whose referenced tags intersect the selected row's
  `Writes` fails differently: the evaluator cannot compute
  that intersection without the atom structure A10 shows it
  does not have.

## Trade-offs

### Consequences

- Positive: missing artifact state becomes unmaskable
  behind an escape — the user learns "I could not tell" as a
  distinct, non-escapable refusal. This is a property of the
  *contract*, realized when RDR 0003's evaluator implements
  it; this RDR ships the rule, the doc narrowing, and the
  conformance vectors, not a behavior change (Phase 1 is
  doc-comments-only, and the kernel consumer side already
  exists — A3). The masking path recorded under *Background*
  stays open until that evaluator is built and runs the
  Phase 2 harness. **Consequence for status:** this RDR's
  own probe still reproduces on the shipped kernel after
  Phases 1–3 land, because no code path changes. The
  tracking kata (`xg7p`) is resolved by the contract, not by
  the behavior, and it MUST be closed against that
  distinction explicitly rather than against a green test
  run — otherwise the same finding is re-derived later
  against a kernel that never changed. The behavior closes
  when RDR 0003's evaluator lands and passes the Phase 2
  harness; until then the honest status of the *defect* is
  open even when the status of the *RDR* is Implemented.
- Positive: no kernel code change; D8/D5 and the ADV suite
  stand as shipped; kata `xg7p`'s conflation resolves by
  doc-contract narrowing.
- Positive: absence remains expressible — authors state it
  intentionally via existence atoms instead of relying on
  implicit closed-world reads.
- Negative: tables whose authors expected closed-world
  reads ("x == 1 just fails when x is missing") now refuse
  where they previously escaped or pruned; those guards
  must be rewritten with explicit existence atoms. The
  exposure today is bounded — no `GuardEvaluator`
  implementation exists yet, so no shipped table depends on
  closed-world reads and there is no *executing* installed
  base to migrate. The rule therefore lands *before* the
  first evaluator rather than changing behavior under
  existing authors, which is why this RDR must be Final
  before RDR 0003 is implemented (Prerequisites). If that
  sequencing slips and an evaluator ships first, this bullet
  becomes a real migration with no inventory tooling behind
  it.
  **Authored guards already exist, and they are the first
  migration case.** "No installed base" counts
  implementations, not authored intent:
  `0003-…/evidence/spikes/guard-fixture.toml` authors guards
  over `cluster_eligible` (`provenance = "observed"`) and
  `prelock_iterations` (`int`, `min = 0 / max = 3`), and
  `0002-…/evidence/spikes/rdr-fixture.toml` authors guards
  over tags declared with no domain at all. Under the domain
  rule, every value-comparing atom in those fixtures becomes
  `guard_unevaluable` the moment its tag is absent from a
  view. Phase 2 MUST classify each authored guard in the
  tree as safe-or-migration when it encodes the vectors —
  the reference fixtures are the smallest honest inventory,
  and reclassifying them is cheap now and expensive after an
  evaluator exists.
- Negative: `guard_unevaluable` is non-escapable by
  design, so an artifact in genuinely degraded state has no
  modeled-escape hatch through a guard that cannot read it;
  recovery is fixing the artifact state or the accessor,
  not the table.

### Risks and Mitigations

- **Risk**: RDR 0003's evaluator is implemented later by a
  session that drifts from this rule (e.g. folds absence
  into false for one operator).
  **Mitigation**: partial, and the residue is real. The
  Phase 2 conformance vectors are normative and live in the
  kernel's test tree, but sequencing is all the Prerequisites
  bind: "0007 Final before 0003 implement" does not make
  0003's build *run* the vectors. Because the vectors are
  parameterized over a `GuardEvaluator` and 0003's evaluator
  is a separate future package, `go test ./internal/resolve/`
  stays green against the stub no matter what 0003 builds —
  the same blind spot as premortem P-10. Phase 2 therefore
  owes an **executable** binding, not a documentary one: the
  vector suite MUST be exported as a reusable conformance
  harness (an exported function taking a `GuardEvaluator`),
  and RDR 0003's implement stage MUST instantiate it against
  its own evaluator. That obligation is recorded in Phase 2
  and handed to RDR 0003 in Phase 3.
  **Residual status: UNMITIGATED, not partially mitigated.**
  Nothing in this repo fails if RDR 0003 never instantiates
  the harness — that is the definition of an unenforced
  obligation, and it is the same hole as the risk it answers
  (P-10). Building the harness is necessary and not
  sufficient; the enforcement point is RDR 0003's implement
  stage accepting it, which this RDR cannot compel. Because
  RDR 0003 is Final and this project does not amend RDRs,
  the acceptance must happen at 0003's *implement* stage
  (which reads its Prerequisites), not by editing 0003's
  text — and if that stage declines, the drift risk is
  carried openly rather than recorded as closed.
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
  **Mitigation**: partial. A6a verifies against
  `resolve.go::assemble` that the *kernel* never drops a
  key present in the input tuple, so this risk cannot
  originate in assembly. A6b remains OPEN: RDR 0004 states
  no read-completeness rule, so a truncated accessor
  snapshot can still reach the kernel indistinguishable
  from genuine absence. Until that seam is stated, the
  recovery path for an unevaluable refusal is operational
  (fix the artifact state or the accessor), never a table
  escape, and the authoring docs must say so.
- **Risk**: Authors defeat partiality by stamping sentinel
  values upstream so keys are never absent, making value
  comparisons decide on placeholders.
  **Mitigation**: conditional on A12. A5's two-row pattern
  is the intended legitimate outlet and authoring guidance
  flags sentinel-stamping as the anti-pattern (premortem
  P-6) — but the pattern's load-time legality is open (A12),
  and sentinel-stamping is what authors fall back to if it is
  rejected. So this mitigation stands or fails with A12: an
  anti-pattern warning with no working alternative is not a
  mitigation, it is a prohibition, and it would make the
  sentinel path the only route to an absence-conditional row.

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
  change. Note the asymmetry this leaves against the
  Problem Statement's promise of being "told plainly that
  the artifact state needed to decide was missing": the
  sibling `owned_state_unavailable` names its keys in
  `MissingOwned`, while `guard_unevaluable` carries
  `Refusal.Guard` (one row's guard text) and
  `Refusal.Rows`. The refusal is therefore *actionable* —
  it names the row and the predicate, and it is
  non-escapable — but it does not name the absent tag
  directly. Closing that last step is the Phase 3 handoff
  above; it is a diagnostic enrichment, not a change to
  the domain rule, and it does not reopen the masking path.
- **Review-time blind spot guarded against**: a typo'd or
  undeclared guard key would be permanently absent —
  forever unevaluable, or forever shadowed when
  FALSE-domination decides the guard from its other atoms.
  The catch is RDR 0003's declared-tag rule at parse/lint
  (A7): the defect fails table load, never reaching a
  production refusal.
- **Refusal defeated by caller-supplied state** (A13, open):
  because presence is provenance-blind and `Input.Observed`
  is caller-supplied, an operator facing a
  `guard_unevaluable` can supply the missing key as an
  observed tag and convert the refusal into a decided
  verdict — producing a plan whose `Writes` derive from
  caller input rather than artifact state, with nothing in
  the output marking the difference. The kernel already
  refuses this substitution on the owned-state path
  (`adversarial_test.go::TestAdv4_ObservedTagsMustNotShadowTheOwnedSnapshot`,
  `::TestAdv5_ObservedTagCannotSatisfyAnOwnedStateRequirement`);
  no equivalent test contends the guard path, and this RDR
  is the first document to state that guards read the view
  provenance-blind. The provenance-blind rule itself is not
  in doubt — it is forced by `TagSet.Lookup` being the only
  exported accessor and by RDR 0003's own fixture guarding
  an observed tag (A11) — but "the refusal is
  non-escapable" is a weaker guarantee than it reads: it is
  unroutable *around within the table*, not unbypassable by
  the caller. A Phase 2 vector should pin whatever A13
  settles.
- **Conflated recovery signal**: `guard_unevaluable`
  cannot itself distinguish transient read failure
  (retryable) from genuine absence (not retryable); until
  A6b settles what the accessor contract reports,
  operators should treat the refusal as "inspect the
  artifact and accessor," not "retry". A6b is confirmed
  open: RDR 0004 requires only "typed tag values or a
  typed refusal", with no completeness rule on the values
  branch, so the conflation is real and currently
  unowned.

## Implementation Plan

### Prerequisites

- [ ] **A9 verified** (empty-`unless` identity against RDR
      0002's normalization) — raised by the cove lens, which
      found both this RDR and RDR 0003 silent on empty atom
      blocks. Must close before lock: the wrong identity
      disables every row carrying an omitted `unless` block.
- [ ] **A10 verified** (guard-structure mapping at the
      `Row.Guard` boundary) — raised by the 3amigo lens as
      its hinge finding. Must close before lock: every
      atom-level normative clause here binds an evaluator
      that must recover atom structure from an opaque
      string, and the Phase 2 vectors for scenarios 4/6/7/8
      cannot be encoded until the mapping is known. If the
      mapping is unspecified anywhere, the resolution is to
      hand it to RDR 0003 as a named obligation, not to
      specify a grammar here — **and the handoff needs a
      destination that exists.** RDR 0003 is Final and this
      project does not amend RDRs, so "state it in 0003"
      cannot mean editing 0003's text. It means 0003's
      *implement* stage owns the mapping as a named
      prerequisite, or a successor RDR states it. Whichever
      it is MUST be named when A10 closes; an obligation
      addressed to a locked document is not a resolution,
      and leaving it there is what turns Phase 2 into
      indefinite deferral.
- [ ] **A12 verified** (two-row absence pattern survives
      RDR 0003's load-time overlap check) — raised by the
      critique lens. Must close before lock: A5's blessed
      pattern is the domain rule's only sanctioned outlet for
      absence-conditional rows, and `ambiguous overlap` is a
      mandatory rejection category, not a downgradeable
      claim. If the pattern is overlap-rejected, the
      authoring story needs a different answer before this
      rule binds authors.
- [ ] **A13 verified** (provenance-blind presence is safe
      against caller-supplied observed tags, or the exposure
      is accepted knowingly) — raised by the critique lens.
      Not a blocker for the domain rule itself; a blocker for
      claiming the refusal cannot be worked around.
- [x] All other Critical Assumptions verified **except
      A6b**, which is carried open by decision: the
      read-failed vs genuinely-absent distinction is not
      stated by RDR 0004 and cannot be resolved inside this
      RDR. It is
      not MVV-critical (the kernel seam cannot express the
      distinction either way), and the contract is sound
      without it; the exposure is recorded under Risks and
      Failure Modes.
- [ ] Read-completeness obligation routed to RDR 0004 (or
      the accessor executor's owning RDR) — A6b. Not a
      blocker for this RDR's lock; a blocker for trusting
      negative existence atoms in production.
- [ ] This RDR Final **before** RDR 0003's implementation
      begins (triage sequencing: D8 ratification precedes
      the evaluator build)

### Minimum Viable Validation

Kernel-level tests (ADV-style, with a stub evaluator
conforming to the domain rule) proving the masking probe
inverts: a guard referencing an absent tag, plus an absent
`RequiresOwned` key, plus a modeled `no_match` escape,
yields a refusal of kind `guard_unevaluable` — while the
same table with the tag present and decided FALSE still
prunes and escapes (D8 preserved). Two named scenarios from
the premortem are in scope: *unevaluable-not-no_match* (one
row, guard over an absent key → `guard_unevaluable`, never
`no_match`) and *unevaluable-blocks-true-sibling* (row A
unevaluable beside row B decided true → refusal, no plan).
Mixed-verdict combination cases (decided-FALSE beside
unevaluable, in `all` and in `unless`) assert the
strong-Kleene selection.

**What the MVV does and does not prove.** Scenarios 1–2 hand
the kernel a verdict and assert its already-frozen mapping,
so they pass against the tree as it stands today — they are
regression pins on the consumer side of the contract (A3's
"no kernel change" claim), not evidence that the domain rule
works. Scenario 3 (*unevaluable-blocks-true-sibling*) is the
only MVV row pinning behavior no shipped test covers. The
clauses that carry this RDR's actual content — the atom
domain rule, `exists` totality, `contains` over an absent
set, strong-Kleene combination — are Testing Strategy rows
4–8, which cannot run until A10 fixes what a vector's input
is. A "green MVV" therefore MUST NOT be read as validating
the domain rule; the honest gate for that is the Phase 2
harness.

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

Two things this phase must build, neither of them reuse:

1. **A view-reading evaluator.** `fixtureGuards` discards
   its `TagSet`, so it cannot express presence/absence;
   scenarios 4/6/7/8 need an evaluator that actually reads
   the view. This is the domain rule's first executable
   expression.
2. **An exported conformance harness**, not an internal
   test. The suite MUST be callable against any
   `GuardEvaluator` (an exported function taking the seam),
   so RDR 0003's implement stage can instantiate it against
   its own evaluator. A suite that only runs against the
   in-tree stub proves nothing about 0003's build and would
   leave P-10 unmitigated.

   Two constraints Phase 2 must satisfy, both currently
   unmet by the tree:
   - **Importability.** The harness MUST live in a
     non-`_test` package so an external caller can import
     it — `fixtures_test.go::fixtureGuards` is in package
     `resolve_test` and is unimportable by construction, so
     it is a model for the stub's *behavior* only, never a
     reuse target for the harness. `GuardEvaluator` lives
     under `internal/`, which is importable anywhere inside
     `github.com/newcoinc/intrastate`; if RDR 0003's
     evaluator lands outside this module, the harness
     placement must be revisited before Phase 2 encodes it.
   - **Signature.** The harness is an exported function
     taking the seam and a `*testing.T`-like reporter, e.g.
     `resolveconform.RunGuardVectors(t, eval)`; the exact
     package name and vector encoding are Phase 2's to fix,
     but the shape MUST be "caller supplies the evaluator",
     not "suite constructs one".

### Phase 3: Diagnostics and authoring handoff

Record with RDR 0003's surfaces: the **guard-structure
mapping** obligation (A10 — what `Row.Guard` carries and how
the evaluator recovers operator, tag key, literal, and
`all`/`unless` placement from it; this RDR's atom-level
clauses are unimplementable until 0003 states it), the "name
the absent tag" diagnostic obligation (its A4 semantic
kinds), confirmation that the declared-tag rule is a
checkable vocabulary lint (A7), the authoring guidance —
blessed two-row absence pattern, sentinel-stamping
anti-pattern (A5) — and the **conformance obligation**: RDR
0003's implement stage must instantiate Phase 2's exported
harness against its own evaluator, which is the only thing
that actually catches drift (Risks). Together these make
unevaluable refusals fully actionable without a kernel
change.

## Validation

### Testing Strategy

The matrix the verified assumptions imply. Rows 1–3 are the
MVV and run against `internal/resolve/fixtures_test.go`'s
`fixtureGuards`, which answers unknown predicates
`GuardUnevaluable` — sufficient because rows 1–3 only need a
verdict handed to the kernel, not a view read.

Rows 4–8 cannot reuse it: `fixtureGuards.Evaluate` discards
its `TagSet` argument and decides on guard *text* alone, so
it cannot express presence/absence at all. Those rows assert
operator × present/absent behavior, which requires a
**view-reading evaluator that does not yet exist**. Phase 2
builds it together with the vectors; it is the domain rule's
first executable expression, not a reused stub. Rows 4–8 are
the Phase 2 golden conformance vectors the RDR 0003
evaluator must pass verbatim. "Done" = every row green, and
rows 1–2 fail loudly if any evaluator folds absence into
`GuardFalse`.

**Prerequisite for rows 4–8 (blocking).** These rows name
operators, atoms, and `all`/`unless` placement, and no Go
representation of any of that exists: `Row.Guard` is an
opaque `string`, there is no atom type, and no `testdata/`
exists repo-wide. The guard-structure mapping the Normative
Contracts' *Scope of the atom vocabulary* clause requires
(A10) MUST be settled before these vectors can be encoded —
it fixes what a vector's input literally *is*. Rows 1–3 are
unaffected: they hand the kernel a verdict and need no atom
structure.

Scenarios 3 and 5 are backed by a Resolve spike that ran
against the shipped kernel; its captured output is the
normative expected value
(`evidence/spikes/aggregation-probe.md`).

1. **Scenario**: *unevaluable-not-no_match* — one candidate
   row whose guard references a tag absent from the view,
   plus an absent `RequiresOwned` key, plus a modeled
   `no_match` escape row (the masking probe from the
   Problem Statement).
   **Expected**: `Result.Refused()` is true with
   `Refusal.Kind == guard_unevaluable`, and `Result.Plan` is
   nil — never `no_match`, never a plan. (`Escaped` is a
   `Plan` field and is unreadable on a refusal; the
   discriminator is `Refusal.Kind` plus a nil `Plan`.)
   `Refusal.Rows` names the unevaluable row. Inverts the
   observed pre-RDR behavior recorded under *Background*.
   Backing: `resolve.go::gate` undecidable branch; A3.

2. **Scenario**: D8 preserved — the same table with the
   guard's tag *present* and the predicate decided FALSE.
   **Expected**: the row prunes and the modeled escape
   rescues. Guards against a fix to scenario 1 that
   over-refuses.
   `adversarial_test.go::TestAdv1b_AllGuardsFalseIsNoMatchNotOwnedStateUnavailable`
   freezes the prune-and-rescue *disposition*, but it is
   not the A/B control this scenario needs: it decides
   `"never"` FALSE from guard text with the `TagSet`
   discarded, and never varies tag presence. The paired
   present/absent run is therefore **new work in Phase 2**
   against the view-reading evaluator, not a reuse of
   ADV-1b.
   Backing: ADV-1/ADV-1b (disposition only); A2 (`F ∧ U = F`
   is witnessed falsity).

3. **Scenario**: *unevaluable-blocks-true-sibling* — row A
   unevaluable beside row B decided TRUE, both matching.
   **Expected**: `guard_unevaluable` naming row A's guard;
   no plan. Row B MUST NOT be selected.
   Backing: `resolve.go::gate` discards `selected` on
   `len(undecidable) > 0`; spike Probe B observed
   `kind="guard_unevaluable" guard="unknown-predicate"`.
   **No shipped test covers this** — `TestAdv3` pairs two
   unevaluable rows (order-stability only); this row
   closes a real gap.

4. **Scenario**: strong-Kleene combination matrix — each
   operator × {tag present, tag absent} × verdict pairs
   across `all`, and the same across `unless` (recalling
   `unless` is ¬(conjunction), not per-atom negation).
   **Expected**: the A2 tables. Specifically `F ∧ U = F`
   (prunes), `T ∧ U = U`, `U ∧ U = U`, and `unless`
   yielding FALSE only when every `unless` atom is TRUE.
   Backing: A2 derivation; A4 (PostgreSQL 16 §9.1).

5. **Scenario**: escape-set scoping — (a) an unevaluable
   escape row over a `no_match` candidate set; (b) an
   unevaluable candidate row beside a *decidable* escape
   row.
   **Expected**: both legs refuse `guard_unevaluable`; the
   discriminator is `Refusal.Rows`, which MUST name the
   escape row in (a) and the *candidate* row in (b) —
   `Refusal.Guard` is single-valued and cannot separate the
   legs when both guards carry the same text. (a) is frozen
   by `adversarial_test.go::TestAdv2` subtest "guard
   UNEVALUABLE must not rescue" and
   `fixup_test.go::TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue`;
   spike Probe A observed it. In (b) the escape path is
   never reached, so give the two rows distinct `RuleID`s
   and assert on `Rows`.
   Backing: `resolve.go::escapeOrRefuse`; A3.

6. **Scenario**: `exists` is total — a positive existence
   atom over an absent key, and a negative existence atom
   (`unless … exists`) over an absent key.
   **Expected**: both decide (FALSE and TRUE respectively);
   neither returns unevaluable. This is the one sanctioned
   route from absence to a decided verdict.
   Backing: A1; RDR 0003's `exists` row ("Tests presence
   or absence, not value equality").

7. **Scenario**: `contains` over an absent set-valued tag.
   **Expected**: unevaluable — NOT false-by-empty-set.
   Pins the one operator whose set-theoretic reading could
   otherwise decide from absence.
   Backing: A1's absent-operand rule, which this RDR fixes
   because RDR 0003 is silent on it.

8. **Scenario**: the blessed two-row absence pattern — row
   1 guarded `unless … exists` on tag X, row 2 guarded
   `X eq v`, evaluated with X absent and then with X = v.
   **Expected**: exactly one row selected in each case;
   never both (they are disjoint by construction), never
   `ambiguous_match`.
   Backing: A5.
   **Load-time precondition (A12).** This scenario asserts
   *runtime* disjointness and presumes the pair loads at all.
   If RDR 0003's overlap proof rejects the pair over a domain
   product with no absent member, the vector is unreachable —
   the table never reaches resolution. A12 MUST close first;
   if it resolves against the pattern, this scenario changes
   with whatever replaces it.

9. **Scenario**: unguarded and empty-block rows — (a) a row
   with no guard at all, evaluated against a view missing
   every tag it would otherwise read; (b) a row whose
   `unless` block is omitted.
   **Expected**: (a) the row is decided TRUE and selectable
   without consulting the view, matching
   `resolve.go::evaluateGuard`'s empty-guard branch; (b) the
   row is *not* disabled — an omitted `unless` block must
   not read as a vacuously-true conjunction. Guards the
   whole-table failure A9 names.
   Backing: A9; `resolve.go::evaluateGuard`.
   **Test level**: (a) is a kernel test. (b) is *not*
   assertable at the kernel seam — a row with an omitted
   `unless` and no `all` reaches the kernel as
   `Guard == ""`, which `evaluateGuard` answers `GuardTrue`
   before the seam is consulted, so the wrong identity is
   invisible there. (b) belongs to the evaluator's own
   vector set (Phase 2) and presumes A10's mapping; if A9
   resolves to RDR 0002 owning the identity, (b) moves there
   with it.

**Owned-state ordering gap (not MVV, flagged for Phase 2).**
No shipped test contends missing owned state *and* an
unevaluable guard within one survivor set, so the
"owned-state reported first" precedence this RDR's
`RequiresOwned` contract relies on is held by
implementation only. A Phase 2 vector should pin it.

**Out of scope here.** A6b's read-failed vs genuinely-absent
distinction is untestable at this seam: `Input.Owned` has no
error channel, so a truncated snapshot and a real absence
are the same input. Testing it requires the accessor
contract A6b leaves open.

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
- PostgreSQL 16 Documentation, Ch. 9 §9.1 "Logical
  Operators", p. 231 — SQL three-valued logic truth tables
  (strong-Kleene; A4). Corpus `DevRef`; quoted in
  `docs/rdr/0007-guard-predicate-totality/evidence/research/resolve-citations.md`
- RDR 0004 — Accessor execution safety model (read
  accessor contract; the read-completeness seam A6b leaves
  open)
- `internal/resolve/resolve.go` (`GuardResult`,
  `GuardEvaluator`, `gate`, `evaluateGuard`, `assemble`,
  `missingOwned`, `escapeOrRefuse`)
- `internal/resolve/adversarial_test.go` (ADV-1, ADV-1b;
  ADV-2 "guard UNEVALUABLE must not rescue")
- `internal/resolve/fixup_test.go`
  (`TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue`,
  `TestFixupGateIsUniformAcrossOrdinaryAndEscapeCandidates`)
- `internal/resolve/fixtures_test.go` (`fixtureGuards` —
  the conforming stub the MVV reuses)
- Stage 4 aggregation spike:
  `docs/rdr/0007-guard-predicate-totality/evidence/spikes/aggregation-probe.md`
- kata `xg7p` (originating finding)
