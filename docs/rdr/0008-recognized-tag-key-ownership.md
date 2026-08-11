# Recommendation 0008: Ownership of the recognized-outcome tag key name

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). **Reference-only** (guidance; never copy into the
instance body). -->

## Metadata

- **Date**: 2026-08-09
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
- **Profile**: foundational — one contract: who owns the
  *name* of the recognized-provenance tag key in the
  assembled view. It is a cross-RDR producer, binding RDR
  0001's kernel carrier to RDR 0002's declared model and
  extending 0002's validation surface.
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
- **Related Issues**: kata `z53t` — "Bind the
  recognized-outcome tag key between RDR 0001's kernel view
  and RDR 0002's tag declarations" (`src:roborev`,
  `severity:medium`, `area:internal-resolve`,
  `release:post-1.0`).
- **Predecessors**: 0001-resolution-kernel,
  0002-transition-table-as-reviewable-data
- **Overrides**: Ratifies RDR 0001 deviation D2 (the
  `recognized` key) from implementation latitude to
  normative contract; extends RDR 0002's tag-declaration
  surface with a name constraint on recognized-provenance
  declarations and one additional data-level validation
  category (`reserved tag key`, within its "including at
  minimum" extensible list). Narrows nothing in either
  peer.
- **Seam Lineage**: `area:internal-resolve` (recognized-tag
  key identity at the RDR 0001 ↔ 0002 boundary) — no prior
  accretion.

## Problem Statement

Someone authoring a transition table wants a row to match on
the outcome a recognizer just produced, using the same tag
vocabulary they use for every other tag. They name their
tags themselves — that is the point of the declared model —
so they name the recognized-provenance one whatever reads
best in their flow. They discover the problem when that row
never fires: the outcome was recognized, the value is right,
and the match still fails, because the key the kernel wrote
it under is not the key they declared. Nothing in the
refusal tells them the two names had to agree, or which name
was supposed to win.

Internally, this RDR must fix **who owns the name** of the
recognized-provenance tag key in the assembled evaluation
view. RDR 0001's kernel asserts a fixed key —
`recognizedTagKey = "recognized"`, an unexported constant at
`internal/resolve/resolve.go::recognizedTagKey` — while RDR
0002 grants table authors naming authority over every
declared tag, including recognized-provenance ones. Neither
RDR states which side is authoritative, so a normalizer and
a kernel can each hold a defensible reading that does not
compose.

The decision is a single fork with three defensible
branches: `"recognized"` is a **reserved kernel keyword**
normalization must conform to (carrier owns identity); the
key is a **declaration the model carries** and the kernel
binds at assembly (model owns identity, carrier
parameterizes); or the **tag-key affordance is withdrawn**
so recognized outcomes are matchable only via `Row.Outcome`
(the identity question dissolves). The choice determines
whether the normalized table gains public surface for this
key. The shadowing rule — whether an owned or observed tag
may collide with that key — is a *consequence* of the
branch, not a separate decision: a breach under the first, a
lint obligation under the second, moot under the third.

## Context

### Background

Raised by a roborev finding against the RDR 0001
implementation (kata `z53t`, `src:roborev`), then re-grounded
during seed triage.

Not reachable at HEAD, and not a correctness break — the
tag-key path is a secondary affordance (Key Discoveries).
RDR 0002 is Final but unimplemented, so no normalizer exists
and no cross-RDR call path can break today. The worst case
at integration time is a silent no-match refusal — typed,
deterministic, diagnosable — not a wrong edge. RDR 0001's
`verification.md` already records the unreserved key as an
explicit non-defect: "deterministic and refusal-shaped, not
a breach."

Constraints inherited from prior rounds. RDR 0001's
deviation **D2** pre-registered exactly this gap and
deferred it: "neither the RDR nor `req-list.md` names the
key it takes there… a future RDR (0002 normalization) must
spell the same key for a table row to match on the
recognized outcome as a tag." RDR 0002's Normative Contracts
require the model to declare every tag it matches or writes,
including each tag's provenance (`owned`, `observed`,
`recognized`), but nowhere bind the recognized tag's *name*.
RDR 0001's deviation **D3** separately owns the precedence
order (`owned` > `observed` > `recognized` in `assemble`)
that produces the shadowing behavior, so this RDR inherits
that order rather than deciding it.

### Technical Environment

Go; `internal/resolve` (the RDR 0001 resolution kernel,
implemented) and the RDR 0002 normalizer surface (Final,
unimplemented). The seam is the assembled evaluation view
the kernel builds and the declared tag model the normalizer
will produce.

## Research Findings

### Investigation

Read the kernel source (`internal/resolve/resolve.go`:
`recognizedTagKey`, `assemble`, `Resolve`'s outcome gate),
RDR 0001's implementation record (deviations D2/D3,
`verification.md`, REQ-17 in `req-list.md`), and RDR 0002's
Technical Design and Normative Contracts (tag declarations,
provenance vocabulary, validation categories, the `<clear>`
sentinel). External prior-art pass (bounded, per the Stage 2
budget) over the `StateMachineRes`/`StateMachineLit` corpora
and the `../state-machines` sibling checkouts; queries,
accepted citations, rejected branches, and negative results
are cached at
`docs/rdr/0008-recognized-tag-key-ownership/evidence/research/prior-art.md`.
Peer proposals RDR 0007 (guard predicate totality) and RDR
0009 (escape-row shape conformance) were read as prior art
at the same `area:internal-resolve` seam.

### Key Discoveries

- **Documented** — The kernel already fixes the key:
  `internal/resolve/resolve.go::recognizedTagKey` is the
  unexported constant `"recognized"`, injected by
  `internal/resolve/resolve.go::assemble` with
  `ProvenanceRecognized` only when `in.Recognized` is
  non-empty. RDR 0001 D2 records this as implementation
  latitude ("adds no public surface"), pinned by the frozen
  `recognizedTagSensitiveTable` fixture.
- **Documented** — Outcome gating never reads the key: the
  gate is `row.Outcome != in.Recognized` plus the
  outcome-alphabet check, so the tag-key path is a secondary
  affordance (REQ-17: matching/guarding on the recognized
  outcome *as a tag* in the assembled view).
- **Documented** — The reserved key is a singleton by
  inheritance, not by this RDR's choice: RDR 0001's `Input`
  carries exactly one `Recognized` string per resolve, so at
  most one recognized outcome exists in any view.
  Multi-recognizer or prior-step outcomes are out of scope
  here — they would first have to extend RDR 0001's `Input`
  contract (premortem P-2).
- **Documented** — RDR 0002's Normative Contracts require
  declaring every matched/written tag with provenance
  (`owned`, `observed`, `recognized`) but never bind the
  recognized declaration's *name*. `recognized` appears in
  RDR 0002 exclusively as one of three provenance values,
  never as a tag key, and no tag key named `outcome` appears
  in the RDR 0002 body — so the two sides do not collide
  today on any concrete spelling; the gap is unbound
  authority, not a live incompatibility.
- **Documented** — RDR 0002's validation-failure contract is
  explicitly extensible: "stable data-level categories …
  including at minimum malformed TOML, unknown schema field,
  … ambiguous overlap." An additive category does not reopen
  the Final contract's closed surface.
- **Documented** (prior art) — Surveyed engines fix the name
  of the engine-injected datum; author vocabulary is a
  disjoint namespace. Both external cites are
  sibling-checkout-relative (`../state-machines/`, not under
  this repo). XState: guards receive
  `{ context, event }`
  (`../state-machines/repos/xstate/packages/core/src/guards.ts::GuardArgs`)
  — the injected event occupies the framework-fixed `event`
  slot; authors name event *types*, never the slot. Inngest:
  "Read-Only Fields … system-managed and cannot be modified
  through write operations"
  (`../state-machines/repos/inngest/pkg/api/v2/README.md`).
  In-repo: RDR 0002
  already reserves the `<clear>` write sentinel — a
  carrier-owned token adjacent to author-named vocabulary.
  Negative result: no surveyed engine lets authors rename
  the injected event/outcome slot.
- **Documented** (prior art, resolved at Stage 4) — W3C
  SCXML §5.10 reserves system-variable names against author
  binding: a conformant document "MUST NOT contain ids
  beginning with '_'", and the Processor "MUST cause any
  attempt to change the value of a system variable to fail".
  Stage 2 could not open the spec locally and demoted it to
  A3; Stage 4 fetched and quoted it (A3, research cache).
  The posture corroborates XState/Inngest; the sigil
  *mechanism* is rejected under Load-Bearing Decisions /
  Naming.

### Critical Assumptions

- **A1 RDR 0002's data-level validation-category list is
  additively extensible without reopening its Final
  contract.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0002 Normative Contracts — "Validation
    failures MUST retain stable data-level categories before
    CLI mapping, **including at minimum** malformed TOML, …
    and ambiguous overlap." The fence is a floor, not a
    closed enum, and RDR 0002 demonstrably knows how to write
    closure when it means it (Technical Design: the outcome
    alphabet is "the **closed set**"; RDR 0001 Normative
    Contracts: the refusal-kind set "is **exactly**"). RDR
    0002 Status is `Final` (Metadata), so the fence is
    locked. No closed-set consumer exists downstream: RDR
    0005 never enumerates RDR 0002's data-level categories at
    all (zero occurrences), its nearest table is a
    different-layer "**Minimum** stable code strings" list of
    CLI codes mapped from kernel refusal kinds, and its
    Decision Rationale keeps the envelope append-only
    ("Error envelopes stay append-only through `CLIError`
    fields"). The only genuinely closed downstream set is RDR
    0001's five refusal kinds — a different layer this RDR
    does not touch. Corroborating (A1d): `recognized` appears
    in RDR 0002 only as a provenance value, never as a tag
    key, and both canonical fixtures name the
    recognized-provenance tag `outcome`
    (`0002-…/evidence/spikes/rdr-fixture.toml`,
    `kata-fixture.toml`: `[tags.outcome]` /
    `provenance = "recognized"`) — the name `recognized` is
    unclaimed vacancy, not a collision.
  - **If wrong**: the `reserved tag key` category cannot be
    added by this RDR; enforcement must be re-homed or RDR
    0002 reopened, surfacing as a Stage 8.1 cross-RDR
    conflict.
- **A2 The assembled evaluation view is the only surface
  through which a table can observe the recognized outcome
  *in tag vocabulary* (`Row.Outcome` reads it as a
  first-class non-tag field and is unaffected), and match
  patterns and guard predicates read the identical assembled
  view in one resolve.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: Package-wide sweep of every
    `Input.Recognized` read in `internal/resolve/resolve.go`
    finds five non-test read statements across four sites,
    exactly one of which reaches
    tag vocabulary: `::assemble` (`view.tags[recognizedTagKey]
    = taggedValue{value: in.Recognized, provenance:
    ProvenanceRecognized}`, the `:151`/`:153` pair) is the sole
    injection point. Three
    compare the raw string and never route
    through `TagSet` — `::Resolve` `in.Table.models(
    in.Recognized)` (`:321`), `::Resolve` (`:332`) and
    `::escapeOrRefuse` (`:479`)
    `row.Outcome != in.Recognized` — confirming `Row.Outcome`
    is a first-class non-tag field unaffected by any reserved
    tag-key rule. The fifth, `::refuse` echoing `Recognized`
    onto the `Refusal` (`:518`), is a diagnosis field, not a
    table-visible surface. Same-view half: `::Resolve`
    assembles once (`view := assemble(in)` — the only
    `assemble` call in the package) and threads that one
    value to both channels — matcher `view.matches(row.Match)`
    and guards `gate(candidates, in.Guards, view)` →
    `::evaluateGuard` → `seam.Evaluate(guard, view)` — with
    `::escapeOrRefuse` reusing the same value. `TagSet` wraps
    a map, so pass-by-value shares one backing map; there is
    no re-derivation anywhere, and premortem P-7's
    guard/matcher view split is structurally impossible at
    HEAD.
  - **Test gap (not a falsification)**: the property is true
    by construction but **unpinned**.
    `internal/resolve/resolve_test.go::TestReq17_OwnedObservedAndRecognizedTagsAllReachSelection`
    pins that the recognized tag reaches *match* selection,
    but the suite's only guard seam discards the view —
    `internal/resolve/fixtures_test.go::fixtureGuards.Evaluate`
    has signature `Evaluate(guard string, _ resolve.TagSet)`.
    So no test would catch a regression that split the two
    views. Phase 3's guard/matcher same-view test is
    therefore net-new coverage, not a confirmation of
    existing coverage.
  - **If wrong**: reserving the view key leaves a second
    unbound name elsewhere (or a guard/matcher view split —
    premortem P-7) and the ownership gap survives this RDR.
- **A3 Established statechart standards reserve
  system-injected variable names (W3C SCXML §5.10:
  `_event`, `_sessionid`; authors must not bind names in
  the reserved `_` namespace).**
  - **Status**: Verified
  - **Method**: Prior Art
  - **Evidence**: W3C SCXML Recommendation §5.10 "System
    Variables" [normative], fetched at Resolve from
    `https://www.w3.org/TR/scxml/` (not in any local corpus —
    Stage 2's demotion reason; now retrieved and quoted in
    the research cache). Reserved namespace: "Variable names
    beginning with '_' are reserved for system use. A
    conformant SCXML document MUST NOT contain ids beginning
    with '_' in the `<data>` element." Enforcement: "The
    Processor MUST cause any attempt to change the value of a
    system variable to fail and MUST place the error
    'error.execution' on the internal event queue when such
    an attempt is made." Injected datum: "The SCXML Processor
    MUST bind the _event variable when an event is pulled off
    the internal or external event queue to be processed."
    Scope note: SCXML reserves a whole sigil *namespace* and
    enforces at write time, where this RDR reserves one bare
    word and enforces at load/lint — the transferable finding
    is the ownership posture (carrier fixes the injected
    datum's name; author vocabulary may not rebind it), not
    the sigil mechanism, which is rejected under Load-Bearing
    Decisions / Naming.
  - **If wrong**: prior-art alignment rests on XState and
    Inngest alone; the choice still stands on blast radius
    and the QOC matrix, but the Decision Rationale's
    prior-art row weakens from "standard + engines" to
    "engines".
- **A4 No fixture or authored table at HEAD declares or
  writes an owned/observed tag named `recognized`, so the
  reserved-key rule invalidates nothing that exists.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: Repo-wide whole-token sweep for
    `recognized` returns zero uses as a tag *key* in
    owned/observed position. In Go table data the literal
    key appears exactly once —
    `internal/resolve/fixtures_test.go::recognizedTagSensitiveTable`,
    where it is a read-only **match pattern**
    (`Match: {{Key: "status", …}, {Key: "recognized", Value:
    "successful"}}`) against the tag `assemble` injects; that
    fixture's `RequiresOwned` is `["status"]` and its
    `NextTags`/`Writes` are `{Key: "status"}`, so the key
    never occupies owned, observed, or write position. The
    three committed TOML fixtures use `recognized` only as a
    provenance value under an author-named key
    (`0002-…/evidence/spikes/rdr-fixture.toml` and
    `kata-fixture.toml`: `[tags.outcome]`;
    `0003-…/evidence/spikes/guard-fixture.toml`:
    `[tags.rewind_target]`). Every remaining hit is prose.
    The reserved-key rule invalidates nothing at HEAD, so the
    mechanical-rename fallback below is not triggered.
  - **If wrong**: implementation must rename the colliding
    fixture/table tags before the lint lands; scope grows by
    a mechanical rename, not a design change.
- **A5 The kernel's D3 precedence (`owned` > `observed` >
  `recognized`) can remain unchanged as the deterministic
  backstop, because a table that passes the reserved-key
  validation can never present a colliding owned/observed
  key at resolve time.**
  - **Status**: Verified
  - **Method**: Derivation
  - **Evidence**: Written out against RDR 0002's exact
    declaration semantics. (P1) RDR 0002 Normative Contracts:
    "The model MUST declare every tag it matches or writes,
    including each tag's provenance: owned, observed, or
    recognized" — a total obligation over both sides that
    binds exactly one provenance per declared tag name. (P2)
    An undeclared use is rejected at load: RDR 0002 Technical
    Design rejects "rules that match on missing tag
    declarations" and "writes to non-owned tags", carried as
    the data-level categories `unknown tag` and `write to
    non-owned tag`; RDR 0002 Failure Modes puts the timing
    before acceptance ("a lint failure before the model is
    accepted"). (P3) This RDR's reserved-key rule forbids an
    owned/observed declaration named `recognized`. Therefore:
    by P3 a conforming table declares no owned/observed
    `recognized`; by P1+P2 it cannot match or write a tag it
    did not declare; so no owned/observed tag keyed
    `recognized` reaches `assemble`, and the
    last-writer-wins overwrite of the recognized entry at
    `internal/resolve/resolve.go::assemble` (the `Observed`
    then `Owned` loops) is unreachable for conforming input.
    D3 precedence therefore needs no change.
  - **Scope of the guarantee (carried into the design)**: the
    derivation is conditional on the table having passed RDR
    0002 validation. Hand-constructed and non-TOML-produced
    input bypasses that path entirely, so the branch stays
    reachable there — which is exactly why D3 must remain as
    the deterministic backstop rather than be removed, and
    why the data channel needs the separate producer
    obligation (A6). RDR 0001 `verification.md` records the
    present unreserved state as "deterministic and
    refusal-shaped, not a breach", consistent with keeping D3
    in place.
  - **If wrong**: shadowing is reachable for conforming
    tables and the silent no-match returns; the kernel would
    then need its own reserved-key breach check, changing
    kernel surface.
- **A6 RDR 0009's producer-obligation seam does not reach
  this RDR's data channel, so the reserved-key input
  precondition needs an enforcement locus of its own.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0009 Normative Contracts scope the
    obligation to "a PRODUCER obligation on every constructor
    of resolve.Row values", with the predicate "exactly
    `len(row.Escape) != 0 && len(row.Writes) != 0`" applied
    to every row "at `Resolve` entry" via a kernel-exported
    check. The string `Input` appears **zero** times in RDR
    0009 (case-sensitive, whole file). The exclusion is
    structural, not merely an omission: 0009's own rationale
    is that its breach is a property of the table value
    "independent of the input tuple" — the opposite trigger
    from an `Input` precondition. So 0009 supplies the
    *pattern* but no boundary this RDR can ride, and the
    Enforcement-locus decision stands as this RDR's own.
    The default lean (b) is confirmed available at the
    kernel: `internal/resolve/resolve.go::Resolve` has
    signature `(Result, error)` and its doc comment reserves
    the channel verbatim — "the error return is reserved for
    programmer mistakes, not for modeled refusals" (RDR 0001
    REQ-6; `Resolve`'s body currently has no non-nil error
    path, so the reserved return is unexercised surface an
    `Input` precondition would land on).
  - **Dependency qualification (RDR 0009 is `Draft`, not
    Final)**: 0009 is an original Draft at pre-Resolve (all
    seven of its assumptions `Pending`). The *negative*
    conclusion above is robust regardless — it rests on what
    0009 excludes, and a Draft is likelier to widen than
    narrow. But the *positive* inference is precedent to
    cite, not a contract to depend on: 0009's predicate form
    is explicitly "sharpened at Pre-Lock", so this RDR must
    not bind to `ValidateTable` or any other symbol name from
    0009 (illustrative there, not fixed). No ordering
    dependency is created — RDR 0009's Background records
    "Disjoint answer spaces, no ordering dependency" against
    this RDR's kata (0009's own self-assessment in prose; it
    keeps no `triage.md` artifact).
  - **If wrong** (0009's seam does reach `Input` after all):
    the obligation rides it as originally drafted and the
    Enforcement-locus decision collapses to a cross-reference
    — strictly less surface than the default lean.
  - **Residual (assumption resolved; the design choice it
    gates is not)**: verification confirms no peer seam
    carries this channel, so the Enforcement-locus decision
    is genuinely this RDR's to make and cannot be deferred to
    0009. If it settles to candidate (a) — unchecked
    documented obligation — resolve-time data keyed
    `recognized` still shadows the recognized outcome under
    D3 (premortem P-1/P-6), detected by nothing. That is the
    cost the lean toward a checked locus exists to avoid;
    Testing Strategy scenario 6 pins whichever branch is
    chosen.

- **A7 Reserving the key in `Row.RequiresOwned` is a name
  constraint this RDR may impose without reopening RDR 0007,
  which owns the field's *meaning*; and the unreserved case
  is reachable — a row naming `recognized` there yields
  `owned_state_unavailable` naming the reserved key.**
  - **Status**: Pending
  - **Method**: Source Search + Peer RDR
  - **Verification plan**: the reachability half is already
    read at Pre-Lock —
    `internal/resolve/resolve.go::missingOwned` iterates
    `row.RequiresOwned` and tests
    `view.has(key, ProvenanceOwned)`, where
    `::TagSet.has` compares `tv.provenance == prov`, so a
    key present under `ProvenanceRecognized` fails the
    owned-only test and is reported missing (a second
    consumer exists at `resolve.go:559`). Confirm by test
    (Testing Strategy scenario 9) rather than by reading
    alone. The ownership half needs RDR 0007's concurrence
    that a name reservation is not a semantics change: 0007
    Metadata claims "the single normative home" of
    `RequiresOwned`'s meaning and its normative block
    forbids conflating `missingOwned`'s owned-only test with
    provenance-blind guard presence — neither of which this
    clause touches, since it constrains only which *name*
    may appear. Reconcile at Stage 7.1 cluster-reconcile.
  - **If wrong** (a name constraint does reach 0007's
    contract): the clause moves to a cross-RDR note or an
    advisory lint rule, and the `RequiresOwned` channel
    reverts to a documented-only residual failure mode —
    strictly the F-8 status quo, not a new defect.

## Proposed Solution

### Approach

`recognized` is a **reserved kernel keyword** — the carrier
owns the identity of the recognized-provenance tag key, and
the model conforms to it rather than choosing it. Concretely:

1. **Ratify D2 as contract.** The assembled evaluation view
   binds the freshly recognized outcome under exactly the
   key `recognized`. What D2 recorded as implementation
   latitude becomes the normative answer to the ownership
   question; the shipped constant and frozen fixture are
   already conforming, so no kernel code changes.
2. **Constrain the declaration, keep the declaration.** RDR
   0002's rule that the model declares every tag it matches
   stands: a table that matches or guards on the recognized
   outcome as a tag still declares it — but the declaration
   with provenance `recognized` MUST be named `recognized`.
   The author declares *that* they use the affordance, not
   what it is called.
3. **Enforce at load/lint, not at resolve.** A
   recognized-provenance declaration under any other name,
   or an owned/observed declaration named `recognized`, is a
   data-level validation failure in a new `reserved tag key`
   category — reported by RDR 0002's normalizer before any
   resolution, which converts the seed's silent no-match
   discovery into an authoring-time diagnosis. No new
   refusal kind; no kernel disposition change for conforming
   tables.
4. **Reserve the key on the input data channel too.** Lint
   sees declarations, not resolve-time data: owned/observed
   input tags arrive at the kernel as data no normalizer
   touches, and D3 precedence would resolve a collision
   *against* the declared owner (premortem P-1/P-6). So
   producers of kernel `Input` — the accessor layer for the
   owned snapshot, the caller for observed context — carry a
   producer obligation not to supply a tag keyed
   `recognized`. This RDR adopts RDR 0009's producer-obligation
   *pattern* (obligation on the producer, breach on the Go
   error path, no new refusal kind); it does **not** inherit
   0009's enforcement locus, which is a different channel —
   0009 constrains `resolve.Row` construction, never `Input`.
   Which locus carries this obligation is the one decision
   this RDR leaves open (Load-Bearing Decisions / Enforcement
   locus, A6).
5. **Reserve the key in `Row.RequiresOwned` too.** An owned-key
   *reference* is a third way the name arrives: a row naming
   `recognized` there fails `missingOwned`'s owned-only
   provenance test and refuses `owned_state_unavailable`
   naming the reserved key. Load/lint rejects the name; the
   field's *meaning* stays RDR 0007's (A7).
6. **Shadowing becomes unreachable on all three channels, not
   re-decided.** The declaration channel is closed by
   load/lint (A5); the data channel by the producer
   obligation (A6); the owned-reference channel by the same
   lint rule (A7). RDR 0001 D3's precedence stays untouched
   as the kernel's deterministic backstop behind all three,
   for non-conforming input only.

### Technical Design

Data flow: the author writes a tag declaration with
provenance `recognized` (name constrained to `recognized`);
the RDR 0002 normalizer validates names at load/lint —
rejecting wrong-named recognized declarations and
reserved-name owned/observed declarations — and emits
normalized rows whose `Match`/`Guard` references to the
recognized outcome use the reserved key; the kernel's
`assemble` continues to inject `in.Recognized` under that
key with `ProvenanceRecognized`, unchanged.

The spelling authority is this RDR's Normative Contracts
(both sides quote the literal). Drift is guarded
behaviorally, not literal-vs-literal: the conformance test
MUST run the real kernel and assert the assembled view binds
the recognized outcome under the normalizer's spelling — a
test comparing two hardcoded literals verifies the doc
against itself and is non-conforming (premortem P-5).
Whether the normalizer additionally references an exported
kernel constant is implementation latitude, sharpened at
Resolve/Pre-Lock — the default lean is the behavioral test
over new exported kernel surface, matching RDR 0001 D2's
"no public surface" posture and RDR 0009's minimal-export
precedent.

#### Normative Contracts

Contract count (Resolve, evidence-grounded): **one**
independent load-bearing contract — the identity of the
recognized-provenance tag key. The blocks below express it
across the three channels that key can arrive through:
declaration (blocks 1–3), resolve-time data (block 4), and
row-level owned-key *references* (block 5). They are not
separable seams — reserving only the
declaration channel leaves the silent-shadowing hole this
RDR exists to close (A5's scope caveat: validation cannot
see non-TOML producers), reserving only the data channel
leaves the authoring-time diagnosis unfixed, and leaving
`RequiresOwned` unreserved re-admits the same
confusing-refusal shape one field over (A7). Block 6 is a
negative contract (what does not change), not an independent
one. So the ≥2 split signal is **not** tripped; `foundational`
is earned on the cross-RDR axis (this RDR binds 0001's
carrier to 0002's model), not on contract count.

```normative
When a resolve carries a freshly recognized outcome, the
assembled evaluation view MUST bind it under exactly the tag
key `recognized`
(`internal/resolve/resolve.go::recognizedTagKey`). The key
is a reserved kernel keyword: table authors conform to it
and MUST NOT rebind it. The reservation holds
unconditionally — the key is reserved whether or not a given
resolve binds it — while the binding obligation is scoped to
resolves that carry an outcome: `assemble` injects only for a
non-empty `Input.Recognized`, so an absent outcome yields a
view with no `recognized` key (unreachable past the
outcome-alphabet gate for any alphabet that excludes the
empty string; the conformance test's input domain is
non-empty outcomes).
```

```normative
A tag declaration with provenance `recognized` MUST be named
`recognized`, and a flow MUST carry at most one such
declaration. A tag declaration with provenance `owned` or
`observed` MUST NOT be named `recognized`. Any violation is
a data-level validation failure in the `reserved tag key`
category, reported at table load/lint before any resolution
— never a kernel refusal. The reserved-key comparison is
byte-exact on the **post-parse** key string: case-sensitive,
no trimming, no folding (so `Recognized` is an ordinary,
unreserved name, and `" recognized"` — a
whitespace-bearing key — is likewise ordinary and
unreserved). TOML quoting is a surface artifact, not a name
variant: `[tags."recognized"]` and `[tags.recognized]` parse
to the identical key string and are therefore both reserved.
```

```normative
Every `reserved tag key` failure — and any undeclared-tag
failure whose offending key is the reserved name — MUST
carry, at the data level, the offending name, the required
name `recognized`, and the rule (the kernel owns this key;
declarations conform to it). A category consumer may map it,
but the guidance travels in the failure data, not the
renderer.
```

```normative
Producers of kernel `Input` MUST NOT supply an owned or
observed tag keyed `recognized`; the reserved key enters the
assembled view only through `Input.Recognized`. A breach is
a producer programmer mistake — not table data, and never a
new `RefusalKind` — and travels the Go error path RDR 0001
reserves for programmer mistakes. The enforcement locus is
open (Load-Bearing Decisions / Enforcement locus): it is NOT
inherited from RDR 0009, whose obligation constrains
`resolve.Row` construction and does not reach `Input`.
```

```normative
A normalized row's `Row.RequiresOwned` MUST NOT name
`recognized`. This is a name reservation only: what
`RequiresOwned` *means* is owned by RDR 0007 (its Metadata
declares it "the single normative home" of the field's
meaning — post-guard write-dependency keys — and pins
`resolve.go::missingOwned`'s provenance-specific
`view.has(key, ProvenanceOwned)` test as deliberately
distinct from provenance-blind guard presence); this RDR
cites that rule rather than restating it, per RDR 0007's own
"peers cite rather than restate" convention. Enforcement is
the same load/lint locus and the same `reserved tag key`
category as the declaration channel. Without this, a row
carrying `RequiresOwned: ["recognized"]` finds the key
present under `ProvenanceRecognized`, fails the owned-only
test, and yields `owned_state_unavailable` naming the
reserved key — the confusing-refusal shape this RDR exists
to eliminate, one field over (A7).
```

```normative
Kernel disposition is unchanged by this RDR for conforming
input: no new `RefusalKind`, no change to RDR 0001 D3's
provenance precedence (`owned` > `observed` > `recognized`),
and no behavior change for any conforming table and
conforming input. D3 remains the deterministic backstop
behind the producer obligation for non-conforming input. If
the enforcement locus resolves to a kernel-side precondition,
breaching input gains a non-nil Go error where it previously
resolved — a programmer-mistake path, not a disposition
change for any conforming caller.
```

#### Load-Bearing Decisions

- **Identity** — the recognized-provenance tag key is the
  exact string `recognized`; equality is byte-exact string
  match on the **post-parse** key string — no case folding, no
  trimming, no aliasing. Near-spellings that parse to a
  different key string (`Recognized`, `RECOGNIZED`, a
  whitespace-bearing `" recognized"`) are by definition
  ordinary unreserved names. Quoting is *not* such a variant:
  `[tags."recognized"]` parses to the same key string as
  `[tags.recognized]`, so it is reserved. Whether lint
  additionally warns on case- or whitespace-variants of the
  reserved word is a Resolve
  question, not identity (premortem P-10).
- **Naming** — canonical name `recognized`, matching the
  shipped `recognizedTagKey` constant and RDR 0001's frozen
  fixture. Rejected: a sigil-guarded name (`_recognized`, or
  an angle-bracket token like RDR 0002's `<clear>`) — it
  would repin the shipped constant and fixture for no added
  safety once the reserved-key validation exists, and it
  diverges from the vocabulary authors already see in
  `Row.Outcome` diagnostics.
- **Enforcement locus (open — settle at Resolve)** — where
  the `Input` producer obligation is *checked*. The
  declaration channel is settled (RDR 0002 load/lint); the
  data channel is not. Three candidates: (a) an unchecked
  documented obligation on producers, with D3 as the only
  backstop — zero new surface, no detection; (b) a
  kernel-entry precondition on `Resolve` returning the Go
  error RDR 0001 reserves, mirroring RDR 0009's shape but on
  `Input` rather than `Row` — detection at the real boundary,
  at the cost of kernel-side surface this RDR's blast-radius
  argument claims to avoid; (c) an exported construction-time
  predicate over `Input` that producers call, matching RDR
  0009's exported-predicate form. Default lean: (b), because
  it is the only locus that sees every producer including
  non-TOML callers, and RDR 0001 already reserves the error
  path for exactly this breach class. The cost is explicit —
  it concedes a narrow kernel-side check, so the "no kernel
  change" claim in the Decision Rationale is scoped to
  *disposition for conforming input*, not to zero surface.
  **A6 (Verified at Resolve) settles the seam question and
  leaves the locus choice standing**: RDR 0009 constrains
  `resolve.Row`/table shape and never mentions `Input`, so
  there is no peer boundary to ride and this RDR must pick.
  Resolve strengthens the lean toward (b) without closing
  it: `internal/resolve/resolve.go::Resolve` already returns
  `(Result, error)` with the error "reserved for programmer
  mistakes" and **no non-nil error path in its body today**,
  so (b) adds a check on reserved-but-unexercised surface
  rather than widening the signature. Note (b) and (c) are
  not exclusive — RDR 0009 pairs exactly these two as "one
  predicate, two call sites" so the enforcement points
  cannot drift; the same pairing is available here. Pick the
  final form at Pre-Lock, and do not bind to any symbol name
  from 0009 (still `Draft`, its predicate form explicitly
  "sharpened at Pre-Lock").

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Reserved key injected into the assembled view | Predecessor (RDR 0001 D2, implemented) | Available | Ratified from latitude to contract; no code change |
| Tag declarations with provenance | Predecessor (RDR 0002, Final unimplemented) | Deferred | Name constraint on `recognized`-provenance declarations |
| Load/lint name validation | Predecessor (RDR 0002 normalizer) | Deferred | One additive validation category: `reserved tag key`; upstream of RDR 0006's normalized-graph lint, not a competing acceptance rule |
| Input-boundary producer enforcement | This RDR (pattern borrowed from RDR 0009) | Open | Locus unsettled — see Load-Bearing Decisions / Enforcement locus (A6) |
| `RequiresOwned` name reservation | This RDR; field meaning stays RDR 0007's | Pending (A7) | Same lint locus and category; cite 0007, do not restate its semantics |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Recognized-outcome injection | `internal/resolve/resolve.go::assemble` | none | Reuse | already conforming |
| Outcome gating | `internal/resolve/resolve.go::Row` `Outcome` field | none | Reuse | tag-key path stays a secondary affordance |
| Reserved-token precedent | RDR 0002 `<clear>` write sentinel | sentinel is syntactically un-declarable; `recognized` is declarable | Reuse pattern (reserve + validate) | validation category instead of sigil syntax |

### Decision Rationale

QOC matrix (question: *who owns the name of the
recognized-provenance tag key in the assembled view?*).
Approaches: **R** = reserved kernel keyword (chosen), **M**
= model-carried declaration, **W** = withdraw the tag-key
affordance.

| Criterion | R — reserved keyword | M — model-carried | W — withdraw |
| --- | --- | --- | --- |
| Correctness fit (fixes silent no-match discovery) | collision impossible after lint; diagnosis at authoring time | correct only if declaration threads everywhere; misdeclaration becomes a runtime concern | dissolves the tag path but removes the capability the user wanted |
| Prior-art alignment | matches the W3C SCXML §5.10 reserved-system-variable rule (A3, verified), XState's fixed `event` slot, Inngest system fields, and RDR 0002's `<clear>` precedent — standard **and** engines | no surveyed engine lets authors rename the injected slot (negative result, research cache) | no engine surveyed withdraws event visibility from guards |
| Blast radius | no disposition change for conforming input; one additive lint category + declaration constraint in 0002; one input-precondition check whose locus is open (A6) and whose worst case is a narrow kernel-entry predicate | new public kernel surface (`Input` key field or parameterized `assemble`), normalized-table surface, threads through 0005; frozen fixture repinned | violates locked REQ-17 of implemented Final RDR 0001; removes guard-visible recognized value that RDR 0007's chosen guard domain reads |
| Reversibility | high — a later RDR could still parameterize; reserving now forecloses nothing | low — public surface, once shipped, is load-bearing | low — deleting shipped kernel behavior and reopening 0001 |
| Consistency with 0007/0009 posture at this seam | same shape as 0009: producer/lint obligation, kernel unchanged for conforming input, no new refusal kind | cuts against both peers: grows kernel surface to solve an authored-data problem | cuts against 0007, whose guard domain assumes the assembled view carries the recognized tag |
| Cost | doc ratification + lint checks inside work 0002 already owes | largest: kernel, table schema, CLI threading, fixtures | medium code deletion + cross-RDR contract reopening |

Deciding rows: **blast radius** and **prior-art alignment**
— R is the only branch with no kernel disposition change and
positive prior art; M has a strict negative prior-art result
and the largest surface growth; W reopens a locked,
implemented contract (REQ-17) and breaks RDR 0007's
guard-domain assumption. The consistency row is
corroborating, not decisive: R independently wins before
peer posture is weighed. Rejected alternatives are analyzed
below.

Sibling-path check: the chosen approach adds one identity
rule (a reserved tag key). Searched for an adjacent path
already making that decision.
`internal/resolve/resolve.go::recognizedTagKey` is the only
reserved key in the kernel (whole-package sweep for key
constants and "reserved" yields nothing else), and RDR
0002's `<clear>` sentinel is the nearest sibling
reserved-token decision — its mechanism (syntactically
un-declarable sigil) was considered and rejected under
Load-Bearing Decisions / Naming; its *pattern* (carrier
reserves a token, validation enforces it) is what this RDR
reuses. Repo-wide (the kernel sweep alone was too narrow to
license a repo-wide claim), one further shipped sibling
exists: `internal/cli/respond/respond.go:22-25` reserves the
terminal `"ok"`/`"failed"` type names against future
emitters. It is a documented, wholly *unenforced*
reservation — `::OK` unconditionally sets `s.Type = "ok"` and
nothing validates a future `Stream` emitter — which makes it
a live in-repo instance of Enforcement-locus candidate (a),
and evidence for how (a) ages rather than for its adequacy.
No parallel signal is invented.

Premortem: hardened (variant) — critic verdict PASS, no
switch forced. Every finding is folded into live text; the
P-numbers are cited inline where each cure lives.

Joint-check: clear (7 peers). One inbound reference, not a
dependency: RDR 0007 (`Final`) quotes this RDR's `Input`
producer sentence in its A13 evidence
(`0007…md:952-958`) as the last item in a five-peer sweep
showing *no* constraint on caller-supplied observed tags. It
corroborates a negative-existential rather than supplying a
premise — A13's "accepted exposure" verdict is unchanged or
reinforced however this RDR's enforcement locus settles, and
0007 elsewhere lists this RDR under "Related but distinct
seams, deliberately not folded in" (`:158-163`) and
disclaims that it owns anything 0007 depends on (`:811-815`).
What the reference does create: A13's claim that the RDR set
holds exactly one input-side producer obligation is
falsifiable by any *future* peer adding a second one — a
cluster-reconcile watch item, not a constraint on this RDR.

## Alternatives Considered

### Alternative 1: Model-carried key declaration (model owns identity)

**Description**: The tag declaration with provenance
`recognized` carries an author-chosen name; the normalized
table surfaces that name; the kernel's `assemble` binds
`in.Recognized` under the table-supplied key (a new `Input`
field or parameter).

**Pros**:

- Preserves full author naming authority — the most
  consistent reading of RDR 0002's "authors name their
  tags" ethos.
- No reserved word in the author's namespace; no new
  validation category.

**Cons**:

- New public kernel surface on an implemented, Final RDR:
  `Input` grows a key field (or `assemble` a parameter),
  RDR 0001 D2's "no public surface" posture is reversed,
  and the frozen `recognizedTagSensitiveTable` fixture is
  repinned.
- The name must thread through the normalized-table surface
  and RDR 0005's CLI path — the widest blast radius of the
  three branches.
- Identity becomes a runtime property of each table: two
  tables in one review can name the same datum differently,
  and a misdeclared name is discovered at resolve time, not
  load time — the original silent-no-match discovery
  problem survives in a new form.
- Strict negative prior-art result: no surveyed engine lets
  authors rename the injected event/outcome slot (research
  cache).

**Reason for rejection**: largest blast radius of the three
branches, spent reversing a shipped posture to solve
statically what a lint rule solves at load time, with zero
supporting prior art.

### Alternative 2: Withdraw the tag-key affordance

**Description**: Remove the recognized tag from the
assembled view; recognized outcomes are matchable only via
`Row.Outcome`; the identity question dissolves.

**Pros**:

- Deletes the ambiguity instead of ruling on it; no
  reserved word, no new validation category.

**Cons**:

- Violates REQ-17 of RDR 0001 — a locked, implemented,
  Final contract ("Merge owned, observed, and freshly
  recognized tags into the evaluation view") — so it
  reopens 0001 rather than binding 0001 to 0002.
- Removes guard visibility of the recognized outcome: RDR
  0007's chosen guard domain evaluates predicates over the
  assembled view, so guards could no longer reference the
  recognized value at all — a capability regression beyond
  this RDR's scope, and a bridge/retire conflict with a
  peer's shipped plan.
- Deletes the exact capability the Problem Statement's user
  wanted (matching the recognized outcome in tag
  vocabulary).

**Reason for rejection**: resolves a naming question by
deleting a locked peer contract and a peer-consumed
capability — the highest-cost, least-reversible branch.

### Briefly Rejected

- **Sigil-guarded rename** (`_recognized` / `<recognized>`):
  see Load-Bearing Decisions / Naming.
- **Dual-name aliasing** (kernel writes both the reserved
  key and the author's declared alias): two names for one
  datum in one view — invites the divergence this RDR
  exists to close.

## Trade-offs

### Consequences

- Positive: the naming collision surfaces at authoring time
  as a typed load/lint failure, replacing the silent
  no-match discovery in the Problem Statement.
- Positive: the shipped constant, fixture, and D3 precedence
  all remain valid — the cheapest branch to implement and
  the easiest to reverse.
- Positive: a consistent trust story at the 0001↔0002 seam
  with RDR 0009 (normalizer enforces, no new refusal kind,
  kernel disposition unchanged for conforming input).
- Negative: authors lose naming freedom for exactly one
  tag; `recognized` becomes a reserved word they must learn
  (mitigated by the lint message naming the rule).
- Negative: RDR 0002's implementation grows one validation
  category and two name constraints it did not
  originally spec — the declaration constraint and the
  `RequiresOwned` reservation (A7) — additive, but still
  scope on an unimplemented Final RDR.

### Risks and Mitigations

- **Risk**: RDR 0002's implementer treats the validation
  category list as closed and rejects the addition.
  **Mitigation**: A1 verifies extensibility at Resolve; the
  fence's "including at minimum" wording marks a floor, and
  the category lands through this RDR's normative contract,
  which 0002's implementation prompt will extract.
- **Risk**: resolve-time owned/observed data keyed
  `recognized` bypasses lint entirely (data channel, no
  declaration) and D3 resolves the collision against the
  declared owner — a silent wrong-value match, worse than
  the original no-match (premortem P-1/P-6).
  **Mitigation**: the producer-obligation fence closes the
  data channel; its enforcement locus is the one open
  decision (Load-Bearing Decisions / Enforcement locus,
  A6) — an unchecked obligation leaves D3 as the sole
  backstop, which is why the default lean is a checked
  locus. D3 stays the deterministic backstop behind it
  either way; the breach test is named in Phase 3.
- **Risk**: the `RequiresOwned` name reservation (A7) is
  judged to encroach on RDR 0007, the declared "single
  normative home" of that field's meaning.
  **Mitigation**: the clause reserves a *name* and cites
  0007 for semantics rather than restating them, per 0007's
  own "peers cite rather than restate" convention; it
  touches neither the post-guard write-dependency definition
  nor the provenance-blind-vs-owned-only test distinction
  0007 pins. Reconciled at Stage 7.1; if it does encroach,
  A7's "If wrong" downgrades it to an advisory rule.
- **Risk**: kernel and normalizer spellings drift (two
  constants, one contract).
  **Mitigation**: the conformance test is normatively
  behavioral — real kernel, assembled-view assertion — so a
  doc-vs-doc literal comparison cannot pass for it
  (premortem P-5); the literal is normative in this RDR, so
  drift is a spec breach with a named test.
- **Risk**: category-enumerating consumers (RDR 0005's
  exit-code map, refusal renderers, remediation docs)
  render the new category as generic and strip the guidance
  it exists to carry (premortem P-4/P-11).
  **Mitigation**: the refusal-text content is normative at
  the data level — offending name, required name, rule —
  so a generic renderer still surfaces the resolution; a
  golden-text test is named in Phase 3.
- **Risk**: near-spellings (`Recognized`, quoted/whitespace
  TOML key forms) read as the reserved word to a human but
  are ordinary names to the byte-exact rule (premortem
  P-10).
  **Mitigation**: the comparison rule is normative
  (byte-exact on the parsed key value); an advisory lint
  warning on case-variants is a Resolve question; the
  normalization fixture set is named in Phase 3.
- **Risk**: pre-existing tables using `recognized` as an
  innocent owned/observed tag break at load with no
  migration story (premortem P-9).
  **Mitigation**: none needed as migration — the project
  carries no back-compat obligation and the reservation
  lands before any normalizer or authored table exists (A4);
  the innocent-collision case is a lint fixture, and the
  normative refusal text tells that author what to rename
  and why.

### Failure Modes

Visible: a wrong-named or duplicated recognized-provenance
declaration, a reserved-name owned/observed declaration, or a
row naming `recognized` in `RequiresOwned` (A7),
fails table load/lint with the `reserved tag key` category —
the failure data carries the rule, the offending name, and
the required name, before any resolution runs. A producer
supplying an owned/observed input tag keyed `recognized`
breaches the producer obligation and surfaces on the Go
error path as a programmer mistake, at whichever locus A6
settles — never as a modeled refusal.

Silent, residual: an author who declares an *observed* tag
under an innocent name (`result`, `outcome`) intending
recognizer semantics passes every name check and the row
never fires — the original silent no-match survives this
narrow path, because intent is invisible to a
name-and-provenance validator (premortem P-3). Diagnosis:
the table dump shows no recognized-provenance declaration
for the flow while rows gate on outcomes; an advisory lint
heuristic for that shape is a Resolve question. Silent,
backstop-only: input that bypasses both lint and the
producer obligation shadows the recognized value under D3
and surfaces as a deterministic `no_match` refusal.
Recovery in every case is renaming or re-provenancing the
colliding declaration. A developer diagnosing kernel
behavior finds the reserved key documented on
`internal/resolve/resolve.go::recognizedTagKey` with a
pointer to this RDR.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified
- [ ] RDR 0002 implementation underway (the enforcement
      point is its normalizer's load/lint path; this RDR's
      checks land inside that work, not before it)

### Minimum Viable Validation

One end-to-end pair over the normalizer + kernel: (a) a
table declaring a recognized-provenance tag named
`recognized` loads, lints clean, and a row matching on that
tag fires against a recognized outcome — asserted through
the real kernel's assembled view (the behavioral
conformance form, premortem P-5); (b) the same table with
the declaration renamed (and separately, with an owned tag
named `recognized`) fails load/lint with the
`reserved tag key` category whose failure data names the
required name — not a silent no-match at resolve time.

### Phase 1: Contract ratification

Promote D2's key to normative: record the reserved-keyword
rule where implementers read it — this RDR's Normative
Contracts as authority, plus a pointer comment on
`internal/resolve/resolve.go::recognizedTagKey` citing RDR
0008. No kernel behavior change.

### Phase 2: Normalizer name validation

Add the reserved-key checks to RDR 0002's load/lint path:
wrong-named recognized-provenance declarations,
reserved-name owned/observed declarations, and a
reserved-name `RequiresOwned` reference (A7) all report the
`reserved tag key` data-level category.

Not RDR 0006's graph lint, and not a competing acceptance
rule under its "MUST NOT define different acceptance rules"
clause: 0006 consumes "a normalized graph value, not Cobra
command state and **not sparse TOML**" (`0006…md:244-246`),
takes declared tags with provenance as already-valid input
(`:251`), and leaves "table source and normalization … in RDR
0002" (`:458-460`). Declaration-name validation is upstream
of 0006's input contract, in the same pre-acceptance layer
where 0002 already homes `unknown tag`, `unknown context`,
and `unknown accessor`.

### Phase 3: Conformance pinning

The MVV fixture pair plus the premortem-derived test set:
the behavioral spelling test (real kernel, assembled-view
assertion — never literal-vs-literal), the
producer-obligation breach test (observed input keyed
`recognized` caught at whichever locus A6 settles), the
guard/matcher same-view test (A2), the duplicate
recognized-declaration fixture, the reserved-key
normalization fixtures (case and whitespace variants stay
unreserved; the quoted form is reserved), the
`RequiresOwned` reservation pair (A7), and the golden
refusal-text check
(offending name, required name, rule).

## Validation

### Testing Strategy

The matrix the verified assumptions imply. Scenarios 1–2 are
the MVV pair; 3–9 are the premortem- and assumption-derived
set named in Phase 3. Scenarios 2–5 and 7 land inside RDR
0002's normalizer work (no normalizer exists at HEAD — reuse
audit); 1, 3, and 6 are kernel-side and runnable against
`internal/resolve` as it stands; 9 has a half on each side.

Done means: every scenario below has a green test, and the
guard/matcher same-view property (scenario 3) is pinned by a
test that would fail if the two views diverged — it is
currently true by construction but unpinned (A2).

1. **Scenario** (MVV a; A2): a table declaring a
   recognized-provenance tag named `recognized` loads, lints
   clean, and a row matching on that tag fires against a
   recognized outcome — asserted through the **real kernel's**
   assembled view, never by comparing two hardcoded literals
   (premortem P-5).
   **Expected**: the row is selected; the recognized outcome
   is readable at key `recognized` with
   `ProvenanceRecognized`. Backed by
   `internal/resolve/resolve.go::assemble` (sole injection
   point, A2) and the existing shape of
   `internal/resolve/fixtures_test.go::recognizedTagSensitiveTable`.

2. **Scenario** (MVV b; A1): the same table with the
   recognized-provenance declaration renamed, and separately
   with an owned tag named `recognized`.
   **Expected**: both fail load/lint in the `reserved tag
   key` data-level category before any resolution — not a
   silent no-match at resolve time. The category is additive
   under RDR 0002's "including at minimum" fence (A1).

3. **Scenario** (A2 — the net-new coverage): a guard
   predicate and a match pattern both read the recognized
   tag in one resolve, through a guard seam that actually
   inspects the `TagSet` it is handed.
   **Expected**: both observe the identical assembled view.
   This is net-new: the suite's only guard seam today
   discards the view
   (`internal/resolve/fixtures_test.go::fixtureGuards.Evaluate`
   takes `_ resolve.TagSet`), so no existing test would catch
   a guard/matcher view split (premortem P-7).

4. **Scenario** (identity/exactness): near-spellings of the
   reserved word as declared tag names — `Recognized`,
   `RECOGNIZED`, a quoted TOML key `"recognized"`, and a
   whitespace-bearing `" recognized"`.
   **Expected**: the byte-exact rule on the post-parse key
   string treats `Recognized`, `RECOGNIZED`, and
   `" recognized"` as ordinary unreserved names (no folding,
   no trimming); the quoted `"recognized"`, which parses to
   the identical key string as the bare form, **is**
   reserved. Pins the
   Load-Bearing Decisions / Identity rule (premortem P-10).

5. **Scenario** (declaration cardinality): a flow carrying
   two recognized-provenance declarations.
   **Expected**: `reserved tag key` failure — the contract
   admits at most one such declaration, matching RDR 0001's
   single `Input.Recognized` per resolve (Key Discoveries:
   singleton by inheritance).

6. **Scenario** (A6 — producer obligation): kernel `Input`
   carrying an owned or observed tag keyed `recognized`,
   constructed directly (bypassing lint, as a non-TOML
   producer would).
   **Expected**: the breach is caught at whichever locus A6's
   Enforcement-locus decision settles, on the Go error path
   `internal/resolve/resolve.go::Resolve` reserves for
   programmer mistakes — never a new `RefusalKind`. If the
   locus resolves to (a) documented-only, this scenario
   instead pins the D3 backstop: the collision resolves
   deterministically to the owned/observed value and surfaces
   as a typed refusal, never a wrong edge.

7. **Scenario** (golden refusal text; premortem P-4/P-11): a
   `reserved tag key` failure rendered by a
   category-enumerating consumer that does not know the new
   category.
   **Expected**: the failure **data** still carries the
   offending name, the required name `recognized`, and the
   rule — so a generic renderer cannot strip the guidance.
   Golden-text assertion at the data level, not the renderer.

8. **Scenario** (A5 — backstop reachability): a conforming
   normalized table proves the `assemble` collision branch
   unreachable; a hand-constructed non-conforming table
   proves D3 still resolves it deterministically.
   **Expected**: no conforming input reaches the overwrite;
   non-conforming input keeps RDR 0001 D3's precedence
   (`owned` > `observed` > `recognized`) unchanged.

9. **Scenario** (A7 — the `RequiresOwned` channel): a row
   carrying `RequiresOwned: []string{"recognized"}` resolved
   against a view whose `recognized` key is present under
   `ProvenanceRecognized`; and the same name rejected at
   load/lint.
   **Expected**: kernel-side, the key is reported missing and
   the resolve refuses `owned_state_unavailable` naming
   `recognized` — pinning that the unreserved case is
   reachable, not hypothetical (`::missingOwned` tests
   `view.has(key, ProvenanceOwned)`; `::TagSet.has` is
   provenance-specific). Lint-side, the same name in
   `RequiresOwned` fails in the `reserved tag key` category
   before resolution. Runnable against `internal/resolve` as
   it stands for the kernel half.

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

- RDR 0001 `0001-resolution-kernel.md` — REQ-17; artifacts
  `deviations.md` D2 (tag key), D3 (precedence);
  `verification.md` ("deterministic and refusal-shaped, not
  a breach").
- RDR 0002 `0002-transition-table-as-reviewable-data.md` —
  Technical Design (tag declarations, provenance
  vocabulary), Normative Contracts (declare-every-tag,
  validation categories, `<clear>` sentinel).
- RDR 0007 `0007-guard-predicate-totality.md`, RDR 0009
  `0009-escape-row-shape-conformance-ownership.md` — peer
  posture at the `area:internal-resolve` seam.
- `internal/resolve/resolve.go` — `recognizedTagKey`,
  `assemble`, `Row`, `GuardEvaluator`.
- Prior-art research cache:
  `docs/rdr/0008-recognized-tag-key-ownership/evidence/research/prior-art.md`
  (XState `GuardArgs`, Inngest read-only fields; SCXML
  demoted to A3). External source paths there and above are
  relative to the `../state-machines` sibling checkout, not
  this repo.
