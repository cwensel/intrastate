package guard_test

// RDR 0003 — failure modes, non-goals, naming, wire ownership, and the
// implementation-phase obligations that are checkable as code.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/guard"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// REQ-112: "Visible failures should be typed load or lint failures: unknown
// operator, operator/tag-kind mismatch, literal parse failure, unknown tag,
// and the two declaration/literal-domain rejections (RDR 0002 load
// categories, §D7(iii)) at load; non-exhaustive finite domain, overlapping
// candidate rows, guard dimension not provable because it lacks a finite
// domain, or finite scoped product too large to prove at lint."
// BOUNDARY
func TestReq112_EveryVisibleFailureIsATypedLoadOrLintFailure(t *testing.T) {
	decls := `
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`
	loadCases := map[string]struct {
		guard string
		want  table.Category
	}{
		"unknown operator":           {"[rule.guard.all.profile]\nneq = \"small\"\n", table.CatMalformedPredicateAtom},
		"operator/tag-kind mismatch": {"[rule.guard.all.profile]\ncontains = [\"small\"]\n", table.CatMalformedPredicateAtom},
		"literal parse failure":      {"[rule.guard.all.profile]\nexists = \"maybe\"\n", table.CatMalformedPredicateAtom},
		"unknown tag":                {"[rule.guard.all.nowhere]\neq = \"x\"\n", table.CatUnknownTag},
		"literal outside domain":     {"[rule.guard.all.profile]\neq = \"huge\"\n", table.CatMalformedPredicateAtom},
	}
	for name, tc := range loadCases {
		t.Run("load/"+name, func(t *testing.T) {
			if got := loadCategory(t, decls, tc.guard); got != tc.want {
				t.Errorf("refused as %q; want the typed load failure %q", got, tc.want)
			}
		})
	}
	// The sixth load rejection: a malformed declaration.
	if got := loadDeclCategory(t, "[tags.t]\nprovenance = \"owned\"\nkind = \"scalar\"\ndomain = [\"a\"]\n"); got != table.CatMalformedTagDeclaration {
		t.Errorf("a malformed declaration refused as %q; want %q", got, table.CatMalformedTagDeclaration)
	}

	lintCases := map[string]struct {
		src  string
		code guard.Code
	}{
		"non-exhaustive finite domain": {productOnlyGapSource(), guard.CodeCoverageGap},
		"overlapping candidate rows":   {twoRuleIdenticalGuardSource(), guard.CodeOverlap},
		"dimension lacking a domain":   {unboundedIntGuardSource(), guard.CodeUnprovableCoverage},
		"product too large":            {overLargeProductSource(), guard.CodeProductTooLarge},
	}
	for name, tc := range lintCases {
		t.Run("lint/"+name, func(t *testing.T) {
			if countCode(guard.Lint(mustLoadSource(t, tc.src)), tc.code) == 0 {
				t.Errorf("produced no %q finding; the failure must be visible "+
					"and typed", tc.code)
			}
		})
	}
}

// REQ-113: "A guard that cannot be decided at runtime surfaces as RDR
// 0007's `guard_unevaluable` refusal, not as a predicate parse kind owned
// here."
// BOUNDARY
func TestReq113_UndecidableAtRuntimeSurfacesAsGuardUnevaluableNotAParseKind(t *testing.T) {
	m := mustLoadSource(t, optionalKeyGroupSource(false))
	res := resolveWith(t, m, m.KernelTable(), guard.View{})

	if !res.Refused() || res.Refusal.Kind != resolve.KindGuardUnevaluable {
		t.Fatalf("an undecidable guard surfaced as %v; want RDR 0007's "+
			"guard_unevaluable refusal", describe(res))
	}
	// It is NOT one of this RDR's predicate semantic kinds: the model loads
	// clean, so no parse kind was minted for it.
	for _, k := range guard.SemanticKinds() {
		if string(k) == string(resolve.KindGuardUnevaluable) {
			t.Errorf("this RDR mints predicate semantic kind %q, duplicating "+
				"RDR 0007's runtime refusal", k)
		}
	}

	// Discriminating leg: this RDR DOES own predicate semantic kinds, so
	// "not one of them" is a boundary rather than an empty set.
	kinds := guard.SemanticKinds()
	if len(kinds) == 0 {
		t.Fatal("this RDR names no predicate semantic kinds at all; the " +
			"clause draws a boundary between two populated sets")
	}
	if !slices.Contains(kinds, guard.SemanticKindLiteralParseFailure) {
		t.Errorf("SemanticKinds() = %v; a literal parse failure IS a "+
			"predicate semantic kind owned here", kinds)
	}
}

// REQ-114: "Silent failure would be a false exhaustiveness claim; the
// recovery path is to keep every exactness claim tied to A2 and the MVV
// fixture."
// ADVERSARIAL
func TestReq114_EveryGreenClaimIsTiedToTheCoverageIdentity(t *testing.T) {
	for _, src := range []string{
		completePartitionSource(),
		existsPartitionSource(),
		optionalKeyGroupSource(true),
	} {
		m := mustLoadSource(t, src)
		// A green claim must be REACHED for the identity to bind it. A lint
		// that never certifies satisfies "every green claim is tied" with an
		// empty quantifier, which is the silent failure this clause names.
		if !hasGreen(guard.Lint(m)) {
			t.Fatalf("no green claim was reached for a group that proves; the "+
				"recovery path is to keep every exactness claim TIED to A2, "+
				"not to make none. reports=%s", renderReports(guard.Lint(m)))
		}
		for _, r := range greenGroups(guard.Lint(m)) {
			g := groupOf(t, m, r.RuleIDs[0])
			union := guard.CoverageUnion(m, g)
			product := guard.Product(m, g)
			if !union.Equal(product) {
				t.Errorf("group %s certified green while its union (%d) does "+
					"not equal its scoped product (%d); every exactness claim "+
					"is tied to the A2 identity", r.Context, union.Len(), product.Len())
			}
		}
	}
}

// REQ-115: "This RDR introduces no encode/decode pair. Parse/render
// fidelity for the sparse TOML source belongs to RDR 0002." — **no
// Round-Trip / Inverse Invariant is declared.**
// BOUNDARY
func TestReq115_NoEncodeDecodePairIsIntroduced(t *testing.T) {
	for _, name := range exportedNames(t, guardPackageDir(t)) {
		lower := strings.ToLower(name)
		for _, forbidden := range []string{"marshal", "unmarshal", "encode", "decode", "parsetoml", "render"} {
			if strings.Contains(lower, forbidden) {
				t.Errorf("this RDR exports %q; it introduces no encode/decode "+
					"pair — parse/render fidelity belongs to RDR 0002", name)
			}
		}
	}

	// No round-trip fidelity requirement is carried, but the surface this
	// RDR DOES own must exist — otherwise "no encode/decode pair" is
	// satisfied by owning nothing at all.
	if len(guard.Operators()) == 0 {
		t.Error("the package exports no operator vocabulary; the absence of " +
			"an encode/decode pair is not the absence of a surface")
	}
	m := mustLoadSource(t, completePartitionSource())
	if len(guard.Lint(m)) == 0 {
		t.Error("the package produces no lint result; this RDR owns semantics, " +
			"not a codec, and the semantics must be present")
	}
}

// REQ-116: "There is no embedded host predicate and no free-form expression
// grammar."
// ADVERSARIAL
func TestReq116_NoHostPredicateAndNoFreeFormExpressionGrammar(t *testing.T) {
	// No exported surface takes or returns a predicate function: a host
	// callback would have to travel as one.
	for name, kind := range exportedKinds(t, guardPackageDir(t)) {
		if kind == "func-typed-field" {
			t.Errorf("exported surface %q carries a func-typed predicate; "+
				"there is no embedded host predicate", name)
		}
	}

	// And no free-form expression is admitted: the operator vocabulary is
	// closed, so an expression string is rejected as an operator.
	for _, expr := range []string{"a && b", "profile != small", "len(x) > 0"} {
		if guard.KnownOperator(expr) {
			t.Errorf("%q is admitted as an operator; there is no free-form "+
				"expression grammar", expr)
		}
	}

	// Discriminating leg: a KnownOperator that rejects everything would pass
	// the loop above while admitting no symbolic vocabulary either. The
	// closed set must actually be admitted.
	for _, op := range []string{"eq", "in", "contains", "exists", "lt", "lte", "gt", "gte"} {
		if !guard.KnownOperator(op) {
			t.Errorf("%q is not admitted as an operator; the vocabulary is "+
				"closed, not empty", op)
		}
	}
}

// REQ-117: "No new third-party dependency is proposed. The predicate
// grammar and finite domain checks should be implemented with local Go code
// unless Resolve proves a small parsing or set library is necessary."
// BOUNDARY
func TestReq117_NoNewThirdPartyDependencyIsIntroduced(t *testing.T) {
	dir := guardPackageDir(t)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}
	fset := token.NewFileSet()
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(dir, e.Name()), nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", e.Name(), err)
		}
		for _, imp := range f.Imports {
			path := strings.Trim(imp.Path.Value, `"`)
			if strings.HasPrefix(path, "github.com/cwensel/intrastate/") {
				continue
			}
			if strings.Contains(strings.SplitN(path, "/", 2)[0], ".") {
				t.Errorf("%s imports third-party package %q; the predicate "+
					"grammar and domain checks are local Go code", e.Name(), path)
			}
		}
	}

	// Discriminating leg: the grammar and the finite-domain checks must
	// actually be implemented here, or "no new dependency" is satisfied by
	// implementing nothing.
	if len(guard.Operators()) == 0 || len(guard.Kinds()) == 0 {
		t.Error("the predicate grammar is not implemented locally")
	}
	if _, ok := guard.AssignmentCount(table.TagDecl{
		Kind: "bool", SingleValued: true, Required: true,
	}); !ok {
		t.Error("the finite domain checks are not implemented locally")
	}
}

// REQ-118: "This RDR creates no persistent resource."
// BOUNDARY
func TestReq118_NoPersistentResourceIsCreated(t *testing.T) {
	dir := guardPackageDir(t)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}
	fset := token.NewFileSet()
	banned := []string{"os", "net", "net/http", "database/sql", "path/filepath", "io/fs"}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(dir, e.Name()), nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", e.Name(), err)
		}
		for _, imp := range f.Imports {
			path := strings.Trim(imp.Path.Value, `"`)
			if slices.Contains(banned, path) {
				t.Errorf("%s imports %q; this RDR creates no persistent "+
					"resource and performs no I/O", e.Name(), path)
			}
		}
	}

	// Discriminating leg: lint runs to completion over an in-memory model
	// with no resource of any kind, which is what "creates no persistent
	// resource" means operationally rather than merely as an import ban.
	m := mustLoadSource(t, completePartitionSource())
	if len(guard.Lint(m)) == 0 {
		t.Error("lint produced nothing over an in-memory model; the absence " +
			"of a persistent resource is not the absence of a result")
	}
}

// REQ-119: "no callback invocation, expression parser, or external engine
// is part of the hot path. If implementation later indexes predicates for
// speed, the optimization must preserve the normalized atom semantics,
// scoped product proof, and exact-one refusal behavior."
// BOUNDARY
func TestReq119_TheHotPathCarriesNoCallbackParserOrExternalEngine(t *testing.T) {
	// The runtime hot path is the seam: one atom plus one present value in,
	// one verdict out. Repeated calls are pure — the same inputs give the
	// same verdict, which an indexed or stateful engine would have to
	// preserve anyway.
	ev := guard.NewEvaluator(map[string]string{"profile": "enum"})
	atom := kernelAtom("profile", "in", `["mid","large"]`)
	first := ev.Evaluate(atom, "large")
	for range 100 {
		if got := ev.Evaluate(atom, "large"); got != first {
			t.Fatalf("Evaluate is not a pure function of (atom, value): %v then %v",
				first, got)
		}
	}
	if first != resolve.GuardTrue {
		t.Errorf("Evaluate(%+v, %q) = %v; want GuardTrue", atom, "large", first)
	}

	// Exact-one refusal behaviour is preserved end to end.
	dup := mustLoadSource(t, twoRuleIdenticalGuardSource())
	res := resolveWith(t, dup, dup.KernelTable(), guard.View{"profile": "large"})
	if !res.Refused() || res.Refusal.Kind != resolve.KindAmbiguousMatch {
		t.Errorf("exact-one refusal behaviour was not preserved: %v", describe(res))
	}
}

// REQ-120: "the canonical name is \"guard predicate\"; rejected alternatives
// are \"condition callback\" and \"guard expression\" because both invite
// opaque host logic."
// BOUNDARY
func TestReq120_TheCanonicalNameIsGuardPredicate(t *testing.T) {
	for _, name := range exportedNames(t, guardPackageDir(t)) {
		lower := strings.ToLower(name)
		for _, rejected := range []string{"conditioncallback", "guardexpression", "expression", "callback"} {
			if strings.Contains(lower, rejected) {
				t.Errorf("exported name %q carries the rejected vocabulary %q; "+
					"the canonical name is \"guard predicate\"", name, rejected)
			}
		}
	}

	// Discriminating leg: the canonical name must be REACHED, not merely
	// unrejected. A package exporting nothing carries neither vocabulary.
	names := exportedNames(t, guardPackageDir(t))
	var carriesPredicateVocabulary bool
	for _, name := range names {
		if strings.Contains(strings.ToLower(name), "predicate") ||
			strings.Contains(strings.ToLower(name), "guard") ||
			strings.Contains(strings.ToLower(name), "atom") {
			carriesPredicateVocabulary = true
		}
	}
	if !carriesPredicateVocabulary {
		t.Errorf("no exported name carries the canonical \"guard predicate\" "+
			"vocabulary; names=%v", names)
	}

	// The name is load-bearing because of what it rules out: a "condition
	// callback" or "guard expression" would carry opaque host logic, and the
	// thing this package actually names must therefore be a SYMBOLIC atom
	// that decides from its own declared fields.
	ev := guard.NewEvaluator(map[string]string{"subject": "enum"})
	atom := kernelAtom("subject", "eq", "Draft")
	if got := ev.Evaluate(atom, "Draft"); got != resolve.GuardTrue {
		t.Errorf("Evaluate(%+v, %q) = %v; the canonically-named guard "+
			"predicate decides from its own symbolic fields, which is the "+
			"property the rejected names would have surrendered",
			atom, "Draft", got)
	}
}

// REQ-121: "RDR 0002 owns the TOML container; this RDR owns the guard atom
// grammar embedded in that container."
// BOUNDARY
func TestReq121_TheContainerIsRDR0002sAndTheAtomGrammarIsThisRDRs(t *testing.T) {
	// The container: this package never reads TOML. It takes an
	// already-loaded model.
	dir := guardPackageDir(t)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}
	fset := token.NewFileSet()
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(dir, e.Name()), nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", e.Name(), err)
		}
		for _, imp := range f.Imports {
			if strings.Contains(imp.Path.Value, "toml") {
				t.Errorf("%s imports a TOML package; RDR 0002 owns the "+
					"container", e.Name())
			}
		}
	}

	// The atom grammar: this RDR owns the closed operator vocabulary and the
	// operator/kind matrix.
	if len(guard.Operators()) == 0 || len(guard.Kinds()) == 0 {
		t.Error("this RDR publishes no operator vocabulary or kind set; it " +
			"owns the guard atom grammar")
	}
}

// REQ-122: "Define the tag-kind/operator compatibility matrix and normalized
// predicate atom shape used by resolver and lint."
// BOUNDARY
func TestReq122_TheMatrixAndAtomShapeAreSharedByResolverAndLint(t *testing.T) {
	// The matrix exists and is total over the published vocabulary.
	for _, op := range guard.Operators() {
		for _, kind := range guard.Kinds() {
			_ = guard.Accepts(op, kind)
		}
	}

	// One atom shape serves both surfaces: the row's normalized atom drives
	// the lint denotation, and the SAME atom crosses to the resolver seam.
	m := mustLoadSource(t, completePartitionSource())
	row := rowByID(t, m, "part-small")
	if !guard.Denotation(m, "profile", row.Atoms[0]).Projectable() {
		t.Error("lint could not project the normalized atom the loader emitted")
	}
	kernelRow := row.KernelRow()
	if len(kernelRow.Guard) == 0 {
		t.Fatal("the normalized row carried no guard atom to the resolver")
	}
	if kernelRow.Guard[0].Key != row.Atoms[0].Key ||
		kernelRow.Guard[0].Operator != row.Atoms[0].Operator {
		t.Errorf("the resolver atom %+v does not carry the lint atom %+v; one "+
			"normalized shape serves both", kernelRow.Guard[0], row.Atoms[0])
	}
}

// REQ-123: "Define how enum, boolean, set-universe, and bounded-int domains
// are converted into scoped row-group coverage and overlap checks, including
// refusal/downgrade behavior for unbounded dimensions and finite products
// too large to prove deterministically."
// HAPPY PATH
func TestReq123_AllFourFiniteKindsReachCoverageAndOverlapChecks(t *testing.T) {
	min0, max1 := 0, 1
	kinds := map[string]table.TagDecl{
		"enum":         {Kind: "enum", Domain: []string{"a", "b"}, SingleValued: true, Required: true},
		"bool":         {Kind: "bool", SingleValued: true, Required: true},
		"set-universe": {Kind: "set", Elements: []string{"x"}, Required: true},
		"bounded-int":  {Kind: "int", Min: &min0, Max: &max1, SingleValued: true, Required: true},
	}
	for name, decl := range kinds {
		if _, ok := guard.AssignmentCount(decl); !ok {
			t.Errorf("%s does not convert into a scoped product dimension", name)
		}
	}

	// Coverage and overlap both run over the converted domains.
	m := mustLoadSource(t, gapAndOverlapSource())
	reports := guard.Lint(m)
	if countCode(reports, guard.CodeCoverageGap) == 0 || countCode(reports, guard.CodeOverlap) == 0 {
		t.Errorf("the converted domains did not reach both checks; findings=%v",
			allFindings(reports))
	}

	// Refusal/downgrade behaviour for both unprovable classes.
	if countCode(guard.Lint(mustLoadSource(t, unboundedIntGuardSource())), guard.CodeUnprovableCoverage) == 0 {
		t.Error("an unbounded dimension did not reach the refusal behaviour")
	}
	if countCode(guard.Lint(mustLoadSource(t, overLargeProductSource())), guard.CodeProductTooLarge) == 0 {
		t.Error("an over-large finite product did not reach the refusal behaviour")
	}
}

// REQ-125: "Connect predicate diagnostics to RDR 0002 source identities, RDR
// 0001 exact-one selection, RDR 0005 CLI output, and RDR 0006 lint
// authority."
// BOUNDARY
func TestReq125_DiagnosticsCarryRDR0002IdentitiesAndRDR0006Codes(t *testing.T) {
	m := mustLoadSource(t, gapAndOverlapSource())
	reports := guard.Lint(m)

	// RDR 0002 source identities travel on every finding.
	for _, f := range allFindings(reports) {
		if len(f.RuleIDs) == 0 {
			t.Errorf("finding %v carries no RDR 0002 rule id", f)
		}
		for _, loc := range f.Locators {
			if loc == "" {
				t.Errorf("finding %v carries an empty source locator", f)
			}
		}
	}

	// RDR 0006 owns the finding codes; this RDR emits them by their names.
	for _, code := range guard.Codes() {
		if !strings.HasPrefix(string(code), "graph-") {
			t.Errorf("code %q is not one of RDR 0006's `graph-` codes; this "+
				"RDR mints none of its own", code)
		}
	}

	// RDR 0001 exact-one selection: the overlap the diagnostic names is the
	// ambiguity the kernel refuses.
	res := resolveWith(t, m, m.KernelTable(), guard.View{"profile": "small", "flag": "true"})
	if !res.Refused() || res.Refusal.Kind != resolve.KindAmbiguousMatch {
		t.Errorf("the kernel gave %v for the assignment lint reports as "+
			"overlapping", describe(res))
	}
}

// REQ-141: "The MVV should become production tests that exercise both
// runtime predicate evaluation and lint-time finite-domain reasoning. Done
// means the same normalized predicate atoms drive exact-one row selection,
// overlap detection, and exhaustiveness proof/refusal without host callbacks
// or source-order priority."
// HAPPY PATH
func TestReq141_OneAtomSetDrivesSelectionOverlapAndProofWithoutOrderPriority(t *testing.T) {
	m := mustLoadSource(t, gapAndOverlapSource())

	// Exact-one selection over the same atoms.
	sel := resolveWith(t, m, m.KernelTable(), guard.View{"profile": "large", "flag": "true"})
	if sel.Refused() {
		t.Errorf("exact-one selection over the normalized atoms refused %v",
			sel.Refusal.Kind)
	}

	// Overlap detection and exhaustiveness refusal over the SAME atoms.
	reports := guard.Lint(m)
	if countCode(reports, guard.CodeOverlap) == 0 || countCode(reports, guard.CodeCoverageGap) == 0 {
		t.Errorf("the same atoms did not drive overlap detection and the "+
			"exhaustiveness verdict; findings=%v", allFindings(reports))
	}

	// Without source-order priority: reordering changes nothing.
	kt := m.KernelTable()
	slices.Reverse(kt.Rows)
	reordered := resolveWith(t, m, kt, guard.View{"profile": "large", "flag": "true"})
	if !sameDisposition(sel, reordered) {
		t.Errorf("reordering rows changed the disposition (%v vs %v)",
			describe(sel), describe(reordered))
	}
}

// --- source-inspection helpers -------------------------------------------

// guardPackageDir locates this package's non-test source directory.
func guardPackageDir(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return wd
}

// exportedNames lists the exported identifiers the package declares.
func exportedNames(t *testing.T, dir string) []string {
	t.Helper()

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", dir, err)
	}
	var out []string
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			for name, obj := range f.Scope.Objects {
				if ast.IsExported(name) {
					out = append(out, name)
					_ = obj
				}
			}
		}
	}
	if len(out) == 0 {
		t.Fatalf("package in %s declares no exported identifiers", dir)
	}
	slices.Sort(out)
	return out
}

// exportedKinds classifies exported struct fields that carry a func type,
// which is how a host callback would have to travel.
func exportedKinds(t *testing.T, dir string) map[string]string {
	t.Helper()

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", dir, err)
	}
	out := map[string]string{}
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				ts, ok := n.(*ast.TypeSpec)
				if !ok || !ast.IsExported(ts.Name.Name) {
					return true
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					return true
				}
				for _, field := range st.Fields.List {
					if _, isFunc := field.Type.(*ast.FuncType); !isFunc {
						continue
					}
					for _, name := range field.Names {
						if ast.IsExported(name.Name) {
							out[ts.Name.Name+"."+name.Name] = "func-typed-field"
						}
					}
				}
				return true
			})
		}
	}
	return out
}
