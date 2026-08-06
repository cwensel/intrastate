# Recommendation 0007: Normalized-table revision binding

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). **Reference-only** (guidance; never copy into the
instance body). -->

## Metadata

- **Date**: 2026-08-06
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
- **Profile**: foundational — defines one revision-binding contract across
  normalized-table production and resolver consumption.
  <!-- Do not paste the matrix below into the field; it is the
  Stage 5 routing latch, provisional on `Draft`, made
  authoritative by Resolve.
  Sized by BLAST RADIUS — the MAX of two axes, not
  contract count or word count.
  (1) contract axis: small = one contract, no user-facing
  surface (skips Stage 5); mid = one contract + user-facing
  surface OR locks a contract; large = locks an enum/hash/
  format/grammar/destructive-op; foundational = cross-RDR
  producer / spans modules.
  (2) accretion axis (HARD floor): if `Seam Lineage` below
  carries ≥2 closed prior point-fixes at this locus, Profile
  is floored at FOUNDATIONAL regardless of the contract axis
  — a seam with prior point-fixes is never small/mid (it
  spans the prior RDRs/patches = the matrix's cross-RDR
  trigger). The only escape is a written accretion disposition
  in the Seam Lineage field. This floor is what stops a
  "one contract → mid" sizing from under-gating an accreting
  seam.
  Matrix: rdr/stages/README.md. Seed estimates from the design
  shape; Resolve overwrites from the verified count; Stage 8
  Gate locks it at Draft → Final. Never skip lenses off a
  Draft Profile until Resolve has run. -->
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
the transition table as reviewable data. Ownership of revision binding and the
failure boundary remain undecided; later stages must evaluate the available
design forks without changing the resolver's established output contracts
accidentally.

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
- **Assumed** — the RDR 0002 normalized semantic value admits one versioned,
  order-independent encoding that covers every behavior- and identity-bearing
  field without depending on TOML source formatting.

### Critical Assumptions

[Required — never omit. Load-bearing assumptions — if
wrong, the approach fails. Each must have a complete
Evidence Record before marking this RDR Final.]

- **A1 The normalized semantic table can be encoded canonically without source
  formatting, map iteration, or caller row order changing its revision.**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: Build the Stage 4 canonicalization spike from RDR 0002's
    normalized candidate-row value and compare reordered, reformatted, and
    one-field-different fixtures.
  - **If wrong**: Semantically identical tables churn revision or distinct
    resolver behavior aliases to one revision, defeating reliable replay.
- **A2 Every field that can change edge selection or selected-rule identity can
  be included in one versioned digest pre-image.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: Audit RDR 0002 `Normative Contracts` and RDR 0008's eventual
    selected-rule carrier against the Stage 4 pre-image field inventory.
  - **If wrong**: A caller can alter an omitted outcome, predicate, action, or
    rule identity while retaining the same revision.
- **A3 Go's standard SHA-256 implementation and lowercase hexadecimal encoding
  are available and sufficient for this table-identity boundary.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: Verify `crypto/sha256::New` or `Sum256` and `encoding/hex`
    against the exact framing chosen at Resolve; no third-party dependency is
    permitted for the hash.
  - **If wrong**: The format needs another algorithm or dependency before its
    revision grammar can be locked.
- **A4 The normalizer and resolver can share one revision function without
  exposing mutable digest state or a second implementation.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: `TestResolveRejectsTableRevisionMismatchBeforeSelection` will
    derive a valid revision, alter one normalized edge, and require the shared
    verifier to refuse the altered table before matching it.
  - **If wrong**: Producer and consumer may accept different revisions for the
    same table, or caller mutation may invalidate a cached binding silently.
- **A5 A table-revision mismatch is a modeled input refusal that can extend RDR
  0001's closed taxonomy without using the Go error path.**
  - **Status**: Pending
  - **Method**: Design Decision
  - **Evidence**: Resolve must reconcile this RDR's `table_revision_mismatch`
    override with RDR 0001 `Normative Contracts` and the downstream RDR 0005
    mapping before lock.
  - **If wrong**: Callers cannot diagnose a replay-identity failure uniformly,
    and may conflate adversarial input with parser or programmer failure.

**Method vocabulary** (Reference-only — guidance for
filling the **Method** field; do NOT copy this list into
the instance body. Pick exactly one per assumption):

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

Make transition-table revision a content-derived identity of the normalized
semantic table. The normalization/loading boundary derives a revision using one
shared, versioned SHA-256 function and returns the revision with the normalized
outcome alphabet and candidate rows as one table value. A replay caller records
that revision as `Input.TableRevision`; `Resolve` recomputes the supplied table's
revision with the same function and compares it before evaluating owned state,
outcomes, guards, or edges.

The digest covers the normalized semantic value, not TOML bytes or a rendered
dump. It therefore binds every behavior- and selected-identity-bearing field
while ignoring source formatting and authored key order. Its external grammar
is algorithm-qualified lowercase hexadecimal, initially `sha256:<64-hex>`; the
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

Canonical encoding operates over typed normalized values, with explicit field
tags and length framing, deterministic ordering for set/map-shaped fields, and
distinct encodings for empty, absent, and present-empty values. The pre-image
begins with a table-identity domain and encoding version. TOML whitespace,
comments, source key order, absolute paths, and rendered dump formatting never
enter the digest. Stage 4 must finish the field inventory and framing spike
before this contract locks.

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
alphabet, candidate selection data, transition actions, and normalized rule
identity/provenance fields. It MUST NOT be an unconstrained caller label, TOML
source checksum, rendered-dump checksum, registry sequence, timestamp, or source
path.

The revision grammar MUST be algorithm-qualified lowercase hexadecimal. Version
1 MUST use `sha256:<64 lowercase hexadecimal digits>` over a domain-separated,
versioned canonical encoding. One shared derivation function MUST be used by the
normalization/loading boundary and resolver validation.

Canonical encoding MUST be independent of TOML whitespace, comments, source key
order, map iteration, and normalized candidate input order. It MUST use explicit
field identity and framing, deterministic ordering for set/map-shaped values,
and distinct representations for absent, empty, and present-empty values. Its
exact version-1 field inventory and framing MUST be locked before Final.

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
  formatting and candidate input order do not distinguish identity; any field
  that can change resolution or selected-rule identity does.
- **Wire / byte format** — the public revision is
  `sha256:<64 lowercase hexadecimal digits>`. The pre-image is a domain-
  separated, versioned, length-framed encoding of the normalized semantic
  value; Stage 4 owns the exact field table and golden vectors required before
  lock.
- **Naming** — the content identity remains `TableRevision` in replay input;
  the mismatch is `table_revision_mismatch`. Rejected: `version`, which RDR 0002
  already uses for source-schema version, and `etag`, which implies transport
  cache semantics.
- **Selection / predicate** — revision equality is a mandatory gate before any
  candidate qualifies. A mismatch selects no ordinary or escape row and cannot
  itself be overridden by a table-authored escape.

### Capability Dependencies

[For each load-bearing behavior, state whether the
enabling capability exists now, is introduced by this
RDR, is provided by a predecessor, or is deferred.]

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Exact-one resolver selection | RDR 0001 | Available | Runs only after revision binding succeeds. |
| Canonical normalized semantic table | RDR 0002 | Predecessor | Supplies the value whose behavior and identity fields are hashed. |
| Selected-rule identity fields | RDR 0002 / RDR 0008 | Deferred peer | Included in the table pre-image without moving handoff ownership here. |
| SHA-256 and hexadecimal encoding | Go standard library | Pending verification | No third-party runtime dependency is intended. |
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

- Replay revision becomes a mechanically checkable claim about the exact
  normalized semantics used for selection and selected-rule identity.
- Semantically equivalent source formatting and author order share one revision;
  behavior- or identity-changing normalized values do not.
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
