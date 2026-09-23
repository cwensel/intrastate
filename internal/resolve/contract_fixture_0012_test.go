package resolve_test

// RDR 0012 — the kernel half of the declared-kind carrier: the published
// conformance fixture (`0012:C3`) and the resolution path's freedom from
// declarations (`0012:C1`).
//
// `resolve.ContractKinds` and the constructor-taking
// `TestGuardEvaluatorContract` do not exist on the tree these tests were
// written against, and naming either as a Go identifier would break the
// build the pre-commit `go vet ./...` gates. Every assertion here therefore
// reads the package SOURCE through `go/parser`: an absent declaration is a
// runtime failure naming the REQ it owes. The behaviour of the fixture —
// fresh copies, the five-kind meta-check, per-kind discrimination — is
// driven from `internal/guard`'s probe suite, which calls it directly.

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// r12KernelFiles parses the kernel's non-test sources.
func r12KernelFiles(t *testing.T) ([]*ast.File, *token.FileSet) {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Fatalf("parsing %s: %v", dir, err)
	}
	var files []*ast.File
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			files = append(files, f)
		}
	}
	return files, fset
}

func r12Print(t *testing.T, fset *token.FileSet, node any) string {
	t.Helper()

	var b strings.Builder
	if err := printer.Fprint(&b, fset, node); err != nil {
		t.Fatalf("rendering: %v", err)
	}
	return b.String()
}

func r12Types(t *testing.T, fset *token.FileSet, list *ast.FieldList) []string {
	t.Helper()

	var out []string
	if list == nil {
		return out
	}
	for _, f := range list.List {
		for range max(len(f.Names), 1) {
			out = append(out, r12Print(t, fset, f.Type))
		}
	}
	return out
}

// r12KernelFunc returns the package-level function `name` from the kernel's
// non-test sources.
func r12KernelFunc(t *testing.T, name string) (*ast.FuncDecl, *token.FileSet) {
	t.Helper()

	files, fset := r12KernelFiles(t)
	for _, f := range files {
		for _, decl := range f.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv == nil && fn.Name.Name == name {
				return fn, fset
			}
		}
	}
	return nil, fset
}

// r12IsNewSeamType reports whether e is `func(map[string]string) <result>`.
func r12IsNewSeamType(t *testing.T, fset *token.FileSet, e ast.Expr, result string) bool {
	t.Helper()

	ft, ok := e.(*ast.FuncType)
	if !ok {
		return false
	}
	return slices.Equal(r12Types(t, fset, ft.Params), []string{"map[string]string"}) &&
		slices.Equal(r12Types(t, fset, ft.Results), []string{result})
}

// REQ-27: "The suite publishes the declaration fixture its cases assume as a
// kernel-owned `map[string]string` (key → kind token; no `internal/table`
// type crosses — A5). It is published as a FUNCTION returning a fresh map:
// func ContractKinds() map[string]string not a package-level `var`."
// BOUNDARY
func TestReq27_0012_ContractKindsIsAFunctionNotAPackageVar(t *testing.T) {
	fn, fset := r12KernelFunc(t, "ContractKinds")
	if fn == nil {
		t.Fatal("REQ-27 owed: the kernel's non-test sources declare no ContractKinds function")
	}
	if got := r12Types(t, fset, fn.Type.Params); len(got) != 0 {
		t.Errorf("ContractKinds takes %v; want no parameters", got)
	}
	if got := r12Types(t, fset, fn.Type.Results); !slices.Equal(got, []string{"map[string]string"}) {
		t.Errorf("ContractKinds returns %v; want map[string]string (no internal/table type crosses)", got)
	}

	// Not a package-level var in disguise: no return hands back a
	// package-scope identifier, which every caller would then share.
	files, _ := r12KernelFiles(t)
	pkgVars := map[string]bool{}
	for _, f := range files {
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.VAR {
				continue
			}
			for _, spec := range gd.Specs {
				for _, name := range spec.(*ast.ValueSpec).Names {
					pkgVars[name.Name] = true
				}
			}
		}
	}
	if pkgVars["ContractKinds"] {
		t.Error("ContractKinds is also a package-level var; it is published as a FUNCTION")
	}
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		ret, ok := n.(*ast.ReturnStmt)
		if !ok {
			return true
		}
		for _, r := range ret.Results {
			if id, ok := r.(*ast.Ident); ok && pkgVars[id.Name] {
				t.Errorf("ContractKinds returns the package-level var %s; each caller gets a "+
					"FRESH map", id.Name)
			}
		}
		return true
	})

	// The fixture must not import a table type into the kernel.
	for _, f := range files {
		for _, imp := range f.Imports {
			if p, _ := strconv.Unquote(imp.Path.Value); strings.HasSuffix(p, "/internal/table") {
				t.Errorf("%s imports internal/table; no table type crosses into the kernel",
					fset.Position(f.Pos()).Filename)
			}
		}
	}
}

// REQ-28: "Delivery is STRUCTURAL, not a documented obligation — the suite
// takes a constructor and builds the seam itself: func
// TestGuardEvaluatorContract(t *testing.T, newSeam func(kinds
// map[string]string) GuardEvaluator)"
// BOUNDARY
func TestReq28_0012_TheSuiteTakesASeamConstructor(t *testing.T) {
	fn, fset := r12KernelFunc(t, "TestGuardEvaluatorContract")
	if fn == nil {
		t.Fatal("TestGuardEvaluatorContract is no longer declared in the kernel's non-test sources")
	}
	params := fn.Type.Params.List
	var types []ast.Expr
	for _, p := range params {
		for range max(len(p.Names), 1) {
			types = append(types, p.Type)
		}
	}
	if len(types) != 2 {
		t.Fatalf("TestGuardEvaluatorContract takes %d parameters; want (t *testing.T, newSeam "+
			"func(kinds map[string]string) GuardEvaluator)", len(types))
	}
	if got := r12Print(t, fset, types[0]); got != "*testing.T" {
		t.Errorf("first parameter is %s; want *testing.T", got)
	}
	if !r12IsNewSeamType(t, fset, types[1], "GuardEvaluator") {
		t.Errorf("REQ-28 owed: second parameter is %s; want a seam CONSTRUCTOR "+
			"func(kinds map[string]string) GuardEvaluator", r12Print(t, fset, types[1]))
	}
}

// r12Reachable returns the names of the kernel's functions and methods
// reachable from `root` by identifier reference. Method calls are resolved
// by name over-approximately, which only widens the set this test forbids.
func r12Reachable(t *testing.T, root string) map[string]bool {
	t.Helper()

	files, _ := r12KernelFiles(t)
	bodies := map[string][]*ast.BlockStmt{}
	for _, f := range files {
		for _, decl := range f.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Body != nil {
				bodies[fn.Name.Name] = append(bodies[fn.Name.Name], fn.Body)
			}
		}
	}
	seen := map[string]bool{root: true}
	queue := []string{root}
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		for _, body := range bodies[name] {
			ast.Inspect(body, func(n ast.Node) bool {
				var ref string
				switch x := n.(type) {
				case *ast.Ident:
					ref = x.Name
				case *ast.SelectorExpr:
					ref = x.Sel.Name
				}
				if _, isFunc := bodies[ref]; isFunc && !seen[ref] {
					seen[ref] = true
					queue = append(queue, ref)
				}
				return true
			})
		}
	}
	return seen
}

// REQ-9: "the kernel's RESOLUTION PATH — `Resolve` and its callees —
// continues to read no declarations. The rule governs that path, not
// `internal/resolve` as a namespace"
// BOUNDARY
func TestReq9_0012_TheResolutionPathReadsNoDeclarations(t *testing.T) {
	// The namespace now carries a declaration fixture — which is exactly why
	// the rule is scoped to the PATH. Its absence means the clause's subject
	// has not landed.
	if fn, _ := r12KernelFunc(t, "ContractKinds"); fn == nil {
		t.Fatal("REQ-9 owed: internal/resolve declares no ContractKinds; the clause scopes the " +
			"no-declarations rule to Resolve's path precisely because the namespace carries one")
	}
	reach := r12Reachable(t, "Resolve")
	if !reach["gate"] {
		t.Fatalf("the reachability walk did not find `gate` from Resolve (%d functions); "+
			"the walk is broken, not the kernel", len(reach))
	}
	for _, forbidden := range []string{"ContractKinds", "TestGuardEvaluatorContract"} {
		if reach[forbidden] {
			t.Errorf("Resolve reaches %s; the resolution path reads no declarations", forbidden)
		}
	}

	files, fset := r12KernelFiles(t)
	for _, f := range files {
		for _, imp := range f.Imports {
			p, _ := strconv.Unquote(imp.Path.Value)
			if strings.HasSuffix(p, "/internal/table") || strings.HasSuffix(p, "/internal/guard") {
				t.Errorf("%s imports %s; the kernel reads no declarations",
					filepath.Base(fset.Position(f.Pos()).Filename), p)
			}
		}
	}
}

// r12TestFile parses one of this package's test files.
func r12TestFile(t *testing.T, name string) (*ast.File, *token.FileSet) {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, name, nil, 0)
	if err != nil {
		t.Fatalf("parsing %s: %v", name, err)
	}
	return file, fset
}

// REQ-30: "Both in-tree implementers are bound, not only `guard.Evaluator`"
// … "Concretely it gains a `kinds map[string]string` field and a matching
// constructor to satisfy the `newSeam` parameter above; `recordingSeam`
// keeps delegating to its `inner` and needs no kinds of its own." ("it" =
// `internal/resolve/guard_mvv_test.go::conformingContractSeam`)
// BOUNDARY
func TestReq30_0012_TheConformingTestSeamIsBoundThroughAConstructor(t *testing.T) {
	file, fset := r12TestFile(t, "guard_mvv_test.go")

	structs := map[string]*ast.StructType{}
	ast.Inspect(file, func(n ast.Node) bool {
		if ts, ok := n.(*ast.TypeSpec); ok {
			if st, ok := ts.Type.(*ast.StructType); ok {
				structs[ts.Name.Name] = st
			}
		}
		return true
	})

	fieldTypes := func(name string) map[string]string {
		out := map[string]string{}
		st := structs[name]
		if st == nil {
			return out
		}
		for _, f := range st.Fields.List {
			for _, n := range f.Names {
				out[n.Name] = r12Print(t, fset, f.Type)
			}
		}
		return out
	}

	conforming := fieldTypes("conformingContractSeam")
	if _, ok := structs["conformingContractSeam"]; !ok {
		t.Fatal("guard_mvv_test.go no longer declares conformingContractSeam")
	}
	hasKinds := false
	for _, typ := range conforming {
		if typ == "map[string]string" {
			hasKinds = true
		}
	}
	if !hasKinds {
		t.Errorf("REQ-30 owed: conformingContractSeam carries fields %v; it gains a "+
			"`kinds map[string]string` field", conforming)
	}

	// The matching constructor: a function taking the mapping and building
	// a conformingContractSeam.
	ctor := false
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil {
			continue
		}
		if !slices.Equal(r12Types(t, fset, fn.Type.Params), []string{"map[string]string"}) {
			continue
		}
		if strings.Contains(r12Print(t, fset, fn.Body), "conformingContractSeam{") {
			ctor = true
		}
	}
	if !ctor {
		t.Error("REQ-30 owed: guard_mvv_test.go declares no constructor " +
			"func(map[string]string) … that builds a conformingContractSeam")
	}

	recording := fieldTypes("recordingSeam")
	if _, ok := recording["inner"]; !ok {
		t.Errorf("recordingSeam fields %v; it keeps delegating to its `inner`", recording)
	}
	for name, typ := range recording {
		if typ == "map[string]string" {
			t.Errorf("recordingSeam gains a kinds field %s; it needs no kinds of its own", name)
		}
	}
}

// REQ-32: "its importable-cross-RDR-surface guarantee is preserved — same
// name, still declared in the kernel's non-test sources (the
// `kernelDeclaresFunc` half at :227 is untouched) — and only the parameter
// shape changes" … "the `var fn` line is re-typed to the new signature
// rather than deleted, so the pin keeps pinning."
// BOUNDARY
func TestReq32_0012_TheImportablePinIsReTypedNotDeleted(t *testing.T) {
	file, fset := r12TestFile(t, "guard_mvv_test.go")

	var body *ast.BlockStmt
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok &&
			fn.Name.Name == "TestReq71_GuardEvaluatorContractIsAnImportableCrossRDRSurface" {
			body = fn.Body
		}
	}
	if body == nil {
		t.Fatal("TestReq71_GuardEvaluatorContractIsAnImportableCrossRDRSurface is gone; the pin " +
			"is re-typed, not deleted")
	}

	pinned, declares := false, false
	ast.Inspect(body, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ValueSpec:
			ft, ok := x.Type.(*ast.FuncType)
			if !ok || len(x.Values) != 1 {
				return true
			}
			if r12Print(t, fset, x.Values[0]) != "resolve.TestGuardEvaluatorContract" {
				return true
			}
			var params []ast.Expr
			for _, p := range ft.Params.List {
				for range max(len(p.Names), 1) {
					params = append(params, p.Type)
				}
			}
			if len(params) == 2 && r12Print(t, fset, params[0]) == "*testing.T" &&
				r12IsNewSeamType(t, fset, params[1], "resolve.GuardEvaluator") {
				pinned = true
			}
		case *ast.CallExpr:
			if id, ok := x.Fun.(*ast.Ident); ok && id.Name == "kernelDeclaresFunc" {
				declares = true
			}
		}
		return true
	})
	if !pinned {
		t.Error("REQ-32 owed: the `var fn` pin is not re-typed to " +
			"func(*testing.T, func(map[string]string) resolve.GuardEvaluator)")
	}
	if !declares {
		t.Error("the kernelDeclaresFunc half of the pin is gone; it stays untouched")
	}
}
