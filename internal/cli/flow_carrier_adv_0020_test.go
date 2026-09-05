package cli

// RDR 0020 — ADVERSARIAL review of the carrier widening.
//
// This file is written against the record's TRADE-OFFS / FAILURE MODES
// section rather than against its Testing Strategy, and it takes the
// reviewer's posture: `0020:C1` WIDENS an admission seam, so the question is
// what a caller can now push through that admission previously refused, and
// whether the compensating controls the Failure Modes section NAMES actually
// exist on the paths where they are needed.
//
// The section makes exactly three load-bearing promises, and each test below
// attacks one:
//
//   - "Visible: every refusal that survives is unchanged and still fires at
//     exit 2 BEFORE ANY ACCESSOR" — ADV-2.
//   - "Silent: a misspelled key (declared or not) passes as a carrier and the
//     intended rule fails to match; resolution then refuses no-match or
//     routes to an escape row" — ADV-3.
//   - "Diagnosis: the resolve payload's `observed` field echoes every carried
//     key byte-for-byte — the stray spelling sits beside the declared keys IN
//     THE SAME ENVELOPE THE REFUSAL RIDES" — ADV-3, and
//     `0020:G-cross-cutting`'s "byte-preserved … no re-encoding" — ADV-1.
//
// The fixture corpus (`flow_carrier_fixtures_0020_test.go`) is REUSED rather
// than duplicated: an adversarial test that authored its own model would be
// attacking a different seam than the one the suite pins.

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// --- ADV-1 ---------------------------------------------------------------

// ADV-1 — `0020:G-cross-cutting` (REQ-6) / `0020:C1` (REQ-2):
// "A carried value is admitted VERBATIM — byte-preserved, with no
// canonicalisation, folding, normalisation, or re-encoding — and echoes in
// the resolve payload's `observed` field as given."
//
// ADVERSARIAL PREMISE. The widening removes every DOMAIN from a carried
// value: a declared `enum` or `set` holds its value to `domain`/`elements`,
// and a carrier is held to nothing. So the widening is precisely what makes
// ARBITRARY BYTES reachable at this seam from ordinary argv, where before
// only domain-conforming text got through for any key an author had thought
// about. "Byte-preserved" is therefore a much stronger claim after C1 than
// before it, and the existing REQ-6 table only exercises printable ASCII
// (sorting, spacing, case, HTML metacharacters). The bytes that actually
// break a JSON echo — C0 control characters, NUL, and INVALID UTF-8 — are
// untested.
//
// The C0/NUL legs are expected to hold (Go's encoder escapes them and the
// decoder restores them). The INVALID UTF-8 leg is the one that discriminates:
// `encoding/json` replaces every ill-formed byte with U+FFFD on the way out,
// which is a RE-ENCODING in C1's own vocabulary, and it is lossy — two
// distinct carried values collapse to the same echo.
//
// SCOPE, stated so this is not read as a regression claim: the mangling is
// PRE-EXISTING and provenance-blind — a DECLARED bare `scalar` (no domain)
// echoes `"a\xffb"` as `"a�b"` at HEAD too, which the second subtest
// pins. C1 widens nothing here; it only makes the pre-existing limit
// reachable for a caller's whole undeclared vocabulary. The finding is a
// bound on the record's REQ-6 wording, not a defect in the seam C1 moved.
func TestAdv1_0020_CarriedControlBytesAndInvalidUTF8AtTheVerbatimBoundary(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)

	cases := []struct {
		name  string
		value string
		// preserved records whether the echo can return the value at all.
		// A false here is the FINDING, not the expectation.
		preserved bool
		what      string
	}{
		{
			name:      "c0-control-newline",
			value:     "a\nb",
			preserved: true,
			what:      "a C0 control the encoder escapes and the decoder restores",
		},
		{
			name:      "c0-control-tab",
			value:     "a\tb",
			preserved: true,
			what:      "a C0 control the encoder escapes and the decoder restores",
		},
		{
			name:      "nul-byte",
			value:     "a\x00b",
			preserved: true,
			what:      "NUL, which rides as \\u0000",
		},
		{
			name:      "del-and-c1-controls",
			value:     "a\x7fb",
			preserved: true,
			what:      "DEL and a C1 control",
		},
		{
			name:      "combining-mark-is-not-nfc-folded",
			value:     "é",
			preserved: true,
			what: "an NFD sequence a Unicode normaliser would fold to " +
				"U+00E9 — C1 forbids normalisation",
		},
		{
			name:      "lone-surrogate-escape-text",
			value:     `\ud800`,
			preserved: true,
			what:      "the literal backslash-u text, which is not decoded",
		},
		{
			name:      "invalid-utf8-lone-continuation",
			value:     "a\xffb",
			preserved: false,
			what: "an ill-formed UTF-8 byte, which `encoding/json` REPLACES " +
				"with U+FFFD — a lossy re-encoding of the carried bytes",
		},
		{
			name:      "invalid-utf8-truncated-sequence",
			value:     "a\xe2\x82b",
			preserved: false,
			what: "a truncated UTF-8 sequence, likewise replaced with U+FFFD " +
				"— two distinct carried values collapse to the same echo",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, err := runCarrier(t, carrierResolve(model, withCarriers(
				carrierUndeclaredScalarKey+"="+tc.value)...)...)
			if err != nil {
				t.Fatalf("the carrier was REFUSED; C1's refusal list is "+
					"closed and no arm inspects the value's ENCODING: %v\n%s",
					err, stdout)
			}

			got := observedEcho(t, stdout)[carrierUndeclaredScalarKey]
			switch {
			case tc.preserved && got != tc.value:
				t.Errorf("observed[%q] = %q; want %q byte-for-byte — the "+
					"echo did not preserve %s",
					carrierUndeclaredScalarKey, got, tc.value, tc.what)
			case !tc.preserved && got == tc.value:
				t.Errorf("observed[%q] round-tripped %q intact; this test "+
					"documents that it does NOT (%s). The seam improved — "+
					"retire this expectation rather than weakening it",
					carrierUndeclaredScalarKey, tc.value, tc.what)
			case !tc.preserved:
				t.Logf("ADV-1 documented limit: observed[%q] = %q, not the "+
					"carried %q — %s",
					carrierUndeclaredScalarKey, got, tc.value, tc.what)
			}
		})
	}
}

// ADV-1b — the scope leg. It proves the ADV-1 mangling is PROVENANCE-BLIND
// and therefore predates C1: a DECLARED bare `scalar` carries no domain, so
// the same ill-formed bytes reach the same echo through the DECLARED arm.
//
// This is what makes ADV-1 a bound on REQ-6's wording and NOT a regression
// the widening introduced. Without this leg a reader could take ADV-1 for a
// defect in the carrier branch, and "fix" the carrier branch alone — which
// would put the two arms out of agreement and break `0020:G-cross-cutting`'s
// actual invariant ("C1 widens what admission accepts and narrows nothing").
func TestAdv1b_0020_TheEchosUTF8LimitIsProvenanceBlindAndPredatesTheCarrier(t *testing.T) {
	const bareScalarModel = `outcomes = ["decide"]

[model]
id = "carrier-adv-scalar-0020"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.freeform]
provenance = "observed"
kind = "scalar"

[[rule]]
id = "only"
[rule.match.recognized]
eq = "decide"
`

	model := writeFlowModel(t, bareScalarModel)
	const illFormed = "a\xffb"

	stdout := requireSuccess(t, "flow", "resolve", "--model", model,
		"--outcome", "decide", "--tag", "freeform="+illFormed, "--as=json")

	declared := observedEcho(t, stdout)["freeform"]
	if declared == illFormed {
		t.Fatalf("the DECLARED bare scalar preserved %q; ADV-1's premise "+
			"was that the echo cannot, so ADV-1 would then be a carrier-"+
			"specific regression rather than a shared limit", illFormed)
	}

	// The two arms must agree: the carried value and the declared value are
	// mangled the SAME way. Disagreement is the real defect.
	carrierModel := writeFlowModel(t, carrierTableModel)
	carried := observedEcho(t, requireSuccess(t, carrierResolve(carrierModel,
		withCarriers(carrierUndeclaredScalarKey+"="+illFormed)...)...))[carrierUndeclaredScalarKey]

	if carried != declared {
		t.Errorf("the echo mangles a CARRIED value as %q but a DECLARED "+
			"one as %q; `0020:G-cross-cutting` says C1 widens and narrows "+
			"nothing, so the two arms must agree byte-for-byte",
			carried, declared)
	}
	t.Logf("ADV-1b: both arms echo %q for the carried bytes %q — the limit "+
		"is shared, not introduced by C1", carried, illFormed)
}

// --- ADV-2 ---------------------------------------------------------------

// ADV-2 — Failure Modes, "Visible" clause:
// "every refusal that survives is unchanged and still fires at exit 2
// BEFORE ANY ACCESSOR — reserved, owned, duplicate and grammar for ANY key,
// and empty-value for any key that is not a declared set".
//
// ADVERSARIAL PREMISE. "Before any accessor" is the clause with a SIDE
// EFFECT behind it, and the widening moves the boundary it describes. An
// invocation carrying an undeclared ARRAY used to exit 2 at admission with
// no accessor run at all; under C1 the same argv is admitted and the verb
// proceeds to run readers, gates and WRITERS. So the seam that keeps a
// refusing invocation side-effect-free is now doing strictly more work, and
// the "before any accessor" promise is worth attacking directly rather than
// inferring from the exit code.
//
// The oracle is a WRITE-BEARING verb over a seeded artifact, with the
// artifact's bytes read before and after. An exit-2 admission refusal that
// had already run the write accessor would leave the artifact changed while
// telling the caller its input was rejected — the worst shape this widening
// could take, and one no exit-code assertion can see.
func TestAdv2_0020_AnAdmissionRefusalLeavesNoAccessorSideEffect(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)

	// Each case is a refusal that C1 keeps in the closed list, paired with
	// the write it would have performed had admission let it through.
	cases := []struct {
		name string
		code string
		tag  string
		why  string
	}{
		{
			name: "undeclared-key-empty-value",
			code: "flow-tag-invalid",
			tag:  "stray=",
			why:  "the hoisted empty-value arm (REQ-10, REQ-15)",
		},
		{
			name: "reserved-key",
			code: "flow-tag-reserved",
			tag:  "recognized=advance",
			why:  "the reserved guard (REQ-9)",
		},
		{
			name: "owned-key",
			code: "flow-tag-owned",
			tag:  carrierOwnedKey + "=final",
			why:  "the owned guard (REQ-9, REQ-29)",
		},
		{
			name: "grammar",
			code: "flow-tag-invalid",
			tag:  "no-equals-sign",
			why:  "the grammar arm (REQ-8)",
		},
		{
			name: "duplicate-undeclared-key",
			code: "flow-tag-duplicate",
			tag:  "stray=plain",
			why:  "the duplicate guard, which C1 keeps ahead of the carrier branch (REQ-23)",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			artifact := seedArtifact(t, model, "status=draft")
			bind := artifactBinding(flowStateRole, artifact)

			before := readArtifactAdv0020(t, artifact)

			args := []string{"flow", "set-state", "--model", model,
				"--artifact", bind, "--tag", tc.tag}
			if tc.name == "duplicate-undeclared-key" {
				args = append(args, "--tag", tc.tag)
			}
			args = append(args, "--write", "status=final", "--as=json")

			requireRefusal(t, tc.code, 2, args...)

			if after := readArtifactAdv0020(t, artifact); after != before {
				t.Errorf("%s refused at admission but the artifact CHANGED "+
					"— the write accessor ran behind an exit-2 refusal, so "+
					"%s no longer precedes every accessor:\nbefore = %q\n"+
					"after  = %q", tc.name, tc.why, before, after)
			}
		})
	}
}

// ADV-2b — the widening's own side-effect leg, stated as the reviewer's
// actual worry rather than as a refusal pin.
//
// An undeclared ARRAY refused at HEAD and is ADMITTED under C1 (REQ-31), so
// an argv that previously guaranteed "no accessor ran" now guarantees the
// opposite. That is a real behavioural widening with a side-effect surface,
// and the record accepts it — but only because the carried key cannot reach
// owned state (Decision Rationale (b): the loader's `unknown_tag` refusals
// make a carrier unreadable by any rule, accessor, write, clear or initial).
//
// This pins the acceptance: the previously-refusing argv now SUCCEEDS and
// performs its declared write, and the carried key appears in NEITHER the
// persisted state nor the write set.
func TestAdv2b_0020_TheAdmittedArrayRunsAccessorsWithoutReachingOwnedState(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	artifact := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, artifact)

	// At HEAD this argv exited 2 at admission on the array literal and no
	// accessor ran. Under C1 it is admitted.
	stdout := requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", bind, "--tag", `stray=["a","b"]`,
		"--write", "status=final", "--as=json")
	t.Logf("previously-refusing argv now succeeds: %s", strings.TrimSpace(stdout))

	// The oracle for owned state goes through the CLI, per the suite's
	// standing rule that no test asserts on the artifact's on-disk format.
	state := requireSuccess(t, "flow", "read-state", "--model", model,
		"--artifact", bind, "--artifact",
		artifactBinding(flowOrphanRole, artifact), "--as=json")

	if !strings.Contains(state, `"status":"final"`) {
		t.Errorf("the declared write did not reach owned state; the "+
			"admitted carrier must not suppress it: %s",
			strings.TrimSpace(state))
	}
	if strings.Contains(state, "stray") {
		t.Errorf("the CARRIED key reached owned state: %s — Decision "+
			"Rationale (b) rests on a carrier being unreadable and "+
			"unwritable by every model-side path", strings.TrimSpace(state))
	}
}

// readArtifactAdv0020 reads the caller-owned artifact's raw bytes.
//
// It is a pure CHANGE DETECTOR across an invocation that must not have
// written anything, and it makes NO claim about the on-disk format: the
// bytes are only ever compared to themselves. Reading the file directly
// rather than round-tripping through `flow read-state` is deliberate — a
// CLI round-trip would embed per-invocation temp paths in the string and
// compare unequal for reasons that have nothing to do with the artifact.
func readArtifactAdv0020(t *testing.T, path string) string {
	t.Helper()

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read the seeded artifact: %v", err)
	}
	return string(raw)
}

// --- ADV-3 ---------------------------------------------------------------

// ADV-3 — Failure Modes, the "Silent" and "Diagnosis" clauses, which are
// one argument and must be attacked as one:
//
//	"Silent: a misspelled key (declared or not) passes as a carrier and the
//	intended rule fails to match; resolution then refuses no-match or routes
//	to an escape row … Diagnosis: the resolve payload's `observed` field
//	echoes every carried key byte-for-byte — the stray spelling sits beside
//	the declared keys IN THE SAME ENVELOPE THE REFUSAL RIDES."
//
// ADVERSARIAL PREMISE, and the reason this is the sharpest of the three.
// The record accepts the silent-typo cost EXPLICITLY, and the Premortem
// says the approach "answers it" with the echo — "which is where the pinned
// MVV test points a diagnostician". But the Silent clause and the Diagnosis
// clause describe DIFFERENT EXITS. The Silent clause's outcome is a
// REFUSAL (no-match) or an escape route; the Diagnosis clause's `observed`
// echo is a field of the SUCCESS payload. The claim that the stray spelling
// rides "the same envelope the refusal rides" is therefore a claim about
// the REFUSAL envelope, and it is the one the operator depends on at
// exactly the moment they are stuck.
//
// It does not hold. `clierr.CLIError` — the whole refusal envelope — is
// `{code, message, param, detail, hint, findings}`. There is no `observed`
// field on it and no path that could populate one, so when the misspelling
// produces the failure the Failure Modes section itself predicts, the
// operator receives an envelope naming the DECLARED key as absent and
// carrying NO trace of the stray key they actually typed.
//
// The finding is documentary: the compensating control the record leans on to
// accept the silent-typo cost is absent from the refusing path. It is not a
// defect in `parseTags` — the admission seam does exactly what C1 says, and no
// REQ in req-list.md obligates the refusal envelope to echo — so the remedy is
// a record-level one (an amended Failure Modes claim, or a follow-on that puts
// the echo on the refusal envelope), not a change to the carrier branch.
//
// The test therefore PINS THE SHIPPED SHAPE rather than asserting the record's
// claim: it fails the moment the refusal envelope gains an `observed` field or
// otherwise names the stray key, which is exactly when the record's Diagnosis
// clause becomes true and this pin must be revisited. See deviations.md D3 and
// verification.md ADV-3.
func TestAdv3_0020_TheStraySpellingDoesNotRideTheRefusalEnvelope(t *testing.T) {
	model := writeFlowModel(t, carrierTableModel)

	// `teir` is a misspelling of the declared, guarded `tier`. Under C1 it
	// is admitted as a carrier; the guard on the real `tier` is then
	// undecidable and the resolve REFUSES — the Silent clause's own
	// narrative, reproduced exactly.
	const misspelled = "teir"

	stdout, err := runCarrier(t, carrierResolve(model,
		misspelled+"=free",
		carrierDeclaredSet+"="+carrierDeclaredSetValue)...)
	if err == nil {
		t.Fatalf("the misspelled key did not lead to a refusal; this test "+
			"reproduces the Failure Modes narrative and needs the refusing "+
			"exit to assert against:\n%s", stdout)
	}

	// The refusal envelope is the artifact the operator actually holds.
	var envelope map[string]json.RawMessage
	if uerr := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &envelope); uerr != nil {
		t.Fatalf("the refusal envelope is not one JSON object: %v\n%s",
			uerr, stdout)
	}

	// PIN, not an assertion of the record's claim. Today the refusal
	// envelope carries no `observed` field and never names the stray key.
	// If either changes, the Diagnosis clause has become true on this exit
	// and D3 is resolved — update this pin then, deliberately.
	if _, ok := envelope["observed"]; ok {
		t.Errorf("the refusal envelope now carries an `observed` field: the "+
			"Failure Modes Diagnosis clause has become true on the refusing "+
			"exit. Resolve deviations.md D3 and update this pin. Keys: %v",
			envelopeKeysAdv0020(envelope))
	}

	if strings.Contains(stdout, misspelled) {
		t.Errorf("the refusal envelope now mentions the misspelled key %q, "+
			"so the operator has a route from the refusal to the typo. "+
			"Resolve deviations.md D3 and update this pin:\n%s",
			misspelled, strings.TrimSpace(stdout))
	}
}

// ADV-3b — the same clause on its OTHER exit, and the reason ADV-3 is not
// merely pedantic about envelope shape.
//
// The Silent clause offers two outcomes: "refuses no-match OR ROUTES TO AN
// ESCAPE ROW". The escape route is the more dangerous of the two, because
// it exits ZERO: the misspelling produces a SUCCESSFUL invocation that took
// a different edge than the caller intended, with no refusal to read at all.
//
// On that path the `observed` echo IS present — the Diagnosis clause holds
// — so this test PASSES and stands as a regression pin on the one path
// where the record's compensating control genuinely works. Keeping it
// beside ADV-3 is what makes the finding legible: the echo is real, it is
// just on the wrong exit.
func TestAdv3b_0020_OnTheEscapeExitTheEchoIsPresentAndCarriesTheTypo(t *testing.T) {
	const escapeModel = `outcomes = ["decide"]

[model]
id = "carrier-adv-escape-0020"
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

[emit.plan]
kind = "enum"
domain = ["basic"]

[[rule]]
id = "free"
[rule.match.recognized]
eq = "decide"
[rule.match.tier]
eq = "free"
[rule.emit]
plan = "basic"

[[rule]]
id = "bail"
escape = ["no_match"]
[rule.match.recognized]
eq = "decide"
`

	model := writeFlowModel(t, escapeModel)
	const misspelled = "teir"

	stdout, err := runCarrier(t, "flow", "resolve", "--model", model,
		"--outcome", "decide", "--tag", misspelled+"=free", "--as=json")
	if err != nil {
		t.Fatalf("the misspelled key did not reach the escape row; this "+
			"test owns the exit-0 leg of the Silent clause and needs it: "+
			"%v\n%s", err, stdout)
	}

	// The escape row was taken: the caller got a SUCCESS on an edge they
	// did not intend, with no refusal to read.
	if sel := selectionOf(t, stdout); sel.Rule != "bail" || !sel.Escaped {
		t.Fatalf("selected rule = %q, escaped = %v; want the escape row "+
			"`bail` — the Silent clause's exit-0 route", sel.Rule, sel.Escaped)
	}

	echo := observedEcho(t, stdout)
	if echo[misspelled] != "free" {
		t.Errorf("observed[%q] = %q; want %q — on the exit-0 escape route "+
			"the echo is the ONLY diagnosis the operator gets, so it must "+
			"carry the stray spelling verbatim",
			misspelled, echo[misspelled], "free")
	}
	if _, ok := echo[carrierDeclaredScalar]; ok {
		t.Errorf("observed carries the DECLARED key %q, which was never "+
			"supplied; the echo must show what the caller actually typed",
			carrierDeclaredScalar)
	}
}

// envelopeKeysAdv0020 lists an envelope's top-level field names, so a
// failure message can show what the refusal DOES carry.
func envelopeKeysAdv0020(envelope map[string]json.RawMessage) []string {
	out := make([]string, 0, len(envelope))
	for key := range envelope {
		out = append(out, key)
	}
	return out
}
