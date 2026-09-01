# Recommendation 0027: Authority surface of the C5 interpreter deny-list

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-08-31
- **Status**: Draft
- **Type**: Architecture
- **Profile**: mid — one contract (the predicate shape and claim wording of 0025:C5's `command_shell_interpreter` deny-list), locking what a reviewer-visible lint promises.
- **Priority**: Medium
- **Related Issues**: intrastate#zvtg (defect tracker; stays open until Stage 8); intrastate#v0hb (closed, RDR 0025's tracker); intrastate#q2q1 (open — the contained `env` option-flag sub-fix, ships independently as a plain kata and forecloses none of this RDR's answers); the charted `stdin = "none"|"envelope"` successor in RDR 0025's `evidence/critique/Charted.md`
- **Predecessors**: 0025-command-invoking-accessor-bindings
- **Overrides**: 0025:C5's "env chains" wording and the deny-list's stated authority surface — to be narrowed, extended to a wrapper-aware argv model, or demoted to a non-load-bearing lint hint, whichever this RDR decides
- **Seam Lineage**: `internal/table/load.go::interpreterForm` (`area:internal-table`) — no prior accretion. Count provenance: no `kata-scope-review §seam-accretion` emission exists on intrastate#zvtg (routed by `rdr-seed-triage`); taken at seed from `kata list --status closed --label area:internal-table` on 2026-08-31 — the one closed C5 fix at this file, intrastate#b84g, is a different symbol (placeholder-category validation, C5 clause 3), and intrastate#q2q1 is open.

## Problem Statement

A reviewer reading a model under `--allow-commands` (RDR 0025) relies on `intrastate
lint`'s `command_shell_interpreter` refusal to make inline shell visible: a `command`
binding that hands code to `sh -c` must be spelled where the review can see it, not
smuggled. Today that reviewer is told `["env","sh","-c",…]` and `["env","A=1","sh",
"-c",…]` are refused, while `["env","-i","sh","-c",…]`, `["env","-u","FOO","sh",
"-c",…]`, `["nice","sh","-c",…]`, `["timeout","5","sh","-c",…]`, `["xargs","sh",
"-c",…]` and `["sh","-s"]` lint clean and execute a shell. They discover it only by
reading argv themselves — the thing the lint exists to spare them — or by an
`env -i sh -c` that reads as deliberate evasion of a reviewer-visible check.

System-internally, `interpreterForm` walks an `env` chain by skipping `env` and bare
`NAME=VALUE` tokens and stopping at the first non-assignment token, so option flags are
invisible to it. Part of that is a plain conformance gap: 0025:C5 names "env chains" as
in scope and 0025's `verification.md` records the `env`-chain form as a Phase 3a PASS,
yet only assignment-only chains are caught (the contained flag-walk fix is
intrastate#q2q1, not this RDR). The rest is not covered by C5's text at all: `nice`,
`timeout`, `xargs`, `nohup`, `setsid`, `stdbuf`, `chpst`, `doas` each re-open an
argv-position-aware wrapper class, and `sh -s` / `sh <script` / `env -S "sh -c …"` are
a different axis — shell-string forms no argv-level check can see without modelling
the wrapper's own parser.

The decision this RDR owns, exactly once: what is the *authority surface* of C5's
deny-list, and therefore what does it promise? The triage names three defensible
answers — (1) an argv0-name heuristic, narrowed in wording and with the wrapper and
shell-string forms declared out of scope by name; (2) an argv-position-aware wrapper
model that closes the class rather than enumerating binaries, with a stated bound on
the wrapper set; (3) retiring argv as the authority surface in favour of the charted
`stdin = "none"|"envelope"` successor, keeping the deny-list as an explicitly
non-load-bearing hint. It is a policy/scope decision about predicate shape plus the
claim wording verification is held to — not an implementation detail — and is not
decided here.

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

Found by the roborev triage of RDR 0025's implementation (`src:roborev`, batch
`rdr-0025`) with real lint runs. Denied, correctly: `env sh -c`, `env A=1 sh -c`,
`bash --login -c`. Admitted: `env -i sh -c`, `env -u FOO sh -c`, `nice sh -c`,
`timeout 5 sh -c`, `xargs sh -c`, `sh -s`; `env -i sh -c 'echo …'` was verified to
actually execute the shell.

Why this is not simply C5 being open: C5's "interpreter set: OPEN (deny-listed, not
closed) — an unlisted interpreter is admitted" covers Phase 3a's `python3` / `nodejs`
non-finding (unlisted argv0 *spellings*). `env -i sh -c` is different — `sh` *is*
listed and C5's own text names `env` chains as in scope — so the option-flag blindness
is a conformance gap against a claim C5 makes, while the `nice`/`timeout`/`xargs`
family is genuinely outside C5's text.

Threat model and impact: low on its own — C5 explicitly declines to make inline shell
impossible and states the deny-list's job as raising the cost, so a bypass voids no
promised guarantee; reachable only under the `--allow-commands` opt-in. The cost is to
reviewer trust in a check whose purpose is visibility.

Constraints carried from the triage: the `env` option-flag walk (`-i`, `-u NAME`,
`-0`, `-C`, `--`) is intrastate#q2q1 and ships independently without foreclosing any
answer; answer 3 may fold into the charted stdin/envelope successor, so that successor
should be paired with, not duplicated by, this RDR; whether `env -S "sh -c …"` is in
scope at all is an open question the chosen answer must settle.

### Technical Environment

Go. `internal/table/load.go`: `interpreterForm` (the `env`-chain walk) and the
`command_shell_interpreter` refusal class it feeds, evaluated by `intrastate lint`
over `command` bindings admitted under `--allow-commands` (RDR 0025). Governing
records: 0025:C5 (the deny-list contract and its "raise the cost" purpose), 0025's
`verification.md` (Phase 3a's "env chains are caught" PASS text), and the charted
`stdin = "none"|"envelope"` successor in 0025's `evidence/critique/Charted.md`, which
already concedes argv is not a sufficient authority surface.

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
