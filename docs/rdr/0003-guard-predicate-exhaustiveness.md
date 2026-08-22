# Recommendation 0003: Guard Predicate Exhaustiveness

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-06-19
- **Status**: Draft [revised from Final 2026-08-12 — re-verify complete
  2026-08-22. A5 re-verified against JDR 0001 §D1, RDR 0007's atom-shape clause,
  RDR 0002's source-identity clause, and the shipped `RuleID`/`SourceLocator`;
  all four anchors hold. The `0003-0006-0007` cluster gate (2026-08-22) closed
  JDR 0001 §JD-4 naming this RDR the recording document and settling that lint
  mints no code and gains no non-blocking tier — **A8 is Verified**. §JD-13
  added the single-valued marker to the declaration model here; §JD-14 confirmed
  this RDR's escape-row reading governs, no edit owed. This RDR owns the tag
  declaration model (value kind, finite domain, optionality, single-valuedness,
  element universe), rehomed from RDR 0002 on 2026-08-21; RDR 0002 owns
  authoring location and normalization carriage and cites this model. Remaining
  open: A10 and A12 are tolerances on RDR 0006's refine, and A14 is a
  normalization-carriage request on RDR 0002; A15 is discharged by MVV
  Scenario 3 at implementation.]

- **Type**: Architecture
- **Profile**: large — locks one guard-predicate contract: symbolic atom grammar plus finite-domain exhaustiveness semantics.
- **Priority**: High
- **Related Issues**: None
- **Predecessors**: 0001-resolution-kernel
- **Overrides**: None
- **Seam Lineage**: no prior accretion

## Problem Statement

A flow author needs conditional edges such as cap-3 handling, profile-to-lens routing, and rewind-target legality to be expressed in a way the lint can prove exhaustive and mutually exclusive. The system-internal requirement is to define how tag predicates are represented without giving up static verification.

## Context

### Background

The model treats guards as predicates on tags, not side inputs. The two target flows need equality, integer comparison, and set-membership; the guard shape has to remain strong enough for the lint to prove determinism.

The real design fork is expressiveness versus verifiability: a fixed enum of operators, a tiny expression grammar, or embedded host predicates.

### Technical Environment

intrastate is a Go CLI wired through `internal/cli`. Guard predicates are consumed by the transition table, the resolver kernel, and the lint that proves the graph is safe before runtime.

## Research Findings

### Investigation

RDR 0001's kernel now ships as `internal/resolve`; there is still no transition
table or guard-predicate package under `internal/`, and no guard evaluator. The
kernel calls the guard seam this RDR owns (`resolve.go::GuardEvaluator`) but
supplies no implementation of it. Remaining consumers are the peer RDRs and the
future lint seam. RDR 0001 requires guard evaluation to feed
exact-one edge selection, RDR 0002 names positive `all` and negative `unless`
guard lists but delegates the operator grammar here, and RDR 0006 will depend on
this RDR for static determinism and exhaustiveness checks.

The Domain priors do not include an in-repo competitor document for guard
predicate exhaustiveness, so the bounded prior-art pass used the sibling
`../state-machines` corpus. The strongest local prior is
`MODEL-transition.md`: state is a tag-set, a transition is predicate matching
over that tag-set, and "guard" is not a side input but a predicate over machine
data. That model also names tag provenance (`owned`, `observed`,
`recognized`) and design-time lint failures for overlapping rows,
non-exhaustive predicates, and owned-tag read-before-write.

The tool corpus sharpens the options. `transitions` has the closest authoring
shape to RDR 0002: positive `conditions` and negative `unless` lists. enetx/fsm,
stateless, and qmuntal-stateless prove the runtime pattern, but their guards are
host callbacks, so they are useful implementation priors and poor static-lint
contracts. qmuntal-stateless explicitly requires same-state guards to be
mutually exclusive, which matches this RDR's target property but leaves the proof
to humans or runtime behavior. SCXML/scxmlcc, Sismic, and StateSmith demonstrate
that extended state, guards, bounded counters, rewinds, and conditional skips are
statechart-shaped and can be specified for validation; the `../state-machines`
contrast docs are clear that intrastate should use that class as
specification/verification prior art, not as a runtime orchestrator. Statewright
is useful contrast for a small external guard DSL (`field`, `op`, `value`) and
per-state `max_iterations`, but its MCP harness model is explicitly outside this
project's "no orchestrator" lens.

Arc corpus searches over `StateMachineOS`, `StateMachineLit`, `DevRefOS`, and
`SpecDrivenDev` did not change the option set, but they strengthened the
boundary. `StateMachineOS` surfaced stateless guard-discrimination tests,
`transitions`' `Condition` helper, Sismic guard evaluation, and SCXML `cond`
fixtures. `StateMachineLit` surfaced Symbolic Guardrails and AgentSpec as the
research form of the same split: symbolic predicates are valuable where the
policy is concrete, and an explicit escape remains necessary where judgment is
not symbolically enforceable. StateFlow confirms the finite-state framing and
finite alphabets, but it is orchestration-adjacent rather than a local predicate
grammar for intrastate.

Sibling-path check for an existing guard identity or predicate signal:

```sh
rg -n "resolve|resolver|transition|state|guard|tag|predicate|recognized|outcome|next legal|illegal|refus" internal cmd docs
```

The search found no implemented guard predicate evaluator under `internal/` —
`internal/resolve` declares the `GuardEvaluator` seam and calls it, but no type
satisfies it. The adjacent design signal is RDR 0002's `all`/`unless` split and
exact-one row selection, so this RDR extends that signal instead of inventing a
parallel guard model.

### Key Discoveries

- **Documented** — RDR 0001 requires guard predicates to expose enough
  structure for deterministic single-edge selection.
- **Documented** — RDR 0002 owns sparse transition-table authoring and
  normalization, and delegates the fixed predicate operator set to this RDR.
- **Documented** — prior-art FSM libraries commonly support guards, but
  callback-style guards are only runtime predicates; they do not give lint a
  symbolic domain to prove coverage or overlap.
- **Documented** — `MODEL-transition.md` already frames guards as tag-set
  predicates and requires design-time rejection for overlapping predicates,
  non-exhaustive predicates, and owned-tag read-before-write.
- **Documented** — statechart/SCXML/Sismic prior art covers the hard workflow
  cases in grammar: extended state, bounded counters, rewinds, and conditional
  skips; the project lens keeps these as validators/specs, not as orchestrators.
- **Documented** — arc `StateMachineLit` results support the symbolic/escape
  boundary: concrete policies can be enforced by symbolic predicates, while
  ambiguous requirements still need a model or human escape path.
- **Verified** — the target RDR and kata flows only need equality, bounded
  integer comparison, enum/set membership, existence, and declared negation via
  `unless`; the Resolve spike encoded representative guard rows with that
  vocabulary.
- **Verified** — requiring finite domains for exhaustive checks is acceptable:
  bounded integer tags such as cap counters can be proven, while unbounded
  values can be evaluated but cannot carry an exhaustiveness guarantee.

### Critical Assumptions

- **A1 The target RDR and kata flows fit a closed typed predicate vocabulary.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: The prior `check.sh` did **not** establish this — it is an awk
    scan over operator *spellings* in TOML section headers whose `coverage=` line
    is a hardcoded `printf` constant and whose report loop is hardcoded to
    `split("eq in lt gte exists")`, so it evaluates no predicate and cannot show
    the vocabulary is *sufficient*. Stage 6 re-ran it as the evaluation harness
    A1's plan named. Command:
    `cd docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes/iter-2/a1-eval-harness && go run .`;
    captured output at
    `docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes/iter-2/a1-eval-harness/output.txt`.
    The harness encodes the fixture's four target-flow rules as normalized atoms
    (key, operator, literal, block), adds a `contains` row over a set-valued tag
    and a `gt`/`lte` band row, and **evaluates** them against eleven
    representative tag views, asserting each row's expected qualify/prune
    verdict: 28 assertions pass, 0 fail. All eight operators of the closed set
    carry at least one evaluated atom — including `lte`, `gt`, and `contains`,
    which the old transcript could not print. Five negative controls confirm the
    vocabulary is closed *and typed*: unknown operator, `lt` on an enum,
    `contains` on an enum, literal/kind mismatch, and `in` with a scalar literal
    are each rejected at well-formedness before evaluation. The harness also
    reproduces the seam split it must not violate — `exists` decides on presence
    alone (an absent `cluster_eligible` yields a decided FALSE, never
    unevaluable) while a value atom over an absent key yields UNEVALUABLE under
    RDR 0007's veto. No target-flow row required an operator outside
    `{eq,in,lt,lte,gt,gte,exists,contains}`.
  - **Scope of this record**: it establishes vocabulary *sufficiency and
    decidability* over the target-flow slice — evaluation consumes tag values,
    not declared domains, so it is independent of the declaration model, which
    governs exhaustiveness (A2) rather than evaluation. The `contains` cell is
    evaluated over a declared set **kind**; the element *universe* the
    exhaustiveness claim needs is a separate fact, declared by A11's model and
    recorded by A9.
  - **If wrong**: The fixed operator set is too small, and authors will need an
    expression grammar or host predicates that weaken static lint.
- **A2 Every exhaustiveness claim can be reduced to scoped finite declared
  domains.**
  - **Status**: Verified
  - **Method**: Derivation
  - **Evidence**: The derivation below is sound, and its inputs are now declared
    by this RDR's own tag declaration model (A11). It consumes "finite domains for the guard
    dimensions" — enum value sets, declared set-element universes, and bounded
    integer `{min..max}` ranges — and the `Normative Contracts` declaration
    clauses make each of those declarable, with domain/kind agreement enforced
    before normalization. The circularity this record was previously corrected
    for is closed at its root: the domain is no longer a fixture-invented key
    read back as evidence, but a normative declaration this document owns.
    A7 (presence dimension) and A9 (element universe) are subfields of that
    model and are likewise `Verified`. The derivation: lint receives the candidate rows that share a selection context from the normalized model plus finite domains for the guard dimensions that participate in that group (every key any row's `all` or `unless` carries, whether or not it discriminates — see the participation clause in `Normative Contracts`): enum/boolean values as declared sets, set-valued tags as a symbolic set over the declared element universe, and bounded integers as `{min..max}`. A row with `all` atoms denotes the intersection of each atom's allowed subset of the scoped product; its `unless` block denotes an excluded intersection that is subtracted from the row's accepted assignments. Coverage is `union(row_i accepted assignments) == scoped product` for that row group, and overlap is any non-empty `row_i accepted assignments intersect row_j accepted assignments`. If any participating dimension lacks a finite domain, or the finite product cannot be represented by the implementation's symbolic/bitset-equivalent proof, lint must refuse or downgrade the exhaustiveness claim.
  - **If wrong**: Lint may falsely claim guard coverage or miss legal gaps in
    cap/profile/lens routing.
- **A3 `all` plus `unless` is enough polarity; inline `not` operators are not
  required.**
  - **Status**: Verified
  - **Method**: Design Decision
  - **Evidence**: This RDR chooses separate positive and negative guard lists: `all` is the required conjunctive predicate set, `unless` is the conjunctive exclusion set, and RDR 0002's normalized candidate-row contract combines both before ambiguity checks. Inline `not` and nested boolean expressions are rejected to keep each diagnostic tied to an authored atom and to preserve finite-domain coverage/overlap derivation.
  - **If wrong**: Authors will duplicate rows or encode confusing inverse
    predicates that make overlap diagnostics harder to understand.
- **A4 Guard predicate errors can use the existing structured CLI failure
  gateway.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/cli/clierr/clierr.go::CLIError` carries stable `Code`, human `Message`, optional `Param`, `Detail`, `Hint`, and exit-code `Group`; `internal/cli/respond/respond.go::Fail` emits the envelope in text/json modes; `internal/cli/config/config.go::Load` already demonstrates stable parse/read error codes. This RDR owns predicate semantic kinds such as unknown operator, type mismatch, literal parse failure, and unsupported operator/tag-kind pairing; the per-row/per-atom `guard_unevaluable` payload is RDR 0007's under JDR 0001 §D4; RDR 0006 owns graph-lint finding codes for non-exhaustive and overlapping row groups; RDR 0005 owns the CLI command/envelope mapping onto the existing gateway.
  - **If wrong**: This RDR or RDR 0005 must add a separate user-facing error
    contract before implementation.
- **A5 The normalized predicate representation can retain source identity for
  actionable diagnostics.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: Re-verified 2026-08-22 (the ID this re-entry names); all four
    anchors resolve and none moved under the declaration-model rehoming, which
    touched declaration semantics, not row carriage. The normalized atom shape
    is fixed at the kernel seam by JDR 0001 §D1
    (`docs/jdr/0001-resolve-kernel-seam.md:64`); RDR 0007 is the landing
    document and states it normatively — a row carries parsed atoms "key,
    operator token, literal, block ∈ {all, unless}"
    (`docs/rdr/0007-guard-predicate-totality.md:1252`), not an opaque predicate
    string, so there is no reconstruction step that could lose identity. RDR
    0007 is `Final`, not yet implemented: the shipped kernel still carries
    `Row.Guard string` (`internal/resolve/resolve.go:185`), and 0007's A3 scopes
    the reshape. The source identity this assumption is about already ships —
    `internal/resolve/resolve.go:173-174` carries `RuleID` and `SourceLocator`,
    commented as "the source identity RDR 0002 requires". RDR 0002 `Normative
    Contracts` require each normalized candidate row to "retain its source rule
    id and source locator"
    (`docs/rdr/0002-transition-table-as-reviewable-data.md:317`), and its
    `Validation / Testing Strategy` carries those through `all`/`unless`
    expansion; RDR 0006 consumes the same source rule ids/spans for graph lint
    findings.
  - **If wrong**: The trigger is no longer a normalization step that loses
    identity — the kernel-fixed shape carries parsed atoms, so no such step
    exists. It is RDR 0007's unshipped reshape landing with a shape that drops
    a field this RDR's diagnostics name, in which case lint detects an error
    but cannot point reviewers at the specific guard atom to fix.
- **A6 Tag provenance is available to predicate lint.**
  - **Status**: Verified (labels only — see A12 for the reachability half)
  - **Method**: Peer RDR
  - **Evidence**: RDR 0002 `Normative Contracts` require the model to declare every matched or written tag, including `owned`, `observed`, or `recognized` provenance. RDR 0006 `Technical Design` consumes tag provenance and owned-tag write effects from normalized rows for owned-set-before-match and coverage checks. **Scope of this record**: it establishes that provenance *labels* reach lint, and nothing more.
  - **If wrong**: Predicate lint cannot tell an owned tag from an observed one,
    and the owned-tag clause below has no input to test.
  - **Note**: This RDR's owned-tag clause and its "can refuse" clause both quantify
    over *every reachable predecessor*, which needs a predecessor relation, not a
    label. That half is **A12**, not this record — A6 was previously stamped
    `Verified` on evidence that does not reach it.
- **A7 An existence atom projects onto the declared-domain product as a
  per-key presence dimension, so a row group carrying `exists` atoms stays
  provable.**
  - **Status**: Verified
  - **Method**: Design Decision
  - **Evidence**: Stated as a normative clause in `Normative Contracts` (the
    presence-dimension projection), on the derivation recorded in
    `docs/rdr/0003-guard-predicate-exhaustiveness/evidence/research/iter-2-projection-derivation.md`.
    The derivation closes the row-group branch of RDR 0007 A12's disjunction:
    RDR 0007's two-row absence pattern answers one source-state/outcome pair, so
    both rows share one selection group and scoping cannot separate them. RDR
    0007 fixes `exists` as the sole TOTAL operator (`presence == literal`), and
    A2 requires every atom to denote a subset of the scoped product; an `exists`
    atom denotes `{present}` or `{absent}`, which are subsets of a presence
    dimension. The bounded prior-art pass is a recorded clean negative — no
    decision-table/DMN/rule-base verification material in `PapersFast`,
    `StateMachineLit`, or `StateMachineRes` — so this resolves as a Design
    Decision, never as Prior Art. The rejected alternative is dropping `exists`
    atoms from the product, which certifies groups exhaustive that refuse
    `guard_unevaluable` at runtime.
    The optionality input this clause reads is now declared by this RDR's own
    tag declaration model (A11), not requested from a peer — which is what
    unblocked it. RDR 0007 A12 routes the question here and it is answered here.
  - **If wrong**: lint either drops `exists` atoms from the product — certifying
    a row group exhaustive that refuses `guard_unevaluable` at runtime, which
    breaches the §JD-4 narrowing this RDR adopted — or refuses every group
    containing an `exists` atom, which disqualifies this RDR's own representative
    fixture row `foundational-to-cove` and leaves RDR 0007's sanctioned two-row
    absence pattern permanently unprovable, pushing authors back to the
    sentinel-stamping anti-pattern.
- **A8 RDR 0006 can carry the same `guard_unevaluable` narrowing this RDR
  states, so the two documents agree on one lint promise.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: JDR 0001 §JD-4 is **CLOSED (2026-08-22)** by the
    `0003-0006-0007` cluster gate — the venue this record's plan named — and
    both halves are decided at the home. **Recording:** "RDR 0003 is the
    recording document; RDR 0006 cites it and mints no code"
    (`docs/jdr/0001-resolve-kernel-seam.md` §JD-4), on the ground that the
    clause must name the refusing atom and `atom` occurs once in RDR 0006, only
    to delegate atoms to this RDR. The clause it assigns here is already written
    — the narrowing block in `Normative Contracts` above. The two triggers
    genuinely differ and both survive: RDR 0006's fires on a *non-finite
    dimension*, this RDR's on a *fully-finite product whose participating row
    can still refuse*. **Finding code:** no new code and no non-blocking tier —
    RDR 0006 reuses `graph-unprovable-coverage`, which already sits in its
    mandatory blocking table (`docs/rdr/0006-graph-lint-authority-and-guarantees.md:317`)
    and its MVV matrix (`:642`), so reuse costs a widened trigger and no new
    test scaffolding. RDR 0006 records the same assignment on its own Status
    line (`:12`): reuse `graph-unprovable-coverage` "with no new code and no
    non-blocking tier". Gate evidence:
    `docs/rdr/cluster-reconcile/0003-0006-0007/reconcile-report.md`.
  - **Consequent duty on RDR 0006 (not this RDR's edit)**: extend its finding
    contract with an atom-level field, since this RDR's clause requires the
    finding to name the refusing atom. Recorded at the home as a §JD-4 duty and
    discharged at RDR 0006's refine. It does not gate this RDR's lock — the
    clause here is complete and cited, and a missing producer field surfaces as
    a RDR 0006 refine item, not as a false claim here.
  - **If wrong**: the two documents state one lint promise two ways — the
    single-source failure §JD-4 exists to prevent — and an implementer reading
    only RDR 0006 certifies a row group exhaustive that this RDR forbids.

- **A9 A declared set-element universe exists for `contains`, so the
  closed operator vocabulary can be accepted as specified.**
  - **Status**: Verified
  - **Method**: Design Decision
  - **Evidence**: The tag declaration model (A11, `Normative Contracts`) makes
    the element universe a declarable property of a `set`-kind tag, and requires
    domain/kind agreement. `contains` is therefore a finite-domain operator
    whenever its tag declares that universe, which is what Phase 3's acceptance
    gate needs. The Stage 6 A1 harness evaluated a `contains` atom over a
    set-valued tag
    (`evidence/spikes/iter-2/a1-eval-harness/output.txt`), establishing the
    *evaluation* half; this record adds the *declaration* half that makes the
    exhaustiveness claim authorable. The rejected alternative is dropping
    `contains` from the closed vocabulary, which the target-flow lens-row
    predicate needs.
  - **If wrong**: `contains` cannot carry an exhaustiveness claim, the closed
    vocabulary ships one operator it cannot prove, and Phase 3's acceptance
    gate cannot be met with the fixture format as it stands.
- **A10 RDR 0006 accepts that the scoped row group is defined here, and supplies
  only reachability of selection contexts.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence needed**: This RDR's Technical Design previously took the grouping
    context *from* RDR 0006, while RDR 0006's A2 takes the coverage/overlap
    derivation *from* this RDR — a cycle with no floor, and the term `row group`
    occurs zero times in RDR 0006. This pass breaks the cycle by defining the
    group here (rows sharing one source state and one recognized outcome, the set
    RDR 0001 resolves exact-one over). Verification is RDR 0006 confirming it
    reads the same division of labour.
  - **Plan**: raise with A8 at cluster reconcile — same venue, same peer, same
    `Final`-document constraint.
  - **Stage 6 disposition — DOWNGRADED to the cluster-reconcile venue.** RDR
    0006 is now `Draft` and its Refinement Context Direction carries this item,
    so the confirmation is a scheduled edit on an open peer rather than a
    route-back. The cycle this record was opened to break is already broken
    *in this document*:
    the row group is defined here as rows sharing one source state and one
    recognized outcome, so this RDR no longer reads its grouping back from RDR
    0006. What remains is peer confirmation of that division of labour, which a
    peer cannot give from inside this document. Survivable because a divergence would
    surface as a cluster-reconcile finding before either document is
    implemented, and the MVV does not consume the peer's agreement. Travels with
    A8 and A12 — same venue, same peer.
  - **If wrong**: the two documents group rows differently, so a gap proved
    absent in one grouping is present in the other, and the exhaustiveness claim
    means different things on each side of the seam.
- **A11 The tag declaration model carries a finite domain, so any
  exhaustiveness claim has a producer.**
  - **Status**: Verified
  - **Method**: Design Decision
  - **Evidence**: This RDR now **owns** the tag declaration model rather than
    requesting it. The `Normative Contracts` declaration clauses state the value
    kinds, make a finite domain declarable for every kind an exhaustiveness
    claim can range over (enum value set, `bool` by construction, `int`
    `{min..max}`, `set` element universe), and require domain/kind agreement
    with rejection before normalization. A2's derivation therefore consumes
    inputs this document defines, not inputs it hopes a peer will add.
    **Why the home moved here.** The field was previously booked as a producer
    request to RDR 0002 on the belief that RDR 0002 owned the tag declaration.
    It does not: RDR 0002's only normative tag clause requires *provenance*
    alone, `value kind` is normative nowhere in the cluster, and RDR 0002's
    validation vocabulary is entirely structural (malformed, unknown, missing)
    with no type or value-domain category. RDR 0007 (`Final`) states the same
    division from its side — "typed literals and value kinds (bounded integers,
    set universes) are RDR 0003's declarations". The declaration model is the
    guard algebra's input alphabet and belongs with the algebra; RDR 0002 keeps
    authoring location and normalization carriage and cites this model.
    The rejected alternative is filing four fields into RDR 0002, which would
    make that document the de facto owner of a type system it has no normative
    vocabulary, validation category, or charter for.
  - **If wrong**: no row group is ever exhaustiveness-eligible, every group takes
    the blocking inability-to-prove outcome, and this RDR's central guarantee is
    unreachable in practice while remaining true in theory.
- **A12 Predecessor reachability is decidable from the declared model, so
  "every reachable predecessor" is a total test.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence needed**: Two clauses quantify over every reachable predecessor —
    the owned-tag clause and the "can refuse" decision procedure. A6 establishes
    only that provenance *labels* reach lint. The predecessor relation itself is
    a graph property; RDR 0006 owns graph traversal but states no reachability
    contract this RDR can cite, and `reachable predecessor` appears once there and
    not at all in RDR 0002.
  - **Plan**: confirm with RDR 0006 at cluster reconcile alongside A8 and A10 that
    it exposes a predecessor relation over selection contexts, and that lint may
    decide the predicate syntactically over declarations.
  - **Stage 6 disposition — DOWNGRADED to the cluster-reconcile venue.** The
    predecessor relation is RDR 0006's graph property, so this cannot close on
    this RDR's initiative — but RDR 0006 is now `Draft` and its Direction
    carries the reachability contract, so the input has a scheduled producer. Survivable to that venue because the demotion already did the
    protective work — the reachability quantifier was removed from assertion and
    booked here, so no clause in this RDR now claims a decision procedure it
    cannot cite. The MVV's possibly-absent-key case is gated by A7's optionality
    field, not by this record. Travels with A8 and A10.
  - **If wrong**: both clauses are unimplementable as written — lint cannot decide
    whether an owned tag is set before match, nor whether a row "can refuse", so
    the narrowing has no decision procedure.
- **A13 The canonical byte spelling of a set literal is fixed by this RDR for
  both `in` and `contains`.**
  - **Status**: Verified
  - **Method**: Design Decision
  - **Evidence**: RDR 0007 (`docs/rdr/0007-guard-predicate-totality.md:1585-1595`)
    assigns the canonical spelling duty to this RDR and names `in` explicitly.
    A9 books only the *element universe* for `contains`, which is a different
    thing from the literal's encoding: `in` takes a typed literal **set** and is
    already inside the spike/MVV operator subset, so the gap affects an operator
    this RDR claims as proven, not just the deferred one. The duty is now
    discharged: the `Normative Contracts` set-literal clause fixes the spelling
    as an unordered, duplicate-free typed set canonicalized before it enters the
    identity tuple, with repeated elements rejected at parse. The rejected
    alternative is an ordered slice preserving authored order, which would make
    `["mid","large"]` and `["large","mid"]` distinct atoms and break MVV
    Scenario 5's reorder invariant.
  - **If wrong**: two implementations spell the same set literal differently, so
    a fixture that parses under one build fails under another, and the identity
    tuple's `literal` component is not stable across producers.
- **A14 Normalization preserves each atom's `block`, so `all` and `unless` stay
  distinguishable downstream.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence needed**: This one record genuinely belongs to RDR 0002 and does
    **not** travel with the declaration model. A11/A7/A9 moved here because they
    are declaration *semantics*; per-atom `block` retention is normalization
    *carriage* — what survives the sparse-to-normalized transform — which is
    RDR 0002's charter proper. RDR 0007 fixes `Block` as an atom field and this
    RDR's identity tuple and `unless` semantics both depend on it, but RDR 0002
    requires normalization to "combine both into one candidate-row predicate
    set" and guarantees only that a row retains its rule id and source locator.
    The container-level combination is compatible with per-atom `block`
    retention, so this is a gap to close, not a contradiction to resolve.
  - **Plan**: single-field request to RDR 0002 at its refine pass — a
    normalization clause stating that each atom retains the block it was
    authored in. RDR 0002 is `Draft`, so this is a live request rather than a
    route-back. Corroborating evidence that the shape is workable: the Stage 6
    A1 harness carries `block` per atom through its own encoding and shows
    `unless` behaving as a conjunctive exclusion block distinct from `all`
    (`evidence/spikes/iter-2/a1-eval-harness/`) — this RDR's encoding, not RDR
    0002's, so it corroborates rather than discharges.
  - **Stage 6 disposition — DOWNGRADED.** Survivable: no peer asserts the
    opposite, the failure is detectable at the first normalization test, and the
    MVV's authorability does not turn on it (the MVV authors `all`/`unless`
    directly). Not lock-blocking.
  - **If wrong**: `unless` collapses into `all` after normalization, the identity
    tuple loses a component that makes it total, and a row's excluded
    intersection can no longer be subtracted from its accepted assignments.
- **A15 A single published cardinality bound is sufficient for the
  too-large-to-prove refusal to be model-independent.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence needed**: the `Normative Contracts` bound clause now fixes the
    *unit* (cardinality of the scoped product) and the *shape* (one published
    integer constant, reported beside the computed cardinality), which is what
    the clause's cross-implementation determinism promise requires. What is not
    yet established is that a single scalar cardinality is the right
    discriminator: the proof cost of a symbolic or bitset representation may
    depend on dimension count or per-dimension width, not on cardinality alone,
    in which case two products of equal cardinality could differ in provability
    and the bound would not predict the verdict it gates.
  - **Plan**: MVV Scenario 3 already lints "a declared finite product too large
    for the deterministic proof representation" — extend it to assert the
    refusal reports both figures, and to compare two equal-cardinality products
    of differing shape (few wide dimensions vs. many narrow ones) for the same
    verdict. If they diverge, the bound needs a second term and this clause is
    restated before lock.
  - **Stage 6 disposition — DOWNGRADED to the named MVV test.** The clause's
    *unit* and *shape* are pinned in `Normative Contracts` (cardinality of the
    scoped product; one published integer constant reported beside the computed
    cardinality), so what remains open is whether that scalar predicts
    provability — a comparison that is not runnable at draft time because it
    needs finite products built and measured, which is implementation work
    (COMPUTE-DON'T-ARGUE), not a drafting question. Survivable because the
    "If wrong" is bounded and stated: a divergence adds a second bound term and
    restates one clause, it does not disturb the grammar, the identity tuple, or
    the coverage derivation. Named plan: MVV Scenario 3, which already lints an
    over-large product and now also compares two equal-cardinality products of
    differing shape (few wide dimensions vs. many narrow ones) for the same
    verdict.
  - **If wrong**: the published bound does not predict provability, so two
    conforming implementations disagree on the same model and the
    model-independence the clause promises is unmet.

**Method vocabulary** (pick exactly one per assumption):

- **Source Search** — verified against dependency
  source code. Evidence: a greppable `path::Symbol`
  (function/type/const name), **not a bare `file:line`**;
  a commit-SHA permalink only for audit/traceability.
  Standard for libraries. (Why symbol not line: flow
  README *Doctrine*.)
- **Spike** — verified by running code against a live
  service or fixture. Evidence: command run + path to
  captured output.
- **Prior Art** — same property holds in ≥1 named
  external system. Evidence: system + section/page.
- **Derivation** — pure math or proof. Evidence: the
  derivation, shown inline.
- **Design Decision** — a scoping choice this RDR is
  *making* (not *verifying*). Evidence: the decision
  and the alternative explicitly rejected.
- **Peer RDR** — relies on a property defined in
  another RDR. Evidence: RDR ID + section.
- **MVV Test** — the property is testable via the
  Minimum Viable Validation, and the test
  is named in this RDR's Validation section (pending
  implementation at lock time). Evidence: test name.
- **Docs Only** — documentation reading alone.
  **Insufficient** for load-bearing assumptions; allowed
  only when paired with a Spike or Source Search plan
  in the Evidence line.

A `Method: Source Search` whose Evidence cites this
same RDR file — or any path under the RDR's artifact
directory — is self-reference and not Verified. The
cited proof must also support **the specific claim**,
not an adjacent one: confirming a neighboring fact and
stamping the assumption `Verified` is not verification.
The cited symbol must resolve on `main` (a renamed,
deleted, or never-built symbol fails the check).

Any exactness claim such as all/every, first/nearest,
byte-identical, lossless, canonical, deterministic, or
stable order must be covered by a Critical Assumption
Evidence Record or by the Minimum Viable Validation.

## Proposed Solution

### Approach

Use a closed, symbolic predicate-atom model over declared tags, and own the tag
declaration model those atoms are written against. This is the middle path
between table-as-opaque-callback and full statechart adoption: keep RDR 0002's
sparse TOML table as the source, but make every guard a symbolic tag-set
predicate that lint can reason about. The typed alphabet — value kind, finite
domain, optionality, single-valuedness, set-element universe — is stated here because it is what
makes coverage and overlap provable; RDR 0002 owns where a declaration is
authored and how it is carried through normalization, and cites this model for
what it means. Authors keep the RDR 0002 shape of
positive `all` guards and negative `unless` guards; each guard entry is a typed
atom: tag reference, fixed operator, and typed literal or literal set. The
operator set is deliberately small: equality, membership, bounded integer
comparison, existence, and set containment. There is no embedded host predicate
and no free-form expression grammar.

Enforcement of the guard domain is split, per JDR 0001 §D4 (normative in RDR
0007, its landing document): the kernel decides key presence, decides existence atoms from presence
alone, marks a value atom over an absent key unevaluable without consulting
this RDR's evaluator, and combines per-atom verdicts. This RDR's evaluator owns
**value semantics over a present value** — nothing more; it never sees the tag
view.

Each tag declaration supplies the value kind and, when lint must prove
exhaustiveness, the finite domain: enum values, boolean values, set element
universe, or bounded integer range. Runtime guard evaluation is simple predicate
evaluation over the assembled tag-set. Static lint reasons over the same atoms
inside row groups supplied by the normalized model: within one selection
context, it checks the scoped product of participating guard domains, proves
whether candidate rows are mutually exclusive, and proves whether they cover the
declared product.

The initial operator/kind matrix is:

| Operator | Accepted tag value kinds | Literal shape | Lint proof role |
| --- | --- | --- | --- |
| `eq` | enum, boolean, integer, string-like scalar | one typed scalar | Narrows the tag domain to one value. |
| `in` | enum, boolean, integer, string-like scalar | non-empty typed scalar set | Narrows the tag domain to the listed values. |
| `lt`, `lte`, `gt`, `gte` | integer | one typed integer | Narrows a bounded integer domain by comparison; remains runtime-only for an unbounded integer. |
| `exists` | optional scalar or optional set-valued tag | boolean | Tests presence or absence, not value equality. Decided by the kernel from presence alone (JDR 0001 §D4); it never reaches this RDR's evaluator. |
| `contains` | set-valued tag with a declared element universe | non-empty typed element set | Narrows the set-valued domain to assignments containing every listed element. |

This intentionally aligns with prior art that separates positive and negative
guard lists, and deliberately diverges from callback-driven FSM libraries.
Callbacks are ergonomic for application code, but intrastate needs the guard
surface to be reviewable, serializable, and mechanically analyzable.

### Technical Design

Guard predicates are part of the normalized transition model, not a separate
runtime plug-in system. RDR 0002 owns the sparse table source, tag declarations,
provenance, and normalization; this RDR owns the guard predicate grammar and
semantics that normalized rows carry. RDR 0001 consumes evaluated predicates for
exact-one edge selection, and RDR 0006 consumes the symbolic predicate
constraints for graph lint.

The atom shape is fixed at the kernel seam by JDR 0001 §D1 and stated
normatively in RDR 0007 (`docs/rdr/0007-guard-predicate-totality.md:1249-1266`):
four fields — `Key`, `Operator`, `Literal`, `Block`. This RDR cites that shape
and does not restate it. What this RDR adds are the well-formedness rules over
it: `Key` must resolve to a declared tag, `Operator` must be allowed by the
operator/kind matrix above, and `Literal` must parse to the operator's literal
shape. **Source identity is not an atom field** — it is carried by the enclosing
normalized row (`RuleID`, `SourceLocator`, per RDR 0002 and
`internal/resolve/resolve.go::Row`), which is why this RDR's atom identity tuple
is the row's identity joined to the atom's own four fields rather than an
identity stored on the atom.

Because the kernel decides presence, a value atom over an absent key is
unevaluable rather than false, and an existence atom over an absent key is
decided — not unevaluable. This RDR's requirement that every guard resolve to a
declared tag is therefore read against `exists = false`: an author expressing
"this row applies when tag X is absent" writes an existence atom, not a value
atom that happens to miss.

Positive `all` predicates are conjunctive requirements. Negative `unless`
predicates are also conjunctive within the excluded predicate set: if all
`unless` predicates hold, the candidate row is disabled. Normalization may
combine them into one internal constraint object, but authoring keeps the
separation because it reads better and mirrors prior art.
For lint, RDR 0002 supplies the normalized candidate rows. **The row group is
defined here**, not deferred: a scoped row group is the set of normalized
candidate rows sharing one selection context — the same source state and the
same recognized outcome — since that is exactly the set RDR 0001 resolves
exact-one over. RDR 0006 supplies the graph traversal that enumerates which
selection contexts are reachable; it does not define the grouping predicate, and
this RDR does not read one back from it. (The term `row group` appears nowhere
in RDR 0006, and its A2 delegates the coverage derivation here — so a definition
deferred to RDR 0006 would close a cycle with no floor. A10 books the
confirmation that RDR 0006 accepts this division.)

RDR 0006 gates its blocking finding on a contract "that claims closed coverage",
and no document defines how a group makes that claim: this RDR reads it as
**default-on for every scoped row group whose participating dimensions are all
finitely declared**, since an opt-in flag would let the exhaustiveness guarantee
be silently skipped exactly where it matters. If RDR 0006 intends an explicit
per-group annotation instead, that is a divergence to settle at cluster
reconcile alongside A8.

**Leaving a dimension undeclared is not an opt-out.** An undeclared dimension
does not remove the row group from proof; it makes the proof unavailable, which
is the Loud blocking inability-to-prove outcome the `disposition` table assigns
to that input class — never a silent downgrade to unchecked. No authoring
gesture quietly exempts a group: a group is either proved, or it carries a
blocking finding naming the dimension that cannot be proved. Within that group, a row's accepted assignments are the
intersection of all positive `all` atom domains minus the single conjunctive
assignment set matched by the row's full `unless` block. `unless` is not
per-atom negation, and it does not create source-order priority.

Exhaustiveness is only asserted over finite declared domains. Enum, boolean,
and set-element universes are finite by declaration; set-valued tags must be
represented symbolically or as an implementation-equivalent bitset rather than
requiring literal powerset materialization. Integer tags are exhaustive only
when they declare a bounded range; this covers cap counters and iteration limits
without pretending arbitrary integers can be fully partitioned. A guard may
still compare an unbounded integer at runtime, but lint must report that it
cannot prove exhaustive coverage for that dimension. If a finite scoped product
is too large for the implementation to prove deterministically, lint must also
refuse or downgrade the exhaustiveness claim instead of silently capping the
proof. Throughout this RDR "refuse or downgrade" names **one** blocking
inability-to-prove outcome, not a choice between a hard and a soft one; the
normative clauses below fix that reading.

Where this exhaustiveness proof and RDR 0007's aggregation veto disagree, the
**promise narrows**: a green exhaustiveness result must mean resolution
succeeds. Lint therefore cannot certify a row group exhaustive when any
participating row could refuse `guard_unevaluable` at runtime; the runtime veto
is not weakened to make lint's claim true.

JDR 0001 §JD-4 decides that substance ("the *promise* narrows — P5 decides")
but leaves open **which document records it**, and its head clause names RDR
0006's proof rather than this RDR's. This RDR records the narrowing for the
proof it defines, because the constraint binds what its own finite-domain
product may claim. RDR 0006's exhaustiveness clause covers only the
non-finite-dimension refusal, so it does not yet carry this case; RDR 0006 is
now `Draft` and carries the §JD-4 recording question in its Refinement Context.
A8 tracks the agreement the two documents owe each other. Provenance affects lint: recognized tags are fresh event inputs, observed
tags are re-read before matching, and owned tags must have a reachable
predecessor write before a row may match them.

#### Normative Contracts

```normative
A guard predicate MUST be a symbolic atom over a declared tag, not a host
language callback and not a free-form expression string.
```

```normative
This RDR owns the **tag declaration model** — the typed alphabet every guard
atom is written against. A tag declaration MUST carry a value kind, and MAY
carry a finite domain, an optionality marker, a single-valued marker, and (for
set-valued kinds) an element universe. The value kinds are `enum`, `bool`,
`int`, `set`, and opaque scalar. RDR 0002 owns where a declaration is authored and how it is
carried through normalization; this RDR owns what a declaration means. Neither
document restates the other (JDR 0001 P6).
```

```normative
A finite domain MUST be declarable for any kind an exhaustiveness claim can
range over: an `enum` declares its value set, a `bool` is finite by
construction, an `int` declares a `{min..max}` bound, and a `set` declares the
element universe its members are drawn from. A declaration carrying no finite
domain is well-formed — the tag remains runtime-evaluable — but a guard
dimension over it cannot carry an exhaustiveness claim, and lint MUST take the
blocking inability-to-prove outcome for that dimension.
```

```normative
A tag declaration MUST be able to state whether the key may be absent. A key
declared always-present MUST NOT be absent from a conforming evaluation view;
a key declared optional MAY be. Presence is a declared property of the tag, not
an observation of one view: lint decides `exists` projection and the
"can refuse" narrowing from this declaration, never from a runtime trace. The
kernel's runtime presence decision (JDR 0001 §D4, normative in RDR 0007) is a
separate question over a single view and is unaffected by this clause.
```

```normative
A tag declaration MUST be able to state that the tag is **single-valued**: at
most one of its declared domain values holds in any conforming evaluation view.
The marker is meaningful only for a kind carrying a finite domain, and a
`set`-valued kind MUST NOT carry it. Single-valuedness is a declared property of
the tag, not an observation of one view. It is what licenses a consumer to treat
the tag's domain as a partition — mutually exclusive values whose coverage the
scoped product can be grouped over — rather than as independent dimensions. RDR
0006's grouping-dependent lint findings read this field from the declaration and
MUST NOT infer it from a tag's name, its value spelling, or a fixture (JDR 0001
§JD-13).
```

```normative
A declared domain MUST agree with its value kind — an `{min..max}` bound on an
`enum`, an element universe on a scalar kind, or a single-valued marker on a
`set` kind or on a kind carrying no finite domain, is a declaration error and
MUST be rejected before normalization completes. A value appearing in a guard
literal that lies outside its tag's declared domain MUST likewise be rejected
before resolution, so an unsatisfiable atom is a load-time error rather than a
silently-never-matching row.
```

```normative
The initial operator vocabulary MUST be closed and typed: equality,
membership, bounded integer comparison, existence, and set containment. Unknown
operators MUST be rejected during parse or lint before resolution.
```

```normative
Each operator MUST declare which tag value kinds it accepts. A predicate whose
literal cannot be parsed as the declared tag kind MUST be rejected before
resolution.
```

```normative
This RDR's guard evaluator MUST decide value semantics over a present value
only. Key presence, existence atoms, absent-key unevaluability, and the
combination of per-atom verdicts belong to the kernel (JDR 0001 §D4, normative
in RDR 0007); the evaluator MUST NOT read the tag view.
```

```normative
Positive guard atoms MUST live in `all`; negative guard atoms MUST live in
`unless`. Successful row matching MUST NOT depend on source order or
first-match priority.
```

```normative
Lint MAY claim guard exhaustiveness only for finite declared domains: enum
values, booleans, declared set element universes, or bounded integer ranges.
```

```normative
If a guard dimension lacks a finite declared domain, lint MUST refuse or
downgrade an exhaustiveness claim for that dimension rather than treating the
covered examples as complete.
```

```normative
Coverage and overlap checks MUST be scoped to a normalized row group supplied by
the transition/lint model, and MUST evaluate the participating guard dimensions
as one product rather than as independent one-dimensional checks.

A dimension **participates** in a row group when any row in that group carries
an atom over that key, in either `all` or `unless` — not only when the rows
constrain it differently. A key every row constrains identically still bounds
the product and still requires a finite declared domain; it MUST NOT be dropped
from the product because it does not discriminate.
```

```normative
An `exists` atom projects onto the scoped product as a **per-key presence
dimension**: a two-valued dimension `{present, absent}` for that key, alongside
the key's value dimension. An `exists` atom denotes `{present}` or `{absent}`
on it; a value atom denotes a subset of the value dimension and, because a
value atom over an absent key is unevaluable rather than false, implicitly
`{present}`. A key declared always-present contributes no presence dimension —
its `{absent}` assignment is not in the product — which is what makes the
narrowing's negative control provable rather than vacuous. Lint MUST NOT drop
`exists` atoms from the product: a group carrying one stays provable, and
certifying it exhaustive while ignoring the presence dimension is the
false-green this RDR's narrowing forbids.
```

```normative
A declared escape row participates in the coverage identity like any other row:
its accepted assignments are computed from its guard atoms and unioned with its
peers'. An escape row carrying no guard atoms denotes the whole scoped product
and therefore closes coverage by itself. Lint MUST NOT treat "an escape row
exists" as a separate coverage-satisfying fact outside the union, and MUST NOT
exclude escape rows from overlap checks — an escape row that overlaps a guarded
row is the ambiguity RDR 0001 refuses at runtime.
```

```normative
An exhaustiveness claim MUST NOT be stronger than the runtime it describes: lint
MUST NOT certify a row group exhaustive when a participating row can refuse
`guard_unevaluable` under RDR 0007's aggregation veto. Where the two disagree
the lint promise narrows; the runtime veto MUST NOT be weakened.
```

```normative
A withheld exhaustiveness claim MUST be observable, not silent. It takes the
same blocking inability-to-prove form an unprovable dimension already takes —
RDR 0006's `graph-unprovable-coverage` — and MUST name the participating row
and the atom that can refuse, using the source rule/context id every other
predicate diagnostic names. Emitting nothing MUST NOT satisfy this clause: an
exit code alone cannot distinguish a withheld claim from a proved one.
```

> **Peer obligation, not yet agreed (A8).** RDR 0006's finding contract carries
> a stable code, model identity, severity, message, and "the source rule/context
> id or source span" — it has **no atom-level field**. Naming *the refusing
> atom* is therefore a payload extension this RDR requests of a `Final`
> document, not a capability it can assume. The route-back in A8 must carry it
> explicitly; until then the atom-naming half of the clause above is a stated
> requirement with no producer, and the row-naming half is satisfiable today.

```normative
Lint MUST report every defect it can decide in one pass over a row group, not
the first one it encounters. Withholding a group's exhaustiveness claim MUST NOT
suppress overlap findings, coverage findings, or further withholding reasons for
that same group: each unprovable dimension, each refusing row, each overlapping
pair, and any coverage gap over a provable product is its own finding. A group
with two refusing rows MUST emit a finding for each, so the emitted set does not
depend on row or dimension iteration order — the same source-order independence
matching already requires. A withheld claim and a coverage or overlap finding
MAY be emitted together; what MUST NOT be emitted is a green exhaustiveness
result alongside any withholding reason.
```

```normative
"Refuse" and "downgrade" are one outcome, not an author's choice: every case
this RDR sends to refuse-or-downgrade — a non-finite dimension, a finite
product too large to prove deterministically, and a withheld claim under the
narrowing above — MUST produce the blocking inability-to-prove finding. This
RDR MUST NOT mint a non-blocking warning category for these cases. JDR 0001
§JD-4 (closed 2026-08-22) decides the same for lint: RDR 0006 mints no code and
gains no non-blocking tier for this class, reusing `graph-unprovable-coverage`.
"Downgrade" therefore names the same blocking outcome as "refuse" on both sides
of the seam, and never a silent or advisory one.
```

```normative
A row "can refuse" when the row carries a value atom whose key is not declared
present-for-every-reachable-predecessor — a property decided from the declared
model, not from a runtime trace or a witness input. Lint MUST decide this
syntactically over declarations so the test is total; it MUST NOT withhold a
claim merely because some assignment in the product is unreached.
```

```normative
A set literal — the right-hand side of `in` and of `contains` — has one
canonical spelling: an unordered set of typed elements, duplicate-free, and
compared as a set. Two authored spellings differing only in element order or in
repeated elements MUST parse to the same literal and therefore to the same atom
under the identity tuple; a repeated element MUST be rejected at parse rather
than silently collapsed. Implementations MUST canonicalize before the literal
enters the identity tuple, so reordering a set literal cannot change a
diagnostic — the source-order independence MVV Scenario 5 requires.
```

```normative
Set-valued guard domains MUST be proved with a deterministic symbolic or
bitset-equivalent representation. If the finite product is too large for that
proof, lint MUST refuse or downgrade the exhaustiveness claim rather than
silently capping enumeration.
```

```normative
"Too large to prove" MUST be a declared, model-independent bound, not an
implementation's incidental limit: the implementation MUST publish the bound it
enforces, and the same model MUST receive the same verdict on every conforming
implementation. A bound discovered by exhausting memory or wall-clock is not a
conforming bound. The refusal diagnostic MUST report the product size it
computed and the bound it exceeded, so an author can tell an over-large product
from an undeclared dimension.

The quantity both figures report is the **cardinality of the scoped product** —
the number of assignments in it, the product of every participating dimension's
declared domain size — not a bitset width, byte size, or row count. The bound is
a single integer constant published by the implementation and reported in the
diagnostic beside the computed cardinality; it is not per-model, per-group, or
configurable per run, since either would make the verdict model-dependent. Two
conforming implementations MAY publish different bounds, but each MUST return
the same verdict for the same model and MUST report which bound it applied.
```

```normative
Overlap and coverage diagnostics MUST name the source rule id or context id
that contributed each predicate involved in the finding.

A coverage gap has no contributing predicate — it is an absence — so it is
attributed differently: a gap finding MUST name the selection context that
scopes the group, every source rule id in that group, and at least one concrete
uncovered assignment from the scoped product, so the author can see which
combination no row accepts. Naming the group without a witness assignment MUST
NOT satisfy this clause.
```

```normative
Predicate lint MUST distinguish owned, observed, and recognized tags. A row
that matches an owned tag MUST be rejected unless every reachable predecessor
sets or preserves that tag before the match.
```

#### `authority` — source-authority census

Cue: the exhaustiveness verdict has more than one candidate source of truth
(this RDR and RDR 0006), and guard-domain enforcement is split across the kernel
and this RDR's evaluator.

| Input / decision | Writer (canonical) | Readers | Sibling arms | Which is canonical |
| --- | --- | --- | --- | --- |
| Key presence | Kernel (RDR 0007, JDR 0001 §D4) | This RDR's lint; RDR 0001 selection | This RDR's evaluator — **explicitly not an arm**; it never reads the tag view | Kernel |
| Existence-atom verdict | Kernel, from presence alone (`presence == literal`) | Lint's presence dimension (A7) | Value atoms over the same key (unevaluable on absence) | Kernel |
| Value-atom verdict over a present value | This RDR's evaluator | Kernel combinator | None | This RDR |
| Per-atom verdict combination | Kernel (strong Kleene) | RDR 0001 exact-one selection | None | Kernel |
| Exhaustiveness verdict for a row group | **This RDR** (finite-domain product proof) | RDR 0006 lint findings; RDR 0005 envelope | RDR 0006's own exhaustiveness clause — non-finite-dimension case only | **Contested — A8.** This RDR states the `guard_unevaluable` narrowing; RDR 0006 holds the §JD-4 tolerance but does not carry it. A8 closes on §JD-4 assigning ONE recording document, not on both restating it. |
| Guard atom shape | Kernel seam (JDR 0001 §D1, normative in RDR 0007) | This RDR's grammar; RDR 0002's normalizer | Shipped `Row.Guard string` — the pre-reshape form | Kernel seam (specified); the shipped string form is superseded, not an arm |
| Source identity for diagnostics | RDR 0002 normalization (`RuleID`, `SourceLocator`) | This RDR's diagnostics; RDR 0006 findings | None | RDR 0002 |
| Tag declaration model (value kind, finite domain, optionality, single-valuedness, element universe) | **This RDR** (`Normative Contracts` declaration clauses) | RDR 0002's authoring schema and normalizer; RDR 0006's lint input contract; this RDR's own product proof | RDR 0002's prose schema list, which names `value kind` non-normatively | **This RDR.** RDR 0002 declares provenance normatively and nothing else about a tag's type; RDR 0007 (`Final`) states value kinds and set universes "are RDR 0003's declarations". RDR 0002 cites this model for meaning and owns where it is authored. |
| Per-tag optionality (which keys may be absent) | **This RDR** (declaration model) | A7's presence dimension; the "can refuse" narrowing | Kernel runtime presence (JDR 0001 §D4) — a different question over one view, not an arm | This RDR |
| Set-element universe for `contains` | **This RDR** (declaration model) | This RDR's finite-domain proof; Phase 3 acceptance gate | None | This RDR |
| Single-valued marker (domain is a partition) | **This RDR** (declaration model) | RDR 0006's `graph-single-valued-state` and its MVV assertion | RDR 0006's lint input contract, which attributed the field to this model before it declared one | **This RDR** — JDR 0001 §JD-13 (2026-08-22); RDR 0006 cites it like the rest of the model |
| Tag provenance (`owned`/`observed`/`recognized`) | RDR 0002 (its one normative tag clause) | This RDR's owned-tag clause (A6); RDR 0004; RDR 0006 | None | RDR 0002 — unchanged by the declaration-model rehoming |
| Withheld-claim lint artifact | **This RDR** (states the blocking form) | RDR 0006 emits it; RDR 0005 envelopes it | RDR 0006's `graph-unprovable-coverage`, previously scoped to non-finite dimensions | **This RDR** — JDR 0001 §JD-4 (closed 2026-08-22) names it the recording document; RDR 0006 widens `graph-unprovable-coverage`'s trigger and mints no code |

#### `disposition` — input class × outcome

Cue: the draft sets outcomes over input classes (refuse, downgrade, prune,
unevaluable).

| Input class | Runtime outcome | Lint outcome | Diagnostic minted | Silent or loud |
| --- | --- | --- | --- | --- |
| Every `all` atom decides true; `unless` block not fully true | Row qualifies | Row contributes its accepted assignments | none | — |
| `all` atom decides false | Row pruned | Row contributes nothing | none | Silent by design — a decided false is not a defect |
| Full `unless` block decides true | Row disabled | Excluded intersection subtracted | none | Silent by design |
| Value atom over an absent key | `guard_unevaluable` refusal (RDR 0007 veto) | Exhaustiveness claim **withheld** for that group | Runtime: RDR 0007's per-row/per-atom payload (key, block, reason `absent`). Lint: the blocking inability-to-prove finding, naming the row and the refusing atom | Loud — both surfaces mint an artifact |
| Existence atom over an absent key | Decided (`presence == literal`) — never unevaluable | Selects `{absent}` on the presence dimension (A7) | none | — |
| Zero rows qualify | RDR 0001 refuses | Coverage gap if the product is provable | `graph-coverage-gap` (RDR 0006), naming the selection context, every rule id in the group, and one uncovered assignment | Loud |
| Two or more rows qualify | RDR 0001 refuses — never first-match | Overlap finding | Overlap code (RDR 0006), naming both source rule ids | Loud |
| Guard dimension lacks a finite declared domain | Evaluable at runtime | **Refuse or downgrade** the group's claim, one finding per unprovable dimension — never treat examples as complete, never stop at the first | Inability-to-prove finding naming the dimension | Loud |
| Finite product too large to prove deterministically | Evaluable at runtime | **Refuse or downgrade** — never silently cap enumeration | Inability-to-prove finding reporting the computed product cardinality and the published bound | Loud |
| Unknown operator / unknown tag / operator–kind mismatch / literal parse failure | Rejected before resolution | Rejected at load | Predicate semantic kind (this RDR) → RDR 0006 finding → RDR 0005 envelope | Loud |
| Row group is domain-exhaustive but a participating row can refuse | Refusal stands | Claim withheld — the narrowing; one finding per refusing row, never just the first | RDR 0007 payload at runtime **and** a blocking lint finding naming the refusing atom — absence of a green result is not the artifact | Loud |

#### `trace` — desk trace over the MVV

Cue: ten normative clauses bear on one output surface — the exhaustiveness
verdict for a row group. Walked stepwise against the MVV, with witnesses from
`evidence/spikes/guard-fixture.toml`.

| Step | Assertions in force | Witness | Verdict |
| --- | --- | --- | --- |
| 1. Load the RDR/kata slice as normalized candidate rows | atoms-not-callbacks; closed typed operator vocabulary; operator declares accepted kinds | `profile-to-grounding` parses to `profile in [mid,large]` + `unless prelock_iterations gte 3`; operators all in the matrix | OK |
| 2. Group rows by selection context | coverage/overlap scoped to a normalized row group, evaluated as one product | `profile-to-grounding` and `foundational-to-cove` share `match status eq Draft` → one group | OK |
| 3. Build the scoped product from declared domains | exhaustiveness only over finite declared domains; every *participating* dimension enters the product, including one every row constrains identically | `profile` = 4 enum values; `prelock_iterations` = `{0..3}`; `status` participates too — the group's shared `status eq "Draft"` atom does not discriminate but still requires a declared domain | OK, and it is the participation clause that puts `status` in the product — under a discriminates-only reading it would silently drop out |
| 4. Project each atom onto the product | every atom denotes a subset of the scoped product (A2) | `profile eq "foundational"` → `{foundational}`; `profile in [mid,large]` → `{mid,large}` | OK |
| 5. Project the `exists` atom | A7 presence dimension — stated as a normative clause | `foundational-to-cove` carries `cluster_eligible exists = true`; `cluster_eligible`'s declared domain `{true,false}` is complete, so no *value* element selects presence — the atom lands on the separate `{present, absent}` presence dimension the projection clause adds | OK — witness: the presence-dimension clause in `Normative Contracts`, on the derivation in `evidence/research/iter-2-projection-derivation.md` |
| 6. Compute coverage | `union(row_i accepted) == scoped product` | Both kinds of group compute: `exists`-bearing groups over the product including the presence dimension (step 5), `exists`-free groups over the value dimensions alone (`continue-prelock-lenses`, `reconcile-rewind-legality`) | OK |
| 7. Compute overlap | any non-empty pairwise intersection; no source-order priority | `profile in [mid,large]` ∩ `profile eq foundational` = ∅ → rows disjoint on that dimension | OK |
| 8. Apply the runtime-veto narrowing | claim MUST NOT be stronger than the runtime; withhold if a participating row can refuse | MVV Scenario 6's row group: domain-exhaustive, one value atom over a possibly-absent key → claim withheld | OK — and the reason A7's presence dimension must not silently drop `exists` atoms |
| 9. Emit the verdict | diagnostics name the contributing source rule/context id; every decidable defect is reported, not the first; a gap additionally names the context, all group rule ids, and one uncovered assignment | `RuleID` + `SourceLocator` ship on `internal/resolve/resolve.go::Row`; the uncovered assignment is computed from the scoped product built at step 3 | OK |
| 10. Cross-document agreement | exactly one document records the narrowing; the other cites it | RDR 0006 carries no `guard_unevaluable` narrowing | **GAP — booked as A8.** Not a contradiction inside this draft; a divergence with an open peer (RDR 0006 is `Draft` since 2026-08-21) that needs a §JD-4 assignment, not duplicated prose. |

No CONTRADICTION row. One GAP row remains — A8, the cross-document assignment —
carrying a named plan and not asserted as already true. Step 5's former gap
closed when the presence-dimension projection was stated; A9 and A11 have no row
here because the declaration model they record is an input to step 3, which the
trace already exercises.

#### Load-Bearing Decisions

- **Identity** — a guard atom is identified by the total tuple
  `(RuleID, SourceLocator, key, block, operator, literal)`, never by its index
  within `all` or `unless`. Position is not a field of the kernel-fixed atom
  shape (JDR 0001 §D1: key, operator token, literal, block), and an index-based
  handle would contradict this RDR's own source-order independence: the same
  tuple must identify the same atom under the reordering MVV Scenario 5
  requires. The tuple is total precisely because one row may carry two atoms
  over one key (RDR 0007 A5's conjoined value row), which `(key, block)` alone
  cannot separate. Semantic equality is `(tag, operator, literal)`.
  **Which equality each operation uses is fixed, not left to the implementer**:
  diagnostics, deduplication, and any "same atom" claim use the **identity
  tuple**; only domain computation — deciding what subset of the product an atom
  denotes — uses **semantic equality**, because two atoms from different rules
  that constrain a dimension identically denote the same subset. Overlap
  detection is therefore *not* an equality test at all: it intersects the rows'
  accepted assignment sets, so two byte-identical guards in two different rule
  ids correctly produce an overlap finding naming both rule ids, rather than
  being deduplicated into one row.
- **Wire / byte format** — RDR 0002 owns the TOML container; this RDR owns the
  guard atom grammar embedded in that container.
- **Naming** — the canonical name is "guard predicate"; rejected alternatives
  are "condition callback" and "guard expression" because both invite opaque
  host logic.
- **Operator semantics** — equality compares a tag value to one typed literal;
  membership checks a scalar tag against a typed literal set; bounded integer
  comparison uses `lt`, `lte`, `gt`, and `gte`; existence checks presence of an
  optional tag value; set containment checks declared set-valued tags against a
  typed element set.
- **Finite-domain proof** — exhaustive coverage is a lint claim over the
  declared tag-domain product for a normalized row group, not over examples
  observed in fixtures. Unbounded dimensions and finite products that cannot be
  represented deterministically remain runtime-evaluable but cannot satisfy an
  exhaustiveness proof.
- **Error ownership** — this RDR names predicate semantic kinds over present
  values; RDR 0007 owns the per-row/per-atom `guard_unevaluable` payload (JDR
  0001 §D4), RDR 0006 maps graph-level predicate failures to lint finding codes,
  and RDR 0005 maps those findings/refusals to the CLI envelope.
- **Selection / predicate** — a row qualifies only when every `all` atom is true
  and the `unless` predicate set is not fully true; if multiple rows qualify,
  RDR 0001's exact-one resolver refuses instead of choosing by priority.

#### Round-Trip / Inverse Invariants

This RDR introduces no encode/decode pair. Parse/render fidelity for the sparse
TOML source belongs to RDR 0002.

#### Illustrative Code

Illustrative predicate shape only:

```toml
[rule.guard.all]
status.eq = "Draft"
profile.in = ["mid", "large", "foundational"]

[rule.guard.unless]
prelock_iterations.gte = 3
```

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Sparse transition-table container | RDR 0002 | Pending | This RDR assumes guards live inside RDR 0002's `all`/`unless` blocks. |
| Deterministic exact-one resolver | RDR 0001 | Pending | Runtime selection refuses zero or multiple matching rows. |
| Guard predicate grammar and finite-domain semantics | This RDR | Introduced | Lint and runtime share one symbolic predicate model. |
| **Tag declaration model** (value kind, finite domain, optionality, set-element universe) | **This RDR** | **Introduced** | The typed alphabet guard atoms are written against. Previously booked as a producer request to RDR 0002 (A11/A7/A9); rehomed here because RDR 0002's only normative tag clause requires provenance alone, `value kind` is normative nowhere, RDR 0002's validation vocabulary has no type axis, and RDR 0007 (`Final`) states that value kinds and set universes "are RDR 0003's declarations". RDR 0002 owns authoring location and normalization carriage and cites this model. |
| Accessor read/write safety | RDR 0004 | Pending | Guard evaluation consumes tag values after accessor binding; it does not execute accessors. |
| Graph lint authority | RDR 0006 | Pending | Exhaustiveness and overlap findings become blocking lint there. |
| Kernel guard-domain enforcement and `guard_unevaluable` payload | RDR 0007 | Pending | The kernel decides presence and existence atoms; this RDR's evaluator narrows to value semantics over a present value. |
| Atom `block` retention through normalization | RDR 0002 | Requested — **A14** | The one field that stays a request: per-atom `block` retention is normalization *carriage*, RDR 0002's charter proper, not declaration *semantics*. RDR 0007 requires `Block` on every atom and this RDR's identity tuple uses it, but RDR 0002's normalization contract promises only that a candidate row retains rule id and locator. Not a contradiction; a retention guarantee this RDR needs stated. |
| Authoring location + normalization carriage for tag declarations | RDR 0002 | Pending | RDR 0002 owns `[tags.<tag>]` placement and carrying declarations through to candidate rows; it cites this RDR's declaration model for what the fields mean. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| User-facing failures | `internal/cli/clierr`, `internal/cli/respond` | Must verify error-code coverage | Reuse | Predicate parse/lint failures should use existing CLI envelopes. |
| Guard evaluator | None found under `internal/` | New semantic surface | Introduce | This RDR owns the evaluator semantics, not command I/O. |
| Transition model source | Pending RDR 0002 | Not implemented | Extend peer | Guard atoms are embedded in table rows/contexts. |
| Tag declaration model | RDR 0002 declares provenance normatively; `value kind` is normative nowhere; no domain, optionality, or element-universe field exists in any RDR | The typed alphabet the guard grammar quantifies over had no normative home | Introduce here | This RDR states the model; RDR 0002 cites it and owns authoring location plus normalization carriage. |

### Decision Rationale

Joint-check: fired → RDR 0007 is the normative home of the guard seam (JDR 0001
§D1, §D4). It fixes the parsed-atom row shape, the value-only per-atom
`GuardEvaluator` seam, the kernel-exported existence operator token and boolean
literal forms RDR 0002's normalizer emits, exact key identity at the kernel, and
the per-atom `guard_unevaluable` payload. This RDR cites those rather than
restating them. It also records the §JD-4 narrowing for its own proof — §JD-4
decides the substance but leaves the recording document open, and names RDR
0006's proof in its head clause; A8 carries the resulting sibling obligation.

The fixed symbolic atom model best matches the user's outcome: flow authors can
write conditional edges, and lint can still prove whether those edges are
complete and mutually exclusive — with the scope A8 currently bounds. One class
is not yet provable: a group whose narrowing verdict depends on the
cross-document agreement A8 tracks. The `exists`-bearing class was the other,
and it closed when the presence-dimension projection was stated as a normative
clause (A7). Consumers inheriting
this promise — RDR 0006's findings and RDR 0005's envelope — inherit those
bounds with it. It preserves RDR 0002's readable `all`/`unless`
authoring form while giving RDR 0006 finite-domain constraints to analyze. It
also follows the `MODEL-transition.md` thesis directly: a guard is a predicate
over tags, not a second decision channel.

This choice aligns with `transitions` on positive/negative guard factoring,
with statechart/SCXML/Sismic on treating counters and rewinds as extended-state
guarded transitions, and with Statewright-style guard objects on keeping
predicate data small. It deliberately diverges from those systems where they
would undermine the local contract: callback guards are not statically
exhaustive, a full statechart runtime would replace RDR 0002's sparse table
source, and an external guardrail harness would violate the no-orchestrator
project lens.

The chosen approach survived the premortem. If it shipped and failed, the most
likely failure would be an operator set too small for a real RDR/kata guard,
causing authors to request escape hatches. The mitigation is to verify the
target-flow fixtures in Resolve before lock and add only concrete operators
proven necessary there. The recommendation still holds because the alternatives
that solve expressiveness up front lose the static proof this RDR exists to
protect.

## Alternatives Considered

### Alternative 1: Fixed Typed Predicate Atoms Inside RDR 0002

**Description**: Represent every guard as a symbolic atom over declared tags,
with a closed operator set, tag provenance, and finite-domain declarations for
exhaustive lint. Author predicates in RDR 0002's positive `all` and negative
`unless` blocks.

**Prior-art alignment**: `MODEL-transition.md` defines transition resolution as
tag-set predicate matching. `transitions` supplies the `conditions`/`unless`
factoring. Statewright's guard objects show the useful small shape (`field`,
`op`, `value`) while the project rejects its external harness model.

**Pros**:

- Keeps guards reviewable and serializable as data.
- Lets lint reason about overlap, coverage, and owned-tag read-before-write
  without executing user code.
- Covers the seed's equality, integer comparison, and membership needs.
- Aligns with RDR 0002's sparse table shape and RDR 0001's exact-one resolver.
- Leaves SCXML/statechart tools available as external validators rather than
  replacing the local source format.

**Cons**:

- Requires finite-domain declarations for strong exhaustiveness claims.
- May need later operator additions if Resolve finds a target-flow gap.
- Requires RDR 0002 to carry tag provenance and source identity through
  normalization.

**Reason for selection**: This is the smallest model that can express the known
guards while preserving static verification and the no-orchestrator project
lens.

### Alternative 2: Table-Level Tag Predicates Without Separate Guard Syntax

**Description**: Remove the conceptual distinction between "match" and "guard":
every transition row is just a set of tag predicates on the left-hand side and
tag writes on the right-hand side. `all`/`unless` becomes documentation or
syntactic sugar over row predicates.

**Prior-art alignment**: This is the purest reading of `MODEL-transition.md`:
there are no separate guards, only predicates on tags.

**Pros**:

- Conceptually clean: one predicate language for all row matching.
- Makes lint rules uniform for outcome selection, profile routing, cap
  counters, and rewind legality.
- Avoids a second nested guard namespace in the table format.

**Cons**:

- Conflicts with RDR 0002's existing `all`/`unless` authoring split.
- Loses the readability of positive requirements versus negative exclusions.
- Makes it harder to align with `transitions` prior art and with the way users
  usually discuss guards.

**Reason for rejection**: This is semantically attractive but too disruptive to
the neighboring RDR. The chosen approach keeps `all`/`unless` as the authoring
surface while normalizing to tag predicates internally.

### Alternative 3: Tiny Boolean Expression Grammar

**Description**: Put a small expression language in guard fields, such as
`status == "Draft" && profile in ["mid", "large"] && iterations < 3`.

**Prior-art alignment**: SCXML `cond` expressions, XState parameterized guards,
and Statewright guard objects all demonstrate guard expressions or structured
guard references, but they do not all provide static exhaustiveness over the
local tag domains.

**Pros**:

- Compact for authors who already know expression syntax.
- Can express nested boolean logic without row factoring.
- Familiar from SCXML/statechart tools.

**Cons**:

- Requires parser, precedence, type-checking, formatting, and diagnostics that
  are larger than the known problem.
- Makes proof/debug output harder to map back to individual authored atoms.
- Encourages clever predicates instead of simple reviewable guard rows.
- Reintroduces much of SCXML's expression machinery without adopting SCXML's
  full validation model.

**Reason for rejection**: The extra expressiveness is not needed for the target
flows and increases the surface that lint must prove.

### Alternative 4: Embedded Host Predicates / FSM Callback Guards

**Description**: Let table rows reference Go functions or callbacks that return
true or false at runtime.

**Prior-art alignment**: enetx/fsm, stateless/qmuntal-stateless, looplab-fsm,
and XState all prove this as a runtime pattern. qmuntal-stateless even states
same-state guards must be mutually exclusive.

**Pros**:

- Maximum expressiveness.
- Matches many FSM libraries' guard pattern.
- Easy to add without designing a predicate grammar.
- Can reuse Go types and helper functions directly.

**Cons**:

- Static lint cannot prove exhaustiveness or mutual exclusion without executing
  arbitrary code.
- Reviewers must inspect scattered Go functions instead of one transition
  artifact.
- Replay determinism depends on callback purity and hidden dependencies.
- The prior art's mutual-exclusion requirement becomes a runtime convention, not
  a design-time proof.

**Reason for rejection**: It gives up the static verification requirement in
the problem statement.

### Alternative 5: Full Statechart/SCXML/Sismic Guard Semantics

**Description**: Adopt a statechart-style guard language or runtime/spec with
hierarchy, datamodel expressions, Design-by-Contract invariants, and conditional
transitions.

**Prior-art alignment**: The `../state-machines` compiler eval found scxmlcc
and Sismic are strong at expressing the RDR hard parts in grammar: cap counters,
rewinds, coupled status, and conditional skips.

**Pros**:

- Strong prior art for guarded transitions and graph tooling.
- Can model hierarchy and extended state directly.
- Could validate the legal graph outside intrastate's own implementation.
- Captures bounded counters and rewind edges more naturally than a flat FSM.

**Cons**:

- Imports a larger conceptual model than intrastate needs.
- Risks making intrastate a workflow engine instead of a thin deterministic
  resolver/lint CLI.
- Still needs local restrictions to make guard predicates statically provable.
- Adds another source-of-truth format beside RDR 0002's sparse TOML model.

**Reason for rejection**: RDR 0002 already chooses a sparse local table format;
this RDR should define the predicate seam inside that format, not replace the
format with a statechart runtime. SCXML/Sismic remain good validator or
cross-check prior art for RDR 0006.

### Alternative 6: External Guardrail DSL / Harness Model

**Description**: Move guard decisions into a workflow harness that constrains
tools, commands, state transitions, and approvals around the agent.

**Prior-art alignment**: Statewright exposes conditional transitions as
structured guards over context data and supports per-state max iterations.

**Pros**:

- Strong for external enforcement of phase behavior and tool access.
- Makes guard conditions explicit and inspectable.
- Already includes workflow-level controls such as approval gates and iteration
  caps.

**Cons**:

- It is an orchestrator/harness, which the `../state-machines` lens explicitly
  rejects for intrastate.
- Moves the source of truth outside the RDR transition model.
- Solves agent/tool enforcement, not the internal guard predicate grammar needed
  by RDR 0001/0002/0006.

**Reason for rejection**: Useful contrast only. Intrastate needs a local
predicate contract inside its resolver/lint data, not a harness wrapped around
skills.

### Briefly Rejected

- **First-match priority guards**: Rejected because source order would become
  behavior, conflicting with exact-one selection and review-safe reordering.
- **Regular-language FSM compilers such as Ragel/re2c**: Rejected by the sibling
  POC as the wrong category; they recognize byte/token streams and push guards,
  counters, and rewinds into host code.

## Trade-offs

### Consequences

- Guard authoring stays data-first and reviewable.
- Static lint can produce gap and overlap diagnostics from the same predicate
  model runtime resolution uses.
- Unbounded integer or free-form string dimensions cannot receive silent
  exhaustiveness claims; authors must declare finite domains or accept a lint
  limitation.
- **`exists` is in the closed vocabulary and in the proof.** The
  presence-dimension projection (A7) puts an `exists` atom on a `{present,
  absent}` dimension of the scoped product, so a row group carrying one stays
  provable — including this RDR's own fixture group (`profile-to-grounding` ‖
  `foundational-to-cove`) and RDR 0007's sanctioned two-row absence pattern. A
  key declared always-present contributes no presence dimension, which is what
  keeps the narrowing's negative control tight rather than vacuous.

### Risks and Mitigations

- **Risk**: The closed operator set misses a real target-flow guard.
  **Mitigation**: Resolve must encode representative RDR and kata guards before
  lock; add only operators proven necessary by that fixture.
- **Risk**: Finite-domain declarations feel like boilerplate.
  **Mitigation**: Keep declarations close to tag definitions in RDR 0002's
  model and make diagnostics explain when a missing domain blocks proof.
- **Risk**: Multi-dimensional coverage is accidentally checked one dimension at
  a time.
  **Mitigation**: Require MVV cases with a gap and overlap that only appear in a
  scoped row-group product.
- **Risk**: Set-valued domains are implemented by naive powerset enumeration.
  **Mitigation**: Require symbolic or bitset-equivalent proof and explicit
  refusal/downgrade for finite products that are too large to prove.
- **Risk**: `unless` semantics are misunderstood as per-atom negation.
  **Mitigation**: Normatively define it as an excluded predicate set and require
  fixtures that demonstrate mixed `all`/`unless` behavior.

### Failure Modes

Visible failures should be typed parse or lint failures: unknown operator,
operator/tag-kind mismatch, literal parse failure, unknown tag, non-exhaustive
finite domain, overlapping candidate rows, guard dimension not provable because
it lacks a finite domain, or finite scoped product too large to prove. A guard
that cannot be decided at runtime surfaces as RDR 0007's `guard_unevaluable`
refusal, not as a predicate parse kind owned here. Silent
failure would be a false exhaustiveness claim; the recovery path is to keep
every exactness claim tied to A2 and the MVV fixture before Final.

## Implementation Plan

### Prerequisites

- [x] **A1 closed.** The Stage 6 evaluation harness ran: representative rows are
  *evaluated* against tag inputs with asserted qualify/prune verdicts (28 pass /
  0 fail), all eight operators carry an evaluated atom, and five typed-rejection
  negative controls confirm the vocabulary is closed and typed. Evidence:
  `evidence/spikes/iter-2/a1-eval-harness/{main.go,output.txt}`. Owned here;
  needed no peer. Vocabulary sufficiency is independent of the declaration
  model — evaluation consumes tag values, not declared domains.
- [x] **A2 closed.** The derivation's finite-domain inputs are declared by this
  RDR's own tag declaration model (A11), so the inputs are authorable and the
  derivation stamps without circularity.
- [x] A3, A4, A5 verified (A5 re-verified against the kernel-fixed atom shape,
  JDR 0001 §D1). A6 verified for provenance *labels* only; its reachability half
  is A12.
- [x] **A11 closed** (parent of A7 and A9). The tag declaration model is stated
  here as normative contracts — value kinds, finite domains per kind,
  optionality, element universe, and domain/kind agreement. Rehomed from a
  producer request to RDR 0002 after that document was found to own no normative
  tag-type vocabulary at all; RDR 0007 (`Final`) independently assigns value
  kinds and set universes to this RDR.
- [ ] **A12 closed.** Route: RDR 0006 confirms it exposes a predecessor relation
  over selection contexts. Done-condition: the owned-tag clause and the "can
  refuse" clause each cite a reachability contract instead of asserting one.
  Venue: cluster reconcile.
- [x] **A13 closed.** The canonical set-literal spelling is stated as a
  normative clause: unordered, duplicate-free, canonicalized before entering the
  identity tuple, repeats rejected at parse. Owned here — RDR 0007 assigned this
  duty to this RDR. Affects `in`, which is inside the claimed-proven operator
  subset.
- [ ] **A15 closed.** Route: MVV Scenario 3 compares two equal-cardinality
  products of differing shape for the same verdict, confirming a scalar
  cardinality bound predicts provability. Owned here; discharged by the MVV
  rather than a peer.
- [ ] **A10 closed.** Route: RDR 0006 confirms the row-group division of labour
  this RDR now states. Venue: cluster reconcile, with A8 and A12.
- [x] **A7 closed.** The presence-dimension projection is stated as a normative
  clause, on the derivation in `evidence/research/iter-2-projection-derivation.md`.
  Its optionality input is declared by this RDR's own declaration model, so the
  clause no longer waits on a peer. RDR 0007 A12 routed this question here and
  it is answered here.
- [x] **A8 closed — 2026-08-22, by the `0003-0006-0007` cluster gate.**
  JDR 0001 §JD-4 is **CLOSED** and names **this RDR as the recording document**;
  RDR 0006 cites the clause and mints no code, reusing
  `graph-unprovable-coverage`. The done-condition is met: §JD-4 names one
  recording document, and the other carries a citation rather than a
  restatement. The clause this RDR already states (`Normative Contracts`, the
  narrowing plus the withheld-claim form) is the record and needs no new prose.
  §JD-4 also settled the second half: **lint gains no warning category and no
  non-blocking tier for this class**, which retires this RDR's standing
  dependency on that question. Consequent duty on RDR 0006, recorded at the
  home: extend its finding contract with an atom-level field, since this RDR's
  clause requires the finding to name the refusing atom. Evidence:
  `docs/rdr/cluster-reconcile/0003-0006-0007/`.
- [x] **§JD-13 discharged (2026-08-22).** The `0003-0006-0007` cluster gate
  found that RDR 0006's lint input contract attributes "single-valued grouping"
  to *this RDR's* tag declaration model, which never declares it — `single-valued`
  occurs zero times here and zero times in RDR 0002 — while a mandatory blocking
  code (`graph-single-valued-state`) and an MVV assertion depend on it. JDR 0001
  §JD-13 decides that **this RDR's declaration model gains the field**, beside
  value kind, finite domain, optionality, and set-element universe. **Done** —
  the declaration model now carries a single-valued marker and a normative
  clause states what it licenses; the domain/kind agreement clause rejects it on
  a `set` kind or a kind with no finite domain. RDR 0006 cites it at its refine.
- [x] **§JD-14 noted (2026-08-22).** The same gate found RDR 0006's
  invariant 3 granting an escape-row overlap exemption this RDR normatively
  forbids. §JD-14 decides **this RDR's reading governs** — escape rows
  participate in the coverage union and are never excluded from overlap checks —
  and that "claims closed coverage" is **default-on**, as this RDR reads it. The
  repair is RDR 0006's; no edit is owed here, and the clauses stand as written.
- [ ] **A14 closed.** Route: single-field request to RDR 0002 for a
  normalization clause retaining each atom's authored block. RDR 0002 is
  `Draft`; this is the one field that stays a peer request, because per-atom
  retention is normalization carriage rather than declaration semantics.
- [x] RDR 0002's sparse table/container contract is stable enough to host guard
  atoms.
- [ ] RDR 0007 is the normative home of the guard seam, the domain rule, and
  the `guard_unevaluable` payload this RDR cites. RDR 0007 is `Final`, so the
  *specification* is settled — but its kernel reshape is unimplemented: the
  shipped kernel still carries `internal/resolve/resolve.go::Row.Guard` as a
  string and `resolve.go::GuardEvaluator.Evaluate(guard string, view TagSet)`.
  Phase 1 cannot begin against the atom-slice shape until that reshape lands,
  so this item gates implementation sequencing, not lock.

### Minimum Viable Validation

Encode one RDR flow slice and one kata flow slice as normalized candidate rows,
including equality, enum membership, set containment, bounded integer
comparison, existence, and mixed `all`/`unless` guards. Lint must prove one
exhaustive and mutually exclusive scoped row group, then detect one intentional
gap and one intentional overlap that only appear in the multi-dimensional
product, with source rule/context ids in the diagnostic. The fixture must also
include one row group that is domain-exhaustive yet contains a possibly-absent
guard key, and assert **positively**: lint emits the blocking
inability-to-prove finding naming that row and the refusing atom. Asserting only
the absence of a green result does not discharge this — a run that emits nothing
must fail the test. The MVV must also carry a **negative control**: a row group
whose keys are all declared always-present, which lint certifies green, proving
the narrowing is tight rather than blanket.

**Authorability.** Every case above is authorable against this RDR's own tag
declaration model: the `contains` case declares a `set` tag with an element
universe, the possibly-absent-key case declares an optional key, its
always-present negative control declares the same key always-present, the
finite-domain cases declare enum sets and `{min..max}` bounds, and a
single-valued enum tag declares the marker §JD-13 homed here — the field RDR
0006's `graph-single-valued-state` reads, asserted here as declarable and
rejected on a `set` kind. This closes the
gate that previously bound the MVV's schedule to a peer's producer request — the
fields are no longer fixture-invented keys read back as evidence but normative
declarations this document owns, which is what removes the circularity A2's
stamp was once corrected for. The MVV's remaining dependency is sequencing, not
authorability: RDR 0007's kernel reshape must land before Phase 1 builds against
the atom-slice shape.

**Phase gating.** The phases below are not a linear sequence — each names the
assumptions that block it, so no phase is picked up on document order alone:

| Phase | Blocked on | Startable today |
| --- | --- | --- |
| 1 Predicate Model | RDR 0007's kernel reshape (shipped kernel still has `Row.Guard string`) — A13's set-literal spelling and the tag declaration model are stated | Partially — the matrix, the declaration model, and the set-literal canonicalization, not the atom slice |
| 2 Finite-Domain Lint Semantics | A10, A12 (RDR 0006 agreements) — A11, A2 and A8 are closed, the domain producer is this RDR | Yes for the product/proof core; the peer agreements gate integration, not construction |
| 3 Target-Flow Fixture | A14 (per-atom `block` retention, RDR 0002) — A1's harness has run and A7/A9 are closed | Yes — the gating predicates are now authorable |
| 4 Integration With Peer RDRs | A10, A12 (tolerances on RDR 0006's refine) — A8 closed 2026-08-22 by §JD-4 | No |

Building Phase 1 against the *current* `Row.Guard string` shape produces work
that must be redone when RDR 0007's reshape lands; that is a sequencing
decision to make deliberately, not by default.

### Phase 1: Predicate Model

Define the tag-kind/operator compatibility matrix and normalized predicate atom
shape used by resolver and lint.

### Phase 2: Finite-Domain Lint Semantics

Define how enum, boolean, set-universe, and bounded-int domains are converted
into scoped row-group coverage and overlap checks, including
refusal/downgrade behavior for unbounded dimensions and finite products too
large to prove deterministically.

### Phase 3: Target-Flow Fixture

Build the MVV fixture against representative RDR and kata guards and use the
result to confirm or adjust the initial operator set. The Resolve spike already
covered the target-flow subset `eq`, `in`, `lt`, `gte`, and `exists`; the
implementation MVV must add at least one `contains` predicate over a declared
set-valued tag before the full operator vocabulary is accepted.

### Phase 4: Integration With Peer RDRs

Connect predicate diagnostics to RDR 0002 source identities, RDR 0001 exact-one
selection, RDR 0005 CLI output, and RDR 0006 lint authority.

### Day 2 Operations

This RDR creates no persistent resource. Day-2 management belongs to the
transition model artifact owned by RDR 0002.

### New Dependencies

No new third-party dependency is proposed. The predicate grammar and finite
domain checks should be implemented with local Go code unless Resolve proves a
small parsing or set library is necessary.

## Validation

### Testing Strategy

The MVV should become production tests that exercise both runtime predicate
evaluation and lint-time finite-domain reasoning. Done means the same normalized
predicate atoms drive exact-one row selection, overlap detection, and
exhaustiveness proof/refusal without host callbacks or source-order priority.
Resolve evidence for the representative authoring shape lives in
`docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes/`: the fixture
covers profile routing, cap-3 handling, prelock lens sets, cluster eligibility,
and rewind legality with the target-flow subset `eq`, `in`, `lt`, `gte`, and
`exists`. The Stage 6 A1 harness
(`evidence/spikes/iter-2/a1-eval-harness/`) evaluates those rows against
representative tag views and asserts each qualify/prune verdict, and exercises
`lte`, `gt`, and `contains` besides — so the closed vocabulary is decided, not
merely spelled. It is a spike, not a production test: it carries its own TOML-to-
atom encoding because the kernel reshape is unshipped, and it consumes no
declared domain, so it proves evaluation and not exhaustiveness. The
implementation MVV must still exercise `contains` over a **declared** set-valued
tag with an element universe before the full closed operator vocabulary is
accepted.

1. **Scenario**: Evaluate representative RDR and kata rows that use equality,
   membership, set containment, bounded integer comparison, existence, and mixed
   `all`/`unless` guards.
   **Expected**: Exactly one qualifying row resolves for the legal input; a row
   whose `all` predicates match is disabled when its full conjunctive `unless`
   block also matches; zero or multiple qualifying rows become typed refusals.
2. **Scenario**: Lint finite enum, boolean, set-universe, and bounded-int tag
   domains inside one normalized row group with one complete partition, one
   intentional gap visible only in the multi-dimensional product, and one
   intentional overlap visible only in the multi-dimensional product.
   **Expected**: Complete partitions pass; product-level gaps and overlaps fail
   with source rule/context ids. The gap and the overlap are both reported from
   one run over the group — detecting only the first is a failure. The gap
   finding names the selection context, every rule id in the group, and one
   concrete uncovered assignment; the overlap names both contributing rule ids.
3. **Scenario**: Lint an otherwise valid guard over an unbounded integer or
   undeclared finite domain, and lint a declared finite product too large for
   the deterministic proof representation. Include a group with two separately
   unprovable dimensions, and two equal-cardinality products of differing shape
   (few wide dimensions vs. many narrow ones).
   **Expected**: Runtime evaluation remains available, but lint refuses or
   downgrades the exhaustiveness claim for that dimension/product. The
   two-dimension group emits two findings, not one. The over-large refusal
   reports both the computed product cardinality and the published bound. The
   two equal-cardinality products receive the same verdict — a divergence
   refutes A15 and the bound clause is restated before lock.
4. **Scenario**: Parse malformed guard atoms: unknown tag, unknown operator,
   unsupported operator/tag-kind pair, and literal parse mismatch. Include the
   malformed *declarations* the domain/kind agreement clause rejects: a
   `{min..max}` bound on an `enum`, an element universe on a scalar kind, and a
   single-valued marker on a `set` kind or on a kind carrying no finite domain.
   **Expected**: Each failure is rejected before resolution with a predicate
   semantic kind that RDR 0006 can map to a lint finding and RDR 0005 can map to
   the structured CLI gateway. The declaration errors are rejected before
   normalization completes, so a consumer reading the model — including RDR
   0006's `graph-single-valued-state` — never sees a marker its kind cannot
   carry.
5. **Scenario**: Reorder authored rows and guard atoms without changing their
   semantics, including reordering the elements inside an `in` set literal.
   **Expected**: Successful matching and lint findings are unchanged because
   source order is not a selection mechanism; an ambiguous pair remains a
   multiple-match refusal instead of becoming a first-match success.
6. **Scenario**: Evaluate a row group whose guards are exhaustive over their
   declared domains but where one participating row carries a value atom over a
   key that can be absent.
   **Expected**: The value atom is unevaluable rather than false, resolution
   refuses `guard_unevaluable` under RDR 0007's veto, and lint emits the
   blocking inability-to-prove finding naming that row and atom — not merely an
   absent green result — so the green-lint-implies-resolution-succeeds promise
   holds and is observable. Paired negative control: an otherwise identical
   group over always-present keys certifies green. An existence atom over the
   same absent key decides instead of refusing, and this RDR's evaluator is
   never consulted for either.

### Performance Expectations

Resolve measured representative fixture size rather than setting a throughput
target: the spike uses four rows and the target-flow operator subset. The
intended implementation uses local typed comparisons and deterministic symbolic
or bitset-equivalent finite-domain proof over declared domains; no callback
invocation, expression parser, or external engine is part of the hot path. If
implementation later indexes predicates for speed, the optimization must
preserve the normalized atom semantics, scoped product proof, and exact-one
refusal behavior.

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

The research supports a symbolic guard-atom model over declared tag domains, and
the proposed solution keeps host callbacks and free-form expression strings out
of the contract so lint can prove finite-domain coverage and overlap.

**The declared tag domains this model rests on are declared here.** They were
previously booked as a producer request to RDR 0002 (A11, with A7 and A9 as
subfields), and that request had no prospect of being met: RDR 0002's only
normative tag clause requires *provenance*, `value kind` is normative nowhere in
the cluster, and RDR 0002's validation vocabulary carries no type or
value-domain category. The gap was an unhomed contract, not an unmet request.
This RDR now states the tag declaration model as normative contracts, so A11,
A7, A9, and A2 behind them are `Verified` and every input the exhaustiveness
proof consumes is defined in the same document that consumes it. RDR 0007
(`Final`) records the same division independently. No internal contradiction
arises: RDR 0002 keeps provenance, authoring location, and normalization
carriage, and the two documents cite rather than restate each other.

Guard-domain enforcement is resolved rather than open: the kernel decides
presence and existence atoms, so this RDR's evaluator is scoped to value
semantics over a present value (JDR 0001 §D4, a Closed entry).

The strength of the exhaustiveness claim is decided in substance but not yet
agreed across documents. §JD-4 settles that the lint promise narrows and the
runtime veto stands, and this RDR now states that for its own proof. §JD-4
remains open as to which document records it and names RDR 0006's proof; RDR
0006 carries no `guard_unevaluable` narrowing. RDR 0006 was demoted to `Draft`
(2026-08-21) and carries the recording question in its Refinement Context, so
the divergence closes by a scheduled edit on an open peer rather than a
route-back. A8 tracks it, and it must close before lock.

### Assumption Verification

Each record uses an allowed Method label, carries concrete Evidence or an
Evidence-needed line, and has a non-empty "If wrong" consequence. No record uses
`Docs Only`. **Three records were demoted by the critique pass**: A1 and A2 from
`Verified` to `Pending`, and A6 narrowed to labels-only with its reachability
half split out as A12. Verified: A1, A2, A3, A4, A5, A6 (labels), A7, A9, A11,
and A13. Pending: A8, A10, A12, A14, A15 — each with a named plan. None is
`Unverified`. The repeatability pass closed A13 by stating the set-literal
clause it owed and opened A15 on the cardinality bound the too-large clause now
declares. The Stage 6 reconcile closed A1 by running the evaluation harness the
critique pass's demotion demanded.

**The declaration-model rehoming closed four records at once (A11, A7, A9, and
A2 behind them).** Stage 6 had them as BLOCKERs: the MVV was unauthorable
because RDR 0002 declared no domain field. Re-examining *whether RDR 0002 was
the right home* dissolved the blocker rather than deferring it. RDR 0002's only
normative tag clause requires provenance alone; `value kind` is normative
nowhere in the cluster (it appears in RDR 0002's prose schema list and in a
spike fixture, and every peer citing "name, provenance, value kind, accessor
reference" is quoting that prose as if it were a contract); RDR 0002's normative
validation categories are entirely structural, with no type or value-domain
axis; and RDR 0007 (`Final`) states in its own rejected-alternative analysis
that "typed literals and value kinds (bounded integers, set universes) are RDR
0003's declarations". The declaration model is the guard algebra's input
alphabet, and it now lives with the algebra. RDR 0002 retains authoring location
and normalization carriage — including provenance, which its consumers RDR 0004
and RDR 0006 read — and cites this model for meaning. Only A14 stays a request
to RDR 0002, because per-atom `block` retention is carriage, not semantics.

The demotions matter because this section previously reported label hygiene as
verification: every record *had* a Method and an Evidence line, so the audit
passed while A1's cited spike proved nothing about the claim and A2's derivation
had no producer for its inputs. A conforming record is not a verified one.

A1 was stamped on a spike that scans operator *spellings* and cannot evaluate a
predicate. The demotion was correct and the remedy has now run: Stage 6 replaced
that spike with an evaluation harness that decides the fixture's rows against
representative tag views and asserts every qualify/prune verdict, so A1 is
`Verified` on a spike that measures rather than describes. A2's derivation is
sound and its inputs are now declared here, so it stamps as a Derivation over
owned declarations rather than borrowed ones. A3 is an explicit design
decision that rejects inline negation and nested boolean expressions. A4 cites
source-search anchors that resolve now:
`internal/cli/clierr/clierr.go::CLIError`,
`internal/cli/respond/respond.go::Fail`, and
`internal/cli/config/config.go::Load`. A5 rests on the kernel-fixed parsed-atom shape (JDR 0001 §D1, normative in RDR
0007) plus RDR 0002's source-identity contract, and on
`internal/resolve/resolve.go::Row`, which already carries `RuleID` and
`SourceLocator`; A6 relies on peer RDR contracts in RDR 0002 and RDR 0006 for
provenance *labels* only — the write-reachability property its old "If wrong"
implied is now A12. None of these is self-reference. No `Source Search` Evidence
cites this RDR or its artifact directory, and no record cites this RDR's own
spike fixture as evidence for a schema RDR 0002 does not declare.

**A7, A9, and A11 are Verified and no longer block lock.** A11 states the tag
declaration model; A7 states the existence-atom presence-dimension projection on
a derivation recorded in `evidence/research/iter-2-projection-derivation.md`
(RDR 0007 A12 routed the question here and it is answered here); A9 rests on the
element universe the same model declares. The MVV's authorability gate is
correspondingly closed: every required case is writable against declarations
this document owns.

**A8 is Pending and is the one record that still gates lock.** This RDR records
the §JD-4 narrowing for its own proof, but §JD-4 is an open ledger entry naming
RDR 0006's proof and leaving the recording document unassigned, and RDR 0006
carries no `guard_unevaluable` narrowing. Closure needs a §JD-4 disposition, not
an edit here. Recommended arm: record the narrowing here — RDR 0006's clause
fires only on a non-finite dimension, while this case is a fully-finite product
whose participating row can still refuse, and the clause must name the refusing
atom, a vocabulary RDR 0006 does not have (`atom` occurs once there, delegating
atoms to this RDR). RDR 0006 then cites it and reuses `graph-unprovable-coverage`,
whose stated meaning — "Required finite-domain proof unavailable" — already
covers a withheld claim.

**A14 is Pending against RDR 0002 and does not block lock.** It is the single
surviving producer request: a normalization clause retaining each atom's
authored block. Survivable because no peer asserts the opposite and the failure
is caught by the first normalization test.

**A10 and A12 are Pending.** Both join A8 at cluster reconcile — A10 (row-group
division of labour) and A12 (predecessor reachability) are RDR 0006 agreements
this RDR cannot make unilaterally. A13 (canonical set-literal spelling) was the
one record closable on this RDR's own initiative and is now `Verified`: the
repeatability pass stated the clause RDR 0007 assigned here, which matters
because it affects `in`, an operator inside the subset this RDR treats as
proven. A15 (cardinality bound) replaces it as the self-owned open record, and
is discharged by the MVV rather than a peer.

### Scope Verification

The Minimum Viable Validation is in scope for implementation, not deferred, and
every required case is now authorable against this RDR's own tag declaration
model. **Scope grew deliberately at Stage 6**: this RDR now owns the tag
declaration model (value kind, finite domain, optionality, element universe) in
addition to the guard grammar and exhaustiveness semantics. That is one contract,
not two — the declaration model is the alphabet the grammar is written against,
and it was previously homeless rather than owned elsewhere. The Proportionality
note below records why this does not overreach. The
specific proof is a production test fixture that encodes one RDR flow slice and
one kata flow slice as normalized candidate rows, including equality, enum
membership, set containment, bounded integer comparison, existence, mixed
`all`/`unless`, one complete partition, one intentional product-level gap, one
intentional product-level overlap, and refusal/downgrade cases for unbounded or
too-large finite domains.

### Cross-Cutting Concerns

- **Incremental adoption**: guards over unbounded dimensions remain runtime
  evaluable, but lint must refuse or downgrade exhaustiveness for those
  dimensions rather than blocking all use of predicates.
- **Memory management**: set-valued finite domains must use deterministic
  symbolic or bitset-equivalent proof; the RDR explicitly rejects naive powerset
  materialization when the finite product is too large to prove.
- **Canonical-form / determinism**: successful matching and lint findings cannot
  depend on source order. Exact-one selection belongs to RDR 0001, source-row
  identity comes from RDR 0002, CLI envelopes belong to RDR 0005, and blocking
  graph-lint findings belong to RDR 0006.

### Proportionality

The document is right-sized for lock. It owns one independent load-bearing
contract: the guard-predicate grammar and finite-domain exhaustiveness
semantics, together with the tag declaration model that grammar quantifies over.
Neighboring surfaces are explicitly delegated to peer RDRs rather than locked
here: sparse table shape and normalization carriage to RDR 0002, exact-one
resolution to RDR 0001, CLI envelope mapping to RDR 0005, graph-lint authority
to RDR 0006, and the kernel seam, domain rule, and `guard_unevaluable` payload
to RDR 0007.

**Why absorbing the declaration model is not scope creep.** Coverage and overlap
are proved over a product of declared domains; a grammar that cannot say what a
tag may hold cannot prove anything about it. The alternative placements were
tested and rejected: RDR 0002 has no normative tag-type vocabulary (its one tag
clause requires provenance alone), no type or value-domain validation category,
and a charter — "the on-disk sparse representation and the normalized
expanded-table view" — that the declaration's *meaning* falls outside of; a new
RDR would leave one tag's declaration split across two documents and add an edge
from all of 0002/0003/0006/0007. RDR 0007 (`Final`) independently records the
same division: value kinds and set universes "are RDR 0003's declarations". The
model was homeless, not homed elsewhere, so this is a gap closed rather than a
boundary crossed — and the dependency direction now matches the rest of the
cluster, where RDR 0002's own Approach already says "RDR 0003 owns the fixed
predicate operators".

The `large` Profile still matches because this RDR locks a grammar and
finite-domain proof semantics. The required large-profile lenses ran
(`grounding`, `3amigo`, and `critique`); every finding they raised is either
resolved in the live text above or carried as a named obligation in the Minimum
Viable Validation and Testing Strategy. The re-entry re-runs that row against
the revised text: `grounding` iteration 2 is complete, and its `authority`,
`disposition`, and `trace` mini-check tables are carried above.

## References

- RDR 0001, Resolution Kernel Contract.
- RDR 0002, Transition Table As Reviewable Data.
- RDR 0006, Graph Lint Authority and Guarantees.
- RDR 0007, Guard Predicate Totality — normative home of the guard seam, the
  domain rule, and the `guard_unevaluable` payload.
- JDR 0001, Resolve Kernel Seam — §D1 (parsed-atom row shape), §D4 (kernel
  enforces the guard domain), §JD-4 (the lint promise narrows — substance
  decided, recording document still open, head clause names RDR 0006).
- `docs/cli-output-contract.md`.
- Resource index: `.rdr/resources.md`.
- Seed prior: `../state-machines/BUILD-SEEDS.md`, especially the guard
  expression and transition-table seeds.
- Transition model prior: "The transition model — inputs, outputs, error
  conditions."
- Prior-art corpus: `../state-machines` audits, evals, contrasts, and checked
  repositories for `transitions`, stateless/qmuntal-stateless, SCXML/scxmlcc,
  Sismic, StateSmith, and Statewright.
