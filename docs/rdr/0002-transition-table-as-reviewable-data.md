# Recommendation 0002: Transition Table As Reviewable Data

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-06-19
- **Status**: Draft
- **Type**: Architecture
- **Profile**: foundational — cross-RDR producer: the sparse TOML wire format,
  the normalization semantics that mint kernel rows (outcome binding,
  `RequiresOwned`, next-state tags and writes, existence constants, per-atom
  block), the expanded-table dump format and its total ordering, and the
  validation category taxonomy — consumed by RDRs 0003, 0006, 0007, 0008,
  and 0009.
- **Priority**: High
- **Related Issues**: None
- **Predecessors**: 0001-resolution-kernel; JDR 0001
  (`docs/jdr/0001-resolve-kernel-seam.md`) §D2, §D4, §JD-3, §JD-10
- **Overrides**: None
- **Seam Lineage**: no prior accretion

## Problem Statement

A flow author needs every legal edge to live in one reviewable artifact, so "is this transition legal" has one answer instead of being reconstructed from scattered prose. The system-internal requirement is to choose a table representation that can encode match predicates, multi-tag writes, and accessor references for both RDR and kata flows.

## Context

### Background

The transition table is the flow's design and must be hand-authored from the legal graph audits, not generated. The RDR table and kata table are the test workload, including multi-tag writes that a scalar state token cannot represent.

The real design fork is representation: TOML, JSON, CSV, or a small DSL, with trade-offs among diff reviewability, parse simplicity, expressiveness, and lintability.

### Technical Environment

intrastate is a Go CLI wired through `internal/cli`. The table format feeds the resolver kernel and the static lint, and must remain practical to parse, validate, and review in normal code review. The kernel this RDR produces rows for is implemented (`internal/resolve/resolve.go::Row`, `::Resolve`); it has no production consumer yet, so the normalizer is being specified against test-covered code rather than a shipped surface.

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

The search found no implemented transition table or normalizer package under
`internal/`; it found the kernel (`internal/resolve`), the existing CLI refusal
plumbing, and the peer RDR drafts. Prior-art reading favors a data table that
looks like "match predicates -> tag writes" and reserves external FSM libraries
for graph validation, not runtime orchestration.

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
- **Documented** — the kernel row (`internal/resolve/resolve.go::Row`) carries
  `RuleID`, `SourceLocator`, `Outcome`, `Match`, `RequiresOwned`, a guard,
  `Writes`, and `Escape`; `Resolve` refuses `unmodeled_outcome` before any row
  is consulted and filters candidates on `Row.Outcome`, so every normalized row
  must bind one outcome from the alphabet.
- **Verified** — Go's TOML tooling can preserve a sparse authoring schema's
  ergonomics for representative malformed-row diagnostics; the Resolve spike
  identified root-key placement as the exact field-layout constraint to lock.
- **Verified** — the RDR and kata legal graphs can be expressed as sparse rules
  that normalize to explicit row candidates, including explicit escape rows for
  modeled `no_match` behavior, without requiring host-code callbacks or an
  embedded expression language.
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
  - **Method**: Source Search
  - **Evidence**: The shipped kernel already implements gate-then-count, so
    the ordering is verified against code rather than merely decided.
    `internal/resolve/resolve.go::Resolve` runs
    `selected, blocked := gate(candidates, in.Guards, view)` and returns on
    `blocked != nil` **before** `switch len(selected)`;
    `::gate` prunes `GuardFalse`, then reports `missingOwned`, then vetoes on
    an undecidable survivor; `::escapeOrRefuse` delegates to the same `gate`
    before counting `viable`. `::missingOwned` sorts with `slices.Sort` and
    documents that "Row order is a normalization detail RDR 0002 owns and
    must not reach the reported diagnosis";
    `adversarial_test.go::TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder`
    freezes the property. The Normative Contracts below restate this flow and
    explicitly reject first-match semantics and count-first pruning
    (JDR 0001 §D2, whose "Lands in 0002" clause this discharges; RDR 0007
    states the same ordering as the kernel's and verifies it at its A20).
  - **If wrong**: Reviewers would have to reason about hidden priority, and
    reordering rows could silently change resolver behavior.
- **A3 RDR and kata flow edges can be encoded sparsely with fixed predicate
  operators rather than host-code callbacks or Cartesian-product row
  enumeration.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `cd docs/rdr/0002-transition-table-as-reviewable-data/evidence/spikes && GOCACHE=/private/tmp/intrastate-rdr0002-gocache GONOSUMDB='*' GOPROXY=off go run . rdr-fixture.toml kata-fixture.toml` parsed and normalized the RDR and kata fixtures. The fixtures cover status, profile, prelock iteration, equality, set membership, integer comparison, existence, self-loop, rewind, positive/negative guards, multi-tag writes, an `in` atom on `recognized` that expands into two candidate rows, and an explicit RDR `escape = ["no_match"]` row without host-code callbacks; `output.txt` captures the expanded candidate rows.
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
  - **Method**: Spike
  - **Evidence**: Normative Contracts define the expanded table as the
    normalized candidate-row value and fix the row order as the identity
    tuple `(model id, rule id, expansion suffix)` compared field by field,
    byte-lexicographically, with the source locator excluded; the same field
    list is what the Round-Trip invariant preserves. The normalizer spike
    witnesses it:
    `cd docs/rdr/0002-transition-table-as-reviewable-data/evidence/spikes &&
    GOFLAGS=-mod=mod GOPROXY=off go run . rdr-fixture.toml kata-fixture.toml`
    emits every field of the row value — identity with `#suffix` on expanded
    rows, locator, kind, lifted outcome, per-atom `@block`, next tags, writes
    including `<clear>`, derived `requires_owned`, escape classes — with
    atoms sorted by (key, block, operator, literal). Three consecutive runs
    produced byte-identical output despite Go's randomized map iteration
    (SHA-256 `c4be7447a241fb724632c53a9e5f39c7a9b5c7cc78e0a1e74ae4132271432a1b`,
    `evidence/spikes/output.txt`), and two semantically identical fixtures
    whose TOML keys are authored in different orders dumped byte-identically
    (`evidence/spikes/negative-cases.txt`). Byte-lexicographic comparison is
    a Go spec guarantee ("Two string values are compared lexically
    byte-wise"), so the ordering is locale- and platform-independent. This
    explicitly rejects source-order, map-iteration, and renderer-specific
    ordering, and rejects ordering on the source locator, which may carry
    line/column detail that changes when unrelated source text is edited.
  - **If wrong**: Golden tests and review dumps could churn across machines or
    refactors even when the transition semantics are unchanged.
- **A8 Every normalized row can bind exactly one outcome from the alphabet, so
  the kernel's `Row.Outcome` filter admits the row.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/resolve/resolve.go::Resolve` skips any row whose
    `Outcome` differs from `Input.Recognized` (`if row.Outcome !=
    in.Recognized { continue }`), before the match check;
    `::escapeOrRefuse` applies the same filter to escape rows; `::assemble`
    binds the outcome under `const recognizedTagKey = "recognized"` only when
    `in.Recognized != ""`; and `::Resolve` refuses `unmodeled_outcome` via
    `!in.Table.models(in.Recognized)` before any row is consulted. So a row
    carrying the lifted outcome is admitted, and with a non-empty
    `Input.Recognized` an outcome-less row never matches on either path.
    **Residual, and the reason this RDR is load-bearing for it:** the kernel
    validates no alphabet well-formedness — `::Table.models` is a plain
    `slices.Contains` — so with `Input.Recognized == ""` *and* the empty
    string in the alphabet, a row with an empty `Outcome` matches and emits a
    plan against a view carrying no `recognized` key. The recognized-totality
    clause below is the sole barrier; RDR 0008 observed the same path and
    left it "unforbidden rather than silently assumed away," and its
    accompanying "nothing in this RDR or RDR 0002 forbids that alphabet
    entry" predates that clause. JDR 0001 §JD-10 (whether a declared
    `recognized` tag is total when no outcome is in flight) is a different
    question and remains open.
  - **If wrong**: Normalized rows never match (every resolve refuses
    `no_match`), or the kernel needs an eventless-row concept this RDR does not
    model.

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
must provide enough structure for RDR 0001's exact-one resolver contract under
the gate-then-count ordering JDR 0001 §D2 fixes: guards are evaluated over every
candidate first, an unevaluable survivor vetoes the resolution, and only then is
exact-one counted; zero or multiple survivors are value-level refusals unless
the table also defines exactly one matching escape row for that failure class.
Unknown outcomes and unavailable accessors remain resolver refusals; malformed
TOML or schema-invalid rules are load/lint failures before the resolver sees a
table.

The model is data, not generated code and not a runtime FSM engine. RDR 0001's
resolver consumes the normalized representation, RDR 0003 owns the fixed
predicate operators and the tag type model, RDR 0004 owns accessor execution and
read-back safety, RDR 0005 exposes the CLI, RDR 0006 owns graph lint, RDR 0007
owns the guard seam and domain rule, RDR 0008 owns the reserved recognized tag
key, and RDR 0009 owns escape-row shape conformance. This RDR owns the on-disk
sparse representation, the normalizer that produces kernel rows from it, and the
normalized expanded-table view those peers consume.

### Technical Design

The source schema has six conceptual parts:

1. Flow metadata: table id, version, and human description.
2. Tag declarations: tag name, provenance (`owned`, `observed`, `recognized`),
   the RDR 0003 type model (value kind, and optionally finite domain,
   optionality, set-element universe, and single-valued marker), and optional
   accessor reference for observed or owned read-back. Provenance and authoring
   location are this RDR's; the type model is RDR 0003's and is cited, not
   restated. The recognized-provenance declaration is named `recognized` (RDR
   0008).
3. Recognized outcome alphabet: the closed set of outcome tags that the
   recognizer may emit for the flow.
4. Shared match contexts: named predicate blocks for dimensions such as RDR
   `Status`, `Profile`, stage, prelock iteration, cluster eligibility, or kata
   lifecycle. Contexts may inherit from another context to model statechart-like
   hierarchy without adopting a statechart runtime.
5. Sparse transition rules: reviewable rule ids, optional context references, a
   local match block, positive `all` guards, negative `unless` guards, and
   either a write block containing one or more tag assignments plus an optional
   explicit clear list, or an explicit `escape` list naming the resolver
   failure class the row models.
6. Render settings: deterministic ordering and field selection for the expanded
   table dump.

The parser turns TOML into typed source data, then normalizes it into explicit
candidate rows. Validation rejects malformed tags, unsupported model versions,
malformed predicate atoms (unknown operator, or a literal ill-formed for its
operator — including any existence literal other than the kernel's two boolean
forms), writes to non-owned tags, rules that match on missing tag declarations,
unresolvable context references, unknown accessor names, malformed escape
declarations, a rule that binds zero or more than one outcome, a
recognized-provenance declaration not named `recognized`, and ambiguous overlaps
between candidate rows.

Runtime matching is deliberately priority-free and follows the kernel's
gate-then-count ordering (JDR 0001 §D2; RDR 0007 states it as the kernel's):
the resolver collects the non-escape candidate rows for the recognized outcome
whose match pattern holds, evaluates every candidate's guard, prunes guard-FALSE
rows, reports absent owned state among survivors, refuses `guard_unevaluable`
if any survivor is undecidable, and only then counts. Exactly one survivor is
the plan. If zero or multiple survive, the resolver may return a modeled escape
disposition only when exactly one escape row for that failure class survives
the same gate; otherwise it returns the kernel-owned typed refusal.
`guard_unevaluable` and `owned_state_unavailable` are never escapable. Multi-tag
writes are first-class because RDR rewinds and kata lifecycle moves need to set
both the next stage/state and side-channel scope tags in one edge.

Guard factoring follows the `transitions` prior art: positive predicates and
negative predicates are authored separately and normalized into one predicate
set, in which each atom still records the block it was authored in. Contract
factoring follows the Sismic prior art: entry preconditions, postconditions,
and invariants are different validation classes, not free-form comments.
Rendered dumps follow the Stateless/Sismic export pattern: they are symbolic
views derived from model metadata and must carry enough source ids to send
diagnostics back to the authored sparse rule.

The normalized candidate row is the kernel row: it carries the source identity
(`(model id, rule id)` plus any expansion suffix, and the source locator), the
row kind (`transition` or `escape`), the single outcome the row responds to, the
predicate atoms (key, operator token, literal, block), the next-state tags and
the writes, both including rendered `<clear>` entries, the required-owned key
set derived from those writes, and the escape failure-class list. The
identity tuple `(model id, rule id, expansion suffix)` is what the dump sorts
on; the source locator is diagnostic and does not participate in that order.
The internal representation may be
indexed as a decision tree, trie, or decision DAG for efficient lookup, but that
is an implementation detail; the normative semantic object is the normalized
candidate-row set plus its source locator back to the sparse TOML rule. The
locator must identify at least the model id and rule id; byte line/column
coordinates are optional diagnostic detail.

#### Normative Contracts

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
live in a rule-level `clear` list; modeled escape rows live in a rule-level
`escape` list.
```

```normative
`[model]` MUST contain `id` and `version`. Version `1` is the only version this
RDR accepts; any other version MUST be refused before normalization.
```

```normative
Rule ids MUST be unique within a model, compared by exact byte equality; a
duplicate rule id is a load failure. Row identity, and therefore the dump's
total row ordering, rests on this.

Each transition rule MUST contain a stable rule id, zero or more shared-context
references, and a local match block. An ordinary transition rule MUST contain a
write block and MAY contain a rule-level explicit clear list. A write block MAY
assign more than one tag. An escape rule MUST contain an `escape` list and MUST
NOT contain a write block or clear list, even an empty one (RDR 0009 binds the
same obligation at the kernel boundary).
```

```normative
An `escape` list MUST contain only resolver failure classes that RDR 0001
allows the table to model: `no_match` and `ambiguous_match`. Normalization MUST
render an escape rule as a candidate row with row kind `escape`, its normal
predicate set, outcome, source rule id, source locator, and modeled failure
class list.
```

```normative
Shared contexts MAY inherit from other contexts, but inheritance MUST normalize
to an explicit predicate set before lint or resolution.
```

```normative
Guard predicates MUST be represented as positive `all` predicates and negative
`unless` predicates. Normalization MUST combine both into one candidate-row
predicate set before ambiguity checks, and each atom in that set MUST retain
the key, operator token, literal, and the block (`all` or `unless`) it was
authored in — the atom shape JDR 0001 §D1 fixes and RDR 0007 spells. Block
retention is carriage: RDR 0003's atom identity tuple and `unless` semantics
read it downstream, and normalization MUST NOT fold `unless` atoms into `all`.
```

```normative
For an existence atom the normalizer MUST emit the kernel's exported constants
verbatim — operator token `OpExists` and literal `LiteralTrue` or
`LiteralFalse` (RDR 0007, existence-operator clause; JDR 0001 §D4) — and MUST
reject at load, as a malformed predicate atom, any existence literal that is
not one of those two forms. The kernel treats a foreign token as a value atom
and a foreign literal as unevaluable; this load rejection is upstream of that
fail-closed backstop, not a substitute for it.
```

```normative
Tag-key identity is exact byte equality on the post-parse key string at every
stage — declaration lookup, context and rule predicate references, write and
clear targets, and the keys a normalized row carries. The normalizer performs
no case folding, trimming, or namespace rewriting (matching RDR 0008's
reserved-key comparison, which the same rule governs). The canonical spelling
of a tag is its `[tags.<tag>]` declaration key; every reference resolves to a
declaration by exact match or fails `unknown tag`, so a normalized row carries
only declared spellings and an atom's key agrees byte-for-byte with the key the
kernel assembles (RDR 0007 A22: the kernel canonicalizes nothing).
```

```normative
A tag declaration with provenance `recognized` MUST be named `recognized`, and
no owned or observed declaration may take that name; violations fail in RDR
0008's `reserved_tag_key` category, which participates in this RDR's data-level
category set. The `recognized` tag is total over matching: the kernel refuses
`unmodeled_outcome` before any row is consulted unless the resolve carries an
outcome in the declared alphabet, so every view that reaches a row binds
`recognized` (JDR 0001 §JD-10). Consequently the `outcomes` alphabet MUST be
non-empty, duplicate-free, and MUST NOT contain the empty string; a rule whose
predicate set requires `recognized` to be absent is dead, which is RDR 0006's
unreachable-rule finding, not a load failure here.
```

```normative
Every rule — ordinary or escape — MUST bind exactly one outcome: its combined
predicate set (local match block plus inherited contexts) MUST contain exactly
one atom on `recognized`, using `eq` or `in`, whose literal(s) are members of
the `outcomes` alphabet. Normalization lifts that atom out of the predicate
set into the row's outcome field; an `in` atom expands into one candidate row
per member, each identified by the rule id plus the outcome literal as its
expansion suffix. A rule binding zero outcomes, more than one `recognized`
atom, or a literal outside the alphabet is a load failure.
```

```normative
`Row.RequiresOwned` has no authored form. The normalizer MUST derive it for
every row as the sorted, duplicate-free set of tag keys named by the rule's
write block and clear list; by the write-to-non-owned-tag rule every such key
is an owned tag, and by RDR 0008 none is `recognized`. An escape row therefore
carries an empty set (RDR 0007 A21). What the field means is RDR 0007's
(post-guard write dependencies); this RDR is its producer (JDR 0001 §JD-3) and
does not add guard-read keys to it.
```

```normative
The kernel row carries the next state and the accessor-facing writes as two
distinct fields, and both have this RDR as their producer. Normalization MUST
populate the next-state tags with the tag values the rule's write block and
clear list produce, and the writes with the owned-tag writes the accessor
layer applies — the same rendered set, including `<clear>` entries. RDR 0009
A4 fixes that these are distinct fields and that only the writes reach the
accessor layer, so an escape row, which carries neither a write block nor a
clear list, normalizes to a row with both empty.
```

```normative
The tool MUST normalize the sparse source into deterministic candidate rows for
lint, resolver lookup, diagnostics, and table dumps. Each candidate row MUST
retain its source rule id and source locator.
```

```normative
The expanded table dump MUST be derived from the normalized candidate-row value
and MUST carry every field of it: row identity, source locator, row kind,
outcome, predicate atoms (with block), next-state tags, writes including
`<clear>` entries, required-owned keys, and escape failure classes.

Dump ordering MUST be deterministic across source key order. Rows sort by row
identity, compared field by field as the tuple `(model id, rule id, expansion
suffix)`, each field compared byte-lexicographically on the post-parse string;
an absent expansion suffix sorts before any present one. Within each row,
atoms sort by (key, block, operator token, literal), and next-state tags,
writes, required-owned keys, and escape classes sort by key.

The ordering is total: rule ids are unique within a model, and an expansion
suffix is a member of the `outcomes` alphabet, which is non-empty and
duplicate-free — so no two distinct rows compare equal and no positional
tiebreak is needed (RDR 0009 A8 establishes that a tiebreak on table position
is forbidden; RDR 0001's identity rule makes two orderings of one row set the
same input).

The **source locator MUST NOT participate in row ordering**. It may carry
optional line/column detail, so ordering on it would make the dump reorder
when unrelated source text is edited — the churn this contract exists to
prevent. This is a deliberate divergence from RDR 0007's refusal-payload
tuple, which orders on `(RuleID, SourceLocator, …)` because a refusal payload
is not a stable artifact under review.

The dump MUST be emitted from a pre-sorted sequence of normalized rows, never
by iterating a map and never by delegating key order to an encoder. Go's map
iteration order is deliberately randomized, and encoder guarantees vary:
`encoding/json` v1 sorts map keys but `encoding/json/v2` does not without an
explicit option, and TOML encoders may sort only within key groups.
```

```normative
Source order and rendered-row order MUST NOT decide a successful transition.
Selection is gate-then-count (JDR 0001 §D2): every candidate's guard is
evaluated first; guard-FALSE rows are pruned; absent owned state among
survivors refuses `owned_state_unavailable`; an unevaluable survivor refuses
`guard_unevaluable`; only then does exact-one counting run over the survivors.
Ordinary transition success requires exactly one surviving non-escape
normalized candidate row. If zero or multiple survive, the resolver MAY return
a modeled escape disposition only when exactly one escape row for that failure
class survives the same gate; otherwise zero or multiple matches remain
kernel-owned typed refusals. `guard_unevaluable` and `owned_state_unavailable`
are not escapable and MUST NOT be masked by a decidable escape row.
```

```normative
The model MUST declare every tag it matches or writes, including each tag's
provenance: owned, observed, or recognized.
```

```normative
A tag declaration also carries its **type model** — value kind, and optionally a
finite domain, an optionality marker, a set-element universe, and a
single-valued marker. RDR 0003 is the normative home of that model: what those
fields mean, which kinds admit a finite domain, and how domain/kind
disagreement is rejected are stated there and MUST NOT be restated here. This
RDR owns where a declaration is authored — under `[tags.<tag>]`, beside
`provenance` and the optional accessor reference — and requires normalization
to carry every declared field through to the normalized model without loss, so
lint (RDR 0006) and the guard proof (RDR 0003) read the same declaration the
author wrote.
```

```normative
Clearing a tag MUST be represented by an explicit rule-level `clear` entry that
normalization renders as a `<clear>` write. Absence from both the write block and
the clear list MUST NOT imply deletion.
```

```normative
Validation failures MUST retain stable data-level categories before CLI mapping,
including at minimum malformed TOML, unknown schema field, missing or malformed
recognized outcome alphabet, unknown tag, unknown context, write to non-owned
tag, unknown accessor, unsupported version, malformed predicate atom,
malformed escape declaration, malformed outcome binding, duplicate rule id,
`reserved_tag_key`
(RDR 0008), and ambiguous overlap.
```

#### Load-Bearing Decisions

- **Identity** — a transition rule is identified by `(model id, rule id)`.
  Normalized candidate rows inherit that identity plus a deterministic expansion
  suffix — the outcome literal when an `in` atom on `recognized` expands the
  rule. Rule ids are stable review anchors and must not be reused for a
  different edge. The identity tuple `(model id, rule id, expansion suffix)`
  is also the dump's row-sort key, and it is total; the source locator is
  diagnostic and deliberately excluded from it.
- **Wire / byte format** — TOML is the on-disk carrier. The exact field names
  are the Resolve spike layout: root `outcomes`, `[model]`, `[tags.<tag>]`,
  `[accessors.<id>]`, `[context.<id>]`, `[[rule]]`, `[rule.write]`,
  rule-level `clear`, rule-level `escape`, and `[dump]`. `[model].version = 1`
  is the only accepted format version. The RDR and kata spike fixtures are the
  canonical examples implementation tests must promote.
- **Naming** — the canonical source artifact name is "transition model"; the
  canonical rendered view is "expanded transition table." Tag keys are exact
  byte strings with no folding, and the recognized-provenance tag is the
  reserved name `recognized`. Rejected names: "state machine config" because
  it suggests a runtime driver, and "workflow graph" because this RDR owns
  sparse transition data, not orchestration. Rejected: case-insensitive or
  trimmed key matching, because it would let two authored spellings reach the
  kernel as distinct keys and would make `Recognized` collide with the reserved
  key RDR 0008 compares byte-exactly.
- **Selection / predicate** — the only successful transition selection is
  exact-one surviving non-escape row after the kernel's gate (guard-FALSE
  pruned, owned-state and unevaluable refusals first). Count-first pruning is
  rejected because it launders missing state into an escapable `no_match`
  (JDR 0001 §D2). Lint rejects ambiguous ordinary overlaps and ambiguous escape
  overlaps before runtime.
- **Outcome binding** — a rule's outcome is the single `recognized` atom in its
  combined predicate set, lifted into the row's outcome field. Rejected: an
  eventless row with no outcome, because the kernel has no such concept and
  would never match it; rejected: leaving the atom in the predicate set only,
  because the kernel filters on `Row.Outcome` before matching.

#### Round-Trip / Inverse Invariants

`parse ∘ normalize ∘ dump = expanded-table value identity` on valid transition
model fixtures: dumping the normalized model and reading the dump as a table
view must preserve the candidate-row set — the same field list the dump
contract carries: row identity, source locator, row kind, outcome, predicate
atoms with block, next-state tags, writes, required-owned keys, and escape
failure classes.
Source rewrite is out of scope for this RDR; a later rewrite-capability RDR
must define its own source-preservation invariant before mutating authored TOML.

#### Illustrative Code

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
recognized.eq = "verdict-flapping"

[rule.guard.all]
iter.lt = 3

[rule.guard.unless]
profile.eq = "small"

[rule.write]
stage = "prelock"

[[rule]]
id = "draft-no-match-escape"
use = ["draft"]
escape = ["no_match"]

[rule.match]
recognized.eq = "round-clean"
```

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Stateless exact-one resolution | RDR 0001 | Implemented | This RDR must provide normalized candidate rows, outcome binding, and tag writes the kernel can evaluate. |
| Fixed predicate operators and tag type model | RDR 0003 | Pending | This RDR names predicate slots and carries declarations but does not own the operator grammar or declaration semantics. |
| Accessor references and safe read-back | RDR 0004 | Pending | This RDR may reference accessors but does not execute them. |
| CLI parse/lint output | RDR 0005 plus existing respond gateway | Pending | Failures must map to the CLI output contract. |
| Graph lint over normalized rows | RDR 0006 | Pending | This RDR must expose enough structure for determinism and reachability checks. |
| Guard seam, atom shape, existence constants, gate ordering | RDR 0007 | Final, kernel reshape unimplemented | Normalizer emits per-atom block and the kernel's `OpExists`/`LiteralTrue`/`LiteralFalse`; derives `RequiresOwned`. The constants and the atom-shaped row do not exist in `internal/resolve` yet, so Phases 2–3 sequence behind that reshape (Prerequisites). |
| Reserved recognized tag key | RDR 0008 | Final | Declaration named `recognized`; `reserved_tag_key` joins this RDR's category set. |
| Escape-row shape conformance | RDR 0009 | Final | Escape rules carry no write block or clear list; normalized escape rows render write-free. |
| Expanded table dump | This RDR | Introduced | Reviewers can inspect the full table without maintaining it by hand. |
| Graph/render export | RDR 0006 | Pending | This RDR exposes normalized rows and source ids; graph-specific rendering remains with graph lint. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Structured CLI failures | `internal/cli/clierr::CLIError` | No table-specific codes yet | Extend | Add stable parse/lint refusal codes later. |
| Text/json output gateway | `internal/cli/respond::Fail` | Gateway is CLI-only, not kernel behavior | Reuse | Parser/lint commands must report through existing gateway. |
| Kernel row and resolver | `internal/resolve::Row`, `::Resolve` | Consumes rows; has no loader or normalizer | Reuse | The normalizer targets this row shape; no new kernel surface. |
| Normalizer/table package | `internal/` search | Not implemented | Introduce | New internal package can own sparse source structs and normalization. |
| Config discovery | `internal/cli/config::Load` | Project config exists, table discovery not designed | Reuse later | Table path binding belongs with CLI integration, not this RDR's data format. |

### Decision Rationale

Joint-check: fired → 0007 (home: JDR 0001 §D4 / §JD-12) — the kernel enforces the guard domain; the normalizer obligations §D4 places on this RDR (existence constants, exact key identity, per-atom block) are Normative Contracts above, citing RDR 0007 rather than restating it.

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
Choice layer, and ambiguous enabled edges should be lint-rejected rather than
resolved by document order. The Ragel POC shows why flat compiled topology is
insufficient for the RDR model: the hard parts are the coupled status register,
cap-3 counter, guards, and readable current state. Therefore the source should
be sparse and semantic, while normalization can render explicit rows for tools.

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
reordering change behavior. Gate-then-count is the kernel's ordering and the
only one under which missing artifact state cannot hide behind an escapable
refusal. Runtime FSM libraries are kept out of the core because they would
either drive orchestration or hide the contract in host callbacks; their useful
role is vocabulary, validation, and visualization.

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
- Every rule must name its outcome; the model has no eventless rows because the
  kernel has none.

### Risks and Mitigations

- **Risk**: Sparse TOML rules become too verbose or too magical for large
  graphs.
  **Mitigation**: Resolve must encode representative RDR and kata fixtures and
  reject the shape if reviewers cannot scan the source or explain the dump.
- **Risk**: Sparse contexts hide an accidental Cartesian product.
  **Mitigation**: The dump must show normalized candidate rows with source rule
  ids, and lint must report expansion counts per rule.
- **Risk**: Future contributors treat row order as priority.
  **Mitigation**: Normative exact-one semantics require explicit escape rows for
  modeled no-match or ambiguous-match behavior; graph lint rejects overlapping
  ordinary or escape rows instead of picking the first match.
- **Risk**: Accessor references pull execution semantics into the table.
  **Mitigation**: The table only names accessor bindings; RDR 0004 owns execution
  and read-back behavior.
- **Risk**: Normalizer output drifts from the kernel's exported constants or
  key identity, so a valid-looking atom fails closed at runtime.
  **Mitigation**: The MVV asserts emitted existence atoms byte-for-byte against
  `resolve.OpExists`/`LiteralTrue`/`LiteralFalse` and asserts row keys equal
  declaration keys.
- **Risk**: Dump order tracks byte positions in the authored source, so adding
  a comment or reordering unrelated rules churns every downstream golden test.
  **Mitigation**: The source locator is excluded from the row-sort key, and
  the identity tuple that replaces it carries no positional component. The MVV
  asserts the comparator is total, so a later identity change cannot
  reintroduce a positional tiebreak unnoticed.

### Failure Modes

Malformed TOML, unknown schema fields, unresolvable context references,
malformed predicate atoms, malformed escape declarations, malformed outcome
bindings, a misnamed recognized declaration, or normalization explosions fail at
load/validation time with stable CLI errors. A row gap, overlap, write to an
undeclared tag, read-before-write condition, dead row, or ambiguous expansion is
a lint failure before the model is accepted. At runtime, unknown outcomes,
unavailable accessor inputs, absent owned state, undecidable guards, and
unmodeled zero/multiple survivors are typed resolver refusals rather than
guessed edges; modeled zero/multiple-match behavior must come from exactly one
surviving escape row, and the gate refusals are never modeled.

## Implementation Plan

### Prerequisites

- [x] All Critical Assumptions verified (A2, A7 re-verified; A8 resolved).
- [ ] RDR 0007 is the normative home of the atom shape and the existence
  constants this RDR's normalizer emits. RDR 0007 is `Final`, so the
  *specification* is settled — but its kernel reshape is unimplemented: the
  shipped kernel still carries `internal/resolve/resolve.go::Row.Guard` as a
  string and exports no `OpExists` / `LiteralTrue` / `LiteralFalse`. Phases 2
  and 3 cannot emit atom-shaped rows against that surface until the reshape
  lands, so this item gates implementation sequencing, not lock. It does not
  invert the cluster's ordering: RDR 0008 and RDR 0009 both name this RDR's
  normalizer as the enforcement point their own checks land inside, and this
  RDR's normalized rows target RDR 0001's shipped kernel row.
- [x] RDR 0001 remains aligned on exact-one stateless resolution.
- [x] RDR 0003 is coherent enough for this RDR to defer the fixed predicate
  operator set and tag type model without importing that grammar.
- [x] RDR 0007, 0008, 0009 are Final; their producer obligations on this RDR are
  bound above.

### Minimum Viable Validation

Parse two hand-authored sparse TOML fixtures, one for a representative RDR flow
slice and one for a representative kata flow slice, into typed source data;
normalize them into candidate rows; dump the expanded table; validate tag
declarations, context references, predicate references, recognized outcomes,
supported model version, multi-tag writes, outcome binding, and escape rows;
then prove by unit test that one sample tag-set resolves through
`internal/resolve::Resolve` to exactly one ordinary row, one unmatched sample
tag-set resolves to exactly one modeled escape row when the table declares one,
one unsupported-version variant is refused before normalization, and one
deliberately overlapping malformed variant is refused as ambiguous. The RDR
fixture must cover at least `Status`, `Profile`, prelock iteration, one rewind
or cluster guard, one existence atom, one `in` atom on `recognized` that
expands into more than one row, and one explicit `no_match` escape row.

The dump ordering is MVV-covered, not left to implementation: assert that the
expanded-table value is unchanged when the same model is authored with its
TOML keys in a different order, that repeated dumps of one model are
byte-identical, and that the row-sort comparator is total — sort, then assert
no adjacent row pair compares equal, so a future identity change cannot
silently reintroduce a positional tiebreak.

### Phase 1: Fixture and Schema Spike

Name the minimal sparse TOML field layout and encode representative RDR and kata
rules, including shared contexts, one self-loop, one rewind, one accessor
reference, one profile-dependent branch, one explicit escape row, and one
multi-tag write.

### Phase 2: Normalizer and Dump

Introduce typed source structures, normalization to `internal/resolve::Row`
values (outcome lifted, per-atom block retained, `RequiresOwned` derived,
existence constants emitted), and an expanded-table dump with deterministic
ordering and source rule ids.

### Phase 3: Parser and Validation Skeleton

Add validation rules for declarations, the reserved `recognized` name, rule
ids, context references, predicate atoms, outcome bindings, write targets,
explicit clears, escape declarations, and expansion-count diagnostics.

### Phase 4: Resolver Handshake

Connect the normalized rows to RDR 0001's gate-then-count exact-one contract
without adding runtime ordering or host-code predicate callbacks. Internal
indexes may be decision trees or tries, but they must preserve the normalized
semantics.

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
   **Expected**: Tag declarations (including `[tags.recognized]`), root recognized-outcome alphabets, shared-context inheritance, accessor references, positive/negative guards, explicit clears, escape declarations, and multi-tag writes decode without ambiguous field placement.
2. **Scenario**: Normalize the RDR fixture's `continue-prelock`, `reconcile-rewind`, `terminal-archive`, and `draft-no-match-escape` rules and the kata fixture's `review-accepted` and `review-needs-work` rules.
   **Expected**: Candidate rows retain source rule ids/source locators, inherited predicates are expanded, each atom reports its authored block (`match`/`all`/`unless`), the single `recognized` atom is lifted into the row's outcome field and absent from the predicate set, `RequiresOwned` equals the sorted write-plus-clear key set (empty on the escape row), escape rows retain row kind plus modeled failure class list and carry neither writes nor next-state tags, and writes are deterministic. `terminal-archive`'s `in` atom expands to two rows carrying the outcome literal as their expansion suffix. The normative fixture for this scenario is `evidence/spikes/output.txt`.
3. **Scenario**: Validate malformed variants for unknown tags, unknown contexts, writes to non-owned tags, unknown accessors, unsupported versions, malformed predicate atoms (unknown operator; `exists` with a non-boolean literal), malformed escape declarations (including an empty write block on an escape rule), a rule with zero or two `recognized` atoms, an outcome literal outside the alphabet, an alphabet containing the empty string, an owned declaration named `recognized` and a recognized declaration named `outcome`, a duplicate rule id, ambiguous overlaps, and missing root outcome alphabets.
   **Expected**: Each failure retains a stable data-level category and becomes a stable `CLIError` through the existing respond gateway when surfaced by CLI commands. Nine of these are already witnessed by the spike, one refusal per mutated fixture, in `evidence/spikes/negative-cases.txt`: unsupported version, empty-string alphabet member, duplicate rule id, unknown written tag, write to an observed tag, escape rule carrying a write block, `exists` with a non-boolean literal, an outcome literal outside the alphabet, and a recognized declaration misnamed `outcome`.
4. **Scenario**: Run `internal/resolve::Resolve` over one matching ordinary tag-set in which every sibling candidate's guard is decidable, one tag-set with no ordinary match but one matching `no_match` escape row, one tag-set in which a sibling candidate's guard is unevaluable, and one deliberately overlapping/ambiguous negative fixture variant.
   **Expected**: The matching ordinary tag-set resolves to one transition row; the no-match tag-set resolves to the modeled escape disposition; the unevaluable-sibling tag-set refuses `guard_unevaluable` even though a decidable sibling and a `no_match` escape row exist; zero or multiple survivors without exactly one surviving escape row are refusals and never fall back to row order.
5. **Scenario**: Normalize and dump two semantically identical fixtures whose TOML keys are authored in different orders, and dump the same fixture repeatedly in one process and across processes.
   **Expected**: The expanded-table value is identical because rows sort by the total identity tuple and atoms/next tags/writes/required-owned keys/escape classes sort by key; repeated dumps are byte-identical despite Go's randomized map iteration. Witnessed in `evidence/spikes/negative-cases.txt` (reordered kata fixture, matching SHA-256) and by three consecutive byte-identical runs recorded under Performance Expectations.
6. **Scenario**: Normalize a rule carrying `exists = true` and one carrying `exists = false`, one carrying a non-boolean existence literal, and a variant spelling a declared tag `Status` where the declaration is `status`.
   **Expected**: Emitted atoms carry `resolve.OpExists` and `resolve.LiteralTrue`/`LiteralFalse` byte-for-byte; the non-boolean literal is refused at load as a malformed predicate atom; the mis-cased reference fails `unknown tag` rather than folding. The RDR fixture's `continue-prelock` guard carries the `finalized_at.exists=false@all` atom the spike emits.

### Performance Expectations

Resolve evidence is functional rather than throughput-oriented. The spike
normalizes representative RDR and kata sparse fixtures into seven deterministic
rows — one explicit escape row and one rule that expands into two rows through
an `in` atom on `recognized`. Three consecutive runs produced byte-identical
output with SHA-256
`c4be7447a241fb724632c53a9e5f39c7a9b5c7cc78e0a1e74ae4132271432a1b`, and a
fixture with the same semantics authored in a different key order produced a
matching digest.

Determinism checklist, as run: source key order is neutralized by sorting every
emitted sequence; Go's map iteration is never the emission order (rows, atoms,
next tags, writes, required-owned keys, and escape classes are all sorted
slices before rendering); string comparison is byte-lexicographic per the Go
spec, so ordering is locale- and platform-independent; set-valued literals are
rendered with their members sorted so two authored orderings of one set are one
literal; no case folding, trimming, or namespace rewriting is applied at any
stage; the expansion suffix is absent or an alphabet member, never the empty
string, so the identity tuple is total. No hash is part of the contract: the
SHA is evidence for this spike output only, and production golden tests must
assert the normalized expanded-table value the normative contract defines.
Runtime lookup may index rows later, but that optimization must preserve the
normalized candidate-row semantics.

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

### Proportionality

[Is the document right-sized for the change? Flag
any sections that should be trimmed before locking.
The split test is **contract count, not word count**:
confirm this RDR is the sole author of at most one
independent load-bearing contract. Re-validate the
**Profile** Metadata field against the contracts
counted.]

## References

- RDR 0001, Resolution Kernel Contract; `internal/resolve/resolve.go::Row`,
  `::Resolve`, `::assemble`.
- JDR 0001, `docs/jdr/0001-resolve-kernel-seam.md` — §D2 (gate-then-count),
  §D4 (kernel-enforced guard domain; normalizer emits existence constants),
  §JD-3 (`RequiresOwned` producer), §JD-10 (reserved-key fixture rename;
  recognized-tag totality).
- RDR 0007, Guard Predicate Totality — atom shape, `OpExists`/`LiteralTrue`/
  `LiteralFalse`, A21/A22, gate ordering.
- RDR 0008, Recognized Tag Key Ownership — reserved `recognized` name and
  `reserved_tag_key` category.
- RDR 0009, Escape-Row Shape Conformance — write-free escape rows.
- `docs/cli-output-contract.md`.
- Resource index: `.rdr/resources.md`.
- Seed prior: `../state-machines/BUILD-SEEDS.md`, especially Seed 2
  transition-table representation and Seed 3 guard expression.
- Transition model prior: "The transition model — inputs, outputs, error
  conditions."
- Kernel vocabulary prior: `../state-machines/ANALYSIS-kernel-vocabulary.md`,
  especially Deferred Choice, Exclusive Choice, and lint-rejecting
  document-order ambiguity.
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
