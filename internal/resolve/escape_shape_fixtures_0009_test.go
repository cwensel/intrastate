package resolve_test

// RDR 0009 — shared fixture builders for the escape-row shape conformance
// suite.
//
// Nothing here mocks the unit under test. Every value below is a table a
// producer could hand the kernel; the kernel itself is always the real
// `resolve.Resolve` / `resolve.Table.CheckValid`.
//
// The normative fixture values (rule id `rdr.escape.needsowned`, locator
// `flows/rdr.toml:90`, the `status=Escaped` write) are read from the record's
// Minimum Viable Validation, not invented (REQ-71).

import (
	"errors"
	"testing"

	"github.com/newcoinc/intrastate/internal/resolve"
)

const (
	// breachRuleID and breachLocator are the normative fixture identity
	// REQ-72 fixes.
	breachRuleID  = "rdr.escape.needsowned"
	breachLocator = "flows/rdr.toml:90"
)

// breachWrites is the normative illegal write REQ-72 fixes.
func breachWrites() []resolve.Tag {
	return []resolve.Tag{{Key: "status", Value: "Escaped"}}
}

// breachingEscapeRow builds an escape row that BREACHES the conformance
// predicate: a non-empty Escape list alongside a non-empty Writes slice.
//
// It is otherwise a genuine escape candidate for the base fixture tag-set,
// so the row is reachable on the escape path and the breach is not an
// artifact of an unreachable row.
func breachingEscapeRow(ruleID, locator string, class resolve.RefusalKind) resolve.Row {
	row := escapeRow(ruleID, locator, class)
	row.Writes = breachWrites()
	return row
}

// conformedEscapeRow0009 builds the same escape row with the writes removed —
// the `escapeRow` builder's post-Phase-2 shape.
func conformedEscapeRow0009(ruleID, locator string, class resolve.RefusalKind) resolve.Row {
	row := escapeRow(ruleID, locator, class)
	row.Writes = nil
	return row
}

// breachingNoMatchInput is the scenario-1 table: a no-match table plus a
// breaching escape row for `no_match`, so the escape path SELECTS the
// offending row when the precondition is absent.
func breachingNoMatchInput() resolve.Input {
	in := noMatchInput()
	in.Table.Rows = append(in.Table.Rows,
		breachingEscapeRow(breachRuleID, breachLocator, resolve.KindNoMatch))
	return in
}

// conformingNoMatchInput is the same table with the writes removed: the
// baseline that must still resolve, escape, and refuse exactly as before.
func conformingNoMatchInput() resolve.Input {
	in := noMatchInput()
	in.Table.Rows = append(in.Table.Rows,
		conformedEscapeRow0009(breachRuleID, breachLocator, resolve.KindNoMatch))
	return in
}

// dormantBreachInput is the scenario-2 table: a table that resolves cleanly
// on its own, plus a breaching escape row NO resolution path reaches.
//
// The dormant row escapes `ambiguous_match`, a class this input never
// raises (the base table matches exactly one row), and it binds an outcome
// outside the request's reach, so neither the candidate partition nor the
// rescue phase visits it.
func dormantBreachInput() resolve.Input {
	in := legalInput()
	dormant := breachingEscapeRow(breachRuleID, breachLocator,
		resolve.KindAmbiguousMatch)
	dormant.Outcome = "never-recognized"
	dormant.Match = []resolve.Tag{{Key: "status", Value: "NeverHeldValue"}}
	in.Table.Rows = append(in.Table.Rows, dormant)
	return in
}

// emptyNotNilEscapeInput is the scenario-3 table: an escape row whose
// Writes is a NON-NIL, ZERO-LENGTH slice. The predicate tests emptiness,
// never nil-ness, so this table conforms.
func emptyNotNilEscapeInput() resolve.Input {
	in := noMatchInput()
	row := escapeRow(breachRuleID, breachLocator, resolve.KindNoMatch)
	row.Writes = []resolve.Tag{}
	in.Table.Rows = append(in.Table.Rows, row)
	return in
}

// multiBreachInput carries several breaching rows with the given identities,
// in the order supplied, on top of a table that would otherwise refuse
// `no_match`.
func multiBreachInput(refs ...resolve.RowRef) resolve.Input {
	in := noMatchInput()
	for _, ref := range refs {
		in.Table.Rows = append(in.Table.Rows,
			breachingEscapeRow(ref.RuleID, ref.SourceLocator,
				resolve.KindNoMatch))
	}
	return in
}

// --- error-shape readers -------------------------------------------------
//
// Every reader below reads STRUCTURE — never message prose. REQ-27, REQ-44,
// and REQ-89 forbid recovering a row identity or a count from formatted
// text, so no helper here calls Error() for anything but a failure report.

// breachElements traverses the aggregate exactly ONE level deep and returns
// its elements, failing the test if any element is not a
// *resolve.EscapeShapeBreachError (REQ-35, REQ-86).
func breachElements(t *testing.T, err error) []*resolve.EscapeShapeBreachError {
	t.Helper()

	if err == nil {
		t.Fatalf("no error to traverse; want a breach aggregate")
	}
	joined, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("the returned error does not expose Unwrap() []error; "+
			"REQ-34 fixes errors.Join as the uniform return shape (%T)", err)
	}
	var out []*resolve.EscapeShapeBreachError
	for i, elem := range joined.Unwrap() {
		var e *resolve.EscapeShapeBreachError
		if !errors.As(elem, &e) {
			t.Fatalf("aggregate element %d is %T; every element must be a "+
				"*EscapeShapeBreachError, never a nested join (REQ-35)", i, elem)
		}
		// One level deep: the element must not itself be a join.
		if _, nested := elem.(interface{ Unwrap() []error }); nested {
			t.Fatalf("aggregate element %d is itself a join; the traversal "+
				"must be EXACTLY one level deep (REQ-35)", i)
		}
		out = append(out, e)
	}
	return out
}

// breachRefs extracts the reported row identities in reported order.
func breachRefs(t *testing.T, err error) []resolve.RowRef {
	t.Helper()

	elems := breachElements(t, err)
	out := make([]resolve.RowRef, 0, len(elems))
	for _, e := range elems {
		out = append(out, e.Ref)
	}
	return out
}

// breachCounts extracts the per-identity Count values in reported order,
// read off the struct field — never parsed from prose (REQ-44, REQ-89).
func breachCounts(t *testing.T, err error) []int {
	t.Helper()

	elems := breachElements(t, err)
	out := make([]int, 0, len(elems))
	for _, e := range elems {
		out = append(out, e.Count)
	}
	return out
}

// sameRefs reports whether two identity sequences are equal, position by
// position. The comparison is over extracted []RowRef, never
// reflect.DeepEqual over two errors.Join values (REQ-83, REQ-86).
func sameRefs(a, b []resolve.RowRef) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// sameCounts reports whether two count sequences are equal position by
// position.
func sameCounts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
