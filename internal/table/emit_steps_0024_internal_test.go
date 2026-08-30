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

	order := []string{"l.loadTags", "l.loadEmitDecls", "l.checkRuleEmit"}
	at := make([]int, len(order))
	for i, step := range order {
		at[i] = strings.Index(slice, step)
		if at[i] < 0 {
			t.Fatalf("`run`'s step slice does not hold `%s`", step)
		}
	}
	if !(at[0] < at[1] && at[1] < at[2]) {
		t.Errorf("the step slice orders loadTags/loadEmitDecls/checkRuleEmit "+
			"at %v; both are inserted IMMEDIATELY AFTER `loadTags`, and "+
			"`loadEmitDecls` MUST precede `checkRuleEmit`", at)
	}

	// "Immediately after": nothing else sits between them.
	between := slice[at[0]:at[2]]
	for _, other := range []string{
		"l.loadAccessors", "l.loadDump", "l.loadContexts", "l.loadInitial",
		"l.loadTerminal", "l.normalizeRules",
	} {
		if strings.Contains(between, other) {
			t.Errorf("`%s` sits between `loadTags` and `checkRuleEmit`; the "+
				"two new steps are inserted IMMEDIATELY after `loadTags`",
				other)
		}
	}

	// REQ-27: the slice is uniform — no conditional call wraps either step.
	if strings.Contains(slice, "if ") {
		t.Error("`run`'s step slice carries a conditional; the zero-" +
			"declaration opt-in gate is evaluated INSIDE each step")
	}
}

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
	if strings.Count(src, "atLine(") < 3 {
		t.Errorf("`load.go` calls `atLine` %d times; all three emit "+
			"categories MUST carry a source line stamped through it, "+
			"beside the tag refusal that already does",
			strings.Count(src, "atLine("))
	}
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
