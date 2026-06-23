# Recommendation 0003: Guard Predicate Exhaustiveness

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-06-19
- **Status**: Draft
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

The seed is still current: there is no implemented resolver, transition table,
or guard package under `internal/`; the only current consumers are the peer RDRs
and the future lint/resolver seams. RDR 0001 requires guard evaluation to feed
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

The search found no implemented guard predicate evaluator under `internal/`.
The adjacent design signal is RDR 0002's `all`/`unless` split and exact-one row
selection, so this RDR extends that signal instead of inventing a parallel
guard model.

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
  - **Evidence**: `cd docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes && sh check.sh guard-fixture.toml` validated four representative guard rows covering status/profile routing, cap-3 handling, prelock lens sets, cluster eligibility, and rewind legality with only `eq`, `in`, `lt`, `gte`, and `exists`; transcript captured in `docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes/output.txt`.
  - **If wrong**: The fixed operator set is too small, and authors will need an
    expression grammar or host predicates that weaken static lint.
- **A2 Every exhaustiveness claim can be reduced to finite declared domains.**
  - **Status**: Verified
  - **Method**: Derivation
  - **Evidence**: For every exhaustiveness-eligible tag, lint receives a finite domain `D`: enum/boolean values as declared sets, set-valued tags as the declared element universe's finite powerset or an implementation-equivalent bitset, and bounded integers as `{min..max}`. A row with `all` atoms denotes the intersection of each atom's allowed subset of `D`; its `unless` block denotes an excluded intersection that is subtracted from the row's accepted assignments. Coverage is `union(row_i accepted assignments) == D` for the checked dimension/product, and overlap is any non-empty `row_i accepted assignments intersect row_j accepted assignments`. If any participating dimension lacks finite `D`, the proof cannot enumerate coverage and lint must refuse or downgrade the exhaustiveness claim.
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
  - **Evidence**: `internal/cli/clierr/clierr.go::CLIError` carries stable `Code`, human `Message`, optional `Param`, `Detail`, `Hint`, and exit-code `Group`; `internal/cli/respond/respond.go::Fail` emits the envelope in text/json modes; `internal/cli/config/config.go::Load` already demonstrates stable parse/read error codes. Predicate parse, unknown-operator, type-mismatch, non-exhaustive, overlap, and unevaluable-guard failures can add stable codes on this existing gateway.
  - **If wrong**: This RDR or RDR 0005 must add a separate user-facing error
    contract before implementation.
- **A5 The normalized predicate representation can retain source identity for
  actionable diagnostics.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0002 `Normative Contracts` require each normalized candidate row to retain source rule id and source locator; its `Validation / Testing Strategy` requires normalized rows to retain source rule ids/source locators while inherited predicates and `all`/`unless` guards are expanded. RDR 0006 also consumes source rule ids/spans for graph lint findings.
  - **If wrong**: Lint may detect an error but fail to point reviewers at the
    guard to fix.
- **A6 Tag provenance is available to predicate lint.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0002 `Normative Contracts` require the model to declare every matched or written tag, including `owned`, `observed`, or `recognized` provenance. RDR 0006 `Technical Design` consumes tag provenance and owned-tag write effects from normalized rows for owned-set-before-match and coverage checks.
  - **If wrong**: Predicate lint can still evaluate runtime truth, but it cannot
    prove that owned tags are set before they are matched.

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

Use a closed, symbolic predicate-atom model over declared tags. This is the
middle path between table-as-opaque-callback and full statechart adoption: keep
RDR 0002's sparse TOML table as the source, but make every guard a symbolic
tag-set predicate that lint can reason about. Authors keep the RDR 0002 shape of
positive `all` guards and negative `unless` guards; each guard entry is a typed
atom: tag reference, fixed operator, and typed literal or literal set. The
operator set is deliberately small: equality, membership, bounded integer
comparison, existence, and set containment. There is no embedded host predicate
and no free-form expression grammar.

Each tag declaration supplies the value kind and, when lint must prove
exhaustiveness, the finite domain: enum values, boolean values, set element
universe, or bounded integer range. Runtime guard evaluation is simple predicate
evaluation over the assembled tag-set. Static lint reasons over the same atoms
by expanding each predicate into a finite domain constraint; it proves whether
candidate rows for an outcome are mutually exclusive and whether they cover the
declared domain.

The initial operator/kind matrix is:

| Operator | Accepted tag value kinds | Literal shape | Lint proof role |
| --- | --- | --- | --- |
| `eq` | enum, boolean, integer, string-like scalar | one typed scalar | Narrows the tag domain to one value. |
| `in` | enum, boolean, integer, string-like scalar | non-empty typed scalar set | Narrows the tag domain to the listed values. |
| `lt`, `lte`, `gt`, `gte` | integer | one typed integer | Narrows a bounded integer domain by comparison; remains runtime-only for an unbounded integer. |
| `exists` | optional scalar or optional set-valued tag | boolean | Tests presence or absence, not value equality. |
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

A predicate atom has four conceptual fields: tag name, operator, expected value,
and source identity. The tag name must resolve to a declared tag. The operator
must be allowed by the operator/kind matrix above. The expected value must parse
to the operator's literal shape. Source identity is inherited from the RDR 0002
row/context so diagnostics can point back to the authored guard.

Positive `all` predicates are conjunctive requirements. Negative `unless`
predicates are also conjunctive within the excluded predicate set: if all
`unless` predicates hold, the candidate row is disabled. Normalization may
combine them into one internal constraint object, but authoring keeps the
separation because it reads better and mirrors prior art.
For lint, a row's accepted assignments are the intersection of all positive
`all` atom domains minus the single conjunctive assignment set matched by the
row's full `unless` block. `unless` is not per-atom negation, and it does not
create source-order priority.

Exhaustiveness is only asserted over finite declared domains. Enum, boolean,
and set-element universes are finite by declaration. Integer tags are exhaustive
only when they declare a bounded range; this covers cap counters and iteration
limits without pretending arbitrary integers can be fully partitioned. A guard
may still compare an unbounded integer at runtime, but lint must report that it
cannot prove exhaustive coverage for that dimension. Provenance affects lint:
recognized tags are fresh event inputs, observed tags are re-read before
matching, and owned tags must have a reachable predecessor write before a row
may match them.

#### Normative Contracts

```normative
A guard predicate MUST be a symbolic atom over a declared tag, not a host
language callback and not a free-form expression string.
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
Overlap and coverage diagnostics MUST name the source rule id or context id
that contributed each predicate involved in the finding.
```

```normative
Predicate lint MUST distinguish owned, observed, and recognized tags. A row
that matches an owned tag MUST be rejected unless every reachable predecessor
sets or preserves that tag before the match.
```

#### Load-Bearing Decisions

- **Identity** — a guard predicate is identified by its source rule/context id
  plus its position within `all` or `unless`; semantic equality is the normalized
  tuple `(tag, operator, literal)`.
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
  declared tag domain, not over examples observed in fixtures. Unbounded
  dimensions remain runtime-evaluable but cannot satisfy an exhaustiveness
  proof.
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
| Accessor read/write safety | RDR 0004 | Pending | Guard evaluation consumes tag values after accessor binding; it does not execute accessors. |
| Graph lint authority | RDR 0006 | Pending | Exhaustiveness and overlap findings become blocking lint there. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| User-facing failures | `internal/cli/clierr`, `internal/cli/respond` | Must verify error-code coverage | Reuse | Predicate parse/lint failures should use existing CLI envelopes. |
| Guard evaluator | None found under `internal/` | New semantic surface | Introduce | This RDR owns the evaluator semantics, not command I/O. |
| Transition model source | Pending RDR 0002 | Not implemented | Extend peer | Guard atoms are embedded in table rows/contexts. |

### Decision Rationale

The fixed symbolic atom model best matches the user's outcome: flow authors can
write conditional edges, and lint can still prove whether those edges are
complete and mutually exclusive. It preserves RDR 0002's readable `all`/`unless`
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

### Risks and Mitigations

- **Risk**: The closed operator set misses a real target-flow guard.
  **Mitigation**: Resolve must encode representative RDR and kata guards before
  lock; add only operators proven necessary by that fixture.
- **Risk**: Finite-domain declarations feel like boilerplate.
  **Mitigation**: Keep declarations close to tag definitions in RDR 0002's
  model and make diagnostics explain when a missing domain blocks proof.
- **Risk**: `unless` semantics are misunderstood as per-atom negation.
  **Mitigation**: Normatively define it as an excluded predicate set and require
  fixtures that demonstrate mixed `all`/`unless` behavior.

### Failure Modes

Visible failures should be typed parse or lint failures: unknown operator,
operator/tag-kind mismatch, literal parse failure, unknown tag, non-exhaustive
finite domain, overlapping candidate rows, or guard dimension not provable
because it lacks a finite domain. Silent failure would be a false
exhaustiveness claim; the recovery path is to keep every exactness claim tied to
A2 and the MVV fixture before Final.

## Implementation Plan

### Prerequisites

- [x] All Critical Assumptions verified
- [x] RDR 0002's sparse table/container contract is stable enough to host guard
  atoms.

### Minimum Viable Validation

Encode one RDR flow slice and one kata flow slice as normalized candidate rows,
including equality, enum membership, set containment, bounded integer
comparison, and mixed `all`/`unless` guards. Lint must prove one exhaustive and
mutually exclusive route, then detect one intentional gap and one intentional
overlap with source rule/context ids in the diagnostic.

### Phase 1: Predicate Model

Define the tag-kind/operator compatibility matrix and normalized predicate atom
shape used by resolver and lint.

### Phase 2: Finite-Domain Lint Semantics

Define how enum, boolean, set-universe, and bounded-int domains are converted
into coverage and overlap checks, including refusal/downgrade behavior for
unbounded dimensions.

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
`exists`. The implementation MVV must extend that coverage with `contains`
before the full closed operator vocabulary is accepted.

1. **Scenario**: Evaluate representative RDR and kata rows that use equality,
   membership, set containment, bounded integer comparison, existence, and mixed
   `all`/`unless` guards.
   **Expected**: Exactly one qualifying row resolves for the legal input; a row
   whose `all` predicates match is disabled when its full conjunctive `unless`
   block also matches; zero or multiple qualifying rows become typed refusals.
2. **Scenario**: Lint finite enum, boolean, set-universe, and bounded-int tag
   domains with one complete partition, one intentional gap, and one intentional
   overlap.
   **Expected**: Complete partitions pass; gaps and overlaps fail with source
   rule/context ids.
3. **Scenario**: Lint an otherwise valid guard over an unbounded integer or
   undeclared finite domain.
   **Expected**: Runtime evaluation remains available, but lint refuses or
   downgrades the exhaustiveness claim for that dimension.
4. **Scenario**: Parse malformed guard atoms: unknown tag, unknown operator,
   unsupported operator/tag-kind pair, and literal parse mismatch.
   **Expected**: Each failure is rejected before resolution with a stable
   structured error category for the CLI gateway.
5. **Scenario**: Reorder authored rows and guard atoms without changing their
   semantics.
   **Expected**: Successful matching and lint findings are unchanged because
   source order is not a selection mechanism; an ambiguous pair remains a
   multiple-match refusal instead of becoming a first-match success.

### Performance Expectations

Resolve measured representative fixture size rather than setting a throughput
target: the spike uses four rows and the target-flow operator subset. The
intended implementation uses local typed comparisons and finite-set expansion
over declared domains; no callback invocation, expression parser, or external
engine is part of the hot path. If implementation later indexes predicates for
speed, the optimization must preserve the normalized atom semantics and
exact-one refusal behavior.

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
past the lens battery undetected. Also confirm form:
value + one clause naming the contract(s); strip any
matrix/provenance prose left from the template or Seed
(it belongs in the template comment, not the instance).]

## References

- RDR 0001, Resolution Kernel Contract.
- RDR 0002, Transition Table As Reviewable Data.
- RDR 0006, Graph Lint Authority and Guarantees.
- `docs/cli-output-contract.md`.
- Resource index: `.rdr/resources.md`.
- Seed prior: `../state-machines/BUILD-SEEDS.md`, especially the guard
  expression and transition-table seeds.
- Transition model prior: "The transition model — inputs, outputs, error
  conditions."
- Prior-art corpus: `../state-machines` audits, evals, contrasts, and checked
  repositories for `transitions`, stateless/qmuntal-stateless, SCXML/scxmlcc,
  Sismic, StateSmith, and Statewright.
