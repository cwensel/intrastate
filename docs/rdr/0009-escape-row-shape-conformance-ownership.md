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
- **Profile**: foundational — one contract (an escape row
  bears no owned-state mutation: no writes, no clears),
  enforced across two modules — RDR 0001's kernel and RDR
  0002's normalizer — and consumed by a third (RDR 0005's
  CLI surfacing).
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

Every escape row in the frozen suite carries `Writes`, from
two sources. The builder
`internal/resolve/fixtures_test.go::escapeRow` sets
`Writes: []resolve.Tag{{Key: "status", Value: "Blocked"}}`,
which all sixteen of its call sites inherit; on top of that,
two call sites *override* the value with `escape.Writes =
[]resolve.Tag{{Key: "status", Value: "Escaped"}}` —
`internal/resolve/adversarial_test.go:229` (ADV-2b) and
`internal/resolve/fixup_test.go:103` (Fixup-1e). Conforming
all three sites is an implementation step of this RDR
(Phase 2), superseding the "Noted, not filed" disposition in
RDR 0001's `triage.md`; A3 confirmed the frozen suite still
discriminates once they are stripped.

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
  model, and it states its own boundary as "Parse/render
  fidelity remains owned by RDR 0002." Row *shape* validation
  is not named on either side of that line, so this RDR
  assigns it rather than inheriting an assignment: the
  authored-path half goes to RDR 0002 as the layer that
  already forbids the construct and owns rendering, which is
  the side of 0006's boundary it falls on. No overlap with
  0006's graph-invariant authority results; both RDRs reuse
  `internal/cli/clierr`'s
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
"parse, don't validate" shape-level doctrine) were from the
model prior and deliberately not load-bearing — demoted to
Resolve assumption A6, where the SCXML claim was researched
and **falsified** (ledger:
`evidence/research/resolve-prior-art.md`).

Stage 4 also ran a DX pass over the *form* of the check —
error shape, aggregate-vs-fail-fast reporting, validation
locus, and naming — against the Go standard library and this
repo's own idiom; findings are cited on A7 and in
Load-Bearing Decisions, ledger at
`evidence/research/resolve-dx-prior-art.md`.

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
- **Researched and falsified** (Resolve) — load-time
  schema/document rejection is **not** an established
  standard enforcement point. W3C SCXML IRP test 313 states
  a processor "MAY reject documents containing syntactically
  ill-formed expressions at document load time, or it MAY
  wait and place error.execution in the internal event queue
  at runtime" — the enforcement point is deliberately left to
  the implementation. A6 restated accordingly; the choice
  never rested on it and is unchanged.

### Critical Assumptions

- **A1 The kernel's Go error return is currently unused by
  `Resolve` for any disposition (it always returns nil
  error today), so adding a table-shape breach error cannot
  collide with an existing error semantic or caller
  branch.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/resolve/resolve.go::Resolve` —
    all five return paths carry a literal `nil` error
    (`KindUnmodeledOutcome`, the `blocked` refusal, the
    exact-one plan, and both `escapeOrRefuse` arms); no
    helper (`gate`, `escapeOrRefuse`, `planOf`, `refuse`)
    returns an error type, so nothing can propagate into
    the slot. Callers: the only importers of
    `internal/resolve` are its own `*_test.go` files plus
    `docs/rdr/0007-guard-predicate-totality/evidence/spikes/aggregation-probe_test.go`
    — no production file imports the package. Existing
    tests do more than assume nil, they *assert* it
    (`internal/resolve/mvv_test.go` "modeled refusal %q
    traveled the Go error path";
    `internal/resolve/resolve_test.go` "Resolve returned a
    Go error for a modeled disposition").
  - **If wrong**: An existing caller treats a non-nil error
    as a different failure class and misreports the breach;
    the error contract must be reconciled before the check
    lands.
- **A2 A Go error from the kernel surfaces through the
  existing CLI envelope (`internal/cli/clierr::GroupInternal`,
  exit 2) with no new refusal code, exit group, or RDR 0005
  contract change.**
  - **Status**: Verified (with a named wiring obligation)
  - **Method**: Source Search
  - **Evidence**: `internal/cli/clierr::ExitCodeFor` —
    `case GroupUserEnv, GroupInternal: return 2`.
    `internal/cli/config/config.go::Load` is the shipped
    precedent for wrapping a Go error into that group
    (`Code: "config-read-error"`, `Detail: err.Error()`,
    `Group: clierr.GroupInternal`, `Cause: err`), and
    `clierr::CLIError.Unwrap` keeps it traversable. RDR
    0005 Key Discoveries ("without requiring new exit-code
    groups beyond the current `clierr.ExitCodeFor`
    taxonomy") and its Technical Design ("Internal parse or
    invariant failures use the existing internal/user-env
    mapping") pre-authorize the routing.
    **Obligation, not a gap:** no `flow` verb exists at
    HEAD (`internal/cli/root.go::NewRootCmd` registers only
    `newVersionCmd()`), so this path is specified, not
    shipped. An *unwrapped* Go error escaping a future
    verb's `RunE` exits **1**, not 2 —
    `clierr::ExitCodeFor` defaults a non-`CLIError` to 1 and
    `root.go::cobraErrorToCLIError` stamps
    `GroupUserEnv`/`command-error`. Exit 2 therefore holds
    only if the verb wraps the kernel error into
    `GroupInternal`; the Normative Contracts below bind
    that wrap.
    **Second obligation (found at Pre-Lock):** the exit code
    routes without an RDR 0005 change, but the *row identity*
    does not travel the wire for free —
    `clierr::CLIError.Cause` is `json:"-"`, so the Go error
    chain carrying the `RowRef`s is not serialized. Reaching
    an operator requires rendering the identities into a
    serialized field (`Detail`, or a new `omitempty` field
    under the type's documented "Extend with new optional
    fields as needed"). That is additive to the envelope and
    leaves the refusal-to-exit-code mapping untouched, so
    A2's "no RDR 0005 contract change" verdict stands.
  - **If wrong**: Surfacing the breach requires reopening
    RDR 0005's envelope contract, and the error-path branch
    loses its no-peer-change advantage.
- **A3 Bringing escape fixtures into conformance (dropping
  `Writes` from write-bearing escape rows) changes no frozen
  assertion outcome: every ADV/MVV/boundary test still
  discriminates, so the entry precondition lands without
  weakening RDR 0001's frozen evidence.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**:
    `docs/rdr/0009-escape-row-shape-conformance-ownership/evidence/spikes/a3-fixture-conformance.md`
    — `go test ./internal/resolve/... -v -count=1` at three
    fixture states: baseline (`a3-baseline.out`), `Writes`
    stripped (`a3-conformed.out`), and `Writes`+`NextTags`
    stripped (`a3-conformed-plus-nexttags.out`). All three:
    154 `--- PASS`, 0 `--- FAIL`, **identical sorted outcome
    sets**. Non-vacuity established two ways: (a) every
    `Plan.Writes` read on an escape-derived plan
    (`adversarial_test.go` ADV-2b, `fixup_test.go` Fixup-1d
    / Fixup-1e) sits in a `t.Fatalf` *format argument* on
    the failure path, never in a condition — the asserted
    properties are `Refusal.Kind` / `Refusal.Guard`; every
    substantive `Writes`-equality assertion (REQ-18, REQ-19,
    REQ-29, MVV leg 1) runs on `singleMatchTable` /
    `legalInput`, untouched; (b) mutation check — bypassing
    the uniform viability gate in `escapeOrRefuse`
    (`viable := escapes`) against the conformed tree still
    fails 11 assertions across ADV-2, ADV-2b, Fixup-1d,
    Fixup-1e, and `TestFixupGateIsUniformAcrossOrdinaryAndEscapeCandidates`
    (`a3-conformed-mutant-gate-bypass.out`). Spike edits
    fully reverted; suite re-verified green.
  - **Scope (measured by this spike)**: conformance must
    edit **three** sites. The builder
    `internal/resolve/fixtures_test.go::escapeRow` sets
    `Writes: []resolve.Tag{{Key: "status", Value: "Blocked"}}`
    (line 257), inherited by all sixteen of its call sites;
    two of those sites additionally override it with
    `escape.Writes = []resolve.Tag{{Key: "status", Value: "Escaped"}}`
    — `internal/resolve/adversarial_test.go:229` (ADV-2b)
    and `internal/resolve/fixup_test.go:103` (Fixup-1e
    "missing owned state"). Stripping only the builder
    leaves those two breaching; stripping only the overrides
    leaves fourteen.
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
  - **Status**: Verified — settled **closed**: the
    predicate does **not** widen; it stays `Writes`-only.
  - **Method**: Peer RDR
  - **Evidence**: Three links, each closed by normative
    text rather than inference. (1) *Only a write accessor
    persists* — RDR 0004 Normative Contracts scope the
    capability set to read/gate/write, with the read
    accessor "MUST NOT mutate authoritative artifacts" and
    the gate accessor returning only allow/deny/
    indeterminate. (2) *That door is scoped to `Writes`, not
    to the plan* — RDR 0004: "A write accessor MUST apply
    **only** planned owned-tag writes produced by a
    successful transition. It MUST NOT write observed or
    recognized tags." An escape plan with empty `Writes`
    hands the accessor an empty work set whatever
    `NextTags` carries. (3) *`NextTags` is not addressed to
    that layer at all* — the string `NextTags` (and "next
    state") appears **nowhere** in RDR 0004; positively,
    `internal/resolve/resolve.go::Plan` separates them in
    the type's own doc contract ("NextTags is the next
    state." vs "Writes are the owned-tag writes **for the
    accessor layer**"), as does RDR 0001's Technical Design
    ("the owned-tag writes **that the accessor layer
    applies**"). RDR 0005's `flow set-state` is likewise
    driven by "planned owned-tag `--write name=value`
    mutations." Corroborating: the A3 spike's
    `NextTags`-stripped variant changed no outcome
    (`a3-conformed-plus-nexttags.out`), so nothing in the
    frozen suite reads `Plan.NextTags` on an escape-derived
    plan either.
  - **Carried forward (peer-contract dependency)**:
    `resolve.go::planOf` copies `NextTags` and `Writes`
    symmetrically with no escape-awareness, so this safety
    rests entirely on RDR 0004's write-accessor scoping,
    not on any kernel-local property. A future RDR 0004
    implementer widening "apply the plan" beyond planned
    writes would breach this RDR's invariant.
  - **Consistency note**: RDR 0002's escape-row rendering
    ("row kind `escape`, its normal predicate set, source
    rule id, source locator, and modeled failure class
    list") names **neither** writes **nor** next-tags, so a
    normalized escape row is authored-path-stricter than
    the kernel predicate requires. Stricter upstream,
    consistent — not a conflict.
  - **If wrong**: The invariant guards the wrong field for
    the stated problem — the user symptom recurs through
    `NextTags` — and the predicate must widen before lock;
    the owner and channel chosen here are unchanged.
- **A5 No production code path constructs `resolve.Row`
  values at HEAD — the only constructors are test fixtures —
  so the entry precondition changes no shipped behavior and
  needs no migration.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: Repo-wide, every reference to
    `resolve.Row` / `resolve.Table` / `resolve.Input` /
    `resolve.Resolve` resolves inside
    `internal/resolve/{fixtures,resolve,adversarial,mvv,fixup}_test.go`.
    No production file imports the package at all:
    `cmd/intrastate/main.go` is `func main() { cli.Execute() }`,
    and `internal/cli/root.go::NewRootCmd` registers only
    `newVersionCmd()` — there is no `flow` verb at HEAD.
    Constructors are exclusively the fixture builders
    (`singleMatchTable`, `noMatchTable`, `ambiguousTable`,
    `missingOwnedTable`, `unevaluableGuardTable`,
    `twoRowsOneGuardFalseTable`, `observedSensitiveTable`,
    `recognizedTagSensitiveTable`, `escapeRow`) plus inline
    literals in `adversarial_test.go` and `fixup_test.go`.
  - **If wrong**: An existing production producer could
    begin erroring at `Resolve` entry; a migration step and
    an RDR 0005 surfacing review become prerequisites.
- **A6 Peer state-machine systems do NOT fix a standard
  enforcement point for table-shape conformance: SCXML
  deliberately leaves the choice between load-time
  rejection and a runtime error channel to the
  implementation.** (Restated at Resolve. The Propose-stage
  claim — that load-time rejection is *the standard* — was
  researched and **falsified**; it was citation-only and
  never load-bearing, so the approach is untouched.)
  - **Status**: Verified (contrary to the original claim)
  - **Method**: Prior Art
  - **Evidence**: W3C SCXML IRP conformance test 313, as
    quoted in uSCXML's conformance matrix
    (`state-machines/study/uscxml/test/w3c/TESTS.md`):
    "The SCXML processor MAY reject documents containing
    syntactically ill-formed expressions at document load
    time, **or** it MAY wait and place error.execution in
    the internal event queue at runtime when the
    expressions are evaluated." The one load-time MUST in
    that matrix (test 301) covers an undownloadable script
    resource, not document shape. Ledger + rejected
    branches:
    `docs/rdr/0009-escape-row-shape-conformance-ownership/evidence/research/resolve-prior-art.md`.
  - **Consequence**: With no external convention to defer
    to, this RDR's two-point enforcement rests on in-repo
    authority — RDR 0002's Final escape-rule prohibition
    for the authored path, RDR 0001's reserved
    programmer-mistake error channel for the backstop —
    which is where it was already argued. SCXML's MAY/or
    structure is *permissive of* that split, not evidence
    for it.
  - **If wrong**: The normalizer-side half loses its
    external precedent and stands on RDR 0002's already-
    Final prohibition alone — the choice survives, the
    citation does not.
- **A7 The breach error can name the offending row
  (`RuleID`, `SourceLocator`) without violating RDR 0005's
  rule that CLI codes are never inferred by inspecting
  error strings: the error is diagnostic prose on the
  internal-error path, not a stable code carrier.**
  - **Status**: Verified — and the structured form is
    **adopted**, not merely permitted (see the DX finding
    below).
  - **Method**: Source Search
  - **Evidence**: Field names confirmed —
    `internal/resolve/resolve.go::Row` really does carry
    `RuleID string` and `SourceLocator string` ("the source
    identity RDR 0002 requires every normalized row to
    retain"), mirrored on `Plan` and `RowRef`. Scope of the
    no-string-inspection rule confirmed: RDR 0001's
    Normative Contracts say "**Each refusal kind** must be
    stable enough for RDR 0005 to map to a CLI error code
    without inspecting error strings," and the same
    sentence is attached to the `RefusalKind` type at
    `resolve.go`; RDR 0001's A5 names the rejected
    alternative as "letting RDR 0005 infer CLI codes from
    Go errors or ad hoc strings." The rule binds the
    refusal discriminator, which this breach explicitly is
    not. `clierr::CLIError` keeps prose and code in
    separate fields (`Detail` "carries hard facts about the
    failure", `Code` stable), and `clierr::ErrorCode` /
    `ExitCodeFor` branch on `Code` only, never on text.
  - **DX finding (Resolve)**: naming the row in *prose*
    would make this the only kernel diagnostic not
    assertable as a value. The kernel already carries row
    identity structurally (`RowRef`, `Refusal.Rows`) and
    its frozen tests assert on it by value
    (`fixup_test.go` `reflect.DeepEqual(forward.Refusal,
    reversed.Refusal)`), with zero nil-sensitive or
    string-matching assertions. Go's own convention agrees:
    `encoding/json/v2` documents the caller idiom as a
    sentinel in an `Err` field plus an exported location
    field, inspected via `errors.AsType`
    (`encoding/json/v2/errors.go::SemanticError`), as do
    `json.UnmarshalTypeError.Field` and `go/types.Error`.
    The Normative Contracts below therefore require a typed
    breach error carrying `[]RowRef`.
  - **If wrong**: The error must carry structure (a typed
    error), which sharpens the contract at Pre-Lock but
    does not move ownership.

- **A8 Collapsing equal `RowRef` identities is the only
  tiebreak available to the multi-breach report, because
  `compareRefs` does not totally order distinct rows and RDR
  0001's REQ-2/REQ-10 forbid breaking the tie by table
  position.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/resolve/resolve.go::compareRefs`
    compares `RuleID` then `SourceLocator` and returns 0 when
    both match; `::rowRefs` sorts with the non-stable
    `slices.SortFunc`. Equal identities are constructible and
    already present: `internal/resolve/fixup_test.go:79`,
    `:101`, `:125` each build
    `escapeRow("rdr.escape.ambiguous", "flows/rdr.toml:99",
    …)`, and hand-built rows without source identity all
    carry `RowRef{"",""}`. The positional tiebreak is
    excluded by the frozen suite, not merely by preference:
    `adversarial_test.go::TestAdv3b_MissingOwnedPayloadMustNotDependOnTableRowOrder`
    and
    `fixup_test.go::TestFixup3c_AmbiguousMatchRowsPayloadMustNotDependOnTableRowOrder`
    both fail a payload that varies with row order, and
    ADV-3b's comment rejects "stable order … with respect to
    row order, which is not the same as stable with respect
    to the input tuple." Fixup-3c grounds it in REQ-2: "two
    orderings of the same row set at the same revision are
    the same input."
  - **If wrong**: If a consumer needs a per-row breach count
    rather than per-identity, the report must carry a count
    field instead of collapsing — the ordering rule is
    unchanged either way.
  - **Raised by**: cove (Pre-Lock), finding L-2.

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
as a `Table` method so any producer can call it at
construction time, and called by `resolve.go::Resolve` at
entry before `assemble` — the same function at both call
sites. When any row carries both a non-empty `Escape` and a
non-empty `Writes`, `Resolve` returns the zero `Result`
(`Plan` and `Refusal` both nil — `Result` is a struct, so
there is no nil `Result` to return) and a non-nil error
naming every such row. Because the zero `Result` reports
`Refused() == false`, callers MUST check the error before
reading the disposition; a caller that branches on
`Refused()` first would read a success-shaped value with a
nil `Plan`. `planOf`
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
kernel boundary.) The predicate is Writes-only and does NOT
extend to NextTags: A4 settled that owned state is reachable
only through a write accessor, which RDR 0004 scopes to
"planned owned-tag writes," so an escape row's NextTags
cannot mutate owned state.
```

```normative
On the authored-table path, RDR 0002's normalizer/lint is the
enforcing implementation: an authored escape rule carrying a
write block or clear list MUST be rejected at table load
under RDR 0002's "malformed escape declaration" validation
class, diagnosable to the source rule. Rejection keys on the
PRESENCE of the block, not its contents — RDR 0002 forbids an
escape rule from containing a write block, so an empty one
(`writes = []`) is rejected too. The authored layer is
therefore stricter than the kernel's length-based predicate,
deliberately: authoring an empty write block on an escape
rule is a statement of intent worth refusing, while an empty
Writes slice reaching the kernel is a Go representation
detail that cannot mutate owned state. Normalized escape rows
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
The breach error MUST be a TYPED error carrying the offending
row identity as the kernel's existing RowRef value, inspectable
via errors.As / errors.AsType without parsing message text, and
MUST wrap a package-level sentinel naming the breach category so
errors.Is can classify it. Both the error type and the sentinel
MUST be EXPORTED: an unexported sentinel would defeat the
errors.Is classification this clause exists to provide, since the
intended callers (a future flow verb, non-kernel table producers)
are outside the package. Row identity MUST NOT be recoverable
only from formatted prose. (Form sharpened at Pre-Lock; the
binding requirement is value-level inspectability, matching the
kernel's existing Refusal.Rows diagnostics and Go's own
convention — encoding/json/v2 SemanticError, json.UnmarshalTypeError,
go/types.Error.)

Because errors.Join returns a wrapper even for a single error,
callers MUST classify and extract with errors.Is / errors.As /
errors.AsType rather than a direct type assertion or equality
against the sentinel. The aggregate is the uniform return shape:
Resolve MUST NOT return the bare per-row error in the
one-breach case and the aggregate otherwise, so caller code has
one shape to handle regardless of breach count.
```

```normative
When a table carries MORE THAN ONE breaching row, the
precondition MUST report every one of them in a single pass —
not the first alone — combined with errors.Join, so each
per-row error stays individually inspectable through the
aggregate's Unwrap() []error. The reported rows MUST be
ordered by RowRef identity using the kernel's existing
compareRefs ordering, never by Table.Rows position, so the
diagnostic payload is a function of the table value rather
than of row order (the REQ-1/REQ-10 rule rowRefs already
follows).

Because compareRefs orders on (RuleID, SourceLocator) only,
two distinct breaching rows sharing that pair compare EQUAL,
so identity alone does not totally order the report. Hand-built
rows are this precondition's whole target population (A5) and
collide trivially: rows built with no source identity all carry
RowRef{"",""}. Table.Rows position MUST NOT be used to break
the tie. RDR 0001 forbids it: REQ-2 makes two orderings of the
same row set THE SAME INPUT, and REQ-10 makes a diagnostic
payload a function of that input — the frozen suite tests the
rule directly (adversarial_test.go ADV-3b, fixup_test.go
Fixup-3c, both named "MustNotDependOnTableRowOrder"), and
ADV-3b's own comment rejects "stable with respect to row
order" as not stable with respect to the input tuple. A
positional tiebreak (SortStableFunc) would therefore
reintroduce exactly the defect those tests freeze.

The report is instead ordered by RowRef identity and made
total by COLLAPSING equal identities: breaching rows sharing
one RowRef contribute ONE reported error, not one per row.
Identity is what the producer fixes by — the report names the
distinct offending (RuleID, SourceLocator) pairs — so
collapsing loses no actionable information, and the payload
stays a function of the table value under any permutation.
Rows carrying no source identity collapse to a single
RowRef{"",""} entry; that is a degenerate producer, and the
error MUST remain diagnostic in that case by stating the
breach count alongside the identities.
```

```normative
The conformance predicate MUST be exported by the kernel
package as a construction-time check callable by any table
producer, and Resolve's entry precondition MUST be that same
function — one predicate, two call sites, so a non-TOML
producer can fail at build or construction time rather than
at first production Resolve, and the two enforcement points
cannot drift. It MUST be a METHOD ON Table taking no
arguments, matching both cited precedents (which are methods
on the value they validate) and the package's own habit of
exporting behavior on the type it belongs to (TagSet.Lookup,
Result.Refused). The whole-table scope is thereby carried by
the receiver rather than restated at each call site. It MUST
return error (nil when valid), and its name MUST follow Go's
error-returning convention — CheckValid or Validate, never
Valid/IsValid/OK, which Go reserves for bool-returning
predicates. Precedent: the
exported-and-also-called-defensively arrangement is
pprof.Profile.CheckValid (exported, re-checked at the parse
boundary) and rsa.PrivateKey.Validate; copy the latter's doc
form — "returns nil if <x> is valid, or else an error
describing a problem."
```

```normative
When a breach reaches the CLI, the verb MUST wrap it into a
*clierr.CLIError carrying a stable Code, Group GroupInternal
(exit 2), the offending row identity, and a Hint stating the
remedy — extending the shipped config.Load wrapping pattern
(Code + Detail + Group + Cause), which carries no Hint of its
own; this RDR adds one. Because clierr.CLIError.Cause is
json:"-", the Go error chain that carries the RowRef values
is NOT wire-visible: the offending row identities MUST reach
the envelope through a serialized field, either the existing
Detail (rendered from the row identities, never re-parsed by
any consumer) or a new omitempty field added under the type's
own "Extend with new optional fields as needed" allowance.
This is an additive envelope change, not a change to RDR
0005's refusal-to-exit-code mapping, so A2 stands. Returning
the kernel error UNWRAPPED is a defect: clierr.ExitCodeFor
defaults a non-CLIError to exit 1 and root.go's
cobraErrorToCLIError stamps GroupUserEnv, so the documented
exit-2 behavior does not hold without the wrap. Adding this
Code does NOT open RDR 0001's refusal taxonomy: the five
refusal kinds stay closed and gain no member — clierr codes
are a separate, explicitly extensible surface ("Add new
codes as needed", AGENTS.md).
```

```normative
The shared resolve.Row type keeps its single shape: escape
identity remains discriminated solely by a non-empty Escape
list. No ordinary/escape type split and no row-kind field is
introduced at the kernel boundary.
```

#### Load-Bearing Decisions

- **Selection / predicate** — the breach predicate over a
  single row is `len(row.Escape) != 0 && len(row.Writes) != 0`.
  A4 settled closed, so it does **not** widen to `NextTags`
  (owned state is reachable only through a write accessor,
  which RDR 0004 scopes to planned writes). Pinned:
  - *Scope* — evaluated over every row of the supplied table
    at `Resolve` entry, not only rows the resolution touches:
    a malformed table is malformed as a value, and checking
    only reached rows would make the breach appear and
    disappear with the input tuple.
  - *Nil-vs-empty* — length-based on purpose: a non-nil but
    empty slice is conforming. The predicate tests emptiness,
    never nil-ness. Three grounds: an empty write set is an
    empty work set — `copyTags` yields `len(Plan.Writes) == 0`
    for both `nil` and `[]Tag{}`, so the two are
    indistinguishable to any consumer that reads the plan's
    writes, and RDR 0004 scopes the write accessor to planned
    writes (a peer contract, not yet implemented — the
    kernel-local half is what is verified here);
    nil-sensitivity would make this the kernel's
    only nil-sensitive contract, against a package that is
    uniformly length-based (`len(row.Escape) != 0` is the
    escape discriminator itself) and a frozen suite whose
    `Writes` assertions are **all four** length-based, with
    **no** nil-ness assertion on `Writes` anywhere (the
    suite's nil checks are `Plan`/`Refusal` presence tests,
    a different question); and `copyTags` is nil-preserving, so
    nil-vs-empty is a representation artifact that survives
    into the Plan without behavioral meaning. Strictness
    about the *authored* form belongs upstream, where RDR
    0002 already forbids a write block on an escape rule
    outright — including an empty one.
  - *Precedence* — the breach error precedes every modeled
    disposition, `unmodeled_outcome` included.
  - *Locus* — one exported function, called at both sites, so
    any future change is a single-site edit that cannot
    desynchronize the two enforcement points.
  - *Reporting* — aggregate, not fail-fast: all breaching
    rows in one pass, sorted by `compareRefs`, with
    equal-identity rows collapsed to one entry (identity is
    not a total order over distinct rows, and RDR 0001's
    REQ-2/REQ-10 forbid breaking the tie by table position).
    Prior art is
    split by audience — `fstest.TestFS` and
    `slogtest.TestHandler` aggregate via `errors.Join` for
    artifacts a human iterates on, while `analysis.Validate`
    is fail-fast for startup wiring. A table author fixing
    rows is the former, and `errors.Join` costs nothing
    (nil when all-nil; degrades to the single error's own
    message for one breach).
- **Naming** — the breach is "escape-row shape breach,"
  carried as a typed Go error wrapping a category sentinel
  (A7). The predicate is `CheckValid`/`Validate`-shaped
  (returns `error`); the exact type/sentinel spelling is
  sharpened at Pre-Lock.

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
| Error envelope | `internal/cli/clierr::GroupInternal` / `ExitCodeFor` | Exit 2 is shared with other internal errors | Reuse | A stable `CLIError.Code` distinguishes the breach within the shared group; `clierr` codes are extensible, so no exit-group change (A2) |
| Escape fixtures | `internal/resolve/fixtures_test.go::escapeRow` (builder) + the `escape.Writes` overrides at `adversarial_test.go:229`, `fixup_test.go:103` | Builder sets `Writes` (inherited by 16 call sites); the two overrides replace the value | Extend (conform all three sites) | Fixtures must satisfy the precondition; supersedes the "Noted, not filed" cleanup |
| Row-identity diagnostics | `internal/resolve/resolve.go::RowRef` + `rowRefs`/`compareRefs` | Sorted for REQ-1/REQ-10 payload stability | Reuse | The typed breach error carries `RowRef`s in the same sorted order |

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
- Cost (measured against the existing path, not assumed):
  every `Resolve` call pays one additional O(rows)
  conformance scan. `Resolve` already walks `Table.Rows` for
  the candidate partition, and again in `escapeOrRefuse` on
  the zero-match and ambiguous arms, plus `gate` over the
  survivors — so the added scan is one more linear pass on
  an already-linear path, not a new cost class. It is not
  free at the margin: on the exact-one-match success path
  `escapeOrRefuse` never runs, so there the precondition
  doubles the number of full `Table.Rows` passes (one to
  two).
  RDR 0001's Performance Expectations set no throughput
  target and direct benchmarking "only if the RDR or kata
  tables become large enough to make table scans visible in
  normal command latency." If measurement ever demands it,
  the precedent is an idempotence short-circuit
  (`rsa.PrivateKey.Validate`), never deleting the check.
- Positive: a breach off the authored path surfaces as a
  typed error carrying the offending rows, wrapped at the
  CLI in a `CLIError` with a stable `Code`, exit 2, and a
  remedy `Hint`. Operators branch on the code and tests
  assert on `RowRef` values — neither reads error prose.
  This costs the refusal taxonomy nothing: the five kinds
  stay closed, and `clierr` codes are a separate,
  explicitly extensible surface.
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
  **Retired at Resolve**: A4 settled it — owned state is
  reachable only through a write accessor, which RDR 0004
  scopes to "planned owned-tag writes," and `NextTags`
  appears nowhere in RDR 0004. The residual is a
  peer-contract dependency, not a predicate gap: an RDR 0004
  implementer who widened "apply the plan" beyond planned
  writes would breach this invariant. Recorded on A4.

### Failure Modes

- **Visible break (programmatic producer)**: `Resolve`
  returns a typed error carrying the offending rows'
  `RowRef`s (every breaching row, not just the first); at
  the CLI it surfaces through a `CLIError` with a stable
  code, exit 2, and a remedy `Hint`. Diagnosis: read the
  row identities off the error value — the rendered code and
  detail for operators; for callers and tests, `errors.Is`
  classifies the category and the identities are recovered by
  walking `Unwrap() []error` over the aggregate. A single
  `errors.As`/`AsType` call reports only the FIRST breach in
  the chain, so it is the wrong instrument for reading a
  multi-breach report — a test asserting on all offending
  rows must traverse the aggregate. Recovery: fix the
  producer — the table value is broken; there is nothing to
  retry.
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
- **Diagnosis (revised at Resolve)**: the breach is
  machine-branchable without opening any closed taxonomy.
  Exit 2 is shared with other internal errors, but the
  stable `CLIError.Code` distinguishes this breach, the
  typed kernel error carries the offending `RowRef`s for
  programmatic callers and tests, and a `Hint` states the
  remedy. The runbook entry is therefore "branch on the
  code; the named rows identify the broken producer," not
  "read the error text." An earlier draft recorded this as
  an accepted diagnostic gap on the reasoning that a
  structured surface would require opening RDR 0001's
  refusal taxonomy; that conflated two distinct surfaces —
  the five refusal kinds (closed, untouched here) and
  `clierr` codes (explicitly extensible, AGENTS.md: "Add
  new codes as needed").
- **Error-channel creep guarded against**: this RDR's
  normative clause states what qualifies for the error path
  (a producer/programmer contract breach, never a modeled
  condition), so the next borrower argues against text, not
  precedent.

## Implementation Plan

### Prerequisites

- [x] All Critical Assumptions verified (Stage 4; A6
      verified contrary to its Propose-stage wording)
- [x] A3 spike green (conformed fixtures, full
      `internal/resolve` suite) before the precondition
      lands —
      `evidence/spikes/a3-fixture-conformance.md`
- [ ] RDR 0002's implement stage binds to this RDR's
      conformance fixtures (Phase 3) once both are Final
- [ ] The `flow` verb that first calls `Resolve` wraps a
      breach into `CLIError{Group: GroupInternal}` with the
      stable code and remedy `Hint` (A2's wiring
      obligation; unwrapped, it would exit 1, not 2), and
      renders the offending row identities into a serialized
      envelope field — `Cause` is `json:"-"`, so identities
      carried only on the Go error chain never reach an
      operator

### Minimum Viable Validation

A hand-constructed table containing an escape row with
populated `Writes`: `Resolve` returns a non-nil error naming
that row's `RuleID` and `SourceLocator`, with no `Plan` and
no `Refusal` — while the same table with the writes removed
resolves, escapes, and refuses exactly as the frozen ADV
suite proves today (zero disposition change for conforming
tables). The exported validator and the entry check are
exercised as the same predicate.

The three scenarios below are **normative fixtures**
(approved at Stage 4). Their values are read from the frozen
suite and the A3 spike runs, not invented: `{status Escaped}`
is the literal override at `adversarial_test.go:229`, and the
conforming baseline is the 154-PASS run captured in
`evidence/spikes/a3-baseline.out`. Identifiers below are
illustrative of shape, not locked spellings — the binding
content is the observable behavior.

1. **breach-yields-error-not-plan** — an escape row
   (`RuleID` `rdr.escape.needsowned`, `SourceLocator`
   `flows/rdr.toml:90`, `Escape` `[no_match]`) carrying
   `Writes: []Tag{{Key: "status", Value: "Escaped"}}`, in a
   table where it is selected via the escape path.
   **Expected**: `Resolve` returns the zero `Result`
   (`Plan: nil, Refusal: nil`) and a non-nil error; the error
   yields that row's `RowRef{"rdr.escape.needsowned",
   "flows/rdr.toml:90"}` through `errors.As`/`AsType`
   without parsing text, and `errors.Is` classifies it as
   the escape-shape-breach category — both through the
   `errors.Join` aggregate, which wraps even this
   single-breach case. The pre-fix behavior —
   a `Plan` with `Escaped: true` carrying those writes — is
   never produced.
2. **dormant-row-still-errors** — a table that resolves
   cleanly on its own, plus an escape row no resolution path
   reaches (an outcome outside the tuple's reach) carrying
   the same `Writes`.
   **Expected**: the same error, naming the dormant row —
   pinning the whole-table decision against an on-path-only
   reading.
3. **empty-not-nil-conforms** — an escape row whose `Writes`
   is a non-nil, zero-length slice (`[]Tag{}`).
   **Expected**: no error; the table resolves, escapes, and
   refuses identically to the A3 baseline, and the emitted
   plan carries `len(Plan.Writes) == 0`. Asserted on length,
   never on nil-ness — `copyTags` is nil-preserving, so a
   nil-ness assertion would pin an incidental representation
   detail rather than the contract.

### Phase 1: Exported predicate + kernel entry precondition

Export the conformance predicate as a construction-time
check — a no-argument method on `Table`,
`CheckValid`/`Validate`-shaped, returning `error`,
doc-commented in the `rsa.PrivateKey.Validate` form — call it
at `Resolve` entry (breach → exported typed error carrying
the offending `RowRef`s and wrapping the exported category
sentinel; conforming tables unchanged), aggregate breaches
with `errors.Join` in `compareRefs` order with equal
identities collapsed, state the producer
obligation on the `Row` doc contract, and record the
kernel's validation surface — exactly which shape property
it checks vs still assumes — so partial validation cannot be
read as general kernel ownership.

### Phase 2: Fixture conformance with discrimination check

Bring the write-bearing escape rows into conformance at all
three sites: drop the `Writes` field from the
`internal/resolve/fixtures_test.go::escapeRow` builder (line
257, inherited by its sixteen call sites), and drop the two
call-site overrides —
`internal/resolve/adversarial_test.go:229` (ADV-2b) and
`internal/resolve/fixup_test.go:103` (Fixup-1e "missing
owned state"). `NextTags` stays on the builder: A4 did not
widen the predicate. A3's spike proves the frozen suite still
discriminates after the strip (154 PASS, identical outcome
set; escape-gate mutation still fails 11 assertions). Because
no remaining fixture breaches, the mutation-style step
(revert the entry check; at least one test must fail) is
load-bearing here rather than a formality — it is what proves
the new check is exercised at all.

### Phase 3: Shared normalizer conformance fixtures

Encode the authored-path half as a named, shared fixture set
binding the future RDR 0002 build: an authored escape rule
with a write block or clear list fails load under "malformed
escape declaration" naming the source rule; an escape rule
carrying an **empty** write block (`writes = []`) fails load
the same way — RDR 0002 forbids an escape rule from
*containing* a write block, so presence is the breach and
emptiness is no defence (this is where authored-layer
strictness lives; the kernel predicate is deliberately
length-based one layer down); normalized escape rows render
write-free; and one canonical authored-clear case pins the
`<clear>`-sentinel-write representation identically for the
kernel suite and the normalizer suite, so the two enforcers
of one invariant cannot drift.

## Validation

### Testing Strategy

Done means: every producer path is closed, no conforming
table changes disposition, and the check itself is proven
non-vacuous. The MVV's three scenarios are the core; the
suite below is the coverage goal.

1. **Scenario**: Escape row with populated `Writes`, selected
   via the escape path. (Normative fixture 1.)
   **Expected**: `Resolve` returns a non-nil error; no `Plan`,
   no `Refusal`. The row's `RowRef` is recovered from the
   error value via `errors.As`/`AsType` and the category via
   `errors.Is` — no assertion parses message text.
2. **Scenario**: Escape row with populated `Writes` that no
   resolution path reaches (dormant). (Normative fixture 2.)
   **Expected**: Same error naming that row — the precondition
   is whole-table, not on-path.
3. **Scenario**: Escape row with a non-nil but empty `Writes`
   slice. (Normative fixture 3.)
   **Expected**: No error; resolves exactly as the A3
   baseline, with `len(Plan.Writes) == 0`. Asserted on length,
   never nil-ness — pins the length-based predicate.
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
   carrying a write block, one carrying a clear list, and one
   carrying an *empty* write block (`writes = []`).
   **Expected**: All three rejected at load under "malformed
   escape declaration," naming the source rule; normalized
   escape rows render write-free. The empty-write-block case
   is the authored-layer strictness that the kernel's
   length-based predicate deliberately does not duplicate.
9. **Scenario**: A table with two or more breaching rows.
   **Expected**: One `Resolve` call reports **every** breaching
   identity, not just the first; each per-row error is
   individually recoverable by traversing the aggregate's
   `Unwrap() []error` — a single `errors.As` finds only the
   first match, so the assertion must walk the slice.
10. **Scenario**: The same multi-breach table, rows carrying
    *distinct* `RowRef` identities, supplied in two different
    orders.
    **Expected**: Identical reported row sequence — ordered by
    `RowRef` identity (`compareRefs`), never by `Table.Rows`
    position. The REQ-1/REQ-10 payload-stability rule that
    `rowRefs` already follows, applied to the breach report.
10b. **Scenario**: A multi-breach table whose breaching rows
    share one `RowRef` identity (including the zero-value
    `RowRef{"",""}` of rows built without source identity),
    supplied in two different orders.
    **Expected**: Identical report under both permutations —
    equal identities collapse to one reported entry, and the
    breach count is stated alongside. Pins the payload as a
    function of the input tuple (REQ-2/REQ-10) rather than of
    row position, the rule ADV-3b and Fixup-3c freeze.
11. **Scenario** (binds the verb that first calls `Resolve`): a
    breach surfaced through the CLI.
    **Expected**: exit 2 via `CLIError{Group: GroupInternal}`
    carrying the stable code, the offending row identity, and
    the remedy `Hint` — not exit 1, which is what an
    unwrapped kernel error would produce. The row identities
    must be readable from the **serialized** envelope, not
    only from the Go error chain: `CLIError.Cause` is
    `json:"-"`, so an assertion that reads identities off the
    JSON output is what proves the wire carrier exists.

### Performance Expectations

No throughput target, inheriting RDR 0001's position that
resolution is bounded by the supplied table and that
benchmarking waits "only if the RDR or kata tables become
large enough to make table scans visible in normal command
latency."

The precondition adds one O(rows) pass over `Table.Rows` per
`Resolve` call. Measured against the existing path rather
than in the abstract, this is a constant factor on an
already-linear function: `resolve.go::Resolve` already walks
`Table.Rows` for the candidate partition, `escapeOrRefuse`
walks it again for escape candidates, and `gate` walks the
survivors. The scan is allocation-free in the conforming case
— it reads `len(row.Escape)` and `len(row.Writes)` and
allocates only when a breach is found and a `RowRef` is
recorded — which keeps the kernel "allocation-light" as RDR
0001 asks.

The cost is per-call because the kernel is stateless by
contract; no caching or memoization is introduced, since a
cache keyed on a caller-supplied value would reintroduce the
ambient state RDR 0001 forbids. If measurement ever shows the
scan mattering, the precedent is an idempotence
short-circuit on an already-validated table
(`crypto/rsa::PrivateKey.Validate` returns early when its
precomputed values are consistent), never removing the check.

No byte-stable hash, canonical serialization, or ordering
guarantee is introduced, so the determinism checklist does
not apply. The one ordering claim this RDR does make — that
a multi-breach report is ordered by `RowRef` identity rather
than table position — is a value-level determinism property
covered by scenario 10 and by the existing `compareRefs`
ordering, not a byte-level one.

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
- `internal/cli/clierr/clierr.go` (`CLIError.Detail`,
  `.Hint`, `.Cause`; `ExitCodeFor` non-`CLIError` → exit 1)
- `internal/cli/config/config.go::Load` (shipped precedent:
  Go error → `GroupInternal` + stable code + `Hint`)
- `internal/cli/root.go` (`NewRootCmd` registers only
  `version`; `cobraErrorToCLIError` stamps `GroupUserEnv`)
- `AGENTS.md` ("Every user-facing failure is a
  `*clierr.CLIError` with a stable `Code` … Add new codes as
  needed")
- Go stdlib DX precedent: `encoding/json/v2::SemanticError`
  (sentinel + location field, `errors.AsType`),
  `encoding/json::UnmarshalTypeError.Field`,
  `go/types::Error`, `testing/fstest::TestFS` and
  `testing/slogtest::TestHandler` (`errors.Join`
  aggregation), `errors::Join`,
  `cmd/vendor/.../pprof/profile::Profile.CheckValid`
  (exported + re-checked at the boundary),
  `crypto/rsa::PrivateKey.Validate` (doc form; idempotence
  short-circuit)
- W3C SCXML IRP test 313 (enforcement point is MAY, not
  standard) via `state-machines/study/uscxml/test/w3c/TESTS.md`
- Prior-art query ledgers:
  `docs/rdr/0009-escape-row-shape-conformance-ownership/evidence/research/propose-prior-art.md`,
  `.../evidence/research/resolve-prior-art.md`,
  `.../evidence/research/resolve-dx-prior-art.md`
- A3 spike:
  `docs/rdr/0009-escape-row-shape-conformance-ownership/evidence/spikes/a3-fixture-conformance.md`
