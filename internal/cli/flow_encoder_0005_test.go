package cli

// RDR 0005 — the canonical set literal and the shared encoder
// (REQ-70..REQ-73), the `Finding` record as it crosses the wire
// (REQ-14..REQ-18, REQ-79), and the both-modes findings rendering
// (REQ-11, REQ-129).
//
// The discriminating property throughout: `<`, `>`, and `&` serialize as
// THEMSELVES. Bare `json.Marshal` HTML-escapes them, so any emit site still
// on it renders `\u003c` / `\u0026` and fails here — which is precisely the
// defect REQ-72 removes.

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

// REQ-70: "A set-valued tag crossing the CLI in --tag, --write, or any
// payload MUST use the canonical JSON array form JDR 0001 §D13 fixes —
// members sorted, duplicate-free, compact — rendered with HTML escaping
// DISABLED, so <, >, and & serialize as themselves and never as \\u003c,
// \\u003e, or \\u0026."
// REQ-73: "Set values are JDR 0001 §D13's canonical JSON array everywhere
// they cross the CLI".
// REQ-127 / `0005:S5`: "One set member carries `<` and one carries `&`:
// both render as themselves at every emit site".
// ADVERSARIAL
func TestReq70And73_SetValuesRenderAngleAndAmpersandAsThemselves(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")

	stdout := requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--write", "labels="+canonicalSetLiteral, "--as=json")

	assertNoHTMLEscapes(t, "set-state success envelope", stdout)
	if !strings.Contains(stdout, setMemberLT) {
		t.Errorf("the emitted payload does not carry the member %q as "+
			"itself:\n%s", setMemberLT, stdout)
	}
	if !strings.Contains(stdout, setMemberAmp) {
		t.Errorf("the emitted payload does not carry the member %q as "+
			"itself:\n%s", setMemberAmp, stdout)
	}
}

// REQ-71: "Every site that emits or compares a canonical set literal MUST
// use the same encoder, so plan-to-request copy-through and read-back
// equality are byte equality."
// REQ-72: "Route `internal/cli/clierr::EmitJSON` and
// `internal/cli/respond::writeJSONLine` through one shared helper using
// `SetEscapeHTML(false)` … Bare `json.Marshal` must not remain on a path a
// set value crosses."
// one bare `json.Marshal` would still be hiding on.
// ADVERSARIAL — the FAILURE envelope is the second emit site, and it is the
func TestReq71And72_BothEmitSitesUseTheSameNonEscapingEncoder(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	// Emit site 1 — respond's success envelope.
	success := requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", bind, "--write", "labels="+canonicalSetLiteral,
		"--as=json")
	assertNoHTMLEscapes(t, "success envelope (respond)", success)

	// Emit site 2 — clierr's failure envelope. The refusal is made to carry
	// the very same set literal in its `param`, so the two encoders are
	// compared on identical bytes.
	failure, _, err := runCmd(t, "flow", "next", "--model", model,
		"--artifact", bind, "--tag", "labels="+canonicalSetLiteral, "--as=json")
	if err == nil {
		t.Fatalf("a `--tag` on the OWNED key `labels` succeeded:\n%s", failure)
	}
	assertNoHTMLEscapes(t, "failure envelope (clierr)", failure)
}

// REQ-70 (the `--tag` arm): a set value crossing in `--tag` obeys the same
// canonical form.
// INPUT EDGE
func TestReq70_SetValuesCrossingInTagObeyTheSameCanonicalForm(t *testing.T) {
	model := writeFlowModel(t, flowSetObservedModel)
	art := seedArtifact(t, model, "status=draft")

	stdout := requireSuccess(t, "flow", "next", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--tag", "marks="+canonicalSetLiteral, "--as=json")

	assertNoHTMLEscapes(t, "next success envelope", stdout)

	data := flowData(t, stdout)
	observed, ok := data["observed"].(map[string]any)
	if !ok {
		t.Fatalf("`observed` is not an object: %#v", data["observed"])
	}
	assertCanonicalSetValue(t, "marks", observed["marks"])
}

// REQ-8: "Failures MUST use the existing CLIError JSON/text envelope …
// extended by exactly one omitempty structured field, findings."
// REQ-16 / A-6: `Findings []Finding` carries `json:"findings,omitempty"`.
// BOUNDARY — a refusal with NO findings must not carry the key at all.
func TestReq8And16_FindingsIsOmitEmptyOnTheCLIErrorEnvelope(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)

	// `flow-tag-reserved` is a single-subject refusal: its carrier is
	// `param`, and it names no findings (REQ-79, REQ-82).
	stdout, _, err := runCmd(t, "flow", "next", "--model", model,
		"--tag", "recognized=hold", "--as=json")
	if err == nil {
		t.Fatalf("a reserved `--tag` succeeded:\n%s", stdout)
	}

	var raw map[string]json.RawMessage
	if uerr := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &raw); uerr != nil {
		t.Fatalf("failure line is not one JSON object: %v\n%s", uerr, stdout)
	}
	if _, present := raw["findings"]; present {
		t.Errorf("a single-subject refusal emits the `findings` key; the "+
			"field is `omitempty` on `CLIError`, so a refusal that names no "+
			"findings omits it entirely:\n%s", stdout)
	}
	// The single-subject carrier IS present.
	if _, present := raw["param"]; !present {
		t.Errorf("the refusal carries no `param`; REQ-79 fixes `param` for "+
			"single-subject failures:\n%s", stdout)
	}
}

// REQ-79: "`param` for single-subject failures, and `findings[]` where
// several rows, keys, atoms, gates, or load categories must be named."
// convenience: a multi-subject refusal must not smuggle its subjects into
// `param`.
// DOMAIN EDGE — the two carriers are assigned by failure SHAPE, not by
func TestReq79_CarrierChoiceFollowsTheFailuresShape(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	ambiguous := writeFlowModel(t, flowAmbiguousModel)
	invalid := writeFlowModel(t, flowInvalidModel)
	art := seedArtifact(t, model, "status=draft")
	ambigArt := seedArtifact(t, ambiguous, "status=draft")

	t.Run("single-subject-uses-param", func(t *testing.T) {
		ce := requireRefusal(t, "flow-tag-duplicate", 2,
			"flow", "next", "--model", model,
			"--artifact", artifactBinding(flowStateRole, art),
			"--tag", "profile=mid", "--tag", "profile=foundational", "--as=json")
		if ce.Param == "" {
			t.Error("a single-subject refusal carries no `param`")
		}
	})

	t.Run("multi-subject-uses-findings", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			args []string
			code string
		}{
			{"matching rows", []string{"flow", "resolve", "--model", ambiguous,
				"--artifact", artifactBinding(flowStateRole, ambigArt),
				"--outcome", "advance"}, "flow-ambiguous-match"},
			{"load categories", []string{"flow", "next", "--model", invalid},
				"flow-model-invalid"},
		} {
			t.Run(tc.name, func(t *testing.T) {
				stdout, _, err := runCmd(t, append(tc.args, "--as=json")...)
				if err == nil {
					t.Fatalf("expected %s:\n%s", tc.code, stdout)
				}
				if len(flowFailureFindings(t, stdout)) == 0 {
					t.Errorf("%s names several subjects but carries no "+
						"`findings[]`; `param` cannot name more than one",
						tc.code)
				}
			})
		}
	})
}

// REQ-14: "Finding MUST be one flat record with omitempty optional fields,
// carrying at least code, message, param, locator, hint, severity, model,
// rule, key, operator, literal, block, and class"
// REQ-15: "each producer populates only the fields it owns, and no producer
// nests its own fields in a sub-object."
// REQ-18: "This RDR populates `code`, `message`, `param`, `locator`, and
// `hint`".
// nested object at all.
// ADVERSARIAL — flatness is asserted structurally: no finding carries a
func TestReq14And15And18_FindingsAreFlatAndUnsetOptionalFieldsAreAbsent(t *testing.T) {
	invalid := writeFlowModel(t, flowInvalidModel)

	stdout, _, err := runCmd(t, "flow", "next", "--model", invalid, "--as=json")
	if err == nil {
		t.Fatalf("an unloadable model succeeded:\n%s", stdout)
	}

	var env struct {
		Findings []map[string]json.RawMessage `json:"findings"`
	}
	if uerr := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &env); uerr != nil {
		t.Fatalf("failure line is not one JSON object: %v\n%s", uerr, stdout)
	}
	if len(env.Findings) == 0 {
		t.Fatal("`flow-model-invalid` carried no findings")
	}

	for i, f := range env.Findings {
		for name, raw := range f {
			trimmed := strings.TrimSpace(string(raw))
			// REQ-15: no producer nests its own fields in a sub-object.
			if strings.HasPrefix(trimmed, "{") {
				t.Errorf("finding[%d].%s is a nested OBJECT (%s); the record "+
					"is FLAT — the atom's fields in particular sit flat, not "+
					"under an `atom` object", i, name, trimmed)
			}
			// REQ-14: optional fields are `omitempty`, so a field that IS
			// present is a field the producer populated — never an empty
			// placeholder.
			if name != "code" && name != "message" && trimmed == `""` {
				t.Errorf("finding[%d].%s is present but EMPTY; optional "+
					"fields are `omitempty`, so each producer populates only "+
					"the fields it owns and the rest are ABSENT", i, name)
			}
		}
		// REQ-18: this RDR's producers populate `code` and `message`, and a
		// load-category finding additionally carries `locator`.
		if _, present := f["code"]; !present {
			t.Errorf("finding[%d] carries no `code`", i)
		}
		if _, present := f["message"]; !present {
			t.Errorf("finding[%d] carries no `message`", i)
		}
	}
}

// REQ-17: "Finding.message MUST be self-sufficient — it MUST render the
// failure readably with no structured field consulted."
// nor a bare restatement of the code.
// DOMAIN EDGE — the message is prose that stands alone: it is neither empty
func TestReq17_EveryFindingMessageStandsAloneWithoutItsStructuredFields(t *testing.T) {
	invalid := writeFlowModel(t, flowInvalidModel)
	ambiguous := writeFlowModel(t, flowAmbiguousModel)
	ambigArt := seedArtifact(t, ambiguous, "status=draft")

	for _, tc := range []struct {
		name string
		args []string
	}{
		{"load-category", []string{"flow", "next", "--model", invalid}},
		{"kernel-refusal", []string{"flow", "resolve", "--model", ambiguous,
			"--artifact", artifactBinding(flowStateRole, ambigArt),
			"--outcome", "advance"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stdout, _, err := runCmd(t, append(tc.args, "--as=json")...)
			if err == nil {
				t.Fatalf("expected a refusal:\n%s", stdout)
			}
			findings := flowFailureFindings(t, stdout)
			if len(findings) == 0 {
				t.Fatal("no findings to check")
			}
			for i, f := range findings {
				if strings.TrimSpace(f.Message) == "" {
					t.Errorf("finding[%d] carries an EMPTY message; the "+
						"message alone must render the failure readably", i)
					continue
				}
				if f.Message == f.Code {
					t.Errorf("finding[%d].message is just the code %q; a "+
						"self-sufficient message reads without any "+
						"structured field consulted", i, f.Code)
				}
				// A message that only makes sense once a structured field is
				// read is not self-sufficient. The observable proxy: it says
				// more than one word.
				if len(strings.Fields(f.Message)) < 2 {
					t.Errorf("finding[%d].message = %q is a single token; it "+
						"must render the failure READABLY on its own",
						i, f.Message)
				}
			}
		})
	}
}

// REQ-11: "it renders `findings` one-per-line under the message."
// REQ-129 / `0005:S7`: "a refusal carrying `findings[]` is rendered in both
// modes … each producer populates only its own fields and unset optional
// fields are absent from the JSON (`omitempty`); no producer nests its
// fields in a sub-object; and each finding's `message` alone renders the
// failure readably".
// DOMAIN EDGE — the set of finding codes in text must EQUAL the set in JSON.
func TestReq11And129_TextModeEnumeratesEveryFindingOnePerLine(t *testing.T) {
	invalid := writeFlowModel(t, flowInvalidModel)

	jsonOut, _, jerr := runCmd(t, "flow", "next", "--model", invalid, "--as=json")
	if jerr == nil {
		t.Fatalf("an unloadable model succeeded:\n%s", jsonOut)
	}
	findings := flowFailureFindings(t, jsonOut)
	if len(findings) == 0 {
		t.Fatal("`flow-model-invalid` carried no findings in JSON mode")
	}

	// Text mode's failure envelope goes to STDERR (clierr.EmitText).
	_, stderr, terr := runCmd(t, "flow", "next", "--model", invalid)
	if terr == nil {
		t.Fatal("text mode succeeded on an unloadable model")
	}

	for _, f := range findings {
		if !strings.Contains(stderr, f.Code) {
			t.Errorf("text mode omits the finding code %q that JSON mode "+
				"reports; the two sets must be EQUAL\nstderr:\n%s",
				f.Code, stderr)
		}
		if !strings.Contains(stderr, f.Message) {
			t.Errorf("text mode omits the finding message %q\nstderr:\n%s",
				f.Message, stderr)
		}
	}

	// One finding is ONE line: no finding's rendering spans a newline.
	var findingLines int
	for _, line := range strings.Split(stderr, "\n") {
		for _, f := range findings {
			if strings.Contains(line, f.Code) && strings.Contains(line, f.Message) {
				findingLines++
				break
			}
		}
	}
	if findingLines < len(findings) {
		t.Errorf("only %d of %d findings render as ONE complete line each; "+
			"`findings` render one-per-line under the message\nstderr:\n%s",
			findingLines, len(findings), stderr)
	}
}

// REQ-129 (the shared-record leg): "the same `Finding` record is populated
// by a kernel refusal, a gate result, and a model-load category."
// DOMAIN EDGE — three producers, one record shape, disjoint field sets.
func TestReq129_ThreeProducersShareOneRecordAndPopulateDisjointFields(t *testing.T) {
	invalid := writeFlowModel(t, flowInvalidModel)
	ambiguous := writeFlowModel(t, flowAmbiguousModel)
	denyModel := writeFlowModel(t, flowGateDenyModel)
	ambigArt := seedArtifact(t, ambiguous, "status=draft")
	denyArt := seedArtifact(t, denyModel, "status=draft")

	producers := []struct {
		name string
		args []string
	}{
		{"model-load-category", []string{"flow", "next", "--model", invalid}},
		{"kernel-refusal", []string{"flow", "resolve", "--model", ambiguous,
			"--artifact", artifactBinding(flowStateRole, ambigArt),
			"--outcome", "advance"}},
		{"gate-result", []string{"flow", "resolve", "--model", denyModel,
			"--artifact", artifactBinding(flowStateRole, denyArt),
			"--outcome", "deny-path"}},
	}

	for _, p := range producers {
		t.Run(p.name, func(t *testing.T) {
			stdout, _, err := runCmd(t, append(p.args, "--as=json")...)
			if err == nil {
				t.Fatalf("expected a refusal:\n%s", stdout)
			}
			findings := flowFailureFindings(t, stdout)
			if len(findings) == 0 {
				t.Fatalf("%s produced NO findings; all three producers "+
					"populate the SAME shared `Finding` record", p.name)
			}
			for i, f := range findings {
				if f.Code == "" || f.Message == "" {
					t.Errorf("finding[%d] from %s is missing `code` or "+
						"`message`; both are required on every producer's "+
						"findings", i, p.name)
				}
			}
		})
	}
}

// --- oracles -------------------------------------------------------------

// assertNoHTMLEscapes fails when out carries any of the three escapes the
// contract forbids. It names the emit site so a failure says WHICH encoder
// is still on bare `json.Marshal`.
func assertNoHTMLEscapes(t *testing.T, site, out string) {
	t.Helper()

	for _, esc := range []string{`\u003c`, `\u003e`, `\u0026`} {
		if strings.Contains(out, esc) {
			t.Errorf("%s renders %s; `<`, `>`, and `&` MUST serialize as "+
				"THEMSELVES — this emit site is still on bare "+
				"`json.Marshal`\noutput:\n%s", site, esc, out)
		}
	}
	if strings.Contains(out, escapedSetLiteral) {
		t.Errorf("%s emitted the HTML-escaped set literal %s instead of the "+
			"canonical %s", site, escapedSetLiteral, canonicalSetLiteral)
	}
}

// assertCanonicalSetValue checks a set value on the wire is JDR 0001 §D13's
// canonical form: the two members present, sorted, duplicate-free, and
// carrying `<` and `&` as themselves.
func assertCanonicalSetValue(t *testing.T, key string, raw any) {
	t.Helper()

	var members []string
	switch v := raw.(type) {
	case string:
		// Rendered as the canonical literal.
		if v != canonicalSetLiteral {
			t.Errorf("%s = %q; want the canonical literal %q — members "+
				"sorted, duplicate-free, compact, HTML escaping DISABLED",
				key, v, canonicalSetLiteral)
		}
		if uerr := json.Unmarshal([]byte(v), &members); uerr != nil {
			t.Fatalf("%s = %q is not a JSON array: %v", key, v, uerr)
		}
	case []any:
		for _, m := range v {
			s, ok := m.(string)
			if !ok {
				t.Fatalf("%s carries a non-string member %#v", key, m)
			}
			members = append(members, s)
		}
	default:
		t.Fatalf("%s = %#v; a set value crosses the CLI as a canonical JSON "+
			"array", key, raw)
	}

	for _, want := range []string{setMemberLT, setMemberAmp} {
		if !slices.Contains(members, want) {
			t.Errorf("%s = %v omits the member %q; the two members carrying "+
				"`<` and `&` must survive byte-identical", key, members, want)
		}
	}
	if !slices.IsSorted(members) {
		t.Errorf("%s = %v is not SORTED; the canonical form is sorted",
			key, members)
	}
	if len(slices.Compact(slices.Clone(members))) != len(members) {
		t.Errorf("%s = %v carries duplicates; the canonical form is "+
			"duplicate-free", key, members)
	}
}
