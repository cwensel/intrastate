model: claude-sonnet-5
variant: lite (profile: large)

# Repeatability reconstruction — RDR 0021 (lint normalized-graph export), run 1

Grounded against the RDR contracts (C1-C5, MVV, S1-S9, RT1-RT3, D-identity/D-wire-byte-format/D-naming/D-selection-predicate, §illustrative-code, §pre-lock-mini-checks, §implementation-plan, §decision-rationale) AND the actual `internal/graphlint` package on disk (`reach.go`, `analysis.go`, `taxonomy.go`), `internal/cli/lint.go::runLint`, and `internal/cli/respond`, `internal/cli/clierr`, `internal/guard` packages, since the RDR names these symbols directly and they already exist pre-implementation. Widened past C1-C5 into §pre-lock-mini-checks (`authority`/`disposition` tables, 600-662) and §implementation-plan (Phase 1-4, 973-994) to resolve the exporter's internal decomposition, which the contracts alone leave as a silence (the RDR fixes the *value* and the *wire*, not the exporter's internal call graph).

## 1. Public API

```go
// internal/cli/graph.go — new sibling verb, mirrors lint.go's arm set (C1).

// newGraphCmd registers the root verb `graph` beside `lint`. GUESS: exact
// constructor name — the RDR never spells it; mirrors lint.go's
// newLintCmd() naming convention 1:1.
func newGraphCmd() *cobra.Command

// runGraph is graph's RunE. GUESS: name only (mirrors runLint); the
// control flow below is NOT a guess — C1 fixes it verbatim against
// lint.go::runLint's literal arm order.
func runGraph(cmd *cobra.Command, args []string) error

// Flags (C1, D-naming):
//   --model <path>       string, mutually exclusive with --flow
//   --flow <id>           string, mutually exclusive with --model;
//                          this build resolves no ids (flag-invalid-value)
//   --emit <format>       string, "json" (default) | "dot" (D-selection-predicate)
//   --as <mode>            persistent root flag from respond.FlagName,
//                          "text" (default) | "json" (C5)

// Error codes (verbatim from lint.go::runLint, C1):
//   both/neither --model/--flow -> "flag-mutually-exclusive" / "flag-required"
//   --flow alone                -> "flag-invalid-value" (Param: "flow")
//   unreadable model file        -> "model-unreadable"
//   load failure                 -> "model-invalid" (Findings: one per load category)
//   unknown --emit value         -> "flag-invalid-value" (Param: "emit") — GUESS: Param
//     spelling "emit" is read from C1's "naming `emit`" phrase; the RDR
//     does not show the literal CLIError{} construction for this arm.
//   traversal incomplete         -> "graph-export-too-large" (GroupUserEnv,
//                                    exit 2), Detail/Hint name the ceiling
//                                    and "narrow a declared domain" remedy (C4)
// All other arms: existing 0/2 exit mapping via respond.OK / respond.Fail.

// internal/graphlint — extends the existing package (already on disk).

// Existing, unchanged (grounded from internal/graphlint/reach.go):
func Reach(m *table.Model) []Node                 // public wrapper, tests only
func reach(m *table.Model) (nodes []Node, complete bool)  // private, unexported

type Node struct {
    Values map[string][]string
}
const OpaqueValue = "<opaque>"
const nodeCeiling = 4096   // internal/graphlint/taxonomy.go:136
func NodeCeiling() int     // GUESS-shaped export: C4 requires "an EXPORTED
                            // graphlint surface" publishing this constant;
                            // the RDR names it graphlint.NodeCeiling() by
                            // literal reference but the export function
                            // itself does not exist in the current package
                            // (only the unexported const does) — so the
                            // wrapper's existence is fixed by C4, its exact
                            // body (`return nodeCeiling`) is a GUESS.
func ProductBound() int    // already exists as guard.Bound(); C4 says the
                            // export's refusal path never touches this —
                            // included here only to mark the boundary.

// NEW — the edge-carrying completeness surface C4 requires (Q3(c)):
//   "the edges-carrying function Q3(c) already adds beside it returns
//    completeness with the nodes and edges" — RDR names the requirement,
//   not the signature. GUESS for exact name/shape, built by symmetry with
//   the existing reach()/Reach() split:
type Edge struct {
    From string // Node.key() of source
    To   string // Node.key() of successor
    Rule string // row identity of the edge-producing rule
}

// ReachWithEdges is the new exported entry point export uses instead of
// Reach(); GUESS name (RDR calls it only "the new exported function" /
// "the new sibling exporter" / "post-hoc recovery over final nodes").
func ReachWithEdges(m *table.Model) (nodes []Node, edges []Edge, complete bool)

// internal/graphlint/export.go — NEW, the document-assembly component
// (§implementation-plan Phase 1: "Build the document-assembly component
// over *table.Model + Reach-with-edges"). GUESS for exact file/type names;
// the RDR fixes the wire shape (C2) and the entry point's inputs/outputs,
// not Go identifiers here.

type Document struct {
    Schema   string            `json:"schema"`             // "intrastate.graph/1"
    Model    string             `json:"model"`
    Class    string             `json:"class"`
    Tags     []TagDoc           `json:"tags"`
    Initial  []Assignment       `json:"initial"`
    Terminal [][]Atom           `json:"terminal"`
    Rows     []RowDoc           `json:"rows"`
    Groups   []GroupDoc         `json:"groups"`
    Reach    ReachDoc           `json:"reach"`
}

type TagDoc struct {
    Name         string   `json:"name"`
    Provenance   string   `json:"provenance"`
    Kind         string   `json:"kind"`
    Required     bool     `json:"required"`
    SingleValued bool     `json:"single_valued"`
    Domain       []string `json:"domain,omitempty"` // present iff guard.AssignmentCount finite
}

type ReachDoc struct {
    Abstraction string       `json:"abstraction"` // "declared-over-approximation"
    Nodes       []NodeDoc    `json:"nodes"`        // sorted by node key
    Edges       []EdgeDoc    `json:"edges"`        // sorted by (from, to, rule)
}
type NodeDoc struct {
    ID     string              `json:"id"`
    Values map[string][]string `json:"values"`
}
type EdgeDoc struct {
    From string `json:"from"`
    To   string `json:"to"`
    Rule string `json:"rule"`
}

// BuildDocument assembles the wire document from a loaded model.
// GUESS name/signature; C2 fixes every field this returns, not the Go
// function that produces it.
func BuildDocument(m *table.Model) (Document, error)

// internal/graphlint/dot.go — NEW, Phase 3's "pure renderer over the
// export value". GUESS name/signature; C2/S6 fix its output's node/edge
// SET and the abstraction-marker placement, not the function shape.
func RenderDOT(doc Document) (string, error)
```

## 2. Three most important internal helpers

1. **`ReachWithEdges` / edge recovery over final nodes** (`internal/graphlint`, new). Responsibility: re-run (or wrap) the existing `reach()` fixpoint and additionally recover the edge relation `(from, to, rule)` by re-deriving `successorsOf` transitions against the FINAL merged node set, rather than recording edges during traversal (A2's rejected in-traversal-observer alternative). This is the one piece of new traversal-adjacent logic; C4 requires it stay the "SAME package-internal traversal lint reaches" so lint/export drift is structural. GUESS on internal decomposition (whether it re-walks `successorsOf` post-fixpoint or threads an observer some other way) — the RDR fixes only the equality bar ("post-hoc recovery; equality to the traversal is the bar (A2)"), not the recovery mechanism's shape.

2. **`BuildDocument`** (`internal/graphlint`, new). Responsibility: project a loaded `*table.Model` plus the `ReachWithEdges` result into the `Document` value — mapping tag declarations (with `guard.AssignmentCount` deciding `Domain` presence), initial/terminal declarations, RDR 0002's canonical row/atom order, selection-context groups, and the sorted/deduplicated reach node-values object — while emitting `[]`/`{}` never `nil` for empty declared collections and omitting (never null-ing) inapplicable optional members. This is the sole point translating internal types to wire types, so it is where C2's field-spelling and empty-vs-absent rules live.

3. **`runGraph`'s arm dispatch** (`internal/cli/graph.go`, new). Responsibility: mirror `runLint`'s exact selection-arm order and error codes (mutual exclusion -> flag-required -> file-read -> load -> now, uniquely to this verb, `--emit` validation and the `graph-export-too-large` refusal keyed off `complete`), then route success through `respond.OK` with `Data` set to either the JSON `Document` or an object `{"dot": <string>}` per C5's mode/format 2x2, and refuse via `respond.Fail` with never a partial document on stdout on any refusing arm.

## 3. Data model (persisted / passed across the boundary)

The wire format IS the persisted/exported artifact (C2); nothing else crosses this boundary (no config file, no cache). Canonical JSON shape (grounded verbatim in §illustrative-code and C2's field list):

```json
{
  "schema": "intrastate.graph/1",
  "model": "<model identity>",
  "class": "<state-machine|decision-table>",
  "tags": [
    {"name": "...", "provenance": "...", "kind": "...",
     "required": true, "single_valued": true, "domain": ["..."]}
  ],
  "initial": [{"key": "...", "value": ["..."]}],
  "terminal": [[{"key": "...", "operator": "eq", "literal": ["..."], "block": "match"}]],
  "rows": [
    {"identity": "...", "source": "...", "kind": "...", "outcome": "...",
     "atoms": [{"key": "...", "operator": "...", "literal": ["..."], "block": "..."}],
     "next": "...", "writes": [{"key": "...", "value": ["..."]}],
     "requires_owned": ["..."], "gate": "...", "escape": false, "emit": "..."}
  ],
  "groups": [{"context": "...", "rules": ["..."]}],
  "reach": {
    "abstraction": "declared-over-approximation",
    "nodes": [{"id": "<node key>", "values": {"tagname": ["v1", "v2"]}}],
    "edges": [{"from": "<node key>", "to": "<node key>", "rule": "<row identity>"}]
  }
}
```

Ordering invariants (C3, C2): `rows` and each row's `atoms` in RDR 0002's canonical order; `reach.nodes` sorted by node key; `reach.edges` sorted by `(from, to, rule)`; every set-valued `values[tag]` array sorted and deduplicated. A tag with no finite declared domain carries `["<opaque>"]` verbatim (not an authored value).

DOT sibling artifact (not JSON, C2/C5): plain text, one node per reach node, one edge per reach edge labeled with `rule`, initial/terminal nodes marked, abstraction marker in header comment + graph label; under `--as=json` it rides inside `{"dot": "<string>"}` as the single required member.

`rows[]` field list is "RDR 0002's dump field list" — GUESS on exact JSON key spellings for `source`, `next`, `requires_owned`, `gate`, `escape`, `emit` (C2 names these as English labels: "identity, source, kind, outcome, atoms, next, writes, requires_owned, gate, escape, emit" but only gives the full worked JSON shape for `identity/kind/outcome/atoms/writes` in §illustrative-code's partial example). snake_case chosen by symmetry with `single_valued`/`requires_owned` (the one multi-word field C2 does spell, confirming snake_case), but the remaining unshown fields are a GUESS.

## 4. Top-level pseudo-code of the main operation (~30 lines)

```
func runGraph(cmd, args):
    if ce := respond.ValidateMode(cmd); ce != nil:
        return respond.Fail(cmd, ce)

    // Selection arms — verbatim from lint.go::runLint (C1)
    if model_flag_set && flow_flag_set:
        return respond.Fail(cmd, CLIError{Code: "flag-mutually-exclusive"})
    if flow_flag_set:
        return respond.Fail(cmd, CLIError{Code: "flag-invalid-value", Param: "flow"})
    if !model_flag_set:
        return respond.Fail(cmd, CLIError{Code: "flag-required"})

    bytes, err := os.ReadFile(modelPath)
    if err != nil:
        return respond.Fail(cmd, CLIError{Code: "model-unreadable"})

    model, loadErr := table.Load(bytes)   // GUESS: exact loader call
    if loadErr != nil:
        return respond.Fail(cmd, CLIError{Code: "model-invalid", Findings: ...})

    emit := cmd.Flags().GetString("emit")
    if emit != "json" && emit != "dot":
        return respond.Fail(cmd, CLIError{Code: "flag-invalid-value", Param: "emit"})

    nodes, edges, complete := graphlint.ReachWithEdges(model)
    if !complete:
        return respond.Fail(cmd, CLIError{
            Code: "graph-export-too-large", Group: GroupUserEnv,
            Detail: fmt.Sprintf("ceiling %d", graphlint.NodeCeiling()),
            Hint: "narrow a declared domain",
        })

    doc := graphlint.BuildDocument(model, nodes, edges)  // C2 assembly,
                                                          // deterministic sort order (C3)

    switch emit:
        case "json":
            payload = doc
        case "dot":
            dotText, _ := graphlint.RenderDOT(doc)
            payload = dotText   // text mode: raw string; json mode: wrapped below

    // Mode fan-out (C5) — one export value, two renderings
    if asMode == "text":
        clierr.WriteJSONLine(stdout, doc)      // --emit json: raw document
        // OR: fmt.Fprintln(stdout, dotText)   // --emit dot: raw DOT + \n
        return nil
    // asMode == "json"
    data := payload
    if emit == "dot":
        data = map[string]string{"dot": dotText}
    return respond.OK(cmd, respond.Success{Data: data})
```

GUESS markers: `table.Load` call name/shape in the load step (not spelled in the RDR — only the resulting error codes are normative); the exact `payload`/`data` variable wiring for the text-vs-json branch is my reconstruction of C5's "one export value, two renderings" rule, not a literal code excerpt in the record.

## Widened spans

Widened past C1-C5/MVV/S1-S9 into: §pre-lock-mini-checks (600-662, `authority`/`oracle`/`disposition`/`fidelity`/`trace` tables) to resolve the exporter's writer/reader/call-site decomposition and exact disposition-table exit/code pairing; §implementation-plan (941-994) for the Phase 1-4 breakdown (document-assembly component, CLI wiring, DOT renderer, contract-surface docs) and the A8-gates-Phase-2 dependency; §decision-rationale (682-735) for the O1 vs O2/O3/O4 scoring that explains why `graph` is a sibling verb and not a `lint --emit` flag. Also read the live `internal/graphlint/reach.go`, `analysis.go`, `taxonomy.go`, `internal/cli/lint.go`, `internal/cli/respond/respond.go`, `internal/cli/clierr/clierr.go`, and `internal/guard/{declaration,product}.go` on disk (pre-implementation, but already present) to ground real Go signatures rather than invent them; every signature reused from these files is NOT a guess, only the new symbols layered on top are.

## GUESS summary

12 explicit GUESS markers, clustering on: (a) exact new-symbol names (`newGraphCmd`, `runGraph`, `ReachWithEdges`, `Edge`, `BuildDocument`, `RenderDOT`, `NodeCeiling` wrapper body) since C1-C4 mandate the *capability* but never spell the Go identifier; (b) the `rows[]` snake_case field spellings C2 lists in English but never shows in full JSON; (c) the `--emit` unknown-value CLIError's `Param` field; (d) the internal edge-recovery mechanism's shape (post-fixpoint re-walk vs. other), which the RDR deliberately leaves to A2/Resolve per C4's "mechanism-independent" oracle language.
