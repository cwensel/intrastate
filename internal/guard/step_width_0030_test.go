package guard_test

// RDR 0030 `0030:C1` — the ONE int-width rule shared by the loader and lint,
// and the byte agreement between the loader's emitted cell and the cells
// lint enumerates.
//
// The width rule's new home is a symbol the tree may not yet declare, so it
// is asserted through the package SOURCE (`go/parser`): naming it directly
// would break the build the pre-commit `go vet ./...` gates.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/guard"
	"github.com/cwensel/intrastate/internal/table"
)

// s30Funcs parses a package directory's non-test sources and returns every
// function and method declaration, keyed by name.
func s30Funcs(t *testing.T, dir string) map[string][]*ast.FuncDecl {
	t.Helper()

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading %s: %v", dir, err)
	}
	fset := token.NewFileSet()
	out := map[string][]*ast.FuncDecl{}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(dir, name), nil, 0)
		if err != nil {
			t.Fatalf("parsing %s: %v", name, err)
		}
		for _, d := range f.Decls {
			if fn, ok := d.(*ast.FuncDecl); ok {
				out[fn.Name.Name] = append(out[fn.Name.Name], fn)
			}
		}
	}
	return out
}

// s30Callees returns what a body calls: same-package names (functions, and
// method names for selector calls) and `pkg.Name` for package selectors.
func s30Callees(body *ast.BlockStmt, pkgs ...string) []string {
	var out []string
	if body == nil {
		return out
	}
	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		switch fn := call.Fun.(type) {
		case *ast.Ident:
			out = append(out, fn.Name)
		case *ast.SelectorExpr:
			if id, ok := fn.X.(*ast.Ident); ok {
				for _, p := range pkgs {
					if id.Name == p {
						out = append(out, p+"."+fn.Sel.Name)
						return true
					}
				}
			}
			out = append(out, fn.Sel.Name)
		}
		return true
	})
	return out
}

// s30Reaches reports whether start (a declaration name in funcs) reaches the
// qualified callee target, following same-package calls transitively.
func s30Reaches(funcs map[string][]*ast.FuncDecl, start, target string) bool {
	seen := map[string]bool{start: true}
	queue := []string{start}
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		for _, fn := range funcs[name] {
			for _, c := range s30Callees(fn.Body, "resolve") {
				if c == target {
					return true
				}
				if _, local := funcs[c]; local && !seen[c] {
					seen[c] = true
					queue = append(queue, c)
				}
			}
		}
	}
	return false
}

// s30IsIdent reports whether e is the bare identifier name.
func s30IsIdent(e ast.Expr, name string) bool {
	id, ok := e.(*ast.Ident)
	return ok && id.Name == name
}

// REQ-4: "the rule MOVES to the shim in the EXISTING `internal/resolve`
// package beside the relocated `Evaluator` rather than being copied" … "it is
// EXPORTED at its new home as `resolve.IntWidth(minV, maxV int) (int, bool)`."
// … "`::IntDomain(d table.TagDecl) []int` does NOT move" … "`guard`'s
// in-package re-inline (`assignment.go::valueAssignments` walks the same
// counted loop) and `::domainSize` both repoint through the moved
// `resolve.IntWidth`." Observable: exactly one width rule in the module.
// BOUNDARY
func TestReq4_0030_OneIntWidthRuleLivesInResolveAndGuardReachesIt(t *testing.T) {
	resolveFuncs := s30Funcs(t, filepath.Join("..", "resolve"))
	var found bool
	for _, fn := range resolveFuncs["IntWidth"] {
		if fn.Recv != nil {
			continue
		}
		found = true
		var params []ast.Expr
		for _, p := range fn.Type.Params.List {
			n := len(p.Names)
			if n == 0 {
				n = 1
			}
			for range n {
				params = append(params, p.Type)
			}
		}
		var results []ast.Expr
		if fn.Type.Results != nil {
			for _, r := range fn.Type.Results.List {
				n := len(r.Names)
				if n == 0 {
					n = 1
				}
				for range n {
					results = append(results, r.Type)
				}
			}
		}
		if len(params) != 2 || !s30IsIdent(params[0], "int") || !s30IsIdent(params[1], "int") ||
			len(results) != 2 || !s30IsIdent(results[0], "int") || !s30IsIdent(results[1], "bool") {
			t.Errorf("resolve.IntWidth is not `func(minV, maxV int) (int, bool)`")
		}
	}
	if !found {
		t.Error("REQ-4 owed: internal/resolve declares no exported `IntWidth(minV, maxV int) (int, bool)`")
	}

	guardFuncs := s30Funcs(t, ".")
	for _, name := range []string{"intWidth", "IntWidth"} {
		if len(guardFuncs[name]) > 0 {
			t.Errorf("internal/guard still declares %s — the width rule was copied, not moved", name)
		}
	}
	for _, name := range []string{"domainSize", "IntDomain", "valueAssignments"} {
		if len(guardFuncs[name]) == 0 {
			t.Errorf("internal/guard declares no %s; the record keeps it there", name)
			continue
		}
		if !s30Reaches(guardFuncs, name, "resolve.IntWidth") {
			t.Errorf("guard::%s does not reach resolve.IntWidth — a second width rule survives", name)
		}
	}
}

// REQ-5: "the loader's emitted `guard.all eq = <cell>` must render
// byte-identically to the cells `valueAssignments` produces or the atom never
// satisfies at lint."
// DOMAIN EDGE
func TestReq5_0030_TheEmittedCellDenotesExactlyOneLintCell(t *testing.T) {
	src, err := os.ReadFile(filepath.Join("..", "table", "testdata", "ladder-step.toml"))
	if err != nil {
		t.Fatalf("reading the ladder fixture: %v", err)
	}
	m, err := table.Load(src, "ladder-step.toml")
	if err != nil {
		t.Fatalf("ladder-step.toml must load; refused: %v", err)
	}

	stepped := map[string]string{"retry": "attempt", "escalate": "tier"}
	var checked int
	for _, row := range m.Rows {
		key, ok := stepped[row.RuleID]
		if !ok {
			continue
		}
		cell := row.Suffix[len(row.Suffix)-1]
		var eq *table.Atom
		for i, a := range row.Atoms {
			if a.Key == key && a.Block == table.BlockAll && a.Operator == "eq" {
				eq = &row.Atoms[i]
			}
		}
		if eq == nil {
			t.Errorf("%s carries no guard.all eq on %s", row.Identity(), key)
			continue
		}
		if key == "attempt" {
			n, err := strconv.Atoi(cell)
			if err != nil || strconv.Itoa(n) != eq.Literal[0] {
				t.Errorf("%s: the emitted literal %q is not the canonical decimal of cell %q",
					row.Identity(), eq.Literal[0], cell)
			}
		}
		if got := guard.Denotation(m, key, *eq).Len(); got != 1 {
			t.Errorf("%s: `guard.all eq = %v` denotes %d lint cells; want exactly 1 — the literal "+
				"does not render as a cell lint enumerates", row.Identity(), eq.Literal, got)
		}
		if guard.AcceptedAssignments(m, row).Len() == 0 {
			t.Errorf("%s accepts no assignment at lint; its cell atom never satisfies", row.Identity())
		}
		checked++
	}
	if checked != 7 {
		t.Errorf("checked %d stepped rows; want 7 (retry#0..4, escalate#small, escalate#mid)", checked)
	}
}
