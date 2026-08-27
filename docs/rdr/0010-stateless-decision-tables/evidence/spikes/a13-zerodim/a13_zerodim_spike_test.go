package a13zerodim_test

// A13 spike (RDR 0010). Read-only over committed source; nothing here is
// imported by internal/.
//
// Q1: does a group with zero participating guard dimensions and a bare
// escape row take `graph-coverage-closed-by-escape` from emitCoverageArms
// today? If so, an added `graph-unprovable-coverage` emission placed
// inside the `len(dims) == 0` branch AHEAD of emitCoverageArms would
// DOUBLE-REPORT unless the closure arm is suppressed, which is exactly the
// precedence C5 states.
//
// Q2: fixture sweep — which checked-in model fixtures carry a group with
// zero participating guard dimensions?

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/graphlint"
	"github.com/newcoinc/intrastate/internal/guard"
	"github.com/newcoinc/intrastate/internal/table"
)

const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true
`

const header = `
[model]
id = "t"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true
`

const accessors = `
[read.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"
read_back = true
`

func source(body string) string {
	var root, rest strings.Builder
	root.WriteString("\noutcomes = [\"go\", \"stop\"]\n")
	for _, line := range strings.SplitAfter(body, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "terminal =") {
			root.WriteString(line)
			continue
		}
		rest.WriteString(line)
	}
	return root.String() + header + decls + accessors + rest.String()
}

func load(t *testing.T, body string) *table.Model {
	t.Helper()
	m, err := table.Load([]byte(source(body)), "a13.toml")
	if err != nil {
		t.Fatalf("fixture must load clean; refused: %v\n---\n%s", err, source(body))
	}
	return m
}

func render(r graphlint.Report) string {
	if len(r.Findings) == 0 {
		return " (no findings)"
	}
	var b strings.Builder
	for _, f := range r.Findings {
		fmt.Fprintf(&b, "\n  %s sev=%s rule=%q element=%q reason=%q dim=%q class=%q\n    msg=%s",
			f.Code, f.Severity, f.Rule, f.Element, f.Reason, f.Dimension, f.Class, f.Message)
	}
	return b.String()
}

// dimsReport prints every group in a model with its participating guard
// dimensions, so a zero-dimension group is visible directly.
func dimsReport(m *table.Model) []string {
	var out []string
	for _, g := range guard.Groups(m) {
		d := guard.Dimensions(m, g)
		var ids []string
		for _, row := range g.Rows {
			ids = append(ids, row.RuleID)
		}
		sort.Strings(ids)
		out = append(out, fmt.Sprintf("group %q rows=%v dims=%v zero=%v",
			g.Context.String(), ids, d, len(d) == 0))
	}
	return out
}

// --- Q1 ------------------------------------------------------------------

// A zero-dimension group carrying a BARE escape row. This is the exact
// shape C5's precedence clause is about.
const zeroDimWithBareEscape = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "ordinary"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "rescue"
escape = ["no_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
`

// The same group with NO escape row, for contrast.
const zeroDimNoEscape = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "ordinary"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`

func TestQ1_ZeroDimGroupWithBareEscapeTakesClosedByEscapeToday(t *testing.T) {
	m := load(t, zeroDimWithBareEscape)

	t.Logf("groups:\n  %s", strings.Join(dimsReport(m), "\n  "))

	var zero int
	for _, g := range guard.Groups(m) {
		if len(guard.Dimensions(m, g)) == 0 {
			zero++
		}
	}
	if zero == 0 {
		t.Fatalf("fixture intended to carry a zero-dimension group, carries none")
	}

	r := graphlint.Run(graphlint.NewRequest(m))
	t.Logf("report:%s", render(r))

	var closed int
	for _, f := range r.Findings {
		if f.Code == graphlint.CodeCoverageClosedByEscape {
			closed++
		}
	}
	if closed == 0 {
		t.Errorf("Q1: expected graph-coverage-closed-by-escape over a "+
			"zero-dimension group carrying a bare escape row; got none:%s",
			render(r))
	}
	t.Logf("Q1: closed-by-escape count over zero-dim group with bare escape = %d", closed)
}

func TestQ1_ZeroDimGroupWithoutEscapeIsSilent(t *testing.T) {
	m := load(t, zeroDimNoEscape)
	t.Logf("groups:\n  %s", strings.Join(dimsReport(m), "\n  "))
	r := graphlint.Run(graphlint.NewRequest(m))
	t.Logf("report:%s", render(r))
	for _, f := range r.Findings {
		if f.Code == graphlint.CodeCoverageClosedByEscape {
			t.Errorf("unexpected closure finding without an escape row:%s", render(r))
		}
	}
}

// --- Q2: fixture sweep ---------------------------------------------------

// sweepRoots are the checked-in fixture trees. Every *.toml under them is
// attempted; ones the loader refuses (deliberate negative fixtures) are
// counted separately and never treated as zero-dimension hits.
var sweepRoots = []string{
	"internal/table/testdata",
	"models",
	"internal/graphlint/testdata",
	"internal/guard/testdata",
	"internal/resolve/testdata",
}

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for dir := wd; dir != "/"; dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
	}
	t.Fatal("repo root not found")
	return ""
}

func TestQ2_FixtureSweepForZeroDimensionGroups(t *testing.T) {
	root := repoRoot(t)

	var (
		loaded, refused, noRows int
		zeroHits                []string
		classOf                 = map[string]string{}
	)

	for _, rel := range sweepRoots {
		base := filepath.Join(root, rel)
		if _, err := os.Stat(base); err != nil {
			t.Logf("sweep root absent, skipped: %s", rel)
			continue
		}
		err := filepath.Walk(base, func(p string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || filepath.Ext(p) != ".toml" {
				return nil
			}
			src, err := os.ReadFile(p)
			if err != nil {
				return nil
			}
			m, err := table.Load(src, p)
			if err != nil || m == nil {
				refused++
				return nil
			}
			loaded++
			rel, _ := filepath.Rel(root, p)
			classOf[rel] = declaredClass(string(src))

			groups := guard.Groups(m)
			if len(groups) == 0 {
				noRows++
			}
			for _, g := range groups {
				if len(guard.Dimensions(m, g)) == 0 {
					var ids []string
					for _, row := range g.Rows {
						ids = append(ids, row.RuleID)
					}
					sort.Strings(ids)
					zeroHits = append(zeroHits, fmt.Sprintf(
						"%s  class=%s  group=%q rows=%v",
						rel, classOf[rel], g.Context.String(), ids))
				}
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	sort.Strings(zeroHits)
	t.Logf("sweep: loaded=%d refused=%d models-with-no-groups=%d "+
		"zero-dimension-groups=%d", loaded, refused, noRows, len(zeroHits))
	if len(zeroHits) == 0 {
		t.Logf("zero-dimension groups: NONE")
		return
	}
	t.Logf("zero-dimension groups:\n  %s", strings.Join(zeroHits, "\n  "))
}

// declaredClass reads the `[model].class` line if the fixture authors one.
// Nothing in the committed loader is assumed to carry it; absent means the
// state-machine default.
func declaredClass(src string) string {
	for _, line := range strings.Split(src, "\n") {
		l := strings.TrimSpace(line)
		if strings.HasPrefix(l, "class") && strings.Contains(l, "=") {
			return strings.TrimSpace(strings.SplitN(l, "=", 2)[1])
		}
	}
	return "(none: state-machine default)"
}

// --- Q1b: is C5's precedence achievable at that site? --------------------

// The proposed patch is, inside checkCoverage's `len(dims) == 0` branch and
// AHEAD of emitCoverageArms:
//
//	if len(dims) == 0 {
//	    if a.model.Class == "decision-table" {          // class-keyed arm
//	        a.emit(Finding{Code: CodeUnprovableCoverage,
//	                       Reason: ReasonNoParticipatingDimension, ...})
//	        return                                       // <- the precedence
//	    }
//	    a.emitCoverageArms(g)
//	    return
//	}
//
// Q1's finding is that emitCoverageArms DOES emit
// `graph-coverage-closed-by-escape` over a zero-dimension group carrying a
// bare escape row (TestQ1_... above, count = 1). So an added emission that
// merely PRECEDES emitCoverageArms without returning would produce TWO
// findings for the one group — the new unprovable-coverage arm and the
// closure advisory — which is exactly the double-report C5 forbids.
//
// The precedence is nevertheless achievable at that site, because the
// branch is a self-contained `emitCoverageArms(g); return` pair: replacing
// the call with a return SUPPRESSES the closure arm for this group only,
// touches no other invariant, and leaves every other group's path
// unchanged. This test asserts the two structural facts the argument rests
// on: the closure finding is the ONLY thing emitCoverageArms produces over
// a zero-dimension group, and it is produced solely from the bare-escape
// path (no coverage-gap can accompany it, because the empty product is
// closed by any row's membership).
func TestQ1b_ZeroDimBranchEmitsOnlyTheClosureArm(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
	}{
		{"bare-escape", zeroDimWithBareEscape},
		{"no-escape", zeroDimNoEscape},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := load(t, tc.body)
			r := graphlint.Run(graphlint.NewRequest(m))
			for _, f := range r.Findings {
				switch f.Code {
				case graphlint.CodeCoverageClosedByEscape:
					// expected in the bare-escape case only
				case graphlint.CodeCoverageGap, graphlint.CodeUnprovableCoverage,
					graphlint.CodeProductTooLarge:
					t.Errorf("zero-dim group produced an invariant-4 finding "+
						"other than the closure arm: %s", render(r))
				}
			}
			t.Logf("%s report:%s", tc.name, render(r))
		})
	}
}
