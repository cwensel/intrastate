package graphlint_test

// RDR 0021 — the `graphlint` surface the export requires (`0021:C4`, A8).
//
// C4 needs completeness to reach the verb through an EXPORTED surface:
// today's `Reach` discards the private `reach`'s `complete` bool, so the
// edges-carrying function adds beside it returns nodes, edges AND
// completeness. A8 records the addition as ADDITIVE — `Reach` and `reach`
// keep their present signatures and lint's path through them is untouched.
//
// That function does not exist yet, so naming it as a Go identifier would
// not compile and `.githooks/pre-commit` (`gofmt -l`, `go vet ./...`, under
// `set -e`) would refuse the commit. These probes therefore parse the
// package SOURCE with `go/parser`, which reads the tree as text: an absent
// function is a `false` return rather than a build break, and every failure
// here is a runtime assertion.
//
// Helper names carry a `graphExport` prefix because this package's existing
// test files already own `render`, `source`, `mustLoad`, `withCode`,
// `requireCode`, `nodeString`, `repoRoot` and friends.

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/graphlint"
	"github.com/cwensel/intrastate/internal/table"
)

// graphExportPkgDir locates this package's own source directory.
func graphExportPkgDir(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return dir
}

// graphExportFuncs returns the exported package-level functions declared in
// the production source of the package rooted at dir.
func graphExportFuncs(t *testing.T, dir string) map[string]*ast.FuncDecl {
	t.Helper()

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Fatalf("parsing %s: %v", dir, err)
	}

	out := map[string]*ast.FuncDecl{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Name == nil || fn.Recv != nil {
					continue
				}
				if fn.Name.IsExported() {
					out[fn.Name.Name] = fn
				}
			}
		}
	}
	return out
}

// graphExportSignature renders a function's parameter and result types as
// source text, so a declaration can be compared without linking to it.
func graphExportSignature(t *testing.T, fn *ast.FuncDecl) (params, results []string) {
	t.Helper()

	fset := token.NewFileSet()
	render := func(expr ast.Expr) string {
		var b strings.Builder
		if err := printer.Fprint(&b, fset, expr); err != nil {
			t.Fatalf("rendering a type: %v", err)
		}
		return b.String()
	}
	if fn.Type.Params != nil {
		for _, field := range fn.Type.Params.List {
			n := max(len(field.Names), 1)
			for range n {
				params = append(params, render(field.Type))
			}
		}
	}
	if fn.Type.Results != nil {
		for _, field := range fn.Type.Results.List {
			n := max(len(field.Names), 1)
			for range n {
				results = append(results, render(field.Type))
			}
		}
	}
	return params, results
}

// REQ-53: "Completeness MUST reach the verb through an EXPORTED `graphlint`
// surface: today's `Reach` discards the private `reach`'s `complete` bool
// (`nodes, _ := reach(m)`), so the edges-carrying function Q3(c) already
// adds beside it returns completeness with the nodes and edges."
// ASSUMPTION A-7: the exported completeness surface is ONE new function,
// not two.
// HAPPY PATH
func TestReq53_AnExportedFunctionCarriesNodesEdgesAndCompleteness(t *testing.T) {
	funcs := graphExportFuncs(t, graphExportPkgDir(t))

	// The Go spelling is not fixed by the record, so the probe accepts the
	// names a reasonable implementation would choose and fails when NONE is
	// declared. Pinning one spelling would make this a naming test rather
	// than a surface test.
	candidates := []string{
		"ReachWithEdges", "ReachGraph", "ReachRelation", "ReachEdges",
		"Relation", "Graph", "Export",
	}
	var found string
	var decl *ast.FuncDecl
	for _, name := range candidates {
		if fn, ok := funcs[name]; ok {
			found, decl = name, fn
			break
		}
	}
	if found == "" {
		var have []string
		for name := range funcs {
			have = append(have, name)
		}
		t.Fatalf("internal/graphlint declares no exported edges-and-"+
			"completeness function (looked for %v; it exports %v). C4 "+
			"requires completeness reach the verb through an EXPORTED "+
			"surface: `Reach` discards the private `reach`'s `complete` "+
			"bool, so without this function the `%s` refusal is "+
			"unimplementable without widening `Reach`'s signature — which "+
			"A8 and C4's neutrality rule both forbid",
			candidates, have, "graph-export-too-large")
	}

	// It carries all three: nodes, edges, and completeness. A separate
	// completeness accessor would be two functions where A8 describes one.
	_, results := graphExportSignature(t, decl)
	if len(results) < 3 {
		t.Errorf("%s returns %v; A8 describes ONE new exported function "+
			"carrying nodes, edges, AND completeness — C4 has it \"return "+
			"completeness with the nodes and edges\"", found, results)
	}
	var hasBool bool
	for _, r := range results {
		if r == "bool" {
			hasBool = true
		}
	}
	if !hasBool {
		t.Errorf("%s returns %v, carrying no boolean completeness; the raw "+
			"typed bool is what the new function carries, and it is the "+
			"only thing `Reach` discards today", found, results)
	}
}

// REQ-56: "Add a sibling exported function (`Reach` and `reach()` both
// untouched)"
// REQ-100: "A8 verified — the exported completeness-carrying surface
// exists. It gates Phase 2"
// ADVERSARIAL — the addition is ADDITIVE. Widening `Reach`'s signature
// instead would change a surface whose other callers are tests, which is
// what A8 and C4's neutrality rule forbid.
func TestReq56And100_ReachKeepsItsPresentSignature(t *testing.T) {
	funcs := graphExportFuncs(t, graphExportPkgDir(t))

	reach, ok := funcs["Reach"]
	if !ok {
		t.Fatal("internal/graphlint no longer exports `Reach`; A8 makes the " +
			"new surface ADDITIVE — `Reach` and `reach` are both untouched")
	}
	params, results := graphExportSignature(t, reach)

	if len(params) != 1 || params[0] != "*table.Model" {
		t.Errorf("Reach takes %v; its present signature is "+
			"`Reach(*table.Model)`, which A8 leaves untouched", params)
	}
	if len(results) != 1 || results[0] != "[]Node" {
		t.Errorf("Reach returns %v; its present signature returns `[]Node`. "+
			"Widening it to carry edges or completeness is exactly the "+
			"change A8 rules out: its other callers are tests, and C4's "+
			"neutrality rule bears on the traversal surface", results)
	}

	// And it still behaves: the additive change leaves the relation itself
	// alone, which is the half a signature check cannot see.
	m, err := table.Load([]byte(graphExportFixture), "fixture.toml")
	if err != nil {
		t.Fatalf("fixture must load clean: %v", err)
	}
	if len(graphlint.Reach(m)) == 0 {
		t.Error("Reach returns no nodes for a model with a declared root; " +
			"the additive surface must leave the existing relation intact")
	}
}

// REQ-54: "The ceiling this refusal names is `reach`'s node-count
// completeness — the SAME bound `analysis.go`'s `checkNodeCeiling` already
// fires on (both read one `reach()` and one `nodeCeiling` constant,
// published as `graphlint.NodeCeiling()`), observed at two call sites"
// REQ-55: "The distinct bound this refusal does NOT involve is the
// guard-product one, `graphlint.ProductBound()` (`guard.Bound()`)"
// BOUNDARY — the two published bounds are distinct values from distinct
// sources, so a refusal naming one cannot be satisfied by the other.
func TestReq54And55_TheTwoPublishedBoundsAreDistinct(t *testing.T) {
	ceiling := graphlint.NodeCeiling()
	product := graphlint.ProductBound()

	if ceiling <= 0 {
		t.Errorf("NodeCeiling() = %d; it is a positive implementation "+
			"constant this refusal names", ceiling)
	}
	if product <= 0 {
		t.Errorf("ProductBound() = %d; it is a positive implementation "+
			"constant", product)
	}
	if ceiling == product {
		t.Errorf("NodeCeiling() and ProductBound() are both %d; C4 names the "+
			"NODE-COUNT completeness bound and explicitly NOT the "+
			"guard-product one, so a test cannot tell which bound a "+
			"refusal named while the two coincide", ceiling)
	}

	// Both are model-independent: repeated calls agree.
	if graphlint.NodeCeiling() != ceiling || graphlint.ProductBound() != product {
		t.Error("a published bound varied between calls; both are " +
			"model-independent implementation constants")
	}
}

// graphExportFixture is a minimal loadable state machine with a declared
// root, so `Reach` has something to traverse.
const graphExportFixture = `
outcomes = ["go"]
terminal = ["done"]

[model]
id = "export-probe"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[read.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"
read_back = true

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "advance-otherwise"
escape = ["no_match"]
[rule.match.recognized]
eq = "go"
`

// graphExportRepoRoot walks up to the module root, for the probes that read
// files outside this package.
func graphExportRepoRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, serr := os.Stat(filepath.Join(dir, "go.mod")); serr == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("no go.mod above the test working directory")
		}
		dir = parent
	}
}

// REQ-48 (the shared-entry-surface half, from the traversal's side).
// REQ-93: "Reachability relation | `reach.go::reach` (private) | lint via
// `graphlint.Run`; export via the new sibling exporter | `analysis.go:46`
// (lint); the new exporter | … | `reach()` — both arms reach the one
// traversal (C4)"
// ADVERSARIAL — the defect is a second traversal implementation for the
// export, which would let the exported relation drift from the certified
// one. The observable: the package declares exactly ONE fixpoint loop.
func TestReq48And93_ThePackageDeclaresOneTraversalNotTwo(t *testing.T) {
	src, err := os.ReadFile(filepath.Join(
		graphExportRepoRoot(t), "internal", "graphlint", "reach.go"))
	if err != nil {
		t.Fatalf("reading reach.go: %v", err)
	}

	// `reach` is the one private fixpoint. A second private traversal
	// declared beside it is the drift C4 forbids structurally.
	if got := strings.Count(string(src), "func reach("); got != 1 {
		t.Errorf("reach.go declares %d `reach(` functions; C4 has BOTH arms "+
			"reach the ONE traversal, so verb/lint drift is structural "+
			"rather than disciplinary", got)
	}
}
