package resolve_test

// The MATCH leg of the duplicate-key positionality defect (kata `qdnj`).
//
// Phase 3b's ADV-3 found that a tag key supplied more than once WITHIN one
// provenance with differing values made the disposition a function of
// SLICE POSITION rather than of the input tuple, which `0007:C8` and RDR
// 0001 REQ-1 forbid. `assemble` was fixed to mark such a key CONFLICTED,
// and deviation D9 routed the guard path through `TagSet.conflicting`.
//
// `Row.Match` selection is a SECOND consumer of the same root cause that
// D9's ruling never had in view: `TagSet.matches` compared the retained
// last-merged value and never consulted the conflicted mark, so the
// candidate filter stayed positional even after the guard path was fixed.
// These tests pin the match leg closed and fence the two narrowings D9
// declared deliberate.
//
// Every fixture drives the real kernel; no row verdict is injected.

import (
	"testing"

	"github.com/cwensel/intrastate/internal/resolve"
)

// dupKey is the key these fixtures duplicate within the observed
// provenance. A13 leaves caller-supplied Observed unconstrained, so a
// caller may legally hand the kernel both values.
const dupKey = "gate"

// matchOnGateTable models one guard-free ordinary edge whose Match names
// `gate`. With no guard atoms the row's disposition is decided entirely by
// the match filter, so these fixtures isolate `TagSet.matches` from the
// guard path D9 already closed.
func matchOnGateTable(rows ...resolve.Row) resolve.Table {
	return resolve.Table{
		Revision: "rev-match-conflicted",
		Outcomes: []string{"successful", "failed"},
		Rows:     rows,
	}
}

// gateMatchRow is the guard-free candidate: it matches only when the view
// carries gate=closed.
func gateMatchRow() resolve.Row {
	return resolve.Row{
		RuleID:        "A",
		SourceLocator: "flows/rdr.toml:300",
		Outcome:       "successful",
		Match: []resolve.Tag{
			{Key: "status", Value: "Draft"},
			{Key: dupKey, Value: "closed"},
		},
		RequiresOwned: []string{"status"},
		NextTags:      []resolve.Tag{{Key: "status", Value: "Final"}},
		Writes:        []resolve.Tag{{Key: "status", Value: "Final"}},
	}
}

// gateInput builds the tuple over the canonical fixture view, appending
// the given observed `gate` values in the order supplied.
func gateInput(table resolve.Table, values ...string) resolve.Input {
	in := resolve.Input{
		Flow:       "rdr",
		Table:      table,
		Owned:      []resolve.Tag{{Key: "status", Value: "Draft"}},
		Observed:   []resolve.Tag{{Key: "reviews", Value: "2"}},
		Recognized: "successful",
		Guards:     eqSeam{},
	}
	for _, v := range values {
		in.Observed = append(in.Observed, resolve.Tag{Key: dupKey, Value: v})
	}
	return in
}

// describe renders a disposition compactly for permutation diffs.
func describe(r resolve.Result) string {
	if r.Refused() {
		return "refusal " + string(r.Refusal.Kind)
	}
	return "plan " + r.Plan.RuleID + " (escaped=" + boolWord(r.Plan.Escaped) + ")"
}

func boolWord(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// assertSameDisposition fails unless the two orderings of one tag multiset
// produce the identical disposition.
func assertSameDisposition(t *testing.T, forward, reversed resolve.Result) {
	t.Helper()
	if describe(forward) != describe(reversed) {
		t.Fatalf("[gate=open, gate=closed] yields %s but [gate=closed, gate=open] "+
			"yields %s. The two inputs carry the SAME tag multiset and differ only in "+
			"slice position, so RDR 0001 REQ-1 and `0007:C8` make them one input tuple "+
			"with one disposition. `assemble` marks the key CONFLICTED; the match "+
			"filter must consume that mark the way the guard path does (D9)",
			describe(forward), describe(reversed))
	}
}

// TestAdv0007_3_DuplicateKeyMakesRowSelectionAFunctionOfSlicePosition is
// the match leg of ADV-3. The guard leg is closed at
// `guard_adversarial_0007_test.go`; this one runs a table with NO guard
// atoms at all, so only `TagSet.matches` can decide the row.
func TestAdv0007_3_DuplicateKeyMakesRowSelectionAFunctionOfSlicePosition(t *testing.T) {
	table := matchOnGateTable(gateMatchRow())

	forward := mustResolve(t, gateInput(table, "open", "closed"))
	reversed := mustResolve(t, gateInput(table, "closed", "open"))

	assertSameDisposition(t, forward, reversed)

	// Agreement alone is not the contract: a fix that made BOTH orderings
	// MATCH would agree and would still have invented a value for a key
	// the view carries no single value for. A conflicted key must fail the
	// match the way it fails a value-comparing atom — conservatively, into
	// a refusal that can only turn a plan into a refusal and never the
	// reverse.
	for name, got := range map[string]resolve.Result{
		"open first":   forward,
		"closed first": reversed,
	} {
		t.Run(name, func(t *testing.T) {
			if !got.Refused() {
				t.Fatalf("kernel emitted %s for a row whose Match names a CONFLICTED "+
					"key: the view carries no single value for %q, so comparing the "+
					"last-merged one picks a value by slice position. `0007:C8` wants "+
					"the undecidable input to be refusal-class, never a false allow",
					describe(got), dupKey)
			}
			if got.Refusal.Kind != resolve.KindNoMatch {
				t.Errorf("refusal kind = %q; want %q — a conflicted key is non-matching, "+
					"which is an ordinary zero-match and mints no new refusal kind (REQ-7 "+
					"pins the set at five)", got.Refusal.Kind, resolve.KindNoMatch)
			}
		})
	}
}

// TestAdv0007_3_ConflictedKeyIsNotEscapableByNamingIt closes the
// escapability question the kata carried. `escapeOrRefuse` filters escape
// candidates through the SAME `matches()` call as ordinary candidates, so
// an escape row that NAMES the conflicted key must fail the match too. If
// it did not, an operator could rescue the positional plan through the
// kernel's most permissive path and reintroduce the defect one layer down.
func TestAdv0007_3_ConflictedKeyIsNotEscapableByNamingIt(t *testing.T) {
	escape := resolve.Row{
		RuleID:        "ESCAPE-NAMES-GATE",
		SourceLocator: "flows/rdr.toml:301",
		Outcome:       "successful",
		Match: []resolve.Tag{
			{Key: "status", Value: "Draft"},
			{Key: dupKey, Value: "open"},
		},
		RequiresOwned: nil,
		NextTags:      []resolve.Tag{{Key: "status", Value: "Blocked"}},
		Writes:        nil,
		Escape:        []resolve.RefusalKind{resolve.KindNoMatch},
	}
	table := matchOnGateTable(gateMatchRow(), escape)

	forward := mustResolve(t, gateInput(table, "open", "closed"))
	reversed := mustResolve(t, gateInput(table, "closed", "open"))

	assertSameDisposition(t, forward, reversed)

	for name, got := range map[string]resolve.Result{
		"open first":   forward,
		"closed first": reversed,
	} {
		t.Run(name, func(t *testing.T) {
			if !got.Refused() {
				t.Fatalf("an escape row NAMING the conflicted key rescued with %s: the "+
					"escape path runs the same match filter, so naming a key the view "+
					"carries no single value for must not select the escape either",
					describe(got))
			}
			if got.Refusal.Kind != resolve.KindNoMatch {
				t.Errorf("refusal kind = %q; want %q", got.Refusal.Kind, resolve.KindNoMatch)
			}
		})
	}
}

// TestAdv0007_3_EscapeNotNamingTheConflictedKeyStillRescues is the
// positive control for the test above. Pruning conflicted-key matches must
// prune ONLY those: an escape edge that does not name the conflicted key
// is the modeled reviewable escape RDR 0002 intends, and it must still
// rescue — identically under both orderings.
func TestAdv0007_3_EscapeNotNamingTheConflictedKeyStillRescues(t *testing.T) {
	escape := conformingEscapeRow("ESCAPE-CLEAN", "flows/rdr.toml:302", resolve.KindNoMatch)
	table := matchOnGateTable(gateMatchRow(), escape)

	forward := mustResolve(t, gateInput(table, "open", "closed"))
	reversed := mustResolve(t, gateInput(table, "closed", "open"))

	assertSameDisposition(t, forward, reversed)

	for name, got := range map[string]resolve.Result{
		"open first":   forward,
		"closed first": reversed,
	} {
		t.Run(name, func(t *testing.T) {
			if got.Refused() {
				t.Fatalf("the modeled escape did not rescue (refusal %q): its Match names "+
					"only `status`, which is not conflicted, so treating a conflicted key "+
					"as non-matching must not reach it", got.Refusal.Kind)
			}
			if got.Plan.RuleID != "ESCAPE-CLEAN" || !got.Plan.Escaped {
				t.Errorf("rescued with %s; want the clean escape row with Escaped=true",
					describe(got))
			}
		})
	}
}

// TestAdv0007_3_MatchNarrowingsHold fences the two narrowings D9 declared
// deliberate, plus RDR 0001 REQ-3, against the match filter. A fix that
// pruned more than conflicted keys would break one of these.
func TestAdv0007_3_MatchNarrowingsHold(t *testing.T) {
	table := matchOnGateTable(gateMatchRow())

	t.Run("a repeat carrying the SAME value still matches", func(t *testing.T) {
		// D9's first narrowing: either resolution is the same value, so the
		// result is already a function of the tuple and no conflict exists.
		plan := mustPlan(t, gateInput(table, "closed", "closed"))
		if plan.RuleID != "A" {
			t.Errorf("selected rule %q; want A", plan.RuleID)
		}
	})

	t.Run("a key crossing provenances still matches by precedence", func(t *testing.T) {
		// D9's second narrowing: owned over observed is a property of the
		// tuple, not of slice order, and is untouched. Observed says `open`,
		// owned says `closed`; owned wins and the row matches under BOTH
		// observed orderings relative to the owned tag.
		for name, observed := range map[string][]string{
			"observed open": {"open"},
			"observed both": {"open", "open"},
		} {
			t.Run(name, func(t *testing.T) {
				in := gateInput(table, observed...)
				in.Owned = append(in.Owned, resolve.Tag{Key: dupKey, Value: "closed"})
				plan := mustPlan(t, in)
				if plan.RuleID != "A" {
					t.Errorf("selected rule %q; want A — owned `gate=closed` overrides the "+
						"observed value by precedence", plan.RuleID)
				}
			})
		}
	})

	t.Run("a permuted non-duplicate input selects identically", func(t *testing.T) {
		// RDR 0001 REQ-3: with no key repeated, permuting the observed slice
		// assembles an identical view and must select the identical row.
		base := gateInput(table, "closed")
		base.Observed = append(base.Observed, resolve.Tag{Key: "iterations", Value: "3"})

		permuted := gateInput(table)
		permuted.Observed = []resolve.Tag{
			{Key: "iterations", Value: "3"},
			{Key: dupKey, Value: "closed"},
			{Key: "reviews", Value: "2"},
		}

		if got, want := describe(mustResolve(t, permuted)), describe(mustResolve(t, base)); got != want {
			t.Errorf("permuting a non-duplicate observed slice changed the disposition: "+
				"%s vs %s", got, want)
		}
	})
}
