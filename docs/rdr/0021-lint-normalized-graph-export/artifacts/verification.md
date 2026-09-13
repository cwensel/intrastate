## Phase 3a — CoVe

Method: for each REQ-N an input was named that would make a correct
implementation visibly violate it, then run against the built binary
(`go build ./cmd/...`) and, where the obligation names an internal authority,
against that authority directly through a scratch probe. Inputs were the
record and `req-list.md` only; Phase 1 tests, `coverage.md` and
`deviations.md` were not read.

FAIL-20 — `domain` is absent for a tag whose declared domain IS finite.

REQ-20 binds `domain` to be "present exactly when `guard.AssignmentCount`
reports it finite". The producer instead gates on
`finite && len(decl.Domain) > 0` (`internal/cli/graph_document.go:185`), and
a `set` declares its members in `Elements` while a `bool` declares none at
all — so both report finite yet ship with the key omitted.

Failing input (a): `internal/table/testdata/kata-fixture.toml`, whose
`[tags.labels]` is `kind = "set"` with `elements = ["bug", "chore",
"needs work"]`.
Observed: `intrastate graph --model internal/table/testdata/kata-fixture.toml`
emits `{"name":"labels","provenance":"owned","kind":"set","required":false,
"single_valued":false}` — no `domain` member. A direct probe over the named
authority returns `AssignmentCount({Kind:"set", Elements:[bug chore
"needs work"]}) = (16, true)` — FINITE — so REQ-20 requires the key present.

Failing input (b): the same fixture with an added observed `bool` tag
(`[tags.flagged]` / `provenance = "observed"` / `kind = "bool"`).
Observed: the tag emits as
`{"name":"flagged","provenance":"observed","kind":"bool","required":false,
"single_valued":false}` — no `domain`, while
`AssignmentCount({Kind:"bool"}) = (2, true)`, FINITE.

The non-finite arm is correct: `scalar` (`owner`) and a domain-less `enum`
(`recognized`) both report `finite=false` and both correctly omit the key,
and a non-finite `set` correctly abstracts to `<opaque>` in the reach values
(REQ-27 holds). The defect is confined to the finite-but-not-`Domain`-backed
kinds, where the code's own doc comment
(`graph_document.go:51-54`, "present exactly when `guard.AssignmentCount`
reports the declaration's domain finite … the finite arm always carries at
least one member") states the binding the gate then fails to implement.

---

Inputs run and NOT violated (each exercised against the built binary):

REQ-8/A-1 default `--emit` byte-equals explicit `--emit=json`. REQ-10
`EmitFormats()` accessor present. REQ-11 `--model <unreadable> --emit=xml`
refuses `flag-invalid-value` naming `emit`, never `model-unreadable`.
REQ-16/A-2 `schema` is the first key on the wire. REQ-21/REQ-32 `terminal`
carries the dereferenced declared predicate sets; no node is marked terminal.
REQ-26/REQ-107 nodes, edges, rows, tags and groups all emit pre-sorted.
REQ-27 a non-finite `set` abstracts to `<opaque>`. REQ-28/REQ-94 no `null`
reaches the wire on any fixture tried. REQ-31 no verdict/finding key.
REQ-37 `schema_version` is envelope-only, absent from `data`. REQ-42/REQ-73
all four `--emit`x`--as` cells are byte-stable, including under
`GODEBUG=randmapiter=1`; REQ-66 20/20 runs byte-identical. REQ-43/REQ-109
zero HTML escapes with `<`, `>`, `&` in tag values, in all three arms, the
raw bytes reaching the wire. REQ-46 the same model exported from two
unrelated paths is byte-identical and leaks no path substring. REQ-51/REQ-77/
REQ-87 a model lint refuses (`graph-lint-failed`) still exports exit 0 with a
defined document. REQ-52/REQ-81/REQ-88 a 4200-state chain refuses
`graph-export-too-large`, exit 2, naming the 4096 ceiling and the
narrow-a-domain remedy, with zero stdout bytes. REQ-53/REQ-56/REQ-100
`ReachGraph` is a new sibling in `export.go`; `reach.go` is untouched by the
implementation commit. REQ-57/F1 text-mode stdout is the document plus
exactly one newline, stderr empty. REQ-58/A-5 `--as=json --emit dot` carries
`data` as an object whose single member is `dot`. REQ-60 all four cells exit
0. REQ-62 no projection flag on the verb. REQ-65/REQ-84 every C2 field
decodes value-for-value. REQ-71/REQ-78 `jq .data` equals the text document,
and the DOT unwrap matches after the one-newline accounting. REQ-80 every
hostile node id (`q"uote`) and rule label is recovered exactly by inverting
A6's escaping order. REQ-83 stdout carries zero document bytes on all five
refusing arms. REQ-108 field names are lowercase snake_case matching
`dumpColumns`. REQ-110 `lint` gained no `--emit` flag.

## Verdict — one violation (FAIL-20)
## Phase 3b — Adversarial

Inputs: the record's Failure Modes section, `req-list.md`, and the
implementation source. Phase 1 tests and Phase 3a findings were not read.

The Failure Modes section enumerates the export's silent modes as guarded
or accepted: "nondeterministic bytes across runs — caught by the MVV's
double-emit byte compare, never shipped silently; a partial graph after an
incomplete traversal — structurally impossible, C4 refuses instead." That
enumeration reasons only about edges the traversal took but failed to
publish. It never considers the opposite direction — edges the export
publishes that the traversal never took — which is the direction A2 fixes
as "EXACTLY the edges the traversal itself took ... equality, not
plausibility, is the bar."

ADV-1 — The exported edge relation carries edges the traversal never took.
`internal/graphlint/reach.go::successorsOf` skips every escape row
(`if row.Kind() == table.KindEscape { continue }`), so the fixpoint takes
no edge from one, and `Reach`'s own contract fixes the rule: "an edge is a
normalized NON-ESCAPE row whose match pattern over owned tags is
satisfiable in the source node." `internal/graphlint/export.go::edgesFrom`
instead admits each satisfiable escape row as a SELF-LOOP, so the exported
relation is strictly larger than the one lint certifies — the verb/lint
drift `0021:C4` states is structural rather than disciplinary. Confirmed at
both layers: `ReachGraph` returns the phantom edge, and the shipped verb
emits it into the document's `reach.edges`. The failure is SILENT: the
phantom edge carries a real rule id and two real node ids, is correctly
sorted, and names nodes the document also carries, so every shipped oracle
accepts it — `TestReq25_ReachEdgesCarryFromToAndRule` checks only an edge's
shape and endpoint membership, `TestReq26And107_...` only the sort order,
and the DOT arm compares itself against the same `reach` block. A consumer
proving a universal claim ("no path does X") over this relation reasons
about a transition the certified graph does not contain, which is the
soundness direction C2's abstraction marker promises is safe.
Test: `internal/cli/graph_adversarial_0021_test.go::TestAdv1_TheExportedEdgeSetCarriesNoEdgeTheTraversalNeverTook`
— FAILS on both fixtures against the current implementation.

ADV-2 — The phantom edge count scales with the reachable node set. A
catch-all escape row constrained only on the recognized outcome is
satisfiable over owned tags in every node, so the export mints one phantom
self-loop PER NODE (observed: 3 self-loops over 3 nodes). The blast radius
therefore grows with the model rather than staying a fixed cost, and
`0021:C2`'s DOT arm renders each phantom into the diagram a human reviewer
reads, so the misreading reaches both machine and human consumers. A
reviewer diffing the exported graph across a commit that merely adds a
state sees the phantom set grow with it.
Test: `internal/cli/graph_adversarial_0021_test.go::TestAdv2_NoPhantomSelfLoopIsMintedPerReachableNode`
— FAILS against the current implementation.

Hypotheses probed and REFUTED (no test added; a test that passes catches
nothing):

- Dropped edges via `edgesFrom`/`successorsOf` join asymmetry. `edgesFrom`
  computes each row's successor unjoined while `successorsOf` joins
  successors sharing a presence footprint, so `indexOf`'s `subsumes` test
  looked able to fail and silently drop an edge with `complete == true`.
  Probed over value-widening joins and over key-shrinking `clear` rows:
  every node retained an inbound edge in both. `subsumes` is
  containment-based and `presenceKey` keeps differing footprints as
  distinct nodes, so the `indexOf < 0` branch is defensive, not live.
- DOT initial-node marking under an opaque domain. `sameMembers` returns
  true unconditionally on an opaque held value and `initialNodeID` takes
  the first match in key order, so a non-root node looked markable. An
  owned `scalar` model collapses to a single node, correctly marked; not
  reproduced.
- Determinism under provoked map order and no-document-on-refusal. Both
  are already covered and green: `TestReq43And73_...` runs under
  `GODEBUG=randmapiter=1`, and every refusing arm asserts
  `assertNoDocumentOnStdout`.

Baseline: the only other reds in `internal/cli` are the three known
TEST-FIXTURE items (D2/REQ-109, D3/REQ-67, D4/REQ-85);
`internal/graphlint` is green. `gofmt` and `go vet` are clean.
