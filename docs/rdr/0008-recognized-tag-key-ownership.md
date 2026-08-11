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
- **Profile**: foundational — provisional; one contract
  (which side owns the *name* of the recognized-provenance
  tag key in the assembled view), but it binds RDR 0001's
  kernel carrier to RDR 0002's declared model, so it spans
  two already-Final RDRs. Resolve overwrites from the
  verified count.
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
  disjoint namespace. XState: guards receive
  `{ context, event }`
  (`repos/xstate/packages/core/src/guards.ts::GuardArgs`) —
  the injected event occupies the framework-fixed `event`
  slot; authors name event *types*, never the slot. Inngest:
  "Read-Only Fields … system-managed and cannot be modified
  through write operations"
  (`repos/inngest/pkg/api/v2/README.md`). In-repo: RDR 0002
  already reserves the `<clear>` write sentinel — a
  carrier-owned token adjacent to author-named vocabulary.
  Negative result: no surveyed engine lets authors rename
  the injected event/outcome slot.
- **Assumed** — W3C SCXML §5.10 reserves system-variable
  names (`_event` etc.) against author binding. The spec
  text is in no local corpus/checkout; demoted to A3 rather
  than leaned on.

### Critical Assumptions

- **A1 RDR 0002's data-level validation-category list is
  additively extensible without reopening its Final
  contract.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: RDR 0002 Normative Contracts, the
    "Validation failures MUST retain stable data-level
    categories … including at minimum" fence — "at minimum"
    marks the list as a floor, not a closed enum.
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
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: `internal/resolve/resolve.go::assemble`
    (sole injection point) and
    `internal/resolve/resolve.go::GuardEvaluator` (guards
    evaluate over the same `TagSet` the matcher reads); to
    be confirmed by sweeping `internal/resolve` for any
    other read of `Input.Recognized` exposed to table data
    and pinned by the same-view test named in Phase 3.
  - **If wrong**: reserving the view key leaves a second
    unbound name elsewhere (or a guard/matcher view split —
    premortem P-7) and the ownership gap survives this RDR.
- **A3 Established statechart standards reserve
  system-injected variable names (W3C SCXML §5.10:
  `_event`, `_sessionid`; authors must not bind names in
  the reserved `_` namespace).**
  - **Status**: Pending
  - **Method**: Prior Art
  - **Evidence**: W3C SCXML Recommendation §5.10 "System
    Variables" — to be fetched and quoted at Resolve; not
    present in local corpora (see research cache).
  - **If wrong**: prior-art alignment rests on XState and
    Inngest alone; the choice still stands on blast radius
    and the QOC matrix, but the Decision Rationale's
    prior-art row weakens from "standard + engines" to
    "engines".
- **A4 No fixture or authored table at HEAD declares or
  writes an owned/observed tag named `recognized`, so the
  reserved-key rule invalidates nothing that exists.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: to be confirmed by a whole-token sweep of
    `internal/resolve` fixtures and any committed table data
    for `recognized` used as an owned/observed key
    (`recognizedTagSensitiveTable` uses it only with
    recognized provenance).
  - **If wrong**: implementation must rename the colliding
    fixture/table tags before the lint lands; scope grows by
    a mechanical rename, not a design change.
- **A5 The kernel's D3 precedence (`owned` > `observed` >
  `recognized`) can remain unchanged as the deterministic
  backstop, because a table that passes the reserved-key
  validation can never present a colliding owned/observed
  key at resolve time.**
  - **Status**: Pending
  - **Method**: Derivation
  - **Evidence**: the reserved-key rule forbids declaring
    `recognized` with owned/observed provenance, and RDR
    0002 forbids matching/writing undeclared tags; therefore
    a conforming normalized table has no owned/observed tag
    keyed `recognized` — the collision branch in `assemble`
    is unreachable for conforming input. Derivation to be
    written out at Resolve against RDR 0002's exact
    declaration semantics.
  - **If wrong**: shadowing is reachable for conforming
    tables and the silent no-match returns; the kernel would
    then need its own reserved-key breach check, changing
    kernel surface.
- **A6 RDR 0009's producer-obligation seam does not reach
  this RDR's data channel, so the reserved-key input
  precondition needs an enforcement locus of its own.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: RDR 0009's Normative Contracts scope the
    obligation to "every constructor of `resolve.Row`
    values" and enforce it as a `Resolve` entry precondition
    via a kernel-exported predicate; the string `Input`
    appears nowhere in RDR 0009. So 0009 supplies the
    *pattern* (producer obligation, Go error path, no new
    refusal kind) but not a boundary this RDR can ride.
    Resolve settles which locus in Load-Bearing Decisions /
    Enforcement locus carries it, and confirms the choice
    against `internal/resolve/resolve.go::Resolve`'s doc
    comment reserving the error return.
  - **If wrong** (0009's seam does reach `Input` after all):
    the obligation rides it as originally drafted and the
    Enforcement-locus decision collapses to a cross-reference
    — strictly less surface than the default lean.
  - **Consequence if unresolved**: resolve-time data keyed
    `recognized` silently shadows the recognized outcome
    under D3 (premortem P-1/P-6), and the RDR's blast-radius
    claim rests on a seam that does not carry it.

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
5. **Shadowing becomes unreachable on both channels, not
   re-decided.** The declaration channel is closed by
   load/lint (A5); the data channel by the producer
   obligation (A6). RDR 0001 D3's precedence stays untouched
   as the kernel's deterministic backstop behind both, for
   non-conforming input only.

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

```normative
The assembled evaluation view MUST bind the freshly
recognized outcome under exactly the tag key `recognized`
(`internal/resolve/resolve.go::recognizedTagKey`). The key
is a reserved kernel keyword: table authors conform to it
and MUST NOT rebind it.
```

```normative
A tag declaration with provenance `recognized` MUST be named
`recognized`, and a flow MUST carry at most one such
declaration. A tag declaration with provenance `owned` or
`observed` MUST NOT be named `recognized`. Any violation is
a data-level validation failure in the `reserved tag key`
category, reported at table load/lint before any resolution
— never a kernel refusal. The reserved-key comparison is
byte-exact on the parsed TOML key value: case-sensitive, no
trimming, no folding (so `Recognized` is an ordinary,
unreserved name).
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
  match on the parsed TOML key value — no case folding, no
  trimming, no aliasing. Near-spellings (`Recognized`, a
  quoted or whitespace-bearing variant) are by definition
  ordinary unreserved names; whether lint additionally
  warns on case-variants of the reserved word is a Resolve
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
  A6 verifies the seam before this is fixed.

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Reserved key injected into the assembled view | Predecessor (RDR 0001 D2, implemented) | Available | Ratified from latitude to contract; no code change |
| Tag declarations with provenance | Predecessor (RDR 0002, Final unimplemented) | Deferred | Name constraint on `recognized`-provenance declarations |
| Load/lint name validation | Predecessor (RDR 0002 normalizer) | Deferred | One additive validation category: `reserved tag key` |
| Input-boundary producer enforcement | This RDR (pattern borrowed from RDR 0009) | Open | Locus unsettled — see Load-Bearing Decisions / Enforcement locus (A6) |

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
| Prior-art alignment | matches XState fixed `event` slot, Inngest system fields, RDR 0002 `<clear>` precedent | no surveyed engine lets authors rename the injected slot (negative result, research cache) | no engine surveyed withdraws event visibility from guards |
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
already making that decision:
`internal/resolve/resolve.go::recognizedTagKey` is the only
reserved key in the kernel (whole-package sweep for key
constants and "reserved" yields nothing else), and RDR
0002's `<clear>` sentinel is the nearest sibling
reserved-token decision — its mechanism (syntactically
un-declarable sigil) was considered and rejected under
Load-Bearing Decisions / Naming; its *pattern* (carrier
reserves a token, validation enforces it) is what this RDR
reuses. No parallel signal is invented.

Premortem: hardened (variant) — critic verdict PASS, no
switch forced. Every finding is folded into live text; the
P-numbers are cited inline where each cure lives.

Joint-check: clear (7 peers)

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
  category and one declaration constraint it did not
  originally spec (additive, but still scope on an
  unimplemented Final RDR).

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
declaration, or a reserved-name owned/observed declaration,
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
wrong-named recognized-provenance declarations and
reserved-name owned/observed declarations both report the
`reserved tag key` data-level category.

### Phase 3: Conformance pinning

The MVV fixture pair plus the premortem-derived test set:
the behavioral spelling test (real kernel, assembled-view
assertion — never literal-vs-literal), the
producer-obligation breach test (observed input keyed
`recognized` caught at whichever locus A6 settles), the
guard/matcher same-view test (A2), the duplicate
recognized-declaration fixture, the reserved-key
normalization fixtures (case/quoting/whitespace variants
stay unreserved), and the golden refusal-text check
(offending name, required name, rule).

## Validation

### Testing Strategy

[Required — never omit. Test scenarios and coverage goals — what to test and
what constitutes "done." For non-functional concerns
(performance, security): state measurement strategy,
not estimates.]

1. **Scenario**: [Description]
   **Expected**: [Result]

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
  demoted to A3).
