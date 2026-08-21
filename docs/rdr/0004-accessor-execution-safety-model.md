# Recommendation 0004: Accessor Execution Safety Model

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-06-19
- **Status**: Draft [revised from Final 2026-08-12; re-verified A8 —
  JDR 0001 §D3's read-completeness clause carried its own exactness claim and is
  now backed by an extended Resolve spike; A1-A7 carried forward Verified]
- **Type**: Architecture
- **Profile**: large — one contract: accessor execution safety (capability, refusal classes including read completeness, timeout, and write read-back) governing authoritative artifact mutation.
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
  - **Evidence**: `cd docs/rdr/0004-accessor-execution-safety-model/evidence/spikes && go run .` binds `state.read`, `state.gate`, and `state.persist` as declared read/gate/write accessors over caller-supplied `state` artifacts (`docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/main.go::main`); transcript lines 1-8 show read success, gate allow, typed gate/refusal cases, capability mismatch, and write success without raw shell strings or callbacks (`output.txt:1-8`).
  - **If wrong**: The capability vocabulary is too small, and authors will
    pressure the model toward unsafe command execution.
- **A2 Write accessors can verify their intended owned-tag effect, and the
  absence of collateral change to non-owned tags, by re-reading the same
  artifact boundary after the write.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: The spike writes planned owned tags, clones the same role's observed tags, and returns `read_back_mismatch` when observed values differ (`docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/main.go::write`); transcript line 8 shows matching `status=Final`, and line 9 captures the mismatch with expected and observed values (`output.txt:8-9`).
  - **If wrong**: A successful write command could silently corrupt or fail to
    update authoritative state.
- **A3 Timeout, execution failure, gate indeterminate, and read-back mismatch
  can be represented as stable accessor refusal classes and mapped through the
  existing CLI failure gateway.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/cli/clierr::CLIError` defines append-only structured codes/messages with optional detail/hint and non-serialized exit group, `internal/cli/clierr::ExitCodeFor` maps groups to stable exits, and `internal/cli/respond::Fail` emits failures centrally. The spike's refusal enum and output lines cover timeout, execution failure, gate indeterminate, capability mismatch, unknown accessor, read-back mismatch, and incomplete read (`docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/main.go::refusal`, `output.txt:3-9`, `output.txt:12`) without accessor-level printing beyond the test harness.
  - **If wrong**: Accessor errors would need a separate user-facing output
    contract or would leak implementation errors to callers.
- **A4 Accessor execution can be deterministic enough for resolver replay when
  the model records artifact role, accessor name, capability, timeout, and
  returned tag values.**
  - **Status**: Verified
  - **Method**: MVV Test
  - **Evidence**: MVV Scenario 4 is `TestAccessorReplayDisposition`: the spike's `replay` function rebuilds the same fixture artifacts and records accessor name/capability/role/timeout outcomes (`docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/main.go::replay`), while sorted map formatting prevents map-order drift in the transcript (`docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/main.go::formatMap`). Transcript line 14 shows identical success disposition for two runs; line 15 shows an injected gate-indeterminate refusal remains stable (`output.txt:14-15`).
  - **If wrong**: Resolver replay could depend on ambient process state rather
    than declared model inputs.
- **A5 External API accessors can be constrained by declared capability and
  timeout without requiring intrastate to own credentials or remote lifecycle.**
  - **Status**: Verified
  - **Method**: Design Decision
  - **Evidence**: This RDR explicitly scopes credentials and remote resource lifecycle outside intrastate in Cross-Cutting Concerns; accessors receive caller-provided artifacts/environment, declare capability and timeout in model data, and return typed success/refusal only. The selected contract rejects ambient artifact discovery and shell/callback authority in Normative Contracts and Alternatives.
  - **If wrong**: The accessor layer becomes an orchestrator and should be split
    into a separate integration contract.
- **A6 Accessor definition validation can reject unsafe execution metadata
  before the resolver runs.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `cd docs/rdr/0004-accessor-execution-safety-model/evidence/spikes && go run .` runs `validateDefinitions`, which rejects missing accessors, multiply-bound identities, capability mismatches, missing/non-positive timeout metadata, missing write read-back metadata, ambient artifact discovery, and non-owned writes before runtime (`docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/main.go::validateDefinitions`, `docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/main.go::validationCases`); transcript lines 16-23 capture every validation disposition (`output.txt:16-23`).
  - **If wrong**: Unsafe or ambiguous accessor definitions could reach runtime
    and turn typed refusals into late execution surprises.
- **A7 Write read-back verification can detect unintended mutation of
  non-owned observed or recognized tags on the same artifact role.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `cd docs/rdr/0004-accessor-execution-safety-model/evidence/spikes && go run .` snapshots pre-write tag values, skips the planned owned tag during the non-owned comparison, and returns `read_back_mismatch` when any other observed tag changes (`docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/main.go::write`). Transcript line 10 shows `status=Final` was written as planned while non-owned `profile` changed from `large` to `small`, producing `read_back_mismatch` (`output.txt:10`).
  - **If wrong**: A write accessor could corrupt caller-observed state while
    still passing the owned-tag success check.
- **A8 A read accessor can resolve every requested key or refuse, such that a
  key the artifact genuinely lacks and a key the accessor could not read take
  different branches.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `cd docs/rdr/0004-accessor-execution-safety-model/evidence/spikes && go run .` resolves reads against an explicit requested key set, returning `incomplete_read` when any requested key is unreadable and `<absent>` as a value when the artifact legitimately lacks it (`docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/main.go::read`, `docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/main.go::absentValue`, `docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/main.go::newArtifacts`, `docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/main.go::newPartialArtifacts`, `docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/main.go::newSparseArtifacts`). Transcript line 11 shows a complete read, line 12 the unreadable-key refusal, line 13 the genuine-absence success carrying `profile=<absent>` (`output.txt:11-13`) — the two failure shapes are distinguishable by branch, not by inspecting a thinned value set. The kernel cannot recover the distinction downstream: `internal/resolve/resolve.go::missingOwned` decides on `internal/resolve/resolve.go::TagSet.has` (map presence) alone, so an absent and an unread owned key both become `owned_state_unavailable`, which is why the rule binds at the accessor boundary.
  - **If wrong**: A truncated read reaches a consumer shaped like genuine
    absence, so a presence/absence predicate decides on state that was never
    read.
- **A9 The read-completeness contract's seam and boundary rules are
  implementable as stated: absence crosses to the resolver as omission from the
  owned snapshot, a validated requested-key set is rejectable before execution,
  `timeout` outranks `incomplete_read` on a partial read, and a read-back re-read
  that cannot read a compared key is reportable as `read_back_incomplete`.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: To be verified by the MVV — Scenario 1's eighth validation arm
    (missing/empty requested key set), Scenario 2's absent-required-key case
    carried through to the resolver asserting `owned_state_unavailable`, its
    timeout-with-unresolved-keys case asserting `timeout`, and Scenario 3's
    read-back-with-unreadable-key case asserting `read_back_incomplete`. The
    existing Resolve spike does not witness any of the four: it derives its
    default key set from the artifact, holds absence as an in-map sentinel,
    sleeps before its key loop so timeout and truncation never overlap, and
    re-reads by cloning the tag map without going through `read`.
  - **If wrong**: One of the four rules is unimplementable at a real binding and
    the contract must weaken to what the boundary can observe, most likely by
    moving the absence encoding into RDR 0001's `Input` shape (JDR 0001 §D3
    option (c)) rather than the accessor's output.

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
or `write`. A read accessor returns typed tag values for exactly the
keys it was asked for from a caller-supplied artifact role, or refuses; a
partial read is a refusal, never a thinned value set, and a key the artifact
genuinely lacks crosses the seam as an omission rather than a placeholder. A gate accessor returns allow, deny, or indeterminate with a
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
missing or multiply-bound accessor identities, capability mismatches, writes to
non-owned tags, missing timeout/read-back metadata, non-positive timeouts, a
missing or empty read requested-key set, and ambient artifact discovery before
resolution. Second, runtime invocation applies
read and gate accessors to caller-supplied artifacts and classifies timeout,
unavailable artifact, execution error, incomplete read, gate denied, and
gate-indeterminate outcomes. A read that cannot produce every requested key
takes the refusal branch, so a downstream consumer never sees a truncated
snapshot that is shaped like a genuine absence; when a read exceeds its timeout
before resolving every key, `timeout` is the reported class. The requested keys
come from the definition's validated metadata, never from the keys a read
happened to resolve. Third, write accessors execute only from a successful transition plan.
Before the write, the executor records the same caller-supplied artifact role's
observed and recognized tag values. After command-level success, it immediately
re-reads that same role from the write binding, compares expected owned-tag
values, and verifies the pre-write observed and recognized tag values are
unchanged. A read-back mismatch is a failure even if the write command exited
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
Within one flow, each `(flow id, accessor name, capability)` identity MUST
resolve to exactly one accessor binding. Missing and multiply-bound identities
MUST fail validation before resolution.
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
Every read accessor definition MUST declare the requested key set as validated
metadata. The set MUST NOT be derived from the keys a read actually resolved. A
missing or empty requested key set MUST fail validation before execution.
```

```normative
The typed-tag-values branch carries a completeness guarantee: a read accessor
MUST return the tag set for exactly the keys it was asked for — no requested key
missing, no unrequested key added — or take the refusal branch. A partial or
truncated read is a refusal, not a value. A missing key MUST be distinguishable
from an unread key only by which branch is taken — absence is a value,
unreadability is a refusal. One unreadable requested key MUST refuse the whole
read; the executor MUST NOT return the keys that did resolve.
```

```normative
A read that cannot resolve every requested key MUST be reported as
`incomplete_read`, its own refusal class, distinct from execution failure and
from timeout. The refusal MUST name the requested keys it could not read. When a
read exceeds its timeout before resolving every requested key, `timeout` takes
precedence over `incomplete_read`.
```

```normative
A requested key the artifact genuinely does not carry MUST NOT be presented to
the resolver as a present owned tag. The accessor layer MAY represent absence
however it chooses internally, but what crosses the seam MUST leave the key
absent from the owned snapshot, so that a required absent key resolves as
`owned_state_unavailable` rather than matching against a placeholder value.
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
the same caller-supplied artifact role named by the write binding and verify the
expected owned-tag values and that observed and recognized tag values present
before the write are unchanged. It MUST NOT satisfy read-back verification by
discovering an ambient artifact or by reading an unrelated role. A read-back
mismatch MUST be reported as a write failure.
```

```normative
The read-back re-read is subject to read completeness. If it cannot read a key
it must compare, the write MUST be reported as `read_back_incomplete` — the
verification did not run — and MUST NOT be reported as `read_back_mismatch`,
which asserts the artifact is wrong, nor as success.
```

```normative
Every accessor invocation MUST have a bounded timeout. Timeout MUST be reported
as its own refusal class, distinct from execution failure and read-back
mismatch.
```

```normative
Missing or non-positive timeout metadata MUST fail validation before execution.
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
- **Read completeness** — the success predicate for a read accessor is "every
  requested key resolved", not "at least one key resolved". A key the artifact
  genuinely does not carry resolves as an absent value; a key the accessor could
  not read is an `incomplete_read` refusal. Consumers may therefore treat a
  returned tag set as exactly the requested keys.
- **Requested key set** — the requested keys are the read definition's declared
  metadata, validated before execution. They are never derived from what a read
  resolved: a set computed from the keys that came back makes completeness
  self-fulfilling, so a read that lost keys still reports success over a smaller
  set. The Resolve spike's default path
  (`docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/main.go::expectedTagKeys`)
  is exactly that circular shape — fixture convenience, not the contract.
- **Absence crosses the seam as omission** — how the accessor layer represents an
  absent value internally is free, but what reaches the resolver must leave the
  key absent from the owned snapshot. RDR 0001 already owns that channel:
  `internal/resolve/resolve.go::Refusal.MissingOwned` names "owned tag keys
  absent from the snapshot", and `internal/resolve/resolve.go::missingOwned`
  decides on map presence through `internal/resolve/resolve.go::TagSet.has`. A
  sentinel value would make an absent key read as *present* and silently retire
  `owned_state_unavailable` for it — the one encoding that turns this RDR's
  safety rule into a regression. The spike's `<absent>` string is fixture
  shorthand, not that seam value.
- **Read refusal granularity** — one unreadable requested key refuses the entire
  read rather than returning the keys that did resolve. This is the conservative
  choice and matches JDR 0001's rule that the kernel refuses rather than guesses:
  a decision, not an artifact of the spike's early return.
- **Absent vs unreadable is the binding's call, defaulting to unreadable** — a
  key is *absent* only when the binding read the artifact successfully and the
  key was not there. Every other outcome — the artifact did not parse, the
  transport truncated, the key's value failed type coercion, permission was
  denied — is *unreadable*, because in each the binding did not establish that
  the key is missing. A binding that cannot tell the two apart at its own
  boundary MUST report `incomplete_read`; guessing absence is the failure this
  contract exists to prevent. The Resolve spike declares unreadability as a
  fixture field rather than deriving it, so it proves the branch machinery, not
  a binding's classification.

#### Round-Trip / Inverse Invariants

`write -> read = expected owned-tag value identity + protected non-owned tag
identity`: after a write accessor reports command-level success, reading the
same caller-supplied artifact role named by the write binding must return the
expected owned-tag values, and any observed or recognized tag values present
before the write must remain unchanged. This is not an undo or byte-for-byte
artifact invariant; it is the minimum safety invariant for the authoritative tag
values and protected non-owned tag values this RDR can observe at the accessor
boundary.

#### Disposition Table

Every input class the accessor boundary can meet, and the outcome it takes.
Witness column cites the Resolve spike transcript
(`docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/output.txt`).

| Input class | Branch | Refusal class minted | Silent or loud | Witness |
| --- | --- | --- | --- | --- |
| All requested keys read | success | — (typed tag values) | loud (values returned) | `output.txt:11` |
| Requested key unreadable | refusal | `incomplete_read` | loud | `output.txt:12` |
| Requested key absent from artifact | success | — (absent value; key omitted from the owned snapshot at the seam) | loud one layer down: a row requiring the key refuses `owned_state_unavailable` naming it | `output.txt:13` (spike sentinel; the seam encoding is A9) |
| Accessor exceeds its timeout | refusal | `timeout` | loud | `output.txt:5` |
| Read times out with keys still unresolved | refusal | `timeout` (takes precedence over `incomplete_read`) | loud | not witnessed — see A9 |
| Accessor name not bound | refusal | `unknown_accessor` | loud | `output.txt:6` |
| Accessor used off-capability | refusal | `capability_mismatch` | loud | `output.txt:7` |
| Artifact role not supplied | refusal | `execution_failure` | loud | `output.txt:4` |
| Gate returns indeterminate | refusal | `gate_indeterminate` | loud | `output.txt:3` |
| Gate returns deny | typed gate result | — (deny, with reason) | loud | `output.txt:2` (allow arm) |
| Write succeeds, owned tag matches | success | — | loud | `output.txt:8` |
| Write succeeds, owned tag differs | refusal | `read_back_mismatch` | loud | `output.txt:9` |
| Write succeeds, non-owned tag changed | refusal | `read_back_mismatch` | loud | `output.txt:10` |
| Write succeeds, read-back cannot read a compared key | refusal | `read_back_incomplete` | loud | not witnessed — see A9 |
| Definition invalid (8 shapes) | validation failure | pre-runtime, per shape | loud | `output.txt:17-23`; the missing/empty requested-key-set arm is not witnessed — see A9 |

No input class exits silently: every row either returns a value, mints a named
refusal, or — for a genuinely-absent key — omits that key from the owned snapshot
so the resolver refuses `owned_state_unavailable` naming it. Omission is loud one
layer down, which is the point of the seam clause: the alternative encoding, a
placeholder value, is the one shape that would exit this table silently.
`exit code` is out of scope here — RDR 0005 owns the CLI mapping; this
table stops at the structured value the accessor package returns.

#### Oracle Discriminability

Each MVV scenario against the question "fails if X is wrong because Y", with the
negative control that makes the pass meaningful.

| MVV scenario | Fails if X is wrong because Y | Negative control |
| --- | --- | --- |
| 1 — definition validation | Fails if a validator arm is dropped, because each of the seven invalid fixtures asserts its *own named* code (`missing_accessor`, `multiply_bound_accessor`, `capability_mismatch`, `missing_or_non_positive_timeout`, `missing_write_read_back`, `ambient_artifact_discovery`, `write_non_owned_tag`) — not merely that validation returned non-empty | `validation-ok=` — the valid fixture must return the empty set (`output.txt:16`), so a validator that rejects everything fails |
| 2 — read / gate dispositions | Fails if completeness collapses, because the unreadable-key and genuine-absence fixtures differ in *branch* (`refusal=incomplete_read` vs the absence line), so an implementation that thins the value set instead of refusing produces the absence line for the truncation fixture. The truncation assertion is on the refusal **with an empty value set**, and the requested key set is pinned in the definition independently of artifact contents — otherwise an implementation that derives the request from what it read returns a one-key *success* and passes | `output.txt:11` complete read — an implementation that refuses whenever a key is interesting fails this row. Second control: the derived-key-set implementation must fail the truncation row rather than reporting success over a smaller set |
| 3 — write read-back | Fails if read-back is skipped or scoped to owned tags only, because the two mismatch fixtures differ in *which* tag moved: owned (`status`) at `output.txt:9`, non-owned (`profile`) at `output.txt:10` | `output.txt:8` write success — a read-back that always reports mismatch fails this row |
| 4 — replay stability | Fails if disposition depends on ambient state or map order, because the two identical runs are compared for equality (`replay-identical=true`) rather than for absence of error | The injected-refusal run (`output.txt:15`) — a replay that returns success unconditionally fails it |
| 6 — absence reaches the resolver (A9) | Fails if the accessor emits an absent key as a present tag, because the assertion is on the *resolver's* disposition (`owned_state_unavailable` naming the key), not on the accessor's branch — a placeholder value makes `TagSet.has` true and the refusal never fires | A key the artifact *does* carry must resolve normally through the same path — an implementation that omits every key passes the refusal assertion vacuously |
| 7 — timeout precedence (A9) | Fails if the executor classifies on first-unreadable-key rather than on deadline, because the fixture makes both classes true at once and asserts the name `timeout` | The truncation fixture from scenario 2, which must still return `incomplete_read` — an implementation that always reports `timeout` fails it |
| 5 — no package prints | Fails only by absence-of-output, which is weak. Strengthened at implementation: the accessor package test asserts the returned structured value carries the refusal class, and the no-print property is asserted by capturing stdout/stderr around the call and requiring both empty | A test that deliberately prints must fail the capture assertion; without that control this scenario passes vacuously |

Scenario 5 is the one absence-of-error oracle; the named capture control is what
makes it discriminating and is normative for the implementation test.

#### Fidelity Table

The Round-Trip / Inverse Invariant above, stated as the exact equality each
operation owes and the exemptions that are deliberate.

| Operation | Invariant | Strength | Lossy exemptions |
| --- | --- | --- | --- |
| `write -> read`, owned tags | Every planned owned-tag key/value present in the re-read equals the transition plan's expected value | value equality over the planned key set | none |
| `write -> read`, non-owned tags | Every observed/recognized tag value present before the write is present and unchanged after | value equality over the pre-write snapshot | tags absent before the write are unconstrained; the write binding's own planned owned tags are excluded from this comparison by construction |
| `write -> read`, artifact as a whole | **Not** claimed | — | byte-for-byte artifact identity, formatting, key order, and comments are explicitly out of scope; the accessor observes tag values, not the artifact encoding |
| `read` over a requested key set | The returned set is exactly the requested keys — none missing, none added | branch equality plus key-set equality | how an absent value is represented *inside* the accessor layer is not pinned — the spike's `<absent>` is fixture shorthand; what crosses the seam is pinned: the key is omitted from the owned snapshot |
| `resolve -> replay` | Two runs over the same model and fixture results produce the same disposition | disposition equality | not byte-identical output; ordering is normalized by sorted map formatting |

The deliberate weakening is the third row: this RDR cannot claim artifact-level
fidelity because it observes the boundary at tag granularity. That is the
recorded limit, not an oversight.

#### Desk Trace

The MVV's end state walked stepwise, with every normative assertion in force at
each step and a concrete witness. Read the assertions as the contract clauses
above.

| Step | Assertions in force | Witness |
| --- | --- | --- |
| 1. Validate definitions | exactly-one-capability; identity resolves to exactly one binding; caller-supplied roles (no ambient discovery); timeout metadata present and positive; write read-back metadata present; writes only to owned tags | `validation-ok=` for the valid fixture, seven named codes for the invalid ones (`output.txt:16-23`) |
| 2. Invoke read accessor | read returns typed values or a typed refusal; completeness — exactly the requested keys or refuse; requested set is validated definition metadata, never derived from what resolved; one unreadable key refuses the whole read; absence crosses the seam as omission, not a placeholder value; timeout outranks incomplete read; no mutation of authoritative artifacts; bounded timeout; no direct stdout/stderr | `tags={profile=large,status=Draft}` (`output.txt:11`); `refusal=incomplete_read` (`output.txt:12`); `tags={profile=<absent>,status=Draft}` (`output.txt:13`) |
| 3. Invoke gate accessor | gate returns allow, deny, or indeterminate; indeterminate is refusal-class, never a false allow or deny; bounded timeout | `gate=allow` (`output.txt:2`); `refusal=gate_indeterminate` (`output.txt:3`) |
| 4. Execute planned write | write applies only planned owned-tag writes from a successful transition; never observed or recognized tags | `tags={status=Final}` (`output.txt:8`) |
| 5. Read back the same role | re-read the *same* caller-supplied role named by the write binding; the re-read is itself subject to read completeness — an unreadable compared key is `read_back_incomplete`, not mismatch and not success; verify expected owned-tag values; verify pre-write observed/recognized values unchanged; no ambient or unrelated-role read; mismatch is a write failure | owned mismatch `expected={status=Final} observed={profile=large,status=corrupt}` (`output.txt:9`); non-owned mismatch `expected={profile=large,status=Draft} observed={profile=small,status=Final}` (`output.txt:10`) |
| 6. Return to caller | structured success/refusal values only; timeout is its own class, distinct from execution failure and read-back mismatch | distinct lines for `timeout` (`output.txt:5`), `execution_failure` (`output.txt:4`), `read_back_mismatch` (`output.txt:9`) |

No CONTRADICTION row, on two pairs. The first — "a read accessor MUST return
typed tag values or a typed refusal" versus the completeness guarantee — is
jointly satisfiable because completeness constrains *which* branch the
disjunction takes rather than adding a third branch; step 2's three witnesses
exercise both without conflict. The second is the cross-layer pair the
completeness rule creates: step 2's genuine-absence success against the kernel
behavior A8 cites (`internal/resolve/resolve.go::missingOwned` deciding on map
presence). These collide only if absence reaches the resolver as a present tag,
which the seam clause now forbids — absence crosses as omission, so the success
branch and `owned_state_unavailable` agree rather than compete.

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

**Trade-off**: Most flexible integration surface and the easiest to prototype,
but it grants ambient process authority unrelated to the transition model and
makes static review nearly impossible — the shell script becomes the real
contract.

**Reason for rejection**: It solves integration speed by accepting the exact
blast radius the problem statement asks this RDR to control.

### Alternative 3: Global Allowlisted Commands

**Description**: The model may call only configured executable names or command
prefixes.

**Trade-off**: Safer than raw shell-out and coarsely auditable, but an
executable allowlist does not prove a command has read, gate, or write authority
over the intended artifact role, and one allowlist accretes broad authority as
integrations appear.

**Reason for rejection**: It constrains mechanism but not semantic power.

### Alternative 4: No Write Accessors

**Description**: Intrastate only reads and gates; callers persist all owned-tag
writes themselves.

**Trade-off**: Lowest mutation risk and the simplest boundary, but it
contradicts the requirement that write accessors mutate authoritative artifacts
as reliable passthroughs, pushes persistence safety into every caller, and
leaves RDR 0001's returned transition plan with no standard application path.

**Reason for rejection**: It is safe by omission, but incomplete.

### Alternative 5: FSM Action Callback Hooks

**Description**: Treat accessors like state-machine entry/exit/transition
actions and execute host-language callbacks during transition handling.

**Trade-off**: Familiar and very expressive, but it hides authority in code
rather than model data, leaves static lint unable to reason about side effects,
timeouts, or read-back behavior, and couples transition selection to action
execution.

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
- **Risk**: A read accessor returns a truncated tag set that a consumer cannot
  distinguish from genuine absence, so a presence/absence predicate decides on
  state that was never read.
  **Mitigation**: Completeness is normative — a read that cannot resolve every
  requested key refuses, and the MVV asserts the truncated case takes the
  refusal branch.
- **Risk**: A genuinely-absent key crosses the seam as a present tag with a
  placeholder value, retiring `owned_state_unavailable` for that key and letting
  a guard match the placeholder.
  **Mitigation**: Absence crosses as omission from the owned snapshot; the MVV
  carries an absent required key through to the resolver and asserts the refusal
  still fires.
- **Risk**: Timeout and gate-indeterminate failures are collapsed into generic
  execution errors.
  **Mitigation**: Make them separate refusal classes and verify CLI mapping.
- **Risk**: Accessor bindings smuggle broad shell execution behind typed names.
  **Mitigation**: Validation records capability, artifact role, and tag keys;
  Resolve must inspect the initial binding interface before lock.

### Failure Modes

Visible failures are typed refusals: unknown accessor, capability mismatch,
artifact unavailable, timeout, execution failure, incomplete read, gate denied,
gate indeterminate, write attempted for a non-owned tag, read-back mismatch, and
read-back incomplete.
There are three silent-failure shapes, each with a mandatory guard. A write
command
reports success but the owned tag did not change as expected, or a non-owned tag
changed alongside it — the read-back check covers both. A read returns fewer keys
than requested and the shortfall reads downstream as genuine absence — the
completeness requirement turns that into an `incomplete_read` refusal. A
genuinely-absent key reaches the resolver as a *present* tag carrying a
placeholder, so a key the artifact never carried reads as answered state and
`owned_state_unavailable` never fires — the seam clause requiring absence to
cross as omission is that guard. Diagnosis starts with the accessor identity, capability, artifact role,
timeout, and expected versus observed tag values.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified (A1-A8 Verified; **A9 Pending** — the
  seam/boundary rules the 3amigo pass added, verified by the MVV; see Assumption
  Verification)
- [ ] RDR 0001 keeps the resolver stateless and returns planned owned-tag
  writes instead of executing persistence.
- [ ] RDR 0002 carries accessor references, tag provenance, and artifact roles
  through normalization.
- [ ] RDR 0003 consumes accessor-produced tag values without executing
  accessors during predicate evaluation.

### Minimum Viable Validation

Build a fixture flow with one read accessor, one gate accessor, and one write
accessor over caller-supplied artifact roles. Prove success, timeout, gate
denied, gate-indeterminate, execution failure, incomplete read, capability
mismatch, unsafe definition validation, write read-back-mismatch, and
`read_back_incomplete` dispositions, plus the two A9 boundary cases: an absent
required key refused as `owned_state_unavailable` at the resolver, and a
timed-out partial read classified as `timeout`. The write success test must assert the re-read owned-tag value
equals the transition plan's expected value. The read test must assert that an
accessor which can resolve only some of the requested keys takes the refusal
branch and is distinguishable from one whose artifact genuinely lacks those
keys.

### Phase 1: Accessor Model

Define the accessor definition structs, capability enum, refusal classes, and
validation rules that connect RDR 0002 accessor references to declared
read/gate/write bindings.

### Phase 2: Executor Boundary

Implement context-bound invocation for typed accessor bindings. The executor
returns structured success/refusal values and performs no direct output.

### Phase 3: Write Read-Back Verification

Apply planned owned-tag writes through write accessors, then re-read the same
artifact role and compare the expected owned-tag values plus the pre-write
observed and recognized tag values. Treat mismatch as a write failure.

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
read completeness over an explicit requested key set, unsafe definition
validation, write read-back verification, and stable replay.
Done means those spike cases become package tests without direct stdout/stderr
output from the accessor package. The write success tests must assert the
re-read owned-tag value equals the transition plan's expected value and that
pre-write observed and recognized tag values on the same artifact role remain
unchanged.

1. **Scenario**: Validate a fixture flow with one read accessor, one gate
   accessor, and one write accessor bound to caller-supplied artifact roles.
   **Expected**: Matching capability references pass; unknown accessors,
   duplicate bindings, capability mismatches, missing timeout metadata, and
   write attempts against non-owned tags fail before resolution. Missing or
   non-positive timeouts, missing write read-back metadata, ambient artifact
   discovery attempts, and a missing or empty read requested-key set also fail
   before resolution — eight validation arms, each asserting its own named code.
2. **Scenario**: Invoke read and gate accessors that succeed, time out, return
   an execution failure, return only a subset of the requested tag keys, or
   return gate indeterminate.
   **Expected**: Successful reads return the complete typed tag values for every
   requested key, gate allow/deny returns typed gate results, and timeout,
   execution failure, incomplete read, and gate indeterminate remain distinct
   refusal classes. A read missing any requested key refuses rather than
   returning a partial set, and the refusal carries no values and names the keys
   it could not read; an artifact that genuinely lacks a requested key returns it
   as an absent value, not a refusal. The requested key set is pinned in the
   accessor definition, so an implementation deriving it from what it read fails
   the truncation case rather than reporting success over a smaller set.
   The read dispositions are proven at `output.txt:11-13`: a complete read, an
   `incomplete_read` refusal for an unreadable key, and a genuine absence
   returned as a value (`profile=<absent>`). The refusal payload, the pinned key
   set, and the empty-value assertion are new and unwitnessed — they are A9's.
3. **Scenario**: Execute a successful transition plan through a write accessor,
   then re-read the same artifact role.
   **Expected**: Matching expected owned-tag values report success; mismatched
   owned-tag values or mutated non-owned observed/recognized tag values report
   read-back mismatch even when command-level write invocation succeeded. A
   read-back whose re-read cannot read a key it must compare reports
   `read_back_incomplete` — neither success nor mismatch (A9).
4. **Scenario**: Run the same transition model and fixture artifacts twice with
   identical accessor results, then once with an injected refusal.
   **Expected**: The first two runs produce the same disposition; the injected
   failure produces the same stable refusal class and accessor identity.
5. **Scenario**: Route accessor refusals through CLI integration without package
   prints.
   **Expected**: The accessor package returns structured values only; the CLI
   layer can map them through `respond.Fail` and `clierr.ExitCodeFor`. The test
   asserts this positively — the returned structured value carries the refusal
   class — and captures stdout and stderr around the call, requiring both empty,
   so the scenario cannot pass by absence of output alone.
6. **Scenario**: Read an artifact that genuinely lacks a required owned key, then
   pass the resulting owned snapshot to the resolver (A9).
   **Expected**: The read succeeds at the accessor boundary, the absent key is
   omitted from the owned snapshot rather than carried as a placeholder value,
   and the resolver refuses `owned_state_unavailable` naming that key. This is
   the absence half of read completeness; scenario 2 proves only which branch the
   accessor took.
7. **Scenario**: Invoke a read accessor that resolves some requested keys and
   then exceeds its timeout (A9).
   **Expected**: The refusal is `timeout`, not `incomplete_read` — the two input
   classes overlap and timeout takes precedence.

### Performance Expectations

No throughput target is set. The relevant non-functional check is bounded
execution: the Resolve spike wraps every invocation in `context.WithTimeout`
(`docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/main.go::withTimeout`), demonstrates timeout as its own refusal (`output.txt:5`),
and shows the write path performs one same-role read-back comparison after a
command-level success (`docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/main.go::write`, `output.txt:8-9`). That cost is
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

No contradictions remain between research findings, design principles, and the
proposed solution. The prior-art callback model is cited only as a contrast; the
selected design remains data-declared accessors with typed capabilities and
refusals. The read-accessor clause previously allowed a partial read to take the
typed-values branch, which contradicted this RDR's own rule that an indeterminate
result is a refusal rather than a guess; read completeness is now normative, so
the read and gate clauses state the same rule at both layers.

### Assumption Verification

Critical Assumptions A1-A8 are verified. A1 and A2 are backed by the Resolve
spike and transcript under
`docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/`; A3 is backed
by the existing `clierr.CLIError`, `clierr.ExitCodeFor`, and `respond.Fail`
source; A4 is backed by MVV Scenario 4 and the spike replay transcript; A5 is
the explicit scoping decision that credentials and remote resource lifecycle
remain outside intrastate; A6 is backed by the same Resolve spike's
definition-validation harness and transcript; A7 is backed by the same spike's
collateral-mutation read-back mismatch case, which is also what carries A2's
non-owned-tag half; A8 is backed by that spike's three read dispositions —
complete, unreadable-key refusal, and genuine absence as a value — and grounded
against `internal/resolve/resolve.go::missingOwned`, which shows the kernel
cannot recover the distinction after the fact. None of the verified evidence cites
this RDR or its artifact directory as self-proof.

A9 is **Pending** and is the one assumption this RDR carries unverified into
lock. It books the four rules the pre-lock persona pass added to close the
absence half of read completeness — seam omission, a validated requested-key set,
timeout precedence, and `read_back_incomplete` — none of which the Resolve spike
witnesses. Its method is the MVV rather than a spike extension because each rule
binds at the accessor→resolver boundary the implementation builds, and the
existing fixture double cannot exercise a real binding's classification.

### Scope Verification

The MVV is in scope: a fixture flow with read, gate, and write accessors must
prove success, timeout, gate-indeterminate, execution failure, incomplete read,
capability mismatch, and write read-back mismatch dispositions. The implementation tests
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
Read completeness is part of that same contract — it is the success predicate of
the read capability, not a separate obligation — so carrying it here does not
widen the RDR. The same holds for what the success branch hands across the seam:
naming the resolver-visible encoding of an absent key is what makes the branch
rule mean anything, and it constrains this RDR's own output rather than
reopening RDR 0001's `Input` shape, which is unchanged. RDR 0002 owns the table carrier, RDR 0003 owns predicate
semantics, and RDR 0005 owns the user-facing CLI mapping. The `large` Profile is retained because this
contract governs authoritative artifact mutation.

## References

- `docs/cli-output-contract.md`
- JDR 0001 §D3: read-accessor completeness (`docs/jdr/0001-resolve-kernel-seam.md`)
- RDR 0007: Guard Predicate Totality — its A6b (Pending, downgraded at Stage 6)
  routed the read-completeness obligation to this RDR; the read-completeness
  clause and the seam-omission clause are what discharge it
- RDR 0001: Resolution Kernel
- RDR 0002: Transition Table as Reviewable Data
- RDR 0003: Guard Predicate Exhaustiveness
- Stateless callback prior art: `StateRepresentation::ExecuteEntryActions`,
  `StateRepresentation.Async::ExecuteEntryActionsAsync`
- OpenTofu state-write prior art: `TaintCommand::Run`,
  `stateMgr.WriteState`, `stateMgr.PersistState`
- Local ADO transition helper prior art: `Client::transitionWorkItem`

