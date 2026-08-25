package table_test

// Shared fixture helpers for the RDR 0002 loader/normalizer suite.
//
// Nothing here mocks the unit under test: every helper reads a promoted
// fixture off disk and hands the bytes to the real table.Load. The
// fixtures are the Stage-4/6 approved set from
// docs/rdr/0002-transition-table-as-reviewable-data/evidence/spikes/iter-2/,
// promoted per REQ-118 (extended, never narrowed).

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/newcoinc/intrastate/internal/table"
)

// Fixture source ids. The locator is derived from (model id, rule id), so
// the source id a caller supplies is the document identity the loader
// carries for diagnosis, not the locator itself (REQ-8, REQ-92).
const (
	rdrFixture  = "rdr-fixture.toml"
	kataFixture = "kata-fixture.toml"
)

// readFixture returns the bytes of a promoted fixture. Load takes bytes
// plus a source id and performs no file I/O of its own (REQ-8), so every
// test reads the file itself.
func readFixture(t *testing.T, rel string) []byte {
	t.Helper()

	b, err := os.ReadFile(filepath.Join("testdata", rel))
	if err != nil {
		t.Fatalf("read fixture %s: %v", rel, err)
	}
	return b
}

// mustLoad loads a fixture that the spec says must load clean and fails
// the test if the loader refuses it.
func mustLoad(t *testing.T, rel string) *table.Model {
	t.Helper()

	m, err := table.Load(readFixture(t, rel), rel)
	if err != nil {
		t.Fatalf("fixture %s must load clean; refused: %v", rel, err)
	}
	if m == nil {
		t.Fatalf("fixture %s loaded nil model with no error", rel)
	}
	return m
}

// loadCategory loads a fixture expected to refuse and returns the stable
// data-level category the refusal carries. The oracle is the category,
// never the message text (REQ-110, REQ-119, REQ-131).
func loadCategory(t *testing.T, rel string) table.Category {
	t.Helper()

	_, err := table.Load(readFixture(t, rel), rel)
	if err == nil {
		t.Fatalf("fixture %s loaded clean; want a refusal", rel)
	}
	cat, ok := table.CategoryOf(err)
	if !ok {
		t.Fatalf("fixture %s refused with a non-categorized error: %v", rel, err)
	}
	return cat
}

// rowByID returns the single normalized row whose full identity — the
// (model id, rule id, expansion suffix) tuple rendered as the row's
// identity string — matches want.
func rowByID(t *testing.T, m *table.Model, want string) table.Row {
	t.Helper()

	var found []table.Row
	for _, r := range m.Rows {
		if r.Identity() == want {
			found = append(found, r)
		}
	}
	switch len(found) {
	case 1:
		return found[0]
	case 0:
		t.Fatalf("no normalized row with identity %q; have %v", want, rowIdentities(m))
	default:
		t.Fatalf("%d normalized rows share identity %q; identities must be unique",
			len(found), want)
	}
	return table.Row{}
}

// rowsByRuleID returns every normalized row expanded from one source rule.
func rowsByRuleID(m *table.Model, ruleID string) []table.Row {
	var out []table.Row
	for _, r := range m.Rows {
		if r.RuleID == ruleID {
			out = append(out, r)
		}
	}
	return out
}

func rowIdentities(m *table.Model) []string {
	out := make([]string, 0, len(m.Rows))
	for _, r := range m.Rows {
		out = append(out, r.Identity())
	}
	return out
}

// atomsOn returns a row's normalized atoms on one tag key, in the row's
// own atom order.
func atomsOn(r table.Row, key string) []table.Atom {
	var out []table.Atom
	for _, a := range r.Atoms {
		if a.Key == key {
			out = append(out, a)
		}
	}
	return out
}

// atomsInBlock returns a row's normalized atoms sitting in one block.
func atomsInBlock(r table.Row, b table.Block) []table.Atom {
	var out []table.Atom
	for _, a := range r.Atoms {
		if a.Block == b {
			out = append(out, a)
		}
	}
	return out
}

// tagValue returns the value assigned to key in a []table.TagValue
// sequence (next-state tags or writes), and whether the key is present.
func tagValue(tags []table.TagValue, key string) ([]string, bool) {
	for _, t := range tags {
		if t.Key == key {
			return t.Value, true
		}
	}
	return nil, false
}
