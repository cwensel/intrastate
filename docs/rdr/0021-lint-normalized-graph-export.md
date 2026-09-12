# Recommendation 0021: Lint's normalized-graph export

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-08-28
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
- **Type**: Feature
- **Profile**: large — one contract, the deterministic export
  grammar of the normalized graph, stated as five clauses (C1–C5)
  of one seam; user-facing yes; locks format (the
  `intrastate.graph/1` wire format and its DOT projection).
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
    (`evidence/research/a3-a4.md`). The edge relation is NOT carried by
    the model value: `reach()` computes nodes only and `successorsOf`
    reads `row.RuleID` but discards it, so no edge carrier exists
    today. It is constructible with no re-parse (the rule id is on
    `table.Row.RuleID`) and is recovered in the export path per A2.
    This narrows the pre-edit wording ("carries everything the
    document needs"), which overstated the model value's coverage;
    narrowed at Resolve 2026-09-12. C2 is unchanged — it predicates
    *carries* of the document, which still carries edges.
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
    into that list (re-anchored at Resolve 2026-09-12,
    `evidence/research/a3-a4.md`).
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
  alone — nodes, edges, initial, and terminal-satisfaction are enough —
  with no reach or analysis internals consulted.**
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
and codes verbatim (`internal/cli/lint.go::runLint`): exactly one of
`--model <path>` / `--flow <id>`; both → `flag-mutually-exclusive`;
`--flow` alone → `flag-invalid-value` (this build resolves no ids);
neither → `flag-required`; unreadable file → `model-unreadable`; load
failure → `model-invalid` with one findings[] entry per load category.
`--emit <format>` selects the document: `json` (default) or `dot`; any
other value → `flag-invalid-value` naming `emit`. RunE starts with
`respond.ValidateMode`, success routes through `respond.OK`, failure
through `respond.Fail`; exit codes are the existing 0/2 mapping — no
new exit group.
```

**C2**

```normative
DOCUMENT. The JSON document is this RDR's wire format, versioned by a
required leading `schema` field, initial value `intrastate.graph/1`;
evolution within `/1` is strictly additive (a consumer ignoring
unknown fields keeps working). It carries, at minimum: model identity
and class; the tag declarations (name, provenance, kind, required,
single-valued, and the declared domain exactly when
`guard.AssignmentCount` reports it finite); the declared `[initial]`
assignments; the declared `terminal` predicate sets; the normalized
rows carrying RDR 0002's dump field list as structured values —
identity, source, kind, outcome, atoms (each `{key, operator,
literal[], block}`), next, writes, requires_owned, gate, escape, emit
— in 0002's canonical row order with atoms in 0002's canonical atom
order; the selection-context groups (context plus member row
identities); and the reachability relation — merged fixpoint nodes
(`{id, values}`, id = the canonical node key) and edges (`{from, to,
rule}`) — with nodes sorted by node key and edges by (from, to, rule),
so construction order is unobservable. Every declared collection
renders as an empty JSON array `[]` (or object `{}`) when it has no
members — never `null`; an optional member that does not apply is
ABSENT, its key omitted, never `null`. The `reach` block carries a
REQUIRED abstraction marker, spelled `abstraction` with the token
value `declared-over-approximation`, stating the relation
is the DECLARED over-approximation, not the runtime — merged nodes,
guard/observed atoms unpruned (`reach.go::Reach` doc) — so a formal
consumer can tell which property classes are sound over it. The
document carries NO verdict or finding field: an export is never a
lint pass, and the schema docs say so. Set-valued members are JSON
arrays, closing `0002:§round-trip-inverse-invariants`'s lossy
set-literal rendering for this document; the document is NOT a model
source and no export→load inverse is claimed. Exact field spellings
are normative as the Illustrative Code spells them: `schema`, `model`,
`class`, `tags[{name, provenance, kind, required, single_valued,
domain}]`, `initial`, `terminal`, `rows[…0002's field list…]`,
`groups[{context, rules}]`, and `reach{abstraction, nodes[{id,
values}], edges[{from, to, rule}]}`. The schema docs state the
soundness rule in one sentence: universal claims ("no path does X")
proved over this relation hold at runtime; existence claims ("some
path reaches X") may be spurious. The DOT
document renders the same value: one node per reachability node, one
edge per reachability edge labeled with its rule id, the initial node
and terminal-satisfying nodes marked, and the abstraction marker
rendered in the graph header comment/label so the diagram carries it
too (premortem P-7); its node/edge SET and the marker are normative,
its styling/attributes are not.
```

**C3**

```normative
DETERMINISM. For one model input and one build, emission is
byte-for-byte identical across invocations, in every `--emit` and
`--as` combination. Every sequence on the wire is pre-sorted by C2's
orders before marshaling; no Go map iteration reaches the wire; JSON
is emitted through the one shared non-HTML-escaping encoder
(`clierr.WriteJSONLine`, per `0005:C1`'s one-encoder rule). This is a
DELIBERATE NARROWING of the seed's "deterministic for the same model"
to (model, build) — stated, not silent (premortem P-3): across builds
the JSON document changes only by C2's additive schema rule (a
cross-build golden pins it), and DOT styling carries no cross-build
stability promise, so a build bump may re-baseline DOT diffs and may
only ADD to JSON ones.
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
present (the MVV asserts this), and the oracle is
MECHANISM-INDEPENDENT — it binds equally if Resolve picks A2's
in-traversal edge observer (premortem P-14). The verb succeeds for
ANY model that loads, including a model lint refuses; the
model-loads-but-lint-refuses case is a dedicated fixture with an
asserted, defined document — never whatever the traversal happens to
do (premortem P-6).
When the traversal is incomplete under the published node ceiling
(`reach.go::reach` returns `complete == false`), the verb refuses with
the scalar code `graph-export-too-large` (GroupUserEnv, exit 2) naming
the ceiling and the narrow-a-domain remedy — never a partial document,
because a partial graph diffs as a graph change.
```

**C5**

```normative
MODE COEXISTENCE. Under `--as=text` (the default) stdout carries the
selected document verbatim and nothing else — the payload satisfies
the gateway's `TextLiner` seam, so the DOT stream pipes to `dot` and
the JSON document diffs raw in CI; advisories stay on stderr. Under
`--as=json` stdout carries exactly one terminal `ok` envelope whose
`data` embeds the same document: the document object for `--emit
json`, a single string field carrying the DOT text for `--emit dot`
(the documented unwrap is one `jq -r` step). Both modes derive from
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

#### Load-Bearing Decisions

- **Identity** — a node is its canonical node key
  (`reach.go::(Node).key` — injective by escaping, ⇒ two exported
  nodes never collide); a row is RDR 0002's row identity; an edge is
  the `(from, to, rule)` triple.
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

- `json-decode ∘ export = value identity on the exported projection`:
  decoding the JSON document reconstructs, value-for-value, every
  field C2 lists — including exact set members, the recoverability
  0002's text dump deliberately declined
  (`0002:§round-trip-inverse-invariants`).
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
| Selection arms | `internal/cli/lint.go::runLint` | Codes are lint's own by contract (`0006:C19`, `0006:C20`) | Reuse the arm set and code spellings | C1 mirrors them verbatim |

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
grepped for this RDR's modify-anchors (`reach.go::Reach` extension,
new `graph` verb/root registration) and contract literals (`--emit`,
`intrastate.graph/1`, `graph-export-too-large`, the mirrored C1
codes): no whole-token hit shares an undecided contract. Context
beside the clear: 0015 names "the 0021 export" only inside its
REJECTED alternative's cons, under the merged-node doctrine settled
at JDR 0001 §JD-23 (0022/0015's home) — this export emits the merged
fixpoint relation as-is and takes no side of that doctrine; 0023/0024
couple with this RDR only through JDR 0002 §D1, cited by C5; 0024's
load-refusal appends reach this verb by construction (C1 consumes the
one load pipeline, `0002:C24`'s owners). Bridge sub-check: n/a — no
sibling plan schedules deletion/replacement of any surface this plan
introduces, and this plan retires nothing.

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
  no edge list ⇒ edges are this RDR's one extension to the traversal
  surface (A2).
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
- Negative: `intrastate.graph/1` is a locked wire format; every later
  producer (declared-emit metadata, new invariant surfaces) must land
  additively or version the schema.
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
   fixture repeatedly across `--emit`×`--as` cells.
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
   **Expected**: value-for-value equality for `json`; the `jq -r` unwrap
   reproduces the DOT text byte-for-byte for `dot` (C5).
6. **Scenario**: DOT arm equality and hostile content — tag values
   carrying quotes, newlines, and non-ASCII.
   **Expected**: DOT node/edge id set equals the JSON document's
   `reach` block, the abstraction marker is present in the header, and
   every fixture survives `dot -Tsvg` (A6, C2, premortem P-16).
7. **Scenario**: Traversal incomplete under the published node ceiling.
   **Expected**: refusal with `graph-export-too-large` (GroupUserEnv,
   exit 2) naming the ceiling and remedy; no document on stdout (C4).
8. **Scenario**: Selection-arm refusals — both/neither of
   `--model`/`--flow`, `--flow` alone, unreadable file, load failure,
   and an unknown `--emit` value.
   **Expected**: lint's code spellings verbatim plus
   `flag-invalid-value` naming `emit`; exit codes stay the existing 0/2
   mapping (C1).
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

### Contradiction Check

[Gate key: contradiction — a gate response is cited as
`cli/NNNN:G-<key>`, so the key is a stable id and is
not derived from this heading, which may be reworded.]

[State any conflicts between Research Findings and
the Proposed Solution. If none exist, state
"No contradictions found between research findings,
design principles, and proposed solution."]

### Assumption Verification

[Gate key: assumptions]

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

[Gate key: scope]

[Confirm the Minimum Viable Validation is in scope
and will be executed during implementation, not
deferred. State the specific test or proof.]

### Cross-Cutting Concerns

[Gate key: cross-cutting]

[Retained at lock — this sub-section stays in the RDR
when the other gate responses move to gate.md, because
peer RDRs cite it as `cli/NNNN:G-cross-cutting` and an
element that is not projected cannot be cited.]

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

[Gate key: proportionality]

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
