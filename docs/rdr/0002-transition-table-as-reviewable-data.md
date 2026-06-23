# Recommendation 0002: Transition Table As Reviewable Data

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
- **Profile**: large — locks the sparse transition-model data format.
- **Priority**: High
- **Related Issues**: None
- **Predecessors**: 0001-resolution-kernel
- **Overrides**: None
- **Seam Lineage**: no prior accretion

## Problem Statement

A flow author needs every legal edge to live in one reviewable artifact, so "is this transition legal" has one answer instead of being reconstructed from scattered prose. The system-internal requirement is to choose a table representation that can encode match predicates, multi-tag writes, and accessor references for both RDR and kata flows.

## Context

### Background

The transition table is the flow's design and must be hand-authored from the legal graph audits, not generated. The RDR table and kata table are the test workload, including multi-tag writes that a scalar state token cannot represent.

The real design fork is representation: TOML, JSON, CSV, or a small DSL, with trade-offs among diff reviewability, parse simplicity, expressiveness, and lintability.

### Technical Environment

intrastate is a Go CLI wired through `internal/cli`. The table format feeds the resolver kernel and the static lint, and must remain practical to parse, validate, and review in normal code review.

## Research Findings

### Investigation

The proposal is shaped by RDR 0001's stateless resolver split, the current CLI
output contract, and the transition-model prior that treats state as a tag-set
matched by reviewable rows. The design has to keep table authors in normal code
review instead of making them read generated code or scattered prose.
Sibling-path check for an existing table/selection signal:

```sh
rg -n "resolve|resolver|transition|state|guard|tag|predicate|recognized|outcome|next legal|illegal|refus" internal cmd docs
```

The search found no implemented transition table or resolver package under
`internal/`; it found only the existing CLI refusal plumbing and the peer RDR
drafts. Prior-art reading favors a data table that looks like "match predicates
-> tag writes" and reserves external FSM libraries for graph validation, not
runtime orchestration.

Direct `arc` searches over `StateMachineOS`, `StateMachineLit`, and `DevRef`
added four design constraints. Sismic keeps nested state source separate from
rendered PlantUML output and carries transition contracts as preconditions,
postconditions, and invariants. `transitions` models guards as positive
`conditions` and negative `unless` lists evaluated as one predicate set.
Stateless builds a symbolic `StateGraph` for diagramming from machine metadata,
including superstates, stay transitions, and decision nodes, instead of making
the graph renderer the runtime. The literature search surfaced statecharts as
the formal answer to dimensional state explosion: hierarchy and extended state
factor common context instead of enumerating every Cartesian row.

### Key Discoveries

- **Documented** — RDR 0001 delegates the reviewable transition-table contract
  to this RDR and requires enough structure for deterministic single-edge
  selection.
- **Documented** — the CLI output contract already establishes refusal-first
  behavior; table parse and lint failures can flow through the existing
  structured error gateway rather than inventing table-specific output.
- **Verified** — Go's TOML tooling can preserve a sparse authoring schema's
  ergonomics for representative malformed-row diagnostics; the Resolve spike
  identified root-key placement as the exact field-layout constraint to lock.
- **Verified** — the RDR and kata legal graphs can be expressed as sparse rules
  that normalize to explicit row candidates without requiring host-code
  callbacks or an embedded expression language.
- **Verified** — hierarchical/shared contexts plus positive/negative guard lists
  are enough factoring to avoid RDR's status/profile/prelock Cartesian explosion
  without importing a full statechart runtime.

### Critical Assumptions

- **A1 TOML can represent sparse transition rules with nested match predicates,
  multi-tag writes, and accessor references without ambiguous decoding.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `cd docs/rdr/0002-transition-table-as-reviewable-data/evidence/spikes && go run . rdr-fixture.toml kata-fixture.toml` parsed both fixtures with `github.com/pelletier/go-toml/v2`, including inherited match contexts, nested predicates, accessor references, multi-tag writes, and explicit clears; transcript captured in `docs/rdr/0002-transition-table-as-reviewable-data/evidence/spikes/output.txt`.
  - **If wrong**: The chosen carrier either loses table semantics or forces a
    custom parser earlier than intended.
- **A2 Row order is not part of successful edge selection.**
  - **Status**: Verified
  - **Method**: Design Decision
  - **Evidence**: Normative Contracts state that source order and rendered-row order MUST NOT decide transition success; the only successful selection is exactly one normalized candidate row, with zero or multiple matches as refusals. This explicitly rejects first-match semantics.
  - **If wrong**: Reviewers would have to reason about hidden priority, and
    reordering rows could silently change resolver behavior.
- **A3 RDR and kata flow edges can be encoded sparsely with fixed predicate
  operators rather than host-code callbacks or Cartesian-product row
  enumeration.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `docs/rdr/0002-transition-table-as-reviewable-data/evidence/spikes/rdr-fixture.toml` and `kata-fixture.toml` encode representative rules covering status, profile, prelock iteration, equality, set membership, integer comparison, self-loop, rewind, positive/negative guards, and multi-tag writes without host-code callbacks; `output.txt` dumps the expanded candidate rows.
  - **If wrong**: RDR 0003 must expand the predicate grammar or this table format
    becomes too weak for the target flows.
- **A4 The model data can carry enough provenance to separate owned, observed,
  and recognized tags for lint and accessor binding.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0003 `Technical Design` requires declared tag provenance for predicate lint and owned-tag read-before-write checks; RDR 0004 `Technical Design` makes RDR 0002 responsible for table-carried accessor references while accessor execution validates capabilities and writes only owned tags.
  - **If wrong**: The resolver cannot prove read-before-write or accessor safety
    from the model alone.
- **A5 Standard parse and validation failures can be surfaced through the
  existing CLI error envelope.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/cli/clierr/clierr.go::CLIError` already carries stable `Code`, `Message`, optional `Param`, diagnostic `Detail`, `Hint`, exit-code `Group`, and `Cause`; `internal/cli/respond/respond.go::Fail` emits that envelope in text/json modes; `internal/cli/config/config.go::Load` already uses stable config load/read error codes (`config-not-found`, `config-read-error`) and marks `config-invalid` as the planned parse-validation path. Table-specific parse and lint failures must add their stable codes during implementation.
  - **If wrong**: Table loading would need a separate user-facing error contract
    owned by this RDR or RDR 0005.
- **A6 Shared contexts and positive/negative guard lists are sufficient to keep
  the RDR model sparse without hiding ambiguity.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `rdr-fixture.toml` encodes `draft -> prelock -> large-prelock` inherited contexts plus `all.iter.lt = 3` and `unless.profile.eq = "small"` guards; `output.txt` shows the normalized row with inherited status/stage/profile predicates and combined `all`/`unless` predicates, preserving ambiguity visibility in the expanded row.
  - **If wrong**: The model either needs a richer statechart-like hierarchy or
    the table becomes too repetitive for reliable human review.
- **A7 Deterministic expanded-table ordering is a format contract, not an
  implementation accident.**
  - **Status**: Verified
  - **Method**: Design Decision
  - **Evidence**: Normative Contracts define the expanded table as the
    normalized candidate-row value containing model id, row identity, source
    locator, predicates, and writes; they require dump ordering to sort rows by
    row identity and then sort predicate and write keys within each row. This
    explicitly rejects source-order, map-iteration, and renderer-specific
    ordering.
  - **If wrong**: Golden tests and review dumps could churn across machines or
    refactors even when the transition semantics are unchanged.

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

Use a sparse, hand-authored TOML transition model that normalizes to explicit
candidate rows. Authors write shared tag declarations, reusable match contexts,
outcome groups, guarded rules, and tag writes; the tool normalizes that source
into a full row set for lint, resolver lookup, and table dumps.
The rendered table is review support, not the source authors maintain.

Each flow has a named model with flow-level tag declarations, optional accessor
references, a legal recognized-outcome alphabet, and sparse rules that encode
the match predicates and tag writes produced by legal edges. The source model
must provide enough structure for RDR 0001's exact-one resolver contract: zero
matches, multiple matches, unknown outcomes, unavailable accessors, and
malformed rules remain typed refusals.

The model is data, not generated code and not a runtime FSM engine. RDR 0001's
resolver consumes the normalized representation, RDR 0003 owns the fixed
predicate operators, RDR 0004 owns accessor execution and read-back safety, RDR
0005 exposes the CLI, and RDR 0006 owns graph lint. This RDR owns the on-disk
sparse representation and the normalized expanded-table view those peers
consume.

### Technical Design

The source schema has six conceptual parts:

1. Flow metadata: table id, version, and human description.
2. Tag declarations: tag name, provenance (`owned`, `observed`, `recognized`),
   value kind, and optional accessor reference for observed or owned read-back.
3. Recognized outcome alphabet: the closed set of outcome tags that the
   recognizer may emit for the flow.
4. Shared match contexts: named predicate blocks for dimensions such as RDR
   `Status`, `Profile`, stage, prelock iteration, cluster eligibility, or kata
   lifecycle. Contexts may inherit from another context to model statechart-like
   hierarchy without adopting a statechart runtime.
5. Sparse transition rules: reviewable rule ids, optional context references, a
   local match block, positive `all` guards, negative `unless` guards, and a
   write block containing one or more tag assignments, plus an optional
   explicit clear list.
6. Render settings: deterministic ordering and field selection for the expanded
   table dump.

The parser turns TOML into typed source data, then normalizes it into explicit
candidate rows. Validation rejects malformed tags, unsupported model versions,
unknown predicate operators, writes to non-owned tags, rules that match on
missing tag declarations, unresolvable context references, unknown accessor
names, and ambiguous overlaps between candidate rows. Runtime matching is
deliberately priority-free: the resolver evaluates candidate rows against the
supplied tag-set and succeeds only when exactly one row matches. Multi-tag writes
are first-class because RDR rewinds and kata lifecycle moves need to set both the
next stage/state and side-channel scope tags in one edge.

Guard factoring follows the `transitions` prior art: positive predicates and
negative predicates are authored separately but normalized into one predicate
set. Contract factoring follows the Sismic prior art: entry preconditions,
postconditions, and invariants are different validation classes, not free-form
comments. Rendered dumps follow the Stateless/Sismic export pattern: they are
symbolic views derived from model metadata and must carry enough source ids to
send diagnostics back to the authored sparse rule.

The internal representation may be indexed as a decision tree, trie, or decision
DAG for efficient lookup, but that is an implementation detail. The normative
semantic object is the normalized candidate-row set plus its source locator back
to the sparse TOML rule. The locator must identify at least the model id and rule
id; byte line/column coordinates are optional diagnostic detail. This keeps
diagnostics tied to the authored source while letting the resolver avoid
scanning irrelevant dimensions.

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

```normative
The transition model MUST be authored as sparse TOML data, not generated code
and not a fully expanded Cartesian-product table.
```

```normative
The source schema MUST use the Resolve spike field layout: root `outcomes`,
`[model]`, `[tags.<tag>]`, `[accessors.<id>]`, `[context.<id>]`, `[[rule]]`,
and `[dump]`. Context predicates live under `[context.<id>.match.<tag>]`; rule
predicates live under `[rule.match.<tag>]`, `[rule.guard.all.<tag>]`, and
`[rule.guard.unless.<tag>]`; writes live under `[rule.write]`; explicit clears
live in a rule-level `clear` list.
```

```normative
`[model]` MUST contain `id` and `version`. Version `1` is the only version this
RDR accepts; any other version MUST be refused before normalization.
```

```normative
Each transition rule MUST contain a stable rule id, zero or more shared-context
references, a local match block, and a write block. A rule MAY contain a
rule-level explicit clear list. A write block MAY assign more than one tag.
```

```normative
Shared contexts MAY inherit from other contexts, but inheritance MUST normalize
to an explicit predicate set before lint or resolution.
```

```normative
Guard predicates MUST be represented as positive `all` predicates and negative
`unless` predicates. Normalization MUST combine both into one candidate-row
predicate set before ambiguity checks.
```

```normative
The tool MUST normalize the sparse source into deterministic candidate rows for
lint, resolver lookup, diagnostics, and table dumps. Each candidate row MUST
retain its source rule id and source locator.
```

```normative
The expanded table dump MUST be derived from the normalized candidate-row value:
model id, row identity, source locator, predicates, and writes. Dump ordering
MUST be deterministic across source key order by sorting rows by row identity,
then sorting predicate and write keys within each row.
```

```normative
Source order and rendered-row order MUST NOT decide a successful transition. A
tag-set resolves only when exactly one normalized candidate row matches; zero or
multiple matches are refusals.
```

```normative
The model MUST declare every tag it matches or writes, including each tag's
provenance: owned, observed, or recognized.
```

```normative
Clearing a tag MUST be represented by an explicit rule-level `clear` entry that
normalization renders as a `<clear>` write. Absence from both the write block and
the clear list MUST NOT imply deletion.
```

```normative
Validation failures MUST retain stable data-level categories before CLI mapping,
including at minimum unknown tag, unknown context, write to non-owned tag,
unknown accessor, unsupported version, and ambiguous overlap.
```

#### Load-Bearing Decisions

[Conditional — include only the classes this RDR
touches; omit (don't N/A-bullet) the rest. These four
decision classes are the ones implementation otherwise
invents silently, so each must carry **one explicit
answer** here when in play. This is targeted rigor on
the churn-prone decisions, not blanket detail.]

- **Identity** — a transition rule is identified by `(model id, rule id)`.
  Normalized candidate rows inherit that identity plus a deterministic expansion
  suffix. Rule ids are stable review anchors and must not be reused for a
  different edge.
- **Wire / byte format** — TOML is the on-disk carrier. The exact field names
  are the Resolve spike layout: root `outcomes`, `[model]`, `[tags.<tag>]`,
  `[accessors.<id>]`, `[context.<id>]`, `[[rule]]`, `[rule.write]`,
  rule-level `clear`, and `[dump]`. `[model].version = 1` is the only accepted
  format version. The RDR and kata spike fixtures are the canonical examples
  implementation tests must promote.
- **Naming** — the canonical source artifact name is "transition model"; the
  canonical rendered view is "expanded transition table." Rejected
  names: "state machine config" because it suggests a runtime driver, and
  "workflow graph" because this RDR owns sparse transition data, not
  orchestration.
- **Selection / predicate** — the only successful selection is exact-one row
  match after normalization. If multiple candidate rows qualify, the model is
  ambiguous; the resolver refuses and lint should reject the overlap before
  runtime.

#### Round-Trip / Inverse Invariants

`parse ∘ normalize ∘ dump = expanded-table value identity` on valid transition
model fixtures: dumping the normalized model and reading the dump as a table
view must preserve the candidate-row set, including row identity, source
locator, predicates, and writes. Source rewrite is out of scope for this RDR; a
later rewrite-capability RDR must define its own source-preservation invariant
before mutating authored TOML.

#### Illustrative Code

[Shape only — not load-bearing. Use sparingly; prose
is usually clearer.]

Illustrative sparse source shape, not a locked schema:

```toml
[context.draft]
status.eq = "Draft"

[context.prelock]
inherits = "draft"
stage.eq = "prelock"

[[rule]]
id = "prelock-flapping-cap"
use = ["prelock"]

[rule.match]
outcome.eq = "verdict-flapping"

[rule.guard.all]
iter.lt = 3

[rule.guard.unless]
profile.eq = "small"

[rule.write]
stage = "prelock"
```

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Stateless exact-one resolution | RDR 0001 | Pending | This RDR must provide normalized candidate rows and tag writes the kernel can evaluate. |
| Fixed predicate operators | RDR 0003 | Pending | This RDR names predicate slots but does not own the operator grammar. |
| Accessor references and safe read-back | RDR 0004 | Pending | This RDR may reference accessors but does not execute them. |
| CLI parse/lint output | RDR 0005 plus existing respond gateway | Pending | Failures must map to the CLI output contract. |
| Graph lint over normalized rows | RDR 0006 | Pending | This RDR must expose enough structure for determinism and reachability checks. |
| Expanded table dump | This RDR | Introduced | Reviewers can inspect the full table without maintaining it by hand. |
| Graph/render export | RDR 0006 | Pending | This RDR exposes normalized rows and source ids; graph-specific rendering remains with graph lint. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Structured CLI failures | `internal/cli/clierr::CLIError` | No table-specific codes yet | Extend | Add stable parse/lint refusal codes later. |
| Text/json output gateway | `internal/cli/respond::Fail` | Gateway is CLI-only, not kernel behavior | Reuse | Parser/lint commands must report through existing gateway. |
| Resolver/table package | `internal/` search | Not implemented | Introduce | New internal package can own sparse source structs and normalized row structs. |
| Config discovery | `internal/cli/config::Load` | Project config exists, table discovery not designed | Reuse later | Table path binding belongs with CLI integration, not this RDR's data format. |

### Decision Rationale

Sparse TOML is the best fit because the source is meant to be reviewed and
edited by humans, while the expanded table is a mechanical view for lint,
debugging, and documentation. A fully expanded table would make the RDR model's
dimensions multiply: status, profile, stage, prelock iteration, cluster
eligibility, rewind scope, and guards would force authors to copy the same
predicate fragments across rows. That is exactly the DX failure this RDR must
avoid.

The resource corpus points to the same split. `BUILD-SEEDS.md` says the table is
the flow's design and is hand-authored, but also says the resolver is table +
thin CLI + lint, not a runtime. `ANALYSIS-kernel-vocabulary.md` supplies the
vocabulary: recognized outcome is a Deferred Choice, guards are an Exclusive
Choice layer, unconditional rows are eventless/automatic transitions, and
ambiguous enabled edges should be lint-rejected rather than resolved by document
order. The Ragel POC shows why flat compiled topology is insufficient for the
RDR model: the hard parts are the coupled status register, cap-3 counter,
guards, and readable current state. Therefore the source should be sparse and
semantic, while normalization can render explicit rows for tools.

The direct `arc` corpus checks sharpen the sparse shape. Sismic demonstrates a
source model with nested states, transition guards, and contract classes that can
be exported to another view; that supports inherited contexts plus separate
precondition/invariant/postcondition diagnostics. `transitions` demonstrates
positive `conditions` and negative `unless` guard factoring; that supports
`all`/`unless` blocks instead of forcing every guard into one expression string.
Stateless demonstrates a symbolic graph object built from machine metadata; that
supports an expanded table/graph dump as a rendered view rather than the
authoring source. The statechart literature search supports hierarchy and
extended state as the known way to contain dimensional state explosion, while
the project still rejects adopting a statechart runtime.

Choosing exact-one candidate-row matching aligns with RDR 0001's deterministic
kernel and deliberately diverges from first-match FSM engines: priority order is
convenient in code, but it makes review harder and lets source or dump
reordering change behavior. Runtime FSM libraries are kept out of the core
because they would either drive orchestration or hide the contract in host
callbacks; their useful role is vocabulary, validation, and visualization.

Premortem: this could fail if the sparse source becomes a verbose
pseudo-language or if expansion hides surprising implicit rows. The
recommendation survives only if Resolve proves the RDR and kata examples with a
parse-normalize-dump spike. If that fails, the source schema must be reduced
before lock.

## Alternatives Considered

### Alternative 1: Sparse TOML Model With Expanded Table Dump

**Description**: Hand-authored TOML files declare tags, recognized outcomes,
shared contexts, accessor references, and sparse transition rules. The tool
normalizes rules into explicit candidate rows and can dump that expanded table.

**Pros**:

- Avoids Cartesian-product authoring for dimensional models such as RDR status
  plus profile plus prelock iteration plus guards.
- Keeps the source reviewable while still giving lint and resolver code a fully
  explicit row set.
- Supports mechanical dumps, diagnostics tied back to stable source rule ids,
  and future graph export by RDR 0006.
- Fits Go CLI implementation with ordinary typed decoding and validation.

**Cons**:

- Normalization is a real contract, not just parsing.
- Exact field layout needs spikes before lock.
- TOML is not a formal state-machine standard, so graph validation must be built
  over the parsed representation.

**Reason for selection**: Best balance of reviewability, parse simplicity, and
structured expressiveness for the target RDR and kata models.

### Alternative 2: Fully Expanded TOML Row Table

**Description**: Authors maintain one TOML row per candidate edge, with all
dimensions repeated inline.

**Pros**:

- Simplest parser and easiest mental model for tiny graphs.
- The source file is already the table the resolver sees.

**Cons**:

- Explodes for the RDR model once status, profile, prelock iteration, cluster
  gates, rewind scope, and guard dimensions interact.
- Repetition makes edits risky: changing one shared condition requires finding
  every copied row.
- The source becomes a generated-looking artifact even though humans are
  expected to own it.

**Reason for rejection**: It is acceptable as a dump format, not as the
authoring format.

### Alternative 3: JSON Model

**Description**: Store the same sparse model as JSON.

**Pros**:

- Simple to parse and strict about data types.
- Easy for tools to generate and consume.

**Cons**:

- Poor hand-review ergonomics: comments are unavailable, trailing comma churn is
  common, and nested objects become noisy for prose-heavy transition data.
- Encourages machine-generated artifacts, which conflicts with the goal that the
  legal graph be authored and reviewed directly.

**Reason for rejection**: It optimizes interchange over the human review loop
that is the central user outcome.

### Alternative 4: CSV / Matrix Table

**Description**: Store transitions as rows with columns for current state,
outcome, guard columns, and writes.

**Pros**:

- Compact and easy to scan for small state machines.
- Familiar representation for simple transition matrices.

**Cons**:

- Weak fit for nested predicates, typed values, accessor references, and
  multi-tag writes.
- Escaping and comments become awkward exactly where RDR/kata examples need
  explanation.

**Reason for rejection**: Too scalar for the required tag-set model.

### Alternative 5: Small DSL

**Description**: Define a custom text grammar such as `match -> writes` with
inline predicates.

**Pros**:

- Can be concise and domain-specific.
- Could make graph-like edges visually obvious.

**Cons**:

- Requires custom parsing, error recovery, formatting, editor support, and
  long-term grammar ownership.
- Pushes this RDR into language design before the target tables are proven.

**Reason for rejection**: The parser and tooling burden is not justified while
TOML can carry the same semantics.

### Briefly Rejected

- **Generated Go tables**: Fast and type-safe, but the reviewable artifact would
  be code, not the legal graph as data.
- **SCXML/XState as the source format**: Strong standard/tooling story, but a
  poor fit for tag provenance and the project's non-orchestrating resolver
  boundary.
- **Embedded host predicates**: Expressive, but defeats static lint and makes
  graph review depend on reading arbitrary code.

## Trade-offs

### Consequences

- The legal graph becomes a first-class reviewed model with stable rule ids.
- The expanded table becomes a generated diagnostic view, not hand-maintained
  source.
- Parser and lint errors become part of the CLI surface even though this RDR is
  primarily an internal data-format decision.
- Some expressiveness is intentionally deferred to RDR 0003 so this table stays
  statically checkable.

### Risks and Mitigations

- **Risk**: Sparse TOML rules become too verbose or too magical for large
  graphs.
  **Mitigation**: Resolve must encode representative RDR and kata fixtures and
  reject the shape if reviewers cannot scan the source or explain the dump.
- **Risk**: Sparse contexts hide an accidental Cartesian product.
  **Mitigation**: The dump must show normalized candidate rows with source rule
  ids, and lint must report expansion counts per rule.
- **Risk**: Future contributors treat row order as priority.
  **Mitigation**: Normative exact-one semantics and graph lint both reject
  overlapping rows instead of picking the first match.
- **Risk**: Accessor references pull execution semantics into the table.
  **Mitigation**: The table only names accessor bindings; RDR 0004 owns execution
  and read-back behavior.

### Failure Modes

Malformed TOML, unknown schema fields, unresolvable context references, or
normalization explosions fail at load/validation time with stable CLI errors. A
row gap, overlap, write to an undeclared tag, read-before-write condition, or
ambiguous expansion is a lint failure before the model is accepted. At runtime,
unknown outcomes, zero matches, multiple matches, and unavailable accessor
inputs are typed resolver refusals rather than guessed edges.

## Implementation Plan

### Prerequisites

- [x] All Critical Assumptions verified
- [ ] RDR 0001 remains aligned on exact-one stateless resolution.
- [ ] RDR 0003 confirms the fixed predicate operator set.

### Minimum Viable Validation

Parse two hand-authored sparse TOML fixtures, one for a representative RDR flow
slice and one for a representative kata flow slice, into typed source data;
normalize them into candidate rows; dump the expanded table; validate tag
declarations, context references, predicate references, recognized outcomes,
supported model version, and multi-tag writes; then prove by unit test that one
sample tag-set resolves to exactly one row, one unsupported-version variant is
refused before normalization, and one deliberately overlapping malformed variant
is refused as ambiguous. The RDR fixture must cover at least `Status`, `Profile`,
prelock iteration, and one rewind or cluster guard.

### Phase 1: Fixture and Schema Spike

Name the minimal sparse TOML field layout and encode representative RDR and kata
rules, including shared contexts, one self-loop, one rewind, one accessor
reference, one profile-dependent branch, and one multi-tag write.

### Phase 2: Normalizer and Dump

Introduce typed source structures, normalized candidate-row structures, and an
expanded-table dump with deterministic ordering and source rule ids.

### Phase 3: Parser and Validation Skeleton

Add validation rules for declarations, rule ids, context references, predicate
references, write targets, explicit clears, and expansion-count diagnostics.

### Phase 4: Resolver Handshake

Connect the normalized rows to RDR 0001's exact-one matching contract without
adding runtime ordering or host-code predicate callbacks. Internal indexes may
be decision trees or tries, but they must preserve the normalized semantics.

### Phase 5: Lint Handshake

Expose the parsed representation needed by RDR 0006 for graph determinism,
reachability, and read-before-write checks.

### Day 2 Operations

| Resource | List | Info | Delete | Verify | Backup |
| --- | --- | --- | --- | --- | --- |
| Transition model files | In scope via normal repository listing | In scope via parse/lint/dump output | N/A; source-controlled files | In scope via lint | N/A; source control is backup |
| Expanded table dumps | In scope through dump command | In scope through source rule ids | Delete/regenerate | In scope via dump tests | N/A; generated from source |

No runtime persistent resource is introduced by this RDR.

### New Dependencies

Use `github.com/pelletier/go-toml/v2` as the TOML parser candidate. The Resolve
spike ran against v2.3.1 from the local module cache, and the module license is
MIT. No production dependency is added until implementation.

## Validation

### Testing Strategy

Implementation tests must promote the Resolve spike into production fixtures:

1. **Scenario**: Parse the RDR and kata sparse TOML fixtures from `docs/rdr/0002-transition-table-as-reviewable-data/evidence/spikes/` into typed source structs.
   **Expected**: Tag declarations, root recognized-outcome alphabets, shared-context inheritance, accessor references, positive/negative guards, explicit clears, and multi-tag writes decode without ambiguous field placement.
2. **Scenario**: Normalize the RDR fixture's `continue-prelock` and `reconcile-rewind` rules and the kata fixture's `review-accepted` and `review-needs-work` rules.
   **Expected**: Candidate rows retain source rule ids/source locators, inherited predicates are expanded, `all`/`unless` predicates are visible in the row predicate set, and writes are deterministic.
3. **Scenario**: Validate malformed variants for unknown tags, unknown contexts, writes to non-owned tags, unknown accessors, unsupported versions, ambiguous overlaps, and missing root outcome alphabets.
   **Expected**: Each failure retains a stable data-level category and becomes a stable `CLIError` through the existing respond gateway when surfaced by CLI commands.
4. **Scenario**: Run exact-one selection over one matching tag-set and one deliberately overlapping/ambiguous negative fixture variant.
   **Expected**: The matching tag-set resolves to one row; zero or multiple matches are refusals and never fall back to row order.
5. **Scenario**: Normalize and dump two semantically identical fixtures whose TOML keys are authored in different orders.
   **Expected**: The expanded-table value is identical because rows sort by row identity and predicates/writes sort by key.

### Performance Expectations

Resolve evidence is functional rather than throughput-oriented. The spike
normalizes representative RDR and kata sparse fixtures into four deterministic
rows, and a repeated run produced byte-identical output with SHA-256
`3041e9e6678510e203a0213d5410df0852971416d64727c767345c4ca4725b24`.
Production code should preserve deterministic dump ordering by sorting stable
model/rule/predicate/write keys rather than relying on map iteration or source
order. The SHA is evidence for the spike output only; production golden tests
must assert the normalized expanded-table value defined by the normative
contract. Runtime lookup may index rows later, but that optimization must
preserve the normalized candidate-row semantics.

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

- RDR 0001, Resolution Kernel Contract.
- `docs/cli-output-contract.md`.
- Resource index: `.rdr/resources.md`.
- Seed prior: `../state-machines/BUILD-SEEDS.md`, especially Seed 2
  transition-table representation and Seed 3 guard expression.
- Transition model prior: "The transition model — inputs, outputs, error
  conditions."
- Kernel vocabulary prior: `../state-machines/ANALYSIS-kernel-vocabulary.md`,
  especially Deferred Choice, Exclusive Choice, eventless/automatic transitions,
  and lint-rejecting document-order ambiguity.
- Direct `arc` corpus checks:
  - `StateMachineOS`: `sismic/sismic/io/datadict.py::import_from_dict` and
    `export_to_dict` for nested source models and contracts.
  - `StateMachineOS`: `transitions/transitions/core.py::Transition` for
    positive `conditions` and negative `unless` guard lists.
  - `StateMachineOS`: `stateless/src/Stateless/Graph/StateGraph.cs::StateGraph`
    for symbolic graph generation from machine metadata.
  - `StateMachineOS`: `sismic/sismic/io/plantuml.py::PlantUMLExporter` for
    rendered graph output as a view over model data.
  - `StateMachineLit`: statechart hierarchy / extended-state literature hits as
    the prior-art answer to dimensional state explosion.
- RDR Ragel POC contrast: `../state-machines/contrast/poc-rdr-ragel/REVIEW.md`,
  especially the coupled status register, cap-3 counter, guard, and readable
  current-state limitations.
- RDR and kata flow audits from the state-machine prior-art corpus.
- Tool-fit assessment for FSM libraries as validation/visualization tools, not
  runtime orchestrators.
