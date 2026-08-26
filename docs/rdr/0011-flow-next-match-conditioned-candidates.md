# Recommendation 0011: flow next selects by match; --all enumerates the alphabet

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-08-26
- **Status**: Draft
- **Type**: Feature
- **Profile**: mid — one user-facing verb predicate (`flow next` candidate set + `--all`) that overrides a locked 0005 contract clause
- **Priority**: High
- **Related Issues**: kata `1mv1` (defect tracker); umbrella rdr#thsc; consumer model rdr#tmxk
- **Predecessors**: 0005-skill-integration-cli-contract
- **Overrides**: 0005 "next enumerates rather than selects" clause (`flow next MUST return the legal recognized-outcome alphabet …`) — superseded by a match-conditioned candidate predicate; `flow_next_0005_test.go` tests move with it
- **Seam Lineage**: no prior accretion

## Problem Statement

A skill author (or the human driving one) runs `flow next --model <model> --artifact <state>` to ask "what is legal from here?" so the next action can be chosen without re-deriving the transition table by hand. Today the answer is wrong for that question: RDR 0005 made `next` enumerate rather than select, stripping each row's `match` (and `escape`) before probing, so every row not pruned by a guard is reported as a candidate — 21 of 21 rules on `models/rdr.toml` at `stage=resolved`. The caller discovers this when the candidate list is the whole alphabet regardless of the supplied state, and cannot constrain a skill's next step from it (0005's own premortem, "next omits enough condition detail to constrain a skill", materialized).

The system-internal requirement is a single decision on the candidate predicate of `flow next`: whether "the legal recognized-outcome alphabet for the supplied state" means guard-conditioned (0005's predicate: candidate = not guard-excluded, match is the kernel's selection pattern and `next` must not select) or match-conditioned (candidate = match holds AND not guard-excluded, evaluated by handing the kernel the row with its match intact so selection semantics stay in the kernel); which of the two is the default and which rides `--all`; and whether guard verdicts on match-excluded rows are still reported. A sub-question inside the same contract: whether `ambiguous_match` / `guard_unevaluable` on a match-retained probe leaves the row a candidate (today only `no_match` excludes). `--all` is the same contract's second surface, not a second contract. Help text follows the decision.

## Critical Assumptions

[Required — never omit. Load-bearing assumptions — if
wrong, the approach fails. Each must have a complete
Evidence Record before marking this RDR Final.]

- **A1 [Statement]**
  - **Status**: Verified | Pending | Unverified
  - **Method**: `one of the eight — README
    §Verifying load-bearing claims`
  - **Evidence**: [single sentence — concrete artifact;
    per-method form in README §Verifying load-bearing
    claims. Prefer a stable anchor: `path::Symbol`,
    section heading, REQ/assumption/test ID, grepable
    literal snippet, or artifact path. A bare `file:line` or peer-RDR
    `~line N` is non-normative — drop or rewrite to a
    stable anchor unless the line number **is** the
    behavior under test.
    **Method: Peer RDR cites an element ID, not a record**:
    `cli/0055:C4`, `0055:A3` — the element the claim rests
    on, never the whole file. `rdr inspect NNNN` lists them.
    A filename or heading-text reference is a *mention*:
    fine for context, not for a load-bearing claim.]
  - **If wrong**: [single sentence — what fails; how
    it surfaces to a user or test]
- **A2 [Statement]** — (same shape)

## Proposed Solution

### Approach

[Detailed description of the recommended solution.]

### Technical Design

[Architecture, component relationships, data flow,
extension points.]

#### Normative Contracts

[Required — never omit. Load-bearing — implementers must match exactly.
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

> **Transient marker (bridge surfaces).** A contract block
> for bridge code may carry one line: `Transient — scheduled
> deletion by <sibling NNNN-slug>, <phase/anchor>;
> <one-clause disposition>`. The surface stays named here —
> Profile sizes by blast radius; the marker caps rigor for a
> surface with a scheduled deletion. A `Transient`-marked
> contract counts toward neither the Profile contract axis
> (blast-radius sizing stays on the durable contracts) nor
> the >1-independent-contract split signal above (that
> signal counts *sole-authored* contracts — a bridge whose
> replacement a sibling owns is not sole-authored).

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

State each Normative item in a clearly labeled block.
**Label every block `**C1**`, `**C2**`, … in document
order** — the label is the contract's name for life: peers
cite `NNNN:C2`, and it survives a heading rewrite, a split,
or the contract moving to another RDR. Never reuse a number,
never renumber (a deleted C2 leaves a gap).

**C1**

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

[Conditional — required whenever a load-bearing behavior
depends on a capability not already available (introduced
here, by a predecessor, or deferred); omit (don't
N/A-bullet) this whole section only if every capability
this RDR relies on already exists. For each load-bearing
behavior, state whether the enabling capability exists
now, is introduced by this RDR, is provided by a
predecessor, or is deferred.]

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| [Capability] | Existing / This RDR / Predecessor / Future | Available / Introduced / Deferred | [Impact] |

### Existing Infrastructure Audit

[Conditional — required whenever this RDR proposes a
component that overlaps an existing module; omit (don't
N/A-bullet) this whole section only if this RDR touches no
existing infrastructure. List existing modules that
overlap with proposed components. For each, state whether
to reuse, extend, or replace, and name any known limit
that affects the spec.]

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| [Capability] | [Module/path] | [Limit or none] | Reuse / Extend / Replace | [Impact] |

### Decision Rationale

[Why this approach over alternatives. Key factors,
how it addresses the problem, why alternatives were
ruled out. Closes with Stage 2's two greppable verdict
lines — `Premortem:` and `Joint-check:` — whose absence
means the check never ran.]

## Alternatives Considered

[Full analysis for seriously evaluated alternatives.
One-sentence rejection for trivially eliminated options.]

[Conditional scaffold — omit (don't N/A-bullet) the
`Alternative 1` block below if no alternative warranted
full analysis; the `Briefly Rejected` list alone is fine.]

### Alternative 1: [Name]

**Description**: [Brief description]

**Pros**:

- [Advantage 1]

**Cons**:

- [Disadvantage 1]

**Reason for rejection**: [Why this wasn't chosen]

### Briefly Rejected

- **[Alternative N]**: [One-sentence rejection]

## Context

### Background

Observed 2026-08-26 while driving the consumer model (rdr#tmxk): `flow next --model models/rdr.toml --artifact rdr=<state with stage=resolved>` lists all 21 rules. The cause is by design in RDR 0005: `internal/cli/flow_next.go::excluded` nils `probe.Match` / `probe.Escape` before probing, under the comment "dropping a row whose match pattern the supplied facts do not satisfy would be a selection this verb was not asked to make" — the sentence this RDR revisits. 0005 rationale A3 holds that `next` exposes the alphabet without owning guard evaluation; selection (gate-then-count, exact-one survivor) belongs to the kernel. 0005 also rejected a `resolve` mega-command because it "makes conditional alphabet queries harder" — anticipating this seed. Constraints: RDRs are never amended in content (this is a new RDR that overrides 0005's clause); intrastate stays generic — no consumer (RDR-process) knowledge in code, docs, or fixtures; `make check` must pass. Not a facet of kata `zdat` (owned-state optionality is a model-class decision; this is a verb-predicate decision) — cross-cite only.

### Technical Environment

Go CLI (`bin/intrastate`); `internal/cli/flow_next.go` (`excluded`, candidate/guard reporting), kernel probe API that today receives rows with match stripped, `flow_next_0005_test.go`. Conventions in `AGENTS.md` (respond gateway, CLIError codes, SilenceUsage). Design history: `docs/rdr/0001–0010`, `docs/jdr/0001`.

## Research Findings

### Investigation

[What was analyzed? Code, docs, source, experiments,
standards. Cite specific locations.]

### Key Discoveries

[Label each finding's evidence basis:

- **Verified** — confirmed by spike/POC/experiment
- **Documented** — from official docs or source reading
- **Assumed** — needs validation before implementation]

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

[Required — never omit. What breaks visibly? What fails
silently? Recovery path? How does a developer diagnose
the problem?]

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified
- [ ] [Other prerequisites]

### Minimum Viable Validation

[Required — never omit. The single end-to-end proof that
the approach works. Must be in scope — not deferred.
State it as a stepwise scenario — numbered steps plus the
expected end-state — so the pre-lock desk trace can walk
it.]

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

[Conditional — omit (don't N/A-bullet) this whole section
if this RDR creates no persistent resource. For every
persistent resource this RDR creates (collection, index,
data store, config entry), address management operations:]

| Resource | List | Info | Delete | Verify | Backup |
| --- | --- | --- | --- | --- | --- |
| [Resource] | In scope / Deferred / N/A | ... | ... | ... | ... |

[If any operation is marked "Deferred," justify why
it is not needed for initial usability.]

### New Dependencies

[Conditional — omit (don't N/A-bullet) this section if no
dependency is added or updated. Dependencies to add/update.
For third-party: note license and whether legal review is
required.]

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

> Complete each item with a written response in
> `{ARTIFACT_DIR}/gate.md` before marking this RDR as
> **Final**. Written responses prevent rubber-stamping
> and produce a review record.
>
> First run the mechanical pre-sweep
> (`prompts/gate/tooling-pass.md`): TEMPLATE section
> coverage, Method-label vocabulary, `Source Search`
> self-reference, `Docs Only` on load-bearing claims. It
> catches what the review rounds disturbed; resolve any
> BLOCK before the written responses.
>
> At lock, replace this section's body with the
> one-line pointer to gate.md — responses are never
> inlined. The sub-sections below spec gate.md's
> content.

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
past the lens battery undetected. A `Transient`-marked
contract with a named deleting sibling and schedule is a
recorded lifespan disposition, not an under-sized
Profile — do not count it when re-deriving. Also confirm form:
value + one clause naming the contract(s); strip any
matrix/provenance prose left from the template or Seed
(it belongs in the template comment, not the instance).]

## References

- [Requirements/standards with section numbers]
- [Dependency docs, source paths reviewed]
- [Dependency repos searched (clone + code search)]
- [Related issues, articles, discussions]
