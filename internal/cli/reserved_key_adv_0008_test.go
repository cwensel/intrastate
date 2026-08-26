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
//
// PHASE 3c ADJUDICATION (see `artifacts/deviations.md` D9). Both seams are
// real, and both are DELIVERY gaps on a carrier RDR 0008 does not own:
//
//   - `0008:Approach` item 3 scopes this RDR to guaranteeing the failure
//     "exists and carries its guidance at load/lint; which user-facing
//     command surfaces it, and under which exit code, is RDR 0005's
//     mapping decision and is not settled here."
//   - REQ-26/28 bind the three fields "at the data level" — `table.Failure`
//     carries them — and REQ-93 says outright that this RDR's consumer "is
//     a test stub, not RDR 0005's exit-code map."
//   - Deviation D5 already settled the advisory as "a table-side advisory
//     list returned separately from the error; CLI delivery left to
//     0005/0006."
//   - Both call sites belong to already-Implemented RDRs (BUILD-ORDER runs
//     4 and 6, ahead of 0008 at run 7): `internal/cli/flow_input.go` is
//     RDR 0005's input grammar, and `internal/cli/lint.go`'s bare
//     `model-invalid` block traces to RDR 0006's own skeleton commit
//     (`d0194d9`) and is refused by 0006's OWN REQ-91, which forbids
//     collapsing findings into "a text-only `Detail` string."
//
// Those wire-level halves are now DELIVERED on the 0005/0006 carrier JDR
// 0001 §D10 item 3 fixes — `internal/cli/loadFindings` maps every load
// category onto `findings[]`, and `runLint` reads the advisory channel onto
// the success payload — so the skips are gone and the assertions run. What
// RDR 0008 owns — the three-field payload on `table.Failure` and the
// separate advisory channel — is still asserted at the data level first, so
// a regression in 0008's own surface fails here before the wire assertion
// can mask it.

import (
	"errors"
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
// FAILURE MODE GUARDED: `runLint` once mapped every `table.Load` refusal
// onto a bare `model-invalid` CLIError carrying `Detail: err.Error()` and NO
// `Findings` at all, dropping the offending name, remedy name, and rule
// identifier the loader populates at the only surface a user reaches. A
// consumer could not branch on `reserved-tag-key/kernel-owned` vs
// `reserved-tag-key/author-must-rename`, which made REQ-30's "a renderer
// MUST NOT present the reserved key as the required name in this direction"
// unenforceable — the direction was not on the wire.
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
	loadErr := func() error {
		_, err := table.Load([]byte(misnamed), path)
		return err
	}()
	if loadErr == nil {
		t.Fatalf("fixture loads clean; it must breach reserved_tag_key")
	} else if cat, ok := table.CategoryOf(loadErr); !ok ||
		cat != table.CatReservedTagKey {
		t.Fatalf("fixture category = %q (ok=%v); want %q",
			cat, ok, table.CatReservedTagKey)
	}

	// --- what RDR 0008 OWNS: the three-field payload, at the data level.
	//
	// REQ-26/28 bind the offending name, the remedy name, and the stable
	// rule identifier as DATA on the failure. This half is 0008's and is
	// asserted unconditionally.
	var f *table.Failure
	if !errors.As(loadErr, &f) {
		t.Fatalf("reserved_tag_key refusal is not a *table.Failure: %T", loadErr)
	}
	if f.Offending != "outcome" {
		t.Errorf("Offending = %q; want the name as authored %q (REQ-26/92)",
			f.Offending, "outcome")
	}
	if f.Remedy != table.RecognizedTagKey {
		t.Errorf("Remedy = %q; the kernel-owned direction renames TO the "+
			"reserved key %q (REQ-29/92)", f.Remedy, table.RecognizedTagKey)
	}
	if f.Rule != table.RuleKernelOwned {
		t.Errorf("Rule = %q; want the stable identifier %q asserted "+
			"byte-for-byte (REQ-27/29/92)", f.Rule, table.RuleKernelOwned)
	}

	// --- what RDR 0005/0006 OWN: getting that payload onto the wire.
	//
	// Now WIRED (JDR 0001 §D10 item 3): `internal/cli/loadFindings` maps every
	// load category onto `findings[]` under the one `model-invalid` code, and
	// carries RDR 0008's per-key payload on the same record.

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
		if f.Param != "outcome" {
			t.Errorf("finding param = %q; the offending name as authored "+
				"(%q) must reach the consumer (REQ-26, REQ-92)",
				f.Param, "outcome")
		}
		if f.Locator == "" {
			t.Errorf("finding carries no locator; a load category attributes " +
				"to file:line (REQ-14, REQ-18, JDR 0001 §D10 item 3)")
		}
	}
	if !found {
		t.Errorf("no finding carries code %q; REQ-103 maps the category "+
			"onto `findings[]` under one CLI code\nfindings: %+v",
			table.CatReservedTagKey, findings)
	}
}

// TestAdv0008_LintNeverNamesTheReservedKeyAsTheRequiredName pins the OTHER
// direction of the reserved-key rule at the wire.
//
// REQ-30 forbids a renderer presenting the reserved key as the required name
// when an owned or observed declaration is named `recognized` — the remedy
// there is "choose any other name", which is why `table.Failure.Remedy` is
// deliberately the EMPTY STRING for RuleAuthorMustRename. A wire hint that
// filled that gap with the reserved key would restate, in the renderer,
// exactly the advice the data withholds.
func TestAdv0008_LintNeverNamesTheReservedKeyAsTheRequiredName(t *testing.T) {
	// The author-must-rename direction: an OWNED declaration named
	// `recognized`. The sorted scan reports the FIRST violation it meets, and
	// the kernel-owned direction would win on a model breaching both, so this
	// fixture retypes the sole `recognized` declaration as owned rather than
	// introducing a second recognized-provenance name.
	misnamed := strings.Replace(legalModel,
		"[tags.recognized]\nprovenance = \"recognized\"\nkind = \"enum\"",
		"[tags.recognized]\nprovenance = \"owned\"\nkind = \"enum\"\n"+
			"domain = [\"go\", \"stop\"]",
		1)
	if misnamed == legalModel {
		t.Fatalf("fixture edit did not apply; legalModel changed shape")
	}

	path := advModelPath(t, "author-must-rename.toml", misnamed)

	// Confirm the loader takes the author-must-rename direction, so the wire
	// assertion below is about the direction it claims to be about.
	_, loadErr := table.Load([]byte(misnamed), path)
	var f *table.Failure
	if !errors.As(loadErr, &f) {
		t.Fatalf("fixture did not produce a *table.Failure: %v", loadErr)
	}
	if f.Rule != table.RuleAuthorMustRename {
		t.Fatalf("fixture rule = %q; want %q — the fixture must exercise the "+
			"rename-AWAY direction", f.Rule, table.RuleAuthorMustRename)
	}
	if f.Remedy != "" {
		t.Fatalf("fixture remedy = %q; this direction carries the EMPTY "+
			"remedy (REQ-30)", f.Remedy)
	}

	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err == nil {
		t.Fatalf("lint succeeded on a reserved_tag_key model:\n%s", stdout)
	}

	var found bool
	for _, got := range flowFailureFindings(t, stdout) {
		if got.Code != string(table.CatReservedTagKey) {
			continue
		}
		found = true
		if got.Rule != table.RuleAuthorMustRename {
			t.Errorf("finding rule = %q; want %q — the direction must be on "+
				"the wire, or REQ-30 is unenforceable by a consumer",
				got.Rule, table.RuleAuthorMustRename)
		}
		if strings.Contains(got.Hint, "`"+table.RecognizedTagKey+"`") {
			t.Errorf("finding hint = %q; it presents the reserved key %q as "+
				"the required name, which REQ-30 forbids in this direction",
				got.Hint, table.RecognizedTagKey)
		}
	}
	if !found {
		t.Errorf("no finding carries code %q\nstdout:\n%s",
			table.CatReservedTagKey, stdout)
	}
}

// TestAdv0008_LintCarriesEveryLoadCategoryAsAFinding pins REQ-91's actual
// blast radius: the collapse this fixes was never reserved-key-specific.
//
// `internal/cli/lint.go` mapped EVERY `table.Load` refusal — all of RDR
// 0002's load categories — onto one prose `detail` string. REQ-91 forbids
// that for the whole family, and JDR 0001 §D10 item 3 fixes the mapping as
// `code` = category slug, `locator` = file:line. A test that only covered
// the reserved-key category would spot-check the fix rather than pin it.
func TestAdv0008_LintCarriesEveryLoadCategoryAsAFinding(t *testing.T) {
	for _, tc := range []struct {
		name string
		src  string
		want table.Category
	}{
		{"malformed-toml", "this is not ][ toml", table.CatMalformedTOML},
		{
			"unsupported-version",
			"[model]\nid = \"x\"\nversion = 99\n",
			table.CatUnsupportedVersion,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := advModelPath(t, "cat.toml", tc.src)

			// The category the loader reports is the oracle; the wire must
			// carry the SAME slug, not a renderer's paraphrase of it.
			_, loadErr := table.Load([]byte(tc.src), path)
			if cat, ok := table.CategoryOf(loadErr); !ok || cat != tc.want {
				t.Fatalf("fixture category = %q (ok=%v); want %q",
					cat, ok, tc.want)
			}

			stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
			if err == nil {
				t.Fatalf("lint succeeded on a non-conforming model:\n%s",
					stdout)
			}

			findings := flowFailureFindings(t, stdout)
			if len(findings) == 0 {
				t.Fatalf("lint emitted NO findings for a %q refusal; REQ-91 "+
					"forbids a text-only `detail` as the machine-readable "+
					"route for ANY load category\nstdout:\n%s",
					tc.want, stdout)
			}
			if got := findings[0].Code; got != string(tc.want) {
				t.Errorf("finding code = %q; want the category slug %q "+
					"(JDR 0001 §D10 item 3)", got, tc.want)
			}
			if got := findings[0].Locator; got != path+":1" {
				t.Errorf("finding locator = %q; want %q — a load category "+
					"attributes to file:line (REQ-18)", got, path+":1")
			}
		})
	}
}

// TestAdv0008_LintSurfacesTheNearMissAdvisory pins the RDR's
// "**Advisory warning — normative, not deferred**" (REQ-58) at the surface
// that makes it non-deferred: a near-miss declaration must reach a user.
//
// FAILURE MODE GUARDED: the advisory channel once existed only on
// `table.LoadWithAdvisories`, which nothing outside `internal/table`'s own
// tests called. A model declaring `[tags.Recognized]` linted entirely clean
// with no advisory anywhere a user could observe, making REQ-58's
// "normative, not deferred" a statement about an unreachable function.
// REQ-61 fixes the channel as advisory precisely so it can be SHOWN without
// blocking; showing it nowhere discharges the requirement by assertion.
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

	// --- what RDR 0008 OWNS: the advisory's values and its separateness.
	//
	// REQ-62 fixes the two spellings plus the stable rule identifier;
	// REQ-63 fixes the channel as distinct from the validation-failure list
	// and forbids it altering the verdict (already asserted: lerr == nil).
	// This half is 0008's and is asserted unconditionally.
	var advised bool
	for _, a := range advisories {
		if a.Authored != "Recognized" {
			continue
		}
		advised = true
		if a.Reserved != table.RecognizedTagKey {
			t.Errorf("advisory Reserved = %q; want the reserved spelling %q "+
				"(REQ-62)", a.Reserved, table.RecognizedTagKey)
		}
		if a.Rule != table.AdvisoryNearMiss {
			t.Errorf("advisory Rule = %q; want the literal %q (REQ-62)",
				a.Rule, table.AdvisoryNearMiss)
		}
	}
	if !advised {
		t.Errorf("no advisory names the authored spelling `Recognized`; "+
			"REQ-59/62 require the advisory to name BOTH spellings\n"+
			"advisories: %+v", advisories)
	}

	// --- what RDR 0005/0006 OWN: surfacing it to a user.
	//
	// Now WIRED: `runLint` calls `table.LoadWithAdvisories` and appends each
	// advisory to the SUCCESS payload's `data.findings` — the carrier REQ-93
	// already fixes — so the verdict stays untouched (REQ-63).

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
