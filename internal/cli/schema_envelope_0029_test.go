package cli

// RDR 0029 — the wire observations: `schema_version` on both terminal
// records, and the properties C1 fixes about how it rides them.
//
// These drive the root command through ExecuteAndEmit, the production
// emission path, and read the RAW JSON so an absent key fails rather than
// decoding to a satisfied zero value. Every assertion here is ADDITIVE per
// S1/S2 — presence and value of the one added key, never an exact key set
// over a record whose optional siblings vary by run.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// rawObject decodes one emitted NDJSON line into its top-level key/value
// map. The raw form is the oracle: a struct decode would turn an omitted
// key into a zero value and hide the very absence under test.
func rawObject(t *testing.T, stdout string) map[string]json.RawMessage {
	t.Helper()

	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &raw); err != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", err, stdout)
	}
	return raw
}

// schemaVersionOf returns the decoded top-level `schema_version` string and
// whether the key was present at all.
func schemaVersionOf(t *testing.T, raw map[string]json.RawMessage) (string, bool) {
	t.Helper()

	v, ok := raw["schema_version"]
	if !ok {
		return "", false
	}
	var s string
	if err := json.Unmarshal(v, &s); err != nil {
		t.Fatalf("schema_version is not a JSON string: %v (raw %s)", err, v)
	}
	return s, true
}

// REQ-1: "The `--as=json` terminal envelope carries a `schema_version`
// string field of the form `MAJOR.MINOR`, present on both terminal records:
// the `ok` envelope (`internal/cli/respond::Success`) and the refusal
// record (`internal/cli/clierr::CLIError`)."
// REQ-42: "`intrastate lint --model <clean-model> --as=json` emits
// top-level keys `data,schema_version,type` with `\"schema_version\":\"0.1\"`.
// The criterion is ADDITIVE, not an exact key set"
// REQ-8: "The schema version begins at `\"0.1\"` and tracks the wire, not
// the binary."
// HAPPY PATH
func TestReq1And8And42_TheOkEnvelopeCarriesSchemaVersion(t *testing.T) {
	path := writeModel(t, legalModel)

	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err != nil {
		t.Fatalf("the clean model failed: %v\nstdout:\n%s", err, stdout)
	}

	raw := rawObject(t, stdout)
	got, present := schemaVersionOf(t, raw)
	if !present {
		t.Fatalf("the `ok` envelope carries no top-level `schema_version` "+
			"key. C1 puts the field on BOTH terminal records, and S1 makes "+
			"the criterion additive: this is the one key the change adds "+
			"over a captured `data,type` baseline.\nstdout:\n%s", stdout)
	}
	if got != "0.1" {
		t.Errorf("schema_version = %q; want %q — C1 begins the schema "+
			"version at \"0.1\", tracking the wire and not the binary",
			got, "0.1")
	}

	// ADDITIVE, per S1: the keys this RDR names must be present; the key
	// set is deliberately NOT asserted closed, because `notes` and
	// `warnings` are `omitempty` siblings that appear on runs unrelated to
	// this record.
	for _, k := range []string{"data", "schema_version", "type"} {
		if _, ok := raw[k]; !ok {
			t.Errorf("the `ok` envelope omits the top-level key %q; S1 names "+
				"data,schema_version,type", k)
		}
	}
}

// REQ-1 / REQ-43: "the assertion is `schema_version` PRESENT on whichever
// failure is provoked — never an exact key-set match over a set that varies
// by failure … No `type` key appears on any of them; the record stays the
// bare `CLIError`."
// HAPPY PATH — the refusal half of C1's "both terminal records".
func TestReq1And43_TheRefusalRecordCarriesSchemaVersionAndStaysBare(t *testing.T) {
	path := writeModel(t, illegalModel)

	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err == nil {
		t.Fatalf("the illegal model succeeded; stdout:\n%s", stdout)
	}

	raw := rawObject(t, stdout)
	got, present := schemaVersionOf(t, raw)
	if !present {
		t.Fatalf("the refusal record carries no top-level `schema_version` "+
			"key. C1 puts the field on the bare `CLIError` too — it is the "+
			"one field the two structurally asymmetric terminal records "+
			"share.\nstdout:\n%s", stdout)
	}
	if got != "0.1" {
		t.Errorf("refusal schema_version = %q; want %q", got, "0.1")
	}

	// The record stays the bare CLIError: a refusal carries no `type` and
	// no wrapper (`0005:C1`). Adding the version must not have wrapped it.
	if _, ok := raw["type"]; ok {
		t.Errorf("the refusal record carries a `type` key; S2 fixes that no "+
			"`type` appears on any provoked failure — the record stays the "+
			"bare `CLIError`.\nstdout:\n%s", stdout)
	}
	if _, ok := raw["code"]; !ok {
		t.Errorf("the refusal record carries no `code`; it is the bare "+
			"CLIError and `code` is what discriminates it from the `ok` "+
			"envelope.\nstdout:\n%s", stdout)
	}
}

// REQ-43: the same additive criterion across DIFFERENT provocations — "the
// baseline is per-provocation, not one set".
// INPUT EDGE
func TestReq43_SchemaVersionRidesEveryProvokedFailureShape(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"unreadable-model", []string{"lint", "--model",
			filepath.Join(t.TempDir(), "absent.toml"), "--as=json"}},
		{"blocking-findings", nil}, // filled below; needs a written fixture
		{"invalid-flag-value", []string{"version", "--as=json", "--as=yaml"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			args := tc.args
			if args == nil {
				args = []string{"lint", "--model",
					writeModel(t, illegalModel), "--as=json"}
			}

			stdout, _, err := runCmd(t, args...)
			if err == nil {
				t.Skipf("this provocation did not fail; stdout:\n%s", stdout)
			}
			if strings.TrimSpace(stdout) == "" {
				t.Skipf("this provocation emitted nothing on stdout")
			}

			raw := rawObject(t, stdout)
			if _, present := schemaVersionOf(t, raw); !present {
				t.Errorf("this failure shape carries no `schema_version`; "+
					"S2 asserts presence on WHICHEVER failure is provoked, "+
					"never an exact key-set match over a set that varies by "+
					"failure.\nstdout:\n%s", stdout)
			}
		})
	}
}

// REQ-5: "`schema_version` rides the TOP LEVEL of each terminal record
// only. It is never projected into `data`"
// ADVERSARIAL — projecting it into `data` would collide with a verb's own
// payload; `intrastate version --as=json` already emits `data.version`
// meaning build identity.
func TestReq5_SchemaVersionIsNeverProjectedIntoData(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"lint", []string{"lint", "--model", writeModel(t, legalModel), "--as=json"}},
		{"version", []string{"version", "--as=json"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stdout, _, err := runCmd(t, tc.args...)
			if err != nil {
				t.Fatalf("%s failed: %v\nstdout:\n%s", tc.name, err, stdout)
			}

			raw := rawObject(t, stdout)

			// The discriminating precondition. Until the field rides the
			// TOP LEVEL there is no projection to rule out, and a bare
			// "absent from `data`" assertion would pass on a tree that
			// carries the field nowhere at all.
			if _, present := schemaVersionOf(t, raw); !present {
				t.Fatalf("the %s envelope carries no top-level "+
					"`schema_version`, so \"never projected into `data`\" is "+
					"unasserted: absent-from-`data` and absent-everywhere are "+
					"indistinguishable here.\nstdout:\n%s", tc.name, stdout)
			}

			dataRaw, ok := raw["data"]
			if !ok {
				t.Fatalf("no `data` key on the %s envelope:\n%s", tc.name, stdout)
			}
			var data map[string]json.RawMessage
			if err := json.Unmarshal(dataRaw, &data); err != nil {
				t.Fatalf("`data` is not an object: %v\n%s", err, stdout)
			}
			if _, leaked := data["schema_version"]; leaked {
				t.Errorf("`data.schema_version` is present on the %s payload; "+
					"C1 fixes that the field rides the TOP LEVEL only and is "+
					"never projected into `data`, where it would collide with "+
					"a verb's own payload.\nstdout:\n%s", tc.name, stdout)
			}
		})
	}
}

// REQ-27: "The `version` verb's payload field names (`version`, `commit`,
// `date` — `internal/version::Info`) are `frozen`"
// BOUNDARY — the keys an agent reads to pin a build.
func TestReq27_TheVersionPayloadFieldNamesAreFrozen(t *testing.T) {
	stdout, _, err := runCmd(t, "version", "--as=json")
	if err != nil {
		t.Fatalf("version --as=json failed: %v\nstdout:\n%s", err, stdout)
	}

	raw := rawObject(t, stdout)
	var data map[string]json.RawMessage
	if err := json.Unmarshal(raw["data"], &data); err != nil {
		t.Fatalf("`data` is not an object: %v\n%s", err, stdout)
	}

	for _, k := range []string{"version", "commit", "date"} {
		if _, ok := data[k]; !ok {
			t.Errorf("the version payload omits %q; C4 freezes these three "+
				"field names — they are the keys an agent reads to pin a "+
				"build, so a rename breaks the surface a consumer uses to "+
				"decide what it is running.\nstdout:\n%s", k, stdout)
		}
	}
}

// REQ-6: "`schema_version` is NOT `omitempty` on either record"
// REQ-7: "The refusal envelope still has exactly one structured field,
// `findings`."
// BOUNDARY — read off the struct tags, which is where `omitempty` lives.
func TestReq6And7_SchemaVersionIsNotOmitemptyAndFindingsStaysTheOneStructuredField(t *testing.T) {
	for _, tc := range []struct {
		name string
		path []string
	}{
		{"respond.Success", []string{"internal", "cli", "respond", "respond.go"}},
		{"clierr.CLIError", []string{"internal", "cli", "clierr", "clierr.go"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			src := readSource(t, pkgDir(t, tc.path...))

			idx := strings.Index(src, `json:"schema_version`)
			if idx < 0 {
				t.Fatalf("%s declares no `schema_version` json tag; C1 puts "+
					"the field on both terminal records", tc.name)
			}
			// The tag runs to the closing quote; `omitempty` inside it is
			// the defect — the field must always be on the wire, which is
			// what makes it a version a consumer can rely on rather than
			// one it must handle the absence of.
			tag := src[idx:]
			if end := strings.Index(tag, `"`+"`"); end >= 0 {
				tag = tag[:end]
			}
			if strings.Contains(tag, "omitempty") {
				t.Errorf("%s tags schema_version %q; C1 fixes it as NOT "+
					"`omitempty` on either record", tc.name, tag)
			}
		})
	}

	// REQ-7: the refusal envelope's structured-field budget is unchanged.
	// `schema_version` is a scalar version marker carrying no discriminator
	// and no payload, so it does not consume that budget.
	clierrSrc := readSource(t, pkgDir(t, "internal", "cli", "clierr", "clierr.go"))
	structured := 0
	for _, line := range strings.Split(clierrSrc, "\n") {
		if strings.Contains(line, "[]Finding") && strings.Contains(line, "json:") {
			structured++
		}
	}
	if structured != 1 {
		t.Errorf("the CLIError carries %d structured `[]Finding` fields; "+
			"C1 keeps the refusal envelope at exactly one, `findings`",
			structured)
	}
}

// REQ-58: "`schema_version` is a string (`\"1.0\"`), not a pair of integers
// and not a number"
// REQ-1 (form): "a `schema_version` string field of the form `MAJOR.MINOR`"
// INPUT EDGE — a JSON number would silently lose a trailing zero and make
// "1.10" and "1.1" the same value.
func TestReq1And58_SchemaVersionIsAStringOfTheFormMajorMinor(t *testing.T) {
	path := writeModel(t, legalModel)
	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err != nil {
		t.Fatalf("the clean model failed: %v\nstdout:\n%s", err, stdout)
	}

	raw := rawObject(t, stdout)
	v, ok := raw["schema_version"]
	if !ok {
		t.Fatalf("no `schema_version` key:\n%s", stdout)
	}

	// The RAW bytes must be a quoted JSON string, not a bare number.
	if !strings.HasPrefix(strings.TrimSpace(string(v)), `"`) {
		t.Fatalf("schema_version is emitted as %s; C1 and D-wire-byte-format "+
			"fix it as a STRING, not a pair of integers and not a number — a "+
			"number cannot hold `MAJOR.MINOR` without losing a trailing zero",
			v)
	}

	var s string
	if err := json.Unmarshal(v, &s); err != nil {
		t.Fatalf("schema_version does not decode as a string: %v", err)
	}
	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		t.Fatalf("schema_version = %q; want the form MAJOR.MINOR (exactly "+
			"two dot-separated components)", s)
	}
	for _, p := range parts {
		if p == "" || strings.TrimLeft(p, "0123456789") != "" {
			t.Errorf("schema_version = %q; each component of MAJOR.MINOR must "+
				"be numeric", s)
		}
	}
}

// REQ-60: "Add the field at the single output gateway — `respond.Success`
// and `clierr.CLIError` — so every verb inherits it without per-verb work"
// DOMAIN EDGE — the property is that a verb the RDR never names still
// carries the field, which is what "inherits without per-verb work" means.
func TestReq60_EveryVerbInheritsTheFieldFromTheGateway(t *testing.T) {
	// `version` is a verb this RDR does not touch and whose payload has
	// nothing to do with the wire schema. If it carries the field, the
	// field was added at the gateway rather than per verb.
	stdout, _, err := runCmd(t, "version", "--as=json")
	if err != nil {
		t.Fatalf("version --as=json failed: %v\nstdout:\n%s", err, stdout)
	}

	raw := rawObject(t, stdout)
	if _, present := schemaVersionOf(t, raw); !present {
		t.Errorf("`intrastate version --as=json` carries no `schema_version`. "+
			"The field is added at the SINGLE output gateway so every verb "+
			"inherits it without per-verb work; a verb missing it means the "+
			"field was bolted onto lint rather than onto `respond.Success`."+
			"\nstdout:\n%s", stdout)
	}
}

// REQ-44: "`--plan` still decodes a `flow resolve --as json` envelope that
// now carries `schema_version`, unchanged."
// REQ-10: "It MUST NOT be made strict: this is an emitted envelope arriving
// back, not an authored document"
// ADVERSARIAL — the A3 regression guard for the strictest DECODER the repo
// holds. A `DisallowUnknownFields` here would break the `--plan` pipe on
// the very release that adds the field.
func TestReq10And44_ThePlanDecoderToleratesTheAddedField(t *testing.T) {
	// The decoder must not be strict. Strictness is a source property of
	// the reader, and it is the one thing C1 forbids changing here.
	src := readSource(t, pkgDir(t, "internal", "cli", "flow_input.go"))
	if strings.Contains(src, "DisallowUnknownFields") {
		t.Errorf("flow_input.go calls DisallowUnknownFields; C1 fixes that " +
			"`planEnvelope` MUST NOT be made strict — it reads an emitted " +
			"envelope arriving back, not an authored document, so adding " +
			"strictness breaks the --plan pipe on the very release that " +
			"adds schema_version")
	}

	// And the behavioural half: an envelope carrying the field decodes to
	// the same plan as one without it.
	const withField = `{"type":"ok","schema_version":"0.1","data":` +
		`{"writes":{"status":"b"},"clear":[]}}`
	const withoutField = `{"type":"ok","data":` +
		`{"writes":{"status":"b"},"clear":[]}}`

	dir := t.TempDir()
	withPath := filepath.Join(dir, "with.json")
	withoutPath := filepath.Join(dir, "without.json")
	for path, body := range map[string]string{
		withPath:    withField,
		withoutPath: withoutField,
	} {
		if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}

	plans := map[string]*carriedPlan{}
	for name, path := range map[string]string{
		"with": withPath, "without": withoutPath,
	} {
		cmd := NewRootCmd()
		got, ce := readPlan(cmd, path)
		if ce != nil {
			t.Fatalf("readPlan refused the %s-schema_version envelope: %v; "+
				"S3 makes this the A3 regression guard — the reader is "+
				"tolerant by construction and must survive the added key",
				name, ce)
		}
		plans[name] = got
	}

	if plans["with"] == nil || plans["without"] == nil {
		t.Fatal("readPlan returned a nil plan without refusing")
	}
	if !slices.Equal(plans["with"].Clear, plans["without"].Clear) ||
		len(plans["with"].Writes) != len(plans["without"].Writes) {
		t.Errorf("the decoded plan differs with the field (%+v) and without "+
			"it (%+v); the added key must be inert to this reader",
			plans["with"], plans["without"])
	}
	for k, v := range plans["without"].Writes {
		if plans["with"].Writes[k] != v {
			t.Errorf("write %q decoded as %q with the field and %q without",
				k, plans["with"].Writes[k], v)
		}
	}
}
