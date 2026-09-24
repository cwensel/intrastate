package cli

// Shared fixtures and envelope helpers for the RDR 0030 CLI suites
// (`lint_mvv_0030_test.go`, `flow_mvv_0030_test.go`).
//
// Every stepped fixture has an UNROLLED twin authored the way a model is
// written today, and the oracle is agreement with that twin — not a
// re-derivation of what the expansion should produce. A twin-side assertion
// runs against today's loader, which is what proves each fixture is the
// shape its scenario claims before the stepped side can load at all.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/table"
)

const (
	c30Step          = "ladder-step.toml"
	c30Literal       = "ladder-literal.toml"
	c30StepBlocking  = "ladder-step-blocking.toml"
	c30LitBlocking   = "ladder-literal-blocking.toml"
	c30Unguarded     = "ladder-step-unguarded.toml"
	c30UnguardedEnum = "ladder-step-unguarded-enum.toml"
)

// c30Path returns a committed ladder fixture's path.
func c30Path(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join(repoRootFor(t), "internal", "table", "testdata", name)
}

// c30Source returns a committed ladder fixture's text.
func c30Source(t *testing.T, name string) string {
	t.Helper()

	b, err := os.ReadFile(c30Path(t, name))
	if err != nil {
		t.Fatalf("reading %s: %v", name, err)
	}
	return string(b)
}

// c30Edit replaces old with repl in src, failing when old is absent so a
// fixture edit can never silently not happen.
func c30Edit(t *testing.T, src, old, repl string) string {
	t.Helper()

	if !strings.Contains(src, old) {
		t.Fatalf("fixture edit anchor not found: %q", old)
	}
	return strings.Replace(src, old, repl, 1)
}

// c30EditAll replaces every occurrence of old.
func c30EditAll(t *testing.T, src, old, repl string) string {
	t.Helper()

	if !strings.Contains(src, old) {
		t.Fatalf("fixture edit anchor not found: %q", old)
	}
	return strings.ReplaceAll(src, old, repl)
}

// c30Model writes src to a temp model file and returns its path.
func c30Model(t *testing.T, src string) string {
	t.Helper()
	return writeModel(t, src)
}

// c30MustLoad loads src through the real loader.
func c30MustLoad(t *testing.T, what, src string) *table.Model {
	t.Helper()

	m, err := table.Load([]byte(src), what)
	if err != nil {
		t.Fatalf("%s must load; refused: %v", what, err)
	}
	return m
}

// c30Lint is one `lint --model <path> --as=json` run.
type c30Lint struct {
	failed   bool
	code     string
	findings []clierr.Finding
	stdout   string
}

func c30RunLint(t *testing.T, path string) c30Lint {
	t.Helper()

	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	r := c30Lint{failed: err != nil, stdout: stdout}
	if err != nil {
		env := parseFailure(t, stdout)
		r.code = env.Code
		if env.Findings != nil {
			r.findings = *env.Findings
		}
		return r
	}
	env := parseSuccess(t, stdout)
	if env.Data != nil && env.Data.Findings != nil {
		r.findings = *env.Data.Findings
	}
	return r
}

// c30LintSrc lints inline model source.
func c30LintSrc(t *testing.T, src string) c30Lint {
	t.Helper()
	return c30RunLint(t, c30Model(t, src))
}

// c30Loaded asserts a lint run got past the loader (it ran lint at all).
func c30Loaded(t *testing.T, what string, r c30Lint) {
	t.Helper()

	for _, f := range r.findings {
		if slices.Contains(table.Categories(), table.Category(f.Code)) {
			t.Fatalf("%s did not load: %s — %s", what, f.Code, f.Message)
		}
	}
	if r.failed && r.code != "graph-lint-failed" {
		t.Fatalf("%s: lint refused with %q before linting:\n%s", what, r.code, r.stdout)
	}
}

// c30LoadRefusal asserts lint refused src at LOAD under category and returns
// the finding carrying it.
func c30LoadRefusal(t *testing.T, what, src string, category table.Category) clierr.Finding {
	t.Helper()

	r := c30LintSrc(t, src)
	if !r.failed {
		t.Fatalf("%s: lint accepted the model; want the %s load refusal", what, category)
	}
	for _, f := range r.findings {
		if f.Code == string(category) {
			return f
		}
	}
	t.Fatalf("%s: lint refused with no %s finding:\n%s", what, category, r.stdout)
	return clierr.Finding{}
}

// c30Proj is S1's projection of one finding: (code, key, dimension, class,
// reason).
func c30Proj(f clierr.Finding) string {
	return strings.Join([]string{f.Code, f.Key, f.Dimension, f.Class, f.Reason}, "|")
}

// c30ProjSet is the sorted, duplicate-free projection of a finding list.
func c30ProjSet(fs []clierr.Finding) []string {
	var out []string
	for _, f := range fs {
		p := c30Proj(f)
		if !slices.Contains(out, p) {
			out = append(out, p)
		}
	}
	slices.Sort(out)
	return out
}

// c30Count counts findings carrying code.
func c30Count(fs []clierr.Finding, code string) int {
	n := 0
	for _, f := range fs {
		if f.Code == code {
			n++
		}
	}
	return n
}

// c30WithCode returns the findings carrying code.
func c30WithCode(fs []clierr.Finding, code string) []clierr.Finding {
	var out []clierr.Finding
	for _, f := range fs {
		if f.Code == code {
			out = append(out, f)
		}
	}
	return out
}

// c30JSON renders v for a failure message.
func c30JSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// c30StepRules returns, for a parsed model source, the ids of rules whose
// write block carries a `{ step = … }` value, keyed by (stepped key,
// outcome), plus the total rule count.
func c30StepRules(t *testing.T, doc map[string]any) (map[string][]string, int) {
	t.Helper()

	rules, _ := doc["rule"].([]any)
	out := map[string][]string{}
	for _, raw := range rules {
		rule, _ := raw.(map[string]any)
		id, _ := rule["id"].(string)
		write, _ := rule["write"].(map[string]any)
		match, _ := rule["match"].(map[string]any)
		rec, _ := match["recognized"].(map[string]any)
		var outcomes []string
		if eq, ok := rec["eq"].(string); ok {
			outcomes = append(outcomes, eq)
		}
		if in, ok := rec["in"].([]any); ok {
			for _, o := range in {
				s, _ := o.(string)
				outcomes = append(outcomes, s)
			}
		}
		for key, v := range write {
			if tbl, ok := v.(map[string]any); ok {
				if _, step := tbl["step"]; step {
					for _, o := range outcomes {
						out[key+"@"+o] = append(out[key+"@"+o], id)
					}
				}
			}
		}
	}
	return out, len(rules)
}
