# Recommendation 0026: Deliver the command timeout even when a detached grandchild holds the output pipes

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-08-31
- **Status**: Draft
- **Type**: Bug Fix
- **Profile**: foundational — accretion floor (3 prior point-fixes at the C4 seam, Seam Lineage below); the contract axis alone reads mid: one locked contract amended (0025:C4's execution-safety triple) on a user-facing surface (whether a `timeout` refusal is delivered at all).
- **Priority**: High
- **Related Issues**: intrastate#936e (defect tracker; stays open until Stage 8); intrastate#v0hb (closed, RDR 0025's tracker); closed point-fixes at the seam: intrastate#1drn, intrastate#x5vq, intrastate#878e
- **Predecessors**: 0025-command-invoking-accessor-bindings, 0004-accessor-execution-safety-model
- **Overrides**: 0025:C4 (the "necessary and sufficient" deadline triple — to be amended to a bounded drain with a stated precedence) and 0025:F4's residue class (leakage of a process, never loss of the refusal)
- **Seam Lineage**: `internal/cli/cmdbind::reapGroup` / `internal/cli/cmdbind::drains.Wait` (`area:internal-cli`) — 3 prior closed point-fixes; trail: cd9a09d + intrastate#1drn (process-group syscall build tags), d22cd7a + intrastate#x5vq (grandchild-leak oracle), b4b7431 + intrastate#878e (cancel oracle); plus 76c121b, the manual `os.Pipe` ownership change inside 0025's own implementation (ADV-2 / FAIL-2), which is the change that opened this gap. Count provenance: no `kata-scope-review §seam-accretion` emission exists on intrastate#936e (it was routed by `rdr-seed-triage`); the count above was taken at seed from `kata list --status closed --label area:internal-cli` on 2026-08-31 and is to be confirmed, not re-derived, at Resolve.

## Problem Statement

A model author who declares a `command` reader with a `timeout` (RDR 0025) expects that
when the child misbehaves, `intrastate` returns the `timeout` / `execution_failure`
refusal within a bound they can reason about, and their agent caller moves on. Today a
child that spawns a grandchild which leaves the process group via `setsid(2)` while
holding the inherited stdout/stderr write ends makes the CLI hang **unboundedly**: no
refusal is ever delivered, so a caller (an agent driving `flow resolve` in a pipeline)
waits forever with nothing to branch on. They discover it as a stuck invocation, not as
a refusal.

System-internally, 0025:C4 states a deadline triple (context deadline, `WaitDelay`,
group-scoped `SIGKILL` via `reapGroup`) and calls it "necessary and sufficient", while
also guaranteeing that drains reach a true EOF with nothing truncated. The drain join
(`drains.Wait()`) is unbounded, and the pipes are manually owned (`os.Pipe`, chosen in
76c121b so `Cmd.Wait` cannot truncate an in-flight read), which puts them outside
`Cmd.WaitDelay`'s reach. An escaped grandchild therefore survives `reapGroup`'s
group-scoped kill and keeps the write ends open; nothing bounds the join. 0025:F4's
residue admits only that a *process* may outlive the refusal — never that the
*refusal* is withheld — so this is a liveness hole in C4, not an accepted residue.

The decision this RDR owns: what C4 promises when "drains reach a true EOF, nothing
truncated" and "no unbounded wait" conflict — i.e. the precedence between a bounded
wait and a whole read, and the restated residue class that follows from it. Every
candidate repair changes what C4 says about output completeness, so it is a
contract-level amendment, not a patch.

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
`rdr-0025`) and reproduced at HEAD through the real `cmdbind.Reader.Read`: a Go helper
spawning a `SysProcAttr{Setsid: true}` grandchild that inherits the pipes and sleeps
60 s, driven with a declared `timeout = "2s"`, was still blocked after 15 s — deadline
(2 s) and `WaitDelay` (500 ms) long expired.

Reproduction caveat worth carrying: macOS ships no `setsid` *binary*, so a shell fixture
using `setsid sh -c …` does **not** escape the process group — `reapGroup` kills it, the
drains close, and the read returns bounded in ~200 ms. A first probe of that shape
wrongly showed "no hang". The escape must be made with a real `setsid(2)` call.

Impact: a CLI hang with no refusal, reachable only under the `--allow-commands` opt-in
(0025), by any child whose descendants detach from the group while holding the pipes.
Constraint: the manual `os.Pipe` ownership was introduced deliberately (76c121b, fixing
the sequential-drain deadlock ADV-2 / FAIL-2) so that `Cmd.Wait` cannot truncate an
in-flight read; a repair that hands the pipes back to `Cmd` re-opens that deadlock.
Closing one C4 liveness hole opened another, which is why the precedence must be
decided once rather than patched a third time (see Seam Lineage).

The triage comment on intrastate#936e names the fork the successor RDR must weigh
(timed join accepting partial output; deadline-capable pipe reads; session-level
reaping) and the correlated clauses (F4's residue restated as honest leakage; an
`FX-deadline` scenario in which a `setsid(2)` grandchild holding both pipes must still
refuse within `timeout + bound`). None of that is decided here.

### Technical Environment

Go, `os/exec` with `SysProcAttr` process-group handling behind build tags
(`internal/cli/cmdbind`, Unix-only syscalls tagged after intrastate#1drn; Windows takes
C4's platform refusal). `internal/cli/cmdbind/cmdbind.go`: `Reader.Read`, the
manually-owned `os.Pipe` drains, `drains.Wait()`, `reapGroup` (`syscall.Kill(-pid,
SIGKILL)`). Output caps in play: the 1 MiB stdout cap and 4 KiB stderr tail that C4's
"partial output" wording would have to be read against. Governing records: 0004
(accessor execution safety model), 0025:C4 / 0025:F4, and 0025's implementation
artifacts (`deviations.md`, `verification.md`), which do not record the `os.Pipe`
ownership change against C4.

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
