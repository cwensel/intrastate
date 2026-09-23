package cli

// RDR 0012 — the declared-kind carrier, end to end through the CLI.
//
// The seam-level matrix is pinned in `internal/guard`. What this file owns
// is the INGRESS: typed comparison reached through the real command tree —
// the owned door (`flow set-state` → artifact → reader → `OwnedSnapshot` →
// `resolve::assemble` → seam), the `--tag` door, `flow next`'s probe, and
// the model loader's canonical-int refusal at `lint --model`. A seam-level
// unit test would satisfy these assertions while testing none of the
// ingress, which is the whole claim (`0012:MVV`).
//
// Every oracle names an exact refusal code and payload field, so an
// unrelated failure cannot pass for the contracted one. Where an assertion
// needs the code shape of an unwritten construction site, it reads the
// SOURCE through `go/parser`: naming a symbol the tree lacks would break the
// build the pre-commit `go vet ./...` gates.

import (
	"bufio"
	"bytes"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
)

// --- fixtures ---------------------------------------------------------------

// r12OwnedModel is the MVV's `mvv-int` model on the OWNED door: `iter` is
// an owned `int`, single-valued, read by the `state` accessor; row
// `guarded` is guarded `iter eq 7` and row `fallback` is unguarded on the
// same outcome. The owned ingress is the venue because it is the only one
// that reaches the seam unconformed.
const r12OwnedModel = `outcomes = ["step"]
terminal = ["done"]

[model]
id = "mvvint"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true
required = true

[tags.iter]
provenance = "owned"
kind = "int"
min = 0
max = 9
single_valued = true

[read.state]
role = "state"
path = "flow.state"
keys = ["status", "iter"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status", "iter"]
timeout = "2s"
read_back = true

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "guarded"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "step"
[rule.guard.all.iter]
eq = "7"
[rule.write]
status = "final"

[[rule]]
id = "fallback"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "step"
[rule.write]
status = "final"
`

// r12SeedModel is r12OwnedModel with `iter` declared `scalar`. It exists
// only to SEED the artifact: the CLI's own write path persists whatever
// value the test names (`many`, `07`, `4`), and the assertions all run
// against r12OwnedModel, whose reader then returns that value unconformed.
// No test here claims to know the artifact's on-disk format.
var r12SeedModel = strings.Replace(
	strings.Replace(r12OwnedModel, "kind = \"int\"\nmin = 0\nmax = 9\nsingle_valued = true\n",
		"kind = \"scalar\"\n", 1),
	"id = \"mvvint\"", "id = \"mvvseed\"", 1)

// r12ObservedModel is the S8 fixture: the same `mvv-int` guard, with `iter`
// an OBSERVED int supplied on the `--tag` door.
const r12ObservedModel = `outcomes = ["step"]

[model]
id = "mvv-int"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.iter]
provenance = "observed"
kind = "int"
min = 0
max = 9
single_valued = true
required = true

[[rule]]
id = "guarded"
[rule.match.recognized]
eq = "step"
[rule.guard.all.iter]
eq = 7
[rule.emit]
route = "guarded"

[[rule]]
id = "fallback"
[rule.match.recognized]
eq = "step"
[rule.emit]
route = "fallback"
`

// r12SplitModel carries ONE int dimension `n` read by both deciders: the
// `g` outcome discriminates it with GUARD atoms, the `m` outcome with MATCH
// atoms. It is the witness for the record's one asymmetry.
const r12SplitModel = `outcomes = ["g", "m"]

[model]
id = "split"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.n]
provenance = "observed"
kind = "int"
min = 0
max = 1
single_valued = true
required = true

[[rule]]
id = "g-zero"
[rule.match.recognized]
eq = "g"
[rule.guard.all.n]
eq = 0
[rule.emit]
verdict = "zero"

[[rule]]
id = "g-one"
[rule.match.recognized]
eq = "g"
[rule.guard.all.n]
eq = 1
[rule.emit]
verdict = "one"

[[rule]]
id = "m-zero"
[rule.match.recognized]
eq = "m"
[rule.match.n]
eq = 0
[rule.emit]
verdict = "zero"

[[rule]]
id = "m-one"
[rule.match.recognized]
eq = "m"
[rule.match.n]
eq = 1
[rule.emit]
verdict = "one"
`

// r12Cover07 is the S7 normative fixture: `[tags.n]` int min 0 max 1, row
// `r-zero` guarded `n eq "00"` emitting `zero`, row `r-one` guarded `n eq 1`
// emitting `one`.
const r12Cover07 = `outcomes = ["decide"]

[model]
id = "cover07"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.n]
provenance = "observed"
kind = "int"
min = 0
max = 1
single_valued = true
required = true

[[rule]]
id = "r-zero"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.n]
eq = "00"
[rule.emit]
verdict = "zero"

[[rule]]
id = "r-one"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.n]
eq = 1
[rule.emit]
verdict = "one"
`

// r12Cover07With replaces r-zero's authored `n` atom block.
func r12Cover07With(block, atom string) string {
	src := strings.Replace(r12Cover07, "[rule.guard.all.n]\neq = \"00\"\n", block+"\n"+atom+"\n", 1)
	return src
}

// --- helpers ----------------------------------------------------------------

// r12Envelope is a decoded `--as=json` failure line with its top-level keys.
type r12Envelope struct {
	Code     string        `json:"code"`
	Message  string        `json:"message"`
	Findings []flowFinding `json:"findings"`
	keys     []string
	raw      string
}

func r12Decode(t *testing.T, stdout string) r12Envelope {
	t.Helper()

	var env r12Envelope
	line := strings.TrimSpace(stdout)
	if err := json.Unmarshal([]byte(line), &env); err != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", err, stdout)
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &top); err != nil {
		t.Fatalf("stdout is not one JSON object: %v", err)
	}
	for k := range top {
		env.keys = append(env.keys, k)
	}
	sort.Strings(env.keys)
	env.raw = line
	return env
}

// r12Run drives args and returns the decoded failure envelope, or nil and
// the success stdout.
func r12Run(t *testing.T, args ...string) (*r12Envelope, string, error) {
	t.Helper()

	stdout, _, err := runCmd(t, append(args, "--as=json")...)
	if err == nil {
		return nil, stdout, nil
	}
	env := r12Decode(t, stdout)
	return &env, stdout, err
}

// r12OwnedBinding seeds an artifact whose `iter` holds value, through the
// CLI's own write path, and returns the `--artifact` binding.
func r12OwnedBinding(t *testing.T, value string) string {
	t.Helper()

	seed := writeFlowModel(t, r12SeedModel)
	art := seedArtifact(t, seed, "status=draft", "iter="+value)
	return artifactBinding(flowStateRole, art)
}

// r12ResolveOwned runs `flow resolve --outcome step --plan-only` on the
// owned door with `iter` holding value.
func r12ResolveOwned(t *testing.T, value string) (*r12Envelope, string, error) {
	t.Helper()

	model := writeFlowModel(t, r12OwnedModel)
	return r12Run(t, "flow", "resolve", "--model", model,
		"--artifact", r12OwnedBinding(t, value), "--outcome", "step", "--plan-only")
}

// r12NextOwned runs `flow next` on the owned door with `iter` holding value
// and returns the candidates by rule id.
func r12NextOwned(t *testing.T, value string) map[string]map[string]any {
	t.Helper()

	model := writeFlowModel(t, r12OwnedModel)
	stdout := requireSuccess(t, "flow", "next", "--model", model,
		"--artifact", r12OwnedBinding(t, value), "--as=json")
	data := flowData(t, stdout)
	cands, ok := objectsAt(data, "candidates")
	if !ok {
		t.Fatalf("flow next payload carries no `candidates` array; keys: %v", keysOf(data))
	}
	out := map[string]map[string]any{}
	for _, c := range cands {
		if id, _ := c["rule"].(string); id != "" {
			out[id] = c
		}
	}
	return out
}

// r12Unknown returns a candidate's `unknown` pairs as key → reason.
func r12Unknown(c map[string]any) map[string]string {
	out := map[string]string{}
	pairs, _ := objectsAt(c, "unknown")
	for _, p := range pairs {
		k, _ := p["key"].(string)
		r, _ := p["reason"].(string)
		out[k] = r
	}
	return out
}

// r12Rules returns the rule ids a refusal's findings name.
func r12Rules(env *r12Envelope) []string {
	var out []string
	for _, f := range env.Findings {
		if f.Rule != "" {
			out = append(out, f.Rule)
		}
	}
	sort.Strings(out)
	return out
}

// r12SelectedRule runs a plan and returns the selected rule, failing on any
// refusal.
func r12SelectedRule(t *testing.T, stdout string) string {
	t.Helper()

	data := flowData(t, stdout)
	rule, _ := data["rule"].(string)
	return rule
}

// r12Lint runs `lint --model` over src and returns the envelope (nil on
// success) and the process exit code.
func r12Lint(t *testing.T, src string) (*r12Envelope, int) {
	t.Helper()

	path := writeModel(t, src)
	env, _, err := r12Run(t, "lint", "--model", path)
	if err == nil {
		return nil, 0
	}
	return env, clierr.ExitCodeFor(err)
}

// r12CanonicalRefusal asserts the loader refused src with C5's diagnostic.
func r12CanonicalRefusal(t *testing.T, src, authored, rewrite string) *r12Envelope {
	t.Helper()

	env, exit := r12Lint(t, src)
	if env == nil {
		t.Fatalf("lint accepted a predicate literal %q over an int tag; want the canonical-"+
			"spelling refusal (write %s)", authored, rewrite)
	}
	if exit != 2 {
		t.Errorf("exit = %d; want 2", exit)
	}
	want := `"` + authored + `" is not the canonical spelling of int tag n; write ` + rewrite
	for _, f := range env.Findings {
		if f.Code == "malformed_predicate_atom" && strings.Contains(f.Message, want) {
			return env
		}
	}
	t.Fatalf("lint refused %s with no malformed_predicate_atom finding carrying %q:\n%s",
		env.Code, want, env.raw)
	return env
}

// --- source helpers ---------------------------------------------------------

func r12Root(t *testing.T) string {
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

func r12ParseFile(t *testing.T, rel string) (*ast.File, *token.FileSet) {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filepath.Join(r12Root(t), rel), nil, 0)
	if err != nil {
		t.Fatalf("parsing %s: %v", rel, err)
	}
	return file, fset
}

func r12Func(file *ast.File, name string) *ast.FuncDecl {
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv == nil && fn.Name.Name == name {
			return fn
		}
	}
	return nil
}

func r12ParamTypes(t *testing.T, fset *token.FileSet, fn *ast.FuncDecl) []string {
	t.Helper()

	var out []string
	for _, f := range fn.Type.Params.List {
		var b strings.Builder
		if err := printer.Fprint(&b, fset, f.Type); err != nil {
			t.Fatalf("rendering: %v", err)
		}
		for range max(len(f.Names), 1) {
			out = append(out, b.String())
		}
	}
	return out
}

// r12ConstructsTyped reports whether node calls NewEvaluator(DeclaredKinds(x)).
func r12ConstructsTyped(node ast.Node) bool {
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
		if inner, ok := call.Args[0].(*ast.CallExpr); ok && name(inner.Fun) == "DeclaredKinds" {
			found = true
		}
		return true
	})
	return found
}

// --- C2 / C4: typed comparison at the ingress --------------------------------

// REQ-11: "TYPED COMPARISON APPLIES ON THE GUARD PATH. `eq` and `in` GUARD
// atoms resolve the atom's key against the constructed kind mapping and
// compare under the declared kind. MATCH atoms do NOT: the kernel
// byte-compares them at `internal/resolve/resolve.go::TagSet.matches`"
// HAPPY PATH
func TestReq11_0012_GuardAtomsCompareTypedAndMatchAtomsCompareBytes(t *testing.T) {
	model := writeFlowModel(t, r12SplitModel)

	// Guard side: `n = 01` against the guard `n eq 1` is parsed 1 == 1.
	env, stdout, _ := r12Run(t, "flow", "resolve", "--model", model,
		"--outcome", "g", "--tag", "n=01", "--plan-only")
	if env != nil {
		t.Errorf("guard path: --tag n=01 refused %s; the guard `n eq 1` compares under the "+
			"declared int kind, so row g-one matches:\n%s", env.Code, env.raw)
	} else if got := r12SelectedRule(t, stdout); got != "g-one" {
		t.Errorf("guard path: selected %q; want g-one (parsed 1 == 1)", got)
	}

	// Match side: the SAME held bytes against the match `n eq 1` are
	// byte-compared, so no row matches.
	env, stdout, _ = r12Run(t, "flow", "resolve", "--model", model,
		"--outcome", "m", "--tag", "n=01", "--plan-only")
	if env == nil {
		t.Errorf("match path: --tag n=01 selected %q; MATCH atoms byte-compare, so \"01\" "+
			"does not match `n eq 1`", r12SelectedRule(t, stdout))
	} else if env.Code != "flow-no-match" {
		t.Errorf("match path: refused %s; want flow-no-match (byte comparison)", env.Code)
	}
}

// REQ-33: "EVERY non-test evaluator construction site … constructs over the
// declarations of the SAME loaded model whose rows it evaluates: one model,
// one evaluator, no mixed-model evaluation. … On `main` that set is exactly
// two: `internal/cli/flow_resolve.go::guardSeam`,
// `internal/guard/product.go::valueSatisfies` (A3)." Testable half: the CLI
// site constructs through `NewEvaluator(DeclaredKinds(m))`, witnessed by a
// typed verdict at `flow resolve` and `flow next`.
// HAPPY PATH
func TestReq33_0012_TheCLISiteConstructsOverTheRequestsModel(t *testing.T) {
	file, _ := r12ParseFile(t, "internal/cli/flow_resolve.go")
	fn := r12Func(file, "guardSeam")
	if fn == nil {
		t.Fatal("internal/cli/flow_resolve.go no longer declares guardSeam")
	}
	if !r12ConstructsTyped(fn) {
		t.Error("guardSeam does not construct guard.NewEvaluator(guard.DeclaredKinds(m)); the " +
			"CLI site builds its evaluator over the loaded model's declarations")
	}

	// `flow resolve` witness: the typed verdict on the --tag door.
	model := writeFlowModel(t, r12ObservedModel)
	env, _, _ := r12Run(t, "flow", "resolve", "--model", model,
		"--outcome", "step", "--tag", "iter=07", "--plan-only")
	if env == nil || env.Code != "flow-ambiguous-match" ||
		!slices.Contains(r12Rules(env), "guarded") {
		t.Errorf("flow resolve --tag iter=07: the guarded row does not match under the " +
			"declared int kind (want flow-ambiguous-match naming `guarded` beside `fallback`)")
	}

	// `flow next` witness: the probe seam is typed too.
	cands := r12NextOwned(t, "07")
	if _, ok := cands["guarded"]; !ok {
		t.Error("flow next with owned iter=07: row `guarded` is excluded; the probe seam " +
			"compares under the declared kind, so parsed 7 == 7 keeps it a candidate")
	}
}

// REQ-34: "At the CLI the constructed evaluator is THREADED DOWN, not
// re-derived: `guardSeam` takes the model and is called once per request,
// and `internal/cli/flow_next.go::probeRow` receives the constructed
// `resolve.GuardEvaluator` rather than widening to take a `*table.Model`."
// BOUNDARY
func TestReq34_0012_TheEvaluatorIsThreadedDownNotReDerived(t *testing.T) {
	file, fset := r12ParseFile(t, "internal/cli/flow_resolve.go")
	seam := r12Func(file, "guardSeam")
	if seam == nil {
		t.Fatal("guardSeam is gone")
	}
	if params := r12ParamTypes(t, fset, seam); !slices.Contains(params, "*table.Model") {
		t.Errorf("guardSeam takes %v; it takes the model", params)
	}

	file, fset = r12ParseFile(t, "internal/cli/flow_next.go")
	probe := r12Func(file, "probeRow")
	if probe == nil {
		t.Fatal("internal/cli/flow_next.go no longer declares probeRow")
	}
	params := r12ParamTypes(t, fset, probe)
	if !slices.Contains(params, "resolve.GuardEvaluator") {
		t.Errorf("probeRow takes %v; it receives the constructed resolve.GuardEvaluator", params)
	}
	if slices.Contains(params, "*table.Model") {
		t.Errorf("probeRow takes %v; it does not widen to take a *table.Model", params)
	}
	ast.Inspect(probe.Body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if id, ok := call.Fun.(*ast.Ident); ok && id.Name == "guardSeam" {
				t.Error("probeRow calls guardSeam; the evaluator is threaded down, not re-derived per row")
			}
		}
		return true
	})
}

// REQ-36: "Every consumer of a seam verdict handles GuardUnevaluable as its
// OWN case, distinct from both GuardTrue and GuardFalse; a `!= GuardFalse` /
// `== GuardTrue` two-valued collapse at any consumer is a defect this clause
// forbids."
// ADVERSARIAL
func TestReq36_0012_ATypedUnevaluableIsNeitherPrunedNorMatched(t *testing.T) {
	// `flow resolve`: neither collapse — not pruned to the fallback plan
	// (== GuardTrue collapse), not selected (!= GuardFalse collapse).
	env, stdout, _ := r12ResolveOwned(t, "many")
	switch {
	case env == nil:
		t.Errorf("flow resolve with owned iter=many planned %q; an unevaluable guard is its "+
			"own case — the row is neither pruned nor matched", r12SelectedRule(t, stdout))
	case env.Code != "flow-guard-unevaluable":
		t.Errorf("flow resolve with owned iter=many refused %s; want flow-guard-unevaluable",
			env.Code)
	}

	// `flow next`: the row stays a candidate with the atom reported unknown,
	// rather than being excluded (a FALSE reading) or reported clean (a
	// TRUE reading).
	cands := r12NextOwned(t, "many")
	row, ok := cands["guarded"]
	if !ok {
		t.Fatal("flow next with owned iter=many excluded row `guarded`; an unevaluable " +
			"verdict is not GuardFalse")
	}
	if _, unknown := r12Unknown(row)["iter"]; !unknown {
		t.Errorf("flow next reports row `guarded` with unknown %v; the undecidable `iter` "+
			"atom must be reported, not read as GuardTrue", r12Unknown(row))
	}
}

// REQ-37: "That split is the record's one asymmetry and it is named on
// purpose (C2 states it on the comparison side, and the user-facing docs
// carry it): one declared kind, two comparison rules, divided by whether the
// two deciders meet."
// HAPPY PATH
func TestReq37_0012_TheUserFacingDocsStateTheGuardMatchSplit(t *testing.T) {
	src, err := os.ReadFile(filepath.Join(r12Root(t), "docs", "model-authoring.md"))
	if err != nil {
		t.Fatalf("reading docs/model-authoring.md: %v", err)
	}
	for para := range strings.SplitSeq(string(src), "\n\n") {
		p := strings.ToLower(strings.Join(strings.Fields(para), " "))
		if strings.Contains(p, "guard") && strings.Contains(p, "match") &&
			strings.Contains(p, "declared kind") && strings.Contains(p, "byte") {
			return
		}
	}
	t.Error("docs/model-authoring.md has no paragraph stating the split: guard `eq`/`in` " +
		"compare under the declared kind while match atoms byte-compare")
}

// REQ-38: "No new refusal kind, error code, or envelope field for the SEAM:
// a typed `GuardUnevaluable` reaches the user through the existing
// `guard_unevaluable` refusal and its payload (`0007:C8`), reason
// `uncomparable` … — the CLI already maps it to `flow-guard-unevaluable`."
// BOUNDARY
func TestReq38_0012_ATypedUnevaluableRidesTheExistingRefusal(t *testing.T) {
	// The baseline: the same model's EXISTING guard_unevaluable envelope,
	// provoked by an absent owned `iter`.
	model := writeFlowModel(t, r12OwnedModel)
	seed := writeFlowModel(t, r12SeedModel)
	absentArt := seedArtifact(t, seed, "status=draft")
	base, _, _ := r12Run(t, "flow", "resolve", "--model", model,
		"--artifact", artifactBinding(flowStateRole, absentArt), "--outcome", "step", "--plan-only")
	if base == nil || base.Code != "flow-guard-unevaluable" {
		t.Fatalf("baseline (absent iter) did not refuse flow-guard-unevaluable: %+v", base)
	}

	env, _, err := r12ResolveOwned(t, "many")
	if env == nil {
		t.Fatal("flow resolve with owned iter=many succeeded; want flow-guard-unevaluable")
	}
	if env.Code != "flow-guard-unevaluable" {
		t.Fatalf("refusal code = %q; want flow-guard-unevaluable (no new code)", env.Code)
	}
	if got := clierr.ExitCodeFor(err); got != 2 {
		t.Errorf("exit = %d; want 2", got)
	}
	if !slices.Equal(env.keys, base.keys) {
		t.Errorf("envelope keys = %v; the existing guard_unevaluable envelope carries %v — "+
			"no new envelope field", env.keys, base.keys)
	}
	for _, f := range env.Findings {
		if f.Code != "flow-guard-unevaluable" {
			t.Errorf("finding code %q; the payload is the existing guard_unevaluable one", f.Code)
		}
		if f.Rule == "guarded" && !strings.Contains(f.Message, "(uncomparable)") {
			t.Errorf("finding message %q; the reason is `uncomparable`", f.Message)
		}
	}
}

// --- C5: canonical int spelling at the predicate ingress ---------------------

// REQ-39: "CANONICAL INT SPELLING. A value authored against an
// `int`-declared tag is admitted only in its canonical decimal spelling:
// after `strconv.Atoi` succeeds, `strconv.Itoa(n)` MUST equal the authored
// string. Non-canonical spellings (`\"00\"`, `\"01\"`, `\"+1\"`, `\"-0\"`)
// are REFUSED at the ingress that admits them. The rule's scope is the
// PREDICATE ingress and only it: guard and match atom literals"
// ADVERSARIAL
func TestReq39_0012_NonCanonicalPredicateIntLiteralsAreRefused(t *testing.T) {
	cases := []struct {
		name, block, atom, authored, rewrite string
	}{
		{"guard eq 00", "[rule.guard.all.n]", `eq = "00"`, "00", "0"},
		{"guard eq 01", "[rule.guard.all.n]", `eq = "01"`, "01", "1"},
		{"guard eq +1", "[rule.guard.all.n]", `eq = "+1"`, "+1", "1"},
		{"guard eq -0", "[rule.guard.all.n]", `eq = "-0"`, "-0", "0"},
		{"guard in member", "[rule.guard.all.n]", `in = ["0", "01"]`, "01", "1"},
		{"guard gte bound", "[rule.guard.all.n]", `gte = "01"`, "01", "1"},
		{"guard unless", "[rule.guard.unless.n]", `eq = "01"`, "01", "1"},
		{"match eq", "[rule.match.n]", `eq = "01"`, "01", "1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r12CanonicalRefusal(t, r12Cover07With(tc.block, tc.atom), tc.authored, tc.rewrite)
		})
	}
}

// REQ-40: "It also covers the bare-float path there, where
// `internal/table/load.go::valueMembers` renders a TOML `eq = -0.0` to the
// literal `\"-0\"` via `strconv.FormatFloat` before any kind check."
// INPUT EDGE
func TestReq40_0012_TheBareFloatNegativeZeroIsRefused(t *testing.T) {
	r12CanonicalRefusal(t, r12Cover07With("[rule.guard.all.n]", "eq = -0.0"), "-0", "0")
}

// REQ-41: "The canonicality check is therefore applied at the atom-building
// site … leaving `conformKind` itself unchanged." Observable:
// `[emit]`/`[initial]`/`[rule.write]` non-canonical int values stay ADMITTED.
// BOUNDARY
func TestReq41_0012_TheCheckReachesPredicateLiteralsOnly(t *testing.T) {
	// Non-canonical ints in [initial] and [rule.write] load.
	admitted := strings.Replace(r12OwnedModel, "[initial]\nstatus = \"draft\"\n",
		"[initial]\nstatus = \"draft\"\niter = \"07\"\n", 1)
	admitted = strings.Replace(admitted, "id = \"fallback\"\n[rule.match.status]\neq = \"draft\"\n"+
		"[rule.match.recognized]\neq = \"step\"\n[rule.write]\nstatus = \"final\"\n",
		"id = \"fallback\"\n[rule.match.status]\neq = \"draft\"\n"+
			"[rule.match.recognized]\neq = \"step\"\n[rule.write]\nstatus = \"final\"\niter = \"+1\"\n", 1)
	if !strings.Contains(admitted, `iter = "07"`) || !strings.Contains(admitted, `iter = "+1"`) {
		t.Fatal("fixture substitution failed")
	}
	if env, _ := r12Lint(t, admitted); env != nil && env.Code == "model-invalid" {
		t.Errorf("a non-canonical int in [initial]/[rule.write] was refused at load; C5 reaches "+
			"predicate literals only:\n%s", env.raw)
	}

	// The same spelling on the predicate ingress of the same model refuses.
	guarded := strings.Replace(r12OwnedModel, "[rule.guard.all.iter]\neq = \"7\"\n",
		"[rule.guard.all.iter]\neq = \"07\"\n", 1)
	env, _ := r12Lint(t, guarded)
	if env == nil || env.Code != "model-invalid" ||
		!strings.Contains(env.raw, `\"07\" is not the canonical spelling of int tag iter; write 7`) {
		t.Errorf("guard literal \"07\" over int tag iter was not refused with the canonical-"+
			"spelling diagnostic; got %+v", env)
	}
}

// REQ-42: "The CLI is OUT of scope: `--tag n=07` / `--write n=+1` stay
// accepted."
// BOUNDARY
func TestReq42_0012_CLIValuesStayAcceptedAndCompareTyped(t *testing.T) {
	// --write iter=+1 is accepted on the owned door.
	model := writeFlowModel(t, r12OwnedModel)
	art := newFlowArtifact(t, "state.artifact")
	if env, _, _ := r12Run(t, "flow", "set-state", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--write", "status=draft", "--write", "iter=+1"); env != nil {
		t.Errorf("--write iter=+1 refused %s; the CLI is out of C5's scope:\n%s", env.Code, env.raw)
	}

	// --tag n=01 / n=+1 are accepted AND, on the guard path, decide typed.
	split := writeFlowModel(t, r12SplitModel)
	for _, v := range []string{"01", "+1"} {
		env, stdout, _ := r12Run(t, "flow", "resolve", "--model", split,
			"--outcome", "g", "--tag", "n="+v, "--plan-only")
		if env != nil {
			if env.Code == codeTagInvalid {
				t.Errorf("--tag n=%s refused flow-tag-invalid; the CLI stays accepting", v)
			} else {
				t.Errorf("--tag n=%s refused %s; the accepted value compares under the declared "+
					"int kind, so g-one matches", v, env.Code)
			}
			continue
		}
		if got := r12SelectedRule(t, stdout); got != "g-one" {
			t.Errorf("--tag n=%s selected %q; want g-one", v, got)
		}
	}
}

// REQ-43: "The refusal is user-visible and DIAGNOSTIC: it names the
// offending site and the canonical rewrite, e.g. `<site>: \"00\" is not the
// canonical spelling of int tag n; write 0`. It travels as whatever error
// value `(*loader).atom`'s existing `conform` call already returns to its
// caller … adding no new error type and no new sentinel."
// HAPPY PATH
func TestReq43_0012_TheRefusalNamesTheSiteAndTheRewrite(t *testing.T) {
	env := r12CanonicalRefusal(t, r12Cover07, "00", "0")

	// The site: the same location prefix the loader's existing
	// malformed_predicate_atom errors for this atom carry.
	bad, _ := r12Lint(t, r12Cover07With("[rule.guard.all.n]", `eq = "x"`))
	if bad == nil || len(bad.Findings) == 0 {
		t.Fatal("baseline: a non-int literal was not refused")
	}
	prefix, _, ok := strings.Cut(bad.Findings[0].Message, `"x" is not an int`)
	if !ok || !strings.Contains(prefix, "r-zero") {
		t.Fatalf("baseline message %q does not carry the site prefix", bad.Findings[0].Message)
	}
	want := prefix + `"00" is not the canonical spelling of int tag n; write 0`
	if env.Findings[0].Message != want {
		t.Errorf("message = %q; want %q — the site the existing conform error names, then "+
			"the canonical rewrite", env.Findings[0].Message, want)
	}
}

// REQ-44: "It reuses the one refusal code its ingress already carries —
// `malformed_predicate_atom` … — and adds no envelope field and no sixth
// refusal kind"
// BOUNDARY
func TestReq44_0012_TheRefusalReusesMalformedPredicateAtom(t *testing.T) {
	env := r12CanonicalRefusal(t, r12Cover07, "00", "0")
	bad, _ := r12Lint(t, r12Cover07With("[rule.guard.all.n]", `eq = "x"`))
	if bad == nil {
		t.Fatal("baseline: a non-int literal was not refused")
	}
	if env.Code != bad.Code {
		t.Errorf("refusal code = %q; the ingress's existing load refusal is %q", env.Code, bad.Code)
	}
	if !slices.Equal(env.keys, bad.keys) {
		t.Errorf("envelope keys = %v; the existing malformed_predicate_atom envelope carries %v",
			env.keys, bad.keys)
	}
	for _, f := range env.Findings {
		if f.Code != "malformed_predicate_atom" {
			t.Errorf("finding code %q; the refusal reuses malformed_predicate_atom", f.Code)
		}
	}

	// The flow door's load refusal carries the same finding code.
	model := writeFlowModel(t, r12Cover07)
	flowEnv, _, _ := r12Run(t, "flow", "resolve", "--model", model,
		"--outcome", "decide", "--tag", "n=0", "--plan-only")
	if flowEnv == nil || flowEnv.Code != "flow-model-invalid" ||
		len(flowEnv.Findings) == 0 || flowEnv.Findings[0].Code != "malformed_predicate_atom" {
		t.Errorf("flow resolve over cover07 did not refuse flow-model-invalid / "+
			"malformed_predicate_atom: %+v", flowEnv)
	}
}

// REQ-45: "\"Canonical\" here is exactly the `strconv.Itoa` round-trip
// stated above and nothing wider — no case folding, whitespace or unicode
// normalization is claimed or performed."
// DOMAIN EDGE
func TestReq45_0012_CanonicalIsExactlyTheItoaRoundTrip(t *testing.T) {
	wide := strings.Replace(r12Cover07, "min = 0\nmax = 1\n", "min = -5\nmax = 10\n", 1)
	with := func(atom string) string {
		return strings.Replace(wide, "[rule.guard.all.n]\neq = \"00\"\n", "[rule.guard.all.n]\n"+atom+"\n", 1)
	}

	// The round-trip's own fixed points are admitted.
	for _, atom := range []string{`eq = "-1"`, `eq = "0"`, `eq = "10"`, `eq = 7`} {
		if env, _ := r12Lint(t, with(atom)); env != nil && env.Code == "model-invalid" {
			t.Errorf("%s was refused at load; it is canonical:\n%s", atom, env.raw)
		}
	}

	// Its non-fixed points are refused as non-canonical.
	r12CanonicalRefusal(t, with(`eq = "+3"`), "+3", "3")

	// Spellings Atoi itself rejects are NOT normalized into canonical ones:
	// they keep the existing not-an-int diagnostic.
	for _, spelled := range []string{" 1", "1 ", "１", "0x1"} {
		env, _ := r12Lint(t, with(`eq = "`+spelled+`"`))
		if env == nil {
			t.Errorf("literal %q loaded; strconv.Atoi rejects it", spelled)
			continue
		}
		want := strconv.Quote(spelled) + " is not an int"
		kept := env.Code == "model-invalid" && slices.ContainsFunc(env.Findings, func(f flowFinding) bool {
			return f.Code == "malformed_predicate_atom" && strings.Contains(f.Message, want)
		})
		if !kept || strings.Contains(env.raw, "canonical spelling") {
			t.Errorf("literal %q did not keep the existing not-an-int load refusal (%s); no "+
				"whitespace or unicode normalization is performed:\n%s", spelled, want, env.raw)
		}
	}
}

// REQ-51 (`0012:S7`): "model `cover07` … invoked `intrastate lint --model
// cover07.toml --as json`. **Expected** (normative fixture, C5): exit 2 with
// the canonical-spelling load refusal naming the site and the rewrite —
// `\"00\" is not the canonical spelling of int tag n; write 0` — under the
// ingress's existing `malformed_predicate_atom` code." … "Also covers the
// bare-float spelling `eq = -0.0`"
// HAPPY PATH
func TestReq51_0012_Cover07RefusesWithTheCanonicalSpellingDiagnostic(t *testing.T) {
	env := r12CanonicalRefusal(t, r12Cover07, "00", "0")
	if !strings.Contains(env.Findings[0].Message, "r-zero") {
		t.Errorf("message %q does not name the site (row r-zero)", env.Findings[0].Message)
	}
	r12CanonicalRefusal(t, strings.Replace(r12Cover07, `eq = "00"`, "eq = -0.0", 1), "-0", "0")
}

// REQ-52 (`0012:S8`): "invoked `intrastate flow resolve --model mvv-int.toml
// --outcome step --tag iter=many --plan-only --as json`. **Expected**
// (normative fixture), unchanged before and after: exit 2," with the
// `flow-tag-invalid` envelope — the standing reason the CLI tag path cannot
// host the MVV's unevaluable step, which the OWNED door reaches.
// BOUNDARY
func TestReq52_0012_TheTagDoorRefusesUpstreamWhileTheOwnedDoorReachesTheSeam(t *testing.T) {
	model := writeFlowModel(t, r12ObservedModel)
	stdout, _, err := runCmd(t, "flow", "resolve", "--model", model, "--outcome", "step",
		"--tag", "iter=many", "--plan-only", "--as", "json")
	if err == nil {
		t.Fatalf("--tag iter=many succeeded; want flow-tag-invalid\n%s", stdout)
	}
	if got := clierr.ExitCodeFor(err); got != 2 {
		t.Errorf("exit = %d; want 2", got)
	}
	want := `{"code":"flow-tag-invalid","message":"the value for ` + "`iter`" +
		` does not conform to its declaration: \"many\" is not an int",` +
		`"schema_version":"0.1","param":"iter"}`
	if got := strings.TrimSpace(stdout); got != want {
		t.Errorf("envelope =\n%s\nwant (normative, unchanged)\n%s", got, want)
	}

	// The contrast the fixture exists for: the owned door carries the same
	// malformed value past ingress, and the typed seam refuses it.
	env, _, _ := r12ResolveOwned(t, "many")
	if env == nil || env.Code != "flow-guard-unevaluable" {
		t.Errorf("owned iter=many did not reach the seam as flow-guard-unevaluable: %+v", env)
	}
}

// REQ-53 (`0012:S9`): "a reader returning `07` for `iter` against the guard
// `iter eq 7`, on the `mvv-int` model. **Expected**: the guarded row
// MATCHES (parsed 7 == 7) where today's raw-string seam prunes it."
// HAPPY PATH
func TestReq53_0012_AnOwnedNonCanonicalValueMatchesUnderParsedComparison(t *testing.T) {
	env, stdout, _ := r12ResolveOwned(t, "07")
	if env == nil {
		t.Fatalf("owned iter=07 planned %q; the guarded row MATCHES (parsed 7 == 7), so it "+
			"conflicts with the unguarded fallback", r12SelectedRule(t, stdout))
	}
	if env.Code != "flow-ambiguous-match" {
		t.Fatalf("owned iter=07 refused %s; want flow-ambiguous-match naming guarded", env.Code)
	}
	if got := r12Rules(env); !slices.Equal(got, []string{"fallback", "guarded"}) {
		t.Errorf("ambiguity names %v; want [fallback guarded]", got)
	}
}

// REQ-54: "`flow next` derives its `unknown` pairs from the same probe seam
// … so a malformed owned int is now reported `uncomparable` there too, and a
// non-canonical held value (`\"07\"` against `n eq 7`) moves a row between
// excluded and candidate."
// ADVERSARIAL
func TestReq54_0012_FlowNextReportsTypedVerdicts(t *testing.T) {
	cands := r12NextOwned(t, "many")
	row, ok := cands["guarded"]
	if !ok {
		t.Error("flow next, owned iter=many: row `guarded` excluded; want it a candidate " +
			"with iter unknown `uncomparable`")
	} else if got := r12Unknown(row)["iter"]; got != "uncomparable" {
		t.Errorf("flow next, owned iter=many: unknown[iter] = %q; want uncomparable", got)
	}

	if _, ok := r12NextOwned(t, "07")["guarded"]; !ok {
		t.Error("flow next, owned iter=07: row `guarded` excluded; parsed 7 == 7 moves it to " +
			"the candidates")
	}

	// The control: a parsed value that differs keeps the row excluded.
	if _, ok := r12NextOwned(t, "4")["guarded"]; ok {
		t.Error("flow next, owned iter=4: row `guarded` is a candidate; parsed 4 != 7 excludes it")
	}
}

// --- REQ-50: the committed corpus does not flip ---------------------------------

const r12CorpusGolden = "testdata/rdr0012_corpus_verdicts.json"

// r12Verdict is one model's lint verdict: its outcome and finding codes.
type r12Verdict struct {
	Code     string   `json:"code"`
	Findings []string `json:"findings"`
}

// r12CorpusModels lists the committed corpus: models/ (the authored
// examples) plus every *.toml under internal/*/testdata, minus rdr.toml.
func r12CorpusModels(t *testing.T) []string {
	t.Helper()

	root := r12Root(t)
	var out []string
	for _, base := range []string{"models", "internal"} {
		err := filepath.WalkDir(filepath.Join(root, base), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".toml") || d.Name() == "rdr.toml" {
				return nil
			}
			rel, _ := filepath.Rel(root, path)
			rel = filepath.ToSlash(rel)
			if base == "internal" && !strings.Contains(rel, "/testdata/") {
				return nil
			}
			out = append(out, rel)
			return nil
		})
		if err != nil {
			t.Fatalf("walking %s: %v", base, err)
		}
	}
	sort.Strings(out)
	return out
}

func r12CorpusVerdicts(t *testing.T) map[string]r12Verdict {
	t.Helper()

	root := r12Root(t)
	out := map[string]r12Verdict{}
	for _, rel := range r12CorpusModels(t) {
		stdout, _, err := runCmd(t, "lint", "--model", filepath.Join(root, rel), "--as=json")
		var env struct {
			Type     string `json:"type"`
			Code     string `json:"code"`
			Findings []struct {
				Code string `json:"code"`
			} `json:"findings"`
			Data struct {
				Findings []struct {
					Code string `json:"code"`
				} `json:"findings"`
			} `json:"data"`
		}
		if jerr := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &env); jerr != nil {
			t.Fatalf("%s: lint stdout is not one JSON object: %v", rel, jerr)
		}
		v := r12Verdict{Code: "ok"}
		findings := env.Data.Findings
		if err != nil {
			v.Code = env.Code
			findings = env.Findings
		}
		for _, f := range findings {
			v.Findings = append(v.Findings, f.Code)
		}
		sort.Strings(v.Findings)
		out[rel] = v
	}
	return out
}

// REQ-50: "Lint/runtime agreement over the committed corpus — before/after
// verdict diff. **Expected**: no committed model's verdict flips." With "A
// green `make check` after this phase is the expected outcome — no existing
// corpus or testdata model should start failing"
// BOUNDARY
//
// The "before" is the golden recorded from the pre-change tree
// (`RDR0012_WRITE_CORPUS_GOLDEN=1` rewrites it). The "after" exists only
// once lint runs the typed seam, so that is asserted too.
func TestReq50_0012_NoCommittedModelsLintVerdictFlips(t *testing.T) {
	got := r12CorpusVerdicts(t)
	golden := filepath.Join(r12Root(t), "internal", "cli", r12CorpusGolden)

	if os.Getenv("RDR0012_WRITE_CORPUS_GOLDEN") != "" {
		buf, err := json.MarshalIndent(got, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, append(buf, '\n'), 0o600); err != nil {
			t.Fatal(err)
		}
		t.Logf("wrote %d verdicts to %s", len(got), golden)
	}

	raw, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("reading the before-verdict golden: %v", err)
	}
	var before map[string]r12Verdict
	if err := json.Unmarshal(raw, &before); err != nil {
		t.Fatalf("decoding the golden: %v", err)
	}
	for rel, b := range before {
		a, ok := got[rel]
		if !ok {
			t.Errorf("%s: in the before-corpus but no longer linted", rel)
			continue
		}
		if a.Code != b.Code || !slices.Equal(a.Findings, b.Findings) {
			t.Errorf("%s: lint verdict flipped: before %+v, after %+v", rel, b, a)
		}
	}

	file, _ := r12ParseFile(t, "internal/guard/product.go")
	if !r12ConstructsTyped(file) {
		t.Error("REQ-50 owed: lint's value-satisfaction site does not yet construct " +
			"NewEvaluator(DeclaredKinds(m)), so there is no typed \"after\" to diff")
	}
}

// --- REQ-MVV ------------------------------------------------------------------

// r12ProbeStep4 runs the guard package's MVV step-4 probes in a tagged
// child and returns each probe's terminal action and the child's non-event
// output.
func r12ProbeStep4(t *testing.T) (map[string]string, string) {
	t.Helper()

	cmd := exec.Command("go", "test", "-tags", "rdr0012probe", "-json",
		"-run", "^TestRDR0012Probe_ReqMVV_", "./internal/guard")
	cmd.Dir = r12Root(t)
	cmd.Env = append(os.Environ(), "RDR0012_PROBE_CHILD=1")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	_ = cmd.Run()

	results := map[string]string{}
	var build strings.Builder
	build.WriteString(stderr.String())
	sc := bufio.NewScanner(&stdout)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	for sc.Scan() {
		var ev struct{ Action, Test, Output string }
		if json.Unmarshal(sc.Bytes(), &ev) != nil {
			build.Write(sc.Bytes())
			build.WriteByte('\n')
			continue
		}
		switch {
		case ev.Test == "" && (ev.Action == "build-output" || ev.Action == "output"):
			build.WriteString(ev.Output)
		case ev.Test != "" && (ev.Action == "pass" || ev.Action == "fail" || ev.Action == "skip"):
			results[ev.Test] = ev.Action
		}
	}
	return results, build.String()
}

// REQ-MVV: the Minimum Viable Validation, end to end on the `mvv-int` model.
// The reader is driven by the persisted artifact, not by a test double:
// each step sets `iter`'s owned value in the artifact the accessor reads
// and re-runs the command.
// HAPPY PATH
func TestReqMVV_0012_DeclaredKindReachesTheSeamThroughTheOwnedDoor(t *testing.T) {
	// Step 2: `many` → flow-guard-unevaluable naming the guarded row and the
	// `iter` atom with reason `uncomparable`; the envelope carries the
	// existing code/message/schema_version fields and NO new top-level key.
	t.Run("step 2: many refuses uncomparable", func(t *testing.T) {
		env, stdout, err := r12ResolveOwned(t, "many")
		if env == nil {
			t.Fatalf("owned iter=many planned %q; want flow-guard-unevaluable (today's silent "+
				"prune to the fallback is the seed defect)", r12SelectedRule(t, stdout))
		}
		if env.Code != "flow-guard-unevaluable" {
			t.Fatalf("refused %s; want flow-guard-unevaluable", env.Code)
		}
		if got := clierr.ExitCodeFor(err); got != 2 {
			t.Errorf("exit = %d; want 2", got)
		}
		allowed := []string{"code", "findings", "message", "schema_version"}
		for _, k := range env.keys {
			if !slices.Contains(allowed, k) {
				t.Errorf("envelope carries top-level key %q; want only %v", k, allowed)
			}
		}
		named := false
		for _, f := range env.Findings {
			if f.Rule == "guarded" && f.Key == "iter" && f.Operator == "eq" && f.Literal == "7" &&
				f.Block == "all" && strings.Contains(f.Message, "guarded") &&
				strings.Contains(f.Message, "iter") && strings.Contains(f.Message, "uncomparable") {
				named = true
			}
		}
		if !named {
			t.Errorf("no finding names row `guarded`, the `iter eq 7` atom, and reason "+
				"`uncomparable`:\n%s", env.raw)
		}
	})

	// Step 3: `07` → the guarded row MATCHES under parsed comparison.
	t.Run("step 3: 07 matches parsed", func(t *testing.T) {
		env, stdout, _ := r12ResolveOwned(t, "07")
		if env == nil {
			t.Fatalf("owned iter=07 planned %q; the guarded row matches (parsed 7 == 7) and "+
				"conflicts with the fallback", r12SelectedRule(t, stdout))
		}
		if env.Code != "flow-ambiguous-match" || !slices.Contains(r12Rules(env), "guarded") {
			t.Errorf("owned iter=07 refused %s naming %v; want flow-ambiguous-match naming "+
				"guarded", env.Code, r12Rules(env))
		}
	})

	// Step 3 control: `4` → the guarded row prunes and the fallback plans.
	t.Run("step 3 control: 4 prunes", func(t *testing.T) {
		env, stdout, _ := r12ResolveOwned(t, "4")
		if env != nil {
			t.Fatalf("owned iter=4 refused %s; want the fallback plan:\n%s", env.Code, env.raw)
		}
		if got := r12SelectedRule(t, stdout); got != "fallback" {
			t.Errorf("owned iter=4 selected %q; want fallback", got)
		}
	})

	// Step 4: the extended suite over guard.NewEvaluator is green, and its
	// want-Unevaluable legs fail today's raw-string arms.
	t.Run("step 4: the extended suite", func(t *testing.T) {
		if os.Getenv("RDR0012_PROBE_CHILD") != "" {
			t.Skip("inside a probe child")
		}
		results, build := r12ProbeStep4(t)
		green := "TestRDR0012Probe_ReqMVV_ContractSuiteIsGreenOverNewEvaluator"
		raw := "TestRDR0012Probe_ReqMVV_Mutant_rawStringSeam"
		if results[green] == "" {
			t.Fatalf("REQ-MVV owed: the step-4 probe does not compile against this tree:\n%s", build)
		}
		if results[green] != "pass" {
			t.Errorf("TestGuardEvaluatorContract over guard.NewEvaluator: %s; want green", results[green])
		}
		if results[raw] != "fail" {
			t.Errorf("the extended suite %sed against today's raw-string arms; its int/bool "+
				"want-Unevaluable legs must fail them", results[raw])
		}
	})

	// Step 5: `lint --model` on `n eq "00"` over an int tag refuses with the
	// canonical-spelling diagnostic naming the site and the rewrite.
	t.Run("step 5: lint refuses 00", func(t *testing.T) {
		env := r12CanonicalRefusal(t, r12Cover07, "00", "0")
		if !strings.Contains(env.Findings[0].Message, "r-zero") {
			t.Errorf("the diagnostic %q does not name the site", env.Findings[0].Message)
		}
	})
}
