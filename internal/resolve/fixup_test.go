package resolve_test

// Phase 3c regression coverage for RDR 0001.
//
// Phase 3a's chain-of-verification widened two of Phase 3b's adversarial
// findings with sub-cases the adversarial suite did not construct. Those
// sub-cases are pinned here so the uniform viability gate cannot regress
// on the paths that were only ever reached by a spec-derived probe:
//
//   - FAIL-1 sub-case 1d: a guarded escape edge with a NIL guard seam. The
//     bypass was not specific to a live seam, so a fix that only routed
//     escape rows through an existing evaluator would still have leaked.
//   - FAIL-1 sub-case 1e: an ambiguous-class escape edge. The bypass was
//     not specific to the no_match class, so a fix applied to one rescue
//     class would still have leaked on the other.
//   - FAIL-3 third surface: the ambiguous_match Refusal.Rows order. ADV-3
//     recorded only the guard_unevaluable and owned_state_unavailable
//     payloads; rowRefs feeds this one too.
//
// Anchored in RDR 0001 Trade-offs / Failure Modes ("Silent failure would
// mean the kernel guessed a transition or executed persistence directly";
// "Diagnosis starts with the input tuple, the refusal kind, and the
// transition table revision used for that resolution") read with REQ-5,
// REQ-12, REQ-15, REQ-23 and REQ-1/REQ-2/REQ-10.

import (
	"reflect"
	"testing"

	"github.com/newcoinc/intrastate/internal/resolve"
)

// FAIL-1 sub-case 1d: a guarded escape edge must not rescue when the guard
// seam is absent entirely.
//
// An ordinary guarded row with a nil seam already refuses guard_unevaluable
// (pinned by TestReq33's "nil guard seam does not fall back to ambient
// evaluation"). The escape path must answer identically: a nil seam means
// the predicate is undecidable, and the kernel cannot know the edge holds.
// Emitting the escape plan anyway would be REQ-12's "guessed a transition"
// on the one path where no ordinary edge was available to check it.
func TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue(t *testing.T) {
	escape := escapeRow("rdr.escape.nilseam", "flows/rdr.toml:90", resolve.KindNoMatch)
	escape.Guard = "iterations >= 3"

	in := noMatchInput()
	in.Table.Revision = "rev-fixup-1d"
	in.Table.Rows = append(in.Table.Rows, escape)
	in.Guards = nil // no seam at all, not merely a seam that cannot decide

	got := mustResolve(t, in)

	if !got.Refused() {
		t.Fatalf("kernel emitted a plan for rule %q (escaped=%v, writes=%v); want "+
			"guard_unevaluable — the escape edge is guarded and no guard seam was "+
			"supplied, so the kernel cannot know the predicate holds",
			got.Plan.RuleID, got.Plan.Escaped, got.Plan.Writes)
	}
	if got.Refusal.Kind != resolve.KindGuardUnevaluable {
		t.Errorf("refusal kind = %q; want %q — an absent seam must refuse, not "+
			"let the escape edge through unevaluated",
			got.Refusal.Kind, resolve.KindGuardUnevaluable)
	}
	if got.Refusal.Guard != escape.Guard {
		t.Errorf("refusal names guard %q; want %q for diagnosis",
			got.Refusal.Guard, escape.Guard)
	}
}

// FAIL-1 sub-case 1e: the ambiguous_match rescue class must be gated too.
//
// The bypass was not specific to no_match. Here two ordinary rows match, so
// the condition is a genuine ambiguity, and the modeled ambiguous-class
// escape is both guard-FALSE and missing required owned state. Neither
// defect may be waived because the row happens to be an escape: the
// disposition must stay the underlying ambiguous_match refusal.
func TestFixup1e_AmbiguousClassEscapeIsGatedLikeAnyOtherCandidate(t *testing.T) {
	t.Run("guard FALSE must not rescue an ambiguity", func(t *testing.T) {
		escape := escapeRow("rdr.escape.ambiguous", "flows/rdr.toml:99", resolve.KindAmbiguousMatch)
		escape.Guard = "never"

		in := ambiguousInput()
		in.Table.Revision = "rev-fixup-1e-false"
		in.Table.Rows = append(in.Table.Rows, escape)
		in.Guards = fixtureGuards{decided: map[string]bool{"never": false}}

		got := mustResolve(t, in)

		if !got.Refused() {
			t.Fatalf("kernel emitted a plan for rule %q (escaped=%v); want an "+
				"ambiguous_match refusal — the escape edge's own guard was decided "+
				"FALSE, so selecting it is a guessed transition",
				got.Plan.RuleID, got.Plan.Escaped)
		}
		if got.Refusal.Kind != resolve.KindAmbiguousMatch {
			t.Errorf("refusal kind = %q; want %q", got.Refusal.Kind, resolve.KindAmbiguousMatch)
		}
	})

	t.Run("missing owned state must not rescue an ambiguity", func(t *testing.T) {
		escape := escapeRow("rdr.escape.ambiguous", "flows/rdr.toml:99", resolve.KindAmbiguousMatch)
		escape.RequiresOwned = []string{"never-present"}
		escape.Writes = []resolve.Tag{{Key: "status", Value: "Escaped"}}

		in := ambiguousInput()
		in.Table.Revision = "rev-fixup-1e-owned"
		in.Table.Rows = append(in.Table.Rows, escape)

		got := mustResolve(t, in)

		if !got.Refused() {
			t.Fatalf("kernel emitted a plan for rule %q describing writes %v; want "+
				"owned_state_unavailable — the ambiguous-class escape edge required "+
				"owned tag %q, which the accessor snapshot does not carry",
				got.Plan.RuleID, got.Plan.Writes, "never-present")
		}
		if got.Refusal.Kind != resolve.KindOwnedStateUnavailable {
			t.Errorf("refusal kind = %q; want %q",
				got.Refusal.Kind, resolve.KindOwnedStateUnavailable)
		}
	})

	t.Run("a viable ambiguous-class escape still rescues", func(t *testing.T) {
		// Control: the gate must prune and refuse, not disable the rescue.
		escape := escapeRow("rdr.escape.ambiguous", "flows/rdr.toml:99", resolve.KindAmbiguousMatch)
		escape.Guard = "always"

		in := ambiguousInput()
		in.Table.Revision = "rev-fixup-1e-control"
		in.Table.Rows = append(in.Table.Rows, escape)
		in.Guards = fixtureGuards{decided: map[string]bool{"always": true}}

		p := planOf(t, mustResolve(t, in))
		if p.RuleID != "rdr.escape.ambiguous" || !p.Escaped {
			t.Errorf("selected rule = %q (escaped=%v); want the modeled escape edge "+
				"— gating must reject unviable escapes, not every escape",
				p.RuleID, p.Escaped)
		}
	})
}

// FAIL-3 third surface: the ambiguous_match Refusal.Rows payload must not
// depend on Table.Rows order.
//
// REQ-2 identifies a resolution by the transition table revision, not by
// row sequence, so two orderings of the same row set at the same revision
// are the same input and must replay the same disposition (REQ-1). The
// conflicting-row list is the diagnosis RDR 0005 reports (REQ-10), so its
// order is part of the disposition, not a rendering detail.
func TestFixup3c_AmbiguousMatchRowsPayloadMustNotDependOnTableRowOrder(t *testing.T) {
	rowA := resolve.Row{
		RuleID:        "rdr.amb.a",
		SourceLocator: "flows/rdr.toml:10",
		Outcome:       "successful",
		Match:         []resolve.Tag{{Key: "status", Value: "Draft"}},
		RequiresOwned: []string{"status"},
		NextTags:      []resolve.Tag{{Key: "status", Value: "A"}},
		Writes:        []resolve.Tag{{Key: "status", Value: "A"}},
	}
	rowB := resolve.Row{
		RuleID:        "rdr.amb.b",
		SourceLocator: "flows/rdr.toml:20",
		Outcome:       "successful",
		Match:         []resolve.Tag{{Key: "status", Value: "Draft"}},
		RequiresOwned: []string{"status"},
		NextTags:      []resolve.Tag{{Key: "status", Value: "B"}},
		Writes:        []resolve.Tag{{Key: "status", Value: "B"}},
	}

	newInput := func(rows ...resolve.Row) resolve.Input {
		return resolve.Input{
			Flow: "rdr",
			Table: resolve.Table{
				Revision: "rev-fixup-3c",
				Outcomes: []string{"successful"},
				Rows:     rows,
			},
			Owned:      []resolve.Tag{{Key: "status", Value: "Draft"}},
			Recognized: "successful",
			Guards:     allGuardsTrue(),
		}
	}

	forward := mustResolve(t, newInput(rowA, rowB))
	reversed := mustResolve(t, newInput(rowB, rowA))

	if !forward.Refused() || !reversed.Refused() {
		t.Fatalf("both orderings must refuse; forward refused=%v reversed refused=%v",
			forward.Refused(), reversed.Refused())
	}
	if forward.Refusal.Kind != resolve.KindAmbiguousMatch {
		t.Fatalf("forward refusal kind = %q; want %q",
			forward.Refusal.Kind, resolve.KindAmbiguousMatch)
	}
	if len(forward.Refusal.Rows) != 2 {
		t.Fatalf("ambiguous refusal names %d rows; want both conflicting rows",
			len(forward.Refusal.Rows))
	}

	if !reflect.DeepEqual(forward.Refusal, reversed.Refusal) {
		t.Errorf("ambiguous_match payload depends on Table.Rows order.\n"+
			" rows [a,b] -> rows=%v\n"+
			" rows [b,a] -> rows=%v\n"+
			"Diagnosis (REQ-10) must be a function of the input tuple, not of "+
			"row order.",
			forward.Refusal.Rows, reversed.Refusal.Rows)
	}
}

// FAIL-3, escape surface: the degraded-ambiguity payload produced when two
// escape edges both match must be row-order-stable for the same reason.
// This path builds its Refusal.Rows from the escape candidate set rather
// than the ordinary one, so it is a distinct call site.
func TestFixup3c_DegradedEscapeAmbiguityPayloadMustNotDependOnRowOrder(t *testing.T) {
	escA := escapeRow("rdr.escape.a", "flows/rdr.toml:90", resolve.KindNoMatch)
	escB := escapeRow("rdr.escape.b", "flows/rdr.toml:95", resolve.KindNoMatch)

	newInput := func(rows ...resolve.Row) resolve.Input {
		in := noMatchInput()
		in.Table.Revision = "rev-fixup-3c-escape"
		in.Table.Rows = append(append([]resolve.Row(nil), in.Table.Rows...), rows...)
		return in
	}

	forward := mustResolve(t, newInput(escA, escB))
	reversed := mustResolve(t, newInput(escB, escA))

	forwardRefusal := refusalOf(t, forward)
	reversedRefusal := refusalOf(t, reversed)

	if forwardRefusal.Kind != resolve.KindAmbiguousMatch {
		t.Fatalf("forward refusal kind = %q; want %q — two matching escape edges "+
			"are not an exact-one rescue", forwardRefusal.Kind, resolve.KindAmbiguousMatch)
	}

	if !reflect.DeepEqual(forwardRefusal, reversedRefusal) {
		t.Errorf("degraded escape ambiguity payload depends on row order.\n"+
			" escapes [a,b] -> rows=%v\n"+
			" escapes [b,a] -> rows=%v",
			forwardRefusal.Rows, reversedRefusal.Rows)
	}
}

// The gate must be uniform, not merely present on both paths: an escape
// edge and an ordinary edge presented with the same blocking condition must
// receive the same refusal kind. This is the property that would have
// caught FAIL-1 and FAIL-2 together, and it is what a future refactor is
// most likely to break by re-introducing a path-specific shortcut.
func TestFixupGateIsUniformAcrossOrdinaryAndEscapeCandidates(t *testing.T) {
	cases := map[string]struct {
		mutate func(row *resolve.Row)
		guards resolve.GuardEvaluator
		want   resolve.RefusalKind
	}{
		"missing owned state": {
			mutate: func(row *resolve.Row) { row.RequiresOwned = []string{"never-present"} },
			guards: allGuardsTrue(),
			want:   resolve.KindOwnedStateUnavailable,
		},
		"undecidable guard": {
			mutate: func(row *resolve.Row) { row.Guard = "unknown-predicate" },
			guards: fixtureGuards{},
			want:   resolve.KindGuardUnevaluable,
		},
		"absent guard seam": {
			mutate: func(row *resolve.Row) { row.Guard = "any-predicate" },
			guards: nil,
			want:   resolve.KindGuardUnevaluable,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// Ordinary candidate: the sole row of a matching table.
			ordinary := legalInput()
			ordinary.Table.Revision = "rev-fixup-uniform-ordinary"
			ordinary.Guards = tc.guards
			row := ordinary.Table.Rows[0]
			tc.mutate(&row)
			ordinary.Table.Rows = []resolve.Row{row}

			// Escape candidate: the same blocking condition on an escape
			// edge rescuing a genuine no_match.
			escape := escapeRow("rdr.escape.uniform", "flows/rdr.toml:90", resolve.KindNoMatch)
			tc.mutate(&escape)
			escaped := noMatchInput()
			escaped.Table.Revision = "rev-fixup-uniform-escape"
			escaped.Guards = tc.guards
			escaped.Table.Rows = append(escaped.Table.Rows, escape)

			ordinaryKind := refusalOf(t, mustResolve(t, ordinary)).Kind
			escapeKind := refusalOf(t, mustResolve(t, escaped)).Kind

			if ordinaryKind != tc.want {
				t.Errorf("ordinary candidate refusal kind = %q; want %q",
					ordinaryKind, tc.want)
			}
			if escapeKind != tc.want {
				t.Errorf("escape candidate refusal kind = %q; want %q — the "+
					"viability gate must treat an escape edge exactly as it treats "+
					"an ordinary edge", escapeKind, tc.want)
			}
		})
	}
}
