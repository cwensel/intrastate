package cli

// RDR 0020 — `0020:C1`, the carrier rule and its closed refusal list.
//
// C1 decides ONE thing: what the zero `TagDecl` means at `--tag` admission.
// Every test here is derived from a clause of that fence (or from the
// Load-Bearing Decisions that fix its identity rule and its step order),
// and each is written to FAIL if that clause is broken — not merely to
// observe what the seam happens to do.
//
// The suite's discriminating choice: declaredness is varied by AUTHORING
// it, never by inspecting `m.Tags`. C1's identity rule is a statement about
// which INPUTS reach which arm, so a test that read the map would be
// pinning the predicate's implementation instead of the partition it
// induces.

import (
	"strings"
	"testing"
)

// REQ-1 / `0020:C1`: "PURE CARRIER. At `--tag` admission, a key ABSENT from
// the loaded model's normalized tag table (the two-value `m.Tags[key]`
// lookup) is a pure carrier: the CLI admits its value VERBATIM —
// byte-preserved, JSON array literals included — and never canonicalises,
// conforms, kind-checks, or compares it."
// REQ-31 / DISP: "undeclared key, array literal | 0 (**changed**; refuses
// at HEAD) | — | `observed` echo, verbatim | loud (echoed)"
// REQ-49 / `0020:S1`: "undeclared scalar and undeclared array literal
// admitted (MVV 2). Expected: exit 0; `observed` echoes both
// byte-for-byte."
// HAPPY PATH — the change itself, with the scalar leg beside the array leg
// in ONE invocation, per the Testing Strategy's "adjacent in one test so
// the asymmetry C1 fixes is pinned by construction".
func TestReq1And31And49_0020_UndeclaredScalarAndArrayAreAdmittedVerbatim(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)

	stdout := requireSuccess(t, carrierResolve(model, withCarriers(
		carrierUndeclaredScalarKey+"="+carrierScalarValue,
		carrierUndeclaredArrayKey+"="+carrierArrayValue,
	)...)...)

	observed := observedEcho(t, stdout)

	// The scalar leg — shipped behaviour, pinned so the carrier branch
	// cannot regress it while widening.
	if got := observed[carrierUndeclaredScalarKey]; got != carrierScalarValue {
		t.Errorf("observed[%q] = %q; want %q — an undeclared key is a PURE "+
			"CARRIER and its value crosses VERBATIM",
			carrierUndeclaredScalarKey, got, carrierScalarValue)
	}

	// The array leg — the whole change. `["a","b"]` must arrive as the
	// caller spelled it: not canonicalised, not re-encoded, not promoted to
	// a JSON array.
	if got := observed[carrierUndeclaredArrayKey]; got != carrierArrayValue {
		t.Errorf("observed[%q] = %q; want %q byte-for-byte — an undeclared "+
			"array literal is admitted VERBATIM, never canonicalised",
			carrierUndeclaredArrayKey, got, carrierArrayValue)
	}
}

// REQ-2 / `0020:C1`: "The carried value rides the kernel's Observed view
// and echoes in the resolve payload's `observed` field byte-for-byte as
// given"
// REQ-6 / `0020:G-cross-cutting`: "A carried value is admitted VERBATIM —
// byte-preserved, with no canonicalisation, folding, normalisation, or
// re-encoding — and echoes in the resolve payload's `observed` field as
// given."
// REQ-26 / TD: "The carried value flows exactly where an undeclared scalar
// already flows today … out through the resolve payload's `observed` echo,
// byte-for-byte as given."
// ADVERSARIAL — the four transformations C1 names by name, each given a
// value that a canonicaliser, a folder, a normaliser, or a re-encoder would
// visibly change. A test using only `plain` would pass against every one of
// them.
func TestReq2And6And26_0020_NoCanonicalisationFoldingNormalisationOrReencoding(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)

	cases := []struct {
		name  string
		value string
		what  string
	}{
		{
			name:  "unsorted-array-is-not-sorted",
			value: `["b","a"]`,
			what: "canonicalisation (the canonical set form is SORTED, so a " +
				"canonicaliser would return [\"a\",\"b\"])",
		},
		{
			name:  "duplicated-array-is-not-deduplicated",
			value: `["a","a"]`,
			what: "canonicalisation (the canonical set form is " +
				"DUPLICATE-FREE, so a canonicaliser would return [\"a\"])",
		},
		{
			name:  "spaced-array-keeps-its-whitespace",
			value: `[ "a" , "b" ]`,
			what: "canonicalisation (the canonical set form is COMPACT, so " +
				"a canonicaliser would strip the interior spaces)",
		},
		{
			name:  "case-is-not-folded",
			value: "MiXeD-CaSe",
			what:  "case folding",
		},
		{
			name:  "html-metacharacters-are-not-escaped",
			value: `["a<b","x&y"]`,
			what: "re-encoding (bare json.Marshal would emit " +
				"\\u003c and \\u0026)",
		},
		{
			name:  "surrounding-whitespace-is-not-trimmed",
			value: "  padded  ",
			what:  "normalisation (a trimmer would strip the padding)",
		},
		{
			name:  "a-bracket-prefixed-scalar-is-still-a-scalar",
			value: "[not json at all",
			what: "kind-checking on a first-byte sniff — C1 forbids " +
				"interpreting bytes the system pledges not to interpret",
		},
		{
			name:  "an-empty-json-array-is-not-elided",
			value: "[]",
			what:  "canonicalisation of the empty set",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stdout := requireSuccess(t, carrierResolve(model, withCarriers(
				carrierUndeclaredScalarKey+"="+tc.value)...)...)

			got := observedEcho(t, stdout)[carrierUndeclaredScalarKey]
			if got != tc.value {
				t.Errorf("observed[%q] = %q; want %q byte-for-byte — the "+
					"echo shows %s",
					carrierUndeclaredScalarKey, got, tc.value, tc.what)
			}
		})
	}
}

// REQ-3 / `0020:D-identity`: "\"declared\" means the key is present,
// byte-exact, in the normalized model's tag table (`m.Tags`, two-value
// lookup) … No parallel registry, no case folding, no zero-decl sentinel."
// DOMAIN EDGE — "byte-exact" and "no case folding" are assertions about the
// KEY, and the only way to observe them from the CLI is to hand admission a
// key that differs from a declared one by case or by a byte: it must be
// treated as UNDECLARED (carried) and not as the declared key it resembles.
func TestReq3_0020_DeclarednessIsByteExactOnTheKeyWithNoCaseFolding(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)

	// Each key differs from a DECLARED key (`tier`, a single-valued enum
	// with domain ["free","paid"]) only in case or by one byte. If
	// declaredness folded case or matched loosely, admission would run the
	// declared enum's domain check and refuse `bogus`.
	nearMisses := []string{"TIER", "Tier", "tier ", "tie", "tierr", "LABELS"}

	for _, key := range nearMisses {
		t.Run(key, func(t *testing.T) {
			stdout, err := runCarrier(t, carrierResolve(model, withCarriers(
				key+"=bogus")...)...)
			if err != nil {
				t.Fatalf("key %q was treated as DECLARED; C1's identity rule "+
					"is BYTE-EXACT presence in `m.Tags` with no case "+
					"folding, so %q is undeclared and must be CARRIED: %v\n%s",
					key, key, err, stdout)
			}
			if got := observedEcho(t, stdout)[key]; got != "bogus" {
				t.Errorf("observed[%q] = %q; want %q — a near-miss key is "+
					"an ordinary carrier", key, got, "bogus")
			}
		})
	}
}

// REQ-4 / `0020:D-wire-byte-format`: "an undeclared value crosses verbatim,
// never re-canonicalised: canonicalising would interpret as a set a value
// whose declaration never asserted set-ness."
// REQ-5 / `0020:C1`: "declared set values keep the canonical-array form of
// `docs/cli-output-contract.md` §Set values on the wire"
// BOUNDARY — the two halves of the wire rule in ONE invocation, so the
// asymmetry is pinned by construction: the DECLARED set is re-canonicalised
// (sorted, deduplicated, compact) and the UNDECLARED array beside it is
// not. A change that canonicalised both, or neither, fails here.
func TestReq4And5_0020_DeclaredSetIsCanonicalisedAndTheCarrierIsNot(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)

	// Both values are unsorted and duplicated in the SAME way. Only the
	// declared one may come back changed.
	const messy = `["docs","security","docs"]`
	const canonical = `["docs","security"]`

	stdout := requireSuccess(t, carrierResolve(model,
		carrierDeclaredScalar+"=free",
		carrierDeclaredSet+"="+messy,
		carrierUndeclaredArrayKey+"="+messy,
	)...)
	observed := observedEcho(t, stdout)

	if got := observed[carrierDeclaredSet]; got != canonical {
		t.Errorf("observed[%q] = %q; want the CANONICAL form %q — a "+
			"DECLARED set keeps the §Set-values-on-the-wire canonical array",
			carrierDeclaredSet, got, canonical)
	}
	if got := observed[carrierUndeclaredArrayKey]; got != messy {
		t.Errorf("observed[%q] = %q; want %q unchanged — canonicalising an "+
			"undeclared value would interpret as a set a value whose "+
			"declaration never asserted set-ness",
			carrierUndeclaredArrayKey, got, messy)
	}
}

// REQ-7 / `0020:C1`: "The only refusals reachable for an undeclared key
// are, unchanged and still preceding any accessor (REQ-27):" — the list is
// EXHAUSTIVE; any other refusal firing over an undeclared key is a
// violation.
// ASSUMPTION-4 scopes the list to the ADMISSION path's own refusals for
// that key, which is why every case below differs from a passing
// invocation ONLY in the carrier's value.
// ADVERSARIAL — the closure claim, asserted in the direction that can
// actually fail: values that would trip a kind, shape, or domain arm if one
// were reachable. Every one must be ADMITTED.
func TestReq7_0020_TheOnlyRefusalsForAnUndeclaredKeyAreTheClosedList(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)

	// Each value is chosen to trip an arm that exists at this seam for
	// DECLARED keys. None of them may fire for a key the model does not
	// declare.
	outsideTheList := []struct {
		name  string
		value string
		arm   string
	}{
		{"array-literal", `["a","b"]`, "the set/scalar shape arm"},
		{"malformed-array", `["a",`, "the JSON-array-of-strings parse arm"},
		{"array-of-non-strings", `[1,2]`, "the set-member parse arm"},
		{"out-of-enum-domain", "not-free-not-paid", "the enum domain arm"},
		{"non-integer", "not-a-number", "the int min/max arm"},
		{"non-boolean", "maybe", "the bool literal arm"},
		{"clear-sentinel-text", "<clear>", "the write sentinel arm"},
		{"leading-bracket-scalar", "[unclosed", "the looksArray sniff"},
	}

	for _, tc := range outsideTheList {
		t.Run(tc.name, func(t *testing.T) {
			stdout, err := runCarrier(t, carrierResolve(model, withCarriers(
				carrierUndeclaredScalarKey+"="+tc.value)...)...)
			if err != nil {
				t.Fatalf("an undeclared key carrying %q was REFUSED; C1's "+
					"refusal list is closed and does not include %s: %v\n%s",
					tc.value, tc.arm, err, stdout)
			}
			if got := observedEcho(t, stdout)[carrierUndeclaredScalarKey]; got != tc.value {
				t.Errorf("observed[%q] = %q; want %q verbatim",
					carrierUndeclaredScalarKey, got, tc.value)
			}
		})
	}
}

// REQ-8 / `0020:C1`: "the flag grammar (`--tag` takes `name=value`),
// `flow-tag-invalid`;"
// ASSUMPTION-6: "'Verbatim' is scoped to what admission received after the
// grammar split, not to the whole argv token" — `--tag k=v` is split on the
// FIRST `=`, so a value containing `=` carries its remainder verbatim.
// INPUT EDGE — the grammar arm is IN the closed list and must still fire,
// and the split it performs must not move.
func TestReq8_0020_TheGrammarArmStillRefusesAndTheFirstEqualsSplitHolds(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)

	t.Run("no-equals-refuses-flow-tag-invalid", func(t *testing.T) {
		ce := requireRefusal(t, "flow-tag-invalid", 2,
			carrierResolve(model, withCarriers("noequalssign")...)...)
		if ce.Param != "tag" {
			t.Errorf("param = %q; want %q — the grammar arm names the FLAG, "+
				"not a key", ce.Param, "tag")
		}
	})

	t.Run("empty-key-refuses-flow-tag-invalid", func(t *testing.T) {
		requireRefusal(t, "flow-tag-invalid", 2,
			carrierResolve(model, withCarriers("=orphaned-value")...)...)
	})

	// The split point is the FIRST `=`: the remainder — including further
	// `=` bytes — is the value, and it crosses verbatim.
	t.Run("value-keeps-its-own-equals-signs", func(t *testing.T) {
		const value = "a=b=c"
		stdout := requireSuccess(t, carrierResolve(model, withCarriers(
			carrierUndeclaredScalarKey+"="+value)...)...)
		if got := observedEcho(t, stdout)[carrierUndeclaredScalarKey]; got != value {
			t.Errorf("observed[%q] = %q; want %q — the value is everything "+
				"after the FIRST `=`, carried verbatim",
				carrierUndeclaredScalarKey, got, value)
		}
	})
}

// REQ-9 / `0020:C1`: "`flow-tag-reserved` (REQ-26), `flow-tag-owned`
// (REQ-27), `flow-tag-duplicate` (REQ-28);"
// REQ-29 / AP: "The provenance guards are untouched and still precede
// admission: reserved key, owned key, duplicate (REQ-26/27/28)."
// REQ-36 / DISP: "reserved / owned / duplicate key, any value | 2 |
// `flow-tag-reserved` / `-owned` / `-duplicate` | none | loud (unchanged,
// still first)"
// ADVERSARIAL — "any value" is the load-bearing word: each guard is
// provoked with an ARRAY LITERAL, the one value class C1 newly admits. A
// carrier branch spliced in ahead of the guards would let the array through
// and lose the provenance refusal.
func TestReq9And29And36_0020_ProvenanceGuardsStillFireFirstForAnyValue(t *testing.T) {
	tableModel := writeFlowModel(t, carrierTableModel)

	t.Run("reserved", func(t *testing.T) {
		for _, value := range []string{"decide", `["a","b"]`, ""} {
			ce := requireRefusal(t, "flow-tag-reserved", 2,
				carrierResolve(tableModel, withCarriers(
					"recognized="+value)...)...)
			if ce.Param != "recognized" {
				t.Errorf("param = %q; want %q", ce.Param, "recognized")
			}
		}
	})

	t.Run("duplicate", func(t *testing.T) {
		// A repeated UNDECLARED key: the duplicate arm must still see it,
		// which it can only do if the carrier branch sits after the
		// `seen[key]` mark. The ARRAY value is the discriminating one — at
		// HEAD the first occurrence never survives to be duplicated,
		// because the conformance arm refuses it before the second is
		// read; under C1 the first is CARRIED and the second is the
		// duplicate.
		for _, value := range []string{"plain", `["a","b"]`} {
			ce := requireRefusal(t, "flow-tag-duplicate", 2,
				carrierResolve(tableModel, withCarriers(
					carrierUndeclaredScalarKey+"="+value,
					carrierUndeclaredScalarKey+"="+value)...)...)
			if ce.Param != carrierUndeclaredScalarKey {
				t.Errorf("param = %q; want %q",
					ce.Param, carrierUndeclaredScalarKey)
			}
		}
	})

	t.Run("owned", func(t *testing.T) {
		ownedModel := writeFlowModel(t, flowMVVModel)
		bind := artifactBinding(flowStateRole,
			seedArtifact(t, ownedModel, "status=draft"))

		for _, value := range []string{"draft", `["draft"]`, ""} {
			ce := requireRefusal(t, "flow-tag-owned", 2,
				"flow", "resolve", "--model", ownedModel,
				"--artifact", bind, "--outcome", "hold",
				"--tag", carrierOwnedKey+"="+value, "--as=json")
			if ce.Param != carrierOwnedKey {
				t.Errorf("param = %q; want %q", ce.Param, carrierOwnedKey)
			}
		}
	})
}

// REQ-10 / `0020:C1`: "the empty-value arm, `flow-tag-invalid`"
// REQ-32 / DISP: "undeclared key, empty value | 2 | `flow-tag-invalid`
// \"was given an empty value\" | none | loud"
// REQ-54 / `0020:S4` fixture F4: "`--tag extra=` ⇒ exit 2,
// `{\"code\":\"flow-tag-invalid\",\"message\":\"the tag `extra` was given
// an empty value\",\"param\":\"extra\"}` — byte-identical to HEAD".
// BOUNDARY — the one arm the carrier does NOT escape, pinned on all three
// fields (code, exit, message) so the hoist out of `canonicalValue` cannot
// silently drop it or reshape its message.
func TestReq10And32And54_0020_AnUndeclaredKeysEmptyValueStillRefuses(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)

	stdout, _ := runCarrier(t, carrierResolve(model, withCarriers(
		carrierUndeclaredScalarKey+"=")...)...)

	ce := requireRefusal(t, "flow-tag-invalid", 2,
		carrierResolve(model, withCarriers(
			carrierUndeclaredScalarKey+"=")...)...)
	if ce.Param != carrierUndeclaredScalarKey {
		t.Errorf("param = %q; want %q",
			ce.Param, carrierUndeclaredScalarKey)
	}

	// Fixture F4's message, byte-for-byte.
	const wantF4 = "the tag `" + carrierUndeclaredScalarKey +
		"` was given an empty value"
	if got := refusalMessage(t, stdout); got != wantF4 {
		t.Errorf("message = %q; want %q (normative fixture F4, "+
			"byte-identical to HEAD — the hoist is behaviour-preserving)",
			got, wantF4)
	}
}

// REQ-11 / `0020:C1`: "`flow-tag-invalid`'s kind, shape, and domain arms —
// including \"is not set-valued\" — are reachable only for a DECLARED key,
// whose loaded declaration (kind required at load, one of RDR 0003's five
// tokens) is what the message then truthfully reports."
// REQ-14 / `0020:C1`: "Declaring a key later is the opt-in tightening:
// admission then enforces that declaration's kind and domain."
// REQ-34 / DISP: "declared scalar/enum key, array literal | 2 |
// `flow-tag-invalid` \"is not set-valued\" — now truthful | none | loud"
// BOUNDARY — the SAME key and the SAME value, run against two models that
// differ only in whether the key is declared. That is the only construction
// that can prove reachability is governed by declaredness rather than by
// the value: one leg passes, the other refuses, and neither leg alone would
// discriminate.
func TestReq11And14And34_0020_TheKindArmIsReachableOnlyUnderARealDeclaration(t *testing.T) {
	// `extras` is declared NOWHERE in the carrier fixture.
	undeclaredModel := writeFlowModel(t, carrierTableModel)

	// The same fixture with `extras` DECLARED as a single-valued enum — the
	// opt-in tightening REQ-14 names.
	declaredModel := writeFlowModel(t, carrierTableModel+`
[tags.`+carrierUndeclaredArrayKey+`]
provenance = "observed"
kind = "enum"
domain = ["a", "b"]
single_valued = true
`)

	tag := carrierUndeclaredArrayKey + "=" + carrierArrayValue

	t.Run("undeclared/admitted", func(t *testing.T) {
		stdout := requireSuccess(t,
			carrierResolve(undeclaredModel, withCarriers(tag)...)...)
		if got := observedEcho(t, stdout)[carrierUndeclaredArrayKey]; got != carrierArrayValue {
			t.Errorf("observed[%q] = %q; want %q — with no declaration "+
				"there is no kind to check",
				carrierUndeclaredArrayKey, got, carrierArrayValue)
		}
	})

	t.Run("declared/refused-truthfully", func(t *testing.T) {
		stdout, _ := runCarrier(t,
			carrierResolve(declaredModel, withCarriers(tag)...)...)

		ce := requireRefusal(t, "flow-tag-invalid", 2,
			carrierResolve(declaredModel, withCarriers(tag)...)...)
		if ce.Param != carrierUndeclaredArrayKey {
			t.Errorf("param = %q; want %q",
				ce.Param, carrierUndeclaredArrayKey)
		}
		// The message is TRUTHFUL only because a declaration now says so.
		if msg := refusalMessage(t, stdout); !strings.Contains(msg, "is not set-valued") {
			t.Errorf("message = %q; want it to carry \"is not set-valued\" "+
				"— declaring the key is the opt-in tightening, and the "+
				"message now truthfully reports the loaded declaration", msg)
		}
	})
}

// REQ-12 / `0020:C1`: "No new refusal code is minted; `flow-tag-undeclared`
// does not exist."
// `0020:D-naming`: "no new code; the rejected alternative's
// `flow-tag-undeclared` is deliberately not minted"
// ADVERSARIAL — the negative, asserted two ways because either alone is
// weak: the spelling must not appear in the shipped code-string surface,
// and no invocation over an undeclared key may produce it.
func TestReq12_0020_NoFlowTagUndeclaredCodeIsMinted(t *testing.T) {
	const rejected = "flow-tag-undeclared"

	t.Run("the-spelling-is-absent-from-the-code-table", func(t *testing.T) {
		// The stable code spellings live in the CLI's constant block; a
		// minted code would have to be spelled somewhere in this package.
		src := readRepoFile(t, "internal/cli/flow_input.go")
		if strings.Contains(src, rejected) {
			t.Errorf("`internal/cli/flow_input.go` mints %q; C1 mints NO "+
				"new refusal code — the rejected alternative's code is "+
				"deliberately not minted", rejected)
		}
	})

	t.Run("no-undeclared-key-invocation-produces-it", func(t *testing.T) {
		model := writeFlowModel(t, carrierTableModel)
		// Every input class an undeclared key can present, including the
		// one class that DOES refuse.
		for _, value := range []string{"plain", `["a","b"]`, "", "[bad"} {
			stdout, err := runCarrier(t, carrierResolve(model, withCarriers(
				carrierUndeclaredScalarKey+"="+value)...)...)
			if err != nil && strings.Contains(stdout, rejected) {
				t.Errorf("value %q produced %q; C1 mints no new code",
					value, rejected)
			}
		}
	})
}

// REQ-13 / `0020:C1`: "The model-side closed world is untouched: a rule
// atom, accessor key, write target, clear target, or `[initial]` assignment
// naming an undeclared tag still refuses at load (`unknown_tag`)."
// DOMAIN EDGE — the widening is CALLER-side only. The safety argument for
// the carrier arm (A2: nothing that routes can read the carried bytes)
// rests entirely on this closure, so a change that widened the model side
// too would silently make carriers semantically live.
func TestReq13_0020_TheModelSideClosedWorldStillRefusesAnUndeclaredTag(t *testing.T) {
	// A rule guard atom naming `extra`, which the model does not declare.
	model := writeFlowModel(t, carrierUnknownTagAtomModel)

	requireRefusal(t, "flow-model-invalid", 2,
		"flow", "resolve", "--model", model,
		"--outcome", "decide", "--as=json")
}
