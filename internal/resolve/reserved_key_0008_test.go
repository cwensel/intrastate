package resolve_test

// RDR 0008 — ownership of the recognized-outcome tag key name, kernel half.
//
// Two surfaces live here. Block 1's reservation is behavior the kernel
// ALREADY carries (`assemble` binds `Input.Recognized` at the reserved key),
// so those tests pin shipped behavior against regression. Blocks 4, 5 and 6
// are net-new: the exported construction-time `Input` predicate and the
// `Resolve`-entry call that applies the same predicate.
//
// Nothing here mocks the unit under test. The kernel is always the real
// resolve.Resolve, and the only stub is RDR 0003's per-atom value seam,
// which RDR 0007 fences as a collaborator that never sees the view.

import (
	"errors"
	"testing"

	"github.com/cwensel/intrastate/internal/resolve"
)

// reservedKey is the literal this whole RDR is about. It is written out
// rather than read from a kernel constant on purpose: a test comparing two
// kernel-supplied spellings verifies the kernel against itself (premortem
// P-5, REQ-6).
const reservedKey = "recognized"

// --- helpers -------------------------------------------------------------

// recognizedMatchTable is a table whose single ordinary row matches on the
// reserved key directly. This is the shape scenario 1's kernel half needs:
// a row selected because the assembled view carries the recognized outcome
// under `recognized`, not because a literal was compared to a literal.
func recognizedMatchTable() resolve.Table {
	return resolve.Table{
		Revision: "rev-0008",
		Outcomes: []string{"round-clean", "reconcile-block"},
		Rows: []resolve.Row{{
			RuleID:        "rdr.round-clean",
			SourceLocator: "flows/rdr.toml:120",
			Outcome:       "round-clean",
			Match: []resolve.Tag{
				{Key: "status", Value: "Draft"},
				{Key: reservedKey, Value: "round-clean"},
			},
			RequiresOwned: []string{"status"},
			NextTags:      []resolve.Tag{{Key: "stage", Value: "prelock"}},
			Writes:        []resolve.Tag{{Key: "status", Value: "Final"}},
		}},
	}
}

// recognizedMatchInput carries a conforming tuple over that table.
func recognizedMatchInput() resolve.Input {
	return resolve.Input{
		Flow:       "rdr",
		Table:      recognizedMatchTable(),
		Owned:      []resolve.Tag{{Key: "status", Value: "Draft"}},
		Observed:   []resolve.Tag{{Key: "reviews", Value: "2"}},
		Recognized: "round-clean",
	}
}

// --- Block 1: the reserved key and the assembled view --------------------

// REQ-1: "When a resolve carries a freshly recognized outcome, the assembled
// evaluation view MUST bind it under exactly the tag key `recognized`
// (`internal/resolve/resolve.go::recognizedTagKey`)."
// HAPPY PATH
//
// Behavioral, per REQ-6: the assertion runs the real kernel and observes the
// binding through selection of a row whose Match names the literal key. If
// the kernel bound the outcome anywhere else, no candidate would match and
// the resolve would refuse rather than plan.
func TestReq1_AssembledViewBindsRecognizedOutcomeAtTheReservedKey(t *testing.T) {
	got, err := resolve.Resolve(recognizedMatchInput())
	if err != nil {
		t.Fatalf("conforming input traveled the Go error path: %v", err)
	}
	if got.Refusal != nil {
		t.Fatalf("row matching on %q refused %q; the outcome is not bound at "+
			"the reserved key", reservedKey, got.Refusal.Kind)
	}
	if got.Plan == nil {
		t.Fatal("no disposition")
	}
	if got.Plan.RuleID != "rdr.round-clean" {
		t.Errorf("selected rule = %q; want rdr.round-clean", got.Plan.RuleID)
	}
}

// REQ-2: "The key is a reserved kernel keyword: table authors conform to it
// and MUST NOT rebind it."
// ADVERSARIAL
//
// The kernel's half of the reservation: a caller-supplied OBSERVED tag on the
// reserved key must not be able to stand in for the recognized outcome. The
// producer obligation of block 4 forbids constructing this at all; here the
// point is the reservation's substance — a rebinding attempt does not get to
// redefine what the key means.
func TestReq2_ObservedTagCannotRebindTheReservedKey(t *testing.T) {
	in := recognizedMatchInput()
	in.Observed = append(in.Observed,
		resolve.Tag{Key: reservedKey, Value: "reconcile-block"})

	got, err := resolve.Resolve(in)
	if err == nil && got.Plan != nil && got.Plan.RuleID != "rdr.round-clean" {
		t.Fatalf("an observed tag on %q rebound the reserved key and selected %q",
			reservedKey, got.Plan.RuleID)
	}
	// Either the block-4 precondition refuses this input outright (the
	// conforming outcome once the predicate lands) or D3's precedence keeps
	// the recognized binding authoritative. What is forbidden is the observed
	// value silently becoming the recognized outcome.
	if err == nil && got.Refusal == nil && got.Plan == nil {
		t.Fatal("neither disposition nor error")
	}
}

// REQ-3: "The reservation holds unconditionally — the key is reserved whether
// or not a given resolve binds it — while the binding obligation is scoped to
// resolves that carry an outcome: `assemble` injects only for a non-empty
// `Input.Recognized`, so an absent outcome yields a view with no `recognized`
// key."
// BOUNDARY
//
// Observed behaviorally: with an empty `Input.Recognized` a row whose Match
// names the reserved key cannot match, because the key is absent from the
// view. A kernel that injected an empty binding would let the row match on
// the empty string.
func TestReq3_AbsentOutcomeYieldsAViewWithNoRecognizedKey(t *testing.T) {
	in := recognizedMatchInput()
	in.Recognized = ""
	in.Table.Outcomes = append(in.Table.Outcomes, "")
	in.Table.Rows[0].Outcome = ""
	in.Table.Rows[0].Match = []resolve.Tag{
		{Key: "status", Value: "Draft"},
		{Key: reservedKey, Value: ""},
	}

	got, err := resolve.Resolve(in)
	if err != nil {
		t.Fatalf("Go error path: %v", err)
	}
	if got.Plan != nil {
		t.Fatalf("a row matching %q = \"\" was selected with no outcome carried; "+
			"assemble injected an empty binding", reservedKey)
	}
	if got.Refusal == nil || got.Refusal.Kind != resolve.KindNoMatch {
		t.Errorf("disposition = %+v; want no_match — the key is absent from the view", got)
	}
}

// REQ-4: "The conformance test's input domain is non-empty outcomes."
// REQ-5: "The binding obligation above is scoped so that such a resolve is
// outside it, not in breach of it — an empty-string outcome is degenerate for
// reasons that are not this RDR's to rule on."
// DOMAIN EDGE
//
// A NEGATIVE requirement: this RDR must NOT close the empty-string-outcome
// path. The assertion is that an empty outcome still travels the value path —
// a modeled disposition with a nil error — rather than acquiring a refusal or
// a Go error this RDR would have had to invent.
func TestReq4And5_EmptyOutcomeIsNotClosedByThisRDR(t *testing.T) {
	in := recognizedMatchInput()
	in.Recognized = ""

	got, err := resolve.Resolve(in)
	if err != nil {
		t.Fatalf("an empty outcome acquired a Go error path this RDR does not "+
			"rule on: %v", err)
	}
	if got.Refusal == nil && got.Plan == nil {
		t.Fatal("empty outcome yielded neither disposition")
	}
	if got.Refusal != nil && got.Refusal.Kind != resolve.KindUnmodeledOutcome {
		t.Errorf("refusal kind = %q; want the pre-existing unmodeled_outcome — "+
			"this RDR mints no refusal for the degenerate outcome", got.Refusal.Kind)
	}
}

// REQ-7: "Whether the normalizer additionally references an exported kernel
// constant for the *spelling* is implementation latitude the implementer may
// resolve; this RDR requires only the behavioral test"
// DOMAIN EDGE
//
// The latitude is real, but the two spellings must AGREE. `table.RecognizedTagKey`
// and the kernel's unexported `recognizedTagKey` being two constants is
// conforming; them diverging is not. Asserted behaviorally: the table
// package's exported spelling, used as a Match key, selects the row the
// kernel bound.
func TestReq7_TableSideSpellingAgreesWithTheKernelBinding(t *testing.T) {
	in := recognizedMatchInput()
	in.Table.Rows[0].Match = []resolve.Tag{
		{Key: "status", Value: "Draft"},
		{Key: tableRecognizedTagKey, Value: "round-clean"},
	}

	got, err := resolve.Resolve(in)
	if err != nil {
		t.Fatalf("Go error path: %v", err)
	}
	if got.Plan == nil {
		t.Fatalf("the table package's spelling %q did not select the kernel-bound "+
			"row; disposition %+v", tableRecognizedTagKey, got)
	}
}

// REQ-9: "record the reserved-keyword rule where implementers read it — this
// RDR's Normative Contracts as authority, plus a pointer comment on
// `internal/resolve/resolve.go::recognizedTagKey` citing RDR 0008. No kernel
// behavior change."
// DOMAIN EDGE
//
// Phase 1's deliverable is a documentation pointer at a named site. The
// artifact is the comment on that constant, and the assertion is that it
// cites RDR 0008 — without a citation the reserved-keyword rule is invisible
// to the next implementer reading the constant.
func TestReq9_RecognizedTagKeyConstantCitesRDR0008(t *testing.T) {
	doc := constDoc(t, "recognizedTagKey")
	if doc == "" {
		t.Fatal("recognizedTagKey carries no doc comment")
	}
	if !mentionsRDR0008(doc) {
		t.Errorf("recognizedTagKey's comment does not cite RDR 0008:\n%s", doc)
	}
}

// REQ-6: "the conformance test MUST run the real kernel and assert the
// assembled view binds the recognized outcome under the normalizer's spelling
// — a test comparing two hardcoded literals verifies the doc against itself
// and is non-conforming (premortem P-5)."
// REQ-8: "That latitude is about the spelling constant only — it is distinct
// from the `Input` predicate the Enforcement locus settles, which this RDR does
// require the kernel to export."
// ADVERSARIAL
//
// REQ-6 constrains the FORM of the conformance test, so its own witness is
// structural: the conformance assertions must run the real kernel, not compare
// constants. Asserted by demonstrating the kernel is load-bearing — mutating
// only the view (the outcome the tuple carries) flips the disposition, which a
// two-literal comparison could never do. REQ-8's half is that the spelling
// latitude does not extend to the predicate: the predicate is required to be
// exported regardless of how the spelling is spelled.
func TestReq6And8_ConformanceIsBehavioralAndThePredicateIsStillRequired(t *testing.T) {
	// The kernel is load-bearing: same table, different carried outcome,
	// different disposition. Two hardcoded literals cannot produce this.
	fires, err := resolve.Resolve(recognizedMatchInput())
	if err != nil {
		t.Fatalf("Go error path: %v", err)
	}
	if fires.Plan == nil {
		t.Fatal("the conforming tuple did not fire")
	}

	other := recognizedMatchInput()
	other.Recognized = "reconcile-block"
	other.Table.Rows[0].Outcome = "reconcile-block"
	missed, err := resolve.Resolve(other)
	if err != nil {
		t.Fatalf("Go error path: %v", err)
	}
	if missed.Plan != nil {
		t.Error("changing only the carried outcome did not change the " +
			"disposition; the assertion is not reading the assembled view")
	}

	// REQ-8: the spelling latitude does not excuse the exported predicate.
	if len(exportedInputPredicates(t)) == 0 {
		t.Error("the kernel exports no Input predicate; REQ-8 says the spelling " +
			"latitude is distinct from the predicate this RDR requires")
	}
}

// --- Block 4: the exported Input predicate -------------------------------

// REQ-66: "**Chosen: (b) and (c) together — one predicate, two call sites.**"
// REQ-68: "the exported predicate is the single definition, and `Resolve`
// calls it at entry."
// ADVERSARIAL
//
// A NEGATIVE requirement: no second, independently-written check. Asserted
// behaviorally — the two call sites must agree on EVERY input, conforming and
// breaching alike. Two independent checks would diverge on some input; one
// definition called twice cannot.
func TestReq66And68_TheTwoCallSitesAgreeOnEveryInput(t *testing.T) {
	inputs := map[string]resolve.Input{
		"conforming":      recognizedMatchInput(),
		"empty Input{}":   {},
		"legal fixture":   legalInput(),
		"owned breach":    breaching(func(in *resolve.Input) { in.Owned = append(in.Owned, resolve.Tag{Key: reservedKey, Value: "x"}) }),
		"observed breach": breaching(func(in *resolve.Input) { in.Observed = append(in.Observed, resolve.Tag{Key: reservedKey, Value: "x"}) }),
		"RequiresOwned breach": breaching(func(in *resolve.Input) {
			in.Table.Rows[0].RequiresOwned = append(in.Table.Rows[0].RequiresOwned, reservedKey)
		}),
		"breach with no outcome": breaching(func(in *resolve.Input) {
			in.Recognized = ""
			in.Owned = append(in.Owned, resolve.Tag{Key: reservedKey, Value: "x"})
		}),
		"reserved key twice": breaching(func(in *resolve.Input) {
			in.Observed = append(in.Observed, resolve.Tag{Key: reservedKey, Value: "a"}, resolve.Tag{Key: reservedKey, Value: "b"})
		}),
		"reserved key in Match": breaching(func(in *resolve.Input) {
			in.Table.Rows[0].Match = []resolve.Tag{{Key: reservedKey, Value: "round-clean"}}
		}),
	}

	for name, in := range inputs {
		t.Run(name, func(t *testing.T) {
			viaPredicate := resolve.CheckInput(in) != nil
			_, entryErr := resolve.Resolve(in)
			viaEntry := entryErr != nil

			if viaPredicate != viaEntry {
				t.Errorf("the predicate says breach=%v and Resolve's entry says "+
					"breach=%v; one definition called at two sites cannot "+
					"disagree", viaPredicate, viaEntry)
			}
		})
	}
}

// breaching builds an Input from the conforming baseline with one mutation.
func breaching(mutate func(*resolve.Input)) resolve.Input {
	in := recognizedMatchInput()
	mutate(&in)
	return in
}

// REQ-36: "Producers of kernel `Input` MUST NOT supply an owned or observed
// tag keyed `recognized`; the reserved key enters the assembled view only
// through `Input.Recognized`."
// REQ-38: "Enforcement is one predicate at two call sites: the kernel MUST
// export a construction-time predicate over `Input` that producers may call,
// and `Resolve` MUST apply that same predicate at entry, returning a non-nil
// error and no `Result` disposition on breach."
// ADVERSARIAL
//
// The net-new core. Both call sites are asserted for each breach channel, so
// a predicate that exists but is never called at entry — or an entry check
// written independently of the predicate — fails here.
func TestReq36And38_ReservedKeyedProducerTagBreachesAtBothCallSites(t *testing.T) {
	cases := map[string]func(*resolve.Input){
		"owned tag on the reserved key": func(in *resolve.Input) {
			in.Owned = append(in.Owned, resolve.Tag{Key: reservedKey, Value: "x"})
		},
		"observed tag on the reserved key": func(in *resolve.Input) {
			in.Observed = append(in.Observed, resolve.Tag{Key: reservedKey, Value: "x"})
		},
	}

	for name, breach := range cases {
		t.Run(name, func(t *testing.T) {
			in := recognizedMatchInput()
			breach(&in)

			if err := resolve.CheckInput(in); err == nil {
				t.Error("the exported predicate accepted a breaching Input")
			}

			got, err := resolve.Resolve(in)
			if err == nil {
				t.Fatal("Resolve did not apply the predicate at entry")
			}
			if got.Plan != nil || got.Refusal != nil {
				t.Errorf("Resolve returned a disposition alongside the breach error: %+v", got)
			}
		})
	}
}

// REQ-37: "A breach is a producer programmer mistake — not table data, and
// never a new `RefusalKind` — and travels the Go error path RDR 0001 reserves
// for programmer mistakes."
// DOMAIN EDGE
//
// The breach must not leak into the modeled vocabulary: no new kind, and no
// refusal value at all.
func TestReq37_BreachIsAGoErrorAndMintsNoRefusalKind(t *testing.T) {
	in := recognizedMatchInput()
	in.Observed = append(in.Observed, resolve.Tag{Key: reservedKey, Value: "x"})

	got, err := resolve.Resolve(in)
	if err == nil {
		t.Fatal("no Go error for a breaching Input")
	}
	if got.Refusal != nil {
		t.Fatalf("the breach was modeled as refusal %q; block 4 forbids a "+
			"refusal for a producer mistake", got.Refusal.Kind)
	}
	for _, kind := range resolve.RefusalKinds() {
		if err.Error() == string(kind) {
			t.Errorf("the breach error spells a RefusalKind (%q)", kind)
		}
	}
}

// REQ-39: "That predicate carries **both** reserved-key obligations — the
// owned/observed tag keys of this block, and block 5's `Row.RequiresOwned`
// reservation — so the reserved name has exactly one enforcement point across
// every channel it can arrive through."
// BOUNDARY
//
// All three channels through one predicate. A predicate covering only the tag
// sequences leaves the RequiresOwned channel unguarded and fails the third leg.
func TestReq39_OnePredicateCarriesAllThreeReservedKeyChannels(t *testing.T) {
	channels := map[string]func(*resolve.Input){
		"Owned": func(in *resolve.Input) {
			in.Owned = append(in.Owned, resolve.Tag{Key: reservedKey, Value: "x"})
		},
		"Observed": func(in *resolve.Input) {
			in.Observed = append(in.Observed, resolve.Tag{Key: reservedKey, Value: "x"})
		},
		"Row.RequiresOwned": func(in *resolve.Input) {
			in.Table.Rows[0].RequiresOwned =
				append(in.Table.Rows[0].RequiresOwned, reservedKey)
		},
	}

	for name, breach := range channels {
		t.Run(name, func(t *testing.T) {
			in := recognizedMatchInput()
			breach(&in)
			if err := resolve.CheckInput(in); err == nil {
				t.Errorf("the predicate does not cover the %s channel", name)
			}
		})
	}
}

// REQ-41: "The predicate's read-domain … it reads `Input.Owned`,
// `Input.Observed`, and each row's `RequiresOwned` reached through
// `Input.Table.Rows` — three sequences, no other `Input` field."
// BOUNDARY
//
// Read-domain is asserted from the outside: varying every OTHER `Input` field
// while holding the three sequences conforming must never produce a breach,
// and varying them while holding a breach in place must never clear it.
func TestReq41_PredicateReadDomainIsExactlyTheThreeSequences(t *testing.T) {
	outside := map[string]func(*resolve.Input){
		"Flow":              func(in *resolve.Input) { in.Flow = reservedKey },
		"Recognized":        func(in *resolve.Input) { in.Recognized = "reconcile-block" },
		"Table.Revision":    func(in *resolve.Input) { in.Table.Revision = reservedKey },
		"Table.Outcomes":    func(in *resolve.Input) { in.Table.Outcomes = []string{reservedKey} },
		"Row.Match":         func(in *resolve.Input) { in.Table.Rows[0].Match = []resolve.Tag{{Key: reservedKey, Value: "x"}} },
		"Row.Writes":        func(in *resolve.Input) { in.Table.Rows[0].Writes = []resolve.Tag{{Key: reservedKey, Value: "x"}} },
		"Row.NextTags":      func(in *resolve.Input) { in.Table.Rows[0].NextTags = []resolve.Tag{{Key: reservedKey, Value: "x"}} },
		"Row.Guard":         func(in *resolve.Input) { in.Table.Rows[0].Guard = []resolve.GuardAtom{allAtom(reservedKey, opEq, "x")} },
		"Row.RuleID":        func(in *resolve.Input) { in.Table.Rows[0].RuleID = reservedKey },
		"Row.SourceLocator": func(in *resolve.Input) { in.Table.Rows[0].SourceLocator = reservedKey },
		"Row.Outcome":       func(in *resolve.Input) { in.Table.Rows[0].Outcome = reservedKey },
	}

	for name, mutate := range outside {
		t.Run("conforming/"+name, func(t *testing.T) {
			in := recognizedMatchInput()
			mutate(&in)
			if err := resolve.CheckInput(in); err != nil {
				t.Errorf("the predicate read %s, which is outside its declared "+
					"read-domain: %v", name, err)
			}
		})
		t.Run("breaching/"+name, func(t *testing.T) {
			in := recognizedMatchInput()
			in.Owned = append(in.Owned, resolve.Tag{Key: reservedKey, Value: "x"})
			mutate(&in)
			if err := resolve.CheckInput(in); err == nil {
				t.Errorf("mutating %s cleared a breach on Input.Owned", name)
			}
		})
	}
}

// REQ-42: "`Owned` and `Observed` are **sequences of key/value tags, not
// maps** … so the check is a scan and a duplicate reserved key is admissible
// input rather than a structural impossibility."
// INPUT EDGE
//
// The same reserved key twice in one sequence is a constructible input, not a
// map collision. It must breach.
func TestReq42_DuplicateReservedKeyInOneSequenceIsAdmissibleAndBreaches(t *testing.T) {
	in := recognizedMatchInput()
	in.Observed = append(in.Observed,
		resolve.Tag{Key: reservedKey, Value: "first"},
		resolve.Tag{Key: reservedKey, Value: "second"},
	)
	if len(in.Observed) != 3 {
		t.Fatalf("Observed collapsed to %d entries; it is not a sequence", len(in.Observed))
	}
	if err := resolve.CheckInput(in); err == nil {
		t.Error("a doubled reserved key did not breach")
	}
}

// REQ-43: "On multiplicity, the predicate reports **the first breach it finds
// and returns a single non-nil error**; it does not aggregate, and the order
// in which it scans the three sequences is implementation latitude."
// ADVERSARIAL
//
// Every channel breaching at once still yields ONE error. Scan order is
// latitude, so nothing here asserts which breach is named — only that the
// return is a single non-nil error and not a join of several.
func TestReq43_MultipleBreachesYieldOneErrorNotAnAggregate(t *testing.T) {
	in := recognizedMatchInput()
	in.Owned = append(in.Owned, resolve.Tag{Key: reservedKey, Value: "a"})
	in.Observed = append(in.Observed, resolve.Tag{Key: reservedKey, Value: "b"})
	in.Table.Rows[0].RequiresOwned = append(in.Table.Rows[0].RequiresOwned, reservedKey)

	err := resolve.CheckInput(in)
	if err == nil {
		t.Fatal("a triply-breaching Input did not breach")
	}
	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		t.Errorf("the predicate aggregated %d errors; block 4 requires a single "+
			"non-nil error", len(joined.Unwrap()))
	}
}

// REQ-44: "What this RDR requires of the implementation is therefore only
// that its own breach be **detected whenever present** — never skipped
// because another precondition also fired."
// REQ-45: "the implementable rule is: **apply this predicate at entry and
// report its breach; if the table also breaches RDR 0009's shape rule, either
// error is conforming.**"
// REQ-90: TS-6's doubly-breaching variant.
// ADVERSARIAL
//
// A table that ALSO carries RDR 0009's escape-row shape breach (non-empty
// Escape with non-empty Writes) must still be refused. JDR 0001 §JD-5 is
// open, so which error is reported is deliberately unasserted — only that the
// input does not resolve.
func TestReq44And45And90_ReservedKeyBreachIsNotSkippedByACoincidentShapeBreach(t *testing.T) {
	in := recognizedMatchInput()
	in.Observed = append(in.Observed, resolve.Tag{Key: reservedKey, Value: "x"})
	in.Table.Rows = append(in.Table.Rows, resolve.Row{
		RuleID:        "rdr.escape",
		SourceLocator: "flows/rdr.toml:200",
		Outcome:       "round-clean",
		Escape:        []resolve.RefusalKind{resolve.KindNoMatch},
		Writes:        []resolve.Tag{{Key: "status", Value: "Final"}},
	})

	if err := resolve.CheckInput(in); err == nil {
		t.Error("the reserved-key breach was skipped on a doubly-breaching table")
	}
	got, err := resolve.Resolve(in)
	if err == nil {
		t.Fatal("Resolve admitted a doubly-breaching input")
	}
	if got.Plan != nil || got.Refusal != nil {
		t.Errorf("Resolve returned a disposition alongside the error: %+v", got)
	}
}

// REQ-46: "The check is unconditional on `Input.Recognized`: a reserved-keyed
// owned or observed tag is a breach whether or not the resolve carries an
// outcome, matching block 1's unconditional reservation."
// REQ-89: TS-6's empty-`Input.Recognized` variant "is expected to breach on
// the reserved key alone".
// BOUNDARY
//
// The distinguishing pair: an empty `Input.Recognized` carrying the reserved
// key breaches, while the empty `Input{}` carrying no reserved key does not.
func TestReq46And89_ReservedKeyBreachesEvenWithNoRecognizedOutcome(t *testing.T) {
	in := recognizedMatchInput()
	in.Recognized = ""
	in.Owned = append(in.Owned, resolve.Tag{Key: reservedKey, Value: "x"})

	if err := resolve.CheckInput(in); err == nil {
		t.Error("the predicate conditioned the reserved-key check on a non-empty " +
			"Input.Recognized")
	}
	if _, err := resolve.Resolve(in); err == nil {
		t.Error("Resolve admitted the breach because no outcome was carried")
	}

	if err := resolve.CheckInput(resolve.Input{}); err != nil {
		t.Errorf("the empty Input{} — which carries no reserved key at all — "+
			"breached: %v", err)
	}
}

// REQ-40: "The obligation is NOT inherited from RDR 0009 … the predicate's
// exported name MUST NOT be bound to any symbol name from RDR 0009."
// REQ-70: "The predicate's exported name is implementation latitude — do
// **not** bind it to any symbol name from RDR 0009 (`Final`)".
// REQ-69: "Surface conceded, stated plainly: one exported predicate over
// `Input` plus one kernel-entry call."
// ADVERSARIAL
//
// The scope ceiling and the name ban together. `internal/resolve` gains
// exactly one new exported function-shaped symbol taking an `Input`, and it is
// not named `Final`.
func TestReq40And69And70_ExactlyOneNewExportedInputPredicateNotNamedFinal(t *testing.T) {
	preds := exportedInputPredicates(t)
	if len(preds) == 0 {
		t.Fatal("internal/resolve exports no predicate over Input")
	}
	if len(preds) > 1 {
		t.Errorf("internal/resolve exports %d predicates over Input (%v); REQ-69 "+
			"concedes exactly one", len(preds), preds)
	}
	for _, name := range preds {
		if name == "Final" {
			t.Errorf("the predicate is bound to %q, a symbol name RDR 0009 owns", name)
		}
	}
}

// --- Block 5: Row.RequiresOwned name reservation -------------------------

// REQ-47: "A normalized row's `Row.RequiresOwned` MUST NOT name `recognized`."
// REQ-50: "the **exported `Input` predicate of block 4 MUST also reject a
// `Row.RequiresOwned` entry naming `recognized`** on the rows of the supplied
// table, on the same Go error path and with the same producer-breach
// semantics."
// REQ-97: TS-9's second half — "Both are asserted; the pair is what makes
// block 5 falsifiable rather than discharged-by-assertion."
// ADVERSARIAL
func TestReq47And50And97_RequiresOwnedNamingTheReservedKeyBreachesAtBothSites(t *testing.T) {
	in := recognizedMatchInput()
	in.Table.Rows[0].RequiresOwned = []string{"status", reservedKey}

	if err := resolve.CheckInput(in); err == nil {
		t.Error("the predicate accepted a row naming the reserved key in RequiresOwned")
	}

	got, err := resolve.Resolve(in)
	if err == nil {
		t.Fatal("Resolve admitted a row naming the reserved key in RequiresOwned")
	}
	if got.Refusal != nil {
		t.Errorf("Resolve refused %q instead of returning the producer-breach "+
			"error; block 5's whole point is that the entry precondition "+
			"outranks the owned_state_unavailable residual", got.Refusal.Kind)
	}
	if got.Plan != nil {
		t.Errorf("Resolve planned alongside the breach error: %+v", *got.Plan)
	}
}

// REQ-51: "The residual this closes is kernel-side: a row carrying
// `RequiresOwned: [\"recognized\"]` finds the key present under
// `ProvenanceRecognized`, fails the owned-only test, and yields
// `owned_state_unavailable` naming the reserved key … It stays a documented
// residual on that narrowed path"
// REQ-96: TS-9's first half.
// REQ-48: "This is a name reservation only: what `RequiresOwned` *means* is
// owned by RDR 0007 … this RDR cites that rule rather than restating it"
// DOMAIN EDGE
//
// The residual must remain REACHABLE on the package-internal path that
// bypasses the entry precondition. An implementation that "fixed" the
// owned-only test — say by treating ProvenanceRecognized as satisfying an
// owned requirement — would alter RDR 0007's rule, which REQ-48 forbids.
// Reached through the in-package test path, not through Resolve.
func TestReq48And51And96_OwnedOnlyResidualRemainsReachableOnTheInternalPath(t *testing.T) {
	kind, missing := resolveBypassingPrecondition(t, recognizedRequiresOwnedInput())

	if kind != resolve.KindOwnedStateUnavailable {
		t.Fatalf("the internal path yielded %q; the documented residual is "+
			"owned_state_unavailable", kind)
	}
	if !containsString(missing, reservedKey) {
		t.Errorf("MissingOwned = %v; want it to name %q", missing, reservedKey)
	}
}

// REQ-52: "The `RequiresOwned` reservation needs no source-lint check here
// (A7, Verified) — the field has no authored source form"
// REQ-98: "There is **no source-lint half** (A7, Verified) … the corresponding
// negative test belongs to RDR 0002's `write to non-owned tag` suite, not this
// RDR's."
// DOMAIN EDGE
//
// A NEGATIVE requirement. The reservation must not be re-homed into the
// declaration channel: no `reserved_tag_key` failure may be minted for a
// RequiresOwned entry (REQ-49). Asserted at the seam that would carry it —
// the predicate's error is a plain Go error, not a categorized table failure.
func TestReq49And52And98_RequiresOwnedBreachIsNotADataLevelCategory(t *testing.T) {
	in := recognizedMatchInput()
	in.Table.Rows[0].RequiresOwned = []string{reservedKey}

	err := resolve.CheckInput(in)
	if err == nil {
		t.Fatal("no breach")
	}
	if carriesLoadCategory(err) {
		t.Errorf("the RequiresOwned breach carries a data-level load category; "+
			"REQ-49 says it does not share the declaration channel: %v", err)
	}
}

// --- Block 6: unchanged kernel disposition -------------------------------

// REQ-53: "Kernel disposition is unchanged by this RDR for conforming input:
// no new `RefusalKind`, no change to RDR 0001 D3's provenance precedence
// (`owned` > `observed` > `recognized`), and no behavior change for any
// conforming table and conforming input."
// BOUNDARY
func TestReq53_NoNewRefusalKindAndTheClosedSetIsUnchanged(t *testing.T) {
	want := []resolve.RefusalKind{
		resolve.KindNoMatch,
		resolve.KindAmbiguousMatch,
		resolve.KindOwnedStateUnavailable,
		resolve.KindGuardUnevaluable,
		resolve.KindUnmodeledOutcome,
	}
	got := resolve.RefusalKinds()
	if len(got) != len(want) {
		t.Fatalf("RefusalKinds() has %d members; RDR 0008 mints none, so the set "+
			"is still the frozen five: %v", len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("RefusalKinds()[%d] = %q; want %q", i, got[i], want[i])
		}
	}
}

// REQ-54: "D3 remains the deterministic backstop behind the producer
// obligation for input that bypasses the precondition."
// REQ-94: TS-8 — "D3's precedence (`owned` > `observed` > `recognized`)
// resolves the collision deterministically — the same input yields the same
// disposition on repeat runs."
// ADVERSARIAL
//
// A non-conforming Input whose owned and observed tags BOTH carry the reserved
// key, resolved through the package-internal path that bypasses the
// precondition. Precedence must decide, and it must decide the same way twice.
func TestReq54And94_D3PrecedenceDeterministicallyResolvesTheCollision(t *testing.T) {
	in := collidingReservedKeyInput()

	first := assembledReservedBinding(t, in)
	second := assembledReservedBinding(t, in)

	if first != second {
		t.Fatalf("the same non-conforming input assembled differently on replay: "+
			"%+v then %+v", first, second)
	}
	if first.provenance != resolve.ProvenanceOwned {
		t.Errorf("provenance = %v; D3 ranks owned above observed above recognized, "+
			"so the owned tag wins the collision", first.provenance)
	}
	if first.value != "owned-wins" {
		t.Errorf("value = %q; want the owned tag's value", first.value)
	}
}

// REQ-95: TS-8's "Scope note: the unreachability half of A5 is **not** a
// test."
// DOMAIN EDGE
//
// A NEGATIVE requirement recorded as an executable guard against
// over-implementation: the collision path must stay reachable in-package.
// This test fails if a future change makes the colliding input unconstructible
// or the internal path unreachable, which is what "documented residual"
// forbids.
func TestReq95_TheCollisionPathStaysConstructibleInPackage(t *testing.T) {
	in := collidingReservedKeyInput()
	if len(in.Owned) == 0 || len(in.Observed) == 0 {
		t.Fatal("the colliding fixture lost one of its two channels")
	}
	// Reached without going through Resolve: the precondition is bypassed by
	// construction, which is exactly what makes the residual documented rather
	// than eliminated.
	if _, err := resolve.Resolve(in); err == nil {
		t.Error("Resolve admitted the colliding input; the entry precondition " +
			"must refuse it while the in-package path still reaches assemble")
	}
}

// REQ-55: "Breaching input gains a non-nil Go error where it previously
// resolved — a programmer-mistake path, not a disposition change for any
// conforming caller."
// REQ-67: "`internal/resolve/resolve.go::Resolve` already returns `(Result,
// error)` … so (b) adds a check without widening the signature."
// REQ-71: "Breach of the precondition returns a non-nil error and no `Result`
// disposition; conforming callers see no behavior change."
// BOUNDARY
func TestReq55And67And71_ResolveSignatureIsUnchangedAndBreachYieldsNoDisposition(t *testing.T) {
	if !resolveSignatureIsResultError(t) {
		t.Error("Resolve's signature is no longer func(Input) (Result, error)")
	}

	in := recognizedMatchInput()
	in.Owned = append(in.Owned, resolve.Tag{Key: reservedKey, Value: "x"})
	got, err := resolve.Resolve(in)
	if err == nil {
		t.Fatal("no error for breaching input")
	}
	if (got != resolve.Result{}) {
		t.Errorf("Resolve returned a non-zero Result alongside the breach: %+v", got)
	}
}

// REQ-88: "A conforming `Input` returns a nil error, pinning that the check
// adds no behavior change for conforming callers — including the empty
// `Input{}` that `internal/resolve/resolve_test.go:748` already pins as a
// nil-error resolve, which the new precondition MUST keep green."
// REQ-99: Phase 3's conforming-`Input` nil-error pin.
// HAPPY PATH
func TestReq88And99_ConformingInputsIncludingTheEmptyTupleStayNilError(t *testing.T) {
	conforming := map[string]resolve.Input{
		"empty Input{}":                 {},
		"canonical legal input":         legalInput(),
		"row matching the reserved key": recognizedMatchInput(),
		"no outcome, no reserved key": func() resolve.Input {
			in := recognizedMatchInput()
			in.Recognized = ""
			return in
		}(),
	}

	for name, in := range conforming {
		t.Run(name, func(t *testing.T) {
			if err := resolve.CheckInput(in); err != nil {
				t.Errorf("the predicate rejected a conforming Input: %v", err)
			}
			got, err := resolve.Resolve(in)
			if err != nil {
				t.Fatalf("the new precondition broke a conforming resolve: %v", err)
			}
			if got.Plan == nil && got.Refusal == nil {
				t.Error("neither disposition")
			}
		})
	}
}

// REQ-86: TS-6 — "it pins block 4's first-breach rule — the predicate returns
// one non-nil error, not two, and the test asserts a single error rather than
// an aggregate."
// REQ-87: TS-6 — "**Expected**: all three variants breach. `Resolve` returns a
// non-nil `error` and a zero-valued `Result` (no `Plan`, no `Refusal`) —
// never a new `RefusalKind`, and never a modeled disposition. The exported
// predicate, called directly on the same `Input`, reports the same breach: one
// predicate, two call sites, asserted at both."
// REQ-85: TS-6's three variants.
// ADVERSARIAL
func TestReq85And86And87_AllThreeTS6VariantsBreachAtBothSitesWithZeroResult(t *testing.T) {
	variants := map[string]resolve.Input{
		"non-empty Input.Recognized": func() resolve.Input {
			in := recognizedMatchInput()
			in.Observed = append(in.Observed, resolve.Tag{Key: reservedKey, Value: "x"})
			return in
		}(),
		"empty Input.Recognized": func() resolve.Input {
			in := recognizedMatchInput()
			in.Recognized = ""
			in.Owned = append(in.Owned, resolve.Tag{Key: reservedKey, Value: "x"})
			return in
		}(),
		"reserved key twice in one sequence": func() resolve.Input {
			in := recognizedMatchInput()
			in.Observed = append(in.Observed,
				resolve.Tag{Key: reservedKey, Value: "a"},
				resolve.Tag{Key: reservedKey, Value: "b"})
			return in
		}(),
	}

	for name, in := range variants {
		t.Run(name, func(t *testing.T) {
			perr := resolve.CheckInput(in)
			if perr == nil {
				t.Fatal("the exported predicate did not report the breach")
			}
			if joined, ok := perr.(interface{ Unwrap() []error }); ok {
				t.Errorf("the predicate aggregated %d errors; one is required",
					len(joined.Unwrap()))
			}

			got, rerr := resolve.Resolve(in)
			if rerr == nil {
				t.Fatal("Resolve did not report the breach at entry")
			}
			if (got != resolve.Result{}) {
				t.Errorf("Result is not zero-valued: %+v", got)
			}
			if got.Refusal != nil {
				t.Errorf("the breach was modeled as refusal %q", got.Refusal.Kind)
			}
		})
	}
}

// --- helpers used above --------------------------------------------------

func containsString(in []string, want string) bool {
	for _, s := range in {
		if s == want {
			return true
		}
	}
	return false
}

// carriesLoadCategory reports whether err carries a data-level load-failure
// category in its chain. RDR 0002's categorized refusal is any error in the
// chain exposing a category discriminator; the predicate's producer-breach
// error must expose none, because REQ-49 keeps the RequiresOwned reservation
// off the declaration channel. `internal/resolve` cannot import
// `internal/table` (the dependency runs the other way), so the probe is
// structural: a category-bearing error answers one of these shapes.
func carriesLoadCategory(err error) bool {
	var byMethod interface{ Category() string }
	if errors.As(err, &byMethod) {
		return true
	}
	var byField interface{ LoadCategory() string }
	return errors.As(err, &byField)
}
