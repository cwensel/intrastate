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
escape path; they discover the gap when a successful resolution reports an
anonymous plan that cannot be traced to its rule or source. Internally, the
resolver must preserve stable normalized rule identity and source location
through its disposition boundary so downstream CLI output and replay tests can
consume that identity unchanged.

## Context

### Background

Review of the resolution-kernel, transition-table, and CLI-integration RDRs
found that `Edge`, `TransitionPlan`, and `plan` discard normalized rule identity
and source location. The gap is not exercised by a production caller yet, but
it becomes downstream-reachable when the RDR 0005 `flow resolve` payload is
wired. The design fork is which immutable resolver disposition carrier owns
model ID, rule ID, expansion suffix, and source locator, and whether modeled
escape selection carries a transition action or an identity-bearing non-success
disposition.

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

Prior art was read before alternatives were named. RDR 0002 `Normative
Contracts` requires: "Each candidate row MUST retain its source rule id and
source locator." Its `Load-Bearing Decisions / Identity` defines `(model id,
rule id)` plus a deterministic expansion suffix. RDR 0005 `Technical Design`
requires `flow resolve` data to contain "matched rule identity." In external
source, Stateless `src/Stateless/Transition.cs::Transition` describes one
successful transition with get-only `Source`, `Destination`, and `Trigger`
values. uscxml `src/uscxml/interpreter/LargeMicroStep.h::Transition` retains the
source transition element, source, targets, event, and condition; the matching
`LargeMicroStep.cpp` take-transitions block passes that same selected transition
element to before/after monitor callbacks. These sources favor carrying
selection identity on the successful transition value rather than
reconstructing it downstream. Queries and rejected branches are recorded under
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
- **Verified on the integration branch** — `Edge`, `TransitionPlan`, and `plan`
  discard all four identity fields, for both ordinary and modeled-escape
  selections. The symbols do not yet resolve on `main`; landing the RDR 0001
  implementation is an explicit implementation prerequisite below.
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
  - **If wrong**: The resolver cannot make identity and action one atomic
    result, so CLI diagnostics may report a different rule than the one used.
- **A2 Copying a successful plan can preserve selection identity, typed source
  locator, and action without aliasing mutable normalized-table storage.**
  - **Status**: Pending
  - **Method**: Spike + Peer RDR
  - **Evidence**: The existing
    `TestPlanOwnsSelectedIdentityAndAction` spike proves isolation for string
    identity fields, a string locator, a next-tag map, and a writes slice. It
    does not prove the revised RDR 0002-owned typed locator is deeply immutable
    or cloned, nor that the validated table owns caller storage before
    `Resolve`. Stage 6 must rerun the spike against the concrete A4/A5 types,
    mutating constructor inputs before resolution and source storage after
    resolution; an external-package API test must show row components cannot be
    replaced after validation.
  - **If wrong**: Audit output or replay assertions could drift after the
    caller reuses or mutates the supplied table.
- **A3 RDR 0005 can require and test direct projection of the selected-rule
  value without reconstructing identity from flow/model inputs.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: RDR 0005 `Technical Design` includes matched rule identity in
    the `flow resolve` payload, but its current `Normative Contracts` require
    only a resolved next tag-set or refusal. Stage 6 must verify that RDR 0005 is
    revised to require `SelectedRule` and `Action` projection in JSON and text
    acceptance tests, without a candidate-table lookup or accessor execution,
    before this RDR locks.
  - **If wrong**: This RDR must revise the carrier before lock or RDR 0005 will
    omit identity or retain a parallel identity lookup.
- **A4 The normalized-table boundary can construct a typed source locator that
  remains bound to the selected model and rule.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: Stage 6 must verify against RDR 0002 `Technical Design` and
    `Normative Contracts` that the locator retained on each normalized row is a
    typed value constructed with that row's model id and rule id, rather than
    an opaque string checked only for presence. Before lock, that verification
    must name the normalized-table constructor or validator, its zero-value and
    mismatch behavior, and the mismatch test at that callable boundary; `plan`
    remains a copying consumer of the validated value.
  - **If wrong**: A successful plan could report the selected identity while
    directing diagnostics to a different authored rule.
- **A5 The production resolver entry point can accept only normalized-table
  rows that bind match inputs, selected identity, locator, and action at one
  validated construction boundary.**
  - **Status**: Pending
  - **Method**: Source Search + Peer RDR
  - **Evidence**: Current `internal/resolver/resolver.go::{Input,Edge}` exposes a
    caller-constructible `[]Edge`, while RDR 0007 requires outcomes and rows to
    become one normalized semantic table value. Stage 6 must name the concrete
    constructor/validator and production input type that prevent callers from
    independently composing rule B identity with rule A predicates or action,
    and must name ordinary and escape mismatch tests at `Resolve`.
  - **If wrong**: The resolver can execute one row's action while returning
    authoritative-looking identity for another row.
- **A6 Modeled escape selection has one action/disposition meaning shared by
  the source schema, normalized row, resolver, and CLI.**
  - **Status**: Pending
  - **Method**: Peer RDR + Source Search
  - **Evidence**: RDR 0002 currently forbids `write` and `clear` on an escape
    rule, while RDR 0001 and `internal/resolver/resolver.go::Resolve` treat an
    exactly-one escape edge as a successful plan; current resolver tests
    hand-construct escape edges with next tags and writes. Stage 6 must reconcile
    those contracts at the production parse → normalize → resolve boundary:
    either define how a selected escape obtains a complete action, or define it
    as an identity-bearing non-success disposition with no action.
  - **If wrong**: `flow resolve` can report a successful selected escape with
    an empty or invented next tag-set, causing automation to repeat the same
    state or violate the authored schema.

## Proposed Solution

### Approach

Make `TransitionPlan`, the successful branch of `Disposition`, own two nested
values: a `SelectedRule` and a separately copyable `TransitionAction`. RDR 0008
owns this Go carrier shape; RDR 0002 continues to own identity semantics and
normalized-table population. `SelectedRule` is the one value shared by a
validated normalized row and `TransitionPlan`, not a second resolver projection
or a value reconstructed during selection. It contains model id, rule id,
expansion suffix, and a typed source locator. `TransitionAction` contains next
tags and owned-tag writes. The normalized-table constructor owns caller storage
and binds match inputs, selected rule, locator, and action into one inaccessible
row; the production resolver entry point accepts one stable validated-table
snapshot and `plan` copies or clones both nested values from the exact row.

Ordinary matches use the successful plan carrier. A6 must settle modeled-escape
semantics before lock: RDR 0002 currently forbids action fields on authored
escape rules, while RDR 0001 and the resolver treat an exactly-one escape edge
as a successful plan. If escape remains a success, normalization must supply a
complete, explicitly defined action; otherwise the selected escape needs an
identity-bearing non-success disposition with no action. Kernel refusals for
which no row was selected carry no selected-rule value. Downstream callers
consume selection for diagnostics/audit and action, when present, for state
application without receiving guard or candidate-table internals.

### Technical Design

The normalized-table boundary constructs an opaque, snapshot-owned validated
table whose internal ordinary rows bind match inputs to
`SelectedRule SelectedRule` and `Action TransitionAction`. Production callers
cannot independently assemble,
replace, or mutate those row components after validation; inspection and dump
APIs return immutable values or defensive copies. `SelectedRule` has string fields
`ModelID`, `RuleID`, and `ExpansionSuffix`, plus a typed `SourceLocator` value
owned by the RDR 0002 normalization boundary. `TransitionAction` has
`NextTags TagSet` and `Writes []OwnedTagWrite`. `TransitionPlan` exposes the
same selected rule and action. Resolver matching reads only match inputs. On
exact-one ordinary selection, `plan` copies the matched row's selected rule and
action without retaining mutable aliases. A2/A4 must determine whether the
typed locator is deeply immutable and value-copyable or requires an explicit
clone. A6 determines the disposition and action-copy rule for an exactly-one
modeled escape.

Before returning success, `plan` treats an empty model id or rule id, or a zero
or invalid source locator, as a programmer/invariant error on the Go error
path, not as a modeled refusal. An expansion suffix may be empty for an
unexpanded source rule; its value still participates in identity. The source
locator is an RDR 0002 typed value that this handoff copies without parsing or
reformatting. Before a row reaches the resolver, the normalized-table boundary
constructs the locator with the same model id and rule id carried by
`SelectedRule`; `plan` enforces required presence but does not reconstruct or
reinterpret provenance. RDR 0005 must project the selected rule and any action
carried by the resolved disposition into `flow resolve` output without
reconstructing identity or executing accessors. User-authored model validation
failures are load/config failures that RDR 0005 maps to a stable `CLIError`; an
invariant breach in an already validated table remains an internal Go error. A
later state operation may pass a successful action to the RDR 0004 accessor
boundary.

The selection value is descriptive, not a second authority: it is copied from
the normalized row that was actually matched and identifies that logical row
within the revision-bound resolver input. It is not a standalone historical
event identity. It does not recalculate or enforce table revision (RDR 0007 and
its charted enforcement successor), redefine row identity (RDR 0002), change
matching/refusal semantics (RDR 0001), or define CLI rendering (RDR 0005).

The existing RDR 0001 replay test is the replay consumer for this handoff. This
RDR extends that test so two resolutions of the same input compare complete
`SelectedRule` and `Action` values; replay does not gain a second identity
lookup or a separate carrier.

#### Normative Contracts

```normative
The production resolver entry point MUST consume one validated normalized-table
value whose construction binds each row's match inputs, `SelectedRule`,
`SourceLocator`, and `Action` into a stable owned snapshot. Production callers
MUST NOT be able to compose, replace, or mutate those row components after
validation, and inspection APIs MUST NOT expose mutable backing storage.
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

Before lock, the RDR 0001 resolver disposition and RDR 0002 source-schema
contracts MUST define one modeled-escape meaning. If an exactly-one modeled
escape is a successful transition, normalization MUST supply its complete
action and define `NextTags` rather than inventing it in `plan`. If it is a
selected non-success disposition, it MUST retain selected-rule identity and
MUST NOT claim an action. A kernel refusal for which no row was selected MUST
NOT claim selected-rule identity.

The returned `SelectedRule` and `Action` MUST remain unchanged if the caller
mutates or reuses constructor inputs before `Resolve`, later reuses source
storage, or inspects the validated table through a public API.

`flow resolve` MUST project the resolver disposition's `SelectedRule` and any
successful `Action` in both JSON and text terminal output without a
candidate-table lookup or accessor execution. RDR 0005 owns rendering, stable
load/config error codes, and CLI wiring, but its consumer contract and
acceptance tests MUST preserve this handoff before this RDR locks.
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
  `TransitionPlan.SelectedRule` and `TransitionPlan.Action`. Rejected:
  "provenance" alone because the value is operational selection identity, and
  "edge" because callers must not receive matching internals.
- **Selection / predicate** — identity is copied only from the exact edge RDR
  0001 already selects; this RDR adds no tie-break or second lookup.

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Exact-one resolver disposition and escape meaning | RDR 0001 / RDR 0002 | A6 Pending | Must reconcile whether a selected escape is a success with a complete action or an identity-bearing non-success. |
| Validated normalized rows, rule identity, and typed source locator | RDR 0002 | Predecessor; A2/A4/A5 Pending | Supplies one snapshot-owned row value; this RDR does not redefine identity semantics. |
| Table-revision derivation and binding | RDR 0007 + charted successor | Deferred peers | Remain outside this logical identity-handoff contract. |
| Matched-rule CLI payload | RDR 0005 | Predecessor; A3/A6 Pending | Must normatively project selected-rule identity and any successful action directly. |
| Successful selected-rule/action carrier | This RDR | Introduced | Extends the resolver plan across the producer/consumer seam. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Normalized candidate | `internal/resolver/resolver.go::{Input,Edge}` | Exported raw rows and backing slices let callers compose or mutate unrelated identity and action | Replace boundary | Accept an opaque, snapshot-owned normalized-table value; return copies from inspection APIs. |
| Successful result | `internal/resolver/resolver.go::TransitionPlan` | Carries action only | Extend | Own selected rule and separately usable action. |
| Result construction | `internal/resolver/resolver.go::plan` | Drops edge identity and copies action fields only | Extend | Copy both values defensively from the matched edge. |
| CLI output | `internal/cli/respond` and RDR 0005 | No resolver verb wired yet | Reuse later | Rendering stays outside the kernel. |

### Decision Rationale

The Questions-Options-Criteria matrix scores 1 (poor) through 5 (strong):

| Approach | Correctness fit | Prior-art alignment | Reversibility | Blast radius | Cost | Total |
| --- | --- | --- | --- | --- | --- | ---: |
| Plan owns selected rule + separate action | 5 — success makes both atomic | 5 — mirrors successful transition values | 4 — nested values can evolve | 4 — resolver types and tests only | 4 — one propagation path | **22** |
| Disposition owns optional selected rule | 3 — permits plan/identity mismatch | 4 — transition context stays near result | 5 — plan stays unchanged | 4 — one extra disposition field | 3 — every branch must enforce presence | **19** |
| Return the matched Edge | 2 — exposes match internals and aliases | 3 — preserves transition object | 2 — hard to narrow later | 2 — downstream couples to resolver input | 4 — minimal initial code | **13** |
| Reconstruct identity in CLI | 1 — second lookup can disagree | 1 — breaks selected-object handoff | 3 — isolated but removable | 1 — duplicates model selection | 2 — lookup/index plumbing | **8** |

The plan-owned carrier wins for ordinary success on the deciding correctness
and prior-art criteria while keeping blast radius bounded. A6 must revisit only
the modeled-escape row: if RDR 0002 continues to define an escape without an
action, an identity-bearing non-success disposition is no longer merely the
rejected optional-field alternative. Returning `Edge` leaks guard and candidate
storage. CLI reconstruction violates the user's audit outcome because reported
identity could diverge from the edge actually selected.

Premortem: this ships and diagnostics sometimes show an empty or stale rule id
because one success path constructs the action without copying selection, or a
caller mutates the source edge after resolution. The recommendation survives:
one ordinary `TransitionPlan` constructor path can make selected rule and action
inseparable, and A2's MVV requires snapshot ownership and defensive-copy
behavior. A6 must define the escape path before the recommendation can lock.

## Alternatives Considered

### Alternative 1: Selection on Disposition

Add an optional selected-rule value beside `Plan` and `Refusal`, leaving
`TransitionPlan` action-only. As a general success shape this permits plan
without selection, refusal with selection, and selection without either branch;
the ordinary successful plan is the narrower atomic carrier. A6 may still choose
a distinct identity-bearing escape disposition if the source schema continues
to define selected escapes without actions.

### Alternative 2: Return the Matched Edge

Put the selected normalized `Edge` directly in the successful disposition.
This preserves every source field with little projection code, but exposes
guards and matching internals, risks slice/map aliasing, and couples downstream
APIs to candidate storage. The resolver should return the minimum successful
projection instead.

### Briefly Rejected

- **CLI-side identity reconstruction**: Rejected because a second model lookup
  can report a row other than the one the kernel selected.

## Trade-offs

### Consequences

- Ordinary successful resolution becomes self-describing for diagnostics and
  replay assertions; A6 determines the corresponding modeled-escape shape.
- Action consumers can apply next tags/writes without depending on rule identity
  or guard internals.
- `Edge`, `TransitionPlan`, construction code, and fixtures must all grow in
  lockstep; incomplete selected-rule identity becomes a programmer error to
  reject and test.

### Risks and Mitigations

- **Risk**: One success path omits or partially fills selected-rule identity.
  **Mitigation**: Use one ordinary plan-construction path and make the MVV cover
  every disposition shape A6 retains.
- **Risk**: A caller pairs one row's identity or locator with another row's
  predicates or action.
  **Mitigation**: Make `Resolve` consume only the validated normalized-table
  value and test that malformed ordinary and escape rows cannot reach success.
- **Risk**: Returned maps or slices alias caller-owned edge data.
  **Mitigation**: Make validated-table construction own caller storage; clone
  the action and clone or value-copy only deeply immutable identity fields.
  Mutate source inputs before and after resolution and run concurrent snapshot
  use under the race detector.
- **Risk**: A selected modeled escape produces a success with no defined state
  transition.
  **Mitigation**: Reconcile RDR 0001 and RDR 0002 at A6; test authored TOML
  through production normalization, resolution, and CLI projection.
- **Risk**: This RDR accidentally absorbs table revision or CLI wire-format
  ownership.
  **Mitigation**: Keep revision validation in RDR 0007 and rendering in RDR 0005.

### Failure Modes

- **Incomplete identity on a selected edge**: `plan` returns a programmer error
  before producing a successful disposition; the MVV covers every required
  field. Diagnose at `resolver.plan` construction.
- **Identity/provenance mismatch on a normalized edge**: normalization rejects
  the row before resolver selection. Diagnose at the normalized-table
  validation boundary; `plan` must not guess or repair provenance.
- **Identity/action mismatch**: audit names one rule while next tags/writes came
  from another. Reject it at the validated normalized-table construction
  boundary before `Resolve`; do not reconstruct or guess downstream.
- **Aliased result**: replay output changes after table reuse. The mutation
  fixture identifies which selected-rule/action field was not copied.
- **Undefined escape action**: a selected escape returns success with empty or
  invented next tags. Reject the unresolved contract at A6; the production
  parser-to-CLI test fixes one disposition and action meaning.
- **Refusal claims a rule**: a no-match or ambiguous refusal falsely implies a
  selection. Reject the invalid disposition in tests and keep refusal identity
  absent.

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
  proves callers cannot independently compose or mutate row identity and action.
- [ ] A6 reconciles modeled-escape disposition and action semantics across RDR
  0001, RDR 0002, production normalization, resolver tests, and CLI projection.
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
CLI test locks the selected-escape disposition and action meaning chosen by A6;
it must not hand-construct `Edge` values.

### Phase 1: Normalize Selection Identity

Define the resolver-facing `SelectedRule` carrier once and have the validated
normalized-table constructor populate that same value, plus the distinct
action, on every ordinary row. A6 defines the normalized selected-escape value.
The production resolver input uses the opaque, snapshot-owned table rather than
caller-assembled raw edges. Do not introduce a parallel normalizer identity
struct or a selection-time conversion.

### Phase 2: Preserve the Successful Handoff

Extend the ordinary plan-construction path to return defensive copies of
selected rule and action without changing matching semantics. Implement the
separately specified A6 escape disposition without inventing action data in
`plan`.

### Phase 3: Hand Off the Consumer Contract

Expose the resolved disposition shape for RDR 0005 without adding CLI rendering
in this implementation. RDR 0005's implementation owns the downstream payload
fixture proving `flow resolve` projects the selected rule and any successful
action without a second identity lookup or accessor execution, and maps
user-authored validation failure to a stable load/config `CLIError`. A3/A6 must
align RDR 0005's normative contract and text/JSON acceptance before this RDR
locks. RDR 0008 implementation acceptance ends at the carrier and
resolver/replay tests; the release-level operator outcome additionally requires
the RDR 0005 consumer acceptance to land.

## Validation

### Testing Strategy

The source audit grounds ordinary and modeled-escape selection in
`internal/resolver/resolver.go::Resolve` and the current shared construction
boundary in `internal/resolver/resolver.go::plan`; A6 decides whether that
sharing survives. The existing A2 spike under `evidence/spikes/` grounds only
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
   **Expected**: the result has the one selected-rule/disposition meaning A6
   locks. A success has a source-defined complete action; a selected non-success
   has no action. The test never hand-constructs an `Edge`.
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
   **Expected**: each disposition has no plan and therefore claims no selected
   rule.
7. **Locator agreement (normalized-table acceptance)** — after A4 pins the
   constructor or validator boundary, supply a locator that designates a
   different model id or rule id from `SelectedRule`.
   **Expected**: the normalized-table boundary rejects the row before
   `Resolve`; `plan` does not parse or repair the locator.
8. **Replay identity** — extend RDR 0001's
   `TestResolve_ReplayReturnsValueIdenticalPlans` with selected identity.
   **Expected**: two resolutions of the same input return value-identical
   `SelectedRule` and `Action`, including a non-empty expansion suffix, without
   a second identity lookup.
9. **CLI projection (RDR 0005 acceptance)** — for an ordinary plan and the A6
   selected-escape case, run production `flow resolve` in JSON and text modes.
   **Expected**: both modes expose the exact selected model id, rule id,
   expansion suffix, source locator, and any successful next tags/writes from
   the resolver disposition; the test observes no model lookup or accessor
   execution. Absence of any required selected-rule field or action, or an
   unexplained action on a selected non-success, fails consumer acceptance. This
   is a release-level consumer obligation, not code owned by this RDR.

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
