package resolve_test

// RDR 0030 IP Phase 2 / `0030:C1` — the SHARED comparison: the `Evaluator`
// moves here from `internal/guard`, and the loader's admitted-cell filter,
// `guard::valueSatisfies` and `graphlint::atomAdmitsValue` reach one
// comparison in this package, each keeping its own undecided-arm policy.
//
// Every assertion reads the module SOURCE through `go/parser`: the relocated
// symbols do not exist before the implementation, and a suite naming them
// would break the build the pre-commit `go vet ./...` gates.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// s30Pkg is one parsed package: its files, function/method declarations by
// name, and type specs by name.
type s30Pkg struct {
	files []*ast.File
	funcs map[string][]*ast.FuncDecl
	types map[string]*ast.TypeSpec
}

func s30Parse(t *testing.T, dir string) s30Pkg {
	t.Helper()

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading %s: %v", dir, err)
	}
	fset := token.NewFileSet()
	p := s30Pkg{funcs: map[string][]*ast.FuncDecl{}, types: map[string]*ast.TypeSpec{}}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(dir, name), nil, 0)
		if err != nil {
			t.Fatalf("parsing %s: %v", name, err)
		}
		p.files = append(p.files, f)
		for _, d := range f.Decls {
			switch x := d.(type) {
			case *ast.FuncDecl:
				p.funcs[x.Name.Name] = append(p.funcs[x.Name.Name], x)
			case *ast.GenDecl:
				for _, s := range x.Specs {
					if ts, ok := s.(*ast.TypeSpec); ok {
						p.types[ts.Name.Name] = ts
					}
				}
			}
		}
	}
	return p
}

// s30Calls returns the calls a node makes: bare names for same-package
// calls and method selectors, `resolve.Name` for this package's selectors.
func s30Calls(n ast.Node) []string {
	var out []string
	if n == nil {
		return out
	}
	ast.Inspect(n, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		switch fn := call.Fun.(type) {
		case *ast.Ident:
			out = append(out, fn.Name)
		case *ast.SelectorExpr:
			if id, ok := fn.X.(*ast.Ident); ok && id.Name == "resolve" {
				out = append(out, "resolve."+fn.Sel.Name)
				return true
			}
			out = append(out, fn.Sel.Name)
		}
		return true
	})
	return out
}

// s30ResolveCalls returns the `resolve.Name` calls reachable from start in
// pkg, following same-package calls transitively.
func s30ResolveCalls(p s30Pkg, start string) map[string]bool {
	out := map[string]bool{}
	seen := map[string]bool{start: true}
	queue := []string{start}
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		for _, fn := range p.funcs[name] {
			for _, c := range s30Calls(fn.Body) {
				if strings.HasPrefix(c, "resolve.") {
					out[c] = true
					continue
				}
				if _, local := p.funcs[c]; local && !seen[c] {
					seen[c] = true
					queue = append(queue, c)
				}
			}
		}
	}
	return out
}

// s30AllResolveCalls returns every `resolve.Name` call anywhere in pkg.
func s30AllResolveCalls(p s30Pkg) map[string]bool {
	out := map[string]bool{}
	for _, f := range p.files {
		for _, c := range s30Calls(f) {
			if strings.HasPrefix(c, "resolve.") {
				out[c] = true
			}
		}
	}
	return out
}

func s30Ident(e ast.Expr, name string) bool {
	id, ok := e.(*ast.Ident)
	return ok && id.Name == name
}

func s30Sel(e ast.Expr, pkg, name string) bool {
	sel, ok := e.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	id, ok := sel.X.(*ast.Ident)
	return ok && id.Name == pkg && sel.Sel.Name == name
}

// s30Flatten expands a field list into one type expression per value.
func s30Flatten(fl *ast.FieldList) []ast.Expr {
	var out []ast.Expr
	if fl == nil {
		return out
	}
	for _, f := range fl.List {
		n := len(f.Names)
		if n == 0 {
			n = 1
		}
		for range n {
			out = append(out, f.Type)
		}
	}
	return out
}

// s30RecvName returns a method's receiver type name.
func s30RecvName(fn *ast.FuncDecl) string {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return ""
	}
	e := fn.Recv.List[0].Type
	if st, ok := e.(*ast.StarExpr); ok {
		e = st.X
	}
	if id, ok := e.(*ast.Ident); ok {
		return id.Name
	}
	return ""
}

// s30ConcreteEvaluators names this package's non-interface types carrying a
// method `Evaluate(GuardAtom, string) GuardResult`.
func s30ConcreteEvaluators(p s30Pkg) []string {
	var out []string
	for _, fn := range p.funcs["Evaluate"] {
		recv := s30RecvName(fn)
		ts, ok := p.types[recv]
		if recv == "" || !ok {
			continue
		}
		if _, iface := ts.Type.(*ast.InterfaceType); iface {
			continue
		}
		params, results := s30Flatten(fn.Type.Params), s30Flatten(fn.Type.Results)
		if len(params) == 2 && s30Ident(params[0], "GuardAtom") && s30Ident(params[1], "string") &&
			len(results) == 1 && s30Ident(results[0], "GuardResult") {
			out = append(out, recv)
		}
	}
	slices.Sort(out)
	return out
}

// s30Exported names this package's exported package-level functions.
func s30Exported(p s30Pkg) map[string]bool {
	out := map[string]bool{}
	for name, fns := range p.funcs {
		for _, fn := range fns {
			if fn.Recv == nil && ast.IsExported(name) {
				out["resolve."+name] = true
			}
		}
	}
	return out
}

// s30Shared returns the exported resolve functions all three comparison call
// sites reach: the loader (anywhere in internal/table), guard::valueSatisfies
// and graphlint::atomAdmitsValue.
func s30Shared(t *testing.T) []string {
	t.Helper()

	here := s30Parse(t, ".")
	exported := s30Exported(here)
	tablePkg := s30Parse(t, filepath.Join("..", "table"))
	guardPkg := s30Parse(t, filepath.Join("..", "guard"))
	lintPkg := s30Parse(t, filepath.Join("..", "graphlint"))

	for _, c := range []struct {
		p    s30Pkg
		name string
	}{{guardPkg, "valueSatisfies"}, {lintPkg, "atomAdmitsValue"}} {
		if len(c.p.funcs[c.name]) == 0 {
			t.Fatalf("no %s declared; the record names it as a comparison call site", c.name)
		}
	}

	loader := s30AllResolveCalls(tablePkg)
	guardSite := s30ResolveCalls(guardPkg, "valueSatisfies")
	lintSite := s30ResolveCalls(lintPkg, "atomAdmitsValue")
	var out []string
	for name := range exported {
		if loader[name] && guardSite[name] && lintSite[name] {
			out = append(out, name)
		}
	}
	slices.Sort(out)
	return out
}

// REQ-14: "the shim takes a `resolve.GuardAtom` and the held value as a
// STRING — the cell, never the declared kind — and returns
// `resolve.GuardResult`, whose members are `GuardTrue`, `GuardFalse` and
// `GuardUnevaluable`."
// HAPPY PATH
func TestReq14_0030_TheShimIsAConcreteEvaluatorHomedInResolve(t *testing.T) {
	here := s30Parse(t, ".")
	evs := s30ConcreteEvaluators(here)
	if len(evs) == 0 {
		t.Fatal("REQ-14 owed: internal/resolve declares no concrete type with " +
			"`Evaluate(GuardAtom, string) GuardResult` — the relocated evaluator the shim is")
	}
	// It is constructed, never zero-valued: some exported constructor
	// returns it.
	var constructed bool
	for name, fns := range here.funcs {
		for _, fn := range fns {
			if fn.Recv != nil || !ast.IsExported(name) {
				continue
			}
			for _, r := range s30Flatten(fn.Type.Results) {
				for _, ev := range evs {
					if s30Ident(r, ev) {
						constructed = true
					}
				}
			}
		}
	}
	if !constructed {
		t.Errorf("no exported internal/resolve function returns %v; the loader cannot construct it", evs)
	}
	for _, member := range []string{"GuardTrue", "GuardFalse", "GuardUnevaluable"} {
		var declared bool
		for _, f := range here.files {
			ast.Inspect(f, func(n ast.Node) bool {
				if vs, ok := n.(*ast.ValueSpec); ok {
					for _, id := range vs.Names {
						if id.Name == member {
							declared = true
						}
					}
				}
				return true
			})
		}
		if !declared {
			t.Errorf("GuardResult member %s is not declared in internal/resolve", member)
		}
	}
}

// REQ-40: "Move `internal/guard/grammar.go`'s `Evaluator` to
// `internal/resolve`, then extract the two-valued core of
// `internal/guard/product.go::valueSatisfies` — render the atom's literal,
// call the evaluator, return its three-valued verdict — to
// `internal/resolve` alongside it, carrying the canonicalizing set
// renderer." With "The loader's admits filter then calls the same
// comparison the runtime does … rather than minting a fourth copy."
// Observable: the loader's admits filter, `guard::valueSatisfies` and
// `graphlint::atomAdmitsValue` reach one comparison in `internal/resolve`;
// the loader holds no reimplementation.
// BOUNDARY
func TestReq40_0030_ThreeCallSitesReachOneComparisonInResolve(t *testing.T) {
	here := s30Parse(t, ".")
	evs := s30ConcreteEvaluators(here)
	if len(evs) == 0 {
		t.Error("REQ-40 owed: the Evaluator has not moved — internal/resolve declares no concrete evaluator")
	}

	// guard keeps `Evaluator` as a re-export of the moved type.
	guardPkg := s30Parse(t, filepath.Join("..", "guard"))
	ts, ok := guardPkg.types["Evaluator"]
	switch {
	case !ok:
		t.Error("internal/guard no longer names Evaluator; 0012's pins keep it as a re-export")
	case !ts.Assign.IsValid():
		t.Error("internal/guard still DECLARES Evaluator; it moves to internal/resolve and guard aliases it")
	default:
		sel, isSel := ts.Type.(*ast.SelectorExpr)
		if !isSel || !s30Sel(sel, "resolve", sel.Sel.Name) || !slices.Contains(evs, sel.Sel.Name) {
			t.Errorf("guard.Evaluator aliases something other than the relocated resolve evaluator %v", evs)
		}
	}

	if shared := s30Shared(t); len(shared) == 0 {
		t.Error("REQ-40 owed: no exported internal/resolve function is reached by all three " +
			"comparison call sites (the loader's admits filter, guard::valueSatisfies, " +
			"graphlint::atomAdmitsValue)")
	}

	// No fourth copy: the loader declares no evaluator of its own.
	tablePkg := s30Parse(t, filepath.Join("..", "table"))
	if len(tablePkg.funcs["Evaluate"]) > 0 {
		t.Error("internal/table declares its own Evaluate — a fourth copy of the comparison")
	}
}

// REQ-41: "Each caller's undecided-arm disposition stays its own — the shim
// is the comparison, never the policy — so the repoint is
// behaviour-preserving in all three." With "`guard` passes the verdict
// through, `graphlint` admits, the loader excludes per C1"
// DOMAIN EDGE
func TestReq41_0030_EachCallerKeepsItsOwnUndecidedArmPolicy(t *testing.T) {
	// The repoint happened, or the policy claims below are about the old
	// copies.
	if shared := s30Shared(t); len(shared) == 0 {
		t.Error("REQ-41 owed: the three call sites are not repointed through one resolve comparison")
	}

	// guard passes the three-valued verdict through.
	guardPkg := s30Parse(t, filepath.Join("..", "guard"))
	for _, fn := range guardPkg.funcs["valueSatisfies"] {
		res := s30Flatten(fn.Type.Results)
		if len(res) != 1 || !s30Sel(res[0], "resolve", "GuardResult") {
			t.Error("guard::valueSatisfies no longer returns resolve.GuardResult — it collapsed the verdict")
		}
	}

	// graphlint admits the undecided arm: a bool collapse against GuardFalse.
	lintPkg := s30Parse(t, filepath.Join("..", "graphlint"))
	for _, fn := range lintPkg.funcs["atomAdmitsValue"] {
		res := s30Flatten(fn.Type.Results)
		if len(res) != 1 || !s30Ident(res[0], "bool") {
			t.Error("graphlint::atomAdmitsValue no longer returns bool")
		}
		var admits bool
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			if b, ok := n.(*ast.BinaryExpr); ok && b.Op == token.NEQ &&
				(s30Sel(b.X, "resolve", "GuardFalse") || s30Sel(b.Y, "resolve", "GuardFalse")) {
				admits = true
			}
			return true
		})
		if !admits {
			t.Error("graphlint::atomAdmitsValue no longer admits the undecided arm (`!= resolve.GuardFalse`)")
		}
	}

	// The loader excludes it: admission on GuardTrue, never a GuardFalse-only
	// comparison.
	tablePkg := s30Parse(t, filepath.Join("..", "table"))
	var onTrue bool
	for _, f := range tablePkg.files {
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.BinaryExpr:
				if s30Sel(x.X, "resolve", "GuardTrue") || s30Sel(x.Y, "resolve", "GuardTrue") {
					onTrue = true
				}
				if s30Sel(x.X, "resolve", "GuardFalse") || s30Sel(x.Y, "resolve", "GuardFalse") {
					t.Error("internal/table compares a verdict against resolve.GuardFalse alone; " +
						"the loader's policy EXCLUDES the undecided arm")
				}
			case *ast.CaseClause:
				for _, e := range x.List {
					if s30Sel(e, "resolve", "GuardTrue") {
						onTrue = true
					}
				}
			}
			return true
		})
	}
	if !onTrue {
		t.Error("REQ-41 owed: internal/table decides no admission on resolve.GuardTrue")
	}
}
