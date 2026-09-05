package cli

// RDR 0020 — `0020:MVV`, run end to end as one test.
//
// The five MVV steps are ONE runnable sequence here rather than five
// independent tests, because the record's Testing Strategy fixes the
// construction: "done is all six green, with the undeclared/declared pair
// adjacent in one test so the asymmetry C1 fixes is pinned by construction
// rather than by two tests that could drift". Splitting the pair would
// satisfy the letter and lose the property.
//
// The MVV's fidelity claim is a ROUND TRIP through the CLI: a value the
// caller spells on argv must come back out of the payload's `observed` echo
// as the SAME BYTES. A green exit code proves nothing here — an
// implementation that admitted the array and canonicalised it would exit 0
// and lose the caller's spelling silently. So every carrier leg compares
// the echoed value to the argv token, byte for byte.

import (
	"strings"
	"testing"
)

// REQ-MVV / `0020:MVV` — the five steps, end to end.
//
// REQ-42 / MVV 1: "Author a fixture model declaring one scalar tag and one
// set tag, with no declaration for `extra` or `extras`."
// REQ-43 / MVV 2: "Run `flow resolve` with `--tag extra=plain` and `--tag
// extras=[\"a\",\"b\"]` beside the declared tags: the invocation exits 0
// and the payload's `observed` field carries both values byte-for-byte as
// given."
// REQ-44 / MVV 3: "Run the same invocation without the two carrier flags:
// the selected rule and outcome are identical — the carried keys influenced
// nothing (A2's runtime leg)."
// REQ-45 / MVV 4: "Hand the DECLARED scalar tag an array literal: still
// refused `flow-tag-invalid` \"is not set-valued\" — now truthfully, and
// the red test pins this beside step 2 so the asymmetry is authored, not
// incidental."
// REQ-46 / MVV 5: "Run `--tag extra=` (undeclared, empty): still refused
// `flow-tag-invalid` \"was given an empty value\" — the one arm the carrier
// does not escape, pinned so the hoist out of `canonicalValue` cannot
// silently drop it."
// REQ-47 / IP Step 3: "Land the red test of MVV steps 2–5 (undeclared
// scalar AND array pass, declared-scalar-given-array still refuses,
// undeclared-empty still refuses)"
// REQ-48 / TS: "Coverage goal — done is all six green, with the
// undeclared/declared pair adjacent in one test so the asymmetry C1 fixes
// is pinned by construction rather than by two tests that could drift"
// HAPPY PATH — the runnable MVV.
func TestMVV_0020_UndeclaredTagKeyAdmission(t *testing.T) {
	// --- MVV 1 — the fixture model ------------------------------------
	//
	// One declared scalar (`tier`), one declared set (`labels`), and no
	// declaration for `extra` or `extras`. The absence IS C1's identity
	// rule, so it is asserted rather than assumed: a model that declared
	// either key would make every step below vacuous.
	model := writeFlowModel(t, carrierTableModel)
	for _, key := range []string{
		carrierUndeclaredScalarKey, carrierUndeclaredArrayKey,
	} {
		if declaresTag(t, carrierTableModel, key) {
			t.Fatalf("MVV 1: the fixture model DECLARES %q; the whole MVV "+
				"rests on it being absent from the tag table", key)
		}
	}
	for _, key := range []string{carrierDeclaredScalar, carrierDeclaredSet} {
		if !declaresTag(t, carrierTableModel, key) {
			t.Fatalf("MVV 1: the fixture model does not declare %q; steps "+
				"3 and 4 need one declared scalar and one declared set", key)
		}
	}

	// --- MVV 2 — both carriers admitted, echoed byte-for-byte ---------
	carrierArgs := carrierResolve(model, withCarriers(
		carrierUndeclaredScalarKey+"="+carrierScalarValue,
		carrierUndeclaredArrayKey+"="+carrierArrayValue,
	)...)

	withStdout := requireSuccess(t, carrierArgs...)
	observed := observedEcho(t, withStdout)

	// The round trip: argv token in, the SAME bytes out. Comparing against
	// the constant the argv was built from is what makes this a fidelity
	// assertion rather than a liveness one.
	roundTrips := map[string]string{
		carrierUndeclaredScalarKey: carrierScalarValue,
		carrierUndeclaredArrayKey:  carrierArrayValue,
	}
	for key, given := range roundTrips {
		got, ok := observed[key]
		if !ok {
			t.Fatalf("MVV 2: `observed` does not carry %q at all; the "+
				"carried value must ride the kernel's Observed view and "+
				"echo in the payload\nobserved: %v", key, observed)
		}
		if got != given {
			t.Errorf("MVV 2: observed[%q] = %q; want %q BYTE-FOR-BYTE as "+
				"given — a value that merely round-trips 'equivalently' "+
				"has lost the caller's spelling, which is the fidelity "+
				"C1 pledges", key, got, given)
		}
	}

	// --- MVV 3 — the same resolve without the carriers ------------------
	//
	// A2's runtime leg: the carried keys influenced nothing. The comparison
	// is over the WHOLE projected selection, not just `rule` — a carrier
	// that perturbed `emit`, `gates`, `writes` or `escaped` while leaving
	// the rule alone would pass a rule-only assertion.
	withoutStdout := requireSuccess(t, carrierResolve(model, declaredPair()...)...)

	withSel := selectionDigest(t, selectionOf(t, withStdout))
	withoutSel := selectionDigest(t, selectionOf(t, withoutStdout))

	if withSel != withoutSel {
		t.Errorf("MVV 3: the selection differs with and without the "+
			"carrier flags\n  with:    %s\n  without: %s\nThe carried "+
			"keys must influence NOTHING — every model-side reference to "+
			"an undeclared tag refuses at load, so nothing that routes can "+
			"read them", withSel, withoutSel)
	}

	// The sole payload delta is `observed` gaining the carried keys
	// (fixture F2). Without that leg, an implementation that dropped the
	// carriers entirely would also pass the equality above.
	baseObserved := observedEcho(t, withoutStdout)
	for key := range roundTrips {
		if _, present := baseObserved[key]; present {
			t.Errorf("MVV 3: the carrier-free run echoes %q; the sole "+
				"payload delta between the two runs must be `observed` "+
				"GAINING the carried keys", key)
		}
	}
	if len(observed) != len(baseObserved)+len(roundTrips) {
		t.Errorf("MVV 3: `observed` went from %d keys to %d; the carrier "+
			"run adds exactly the %d carried keys and changes nothing else"+
			"\n  with:    %v\n  without: %v",
			len(baseObserved), len(observed), len(roundTrips),
			observed, baseObserved)
	}
	for key, want := range baseObserved {
		if got := observed[key]; got != want {
			t.Errorf("MVV 3: observed[%q] = %q with carriers but %q "+
				"without; a carrier must not perturb a declared key's echo",
				key, got, want)
		}
	}

	// --- MVV 4 — the DECLARED scalar handed an array literal -----------
	//
	// Beside step 2, per the record: the asymmetry is authored, not
	// incidental. The same array literal that step 2 CARRIES is refused
	// here because this key is declared.
	declaredArrayArgs := carrierResolve(model,
		carrierDeclaredScalar+"="+carrierArrayValue,
		carrierDeclaredSet+"="+carrierDeclaredSetValue,
	)
	step4Stdout, _ := runCarrier(t, declaredArrayArgs...)
	ce4 := requireRefusal(t, "flow-tag-invalid", 2, declaredArrayArgs...)
	if ce4.Param != carrierDeclaredScalar {
		t.Errorf("MVV 4: param = %q; want %q",
			ce4.Param, carrierDeclaredScalar)
	}
	// Fixture F3's message, byte-for-byte. It is TRUTHFUL now: `tier` is
	// declared, and its declaration is what the message reports.
	const wantF3 = "the tag `" + carrierDeclaredScalar +
		"` is not set-valued; got the array literal " + carrierArrayValue
	if got := refusalMessage(t, step4Stdout); got != wantF3 {
		t.Errorf("MVV 4: message = %q; want %q (normative fixture F3)",
			got, wantF3)
	}

	// --- MVV 5 — the undeclared key given an empty value ---------------
	emptyArgs := carrierResolve(model, withCarriers(
		carrierUndeclaredScalarKey+"=")...)
	step5Stdout, _ := runCarrier(t, emptyArgs...)
	ce5 := requireRefusal(t, "flow-tag-invalid", 2, emptyArgs...)
	if ce5.Param != carrierUndeclaredScalarKey {
		t.Errorf("MVV 5: param = %q; want %q",
			ce5.Param, carrierUndeclaredScalarKey)
	}
	if got := refusalMessage(t, step5Stdout); got != carrierFixtureF4Message {
		t.Errorf("MVV 5: message = %q; want %q (normative fixture F4, "+
			"byte-identical to HEAD) — the hoist out of `canonicalValue` "+
			"must not silently drop or reshape this arm",
			got, carrierFixtureF4Message)
	}
}

// REQ-50 / `0020:S1` normative fixture F1: "for a model declaring scalar
// `tier` and set `labels` and nothing named `extra`: `--tag tier=free
// --tag labels=[\"security\"] --tag extra=plain` ⇒
// `\"observed\":{\"extra\":\"plain\",\"labels\":\"[\\\"security\\\"]\",\"tier\":\"free\"}`."
// REQ-51 / `0020:S1`: "the asserted wire form follows from
// `resolve.Input.Observed`'s `map[string]string` type — the raw argument
// text as a string value, `\"extras\":\"[\\\"a\\\",\\\"b\\\"]\"`, the same
// shape F1 already witnesses for the declared set key `labels`."
// HAPPY PATH — fixture F1 as a WHOLE map, plus the array leg's wire form.
// Asserting the whole map (not three lookups) is what pins "nothing else
// appears in the echo".
func TestReq50And51_0020_FixtureF1AndTheArrayLegsWireForm(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)

	t.Run("F1/scalar-leg", func(t *testing.T) {
		stdout := requireSuccess(t, carrierResolve(model,
			carrierDeclaredScalar+"=free",
			carrierDeclaredSet+"="+carrierDeclaredSetValue,
			carrierUndeclaredScalarKey+"="+carrierScalarValue,
		)...)

		want := map[string]string{
			carrierUndeclaredScalarKey: carrierScalarValue,
			carrierDeclaredSet:         carrierDeclaredSetValue,
			carrierDeclaredScalar:      "free",
		}
		assertObservedEquals(t, observedEcho(t, stdout), want, "fixture F1")
	})

	t.Run("array-leg/raw-argument-text-as-a-string", func(t *testing.T) {
		stdout := requireSuccess(t, carrierResolve(model,
			carrierDeclaredScalar+"=free",
			carrierDeclaredSet+"="+carrierDeclaredSetValue,
			carrierUndeclaredArrayKey+"="+carrierArrayValue,
		)...)

		// `observedEcho` decodes into `map[string]string`, so this fails
		// outright if a carried array were promoted to a JSON array —
		// exactly the wire form REQ-51 rules out.
		want := map[string]string{
			carrierUndeclaredArrayKey: carrierArrayValue,
			carrierDeclaredSet:        carrierDeclaredSetValue,
			carrierDeclaredScalar:     "free",
		}
		assertObservedEquals(t, observedEcho(t, stdout), want,
			"the array leg's wire form (the same STRING shape F1 already "+
				"witnesses for the declared set key)")
	})
}

// REQ-52 / `0020:S2` fixture F2: "`\"rule\":\"free\"`, `\"emit\":{\"plan\":
// \"basic\"}`, with `gates`, `next`, `writes`, `clear` empty and
// `escaped:false` in both runs — every projected field diffs empty and the
// sole payload delta is `observed` gaining the carried key."
// DOMAIN EDGE — F2 pins the ABSOLUTE selection, not merely that two runs
// agree. Two runs that both mis-selected `paid` would satisfy an
// equality-only assertion; this one names the row.
func TestReq52_0020_FixtureF2PinsTheAbsoluteSelectionInBothRuns(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)

	runs := map[string][]string{
		"with-carriers": withCarriers(
			carrierUndeclaredScalarKey+"="+carrierScalarValue,
			carrierUndeclaredArrayKey+"="+carrierArrayValue),
		"without-carriers": declaredPair(),
	}

	for name, tags := range runs {
		t.Run(name, func(t *testing.T) {
			sel := selectionOf(t, requireSuccess(t,
				carrierResolve(model, tags...)...))

			if sel.Rule != "free" {
				t.Errorf("rule = %q; want %q (fixture F2)", sel.Rule, "free")
			}
			if sel.Outcome != "decide" {
				t.Errorf("outcome = %q; want %q", sel.Outcome, "decide")
			}
			if got := sel.Emit["plan"]; got != "basic" {
				t.Errorf("emit[plan] = %q; want %q (fixture F2)",
					got, "basic")
			}
			if sel.Escaped {
				t.Error("escaped = true; fixture F2 fixes it false in both runs")
			}
			if n := len(sel.Gates); n != 0 {
				t.Errorf("gates has %d entries; fixture F2 fixes it empty", n)
			}
			if n := len(sel.Next); n != 0 {
				t.Errorf("next has %d entries; fixture F2 fixes it empty", n)
			}
			if n := len(sel.Writes); n != 0 {
				t.Errorf("writes has %d entries; fixture F2 fixes it empty", n)
			}
			if n := len(sel.Clear); n != 0 {
				t.Errorf("clear has %d entries; fixture F2 fixes it empty", n)
			}
		})
	}
}

// REQ-53 / `0020:S3` fixture F3: "`--tag tier=[\"a\",\"b\"]` ⇒ exit 2,
// `{\"code\":\"flow-tag-invalid\",\"message\":\"the tag `tier` is not
// set-valued; got the array literal [\\\"a\\\",\\\"b\\\"]\",\"param\":
// \"tier\"}`."
// BOUNDARY — fixture F3 on all four wire fields. The message interpolates
// the array literal AS THE CALLER SPELLED IT, which is itself a verbatim
// claim: a refusal that reported a canonicalised literal would misreport
// the input.
func TestReq53_0020_FixtureF3_DeclaredScalarGivenAnArrayLiteral(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)
	args := carrierResolve(model,
		carrierDeclaredScalar+"="+carrierArrayValue,
		carrierDeclaredSet+"="+carrierDeclaredSetValue)

	stdout, _ := runCarrier(t, args...)

	ce := requireRefusal(t, "flow-tag-invalid", 2, args...)
	if ce.Param != carrierDeclaredScalar {
		t.Errorf("param = %q; want %q", ce.Param, carrierDeclaredScalar)
	}
	const want = "the tag `" + carrierDeclaredScalar +
		"` is not set-valued; got the array literal " + carrierArrayValue
	if got := refusalMessage(t, stdout); got != want {
		t.Errorf("message = %q; want %q (normative fixture F3) — the "+
			"literal is reported as the caller spelled it", got, want)
	}
}

// REQ-30 / DISP: "undeclared key, scalar value | 0 | — | `observed` echo,
// verbatim | loud (echoed)"
// REQ-38 / FM: "Diagnosis: the resolve payload's `observed` field echoes
// every carried key byte-for-byte — the stray spelling sits beside the
// declared keys in the same envelope the refusal rides."
// REQ-57 / `0020:G-cross-cutting`: "C1 widens what admission accepts and
// narrows nothing, so no shipped invocation changes meaning: undeclared
// scalars passed before and pass now"
// HAPPY PATH — the shipped behaviour C1 must not narrow, plus the
// diagnostic property the accepted open-world cost is answered with: the
// stray key appears BESIDE the declared keys in one envelope.
func TestReq30And38And57_0020_UndeclaredScalarsStillPassAndSitBesideDeclaredKeys(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)

	stdout := requireSuccess(t, carrierResolve(model, withCarriers(
		carrierUndeclaredScalarKey+"="+carrierScalarValue)...)...)
	observed := observedEcho(t, stdout)

	if got := observed[carrierUndeclaredScalarKey]; got != carrierScalarValue {
		t.Errorf("observed[%q] = %q; want %q — undeclared scalars passed "+
			"before and pass now",
			carrierUndeclaredScalarKey, got, carrierScalarValue)
	}
	// The diagnosis leg: the carrier and the declared keys ride the SAME
	// envelope, which is what makes the echo a usable diagnostic.
	for _, declared := range []string{carrierDeclaredScalar, carrierDeclaredSet} {
		if _, present := observed[declared]; !present {
			t.Errorf("`observed` omits the declared key %q; the stray "+
				"spelling must sit BESIDE the declared keys in the same "+
				"envelope for the echo to be the diagnosis", declared)
		}
	}
}

// REQ-33 / DISP: "undeclared key that is a MISSPELLING of a declared one |
// 0 | — | `observed` echo, verbatim | **silent** — accepted open-world
// cost; diagnosis is the echo (Failure Modes)"
// Q2/D2: OPTIONAL coverage, never an MVV gate. Recorded as the accepted
// cost's witness — its behaviour is fully determined by REQ-1/REQ-31.
// DOMAIN EDGE — the misspelling of a declared SET key with a SET-shaped
// value: the class the record calls out as newly silent (it refused with a
// wrong message at HEAD). "Silent" is asserted as the absence of any
// finding, not merely as exit 0.
func TestReq33_0020_AMisspelledDeclaredKeyPassesSilentlyAsACarrier(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)

	// `labelz` is one byte off the declared set `labels`.
	const misspelled = "labelz"

	stdout := requireSuccess(t, carrierResolve(model, withCarriers(
		misspelled+"="+carrierDeclaredSetValue)...)...)
	observed := observedEcho(t, stdout)

	if got := observed[misspelled]; got != carrierDeclaredSetValue {
		t.Errorf("observed[%q] = %q; want %q verbatim — a misspelling IS "+
			"an undeclared key and is carried, not conformed",
			misspelled, got, carrierDeclaredSetValue)
	}
	// Silent: the declared key it resembles is untouched, and the carrier
	// did not become it.
	if got := observed[carrierDeclaredSet]; got != carrierDeclaredSetValue {
		t.Errorf("observed[%q] = %q; a misspelling must not be folded onto "+
			"the key it resembles", carrierDeclaredSet, got)
	}
}

// --- local oracles -------------------------------------------------------

// declaresTag reports whether the fixture model SOURCE declares a
// `[tags.<key>]` block. MVV 1 is a statement about the fixture, so it is
// asserted against the model document — the input the loader reads — and
// never against the loaded model's internals.
func declaresTag(t *testing.T, src, key string) bool {
	t.Helper()

	return strings.Contains("\n"+src, "\n[tags."+key+"]\n")
}

// assertObservedEquals compares a whole `observed` map against a fixture,
// reporting extras and missing keys separately so a failure names what
// actually diverged.
func assertObservedEquals(t *testing.T, got, want map[string]string, fixture string) {
	t.Helper()

	for key, wantValue := range want {
		gotValue, present := got[key]
		if !present {
			t.Errorf("%s: `observed` is missing %q\n  got: %v",
				fixture, key, got)
			continue
		}
		if gotValue != wantValue {
			t.Errorf("%s: observed[%q] = %q; want %q byte-for-byte",
				fixture, key, gotValue, wantValue)
		}
	}
	for key := range got {
		if _, expected := want[key]; !expected {
			t.Errorf("%s: `observed` carries an unexpected key %q = %q",
				fixture, key, got[key])
		}
	}
}
