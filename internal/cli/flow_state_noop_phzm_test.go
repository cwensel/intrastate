package cli

// kata phzm — `flow set-state`'s help states how a no-op is told from an
// applied write.
//
// The kata was filed as a false-green bug: a `stopped:*` plan piped through
// `--plan` was said to be indistinguishable from an applied write at exit 0.
// It is not — `writers` is appended to only as each writer's read-back
// confirms, so it is `[]` exactly when nothing was applied, and the no-op
// shape itself is pinned by
// `TestC3xz_AnEmptyPlanIsANoOpSuccessAtExitZero`. What was genuinely
// missing is that NOTHING SAID SO: the consumer that filed this reached for
// `dispositions`, found it absent, and read the seam as broken.
//
// `dispositions` stays absent by decision, not oversight: it is joined from
// the SELECTED ROW's authored emit (`0024:C4` — cited here in source, never
// in the help itself), and `set-state` selects no row, so echoing a caller-supplied blob under that name would give one wire
// key two provenances — the same reason `flow next` carries none.
//
// This oracle fails if either half of that statement is dropped from the
// help, which is the artifact the next consumer reads.

import (
	"strings"
	"testing"
)

// The help must name the emptiness check AND say why `dispositions` is
// absent. Naming only the first leaves the next consumer to conclude the
// absence is a defect and re-file; naming only the second leaves them with
// no check to use instead.
func TestPhzm_SetStateHelpStatesHowANoOpIsToldFromAnAppliedWrite(t *testing.T) {
	cmd := NewRootCmd()
	sub, _, err := cmd.Find([]string{"flow", "set-state"})
	if err != nil || sub == nil || sub.Name() != "set-state" {
		t.Fatalf("the `flow` group registers no `set-state` verb: %v", err)
	}

	// Both widths matter: `Long` is what a bare `--help` shows, and the
	// extended description is where the `--plan` idiom that raises the
	// question is documented.
	for _, tc := range []struct {
		width string
		help  string
	}{
		{"long", sub.Long},
		{"extended", flowSetStateExtendedDesc},
	} {
		t.Run(tc.width, func(t *testing.T) {
			lower := strings.ToLower(tc.help)

			// (i) the check the caller actually runs.
			if !strings.Contains(lower, "writers") {
				t.Errorf("the %s help never names `writers`; it is the field "+
					"an applied write is told from a no-op by, and the "+
					"caller has no other sanctioned check:\n%s",
					tc.width, tc.help)
			}
			var stated bool
			for _, phrase := range []string{
				"emptiness", "empty", "[]",
			} {
				if strings.Contains(lower, phrase) {
					stated = true
					break
				}
			}
			if !stated {
				t.Errorf("the %s help names `writers` without saying that "+
					"EMPTINESS is the signal; a caller cannot act on the "+
					"field without knowing which value means nothing "+
					"landed:\n%s", tc.width, tc.help)
			}

			// (ii) why `dispositions` is not there, cited to the contract
			// rather than to this kata — the absence is a decision, and an
			// unexplained absence is what got re-filed.
			if !strings.Contains(lower, "dispositions") {
				t.Errorf("the %s help never mentions `dispositions`; its "+
					"absence from this payload is deliberate (`0024:C4`), "+
					"and an unexplained absence reads as a defect:\n%s",
					tc.width, tc.help)
			}
			// The REASON, not an RDR id: internal contract ids are never
			// surfaced in user-facing help, so what makes the absence
			// checkable here is the stated selected-row provenance.
			if !strings.Contains(lower, "selects no row") &&
				!strings.Contains(lower, "selected row") {
				t.Errorf("the %s help asserts `dispositions` is absent "+
					"without stating the REASON — a disposition joins the "+
					"SELECTED ROW's authored emit, and `set-state` selects "+
					"no row:\n%s", tc.width, tc.help)
			}
			// Help is user-facing: an internal contract id leaking into it
			// is the defect this guards, so assert the id is NOT present.
			if strings.Contains(tc.help, "0024") {
				t.Errorf("the %s help leaks the internal contract id "+
					"`0024:C4`; contract ids stay in source comments and "+
					"never reach user-facing output:\n%s",
					tc.width, tc.help)
			}
		})
	}
}
