# Recommendation 0008: Resolver disposition identity handoff

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). **Reference-only** (guidance; never copy into the
instance body). -->

## Metadata

- **Date**: 2026-08-06
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
- **Profile**: foundational — one resolver disposition handoff consumed across
  replay validation and CLI projection.
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
- **Related Issues**: kata `zy0m`; RDR 0005; RDR 0007
- **Predecessors**: 0001-resolution-kernel,
  0002-transition-table-as-reviewable-data
- **Seam Lineage**: no prior accretion

## Problem Statement

Operators running flow resolution need diagnostics and audit output to identify
the exact normalized rule selected, including when selection takes a write-free
escape path; they discover the gap when a resolution reports an anonymous plan
or refusal that cannot be traced to its rule or source. Internally, the resolver
must preserve stable normalized rule identity and source location through its
disposition boundary so downstream CLI output and replay tests can consume that
identity unchanged without treating an actionless escape as success.

## Context

### Background

Review of the resolution-kernel, transition-table, and CLI-integration RDRs
found that `Edge`, `TransitionPlan`, and `plan` discard normalized rule identity
and source location. The gap is not exercised by a production caller yet, but
it becomes downstream-reachable when the RDR 0005 `flow resolve` payload is
wired. The design fork is which immutable resolver disposition branch owns
model ID, rule ID, expansion suffix, and source locator while preserving the
schema's distinction between an ordinary row with an action and an escape row
that models a refusal class without an action.

### Technical Environment

The affected seam is the Go resolver in `internal/resolver/resolver.go`, whose
selection result feeds replay validation and the planned CLI resolution payload.
That seam exists on the current integration branch but not yet on `main`; the
RDR 0001 implementation must be present on the implementation baseline before
this RDR can change it.
RDR 0001 defines the resolution kernel requirement to expose downstream values;
RDR 0002 defines normalized model/rule identity and source locators; RDR 0005
defines the consuming `flow resolve` payload contract. Draft RDR 0007 separately
owns normalized-table revision binding; this RDR does not define or validate the
table revision.

## Research Findings

### Investigation

RDRs 0001, 0002, and 0005 are Final. The named resolver seam exists on the
current integration branch at `internal/resolver/resolver.go`, but not yet on
`main`. On that branch, `Edge` carries only match/action data, `plan` copies
only `NextTags` and `Writes`, and `TransitionPlan` cannot satisfy RDR 0005's
selected-rule payload. RDR 0007 owns the separate table-revision binding
question.

RDR 0002 `Normative Contracts` requires every candidate row to retain its
source rule id and source locator, and its `Load-Bearing Decisions / Identity`
defines `(model id, rule id)` plus a deterministic expansion suffix. RDR 0005
`Technical Design` requires `flow resolve` data to contain matched-rule
identity. In surviving local prior art,
qmuntal-stateless `statemachine.go::internalFireOne` constructs a
`Transition{Source, Destination, Trigger}` after handler selection and passes
that value through `handleTransitioningTrigger` to transition callbacks.
scxmlcc `doc/user-manual.md::Transition` models targetless transitions whose
executable content is optional. Those examples support carrying selected
transition context and show that no target does not itself define whether an
action exists; they do not prove the exact selected-rule/source-locator handoff
chosen here. That bounded contrast does not justify synthesizing an action for
an RDR 0002 escape row that names a refusal class and forbids `write` and
`clear`. Queries and rejected branches are recorded under
`docs/rdr/0008-resolver-disposition-identity-handoff/evidence/research/`.

Sibling-path check: source search found no adjacent selected-rule carrier or
identity discriminator in `internal/resolver`; the existing decision boundary
is `internal/resolver/resolver.go::plan`, so this proposal extends that boundary
instead of inventing a parallel CLI-side lookup.

### Key Discoveries

- **Documented** — RDR 0002 already owns the complete normalized selection
  identity; this RDR only owns preserving it across resolver dispositions.
- **Documented** — RDR 0005 consumes matched-rule identity but does not own
  reconstructing it from model data.
- **Prior art, bounded** — qmuntal-stateless carries source, destination, and
  trigger from selection through transition callbacks. This supports direct
  transition-context handoff but not this RDR's exact identity fields or
  branch taxonomy.
- **Verified on the integration branch** — `Edge`, `TransitionPlan`, `Refusal`,
  and their constructors discard all four identity fields. `Resolve` currently
  turns an exactly-one escape into a plan even though RDR 0002 forbids action
  fields on escape rules. The symbols do not yet resolve on `main`; landing the
  RDR 0001 implementation is an explicit implementation prerequisite below.
- **Verified** — RDR 0002 requires the future normalizer to populate the
  complete selected-rule value on every normalized ordinary and escape row,
  without a downstream lookup.

### Critical Assumptions

- **A1 Every normalized ordinary and escape edge can carry the complete RDR
  0002 selection identity without a downstream model lookup.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0002 `Normative Contracts` requires every normalized
    candidate row to retain source rule id and source locator;
    `Load-Bearing Decisions / Identity` defines the same row's identity as
    model id, rule id, and deterministic expansion suffix, and its modeled-
    escape contract keeps escape selection in that normalized row set.
  - **If wrong**: The resolver cannot make identity and the selected row's
    ordinary action or modeled refusal one coherent result, so CLI diagnostics
    may report a different rule than the one used.
- **A2 Copying a selected disposition can preserve selection identity, typed
  source locator, and any action without aliasing mutable normalized-table
  storage.**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: `TestPlanOwnsSelectedIdentityAndAction` is green and
    proves isolation for string identity fields, a string locator, a cloned
    next-tag map, and a copied writes slice. It remains a local surrogate:
    production `internal/resolver/resolver.go::{Edge,Input}` has no
    `SelectedRule` or typed locator and still accepts caller-built `[]Edge`.
    Before Resolve can pass, rerun the spike against the concrete A4/A5 types
    for ordinary and modeled-escape selection, mutating constructor inputs
    before resolution and source or inspection storage afterward; an external-
    package API test must show row components cannot be independently
    constructed, replaced, or mutably exposed.
  - **If wrong**: Audit output or replay assertions could drift after the
    caller reuses or mutates the supplied table.
- **A3 RDR 0005 can require and test direct projection of the selected-rule
  value without reconstructing identity from flow/model inputs.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: Final RDR 0005 `Technical Design` includes matched rule
    identity and action in the `flow resolve` payload, but its `Normative
    Contracts` require only a resolved next tag-set or refusal. Before Resolve
    can pass, RDR 0005 must require `SelectedRule` projection on ordinary and
    modeled-refusal paths, and `Action` projection only on ordinary success, in
    JSON and text acceptance tests without a candidate-table lookup or accessor
    execution.
  - **If wrong**: This RDR must revise the carrier before lock or RDR 0005 will
    omit identity or retain a parallel identity lookup.
- **A4 The normalized-table boundary can construct a typed source locator that
  remains bound to the selected model and rule.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: RDR 0002 `Technical Design` says a locator identifies the
    model and rule and its `Normative Contracts` require retention, but neither
    section defines a typed constructor or callable validator, zero-value
    behavior, model/rule mismatch behavior, or a boundary test. Before Resolve
    can pass, that contract must be named at the normalized-table construction
    boundary; `plan` remains only a copying consumer of the validated value.
  - **If wrong**: A plan or modeled refusal could report the selected identity
    while directing diagnostics to a different authored rule.
- **A5 The production resolver entry point can accept only normalized-table
  rows that bind match inputs and selected identity to the row-kind payload at
  one validated construction boundary.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: Current `internal/resolver/resolver.go::{Input,Edge}` exposes a
    caller-constructible `[]Edge`; RDR 0007 is Draft and explicitly records that
    no production normalizer or validated-table type exists. Before Resolve can
    pass, the owning RDR must name the concrete constructor or validator and
    production input type that prevent callers from composing rule B identity
    with rule A predicates, ordinary action, or escape failure classes, plus
    ordinary and escape mismatch tests.
  - **If wrong**: The resolver can execute one row's action or report one row's
    refusal classes while returning authoritative-looking identity for another
    row.
- **A6 An exactly-one modeled escape can return its underlying refusal kind
  with selected-rule identity and no action across the source schema,
  normalized row, resolver, and CLI.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: RDR 0002 `Normative Contracts` make an escape row name
    `no_match` or `ambiguous_match` and forbid `write` and `clear`; RDR 0005
    permits `flow resolve` to return a refusal. RDR 0001 `Technical Design` and
    `internal/resolver/resolver.go::Resolve` currently turn an exactly-one
    escape into a successful plan, so Resolve must verify that the peer
    contracts can converge on an identity-bearing refusal and update their
    acceptance tests before this RDR locks.
  - **If wrong**: The source schema must gain an explicit complete no-change
    action or this RDR needs a distinct modeled-escape disposition; `plan` must
    not invent next tags from an actionless row.

## Proposed Solution

### Approach

Carry one row-owned `SelectedRule` into the disposition branch that selection
actually produces. An exact-one ordinary selection returns a `TransitionPlan`
with that identity and a separately copyable `TransitionAction`. An exact-one
modeled escape returns the underlying typed `Refusal` with the same identity
and no action. A kernel refusal for which no row was selected has no selected
rule. This keeps "selected" distinct from "successful": diagnostics can name
the authored escape without pretending that an actionless row produced a next
state.

RDR 0008 owns only this Go disposition handoff. RDR 0002 owns identity semantics
and must supply typed source-locator validation, coherent normalized-row
construction, and snapshot ownership; RDR 0005 owns CLI projection; RDR 0007
owns table-revision binding. The handoff copies the `SelectedRule` and, for
ordinary selection, the `TransitionAction` directly from the selected normalized
row. It never reconstructs identity or exposes matching internals.

### Technical Design

The resolver consumes the tagged normalized row supplied by the RDR 0002
boundary: common match inputs and `SelectedRule`, plus either an ordinary
`Action` or modeled refusal classes. This RDR relies on that boundary's
validation and snapshot ownership; it does not define its constructor or
inspection API.

`SelectedRule` has `ModelID`, `RuleID`, `ExpansionSuffix`, and the typed
`SourceLocator` owned by RDR 0002. `TransitionAction` has `NextTags TagSet` and
`Writes []OwnedTagWrite`. `TransitionPlan` exposes both values. `Refusal`
exposes `Kind` and an optional `SelectedRule`, present only for a selected
modeled escape. Resolver matching reads only match inputs. `plan` copies the
ordinary row's selected rule and action; `modeledRefusal` copies the escape
row's selected rule and handled `no_match` or `ambiguous_match` kind; kernel
`refuse` attaches no identity. A2/A4 determine whether the typed locator is
deeply immutable and value-copyable or requires an explicit clone.

`plan` and `modeledRefusal` treat an empty model id or rule id, or a zero or
invalid source locator, as an internal invariant error rather than a modeled
refusal. An expansion suffix may be empty for an unexpanded rule but still
participates in identity. The handoff copies the typed locator without parsing
or reformatting it. RDR 0005 projects the resulting selection and successful
action; user-authored validation failures remain load/config failures under its
CLI contract. A later state operation may pass only a successful action to the
RDR 0004 accessor boundary.

The selection value is descriptive, not a second authority. It identifies the
matched logical row within the revision-bound resolver input; it neither
recalculates table revision nor redefines identity, matching semantics, or CLI
rendering.

The existing RDR 0001 replay test is the replay consumer for this handoff. This
RDR extends that test so two resolutions of the same input compare the complete
disposition: ordinary plans compare `SelectedRule` and `Action`; modeled
refusals compare `Kind` and `SelectedRule`; kernel refusals compare `Kind` and
the absence of selection. Replay does not gain a second identity lookup.

#### Normative Contracts

```normative
Given one selected row from the validated normalized-table boundary, the
resolver MUST copy that row's `SelectedRule` into the disposition branch that
selection produces. `SelectedRule` contains `ModelID`, `RuleID`,
`ExpansionSuffix`, and the typed `SourceLocator`. `ModelID` and `RuleID` MUST be
non-empty and `SourceLocator` MUST be valid; a missing required field MUST use
the Go programmer/invariant error path, not a modeled refusal.

An exact-one ordinary selection MUST return exactly one `TransitionPlan` with
the matched row's `SelectedRule` and `Action`.

An ordinary transition plan MUST expose `Action` separately from `SelectedRule`.
`Action` contains `NextTags` and `Writes`; callers MUST NOT need guard predicates
or the matched `Edge` to apply it.

When zero or multiple ordinary rows match and exactly one eligible escape row
matches, the resolver MUST return the underlying `no_match` or
`ambiguous_match` `Refusal` with that row's `SelectedRule` and MUST NOT return a
`TransitionPlan`, `NextTags`, or `Writes`. A kernel refusal for which no escape
row was selected MUST NOT claim selected-rule identity. The `Refusal`
constructor and tests MUST enforce that selected-rule presence means exactly
one modeled escape row was selected; no other refusal may claim it.

The returned `SelectedRule` and any `Action` MUST be owned values that remain
unchanged after `Resolve` when producer or inspection storage is reused or
mutated.
```

Before this RDR locks, RDR 0002's validated-table contract must separately
prevent callers from composing or replacing row identity and payload, and RDR
0005's consumer contract must preserve the disposition in text and JSON without
a candidate-table lookup or accessor execution.

#### Load-Bearing Decisions

- **Identity** — selected-rule equality within one revision-bound resolver
  input is the RDR 0002 tuple `(model id, rule id, expansion suffix)`; source
  locator is carried diagnostic provenance and does not create a different
  logical rule identity. RDR 0007 owns durable table revision and its successor
  owns binding that revision to the normalized table; this carrier is not a
  standalone historical audit-event identity.
- **Naming** — the successful carrier remains `TransitionPlan`; the nested
  types are `SelectedRule` and `TransitionAction`, exposed as
  `TransitionPlan.SelectedRule` and `TransitionPlan.Action`. A modeled escape
  remains a `Refusal` and exposes the same type as `Refusal.SelectedRule`.
  Rejected: "provenance" alone because the value is operational selection
  identity, and "edge" because callers must not receive matching internals.
- **Selection / predicate** — identity is copied only from the exact edge RDR
  0001 already selects; this RDR adds no tie-break or second lookup.
- **Boundary ownership** — this RDR defines the disposition carrier and branch
  invariants only. RDR 0002 owns producer validation and RDR 0005 owns terminal
  projection; their acceptance is a lock prerequisite, not an additional
  contract authored here.

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Exact-one resolver disposition and escape meaning | RDR 0001 / RDR 0002 | A6 Pending | Must align peer contracts on an identity-bearing refusal with no action. |
| Validated normalized rows, rule identity, and typed source locator | RDR 0002 | Predecessor; A2/A4/A5 Pending | Must supply one snapshot-owned row value; this RDR consumes but does not specify its construction API. |
| Table-revision derivation and binding | RDR 0007 + charted successor | Deferred peers | Remain outside this logical identity-handoff contract. |
| Matched-rule CLI payload | RDR 0005 | Predecessor; A3/A6 Pending | Must project selected-rule identity on success and modeled refusal, with action only on success. |
| Disposition identity handoff | This RDR | Introduced | Extends ordinary plans and modeled refusals across the producer/consumer seam. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Normalized candidate | `internal/resolver/resolver.go::{Input,Edge}` | Exported raw rows and backing slices cannot provide the required coherent selected value | Predecessor change | RDR 0002/A5 must supply an opaque, snapshot-owned row; this RDR consumes it. |
| Successful result | `internal/resolver/resolver.go::TransitionPlan` | Carries action only | Extend | Own selected rule and separately usable action. |
| Modeled refusal | `internal/resolver/resolver.go::{Refusal,Resolve}` | Exactly-one escape is converted to an action-bearing plan | Correct | Return the underlying refusal kind with selected-rule identity and no action. |
| Result construction | `internal/resolver/resolver.go::{plan,refuse}` | Drops identity on every branch | Extend | Add defensive ordinary-plan and modeled-refusal construction paths. |
| CLI output | `internal/cli/respond` and RDR 0005 | No resolver verb wired yet | Reuse later | Rendering stays outside the kernel. |

### Decision Rationale

The Questions-Options-Criteria matrix scores 1 (poor) through 5 (strong):

| Approach | Correctness fit | Prior-art alignment | Reversibility | Blast radius | Cost | Total |
| --- | --- | --- | --- | --- | --- | ---: |
| Plan on ordinary selection; identity-bearing Refusal on escape | 5 — action exists only where authored | 4 — transition context is handed through directly; identity/refusal split is a local decision | 4 — nested carriers can evolve | 4 — resolver result and peer tests | 4 — two explicit constructors | **21** |
| Successful plan with explicit no-change escape action | 2 — invents action absent from source | 3 — SCXML permits targetless success, but with transition semantics | 2 — changes schema and consumers | 2 — spans parser through CLI | 2 — new action variant and projection | **11** |
| Dedicated ModeledEscape disposition branch | 5 — invalid combinations are unrepresentable | 3 — explicit but adds a project-specific third branch | 3 — removable through migration | 2 — every disposition consumer changes | 3 — new branch and mapping | **16** |
| Disposition-level optional selected rule | 3 — permits plan/identity mismatch | 4 — context stays near result | 5 — leaves branch types unchanged | 4 — one extra outer field | 3 — every branch enforces presence | **19** |

The split carrier wins on correctness and local design fit: ordinary rows
already author actions, while escape rows author refusal classes and forbid
actions. It preserves the direct transition-context handoff seen in local prior
art without treating that prior art as proof of the selected-rule schema or
adding a third public disposition branch. A successful no-change plan
would require RDR 0002 to author a real action rather than letting `plan` infer
one. A dedicated branch makes the invariant clearer in the type shape but
forces every RDR 0001/RDR 0005 consumer to learn a third terminal outcome.
Putting selection on the outer disposition admits invalid combinations and
separates ordinary identity from its action.

Premortem: this ships and diagnostics sometimes show an empty or stale rule id,
or a selected escape is accidentally rendered as success because one constructor
copies action but not identity and another attaches identity to the wrong
refusal. The recommendation survives: separate `plan`, `modeledRefusal`, and
kernel `refuse` constructors give each branch one invariant, while A2's MVV
requires snapshot ownership and defensive-copy behavior. The parser-to-CLI
escape test makes an invented action or anonymous modeled refusal observable.

## Alternatives Considered

### Alternative 1: Successful No-Change Escape Action

Treat an exactly-one escape as a successful `TransitionPlan` with an explicit
no-change action. scxmlcc's SCXML model shows that a targetless transition can
still have explicitly modeled executable content. Intrastate's escape schema
instead names a refusal class and forbids `write` and `clear`; adopting this
alternative would require a source and normalized action variant, not empty or
input-derived `NextTags` in `plan`.

### Alternative 2: Dedicated Modeled-Escape Disposition

Add a third disposition branch containing the handled refusal kind and
`SelectedRule`, separate from both `TransitionPlan` and `Refusal`. This makes
its shape unambiguous but turns a modeled form of an existing refusal into a
new terminal taxonomy that every resolver and CLI consumer must handle.

### Alternative 3: Selection on Disposition

Add an optional selected-rule value beside `Plan` and `Refusal`, leaving both
branch values unchanged. This permits plan without selection, kernel refusal
with selection, and selection without either branch; placing identity inside
the branch whose constructor selected the row gives a narrower invariant.

### Briefly Rejected

- **CLI-side identity reconstruction**: Rejected because a second model lookup
  can report a row other than the one the kernel selected.
- **Return the matched `Edge`**: Rejected because it exposes guards and matching
  internals, risks aliases, and couples consumers to candidate storage.

## Trade-offs

### Consequences

- Ordinary success and selected modeled refusal both become self-describing for
  diagnostics and replay assertions.
- Action consumers can apply next tags/writes without depending on rule identity
  or guard internals.
- `TransitionPlan`, `Refusal`, their constructors, and replay fixtures must grow
  in lockstep; incomplete selected-rule identity becomes an invariant error.

### Risks and Mitigations

- **Risk**: One selected path omits or partially fills selected-rule identity.
  **Mitigation**: Use one ordinary plan constructor and one modeled-refusal
  constructor; make the MVV cover both plus anonymous kernel refusals.
- **Risk**: The producer supplies mismatched identity, locator, or row payload.
  **Mitigation**: Make RDR 0002's validated-table acceptance an A4/A5 lock
  prerequisite; the resolver never repairs or reconstructs the value.
- **Risk**: Returned maps or slices alias producer storage.
  **Mitigation**: Clone the action and clone or value-copy identity fields as A2
  requires; mutate inspection storage after resolution and exercise concurrent
  snapshot use under the race detector.
- **Risk**: A selected modeled escape is rendered as success or gains an
  invented next state.
  **Mitigation**: Reconcile RDR 0001, RDR 0002, and RDR 0005 at A6; test
  authored TOML through production normalization, resolution, and CLI
  projection.
- **Risk**: This RDR accidentally absorbs table revision or CLI wire-format
  ownership.
  **Mitigation**: Keep revision validation in RDR 0007 and rendering in RDR 0005.

### Failure Modes

- **Incomplete identity on a selected edge**: `plan` or `modeledRefusal` returns
  a programmer error before producing a disposition; the MVV covers every
  required field. Diagnose at the matching resolver constructor.
- **Invalid producer row**: identity, locator, predicates, and payload do not
  describe one rule. The RDR 0002 boundary rejects it before selection; `plan`
  never guesses or repairs provenance.
- **Aliased result**: replay output changes after table reuse. The mutation
  fixture identifies which selected-rule/action field was not copied.
- **Undefined escape action**: a selected escape returns success with empty or
  invented next tags. Reject it in `Resolve`; the production parser-to-CLI test
  requires an identity-bearing refusal with no action.
- **Kernel refusal claims a rule**: a no-match or ambiguous refusal falsely
  implies an escape selection. Reject the invalid disposition in tests;
  `refuse` keeps identity absent and only `modeledRefusal` may populate it.

## Implementation Plan

### Prerequisites

- [ ] A2 through A6 are verified and reconciled.
- [ ] RDR 0001's resolver implementation is present on the implementation
  baseline; the reviewed symbols currently exist only on the integration branch.
- [ ] RDR 0002 supplies the A4/A5 typed locator and validated normalized-row
  boundary, including mismatch and storage-ownership acceptance tests.
- [ ] RDRs 0001, 0002, and 0005 align under A3/A6 on an identity-bearing modeled
  refusal with no action; RDR 0005 acceptance directly projects ordinary and
  modeled-refusal dispositions in text and JSON.

### Minimum Viable Validation

`TestResolvePreservesSelectedRuleAndAction` resolves a production normalized
ordinary row and asserts that the returned selected-rule tuple and action match
that exact row. The fixture carries a non-empty expansion suffix and write, then
mutates exposed inspection storage after `Resolve` and asserts the returned plan
is unchanged. A companion resolver test selects a normalized escape and asserts
its failure kind and complete selected-rule identity with no plan or action.
RDR 0002 separately proves producer coherence and storage ownership; RDR 0005
separately proves terminal projection.

### Phase 1: Adopt the Producer Selection Value

Wire the RDR 0002 `SelectedRule` value into the resolver-facing normalized row.
Ordinary rows supply it beside an action; escape rows supply it beside modeled
refusal classes. Do not introduce a parallel identity type or a selection-time
conversion.

### Phase 2: Preserve the Disposition Handoff

Extend the ordinary plan-construction path to return defensive copies of
selected rule and action without changing matching semantics. Add a modeled-
refusal constructor that returns the selected escape's underlying failure kind
and defensive selected-rule copy without inventing action data in `plan`.

## Validation

### Testing Strategy

The source audit grounds ordinary and modeled-escape selection in
`internal/resolver/resolver.go::Resolve` and the current shared construction
boundary in `internal/resolver/resolver.go::plan`; the proposal splits that
boundary into ordinary-plan, modeled-refusal, and kernel-refusal constructors.
The A2 spike under `evidence/spikes/` grounds only string-identity/action
copying. RDR 0008 implementation acceptance covers:

1. **Ordinary selection** — resolve an exact-one ordinary row through
   `Resolve` and `plan`.
   **Expected**: `TransitionPlan.SelectedRule` equals that row's complete
   identity, including a non-empty expansion suffix, and `Action` equals its
   next tags and writes.
2. **Modeled escape** — resolve a normalized escape with no action.
   **Expected**: the result is the modeled `no_match` or `ambiguous_match`
   refusal with the escape row's complete `SelectedRule`, no plan, and no
   action.
3. **Copy isolation** — mutate any storage exposed for inspection after
   resolution and concurrently resolve one stable snapshot under `go test
   -race`.
   **Expected**: the returned selected rule and action remain unchanged.
4. **Incomplete identity** — inject an impossible empty model id, rule id, or
   invalid locator inside a validated row in a package test.
   **Expected**: the programmer/invariant error path returns no disposition.
5. **Refusal branches** — exercise no match, ambiguous match, unavailable owned
   state, unevaluable guard, and unmodeled outcome without a selected escape.
   **Expected**: each disposition has no plan and no selected rule; only the
   modeled-escape fixture may attach selected identity to a refusal.
6. **Replay identity** — extend RDR 0001's
   `TestResolve_ReplayReturnsValueIdenticalPlans` with ordinary and escape
   selected identity.
   **Expected**: two resolutions of the same input return value-identical
   dispositions: complete `SelectedRule` and `Action` for ordinary success,
   complete `SelectedRule` and `Kind` for modeled refusal, without a second
   identity lookup.

Lock additionally requires peer acceptance: RDR 0002 rejects mismatched
identity/locator/payload combinations and mutable backing storage before
selection; RDR 0005 projects ordinary and modeled-refusal dispositions in JSON
and text without lookup or accessor execution. Those tests verify dependencies,
not code owned by RDR 0008.

### Performance Expectations

The existing A2 spike supports bounded copying for string identity and action
fields: the next-tag map and writes slice are cloned once in `plan`, so time and
allocation scale linearly with the selected action's tag and write counts. A2
must verify the concrete typed locator and snapshot-owned table before lock.
The spike is a correctness proof, not a latency benchmark; this RDR makes no
throughput claim. Add a benchmark only if real normalized actions make copy cost
material.

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

- RDR 0001 `Normative Contracts` — resolver disposition and refusal ownership.
- RDR 0002 `Normative Contracts`; `Load-Bearing Decisions / Identity` — source
  rule id, source locator, model id, and expansion suffix.
- RDR 0005 `Technical Design`; `Normative Contracts` — matched-rule CLI payload.
- RDR 0007 — normalized-table revision binding, explicitly separate.
- qmuntal-stateless `statemachine.go::internalFireOne` and
  `::handleTransitioningTrigger` — direct transition-context handoff.
- scxmlcc `doc/user-manual.md::Transition` and
  `::Custom Actions and Conditions` — targetless-transition/action contrast.
