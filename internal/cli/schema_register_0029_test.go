package cli

// RDR 0029 — C2's published register and C3's mechanical backing.
//
// C2 requires every machine-readable vocabulary to carry exactly one
// declared tier, recorded in `docs/cli-output-contract.md` beside that
// vocabulary. C3 requires the disclosure obligation not to rest on review
// alone: the tiered-vocabulary seams are enumerable, so CI records each
// one's members and each finding code's severity in a committed snapshot
// and fails when the tree differs from it.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// outputContract reads the published register C2 names.
func outputContract(t *testing.T) string {
	t.Helper()

	body, err := os.ReadFile(filepath.Join(repoRootFor(t),
		"docs", "cli-output-contract.md"))
	if err != nil {
		t.Fatalf("read docs/cli-output-contract.md: %v", err)
	}
	return string(body)
}

// REQ-11: "Every machine-readable vocabulary this CLI emits carries exactly
// one declared stability tier, and the tier is recorded in
// `docs/cli-output-contract.md` beside that vocabulary"
// REQ-12: "`frozen` — no member added or removed within a major."
// REQ-13: "`append-only` — members MAY be added in a minor; none removed or
// renamed within a major."
// REQ-14: "`growing` — as `append-only`, and additionally a new member MAY
// fire on input that previously produced no such finding, subject to C3."
// HAPPY PATH
func TestReq11And12And13And14_TheTierRegisterIsPublishedInTheOutputContract(t *testing.T) {
	doc := outputContract(t)

	// The three tier NAMES must be defined in the register, or a consumer
	// reading a tier word has nowhere to learn what it promises.
	for _, tier := range []string{"frozen", "append-only", "growing"} {
		if !strings.Contains(doc, tier) {
			t.Errorf("docs/cli-output-contract.md defines no `%s` tier. C2 "+
				"requires every emitted vocabulary to carry exactly one "+
				"declared tier RECORDED HERE, beside that vocabulary — the "+
				"tier names are the consumer-facing half of this contract",
				tier)
		}
	}

	// And the vocabularies C4 tiers must appear beside a tier, so the
	// register is a table a consumer can look a vocabulary up in rather
	// than three definitions with nothing assigned.
	for _, vocab := range []string{
		"schema_version", "severity", "verdict", "escape_class",
	} {
		if !strings.Contains(doc, vocab) {
			t.Errorf("docs/cli-output-contract.md does not record a tier "+
				"beside `%s`; C2's register is per-vocabulary", vocab)
		}
	}
}

// REQ-11 (the qualifier C2 attaches to the register): "While the binary's
// version is `0.x`, a tier is a DECLARATION OF INTENT … The register
// (Activation Step 1, item 4) states the same rule in the same section as
// the tier table, so the tier word is never published without the qualifier
// attached to it."
// DOMAIN EDGE — a tier word published bare would be read as a guarantee
// already in force.
func TestReq11_TheRegisterCarriesThe0xIntentQualifierBesideTheTierTable(t *testing.T) {
	doc := outputContract(t)

	if !strings.Contains(doc, "frozen") {
		t.Fatal("no tier table is published in docs/cli-output-contract.md, " +
			"so the `0.x` qualifier has nothing to attach to. C2 requires " +
			"the qualifier in the SAME section as the table precisely so a " +
			"tier word is never published bare — the table and its qualifier " +
			"land together or neither does.")
	}

	lower := strings.ToLower(doc)
	var qualified bool
	for _, phrase := range []string{
		"declaration of intent", "intent", "not a guarantee",
		"already in force", "0.x",
	} {
		if strings.Contains(lower, phrase) {
			qualified = true
			break
		}
	}
	if !qualified {
		t.Error("the tier table is published without the `0.x` qualifier. " +
			"While the binary is 0.x a tier is a DECLARATION OF INTENT and " +
			"MUST NOT be read as a guarantee already in force; C2 requires " +
			"that rule in the SAME section as the table, so the tier word is " +
			"never published bare")
	}
}

// REQ-21: "CI records each one's members and each finding code's severity
// in a committed snapshot, and fails when the current tree differs from it"
// HAPPY PATH — A7 fixes the SHAPE as a golden-file `go test` riding the
// existing `test` job, not a new CI job and not a `go:generate` directive.
func TestReq21_ACommittedSnapshotOfTheSeamMembersAndSeveritiesIsDiffed(t *testing.T) {
	root := repoRootFor(t)

	// The committed artifact. A7 rules the idiom a `testdata/` golden,
	// which the repo already carries at internal/table/testdata.
	var found string
	for _, candidate := range []string{
		filepath.Join("internal", "graphlint", "testdata", "vocabularies.txt"),
		filepath.Join("internal", "graphlint", "testdata", "vocabulary-snapshot.txt"),
		filepath.Join("internal", "graphlint", "testdata", "seams.txt"),
		filepath.Join("internal", "cli", "testdata", "vocabularies.txt"),
		filepath.Join("internal", "cli", "testdata", "vocabulary-snapshot.txt"),
	} {
		if _, err := os.Stat(filepath.Join(root, candidate)); err == nil {
			found = candidate
			break
		}
	}
	if found == "" {
		t.Fatal("no committed vocabulary snapshot exists. C3 refuses to rest " +
			"the disclosure obligation on review alone: the tiered-vocabulary " +
			"seams C4 obliges are enumerable, so CI records each one's " +
			"members and each finding code's SEVERITY in a committed " +
			"snapshot and fails when the current tree differs. A7 verified " +
			"the shape — a golden-file `go test` riding the existing `test` " +
			"job — by demonstrating the undisclosed promotion of " +
			"`graph-vacuous-atom` red and then green again.")
	}

	body, err := os.ReadFile(filepath.Join(root, found))
	if err != nil {
		t.Fatalf("read %s: %v", found, err)
	}
	snapshot := string(body)

	// Its subject is the SEAM MEMBERS and their severities. A snapshot that
	// recorded only members could not catch a promotion, which is the one
	// event C3 exists to make loud.
	for _, want := range []string{"blocking", "info"} {
		if !strings.Contains(snapshot, want) {
			t.Errorf("the snapshot records no %q severity; a promotion "+
				"`info`→`blocking` is exactly what it must turn red", want)
		}
	}
	// And a promotion-detecting snapshot must name the codes whose severity
	// can move.
	if !strings.Contains(snapshot, "graph-vacuous-atom") {
		t.Errorf("the snapshot does not record `graph-vacuous-atom`'s " +
			"severity; A7's demonstrated failing case is moving exactly that " +
			"code between tiers")
	}
}

// REQ-21 (the diff half): the snapshot is only a check if something
// compares the tree against it.
// ADVERSARIAL — a committed file nothing reads is a decoration, and the
// obligation would silently be back to review-only.
func TestReq21_TheSnapshotIsComparedAgainstTheCurrentTree(t *testing.T) {
	root := repoRootFor(t)

	// The artifact must exist before asking whether anything reads it;
	// otherwise an unrelated test mentioning `testdata` satisfies the
	// substring probe and this passes with no snapshot in the tree.
	var snapshotExists bool
	for _, candidate := range []string{
		filepath.Join("internal", "graphlint", "testdata", "vocabularies.txt"),
		filepath.Join("internal", "graphlint", "testdata", "vocabulary-snapshot.txt"),
		filepath.Join("internal", "graphlint", "testdata", "seams.txt"),
		filepath.Join("internal", "cli", "testdata", "vocabularies.txt"),
		filepath.Join("internal", "cli", "testdata", "vocabulary-snapshot.txt"),
	} {
		if _, err := os.Stat(filepath.Join(root, candidate)); err == nil {
			snapshotExists = true
			break
		}
	}
	if !snapshotExists {
		t.Fatal("no committed vocabulary snapshot exists, so \"something " +
			"compares the tree against it\" has no subject")
	}

	// Some test must READ the snapshot and compare. A7 hosts it beside
	// internal/graphlint in an external test package; the search is over
	// the tree rather than one path so the host can move.
	var reader bool
	for _, dir := range []string{
		filepath.Join("internal", "graphlint"),
		filepath.Join("internal", "cli"),
	} {
		entries, err := os.ReadDir(filepath.Join(root, dir))
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !strings.HasSuffix(e.Name(), "_test.go") {
				continue
			}
			body, rerr := os.ReadFile(filepath.Join(root, dir, e.Name()))
			if rerr != nil {
				continue
			}
			if strings.Contains(string(body), "testdata") &&
				(strings.Contains(string(body), "vocabular") ||
					strings.Contains(string(body), "seams")) {
				reader = true
			}
		}
	}
	if !reader {
		t.Error("no test reads the vocabulary snapshot, so nothing fails " +
			"when the tree differs from it. A diff must be a DELIBERATE act " +
			"— the PR either updates the snapshot, which is the reviewable " +
			"artifact naming the moved code, or goes red.")
	}
}

// REQ-22: "A code MAY be introduced directly at `blocking` only when it
// reports a condition that was already refused by some other code"
// REQ-9: "the minor component still increments on every change so a
// consumer can detect movement even while it cannot rely on compatibility"
// DOMAIN EDGE — the register is where both rules are published, since
// neither has a runtime witness in this build.
func TestReq9And22_TheIntroductionAndMovementRulesArePublished(t *testing.T) {
	doc := outputContract(t)
	lower := strings.ToLower(doc)

	if !strings.Contains(lower, "schema_version") {
		t.Fatal("the output contract does not document `schema_version`; " +
			"C1's version marker is the field a consumer gates on")
	}

	// REQ-9: movement is detectable even while the major is 0.
	var movement bool
	for _, phrase := range []string{
		"minor", "increment", "every change", "detect movement",
	} {
		if strings.Contains(lower, phrase) {
			movement = true
			break
		}
	}
	if !movement {
		t.Error("the register does not state that the minor increments on " +
			"every change. While the major is `0` the schema may change " +
			"incompatibly, so the moving minor is the ONLY signal a consumer " +
			"has that the wire moved at all")
	}

	// REQ-22: the one escape hatch, and the evidence that discharges it.
	var reattribution bool
	for _, phrase := range []string{
		"already refused", "re-attribut", "reattribut",
	} {
		if strings.Contains(lower, phrase) {
			reattribution = true
			break
		}
	}
	if !reattribution {
		t.Error("the register does not publish C3's direct-at-`blocking` " +
			"exception. A code may be introduced at `blocking` ONLY when it " +
			"re-attributes an existing refusal — same input, same verdict, " +
			"different code — and the exception is discharged by evidence, " +
			"not by judgement, which keeps it from widening into \"we " +
			"thought it was serious enough\"")
	}
}
