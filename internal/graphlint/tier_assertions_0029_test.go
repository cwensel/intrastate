package graphlint_test

// RDR 0029 — S8's three replacements.
//
// `TestReq73_TheBlockingCodeSetIsExactlyTheTenNamed`,
// `TestReq74_TheAdvisoryTierIsClosedAtExactlyFourMembers` and
// `TestReq80_UnprovableCoverageCarriesAReasonFromTheClosedSet` each pin a
// non-`frozen` set by exact value, two of them with the cardinality in the
// name. C2 forbids a consumer to assert on cardinality or ordinal position
// over `append-only` and `growing` sets, so those three are the ANTI-pattern
// S7 retires rather than incumbents it tolerates.
//
// The tests below are their replacements: membership and uniqueness over
// the same seams, renamed off the count. They fail while the originals
// stand, because S7's repo-wide negative is false in a tree that still
// carries an exact-set assertion over a tiered set.

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/graphlint"
)

// `repoRoot` is declared in authority_0006_test.go, in this same external
// test package, and is reused here rather than respelled.

// REQ-53: "`internal/graphlint/findings_0006_test.go::TestReq74_TheAdvisoryTierIsClosedAtExactlyFourMembers`
// is REPLACED rather than updated … It becomes a membership-and-uniqueness
// assertion over `AdvisoryCodes()`, renamed off the count."
// REQ-51: "`TestReq73_TheBlockingCodeSetIsExactlyTheTenNamed` and
// `TestReq80_UnprovableCoverageCarriesAReasonFromTheClosedSet` … Both are
// rewritten with `TestReq74` under S8"
// ADVERSARIAL — the retired assertions must be GONE, not re-pinned at a
// higher number: re-pinning four at five reproduces the defect at six.
func TestReq51And53_TheCardinalityPinningAssertionsAreReplacedNotUpdated(t *testing.T) {
	src, err := os.ReadFile(filepath.Join(repoRoot(t),
		"internal", "graphlint", "findings_0006_test.go"))
	if err != nil {
		t.Fatalf("read findings_0006_test.go: %v", err)
	}
	body := string(src)

	for _, retired := range []string{
		"TestReq74_TheAdvisoryTierIsClosedAtExactlyFourMembers",
		"TestReq73_TheBlockingCodeSetIsExactlyTheTenNamed",
		"TestReq80_UnprovableCoverageCarriesAReasonFromTheClosedSet",
	} {
		if strings.Contains(body, retired) {
			t.Errorf("%s still stands. Its name and its exact-value `want` "+
				"literal encode a cardinality C2 forbids asserting over a "+
				"non-frozen set, so it is REPLACED rather than updated — "+
				"re-pinning it at the next number reproduces the defect at "+
				"the one after. Left standing, S7's repo-wide negative is "+
				"false in the tree the day this RDR is declared done", retired)
		}
	}
}

// REQ-26: "`growing`: the graph-lint ADVISORY finding codes. `0006:C17` is
// amended from \"closed at\" its four members to append-only as part of
// this RDR's implementation (A1)"
// BOUNDARY — the amendment is in the peer record's own text.
func TestReq26_The0006AdvisoryTierTextIsAmendedOffItsClosure(t *testing.T) {
	matches, err := filepath.Glob(filepath.Join(repoRoot(t),
		"docs", "rdr", "0006-*.md"))
	if err != nil || len(matches) == 0 {
		t.Fatalf("RDR 0006 is not on disk: %v", err)
	}
	src, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatalf("read 0006: %v", err)
	}

	if strings.Contains(string(src), "The advisory tier is closed at") {
		t.Errorf("`0006:C17` still reads \"The advisory tier is closed at\" " +
			"its four members. A1 rules that closure incidental-as-then-stood " +
			"and opens it to append-only as part of THIS RDR's " +
			"implementation; C17's load-bearing guarantee is its last " +
			"sentence, \"Advisory findings MUST NOT change the success " +
			"disposition\", which adding a member does not touch")
	}
}

// REQ-50: "membership-and-uniqueness assertions only, each over that
// vocabulary's enumeration seam"
// REQ-14: "`growing` — as `append-only`, and additionally a new member MAY
// fire on input that previously produced no such finding"
// HAPPY PATH — the advisory tier's replacement assertion, renamed off the
// count. A fifth member is licensed by the tier and must not fail here.
func TestReq14And50_TheAdvisoryTierCarriesItsNamedMembersAndNoDuplicate(t *testing.T) {
	got := graphlint.AdvisoryCodes()
	if len(got) == 0 {
		t.Fatal("AdvisoryCodes() is empty")
	}

	// MEMBERSHIP: the four members `0006:C17` names must each be present.
	// Their presence is the contract; the set's SIZE is not.
	for _, want := range []string{
		graphlint.CodeCoverageClosedByEscape,
		graphlint.CodeRedundantRow,
		graphlint.CodeUnreachableRule,
		graphlint.CodeVacuousAtom,
	} {
		if !slices.Contains(got, want) {
			t.Errorf("AdvisoryCodes() omits %q", want)
		}
	}

	// UNIQUENESS.
	seen := map[string]bool{}
	for _, c := range got {
		if seen[c] {
			t.Errorf("advisory code %q appears twice; a duplicate makes "+
				"membership ambiguous", c)
		}
		seen[c] = true
	}

	// The tiers stay disjoint, and no advisory code is blocking — that is
	// C17's load-bearing guarantee, which the amendment does not touch.
	for _, c := range got {
		if graphlint.IsBlocking(c) {
			t.Errorf("%q is reported blocking; an advisory finding never "+
				"changes the success disposition", c)
		}
	}
}

// REQ-25: "`append-only`: … the graph-lint BLOCKING finding codes"
// REQ-50: membership-and-uniqueness only.
// HAPPY PATH — TestReq73's replacement, renamed off the count.
func TestReq25And50_TheBlockingTierCarriesItsNamedMembersAndNoDuplicate(t *testing.T) {
	got := graphlint.BlockingCodes()
	if len(got) == 0 {
		t.Fatal("BlockingCodes() is empty")
	}

	for _, want := range []string{
		graphlint.CodeDanglingEdge,
		graphlint.CodeDeadEnd,
		graphlint.CodeOverlap,
		graphlint.CodeCoverageGap,
		graphlint.CodeUnprovableCoverage,
		graphlint.CodeSingleValuedState,
		graphlint.CodeAlwaysPresentOwned,
		graphlint.CodeOwnedBeforeWrite,
		graphlint.CodeTerminalEscape,
		graphlint.CodeProductTooLarge,
	} {
		if !slices.Contains(got, want) {
			t.Errorf("BlockingCodes() omits %q", want)
		}
		if !graphlint.IsBlocking(want) {
			t.Errorf("%q is not reported blocking", want)
		}
	}

	seen := map[string]bool{}
	for _, c := range got {
		if seen[c] {
			t.Errorf("blocking code %q appears twice", c)
		}
		seen[c] = true
	}
}

// REQ-25: "the `graph-unprovable-coverage` `reason` set" is append-only.
// REQ-50: membership-and-uniqueness only.
// HAPPY PATH — TestReq80's replacement. Its predecessor's own comment
// defended an exact-set assertion on a set it called append-only, which is
// the ambiguity C2 retires arguing for itself in a test.
func TestReq25And50_TheReasonSetCarriesItsNamedMembersAndNoDuplicate(t *testing.T) {
	got := graphlint.Reasons()
	if len(got) == 0 {
		t.Fatal("Reasons() is empty")
	}

	for _, want := range []string{
		graphlint.ReasonDimensionNotFinite,
		graphlint.ReasonTagNotSingleValued,
		graphlint.ReasonRowCanRefuse,
		graphlint.ReasonNoParticipatingDimension,
	} {
		if !slices.Contains(got, want) {
			t.Errorf("Reasons() omits %q", want)
		}
	}

	seen := map[string]bool{}
	for _, r := range got {
		if seen[r] {
			t.Errorf("reason %q appears twice", r)
		}
		seen[r] = true
	}

	// The discriminator stays scoped to its one code: every emitted
	// `graph-unprovable-coverage` carries a reason FROM the set, and no
	// other code carries one. This is membership, never the set's size.
	r := lint(t, optionalGuardDecls, multiDefectBody)
	for _, f := range r.Findings {
		if f.Code == graphlint.CodeUnprovableCoverage {
			if f.Reason != "" && !slices.Contains(got, f.Reason) {
				t.Errorf("%s carries reason %q, outside the enumerated set %v",
					f.Code, f.Reason, got)
			}
			continue
		}
		if f.Reason != "" {
			t.Errorf("%s carries the reason discriminator %q; it is scoped "+
				"to %s alone", f.Code, f.Reason,
				graphlint.CodeUnprovableCoverage)
		}
	}
}

// REQ-18: "A lint finding code is introduced at severity `info`."
// REQ-59: "when a new finding could be introduced at either severity,
// `info` is chosen unless C3's re-attribution clause applies"
// DOMAIN EDGE — severity is DERIVED, not carried per-code, so a new code
// that is not in the blocking tier is `info` with no severity plumbing.
func TestReq18And59_ACodeOutsideTheBlockingTierTakesInfoSeverity(t *testing.T) {
	// A code the taxonomy does not know is not blocking: the two tiers are
	// the whole vocabulary and a stranger cannot refuse a model. That is
	// the mechanism by which introduction lands at `info`.
	for _, stranger := range []string{
		"graph-some-new-advisory", "graph-not-yet-minted", "",
	} {
		if graphlint.IsBlocking(stranger) {
			t.Errorf("IsBlocking(%q) = true; a code outside the blocking "+
				"tier must not refuse a model — this is what makes `info` "+
				"the default landing for an introduction", stranger)
		}
	}

	// And the advisory tier, which is where an introduced code lands, is
	// disjoint from the blocking one.
	for _, c := range graphlint.AdvisoryCodes() {
		if graphlint.IsBlocking(c) {
			t.Errorf("%q is in both tiers", c)
		}
	}
}

// REQ-19: "Introduction of an `info` code is a minor-release change and
// MUST NOT alter the success disposition of any input"
// REQ-46: "exit 0, `type` still `\"ok\"`, the finding present in
// `data.findings` with `\"severity\":\"info\"`"
// HAPPY PATH — at the engine boundary; the CLI half is asserted end to end
// in the MVV runner.
func TestReq19And46_AdvisoryFindingsLeaveTheSuccessDispositionUntouched(t *testing.T) {
	r := lint(t, twoStateDecls, advisoryBody)

	advisory := r.Advisory()
	if len(advisory) == 0 {
		t.Fatalf("the advisory fixture produced no advisory findings; "+
			"report:%s", render(r))
	}
	for _, f := range advisory {
		if f.Severity != graphlint.SeverityInfo {
			t.Errorf("%s severity = %q; want %q — an introduced code is "+
				"`info`", f.Code, f.Severity, graphlint.SeverityInfo)
		}
	}

	// The disposition: no blocking finding, so the run succeeds. An
	// advisory member firing must never move that verdict.
	if blocking := r.Blocking(); len(blocking) != 0 {
		var codes []string
		for _, f := range blocking {
			codes = append(codes, f.Code)
		}
		t.Errorf("the advisory fixture reported blocking findings %v; "+
			"advisory findings MUST NOT change the success disposition",
			codes)
	}
}
