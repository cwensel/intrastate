Model: claude-sonnet-5

# A8 verification: graphlint surface addition without altering reach() or lint's path

Source examined directly on disk in /Users/cwensel/sandbox/newcoinc/intrastate (no design
documents consulted for the facts below).

---

## Fact 1 — unexported `reach()` signature

File: `internal/graphlint/reach.go:90`

```go
func reach(m *table.Model) (nodes []Node, complete bool) {
```

- Returns `nodes []Node` — only nodes, no edges, no edge type at all (`grep -n 'type Edge'`
  over the package returns nothing).
- Returns `complete bool` — a second, named return value recording whether the traversal
  finished under the published node ceiling (doc comment at reach.go:78-80: "reach returns
  the reachable node set and whether the traversal completed under the published node
  ceiling. An incomplete traversal returns what it found").
- `complete` is set `true` at reach.go:128, flipped to `false` only at reach.go:131-133 when
  `len(nodes) > nodeCeiling` inside the worklist loop, and returned unchanged at reach.go:160
  (`return nodes, complete`). Two early returns (nil model, no declared root) also return
  `true` explicitly (reach.go:104, reach.go:110).

Confirmed: `reach()` returns nodes-only (no edges) plus a completeness bool.

---

## Fact 2 — exported `Reach` signature and body

File: `internal/graphlint/reach.go:73-76`

```go
func Reach(m *table.Model) []Node {
	nodes, _ := reach(m)
	return nodes
}
```

This is exactly the shape the claim describes: `Reach` calls `reach(m)`, binds the second
return value to `_`, and returns only `nodes`. The completeness value computed inside
`reach()` is discarded at this call site and is unreachable from `Reach`'s return type
(`[]Node` carries no completeness field).

Confirmed: `Reach` today discards completeness exactly as the claim states.

---

## Fact 3 — package `graphlint`'s exported surface (does anything else surface completeness?)

Full exported surface of package `graphlint` (non-test `.go` files), gathered via
`rg -n '^func [A-Z]|^type [A-Z]|^const [A-Z]|^var [A-Z]'`:

- `internal/graphlint/reach.go:18` — `type Node struct { Values map[string][]string }`
  (one field, `Values`; no completeness-shaped field)
- `internal/graphlint/reach.go:28` — `const OpaqueValue = "<opaque>"` (unrelated: the
  abstract-value sentinel for infinite-domain tags)
- `internal/graphlint/reach.go:73` — `func Reach(m *table.Model) []Node` (see Fact 2)
- `internal/graphlint/engine.go:13` — `type Request struct { Model *table.Model }`
- `internal/graphlint/engine.go:20` — `func NewRequest(m *table.Model) Request`
- `internal/graphlint/engine.go:23` — `type Report struct { ModelID string; Findings []clierr.Finding }`
  (two fields; no completeness/truncation field)
- `internal/graphlint/engine.go:32` — `func (r Report) Blocking() []clierr.Finding`
- `internal/graphlint/engine.go:35` — `func (r Report) Advisory() []clierr.Finding`
- `internal/graphlint/engine.go:59` — `func Run(req Request) Report`
- `internal/graphlint/engine.go:131` — `func Fingerprint(row table.Row) string`
- `internal/graphlint/taxonomy.go:48` — `const AggregateCode = "graph-lint-failed"`
- `internal/graphlint/taxonomy.go:106` — `func BlockingCodes() []string`
- `internal/graphlint/taxonomy.go:109` — `func AdvisoryCodes() []string`
- `internal/graphlint/taxonomy.go:112` — `func Reasons() []string`
- `internal/graphlint/taxonomy.go:115` — `func Severities() []string`
- `internal/graphlint/taxonomy.go:120` — `func IsBlocking(code string) bool`
- `internal/graphlint/taxonomy.go:147` — `func ProductBound() int`
- `internal/graphlint/taxonomy.go:152` — `func NodeCeiling() int` (returns the model-independent
  ceiling constant `4096`, `taxonomy.go:136` — a fixed config value, not a per-run result)

No struct field, method, sentinel error, or const in this list carries a per-run
completeness/truncation boolean. The only place a per-run "traversal did not complete"
signal reaches an exported type is indirectly, as a **derived lint finding**:

- `internal/graphlint/analysis.go:28-43` — the unexported `analysis` struct holds an
  unexported `complete bool` field (line 34), populated at `analysis.go:46`
  (`nodes, complete := reach(m)`) inside the unexported constructor `newAnalysis` (line 45).
- `internal/graphlint/analysis.go:177-188` — `checkNodeCeiling` reads `a.complete` (an
  unexported field on an unexported type) and, only when it is `false`, emits a
  `clierr.Finding{Code: CodeProductTooLarge, Element: elementTraversal, Message: ...}`
  into the exported `Report.Findings` slice via `Run` (`engine.go:59-80`).

So a caller CAN learn "the traversal was incomplete" only by scanning `Report.Findings` for
`CodeProductTooLarge` — a lossy, string-code-shaped side channel — never by reading a `bool`
off any exported type. No exported symbol returns or stores the raw `complete` value itself.

Confirmed (crux): no caller can observe the raw completeness value today; the only
observable trace is the derived `CodeProductTooLarge` finding, which is a different (lossier)
signal than the boolean the claim is about.

---

## Fact 4 — lint's call path

File: `internal/cli/lint.go:212` (inside `runLint`, defined at `lint.go:134`)

```go
report := graphlint.Run(graphlint.NewRequest(m))
```

`runLint` never references `graphlint.Reach`, `graphlint.reach` (unexported, inaccessible
outside the package anyway), or `graphlint.Node` directly. It goes exclusively through
`graphlint.Run(graphlint.NewRequest(m))` and then reads `report.Blocking()`
(`lint.go:213`) and iterates advisories later in the function.

Inside package `graphlint`, the ONLY call site of `reach()` is `internal/graphlint/analysis.go:46`,
inside `newAnalysis`, which `Run` invokes at `engine.go:65` (`a := newAnalysis(m)`). `Reach`
(exported, reach.go:73) is a second, independent call site of `reach()` used by nothing else
in this repo (`rg -n 'graphlint\.Reach\b' internal/` outside the package returns no hits —
`Reach` currently has no caller in `internal/cli` or elsewhere).

Confirmed: lint reaches the traversal through `Run` → `newAnalysis` → `reach()`, never through
exported `Reach`, and never directly. A new exported function added beside `Reach` in
`reach.go` (e.g. calling its own traversal or wrapping `reach()`) would not touch
`analysis.go:46`'s call site or `lint.go:212`'s call site at all.

---

## Edge computation: does `reach()` already compute edges internally?

File: `internal/graphlint/reach.go:200-223`

```go
func successorsOf(m *table.Model, src Node) []Node {
	index := map[string]int{}
	var out []Node
	for _, row := range m.Rows {
		if row.Kind() == table.KindEscape {
			continue
		}
		if !matchSatisfiable(m, src, row) {
			continue
		}
		produced := successor(m, src, row)
		id := presenceKey(produced)
		if at, seen := index[id]; seen {
			out[at] = joinNodes(out[at], produced)
			continue
		}
		index[id] = len(out)
		out = append(out, produced)
	}
	return out
}
```

`successorsOf` DOES walk row-by-row (each `row` here is conceptually one edge — it is
matched via `matchSatisfiable` and applied via `successor`), so the traversal already knows,
transiently, which row produced which successor node. But the function's return type is
`[]Node` only: the row/edge identity is read, used to compute the successor, and then
discarded — never retained in any struct, slice, or map that survives the call. There is no
`Edge` type anywhere in the package (`grep -n 'type Edge' internal/graphlint/*.go` → no
matches).

`reach()`'s worklist loop (reach.go:130-159) calls `successorsOf(m, nodes[i])` at line 145
and only ever consumes the returned `[]Node`.

Conclusion: edges are computed and immediately discarded during traversal, not retained.
Recovering them for a new exported function does NOT require changing the fixpoint/worklist
algorithm in `reach()` (lines 90-161) — a new sibling function can run its own equivalent
traversal (reusing `successorsOf`'s row-walk logic, or a variant of it that also carries the
producing `row`/`table.Row` alongside each successor `Node`) entirely additively, in a new
function, without modifying `reach()`'s existing signature or body. `reach()` and its single
call site in `analysis.go:46`, and `Reach`'s existing body, need no edits.

---

## Overall verdict on assumption A8

All three limbs of the claim hold on the code as it stands today:

1. `reach()` returns nodes-only plus a completeness bool (Fact 1) — confirmed.
2. `Reach` discards completeness via `nodes, _ := reach(m)` (Fact 2) — confirmed verbatim.
3. No exported symbol in package `graphlint` surfaces the raw completeness value; the only
   observable trace is the lossy `CodeProductTooLarge` finding (Fact 3) — confirmed.
4. Lint (`internal/cli/lint.go::runLint`) reaches the traversal only through
   `graphlint.Run` → `newAnalysis` → `reach()` (analysis.go:46) and `lint.go:212`, never
   through `Reach` and never directly; a new sibling exported function beside `Reach` would
   not touch either call site (Fact 4) — confirmed.
5. Edge recovery does not require altering `reach()`'s traversal algorithm: `successorsOf`
   already walks rows/edges transiently but discards them; a new function can carry them out
   additively without changing `reach()`'s signature, body, or callers.

The claim as stated is supported by the code on disk.

---

## Follow-up: does CodeProductTooLarge make completeness observable in practice?

Coordinator asked for narrowing on the "lossy derived finding" parenthetical above, answered
from the code already read in this file (no fresh search pass).

### Q1 — same condition as `complete == false`, or narrower/different?

Same condition exactly, not a separate pre-traversal estimate.

- `internal/graphlint/analysis.go:177-180`:
  ```go
  func (a *analysis) checkNodeCeiling() {
  	if a.complete {
  		return
  	}
  	a.emit(clierr.Finding{Code: CodeProductTooLarge, ...})
  ```
- `a.complete` is assigned once, verbatim from `reach()`'s own return, at
  `internal/graphlint/analysis.go:46` (`nodes, complete := reach(m)`).
- `reach()` sets `complete = false` in exactly one place, inside the fixpoint worklist loop,
  at `internal/graphlint/reach.go:131-133` (`if len(nodes) > nodeCeiling { complete = false;
  break }`) — not from any independent pre-traversal size estimate or product-of-domains
  computation. `taxonomy.go:147,152` (`ProductBound`, `NodeCeiling`) expose only the static
  configured ceilings, not a computed per-run estimate.

Conclusion: `CodeProductTooLarge` fires on exactly the same condition `reach()` reports as
`complete == false` — not a different or narrower one.

### Q2 — reachable programmatically, or human/rendered-output only?

Reachable programmatically, not merely via lint's rendered CLI output.

- `internal/graphlint/engine.go:59-80` (`Run`) returns the exported `Report` type
  (`engine.go:23-29`: `Report{ModelID string; Findings []clierr.Finding}`), built at
  `engine.go:80` from `a.findings`, which `checkNodeCeiling` appended into.
- Any caller — not just `internal/cli/lint.go` — can call
  `graphlint.Run(graphlint.NewRequest(m))` directly and scan `report.Findings` (or
  `report.Blocking()`, `engine.go:32`) for the exported `CodeProductTooLarge` constant
  (referenced at `analysis.go:182`, declared as an exported taxonomy code). No dependency on
  lint's CLI rendering or on human-readable text — it is a typed, exported `[]clierr.Finding`
  slice on an exported struct.

### Q3 — faithful proxy, recoverable without re-running the traversal?

Yes, one-to-one in both directions, on the code read so far.

- `checkNodeCeiling` emits the finding if and only if `a.complete == false` (guard-and-return
  at `analysis.go:178-180` — no other emission path, no other guard).
- `a.complete` is a direct, untransformed copy of `reach()`'s `complete` return
  (`analysis.go:46`).
- No code path examined has `complete == false` with the finding absent, or `complete ==
  true` with the finding present — presence of `CodeProductTooLarge` in `Report.Findings` and
  `!complete` are the same boolean, encoded as a finding code instead of a `bool`.
- A caller can recover it without re-running anything:
  `complete == !containsCode(report.Findings, CodeProductTooLarge)`.

### Narrowing implication

The raw `bool` `reach()`/`Reach` compute is not exported as a `bool` on any type (Fact 3
above still holds literally). But the CONDITION it represents — "did the traversal complete
under the node ceiling" — is already observable today, faithfully and programmatically, via
`Report.Findings`'s `CodeProductTooLarge` code. A8's wording ("names a condition no caller can
currently observe") is stated more broadly than the code supports: the condition is
observable today, just not as a directly-typed boolean value. A8 should be narrowed to say
the DIRECT/TYPED completeness value is unobservable (true), not that the condition itself is
unobservable (not true, given `CodeProductTooLarge`) — a `NEEDS_DECISION` on the exact wording
before A8 is marked Verified.
