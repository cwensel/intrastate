package cli

// RDR 0020 — the hoisted empty-value arm, its GUARD, and admission
// precedence.
//
// C1 moves exactly one arm: the empty-value refusal lifts out of
// `canonicalValue`'s `!isSet` branch onto the shared admission path, so it
// still binds a carrier that no longer enters that function. The move is
// claimed BEHAVIOUR-PRESERVING, and the record names the single input class
// that could falsify the claim — a DECLARED SET key given an empty value,
// which must stay on its set-specific conformance message.
//
// The suite is therefore built around a three-way partition of `value ==
// ""` — undeclared, declared non-set, declared set — because a test that
// only exercised the first two would pass against the unconditional
// `value == ""` hoist that C1 explicitly forbids.

import (
	"testing"
)

// The three normative message fixtures the hoist is pinned against, spelled
// as the record spells them. `0020:G-cross-cutting` (REQ-59) scopes C1's
// "byte-identical" claims to MESSAGE-STRING equality between two refusal
// sites — never to a payload hash — so these are compared as strings.
const (
	// F4 — undeclared key, empty value (`evidence/spikes/d-undeclared-empty.txt`).
	carrierFixtureF4Message = "the tag `extra` was given an empty value"
	// E — declared SCALAR key, empty value (`evidence/spikes/e-declared-empty.txt`).
	carrierFixtureEMessage = "the tag `tier` was given an empty value"
	// G — declared SET key, empty value (`evidence/spikes/g-declared-set-empty.txt`).
	// Note the TRAILING SPACE: the empty value is interpolated after "got ".
	carrierFixtureGMessage = "the tag `labels` is set-valued and takes a " +
		"JSON array literal; got "
)

// REQ-15 / `0020:C1`: "an empty observed value is indistinguishable from
// unset, which is a grammar-level fact about the value that holds with or
// without a declaration. It therefore binds the carrier too, and MOVES to
// the admission path ahead of the carrier branch; today it sits inside
// `canonicalValue`'s `!isSet` arm, which the carrier no longer enters."
// REQ-16 / `0020:C1`: "The hoisted arm carries the SCALAR-shaped message
// for a carrier — “the tag `X` was given an empty value“ — which is
// byte-identical to what an undeclared key already receives at HEAD
// (normative fixture **F4** …), so the hoist is behaviour-preserving."
// REQ-18 / `0020:C1` (the guard): "The hoisted arm is therefore NOT
// unconditional on `value == \"\"`: it fires for a key that is undeclared,
// or declared with a NON-set kind, and it must not intercept a declared SET
// key, whose empty value keeps the set-specific conformance message
// (normative fixture **G** …)"
// REQ-19 / `0020:C1`: "A declared scalar's empty value keeps the scalar
// message it has at HEAD (`evidence/spikes/e-declared-empty.txt`), which is
// the same string the hoisted arm emits, so routing it through either site
// is byte-identical."
// REQ-56 / `0020:S5`: "An unconditional `value == \"\"` test ahead of the
// declared lookup fails this scenario, which is what makes it the guard's
// witness."
// REQ-58 / `0020:G-cross-cutting`: "The one arm that does not widen, the
// empty-value refusal, is hoisted rather than changed and is byte-identical
// for both affected input classes (fixtures F4, E, G), so it needs no
// migration either."
// BOUNDARY — the three-way partition of `value == ""` in ONE table. This is
// the guard's witness: an implementation that hoists on the bare value test
// fails the third row and only the third row.
func TestReq15Through19And56And58_0020_TheHoistedEmptyValueArmIsGuardedByKind(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)

	cases := []struct {
		name    string
		key     string
		tags    []string
		message string
		why     string
	}{
		{
			name:    "undeclared/scalar-shaped-message",
			key:     carrierUndeclaredScalarKey,
			tags:    withCarriers(carrierUndeclaredScalarKey + "="),
			message: carrierFixtureF4Message,
			why: "fixture F4 — the carrier does not escape this arm, and " +
				"the hoisted message is byte-identical to HEAD's",
		},
		{
			name:    "declared-non-set/same-scalar-shaped-message",
			key:     carrierDeclaredScalar,
			tags:    []string{carrierDeclaredScalar + "=", carrierDeclaredSet + "=" + carrierDeclaredSetValue},
			message: carrierFixtureEMessage,
			why: "fixture E — a declared scalar keeps the scalar message, " +
				"so routing it through either site is byte-identical",
		},
		{
			name:    "declared-SET/keeps-its-conformance-message",
			key:     carrierDeclaredSet,
			tags:    []string{carrierDeclaredScalar + "=free", carrierDeclaredSet + "="},
			message: carrierFixtureGMessage,
			why: "fixture G — the GUARD's witness. The set-specific " +
				"empty-value message is a CONFORMANCE message and stays " +
				"declared-only; an unconditional `value == \"\"` hoist " +
				"would move this key off it and regress a declared key, " +
				"a change this RDR does not make",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, _ := runCarrier(t, carrierResolve(model, tc.tags...)...)

			ce := requireRefusal(t, "flow-tag-invalid", 2,
				carrierResolve(model, tc.tags...)...)
			if ce.Param != tc.key {
				t.Errorf("param = %q; want %q", ce.Param, tc.key)
			}
			if got := refusalMessage(t, stdout); got != tc.message {
				t.Errorf("message = %q;\n   want %q\n   — %s",
					got, tc.message, tc.why)
			}
		})
	}
}

// REQ-17 / `0020:C1`: "The set-specific empty-value message (\"is
// set-valued and takes a JSON array literal; got \") is a CONFORMANCE
// message and stays declared-only, inside the arms below: it presupposes a
// declared kind, which a carrier by definition has none of."
// REQ-20 / `0020:C1` (negative): "Hoisting on the bare value test alone
// would move a declared set key off its conformance message — a
// declared-key change this RDR does not make." — a test asserting the
// unconditional form FAILS.
// ADVERSARIAL — the converse of the guard, asserted from the carrier's
// side: a CARRIER must never receive the conformance message. The message
// presupposes a declared kind, so emitting it for a key with no declaration
// would be the CLI reporting a declaration that does not exist — the exact
// defect this RDR exists to fix.
func TestReq17And20_0020_ACarrierNeverReceivesTheSetConformanceMessage(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)

	// The value classes that reach the set-conformance arms for a DECLARED
	// set: empty (the "got " arm), a bare scalar (the same arm), and a
	// malformed array (the parse arm). None may fire for a carrier.
	for _, value := range []string{"", "bare-scalar", `["a",`, `[1,2]`} {
		t.Run("value="+value, func(t *testing.T) {
			stdout, err := runCarrier(t, carrierResolve(model, withCarriers(
				carrierUndeclaredScalarKey+"="+value)...)...)

			if value == "" {
				// The one refusal a carrier does not escape — but on the
				// SCALAR-shaped message, never the conformance one.
				if err == nil {
					t.Fatalf("an empty carrier value was ADMITTED; the "+
						"empty-value arm binds the carrier too\n%s", stdout)
				}
				if got := refusalMessage(t, stdout); got != carrierFixtureF4Message {
					t.Errorf("message = %q; want %q — the CONFORMANCE "+
						"message presupposes a declared kind, which a "+
						"carrier by definition has none of",
						got, carrierFixtureF4Message)
				}
				return
			}

			if err != nil {
				t.Fatalf("a carrier value %q was REFUSED; the conformance "+
					"arms are declared-only: %v\n%s", value, err, stdout)
			}
		})
	}
}

// REQ-21 / `0020:D-selection-predicate`: "refusal precedence at admission is
// unchanged in order: grammar → reserved → owned → duplicate → empty-value
// (undeclared or non-set-declared keys) → (declared keys only) kind/shape/
// domain conformance;"
// REQ-24 / `0020:D-selection-predicate`: "The order above is the whole
// constraint on the splice point; where the arm sits within it is
// implementation latitude."
// REQ-37 / FM: "every refusal that survives is unchanged and still fires at
// exit 2 before any accessor"
// ADVERSARIAL — precedence is only observable when TWO arms are armed at
// once and exactly one code comes back. Each row below arms an earlier arm
// and a later one in a single invocation; asserting the later code would
// pass a reordering, so the earlier one is what is pinned. The order is the
// whole constraint — nothing here asserts a splice POSITION.
func TestReq21And24And37_0020_AdmissionRefusalPrecedenceIsUnchanged(t *testing.T) {
	tableModel := writeFlowModel(t, carrierTableModel)

	cases := []struct {
		name string
		tags []string
		code string
		beat string
	}{
		{
			name: "grammar-beats-reserved",
			tags: withCarriers("recognized-with-no-equals"),
			code: "flow-tag-invalid",
			beat: "grammar precedes reserved",
		},
		{
			name: "reserved-beats-duplicate",
			tags: withCarriers("recognized=decide", "recognized=decide"),
			code: "flow-tag-reserved",
			beat: "reserved precedes duplicate",
		},
		{
			name: "reserved-beats-empty-value",
			tags: withCarriers("recognized="),
			code: "flow-tag-reserved",
			beat: "reserved precedes the empty-value arm",
		},
		{
			name: "duplicate-beats-empty-value/undeclared",
			tags: withCarriers(
				carrierUndeclaredScalarKey+"=plain",
				carrierUndeclaredScalarKey+"="),
			code: "flow-tag-duplicate",
			beat: "duplicate precedes the empty-value arm",
		},
		{
			name: "duplicate-beats-conformance/declared-scalar",
			tags: []string{
				carrierDeclaredScalar + "=free",
				carrierDeclaredSet + "=" + carrierDeclaredSetValue,
				carrierDeclaredScalar + `=["a","b"]`,
			},
			code: "flow-tag-duplicate",
			beat: "duplicate precedes the declared-key conformance arms",
		},
		{
			name: "empty-value-beats-conformance/declared-scalar",
			tags: []string{
				carrierDeclaredScalar + "=",
				carrierDeclaredSet + "=" + carrierDeclaredSetValue,
			},
			code: "flow-tag-invalid",
			beat: "the empty-value arm precedes the domain check " +
				"(an empty value is not in the enum's domain either)",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			requireRefusal(t, tc.code, 2,
				carrierResolve(tableModel, tc.tags...)...)
			t.Logf("precedence held: %s", tc.beat)
		})
	}

	// Owned sits between reserved and duplicate, and needs a model that
	// declares an owned tag — RDR 0005's shipped MVV fixture.
	ownedModel := writeFlowModel(t, flowMVVModel)
	ownedBind := artifactBinding(flowStateRole,
		seedArtifact(t, ownedModel, "status=draft"))

	t.Run("owned-beats-duplicate", func(t *testing.T) {
		requireRefusal(t, "flow-tag-owned", 2,
			"flow", "resolve", "--model", ownedModel,
			"--artifact", ownedBind, "--outcome", "hold",
			"--tag", carrierOwnedKey+"=draft",
			"--tag", carrierOwnedKey+"=final", "--as=json")
	})

	t.Run("owned-beats-empty-value", func(t *testing.T) {
		requireRefusal(t, "flow-tag-owned", 2,
			"flow", "resolve", "--model", ownedModel,
			"--artifact", ownedBind, "--outcome", "hold",
			"--tag", carrierOwnedKey+"=", "--as=json")
	})

	t.Run("reserved-beats-owned", func(t *testing.T) {
		// Both arms armed, reserved second in argv order so a
		// first-match-wins loop cannot fake the result. Reserved must win.
		requireRefusal(t, "flow-tag-reserved", 2,
			"flow", "resolve", "--model", ownedModel,
			"--artifact", ownedBind, "--outcome", "hold",
			"--tag", "recognized=hold",
			"--tag", carrierOwnedKey+"=draft", "--as=json")
	})
}

// REQ-22 / `0020:D-selection-predicate`: "the carrier branch sits where the
// conformance arms would have run."
// DOMAIN EDGE — a splice position is not directly observable, but its
// CONSEQUENCE is: the carrier branch is downstream of every guard and
// upstream of nothing, so a carrier that would have tripped a conformance
// arm is admitted while every guard that precedes those arms still fires on
// the same invocation. Asserted as the pair — one argv where the carrier
// passes and a guard on a DIFFERENT key still refuses.
func TestReq22_0020_TheCarrierBranchSitsWhereTheConformanceArmsWouldHaveRun(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)

	// The carrier's value would trip the shape arm if the conformance arms
	// ran for it; the duplicate on a declared key is upstream of those arms
	// and must still be what refuses.
	requireRefusal(t, "flow-tag-duplicate", 2,
		carrierResolve(model,
			carrierDeclaredScalar+"=free",
			carrierDeclaredSet+"="+carrierDeclaredSetValue,
			carrierUndeclaredArrayKey+"="+carrierArrayValue,
			carrierDeclaredScalar+"=paid",
		)...)

	// The same carrier, with no guard armed, is admitted — so the refusal
	// above was genuinely the guard's and not the carrier's.
	stdout := requireSuccess(t, carrierResolve(model, withCarriers(
		carrierUndeclaredArrayKey+"="+carrierArrayValue)...)...)
	if got := observedEcho(t, stdout)[carrierUndeclaredArrayKey]; got != carrierArrayValue {
		t.Errorf("observed[%q] = %q; want %q — the carrier branch replaces "+
			"the conformance arms for an undeclared key",
			carrierUndeclaredArrayKey, got, carrierArrayValue)
	}
}

// REQ-23 / `0020:D-selection-predicate` (the anti-inversion rule): "It lands
// AFTER the duplicate arm and the `seen[key]` mark, not before: at HEAD the
// duplicate check already precedes the empty-value site … so a repeated key
// refuses `flow-tag-duplicate` today and must still do so — hoisting the
// arm ahead of the duplicate case would invert that precedence and change
// behaviour this RDR does not touch."
// ADVERSARIAL — the one inversion the hoist could introduce, isolated. The
// discriminating construction is a repeat whose SECOND occurrence is empty:
// the duplicate arm can only reach it if the hoisted empty-value arm sits
// after `seen[key]`. An implementation that hoisted to the top of the loop
// body would answer `flow-tag-invalid` here.
func TestReq23_0020_TheHoistedArmLandsAfterTheDuplicateMarkNotBefore(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)

	cases := []struct {
		name string
		tags []string
	}{
		{
			name: "undeclared-key/second-occurrence-empty",
			tags: withCarriers(
				carrierUndeclaredScalarKey+"=plain",
				carrierUndeclaredScalarKey+"="),
		},
		{
			name: "declared-scalar/second-occurrence-empty",
			tags: []string{
				carrierDeclaredScalar + "=free",
				carrierDeclaredSet + "=" + carrierDeclaredSetValue,
				carrierDeclaredScalar + "=",
			},
		},
		{
			name: "declared-set/second-occurrence-empty",
			tags: []string{
				carrierDeclaredScalar + "=free",
				carrierDeclaredSet + "=" + carrierDeclaredSetValue,
				carrierDeclaredSet + "=",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ce := requireRefusal(t, "flow-tag-duplicate", 2,
				carrierResolve(model, tc.tags...)...)
			if ce.Code != "flow-tag-duplicate" {
				t.Errorf("code = %q; hoisting the empty-value arm ahead of "+
					"the duplicate case inverts a precedence this RDR does "+
					"not touch", ce.Code)
			}
		})
	}
}

// REQ-35 / DISP (the guard's witness): "declared set key, empty value | 2 |
// `flow-tag-invalid` \"is set-valued and takes a JSON array literal; got \"
// | none | loud"
// REQ-55 / `0020:S5` fixture G: "`--tag labels=` ⇒ exit 2,
// `{\"code\":\"flow-tag-invalid\",\"message\":\"the tag `labels` is
// set-valued and takes a JSON array literal; got \",\"param\":\"labels\"}`
// — byte-identical to HEAD."
// BOUNDARY — fixture G on all four wire fields, standing alone so the
// regression it guards against fails a test that names it. The trailing
// space in the message is load-bearing: the empty value is interpolated,
// and a message trimmed or reworded is a declared-key change.
func TestReq35And55_0020_FixtureG_DeclaredSetEmptyKeepsItsConformanceMessage(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)
	args := carrierResolve(model,
		carrierDeclaredScalar+"=free", carrierDeclaredSet+"=")

	stdout, _ := runCarrier(t, args...)

	ce := requireRefusal(t, "flow-tag-invalid", 2, args...)
	if ce.Param != carrierDeclaredSet {
		t.Errorf("param = %q; want %q", ce.Param, carrierDeclaredSet)
	}
	if got := refusalMessage(t, stdout); got != carrierFixtureGMessage {
		t.Errorf("message = %q;\n   want %q (normative fixture G, "+
			"byte-identical to HEAD — note the TRAILING SPACE, the empty "+
			"value interpolated after \"got \")", got, carrierFixtureGMessage)
	}
}
