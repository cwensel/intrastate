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
- **Profile**: foundational — one cross-RDR selected-rule/action handoff
  contract spanning the resolver, replay validation, and CLI projection.
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

The freshness check confirmed that RDRs 0001, 0002, and 0005 remain Final. The
named resolver seam exists on the current integration branch at
`internal/resolver/resolver.go`, but not yet on `main`. On that branch, `Edge`
carries only match/action data, `plan` copies only `NextTags` and `Writes`, and
`TransitionPlan` therefore cannot satisfy RDR 0005's selected-rule payload.
RDR 0007 now owns the adjacent table-revision binding question and is explicitly
outside this RDR.

Prior art was read before alternatives were revised. RDR 0002 `Normative
Contracts` requires: "Each candidate row MUST retain its source rule id and
source locator." Its `Load-Bearing Decisions / Identity` defines `(model id,
rule id)` plus a deterministic expansion suffix. RDR 0005 `Technical Design`
requires `flow resolve` data to contain "matched rule identity." In external
source, Stateless `src/Stateless/Transition.cs::Transition` and uscxml
`src/uscxml/interpreter/LargeMicroStep.h::Transition` retain selected-transition
context through successful execution and monitoring, favoring a direct carrier
over downstream reconstruction. W3C SCXML 1.0 `3.1.5 Type and Transitions`
says a targetless transition does not change the state configuration but "does
invoke the executable content" in the transition. That contrast supports a
no-state-change success only when transition behavior is explicitly modeled;
it does not justify synthesizing an action for an RDR 0002 escape row that names
a refusal class and forbids `write` and `clear`. Queries and rejected branches
are recorded under
`docs/rdr/0008-resolver-disposition-identity-handoff/evidence/research/`.

Sibling-path check: source search found no adjacent selected-rule carrier or
identity discriminator in `internal/resolver`; the existing decision boundary
is `internal/resolver/resolver.go::plan`, so this proposal extends that boundary
instead of inventing a parallel CLI-side lookup.

### Key Discoveries

- **Documented** — RDR 0002 already owns the complete normalized selection
  identity; this RDR only owns preserving it across resolver success.
- **Documented** — RDR 0005 consumes matched-rule identity but does not own
  reconstructing it from model data.
- **Documented** — Stateless and uscxml keep selected-transition context on the
  successful transition object used by execution and diagnostics.
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
  - **Method**: Spike + Peer RDR
  - **Evidence**: `TestPlanOwnsSelectedIdentityAndAction` reruns green and
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
  - **Method**: Source Search + Peer RDR
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
  - **Method**: Peer RDR + Source Search
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

Make the resolver disposition preserve selected-rule identity on the semantic
branch that actually occurred. An ordinary exact-one selection returns a
`TransitionPlan` containing a `SelectedRule` and a separately copyable
`TransitionAction`. An exactly-one modeled escape returns the underlying typed
`Refusal` with `SelectedRule` populated and no action. A kernel refusal for which
no row was selected has no selected rule. This keeps "selected" distinct from
"successful": diagnostics can name the authored escape without pretending that
an actionless row produced a next state.

RDR 0008 owns this Go handoff shape; RDR 0002 continues to own identity
semantics and normalized-table population. `SelectedRule` is the one value
shared by a validated normalized row and whichever disposition branch consumes
that row, not a second resolver projection or a value reconstructed during
selection. It contains model id, rule id, expansion suffix, and a typed source
locator. `TransitionAction` contains next tags and owned-tag writes and exists
only on an ordinary row and its successful plan. The normalized-table
constructor owns caller storage and builds a tagged row: common match inputs and
selected rule plus either an ordinary action or modeled refusal classes. The
production resolver accepts one stable validated-table snapshot; `plan` copies
the ordinary row's identity and action, while `modeledRefusal` copies the escape
row's identity and the failure kind it handled. Downstream callers consume
selection for diagnostics/audit and successful action, when present, without
receiving guard or candidate-table internals.

### Technical Design

The normalized-table boundary constructs an opaque, snapshot-owned validated
table whose internal tagged rows bind match inputs and `SelectedRule` to exactly
one row payload: `Action TransitionAction` for an ordinary row or a modeled
failure-class set for an escape row. Production callers cannot independently
assemble, replace, or mutate those row components after validation; inspection
and dump APIs return immutable values or defensive copies. `SelectedRule` has
string fields `ModelID`, `RuleID`, and `ExpansionSuffix`, plus a typed
`SourceLocator` value owned by the RDR 0002 normalization boundary.
`TransitionAction` has `NextTags TagSet` and `Writes []OwnedTagWrite`.
`TransitionPlan` exposes the selected rule and action. `Refusal` exposes its
`Kind` and an optional `SelectedRule`; that field MUST be present exactly when a
modeled escape row was selected. Resolver matching reads only match inputs. On
exact-one ordinary selection, `plan` copies the matched row's selected rule and
action without retaining mutable aliases. On exact-one escape selection,
`modeledRefusal` copies the matched row's selected rule and the underlying
`no_match` or `ambiguous_match` kind without constructing an action. A2/A4 must
determine whether the typed locator is deeply immutable and value-copyable or
requires an explicit clone.

Before returning a selected disposition, `plan` and `modeledRefusal` treat an
empty model id or rule id, or a zero or invalid source locator, as a programmer/
invariant error on the Go error path, not as a modeled refusal. An expansion
suffix may be empty for an unexpanded source rule; its value still participates
in identity. The source locator is an RDR 0002 typed value that this handoff
copies without parsing or reformatting. Before a row reaches the resolver, the
normalized-table boundary constructs the locator with the same model id and
rule id carried by `SelectedRule`; both constructors enforce required presence
but do not reconstruct or reinterpret provenance. RDR 0005 must project the
selected rule and any action carried by the resolved disposition into
`flow resolve` output without
reconstructing identity or executing accessors. A modeled escape maps to the
stable CLI refusal for its `Kind` while retaining matched-rule identity in the
RDR 0005 refusal payload; it does not claim next tags or writes. User-authored
model validation failures are load/config failures that RDR 0005 maps to a
stable `CLIError`; an invariant breach in an already validated table remains an
internal Go error. A later state operation may pass only a successful action to
the RDR 0004 accessor boundary.

The selection value is descriptive, not a second authority: it is copied from
the normalized row that was actually matched and identifies that logical row
within the revision-bound resolver input. It is not a standalone historical
event identity. It does not recalculate or enforce table revision (RDR 0007 and
its charted enforcement successor), redefine row identity (RDR 0002), change
matching/refusal semantics (RDR 0001), or define CLI rendering (RDR 0005).

The existing RDR 0001 replay test is the replay consumer for this handoff. This
RDR extends that test so two resolutions of the same input compare the complete
disposition: ordinary plans compare `SelectedRule` and `Action`; modeled
refusals compare `Kind` and `SelectedRule`; kernel refusals compare `Kind` and
the absence of selection. Replay does not gain a second identity lookup.

#### Normative Contracts

```normative
The production resolver entry point MUST consume one validated normalized-table
value whose construction binds each row's match inputs and `SelectedRule` to
exactly one tagged payload: `Action` for an ordinary row or modeled refusal
classes for an escape row. Production callers MUST NOT be able to compose,
replace, or mutate those row components after validation, and inspection APIs
MUST NOT expose mutable backing storage.
Invalid user-authored model construction MUST fail as a load/config error before
selection; an impossible value observed inside an already validated table MUST
use the internal programmer/invariant error path.

Every successful resolver disposition MUST expose exactly one `TransitionPlan`
whose `SelectedRule` contains `ModelID`, `RuleID`, `ExpansionSuffix`, and
the typed `SourceLocator` copied from the normalized row that matched.
`ModelID` and `RuleID` MUST be non-empty and `SourceLocator` MUST be valid; a
missing required field MUST use the Go error path for a programmer/invariant
failure, not a modeled refusal. The normalized-table boundary MUST reject a row
whose `SourceLocator` designates a model id or rule id different from that
row's `SelectedRule`.

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

The returned `SelectedRule` and any `Action` MUST remain unchanged if the
caller mutates or reuses constructor inputs before `Resolve`, later reuses
source storage, or inspects the validated table through a public API.

`flow resolve` MUST project the resolver disposition's `SelectedRule` and any
successful `Action` in both JSON and text terminal output without a
candidate-table lookup or accessor execution. A modeled escape MUST map to the
stable CLI refusal for its `Kind`, include its selected-rule identity, and omit
next tags and writes. RDR 0005 owns rendering, stable load/config error codes,
and CLI wiring, but its consumer contract and acceptance tests MUST preserve
this handoff before this RDR locks.
```

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

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Exact-one resolver disposition and escape meaning | RDR 0001 / RDR 0002 | A6 Pending | Must align peer contracts on an identity-bearing refusal with no action. |
| Validated normalized rows, rule identity, and typed source locator | RDR 0002 | Predecessor; A2/A4/A5 Pending | Supplies one snapshot-owned row value; this RDR does not redefine identity semantics. |
| Table-revision derivation and binding | RDR 0007 + charted successor | Deferred peers | Remain outside this logical identity-handoff contract. |
| Matched-rule CLI payload | RDR 0005 | Predecessor; A3/A6 Pending | Must project selected-rule identity on success and modeled refusal, with action only on success. |
| Disposition identity handoff | This RDR | Introduced | Extends ordinary plans and modeled refusals across the producer/consumer seam. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Normalized candidate | `internal/resolver/resolver.go::{Input,Edge}` | Exported raw rows and backing slices let callers compose or mutate unrelated identity and row payload | Replace boundary | Accept an opaque, snapshot-owned normalized-table value; return copies from inspection APIs. |
| Successful result | `internal/resolver/resolver.go::TransitionPlan` | Carries action only | Extend | Own selected rule and separately usable action. |
| Modeled refusal | `internal/resolver/resolver.go::{Refusal,Resolve}` | Exactly-one escape is converted to an action-bearing plan | Correct | Return the underlying refusal kind with selected-rule identity and no action. |
| Result construction | `internal/resolver/resolver.go::{plan,refuse}` | Drops identity on every branch | Extend | Add defensive ordinary-plan and modeled-refusal construction paths. |
| CLI output | `internal/cli/respond` and RDR 0005 | No resolver verb wired yet | Reuse later | Rendering stays outside the kernel. |

### Decision Rationale

The Questions-Options-Criteria matrix scores 1 (poor) through 5 (strong):

| Approach | Correctness fit | Prior-art alignment | Reversibility | Blast radius | Cost | Total |
| --- | --- | --- | --- | --- | --- | ---: |
| Plan on ordinary selection; identity-bearing Refusal on escape | 5 — action exists only where authored | 5 — direct selected-object handoff; preserves local refusal model | 4 — nested carriers can evolve | 4 — resolver result and peer tests | 4 — two explicit constructors | **22** |
| Successful plan with explicit no-change escape action | 2 — invents action absent from source | 3 — SCXML permits targetless success, but with transition semantics | 2 — changes schema and consumers | 2 — spans parser through CLI | 2 — new action variant and projection | **11** |
| Dedicated ModeledEscape disposition branch | 5 — invalid combinations are unrepresentable | 3 — explicit but adds a project-specific third branch | 3 — removable through migration | 2 — every disposition consumer changes | 3 — new branch and mapping | **16** |
| Disposition-level optional selected rule | 3 — permits plan/identity mismatch | 4 — context stays near result | 5 — leaves branch types unchanged | 4 — one extra outer field | 3 — every branch enforces presence | **19** |

The split carrier wins on correctness and local prior-art alignment: ordinary
rows already author actions, while escape rows author refusal classes and
forbid actions. It preserves the direct selected-object handoff on both paths
without adding a third public disposition branch. A successful no-change plan
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
no-change action. W3C SCXML shows that a targetless transition can retain
transition semantics without changing state, but it still executes explicitly
modeled content. Intrastate's escape schema instead names a refusal class and
forbids `write` and `clear`; adopting this alternative would require a source
and normalized action variant, not empty or input-derived `NextTags` in `plan`.

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
- `Edge`, `TransitionPlan`, `Refusal`, construction code, and fixtures must all
  grow in lockstep; incomplete selected-rule identity becomes a programmer
  error to reject and test.

### Risks and Mitigations

- **Risk**: One selected path omits or partially fills selected-rule identity.
  **Mitigation**: Use one ordinary plan constructor and one modeled-refusal
  constructor; make the MVV cover both plus anonymous kernel refusals.
- **Risk**: A caller pairs one row's identity or locator with another row's
  predicates or row-kind payload.
  **Mitigation**: Make `Resolve` consume only the validated normalized-table
  value and test that malformed ordinary and escape rows cannot produce a
  disposition.
- **Risk**: Returned maps or slices alias caller-owned edge data.
  **Mitigation**: Make validated-table construction own caller storage; clone
  the action and clone or value-copy only deeply immutable identity fields.
  Mutate source inputs before and after resolution and run concurrent snapshot
  use under the race detector.
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
- **Identity/provenance mismatch on a normalized edge**: normalization rejects
  the row before resolver selection. Diagnose at the normalized-table
  validation boundary; `plan` must not guess or repair provenance.
- **Identity/action mismatch**: audit names one rule while next tags/writes came
  from another. Reject it at the validated normalized-table construction
  boundary before `Resolve`; do not reconstruct or guess downstream.
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

- [ ] All Critical Assumptions reconciled; A2 through A6 are Pending for Stage
  6.
- [ ] RDR 0001's resolver implementation is present on the implementation
  baseline; the reviewed symbols currently exist only on the integration branch.
- [ ] RDR 0002 selected-rule fields are available at normalized edge creation.
- [ ] A4 names the RDR 0002 typed locator constructor and its model/rule
  agreement test.
- [ ] A5 names the validated normalized-table production type/constructor and
  proves callers cannot independently compose or mutate row identity and its
  ordinary-action or escape-refusal payload.
- [ ] A6 aligns RDR 0001, RDR 0002, and RDR 0005 on the selected escape as an
  identity-bearing refusal with no action, including production normalization,
  resolver tests, and CLI projection.
- [ ] RDR 0007 and its charted enforcement successor own revision derivation and
  table binding so this RDR does not duplicate historical event identity.
- [ ] RDR 0005 `Normative Contracts` and text/JSON acceptance require direct
  projection of the returned `SelectedRule` and any successful `Action`, plus a
  stable load/config `CLIError` for malformed authored identity.

### Minimum Viable Validation

`TestResolvePreservesSelectedRuleAndAction` parses and normalizes a production
ordinary rule, resolves it, and asserts its selected-rule tuple and action match
the exact row. The fixture carries a non-empty expansion suffix and write. It
then mutates caller-owned row slices, write backing arrays, next-tag maps,
identity strings, and locator-constructor inputs before and after `Resolve` and
asserts the validated snapshot and returned plan remain unchanged. An
external-package API test proves callers cannot construct or replace row
components independently. A companion production TOML → normalize → resolve →
CLI test asserts that an authored escape returns its failure kind and complete
selected-rule identity with no next tags or writes; it must not hand-construct
`Edge` values.

### Phase 1: Normalize Selection Identity

Define the resolver-facing `SelectedRule` carrier once and have the validated
normalized-table constructor populate that same value, plus the distinct
action, on every ordinary row. Escape rows carry the same selected-rule value
and modeled refusal classes but no action.
The production resolver input uses the opaque, snapshot-owned table rather than
caller-assembled raw edges. Do not introduce a parallel normalizer identity
struct or a selection-time conversion.

### Phase 2: Preserve the Disposition Handoff

Extend the ordinary plan-construction path to return defensive copies of
selected rule and action without changing matching semantics. Add a modeled-
refusal constructor that returns the selected escape's underlying failure kind
and defensive selected-rule copy without inventing action data in `plan`.

### Phase 3: Hand Off the Consumer Contract

Expose the resolved disposition shape for RDR 0005 without adding CLI rendering
in this implementation. RDR 0005's implementation owns the downstream payload
fixture proving that `flow resolve` projects the selected rule and any
successful action without a second identity lookup or accessor execution, and
that it maps the modeled escape to its stable refusal with selected identity
and no action. It also maps user-authored validation failure to a stable
load/config `CLIError`. A3/A6 must align RDR 0005's normative contract and
text/JSON acceptance before this RDR locks. RDR 0008 implementation acceptance
ends at the carrier and resolver/replay tests; the release-level operator
outcome additionally requires the RDR 0005 consumer acceptance to land.

## Validation

### Testing Strategy

The source audit grounds ordinary and modeled-escape selection in
`internal/resolver/resolver.go::Resolve` and the current shared construction
boundary in `internal/resolver/resolver.go::plan`; the proposal splits that
boundary into ordinary plan, modeled-refusal, and kernel-refusal constructors.
The existing A2 spike under `evidence/spikes/` grounds only
the string-identity/action copy mechanism. Implementation must turn the
resolved A2/A6 evidence into this matrix:

1. **Validated row provenance** — try to construct ordinary and escape rows
   that pair rule B identity/locator with rule A predicates/action, including
   through the public production resolver API.
   **Expected**: the named normalized-table constructor rejects each malformed
   row, or the API makes the composition unrepresentable; no successful plan is
   returned.
2. **Ordinary selection** — resolve an exact-one ordinary edge through
   `Resolve` and `plan`.
   **Expected**: `TransitionPlan.SelectedRule` equals that edge's complete
   identity, including a non-empty expansion suffix, and `Action` equals its
   next tags and writes.
3. **Modeled escape end to end** — author a valid RDR 0002 TOML escape with no
   write or clear block, then run the production parser, normalizer, `Resolve`,
   and CLI projection.
   **Expected**: the result is the modeled `no_match` or `ambiguous_match`
   refusal with the escape row's complete `SelectedRule`, no plan, and no
   action. The test never hand-constructs an `Edge`.
4. **Copy isolation** — mutate caller-owned row slices, identity/locator inputs,
   next-tag maps, and write backing arrays after constructing the table but
   before `Resolve`, then again after resolution; concurrently resolve one
   stable snapshot under `go test -race`.
   **Expected**: the selected rule and action remain coherent and unchanged;
   inspection APIs expose no mutable backing storage.
5. **Incomplete identity** — load authored rows with an empty model id or rule
   id, zero/invalid source locator, or locator/identity mismatch; separately
   inject an impossible invalid value inside a validated table in a package test.
   **Expected**: authored invalidity fails before resolution as a stable
   load/config error; the impossible internal value uses the programmer error
   path. Neither returns a resolver disposition.
6. **Refusal branches** — exercise no match, ambiguous match, unavailable owned
   state, unevaluable guard, and unmodeled outcome without a selected escape.
   **Expected**: each disposition has no plan and no selected rule; only the
   modeled-escape fixture may attach selected identity to a refusal.
7. **Locator agreement (normalized-table acceptance)** — after A4 pins the
   constructor or validator boundary, supply a locator that designates a
   different model id or rule id from `SelectedRule`.
   **Expected**: the normalized-table boundary rejects the row before
   `Resolve`; `plan` does not parse or repair the locator.
8. **Replay identity** — extend RDR 0001's
   `TestResolve_ReplayReturnsValueIdenticalPlans` with ordinary and escape
   selected identity.
   **Expected**: two resolutions of the same input return value-identical
   dispositions: complete `SelectedRule` and `Action` for ordinary success,
   complete `SelectedRule` and `Kind` for modeled refusal, without a second
   identity lookup.
9. **CLI projection (RDR 0005 acceptance)** — for an ordinary plan and the A6
   selected-escape case, run production `flow resolve` in JSON and text modes.
   **Expected**: both modes expose the exact selected model id, rule id,
   expansion suffix, source locator, and any successful next tags/writes from
   the resolver disposition; the test observes no model lookup or accessor
   execution. The escape case maps to the stable CLI refusal for its kind,
   includes selected-rule identity, and omits next tags/writes. Missing ordinary
   action, missing selected identity, or any action on the modeled refusal fails
   consumer acceptance. This is a release-level consumer obligation, not code
   owned by this RDR.

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
- Stateless `src/Stateless/Transition.cs::Transition`.
- uscxml `src/uscxml/interpreter/LargeMicroStep.h::Transition`,
  `LargeMicroStep.cpp` take-transitions block, and
  `src/bindings/swig/wrapped/WrappedInterpreterMonitor.cpp::afterTakingTransition`.
- W3C SCXML 1.0 `3.1.5 Type and Transitions` — targetless transition semantics.
