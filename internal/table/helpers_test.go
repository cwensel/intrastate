package table_test

// Shared fixture helpers for the RDR 0002 loader/normalizer suite.
//
// Nothing here mocks the unit under test: every helper reads a promoted
// fixture off disk and hands the bytes to the real table.Load. The
// fixtures are the Stage-4/6 approved set from
// docs/rdr/0002-transition-table-as-reviewable-data/evidence/spikes/iter-2/,
// promoted per REQ-118 (extended, never narrowed).

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/newcoinc/intrastate/internal/resolve"
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
	// A loader returning an empty model with no error would let every
	// "for each row" assertion below pass vacuously. No fixture in the
	// promoted set normalizes to zero rows.
	if len(m.Rows) == 0 {
		t.Fatalf("fixture %s loaded clean but normalized to no candidate rows", rel)
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

// expectedSeamLiteral re-derives the byte form a normalized atom's literal
// must take when it crosses the kernel seam as resolve.GuardAtom.Literal.
//
// This restates D11's rule test-side rather than reaching for
// model.go::seamValue and ::setValuedLiteral, which are unexported while
// this suite is package table_test. Restating it is the point: an oracle
// that called the production helper would agree with any carriage the
// implementation chose, including a wrong one. Carriage is keyed on the
// OPERATOR ALONE (D8, D11) — `in`/`contains` cross as the JDR 0001 §D13
// canonical JSON array (members sorted, duplicate-free, compact), and
// every other operator — `exists`, `eq`, the integer comparisons —
// crosses as its bare single value. If a future operator gains set
// carriage this derivation fails loudly against correct production code,
// which is the right direction for a rule this oracle exists to pin.
func expectedSeamLiteral(t *testing.T, a table.Atom) string {
	t.Helper()

	if a.Operator != "in" && a.Operator != "contains" {
		if len(a.Literal) == 0 {
			return ""
		}
		return a.Literal[0]
	}
	b, err := json.Marshal(slices.Compact(slices.Sorted(slices.Values(a.Literal))))
	if err != nil {
		t.Fatalf("marshal set literal %v: %v", a.Literal, err)
	}
	return string(b)
}

// assertGuardBijection asserts that a row's kernel guard payload is a
// BIJECTION onto the row's normalized all/unless atoms: every guard entry
// consumes exactly one atom and every atom is consumed.
//
// Cardinality alone does not carry REQ-85's "exhaustive and disjoint": an
// implementation emitting one guard atom twice while dropping another has
// the right |Guard|, carries no BlockMatch, and traces every Match entry
// correctly. The consumed-at-most-once bookkeeping below is what rejects
// it — matching without it would let the duplicate pair with the same
// atom twice.
func assertGuardBijection(t *testing.T, row table.Row, guard []resolve.GuardAtom) {
	t.Helper()

	atoms := append(
		slices.Clone(atomsInBlock(row, table.BlockAll)),
		atomsInBlock(row, table.BlockUnless)...)
	consumed := make([]bool, len(atoms))

	for _, g := range guard {
		matched := -1
		for i, a := range atoms {
			if consumed[i] {
				continue
			}
			if a.Key == g.Key && a.Operator == g.Operator &&
				a.Block == g.Block && expectedSeamLiteral(t, a) == g.Literal {
				matched = i
				break
			}
		}
		if matched < 0 {
			t.Errorf("%s: guard entry %+v consumes no unconsumed all/unless atom; "+
				"the routing is not injective over %+v", row.Identity(), g, atoms)
			continue
		}
		consumed[matched] = true
	}

	for i, a := range atoms {
		if !consumed[i] {
			t.Errorf("%s: all/unless atom %+v reached no guard entry; the routing "+
				"is not exhaustive (guard = %+v)", row.Identity(), a, guard)
		}
	}
}
