package table_test

// RDR 0030 Phase 3b — adversarial coverage anchored in the record's
// `Trade-offs / Failure Modes` section, first bullet:
//
//	A step write on a `scalar`, `set`, `bool`, or unbounded `int` — load
//	refusal, `malformed_tag_declaration`, detail names the admitted kinds.
//	An `int` whose width is not representable (`max - min` negative or
//	`math.MaxInt`) refuses the same way … so this is an explicit refusal
//	rather than an accidental fall-through to the zero-cell rule.
//
// That bullet draws a line at representability and says both sides of it
// behave deliberately. Two tests probe the line.
//
// ADV-1 — the side just INSIDE the line. `max - min == math.MaxInt - 1` is
// representable, so C1 admits it ("a wide but representable domain is
// admitted and merely expensive"), and A7 bounds the rows by the admitted
// cells. `step.go::stepCells` sizes its slice by the declared width
// (`make([]string, 0, width)`), so the loader panics with `makeslice: cap
// out of range` before any atom is consulted — a crash, not a refusal and
// not an expense. The test runs the load in a child process so a panic, or
// a walk that never returns, fails the test instead of the binary.
//
// ADV-3 — "refuses the same way". The kind refusal's detail names the
// admitted kinds; the width refusal's detail names neither admitted kind
// nor the admitted form, so the author is told what is wrong and not what
// is accepted (S5: "whose detail names the admitted form or kinds").

import (
	"context"
	"errors"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cwensel/intrastate/internal/table"
)

// advStepModel assembles a minimal state-machine model around one stepped
// `attempt` declaration and one `retry` rule stepping it by one under the
// rule's own positive atoms.
func advStepModel(attemptDecl, guard string) string {
	return `outcomes = ["retry", "succeed"]
terminal = ["done"]

[model]
id = "adv"
version = 1

[initial]
status = "open"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["open", "done"]
single_valued = true
required = true

[tags.attempt]
provenance = "owned"
kind = "int"
` + attemptDecl + `
single_valued = true

[read.state]
role = "state"
path = "flow.state"
keys = ["status", "attempt"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status", "attempt"]
timeout = "2s"
read_back = true

[context.open.match.status]
eq = "open"

[context.done.match.status]
eq = "done"

[[rule]]
id = "retry"
use = ["open"]
[rule.match.recognized]
eq = "retry"
[rule.guard.all.attempt]
` + guard + `
[rule.write]
attempt = { step = 1 }

[[rule]]
id = "succeed"
use = ["open"]
[rule.match.recognized]
eq = "succeed"
[rule.write]
status = "done"
`
}

const adv1ChildEnv = "INTRASTATE_ADV1_0030_CHILD"

// TestAdv1_0030_RepresentableWidthNeverPanics loads a step over the widest
// REPRESENTABLE int declaration — one below the width F1 refuses — whose
// rule admits five cells. C1 admits the declaration, so the load returns
// the five stepped rows; it must never panic and must return.
func TestAdv1_0030_RepresentableWidthNeverPanics(t *testing.T) {
	if os.Getenv(adv1ChildEnv) == "1" {
		adv1Child(t)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0],
		"-test.run=^TestAdv1_0030_RepresentableWidthNeverPanics$", "-test.count=1", "-test.v")
	cmd.Env = append(os.Environ(), adv1ChildEnv+"=1")
	out, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		t.Fatalf("load of a representable-width step admitting 5 cells did not return "+
			"within 30s: the admitted-cell walk is sized by the declared width\n%s", out)
	}
	if err != nil {
		t.Fatalf("load of a representable-width step (max-min == MaxInt-1) crashed or "+
			"failed; C1 admits it (representability, not size, is the bound): %v\n%s", err, out)
	}
}

func adv1Child(t *testing.T) {
	minV, maxV := 0, math.MaxInt-1
	if _, ok := resolveWidth(minV, maxV); !ok {
		t.Fatalf("fixture drifted: {%d..%d} is not representable", minV, maxV)
	}
	decl := "min = " + strconv.Itoa(minV) + "\nmax = " + strconv.Itoa(maxV)
	m, err := table.Load([]byte(advStepModel(decl, "gte = 0\nlt = 5")), "adv1.toml")
	if err != nil {
		t.Fatalf("load refused a representable width C1 admits: %v", err)
	}
	var got []string
	for _, r := range m.Rows {
		if r.RuleID == "retry" {
			got = append(got, strings.Join(r.Suffix, "#"))
		}
	}
	if want := []string{"0", "1", "2", "3", "4"}; strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("retry rows carry suffixes %v; want the admitted cells %v", got, want)
	}
}

// resolveWidth restates F1's representability line (`max - min` neither
// negative nor `math.MaxInt`) so the fixture proves it sits inside it.
func resolveWidth(minV, maxV int) (int, bool) {
	span := maxV - minV
	if minV > maxV || span < 0 || span == math.MaxInt {
		return 0, false
	}
	return span + 1, true
}

// TestAdv3_0030_WidthRefusalNamesAdmittedKinds pins F1's "refuses the same
// way": the non-representable width refusal names the admitted kinds, as
// the kind refusal beside it does.
func TestAdv3_0030_WidthRefusalNamesAdmittedKinds(t *testing.T) {
	cases := map[string]string{
		"kind (unbounded int)":           "max = 9",
		"width (max - min == MaxInt)":    "min = 0\nmax = " + strconv.Itoa(math.MaxInt),
		"width (wraps from math.MinInt)": "min = " + strconv.Itoa(math.MinInt) + "\nmax = -1",
	}
	for name, decl := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := table.Load([]byte(advStepModel(decl, "lt = 5")), "adv3.toml")
			var f *table.Failure
			if !errors.As(err, &f) {
				t.Fatalf("want a load refusal, got %v", err)
			}
			if f.Category != table.CatMalformedTagDeclaration {
				t.Fatalf("category %s; want %s (%s)", f.Category, table.CatMalformedTagDeclaration, f.Detail)
			}
			// Whole phrases: the width detail's "enumerate" would satisfy a
			// bare "enum" substring.
			if !strings.Contains(f.Detail, "an int") || !strings.Contains(f.Detail, "an enum") {
				t.Errorf("detail %q does not name the admitted kinds (an int declaring both "+
					"bounds, an enum with a non-empty domain); F1 has the width arm refuse "+
					"\"the same way\" as the kind arm", f.Detail)
			}
		})
	}
}
