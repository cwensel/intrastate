# Recommendation 0029: What a released version number promises an agent about the machine-readable output

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-09-11
- **Status**: Draft
  <!--
  - `Deferred` is the parked-with-a-revisit-trigger status for a
    Draft that cannot proceed because **no acceptable mechanism
    exists yet** — every in-our-control path is ruled out and the
    one that would work is outside our control. It is a *pause in
    the lifecycle*, not an exit from it: the RDR stays intact and
    re-enters at the stage it stopped when the trigger fires.
    Carry the condition on the live value:
    `Deferred [revisit when <condition>]`, and say in the same
    field what was ruled out and why (Alternatives Considered
    carries the long form). Distinct from `Abandoned`, which is
    terminal — an Abandoned RDR is closed, owes a post-mortem, and
    never re-enters. A Deferred RDR owes **no** post-mortem
    (nothing was implemented), and its `Priority` records what the
    fix is *worth*, not what is scheduled. Do not defer merely to
    park work that is possible but unfunded — that is a Priority,
    not a Status.
  - `Demoted` is the terminal status for an RDR judged
    *not RDR-shaped* — the decision was never a real
    design fork, so it leaves the RDR lifecycle and is
    refiled as a plain issue. Carry the destination on the
    live value: `Demoted [→ <issue link>]`, and record the
    same link under **Related Issues**. A `Demoted` RDR runs
    no further stages. (Distinct from the 07.1 *demotion*
    below, which is a `Final → Draft` flip that keeps the
    RDR in the lifecycle — that flip never writes
    `Status: Demoted`; see the disambiguation note there.)
  - A Draft demoted from Final by the 07.1 cluster gate
    carries a qualifier on the live value:
    `Draft [revised from Final YYYY-MM-DD; re-verify A2,A4
    @refine — <one-line reason>]` — the one place this
    grammar is spelled: the date, `;`, `re-verify <IDs>` or
    `re-verify none`, then `@<stage>` naming the target
    re-entry stage (`@propose` = Stage 2, `@refine` = 3,
    `@resolve` = 4, `@finalize` = a RE-LOCK-ONLY re-entry
    whose note-listed wording fixes land in the lock pass),
    then `— <reason>`. Every clause after
    the date is optional; `/rdr-status` routes on `@<stage>`
    and falls back to `/rdr-resolve` when it is absent. It
    is still a `Draft` for every binary Draft/Final gate;
    only Stage 4 (scoped re-verify) and Stage 7 (re-lock)
    parse the rest of the qualifier.
    The Stage 7 flip to `Final` overwrites the whole value,
    so the qualifier self-clears at re-lock — no separate
    cleanup. This 07.1 "demotion" is a *verb* describing the
    Final→Draft flip; it is **not** the `Demoted` status
    above (which exits the lifecycle to an issue) — do not
    conflate the two. (`Reverted` above is the unrelated
    terminal "implementation rolled back" status — also do
    not conflate.)
  - A Draft sent BACKWARD by a mid-flow stage carries the
    sibling qualifier:
    `Draft [routed back from resolve YYYY-MM-DD; re-verify
    A5,A6 @propose — <one-line reason>]` — the origin stage
    that sent it back, the date, `;`, `re-verify <IDs>` or
    `re-verify none`, then the SAME `@<stage>` slot spelled
    above (here also `@prelock` = Stage 5 and `@reconcile`
    = 6, which a demotion can never name because a Final has
    closed them), then `— <reason>`. Every clause after the
    date is optional, and `/rdr-status` routes it through the
    same rules. The **origin** is carried because the rework
    owed differs by where it came from; without the qualifier
    the record's evidence reads FORWARD and the receiving stage
    refuses it as already-passed. Unlike the demotion form it
    does not wait for a re-lock — **the stage named by
    `@<stage>` clears it as its first act**. A second
    route-back overwrites the value rather than stacking.
  - A Final tolerated at the 07.1 gate under a JOINT-DECISION
    carries `Final [joint decision → <home §-anchor>: <the
    open question>]`. It is still a `Final` for every binary
    gate. The qualifier is an **open obligation, not a
    coherence claim**: it says the named question is
    unanswered here, not that this RDR agrees with the answer.
    So it does not self-clear. When the home answers, this RDR
    owes a scoped check of that answer against its own
    normative fences before it re-locks or implements —
    consistent → drop the qualifier and record the clearing;
    contradicts fenced text → a 07.1 SPEC-DEFECT. A re-lock
    that comes first carries the qualifier forward unchanged;
    it is never silently dropped.
  -->
- **Type**: Architecture
- **Profile**: mid — provisional; locks the stability class of the
  `--as=json` envelope, the refusal-code vocabulary, and the lint
  finding taxonomy; user-facing yes; locks contract.
  <!-- Do not paste the matrix below into the field; it is the
  Stage 5 routing latch, provisional on `Draft`, made
  authoritative by Resolve.
  Sized by BLAST RADIUS — the MAX of two axes, not
  contract count or word count.
  (1) contract axis: resolved by rdr-write.toml's `profile`
  rows (the rule's one home) from the durable contract
  count and two dispositions written in the clause:
  `user-facing <yes|no>; locks <none|contract|format|cross-rdr>`.
  (2) accretion axis (a one-tier raise): resolved from `Seam
  Lineage` below by the routing model (rdr-status.toml's
  `floor` group — the rule's one home), never re-read here.
  Matrix: rdr/stages/README.md. Seed estimates from the design
  shape; Resolve overwrites from the verified count; Stage 8
  Gate locks it at Draft → Final. Never skip lenses off a
  Draft Profile until Resolve has run. -->
- **Priority**: High
- **Related Issues**: none
- **Seam Lineage**: no prior accretion

## Problem Statement

> Synthesized from a free-form idea rather than a kata — flagged for author
> review.

An agent author pins a released `intrastate` version in a skill, a CI job, or
a harness, and writes code that parses `--as=json`: it branches on the
envelope's `type`, matches a refusal `code`, checks an exit status, and —
if it lints models — expects a model that passes today to pass tomorrow. The
agent then upgrades, because an install step resolved a newer release or a
package manager pulled one. It needs to know, before upgrading, what that
version number said was allowed to change. It discovers the gap the hard way:
its parse silently takes the wrong branch, or its CI goes red on a model
nobody edited.

The audience makes this sharper than it would be for a human-facing CLI. A
human reader degrades gracefully when a field is renamed or a new diagnostic
appears — they read the message and adapt. An agent hard-fails, or worse,
mis-parses without failing. So the promise a version number carries has to be
stated over the *machine-readable* surfaces specifically, and stated somewhere
an agent can find it.

Internally, `intrastate` already has three such surfaces, and they already move
independently in code:

- the `--as=json` envelope (`type`, `data`, `notes`, `warnings`), the refusal
  `code` vocabulary, and the exit-code classes;
- the model TOML's `[model] version` integer, pinned at `1`, with a load
  refusing anything else;
- the lint finding taxonomy, whose per-finding code identity RDR 0017 fixed.

No record says what a release version promises about any of them, nor which
mechanism proves the promise held. The open fork this RDR must settle: a
**new** lint finding code removes nothing and renames nothing, yet a model
that lints clean at one version starts refusing at the next — breaking a
consumer's pipeline. Semantic versioning has no opinion on whether that is a
minor or a major; this project must have one, and every future change to the
envelope or the taxonomy will need to read it.

Out of scope: the release *mechanism* — tag format, goreleaser, semver
arithmetic, and starting from a `0.x` series. Those are settled by existing
tooling and wider convention, and belong in ordinary documentation rather than
a decision record.

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
Phase 2). Stage 5 closes this section with one line,
outside the fences: either `Determinacy: fired — <contracts>`
or else `Determinacy: n/a — <reason>` (prompts/pre-lock/
3-repeatability.md §Determinacy trigger).]

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
never renumber (a deleted C2 leaves a gap). A clause labelled
at column zero inside the fence (`L-3  …`) is cited as
`NNNN:L-3`; keep clause labels unique across the record.

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
`Method: Peer RDR` so Stage 7.1 asserts it across the
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

[Conditional scaffold — this block is a per-instance slot, not a
section every RDR owes: the heading is the author's own and the
block is omitted (never N/A-bulleted) when unused.]

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

Raised while judging whether the project is stable enough to cut its first
release. The build and publish machinery is already in place and unexercised:
`.goreleaser.yaml` builds five targets, `.github/workflows/release.yml`
publishes a GitHub Release on a `v*` tag after re-running `make check` against
the tagged commit, `ci.yml` carries a gated snapshot job, and
`internal/version` stamps build identity via ldflags with a VCS fallback that
`make release-check` guards. No tag has ever been cut, so
`intrastate version --as=json` currently reports a pseudo-version.

What is absent is not machinery but a stated promise. The distribution surface
is deliberately consumer-facing already — `.goreleaser.yaml` calls the archive
`name_template` "a compatibility surface, not cosmetics" — and `llms.txt`
tells agents the binary is authoritative and to prefer asking it over reading
any file. Neither says what stays true across versions.

The seed of the fork is in the existing records: RDR 0017 gave each lint
finding its own code identity, which is what makes "a new finding code
appeared" a distinguishable event a consumer can trip over. RDR 0023 and 0028
moved envelope shapes recently enough that the wire contract is visibly still
in motion.

### Technical Environment

Go 1.26.3, `github.com/cwensel/intrastate`, a cobra command tree over a pure-Go
binary (no cgo). The relevant surfaces:

- `internal/cli/respond/` — the output gateway every verb routes through;
  `Success{Type, Notes, Warnings, Data}` is the terminal JSON envelope.
- `internal/cli/clierr/` — structured `CLIError` plus exit-code mapping; the
  refusal-code vocabulary and the `findings` carrier.
- `internal/table/load.go` — refuses a model whose `[model] version` is not
  `1`, reporting the `unsupported_version` category from
  `internal/table/category.go`.
- `internal/version/` — ldflags-stamped build identity, with a
  `debug.ReadBuildInfo` fallback.
- `docs/cli-output-contract.md` — the worked payloads and envelope rationale;
  `docs/cli-reference.md` and `llms.txt` are generated from the command tree
  and gated against staleness by `make docs-check`.

Distribution today is GitHub Releases plus `go install`. Package-manager reach
(a Homebrew tap, an install script, any npm wrapper) is a separate, dependent
question — it widens the audience that would be holding whatever promise this
RDR settles, but it does not decide the promise.

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

[Conditional scaffold]

[Instructions]

#### Step 2: [Title]

[Conditional scaffold]

[Instructions]

### Phase 2: Operational Activation

[Conditional scaffold]

[Deployment, CI/CD, credentials, shared infrastructure.
Omit if not applicable.]

#### Activation Step 1: [Title]

[Conditional scaffold]

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
> At lock, replace Contradiction Check, Assumption
> Verification, Scope Verification and Proportionality
> with the one-line pointer to gate.md — those four
> judge THIS record at THIS lock and no peer cites
> them. **Cross-Cutting Concerns stays here**, below
> the pointer: it names the project-wide policy other
> RDRs conform to, so it must stay projected and
> citable as `cli/NNNN:G-cross-cutting`. Cite it that
> way, not by section name.

### Contradiction Check

[Gate key: contradiction — a gate response is cited as
`cli/NNNN:G-<key>`, so the key is a stable id and is
not derived from this heading, which may be reworded.]

[State any conflicts between Research Findings and
the Proposed Solution. If none exist, state
"No contradictions found between research findings,
design principles, and proposed solution."]

### Assumption Verification

[Gate key: assumptions]

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

[Gate key: scope]

[Confirm the Minimum Viable Validation is in scope
and will be executed during implementation, not
deferred. State the specific test or proof.]

### Cross-Cutting Concerns

[Gate key: cross-cutting]

[Retained at lock — this sub-section stays in the RDR
when the other gate responses move to gate.md, because
peer RDRs cite it as `cli/NNNN:G-cross-cutting` and an
element that is not projected cannot be cited.]

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

[Gate key: proportionality]

[Is the document right-sized for the change? Flag
any sections that should be trimmed before locking.
The split test is **contract count, not word count**:
confirm this RDR is the sole author of at most one
independent load-bearing contract (per the Normative
Contracts split signal). If it owns more than one
seam, flag it for splitting rather than locking the
seams together.

Re-validate the **Profile** Metadata field: re-run
rdr-write's `--outcome profile` with the clause's own
dispositions and confirm the value Resolve wrote is
what it emits (a stop is not a match).
If the lenses that actually ran disagree with the
Profile (e.g. Profile says `small` but the change locks
a contract that warranted `mid`+ lenses, or the lenses
were skipped on a wrong `small`), correct the field and
do not lock until the missing lenses have run. This is
the latch's backstop — a wrong Profile cannot route
past the lens battery undetected. (The row already
excludes `Transient`-marked contracts.) Also confirm form:
value + one clause naming the contract(s) and its two
dispositions; strip any
matrix/provenance prose left from the template or Seed
(it belongs in the template comment, not the instance).]

## References

- [Requirements/standards with section numbers]
- [Dependency docs, source paths reviewed]
- [Dependency repos searched (clone + code search)]
- [Related issues, articles, discussions]
