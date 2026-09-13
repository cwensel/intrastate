package cli

// RDR 0029 — Phase 3b adversarial: the Failure Modes section's two
// SILENT modes, read against the near-miss advisory.
//
// Failure Modes names "an `info` code is promoted to `blocking` without a
// release note … the only way to diagnose it is to diff the severity
// constants between two builds" and "a new machine-readable surface ships
// without a tier assignment, so it has no stated promise and nobody
// notices until a consumer depends on the wrong assumption."
//
// C3 backs the first with the committed snapshot and concedes exactly ONE
// blind spot: "The check reaches every vocabulary that HAS a seam; the
// CLIError `code` row has none buildable (A5), so its `append-only` tier
// stays a prose promise and a new code there is the one promotion-shaped
// event this mechanism cannot catch." One — not two.
//
// `internal/cli/lint.go::nearMissFindings` emits `table.AdvisoryNearMiss`
// as a `findings[].code` carrying `"severity":"info"` on the SUCCESS wire.
// That code is in neither `graphlint.BlockingCodes()` nor
// `AdvisoryCodes()`, so the snapshot — which ranges over exactly those two
// — never records it, and the severity is hardcoded at the emit site
// rather than derived from a taxonomy the snapshot can diff.

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/graphlint"
	"github.com/cwensel/intrastate/internal/table"
)

// nearMissModel builds the fixture that lints CLEAN and still raises the
// near-miss advisory: `Recognized` is a legal, unreserved name that becomes
// the reserved key under simple case folding alone.
func nearMissModel(t *testing.T) string {
	t.Helper()

	src := strings.Replace(legalModel,
		"[tags.flag]\nprovenance = \"owned\"\nkind = \"bool\"",
		"[tags.Recognized]\nprovenance = \"owned\"\nkind = \"bool\"",
		1)
	if src == legalModel {
		t.Fatalf("fixture edit did not apply; legalModel changed shape")
	}
	src = strings.ReplaceAll(src, `keys = ["flag", "status"]`,
		`keys = ["Recognized", "status"]`)
	src = strings.ReplaceAll(src, "[rule.guard.all.flag]",
		"[rule.guard.all.Recognized]")
	src = strings.ReplaceAll(src, "flag = \"false\"",
		"Recognized = \"false\"")
	return src
}

// emittedNearMissFinding runs the lint verb on the near-miss fixture and
// returns the emitted finding whose code is the advisory identifier.
func emittedNearMissFinding(t *testing.T) map[string]any {
	t.Helper()

	path := writeModel(t, nearMissModel(t))

	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err != nil {
		t.Fatalf("the near-miss model did not lint clean: %v\nstdout:\n%s",
			err, stdout)
	}

	var env struct {
		Data struct {
			Findings []map[string]any `json:"findings"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &env); err != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", err, stdout)
	}

	for _, f := range env.Data.Findings {
		if code, _ := f["code"].(string); code == table.AdvisoryNearMiss {
			return f
		}
	}
	t.Fatalf("the near-miss fixture emitted no %q finding; the fixture is "+
		"not exercising the advisory\nfindings: %+v",
		table.AdvisoryNearMiss, env.Data.Findings)
	return nil
}

// ADV-1 — FAILURE MODE: "an `info` code is promoted to `blocking` without a
// release note … the only way to diagnose it is to diff the severity
// constants between two builds" (Failure Modes, second bullet).
//
// C3 rests that obligation on the committed snapshot and concedes ONE
// uncovered row (the CLIError `code` registry, A5). The near-miss advisory
// is a SECOND, unconceded one: it emits a severity on the wire that the
// snapshot's range — `BlockingCodes() + AdvisoryCodes()` — does not
// include, so promoting it to `blocking` moves no byte of
// `testdata/vocabularies.txt` and merges green.
//
// The test catches the mode by requiring that any code emitted with a
// severity on the wire is a code the snapshot's range records.
func TestAdv0029_EveryWireSeverityIsCarriedByTheSnapshotRange(t *testing.T) {
	finding := emittedNearMissFinding(t)

	code, _ := finding["code"].(string)
	severity, _ := finding["severity"].(string)

	if severity == "" {
		t.Fatalf("the near-miss finding %q carries no `severity` on the "+
			"wire; this test's premise is that it does", code)
	}

	// The snapshot (internal/graphlint/testdata/vocabularies.txt, rendered
	// by vocabulary_snapshot_0029_test.go) records severities for exactly
	// the union of these two seams.
	covered := false
	for _, c := range append(graphlint.BlockingCodes(),
		graphlint.AdvisoryCodes()...) {
		if c == code {
			covered = true
			break
		}
	}

	if !covered {
		t.Errorf("the finding code %q is emitted with \"severity\":%q on the "+
			"success wire, but it is a member of neither "+
			"graphlint.BlockingCodes() nor graphlint.AdvisoryCodes() — the "+
			"two seams the C3/A7 committed snapshot ranges over.\n\n"+
			"C3 backs its disclosure obligation with that snapshot and "+
			"concedes exactly ONE uncovered row: \"the CLIError `code` row "+
			"has none buildable (A5), so its `append-only` tier stays a "+
			"prose promise and a new code there is the one promotion-shaped "+
			"event this mechanism cannot catch.\" This is a SECOND "+
			"uncovered promotion-shaped surface, and it is not conceded.\n\n"+
			"Concretely: `internal/cli/lint.go::nearMissFindings` hardcodes "+
			"`Severity: graphlint.SeverityInfo` at the emit site rather than "+
			"deriving it from a taxonomy the snapshot diffs. Editing that "+
			"literal to SeverityBlocking changes what a consumer's CI sees "+
			"while leaving testdata/vocabularies.txt byte-identical — the "+
			"Failure Modes bullet's \"consumer's CI goes red on a model "+
			"nobody edited\", with no artifact naming the moved code.",
			code, severity)
	}
}

// ADV-2 — FAILURE MODE: "a new machine-readable surface ships without a
// tier assignment, so it has no stated promise and nobody notices until a
// consumer depends on the wrong assumption" (Failure Modes, third bullet).
//
// REQ-29 (C4): "A vocabulary is tiered wherever it is EMITTED, including on
// the fields of a `findings[]` element — not only at the envelope's top
// level and under `data` … an unassigned machine-readable surface is a
// defect." REQ-30: "A tier assignment obliges an ENUMERATION SEAM: an
// exported accessor returning the vocabulary's members, in the package that
// owns them."
//
// `table.AdvisoryNearMiss` is emitted as `findings[].code`. C4's census
// enumerates the graph-lint blocking codes, the advisory codes and the
// CLIError `code` vocabulary — the near-miss advisory identifier is none of
// those three: it is minted in `internal/table`, rides the SUCCESS payload,
// and `internal/table` exports no accessor over it.
//
// The test catches the mode by requiring an enumeration seam in the owning
// package for the vocabulary this code belongs to.
func TestAdv0029_TheEmittedAdvisoryCodeVocabularyHasATierSeam(t *testing.T) {
	finding := emittedNearMissFinding(t)
	code, _ := finding["code"].(string)

	// The three enumerable code vocabularies C4's census names. If the
	// emitted code belongs to none of them, no seam in any package covers
	// it and its tier is unstated.
	for _, c := range append(graphlint.BlockingCodes(),
		graphlint.AdvisoryCodes()...) {
		if c == code {
			return // covered by a censused seam
		}
	}

	t.Errorf("the finding code %q crosses the wire on `data.findings[].code` "+
		"but belongs to no vocabulary carrying an enumeration seam.\n\n"+
		"REQ-29 fixes the scope: \"A vocabulary is tiered wherever it is "+
		"EMITTED, including on the fields of a `findings[]` element … an "+
		"unassigned machine-readable surface is a defect.\" REQ-30 obliges "+
		"the seam: \"an exported accessor returning the vocabulary's "+
		"members, in the package that owns them.\"\n\n"+
		"`internal/table` owns this identifier (table.AdvisoryNearMiss, a "+
		"single exported const) and exports no accessor enumerating the "+
		"advisory-rule vocabulary it belongs to, so a second near-miss-class "+
		"rule can be minted with no tier, no seam and no snapshot row. C4's "+
		"census covers the graph-lint blocking and advisory code sets and "+
		"the CLIError `code` registry; this code is in none of them, and C4 "+
		"does not list it among the surfaces deliberately left untiered "+
		"(those are `findings[].class` and `data.dispositions`).",
		code)
}
