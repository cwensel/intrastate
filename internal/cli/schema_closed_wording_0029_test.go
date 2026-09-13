package cli

// RDR 0029 — C2's last clause: the retired `closed` wording.
//
// The predicate is VOCABULARY-scoped, not package-scoped. `closed` stays
// wherever it describes something C4 does not tier (the accessor capability
// set, the predicate-semantic-kind set, `DumpColumns`, the `cmdbind`
// placeholder vocabulary), and everywhere it means something else entirely
// (`fail closed`, closed-world matching, a closed pipe, an unclosed brace).
// C2 retires an ambiguity between two tier senses; it does not reserve an
// English word.
//
// S9 therefore owes an ENUMERATED census rather than a grep for the word:
// "A grep over that enumerated site list is the assertion; an unscoped grep
// for the word is not, and would fail on dozens of sites this RDR does not
// intend to touch."

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// s9Census is the in-scope site list S9 enumerates, each paired with the
// tiered vocabulary whose description made it in-scope. Nothing outside
// this list is asserted, which is the point: the census is owed with the
// fix rather than discovered by grep.
var s9Census = []struct {
	path       string
	vocabulary string
}{
	{filepath.Join("internal", "guard", "grammar.go"), "Operators (frozen)"},
	{filepath.Join("internal", "guard", "doc.go"), "Operators (frozen)"},
	{filepath.Join("internal", "table", "model.go"), "the Operators mirror (frozen)"},
	{filepath.Join("internal", "table", "normalize.go"), "the Operators mirror (frozen)"},
	{filepath.Join("internal", "accessor", "model.go"), "the verdict set (frozen)"},
	{filepath.Join("internal", "resolve", "resolve.go"), "RefusalKinds (frozen)"},
	{filepath.Join("internal", "graphlint", "taxonomy.go"), "the advisory and reason sets"},
	{filepath.Join("internal", "table", "category.go"), "Categories (append-only)"},
	{filepath.Join("internal", "cli", "flow_input.go"), "the CLIError code vocabulary (append-only)"},
}

// exemptLine reports whether a line carrying `closed` is one C2 exempts:
// the emitted identifier and the Go identifiers that produce it, and
// verbatim quotations of a peer contract that predates this RDR.
func exemptLine(line string) bool {
	for _, exempt := range []string{
		// C2: "It does not reach an emitted identifier that merely contains
		// the word" — and the Go identifiers that produce it.
		"graph-coverage-closed-by-escape",
		"ClosedByEscape",
		"closedBy",
		// C2: "Nor does it reach a verbatim quotation of a peer contract
		// that predates this RDR".
		"0003:C7",
	} {
		if strings.Contains(line, exempt) {
			return true
		}
	}
	// Non-tier senses of the word: a closed pipe, closed-world matching,
	// fail-closed, an unclosed brace. C2 retires an ambiguity between two
	// TIER senses; it does not reserve an English word.
	for _, sense := range []string{
		"closed-world", "closed world", "fail closed", "fail-closed",
		"unclosed", "closed pipe", "closed the pipe", "closed layout",
	} {
		if strings.Contains(strings.ToLower(line), sense) {
			return true
		}
	}
	return false
}

// REQ-54: "no occurrence of `closed` describing a vocabulary C4 TIERS"
// REQ-55: the enumerated in-scope site census.
// REQ-57: "A grep over that enumerated site list is the assertion; an
// unscoped grep for the word is not"
// REQ-15: "The term `closed` MUST NOT be used to DESCRIBE any of these
// tiers, in code comments or documentation"
// BOUNDARY
func TestReq15And54And55And57_TheClosedWordingIsRetiredAcrossTheEnumeratedCensus(t *testing.T) {
	root := repoRootFor(t)

	for _, site := range s9Census {
		t.Run(site.path, func(t *testing.T) {
			body, err := os.ReadFile(filepath.Join(root, site.path))
			if err != nil {
				t.Fatalf("read %s: %v", site.path, err)
			}

			for i, line := range strings.Split(string(body), "\n") {
				if !strings.Contains(strings.ToLower(line), "closed") {
					continue
				}
				if exemptLine(line) {
					continue
				}
				t.Errorf("%s:%d describes %s with the retired word:\n  %s\n"+
					"C2 retires `closed` as a TIER description because it is "+
					"live in two incompatible senses; the declared tier name "+
					"(`frozen`, `append-only`, `growing`) says which promise "+
					"is meant.", site.path, i+1, site.vocabulary,
					strings.TrimSpace(line))
			}
		})
	}
}

// REQ-16: "It does not reach an emitted identifier that merely contains the
// word: `graph-coverage-closed-by-escape` is a member of a `growing`
// vocabulary"
// REQ-56: "Exempt, per C2: the emitted identifier
// `graph-coverage-closed-by-escape` together with the Go identifiers that
// produce it (`ClosedByEscape`, `closedBy`), and verbatim quotations of
// `0003:C7`."
// ADVERSARIAL — renaming the emitted identifier would be the breaking
// change this contract exists to prevent.
func TestReq16And56_TheEmittedIdentifierCarryingTheWordSurvives(t *testing.T) {
	body, err := os.ReadFile(filepath.Join(repoRootFor(t),
		"internal", "graphlint", "taxonomy.go"))
	if err != nil {
		t.Fatalf("read taxonomy.go: %v", err)
	}

	// The wire string is a member of a `growing` vocabulary. A retirement
	// pass that renamed it would break every consumer branching on it.
	if !strings.Contains(string(body), `"graph-coverage-closed-by-escape"`) {
		t.Error("the emitted identifier `graph-coverage-closed-by-escape` is " +
			"gone. C2's prohibition does not reach an emitted identifier " +
			"that merely CONTAINS the word — it is a member of a `growing` " +
			"vocabulary, so renaming it is the breaking change this contract " +
			"exists to prevent")
	}
	// And the Go identifier that produces it.
	if !strings.Contains(string(body), "CodeCoverageClosedByEscape") {
		t.Error("the Go identifier `CodeCoverageClosedByEscape` is gone; C2 " +
			"exempts the identifiers that produce the emitted member")
	}
}

// REQ-17: "Nor does it reach a verbatim quotation of a peer contract that
// predates this RDR (C4 quotes `0003:C7`)"
// REQ-63: "Sites describing an untiered vocabulary, and the non-tier senses
// of the word, are out of scope and stay; so do C2's exemptions."
// DOMAIN EDGE — the negative: an over-eager retirement that stripped the
// word from an UNTIERED vocabulary has over-reached.
func TestReq17And63_UntieredVocabulariesAndNonTierSensesKeepTheWord(t *testing.T) {
	root := repoRootFor(t)

	// C4 tiers none of these, and S9 names them as sites where `closed`
	// stays: the accessor capability set, the validation-code set, and the
	// dump column vocabulary.
	for _, tc := range []struct {
		path    string
		snippet string
		why     string
	}{
		{filepath.Join("internal", "accessor", "model.go"),
			"closed three-member capability vocabulary",
			"the accessor capability set is untiered by C4"},
		{filepath.Join("internal", "accessor", "model.go"),
			"closed eight-member validation-code set",
			"the validation-code set is untiered by C4"},
		{filepath.Join("internal", "table", "dump.go"),
			"closed column vocabulary",
			"DumpColumns is untiered by C4"},
	} {
		body, err := os.ReadFile(filepath.Join(root, tc.path))
		if err != nil {
			t.Fatalf("read %s: %v", tc.path, err)
		}
		if !strings.Contains(string(body), tc.snippet) {
			t.Errorf("%s no longer carries %q. %s, so the word stays: the "+
				"prohibition is VOCABULARY-scoped, not package-scoped, and a "+
				"retirement that reached here has over-applied C2",
				tc.path, tc.snippet, tc.why)
		}
	}
}

// REQ-37: "Three doc comments in `respond.go` and `readPlan` still describe
// a `{\"type\":\"failed\",…}` record; they are stale today and Step 2
// retires them"
// REQ-62: "Also correct the two doc comments that describe the failure
// envelope as `{\"type\":\"failed\",…}` — `internal/cli/respond`'s package
// doc and `readPlan`'s comment"
// ADVERSARIAL — the comments contradict what `Fail` actually emits, and
// would make S2 read as a regression.
func TestReq37And62_TheStaleFailedRecordDocCommentsAreRetired(t *testing.T) {
	root := repoRootFor(t)

	for _, tc := range []struct {
		path string
		what string
	}{
		{filepath.Join("internal", "cli", "respond", "respond.go"),
			"respond's package doc and Success's comment"},
		{filepath.Join("internal", "cli", "flow_input.go"),
			"readPlan's comment"},
	} {
		body, err := os.ReadFile(filepath.Join(root, tc.path))
		if err != nil {
			t.Fatalf("read %s: %v", tc.path, err)
		}

		for i, line := range strings.Split(string(body), "\n") {
			trimmed := strings.TrimSpace(line)
			if !strings.HasPrefix(trimmed, "//") {
				continue
			}
			if !strings.Contains(line, `"failed"`) &&
				!strings.Contains(line, `failed",`) {
				continue
			}
			t.Errorf("%s:%d (%s) still describes a `{\"type\":\"failed\",…}` "+
				"record:\n  %s\nA refusal carries no `type` and no wrapper — "+
				"`respond.Fail` writes the bare *CLIError (`0005:C1`). The "+
				"comment is stale today, and left standing it makes S2's "+
				"assertion read as a regression rather than a confirmation.",
				tc.path, i+1, tc.what, trimmed)
		}
	}
}
