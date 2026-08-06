# Recommendation 0007: Normalized-table revision identity

## Metadata

- **Date**: 2026-08-06
- **Status**: Draft
- **Type**: Architecture
- **Profile**: foundational — defines one revision-identity contract consumed
  across normalization, replay, and resolver enforcement.
- **Priority**: Medium
- **Related Issues**: kata jv85; RDR 0001; RDR 0002
- **Predecessors**: 0001-resolution-kernel,
  0002-transition-table-as-reviewable-data
- **Overrides**: RDR 0002's explicit exclusion of content-addressed table
  identity
- **Seam Lineage**: area:resolver — no prior accretion

## Problem Statement

Developers replaying resolver decisions need a compact revision that identifies
one normalized transition-table value, so tests and audits can distinguish
behaviorally different tables without depending on source formatting or row
order. They discover the gap when the same nominal revision can name different
normalized contents. The system therefore needs one canonical,
content-derived identity contract for normalized tables.

## Context

### Background

RDR 0001 includes transition-table revision in replay identity, but RDR 0002
deliberately leaves content-addressed normalized-table identity unspecified. A
buggy or adversarial caller can therefore reuse a revision string with altered
normalized rows, undermining deterministic replay and audit claims. This RDR
defines the identity that closes that naming gap. Resolver comparison,
mismatch classification and ordering, and downstream CLI mapping are a separate
enforcement contract and require a follow-up RDR.

### Technical Environment

The affected seam is the loading or normalization boundary that produces the
normalized transition-table value consumed by
`internal/resolver/resolver.go::Resolve`. RDR 0001 defines the replay tuple and
RDR 0002 defines reviewable transition data and normalized row semantics. This
RDR defines only the canonical public identity of that normalized value; it
does not change resolver selection or the CLI output contract.

## Research Findings

### Investigation

RDRs 0001 and 0002 are Final, kata jv85 tracks this gap, and
`internal/resolver/resolver.go::Input` carries `TableRevision` beside caller-
supplied `Outcomes` and `Table`. `Resolve` does not read `TableRevision`, and
current replay tests reuse the same table value. Draft RDR 0008 owns
selected-rule identity after selection and leaves table revision identity here.

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
- **Verified** — the canonicalization spike produces the same pre-image
  and revision after semantic-set reordering, map iteration, source-locator
  changes, source-schema-version changes, and nil/empty construction changes;
  every projected one-field mutation and clear-versus-empty boundary changes
  both.
- **Verified** — Go's `crypto/sha256::Sum256` and
  `encoding/hex::EncodeToString` provide the intended standard-library digest
  and lowercase 64-digit encoding without a third-party dependency. Revision
  parsing must still reject uppercase explicitly because hex decoding accepts
  it.
- **Pending** — the version-1 projection now includes explicit dispositions for
  row outcome, predicate provenance/name/operator, write role, and guard
  unevaluability, but Resolve must verify the inventory against RDR 0002 and
  the concrete normalized value.
- **Documented** — resolver mismatch classification and ordering alter RDR
  0001's closed refusal policy independently of canonical identity; they are
  outside this RDR and must be specified before resolver enforcement ships.

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
  - **Evidence**: RDR 0008 `Load-Bearing Decisions / Identity` makes
    `SourceLocator` diagnostic-only. The explicit inventory in Technical Design
    covers row outcome, predicate provenance/name/operator, write role, and
    guard unevaluability; Resolve must audit it against the RDR 0002 and concrete
    normalized values and extend the one-field-change vectors.
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
- **A4 The normalizer and later enforcement consumers can share one pure
  revision function without exposing mutable digest state or a second
  implementation.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: The peer/source audit found no existing production digest or
    binding implementation. Resolve must identify the one exported derivation
    symbol and the normalization call site; the follow-up enforcement RDR must
    require its consumers to call that symbol rather than reimplement it.
  - **If wrong**: Producer and consumer may accept different revisions for the
    same table, or caller mutation may invalidate a cached binding silently.

## Proposed Solution

### Approach

Make transition-table revision a content-derived identity of the normalized
semantic table. The normalization/loading boundary derives the revision through
one exported, versioned SHA-256 function over the normalized outcome alphabet
and candidate rows. Replay records store that derived value as
`TableRevision`; any boundary that compares a table to a claimed revision must
reuse the same derivation function under a separate enforcement contract.

The digest covers the normalized semantic value, not TOML bytes or a rendered
dump. It binds every match, action, and logical selected-rule identity field,
but excludes `SourceLocator`: RDR 0008 defines that field as diagnostic
provenance rather than logical rule identity. Its external grammar is
algorithm-qualified lowercase hexadecimal, initially `sha256:<64-hex>`; the
pre-image is domain-separated and versioned so later semantic formats cannot
accidentally share an identity namespace.

Revision parsing accepts only the locked public grammar. Failure to derive a
revision from an invalid normalized value, or failure to parse revision syntax,
uses the load/programmer error path. This RDR does not choose how a resolver
reports or orders a comparison failure.

### Technical Design

The normalized-table producer owns revision derivation because it owns the
semantic value being identified. It uses one package-level canonical encoder
and exported digest function; enforcement and CLI code must not implement a
second hash path. RDR 0002 remains the owner of normalized values' semantics
and source identity. This RDR owns only their content-derived revision.

Canonical encoding operates over a version-1 semantic projection, not Go field
layout. The pre-image begins with the ASCII bytes
`intrastate.transition-table-revision`, a NUL byte, and the ASCII encoding label
`v1`. Every projected value is framed as a one-byte type tag, an unsigned
64-bit big-endian payload length, and the payload: `0x01` is a UTF-8 string,
`0x02` an unsigned 64-bit big-endian integer, `0x03` a list, `0x04` a record,
and `0x05` a Boolean encoded as the single byte `0x00` or `0x01`. A list
payload begins with an unsigned 64-bit element count followed by framed
elements. A record payload begins with a field count followed by pairs of
framed field-name strings and framed values, sorted by field-name bytes.

The top-level record fields are `model_id`, `outcomes`, and `rows`; RDR 0002's
source-schema version is excluded because it is not part of the normalized
semantic value. Each row record contains `rule_id`, `expansion_suffix`, `kind`,
`outcome`, `predicates`, `escape_classes`, `guard_unevaluable`, `next_tags`, and
`writes`. Predicate entries are records of `provenance`, `name`, `operator`, and
`value`, so owned/observed/recognized lookup, fixed predicate operation, and
positive/negative semantics cannot alias. Next-tag entries are records of `tag`
and `value`. Write entries are records of `role`, `tag`, `operation`, and, for
`set`, `value`; `operation` is either `set` or `clear`.
`guard_unevaluable` is included while it remains a normalized behavior-bearing
input; removing it requires the normalized type to eliminate that behavior,
not merely omit it from the projection. Outcomes, rows, predicates, escape
classes, next tags, and writes are semantic sets and are sorted by their
complete framed element bytes before list framing.

`SourceLocator` and every other source-formatting or diagnostic-only field are
excluded. Strings are their validated UTF-8 bytes without additional Unicode
normalization. The explicit `operation` field distinguishes clear from setting
an empty string; authored absent and present-empty forms that normalize to the
same typed value cannot produce different revisions. Resolve must verify this
projection against the concrete RDR 0002/RDR 0008 types and produce golden
vectors before the contract locks; it may not change these inclusion,
exclusion, ordering, framing, or primitive-encoding rules.

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
versioned canonical encoding. One exported derivation function MUST be used by
the normalization/loading boundary and every later identity consumer.

Canonical encoding MUST use the version-1 domain prefix, recursively typed
length framing, semantic field projection, and encoded-byte sorting defined in
Technical Design. It MUST be independent of TOML whitespace,
comments, source key order, map iteration, normalized candidate input order,
and diagnostic source location. The version-1 projection MUST include row
outcome, predicate provenance/name/operator, write role, and any normalized
guard-unevaluability state in addition to the other fields listed in Technical
Design. The projection-to-normalized-type inventory and golden vectors MUST be
verified before Final.

Revision parsing MUST reject any value outside `sha256:<64 lowercase
hexadecimal digits>`. Malformed revision syntax and failure to construct a
valid normalized semantic table MUST use the load/programmer error path.
Resolver mismatch disposition and evaluation ordering are outside this RDR.
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
  normalized semantic value. Resolve verifies the projection and golden
  vectors; it does not choose different fields, ordering, primitives, or
  framing.
- **Naming** — the content identity remains `TableRevision` in replay input;
  rejected: `version`, which RDR 0002 already uses for source-schema version,
  and `etag`, which implies transport cache semantics.
- **Selection / predicate** — every behavior-bearing predicate component is
  part of the semantic projection; this RDR does not alter predicate evaluation
  or resolver selection ordering.

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Canonical normalized semantic table | RDR 0002 | Predecessor | Supplies the value whose behavior and identity fields are hashed. |
| Logical selected-rule identity fields | RDR 0002 / RDR 0008 | Deferred peer | Included in the table pre-image; diagnostic source location is excluded. |
| SHA-256 and hexadecimal encoding | Go standard library | Available | A3 verifies `crypto/sha256::Sum256` and lowercase `encoding/hex` output; grammar validation remains explicit. |
| Table revision derivation | This RDR | Introduced | One exported function gives producer output its canonical replay identity. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Replay revision claim | `internal/resolver/resolver.go::Input.TableRevision` | Unconstrained caller label | Preserve consumer | A follow-up enforcement RDR binds this claim to the derived identity. |
| Normalized table input | `internal/resolver/resolver.go::Input.Outcomes` and `Table` | No canonical identity | Consume | Derivation projects the complete normalized semantic value without owning its type shape. |
| Revision/digest helper | `internal/` sibling-path search | None exists | Introduce | Keep canonicalization and hashing on one shared path. |

### Decision Rationale

The Questions-Options-Criteria matrix scores 1 (poor) through 5 (strong):

| Approach | Correctness fit | Prior-art alignment | Reversibility | Blast radius | Cost | Total |
| --- | --- | --- | --- | --- | --- | ---: |
| Content-derived semantic revision | 5 — identity is mechanically derived | 5 — content addressing and immutable values | 3 — locks a hash/pre-image contract | 4 — normalizer and replay consumers | 3 — canonicalizer and golden vectors | **20** |
| Loader-issued revision + registry binding | 5 — registry can enforce one-to-one binding | 3 — version registries are established | 2 — introduces persistent authority | 1 — loader, registry, discovery, failure policy | 2 — lifecycle and lookup plumbing | **13** |
| Full normalized table in replay identity | 4 — replay stores the actual value | 2 — no compact stable identifier | 4 — avoids a hash contract | 3 — replay/audit schemas expand | 2 — large equality and storage surface | **15** |

The content-derived approach wins on the deciding correctness and prior-art
criteria while preserving the stateless resolver boundary. A loader registry
can bind friendly revisions, but adds ambient persistence, discovery, and
availability failures that RDR 0001 deliberately excludes. Recording the full
table makes the identity truthful but turns every replay and audit record into a
copy of the model and still needs canonical value equality.

Premortem: this ships and two behaviorally different tables receive one
revision because the encoder omitted a newly added edge field, or a consumer
reimplemented the framing rules. Replay can then treat different tables as one
identity. The recommendation survives: one exported encoder, an explicit
versioned field inventory, golden cross-boundary vectors, and A4's call-site
audit make omission or divergence a lock-blocking failure.

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
  derivation rule prevents the same label from naming different rows.
- **Hash sparse TOML bytes**: Formatting, comments, and source key order would
  churn replay identity without changing resolver behavior.

## Trade-offs

### Consequences

- Replay revision becomes a mechanically derived identity of the normalized
  match, action, and logical selected-rule value.
- Semantically equivalent source formatting and author order share one revision;
  behavior- or logical-identity-changing normalized values do not. Relocating
  unchanged source preserves the revision.
- Identity derivation pays canonicalization and SHA-256 cost, and normalized
  field evolution must update a versioned pre-image contract.
- RDR 0002's no-hash disposition is explicitly superseded at this seam.

### Risks and Mitigations

- **Risk**: A normalized field is omitted from the digest.
  **Mitigation**: Maintain an explicit field inventory, golden one-field-change
  vectors, and an exhaustiveness test beside the normalized type definitions.
- **Risk**: Order, nil, or empty-value handling makes equivalent tables hash
  differently across construction paths.
  **Mitigation**: Canonicalize typed values with explicit ordering and framing;
  lock cross-path golden vectors before Final.
- **Risk**: Repeated derivation becomes visible for larger tables.
  **Mitigation**: Benchmark after correctness is verified; an immutable cached
  digest is an allowed later optimization only when the normalized value is
  immutable and derivation-equivalence tests still pass.

### Failure Modes

- **Malformed revision syntax**: parsing fails on the load/programmer error
  path; correct the revision grammar rather than accepting a non-canonical
  alias.
- **Canonical encoder omits or aliases a field**: adversarial one-field-change
  or golden-vector tests fail. Treat as a spec/implementation defect and revise
  the version-1 field inventory before release.
- **Producer/consumer derivation drift**: cross-boundary golden vectors or the
  call-site audit fail. Every boundary must call the exported derivation
  function; no compatibility path may implement its own encoder.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified
- [ ] RDR 0002's complete normalized semantic field inventory is available.
- [ ] The version-1 pre-image grammar and golden vectors are locked.

### Minimum Viable Validation

`TestNormalizedTableRevisionCanonicalIdentity` derives revisions for one valid
normalized table, semantically equivalent reordered/reformatted constructions,
and one-field mutations spanning every projected field. Equivalent values must
produce the same pre-image and revision; every behavior- or logical-identity-
changing mutation must produce a different pre-image and revision. The test
also rejects uppercase, wrong-length, wrong-prefix, and non-hex revision text.

### Phase 1: Lock Canonical Revision Semantics

Specify and verify the version-1 normalized field inventory, ordering, framing,
domain separator, revision grammar, and golden vectors.

### Phase 2: Bind Producer Output

Make normalization/loading derive revision with the exported function and
return it beside the complete normalized semantic table value. Implement the
strict public revision parser on the same identity path.

## Validation

### Testing Strategy

1. **Scenario**: Normalize tables that differ only in source formatting, source
   key order, candidate input order, map iteration, or `SourceLocator`.
   **Expected**: Every table produces the same version-1 pre-image and revision.
2. **Scenario**: Change each outcome, row kind, predicate provenance/name/
   operator/value, escape class, guard-unevaluability value, next tag, write
   role/tag/operation/value, or logical selected-rule identity field
   individually.
   **Expected**: Every change produces a different golden pre-image and revision;
   the projection inventory contains no unclassified semantic field.
3. **Scenario**: Supply uppercase, wrong-length, wrong-prefix, non-hex revision
   text, or an invalid normalized
   value.
   **Expected**: The load/programmer error path rejects it.
4. **Scenario**: Audit normalization and identity consumers.
   **Expected**: They call the one exported derivation symbol; no second encoder
   or digest implementation exists.

## Finalization Gate

### Contradiction Check

Pending Stage 7 response after assumption verification and pre-lock review.

### Assumption Verification

Pending Stage 7 response after A1–A4 reach terminal dispositions.

### Scope Verification

Pending Stage 7 response against
`TestNormalizedTableRevisionCanonicalIdentity` and the version-1 golden
vectors.

### Cross-Cutting Concerns

Pending Stage 7 response for versioning, character encoding, canonical form,
determinism, and mutation safety.

### Proportionality

Pending Stage 7 response revalidating the single revision-identity contract and
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
- `internal/resolver/resolver.go::Input`, `Edge`, `TagPredicate`, `Guard`, and
  `OwnedTagWrite`; normalized semantic inputs and replay revision claim.
- Pro Git, chapter 10.2, *Git Internals — Git Objects*.
- Eric Evans, *Domain-Driven Design*, chapter 5, *A Model Expressed in Software*.
- Srinivasan, *A Methodology for Selecting and Composing Runtime Architecture
  Patterns for Production LLM Agents*, pages 7 and 11.
- `../state-machines/repos/README.md`, *No new clone needed — the resolver is a
  table + CLI*.
- Kata jv85, *Define normalized-table revision binding for resolver replay*.
