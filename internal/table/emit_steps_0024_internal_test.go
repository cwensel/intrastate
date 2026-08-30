package table

// RDR 0024 `0024:C2` — the two step signatures and their position.
//
// These clauses are about the SHAPE of the loader, not about a refusal, so
// they are asserted from inside the package. The step slice holds bound
// method values, which is the reason C2 states the signatures inside a
// normative fence (ASSUMPTION-2 reads them as binding).

import (
	"os"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// REQ-26: "**Both are methods on `*loader` taking no arguments and
// returning `error`** — `func (l *loader) loadEmitDecls() error` and `func
// (l *loader) checkRuleEmit() error`"
// ASSUMPTION-2: read as BINDING, not illustrative, because C2 states them
// inside a normative fence with the reason.
// DOMAIN EDGE
func TestReq26_0024_BothStepsAreNoArgErrorReturningLoaderMethods(t *testing.T) {
	l := &loader{}

	// The oracle is a compile-time assignment to `func() error`, not
	// reflection: `reflect.MethodByName` cannot see an unexported method,
	// so a reflective test could never go green. A bound method value only
	// assigns to this type when the receiver, arity, and result match
	// exactly — which is precisely what the step slice requires.
	var steps []func() error
	steps = append(steps, l.loadEmitDecls, l.checkRuleEmit)

	if len(steps) != 2 {
		t.Fatalf("bound %d step method values; C2 names two", len(steps))
	}
	for i, step := range steps {
		if step == nil {
			t.Errorf("step %d bound to nil", i)
		}
	}

	// And the step slice's own element type is that same signature, which
	// is the reason C2 states the signatures normatively.
	if got := reflect.TypeOf(steps).Elem().String(); got != "func() error" {
		t.Errorf("the step slice's element type is %s; both steps are "+
			"`*loader` methods taking NO arguments and returning `error`",
			got)
	}
}

// REQ-25: "both are inserted immediately after `loadTags` in
// `load.go::run`'s step slice", and "`loadEmitDecls` MUST precede
// `checkRuleEmit`".
// REQ-27: "the step slice is uniform and has no room for a conditional
// call, so both steps run unconditionally".
// DOMAIN EDGE
func TestReq25_0024_BothStepsSitImmediatelyAfterLoadTagsInTheStepSlice(t *testing.T) {
	src := readLoadGoSource(t)

	start := strings.Index(src, "[]func() error{")
	if start < 0 {
		t.Fatal("`load.go::run` no longer builds a `[]func() error` step " +
			"slice; C2's position clause is stated against that slice")
	}
	end := strings.Index(src[start:], "\n\t} {")
	if end < 0 {
		t.Fatal("could not find the end of `run`'s step slice")
	}
	slice := src[start : start+end]

	// The slice is parsed into its ELEMENTS once, and every clause below is
	// stated against that ordered list. An index-of-substring reading
	// cannot say "immediately after": it is satisfied by any unlisted bare
	// step sitting between `loadTags` and the emit steps, since such a step
	// appears in no exclusion list.
	//
	// REQ-27: the slice is uniform — every element is a bare bound method
	// value, so there is no room for a conditional call. A
	// `strings.Contains(slice, "if ")` guard is not that claim either: a
	// wrapper such as `conditionally(l.checkRuleEmit)` carries no `if `.
	// The shape of each element is what the clause is about, so each
	// element is matched, and a non-conforming one is what makes the
	// adjacency reading below trustworthy.
	//
	// The span opens on the composite-literal token itself; every LATER
	// line is an element or a comment.
	var steps []string
	for _, raw := range strings.Split(slice, "\n")[1:] {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		m := stepElement.FindStringSubmatch(raw)
		if m == nil {
			t.Fatalf("`run`'s step slice carries the element %q; the slice "+
				"is UNIFORM — every element is a bare `l.<method>,` bound "+
				"method value, with the zero-declaration opt-in gate "+
				"evaluated INSIDE each step", line)
		}
		steps = append(steps, m[1])
	}

	// REQ-25: `loadTags`, `loadEmitDecls` and `checkRuleEmit` occupy three
	// CONSECUTIVE positions, in that order.
	want := []string{"loadTags", "loadEmitDecls", "checkRuleEmit"}
	tags := slices.Index(steps, "loadTags")
	if tags < 0 {
		t.Fatalf("`run`'s step slice does not hold `loadTags`; it holds %v",
			steps)
	}
	if tags+len(want) > len(steps) ||
		!slices.Equal(steps[tags:tags+len(want)], want) {
		t.Fatalf("`run`'s step slice reads %v; both emit steps are inserted "+
			"IMMEDIATELY AFTER `loadTags` — the three occupy consecutive "+
			"positions %v — and `loadEmitDecls` MUST precede `checkRuleEmit`",
			steps, want)
	}

	// And `normalizeRules` FOLLOWS all three, positively: C2 places the
	// emit steps ahead of it so both read the SOURCE rules. Adjacency alone
	// cannot say this — it is satisfied by the whole trio sitting AFTER
	// `normalizeRules`.
	norm := slices.Index(steps, "normalizeRules")
	if norm < 0 {
		t.Fatalf("`run`'s step slice does not hold `normalizeRules`; it "+
			"holds %v", steps)
	}
	if norm < tags+len(want) {
		t.Errorf("`run`'s step slice reads %v; both emit steps sit AHEAD of "+
			"`normalizeRules` so they read the SOURCE rules rather than "+
			"normalized rows", steps)
	}
}

// stepElement matches one element of `run`'s step slice — a bare bound
// method value at the slice's own indentation, and nothing else — and
// captures the method name.
var stepElement = regexp.MustCompile(`^\t\tl\.([A-Za-z]+),$`)

// REQ-25 behaviourally: `checkRuleEmit` runs BEFORE `normalizeRules`.
//
// The positional oracle above reads `load.go` as text, so it proves where
// the tokens sit, not what the loader does. This one drives a fixture
// carrying BOTH an undeclared emit key AND a defect only `normalizeRules`
// sees (a duplicate rule id), through the public `Load`. Only one refusal
// comes back — the tier is fail-fast — and it is `CatUnknownEmitKey`
// exactly when the emit step ran first. Reorder the slice and this goes to
// `CatDuplicateRuleID`; no amount of text in `load.go` can hold it green.
// ADVERSARIAL
func TestReq25_0024_TheEmitStepRefusesAheadOfNormalization(t *testing.T) {
	_, err := Load([]byte(emitAheadOfNormalizeSrc0024), "emit-ahead-of-normalize")
	if err == nil {
		t.Fatal("the doubly-defective fixture loaded clean; it carries both " +
			"an undeclared emit key and a duplicate rule id")
	}
	cat, ok := CategoryOf(err)
	if !ok {
		t.Fatalf("the refusal carries no category: %v", err)
	}
	if cat != CatUnknownEmitKey {
		t.Errorf("the fixture refuses with %s; `checkRuleEmit` is inserted "+
			"AHEAD of `normalizeRules` (`0024:C2`), so the emit defect is "+
			"the one the fail-fast tier reports, not %s", cat, CatDuplicateRuleID)
	}
}

// emitAheadOfNormalizeSrc0024 is doubly defective ON PURPOSE, which is why
// it is not a promoted `testdata/neg` witness: those are single-defect by
// construction (REQ-105 exclusivity). `nxet` is declared nowhere, and
// `cell-x` is declared twice.
const emitAheadOfNormalizeSrc0024 = `outcomes = ["decide"]

[model]
id = "emit-ahead-of-normalize"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.a]
provenance = "observed"
kind = "enum"
domain = ["x", "y"]
single_valued = true
required = true

[emit.verdict]
kind = "enum"
domain = ["alpha", "beta"]

[[rule]]
id = "cell-x"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "x"
[rule.emit]
verdict = "alpha"
nxet = "typo"

[[rule]]
id = "cell-x"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "y"
[rule.emit]
verdict = "beta"
`

// REQ-33: "The decoder is not an alternative route — `pelletier/go-toml/v2`
// exposes position only on `DecodeError` with unexported fields"
// REQ-76 / `PH1` / ASSUMPTION-3: the position is recovered by a source
// scan, and "recovering BOTH lines" is the obligation; the helper's name
// and arity are the implementer's.
// ADVERSARIAL — negative REQ: no decoder-position plumbing.
func TestReq33_0024_NoDecoderPositionPlumbingIsIntroduced(t *testing.T) {
	src := readLoadGoSource(t)

	for _, token := range []string{
		"toml.DecodeError", "DecodeError", "*toml.StrictMissingError",
	} {
		if !strings.Contains(src, token) {
			continue
		}
		// A mention is only a defect when it is used to RECOVER A POSITION.
		for _, positional := range []string{".Row(", ".Column(", ".Position"} {
			if strings.Contains(src, positional) {
				t.Errorf("`load.go` reads %q off %q; the decoder is not an "+
					"alternative route — it exposes position only on "+
					"`DecodeError` with unexported fields", positional, token)
			}
		}
	}

	// The positive leg: both lines are recovered by a SOURCE SCAN, the way
	// `tagHeaderLine` already does it. `atLine` is the one stamping seam.
	if !strings.Contains(src, "func tagHeaderLine(") {
		t.Error("`load.go` no longer declares `tagHeaderLine`; C2 says the " +
			"declaration-side locator reaches its position by that " +
			"technique, unchanged")
	}
	// The count of textual `atLine(` occurrences is NOT asserted here.
	// Comments and a fourth call site perturb it in both directions, and
	// `TestReq30_0024` already asserts every one of the three categories
	// carries a non-zero `Line` — the behavioural form of the same claim,
	// which strictly dominates a token census.
}

// readLoadGoSource returns `internal/table/load.go`'s text. The two clauses
// above are about the loader's SHAPE — a step slice's order and the absence
// of a decoder-position route — neither of which has a runtime observable.
func readLoadGoSource(t *testing.T) string {
	t.Helper()

	body, err := os.ReadFile("load.go")
	if err != nil {
		t.Fatalf("read load.go: %v", err)
	}
	return string(body)
}
