# Triage — RDR 0021 Lint Normalized Graph Export

Single-pass unattended triage of the roborev findings left by the Stage-8
launch. Window frozen once: `BASE=86f8c24`, `HEAD=b409d65` (6 commits).
Batch label: `batch:rdr-0021`.

- **Spine** — 2 in-window per-commit auto-reviews carrying findings (jobs
  7792, 7793), 5 distinct findings. Jobs 7794, 7795 and 7796 reviewed clean
  ("No issues found"; 7796's diff is docs-only and empty). Job 7791 reviewed
  `b409d65` before the rebase and its ref is not an ancestor of `HEAD` — out
  of window, untouched. Ancestry confirmed mechanically per job ref, not by
  subject.
- **Range net** — one `--since BASE` review (job 7797), 4 findings. Two are
  duplicates of the spine, charged to job 7793. One re-raised job 7792's
  round-trip oracle at its post-fixup line (`:915`, was `:875`), proving it
  was **not** stale and charging it here. One is genuinely new and appears in
  no per-commit review — the `export.go` edge/successor join asymmetry, an
  interaction defect that only the cross-commit net could surface.
- **Pre-filter** — 2 findings dropped without a leaf (Phase 1.3): the cited
  hunks were rewritten by `b409d65`.

Every surviving finding was grounded in a read-only leaf against the RDR +
`{RDR_RESOURCES}`, which roborev's repo-sandboxed prompt never sees. Three of
the four IN-SCOPE verdicts were confirmed by falsification probe in a detached
scratch worktree, never in the launch worktree.

## Dispositions

| # | Source | Location | Verdict | odc-type | odc-trigger | Outcome |
|---|---|---|---|---|---|---|
| 1 | 7792 | `graph_document_0021_test.go:1101` — escape loops mutually unsatisfiable | DROP superseded-at-HEAD | n/a-not-a-defect | design-conformance | `drop:superseded-at-HEAD` — `b409d65` (D2) repointed loop 1's needles at `<`/`>`/`&` |
| 2 | 7792 | `graph_stability_0021_test.go:176` — inverse-load scanner matches its own source | DROP superseded-at-HEAD | n/a-not-a-defect | design-conformance | `drop:superseded-at-HEAD` — `b409d65` (D3) excludes the scanning file from its own sweep |
| 3 | 7792 + 7797 | `graph_document_0021_test.go:915` (was `:875`) — RT1 round-trip oracle is a self-consistency tautology | IN-SCOPE (`0021:RT1`) | test-oracle | design-conformance | `fixed:3a1b722` — the oracle now asserts the decoded document against the fixture's AUTHORED values (initial `flavors == ["x","y"]`, the `advance` row's write `flavors == ["y","z"]`); `Initial` and `Writes` added to the decode view |
| 4 | 7793 + 7797 | `graph_dot.go:81` — `nodeMatchesInitial` skips a cleared initial key, marking a successor as the root | IN-SCOPE (`0021:C2` DOT / REQ-39) | algorithm | rare-situation | `fixed:3a1b722` — a missing declared initial key now returns false and the arity test is `len(initial) == len(n.Values)`; regression asserts WHICH node carries the marker |
| 5 | 7793 + 7797 | `graph_document.go:297` — empty owned `set` exports `"values":{"key":null}` | IN-SCOPE (`0021:C2` / REQ-28, REQ-94) | assignment | boundary | `fixed:3a1b722` — routed through the existing `stringsOrEmpty` helper; regression covers an empty finite `set` in BOTH JSON output modes |
| 6 | 7797 | `internal/graphlint/export.go:60` — edge recovery resolves successors before `successorsOf`'s presence-footprint join, orphaning merged nodes | IN-SCOPE (`0021:A2`, `0021:C2` / REQ-25) | algorithm | rare-situation | `fixed:3a1b722` — new `joinedEdgesFrom` joins each source's produced edges by presence footprint (mirroring `successorsOf`) retaining rule ids, then resolves the JOINED target; regression compares endpoints by EXACT KEY, never through `indexOf` |

No finding routed OUT-OF-SCOPE, so no kata was filed and no `## Open question`
survives. Nothing was undecided: every verdict rests on positive evidence.

## The four IN-SCOPE findings

**#3 RT1 oracle (`test-oracle`).** The round-trip test decodes the exported
document into `first`, re-encodes it, decodes that into `second`, and compares
`first` to `second` — a comparison of `f(x)` to `f(x)` that no input can
falsify. RT1 (RDR `:631`/`:635`) states value identity on every field **C2
lists**, and says the quantifier is C2's field list, *not* "whatever was
exported", precisely so the invariant constrains the emitter rather than
restating itself. Falsified empirically: mutating the exporter to drop the last
member of every multi-member sequence made it emit initial `["x"]` for authored
`["x","y"]` and write `["y"]` for authored `["y","z"]`, and the test still
passed with the whole `internal/cli` package green. Not redundant — no other
test compares set members to authored values, and neither golden model contains
`flavors`. Fix: assert decoded values against the fixture's authored values,
which requires adding `Initial` to `graphDocument` and `Writes` to `graphRow`
in `graph_probe_0021_test.go` (both fields absent today).

**#4 initial-node marker (`algorithm`).** `nodeMatchesInitial` `continue`s on a
declared initial key the node does not hold, so a cleared key vanishes and
`declared == len(n.Values)` holds vacuously; node ids sort lexically and the
shorter key is a strict prefix, so the successor sorts first and takes the
marker. Reproduced at HEAD: a model with owned `flag`+`status` and a rule
clearing `status` marks the successor `"flag=false,;"` while the true root
`"flag=false,;status=a,;"` goes unmarked — and the marked node has an inbound
edge from the root. The finding's literal "all initial tags cleared" case is
unreachable (`normalize.go:498` forces a write block), but the weaker
cleared-subset path reproduces, so this is not a DROP. The `!ok` guard is dead
code: `load.go:1714` makes a non-owned `[initial]` key unauthorable. The
existing test only greps for the substring `initial`, so it passes against the
buggy render.

**#5 empty set renders `null` (`assignment`).** `graphReachDocOf` builds
`values[key] = slices.Compact(slices.Sorted(slices.Values(held)))`, which is
nil for an empty held sequence. C2 names `values` as an object "whose every
value is that tag's sorted, deduplicated value ARRAY" and states the blanket
rule in the same clause — every declared collection renders `[]` when empty,
"never `null`". REQ-94 binds the producer specifically, and `graphReachDocOf`
is that producer. Reproduced at HEAD in both output modes:
`"values":{"labels":null,"status":["a"]}` while `initial` and `writes` render
`[]` in the same document. The existing null-check test is vacuous because its
fixture declares no `set`. Fix is one call to the helper that already exists
for this purpose: `stringsOrEmpty(...)`.

**#6 edge/successor join asymmetry (`algorithm`).** `ReachGraph` recovers edges
over the settled node set, but resolves each row successor through `indexOf`
individually, while the fixpoint's `successorsOf` joins successors by presence
footprint. A cycle that re-enters an already-widened node therefore publishes
edges to singletons only, leaving the merged nodes with no incoming edges.
Confirmed on the shipped CLI wire, not just the package API: on an
`a→b, b→c, c→a, c→b` model the document carries 6 nodes and 13 edges, and three
of the six nodes — `s=a,b,;`, `s=a,b,c,;`, `s=b,c,;` — have zero incoming
edges. Half the published node set is unreachable in the published relation,
and the DOT arm renders the same value, so a reviewer sees floating states.
A2 requires "EXACTLY the edges the traversal itself took — equality, not
plausibility"; commit `8e8dc0f` (D5) adjudicated the opposite direction —
escape rows admitted as edges — and never reached the join asymmetry. The
soundness direction is safe (edges carry real rule ids), which is why every
shipped oracle accepts it.

Note for the fixup: `verification.md` Phase 3b lists this hypothesis under
"Hypotheses probed and REFUTED", but that probe tested the wrong failure mode.
It is correct that the `indexOf < 0` branch is dead — `indexOf` misses are zero
— yet the defect is a mis-targeted edge via a subsumption **hit** on the
narrower node. Phase 3b's own oracle, "every node retained an inbound edge", is
what this cycle breaks, so the entry is not an adjudication of this finding.
A2's spike is likewise not a witness: its differential mapped both sides
through `indexOf`, which is self-confirming. Any regression added here must
compare targets by exact key, never through `indexOf`.

## Fixup outcome (`3a1b722`)

All four IN-SCOPE rows are `fixed:3a1b722`; the two DROP rows were terminal
already and are untouched. Each fix was proven non-vacuous by reverting the
three production files with the new tests kept: every regression FAILED
against the pre-fix code with the diagnosed symptom — the DOT marked
`flag=false,;` (the successor) rather than the root, the empty owned `set`
rendered `"labels":null`, and the widened cycle left 3 of 6 nodes with zero
incoming edges. The production files were then restored and the full suite
run green.

Two record-keeping consequences:

- `verification.md`'s Phase 3b entry listing the export join asymmetry under
  "Hypotheses probed and REFUTED" is superseded for this defect. That probe's
  lemma — the `indexOf < 0` branch is dead — still holds (the fix does not
  touch it), but its oracle "every node retained an inbound edge" is exactly
  what the cited cycle breaks, so it never adjudicated this failure mode.
  The new regression therefore compares endpoints by exact node key rather
  than through `indexOf`, which is self-confirming.
- REQ-18 was demoted to `req-list.md`'s `## EXCLUDED` section and its empty
  `coverage.md` row removed, on deviation D1's existing rationale: the clause
  has no observable of its own and a direct test would assert the ABSENCE of
  an assertion. No vacuous test was written.

## Why two spine findings were dropped

Both cited hunks were rewritten by `b409d65` while resolving the launch's own
author decisions, so neither reproduces at HEAD. D2 replaced loop 1's raw `<`,
`>`, `&` needles with the escape sequences REQ-109 actually forbids, resolving
the contradiction against loop 2's `a<b`/`x&y` expectations. D3 added the
scanning file to its own exclusion list, so the inverse-load scanner no longer
matches its own predicates. These are the flow's own red-then-green rhythm, not
surviving defects.
