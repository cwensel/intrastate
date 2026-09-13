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
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/graphlint"
	"github.com/cwensel/intrastate/internal/table"
)

// declaredEnumerationSeams is the set of enumeration seams REQ-30 obliges for
// the code vocabularies that cross the wire on `data.findings[].code`.
//
// REQ-30 places a seam "in the package that owns them", so this is a table of
// seams across packages, NOT one package's accessors: `graph-` codes are
// minted in `internal/graphlint`, and `table.AdvisoryNearMiss` is minted in
// `internal/table`, whose members 0006's Naming decision bars from the
// graph-lint tier. A predicate that ranged over `graphlint` alone would be
// blind to every seam the owning-package rule requires, which is the
// over-narrow assertion this table replaces.
//
// Adding an emitted code vocabulary without adding its seam here is the
// defect ADV-2 names; the test below is what makes that omission visible.
func declaredEnumerationSeams() map[string][]string {
	return map[string][]string{
		"graphlint.BlockingCodes": graphlint.BlockingCodes(),
		"graphlint.AdvisoryCodes": graphlint.AdvisoryCodes(),
		"table.AdvisoryRules":     table.AdvisoryRules(),
	}
}

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

	// The subject is the committed snapshot itself
	// (internal/graphlint/testdata/vocabularies.txt, rendered by
	// vocabulary_snapshot_0029_test.go), not any one package's accessors:
	// C3 rests its disclosure obligation on that artifact, so the honest
	// question is whether the artifact carries a severity row for this code.
	// Asserting on the file keeps the check falsifiable — dropping the row
	// turns this red even while every Go seam still enumerates the code.
	snapshot := readVocabularySnapshot(t)

	if !snapshotRecordsSeverity(snapshot, code, severity) {
		t.Errorf("the finding code %q is emitted with \"severity\":%q on the "+
			"success wire, but the C3/A7 committed snapshot records no "+
			"`%s\\t%s` severity row for it.\n\n"+
			"C3 backs its disclosure obligation with that snapshot and "+
			"concedes exactly ONE uncovered row: \"the CLIError `code` row "+
			"has none buildable (A5), so its `append-only` tier stays a "+
			"prose promise and a new code there is the one promotion-shaped "+
			"event this mechanism cannot catch.\" An emitted severity the "+
			"snapshot does not carry is a SECOND uncovered promotion-shaped "+
			"surface, and it is not conceded.\n\n"+
			"Concretely: if the emit site hardcodes a severity literal "+
			"rather than deriving it from a declaration the snapshot "+
			"renders, editing that literal changes what a consumer's CI "+
			"sees while leaving testdata/vocabularies.txt byte-identical — "+
			"the Failure Modes bullet's \"consumer's CI goes red on a model "+
			"nobody edited\", with no artifact naming the moved code.\n\n"+
			"snapshot:\n%s",
			code, severity, code, severity, snapshot)
	}
}

// readVocabularySnapshot returns the committed vocabulary snapshot's bytes.
func readVocabularySnapshot(t *testing.T) string {
	t.Helper()

	path := filepath.Join("..", "graphlint", "testdata", "vocabularies.txt")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("the C3/A7 committed snapshot is unreadable at %s: %v", path, err)
	}
	return string(b)
}

// snapshotRecordsSeverity reports whether the snapshot carries the
// `<code>\t<severity>` mapping row the promotion check diffs.
func snapshotRecordsSeverity(snapshot, code, severity string) bool {
	for line := range strings.SplitSeq(snapshot, "\n") {
		if strings.TrimSpace(line) == code+"\t"+severity {
			return true
		}
	}
	return false
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

	// Every DECLARED enumeration seam, across packages — REQ-30 homes a seam
	// in the package that OWNS the identifier, so a per-package seam must be
	// visible here. If the emitted code belongs to no seam in the table, its
	// tier is unstated and the vocabulary is unenumerable.
	for name, members := range declaredEnumerationSeams() {
		for _, c := range members {
			if c == code {
				t.Logf("code %q is enumerated by the %s seam", code, name)
				return
			}
		}
	}

	t.Errorf("the finding code %q crosses the wire on `data.findings[].code` "+
		"but belongs to no vocabulary carrying an enumeration seam.\n\n"+
		"REQ-29 fixes the scope: \"A vocabulary is tiered wherever it is "+
		"EMITTED, including on the fields of a `findings[]` element … an "+
		"unassigned machine-readable surface is a defect.\" REQ-30 obliges "+
		"the seam: \"an exported accessor returning the vocabulary's "+
		"members, in the package that owns them.\"\n\n"+
		"No seam in declaredEnumerationSeams() enumerates it, so the "+
		"identifier's owning package exports no accessor over the "+
		"vocabulary it belongs to and a second rule of its class can be "+
		"minted with no tier, no seam and no snapshot row. C4 does not list "+
		"this surface among the ones deliberately left untiered (those are "+
		"`findings[].class` and `data.dispositions`).\n\n"+
		"seams consulted: %v",
		code, seamNames())
}

// seamNames lists the declared seams, for the failure message.
func seamNames() []string {
	var names []string
	for name := range declaredEnumerationSeams() {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
