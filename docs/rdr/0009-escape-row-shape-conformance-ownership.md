# Recommendation 0009: Ownership of escape-row shape conformance

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

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
- **Profile**: foundational — one contract (which layer
  guarantees that an escape row bears no writes and no
  clears, and in what form), enforced across two modules:
  RDR 0001's kernel and RDR 0002's normalizer.
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
  is floored at FOUNDATIONAL regardless of the contract axis.
  Matrix: rdr/stages/README.md. Seed estimates from the design
  shape; Resolve overwrites from the verified count; Stage 8
  Gate locks it at Draft → Final. Never skip lenses off a
  Draft Profile until Resolve has run. -->
- **Priority**: Medium
- **Related Issues**: kata `tgzh` — "Escape-row shape
  conformance is unowned at the RDR 0001/0002 seam"
  (`src:roborev`, `severity:medium`,
  `area:internal-resolve`, `release:post-1.0`).
- **Predecessors**: 0001-resolution-kernel,
  0002-transition-table-as-reviewable-data
- **Overrides**: None superseded. This RDR *narrows* RDR
  0001's recorded deferral ("Kernel assumes a parsed,
  reviewable table shape") by adding one entry precondition
  on the error path RDR 0001 already reserves for programmer
  mistakes, and makes RDR 0002's escape-rule prohibition an
  explicitly owned producer obligation. It does **not**
  reopen RDR 0001's closed refusal taxonomy and does not
  change RDR 0002's authored grammar. It is the single
  normative home of the escape-row shape rule; both peers
  cite it rather than restate it.
- **Seam Lineage**: `area:internal-resolve` (escape-row shape
  conformance at the RDR 0001 ↔ 0002 boundary) — no prior
  accretion.

## Problem Statement

Someone authoring a transition table writes an escape rule
to say "when this failure class happens, leave by this
route." They expect an escape to be purely an exit — it
names where control goes, not what state gets written down.
They discover the problem only if something downstream
persists a tag they never authorized: the escape fires, the
route is right, and afterwards an owned tag has a value no
rule in their table ever set. Nothing in the disposition
tells them the write came along for the ride, because as far
as their table is concerned the write does not exist.

Internally, this RDR must fix **which layer guarantees the
invariant "an escape row bears no writes and no clears,"
and in what form that guarantee is carried**. RDR 0002
constrains the authored rule (an escape rule MUST NOT carry
a write block or clear list) and omits writes from the
escape row it renders, but ships no artifact that makes the
combination unrepresentable. RDR 0001's kernel accepts a
`Row` with both `Escape` and `Writes` populated and copies
whatever the row carries into the emitted `Plan`, having
explicitly deferred table shape to its peer ("Kernel assumes
a parsed, reviewable table shape"). The shared
`internal/resolve` `Row` type is the vocabulary both sides
speak, and it permits the illegal combination. Neither RDR
assigns the guarantee to anyone.

The decision is a single fork — **where trust is placed at
the 0001/0002 seam** — with three defensible branches: the
invariant is a **value-level obligation** discharged by RDR
0002's normalizer and lint, with RDR 0001's `Row` documenting
the precondition it trusts; it is a **shape-level
obligation**, with `Row` splitting into ordinary and escape
variants so the combination cannot be constructed; or it is
a **runtime check**, with the kernel rejecting malformed
rows at evaluation. The third branch splits in two: rejection
as a *modeled refusal* reaches past this seam into RDR 0001's
closed five-kind taxonomy and RDR 0005's exit-code mapping,
while rejection on the *reserved error path* does not.

## Context

### Background

Raised by a roborev finding against the RDR 0001
implementation and triaged as kata `tgzh`. The reviewer read
RDR 0002's escape-rule prohibition, then read RDR 0001's
kernel, and found the prohibition had no enforcer on either
side of the seam.

Current behavior: `internal/resolve` defines a single `Row`
type carrying `Escape`, `Writes`, and `NextTags` together,
and the plan constructor copies `Writes` onto the emitted
plan without consulting whether the selection came from an
escape edge. A row populating both fields therefore produces
a write-bearing escape disposition.

Reachability today is narrow: the RDR 0002 normalizer is not
yet implemented, so the combination is reachable only from
hand-constructed `resolve.Row` values — tests and any future
non-TOML table producer (A5). What makes it worth
adjudicating rather than dropping is the blast radius if it
did become reachable: an escape disposition describing
owned-tag persistence that no authored rule sanctioned,
handed to the RDR 0004 accessor layer — precisely the
"kernel guessing at persistence" that RDR 0001's ADV-2b
names.

The escape fixtures in
`internal/resolve/fixtures_test.go::escapeRow` populate
`Writes`/`NextTags` on escape rows. Conforming them is an
implementation step of this RDR (Phase 2), superseding the
"Noted, not filed" disposition in RDR 0001's `triage.md`;
whether the frozen suite still discriminates once they are
stripped is A3.

Triage also examined and rejected merging this with kata
`z53t` (RDR 0008): both sit at the 0001↔0002 seam, but this
decides the structural validity of a row *value* while 0008
decides the identity of a *name*. Disjoint answer spaces, no
ordering dependency. RDR 0007 (guard predicate totality) is
the adjacent 0001↔0003 seam decision; its precedent on the
refusal taxonomy is weighed in Decision Rationale.

### Technical Environment

Go CLI (`intrastate`). The seam spans:

- `internal/resolve` — the RDR 0001 resolution kernel.
  Owns the `Row` and `Plan` types, the escape-candidate
  path, and the plan constructor that copies a row's tag
  fields onto the emitted plan. Implemented (RDR 0001 is
  `Implemented`). `Resolve` returns `(Result, error)`;
  modeled refusals travel the `Result` value with a nil
  error, and the error return is documented as "reserved
  for programmer mistakes, not for modeled refusals".
- RDR 0002 (`0002-transition-table-as-reviewable-data`,
  `Final`, not yet implemented) — owns the authored TOML
  rule grammar, the escape-rule prohibition on write blocks
  and clear lists, the normalized row rendering, and the
  validation categories including "malformed escape
  declaration."
- RDR 0001's refusal taxonomy — exactly five behavioral
  refusal kinds, deliberately closed (design decision A5),
  with no "malformed row" member.
- RDR 0005 (`0005-skill-integration-cli-contract`, `Final`)
  — maps refusal kinds to `CLIError` codes; Go errors travel
  the existing envelope (`internal/cli/clierr::GroupInternal`,
  exit 2), so it is touched only if a new refusal kind is
  introduced.
- RDR 0004 (`0004-accessor-execution-safety-model`,
  `Final`) — the downstream consumer of a plan's writes, and
  the layer where an unsanctioned write would take effect.
- RDR 0006 (`0006-graph-lint-authority-and-guarantees`,
  `Final`) — the graph-level lint neighbor. Its blocking
  authority covers graph invariants over the normalized
  model and, by its own boundary, leaves parse/render
  fidelity and row-shape validation with RDR 0002 — so this
  RDR's load-time enforcement assignment does not overlap
  its authority; both RDRs reuse `internal/cli/clierr`'s
  implemented exit-code surface, whose normative home is
  RDR 0005.

## Research Findings

### Investigation

Prior art was read before enumerating approaches; the query
ledger and citations are cached at
`docs/rdr/0009-escape-row-shape-conformance-ownership/evidence/research/propose-prior-art.md`.
⚠ no prior-art coverage for the structural row-shape
conformance problem class in the arc corpora
(`StateMachineRes` ×2 queries, `StateMachineLit` ×1);
external claims (SCXML load-time document conformance,
"parse, don't validate" shape-level doctrine) are from the
model prior and are deliberately not load-bearing — demoted
to Resolve assumption A6.

The load-bearing findings are in-repo. RDR 0002's Normative
Contracts already state both halves of the authored-path
guarantee: "An escape rule MUST contain an `escape` list and
MUST NOT contain a write block or clear list," and
"Normalization MUST render an escape rule as a candidate row
with row kind `escape`, its normal predicate set, source
rule id, source locator, and modeled failure class list" —
a rendering that names no writes. Its validation categories
include "malformed escape declarations." So the value-level
branch has a normative home already; what is missing is the
ownership statement and any guarantee off the authored path.

The decisive discovery is that a taxonomy-neutral runtime
channel already exists. RDR 0001 states: "Modeled refusal is
a value-level resolver disposition, not a CLI error and not
the Go error path for parser bugs, IO failures, or
programmer mistakes," and
`internal/resolve/resolve.go::Resolve` documents "the error
return is reserved for programmer mistakes, not for modeled
refusals." A producer handing the kernel an illegal row
shape *is* a programmer mistake — the breach class RDR 0001
pre-allocated a channel for. Runtime rejection therefore
does **not** require opening the five-kind taxonomy: only
the modeled-refusal variant does, and that variant is
rejected (Alternative 2).

Sibling-path check (discriminator reuse): the escape/ordinary
discrimination already exists — `resolve.go::rescues` and
the candidate loop's `len(row.Escape) != 0` partition in
`resolve.go::Resolve`. The chosen approach reuses that
signal; no row-kind field or parallel discriminator is
introduced. Searched for an existing table-shape validation
path in the kernel: none exists (`Resolve` performs no
structural check on `Table.Rows` before evaluation).

### Key Discoveries

- **Documented** — RDR 0002 Normative Contracts: escape
  rules MUST NOT carry a write block or clear list; the
  normalized escape-row rendering carries no writes;
  "malformed escape declaration" is an existing load-time
  validation category. The authored path has an enforcer
  specified; the guarantee lacks an owner and an off-path
  backstop.
- **Documented** — RDR 0001 reserves the kernel's Go error
  return for programmer mistakes
  (`internal/resolve/resolve.go::Resolve` doc contract;
  RDR 0001 "Modeled refusal is a value-level resolver
  disposition…"). A taxonomy-neutral runtime channel for
  this breach class already exists.
- **Documented** — the CLI already surfaces Go errors
  through the existing envelope:
  `internal/cli/clierr::GroupInternal` maps to exit 2, and
  RDR 0005 verified that refusal classes map "without
  requiring new exit-code groups beyond the current
  `clierr.ExitCodeFor` taxonomy." No RDR 0005 change is
  implied by an error-path breach (A2).
- **Documented** — clears normalize into writes (RDR 0002:
  a rule-level `clear` entry "renders as a `<clear>` write"),
  so at the kernel boundary "no writes and no clears"
  reduces to an empty `Writes` slice on escape rows.
- **Documented** — `resolve.go::planOf` copies
  `row.Writes` onto the emitted `Plan` unconditionally, for
  escape selections too (`Escaped: true`); ADV-2b
  (`internal/resolve/adversarial_test.go`) names the
  resulting class "the kernel guessing."
- **Assumed** — load-time schema/document rejection as the
  standard enforcement point in peer state-machine systems
  (SCXML conformance): model-prior, to be cited at Resolve
  (A6); the choice does not rest on it.

### Critical Assumptions

- **A1 The kernel's Go error return is currently unused by
  `Resolve` for any disposition (it always returns nil
  error today), so adding a table-shape breach error cannot
  collide with an existing error semantic or caller
  branch.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: `internal/resolve/resolve.go::Resolve`
    (every return path today carries nil error); grep all
    callers of `resolve.Resolve` for error handling that
    assumes nil.
  - **If wrong**: An existing caller treats a non-nil error
    as a different failure class and misreports the breach;
    the error contract must be reconciled before the check
    lands.
- **A2 A Go error from the kernel surfaces through the
  existing CLI envelope (`internal/cli/clierr::GroupInternal`,
  exit 2) with no new refusal code, exit group, or RDR 0005
  contract change.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: `internal/cli/clierr::ExitCodeFor` and
    `GroupInternal`; RDR 0005 Key Discoveries ("without
    requiring new exit-code groups beyond the current
    `clierr.ExitCodeFor` taxonomy").
  - **If wrong**: Surfacing the breach requires reopening
    RDR 0005's envelope contract, and the error-path branch
    loses its no-peer-change advantage.
- **A3 Bringing escape fixtures into conformance (dropping
  `Writes` from `escapeRow`-built rows) changes no frozen
  assertion outcome: every ADV/MVV/boundary test still
  discriminates, so the entry precondition lands without
  weakening RDR 0001's frozen evidence.**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: Run the `internal/resolve` test suite
    with conformed fixtures at Resolve; triage's claim
    ("changes no assertion outcome" — RDR 0001
    `artifacts/triage.md`, "Noted, not filed") is the prior
    to confirm, not the proof.
  - **If wrong**: A frozen test actually depends on a
    write-bearing escape row; the precondition would flip
    that test to an error and this RDR must re-open an RDR
    0001 deviation instead of shipping as check-plus-
    fixtures.
- **A4 Guarding `Writes` guards the user's *outcome*, not
  just a field: nothing downstream persists an escape plan's
  `NextTags` into owned tag state, so an escape row cannot
  mutate owned state through the blessed exit-route field.
  If it can, the conformance predicate must also require an
  empty `NextTags` on escape rows.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: RDR 0002 Normative Contracts, escape-row
    rendering block ("row kind `escape`, its normal
    predicate set, source rule id, source locator, and
    modeled failure class list"); RDR 0001's `Plan.NextTags`
    vs `Plan.Writes` doc contracts (writes are what "the
    accessor layer applies"); RDR 0004's contract for what
    an accessor may persist. The miss this must rule out:
    the original complaint ("an owned tag has a value no
    rule set") reproduced via NextTags persistence while a
    Writes-only invariant reports the table conforming.
  - **If wrong**: The invariant guards the wrong field for
    the stated problem — the user symptom recurs through
    `NextTags` — and the predicate must widen before lock;
    the owner and channel chosen here are unchanged.
- **A5 No production code path constructs `resolve.Row`
  values at HEAD — the only constructors are test fixtures —
  so the entry precondition changes no shipped behavior and
  needs no migration.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: Grep constructors/callers of
    `resolve.Row` and `resolve.Resolve` outside
    `internal/resolve` tests.
  - **If wrong**: An existing production producer could
    begin erroring at `Resolve` entry; a migration step and
    an RDR 0005 surfacing review become prerequisites.
- **A6 Load-time document/schema rejection is the standard
  enforcement point for table-shape conformance in peer
  state-machine systems (SCXML conformance model).**
  - **Status**: Pending
  - **Method**: Prior Art
  - **Evidence**: To cite at Resolve from the W3C SCXML
    Recommendation's conformance section (corpus pass found
    no coverage; claim is from the model prior and is
    deliberately not load-bearing for the choice).
  - **If wrong**: The normalizer-side half loses its
    external precedent and stands on RDR 0002's already-
    Final prohibition alone — the choice survives, the
    citation does not.
- **A7 The breach error can name the offending row
  (`RuleID`, `SourceLocator`) without violating RDR 0005's
  rule that CLI codes are never inferred by inspecting
  error strings: the error is diagnostic prose on the
  internal-error path, not a stable code carrier.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: RDR 0001 ("RDR 0005 maps each to a CLI
    error code without inspecting error strings" — the rule
    binds refusal kinds, not internal errors);
    `internal/cli/clierr::CLIError` construction sites for
    how wrapped internal errors carry detail.
  - **If wrong**: The error must carry structure (a typed
    error), which sharpens the contract at Pre-Lock but
    does not move ownership.

## Proposed Solution

### Approach

**Escape-row shape conformance is a producer-side value
obligation, owned by this RDR as a single rule, enforced at
table load by RDR 0002's normalizer/lint on the authored
path, and backstopped by the kernel as an entry precondition
whose breach travels RDR 0001's reserved programmer-mistake
error path — never a modeled refusal, never a new refusal
kind, never a `Row` type split.**

Concretely: any component that constructs `resolve.Row`
values must not populate `Writes` on a row whose `Escape`
list is non-empty. On the
authored path, RDR 0002's normalizer rejects the rule at
load — its existing "malformed escape declaration"
validation class, diagnosable to the source rule. Off the
authored path — hand-constructed rows in tests, any future
non-TOML table producer — the kernel's `Resolve` checks the
table's rows at entry and returns a Go error naming the
offending row's `RuleID` and `SourceLocator`. The error is
the channel RDR 0001 explicitly reserves for "programmer
mistakes": a malformed table is a broken producer, not a
behavioral disposition of the input tuple, so the five-kind
refusal taxonomy stays closed (Alternative 2).

Ownership: **this RDR is the single normative home of the
escape-row shape rule.** RDR 0002's normalizer *implements*
the load-time half (its prohibition and validation category
already exist; this RDR binds them with conformance
fixtures); RDR 0001's kernel *implements* the backstop half
(one entry check on a channel it already reserved). Neither
peer's normative surface is reopened: RDR 0001's "assumes a
parsed, reviewable table shape" deferral is narrowed to
"assumes, and cheaply asserts, …" — recorded in Overrides.

### Technical Design

Data flow is unchanged for every conforming table:
`assemble` → candidate partition (`len(row.Escape) != 0`) →
`gate` → selection → `planOf`. The one addition is a
structural conformance predicate over `Table.Rows`, exported
so any producer can call it at construction time, and called
by `resolve.go::Resolve` at entry before `assemble` — the
same function at both call sites. On the first row with
non-empty `Escape` and non-empty `Writes`, `Resolve` returns
a nil `Result` and a non-nil error naming the row. `planOf`
stays unconditional — under the precondition it can no
longer copy writes onto an escaped plan, so no stripping
logic is added (stripping was rejected as silent divergence;
see Briefly Rejected). The escape discriminator stays
`len(Escape) != 0` / `rescues` — the existing sibling
signal, reused.

#### Normative Contracts

```normative
Escape-row shape conformance — an escape row carries no
owned-state mutation — is a PRODUCER obligation on every
constructor of resolve.Row values: a Row with a non-empty
Escape list MUST have an empty Writes slice. (Authored clears
normalize to `<clear>` writes per RDR 0002, so the Writes
predicate carries both "no writes" and "no clears" at the
kernel boundary.) Whether the conformance predicate also
requires an empty NextTags is OPEN pending A4, which settles
it against RDR 0002's escape-row rendering and RDR 0004's
accessor write contract; the obligation, its owner, and its
enforcement channel do not depend on that outcome.
```

```normative
On the authored-table path, RDR 0002's normalizer/lint is the
enforcing implementation: an authored escape rule carrying a
write block or clear list MUST be rejected at table load
under RDR 0002's "malformed escape declaration" validation
class, diagnosable to the source rule. Normalized escape rows
MUST render write-free. This RDR binds that duty with a
conformance fixture set; RDR 0002's grammar is unchanged.
```

```normative
The kernel enforces the same obligation as an entry
precondition of Resolve: a table containing a row that
breaches the conformance predicate above MUST cause Resolve
to return a non-nil Go error identifying the offending row by
RuleID and SourceLocator, with no Result disposition. The
precondition is evaluated over the WHOLE table before any
evaluation step: a breach surfaces even when no resolution
path reaches the offending row, and it precedes every
modeled disposition (a malformed table is malformed as a
value, independent of the input tuple). The breach travels
the error path RDR 0001 reserves for programmer mistakes:
it MUST NOT be a modeled refusal, MUST NOT introduce a
refusal kind, and MUST NOT alter any disposition of a
conforming table.
```

```normative
The conformance predicate MUST be exported by the kernel
package as a construction-time check callable by any table
producer (form sharpened at Pre-Lock; e.g. a
ValidateTable-style function), and Resolve's entry
precondition MUST be that same function — one predicate,
two call sites, so a non-TOML producer can fail at build or
construction time rather than at first production Resolve,
and the two enforcement points cannot drift.
```

```normative
The shared resolve.Row type keeps its single shape: escape
identity remains discriminated solely by a non-empty Escape
list. No ordinary/escape type split and no row-kind field is
introduced at the kernel boundary.
```

#### Load-Bearing Decisions

- **Selection / predicate** — the breach predicate over a
  single row is `len(row.Escape) != 0 && <mutation test>`,
  where the mutation test is `len(row.Writes) != 0` today and
  becomes `(len(row.Writes) != 0 || len(row.NextTags) != 0)`
  if A4 finds that `NextTags` on an escape row can reach
  owned state. **A4 is the only open input**; these are
  pinned regardless of its outcome:
  - *Scope* — evaluated over every row of the supplied table
    at `Resolve` entry, not only rows the resolution touches:
    a malformed table is malformed as a value, and checking
    only reached rows would make the breach appear and
    disappear with the input tuple.
  - *Nil-vs-empty* — length-based on purpose: a non-nil but
    empty slice is conforming. The predicate tests emptiness,
    never nil-ness.
  - *Precedence* — the breach error precedes every modeled
    disposition, `unmodeled_outcome` included.
  - *Locus* — one exported function, called at both sites, so
    widening under A4 is a single-site edit that cannot
    desynchronize the two enforcement points.
- **Naming** — the breach is "escape-row shape breach,"
  carried as a Go error. Its form (opaque vs sentinel/typed)
  is open, sharpened at Pre-Lock and constrained by A7.

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Reserved programmer-mistake error path on `Resolve` | Predecessor (RDR 0001, implemented) | Available | The backstop's channel exists; no taxonomy or signature change |
| Load-time rejection of malformed escape declarations | Predecessor (RDR 0002, Final, unimplemented) | Deferred | Conformance fixtures from this RDR bind the normalizer build |
| CLI surfacing of internal errors (exit 2) | Predecessor (RDR 0005 + `internal/cli/clierr`, implemented) | Available | Breach reaches operators without an RDR 0005 change (A2) |
| Escape/ordinary row discrimination | Predecessor (RDR 0001, implemented) | Available | Reused (`rescues`, `len(Escape) != 0`); no new discriminator |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Escape discrimination | `internal/resolve/resolve.go::rescues` + candidate partition in `::Resolve` | None | Reuse | The precondition keys on the same field; no row-kind added |
| Plan construction | `internal/resolve/resolve.go::planOf` | Copies `Writes` unconditionally | Reuse | Unchanged — the precondition makes stripping unnecessary |
| Error envelope | `internal/cli/clierr::GroupInternal` / `ExitCodeFor` | Exit 2 is shared with other internal errors | Reuse | Breach is diagnosable from the error text (A7), not a new code |
| Escape fixtures | `internal/resolve/fixtures_test.go::escapeRow` | Populates `Writes`/`NextTags` on escape rows | Extend (conform) | Fixtures must satisfy the precondition; supersedes the "Noted, not filed" cleanup |

### Decision Rationale

The undecided contract (accretion-gate discipline, applied
without recorded accretion): **"which layer guarantees that
an escape row bears no writes and no clears, and in what
form is that guarantee carried?"** Approaches were
enumerated as answers to it.

Scored matrix (criteria weighted by the Problem Statement's
user outcome — no unsanctioned write ever rides an escape,
and every failure is diagnosable to the responsible party):

| Criterion | A: Normalizer-only (trust + doc'd precondition) | B: Shape split (unrepresentable) | C: Runtime modeled refusal (new kind) | D: Obligation + kernel error-path precondition (chosen) |
| --- | --- | --- | --- | --- |
| Correctness fit (every producer path closed?) | ✗ — guards only the authored path; the sole *reachable* path today (hand-built rows, future non-TOML producers) stays open | ✓ — unconstructable | ✓ — caught at evaluation | ✓ — load-time on the authored path, entry check on every other |
| Taxonomy / precedent alignment | ✓ neutral | ✓ neutral | ✗ — opens RDR 0001's closed five-kind set (A5); collides with RDR 0007's recorded rejection of a sixth kind | ✓ — uses the channel RDR 0001 pre-reserved for exactly this breach class |
| Reversibility | ✓ — docs + fixtures | ✗ — permanent public type surface; reopens RDR 0002's single-shape row rendering and RDR 0001's implemented types | ~ — refusal kinds are forever (RDR 0005 mapping) | ✓ — one check + one error, removable without breaking any modeled contract |
| Blast radius | Minimal now; unbounded at the first unsanctioned write (RDR 0004 applies it) | Kernel `Row`/`Table`/frozen ADV fixtures + 0002 rendering reopen | 0001 taxonomy + 0005 exit-code mapping | Small: one entry scan + fixture conformance; zero disposition change for conforming tables |
| Cost | Low | High | Medium | Low — discriminator and channel both already shipped |

The deciding rows are correctness fit and taxonomy
alignment. D is the only approach that closes every producer
path without opening the taxonomy or any peer's locked
surface, and the sibling-path check shows both of its
ingredients already exist (`rescues`/`len(Escape)!=0` as the
discriminator, the reserved error return as the channel) —
it reuses existing signals where B and C would mint new
ones. A is rejected because it guards the path that cannot
yet misbehave and leaves open the only path that can. B is
rejected on blast radius and reversibility: Go's lack of sum
types makes the split a permanent two-type public surface
that reopens two locked RDRs to prevent a state one
predicate detects. C is rejected on taxonomy alignment
(Alternative 2): a malformed table is a broken producer, not
a behavioral disposition of the input tuple, and RDR 0001
both closed the behavioral taxonomy (A5) and pre-allocated
the programmer-mistake channel.

Premortem: hardened — critic verdict PASS, no switch forced.
The ownership structure survived every negation; the
findings landed at unwritten seams, not at the choice, and
are discharged in the live text above (the exported-predicate
clause, the whole-table and precedence pins, A4's outcome
framing, the shared Phase 3 fixtures, and the recorded
exit-2 diagnostic decision in Failure Modes). Ledger:
`docs/rdr/0009-escape-row-shape-conformance-ownership/evidence/propose-premortem/critic.md`.

Joint-check: clear (7 peers)

## Alternatives Considered

### Alternative 1: Shape-level unrepresentability (`Row` split)

**Description**: Split the shared row vocabulary into
ordinary and escape variants (two types, or a kind-tagged
sum emulation), so a value carrying both `Escape` and
`Writes` cannot be constructed. "Parse, don't validate" as
type design.

**Pros**:

- The invariant holds by construction; no producer can
  breach it and no check can be forgotten.
- Self-documenting vocabulary: the type says what an escape
  row may carry.

**Cons**:

- Go has no sum types: the split becomes two public types
  plus an interface or a two-slice `Table`, a permanent
  surface reworking the implemented kernel (`Resolve`'s
  candidate loop, `gate`, `escapeOrRefuse`, `planOf`) and
  its frozen ADV fixtures.
- Reopens RDR 0002's rendering contract, which normatively
  renders *one* candidate-row shape with a row-kind marker,
  and RDR 0001's implemented type vocabulary — two locked
  RDRs reshaped to prevent a state one predicate detects.
- Poor reversibility: type surface, once public, is forever;
  a predicate can be moved or removed.

**Reason for rejection**: highest blast radius and worst
reversibility for a guarantee D reaches with one entry check
on an existing channel.

### Alternative 2: Runtime rejection as a modeled refusal (new kind)

**Description**: The kernel treats a write-bearing escape
row as a value-level refusal — a sixth kind (e.g.
`malformed_row`) alongside the behavioral five — mapped by
RDR 0005 to a CLI code.

**Pros**:

- Structured, machine-branchable surfacing of the breach to
  CLI callers, symmetrical with other refusals.
- Enforcement point identical to D's (kernel, at
  evaluation).

**Cons**:

- Opens RDR 0001's deliberately closed five-kind taxonomy
  (design decision A5) and RDR 0005's exit-code mapping —
  two Final/Implemented contracts — for a condition that is
  not a behavioral disposition of the input tuple but a
  broken producer.
- Category error the taxonomy was closed against: refusals
  describe the resolution of a well-formed question;
  a malformed table means the question itself was
  ill-posed. RDR 0001 already routes that class to the
  error path ("parser bugs, IO failures, or programmer
  mistakes").
- RDR 0007 weighed and rejected the same move
  (`guard_input_missing`) at the adjacent seam; two
  adjacent seams independently minting kinds is exactly how
  a closed taxonomy erodes.

**Reason for rejection**: fails the taxonomy-alignment
criterion; the reserved error path delivers the same
enforcement point without opening any peer contract.

### Briefly Rejected

- **Normalizer-only trust (matrix column A)**: enforces
  where the combination already cannot arise once RDR 0002
  is built, and leaves the only reachable construction path
  (hand-built rows, future non-TOML producers) guarded by
  documentation alone.
- **Kernel sanitization (strip `Writes` when the selection
  escapes, in `planOf`)**: silently diverges behavior from
  the reviewed table — the row says writes, the plan says
  none — masking the producer bug and violating RDR 0002's
  reviewable-data doctrine that behavior be derivable from
  the reviewed artifact.
- **Test-side lint only (assert fixtures conform)**: binds
  this repo's tests, not future producers; the seam stays
  unowned.

## Trade-offs

### Consequences

- Positive: no write-bearing escape plan can be emitted by
  any producer path — the ADV-2b "kernel guessing at
  persistence" class is closed before RDR 0004 can apply an
  unsanctioned write.
- Positive: the five-kind taxonomy, RDR 0005's mapping, RDR
  0002's grammar, and the kernel's type vocabulary all stand
  unchanged; the check reuses a shipped discriminator and a
  reserved channel.
- Positive: authored-path failures stay author-diagnosable
  (load-time, named source rule); programmatic-path failures
  name the offending row.
- Negative: `Resolve` takes on a narrow validation duty RDR
  0001 had deferred wholesale to its peer — the deferral is
  narrowed, not honored verbatim (recorded in Overrides).
- Negative: every `Resolve` call pays an O(rows) conformance
  scan; acceptable at table scale, and hoisting is an
  implementation detail, but the cost is per-call because
  the kernel is stateless.
- Negative: a breach off the authored path surfaces as a
  generic exit-2 internal error, not a structured refusal —
  operators diagnose from error text (A7), which is the
  price of keeping the taxonomy closed.
- Negative (behavior change, deliberate): a hand-built table
  carrying a *dormant* malformed row — one no resolution
  path reaches — previously resolved fine and now errors on
  every `Resolve` (whole-table precondition). The MVV pins
  the new outcome. Whether any deployed table can regress
  turns on A5 (no production `resolve.Row` constructor at
  HEAD); if A5 fails, a migration step becomes a
  prerequisite.

### Risks and Mitigations

- **Risk**: RDR 0002's implementer treats the kernel
  precondition as *the* enforcer and skips or weakens
  load-time rejection, so table authors get an exit-2
  internal error instead of a load diagnostic naming their
  rule.
  **Mitigation**: The conformance fixture set (Phase 3) is
  normative — authored escape rule with a write block →
  load-time rejection — and Prerequisites bind RDR 0002's
  implement stage to it; the kernel error is documented as a
  backstop for programmatic producers, not the authoring
  UX.
- **Risk**: The error path accretes semantics — a future
  change routes some modeled condition to the error return
  because this RDR normalized its use.
  **Mitigation**: The normative clause states the breach
  MUST NOT be a modeled refusal and the error path remains
  reserved for producer/programmer mistakes; the contract is
  one predicate, not an open validation framework.
- **Risk**: A frozen ADV test secretly depends on a
  write-bearing escape fixture, so the precondition breaks
  RDR 0001's evidence base.
  **Mitigation**: A3 is a Spike — run the suite with
  conformed fixtures before lock; if it fails, the finding
  routes back as an RDR 0001 deviation question rather than
  shipping silently.
- **Risk**: The predicate under-covers (`NextTags` rides an
  escape unexamined).
  **Mitigation**: A4 settles NextTags against RDR 0002's
  rendering contract at Resolve; widening the predicate is
  mechanical and moves no ownership.

### Failure Modes

- **Visible break (programmatic producer)**: `Resolve`
  returns an error naming the offending row's `RuleID` and
  `SourceLocator`; at the CLI it surfaces as an internal
  error (exit 2). Diagnosis: the row identity in the error
  text points at the producing code or fixture. Recovery:
  fix the producer — the table value is broken; there is
  nothing to retry.
- **Visible break (table author)**: the RDR 0002 normalizer
  rejects the rule at load under "malformed escape
  declaration," naming the source rule. Recovery: remove the
  write block/clear list from the escape rule, or make it an
  ordinary rule.
- **Silent failure guarded against**: an escape plan
  carrying writes no rule sanctioned — previously
  constructible and downstream-applied by RDR 0004 with no
  symptom until owned state diverged. Under the
  precondition it cannot be emitted; the conformance
  fixtures are the tripwire against regression.
- **Diagnostic gap (recorded decision)**: exit 2 does not
  distinguish this breach from other internal errors by
  code; only the error text — which names the row — does
  (A7). An operator runbook entry therefore reads "exit 2
  from a resolve: read the error text; a row-identified
  escape-shape breach means a broken table producer, not a
  tool defect". If operational experience
  needs a machine-branchable code, that is an RDR 0005
  envelope question — explicitly out of scope here, and the
  recorded cost of keeping the taxonomy closed.
- **Error-channel creep guarded against**: this RDR's
  normative clause states what qualifies for the error path
  (a producer/programmer contract breach, never a modeled
  condition), so the next borrower argues against text, not
  precedent.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified
- [ ] A3 spike green (conformed fixtures, full
      `internal/resolve` suite) before the precondition
      lands
- [ ] RDR 0002's implement stage binds to this RDR's
      conformance fixtures (Phase 3) once both are Final

### Minimum Viable Validation

A hand-constructed table containing an escape row with
populated `Writes`: `Resolve` returns a non-nil error naming
that row's `RuleID` and `SourceLocator`, with no `Plan` and
no `Refusal` — while the same table with the writes removed
resolves, escapes, and refuses exactly as the frozen ADV
suite proves today (zero disposition change for conforming
tables). Three scenarios are in scope:
*breach-yields-error-not-plan* (the pre-fix probe — a
write-bearing escape row selected via escape → plan with
`Escaped:true` carrying writes — now yields the error, never
that plan); *dormant-row-still-errors* (a malformed row no
resolution path reaches → error on every `Resolve`, pinning
the whole-table decision); and *empty-not-nil-conforms* (a
non-nil empty `Writes` slice on an escape row → no error,
pinning the length-based predicate). The exported validator
and the entry check are exercised as the same predicate.

### Phase 1: Exported predicate + kernel entry precondition

Export the conformance predicate as a construction-time
check, call it at `Resolve` entry (breach → Go error naming
the row; conforming tables unchanged), state the producer
obligation on the `Row` doc contract, and record the
kernel's validation surface — exactly which shape property
it checks vs still assumes — so partial validation cannot be
read as general kernel ownership.

### Phase 2: Fixture conformance with discrimination check

Bring `escapeRow` fixtures into conformance: drop `Writes`,
and drop `NextTags` only if A4 widened the predicate. A3's
spike proves the frozen suite still discriminates, and one
mutation-style step (revert the entry check; at least one
test must fail) proves the stripped fixtures did not leave
the check untested.

### Phase 3: Shared normalizer conformance fixtures

Encode the authored-path half as a named, shared fixture set
binding the future RDR 0002 build: an authored escape rule
with a write block or clear list fails load under "malformed
escape declaration" naming the source rule; normalized
escape rows render write-free; and one canonical
authored-clear case pins the `<clear>`-sentinel-write
representation identically for the kernel suite and the
normalizer suite, so the two enforcers of one invariant
cannot drift.

## Validation

### Testing Strategy

Done means: every producer path is closed, no conforming
table changes disposition, and the check itself is proven
non-vacuous. The MVV's three scenarios are the core; the
suite below is the coverage goal.

1. **Scenario**: Escape row with populated `Writes`, selected
   via the escape path.
   **Expected**: `Resolve` returns a non-nil error naming the
   row's `RuleID` and `SourceLocator`; no `Plan`, no
   `Refusal`.
2. **Scenario**: Escape row with populated `Writes` that no
   resolution path reaches (dormant).
   **Expected**: Same error — the precondition is whole-table,
   not on-path.
3. **Scenario**: Escape row with a non-nil but empty `Writes`
   slice.
   **Expected**: No error; resolves normally. Pins the
   length-based predicate against a nil-ness reading.
4. **Scenario**: Breach coexists with a table state that would
   otherwise yield `unmodeled_outcome`.
   **Expected**: The breach error wins — precedence over every
   modeled disposition.
5. **Scenario**: The full frozen `internal/resolve` suite
   (ADV/MVV/boundary) run against conformed escape fixtures.
   **Expected**: Every assertion still passes and still
   discriminates — no test silently becomes vacuous (A3).
6. **Scenario**: Mutation check — revert the entry check and
   re-run.
   **Expected**: At least one test fails, proving the
   conformed fixtures did not leave the check untested.
7. **Scenario**: The exported predicate called directly by a
   producer at construction time, on the same tables as
   scenarios 1–3.
   **Expected**: Verdicts identical to `Resolve`'s entry check
   — one predicate, two call sites, no drift.
8. **Scenario** (binds RDR 0002's build): authored escape rule
   carrying a write block, and one carrying a clear list.
   **Expected**: Both rejected at load under "malformed escape
   declaration," naming the source rule; normalized escape
   rows render write-free.

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
Profile, correct the field and do not lock until the
missing lenses have run. Also confirm form: value +
one clause naming the contract(s); strip any
matrix/provenance prose left from the template or Seed
(it belongs in the template comment, not the instance).]

## References

- RDR 0001 — Resolution kernel (closed refusal taxonomy /
  A5; "Modeled refusal is a value-level resolver
  disposition, not a CLI error and not the Go error path
  for parser bugs, IO failures, or programmer mistakes";
  Existing Infrastructure Audit deferral "Kernel assumes a
  parsed, reviewable table shape"; `artifacts/triage.md`
  "Noted, not filed" fixture entry)
- RDR 0002 — Transition table as reviewable data
  (escape-rule prohibition; escape-row rendering; "malformed
  escape declarations" validation class; clear-as-`<clear>`-
  write normalization)
- RDR 0005 — Skill integration CLI contract (refusal→code
  mapping; "without requiring new exit-code groups")
- RDR 0006 — Graph lint authority and guarantees (boundary:
  graph-level lint authority; row-shape/parse validation
  stays with RDR 0002)
- RDR 0007 — Guard predicate totality (precedent: rejected
  minting `guard_input_missing` to keep the taxonomy closed)
- `internal/resolve/resolve.go` (`Row`, `Plan`, `Resolve`,
  `planOf`, `rescues`, `escapeOrRefuse`)
- `internal/resolve/adversarial_test.go` (ADV-2b — "kernel
  guessing" comment)
- `internal/resolve/fixtures_test.go::escapeRow`
- `internal/cli/clierr` (`GroupInternal`, `ExitCodeFor`)
- kata `tgzh` (originating finding)
- Prior-art query ledger:
  `docs/rdr/0009-escape-row-shape-conformance-ownership/evidence/research/propose-prior-art.md`
