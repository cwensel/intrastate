# Recommendation 0008: Ownership of the recognized-outcome tag key name

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). **Reference-only** (guidance; never copy into the
instance body). -->

## Metadata

- **Date**: 2026-08-09
- **Status**: Final
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

**The outcome this RDR commits to is "my row fires, or I am
told why" — not "I get to name this tag."** Those diverge,
and the divergence is a deliberate product stance rather
than an incidental constraint. The author's underlying want
is a working match with a diagnosable failure; the naming
freedom was a *means* they assumed, not the end. This RDR
satisfies the end and withdraws that particular means for
exactly one tag, on the same footing as XState's fixed
`event` slot: the engine-injected datum has an
engine-owned name, and that is a rule authors learn once
rather than a limitation apologized for. The `reserved tag
key` failure text is where that stance is taught, which is
why its content is normative (Normative Contracts, block 3).

The XState analogy is weaker than it looks, and the gap is
conceded rather than papered over: `event` is a field in a
struct the engine owns (`GuardArgs`), disjoint from the
author's `context` namespace, so nothing an author names can
collide with it. Here `recognized` is a bare word carved out
of the single flat `[tags.*]` namespace the author otherwise
owns entirely — an author *can* collide with it, which is
exactly why this RDR needs a validation category where XState
needs none. The prior art supports the *posture* (the carrier
fixes the injected datum's name), not the *mechanism*. The
in-repo evidence also cuts against treating the name as
incidental: all three committed recognized-provenance
declarations chose an author name (`outcome`, `outcome`,
`rewind_target` — A4), so every author observed so far
exercised the freedom this RDR withdraws. The stance still
holds, because a firing row with a diagnosable failure is what
the author is ultimately after and a fixed name is the only
branch that delivers it (Decision Rationale) — but it is a
real cost against an observed convention, not a costless
clarification.

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
  - **Scope limit — this covers the category list only**: the
    "including at minimum" fence lives in RDR 0002's
    *validation-category* normative block (`0002…md:342-348`).
    This RDR's name constraint does not land there; it
    constrains the **`[tags.<tag>]` declaration grammar**,
    which is a *separate* normative block ("The source schema
    MUST use the Resolve spike field layout: root `outcomes`,
    `[model]`, `[tags.<tag>]`, …", `0002…md:266-274`) carrying
    no extensibility fence of its own. So A1 licenses the new
    *category* but not the new *constraint on tag names*; that
    half is licensed by this RDR's Overrides field, which
    states the extension explicitly and is what a Stage 7.1
    cross-RDR read must adjudicate. The distinction matters
    because an 0002 implementer reading 0002 alone sees a field
    layout that admits any `[tags.<name>]` and no pointer to
    this constraint — which is why Phase 2 lands the check
    inside 0002's own normalizer work rather than as a
    detached rule. A12 has since confirmed that the missing
    pointer is a real gap and not merely a stylistic one: no
    engine mechanism carries this RDR's `Overrides` to a
    `Final` peer's implementation, so the handoff is now an
    explicit Prerequisite rather than an assumed traversal.
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
- **A4 The reserved-key rule invalidates nothing at HEAD in
  its *second* clause (no owned/observed tag is named
  `recognized`) but its *first* clause invalidates every
  committed recognized-provenance declaration in the repo —
  three fixtures, each of which RDR 0002 designates a
  canonical example.**
  - **Status**: Verified (both clauses swept separately)
  - **Method**: Source Search
  - **Evidence**: The rule has two clauses with different
    blast radii, so each is swept separately — a sweep of the
    second alone licenses no conclusion about the first.
    *Second clause* (no owned/observed named `recognized`) —
    repo-wide whole-token sweep returns zero uses as a tag
    *key* in owned/observed position. In Go table data the
    literal key appears exactly once —
    `internal/resolve/fixtures_test.go::recognizedTagSensitiveTable`,
    where it is a read-only **match pattern**
    (`Match: {{Key: "status", …}, {Key: "recognized", Value:
    "successful"}}`) against the tag `assemble` injects; that
    fixture's `RequiresOwned` is `["status"]` and its
    `NextTags`/`Writes` are `{Key: "status"}`, so the key
    never occupies owned, observed, or write position.
    *First clause* (a recognized-provenance declaration MUST
    be named `recognized`) — **three violations at HEAD**, the
    complete set of committed recognized-provenance
    declarations: `0002-…/evidence/spikes/rdr-fixture.toml:25`
    and `kata-fixture.toml:22` (both `[tags.outcome]`), and
    `0003-…/evidence/spikes/guard-fixture.toml:36`
    (`[tags.rewind_target]`). RDR 0002 Load-Bearing Decisions
    names these "the canonical examples implementation tests
    must promote" (`0002…md:361`), so the rule's first
    enforcement lands on the peer's own worked examples.
  - **Migration inventory**: the three declarations above plus
    **six** reference sites across all three fixtures —
    `[rule.match.outcome]` at `rdr-fixture.toml:59,74,86` and
    `kata-fixture.toml:47,57`, and
    `[rule.guard.all.rewind_target]` at `guard-fixture.toml:72`.
    The guard site belongs in the inventory for the same reason
    the match sites do: RDR 0002's grammar puts
    `[rule.match.<tag>]`, `[rule.guard.all.<tag>]`, and
    `[rule.guard.unless.<tag>]` in one tag-predicate family
    (`0002…md:266-274`), so a renamed declaration carries every
    one of them. Renaming a declaration without its guard
    reference would leave an undeclared-tag reference that
    trips RDR 0002's pre-existing `unknown tag` category — a
    regression introduced *by* the migration, in a category
    this RDR does not own.
  - **Sub-question this surfaced, since closed (A9,
    Verified)**: in the canonical fixture the
    recognized-provenance tag is *also* the outcome-gating
    predicate — `[tags.outcome]` is matched by
    `[rule.match.outcome] eq = "round-clean"` whose values are
    exactly the root `outcomes` alphabet
    (`rdr-fixture.toml:1,25,59`). Whether that predicate
    normalizes into `Row.Match` or into `Row.Outcome` is not
    stated in RDR 0002, but it *is* decidable against 0002's
    own normalizer spike, which A9 ran: the predicate lands in
    the candidate row's predicate set and the normalized row
    has no outcome field at all. The rename is therefore
    mechanical, and this migration inventory stands as
    written.
  - **If wrong** (the spike's normalization is not what 0002
    implements): the fixture rename becomes a cross-RDR edit
    routed to RDR 0002 rather than fixture work landing under
    this RDR's Phase 2, and block 2's reference-position
    fall-through needs the corresponding rewrite.
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
    reachable there — which is why the data channel needs the
    separate producer obligation (A6).
  - **Re-derived after the Enforcement locus settled**: with
    the predicate applied by `Resolve` at entry, the set of
    input reaching `assemble` unchecked is *narrower* than
    this scope note originally assumed — a non-TOML producer
    no longer bypasses detection, it trips the precondition.
    D3's role as backstop is therefore reduced but not empty,
    and it is retained for two live paths: package-internal
    callers of `assemble` that do not go through `Resolve`
    (the path scenario 8 exercises), and any future entry
    point added without the precondition. Removing D3 would
    also reopen RDR 0001's locked precedence for no gain,
    which is the decisive reason it stays. What is *not* still
    true is the pre-settlement framing that D3 is the only
    thing standing between a non-TOML producer and a silent
    shadow. RDR 0001 `verification.md` records the
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
    (This is a claim about *obligations*, not about which
    values each predicate reads: block 5's widening means both
    predicates traverse `in.Table.Rows`, over disjoint fields.
    That co-residency is what A11 reconciles; it does not make
    0009's obligation reach this RDR's subject.)
    Candidate (b) is confirmed available at the
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
    — strictly less surface than the chosen pairing.
  - **Residual closed at Pre-Lock**: verification confirmed no
    peer seam carries this channel, so the Enforcement-locus
    decision was genuinely this RDR's to make. It is now made
    — (b)+(c) paired, one predicate at two call sites (see
    Load-Bearing Decisions / Enforcement locus). Candidate (a)
    was rejected on the user-outcome criterion: it would leave
    resolve-time data keyed `recognized` shadowing the
    recognized outcome under D3 (premortem P-1/P-6), detected
    by nothing. The surface that choice concedes is carried by
    A8.

- **A7 Reserving the key in `Row.RequiresOwned` is a name
  constraint this RDR may impose without reopening RDR 0007,
  which owns the field's *meaning*; and the unreserved case
  is reachable — a row naming `recognized` there yields
  `owned_state_unavailable` naming the reserved key.**
  - **Status**: Verified
  - **Method**: Source Search + Peer RDR
  - **Evidence**: both halves closed at Reconcile.
    *Reachability* — `internal/resolve/resolve.go::missingOwned`
    iterates `row.RequiresOwned` and tests
    `view.has(key, ProvenanceOwned)`, where
    `::TagSet.has` compares `tv.provenance == prov`; the
    reserved key enters the view under `ProvenanceRecognized`
    (`::assemble`), so it fails the owned-only test, is
    collected into `missing`, and `::gate` refuses
    `KindOwnedStateUnavailable` naming it. Reachable, not
    hypothetical; still pinned by test (scenario 9) rather
    than by reading alone.
    *Ownership* — RDR 0007 (`Final`) scopes its claim to
    meaning, never to names: its "single normative home"
    clause covers the *guard-decidability domain rule*, and
    its `RequiresOwned` normative block is wholly semantic
    (what the field means, what it is not responsible for,
    what diagnosis a listed key yields). Decisively, 0007
    *itself* cites this RDR's rule approvingly as a peer's
    name constraint — "The single input-side producer
    obligation in the whole RDR set constrains one reserved
    key *name*, not provenance origin: RDR 0008's 'Producers
    of kernel `Input` MUST NOT supply an owned or observed tag
    keyed `recognized`'" (`0007…md:951-955`). 0007 draws the
    name/meaning line itself and places this clause on the
    name side, so no concurrence is outstanding and no
    reopening is implied.
  - **Residual for Stage 7.1**: none for ownership. The
    cluster read carries only the standing watch item (0007
    A13's negative-existential, falsifiable by a *future*
    peer).
  - **Enforcement locus corrected at Pre-Lock (3amigo
    IMP-1)**: `RequiresOwned` has **no authored source
    form** — RDR 0002's normative field layout
    (`0002…md`, Technical Design: root `outcomes`,
    `[model]`, `[tags.<tag>]`, `[accessors.<id>]`,
    `[context.<id>]`, `[[rule]]`, `[dump]`) carries no
    `requires_owned` key, and the token has zero occurrences
    repo-wide. The field is a derived normalizer output, so
    the reservation is a producer obligation on the
    normalizer, not a source-lint rule in the `reserved tag
    key` category. Under RDR 0007's narrowing it derives from
    `[rule.write]`, where RDR 0002's pre-existing `write to
    non-owned tag` category already rejects a `recognized`
    entry — so the obligation is discharged by construction
    on the current derivation path and binds any future one.
  - **If wrong** (a name constraint does reach 0007's
    contract): the clause moves to a cross-RDR note, and the
    `RequiresOwned` channel reverts to a documented-only
    residual failure mode — strictly the F-8 status quo, not
    a new defect.

- **A8 Adding an exported `Input` predicate plus a
  `Resolve`-entry call introduces the kernel package's first
  non-nil error path without disturbing RDR 0001's locked
  contracts or the existing suite's assertions.**
  - **Status**: Verified
  - **Method**: Source Search + Peer RDR
  - **Evidence**: Two halves, both read at Pre-Lock.
    *Permission*: RDR 0001 REQ-6 states verbatim — "Modeled
    refusal is a value-level resolver disposition, not a CLI
    error and not the Go error path for parser bugs, IO
    failures, or programmer mistakes"
    (`0001-…/artifacts/req-list.md:34`, `(NC)`). The clause
    affirmatively reserves the Go error path *for* programmer
    mistakes, so a producer breach travelling that path
    conforms to REQ-6 rather than contradicting it; the
    prohibition runs the other way (a modeled refusal must
    not use it). `internal/resolve/resolve.go::Resolve` has
    signature `(Result, error)`, so no signature widens.
    *Blast radius, in-repo*: a sweep of
    `internal/resolve/*_test.go` for the literal
    `Key: "recognized"` returns exactly one hit —
    `fixtures_test.go:237`, inside
    `recognizedTagSensitiveTable`'s `Row.Match` — which is a
    table match pattern, not `Input.Owned`/`Input.Observed`.
    No existing fixture constructs an `Input` that would trip
    the precondition. The kernel has **zero non-test callers**
    at HEAD (`Resolve` is referenced only by its definition at
    `resolve.go:318` and by the suites), and every test call
    site routes through
    `internal/resolve/resolve_test.go::mustResolve`, which
    `t.Fatalf`s on a non-nil error — so a tripped precondition
    surfaces as a red test, never as a silently-ignored error.
    `resolve_test.go:748` additionally pins
    `Resolve(resolve.Input{})` returning a nil error, which the
    predicate must keep true (an empty `Input` carries no
    reserved key).
  - **Scope limit of this sweep (does NOT clear the data
    channel)**: the sweep finds Go *literals*. `Input.Owned` is
    "an accessor-produced owned tag snapshot"
    (`internal/resolve/resolve.go` package doc) — RDR 0004's
    accessor layer derives owned tag keys from artifact
    content at runtime, so a key spelled `recognized` can
    arrive as data with no literal anywhere. A literal sweep is
    structurally incapable of refuting that case; it is
    governed by A10, not by this assumption.
  - **If wrong** (0001 forbids the error path, or existing
    fixtures breach): the locus falls back to (c) alone — the
    exported predicate without the `Resolve` entry call —
    which loses detection for producers that never call it,
    and the Enforcement-locus decision must record that
    reduced guarantee.

- **A9 Every tag-predicate reference to a
  recognized-provenance declaration normalizes into a
  view-read row field — `Row.Match` for
  `[rule.match.<tag>]`, the guard field for
  `[rule.guard.all.<tag>]` / `[rule.guard.unless.<tag>]` —
  rather than into `Row.Outcome`, so renaming the declaration
  carries all of them with it and changes no gating
  semantics.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: settled by re-running RDR 0002's own
    normalizer spike — the peer's executable artifact, not a
    reading of its prose — over all three committed fixtures.
    Command and verbatim output:
    `docs/rdr/0008-recognized-tag-key-ownership/evidence/spikes/a9-normalization.md`.
    Three findings, each decisive for one position. (1) The
    normalized row has **no outcome field to receive a
    predicate**: `0002-…/evidence/spikes/main.go::Row` is
    `{Model, Rule, Kind, Source string; Match, Write []string}`.
    (2) `[rule.match.<tag>]` lands in the predicate set —
    `outcome.eq=round-clean` appears inside `match=[…]`,
    unprefixed and sorted among `stage.eq`, `status.eq`,
    `profile.in`, with no special-casing on the name. (3) Both
    guard positions land in that same set, prefixed:
    `::normalize` merges `rule.Match` (`""`),
    `rule.Guard.All` (`"all:"`), and `rule.Guard.Unless`
    (`"unless:"`) into one map, rendered as `all:iter.lt=3` /
    `unless:profile.eq=small`. The root `outcomes` alphabet is
    model metadata printed on the MODEL line, never routed into
    a row. All three tag-predicate positions are therefore
    view-read; none produces `Row.Outcome`.
  - **Why `Row.Outcome` is not the destination**: it is not
    RDR 0002's field to populate — 0002 never names it. It is
    a kernel field fed from `Input.Recognized`, gated by
    `row.Outcome != in.Recognized`
    (`internal/resolve/resolve.go::Resolve`), while
    `::assemble` separately injects the same value as a tag so
    a row may *also* match it (RDR 0001 REQ-17; D2: "a future
    RDR (0002 normalization) must spell the same key for a
    table row to match on the recognized outcome **as a
    tag**"). The two paths are complementary, which is why
    both readings looked live: the fixture's
    `[rule.match.outcome]` is a tag-level restatement of a
    value the outcome gate reads independently.
  - **Consequence**: the first reading holds, so A4's
    migration inventory, block 2's reference-position
    fall-through, and Phase 2's rename all stand as written.
    The rename is mechanical and changes no gating semantics;
    it does not route to RDR 0002. Noted rather than left
    implicit: the rename makes those fixtures *newly
    functional* against the kernel — the predicate key
    `outcome` cannot match `assemble`'s injected key
    `recognized` today and can afterward. That is a strict
    improvement with nothing depending on the current
    behavior, since no normalizer exists under `internal/` to
    turn these fixtures into a `resolve.Table`.

- **A10 No accessor-produced owned tag key at HEAD, and no
  key derivable from a committed artifact shape, spells
  `recognized` — so the `Input`-precondition breach is not
  reachable from real data today.**
  - **Status**: Verified
  - **Method**: Source Search + Peer RDR
  - **Evidence**: the accessor never names a key of its own —
    the key is the `[tags.<tag>]` table name, so artifact
    content supplies only the *value*. The reference direction
    is one-way (tag → accessor): RDR 0002 Technical Design
    declares "tag name, provenance …, value kind, and
    **optional accessor reference** for observed or owned
    read-back" (`0002…md:215-216`), and the committed fixture
    shows the shape — `[tags.status]` carries
    `accessor = "rdr-status"` while `[accessors.rdr-status]`
    carries only `mode` and `path`
    (`0002-…/evidence/spikes/rdr-fixture.toml`). The peer's
    loader types confirm there is no key channel to carry a
    spelling: `0002-…/evidence/spikes/main.go::Accessor` is
    `{Mode, Path}` with **no key field**, while tags are the
    keys of `Tags map[string]Tag`. RDR 0004 (`Final`) pins the
    same binding from its side — an accessor definition
    declares "a stable name, capability, artifact role,
    **expected tag keys**, timeout policy …"
    (`0004…md:208-212`), its identity is
    `(flow id, accessor name, capability)` — a name, never a
    data-derived key — and a read accessor "MUST return typed
    tag values or a typed refusal" (`0004…md:257-260`): values,
    not a key-space. RDR 0002's `unknown tag` category then
    rejects any key that did arrive undeclared
    (`0002…md:341-347`). So the data channel carries no
    key-spelling authority; the declaration channel
    (blocks 1–3) is where every spelling is decided, and block
    4's breach is correctly a producer programmer mistake.
  - **Scope of the guarantee (contractual, not yet
    mechanical)**: RDR 0004's accessor layer is **unimplemented
    at HEAD** — `internal/` holds only `cli`, `resolve`, and
    `version`, and 0004's own audit records "Accessor executor
    | None found under `internal/`" (`0004…md:350`). The claim
    therefore rests on two `Final` peer contracts plus their
    committed spikes, not on shipped code. One implementer
    note follows from that: 0004's spike `read` returns the
    whole fixture tag map unfiltered, so the "expected tag
    keys" filter (`0004…md:208-212`) must actually be applied
    when binding a read result into `Input.Owned` for this
    guarantee to survive implementation.
  - **Consequence for A8**: A8's literal sweep was not merely
    incomplete for the data channel — it was sufficient, because
    there was never a second spelling source to sweep. A8's
    scope limit stands as written; A10 is what closes it.

- **A11 RDR 0009's `Resolve`-entry precondition and this
  RDR's compose without either RDR having to reopen: they are
  decidable independently, so only the *reported* error for a
  doubly-breaching input is unsettled — which block 4 leaves
  deliberately open rather than blocking on.**
  - **Status**: Verified (independence); report order remains
    a Stage 7.1 item by design, not an open verification
  - **Method**: Peer RDR
  - **Evidence**: the independence half — the half that could
    have forced a design change — is closed. The predicates
    read **disjoint fields**: this RDR's read-domain is
    `Input.Owned`, `Input.Observed`, and each row's
    `RequiresOwned` (block 4, stated); RDR 0009's predicate is
    "exactly `len(row.Escape) != 0 && len(row.Writes) != 0`,
    applied to every row of the supplied table at `Resolve`
    entry" (`0009…md:536-540`). Both traverse the row set;
    neither touches a field the other reads, and both are pure
    reads that mutate nothing — so neither can change or
    suppress the other's verdict, and both are decidable in
    one pass. Corroborating the disjointness at the document
    level: the string `Input` occurs **zero** times in RDR
    0009 (case-sensitive, whole file, re-counted at
    Reconcile), so 0009's contracts speak only of "the
    supplied table," never of this RDR's subject.
    Neither RDR fixes a total order: 0009 pins its breach
    ahead of *dispositions* ("precedes every modeled
    disposition, `unmodeled_outcome` included") and is silent
    on a peer precondition, exactly as this RDR is silent on
    0009's.
  - **What remains, and why it is not a Pending
    verification**: only which of two errors a
    doubly-breaching caller *sees*. That is a cross-RDR
    decision neither RDR may take unilaterally — not a fact
    either could look up — and block 4 already carries a
    conforming interim rule ("either error is conforming;
    what is *not* conforming is skipping this check because
    another fired"). RDR 0009 is `Draft` at Stage 2/2.5 (its
    evidence tree holds only `research/` and
    `propose-premortem/`), with its predicate form explicitly
    "sharpened at Pre-Lock", so the ordering clause may land
    in either RDR and this one must not bind 0009's symbol
    name. Carried to Stage 7.1 as a reconciliation item.
  - **If wrong** (the predicates are *not* independent — one
    can suppress the other's verdict, rather than merely
    reporting first): the breach this RDR detects could be
    skipped entirely for some inputs, which block 4 forbids,
    and the interim "either error is conforming" licence is
    unsafe. The repair is then a single normative ordering
    clause owned by whichever RDR locks last. Mere
    report-order variation is *not* a falsification — block 4
    already licenses it pending this reconciliation.

- **A12 An RDR 0002 implementer reading 0002 alone will
  encounter this RDR's name constraint, because the
  implementation prompt extracts normative contracts from the
  RDR set rather than from one document.**
  - **Status**: **Refuted as stated** — the assumed mechanism
    does not exist. Superseded by the explicit prerequisite
    below, which is what this RDR now relies on instead.
  - **Method**: Peer RDR + Source Search
  - **Evidence**: three independent reads, all negative.
    (1) The Stage 8 launch prompt implements **one** RDR and
    disclaims the set: "RDRs are run one at a time… **Cross-RDR
    orchestration is out of scope**"
    (`$RDR_HOME/prompts/implementation/launch.md:14-15`); the
    orchestrator's input is a single `{RDR_PATH}`, and Phase 0
    writes `req-list.md` from that one document's testable
    clauses. (2) The only cross-RDR traversal is the
    `Predecessors` field, which is a *completion gate* pointing
    **backward** — it cannot reach a higher-numbered RDR, and
    0002 numerically precedes this one. (3) The `Overrides`
    field is **never read by the engine**: `grep -rn "Overrides"`
    over `$RDR_HOME/prompts/ stages/ skills/` returns no
    matches, so this RDR's Overrides is a durable pointer only
    inside *this* document — a 0002 implementer never opens it.
    Confirming the gap at the destination: RDR 0002 carries
    **zero** occurrences of `0008`, its `[tags.<tag>]` grammar
    block admits any tag name and points nowhere, and
    `reserved_tag_key` is absent from its enumerated category
    list (extensible as *license*, but not named).
  - **Why this is a downgrade, not a blocker**: A12 asserted a
    *discoverability* mechanism, not a design premise. No
    normative contract in this RDR changes — the rule, the
    category, the payload, and the predicate all stand exactly
    as written. What its falsification changes is enforcement
    *reach*: a 0002 implementation that never learns of the
    constraint ships a normalizer accepting any
    `[tags.<name>]`, so the declaration-channel diagnosis
    (blocks 2–3) would be missing while the kernel-side
    predicate (blocks 4–5) still catches the breach at
    `Resolve` entry. The authoring-time diagnosis is the part
    at risk, which is precisely the Problem Statement's
    outcome — so the repair is mandatory, not optional.
  - **Repair, adopted here (this RDR's own "If wrong" branch)**:
    an explicit implementation-ordering prerequisite rather
    than a silent reliance on the launch prompt. Recorded in
    Prerequisites: this RDR's carried half MUST be handed to
    0002's implementation explicitly — 0002's Stage 8 run takes
    this RDR's Normative Contracts as a named input, or a
    pointer lands in 0002's own surface as a cross-RDR edit.
    Which of the two is a cluster-level call, so the choice is
    routed to Stage 7.1; the *obligation to do one of them* is
    settled here and is no longer assumed away.

- **A13 `Input.Owned` and `Input.Observed` are ordered
  sequences of key/value tags, not keyed maps, so duplicate tag
  keys are admissible input the reserved-key predicate must
  scan for rather than look up.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/resolve/resolve.go::Input` declares
    `Owned []Tag` and `Observed []Tag`, where
    `internal/resolve/resolve.go::Tag` is
    `struct { Key, Value string }` — no map anywhere on the
    input surface. `::assemble` consumes them as sequences
    (`for _, t := range in.Observed` then
    `for _, t := range in.Owned`, writing
    `view.tags[t.Key]`), which is where the package doc's
    "within one provenance the last tag wins" behavior comes
    from: a duplicate key is not rejected, it is overwritten.
    The `TagSet` the view exposes *is* map-backed
    (`TagSet{tags: make(map[string]taggedValue, …)}`), which is
    the likely source of the confusion — the assembled view is
    keyed, the input is not.
  - **Why it is load-bearing here**: block 4's predicate is
    specified over `Input.Owned`/`Input.Observed`. Under a map
    reading the reserved-key check is one lookup and
    multiplicity cannot arise; under the actual slice reading
    the same `Input` can carry the reserved key more than once,
    which is precisely the case block 4's first-breach rule and
    scenario 6's third variant now settle. Surfaced by the
    repeatability lens: all three runs independently modeled
    these as maps and dismissed the shape as irrelevant.
  - **If wrong** (the shape changes under RDR 0001): the
    duplicate-key variant of scenario 6 becomes unreachable and
    block 4's scan wording reverts to a lookup — a
    simplification, not a redesign. The first-breach rule is
    unaffected either way, since block 5's `RequiresOwned`
    channel is a `[]string` regardless.

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
   data-level validation failure in a new `reserved_tag_key`
   category — reported by RDR 0002's normalizer before any
   resolution. This converts the seed's silent no-match into
   a typed, named failure at the earliest point the model is
   read. Scoped precisely: this RDR guarantees the failure
   *exists and carries its guidance* at load/lint; which
   user-facing command surfaces it, and under which exit
   code, is RDR 0005's mapping decision and is not settled
   here (charted). No new refusal kind; no kernel
   disposition change for conforming tables.
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
   The locus is settled here: a kernel-exported predicate over
   `Input` that producers may call, applied by `Resolve` at
   entry — one predicate, two call sites (Load-Bearing
   Decisions / Enforcement locus; A6, A8).
5. **Reserve the key in `Row.RequiresOwned` too.** An owned-key
   *reference* is a third way the name arrives: a row naming
   `recognized` there fails `missingOwned`'s owned-only
   provenance test and refuses `owned_state_unavailable`
   naming the reserved key. Because `RequiresOwned` is a
   derived normalizer output with no authored TOML form, this
   is a producer obligation on the normalizer rather than a
   source-lint rule. Under RDR 0007's narrowing to
   `[rule.write]`-derived keys (A7, Verified), RDR 0002's
   `write to non-owned tag` rule already discharges it on that
   path; that discharge is a consequence of the current
   derivation path and does not remove the obligation for any
   other one. So the
   obligation is *also* carried by the same exported `Input`
   predicate item 4 settles, which checks each row's
   `RequiresOwned` — one predicate covering both the data and
   owned-reference channels, so neither rests on a derivation
   argument alone. The field's *meaning* stays RDR 0007's (A7).
6. **Shadowing becomes unreachable on all three channels, not
   re-decided.** The declaration channel is closed by
   load/lint (A5); the data and owned-reference channels by
   the one exported predicate (A6, A7). RDR 0001 D3's
   precedence stays untouched as the kernel's deterministic
   backstop behind all three, for input that reaches
   `assemble` without passing the precondition (A5's
   re-derived scope note).

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
kernel constant for the *spelling* is implementation
latitude the implementer may resolve; this RDR requires only
the behavioral test, matching RDR 0001 D2's "no public
surface" posture and RDR 0009's minimal-export precedent.
That latitude is about the spelling constant only — it is
distinct from the `Input` predicate the Enforcement locus
settles, which this RDR does require the kernel to export.

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
see non-TOML producers), while reserving only the data
channel leaves the load-time diagnosis unfixed. Block 5 is
the same identity rule carried into a derived field. Block 6
is a negative contract (what does not change), not an
independent one.

**Closure argument for "three channels" (not an assertion —
the enumeration was extended twice under review, so it owes
one).** The channels are derived from the code, not
enumerated by inspection of this document: a key becomes
visible to rule evaluation only by entering the assembled
view or by being named in a row field the kernel reads
against that view. `internal/resolve/resolve.go::assemble` is
the sole constructor of the view (A2, whole-package sweep)
and it writes exactly three sources — `in.Recognized`,
`in.Observed`, `in.Owned` (`:151-162`). The first is the
carrier's own injection (block 1); the other two are the data
channel (block 4). Rows reference view keys through
`Row.Match`, `Row.RequiresOwned`, and the guard seam: `Match`
and guard references are name-checked at declaration
(blocks 2–3, since a reference must resolve to a declared
tag), and `RequiresOwned` is the one row field read against
the view *by provenance* rather than by value, which is what
makes it a separate channel (block 5). That is the closure:
three write-sources into the view, plus the one row field
whose read is provenance-sensitive. A fourth channel would
require either a second view constructor or a new row field
the kernel reads against the view — both of which are kernel
changes this RDR would have to re-open anyway.

So the ≥2 split signal is **not** tripped;
`foundational` is earned on the cross-RDR axis (this RDR
binds 0001's carrier to 0002's model), not on contract count.

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
view with no `recognized` key. The conformance test's input
domain is non-empty outcomes. A table whose declared
`outcomes` alphabet contains the empty string would pass the
alphabet gate (`Table.models` is a plain
`slices.Contains`, `internal/resolve/resolve.go:216-218`, and
`assemble` runs before it at `:319-321`) and resolve with no
reserved key in the view; nothing in this RDR or RDR 0002
forbids that alphabet entry. It is left unforbidden rather
than silently assumed away — an empty-string outcome is
degenerate for reasons that are not this RDR's to rule on, and
the binding obligation above is scoped so that such a resolve
is simply outside it, not in breach of it.
```

```normative
A tag declaration with provenance `recognized` MUST be named
`recognized`. A tag declaration with provenance `owned` or
`observed` MUST NOT be named `recognized`. Together these
bound the count: because `[tags.<tag>]` is keyed by tag name,
one model admits at most one declaration named `recognized`,
so at most one recognized-provenance declaration survives
validation — cardinality is a consequence of the naming rule,
not a separate check. A second recognized-provenance
declaration necessarily carries a different name and fails
the naming rule; two declarations spelled `recognized` are a
TOML duplicate-key error in the `malformed TOML` category
before this rule is reached. The scope unit is one `[model]`
(RDR 0002's schema unit); this RDR defines no cross-model or
cross-file cardinality rule.

The bound is an upper one only, deliberately: this RDR imposes
**no lower bound** — a model that declares no
recognized-provenance tag at all is conforming here. Matching
the recognized outcome as a tag is an affordance, not an
obligation (Key Discoveries: outcome gating never reads the
key), so requiring the declaration would break every model
that gates only on `Row.Outcome`. The consequence is stated so
it is not mistaken for coverage: a model whose rows were meant
to match the recognized tag but which declares no
recognized-provenance tag lints clean and refuses `no_match`
at resolve time. That is the same intent-channel gap the
reference-position paragraph below describes and the Failure
Modes section charts — not a naming defect, and not closed
here.

Any violation is a data-level validation failure in the
`reserved_tag_key` category, reported at table load/lint
before any resolution — never a kernel refusal. The
reserved-key comparison is byte-exact on the **post-parse**
key string: case-sensitive, no trimming, no folding (so
`Recognized` is an ordinary, unreserved name, and
`" recognized"` — a whitespace-bearing key — is likewise
ordinary and unreserved). TOML quoting is a surface artifact,
not a name variant: `[tags."recognized"]` and
`[tags.recognized]` parse to the identical key string and are
therefore both reserved. The checked positions are the
`[tags.<tag>]` declaration keys only. Reserved-name uses in
predicate positions (`[rule.match.<tag>]`,
`[rule.guard.all.<tag>]`, `[rule.guard.unless.<tag>]`) need no
separate reserved-key check: a conforming model declares
`recognized` and those references resolve to it, while a model
that does not declare it already fails RDR 0002's `unknown
tag` rule. `[rule.write]` is covered by a *different*
pre-existing rule and must not be folded into the sentence
above: once a model legally declares `[tags.recognized]` the
name is a **known** tag, so `unknown tag` no longer fires on
`[rule.write] recognized = …`; what rejects it is RDR 0002's
`write to non-owned tag` category (`0002…md:346`), because a
recognized-provenance tag is not owned. This RDR adds no write
-position check and depends on that category holding.

What this fall-through does **not** cover, stated rather than
left implicit: a model that declares `[tags.recognized]`
correctly *and* declares some other tag (say `[tags.outcome]`,
provenance `owned`) *and* writes `[rule.match.outcome]` where
it meant the recognized one. Every rule in this RDR passes —
the recognized declaration is named correctly, no owned tag is
named `recognized`, and `outcome` is declared so `unknown tag`
is silent — and the row still never fires. That is the
Problem Statement's symptom surviving inside a conforming
model. It is **not** closed here, and it is not a naming
defect: no name is wrong, a reference points at the wrong
declared tag. Closing it needs a rule over declaration
*intent* rather than declaration *name*, which is the same
out-of-scope heuristic the Failure Modes section charts
(`3amigo/charted.md` C-1). This RDR closes the naming channel
and says so; a reader must not read block 2 as a guarantee
that a conforming model's outcome row fires.
```

```normative
Every `reserved_tag_key` failure — and any undeclared-tag
failure whose offending key is the reserved name — MUST
carry, at the data level, three distinct machine-readable
fields: the offending name as authored, a **remedy name**, and
a stable rule identifier. The rule identifier is a comparable
token, not prose: a golden test asserts it byte-for-byte, and
human-readable wording is the renderer's to choose. A category
consumer may map the failure, but the guidance travels in the
failure data, not the renderer.

The remedy name and the rule identifier are **direction-
specific**, because block 2's two naming rules have opposite
remedies and a single pair would give one of them actively
wrong advice:

- A recognized-provenance declaration under a wrong name must
  be renamed **to** the reserved key. Remedy name: the literal
  `recognized`. Rule identifier:
  `reserved-tag-key/kernel-owned`.
- An owned or observed declaration named `recognized` must be
  renamed **away from** the reserved key — no specific name is
  required, since any non-reserved name conforms. Remedy name:
  the **offending name repeated is forbidden**; the field
  carries the empty string, and the rule identifier
  `reserved-tag-key/author-must-rename` is what tells a
  consumer the remedy is "choose any other name" rather than
  "use this one". A renderer MUST NOT present the reserved key
  as the required name in this direction.

Both rule identifiers sit inside the one `reserved_tag_key`
category — they are the finer key the category-vs-rule
distinction below already contemplates ("A category gains rules
over time"). Scenario 7 asserts the first pair byte-for-byte;
scenario 2's owned-tag half asserts the second.

The category's stable data-level value is the token
`reserved_tag_key`, matching the snake_case discriminator
grammar RDR 0001 uses for `RefusalKind` values
(`no_match`, `owned_state_unavailable`); the spaced form
`reserved tag key` is prose for this document only.

The two tokens are not alternatives and a consumer MUST NOT
choose between them: `reserved_tag_key` is the **category**
discriminator — the value a category-mapping consumer (RDR
0005's exit-code map) keys on, and the only one that
participates in RDR 0002's data-level category set.
The **rule identifier** (`reserved-tag-key/kernel-owned` or
`reserved-tag-key/author-must-rename`, per the direction above)
is carried inside the failure payload, identifying which rule
within the category fired; it is for golden assertions and
remediation lookup, never for category dispatch. A category
gains rules over time, so the rule id is the finer key and the
category id is the stable one. Both are asserted byte-for-byte
by scenario 7 in their respective positions. The Go
type, package, and field names carrying these values are RDR
0002's to choose — this RDR constrains the values and their
distinctness, not the struct. Where the offending key is the
reserved name, this payload requirement extends RDR 0002's
pre-existing `unknown tag` failure; that extension is
authored by this RDR and lands in 0002's implementation.
```

```normative
Producers of kernel `Input` MUST NOT supply an owned or
observed tag keyed `recognized`; the reserved key enters the
assembled view only through `Input.Recognized`. A breach is
a producer programmer mistake — not table data, and never a
new `RefusalKind` — and travels the Go error path RDR 0001
reserves for programmer mistakes.

Enforcement is one predicate at two call sites: the kernel
MUST export a construction-time predicate over `Input` that
producers may call, and `Resolve` MUST apply that same
predicate at entry, returning a non-nil error and no
`Result` disposition on breach. That predicate carries **both**
reserved-key obligations — the owned/observed tag keys of this
block, and block 5's `Row.RequiresOwned` reservation — so the
reserved name has exactly one enforcement point across every
channel it can arrive through. The obligation is NOT inherited
from RDR 0009, whose obligation is a *different rule* (escape
-row shape) even where it reads the same `resolve.Row` values;
the predicate's exported name MUST NOT be bound to any symbol
name from RDR 0009.

The predicate's read-domain is stated here rather than left to
be recovered from block 5's ordering aside: it reads
`Input.Owned`, `Input.Observed`, and each row's
`RequiresOwned` reached through `Input.Table.Rows` — three
sequences, no other `Input` field. `Owned` and `Observed` are
**sequences of key/value tags, not maps**
(`internal/resolve/resolve.go::Input` carries `Owned []Tag` /
`Observed []Tag`, `Tag` being `{Key, Value string}`), so the
check is a scan and a duplicate reserved key is admissible
input rather than a structural impossibility. This is stated
because a map reading makes the reserved-key check a single
lookup and silently forecloses the multiplicity question the
next paragraph settles.

On multiplicity, the predicate reports **the first breach it
finds and returns a single non-nil error**; it does not
aggregate, and the order in which it scans the three sequences
is implementation latitude. One breach is sufficient to reject
the `Input`, and a producer holding a programmer mistake is
not owed an exhaustive list. This is distinct from — and must
not be read out of — the "detected whenever present" clause
below, which governs this predicate against *RDR 0009's*
separate precondition, not the reporting of multiple breaches
within this one.

RDR 0009 writes a structurally identical clause over the same
entry point — "the conformance predicate MUST be exported by
the kernel package … and Resolve's entry precondition MUST be
that same function — one predicate, two call sites"
(`0009…md:510-517`) — over a disjoint subject (`resolve.Row`
shape, not `Input`). Two entry preconditions therefore land on
one `Resolve`, and neither RDR alone can settle their order.
This RDR binds only what it owns. The two predicates decide
**different questions over different fields** (this one:
reserved-key occupancy in `Input.Owned`/`Input.Observed` and
in each row's `RequiresOwned`; RDR 0009's: escape/writes
co-occupancy in each row's `Escape`/`Writes`), so neither can
change the other's verdict and both are decidable in one pass.
What this RDR requires of the implementation is therefore only
that its own breach be **detected whenever present** — never
skipped because another precondition also fired.

Which breach a doubly-breaching caller *sees* is an ordering
question, and `Resolve` returns a single `error`, so one of
the two necessarily reports. That order is a cross-RDR
decision neither RDR may take unilaterally; it is reconciled
at Stage 7.1 (A11). Until then the implementable rule is:
**apply this predicate at entry and report its breach; if the
table also breaches RDR 0009's shape rule, either error is
conforming.** A doubly-breaching input is a producer with two
programmer mistakes, not a case whose exact error text this
RDR owes a guarantee on. What is *not* conforming is skipping
this check because another fired. Stage 7.1 may narrow this
to a fixed order; nothing here forecloses that, and no
implementer needs to invent one to proceed.

The check is unconditional on `Input.Recognized`: a
reserved-keyed owned or observed tag is a breach whether or
not the resolve carries an outcome, matching block 1's
unconditional reservation. The cost of that choice is stated
plainly rather than assumed away — when `Input.Recognized` is
empty `assemble` injects nothing
(`internal/resolve/resolve.go:151`), so there is no collision
to prevent and the error is a *reservation* being enforced,
not a shadowing being averted. It is kept unconditional
because a key whose reservation lapses per-call is not
reserved: a producer would have to know whether an outcome is
in flight to know whether its own tag key is legal, which is
exactly the cross-side coupling this RDR exists to remove.
The empty-outcome variant is pinned as expected by scenario 6
for that reason, not by symmetry with block 1's prose.
```

```normative
A normalized row's `Row.RequiresOwned` MUST NOT name
`recognized`. This is a name reservation only: what
`RequiresOwned` *means* is owned by RDR 0007, which narrows
it to post-guard write-dependency keys and pins
`resolve.go::missingOwned`'s provenance-specific
`view.has(key, ProvenanceOwned)` test as deliberately
distinct from provenance-blind guard presence; this RDR
cites that rule rather than restating it, per RDR 0007's own
"peers cite rather than restate" convention.

`RequiresOwned` has no authored source form — it is a
*derived* normalizer output, and RDR 0002's field layout
carries no `requires_owned` key — so this reservation is a
**producer obligation on the normalizer**, not a source-lint
rule over authored TOML, and it does not share the
declaration channel's `reserved_tag_key` category. Under RDR
0007's narrowing the field derives from `[rule.write]` keys,
where a `recognized` entry is already rejected by RDR 0002's
pre-existing `write to non-owned tag` category (A7,
Verified), so the obligation is discharged by construction on
today's derivation path, and it binds any future path that
does not route through that rule.

An obligation discharged by construction still needs an
artifact that fails when the construction changes, or it is
text nothing executes. So this block carries its own
enforcement rather than resting on the derivation: the
**exported `Input` predicate of block 4 MUST also reject a
`Row.RequiresOwned` entry naming `recognized`** on the rows of
the supplied table, on the same Go error path and with the
same producer-breach semantics. This costs no new surface —
it is the predicate block 4 already requires — and it converts
block 5 from an unfalsifiable MUST into a check that fails the
moment any derivation path, present or future, emits the
reserved name. Scenario 9 asserts it alongside the kernel-side
residual. Note the ordering consequence for A11: this widens
the `Input` predicate to read `in.Table.Rows`, so it is no
longer decidable without reading the table, and the
independence half of A11 must be re-checked against RDR
0009's row-shape precondition rather than assumed.

The residual this closes is kernel-side: a row carrying
`RequiresOwned: ["recognized"]` finds the key present under
`ProvenanceRecognized`, fails the owned-only test, and yields
`owned_state_unavailable` naming the reserved key — the
confusing-refusal shape this RDR exists to eliminate, one
field over (A7). With the predicate applied at `Resolve`
entry that refusal is now reachable only by a caller that
reaches `assemble` without passing the precondition (the
package-internal path scenario 9 exercises), not by any
caller going through `Resolve`. It stays a documented
residual on that narrowed path, consistent with the
non-conforming-input posture blocks 4 and 6 take.
```

```normative
Kernel disposition is unchanged by this RDR for conforming
input: no new `RefusalKind`, no change to RDR 0001 D3's
provenance precedence (`owned` > `observed` > `recognized`),
and no behavior change for any conforming table and
conforming input. D3 remains the deterministic backstop
behind the producer obligation for input that bypasses the
precondition. Breaching input gains a non-nil Go error where
it previously resolved — a programmer-mistake path, not a
disposition change for any conforming caller.
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
  `[tags.recognized]`, so it is reserved.

  **Advisory warning — normative, not deferred.** The
  byte-exact rule *creates* the hazard an advisory mitigates,
  so the two are one decision and are settled together.
  A declaration named `Recognized` or `" recognized"` is a
  legal ordinary tag that a human reader cannot distinguish
  from the reserved word — the whitespace form is invisible in
  a diff — so an author can declare what looks like the
  reserved key, lint clean, and watch the row never fire: the
  Problem Statement's failure reproduced *by* the fix. The
  identity rule stays byte-exact (a validator that folds cannot
  tell the author which spelling it wants, and folding would
  make the reserved set unbounded). The mitigation is
  therefore additive rather than a change to identity: a
  declaration whose post-parse key is not `recognized` but
  becomes `recognized` under **either** Unicode-simple case
  folding **or** trimming of leading/trailing whitespace, or
  both applied together, MUST raise a **non-blocking advisory**
  naming both spellings. The trigger is deliberately
  disjunctive: `Recognized` (folding alone) and `" recognized"`
  (trimming alone) are each near-misses on their own, and
  scenario 4 requires the advisory on both, so a conjunctive
  reading would fire on neither. Advisory, not a failure,
  because the name is legal and this RDR must not reject a tag
  it does not own. Scenario 4 pins the dispositions; the
  advisory is asserted there alongside them.

  **Advisory payload** — the advisory carries the same
  machine-readable discipline as the failure payload (block 3),
  because a normative artifact a golden test cannot compare
  byte-for-byte is not testable: two distinct fields, the
  **authored spelling** and the **reserved spelling** (the
  literal `recognized`), plus a stable rule identifier whose
  value is the literal `reserved-tag-key/near-miss`. It travels
  on an advisory channel distinct from the validation-failure
  list — it carries **no** `reserved_tag_key` category
  discriminator, because it is not a validation failure and
  MUST NOT participate in category dispatch or alter the
  load/lint verdict. The carrier's Go type and field names are
  RDR 0002's to choose, on the same footing as the failure
  payload; this RDR constrains the values and the channel's
  separateness.
- **Naming** — canonical name `recognized`, matching the
  shipped `recognizedTagKey` constant and RDR 0001's frozen
  fixture. Rejected: a sigil-guarded name (`_recognized`, or
  an angle-bracket token like RDR 0002's `<clear>`) — it
  would repin the shipped constant and fixture for no added
  safety once the reserved-key validation exists, and it
  diverges from the vocabulary authors already see in
  `Row.Outcome` diagnostics.
- **Enforcement locus (settled at Pre-Lock: (b)+(c) paired)** —
  where the `Input` producer obligation is *checked*. The
  declaration channel is settled (RDR 0002 load/lint); the
  data channel is now settled here. Three candidates were
  weighed: (a) an unchecked documented obligation on
  producers, with D3 as the only backstop — zero new surface,
  no detection; (b) a kernel-entry precondition on `Resolve`
  returning the Go error RDR 0001 reserves — detection at the
  real boundary, at the cost of a narrow kernel-side check;
  (c) an exported construction-time predicate over `Input`
  that producers call, matching RDR 0009's exported-predicate
  form.

  **Chosen: (b) and (c) together — one predicate, two call
  sites.** (a) is refuted by in-repo evidence rather than by
  preference: `internal/cli/respond/respond.go:22-25`
  reserves the terminal `"ok"`/`"failed"` type names with no
  enforcement whatever, and nothing validates a future
  `Stream` emitter — a live instance of how (a) ages (see the
  sibling-path check). A6's Residual makes the cost concrete:
  under (a) resolve-time data keyed `recognized` still
  shadows the outcome under D3, detected by nothing. Between
  the survivors, (b) is the only locus that sees *every*
  producer including non-TOML callers, and it lands on
  reserved-but-unexercised surface —
  `internal/resolve/resolve.go::Resolve` already returns
  `(Result, error)` whose error is doc-reserved for
  programmer mistakes with no non-nil path in its body today,
  so (b) adds a check without widening the signature. (c)
  alone cannot see a producer that never calls it. Pairing
  them is RDR 0009's own shape ("one predicate, two call
  sites") and is what stops the two enforcement points from
  drifting: the exported predicate is the single definition,
  and `Resolve` calls it at entry.

  Surface conceded, stated plainly: one exported predicate
  over `Input` plus one kernel-entry call. The "no kernel
  change" claim in the Decision Rationale is therefore scoped
  to *disposition for conforming input*, not to zero surface.
  The predicate's exported name is implementation latitude —
  do **not** bind it to any symbol name from RDR 0009, which
  is still `Draft` with its predicate form explicitly
  "sharpened at Pre-Lock". Breach of the precondition returns
  a non-nil error and no `Result` disposition; conforming
  callers see no behavior change.

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Reserved key injected into the assembled view | Predecessor (RDR 0001 D2, implemented) | Available | Ratified from latitude to contract; no code change |
| Tag declarations with provenance | Predecessor (RDR 0002, Final unimplemented) | Deferred | Name constraint on `recognized`-provenance declarations |
| Load/lint name validation | Predecessor (RDR 0002 normalizer) | Deferred | One additive validation category: `reserved_tag_key`; upstream of RDR 0006's normalized-graph lint, not a competing acceptance rule |
| Input-boundary producer enforcement | This RDR (pattern borrowed from RDR 0009) | Available | Settled at Pre-Lock: exported `Input` predicate + `Resolve` entry call, covering owned/observed tag keys **and** each row's `RequiresOwned` (Load-Bearing Decisions / Enforcement locus; A6, A8). Ordering vs RDR 0009's row-shape precondition is A11 |
| `RequiresOwned` name reservation | This RDR; field meaning stays RDR 0007's | Available (A7 Verified) | Producer obligation on the normalizer, not a source-lint rule — the field has no authored form. Discharged on the current derivation path by 0002's `write to non-owned tag`; the exported `Input` predicate enforces it on any other path. Cite 0007, do not restate its semantics |

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
| Blast radius | no disposition change for conforming input, but **not zero surface**: one exported `Input` predicate plus its `Resolve` entry call (the kernel's first non-nil error path, on error surface RDR 0001 already reserves — A6, A8), one additive lint category + declaration constraint in 0002, and a rename of three committed fixtures plus six reference sites (A4) | new public kernel surface (`Input` key field or parameterized `assemble`), normalized-table surface, threads through 0005; frozen fixture repinned | violates locked REQ-17 of implemented Final RDR 0001; removes guard-visible recognized value that RDR 0007's chosen guard domain reads |
| Reversibility | high — a later RDR could still parameterize; reserving now forecloses nothing | low — public surface, once shipped, is load-bearing | low — deleting shipped kernel behavior and reopening 0001 |
| Consistency with 0007/0009 posture at this seam | same shape as 0009: producer/lint obligation, kernel unchanged for conforming input, no new refusal kind | cuts against both peers: grows kernel surface to solve an authored-data problem | cuts against 0007, whose guard domain assumes the assembled view carries the recognized tag |
| Cost | doc ratification + lint checks inside work 0002 already owes | largest: kernel, table schema, CLI threading, fixtures | medium code deletion + cross-RDR contract reopening |

Deciding rows: **correctness fit**, **blast radius**, and
**prior-art alignment**. Correctness fit is named first
because it carries the user outcome ("my row fires, or I am
told why", Problem Statement): R is the only branch that
both keeps the match working and makes the failure
diagnosable — M defers misdeclaration to resolve time,
reproducing the silent discovery in a new form, and W deletes
the match outright. R also has no kernel disposition change
and positive prior art, where M has a strict negative
prior-art result and the largest surface growth, and W
reopens a locked, implemented contract (REQ-17) and breaks
RDR 0007's guard-domain assumption. R therefore wins on the
user outcome *and* on cost, rather than on cost alone; the
consistency row is corroborating, not decisive.

The blast-radius row was **re-scored after** the Enforcement
locus settled to (b)+(c), since the fork was originally
weighed against an R that changed no kernel code. R now
carries one exported predicate, the kernel's first non-nil
error path, and a three-fixture rename (A4). The ranking is
unchanged — M still grows `Input` or `assemble` *and* the
normalized-table surface *and* threads a name through 0005,
strictly dominating R's addition, and W still reopens a
locked implemented contract — but the margin is narrower than
the pre-settlement row implied, and R is no longer the
zero-kernel-change branch it was chosen as. Recorded here so a
reader weighing the matrix sees the cost that actually ships.

The same user-outcome criterion settles the Enforcement
locus below: candidate (a) is the cheapest on blast radius
but leaves the data-channel shadowing detected by nothing
(A6's Residual), which fails the outcome this RDR commits
to — so cost does not get the last word there either. Rejected alternatives are analyzed
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
- Deletes the capability itself — matching the recognized
  outcome in tag vocabulary *at all*. This is a strictly
  larger loss than the chosen branch's: R withdraws the
  author's choice of *name* for that tag while keeping the
  match working, whereas W removes the match. Against the
  outcome this RDR commits to ("my row fires, or I am told
  why"), R preserves it and W does not.

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
- Positive: the shipped constant and D3 precedence remain
  valid, and the branch stays the easiest to reverse (a later
  RDR can still parameterize the key). Not "cheapest to
  implement" without qualification: the settled enforcement
  locus adds an exported predicate and the kernel's first
  non-nil error path, and A4's inventory adds a rename of
  three committed fixtures and six reference sites. Cheapest
  *relative to M and W*, which is what the QOC matrix ranks.
- Positive: a consistent trust story at the 0001↔0002 seam
  with RDR 0009 (normalizer enforces, no new refusal kind,
  kernel disposition unchanged for conforming input).
- Negative: authors lose naming freedom for exactly one
  tag; `recognized` becomes a reserved word they must learn.
  Mitigated by the failure *data* — offending name, required
  name, rule identifier — which is normative and tested
  (scenario 7). The human wording built from those fields is
  the renderer's, and this RDR neither specifies nor tests
  it: a first-time author's comprehension is therefore
  asserted at the data level only. Whether that suffices is a
  question for whichever CLI surface renders it (charted with
  the user-facing-surface item).
- Negative: RDR 0002's implementation grows one validation
  category and the declaration name constraint it did not
  originally spec — additive, but still scope on an
  unimplemented Final RDR.
- Negative: the kernel gains a narrow public surface it did
  not have — one exported `Input` predicate plus its
  `Resolve`-entry call (A8) — so the "no kernel change" story
  holds for *disposition on conforming input*, not for zero
  surface. This is the price of detecting the data-channel
  breach at all; the alternative (a documented-only
  obligation) was rejected because nothing would catch it.

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
  data channel, checked at a settled locus — a kernel-exported
  `Input` predicate applied by `Resolve` at entry (A6, A8), so
  the breach is detected rather than left to D3 alone. D3
  stays the deterministic backstop behind it for input that
  bypasses the check; the breach test is scenario 6.
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
  the data level — offending name, remedy name, rule id —
  so a generic renderer still surfaces the resolution; a
  golden-text test is named in Phase 3.
- **Risk**: near-spellings (`Recognized`, quoted/whitespace
  TOML key forms) read as the reserved word to a human but
  are ordinary names to the byte-exact rule (premortem
  P-10).
  **Mitigation**: the comparison rule is normative
  (byte-exact on the post-parse key string) *and* the
  non-blocking near-miss advisory is now normative alongside
  it (Load-Bearing Decisions / Identity) — settled at Pre-Lock
  rather than left as a Resolve question, because the
  byte-exact rule is what creates the hazard. The whitespace
  form is the dangerous half: invisible in a diff. Scenario 4
  pins both the dispositions and the advisory.
- **Risk**: pre-existing tables using `recognized` as an
  innocent owned/observed tag break at load with no
  migration story (premortem P-9).
  **Mitigation**: none needed as migration for *that* clause —
  no owned/observed tag is named `recognized` at HEAD (A4,
  second clause), the project carries no back-compat
  obligation, and the innocent-collision case is a lint
  fixture whose refusal text tells the author what to rename.

- **Risk**: the rule's *first* clause breaks the peer's own
  worked examples on day one — all three committed
  recognized-provenance declarations are wrong-named under it,
  and RDR 0002 designates them "the canonical examples
  implementation tests must promote" (`0002…md:361`). The
  reserved-key rule's debut is a red test on the reference
  model (A4, first clause).
  **Mitigation**: the migration is named work, not a
  discovered surprise — A4 carries the inventory (three
  declarations, six tag-predicate reference sites)
  and Phase 2 carries the edits. The rename is **mechanical**,
  settled at Reconcile: A9 is Verified against RDR 0002's own
  normalizer spike, which puts every tag predicate — match and
  both guard positions — into the candidate row's predicate
  set and produces no outcome field, so no gating semantics
  move with the rename. The fixture work is no longer gated.

### Failure Modes

Visible: a wrong-named recognized-provenance declaration (a
second such declaration is wrong-named by construction), or
a reserved-name owned/observed declaration, fails table
load/lint in the `reserved_tag_key` category — the failure
data carries the direction-specific rule identifier, the
offending name, and the remedy name, before any resolution
runs. A producer
supplying an owned/observed input tag keyed `recognized`, or a
row naming it in `RequiresOwned`, breaches the producer
obligation and surfaces on the Go error path as a programmer
mistake — caught by the exported `Input` predicate or by
`Resolve` at entry (A6), never as a modeled refusal. The
programmer-mistake channel is the right one for that breach:
an *accessor-produced* owned key cannot introduce a spelling
of its own (A10, Verified) — the key is always the declared
`[tags.<tag>]` name and the artifact supplies only the value —
so no user data is reclassified as a programmer mistake.

Silent, residual — **and not a narrow path**: an author who
declares an *observed* tag under an innocent name (`result`,
`outcome`) intending recognizer semantics passes every name
check and the row never fires, because intent is invisible to
a name-and-provenance validator (premortem P-3). This RDR's
own evidence shows the shape is common rather than
exotic: both canonical RDR 0002 fixtures name the
recognized-provenance tag `outcome`
(`0002-…/evidence/spikes/rdr-fixture.toml`,
`kata-fixture.toml`), so every existing example of an author
naming this datum chose an innocent name. The reserved-key
rule catches the author who got the *provenance* right and
only the *name* wrong; it does not catch the author who
reached for `provenance = "observed"` instead.

Scoped honestly: this RDR closes the naming channel, not the
intent channel. The intent channel needs a heuristic over
declaration *shape* — a model whose rows gate on outcomes but
which declares no recognized-provenance tag — which is a
different rule with a different false-positive profile, and
is **out of scope here** (charted; see the 3amigo evidence
folder). Diagnosis meanwhile: the table dump shows no
recognized-provenance declaration for the model while rows
gate on outcomes. Silent,
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
- [ ] **The carried half is handed to 0002 explicitly** — this
      RDR's Normative Contracts are a named input to RDR
      0002's Stage 8 run, or a pointer to this RDR lands in
      0002's own surface as a cross-RDR edit. Not optional and
      not automatic: A12 established that no engine mechanism
      carries a later RDR's `Overrides` to a `Final` peer's
      implementation (the launch prompt implements one RDR and
      puts cross-RDR orchestration out of scope; `Predecessors`
      points backward only; `Overrides` is read by nothing).
      Without this, 0002 ships a normalizer that accepts any
      `[tags.<name>]` and the authoring-time diagnosis — the
      Problem Statement's outcome — never arrives. Which of
      the two forms is used is a Stage 7.1 call; that one of
      them happens is a gate here.

### Minimum Viable Validation

Two halves, matching the Done split in Testing Strategy.

**Kernel half — executed during this RDR's implementation.**
A row matching on the tag `recognized` fires against a
recognized outcome, asserted through the real kernel's
assembled view (the behavioral conformance form, premortem
P-5), and an `Input` supplying an owned or observed tag keyed
`recognized` — or a row naming it in `RequiresOwned` — is
rejected by the exported predicate and by `Resolve` at entry.
Both run against `internal/resolve` as it stands. This half
alone does not deliver the Problem Statement's outcome; see
Testing Strategy's note on the split.

**Normalizer half — executed inside RDR 0002's
implementation.** A table declaring a recognized-provenance
tag named `recognized` loads and lints clean; the same table
with the declaration renamed (and separately, with an owned
tag named `recognized`) fails load/lint in the
`reserved_tag_key` category whose failure data carries the
direction-appropriate remedy name and rule identifier — not a
silent no-match at resolve time. This
half is gated by the Prerequisite below, not deferred by
choice.

### Phase 1: Contract ratification

Promote D2's key to normative: record the reserved-keyword
rule where implementers read it — this RDR's Normative
Contracts as authority, plus a pointer comment on
`internal/resolve/resolve.go::recognizedTagKey` citing RDR
0008. No kernel behavior change.

### Phase 2: Normalizer name validation

Add the reserved-key checks to RDR 0002's load/lint path
over the `[tags.<tag>]` declaration keys: wrong-named
recognized-provenance declarations and reserved-name
owned/observed declarations report the `reserved_tag_key`
data-level category with its three-field payload, plus the
non-blocking near-miss advisory (Load-Bearing Decisions /
Identity). The `RequiresOwned` reservation needs no
source-lint check here (A7, Verified) — the field has no
authored source form, and on RDR 0007's `[rule.write]`
derivation it is already covered by 0002's `write to non-owned
tag` rule. Its enforcement is the exported `Input` predicate
(block 5), which binds any future derivation path.

Also in Phase 2: rename the three committed
recognized-provenance declarations and **all six** of their
tag-predicate reference sites (A4's inventory — five
`[rule.match.outcome]` plus `[rule.guard.all.rewind_target]`
at `guard-fixture.toml:72`). A declaration renamed without
every one of its references leaves an undeclared-tag
reference and trips 0002's `unknown tag`. The rename is
mechanical and no longer gated: A9 is Verified, so every one
of those references normalizes into the candidate row's
predicate set and no gating semantics move. Expect the
renamed fixtures to become *newly* matchable against the
kernel (the key `outcome` cannot match `assemble`'s injected
`recognized` today) — an improvement, and safe because no
code path at HEAD normalizes these fixtures into a
`resolve.Table`.

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
`recognized`, asserted at both the exported predicate and
`Resolve`'s entry call — A6), the guard/matcher same-view
test (A2), the two-recognized-declaration fixture (two
naming failures), the reserved-key normalization fixtures
(case and whitespace variants stay unreserved; the quoted
form is reserved) **plus the near-miss advisory assertions**
(Load-Bearing Decisions / Identity), the kernel-side
`RequiresOwned` refusal test **and its predicate-side breach
assertion** (block 5; A7), the conforming-`Input` nil-error
pin including the existing empty-`Input{}` case
(`resolve_test.go:748`), and the golden failure-data check
(offending name, remedy name, rule identifier — both
directions, since block 3 pins a distinct remedy/rule pair for
each).

## Validation

### Testing Strategy

The matrix the verified assumptions imply. Scenarios 1–2 are
the MVV pair; 3–9 are the premortem- and assumption-derived
set named in Phase 3. Scenarios 2, 4, 5, and 7 land inside
RDR 0002's normalizer work (no normalizer exists at HEAD —
reuse audit); 3, 6, 8, and 9 are kernel-side and runnable
against `internal/resolve` as it stands; 1 has a half on each
side (its kernel assertion runs today, its load/lint half
does not).

Done is defined in two parts, because this RDR's contract
spans a kernel that exists and a normalizer that does not.

**Done for this RDR's own implementation** (all runnable at
HEAD): scenarios 1's kernel half, 3, 6, 8, and 9 are green,
and Phase 1's pointer comment has landed. Scenario 3's
guard/matcher same-view property is pinned by a test that
would fail if the two views diverged — currently true by
construction but unpinned (A2). Scenario 6 requires the
exported `Input` predicate and the `Resolve` entry call,
both settled at Pre-Lock (Load-Bearing Decisions /
Enforcement locus).

**Carried into RDR 0002's implementation** as a written
obligation this RDR authors and 0002 executes: scenarios 2,
4, 5, 7, and the lint half of the MVV. These cannot go green
before a normalizer exists, and the Prerequisites gate says
so. They are not deferred in the Finalization Gate's sense —
they are *scoped to the peer that owns the code*, and this
RDR's Normative Contracts are the authority 0002's
implementation prompt extracts.

**What this split means, stated rather than left to be
discovered**: every scenario that delivers the Problem
Statement's outcome — the typed load/lint failure replacing
the silent no-match — is in the carried half. This RDR's own
half ratifies the key, adds the `Input`/`RequiresOwned`
predicate, and pins view plumbing; none of it is perceptible
to a table author. So this RDR can reach "implemented" with
its user-facing outcome still unshipped, and that is a real
property of the split, not an accounting artifact.

Two consequences are binding rather than advisory. First, the
Prerequisite "RDR 0002 implementation underway" is a **gate on
declaring this RDR's problem solved**, not merely on running
the carried tests: this RDR's Close MUST NOT claim the
Problem Statement outcome until the carried half is green.
Second, the implementation-ordering fact — that 0008 alone
changes nothing an author perceives — belongs in the Stage 7.1
cluster read, because it determines whether locking 0008
before 0002 is implemented buys anything. If cluster-reconcile
finds it does not, the honest disposition is to hold this RDR
at Final-unimplemented until 0002's work starts rather than
implement a kernel predicate in isolation.

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
   **Expected**: both fail load/lint in the
   `reserved_tag_key` data-level category before any
   resolution — not a silent no-match at resolve time. The
   category is additive under RDR 0002's "including at
   minimum" fence (A1).

3. **Scenario** (A2 — the net-new coverage): a guard
   predicate and a match pattern both read the recognized
   tag in one resolve, through a **new** view-capturing guard
   seam added alongside
   `internal/resolve/fixtures_test.go::fixtureGuards` — not a
   change to it, since ~10 existing tests depend on its
   current shape.
   **Expected**: the captured guard-side view satisfies
   `Lookup("recognized") == (in.Recognized,
   ProvenanceRecognized, true)`, and `Len()` equals the key
   count the matcher's row pattern implies — the two
   observables `TagSet` exports. A row whose `Match` names
   the recognized tag is selected in the same resolve, so a
   divergence between the guard's view and the matcher's
   would fail one of the two assertions. This is net-new: the
   suite's only guard seam today discards the view
   (`fixtureGuards.Evaluate` takes `_ resolve.TagSet`), so no
   existing test would catch a guard/matcher view split
   (premortem P-7).

4. **Scenario** (identity/exactness): near-spellings of the
   reserved word as declared tag names — `Recognized`,
   `RECOGNIZED`, a quoted TOML key `"recognized"`, and a
   whitespace-bearing `" recognized"`.
   **Expected**: the byte-exact rule on the post-parse key
   string treats `Recognized`, `RECOGNIZED`, and
   `" recognized"` as ordinary unreserved names (no folding,
   no trimming); the quoted `"recognized"`, which parses to
   the identical key string as the bare form, **is**
   reserved. Additionally, each of the three unreserved
   near-misses raises the **non-blocking advisory** while
   `[tags.result]` — an ordinary name that is not a near-miss —
   raises none, so the advisory is pinned as targeted rather
   than blanket. The three cover the trigger's disjunction
   deliberately: `Recognized` and `RECOGNIZED` qualify by case
   folding alone and `" recognized"` by trimming alone, so a
   conjunctive (fold-*and*-trim) implementation fails this
   scenario on all three. The advisory's payload is asserted
   byte-for-byte like the failure payload: authored spelling,
   reserved spelling `recognized`, rule identifier
   `reserved-tag-key/near-miss`, and **no** `reserved_tag_key`
   category discriminator. The advisory MUST NOT change the
   load/lint verdict for any of them. Pins the Load-Bearing
   Decisions / Identity rule including its advisory clause
   (premortem P-10).

5. **Scenario** (declaration cardinality, as a consequence of
   the naming rule): a model carrying two
   recognized-provenance declarations — necessarily two
   different names, e.g. `[tags.outcome]` and
   `[tags.result]`, since `[tags.<tag>]` is name-keyed.
   **Expected**: both fail in the `reserved_tag_key` category
   on the naming rule (two failures, not one), leaving at most
   one surviving recognized-provenance declaration — which
   matches RDR 0001's single `Input.Recognized` per resolve
   (Key Discoveries: singleton by inheritance). The
   same-name variant is out of scope for this category:
   `[tags.recognized]` twice is a TOML duplicate-key error in
   `malformed TOML`, asserted as such.

6. **Scenario** (A6 — producer obligation): kernel `Input`
   carrying an owned or observed tag keyed `recognized`,
   constructed directly (bypassing lint, as a non-TOML
   producer would), in three variants — one with a non-empty
   `Input.Recognized`, one with it empty, and one carrying the
   reserved key **twice** in the same tag sequence. The third
   variant exists because `Input.Owned`/`Observed` are
   `[]Tag` slices, not maps
   (`internal/resolve/resolve.go::Input`), so duplicate keys
   are admissible input; it pins block 4's first-breach rule —
   the predicate returns one non-nil error, not two, and the
   test asserts a single error rather than an aggregate.
   **Expected**: all three variants breach. `Resolve` returns a
   non-nil `error` and a zero-valued `Result` (no `Plan`, no
   `Refusal`) — never a new `RefusalKind`, and never a
   modeled disposition. The exported predicate, called
   directly on the same `Input`, reports the same breach:
   one predicate, two call sites, asserted at both. A
   conforming `Input` returns a nil error, pinning that the
   check adds no behavior change for conforming callers —
   including the empty `Input{}` that
   `internal/resolve/resolve_test.go:748` already pins as a
   nil-error resolve, which the new precondition MUST keep
   green. The empty-`Input.Recognized` variant is expected to
   breach on the reserved key alone (block 4's unconditional
   clause), which is what distinguishes it from the empty
   `Input{}` case that carries no reserved key at all.

7. **Scenario** (golden failure data; premortem P-4/P-11):
   **both directions** of the naming rule, each passed to a
   synthetic consumer that maps only the categories RDR 0002
   enumerates today (i.e. treats this one as unknown and falls
   through to a generic branch).
   **Expected**: the failure **data** still carries all three
   fields, asserted byte-for-byte and independent of anything
   the consumer renders. Wrong-named recognized-provenance
   declaration: offending name as authored, remedy name
   `recognized`, rule identifier
   `reserved-tag-key/kernel-owned`. Owned/observed declaration
   named `recognized`: offending name `recognized`, remedy name
   the empty string, rule identifier
   `reserved-tag-key/author-must-rename` — pinning that the
   reserved key is **not** presented as a required name in the
   rename-away direction (block 3). Both carry the same
   category discriminator `reserved_tag_key`. The consumer
   is a test stub, not RDR 0005's exit-code map: this RDR
   asserts the payload contract, and 0005 owns whatever
   mapping it later adds.

8. **Scenario** (A5 — backstop determinism): a
   hand-constructed non-conforming `Input` whose observed or
   owned tags include the key `recognized`, resolved with the
   producer precondition bypassed (calling `assemble`'s
   behavior through the package-internal test path).
   **Expected**: D3's precedence (`owned` > `observed` >
   `recognized`) resolves the collision deterministically —
   the same input yields the same disposition on repeat runs.
   Scope note: the unreachability half of A5 is **not** a
   test. It is a derivation over RDR 0002's declaration
   semantics (A5's P1+P2+P3), and the overwrite in
   `internal/resolve/resolve.go::assemble` has no observable
   at the package boundary — no counter, no error, no
   distinguishing output. A5's derivation plus scenario 2's
   lint rejection are the evidence that conforming input
   cannot reach it; this scenario pins only what a test can
   actually observe.

9. **Scenario** (A7 — the `RequiresOwned` channel): a
   hand-constructed row carrying
   `RequiresOwned: []string{"recognized"}` evaluated against a
   view whose `recognized` key is present under
   `ProvenanceRecognized`, reached through the
   package-internal path that bypasses the entry precondition
   (as scenario 8 does) — since a caller going through
   `Resolve` is now stopped by the predicate before this
   refusal can occur.
   **Expected**: the key is reported missing and the resolution
   refuses `owned_state_unavailable` with `MissingOwned`
   containing `recognized` — pinning that the unreserved case
   is reachable, not hypothetical (`::missingOwned` tests
   `view.has(key, ProvenanceOwned)`; `::TagSet.has` is
   provenance-specific). Runnable against `internal/resolve`
   as it stands. **Second half (block 5's enforcement)**: the
   same row passed to the exported `Input` predicate reports a
   breach, and `Resolve` at entry returns a non-nil error
   rather than the `owned_state_unavailable` refusal above —
   so the refusal shape is what a *bypassing* caller sees,
   while a caller going through the predicate is stopped
   first. Both are asserted; the pair is what makes block 5
   falsifiable rather than discharged-by-assertion. There is
   **no source-lint half** (A7, Verified): `RequiresOwned` has
   no authored source form, so the normalizer-side source
   check is unnecessary and the corresponding negative test
   belongs to RDR 0002's `write to non-owned tag` suite, not
   this RDR's.

## Finalization Gate

Responses: `0008-recognized-tag-key-ownership/artifacts/gate.md`
(Gate PASS 2026-08-11)

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
