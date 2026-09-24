package table_test

// Shared fixtures and oracles for the RDR 0030 step-write suite.
//
// Every model here is either a committed `ladder-*.toml` fixture or an
// inline variant built from one base, so a refusal fixture and its admitted
// twin differ ONLY in the defect under test. The admitted twin is what makes
// a refusal assertion discriminating: today's loader refuses every table
// value on the write path, so "it refused" alone proves nothing about the
// step grammar — "it refused, AND the same model with the defect removed
// loads" does.

import (
	"errors"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// s30Opts shapes one inline step model over the base ladder header.
type s30Opts struct {
	// attempt and tier replace the default declarations when non-empty.
	attempt string
	tier    string
	// decls are extra `[tags.*]` blocks; keys are extra OWNED keys the
	// state accessors serve.
	decls string
	keys  []string
	// initial replaces the default `[initial]` body when non-empty.
	initial string
	// contexts are extra `[context.*]` blocks.
	contexts string
	// head is extra top-level text placed before the tag declarations.
	head string
	// rules is the `[[rule]]` text.
	rules string
}

const s30DefaultAttempt = `[tags.attempt]
provenance = "owned"
kind = "int"
min = 0
max = 9
single_valued = true
required = true
`

const s30DefaultTier = `[tags.tier]
provenance = "owned"
kind = "enum"
domain = ["small", "mid", "large"]
single_valued = true
required = true
`

const s30DefaultInitial = `status = "open"
attempt = 0
tier = "small"
`

// s30Src renders an inline model. Its id is `m`, so a row identity reads
// `m.<rule>#<suffix…>`.
func s30Src(o s30Opts) string {
	attempt := o.attempt
	if attempt == "" {
		attempt = s30DefaultAttempt
	}
	tier := o.tier
	if tier == "" {
		tier = s30DefaultTier
	}
	initial := o.initial
	if initial == "" {
		initial = s30DefaultInitial
	}
	keys := append([]string{"status", "attempt", "tier"}, o.keys...)
	quoted := make([]string, len(keys))
	for i, k := range keys {
		quoted[i] = `"` + k + `"`
	}
	list := "[" + strings.Join(quoted, ", ") + "]"

	var b strings.Builder
	b.WriteString(`outcomes = ["retry", "escalate", "succeed"]
terminal = ["done"]
`)
	b.WriteString(o.head)
	b.WriteString(`
[model]
id = "m"
version = 1

[initial]
`)
	b.WriteString(initial)
	b.WriteString(`
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

`)
	b.WriteString(attempt)
	b.WriteString("\n")
	b.WriteString(tier)
	b.WriteString("\n")
	b.WriteString(o.decls)
	b.WriteString(`
[read.state]
role = "state"
path = "flow.state"
keys = ` + list + `
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ` + list + `
timeout = "2s"
read_back = true

[context.open.match.status]
eq = "open"

[context.done.match.status]
eq = "done"
`)
	b.WriteString(o.contexts)
	b.WriteString("\n")
	b.WriteString(o.rules)
	return b.String()
}

// s30Retry renders the `retry` rule: outcome `retry`, the given guard body
// under `[rule.guard.all.attempt]` (omitted when empty), and the given write
// body.
func s30Retry(guard, write string) string {
	var b strings.Builder
	b.WriteString("[[rule]]\nid = \"retry\"\nuse = [\"open\"]\n[rule.match.recognized]\neq = \"retry\"\n")
	if guard != "" {
		b.WriteString("[rule.guard.all.attempt]\n" + guard + "\n")
	}
	b.WriteString("[rule.write]\n" + write + "\n")
	return b.String()
}

// s30Escalate renders the `escalate` rule over tier.
func s30Escalate(guard, write string) string {
	var b strings.Builder
	b.WriteString("[[rule]]\nid = \"escalate\"\nuse = [\"open\"]\n[rule.match.recognized]\neq = \"escalate\"\n")
	if guard != "" {
		b.WriteString("[rule.guard.all.tier]\n" + guard + "\n")
	}
	b.WriteString("[rule.write]\n" + write + "\n")
	return b.String()
}

// s30Load loads inline source.
func s30Load(src string) (*table.Model, error) {
	return table.Load([]byte(src), "inline-0030.toml")
}

// s30MustLoad loads src and fails the test when the loader refuses it.
func s30MustLoad(t *testing.T, what, src string) *table.Model {
	t.Helper()

	m, err := s30Load(src)
	if err != nil {
		t.Fatalf("%s must LOAD (it is the admitted twin / expected-to-load "+
			"model); the loader refused: %v\n--- source ---\n%s", what, err, src)
	}
	return m
}

// s30Refusal loads src, which must refuse, and returns the categorized
// failure.
func s30Refusal(t *testing.T, what, src string) *table.Failure {
	t.Helper()

	_, err := s30Load(src)
	if err == nil {
		t.Fatalf("%s must be REFUSED at load; it loaded clean\n--- source ---\n%s", what, src)
	}
	var f *table.Failure
	if !errors.As(err, &f) {
		t.Fatalf("%s refused with an uncategorized error %v; every load refusal "+
			"carries a table.Failure", what, err)
	}
	return f
}

// s30RequireCategory asserts a refusal's category.
func s30RequireCategory(t *testing.T, what string, f *table.Failure, want table.Category) {
	t.Helper()

	if f.Category != want {
		t.Errorf("%s refused under %q; want %q (detail: %s)", what, f.Category, want, f.Detail)
	}
}

// s30RowIDs returns the identities of the rows one rule expands to, sorted.
func s30RowIDs(m *table.Model, rule string) []string {
	var out []string
	for _, r := range m.Rows {
		if r.RuleID == rule {
			out = append(out, r.Identity())
		}
	}
	slices.Sort(out)
	return out
}

// s30Row returns the single row with identity id.
func s30Row(t *testing.T, m *table.Model, id string) table.Row {
	t.Helper()

	for _, r := range m.Rows {
		if r.Identity() == id {
			return r
		}
	}
	t.Fatalf("no row %q; rows are %v", id, rowIdentities(m))
	return table.Row{}
}

// s30Value returns the single member a TagValue sequence binds to key.
func s30Value(t *testing.T, where string, tags []table.TagValue, key string) string {
	t.Helper()

	v, ok := tagValue(tags, key)
	if !ok {
		t.Fatalf("%s: no value for %q in %+v", where, key, tags)
	}
	if len(v) != 1 {
		t.Fatalf("%s: %q holds %d members %v; want exactly one", where, key, len(v), v)
	}
	return v[0]
}

var s30IntToken = regexp.MustCompile(`-?\d+`)

// s30FirstIntAfter returns the first integer token in detail that follows
// the first occurrence of anchor, or "" when there is none. "The detail names
// the cell FIRST" is asserted positionally with it: the anchor is the tag key
// the existing `rule <id> write <key>: ` prefix already names, so the first
// integer after it is the first number the record's new text contributes.
func s30FirstIntAfter(detail, anchor string) string {
	i := strings.Index(detail, anchor)
	if i < 0 {
		return ""
	}
	return s30IntToken.FindString(detail[i+len(anchor):])
}

// s30Contains asserts every needle appears in the refusal detail.
func s30Contains(t *testing.T, what, detail string, needles ...string) {
	t.Helper()

	for _, n := range needles {
		if !strings.Contains(detail, n) {
			t.Errorf("%s: the refusal detail does not name %q:\n  %s", what, n, detail)
		}
	}
}
