# Recommendation 0006: Graph Lint Authority And Guarantees

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-06-19
- **Status**: Draft
- **Type**: Feature
- **Profile**: large — locks one graph-lint acceptance contract: blocking authority plus invariant taxonomy.
- **Priority**: High
- **Related Issues**: None
- **Predecessors**: 0001-resolution-kernel, 0002-transition-table-as-reviewable-data, 0003-guard-predicate-exhaustiveness
- **Overrides**: None
- **Seam Lineage**: no prior accretion

## Problem Statement

A flow maintainer needs illegal or incomplete transition graphs caught at design time instead of discovered at runtime. The system-internal requirement is to define where lint runs, whether it is advisory or blocking, and which invariants it must prove before a graph is accepted.

## Context

### Background

The lint is the static twin of the kernel's runtime refusal. It must prove the legal graph before use, while the kernel rejects illegal calls at runtime; those guarantees need to stay distinct.

The real design fork is placement and authority: standalone CI check, `resolve --lint` subcommand, or pre-commit hook; advisory versus blocking; and the mandatory invariant set such as dangling edges, dead ends, determinism, guard exhaustiveness, single-valued state, and owned-set-before-match.

### Technical Environment

intrastate is a Go CLI wired through `internal/cli`. The lint must be compatible with the chosen table representation, guard predicate representation, and output contract for CLI verbs.

## Research Findings

### Investigation

Source search found no implemented transition graph, resolver, table loader,
guard package, or graph lint command under `internal/`; it found only the peer
RDR drafts, the existing CLI output/error plumbing, and normal Go tooling lint.
The design therefore targets the RDR cluster contract rather than an
implemented command.

Prior art was read before naming approaches. The resource index names
`../state-machines` as the current prior-system source, and its Seed 6 states
the lint problem directly: illegal or incomplete graphs should be caught at
design time, with the design fork of "standalone CI check vs `resolve --lint`
subcommand vs pre-commit hook" and an invariant set including dangling edges,
dead ends, determinism, guard exhaustiveness, single-valued state, and
owned-set-before-match. `../state-machines/MODEL-transition.md` gives the
closest semantic prior: design-time lint rejects overlapping predicates,
non-exhaustive predicates, and owned-tag reads before predecessor writes, and
says provenance feeds the determinism check. `../state-machines/attic/RESOLVER-DESIGN.md`
states the old thin-lint form as "Run once at design time over the table" and
lists no dangling edges, no dead ends, determinism, guard coverage, and
single-valued-state invariants. `../state-machines/repos/README.md` records the
superseded formal-tool conclusion: the RDR/kata graphs are small enough for a
declarative table plus a short lint instead of a model checker. The local CLI
contract adds that every user-facing graceful exit must route through
`respond`/`clierr`, not direct printing.

The peer RDR split is load-bearing. RDR 0001 owns runtime exact-one resolution
and keeps runtime refusal necessary after lint. RDR 0002 owns sparse TOML
source, normalization, stable rule identity, source spans, and graph/render
views. RDR 0003 owns symbolic predicate atoms, finite-domain exhaustiveness,
overlap diagnostics, and owned/observed/recognized provenance. RDR 0004 owns
accessor execution and read-back safety. RDR 0005 owns the resolver CLI surface
and may expose lint, but it explicitly leaves static graph lint authority to
this RDR. This RDR should therefore define when a normalized graph is accepted,
which invariant failures are blocking, and how those failures map to the
existing CLI output contract.

### Key Discoveries

- **Documented** — RDR 0001 makes runtime resolution exact-one-or-refusal and
  keeps runtime refusal distinct from design-time graph acceptance.
- **Documented** — RDR 0002's normalized candidate rows retain source rule ids
  and source locators, which lint needs for actionable diagnostics.
- **Documented** — RDR 0003 gives lint finite-domain predicate semantics for
  overlap, coverage, and owned-tag read-before-write checks.
- **Documented** — `docs/cli-output-contract.md` makes `respond`/`clierr` the
  route for text/json success and structured failure output.
- **Documented** — prior-system notes reject a separate formal runtime or model
  checker for these small graphs and preserve a design-time validator/lint
  carve-out.
- **Verified** — the normalized table and predicate contracts expose enough
  graph structure for lint to prove all mandatory invariants without reading
  sparse TOML directly.
- **Verified** — a blocking root command, `intrastate lint`, is the right
  authority surface; local hooks and resolver-adjacent helpers may call it but
  cannot be the source of truth.
- **Verified** — the initial invariant set can be expressed as deterministic
  checks over normalized rows, declared tags/domains, declared terminals, and
  predecessor/write reachability.

### Critical Assumptions

- **A1 Normalized rows expose every graph edge, guard constraint, write, source
  rule id, and source locator needed for lint diagnostics.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0002 `Technical Design` defines normalized candidate
    rows with source locator, predicates, and writes; RDR 0002 `Normative
    Contracts` require each candidate row to retain source rule id and source
    locator, and require tag declarations, explicit writes/clears, accessor
    references, and deterministic candidate-row dumps. RDR 0003 A5 confirms
    normalized predicates retain source identity for diagnostics.
  - **If wrong**: Lint may find a graph defect but fail to locate the authored
    rule or may need to parse sparse source through a parallel model.
- **A2 Predicate lint can decide overlap and coverage for the finite domains
  this graph claims exhaustive.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0003 A2 derives coverage as
    `union(row_i accepted assignments) == scoped product` and overlap as any
    non-empty row intersection over finite enum/boolean, declared set-universe,
    and bounded-int domains; its `Normative Contracts` require finite domains
    for exhaustiveness claims and blocking refusal/downgrade when proof is not
    possible.
  - **If wrong**: The lint would either miss ambiguous/gap cases or block valid
    guarded edges with false positives.
- **A3 Owned-tag read-before-write can be checked over the normalized graph
  without executing accessors.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0002 `Normative Contracts` require every matched or
    written tag to declare provenance (`owned`, `observed`, `recognized`) and
    preserve explicit writes/clears in normalized candidate rows. RDR 0003 A6
    verifies provenance is available to predicate lint, and its `Normative
    Contracts` reject rows that match owned tags unless reachable predecessors
    write or preserve them. RDR 0004 `Technical Design` keeps accessor
    execution outside predicate evaluation while validating write capability
    and owned-tag effects at the accessor boundary.
  - **If wrong**: Lint cannot prove that trusted owned state exists before a row
    matches it, so a class of runtime missing-state refusals remains design-time
    invisible.
- **A4 Blocking lint failures can map to the existing CLI failure envelope and
  exit-code groups.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/cli/clierr/clierr.go::CLIError` carries stable
    `Code`, `Message`, optional `Param`, `Detail`, `Hint`, and exit-code
    `Group`; `internal/cli/clierr/clierr.go::ExitCodeFor` maps existing groups
    to stable process exits; `internal/cli/respond/respond.go::Fail` emits the
    structured envelope in text/json modes; `internal/cli/root.go::ExecuteAndEmit`
    converts Cobra-level failures through the same gateway. Graph lint needs
    new stable codes, not a new envelope or direct output path.
  - **If wrong**: The lint command needs a separate output contract or new
    exit-code group before it can be authoritative in CI.
- **A5 CI can run `intrastate lint` as the blocking graph-acceptance authority
  for transition model changes.**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: Verification must add the repository gate after the command
    and fixture corpus exist: `make check` or `.github/workflows/ci.yml` must
    invoke the built `intrastate lint` command over the checked-in transition
    model or fixture corpus, and the captured run must show one legal model
    passing plus illegal fixtures for every blocking invariant class failing
    through the production command shape.
  - **If wrong**: The graph may be lintable locally but not enforced at the
    design-time boundary maintainers actually rely on.

**Method vocabulary** (pick exactly one per assumption):

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

Make graph lint a first-class, blocking acceptance gate over the normalized
transition graph. The lint consumes the normalized model produced by RDR 0002
and the symbolic predicate/domain semantics from RDR 0003; it does not parse
sparse TOML independently, execute accessors, run the resolver, or become a
runtime state-machine engine. A model is accepted only when every mandatory
static invariant passes. A model with an invariant failure is rejected before it
can be used by the resolver or accepted by CI.

The authoritative user-facing surface is the root command `intrastate lint`.
A pre-commit hook may call the same command for ergonomics, and a future
resolver-adjacent helper such as `intrastate flow lint` may exist for local
convenience, but any helper or alias must share the same command/request builder
and graph-lint engine. The authority is the blocking lint result over model
files using the existing CLI output envelope.

Runtime refusal remains separate. Lint proves that the designed graph has no
known static defects; RDR 0001's kernel still refuses bad runtime inputs such
as missing artifacts, unavailable observed tags, unmodeled recognized outcomes,
or guard values that cannot be evaluated for a live call.

### Technical Design

The lint pipeline has four conceptual stages. First, load and normalize the
transition model through the RDR 0002 table contract. Second, derive a graph
view from normalized candidate rows: source rule ids/spans, source and target
tag writes, recognized outcomes, guard constraints, terminal declarations, tag
provenance, finite domains, and owned-tag write effects. Third, run invariant
checks over that graph and collect typed findings. Fourth, emit either a success
report or one aggregate structured `CLIError` failure through `respond`.

The lint engine boundary is an internal graph-lint package that receives a
normalized graph value, not Cobra command state and not sparse TOML. The minimum
input contract is:

- model identity and version;
- normalized candidate rows with deterministic row identity, source rule id,
  optional source span, match predicates, guard predicates, writes, and clears;
- declared tags with provenance, value kind, finite-domain metadata when
  exhaustiveness is claimed, and single-valued grouping when applicable;
- recognized outcome alphabet, declared terminal states, and declared escape
  rows;
- accessor references and context references only as normalized identifiers
  needed for dangling-reference diagnostics.

The authoritative command surface is the root command `intrastate lint` over
model files using the same graph-lint package. This command is intentionally not
hidden under RDR 0005's `flow` resolver group: `flow next` and `flow resolve`
answer runtime resolver questions, while root `lint` is a design-time graph
acceptance gate. CI must run that production command shape for transition-model
changes; hooks, wrappers, aliases, and future resolver-local validation flags
may call it, but they do not define acceptance and must not use a separate rule
set.

The mandatory invariant set is:

1. **Dangling edge** — every transition target, context reference, tag, outcome,
   accessor reference, and terminal state named by a row must resolve to a
   declared model element.
2. **Dead end** — every non-terminal reachable state must have at least one
   outgoing modeled edge for a legal recognized outcome.
3. **Determinism / overlap** — no finite-domain input assignment may enable two
   candidate rows for the same state/outcome unless the model explicitly routes
   to one deterministic escape row.
4. **Guard exhaustiveness / gap** — for each state/outcome pair that claims
   closed coverage, finite-domain input assignments must either match exactly
   one modeled row or a declared escape row.
5. **Single-valued state** — writes must not produce two values for a tag class
   that the model declares single-valued, such as one lifecycle/status value.
6. **Owned-set-before-match** — every row matching an owned tag must be reachable
   only after a predecessor writes or preserves that owned tag, or after an
   initial-state declaration supplies it.
7. **Declared terminal/escape handling** — terminal states and intentional
   escapes are explicit model data; lint must not infer them from missing rows.

Each finding carries a stable code, severity, model id, rule/context id when
available, source span when available, and a concise human message. Blocking
findings make the command fail. Non-blocking advisories may exist for redundant
rows or unreachable rules only if they do not weaken the mandatory acceptance
gate.

Mandatory blocking finding codes:

| Invariant | Stable Code |
| --- | --- |
| Dangling edge/reference | `graph-dangling-edge` |
| Dead end | `graph-dead-end` |
| Determinism / overlap | `graph-overlap` |
| Guard exhaustiveness / gap | `graph-coverage-gap` |
| Required finite-domain proof unavailable | `graph-unprovable-coverage` |
| Single-valued state violation | `graph-single-valued-state` |
| Owned-set-before-match | `graph-owned-before-write` |
| Declared terminal/escape handling | `graph-terminal-escape` |

When one or more blocking findings exist, the command returns one aggregate
`CLIError` with `Code: graph-lint-failed` and `Group: GroupUserEnv`. The
implementation must extend `internal/cli/clierr.CLIError` with
an optional typed `Findings` field serialized as `findings`, or an equivalently
named respond/clierr-owned typed field, so JSON mode carries individual
findings as structured data. Each finding record carries its own stable code and
identity.
Text mode may summarize the same findings in `Detail`, but JSON mode must keep
findings machine-readable. Success means no blocking findings for the supplied
normalized model. Non-blocking advisories, if implemented, are warnings and do
not change the success disposition.

#### Normative Contracts

```normative
Graph lint MUST be a blocking acceptance gate over the normalized transition
model. A model with any blocking lint finding MUST NOT be accepted for resolver
use or CI success.
```

```normative
Graph lint MUST consume the normalized candidate-row graph from the transition
model contract. It MUST NOT define a second sparse-source parser or a parallel
transition semantics.
```

```normative
Graph lint MUST check at least these blocking invariant classes: dangling edge,
dead end, determinism/overlap, guard exhaustiveness/gap, single-valued state,
owned-set-before-match, and declared terminal/escape handling.
```

```normative
Graph lint MUST reject ambiguity instead of relying on source order,
rendered-row order, or first-match priority to choose between enabled rows.
```

```normative
Graph lint MAY claim exhaustiveness only over finite declared domains supplied
by the predicate/tag model. If a required dimension is not finite, lint MUST
emit a blocking inability-to-prove finding for any contract that depends on
closed coverage.
```

```normative
Every blocking finding MUST carry a stable code, model identity, severity,
human-readable message, and the source rule/context id or source span when the
normalized model can provide one.
```

```normative
Graph lint failure MUST return one aggregate `CLIError` with code
`graph-lint-failed` and `GroupUserEnv`; the individual blocking findings MUST
remain machine-readable in JSON mode through an append-only optional typed
`findings` envelope field owned by `clierr`/`respond`, not through a verb-local
wrapper or a text-only `Detail` string.
```

```normative
Graph lint findings MUST be emitted in deterministic order by finding identity:
model id, invariant code, source rule/context id or graph element id, then
normalized predicate/write fingerprint.
```

```normative
The authoritative CLI surface for graph acceptance MUST be the root command
`intrastate lint` or a same-engine CI invocation of that command. Pre-commit
hooks, aliases, and resolver-local validation flags MAY call that engine, but
MUST NOT define different acceptance rules.
```

```normative
Lint command success and failure MUST route through `respond.OK`,
`respond.Fail`, and `CLIError`; the command MUST NOT write directly to stdout
or stderr.
```

#### Load-Bearing Decisions

- **Identity** — a lint finding is identified by `(model id, invariant code,
  source rule/context id or graph element id, normalized predicate/write
  fingerprint)`. This makes repeated runs stable enough for review and tests
  without depending on source line numbers alone.
- **Wire / byte format** — graph-lint failure uses aggregate
  `CLIError.Code = "graph-lint-failed"` and individual finding codes from the
  mandatory taxonomy above. JSON mode carries findings as append-only structured
  data in a typed `findings` field on the error envelope; text mode may render a
  concise detail summary.
- **Naming** — the canonical command and subsystem name is "lint". Rejected:
  "validate" because it is too broad and collides with parse/schema validation;
  "`resolve --lint`" as the authority because it hides graph acceptance under a
  runtime verb; and "pre-commit check" because hooks are optional ergonomics,
  not an acceptance boundary.
- **Selection / predicate** — when two rows qualify for the same state/outcome,
  lint rejects the model. It never selects by order; explicit escape rows are
  modeled graph edges, not a tie-breaker.

#### Round-Trip / Inverse Invariants

This RDR introduces no encode/decode, import/export, or inverse operation.
Parse/render fidelity remains owned by RDR 0002. Lint determinism is covered by
the finding identity decision and the Minimum Viable Validation.

#### Illustrative Code

Illustrative command shape only; RDR 0005 may still adjust flag placement:

```sh
intrastate lint --flow rdr --model ./path/to/rdr-transition-model.toml
```

Illustrative finding shape only:

```json
{
  "code": "graph-overlap",
  "model": "rdr",
  "severity": "blocking",
  "rule": "prelock-flapping-cap",
  "message": "two rows can match status=Draft outcome=reviewed profile=large"
}
```

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Normalized candidate rows and source identity | RDR 0002 | Verified peer RDR | Lint consumes normalized graph data and reports authored source ids/spans. |
| Symbolic predicate finite-domain semantics | RDR 0003 | Verified peer RDR | Enables overlap and exhaustiveness proofs. |
| Owned/observed/recognized tag provenance | RDR 0002 / RDR 0003 | Verified peer RDR | Required for owned-set-before-match and coverage checks. |
| Accessor write/read-back semantics | RDR 0004 | Verified peer RDR | Lint reasons about declared owned writes without executing accessors. |
| CLI command exposure | RDR 0005 plus this RDR | Verified / introduced | RDR 0005 wires flow verbs; this RDR owns the lint acceptance contract and command authority. |
| CLI failure/output gateway | Existing `respond` / `clierr` | Available | Lint results must use existing text/json and exit-code behavior. |
| Graph lint invariant taxonomy | This RDR | Introduced | Defines blocking graph acceptance before resolver use. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Structured CLI failures | `internal/cli/clierr` | No graph-lint codes yet | Extend | Add aggregate `graph-lint-failed` plus structured finding codes; model defects map to `GroupUserEnv`. |
| Text/json output routing | `internal/cli/respond` | Existing gateway, no lint payload yet | Reuse | Success/failure output goes through `respond.OK` / `respond.Fail`. |
| Command tree | `internal/cli` | No resolver/lint commands yet | Extend via RDR 0005 | Add lint command without bypassing command conventions. |
| Transition graph model | none under `internal/` | Pending peer implementation | Introduce via RDR 0002 | Lint package depends on normalized graph API, not source TOML parsing. |
| Guard reasoning | none under `internal/` | Pending peer implementation | Introduce via RDR 0003 | Lint package depends on finite-domain predicate API. |

### Decision Rationale

The blocking `intrastate lint` approach is the smallest authority that solves
the user's problem: illegal or incomplete graphs fail before use, while runtime
refusal still handles live bad inputs. It follows the local prior that the
legal graph should be specified and checked as data, not driven by a framework
or reinterpreted by a hook. It also preserves the peer seams: table source and
normalization stay in RDR 0002, predicate semantics stay in RDR 0003, accessor
safety stays in RDR 0004, and CLI presentation stays in RDR 0005.

Q-O-C matrix for the large-profile decision:

| Criterion | Standalone blocking `intrastate lint` + CI | `resolve --lint` validation mode | Pre-commit hook authority | External model checker / FSM validator |
| --- | --- | --- | --- | --- |
| Correctness fit | Directly gates graph acceptance before resolver use. | Couples design-time proof to runtime verb and can be skipped by non-resolve paths. | Local-only; hooks are easy to skip or absent in CI. | Strong proof surface, but oversized for small declared tables. |
| Prior-art alignment | Matches sibling "table + thin CLI + lint" and "run once at design time" notes. | Partly matches thin CLI, but weakens the design-time/runtime split. | Useful ergonomics, not an authority in prior art. | Prior notes keep validator carve-out but supersede model-checker adoption. |
| Reversibility | Can add hook wrappers or resolver-local flags later without changing lint semantics. | Harder to extract once runtime and lint flags share behavior. | Easy to add/remove, but cannot be the only gate. | Hard to remove after specs/tests depend on tool language. |
| Blast radius | New lint package/verb plus CI path; no runtime engine expansion. | Runtime CLI grows acceptance semantics and may obscure refusal boundaries. | Low code blast radius but high process risk. | High dependency/tooling blast radius. |
| Cost | Implementable over normalized rows and finite-domain predicates. | Similar engine cost plus confusing command semantics. | Cheap wrapper but insufficient guarantee. | Expensive setup and translation with little benefit for tiny graphs. |

Premortem: this approach fails if normalized rows cannot preserve source
identity, if finite-domain predicate reasoning is too weak for the RDR/kata
graphs, or if CI never runs the lint. The recommendation survives because those
are explicit assumptions for Resolve and because each failure has a bounded
response: fix the peer normalized graph contract, restrict which exhaustiveness
claims lint may make, or wire CI to the same command. A runtime-only or hook-only
approach would hide those risks until after the graph is already in use.

## Alternatives Considered

### Alternative 1: Standalone Blocking `intrastate lint` And CI

**Description**: Add a lint engine and command that consumes normalized
transition models and fails on blocking graph invariant violations. CI invokes
the same command for model changes.

**Pros**:

- Clear design-time authority.
- Keeps runtime resolver refusal separate from graph acceptance.
- Lets hooks and future validation flags reuse one engine.
- Fits existing CLI output and error conventions.

**Cons**:

- Requires CI wiring and stable fixture coverage before the guarantee is real.
- Depends on peer normalized-row and predicate contracts being strong enough.

**Reason for selection**: It is the only option that is both enforceable and
properly scoped to design-time graph acceptance.

### Alternative 2: `resolve --lint` As The Authority

**Description**: Put graph lint behind the resolver command as a validation mode
or flag.

**Pros**:

- Keeps graph validation close to the command that consumes the graph.
- May be convenient for users already invoking resolver flows.

**Cons**:

- Blurs the static proof/runtime refusal boundary.
- Makes lint look optional or call-specific instead of a model acceptance gate.
- Risks forcing the resolver CLI to carry graph-diagnostic concerns owned here.

**Reason for rejection**: It is useful as an optional wrapper later, but it is
the wrong authority surface.

### Alternative 3: Pre-Commit Hook Authority

**Description**: Run graph lint only as a repository hook before commits.

**Pros**:

- Fast local feedback.
- Low command-surface complexity if it shells out to an internal package.

**Cons**:

- Hooks are local, mutable, and easy to bypass.
- CI and automated generation paths can miss the gate.
- Hook output tends to drift from the project CLI output contract.

**Reason for rejection**: A hook may call `intrastate lint`, but it cannot be
the acceptance authority.

### Alternative 4: External Model Checker Or FSM Validator

**Description**: Translate the legal graph to Quint/TLA+/SCXML/Sismic or a
similar validator and use that tool as the proof authority.

**Pros**:

- Strong prior-art fit for larger safety-critical graphs.
- Could express richer temporal properties later.

**Cons**:

- Adds a second graph language and toolchain.
- Duplicates the normalized model semantics this RDR needs to check directly.
- Prior notes already judged it overkill for these small RDR/kata graphs.

**Reason for rejection**: Keep external validators as future contrast or audit
tools; the authoritative gate should be native lint over the model intrastate
actually consumes.

### Briefly Rejected

- **Advisory-only lint**: Rejected because it does not catch illegal graphs
  before use.
- **Runtime-only resolver refusal**: Rejected because it discovers design
  defects one call at a time instead of rejecting the graph.
- **Generated exhaustive table as the authority**: Rejected because RDR 0002
  makes the sparse source the authored model and the expanded table a view.

## Trade-offs

### Consequences

- Transition model changes get a single blocking acceptance gate that can be
  run locally and in CI.
- Runtime resolver code remains simpler because graph-design defects should be
  rejected before models reach it.
- Lint must depend on peer model/predicate APIs; if those contracts drift, lint
  failure quality degrades.
- Some live failures remain runtime refusals by design; lint is not a promise
  that every future call has complete artifacts or observed context.

### Risks and Mitigations

- **Risk**: Lint emits correct findings without useful source locations.
  **Mitigation**: Make source rule/context identity a peer-RDR assumption and
  block lock until diagnostics can point to authored rules.
- **Risk**: Exhaustiveness checks overclaim on open domains.
  **Mitigation**: Require finite declared domains for closed-coverage claims and
  make inability-to-prove a blocking finding when acceptance depends on it.
- **Risk**: CI wiring lags behind command implementation.
  **Mitigation**: Include CI-shaped command invocation in the Minimum Viable
  Validation rather than treating it as Day 2 work.

### Failure Modes

Visible failures are structured lint failures: dangling references, dead ends,
overlapping rows, coverage gaps, multi-valued state writes, owned-tag
read-before-write, undeclared terminals, and inability to prove a required
finite-domain guarantee. Silent failure would mean accepting a model with a
blocking invariant defect or letting a hook/alternate command use different
rules; the normative command/CI authority and shared lint engine are meant to
prevent that. Diagnosis starts from the finding code plus source rule/context id
or source span.

## Implementation Plan

### Prerequisites

- [ ] A5 CI gate verification pending: once `intrastate lint` and the fixture
  corpus exist, prove the repository gate invokes the production command over
  the checked-in transition model or lint fixture corpus.
- [ ] RDR 0002 normalized-row identity and RDR 0003 finite-domain predicate
  semantics are coherent enough to implement checks against.
- [x] RDR 0005 command placement is coherent enough to expose root
  `intrastate lint`: RDR 0005 owns runtime `flow` verbs and leaves graph lint
  authority to this RDR.

### Minimum Viable Validation

Add a fixture-backed lint invocation that uses the same command shape intended
for CI. It must pass one legal transition model and fail one illegal model for
each blocking invariant class named in this RDR, asserting stable finding codes
and source rule/context identity in JSON mode. The illegal fixture matrix must
assert `graph-dangling-edge`, `graph-dead-end`, `graph-overlap`,
`graph-coverage-gap`, `graph-unprovable-coverage`,
`graph-single-valued-state`, `graph-owned-before-write`, and
`graph-terminal-escape`. Every blocking run must return aggregate
`CLIError.Code = graph-lint-failed` with `GroupUserEnv` exit behavior and a
machine-readable finding list.

### Phase 1: Lint Boundary

Define the lint package boundary over the normalized graph API and the finding
taxonomy without introducing a second source parser. The boundary accepts only a
normalized graph value with the minimum fields listed in Technical Design.

### Phase 2: Invariant Engine

Implement the mandatory graph checks over normalized rows, finite domains, tag
provenance, declared terminals, and owned writes.

### Phase 3: CLI And CI Surface

Expose root `intrastate lint` through the existing Cobra/respond/clierr gateway.
If any `flow lint` alias or resolver-local validation flag is added later, it
must reuse the same request builder and graph-lint engine. Extend
`clierr.CLIError`/`respond` with the typed optional findings field used by JSON
mode. Add the CI-shaped production command invocation to `make check` or the
GitHub workflow once the checked-in transition model or fixture corpus exists.

### Phase 4: Fixture Corpus

Add compact legal and illegal model fixtures that exercise every blocking
finding and preserve stable source identity for diagnostics.

### Day 2 Operations

| Resource | List | Info | Delete | Verify | Backup |
| --- | --- | --- | --- | --- | --- |
| Transition model files | Covered by repository tools | Covered by lint output and table dumps | Covered by version control | In scope through `intrastate lint` | Covered by version control |
| Lint fixtures | Covered by repository tools | Covered by test names and fixture paths | Covered by version control | In scope through tests | Covered by version control |

### New Dependencies

No new third-party dependency is selected at Propose. The chosen approach should
first use the normalized model and predicate packages introduced by peer RDRs.

## Validation

### Testing Strategy

The verified assumptions imply a production-command test matrix: exercise
`intrastate lint` through the Cobra path and the same output gateway intended
for CI. Coverage must include success output, blocking failure output, stable
finding codes, and source rule/context identity in JSON mode. A1 supplies the
source identity and write/predicate graph data, A2 supplies finite-domain
overlap and coverage proof obligations, A3 supplies owned-tag
read-before-write checks, and A4 supplies the `respond`/`clierr` output path.

1. **Scenario**: legal fixture model with all mandatory declarations and
   finite-domain coverage.
   **Expected**: `intrastate lint --as=json` succeeds through `respond.OK` and
   reports no blocking findings. Non-blocking advisories, if present, are
   warnings and do not change exit behavior.
2. **Scenario**: illegal fixture models covering dangling edge, dead end,
   overlap, coverage gap, multi-valued state, owned-tag read-before-write, and
   undeclared terminal/escape handling.
   **Expected**: each run fails through `respond.Fail` with aggregate
   `CLIError.Code = graph-lint-failed`, exit code 2, and a machine-readable
   finding containing the exact stable code for that invariant, blocking
   severity, model id, and source rule/context id or span when the normalized
   model provides one.
3. **Scenario**: a model claims closed coverage over an input dimension that is
   not finite under the predicate contract.
   **Expected**: lint emits `graph-unprovable-coverage` instead of silently
   accepting or weakening the coverage guarantee.
4. **Scenario**: text and JSON modes for the same legal and illegal fixtures.
   **Expected**: both modes return the same semantic result and exit behavior
   without direct stdout/stderr writes from the command. JSON finding order is
   deterministic by finding identity; text tests may assert equivalent
   semantics without depending on paragraph wrapping.
5. **Scenario**: a finding whose normalized row lacks a rule id but has a source
   span or graph element id.
   **Expected**: the finding remains actionable and stable by carrying the
   available source span or graph element id in the identity fields.
6. **Scenario**: command placement.
   **Expected**: root `intrastate lint` is the canonical tested command path.
   Any alias or resolver-local helper uses the same request builder and
   graph-lint engine, and tests fail if it accepts different flags or returns a
   different semantic result.
7. **Scenario**: CI-shaped invocation.
   **Expected**: `make check` or `.github/workflows/ci.yml` invokes the
   production `intrastate lint` command over the checked-in transition model or
   fixture corpus, not a hook-only wrapper, unit-test-only engine path, or
   alternate rule implementation.

### Performance Expectations

No throughput target is load-bearing for this RDR. Evidence bounds lint to a
single-invocation static check over the normalized model: load model data,
derive the graph view, run deterministic invariant checks over declared finite
domains and predecessor/write reachability, and render one terminal result.
Determinism depends on RDR 0002's candidate-row/source identity and stable dump
ordering plus this RDR's finding identity tuple, not on source order, map order,
or first-match priority. If fixture runtime becomes material, implementation
should profile the invariant engine before changing the contract.

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

Pending Stage 7. Refine found no internal contradiction between the research and
the proposed solution: the research points to a design-time validator over the
declared table, peer RDRs own the source/predicate/runtime seams, and the
solution keeps graph acceptance in a blocking lint command.

### Assumption Verification

A1-A4 are verified. None uses `Docs Only`, and none is stamped `Verified` on
self-reference. A1-A3 are verified against peer RDR contracts, and A4 is
verified by source search against the CLI failure gateway. A5 is pending
because CI authority cannot be verified until the command and checked-in model
or fixture corpus exist; its verification plan is the repository gate spike
named in A5 and the Validation scenarios.

### Scope Verification

The Minimum Viable Validation is in scope: fixture-backed `intrastate lint`
invocations must prove one legal model and one illegal model per blocking
invariant through the production command path.

### Cross-Cutting Concerns

- **Versioning**: lint finding codes and JSON payload fields must be stable and
  append-only under the existing CLI output envelope.
- **Build tool compatibility**: the authoritative check is the same
  root `intrastate lint` command CI can run after `make build`; A5 remains
  pending until that gate is wired and captured.
- **Incremental adoption**: local hooks and resolver-local validation flags may
  call the lint engine later, but only the command/CI gate defines acceptance.
- **Canonical-form / determinism**: deterministic claims are semantic finding
  identity and invariant results for the same normalized model, not
  byte-identical output or content-addressed hashes.

### Proportionality

This RDR owns one load-bearing contract: blocking static graph-lint authority
and its mandatory invariant set over the normalized model. It does not own the
source table format, predicate grammar, runtime resolver, accessor execution, or
CLI output envelope. The `large` profile remains appropriate because the
contract locks graph-acceptance invariants and CI authority with no prior
accretion in Seam Lineage.

## References

- `docs/cli-output-contract.md`
- `docs/rdr/0001-resolution-kernel.md`
- `docs/rdr/0002-transition-table-as-reviewable-data.md`
- `docs/rdr/0003-guard-predicate-exhaustiveness.md`
- `docs/rdr/0004-accessor-execution-safety-model.md`
- `docs/rdr/0005-skill-integration-cli-contract.md`
- `internal/cli/clierr`
- `internal/cli/respond`
