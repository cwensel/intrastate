# Recommendation 0004: Accessor Execution Safety Model

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-06-19
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
- **Profile**: large — locks write-accessor safety semantics for authoritative artifacts.
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
- **Priority**: High
- **Related Issues**: None
- **Predecessors**: 0001-resolution-kernel, 0002-transition-table-as-reviewable-data, 0003-guard-predicate-exhaustiveness
- **Overrides**: None
- **Seam Lineage**: no prior accretion

## Problem Statement

A skill needs read, write, and gate accessors to run as reliable passthroughs: failures must surface, writes must not silently corrupt authoritative artifacts, and external calls must not overstep their intended power. The system-internal requirement is to define accessor execution, failure, timeout, gate-indeterminate, and read-back-verify semantics.

## Context

### Background

The resolver is stateless and non-orchestrating, so accessors are the injected I/O seam for owned tags and world facts. Write accessors mutate authoritative artifacts with no ledger or undo, making the safety model part of the core contract rather than an implementation detail.

The real design fork is power versus blast radius: raw shell-out, an allowlisted command set, or a declared-capability model such as read-only, persist-decision, and forbidden.

### Technical Environment

intrastate is a Go CLI wired through `internal/cli`. All command output must route through `internal/cli/respond`, and user-facing failures must use structured `CLIError` values through `internal/cli/clierr`.

## Research Findings

### Investigation

The proposal is shaped by the peer RDR split, the CLI output contract, and a
bounded prior-art pass over state-machine and infrastructure state-write
systems. RDR 0001 requires a stateless, non-orchestrating resolver that reads
and persists through accessors; RDR 0002 names accessor references in the sparse
transition model but delegates execution; RDR 0003 keeps guards as symbolic
predicates over tag values after accessor binding. `docs/cli-output-contract.md`
requires every graceful failure to route through `respond` / `clierr`, so
accessor failures must become typed refusals rather than direct prints.

Prior art was read before choosing an approach. Stateless stores entry, exit,
activate, and deactivate actions as callback collections and executes them from
state transitions (`StateRepresentation::ExecuteEntryActions`,
`StateRepresentation.Async::ExecuteEntryActionsAsync`); it is useful evidence
that action hooks are common, but also shows why unrestricted callbacks are too
powerful for intrastate's reviewable-data model. OpenTofu's taint command writes
and then persists state with explicit error diagnostics
(`TaintCommand::Run`, `stateMgr.WriteState`, `stateMgr.PersistState`), which
supports treating state mutation as an explicit checked operation rather than a
hidden transition side effect. The local ADO transition helper first tries a
direct external update, then walks known intermediate states only for a specific
400 response (`Client::transitionWorkItem`); that supports typed external
failure classification instead of swallowing remote errors.

Sibling-path check:

```sh
rg -n "accessor|starlark|eval|exec|script|sandbox|determin|Config|toml|intrastate.toml|policy" .
```

The search found no implemented accessor executor or capability discriminator
under `internal/`. Existing code only provides the CLI config loader and the
`respond` / `clierr` output/failure gateway, so this RDR introduces the
accessor safety contract and reuses existing CLI failure plumbing later.

### Key Discoveries

- **Documented** — RDR 0001 depends on RDR 0004 for artifact-bound accessor
  execution and read-back semantics compatible with a stateless resolver.
- **Documented** — RDR 0002 may reference accessors in the transition model but
  must not execute them; this RDR owns the execution boundary.
- **Documented** — RDR 0003 consumes tag values after accessor binding and does
  not execute accessors during guard evaluation.
- **Documented** — the CLI output contract requires graceful failures to use the
  existing structured envelope and exit-code mapping.
- **Documented** — callback-based FSM prior art executes action hooks directly
  during transitions; useful as a contrast, but not safe enough for a
  reviewable static model.
- **Assumed** — a small declared-capability vocabulary can express the RDR and
  kata accessors without falling back to raw shell strings or host callbacks.
- **Assumed** — read-back verification can prove the expected owned-tag effect
  for initial write accessors without needing a full undo log.

### Critical Assumptions

- **A1 The target RDR and kata flows only need declared read, gate, and write
  accessors over caller-supplied artifact roles.**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: Pending: Resolve must bind representative RDR and kata flow
    fixtures to read, gate, and write accessors without raw shell strings or
    host-language transition callbacks.
  - **If wrong**: The capability vocabulary is too small, and authors will
    pressure the model toward unsafe command execution.
- **A2 Write accessors can verify their intended owned-tag effect by re-reading
  the same artifact boundary after the write.**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: Pending: Resolve must run a write-then-read-back fixture that
    changes an owned tag, re-reads it through the accessor seam, and captures a
    mismatch as a typed failure.
  - **If wrong**: A successful write command could silently corrupt or fail to
    update authoritative state.
- **A3 Timeout, execution failure, gate indeterminate, and read-back mismatch
  can be represented as stable accessor refusal classes and mapped through the
  existing CLI failure gateway.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: Pending: Resolve must confirm `internal/cli/clierr::CLIError`,
    `internal/cli/clierr::ExitCodeFor`, and `internal/cli/respond::Fail` can
    carry accessor refusal classes without adding direct output from the
    accessor package.
  - **If wrong**: Accessor errors would need a separate user-facing output
    contract or would leak implementation errors to callers.
- **A4 Accessor execution can be deterministic enough for resolver replay when
  the model records artifact role, accessor name, capability, timeout, and
  returned tag values.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: Pending: Resolve must name a replay test where the same
    transition model and fixture artifacts produce the same read/gate/write
    disposition or the same typed refusal.
  - **If wrong**: Resolver replay could depend on ambient process state rather
    than declared model inputs.
- **A5 External API accessors can be constrained by declared capability and
  timeout without requiring intrastate to own credentials or remote lifecycle.**
  - **Status**: Pending
  - **Method**: Design Decision
  - **Evidence**: Pending: this RDR scopes credentials and remote resource
    lifecycle outside intrastate; accessors receive caller-provided environment
    and return typed success/failure only.
  - **If wrong**: The accessor layer becomes an orchestrator and should be split
    into a separate integration contract.

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

Use a declared-capability accessor model. The transition model names accessors;
the accessor registry binds each name to one capability class: `read`, `gate`,
or `write`. A read accessor returns typed tag values from a caller-supplied
artifact role. A gate accessor returns allow, deny, or indeterminate with a
reason. A write accessor applies a planned owned-tag mutation to a
caller-supplied artifact role, then re-reads the same role and verifies that the
expected owned-tag values are present.

The executor is not a shell runner and not a state-machine action callback
surface. It invokes typed bindings selected by accessor name and capability,
enforces a per-accessor timeout, converts execution errors into stable refusal
classes, and never prints directly. The resolver stays stateless: it receives
the accessor-read values and write plan disposition, while RDR 0005 later maps
those results to the user-facing CLI.

### Technical Design

Accessor definitions live beside the transition model as data consumed by the
loader/validator. Each definition declares a stable name, capability, artifact
role, expected tag keys, timeout policy, and whether read-back verification is
required. RDR 0002 owns the table carrier and accessor references; this RDR owns
what it means to execute a referenced accessor safely.

Execution has three phases. First, validation rejects unknown accessor names,
capability mismatches, writes to non-owned tags, and missing timeout/read-back
metadata before resolution. Second, runtime invocation applies read and gate
accessors to caller-supplied artifacts and classifies timeout, unavailable
artifact, execution error, and gate-indeterminate outcomes. Third, write
accessors execute only from a successful transition plan, then immediately
re-read through the same accessor boundary and compare expected owned-tag
values. A read-back mismatch is a failure even if the write command exited
successfully.

Large-profile Q-O-C matrix:

| Approach | Correctness fit | Prior-art alignment | Reversibility | Blast radius | Cost |
| --- | --- | --- | --- | --- | --- |
| Raw shell-out accessors | Weak: any command can mutate state outside the model. | Diverges from RDR 0002/0003 static data; resembles opaque callbacks. | Poor: side effects are uncontrolled. | Highest: grants ambient process power. | Low upfront, high debugging cost. |
| Global allowlisted command set | Medium: constrains executable names but not semantic authority. | Partial fit with CLI tools, weak fit with typed FSM semantics. | Medium: still hard to prove the intended tag changed. | High: one allowlist applies across accessors. | Medium. |
| Declared read/gate/write capabilities | Strong: capability, artifact role, timeout, and read-back are all model data. | Aligns with peer RDR split and with state-write prior art that checks persistence errors. | Good: write success is checked by re-read; undo is not claimed. | Bounded per accessor and artifact role. | Medium; needs validator and fixture spikes. |
| No write accessors | Strong for read safety, fails persistence outcome. | Aligns with pure resolver but contradicts RDR 0001's accessor persistence dependency. | Best because there are no writes. | Low. | Low, but incomplete. |
| FSM action callbacks | Weak for static proof; action code decides behavior. | Common in Stateless (`ExecuteEntryActions`), but deliberately rejected by RDR 0003's symbolic guard model. | Poor unless every callback self-audits. | High: arbitrary host code runs during transitions. | Low initially, high review cost. |

#### Normative Contracts

```normative
Every accessor definition MUST declare exactly one capability: read, gate, or
write. Runtime execution MUST reject any attempt to use an accessor for a
different capability than the one declared.
```

```normative
Accessors MUST operate on caller-supplied artifact roles. The accessor executor
MUST NOT discover authoritative artifacts from ambient process state.
```

```normative
A read accessor MUST return typed tag values or a typed refusal. It MUST NOT
mutate authoritative artifacts.
```

```normative
A gate accessor MUST return allow, deny, or indeterminate. Indeterminate MUST be
a refusal-class result, not a false allow and not a false deny.
```

```normative
A write accessor MUST apply only planned owned-tag writes produced by a
successful transition. It MUST NOT write observed or recognized tags.
```

```normative
After a write accessor reports command-level success, the executor MUST re-read
the same artifact role and verify the expected owned-tag values. A read-back
mismatch MUST be reported as a write failure.
```

```normative
Every accessor invocation MUST have a bounded timeout. Timeout MUST be reported
as its own refusal class, distinct from execution failure and read-back
mismatch.
```

```normative
The accessor package MUST return structured success/refusal values and MUST NOT
write to stdout or stderr directly.
```

#### Load-Bearing Decisions

- **Identity** — an accessor is identified by `(flow id, accessor name,
  capability)`. The same name cannot be rebound to a different capability within
  one flow.
- **Wire / byte format** — RDR 0002 owns the TOML carrier. This RDR owns the
  accessor execution semantics embedded behind accessor references.
- **Naming** — the canonical names are "read accessor", "gate accessor", and
  "write accessor". Rejected names: "hook" and "action" because they imply
  arbitrary transition callbacks.
- **Selection / predicate** — when an accessor reference names a capability, the
  executor selects only a binding with the same accessor identity and capability.
  Missing or multiply-bound accessors are validation failures.

#### Round-Trip / Inverse Invariants

`write -> read = expected owned-tag value identity` for the written tag subset:
after a write accessor reports command-level success, reading through the same
artifact role must return the expected owned-tag values. This is not an undo or
byte-for-byte artifact invariant; it is the minimum safety invariant for the
authoritative tag values this RDR owns.

#### Illustrative Code

Illustrative execution shape only:

1. Validate that the transition model names `state.read` as read and
   `state.persist` as write.
2. Read owned tags from the caller's `state` artifact role.
3. Resolve the transition using RDR 0001 and RDR 0003.
4. Apply the planned owned-tag write through `state.persist`.
5. Re-read `state` and compare the expected owned-tag values.
6. Return success or a typed refusal without printing.

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Transition model accessor references | RDR 0002 | Pending | This RDR assumes the model can name accessors and tag provenance. |
| Guard evaluation after binding | RDR 0003 | Pending | Accessors provide tag values; guards consume them symbolically. |
| Stateless transition plan | RDR 0001 | Pending | Write accessors execute only a successful planned transition. |
| Accessor capability taxonomy | This RDR | Introduced | Defines read, gate, write, and refusal classes. |
| CLI output mapping | Existing `respond` / `clierr`; RDR 0005 for user surface | Available / deferred | Accessor package returns values; CLI maps them later. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Structured CLI failure | `internal/cli/clierr::CLIError` | No accessor-specific codes yet | Extend | Add stable refusal codes during implementation or RDR 0005. |
| Output routing | `internal/cli/respond::Fail` | CLI-only; not an internal executor API | Reuse | Accessor package must not print. |
| Config loading | `internal/cli/config::Load` | Project config exists, accessor config not designed | Reuse later | Accessor/table path discovery belongs to CLI integration. |
| Accessor executor | None found under `internal/` | New safety seam | Introduce | New internal package should own capability validation and invocation. |

### Decision Rationale

Declared capabilities best match the user outcome: skills get reliable
passthrough reads, gates, and writes, while the model still prevents accessors
from becoming arbitrary orchestration hooks. The choice aligns with RDR 0001 by
keeping artifact discovery and persistence outside the resolver, with RDR 0002
by treating accessors as declared data, and with RDR 0003 by feeding symbolic
guards rather than callback predicates.

The deciding matrix rejects raw shell-out and FSM action callbacks because they
hide authority in executable code. It also rejects a global allowlist because
allowing `git` or `gh` by name says little about whether a specific accessor may
mutate an authoritative artifact. "No write accessors" would be safer but fails
the system requirement that owned tags can be persisted. Declared read/gate/write
capabilities put the permission decision at the seam where the table references
the external world.

Premortem: this approach ships and fails if read-back checks are too weak, if
every interesting integration demands bespoke command execution, or if timeout
classification collapses distinct failures into one vague error. The
recommendation survives because these are concrete Resolve tasks: fixture the
initial RDR/kata accessors, prove write/read-back mismatch handling, and verify
CLI error mapping before lock. If those fail, the correct rework is to narrow
write accessors further, not to grant raw shell authority.

## Alternatives Considered

### Alternative 1: Declared Read/Gate/Write Capability Accessors

**Description**: Accessor definitions declare a stable name, artifact role,
capability, timeout, and expected tag contract. Runtime invocation enforces that
capability and write accessors must pass read-back verification.

**Pros**:

- Bounds accessor authority per model reference instead of per executable name.
- Keeps resolver behavior replayable and non-orchestrating.
- Makes timeout, gate-indeterminate, execution failure, and read-back mismatch
  visible refusal classes.
- Fits the peer RDR split: table names accessors, predicates consume values,
  this RDR executes safely.

**Cons**:

- Requires a validator and executor instead of a trivial command runner.
- Requires Resolve spikes to prove the vocabulary covers real RDR/kata
  accessors.
- Does not provide general undo; it only verifies the intended owned-tag effect.

**Reason for selection**: Best balance of write safety, reviewability, and
bounded power for authoritative artifacts.

### Alternative 2: Raw Shell-Out Accessors

**Description**: Accessor references contain shell commands or scripts that the
runtime executes for reads, gates, and writes.

**Pros**:

- Most flexible integration surface.
- Easy to prototype against existing CLI tools.
- Can express remote calls without adding typed bindings first.

**Cons**:

- Grants ambient process authority unrelated to the transition model.
- Makes static review nearly impossible; the shell script becomes the real
  contract.
- Timeout and read-back can be bolted on, but command semantics remain opaque.

**Reason for rejection**: It solves integration speed by accepting the exact
blast radius the problem statement asks this RDR to control.

### Alternative 3: Global Allowlisted Commands

**Description**: The model may call only configured executable names or command
prefixes.

**Pros**:

- Safer than raw shell-out.
- Simple to explain and audit at a coarse level.
- Compatible with external tools already installed on user machines.

**Cons**:

- An executable allowlist does not prove the command is read-only, gated, or a
  write to the intended artifact role.
- One allowlist tends to accrete broad authority as new integrations appear.
- Read-back verification is still a separate discipline rather than part of the
  accessor identity.

**Reason for rejection**: It constrains mechanism but not semantic power.

### Alternative 4: No Write Accessors

**Description**: Intrastate only reads and gates; callers persist all owned-tag
writes themselves.

**Pros**:

- Lowest mutation risk inside intrastate.
- Keeps the resolver/accessor boundary very simple.
- Avoids partial-write and read-back mismatch handling.

**Cons**:

- Contradicts the seeded requirement that write accessors mutate authoritative
  artifacts as reliable passthroughs.
- Pushes persistence safety back into every skill or caller.
- Weakens RDR 0001's plan/write split because the returned transition plan would
  have no standard application path.

**Reason for rejection**: It is safe by omission, but incomplete.

### Alternative 5: FSM Action Callback Hooks

**Description**: Treat accessors like state-machine entry/exit/transition
actions and execute host-language callbacks during transition handling.

**Pros**:

- Familiar from state-machine libraries.
- Very expressive and easy to extend in Go code.
- Can share ordinary application helpers.

**Cons**:

- Hides authority in code rather than model data.
- Static lint cannot reason about side effects, timeouts, or read-back behavior.
- Couples transition selection to action execution, undermining the stateless
  resolver boundary.

**Reason for rejection**: It imports the callback power of FSM libraries without
their runtime ownership model, and it breaks the reviewable-data premise.

### Briefly Rejected

- **External workflow harness**: Rejected because intrastate is a thin resolver
  and lint CLI, not an orchestrator that owns credential and tool lifecycle.
- **Transactional store owned by intrastate**: Rejected because RDR 0001 keeps
  authoritative artifacts caller-supplied, and this RDR should not create a
  parallel state store.

## Trade-offs

### Consequences

- Accessor authority becomes explicit model data rather than implicit code.
- Write accessors are slower and more complex than fire-and-forget commands
  because every successful write performs read-back verification.
- The model can refuse gate-indeterminate and timeout cases without guessing a
  transition.
- External integrations remain possible, but only through typed bindings whose
  capability is visible to validation.

### Risks and Mitigations

- **Risk**: The capability vocabulary is too small for real integrations.
  **Mitigation**: Resolve must fixture representative RDR and kata accessors
  before lock and add only proven-needed capabilities.
- **Risk**: A write command succeeds but changes the wrong artifact or wrong
  tag.
  **Mitigation**: Require same-role read-back verification against expected
  owned-tag values.
- **Risk**: Timeout and gate-indeterminate failures are collapsed into generic
  execution errors.
  **Mitigation**: Make them separate refusal classes and verify CLI mapping.
- **Risk**: Accessor bindings smuggle broad shell execution behind typed names.
  **Mitigation**: Validation records capability, artifact role, and tag keys;
  Resolve must inspect the initial binding interface before lock.

### Failure Modes

Visible failures are typed refusals: unknown accessor, capability mismatch,
artifact unavailable, timeout, execution failure, gate denied, gate
indeterminate, write attempted for a non-owned tag, and read-back mismatch.
Silent failure would mean a write command reported success but the owned tag did
not change as expected; the mandatory read-back check is the guard against that
case. Diagnosis starts with the accessor identity, capability, artifact role,
timeout, and expected versus observed tag values.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified
- [ ] RDR 0001 keeps the resolver stateless and returns planned owned-tag
  writes instead of executing persistence.
- [ ] RDR 0002 carries accessor references, tag provenance, and artifact roles
  through normalization.
- [ ] RDR 0003 consumes accessor-produced tag values without executing
  accessors during predicate evaluation.

### Minimum Viable Validation

Build a fixture flow with one read accessor, one gate accessor, and one write
accessor over caller-supplied artifact roles. Prove success, timeout,
gate-indeterminate, execution failure, capability mismatch, and write
read-back-mismatch dispositions. The write success test must assert the
re-read owned-tag value equals the transition plan's expected value.

### Phase 1: Accessor Model

Define the accessor definition structs, capability enum, refusal classes, and
validation rules that connect RDR 0002 accessor references to declared
read/gate/write bindings.

### Phase 2: Executor Boundary

Implement context-bound invocation for typed accessor bindings. The executor
returns structured success/refusal values and performs no direct output.

### Phase 3: Write Read-Back Verification

Apply planned owned-tag writes through write accessors, then re-read the same
artifact role and compare the expected owned-tag values. Treat mismatch as a
write failure.

### Phase 4: CLI Integration Hook

Expose accessor refusal classes to the CLI layer so RDR 0005 can map them
through `respond.Fail` and `clierr.ExitCodeFor` without package-level prints.

### Day 2 Operations

| Resource | List | Info | Delete | Verify | Backup |
| --- | --- | --- | --- | --- | --- |
| Accessor definitions in transition model | Covered by table/model inspection | Covered by table/model inspection | N/A | In scope through validation and MVV fixtures | Covered by caller's artifact/version-control workflow |
| Authoritative artifacts mutated by write accessors | Caller-owned | Caller-owned | Caller-owned | In scope through read-back verification | Caller-owned |

### New Dependencies

No new third-party dependency is selected at Propose. Resolve may choose a Go
test helper or TOML library only if RDR 0002 has not already selected one.

## Validation

### Testing Strategy

[Test scenarios and coverage goals — what to test and
what constitutes "done." For non-functional concerns
(performance, security): state measurement strategy,
not estimates.]

1. **Scenario**: [Description]
   **Expected**: [Result]

### Performance Expectations

[Do not include effort estimates or speculative
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

- [Requirements/standards with section numbers]
- [Dependency docs, source paths reviewed]
- [Dependency repos searched (clone + code search)]
- [Related issues, articles, discussions]
