# Recommendation 0021: Lint's normalized-graph export

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-08-28
- **Status**: Draft [revised from Final 2026-09-12; re-verify A4 @refine — 0029:C4 obliges tier assignments for `--emit`, `graph-export-too-large` and C2's field spellings that this record never makes]
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
- **Type**: Feature
- **Profile**: large — one contract, the export of the normalized
  graph, stated as five clauses (C1–C5) of one seam; user-facing yes;
  locks the `intrastate.graph/1` document's field list and marker
  (its DOT projection is a documentation rendering, styling
  non-normative).
- **Priority**: Low
- **Related Issues**: kata `intrastate#jjkh` (1602); kata `4hps`
  (terminal-reachability invariant — now RDR 0022, a potential
  consumer of this export, neither blocking the other)
- **Predecessors**: 0002-transition-table-as-reviewable-data,
  0005-skill-integration-cli-contract,
  0006-graph-lint-authority-and-guarantees
- **Seam Lineage**: no prior accretion

## Problem Statement

A reviewer or a CI pipeline wants to see the graph lint actually
certifies — review it as a diagram, diff it across commits, or hand it
to an external formal checker. `intrastate lint` already computes the
normalized row set, the reachability relation from `[initial]`, and
the terminal set (`internal/graphlint/analysis.go`, `reach.go`) but
exposes none of it: pass/fail findings are the only output. Prior art
assigned intrastate the role of spec/validator ("spec the legal graph
once, verify no illegal edge — never a runtime"), and a hand-written
TLA+/Alloy spec was rejected in review as a second artifact that
drifts from the source of truth; a first-class export closes that gap
by letting any formal model be *derived* from the code instead.

The decisions: the surface shape (`--emit=json|dot` on the lint verb
vs a sibling verb); JSON schema ownership relative to RDR 0002's dump
identity/ordering; and how a non-JSON DOT rendering coexists with
RDR 0005's one-envelope-per-invocation `--as` contract (raw stream,
string-in-data, or file target) — all under two invariants: emission
is byte-for-byte deterministic for the same model (replay-safe,
diffable in CI), and emission can never alter lint's verdict or
finding set. The payload covers what lint already sees: normalized
rows (guards after `all`/`unless` combination), declared tags and
domains, initial state, terminal set, reachability edges, and row
groups. No prior record decides a lint export surface: RDR 0006 emits
only findings, RDR 0005 fixes one JSON envelope per invocation,
RDR 0002 owns only row-dump ordering.

## Critical Assumptions

- **A1 The respond gateway's `TextLiner` seam prints a multi-line
  `Data` string verbatim on stdout in text mode, with no decoration,
  reordering, or trailing content beyond one final newline — including
  under provoked notes/warnings, which stay on stderr.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `internal/cli/respond/respond.go::TextLiner` exists
    and `OK` prints it via `Fprintln`; the spike byte-compared
    `stdout == document + "\n"` for a multi-kilobyte multi-line
    payload with advisories provoked (premortem P-5/P-15) —
    byte-identical in both arms, notes/warnings confined to stderr
    (`evidence/spikes/a1-textliner.md`). Normative fixture F1:
    text-mode stdout is exactly `document + "\n"`, stderr empty
    absent advisories.
  - **If wrong**: the bare-document text mode gains stray bytes, CI
    diffs of the export break, and C5 needs a gateway extension
    instead of a ride on the existing seam.
- **A2 The reachable node set and its successor edges are a
  deterministic function of the model value, and the edge list is
  recoverable after the fixpoint by re-running `successorsOf` over the
  final nodes with subsumption mapping (`indexOf`), yielding EXACTLY
  the edges the traversal itself took — over-approximation by merged
  re-run is the premortem's central defect (P-1), so equality, not
  plausibility, is the bar.**
  - **Status**: Verified — with a recorded coupling: holds because
    `indexOf` uses `subsumes`, which makes `reach()`'s in-place
    widening arm unreachable (exhaustive lemma + zero-event audit,
    `evidence/spikes/a2-edge-recovery.md`)
  - **Method**: Spike
  - **Evidence**: `internal/graphlint/reach.go::reach` iterates slices,
    not maps, and `successorsOf` appends in first-seen row order. The
    spike showed (a) repeated runs byte-identical including under
    `GODEBUG=randmapiter=1`, and (b) a DIFFERENTIAL test — an observer
    recording edges during the traversal vs the post-hoc recovery —
    giving EXACT multiset equality over synthetic fixtures A–I and all
    36 authored models (`EQUAL=36, NOT-EQUAL=0`), with a negative
    control (injected phantom/dropped edge) and a positive control (a
    weakened `indexOf` made the arm live and diverge by 2 edges), so
    the harness detects divergence rather than defaulting green. The
    coupling: `subsumes` requires an identical key set plus value
    coverage, so `joinNodes(nodes[j], next)` is always key-equal to
    `nodes[j]` and the `merged.key() == nodes[j].key()` guard fires
    first (`a2-edge-recovery.md:59-64`). The property is pinned in
    CODE, not restated here — kata `yybx` ships the behaviour
    regression test and the stale-comment fix in isolation from this
    RDR; once landed, that test is the property's home. Graded
    `constraint` (answer fixed by shipped code, `5e57b79`); it does
    not gate the lock.
  - **If wrong**: the named fallback is the in-traversal edge observer
    (record edges as `reach` takes them); that arm is admissible
    because C4's neutrality oracle — lint byte-identical with and
    without the export code — is stated mechanism-independently and
    the MVV asserts it either way (P-14).
- **A3 The normalized model value carries everything the document
  needs — declarations with finite domains (via
  `guard.AssignmentCount`), `[initial]`, `terminal`, normalized rows
  with blocked atoms, and `guard.Groups` — with no re-parse of the
  authored TOML.**
  - **Status**: Verified (declarations, tags, initial, terminal, rows,
    groups, nodes); the edge relation is derived in the export path as
    a pure function of the model value — see A2
  - **Method**: Source Search
  - **Evidence**: `internal/table/model.go` field walk against the C2
    field list (`Model.ID`/`Class`/`Tags`/`Initial`/`Terminal`/`Rows`),
    `internal/guard/declaration.go::AssignmentCount` for finite
    domains, `internal/guard/product.go::Groups` for the partition,
    and `internal/graphlint/reach.go::Reach` for nodes — every one a
    pure function of the model value with no TOML re-parse
    (`evidence/research/a3-a4.md`). The two `guard` functions are read
    as-is: this record changes nothing in that package, whose
    observation surface is RDR 0013's to evolve. The edge relation is
    NOT carried by the model value: `reach()` computes nodes only and
    `successorsOf` reads `row.RuleID` but discards it, so it is
    recovered in the export path per A2 (the rule id is on
    `table.Row.RuleID`, no re-parse). C2 predicates *carries* of the
    document, which does carry edges.
  - **If wrong**: the schema shrinks, or a loader extension becomes a
    prerequisite and the blast radius grows past this RDR.
- **A4 A new root export verb requires no amendment to RDR 0005's
  envelope contract.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: `0005:C1` (0005 is Implemented, gate not stale)
    scopes itself in its opening sentence — "The CLI MUST expose one
    command group for skill integration with these verbs: next,
    resolve, read-state, and set-state" — so every MUST-clause that
    follows binds the `flow` skill-integration group and has no
    jurisdiction over a root `graph` verb. The carve-out sentence
    ("Other command groups (lint, dump, parse) are outside this
    contract and are owned by the RDR that names them") is a closed
    three-item enumeration that does not itself name `graph`; the
    conclusion rests on the scope sentence, not on reading `graph`
    into that list (`evidence/research/a3-a4.md`).
  - **If wrong**: the surface must be renegotiated at the envelope
    home before Phase 2 can land.
- **A5 Marshaling the document through the shared non-HTML-escaping
  encoder is byte-stable: struct field order is fixed, every sequence
  is pre-sorted, and no Go map reaches the wire.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `internal/cli/clierr/clierr.go:174::WriteJSONLine` is
    the one encoder (`0005:C1`); the spike double-emitted a fixture
    byte-identically and held one sha256 across 25 separate processes,
    15 of them under `GODEBUG=randmapiter=1`
    (`evidence/spikes/a5-encoder-stability.md`). Normative fixture F2
    is the ENCODER PROPERTY only — double-emit byte identity, process
    stability, `<`/`>`/`&` unescaped, compact with one trailing `\n`,
    `schema` first, `[]` for empties — NOT a golden hash and NOT a
    canonical `intrastate.graph/1` emission: the spike improvised field
    spellings that differ from C2's normative list (and from A6's
    spike), so pinning that sha256 would fix C2's spellings by
    accident. The Phase 1 golden minted by the real exporter is the
    byte fixture for Testing Strategy scenario 2. Latent risk, not a
    falsifier: four call sites build ad hoc `SetEscapeHTML(false)`
    encoders instead of routing through `WriteJSONLine`, so the
    implementation must reuse the shared encoder rather than add a
    fifth.
  - **If wrong**: C3's replay invariant fails and the document needs a
    custom marshaler with its own ordering proof.
- **A6 The DOT rendering is derivable from the exported document value
  alone — nodes, edges, and initial are enough — with no reach or
  analysis internals consulted.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: the spike renderer imports only stdlib — zero
    `graphlint`/analysis imports — and its DOT rendered via `dot -Tsvg`
    (graphviz 12.2.1, exit 0) for a canonical fixture and a hostile one
    carrying quotes, newlines, and non-ASCII; the emitted node/edge
    statements were an exact set match against the fixture's `reach`
    block (`evidence/spikes/a6-dot-derivability.md`). Normative fixture
    F3 is scoped to the node/edge SET and the abstraction-marker
    placement (header comment + graph `label`), not styling: the
    spike's label composition has a real escaping bug (the separator is
    injected as two bytes before `dotQuote` doubles the backslash),
    which is styling and outside the fixture, while the escaping ORDER
    (backslash, then quote, then newline) is sound and carries into
    implementation. The spike's field spellings are provisional
    placeholders (`abstraction_marker`, `terminal_satisfying`); C2's
    spellings govern.
  - **If wrong**: the renderer couples to graphlint internals and the
    two `--emit` arms stop being projections of one value.

- **A7 The `values` member of a `reach` node has one declared shape,
  and C2's field list fixes it.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/graphlint/reach.go::Node` carries `Values
    map[string][]string` — tag name to a sorted, deduplicated value
    list (`canonicalValues` is
    `slices.Compact(slices.Sorted(slices.Values(value)))`), intact at
    the point `reach`/`Reach` return nodes to a caller. So the OBJECT
    keyed by tag is the invention-free projection (map → object,
    `[]string` → array), which C2 fixes; joined `key=value` strings
    exist only inside `(Node).key`'s escaped fingerprint, for internal
    dedup identity, never as an exportable value.
  - **If wrong**: the wire shape of every exported node is
    underdetermined at lock; two implementers read one field list two
    ways and the golden fixture pins whichever shipped first.

- **A8 The exported `graphlint` surface C4 requires — the
  edges-carrying function that also returns the traversal's
  completeness — can be added without altering `reach()` or lint's
  path through it.**
  - **Status**: Verified — the addition is additive
  - **Method**: Source Search
  - **Evidence**: `evidence/spikes/a8-graphlint-surface.md`.
    Unexported `reach()` returns `(nodes []Node, complete bool)`
    (`internal/graphlint/reach.go::reach`) and exported `Reach`
    discards the bool (`nodes, _ := reach(m)`,
    `internal/graphlint/reach.go::Reach`). Lint never calls `Reach`:
    it reaches `reach()` through `Run`→`newAnalysis`
    (`internal/graphlint/analysis.go::newAnalysis`,
    `internal/graphlint/engine.go::Run`), and `successorsOf` already
    computes the edges transiently and discards them — so a NEW
    exported function beside `Reach` carrying nodes, edges, and
    completeness leaves `reach()`, `Reach`, and lint's call site
    untouched (C4's neutrality). Q3(c) already sanctions that function
    for edge recovery. The CONDITION is already observable —
    `internal/graphlint/taxonomy.go::CodeProductTooLarge` fires on
    exactly `complete == false`
    (`internal/graphlint/analysis.go::newAnalysis`) and rides exported
    `Report.Findings` (`internal/graphlint/engine.go::Run`), a faithful
    1:1 proxy with no divergent case found; only the raw typed bool is
    unexported, which is what the new function carries.
  - **If wrong**: C4's `graph-export-too-large` arm is unimplementable
    without widening `Reach`'s public signature — a change to a surface
    whose only other callers are tests, which C4's neutrality rule and
    A2's recorded coupling both bear on.

## Proposed Solution

### Approach

A new root verb, `intrastate graph`, exports the normalized graph lint
certifies as one deterministic document. `--emit json|dot` (default
`json`) selects the document format — the peer pattern is one exporter
surface with a format-selector flag, JSON and DOT as sibling output
types (`../state-machines/repos/state-machine-cat/README.md`,
`--output-type … dot|…|json|…`). `--as` keeps its existing envelope
meaning untouched: under `--as=text` stdout carries the selected
document verbatim (riding the gateway's existing
`respond.go::TextLiner` seam — the seam for a payload that is one
canonical string a caller pipes; its writer is the verb's payload
constructor, set per invocation); under `--as=json` stdout carries
exactly one `ok` envelope embedding the same document. The verb runs
load + normalize + reachability only: it never runs the lint
invariants, emits no findings, and exports any model that loads —
including one lint would refuse — so emission can never alter lint's
verdict and a failing model can still be inspected as a graph. The
JSON document is the machine-readable follow-up RDR 0002 seeded and
declined to define, so this RDR owns the schema without touching
0002's dump text (Key Discoveries).

### Technical Design

One new document-assembly component (package boundary settled at
Resolve; it sits beside the relation it exports) builds an export
value from the loaded `*table.Model` plus the reachability fixpoint
(`internal/graphlint/reach.go`, via a new exported edge-recovery
function beside `Reach` per
A2). The CLI verb is a thin arm over it: selection flags → load
(`table.LoadWithAdvisories`'s underlying load path — advisories are
ignored here; they are lint's advisory channel, not graph data) →
assemble → render (`--emit`) → respond gateway. The DOT renderer is a
pure function of the export value (A6). Extension point: the document
is schema-versioned, so later producers (a declared-emit vocabulary,
new invariant metadata) join additively without a second export
surface.

#### Normative Contracts

These five blocks are facets of the one contract this RDR owns — the
export surface and its document — stated separately by concern
(surface, document, determinism, neutrality, mode coexistence).

**C1**

```normative
SURFACE. A new root verb `graph` is registered beside `lint` — outside
the `flow` group, under `0005:C1`'s carve-out for command groups "owned
by the RDR that names them". Model selection mirrors `lint`'s arm set
— mirroring is scoped to the ARM SET and its CODE SPELLINGS, not to
lint's help/usage text, which each verb words for itself —
and codes verbatim (`internal/cli/lint.go::runLint`): exactly one of
`--model <path>` / `--flow <id>`; both → `flag-mutually-exclusive`;
`--flow` alone → `flag-invalid-value` (this build resolves no ids);
neither → `flag-required`; unreadable file → `model-unreadable`; load
failure → `model-invalid` with one findings[] entry per load category.
`--emit <format>` selects the document: `json` (default) or `dot`; any
other value → `flag-invalid-value` naming `emit`. `--emit` is NEW to
this verb — `lint` has no such flag, so mirroring supplies no position
for it and the order is fixed here: `--emit` validity is checked with
the argument-shaped arms, AFTER the `--model`/`--flow` selection arms
and BEFORE any file I/O. So `graph --model <unreadable> --emit=xml`
refuses `flag-invalid-value` naming `emit`, never `model-unreadable`;
the request is wrong independent of the environment. This is the one
value-checked flag on the verb: `runLint` performs no enum-value check
at all (its `--flow`-alone arm emits `flag-invalid-value` but is a
presence arm, and the only load-time enum — `class` — is inside the
loader, i.e. inside the `model-invalid` arm). RunE starts with
`respond.ValidateMode`, success routes through `respond.OK`, failure
through `respond.Fail`; exit codes are the existing 0/2 mapping — no
new exit group (`0005:C1`: exit 3 is environment-not-consulted only,
and an unreadable model file is a wrong request, so it stays exit 2).
```

**C2**

```normative
DOCUMENT. The JSON document is a documentation artifact, versioned by
a required leading `schema` field, initial value `intrastate.graph/1`;
evolution within `/1` is additive (a consumer ignoring unknown fields
keeps working), and an incompatible change bumps the marker under the
`0.x` promise RDR 0029 governs. Exact field spellings are normative AS
SPELLED HERE — the Illustrative Code is an exhibit that disclaims
literal assertion and shows only some members, so it binds nothing:
`schema`; `model` (`table.Model.ID`, the AUTHORED `[model] id`, never
the `--model <path>` argument or any path-derived string); `class`;
`tags[{name, provenance, kind, required, single_valued, domain}]`
(`domain` present exactly when `guard.AssignmentCount` reports it
finite); `initial`; `terminal` (the declared predicate sets, carried
as declared); `rows[{identity, source, kind, outcome, atoms, next,
writes, requires_owned, gate, escape, emit}]` — RDR 0002's dump field
list in its canonical row order, atoms (each `{key, operator,
literal[], block}`) in 0002's canonical atom order, lowercase
snake_case (`internal/table/dump.go::dumpColumns`, a closed 11-member
list with `emit` appended last per `0010:C3`); `groups[{context,
rules}]`; and `reach{abstraction, nodes[{id, values}], edges[{from,
to, rule}]}` — the merged fixpoint relation, `id` the canonical node
key, `values` an OBJECT keyed by tag name whose every value is that
tag's sorted, deduplicated value ARRAY (the shape `reach.go::Node`'s
`Values map[string][]string` projects without invention), nodes sorted
by node key and edges by (from, to, rule) so construction order is
unobservable. A tag with no finite declared domain carries the single
value `<opaque>` (`reach.go::OpaqueValue`) in its array, passed
through verbatim; this record attaches no meaning to that spelling.
Every declared collection renders as an empty JSON array `[]` (or
object `{}`) when it has no members — never `null`; an optional member
that does not apply is ABSENT, its key omitted, never `null`. The
`reach` block carries a REQUIRED marker, `abstraction` with the token
value `declared-over-approximation`, stating the relation is the
DECLARED over-approximation, not the runtime — merged nodes,
guard/observed atoms unpruned (`reach.go::Reach` doc); the schema docs
state the soundness rule in one sentence: universal claims ("no path
does X") proved over this relation hold at runtime; existence claims
("some path reaches X") may be spurious. The document carries NO
verdict or finding field — an export is never a lint pass — and NO
per-node terminal marking: which merged nodes satisfy a `terminal`
predicate set is the dead-end quantifier RDR 0015 owns (JDR 0001
§JD-23), and this record declares no evaluator of its own; the
declared sets travel in the document for a consumer to evaluate.
Set-valued members are JSON arrays, closing
`0002:§round-trip-inverse-invariants`'s lossy set-literal rendering
for this document; the document is NOT a model source and no
export→load inverse is claimed. The DOT document renders the same
value: one node per reachability node, one edge per reachability edge
labeled with its rule id, the initial node marked, and the abstraction
marker rendered in the graph header comment/label so the diagram
carries it too (premortem P-7); its node/edge SET and the marker are
normative, its styling/attributes are not.
```

**C3**

```normative
DETERMINISM. For one model input and one build, emission is
byte-for-byte identical across invocations, in every `--emit` and
`--as` combination: every sequence on the wire is pre-sorted by C2's
orders before marshaling, no Go map iteration reaches the wire, and
JSON is emitted through the one shared non-HTML-escaping encoder
(`clierr.WriteJSONLine`, per `0005:C1`'s one-encoder rule). The
promise is scoped to (model, build), not across builds (premortem
P-3): a build bump may change the JSON only by C2's additive rule and
may re-baseline DOT diffs freely, since DOT styling is non-normative.
```

**C4**

```normative
NEUTRALITY. The export runs load, normalization, grouping, and the
reachability traversal only, reaching the traversal through the same
`graphlint` entry surface lint uses (`0006:§approach`'s
same-command/request-builder precedent), so verb/lint drift is
structural, not disciplinary. It
MUST NOT run the lint invariants, MUST NOT emit findings, and MUST
NOT alter any input it shares with lint: `intrastate lint`'s verdict,
finding set, and bytes are identical with and without the export code
present (the MVV asserts this; what `intrastate lint` promises as a
repository gate is RDR 0014's contract, which this record only leaves
unchanged), and the oracle is
MECHANISM-INDEPENDENT — it binds equally if Resolve picks A2's
in-traversal edge observer (premortem P-14). The verb succeeds for
ANY model that loads, including a model lint refuses; the
model-loads-but-lint-refuses case is a dedicated fixture with an
asserted, defined document — never whatever the traversal happens to
do (premortem P-6).
When the traversal is incomplete under the published node ceiling
(`graphlint.NodeCeiling()`), the verb refuses with the scalar code
`graph-export-too-large` (GroupUserEnv, exit 2) naming the ceiling and
the narrow-a-domain remedy — never a partial document, because a
partial graph diffs as a graph change. Completeness MUST reach the verb
through an EXPORTED `graphlint` surface: today's `Reach` discards the
private `reach`'s `complete` bool (`nodes, _ := reach(m)`), so the
edges-carrying function Q3(c) already adds beside it returns
completeness with the nodes and edges. This is the same
package-internal traversal lint reaches, so the neutrality rule above
is unaffected. The ceiling this refusal names is `reach`'s node-count
completeness — the SAME bound `analysis.go`'s `checkNodeCeiling`
already fires on (both read one `reach()` and one `nodeCeiling`
constant, published as `graphlint.NodeCeiling()`), observed at two call
sites, which is the neutrality rule above rather than an exception to
it. The distinct bound this refusal does NOT involve is the
guard-product one, `graphlint.ProductBound()` (`guard.Bound()`); no
guard-product input reaches the export's refusal path.
```

**C5**

```normative
MODE COEXISTENCE. Under `--as=text` (the default) stdout carries the
selected document verbatim and nothing else — the payload satisfies
the gateway's `TextLiner` seam, so the DOT stream pipes to `dot` and
the JSON document diffs raw in CI; advisories stay on stderr. Under
`--as=json` stdout carries exactly one terminal `ok` envelope whose
`data` embeds the same document: the document object for `--emit
json`, and for `--emit dot` an object carrying the DOT text as the
single REQUIRED string member `dot` — spelled normatively here, so the
documented unwrap is exactly `jq -r .data.dot`. Both modes derive from
one export value (0005's two-modes agreement). All four `--as`×
`--emit` cells are defined — none refused, none dead — and a consumer
that parses stdout as an envelope MUST use `--as=json`: the text-mode
stream is the bare document by contract, exactly as `version`'s
`TextLine` is its identity string (premortem P-10).
The verb offers NO caller-controlled projection of its success
payload; if a later revision adds one it MUST conform to JDR 0002 §D1
(projection before respond.OK, echo-group-only, enforced partition,
always-keep core).
```

Determinacy: fired — C3 states byte-for-byte identity across
invocations in every `--emit`×`--as` cell, and C2 fixes exact field
spellings, sort orders, and empty/absent renderings that Scenario 2's
goldens pin. An under-specified cell would let two implementers ship
different bytes that both read as conforming, which is the trigger's
subject.

#### Load-Bearing Decisions

- **Identity** — a node is its canonical node key
  (`reach.go::(Node).key` — injective by escaping, ⇒ two exported
  nodes never collide); a row is RDR 0002's row identity; an edge is
  the `(from, to, rule)` triple; and the document's `model` member is
  `table.Model.ID` — the AUTHORED `[model] id` declaration
  (`internal/table/load.go` assigns it from the decoded `[model]`
  table; an absent id refuses the load), never the `--model <path>`
  argument or any path-derived string. Identity is therefore
  invocation-independent: the same model exported from two checkouts
  is byte-identical, which is what C3's (model, build) narrowing
  assumes and what the CI-diffability outcome rests on.
- **Wire / byte format** — C2's schema-versioned JSON document; exact
  field spellings normative per C2's list; empty declared collections
  render `[]`/`{}` and an inapplicable optional member is absent, never
  `null`; DOT styling explicitly non-normative.
- **Naming** — verb `graph`, flag `--emit`. Rejected: `dump` (0002's
  dump is a distinct, rows-only review rendering — reusing the name
  would merge two contracts), `export` (names the act, not the
  artifact), overloading `--as` for format selection (conflates
  envelope mode with document format and leaves no room for a third
  format).
- **Selection / predicate** — `--emit` selects the DOCUMENT, `--as`
  selects the ENVELOPE; the 2×2 composes with no refused or dead cell,
  which is what kept a second stream/exception out of the gateway.

#### Round-Trip / Inverse Invariants

- `json-decode ∘ export = value identity on every field C2 lists`:
  decoding the JSON document reconstructs, value-for-value, each of
  those fields — including exact set members, the recoverability
  0002's text dump deliberately declined
  (`0002:§round-trip-inverse-invariants`). The quantifier is C2's field
  list, not "whatever was exported", so the invariant constrains the
  emitter rather than restating itself.
- `export ∘ export = byte identity on any loadable model` (C3's replay
  form; the MVV asserts it as byte equality, not exit-code green).
- No `load ∘ export` inverse is claimed: the document is derived
  output, never a model source.

#### Illustrative Code

Illustrative only — tests must not assert these literally. The JSON
below is a PARTIAL document: the row and tag objects show a few of C2's
members, not all. C2 is the normative field list; this example never
narrows it, and the field spellings below are the normative ones C2
cites.

```sh
# CI: diff the certified graph across two commits.
intrastate graph --model flow.toml > graph.json

# Review: render the diagram.
intrastate graph --model flow.toml --emit dot | dot -Tsvg -o flow.svg

# Machine caller: same document, enveloped.
intrastate graph --model flow.toml --as=json | jq .data.schema
```

```json
{"schema":"intrastate.graph/1","model":"flow","class":"state-machine",
 "tags":[{"name":"stage","provenance":"owned","kind":"enum",
          "required":true,"single_valued":true,
          "domain":["draft","final"]}],
 "initial":[{"key":"stage","value":["draft"]}],
 "terminal":[[{"key":"stage","operator":"eq","literal":["final"],
               "block":"match"}]],
 "rows":[{"identity":"flow/lock","kind":"ordinary","outcome":"lock",
          "atoms":[{"key":"stage","operator":"eq","literal":["draft"],
                    "block":"match"}],
          "writes":[{"key":"stage","value":["final"]}]}],
 "groups":[{"context":"…","rules":["flow/lock"]}],
 "reach":{"abstraction":"declared-over-approximation",
          "nodes":[{"id":"stage=draft,;","values":{"stage":["draft"]}},
                   {"id":"stage=final,;","values":{"stage":["final"]}}],
          "edges":[{"from":"stage=draft,;","to":"stage=final,;",
                    "rule":"flow/lock"}]}}
```

#### Pre-Lock Mini-Checks

Cue-fired tables (Stage 5). Witnesses are spike output, not exporter
output — the exporter is unbuilt — and are cited by spike, never
promoted to goldens (`F2`: the A5 sha256 is the encoder property, not
the canonical `intrastate.graph/1` hash; `F3`: the A6 arm is scoped to
the node/edge SET and marker placement, not styling).

`authority` — input/decision × writer · readers · call sites · sibling arms · canonical

| Input / decision | Writer | Readers | Call sites | Sibling arms | Canonical |
| --- | --- | --- | --- | --- | --- |
| Reachability relation | `reach.go::reach` (private) | lint via `graphlint.Run`; export via the new sibling exporter | `analysis.go:46` (lint); the new exporter | `graphlint.Reach` (public wrapper, tests only) | `reach()` — both arms reach the one traversal (C4) |
| Edge list | the new exported function (post-hoc recovery over final nodes) | export only | the new exporter | A2's rejected in-traversal observer | post-hoc recovery; equality to the traversal is the bar (A2) |
| Empty collection on the wire | producer code choosing `[]T{}`, never nil | every document consumer | every C2 collection field | `null` (rejected) | `[]`/`{}`; optional member ABSENT (Q2 ruling) |
| Selection-arm refusal codes | `lint.go::runLint` | `graph` verb mirrors the arm set and code spellings verbatim; help text is the verb's own (C1) | C1's arm set | lint's own codes (`0006:C19`/`C20`) | lint — C1 mirrors, never redefines |
| Document format selection | `--emit` (this verb only) | export path | the `graph` verb | `--as` (envelope mode, `respond.go::FlagName`) | `--emit` selects DOCUMENT, `--as` selects ENVELOPE |

`oracle` — each MVV row × "fails if X is wrong because Y" + the negative control

| MVV step | Fails if… because | Negative control |
| --- | --- | --- |
| 2 — double emit byte-identical | map iteration or field order reaches the wire; byte compare diverges | A5 §5a: a map probe emitted under `randmapiter=1`, 25 processes — sorted, stable; A2 positive control: weakened `indexOf` diverged by 2 edges |
| 2 — carries every C2 field | a field is dropped or renamed; the field-presence assertion fails | A5 fixture exercises every C2 member incl. `empty_probe: []` |
| 3 — DOT set equals `reach` block | the DOT arm derives from analysis internals rather than the document; id sets diverge, or an identifier is mis-escaped into different bytes that still parse | A6 hostile fixture: quotes/backslash/newline/non-ASCII survive to `dot -Tsvg` exit 0 with exact set match; exit 0 alone is NOT the bar — each identifier is unescaped by inverting A6's escaping order and compared to its source (S6) |
| 4 — `jq .data` equals the text document | the two modes derive from different values | A1: text-mode stdout is exactly `document + "\n"`, one `Fprintln` |
| 5 — lint refuses identically | export code perturbs lint's verdict/finding set/bytes | byte-compare against a pre-change capture; oracle is mechanism-independent (C4) |

`fidelity` — operation × invariant + lossy exemptions

| Operation | Invariant | Witness / exemption |
| --- | --- | --- |
| `export ∘ export` | byte identity on any loadable model (RT2, C3) | A5: double-emit sha256 match; 25 fresh processes identical; A2: 2000 repeats → 1 digest |
| `json-decode ∘ export` | value identity on every C2 field incl. exact set members (RT1) | set-valued atoms are JSON arrays — closes `0002`'s lossy set-literal rendering |
| `load ∘ export` | **no inverse claimed** (RT3) | derived output, never a model source — declared exemption |
| DOT ← document | node/edge SET + marker equality; styling NON-normative | A6: exact set match, marker in header comment + graph `label`; escaping order backslash→quote→newline |
| across builds | JSON additive-only; DOT carries no cross-build promise | C3's stated narrowing to (model, build); cross-build golden makes the moment visible |

`disposition` — input class × exit · error · artifact · silent-vs-loud

| Input class | Exit | Error / code | Artifact | Silent or loud |
| --- | --- | --- | --- | --- |
| Model loads, lint clean | 0 | — | document on stdout | loud |
| Model loads, lint would refuse | 0 | — | ASSERTED document (defined content) | loud — the debugging half of the outcome (C4) |
| Traversal incomplete at node ceiling (`graphlint.NodeCeiling()`, via the exported completeness surface C4 requires) | 2 | `graph-export-too-large` (GroupUserEnv), names ceiling + narrow-a-domain remedy | **none** — never a partial document | loud (C4) |
| Both/neither `--model`/`--flow` | 2 | `flag-mutually-exclusive` / `flag-required` | none | loud (C1) |
| `--flow` alone | 2 | `flag-invalid-value` (no ids resolve this build) | none | loud |
| Unreadable file / load failure | 2 | `model-unreadable` / `model-invalid` (one findings[] entry per load category) | none | loud |
| Unknown `--emit` value | 2 | `flag-invalid-value` naming `emit` | none | loud |
| Declared collection with no members | 0 | — | `[]`/`{}` — never `null` | loud (Q2) |
| Inapplicable optional member | 0 | — | key ABSENT | silent by contract (Q2) |

`trace` — MVV walked stepwise × assertions in force × witness

| Step | Assertions in force | Witness |
| --- | --- | --- |
| 1 — author fixtures | C2 (field list) | spike fixtures exercise every C2 member |
| 2 — emit twice, compare | C3 byte identity × C2 field list × C5 text-mode bare document | A5 sha256 `dcb0259a…` both emits; 25 processes identical; `schema` first; compact + one trailing `\n` |
| 2 — parse as JSON | C2 `[]`-not-`null` × Q2 ABSENT rule | A5: `empty_slice_literal: []` vs `nil_slice_literal: null` — distinguishable on the wire |
| 3 — `--emit dot \| dot -Tsvg` | C2 DOT set normativity × A6 derivability × S6 hostile content | A6: exact node/edge set match, marker in header + label, exit 0 on hostile fixture, and each identifier unescapes (A6 order inverted) back to its source — exit 0 alone would pass a mis-escaped identifier |
| 4 — `--as=json \| jq .data` | C5 mode agreement × C3 determinism | A1: text stdout is exactly `document + "\n"`; both modes derive from one export value |
| 5 — lint neutrality | C4 MUST-NOT-alter × C1 shared arm codes | A2: 36/36 real models EQUAL; lint reaches `reach()` via `analysis.go:46` unchanged |
| — `reach.values` shape | C2 now fixes `values` as a tag-keyed object of value arrays | RESOLVED (was CONTRADICTION): `reach.go::Node.Values` is `map[string][]string`, so the OBJECT exhibit projects without invention and the A5 fixture's array-of-strings rendering was the wrong exhibit — that flattening lives only in `(Node).key`'s internal fingerprint. C2 pins the shape; A7 Verified. |

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Reachability edge list (nodes exist; edges are not returned today) | This RDR (extends `reach.go`) | Introduced | A2's recovery must equal the traversal's own edges |
| Bare-document text mode | Existing (`respond.go::TextLiner`) | Available | A1 verifies verbatim multi-line pass-through |
| Canonical row/atom orders | Predecessor (RDR 0002 normalize/dump) | Available | C2 cites, never re-sorts differently |
| DOT rendering | This RDR | Introduced | Pure function of the export value (A6) |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Row rendering | `internal/table/dump.go::Dump` | Text-only, set literals lossy by decision (`0002:§round-trip-inverse-invariants`) | Reuse the field list and order, not the renderer | C2 carries the dump vocabulary as structured JSON |
| Reachability | `internal/graphlint/reach.go::Reach` | Returns nodes only, no edges | Add a sibling exported function (`Reach` and `reach()` both untouched) | A2; no behavior change to lint's relation |
| Output gateway | `internal/cli/respond` | One terminal envelope under `--as=json` | Reuse (`TextLiner` + `OK`) | C5; no gateway exception, no second stream |
| Selection arms | `internal/cli/lint.go::runLint` | Codes are lint's own by contract (`0006:C19`, `0006:C20`) | Reuse the arm set and code spellings | C1 mirrors the arm set and code spellings verbatim; help text is not mirrored |

### Decision Rationale

Scored matrix (large profile), approaches × deciding criteria — one
clause per cell; ✓ favorable, ✗ disqualifying, ~ workable with cost:

| Criterion | O1 sibling verb `graph --emit` | O2 `--emit` on `lint` | O3 file target `--out` | O4 external generator over dump text |
| --- | --- | --- | --- | --- |
| Verdict neutrality provable | ✓ structural — lint untouched | ✗ shares one RunE with the verdict | ~ side-effect write inside a verify verb | ✓ out of process |
| Export of a failing model | ✓ load-success only | ✗ blocking model exits 2 with a CLIError that has no graph carrier | ✗ same entanglement | ✓ but only rows, no reachability |
| Envelope-contract fit | ✓ TextLiner + one envelope, no exception | ✗ second document on one stream breaks one-envelope | ~ receipt envelope, but two output places | ✓ n/a |
| CI diffability | ✓ raw stdout document | ~ must survive beside findings | ~ temp-file management in CI | ✗ dump text lossy (`0002:§round-trip-inverse-invariants`) |
| Prior-art alignment | ✓ smcat/`--output-type` stdout stream | ✗ no opened peer folds export into its verifier verdict | ~ smcat supports `-o`, but stdout is its documented pipe form | ✗ peers export structured, not scraped |
| Blast radius / reversibility | ✓ one new verb, deletable | ~ flags on the authoritative acceptance surface | ~ new I/O side-effect class | ✓ none, but outcome unmet |

O1 wins on the two disqualifying rows: only a surface that does not
share `lint`'s RunE can export a model lint refuses (the debugging
half of the user outcome), and only a document that IS the stdout
stream in text mode is diffable and pipeable without envelope
surgery. O2 fails both; O3 fails the failing-model row and adds an
I/O class for no gain over a shell redirect; O4 cannot carry the
reachability relation at all. Sibling-path check: the repo's only
existing format discriminator is the `--as` envelope mode
(`respond.go::FlagName`); searched `internal/cli` for an existing
document-format flag — none exists, so `--emit` is a new
discriminator, deliberately scoped to this one verb.

Premortem: hardened (hardened) — one fresh-context critic, queried
once (`evidence/propose-premortem/critic.md`, 16 findings): the
approach survived; the central mechanism finding (P-1, post-fixpoint
edge recovery) hardened A2 into an equality bar with a named
in-traversal-observer fallback, and P-3/P-6/P-7/P-10/P-14/P-16 folded
into C2–C5, Risks, and Failure Modes; the three citation findings
(P-2/P-11/P-12) resolved against the already-opened sources.

Ground-sweep: clean (20 anchors) — one fresh-context checker over the
cited anchors only; every Go symbol, refusal code, dump column, peer
passage (0005:C1, `0002:§round-trip-inverse-invariants`, `0006:C19`
via `lint.go`'s comment), JDR 0002 §D1, and both `../state-machines`
passages CONFIRMED on `main`.

Joint-check: clear (12 peers) — open peers 0012–0020, 0022–0024
checked on both arms: modify-anchors (`reach.go::Reach` extension,
new `graph` verb/root registration) and contract literals (`--emit`,
`intrastate.graph/1`, `graph-export-too-large`, the mirrored C1
codes) share no undecided contract. The three peers whose anchors or
literals this record also touches are cited where they touch, each a
read-only use with the decision left to its owner: RDR 0013 owns the
`guard` observation surface A3 reads; RDR 0014 owns what `intrastate
lint` promises as a repository gate, which C4 only holds unchanged;
RDR 0015 owns terminal satisfaction over merged nodes (JDR 0001
§JD-23), on which C2 declares nothing — the export publishes the
merged relation and the declared `terminal` sets, marks no node
terminal, and takes no side of the split quantifier. RDR 0022 cites
this record as a prospective consumer; 0023/0024 couple only through
JDR 0002 §D1,
cited by C5; 0024's load-refusal appends reach this verb by
construction (C1 consumes the one load pipeline, `0002:C24`'s owners).
Bridge sub-check: n/a — no sibling plan schedules deletion/replacement
of any surface this plan introduces, and this plan retires nothing.

## Alternatives Considered

### Alternative 1: `--emit` on the `lint` verb (O2)

**Description**: The seed's first-named shape — `intrastate lint
--emit=json|dot` exports from the same invocation that verdicts, since
`newAnalysis` already computes the declarations, rows, groups and
reachable node set the document carries (the edge relation is carried
by neither shape — it is recovered in the export path either way —
A2, A3).

**Pros**:

- One invocation lints and exports; CI runs one command.
- No new verb to document.

**Cons**:

- A blocking model exits 2 through `respond.Fail` with a `CLIError`
  envelope that has no graph carrier — `0005:C1` fixes the failure
  envelope as "extended by exactly one omitempty structured field,
  findings, per JDR 0001 §D10" (opened and quoted, not assumed —
  premortem P-11), so carrying a graph on failure means amending a
  locked envelope contract or breaking one-envelope with a second
  document on the stream.
- Verdict-neutrality becomes a property to argue about a shared RunE
  instead of a property the structure gives for free.
- Flag surface lands on the authoritative acceptance verb
  (`0006:C19`), so every future export change re-opens lint's
  contract.

**Reason for rejection**: fails both disqualifying matrix rows
(failing-model export, envelope fit); "lint also exports" couples two
contracts one of which must never influence the other. Its one real
advantage — same-run export structurally IS the certified traversal —
is preserved in the chosen shape by C4's same-entry-surface rule plus
A2's traversal-equality bar, not lost (premortem P-11's drift
concern).

### Alternative 2: file-target export (O3)

**Description**: `--out <path>` (on `lint` or a sibling) writes the
document to a file; stdout keeps its envelope; the seed's third DOT
option.

**Pros**: stdout contract untouched; multi-format runs write several
files.

**Cons**: introduces a side-effecting write to a read-only surface;
CI must manage paths the shell redirect already manages; peers
document the pipe form (`smcat -T dot … -o - | dot …`) as the primary
composition.

**Reason for rejection**: adds an I/O class for nothing a `>` redirect
does not already provide, and still needs a stdout answer for pipes.

### Briefly Rejected

- **JSON-only, DOT via an external converter**: the diagram is half
  the stated user outcome, DOT derives from the same value at the cost
  of one pure renderer, and every opened peer exporter ships dot and
  json as sibling output types of one surface.
- **External generator over `table.Dump` text (O4)**: the dump's
  rendered set literals are non-recoverable by 0002's own decision,
  and the dump carries no reachability relation.
- **Overloading `--as` (text→DOT, json→document)**: conflates envelope
  mode with document format; no cell left for a third format or for a
  raw-JSON pipe.
- **Emitting the graph inside lint's success payload always**: bloats
  every lint call's envelope with data findings-consumers never asked
  for, and still fails the failing-model case.

## Context

### Background

Tracked as kata `intrastate#jjkh` (`release:post-1.0`). Prior art:
the state-machines research audit (CW8) and research notes (§7d)
assigned intrastate the spec/validator role, and review rejected a
hand-maintained formal model as a drifting second artifact. Out of
scope by the seed's own framing: any new lint invariant —
liveness/terminal-reachability is kata `4hps` (now RDR 0022), which
may consume this export's reachability relation but is not blocked by
it, and each ships with the other unlanded; likewise shipping or
maintaining a TLA+/SMV model in-repo — a generator script over the
JSON is a downstream consumer, not part of this work.

### Technical Environment

Go module `github.com/cwensel/intrastate`. Surfaces:
`internal/graphlint/analysis.go`, `internal/graphlint/reach.go`,
`internal/graphlint/taxonomy.go` (finding codes — the only present
output). Governing records: RDR 0002 (normalized rows, dump
identity/ordering), RDR 0003 (the finite product lint enumerates),
RDR 0005 (`--as` envelope contract), RDR 0006 (`0006:A6` —
initial/terminal give reachability a root and stop set).

## Research Findings

### Investigation

Read before enumerating: the in-repo priors
(`docs/cli-output-contract.md`, `internal/cli/respond/respond.go`,
`internal/cli/lint.go`, `internal/table/dump.go`,
`internal/graphlint/{engine,analysis,reach,groups}.go`), the governing
records (`0005:C1`, `0002:§round-trip-inverse-invariants`), and the
external prior art: the seed's role citations re-validated in the
`../state-machines` audit (CW8: "spec/validator, not a runtime") and a
bounded StateMachineRes corpus pass whose one strong instance is
state-machine-cat's `--output-type … dot|…|json|…` exporter (queries,
opened hits, and rejected branches recorded in
`evidence/research/prior-art.md`). Constraint that shaped the choice:
the respond gateway already has a seam (`TextLiner`) for a payload
that IS one canonical string, so a bare-document text mode needs no
gateway exception.

### Key Discoveries

- **Documented** — `0005:C1` carves `lint`/`dump`/`parse`-class
  command groups out of the flow contract, "owned by the RDR that
  names them" ⇒ a root `graph` verb is claimable here without
  amending 0005.
- **Documented** — 0002 fixes the dump's field list and row order but
  "does not define a dump grammar" and names a re-readable form as
  follow-up work ⇒ the JSON schema is this RDR's to own, and reusing
  the dump's field vocabulary keeps one row contract, not two.
- **Documented** — peer exporters (state-machine-cat) treat json and
  dot as sibling output types of one surface streamed to stdout ⇒ the
  format flag + raw stream shape; no opened peer wraps DOT in a JSON
  envelope (⚠ negative corpus result, recorded).
- **Documented** — `reach()` returns nodes and a completeness bit but
  no edge list, and exported `Reach` discards the completeness bit
  (`nodes, _ := reach(m)`) ⇒ edges AND an exported completeness surface
  are this RDR's one extension to the traversal surface (A2, C4).
- **Assumed** — TextLiner passes a multi-line document byte-exact
  (A1); post-fixpoint edge recovery equals the traversal's edges (A2);
  shared-encoder marshaling is byte-stable (A5).

## Trade-offs

### Consequences

- Positive: the certified graph becomes a first-class, diffable CI
  artifact; formal models (TLA+/Alloy/SMV) become derivable downstream
  consumers instead of drifting second artifacts; a refused model can
  still be inspected as a graph.
- Positive: one row contract — the document reuses 0002's dump
  vocabulary, so a column added there (the `emit` precedent) has one
  obvious landing in the schema.
- Negative: `intrastate.graph/1` has a fixed field list; a later
  producer (declared-emit metadata, new invariant surfaces) lands
  additively or bumps the marker under RDR 0029's `0.x` promise.
- Negative: the exported reachability relation exposes the
  OVER-APPROXIMATION lint reasons over (merged nodes, unpruned
  guard/observed edges — `reach.go::Reach` doc); a consumer reading it
  as the runtime relation will over-count edges.

### Risks and Mitigations

- **Risk**: consumers mistake the merged, over-approximate relation
  for runtime behavior and "verify" properties the runtime lacks.
  **Mitigation**: the document self-describes — the `reach` block
  carries an explicit `abstraction` marker (C2's spelling), and the
  docs state the over-approximation contract, and the sound/unsound
  property classes, where the schema is described.
- **Risk**: schema drift — a later field added ad hoc breaks C2's
  additive rule.
  **Mitigation**: golden-fixture byte tests pin `/1`; a failing pin is
  the tripwire that forces the additive-or-version decision.
- **Risk**: the DOT arm quietly diverges from the JSON arm's graph,
  or breaks on hostile content (quotes/newlines/non-ASCII in tag
  values — DOT quoting is its own escaping surface, premortem P-16).
  **Mitigation**: A6's design — DOT renders the export value, not the
  analysis — plus a test asserting DOT node/edge ids equal the JSON
  document's, plus escaping fixtures that must survive `dot -Tsvg`.
- **Risk**: a clean, schema-stamped export reads as certification and
  a PR merges on the diagram while lint failed in another job
  (premortem P-6).
  **Mitigation**: C2 — the document carries no verdict field and the
  docs state an export is never a lint pass; the lint gate, not the
  export, stays the acceptance surface (`0006:C19`).

### Failure Modes

- Visible: refusals reuse lint's arm codes (C1) and the new
  `graph-export-too-large` refusal names the ceiling and remedy
  (accepted cost, premortem P-13: the over-ceiling debugger gets the
  same narrow-a-domain remedy lint gives, not a partial graph that
  diffs as a graph change); a killed process leaves no terminal
  envelope (existing contract).
- Visible: an envelope-sniffing wrapper pointed at text-mode output
  mis-parses the bare document — by contract it must use `--as=json`
  (C5); the four-cell behavior is documented and tested.
- Silent (accepted, marked): a build bump re-baselines DOT diffs and
  may add JSON fields (C3's stated narrowing); the cross-build golden
  makes the moment visible.
- Silent (guarded): nondeterministic bytes across runs — caught by the
  MVV's double-emit byte compare, never shipped silently; a partial
  graph after an incomplete traversal — structurally impossible, C4
  refuses instead.
- Diagnosis: byte-diff two runs (determinism), `jq .data` vs text
  output (mode agreement), `intrastate lint` before/after (neutrality
  oracle).

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified (A1/A2/A5/A6 are spikes; A2 —
  edge recovery equals the traversal's edges — gates Phase 1)
- [ ] A8 verified — the exported completeness-carrying surface exists.
  It gates Phase 2: C4's `graph-export-too-large` MUST, S7, and the
  `disposition` table's ceiling row all depend on it, so Phase 2 cannot
  ship while it is Pending.

### Minimum Viable Validation

1. Author a small state-machine fixture (two owned states, one
   terminal, one escape row) and a decision-table fixture.
2. `intrastate graph --model <fixture>` twice → the two stdouts are
   BYTE-identical, parse as JSON, and carry every C2 field (schema,
   model identity and class, tags, initial, terminal, rows, groups,
   reach.nodes, reach.edges, and the `reach` abstraction marker).
3. `intrastate graph --model <fixture> --emit dot | dot -Tsvg` renders;
   the DOT node/edge id set equals step 2's `reach` block.
4. `intrastate graph --model <fixture> --as=json | jq .data` equals
   step 2's document, value-for-value.
5. `intrastate lint` over a fixture WITH blocking findings still
   refuses identically (byte-compared against a pre-change capture),
   while `intrastate graph` over the same model succeeds with an
   ASSERTED document (defined content, not incidental) — the
   neutrality oracle, mechanism-independent per C4.

End-state: one loadable model, four invocations, byte-level oracles
for determinism, mode agreement, and lint neutrality.

### Phase 1: Export value and edge recovery

Build the document-assembly component over `*table.Model` +
`Reach`-with-edges; goldens pin `intrastate.graph/1`; the A2
differential test (observer vs recovery, subsumption-merge fixtures)
decides the recovery mechanism before the schema freezes.

### Phase 2: CLI verb and gateway wiring

Register root `graph` with C1's arm set, `--emit`, the `TextLiner`
ride, and the `graph-export-too-large` refusal.

### Phase 3: DOT renderer

A pure renderer over the export value; equality test against the JSON
arm's node/edge set.

### Phase 4: Contract surfaces

`docs/cli-output-contract.md` gains the document section; `--help-all`
gains the verb; the schema docs carry C2's `abstraction` marker and its
soundness sentence.

## Validation

### Testing Strategy

The MVV is the acceptance floor; these scenarios are what "done" adds
beyond it. Every oracle is byte- or set-equality — none asserts an exit
code alone.

1. **Scenario**: Determinism under map-seed variation — emit the same
   fixture repeatedly across `--emit`×`--as` cells, with map order
   PROVOKED under `GODEBUG=randmapiter=1` (the mechanism A2's spike
   used; without it a small fixture can pass by luck).
   **Expected**: byte-identical stdout per cell (C3); no cell depends on
   Go map iteration order.
2. **Scenario**: Golden fixtures pin `intrastate.graph/1` for a
   state-machine model and a decision-table model.
   **Expected**: the goldens hold; a field added without a schema
   decision fails the pin (C2's additive rule tripwire).
3. **Scenario**: Lint neutrality — run `intrastate lint` over blocking
   and clean fixtures with the export code present, byte-compared
   against a pre-change capture.
   **Expected**: identical verdict, finding set, and bytes (C4), and the
   oracle holds whichever edge mechanism A2 settles on.
4. **Scenario**: Model that loads but lint refuses.
   **Expected**: `graph` succeeds with an ASSERTED document (defined
   content, not incidental), carrying no verdict or finding field (C2,
   C4).
5. **Scenario**: Mode agreement — `--as=json | jq .data` against the
   `--as=text` document, for both `--emit` values.
   **Expected**: value-for-value equality for `json`; for `dot`, the
   `jq -r .data.dot` unwrap reproduces the DOT text byte-for-byte AFTER
   accounting for the one trailing newline F1 pins on text-mode stdout
   (`document + "\n"`) — the `dot` string member carries the
   document without that gateway newline, so the comparison is against
   the document, not the stream (C5, A1).
6. **Scenario**: DOT arm equality and hostile content — tag values
   carrying quotes, newlines, and non-ASCII.
   **Expected**: DOT node/edge id set equals the JSON document's
   `reach` block, the abstraction marker is present in the header, and
   every fixture survives `dot -Tsvg` (A6, C2, premortem P-16). Exit 0
   is NOT sufficient: each hostile node and edge identifier's emitted
   quoted string is unescaped by inverting A6's recorded escaping ORDER
   (undo the `\n` line-break escape, then `\"`, then `\\`) and asserted
   EQUAL to its source id, so a well-formed but mis-escaped identifier
   fails — the bug class A6's spike hit. Identifiers only; DOT label
   STYLING stays non-normative (F3).
7. **Scenario**: Traversal incomplete under the published node ceiling.
   **Expected**: refusal with `graph-export-too-large` (GroupUserEnv,
   exit 2) naming the ceiling and remedy; no document on stdout (C4).
8. **Scenario**: Selection-arm refusals — both/neither of
   `--model`/`--flow`, `--flow` alone, unreadable file, load failure,
   and an unknown `--emit` value.
   **Expected**: lint's code spellings verbatim plus
   `flag-invalid-value` naming `emit`; exit codes stay the existing 0/2
   mapping (C1); and stdout carries NO document on every refusing arm
   (C4's never-a-partial-document rule, asserted here and not only in
   the mini-check disposition table).
9. **Scenario**: JSON round-trip — decode the exported document.
   **Expected**: value identity on every C2 field including exact set
   members (RT1); no `load ∘ export` inverse is exercised (RT3).

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
> At lock, replace Contradiction Check, Assumption
> Verification, Scope Verification and Proportionality
> with the one-line pointer to gate.md — those four
> judge THIS record at THIS lock and no peer cites
> them. **Cross-Cutting Concerns stays here**, below
> the pointer: it names the project-wide policy other
> RDRs conform to, so it must stay projected and
> citable as `cli/NNNN:G-cross-cutting`. Cite it that
> way, not by section name.

Responses: 0021-lint-normalized-graph-export/artifacts/gate.md (Gate PASS 2026-09-12)

### Cross-Cutting Concerns

[Gate key: cross-cutting]

**Versioning.** The document carries a required leading `schema` field,
initial value `intrastate.graph/1` (C2). Evolution within `/1` is
additive — a consumer ignoring unknown fields keeps working — and an
incompatible change bumps the marker under the `0.x` promise RDR 0029
governs. The exported field spellings are normative as spelled in C2;
the Illustrative Code disclaims literal assertion and binds nothing.

**Canonical form / determinism.** Claimed, and confirmed field by
field:

- *Hash function + library* — none. This record claims no
  content-addressed identity and no replay-stable hash; the promise is
  byte identity of the emitted stream (C3, RT2), not a digest.
- *Pre-image byte layout / primitive encodings* — JSON is emitted
  through the one shared non-HTML-escaping encoder,
  `clierr.WriteJSONLine`, per `0005:C1`'s one-encoder rule (C3).
- *Map iteration order* — no Go map iteration reaches the wire. Every
  sequence is pre-sorted by C2's declared orders before marshaling:
  `rows` in RDR 0002's canonical row order, atoms in 0002's canonical
  atom order, `reach.nodes` by node key, `reach.edges` by
  (from, to, rule), and each tag's `values` array sorted and
  deduplicated — so construction order is unobservable.
- *Whitespace policy / case folding* — field names are lowercase
  snake_case against `internal/table/dump.go::dumpColumns`, a closed
  11-member list with `emit` appended last per `0010:C3`; the
  one-encoder rule fixes the remaining byte layout, and no case
  folding is applied to values.
- *Empty / null / absent distinguishability* — every declared
  collection renders as an empty array `[]` (or object `{}`) when it
  has no members, never `null`; an optional member that does not
  apply is ABSENT, its key omitted, never `null` (C2).
- *Scope* — the promise is (model, build), not across builds
  (premortem P-3). A build bump may change the JSON only by C2's
  additive rule, and may re-baseline DOT diffs freely since DOT
  styling is non-normative.

Oracles: MVV step 2 (two invocations byte-identical), S1 (determinism
under map-seed variation), S2 (golden fixtures pinning
`intrastate.graph/1`), S9 (JSON round-trip), RT2 (`export ∘ export`
byte identity).

**Character encoding.** The shared encoder is non-HTML-escaping, so
tag values and rule ids reach the wire unescaped; S6 exercises hostile
tag content against the DOT arm, where quoting — not escaping-by-
default — is the containment.

**Incremental adoption.** The export is a new root verb beside `lint`,
adding no flag to `lint` and altering no existing output. C4 makes
neutrality a normative obligation with a mechanism-independent oracle:
`intrastate lint`'s verdict, finding set, and bytes are identical with
and without the export code present. Consumers adopt `graph` when they
want it; nothing is required to migrate.

Omitted as not applicable: build tool compatibility, licensing,
deployment model, IDE compatibility, secret/credential lifecycle,
memory management, concurrency model.

## References

- `0005:C1` (envelope contract; command-group carve-out); `0002:§round-trip-inverse-invariants`
  (dump grammar declined, follow-up seeded); `0006:C19` (lint as the
  acceptance surface); JDR 0002 §D1 (projection doctrine, cited by C5).
- Source reviewed: `internal/cli/{lint.go,root.go}`,
  `internal/cli/respond/respond.go`, `internal/table/dump.go`,
  `internal/graphlint/{engine,analysis,reach,groups,taxonomy}.go`.
- Prior art: `../state-machines/AUDIT-rdr-flow.md` (CW8),
  `../state-machines/research/state-machines-research.md` (§7d),
  `../state-machines/repos/state-machine-cat/README.md`
  (`--output-type`); search record in
  `0021-lint-normalized-graph-export/evidence/research/prior-art.md`.
- Related: kata `intrastate#jjkh`; RDR 0022 (prospective consumer of
  the reachability relation).

## Refinement Context (cluster re-entry — delete on re-lock)

Cluster `0021-0029`, reconciled 2026-09-12 (Stage 7.1, iteration 1).
Peer pair: 0021 ↔ 0029. Report:
`docs/rdr/cluster-reconcile/0021-0029/reconcile-report.md`.

**TARGET RE-ENTRY STAGE**: 3 (refine).
**RE-ENTRY SCOPE**: STAGE-SCOPED — the export approach, the document
shape and every byte-determinism proof stand; what is missing is a
stability declaration over surfaces this record already defines, plus
two peer citations. No alternative is reopened.

### Defect 1 — the tiered-surface obligation is never discharged (PW-1, C-1, C-2, C-3)

`0029:C4` is the normative home for "a machine-readable surface takes
its tier in the record that adds it", and it names this record's three
surfaces explicitly:

> the vocabularies `cli/0021` emits under `data` — its `--emit` format
> set, its `graph-export-too-large` refusal code, and the field names
> its C2 fixes at Resolve — are machine-readable surfaces this census
> does not enumerate, because they do not exist on `main` yet. They
> take their tier assignments in `cli/0021` itself, in the change that
> adds them, per the rule above.

This record assigns none. The tier words `frozen`, `append-only` and
`growing` appear nowhere in its text, and `0029:C1`/`0029:C4` are
cited nowhere. By `0029:C4`'s own words — "an unassigned
machine-readable surface is a defect" — the surfaces ship untiered.

**Resolution direction**: in C1 and C2, assign exactly one tier to each
of the three surfaces (the `--emit` value set, `graph-export-too-large`,
and C2's field spellings), and carry one citation noting the envelope
around the document is versioned separately by `0029:C1`. Per `0029:C4`
a tier assignment also obliges an ENUMERATION SEAM — an exported
accessor returning the vocabulary's members — so name the seam for the
`--emit` set, or record `seam: none (prose-only)` with the reason, as
`0029:C4` does for the CLIError `code` row.

### Defect 2 — C2 states a cardinality its peer forbids consumers to assert (C-14)

`0021:C2` fixes `rows[]` as:

> a closed 11-member list with `emit` appended last per `0010:C3`

`0029:C2` forbids a consumer asserting "the set's cardinality, a
member's ordinal position, or a tail position" on an `append-only`
vocabulary. Whether this wording is a legitimate `frozen` declaration
or the exact anti-pattern `0029:C2` retires cannot be decided until
Defect 1 assigns the tier. Settle it in the same pass.

(The descriptive word "closed" itself is NOT in scope: `0029:S9`
exempts sites describing something C4 does not tier, naming
`DumpColumns` — this record's usage — explicitly.)

### Defect 3 — the joint-check never saw 0029 (PW-3, C-9)

Decision Rationale records:

> Joint-check: clear (12 peers) — open peers 0012–0020, 0022–0024

0029 is in neither range, yet `0029` fired a joint decision AT this
record at its lock fence (home `cli/0029 §Normative Contracts` C4).
This record locked 2026-08-28; 0029 is dated 2026-09-11, so the check
could not have seen it and was never re-run.

**Resolution direction**: re-run the propose-time joint-decision check
against 0029 and record the result. Defects 1 and 2 are what a
correctly scoped check would have surfaced.

### Re-verify on re-entry — A4 (C-11)

`0021:A4` ("a new root export verb requires no amendment to RDR 0005's
envelope contract") is Verified against the envelope as it stood.
`0029:C1` adds a non-`omitempty` `schema_version` to every terminal
record at the gateway, including this verb's. A4's conclusion is
expected to survive — it reasons from `0005:C1`'s scope sentence, not
from envelope immutability — but it must be re-verified against the
post-0029 envelope, and any assertion over an exact envelope key set
re-examined.

### Carried forward — scope 0021:S2's goldens to the document (C-8)

`0021:S2` pins golden fixtures for `intrastate.graph/1`. Captured under
`--as=json` they contain the envelope and therefore `schema_version`,
which moves on 0029's minor-increment schedule for reasons unrelated to
this document — training the additive tripwire to be re-baselined.
Scope those goldens to the DOCUMENT (`--as=text`, or `jq .data`).
No golden exists on disk yet, so this is a wording fix, not a
re-capture.

### Standing joint decision — two-marker read order (C-7)

Neither record sequences the two version markers: `0029:C1` governs the
envelope's `schema_version`, `0021:C2` the document's `schema`, and
neither states which an agent reads first, nor the behaviour when the
envelope major is supported and the document marker is not (or the
inverse). Hoisted to `cli/0029 §Normative Contracts` C4 (already the
home for the envelope-versus-document boundary). Open question: **which
marker does a consumer check first, and what does it do when one is
supported and the other is not?** This record does not answer it; it
cites the home once C4 does.
