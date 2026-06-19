# Recommendation 0003: Guard Predicate Exhaustiveness

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
- **Profile**: large — locks the guard predicate grammar and exhaustiveness contract.
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

[What was analyzed? Code, docs, source, experiments,
standards. Cite specific locations.]

### Key Discoveries

[Label each finding's evidence basis:

- **Verified** — confirmed by spike/POC/experiment
- **Documented** — from official docs or source reading
- **Assumed** — needs validation before implementation]

### Critical Assumptions

[Load-bearing assumptions — if wrong, the approach
fails. Each must have a complete Evidence Record
before marking this RDR Final.]

- **A1 [Statement]**
  - **Status**: Verified | Pending | Unverified
  - **Method**: `one of the eight below`
  - **Evidence**: [single sentence — concrete artifact;
    see method-specific guidance below]
  - **If wrong**: [single sentence — what fails; how
    it surfaces to a user or test]
- **A2 [Statement]** — (same shape)

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

[Detailed description of the recommended solution.]

### Technical Design

[Architecture, component relationships, data flow,
extension points.]

#### Normative Contracts

[Load-bearing — implementers must match exactly.
The implementation prompt extracts REQ-N quotes from
this section. This section is also the **authoritative
list of the contracts this RDR owns**: a surface not
named here has no spec to test against, so during
implementation an un-named surface is a deviation, not
free latitude (see `prompts/implementation/launch.md`
Phase 2).]

> **Proportionality (split signal).** Count the
> *independent* load-bearing contracts this RDR is the
> sole author of (a distinct type design, a hash, a wire
> format, a taxonomy, a destructive-op policy each count
> as one). If an implementer would have to hold **more
> than one** such contract in working memory at once,
> this RDR spans more than one seam — split it along those
> seams rather than locking them together. The split test
> is **contract count, not word count**.

- Function/method signatures and type definitions for
  values that cross module boundaries
- Wire-format / on-disk / serialization grammars
- Error envelope shapes and error code enums
- For every introduced user-facing or system-facing
  surface, specify the I/O contract:
  - **Success output**: silent | single value | named
    structured format (link to grammar)
  - **Failure output**: human-readable | structured |
    both (give field-level shape if structured)
  - **Status / sentinel errors**: every distinct code or
    state with one-line user-visible meaning
  - **Preview / dry-run / validation-only mode**: exact
    shape; how it differs from committed success output
  - **Environment divergence**: what changes across
    interactive vs non-interactive, local vs remote,
    batch vs streaming, or equivalent execution modes

State each Normative item in a clearly labeled block,
e.g.:

```normative
func Check(sealed []op.Op, proposed []op.Op) Report
type Report struct { ... }
```

Every external API call inside a Normative block must
have a corresponding Critical Assumption Evidence
Record above (Method: Source Search or Spike, with a
greppable `path::Symbol` or command + output).

#### Load-Bearing Decisions

[Conditional — include only the classes this RDR
touches; omit (don't N/A-bullet) the rest. These four
decision classes are the ones implementation otherwise
invents silently, so each must carry **one explicit
answer** here when in play. This is targeted rigor on
the churn-prone decisions, not blanket detail.]

- **Identity** — what makes two of these things "the
  same"? (the equality/dedup/merge key)
- **Wire / byte format** — the exact layout, or
  explicitly deferred with the named owner.
- **Naming** — the canonical name, and the rejected
  alternatives.
- **Selection / predicate** — when N candidates qualify,
  *which one* is chosen and *why*.

#### Round-Trip / Inverse Invariants

[Conditional — include only if this RDR introduces a
pair of operations expected to compose to identity
(encode/decode, serialize/parse, import/export,
migrate/rollback, snapshot/restore, undo/redo). Omit
otherwise.]

State each invariant explicitly as `X ∘ Y = identity on
input class Z`, and specify the equality as **byte- or
value-for-byte fidelity** — *not* "does not error." A
green exit code does not prove the round-trip preserved
the input; the validation must assert the reconstructed
value equals the original. If the pair spans two RDRs,
also record it as a Critical Assumption with
`Method: Peer RDR` so Stage 8.1 asserts it across the
seam.

#### Illustrative Code

[Shape only — not load-bearing. Use sparingly; prose
is usually clearer.]

- Pseudocode showing algorithmic structure
- Sample invocations showing user-side syntax
- Examples of canonical-form output

Every example, fixture, sample input/output, numeric
count, and platform path is either **Normative** (tests
may assert it; cite the artifact or derivation) or
**Illustrative** (intent only; tests must not assert it
literally).

Do not include full class implementations,
config/schema definitions, or code for deferred
features. Do not annotate Verified/Assumed inside
Illustrative blocks; the surrounding prose makes
assumptions explicit.

### Capability Dependencies

[For each load-bearing behavior, state whether the
enabling capability exists now, is introduced by this
RDR, is provided by a predecessor, or is deferred.]

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| [Capability] | Existing / This RDR / Predecessor / Future | Available / Introduced / Deferred | [Impact] |

### Existing Infrastructure Audit

[List existing modules that overlap with proposed
components. For each, state whether to reuse, extend,
or replace, and name any known limit that affects the
spec.]

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| [Capability] | [Module/path] | [Limit or none] | Reuse / Extend / Replace | [Impact] |

### Decision Rationale

[Why this approach over alternatives. Key factors,
how it addresses the problem, why alternatives were
ruled out.]

## Alternatives Considered

[Full analysis for seriously evaluated alternatives.
One-sentence rejection for trivially eliminated options.]

### Alternative 1: [Name]

**Description**: [Brief description]

**Pros**:

- [Advantage 1]

**Cons**:

- [Disadvantage 1]

**Reason for rejection**: [Why this wasn't chosen]

### Briefly Rejected

- **[Alternative N]**: [One-sentence rejection]

## Trade-offs

### Consequences

[Positive and negative consequences of the chosen
approach.]

- [Consequence 1 — positive or negative]
- [Consequence 2 — positive or negative]

### Risks and Mitigations

- **Risk**: [Description]
  **Mitigation**: [How to address]

### Failure Modes

[What breaks visibly? What fails silently? Recovery
path? How does a developer diagnose the problem?]

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified
- [ ] [Other prerequisites]

### Minimum Viable Validation

[The single end-to-end proof that the approach works.
Must be in scope — not deferred.]

### Phase 1: Code Implementation

#### Step 1: [Title]

[Instructions]

#### Step 2: [Title]

[Instructions]

### Phase 2: Operational Activation

[Deployment, CI/CD, credentials, shared infrastructure.
Omit if not applicable.]

#### Activation Step 1: [Title]

[Instructions]

### Day 2 Operations

[For every persistent resource this RDR creates
(collection, index, data store, config entry),
address management operations:]

| Resource | List | Info | Delete | Verify | Backup |
| --- | --- | --- | --- | --- | --- |
| [Resource] | In scope / Deferred / N/A | ... | ... | ... | ... |

[If any operation is marked "Deferred," justify why
it is not needed for initial usability.]

### New Dependencies

[Dependencies to add/update. For third-party: note
license and whether legal review is required.]

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
