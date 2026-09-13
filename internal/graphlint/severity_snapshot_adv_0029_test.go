package graphlint_test

// RDR 0029 — Phase 3b adversarial: the frozen severity vocabulary is a
// tiered surface the committed snapshot does not record.
//
// Failure Modes, third bullet: "a new machine-readable surface ships
// without a tier assignment, so it has no stated promise and nobody
// notices until a consumer depends on the wrong assumption. Diagnosed by
// auditing C4's list against the emitted vocabularies."
//
// This test IS that audit, for the one C4-frozen vocabulary that has a
// seam and is missing from the artifact.

import (
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/graphlint"
)

// ADV-3 — C4 lists "the severity vocabulary (`blocking`, `info`)" among the
// `frozen` tier's members, and `graphlint.Severities()` is the enumeration
// seam REQ-30 obliges for it. A7's snapshot is the mechanism C3 names for
// catching an unassigned or moved surface: "CI records each one's members
// … in a committed snapshot, and fails when the current tree differs from
// it", reaching "every vocabulary that HAS a seam".
//
// The severity vocabulary HAS a seam and is absent from the rendered
// snapshot. `[graphlint.FindingSeverities]` records each CODE's severity —
// the code→severity mapping — which is a different subject: it would not
// move if a third severity member were added to `severities` and never yet
// attached to a code. A frozen set gaining a member is exactly what the
// snapshot exists to make loud, and for this one vocabulary it is silent.
func TestAdv0029_TheFrozenSeverityVocabularyIsRecordedInTheSnapshot(t *testing.T) {
	rendered := renderVocabularySnapshot()

	if !strings.Contains(rendered, "[graphlint.Severities]") {
		t.Errorf("the committed vocabulary snapshot renders no "+
			"[graphlint.Severities] section, so the frozen severity "+
			"vocabulary's MEMBERS are not recorded anywhere the tree is "+
			"diffed against.\n\n"+
			"C4 tiers \"the severity vocabulary (`blocking`, `info`)\" as "+
			"`frozen`, and graphlint.Severities() is its enumeration seam "+
			"(currently %v). C3 says the check \"reaches every vocabulary "+
			"that HAS a seam\" and concedes exactly one that does not (the "+
			"CLIError `code` registry, A5) — this vocabulary has a seam and "+
			"is still unreached.\n\n"+
			"The existing [graphlint.FindingSeverities] section is not this "+
			"subject: it records the code→severity MAPPING, so adding a "+
			"third member to `severities` — a frozen set gaining a member, "+
			"the loudest event the tier admits — leaves the snapshot "+
			"byte-identical until some code is also moved onto it.",
			graphlint.Severities())
	}
}
