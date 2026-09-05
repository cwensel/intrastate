package cli

// RDR 0020 — the fixture corpus for undeclared `--tag` key admission.
//
// `0020:C1` decides what the zero `TagDecl` MEANS at the admission seam, so
// every fixture here exists to make the declared/undeclared distinction
// observable through the CLI and nothing else. Two rules hold:
//
//  1. NOTHING here inspects `m.Tags`, `parseTags`, or `canonicalValue`
//     directly. "Declared" is C1's identity rule — presence in the
//     normalized model's tag table — and a test that reached into the map
//     would be asserting the implementation of the predicate rather than
//     the behaviour the predicate selects. Declaredness is instead varied
//     by AUTHORING it in the fixture model and observed through the exit
//     code, the refusal identity, and the payload's `observed` echo.
//
//  2. The model is a `decision-table`, so no reader, no `--artifact`, and
//     no owned state stand between the argv and the admission seam. A
//     refusal reaching this suite is therefore `parseTags`'s and not an
//     accessor's — which is what lets REQ-7's closed refusal list be
//     asserted as CLOSED.
//
// The corpus mirrors the RDR's normative spike fixture
// (`evidence/spikes/carrier-table.toml`, from which fixtures F1–F4 and G
// were read at HEAD): declared enum scalar `tier` guarded by both rows,
// declared set `labels` unguarded, and NOTHING named `extra` or `extras`.

import (
	"encoding/json"
	"strings"
	"testing"
)

// --- the carrier fixture identities -------------------------------------
//
// These are the four key classes C1's identity rule partitions, named once
// so a test reads as the class it exercises rather than as a string.
const (
	// carrierDeclaredScalar is DECLARED, kind `enum`, single-valued. Its
	// declaration is what makes "is not set-valued" truthful (REQ-11).
	carrierDeclaredScalar = "tier"
	// carrierDeclaredSet is DECLARED, kind `set`. It is the hoisted
	// empty-value arm's guard witness (REQ-18, fixture G).
	carrierDeclaredSet = "labels"
	// carrierUndeclaredScalarKey and carrierUndeclaredArrayKey are ABSENT
	// from the model's tag table. C1 makes both pure carriers.
	carrierUndeclaredScalarKey = "extra"
	carrierUndeclaredArrayKey  = "extras"
)

// The verbatim wire values the MVV names. Written as Go literals of the
// exact bytes the caller spells, so an assertion of "byte-for-byte as
// given" compares against the argv token and not against a re-rendering.
const (
	// carrierScalarValue is MVV 2's `--tag extra=plain`.
	carrierScalarValue = "plain"
	// carrierArrayValue is MVV 2's `--tag extras=["a","b"]`. It is
	// deliberately NOT canonical-set-shaped in the sense that matters: a
	// canonicaliser would be free to re-render it, and C1 forbids that.
	carrierArrayValue = `["a","b"]`
	// carrierDeclaredSetValue is the DECLARED set's value. It is
	// re-canonicalised on the wire (REQ-5) and so is asserted separately.
	carrierDeclaredSetValue = `["security"]`
)

// carrierTableModel is the RDR 0020 fixture model: MVV step 1 —
// "a fixture model declaring one scalar tag and one set tag, with no
// declaration for `extra` or `extras`."
//
// Structure, keyed to the clause each element discriminates:
//
//   - `[tags.tier]`   declared enum scalar, GUARDED by both rows, so a
//     resolve actually selects and MVV 3's with/without comparison has a
//     selected rule to compare.
//   - `[tags.labels]` declared set, UNGUARDED, so it proves a declared set
//     parses an array literal without contributing a dimension.
//   - no `[tags.extra]`, no `[tags.extras]` — C1's identity rule ("ABSENT
//     from the loaded model's normalized tag table") is established by
//     this absence and by nothing else.
const carrierTableModel = `outcomes = ["decide"]

[model]
id = "carrier-0020"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.tier]
provenance = "observed"
kind = "enum"
domain = ["free", "paid"]
single_valued = true
required = true

[tags.labels]
provenance = "observed"
kind = "set"
elements = ["security", "docs"]
required = true

[emit.plan]
kind = "enum"
domain = ["basic", "pro"]

[[rule]]
id = "free"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.tier]
eq = "free"
[rule.emit]
plan = "basic"

[[rule]]
id = "paid"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.tier]
eq = "paid"
[rule.emit]
plan = "pro"
`

// The owned/write surface REUSES RDR 0005's shipped MVV fixture rather
// than authoring a second accessor-bearing model. Two reasons, both
// load-bearing: `flow-tag-owned` and `flow-write-unbound` are 0005's
// refusals and C1 leaves them unchanged, so asserting them against 0005's
// own fixture is the strongest form of "unchanged"; and a fresh
// accessor model would put this suite's fixture — not the seam — on the
// hook for every loader rule about read/write pairing.
//
// From `flowMVVModel`: `status` is OWNED (the `flow-tag-owned` witness),
// `profile` is an OBSERVED enum, `labels` is an OWNED set, and no tag is
// named `extra` or `extras`.
// carrierOwnedKey is DECLARED with provenance `owned` in 0005's MVV
// fixture, and has a bound writer. A `--tag` naming it refuses
// `flow-tag-owned`; a `--write` naming it reaches the conformance arms
// whose copy `canonicalValue` retains (REQ-28).
const carrierOwnedKey = "status"

// carrierUnknownTagAtomModel names an UNDECLARED tag in a rule's guard
// atom. C1 leaves the model-side closed world untouched (REQ-13), so the
// LOADER must still refuse this — the carrier arm widens the CALLER's
// vocabulary only. Reaching the CLI as `flow-model-invalid` is the
// observable form of `unknown_tag`.
const carrierUnknownTagAtomModel = `outcomes = ["decide"]

[model]
id = "carrier-unknown-atom-0020"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[[rule]]
id = "free"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.extra]
eq = "plain"
`

// --- invocation builders -------------------------------------------------

// carrierResolve builds the `flow resolve` argv for the carrier fixture,
// with the declared pair always supplied (both are `required`) and the
// caller's extra `--tag` flags appended. `--as=json` is always last, so
// every oracle in this suite reads the machine envelope.
func carrierResolve(model string, tags ...string) []string {
	args := []string{
		"flow", "resolve", "--model", model, "--outcome", "decide",
	}
	for _, tag := range tags {
		args = append(args, "--tag", tag)
	}
	return append(args, "--as=json")
}

// declaredPair is the two DECLARED tags the fixture model requires, spelled
// as `--tag` values. Every carrier invocation carries them, so a carrier
// key is admitted BESIDE a real declaration and not in isolation.
func declaredPair() []string {
	return []string{
		carrierDeclaredScalar + "=free",
		carrierDeclaredSet + "=" + carrierDeclaredSetValue,
	}
}

// withCarriers returns the declared pair plus the given carrier tags.
func withCarriers(tags ...string) []string {
	return append(declaredPair(), tags...)
}

// --- oracles -------------------------------------------------------------

// observedEcho decodes the resolve payload's `observed` field as the
// `map[string]string` the wire carries (REQ-51: `resolve.Input.Observed` is
// a `map[string]string`, so every value is a JSON STRING — an array literal
// rides as the raw argument text, not as a JSON array).
//
// Decoding into `map[string]string` rather than `map[string]any` is
// load-bearing: it FAILS if an implementation ever promoted a carried array
// literal to a real JSON array, which would be a re-encoding C1 forbids.
func observedEcho(t *testing.T, stdout string) map[string]string {
	t.Helper()

	var env struct {
		Type string `json:"type"`
		Data struct {
			Observed map[string]string `json:"observed"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &env); err != nil {
		t.Fatalf("stdout is not one JSON object, or `observed` is not a "+
			"map of STRING values (C1 forbids re-encoding a carried array "+
			"literal into a JSON array): %v\n%s", err, stdout)
	}
	if env.Type != "ok" {
		t.Fatalf("envelope type = %q; want %q\n%s", env.Type, "ok", stdout)
	}
	return env.Data.Observed
}

// carrierSelection is the projected SELECTION of one resolve payload —
// every field `0020:S2`'s fixture F2 names — with `observed` deliberately
// EXCLUDED. It is what "identical selected rule and outcome … the sole
// payload delta is `observed`" (REQ-52) is asserted over.
type carrierSelection struct {
	Outcome      string            `json:"outcome"`
	Rule         string            `json:"rule"`
	Emit         map[string]string `json:"emit"`
	Dispositions map[string]any    `json:"dispositions"`
	Gates        []any             `json:"gates"`
	Next         map[string]any    `json:"next"`
	Writes       map[string]any    `json:"writes"`
	Clear        []any             `json:"clear"`
	Escaped      bool              `json:"escaped"`
}

// selectionOf decodes the projected selection off a success envelope.
func selectionOf(t *testing.T, stdout string) carrierSelection {
	t.Helper()

	var env struct {
		Type string           `json:"type"`
		Data carrierSelection `json:"data"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &env); err != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", err, stdout)
	}
	if env.Type != "ok" {
		t.Fatalf("envelope type = %q; want %q\n%s", env.Type, "ok", stdout)
	}
	return env.Data
}

// selectionDigest renders a decoded selection as canonical JSON, so two
// runs can be compared as ONE string. Go's `encoding/json` sorts map keys,
// which makes the rendering order-stable; this is a comparison device for
// the test's own diffing and NOT a claim about output determinism, which
// `0020:G-cross-cutting` (REQ-59) explicitly disclaims.
func selectionDigest(t *testing.T, sel carrierSelection) string {
	t.Helper()

	body, err := json.Marshal(sel)
	if err != nil {
		t.Fatalf("re-render the projected selection: %v", err)
	}
	return string(body)
}

// refusalMessage decodes the `message` off a `--as=json` failure line.
// `0020:G-cross-cutting` scopes C1's "byte-identical" claims to MESSAGE
// STRINGS between two refusal sites (REQ-59) — never to a payload hash — so
// message equality, and never a digest, is what this suite compares.
func refusalMessage(t *testing.T, stdout string) string {
	t.Helper()

	var env struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Param   string `json:"param"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &env); err != nil {
		t.Fatalf("failure stdout is not one JSON object: %v\n%s", err, stdout)
	}
	return env.Message
}

// runCarrier drives an invocation and returns stdout plus the error, with
// no assertion of its own. Tests that must inspect BOTH legs of a
// pass/refuse pair use it directly.
func runCarrier(t *testing.T, args ...string) (string, error) {
	t.Helper()

	stdout, _, err := runCmd(t, args...)
	return stdout, err
}
