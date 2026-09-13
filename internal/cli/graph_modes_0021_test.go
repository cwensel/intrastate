package cli

// RDR 0021 — mode coexistence (`0021:C5`) and the refusal dispositions
// (`0021:S8`, the `disposition` mini-check table).
//
// C5 defines all four `--as`×`--emit` cells: none refused, none dead. Under
// `--as=text` stdout is the bare document, riding the gateway's existing
// `TextLiner` seam; under `--as=json` stdout is exactly one `ok` envelope
// whose `data` embeds the same document.

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/respond"
)

// REQ-57: "Under `--as=text` (the default) stdout carries the selected
// document verbatim and nothing else — the payload satisfies the gateway's
// `TextLiner` seam, so the DOT stream pipes to `dot` and the JSON document
// diffs raw in CI; advisories stay on stderr."
// REQ-61: "a consumer that parses stdout as an envelope MUST use
// `--as=json`: the text-mode stream is the bare document by contract"
// REQ-64: "Reuse (`TextLiner` + `OK`)"
// HAPPY PATH — A1's normative fixture F1: text-mode stdout is exactly
// `document + "\n"`, stderr empty absent advisories.
func TestReq57And61And64_TextModeStdoutIsTheBareDocument(t *testing.T) {
	requireGraphVerb(t)

	path := writeModel(t, legalModel)
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"json document", []string{"--model", path}},
		{"dot document", []string{"--model", path, "--emit", "dot"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, err := runGraph(t, tc.args...)
			if err != nil {
				t.Fatalf("export refused: %v", err)
			}

			// NOT an envelope: the bare document by contract, exactly as
			// `version`'s TextLine is its identity string (premortem P-10).
			var envelope map[string]any
			if json.Unmarshal([]byte(strings.TrimSpace(stdout)), &envelope) == nil {
				if _, hasType := envelope["type"]; hasType {
					t.Errorf("text-mode stdout carries an envelope `type` "+
						"key; C5 makes the text stream the BARE document — a "+
						"consumer that wants an envelope must use "+
						"`--as=json`\n%s", stdout)
				}
				if _, hasVersion := envelope["schema_version"]; hasVersion {
					t.Errorf("text-mode stdout carries the envelope's "+
						"`schema_version`; the text stream is the document, "+
						"not the envelope\n%s", stdout)
				}
			}

			// F1: exactly `document + "\n"`. One trailing newline, no
			// decoration, nothing after it.
			if !strings.HasSuffix(stdout, "\n") {
				t.Errorf("text-mode stdout does not end in a newline; F1 " +
					"pins it as exactly `document + \"\\n\"`")
			}
			if strings.HasSuffix(stdout, "\n\n") {
				t.Errorf("text-mode stdout ends in TWO newlines; F1 pins " +
					"exactly one — the gateway's single Fprintln")
			}
			// Absent provoked advisories, stderr is empty.
			if strings.TrimSpace(stderr) != "" {
				t.Errorf("text-mode stderr is not empty: %q — advisories "+
					"stay on stderr and this run provokes none", stderr)
			}
		})
	}
}

// REQ-58: "Under `--as=json` stdout carries exactly one terminal `ok`
// envelope whose `data` embeds the same document: the document object for
// `--emit json`, and for `--emit dot` an object carrying the DOT text as
// the single REQUIRED string member `dot` — spelled normatively here, so
// the documented unwrap is exactly `jq -r .data.dot`."
// ASSUMPTION A-5: `data` is an object `{"dot": "<text>"}`, never a bare
// JSON string.
// HAPPY PATH
func TestReq58_JSONModeEmbedsTheDocumentUnderData(t *testing.T) {
	requireGraphVerb(t)

	path := writeModel(t, legalModel)

	t.Run("emit json", func(t *testing.T) {
		stdout, _, err := runGraph(t, "--model", path, "--as=json")
		if err != nil {
			t.Fatalf("export refused: %v", err)
		}
		assertExactlyOneEnvelopeLine(t, stdout)

		var env struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data"`
		}
		if uerr := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &env); uerr != nil {
			t.Fatalf("stdout is not one JSON object: %v\n%s", uerr, stdout)
		}
		if env.Type != "ok" {
			t.Errorf("envelope type = %q; want %q", env.Type, "ok")
		}
		var doc map[string]any
		if uerr := json.Unmarshal(env.Data, &doc); uerr != nil {
			t.Fatalf("`data` is not the document OBJECT: %v", uerr)
		}
		if got, _ := doc["schema"].(string); got != schemaMarker {
			t.Errorf("data.schema = %q; want %q — `data` embeds the same "+
				"document the text mode emits", got, schemaMarker)
		}
	})

	t.Run("emit dot", func(t *testing.T) {
		stdout, _, err := runGraph(t, "--model", path, "--emit", "dot", "--as=json")
		if err != nil {
			t.Fatalf("export refused: %v", err)
		}
		assertExactlyOneEnvelopeLine(t, stdout)

		var env struct {
			Data json.RawMessage `json:"data"`
		}
		if uerr := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &env); uerr != nil {
			t.Fatalf("stdout is not one JSON object: %v\n%s", uerr, stdout)
		}

		// `data` is an OBJECT with one REQUIRED string member `dot` — not a
		// bare JSON string. The documented unwrap is exactly
		// `jq -r .data.dot`, which only works against an object.
		var asString string
		if json.Unmarshal(env.Data, &asString) == nil {
			t.Fatalf("`data` is a bare JSON string under `--emit dot`; C5 " +
				"spells it as an OBJECT carrying the DOT text as the single " +
				"REQUIRED string member `dot`, so the documented unwrap is " +
				"`jq -r .data.dot`")
		}
		var obj map[string]any
		if uerr := json.Unmarshal(env.Data, &obj); uerr != nil {
			t.Fatalf("`data` is not an object: %v", uerr)
		}
		text, ok := obj["dot"].(string)
		if !ok {
			t.Fatalf("`data.dot` is absent or not a string; C5 spells `dot` "+
				"as the single REQUIRED string member: %v", obj)
		}
		if !strings.Contains(text, "digraph") &&
			!strings.Contains(text, "graph") {
			t.Errorf("`data.dot` does not carry a DOT document: %q", text)
		}
	})
}

// REQ-59: "Both modes derive from one export value (0005's two-modes
// agreement)."
// REQ-71: "`intrastate graph --model <fixture> --as=json | jq .data` equals
// step 2's document, value-for-value."
// REQ-78: "Mode agreement — `--as=json | jq .data` against the `--as=text`
// document, for both `--emit` values. Expected: value-for-value equality
// for `json`; for `dot`, the `jq -r .data.dot` unwrap reproduces the DOT
// text byte-for-byte AFTER accounting for the one trailing newline F1 pins
// on text-mode stdout (`document + \"\\n\"`)"
// HAPPY PATH — the trailing-newline accounting is normative, not incidental.
func TestReq59And71And78_TheTwoModesAgreeOnOneExportValue(t *testing.T) {
	requireGraphVerb(t)

	path := writeModel(t, legalModel)

	t.Run("emit json is value-for-value equal", func(t *testing.T) {
		text, _, err := runGraph(t, "--model", path)
		if err != nil {
			t.Fatalf("text mode refused: %v", err)
		}
		enveloped, _, err := runGraph(t, "--model", path, "--as=json")
		if err != nil {
			t.Fatalf("json mode refused: %v", err)
		}

		var env struct {
			Data json.RawMessage `json:"data"`
		}
		if uerr := json.Unmarshal([]byte(strings.TrimSpace(enveloped)), &env); uerr != nil {
			t.Fatalf("stdout is not one JSON object: %v", uerr)
		}

		var fromText, fromData any
		if uerr := json.Unmarshal([]byte(strings.TrimSpace(text)), &fromText); uerr != nil {
			t.Fatalf("the text-mode document is not JSON: %v", uerr)
		}
		if uerr := json.Unmarshal(env.Data, &fromData); uerr != nil {
			t.Fatalf("`data` is not JSON: %v", uerr)
		}
		if !jsonDeepEqual(fromText, fromData) {
			t.Errorf("`jq .data` does not equal the text-mode document "+
				"value-for-value; both modes derive from ONE export value "+
				"(C5)\n--- text ---\n%s\n--- data ---\n%s", text, env.Data)
		}
	})

	t.Run("emit dot reproduces the text byte-for-byte", func(t *testing.T) {
		text, _, err := runGraph(t, "--model", path, "--emit", "dot")
		if err != nil {
			t.Fatalf("text mode refused: %v", err)
		}
		enveloped, _, err := runGraph(t, "--model", path, "--emit", "dot", "--as=json")
		if err != nil {
			t.Fatalf("json mode refused: %v", err)
		}

		var env struct {
			Data struct {
				Dot string `json:"dot"`
			} `json:"data"`
		}
		if uerr := json.Unmarshal([]byte(strings.TrimSpace(enveloped)), &env); uerr != nil {
			t.Fatalf("stdout is not one JSON object: %v", uerr)
		}

		// S5's normative accounting: text-mode stdout is `document + "\n"`,
		// and the `dot` member carries the document WITHOUT that gateway
		// newline. So the comparison is against the document, not the
		// stream.
		wantDocument := strings.TrimSuffix(text, "\n")
		if env.Data.Dot != wantDocument {
			t.Errorf("the `jq -r .data.dot` unwrap does not reproduce the "+
				"text-mode DOT document byte-for-byte after accounting for "+
				"the one trailing newline F1 pins.\n--- data.dot ---\n%q\n"+
				"--- text minus newline ---\n%q", env.Data.Dot, wantDocument)
		}
	})
}

// REQ-60: "All four `--as`×`--emit` cells are defined — none refused, none
// dead"
// REQ-63: "`--emit` selects the DOCUMENT, `--as` selects the ENVELOPE; the
// 2×2 composes with no refused or dead cell"
// REQ-96 (the same authority row).
// HAPPY PATH — each cell produces its own defined output, and the two
// selectors are independent.
func TestReq60And63And96_AllFourCellsAreDefinedAndIndependent(t *testing.T) {
	requireGraphVerb(t)

	path := writeModel(t, legalModel)
	got := map[string]string{}
	for _, cell := range emitCells {
		args := append([]string{"--model", path}, cell.args...)
		stdout, _, err := runGraph(t, args...)
		if err != nil {
			t.Errorf("the %s cell refused: %v — C5 defines all four "+
				"`--as`×`--emit` cells: none refused, none dead",
				cell.name, err)
			continue
		}
		if strings.TrimSpace(stdout) == "" {
			t.Errorf("the %s cell produced no output; no cell is DEAD",
				cell.name)
		}
		got[cell.name] = stdout
	}

	// The selectors are independent: changing `--emit` changes the
	// DOCUMENT, changing `--as` changes the ENVELOPE. So no two cells
	// coincide.
	for a, left := range got {
		for b, right := range got {
			if a < b && left == right {
				t.Errorf("the %s and %s cells produced identical output; "+
					"`--emit` selects the DOCUMENT and `--as` the ENVELOPE, "+
					"so the 2×2 has four distinct cells", a, b)
			}
		}
	}
}

// REQ-62: "The verb offers NO caller-controlled projection of its success
// payload"
// BOUNDARY — the present-tense half is a testable negative: no projection
// flag on the verb. (The conditional half binds a FUTURE revision and is
// EXCLUDED from this suite by the Phase 0 audit.)
func TestReq62_TheVerbOffersNoCallerControlledProjection(t *testing.T) {
	requireGraphVerb(t)

	cmd := lookupGraphCmd(t)
	for _, projection := range []string{
		"only", "fields", "select", "project", "pick", "omit", "format",
	} {
		if f := cmd.Flags().Lookup(projection); f != nil {
			t.Errorf("the verb registers a --%s flag; C5 gives it NO "+
				"caller-controlled projection of its success payload. A "+
				"later revision that adds one must conform to JDR 0002 §D1",
				projection)
		}
	}
}

// REQ-83: "stdout carries NO document on every refusing arm (C4's
// never-a-partial-document rule, asserted here and not only in the
// mini-check disposition table)."
// REQ-82: "Selection-arm refusals … Expected: lint's code spellings
// verbatim plus `flag-invalid-value` naming `emit`; exit codes stay the
// existing 0/2 mapping (C1)"
// ADVERSARIAL — EVERY refusing arm, not just the ceiling one.
func TestReq82And83_NoRefusingArmPutsADocumentOnStdout(t *testing.T) {
	requireGraphVerb(t)

	good := writeModel(t, legalModel)
	missing := writeModel(t, legalModel) + ".absent"

	for _, tc := range []struct {
		name string
		args []string
		code string
	}{
		{"both flags", []string{"--model", good, "--flow", "x"}, "flag-mutually-exclusive"},
		{"neither flag", nil, "flag-required"},
		{"flow alone", []string{"--flow", "x"}, "flag-invalid-value"},
		{"unreadable", []string{"--model", missing}, "model-unreadable"},
		{"unknown emit", []string{"--model", good, "--emit", "xml"}, "flag-invalid-value"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, mode := range []string{"--as=text", "--as=json"} {
				stdout, _, err := runGraph(t, append(append([]string{}, tc.args...), mode)...)
				assertRefusalCode(t, err, tc.code)
				assertExitTwo(t, err)
				assertNoDocumentOnStdout(t, stdout, tc.name+" "+mode)
			}
		})
	}
}

// REQ-86: "Model loads, lint clean | 0 | — | document on stdout | loud"
// HAPPY PATH — the success disposition row.
func TestReq86_ACleanModelExportsADocumentAtExitZero(t *testing.T) {
	requireGraphVerb(t)

	path := writeModel(t, legalModel)
	stdout, _, err := runGraph(t, "--model", path)
	if err != nil {
		t.Fatalf("a clean model refused: %v", err)
	}
	if !strings.Contains(stdout, schemaMarker) {
		t.Errorf("stdout carries no document: %s", stdout)
	}
}

// REQ-97: "refusals reuse lint's arm codes (C1) and the new
// `graph-export-too-large` refusal names the ceiling and remedy … a killed
// process leaves no terminal envelope (existing contract)."
// REQ-98: "an envelope-sniffing wrapper pointed at text-mode output
// mis-parses the bare document — by contract it must use `--as=json` (C5);
// the four-cell behavior is documented and tested."
// REQ-99: "nondeterministic bytes across runs — caught by the MVV's
// double-emit byte compare, never shipped silently; a partial graph after
// an incomplete traversal — structurally impossible, C4 refuses instead."
// HAPPY PATH — the failure-mode clauses whose observable is that the
// documented behaviour holds. Each is asserted by the suites those clauses
// name; this test pins the one residue they share: under `--as=json`
// EXACTLY ONE terminal record reaches stdout, success or refusal.
func TestReq97And98And99_ExactlyOneTerminalRecordUnderJSONMode(t *testing.T) {
	requireGraphVerb(t)

	good := writeModel(t, legalModel)
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"success", []string{"--model", good, "--as=json"}},
		{"success dot", []string{"--model", good, "--emit", "dot", "--as=json"}},
		{"refusal", []string{"--as=json"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stdout, _, _ := runGraph(t, tc.args...)
			assertExactlyOneEnvelopeLine(t, stdout)
		})
	}
}

// assertExactlyOneEnvelopeLine asserts stdout carries exactly one NDJSON
// terminal record — the one-envelope-per-invocation rule `0005:C1` fixes
// and this record rides rather than amends.
func assertExactlyOneEnvelopeLine(t *testing.T, stdout string) {
	t.Helper()

	lines := 0
	for line := range strings.Lines(stdout) {
		if strings.TrimSpace(line) != "" {
			lines++
		}
	}
	if lines != 1 {
		t.Errorf("stdout carries %d non-empty lines under `--as=json`; the "+
			"envelope contract is EXACTLY ONE terminal record per "+
			"invocation (`0005:C1`), which this record rides rather than "+
			"amends\n%s", lines, stdout)
	}
}

// REQ-13 (the gateway seams the verb routes through), REQ-64.
// HAPPY PATH — the payload satisfies the gateway's TextLiner seam, which
// is what makes the text branch emit one canonical string rather than the
// generic field-per-line rendering.
func TestReq64_TheSuccessPayloadRidesTheTextLinerSeam(t *testing.T) {
	requireGraphVerb(t)

	// The seam is an interface the gateway asks for. Its observable is the
	// text branch's output shape: one canonical document, not a sorted
	// leaf-per-field spread.
	var _ respond.TextLiner = (*textLinerProbe)(nil)

	path := writeModel(t, legalModel)
	stdout, _, err := runGraph(t, "--model", path)
	if err != nil {
		t.Fatalf("export refused: %v", err)
	}
	// The generic renderer spreads a payload over `key: value` lines; the
	// TextLiner branch emits the document itself.
	if !strings.HasPrefix(strings.TrimSpace(stdout), "{") {
		t.Errorf("text-mode stdout does not open with the document; a "+
			"payload that did NOT satisfy TextLiner would be spread "+
			"field-per-line by the gateway's generic renderer\n%s", stdout)
	}
}

// textLinerProbe exists only to pin, at compile time, the seam's shape as
// this record relies on it.
type textLinerProbe struct{}

func (textLinerProbe) TextLine() string { return "" }
