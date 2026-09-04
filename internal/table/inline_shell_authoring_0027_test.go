package table_test

import (
	"os"
	"strings"
	"testing"
)

// Authoring obligations from `0027`'s implementation plan and prerequisites.
//
// These two REQs constrain the SHAPE OF THE DIFF rather than runtime
// behaviour, so Phase 1 recorded them as coverage orphans. They are not in
// fact unobservable: each names a fact about the shipped tree that a test can
// read directly. Asserting them here closes the orphan rows honestly — the
// test fails if a later edit violates the obligation — rather than pointing a
// row at a behavioural test that does not actually constrain the clause.
//
// Each assertion is scoped to the SITE the clause names, not to the whole
// file: a file-wide `strings.Contains` is satisfied by any unrelated
// occurrence elsewhere, so it stays green through the very edit the guard
// exists to catch (deviation D9). Negative checks — "this stale phrasing
// appears NOWHERE" — are correctly file-wide and stay that way.
//
// REQ-32 stays an orphan by construction: it forbids a test from asserting
// the record's illustrative Go literally, and the compliance evidence is the
// ABSENCE of such a test, which no assertion can express.
//
// REQ-35 ("keep the four existing probes as-is") is NOT asserted here. It
// protects the predecessor's REFUSAL BEHAVIOUR, and
// `internal/cli/inline_shell_mvv_0027_test.go`'s
// `ReqMVV0027/step_4_the_0025_probes_still_refuse_and_the_category_order_holds`
// drives the same four argv through `table.Load` asserting
// `CatCommandShellInterpreter`. Re-deriving the property beats reading the
// predecessor test file's bytes, which survives any respelling of the probes
// AND any gutting of the loop that drives them (deviation D9).

// commentBlockContaining returns the contiguous run of `//` comment lines that
// holds anchor, so an assertion can be scoped to ONE comment rather than to
// the whole file. ok is false when no comment line carries the anchor.
func commentBlockContaining(src, anchor string) (block string, ok bool) {
	lines := strings.Split(src, "\n")
	at := -1
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "//") &&
			strings.Contains(line, anchor) {
			at = i
			break
		}
	}
	if at < 0 {
		return "", false
	}
	lo, hi := at, at
	for lo > 0 && strings.HasPrefix(strings.TrimSpace(lines[lo-1]), "//") {
		lo--
	}
	for hi+1 < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[hi+1]), "//") {
		hi++
	}
	return strings.Join(lines[lo:hi+1], "\n"), true
}

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

	// The two comments the clause names, each identified by the SUBSTANTIVE
	// claim it makes rather than by line number or by the citation itself —
	// so rewriting either site's citation back to the predecessor fails here
	// even though `0027:C1` still appears at the file's other cite sites.
	for _, site := range []struct{ what, anchor string }{
		{
			what:   "clause 3's exemption comment",
			anchor: "exemption belongs to the INTERPRETER FORM",
		},
		{
			what:   "clause 4's superseding comment",
			anchor: "argv0 line",
		},
	} {
		block, ok := commentBlockContaining(body, site.anchor)
		if !ok {
			t.Errorf("load.go has no comment carrying %q, so %s cannot be "+
				"checked; the C5 comments must still explain why C1 supersedes "+
				"the predecessor (REQ-0b)", site.anchor, site.what)
			continue
		}
		if !strings.Contains(block, "0027:C1") {
			t.Errorf("%s cites no `0027:C1`:\n%s\nthe C5 comments must re-cite "+
				"this record's successor clause, since 0025 itself is never "+
				"edited (REQ-0b)", site.what, block)
		}
	}

	// The superseded line may still be cited HISTORICALLY ("succeeding
	// 0025:C5's argv0 line"), but no comment may still describe the live
	// predicate as argv0-anchored. Guard the specific superseded phrasing.
	// This one is correctly FILE-WIDE: it forbids a phrase everywhere.
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

// namedAfter reports whether line opens with prefix+name AND the name ends
// there — the next character must not extend it into a longer name. Without
// this, `## D6` also selects `## D60` and `**Disposition` also selects
// `**Dispositions`, so the guard would read a DIFFERENT record's text and pass.
// What follows the name is otherwise unconstrained: the headings here run
// `## D6 — title` and the fields run both `**Disposition.**` and
// `**Disposition (Stage 8).**`, so retitling or restyling must not matter.
func namedAfter(line, prefix, name string) bool {
	if !strings.HasPrefix(line, prefix+name) {
		return false
	}
	rest := line[len(prefix+name):]
	if rest == "" {
		return true
	}
	r := rest[0]
	return !('0' <= r && r <= '9') && !('a' <= r && r <= 'z') && !('A' <= r && r <= 'Z')
}

// markdownSection returns the body of the `## <name>` section — heading line
// through the line before the next `## ` heading, or EOF. ok is false when no
// such heading exists.
func markdownSection(doc, name string) (section string, ok bool) {
	lines := strings.Split(doc, "\n")
	at := -1
	for i, line := range lines {
		if namedAfter(line, "## ", name) {
			at = i
			break
		}
	}
	if at < 0 {
		return "", false
	}
	end := len(lines)
	for i := at + 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "## ") {
			end = i
			break
		}
	}
	return strings.Join(lines[at:end], "\n"), true
}

// boldField returns the `**<name>` paragraph within section: the line opening
// the field through the last line before the next blank line. ok is false when
// the section declares no such field.
func boldField(section, name string) (field string, ok bool) {
	lines := strings.Split(section, "\n")
	at := -1
	for i, line := range lines {
		if namedAfter(line, "**", name) {
			at = i
			break
		}
	}
	if at < 0 {
		return "", false
	}
	end := len(lines)
	for i := at + 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "" {
			end = i
			break
		}
	}
	return strings.Join(lines[at:end], "\n"), true
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

	// Scoped to D6, the deviation that OWNS this prerequisite. File-wide, an
	// unrelated deviation's "landed first" (D3 records 0027 landing before a
	// peer) launders a D6 whose disposition has been gutted to TBD.
	const owner = "D6"
	section, ok := markdownSection(string(src), owner)
	if !ok {
		t.Fatalf("deviations.md has no `## %s` section; the q2q1 prerequisite "+
			"is discharged by recording its disposition there (REQ-40)", owner)
	}

	if !strings.Contains(section, "q2q1") {
		t.Fatalf("deviations.md %s records no q2q1 disposition:\n%s\nthe "+
			"prerequisite is discharged by recording it (REQ-40)", owner, section)
	}

	// "landed first, or closed as subsumed" — the record must name one, and
	// name it in the DISPOSITION FIELD. Scoped tighter than the section
	// because D6 quotes the clause's own "landed first, or closed as
	// subsumed" wording in its Situation, so a section-wide check would read
	// the requirement back to itself and pass over a disposition gutted to
	// TBD.
	disposition, ok := boldField(section, "Disposition")
	if !ok {
		t.Fatalf("deviations.md %s declares no **Disposition** field:\n%s\nthe "+
			"prerequisite is discharged by recording which way q2q1 went "+
			"(REQ-40)", owner, section)
	}
	if !strings.Contains(disposition, "SUBSUMED") && !strings.Contains(disposition, "subsumed") &&
		!strings.Contains(disposition, "landed first") {
		t.Errorf("deviations.md %s names neither disposition (landed first / "+
			"closed as subsumed):\n%s (REQ-40)", owner, disposition)
	}
}
