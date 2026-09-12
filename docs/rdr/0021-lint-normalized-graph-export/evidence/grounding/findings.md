Model: claude-sonnet-5

## Refuted / Not-Found

None. Every claim swept (source-anchored and unanchored) confirmed against `main`.

## Inverse-rule sweep (new discriminators vs. existing siblings)

- New root verb `graph`: searched `internal/cli/root.go::NewRootCmd` — only `newVersionCmd`, `newLintCmd`, `newFlowCmd`, `newDocsCmd` registered. No existing `graph.go` or export/emit verb in production code (only test-file string mentions of "graph"). No sibling exists — confirms RDR's own claim.
- New `--emit` format-selector flag: searched all `cmd.Flags().String(...)` registrations across `internal/cli/*.go` (non-test) — `root.go:188` registers `--as` (`respond.FlagName`) as the only output-mode/format flag in the whole CLI; `lint.go`/`flow.go` register `--model`/`--flow` (selection, not format). No existing format-selector flag anywhere else. Confirms decision-rationale's "sibling-path check: none exists" claim (`§decision-rationale`).
- New schema-versioned document (`intrastate.graph/1`): grepped for `"schema":` / versioned-document literal patterns repo-wide (non-test Go) — no hit. No existing sibling document-versioning convention to point to; genuinely new.
- New DOT rendering: `find` for `*dot*.go` in the repo (excluding the external `state-machines` corpus) — no hit. No existing DOT/graphviz writer package anywhere in intrastate. The only DOT precedent is the external `state-machine-cat` prior art, correctly cited as external, not as an in-repo sibling.
- New selection predicate (any model that loads, including one lint refuses): `graphlint.Reach`/`reach` and `lint.go::runLint` share the same load path; no existing verb currently exports a load-succeeds-lint-fails model, so this is genuinely new behavior, not a rediscovery of an existing arm.

## Confirmed (coverage only, no findings needed)

- A1 TextLiner seam: `internal/cli/respond/respond.go::OK` (lines 165-167) — `if liner, ok := s.Data.(TextLiner); ok { fmt.Fprintln(cmd.OutOrStdout(), liner.TextLine()) }`; notes/warnings written to `cmd.ErrOrStderr()` (lines 141-147). Matches "verbatim, no decoration, notes/warnings on stderr."
- A2 `reach`/`successorsOf`/`indexOf`/`subsumes`: `internal/graphlint/reach.go` — `reach` (90), `successorsOf` (200), `indexOf` (249), `subsumes` (262) all exist with claimed shapes; `indexOf` calls `subsumes` (251).
- A3 model field walk: `internal/table/model.go::Model` struct (496) carries `ID, Class, Tags, Initial, Terminal` etc.; `internal/guard/declaration.go::AssignmentCount` (80) and `internal/guard/product.go::Groups` (66) exist as claimed.
- A4 RDR 0005:C1 scope: read `0005:C1` directly — opening sentence scopes the MUST-clauses to the `flow` skill-integration group (`next, resolve, read-state, set-state`); carve-out sentence names "lint, dump, parse" as a closed enumeration not including `graph`. Matches RDR's reading exactly.
- A5 `WriteJSONLine` + ad hoc encoder count: `internal/cli/clierr/clierr.go:174-176` (`enc.SetEscapeHTML(false)`); grepped `SetEscapeHTML` repo-wide (non-test) — exactly 5 hits: `clierr.go` (the shared one) plus 4 ad hoc sites (`internal/table/model.go:466`, `internal/cli/cmdbind/cmdbind.go:1097`, `internal/cli/flow_input.go:800`, `internal/cli/respond/text.go:134`). Matches "four ad hoc call sites" exactly.
- C1 SURFACE: verb registration site `internal/cli/root.go::NewRootCmd` (173, `AddCommand(newLintCmd())` at 195); arm codes verbatim in `internal/cli/lint.go::runLint` — `flag-mutually-exclusive` (153), `flag-required` (173), `model-unreadable` (183), `model-invalid` (202); `graph-export-too-large` grepped repo-wide — no existing use, confirmed new/uncollided.
- C4 NEUTRALITY / shared entry surface: `internal/cli/lint.go::runLint` calls `graphlint.Run(graphlint.NewRequest(m))` (212); `graphlint.Run` → `internal/graphlint/analysis.go:46` calls the same private `reach(m)` that `graphlint.Reach` (public wrapper, `reach.go:73-76`) also calls. The "same entry surface" claim resolves to a real shared call chain, not merely a shared package name.
- C5 MODE COEXISTENCE / `version`'s TextLine precedent: the RDR cites "`version`'s `TextLine`" without a symbol — resolved to `internal/version/version.go:130::Info.TextLine` (`func (i Info) TextLine() string { return i.String() }`), consumed by `internal/cli/version.go::newVersionCmd` via `respond.OK(cmd, respond.Success{Data: version.Get()})`. Confirmed as the cited identity-string precedent.
- `0006:C19`/`C20` mirrored codes and `0006:§approach` shared-engine precedent: read directly, matches RDR's paraphrase (single authoritative surface, route through `respond.OK`/`Fail`/`CLIError`, no direct stdout/stderr writes).
- `0002:§round-trip-inverse-invariants` lossy set-literal rendering: read directly — unescaped separators make set-valued atom literals non-recoverable in dump text; matches C2's "closing this document's lossy rendering" claim.
- JDR 0002 §D1 exists at `docs/jdr/0002-success-envelope-projection.md:40` ("D1 — The caller-projection partition doctrine"), cited correctly by C5.
- External prior art: `state-machines/repos/state-machine-cat/README.md:66` — `-T --output-type <type> svg|eps|ps|ps2|dot|smcat|json|ast|scxml|oldsvg|scjson|pdf|png` — confirms the cited single-flag-selects-sibling-formats pattern (json and dot as siblings, no JSON-wrapped-DOT).

## Widened past scoped spans

- Read `internal/graphlint/analysis.go:46` (not in scope.json) to verify C4's "same entry surface" claim at the call-flow level rather than trusting the package name — this is the call site the RDR's neutrality claim depends on and the record cites no symbol for it.
- Read `internal/cli/version.go` and `internal/version/version.go` (not in scope.json) to resolve C5's unnamed "`version`'s `TextLine`" citation to a concrete symbol.
- Read `internal/graphlint/reach.go:60-100` beyond the single-line anchor to confirm the `complete bool` return this RDR's C4 ceiling-refusal mechanism depends on.
- Grepped repo-wide for `SetEscapeHTML`, `"schema":`, `graph-export-too-large`, `--emit`, `*dot*.go`, and `Flags().String(` to run the inverse-rule sweep (Task 2) and confirm A5's exact call-site count — none of these are scope.json edges since they are negative/counting claims that mint no anchor.
