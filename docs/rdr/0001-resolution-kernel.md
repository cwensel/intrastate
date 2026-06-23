# Recommendation 0001: Resolution Kernel Contract

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-06-19
- **Status**: Draft
- **Type**: Architecture
- **Profile**: large — locks one resolver-kernel contract with deterministic disposition and typed-refusal semantics.
- **Priority**: High
- **Related Issues**: None
- **Predecessors**: None
- **Overrides**: None
- **Seam Lineage**: no prior accretion

## Problem Statement

A skill that recognized a typed outcome needs a contracted resolver that returns the next legal state predictably: the same input must produce the same legal output, unmodeled matches must be refused instead of guessed, and escape behavior must be explicit table data rather than a confident wrong edge. The system-internal requirement is to define the deterministic `resolve` half of recognize-to-resolve as a kernel contract.

## Context

### Background

The build model needs a replay-safe transition kernel for RDR and kata flows. The kernel must treat state as a tag-set, with guards expressed as predicates over tags and tag provenance distinguished as owned, observed, or freshly recognized.

The real design fork is where owned state lives for replay safety: the caller can pass the whole tag-set each call, the kernel can own state storage, or a hybrid can use an accessor layer to read owned tags while the caller supplies non-owned facts.

### Technical Environment

intrastate is a Go CLI wired through `internal/cli`. The resolver belongs in an internal package behind the CLI: it must stay stateless and non-orchestrating, consume owned-state snapshots and write targets produced by the accessor layer, and refuse illegal or incomplete transition inputs rather than initiating work.

## Research Findings

### Investigation

The proposal is shaped by the current CLI output contract, the seed's
stateless-kernel constraint, and the sibling RDR split already present in this
cluster. `docs/cli-output-contract.md` establishes refusal-first CLI behavior;
RDR 0002 owns the reviewable transition table, RDR 0003 owns guard predicate
shape, RDR 0004 owns accessor safety, RDR 0005 owns the user-facing CLI, and
RDR 0006 owns graph lint. Reuse audit against the Stage 4 paths found no
implemented resolver kernel, transition table, guard evaluator, or accessor
executor under `internal/`; the only reusable surfaces are the existing
`respond`/`clierr` refusal plumbing and the peer RDR drafts. Audit command:
`rg -n "resolve|resolver|transition|state|guard|tag|predicate|recognized|outcome|next legal|illegal|refus" internal cmd docs`.

### Key Discoveries

- **Documented** — `docs/cli-output-contract.md` already treats unknown CLI
  values as refusals and routes structured failures through the CLI gateway.
- **Documented** — peer RDRs 0002 through 0006 are already seeded around the
  adjacent seams this RDR must not absorb.
- **Verified** — the accessor, table, and guard seams define the structures this
  kernel consumes: RDR 0002 defines normalized exact-one candidate rows, RDR
  0003 defines symbolic guard predicates without host callbacks, and RDR 0004
  defines caller-artifact accessors with read-back verification.

### Critical Assumptions

- **A1 Caller-supplied observed and freshly recognized tags are sufficient
  inputs for non-owned facts.**
  - **Status**: Verified
  - **Method**: MVV Test
  - **Evidence**: Minimum Viable Validation names `Replay the same legal input
    tuple twice`: identical table, owned snapshot, observed tags, and recognized
    outcome must return value-identical transition plans without hidden reads.
  - **If wrong**: The resolver would need orchestration authority or ambient
    discovery, breaking replay safety.
- **A2 Owned tags can be read from and persisted to caller-provided artifacts
  through the accessor layer without making the resolver itself stateful.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0004 `Approach`, `Technical Design`, and `Normative
    Contracts` define read/write/gate accessors over caller-supplied artifact
    roles, planned owned-tag writes only, same-role read-back verification, and
    structured success/refusal values with no direct output.
  - **If wrong**: The resolver cannot rely on explicit owned-state snapshots and
    planned writes without owning storage, collapsing the resolver/accessor
    split.
- **A3 The transition table and guard predicates can expose enough structure for
  deterministic single-edge selection.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0002 `Approach` and `Normative Contracts` define
    normalized candidate rows, source identity, and exact-one row matching; RDR
    0003 `Approach` and `Normative Contracts` define closed symbolic predicate
    atoms over declared tags and reject host predicates before resolution.
  - **If wrong**: The kernel cannot prove whether zero, one, or multiple edges
    match, so refusal semantics become guesswork.
- **A4 Refusal outcomes can be represented as typed kernel results and mapped to
  CLI errors outside the kernel.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/cli/clierr/clierr.go::CLIError` defines extensible
    error codes and exit groups; `internal/cli/respond/respond.go::Fail` emits
    failures through the CLI gateway; `internal/cli/root.go::ExecuteAndEmit`
    converts command-level failures through the same gateway.
  - **If wrong**: The kernel would need CLI-specific behavior or the CLI would
    lose stable error mapping.

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

Use a hybrid, stateless resolver kernel: the caller supplies the transition
table, the recognized typed outcome, all non-owned context tags, and the owned
tag snapshot already read from caller-provided artifacts by the accessor layer.
The kernel evaluates candidate edges deterministically and returns either one
legal transition plan or a typed refusal. The kernel does not discover
artifacts, execute accessors, discover work, run subflows, print output, or own
persistence. A successful resolution returns the next state tags and the
owned-tag writes that the accessor layer applies back to the same
caller-provided artifact boundary; an illegal, ambiguous, incomplete, or
unmodeled input returns an explicit refusal class.

### Technical Design

The kernel is the pure decision boundary between recognition and action. Inputs
are: flow identity, transition table revision, owned-state snapshot values
produced by the accessor layer, observed tags supplied by the caller, the
freshly recognized outcome tag, and the reviewable transition table. Evaluation
builds a single tag-set view, selects matching candidate edges, refuses zero or
multiple matches unless the table contract explicitly models an escape edge,
and emits a transition plan.

Artifact selection and accessor execution stay outside this RDR's contract. The
kernel may describe owned-tag writes in the returned plan, and the accessor
layer knows how to apply those writes to caller-provided artifacts, but RDR 0004
owns how reads, writes, and read-back verification execute. The CLI surface,
including text/json output and exit-code mapping, is owned by RDR 0005; the
kernel only exposes structured success/refusal values that the CLI can map.

#### Normative Contracts

```normative
Resolver kernel contract:
Given the same flow identity, transition table revision, accessor-produced owned
tag snapshot, caller-supplied observed tags, and freshly recognized outcome tag,
resolve returns the same disposition: exactly one transition plan or exactly
one typed refusal.

The kernel MUST refuse instead of guessing when no edge matches, more than one
edge matches, required owned state is unavailable, a guard cannot be evaluated,
or the recognized outcome is not modeled by the table.

Modeled refusal is a value-level resolver disposition, not a CLI error and not
the Go error path for parser bugs, IO failures, or programmer mistakes. The
kernel-owned refusal kind set is exactly:

- `no_match`
- `ambiguous_match`
- `owned_state_unavailable`
- `guard_unevaluable`
- `unmodeled_outcome`

Each refusal kind must be stable enough for RDR 0005 to map to a CLI error code
without inspecting error strings.

The kernel MUST NOT print output, inspect CLI flags, discover ambient state,
choose artifacts on behalf of the caller, initiate work, or execute persistence
side effects directly.
```

#### Load-Bearing Decisions

- **Identity** — a resolution input is the tuple of flow identity, transition
  table revision, accessor-produced owned tag snapshot, observed tag-set, and
  freshly recognized outcome tag. Replaying that tuple must replay the
  disposition.
- **Naming** — the internal component is the resolver kernel. Rejected names:
  "orchestrator" because it implies initiating work, and "state machine runner"
  because it implies owning persistence.
- **Selection / predicate** — the only successful selection is exactly one
  matching edge after guard evaluation. Zero, multiple, unavailable, or
  unevaluable candidates are refusals unless the table contains a modeled escape
  edge that itself matches exactly once.

#### Round-Trip / Inverse Invariants

No encode/decode, import/export, or inverse operation is introduced by this
RDR. Replay determinism is covered instead by A1 and the Minimum Viable
Validation.

#### Illustrative Code

Illustrative algorithm only:

1. Receive an owned tag snapshot from the accessor layer.
2. Merge owned, observed, and freshly recognized tags into the evaluation view.
3. Evaluate transition-table guards against that view.
4. Return one transition plan if exactly one legal edge matches.
5. Return a typed refusal for no match, many matches, unavailable owned state,
   unevaluable guard, or unmodeled recognized outcome.

Illustrative CLI flow only; RDR 0005 owns the final command syntax and may split
artifact reads from pure resolution:

```sh
intrastate read-state \
  --flow rdr \
  --artifact state:./docs/rdr/0001-resolution-kernel.md

intrastate resolve \
  --flow rdr \
  --owned status:Draft \
  --recognized successful
```

When more than one artifact participates, callers pass role-qualified artifacts
rather than implementation labels such as `owned` or `observed` to the accessor
verbs:

```sh
intrastate read-state \
  --flow rdr \
  --artifact state:./docs/rdr/0001-resolution-kernel.md \
  --artifact evidence:./docs/rdr/0001-resolution-kernel/evidence/repeatability/diff.md
```

In this shape, the flow definition maps artifact roles to accessors. The
`state` role can read and write owned RDR state; the `evidence` role can expose
read-only observed tags such as repeatability iteration count. The kernel sees
the resolved tags and write targets, not user-facing `owned`/`observed` flags.

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| CLI output and error envelope | Existing `respond` / `clierr` | Available | Kernel refusal values must stay mappable to CLI errors without printing. |
| Transition table | RDR 0002 | Deferred to peer | Kernel assumes a parsed, reviewable table shape. |
| Guard predicate evaluation | RDR 0003 | Deferred to peer | Kernel assumes guards are evaluable without arbitrary callbacks. |
| Accessor execution | RDR 0004 | Deferred to peer | Accessor layer reads owned tags before resolution and applies planned writes after a successful transition. |
| CLI surface | RDR 0005 | Deferred to peer | User-facing commands wrap the kernel but do not define kernel behavior. |
| Static lint | RDR 0006 | Deferred to peer | Lint proves graph safety before runtime; kernel still refuses bad runtime inputs. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Refusal mapping | `internal/cli/clierr` | No resolver-specific CLI codes yet | Extend | RDR 0005 maps kernel refusal kinds to CLI error codes outside the kernel package. |
| Output routing | `internal/cli/respond` | CLI-only gateway | Reuse | Kernel must return values, not write output. |
| Command wiring | `internal/cli` | Only version/config scaffolding exists | Extend later | RDR 0005 owns command shape. |
| Resolver implementation | none found under `internal/` | No existing kernel | Introduce | New internal package can own pure resolution logic. |

### Decision Rationale

The hybrid stateless kernel preserves replayability without making callers
reconstruct authoritative owned state themselves. It aligns with the existing
CLI output prior by treating illegal inputs as structured refusals, while
keeping text/json and exit-code concerns outside the kernel. It also keeps the
cluster split clean: RDR 0001 owns "given inputs, choose or refuse"; RDRs 0002,
0003, 0004, 0005, and 0006 own data shape, predicates, accessors, CLI, and lint.

Premortem: this approach fails if the accessor seam makes owned-state reads
non-repeatable, if the transition table cannot expose enough structure to prove
single-edge matches, or if callers omit observed context that the kernel cannot
discover. The recommendation survives because those are exactly the assumptions
handed to Resolve and the peer RDRs; moving the state store into the kernel would
hide the same risks while increasing blast radius.

## Alternatives Considered

### Alternative 1: Caller Supplies The Whole Tag-Set

**Description**: The caller passes owned, observed, and freshly recognized tags
on every call; the kernel does no accessor reads and only evaluates the table.

**Pros**:

- Maximally pure and simple to replay.
- Keeps persistence completely outside the resolver.

**Cons**:

- Makes every caller responsible for reconstructing authoritative owned state.
- Increases the chance that skill code drifts into local state interpretation.

**Reason for rejection**: It pushes the hardest replay-safety decision to every
caller instead of centralizing it behind the accessor seam.

### Alternative 2: Kernel Owns State Storage

**Description**: The resolver reads and writes authoritative flow state directly,
using its own storage model as part of resolution.

**Pros**:

- Centralizes owned-state consistency.
- Gives the kernel complete information for every decision.

**Cons**:

- Makes the kernel stateful and non-replayable without storage fixtures.
- Couples resolver semantics to storage, migration, and corruption handling.
- Overlaps RDR 0004's accessor safety model.

**Reason for rejection**: It solves missing context by expanding kernel authority
past the seed's stateless, non-orchestrating boundary.

### Alternative 3: Hybrid Accessor-Backed Kernel

**Description**: The caller supplies non-owned context and the recognized
outcome; the accessor layer supplies the owned state snapshot, and the kernel
evaluates the table and returns a plan or typed refusal.

**Pros**:

- Keeps the resolver stateless while avoiding duplicated owned-state assembly in
  every caller.
- Preserves deterministic replay by making the owned snapshot an explicit input
  to validation.
- Leaves persistence safety, CLI output, and static lint in their peer seams.

**Cons**:

- Depends on the accessor contract being stable and read-back-verifiable.
- Requires tests to capture the exact owned snapshot used for replay.

**Reason for selection**: It best matches the user's outcome: deterministic legal
next-state selection with explicit refusal and without turning the kernel into
an orchestrator.

### Briefly Rejected

- **Skill-local transition logic**: Rejected because it recreates the original
  problem of scattered, guessed edges.
- **Host-language guard callbacks in the kernel**: Rejected here because RDR
  0003 must first define a predicate shape that lint can reason about.

## Trade-offs

### Consequences

- Skills get one deterministic resolver behavior instead of re-implementing
  transition choices.
- Runtime refusal remains necessary even after lint, because live inputs can be
  incomplete or unavailable.
- The kernel cannot answer questions that require ambient discovery; callers or
  accessors must provide those facts explicitly.

### Risks and Mitigations

- **Risk**: Accessor reads are not stable enough for replay.
  **Mitigation**: Resolve A2 against RDR 0004 before locking this RDR.
- **Risk**: Multiple edges match and callers expect priority ordering.
  **Mitigation**: Make multiple matches a refusal; any priority or escape must
  be explicit table data owned by RDR 0002.
- **Risk**: CLI callers treat refusal as a crash.
  **Mitigation**: RDR 0005 maps kernel refusals through the established
  `CLIError`/`respond.Fail` contract.

### Failure Modes

Visible failures are typed refusal values: `no_match`, `ambiguous_match`,
`owned_state_unavailable`, `guard_unevaluable`, or `unmodeled_outcome`. Silent
failure would mean the kernel guessed a transition or executed persistence
directly; both are prohibited by the normative contract. Diagnosis starts with
the input tuple, the refusal kind, and the transition table revision used for
that resolution.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified
- [ ] RDR 0002, RDR 0003, and RDR 0004 are at least coherent enough to verify the
  table, predicate, and accessor dependencies.

### Minimum Viable Validation

Resolve must name and implementation must add a replay test that feeds the same
table, owned snapshot, observed tags, and recognized outcome to the kernel twice
and asserts value-identical dispositions. The same validation must include at
least one value-level refusal each for `no_match`, `ambiguous_match`,
`owned_state_unavailable`, `guard_unevaluable`, and `unmodeled_outcome`, and
must assert those modeled refusals do not use the CLI or Go error path.

### Phase 1: Kernel Boundary

Define the internal resolver package boundary, result taxonomy, and pure
resolution entry point without CLI output or persistence side effects.

### Phase 2: Evaluation Semantics

Implement tag-set assembly, guard evaluation delegation, exact-one edge
selection, and typed refusal behavior.

### Phase 3: Replay Validation

Add focused kernel tests for deterministic replay and refusal classes, using
fixtures that exercise owned, observed, and freshly recognized tags.

### Phase 4: CLI Integration Handoff

Expose only the kernel values needed by RDR 0005; do not add command output or
state mutation semantics in this RDR.

### New Dependencies

No third-party dependency is proposed at this stage.

## Validation

### Testing Strategy

Implementation must add focused resolver-kernel tests before any CLI command
wraps the kernel. The tests exercise the pure decision boundary: fixture
transition tables, owned snapshots from accessor fixtures, observed tags
supplied by the caller, and freshly recognized outcome tags. The matrix is
grounded in A1's MVV replay evidence, A2's RDR 0004 accessor seam, A3's RDR
0002/0003 exact-one predicate seam, and A4's existing CLI refusal gateway.

1. **Scenario**: Replay the same legal input tuple twice.
   **Expected**: Both calls return value-identical transition plans, including
   next tags and accessor write descriptions.
2. **Scenario**: No table edge matches the assembled tag-set.
   **Expected**: The kernel returns the typed no-match refusal and performs no
   persistence side effect.
3. **Scenario**: More than one table edge matches after guard evaluation.
   **Expected**: The kernel returns the typed ambiguous-match refusal unless one
   explicit escape edge matches exactly once.
4. **Scenario**: The freshly recognized outcome is not modeled by the table.
   **Expected**: The kernel returns the typed unmodeled-outcome refusal and
   performs no persistence side effect.
5. **Scenario**: Owned state is unavailable or a guard cannot be evaluated.
   **Expected**: The kernel returns the corresponding value-level typed refusal
   and does not fall back to ambient discovery or the CLI/Go error path.

### Performance Expectations

No throughput target is part of this RDR. Resolution is bounded by the supplied
transition table and tag snapshots. No byte-stable hash or
canonical serialization is introduced; determinism is value-level replay of the
input tuple named in A1. Implementation should keep the kernel allocation-light
and deterministic, then benchmark only if the RDR or kata tables become large
enough to make table scans visible in normal command latency.

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

- `docs/cli-output-contract.md`
- `docs/rdr/0002-transition-table-as-reviewable-data.md`
- `docs/rdr/0003-guard-predicate-exhaustiveness.md`
- `docs/rdr/0004-accessor-execution-safety-model.md`
- `docs/rdr/0005-skill-integration-cli-contract.md`
- `docs/rdr/0006-graph-lint-authority-and-guarantees.md`
- `internal/cli/clierr/clierr.go`
- `internal/cli/respond/respond.go`
