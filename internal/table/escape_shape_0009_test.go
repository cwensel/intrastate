package table_test

// RDR 0009 — the authored-path half of escape-row shape conformance, plus
// the Phase 3 shared conformance fixture set.
//
// Ownership: this RDR is the single normative home of the escape-row shape
// rule; RDR 0002's normalizer/lint is the ENFORCING implementation on the
// authored path. Several clauses in this file are already satisfied at HEAD
// (`internal/table/normalize.go` rejects an escape rule carrying a write
// block, a clear list, or a gate list, keyed on key PRESENCE). Those tests
// pin shipped behaviour so a future change cannot relax it; the file records
// which ones those are in each test's own comment.
//
// The Phase 3 fixtures at the foot of this file are the SHARED conformance
// set REQ-97/REQ-98 ask this RDR to ship: one authored-clear case whose
// `<clear>`-sentinel representation is pinned identically for the kernel
// suite and the normalizer suite, so the two enforcers of one invariant
// cannot drift.

import (
	"slices"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/resolve"
	"github.com/newcoinc/intrastate/internal/table"
)

// escapeShapeBaseModel is the minimal authored document every case below
// extends. It declares one owned tag, one recognized outcome, and one
// ordinary rule, so the only thing the appended escape rule varies is its
// own shape.
const escapeShapeBaseModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "escapeshape"
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

[read.state]
role = "state"
path = "flow.state"
keys = ["status"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status"]
timeout = "2s"
read_back = true

[gate.approval]
role = "state"
path = "flow.gate"
keys = ["status"]
timeout = "2s"

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "advance-draft"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"
`

// loadEscapeShape loads the base model plus the supplied rule fragment.
func loadEscapeShape(t *testing.T, fragment string) (*table.Model, error) {
	t.Helper()

	src := escapeShapeBaseModel + fragment
	return table.Load([]byte(src), "escape-shape-0009.toml")
}

// The three malformed escape declarations REQ-11/REQ-12 name, plus the
// conforming control.
const (
	escapeWithWriteBlock = `
[[rule]]
id = "escape-with-write"
escape = ["no_match"]
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"
`
	escapeWithClearList = `
[[rule]]
id = "escape-with-clear"
escape = ["no_match"]
clear = ["status"]
[rule.match.recognized]
eq = "advance"
`
	escapeWithEmptyWriteBlock = `
[[rule]]
id = "escape-with-empty-write"
escape = ["no_match"]
[rule.match.recognized]
eq = "advance"
[rule.write]
`
	escapeWithEmptyClearList = `
[[rule]]
id = "escape-with-empty-clear"
escape = ["no_match"]
clear = []
[rule.match.recognized]
eq = "advance"
`
	conformingEscapeRule = `
[[rule]]
id = "escape-conforming"
escape = ["no_match"]
[rule.match.recognized]
eq = "advance"
`
)

// REQ-11: "On the authored-table path, RDR 0002's normalizer/lint is the
// enforcing implementation: an authored escape rule carrying a write block or
// clear list MUST be rejected at table load under RDR 0002's "malformed escape
// declaration" validation class, diagnosable to the source rule."
// HAPPY PATH
//
// ALREADY SATISFIED AT HEAD — `normalize.go::normalizeRule` rejects both.
// This test pins it: the rejection must stay at LOAD, must carry the named
// category (never a generic one), and must remain diagnosable to the rule id.
func TestReq11_AnAuthoredEscapeRuleCarryingWritesOrClearsIsRejectedAtLoad(t *testing.T) {
	cases := map[string]struct {
		fragment string
		ruleID   string
	}{
		"write_block": {escapeWithWriteBlock, "escape-with-write"},
		"clear_list":  {escapeWithClearList, "escape-with-clear"},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			m, err := loadEscapeShape(t, tc.fragment)
			if err == nil {
				t.Fatalf("the escape rule loaded clean (%d rows); an escape "+
					"rule carrying a write block or clear list MUST be "+
					"rejected at table load", len(m.Rows))
			}
			cat, ok := table.CategoryOf(err)
			if !ok {
				t.Fatalf("the refusal carries no validation category: %v", err)
			}
			if cat != table.CatMalformedEscapeDeclaration {
				t.Errorf("category = %q; want %q — the rejection is under "+
					"RDR 0002's malformed escape declaration class",
					cat, table.CatMalformedEscapeDeclaration)
			}
			// Diagnosable to the SOURCE RULE: the offending rule id must be
			// recoverable from the refusal, not merely "some rule".
			if !strings.Contains(err.Error(), tc.ruleID) {
				t.Errorf("the refusal does not name the source rule %q: %v",
					tc.ruleID, err)
			}
		})
	}
}

// REQ-12: "Rejection keys on the PRESENCE of the block, not its contents — RDR
// 0002 forbids an escape rule from containing a write block, so an empty one
// (`writes = []`) is rejected too."
// INPUT EDGE
//
// ALREADY SATISFIED AT HEAD — the loader switches on `rule.Write != nil`, so
// a present-but-empty block decodes non-nil and is rejected. This is the
// clause a length-based reading would break, so it is pinned for both the
// empty write block and the empty clear list.
func TestReq12_AnEmptyWriteBlockOnAnEscapeRuleIsRejectedToo(t *testing.T) {
	cases := map[string]struct {
		fragment string
		ruleID   string
	}{
		"empty_write_block": {escapeWithEmptyWriteBlock, "escape-with-empty-write"},
		"empty_clear_list":  {escapeWithEmptyClearList, "escape-with-empty-clear"},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			m, err := loadEscapeShape(t, tc.fragment)
			if err == nil {
				t.Fatalf("an escape rule carrying an EMPTY block loaded "+
					"clean (%d rows); rejection keys on the PRESENCE of the "+
					"block, not on its contents", len(m.Rows))
			}
			if cat, _ := table.CategoryOf(err); cat != table.CatMalformedEscapeDeclaration {
				t.Errorf("category = %q; want %q", cat,
					table.CatMalformedEscapeDeclaration)
			}
			// Like REQ-11's rejection, the refusal on a present-but-empty
			// block stays diagnosable to the SOURCE RULE: an operator must
			// be able to find which rule was refused, not merely that one
			// was.
			if !strings.Contains(err.Error(), tc.ruleID) {
				t.Errorf("the refusal does not name the source rule %q: %v",
					tc.ruleID, err)
			}
		})
	}
}

// REQ-13: "The authored layer is therefore stricter than the kernel's
// length-based predicate, deliberately"
// ADVERSARIAL
//
// The strictness GAP is the claim, and it is only assertable by showing the
// two layers disagree on one input: an escape rule with an empty write block
// is REJECTED by the authored layer, while the equivalent kernel row (empty
// Writes) CONFORMS. A future change that aligned the two — by relaxing the
// loader to a length check — would close the gap and fail here.
func TestReq13_TheAuthoredLayerIsStricterThanTheKernelPredicate(t *testing.T) {
	// Authored layer: rejected.
	if _, err := loadEscapeShape(t, escapeWithEmptyWriteBlock); err == nil {
		t.Errorf("the authored layer accepted an empty write block on an " +
			"escape rule; it is deliberately stricter than the kernel")
	}

	// Kernel layer: the same shape, expressed as a normalized row, conforms.
	kernelEquivalent := resolve.Table{
		Revision: "escapeshape",
		Outcomes: []string{"advance"},
		Rows: []resolve.Row{{
			RuleID:        "escape-with-empty-write",
			SourceLocator: "escape-shape-0009.toml:1",
			Outcome:       "advance",
			Escape:        []resolve.RefusalKind{resolve.KindNoMatch},
			Writes:        []resolve.Tag{},
		}},
	}
	if err := kernelEquivalent.CheckValid(); err != nil {
		t.Errorf("the kernel rejected an empty Writes slice (%v); the "+
			"kernel predicate is length-based and the authored layer is the "+
			"stricter of the two", err)
	}
}

// REQ-14: "Normalized escape rows MUST render write-free. This RDR binds that
// duty with a conformance fixture set; RDR 0002's grammar is unchanged."
// HAPPY PATH
//
// The positive obligation the loader owes the kernel: a conforming authored
// escape rule must normalize to a row the kernel's predicate accepts — both
// as `table.Row` and after the `KernelRow`/`KernelTable` handoff.
func TestReq14_NormalizedEscapeRowsRenderWriteFree(t *testing.T) {
	m, err := loadEscapeShape(t, conformingEscapeRule)
	if err != nil {
		t.Fatalf("the conforming escape rule was refused: %v", err)
	}

	var escapes int
	for _, row := range m.Rows {
		if len(row.Escape) == 0 {
			continue
		}
		escapes++
		if len(row.Writes) != 0 {
			t.Errorf("normalized escape row %q renders %d writes; escape "+
				"rows MUST render write-free", row.RuleID, len(row.Writes))
		}
	}
	if escapes == 0 {
		t.Fatalf("the model normalized no escape rows; the fixture is vacuous")
	}

	// And the handoff to the kernel carries the same guarantee, which is
	// what makes the loader's duty and the kernel's predicate one invariant.
	if err := m.KernelTable().CheckValid(); err != nil {
		t.Errorf("the normalized model breaches the kernel predicate: %v", err)
	}
}

// REQ-15: "Ownership: **this RDR is the single normative home of the
// escape-row shape rule.**"
// ADVERSARIAL
//
// Single home means one rule, two enforcers that never disagree in the
// direction that matters: whatever the AUTHORED layer accepts, the KERNEL
// predicate must also accept. (The reverse does not hold — REQ-13 makes the
// authored layer strictly stricter.) Asserted over every authored escape
// shape this file exercises.
func TestReq15_TheAuthoredLayerNeverAcceptsWhatTheKernelRejects(t *testing.T) {
	fragments := map[string]string{
		"conforming":        conformingEscapeRule,
		"write_block":       escapeWithWriteBlock,
		"clear_list":        escapeWithClearList,
		"empty_write_block": escapeWithEmptyWriteBlock,
		"empty_clear_list":  escapeWithEmptyClearList,
	}

	for name, fragment := range fragments {
		t.Run(name, func(t *testing.T) {
			m, err := loadEscapeShape(t, fragment)
			if err != nil {
				// Refused by the authored layer; the kernel is never reached.
				return
			}
			if err := m.KernelTable().CheckValid(); err != nil {
				t.Errorf("the authored layer accepted a document the kernel "+
					"predicate rejects (%v); one rule, two enforcers that "+
					"cannot disagree in this direction", err)
			}
		})
	}
}

// --- K. Phase 3 — shared normalizer conformance fixtures -----------------

// REQ-97: "Encode the authored-path half as a named, shared fixture set
// binding the future RDR 0002 build"
// HAPPY PATH
//
// The fixture set is NAMED and SHARED — one table of cases with a stable
// name per case and a stable expectation per case, exercised here against the
// real loader so the set is executable rather than inert data. Each case's
// name is the binding a successor build cites.
func TestReq97_TheSharedEscapeConformanceFixtureSetBindsTheAuthoredPath(t *testing.T) {
	for _, fx := range EscapeShapeConformanceFixtures() {
		t.Run(fx.Name, func(t *testing.T) {
			m, err := loadEscapeShape(t, fx.Rule)

			if fx.WantCategory == "" {
				if err != nil {
					t.Fatalf("fixture %q must load clean; refused: %v",
						fx.Name, err)
				}
				if err := m.KernelTable().CheckValid(); err != nil {
					t.Errorf("fixture %q loaded a table the kernel "+
						"predicate rejects: %v", fx.Name, err)
				}
				return
			}

			if err == nil {
				t.Fatalf("fixture %q loaded clean; want category %q",
					fx.Name, fx.WantCategory)
			}
			cat, ok := table.CategoryOf(err)
			if !ok {
				t.Fatalf("fixture %q refused with no category: %v",
					fx.Name, err)
			}
			if cat != fx.WantCategory {
				t.Errorf("fixture %q: category = %q; want %q",
					fx.Name, cat, fx.WantCategory)
			}
		})
	}

	// The set must actually cover the three malformed shapes the record
	// names plus a conforming control; a set that lost a case would still
	// pass the loop above.
	names := map[string]bool{}
	for _, fx := range EscapeShapeConformanceFixtures() {
		names[fx.Name] = true
	}
	for _, want := range []string{
		"escape-with-write-block",
		"escape-with-clear-list",
		"escape-with-empty-write-block",
		"escape-conforming",
		"authored-clear-canonical",
	} {
		if !names[want] {
			t.Errorf("the shared fixture set is missing the named case %q",
				want)
		}
	}
}

// REQ-98: "one canonical authored-clear case pins the `<clear>`-sentinel-write
// representation identically for the kernel suite and the normalizer suite, so
// the two enforcers of one invariant cannot drift."
// BOUNDARY
//
// The canonical case is an ORDINARY rule carrying a clear (an escape rule may
// not carry one at all). It pins the representation at both ends of the
// handoff: the normalizer renders the clear as a `<clear>`-sentinel WRITE,
// and the kernel row carries that same sentinel value — which is exactly why
// the kernel's Writes-only predicate catches a clear (REQ-2).
func TestReq98_TheCanonicalAuthoredClearPinsTheSentinelRepresentation(t *testing.T) {
	var canonical EscapeShapeFixture
	for _, fx := range EscapeShapeConformanceFixtures() {
		if fx.Name == "authored-clear-canonical" {
			canonical = fx
		}
	}
	if canonical.Name == "" {
		t.Fatalf("the shared fixture set carries no canonical authored-clear " +
			"case")
	}

	// The expectation both halves below assert against is carried by the
	// fixture itself, so the normalizer half and the kernel half cannot
	// drift onto two independently written literals.
	if canonical.ClearRuleID == "" || canonical.ClearKey == "" {
		t.Fatalf("the canonical authored-clear fixture carries no pinned "+
			"expectation (rule %q, key %q); the two halves would then have "+
			"no shared source of truth",
			canonical.ClearRuleID, canonical.ClearKey)
	}

	m, err := loadEscapeShape(t, canonical.Rule)
	if err != nil {
		t.Fatalf("the canonical authored-clear fixture was refused: %v", err)
	}

	// Normalizer side: on the identified row, the clear is rendered as the
	// row's SOLE write, on the cleared key, whose value is exactly the
	// one-member `<clear>` sentinel. A prefixed, multi-member, differently
	// keyed, or differently placed rendering is a different representation.
	row := onlyRowWithRuleID(t, m.Rows,
		func(r table.Row) string { return r.RuleID },
		canonical.ClearRuleID)
	if len(row.Writes) != 1 {
		t.Fatalf("normalized row %q renders %d writes; the authored clear "+
			"overwrites the rule's write block and is the row's sole write",
			canonical.ClearRuleID, len(row.Writes))
	}
	if row.Writes[0].Key != canonical.ClearKey {
		t.Errorf("normalized row %q renders its clear on key %q; want %q",
			canonical.ClearRuleID, row.Writes[0].Key, canonical.ClearKey)
	}
	if want := []string{table.ClearSentinel}; !slices.Equal(row.Writes[0].Value, want) {
		t.Errorf("normalized row %q renders %q = %q; want exactly %q — the "+
			"two enforcers read ONE representation, not merely a value that "+
			"contains the sentinel",
			canonical.ClearRuleID, row.Writes[0].Key, row.Writes[0].Value, want)
	}

	// Kernel side: the SAME representation survives the handoff, which is
	// what makes the kernel's Writes-only predicate cover clears (REQ-2).
	// `resolve.Tag.Value` is one string, so the sentinel must be the WHOLE
	// value — a value merely containing it would leave the predicate
	// covering something other than a clear.
	kernelRow := onlyRowWithRuleID(t, m.KernelTable().Rows,
		func(r resolve.Row) string { return r.RuleID },
		canonical.ClearRuleID)
	if len(kernelRow.Writes) != 1 {
		t.Fatalf("kernel row %q carries %d writes; want the single clear "+
			"write", canonical.ClearRuleID, len(kernelRow.Writes))
	}
	if kernelRow.Writes[0].Key != canonical.ClearKey {
		t.Errorf("kernel row %q carries its clear on key %q; want %q",
			canonical.ClearRuleID, kernelRow.Writes[0].Key, canonical.ClearKey)
	}
	if kernelRow.Writes[0].Value != table.ClearSentinel {
		t.Errorf("kernel row %q carries %q = %q; want exactly %q — the "+
			"kernel suite and the normalizer suite would otherwise pin "+
			"different representations of one invariant",
			canonical.ClearRuleID, kernelRow.Writes[0].Key,
			kernelRow.Writes[0].Value, table.ClearSentinel)
	}
}

// onlyRowWithRuleID returns the single row in rows whose rule id is want,
// failing the test if none or more than one matches. Selecting by identity
// is what keeps the REQ-98 expectation pinned to ONE row rather than to
// "somewhere in the table".
func onlyRowWithRuleID[R any](t *testing.T, rows []R, ruleID func(R) string, want string) R {
	t.Helper()

	var found []R
	for _, row := range rows {
		if ruleID(row) == want {
			found = append(found, row)
		}
	}
	switch len(found) {
	case 1:
		return found[0]
	case 0:
		t.Fatalf("no row carries the rule id %q; the fixture no longer "+
			"normalizes to the row the expectation names", want)
	default:
		t.Fatalf("%d rows carry the rule id %q; the expectation names one "+
			"row", len(found), want)
	}
	panic("unreachable")
}

// --- the shared fixture set itself ---------------------------------------

// EscapeShapeFixture is one named case in the shared conformance set.
//
// WantCategory is "" for a case that must LOAD CLEAN; otherwise it is the
// validation category the load must refuse under.
type EscapeShapeFixture struct {
	Name         string
	Rule         string
	WantCategory table.Category
	// ClearRuleID and ClearKey carry the canonical authored-clear
	// expectation — which normalized row bears the `<clear>` sentinel, and
	// on which tag key. They are the ONE source of truth REQ-98 asks for:
	// the normalizer half and the kernel half both read them, so neither
	// can be relaxed without relaxing the other. Empty on every case but
	// the canonical authored clear.
	ClearRuleID string
	ClearKey    string
	// Why records what the case binds, so a successor build reading the set
	// does not have to re-derive the intent.
	Why string
}

// EscapeShapeConformanceFixtures is the named, shared fixture set REQ-97
// asks this RDR to ship: the authored-path half of the escape-row shape
// invariant, encoded once and exercised by both the normalizer suite (here)
// and, on the kernel side, by the conformance predicate.
func EscapeShapeConformanceFixtures() []EscapeShapeFixture {
	return []EscapeShapeFixture{
		{
			Name:         "escape-with-write-block",
			Rule:         escapeWithWriteBlock,
			WantCategory: table.CatMalformedEscapeDeclaration,
			Why:          "an escape rule carrying a write block (REQ-11)",
		},
		{
			Name:         "escape-with-clear-list",
			Rule:         escapeWithClearList,
			WantCategory: table.CatMalformedEscapeDeclaration,
			Why:          "an escape rule carrying a clear list (REQ-11)",
		},
		{
			Name:         "escape-with-empty-write-block",
			Rule:         escapeWithEmptyWriteBlock,
			WantCategory: table.CatMalformedEscapeDeclaration,
			Why: "rejection keys on the PRESENCE of the block, so an " +
				"empty one is rejected too (REQ-12)",
		},
		{
			Name: "escape-conforming",
			Rule: conformingEscapeRule,
			Why:  "a conforming escape rule normalizes write-free (REQ-14)",
		},
		{
			Name:        "authored-clear-canonical",
			ClearRuleID: "ordinary-clear",
			ClearKey:    "status",
			Rule: `
[[rule]]
id = "ordinary-clear"
clear = ["status"]
[rule.match.status]
eq = "final"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "draft"
`,
			Why: "the canonical authored clear, whose `<clear>`-sentinel " +
				"write representation the kernel suite and the normalizer " +
				"suite pin identically (REQ-98, REQ-2)",
		},
	}
}
