package resolve_test

// RDR 0008 scenario 3 — A2's same-view property, written against RDR 0007's
// per-atom guard seam per deviation D1.
//
// The record's scenario 3 expected a view-capturing guard seam that would
// assert `Lookup("recognized")` and `Len()` on the guard side. RDR 0007
// (implemented, run 1) fenced `Evaluate(atom, value)` as a seam that "never
// sees the view", so that capture point does not exist and must not be added.
//
// D1's binding disposition moves the capture point without moving the
// contract: capture the VALUE handed to the guard atom over `recognized` and
// assert it equals `(in.Recognized, ProvenanceRecognized)`, while a `Match` on
// the same key selects the row in the SAME resolve. A divergence between what
// the guard seam is handed and what the matcher compares would fail one of the
// two assertions, which is the property A2 actually asserts.

import (
	"testing"

	"github.com/cwensel/intrastate/internal/resolve"
)

// REQ-76: TS-3 "a guard predicate and a match pattern both read the recognized
// tag in one resolve, through a **new** view-capturing guard seam added
// alongside `internal/resolve/fixtures_test.go::fixtureGuards` — not a change
// to it, since ~10 existing tests depend on its current shape."
// REQ-77: TS-3 "A row whose `Match` names the recognized tag is selected in the
// same resolve, so a divergence between the guard's view and the matcher's
// would fail one of the two assertions."
// REQ-78: TS-3 "The rows here are constructed directly: Final RDR 0002 mints no
// guard atom on `recognized` — it refuses one at load (A9) — so this is a
// kernel-plumbing test of A2's same-view property, not a shape any 0002-loaded
// table produces"
// ADVERSARIAL
//
// Two rows in one table, resolved once:
//
//   - the GUARDED row carries a guard atom over `recognized`, so the kernel
//     hands the seam the value the assembled view holds at that key;
//   - the MATCHED row's Match names `recognized` with the same expected value,
//     so it is selected only if the matcher reads the same binding.
//
// The guarded row's guard is decided FALSE, so it is pruned and the matched
// row is the exact-one selection. Both halves therefore run in the same
// resolve over the same view, and either half diverging fails.
//
// The seam is `valueSeam` from guard_fixtures_test.go — RDR 0007's per-atom
// stand-in with a call recorder. It is a NEW seam alongside `fixtureGuards`,
// not a change to it, which is what REQ-76 asks for.
func TestReq76And77And78_GuardSeamAndMatcherReadTheSameRecognizedBinding(t *testing.T) {
	const outcome = "round-clean"

	guarded := resolve.Row{
		RuleID:        "rdr.guarded",
		SourceLocator: "flows/rdr.toml:10",
		Outcome:       outcome,
		Match:         []resolve.Tag{{Key: "status", Value: "Draft"}},
		Guard: []resolve.GuardAtom{
			allAtom(reservedKey, opEq, "never-equal"),
		},
		NextTags: []resolve.Tag{{Key: "stage", Value: "guarded"}},
	}
	matched := resolve.Row{
		RuleID:        "rdr.matched",
		SourceLocator: "flows/rdr.toml:20",
		Outcome:       outcome,
		Match: []resolve.Tag{
			{Key: "status", Value: "Draft"},
			{Key: reservedKey, Value: outcome},
		},
		NextTags: []resolve.Tag{{Key: "stage", Value: "matched"}},
	}

	seam := newValueSeam()
	// The guarded row's atom is decided FALSE for the value the view actually
	// holds. Programming it against that exact value is what makes the seam's
	// verdict evidence of what it was handed: any other value reaches the
	// seam's default and answers GuardUnevaluable, which would refuse the
	// resolve rather than prune the row.
	seam.decide(allAtom(reservedKey, opEq, "never-equal"), outcome, resolve.GuardFalse)

	in := resolve.Input{
		Flow: "rdr",
		Table: resolve.Table{
			Revision: "rev-0008",
			Outcomes: []string{outcome},
			Rows:     []resolve.Row{guarded, matched},
		},
		Owned:      []resolve.Tag{{Key: "status", Value: "Draft"}},
		Recognized: outcome,
		Guards:     seam,
	}

	got, err := resolve.Resolve(in)
	if err != nil {
		t.Fatalf("conforming input traveled the Go error path: %v", err)
	}
	if got.Refusal != nil {
		t.Fatalf("resolve refused %q; both halves must run in one resolve",
			got.Refusal.Kind)
	}

	// Half 1 — the MATCHER read the recognized binding.
	if got.Plan == nil || got.Plan.RuleID != "rdr.matched" {
		t.Fatalf("selected %+v; the matcher did not read the recognized outcome "+
			"at %q", got.Plan, reservedKey)
	}

	// Half 2 — the GUARD SEAM was handed the same binding's value.
	var handed []string
	for _, call := range seam.seen() {
		if call.Key == reservedKey {
			handed = append(handed, call.Value)
		}
	}
	if len(handed) == 0 {
		t.Fatalf("the seam was never consulted for an atom over %q; calls: %+v",
			reservedKey, seam.seen())
	}
	for i, value := range handed {
		if value != in.Recognized {
			t.Errorf("seam call %d over %q was handed value %q; the assembled "+
				"view binds Input.Recognized = %q there — the guard's view and "+
				"the matcher's have diverged", i, reservedKey, value, in.Recognized)
		}
	}

	// Half 2b — the binding's PROVENANCE. `Evaluate` is fenced from the view,
	// so provenance is read on the package-internal path over the identical
	// Input, which is D1's "equals `(in.Recognized, ProvenanceRecognized)`".
	value, prov, ok := resolve.AssembledBindingForTest(in, reservedKey)
	if !ok {
		t.Fatalf("the assembled view carries no %q key", reservedKey)
	}
	if value != in.Recognized || prov != resolve.ProvenanceRecognized {
		t.Errorf("the assembled view binds %q = (%q, %v); want (%q, %v)",
			reservedKey, value, prov, in.Recognized, resolve.ProvenanceRecognized)
	}
}
