# Recommendation 0006: Graph Lint Authority And Guarantees

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-06-19
- **Status**: Implemented
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

## Critical Assumptions

- **A1 Normalized rows expose every graph edge, guard constraint, write, source
  rule id, and source locator needed for lint diagnostics.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0002 `Technical Design` defines normalized candidate
    rows with source locator, predicates, and writes; RDR 0002 `Normative
    Contracts` require each candidate row to retain source rule id and source
    locator, and require tag declarations, explicit writes/clears, accessor
    references, and deterministic candidate-row dumps. RDR 0003:A5 confirms
    normalized predicates retain source identity for diagnostics.
  - **If wrong**: Lint may find a graph defect but fail to locate the authored
    rule or may need to parse sparse source through a parallel model.
- **A2 Predicate lint can decide overlap and coverage for the finite domains
  this graph claims exhaustive.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0003:A2 derives coverage as
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
    contexts (`Load-Bearing Decisions`) — group *membership* is the authored
    match pattern, per RDR 0003, and reachability decides only which groups are
    proven. That is the confirmation `0003::A10` awaits and the published
    relation `0003::A12` awaits; the cluster gate (`0003-0006-0007`, F4) scoped
    A10's discharge to exactly this citation, "not a redefinition".
  - **Peer records this draft supplies a producer for** (they close on the peer,
    not here — a Draft cannot flip a locked peer's record): **A17** by invariant
    3's two-population overlap, **A18** by invariant 5 plus
    `graph-always-present-owned` for the owned half only (A7), **A19** by
    `graph-coverage-closed-by-escape`, **A20**'s lint half by the atom field on
    the finding contract. RDR 0003 is `Final [locked 2026-08-22]`, so each is a
    route-back for Stage 6 to route, tracked on A7 and the needs-verification
    list — not closed by assertion from this side.
  - **If wrong**: The lint would either miss ambiguous/gap cases or block valid
    guarded edges with false positives.
- **A3 Owned-tag read-before-write can be checked over the normalized graph
  without executing accessors.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0002 `Normative Contracts` require every matched or
    written tag to declare provenance (`owned`, `observed`, `recognized`) and
    preserve explicit writes/clears in normalized candidate rows. RDR 0003:A6
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
  - **Gate target (decided here, so the scenario has one oracle)**: a **new
    `graph-lint` job in `.github/workflows/ci.yml`** — the key `lint` is taken
    by the shipped golangci-lint/gofmt job (`ci.yml::jobs.lint`, `name: Lint`)
    and MUST NOT be reused or clobbered — running `make build` then the built
    `./bin/intrastate lint --as=json` over the **checked-in transition model**.
    Not the fixture corpus, which tests the engine rather than the shipped
    model. `make check` is additionally wired to the same command for local
    parity, which requires adding the `build` edge `check` currently lacks. The
    gate asserts on the JSON `code` field, never the exit integer alone, because
    `ExitCodeFor` maps `GroupUserEnv` and `GroupInternal` both to 2. The
    disjunctions the earlier wording carried ("`make check` **or** CI", "model
    **or** corpus") are resolved: CI is the authority, the model is the subject,
    and local `make check` is convenience.
  - **Subject the gate runs over (A8)**: no transition model is checked in today
    — the repo's only `.toml` files are `.roborev.toml`, `.kata.toml`, and RDR
    spike fixtures, which this gate excludes by name. Authoring and homing that
    model is booked as **A8**; until it lands, scenario 7 is uninstantiable and
    A5 cannot flip.
  - **Evidence**: Validation scenario 7 must capture that job invoking the
    production `intrastate lint` command over the checked-in transition
    model. Source
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
  - **Stage 6 disposition — DOWNGRADED.** Not verifiable before lock: the gate
    proves a command that does not exist yet, over a model that does not exist
    yet (A8). Named plan that WILL run during implementation: **Validation
    scenario 7**, whose oracle pre-lock made mechanical — the `graph-lint` job
    asserting the JSON `code`. Survivable because the MVV is fixture-backed and
    does not consume this gate: every blocking invariant class is proven against
    authored fixtures through the production command path, so the contract this
    RDR locks is provable without the CI wiring. What defers is *enforcement*,
    not the proof. Not MVV-critical, so HARD RULE 2 does not fire.
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
    0001's kernel is normatively stateless and returns one legal transition
    plan per supplied snapshot — single-step is this RDR's gloss on that, not
    an RDR 0001 clause — so multi-step reachability
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
  - **Shape requested (this RDR's consumer requirement)**: both declarations are
    **tag predicates, not state names** — the cluster models state as a tag-set
    (RDR 0001) and RDR 0002 emits next-state tags, so there is nothing to name.
    Concretely: an `initial` table of owned `tag = value` assignments fixing the
    root node, and a `terminal` list of predicates over owned tags in the same
    atom shape rules already use. RDR 0002 owns the final spelling; what this
    RDR requires is that neither is a bare identifier a row references by name,
    since invariant 1 checks them as tag keys and values and invariant 2
    evaluates a terminal as a predicate over a node.
  - **Behavior before it lands**: not "vacuous" — a model declaring no initial
    owned state is rejected with a blocking finding (disposition table), so the
    unlanded schema cannot be mistaken for a clean model.
  - **If wrong**: reachability has no root, so invariant 6 is vacuous and
    invariant 2 cannot distinguish a designed stop from a dead end; lint would
    have to infer terminals from missing rows, which invariant 7 forbids.
  - **Stage 6 disposition — DOWNGRADED.** Confirmed at source rather than
    assumed: neither `initial` nor a terminal declaration is normative in RDR
    0002, RDR 0001 and RDR 0003 declare neither, and both `[model]` blocks in
    RDR 0002's spike fixtures carry only `id`, `version`, `description`. The
    absence is real and this RDR states it inline. Named plan: the additive
    clause on **RDR 0002's authoring schema**, which is `Draft` — a scheduled
    edit on an open peer, not a route-back. Survivable because the window is
    not silent: a model declaring no initial owned state takes a blocking
    `graph-dangling-edge` (disposition table), so an unlanded schema cannot read
    as a clean model. The MVV's fixtures declare their own roots and terminals,
    so the invariants are provable before the peer clause lands.
- **A7 Invariant 5 plus `graph-always-present-owned` discharge the owned half of
  RDR 0003:A18's conformance premise; the observed/recognized half stays
  unowned.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence needed**: RDR 0003's conformance premise routes the view-level
    check to "the kernel's view assembly", but its A18 (`Status: Pending`)
    records the premise as "currently held by no component" — RDR 0007 owns the
    kernel's view handling and states no conformance obligation, and
    `internal/resolve/resolve.go::assemble` merges owned, observed, and
    recognized tags into a `TagSet` while reading no declaration (`conform`
    occurs nowhere in the non-test kernel). A18 names two sufficient producers
    and requires only one: RDR 0007's assembly rejecting a non-conforming view,
    **or** this RDR's lint reporting a model-level conformance violation.
    Invariant 5 is that second producer for the **single-valued** conjunct, and
    `graph-always-present-owned` for the **always-present** conjunct over owned
    keys only. RDR 0003 defines a conforming view as *every always-present key
    present* **and** *every single-valued tag holding at most one declared
    value*; an always-present observed or recognized key is not model-decidable,
    since those arrive from accessors and the caller and `assemble` reads no
    declaration. The discharge is therefore partial by construction, not by
    omission. It flips to `Verified` when RDR 0003:A18 records this RDR's two
    codes as its model-level producer and records the observed/recognized
    residue as still open. RDR 0003 is `Final [locked 2026-08-22]`, so that edit
    is a route-back on a locked peer, not a Draft amendment — booked here rather
    than assumed.
  - **If wrong**: if RDR 0007 later adopts a view-level check with different
    semantics, the two producers could disagree on which views are conforming;
    the model-level finding stays correct, but the pair needs reconciling. If
    the observed/recognized residue is claimed by no component, a non-conforming
    view can still reach the kernel — lint's green means the *model* conforms,
    never that every assembled view will.
  - **Stage 6 disposition — DOWNGRADED.** The peer's own record settles the
    routing: RDR 0003:A18 is `Status: Pending` and was itself **DOWNGRADED at
    RDR 0003's Stage 6** as *not lock-blocking*, with its venue named as "RDR
    0007's next touch ... and/or RDR 0006's refine — either discharges it, both
    is better". RDR 0003's own `Status` line lists A18 among eight open records,
    "none lock-blocking". This RDR has now supplied that second producer, so the
    edit is a **scheduled peer confirmation**, not a defect either document must
    resolve before locking. Survivable because this RDR never asserts the
    view-level enforcement: invariant 5 states the discharge as *partial* and
    names the observed/recognized residue as unowned, so no clause here claims
    an enforcement it cannot cite. Route-back on a locked peer, tracked for RDR
    0003's next touch alongside A10/A12/A17/A19/A20.
- **A8 A transition model is checked into this repo for the CI gate to lint.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence needed**: A5's gate target names the "checked-in transition
    model" as its subject and excludes the fixture corpus. Verified absent:
    the repo's only `.toml` files are `.roborev.toml`, `.kata.toml`, and the
    spike fixtures under `docs/rdr/*/evidence/spikes/`. No RDR in the cluster
    books who authors the shipped model or where it lives — RDR 0002 owns its
    *format*, not its instantiation. Flips when a model file exists at a named
    path and `intrastate lint` resolves it through RDR 0005's discovery.
  - **If wrong**: scenario 7 has no subject, so A5 cannot flip and the gate
    either ships pointing at a spike fixture (undoing A5's decided disjunction)
    or does not ship at all.
  - **Stage 6 disposition — DOWNGRADED.** Re-verified at this stage against the
    working tree: the repo's complete `.toml` set is `.roborev.toml`,
    `.kata.toml`, and three RDR spike fixtures under
    `docs/rdr/000{2,3}-*/evidence/spikes/`. No transition model is checked in,
    exactly as recorded. Named plan: the model is authored and homed during
    implementation, as **Prerequisites** now books, and scenario 7 instantiates
    against it. Survivable because it gates only scenarios 7 and 23 — the CI
    gate and the false-positive census — while the MVV proper is fixture-backed
    and needs no shipped model. Not MVV-critical, so HARD RULE 2 does not fire.
- **A9 `--flow <id>` can resolve a model through RDR 0005's config discovery.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence needed**: `Illustrative Code` makes `--flow` the ergonomic entry
    point, but RDR 0005's own audit defers transition-model config discovery
    ("Parser placeholder; no transition-model config yet"). Until RDR 0005 lands
    discovery, only `--model <path>` is instantiable.
  - **If wrong**: the documented invocation does not work; the gate and the MVV
    must use `--model <path>` explicitly, which is sufficient but less
    ergonomic.
  - **Stage 6 disposition — DOWNGRADED.** Confirmed at source: RDR 0005 is
    `Final` and its own audit row reads "Parser placeholder; no transition-model
    config yet. | Extend later". Named plan: `--flow` flips to instantiable when
    RDR 0005's config discovery lands; until then `--model <path>` is the form
    the gate and the MVV use, which **Prerequisites** now records. Survivable
    because the RDR's own `If wrong` is the shipping behavior, not a failure —
    `--model <path>` is sufficient for every scenario, and the two flags are
    already pinned mutually exclusive. Ergonomics, not authority.
- **A10 A write to a single-valued tag replaces its prior value, so no path
  accumulates a second value.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence needed**: invariant 5 decides single-valuedness **per row,
    syntactically**, and drops the path-accumulation case on this semantics. No
    peer states it: RDR 0002's write/clear clauses fix only that "absence from
    both the write block and the clear list MUST NOT imply deletion" (no
    replace-vs-accumulate rule), and RDR 0003 owns the single-valued marker, not
    the update rule. Flips when RDR 0002 states write-replaces for single-valued
    tags. RDR 0002 is `Draft`, so this is a scheduled edit on an open peer, not
    a route-back.
  - **If wrong**: under accumulate or set-union semantics a single path can
    reach two values for a single-valued tag, invariant 5's per-row reading
    misses it, and the multi-value state reaches the kernel where
    `assemble` resolves it silently — lint reported green.
  - **Stage 6 disposition — DOWNGRADED.** Confirmed at source rather than
    assumed: RDR 0002's write/clear clause fixes only that "Absence from both
    the write block and the clear list MUST NOT imply deletion" — no
    replace-vs-accumulate rule is stated anywhere in it — and RDR 0003 owns the
    single-valued *marker*, not the update rule. Named plan: RDR 0002 states
    write-replaces for single-valued tags; RDR 0002 is `Draft`, so this is a
    scheduled edit on an open peer, now booked in **Prerequisites**. Survivable
    because the failure is bounded and named in invariant 5's own text: if a
    peer lands accumulate or set-union semantics instead, the per-row reading is
    unsound and the path-accumulation case returns — a re-widening of one
    invariant, not a redesign of the contract. The MVV's
    `graph-single-valued-state` fixture proves the per-row check regardless;
    A10 governs whether that check is *complete*, not whether it runs, so it is
    not MVV-critical and HARD RULE 2 does not fire.

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
can be used by the resolver or accepted by CI. "Before use" is enforced **at the
merge boundary, not at resolve time**: the CI gate is what stops an illegal model
from landing. The kernel does not consult a lint verdict at runtime — RDR 0001's
kernel is stateless and refuses bad inputs on its own — so a locally-edited
illegal model can still be resolved against on a working tree. That is accepted:
coupling the resolver to a lint result would put a design-time proof on the
runtime path, which this RDR's whole authority split rejects.

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
  model — value kind, finite domain (enum value set or `{min..max}` bound),
  optionality marker, single-valued marker, and the set-element universe a
  set-valued kind declares. Lint
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
and 4 run once per reachable group. This RDR reads that division of labour as
RDR 0003 states it, which closes `0003::A10`.

**Source state is the authored match pattern.** The cluster has no named
states — RDR 0001 models state as a tag-set and RDR 0002 emits *next-state tags*
rather than a transition target — so RDR 0003's "same source state" needs an
identity. It is the row's **authored match pattern**, the field
`internal/resolve/resolve.go::Resolve` compares (`view.matches(row.Match)`)
before the guard gate: rows sharing a match pattern and a recognized outcome are
one group. Two rows whose patterns differ are two groups even when some abstract
node satisfies both, because no concrete view makes them candidates together.

The reachability relation below decides **which** groups are proven, never
**which rows are in one**. Membership is syntactic and per-row; reachability is
a filter over contexts. Keying groups on abstract nodes instead would pool rows
the runtime never evaluates together — see the coverage direction in
`Load-Bearing Decisions`.

The mandatory invariant set is:

1. **Dangling edge** — every context reference, tag key, tag value, outcome, and
   accessor reference named by a row must resolve to a declared model element,
   and the model must declare an initial owned state. There is no "transition
   target" term: RDR 0002 emits *next-state tags*, not a named target, so the
   edge's destination is checked as tag keys and values against their
   `[tags.<tag>]` declarations and declared domains. Terminal declarations are
   checked the same way — as tag predicates, not as state names.
2. **Dead end** — every reachable owned-state node that satisfies no declared
   terminal must be the source of at least one modeled non-escape row (escape
   self-loops do not count as progress). A terminal is a predicate over owned
   tags, and a node *satisfies* it when **every** value in each of the node's
   per-tag value sets meets it. Partial satisfaction is deliberately not enough:
   accepting it would let a merged node close on a path that has not actually
   terminated.
   **The test runs on split nodes, not merged ones.** Universal satisfaction is
   anti-monotone in merging — the more paths reach a node, the less likely any
   single terminal covers its whole value set — so testing it against a merged
   node would accuse every flow with two routes to completion. Before evaluating
   this invariant, lint splits each reachable node on the terminal-participating
   keys: a node whose `status` set is `{done, review}` against a terminal
   `status=done` is tested as `{done}` (satisfies the terminal, no outgoing row
   required) and `{review}` (does not, so it must have one) — independently.
   Splitting is exact, not an approximation: each split node is a set of
   concrete views the merged node already stood for, so nothing reachable is
   dropped and no false green is introduced. The split is bounded by the
   declared domains of terminal-participating keys only, never the whole
   lattice.
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
   declared single-valued (RDR 0003's single-valued clause) holds two values.
   This is decided **per row, syntactically**: no row's write block may assign a
   single-valued tag two values. It is deliberately *not* decided over the
   reachability nodes: a merged node's value set has cardinality > 1 whenever two
   paths write different single values, which is the legal shape of a two-path
   merge, not a violation — reporting it would make the abstraction's join rule
   into a false positive generator. Dropping the path-accumulation check rests on
   **write-replaces semantics** for single-valued tags — a write supplants the
   prior value, no preceding `clear` required — so no single path accumulates a
   second value. No peer states that rule today: RDR 0002 fixes only that
   "absence from both the write block and the clear list MUST NOT imply
   deletion", and RDR 0003 owns the single-valued *marker*, not the update rule.
   It is booked as **A10** against RDR 0002 rather than assumed; if a peer lands
   accumulate or set-union semantics instead, this per-row reading is unsound
   and the path-accumulation case returns. This is the
   model-level half of RDR 0003's **single-valued** conformance conjunct. RDR
   0003 defines a conforming view as one where *every always-present key is
   present* **and** *every single-valued tag holds at most one declared value*;
   invariant 5 covers only the second. The first is not model-decidable for
   observed or recognized keys — they arrive from accessors and the caller at
   runtime, and `internal/resolve/resolve.go::assemble` merges them without
   reading any declaration — so lint checks always-presence only for **owned**
   keys, where the reachability root and writes make it decidable: a key
   declared always-present must be held in every reachable owned-state node.
   That is reported as `graph-single-valued-state`'s sibling code
   `graph-always-present-owned` (blocking).
   **The root is included, and that is a requirement on the `initial`
   declaration, not an exemption.** The declared initial owned state is a
   reachable node by construction, so an always-present owned key the `initial`
   table omits is absent at the root and takes the finding there. This is
   deliberate: a key declared always-present but not established at the root is
   absent for any row matching before the first write to it — exactly the defect
   the code exists to catch — so the correct fix is to declare it in `initial`,
   not to exempt the root. A6's requested shape (an `initial` table of owned
   `tag = value` assignments) must therefore cover every always-present owned
   key; lint reports the omission against the `initial` declaration rather than
   against an arbitrary downstream node, so the diagnostic names the authored
   site. The observed/recognized half stays
   unowned, so this RDR discharges RDR 0003:A18 **in part, not on its own**, and
   must not assume a view-level check runs — RDR 0007 states no conformance
   obligation and `assemble` performs none (A7).
6. **Owned-set-before-match** — a row that reads an owned tag (match key or
   guard atom) must find it held in every reachable owned-state that satisfies
   the row's match pattern; the declared initial owned state counts as a write.
   **The read-set is computed by lint from the row's atoms, never read from
   `Row.RequiresOwned`.** RDR 0007 (`Final`) narrows `RequiresOwned` to
   post-guard *write* dependencies and explicitly excludes guard-read keys
   ("Guard decidability is not `RequiresOwned`'s job: guard-input coverage
   is enforced by the domain rule above"), and RDR 0002 fixes
   it as "the sorted, duplicate-free set of tag keys named by the rule's write
   block and clear list". Reading that field here would check
   writes-before-writes instead of reads-before-writes and silently never fire
   on the defect this invariant exists for. Lint therefore derives the owned
   read-set as the owned-provenance keys named by the row's match pattern and
   its `all`/`unless` guard atoms — an obligation on this RDR's minimum input
   contract, which must supply per-atom keys with provenance (A1, A3).
7. **Declared terminal/escape handling** — terminal states and escape rows are
   explicit model data; lint must not infer them from missing rows. The
   input-observable defect this mints `graph-terminal-escape` for is a model
   that *relies* on such inference: a reachable non-terminal node with no
   outgoing non-escape row and no terminal declaration covering it (an implied
   terminal), or a group relying on an undeclared escape row to close coverage.
   The prohibition on lint's own reasoning is what makes those models defects
   rather than silently-accepted shapes; a fixture is authored by omitting the
   declaration the model depends on.

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
informational: the advisory tier is closed at four members — redundant rows,
unreachable rules, the bare-escape coverage closure above, and vacuous `exists`
atoms; it never changes the success disposition and must not absorb a withheld
claim. The two advisory
classes the taxonomy previously left unnamed are defined here, because RDR 0002
routes the second to this RDR by name ("a rule whose predicate set requires
`recognized` to be absent is dead, which is RDR 0006's unreachable-rule
finding"):

- **Redundant row** (`graph-redundant-row`) — a row whose accepted assignments
  are a *proper subset* of a sibling's in the same group, so it can never be the
  exact-one match. Distinct from overlap, which is a partial intersection
  between two rows neither of which subsumes the other, and is blocking.
- **Unreachable rule** (`graph-unreachable-rule`) — a row no reachable
  owned-state node satisfies, including RDR 0002's dead-rule case (a predicate
  set requiring `recognized` to be absent, which no view reaching a row can
  satisfy).

**Emission is complete, not first-failure.** RDR 0003's reporting clause governs
and this RDR adopts it by citation: lint MUST report every defect it can decide
in one pass over a row group, not the first it encounters, and withholding a
group's exhaustiveness claim MUST NOT suppress overlap, coverage, or further
withholding findings for that same group. Each unprovable dimension, each
refusing row, each overlapping pair, and any coverage gap over a provable
product is its own finding, so the emitted set never depends on row or dimension
iteration order. An overlapping pair yields **one** finding naming both rows,
not one per row; two escape rows sharing two failure classes yield one finding
per shared class.

**Two further RDR 0003 blocking triggers this RDR consumes.** Its A21 makes the
single-valued marker a *precondition for provability*: an `eq`/`in`/comparison
atom over a tag not declared single-valued has no projection and takes the
blocking outcome — that case is `graph-unprovable-coverage`, the same code as a
non-finite dimension. Its product-size clause requires "too large to prove" to
be a declared, model-independent bound the implementation MUST publish; a group
whose declared finite product exceeds it is `graph-product-too-large`. The bound
is an implementation constant, published in the command's help output and
asserted by the MVV, not a per-model input.

Finding codes:

| Invariant | Stable Code | Severity |
| --- | --- | --- |
| Dangling edge/reference | `graph-dangling-edge` | blocking |
| Dead end | `graph-dead-end` | blocking |
| Determinism / overlap (either population) | `graph-overlap` | blocking |
| Guard exhaustiveness / gap over a provable product | `graph-coverage-gap` | blocking |
| Coverage claim withheld (three triggers — see `reason` below) | `graph-unprovable-coverage` | blocking |
| Single-valued state violation | `graph-single-valued-state` | blocking |
| Always-present owned key absent from a reachable owned-state | `graph-always-present-owned` | blocking |
| Owned-set-before-match | `graph-owned-before-write` | blocking |
| Declared terminal/escape handling | `graph-terminal-escape` | blocking |
| Declared finite product exceeds the published proof bound | `graph-product-too-large` | blocking |
| Coverage closed by a bare escape row | `graph-coverage-closed-by-escape` | info |
| Row whose accepted assignments are a subset of a sibling's | `graph-redundant-row` | info |
| Rule whose selection context no reachable owned-state satisfies | `graph-unreachable-rule` | info |
| Vacuous `exists` atom over an always-present key | `graph-vacuous-atom` | info |

**`graph-unprovable-coverage` carries a `reason` discriminator.** One code, three
triggers, and each has a *different* remedy, so a machine consumer branching on
`code` alone cannot tell the author what to do. The finding MUST carry a stable
`reason` from this closed, append-only set, naming the key, dimension, or row at
fault:

| `reason` | Trigger | The author's remedy |
| --- | --- | --- |
| `dimension-not-finite` | a participating dimension has no finite declared domain or is not projectable | declare the domain |
| `tag-not-single-valued` | an `eq`/`in`/comparison atom over a tag lacking the single-valued marker (`0003::A21`) | add the marker |
| `row-can-refuse` | a participating row can refuse `guard_unevaluable` (RDR 0003's narrowing) | nothing — the model is fine and lint declines to promise |

`row-can-refuse` is *not* a model defect — it is an honest withholding — and the message MUST say so rather than reading as an authoring
error.

**`graph-vacuous-atom` discharges a locked-peer obligation.** RDR 0003 requires
an `exists` atom over an always-present key be reported as "well-formed but
vacuous … rather than rejecting it". It is `info`, so it does not change the
success disposition, and it is the fourth member of the advisory tier.

When one or more blocking findings exist, the command returns one aggregate
`CLIError` with `Code: graph-lint-failed` and `Group: GroupUserEnv`. The
implementation must extend `internal/cli/clierr.CLIError` with an optional typed
`Findings` field serialized as `findings`, so JSON mode carries individual
findings as structured data. Text mode MUST enumerate every finding's code and
message, not merely an aggregate summary; layout is free but no finding may be
dropped.

The two envelopes carry findings at **different depths, deliberately**, because
the shipped types differ.

**Failure: a top-level `findings` key on the `CLIError` object.** There is no
`{"type":"failed","error":{…}}` wrapper to hang it under —
`respond.Fail` calls `clierr.EmitJSON`, which does `json.Marshal(e)` on the
`*CLIError` and prints it, so the failure line is
`{"code":…,"message":…}` at top level. `docs/cli-output-contract.md` states this
explicitly: "`{"type":"failed", ...}` is not emitted by `Fail`; instead the
structured `CLIError` envelope … is written to stdout." The findings list is
therefore a **sibling of `code`**, reached as `.findings`, not `.error.findings`:

    {"code":"graph-lint-failed","message":…,"findings":[…]}

Any future `type` discriminator is that document's revisit note, not this RDR's
to add — this RDR MUST NOT restructure the shipped failure envelope.

**Success: inside the existing `Data any` payload**
(`{"type":"ok","data":{"findings":[…]}}`) rather than growing `respond.Success` a
verb-specific field. On success the list holds only non-blocking findings, so the
bare-escape closure is observable without inspecting the model.

On both surfaces the key is **always emitted, empty list included** — the empty
list is the proof's receipt. Two distinct `omitempty` hazards sit on that
promise, and both must be closed: the `Findings` field itself must not carry
`omitempty` (every other optional `CLIError` field does), **and**
`respond.Success.Data` is itself declared `json:"data,omitempty"`
(`respond.go::Success`), so the lint verb MUST assign `Data` a non-nil struct
value carrying the (possibly empty) list — a nil or empty-map `Data` drops the
`data` key entirely and takes the receipt with it.

**Type ownership**: `clierr` is the leaf package other internal packages import
without pulling `internal/cli`, so it MUST NOT depend on the graph-lint package.
The `Finding` struct is therefore defined **in `clierr`** as a
subsystem-agnostic record (code, severity, message, and the optional identity
fields below); the lint engine constructs values of it.

**How the atom payload stays subsystem-agnostic.** The contract requires an
atom-attributed finding to carry `Key`, `Operator`, `Literal`, and `Block`, and
an escape-scoped one to carry its failure class — but `Block` is RDR 0002
authoring vocabulary (`match`/`all`/`unless`) and the failure class is RDR
0001's `RefusalKind`. Importing either type would give `clierr` the dependency
this clause exists to prevent. So `clierr.Finding` carries those as **declared
string-typed fields, not enums**: `Key`, `Operator`, `Literal`, `Block`, and
`Class` are `string`, and `clierr` ascribes them no meaning. The producing
package owns the vocabulary and MUST populate them from its own typed values;
`clierr` only transports and serializes. This is a deliberate trade — the type
system does not enforce the "MUST carry" clause, so the MVV does (scenario 8
asserts the exact atom fields on a `graph-unprovable-coverage` finding), and
RDR 0003:A20's carrier is that field set, not a shared enum.

Disposition of every model class lint can meet:

| Input class | Exit | Envelope / error | Finding minted | Silent or loud |
| --- | --- | --- | --- | --- |
| Model with ≥1 blocking defect | 2 (`GroupUserEnv`) | `respond.Fail`, aggregate `CLIError.Code = graph-lint-failed` | every blocking finding, in the typed `findings` list | Loud |
| Clean model, coverage proved over declared domains | 0 | `respond.OK` | `data.findings` present and empty | Loud (an empty list is the proof's receipt) |
| Redundant row (assignments a proper subset of a sibling's) | 0 | `respond.OK` | `graph-redundant-row` naming both rows | Loud, never changes the success disposition |
| Rule no reachable owned-state satisfies | 0 | `respond.OK` | `graph-unreachable-rule` naming the row | Loud, never changes the success disposition |
| Always-present owned key absent from a reachable node | 2 (`GroupUserEnv`) | same aggregate failure | `graph-always-present-owned` naming key + node | Loud |
| Declared finite product over the published bound | 2 (`GroupUserEnv`) | same aggregate failure | `graph-product-too-large` naming group + bound | Loud |
| Model with no declared initial owned state | 2 (`GroupUserEnv`) | same aggregate failure | `graph-dangling-edge` naming the missing `[model]` declaration | Loud — never "nothing reachable, therefore clean" |
| Clean model, a group closed by a bare escape row | 0 | `respond.OK` | `graph-coverage-closed-by-escape` naming the row | Loud — a bare green MUST NOT satisfy this class |
| Escape row overlapping an ordinary row | 0 | `respond.OK` | none | **Silent by design** — never a runtime ambiguity (RDR 0003's two-population clause); reporting it would fail a model the runtime accepts |
| Redundant row / unreachable rule | 0 | `respond.OK` | informational entry | Loud, but never changes the success disposition |
| Group whose participating row can refuse `guard_unevaluable` | 2 (`GroupUserEnv`) | same aggregate failure | `graph-unprovable-coverage`, `reason: row-can-refuse`, naming row + atom | Loud — no non-blocking tier for this class |
| Vacuous `exists` atom over an always-present key | 0 | `respond.OK` | `graph-vacuous-atom` naming the atom | Loud, advisory — RDR 0003 requires the report, not a rejection |
| Overlap-free group with no `ambiguous_match` escape row | 0 | `respond.OK` | none | **Silent by design** — that arm is unreachable when the group cannot overlap; demanding a row would mint one `graph-unreachable-rule` flags |
| Reachable node set over the published node ceiling | 2 (`GroupUserEnv`) | same aggregate failure | `graph-product-too-large` naming the traversal + ceiling | Loud — never an unbounded run |
| Model unreadable / not conforming to RDR 0002's schema | per RDR 0002 | refused before normalization; lint never runs | none of this RDR's codes | Loud, upstream |

#### Normative Contracts

**C1**
```normative
Graph lint MUST be a blocking acceptance gate over the normalized transition
model. A model with any blocking lint finding MUST NOT be accepted for CI
success, and MUST NOT reach the resolver **through the merge boundary** the gate
guards. That boundary is the enforcement site and the only one: the kernel does
not consult a lint verdict at runtime, so a locally-edited illegal model on a
working tree can still be resolved against, and this clause does not claim
otherwise. Coupling the resolver to a lint result would put a design-time proof
on the runtime path, which this RDR's authority split rejects; RDR 0001's
stateless refusal is what protects that path instead.
```

**C2**
```normative
Graph lint MUST consume the normalized candidate-row graph from the transition
model contract. It MUST NOT define a second sparse-source parser or a parallel
transition semantics.
```

**C3**
```normative
Graph lint MUST check at least these blocking invariant classes: dangling edge,
dead end, determinism/overlap, guard exhaustiveness/gap, single-valued state,
owned-set-before-match, and declared terminal/escape handling.
```

**C4**
```normative
Graph lint MUST reject ambiguity instead of relying on source order,
rendered-row order, or first-match priority to choose between enabled rows.
```

**C5**
```normative
Coverage, overlap, and withholding MUST be decided per scoped row group as RDR
0003 defines it (`0003::Technical Design`, row-group paragraph; `0003::Normative
Contracts`, participation clause). This RDR MUST NOT define a second grouping
predicate; it supplies only which selection contexts are reachable, per the
reachability relation in `Load-Bearing Decisions`.
```

**C6**
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

**C7**
```normative
Coverage is a universal claim, so it MUST be computed from the group's authored
rows and their declared guard domains alone — never from a reachability node.
The union `union(row_i accepted assignments)` and the scoped product both range
over the group's *participating guard dimensions* as RDR 0003 defines
participation; match keys are not product dimensions
(`0003::Normative Contracts`) and contribute no assignment to either side. Lint
MUST NOT widen a group's row set by pooling rows satisfiable at a shared
abstract node: doing so makes the union a superset of any concrete view's and
certifies coverage a real state lacks. Reachability decides only whether a group
is proven at all; it never enters the coverage computation. See
`Load-Bearing Decisions`, *Soundness direction is per invariant*.
```

**C8**
```normative
Graph lint's exhaustiveness promise is narrowed exactly as RDR 0003's narrowing
clause states (`0003::Normative Contracts`, "An exhaustiveness claim MUST NOT be
stronger than the runtime it describes"): lint MUST NOT certify a row group
exhaustive when a participating row can refuse `guard_unevaluable`. A withheld
claim MUST be emitted as `graph-unprovable-coverage`, naming the participating
row and the refusing atom. This RDR cites that clause and MUST NOT restate it,
mint a second code for it, or carry it in a non-blocking tier.
```

**C9**
```normative
Escape rows MUST participate in the coverage union and MUST be overlap-checked
in one population per declared failure class, never against ordinary rows, as
RDR 0003's escape-row clause states (`0003::Normative Contracts`, "A declared
escape row participates in the coverage identity"). A group whose coverage is
closed by a bare escape row MUST emit `graph-coverage-closed-by-escape` naming
that row; a bare green MUST NOT satisfy this clause.
```

**C10**
```normative
An escape row closes coverage only for the failure classes it declares. RDR 0003
fixes that an escape row cannot rescue `guard_unevaluable` or
`owned_state_unavailable` "however bare its guard" (`0003::Normative Contracts`),
because the kernel's gate returns before `escapeOrRefuse`. Lint MUST therefore
compute the coverage union per (scoped row group × declared rescuable class) —
`no_match` and `ambiguous_match` only — and MUST NOT let a row declaring one
class close the group's other arms. Each class is separately reachable at
runtime: `escapeOrRefuse` selects escapes by `row.rescues(r.Kind)` for the kind
that actually occurred, so a group rescued for `no_match` alone is genuinely
unrescued when the other kind arises.

**The `ambiguous_match` arm is only demanded where it is reachable.** Overlap is
itself blocking, so a group whose ordinary rows pass invariant 3 cannot produce
`ambiguous_match` at runtime, and requiring an escape row for that arm would
demand a row the kernel can never reach — which invariant's own
`graph-unreachable-rule` would then flag. Lint MUST therefore check the
`ambiguous_match` arm's coverage only for a group that carries a
`graph-overlap` finding, and MUST treat the arm as vacuously closed for a group
whose ordinary population is overlap-free. `owned_state_unavailable` is a
runtime refusal class no escape row rescues and no lint finding mints;
invariant 6's `graph-owned-before-write` is the design-time check whose runtime
counterpart it is, and the two MUST NOT be conflated.
```

**C11**
```normative
An owned key declared always-present MUST be held in every reachable owned-state
node; a violation is `graph-always-present-owned`. Lint MUST NOT extend this
check to observed or recognized keys — those arrive at runtime and
`internal/resolve/resolve.go::assemble` reads no declaration — so this RDR
discharges only the owned half of RDR 0003's conformance premise.
```

**C12**
```normative
Lint MUST publish the model-independent product bound above which it declines to
prove coverage, and MUST emit `graph-product-too-large` for a group whose
declared finite product exceeds it, per RDR 0003's product-size clause. An atom
over a tag not declared single-valued has no projection and MUST take
`graph-unprovable-coverage` (`0003::A21`).
```

**C13**
```normative
Every blocking finding MUST carry a stable code, model identity, severity,
human-readable message, and the source rule/context id or source span when the
normalized model can provide one. A finding attributed to one guard atom MUST
carry that atom's `Key`, `Operator`, `Literal`, and `Block`; a finding scoped to
an escape population MUST carry the failure class.
```

**C14**
```normative
Graph lint failure MUST return one aggregate `CLIError` with code
`graph-lint-failed` and `GroupUserEnv`; the individual blocking findings MUST
remain machine-readable in JSON mode through an append-only typed `findings`
field owned by `clierr`, not through a verb-local wrapper or a text-only
`Detail` string. That field is a **top-level sibling of `code`** on the marshalled
`CLIError` — reached as `.findings` — because `clierr.EmitJSON` marshals the
`*CLIError` itself and emits no `{"type":"failed","error":{…}}` wrapper; lint
MUST NOT restructure that envelope to create one. Lint success MUST carry
non-blocking findings under the existing `respond.Success.Data` payload as
`data.findings`. On both surfaces the key MUST be emitted even when the list is
empty — the empty list is the proof's receipt — so the `findings` field MUST NOT
be `omitempty`, **and** the success path MUST assign `respond.Success.Data` a
non-nil struct value, since `Data` itself is `json:"data,omitempty"` and a nil
payload would drop the key. The `Finding` type MUST be defined in `clierr` as a
subsystem-agnostic record so `clierr` gains no dependency on the graph-lint
package. Text mode MUST enumerate every finding's code and message, which
requires extending `clierr.EmitText` (failure) and the `respond.OK` text branch
(success, which renders only `Notes`/`Warnings` today and drops `Data`) — both
are booked edits, not reuse.
```

**C15**
```normative
Graph lint findings MUST be emitted in deterministic order by finding identity:
model id, invariant code, source rule/context id or graph element id, then
normalized predicate/write fingerprint. Model id leads for identity stability,
not because a run spans models — exactly one model is linted per invocation, so
within a run it is constant and the invariant code is the first discriminating
key. The fingerprint MUST be a canonical
sortable serialization, never a hash: RDR 0002's canonical atom sort
(`(key, block, operator token, literal)`) over the row's predicate atoms and
next-state tags, with RDR 0003's canonical set-literal form. When one run mixes
identity namespaces, a source rule/context id sorts before any graph element id.
```

**C16**
```normative
Graph lint MUST report every defect it can decide in one pass over a row group,
not the first it encounters, exactly as RDR 0003's reporting clause states
(`0003::Normative Contracts`). Withholding a group's exhaustiveness claim MUST
NOT suppress overlap, coverage, or further withholding findings for that group.
One overlapping pair MUST yield one finding naming both rows; escape rows
overlapping across several shared failure classes MUST yield one finding per
shared class.
```

**C17**
```normative
The advisory tier is closed at `graph-coverage-closed-by-escape`,
`graph-redundant-row`, `graph-unreachable-rule`, and `graph-vacuous-atom`. A
redundant row is one whose accepted assignments are a proper subset of a
sibling's in the same group; an unreachable rule is one no reachable owned-state
node satisfies, including RDR 0002's dead-rule case; a vacuous atom is an
`exists` atom over an always-present key, which RDR 0003 requires be reported as
well-formed-but-vacuous rather than rejected. Advisory findings MUST NOT change
the success disposition.
```

**C18**
```normative
A model that declares no initial owned state MUST be rejected with a blocking
finding; lint MUST NOT treat an absent root as an empty reachable set and report
a clean model.
```

**C19**
```normative
The authoritative CLI surface for graph acceptance MUST be the root command
`intrastate lint` or a same-engine CI invocation of that command. Pre-commit
hooks, aliases, and resolver-local validation flags MAY call that engine, but
MUST NOT define different acceptance rules.
```

**C20**
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
  writes applied and clears removed. **Escape rows are edges too, and they are
  self-loops**: RDR 0002 fixes that an escape row "carries neither a write block
  nor a clear list, normalizes to a row with both empty"
  (`0002::Normative Contracts`), so an escape rescue is a legal transition whose
  successor node equals its source. Lint MUST model them as such rather than
  omitting them, because a node reachable only through an escape rescue is
  reachable at runtime; omitting the edge would let invariant 2 accuse the
  recovery arm of being a dead end and invariant 6 accuse rows downstream of it.
  Escape self-loops are excluded from invariant 2's "at least one outgoing row"
  test — a self-loop is not progress — so a node whose only exit is an escape
  rescue is still a dead end, and from invariant 3's ordinary-row overlap
  population, which RDR 0003's two-population rule already separates. Match atoms over observed/recognized tags
  and all guard atoms are **not** pruned — every such edge is taken as
  traversable. A selection context is *reachable* when some reachable
  owned-state satisfies its match pattern; a row's *reachable predecessors* are
  the rows on any root-to-source path. This is a syntactic dataflow over
  declarations, writes, and clears — no accessor execution, no runtime trace,
  no guard evaluation — so it is total, and it over-approximates the runtime
  (a path lint walks may be infeasible at runtime). What that over-approximation
  guarantees is **per check, not global** — see *Soundness direction* below. It
  is a different predicate from RDR 0003's "can refuse" test, which is decided
  over the optionality field alone and never consults this relation.
- **Join rule and termination** — the traversal is a **fixpoint over merged
  nodes, never path-sensitive**: two edges reaching the same successor produce
  one node whose per-tag value sets are the union of theirs (a lattice
  widening), and the worklist runs to a fixpoint. The lattice is finite —
  finitely-declared tags range over their declared domain's subsets, a tag with
  no finite domain over `{held, absent}` — so the fixpoint terminates on cyclic
  graphs, including the self-loops RDR 0002 models. A merged node holds every
  value some path brings, so a check against it can accuse a path the runtime
  never walks. A path-sensitive reading would be exponential and is rejected.
- **Soundness direction is per invariant, not global** — merging widens the set
  of concrete views a node stands for. Whether that is safe depends on the
  *quantifier* of the check, and this lint makes both kinds. The classification
  is what decides where a merged node may be read directly and where it may not:
  - **Existential checks read merged nodes directly.** Overlap ("do two rows
    match together?"), dangling edge, and owned-set-before-match ask whether
    *some* witness exists. A merged node admits a superset of concrete views, so
    a witness real at a concrete view is still visible at the node: false
    positives possible, missed defects not.
  - **Universal checks must not.** A ∀-claim evaluated against a widened domain
    gets *easier* to satisfy, which is the false-green direction. Two invariants
    are universal, and each needs its own remedy:
    - **Coverage** (`union(row accepted assignments) == scoped product`,
      `0003::Technical Design`) asks whether *every* assignment is matched. The
      remedy is that it never reads a node at all: membership is the authored
      match pattern, so the union is computed from the group's own rows and
      reachability only decides whether the group is proven. Pooling rows
      satisfiable at a shared node would make the union a superset of any
      concrete view's — a gap real at one view closed by a row that is not a
      candidate there, green from lint and `no_match` from `Resolve`.
    - **Dead end** (invariant 2) asks whether *every* value in a node's per-tag
      set meets a declared terminal. It is anti-monotone in merging for the same
      reason: the more paths converge, the less likely one terminal covers the
      whole set — so an unsplit merged-node test accuses ordinary converging
      flows. The remedy here cannot be "don't read the node", since the check is
      *about* nodes; it is to **split** the node on terminal-participating keys
      first, recovering exactness. See invariant 2.

  Merging is therefore load-bearing only for the existential checks. The two
  universal ones each pay for their soundness — coverage by never consulting a
  node, dead end by splitting one — and neither is left resting on the
  over-approximation.
- **Owned-set-before-match is decided over nodes, not paths** — invariant 6's
  "every reachable owned-state that satisfies the row's match pattern" is the
  normative form, evaluated against merged fixpoint nodes. "Reachable
  predecessors" is descriptive prose for the same relation, not a second
  per-path enumeration; where the two readings differ the node form governs. A
  row *preserves* a tag when it neither writes nor clears it, so the tag's value
  set passes through the edge unchanged.
- **Wire / byte format** — graph-lint failure uses aggregate
  `CLIError.Code = "graph-lint-failed"` and individual finding codes from the
  taxonomy above. JSON mode carries findings as append-only structured data at
  two different depths, because the shipped types differ: **`.findings`** as a
  top-level sibling of `code` on the failure line (a new typed `clierr` field —
  `EmitJSON` marshals the `CLIError` itself, so there is no `error` wrapper to
  nest under) and `data.findings` inside `respond.Success.Data` on success. Both
  always emit the key, empty list included, which requires the `findings` field
  to skip `omitempty` and the success path to assign `Data` a non-nil struct.
  Text mode enumerates every finding's code and message, via extensions to
  `clierr.EmitText` and the `respond.OK` text branch.
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

The command's **input contract is normative** (the name cannot be authoritative
while the way to invoke it is illustrative); only cosmetic flag spelling defers
to RDR 0005. `--flow <id>` selects the model through RDR 0005's config
discovery; `--model <path>` names one explicitly. They are mutually exclusive —
supplying both is a `GroupUserEnv` usage error, not a precedence rule. Exactly
one model is linted per invocation: a corpus is linted by invoking the command
once per model, so a run's findings never span models. The finding's `model`
field is the `[model].id` from the model itself, never the path.

```sh
intrastate lint --flow rdr --as=json
intrastate lint --model ./path/to/rdr-transition-model.toml --as=json
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
| Tag declaration model, scoped row group, finite-domain and escape-row semantics, narrowing, complete-emission, single-valued-as-provability (A21), product bound, conformance definition | RDR 0003 | Verified peer RDR | Enables overlap and exhaustiveness proofs; cited, not restated. Escape rows close coverage per declared rescuable class only. |
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
| Structured CLI failures | `internal/cli/clierr` | No graph-lint codes yet; `EmitJSON` marshals `CLIError` flat (no `error` wrapper); `EmitText` renders code/message/detail/hint only | **Extend** | Add aggregate `graph-lint-failed` plus finding codes; add the non-`omitempty` typed `Findings` field as a **top-level** sibling of `code`; extend `EmitText` to enumerate findings. Model defects map to `GroupUserEnv`. |
| Text/json output routing | `internal/cli/respond` | `Success.Data` is `json:"data,omitempty"`; the `OK` text branch renders only `Notes`/`Warnings` and drops `Data` | **Extend** (not reuse) | Routing reuses `respond.OK` / `respond.Fail`, but the empty-list receipt needs `Data` assigned a non-nil struct, and text-mode enumeration needs the `OK` text branch to render the findings list. Both are booked edits. |
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
  hit; the cure is a clearer model, never a weaker lint. **There is no
  suppression or waiver mechanism** — no inline ignore comment, no allowlist. A
  waiver channel would become the way blocking findings get skipped, which is
  the failure this RDR's blocking authority exists to prevent.

  **The rate is not asserted; it is bounded by design and measured by the MVV.**
  Three structural false-positive generators were identified at pre-lock and
  each is closed rather than tolerated: dead-end is tested on *split* nodes, not
  merged ones (invariant 2), so converging flows are not accused; coverage is
  computed from the group's authored rows alone, so no node's imprecision can
  produce a spurious gap; and always-present is reported against the `initial`
  declaration, whose fix is one authoring edit. What remains is the
  guard-infeasible path, which is inherent to a guard-blind traversal.
  Withholding rather than failing is the general shape: where the abstraction is
  imprecise, lint declines to promise instead of accusing.

  The residual rate is an MVV measurement, not a prediction — scenario 23 lints
  the checked-in model (A8) and records the count of blocking findings a model
  its authors consider correct receives. **The trigger is mechanical: if that
  count is non-zero on a model the maintainers accept, the response is a
  successor RDR on guard-aware pruning** — still not a suppression flag. The
  accepted cost is real: lint can push a model toward explicitness it would not
  otherwise need.
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
absent always-present owned keys, owned-tag read-before-write, undeclared
terminals, and over-large products. Silent failure would mean accepting a model
with a blocking invariant defect, letting a hook/alternate command use different
rules, reading the same model under a different grouping or escape-row
population than RDR 0003, reporting only the first defect per group so a fixed
model still fails, or treating a model with no declared root as clean because
nothing is reachable — the normative command/CI authority, shared lint engine,
complete-emission clause, missing-root rejection, and by-citation adoption of
RDR 0003's clauses are meant to prevent that. Diagnosis starts from the finding code plus source
rule/context id or source span, and the atom or failure class when carried.

## Implementation Plan

### Prerequisites

- [x] A5 CI gate surface verified: `Makefile` and GitHub Actions can host the
  production `intrastate lint` gate once the command and its subject exist.
- [ ] A5 production gate proof pending MVV scenario 7: the new `graph-lint`
  job in `.github/workflows/ci.yml` must invoke the built `intrastate lint`
  command over the checked-in transition model, asserting the JSON `code`. The
  shipped `lint` job key is the golangci-lint gate and is not reusable; local
  `make check` parity is convenience and needs the `build` edge it lacks.
- [ ] A8: a transition model is checked into this repo for that gate to lint.
  No such model exists today, so scenario 7 has no subject until one is
  authored and homed.
- [ ] A6: RDR 0002's schema declares the initial owned state and terminal
  states. Decided at Stage 4 in favour of the RDR 0002 arm over a sidecar; the
  clause amends RDR 0002's layout enumeration and its `[model]` contents
  clause. RDR 0002 is `Draft`, so this schedules an edit on an open peer rather
  than reopening a locked one.
- [ ] A10: RDR 0002 states write-replaces semantics for single-valued tags.
  Invariant 5's per-row reading drops the path-accumulation case on that
  semantics; RDR 0002 is `Draft`, so this too is a scheduled peer edit.
- [ ] RDR 0002 normalized-row identity and RDR 0003 finite-domain predicate
  semantics are coherent enough to implement checks against.
- [x] RDR 0005 command placement is coherent enough to expose root
  `intrastate lint`: RDR 0005 owns runtime `flow` verbs and leaves graph lint
  authority to this RDR. A9 narrows the entry point: RDR 0005 defers
  transition-model config discovery, so `--model <path>` is the instantiable
  form and `--flow <id>` waits on that discovery landing.

### Minimum Viable Validation

Add a fixture-backed lint invocation that uses the same command shape intended
for CI. It must pass one legal transition model and fail one illegal model for
each blocking invariant class named in this RDR, asserting stable finding codes
and source rule/context identity in JSON mode. The illegal fixture matrix must
assert `graph-dangling-edge`, `graph-dead-end`, `graph-overlap` (both an
ordinary-row pair and an escape-row pair sharing a failure class),
`graph-coverage-gap`, `graph-unprovable-coverage` (both a non-finite dimension
and a withheld claim naming row and atom), `graph-single-valued-state`,
`graph-always-present-owned`, `graph-owned-before-write`,
`graph-terminal-escape`, and `graph-product-too-large`. It must also include one
**multi-defect group** asserting every expected code from a single run, so a
first-failure engine cannot pass. The legal matrix must include a group closed by
a bare escape row, a redundant row, and an unreachable rule, asserting
`graph-coverage-closed-by-escape`, `graph-redundant-row`, and
`graph-unreachable-rule` on success. Every blocking run must return aggregate
`CLIError.Code = graph-lint-failed` with `GroupUserEnv` exit behavior and a
machine-readable finding list; every run, clean or not, must emit the `findings`
key even when empty.

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
mode on both success and failure. Add the production command
invocation as the new `graph-lint` job in `.github/workflows/ci.yml` once the
checked-in transition model exists (A8), wiring `make check` to the same
command for local parity.

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
production command, A6 supplies the reachability root and terminal set, and A7
records that invariant 5's model-level check stands alone.

The code paths these tests exercise were reviewed at Stage 4 and none exists
yet: `internal/` and `cmd/` hold eight non-test Go files whose only registered
verb is `internal/cli/version.go::newVersionCmd`, and no graph, lint verb,
invariant engine, or reachability dataflow is present. The output path the
matrix asserts against does ship — `internal/cli/clierr/clierr.go::CLIError`,
`::ExitCodeFor`, `::GroupUserEnv`, `internal/cli/respond/respond.go::OK` and
`::Fail`, and `internal/cli/root.go::ExecuteAndEmit` — and every existing
optional `CLIError` field is `omitempty`, so appending the typed `findings`
field this RDR requires is wire-compatible with the shipped envelope. That field
is the one exception to the `omitempty` habit: it must serialize even when
empty, since an empty list is the proof's receipt.
`internal/cli/respond/respond.go::Success` already carries `Data any`, the host
for the success-side findings list, so the success half needs no new field on
`Success` itself.

Desk trace — one lint invocation, with every normative clause in force at each
step and a witness from the fixture corpus:

| Step | Clauses in force | Witness / outcome |
| --- | --- | --- |
| 1. Load + normalize | RDR 0002 schema and `[model]` clause; "MUST NOT define a second sparse-source parser" | Non-conforming source is refused before normalization; lint never runs (upstream disposition row). |
| 2. Derive graph view | reachability relation (`Load-Bearing Decisions`); A6 initial/terminal declarations; node-ceiling clause | Root = declared initial owned state; escape rows are self-loop edges. Over-approximating, so total — but the guarantee is **per invariant**: safe for existential checks, not for coverage (step 6). Exceeding the published node ceiling is `graph-product-too-large` on the traversal, never an unbounded run. |
| 3. Scope row groups | row-group clause (RDR 0003's, cited not restated); "MUST NOT define a second grouping"; source-state clause | Membership is the **authored match pattern** + recognized outcome, per RDR 0003 and matching `Resolve`'s `view.matches(row.Match)`. Reachability decides only *which* groups are proven, never which rows are in one — so the group count is bounded by the authored rule count, not the lattice. |
| 4. Overlap | invariant 3; RDR 0003's two-population clause | Ordinary×ordinary → `graph-overlap`. Escape×escape sharing a class → `graph-overlap`, once per shared class. Escape×ordinary → **no finding** (scenario 9a). |
| 5. Withholding test | narrowing clause; "can refuse" as a syntactic test over the declared optionality field | Decided over **every** participating row, escape rows included; total because it reads declarations, never the graph. |
| 6. Coverage | invariant 4; presence-dimension projection clause; escape class-scoping clause; **authored-rows-only clause** | `union(row_i accepted assignments) == scoped product`, per (group × declared rescuable class), over the group's **participating guard dimensions** — match keys are not product dimensions (`0003`) — product including the `{absent}` assignment of an `exists` atom over an optional key (scenario 11). Computed from the group's authored rows alone; no node enters the computation, so pooling rows at a shared abstract node (the false green of scenario 20) cannot arise. The `ambiguous_match` arm is checked only for a group that overlaps; vacuously closed otherwise. |
| 7. Precedence when 5 and 6 both bear | narrowing clause governs: "no non-blocking tier for this class"; the advisory tier "must not absorb a withheld claim" | A group both closable by a bare escape row **and** carrying a row that can refuse resolves to `graph-unprovable-coverage` (blocking, exit 2) — **not** a success with `graph-coverage-closed-by-escape`. Withholding dominates closure, since at runtime the refusal returns before the escape row is consulted (scenario 8). |
| 8. Emit | aggregate-`CLIError` clause; deterministic-order clause; `respond`-only clause; complete-emission clause | Every decidable defect in the group is emitted, not the first. Blocking → one `graph-lint-failed` / `GroupUserEnv` / exit 2, findings at top-level `.findings` on the marshalled `CLIError` (no `error` wrapper exists — `EmitJSON` marshals the error itself). Clean → `respond.OK` with the advisory list at `data.findings`, which needs `Data` assigned a non-nil struct since `Data` is itself `omitempty`. Both emit the key even when empty. Order by finding identity, never source order. |

No CONTRADICTION row. Step 7 is the one place two clauses bear on the same
verdict; the narrowing clause's own text settles the precedence, and this trace
records it so an implementer does not read the info code as an escape hatch.

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
   accepting or weakening the coverage guarantee, carrying the tag key of the
   unprovable dimension. The atom fields are required only when the finding is
   attributed to one guard atom (the withheld-claim arm, scenario 8); the
   dimension arm must still name the dimension, so no arm of this code may carry
   the code alone.
4. **Scenario**: text and JSON modes for the same legal and illegal fixtures.
   **Expected**: both modes return the same exit behavior with no direct
   stdout/stderr writes. JSON finding order is deterministic by finding
   identity, asserted as full-list equality against a golden file, not merely as
   a sorted property. Text mode's oracle is concrete: the set of finding codes
   appearing in text output equals the set in JSON output for the same fixture.
   Layout and wrapping are unasserted; a text renderer that drops any finding
   fails.
5. **Scenario**: a finding whose normalized row lacks a rule id but has a source
   span or graph element id.
   **Expected**: the finding remains actionable and stable by carrying the
   available source span or graph element id in the identity fields.
6. **Scenario**: command placement.
   **Expected**: root `intrastate lint` is the canonical tested command path.
   The assertable-now form, since no alias exists at lock: the lint package
   exposes exactly **one** exported engine entry point and one request builder,
   and a test asserts the root command calls them — so a later alias has no
   second constructor to diverge through. If an alias is added, it must route
   through those same two symbols and a test must assert flag and result parity.
7. **Scenario**: CI-shaped invocation.
   **Expected**: `.github/workflows/ci.yml` contains a **`graph-lint`** job —
   not `lint`, which is the shipped golangci-lint job and must survive
   untouched — that runs `make build` and then `./bin/intrastate lint --as=json`
   over the checked-in transition model, with the assertion reading the JSON
   `code` field, not the exit integer alone. Mechanically checkable: both jobs
   exist in that file, the new one names that command, and the model path it
   passes is the checked-in model (A8) — not a hook wrapper, unit-test-only
   engine path, or fixture-only corpus. A deliberately illegal commit to the
   checked-in model must fail that job.
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
    **Expected**: `graph-owned-before-write` naming the row. Paired dead-end
    case, covering the other half of the same over-approximation contract: a
    node reaching a declared terminal on one path and a sink on another emits
    `graph-dead-end`, since the merged node does not wholly satisfy the
    terminal. Both are accepted false positives whose cure is an explicit write,
    clear, or terminal declaration on the model, never a guard-aware lint.
11. **Scenario**: presence dimension — a scoped row group carrying an `exists`
    atom over a key declared optional, whose rows cover every value assignment
    but leave the `{absent}` assignment uncovered; plus a control group whose
    `exists` atom is over a key declared always-present.
    **Expected**: the optional-key group emits `graph-coverage-gap` rather than
    passing, proving the presence dimension was not dropped from the product;
    the always-present control passes, since such a key contributes no presence
    dimension. The control's key must carry RDR 0003's **explicit**
    always-present marker: an unmarked key defaults to optional and takes the ×2
    presence row, so an unmarked control would not test the distinction.
12. **Scenario**: precedence — a group both closable by a bare escape row **and**
    carrying a row that can refuse `guard_unevaluable` (desk-trace step 7's
    conjunction).
    **Expected**: `graph-unprovable-coverage`, exit 2, and
    `graph-coverage-closed-by-escape` **absent** from the payload. The negative
    half is the assertion that matters — withholding dominates closure.
13. **Scenario**: complete emission — one group carrying an overlapping ordinary
    pair, a coverage gap, and two rows that can refuse.
    **Expected**: findings for all of them in one run — one `graph-overlap`
    naming both rows (not one per row), the coverage finding, and one
    `graph-unprovable-coverage` per refusing row. A first-failure engine
    emitting only the first fails this scenario.
14. **Scenario**: advisory tier — a row whose accepted assignments are a proper
    subset of a sibling's, and a rule whose predicate set requires `recognized`
    to be absent.
    **Expected**: exit 0 with `graph-redundant-row` and `graph-unreachable-rule`
    in `data.findings`. A lint that silently drops the advisory tier fails.
15. **Scenario**: missing root — a model declaring no initial owned state.
    **Expected**: blocking `graph-dangling-edge` naming the absent declaration,
    exit 2 — never exit 0 on an empty reachable set.
16. **Scenario**: empty-list receipt — the legal fixture in JSON mode.
    **Expected**: the `findings` key is present with value `[]` on success
    (`data.findings`), asserted against a golden JSON file so an omitted key
    fails, and present on the failure envelope alongside blocking findings.
17. **Scenario**: always-present owned key — a key declared always-present that
    some reachable owned-state node does not hold, plus an *observed* key
    declared always-present.
    **Expected**: `graph-always-present-owned` for the owned key and **no**
    finding for the observed one, proving lint claims only the owned half of
    RDR 0003's conformance premise.
18. **Scenario**: product bound — a group whose declared finite product exceeds
    the published bound, and an atom over a tag not declared single-valued.
    **Expected**: `graph-product-too-large` for the first,
    `graph-unprovable-coverage` for the second, and the enforced bound visible
    in the command's help output.
19. **Scenario**: escape-row class scoping — a group whose only escape row
    declares `no_match` alone, in two variants: (a) the group's ordinary rows
    overlap; (b) they do not.
    **Expected**: the coverage union is computed per (group × declared class).
    In (a) the `ambiguous_match` arm is reachable, so its uncovered state emits
    its coverage finding alongside `graph-overlap`. In (b) that arm is
    **vacuously closed** and MUST NOT emit a finding — an overlap-free group
    cannot produce `ambiguous_match`, and demanding an escape row for it would
    require a row `graph-unreachable-rule` then flags. Both variants emit
    `graph-coverage-closed-by-escape` for the `no_match` arm.
20. **Scenario**: grouping soundness — two rows with **different** match
    patterns (`status=a` and `status=b`), each individually non-exhaustive over
    a shared guard dimension, whose accepted assignments are complementary so
    their union would cover the product. Two paths write `a` and `b`, converging
    on a node whose value set for `status` is `{a, b}`, so both rows are
    satisfiable there.
    **Expected**: **two groups, and `graph-coverage-gap` on each** — membership
    is the authored match pattern, so the rows never pool. The negative control
    is the point of the fixture: **no run may exit 0 with an empty blocking
    list**, since a green here is exactly the false green a merged-node union
    would produce (the union closes, but at the concrete view `status=a` only
    the first row is a candidate and the dimension is uncovered). Paired
    positive control: give both rows the *same* match pattern and they are one
    group whose union genuinely closes — lints clean.
21. **Scenario**: converging terminal — two paths reaching one node, one having
    written a terminal-satisfying value and the other a live value.
    **Expected**: `graph-dead-end` is **not** emitted for the terminated branch
    (invariant 2 splits before testing); it *is* emitted if the live branch has
    no outgoing row. A single merged-node test would fail both, so this fixture
    is the regression guard on the split rule.
22. **Scenario**: node ceiling — a model whose declared owned tags and domains
    make the reachable owned-state set exceed the published ceiling.
    **Expected**: `graph-product-too-large` naming the traversal and the
    ceiling, at `GroupUserEnv` / exit 2 — never an unbounded run and never a
    partial green. The ceiling is asserted to be the published constant from
    the command's help output, not a test-local value.
23. **Scenario**: false-positive census — lint the checked-in transition model
    (A8) in its accepted state.
    **Expected**: the blocking finding count is recorded. Zero is the pass
    condition; a non-zero count on a model its maintainers accept is the
    mechanical trigger for the guard-aware-pruning successor RDR named in
    `Consequences`, and MUST be reported rather than waived.

### Performance Expectations

No *throughput* target is load-bearing: lint is a single-invocation static check
— load model data, build the owned-state graph, run deterministic invariant
checks, render one terminal result. Determinism depends on RDR 0002's
candidate-row/source identity and stable dump ordering plus this RDR's finding
identity tuple, not on source order, map order, or first-match priority.

**Node count is bounded, and that bound is load-bearing.** Total work is
(reachable nodes) × (groups proven per node) × (product proof up to the
published bound). The published bound covers only the last factor, so the node
count needs its own ceiling: the owned-state lattice is finite but its size is
exponential in owned tags × declared domain sizes, and merging shrinks it only
where paths converge. Lint MUST therefore publish a **model-independent node
ceiling** alongside the product bound, and MUST emit `graph-product-too-large`
naming the traversal (not a group) when a model's reachable node set exceeds it,
rather than running unboundedly. Both bounds are implementation constants in the
command's help output and asserted by the MVV. Group *membership* is syntactic
(the authored match pattern), so the group count is bounded by the authored rule
count regardless of the lattice; only the reachability filter scales with nodes.

## Finalization Gate

Responses: 0006-graph-lint-authority-and-guarantees/artifacts/gate.md (Gate PASS 2026-08-23)

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
