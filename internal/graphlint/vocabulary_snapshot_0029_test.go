package graphlint_test

// RDR 0029 — C3's mechanical backing, in the shape A7 verified.
//
// C3 refuses to rest the disclosure obligation on review alone: the
// tiered-vocabulary seams C4 obliges are enumerable, so each one's members
// and each finding code's SEVERITY are recorded in a committed snapshot and
// the tree is diffed against it. A severity moving `info`→`blocking` cannot
// merge undisclosed — the PR either updates the snapshot, which is the
// reviewable artifact naming the moved code, or goes red.
//
// The shape is A7's: a golden-file `go test` riding the existing `test`
// job, hosted in an external test package so the downward-only imports
// (`accessor`, `resolve`, `table`, `guard`, `respond`, `clierr`) raise no
// cycle.

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/cli/respond"
	"github.com/cwensel/intrastate/internal/graphlint"
	"github.com/cwensel/intrastate/internal/guard"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// snapshotPath0029 is the committed artifact the tree is compared against.
const snapshotPath0029 = "testdata/vocabularies.txt"

// renderVocabularySnapshot serializes every tiered vocabulary's members and
// every finding code's severity, in a stable order.
//
// Members are sorted rather than emitted in declaration order: the snapshot
// must go red on a member ADDED, REMOVED or RENAMED, never merely
// reordered, because `append-only` explicitly licenses a consumer nothing
// about ordinal position.
func renderVocabularySnapshot() string {
	var b strings.Builder

	b.WriteString("# Vocabulary and severity snapshot (RDR 0029 C3/A7).\n")
	b.WriteString("# Regenerate deliberately; a diff names the moved member.\n\n")

	section := func(name string, members []string) {
		fmt.Fprintf(&b, "[%s]\n", name)
		sorted := append([]string(nil), members...)
		sort.Strings(sorted)
		for _, m := range sorted {
			fmt.Fprintln(&b, m)
		}
		b.WriteString("\n")
	}

	section("respond.Types", respond.Types())
	section("respond.Levels", respond.Levels())

	var exits []string
	for _, c := range clierr.ExitCodes() {
		exits = append(exits, fmt.Sprint(c))
	}
	section("clierr.ExitCodes", exits)

	var verdicts []string
	for _, v := range accessor.Verdicts() {
		verdicts = append(verdicts, string(v))
	}
	section("accessor.Verdicts", verdicts)

	var kinds []string
	for _, k := range resolve.RefusalKinds() {
		kinds = append(kinds, string(k))
	}
	section("resolve.RefusalKinds", kinds)

	var blocks []string
	for _, bl := range resolve.Blocks() {
		blocks = append(blocks, string(bl))
	}
	section("resolve.Blocks", blocks)

	section("guard.Operators", guard.Operators())
	section("table.Operators", table.Operators())

	var categories []string
	for _, c := range table.Categories() {
		categories = append(categories, string(c))
	}
	section("table.Categories", categories)

	section("graphlint.Reasons", graphlint.Reasons())

	// The severities are the half a member list alone cannot carry: a
	// promotion moves no member, only the tier a code sits in, so the
	// snapshot records the code→severity mapping explicitly.
	b.WriteString("[graphlint.FindingSeverities]\n")
	var codes []string
	codes = append(codes, graphlint.BlockingCodes()...)
	codes = append(codes, graphlint.AdvisoryCodes()...)
	sort.Strings(codes)
	for _, c := range codes {
		severity := graphlint.SeverityInfo
		if graphlint.IsBlocking(c) {
			severity = graphlint.SeverityBlocking
		}
		fmt.Fprintf(&b, "%s\t%s\n", c, severity)
	}

	return b.String()
}

// REQ-21: "CI records each one's members and each finding code's severity
// in a committed snapshot, and fails when the current tree differs from it"
// HAPPY PATH — the diff itself.
func TestReq21_TheVocabularySnapshotMatchesTheTree(t *testing.T) {
	want, err := os.ReadFile(filepath.FromSlash(snapshotPath0029))
	if err != nil {
		t.Fatalf("the committed vocabulary snapshot is missing (%s): %v\n\n"+
			"C3 backs the disclosure obligation with this artifact: without "+
			"it an undisclosed `info`→`blocking` promotion merges silently, "+
			"which is the one event the clause exists to make loud.",
			snapshotPath0029, err)
	}

	got := renderVocabularySnapshot()
	if got != string(want) {
		t.Errorf("the tiered vocabularies differ from the committed "+
			"snapshot.\n\n--- committed (%s)\n%s\n--- current tree\n%s\n"+
			"A diff here is a DELIBERATE act: either the change discloses a "+
			"moved member or severity and updates this file — the reviewable "+
			"artifact naming what moved — or it goes red. A severity moving "+
			"`info`→`blocking` is a promotion C3 requires be disclosed in the "+
			"release notes naming the code.",
			snapshotPath0029, string(want), got)
	}
}

// REQ-21 (the severity half): the snapshot's subject includes each finding
// code's severity, so a PROMOTION — which moves no member between
// vocabularies, only a code between tiers — is caught.
// ADVERSARIAL — A7's demonstrated failing case is moving `graph-vacuous-atom`
// between tiers.
func TestReq21_TheSnapshotRecordsEverySeveritySoAPromotionIsCaught(t *testing.T) {
	rendered := renderVocabularySnapshot()

	for _, code := range append(graphlint.BlockingCodes(),
		graphlint.AdvisoryCodes()...) {
		severity := graphlint.SeverityInfo
		if graphlint.IsBlocking(code) {
			severity = graphlint.SeverityBlocking
		}
		if !strings.Contains(rendered, code+"\t"+severity) {
			t.Errorf("the snapshot does not record %q at severity %q; a "+
				"promotion is invisible to a snapshot that records only "+
				"membership, because the code is a member of the taxonomy "+
				"either way", code, severity)
		}
	}
}
