# RDR 0007 prior-art selection record

## Accepted anchors

- RDR 0001 `Normative Contracts`: "Given the same flow identity, transition
  table revision ... resolve returns the same disposition." Its
  `Cross-Cutting Concerns / Canonical-form` explicitly limits the predecessor to
  value-level replay and introduces no hash or canonical serialization.
- RDR 0002 `Normative Contracts` makes the normalized candidate-row value the
  source for resolver lookup and deterministic dumps. Its `Cross-Cutting
  Concerns / Canonical-form` calls that value canonical while explicitly saying
  the spike SHA is not a content-addressed contract.
- Pro Git, chapter 10.2, page 483: `hash-object` produces a key from a
  "checksum of the content you're storing plus a header." This supports a
  domain-separated content identity rather than a free revision label.
- Eric Evans, *Domain-Driven Design*, chapter 5, page 56: for shared value
  safety, an object "must be immutable: it cannot be changed except by full
  replacement." This supports binding revision and normalized semantics as one
  value and defending against post-derivation mutation.
- Srinivasan, *A Methodology for Selecting and Composing Runtime Architecture
  Patterns for Production LLM Agents*, page 7: under model or prompt revision,
  the same event can produce different replay output. This supports making the
  deciding model revision mechanically meaningful at the deterministic
  boundary.
- `../state-machines/repos/README.md`, *No new clone needed — the resolver is a
  table + CLI*: the runtime is a declarative table and stateless resolver, which
  weighs against adding a table registry to this seam.

## Queries

1. Semble, intrastate: `resolver Resolve function uses normalized transition
   table revision replay identity and selects transition edges`.
2. Semble, intrastate: `normalized transition table creation loading validation
   revision identifier content digest canonical rows`.
3. Semble, state-machines: `transition table revision identity content hash
   canonical normalized table resolver replay`.
4. Semble, state-machines: `resolver design stateless transition table model
   version immutable revision provenance replay deterministic same table`.
5. Arc, `DevRef`: `content-addressed object identifier is cryptographic hash of
   canonical object contents verification detects identifier content mismatch`.
6. Arc, `StateMachineLit`: `state machine transition relation version revision
   identity deterministic replay uses same model configuration`.
7. Arc, `DevRef`: `immutable value object bundles version identity with content
   prevent mismatched revision and data at consumer boundary`.

## Rejected branches

- The configured `OpenSource` and `DevRefOS` corpus names recorded by the RDR
  resource index no longer exist; their attempted queries produced no evidence.
  Current `DevRef` and `StateMachineLit` results were used instead.
- The `StateMachineLit` search returned replay-divergence evidence but no
  load-bearing canonical byte layout; the exact pre-image remains a Pending
  Stage 4 assumption rather than being inferred from the paper.
- A loader-issued registry revision was not adopted because it introduces
  persistent authority, discovery, and availability semantics into RDR 0001's
  explicit-input stateless resolver boundary.
- Persisting the full normalized table as replay identity was not adopted
  because it duplicates model data in replay/audit records and still requires a
  stable semantic equality contract.
- Hashing sparse TOML bytes was rejected because whitespace, comments, and
  source ordering are not resolver semantics.

## Stage 4 resolve audit

### Accepted evidence

- `GOCACHE=/tmp/intrastate-rdr7-go-cache go run
  ./docs/rdr/0007-normalized-table-revision-binding/evidence/spikes` was rerun
  on 2026-08-06 and byte-matched `../spikes/output.txt`. The live spike verifies
  the proposed framing and ordering mechanics:
  reordered semantic sets, map iteration, diagnostic locator changes,
  source-schema-version changes, and nil/empty construction preserve the
  pre-image and revision; projected one-field changes, clear-versus-empty, and
  length-framing boundaries do not.
- Go 1.26 `crypto/sha256::Sum256` supplies the standardized 32-byte digest, and
  `encoding/hex::EncodeToString` emits 64 lowercase hexadecimal digits. The
  decoder accepts uppercase, so revision grammar validation cannot rely on
  successful hex decoding alone.
- RDR 0002 `Normative Contracts` and `Load-Bearing Decisions / Identity` define
  the normalized row's model/rule identity, expansion suffix, row kind,
  predicates, writes, escape classes, and diagnostic source locator. RDR 0008
  `Load-Bearing Decisions / Identity` confirms that `SourceLocator` is not
  logical selected-rule identity.
- RDR 0003 `Technical Design` and `Load-Bearing Decisions / Identity` define a
  predicate's behavior-bearing typed literal and distinguish positive `all`
  predicates from one conjunctive `unless` block.

### Queries and negative findings

1. Semble, intrastate: `existing implementation for transition table revision
   canonical encoding digest SHA-256 validation and resolver binding mismatch
   refusal`.
2. Semble, intrastate: `normalized transition table concrete Go types candidate
   rows actions source locator model id`.
3. Semble, intrastate: `resolver refusal kind modeled value and CLI mapping
   structured output error path`.
4. Semble, intrastate: `tests for resolver refusal handling and transition
   table revision replay identity`.
5. No production normalized-table constructor, revision digest, parser, or
   binding call site exists to reuse. `internal/resolver/resolver.go::Input`
   contains the only `TableRevision`, and `Resolve` does not consult it. RDR
   0002's evidence contains only a normalization/sorting prototype.

### Blocking inventory finding

The current version-1 projection remains insufficient evidence for A2. It now
names `Edge.Outcome`, predicate provenance/name/operator, write role, and guard
unevaluability, and it excludes `SourceLocator` consistently with RDR 0008.
However, it does not encode RDR 0003's typed predicate literal shape or the
semantic distinction between positive `all` predicates and the conjunctive
`unless` block. Production also has no validated normalized-table value against
which to close the inventory. Treating the projection as complete was rejected;
Refine must settle those dimensions before Stage 4 extends and reruns the
one-field vectors and peer audit.

### Profile recount

The Normative Contracts currently lock at least two independent load-bearing
contracts: the canonical semantic pre-image/SHA-256 identity and the public
revision grammar/parser contract. Resolver mismatch disposition and ordering
have correctly moved outside this RDR, but the remaining count is still a
Resolve split signal. Stage 4 therefore did not overwrite the provisional
`foundational` profile. Refine must either demonstrate that these are one
inseparable contract or split them before the routing latch can be set.

### Stage 4 cold recheck (2026-08-06)

- **A1 — verified**: `GOCACHE=/tmp/intrastate-rdr7-go-cache-a1 go run
  ./docs/rdr/0007-normalized-table-revision-binding/evidence/spikes` reproduced
  `../spikes/output.txt`. The spike's `canonicalPreimage` and `revision`
  functions preserve identity across semantic-set order, map iteration,
  diagnostic-locator/source-schema changes, and nil/empty construction; its
  projected mutations and clear-versus-empty boundary change the identity.
- **A3 — verified**: Go 1.26
  `crypto/sha256/sha256.go::{Size,Sum256}` supplies a 32-byte SHA-256 result;
  `encoding/hex/hex.go::{hextable,EncodedLen,EncodeToString}` produces its
  64-character lowercase form. `reverseHexTable` accepts uppercase, so the
  parser must check lowercase bytes before decoding.
- **A2 and A4 — incomplete**: Semble searches for the normalized-table
  constructor, validation boundary, canonical encoder, digest helper, and
  revision binding find only this RDR's spike and
  `internal/resolver/resolver.go::{Input,Resolve}`. The production input still
  accepts raw `[]Edge` plus an unused caller `TableRevision`; there is no
  normalized-table, `SelectedRule`, action, or typed-`SourceLocator` inventory
  against which to prove completeness or shared-consumer reuse. RDR 0002/0008
  must first make those concrete types and their construction boundary
  available, then this RDR must audit the projection and extend the one-field
  mutation vectors.
