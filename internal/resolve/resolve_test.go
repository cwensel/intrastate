package resolve_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/newcoinc/intrastate/internal/resolve"
)

// mustResolve calls the kernel and fails the test if the Go error path was
// used. Every modeled disposition — plan or refusal — travels the Result
// value; the error return is reserved for programmer mistakes.
func mustResolve(t *testing.T, in resolve.Input) resolve.Result {
	t.Helper()
	got, err := resolve.Resolve(in)
	if err != nil {
		t.Fatalf("Resolve returned a Go error for a modeled disposition: %v", err)
	}
	return got
}

// refusalOf asserts the disposition is a refusal and returns it.
func refusalOf(t *testing.T, got resolve.Result) *resolve.Refusal {
	t.Helper()
	if got.Plan != nil {
		t.Fatalf("expected a refusal disposition; got plan %+v", *got.Plan)
	}
	if got.Refusal == nil {
		t.Fatal("expected a refusal disposition; got neither plan nor refusal")
	}
	return got.Refusal
}

// planOf asserts the disposition is a plan and returns it.
func planOf(t *testing.T, got resolve.Result) *resolve.Plan {
	t.Helper()
	if got.Refusal != nil {
		t.Fatalf("expected a transition plan; got refusal kind %q", got.Refusal.Kind)
	}
	if got.Plan == nil {
		t.Fatal("expected a transition plan; got neither plan nor refusal")
	}
	return got.Plan
}

// ---------------------------------------------------------------------
// Determinism and disposition shape
// ---------------------------------------------------------------------

// REQ-1: "Given the same flow identity, transition table revision,
// accessor-produced owned tag snapshot, caller-supplied observed tags, and
// freshly recognized outcome tag, resolve returns the same disposition:
// exactly one transition plan or exactly one typed refusal."
// HAPPY PATH
func TestReq1_SameInputTupleReturnsSameDisposition(t *testing.T) {
	first := mustResolve(t, legalInput())
	second := mustResolve(t, legalInput())

	if !reflect.DeepEqual(first, second) {
		t.Errorf("disposition differs across replay:\nfirst  = %+v\nsecond = %+v", first, second)
	}
	planOf(t, first)
}

// REQ-1: "resolve returns the same disposition: exactly one transition plan
// or exactly one typed refusal."
// BOUNDARY
//
// "Exactly one" is a cardinality claim on the disposition itself: a Result
// must never carry both a plan and a refusal, and never neither. Checked
// across the legal case and all five refusal conditions.
func TestReq1_DispositionIsExactlyOneOfPlanOrRefusal(t *testing.T) {
	for name, in := range allDispositionInputs() {
		t.Run(name, func(t *testing.T) {
			got := mustResolve(t, in)
			switch {
			case got.Plan != nil && got.Refusal != nil:
				t.Errorf("disposition carries both a plan and a refusal: %+v", got)
			case got.Plan == nil && got.Refusal == nil:
				t.Error("disposition carries neither a plan nor a refusal")
			}
			if (got.Refusal != nil) != got.Refused() {
				t.Errorf("Refused() = %v; want %v", got.Refused(), got.Refusal != nil)
			}
		})
	}
}

// REQ-2: "a resolution input is the tuple of flow identity, transition table
// revision, accessor-produced owned tag snapshot, observed tag-set, and
// freshly recognized outcome tag. Replaying that tuple must replay the
// disposition."
// HAPPY PATH
//
// Each named tuple member is load-bearing: perturbing any one of them must
// be able to change the disposition, otherwise it is not part of the
// identity the RDR claims.
func TestReq2_EveryTupleMemberParticipatesInTheDisposition(t *testing.T) {
	base := legalInput()
	basePlan := planOf(t, mustResolve(t, base))

	t.Run("owned snapshot", func(t *testing.T) {
		in := legalInput()
		in.Owned = []resolve.Tag{{Key: "status", Value: "Final"}}
		got := mustResolve(t, in)
		if got.Plan != nil && reflect.DeepEqual(*got.Plan, *basePlan) {
			t.Error("changing the owned snapshot left the disposition unchanged")
		}
	})

	t.Run("recognized outcome", func(t *testing.T) {
		in := legalInput()
		in.Recognized = "failed"
		got := mustResolve(t, in)
		if got.Plan != nil && reflect.DeepEqual(*got.Plan, *basePlan) {
			t.Error("changing the recognized outcome left the disposition unchanged")
		}
	})

	t.Run("table revision travels with the disposition", func(t *testing.T) {
		in := legalInput()
		in.Table.Revision = "rev-2"
		p := planOf(t, mustResolve(t, in))
		if p.Revision != "rev-2" {
			t.Errorf("plan revision = %q; want %q", p.Revision, "rev-2")
		}
	})

	t.Run("flow identity travels with a refusal", func(t *testing.T) {
		in := legalInput()
		in.Flow = "kata"
		in.Table = noMatchTable()
		r := refusalOf(t, mustResolve(t, in))
		if r.Flow != "kata" {
			t.Errorf("refusal flow = %q; want %q", r.Flow, "kata")
		}
	})
}

// REQ-3: "this RDR claims value-level replay determinism, not byte-identical
// output, hashes, or serialized canonical form. Determinism is covered by A1
// and the MVV replay test over the explicit input tuple."
// DOMAIN EDGE
//
// Value-level determinism must hold even when the caller presents the same
// tag-set in a different slice order: order of the input slices is not part
// of the value identity of a tag-set.
func TestReq3_ValueLevelReplayIsIndependentOfInputSliceOrder(t *testing.T) {
	a := legalInput()
	a.Observed = []resolve.Tag{{Key: "reviews", Value: "2"}, {Key: "lane", Value: "fast"}}

	b := legalInput()
	b.Observed = []resolve.Tag{{Key: "lane", Value: "fast"}, {Key: "reviews", Value: "2"}}

	gotA := mustResolve(t, a)
	gotB := mustResolve(t, b)

	// Anchor: both orderings must select the real edge. Comparing two
	// empty dispositions would prove nothing.
	planOf(t, gotA)
	planOf(t, gotB)

	if !reflect.DeepEqual(gotA, gotB) {
		t.Error("disposition depends on observed-tag slice order; determinism is value-level, not positional")
	}
}

// REQ-4: "identical table, owned snapshot, observed tags, and recognized
// outcome must return value-identical transition plans without hidden reads."
// HAPPY PATH
//
// "Without hidden reads" is asserted structurally: the kernel is given no
// ambient collaborator, and the caller's input values are not mutated by the
// call, so the second call sees exactly what the first saw.
func TestReq4_ReplayIsValueIdenticalAndDoesNotMutateCallerInput(t *testing.T) {
	in := legalInput()
	snapshot := deepCopyInput(in)

	first := planOf(t, mustResolve(t, in))
	second := planOf(t, mustResolve(t, in))

	if !reflect.DeepEqual(*first, *second) {
		t.Errorf("plans are not value-identical:\nfirst  = %+v\nsecond = %+v", *first, *second)
	}
	if !reflect.DeepEqual(in, snapshot) {
		t.Error("Resolve mutated its caller-supplied input; replay safety requires the input tuple to survive the call")
	}
}

// ---------------------------------------------------------------------
// Refusal obligations
// ---------------------------------------------------------------------

// REQ-5: "The kernel MUST refuse instead of guessing when no edge matches,
// more than one edge matches, required owned state is unavailable, a guard
// cannot be evaluated, or the recognized outcome is not modeled by the table."
// ADVERSARIAL
func TestReq5_EachRefusalConditionRefusesInsteadOfGuessing(t *testing.T) {
	cases := map[string]struct {
		in   resolve.Input
		kind resolve.RefusalKind
	}{
		"no edge matches":              {noMatchInput(), resolve.KindNoMatch},
		"more than one edge matches":   {ambiguousInput(), resolve.KindAmbiguousMatch},
		"required owned state missing": {missingOwnedInput(), resolve.KindOwnedStateUnavailable},
		"guard cannot be evaluated":    {unevaluableGuardInput(), resolve.KindGuardUnevaluable},
		"outcome not modeled":          {unmodeledOutcomeInput(), resolve.KindUnmodeledOutcome},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			r := refusalOf(t, mustResolve(t, tc.in))
			if r.Kind != tc.kind {
				t.Errorf("refusal kind = %q; want %q", r.Kind, tc.kind)
			}
		})
	}
}

// REQ-6: "Modeled refusal is a value-level resolver disposition, not a CLI
// error and not the Go error path for parser bugs, IO failures, or programmer
// mistakes."
// ADVERSARIAL
func TestReq6_ModeledRefusalsDoNotUseTheGoErrorPath(t *testing.T) {
	for name, in := range allRefusalInputs() {
		t.Run(name, func(t *testing.T) {
			got, err := resolve.Resolve(in)
			if err != nil {
				t.Fatalf("modeled refusal traveled the Go error path: %v", err)
			}
			if got.Refusal == nil {
				t.Fatal("modeled refusal did not travel the Result value")
			}
		})
	}
}

// REQ-7: "The kernel-owned refusal kind set is exactly:" no_match,
// ambiguous_match, owned_state_unavailable, guard_unevaluable,
// unmodeled_outcome
// BOUNDARY
//
// "Exactly" is a closure claim: no sixth kind, and each constant carries the
// stable wire spelling RDR 0005 maps.
func TestReq7_RefusalKindSetIsExactlyTheFiveNamedKinds(t *testing.T) {
	want := []string{
		"ambiguous_match",
		"guard_unevaluable",
		"no_match",
		"owned_state_unavailable",
		"unmodeled_outcome",
	}

	got := make([]string, 0, len(resolve.RefusalKinds()))
	for _, k := range resolve.RefusalKinds() {
		got = append(got, string(k))
	}
	sort.Strings(got)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("refusal kind set = %v; want exactly %v", got, want)
	}

	named := map[string]resolve.RefusalKind{
		"no_match":                resolve.KindNoMatch,
		"ambiguous_match":         resolve.KindAmbiguousMatch,
		"owned_state_unavailable": resolve.KindOwnedStateUnavailable,
		"guard_unevaluable":       resolve.KindGuardUnevaluable,
		"unmodeled_outcome":       resolve.KindUnmodeledOutcome,
	}
	for spelling, k := range named {
		if string(k) != spelling {
			t.Errorf("constant for %q spells %q", spelling, string(k))
		}
	}
}

// REQ-8: "Each refusal kind must be stable enough for RDR 0005 to map to a
// CLI error code without inspecting error strings."
// DOMAIN EDGE
//
// Stands in for the RDR 0005 mapping layer: a total switch over the kind
// value alone, reading no error and no message text, must classify every
// refusal the kernel can produce.
func TestReq8_RefusalKindAloneDrivesCLIMappingWithoutErrorStrings(t *testing.T) {
	mapCode := func(k resolve.RefusalKind) string {
		switch k {
		case resolve.KindNoMatch:
			return "resolver-no-match"
		case resolve.KindAmbiguousMatch:
			return "resolver-ambiguous-match"
		case resolve.KindOwnedStateUnavailable:
			return "resolver-owned-state-unavailable"
		case resolve.KindGuardUnevaluable:
			return "resolver-guard-unevaluable"
		case resolve.KindUnmodeledOutcome:
			return "resolver-unmodeled-outcome"
		default:
			return ""
		}
	}

	seen := make(map[string]bool)
	for name, in := range allRefusalInputs() {
		r := refusalOf(t, mustResolve(t, in))
		code := mapCode(r.Kind)
		if code == "" {
			t.Errorf("%s: refusal kind %q is not mappable by kind alone", name, r.Kind)
			continue
		}
		if seen[code] {
			t.Errorf("%s: kind %q collides with an already-mapped CLI code %q", name, r.Kind, code)
		}
		seen[code] = true
	}
	if len(seen) != 5 {
		t.Errorf("mapped %d distinct CLI codes; want 5 (one per kernel refusal kind)", len(seen))
	}
}

// REQ-9: "Visible failures are typed refusal values: no_match,
// ambiguous_match, owned_state_unavailable, guard_unevaluable, or
// unmodeled_outcome."
// ADVERSARIAL
//
// Every refusal the kernel produces must carry a kind drawn from the closed
// set — never the zero value of an unset field, never an off-list kind.
func TestReq9_EveryRefusalCarriesAKindFromTheClosedSet(t *testing.T) {
	inSet := func(k resolve.RefusalKind) bool {
		for _, known := range resolve.RefusalKinds() {
			if k == known {
				return true
			}
		}
		return false
	}
	for name, in := range allRefusalInputs() {
		t.Run(name, func(t *testing.T) {
			r := refusalOf(t, mustResolve(t, in))
			if !inSet(r.Kind) {
				t.Errorf("refusal kind %q is not a member of the closed kernel set", r.Kind)
			}
		})
	}
}

// REQ-10: "Diagnosis starts with the input tuple, the refusal kind, and the
// transition table revision used for that resolution."
// DOMAIN EDGE
func TestReq10_RefusalCarriesInputTupleIdentityAndTableRevision(t *testing.T) {
	for name, in := range allRefusalInputs() {
		t.Run(name, func(t *testing.T) {
			r := refusalOf(t, mustResolve(t, in))
			if r.Revision != in.Table.Revision {
				t.Errorf("refusal revision = %q; want %q", r.Revision, in.Table.Revision)
			}
			if r.Flow != in.Flow {
				t.Errorf("refusal flow = %q; want %q", r.Flow, in.Flow)
			}
			if r.Recognized != in.Recognized {
				t.Errorf("refusal recognized outcome = %q; want %q", r.Recognized, in.Recognized)
			}
			if r.Kind == "" {
				t.Error("refusal carries no kind")
			}
		})
	}
}

// ---------------------------------------------------------------------
// Prohibitions (negative contract)
// ---------------------------------------------------------------------

// REQ-11: "The kernel MUST NOT print output, inspect CLI flags, discover
// ambient state, choose artifacts on behalf of the caller, initiate work, or
// execute persistence side effects directly."
// ADVERSARIAL
//
// Asserted at the dependency boundary: a kernel that printed output, read
// flags, or executed persistence would have to reach the CLI, output, or
// filesystem packages. The import graph is the checkable trace of that.
func TestReq11_KernelImportsNoCLIOutputOrPersistenceFacility(t *testing.T) {
	forbidden := []string{
		"github.com/spf13/cobra",
		"github.com/spf13/pflag",
		"os",
		"os/exec",
		"io/ioutil",
		"path/filepath",
		"log",
	}
	imports := kernelPackageImports(t)
	for _, f := range forbidden {
		if imports[f] {
			t.Errorf("resolver kernel imports %q; the kernel must not print output, inspect CLI flags, discover ambient state, or execute persistence", f)
		}
	}
	// The CLI tree is checked by prefix rather than by name: an enumeration
	// goes stale as CLI packages are added or deleted, and REQ-11 prohibits a
	// class of behaviour, not a fixed list of packages.
	for imp := range imports {
		if isCLIPackage(imp) {
			t.Errorf("resolver kernel imports CLI package %q; the kernel must not print output or inspect CLI flags", imp)
		}
	}
}

// REQ-12: "Silent failure would mean the kernel guessed a transition or
// executed persistence directly; both are prohibited by the normative
// contract."
// ADVERSARIAL
//
// The observable form of "did not guess": a refusing disposition carries no
// next state and no write description, so no caller can mistake it for a
// transition.
func TestReq12_RefusalNeverCarriesAGuessedTransition(t *testing.T) {
	for name, in := range allRefusalInputs() {
		t.Run(name, func(t *testing.T) {
			got := mustResolve(t, in)
			refusalOf(t, got)
			if got.Plan != nil {
				t.Fatalf("refusal disposition also carries a plan: %+v", *got.Plan)
			}
		})
	}
}

// REQ-13: "The kernel does not discover artifacts, execute accessors,
// discover work, run subflows, print output, or own persistence."
// ADVERSARIAL
//
// The accessor seam is not a kernel collaborator at all: the kernel's input
// carries an already-read owned snapshot, never an accessor to call. If the
// kernel could execute an accessor, an owned tag absent from the snapshot
// would be fetchable rather than a refusal.
func TestReq13_KernelDoesNotExecuteAccessorsToFillMissingOwnedState(t *testing.T) {
	in := missingOwnedInput()
	r := refusalOf(t, mustResolve(t, in))
	if r.Kind != resolve.KindOwnedStateUnavailable {
		t.Errorf("refusal kind = %q; want %q — a missing owned tag must refuse, not trigger a read",
			r.Kind, resolve.KindOwnedStateUnavailable)
	}
}

// REQ-14: "The resolver belongs in an internal package behind the CLI: it
// must stay stateless and non-orchestrating, consume owned-state snapshots
// and write targets produced by the accessor layer, and refuse illegal or
// incomplete transition inputs rather than initiating work."
// ADVERSARIAL
//
// Statelessness is observable as call-order independence: interleaving other
// resolutions between two replays of the same tuple must not change the
// disposition, because the kernel retains nothing across calls.
func TestReq14_KernelIsStatelessAcrossInterleavedCalls(t *testing.T) {
	first := mustResolve(t, legalInput())
	planOf(t, first)

	for name, other := range allRefusalInputs() {
		got := mustResolve(t, other)
		if got.Refusal == nil {
			t.Fatalf("interleaved %s input did not refuse; the interleaving is not exercising the kernel", name)
		}
	}
	mustResolve(t, legalInput())

	last := mustResolve(t, legalInput())
	planOf(t, last)
	if !reflect.DeepEqual(first, last) {
		t.Errorf("disposition changed after interleaved calls; the kernel is carrying state:\nfirst = %+v\nlast  = %+v", first, last)
	}
}

// ---------------------------------------------------------------------
// Selection semantics
// ---------------------------------------------------------------------

// REQ-15: "the only successful selection is exactly one matching edge after
// guard evaluation. Zero, multiple, unavailable, or unevaluable candidates
// are refusals unless the table contains a modeled escape edge that itself
// matches exactly once."
// BOUNDARY
//
// The positive half: exactly one match after guards succeeds.
func TestReq15_ExactlyOneMatchAfterGuardsIsTheOnlySuccess(t *testing.T) {
	in := legalInput()
	in.Table = twoRowsOneGuardFalseTable()
	in.Guards = fixtureGuards{decided: map[string]bool{
		"is-fast-lane": false,
		"is-slow-lane": true,
	}}

	p := planOf(t, mustResolve(t, in))
	if p.RuleID != "rdr.draft.successful.slow" {
		t.Errorf("selected rule = %q; want the single row surviving guard evaluation", p.RuleID)
	}
}

// REQ-15: "Zero, multiple, unavailable, or unevaluable candidates are
// refusals unless the table contains a modeled escape edge that itself
// matches exactly once."
// DOMAIN EDGE
//
// The negative half, read as RDR 0001 states it: each of the four listed
// conditions refuses when no escape edge is modeled. (RDR 0002 narrows
// escapes to no_match/ambiguous_match; RDR 0001 is the governing spec here
// and its clause lists all four, so all four are asserted to refuse absent
// an escape.)
func TestReq15_AllFourListedConditionsRefuseWithoutAnEscapeEdge(t *testing.T) {
	cases := map[string]struct {
		in   resolve.Input
		kind resolve.RefusalKind
	}{
		"zero candidates":       {noMatchInput(), resolve.KindNoMatch},
		"multiple candidates":   {ambiguousInput(), resolve.KindAmbiguousMatch},
		"unavailable candidate": {missingOwnedInput(), resolve.KindOwnedStateUnavailable},
		"unevaluable candidate": {unevaluableGuardInput(), resolve.KindGuardUnevaluable},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			for _, row := range tc.in.Table.Rows {
				if len(row.Escape) != 0 {
					t.Fatalf("fixture models an escape edge on %q; this case must have none", row.RuleID)
				}
			}
			r := refusalOf(t, mustResolve(t, tc.in))
			if r.Kind != tc.kind {
				t.Errorf("refusal kind = %q; want %q", r.Kind, tc.kind)
			}
		})
	}
}

// REQ-15: "unless the table contains a modeled escape edge that itself
// matches exactly once."
// DOMAIN EDGE
//
// An escape edge rescues its modeled failure class, and the rescue is
// itself exact-one: a table with two matching escape rows for the same
// class must not silently pick one.
func TestReq15_EscapeEdgeRescuesOnlyWhenItMatchesExactlyOnce(t *testing.T) {
	t.Run("single escape edge rescues no_match", func(t *testing.T) {
		in := noMatchInput()
		in.Table.Rows = append(in.Table.Rows, escapeRow("rdr.escape.nomatch", "flows/rdr.toml:90", resolve.KindNoMatch))

		p := planOf(t, mustResolve(t, in))
		if p.RuleID != "rdr.escape.nomatch" {
			t.Errorf("selected rule = %q; want the modeled escape edge", p.RuleID)
		}
		if !p.Escaped {
			t.Error("plan does not report that it came from a modeled escape edge")
		}
	})

	t.Run("two matching escape edges do not rescue", func(t *testing.T) {
		in := noMatchInput()
		in.Table.Rows = append(in.Table.Rows,
			escapeRow("rdr.escape.a", "flows/rdr.toml:90", resolve.KindNoMatch),
			escapeRow("rdr.escape.b", "flows/rdr.toml:95", resolve.KindNoMatch),
		)
		r := refusalOf(t, mustResolve(t, in))
		if r.Kind != resolve.KindAmbiguousMatch {
			t.Errorf("refusal kind = %q; want %q — an escape that matches twice is not an exact-one rescue",
				r.Kind, resolve.KindAmbiguousMatch)
		}
	})

	t.Run("escape modeled for another class does not rescue", func(t *testing.T) {
		in := noMatchInput()
		in.Table.Rows = append(in.Table.Rows,
			escapeRow("rdr.escape.wrongclass", "flows/rdr.toml:90", resolve.KindAmbiguousMatch))

		r := refusalOf(t, mustResolve(t, in))
		if r.Kind != resolve.KindNoMatch {
			t.Errorf("refusal kind = %q; want %q — an escape modeled for a different class must not rescue",
				r.Kind, resolve.KindNoMatch)
		}
	})
}

// REQ-16: "Evaluation builds a single tag-set view, selects matching
// candidate edges, refuses zero or multiple matches unless the table contract
// explicitly models an escape edge, and emits a transition plan."
// BOUNDARY
func TestReq16_ZeroAndMultipleMatchesRefuseUnlessEscapeIsModeled(t *testing.T) {
	t.Run("zero without escape refuses", func(t *testing.T) {
		r := refusalOf(t, mustResolve(t, noMatchInput()))
		if r.Kind != resolve.KindNoMatch {
			t.Errorf("refusal kind = %q; want %q", r.Kind, resolve.KindNoMatch)
		}
	})

	t.Run("multiple without escape refuses", func(t *testing.T) {
		r := refusalOf(t, mustResolve(t, ambiguousInput()))
		if r.Kind != resolve.KindAmbiguousMatch {
			t.Errorf("refusal kind = %q; want %q", r.Kind, resolve.KindAmbiguousMatch)
		}
	})

	t.Run("multiple with a modeled ambiguous escape emits a plan", func(t *testing.T) {
		in := ambiguousInput()
		in.Table.Rows = append(in.Table.Rows,
			escapeRow("rdr.escape.ambiguous", "flows/rdr.toml:99", resolve.KindAmbiguousMatch))

		p := planOf(t, mustResolve(t, in))
		if p.RuleID != "rdr.escape.ambiguous" {
			t.Errorf("selected rule = %q; want the modeled escape edge", p.RuleID)
		}
	})
}

// REQ-17: "Merge owned, observed, and freshly recognized tags into the
// evaluation view."
// DOMAIN EDGE
//
// All three provenances must reach selection: a row that can only match by
// reading an observed tag, and a row that can only match by reading the
// freshly recognized tag, must both be selectable.
func TestReq17_OwnedObservedAndRecognizedTagsAllReachSelection(t *testing.T) {
	t.Run("observed tag participates in matching", func(t *testing.T) {
		in := legalInput()
		in.Table = observedSensitiveTable()

		p := planOf(t, mustResolve(t, in))
		if p.RuleID != "rdr.draft.reviewed" {
			t.Errorf("selected rule = %q; want the row matching on the observed tag", p.RuleID)
		}

		// Removing only the observed tag must change the disposition.
		in2 := legalInput()
		in2.Table = observedSensitiveTable()
		in2.Observed = nil
		if got := mustResolve(t, in2); got.Refusal == nil {
			t.Error("dropping the observed tag still produced a plan; observed tags are not in the evaluation view")
		}
	})

	t.Run("freshly recognized tag participates in matching", func(t *testing.T) {
		in := legalInput()
		in.Table = recognizedTagSensitiveTable()

		p := planOf(t, mustResolve(t, in))
		if p.RuleID != "rdr.draft.on-recognized-tag" {
			t.Errorf("selected rule = %q; want the row matching on the freshly recognized tag", p.RuleID)
		}
	})

	t.Run("owned tag participates in matching", func(t *testing.T) {
		in := legalInput()
		in.Owned = []resolve.Tag{{Key: "status", Value: "Review"}}
		if got := mustResolve(t, in); got.Refusal == nil {
			t.Error("changing the owned tag still matched the Draft row; owned tags are not in the evaluation view")
		}
	})
}

// ---------------------------------------------------------------------
// Success plan shape
// ---------------------------------------------------------------------

// REQ-18: "A successful resolution returns the next state tags and the
// owned-tag writes that the accessor layer applies back to the same
// caller-provided artifact boundary; an illegal, ambiguous, incomplete, or
// unmodeled input returns an explicit refusal class."
// HAPPY PATH
func TestReq18_SuccessCarriesNextTagsAndOwnedTagWrites(t *testing.T) {
	p := planOf(t, mustResolve(t, legalInput()))

	wantNext := []resolve.Tag{{Key: "status", Value: "Final"}}
	if !reflect.DeepEqual(p.NextTags, wantNext) {
		t.Errorf("next tags = %+v; want %+v", p.NextTags, wantNext)
	}
	if len(p.Writes) == 0 {
		t.Fatal("plan carries no owned-tag writes for the accessor layer to apply")
	}
	wantWrites := []resolve.Tag{{Key: "status", Value: "Final"}}
	if !reflect.DeepEqual(p.Writes, wantWrites) {
		t.Errorf("writes = %+v; want %+v", p.Writes, wantWrites)
	}
}

// REQ-18: "an illegal, ambiguous, incomplete, or unmodeled input returns an
// explicit refusal class."
// ADVERSARIAL
//
// The refusal counterpart of the clause: a refusal carries no write
// description at all, so nothing downstream can apply a partial transition.
func TestReq18_RefusalCarriesNoWriteDescription(t *testing.T) {
	for name, in := range allRefusalInputs() {
		t.Run(name, func(t *testing.T) {
			got := mustResolve(t, in)
			refusalOf(t, got)
			if got.Plan != nil {
				t.Fatalf("refusal carries a plan with writes %+v", got.Plan.Writes)
			}
		})
	}
}

// REQ-19: "The kernel may describe owned-tag writes in the returned plan, and
// the accessor layer knows how to apply those writes to caller-provided
// artifacts, but RDR 0004 owns how reads, writes, and read-back verification
// execute."
// DOMAIN EDGE
//
// The plan describes writes; it never executes them. Observably: the writes
// are a description in the returned value, and the caller's owned snapshot
// is untouched by the call.
func TestReq19_WritesAreDescribedNotApplied(t *testing.T) {
	in := legalInput()
	ownedBefore := deepCopyTags(in.Owned)

	p := planOf(t, mustResolve(t, in))
	if len(p.Writes) == 0 {
		t.Fatal("plan describes no owned-tag writes")
	}
	if !reflect.DeepEqual(in.Owned, ownedBefore) {
		t.Errorf("Resolve applied writes to the owned snapshot: got %+v, want %+v", in.Owned, ownedBefore)
	}
}

// REQ-19: "The kernel may describe owned-tag writes"
// DOMAIN EDGE
//
// Only owned-provenance tags may appear in a write description — RDR 0004
// forbids the accessor layer from writing observed or recognized tags.
func TestReq19_OnlyOwnedProvenanceTagsAppearInWrites(t *testing.T) {
	in := legalInput()
	in.Table = observedSensitiveTable()

	p := planOf(t, mustResolve(t, in))
	observedKeys := map[string]bool{}
	for _, tg := range in.Observed {
		observedKeys[tg.Key] = true
	}
	for _, w := range p.Writes {
		if observedKeys[w.Key] {
			t.Errorf("write description includes observed tag %q; only owned tags are writable", w.Key)
		}
	}
}

// REQ-20: "Inputs are: flow identity, transition table revision, owned-state
// snapshot values produced by the accessor layer, observed tags supplied by
// the caller, the freshly recognized outcome tag, and the reviewable
// transition table."
// INPUT EDGE
//
// The input surface is exactly these; nothing is discovered. A fully empty
// input tuple must still produce a disposition through the value path — the
// kernel refuses rather than reaching for ambient defaults.
func TestReq20_EmptyInputTupleStillYieldsAValueDisposition(t *testing.T) {
	got, err := resolve.Resolve(resolve.Input{})
	if err != nil {
		t.Fatalf("empty input tuple traveled the Go error path: %v", err)
	}
	if got.Refusal == nil {
		t.Fatal("empty input tuple did not produce a refusal disposition")
	}
	if got.Plan != nil {
		t.Errorf("empty input tuple produced a plan: %+v", *got.Plan)
	}
}

// REQ-21: "the kernel only exposes structured success/refusal values that the
// CLI can map."
// DOMAIN EDGE
//
// Every field a caller needs to map is a structured value, not prose: the
// refusal discriminator is a comparable kind, not a message. Asserted by
// mapping every kernel refusal through a switch on kind alone and requiring
// the kind values be distinct and non-empty.
func TestReq21_DispositionsAreStructuredValuesNotProse(t *testing.T) {
	kinds := map[resolve.RefusalKind]bool{}
	for name, in := range allRefusalInputs() {
		r := refusalOf(t, mustResolve(t, in))
		if r.Kind == "" {
			t.Errorf("%s: refusal has an empty kind", name)
		}
		if kinds[r.Kind] {
			t.Errorf("%s: kind %q already produced by another refusal fixture", name, r.Kind)
		}
		kinds[r.Kind] = true
	}

	p := planOf(t, mustResolve(t, legalInput()))
	if p.RuleID == "" || p.SourceLocator == "" {
		t.Error("plan lacks structured source identity for the CLI to report")
	}
}

// ---------------------------------------------------------------------
// Package boundary and scope
// ---------------------------------------------------------------------

// REQ-22: "Define the internal resolver package boundary, result taxonomy,
// and pure resolution entry point without CLI output or persistence side
// effects."
// BOUNDARY
func TestReq22_PackageIsInternalAndExposesAPureEntryPoint(t *testing.T) {
	pkgPath := kernelImportPath(t)
	if !isInternalPackage(pkgPath) {
		t.Errorf("resolver package %q is not under internal/", pkgPath)
	}

	// The entry point is pure in the observable sense: two calls with the
	// same tuple yield the same value and touch nothing outside the return.
	a := mustResolve(t, legalInput())
	b := mustResolve(t, legalInput())
	// Anchor: purity over a no-op entry point is vacuous.
	planOf(t, a)
	planOf(t, b)
	if !reflect.DeepEqual(a, b) {
		t.Error("entry point is not pure: identical tuples produced different values")
	}
}

// REQ-23: "Implement tag-set assembly, guard evaluation delegation, exact-one
// edge selection, and typed refusal behavior."
// DOMAIN EDGE
//
// Guard evaluation is delegated, not implemented in the kernel: the kernel
// must ask the injected seam and honour its verdict rather than deciding the
// predicate itself.
func TestReq23_GuardEvaluationIsDelegatedToTheInjectedSeam(t *testing.T) {
	in := legalInput()
	in.Table = twoRowsOneGuardFalseTable()

	t.Run("seam says fast lane holds", func(t *testing.T) {
		in := in
		in.Guards = fixtureGuards{decided: map[string]bool{"is-fast-lane": true, "is-slow-lane": false}}
		p := planOf(t, mustResolve(t, in))
		if p.RuleID != "rdr.draft.successful.fast" {
			t.Errorf("selected rule = %q; want the row the seam decided true", p.RuleID)
		}
	})

	t.Run("seam says slow lane holds", func(t *testing.T) {
		in := in
		in.Guards = fixtureGuards{decided: map[string]bool{"is-fast-lane": false, "is-slow-lane": true}}
		p := planOf(t, mustResolve(t, in))
		if p.RuleID != "rdr.draft.successful.slow" {
			t.Errorf("selected rule = %q; want the row the seam decided true", p.RuleID)
		}
	})

	t.Run("seam decides both true so selection is ambiguous", func(t *testing.T) {
		in := in
		in.Guards = fixtureGuards{decided: map[string]bool{"is-fast-lane": true, "is-slow-lane": true}}
		r := refusalOf(t, mustResolve(t, in))
		if r.Kind != resolve.KindAmbiguousMatch {
			t.Errorf("refusal kind = %q; want %q", r.Kind, resolve.KindAmbiguousMatch)
		}
	})

	t.Run("seam decides both false so nothing matches", func(t *testing.T) {
		in := in
		in.Guards = fixtureGuards{decided: map[string]bool{"is-fast-lane": false, "is-slow-lane": false}}
		r := refusalOf(t, mustResolve(t, in))
		if r.Kind != resolve.KindNoMatch {
			t.Errorf("refusal kind = %q; want %q", r.Kind, resolve.KindNoMatch)
		}
	})
}

// REQ-24: "Expose only the kernel values needed by RDR 0005; do not add
// command output or state mutation semantics in this RDR."
// ADVERSARIAL
//
// No state mutation semantics: the kernel mutates neither the owned snapshot,
// the observed tags, nor the table it is handed.
func TestReq24_KernelAddsNoStateMutationSemantics(t *testing.T) {
	for name, in := range allDispositionInputs() {
		t.Run(name, func(t *testing.T) {
			before := deepCopyInput(in)
			got, err := resolve.Resolve(in)
			if err != nil {
				t.Fatalf("Resolve used the Go error path: %v", err)
			}
			if got.Plan == nil && got.Refusal == nil {
				t.Fatal("Resolve produced no disposition; immutability over a no-op call proves nothing")
			}
			if !reflect.DeepEqual(in, before) {
				t.Errorf("Resolve mutated its inputs:\n got  = %+v\n want = %+v", in, before)
			}
		})
	}
}

// REQ-25: "the internal component is the resolver kernel. Rejected names:
// 'orchestrator' because it implies initiating work, and 'state machine
// runner' because it implies owning persistence."
// DOMAIN EDGE
//
// The rejected names are rejected because of the behaviours they imply. The
// checkable form of "does not initiate work / does not own persistence" is
// that no exported kernel symbol offers to run, execute, apply, or persist.
func TestReq25_NoExportedSymbolImpliesOrchestrationOrPersistence(t *testing.T) {
	banned := []string{"Run", "Execute", "Apply", "Persist", "Save", "Commit", "Start", "Orchestrate"}
	names := exportedKernelSymbols(t)
	for _, n := range names {
		for _, b := range banned {
			if n == b {
				t.Errorf("exported symbol %q implies initiating work or owning persistence", n)
			}
		}
	}
}

// REQ-26: "No third-party dependency is proposed at this stage."
// BOUNDARY
func TestReq26_KernelUsesNoThirdPartyDependency(t *testing.T) {
	for imp := range kernelPackageImports(t) {
		if isThirdParty(imp) {
			t.Errorf("resolver kernel imports third-party package %q", imp)
		}
	}
}

// REQ-27: "the resolver is stateless and receives snapshots and table data as
// inputs; it does not own shared persistence or execute accessor writes."
// ADVERSARIAL
//
// Statelessness under concurrency: the same tuple resolved from many
// goroutines must yield the same value, which cannot hold if the kernel owns
// shared mutable state.
func TestReq27_ConcurrentResolutionsAreIndependent(t *testing.T) {
	want := mustResolve(t, legalInput())
	// Anchor: agreeing on an empty disposition would prove nothing.
	planOf(t, want)

	const workers = 16
	results := make(chan resolve.Result, workers)
	errs := make(chan error, workers)
	for i := 0; i < workers; i++ {
		go func() {
			got, err := resolve.Resolve(legalInput())
			results <- got
			errs <- err
		}()
	}
	for i := 0; i < workers; i++ {
		if err := <-errs; err != nil {
			t.Errorf("concurrent Resolve returned a Go error: %v", err)
		}
		if got := <-results; !reflect.DeepEqual(got, want) {
			t.Errorf("concurrent disposition differs:\n got  = %+v\n want = %+v", got, want)
		}
	}
}

// ---------------------------------------------------------------------
// Testing strategy scenarios
// ---------------------------------------------------------------------

// REQ-28: "Implementation must add focused resolver-kernel tests before any
// CLI command wraps the kernel."
// BOUNDARY
//
// The dependency direction is CLI -> kernel only. A kernel that the CLI
// already wrapped in the reverse direction would show as a kernel import of
// the CLI tree.
func TestReq28_KernelIsTestableWithoutAnyCLICommand(t *testing.T) {
	imports := kernelPackageImports(t)
	for imp := range imports {
		if isCLIPackage(imp) {
			t.Errorf("kernel imports CLI package %q; the kernel must be exercisable before any command wraps it", imp)
		}
	}
	// And the kernel resolves with no CLI participation at all.
	planOf(t, mustResolve(t, legalInput()))
}

// REQ-29: "Scenario: Replay the same legal input tuple twice. Expected: Both
// calls return value-identical transition plans, including next tags and
// accessor write descriptions."
// HAPPY PATH
func TestReq29_ReplayReturnsValueIdenticalPlansIncludingWrites(t *testing.T) {
	first := planOf(t, mustResolve(t, legalInput()))
	second := planOf(t, mustResolve(t, legalInput()))

	if !reflect.DeepEqual(first.NextTags, second.NextTags) {
		t.Errorf("next tags differ across replay: %+v vs %+v", first.NextTags, second.NextTags)
	}
	if !reflect.DeepEqual(first.Writes, second.Writes) {
		t.Errorf("write descriptions differ across replay: %+v vs %+v", first.Writes, second.Writes)
	}
	if !reflect.DeepEqual(*first, *second) {
		t.Errorf("plans are not value-identical:\n%+v\n%+v", *first, *second)
	}
	if len(first.Writes) == 0 {
		t.Error("replayed plan carries no accessor write descriptions to compare")
	}
}

// REQ-30: "Scenario: No table edge matches the assembled tag-set. Expected:
// The kernel returns the typed no-match refusal and performs no persistence
// side effect."
// DOMAIN EDGE
func TestReq30_NoMatchRefusesWithoutPersistenceSideEffect(t *testing.T) {
	in := noMatchInput()
	before := deepCopyInput(in)

	got := mustResolve(t, in)
	r := refusalOf(t, got)
	if r.Kind != resolve.KindNoMatch {
		t.Errorf("refusal kind = %q; want %q", r.Kind, resolve.KindNoMatch)
	}
	if got.Plan != nil {
		t.Errorf("no-match refusal carries a plan with writes %+v", got.Plan.Writes)
	}
	if !reflect.DeepEqual(in, before) {
		t.Error("no-match resolution mutated the caller's input")
	}
}

// REQ-31: "Scenario: More than one table edge matches after guard evaluation.
// Expected: The kernel returns the typed ambiguous-match refusal unless one
// explicit escape edge matches exactly once."
// DOMAIN EDGE
func TestReq31_AmbiguousMatchRefusesAndNamesTheConflictingRows(t *testing.T) {
	r := refusalOf(t, mustResolve(t, ambiguousInput()))
	if r.Kind != resolve.KindAmbiguousMatch {
		t.Fatalf("refusal kind = %q; want %q", r.Kind, resolve.KindAmbiguousMatch)
	}

	// The refusal must carry enough source identity for RDR 0005 to report
	// which rows conflicted.
	if len(r.Rows) < 2 {
		t.Fatalf("ambiguous refusal names %d rows; want the 2 conflicting rows", len(r.Rows))
	}
	got := map[string]string{}
	for _, ref := range r.Rows {
		got[ref.RuleID] = ref.SourceLocator
	}
	for _, want := range []struct{ id, loc string }{
		{"rdr.draft.successful.a", "flows/rdr.toml:12"},
		{"rdr.draft.successful.b", "flows/rdr.toml:30"},
	} {
		if got[want.id] != want.loc {
			t.Errorf("row %q locator = %q; want %q", want.id, got[want.id], want.loc)
		}
	}
}

// REQ-32: "Scenario: The freshly recognized outcome is not modeled by the
// table. Expected: The kernel returns the typed unmodeled-outcome refusal and
// performs no persistence side effect."
// DOMAIN EDGE
//
// unmodeled_outcome is decided against the table's declared alphabet and
// takes precedence over a mere zero-match.
func TestReq32_UnmodeledOutcomeRefusesAndIsNotAMereNoMatch(t *testing.T) {
	in := unmodeledOutcomeInput()
	got := mustResolve(t, in)
	r := refusalOf(t, got)

	if r.Kind != resolve.KindUnmodeledOutcome {
		t.Errorf("refusal kind = %q; want %q — the outcome is outside the table's declared alphabet %v",
			r.Kind, resolve.KindUnmodeledOutcome, in.Table.Outcomes)
	}
	if got.Plan != nil {
		t.Errorf("unmodeled-outcome refusal carries a plan with writes %+v", got.Plan.Writes)
	}
}

// REQ-33: "Scenario: Owned state is unavailable or a guard cannot be
// evaluated. Expected: The kernel returns the corresponding value-level typed
// refusal and does not fall back to ambient discovery or the CLI/Go error
// path."
// DOMAIN EDGE
func TestReq33_UnavailableOwnedStateAndUnevaluableGuardAreValueRefusals(t *testing.T) {
	t.Run("owned state unavailable", func(t *testing.T) {
		got, err := resolve.Resolve(missingOwnedInput())
		if err != nil {
			t.Fatalf("owned-state refusal traveled the Go error path: %v", err)
		}
		r := refusalOf(t, got)
		if r.Kind != resolve.KindOwnedStateUnavailable {
			t.Errorf("refusal kind = %q; want %q", r.Kind, resolve.KindOwnedStateUnavailable)
		}
		if len(r.MissingOwned) == 0 {
			t.Error("refusal names no missing owned tag key for diagnosis")
		}
	})

	t.Run("guard unevaluable", func(t *testing.T) {
		got, err := resolve.Resolve(unevaluableGuardInput())
		if err != nil {
			t.Fatalf("guard refusal traveled the Go error path: %v", err)
		}
		r := refusalOf(t, got)
		if r.Kind != resolve.KindGuardUnevaluable {
			t.Errorf("refusal kind = %q; want %q", r.Kind, resolve.KindGuardUnevaluable)
		}
		if len(r.Undecided) == 0 {
			t.Error("refusal names no undecidable guard for diagnosis")
		}
	})

	t.Run("nil guard seam does not fall back to ambient evaluation", func(t *testing.T) {
		in := unevaluableGuardInput()
		in.Guards = nil
		got, err := resolve.Resolve(in)
		if err != nil {
			t.Fatalf("absent guard seam traveled the Go error path: %v", err)
		}
		r := refusalOf(t, got)
		if r.Kind != resolve.KindGuardUnevaluable {
			t.Errorf("refusal kind = %q; want %q — a missing seam must refuse, not evaluate the predicate itself",
				r.Kind, resolve.KindGuardUnevaluable)
		}
	})
}

// REQ-34: "The tests exercise the pure decision boundary: fixture transition
// tables, owned snapshots from accessor fixtures, observed tags supplied by
// the caller, and freshly recognized outcome tags."
// BOUNDARY
//
// The decision boundary is pure in the sense that the disposition is a
// function of the supplied tuple alone: two callers supplying equal tuples
// built independently get equal dispositions.
func TestReq34_DispositionIsAFunctionOfTheSuppliedTupleAlone(t *testing.T) {
	a := legalInput()

	b := resolve.Input{
		Flow:       "rdr",
		Table:      singleMatchTable(),
		Owned:      []resolve.Tag{{Key: "status", Value: "Draft"}},
		Observed:   []resolve.Tag{{Key: "reviews", Value: "2"}},
		Recognized: "successful",
		Guards:     allGuardsTrue(),
	}

	gotA := mustResolve(t, a)
	gotB := mustResolve(t, b)

	// Anchor: both tuples must select the real edge.
	planOf(t, gotA)
	planOf(t, gotB)

	if !reflect.DeepEqual(gotA, gotB) {
		t.Error("independently built but equal input tuples produced different dispositions")
	}
}

// REQ-35: "Add focused kernel tests for deterministic replay and refusal
// classes, using fixtures that exercise owned, observed, and freshly
// recognized tags."
// HAPPY PATH
//
// Replay determinism must hold for refusals too, not only for plans: a
// refusal is a disposition and is subject to the same replay claim.
func TestReq35_RefusalDispositionsReplayValueIdentically(t *testing.T) {
	for name, build := range refusalInputBuilders() {
		t.Run(name, func(t *testing.T) {
			first := mustResolve(t, build())
			second := mustResolve(t, build())

			// Anchor: each replay must be a real refusal of the kind this
			// fixture models. Two empty dispositions are not a replay.
			r1 := refusalOf(t, first)
			r2 := refusalOf(t, second)
			if string(r1.Kind) != name {
				t.Fatalf("refusal kind = %q; want %q", r1.Kind, name)
			}
			if string(r2.Kind) != name {
				t.Fatalf("replayed refusal kind = %q; want %q", r2.Kind, name)
			}

			if !reflect.DeepEqual(first, second) {
				t.Errorf("refusal disposition differs across replay:\nfirst  = %+v\nsecond = %+v", *r1, *r2)
			}
		})
	}
}

// ---------------------------------------------------------------------
// Non-goals asserted as contract
// ---------------------------------------------------------------------

// REQ-36: "No encode/decode, import/export, or inverse operation is
// introduced by this RDR."
// BOUNDARY
func TestReq36_KernelExposesNoEncodeDecodeOrInverseOperation(t *testing.T) {
	banned := []string{
		"Encode", "Decode", "Marshal", "Unmarshal",
		"Import", "Export", "Parse", "Serialize", "Deserialize", "Invert",
	}
	for _, n := range exportedKernelSymbols(t) {
		for _, b := range banned {
			if n == b {
				t.Errorf("exported symbol %q introduces an encode/decode or inverse operation this RDR excludes", n)
			}
		}
	}
}

// REQ-37: "No throughput target is part of this RDR. … No byte-stable hash or
// canonical serialization is introduced"
// BOUNDARY
func TestReq37_KernelIntroducesNoHashOrCanonicalSerialization(t *testing.T) {
	for imp := range kernelPackageImports(t) {
		switch imp {
		case "crypto/sha256", "crypto/sha1", "crypto/md5", "hash", "hash/fnv",
			"hash/crc32", "encoding/json", "encoding/gob", "encoding/hex":
			t.Errorf("resolver kernel imports %q; no byte-stable hash or canonical serialization is introduced by this RDR", imp)
		}
	}

	banned := []string{"Hash", "Checksum", "Canonicalize", "Fingerprint", "Digest"}
	for _, n := range exportedKernelSymbols(t) {
		for _, b := range banned {
			if n == b {
				t.Errorf("exported symbol %q introduces a hash or canonical form this RDR excludes", n)
			}
		}
	}
}
