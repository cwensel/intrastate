# Recommendation 0003: Guard Predicate Exhaustiveness

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-06-19
- **Status**: Implemented
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

## Critical Assumptions

- **A1 The target RDR and kata flows fit a closed typed predicate vocabulary.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: Evaluation harness. Command:
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
    governs exhaustiveness (A2) rather than evaluation. That independence is what
    keeps the record independent of the declaration model this RDR owns, and
    it is also this record's limit: the harness predates the model, uses its own
    TOML-to-atom encoding rather than the shipped shape, and runs against a
    fixture declaring none of the five fields. So `contains` is proven to
    *evaluate* over a set-valued tag and is **not** proven to carry an
    exhaustiveness claim over a *declared* element universe — Phase 3's gate, and
    the one operator whose closed-vocabulary acceptance still rests on a
    condition discharged at implementation. The `contains` cell is
    evaluated over a declared set **kind**; the element *universe* the
    exhaustiveness claim needs is a separate fact, declared by A11's model and
    recorded by A9.
  - **If wrong**: The fixed operator set is too small, and authors will need an
    expression grammar or host predicates that weaken static lint.
- **A2 Every exhaustiveness claim can be reduced to scoped finite declared
  domains.**
  - **Status**: Verified
  - **Method**: Derivation
  - **Evidence**: The derivation below is sound, and its inputs are declared
    by this RDR's own tag declaration model (A11). It consumes "finite domains for the guard
    dimensions" — enum value sets, declared set-element universes, and bounded
    integer `{min..max}` ranges — and the `Normative Contracts` declaration
    clauses make each of those declarable, with domain/kind agreement enforced
    before normalization. The domain is a normative declaration this document
    owns, not a fixture-invented key read back as evidence, so the derivation
    stamps without circularity.
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
  - **Evidence**: `internal/cli/clierr/clierr.go::CLIError` carries stable `Code`, human `Message`, optional `Param`, `Detail`, `Hint`, and exit-code `Group`; `internal/cli/respond/respond.go::Fail` emits the envelope in text/json modes; `internal/cli/config/config.go::Load` already demonstrates stable parse/read error codes. This RDR owns predicate semantic kinds such as unknown operator, type mismatch, literal parse failure, literal-outside-declared-domain, declaration/kind disagreement, and unsupported operator/tag-kind pairing; those load-time rejections surface through RDR 0002's load categories (`malformed tag declaration`, `malformed predicate atom` — JDR 0001 §D7(iii)); the per-row/per-atom `guard_unevaluable` payload is RDR 0007's under JDR 0001 §D4; RDR 0006 owns graph-lint finding codes for non-exhaustive and overlapping row groups; RDR 0005 owns the CLI command/envelope mapping onto the existing gateway (JDR 0001 §JD-8).
  - **If wrong**: This RDR or RDR 0005 must add a separate user-facing error
    contract before implementation.
- **A5 The normalized predicate representation can retain source identity for
  actionable diagnostics.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: The normalized atom shape
    is fixed at the kernel seam by JDR 0001 §D1
    (`docs/jdr/0001-resolve-kernel-seam.md::D1`); RDR 0007 is the landing
    document and states it normatively — a row carries parsed atoms "key,
    operator token, literal, block ∈ {all, unless}"
    (`0007::Normative Contracts`, the SEAM parsed-atom clause), not an opaque predicate
    string, so there is no reconstruction step that could lose identity. RDR
    0007 is `Final`, not yet implemented: the shipped kernel still carries
    `Row.Guard string` (`internal/resolve/resolve.go::Row.Guard`), and 0007's A3 scopes
    the reshape. The source identity this assumption is about already ships —
    `internal/resolve/resolve.go::Row` carries `RuleID` and `SourceLocator`,
    commented as "the source identity RDR 0002 requires". RDR 0002 `Normative
    Contracts` require each normalized candidate row to "retain its source rule
    id and source locator"
    (`0002::Normative Contracts`, the candidate-row source-identity clause), and its
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
  - **Note**: This RDR's owned-tag clause quantifies over *every reachable
    predecessor*, which needs a predecessor relation, not a label. That half is
    **A12**. The "can refuse" clause reads the declared optionality field
    alone, so A12 carries the owned-tag clause only.
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
    The optionality input this clause reads is declared by this RDR's own
    tag declaration model (A11). RDR 0007 A12 routes the question here and it
    is answered here.
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
    RDR 0006 (`Final`, locked 2026-08-23) states it in fenced text: a withheld
    claim "MUST be emitted as `graph-unprovable-coverage`, naming the
    participating row and the refusing atom. This RDR cites that clause and
    MUST NOT restate it, mint a second code for it, or carry it in a
    non-blocking tier" (`0006::Normative Contracts`, the narrowing-citation
    clause); `graph-unprovable-coverage` sits in its blocking finding-code
    table. Gate evidence:
    `docs/rdr/cluster-reconcile/0003-0006-0007/reconcile-report.md`.
  - **Consequent duty on RDR 0006 — discharged.** This RDR's clause requires
    the finding to name the refusing atom, and RDR 0006's finding contract
    carries it: "A finding attributed to one guard atom MUST carry that atom's
    `Key`, `Operator`, `Literal`, and `Block`" (`0006::Normative Contracts`,
    the finding-record clause).
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
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: The group is defined here (rows sharing one source state and
    one recognized outcome — the set RDR 0001 resolves exact-one over), and RDR
    0006 (`Final`) reads the same division of labour in fenced text: "Coverage,
    overlap, and withholding MUST be decided per scoped row group as RDR 0003
    defines it … This RDR MUST NOT define a second grouping predicate; it
    supplies only which selection contexts are reachable"
    (`0006::Normative Contracts`, the row-group clause), and its Technical
    Design states "Row groups are RDR 0003's" and closes this record by name.
    Ledger: `docs/rdr/cluster-reconcile/0002-0009/iter-2/pairwise-0003-0006.md`
    F4. The grouping predicate therefore has one home and no cycle.
  - **If wrong**: the two documents group rows differently, so a gap proved
    absent in one grouping is present in the other, and the exhaustiveness claim
    means different things on each side of the seam.
- **A11 The tag declaration model carries a finite domain, so any
  exhaustiveness claim has a producer.**
  - **Status**: Verified
  - **Method**: Design Decision
  - **Evidence**: This RDR **owns** the tag declaration model. The `Normative
    Contracts` declaration clauses state the value
    kinds, make a finite domain declarable for every kind an exhaustiveness
    claim can range over (enum value set, `bool` by construction, `int`
    `{min..max}`, `set` element universe), and require domain/kind agreement
    with rejection before normalization. A2's derivation therefore consumes
    inputs this document defines.
    **Why the home is here, not RDR 0002.** The declaration model is the guard
    algebra's input alphabet and belongs with the algebra. RDR 0002's own
    normative tag vocabulary is provenance plus the wire-key spellings; it
    cites this RDR for what the keys mean ("RDR 0003 is the normative home of
    that model … MUST NOT be restated here", `0002::Normative Contracts`, the
    type-model clause), and RDR 0007 (`Final`) states the same division from
    its side — "typed literals and value kinds (bounded integers, set
    universes) are RDR 0003's declarations". RDR 0002 keeps authoring location
    and normalization carriage. The rejected alternative is homing the five
    fields' meaning in RDR 0002, which would make that document the de facto
    owner of a type system it has no validation category or charter for.
  - **If wrong**: no row group is ever exhaustiveness-eligible, every group takes
    the blocking inability-to-prove outcome, and this RDR's central guarantee is
    unreachable in practice while remaining true in theory.
- **A12 Predecessor reachability is decidable from the declared model, so
  "every reachable predecessor" is a total test.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: The owned-tag clause quantifies over every reachable
    predecessor; the "can refuse" clause does not — it reads the declared
    optionality field alone, which is what makes it total, so this record
    gates the owned-tag clause only. A6 establishes that provenance *labels*
    reach lint; the predecessor relation is RDR 0006's graph property, and RDR
    0006 (`Final`) publishes it: the **owned-state graph** — nodes are abstract
    owned-states, the root is the declared initial owned state, edges are
    normalized rows (escape rows as self-loops) — with "a row's *reachable
    predecessors* are the rows on any root-to-source path", computed as "a
    syntactic dataflow over declarations, writes, and clears … so it is total,
    and it over-approximates the runtime" (`0006::Load-Bearing Decisions`,
    "Reachability relation"). It also answers which procedure governs each
    clause: reachability for the owned-tag check (`graph-owned-before-write`),
    and "a different predicate from RDR 0003's 'can refuse' test, which is
    decided over the optionality field alone and never consults this
    relation". Where the path wording and RDR 0006's node form differ, the
    node form governs ("Owned-set-before-match is decided over nodes, not
    paths"); the difference is in the false-positive direction only, so no
    model RDR 0006 accepts is one this clause rejects. The owned-tag clause
    below cites that relation rather than defining one. Ledger:
    `docs/rdr/cluster-reconcile/0002-0009/iter-2/pairwise-0003-0006.md` F5.
  - **If wrong**: the owned-tag clause is unimplementable as written — lint
    cannot decide whether an owned tag is set before match.
- **A13 The canonical byte spelling of a set literal is fixed by this RDR for
  both `in` and `contains`.**
  - **Status**: Verified
  - **Method**: Design Decision
  - **Evidence**: RDR 0007 (`0007::Normative Contracts`, the literal-spelling
    condition the totality claim rests on)
    assigns the canonical spelling duty to this RDR and names `in` explicitly.
    A9 books only the *element universe* for `contains`, which is a different
    thing from the literal's encoding: `in` takes a typed literal **set** and is
    already inside the spike/MVV operator subset, so the gap affects an operator
    this RDR claims as proven, not just the deferred one. The duty is
    discharged: the `Normative Contracts` set-literal clause fixes the spelling
    as an unordered, duplicate-free typed set canonicalized before it enters the
    identity tuple, with repeated elements rejected at parse. The rejected
    alternative is an ordered slice preserving authored order, which would make
    `["mid","large"]` and `["large","mid"]` distinct atoms and break MVV
    Scenario 6's reorder invariant.
  - **If wrong**: two implementations spell the same set literal differently, so
    a fixture that parses under one build fails under another, and the identity
    tuple's `literal` component is not stable across producers.
- **A14 Normalization preserves each atom's `block`, so `all` and `unless` stay
  distinguishable downstream.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0002 states the carriage clause normatively. Its
    normalization contract reads "Normalization MUST combine the match atoms and
    both guard blocks into one candidate-row predicate set before ambiguity
    checks, and each atom in that set MUST retain the key, operator token,
    literal, and the block it was authored in — a three-valued domain, `match`,
    `all`, or `unless`", followed by an explicit anti-collapse obligation naming
    this RDR as the reason: "The block is load-bearing, not merely carried: it
    is the key the handoff routes on (JDR 0001 §D6, below), RDR 0003's atom
    identity tuple and `unless` semantics read it downstream, and RDR 0006
    serializes it. Normalization MUST NOT fold `unless` atoms into `all` or
    `match` atoms into either guard block" (`0002::Normative Contracts`, the
    normalization block). The third `match` member is RDR 0002's own widening on
    JDR 0001 §D6's authority — this RDR and RDR 0007 spell the two-valued guard
    domain §D1 fixes, and §D6 routes the handoff on the authored block — so the
    guard-side carriage this record depends on is unaffected. The block also survives into the reviewable surface:
    the dump field list carries "predicate atoms (with block)" and the dump's
    within-row total order sorts atoms by "(key, block, operator token,
    literal)", so `block` is load-bearing in the ordering rather than merely
    present. Corroborated by the Stage 6 A1 harness, which carries `block` per
    atom through its own encoding and shows `unless` behaving as a conjunctive
    exclusion block distinct from `all`
    (`evidence/spikes/iter-2/a1-eval-harness/`) — this RDR's encoding, so it
    corroborates rather than discharges; RDR 0002's clause is what discharges.
    MVV Scenario 7 (each atom still reports its authored block after
    normalization) is the test that keeps it honest, not the plan that closes
    it.
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
  - **Plan**: MVV Scenario 3 lints "a declared finite product too large for
    the deterministic proof representation", asserts the refusal reports both
    figures, and compares two equal-cardinality products of differing shape
    (few wide dimensions vs. many narrow ones) for the same verdict. **That
    comparison alone cannot close this record** — reading the bound `B` from
    the implementation makes the two verdicts comparable but takes the oracle
    from the system under test, so it establishes consistency and never
    adequacy. Scenario 3 therefore also measures whether the proof
    representation actually completes for each shape and compares that against
    the verdict the bound gave; A15 closes only when predicted and observed
    provability agree, and is refuted by a divergence in either direction.
    Not runnable at draft time — it needs finite products built and measured.
    Survivable to implementation because the "If wrong" is bounded: a
    divergence adds a second bound term and restates one clause; it does not
    disturb the grammar, the identity tuple, or the coverage derivation.
  - **If wrong**: the published bound does not predict provability, so two
    conforming implementations disagree on the same model and the
    model-independence the clause promises is unmet.
- **A16 RDR 0002's authoring surface carries the single-valued marker, so the
  field §JD-13 homed here is writable by an author.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0002 (`Final`) carries the field on both surfaces, each
    citing this RDR for meaning rather than restating it. Its `[tags.<tag>]`
    schema enumeration reads "the RDR 0003 type model spelled as the wire keys
    `kind`, `domain`, `min`, `max`, `elements`, `single_valued`, `required`"
    (`0002::Technical Design`, the "Tag declarations" schema part), and its
    normative type-model clause reads "A tag declaration also carries its
    **type model**, spelled as the wire keys `kind`, `domain`, `min`, `max`,
    `elements`, `single_valued`, and `required` (JDR 0001 §D7(iii) …). RDR
    0003 is the normative home of that model" (`0002::Normative Contracts`,
    the type-model clause). The same clause states the authoring location
    normatively ("under `[tags.<tag>]`, beside `provenance`, with no accessor
    reference") and the carriage obligation ("Normalization MUST carry every
    declared field through to the normalized model without loss, so lint (RDR
    0006) and the guard proof (RDR 0003) read the same declaration the author
    wrote"). The marker is therefore writable by an author and reaches RDR
    0006's lint input, which gates `graph-single-valued-state` on it. This
    record matters beyond one lint code: the single-value-operator projection
    clause makes the marker a precondition for projecting any
    `eq`/`in`/comparison atom, so an unwritable marker would leave MVV Scenario
    2's passing partition and Scenario 8's green negative control unauthorable
    and the MVV with no positive proof outcome.
  - **Scope of this record**: it establishes that the field is *writable and
    carried*, not that the six target-flow dimensions are correctly marked —
    that authoring-coverage half is A21, measured at MVV Scenario 4.
  - **If wrong**: the single-valued marker is a declared-but-unwritable field,
    `graph-single-valued-state` keeps an unhomed producer, and no guard using a
    single-value operator can carry an exhaustiveness claim — the whole proof
    surface, not one lint code.
- **A17 Escape-row overlap is checked among escape rows for one failure class,
  not between an escape row and a guarded row.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: The two-population reading follows from the runtime: JDR
    0001 §D2's gate-then-count selection (RDR 0002 defers to it — "the only
    successful transition selection is exact-one surviving non-escape row
    after the kernel's gate", `0002::Load-Bearing Decisions`; a modeled escape
    fires only when exactly one escape row for that failure class survives)
    and the shipped kernel
    (`internal/resolve/resolve.go::Resolve` step 2 builds candidates from
    **non-escape** rows; escape rows are consulted only to rescue `no_match` or
    `ambiguous_match`, matched membership-wise by `resolve.go::Row.rescues`).
    An escape row overlapping a guarded row is therefore never a runtime
    ambiguity, and reporting it would fail a fixture the runtime accepts; two
    escape rows matching one failure class are a genuine ambiguity. The cluster
    confirmed that reading on both sides of the seam: JDR 0001 §JD-14 was
    **corrected 2026-08-23** ("overlap is checked in **two populations** …
    never escape-vs-ordinary … This entry's decision is amended to that
    reading; 0003 A17 closes on it"), and RDR 0006 (`Final`) states it in
    fenced text — "Escape rows MUST participate in the coverage union and MUST
    be overlap-checked in one population per declared failure class, never
    against ordinary rows, as RDR 0003's escape-row clause states"
    (`0006::Normative Contracts`, the escape-row clause; invariant 3 in its
    `Technical Design`). The coverage half of §JD-14 stands unchanged.
  - **If wrong**: this RDR's two-population overlap clause reverts to a single
    population and RDR 0006's invariant 3 exemption is removed; the coverage
    clause and the union identity are untouched either way.

- **A18 A conforming-view check has a producer, so conformance is enforced
  rather than assumed.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence needed**: every lint claim in this RDR is conditional on the
    evaluation view conforming to the declared model. The **owned half** has a
    producer: RDR 0006 (`Final`) requires "An owned key declared always-present
    MUST be held in every reachable owned-state node; a violation is
    `graph-always-present-owned`" and its invariant 5 rejects a write assigning
    a single-valued tag two values — while stating that it "discharges only the
    owned half of RDR 0003's conformance premise" and MUST NOT extend the check
    to observed or recognized keys (`0006::Normative Contracts`, the
    always-present-owned clause). The **view-level half** — an observed or
    recognized key declared always-present but absent at runtime — is held by
    no component: the shipped assembly (`internal/resolve/resolve.go::assemble`)
    merges owned, observed, and recognized tags into a `TagSet` and reads no
    declaration, `conform` occurs nowhere in the kernel, and RDR 0007
    (`Final`) states no conformance obligation.
  - **Plan**: JDR 0001 **§JD-18** (open; siblings 0003 and 0007) decides which
    document states the view-level check and where it runs — view assembly
    rejecting a non-conforming view, or exposing it as a typed refusal. This
    record closes when §JD-18 is answered and the answer is cited here.
    Survivable to that venue because every lint claim here is explicitly
    conditional on a conforming view, so no clause asserts an enforcement it
    cannot cite, and the MVV's green cases declare their own conforming
    fixtures, so the property under proof is reachable without the runtime
    enforcer.
  - **If wrong**: a key declared always-present is absent at runtime, lint proved
    coverage over a product carrying no `{absent}` assignment for it, and the
    resolution refuses `guard_unevaluable` on a group certified green.
- **A19 A bare escape row's coverage closure is observable in lint output.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: This RDR requires a group whose coverage is closed by a
    guard-atom-free escape row to say so in its verdict, so a catch-all cannot
    silently stand in for a proof; the carrier is RDR 0006's to choose, and it
    chose a code: "A group whose coverage is closed by a bare escape row MUST
    emit `graph-coverage-closed-by-escape` naming that row; a bare green MUST
    NOT satisfy this clause" (`0006::Normative Contracts`, the escape-row
    clause), listed in its finding-code table at the `info` tier — observable
    without being blocking, which is the distinction this clause asks for.
  - **If wrong**: a bare escape row remains a one-line blocking-free way to close
    any coverage gap, and the exhaustiveness guarantee erodes silently as authors
    reach for it.

- **A20 The atom-naming half of the withheld-claim clause gets a carrier on both
  surfaces.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence needed**: the withheld-claim clause requires the artifact to name
    the participating row **and the atom that can refuse**. The **lint half is
    closed**: RDR 0006's finding contract carries the atom — "A finding
    attributed to one guard atom MUST carry that atom's `Key`, `Operator`,
    `Literal`, and `Block`" (`0006::Normative Contracts`, the finding-record
    clause). The **runtime half is open**: the shipped refusal carries one
    opaque `Guard string` plus `Rows []RowRef`
    (`internal/resolve/resolve.go::Refusal`), so the atom is not addressable at
    runtime, and RDR 0007 (`Final`) states no atom carrier.
  - **Plan**: the runtime carrier is JDR 0001 **§JD-18**'s subject (open;
    siblings 0003 and 0007), alongside A18's view-level enforcer. This record
    closes when §JD-18 lands an atom-addressable refusal payload and it is
    cited here. Survivable to that venue: MVV Scenario 8's other assertions
    (the withheld claim, the `unless` variant, the escape-row non-rescue) all
    stand without the runtime atom field, so the scenario weakens rather than
    disappears if it slips.
  - **If wrong**: MVV Scenario 8 — this RDR's flagship false-green defense —
    cannot be written as specified, and is either deleted or weakened to assert
    only the row, which is the assertion the clause was strengthened to replace.

- **A21 Requiring `single_valued` for a provable `eq`/`in` atom leaves the
  target flows authorable.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence needed**: the operator/kind agreement clause now makes the
    single-valued marker a **precondition for provability**, not merely a product-size
    optimization: an `eq`/`in`/comparison atom over a tag not declared
    single-valued has no projection and takes the blocking outcome. This is a
    strictly stronger demand on the authoring surface than A16 was written
    against — A16 asks that the marker be *writable*, this asks that every guard
    dimension in the two target flows actually *carry* it. The desk trace shows
    the cost of getting it wrong: the representative fixture declares the marker
    nowhere, so all four of its rules become unprovable rather than merely
    gapped. Unverified today because no fixture declares the marker at all.
  - **Plan**: MVV Scenario 4 already covers single-valued cases — extend it to
    assert the negative control (an unmarked tag under `eq` yields the blocking
    inability-to-prove finding naming the dimension, not a coverage gap), and
    have Phase 3's fixture declare `single_valued` on `profile`,
    `prelock_iterations`, `cluster_eligible`, `stage`, `lens` and
    `rewind_target`. If any target-flow dimension turns out to need co-occurring
    values, it is a `set` tag under `contains`, not an unmarked enum — record
    which, since that changes RDR 0002's authoring request.
  - **If wrong**: the marker is required on essentially every guard tag, which
    makes it a de-facto default rather than an opt-in, and the declaration model
    should invert it — declare *multi-valued* explicitly and let the absent
    marker mean single-valued. That inversion contradicts the conservative-default
    reasoning the optionality clause is built on, so it is a real fork, not a
    tweak: it would be decided here and carried to RDR 0002 and RDR 0006.
    Not runnable at draft time — it asks whether the two target flows' guard
    dimensions can all carry the marker, a fact about fixtures rather than
    about this contract. Survivable to implementation because the consequence
    is a *default inversion*, not a defect in the proof: the projection clause,
    the assignment-count table, and the coverage derivation all stand under
    either polarity — only which declaration an author must write changes.
    Nothing sequences ahead of it (A16 closed), but a Phase 3 result that
    forces the inversion is a fork decided here and carried to RDR 0002 and
    RDR 0006, so it is worth measuring early in Phase 3 rather than late.

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
| `eq` | `enum`, `bool`, `int`, `scalar` | one typed scalar | Narrows the tag domain to one value. |
| `in` | `enum`, `bool`, `int`, `scalar` | non-empty typed scalar set | Narrows the tag domain to the listed values. |
| `lt`, `lte`, `gt`, `gte` | `int` | one typed integer | Narrows a bounded integer domain by comparison; remains runtime-only for an unbounded integer. |
| `exists` | any kind, provided the tag is declared optional (a key declared always-present contributes no `{absent}` assignment, so an `exists` atom over it is well-formed but vacuous — lint reports it as such rather than rejecting it) | boolean | Tests presence or absence, not value equality. Decided by the kernel from presence alone (JDR 0001 §D4); it never reaches this RDR's evaluator. |
| `contains` | `set` with a declared element universe | non-empty typed element set | Narrows the set-valued domain to assignments containing every listed element. |

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
normatively in RDR 0007 (`0007::Normative Contracts`, the SEAM parsed-atom
clause): four fields — `Key`, `Operator`, `Literal`, `Block`. This RDR cites that shape
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
this RDR does not read one back from it. RDR 0006 reads the same division —
it "MUST NOT define a second grouping predicate; it supplies only which
selection contexts are reachable" (`0006::Normative Contracts`, the row-group
clause; A10) — so the grouping predicate has one home and no cycle.

The exhaustiveness claim is **default-on for every scoped row group whose
participating dimensions are all finitely declared**: an opt-in flag would let
the guarantee be silently skipped exactly where it matters. RDR 0006 gates its
blocking finding on the same default (`0006::Normative Contracts`, the
exhaustiveness clause — "default-on … never an opt-in annotation").

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
and names this RDR the recording document; RDR 0006 cites the clause and mints
no code (A8). It is recorded here because the constraint binds what this RDR's
own finite-domain product may claim; RDR 0006's exhaustiveness clause covers
only the non-finite-dimension refusal, so the two triggers differ and both
survive.

Provenance affects lint: recognized tags are fresh event inputs, observed
tags are re-read before matching, and owned tags must have a reachable
predecessor write before a row may match them.

#### Normative Contracts

**C1**
```normative
A guard predicate MUST be a symbolic atom over a declared tag, not a host
language callback and not a free-form expression string.
```

**C2**
```normative
This RDR owns the **tag declaration model** — the typed alphabet every guard
atom is written against. A tag declaration MUST carry a value kind, and — as its
kind admits, per the agreement clause below — MAY carry a finite domain, an
optionality marker, a single-valued marker, and (for set-valued kinds) an
element universe.

The value kinds are exactly five, spelled with these tokens wherever a kind is
named: in a declaration, in the operator/kind matrix, and in a diagnostic —
`enum`, `bool`, `int`, `set`, and `scalar`. The `scalar` kind is the opaque
scalar: a typed value compared only by equality and membership. It carries **no**
finite domain and can never bear an exhaustiveness claim, so a guard dimension
over a `scalar` takes the blocking inability-to-prove outcome, and a `scalar`
declaration carrying a finite domain, an element universe, or a single-valued
marker is a declaration error.

RDR 0002 owns where a declaration is authored and how it is carried through
normalization; this RDR owns what a declaration means. Neither document restates
the other (JDR 0001 P6).
```

**C3**
```normative
A finite domain MUST be declarable for any kind an exhaustiveness claim can
range over: an `enum` declares its value set, a `bool` is finite by
construction, an `int` declares a `{min..max}` bound — notation, not wire
spelling: RDR 0002 spells the authored form as the two wire keys `min` and
`max` (`0002::Normative Contracts`, the type-model clause; JDR 0001 §D7(iii)),
and this RDR fixes their meaning — both endpoints are **inclusive**, so
`{0..3}` has cardinality 4 — and a `set` declares the
element universe its members are drawn from. A declaration carrying no finite
domain is well-formed — the tag remains runtime-evaluable — but a guard
dimension over it cannot carry an exhaustiveness claim, and lint MUST take the
blocking inability-to-prove outcome for that dimension.
```

**C4**
```normative
A tag declaration MUST be able to state whether the key may be absent. A
declaration carrying **no optionality marker declares the key optional** — the
conservative default, since assuming always-present would let lint certify a
group green on an undeclared property. (Symmetric with the domain rule above:
leaving presence undeclared is not an opt-out from the withholding it implies.)
A key declared always-present MUST NOT be absent from a conforming evaluation
view; a key declared optional MAY be. Presence is a declared property of the tag, not
an observation of one view: lint decides `exists` projection and the
"can refuse" narrowing from this declaration, never from a runtime trace. The
kernel's runtime presence decision (JDR 0001 §D4, normative in RDR 0007) is a
separate question over a single view and is unaffected by this clause.

An evaluation view **conforms** to the declared model when every always-present
key is present in it and every single-valued tag holds at most one of its
declared domain values. Conformance is the premise every lint claim in this RDR
is conditional on: a green exhaustiveness result asserts coverage over
conforming views only. A non-conforming view is a defect in the model or its
producer, not an input this RDR's evaluator interprets — the evaluator never
reads the tag view (see the evaluator-scope clause above), so the check belongs
to the kernel's view assembly, and this RDR states only what conformance means
and that its claims presuppose it. Whether a non-conforming view is refused at
assembly or reported as a lint finding over the declared model is RDR 0007's and
RDR 0006's respectively; this RDR MUST NOT be read as promising anything about
a view that violates the declarations.

**Conformance is enforced for owned keys and unhomed for the rest — a booked
gap, not a silence.** RDR 0006 enforces the owned half over the declared
model: an owned key declared always-present must be held in every reachable
owned-state node (`graph-always-present-owned`), and no row's write block may
assign a single-valued tag two values (its invariant 5) — and it states that
this "discharges only the owned half of RDR 0003's conformance premise"
(`0006::Normative Contracts`, the always-present-owned clause). The view-level
half — an observed or recognized key declared always-present but absent at
runtime — has no enforcer: the shipped view assembly
(`internal/resolve/resolve.go::assemble`) merges owned, observed, and
recognized tags into a `TagSet` without reading the declaration model, and the
string `conform` appears nowhere in the kernel. Such a key yields
`guard_unevaluable` over a product built with no `{absent}` assignment — green
lint, refused resolution, with this RDR's answer being that the view was out
of scope. Naming the owner is not the same as the check existing: JDR 0001
§JD-18 decides the producer, and **A18** books it.
```

**C5**
```normative
A tag declaration MUST be able to state that the tag is **single-valued**: at
most one of its declared domain values holds in any conforming evaluation view.
The marker is meaningful only for a kind carrying a finite domain, and a
`set`-valued kind MUST NOT carry it. Single-valuedness is a declared property of
the tag, not an observation of one view. It is what licenses a consumer to treat
the tag's domain as a partition — mutually exclusive values whose coverage the
scoped product can be grouped over — rather than as independent dimensions.

**Effect on the scoped product, stated so the verdict is computable**: a
single-valued tag contributes **one dimension of `|domain|` assignments** (for
enum `{a,b,c}`: three), because exactly one value holds per view. Absent the
marker, a tag whose values could co-occur contributes one **independent boolean
dimension per value** (`2^|domain|`, minus nothing — the model does not assume
at least one holds). This is the whole operational content of the marker: it is
the difference between a product a lint can partition and one it must treat as a
power set, and it is why the domain/kind agreement clause rejects it on a `set`
kind, whose values co-occur by construction.

RDR 0006's grouping-dependent lint findings read this field from the declaration
and MUST NOT infer it from a tag's name, its value spelling, or a fixture (JDR
0001 §JD-13). **Shape note**: RDR 0006's invariant 5 constrains *writes* —
"no row's write block may assign a single-valued tag two values", decided per
row, syntactically (`0006::Technical Design`, invariant 5) — while this marker
is a property of the *evaluation view*. The two are compatible, not the same
predicate: a write producing two values is one way a view stops conforming,
and RDR 0006 reconciles its invariant as "the model-level half of RDR 0003's
single-valued conformance conjunct".
```

**C6**
```normative
A declared domain MUST agree with its value kind. Each kind admits exactly
these fields, and any other combination is a declaration error that MUST be
rejected before normalization completes:

| Kind | Finite domain | Element universe | Single-valued marker |
| --- | --- | --- | --- |
| `enum` | its value set (required to claim exhaustiveness) | no | yes |
| `bool` | finite by construction | no | yes |
| `int` | `{min..max}` bound (required to claim exhaustiveness) | no | yes |
| `set` | its element universe | yes (required to claim exhaustiveness) | **no** — values co-occur by construction |
| `scalar` | **no** | no | **no** — it carries no finite domain to partition |

So an `{min..max}` bound on an `enum`, an element universe on any kind but
`set`, or a single-valued marker on a `set` or a `scalar`, is each a declaration
error. A value appearing in a guard
literal that lies outside its tag's declared domain MUST likewise be rejected —
this is the predicate semantic kind **literal-outside-declared-domain**, distinct
from a literal parse failure (the literal parses fine; it is simply not in the
domain) — so an unsatisfiable atom is a load-time error rather than a
silently-never-matching row.

The two rejections in this clause sit at **different phase boundaries, and each
names its owner**: a malformed *declaration* (domain disagreeing with its kind)
is rejected by the declaration loader before normalization completes, since
normalization carries declarations it must first be able to trust; a
*literal-outside-domain* is rejected by guard parsing after the declarations
have loaded and before rows are yielded. Both are this RDR's rejection rules
carried by RDR 0002's load categories (JDR 0001 §D7(iii)): the malformed
declaration is a `malformed tag declaration`, the out-of-domain literal a
`malformed predicate atom`. RDR 0006 mints nothing for either; RDR 0005 maps
both onto the envelope under JDR 0001 §JD-8. The predicate semantic kinds this
RDR owns (A4) name *which* rule fired within those categories.
```

**C7**
```normative
The initial operator vocabulary MUST be closed and typed: equality,
membership, bounded integer comparison, existence, and set containment. Unknown
operators MUST be rejected during parse or lint before resolution.
```

**C8**
```normative
Each operator MUST declare which tag value kinds it accepts. A predicate whose
literal cannot be parsed as the declared tag kind MUST be rejected before
resolution.

**`eq`, `in`, and the integer comparisons are single-value operators**, and how
a value atom projects onto its dimension follows from that. Each is defined over
a tag holding exactly one value — the matrix's "narrows the tag domain to one
value" / "to the listed values" — so an atom denotes the assignments in which
the tag's single held value is the literal (`eq`), is a member of the literal set
(`in`), or satisfies the comparison. `contains` is the operator for a tag whose
values co-occur, and it projects the other way: the assignments whose held set
contains every listed element.

Consequently a value atom over a dimension the model does **not** declare
single-valued is not a differently-projecting atom — it is one lint **cannot
project at all**, because the tag has no single held value for the operator to
narrow. Such an atom MUST take the blocking inability-to-prove outcome for that
dimension, exactly as a dimension with no finite declared domain does. Lint MUST
NOT silently pick a reading — resolving it as "the literal is among the held
values", "the held set equals the literal", or "the held set is contained in the
literal" yields different union cardinalities and different overlap verdicts on
the same model, which the published-bound clause forbids. An author who wants
`eq`/`in` proven over a tag declares it single-valued; an author who means
co-occurring membership declares the tag `set` and uses `contains`.
```

**C9**
```normative
This RDR's guard evaluator MUST decide value semantics over a present value
only. Key presence, existence atoms, absent-key unevaluability, and the
combination of per-atom verdicts belong to the kernel (JDR 0001 §D4, normative
in RDR 0007); the evaluator MUST NOT read the tag view.
```

**C10**
```normative
Positive guard atoms MUST live in `all`; negative guard atoms MUST live in
`unless`. Successful row matching MUST NOT depend on source order or
first-match priority.
```

**C11**
```normative
Lint MAY claim guard exhaustiveness only for finite declared domains: enum
values, booleans, declared set element universes, or bounded integer ranges.
```

**C12**
```normative
If a guard dimension lacks a finite declared domain, lint MUST refuse or
downgrade an exhaustiveness claim for that dimension rather than treating the
covered examples as complete.
```

**C13**
```normative
Coverage and overlap checks MUST be scoped to a normalized row group supplied by
the transition/lint model, and MUST evaluate the participating guard dimensions
as one product rather than as independent one-dimensional checks.

A dimension **participates** in a row group when any row in that group carries
a **guard** atom over that key, in either `all` or `unless` — not only when the
rows constrain it differently. The authored block is what makes an atom a guard
atom (JDR 0001 §D6): every atom under `guard.all` or `guard.unless` is a guard
atom regardless of operator, `eq` included. A key every row constrains identically still
bounds the product and still requires a finite declared domain; it MUST NOT be
dropped from the product because it does not discriminate.

**Match keys are not product dimensions.** A row's match pattern selects which
group the row belongs to — the selection context this RDR groups by — and is a
separate field from its guard both in RDR 0002's authored shape
(`[rule.match.*]` vs. `[rule.guard.*]`, routed by block per JDR 0001 §D6) and in the normalized row the kernel
consumes (`internal/resolve/resolve.go::Row` carries `Match` and `Guard`
distinctly, and `Resolve` filters candidates on match before the guard gate).
A key constrained identically by every row in a group *because it is the
grouping key* contributes exactly one assignment and cannot produce a gap;
folding it into the product would multiply every product by 1 and misreport the
cardinality the too-large bound is measured against.

**The match/guard split is per atom per group, not per key per model.** The same
tag key MAY be a match key in one rule and a guard key in another — the fixture
does exactly this with `status`, matched in `profile-to-grounding` and guarded
in `reconcile-rewind-legality`. Participation is therefore decided by where each
atom is authored **within the group being proved**: a key enters that group's
product when some row in *that* group carries a guard atom over it, regardless of
how the key is used in any other group. Reading the split as a per-key property
of the model would make one key's role ambiguous the moment two groups use it
differently, and two implementations would then build different products — the
same cross-implementation divergence the published bound exists to prevent.
```

**C14**
```normative
An `exists` atom projects onto the scoped product as a **per-key presence
dimension**: a two-valued dimension `{present, absent}` for that key, alongside
the key's value dimension. An `exists` atom denotes `{present}` or `{absent}`
on it; a value atom denotes a subset of the key's value dimension **as the
operator/kind agreement clause projects it** — one held value narrowed for
`eq`/`in`/comparisons, which therefore requires the key be declared
single-valued, or a containing subset for `contains` — and, because a value atom
over an absent key is unevaluable rather than false, implicitly `{present}`.
`exists` is unaffected by that requirement: it reads presence alone, so it
projects over a key whose values co-occur exactly as it does over a
single-valued one. A key declared always-present contributes no presence dimension —
its `{absent}` assignment is not in the product — which is what makes the
narrowing's negative control provable rather than vacuous. Lint MUST NOT drop
`exists` atoms from the product: a group carrying one stays provable, and
certifying it exhaustive while ignoring the presence dimension is the
false-green this RDR's narrowing forbids.
```

**C15**
```normative
A declared escape row participates in the coverage identity like any other row:
its accepted assignments are computed from its guard atoms and unioned with its
peers'. An escape row carrying no guard atoms denotes the whole scoped product
and therefore closes coverage by itself. Lint MUST NOT treat "an escape row
exists" as a separate coverage-satisfying fact outside the union.

**An escape row closes coverage only for the failure classes it can actually
rescue** — those it declares in its rescue list, and of those only `no_match`
and `ambiguous_match`. The kernel consults escape rows solely from the
zero-match and multiple-match arms: `internal/resolve/resolve.go::Resolve`
calls `gate` first and returns `refuse(in, *blocked)` immediately when the
gate raises a typed blocking condition, so `escapeOrRefuse` is **never
reached** for `guard_unevaluable` or `owned_state_unavailable`, and within
those arms it selects escapes by `row.rescues(r.Kind)` for the kind that
occurred. An escape row therefore cannot rescue a class it does not declare,
nor either blocking class, however bare its guard. The coverage union is
consequently computed per (scoped row group × declared rescuable class), and a
row declaring one class MUST NOT close the group's other arm; RDR 0006 states
the arm mechanics (`0006::Normative Contracts`, the per-class coverage clause)
and this RDR cites them.

Consequently a bare escape row MUST NOT be read as discharging the narrowing: a
group whose ordinary row can refuse `guard_unevaluable` still has its claim
withheld, because at runtime that refusal is returned before the escape row is
ever consulted. Counting the escape row's assignments toward the union while a
peer can refuse would certify green a group the runtime refuses — the false-green
P5 forbids, reached by citing a rescue path the kernel does not take. This states
that a bare escape row cannot *discharge* another row's withholding; it is not a
narrower population for the narrowing itself, which quantifies over every
participating row including escape rows (see the narrowing clause below) — a
*guarded* escape row can refuse on its own path and withholds the group too.

Overlap is checked in **two separate populations**, because the runtime never
mixes them: the gate-then-count selection (JDR 0001 §D2; RDR 0002
`Load-Bearing Decisions`) and the kernel both resolve ordinary
candidates over non-escape rows only, and consult escape rows solely to rescue a
`no_match` or `ambiguous_match` refusal (`internal/resolve/resolve.go::Resolve`
step 2 — "candidate rows are the non-escape rows for that outcome"). So an
escape row overlapping a *guarded* row is **not** a runtime ambiguity and MUST
NOT be reported as one. What lint MUST still check is overlap **among escape
rows for the same failure class**, since RDR 0002 admits a rescue only when
*exactly one* escape row matches — two overlapping escape rows are a real
ambiguity that silently disables the rescue. Escape rows are therefore excluded
from the ordinary-row overlap check and subjected to their own; excluding them
from **coverage** is what MUST NOT happen.

**The failure class is read from the row's declared rescue list, and the check
runs once per class.** An escape row declares the classes it rescues as a list,
not a single value — RDR 0002's `escape` field (`0002::Normative Contracts`,
the escape-list clause: an `escape` list "MUST contain only resolver failure
classes that RDR 0001 allows the table to model") is plural and the
kernel matches it membership-wise (`internal/resolve/resolve.go::Row.rescues`,
`slices.Contains(r.Escape, kind)`). So lint MUST partition the escape rows into
**one population per declared failure class**, place a row declaring several
classes in **each** of those populations, and run the pairwise overlap check
within each population independently. A single flat pairing over all escape rows
over-reports — two rows that share no class never compete at runtime, because
the kernel filters by `rescues(r.Kind)` before taking the exact-one count — and a
per-row partition under-reports, missing two rows that collide on one shared
class of several. A pair overlapping in more than one class is one finding per
class, naming the class, so the author can see which rescue path is disabled.

A bare escape row consequently closes coverage without generating an overlap
finding against every guarded peer it subsumes — the two checks read the same
row under different populations, which is why they do not contradict.

**A bare escape row is not a silent opt-out from the guarantee.** Because it
closes coverage by itself and draws no overlap finding, it would otherwise be a
one-line, always-available way to make any `graph-coverage-gap` disappear — the
same opt-out this RDR forbids for declarations, arriving through the escape
population. Lint MUST therefore report a bare escape row — one carrying no guard
atoms — that closes a group's coverage as an **observable** result: the group's
verdict names the escape row that closed it, so a reviewer reading the lint
output can tell a group proved over its declared domains from one closed by a
catch-all. Emitting a bare green for such a group MUST NOT satisfy this clause.
The carrier is RDR 0006's to choose — it emits `graph-coverage-closed-by-escape`
naming the row (A19); what this RDR fixes is that the distinction MUST be
visible and MUST NOT depend on the reader inspecting the model by hand.
```

**C16**
```normative
An exhaustiveness claim MUST NOT be stronger than the runtime it describes: lint
MUST NOT certify a row group exhaustive when a participating row can refuse
`guard_unevaluable` under RDR 0007's aggregation veto. Where the two disagree
the lint promise narrows; the runtime veto MUST NOT be weakened.

**"Participating row" here means every row in the group, escape rows included** —
the participation clause's population, not the ordinary-row population the
overlap check uses. A guarded escape row can refuse at runtime on its own path:
`escapeOrRefuse` runs escape candidates through the same viability gate as
ordinary ones and returns `refuse(in, *blocked)` when one carries an undecidable
guard (`internal/resolve/resolve.go::escapeOrRefuse`, "an escape edge … carrying
an undecidable guard raises that typed refusal"), so a group whose only
possibly-refusing row is an escape row is still a group the runtime can refuse.
Withholding is therefore decided over the whole group, while overlap remains
split into the two populations above — the split exists because the runtime
never *matches* the populations against each other, not because escape rows are
exempt from the veto. The two-population reading applies to overlap only; reading
it into the narrowing would certify green exactly the group the runtime refuses.
```

**C17**
```normative
**`unless` is subtracted two-valued only when its atoms are decided.** The
excluded-intersection subtraction this RDR states for lint is a set operation
over decided atoms, while the runtime computes the row verdict in three-valued
Kleene — RDR 0007 fixes the verdict as `all_result ∧ ¬(unless_conj)` with
`¬U = U`, so an **unevaluable atom inside `unless` makes the whole row
unevaluable**, not merely un-excluded.

Lint MUST therefore treat a value atom in an `unless` block over a key declared
optional exactly as it treats one in `all`: the row **can refuse**, and the
group's claim is withheld under the narrowing above. Subtracting an `unless`
block whose atoms may be unevaluable — as if absence made the exclusion simply
not apply — certifies green a group the runtime refuses. This is the same
false-green the narrowing forbids, reached through the block the subtraction
model does not cover.

**A can-refuse row contributes no assignments to the coverage union**, and the
group it sits in has no provable product. The subtraction model is defined over
decided atoms; a row that can refuse has no decidable accepted-assignment set to
contribute, so lint MUST NOT credit it with its `all`-intersection unsubtracted
— that is the false-green above, reached by counting a row the runtime may
refuse. Nor is the shortfall a coverage gap: the report-every-defect clause
scopes the separate gap finding to a gap **over a provable product**, which this
group does not have. So a withheld group emits the withholding finding for each
refusing row, and MUST NOT additionally emit a `graph-coverage-gap` naming a
witness assignment that a refusing row would in fact accept when its key is
present. Overlap findings among the group's decidable rows are unaffected and
still MUST be emitted.
```

**C18**
```normative
A withheld exhaustiveness claim MUST be observable, not silent. It takes the
same blocking inability-to-prove form an unprovable dimension already takes —
RDR 0006's `graph-unprovable-coverage` — and MUST name the participating row
and the atom that can refuse, using the source rule/context id every other
predicate diagnostic names. Emitting nothing MUST NOT satisfy this clause: an
exit code alone cannot distinguish a withheld claim from a proved one.
```

> **Carriers.** The lint surface carries the atom: RDR 0006's finding contract
> requires "A finding attributed to one guard atom MUST carry that atom's `Key`,
> `Operator`, `Literal`, and `Block`" (`0006::Normative Contracts`, the
> finding-record clause; §JD-4's duty, discharged). **The runtime surface is
> short.** The shipped refusal carries `Guard string` — one opaque predicate
> string for the whole refusal (`internal/resolve/resolve.go::Refusal`),
> populated from the lexicographically-lowest undecidable row — while `Rows
> []RowRef` does carry every implicated row. The runtime can therefore name
> *which rows* refused but not *which atom* within them, and where several
> rows refuse it surfaces a single guard string. MVV Scenario 8 asserts the
> finding names "that row and atom": the row half is satisfiable on both
> surfaces, the atom half at runtime is not. Booked as **A20** (JDR 0001
> §JD-18) rather than left standing as a MUST no producer can meet.

**C19**
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

**C20**
```normative
"Refuse" and "downgrade" are one outcome, not an author's choice: every case
this RDR sends to refuse-or-downgrade — a non-finite dimension, a finite
product too large to prove deterministically, and a withheld claim under the
narrowing above — MUST produce a blocking inability-to-prove finding. This
RDR MUST NOT mint a non-blocking warning category for these cases, and RDR
0006 carries none: `graph-unprovable-coverage` for the non-finite dimension and
the withheld claim (JDR 0001 §JD-4 — no new code, no non-blocking tier) and
`graph-product-too-large` for the over-large product, both blocking
(`0006::Technical Design`, the finding-code table). "Downgrade" therefore names
the same blocking outcome as "refuse" on both sides of the seam, and never a
silent or advisory one.
```

**C21**
```normative
A row "can refuse" when the row carries a value atom over a key **declared
optional** **in either block** — `all` or `unless`, since an unevaluable atom
inside `unless` makes the row unevaluable under `¬U = U` rather than dropping
the exclusion; the block, not the operator, is what makes it a guard atom (JDR
0001 §D6) — the same one declared field the optionality clause defines, not
a second presence property and not a graph query. Lint MUST decide this
syntactically over declarations so the test is total; it MUST NOT withhold a
claim merely because some assignment in the product is unreached.

An **owned** tag carries a second, graph-level presence condition — the owned-tag
clause below requires a reachable predecessor write — but that condition governs
whether a row may *match*, not whether its guard can refuse, and it is A12's
subject. Presence for the withholding decision reads exactly one field, so this
clause is decidable today and does not wait on A12: a declaration-only test is
what makes it total.
```

**C22**
```normative
A set literal — the right-hand side of `in` and of `contains` — has one
canonical spelling: an unordered set of typed elements, duplicate-free, and
compared as a set. Two authored spellings differing only in element order or in
repeated elements MUST parse to the same literal and therefore to the same atom
under the identity tuple; a repeated element MUST be rejected at parse rather
than silently collapsed. Implementations MUST canonicalize before the literal
enters the identity tuple, so reordering a set literal cannot change a
diagnostic — the source-order independence MVV Scenario 6 requires.
```

**C23**
```normative
Set-valued guard domains MUST be proved with a deterministic symbolic or
bitset-equivalent representation. If the finite product is too large for that
proof, lint MUST refuse or downgrade the exhaustiveness claim rather than
silently capping enumeration.
```

**C24**
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
**assignment count** — not a bitset width, byte size, or row count. A
dimension's assignment count is not always its declared domain size, and the
difference is exponential, so it is fixed per kind here:

| Kind | Assignment count for one participating dimension |
| --- | --- |
| `enum`, `int` (single-valued) | `\|domain\|` |
| `bool` (single-valued) | 2 |
| `enum`, `bool`, `int` **without** the single-valued marker | `2^\|domain\|` — one independent boolean dimension per value |
| `set` | `2^\|element universe\|` — a set-valued tag holds any subset, so its assignment space is the powerset, **never** `\|universe\|` |
| any optional key | multiplied by 2 for its `{present, absent}` presence dimension |

**The table reads the declaration's defaults, not an author's intent.** Both
markers this table branches on are optional in the declaration and both default
against the smaller product: a declaration carrying no optionality marker is
optional, so it takes the ×2 presence row; a declaration carrying no
single-valued marker is not single-valued, so a finite kind takes the
`2^|domain|` row rather than the `|domain|` one. An implementation MUST apply
these defaults when computing the cardinality — reading an unmarked declaration
as single-valued or always-present is the same exponential understatement the
`set` row forbids, reached through a marker the model never carried. A model
that wants the smaller product declares the markers that license it.

Reading a `set` dimension as `|universe|` understates the product
exponentially — a model with three 20-element set tags computes 8,000 and
"proves" it, while the real assignment space is `2^60`. That is precisely the
memory-exhaustion the "a bound discovered by exhausting memory is not a
conforming bound" sentence forbids, so the arithmetic above is what makes that
sentence enforceable rather than aspirational. The bound is
a single integer constant published by the implementation and reported in the
diagnostic beside the computed cardinality; it is not per-model, per-group, or
configurable per run, since either would make the verdict model-dependent. Two
conforming implementations MAY publish different bounds, but each MUST return
the same verdict for the same model and MUST report which bound it applied.

**A dimension with no finite declared domain has no assignment count**, so a
product containing one has no cardinality: the table above is total over the
finite kinds and defines no value for a `scalar` dimension or for a finite kind
declaring no domain. The too-large comparison is therefore defined **only over a
fully-provable product** — lint MUST compute the cardinality and test the bound
after every participating dimension is known finite, and MUST NOT report a
computed size for a product carrying an unprovable dimension. The two input
classes stay separately reported, as report-every-defect requires: each
unprovable dimension draws its own blocking finding naming that dimension, and
the over-large refusal — the one that carries `(computed size, bound)` — is
simply not among the findings for such a group, since the figure it must report
does not exist. This is not short-circuiting between the classes; it is the
size-bearing diagnostic being unavailable when its quantity is undefined.
```

**C25**
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

**C26**
```normative
Predicate lint MUST distinguish owned, observed, and recognized tags. A row
that matches an owned tag MUST be rejected unless every reachable predecessor
sets or preserves that tag before the match. The reachable-predecessor
relation is RDR 0006's owned-state reachability relation
(`0006::Load-Bearing Decisions`, "Reachability relation" — a total syntactic
dataflow over declarations, writes, and clears, evaluated in its node form;
its finding is `graph-owned-before-write`); this RDR cites it and defines no
second one (A12).
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
| Exhaustiveness verdict for a row group | **This RDR** (finite-domain product proof) | RDR 0006 lint findings; RDR 0005 envelope | RDR 0006's own exhaustiveness clause — non-finite-dimension case only | **This RDR** — JDR 0001 §JD-4 (closed 2026-08-22) names it the recording document; RDR 0006 cites the clause and does not restate it. The two triggers differ and both survive: RDR 0006's fires on a non-finite dimension, this RDR's on a fully-finite product whose participating row can still refuse. |
| Guard atom shape | Kernel seam (JDR 0001 §D1, normative in RDR 0007) | This RDR's grammar; RDR 0002's normalizer | Shipped `Row.Guard string` — the pre-reshape form | Kernel seam (specified); the shipped string form is superseded, not an arm |
| Source identity for diagnostics | RDR 0002 normalization (`RuleID`, `SourceLocator`) | This RDR's diagnostics; RDR 0006 findings | None | RDR 0002 |
| Tag declaration model (value kind, finite domain, optionality, single-valuedness, element universe) | **This RDR** (`Normative Contracts` declaration clauses) | RDR 0002's authoring schema and normalizer; RDR 0006's lint input contract; this RDR's own product proof | **None — settled bilaterally.** RDR 0002's pre-refine prose schema list named `value kind` non-normatively; its refine replaced that with a normative clause naming all five fields and stating "RDR 0003 is the normative home of that model … MUST NOT be restated here" | **This RDR.** RDR 0007 (`Final`) states value kinds and set universes "are RDR 0003's declarations", and RDR 0002 now cites this model for meaning and owns only where it is authored — the rehoming is ratified by the peer, not merely asserted here. |
| Per-tag optionality (which keys may be absent) | **This RDR** (declaration model) | A7's presence dimension; the "can refuse" narrowing | Kernel runtime presence (JDR 0001 §D4) — a different question over one view, not an arm | This RDR |
| Set-element universe for `contains` | **This RDR** (declaration model) | This RDR's finite-domain proof; Phase 3 acceptance gate | None | This RDR |
| Single-valued marker (domain is a partition) | **This RDR** (declaration model) | RDR 0006's `graph-single-valued-state` and its MVV assertion | RDR 0006's lint input contract, which attributed the field to this model before it declared one | **This RDR** — JDR 0001 §JD-13 (2026-08-22); RDR 0006 cites it like the rest of the model |
| Escape-row overlap population | **This RDR** (two-population overlap clause) | RDR 0006's `graph-overlap` findings (invariant 3) | A single population mixing escape and ordinary rows — rejected: the kernel resolves ordinary candidates over non-escape rows only (`resolve.go::Resolve` step 2) and gate-then-count selection (JDR 0001 §D2) admits a rescue only after ordinary resolution fails | **This RDR** — JDR 0001 §JD-14 as corrected 2026-08-23; RDR 0006 cites the clause (A17) |
| Scoped-product membership (which keys are dimensions) | **This RDR** (participation clause) | RDR 0006 lint; the too-large cardinality bound | Match keys — **explicitly not an arm**: they form the selection context the group is defined by, and are a separate field in RDR 0002's authored shape and in `resolve.go::Row` | This RDR — guard keys only |
| Tag provenance (`owned`/`observed`/`recognized`) | RDR 0002 (its one normative tag clause) | This RDR's owned-tag clause (A6); RDR 0004; RDR 0006 | None | RDR 0002 — unchanged by the declaration-model rehoming |
| Withheld-claim lint artifact | **This RDR** (states the blocking form) | RDR 0006 emits it; RDR 0005 envelopes it | `graph-unprovable-coverage` scoped to non-finite dimensions alone — rejected | **This RDR** — JDR 0001 §JD-4 (closed 2026-08-22) names it the recording document; RDR 0006 widens `graph-unprovable-coverage`'s trigger and mints no code |
| `unless` block semantics under an unevaluable atom | Kernel (RDR 0007 — `all_result ∧ ¬(unless_conj)`, `¬U = U`) | This RDR's withholding decision and lint subtraction | This RDR's excluded-intersection subtraction — a **two-valued** model, valid only over decided atoms | Kernel — the subtraction is the lint projection of the kernel's verdict, never a second semantics |
| Which failure classes an escape row can rescue | Kernel (`resolve.go::Resolve` — gate refusal returns before `escapeOrRefuse`; `Row.rescues` filters by declared class) | This RDR's coverage-union clause; RDR 0006's per-class coverage union | A bare escape row read as closing coverage unconditionally — rejected | Kernel — escape rows rescue `no_match` and `ambiguous_match` only, and only the classes they declare; never `guard_unevaluable` or `owned_state_unavailable` |
| Per-dimension assignment count for the cardinality bound | **This RDR** (the per-kind table in the bound clause) | The too-large refusal (`graph-product-too-large`); RDR 0006's finding | "Declared domain size" — rejected: it understates `set` and unmarked kinds exponentially | **This RDR** — a `set` dimension is `2^\|universe\|`, an unmarked finite kind `2^\|domain\|`, an optional key ×2 |
| Conforming-view enforcement | **Owned half: RDR 0006** (`graph-always-present-owned`; invariant 5). **View-level half: unhomed — A18** | Every lint claim in this RDR is conditional on it | RDR 0007's view assembly (states no obligation); shipped `resolve.go::assemble` (performs no check) | RDR 0006 for owned keys over the declared model; **no current owner** for observed/recognized keys at runtime — JDR 0001 §JD-18 decides, A18 books it |

#### `disposition` — input class × outcome

Cue: the draft sets outcomes over input classes (refuse, downgrade, prune,
unevaluable).

| Input class | Runtime outcome | Lint outcome | Diagnostic minted | Silent or loud |
| --- | --- | --- | --- | --- |
| Every `all` atom decides true; `unless` block not fully true | Row qualifies | Row contributes its accepted assignments | none | — |
| `all` atom decides false | Row pruned | Row contributes nothing | none | Silent by design — a decided false is not a defect |
| Full `unless` block decides true | Row disabled | Excluded intersection subtracted | none | Silent by design |
| Value atom **inside `unless`** over an absent key | `guard_unevaluable` refusal — `¬U = U`, so the row is unevaluable, not merely un-excluded (`0007::Normative Contracts`, the strong-Kleene atom-verdict clause and the row-verdict clause) | Exhaustiveness claim **withheld** for that group; the block is **not** subtracted as if decided | Runtime: RDR 0007's per-row/per-atom payload. Lint: the blocking inability-to-prove finding, naming the row and the refusing atom | Loud — both surfaces mint an artifact |
| Value atom over an absent key | `guard_unevaluable` refusal (RDR 0007 veto) | Exhaustiveness claim **withheld** for that group | Runtime: RDR 0007's per-row/per-atom payload (key, block, reason `absent`). Lint: the blocking inability-to-prove finding, naming the row and the refusing atom | Loud — both surfaces mint an artifact |
| Existence atom over an absent key | Decided (`presence == literal`) — never unevaluable | Selects `{absent}` on the presence dimension (A7) | none | — |
| Zero rows qualify | RDR 0001 refuses | Coverage gap if the product is provable | `graph-coverage-gap` (RDR 0006), naming the selection context, every rule id in the group, and one uncovered assignment | Loud |
| Two or more rows qualify | RDR 0001 refuses — never first-match | Overlap finding | `graph-overlap` (RDR 0006), naming both source rule ids | Loud |
| Guard dimension lacks a finite declared domain | Evaluable at runtime | **Refuse or downgrade** the group's claim, one finding per unprovable dimension — never treat examples as complete, never stop at the first | `graph-unprovable-coverage` (RDR 0006) naming the dimension | Loud |
| Finite product too large to prove deterministically | Evaluable at runtime | **Refuse or downgrade** — never silently cap enumeration | `graph-product-too-large` (RDR 0006) reporting the computed product cardinality and the published bound | Loud |
| Unknown operator / unknown tag / operator–kind mismatch / literal parse failure | Rejected before resolution | Rejected at load | Predicate semantic kind (this RDR) → RDR 0002 load category (`malformed predicate atom`, `unknown tag`) → RDR 0005 envelope (§JD-8) | Loud |
| Guard literal outside its tag's declared domain | Rejected before resolution | Rejected at guard parse, after declarations load and before rows are yielded | `literal-outside-declared-domain` (this RDR) → RDR 0002 `malformed predicate atom` → RDR 0005 envelope (§JD-8) | Loud — an unsatisfiable atom is never a silently-never-matching row |
| Declaration whose domain disagrees with its kind (incl. `single_valued` on a `set` or on a kind with no finite domain) | n/a — never reaches runtime | Rejected by the declaration loader before normalization completes | Declaration/kind disagreement (this RDR) → RDR 0002 `malformed tag declaration` → RDR 0005 envelope (§JD-8) | Loud |
| Two escape rows matching one failure class | RDR 0002 admits a rescue only when exactly one matches — the rescue silently does not fire | Overlap finding within the escape population | `graph-overlap` (RDR 0006), naming both escape rule ids and the failure class | Loud |
| Escape row overlapping a guarded row | Not an ambiguity — the kernel resolves ordinary candidates over non-escape rows only | **No overlap finding**; the escape row still contributes its assignments to the coverage union | none | Silent by design — reporting it would fail a model the runtime accepts (A17) |
| Row group is domain-exhaustive but a participating row can refuse | Refusal stands | Claim withheld — the narrowing; one finding per refusing row, never just the first | RDR 0007 payload at runtime **and** a blocking lint finding naming the refusing atom — absence of a green result is not the artifact | Loud |

#### `trace` — desk trace over the MVV

Cue: ten normative clauses bear on one output surface — the exhaustiveness
verdict for a row group. Walked stepwise against the MVV, with witnesses from
`evidence/spikes/guard-fixture.toml`.

| Step | Assertions in force | Witness | Verdict |
| --- | --- | --- | --- |
| 1. Load the RDR/kata slice as normalized candidate rows | atoms-not-callbacks; closed typed operator vocabulary; operator declares accepted kinds | `profile-to-grounding` parses to `profile in [mid,large]` + `unless prelock_iterations gte 3`; operators all in the matrix | OK |
| 2. Group rows by selection context | coverage/overlap scoped to a normalized row group, evaluated as one product | `profile-to-grounding` and `foundational-to-cove` share `match status eq Draft` → one group | OK |
| 3. Build the scoped product from declared domains | exhaustiveness only over finite declared domains; every *participating* **guard** dimension enters the product, including one only some rows constrain; match keys are the grouping context, not dimensions | Computed with the assignment-count table's declaration defaults, which this fixture triggers throughout: it declares **no** `single_valued` and **no** optionality marker on any tag, so every dimension takes the unmarked `2^\|domain\|` row and every key takes the ×2 presence row. `profile` = `2^4 · 2` = 32; `prelock_iterations` = `{0..3}` inclusive so `2^4 · 2` = 32; `cluster_eligible` = `2^2 · 2` = 8 — **8,192 assignments in all**. `status` does **not** enter the product: its `status eq "Draft"` atom is authored under `[rule.match.status]`, so it is the grouping key that forms this selection context, not a guard dimension within it. `prelock_iterations` is the participation witness — only `profile-to-grounding` carries an atom over it (`unless … gte 3`), and it enters the product for the whole group even though `foundational-to-cove` never mentions it | OK — and it is the participation clause that keeps `prelock_iterations` in the product; under a rows-must-differ reading a key only one row constrains would silently drop out |
| 4. Project each atom onto the product | every atom denotes a subset of the scoped product (A2); `eq`/`in`/comparisons are single-value operators and project only over a single-valued dimension | **Blocked, and this is the step that catches it.** Every guard atom in the group is a single-value operator — `profile in [mid,large]`, `profile eq "foundational"`, `unless prelock_iterations gte 3` — over a tag the fixture does **not** declare single-valued. Under the assignment-count defaults each of those dimensions is a power set (step 3), so the tag has no single held value for the operator to narrow and the atom has no projection. `cluster_eligible exists = true` is the one atom that *does* project: `exists` reads the presence dimension alone, which the defaults give it | **UNPROVABLE — one blocking finding per dimension** (`profile`, `prelock_iterations`), naming the dimension and the unprojectable atom. Not a gap: coverage is not computed for a group whose atoms do not project |
| 5. Project the `exists` atom | A7 presence dimension — stated as a normative clause | `foundational-to-cove` carries `cluster_eligible exists = true`; `cluster_eligible`'s declared domain `{true,false}` is complete, so no *value* element selects presence — the atom lands on the separate `{present, absent}` presence dimension the projection clause adds | OK — witness: the presence-dimension clause in `Normative Contracts`, on the derivation in `evidence/research/iter-2-projection-derivation.md` |
| 6. Compute coverage | `union(row_i accepted) == scoped product`; a gap finding is scoped to a gap over a **provable** product | **Not reached for this fixture, and that is the correct outcome.** Step 4 left two dimensions unprojectable, so this group has no provable product and no union to compare — the report-every-defect clause scopes the separate gap finding to a provable product precisely so lint does not manufacture a witness assignment out of atoms it could not project. The blocking findings from step 4 stand alone | **No coverage verdict** — withheld, not green and not GAP. At a single-valued, always-present reading the fixture never declares, this row would compute product 64, union 32, "GAP"; at the declared defaults the product is 8,192 and the atoms do not project |
| 7. Compute overlap | any non-empty pairwise intersection, computed **within** a population — ordinary rows against ordinary rows, escape rows against escape rows for one failure class; no source-order priority | Also not reached: an intersection is taken between projected assignment sets, and step 4 produced none for the two ordinary rows. The fixture declares no escape row, so the escape population is empty and contributes no finding either | **No overlap verdict.** `profile in [mid,large]` ∩ `profile eq foundational` = ∅ holds only under the single-valued reading; with `profile` unmarked the two atoms have no projections to intersect, and whether they would overlap is exactly the question the projection clause refuses to answer by guess |
| 8. Apply the runtime-veto narrowing | claim MUST NOT be stronger than the runtime; withhold if a participating row can refuse | MVV Scenario 8's row group: domain-exhaustive, one value atom over a possibly-absent key → claim withheld | OK — and the reason A7's presence dimension must not silently drop `exists` atoms |
| 9. Emit the verdict | diagnostics name the contributing source rule/context id; every decidable defect is reported, not the first; a gap additionally names the context, all group rule ids, and one uncovered assignment | `RuleID` + `SourceLocator` ship on `internal/resolve/resolve.go::Row`; the uncovered assignment is computed from the scoped product built at step 3 | OK |
| 10. Cross-document agreement | exactly one document records the narrowing; the other cites it | §JD-4 (closed 2026-08-22) names this RDR the recording document; RDR 0006's narrowing-citation clause (`0006::Normative Contracts`) cites it, reuses `graph-unprovable-coverage`, and "MUST NOT restate it" | OK |

No CONTRADICTION row. **Step 4 is the blocking row**: computing step 3 at the
assignment-count table's declared defaults — rather than assuming the markers the
fixture never wrote — puts every guard dimension at `2^|domain|`, and the group's
`eq`/`in`/`gte` atoms then have no single held value to narrow. The group is
**unprovable**, so steps 6 and 7 are not reached. The defect is in the *fixture*,
which declares kinds and domains but neither marker, and was authored to exercise
operators rather than to carry a proof; the clauses behave correctly on it.
Computing step 3 at the fixture's declared defaults, rather than at the
single-valued, always-present semantics it never declares, is the
uncomputed-arithmetic defect this trace exists to catch. The fixture is
therefore not usable as the MVV's "one complete partition" case without
**declaring `single_valued` on `profile`, `prelock_iterations` and
`cluster_eligible`** and then adding the rows that route `profile = small` and
`prelock_iterations = 3` — both recorded as Phase 3 prerequisites in `Testing
Strategy`. Marker-declaring is the load-bearing half: without it the group
never reaches a coverage verdict at all, so the added rows would close nothing.
A9 and A11 have no row here because the declaration model they record is an
input to step 3, which the trace already exercises.

Two clauses are **unexercised by this fixture** rather than gapped, and each has
a named test instead: the escape-row overlap population (the fixture declares no
escape row — MVV Scenario 2's group construction), and the single-valued
marker's effect on product shape (the fixture declares no marker — MVV
Scenario 5). Step 3's product is built from guard keys only, so `status` is
the grouping context here rather than a dimension.

#### Load-Bearing Decisions

- **Identity** — a guard atom is identified by the total tuple
  `(RuleID, SourceLocator, key, block, operator, literal)`, never by its index
  within `all` or `unless`. Position is not a field of the kernel-fixed atom
  shape (JDR 0001 §D1: key, operator token, literal, block), and an index-based
  handle would contradict this RDR's own source-order independence: the same
  tuple must identify the same atom under the reordering MVV Scenario 6
  requires. The tuple is total precisely because one row may carry two atoms
  over one key (RDR 0007:A5's conjoined value row), which `(key, block)` alone
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
- **Declaration model** — a tag declaration carries five fields: value kind,
  finite domain, optionality, single-valuedness, and (for set kinds) an element
  universe. Single-valuedness is the fifth, homed here by JDR 0001 §JD-13; any
  enumeration of this model in this document states all five.
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

Illustrative only — RDR 0002 owns the authored container in both cases; this
RDR owns what the declaration fields mean and what the guard atoms denote.

A tag declaration carrying the five fields this RDR defines:

```toml
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "mid", "large", "foundational"]
single_valued = true     # at most one value holds per view — the partition licence
required = true          # always-present; contributes no {absent} assignment

[tags.labels]
provenance = "observed"
kind = "set"
elements = ["urgent", "blocked", "external"]   # element universe; no single_valued
```

The guard atoms written against that alphabet:

```toml
[rule.guard.all]
profile.in = ["mid", "large", "foundational"]

[rule.guard.unless]
prelock_iterations.gte = 3
```

The key spellings above are RDR 0002's wire keys (JDR 0001 §D7(iii); A16) —
this block shows the shape, RDR 0002 asserts it. `required` defaults to
optional when omitted, as the optionality clause mandates.

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Sparse transition-table container | RDR 0002 | Pending | This RDR assumes guards live inside RDR 0002's `all`/`unless` blocks. |
| Deterministic exact-one resolver | RDR 0001 | Pending | Runtime selection refuses zero or multiple matching rows. |
| Guard predicate grammar and finite-domain semantics | This RDR | Introduced | Lint and runtime share one symbolic predicate model. |
| **Tag declaration model** (value kind, finite domain, optionality, single-valuedness, set-element universe) | **This RDR** | **Introduced** | The typed alphabet guard atoms are written against. Homed here because it is the guard algebra's input alphabet: RDR 0002's normative tag vocabulary is provenance plus wire-key spelling and cites this model for meaning, and RDR 0007 (`Final`) states that value kinds and set universes "are RDR 0003's declarations". RDR 0002 owns authoring location and normalization carriage (A11). |
| Accessor read/write safety | RDR 0004 | Pending | Guard evaluation consumes tag values after accessor binding; it does not execute accessors. |
| Graph lint authority | RDR 0006 | Pending | Exhaustiveness and overlap findings become blocking lint there. |
| Kernel guard-domain enforcement and `guard_unevaluable` payload | RDR 0007 | Pending | The kernel decides presence and existence atoms; this RDR's evaluator narrows to value semantics over a present value. |
| Atom `block` retention through normalization | RDR 0002 | **Landed — A14 Verified** | Per-atom `block` retention is normalization *carriage*, RDR 0002's charter proper, not declaration *semantics*. RDR 0002 states it normatively — each atom MUST retain key, operator token, literal, and authored block, and normalization MUST NOT fold `unless` atoms into `all`, citing this RDR's identity tuple and `unless` semantics as the reason — and carries `block` into the dump field list and the within-row sort key. |
| Authoring location for the **single-valued marker** | RDR 0002 | **Landed — A16 Verified** | §JD-13 homed the field's meaning here; RDR 0002 carries `single_valued` in both its `[tags.<tag>]` schema enumeration and its normative type-model clause, each citing this RDR for meaning, under a carriage clause that is general ("every declared field through without loss"). |
| Authoring location + normalization carriage for tag declarations | RDR 0002 | Landed (specified; `Final`) | RDR 0002 owns `[tags.<tag>]` placement, the wire keys `kind`/`domain`/`min`/`max`/`elements`/`single_valued`/`required` (JDR 0001 §D7(iii)), the two load categories carrying this RDR's rejection rules, and carrying declarations through to candidate rows; it cites this RDR's declaration model for what the fields mean. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| User-facing failures | `internal/cli/clierr`, `internal/cli/respond` | Must verify error-code coverage | Reuse | Predicate parse/lint failures should use existing CLI envelopes. |
| Guard evaluator | None found under `internal/` | New semantic surface | Introduce | This RDR owns the evaluator semantics, not command I/O. |
| Transition model source | Pending RDR 0002 | Not implemented | Extend peer | Guard atoms are embedded in table rows/contexts. |
| Tag declaration model | RDR 0002 declares provenance and the wire-key spellings normatively, and cites this RDR for meaning | The typed alphabet the guard grammar quantifies over has no other normative home | Introduce here | This RDR states the model (all five fields); RDR 0002 owns authoring location plus normalization carriage and carries every field, the single-valued marker included (A16 Verified). |

### Decision Rationale

Joint-check: fired → RDR 0007 is the normative home of the guard seam (JDR 0001
§D1, §D4). It fixes the parsed-atom row shape, the value-only per-atom
`GuardEvaluator` seam, the kernel-exported existence operator token and boolean
literal forms RDR 0002's normalizer emits, exact key identity at the kernel, and
the per-atom `guard_unevaluable` payload. This RDR cites those rather than
restating them. It also records the §JD-4 narrowing for its own proof: §JD-4 decides the
substance and, closed 2026-08-22, names this RDR the recording document, with
RDR 0006 citing the clause and minting no code (A8).

The fixed symbolic atom model best matches the user's outcome: flow authors can
write conditional edges, and lint can still prove whether those edges are
complete and mutually exclusive. Two classes are provable only because this RDR
states the clause: the cross-document narrowing verdict (§JD-4 names this RDR
the recording document, A8) and the `exists`-bearing class (the
presence-dimension projection, A7). Consumers inheriting
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

Premortem: survived. If it shipped and failed, the most
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

Visible failures should be typed load or lint failures: unknown operator,
operator/tag-kind mismatch, literal parse failure, unknown tag, and the two
declaration/literal-domain rejections (RDR 0002 load categories, §D7(iii)) at
load; non-exhaustive finite domain, overlapping candidate rows, guard
dimension not provable because it lacks a finite domain, or finite scoped
product too large to prove at lint. A guard
that cannot be decided at runtime surfaces as RDR 0007's `guard_unevaluable`
refusal, not as a predicate parse kind owned here. Silent
failure would be a false exhaustiveness claim; the recovery path is to keep
every exactness claim tied to A2 and the MVV fixture.

## Implementation Plan

### Prerequisites

- [x] **A1** — the evaluation harness ran: representative rows are *evaluated*
  against tag inputs with asserted qualify/prune verdicts (28 pass / 0 fail),
  all eight operators carry an evaluated atom, and five typed-rejection
  negative controls confirm the vocabulary is closed and typed
  (`evidence/spikes/iter-2/a1-eval-harness/{main.go,output.txt}`). Vocabulary
  sufficiency is independent of the declaration model.
- [x] **A2** — the derivation's finite-domain inputs are declared by this
  RDR's own tag declaration model (A11), so the derivation stamps without
  circularity.
- [x] A3, A4, A5 verified (A5 against the kernel-fixed atom shape, JDR 0001
  §D1). A6 verified for provenance *labels*; its reachability half is A12.
- [x] **A11** (parent of A7 and A9) — the tag declaration model is stated here
  as normative contracts: value kinds, finite domains per kind, optionality,
  single-valuedness, element universe, and domain/kind agreement. RDR 0002
  spells the wire keys and cites this model; RDR 0007 (`Final`) assigns value
  kinds and set universes to this RDR.
- [x] **A7** — the presence-dimension projection is stated as a normative
  clause, on the derivation in
  `evidence/research/iter-2-projection-derivation.md`; its optionality input
  is this RDR's own declaration. RDR 0007:A12 routes the question here and it
  is answered here.
- [x] **A8** — JDR 0001 §JD-4 (closed 2026-08-22) names this RDR the recording
  document for the narrowing; RDR 0006 cites the clause, reuses
  `graph-unprovable-coverage`, mints no code and no non-blocking tier, and
  carries the atom-level finding field the clause requires. Evidence:
  `docs/rdr/cluster-reconcile/0003-0006-0007/`.
- [x] **A10** — RDR 0006 reads the row-group division of labour this RDR
  states and defines no second grouping predicate (`0006::Normative
  Contracts`, the row-group clause).
- [x] **A12** — RDR 0006 publishes the owned-state reachability relation and
  states which procedure governs each clause (reachability for the owned-tag
  check; a syntactic decision over the optionality field for "can refuse");
  this RDR's owned-tag clause cites it.
- [x] **A13** — the canonical set-literal spelling is stated as a normative
  clause: unordered, duplicate-free, canonicalized before entering the
  identity tuple, repeats rejected at parse. Affects `in`, which is inside the
  claimed-proven operator subset.
- [x] **A14** — RDR 0002's normalization clause retains each atom's authored
  block, with an explicit MUST NOT against folding `unless` into `all`, and
  carries `block` into the dump field list and the within-row sort key. MVV
  Scenario 7 keeps it honest.
- [x] **§JD-13 discharged** — the declaration model carries the single-valued
  marker, a normative clause states what it licenses, and the domain/kind
  agreement clause rejects it on a `set` kind or a kind with no finite domain.
  Authoring half: **A16** closed (RDR 0002 carries `single_valued` on both its
  surfaces). Consumer half: RDR 0006 reads the marker from the declaration and
  reconciles its invariant 5 as the model-level half of this RDR's
  single-valued conformance conjunct.
- [x] **A16** — `single_valued` occurs in RDR 0002's authoring schema and in
  its normative type-model clause, each citing this RDR for meaning. This is
  what keeps the MVV's green cases (Scenario 2's partition, Scenario 8's
  negative control) authorable under the projection clause.
- [x] **A17 / §JD-14** — escape-row overlap is checked among escape rows for
  one failure class, never between an escape row and a guarded row; JDR 0001
  §JD-14 (corrected 2026-08-23) and RDR 0006's escape-row clause both state
  it. The coverage half of §JD-14 stands: escape rows are ordinary
  participants in the union, computed per declared rescuable class.
- [x] **A19** — RDR 0006 emits `graph-coverage-closed-by-escape` naming the
  bare escape row that closed a group's coverage.
- [ ] **A18** — Route: JDR 0001 §JD-18 decides which document states the
  view-level conformance check for observed/recognized always-present keys
  and where it runs (view assembly rejects a non-conforming view, or exposes
  it as a typed refusal). Done-condition: the answer is cited here beside RDR
  0006's owned-half producer (`graph-always-present-owned`, invariant 5).
  Venue: JDR 0001 §JD-18 (siblings 0003, 0007). Not lock-blocking.
- [ ] **A20** — Route: the runtime half of the atom carrier — RDR 0007's
  refusal payload, which today carries one opaque `Guard string` beside `Rows
  []RowRef` (`internal/resolve/resolve.go::Refusal`). Done-condition: the
  refusing atom is addressable on the runtime surface as it already is on
  RDR 0006's finding contract. Venue: JDR 0001 §JD-18, with A18. Not
  lock-blocking.
- [ ] **A15** — Route: MVV Scenario 3 compares two equal-cardinality products
  of differing shape for the same verdict and measures whether the proof
  representation completes for each, confirming a scalar cardinality bound
  predicts provability. Owned here; discharged by the MVV.
- [ ] **A21** — Route: MVV Scenario 4 asserts the negative control (an
  unmarked tag under `eq` yields the blocking inability-to-prove finding naming
  the dimension, not a coverage gap) and Phase 3's fixture declares
  `single_valued` on every guard dimension the two target flows use.
  Done-condition: the target flows are authorable under the projection clause,
  or the dimension that is not is recorded as a `set` tag under `contains`.
  Owned here; discharged by the MVV. A failure opens the marker-inversion fork
  the record names.
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
0006's `graph-single-valued-state` reads, defined here and rejected on a `set`
kind. The marker is load-bearing for the whole positive half of the MVV:
`eq`/`in`/comparisons are single-value operators whose atoms over an unmarked
dimension do not project at all, so **every case whose expected verdict is
green** — Scenario 2's complete partition (the MVV's only passing
exhaustiveness verdict) and Scenario 8's always-present negative control —
depends on `single_valued` being writable, which RDR 0002's schema makes it
(A16). The fields are normative declarations this document owns, spelled by
RDR 0002's wire keys, not fixture-invented keys read back as evidence. The
MVV's remaining dependency is sequencing, not authorability: RDR 0007's kernel
reshape must land before Phase 1 builds against the atom-slice shape.

**Phase gating.** The phases below are not a linear sequence — each names the
assumptions that block it, so no phase is picked up on document order alone:

| Phase | Blocked on | Startable today |
| --- | --- | --- |
| 1 Predicate Model | RDR 0007's kernel reshape (shipped kernel still has `Row.Guard string`) — A13's set-literal spelling and the tag declaration model are stated | Partially — the matrix, the declaration model, and the set-literal canonicalization, not the atom slice |
| 2 Finite-Domain Lint Semantics | Nothing open on peers — A10, A12, A17 and A19 are closed against RDR 0006 (`Final`), and the domain producer is this RDR | Yes — the product arithmetic (per-kind assignment-count table), the withholding decision (a declaration-only test, extended to `unless` atoms), the single-valued partition, group construction, the two-population escape-row overlap, and bare-escape observability |
| 3 Target-Flow Fixture | Nothing open — A14 and A16 closed against RDR 0002, A1's harness has run, and A7/A9 are closed | Yes. Note the fixture's Draft group does not reach a coverage verdict (desk trace step 4: no tag in the *fixture* declares `single_valued`, so its `eq`/`in`/`gte` atoms do not project and the group is unprovable) — a fixture gap, not a schema one, so Scenario 2's partition case needs the markers declared **and** rows added, or a separate group authored |
| 4 Integration With Peer RDRs | RDR 0007's kernel reshape; §JD-18 for the runtime atom carrier and view-level conformance (A18, A20 — neither gates the lint-side integration) | No |

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

**The fixture is not the MVV's complete-partition case, and misses by two
independent margins.** The desk trace's step 4 blocks first: no tag declares
`single_valued`, so at the assignment-count table's defaults every guard
dimension is a power set, the group's `eq`/`in`/`gte` atoms have no single held
value to narrow, and the group is unprovable before coverage is ever computed.
Behind that, the rows do not partition either — nothing routes `profile = small`,
and `prelock_iterations = 3` falls through `profile-to-grounding`'s `unless`
block. The fixture was authored to exercise the operator vocabulary, and it does
that. Phase 3 must therefore declare the markers **and** close the routing, or
author a separate group for Scenario 2's complete partition — reusing it
unchanged makes Scenario 2 fail on first run against a correct lint, and fail at
the projection step rather than the coverage step it is written to test.

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
   unprovable dimensions, and two equal-cardinality products of differing shape.
   **The shape pair is constructed relative to the implementation's published
   bound B, not to an absolute size**: build two products of equal cardinality
   C with C > B — one from few wide dimensions (e.g. two dimensions of
   ~sqrt(C) values each), one from many narrow ones (e.g. log2(C) boolean
   dimensions) — so both must refuse if cardinality alone gates the verdict.
   Repeat the pair just under B, where both must prove. Reading B at fixture
   setup is what makes the case runnable at implementation time.
   **Expected**: Runtime evaluation remains available, but lint refuses or
   downgrades the exhaustiveness claim for that dimension/product. The
   two-dimension group emits two findings, not one. The over-large refusal
   reports both the computed product cardinality and the published bound. The
   two equal-cardinality products receive the same verdict — a divergence
   refutes A15 and the bound clause is restated (a spec defect routed back
   to this RDR).
   **This scenario tests the bound's consistency, not its adequacy.** Reading `B`
   from the implementation makes the verdicts comparable, but it also means the
   test derives its oracle from the system under test and so cannot falsify the
   claim A15 actually makes — that cardinality *predicts provability*. Add the
   measurement that can: for each shape, record whether the proof representation
   actually completes within the implementation's own resource budget, and
   compare that to the verdict the bound gave. A15 is refuted when the bound says
   provable and the representation does not complete, or vice versa — a
   divergence between predicted and observed provability, which is a fact about
   the implementation and not about `B`.
4. **Scenario**: Parse malformed guard atoms: unknown tag, unknown operator,
   unsupported operator/tag-kind pair, literal parse mismatch, and a **literal
   outside its tag's declared domain** (a well-typed value the domain does not
   contain — distinct from a parse failure, and rejected before resolution).
   Include the malformed *declarations* the domain/kind agreement clause
   rejects: a `{min..max}` bound on an `enum`, an element universe on a scalar
   kind, and a single-valued marker on a `set` kind or on a kind carrying no
   finite domain.
   **Expected**: Each failure is rejected before resolution with a predicate
   semantic kind that RDR 0006 can map to a lint finding and RDR 0005 can map to
   the structured CLI gateway. The declaration errors are rejected before
   normalization completes, so a consumer reading the model — including RDR
   0006's `graph-single-valued-state` — never sees a marker its kind cannot
   carry.
5. **Scenario (single-valued acceptance)**: Lint one row group over a
   finite-domain tag **with** and **without** the single-valued marker, holding every other input fixed.
   **Expected**: the scoped product differs by construction — with the marker the
   tag contributes one dimension of `|domain|` assignments; without it, one
   boolean dimension per value (`2^|domain|`) — so a row set that closes coverage
   under the marker leaves an uncovered assignment without it. The two runs must
   therefore return different verdicts on the same rows; identical verdicts mean
   the marker is not reaching the product and `graph-single-valued-state` has no
   effective producer. Assert the marker is read from the declaration and never
   inferred from the tag's name or value spelling.
6. **Scenario**: Reorder authored rows and guard atoms without changing their
   semantics, including reordering the elements inside an `in` set literal.
   **Expected**: Successful matching and lint findings are unchanged because
   source order is not a selection mechanism; an ambiguous pair remains a
   multiple-match refusal instead of becoming a first-match success.
7. **Scenario (block retention)**: Normalize a row carrying atoms over the
   **same key in both blocks** — one in `all`, one in `unless` — and inspect the normalized atoms.
   **Expected**: each atom still reports the block it was authored in
   (`Block == all` / `Block == unless`), and the two remain distinguishable and
   separately identifiable under the identity tuple. This is the "first
   normalization test" A14's survivable-if-wrong disposition names: if carriage
   drops `block`, this scenario fails immediately rather than surfacing as a
   wrong `unless` verdict later.
8. **Scenario**: Evaluate a row group whose guards are exhaustive over their
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
   **Both blocks, and the escape path.** Run the same group a second time with
   the optional-key value atom moved from `all` into `unless`: the verdict must be
   identical — claim withheld — because `¬U = U` makes the row unevaluable rather
   than un-excluded. A run that certifies the `unless` variant green while
   withholding the `all` variant is the two-valued-subtraction defect. Then add a
   bare escape row to the group: the claim must **still** be withheld, since
   `resolve.go::Resolve` returns the gate's refusal before `escapeOrRefuse` is
   reachable, so the escape row cannot rescue `guard_unevaluable`. A run that
   goes green once the escape row is present is the false-green A19 guards.

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

Responses: `0003-guard-predicate-exhaustiveness/artifacts/gate.md` (Gate PASS 2026-08-24)

## References

- RDR 0001, Resolution Kernel Contract.
- RDR 0002, Transition Table As Reviewable Data.
- RDR 0006, Graph Lint Authority and Guarantees.
- RDR 0007, Guard Predicate Totality — normative home of the guard seam, the
  domain rule, and the `guard_unevaluable` payload.
- JDR 0001, Resolve Kernel Seam — §D1 (parsed-atom row shape), §D4 (kernel
  enforces the guard domain), §D6 (the authored block routes an atom to match
  or guard), §D7(iii) (type-model wire keys; the two load categories carrying
  this RDR's rejection rules), §JD-4 (the lint promise narrows — closed
  2026-08-22 naming this RDR the recording document), §JD-8 (envelope mapping
  of refusal and load categories), §JD-13 (single-valued marker homed in this
  RDR's declaration model), §JD-14 (escape rows in the union; overlap in two
  populations, corrected 2026-08-23), §JD-18 (conforming-view enforcer and
  runtime atom carrier — open; A18, A20).
- `docs/cli-output-contract.md`.
- Resource index: `.rdr/resources.md`.
- Seed prior: `../state-machines/BUILD-SEEDS.md`, especially the guard
  expression and transition-table seeds.
- Transition model prior: "The transition model — inputs, outputs, error
  conditions."
- Prior-art corpus: `../state-machines` audits, evals, contrasts, and checked
  repositories for `transitions`, stateless/qmuntal-stateless, SCXML/scxmlcc,
  Sismic, StateSmith, and Statewright.
