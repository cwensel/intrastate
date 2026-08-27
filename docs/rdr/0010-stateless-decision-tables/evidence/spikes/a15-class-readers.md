Model: claude-opus-5[1m]

# A15 — four class readers can reach `Model.Class`; C1 window is real

Read-only source verification. RDR read through
`rdr inspect` only.

## The four readers (from `0010:§authority`)

`normalizeRule`, `reach`, `checkDanglingEdge` root arm, `checkCoverage`
zero-dim arm.

| Reader | file:line | Reaches `Class`? | Via what value |
| --- | --- | --- | --- |
| `normalizeRule` | `internal/table/normalize.go:335` | yes | method on `*loader`; `l.model *Model` (`internal/table/load.go:68`) is the model under construction — `l.model.Class` |
| `reach` | `internal/graphlint/reach.go:90` | yes | parameter `m *table.Model` — `m.Class` directly |
| `checkDanglingEdge` root arm | `internal/graphlint/analysis.go:112` (root arm at `:113`, `len(a.model.Initial) == 0`) | yes | method on `*analysis`; `a.model *table.Model` (`internal/graphlint/analysis.go:29`) |
| `checkCoverage` zero-dim arm | `internal/graphlint/coverage.go:33` (zero-dim arm at `:35`, `len(dims) == 0`) | yes | same `a.model *table.Model` (already read at `:34` via `guard.Dimensions(a.model, g)`) |

## (a) Call-path trace

Table reader:

- `internal/table/load.go:75` `func (l *loader) run() (*Model, error)` sets
  `l.model = &Model{}` at `:76` and runs the fixed step slice.
- `normalizeRules` (`internal/table/normalize.go:282`) is step 9; it calls
  `l.normalizeRule(rule, id, setKeys)` at `:308`. Both are `*loader`
  methods, so `l.model` — and therefore an exported `l.model.Class` — is in
  scope at the reader's body.

Graphlint readers:

- `internal/cli/lint.go:99` `m, advisories, err := table.LoadWithAdvisories(src, modelPath)`
- `internal/cli/lint.go:115` `report := graphlint.Run(graphlint.NewRequest(m))`
- `internal/graphlint/engine.go:19` `NewRequest(m *table.Model) Request { return Request{Model: m} }`
- `internal/graphlint/engine.go:59` `Run(req Request)`: `m := req.Model`,
  then `a := newAnalysis(m)` at `:65`.
- `internal/graphlint/analysis.go:63` `newAnalysis(m *table.Model)` calls
  `reach(m)` and stores `model: m` on the `analysis`.

So `reach` gets the `*table.Model` as its own parameter, and every
`*analysis` method reaches it as `a.model`.

## (b) No reader needs the class before `loadModelHeader`

`loadModelHeader` is `run`'s FIRST step (`internal/table/load.go:79`).

- `normalizeRule` runs under step 9, `normalizeRules`
  (`internal/table/load.go:87`) — strictly after step 1. Confirmed by the
  literal slice order in `run`.
- The three graphlint readers run only after `table.LoadWithAdvisories`
  has RETURNED a `*Model` (`internal/cli/lint.go:99` → `:115`), i.e. after
  every one of `run`'s ten steps completed. `run` returns `l.model` only at
  `internal/table/load.go:94`, after the loop.

No reader runs at or before `loadModelHeader`.

## (c) Graphlint receives the full `*table.Model`, not a projection

This was the load-bearing risk and it is CLEAR.

`internal/graphlint/engine.go:13-15`:

```go
type Request struct {
	Model *table.Model
}
```

`graphlint.NewRequest` is the sole request builder (`engine.go:19`) and it
stores the pointer verbatim. `Model.KernelTable`
(`internal/table/model.go:398`) — the one projection on `Model`, which
would indeed drop a `Class` field — is never called on the lint path. Its
only non-test callers are `internal/cli/flow_resolve.go:108` and
`internal/cli/flow_resolve.go:208`, both on the resolve/exec path, not
lint:

```
$ grep -rn "graphlint.NewRequest|graphlint.Run|KernelTable()" internal cmd | grep -v _test.go
internal/cli/flow_resolve.go:108:  Table:      req.model.KernelTable(),
internal/cli/flow_resolve.go:208:  probe := req.model.KernelTable()
internal/cli/lint.go:115:          report := graphlint.Run(graphlint.NewRequest(m))
internal/table/model.go:398:      func (m *Model) KernelTable() resolve.Table {
```

An exported field on `Model` is therefore sufficient for the graphlint
readers; no projection needs to carry it.

## (2) The C1 agreement-check window is non-empty and correctly stocked

`internal/table/load.go:75-94`, verbatim step slice:

```go
func (l *loader) run() (*Model, error) {
	l.model = &Model{}

	for _, step := range []func() error{
		l.loadModelHeader,
		l.loadOutcomes,
		l.loadTags,
		l.loadAccessors,
		l.loadDump,
		l.loadContexts,
		l.loadInitial,
		l.loadTerminal,
		l.normalizeRules,
		l.checkAccessorBindings,
	} {
		if err := step(); err != nil {
			return nil, err
		}
	}
	return l.model, nil
}
```

Matches C1's stated order exactly. `checkAccessorBindings` is the last
step, so the window `[loadTags .. before checkAccessorBindings]` spans six
insertion points (after `loadTags`, `loadAccessors`, `loadDump`,
`loadContexts`, `loadInitial`, `loadTerminal`, `normalizeRules`).

Step signature is uniform `func() error` — a new step is a
`func (l *loader) checkClassTagAgreement() error` added to the slice, no
signature change.

Shared state available in the window:

- Class: `loadModelHeader` (`internal/table/load.go:99`) mutates
  `l.model.ID`, `l.model.Version`, `l.model.Description`, `l.model.Metadata`
  (`:111-116`); `l.model.Class = ...` would sit alongside them. Set at step
  1, so available at every later step.
- Tags: `loadTags` (`internal/table/load.go:147`) assigns
  `l.model.Tags = decls` at `:201` as its final statement. So
  `l.model.Tags` (a `map[string]TagDecl`, each carrying
  `Provenance`) is populated from step 3 onward, and the owned count is
  `len(k for k, d := range l.model.Tags if d.Provenance == ProvenanceOwned)`.

Both halves of the agreement check are on `l.model` for any step in the
window. C1's claim that `l.model.Tags` is empty during `loadModelHeader` is
also confirmed: `loadTags` is step 3 and nothing before it writes `Tags`.

`checkAccessorBindings` (`internal/table/load.go:360`) reads `l.model.Readers`,
`l.model.Tags`, `l.model.Writers`, `l.model.Rows`, `l.model.Initial` — it
does not touch a class, so inserting a class check before it does not
disturb it, and the ceiling is what keeps the doubly-malformed case
determinate as C1 argues.

## Model shape confirmation

`internal/table/model.go:369-395` — `Model` carries thirteen exported
fields (`ID`, `Version`, `Description`, `Metadata`, `Outcomes`, `Initial`,
`Terminal`, `Tags`, `Readers`, `Writers`, `Gates`, `DumpOrder`, `Rows`) and
no field accessor. Its only method is `KernelTable`
(`internal/table/model.go:398`), a projection. An exported `Class string`
is the type's existing convention.

## Verdict

**PASS.** All four readers reach a `*table.Model` (or the `*loader` that
owns it) at their own site; none runs before `loadModelHeader`; graphlint
takes the full `*table.Model` pointer and never the `KernelTable`
projection, so an exported field is sufficient; and the C1 window is a real
six-slot span with both the class (step 1) and the tags (step 3) present in
shared `l.model` state.
