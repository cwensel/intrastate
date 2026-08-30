package cli

// RDR 0023 — the Minimum Viable Validation, run end to end.
//
// `0023:MVV` is six steps over the CHECKED-IN decision-table fixture, and
// this file runs all six as one battery. Two properties of the battery are
// load-bearing and are asserted rather than assumed:
//
//  1. Step 1 is the ONLY comparison with a pre-change side. "S1–S5 all
//     compare the new build against itself, so none of them can see a
//     default-mode regression that moves both sides together; without a
//     captured golden, C1's headline 'byte-identical to today's' is
//     asserted nowhere." So step 1 compares against a GOLDEN CAPTURED FROM
//     THE PRE-CHANGE BINARY, checked in before Phase 1 edits the struct.
//
//  2. The pass bar is a MEASURED RATIO, not a green exit code. Step 6's bar
//     is "the CHECKED-IN fixtures … each save at least 40%", measured on
//     C1's unit (the full emitted line). A run that merely does not error
//     tells us nothing: fidelity loss and zero-saving projections both exit
//     0. So every reconstructed value is compared byte-for-byte and every
//     width is compared as a number.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/spf13/cobra"
)

// goldenPath0023 is where MVV step 1's pre-change golden is checked in.
//
// `0023:A-8` fixes the ORDERING, not the location: "the capture commit
// precedes the conversion commit, and a golden captured after any Phase 1
// edit does not satisfy REQ-96/REQ-97". This path is that location.
const goldenPath0023 = "docs/rdr/0023-resolve-envelope-projection/artifacts/" +
	"mvv-step1-default-golden.json"

// mvvCall0023 is the MVV's own invocation: the checked-in pricing decision
// table, "a recognized outcome and discriminating `--tag`s".
//
// Discriminating matters: `tier=paid` + `region=eu` selects exactly
// `paid-eu` out of the 2×2, so the golden pins a REAL selection rather than
// a degenerate one.
func mvvCall0023(t *testing.T, extra ...string) []string {
	t.Helper()

	return resolveArgs(pricingModelPath(t), "", "decide",
		append([]string{"--tag", "tier=paid", "--tag", "region=eu"}, extra...)...)
}

// REQ-95 / MVV step 1: "Run `flow resolve` with a recognized outcome and
// discriminating `--tag`s, `--as=json`, without the flag → today's full
// payload, byte-identical to the pre-change fixture."
// REQ-96: "The fixture is a GOLDEN CAPTURED FROM THE PRE-CHANGE BINARY and
// checked in BEFORE Phase 1 edits the struct (Phase 1's first step, `git
// stash`-clean tree) — this is the only comparison in the whole battery
// that has a pre-change side."
// REQ-97: "S1–S5 all compare the new build against itself, so none of them
// can see a default-mode regression that moves both sides together; without
// a captured golden, C1's headline \"byte-identical to today's\" is
// asserted nowhere"
// REQ-27: "Default-mode output is byte-identical to today's: the request
// echo is the correct default and this RDR does not change it."
// REQ-28: the `Overrides` narrowing — "Default-mode output is byte-identical
// to today's; every other 0005 obligation on `flow resolve` … is
// unchanged." — [0005-carried]
// REQ-129: "The flag is additive and opt-in, so no existing caller observes
// a change; default-mode output is byte-identical to the pre-change binary
// (C1, guarded by the MVV's pre-change golden)."
// A-8: the ordering is normative even though the path is not.
// HAPPY PATH — the battery's ONLY pre-change comparison, and the only place
// C1's headline claim is assertable at all.
func TestReq27And28And95And96And97And129_DefaultModeIsByteIdenticalToTheCapturedPreChangeGolden(t *testing.T) {
	path := filepath.Join(repoRootFor(t), goldenPath0023)
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("the pre-change golden is not checked in at %s: %v\n\n"+
			"This is Phase 1's FIRST step, on a `git stash`-clean tree, "+
			"BEFORE the struct is edited. It is the only comparison in the "+
			"whole battery with a pre-change side: S1-S5 all compare the new "+
			"build against itself, so none of them can see a default-mode "+
			"regression that moves BOTH sides together. Without it C1's "+
			"headline \"byte-identical to today's\" is asserted nowhere and "+
			"every existing consumer rides on a manual step.",
			goldenPath0023, err)
	}

	// The observed line's `model` value is the `--model` argument VERBATIM
	// (`flow_input.go::selectModel` returns the path unchanged), and this
	// package's tests run with `cwd = internal/cli`, so `mvvCall0023` has to
	// pass an ABSOLUTE path for the model to resolve at all. Folding the
	// checkout root back out is what lets the golden be the RELATIVE
	// spelling the RDR names — `0023:MVV` step 1 and REQ-107 both write
	// `models/examples/pricing-decision-table.toml`, and that is the
	// spelling `evidence/spikes/a1-byte-width.md` row S1 and
	// `a2-encoder-mechanism.md`'s reference bytes are measured on.
	//
	// This narrows nothing. Byte-identity is still asserted over the whole
	// record, `model` included: the key must be present and its value must
	// be that relative path. Only the prefix that differs between two
	// checkouts of the same commit is folded — a golden carrying it would
	// fail everywhere but the machine that captured it, and would commit an
	// absolute local path into a checked-in artifact.
	got, err := foldEmitDispositions0024(foldCheckoutRoot(t,
		emittedLine(t, requireSuccess(t, append(mvvCall0023(t), "--as=json")...))))
	if err != nil {
		t.Fatalf("%v", err)
	}
	if got != strings.TrimRight(string(want), "\n") {
		t.Errorf("default-mode output differs from the pre-change golden:\n"+
			"  golden = %s\n  now    = %s\n"+
			"The request echo is the CORRECT DEFAULT and this RDR does not "+
			"change it. The flag is additive and opt-in, so no existing "+
			"caller observes a change — and this is the assertion that "+
			"holds that promise.", strings.TrimRight(string(want), "\n"), got)
	}
}

// REQ-107: the MVV runs "Over the checked-in decision-table fixture
// (`models/examples/pricing-decision-table.toml` or the 0005 fixture
// family)"
// REQ-104: the step-6 PASS BAR — "every shape measured saves bytes (the S1
// width oracle already enforces this per-run), and the CHECKED-IN fixtures
// (`models/examples/pricing-decision-table.toml`, `release-grammar.toml`,
// `review-state-machine.toml`) each save at least 40%"
// REQ-105: "a landing below 40% there means the plan group is carrying
// materially more than A1 measured, which is A1's \"saves too little to
// justify a new surface\" condition and routes back rather than recording a
// number."
// REQ-103 / MVV step 6: "Record default vs projected byte counts on the
// motivating-model shape and one gate/write-heavy fixture, measured on the
// full emitted line (C1's unit)"
// REQ-106: "The 79.4% figure is NOT a pass bar, because its shape is not
// checked in" — negative REQ: the 48-fact demonstration is not gated here.
// REQ-109: "the corresponding widths (290 B full line → 152 B projected,
// envelope included — the same unit this scenario asserts) are row S1 of
// `evidence/spikes/a1-byte-width.md`."
// A-5: the 40% bar is MVV-only; S1's per-run oracle carries no threshold.
//
// THE MVV — a RUNNABLE end-to-end battery, all six steps, on the checked-in
// fixtures. The pass bar is a MEASURED RATIO on C1's unit: a green exit
// code is explicitly NOT sufficient, since a zero-saving projection and a
// fidelity-losing one both exit 0.
func TestMVV0023_ResolveEnvelopeProjectionEndToEnd(t *testing.T) {
	// --- step 1: default mode, the full payload -----------------------
	//
	// The golden comparison is its own oracle above; here step 1 supplies
	// the EXPECTED SIDE for steps 2 and 3, which is its other role.
	step1 := requireSuccess(t, append(mvvCall0023(t), "--as=json")...)
	step1Line := emittedLine(t, step1)
	step1Fields := rawFieldsOf(t, step1)
	step1Text := requireSuccess(t, append(mvvCall0023(t), "--as=text")...)

	t.Run("step1_default_payload_is_the_full_width", func(t *testing.T) {
		// All fourteen fields, minus `escape_class` on this non-escaped
		// shape (`0023:A-11`: 14 struct fields, 13 rendered keys here).
		for _, key := range append(slices.Clone(echoGroup0023),
			"revision", "rule", "gates", "emit", "next", "writes", "clear",
			"escaped") {
			if _, held := step1Fields[key]; !held {
				t.Errorf("step 1's default payload carries no %q; it is "+
					"today's FULL payload and `0005:A6`'s envelope list is "+
					"unchanged in the default width", key)
			}
		}
		if _, held := step1Fields["escape_class"]; held {
			t.Errorf("step 1 carries `escape_class` on an unescaped " +
				"selection; its producer's presence rule (`0005:A-3`) omits " +
				"it, and the projected literal's ninth key is conditional " +
				"on exactly this")
		}
	})

	// --- step 2: --plan-only, json ------------------------------------
	//
	// "exactly the normative projected key set …; every carried field
	// byte-identical to step 1's; model/observed/owned/readers/outcome
	// absent (not null, not empty)."
	step2 := requireSuccess(t, append(mvvCall0023(t, "--"+planOnlyFlag), "--as=json")...)
	step2Line := emittedLine(t, step2)
	step2Fields := rawFieldsOf(t, step2)

	t.Run("step2_projected_key_set_and_byte_identical_carried_values", func(t *testing.T) {
		want := slices.Clone(projectedKeyLiteral0023)
		if _, carried := step1Fields["escape_class"]; carried {
			want = append(want, "escape_class")
		}
		slices.Sort(want)

		if got := sortedKeys(step2Fields); !slices.Equal(got, want) {
			t.Errorf("step 2's key set = %v;\nwant              %v", got, want)
		}

		// THE INVERSE INVARIANT, value-for-value: every carried field is
		// byte-identical to step 1's. A green exit code proves nothing —
		// fidelity loss hides behind a passing run, which is why each value
		// is reconstructed and compared against its step-1 original rather
		// than merely counted.
		for key, raw := range step2Fields {
			orig, held := step1Fields[key]
			if !held {
				t.Errorf("step 2 carries %q, which step 1 does not", key)
				continue
			}
			if string(raw) != string(orig) {
				t.Errorf("carried field %q is NOT byte-identical to step "+
					"1's:\n  step 1 = %s\n  step 2 = %s\nProjection is key "+
					"deletion; every surviving field's bytes are step 1's "+
					"own", key, orig, raw)
			}
		}

		for _, key := range echoGroup0023 {
			if raw, held := step2Fields[key]; held {
				t.Errorf("step 2 carries the echo key %q as %s; it must be "+
					"ABSENT — not null, not empty", key, raw)
			}
		}
	})

	// --- step 3: --plan-only, text ------------------------------------
	//
	// REQ-99 / MVV step 3: "Re-run step 2 with `--as=text` → the projected
	// lines are a byte-identical subset of step 1's text lines, stable
	// across repeated runs."
	t.Run("step3_projected_text_is_a_stable_byte_identical_subset", func(t *testing.T) {
		var step3 string
		for i := range 3 {
			out := requireSuccess(t,
				append(mvvCall0023(t, "--"+planOnlyFlag), "--as=text")...)
			if i == 0 {
				step3 = out
				continue
			}
			if out != step3 {
				t.Fatalf("step 3 run %d differs from run 0; the projected "+
					"lines are STABLE ACROSS REPEATED RUNS\nfirst:\n%s\n"+
					"again:\n%s", i, step3, out)
			}
		}

		defLines := textLines(step1Text)
		for _, line := range textLines(step3) {
			if !slices.Contains(defLines, line) {
				t.Errorf("step 3 line %q is not a byte-identical member of "+
					"step 1's text lines\nstep 1:\n%s\nstep 3:\n%s",
					line, step1Text, step3)
			}
		}
		if textWidth(step3) >= textWidth(step1Text) {
			t.Errorf("step 3's text is %d B against step 1's %d B; C1's "+
				"width clause binds text mode too",
				textWidth(step3), textWidth(step1Text))
		}
	})

	// --- step 4: unrecognized outcome, ± the flag ---------------------
	t.Run("step4_refusal_envelopes_and_exits_are_byte_identical", func(t *testing.T) {
		args := resolveArgs(pricingModelPath(t), "", "not-an-outcome",
			"--tag", "tier=paid", "--tag", "region=eu")

		defOut, _, defErr := runCmd(t, append(slices.Clone(args), "--as=json")...)
		projOut, _, projErr := runCmd(t,
			append(slices.Clone(args), "--"+planOnlyFlag, "--as=json")...)

		if defErr == nil {
			t.Fatalf("the unrecognized outcome SUCCEEDED; step 4 exists to "+
				"exercise a refusal\n%s", defOut)
		}
		if code := clierr.ErrorCode(defErr); code != "flow-unmodeled-outcome" {
			t.Fatalf("step 4's refusal code = %q; want "+
				"`flow-unmodeled-outcome`", code)
		}
		if got, want := strings.TrimRight(projOut, "\n"),
			strings.TrimRight(defOut, "\n"); got != want {
			t.Errorf("step 4's refusal envelopes differ ± the flag:\n"+
				"  default   = %s\n  projected = %s", want, got)
		}
		if got, want := clierr.ExitCodeFor(projErr), clierr.ExitCodeFor(defErr); got != want {
			t.Errorf("step 4's exit = %d ± the flag and %d by default", got, want)
		}
	})

	// --- step 5: the sibling verb and the structural walk -------------
	//
	// REQ-101: "Run `flow next --plan-only` → `command-error`, exit 2; and
	// the whole-tree structural oracle … passes."
	// REQ-102: "One sibling verb is run, not all three, deliberately: C1
	// makes the behavioural run CORROBORATION and the structural walk the
	// assertion" — negative REQ: `read-state`/`set-state` owe no
	// behavioural run here.
	t.Run("step5_sibling_verb_refuses_and_the_structural_walk_passes", func(t *testing.T) {
		model := writeFlowModel(t, flowMVVModel)
		art := seedArtifact(t, model, "status=draft")

		_, _, err := runCmd(t, "flow", "next", "--model", model,
			"--artifact", artifactBinding(flowStateRole, art),
			"--"+planOnlyFlag, "--as=json")
		if err == nil {
			t.Fatalf("`flow next --%s` SUCCEEDED", planOnlyFlag)
		}
		if code := clierr.ErrorCode(err); code != "command-error" {
			t.Errorf("code = %q; want `command-error` — the shared usage "+
				"bucket, `0011:A14`'s accepted limit", code)
		}
		if exit := clierr.ExitCodeFor(err); exit != 2 {
			t.Errorf("exit = %d; want 2", exit)
		}

		// The ASSERTION half: the structural walk, over a root with both
		// auto-generated commands materialized and present in the walked
		// set. The behavioural run above is CORROBORATION, never this.
		root := NewRootCmd()
		root.InitDefaultHelpCmd()
		root.InitDefaultCompletionCmd()

		var walked, registrants []string
		var walk func(*cobra.Command)
		walk = func(c *cobra.Command) {
			walked = append(walked, c.Name())
			if planOnlyRegistered(c) {
				registrants = append(registrants, commandPath(c))
			}
			for _, sub := range c.Commands() {
				walk(sub)
			}
		}
		walk(root)

		for _, required := range []string{"help", "completion"} {
			if !slices.Contains(walked, required) {
				t.Fatalf("the walked set lacks %q; the precondition is "+
					"asserted, not assumed — `completion` is the "+
					"discriminating member a bare `NewRootCmd()` tree cannot "+
					"supply\nwalked = %v", required, walked)
			}
		}
		slices.Sort(registrants)
		if want := []string{"intrastate flow resolve"}; !slices.Equal(registrants, want) {
			t.Errorf("commands registering --%s = %v; want %v",
				planOnlyFlag, registrants, want)
		}
	})

	// --- step 6: the byte table and the 40% PASS BAR ------------------
	//
	// The bar is stated on shapes the implementer can actually run: the
	// three CHECKED-IN fixtures. A ratio, measured on C1's unit — not a
	// green exit code, and not the not-checked-in 79.4% demonstration.
	t.Run("step6_checked_in_fixtures_each_save_at_least_forty_percent", func(t *testing.T) {
		reviewModel := reviewModelPath(t)
		reviewArt := newFlowArtifact(t, "review.artifact")
		reviewBind := artifactBinding("review", reviewArt)
		if _, _, err := runCmd(t, "flow", "set-state", "--model", reviewModel,
			"--artifact", reviewBind, "--write", "status=draft", "--as=json",
		); err != nil {
			t.Fatalf("seeding the review fixture failed: %v", err)
		}

		// The release grammar's `begin` row guards on `phase = idle` and
		// writes `build-id`, so BOTH owned keys must be established before
		// the row is decidable — an unseeded artifact refuses
		// `flow-owned-state-unavailable` rather than producing a plan to
		// measure. `evidence/spikes/a1-byte-width.md` §S2b seeds exactly
		// `phase=idle`, `build-id=none`.
		releaseModel := releaseModelPath(t)
		releaseArt := newFlowArtifact(t, "release.artifact")
		releaseBind := artifactBinding("release", releaseArt)
		if _, _, err := runCmd(t, "flow", "set-state", "--model", releaseModel,
			"--artifact", releaseBind, "--write", "phase=idle",
			"--write", "build-id=none", "--as=json",
		); err != nil {
			t.Fatalf("seeding the release fixture failed: %v", err)
		}

		for _, shape := range []struct {
			name string
			args []string
		}{
			{
				// The motivating-model shape: the pricing decision table.
				name: "pricing-decision-table",
				args: mvvCall0023(t),
			},
			{
				// The gate/write-heavy fixture step 6 names beside it.
				//
				// The two observed `--tag`s are the A1 spike's own S2b
				// invocation, not decoration: step 6 says to "compare
				// against A1's baseline", and the baseline row is measured
				// on this argv. Dropping them empties `observed` and
				// measures a DIFFERENT shape — one whose echo group is
				// smaller than the one the 40% bar was set against.
				name: "release-grammar/begin",
				args: resolveArgs(releaseModel, releaseBind, "build",
					"--tag", "risk=0", "--tag", "checks=[]"),
			},
			{
				name: "review-state-machine/approve",
				args: resolveArgs(reviewModel, reviewBind, "approve"),
			},
		} {
			t.Run(shape.name, func(t *testing.T) {
				defOut, _, err := runCmd(t, append(slices.Clone(shape.args), "--as=json")...)
				if err != nil {
					t.Fatalf("the default run failed: %v\n%s", err, defOut)
				}
				projOut, _, err := runCmd(t,
					append(slices.Clone(shape.args), "--"+planOnlyFlag, "--as=json")...)
				if err != nil {
					t.Fatalf("the projected run failed: %v\n%s", err, projOut)
				}

				// The checkout root is folded out of BOTH widths before
				// measuring, for the same reason D-1 folds it out of the
				// golden: `model` carries the `--model` argument verbatim,
				// and the harness must pass an ABSOLUTE path because the
				// package's tests run with `cwd = internal/cli`. That
				// prefix lands only in the DEFAULT payload — `model` is an
				// echo field the projection drops — so leaving it in
				// inflates `def` alone and reports a saving that rises with
				// the length of the capturing machine's checkout path.
				//
				// That is not a cosmetic bias: it measures the wrong thing
				// in the favourable direction. `0023:A1` measured
				// release-grammar at 412 B → 238 B (42.2%), clearing the
				// 40% bar by 2.2 points; an absolute prefix reports the
				// same run at ~50%, so a genuine regression below the bar
				// could still pass here while failing on a short checkout
				// path. Folding restores the RELATIVE spelling the RDR and
				// both spikes name as normative, making the ratio a
				// property of the payload rather than of the filesystem.
				def := len(foldCheckoutRoot(t, emittedLine(t, defOut)))
				proj := len(foldCheckoutRoot(t, emittedLine(t, projOut)))
				saved := float64(def-proj) / float64(def) * 100

				// The pass bar. Not "did not error" — a MEASURED RATIO on
				// the full emitted line, C1's unit.
				if saved < 40 {
					t.Errorf("%s saved %.1f%% (%d B → %d B) on the full "+
						"emitted line; the CHECKED-IN fixtures must each "+
						"save at least 40%%. `0023:A1` measured 42.2-51.3%% "+
						"on them, so a landing below 40%% means the plan "+
						"group is carrying materially more than A1 measured "+
						"— A1's \"saves too little to justify a new "+
						"surface\" condition, which ROUTES BACK rather than "+
						"recording a number", shape.name, saved, def, proj)
				}
				t.Logf("%s: %d B → %d B (%.1f%% saved, full emitted line)",
					shape.name, def, proj, saved)
			})
		}

		// The A1 row-S1 witness, for the record: 290 B → 152 B on the
		// pricing 2×2 call, envelope included. Logged rather than asserted
		// — REQ-114 makes A1's table "recorded evidence, not a test
		// assertion", and pinning the exact byte count would turn an
		// unrelated encoder change into a false failure here.
		t.Logf("A1 row S1 baseline (evidence/spikes/a1-byte-width.md): "+
			"290 B → 152 B, 47.6%% saved; this run measured %d B → %d B",
			len(step1Line), len(step2Line))
	})
}

// REQ-136: "this RDR ships the CAPABILITY, not its adoption. The flag is
// opt-in and no consumer passes it on landing … Phase 2 green plus the MVV
// step-6 table is the completion bar"
// REQ-137: "`intrastate#srz2` (P1) is the ADOPTION ask, so this RDR landing
// green does not close it — it unblocks it." — negative REQ: no consumer
// call site is migrated by this RDR.
// REQ-134: "Recovery: drop the flag — the default width is the full record
// and is contractually unchanged."
// DOMAIN EDGE — the deliverable's boundary, asserted as an absence: no
// checked-in consumer surface passes the flag, and dropping it recovers the
// full record.
func TestReq134And136And137_TheFlagIsOptInAndNoConsumerPassesItOnLanding(t *testing.T) {
	// CONTROL — "no consumer passes the flag" is vacuously true of a tree
	// where the flag does not exist, and would report coverage of a
	// deliverable boundary that is not yet drawn. The capability must SHIP
	// before its non-adoption says anything.
	sub := resolveCmd(t)
	if sub == nil || sub.Flags().Lookup(planOnlyFlag) == nil {
		t.Fatalf("`flow resolve` does not register --%s; the RDR ships the "+
			"CAPABILITY, and until it exists \"no consumer passes it\" is "+
			"true for the wrong reason. The completion bar is Phase 2 green "+
			"PLUS the MVV step-6 table — not an absence that holds because "+
			"nothing was built", planOnlyFlag)
	}

	// No checked-in model, doc-example invocation or skill surface passes
	// the flag: the RDR ships the capability, not its adoption, and the
	// motivating tracker is UNBLOCKED rather than closed by a green build.
	root := repoRootFor(t)
	for _, dir := range []string{"models", "docs/skills"} {
		path := filepath.Join(root, dir)
		if _, err := os.Stat(path); err != nil {
			continue
		}
		err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			src, rerr := os.ReadFile(p)
			if rerr != nil {
				return nil
			}
			if strings.Contains(string(src), "--"+planOnlyFlag) {
				t.Errorf("%s passes --%s; this RDR ships the CAPABILITY, "+
					"not its adoption — no consumer passes it on landing, "+
					"and the realized saving is a follow-on migration's to "+
					"deliver", p, planOnlyFlag)
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", path, err)
		}
	}

	// Recovery: dropping the flag returns the full record, contractually
	// unchanged — which is what makes the opt-in reversible.
	full := rawFieldsOf(t, requireSuccess(t, append(mvvCall0023(t), "--as=json")...))
	for _, key := range echoGroup0023 {
		if _, held := full[key]; !held {
			t.Errorf("dropping the flag does not restore %q; the default "+
				"width is the FULL RECORD and is contractually unchanged — "+
				"that is the whole recovery procedure", key)
		}
	}
}

// REQ-115: "`docs/cli-output-contract.md` gains the projected worked
// payload beside the full one; `flowResolveExtendedDesc` gains the flag
// under \"Reading a successful plan\"; the partition statement (C2) lands
// where the payload fields are documented."
// REQ-116: "`flowResolveExtendedDesc` and the flag's own usage string are
// inputs to the GENERATED artifacts `docs/cli-reference.md` and `llms.txt`
// … So this phase regenerates those files and commits them in the same
// change."
// REQ-117: "because the generated artifacts derive from the flag's
// registration, the commit that registers the flag MUST also carry the
// regenerated `docs/cli-reference.md` and `llms.txt`"
// REQ-118: "The flag's usage string is authored here, once, since it ships
// into a committed artifact rather than staying an implementation detail."
// REQ-119: "the help line states the omitted group by name"
// REQ-120: "Verb help | `flow_resolve.go::flowResolveExtendedDesc` …
// Extend | The help's plan list is the partition's user-facing statement;
// gains the flag line"
// DOMAIN EDGE — the docs surface, which `make docs-check` turns red if it
// drifts. Asserted on content rather than on the Makefile: the flag must be
// documented, its usage string non-empty, and the omitted group NAMED.
func TestReq115And116And117And118And119And120_TheFlagIsDocumentedAndTheOmittedGroupIsNamed(t *testing.T) {
	sub := resolveCmd(t)
	if sub == nil {
		t.Fatal("the `flow` group registers no `resolve` verb")
	}
	f := sub.Flags().Lookup(planOnlyFlag)
	if f == nil {
		t.Fatalf("`flow resolve` does not register --%s", planOnlyFlag)
	}

	// The usage string is AUTHORED, once — it ships into a committed
	// generated artifact rather than staying an implementation detail.
	if strings.TrimSpace(f.Usage) == "" {
		t.Errorf("--%s carries an empty usage string; it is an input to the "+
			"GENERATED `docs/cli-reference.md` and `llms.txt`, so it is "+
			"authored here, once", planOnlyFlag)
	}

	// The verb's extended help gains the flag line, and it STATES THE
	// OMITTED GROUP BY NAME — the mitigation for the model-reference risk,
	// since `model` is precisely the field a projected caller loses.
	help := flowResolveExtendedDesc
	if !strings.Contains(help, planOnlyFlag) {
		t.Errorf("`flowResolveExtendedDesc` does not mention --%s; the "+
			"help's plan list is the partition's USER-FACING STATEMENT and "+
			"gains the flag line under \"Reading a successful plan\"",
			planOnlyFlag)
	}
	var named []string
	for _, key := range echoGroup0023 {
		if !strings.Contains(help, key) {
			named = append(named, key)
		}
	}
	if len(named) > 0 {
		t.Errorf("the verb help does not name the omitted echo fields %v; "+
			"the help line STATES THE OMITTED GROUP BY NAME, which is what "+
			"makes the opt-in safe for a caller who would otherwise lose "+
			"the model reference", named)
	}

	// `docs/cli-reference.md` derives from the flag's REGISTRATION — it
	// renders every command's full flag set — so it must already carry the
	// flag: a commit that registers it without regenerating turns
	// `make check` red on `docs-check`.
	//
	// `llms.txt` is NOT asserted the same way. It is a command INDEX: its
	// generator (`docs.go::writeLLMsTxt`) emits `CommandPath + Short` per
	// command plus static prose, and renders no command's flags at all —
	// `--all`, `--outcome`, `--artifact` and `--model` are all absent from
	// it today. REQ-116's premise that the usage string is an input to that
	// file does not hold of the shipped generator (deviation D-6,
	// SPEC-DEFECT). The obligation REQ-117 actually states — that the
	// registering commit must leave the generated artifacts NOT STALE — is
	// asserted here for both files, by re-rendering and comparing, which is
	// exactly what `make docs-check` does and is strictly stronger than a
	// substring probe for the one file that does carry flags.
	generated := t.TempDir()
	if _, _, err := runCmd(t, "docs", "--dir", generated, "--as=json"); err != nil {
		t.Fatalf("re-rendering the generated artifacts failed: %v", err)
	}
	for _, artifact := range []string{"docs/cli-reference.md", "llms.txt"} {
		committed, err := os.ReadFile(filepath.Join(repoRootFor(t), artifact))
		if err != nil {
			t.Errorf("read %s: %v", artifact, err)
			continue
		}
		fresh, err := os.ReadFile(filepath.Join(generated, artifact))
		if err != nil {
			t.Errorf("read the freshly rendered %s: %v", artifact, err)
			continue
		}
		if string(committed) != string(fresh) {
			t.Errorf("%s is STALE against the command tree this binary "+
				"carries; the generated artifacts derive from the flag's "+
				"REGISTRATION, so the commit that registers it must also carry "+
				"them regenerated — a phase split that leaves the flag "+
				"registered and the artifacts stale is a red build by "+
				"construction. Run `make docs`", artifact)
		}
	}

	// And the file that DOES render flags carries this one, with its
	// authored usage string.
	reference, err := os.ReadFile(
		filepath.Join(repoRootFor(t), "docs/cli-reference.md"))
	if err != nil {
		t.Fatalf("read docs/cli-reference.md: %v", err)
	}
	if !strings.Contains(string(reference), "--"+planOnlyFlag) {
		t.Errorf("`docs/cli-reference.md` does not carry --%s; it renders "+
			"every command's full flag set from the live tree, so the flag's "+
			"authored usage string reaches a committed artifact through it",
			planOnlyFlag)
	}

	// The worked payload lands beside the full one in the output contract.
	contract := filepath.Join(repoRootFor(t), "docs/cli-output-contract.md")
	src, err := os.ReadFile(contract)
	if err != nil {
		t.Fatalf("read docs/cli-output-contract.md: %v", err)
	}
	if !strings.Contains(string(src), planOnlyFlag) {
		t.Errorf("`docs/cli-output-contract.md` does not mention --%s; it "+
			"gains the PROJECTED WORKED PAYLOAD beside the full one, and "+
			"the partition statement (C2) lands where the payload fields "+
			"are documented", planOnlyFlag)
	}
}

// REQ-126: "Wall time unchanged: projection is key deletion on the
// assembled payload before `respond.OK` — no extra I/O, no second marshal
// path"; and "no oracle or scenario asserts a wall-time threshold, and none
// is specified" — negative REQ: no wall-time budget is owed.
// REQ-83: "S4 is the one structural-absence oracle here, and the vacuity
// guard is what stops it passing on an empty command tree."
// REQ-112: "Five oracles, one per Testing Strategy scenario S1–S5, and the
// enumeration here is that list — not a sixth"
// A-9: "the build carries five oracles, one of which (S2) gains the A9 arm
// early and one of which (S4) is a new symbol."
// DOMAIN EDGE — the enumeration's own shape, asserted as the absence of a
// wall-time bar and the presence of exactly the five scenario families.
// This is a structural claim about the suite, so it is asserted over the
// suite's own surface rather than over the binary's behaviour.
func TestReq83And112And126_TheOracleBatteryIsTheFiveScenariosAndOwesNoWallTimeBudget(t *testing.T) {
	// No wall-time threshold is specified, so none may be asserted: a
	// millisecond bar on a 5ms call would measure the harness rather than
	// the change. Assert the absence over this RDR's own test surface.
	dir := filepath.Join(repoRootFor(t), "internal/cli")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), "_0023_test.go") {
			continue
		}
		src, rerr := os.ReadFile(filepath.Join(dir, e.Name()))
		if rerr != nil {
			t.Fatalf("read %s: %v", e.Name(), rerr)
		}
		text := string(src)
		// Assembled at runtime rather than written as literals, so this
		// oracle does not match its own source and report a finding
		// against itself.
		for _, banned := range []string{
			"time" + ".Since", "time" + ".Now(", "testing" + ".Benchmark",
		} {
			if strings.Contains(text, banned) {
				t.Errorf("%s references %s; NO oracle or scenario asserts a "+
					"wall-time threshold and none is specified — the work "+
					"removed (encoding five fewer keys) cannot make the call "+
					"slower, and a millisecond bar on a 5ms call would "+
					"measure the harness rather than the change",
					e.Name(), banned)
			}
		}
	}

	// The vacuity guard S4 rests on: the four flow verbs exist, so the one
	// structural-absence oracle in this RDR cannot pass on an empty tree.
	if names := flowGroupNames(t); len(names) != len(flowVerbs) {
		t.Fatalf("the `flow` group registers %v; S4 is the one "+
			"structural-absence oracle here and the vacuity guard is what "+
			"stops it passing on an empty command tree", names)
	}

	// The enumeration itself: FIVE oracles, one per Testing Strategy
	// scenario, and "the enumeration here is that list — not a sixth". Each
	// scenario's subject must be REACHABLE for its oracle to be more than a
	// name, and the flag is the subject of all five — S1's differential,
	// S2's key set, S3's emitted-keys arm, S4's registrant set, S5's text
	// subset. So the enumeration is only satisfiable once the flag ships.
	sub := resolveCmd(t)
	if sub == nil || sub.Flags().Lookup(planOnlyFlag) == nil {
		t.Fatalf("`flow resolve` does not register --%s; five oracles are "+
			"owed (S1 differential, S2 explicit key set carrying "+
			"absent-not-null, S3 reflective partition-completeness, S4 "+
			"whole-tree walk, S5 text subset + width) and four of the five "+
			"have no reachable subject without it. `0023:A-9` reads the "+
			"battery as five oracles — one (S2) gaining A9's arm in Phase "+
			"1, one (S4) a new symbol — never four plus an aspiration",
			planOnlyFlag)
	}
}

// foldCheckoutRoot removes the capturing machine's checkout prefix from an
// emitted record, so a width or a golden comparison is a property of the
// PAYLOAD rather than of the filesystem it was captured on.
//
// `resolvePayload.Model` carries the `--model` argument verbatim
// (`flow_input.go::selectModel` returns `path` unchanged), and the harness
// must pass an ABSOLUTE path because the package's tests run with
// `cwd = internal/cli`. What lands in the record is therefore
// checkout-specific, and `model` is an echo field the projection drops — so
// an unfolded prefix inflates the DEFAULT side alone.
//
// BOTH spellings are folded, and that is the portability point. The record
// is JSON, so on a platform whose separator needs escaping (Windows,
// `C:\Users\…` → `C:\\Users\\…`) the RAW prefix does not occur in the
// encoded bytes at all and a raw-only replace silently does nothing: the
// golden would fail and the width would stay path-dependent, on the one
// platform least likely to be running this suite. `windows` is a shipped
// release target (`.goreleaser.yaml`), even though CI is Linux-only, so the
// escaped form is folded too rather than left as a latent trap.
//
// Folding narrows nothing. Byte-identity is still asserted over every byte
// of the record, `model` included, in the RELATIVE spelling `0023:MVV` and
// both A1/A2 spikes name as normative; only the prefix that differs between
// two checkouts of the same commit is removed.
// foldEmitDispositions0024 removes RDR 0024 `0024:C4`'s appended
// `dispositions` key from an emitted record, so this comparison keeps
// measuring the thing it exists to measure.
//
// The golden above is a PRE-CHANGE capture: it was taken from a
// `git stash`-clean tree before RDR 0023's struct edit, and that
// provenance is the whole evidence. It cannot be re-captured — no
// post-0023 tree can produce a pre-0023 binary's bytes — so REGENERATING
// it would not update the fixture, it would destroy the only pre-change
// side the battery has and silently convert this assertion into another
// new-build-against-itself comparison, which REQ-97 names as the exact
// hole it exists to close.
//
// 0024 appends `dispositions` immediately after `emit`, a licensed
// additive append this record does not own. Folding it out is the same
// move `foldCheckoutRoot` makes for the checkout prefix and narrows
// nothing 0023 claims: byte-identity is still asserted over every other
// key, in order, and the fold itself is EXACT — it accepts the EMPTY
// OBJECT and nothing else, and only where the key sits immediately after
// `emit`.
//
// The empty-object literal is the whole point of the narrowing. A needle
// spelled `\{[^{}]*\}` on the `dispositions` half excludes only NESTED
// braces, not content, so it would erase `{"verdict":"route"}` exactly as
// silently as it erases `{}` — and this test is REQ-97's only pre-change
// comparison, so a default-mode regression that started POPULATING
// `dispositions` would be absorbed by the fold and seen nowhere else in
// the battery. `0024:C4` ships the key as `{}` in default mode (D7: all 28
// regenerated sites took `"dispositions":{}` exactly), so a populated
// value here is a contract move to be adjudicated, never folded away. The
// `emit` half stays `[^{}]*` — its value legitimately varies and is not
// this fold's subject.
//
// Returning an error rather than calling `t.Fatalf` keeps the three
// outcomes — folded, key absent, key populated — drivable from a test.
func foldEmitDispositions0024(line string) (string, error) {
	// The needle anchors on the PREDECESSOR key, so the fold cannot
	// silently succeed on a record where `dispositions` drifted to some
	// other position — which is the one thing `0024:C4` fixes about it.
	appended := regexp.MustCompile(`("emit":\{[^{}]*\}),"dispositions":\{\}`)
	folded := appended.ReplaceAllString(line, "$1")
	if folded == line {
		if populated0024.MatchString(line) {
			return "", fmt.Errorf("the emitted record carries a NON-EMPTY "+
				"`dispositions` value after `emit`; `0024:C4` ships the key "+
				"as `{}` in default mode, and this comparison is REQ-97's "+
				"only pre-change side — \"S1-S5 all compare the new build "+
				"against itself, so none of them can see a default-mode "+
				"regression that moves both sides together\". Folding a "+
				"populated value would absorb exactly that regression and "+
				"leave it asserted nowhere, so it is reported instead:\n%s",
				line)
		}
		return "", fmt.Errorf("the emitted record carries no `dispositions` "+
			"key immediately after `emit`; `0024:C4` appends one to every "+
			"resolve payload, present as `{}` and never omitted, so its "+
			"absence here means the append regressed rather than that this "+
			"fold is unnecessary:\n%s", line)
	}
	if strings.Contains(folded, "dispositions") {
		return "", fmt.Errorf("the record carries `dispositions` more than "+
			"once; the fold removes exactly the one appended key:\n%s", line)
	}
	return folded, nil
}

// populated0024 discriminates the two ways the exact fold can miss: the
// key is absent after `emit` (the append regressed) versus present but
// carrying a value the fold refuses (a default-mode regression). It is
// deliberately NOT the fold's own needle — it matches any brace-free run,
// including the nested-brace case the fold's `[^{}]*` half would also
// decline, so the diagnostic never claims absence for a key that is there.
var populated0024 = regexp.MustCompile(`"emit":\{[^{}]*\},"dispositions":`)

func foldCheckoutRoot(t *testing.T, line string) string {
	t.Helper()

	raw := repoRootFor(t) + "/"

	esc, err := encodedCheckoutPrefix(raw)
	if err != nil {
		t.Fatalf("encode the checkout prefix: %v", err)
	}

	line = strings.ReplaceAll(line, esc, "")
	if esc != raw {
		line = strings.ReplaceAll(line, raw, "")
	}
	return line
}

// encodedCheckoutPrefix returns the prefix as it appears once encoded into
// the record.
//
// Encoded with the SAME encoder settings the record itself is written with:
// `respond.marshalCanonical` sets `SetEscapeHTML(false)`, so a `<`, `>` or
// `&` in the path renders as itself on the wire. `json.Marshal` escapes all
// three (`&` as `\u0026`), which would make the needle a string that does
// not occur in the record at all — and on a platform whose separator also
// needs escaping the raw fallback would miss too, folding nothing while
// still reporting success.
//
// `Encode` returns the value quoted and newline-terminated; both are
// trimmed to leave the escaped body.
func encodedCheckoutPrefix(raw string) (string, error) {
	var buf strings.Builder
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(raw); err != nil {
		return "", err
	}
	return strings.Trim(strings.TrimRight(buf.String(), "\n"), `"`), nil
}

// The fold's needle must be byte-identical to how a REAL emitted record
// spells the same prefix.
//
// Driven end to end rather than against a locally re-encoded expectation:
// the record below is produced by the shipped emit path, so the arms fail
// if the fold's encoder settings and the wire's diverge — which is exactly
// the divergence `json.Marshal` introduced. The arms are the paths on which
// the two encoders disagree: an HTML metacharacter, which `json.Marshal`
// escapes (`&` as `\u0026`) and the wire encoder does not.
//
// `<` and `>` are RESERVED characters in a Windows filename, so the arms
// carrying them cannot create their fixture there and are skipped rather
// than failed. Skipping is safe because the discrimination does not depend
// on them: `&` is legal in a Windows path and `json.Marshal` escapes it too
// (as `\u0026`), so the `ampersand` arm keeps this oracle red under the
// reverted encoder on EVERY platform. Failing instead would take the whole
// suite red on the one shipped target (`.goreleaser.yaml`) whose separator
// escaping is the reason `foldCheckoutRoot` folds two spellings at all.
func TestFoldCheckoutRootEncodesThePrefixTheWayTheRecordDoes(t *testing.T) {
	for _, tc := range []struct {
		name string
		dir  string
	}{
		{name: "plain", dir: "checkout"},
		{name: "ampersand", dir: "a&b"},
		{name: "angles", dir: "less<greater>"},
		{name: "all-three", dir: "a&b<c>d"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if runtime.GOOS == "windows" &&
				strings.ContainsAny(tc.dir, "<>") {
				t.Skipf("`<` and `>` are reserved in a Windows filename, so "+
					"%q is not creatable here; the `ampersand` arm carries "+
					"the same HTML-escaping discrimination on this platform",
					tc.dir)
			}

			// The normative pricing table, relocated under a directory
			// spelled the way a hostile checkout root would be.
			src, err := os.ReadFile(pricingModelPath(t))
			if err != nil {
				t.Fatalf("read the pricing fixture: %v", err)
			}
			root := filepath.Join(t.TempDir(), tc.dir)
			if err := os.MkdirAll(root, 0o750); err != nil {
				t.Fatalf("make the checkout-shaped dir: %v", err)
			}
			model := filepath.Join(root, "pricing-decision-table.toml")
			if err := os.WriteFile(model, src, 0o600); err != nil {
				t.Fatalf("write the relocated fixture: %v", err)
			}

			// A real record, from the shipped emit path. `model` echoes the
			// `--model` argument verbatim, so the record carries this
			// directory in the wire's own spelling.
			line := requireSuccess(t, resolveArgs(model, "", "decide",
				"--tag", "tier=paid", "--tag", "region=eu", "--as=json")...)

			esc, err := encodedCheckoutPrefix(root + string(os.PathSeparator))
			if err != nil {
				t.Fatalf("encode the checkout prefix: %v", err)
			}
			if !strings.Contains(line, esc) {
				t.Errorf("the folded needle does not occur in the emitted "+
					"record:\n  needle = %s\n  record = %s\n"+
					"the fold must encode the prefix with the RECORD's "+
					"encoder settings (SetEscapeHTML(false)); a needle that "+
					"HTML-escapes `&`/`<`/`>` matches nothing and "+
					"foldCheckoutRoot silently folds nothing", esc, line)
			}
		})
	}
}

// D8 records that RDR 0023's pre-change golden is FOLDED rather than
// regenerated, and states the standard the fold has to meet: "The fold is
// exact rather than permissive, so it strengthens rather than weakens the
// site." This oracle holds the fold to that word on the one axis a
// position-anchored needle does not cover by itself — the VALUE.
//
// A `dispositions` half spelled `\{[^{}]*\}` excludes nested braces, not
// content, so it erases `{"verdict":"route"}` verbatim. Neither of the
// fold's guards catches that: the absence guard fires only on a total miss
// and the duplicate guard only on a second key. The populated value would
// therefore be folded silently, inside the battery's ONLY pre-change
// comparison (REQ-97: "S1-S5 all compare the new build against itself, so
// none of them can see a default-mode regression that moves both sides
// together"). That is the regression this oracle is red for.
//
// The inputs are hand-built records rather than emitted ones on purpose:
// the populated arm cannot be produced by the shipped default-mode path at
// all, which is exactly why it must be asserted here.
func TestD8_0024_TheDispositionsFoldAcceptsTheEmptyObjectAndNothingElse(t *testing.T) {
	const emitted = `{"model":"m","emit":{"columns":["verdict"]}`

	for _, tc := range []struct {
		name    string
		line    string
		want    string
		wantErr string
	}{
		{
			// `0024:C4`'s shipped default-mode shape (D7: all 28
			// regenerated sites took `"dispositions":{}` exactly).
			name: "empty-object-folds",
			line: emitted + `,"dispositions":{},"next":"n"}`,
			want: emitted + `,"next":"n"}`,
		},
		{
			// A default-mode regression that started populating the key.
			// Folding this absorbs the very thing REQ-97 exists to see.
			name:    "populated-object-is-refused",
			line:    emitted + `,"dispositions":{"verdict":"route"},"next":"n"}`,
			wantErr: "NON-EMPTY",
		},
		{
			// The append itself regressed.
			name:    "absent-key-is-refused",
			line:    emitted + `,"next":"n"}`,
			wantErr: "no `dispositions` key",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := foldEmitDispositions0024(tc.line)

			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("the fold refused `0024:C4`'s shipped shape, "+
						"which it exists to accept: %v", err)
				}
				if got != tc.want {
					t.Errorf("the fold removed the wrong span:\n"+
						"  want = %s\n  got  = %s", tc.want, got)
				}
				return
			}

			if err == nil {
				t.Fatalf("the fold ACCEPTED a record it must refuse and "+
					"returned %s\n\nThe fold is the pre-change golden's "+
					"only filter, and this comparison is the battery's only "+
					"pre-change side. A value it swallows is a regression "+
					"no other oracle can see.", got)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("the refusal does not name what moved:\n"+
					"  want substring = %q\n  got            = %v",
					tc.wantErr, err)
			}
		})
	}
}

// The populated-value refusal must name REQ-97, because REQ-97 is the hole
// the golden exists to close and the reason the value cannot be folded: it
// is the requirement that says no other oracle in the battery can see a
// default-mode regression. A message that reports only "unexpected value"
// leaves the reader to rediscover why an exact fold is not pedantry.
func TestD8_0024_ThePopulatedDispositionsRefusalNamesReq97sHole(t *testing.T) {
	_, err := foldEmitDispositions0024(
		`{"emit":{"columns":["verdict"]},"dispositions":{"verdict":"route"}}`)
	if err == nil {
		t.Fatal("the fold accepted a populated `dispositions` value")
	}

	// REQ-97's own words, as quoted in this file's header.
	const hole = "compare the new build against itself"
	if !strings.Contains(err.Error(), hole) {
		t.Errorf("the refusal does not carry REQ-97's hole:\n"+
			"  want substring = %q\n  got            = %v\n"+
			"`0024:C4` ships `dispositions` as `{}` in default mode, so a "+
			"populated value means either the append or the default-mode "+
			"contract moved — both adjudicated, neither folded.", hole, err)
	}
}
