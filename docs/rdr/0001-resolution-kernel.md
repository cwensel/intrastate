# Recommendation 0001: Resolution Kernel Contract

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
- **Profile**: large — locks the resolver's deterministic kernel contract and refusal semantics.
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
- **Predecessors**: None
- **Overrides**: None
- **Seam Lineage**: no prior accretion

## Problem Statement

A skill that recognized a typed outcome needs a contracted resolver that returns the next legal state predictably: the same input must produce the same legal output, unmodeled matches must be refused instead of guessed, and unclassifiable outcomes must route to an explicit escape rather than a confident wrong edge. The system-internal requirement is to define the deterministic `resolve` half of recognize-to-resolve as a kernel contract.

## Context

### Background

The build model needs a replay-safe transition kernel for RDR and kata flows. The kernel must treat state as a tag-set, with guards expressed as predicates over tags and tag provenance distinguished as owned, observed, or freshly recognized.

The real design fork is where owned state lives for replay safety: the caller can pass the whole tag-set each call, the kernel can read owned tags via an accessor, or a hybrid can re-derive position while trusting recorded judgments and re-fetching world facts.

### Technical Environment

intrastate is a Go CLI wired through `internal/cli`. The resolver is constrained to be a table plus thin CLI plus lint; it must stay stateless and non-orchestrating, read and persist through passthrough accessors, and refuse illegal or incomplete transition inputs rather than initiating work.

## Research Findings

### Investigation

The proposal is shaped by the current CLI output contract, the seed's
stateless-kernel constraint, and the sibling RDR split already present in this
cluster. `docs/cli-output-contract.md` establishes refusal-first CLI behavior;
RDR 0002 owns the reviewable transition table, RDR 0003 owns guard predicate
shape, RDR 0004 owns accessor safety, RDR 0005 owns the user-facing CLI, and
RDR 0006 owns graph lint. Sibling-path check for resolver selection rules:
`rg -n "resolve|resolver|transition|state|guard|tag|predicate|recognized|outcome|next legal|illegal|refus" internal cmd docs`
found no implemented resolver kernel under `internal/`; it found only the
existing `respond`/`clierr` refusal plumbing and the peer RDR drafts.

### Key Discoveries

- **Documented** — `docs/cli-output-contract.md` already treats unknown CLI
  values as refusals and routes structured failures through the CLI gateway.
- **Documented** — peer RDRs 0002 through 0006 are already seeded around the
  adjacent seams this RDR must not absorb.
- **Assumed** — the accessor, table, and guard seams can provide stable enough
  inputs for deterministic kernel replay; Critical Assumptions A1-A3 carry that
  verification work.

### Critical Assumptions

- **A1 Caller-supplied observed and freshly recognized tags are sufficient
  inputs for non-owned facts.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: Pending: Resolve must name a kernel test where identical
    table, owned snapshot, observed tags, and recognized outcome produce an
    identical resolution without hidden reads.
  - **If wrong**: The resolver would need orchestration authority or ambient
    discovery, breaking replay safety.
- **A2 Owned tags can be read from and persisted to caller-provided artifacts
  through injected accessors without making the resolver itself stateful.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: Pending: RDR 0004 must define artifact-bound accessor
    execution and read-back semantics compatible with a stateless resolver.
  - **If wrong**: The kernel either cannot persist legal transitions or must own
    storage, collapsing the resolver/accessor split.
- **A3 The transition table and guard predicates can expose enough structure for
  deterministic single-edge selection.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: Pending: RDR 0002 and RDR 0003 must define table and predicate
    contracts that let the kernel evaluate candidate edges without host-code
    callbacks.
  - **If wrong**: The kernel cannot prove whether zero, one, or multiple edges
    match, so refusal semantics become guesswork.
- **A4 Refusal outcomes can be represented as typed kernel results and mapped to
  CLI errors outside the kernel.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: Pending: Resolve must confirm `internal/cli/clierr::CLIError`
    and `internal/cli/respond::Fail` can carry the resolver's refusal classes
    without adding kernel-owned output behavior.
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
table, the recognized typed outcome, all non-owned context tags, and the artifact
handles/references that contain owned state. The kernel applies declared
accessors to those caller-provided artifacts to read the owned tag snapshot,
evaluates candidate edges deterministically, and returns either one legal
transition plan or a typed refusal. The kernel does not discover artifacts,
discover work, run subflows, print output, or own persistence. A successful
resolution returns the next state tags and the owned-tag writes that the
accessor layer applies back to the same caller-provided artifact boundary; an
illegal, ambiguous, incomplete, or unclassifiable input returns an explicit
refusal class.

### Technical Design

The kernel is the pure decision boundary between recognition and action. Inputs
are: flow identity, artifact handles supplied by the caller, current owned-state
snapshot obtained by applying accessors to those artifacts, observed tags
supplied by the caller, the freshly recognized outcome tag, and the reviewable
transition table. Evaluation builds a single tag-set view, selects matching
candidate edges, refuses zero or multiple matches unless the table contract
explicitly models an escape edge, and emits a transition plan.

Artifact selection stays outside this RDR's contract. The kernel may describe
owned-tag writes in the returned plan, and the accessor layer knows how to apply
those writes to the caller-provided artifacts, but RDR 0004 owns how those
writes execute and how read-back verification works. The CLI surface, including
text/json output and exit-code mapping, is owned by RDR 0005; the kernel only
exposes structured success/refusal values that the CLI can map.

#### Normative Contracts

```normative
Resolver kernel contract:
Given the same transition table, caller-supplied artifact references, owned tag
snapshot read from those artifacts, caller-supplied observed tags, and freshly
recognized outcome tag, resolve returns the same disposition: exactly one
transition plan or exactly one typed refusal.

The kernel MUST refuse instead of guessing when no edge matches, more than one
edge matches, required owned state is unavailable, a guard cannot be evaluated,
or the recognized outcome is not modeled by the table.

The kernel MUST NOT print output, inspect CLI flags, discover ambient state,
choose artifacts on behalf of the caller, initiate work, or execute persistence
side effects directly.
```

#### Load-Bearing Decisions

- **Identity** — a resolution input is the tuple of flow identity, transition
  table revision, caller-provided artifact reference, owned tag snapshot,
  observed tag-set, and freshly recognized outcome tag. Replaying that tuple
  must replay the disposition.
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

1. Receive the artifact handle/reference from the caller.
2. Read the owned tag snapshot by applying the accessor seam to that artifact.
3. Merge owned, observed, and freshly recognized tags into the evaluation view.
4. Evaluate transition-table guards against that view.
5. Return one transition plan if exactly one legal edge matches.
6. Return a typed refusal for no match, many matches, unavailable state, or an
   unevaluable guard.

Illustrative CLI shape only; RDR 0005 owns the final command syntax:

```sh
intrastate resolve \
  --flow rdr \
  --artifact state:./docs/rdr/0001-resolution-kernel.md \
  --recognized successful
```

When more than one artifact participates, callers pass role-qualified artifacts
rather than implementation labels such as `owned` or `observed`:

```sh
intrastate resolve \
  --flow rdr \
  --artifact state:./docs/rdr/0001-resolution-kernel.md \
  --artifact evidence:./docs/rdr/0001-resolution-kernel/evidence/repeatability/diff.md \
  --recognized successful
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
| Accessor execution | RDR 0004 | Deferred to peer | Kernel applies accessors to caller-provided artifacts and describes writes through this seam. |
| CLI surface | RDR 0005 | Deferred to peer | User-facing commands wrap the kernel but do not define kernel behavior. |
| Static lint | RDR 0006 | Deferred to peer | Lint proves graph safety before runtime; kernel still refuses bad runtime inputs. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Refusal mapping | `internal/cli/clierr` | No resolver-specific codes yet | Extend | Add resolver refusal codes outside the kernel package. |
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
outcome; the kernel obtains owned state through injected accessors, evaluates the
table, and returns a plan or typed refusal.

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

Visible failures are typed refusals: no modeled edge, ambiguous edge, missing
owned state, unevaluable guard, or unmodeled recognized outcome. Silent failure
would mean the kernel guessed a transition or executed persistence directly;
both are prohibited by the normative contract. Diagnosis starts with the input
tuple and the transition table revision used for that resolution.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified
- [ ] RDR 0002, RDR 0003, and RDR 0004 are at least coherent enough to verify the
  table, predicate, and accessor dependencies.

### Minimum Viable Validation

Resolve must name and implementation must add a replay test that feeds the same
table, owned snapshot, observed tags, and recognized outcome to the kernel twice
and asserts byte/value-identical dispositions. The same validation must include
at least one refusal for no match and one refusal for ambiguous matches.

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

- `docs/cli-output-contract.md`
- `docs/rdr/0002-transition-table-as-reviewable-data.md`
- `docs/rdr/0003-guard-predicate-exhaustiveness.md`
- `docs/rdr/0004-accessor-execution-safety-model.md`
- `docs/rdr/0005-skill-integration-cli-contract.md`
- `docs/rdr/0006-graph-lint-authority-and-guarantees.md`
- `internal/cli/clierr/clierr.go`
- `internal/cli/respond/respond.go`
