# Recommendation 0010: Owned state is optional: stateless decision tables are first-class

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-08-26
- **Status**: Draft
  <!--
  - `Deferred` is the parked-with-a-revisit-trigger status for a
    Draft that cannot proceed because **no acceptable mechanism
    exists yet** — every in-our-control path is ruled out and the
    one that would work is outside our control. It is a *pause in
    the lifecycle*, not an exit from it: the RDR stays intact and
    re-enters at the stage it stopped when the trigger fires.
    Carry the condition on the live value:
    `Deferred [revisit when <condition>]`, and say in the same
    field what was ruled out and why (Alternatives Considered
    carries the long form). Distinct from `Abandoned`, which is
    terminal — an Abandoned RDR is closed, owes a post-mortem, and
    never re-enters. A Deferred RDR owes **no** post-mortem
    (nothing was implemented), and its `Priority` records what the
    fix is *worth*, not what is scheduled. Do not defer merely to
    park work that is possible but unfunded — that is a Priority,
    not a Status.
  - `Demoted` is the terminal status for an RDR judged
    *not RDR-shaped* — the decision was never a real
    design fork, so it leaves the RDR lifecycle and is
    refiled as a plain issue. Carry the destination on the
    live value: `Demoted [→ <issue link>]`, and record the
    same link under **Related Issues**. A `Demoted` RDR runs
    no further stages. (Distinct from the 07.1 *demotion*
    below, which is a `Final → Draft` flip that keeps the
    RDR in the lifecycle — that flip never writes
    `Status: Demoted`; see the disambiguation note there.)
  - A Draft demoted from Final by the 07.1 cluster gate
    carries a qualifier on the live value:
    `Draft [revised from Final YYYY-MM-DD; re-verify A2,A4
    — <one-line reason>]`. It is still a `Draft` for every
    binary Draft/Final gate; only Stage 4 (scoped
    re-verify) and Stage 7 (re-lock) parse the qualifier.
    The Stage 7 flip to `Final` overwrites the whole value,
    so the qualifier self-clears at re-lock — no separate
    cleanup. This 07.1 "demotion" is a *verb* describing the
    Final→Draft flip; it is **not** the `Demoted` status
    above (which exits the lifecycle to an issue) — do not
    conflate the two. (`Reverted` above is the unrelated
    terminal "implementation rolled back" status — also do
    not conflate.)
  - A Final tolerated at the 07.1 gate under a JOINT-DECISION
    carries `Final [joint decision → <home §-anchor>: <the
    open question>]`. It is still a `Final` for every binary
    gate. The qualifier is an **open obligation, not a
    coherence claim**: it says the named question is
    unanswered here, not that this RDR agrees with the answer.
    So it does not self-clear. When the home answers, this RDR
    owes a scoped check of that answer against its own
    normative fences before it re-locks or implements —
    consistent → drop the qualifier and record the clearing;
    contradicts fenced text → a 07.1 SPEC-DEFECT. A re-lock
    that comes first carries the qualifier forward unchanged;
    it is never silently dropped.
  -->
- **Type**: Architecture
- **Profile**: foundational — provisional; redefines what a model *is* across the table grammar (0002), the graph-lint root rule (0006), and the resolve demand set / wire payload (0005), so it spans modules and is a cross-RDR producer.
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
- **Priority**: High
- **Related Issues**: intrastate#zdat (seed; stays open as the defect tracker), rdr#thsc (umbrella), rdr#tmxk (the consumer model that motivates this), intrastate#1mv1 (sibling seed on `flow next` verb semantics — cross-cite, not a facet)
- **Predecessors**: 0002-transition-table-as-reviewable-data, 0005-skill-integration-cli-contract, 0006-graph-lint-authority-and-guarantees
- **Overrides**: 0002:C2 and 0002:C3 (the closed `[model]` layout gains an optional `class` key; the rule layout gains `[rule.emit]`) — additive; 0002:C4 ("an ordinary transition rule MUST contain a write block", `internal/table/normalize.go::normalizeRule`, category `malformed_rule_shape`) — conditioned on the `state-machine` class; 0002:C19 (closed dump column vocabulary gains `emit`); 0005:C1 (the `flow resolve` success payload gains `emit`; `--artifact` unrequired is a consequence of its own reader scoping, not an override); 0006:D-reachability-relation (the root is the declared initial owned state *or*, for the `decision-table` class, the empty owned-state node). Confirmed at propose as NOT overridden: 0006:C18 (missing root stays a blocking finding for the `state-machine` class, `internal/graphlint/analysis.go::checkDanglingEdge`), 0002:C2's write-only-owned-tag clause (unreachable with zero owned tags), 0005 DEV-8's demand set (`internal/cli/flow_exec.go::invokedReaders`, already empty over no owned keys).
- **Seam Lineage**: no prior accretion

## Problem Statement

A model author who wants intrastate only to *decide* — which row of a
decision table applies to the observed dimensions — and has no state of their
own to carry between resolves, cannot write that model as what it is. Today
they must declare a dummy owned tag, an `[initial]`, read/write accessors,
and hand `flow resolve --artifact` a scratch file in an undocumented format.
They discover this when a straightforward decision-only model is refused at
resolve with `flow-owned-state-unavailable` on the write-only key until a
scratch artifact is seeded (verified 2026-08-26 with a proof-of-concept
navigator model). The workaround obscures the model's intent, leaks a
consumer-invented file format, and adds accessor surface the author never
wanted.

The system-internal requirement: the model grammar, lint, and resolve
currently assume every model is a state machine. Every ordinary rule must
carry a write block (RDR 0002 C4, `internal/table/normalize.go::normalizeRule`
→ `CatMalformedRuleShape`); a missing `[initial]` is reported as a dangling
edge (RDR 0006 C18, `internal/graphlint/analysis.go::checkDanglingEdge`); and
resolve demands the write target as owned state (RDR 0005 DEV-8,
`internal/cli/flow_exec.go::invokedReaders` demand set = RequiresOwned ∪
guard keys). The design question is whether a model with zero owned tags is a
legal, first-class model class — and if so, what `flow resolve` returns over
it (the selected row's identity, an authored `emit` block of literal
key/values, or something else), and what each lint rule means when the owned
set is empty. Exhaustiveness/coverage lint over observed dimensions must keep
working unchanged — the PoC surfaced a 208/512 coverage hole in a partial
table, which is the property that makes stateless tables worth supporting at
all. Stateful models are unchanged; intrastate stays generic (no consumer
knowledge, no artifact discovery).

## Critical Assumptions

- **A1 Seeding reachability at the empty owned-state node makes every group of
  a decision-table model reachable, so overlap and coverage produce the same
  findings the scratch-tag PoC produced (the 208/512 hole reproduces with the
  owned tag removed).**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: `internal/graphlint/analysis.go::nodeSatisfiesMatch`
    consults only owned atoms, so an owned-atom-free context is satisfied by
    the ∅ node; the spike re-lints the PoC model with its owned tag deleted
    and `class = "decision-table"` set, asserting the same coverage finding
    set cell-for-cell (not the count), and that neither the product bound
    (`0006:C12`) nor the unprojectable-group size suppression withholds the
    table's single group.
  - **If wrong**: coverage is silently vacuous over decision tables — exit 0
    with an unproven table, the exact outcome 0006:C18 forbids.
- **A2 With zero owned tags the reader demand set is empty, so `flow resolve`
  and `flow next` invoke no reader, raise no `flow-artifact-missing`, and a
  bound `--artifact` for an uninvoked role is ignored rather than refused.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: `internal/cli/flow_exec.go::invokedReaders` unions
    `Row.RequiresOwned` with guard-owned keys; `internal/cli/flow_exec.go::runReaders`
    checks bindings only for the invoked names; the demand is keyed on
    owned provenance, so an observed-tag reader a decision table may still
    declare is never invoked by `next`/`resolve` (only `read-state` runs
    every declared reader, unchanged).
  - **If wrong**: a decision table still needs a scratch artifact — the
    problem statement's refusal survives.
- **A3 The kernel returns a plan over an empty owned view when every
  candidate row carries empty `RequiresOwned`, `Writes`, and `NextTags`.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: `internal/resolve/resolve.go::missingOwned` iterates
    `RequiresOwned` only; escape rows already normalize to both-empty
    (`0002:C14`, `0002:C15`) and resolve today.
  - **If wrong**: `flow resolve` over a decision table refuses on a kernel
    precondition, and the class needs a kernel change RDR 0001 owns.
- **A4 `emit` can join the closed dump column vocabulary under `0002:C19`, and
  every checked-in model or fixture carrying an explicit `[dump]` column list
  is updated in the same change.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: `internal/table/dump.go::dumpColumns` is the vocabulary;
    `internal/table/load.go::loadDump` refuses an omitted column
    (`CatMalformedDumpDeclaration`); 103 fixtures under
    `internal/table/testdata/` declare `[dump]` today (plus 0002's spike
    fixtures under its evidence tree); `models/rdr.toml` declares none.
  - **If wrong**: every existing `[dump]` declaration refuses at load until
    edited, or `emit` is undumpable and the dump stops carrying "every
    field" of the normalized value.
- **A5 The kernel's exact-one selection is the decision-table hit policy the
  motivating consumer needs; no first-hit, priority, or collect policy is
  required.**
  - **Status**: Pending
  - **Method**: Prior Art
  - **Evidence**: model prior only — ⚠ no corpus coverage for decision-table
    hit policies (research/prior-art.md K1); the consumer model (rdr#tmxk)
    asks that "every status×profile×ca×lens cell has exactly one row".
  - **If wrong**: a decision table that relies on row precedence is refused
    `flow-ambiguous-match` and the class needs a selection rule the kernel
    does not have.
- **A6 A flat, string-valued `[rule.emit]` table is sufficient for the
  motivating consumer's answer.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: rdr#tmxk — "the selected row (or an `emit`) is the answer";
    the PoC carried the answer as a single enum value written to `next`.
  - **If wrong**: consumers pack structure into strings (a command with
    arguments is one string today), and the wire shape needs an `any`-typed
    value later (a widening, not a break).
- **A7 `internal/graphlint/engine.go::Fingerprint` and 0006's finding
  identity do not read `emit`, so editing an emit block never changes a
  finding's identity.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: the fingerprint is computed over the row's graph-bearing
    fields; `emit` is not an edge, write, or predicate — asserted by a test
    that edits only an emit block and compares fingerprints, not by reading.
  - **If wrong**: golden lint outputs churn on emit edits, and `0006:C15`'s
    ordering key gains a non-graph input.
- **A8 Strict decoding refuses an unknown `[model]` key, so a
  `class = "decision-table"` model loaded by a binary predating this RDR
  refuses as `unknown_schema_field` rather than silently loading as a
  machine.**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: `internal/table/source.go::sourceModel` is decoded
    strictly (`0002:A1`); `CatUnknownSchemaField` is the category — run the
    pre-change binary over the MVV model and assert that refusal.
  - **If wrong**: an old binary lints a decision table as a rootless machine
    and reports `graph-dangling-edge` — loud, but misattributed.
- **A9 The "otherwise" row of a decision table is an escape row rescuing
  `no_match` that carries `[rule.emit]`; it closes coverage only under
  0006's advisory (`graph-coverage-closed-by-escape`), never as a silent
  default, and the rescued plan reports `escaped = true`.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: `0006:C10` scopes escape closure to declared classes;
    `0002:C5` admits `no_match`; C3 admits emit on escape rows.
  - **If wrong**: authors reach for an overlapping catch-all ordinary row
    and hit `flow-ambiguous-match`, or a default row hides a coverage hole.
- **A10 Every 0006 finding code is either exercised unchanged (overlap,
  coverage, withheld, redundant, vacuous-atom, product bound, node ceiling)
  or provably silent (dangling-edge root arm, dead end, always-present,
  owned-before-match, single-valued, terminal-escape, unreachable-rule) over
  a zero-write model.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: the graph-lint fixture pair asserts the exact finding
    list — one coverage finding on the partial table, an empty list on the
    complete one — against the full taxonomy in
    `internal/graphlint/taxonomy.go`, not a subset.
  - **If wrong**: a machine-only invariant fires on every decision table, or
    stays silent on a defect it should report.
- **A11 `flow next` over a decision-table model succeeds unchanged — no
  owned demand, empty `required`/`next`/`writes` per candidate — and needs
  no refusal; re-verified against cli/0011's candidate predicate once 0011
  locks.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: `0005:C1` next preview reads normalized data only; the
    MVV's `flow next` step asserts exit 0 with empty `required` per
    candidate. The candidate predicate is owned by the sibling seed on
    `flow next` (intrastate#1mv1), proposed after this RDR — its element
    is cited once it has one.
  - **If wrong**: `next` over a decision table refuses or reports nonsense,
    and this RDR owes a class-specific clause on 0011's seam.

## Proposed Solution

### Approach

A model declares its **class** in `[model]`: `class = "state-machine"` (the
default when the key is absent — every existing model is unchanged) or
`class = "decision-table"`. A decision-table model declares **zero owned
tags**, and the loader refuses a class that disagrees with the owned set in
either direction. Everything a state machine needs *because it has owned
state* — `[initial]`, `terminal`, owned readers and writers, the write block
on every ordinary rule, `--artifact` at resolve — is absent from a decision
table by construction, and each absence is enforced by a rule 0002/0005
already have rather than by a new one. The one thing a decision table cannot
express today, an *answer*, is authored as an optional `[rule.emit]` table of
literal string key/values on any rule of either class; the normalized row
carries it, the dump renders it, and the `flow resolve` payload returns it
beside the selected `rule`. Lint treats the decision-table class as a graph
with one node — the empty owned-state — so 0006's reachability filter admits
every group and overlap/coverage run unchanged over the observed dimensions;
the machine-only invariants (root, dead end, always-present, owned-before-
match, terminal handling) are vacuous by construction and stay silent. The
missing-root finding (`0006:C18`) is untouched for the state-machine class.

The class is **declared, not inferred** from the empty owned set. The
sibling discriminator in this codebase — `Row.Kind` (ordinary vs escape,
dump column `kind`) — is inferred from the *presence* of an `escape` list
(`internal/table/normalize.go::normalizeRule`); nothing in the tree infers a
class from an absence, and 0006 rests on exactly the rule that an absent root
is a defect, never an empty graph (`0006:C18`). ⇒ inferring "decision table"
from "no owned tags" would turn a forgotten `provenance = "owned"` into a
lint-clean table; declaring it makes the mismatch a load failure both ways.

### Technical Design

Three seams change, each additively on its owner's grammar:

1. **Grammar and load** (`internal/table`): `sourceModel` gains `class`;
   `sourceRule` gains `emit`; `normalizeRule` conditions the no-write-block
   arm of `CatMalformedRuleShape` on the class; the loader adds a
   class/owned-set agreement check (`CatMalformedModelDeclaration`, new arm);
   `Row` gains `Emit []TagValue`-shaped literals sorted by key; `dumpColumns`
   gains `emit`. Every site that today asks `len(Initial) == 0` or would ask
   `len(owned) == 0` routes through one loaded-model accessor for the class,
   so there is exactly one place the question is answered. The kernel row (`internal/resolve/resolve.go::Row`) is
   untouched — `emit` never crosses the kernel boundary; the CLI joins it
   back on `Plan.RuleID`, which `0002:C4` makes unique.
2. **Lint** (`internal/graphlint`): `reach` seeds the traversal at the ∅ node
   when the model's class is decision-table; `checkDanglingEdge`'s
   missing-root arm keys on class, not on `len(Initial)`. No taxonomy change.
3. **CLI** (`internal/cli`): `resolvePayload` gains `emit`; text mode renders
   it. `invokedReaders` and `runReaders` are unchanged — the empty demand set
   already makes `--artifact` unrequired (0005:C1 "A reader no candidate row
   needs MUST NOT run"). `flow next` is not touched by this RDR (its
   candidate predicate is cli/0011's).

Data flow for a decision table: TOML → loader (class check, no owned tags,
rows with empty `RequiresOwned`/`Writes`/`NextTags`, populated `Emit`) →
lint (root ∅, groups over observed dimensions) → `flow resolve --outcome
<o> --tag k=v…` → kernel exact-one → payload `{rule, emit, next: {},
writes: {}, …}`.

#### Normative Contracts

**C1**

```normative
`[model]` MAY carry `class`, a string whose only admitted values are
`"state-machine"` and `"decision-table"`; an absent `class` MUST be read as
`"state-machine"`. A `class` outside that set is `malformed model
declaration`. The class MUST agree with the declared owned tag set: a
`"decision-table"` model MUST declare zero tags of provenance `owned`, and a
`"state-machine"` model MUST declare at least one; disagreement in either
direction is a `malformed model declaration` load failure whose detail names
the class and the offending count. The class is model data and is carried on
the loaded model; nothing downstream infers it from the owned set.
```

Overrides, additively: `0002:C2` (the closed `[model]` layout) and `0002:C3`
("`[model]` MUST contain `id` and `version`" — still true; `class` is a third,
optional key). ⇒ every existing model loads unchanged; the class becomes a
reviewable line in the artifact 0002 makes reviewable.

**C2**

```normative
A decision-table model MUST NOT declare `[initial]`, `terminal`, or any
accessor whose `keys` name an owned tag, and its ordinary rules MUST NOT
carry a write block or a clear list. None of these is a new refusal: each is
already unauthorable once the owned set is empty — an `[initial]` key or a
write/clear key that is not a declared owned tag refuses as `unknown tag`
(undeclared), `malformed accessor binding`, or `write to non-owned tag`
(declared, non-owned), and a terminal predicate over a non-owned tag is
0006's `graph-dangling-edge` terminal arm. `0002:C4`'s "an ordinary
transition rule MUST contain a write block" is conditioned on class: it binds
the `"state-machine"` class only, and the `malformed rule shape` arm that
enforces it MUST NOT fire for a `"decision-table"` model. Rule ids, the
match-block obligation, shared contexts, guards, gates, escape rules, and
outcome binding are unchanged in both classes.
```

⇒ a decision-table row is the shape an escape row already has at the kernel
boundary — empty `RequiresOwned`, `Writes`, `NextTags` (`0002:C14`,
`0002:C15`) — so the kernel sees nothing new (A3).

**C3**

```normative
Any rule, ordinary or escape, in either class MAY carry `[rule.emit]`: a flat
TOML table whose values MUST be strings. Emit keys are NOT tag keys: they are
undeclared, uninterpreted, compared by exact byte equality, and MUST NOT be
matched, guarded, written, or read by any accessor — the same carried-through
treatment `[model.metadata]` receives under `0002:C2`. Normalization MUST
carry the block on the row as a key-sorted, duplicate-free sequence; an
absent block normalizes to an empty sequence. `emit` MUST join the closed
dump column vocabulary (`0002:C19`), rendered as `key=value` pairs in key
order, so a `[dump]` column list MUST name it, and the default column set
used when `[dump]` is absent MUST include it for both classes. `emit` MUST
NOT be added to the kernel row; it is table data joined back by rule id after selection.
```

⇒ the normalized value, and therefore the dump, is the single reviewable
carrier of the answer (0002's premise); the kernel stays JDR 0001-shaped.

**C4**

```normative
The `flow resolve` success payload MUST carry `emit`: a JSON object of
string values, keys in byte order, present as `{}` — never `null`, never
omitted — when the selected row authored none. `next`, `writes`, `clear`, `owned`, and `readers` keep their
`0005:C1` shapes and are empty over a decision-table model. Under `--as=text`
the plan MUST render each emit pair as one `key=value` line. `--artifact` is
not required when the invoked reader set is empty — a consequence of
`0005:C1`'s reader scoping, restated here, not a new rule — and an
`--artifact` binding for a role no invoked reader needs MUST be ignored.
`--outcome` remains required in both classes. `flow next`, `flow read-state`,
and `flow set-state` are unchanged by this RDR.
```

⇒ a caller's answer is `rule` plus `emit`; a consumer that needs no literal
payload maps `rule` and authors no emit.

**C5**

```normative
Graph lint over a `"decision-table"` model MUST root the reachability
relation at the empty owned-state node; the reachable set is exactly that
node, and every selection context is reachable. Overlap (invariant 3),
coverage and withholding (invariant 4), the redundant-row, vacuous-atom, and
unreachable-rule advisories, the node ceiling, and the product bound MUST run
unchanged over the authored observed dimensions. The missing-root arm of
invariant 1 (`0006:C18`), dead end (2), always-present-owned (5),
owned-set-before-match (6), single-valued state, and declared-terminal
handling (7) are vacuous by construction over this class and MUST NOT emit; lint MUST NOT report the
class as a finding of any severity. Over a `"state-machine"` model every
0006 invariant, `0006:C18` included, is unchanged.
```

⇒ "unchanged coverage over observed dimensions" is a consequence of 0006's
own reachability filter (`internal/graphlint/groups.go::checkGroups` runs
coverage only for reachable groups) once the root exists, not a parallel
coverage path.

#### Load-Bearing Decisions

- **Identity** — row identity is unchanged (`0002:C4` rule id). `emit`
  joins the normalized row *value* and so the dump's byte identity; whether
  it joins `internal/graphlint/engine.go::Fingerprint` is A7's question and
  the answer this RDR wants is *no* — emit is not graph-bearing.
- **Wire / byte format** — `emit` on the wire is a JSON object of
  string→string, keys in byte order, HTML escaping disabled like every other
  payload map (`0005:C1`); in the dump it is one column, `key=value` pairs in
  key order, bracketed like `writes`.
- **Naming** — `class` with values `state-machine` / `decision-table`.
  Rejected: `kind` (collides with `[tags.<tag>] kind`, the type-model key),
  `type` (same), `stateless = true` (names what the model lacks, not what it
  is). `emit` for the answer block. Rejected: `output` (0002 already uses
  "output" for the dump), `result`, `answer`.
- **Selection / predicate** — unchanged: the kernel's exact-one match is the
  decision table's hit policy (A5); two matching rows is
  `flow-ambiguous-match` in both classes, and an escape row rescues it the
  same way.

#### Illustrative Code

Illustrative — a two-dimension decision table and the resolve envelope it
yields; tests must not assert these literally.

```toml
outcomes = ["locate"]

[model]
id = "navigator"
version = 1
class = "decision-table"

[tags.status]
provenance = "observed"
kind = "enum"
domain = ["Draft", "Final"]

[tags.recognized]
provenance = "recognized"
kind = "enum"
domain = ["locate"]

[[rule]]
id = "draft-propose"
[rule.match.status]
eq = "Draft"
[rule.match.recognized]
eq = "locate"
[rule.emit]
next = "propose"
```

```sh
intrastate flow resolve --model navigator.toml --outcome locate \
  --tag status=Draft --as=json
# → {"type":"ok","data":{"rule":"draft-propose","emit":{"next":"propose"},
#    "next":{},"writes":{},"clear":[],"owned":{},"readers":[],…}}
```

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Model class declaration | This RDR | Introduced | C1 |
| Class-conditioned rule shape | This RDR (over `0002:C4`) | Introduced | C2 |
| `[rule.emit]` on the row and in the dump | This RDR (over `0002:C19`) | Introduced | C3 |
| `emit` on the resolve payload | This RDR (over `0005:C1`) | Introduced | C4 |
| ∅-rooted reachability | This RDR (over `0006:D-reachability-relation`) | Introduced | C5 |
| Reader scoping that makes `--artifact` unrequired | Predecessor 0005 (`0005:C1`, `invokedReaders`) | Available | none |
| Coverage over the group's authored dimensions | Predecessor 0006 (`0006:C7`) | Available | none |
| Kernel plan over empty `RequiresOwned` | Predecessor 0002 (`0002:C14`) | Available | A3 |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Uninterpreted literal block on model data | `internal/table/model.go::Model.Metadata` (`[model.metadata]`) | model-level, `any`-typed, not dumped | Reuse the treatment, not the field — `[rule.emit]` is row-level and string-valued so it can be a dump column | C3 |
| Root of the reachability traversal | `internal/graphlint/reach.go::reach` | returns no nodes without `[initial]` | Extend: seed ∅ for the decision-table class | C5 |
| Missing-root finding | `internal/graphlint/analysis.go::checkDanglingEdge` | keys on `len(Initial)` | Extend: key on class | C5 |
| Ordinary-rule write obligation | `internal/table/normalize.go::normalizeRule` (`CatMalformedRuleShape`) | unconditional | Extend: condition on class | C2 |
| Reader demand set | `internal/cli/flow_exec.go::invokedReaders` | none — empty demand already invokes nothing | Reuse unchanged | C4 |
| Resolve payload | `internal/cli/flow_resolve.go::resolvePayload` | no answer field | Extend: `emit` | C4 |
| Dump vocabulary | `internal/table/dump.go::dumpColumns` | closed list | Extend: `emit` | C3 |

### Decision Rationale

The contract question is *whether a model with zero owned tags is a legal
class and what resolve answers over it*. Four answers were scored:

| Criterion | A: mandatory owned state, documented dummy-tag convention | B: class inferred from the empty owned set | **C: class declared in `[model]`, owned set must agree** | D: separate model kind + `flow decide` verb |
| --- | --- | --- | --- | --- |
| Correctness fit (author writes the model as what it is) | fails — the scratch artifact and dummy accessor are the defect | fits | fits | fits, but duplicates the kernel's exact-one at a second verb |
| Silent misclassification | n/a | a forgotten `provenance = "owned"` becomes a lint-clean table | refused at load both ways | refused (separate schema) |
| Prior-art alignment | matches pytransitions' auto-injected `'initial'` state — the workaround peers had to add an opt-out for (K3) | SCXML's targetless transition is per-transition, not per-machine; no peer infers machine class from absence | DMN treats the table as a declared artifact (model prior, A5); 0006:C18's "absent root is a defect" preserved verbatim | no peer splits select from apply at the verb (K2) |
| Overrides required | none | 0002:C4, 0002:C2/C19, 0005:C1, 0006:C18 + D-reachability | 0002:C2/C3/C4/C19 additively, 0005:C1 payload, 0006 D-reachability; **0006:C18 untouched** | a second grammar and verb; 0005:C1's four-verb closure |
| Reversibility | n/a | → declared later is a grammar addition | → inferred later is dropping a key | hard — two schemas to converge |
| Blast radius / cost | zero code; recurring consumer cost | one arm flip in lint + resolve payload | B plus one `[model]` key and one load check | new verb, payload, help, tests |

C wins on the misclassification and overrides rows: it is B plus one
declared line, and that line is what lets `0006:C18` stand unchanged and
makes the class reviewable in the artifact. A is the status quo the problem
statement rejects; D pays for a verb the kernel already provides.

For the answer, three shapes: row identity alone (already in the payload —
kept, insufficient alone because the consumer then needs an out-of-band map),
the PoC's pseudo-owned `next` tag (the defect itself), and an authored
literal `emit` block (chosen; the sibling design doc's own "same row → same
output" framing, K4).

Premortem: hardened (hardened) — critic verdict PASS, no switch forced;
nine mitigations folded: A1 pins the hole cell set and the size
suppression; A2 covers observed readers; A7/A8 become tests; A9 (the
`no_match` escape "otherwise" row), A10 (full-taxonomy silence), A11
(`flow next` over the class, re-verified after 0011) added; C3 fixes the
default dump columns, C4 `{}`-never-`null`; one class accessor for every
`len(...)==0` site. Ledger:
`docs/rdr/0010-stateless-decision-tables/evidence/propose-premortem/critic.md`.
Ground-sweep: clean (33 anchors) — 18 source anchors, 11 peer-element
quotes, 4 external prior-art quotes; two cosmetic corrections folded (the
`[dump]` fixture count in A4; the refusal categories for a declared
non-owned `[initial]` key in C2). Ledger:
`docs/rdr/0010-stateless-decision-tables/evidence/grounding-sweep/sweep.md`.
Joint-check: clear (1 peer). Open peers at depth 1 with Status Draft/Final:
cli/0011 only (0001–0009 are `Implemented`); 0011 shares none of this RDR's
modify-anchors (`invokedReaders`, `runReaders`, `resolvePayload`,
`normalizeRule`, `checkDanglingEdge`, `reach`, `checkGroups`, `dumpColumns`)
nor its contract literals (`class`, `decision-table`, `state-machine`,
`emit`, `malformed model declaration`, `graph-dangling-edge`,
`flow-artifact-missing`), and this RDR deliberately leaves `flow next` to
it. Absence arm: the refusals this RDR converts to acceptances — the
no-write-block `malformed rule shape` arm and the missing-root
`graph-dangling-edge` — are relied on by the `Implemented` predecessors
0002 (C4) and 0006 (C18, S15), which are named under `Overrides` and
conditioned on class here rather than removed; no Final peer exists to
fire on, and the predecessor coupling rides to 7.1 through `Overrides`.
Bridge sub-check: n/a — no sibling plan retires a surface this plan
introduces, and no plan here retires one of 0011's.

## Alternatives Considered

### Alternative 1: Keep owned state mandatory; document the dummy-tag convention

**Description**: Leave 0002/0005/0006 as they are and publish the workaround:
declare one owned enum, an `[initial]`, a read/write accessor pair on a
scratch role, and seed the scratch artifact before the first resolve.

**Pros**:

- No code change; no override of any locked clause.
- Peer precedent: pytransitions injects a default `'initial'` state when
  none is given.

**Cons**:

- The model lies about itself — a reviewer reads a state machine that has
  no state.
- The scratch artifact's format is consumer-invented and undocumented
  (intrastate has no artifact discovery by design), so every consumer
  reinvents it.
- The write-only owned tag is exactly the shape `0002:C2`'s accessor arm
  calls "not merely unbound but unusable".

**Reason for rejection**: it is the problem statement, restated as policy.

### Alternative 2: Infer the class from the empty owned set

**Description**: No new declaration; a model declaring zero owned tags *is* a
decision table, and every machine-only obligation is conditioned on
`len(owned) > 0`.

**Pros**:

- One key fewer in the grammar; nothing to keep in agreement.
- Identical lint and resolve behaviour to the chosen approach once the
  class is known.

**Cons**:

- The class becomes a property of what the author *omitted*. A state
  machine whose author mis-declares its state tag as `observed` silently
  becomes a decision table, lints clean, and resolves with empty writes —
  the "empty reachable set read as clean" failure `0006:C18` was written to
  forbid.
- Overrides `0006:C18` itself, where the chosen approach leaves it verbatim.
- No adjacent path infers from absence: `Row.Kind` is inferred from the
  presence of `escape`.

**Reason for rejection**: one declared line buys a load-time refusal in both
directions and keeps 0006's root rule intact.

### Alternative 3: A separate model kind with its own `flow decide` verb

**Description**: A distinct schema (`[decision]` or a separate file kind) and
a fifth verb that evaluates it, leaving `flow resolve` machine-only.

**Pros**:

- No conditioning inside the existing grammar; the machine grammar never
  learns about tables.

**Cons**:

- The kernel's exact-one selection already *is* the decision-table hit
  policy; a second verb re-exposes it with a second payload and help text.
- `0005:C1` closes the skill-integration surface at four verbs; a fifth is a
  producer change on a locked contract for no new semantics.
- Two schemas share tags, contexts, guards, gates, escape rows, and lint —
  everything but the write block — and would drift.

**Reason for rejection**: it duplicates a seam to avoid one conditional.

### Briefly Rejected

- **Answer = row identity only**: already carried as `rule`; alone it forces
  every consumer to keep an out-of-band rule→action map beside the model.
- **Answer = pseudo-owned `next` tag (the PoC's shape)**: is the defect.
- **`any`-typed emit values**: `[model.metadata]` precedent, but a dump
  column needs a deterministic rendering; raised as A6, widenable later.
- **`emit` in the `flow next` candidate preview**: readable without
  evaluation, but `flow next` is cli/0011's seam this cycle; deferred as a
  rider for the record that next owns the preview.
- **`--outcome` optional when the model declares one outcome**: applies to
  both classes; parked as a plain kata at triage.
- **A new lint finding announcing the class**: the class is reviewable in
  `[model]`; a finding would be noise on every clean table.

## Context

### Background

Discovered while authoring a consumer decision model (rdr#tmxk, under the
rdr#thsc umbrella): the model's only owned tag was the write target, and
`flow resolve --artifact role=<nonexistent>` refused with
`flow-owned-state-unavailable`. Seeding a scratch artifact unblocked it, but
only by inventing state the model does not have. The kata-triage comment on
intrastate#zdat frames the fork: (a) keep owned state mandatory and document
the dummy-tag/scratch-artifact convention; (b) make owned tags optional so a
model with none is a decision table; (c) a separate model kind/verb (e.g.
`decide`) rather than inference from zero owned tags. Triage also parked the
`--outcome` single-outcome default as a separable plain kata (it applies to
stateful models too) and ruled this seed is not a facet of intrastate#1mv1.

Constraints: RDRs are never amended; a decision change lands as this new
RDR that supersedes named clauses. Keep intrastate generic. `make check`
gates every change.

### Technical Environment

Go CLI (`bin/intrastate`; `make check` = fmt/vet/lint/build/graph-lint/test).
Related components: `internal/table` (rule-shape normalization, RDR 0002),
`internal/graphlint` (graph authority, RDR 0006), `internal/cli/flow_exec.go`
(resolve demand set, RDR 0005), `internal/resolve` kernel `Row` and planned
owned-tag `Writes` (RDR 0001/0007). Design history: `docs/rdr/0001–0009`,
`docs/jdr/0001`. Conventions in `AGENTS.md` (respond gateway, CLIError codes,
SilenceUsage).

## Research Findings

### Investigation

Read first, before enumerating (research/prior-art.md): the decision-table
artifact is defined by combination coverage (`developer-testing.pdf` p.142),
peer state-machine engines admit a transition that changes no state
(stateless `InternalTransition`; SCXML "transition without target"), and
peers *require* an initial state, with pytransitions silently injecting a
default one — the dummy-state workaround as prior art. ⚠ no prior-art
coverage for decision-table hit policies (DMN) in the corpora; that claim is
model prior and lives in A5. In-repo: `internal/graphlint/reach.go::reach`
returns no nodes without a root and `internal/graphlint/groups.go::checkGroups`
runs coverage only for reachable groups, so suppressing the missing-root
finding alone would make coverage vacuous — the class needs a root defined;
`internal/graphlint/analysis.go::nodeSatisfiesMatch` reads only owned atoms,
so the ∅ node satisfies every decision-table context;
`internal/cli/flow_exec.go::invokedReaders` already invokes nothing when no
row demands an owned key; `internal/table/model.go::Model.Metadata` is the
uninterpreted-literal precedent; and `internal/table/dump.go::dumpColumns`
is closed, so an emit field is a dump column. Sibling-path check for a class
discriminator: `Row.Kind` is inferred from the presence of `escape`
(`internal/table/normalize.go::normalizeRule`); searched, no path infers a
class from an absence.

### Key Discoveries

- **Documented** — 0006's reachability filter, not its missing-root finding,
  is what would make coverage vacuous over a rootless model
  (`internal/graphlint/groups.go::checkGroups`).
- **Documented** — `nodeSatisfiesMatch` consults owned atoms only, so a
  single ∅ node reaches every owned-atom-free context.
- **Documented** — the seed's "0002 write-only-owned-tag clause" needs no
  override: a decision table has no owned tags at all, so the clause is
  unreachable, not contradicted.
- **Documented** — `--artifact` becomes unrequired with no CLI change:
  `invokedReaders` demands nothing, and `runReaders` checks bindings for
  invoked names only.
- **Verified** (PoC, 2026-08-26, scratch-tag form) — a partial 512-cell
  table reports a 208-cell coverage hole; A1 is whether it reproduces
  without the scratch tag.
- **Assumed** — DMN "unique" hit policy ≡ kernel exact-one (A5).

## Trade-offs

### Consequences

- Positive: a decision-only model is authored as one, with no accessor
  surface, no scratch artifact, and a reviewable `class` line.
- Positive: every 0006 exhaustiveness guarantee applies to decision tables
  through the existing group machinery; no second coverage path.
- Positive: `0006:C18` and every state-machine behaviour are untouched;
  existing models load, lint, and resolve byte-identically except for the
  new `emit` payload field and dump column.
- Negative: one conditional (`class`) enters the loader, the rule-shape
  check, and the lint root — three sites that must agree.
- Negative: `emit` widens the dump vocabulary, so every explicit `[dump]`
  column list must name it (A4).
- Negative: the answer is string-valued; structured answers wait for a
  widening.

### Risks and Mitigations

- **Risk**: the three class-keyed sites drift (loader admits what lint
  refuses).
  **Mitigation**: the agreement check lives in the loader only (C1); lint
  and normalize read the loaded class, never the owned set.
- **Risk**: coverage over a decision table is vacuous through a path the
  ∅ root does not reach (e.g. the product bound or node ceiling short-
  circuits on a one-node graph).
  **Mitigation**: A1's spike asserts the PoC's exact finding set; the MVV
  requires a *positive* coverage finding before the clean run.
- **Risk**: a consumer treats `emit` keys as tags and expects them matched.
  **Mitigation**: C3 forbids it and the loader never declares them; a
  `[rule.match.<emit-key>]` refuses as `unknown tag`.
- **Risk**: an old binary meets a new model.
  **Mitigation**: strict decoding refuses the unknown `[model]` key (A8).
- **Risk**: the single group of a large decision table trips the product
  bound or the unprojectable-size suppression and its claim is withheld
  silently.
  **Mitigation**: A1's spike asserts the hole cell set; a withheld claim
  stays a visible `withheld` finding (`0006:C7`/`C12`), never exit 0.
- **Risk**: authors write an overlapping catch-all ordinary row as the
  "otherwise" and hit `flow-ambiguous-match`.
  **Mitigation**: A9 — the idiom is an escape row rescuing `no_match` with
  an emit; documented beside the class in the authoring docs (Phase 4).

### Failure Modes

- **Visible**: class/owned-set disagreement → `flow-model-invalid` with a
  `malformed model declaration` finding naming the class and count; a write
  block on a decision-table row → `write to non-owned tag` / `unknown tag`;
  a `[dump]` list omitting `emit` → `malformed dump declaration`.
- **Silent (guarded)**: a decision table with no coverage findings when it
  should have them — the ∅-root regression. Diagnosis: `intrastate lint
  --as=json` over a deliberately partial table must report a coverage
  finding; the MVV pins it and the graph-lint fixture set carries a
  decision-table negative control.
- **Recovery**: drop `class` (the model is refused as a rootless machine,
  as today) or add `provenance = "owned"` state and an `[initial]` to make
  it a machine; nothing persistent is written by this RDR.
- **Diagnosis path**: `intrastate dump` shows the `emit` column and empty
  `next`/`writes`; `flow resolve --as=json` shows `readers: []` and
  `owned: {}` — a non-empty `readers` over a declared decision table means
  the class check did not run.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified (A1 spike first — it decides whether
  C5 needs more than the ∅ root)
- [ ] Peer cli/0011 proposed, so `flow next` ownership is settled before
  Phase 3 touches `internal/cli`

### Minimum Viable Validation

1. Author a decision-table model: `class = "decision-table"`, two observed
   enum dimensions of two values each, one recognized outcome, three
   ordinary rules each carrying `[rule.emit]`, no `[initial]`, no
   accessors.
2. `intrastate lint --model m.toml --as=json` → exit 2 with one coverage
   finding naming the uncovered cell (a *positive* finding, proving
   coverage ran with no root declared).
3. Add the fourth rule; lint → exit 0, findings exactly `[]` — no code from
   the 0006 taxonomy present (A10). Variant: replace the fourth rule with an
   escape row rescuing `no_match` carrying `[rule.emit]` → exit 0 with only
   the `graph-coverage-closed-by-escape` advisory (A9).
4. `intrastate flow resolve --model m.toml --outcome <o> --tag a=x --tag
   b=y --as=json` with no `--artifact` → exit 0; `data.rule` is the fourth
   rule's id, `data.emit` equals its authored block, `data.next`,
   `data.writes`, `data.owned` are empty, `data.readers` is `[]`.
5. `intrastate dump --model m.toml` renders the `emit` column for every row
   with no `[dump]` declared; `flow next --model m.toml --as=json` with no
   `--artifact` → exit 0, every candidate with empty `required` (A11).
6. Negative controls: the same model with `class` omitted → lint
   `graph-dangling-edge` naming the absent `[initial]` (0006:C18
   unchanged); with one `provenance = "owned"` tag added → load refusal
   `malformed model declaration`; `models/rdr.toml` lints and resolves as
   before.

End-state: a stateless model lints for coverage and resolves to `rule` +
`emit` with no owned state, no accessor, and no artifact, while every
state-machine model behaves as it did.

### Phase 1: Grammar and Load

Intent: make the class a declared, checked fact and let a row carry an
answer — `class` in `sourceModel`, `emit` in `sourceRule`, the agreement
check, the class-conditioned rule-shape arm, `Row.Emit`, and the `emit` dump
column, with every `[dump]`-carrying fixture updated.

### Phase 2: Lint Root

Intent: root the decision-table class at ∅ so 0006's invariants run
unchanged — `reach` seeds by class, `checkDanglingEdge` keys its root arm on
class, and the graph-lint fixtures gain a decision-table pair (partial →
coverage finding; complete → clean) plus the class-omitted control.

### Phase 3: Resolve Payload

Intent: return the answer — `emit` on `resolvePayload`, joined by rule id
after selection, rendered in text mode; the output contract doc gains the
field and a decision-table invocation.

### Phase 4: Model and Docs

Intent: land the motivating shape without consumer knowledge — a generic
decision-table fixture under `internal/*/testdata`, `docs/cli-output-contract.md`
and the model authoring docs updated; `models/rdr.toml` untouched.

## Validation

### Testing Strategy

[Required — never omit. Test scenarios and coverage goals — what to test and
what constitutes "done." For non-functional concerns
(performance, security): state measurement strategy,
not estimates.]

1. **Scenario**: [Description]
   **Expected**: [Result]

## Finalization Gate

> Complete each item with a written response in
> `{ARTIFACT_DIR}/gate.md` before marking this RDR as
> **Final**. Written responses prevent rubber-stamping
> and produce a review record.
>
> First run the mechanical pre-sweep
> (`prompts/gate/tooling-pass.md`): TEMPLATE section
> coverage, Method-label vocabulary, `Source Search`
> self-reference, `Docs Only` on load-bearing claims. It
> catches what the review rounds disturbed; resolve any
> BLOCK before the written responses.
>
> At lock, replace this section's body with the
> one-line pointer to gate.md — responses are never
> inlined. The sub-sections below spec gate.md's
> content.

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
past the lens battery undetected. A `Transient`-marked
contract with a named deleting sibling and schedule is a
recorded lifespan disposition, not an under-sized
Profile — do not count it when re-deriving. Also confirm form:
value + one clause naming the contract(s); strip any
matrix/provenance prose left from the template or Seed
(it belongs in the template comment, not the instance).]

## References

- RDR 0002 `0002:C2`, `0002:C3`, `0002:C4`, `0002:C14`, `0002:C15`,
  `0002:C19` — the grammar, rule shape, derived owned set, and dump
  vocabulary this RDR conditions or extends.
- RDR 0005 `0005:C1` — reader scoping, the resolve payload, the four-verb
  closure.
- RDR 0006 `0006:C7`, `0006:C18`, `0006:D-reachability-relation` — coverage
  over authored dimensions, the missing-root rule, the rooted traversal.
- Source reviewed: `internal/table/normalize.go`, `internal/table/load.go`,
  `internal/table/source.go`, `internal/table/dump.go`,
  `internal/table/model.go`, `internal/graphlint/reach.go`,
  `internal/graphlint/analysis.go`, `internal/graphlint/groups.go`,
  `internal/graphlint/engine.go`, `internal/cli/flow_exec.go`,
  `internal/cli/flow_resolve.go`, `internal/resolve/resolve.go`.
- Prior art (evidence/research/prior-art.md): `developer-testing.pdf` p.142
  (decision tables); `../state-machines/study/stateless/README.md` §Internal
  transitions; `../state-machines/repos/scxmlcc/doc/user-manual.md`
  §transition; `../state-machines/study/transitions/README.md` (initial
  state; default `'initial'` injection); `../state-machines/MODEL-transition.md`
  §3.
- intrastate#zdat (seed and tracker), rdr#tmxk (motivating consumer model),
  rdr#thsc (umbrella), intrastate#1mv1 / cli/0011 (`flow next`, cross-cited).
