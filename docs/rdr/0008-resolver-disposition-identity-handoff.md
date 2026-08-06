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
- **Profile**: foundational — cross-RDR identity producer spanning resolver
  dispositions and the CLI resolution payload.
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
model ID, rule ID, expansion suffix, and source locator while action shape
remains separately usable for ordinary and modeled escape selections.

### Technical Environment

The affected seam is the Go resolver in `internal/resolver/resolver.go`, whose
selection result feeds replay validation and the planned CLI resolution payload.
RDR 0001 defines the resolution kernel requirement to expose downstream values;
RDR 0002 defines normalized model/rule identity and source locators; RDR 0005
defines the consuming `flow resolve` payload contract. Draft RDR 0007 separately
owns normalized-table revision binding; this RDR does not define or validate the
table revision.

## Research Findings

### Investigation

The freshness check confirmed that RDRs 0001, 0002, and 0005 remain Final and
that the named resolver seam still exists at `internal/resolver/resolver.go`.
`Edge` carries only match/action data, `plan` copies only `NextTags` and
`Writes`, and `TransitionPlan` therefore cannot satisfy RDR 0005's selected-rule
payload. RDR 0007 now owns the adjacent table-revision binding question and is
explicitly outside this RDR.

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
- **Verified** — current `Edge`, `TransitionPlan`, and `plan` discard all four
  identity fields, for both ordinary and modeled-escape selections.
- **Assumed** — the future normalizer can populate a complete selected-rule
  value on every normalized ordinary and escape edge without a second lookup.

### Critical Assumptions

- **A1 Every normalized ordinary and escape edge can carry the complete RDR
  0002 selection identity without a downstream model lookup.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: Verify RDR 0002 `Normative Contracts` and the normalizer's
    eventual edge-construction symbol expose model id, rule id, deterministic
    expansion suffix, and source locator together.
  - **If wrong**: The resolver cannot make identity and action one atomic
    result, so CLI diagnostics may report a different rule than the one used.
- **A2 Copying a successful plan can preserve selection identity and action
  without aliasing mutable normalized-table storage.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: `TestResolvePreservesSelectedRuleForOrdinaryAndEscape` will
    mutate the input edge after resolution and require the returned plan's
    identity and action to remain unchanged.
  - **If wrong**: Audit output or replay assertions could drift after the
    caller reuses or mutates the supplied table.
- **A3 RDR 0005 can render the selected-rule value directly without
  reconstructing identity from flow/model inputs.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: Verify RDR 0005 `Technical Design` maps the resolver plan's
    matched-rule identity into the `flow resolve` payload and needs no richer
    resolver internals.
  - **If wrong**: This RDR must revise the carrier before lock or RDR 0005 will
    retain a parallel identity lookup.

## Proposed Solution

### Approach

Make `TransitionPlan`, the successful branch of `Disposition`, own two nested
values: an immutable selected-rule value and a separately copyable transition
action. The selected-rule value contains the RDR 0002 identity tuple—model id,
rule id, expansion suffix, and source locator. The action contains next tags and
owned-tag writes. Each normalized `Edge` supplies both values; `plan` copies
them into the returned plan.

Ordinary and modeled-escape matches use the same successful plan carrier. A
write-free escape therefore still reports the exact selected rule while its
action legitimately contains no writes. Kernel refusals carry no selected-rule
value because no row was selected. Downstream callers consume selection for
diagnostics/audit and action for state application without receiving guard or
candidate-table internals.

### Technical Design

The normalized-table boundary constructs each `Edge` from three conceptual
parts: match inputs, selected-rule identity, and transition action. Resolver
matching reads only the match inputs. On exact-one ordinary selection, or
exact-one modeled escape selection, `plan` returns a `TransitionPlan` containing
a value-copy of that edge's selected-rule identity and a defensive copy of its
action. RDR 0005 reads the former into `flow resolve` diagnostics and hands the
latter to the accessor/state boundary.

The selection value is descriptive, not a second authority: it is copied from
the normalized row that was actually matched. It does not recalculate table
revision (RDR 0007), redefine row identity (RDR 0002), change matching/refusal
semantics (RDR 0001), or define CLI rendering (RDR 0005).

#### Normative Contracts

```normative
Every successful resolver disposition MUST expose exactly one transition plan
whose selected-rule value contains model id, rule id, expansion suffix, and
source locator copied from the normalized edge that matched.

The transition plan MUST expose its transition action separately from the
selected-rule value. The action contains next tags and owned-tag writes; callers
MUST NOT need guard predicates or the matched Edge to apply it.

Ordinary and modeled-escape selections MUST use the same plan shape. A
write-free modeled escape MUST retain selected-rule identity even when its
action has no owned-tag writes. A kernel refusal MUST NOT claim selected-rule
identity because no normalized row was selected.

The returned selected-rule value and action MUST remain unchanged if the caller
later mutates or reuses the supplied normalized table.
```

#### Load-Bearing Decisions

- **Identity** — selected-rule equality is the RDR 0002 tuple `(model id, rule
  id, expansion suffix)`; source locator is carried diagnostic provenance and
  does not create a different logical rule identity.
- **Naming** — the successful carrier remains `TransitionPlan`; its nested
  values are "selected rule" and "transition action." Rejected: "provenance"
  alone because the value is operational selection identity, and "edge"
  because callers must not receive matching internals.
- **Selection / predicate** — identity is copied only from the exact edge RDR
  0001 already selects; this RDR adds no tie-break or second lookup.

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Exact-one resolver disposition | RDR 0001 | Available | Preserves existing success/refusal selection semantics. |
| Normalized rule identity and source locator | RDR 0002 | Predecessor | Supplies the selected-rule value; this RDR does not redefine it. |
| Table-revision binding | RDR 0007 | Deferred peer | Remains outside this identity-handoff contract. |
| Matched-rule CLI payload | RDR 0005 | Predecessor | Consumes selected-rule identity directly. |
| Successful selected-rule/action carrier | This RDR | Introduced | Extends the resolver plan across the producer/consumer seam. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Normalized candidate | `internal/resolver/resolver.go::Edge` | Carries match and action only | Extend | Add the RDR 0002-selected rule value without changing matching. |
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

The chosen plan-owned carrier wins on the deciding correctness and prior-art
criteria while keeping blast radius bounded. Putting selection on `Disposition`
would preserve the current action-only plan, but creates an optional-field
invariant across success and refusal branches. Returning `Edge` leaks guard and
candidate storage. CLI reconstruction violates the user's audit outcome because
reported identity could diverge from the edge actually selected.

Premortem: this ships and diagnostics sometimes show an empty or stale rule id
because one success path constructs the action without copying selection, or a
caller mutates the source edge after resolution. The recommendation survives:
one successful `TransitionPlan` constructor path can make selected rule and
action inseparable, and A2's MVV requires defensive-copy behavior for both
ordinary and write-free modeled-escape paths.

## Alternatives Considered

### Alternative 1: Selection on Disposition

**Description**: Add an optional selected-rule value beside `Plan` and
`Refusal`, leaving `TransitionPlan` as action-only.

**Pros**:

- Keeps action-only plans unchanged.
- Makes selection visible at the outermost result boundary.

**Cons**:

- Allows invalid states: plan without selection, refusal with selection, or
  selection without either branch.
- Requires every disposition constructor and test to enforce a cross-field
  invariant.

**Reason for rejection**: The successful plan is the narrower carrier that can
make selection and action atomic without weakening refusal semantics.

### Alternative 2: Return the Matched Edge

**Description**: Put the selected normalized `Edge` directly in the successful
disposition and let callers read both identity and action from it.

**Pros**:

- Preserves every source field with almost no projection code.
- Closely resembles uscxml retaining its selected transition object.

**Cons**:

- Exposes guards, outcome matching, and normalized-table storage to action
  consumers.
- Risks slice/map aliasing and turns future `Edge` changes into downstream API
  changes.

**Reason for rejection**: The resolver should return the minimum successful
projection, not its candidate input object.

### Briefly Rejected

- **CLI-side identity reconstruction**: Rejected because a second model lookup
  can report a row other than the one the kernel selected.

## Trade-offs

### Consequences

- Successful resolution becomes self-describing for diagnostics and replay
  assertions, including modeled escapes with empty write sets.
- Action consumers can apply next tags/writes without depending on rule identity
  or guard internals.
- `Edge`, `TransitionPlan`, construction code, and fixtures must all grow in
  lockstep; zero-valued identity becomes a new defect class to refuse or test.

### Risks and Mitigations

- **Risk**: One success path omits or partially fills selected-rule identity.
  **Mitigation**: Use one plan-construction path and make the MVV cover ordinary
  and modeled-escape selections.
- **Risk**: Returned maps or slices alias caller-owned edge data.
  **Mitigation**: Clone the action and value-copy immutable identity; prove it by
  mutating the input fixture after resolution.
- **Risk**: This RDR accidentally absorbs table revision or CLI wire-format
  ownership.
  **Mitigation**: Keep revision validation in RDR 0007 and rendering in RDR 0005.

### Failure Modes

- **Missing identity on success**: CLI output lacks model/rule/source fields;
  the MVV fails before release. Diagnose at `resolver.plan` construction.
- **Identity/action mismatch**: audit names one rule while next tags/writes came
  from another. Treat as an internal invariant failure; do not reconstruct or
  guess downstream.
- **Aliased result**: replay output changes after table reuse. The mutation
  fixture identifies which selected-rule/action field was not copied.
- **Refusal claims a rule**: a no-match or ambiguous refusal falsely implies a
  selection. Reject the invalid disposition in tests and keep refusal identity
  absent.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified
- [ ] RDR 0002 selected-rule fields are available at normalized edge creation.
- [ ] RDR 0007 boundary is confirmed so revision binding is not duplicated.

### Minimum Viable Validation

`TestResolvePreservesSelectedRuleForOrdinaryAndEscape` resolves one ordinary
edge and one write-free modeled escape, asserts the returned selected-rule tuple
matches the exact input edge and the action matches that edge, then mutates the
input table and asserts both returned plans remain value-identical.

### Phase 1: Normalize Selection Identity

Carry the RDR 0002 identity tuple and action as distinct values on every
ordinary and modeled-escape edge.

### Phase 2: Preserve the Successful Handoff

Extend the single plan-construction path to return defensive copies of selected
rule and action without changing matching or refusal semantics.

### Phase 3: Prove the Consumer Boundary

Add ordinary, modeled-escape, mutation, and downstream payload fixtures that
show diagnostics consume the selected rule while state application consumes the
separate action.

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
