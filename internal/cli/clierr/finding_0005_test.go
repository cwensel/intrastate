package clierr_test

// RDR 0005 — the `clierr` type-level obligations this contract lands:
// the two `Finding` fields it ADDS (REQ-14, REQ-16, REQ-18, A-7), the
// `omitempty` amendment to `CLIError.Findings` (REQ-8, REQ-16, A-6), and
// the shared non-HTML-escaping encoder both emit sites route through
// (REQ-71, REQ-72, A-8).
//
// These assertions are made HERE rather than through the CLI because they
// are properties of the shared record and its encoder, which RDR 0006 also
// consumes. A test at the verb boundary could be satisfied by a verb that
// happened to spell the fields itself; a test here cannot.

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
)

// REQ-14: "Finding MUST be one flat record with omitempty optional fields,
// carrying at least code, message, param, locator, hint, severity, model,
// rule, key, operator, literal, block, and class"
// REQ-16: "Add `Findings []Finding` … and the flat `Finding` record to
// `internal/cli/clierr` … fields `Code`, `Message`, plus `omitempty`
// `Param`, `Locator`, `Hint`, `Severity`, `Model`, `Rule`, `Key`,
// `Operator`, `Literal`, `Block`, `Class`."
// A-7: "the `Finding` fields this RDR adds are `Param` and `Locator` only."
// BOUNDARY — the thirteen names REQ-14 fixes must all be reachable on the
// record and must all carry their contract wire spellings.
func TestReq14And16_FindingCarriesTheThirteenNamedFieldsFlat(t *testing.T) {
	// The thirteen Go field names REQ-16 fixes, paired with the wire
	// spellings REQ-14 fixes. `Param` and `Locator` are the two this RDR
	// ADDS (A-7); the other eleven are peer-delivered by RDR 0006.
	//
	// Presence is checked REFLECTIVELY rather than by a struct literal so
	// that a missing field is a readable FAILURE naming the field, not a
	// compile break that would take RDR 0006's shipped `clierr` suite down
	// with it.
	want := []struct{ field, wire string }{
		{"Code", "code"},
		{"Message", "message"},
		{"Param", "param"},
		{"Locator", "locator"},
		{"Hint", "hint"},
		{"Severity", "severity"},
		{"Model", "model"},
		{"Rule", "rule"},
		{"Key", "key"},
		{"Operator", "operator"},
		{"Literal", "literal"},
		{"Block", "block"},
		{"Class", "class"},
	}

	rt := reflect.TypeOf(clierr.Finding{})
	for _, w := range want {
		sf, ok := rt.FieldByName(w.field)
		if !ok {
			t.Errorf("clierr.Finding carries no field %s; REQ-14 fixes it as "+
				"one of the thirteen the record carries AT LEAST, and A-7 "+
				"names Param and Locator as this RDR's two additions",
				w.field)
			continue
		}
		if sf.Type.Kind() != reflect.String {
			t.Errorf("Finding.%s is %s; every named field is a flat string "+
				"the transporting package ascribes no meaning to",
				w.field, sf.Type)
		}
		tag := sf.Tag.Get("json")
		name, opts, _ := strings.Cut(tag, ",")
		if name != w.wire {
			t.Errorf("Finding.%s marshals as %q; the contract fixes the wire "+
				"spelling %q", w.field, name, w.wire)
		}
		// Code and Message are required; the other eleven are `omitempty`
		// so each producer populates only the fields it owns (REQ-15).
		optional := w.field != "Code" && w.field != "Message"
		if optional && !strings.Contains(opts, "omitempty") {
			t.Errorf("Finding.%s is not `omitempty`; REQ-14 fixes the "+
				"optional fields as omitempty", w.field)
		}
	}

	// Flatness: no field on the record is a nested struct or map, so no
	// producer can nest its own fields in a sub-object (REQ-15).
	for i := range rt.NumField() {
		sf := rt.Field(i)
		switch sf.Type.Kind() {
		case reflect.Struct, reflect.Map:
			t.Errorf("Finding.%s is a nested %s; the record is FLAT — the "+
				"guard atom's four fields in particular sit flat on it, not "+
				"under an `atom` object", sf.Name, sf.Type.Kind())
		}
	}
}

// REQ-14: "one flat record with omitempty optional fields"
// REQ-129 / `0005:S7`: "unset optional fields are absent from the JSON
// (`omitempty`)".
// BOUNDARY — the zero value marshals to only the two required fields.
func TestReq14_UnsetOptionalFindingFieldsAreAbsentFromTheJSON(t *testing.T) {
	buf, err := json.Marshal(clierr.Finding{Code: "c", Message: "m"})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var wire map[string]json.RawMessage
	if err := json.Unmarshal(buf, &wire); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	for _, optional := range []string{
		"param", "locator", "hint", "severity", "rule", "key", "operator",
		"literal", "block", "class",
	} {
		if _, present := wire[optional]; present {
			t.Errorf("the unset optional field %q is PRESENT on the wire; "+
				"optional fields are `omitempty`, which is what lets each "+
				"producer populate only the fields it owns", optional)
		}
	}
}

// REQ-8: "extended by exactly one omitempty structured field, findings"
// REQ-16 / A-6: `Findings` moves to `json:"findings,omitempty"`.
// BOUNDARY — a CLIError naming no findings omits the key entirely.
func TestReq8And16_CLIErrorFindingsIsOmitEmpty(t *testing.T) {
	buf, err := json.Marshal(&clierr.CLIError{
		Code: "flow-tag-reserved", Message: "m", Param: "recognized",
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var wire map[string]json.RawMessage
	if err := json.Unmarshal(buf, &wire); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if _, present := wire["findings"]; present {
		t.Errorf("`findings` is present on a CLIError that names none; "+
			"REQ-8 makes it the ONE `omitempty` structured field this RDR "+
			"adds\nwire: %s", buf)
	}
	// A CLIError that DOES carry findings still emits them, so `omitempty`
	// has not made the receipt unreachable.
	buf2, err := json.Marshal(&clierr.CLIError{
		Code: "flow-model-invalid", Message: "m",
		Findings: []clierr.Finding{{Code: "malformed_tag_declaration", Message: "m"}},
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(buf2), `"findings"`) {
		t.Errorf("a CLIError carrying findings does not emit them: %s", buf2)
	}
}

// REQ-71: "Every site that emits or compares a canonical set literal MUST
// use the same encoder, so plan-to-request copy-through and read-back
// equality are byte equality."
// REQ-72: "Route `internal/cli/clierr::EmitJSON` and
// `internal/cli/respond::writeJSONLine` through one shared helper using
// `SetEscapeHTML(false)` … Bare `json.Marshal` must not remain on a path a
// set value crosses."
// ADVERSARIAL — `EmitJSON` is one of the two named sites.
func TestReq71And72_EmitJSONRendersAngleAndAmpersandAsThemselves(t *testing.T) {
	// The canonical set literal rides in `param`, which is the shape a
	// wrong-kind `--tag` refusal takes.
	const literal = `["a<b","x&y"]`

	var buf bytes.Buffer
	clierr.EmitJSON(&buf, &clierr.CLIError{
		Code:    "flow-tag-invalid",
		Message: "the value is not a member of the declared domain",
		Param:   literal,
		Findings: []clierr.Finding{{
			Code:    "wrong_kind",
			Message: "a set literal for a scalar key",
			Literal: literal,
		}},
	})

	out := buf.String()
	// The list is the HTML-ESCAPED spellings, which is what the clause
	// forbids. Looping the RAW characters instead would fire exactly when
	// REQ-72 is SATISFIED — the inverted oracle Phase 2 recorded as DEV-2.
	for _, esc := range []string{`\u003c`, `\u003e`, `\u0026`} {
		if strings.Contains(out, esc) {
			t.Errorf("EmitJSON rendered %s; `<`, `>`, and `&` MUST serialize "+
				"as THEMSELVES — this site is still on bare `json.Marshal`\n"+
				"output: %s", esc, out)
		}
	}
	// The literal rides inside a JSON STRING field, whose inner `"` MUST
	// escape per RFC 8259 §7 — no encoder setting changes that, and
	// `SetEscapeHTML(false)` governs only `<`, `>`, `&`. So the carrier is
	// DECODED before comparison; comparing the raw line against the
	// unencoded form asserted something unsatisfiable by construction
	// (DEV-3). What REQ-71 requires is that the VALUE cross byte-identically.
	if got := decodedStringField(t, out, "param"); got != literal {
		t.Errorf("EmitJSON did not carry the canonical literal %s through "+
			"byte-identically; `param` decodes to %s\noutput: %s",
			literal, got, out)
	}

	// Exactly ONE line: the shared helper must not double-write the newline
	// that `json.Encoder` appends of its own accord (A-8).
	if n := strings.Count(strings.TrimRight(out, "\n"), "\n"); n != 0 {
		t.Errorf("EmitJSON wrote %d embedded newlines; the envelope is ONE "+
			"NDJSON line\noutput: %q", n, out)
	}
	if !strings.HasSuffix(out, "\n") {
		t.Errorf("EmitJSON did not terminate the line: %q", out)
	}
	if strings.HasSuffix(out, "\n\n") {
		t.Errorf("EmitJSON double-wrote the trailing newline; "+
			"`json.Encoder` appends one already (A-8): %q", out)
	}
}

// REQ-71 (the equality leg): the canonical literal crosses the JSON emit
// site byte-identically wherever it rides — `param` and a finding's
// `literal` alike — which is what makes plan-to-request copy-through and
// read-back equality byte equality.
//
// The oracle is written against the SHIPPED emit function rather than any
// particular helper signature: REQ-72 fixes the OBLIGATION (one shared
// non-escaping encoder) and leaves the helper's name and shape to the
// implementation.
//
// Only the JSON site is asserted for byte-identity. `EmitText` deliberately
// renders identity values through `strconv.Quote` so no authored value can
// forge a second finding's identity suffix — that is RDR 0006's shipped
// anti-forgery rule, and quoting is not HTML escaping.
// ADVERSARIAL
func TestReq71_TheCanonicalLiteralSurvivesTheJSONEmitSiteByteIdentically(t *testing.T) {
	const literal = `["a<b","x&y"]`

	var out bytes.Buffer
	clierr.EmitJSON(&out, &clierr.CLIError{
		Code: "flow-tag-invalid", Message: "wrong kind", Param: literal,
		Findings: []clierr.Finding{{
			Code: "wrong_kind", Message: "a set literal for a scalar key",
			Literal: literal,
		}},
	})

	// The literal survives at BOTH carriers on the same line: the envelope's
	// `param` and the finding's `literal`.
	//
	// Both are JSON STRING fields, so each is DECODED before comparison —
	// their inner `"` characters escape per RFC 8259 §7 regardless of the
	// encoder's HTML setting, and counting the unencoded form in the raw
	// line asserted something unsatisfiable by construction (DEV-3). The
	// escaping defect under test is caught by the `\u00xx` sweep below,
	// which runs against the RAW bytes and is not laundered by decoding.
	for _, carrier := range []string{"param", "literal"} {
		if got := decodedStringField(t, out.String(), carrier); got != literal {
			t.Errorf("the canonical literal does not cross `%s` "+
				"byte-identically: want %s, got %s; every site that emits or "+
				"compares one MUST use the same encoder\noutput: %s",
				carrier, literal, got, out.String())
		}
	}
	for _, esc := range []string{`\u003c`, `\u003e`, `\u0026`} {
		if strings.Contains(out.String(), esc) {
			t.Errorf("EmitJSON rendered %s; `<`, `>`, and `&` serialize as "+
				"THEMSELVES\noutput: %s", esc, out.String())
		}
	}
}

// REQ-17: "Finding.message MUST be self-sufficient — it MUST render the
// failure readably with no structured field consulted."
// REQ-11: text mode "renders `findings` one-per-line under the message."
// REQ-18: "This RDR populates `code`, `message`, `param`, `locator`, and
// `hint`" — so the text branch must render THIS RDR's identity fields too,
// not only RDR 0006's.
// DOMAIN EDGE — one finding is ONE line, and a `locator` this RDR
// populates must reach the reader on it. A text renderer that enumerates
// only the shipped identity fields drops the file:line a model-load
// finding exists to convey.
func TestReq11And17And18_TextRendersThisRDRsIdentityFieldsOnePerLine(t *testing.T) {
	// Set via reflection so a not-yet-added field is a readable FAILURE
	// naming it, rather than a compile break taking RDR 0006's shipped
	// `clierr` suite down with it.
	first := newFindingWith(t, map[string]string{
		"Code":    "malformed_tag_declaration",
		"Message": "tag `status` declares an unknown kind",
		"Locator": "flow.toml:12",
	})
	second := newFindingWith(t, map[string]string{
		"Code":    "reserved_tag_key",
		"Message": "tag key `recognized` is reserved",
		"Locator": "flow.toml:20",
		"Param":   "model",
	})
	findings := []clierr.Finding{first, second}

	var buf bytes.Buffer
	clierr.EmitFindingsText(&buf, findings)

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != len(findings) {
		t.Fatalf("rendered %d lines for %d findings; findings render "+
			"ONE-PER-LINE\noutput:\n%s", len(lines), len(findings),
			buf.String())
	}

	for i, f := range findings {
		if !strings.Contains(lines[i], f.Code) {
			t.Errorf("line %d omits the code %q: %q", i, f.Code, lines[i])
		}
		if !strings.Contains(lines[i], f.Message) {
			t.Errorf("line %d omits the message %q, which alone must render "+
				"the failure readably: %q", i, f.Message, lines[i])
		}
	}

	// The `locator` this RDR populates reaches the reader. Without it a
	// model-load finding cannot be acted on from text output at all.
	for i, want := range []string{"flow.toml:12", "flow.toml:20"} {
		if !strings.Contains(lines[i], want) {
			t.Errorf("line %d omits the `locator` %q that this RDR "+
				"populates; the text branch renders the identity fields a "+
				"reader needs to act without re-running in JSON mode: %q",
				i, want, lines[i])
		}
	}
	if !strings.Contains(lines[1], "model") {
		t.Errorf("line 1 omits the `param` %q this RDR populates: %q",
			"model", lines[1])
	}
}

// newFindingWith builds a clierr.Finding by field NAME, so a field this RDR
// has not yet added fails the test with a message naming it instead of
// breaking the package build.
func newFindingWith(t *testing.T, fields map[string]string) clierr.Finding {
	t.Helper()

	var f clierr.Finding
	v := reflect.ValueOf(&f).Elem()
	for name, value := range fields {
		fv := v.FieldByName(name)
		if !fv.IsValid() {
			t.Fatalf("clierr.Finding carries no field %s; REQ-16 fixes it as "+
				"one of the record's fields and A-7 names Param and Locator "+
				"as this RDR's two additions", name)
		}
		fv.SetString(value)
	}
	return f
}

// decodedStringField returns the DECODED value of the first JSON string
// field named key, searched anywhere in one emitted envelope line.
//
// A canonical set literal rides in `CLIError.Param` and `Finding.Literal`,
// both Go `string` fields, so it reaches the wire as a JSON string whose
// inner `"` characters are escaped per RFC 8259 §7. That escaping is
// mandatory and is unrelated to `SetEscapeHTML(false)`, which governs only
// `<`, `>`, and `&`. REQ-71's obligation is that the VALUE cross the emit
// site byte-identically, so the carrier must be decoded before comparison;
// searching the raw line for the unencoded literal asserts something no
// conforming encoder can satisfy (DEV-3).
//
// Decoding here does not launder the escaping defect under test: each
// caller sweeps the RAW line for the `\u00xx` spellings separately.
func decodedStringField(t *testing.T, line, key string) string {
	t.Helper()

	var walk func(raw json.RawMessage) (string, bool)
	walk = func(raw json.RawMessage) (string, bool) {
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(raw, &obj); err == nil {
			if v, held := obj[key]; held {
				var decoded string
				if err := json.Unmarshal(v, &decoded); err == nil {
					return decoded, true
				}
			}
			for _, v := range obj {
				if found, ok := walk(v); ok {
					return found, true
				}
			}
			return "", false
		}
		var arr []json.RawMessage
		if err := json.Unmarshal(raw, &arr); err == nil {
			for _, e := range arr {
				if found, ok := walk(e); ok {
					return found, true
				}
			}
		}
		return "", false
	}

	found, ok := walk(json.RawMessage(strings.TrimSpace(line)))
	if !ok {
		t.Fatalf("the emitted line carries no string field %q:\n%s", key, line)
	}
	return found
}
