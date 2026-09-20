package graphlint_test

// RDR 0006 — authority, placement, and the lint engine boundary
// (REQ-1..REQ-9) plus the minimum input contract and pipeline
// (REQ-15..REQ-18).

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/graphlint"
	"github.com/cwensel/intrastate/internal/guard"
	"github.com/cwensel/intrastate/internal/table"
)

// repoRoot walks up from the test's working directory to the module root.
func repoRoot(t *testing.T) string {
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
			t.Fatal("no go.mod above the test working directory")
		}
		dir = parent
	}
}

// exportedFuncsIn parses every non-test .go file in pkgDir and returns the
// exported top-level function names.
func exportedFuncsIn(t *testing.T, pkgDir string) []string {
	t.Helper()

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgDir, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", pkgDir, err)
	}
	var out []string
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Recv != nil || !fn.Name.IsExported() {
					continue
				}
				out = append(out, fn.Name.Name)
			}
		}
	}
	slices.Sort(out)
	return out
}

// REQ-1: "Graph lint MUST be a blocking acceptance gate over the
// normalized transition model. A model with any blocking lint finding
// MUST NOT be accepted for CI success"
// HAPPY PATH
func TestReq1_AnyBlockingFindingRefusesAcceptance(t *testing.T) {
	// A model with a dangling reference carries a blocking finding, so the
	// gate must refuse it: a blocking finding and an accepted model are
	// mutually exclusive dispositions.
	r := lint(t, twoStateDecls, danglingBody)

	if len(r.Blocking()) == 0 {
		t.Fatalf("a model with a dangling reference produced no blocking "+
			"finding, so the gate accepts it; report:%s", render(r))
	}
	for _, f := range r.Blocking() {
		if f.Severity != graphlint.SeverityBlocking {
			t.Errorf("Blocking() returned %s with severity %q; want %q",
				f.Code, f.Severity, graphlint.SeverityBlocking)
		}
	}
	// The clean counterpart is accepted, so the gate is a gate and not a
	// blanket refusal.
	clean := lint(t, twoStateDecls, legalBody)
	requireClean(t, clean)
}

// REQ-2: "That boundary is the enforcement site and the only one: the
// kernel does not consult a lint verdict at runtime"
// ADVERSARIAL
func TestReq2_ResolverConsultsNoLintVerdict(t *testing.T) {
	// Negative REQ: no resolver-side lint consultation may be implemented.
	// The mechanical oracle is the import graph — internal/resolve must
	// not reach the graph-lint package, in either direction of naming.
	root := repoRoot(t)
	resolveDir := filepath.Join(root, "internal", "resolve")

	entries, err := os.ReadDir(resolveDir)
	if err != nil {
		t.Fatalf("read %s: %v", resolveDir, err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		src, err := os.ReadFile(filepath.Join(resolveDir, e.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", e.Name(), err)
		}
		if strings.Contains(string(src), "internal/graphlint") {
			t.Errorf("internal/resolve/%s imports the graph-lint package; "+
				"the kernel MUST NOT consult a lint verdict at runtime",
				e.Name())
		}
	}
}

// REQ-3: "Graph lint MUST consume the normalized candidate-row graph from
// the transition model contract. It MUST NOT define a second sparse-source
// parser or a parallel transition semantics."
// ADVERSARIAL
func TestReq3_EngineDefinesNoSecondSourceParser(t *testing.T) {
	root := repoRoot(t)
	pkgDir := filepath.Join(root, "internal", "graphlint")

	entries, err := os.ReadDir(pkgDir)
	if err != nil {
		t.Fatalf("read %s: %v", pkgDir, err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") ||
			strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		src, err := os.ReadFile(filepath.Join(pkgDir, e.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", e.Name(), err)
		}
		for _, banned := range []string{"go-toml", "encoding/toml", "toml.Unmarshal", "toml.NewDecoder"} {
			if strings.Contains(string(src), banned) {
				t.Errorf("internal/graphlint/%s references %q; the engine "+
					"MUST NOT define a second sparse-source parser",
					e.Name(), banned)
			}
		}
	}

	// The positive half: the engine consumes the normalized row graph the
	// table contract produces, so a normalized model is all it needs.
	m := mustLoad(t, source(twoStateDecls, legalBody))
	req := graphlint.NewRequest(m)
	if req.Model != m {
		t.Errorf("NewRequest did not carry the normalized model through; "+
			"the engine consumes the normalized candidate-row graph, "+
			"got %v", req.Model)
	}
}

// REQ-4: "The authoritative CLI surface for graph acceptance MUST be the
// root command `intrastate lint` or a same-engine CI invocation of that
// command."
// HAPPY PATH
func TestReq4_RootLintIsTheAuthoritativeSurface(t *testing.T) {
	// Asserted from the CLI side in internal/cli; here the engine half:
	// the authoritative surface routes through this package's one entry
	// point, so the engine must expose it.
	funcs := exportedFuncsIn(t, filepath.Join(repoRoot(t), "internal", "graphlint"))
	if !slices.Contains(funcs, "Run") {
		t.Fatalf("the graph-lint package exposes no Run entry point; "+
			"exported funcs=%v", funcs)
	}
}

// REQ-7: "the lint package exposes exactly **one** exported engine entry
// point and one request builder, and a test asserts the root command calls
// them"
// BOUNDARY
func TestReq7_ExactlyOneEntryPointAndOneRequestBuilder(t *testing.T) {
	funcs := exportedFuncsIn(t, filepath.Join(repoRoot(t), "internal", "graphlint"))

	// The engine's ENTRY POINTS are the functions taking a Request or a
	// model and producing a verdict. The clause is that there is exactly
	// one of each, so a later alias has no second constructor to diverge
	// through. The named pair is Run + NewRequest; any second function
	// whose name reads as an alternative entry point or builder breaks it.
	var entries, builders []string
	for _, name := range funcs {
		switch {
		case name == "Run" || strings.HasPrefix(name, "Lint") ||
			strings.HasPrefix(name, "Check") || strings.HasPrefix(name, "Analyze"):
			entries = append(entries, name)
		case strings.HasPrefix(name, "NewRequest") || strings.HasPrefix(name, "BuildRequest") ||
			strings.HasPrefix(name, "RequestFor"):
			builders = append(builders, name)
		}
	}
	if len(entries) != 1 {
		t.Errorf("the lint package exposes %d engine entry points %v; "+
			"want exactly one", len(entries), entries)
	}
	if len(builders) != 1 {
		t.Errorf("the lint package exposes %d request builders %v; "+
			"want exactly one", len(builders), builders)
	}
}

// REQ-8: "The lint engine boundary is an internal graph-lint package that
// receives a normalized graph value, not Cobra command state and not
// sparse TOML."
// BOUNDARY
func TestReq8_EngineReceivesANormalizedGraphValueNotCommandState(t *testing.T) {
	root := repoRoot(t)
	pkgDir := filepath.Join(root, "internal", "graphlint")

	entries, err := os.ReadDir(pkgDir)
	if err != nil {
		t.Fatalf("read %s: %v", pkgDir, err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") ||
			strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		src, err := os.ReadFile(filepath.Join(pkgDir, e.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", e.Name(), err)
		}
		for _, banned := range []string{"spf13/cobra", "spf13/pflag", "internal/cli/respond"} {
			if strings.Contains(string(src), banned) {
				t.Errorf("internal/graphlint/%s imports %q; the engine "+
					"receives a normalized graph value, not Cobra command "+
					"state", e.Name(), banned)
			}
		}
	}

	// The positive half: a Request is constructible from a normalized
	// model alone, with no command in sight, and the engine runs on it.
	m := mustLoad(t, source(twoStateDecls, legalBody))
	r := graphlint.Run(graphlint.NewRequest(m))
	if r.ModelID != m.ID {
		t.Errorf("report ModelID = %q; want the normalized model's id %q",
			r.ModelID, m.ID)
	}
}

// REQ-9: "the canonical command and subsystem name is \"lint\"."
// HAPPY PATH
func TestReq9_CanonicalSubsystemNameIsLint(t *testing.T) {
	// The subsystem half: every finding code this RDR mints is namespaced
	// under the graph-lint taxonomy, and the aggregate error names lint.
	if graphlint.AggregateCode != "graph-lint-failed" {
		t.Errorf("aggregate code = %q; want %q, which names lint",
			graphlint.AggregateCode, "graph-lint-failed")
	}
	// The rejected alternatives — "validate", "resolve --lint",
	// "pre-commit check" — must not appear as a code namespace.
	for _, c := range append(graphlint.BlockingCodes(), graphlint.AdvisoryCodes()...) {
		if !strings.HasPrefix(c, "graph-") {
			t.Errorf("finding code %q is outside the graph-lint namespace", c)
		}
		for _, rejected := range []string{"validate", "verify", "precommit", "pre-commit"} {
			if strings.Contains(c, rejected) {
				t.Errorf("finding code %q carries the rejected name %q; "+
					"the canonical subsystem name is \"lint\"", c, rejected)
			}
		}
	}
}

// REQ-15: The minimum input contract is: "model identity and version";
// "normalized candidate rows with deterministic row identity, source rule
// id, optional source span, match predicates, guard predicates …, writes,
// clears, and — for escape rows — row kind and the declared list of
// failure classes"; declared tags with provenance; "recognized outcome
// alphabet, declared initial owned state, and declared terminal states"
// HAPPY PATH
func TestReq15_MinimumInputContractIsAvailableToTheEngine(t *testing.T) {
	m := mustLoad(t, source(twoStateDecls, escapeBody))

	// The engine receives the model, so every named field must reach it
	// through the request. A field the normalized model does not carry
	// cannot be checked by any invariant that needs it.
	req := graphlint.NewRequest(m)
	got := req.Model
	if got == nil {
		t.Fatal("NewRequest carried no model; the engine has no input")
	}

	if got.ID == "" {
		t.Error("input contract: model identity is absent")
	}
	if got.Version == 0 {
		t.Error("input contract: model version is absent")
	}
	if len(got.Outcomes) == 0 {
		t.Error("input contract: the recognized outcome alphabet is absent")
	}
	if len(got.Initial) == 0 {
		t.Error("input contract: the declared initial owned state is absent")
	}
	if len(got.Terminal) == 0 {
		t.Error("input contract: the declared terminal states are absent")
	}
	if len(got.Tags) == 0 {
		t.Error("input contract: declared tags with provenance are absent")
	}

	var sawEscape, sawGuardAtom, sawWrite bool
	for _, row := range got.Rows {
		if row.Identity() == "" {
			t.Errorf("input contract: row %q carries no deterministic row identity", row.RuleID)
		}
		if row.RuleID == "" {
			t.Errorf("input contract: a row carries no source rule id")
		}
		if len(row.Escape) > 0 {
			sawEscape = true
			if row.Kind() != table.KindEscape {
				t.Errorf("input contract: row %q declares escape classes %v "+
					"but its derived kind is %q", row.RuleID, row.Escape, row.Kind())
			}
		}
		for _, a := range row.Atoms {
			if a.Block != table.BlockMatch {
				sawGuardAtom = true
			}
			if a.Key == "" || a.Operator == "" {
				t.Errorf("input contract: row %q carries an atom missing "+
					"key/operator: %+v", row.RuleID, a)
			}
		}
		if len(row.Writes) > 0 {
			sawWrite = true
		}
	}
	if !sawEscape {
		t.Error("input contract fixture carries no escape row, so row kind " +
			"and the declared failure-class list are untested")
	}
	if !sawGuardAtom {
		t.Error("input contract fixture carries no guard atom")
	}
	if !sawWrite {
		t.Error("input contract fixture carries no write")
	}
}

// REQ-16: "Lint reads every one of these from the declaration and never
// infers one from a tag's name, value spelling, or a fixture"
// ADVERSARIAL
func TestReq16_DeclaredPropertiesAreReadNeverInferred(t *testing.T) {
	// Two tags whose NAMES and value spellings suggest finiteness and
	// single-valuedness, but whose declarations say otherwise. A lint that
	// infers from the name or the spelling would prove coverage over
	// `status_enum`; a lint that reads the declaration must withhold.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[tags.status_enum_bool_domain]
provenance = "owned"
kind = "scalar"
required = true
`
	const body = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "infers"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.status_enum_bool_domain]
eq = "true"
[rule.write]
status = "b"
`
	r := lint(t, decls, body)

	// The tag is `scalar` with no declared domain, so the dimension is not
	// finite whatever its name spells. Certifying coverage here is the
	// inference the clause forbids.
	f := requireCode(t, r, graphlint.CodeUnprovableCoverage)
	var sawDimensionReason bool
	for _, got := range f {
		if got.Reason == graphlint.ReasonDimensionNotFinite {
			sawDimensionReason = true
		}
	}
	if !sawDimensionReason {
		t.Errorf("no %s finding carries reason %q; lint inferred a finite "+
			"domain from the tag's name or value spelling; report:%s",
			graphlint.CodeUnprovableCoverage,
			graphlint.ReasonDimensionNotFinite, render(r))
	}
}

// REQ-17: The pipeline has four stages: "load and normalize the transition
// model through the RDR 0002 table contract"; "derive a graph view from
// normalized candidate rows"; "run invariant checks over that graph and
// collect typed findings"; "emit either a success report or one aggregate
// structured `CLIError` failure through `respond`."
// HAPPY PATH
func TestReq17_PipelineRunsNormalizeDeriveCheckEmit(t *testing.T) {
	// Stage 1 is RDR 0002's: a non-conforming source is refused BEFORE
	// normalization and lint never runs.
	if _, err := table.Load([]byte("not toml at all ][）"), "fixture.toml"); err == nil {
		t.Error("stage 1: a non-conforming source loaded clean")
	}

	// Stages 2 and 3: the engine derives a graph view and collects TYPED
	// findings over it. The graph view is observable through the
	// reachability relation the derivation produces.
	m := mustLoad(t, source(twoStateDecls, legalBody))
	nodes := graphlint.Reach(m)
	if len(nodes) == 0 {
		t.Errorf("stage 2: the derivation produced no reachable nodes for a "+
			"model with a declared root %v", m.Initial)
	}

	r := graphlint.Run(graphlint.NewRequest(m))
	if r.ModelID == "" {
		t.Error("stage 3: the run collected no report identity")
	}
	// Stage 4's success/failure disposition is the Blocking() partition;
	// the respond wiring is asserted in the CLI suite.
	if len(r.Blocking())+len(r.Advisory()) != len(r.Findings) {
		t.Errorf("stage 3: the findings partition is not exhaustive: "+
			"%d blocking + %d advisory != %d total; report:%s",
			len(r.Blocking()), len(r.Advisory()), len(r.Findings), render(r))
	}
}

// REQ-18: "The read-set is computed by lint from the row's atoms, never
// read from `Row.RequiresOwned`." "Lint therefore derives the owned
// read-set as the owned-provenance keys named by the row's match pattern
// and its `all`/`unless` guard atoms"
// ADVERSARIAL
func TestReq18_OwnedReadSetIsDerivedFromAtomsNotRequiresOwned(t *testing.T) {
	// The fixture separates the two fields: `gate_key` is READ by a guard
	// atom and never written, while `status` is WRITTEN and so is the only
	// key RequiresOwned names. A lint reading RequiresOwned checks
	// writes-before-writes and never fires; a lint deriving the read-set
	// from the atoms reports gate_key.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[tags.gate_key]
provenance = "owned"
kind = "enum"
domain = ["open", "shut"]
single_valued = true
`
	const body = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "reads-unwritten"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.gate_key]
eq = "open"
[rule.write]
status = "b"
`
	m := mustLoad(t, source(decls, body))

	// Precondition on the fixture: RequiresOwned must NOT name gate_key,
	// or the test could not tell the two derivations apart.
	row := rowByRuleID(t, m, "reads-unwritten")
	if slices.Contains(row.RequiresOwned, "gate_key") {
		t.Fatalf("fixture precondition failed: RequiresOwned = %v already "+
			"names the guard-read key, so this test cannot distinguish the "+
			"atom-derived read-set from the RequiresOwned one", row.RequiresOwned)
	}

	r := graphlint.Run(graphlint.NewRequest(m))
	if !namesRule(r, graphlint.CodeOwnedBeforeWrite, "reads-unwritten") {
		t.Errorf("no %s finding names rule %q; lint read the write-derived "+
			"RequiresOwned (%v) instead of deriving the read-set from the "+
			"row's atoms; report:%s", graphlint.CodeOwnedBeforeWrite,
			"reads-unwritten", row.RequiresOwned, render(r))
	}
}

// rowByRuleID returns the single normalized row carrying ruleID.
func rowByRuleID(t *testing.T, m *table.Model, ruleID string) table.Row {
	t.Helper()

	for _, r := range m.Rows {
		if r.RuleID == ruleID {
			return r
		}
	}
	t.Fatalf("model carries no row with rule id %q", ruleID)
	return table.Row{}
}

// RDR 0003's req-list ASSUMPTION: "this RDR ships no lint finding *codes* of
// its own — it states the semantics and RDR 0006 mints the codes". This RDR
// owns the taxonomy, so the check belongs here: every code `internal/guard`
// emits MUST be a member of one of this RDR's two tiers. A guard-side code
// outside the taxonomy is one RDR 0003 minted in a space it disclaimed, and
// a consumer reading the machine-readable output would receive an unknown
// code. `graph-vacuous-exists` was exactly that (kata `fyf4`).
func TestGuardEmitsOnlyCodesThisTaxonomyMints(t *testing.T) {
	minted := map[string]bool{}
	for _, c := range graphlint.BlockingCodes() {
		minted[c] = true
	}
	for _, c := range graphlint.AdvisoryCodes() {
		minted[c] = true
	}

	for _, c := range guard.Codes() {
		if !minted[string(c)] {
			t.Errorf("`internal/guard` emits %q, which this RDR's taxonomy "+
				"does not mint; RDR 0003 ships no codes of its own, so every "+
				"code it emits must be a member of the blocking or advisory "+
				"tier", c)
		}
	}
}

// C9's closure attribution has TWO implementations — this engine's
// `emitCoverageArms` and `internal/guard`'s `coverageFindings` — over one
// clause. This RDR is authoritative for the mechanics, so a guard-side
// reading that diverges is a defect even where `guard.Lint` has no
// non-test caller: guard's suite is the oracle this engine's behavior is
// checked against, and an oracle that disagrees with the authority silently
// stops being one. The divergence this pins was a guard-side gate
// suppressing the attribution whenever the ordinary rows also closed the
// product — invisible to every single-engine test, because each engine was
// self-consistent (kata `dhfj`).
func TestClosureAttributionAgreesAcrossBothEngines(t *testing.T) {
	// The shape that split them: two ordinary rows PARTITIONING the guard
	// dimension, so nothing is uncovered, plus a bare escape row. It is the
	// ordinary "authored rows cover the domain, plus a defensive catch-all"
	// idiom, and the shipped example models carry it.
	const body = `
terminal = ["done"]

[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"

[[rule]]
id = "advance-off"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "false"
[rule.write]
status = "b"

[[rule]]
id = "bare-rescue"
escape = ["no_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
`
	src := source(twoStateDecls, body)
	m := mustLoad(t, src)

	// This engine, the authority.
	var authoritative []string
	for _, f := range graphlint.Run(graphlint.NewRequest(m)).Findings {
		if f.Code == graphlint.CodeCoverageClosedByEscape {
			authoritative = append(authoritative, f.Rule)
		}
	}
	slices.Sort(authoritative)

	// `internal/guard`, over the SAME normalized model — not a parallel
	// fixture, so no fixture drift can explain a disagreement away.
	var mirrored []string
	for _, r := range guard.Lint(m) {
		for _, f := range r.Findings {
			if f.Code == guard.CodeCoverageClosedByEscape {
				mirrored = append(mirrored, f.RuleIDs...)
			}
		}
		// The struct field and the finding are one claim, not two.
		if r.ClosedByEscape != "" && !slices.Contains(mirrored, r.ClosedByEscape) {
			t.Errorf("guard's group %q reports ClosedByEscape=%q with no "+
				"matching %s finding; the field and the finding are one "+
				"claim", r.Context, r.ClosedByEscape,
				guard.CodeCoverageClosedByEscape)
		}
	}
	slices.Sort(mirrored)

	if len(authoritative) == 0 {
		t.Fatalf("the authoritative engine attributed no closure over a "+
			"model carrying a bare escape row; the parity question does "+
			"not arise. src:\n%s", src)
	}
	if !slices.Equal(authoritative, mirrored) {
		t.Errorf("the two implementations of C9's closure attribution "+
			"disagree over one model: this engine names %v, "+
			"`internal/guard` names %v. This RDR is authoritative for the "+
			"mechanics, so guard follows it", authoritative, mirrored)
	}
}
