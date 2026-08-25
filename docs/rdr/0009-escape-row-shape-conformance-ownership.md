# Recommendation 0009: Ownership of escape-row shape conformance

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-08-09
- **Status**: Final [joint decision → JDR 0001 §JD-5: precondition precedence]
  [§JD-8 answered 2026-08-24 by §D10 — the code table, the exit-3 rule, and the
  `omitempty` `findings` list carry the row identities; `GroupInternal` already
  ships. Citations owed at Stage 8, see `artifacts/deviations.md`.] [§JD-15
  answered 2026-08-24 by §D5, checked consistent at the 0002-0009 iteration-3
  gate — 0009 unchanged]
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
    verb's `RunE` still exits **2**, but misclassified:
    `root.go::ExecuteAndEmit` converts every non-`CLIError`
    via `cobraErrorToCLIError`, stamping `GroupUserEnv`
    (exit 2), `Code: "command-error"`, and the generic
    `--help` hint — so the exit code cannot distinguish a
    wrapped breach from a missing wrap (corrected at
    Pre-Lock, critique: the earlier exit-1 reading stopped
    at `ExitCodeFor`'s non-`CLIError` default, a path
    `ExecuteAndEmit` never lets a verb error reach). The
    wrap supplies the internal classification, the stable
    code, the row identities, and the remedy `Hint`; the
    Normative Contracts below bind that wrap, and
    conformance is asserted on the `Code`, not the exit.
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
    **Re-confirmed at Stage 6** (critique carried the
    corrected exit-2 wiring claim forward for
    re-verification): A9's source search closed the additive
    question against `internal/cli/clierr/clierr.go` — the
    `CLIError` type's own doc comment authorizes `omitempty`
    extension, no consumer or test reads an exact field set,
    and `ExitCodeFor`'s refusal-to-exit-code mapping is
    untouched by the addition. The corrected reading
    (`ExecuteAndEmit` → `cobraErrorToCLIError` → exit 2,
    never exit 1) stands.
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
  - **If wrong**: If a consumer needs a per-row breach *entry*
    rather than a per-identity one, the report must carry
    one entry per row instead of collapsing — the ordering
    rule is unchanged either way. (Note this is not the
    per-identity `Count` field, which the Normative Contracts
    require *alongside* collapsing; the rejected branch is
    per-row entries, which would reintroduce the ordering
    defect REQ-2/REQ-10 forbid.)
  - **Raised by**: cove (Pre-Lock), finding L-2; count/collapse
    contradiction sharpened by 3amigo (Pre-Lock), T-1.

- **A9 Adding a new `omitempty` field to `clierr.CLIError` to
  carry the offending row identities is additive to RDR 0005's
  envelope and breaks no shipped consumer — even though
  `Cause`'s own doc comment currently states "the wire-visible
  cause surface is Detail."**
  - **Status**: Verified — with a named amendment obligation
    (the `Cause` doc comment; see (c))
  - **Method**: Source Search
  - **Evidence**: All four plan items closed at Stage 6.
    **(a) No consumer depends on `Detail` being the sole wire
    cause surface.** `internal/cli/` carries exactly one test
    file, `internal/cli/version_test.go`, and no `testdata/`,
    golden files, or snapshot fixtures exist anywhere in the
    CLI tree. Its only JSON assertion,
    `version_test.go::TestVersion_JSON`, unmarshals into a
    *partial anonymous struct* (`Type`, `Data.Version`), which
    ignores unknown fields — and it tests the `respond.Success`
    envelope, not `CLIError`. The two error-path tests
    (`::TestUnknownCommand_IsStructured`,
    `::TestInvalidMode_IsRefused`) assert only on
    `clierr.ErrorCode` / `clierr.ExitCodeFor` and never inspect
    serialized bytes. The shapes that *would* break on an added
    field — whole-JSON-string comparison, `DeepEqual` over an
    unmarshalled map, `DisallowUnknownFields` — exist nowhere in
    the repo (every `reflect.DeepEqual` hit is in
    `internal/resolve/*_test.go` over kernel types). No non-Go
    consumer parses the envelope either. So the addition breaks
    nothing whether or not the field is populated.
    **(b) The allowance is real and conditional.**
    `internal/cli/clierr/clierr.go::CLIError`'s doc comment
    reads verbatim: "Extend with new optional fields as needed —
    keep them `omitempty` so the envelope stays append-only and
    stable for tools." The condition (`omitempty`) is exactly
    what this RDR's clause already requires. Corroborated by RDR
    0005 ("Error envelopes stay append-only through `CLIError`
    fields"); its reuse-audit "no new envelope fields required"
    row scopes 0005's own MVP verbs and is not a prohibition.
    Serialization is a plain `json.Marshal(e)` in
    `clierr.go::EmitJSON` with no custom `MarshalJSON`, so a
    tagged field reaches the wire with no other code change.
    Field tags confirmed: `Code`/`Message` unconditional,
    `Param`/`Detail`/`Hint` already `omitempty`, `Group` and
    `Cause` both `json:"-"`.
    **(c) The `Cause` doc comment must be amended alongside —
    and this RDR authorizes only the amendment, not a contract
    change.** Verbatim at
    `internal/cli/clierr/clierr.go::CLIError.Cause`: "Cause
    preserves the underlying Go error for errors.Is/errors.As
    traversal. Not serialized — the wire-visible cause surface
    is Detail." The definite article makes that a uniqueness
    claim, and row identities are cause information, so the
    sentence becomes false once the field lands. The defect is
    documentary, not behavioral: nothing branches on the claim
    and `Cause` stays `json:"-"`. Recorded as a Prerequisite
    obligation to reword the second clause (and to update
    `docs/cli-output-contract.md`'s envelope field list) in the
    same change. The comment's normative home is RDR 0005
    (which names `internal/cli/clierr::CLIError` in its Touched
    Surfaces); amending a stale *code comment* to match an
    additive field is not reopening 0005's envelope contract,
    which is what A2's verdict turns on.
    **(d) `clierr` stays a leaf.** `go list` gives its complete
    import set as four stdlib packages — `encoding/json errors
    fmt io` — with no intra-repo import at all, and its package
    doc records that as deliberate (it exists "so other internal
    packages … can construct CLIErrors without importing
    internal/cli and forming an import cycle"). Nothing forces a
    `resolve` import; `internal/resolve/resolve.go` imports only
    `slices` and `strings` and could not absorb the reverse
    dependency either. Note the direction guard runs one way —
    `internal/resolve/resolve_test.go::TestReq11_KernelImportsNoCLIOutputOrPersistenceFacility`
    forbids the kernel importing `clierr`, not the converse — so
    the leaf property is preserved by design intent, not by a
    test. **Settled: the carrier is a clierr-local
    representation** (plain strings or a small row-identity
    struct declared in `clierr`), never `resolve.RowRef`; the
    CLI verb layer, which already imports both, does the
    conversion. **Settled: the per-identity `Count` DOES
    serialize alongside the identities** — the kernel carries it
    structurally precisely so no consumer parses prose for it,
    and dropping it at the wire would force exactly that
    re-parse for the degenerate `RowRef{"",""}` case the
    multi-breach clause calls out.
  - **If wrong**: the identities ride `Detail` as rendered
    prose after all (accepting a re-parse consumers are told
    not to perform), or RDR 0005's envelope contract must be
    reopened — which would cost A2 its "no RDR 0005 contract
    change" verdict.
  - **Raised by**: 3amigo (Pre-Lock), finding T-3; closed at
    Stage 6 (reconcile).

- **A10 `Resolve` takes exactly one parameter and the table it
  checks is reached through it, so the entry precondition adds
  no parameter and changes no arity: the signature is
  `func Resolve(in Input) (Result, error)` and the table is
  `in.Table`.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/resolve/resolve.go::Resolve` is
    declared `func Resolve(in Input) (Result, error)` — one
    parameter. `internal/resolve/resolve.go::Input` carries
    `Table Table` alongside `Flow`, `Owned`, `Observed`,
    `Recognized`, and `Guards`, and `Resolve`'s own body already
    reaches the rows as `in.Table.Rows` for the candidate
    partition. So `in.Table.CheckValid()` is a call on a value
    the function already holds; no signature change is implied
    by the precondition.
  - **If wrong**: the precondition cannot be a first statement
    on the existing entry point without changing RDR 0001's
    locked arity, which this RDR does not authorize and which
    would break all sixteen `escapeRow` call sites and the
    frozen suite — the placement decision and the
    no-peer-change claim would both have to be reopened.
  - **Raised by**: repeatability (Pre-Lock), diff finding D-1
    (run-1 reconstructed a two-parameter `Resolve`).

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
precondition of Resolve. Resolve's signature is UNCHANGED —
`func Resolve(in Input) (Result, error)`, one parameter — and
the checked table is the one already reached through the input,
`in.Table`; the precondition adds no parameter and takes no
table argument of its own. (Stated because the RDR otherwise
names the table only as "the supplied table": a reconstruction
that reads the precondition as taking its own Table parameter
changes RDR 0001's locked entry-point arity, which this RDR
does not authorize.) A table containing a row that
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
only from formatted prose. (The binding requirement is
value-level inspectability, matching the kernel's existing
Refusal.Rows diagnostics and Go's own convention —
encoding/json/v2 SemanticError, json.UnmarshalTypeError,
go/types.Error.)

The surface is fixed here, not left to the implementer, because
it lands permanently on RDR 0001's locked package surface and
every test below asserts through it:

  - the sentinel is ErrEscapeShapeBreach (package-level, exported);
  - the typed error is *EscapeShapeBreachError (exported), carrying
    exactly ONE RowRef field, Ref, plus a Count int field (see the
    multi-breach clause);
  - its Unwrap() error returns ErrEscapeShapeBreach, so errors.Is
    classifies EVERY per-row element, not only the aggregate;
  - the predicate is Table.CheckValid() error.

CheckValid is chosen over Validate because the whole-table scan is
a cheap structural check re-run defensively at a boundary, which is
pprof.Profile.CheckValid's exact arrangement; rsa.PrivateKey.Validate
remains the doc-form precedent only. These four spellings collide
with no frozen boundary guard: the kernel's banned exported-name
lists (REQ-25, REQ-36, REQ-37) are exact-match, and errors/fmt are
forbidden by no import guard (REQ-26, REQ-28, REQ-37).

Because errors.Join returns a wrapper even for a single error,
callers MUST classify and extract with errors.Is / errors.As /
errors.AsType rather than a direct type assertion or equality
against the sentinel. The aggregate is the uniform return shape:
Resolve MUST NOT return the bare per-row error in the
one-breach case and the aggregate otherwise, so caller code has
one shape to handle regardless of breach count.

The aggregate's structure is likewise fixed, because the
multi-breach tests traverse it: Unwrap() []error MUST be EXACTLY
ONE LEVEL deep, every element being a *EscapeShapeBreachError —
never a nested join. Resolve MUST return CheckValid's error
VERBATIM, never fmt.Errorf-wrapped: an added layer would expose
Unwrap() error at the outermost level and break the flat
traversal this clause guarantees.
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
follows). REQ-2/REQ-10 bind modeled refusal payloads; the
breach report sits on the error path outside that surface, so
this clause is a DELIBERATE adoption of the same input-tuple
determinism for the diagnostic, not an inherited obligation.

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

The count is carried STRUCTURALLY, as the Count field on each
per-identity *EscapeShapeBreachError — never only in formatted
prose, which would be readable solely by the string matching this
RDR forbids everywhere else and that the frozen suite never does.
Count is PER-IDENTITY and counts PRE-COLLAPSE ROWS: three
breaching rows sharing RowRef{"",""} yield ONE reported error
with Ref == RowRef{"",""} and Count == 3. A single-row breach
carries Count == 1, so the field is uniform rather than present
only in the degenerate case.

Collapsing and the count are COMPANIONS, not alternatives:
collapsing keeps the report a function of the input tuple, and
Count restores the multiplicity collapsing would otherwise
discard. A8 records the narrower contingency — a consumer needing
per-row ENTRIES rather than per-identity ones — which stays the
rejected branch, since per-row entries reintroduce the ordering
defect REQ-2/REQ-10 forbid.
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
describing a problem." The doc comment MUST also name the ONE
property checked — escape-row shape conformance — and state
that nothing else is checked (in particular not the
Escape-class restriction to no_match/ambiguous_match that
Row's doc records for RDR 0002), so a nil return is never
read as general table validity.

A table with no rows (empty or nil Rows) conforms VACUOUSLY:
CheckValid MUST return nil, since no row can breach a predicate
quantified over rows and errors.Join of nothing is nil. Stated
because the kernel already admits that shape — REQ-20's
empty-input region reaches it — and a nil return there is the
vacuous truth, never a signal that validation was skipped.
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
the envelope through a serialized field. The chosen carrier is
a NEW omitempty field added under the type's own "Extend with
new optional fields as needed" allowance — NOT the existing
Detail, because reading identities out of Detail would require
the consumer to re-parse rendered prose, which this RDR forbids
everywhere else. The clause binds the carrier CLASS only; the
field's NAME stays the implementer's. A9 settled the rest at
Stage 6: the field's type MUST be a clierr-local
representation (plain strings or a small row-identity struct
declared in clierr) and MUST NOT be resolve.RowRef, keeping
clierr a leaf that does not import internal/resolve — the CLI
verb layer, which already imports both, does the conversion;
and the per-identity Count MUST serialize alongside the
identities, since carrying it structurally in the kernel and
then dropping it at the wire would force exactly the prose
re-parse this RDR forbids. The stable
Code is "escape-row-shape-breach".
This is an additive envelope change, not a change to RDR
0005's refusal-to-exit-code mapping, so A2 stands. A9 verified
the addition against clierr.CLIError.Cause's shipped doc
comment naming Detail the sole wire-visible cause surface: the
addition is authorized by the CLIError type's own "Extend with
new optional fields as needed — keep them omitempty" allowance,
and the Cause comment's now-stale second clause MUST be
amended in the same change (a Prerequisite; amending a stale
code comment is not reopening RDR 0005's envelope contract).
Returning
the kernel error UNWRAPPED is a defect the exit code does NOT
reveal: root.go's ExecuteAndEmit converts every non-CLIError
via cobraErrorToCLIError — GroupUserEnv, still exit 2, Code
"command-error", the generic --help hint, no row identities.
The wrap supplies the internal classification, the stable
code, the identities, and the remedy; conformance MUST be
asserted on the Code, never on the exit code alone. Adding this
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
  - *Placement* — the call is `Resolve`'s first statement,
    `in.Table.CheckValid()`, above the existing
    `view := assemble(in)`, so a breaching
    table costs no view assembly. This is a **cheapness
    preference, not an observable contract**: `assemble` is
    pure and allocation-only, so placing the check after it
    would produce identical dispositions. Nothing downstream
    may come to depend on the ordering.
  - *Doc-contract amendment* — `resolve.go::Resolve`'s
    numbered evaluation-order list (a REQ-1 contract artifact)
    currently opens at the alphabet check, so the whole-table
    precondition MUST be prepended to it as a new step 0.
    This RDR authorizes that edit: leaving the list unamended
    would ship a doc comment that is actively wrong about
    evaluation order. No frozen test asserts on the comment
    text, so the amendment is documentation-only — the
    behavioral change it describes is the precedence pin
    above, which scenario 4 covers per kind. The same
    authorization covers the two disposition-cardinality doc
    contracts the breach return would otherwise contradict —
    `Result`'s "never both and never neither" and `Resolve`'s
    "returns exactly one disposition" — each amended to scope
    itself to nil-error returns. The frozen cardinality test
    (`resolve_test.go::TestReq1_DispositionIsExactlyOneOfPlanOrRefusal`)
    is unaffected: its inputs are conforming tables and its
    helper checks the error before reading the disposition
    (scenario 7b pins the breach-side behavior).
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
  (A7). Spellings pinned (Pre-Lock, 3amigo T-2): sentinel
  `ErrEscapeShapeBreach`; typed error `*EscapeShapeBreachError`
  with fields `Ref RowRef` and `Count int`; predicate
  `Table.CheckValid() error`. `CheckValid` over `Validate`
  because the check is a cheap structural scan re-run
  defensively at a boundary — `pprof.Profile.CheckValid`'s
  arrangement exactly; `rsa.PrivateKey.Validate` stays the
  doc-form precedent only. All four clear the kernel's frozen
  exported-name and import guards (REQ-25/26/28/36/37), which
  are exact-match lists naming none of them.

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
  **Staged, not immediate.** Phases 1–2 close the path that is
  *reachable today* (hand-built rows; A5 verified no production
  constructor exists at HEAD), which is the whole reachable
  population. The authored path the Problem Statement's table
  author actually travels is closed by RDR 0002's normalizer,
  which is Final but unimplemented — Phase 3 ships the binding
  fixtures, not the enforcement. So at this RDR's own
  completion the invariant holds everywhere a row can be
  constructed; it becomes an *authoring-time* guarantee only
  when RDR 0002 is built. This RDR is done when Phases 1–3
  land; the table author's load-time diagnostic arrives with
  RDR 0002 (Prerequisites).
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
- Cost (author friction, weighed rather than assumed): the
  strict readings — whole-table scope here, and Phase 3's
  rejection of an *empty* authored write block — both refuse
  tables that could not actually misbehave. That is deliberate:
  a dormant malformed row is a producer bug whose latency is
  the hazard, and an authored `writes = []` on an escape rule
  is a statement of intent worth refusing. The friction is
  bounded and lands on the right party — today on nobody (A5:
  no production constructor), and after RDR 0002 ships on the
  table author, at load, naming their rule. A5's protection
  expires exactly when authored tables arrive, which is also
  when the load-time diagnostic that makes the refusal
  actionable arrives, so the strictness never lands without
  its remedy. If the empty-block refusal proves to be pure
  friction in practice, it is RDR 0002's to relax — the kernel
  predicate is length-based one layer down and would not move.

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

- [x] All Critical Assumptions verified — A1–A8 verified at
      Stage 4 (A6 verified contrary to its Propose-stage
      wording), A10 verified at Pre-Lock (repeatability), A9
      verified at Stage 6 (reconcile) with the `Cause`
      doc-comment amendment named below. None Pending
- [x] A3 spike green (conformed fixtures, full
      `internal/resolve` suite) before the precondition
      lands —
      `evidence/spikes/a3-fixture-conformance.md`
- [ ] RDR 0002's implement stage binds to this RDR's
      conformance fixtures (Phase 3) once both are Final
- [ ] RDR 0004's implement stage keeps the write accessor
      scoped to *planned owned-tag writes* — the peer contract
      A4's verdict rests on. This RDR guarantees an escape plan
      carries no `Writes`; it does **not** guarantee anything
      about a plan's `NextTags`, so an implementer who widened
      "apply the plan" beyond planned writes would reintroduce
      the exact user symptom this RDR exists to prevent
      (A4, *Carried forward*). Recorded as a binding obligation
      rather than a note because RDR 0004 is the layer where
      the harm actually occurs. The obligation has a test form
      0004's implement stage copies: a write accessor handed an
      escaped plan with empty `Writes` and non-empty `NextTags`
      applies zero owned-tag mutations
- [ ] The `flow` verb that first calls `Resolve` wraps a
      breach into `CLIError{Group: GroupInternal}` with the
      stable code and remedy `Hint` (A2's wiring
      obligation; unwrapped it still exits 2 but as
      `command-error`/`GroupUserEnv` with no row identities —
      the stable code, not the exit, is the tripwire), and
      renders the offending row identities into a serialized
      envelope field — `Cause` is `json:"-"`, so identities
      carried only on the Go error chain never reach an
      operator
- [ ] The same change that adds the `omitempty` identity
      field to `clierr.CLIError` amends the now-stale second
      clause of `internal/cli/clierr/clierr.go::CLIError.Cause`'s
      doc comment ("Not serialized — the wire-visible cause
      surface is Detail"), which the addition makes false, and
      updates `docs/cli-output-contract.md`'s error-envelope
      field list alongside. A9 verified this is a stale-comment
      correction riding an additive field, not a reopening of
      RDR 0005's envelope contract — but shipping the field
      without the amendment leaves the package documenting the
      opposite of what it does

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
`evidence/spikes/a3-baseline.out`. Fixture *names* below are
illustrative; the *exported identifiers* they assert through
(`ErrEscapeShapeBreach`, `*EscapeShapeBreachError`,
`Table.CheckValid`) are pinned by the Normative Contracts.

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
   `ErrEscapeShapeBreach` — both through the `errors.Join`
   aggregate, which wraps even this single-breach case. The
   aggregate's `Unwrap() []error` has length **1**; its single
   `*EscapeShapeBreachError` carries `Count == 1`. The pre-fix
   behavior —
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
obligation on the `Row` doc contract, prepend the precondition
to `Resolve`'s numbered evaluation-order doc list as step 0
(otherwise that list ships wrong), and record the
kernel's validation surface — exactly which shape property
it checks vs still assumes — so partial validation cannot be
read as general kernel ownership.

Exported surface (pinned above): `ErrEscapeShapeBreach`,
`*EscapeShapeBreachError{Ref RowRef; Count int}`, and
`Table.CheckValid() error`. The package gains `errors` (and
`fmt` for the message); both clear the frozen import guards.

Phases 1 and 2 land as ONE change: the moment the entry check
exists, every write-bearing escape fixture errors at `Resolve`
entry and `mustResolve` fatals across the frozen suite (the
`escapeRow` builder feeds sixteen call sites), so Phase 1
alone turns the suite red mid-plan.

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

Done means, in user terms: **an escape row can no longer carry
an owned-tag write to the accessor layer, so no tag value can
appear in owned state that no authored rule set** — the symptom
the Problem Statement opens with. The checkable form of that is:
every producer path is closed, no conforming table changes
disposition, and the check itself is proven non-vacuous. The
MVV's three scenarios are the core; the suite below is the
coverage goal.

**Scope of this RDR's Done criteria**: scenarios 1–7, 9, 10,
and 10b execute in this implementation. Scenario 8 binds RDR
0002's build and scenario 11 binds the future `flow` verb;
both are **deferred** — specified here so the obligation is
inherited, not run here. This RDR does not wait on them.

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
   **Expected**: No error; the table resolves, escapes, and
   refuses identically to the same table built with `Writes`
   nil, and the emitted plan carries `len(Plan.Writes) == 0`.
   Asserted on length, never nil-ness — pins the length-based
   predicate. The **discriminating** assertion is on the
   *input row* (`row.Writes != nil && len(row.Writes) == 0`
   yields no error): asserting only on the plan would not
   separate the two states this scenario exists to separate,
   since `copyTags` is nil-preserving and both reach the plan
   as `len == 0`. The non-nil-empty variant appears in **no**
   A3 run (the spike ran original / `Writes`-stripped /
   `Writes`+`NextTags`-stripped, all nil-valued), so this is
   new evidence rather than a re-assertion of A3.
4. **Scenario**: Breach precedence over each modeled
   disposition, run once **per refusal kind** — a breaching
   row coexisting with a table state that would otherwise
   yield `unmodeled_outcome`, `no_match`, `ambiguous_match`,
   `owned_state_unavailable`, and `guard_unevaluable`
   respectively.
   **Expected**: The breach error wins in every case, with no
   `Result` disposition. Exhaustive over the five kinds
   because the contract says "precedes **every** modeled
   disposition" and the house style for a five-kind claim is
   exhaustive enumeration (`mvv_test.go` leg 2 asserts
   `len(covered) != 5`). `unmodeled_outcome` alone is the
   weakest case — it is `Resolve`'s first branch; the sharp
   cases are those reached through `escapeOrRefuse`, where
   `gate` could otherwise return a refusal derived from the
   breaching row itself.
5. **Scenario**: The full frozen `internal/resolve` suite
   (ADV/MVV/boundary) run against conformed escape fixtures.
   **Expected**: Every assertion still passes. The comparison
   baseline is the A3 spike's **conformed** run
   (`a3-conformed.out`: 154 PASS, 0 FAIL), restricted to the
   tests existing at that revision — **not** `a3-baseline.out`
   (the pre-conformance tree) and not the raw post-Phase-2
   count, which necessarily exceeds 154 because this RDR adds
   tests. The oracle is the sorted outcome set over that
   pre-existing test set, which must be identical.
   "Still discriminates" is **not** asserted by this scenario —
   `go test` reports pass/fail, not discrimination. That
   property is scenario 6's job. The baseline is regenerable,
   not frozen to the spike artifact: re-run the frozen suite
   at the implementation's base commit immediately before
   Phase 1 lands and diff sorted outcome sets over the tests
   existing at that commit — the spike file records an
   instance of the oracle, not its definition.
6. **Scenario**: Mutation check — the check's non-vacuity.
   Two named mutants, run separately, with this RDR's **own**
   new tests (scenarios 1–4, 7, 9, 10, 10b) **excluded** from
   the oracle — including them makes the result a tautology,
   since they were written to assert exactly this check.
   - *Mutant A* — `Table.CheckValid` returns `nil`
     unconditionally. **Expected**: the frozen suite still
     passes. This is the honest result, and it is the point:
     after Phase 2 no pre-existing fixture breaches, so the
     frozen suite alone cannot detect the check's removal.
     Recording this refutes the weaker "at least one test
     fails" claim rather than hiding behind it.
   - *Mutant B* — re-introduce the breach into the fixtures
     (restore `Writes` on `fixtures_test.go::escapeRow`) with
     the entry check intact. **Expected**: the suite fails,
     proving the check is wired into the real evaluation path
     and not dead code.
   Together these establish what scenario 6 exists to
   establish: the check is live (B) and the conformed frozen
   suite is *not* its oracle (A) — which is precisely why
   scenarios 1–4 must exist.
7. **Scenario**: The exported predicate called directly by a
   producer at construction time, on the same tables as
   scenarios 1–3.
   **Expected**: Verdicts identical to `Resolve`'s entry check
   — one predicate, two call sites, no drift. "Identical" is
   defined as: both nil, or both non-nil with equal extracted
   `[]RowRef` **and** equal per-identity `Count` values, in
   the same order. Not `reflect.DeepEqual` over the two error
   values — `errors.Join` holds a slice, so DeepEqual turns on
   element identity rather than on the contract, and can
   report a false red or a false green depending on
   construction details this RDR deliberately does not pin.
7b. **Scenario**: The zero-`Result` caller trap on a breach.
   **Expected**: `Resolve` returns `Result{Plan: nil, Refusal:
   nil}` alongside the non-nil error, so `Refused() == false`
   on a call that did not succeed. Asserted directly, because
   it is the one newly-introduced hazard for callers: a caller
   branching on `Refused()` before checking the error reads a
   success-shaped value carrying a nil `Plan`. The frozen
   helpers happen to be safe (`mustResolve` fatals on error
   first), so no existing test covers this.
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
   first match, so the assertion must walk the slice. The
   traversal is **exactly one level deep** and every element
   type-asserts to `*EscapeShapeBreachError`; `errors.Is(elem,
   ErrEscapeShapeBreach)` holds on **each element**, not only
   on the aggregate. The assertion compares the extracted
   `[]RowRef` — never `reflect.DeepEqual` over two
   `errors.Join` values, whose equality depends on element
   identity rather than on the contract.
10. **Scenario**: The same multi-breach table, rows carrying
    *distinct* `RowRef` identities, supplied in two different
    orders.
    **Expected**: Identical reported row sequence — ordered by
    `RowRef` identity (`compareRefs`), never by `Table.Rows`
    position. The REQ-1/REQ-10 payload-stability rule that
    `rowRefs` already follows, applied to the breach report.
    The table MUST be **mixed** — at least one row carrying
    source identity and at least one built without it
    (`RowRef{"",""}`) — so the report's ordering is pinned
    across the two identity classes rather than only within
    one. Scenarios 10 and 10b are otherwise homogeneous (all
    distinct, all degenerate), which would leave the relative
    position of the zero-value identity unasserted;
    `compareRefs` compares strings, so it sorts first, and this
    scenario is what holds that.
10b. **Scenario**: A multi-breach table whose breaching rows
    share one `RowRef` identity (including the zero-value
    `RowRef{"",""}` of rows built without source identity),
    supplied in two different orders.
    **Expected**: Identical report under both permutations —
    equal identities collapse to one reported entry. For three
    rows sharing `RowRef{"",""}`: `Unwrap() []error` has
    length **1**, its single element has `Ref ==
    RowRef{"",""}` and `Count == 3` (pre-collapse rows, read
    off the struct field — no assertion parses message text).
    Pins the payload as a function of the input tuple
    (REQ-2/REQ-10) rather than of row position, the rule
    ADV-3b and Fixup-3c freeze.
11. **Scenario** — **DEFERRED: not executable by this RDR**
    (binds the future verb that first calls `Resolve`; no
    `flow` verb exists at HEAD, per A2/A5). A breach surfaced
    through the CLI.
    **Expected**: exit 2 via `CLIError{Group: GroupInternal}`
    carrying `Code: "escape-row-shape-breach"`, the offending
    row identity, and the remedy `Hint` "fix the table
    producer: an escape row must carry no writes" — asserted
    on the `Code`: an unwrapped kernel error also exits 2
    (`ExecuteAndEmit` → `cobraErrorToCLIError`,
    `command-error`, `GroupUserEnv`, generic hint), so the
    exit code cannot detect a missing wrap; the stable code
    can.
    The row identities must be readable from the **serialized**
    envelope, not only from the Go error chain:
    `CLIError.Cause` is `json:"-"`, so the assertion reads the
    identities off the JSON output. The intended carrier is a
    new `omitempty` field on `CLIError` — **not** rendered
    prose in `Detail`, which would force the consumer
    re-parsing this RDR forbids; `Detail` carries the human
    sentence, the new field carries the identities, and the
    per-identity `Count` serializes with them (A9). A9
    **confirmed** the addition is additive to RDR 0005's
    envelope; the field's clierr-local type and the same-change
    amendment of `Cause`'s now-stale doc comment are recorded on
    A9 and in Prerequisites. The kernel cannot serialize
    its own breach: frozen guards REQ-37 and REQ-28 forbid it
    importing `encoding/json` or any CLI package, so the wire
    carrier is necessarily CLI-side work.

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

Responses: `0009-escape-row-shape-conformance-ownership/artifacts/gate.md`
(Gate PASS 2026-08-11)

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
  `.Hint`, `.Cause` — `json:"-"`, doc comment "the
  wire-visible cause surface is Detail"; the type's "Extend
  with new optional fields as needed — keep them `omitempty`"
  allowance; `ExitCodeFor`'s non-`CLIError` default returns
  exit 1 but is **unreachable for a verb error**, which
  `root.go::ExecuteAndEmit` converts via
  `cobraErrorToCLIError` first — exit 2)
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
