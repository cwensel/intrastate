package guard_test

// RDR 0012 — the declared-kind carrier at the guard value seam.
//
// Two kinds of test live here.
//
//  1. Tests over surfaces that exist today — the zero-value `Evaluator`, its
//     reflected field set, and the module's non-test SOURCE read through
//     `go/parser`. These name nothing the tree lacks.
//  2. Drivers for the behavioural probes in `rdr0012_probe_test.go`. Those
//     call `guard.NewEvaluator`, `guard.DeclaredKinds` and
//     `resolve.ContractKinds` directly, so they sit behind the
//     `rdr0012probe` build tag and run in a `go test -tags rdr0012probe`
//     child. A probe that does not compile is reported here as the REQ it
//     owes — a runtime assertion, never a build break, so the red suite
//     survives the pre-commit `go vet ./...`.

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/cwensel/intrastate/internal/guard"
	"github.com/cwensel/intrastate/internal/resolve"
)

// --- the probe child --------------------------------------------------------

const r12ProbeChildEnv = "RDR0012_PROBE_CHILD"

type r12ProbeRun struct {
	// results is each probe test's terminal action: pass, fail, or skip.
	results map[string]string
	// output is each probe test's captured output.
	output map[string]*strings.Builder
	// build is everything the child printed that was not a test event —
	// the compile errors when the probe file does not build.
	build string
}

var (
	r12ProbeOnce sync.Once
	r12Probe     r12ProbeRun
)

// r12RunProbes runs every `TestRDR0012Probe_*` in a tagged child once per
// test binary.
func r12RunProbes(t *testing.T) *r12ProbeRun {
	t.Helper()

	if os.Getenv(r12ProbeChildEnv) != "" {
		t.Skip("inside the probe child; drivers do not recurse")
	}
	r12ProbeOnce.Do(func() {
		r12Probe = r12ProbeRun{results: map[string]string{}, output: map[string]*strings.Builder{}}

		dir, err := os.Getwd()
		if err != nil {
			r12Probe.build = "getwd: " + err.Error()
			return
		}
		cmd := exec.Command("go", "test", "-tags", "rdr0012probe", "-json",
			"-run", "^TestRDR0012Probe_", ".")
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), r12ProbeChildEnv+"=1")
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		_ = cmd.Run() // a failing probe is data, not an error

		var build strings.Builder
		build.WriteString(stderr.String())
		sc := bufio.NewScanner(&stdout)
		sc.Buffer(make([]byte, 1<<20), 1<<24)
		for sc.Scan() {
			line := sc.Bytes()
			var ev struct {
				Action string
				Test   string
				Output string
			}
			if json.Unmarshal(line, &ev) != nil {
				build.Write(line)
				build.WriteByte('\n')
				continue
			}
			if ev.Test == "" {
				if ev.Action == "build-output" || ev.Action == "output" {
					build.WriteString(ev.Output)
				}
				continue
			}
			switch ev.Action {
			case "pass", "fail", "skip":
				r12Probe.results[ev.Test] = ev.Action
			case "output":
				b := r12Probe.output[ev.Test]
				if b == nil {
					b = &strings.Builder{}
					r12Probe.output[ev.Test] = b
				}
				b.WriteString(ev.Output)
			}
		}
		r12Probe.build = build.String()
	})
	return &r12Probe
}

// r12Owed renders the child's build output as the reason a REQ is owed.
func (p *r12ProbeRun) owed(req, probe string) string {
	lines := strings.Split(strings.TrimSpace(p.build), "\n")
	if len(lines) > 25 {
		lines = append(lines[:25], "…")
	}
	return req + " owed: probe " + probe + " did not run — the probe suite does not " +
		"compile against this tree (the RDR 0012 surface is absent):\n" + strings.Join(lines, "\n")
}

// r12RequireProbe asserts the named probe ran and PASSED.
func r12RequireProbe(t *testing.T, req, probe string) {
	t.Helper()

	p := r12RunProbes(t)
	switch p.results[probe] {
	case "pass":
		return
	case "":
		t.Error(p.owed(req, probe))
	default:
		out := ""
		if b := p.output[probe]; b != nil {
			out = b.String()
		}
		t.Errorf("%s: probe %s %sed:\n%s", req, probe, p.results[probe], out)
	}
}

// r12RequireProbeFails asserts the named MUTANT probe ran and FAILED: the
// conformance suite caught the deliberately wrong seam.
func r12RequireProbeFails(t *testing.T, req, probe, why string) {
	t.Helper()

	p := r12RunProbes(t)
	switch p.results[probe] {
	case "fail":
		return
	case "":
		t.Error(p.owed(req, probe))
	default:
		t.Errorf("%s: the conformance suite PASSED mutant %s; %s", req, probe, why)
	}
}

// --- source helpers ---------------------------------------------------------

// r12ModuleRoot walks up from the package directory to go.mod.
func r12ModuleRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("no go.mod above the package directory")
		}
		dir = parent
	}
}

// r12PackageFuncs parses the non-test source of the package in dir and
// returns its package-level functions by name.
func r12PackageFuncs(t *testing.T, dir string) (map[string]*ast.FuncDecl, *token.FileSet) {
	t.Helper()

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, parser.ParseComments)
	if err != nil {
		t.Fatalf("parsing %s: %v", dir, err)
	}
	out := map[string]*ast.FuncDecl{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if ok && fn.Recv == nil {
					out[fn.Name.Name] = fn
				}
			}
		}
	}
	return out, fset
}

// r12Render prints an AST node back to source text.
func r12Render(t *testing.T, fset *token.FileSet, node any) string {
	t.Helper()

	var b strings.Builder
	if err := printer.Fprint(&b, fset, node); err != nil {
		t.Fatalf("rendering: %v", err)
	}
	return b.String()
}

// r12FieldTypes renders a field list's types, one per declared name.
func r12FieldTypes(t *testing.T, fset *token.FileSet, list *ast.FieldList) []string {
	t.Helper()

	var out []string
	if list == nil {
		return out
	}
	for _, f := range list.List {
		n := max(len(f.Names), 1)
		for range n {
			out = append(out, r12Render(t, fset, f.Type))
		}
	}
	return out
}

func r12GuardDir(t *testing.T) string {
	return filepath.Join(r12ModuleRoot(t), "internal", "guard")
}

// --- C1: the carrier --------------------------------------------------------

// REQ-1: "CARRIER. The declared kind travels to the value seam by
// CONSTRUCTION, not by signature, atom field, or upstream refusal. func
// NewEvaluator(kinds map[string]string) Evaluator"
// HAPPY PATH
func TestReq1_0012_NewEvaluatorIsDeclaredWithTheCarrierSignature(t *testing.T) {
	funcs, fset := r12PackageFuncs(t, r12GuardDir(t))
	fn, ok := funcs["NewEvaluator"]
	if !ok {
		t.Fatal("REQ-1 owed: internal/guard declares no package-level NewEvaluator")
	}
	params := r12FieldTypes(t, fset, fn.Type.Params)
	results := r12FieldTypes(t, fset, fn.Type.Results)
	if !slices.Equal(params, []string{"map[string]string"}) {
		t.Errorf("NewEvaluator parameters = %v; want (kinds map[string]string)", params)
	}
	if !slices.Equal(results, []string{"Evaluator"}) {
		t.Errorf("NewEvaluator results = %v; want Evaluator", results)
	}
	r12RequireProbe(t, "REQ-1", "TestRDR0012Probe_Req1_NewEvaluatorConstructsTheConcreteSeam")
}

// REQ-10: "`NewEvaluator`'s return type is NOT widened to the interface to
// avoid the wrapper: C1 fixes the concrete value shape deliberately"
// BOUNDARY
func TestReq10_0012_NewEvaluatorReturnsTheConcreteEvaluator(t *testing.T) {
	funcs, fset := r12PackageFuncs(t, r12GuardDir(t))
	fn, ok := funcs["NewEvaluator"]
	if !ok {
		t.Fatal("REQ-10 owed: internal/guard declares no NewEvaluator")
	}
	results := r12FieldTypes(t, fset, fn.Type.Results)
	for _, r := range results {
		if strings.Contains(r, "GuardEvaluator") {
			t.Errorf("NewEvaluator returns %s; the return is the concrete guard.Evaluator, "+
				"never widened to the interface", r)
		}
	}
	if !slices.Equal(results, []string{"Evaluator"}) {
		t.Errorf("NewEvaluator results = %v; want the concrete Evaluator", results)
	}
}

// REQ-2: "`kinds` maps tag key → declared kind token (the five-kind
// vocabulary `enum | bool | int | set | scalar`, RDR 0003's spelling). Its
// producer is EXPORTED and singular: func DeclaredKinds(m *table.Model)
// map[string]string" … "`DeclaredKinds` is homed in `internal/guard`"
// HAPPY PATH
func TestReq2_0012_DeclaredKindsIsHomedInGuardAndMapsKeyToKind(t *testing.T) {
	funcs, fset := r12PackageFuncs(t, r12GuardDir(t))
	fn, ok := funcs["DeclaredKinds"]
	if !ok {
		t.Fatal("REQ-2 owed: internal/guard declares no package-level DeclaredKinds")
	}
	if got := r12FieldTypes(t, fset, fn.Type.Params); !slices.Equal(got, []string{"*table.Model"}) {
		t.Errorf("DeclaredKinds parameters = %v; want (m *table.Model)", got)
	}
	if got := r12FieldTypes(t, fset, fn.Type.Results); !slices.Equal(got, []string{"map[string]string"}) {
		t.Errorf("DeclaredKinds results = %v; want map[string]string", got)
	}

	// Singular: no second producer of the mapping elsewhere in the module.
	root := r12ModuleRoot(t)
	for _, dir := range []string{"internal/cli", "internal/resolve", "internal/table", "internal/graphlint"} {
		others, _ := r12PackageFuncs(t, filepath.Join(root, dir))
		if _, dup := others["DeclaredKinds"]; dup {
			t.Errorf("%s also declares DeclaredKinds; the producer is singular", dir)
		}
	}
	r12RequireProbe(t, "REQ-2", "TestRDR0012Probe_Req2_DeclaredKindsMapsKeyToKindToken")
}

// REQ-3: "It is TOTAL over `m.Tags`: every declared key gets an entry, and a
// key whose `Kind` is the empty string is OMITTED rather than mapped to
// `\"\"`"
// DOMAIN EDGE
func TestReq3_0012_DeclaredKindsIsTotalAndOmitsEmptyKind(t *testing.T) {
	r12RequireProbe(t, "REQ-3", "TestRDR0012Probe_Req3_DeclaredKindsIsTotalAndOmitsEmptyKind")
}

// REQ-4: "Handed a nil model it returns an empty non-nil map rather than
// panicking"
// INPUT EDGE
func TestReq4_0012_DeclaredKindsOverANilModelIsEmptyNonNil(t *testing.T) {
	r12RequireProbe(t, "REQ-4", "TestRDR0012Probe_Req4_DeclaredKindsOverNilModelIsEmptyNonNil")
}

// r12ForbiddenState reports whether a field type is, or reaches, view-typed
// or runtime-valued state from the kernel.
func r12ForbiddenState(typ reflect.Type, seen map[reflect.Type]bool) (string, bool) {
	if seen[typ] {
		return "", false
	}
	seen[typ] = true
	if typ.PkgPath() == reflect.TypeOf(resolve.TagSet{}).PkgPath() {
		return typ.String(), true
	}
	switch typ.Kind() {
	case reflect.Interface, reflect.Func, reflect.Chan, reflect.UnsafePointer:
		return typ.String(), true
	case reflect.Pointer, reflect.Slice, reflect.Array:
		return r12ForbiddenState(typ.Elem(), seen)
	case reflect.Map:
		if name, bad := r12ForbiddenState(typ.Key(), seen); bad {
			return name, true
		}
		return r12ForbiddenState(typ.Elem(), seen)
	case reflect.Struct:
		for i := range typ.NumField() {
			if name, bad := r12ForbiddenState(typ.Field(i).Type, seen); bad {
				return name, true
			}
		}
	}
	return "", false
}

// REQ-5: "The evaluator holds this mapping and NOTHING else: no
// `resolve.TagSet`, no runtime tag value, no view-typed state"
// REQ-47: "The zero-field reflection assertion in `guard_evaluator_0003_test.go`
// narrowed to \"no view-typed or runtime-valued state\". **Expected**: passes
// with the declaration mapping present, fails if a `resolve.TagSet` or
// runtime value is added."
// BOUNDARY
func TestReq47_0012_EvaluatorHoldsTheDeclarationMappingAndNothingElse(t *testing.T) {
	typ := reflect.TypeOf(guard.Evaluator{})
	mapType := reflect.TypeOf(map[string]string(nil))

	mappings := 0
	for i := range typ.NumField() {
		f := typ.Field(i)
		if f.Type == mapType {
			mappings++
			continue
		}
		if name, bad := r12ForbiddenState(f.Type, map[reflect.Type]bool{}); bad {
			t.Errorf("guard.Evaluator field %s carries %s; the evaluator holds the "+
				"declaration mapping and NO view-typed or runtime-valued state", f.Name, name)
			continue
		}
		t.Errorf("guard.Evaluator field %s (%s) is not the declaration mapping; the "+
			"evaluator holds this mapping and NOTHING else", f.Name, f.Type)
	}
	if mappings != 1 {
		t.Errorf("guard.Evaluator holds %d map[string]string declaration mappings; want "+
			"exactly one — the kind travels to the seam by CONSTRUCTION", mappings)
	}
}

// REQ-7: "a nil-mapping evaluator answers GuardUnevaluable for every
// `eq`/`in` atom (the no-declaration arm) … never silently reverting to
// raw-string comparison."
// REQ-49: "Construct a bare `guard.Evaluator{}` (and the `var` / `new` zero
// values) and evaluate an `eq` atom whose held value equals its literal,
// plus an `in` atom whose held value is a member. **Expected**:
// `GuardUnevaluable` for BOTH — never `GuardTrue`."
// ADVERSARIAL
func TestReq49_0012_AZeroValueEvaluatorNeverRevertsToRawComparison(t *testing.T) {
	var declared guard.Evaluator
	zeros := map[string]guard.Evaluator{
		"composite literal": {},
		"var":               declared,
		"new":               *new(guard.Evaluator),
	}
	atoms := []struct {
		atom  resolve.GuardAtom
		value string
	}{
		{resolve.GuardAtom{Key: "k", Operator: "eq", Literal: "x", Block: resolve.BlockAll}, "x"},
		{resolve.GuardAtom{Key: "k", Operator: "in", Literal: `["x","y"]`, Block: resolve.BlockAll}, "x"},
		{resolve.GuardAtom{Key: "k", Operator: "eq", Literal: "7", Block: resolve.BlockAll}, "7"},
	}
	for form, ev := range zeros {
		for _, a := range atoms {
			if got := ev.Evaluate(a.atom, a.value); got != resolve.GuardUnevaluable {
				t.Errorf("%s zero value: Evaluate(%s %s, %q) = %v; want GuardUnevaluable — "+
					"a nil-mapping evaluator never reverts to raw-string comparison",
					form, a.atom.Operator, a.atom.Literal, a.value, got)
			}
		}
	}
}

// REQ-8: "The seam signature `Evaluate(atom GuardAtom, value string)
// GuardResult` (`0007:C1`, REQ-10) and the atom shape (JDR 0001 §D1: Key,
// Operator, Literal, Block) are UNCHANGED" … "no kind is ever stamped into
// an atom, payload, or hash."
// BOUNDARY
func TestReq8_0012_TheSeamSignatureAndAtomShapeAreUnchanged(t *testing.T) {
	r12RequireProbe(t, "REQ-8", "TestRDR0012Probe_Req8_SeamSignatureAndAtomShapeAreUnchanged")
}

// r12GuardPath is internal/guard's import path.
const r12GuardPath = "github.com/cwensel/intrastate/internal/guard"

// r12ZeroSite is one zero-value construction of guard.Evaluator.
type r12ZeroSite struct {
	file string
	fn   string
	line int
	form string
}

var (
	r12SitesOnce sync.Once
	r12Sites     []r12ZeroSite
	r12SitesErr  error
)

// r12ZeroValueSites type-checks every non-test Go file in the module and
// returns each composite literal, zero-value `var` (a named result is one),
// `new`, `make`, or struct field whose type is guard.Evaluator (or an array
// of it), with its enclosing function. The type decides, not the spelling:
// an alias, a dot import and an elided literal type all resolve through
// go/types (`0012:S4` "a typed check, not a text grep").
func r12ZeroValueSites(t *testing.T) []r12ZeroSite {
	t.Helper()

	root := r12ModuleRoot(t)
	r12SitesOnce.Do(func() { r12Sites, r12SitesErr = r12TypedZeroValueSites(root) })
	if r12SitesErr != nil {
		t.Fatalf("type-checking the module's non-test code: %v", r12SitesErr)
	}
	return r12Sites
}

// r12ListedPackage is the slice of `go list -json` the typed check reads.
type r12ListedPackage struct {
	ImportPath string
	Dir        string
	GoFiles    []string
	Export     string
	DepOnly    bool
	Error      *struct{ Err string }
}

// r12TypedZeroValueSites type-checks the module's packages from source,
// importing their dependencies from the export data `go list -export`
// names, and collects the zero-value sites.
func r12TypedZeroValueSites(root string) ([]r12ZeroSite, error) {
	cmd := exec.Command("go", "list", "-e", "-deps", "-export",
		"-json=ImportPath,Dir,GoFiles,Export,DepOnly,Error", "./...")
	cmd.Dir = root
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("go list: %v\n%s", err, stderr.String())
	}

	exports := map[string]string{}
	var own []r12ListedPackage
	dec := json.NewDecoder(bytes.NewReader(out))
	for {
		var p r12ListedPackage
		if err := dec.Decode(&p); err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("decoding go list: %v", err)
		}
		if p.Error != nil {
			return nil, fmt.Errorf("%s: %s", p.ImportPath, p.Error.Err)
		}
		exports[p.ImportPath] = p.Export
		if !p.DepOnly {
			own = append(own, p)
		}
	}

	fset := token.NewFileSet()
	imp := importer.ForCompiler(fset, "gc", func(path string) (io.ReadCloser, error) {
		file := exports[path]
		if file == "" {
			return nil, fmt.Errorf("no export data for %s", path)
		}
		return os.Open(file)
	})

	var sites []r12ZeroSite
	for _, p := range own {
		var files []*ast.File
		for _, name := range p.GoFiles {
			f, err := parser.ParseFile(fset, filepath.Join(p.Dir, name), nil, parser.SkipObjectResolution)
			if err != nil {
				return nil, err
			}
			files = append(files, f)
		}
		info := &types.Info{
			Types: map[ast.Expr]types.TypeAndValue{},
			Defs:  map[*ast.Ident]types.Object{},
			Uses:  map[*ast.Ident]types.Object{},
		}
		conf := types.Config{Importer: imp}
		if _, err := conf.Check(p.ImportPath, fset, files, info); err != nil {
			return nil, fmt.Errorf("type-checking %s: %v", p.ImportPath, err)
		}
		for _, f := range files {
			sites = append(sites, r12FileZeroValueSites(root, fset, f, info)...)
		}
	}
	return sites, nil
}

// r12IsEvaluator reports whether t is guard.Evaluator, through any alias.
func r12IsEvaluator(t types.Type) bool {
	if t == nil {
		return false
	}
	n, ok := types.Unalias(t).(*types.Named)
	if !ok {
		return false
	}
	obj := n.Obj()
	return obj.Pkg() != nil && obj.Pkg().Path() == r12GuardPath && obj.Name() == "Evaluator"
}

// r12HoldsEvaluator reports whether a zero value of t is, or holds by value,
// a zero-value guard.Evaluator: the type itself or an array of it.
func r12HoldsEvaluator(t types.Type) bool {
	if r12IsEvaluator(t) {
		return true
	}
	if a, ok := types.Unalias(t).Underlying().(*types.Array); ok {
		return r12HoldsEvaluator(a.Elem())
	}
	return false
}

// r12FileZeroValueSites collects one type-checked file's zero-value sites.
func r12FileZeroValueSites(root string, fset *token.FileSet, file *ast.File, info *types.Info) []r12ZeroSite {
	var sites []r12ZeroSite
	record := func(fn string, pos token.Pos, form string) {
		p := fset.Position(pos)
		rel, _ := filepath.Rel(root, p.Filename)
		sites = append(sites, r12ZeroSite{file: filepath.ToSlash(rel), fn: fn, line: p.Line, form: form})
	}
	namedResults := func(fn string, ft *ast.FuncType) {
		if ft.Results == nil {
			return
		}
		for _, f := range ft.Results.List {
			if len(f.Names) > 0 && r12HoldsEvaluator(info.TypeOf(f.Type)) {
				record(fn, f.Pos(), "zero-value named result")
			}
		}
	}
	inspectIn := func(fn string, node ast.Node) {
		ast.Inspect(node, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.CompositeLit:
				t := info.TypeOf(x)
				if r12IsEvaluator(t) {
					record(fn, x.Pos(), "composite literal")
					break
				}
				// An array literal zero-fills every index it does not name;
				// duplicate indices do not compile, so fewer elements than
				// the length means at least one zero-value element.
				if t == nil {
					break
				}
				if a, ok := types.Unalias(t).Underlying().(*types.Array); ok &&
					r12HoldsEvaluator(a.Elem()) && int64(len(x.Elts)) < a.Len() {
					record(fn, x.Pos(), "array literal with omitted elements")
				}
			case *ast.ValueSpec:
				if len(x.Values) > 0 {
					break
				}
				for _, id := range x.Names {
					if obj, ok := info.Defs[id].(*types.Var); ok && r12HoldsEvaluator(obj.Type()) {
						record(fn, x.Pos(), "zero-value var")
						break
					}
				}
			case *ast.FuncLit:
				namedResults(fn, x.Type)
			case *ast.CallExpr:
				id, ok := ast.Unparen(x.Fun).(*ast.Ident)
				if !ok || len(x.Args) == 0 {
					break
				}
				builtin, ok := info.Uses[id].(*types.Builtin)
				arg := info.Types[x.Args[0]]
				if !ok || !arg.IsType() {
					break
				}
				switch builtin.Name() {
				case "new":
					if r12HoldsEvaluator(arg.Type) {
						record(fn, x.Pos(), "new")
					}
				case "make":
					if s, ok := arg.Type.Underlying().(*types.Slice); ok && r12HoldsEvaluator(s.Elem()) {
						record(fn, x.Pos(), "make")
					}
				}
			case *ast.StructType:
				for _, f := range x.Fields.List {
					if r12HoldsEvaluator(info.TypeOf(f.Type)) {
						record(fn, f.Pos(), "struct field")
					}
				}
			}
			return true
		})
	}
	for _, decl := range file.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok {
			name := fd.Name.Name
			if fd.Recv != nil && len(fd.Recv.List) == 1 {
				name = r12RecvName(fd.Recv.List[0].Type) + "." + name
			}
			namedResults(name, fd.Type)
			if fd.Body != nil {
				inspectIn(name, fd.Body)
			}
			continue
		}
		inspectIn("", decl)
	}
	return sites
}

func r12RecvName(e ast.Expr) string {
	switch x := e.(type) {
	case *ast.StarExpr:
		return r12RecvName(x.X)
	case *ast.Ident:
		return x.Name
	case *ast.IndexExpr:
		return r12RecvName(x.X)
	}
	return "?"
}

// r12AllowedZeroSite is the ONE named exception (`0012:C4`, REQ-35), plus
// NewEvaluator's own body.
func r12AllowedZeroSite(s r12ZeroSite) bool {
	if s.fn == "NewEvaluator" && strings.HasPrefix(s.file, "internal/guard/") {
		return true
	}
	return s.file == "internal/graphlint/reach.go" && s.fn == "atomAdmitsValue"
}

// REQ-6: "The zero-value `Evaluator{}` construction form is retired; every
// non-test construction goes through `NewEvaluator`."
// REQ-48: "no zero-value construction of `guard.Evaluator` outside
// `NewEvaluator` survives in non-test code, with ONE allow-listed
// exception: `internal/graphlint/reach.go::atomAdmitsValue`" … "The
// assertion matches the composite-literal form (`Evaluator{}`) *and* the
// `var` / `new` / embedded-field zero values Go equally admits" … "The
// allow-list is ONE named function, asserted by name rather than by
// pattern, so a second untyped site cannot slip in under it and moving or
// renaming that function fails the check"
// ADVERSARIAL
func TestReq48_0012_NoZeroValueEvaluatorSurvivesOutsideTheOneAllowListedFunction(t *testing.T) {
	sites := r12ZeroValueSites(t)

	// The check must not be vacuous: the allow-listed function is still a
	// zero-value site, found by name at its home.
	allowed := false
	for _, s := range sites {
		if s.file == "internal/graphlint/reach.go" && s.fn == "atomAdmitsValue" {
			allowed = true
		}
	}
	if !allowed {
		t.Error("internal/graphlint/reach.go::atomAdmitsValue no longer constructs a zero-value " +
			"guard.Evaluator (or was moved/renamed); the allow-list names that ONE function, so " +
			"moving or renaming it fails the check")
	}

	for _, s := range sites {
		if r12AllowedZeroSite(s) {
			continue
		}
		where := s.fn
		if where == "" {
			where = "package scope"
		}
		t.Errorf("%s:%d (%s): %s of guard.Evaluator; every non-test construction goes "+
			"through NewEvaluator", s.file, s.line, where, s.form)
	}
}

// REQ-35: "`internal/graphlint/reach.go::atomAdmitsValue` is NOT a
// construction site under this clause and keeps its `main` shape: it
// constructs `guard.Evaluator{}` with no mapping and compares raw bytes.
// That is required, not tolerated." … "`atomAdmitsValue` keeps its `bool`
// return and its name."
// BOUNDARY
func TestReq35_0012_AtomAdmitsValueIsTheSoleRawByteSiteAndKeepsItsShape(t *testing.T) {
	dir := filepath.Join(r12ModuleRoot(t), "internal", "graphlint")
	funcs, fset := r12PackageFuncs(t, dir)
	fn, ok := funcs["atomAdmitsValue"]
	if !ok {
		t.Fatal("internal/graphlint declares no atomAdmitsValue; it keeps its name")
	}
	if got := fset.Position(fn.Pos()).Filename; filepath.Base(got) != "reach.go" {
		t.Errorf("atomAdmitsValue lives in %s; it keeps its home in reach.go", got)
	}
	if got := r12FieldTypes(t, fset, fn.Type.Results); !slices.Equal(got, []string{"bool"}) {
		t.Errorf("atomAdmitsValue results = %v; it keeps its bool return", got)
	}
	body := r12Render(t, fset, fn.Body)
	if !strings.Contains(body, "guard.Evaluator{}") {
		t.Error("atomAdmitsValue no longer constructs guard.Evaluator{} with no mapping; it " +
			"is required to keep comparing raw bytes")
	}
	if strings.Contains(body, "NewEvaluator") || strings.Contains(body, "DeclaredKinds") {
		t.Error("atomAdmitsValue constructs over declarations; it is NOT a construction site " +
			"under C4 and keeps its main shape")
	}

	// And it is the ONLY one: every other construction is typed.
	for _, s := range r12ZeroValueSites(t) {
		if !r12AllowedZeroSite(s) {
			t.Errorf("%s:%d (%s) is a second untyped construction; atomAdmitsValue is the "+
				"record's one raw-byte site", s.file, s.line, s.fn)
		}
	}
}

// r12CallsNewEvaluatorOverDeclaredKinds reports whether node contains a
// call NewEvaluator(DeclaredKinds(<arg>)), qualified or not.
func r12CallsNewEvaluatorOverDeclaredKinds(node ast.Node) bool {
	name := func(e ast.Expr) string {
		switch x := e.(type) {
		case *ast.Ident:
			return x.Name
		case *ast.SelectorExpr:
			return x.Sel.Name
		}
		return ""
	}
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok || name(call.Fun) != "NewEvaluator" || len(call.Args) != 1 {
			return true
		}
		inner, ok := call.Args[0].(*ast.CallExpr)
		if ok && name(inner.Fun) == "DeclaredKinds" && len(inner.Args) == 1 {
			found = true
		}
		return true
	})
	return found
}

// REQ-33 (lint half): "EVERY non-test evaluator construction site … constructs
// over the declarations of the SAME loaded model whose rows it evaluates" …
// "On `main` that set is exactly two: `internal/cli/flow_resolve.go::guardSeam`,
// `internal/guard/product.go::valueSatisfies` (A3)."
// HAPPY PATH
func TestReq33_0012_TheLintSiteConstructsOverTheModelsDeclarations(t *testing.T) {
	path := filepath.Join(r12GuardDir(t), "product.go")
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		t.Fatalf("parsing product.go: %v", err)
	}
	if !r12CallsNewEvaluatorOverDeclaredKinds(file) {
		t.Error("internal/guard/product.go never constructs NewEvaluator(DeclaredKinds(m)); the " +
			"lint site evaluates over the loaded model's declarations")
	}
	for _, s := range r12ZeroValueSites(t) {
		if s.file == "internal/guard/product.go" {
			t.Errorf("product.go:%d (%s) still constructs a zero-value evaluator", s.line, s.fn)
		}
	}
}

// --- C2: the typed comparison ----------------------------------------------

// REQ-11: "TYPED COMPARISON APPLIES ON THE GUARD PATH. `eq` and `in` GUARD
// atoms resolve the atom's key against the constructed kind mapping and
// compare under the declared kind." (seam half; the match half is pinned
// end to end in internal/cli)
// HAPPY PATH
func TestReq11_0012_GuardEqAndInCompareUnderTheDeclaredKind(t *testing.T) {
	r12RequireProbe(t, "REQ-11", "TestRDR0012Probe_Req46_PerKindDispatchMatrix")
}

// REQ-12: "kind `int`: held value AND literal (each member, for `in`) must
// parse as integers; comparison is over the PARSED values. A present held
// value that does not parse is GuardUnevaluable — never GuardFalse."
// HAPPY PATH
func TestReq12_0012_IntComparesParsedValues(t *testing.T) {
	r12RequireProbe(t, "REQ-12", "TestRDR0012Probe_Req12_IntComparesParsedValues")
}

// REQ-13: "For `in`, ONE unparseable member poisons the whole list:
// GuardUnevaluable even when the held value equals a member that did
// parse."
// ADVERSARIAL
func TestReq13_0012_OneUnparseableInMemberPoisonsTheList(t *testing.T) {
	r12RequireProbe(t, "REQ-13", "TestRDR0012Probe_Req13_OneUnparseableInMemberPoisonsTheList")
}

// REQ-14: "The parse is `strconv.Atoi` … an out-of-range spelling is
// GuardUnevaluable on the target under test, and the suite pins the
// in-range cases only."
// BOUNDARY
func TestReq14_0012_OverflowIsUnevaluable(t *testing.T) {
	r12RequireProbe(t, "REQ-14", "TestRDR0012Probe_Req14_OverflowIsUnevaluableAndInRangeIsPinned")
}

// REQ-15: "kind `bool`: held value and literal must be the boolean tokens
// (`true` | `false`) … any other held value is GuardUnevaluable. The
// rejected set is every other spelling — `\"1\"`, `\"0\"`, `\"True\"`,
// `\"yes\"`"
// INPUT EDGE
func TestReq15_0012_BoolComparesTokensOnly(t *testing.T) {
	r12RequireProbe(t, "REQ-15", "TestRDR0012Probe_Req15_BoolComparesTokensOnly")
}

// REQ-16: "kind `enum` | `scalar`: every string parses; comparison is exact
// string equality (unchanged behavior). These two SHARE one arm"
// HAPPY PATH
func TestReq16_0012_EnumAndScalarCompareExactStrings(t *testing.T) {
	r12RequireProbe(t, "REQ-16", "TestRDR0012Probe_Req16_EnumAndScalarCompareExactStrings")
}

// REQ-17: "kind `set`: the operator/kind matrix does not admit `eq`/`in`
// over `set`; handed one anyway, the seam answers GuardUnevaluable"
// DOMAIN EDGE
func TestReq17_0012_SetEqAndInAreUnevaluable(t *testing.T) {
	r12RequireProbe(t, "REQ-17", "TestRDR0012Probe_Req17_SetEqAndInAreUnevaluable")
}

// REQ-18: "key absent from the mapping: GuardUnevaluable — the seam cannot
// type the comparison."
// DOMAIN EDGE
func TestReq18_0012_AbsentKeyIsUnevaluable(t *testing.T) {
	r12RequireProbe(t, "REQ-18", "TestRDR0012Probe_Req18_AbsentKeyIsUnevaluable")
}

// REQ-19: "For such a key this is the ordinary arm, not a defensive one:
// `eq`/`in` over it answers GuardUnevaluable where today it compares raw
// strings." ("such a key" = a key declared with `Kind` the empty string)
// DOMAIN EDGE
func TestReq19_0012_EmptyKindDeclaredKeyIsUnevaluable(t *testing.T) {
	r12RequireProbe(t, "REQ-19", "TestRDR0012Probe_Req19_EmptyKindDeclaredKeyIsUnevaluable")
}

// REQ-20: "unparseable LITERAL: GuardUnevaluable (defense in depth …)"
// ADVERSARIAL
func TestReq20_0012_UnparseableLiteralIsUnevaluable(t *testing.T) {
	r12RequireProbe(t, "REQ-20", "TestRDR0012Probe_Req20_UnparseableLiteralIsUnevaluable")
}

// REQ-21: "Dispatch over the kind token is EXHAUSTIVE across the five-kind
// vocabulary, and the fallthrough `default` means GuardUnevaluable" … "so
// the suite asserts an unknown kind token answers GuardUnevaluable rather
// than comparing."
// DOMAIN EDGE
func TestReq21_0012_UnknownKindTokenIsUnevaluable(t *testing.T) {
	r12RequireProbe(t, "REQ-21", "TestRDR0012Probe_Req21_UnknownKindTokenIsUnevaluable")
	r12RequireProbe(t, "REQ-21", "TestRDR0012Probe_Req21_TheSuitePinsAnUnknownKindToken")
}

// REQ-22: "`lt/lte/gt/gte` keep their operator-inferred integer parse …
// `contains` and the §D13 set arms are unchanged … the kind lookup is not
// consulted for those operators at all" + "a nil-mapping evaluator still
// answers them normally"
// BOUNDARY
func TestReq22_0012_OrderingAndContainsDoNotConsultTheMapping(t *testing.T) {
	r12RequireProbe(t, "REQ-22", "TestRDR0012Probe_Req22_OrderingAndContainsDoNotConsultTheMapping")
}

// REQ-23: "an operator the seam does not recognize — `exists` included …
// answers GuardUnevaluable like any other atom it cannot type. The seam
// NEVER panics"
// ADVERSARIAL
func TestReq23_0012_UnknownOperatorIsUnevaluableAndNeverPanics(t *testing.T) {
	r12RequireProbe(t, "REQ-23", "TestRDR0012Probe_Req23_UnknownOperatorIsUnevaluableAndNeverPanics")
}

// REQ-24: "\"Each member, for `in`\" above means the members decoded from
// that array, and parse-all-then-compare quantifies over them. No new
// members field and no delimiter convention is introduced."
// INPUT EDGE
func TestReq24_0012_InMembersAreTheDecodedArray(t *testing.T) {
	r12RequireProbe(t, "REQ-24", "TestRDR0012Probe_Req24_InMembersAreTheDecodedArray")
}

// REQ-46: "**Expected**: the C2 matrix verbatim." over "Per-kind dispatch at
// the seam — `int`, `bool`, `enum`, `scalar`, `set`, and key-absent — over
// `eq` and `in`."
// HAPPY PATH
func TestReq46_0012_PerKindDispatchMatrix(t *testing.T) {
	r12RequireProbe(t, "REQ-46", "TestRDR0012Probe_Req46_PerKindDispatchMatrix")
}

// --- C3: the conformance suite ---------------------------------------------

// REQ-25: "`eq`/`in` over an `int` and a `bool` dimension where the held
// value does not parse (want GuardUnevaluable) and where it parses but
// differs (want GuardFalse); plus the defensive arms — an `in` list with a
// non-integer member on an `int` dimension, an integer-overflow held value,
// `eq` over a `set`-kind key, and a key absent from the fixture (each want
// GuardUnevaluable)."
// ADVERSARIAL
func TestReq25_0012_TheSuiteCarriesTheTypedAndDefensiveLegs(t *testing.T) {
	r12RequireProbe(t, "REQ-25", "TestRDR0012Probe_Req25_TheSuiteAsksTheNewLegs")
	r12RequireProbeFails(t, "REQ-25", "TestRDR0012Probe_ReqMVV_Mutant_rawStringSeam",
		"the int/bool want-Unevaluable legs must fail against today's raw-string arms")
}

// REQ-26: "The suite carries at least one discriminating case per kind
// token, and a meta-check that the fixture's kind tokens are exactly the
// five-kind vocabulary."
// ADVERSARIAL
func TestReq26_0012_EachKindTokenHasADiscriminatingCase(t *testing.T) {
	r12RequireProbe(t, "REQ-26", "TestRDR0012Probe_Req26_FixtureTokensAreExactlyTheVocabulary")
	for _, kind := range []string{"int", "bool", "enum", "scalar", "set"} {
		r12RequireProbeFails(t, "REQ-26", "TestRDR0012Probe_Req26_Mutant_"+kind,
			"no case discriminates the `"+kind+"` arm")
	}
}

// REQ-27 (behaviour): "It is published as a FUNCTION returning a fresh
// map" … "A constructor hands each caller its own copy."
// ADVERSARIAL
func TestReq27_0012_ContractKindsHandsEachCallerAFreshCopy(t *testing.T) {
	r12RequireProbe(t, "REQ-27", "TestRDR0012Probe_Req27_ContractKindsHandsEachCallerItsOwnCopy")
}

// REQ-29: "each case names its own key drawn from the fixture, chosen so the
// case's operator is one the matrix admits for that kind." … "their verdicts
// are unchanged"
// BOUNDARY
func TestReq29_0012_CasesAreKeyedOntoAdmittedKinds(t *testing.T) {
	r12RequireProbe(t, "REQ-29", "TestRDR0012Probe_Req29_CasesAreKeyedOntoAdmittedKinds")
}

// r12ContractCalls returns every call to resolve.TestGuardEvaluatorContract
// in the named guard test file.
func r12ContractCalls(t *testing.T, name string) ([]*ast.CallExpr, *token.FileSet) {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filepath.Join(r12GuardDir(t), name), nil, 0)
	if err != nil {
		t.Fatalf("parsing %s: %v", name, err)
	}
	var calls []*ast.CallExpr
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if ok && sel.Sel.Name == "TestGuardEvaluatorContract" {
			calls = append(calls, call)
		}
		return true
	})
	return calls, fset
}

// REQ-31: "`guard.Evaluator`'s existing suite call is RE-POINTED, not
// duplicated." … "Each migrating site therefore supplies a one-line
// adapter: func(k map[string]string) resolve.GuardEvaluator { return
// guard.NewEvaluator(k) }"
// BOUNDARY
func TestReq31_0012_TheGuardSuiteCallIsRePointedThroughAOneLineAdapter(t *testing.T) {
	calls, fset := r12ContractCalls(t, "guard_evaluator_0003_test.go")
	if len(calls) != 1 {
		t.Fatalf("guard_evaluator_0003_test.go calls TestGuardEvaluatorContract %d times; "+
			"the existing call is RE-POINTED, not duplicated (want exactly 1)", len(calls))
	}
	call := calls[0]
	if len(call.Args) != 2 {
		t.Fatalf("the suite call passes %d arguments; want (t, newSeam)", len(call.Args))
	}
	lit, ok := call.Args[1].(*ast.FuncLit)
	if !ok {
		t.Fatalf("the suite call passes %s; want the one-line adapter "+
			"func(k map[string]string) resolve.GuardEvaluator { return guard.NewEvaluator(k) }",
			r12Render(t, fset, call.Args[1]))
	}
	if got := r12FieldTypes(t, fset, lit.Type.Params); !slices.Equal(got, []string{"map[string]string"}) {
		t.Errorf("adapter parameters = %v; want (k map[string]string)", got)
	}
	if got := r12FieldTypes(t, fset, lit.Type.Results); !slices.Equal(got, []string{"resolve.GuardEvaluator"}) {
		t.Errorf("adapter results = %v; want resolve.GuardEvaluator", got)
	}
	if len(lit.Body.List) != 1 {
		t.Fatalf("adapter body has %d statements; it is ONE line", len(lit.Body.List))
	}
	ret, ok := lit.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 1 {
		t.Fatal("adapter body is not a single return")
	}
	got := r12Render(t, fset, ret.Results[0])
	param := ""
	if names := lit.Type.Params.List[0].Names; len(names) == 1 {
		param = names[0].Name
	}
	if got != "guard.NewEvaluator("+param+")" {
		t.Errorf("adapter returns %s; want guard.NewEvaluator(%s)", got, param)
	}
}

// --- MVV step 4 -------------------------------------------------------------

// REQ-MVV (step 4): "`TestGuardEvaluatorContract`, extended per C3 and
// driven against `guard.NewEvaluator` constructed over the published
// fixture kinds, is green — including the int/bool want-Unevaluable legs
// that fail against today's raw-string arms."
// HAPPY PATH
func TestReqMVV_0012_Step4_TheExtendedSuiteIsGreenOverNewEvaluator(t *testing.T) {
	r12RequireProbe(t, "REQ-MVV", "TestRDR0012Probe_ReqMVV_ContractSuiteIsGreenOverNewEvaluator")
	r12RequireProbeFails(t, "REQ-MVV", "TestRDR0012Probe_ReqMVV_Mutant_rawStringSeam",
		"the want-Unevaluable legs must fail against today's raw-string arms")
}
