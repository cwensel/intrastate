# Recommendation 0007: Normalized-table revision binding

## Metadata

- **Date**: 2026-08-06
- **Status**: Draft
- **Type**: Architecture
- **Profile**: foundational — defines one revision-binding contract across
  normalized-table production and resolver consumption.
- **Priority**: Medium
- **Related Issues**: kata jv85; RDR 0001; RDR 0002
- **Predecessors**: 0001-resolution-kernel,
  0002-transition-table-as-reviewable-data
- **Overrides**: RDR 0001's closed refusal-kind set and RDR 0002's explicit
  exclusion of content-addressed table identity
- **Seam Lineage**: area:resolver — no prior accretion

## Problem Statement

Developers replaying resolver decisions need the same declared replay identity
to select the same transition rule, so tests and audits can trust reproduced
dispositions. They discover the gap when identical nominal revision values are
paired with different normalized transition-table contents and resolution can
silently follow different edges. The system therefore needs a defined ownership
and validation contract between normalized-table revision identity and the table
used by the resolver, including a canonical mismatch outcome before edge
selection.

## Context

### Background

RDR 0001 includes the transition-table revision in replay identity, while the
current resolver selects from caller-provided table content without validating
that the content is bound to the claimed revision. RDR 0002 deliberately leaves
a content-addressed identity contract unspecified. As a result, a buggy or
adversarial caller can reuse a revision string with altered normalized rows,
undermining deterministic replay and audit claims; existing tests only replay
the same table value and do not constrain this case. Kata jv85 was opened from
review findings that exposed this cross-RDR contract gap.

### Technical Environment

The affected system is the Go resolver centered on
`internal/resolver/resolver.go::Resolve`, its caller-provided normalized
transition table, and the loading or normalization boundary that supplies table
revision identity. RDR 0001 defines resolver replay identity and RDR 0002 defines
the transition table as reviewable data. The binding boundary must preserve the
resolver's structured value/error split and the CLI output contract owned by
RDR 0005.

## Research Findings

### Investigation

The freshness check confirmed that RDRs 0001 and 0002 remain Final, kata jv85
remains open and tracks this RDR, and the named resolver seam still exists.
`internal/resolver/resolver.go::Input` carries `TableRevision` beside caller-
supplied `Outcomes` and `Table`, but `Resolve` never reads `FlowID` or
`TableRevision`; current replay tests reuse the same table value. Draft RDR 0008
owns selected-rule identity after selection and explicitly leaves table-revision
binding here.

Prior art was read before alternatives were named. RDR 0001 `Normative
Contracts` says that the same transition-table revision must replay the same
disposition, while RDR 0002 `Cross-Cutting Concerns` calls the normalized
candidate-row value canonical but explicitly says its spike hash is not a
content-addressed contract. Pro Git, chapter 10.2, describes an object key as a
"checksum of the content you're storing plus a header"; Domain-Driven Design,
chapter 5, says that for an object to be shared safely it must be immutable.
Srinivasan, *A Methodology for Selecting and Composing Runtime Architecture
Patterns for Production LLM Agents*, page 7, observes that changing a model
version can make the same event produce different replay output. These anchors
favor a content-derived, versioned identity bound to an immutable semantic
value, not a free revision label beside mutable rows. Queries and rejected
branches are recorded under
`docs/rdr/0007-normalized-table-revision-binding/evidence/research/`.

Sibling-path check: a repository-wide search for `sha256`, `digest`,
`checksum`, content addressing, and `TableRevision` found only the inert
`TableRevision` field and its fixtures under `internal/resolver`; no adjacent
digest, registry binding, or identity discriminator exists to reuse.

### Key Discoveries

- **Documented** — RDR 0001 makes table revision part of replay identity but
  deliberately introduces no hash or canonical serialization.
- **Documented** — RDR 0002 defines the normalized candidate-row value as the
  canonical semantic form while explicitly leaving content-addressed identity
  outside its contract.
- **Documented** — Git-style content addressing binds a compact identifier to
  the bytes identified, and immutable value-object practice prevents identity
  from drifting after construction.
- **Verified** — `internal/resolver/resolver.go::Resolve` selects from caller-
  supplied rows without consulting `TableRevision`, so the same nominal replay
  tuple can select different edges.
- **Verified** — the Stage 4 canonicalization spike produces the same pre-image
  and revision after semantic-set reordering, map iteration, source-locator
  changes, source-schema-version changes, and nil/empty construction changes;
  every projected one-field mutation and clear-versus-empty boundary changes
  both.
- **Verified** — Go's `crypto/sha256::Sum256` and
  `encoding/hex::EncodeToString` provide the intended standard-library digest
  and lowercase 64-digit encoding without a third-party dependency. Revision
  parsing must still reject uppercase explicitly because hex decoding accepts
  it.
- **Blocked** — the proposed version-1 projection is not yet complete against
  RDR 0002 and the concrete resolver value: row outcome, predicate
  provenance/operator, write role, and guard unevaluability lack explicit
  include-or-eliminate dispositions.
- **Verified** — RDR 0001's value-level refusal contract and RDR 0005's stable
  refusal-to-`CLIError` mapping can preserve a distinct revision-mismatch
  outcome without using the Go error path, provided the closed refusal set and
  CLI code set are extended explicitly.
- **Blocked** — the Stage 4 contract recount finds at least three independent
  load-bearing contracts: canonical hash/pre-image, public revision grammar,
  and resolver refusal/ordering policy. The Resolve profile rule treats that as
  a split signal, so the provisional `foundational` profile cannot be latched
  until Refine consolidates or splits those seams.

### Critical Assumptions

- **A1 The normalized semantic table can be encoded canonically without source
  formatting, map iteration, or caller row order changing its revision.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `GOCACHE=/tmp/intrastate-rdr7-go-cache go run
    ./docs/rdr/0007-normalized-table-revision-binding/evidence/spikes` produced
    the golden vector and comparison matrix captured at
    `docs/rdr/0007-normalized-table-revision-binding/evidence/spikes/output.txt`:
    all ordering, map, source-only, and nil/empty variants were byte- and
    revision-identical, while every projected semantic mutation,
    clear-versus-empty, and ambiguous unframed concatenation differed.
  - **If wrong**: Semantically identical tables churn revision or distinct
    resolver behavior aliases to one revision, defeating reliable replay.
- **A2 Every field that can change edge selection, action, or logical
  selected-rule identity can be included in one versioned digest pre-image,
  while diagnostic source location can be excluded.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: The Stage 4 audit confirms RDR 0008 `Load-Bearing Decisions /
    Identity` makes `SourceLocator` diagnostic-only, but the proposed inventory
    does not explicitly project or eliminate RDR 0002/current resolver row
    outcome, predicate provenance/operator, write role, or guard
    unevaluability. Refine the inventory, then extend the one-field-change
    vectors and repeat this audit.
  - **If wrong**: A caller can alter an omitted outcome, predicate, action, or
    rule identity while retaining the same revision.
- **A3 Go's standard SHA-256 implementation and lowercase hexadecimal encoding
  are available and sufficient for this table-identity boundary.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: Go 1.26 `crypto/sha256::Sum256` returns the standardized
    32-byte digest and `encoding/hex::EncodeToString` emits the required 64
    lowercase digits. The revision parser must separately enforce the
    `sha256:` prefix, exact length, and lowercase-only grammar because
    `encoding/hex::DecodeString` accepts uppercase.
  - **If wrong**: The format needs another algorithm or dependency before its
    revision grammar can be locked.
- **A4 The normalizer and resolver can share one revision function without
  exposing mutable digest state or a second implementation.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: The peer/source audit found no existing production digest or
    binding implementation and confirmed a pure shared function is feasible,
    but `TestResolveRejectsTableRevisionMismatchBeforeSelection` alone proves
    recomputation, not single implementation. Pair that mutation-before-
    selection MVV with a source-level definition/call-site audit showing both
    producer and resolver use the same exported derivation symbol.
  - **If wrong**: Producer and consumer may accept different revisions for the
    same table, or caller mutation may invalidate a cached binding silently.
- **A5 Downstream resolver and CLI mapping can preserve
  `table_revision_mismatch` as a distinct modeled refusal without using the Go
  error path.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0001 `Normative Contracts` makes resolver refusal a
    value-level disposition, and RDR 0005 `Approach` / `Technical Design`
    requires `flow resolve` to map each stable refusal to `CLIError`. This RDR
    can extend both closed sets with `table_revision_mismatch` and a distinct
    `flow-table-revision-mismatch` CLI code while reserving malformed syntax and
    invalid normalized values for the Go/load error path.
  - **If wrong**: Callers cannot diagnose a replay-identity failure uniformly,
    and may conflate adversarial input with parser or programmer failure.

## Proposed Solution

### Approach

Make transition-table revision a content-derived identity of the normalized
semantic table. The normalization/loading boundary derives a revision using one
shared, versioned SHA-256 function and returns the revision with the normalized
outcome alphabet and candidate rows as one table value. A replay caller records
that revision as `Input.TableRevision`; `Resolve` recomputes the supplied table's
revision with the same function and compares it before evaluating owned state,
outcomes, guards, or edges.

The digest covers the normalized semantic value, not TOML bytes or a rendered
dump. It binds every match, action, and logical selected-rule identity field,
but excludes `SourceLocator`: RDR 0008 defines that field as diagnostic
provenance rather than logical rule identity. Its external grammar is
algorithm-qualified lowercase hexadecimal, initially `sha256:<64-hex>`; the
pre-image is domain-separated and versioned so later semantic formats cannot
accidentally share an identity namespace.

An equality match allows existing exact-one selection to proceed unchanged. A
mismatch returns a new value-level `table_revision_mismatch` refusal with no
transition plan and never attempts an ordinary or escape edge. Malformed
revision syntax or an impossible normalized value remains a load/programmer
error before the modeled resolver boundary; only a well-formed claimed revision
that does not identify the supplied semantic table produces the refusal.

### Technical Design

The normalized-table producer owns revision derivation because it owns the
semantic value being identified. It uses one package-level canonical encoder
and digest function shared with the resolver; there is no second hash
implementation in CLI code. The table value groups the recognized-outcome
alphabet and normalized candidate rows that RDR 0001 currently receives as
separate `Outcomes` and `Table` fields. RDR 0002 remains the owner of those
values' semantics and source identity; this RDR owns only their content-derived
revision and the binding check.

Canonical encoding operates over a version-1 semantic projection, not Go field
layout. The pre-image begins with the ASCII bytes
`intrastate.transition-table-revision`, a NUL byte, and the ASCII encoding label
`v1`. Every projected value is framed as a one-byte type tag, an unsigned
64-bit big-endian payload length, and the payload: `0x01` is a UTF-8 string,
`0x02` an unsigned 64-bit big-endian integer, `0x03` a list, and `0x04` a
record. A list payload begins with an unsigned 64-bit element count followed by
framed elements. A record payload begins with a field count followed by pairs
of framed field-name strings and framed values, sorted by field-name bytes.

The top-level record fields are `model_id`, `outcomes`, and `rows`; RDR 0002's
source-schema version is excluded because it is not part of the normalized
semantic value. Each row record contains `rule_id`, `expansion_suffix`, `kind`,
`predicates`, `escape_classes`, `next_tags`, and `writes`. Predicate entries are
records of `key` and `value`; next-tag entries are records of `tag` and `value`;
write entries are records of `tag`, `operation`, and, for `set`, `value`.
`operation` is either `set` or `clear`. Outcomes, rows, predicates, escape
classes, next tags, and writes are semantic sets and are sorted by their
complete framed element bytes before list framing.

`SourceLocator` and every other source-formatting or diagnostic-only field are
excluded. Strings are their validated UTF-8 bytes without additional Unicode
normalization. The explicit `operation` field distinguishes clear from setting
an empty string; authored absent and present-empty forms that normalize to the
same typed value cannot produce different revisions. Stage 4 must verify this
projection against the concrete RDR 0002/RDR 0008 types and produce golden
vectors before the contract locks; it may not change these inclusion,
exclusion, ordering, framing, or primitive-encoding rules.

`Resolve` validates revision binding at the start of modeled input handling. It
derives the actual revision from the supplied normalized table and compares it
to the claimed replay revision. On mismatch it returns exactly one refusal and
does not call outcome modeling, guard evaluation, ordinary matching, or escape
matching. On equality, existing resolver semantics continue. RDR 0008 may carry
selected-rule identity through a successful plan, but it neither derives nor
validates table revision.

#### Normative Contracts

```normative
A normalized transition table revision MUST be the content-derived identity of
the complete normalized semantic table, including its recognized-outcome
alphabet, candidate selection data, transition actions, and logical rule
identity. Diagnostic provenance, including `SourceLocator`, MUST NOT participate
in revision identity. The revision MUST NOT be an unconstrained caller label,
TOML source checksum, rendered-dump checksum, registry sequence, timestamp, or
source path.

The revision grammar MUST be algorithm-qualified lowercase hexadecimal. Version
1 MUST use `sha256:<64 lowercase hexadecimal digits>` over a domain-separated,
versioned canonical encoding. One shared derivation function MUST be used by the
normalization/loading boundary and resolver validation.

Canonical encoding MUST use the version-1 domain prefix, recursively typed
length framing, semantic field projection, and encoded-byte sorting defined in
Technical Design. It MUST be independent of TOML whitespace,
comments, source key order, map iteration, normalized candidate input order,
and diagnostic source location. The projection-to-normalized-type inventory and
golden vectors MUST be verified before Final.

Before evaluating owned-state availability, the outcome alphabet, guards,
ordinary edges, or escape edges, `Resolve` MUST derive the supplied normalized
table's revision and compare it to the claimed replay revision. Equality permits
normal resolution. Inequality MUST return a value-level
`table_revision_mismatch` refusal with no transition plan and MUST NOT evaluate
or select any table edge.

Malformed revision syntax and failure to construct a valid normalized semantic
table MUST use the load/programmer error path, not
`table_revision_mismatch`. The mismatch refusal is reserved for a well-formed
claimed revision that does not identify the supplied normalized value.
```

#### Load-Bearing Decisions

- **Identity** — two normalized tables have the same revision only when their
  version-1 canonical semantic encodings have the same SHA-256 digest. Source
  formatting, diagnostic source location, and candidate input order do not
  distinguish identity; any field that can change matching, action, or logical
  selected-rule identity does.
- **Wire / byte format** — the public revision is
  `sha256:<64 lowercase hexadecimal digits>`. The pre-image is a domain-
  separated, versioned, recursively typed and length-framed encoding of the
  normalized semantic value. Stage 4 verifies the projection and golden
  vectors; it does not choose different fields, ordering, primitives, or
  framing.
- **Naming** — the content identity remains `TableRevision` in replay input;
  the mismatch is `table_revision_mismatch`. Rejected: `version`, which RDR 0002
  already uses for source-schema version, and `etag`, which implies transport
  cache semantics.
- **Selection / predicate** — revision equality is a mandatory gate before any
  candidate qualifies. A mismatch selects no ordinary or escape row and cannot
  itself be overridden by a table-authored escape.

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Exact-one resolver selection | RDR 0001 | Available | Runs only after revision binding succeeds. |
| Canonical normalized semantic table | RDR 0002 | Predecessor | Supplies the value whose behavior and identity fields are hashed. |
| Logical selected-rule identity fields | RDR 0002 / RDR 0008 | Deferred peer | Included in the table pre-image; diagnostic source location is excluded. |
| SHA-256 and hexadecimal encoding | Go standard library | Available | A3 verifies `crypto/sha256::Sum256` and lowercase `encoding/hex` output; grammar validation remains explicit. |
| Table revision derivation and validation | This RDR | Introduced | One shared function binds producer output to resolver replay input. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Replay revision claim | `internal/resolver/resolver.go::Input.TableRevision` | Never read | Extend | Treat as the claimed content identity. |
| Normalized table input | `internal/resolver/resolver.go::Input.Outcomes` and `Table` | Separately mutable; no binding | Replace shape | Group as one normalized semantic table value. |
| Resolver gate | `internal/resolver/resolver.go::Resolve` | Starts with owned-state validation | Extend | Validate revision before every existing modeled input check. |
| Refusal taxonomy | `internal/resolver/resolver.go::RefusalKind` | RDR 0001 closed set omits mismatch | Extend | Add one stable value-level refusal and downstream mapping. |
| Revision/digest helper | `internal/` sibling-path search | None exists | Introduce | Keep canonicalization and hashing on one shared path. |

### Decision Rationale

The Questions-Options-Criteria matrix scores 1 (poor) through 5 (strong):

| Approach | Correctness fit | Prior-art alignment | Reversibility | Blast radius | Cost | Total |
| --- | --- | --- | --- | --- | --- | ---: |
| Content-derived semantic revision + resolver check | 5 — claimed identity is mechanically bound | 5 — content addressing and immutable values | 3 — locks a hash/pre-image contract | 4 — normalizer, resolver, tests | 3 — canonicalizer plus per-call check | **20** |
| Loader-issued revision + registry binding | 5 — registry can enforce one-to-one binding | 3 — version registries are established | 2 — introduces persistent authority | 1 — loader, registry, discovery, failure policy | 2 — lifecycle and lookup plumbing | **13** |
| Full normalized table in replay identity | 4 — replay stores the actual value | 2 — no compact stable identifier | 4 — avoids a hash contract | 3 — replay/audit schemas expand | 2 — large equality and storage surface | **15** |

The content-derived approach wins on the deciding correctness and prior-art
criteria while preserving the stateless resolver boundary. A loader registry
can bind friendly revisions, but adds ambient persistence, discovery, and
availability failures that RDR 0001 deliberately excludes. Recording the full
table makes the identity truthful but turns every replay and audit record into a
copy of the model and still needs canonical value equality.

Premortem: this ships and two behaviorally different tables receive one
revision because the encoder omitted a newly added edge field, or the producer
and resolver use subtly different order/framing rules. Replay then selects a
different rule without a mismatch. The recommendation survives: one shared
encoder, an explicit versioned field inventory, golden cross-boundary vectors,
and A4's adversarial MVV make omission or divergence a lock-blocking failure
rather than an implicit convention.

## Alternatives Considered

### Alternative 1: Loader-issued revision with registry binding

Have a table registry assign a human-readable or monotonic revision when a
normalized table is loaded, persist the one-to-one association, and require
replay to fetch the table by revision rather than supply rows independently.

**Pros**:

- Enforces one revision-to-table association without defining canonical hash
  bytes.
- Allows friendly revision names and explicit administrative retirement.

**Cons**:

- Adds persistent registry lifecycle, lookup availability, discovery, and
  concurrency policy to a stateless resolver design.
- Makes offline replay depend on retaining and reaching the issuing registry.

**Reason for rejection**: The correctness is strong, but the new stateful
authority has a much larger blast radius than the missing binding and conflicts
with RDR 0001's explicit-input replay boundary.

### Alternative 2: Full normalized table as replay identity

Remove compact `TableRevision` identity and persist/compare the entire
normalized outcome alphabet and candidate-row value in each replay record.

**Pros**:

- The replay record contains the exact semantic value used.
- Avoids hash algorithm and collision assumptions.

**Cons**:

- Duplicates the model in every replay/audit record and expands downstream wire
  contracts.
- Still requires stable semantic equality and ordering rules for comparison and
  diagnostics.

**Reason for rejection**: It solves aliasing at disproportionate storage and
API cost while losing the compact revision identity already established by RDR
0001.

### Briefly Rejected

- **Trust the caller's label**: This is the current defect; no ownership or
  validation rule prevents the same label from naming different rows.
- **Hash sparse TOML bytes**: Formatting, comments, and source key order would
  churn replay identity without changing resolver behavior.
- **Validate only during load**: Caller mutation or independently constructed
  resolver input could bypass the binding unless `Resolve` enforces it too.

## Trade-offs

### Consequences

- Replay revision becomes a mechanically checkable claim about the normalized
  match, action, and logical selected-rule identity used for selection.
- Semantically equivalent source formatting and author order share one revision;
  behavior- or logical-identity-changing normalized values do not. Relocating
  unchanged source preserves the revision.
- Resolver calls pay canonicalization and SHA-256 cost before matching, and
  normalized field evolution must update a versioned pre-image contract.
- RDR 0001's refusal taxonomy and RDR 0002's no-hash disposition are explicitly
  superseded at this seam.

### Risks and Mitigations

- **Risk**: A normalized field is omitted from the digest.
  **Mitigation**: Maintain an explicit field inventory, golden one-field-change
  vectors, and an exhaustiveness test beside the normalized type definitions.
- **Risk**: Order, nil, or empty-value handling makes equivalent tables hash
  differently across construction paths.
  **Mitigation**: Canonicalize typed values with explicit ordering and framing;
  lock cross-path golden vectors before Final.
- **Risk**: Per-resolve hashing becomes visible for larger tables.
  **Mitigation**: Benchmark after correctness is verified; an immutable cached
  digest is an allowed later optimization only if the adversarial mutation MVV
  still passes.
- **Risk**: Downstream callers treat mismatch as an ordinary no-match.
  **Mitigation**: Give it a distinct refusal kind and require RDR 0005 mapping to
  preserve that distinction.

### Failure Modes

- **Claimed revision does not identify supplied table**: resolution returns
  `table_revision_mismatch`, no plan, and no edge evaluation. Diagnose by
  deriving the table revision at the producer boundary and replacing the stale
  replay claim or restoring the intended table.
- **Malformed revision syntax**: loading/input construction fails before
  modeled resolution; correct the revision grammar rather than retrying edge
  selection.
- **Canonical encoder omits or aliases a field**: adversarial one-field-change
  or golden-vector tests fail. Treat as a spec/implementation defect and revise
  the version-1 field inventory before release.
- **Producer/resolver derivation drift**: a producer-issued revision immediately
  mismatches the same table at `Resolve`. Both boundaries must call the same
  derivation function; no compatibility fallback may trust either side.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified
- [ ] RDR 0002's complete normalized semantic field inventory is available.
- [ ] RDR 0001 and RDR 0005 refusal-taxonomy effects are reconciled.
- [ ] The version-1 pre-image grammar and golden vectors are locked.

### Minimum Viable Validation

`TestResolveRejectsTableRevisionMismatchBeforeSelection` constructs two valid
normalized tables that differ in one behavior-bearing edge field, derives the
revision for the first, and resolves both using that same claimed revision. The
first must reach its expected plan; the second must return only
`table_revision_mismatch` before any ordinary or escape edge is evaluated. A
reordered/reformatted construction of the first table must derive the same
revision and the same plan.

### Phase 1: Lock Canonical Revision Semantics

Specify and verify the version-1 normalized field inventory, ordering, framing,
domain separator, revision grammar, and golden vectors.

### Phase 2: Bind Producer Output

Make normalization/loading derive revision with the shared function and return
it with the complete normalized semantic table value.

### Phase 3: Enforce Resolver Input

Validate the claimed revision before existing modeled input checks and add the
distinct mismatch refusal without changing exact-one selection after equality.

### Phase 4: Prove Replay and Downstream Mapping

Add adversarial same-revision/different-table coverage and preserve the new
refusal through the planned CLI mapping and audit fixtures.

## Validation

### Testing Strategy

1. **Scenario**: Normalize tables that differ only in source formatting, source
   key order, candidate input order, map iteration, or `SourceLocator`.
   **Expected**: Every table produces the same version-1 pre-image, revision,
   and resolver plan.
2. **Scenario**: Change each match, action, outcome, row-kind, or logical
   selected-rule identity field individually.
   **Expected**: Every change produces a different golden pre-image and revision;
   the projection inventory contains no unclassified semantic field.
3. **Scenario**: Resolve a table using the revision derived by the producer.
   **Expected**: Binding succeeds and the existing exact-one resolution result
   is unchanged.
4. **Scenario**: Resolve an altered table using the original table's revision.
   **Expected**: Only `table_revision_mismatch` is returned, with no plan and no
   owned-state, outcome, guard, ordinary-edge, or escape-edge evaluation.
5. **Scenario**: Supply malformed revision syntax or an invalid normalized
   value.
   **Expected**: The Go load/programmer error path is used; neither case is
   reported as a modeled mismatch.

## Finalization Gate

### Contradiction Check

Pending Stage 7 response after assumption verification and pre-lock review.

### Assumption Verification

Pending Stage 7 response after A1–A5 reach terminal dispositions.

### Scope Verification

Pending Stage 7 response against
`TestResolveRejectsTableRevisionMismatchBeforeSelection` and the version-1
golden vectors.

### Cross-Cutting Concerns

Pending Stage 7 response for versioning, character encoding, canonical form,
determinism, mutation safety, and downstream CLI mapping.

### Proportionality

Pending Stage 7 response revalidating the single revision-binding contract and
the `foundational` profile.

## References

- RDR 0001 `Normative Contracts`, `Load-Bearing Decisions / Identity`, and
  `Cross-Cutting Concerns / Canonical-form` — replay tuple and explicit no-hash
  predecessor contract.
- RDR 0002 `Normative Contracts`, `Load-Bearing Decisions / Identity`, and
  `Cross-Cutting Concerns / Canonical-form` — normalized semantic value and
  explicit exclusion of content-addressed identity.
- RDR 0008 `Technical Environment` — selected-rule handoff is adjacent and
  explicitly does not own table revision.
- `internal/resolver/resolver.go::Input`, `Resolve`, `matchingEdges`, and
  `RefusalKind`; replay and adversarial resolver tests.
- Pro Git, chapter 10.2, *Git Internals — Git Objects*.
- Eric Evans, *Domain-Driven Design*, chapter 5, *A Model Expressed in Software*.
- Srinivasan, *A Methodology for Selecting and Composing Runtime Architecture
  Patterns for Production LLM Agents*, pages 7 and 11.
- `../state-machines/repos/README.md`, *No new clone needed — the resolver is a
  table + CLI*.
- Kata jv85, *Define normalized-table revision binding for resolver replay*.
