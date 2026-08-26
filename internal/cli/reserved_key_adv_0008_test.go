package cli

// RDR 0008 Phase 3b — adversarial coverage for the CLI carrier.
//
// The RDR's Failure Modes section rests the whole reserved-key contract on
// one claim: the failure is "Visible" — it "fails table load/lint in the
// `reserved_tag_key` category" and "the failure data carries the
// direction-specific rule identifier, the offending name, and the remedy
// name, before any resolution runs." `0008:C3` sharpens that into REQ-28:
// "A category consumer may map the failure, but the guidance travels in the
// failure data, not the renderer."
//
// These tests attack the seam where that guidance is supposed to survive
// the trip to a consumer: the CLI. A payload that exists on
// `table.Failure` but never reaches the wire is guidance that travels in
// the renderer's prose after all, which is exactly what `0008:C3` forbids.
// JDR 0001 §D10 (REQ-102/REQ-103) names the carrier: 0008's per-key
// `reserved_tag_key` rides `findings[]` under one CLI code, with the
// near-miss advisory in `hint`.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/table"
)

// advModelPath writes src to a temp file and returns its path.
func advModelPath(t *testing.T, name, src string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return path
}

// TestAdv0008_LintCarriesReservedKeyPayloadOnTheWire pins that
// `intrastate lint` on a model breaching the reserved-key naming rule
// surfaces the three-field payload `0008:C3` mandates, as DATA a consumer
// reads — not folded into one prose `detail` string.
//
// FAILURE MODE: `internal/cli/lint.go:94-103` maps every `table.Load`
// refusal onto a bare `model-invalid` CLIError carrying `Detail:
// err.Error()` and NO `Findings` at all. The offending name, remedy name,
// and rule identifier that `internal/table/load.go:165-184` carefully
// populates are dropped at the only surface a user reaches. A consumer
// cannot branch on `reserved-tag-key/kernel-owned` vs
// `reserved-tag-key/author-must-rename`, and REQ-30's "a renderer MUST NOT
// present the reserved key as the required name in this direction" becomes
// unenforceable because the direction is not on the wire.
func TestAdv0008_LintCarriesReservedKeyPayloadOnTheWire(t *testing.T) {
	// The kernel-owned direction: a recognized-provenance declaration under
	// a name other than `recognized`.
	misnamed := strings.Replace(legalModel,
		"[tags.recognized]\nprovenance = \"recognized\"",
		"[tags.outcome]\nprovenance = \"recognized\"",
		1)
	if misnamed == legalModel {
		t.Fatalf("fixture edit did not apply; legalModel changed shape")
	}
	// Repoint the rules' match atoms at the renamed declaration so the ONLY
	// defect is the reserved-key naming rule.
	misnamed = strings.ReplaceAll(misnamed, "[rule.match.recognized]",
		"[rule.match.outcome]")

	path := advModelPath(t, "misnamed.toml", misnamed)

	// Confirm the loader really refuses in this category, so a green test
	// below could only mean the CLI dropped the payload.
	if _, err := table.Load([]byte(misnamed), path); err == nil {
		t.Fatalf("fixture loads clean; it must breach reserved_tag_key")
	} else if cat, ok := table.CategoryOf(err); !ok ||
		cat != table.CatReservedTagKey {
		t.Fatalf("fixture category = %q (ok=%v); want %q",
			cat, ok, table.CatReservedTagKey)
	}

	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err == nil {
		t.Fatalf("lint succeeded on a reserved_tag_key model:\n%s", stdout)
	}

	findings := flowFailureFindings(t, stdout)
	if len(findings) == 0 {
		t.Fatalf("lint emitted NO findings for a reserved_tag_key refusal; "+
			"`0008:C3` requires the offending name, remedy name, and rule "+
			"identifier to travel as failure DATA, not renderer prose "+
			"(REQ-26, REQ-28, REQ-103)\nstdout:\n%s", stdout)
	}

	var found bool
	for _, f := range findings {
		if f.Code != string(table.CatReservedTagKey) {
			continue
		}
		found = true
		if f.Rule != table.RuleKernelOwned {
			t.Errorf("finding rule = %q; want the stable identifier %q "+
				"asserted byte-for-byte (REQ-27, REQ-29, REQ-92)",
				f.Rule, table.RuleKernelOwned)
		}
		if !strings.Contains(f.Hint, table.RecognizedTagKey) {
			t.Errorf("finding hint = %q; the remedy name %q must reach the "+
				"consumer (REQ-29, JDR 0001 §D10 — the advisory/remedy "+
				"rides `hint`)", f.Hint, table.RecognizedTagKey)
		}
	}
	if !found {
		t.Errorf("no finding carries code %q; REQ-103 maps the category "+
			"onto `findings[]` under one CLI code\nfindings: %+v",
			table.CatReservedTagKey, findings)
	}
}

// TestAdv0008_LintSurfacesTheNearMissAdvisory pins the RDR's
// "**Advisory warning — normative, not deferred**" (REQ-58) at the surface
// that makes it non-deferred: a near-miss declaration must reach a user.
//
// FAILURE MODE: the advisory channel exists only on
// `table.LoadWithAdvisories`, and NOTHING outside `internal/table`'s own
// tests calls it — `internal/cli/lint.go:94` and
// `internal/cli/flow_input.go:121` both call `table.Load`. A model
// declaring `[tags.Recognized]` therefore lints entirely clean with no
// advisory anywhere a user can observe, which makes REQ-58's "normative,
// not deferred" a statement about an unreachable function. REQ-61 fixes
// the channel as advisory precisely so it can be SHOWN without blocking;
// showing it nowhere discharges the requirement by assertion.
func TestAdv0008_LintSurfacesTheNearMissAdvisory(t *testing.T) {
	// `Recognized` is a near-miss under folding alone and an ordinary,
	// legal, unreserved name — so the model must lint clean AND advise.
	nearMiss := strings.Replace(legalModel,
		"[tags.flag]\nprovenance = \"owned\"\nkind = \"bool\"",
		"[tags.Recognized]\nprovenance = \"owned\"\nkind = \"bool\"",
		1)
	if nearMiss == legalModel {
		t.Fatalf("fixture edit did not apply; legalModel changed shape")
	}
	nearMiss = strings.ReplaceAll(nearMiss, `keys = ["flag", "status"]`,
		`keys = ["Recognized", "status"]`)
	nearMiss = strings.ReplaceAll(nearMiss, "[rule.guard.all.flag]",
		"[rule.guard.all.Recognized]")
	nearMiss = strings.ReplaceAll(nearMiss, "flag = \"false\"",
		"Recognized = \"false\"")

	path := advModelPath(t, "near-miss.toml", nearMiss)

	// The advisory must not alter the verdict (REQ-63): the model loads.
	m, advisories, lerr := table.LoadWithAdvisories([]byte(nearMiss), path)
	if lerr != nil {
		t.Fatalf("near-miss fixture does not load: %v", lerr)
	}
	if m == nil {
		t.Fatalf("near-miss fixture loaded a nil model")
	}
	if len(advisories) == 0 {
		t.Fatalf("table.LoadWithAdvisories raised no advisory for " +
			"`Recognized`; the fixture is not exercising the near-miss rule")
	}

	stdout := requireSuccess(t, "lint", "--model", path, "--as=json")
	data := flowData(t, stdout)

	raw, ok := objectsAt(data, "findings")
	if !ok {
		t.Fatalf("lint success payload carries no `findings` array; "+
			"keys = %v", keysOf(data))
	}

	for _, f := range raw {
		rule, _ := f["rule"].(string)
		if rule == table.AdvisoryNearMiss {
			return // the advisory reached the user
		}
	}
	t.Errorf("lint reported no `%s` advisory for the near-miss declaration "+
		"`Recognized`; REQ-58 fixes the advisory as \"normative, not "+
		"deferred\", but table.LoadWithAdvisories has no caller outside "+
		"internal/table's own tests — internal/cli/lint.go:94 calls "+
		"table.Load\nfindings: %+v", table.AdvisoryNearMiss, raw)
}
