package cli

// RDR 0029 — the reflective probe layer.
//
// Every surface this RDR adds is ABSENT from the tree these tests are
// committed against: `clierr.SchemaVersion`, `respond::Types`,
// `respond::Levels`, `clierr::ExitCodes`, `cli::UnknownReasons`, and
// `resolve::Blocks`. A test naming any of them directly would not compile,
// and the pre-commit hook runs `go vet ./...` under `set -e`, so a red
// suite that breaks the build cannot be committed at all.
//
// The probes below therefore reach the unwritten surface through the
// PACKAGE SOURCE and through values that exist today, never through an
// identifier that does not. Every failure here is a runtime assertion, and
// every one of them turns green only when the real accessor lands — a
// probe that passed against an absent seam would be tautological.

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// printNode renders one AST node back to source text, which is how a
// declared result type is compared against the census's spelling.
func printNode(w *strings.Builder, fset *token.FileSet, node ast.Node) error {
	return printer.Fprint(w, fset, node)
}

// declaredFuncResult reports the source-declared result type of the
// exported function `name` in the package rooted at `dir`, and whether such
// a function is declared at all.
//
// It parses the package rather than linking against it because the whole
// point is to ask about functions that do not exist yet. `go/parser` reads
// the tree as text, so an absent accessor is a false return rather than a
// compile error.
func declaredFuncResult(t *testing.T, dir, name string) (string, bool) {
	t.Helper()

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Fatalf("parsing %s: %v", dir, err)
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Name == nil || fn.Name.Name != name {
					continue
				}
				// A method is not a package-level accessor.
				if fn.Recv != nil {
					continue
				}
				if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
					return "", true
				}
				var buf strings.Builder
				if err := printNode(&buf, fset, fn.Type.Results.List[0].Type); err != nil {
					t.Fatalf("rendering %s's result: %v", name, err)
				}
				return buf.String(), true
			}
		}
	}
	return "", false
}

// declaredConstNames returns the exported package-level constant names
// declared in the package rooted at dir.
func declaredConstNames(t *testing.T, dir string) map[string]bool {
	t.Helper()

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Fatalf("parsing %s: %v", dir, err)
	}

	out := map[string]bool{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				gd, ok := decl.(*ast.GenDecl)
				if !ok || gd.Tok != token.CONST {
					continue
				}
				for _, spec := range gd.Specs {
					vs, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					for _, id := range vs.Names {
						if id.IsExported() {
							out[id.Name] = true
						}
					}
				}
			}
		}
	}
	return out
}

// pkgDir resolves an internal package directory from the repo root.
func pkgDir(t *testing.T, parts ...string) string {
	t.Helper()
	return filepath.Join(append([]string{repoRootFor(t)}, parts...)...)
}

// REQ-3: "The value has ONE home: a single exported constant, which both
// terminal records read and neither declares."
// REQ-4: "the constant lives in `clierr` (the leaf both can reach) and
// `respond` reads it from there"
// HAPPY PATH
func TestReq3And4_TheSchemaVersionConstantHasOneHomeInClierr(t *testing.T) {
	consts := declaredConstNames(t, pkgDir(t, "internal", "cli", "clierr"))

	// The name is not fixed by the record as a Go identifier, so the probe
	// accepts the spellings a reasonable implementation would choose and
	// fails when NONE of them is declared. Pinning one spelling would make
	// this a naming test rather than a homing test.
	candidates := []string{"SchemaVersion", "SchemaVersionValue", "WireSchemaVersion"}
	var found string
	for _, c := range candidates {
		if consts[c] {
			found = c
			break
		}
	}
	if found == "" {
		t.Fatalf("internal/cli/clierr declares no exported schema-version "+
			"constant (looked for %v). C1 gives the value ONE home: a single "+
			"exported constant in `clierr`, the leaf both terminal records "+
			"can reach, which `respond` reads rather than redeclaring. A "+
			"literal duplicated across the two records could drift, which is "+
			"the \"one key meaning two things on one wire\" defect this RDR "+
			"exists to prevent.", candidates)
	}

	// And `respond` must READ it rather than declare its own. A second
	// declaration in `respond` is exactly the drift C1 forbids.
	respondConsts := declaredConstNames(t, pkgDir(t, "internal", "cli", "respond"))
	for _, c := range candidates {
		if respondConsts[c] {
			t.Errorf("internal/cli/respond declares %q itself; C1 fixes ONE "+
				"home for the value and `respond` reads it from `clierr`", c)
		}
	}

	src := readSource(t, pkgDir(t, "internal", "cli", "respond", "respond.go"))
	if !strings.Contains(src, "clierr."+found) {
		t.Errorf("internal/cli/respond/respond.go does not read clierr.%s; "+
			"the `ok` envelope must carry the value from its one home, not a "+
			"literal of its own", found)
	}
}

// REQ-2: "It is versioned independently of the binary's release version
// and MUST NOT be derived from it."
// ADVERSARIAL — the defect is a schema version wired to build identity,
// which would move the wire contract on every release that changes nothing
// about the wire.
func TestReq2_TheSchemaVersionIsNotDerivedFromTheBinaryVersion(t *testing.T) {
	clierrDir := pkgDir(t, "internal", "cli", "clierr")

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, clierrDir, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Fatalf("parsing clierr: %v", err)
	}
	for _, pkg := range pkgs {
		for path, file := range pkg.Files {
			for _, imp := range file.Imports {
				if strings.Contains(imp.Path.Value, "internal/version") {
					t.Errorf("%s imports internal/version; the schema version "+
						"is versioned independently of the binary's release "+
						"version and MUST NOT be derived from it",
						filepath.Base(path))
				}
			}
		}
	}

	// The subject must EXIST before "not derived from the binary version"
	// says anything: with no schema-version constant at all, the import
	// check above is satisfied by a package that simply has no version to
	// derive.
	consts := declaredConstNames(t, clierrDir)
	var declared bool
	for _, c := range []string{
		"SchemaVersion", "SchemaVersionValue", "WireSchemaVersion",
	} {
		if consts[c] {
			declared = true
			break
		}
	}
	if !declared {
		t.Fatal("internal/cli/clierr declares no schema-version constant, so " +
			"\"versioned independently of the binary's release version\" has " +
			"no subject — the import assertion above passes vacuously on a " +
			"package that carries no version at all")
	}
}

// REQ-31: "`respond::Types() []string` over `{ok}`"
// REQ-33: "`respond::Levels() []string`"
// REQ-32: "`clierr::ExitCodes() []int` over the five values `{0,1,2,3,130}`"
// REQ-34: "`cli::UnknownReasons() []string` over the union, matching
// `graphlint::Reasons`"
// REQ-35: "`resolve::Blocks() []Block`"
// REQ-30: "A tier assignment obliges an ENUMERATION SEAM: an exported
// accessor returning the vocabulary's members, in the package that owns
// them."
// REQ-61: "Build the enumeration seams C4's census obliges for the five
// vocabularies that admit one … The seams are production accessors, not
// test helpers"
// HAPPY PATH
func TestReq30And31Thru35And61_TheFiveOwedEnumerationSeamsExist(t *testing.T) {
	for _, tc := range []struct {
		req    string
		dir    []string
		name   string
		result string
	}{
		{"REQ-31", []string{"internal", "cli", "respond"}, "Types", "[]string"},
		{"REQ-33", []string{"internal", "cli", "respond"}, "Levels", "[]string"},
		{"REQ-32", []string{"internal", "cli", "clierr"}, "ExitCodes", "[]int"},
		{"REQ-34", []string{"internal", "cli"}, "UnknownReasons", "[]string"},
		{"REQ-35", []string{"internal", "resolve"}, "Blocks", "[]Block"},
	} {
		t.Run(tc.req+"/"+tc.name, func(t *testing.T) {
			got, ok := declaredFuncResult(t, pkgDir(t, tc.dir...), tc.name)
			if !ok {
				t.Fatalf("%s declares no exported %s; C4's census owes this "+
					"seam, and without it the tier is unassertable — a "+
					"`frozen` set has nothing to compare by value and an "+
					"`append-only` set has nothing to check membership and "+
					"uniqueness over, so C2's tiers would bind only prose",
					filepath.Join(tc.dir...), tc.name)
			}
			if got != tc.result {
				t.Errorf("%s returns %q; the census fixes %s() %s",
					tc.name, got, tc.name, tc.result)
			}
		})
	}
}

// REQ-61: "The seams are production accessors, not test helpers"
// BOUNDARY — a seam declared in a _test.go file satisfies no consumer.
func TestReq61_TheSeamsAreProductionAccessorsNotTestHelpers(t *testing.T) {
	// declaredFuncResult parses production files only, so a seam found by
	// it is by construction not a test helper. The discriminating half is
	// that the same name must NOT be findable only in test files: this
	// asserts the production tree carries them, which the probe above
	// already ranges over. Here we pin the complementary property — the
	// accessor is exported, so a consumer outside the package can call it.
	for _, tc := range []struct {
		dir  []string
		name string
	}{
		{[]string{"internal", "cli", "respond"}, "Types"},
		{[]string{"internal", "cli", "respond"}, "Levels"},
		{[]string{"internal", "cli", "clierr"}, "ExitCodes"},
		{[]string{"internal", "cli"}, "UnknownReasons"},
		{[]string{"internal", "resolve"}, "Blocks"},
	} {
		if !ast.IsExported(tc.name) {
			t.Fatalf("%s is not an exported name", tc.name)
		}
		if _, ok := declaredFuncResult(t, pkgDir(t, tc.dir...), tc.name); !ok {
			t.Errorf("%s::%s is not declared in production source; C4 requires "+
				"an exported accessor in the package that OWNS the vocabulary, "+
				"not a helper beside a test", filepath.Join(tc.dir...), tc.name)
		}
	}
}

// readSource reads one source file, failing the test if it cannot.
func readSource(t *testing.T, path string) string {
	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}
