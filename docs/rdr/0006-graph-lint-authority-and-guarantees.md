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

Source search found no implemented graph lint command under `internal/`; it
found the shipped resolution kernel (`internal/resolve/resolve.go::Resolve`),
the peer RDRs, the existing CLI output/error plumbing, and normal Go tooling
lint. The design therefore targets the RDR cluster contract rather than an
implemented command.

Prior art was read before naming approaches. `../state-machines` Seed 6 states
the lint problem directly: illegal or incomplete graphs should be caught at
design time, with the design fork of "standalone CI check vs `resolve --lint`
subcommand vs pre-commit hook" and an invariant set including dangling edges,
dead ends, determinism, guard exhaustiveness, single-valued state, and
owned-set-before-match. `../state-machines/MODEL-transition.md` gives the
closest semantic prior: design-time lint rejects overlapping predicates,
non-exhaustive predicates, and owned-tag reads before predecessor writes, and
says provenance feeds the determinism check. `../state-machines/attic/RESOLVER-DESIGN.md`
states the old thin-lint form as "Run once at design time over the table".
`../state-machines/repos/README.md` records the superseded formal-tool
conclusion: the RDR/kata graphs are small enough for a declarative table plus a
short lint instead of a model checker. The local CLI contract adds that every
user-facing graceful exit must route through `respond`/`clierr`, not direct
printing.

The peer RDR split is load-bearing. RDR 0001 owns runtime exact-one resolution
and keeps runtime refusal necessary after lint. RDR 0002 owns sparse TOML
source, normalization, stable rule identity, source spans, escape-row carriage,
and graph/render views. RDR 0003 owns symbolic predicate atoms, the tag
declaration model, the scoped row group, finite-domain exhaustiveness and
overlap semantics (including escape-row participation and the
`guard_unevaluable` narrowing), and provenance labels. RDR 0004 owns accessor
execution and read-back safety. RDR 0005 owns the resolver CLI surface and
leaves static graph lint authority to this RDR. RDR 0007 owns the kernel's
atom shape, presence decision, and refusal payload. This RDR therefore defines
when a normalized graph is accepted, which invariant failures are blocking,
the graph-level reachability relation peers quantify over, and how failures
map to the existing CLI output contract.

### Key Discoveries

- **Documented** — RDR 0001 makes runtime resolution exact-one-or-refusal and
  keeps runtime refusal distinct from design-time graph acceptance.
- **Documented** — RDR 0002's normalized candidate rows retain source rule ids
  and source locators, which lint needs for actionable diagnostics.
- **Documented** — RDR 0003 gives lint finite-domain predicate semantics for
  overlap, coverage, and owned-tag read-before-write checks, and defines the
  scoped row group lint proves over.
- **Documented** — `docs/cli-output-contract.md` makes `respond`/`clierr` the
  route for text/json success and structured failure output.
- **Documented** — prior-system notes reject a separate formal runtime or model
  checker for these small graphs and preserve a design-time validator/lint
  carve-out.
- **Documented** — the shipped kernel consults escape rows only to rescue
  `no_match` / `ambiguous_match` (`internal/resolve/resolve.go::Resolve`), so
  an escape row overlapping a guarded row is never a runtime ambiguity.
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
    for exhaustiveness claims and a blocking inability-to-prove outcome when
    proof is not possible. RDR 0003 owns the **tag declaration model** (value
    kind, finite domain, optionality, single-valued marker, element universe)
    that supplies those domains, and defines the **scoped row group** (rows
    sharing one selection context — source state and recognized outcome) the
    derivation runs over. This RDR adopts both by citation (`Technical Design`,
    `Normative Contracts`) and supplies only the reachability of selection
    contexts (`Load-Bearing Decisions`), which closes RDR 0003 A10 and A12.
  - **If wrong**: The lint would either miss ambiguous/gap cases or block valid
    guarded edges with false positives.
- **A3 Owned-tag read-before-write can be checked over the normalized graph
  without executing accessors.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0002 `Normative Contracts` require every matched or
    written tag to declare provenance (`owned`, `observed`, `recognized`) and
    preserve explicit writes/clears in normalized candidate rows. RDR 0003 A6
    verifies provenance labels are available to predicate lint, and its
    owned-tag clause rejects rows that match owned tags unless every reachable
    predecessor sets or preserves them — quantifying over the reachability
    relation this RDR defines. RDR 0004 `Technical Design` keeps accessor
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
  - **Method**: MVV Test
  - **Evidence**: Validation scenario 7 must capture `make check` or
    `.github/workflows/ci.yml` invoking the production `intrastate lint`
    command over the checked-in transition model or fixture corpus. Source
    search confirms the gate surfaces exist and names what the wiring costs:
    `Makefile::check` is the local aggregate gate (`fmt-check vet lint test`)
    but does **not** depend on `Makefile::build`, so a binary-invoking gate step
    must add that edge; `Makefile::build` builds `./bin/intrastate` from
    `cmd/intrastate`; `.github/workflows/ci.yml::jobs` runs repository gates on
    push and pull request as the discrete jobs `test`, `lint`, and `vuln`, none
    of which invokes `make check` — so the CI gate is a new step or job, not a
    free ride on an existing invocation; and `internal/cli/root.go::NewRootCmd`
    registers root-level commands through the same `ExecuteAndEmit` path.
    `clierr::ExitCodeFor` maps `GroupUserEnv` to exit 2 as scenario 2 asserts,
    but shares that code with `GroupInternal`, so the gate asserts on `Code`
    rather than the exit integer alone.
  - **If wrong**: The implementation cannot complete this RDR's MVV until the
    production gate runs the lint command. If that gate is bypassed after
    implementation, the graph may be lintable locally but not enforced at the
    design-time boundary maintainers actually rely on.
- **A6 The transition model declares an initial owned state and its terminal
  states, so reachability has a root and dead-end detection has a stop set.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence needed**: the reachability relation (`Load-Bearing Decisions`)
    is rooted at a declared initial owned state, and invariants 2 and 7 read
    declared terminals. Verified absent at Stage 4: neither `initial` nor a
    terminal declaration is normative in RDR 0002 (its only `terminal`
    occurrences name the fixture rule id `terminal-archive` in its Validation
    scenario 2), RDR 0001 and RDR 0003 declare neither, and both `[model]`
    blocks in RDR 0002's spike fixtures carry only `id`, `version`, and
    `description`. No peer declares a root because no peer needs one — RDR
    0001's kernel is stateless and single-step, so multi-step reachability
    exists nowhere in the cluster except this RDR's lint.
  - **Producer (decided at Stage 4)**: **RDR 0002's authoring schema**, not a
    sidecar. The declarations are model data about the legal graph, and RDR
    0002's premise is that every legal edge lives in one reviewable artifact; a
    sidecar would split model authoring across two files and give this RDR a
    slice of the wire-format authority the cluster split assigns to RDR 0002.
    RDR 0002 already carries this RDR as a Pending consumer for "determinism and
    reachability checks" (`0002::Capability Dependencies`). The edit is additive
    but touches two of RDR 0002's normative clauses together — the closed layout
    enumeration (root `outcomes`, `[model]`, `[tags.<tag>]`, `[accessors.<id>]`,
    `[context.<id>]`, `[[rule]]`, `[dump]`) and "`[model]` MUST contain `id` and
    `version`". RDR 0002 is `Draft` with every assumption Verified, so this
    reopens no locked document. This assumption flips to `Verified` by citation
    once that clause lands; it does not gate this RDR's remaining pre-lock
    lenses.
  - **If wrong**: reachability has no root, so invariant 6 is vacuous and
    invariant 2 cannot distinguish a designed stop from a dead end; lint would
    have to infer terminals from missing rows, which invariant 7 forbids.

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

The authoritative user-facing surface is the root command `intrastate lint`,
using the existing CLI output envelope. A pre-commit hook, a future
resolver-adjacent helper such as `intrastate flow lint`, or a resolver-local
validation flag may call it for ergonomics, but every such path must share the
same command/request builder and graph-lint engine; none defines acceptance.
Root `lint` is deliberately not under RDR 0005's `flow` group: `flow next` and
`flow resolve` answer runtime resolver questions, while `lint` is a design-time
graph acceptance gate.

Runtime refusal remains separate. Lint proves that the designed graph has no
known static defects; RDR 0001's kernel still refuses bad runtime inputs such
as missing artifacts, unavailable observed tags, unmodeled recognized outcomes,
or guard values that cannot be evaluated for a live call. Where lint's
exhaustiveness proof and the kernel's `guard_unevaluable` veto disagree, the
lint promise narrows — RDR 0003 records that narrowing and this RDR cites it.

### Technical Design

The lint pipeline has four conceptual stages. First, load and normalize the
transition model through the RDR 0002 table contract. Second, derive a graph
view from normalized candidate rows: source rule ids/spans, match patterns,
guard constraints, writes and clears, recognized outcomes, escape-row kind and
rescue classes, tag declarations, initial/terminal declarations, and the
reachable owned-state graph (`Load-Bearing Decisions`). Third, run invariant
checks over that graph and collect typed findings. Fourth, emit either a success
report or one aggregate structured `CLIError` failure through `respond`.

The lint engine boundary is an internal graph-lint package that receives a
normalized graph value, not Cobra command state and not sparse TOML. The minimum
input contract is:

- model identity and version;
- normalized candidate rows with deterministic row identity, source rule id,
  optional source span, match predicates, guard predicates (atoms in RDR 0007's
  `Key`/`Operator`/`Literal`/`Block` shape), writes, clears, and — for escape
  rows — row kind and the declared list of failure classes the row rescues
  (RDR 0002);
- declared tags with provenance (RDR 0002) and RDR 0003's tag declaration
  model — value kind, finite domain (enum value set, `{min..max}` bound,
  set-element universe), optionality marker, and single-valued marker. Lint
  reads every one of these from the declaration and never infers one from a
  tag's name, value spelling, or a fixture;
- recognized outcome alphabet, declared initial owned state, and declared
  terminal states (A6);
- accessor references and context references only as normalized identifiers
  needed for dangling-reference diagnostics.

**Row groups are RDR 0003's.** Coverage, overlap, and withholding are decided
per *scoped row group* as RDR 0003 defines it: the normalized rows sharing one
selection context — the same source state and the same recognized outcome — the
set RDR 0001 resolves exact-one over. This RDR does not define a second
grouping; it supplies which selection contexts are reachable, and invariants 3
and 4 run once per reachable group.

The mandatory invariant set is:

1. **Dangling edge** — every transition target, context reference, tag, outcome,
   accessor reference, and terminal state named by a row must resolve to a
   declared model element.
2. **Dead end** — every reachable owned-state that satisfies no declared
   terminal must be the source of at least one modeled non-escape row.
3. **Determinism / overlap** — within a scoped row group, no finite-domain
   input assignment may enable two ordinary (non-escape) rows. Escape rows are
   checked for overlap in their own populations, one per declared failure
   class, because the kernel matches them only to rescue a `no_match` /
   `ambiguous_match` refusal and only when exactly one matches; an escape row
   overlapping an ordinary row is not a runtime ambiguity and is not reported.
   Both populations and the per-class partition are RDR 0003's escape-row
   clause, cited not restated.
4. **Guard exhaustiveness / gap** — every scoped row group whose participating
   guard dimensions are all finitely declared claims closed coverage
   (default-on, never an opt-in annotation), and lint must prove
   `union(row_i accepted assignments) == scoped product` over RDR 0003's
   derivation. The scoped product includes the **presence dimension** RDR 0003's
   projection clause contributes: an `exists` atom over a key declared optional
   is a dimension of the product, and lint must not drop it — certifying a group
   exhaustive while ignoring the `{absent}` assignment is the false green RDR
   0003's narrowing forbids. A key declared always-present contributes no
   presence dimension. Escape rows contribute their accepted assignments to the
   union like any other row; there is no separate "an escape row exists"
   disjunct.
   A group whose coverage is closed by a bare escape row (one carrying no guard
   atoms) passes, but the verdict must say so (`graph-coverage-closed-by-escape`).
5. **Single-valued state** — the model must not produce a view in which a tag
   declared single-valued (RDR 0003's single-valued clause) holds two values:
   no row's writes assign it two values, and no reachable path leaves a second
   value held beside one no write or clear on that path removes. This is the
   model-level half of RDR 0003's conformance premise; the view-level check at
   assembly is RDR 0007's.
6. **Owned-set-before-match** — a row that reads an owned tag (match key or
   guard atom) must find it held in every reachable owned-state that satisfies
   the row's match pattern; the declared initial owned state counts as a write.
7. **Declared terminal/escape handling** — terminal states and escape rows are
   explicit model data; lint must not infer them from missing rows.

**Withheld exhaustiveness claims.** Lint must not certify a row group
exhaustive when any participating row — escape rows included — can refuse
`guard_unevaluable` under RDR 0007's aggregation veto; RDR 0003's narrowing
clause and its "can refuse" decision procedure (a syntactic test over the
declared optionality field, not a reachability query) govern, and this RDR
cites them. A withheld claim is emitted as `graph-unprovable-coverage` naming
the participating row and the refusing atom — the same blocking code an
unprovable dimension takes, with a widened trigger. There is no non-blocking
tier for this class.

Each finding carries a stable code, severity, model id, rule/context id when
available, source span when available, the atom (`Key`, `Operator`, `Literal`,
`Block`) when the finding is attributed to one guard atom, the failure class
when the finding is scoped to an escape population, and a concise human
message. Blocking findings make the command fail. Non-blocking findings are
informational: the advisory tier is scoped to redundant rows, unreachable
rules, and the bare-escape coverage closure above; it never changes the
success disposition and must not absorb a withheld claim.

Finding codes:

| Invariant | Stable Code | Severity |
| --- | --- | --- |
| Dangling edge/reference | `graph-dangling-edge` | blocking |
| Dead end | `graph-dead-end` | blocking |
| Determinism / overlap (either population) | `graph-overlap` | blocking |
| Guard exhaustiveness / gap over a provable product | `graph-coverage-gap` | blocking |
| Finite-domain proof unavailable, or claim withheld under RDR 0003's narrowing | `graph-unprovable-coverage` | blocking |
| Single-valued state violation | `graph-single-valued-state` | blocking |
| Owned-set-before-match | `graph-owned-before-write` | blocking |
| Declared terminal/escape handling | `graph-terminal-escape` | blocking |
| Coverage closed by a bare escape row | `graph-coverage-closed-by-escape` | info |

When one or more blocking findings exist, the command returns one aggregate
`CLIError` with `Code: graph-lint-failed` and `Group: GroupUserEnv`. The
implementation must extend `internal/cli/clierr.CLIError` with an optional
typed `Findings` field serialized as `findings`, or an equivalently named
respond/clierr-owned typed field, so JSON mode carries individual findings as
structured data. Text mode may summarize the same findings in `Detail`. On
success the `respond.OK` payload carries the same typed `findings` list holding
only non-blocking findings (possibly empty), so the bare-escape closure is
observable without inspecting the model.

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
Coverage, overlap, and withholding MUST be decided per scoped row group as RDR
0003 defines it (`0003::Technical Design`, row-group paragraph; `0003::Normative
Contracts`, participation clause). This RDR MUST NOT define a second grouping
predicate; it supplies only which selection contexts are reachable, per the
reachability relation in `Load-Bearing Decisions`.
```

```normative
Graph lint MAY claim exhaustiveness only over finite declared domains supplied
by RDR 0003's tag declaration model, and MUST read finite domain, optionality,
and single-valuedness from that declaration — never inferred from a tag's name,
value spelling, or a fixture. If a participating dimension is not finite or not
projectable, lint MUST emit `graph-unprovable-coverage` for that dimension.
The claim is default-on for every scoped row group whose participating
dimensions are all finitely declared. The scoped product MUST include the
presence dimension an `exists` atom over an optional key contributes, as RDR
0003's projection clause states (`0003::Normative Contracts`, "Lint MUST NOT
drop `exists` atoms from the product"); lint MUST NOT certify a group
exhaustive while ignoring that dimension's `{absent}` assignment.
```

```normative
Graph lint's exhaustiveness promise is narrowed exactly as RDR 0003's narrowing
clause states (`0003::Normative Contracts`, "An exhaustiveness claim MUST NOT be
stronger than the runtime it describes"): lint MUST NOT certify a row group
exhaustive when a participating row can refuse `guard_unevaluable`. A withheld
claim MUST be emitted as `graph-unprovable-coverage`, naming the participating
row and the refusing atom. This RDR cites that clause and MUST NOT restate it,
mint a second code for it, or carry it in a non-blocking tier.
```

```normative
Escape rows MUST participate in the coverage union and MUST be overlap-checked
in one population per declared failure class, never against ordinary rows, as
RDR 0003's escape-row clause states (`0003::Normative Contracts`, "A declared
escape row participates in the coverage identity"). A group whose coverage is
closed by a bare escape row MUST emit `graph-coverage-closed-by-escape` naming
that row; a bare green MUST NOT satisfy this clause.
```

```normative
Every blocking finding MUST carry a stable code, model identity, severity,
human-readable message, and the source rule/context id or source span when the
normalized model can provide one. A finding attributed to one guard atom MUST
carry that atom's `Key`, `Operator`, `Literal`, and `Block`; a finding scoped to
an escape population MUST carry the failure class.
```

```normative
Graph lint failure MUST return one aggregate `CLIError` with code
`graph-lint-failed` and `GroupUserEnv`; the individual blocking findings MUST
remain machine-readable in JSON mode through an append-only optional typed
`findings` envelope field owned by `clierr`/`respond`, not through a verb-local
wrapper or a text-only `Detail` string. Lint success MUST carry non-blocking
findings in the same typed field.
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
  fingerprint)`. The fingerprint covers the attributed atom and failure class
  when present. This makes repeated runs stable enough for review and tests
  without depending on source line numbers alone.
- **Reachability relation** — the graph lint traverses is the **owned-state
  graph**: a node is an abstract owned-state (per owned tag: absent, or held
  with its set of possible declared values; a tag with no finite domain
  abstracts to held/absent); the root is the declared initial owned state
  (A6); an edge is a normalized non-escape row whose match pattern over owned
  tags is satisfiable in the source node, producing the node with that row's
  writes applied and clears removed. Match atoms over observed/recognized tags
  and all guard atoms are **not** pruned — every such edge is taken as
  traversable. A selection context is *reachable* when some reachable
  owned-state satisfies its match pattern; a row's *reachable predecessors* are
  the rows on any root-to-source path. This is a syntactic dataflow over
  declarations, writes, and clears — no accessor execution, no runtime trace,
  no guard evaluation — so it is total, and it over-approximates the runtime
  (a path lint walks may be infeasible at runtime), which can only produce a
  false positive, never a false green. It is a different predicate from RDR
  0003's "can refuse" test, which is decided over the optionality field alone
  and never consults this relation.
- **Wire / byte format** — graph-lint failure uses aggregate
  `CLIError.Code = "graph-lint-failed"` and individual finding codes from the
  taxonomy above. JSON mode carries findings as append-only structured data in
  a typed `findings` field on both the error envelope and the success payload;
  text mode may render a concise detail summary.
- **Naming** — the canonical command and subsystem name is "lint". Rejected:
  "validate" because it is too broad and collides with parse/schema validation;
  "`resolve --lint`" as the authority because it hides graph acceptance under a
  runtime verb; and "pre-commit check" because hooks are optional ergonomics,
  not an acceptance boundary.
- **Selection / predicate** — when two ordinary rows qualify for the same
  selection context, lint rejects the model. It never selects by order; escape
  rows are modeled graph edges with their own per-class overlap population, not
  a tie-breaker and not a coverage exemption.

#### Round-Trip / Inverse Invariants

This RDR introduces no encode/decode, import/export, or inverse operation.
Parse/render fidelity remains owned by RDR 0002. Lint determinism is covered by
the finding identity decision and the Minimum Viable Validation.

#### Illustrative Code

Illustrative command shape only; RDR 0005 may still adjust flag placement:

```sh
intrastate lint --flow rdr --model ./path/to/rdr-transition-model.toml
```

Illustrative finding shapes only:

```json
{
  "code": "graph-overlap",
  "model": "rdr",
  "severity": "blocking",
  "rule": "prelock-flapping-cap",
  "message": "two rows can match status=Draft outcome=reviewed profile=large"
}
```

```json
{
  "code": "graph-unprovable-coverage",
  "model": "rdr",
  "severity": "blocking",
  "rule": "reconcile-rewind-legality",
  "atom": {"key": "rewind_target", "operator": "eq", "literal": "propose", "block": "all"},
  "message": "claim withheld: rewind_target is optional, row can refuse guard_unevaluable"
}
```

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Normalized candidate rows, source identity, escape-row kind and rescue classes | RDR 0002 | Verified peer RDR | Lint consumes normalized graph data and reports authored source ids/spans. |
| Tag declaration model, scoped row group, finite-domain and escape-row semantics, narrowing | RDR 0003 | Verified peer RDR | Enables overlap and exhaustiveness proofs; cited, not restated. |
| Owned/observed/recognized tag provenance | RDR 0002 / RDR 0003 | Verified peer RDR | Required for owned-set-before-match and coverage checks. |
| Accessor write/read-back semantics | RDR 0004 | Verified peer RDR | Lint reasons about declared owned writes without executing accessors. |
| Guard atom shape | RDR 0007 | Verified peer RDR | Atom-level finding field cites the four-field shape. |
| Initial owned state and terminal declarations | RDR 0002 (requested; arm decided Stage 4) | Pending (A6) | Roots reachability; stop set for dead-end detection. Lint-only inputs — no peer needs them, since RDR 0001's kernel is single-step. |
| CLI command exposure | RDR 0005 plus this RDR | Verified / introduced | RDR 0005 wires flow verbs; this RDR owns the lint acceptance contract and command authority. |
| CLI failure/output gateway | Existing `respond` / `clierr` | Available | Lint results must use existing text/json and exit-code behavior. |
| Graph lint invariant taxonomy and reachability relation | This RDR | Introduced | Defines blocking graph acceptance before resolver use; peers cite the relation. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Structured CLI failures | `internal/cli/clierr` | No graph-lint codes yet | Extend | Add aggregate `graph-lint-failed` plus structured finding codes; model defects map to `GroupUserEnv`. |
| Text/json output routing | `internal/cli/respond` | Existing gateway, no lint payload yet | Reuse | Success/failure output goes through `respond.OK` / `respond.Fail`. |
| Command tree | `internal/cli` | No lint command yet | Extend via RDR 0005 | Add lint command without bypassing command conventions. |
| Resolution kernel | `internal/resolve/resolve.go::Resolve` | Runtime only; consults escape rows only for `no_match` / `ambiguous_match` | Reuse as semantic reference | Lint's two-population overlap check mirrors the kernel's selection order. |
| Transition graph model | none under `internal/` beyond the kernel's `Row` | Pending peer implementation | Introduce via RDR 0002 | Lint package depends on normalized graph API, not source TOML parsing. |
| Guard reasoning | none under `internal/` | Pending peer implementation | Introduce via RDR 0003 | Lint package depends on finite-domain predicate API. |

### Decision Rationale

The blocking `intrastate lint` approach is the smallest authority that solves
the user's problem: illegal or incomplete graphs fail before use, while runtime
refusal still handles live bad inputs. It follows the local prior that the
legal graph should be specified and checked as data, not driven by a framework
or reinterpreted by a hook. It also preserves the peer seams: table source and
normalization stay in RDR 0002, predicate and row-group semantics stay in RDR
0003, accessor safety stays in RDR 0004, and CLI presentation stays in RDR 0005;
this RDR contributes the one graph property — reachability — those peers cannot
state for themselves.

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
graphs, if the model never declares a root for reachability, or if CI never
runs the lint. The recommendation survives because those are explicit
assumptions for Resolve and because each failure has a bounded response: fix
the peer normalized graph contract, restrict which exhaustiveness claims lint
may make, add the declaration to RDR 0002's schema, or wire CI to the same
command. A runtime-only or hook-only approach would hide those risks until
after the graph is already in use.

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
- **Escape row as an overlap exemption or tie-breaker**: Rejected because the
  kernel never matches escape rows against ordinary rows and refuses ambiguity
  outright; an exemption would certify green a group the runtime refuses.
- **A second finding code or a warning tier for withheld claims**: Rejected
  because a withheld claim is a proof the lint cannot make, which is exactly
  what `graph-unprovable-coverage` already names; a warning tier would let the
  guarantee be skipped where it matters.

## Trade-offs

### Consequences

- Transition model changes get a single blocking acceptance gate that can be
  run locally and in CI.
- Runtime resolver code remains simpler because graph-design defects should be
  rejected before models reach it.
- Lint must depend on peer model/predicate APIs; if those contracts drift, lint
  failure quality degrades.
- Reachability over-approximates the runtime, so a guard-infeasible path can
  produce an owned-before-write or dead-end finding the runtime would never
  hit; the cure is a clearer model, never a weaker lint.
- Some live failures remain runtime refusals by design; lint is not a promise
  that every future call has complete artifacts or observed context.

### Risks and Mitigations

- **Risk**: Lint emits correct findings without useful source locations.
  **Mitigation**: Make source rule/context identity a peer-RDR assumption and
  block lock until diagnostics can point to authored rules.
- **Risk**: Exhaustiveness checks overclaim on open domains or on rows the
  runtime can refuse.
  **Mitigation**: Require finite declared domains for closed-coverage claims,
  cite RDR 0003's narrowing, and make inability-to-prove a blocking finding.
- **Risk**: A bare escape row becomes a silent way to close any coverage gap.
  **Mitigation**: The closure is an observable informational finding naming the
  row, so review can distinguish a proof from a catch-all.
- **Risk**: CI wiring lags behind command implementation.
  **Mitigation**: Include CI-shaped command invocation in the Minimum Viable
  Validation rather than treating it as Day 2 work.

### Failure Modes

Visible failures are structured lint findings: dangling references, dead ends,
overlapping rows, coverage gaps, withheld claims, multi-valued state writes,
owned-tag read-before-write, and undeclared terminals. Silent failure would mean
accepting a model with a blocking invariant defect, letting a hook/alternate
command use different rules, or reading the same model under a different
grouping or escape-row population than RDR 0003 — the normative command/CI
authority, shared lint engine, and by-citation adoption of RDR 0003's clauses
are meant to prevent that. Diagnosis starts from the finding code plus source
rule/context id or source span, and the atom or failure class when carried.

## Implementation Plan

### Prerequisites

- [x] A5 CI gate surface verified: `Makefile` and GitHub Actions can host the
  production `intrastate lint` gate once the command and fixture corpus exist.
- [ ] A5 production gate proof pending MVV scenario 7: `make check` or the
  GitHub workflow must invoke the built `intrastate lint` command over the
  checked-in transition model or fixture corpus.
- [ ] A6: RDR 0002's schema declares the initial owned state and terminal
  states. Decided at Stage 4 in favour of the RDR 0002 arm over a sidecar; the
  clause amends RDR 0002's layout enumeration and its `[model]` contents
  clause. RDR 0002 is `Draft`, so this schedules an edit on an open peer rather
  than reopening a locked one.
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
assert `graph-dangling-edge`, `graph-dead-end`, `graph-overlap` (both an
ordinary-row pair and an escape-row pair sharing a failure class),
`graph-coverage-gap`, `graph-unprovable-coverage` (both a non-finite dimension
and a withheld claim naming row and atom), `graph-single-valued-state`,
`graph-owned-before-write`, and `graph-terminal-escape`. The legal matrix must
include a group closed by a bare escape row and assert
`graph-coverage-closed-by-escape` on success. Every blocking run must return
aggregate `CLIError.Code = graph-lint-failed` with `GroupUserEnv` exit behavior
and a machine-readable finding list.

### Phase 1: Lint Boundary

Define the lint package boundary over the normalized graph API and the finding
taxonomy without introducing a second source parser. The boundary accepts only a
normalized graph value with the minimum fields listed in Technical Design.

### Phase 2: Invariant Engine

Implement the owned-state reachability dataflow, then the mandatory graph
checks over normalized rows, finite domains, tag provenance, declared
terminals, owned writes, and RDR 0003's row-group, escape-row, and narrowing
semantics.

### Phase 3: CLI And CI Surface

Expose root `intrastate lint` through the existing Cobra/respond/clierr gateway.
If any `flow lint` alias or resolver-local validation flag is added later, it
must reuse the same request builder and graph-lint engine. Extend
`clierr.CLIError`/`respond` with the typed optional findings field used by JSON
mode on both success and failure. Add the CI-shaped production command
invocation to `make check` or the GitHub workflow once the checked-in
transition model or fixture corpus exists.

### Phase 4: Fixture Corpus

Add compact legal and illegal model fixtures that exercise every finding code
and preserve stable source identity for diagnostics.

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
overlap and coverage proof obligations including the presence dimension, A3
supplies owned-tag read-before-write checks, A4 supplies the `respond`/`clierr`
output path, A5 supplies the repository gate surface that must run the
production command, and A6 supplies the reachability root and terminal set.

The code paths these tests exercise were reviewed at Stage 4 and none exists
yet: `internal/` and `cmd/` hold eight non-test Go files whose only registered
verb is `internal/cli/version.go::newVersionCmd`, and no graph, lint verb,
invariant engine, or reachability dataflow is present. The output path the
matrix asserts against does ship — `internal/cli/clierr/clierr.go::CLIError`,
`::ExitCodeFor`, `::GroupUserEnv`, `internal/cli/respond/respond.go::OK` and
`::Fail`, and `internal/cli/root.go::ExecuteAndEmit` — and every optional
`CLIError` field is `omitempty`, so the append-only typed `findings` field this
RDR requires is wire-compatible with the shipped envelope.
`internal/cli/respond/respond.go::Success` already carries `Data any`, a host
for the success-side findings list.

1. **Scenario**: legal fixture model with all mandatory declarations and
   finite-domain coverage.
   **Expected**: `intrastate lint --as=json` succeeds through `respond.OK` and
   reports no blocking findings; `findings` carries only informational entries.
2. **Scenario**: illegal fixture models covering dangling edge, dead end,
   ordinary-row overlap, coverage gap, multi-valued state, owned-tag
   read-before-write, and undeclared terminal/escape handling.
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
8. **Scenario**: withheld claim — a fully-finite row group in which one row
   carries a value atom over a tag declared optional (in `all` in one fixture,
   in `unless` in another).
   **Expected**: `graph-unprovable-coverage` naming that row and the refusing
   atom; no `graph-coverage-gap` is emitted for the same group; overlap findings
   among the group's decidable rows are still emitted.
9. **Scenario**: escape rows — (a) an escape row whose accepted assignments
   overlap an ordinary row; (b) two escape rows declaring a shared failure
   class whose assignments overlap; (c) a group whose coverage is closed only
   by a bare escape row.
   **Expected**: (a) no `graph-overlap`; (b) `graph-overlap` naming both rows
   and the shared class, once per shared class; (c) success with
   `graph-coverage-closed-by-escape` naming the escape row.
10. **Scenario**: owned-before-write over an infeasible path — a row reads an
    owned tag that is held on every guard-feasible path but absent on one path
    lint cannot prune.
    **Expected**: `graph-owned-before-write` naming the row; the fixture
    documents that the cure is an explicit write or clear on the model, not a
    guard-aware lint.
11. **Scenario**: presence dimension — a scoped row group carrying an `exists`
    atom over a key declared optional, whose rows cover every value assignment
    but leave the `{absent}` assignment uncovered; plus a control group whose
    `exists` atom is over a key declared always-present.
    **Expected**: the optional-key group emits `graph-coverage-gap` rather than
    passing, proving the presence dimension was not dropped from the product;
    the always-present control passes, since such a key contributes no presence
    dimension.

### Performance Expectations

No throughput target is load-bearing for this RDR. Lint is a single-invocation
static check: load model data, build the owned-state graph, run deterministic
invariant checks, and render one terminal result. Determinism depends on RDR
0002's candidate-row/source identity and stable dump ordering plus this RDR's
finding identity tuple, not on source order, map order, or first-match
priority. If fixture runtime becomes material, implementation should profile
the invariant engine before changing the contract.

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

Refine resolved the escape-row contradiction (invariants 3 and 4 versus the
Load-Bearing Decision and RDR 0003) in favor of RDR 0003's two-population
reading, and aligned invariant 5 with RDR 0003's per-tag single-valued marker.
The research points to a design-time validator over the declared table, peer
RDRs own the source/predicate/runtime seams, and the solution keeps graph
acceptance in a blocking lint command while preserving runtime refusal for live
input failures.

### Assumption Verification

A1-A4 are verified. None uses `Docs Only`, and none is stamped `Verified` on
self-reference. A1-A3 are verified against peer RDR contracts, and A4 is
verified by source search against the CLI failure gateway. A2 was re-verified
at Stage 4 against RDR 0003's locked text clause by clause — the coverage and
overlap derivation, the finite-domain requirement and its blocking
inability-to-prove outcome, all five fields of the tag declaration model, the
scoped row group, the escape-row participation clause, the narrowing clause and
its syntactic "can refuse" test, and the default-on reading — each matching as
quoted. A5 is deliberately Pending by `MVV Test`: implementation must prove that
the command and fixture corpus are wired into the production gate through
Validation scenario 7; Stage 4 re-verified that the gate surfaces exist and
recorded what the wiring costs (CI runs discrete jobs, never `make check`, and
`check` does not depend on `build`). A6 is Pending on RDR 0002 declaring the
initial owned state and terminal states, with that arm decided at Stage 4 over
the sidecar alternative.

### Scope Verification

The Minimum Viable Validation is in scope: fixture-backed `intrastate lint`
invocations must prove one legal model and one illegal model per blocking
invariant through the production command path.

### Cross-Cutting Concerns

- **Versioning**: lint finding codes and JSON payload fields must be stable and
  append-only under the existing CLI output envelope.
- **Build tool compatibility**: the authoritative check is the same
  root `intrastate lint` command CI can run after `make build`; the Validation
  scenarios must capture that gate once the command and fixture corpus exist.
- **Incremental adoption**: local hooks and resolver-local validation flags may
  call the lint engine later, but only the command/CI gate defines acceptance.
- **Canonical-form / determinism**: deterministic claims are semantic finding
  identity and invariant results for the same normalized model, not
  byte-identical output or content-addressed hashes.

### Proportionality

This RDR owns one load-bearing contract: blocking static graph-lint authority
and its mandatory invariant set over the normalized model, including the
reachability relation those invariants quantify over. It does not own the
source table format, predicate grammar, row-group definition, runtime resolver,
accessor execution, or CLI output envelope. The `large` profile remains
appropriate because the contract locks graph-acceptance invariants and CI
authority with no prior accretion in Seam Lineage.

## References

- `docs/cli-output-contract.md`
- `docs/jdr/0001-resolve-kernel-seam.md` (§JD-4, §JD-13, §JD-14)
- `docs/rdr/0001-resolution-kernel.md`
- `docs/rdr/0002-transition-table-as-reviewable-data.md`
- `docs/rdr/0003-guard-predicate-exhaustiveness.md`
- `docs/rdr/0004-accessor-execution-safety-model.md`
- `docs/rdr/0005-skill-integration-cli-contract.md`
- `docs/rdr/0007-guard-predicate-totality.md`
- `internal/cli/clierr`
- `internal/cli/respond`
- `internal/resolve/resolve.go::Resolve`
