# Recommendation 0002: Transition Table As Reviewable Data

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-06-19
- **Status**: Final [joint decision → JDR 0001 §JD-15, §JD-16, §JD-17: `<clear>` at the write accessor; Match/Guard routing key; wire keys for initial/terminal, write-replaces, accessor metadata, type-model fields]
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
  - **Evidence**: `internal/cli/clierr/clierr.go::CLIError` already carries stable `Code`, `Message`, optional `Param`, diagnostic `Detail`, `Hint`, exit-code `Group`, and `Cause`; `internal/cli/respond/respond.go::Fail` emits that envelope in text/json modes. The envelope is the verified half. **No parse-validation code has yet travelled it:** `internal/cli/config/config.go::Load` emits only `config-not-found` and `config-read-error`, and its `config-invalid` exists solely as a doc comment and a `TODO` — `Load` reads bytes, stamps `SchemaVersion` unconditionally without reading the file's own version, and returns a nil error. The codes emitted anywhere in the CLI are exhaustively `config-not-found`, `config-read-error`, `flag-invalid-value`, and `command-error`. This RDR's table parse and lint categories would be the **first** parse-validation codes in the envelope, so their code names, `Group` assignments, and exit mapping are new design work at implementation, not a pattern to copy.
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

    **The barrier's reach is bounded, and this is a known residual.** The
    alphabet ban is a *load-time* check in the normalizer; the kernel does not
    import the normalizer and validates no alphabet well-formedness itself, so
    the barrier covers exactly one producer — tables this RDR's loader built.
    Any other construction of a `resolve.Table` (a test helper, a future
    programmatic producer, a caller assembling rows directly) reaches the same
    unguarded path, and so does a caller passing `Input.Recognized == ""` from
    the CLI edge, which is RDR 0005's boundary and not reachable from here.
    Closing it for all producers means a kernel-side check, which is RDR 0001's
    to make. This RDR closes what it produces and states the rest as residual
    rather than implying table-wide coverage it cannot enforce.
  - **If wrong**: Normalized rows never match (every resolve refuses
    `no_match`), or the kernel needs an eventless-row concept this RDR does not
    model.
- **A9 The normalized candidate row can carry one unified, block-retaining atom
  set and still produce a conforming `resolve.Row` at the kernel handoff.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: Owed. The repeatability lens found the RDR placed the
    operator-based `Match`/`Guard` split on two different artifacts; it is now
    pinned to the handoff so the normalized value keeps each atom's authored
    block, which `Guard string` cannot carry. Verification: the Testing Strategy
    scenario 2 normalization assertions plus scenario 4's `Resolve` run must
    both pass against one normalized value — the first reading the unified atom
    set, the second consuming the split row. Blocked on RDR 0007's kernel
    reshape (Prerequisites), like the rest of Phase 2.
  - **If wrong**: Either block retention is lost (breaking RDR 0003 downstream
    and the Round-Trip field list) or the kernel handoff needs a shape this RDR
    does not model.
- **A10 `duplicate model id` is decidable only across documents, so no
  single-document loader can or should report it.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: The decidability half is verified at source: the Resolve
    spike's `evidence/spikes/main.go::Model` declares `Model` as a **singular
    struct field**, not a slice, mirroring `[model]` as a singular TOML table
    (`evidence/spikes/rdr-fixture.toml` authors one `id = "rdr"` per document),
    so a single-document load structurally cannot observe two model ids. The
    category is therefore scoped to a multi-document dump/lint invocation. The
    behavioral half is owed: the scenario 3 fixture for this category is a
    **pair** of documents sharing a `model id`, and a single-document load of
    either one must not report it.
  - **If wrong**: The category is unreachable as specified, or a multi-document
    surface exists that this RDR has not named.
- **A11 An escape rule binding its outcome with an `in` atom expands into one
  rescuing row per member, each rescuing only its own outcome.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: Owed. The critique lens introduced the clause and booked it
    unwitnessed; the spike fixture's escape rule
    (`evidence/spikes/rdr-fixture.toml` `draft-no-match-escape`) binds with
    `eq = "round-clean"`, so no captured run exercises escape expansion. The
    kernel half is verified — `internal/resolve/resolve.go::escapeOrRefuse`
    filters escape candidates on `row.Outcome != in.Recognized` exactly as
    `::Resolve` does, so a per-outcome escape row is what the kernel admits.
    Verification: extend the scenario 2 fixture with an escape rule binding
    `in` over two alphabet members and assert it expands to two escape rows
    carrying distinct expansion suffixes, each with an empty write set.
  - **If wrong**: One authored escape rule cannot cover several outcomes, so
    covering an alphabet of N outcomes requires N hand-authored escape rules
    and the "expand like any other rule" clause overstates the format.
- **A12 The kernel handoff routes every atom to exactly one of `Match` /
  `Guard`, exhaustively and disjointly.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: Owed. A9 asserts the unified atom set *can* produce a
    conforming `resolve.Row`; it does not assert the routing is total. The
    critique lens booked this separately and it landed in no assertion.
    Verification: a scenario 2 assertion that for the RDR fixture the union of
    the constructed row's `Match` and guard atoms equals the normalized atom
    set and the intersection is empty — including the `unless`-block equality
    atom (`status.eq=closed@unless` in the kata fixture), which routes to the
    guard despite being an equality operator. Blocked on RDR 0007's kernel
    reshape (Prerequisites), like A9.
  - **If wrong**: An atom is silently dropped at the handoff or double-counted,
    which changes which rows match without any dump or normalization test
    observing it.
- **A13 Merging is idempotent on identical atoms: one atom contributed by a
  rule and by one or more inherited contexts collapses to exactly one.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: Owed. Testing Strategy scenario 2 asserts the *differing*-
    literal half (two atoms, both `use` orders); the mirror half landed in no
    assertion, so a normalizer emitting duplicates passes every stated control.
    Verification: a scenario 2 assertion that a rule and an inherited context
    contributing a byte-identical `(block, key, operator, literal)` atom yield
    an atom set of count one for that key.
  - **If wrong**: The normalized atom set carries duplicates, inflating the
    dump and the Round-Trip field list without changing match semantics — or
    a de-duplication keyed too widely drops the A13-adjacent distinct-literal
    atoms that scenario 2 requires to survive.
- **A14 The version gate precedes strict field decoding, so a v2 document
  refuses as an unsupported version rather than an unknown schema field.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: Owed. The clause is normative and the only ordering this RDR
    fixes, but scenario 3's `unsupported version` fixture is a v1-shaped
    document with a bad version value, which trips the category on either
    ordering and therefore cannot witness the precedence. Verification: a
    scenario 3 fixture that is **v2-shaped** — `version = 2` plus a v2-only key
    the strict decoder would reject — asserted to refuse `unsupported version`
    and not `unknown schema field`.
  - **If wrong**: A future v2 file fails with a diagnostic naming an arbitrary
    field and never mentioning the version, which is the confusion the two-pass
    gate exists to prevent.

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
candidate rows. Load-time validation rejects malformed tags, unsupported model
versions, malformed predicate atoms (unknown operator, or a literal ill-formed
for its operator — including any existence literal other than the kernel's two
boolean forms), writes to non-owned tags, rules that match on missing tag
declarations, unresolvable or cyclic context references, unknown accessor names,
malformed escape declarations, a rule that binds zero or more than one outcome,
and a recognized-provenance declaration not named `recognized`.

Load and lint split on **arity, not severity**: a check decidable from one rule
plus the model's declarations is load-time and this RDR's; a check that must
compare normalized rows against each other — overlap, gap, dead row,
read-before-write — is graph lint and RDR 0006's. Ambiguous overlap between
candidate rows is therefore a lint finding, not a load failure.

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

The normalized candidate row is the kernel row up to the predicate split — the
kernel's `Match`/`Guard` routing is applied at the handoff (Normative
Contracts), not stored in the normalized value. It carries the source identity
(`(model id, rule id)` plus any expansion suffix, and the source locator), the
single outcome the row responds to, the
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
`escape` list. An `[accessors.<id>]` entry carries `mode` and `path`, as the
spike fixtures author it; this RDR validates only that a referenced accessor id
resolves to a declared entry (`unknown accessor` otherwise) — RDR 0004 owns the
entry's semantics and execution.
```

```normative
`[model]` MUST contain `id` and `version`. Version `1` is the only version this
RDR accepts; any other version MUST be refused before normalization.

**The version check MUST run before strict field validation, not merely before
normalization.** Loading MUST therefore proceed in two passes: read `[model]`
permissively enough to obtain `version`, refuse on any value but `1`, and only
then decode the document strictly. Ordering these the other way makes a future
v2 file fail as `unknown schema field` on whichever v2-only key the decoder
reaches first — a diagnostic that names an arbitrary field and never mentions the
version, which is precisely the confusion the version gate exists to prevent.

**Load is fail-fast: the first category a document trips is the refusal.**
Beyond the version gate above, the order in which independent defects are
checked is deliberately **unspecified** — an implementation MAY check in any
order, and a document tripping two categories MAY be refused with either. Only
the version gate's precedence is normative. This is stated so the freedom is
visible rather than accidental: the Testing Strategy's one-defect-per-fixture
rule makes order unobservable by construction, so a later implementation MUST
NOT be read as bound to whichever order the first one happened to use, and an
accumulating loader that returns a list is a different contract than this one.

**Strict decoding is an obligation on this format, not a property of a library.**
The decoder MUST reject unmapped keys so an unknown schema field is a stable
refusal rather than a silent no-op. This is stated normatively because it is not
the default behavior of TOML decoders generally, and a parser swap that silently
lost it would retire the `unknown schema field` category without any contract
appearing to change.
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
render an escape rule as a candidate row carrying its normal predicate set,
outcome, source rule id, source locator, and modeled failure class list.

**Escape rescue is per-outcome, and an escape rule binds its outcome by the same
rule as any other.** An escape row rescues only resolves carrying the outcome it
binds: `internal/resolve/resolve.go::escapeOrRefuse` filters escape candidates
with `row.Outcome != in.Recognized` before gating them, exactly as `::Resolve`
filters ordinary rows. There is no table-wide catch-all: covering an alphabet of
N outcomes requires N escape rows. An escape rule MAY bind its outcome with an
`in` atom and expand like any other rule, which is how one authored rule covers
several outcomes — each expansion is a separate row rescuing its own outcome, and
the expansion suffix distinguishes them. This is a consequence of the shipped
kernel, not a choice this RDR makes; it is stated here because the escape clause
is where an author looks, and an author who expects one unbound escape row to
backstop the whole table gets silent unrescued refusals on every other outcome.

Row kind (`transition` / `escape`) is a **derived view property, never a row
field**: escape identity is discriminated solely by a non-empty escape class
list, which RDR 0009 fixes as the single discriminator and which forbids a
row-kind field at the kernel boundary. The dump renders the kind as a column
computed from that list; normalization introduces no field to carry it. The
prohibition binds the **normalized value and the kernel row**; a render-time
view struct MAY materialize the computed column, since the dump contract
requires the column. What it MUST NOT do is let that view feed back into the
normalized value or the kernel row — the discriminator stays the escape class
list, and no code may branch on a stored kind.
```

```normative
Shared contexts MAY inherit from other contexts, but inheritance MUST normalize
to an explicit predicate set before lint or resolution.

The combined predicate set is a **set over the full atom identity** — the tuple
`(key, block, operator token, literal)`, the same tuple the dump sorts on.
Merging MUST NOT key on any proper prefix of it: two atoms agreeing on
`(key, block, operator)` but differing in literal are **distinct atoms** and both
survive the merge. Keying the merge on `(block, key, operator)` alone silently
drops one constraint and makes the survivor depend on map iteration order over
the `use` list, which the dump's sort then certifies as deterministic — the loss
happens before the sort can see it, so ordering determinism cannot detect it.

Merging is idempotent on identical atoms: the same atom contributed by a rule and
by one or more inherited contexts collapses to one. Inheritance therefore never
overrides — it only accumulates. A narrowing "override" of an inherited
constraint is not expressible, and is not silently approximated: authoring two
atoms on one key that no view can satisfy together yields a dead rule, which is
RDR 0006's unreachable-rule finding, not a load failure here.
```

```normative
Guard predicates MUST be represented as positive `all` predicates and negative
`unless` predicates. Normalization MUST combine both into one candidate-row
predicate set before ambiguity checks, and each atom in that set MUST retain
the key, operator token, literal, and the block (`all` or `unless`) it was
authored in — the atom shape JDR 0001 §D1 fixes and RDR 0007 spells. Block
retention is carriage: RDR 0003's atom identity tuple and `unless` semantics
read it downstream, and normalization MUST NOT fold `unless` atoms into `all`.

Because the block is part of the atom identity, one atom authored in **both**
`all` and `unless` is two distinct atoms and both survive normalization; this is
not a load failure. The rule is self-contradictory — no view satisfies a
predicate and its negation — so it is a dead rule, reported by RDR 0006's
unreachable-rule check. It MUST NOT be silently pruned at load, because a load
that dropped one block would turn an authoring mistake into a rule that
matches, which is the more dangerous failure.
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
kernel assembles (RDR 0007 A22: the kernel canonicalizes nothing). The shipped
kernel already makes this decision and is the confirming sibling:
`internal/resolve/resolve.go::TagSet.matches`, `::Lookup`, and `::has` index by
raw map lookup with no folding, and the repo's only `strings.ToLower` is on the
`--as` flag value in `internal/cli/respond/respond.go::ModeOf`.
```

```normative
**Literals carry the same byte-exact identity as keys, and a set literal is a
sequence, not a joined string.** A set-valued literal MUST normalize to an
ordered sequence of its members, each compared byte-exactly; it MUST NOT be
rendered into a single delimiter-joined string as its normalized value. Joining
on any delimiter makes membership ambiguous whenever a member contains that
delimiter — space-joining collapses `["needs work"]` and `["needs", "work"]` into
one literal, so two different predicates become one atom and the dump's atom sort
certifies the collision as canonical. Members sort byte-lexicographically so two
authored orderings of one set are one literal; the sort is over members, after
which the sequence — not a joined rendering — is the value that atom identity and
the round-trip invariant compare.
```

```normative
**Rule ids and outcome literals MUST NOT contain the expansion-suffix
separator.** The identity tuple `(model id, rule id, expansion suffix)` is total
only if its fields cannot bleed into one another. Because a rendered row identity
joins rule id and suffix with `#`, a rule id containing `#` makes the rendered
identity ambiguous: rule `a#x` and rule `a` expanding on outcome `x` render the
same string. Rejecting `#` in both rule ids and alphabet members at load closes
this at the source rather than leaving the rendered form to disambiguate what the
value could not — and it is what lets the Round-Trip section's third lossy site
be a rendering caveat instead of an identity collision. Violations are a
malformed rule id and a malformed recognized outcome alphabet respectively.
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
Every rule — ordinary or escape — MUST bind exactly one outcome. **Outcome
binding reads the match blocks only** — the rule's local `match` block plus the
`match` blocks of its inherited contexts. That set MUST contain exactly one atom
on `recognized`, using `eq` or `in`, whose literal(s) are members of the
`outcomes` alphabet. Normalization lifts that atom out of the predicate set into
the row's outcome field; an `in` atom expands into one candidate row per member,
each identified by the rule id plus the outcome literal as its expansion suffix.
A rule binding zero outcomes, more than one `recognized` atom, or a literal
outside the alphabet is a load failure.

"Combined predicate set" is used in two extents in this document and they are
not interchangeable: **outcome binding** scans the match blocks (above), while
**ambiguity checking and row carriage** scan match plus `guard.all` plus
`guard.unless`. Lifting must use the narrower extent. A `recognized` atom
authored under `guard.all` or `guard.unless` MUST be refused at load as a
malformed outcome binding — never lifted. Lifting it from `unless` would inflict
the exact inversion the author guarded against, turning "this rule does not
apply to outcome X" into "this rule binds outcome X," with no diagnostic and a
normalized row that reads as intentional.

**Expansion suffixes attach only where a rule expands.** A rule binding one
outcome normalizes to an unsuffixed row whether it was authored `eq = "x"` or
`in = ["x"]` — a single-member `in` expands to one row and takes no suffix, so
two spellings of one edge cannot mint two identities. A suffix is present
exactly when the rule produced more than one row. (RDR 0003 rejects a repeated
element at parse, so no two expansions of one rule can share a suffix.)
```

```normative
`Row.RequiresOwned` has no authored form. The normalizer MUST derive it for
every row as the sorted, duplicate-free set of tag keys named by the rule's
write block and clear list; by the write-to-non-owned-tag rule every such key
is an owned tag, and by RDR 0008 none is `recognized`. An escape rule carries
neither a write block nor a clear list, so a normalized escape row MUST carry
an empty set (RDR 0007 A21).

This is a **producer obligation, not a kernel guarantee**: the kernel checks
`RequiresOwned` on escape rows identically to ordinary ones — `escapeOrRefuse`
passes escape candidates through the same gate, and
`adversarial_test.go::TestAdv2b_EscapeEdgeMustNotBypassTheOwnedStateRequirement`
freezes that a populated set on an escape row yields
`owned_state_unavailable`. Nothing downstream rejects a malformed escape row
for this; the normalizer is the only enforcement point, which is why the
obligation is stated as a MUST on it rather than derived.

What the field means is RDR 0007's (post-guard write dependencies); this RDR is
its producer (JDR 0001 §JD-3) and does not add guard-read keys to it.
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

The two fields hold **equal sets under this RDR's authoring surface**, because
every write a rule can author targets an owned tag and therefore reaches the
accessor layer. They stay distinct fields rather than collapsing to one because
RDR 0009 A4 fixes them as distinct and only the writes cross the accessor
boundary; a later RDR that admits a next-state tag which is *not* an owned write
would separate them by value. Normalization MUST populate both explicitly from
the rule and MUST NOT populate one by aliasing the other — an alias would make
that future divergence a silent behavior change rather than a compile-time one.
```

```normative
**The unified predicate set splits across the kernel's two predicate fields by
operator, and this RDR owns the split — at the kernel handoff, not in the
normalized value.** The normalized candidate row carries the atoms as **one
unified set with each atom's authored block retained**; that single field is
what the dump contract lists, what the Round-Trip invariant compares, and what
RDR 0003 reads downstream. The split is applied when a `resolve.Row` is
constructed for the kernel. It cannot be applied earlier without loss:
`Guard string` carries no per-atom block, so a normalized value already split
into `Match`/`Guard` could not satisfy the block-retention obligation above.

The kernel row exposes `Match []Tag` (an equality pattern the kernel tests
directly) and `Guard string` (handed to the guard seam), so the handoff MUST
route each atom to exactly one of them: atoms
whose operator is equality on a declared tag populate `Match`; every other atom —
set membership, comparison, existence, and every `unless` atom regardless of
operator — belongs to the guard predicate. An atom MUST NOT appear in both.

Routing an atom into `Match` that the kernel cannot evaluate as an equality tag
would make it silently fail to match rather than reach the guard seam, and
routing an equality atom into the guard defers a decidable check to a seam that
can report `guard_unevaluable`. Which surface the guard predicate itself presents
— the atom shape and its encoding — is RDR 0007's, and the reshape named in
Prerequisites replaces the `Guard string` field; this clause fixes only which
atoms are the guard's, which is stable across that reshape.
```

```normative
The tool MUST normalize the sparse source into deterministic candidate rows for
lint, resolver lookup, diagnostics, and table dumps. Each candidate row MUST
retain its source rule id and source locator.
```

```normative
The expanded table dump MUST be derived from the normalized candidate-row value
and MUST carry every field of it: row identity, source locator, outcome,
predicate atoms (with block), next-state tags, writes including `<clear>`
entries, required-owned keys, and escape failure classes — plus the derived row
kind column. `[dump]` settings MAY reorder the rendered columns; they MUST NOT
omit a field. A dump missing any field is not an expanded table dump and does
not satisfy the Round-Trip invariant.

`[dump]` is **presentation, and MUST NOT reach the normalized value**: editing it
changes column order in the rendered view and nothing else. Normalization MUST
ignore it entirely, so two models differing only in `[dump]` normalize to
identical candidate-row sets and the Round-Trip invariant — which compares values,
not dump text — is unaffected by it. It lives in the source file because the
column order a table wants reviewed travels with that table; the cost is that a
`[dump]` edit churns golden files asserted over rendered text. Golden tests
SHOULD therefore assert over the normalized value, which no `[dump]` edit can
move, and reserve rendered-text goldens for tests of rendering itself.

Dump ordering MUST be deterministic across source key order. Rows sort by row
identity, compared field by field as the tuple `(model id, rule id, expansion
suffix)`, each field compared byte-lexicographically on the post-parse string;
an absent expansion suffix sorts before any present one. Within each row,
atoms sort by (key, block, operator token, literal), and next-state tags,
writes, required-owned keys, and escape classes sort by key.

The ordering is total **over one model's rows**: rule ids are unique within a
model, and an expansion suffix is a member of the `outcomes` alphabet, which is
non-empty and duplicate-free — so no two distinct rows compare equal and no
positional tiebreak is needed (RDR 0009 A8 establishes that a tiebreak on table
position is forbidden; RDR 0001's identity rule makes two orderings of one row
set the same input). A dump spanning more than one model MUST carry a
model-unique `model id`, which is what makes the leading tuple field
discriminating; two models sharing an id is a `duplicate model id` load
failure. RDR 0006 lints exactly one model per invocation.

**One source document carries exactly one model** — `[model]` is a singular
table, not an array of tables — so `duplicate model id` is never decidable
within a single document. It is scoped to a caller that loads **several
documents into one dump or lint invocation** and MUST be checked there, over the
set of loaded models, before rows are merged for rendering. A loader handed one
document MUST NOT report it. The Testing Strategy's mutated fixture for this
category is therefore a **pair** of documents sharing a `model id`, not a single
malformed file.

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
This is this RDR's clause: the normalized row set is unordered as far as
selection is concerned, so reordering the authored rules or the dump cannot
change a disposition, and normalization MUST NOT emit any positional field a
consumer could tiebreak on.

The selection procedure itself — gate-then-count, and which refusals are
escapable — is **JDR 0001 §D2's and the kernel's, and MUST NOT be restated
here**; RDR 0007 is its normative home for the guard seam. This RDR's obligation
is only to produce rows that procedure can evaluate: exactly one outcome bound
per row, escape rows carrying their modeled failure classes, and no positional
dependence. Ordinary transition success requiring exactly one surviving
non-escape row, and a modeled escape requiring exactly one surviving escape row
for that class, are consequences of that procedure this RDR relies on rather
than defines.
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
Load-time validation failures MUST retain stable data-level categories before
CLI mapping, including at minimum malformed TOML, unknown schema field, missing
recognized outcome alphabet, malformed recognized outcome alphabet (an empty
alphabet, a duplicate member, the empty string as a member, or a member
containing the expansion-suffix separator — each distinguishable), unknown tag,
unknown context, cyclic context inheritance, write to non-owned tag, unknown
accessor, unsupported version, malformed predicate atom, malformed escape
declaration, malformed outcome binding (including a `recognized` atom authored
in a guard block), **malformed rule shape** (an ordinary rule carrying no write
block — the mirror of `malformed escape declaration`, which covers only the
escape-side shapes), malformed rule id (containing the expansion-suffix
separator), duplicate rule id, duplicate model id, and `reserved_tag_key`
(RDR 0008).

These are the checks the load/lint arity split assigns to this RDR. All but one
are single-rule — decidable from one rule plus the model's declarations.
`duplicate model id` is the sole exception and is **not** single-rule: it is
decidable only over the set of models a caller loads into one invocation (the
clause above), so a loader handed one document MUST NOT report it. It is listed
here because it is a load-time refusal this RDR owns, not because the arity
split makes it single-rule.
Cross-row findings — overlap, gap, dead row, read-before-write — carry RDR
0006's lint categories, not these.

**"Load" names the whole source-to-candidate-rows pipeline, not one callable.**
Every category above MUST be refused by that pipeline before it yields candidate
rows, whatever internal decomposition an implementation chooses. This is stated
because several categories are only decidable during normalization — outcome
binding and the lifted literal's alphabet membership are decided while lifting,
and the write-free escape check names the normalizer as its only enforcement
point — so a reading of "load-time" as "before normalization begins" would make
them unreachable. An implementation MAY split parse, validate, and normalize
across several exported calls; it MUST NOT let a category escape because the
stage that decides it ran after the one the caller thinks of as "load".

The categories are **data-level and MUST NOT depend on the CLI envelope**: this
package MUST NOT import `internal/cli`, and CLI code names, `Group` assignments,
and exit mapping remain RDR 0005's design work (A5). `internal/cli/clierr` is a
leaf package for exactly this reason, so a later mapping needs no restructuring
here.

Category identifiers are **snake_case renderings of the prose names above**,
anchored on `reserved_tag_key` — the one member spelled as an identifier, whose
spelling RDR 0008 owns. The Testing Strategy asserts on the category, so these
identifiers are an API surface, not message text.
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
  **match blocks** (local `match` plus inherited context `match`), lifted into
  the row's outcome field. The extent is the narrow one the Normative Contracts
  fix; "combined predicate set" is this document's *wider* term (match plus
  `guard.all` plus `guard.unless`) and MUST NOT be read here. Rejected: an
  eventless row with no outcome, because the kernel has no such concept and
  would never match it; rejected: leaving the atom in the predicate set only,
  because the kernel filters on `Row.Outcome` before matching.

#### Round-Trip / Inverse Invariants

`parse ∘ normalize = expanded-table value identity` on valid transition model
fixtures: two authorings of one model — differing in TOML key order, rule
declaration order, or `eq`-versus-single-member-`in` spelling — must normalize
to the same candidate-row set, compared as values over the full field list the
dump contract carries (row identity, source locator, outcome, predicate atoms
with block, next-state tags, writes, required-owned keys, escape failure
classes), with the source locator's optional line/column detail excluded from
the comparison.

The locator's **presence and its rule-identifying part are compared**, even
though its line/column detail is not: excluding the whole field would let a
normalizer that dropped locators entirely satisfy the invariant, and the locator
is what routes a lint or refusal diagnostic back to the authored rule. Only the
positional detail is exempt, because it moves when unrelated source text is
edited — the same reason the locator is excluded from row ordering.

The invariant is stated over the **normalized value, not the dump text**,
because this RDR fixes the dump's field list and row ordering but does not
define a dump grammar — no delimiter, escaping, quoting, or record separator.
`dump` therefore has no specified inverse, and a textual round-trip is not
claimed. Three sites make the rendered form lossy today and would each have to
be closed by any later RDR that defines a re-readable dump: the `<clear>`
sentinel is not reserved in the tag-value space, so an authored value of
`<clear>` renders identically to a cleared tag; set-valued and multi-entry
fields are rendered with unescaped separators; and the expansion suffix has no
reserved separator, so a rule id containing it makes `(rule id, suffix)`
unrecoverable from the rendered identity. Rendering is a review surface here,
and the normalized value is the contract.
Source rewrite is likewise out of scope for this RDR; a later rewrite-capability
RDR must define its own source-preservation invariant before mutating authored
TOML.

#### Conditional Mini-Checks

**`fidelity`** — the chain's invariant holds over the normalized value; the
rendered form is lossy at three named sites and carries no inverse.

| Operation | Invariant | Lossy exemptions |
| --- | --- | --- |
| parse (TOML → source) | none claimed — comments, key order, and whitespace are dropped by design | source rewrite out of scope |
| normalize (source → rows) | **value identity**: key order, rule order, and `eq`/single-member-`in` spelling all normalize to one candidate-row set | contexts flattened; `recognized` atom lifted out; `in` expands 1→N; set literals member-sorted; `all`/`unless` merged with block retained |
| dump (rows → text) | fixed field list + total row order; every field present, columns reorderable | **no grammar**: `<clear>` unreserved in the value space; separators unescaped; suffix separator unreserved |
| read-back | **not claimed** — no inverse is specified | charted to a successor dump-format RDR |

**`disposition`** — routing is by arity: single-rule → load (this RDR),
cross-row → lint (RDR 0006), evaluation-time → kernel refusal (RDR 0001).
`duplicate model id` is the one load-time check that is cross-**document**
rather than single-rule (Normative Contracts); it stays this RDR's.

| Input class | Bucket | Owner | Silent vs loud |
| --- | --- | --- | --- |
| malformed TOML, unknown schema field, unsupported/duplicate version or model id | load | this RDR | loud, stable category |
| unknown/cyclic context, unknown tag, unknown accessor, write to non-owned tag | load | this RDR | loud, stable category |
| malformed predicate atom, escape declaration, outcome binding, alphabet; `reserved_tag_key` | load | this RDR (0008 for the reserved key) | loud, stable category |
| overlap, gap, dead row, read-before-write, ambiguous expansion | lint | RDR 0006 | loud, blocking |
| expansion counts per rule | lint | RDR 0006 | reported, non-blocking; no threshold here |
| `unmodeled_outcome`, `owned_state_unavailable`, `guard_unevaluable` | kernel refusal | RDR 0001 | loud; the last two are never escapable |
| zero/multiple survivors | kernel refusal, or modeled escape when exactly one escape row survives the gate | RDR 0001 | loud |

**`oracle`** — every MVV assertion carries a failing control; the rail and the
per-item controls are in the Minimum Viable Validation and Testing Strategy.
The three that previously passed against a wrong implementation — the category
oracle (message text, not category), the determinism oracles (a constant
passes), and "the comparator is total" (satisfied by the positional tiebreak it
excludes) — are replaced by category-level assertions, RDR-fixture permutation,
and rule-order permutation respectively.

**`trace`** — desk trace of the MVV against the clauses in force. One row per
step; witness values from `evidence/spikes/output.txt`.

| Step | Assertions in force | Witness / verdict |
| --- | --- | --- |
| load fixture | alphabet non-empty, dup-free, no empty string; version = 1 | `outcomes=round-clean,verdict-flapping,reconcile-block,finalized` — holds |
| normalize | outcome lifted; per-atom block retained; `RequiresOwned` = writes ∪ clears | `continue-prelock … atoms=[…profile.eq=small@unless…] requires_owned=[iter,stage]` — holds |
| expand `in` | one row per member; suffix only where the rule expands | `terminal-archive#finalized` / `#verdict-flapping` — holds; the spike suffixes on expansion count (`main.go` `if len(outcomes) > 1`), matching the contract, so `eq`/single-member-`in` equivalence is witnessed by construction but has no fixture row — the permutation test owes it |
| escape row | no write block or clear list; `RequiresOwned` empty | `draft-no-match-escape … next=[] write=[] requires_owned=[]` — holds |
| dump order | rows by identity tuple; locator excluded | three runs byte-identical — holds |
| resolve, sibling gate | unevaluable survivor refuses despite a decidable sibling and an escape row | **no witness possible in the current fixture** — every outcome binds exactly one transition row; closed by the MVV's two-sibling requirement |

No CONTRADICTION row survives: the sibling gap and the suffix divergence are
booked as needs-verification items for Stage 6, not left as silent assumptions.

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
| Structured CLI failures | `internal/cli/clierr::CLIError` | Envelope exists; no parse-validation code has ever travelled it | Extend | Add stable parse/lint refusal codes later — these are the first of their kind, so names, `Group`, and exit mapping are new design work. |
| Text/json output gateway | `internal/cli/respond::Fail` | Gateway is CLI-only, not kernel behavior | Reuse | Parser/lint commands must report through existing gateway. |
| Kernel row and resolver | `internal/resolve::Row`, `::Resolve` | Consumes rows; has no loader or normalizer | Reuse | The normalizer targets this row shape; no new kernel surface. |
| Normalizer/table package | `internal/` search | Not implemented | Introduce | New internal package can own sparse source structs and normalization. |
| Config discovery | `internal/cli/config::Load` | Discovery and file read exist; no parsing or validation is wired | Reuse later | Table path binding belongs with CLI integration, not this RDR's data format. The load entry therefore takes **already-read bytes plus a source id** for the locator, not a filesystem path — this package performs no file I/O and no path resolution. A path-taking convenience wrapper MAY live at the CLI seam. |

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
  permutes rule declaration order and asserts the normalized value is
  unchanged, so a later identity change cannot reintroduce a positional
  tiebreak unnoticed.
- **Risk**: This RDR is the cluster's only unlocked document — 0001 is
  `Implemented` and 0003, 0004, 0005, 0006, 0007, 0008, 0009 are `Final` — so
  the producer locks *after* its consumers, inverting the dependency order.
  RDR 0009 twice cites "RDR 0002's **Final** escape-rule prohibition" and rests
  its enforcement argument on it, and RDR 0008's normative block still asserts
  "nothing in this RDR or RDR 0002 forbids that alphabet entry" about the
  empty-string outcome, which this RDR now **does** forbid and which A8 leans on
  as its sole barrier. The decisions are compatible — 0008 explicitly scoped
  itself out ("not this RDR's to rule on") — so this is stale peer *fact*, not
  a contract conflict, and RDRs are never amended.
  **Mitigation**: Every contract this pass changed is checked against the Final
  peers' verified assumptions before implementation via
  `/rdr-cluster-reconcile`, which is the mechanism that owns cross-RDR drift.
  Routed there by the cove lens (CV-020, CV-021) and carried here so the
  obligation survives lock.

### Failure Modes

Malformed TOML, unknown schema fields, unresolvable or cyclic context
references, malformed predicate atoms, malformed escape declarations, malformed
outcome bindings, a misnamed recognized declaration, or a write to an undeclared
or non-owned tag fail at load/validation time with stable CLI errors — each is
decidable from one rule plus the declarations. A row gap, overlap, dead row,
read-before-write condition, or ambiguous expansion requires comparing
normalized rows and is therefore an RDR 0006 lint failure before the model is
accepted. Expansion counts are a lint diagnostic reported per rule, not a
refusal: this RDR sets no expansion threshold. At runtime, unknown outcomes,
unavailable accessor inputs, absent owned state, undecidable guards, and
unmodeled zero/multiple survivors are typed resolver refusals rather than
guessed edges; modeled zero/multiple-match behavior must come from exactly one
surviving escape row, and the gate refusals are never modeled.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified (A2, A7 re-verified; A8 resolved).
  A1–A8 are terminal. **A9–A14 are `Pending` by deliberate Stage 6 disposition,
  every one of them a DOWNGRADE with a named MVV assertion**, not an unexamined
  gap: A9/A12 are blocked on RDR 0007's kernel reshape (below) and are
  unverifiable before it lands; A10's decidability half is source-verified and
  only its paired-document behavior is owed; A11/A13/A14 are contract clauses
  the review rounds introduced whose oracles now exist in Testing Strategy
  scenarios 2 and 3 but whose fixtures are owed at implementation. None is
  load-bearing for a *pre-lock* MVV assertion; each is survivable per its
  "If wrong".
- [ ] RDR 0007 is the normative home of the atom shape and the existence
  constants this RDR's normalizer emits. RDR 0007 is `Final`, so the
  *specification* is settled — but its kernel reshape is unimplemented, and the
  gap is the whole atom surface, not only the constants. On `main`
  `internal/resolve/resolve.go::Row.Guard` is a `string`, `Row.Match` is a flat
  `[]Tag` equality list matched by `::TagSet.matches` with no operator
  dispatch, and there is no `Atom` type, no operator-token type, no per-atom
  block field, and no `OpExists` / `LiteralTrue` / `LiteralFalse` anywhere in
  the repo. Of Phase 2's four deliverables only outcome lifting and
  `RequiresOwned` derivation have fields to land in; per-atom block retention
  and verbatim existence constants have none. Phases 2 and 3 therefore sequence
  behind the reshape in their entirety, and the Testing Strategy assertions that
  cite `resolve.OpExists` / `LiteralTrue` / `LiteralFalse` are unsatisfiable
  until it lands. This item gates implementation sequencing, not lock. It does not
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
deliberately overlapping variant is reported by lint as an ambiguous overlap.

The overlap item is the one MVV assertion this RDR cannot discharge alone: the
ambiguous-overlap check is cross-row and therefore RDR 0006's by the arity
split, and RDR 0006 is unimplemented. It is **deferred to the Phase 5 lint
handshake rather than counted as satisfied here** — the rest of the MVV is
executable against this RDR's own load and normalize paths. What this RDR owes
before that handshake is that the normalized rows *carry enough structure* for
the check: overlapping rows must be distinguishable by row identity and
comparable by predicate set, which the identity tuple and the atom-set contract
above already fix. Marking the overlap assertion green before RDR 0006 lands
would be asserting on a stub.

The RDR fixture must cover at least `Status`, `Profile`, prelock iteration, one
guard that actually carries a predicate, one existence atom, one `in` atom on
`recognized` that expands into more than one row, and one explicit `no_match`
escape row. It must additionally carry **two sibling candidate rows binding the
same outcome** — without them the gate's ordering is untestable, since a
one-candidate-per-outcome fixture can never exercise "an unevaluable survivor
refuses even though a decidable sibling and a modeled escape both exist," which
is the property this RDR's row shape exists to support.

**The spike fixtures do not yet carry those siblings, and the MVV requirement
wins over the promote-verbatim mandate.** The current RDR fixture binds
`round-clean` on one ordinary rule and one escape rule, so no two *ordinary*
candidates ever contend — which is why the desk trace records "no witness
possible in the current fixture" for the sibling-gate row. Promotion is
therefore **extend-then-promote**: the production fixture is the spike fixture
plus a second ordinary rule binding an outcome one other ordinary rule already
binds. The Load-Bearing Decisions' "canonical examples" clause fixes the field
layout and the authoring idiom those fixtures demonstrate, not their row census;
a fixture may not be narrowed on promotion, only extended.

**Every oracle must have a failing control.** An assertion that passes by
absence-of-error, exit-0, or fixture-name match does not discharge an MVV item:

- Each load-time validation category is asserted by **category**, not by
  message text, and each has one mutated fixture that trips it and no other —
  otherwise an implementation collapsing all categories into one code passes.
- Determinism and key-order assertions are paired with the value assertion
  above; alone they are satisfied by a normalizer emitting a constant or
  dropping every `unless` atom. The permutation control runs on the RDR
  fixture, which carries the escape row, the `in`-expansion, `<clear>`, and
  inherited contexts — not only on the simpler kata fixture.
- The existence-atom test exercises `exists = true` **and** `exists = false`; a
  fixture carrying only one lets a hardcoded literal pass.
- The byte-equality contract is asserted by a negative control — a reference
  spelling a declared `status` as `Status` must fail `unknown tag` — since a
  normalizer applying case folding passes every positive test.
- Unknown schema fields are asserted by a fixture carrying one; the decoder
  must be configured to reject unmapped keys, which is not the TOML library's
  default and is otherwise a silent no-op.

The dump ordering is MVV-covered, not left to implementation: assert that the
normalized value is unchanged when the same model is authored with its TOML
keys in a different order, **and when its rules are declared in a different
order**, and that repeated dumps of one model are byte-identical. Row-order
permutation is the assertion that actually excludes a positional tiebreak;
sorting and checking that no adjacent pair compares equal does not, because a
comparator falling back to table position produces no equal adjacent pair
either — it is satisfied by the defect it was meant to catch.

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
   **Expected**: Candidate rows retain source rule ids/source locators, inherited predicates are expanded, each atom reports its authored block (`match`/`all`/`unless`), the single `recognized` atom is lifted into the row's outcome field and absent from the predicate set, `RequiresOwned` equals the sorted write-plus-clear key set (empty on the escape row), escape rows retain their modeled failure class list — from which the derived `escape` kind is computed — and carry neither writes nor next-state tags, and writes are deterministic. `terminal-archive`'s `in` atom expands to two rows carrying the outcome literal as their expansion suffix. The normative fixture for this scenario is `evidence/spikes/output.txt`.
   **The no-alias obligation on next-state tags and writes is not assertable by value comparison** — the two fields hold equal sets under this RDR's authoring surface, so an aliased pair and an independently built pair compare equal. On `main` both are `[]Tag` slices, so an alias also shares a backing array and a later mutation of one would silently move the other. Enforcement is by review of the normalizer, or by a test that mutates one field and asserts the other is unchanged.
   **Additionally**, atom-set semantics are asserted positively, since these shapes must *survive* load rather than be refused: two inherited contexts contributing the same key and operator with **different literals** yield **two** atoms in the normalized set, not one — asserted by count and by value, with the contexts listed in both `use` orders to prove the result is order-independent (a merge keyed on `(block, key, operator)` drops one and passes a count-only check on a single ordering); a set-valued literal whose member contains a space normalizes to a member sequence, so `["needs work"]` and `["needs", "work"]` are **distinct** atoms; and a rule authoring one atom in both `all` and `unless` loads successfully as two atoms and is reported by lint as a dead rule, not refused at load.
   **The idempotence mirror of that assertion is required too (A13)**: a rule and an inherited context contributing a **byte-identical** `(block, key, operator, literal)` atom collapse to **exactly one** atom, asserted by count on that key. Without it a normalizer that de-duplicates nothing passes every control above, since they only ever count the two-distinct-literal case.
   **Escape expansion under `in` is asserted here (A11)**: the extended fixture carries an escape rule binding its outcome with `in` over two alphabet members, and normalization must yield **two** escape rows carrying distinct expansion suffixes, each with an empty write set and each retaining the modeled failure-class list. The current spike fixture's `draft-no-match-escape` binds with `eq`, so this shape has no captured witness and is owed at implementation.
   **Handoff routing is asserted as total and disjoint (A12)**: constructing the `resolve.Row` for each normalized RDR-fixture row, the union of `Match` and the guard's atoms equals the normalized atom set and their intersection is empty. The kata fixture's `status.eq=closed@unless` atom is the discriminating case — an equality operator that must route to the guard because its block is `unless`. Like the rest of Phase 2 this assertion is unsatisfiable until RDR 0007's reshape lands (Prerequisites).
3. **Scenario**: Validate one malformed variant per load-time category — unknown tag written, unknown tag matched, unknown context, cyclic context inheritance, writes to non-owned tags, unknown accessors, unsupported versions, unknown schema field, malformed predicate atoms (unknown operator; `exists` with a non-boolean literal), malformed escape declarations (an empty write block on an escape rule, and separately an escape rule carrying a `clear` list — the normalizer is the only enforcement point for a write-free escape row, so each shape needs its own control), a rule with zero or two `recognized` atoms, a `recognized` atom authored under `guard.all` and one under `guard.unless`, a rule id containing `#` and an alphabet member containing `#`, an outcome literal outside the alphabet, an alphabet containing the empty string, a duplicate alphabet member, an empty alphabet, an owned declaration named `recognized` and a recognized declaration named `outcome`, a duplicate rule id, a duplicate model id (a **pair** of documents sharing a `model id`, loaded into one invocation — the only surface on which the category is decidable), an ordinary rule carrying no write block, and a missing root outcome alphabet.
   **Expected**: Each variant is refused with the **one** category its mutation targets and no other — the assertion is on the category, not the message text, so an implementation collapsing several categories into one code fails. Ambiguous overlap is deliberately absent from this scenario: it is cross-row and therefore an RDR 0006 lint finding (scenario 4), not a load failure. Nine categories are already witnessed by the spike, one refusal per mutated fixture, in `evidence/spikes/negative-cases.txt`: unsupported version, empty-string alphabet member, duplicate rule id, unknown written tag, write to an observed tag, escape rule carrying a write block, `exists` with a non-boolean literal, an outcome literal outside the alphabet, and a recognized declaration misnamed `outcome`. The remainder are owed at implementation.
   **The version gate's precedence needs its own fixture (A14), separate from the `unsupported version` category fixture above.** That fixture is a v1-shaped document with a bad version value, which trips the category under *either* check ordering and therefore witnesses nothing about precedence. The precedence control is a **v2-shaped** document — `version = 2` plus a v2-only key the strict decoder would reject as an unknown schema field — asserted to refuse `unsupported version` and **not** `unknown schema field`. This is the only ordering this RDR fixes normatively, so it is the only ordering that gets an oracle.
4. **Scenario**: Run `internal/resolve::Resolve` over one matching ordinary tag-set in which every sibling candidate's guard is decidable, one tag-set with no ordinary match but one matching `no_match` escape row, and one tag-set in which a sibling candidate's guard is unevaluable — all three drawn from the fixture's two same-outcome sibling rows. Separately, run RDR 0006's lint over a deliberately overlapping variant.
   **Expected**: The matching ordinary tag-set resolves to one transition row; the no-match tag-set resolves to the modeled escape disposition; the unevaluable-sibling tag-set refuses `guard_unevaluable` even though a decidable sibling and a `no_match` escape row exist; zero or multiple survivors without exactly one surviving escape row are refusals and never fall back to row order. The overlapping variant is reported by lint as an ambiguous overlap before the model is accepted, naming both rows.
5. **Scenario**: Normalize and dump semantically identical variants of the **RDR** fixture — one with its TOML keys authored in a different order, one with its `[[rule]]` blocks declared in a different order, and one spelling a single-outcome binding `in = ["x"]` where the original spells `eq = "x"` — and dump the same fixture repeatedly in one process and across processes.
   **Expected**: The normalized value is identical in every case because rows sort by the identity tuple and atoms/next tags/writes/required-owned keys/escape classes sort by key, and because a single-member `in` takes no expansion suffix; repeated dumps are byte-identical despite Go's randomized map iteration. Rule-order permutation is what excludes a positional tiebreak. Partially witnessed in `evidence/spikes/negative-cases.txt` (reordered **kata** fixture, matching SHA-256) and by three consecutive byte-identical runs recorded under Performance Expectations; the RDR-fixture permutations and the `in`/`eq` equivalence are owed at implementation.
6. **Scenario**: Normalize a rule carrying `exists = true` and one carrying `exists = false`, one carrying a non-boolean existence literal, and a variant spelling a declared tag `Status` where the declaration is `status`.
   **Expected**: Emitted atoms carry `resolve.OpExists` and `resolve.LiteralTrue`/`LiteralFalse` byte-for-byte, referencing the kernel constants directly rather than a local mirror — so this assertion compiles only once RDR 0007's reshape exports them (Prerequisites), and is owed at that point rather than satisfiable today. The non-boolean literal is refused at load as a malformed predicate atom; the mis-cased reference fails `unknown tag` rather than folding. The RDR fixture's `continue-prelock` guard carries the `finalized_at.exists=false@all` atom the spike emits; the spike mirrors the constants locally because the kernel does not yet export them.

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

Responses: 0002-transition-table-as-reviewable-data/artifacts/gate.md (Gate PASS 2026-08-23)

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
