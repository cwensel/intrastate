# Recommendation 0004: Accessor Execution Safety Model

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-06-19
- **Status**: Draft
- **Type**: Architecture
- **Profile**: large — locks one accessor execution safety contract governing authoritative artifact mutation.
- **Priority**: High
- **Related Issues**: None
- **Predecessors**: 0001-resolution-kernel, 0002-transition-table-as-reviewable-data, 0003-guard-predicate-exhaustiveness
- **Overrides**: None
- **Seam Lineage**: no prior accretion

## Problem Statement

A skill needs read, gate, and write accessors to run as reliable passthroughs: failures must surface, writes must not silently corrupt authoritative artifacts, and external calls must not overstep their intended power. The system-internal requirement is to define accessor execution, failure, timeout, gate-indeterminate, and read-back verification semantics.

## Context

### Background

The resolver is stateless and non-orchestrating, so accessors are the injected I/O seam for owned tags and world facts. Write accessors mutate authoritative artifacts with no ledger or undo, making the safety model part of the core contract rather than an implementation detail.

The real design fork is power versus blast radius: raw shell-out, an allowlisted command set, or a declared-capability model with read, gate, and write authority.

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
- **Verified** — a small declared-capability vocabulary can express the RDR and
  kata accessors without falling back to raw shell strings or host callbacks.
- **Verified** — read-back verification can prove the expected owned-tag effect
  for initial write accessors without needing a full undo log.

### Critical Assumptions

- **A1 The target RDR and kata flows only need declared read, gate, and write
  accessors over caller-supplied artifact roles.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `cd docs/rdr/0004-accessor-execution-safety-model/evidence/spikes && go run .` binds `state.read`, `state.gate`, and `state.persist` as declared read/gate/write accessors over caller-supplied `state` artifacts (`main.go:208-226`); transcript lines 1-8 show read success, gate allow, typed gate/refusal cases, capability mismatch, and write success without raw shell strings or callbacks (`output.txt:1-8`).
  - **If wrong**: The capability vocabulary is too small, and authors will
    pressure the model toward unsafe command execution.
- **A2 Write accessors can verify their intended owned-tag effect by re-reading
  the same artifact boundary after the write.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: The spike writes planned owned tags, clones the same role's observed tags, and returns `read_back_mismatch` when observed values differ (`main.go:121-153`); transcript line 8 shows matching `status=Final`, and line 9 captures the mismatch with expected and observed values (`output.txt:8-9`).
  - **If wrong**: A successful write command could silently corrupt or fail to
    update authoritative state.
- **A3 Timeout, execution failure, gate indeterminate, and read-back mismatch
  can be represented as stable accessor refusal classes and mapped through the
  existing CLI failure gateway.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/cli/clierr::CLIError` defines append-only structured codes/messages with optional detail/hint and non-serialized exit group, `internal/cli/clierr::ExitCodeFor` maps groups to stable exits, and `internal/cli/respond::Fail` emits failures centrally. The spike's refusal enum and output lines cover timeout, execution failure, gate indeterminate, capability mismatch, unknown accessor, and read-back mismatch (`main.go:22-30`, `output.txt:3-9`) without accessor-level printing beyond the test harness.
  - **If wrong**: Accessor errors would need a separate user-facing output
    contract or would leak implementation errors to callers.
- **A4 Accessor execution can be deterministic enough for resolver replay when
  the model records artifact role, accessor name, capability, timeout, and
  returned tag values.**
  - **Status**: Verified
  - **Method**: MVV Test
  - **Evidence**: MVV Scenario 4 is `TestAccessorReplayDisposition`: the spike's `replay` function rebuilds the same fixture artifacts and records accessor name/capability/role/timeout outcomes (`main.go:246-264`), while sorted map formatting prevents map-order drift in the transcript (`main.go:188-206`). Transcript line 10 shows identical success disposition for two runs; line 11 shows an injected gate-indeterminate refusal remains stable (`output.txt:10-11`).
  - **If wrong**: Resolver replay could depend on ambient process state rather
    than declared model inputs.
- **A5 External API accessors can be constrained by declared capability and
  timeout without requiring intrastate to own credentials or remote lifecycle.**
  - **Status**: Verified
  - **Method**: Design Decision
  - **Evidence**: This RDR explicitly scopes credentials and remote resource lifecycle outside intrastate in Cross-Cutting Concerns; accessors receive caller-provided artifacts/environment, declare capability and timeout in model data, and return typed success/refusal only. The selected contract rejects ambient artifact discovery and shell/callback authority in Normative Contracts and Alternatives.
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

- An executable allowlist does not prove the command has read capability, gate
  capability, or write authority over the intended artifact role.
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

- [x] All Critical Assumptions verified
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

The MVV should become production tests around the accessor validator and
executor boundary. The Resolve spike at
`docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/main.go`
already exercises the test matrix as a fixture proof: declared read/gate/write
bindings over caller-supplied artifacts, bounded timeouts, typed refusals,
write read-back verification, and stable replay. Done means those spike cases
become package tests without direct stdout/stderr output from the accessor
package.

1. **Scenario**: Validate a fixture flow with one read accessor, one gate
   accessor, and one write accessor bound to caller-supplied artifact roles.
   **Expected**: Matching capability references pass; unknown accessors,
   duplicate bindings, capability mismatches, missing timeout metadata, and
   write attempts against non-owned tags fail before resolution.
2. **Scenario**: Invoke read and gate accessors that succeed, time out, return
   an execution failure, or return gate indeterminate.
   **Expected**: Successful reads return typed tag values, gate allow/deny
   returns typed gate results, and timeout, execution failure, and gate
   indeterminate remain distinct refusal classes.
3. **Scenario**: Execute a successful transition plan through a write accessor,
   then re-read the same artifact role.
   **Expected**: Matching expected owned-tag values report success; mismatched
   values report read-back mismatch even when command-level write invocation
   succeeded.
4. **Scenario**: Run the same transition model and fixture artifacts twice with
   identical accessor results, then once with an injected refusal.
   **Expected**: The first two runs produce the same disposition; the injected
   failure produces the same stable refusal class and accessor identity.
5. **Scenario**: Route accessor refusals through CLI integration without package
   prints.
   **Expected**: The accessor package returns structured values only; the CLI
   layer can map them through `respond.Fail` and `clierr.ExitCodeFor`.

### Performance Expectations

No throughput target is set. The relevant non-functional check is bounded
execution: the Resolve spike wraps every invocation in `context.WithTimeout`
(`main.go:71-74`), demonstrates timeout as its own refusal (`output.txt:5`),
and shows the write path performs one same-role read-back comparison after a
command-level success (`main.go:121-153`, `output.txt:8-9`). That cost is
acceptable for the representative RDR/kata flow shape because accessors run at
transition boundaries, not inside graph-wide lint loops.

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

No contradictions found between research findings, design principles, and the
proposed solution. The prior-art callback model is cited only as a contrast;
the selected design remains data-declared accessors with typed capabilities and
refusals.

### Assumption Verification

All Critical Assumptions are verified. A1 and A2 are backed by the Resolve spike
and transcript under
`docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/`; A3 is backed
by the existing `clierr.CLIError`, `clierr.ExitCodeFor`, and `respond.Fail`
source; A4 is backed by MVV Scenario 4 and the spike replay transcript; A5 is
the explicit scoping decision that credentials and remote resource lifecycle
remain outside intrastate. None of the evidence cites this RDR or its artifact
directory as self-proof.

### Scope Verification

The MVV is in scope: a fixture flow with read, gate, and write accessors must
prove success, timeout, gate-indeterminate, execution failure, capability
mismatch, and write read-back mismatch dispositions. The implementation tests
must include the named replay scenario and the no-direct-output CLI mapping
scenario.

### Cross-Cutting Concerns

- **Secret/credential lifecycle**: intrastate does not own credentials or remote
  resource lifecycle; external API accessors receive caller-provided
  environment and return typed success/refusal only.
- **Concurrency model**: every invocation is context-bound and has a declared
  timeout; write accessors verify effects through same-role read-back rather
  than relying on fire-and-forget mutation.
- **Determinism**: the RDR claims stable replay disposition, not byte-identical
  output or replay-stable hashes. A4 and the MVV replay scenario must verify
  that identical model inputs and fixture artifacts produce the same success or
  typed refusal.

### Proportionality

This RDR is right-sized by contract count. It owns one load-bearing contract:
accessor execution safety for declared read, gate, and write accessors,
including refusal classes, timeout behavior, and write read-back verification.
RDR 0002 owns the table carrier, RDR 0003 owns predicate semantics, and RDR 0005
owns the user-facing CLI mapping. The `large` Profile is retained because this
contract governs authoritative artifact mutation.

## References

- `docs/cli-output-contract.md`
- RDR 0001: Resolution Kernel
- RDR 0002: Transition Table as Reviewable Data
- RDR 0003: Guard Predicate Exhaustiveness
- Stateless callback prior art: `StateRepresentation::ExecuteEntryActions`,
  `StateRepresentation.Async::ExecuteEntryActionsAsync`
- OpenTofu state-write prior art: `TaintCommand::Run`,
  `stateMgr.WriteState`, `stateMgr.PersistState`
- Local ADO transition helper prior art: `Client::transitionWorkItem`
