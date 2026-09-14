package cli

// RDR 0021 — the probe layer shared by every `graph` export suite.
//
// EVERY surface this RDR adds is ABSENT from the tree these tests are
// committed against: the root `graph` verb, its `--emit` flag, the
// `intrastate.graph/1` document, the `graph-export-too-large` refusal, the
// `--emit` enumeration accessor, and the edges-and-completeness-carrying
// `graphlint` function C4 requires. A test naming any of them as a Go
// identifier would not compile, and `.githooks/pre-commit` runs `gofmt -l`
// and `go vet ./...` under `set -e` — so a red suite that breaks the build
// cannot be committed at all.
//
// The probes below therefore reach the unwritten surface two ways, both of
// which compile against today's tree:
//
//   - BEHAVIOURALLY, by driving `ExecuteAndEmit` with the string "graph".
//     Today cobra refuses it as an unknown command at RUN time, so each
//     assertion fails on its own oracle rather than on a link error.
//   - STRUCTURALLY, by parsing package source with `go/parser`, which reads
//     the tree as text — an absent function is a `false` return, never a
//     compile error. This is the pattern `schema_probe_0029_test.go`
//     established for RDR 0029's unwritten accessors.
//
// Every failure here is a runtime assertion, and every one turns green only
// when the real surface lands. A probe that passed against an absent seam
// would be tautological, which the red-before-green gate exists to prevent.

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- the verb under test -------------------------------------------------

// graphVerb is the root verb name C1 fixes. It is a STRING, not a Go
// identifier, which is what lets this suite compile before the verb exists.
const graphVerb = "graph"

// emitFlag is the document-format flag name C1 fixes (`--emit`).
const emitFlag = "emit"

// schemaMarker is C2's initial document version marker.
const schemaMarker = "intrastate.graph/1"

// abstractionMarker is the REQUIRED token value C2 fixes on the `reach`
// block's `abstraction` member.
const abstractionMarker = "declared-over-approximation"

// tooLargeCode is the scalar refusal code C4 mints for an incomplete
// traversal.
const tooLargeCode = "graph-export-too-large"

// --- structural probes ---------------------------------------------------

// parsePkg parses one package directory's PRODUCTION files. It reads the
// tree as text rather than linking against it, because the whole point is
// to ask about functions that do not exist yet.
func parsePkg(t *testing.T, dir string) map[string]*ast.Package {
	t.Helper()

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, parser.ParseComments)
	if err != nil {
		t.Fatalf("parsing %s: %v", dir, err)
	}
	return pkgs
}

// graphPkgDir resolves an internal package directory from the repo root.
func graphPkgDir(t *testing.T, parts ...string) string {
	t.Helper()
	return filepath.Join(append([]string{repoRootFor(t)}, parts...)...)
}

// exportedFuncs returns the exported package-level function names declared
// in the package rooted at dir, mapped to their declaration.
func exportedFuncs(t *testing.T, dir string) map[string]*ast.FuncDecl {
	t.Helper()

	out := map[string]*ast.FuncDecl{}
	for _, pkg := range parsePkg(t, dir) {
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

// resultTypes renders a function's declared result types as source text.
func resultTypes(t *testing.T, fn *ast.FuncDecl) []string {
	t.Helper()

	if fn.Type.Results == nil {
		return nil
	}
	fset := token.NewFileSet()
	var out []string
	for _, field := range fn.Type.Results.List {
		var buf strings.Builder
		// printNode is the package's existing AST renderer
		// (`schema_probe_0029_test.go`); this suite reuses it rather than
		// declaring a second one in the same package.
		if err := printNode(&buf, fset, field.Type); err != nil {
			t.Fatalf("rendering result type: %v", err)
		}
		// One field may declare several names sharing a type.
		n := len(field.Names)
		if n == 0 {
			n = 1
		}
		for range n {
			out = append(out, buf.String())
		}
	}
	return out
}

// pkgSource concatenates a package's production source, for the coarse
// "does this package mention X at all" probes.
func pkgSource(t *testing.T, dir string) string {
	t.Helper()

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading %s: %v", dir, err)
	}
	var b strings.Builder
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") ||
			strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		src, rerr := os.ReadFile(filepath.Join(dir, e.Name()))
		if rerr != nil {
			t.Fatalf("reading %s: %v", e.Name(), rerr)
		}
		b.Write(src)
		b.WriteString("\n")
	}
	return b.String()
}

// --- behavioural probes --------------------------------------------------

// runGraph drives the production emission path with the `graph` verb and
// the supplied arguments. Before the verb exists cobra refuses it as an
// unknown command, which surfaces here as a returned error — a RUNTIME
// failure each caller asserts against, never a compile break.
func runGraph(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	return runCmd(t, append([]string{graphVerb}, args...)...)
}

// graphRegistered reports whether a root `graph` command is registered.
// Every behavioural test calls this first, so a missing verb fails with
// one legible message rather than as a confusing envelope mismatch.
func graphRegistered(t *testing.T) bool {
	t.Helper()

	for _, c := range NewRootCmd().Commands() {
		if c.Name() == graphVerb {
			return true
		}
	}
	return false
}

// requireGraphVerb fails the calling test when the verb is absent. This is
// the single red gate every behavioural probe in this RDR routes through
// before it can assert anything about the document.
func requireGraphVerb(t *testing.T) {
	t.Helper()

	if graphRegistered(t) {
		return
	}
	var names []string
	for _, c := range NewRootCmd().Commands() {
		names = append(names, c.Name())
	}
	t.Fatalf("no root %q command is registered; root commands = %v. "+
		"`0021:C1` registers a new root verb `graph` beside `lint`, outside "+
		"the `flow` group, under `0005:C1`'s carve-out for command groups "+
		"\"owned by the RDR that names them\"", graphVerb, names)
}

// --- document decoding ---------------------------------------------------

// graphDocument is a DECODING view of the exported document. It deliberately
// does NOT enumerate C2's field set as a closed struct: C2's STABILITY
// clause forbids a consumer from asserting the field set's cardinality, a
// member's ordinal position, or a tail position, so the fields below are
// the ones individual REQs name and nothing more.
type graphDocument struct {
	Schema  string          `json:"schema"`
	Model   string          `json:"model"`
	Class   string          `json:"class"`
	Tags    []graphTag      `json:"tags"`
	Initial []graphTagValue `json:"initial"`
	Rows    []graphRow      `json:"rows"`
	Groups  []graphGroup    `json:"groups"`
	Reach   *graphReach     `json:"reach"`
	Raw     json.RawMessage `json:"-"`
}

// `initial` and a row's `writes` both decode into the producer's own
// `graphTagValue` ({key, value}), so a test can compare the exported members
// to the fixture's AUTHORED values (RT1: the quantifier is C2's field list,
// "not whatever was exported").

type graphTag struct {
	Name         string   `json:"name"`
	Provenance   string   `json:"provenance"`
	Kind         string   `json:"kind"`
	Required     bool     `json:"required"`
	SingleValued bool     `json:"single_valued"`
	Domain       []string `json:"domain"`
}

type graphRow struct {
	Identity string          `json:"identity"`
	Kind     string          `json:"kind"`
	Outcome  string          `json:"outcome"`
	Atoms    []graphAtom     `json:"atoms"`
	Writes   []graphTagValue `json:"writes"`
}

type graphAtom struct {
	Key      string   `json:"key"`
	Operator string   `json:"operator"`
	Literal  []string `json:"literal"`
	Block    string   `json:"block"`
}

type graphGroup struct {
	Context string   `json:"context"`
	Rules   []string `json:"rules"`
}

type graphReach struct {
	Abstraction string      `json:"abstraction"`
	Nodes       []graphNode `json:"nodes"`
	Edges       []graphEdge `json:"edges"`
}

type graphNode struct {
	ID     string              `json:"id"`
	Values map[string][]string `json:"values"`
}

type graphEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
	Rule string `json:"rule"`
}

// decodeDocument parses one exported JSON document.
func decodeDocument(t *testing.T, body string) graphDocument {
	t.Helper()

	var doc graphDocument
	trimmed := strings.TrimSpace(body)
	if err := json.Unmarshal([]byte(trimmed), &doc); err != nil {
		t.Fatalf("the exported document is not one JSON object: %v\n%s",
			err, body)
	}
	doc.Raw = json.RawMessage(trimmed)
	return doc
}

// decodeGeneric parses the document into a generic map, for the probes that
// ask about key PRESENCE and null-versus-absent rather than about values.
func decodeGeneric(t *testing.T, body string) map[string]any {
	t.Helper()

	var out map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(body)), &out); err != nil {
		t.Fatalf("the exported document is not one JSON object: %v\n%s",
			err, body)
	}
	return out
}
