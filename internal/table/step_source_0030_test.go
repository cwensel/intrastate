package table_test

// RDR 0030 — the loader-side shape claims the record states as code
// structure with a named observable: the admits filter's UNDECIDED arm
// (`0030:C1`) and `expand`'s totality (`0030:C1`, IP Phase 2).
//
// These read the package SOURCE through `go/parser`, never a symbol the tree
// may not yet declare: a suite naming an unwritten symbol cannot pass the
// pre-commit `go vet ./...`.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// s30TableFiles parses this package's non-test sources.
func s30TableFiles(t *testing.T) []*ast.File {
	t.Helper()

	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("reading the package directory: %v", err)
	}
	fset := token.NewFileSet()
	var out []*ast.File
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(".", name), nil, 0)
		if err != nil {
			t.Fatalf("parsing %s: %v", name, err)
		}
		out = append(out, f)
	}
	return out
}

// s30IsSel reports whether e is the selector `pkg.name`.
func s30IsSel(e ast.Expr, pkg, name string) bool {
	sel, ok := e.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	id, ok := sel.X.(*ast.Ident)
	return ok && id.Name == pkg && sel.Sel.Name == name
}

// REQ-15: "The runtime evaluator is three-valued; an atom answering neither
// true nor false at a member does NOT admit that member." A soundness fence
// unreachable through a loaded model; asserted at the filter's own seam.
//
// The filter is unexported and the record names no spelling for it, so the
// seam is read as SOURCE: the loader must decide admission against
// `resolve.GuardTrue` (an `==`/`!=` test or a `case`), and must never collapse
// the verdict against `resolve.GuardFalse` alone — `v != GuardFalse` admits
// GuardUnevaluable, which is exactly the member C1 excludes.
// ADVERSARIAL
func TestReq15_0030_TheLoaderAdmitsACellOnlyOnGuardTrue(t *testing.T) {
	var decidesOnTrue bool
	for _, f := range s30TableFiles(t) {
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.BinaryExpr:
				if x.Op != token.EQL && x.Op != token.NEQ {
					return true
				}
				if s30IsSel(x.X, "resolve", "GuardTrue") || s30IsSel(x.Y, "resolve", "GuardTrue") {
					decidesOnTrue = true
				}
				if s30IsSel(x.X, "resolve", "GuardFalse") || s30IsSel(x.Y, "resolve", "GuardFalse") {
					t.Errorf("internal/table compares a verdict against resolve.GuardFalse alone "+
						"(%s); that two-valued collapse admits GuardUnevaluable", x.Op)
				}
			case *ast.CaseClause:
				for _, e := range x.List {
					if s30IsSel(e, "resolve", "GuardTrue") {
						decidesOnTrue = true
					}
				}
			}
			return true
		})
	}
	if !decidesOnTrue {
		t.Error("REQ-15 owed: internal/table decides no admission against resolve.GuardTrue — " +
			"the loader's admitted-cell filter, which admits a member only on a true verdict, " +
			"does not exist yet")
	}
}

// REQ-22: "`::expand` then does only what it does for `in` today: take a list
// of members per key and mint one row per combination, carrying the supplied
// literal and appending the suffix element. It therefore stays TOTAL — it
// keeps its `[]Row` return and gains no error — because every refusal this
// record mints has already fired before it is called." Observable: `expand`'s
// result list is `[]Row` alone.
// BOUNDARY
func TestReq22_0030_ExpandStaysTotalWhileMintingSteppedRows(t *testing.T) {
	// The stepped rows exist — the totality claim is about the expand that
	// mints them, so a tree where nothing steps proves nothing.
	m := s30MustLoad(t, "the admitted step", s30AdmittedRetry())
	if got := len(rowsByRuleID(m, "retry")); got != 5 {
		t.Errorf("the admitted step expands to %d rows; want 5", got)
	}

	var found bool
	for _, f := range s30TableFiles(t) {
		for _, d := range f.Decls {
			fn, ok := d.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || fn.Name.Name != "expand" {
				continue
			}
			found = true
			res := fn.Type.Results
			if res == nil || len(res.List) != 1 || len(res.List[0].Names) > 1 {
				t.Fatalf("expand's result list is not one value: %+v", res)
			}
			arr, ok := res.List[0].Type.(*ast.ArrayType)
			id, idOK := func() (*ast.Ident, bool) {
				if !ok || arr.Len != nil {
					return nil, false
				}
				i, k := arr.Elt.(*ast.Ident)
				return i, k
			}()
			if !idOK || id.Name != "Row" {
				t.Errorf("expand returns something other than []Row alone")
			}
		}
	}
	if !found {
		t.Error("no package-level expand function; the expansion loop the record names is gone")
	}
}
