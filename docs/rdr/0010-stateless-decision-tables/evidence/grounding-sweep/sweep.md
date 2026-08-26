Model: claude-fable-5

# cli/0010 — Stage 2 grounding micro-sweep (factored; checker saw anchors only)

Verdict: PASS — 33 confirmed / 0 refuted / 0 not-found.

| Anchor | Verdict | Where |
| --- | --- | --- |
| G1 `internal/table/normalize.go::normalizeRule` no-write-block arm; kind by presence of `escape` | CONFIRMED | normalize.go:336,366 |
| G2 six `Cat*` categories | CONFIRMED | category.go:18–37 |
| G3 `checkDanglingEdge` root arm + terminal arm | CONFIRMED | analysis.go:112–150 |
| G4 `nodeSatisfiesMatch` owned atoms only | CONFIRMED | analysis.go:78–88 |
| G5 `reach` returns no nodes without `[initial]` | CONFIRMED | reach.go:90–95 |
| G6 `checkGroups` coverage/overlap behind reachable gate | CONFIRMED | groups.go:31–44 |
| G7 `Fingerprint` hashes `Atoms` + `NextTags` only | CONFIRMED | engine.go:131–154 |
| G8 `invokedReaders` demand set | CONFIRMED | flow_exec.go:86–114 |
| G9 `runReaders` artifact-missing scoped to invoked names | CONFIRMED | flow_exec.go:175–195 |
| G10 `resolvePayload` fields, no `emit`; `--outcome` required | CONFIRMED | flow_resolve.go:31–43,94 |
| G11 `missingOwned`, `Row`, `Plan` | CONFIRMED | resolve.go:579,256,338 |
| G12 `dumpColumns` closed list | CONFIRMED | dump.go:13–16 |
| G13 `loadDump` refuses omitted column | CONFIRMED | load.go:447–470 |
| G14 `Model.Metadata` uninterpreted; `Initial`/`Terminal`/accessor maps | CONFIRMED | model.go:370–385 |
| G15 `sourceModel`/`sourceRule`; strict via `DisallowUnknownFields` (source.go, not load.go) | CONFIRMED | source.go:28,74,115,119 |
| G16 reader arity; `[initial]` key categories; write to non-owned | CONFIRMED (note) | load.go:351–373,546,574 — a declared non-owned `[initial]` key has no dedicated arm: `malformed accessor binding` or `write to non-owned tag` |
| G17 `[dump]` fixture count | CONFIRMED (correction) | 103 files under `internal/table/testdata`; none in `models/`; 0002 spike fixtures also carry it |
| G18 `guard.Groups` | CONFIRMED | product.go:66 |
| G19–G24 0002:C4, C2, C14, C15, C19, C3, A1 quotes | CONFIRMED | projector |
| G25 0005:C1 quotes | CONFIRMED | projector |
| G26–G29 0006:C18, D-reachability-relation, C7, C15 | CONFIRMED | projector |
| G30–G33 stateless / scxmlcc / transitions / MODEL-transition quotes | CONFIRMED | sibling checkout |

Folded into the record: A4 fixture count; C2 refusal categories.
