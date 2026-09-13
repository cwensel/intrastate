package cli

// RDR 0029 — the Minimum Viable Validation, run end to end.
//
// "An agent pins version N, parses the envelope, and survives N+1 across
// exactly the event the seed named — a new lint finding code appearing."
//
// The MVV's own oracle table fixes what each step may pass on: step 4 is
// "the promise's load-bearing observation, so it may not pass by absence of
// error". Every assertion below is BY VALUE off the emitted envelope's own
// bytes — exit code, `type`, the finding's `severity`, and the version
// string's two components — never a green exit or a "did not error".

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/graphlint"
)

// goldenPath0029 is the 0023 envelope golden S3 requires be re-captured in
// the same commit that adds the field.
const goldenPath0029 = "docs/rdr/0023-resolve-envelope-projection/artifacts/" +
	"mvv-step1-default-golden.json"

// REQ-45: "that golden is re-captured in the same commit that adds the
// field, and the re-captured file differs from its predecessor by exactly
// the one `schema_version` key"
// ADVERSARIAL — "A diff showing any other key changed means this RDR
// altered an envelope it promised not to touch."
func TestReq45_TheGoldenIsRecapturedDifferingByExactlyTheOneKey(t *testing.T) {
	body, err := os.ReadFile(filepath.Join(repoRootFor(t), goldenPath0029))
	if err != nil {
		t.Fatalf("the 0023 golden is not checked in at %s: %v",
			goldenPath0029, err)
	}

	var envelope map[string]json.RawMessage
	if jerr := json.Unmarshal(body, &envelope); jerr != nil {
		t.Fatalf("the golden is not one JSON object: %v", jerr)
	}

	// The re-captured golden must carry the field. Until it does, the
	// byte-identity assertion at flow_mvv_0023_test.go fails against a
	// non-`omitempty` addition — which is precisely why Step 1 re-captures
	// it in the same commit.
	if _, ok := envelope["schema_version"]; !ok {
		t.Fatalf("%s carries no `schema_version`. `0029:S3` requires the "+
			"golden be RE-CAPTURED in the same commit that adds the field: "+
			"it is asserted by byte-identity, so a non-`omitempty` field "+
			"fails it until the baseline is retaken.", goldenPath0029)
	}

	// And it must differ from its predecessor by EXACTLY that key: the
	// predecessor's two keys are `type` and `data`, and nothing else may
	// have moved.
	got := make([]string, 0, len(envelope))
	for k := range envelope {
		got = append(got, k)
	}
	for _, want := range []string{"type", "data"} {
		if _, ok := envelope[want]; !ok {
			t.Errorf("the re-captured golden lost the key %q; it must differ "+
				"from its predecessor by exactly the one `schema_version` "+
				"key. A diff showing any other key changed means this RDR "+
				"altered an envelope it promised not to touch (got %v)",
				want, got)
		}
	}
	if len(envelope) != 3 {
		t.Errorf("the re-captured golden carries %d top-level keys (%v); the "+
			"predecessor carried `type` and `data`, so exactly three is the "+
			"one-key difference S3 allows", len(envelope), got)
	}
}

// REQ-47: "exit 2 and the record shape changes from the `ok` envelope to
// the bare `CLIError`"
// REQ-20 (the observable half): promotion to `blocking` IS the
// verdict-changing event.
// BOUNDARY — this is the control that makes the `info` run meaningful:
// "same input, same code, different severity, different verdict".
func TestReq47_PromotionToBlockingFlipsTheExitAndReplacesTheRecord(t *testing.T) {
	// The `info` side: a model carrying advisory findings only succeeds.
	advisory := writeModel(t, advisoryModel)
	okOut, _, okErr := runCmd(t, "lint", "--model", advisory, "--as=json")
	if okErr != nil {
		t.Fatalf("the advisory model failed: %v\nstdout:\n%s", okErr, okOut)
	}
	okRaw := rawObject(t, okOut)
	if _, ok := okRaw["type"]; !ok {
		t.Errorf("the `info` run emitted no `type` key:\n%s", okOut)
	}

	// Both records must carry the version, or this comparison is about
	// today's already-correct shapes and says nothing about this RDR: the
	// promotion event is observed ACROSS a versioned wire.
	if _, present := schemaVersionOf(t, okRaw); !present {
		t.Fatalf("the `info` run carries no `schema_version`; the "+
			"introduce/promote contrast is observed across a versioned wire, "+
			"so without the field this test re-asserts the pre-existing "+
			"envelope shapes only.\nstdout:\n%s", okOut)
	}

	// The `blocking` side: the record shape CHANGES — the `ok` envelope is
	// replaced by the bare `CLIError`, which carries `code` and no `type`.
	blocking := writeModel(t, illegalModel)
	failOut, _, failErr := runCmd(t, "lint", "--model", blocking, "--as=json")
	if failErr == nil {
		t.Fatalf("the blocking model succeeded:\n%s", failOut)
	}

	failRaw := rawObject(t, failOut)
	if _, ok := failRaw["type"]; ok {
		t.Errorf("the blocking run kept a `type` key; the shape change IS "+
			"the disclosed event C3 governs — the `ok` envelope is replaced "+
			"by the bare `CLIError`.\nstdout:\n%s", failOut)
	}
	if _, ok := failRaw["code"]; !ok {
		t.Errorf("the blocking run emitted no `code`:\n%s", failOut)
	}
	if _, present := schemaVersionOf(t, failRaw); !present {
		t.Errorf("the blocking run carries no `schema_version`; C1 puts the "+
			"field on BOTH terminal records, so a consumer survives the "+
			"shape change by reading the version off either.\nstdout:\n%s",
			failOut)
	}

	// "That is a more disruptive change for a consumer than a `type` value
	// flipping to `\"failed\"`, which is what makes promotion the
	// verdict-changing event C3 requires be disclosed, and introduction not."
	if got := exitCodeOf(t, failErr); got != 2 {
		t.Errorf("the blocking run exited %d; want 2", got)
	}
}

// exitCodeOf reads the process exit code a returned error maps to.
func exitCodeOf(t *testing.T, err error) int {
	t.Helper()
	return clierr.ExitCodeFor(err)
}

// REQ-MVV: "1. Build at the current HEAD with `schema_version` implemented;
// run `intrastate lint --model <clean-model> --as=json` and record the full
// envelope. It reports `\"schema_version\":\"0.1\"` and no findings. 2. Add a
// new graph-lint finding code at severity `info` that fires on the model
// from step 1. 3. Re-run the same command against the same unmodified
// model. 4. **Expected end state**: the exit code is unchanged (0), the
// envelope's `type` is still `ok`, and the new finding appears in
// `data.findings` carrying `\"severity\":\"info\"`. The schema version's
// MAJOR is unchanged — a growing-vocabulary member is an additive change,
// so C1 moves the minor (`\"0.1\"` → `\"0.2\"`) and never the major."
// HAPPY PATH — the runnable end-to-end battery.
func TestMVV0029_AnAgentPinsAVersionAndSurvivesANewFindingCode(t *testing.T) {
	// --- step 1: the baseline envelope over a clean model --------------
	//
	// `models/rdr.toml` is the checked-in model the graph-lint gate lints
	// and it reports no findings, which is the MVV's "clean model".
	clean := filepath.Join(repoRootFor(t), "models", "rdr.toml")
	if _, err := os.Stat(clean); err != nil {
		t.Fatalf("the checked-in clean model is absent: %v", err)
	}

	step1Out, _, step1Err := runCmd(t, "lint", "--model", clean, "--as=json")
	if step1Err != nil {
		t.Fatalf("the clean model failed: %v\nstdout:\n%s", step1Err, step1Out)
	}

	step1 := rawObject(t, step1Out)
	baseline, present := schemaVersionOf(t, step1)
	if !present {
		t.Fatalf("MVV step 1: the envelope carries no `schema_version`. The "+
			"oracle's own negative control is that the PRE-CHANGE baseline "+
			"has no such key and MUST fail this row.\nstdout:\n%s", step1Out)
	}
	if baseline != "0.1" {
		t.Errorf("MVV step 1: schema_version = %q; want %q", baseline, "0.1")
	}

	// "and no findings" — the receipt is the emitted empty list, read off
	// the raw bytes so an omitted key cannot decode to a satisfied nil.
	var data map[string]json.RawMessage
	if jerr := json.Unmarshal(step1["data"], &data); jerr != nil {
		t.Fatalf("MVV step 1: `data` is not an object: %v\n%s", jerr, step1Out)
	}
	if got := strings.TrimSpace(string(data["findings"])); got != "[]" {
		t.Errorf("MVV step 1: data.findings = %s; want the empty-list "+
			"receipt `[]` on the clean model", got)
	}
	if step1Type := decodeString(t, step1["type"]); step1Type != "ok" {
		t.Errorf("MVV step 1: type = %q; want %q", step1Type, "ok")
	}

	// --- steps 2 and 3: a new `info` code fires on the same model -------
	//
	// Step 2 adds "a new graph-lint finding code at severity `info` that
	// fires on the model from step 1". The advisory tier is `growing`, so
	// a new member MAY fire on input that previously produced no such
	// finding — that licence is the whole point of the tier.
	//
	// The code is drawn from the ADVISORY tier, whose membership is what
	// the implementation extends. Until a fifth member exists, the tier
	// cannot have grown and this step has no subject.
	advisory := graphlint.AdvisoryCodes()
	if len(advisory) <= 4 {
		t.Fatalf("MVV step 2: the advisory tier still carries its original "+
			"four members (%v). The MVV adds a NEW graph-lint finding code "+
			"at severity `info` that fires on the clean model; `0006:C17` is "+
			"amended from \"closed at\" four to append-only (A1) precisely so "+
			"this step has a subject.", advisory)
	}

	step3Out, _, step3Err := runCmd(t, "lint", "--model", clean, "--as=json")

	// --- step 4: the expected end state --------------------------------
	//
	// "the exit code is unchanged (0), the envelope's `type` is still `ok`,
	// and the new finding appears in `data.findings` carrying
	// `\"severity\":\"info\"`" — each asserted BY VALUE, not by absence of
	// error.
	if step3Err != nil {
		t.Fatalf("MVV step 4: the re-run failed (%v); the exit code must be "+
			"UNCHANGED at 0. A growing-vocabulary member is an additive "+
			"change and never moves the verdict.\nstdout:\n%s",
			step3Err, step3Out)
	}

	step3 := rawObject(t, step3Out)
	if got := decodeString(t, step3["type"]); got != "ok" {
		t.Errorf("MVV step 4: type = %q; want %q — a consumer branching on "+
			"`type` and the exit code must observe NO difference", got, "ok")
	}

	var step3Data map[string]json.RawMessage
	if jerr := json.Unmarshal(step3["data"], &step3Data); jerr != nil {
		t.Fatalf("MVV step 4: `data` is not an object: %v\n%s", jerr, step3Out)
	}
	var findings []struct {
		Code     string `json:"code"`
		Severity string `json:"severity"`
	}
	if jerr := json.Unmarshal(step3Data["findings"], &findings); jerr != nil {
		t.Fatalf("MVV step 4: `data.findings` is not an array: %v\n%s",
			jerr, step3Out)
	}
	if len(findings) == 0 {
		t.Fatalf("MVV step 4: the new finding does not appear in "+
			"`data.findings`; a consumer reading findings must see one more "+
			"entry.\nstdout:\n%s", step3Out)
	}
	for _, f := range findings {
		if f.Severity != graphlint.SeverityInfo {
			t.Errorf("MVV step 4: finding %q carries severity %q; want %q — "+
				"the introduced code rides the success payload",
				f.Code, f.Severity, graphlint.SeverityInfo)
		}
		if graphlint.IsBlocking(f.Code) {
			t.Errorf("MVV step 4: finding %q is blocking; introduction MUST "+
				"NOT alter the success disposition of any input", f.Code)
		}
	}

	// --- step 4, the version movement ----------------------------------
	//
	// "The schema version's MAJOR is unchanged — a growing-vocabulary
	// member is an additive change, so C1 moves the minor (`\"0.1\"` →
	// `\"0.2\"`) and never the major." Read OFF THE EMITTED STRING, never
	// inferred from the change class.
	moved, present := schemaVersionOf(t, step3)
	if !present {
		t.Fatalf("MVV step 4: the re-run emitted no `schema_version`")
	}

	baseMajor, baseMinor := splitVersion(t, baseline)
	movedMajor, movedMinor := splitVersion(t, moved)

	if movedMajor != baseMajor {
		t.Errorf("MVV step 4: the MAJOR moved %q→%q. A growing-vocabulary "+
			"member is an ADDITIVE change, so C1 moves the minor and never "+
			"the major — a consumer gating on the major keeps parsing",
			baseMajor, movedMajor)
	}
	if movedMinor == baseMinor {
		t.Errorf("MVV step 4: the minor did not move (%q→%q). The minor "+
			"increments on every change so a consumer can detect movement "+
			"even while it cannot rely on compatibility", baseline, moved)
	}
}

// decodeString decodes one raw JSON string field.
func decodeString(t *testing.T, raw json.RawMessage) string {
	t.Helper()

	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("value %s is not a JSON string: %v", raw, err)
	}
	return s
}

// splitVersion splits a `MAJOR.MINOR` schema version into its components.
func splitVersion(t *testing.T, v string) (major, minor string) {
	t.Helper()

	parts := strings.Split(v, ".")
	if len(parts) != 2 {
		t.Fatalf("schema_version %q is not of the form MAJOR.MINOR", v)
	}
	return parts[0], parts[1]
}
