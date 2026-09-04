package table_test

import (
	"os"
	"strings"
	"testing"
)

// Authoring obligations from `0027`'s implementation plan and prerequisites.
//
// These three REQs constrain the SHAPE OF THE DIFF rather than runtime
// behaviour, so Phase 1 recorded them as coverage orphans. They are not in
// fact unobservable: each names a fact about the shipped tree that a test can
// read directly. Asserting them here closes the orphan rows honestly — the
// test fails if a later edit violates the obligation — rather than pointing a
// row at a behavioural test that does not actually constrain the clause.
//
// REQ-32 stays an orphan by construction: it forbids a test from asserting
// the record's illustrative Go literally, and the compliance evidence is the
// ABSENCE of such a test, which no assertion can express.

// REQ-0b: "0025 is not edited; `load.go`'s C5 comments re-cite this clause."
// DOMAIN EDGE
//
// The record supersedes 0025:C5's argv0 line by writing a successor clause,
// NOT by amending the locked predecessor. A stale C5 comment changes no
// behaviour, which is exactly why it needs an explicit guard: nothing else
// fails when the citation rots.
func TestReq0b_C5CommentsInLoadGoReciteThisRecordsSuccessorClause(t *testing.T) {
	src, err := os.ReadFile("load.go")
	if err != nil {
		t.Fatalf("read load.go: %v", err)
	}
	body := string(src)

	if !strings.Contains(body, "0027:C1") {
		t.Fatal("load.go cites no `0027:C1`: the C5 comments must re-cite this " +
			"record's successor clause, since 0025 itself is never edited (REQ-0b)")
	}

	// The superseded line may still be cited HISTORICALLY ("succeeding
	// 0025:C5's argv0 line"), but no comment may still describe the live
	// predicate as argv0-anchored. Guard the specific superseded phrasing.
	for _, stale := range []string{
		"argv0 + inline-code flag",
		"interpreter at argv0",
		"argv0-anchored",
	} {
		if strings.Contains(body, stale) {
			t.Errorf("load.go still describes the live predicate as %q; C1 is "+
				"position-free and the C5 comments must say so (REQ-0b)", stale)
		}
	}
}

// REQ-35: "Phase 2: Probes — Add the wrapper mutants and the admitted-form
// probes beside the REQ-74 test; keep the four existing probes as-is."
// NEGATIVE REQ / ADVERSARIAL
//
// The widening must not be paid for by loosening the predecessor's probes.
// `TestReq74_KnownInterpreterWithInlineCodeIsALoadTimeDefect` pins four
// vectors that 0025 already refused; all four must survive verbatim, since
// deleting or relaxing one would hide a regression the widening could cause.
func TestReq35_TheFourExistingReq74ProbesSurviveUnweakened(t *testing.T) {
	src, err := os.ReadFile("command_carrier_0025_test.go")
	if err != nil {
		t.Fatalf("read command_carrier_0025_test.go: %v", err)
	}
	body := string(src)

	if !strings.Contains(body, "func TestReq74_KnownInterpreterWithInlineCodeIsALoadTimeDefect(") {
		t.Fatal("the REQ-74 probe test is gone; 0027 keeps the four existing " +
			"probes as-is (REQ-35)")
	}

	// The four probe vectors, by the argv each drives.
	for _, probe := range []struct{ name, argv string }{
		{"sh -c", `command = ["sh", "-c", "cat {artifact}"]`},
		{"bash -c", `command = ["bash", "-c", "echo hi"]`},
		{"python -c", `command = ["python", "-c", "print(1)"]`},
		{"env chain to sh -c", `command = ["env", "sh", "-c", "echo hi"]`},
	} {
		if !strings.Contains(body, probe.argv) {
			t.Errorf("REQ-74 probe %q no longer drives %s; the four existing "+
				"probes are kept as-is (REQ-35)", probe.name, probe.argv)
		}
	}
}

// REQ-40: "intrastate#q2q1's disposition recorded (landed first, or closed as
// subsumed) — either way the `env` walk is removed here."
// BOUNDARY
//
// The prerequisite is discharged by RECORDING the disposition, so the
// artifact is the evidence. The behavioural half — that the `env` walk is
// gone — is covered by Req44; what this asserts is that the record says which
// way q2q1 went, so a later reader is not left guessing.
func TestReq40_TheQ2q1DispositionIsRecordedInDeviations(t *testing.T) {
	const deviations = "../../docs/rdr/0027-inline-shell-detection-scope/artifacts/deviations.md"

	src, err := os.ReadFile(deviations)
	if err != nil {
		t.Fatalf("read %s: %v", deviations, err)
	}
	body := string(src)

	if !strings.Contains(body, "q2q1") {
		t.Fatal("deviations.md records no q2q1 disposition; the prerequisite " +
			"is discharged by recording it (REQ-40)")
	}

	// "landed first, or closed as subsumed" — the record must name one.
	if !strings.Contains(body, "SUBSUMED") && !strings.Contains(body, "subsumed") &&
		!strings.Contains(body, "landed first") {
		t.Error("deviations.md mentions q2q1 but names neither disposition " +
			"(landed first / closed as subsumed) (REQ-40)")
	}
}
