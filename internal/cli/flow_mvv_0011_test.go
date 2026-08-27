package cli

// RDR 0011 — the Minimum Viable Validation, steps 2 through 9, runnable
// end to end.
//
// The MVV's discriminating figure is the 21→3 narrowing over
// `models/rdr.toml` at `stage=resolved`: a stripped-match build reports 21,
// C1's predicate reports 3, and `--all` must still report the 21. That
// witness is one model's, never proof of the general claim — the RDR says
// so itself ("the `21→3` figure is a witness for one model") and closes the
// gap by ENUMERATION instead, through the `disposition` table's per-class
// oracles. Those live beside this file; this one runs the witness.
//
// No round-trip / inverse invariant is asserted here, and that is the
// record's own reading: the mini-check cue register says "Cues absent:
// round-trip / fidelity — this RDR defines no import/export, parse/deparse,
// or serialize inverse." The invariant this MVV pins instead is the
// CANDIDATE SET's exact membership, both directions: every row the state
// can take is reported, and every row it excludes is absent.

import (
	"os/exec"
	"slices"
	"strings"
	"testing"
)

// mvvRdrArtifact seeds `models/rdr.toml`'s owned state to
// `stage=resolved, status=draft, gate_passed=false` — MVV step 2's artifact
// — through the production write path, and returns its binding.
func mvvRdrArtifact(t *testing.T) (model, bind string) {
	t.Helper()

	model = repoModelPath(t)
	art := newFlowArtifact(t, "rdr.artifact")
	bind = artifactBinding("rdr", art)
	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", bind, "--write", "stage=resolved",
		"--write", "status=draft", "--write", "gate_passed=false",
		"--as=json")
	return model, bind
}

// REQ-87 / MVV 3: "Assert: `candidates[]` is exactly `prelock`,
// `resolve-route-back`, `resolve-abandon` (3 rows, one per outcome), and no
// rule whose `match.stage` names another value; `outcomes[]` is still the
// model's full alphabet."
// REQ-95 / `0011:S1`: the same, as the first testing scenario.
// A-10: "no rule whose `match.stage` names another value" is an assertion
// about the REPORTED SET, not about the model file — "the three named rules
// are the whole of `candidates[]` and every other rule of the 22 is absent,
// which is the discriminating form (a stripped-match build reports 21)."
// HAPPY PATH — the MVV's witness.
func TestReq87And95_MVV3NarrowsTheCheckedInModelFrom21RowsToExactlyThree(t *testing.T) {
	model, bind := mvvRdrArtifact(t)

	data := runNext(t, model, []string{bind})

	want := []string{"prelock", "resolve-abandon", "resolve-route-back"}
	got := candidateRules(t, data)
	if !slices.Equal(got, want) {
		t.Errorf("candidates = %v (%d rows); want exactly %v (3 rows, one "+
			"per declared outcome)\nThis is the MVV's discriminating "+
			"figure: a stripped-match build reports 21 of the 22 rules, "+
			"and `grep -n 'eq = \"resolved\"' models/rdr.toml` finds "+
			"exactly three `[rule.match.stage]` hits. Every other rule is "+
			"ABSENT — the assertion is about the REPORTED SET, which is "+
			"what discriminates", got, len(got), want)
	}

	// `outcomes[]` is STILL the model's full alphabet: narrowing the
	// candidate list never narrows what may be requested.
	outcomes, ok := stringsAt(data, "outcomes")
	if !ok {
		t.Fatalf("`outcomes` is not an array of bare tags: %#v",
			data["outcomes"])
	}
	slices.Sort(outcomes)
	if wantOutcomes := []string{"abandon", "advance", "revise"}; !slices.Equal(outcomes, wantOutcomes) {
		t.Errorf("outcomes = %v; want the full declared alphabet %v — the "+
			"alphabet is what may be REQUESTED, the candidates are what "+
			"varies with state", outcomes, wantOutcomes)
	}

	// The narrowed set names the same three outcomes the alphabet carries:
	// `models/rdr.toml` authors one ordinary rule per (stage, outcome), so
	// what the caller GAINS here is the per-outcome ROW instead of 21 rows
	// to sift — not a smaller outcome set. Stated so the witness is not
	// read as more than it is.
	seen := map[string]bool{}
	for _, c := range nextCandidates(t, data) {
		if o, _ := c["outcome"].(string); o != "" {
			seen[o] = true
		}
	}
	if len(seen) != 3 {
		t.Errorf("the three candidates cover %d outcomes; want 3, one per "+
			"declared outcome", len(seen))
	}
}

// REQ-88 / MVV 4: "Re-run with `--all`; assert `candidates[]` has the 21
// rows the baseline reported (the 22nd, `finalize-pass`, is guard-excluded
// on `gate_passed`, not match-excluded) and the same `outcomes[]`."
// BOUNDARY — the `--all` half of the witness. The 21 is a property of
// `gate_passed=false`, not of `stage=resolved`, so `finalize-pass` stays
// excluded in BOTH modes and the two exclusions do not confound.
func TestReq88_MVV4UnderAllTheSameModelReportsThe21RowsTheBaselineDid(t *testing.T) {
	model, bind := mvvRdrArtifact(t)

	data := runNext(t, model, []string{bind}, "--all")
	got := candidateRules(t, data)

	if len(got) != 21 {
		t.Errorf("`--all` reports %d candidates; want 21 — the rows the "+
			"BASELINE (stripped-match) build reported. `--all` restores "+
			"0005's enumeration exactly: every non-escape row whose guard "+
			"the kernel does not decide false.\ngot = %v", len(got), got)
	}

	// The 22nd is guard-excluded on `gate_passed`, not match-excluded — so
	// it is absent under `--all` too, and the two mechanisms do not
	// confound. Without this arm the count above could be satisfied by a
	// build that excluded the wrong row.
	if slices.Contains(got, "finalize-pass") {
		t.Errorf("`finalize-pass` is reported under --all; its " +
			"`[rule.guard.all.gate_passed] eq = \"true\"` is decided FALSE " +
			"over `gate_passed=false`, and C1 leaves guard exclusions " +
			"untouched")
	}
	// And its guard-partitioned sibling IS reported, so the absence above
	// is the guard biting rather than the row being unreachable.
	if !slices.Contains(got, "finalize-blocked") {
		t.Errorf("`finalize-blocked` is absent under --all: %v\nIt is "+
			"`finalize-pass`'s guard-partitioned sibling at the same "+
			"(stage, outcome), so its presence is what proves the exclusion "+
			"above is the GUARD and not the row", got)
	}

	// Same alphabet in both modes.
	def := runNext(t, model, []string{bind})
	defOutcomes, _ := stringsAt(def, "outcomes")
	allOutcomes, _ := stringsAt(data, "outcomes")
	slices.Sort(defOutcomes)
	slices.Sort(allOutcomes)
	if !slices.Equal(defOutcomes, allOutcomes) {
		t.Errorf("outcomes differ across modes: default %v, --all %v",
			defOutcomes, allOutcomes)
	}
}

// REQ-89 / MVV 5: "Run over the 0005 gated fixture (`flowGatedNextModel` …)
// with and without `--all`; assert the `gated-excluded` row is absent in
// both (guard-excluded), paired with its negative control `gated-reported`
// present in both. Then over `flowMVVModel` (`status=draft`), assert the
// row C3 adds … is absent by default and present under `--all`."
// REQ-96: the fixture that AUTHORS the row is where the assertion belongs —
// "asserted over `flowMVVModel`, which carries no such row, the oracle
// passes vacuously in every possible build and pins nothing."
// REQ-61: C3's discriminating fixtures — "a row whose match key is present
// and UNEQUAL (excluded by default, reported under --all)".
// DOMAIN EDGE
func TestReq61And89And96_MVV5GuardExclusionsAndTheAddedDiscriminatingRow(t *testing.T) {
	t.Run("guard-excluded-over-the-fixture-that-authors-it", func(t *testing.T) {
		model := writeFlowModel(t, flowGatedNextModel)
		bind := artifactBinding(flowStateRole,
			seedArtifact(t, model, "status=draft", "flag=true"))

		for _, mode := range []string{"default", "all"} {
			var extra []string
			if mode == "all" {
				extra = []string{"--all"}
			}
			data := runNext(t, model, []string{bind}, extra...)

			if c := candidateNamed(t, data, "gated-excluded"); c != nil {
				t.Errorf("[%s] `gated-excluded` is reported: %#v — it is "+
					"GUARD-excluded, which C1 leaves untouched", mode, c)
			}
			// The negative control, in the same mode.
			requireCandidate(t, data, "gated-reported")
		}
	})

	t.Run("added-row-absent-by-default-present-under-all", func(t *testing.T) {
		// C3's discriminating shape: a row whose `match` key is PRESENT
		// and names a value other than the one the state holds. `phase`
		// is supplied at `alpha`, and `match-beta` names `beta`.
		model := writeFlowModel(t, flowMatchClassesModel)
		bind := seedMatchArtifact(t, model, "status=draft")

		def := runNext(t, model, []string{bind}, "--tag", "phase=alpha")
		if c := candidateNamed(t, def, "match-beta"); c != nil {
			t.Errorf("the added row is reported by DEFAULT: %#v\nIts match "+
				"key is PRESENT and UNEQUAL, so the kernel refuses its "+
				"probe `no_match` — the one disposition that excludes", c)
		}

		all := runNext(t, model, []string{bind}, "--tag", "phase=alpha", "--all")
		requireCandidate(t, all, "match-beta")
	})
}

// REQ-90 / MVV 6: "Run over a model exercising the three reachable match
// cases … and a dead row … Then add `--tag key=<value>` and assert the list
// narrows to the rows whose match holds. Assert the three pins"
// REQ-63: "It MUST also add the three pins that license C1's presence test:
// a model declaring one owned key served by two readers is refused at load
// …; a repeated --tag key is refused `flow-tag-duplicate`; and a --tag on a
// key the model declares owned is refused `flow-tag-owned` — all before any
// probe is built."
// REQ-106 / `0011:S5`: the same three, with the loader's own message.
// ADVERSARIAL — the invariants C1's exactness rests on. "C3 pins them so a
// relaxation fails a named oracle rather than silently reopening the fold."
func TestReq63And90And106_MVV6TheThreePinsThatLicenseThePresenceTest(t *testing.T) {
	t.Run("two-readers-on-one-owned-key-is-refused-at-load", func(t *testing.T) {
		model := writeFlowModel(t, flowTwoReadersOneKeyModel)

		ce := requireRefusal(t, "flow-model-invalid", 2,
			"flow", "next", "--model", model, "--as=json")

		// The refusal is `checkAccessorBindings`', which counts key
		// OCCURRENCES — so one reader declaring a key twice fails the same
		// test. Asserting the CATEGORY rather than the message keeps the
		// oracle on the contract rather than on the loader's wording.
		if ce.Code != "flow-model-invalid" {
			t.Errorf("code = %q; want `flow-model-invalid`", ce.Code)
		}
		// It happens at LOAD, before any probe exists: a build that
		// reached the kernel with a conflicted key would fold it into a
		// silent `no_match`, which is the drop P1 forbids.
		if strings.Contains(ce.Message, "candidate") {
			t.Errorf("the refusal mentions candidates (%q); it precedes any "+
				"probe", ce.Message)
		}
	})

	t.Run("repeated-tag-key-is-refused-flow-tag-duplicate", func(t *testing.T) {
		model := writeFlowModel(t, flowMatchClassesModel)
		bind := seedMatchArtifact(t, model, "status=draft")

		requireRefusal(t, "flow-tag-duplicate", 2,
			"flow", "next", "--model", model, "--artifact", bind,
			"--tag", "phase=alpha", "--tag", "phase=beta", "--as=json")
	})

	t.Run("tag-on-an-owned-key-is-refused-flow-tag-owned", func(t *testing.T) {
		model := writeFlowModel(t, flowMatchClassesModel)
		bind := seedMatchArtifact(t, model, "status=draft")

		// `status` is declared OWNED, so the observed channel is refused —
		// which is also why an owned match key has NO remedy within the
		// invocation and why the demand-set term is the fix rather than
		// `--tag`.
		requireRefusal(t, "flow-tag-owned", 2,
			"flow", "next", "--model", model, "--artifact", bind,
			"--tag", "status=draft", "--as=json")
	})
}

// REQ-64: "There is no conflicted-key fixture for flow next: the case is
// unreachable (A3), and a fixture that reaches the kernel's conflicted fold
// by constructing a `TagSet` in-package tests `internal/resolve`, which this
// RDR does not change."
// REQ-105 / `0011:S4`: "An AUTHORED empty `[rule.match]` is not a fixture
// here and MUST NOT be attempted: `internal/table/normalize.go` refuses it
// at load … the empty match pattern C1 names is minted by the probe builder
// at runtime, never authored."
// REQ-66: "the other three producers are unreachable and MUST NOT be
// fixtured".
// ADVERSARIAL — three NEGATIVE obligations. Each names a fixture the build
// must NOT contain, and the oracle is a scan of this RDR's own fixture
// corpus, because a negative obligation about fixtures is only assertable
// against the fixtures.
func TestReq64And66And105_TheForbiddenFixtureShapesAreAbsentFromTheCorpus(t *testing.T) {
	corpus := map[string]string{
		"flowMatchClassesModel":       flowMatchClassesModel,
		"flowSetKindedMatchModel":     flowSetKindedMatchModel,
		"flowMixedStateRowModel":      flowMixedStateRowModel,
		"flowAbsentMatchKeyModel":     flowAbsentMatchKeyModel,
		"flowFullyResolvedModel":      flowFullyResolvedModel,
		"flowLooseWriterModel":        flowLooseWriterModel,
		"flowUncomparableGuardModel":  flowUncomparableGuardModel,
		"flowOwnedUnavailableModel":   flowOwnedUnavailableModel,
		"flowSortOrderModel":          flowSortOrderModel,
		"flowMatchOnlyOwnedModel":     flowMatchOnlyOwnedModel,
		"flowMatchOnlyOwnedSoloModel": flowMatchOnlyOwnedSoloModel,
		"flowZeroOwnedMatchModel":     flowZeroOwnedMatchModel,
		"flowSideWriterModel":         flowSideWriterModel,
	}

	for name, src := range corpus {
		// REQ-105: no AUTHORED empty match block. The loader refuses it
		// (`rule <id> carries no match block`, `malformed_rule_shape`),
		// keying on key PRESENCE — so an absent and a present-but-empty
		// block are equally malformed. Every model here must LOAD, which
		// is the enforcing property.
		if strings.Contains(src, "[rule.match]\n[") ||
			strings.HasSuffix(strings.TrimSpace(src), "[rule.match]") {
			t.Errorf("%s authors an EMPTY `[rule.match]` block; it is "+
				"refused at load and MUST NOT be attempted — C1's empty "+
				"match pattern is minted by the probe builder at RUNTIME, "+
				"never authored", name)
		}
	}

	// REQ-64: no conflicted-key fixture. The case is unreachable from the
	// CLI over any model that loads, and reaching the kernel's conflicted
	// fold takes an in-package `TagSet` construction — which tests
	// `internal/resolve`, the package this RDR does not change. The
	// assertable form: this RDR's own test files construct no
	// `resolve.TagSet`.
	for _, file := range []string{
		"internal/cli/flow_next_0011_test.go",
		"internal/cli/flow_all_0011_test.go",
		"internal/cli/flow_unknown_0011_test.go",
		"internal/cli/flow_demand_0011_test.go",
		"internal/cli/flow_fixtures_0011_test.go",
		"internal/cli/flow_mvv_0011_test.go",
		"internal/cli/flow_rehome_0011_test.go",
	} {
		src := readRepoFile(t, file)
		// Spelled in halves so this scan does not match itself.
		if strings.Contains(src, "resolve.TagSet"+"{") {
			t.Errorf("%s constructs a `resolve.TagSet`; a fixture that "+
				"reaches the kernel's conflicted fold that way tests "+
				"`internal/resolve`, which this RDR does not change", file)
		}
	}
}

// REQ-91 / MVV 7: "Assert `unknown` is present as `[]` rather than omitted
// on a candidate with nothing undecided, in both modes — the fixture must
// SUPPLY one … and which is run WITH `--evaluate-gates`"
// REQ-112: "Assert a guard atom the seam cannot decide over a present key
// reports `{key, uncomparable}` (A13)."
// BOUNDARY — the two halves of MVV 7.
func TestReq91_MVV7UnknownIsEmptyOnAFullyResolvedCandidateInBothModes(t *testing.T) {
	model := writeFlowModel(t, flowFullyResolvedModel)
	bind := artifactBinding(flowStateRole,
		seedArtifact(t, model, "status=draft", "flag=true"))

	for _, mode := range []string{"default", "all"} {
		extra := []string{"--evaluate-gates"}
		if mode == "all" {
			extra = append(extra, "--all")
		}
		data := runNext(t, model, []string{bind}, extra...)

		c := requireCandidate(t, data, "resolved-row")
		if _, held := c["unknown"]; !held {
			t.Errorf("[%s] the candidate carries no `unknown` field; it is "+
				"emitted as `[]` rather than OMITTED, so ONE consumer "+
				"struct parses either mode. keys = %v", mode, keysOf(c))
			continue
		}
		if pairs := candidateUnknown(t, c); len(pairs) != 0 {
			t.Errorf("[%s] `unknown` = %#v on a candidate with nothing "+
				"undecided\nThe fixture SUPPLIES this shape because no "+
				"shipped 0005 fixture produces a fully-resolved candidate: "+
				"the row's match and guard keys are all established by the "+
				"invoked reader, that reader also serves the owned key it "+
				"writes, and the model declares no gate", mode, pairs)
		}
	}
}

// REQ-92 / MVV 8: "Run over a model declaring an owned key some row MATCHES
// on but no row writes, clears, or guards (A15): assert its reader appears
// in `readers`, the key is in the view, and the row is match-decided …
// assert the same `readers` set under `--all`."
// HAPPY PATH — MVV 8's first half, run end to end as the MVV states it.
// (The `flow resolve` half is REQ-93 and lives beside the S8 oracles.)
func TestReq92_MVV8TheMatchOnlyOwnedKeyIsReadAndDecidedInBothModes(t *testing.T) {
	model := writeFlowModel(t, flowMatchOnlyOwnedSoloModel)

	state := newFlowArtifact(t, "state.artifact")
	stateBind := artifactBinding(flowMatchRole, state)
	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", stateBind, "--write", "status=draft", "--as=json")

	side := newFlowArtifact(t, "side.artifact")
	sideBind := artifactBinding(flowSideRole, side)
	requireSuccess(t, "flow", "set-state", "--model",
		writeFlowModel(t, flowSideWriterModel), "--artifact", sideBind,
		"--write", "mode=fast", "--as=json")

	binds := []string{stateBind, sideBind}

	def := runNext(t, model, binds)
	if readers := readersOf(t, def); !slices.Contains(readers, "side") {
		t.Errorf("readers = %v; `read.side` serves `mode`, which `mode-row` "+
			"MATCHES on and no row writes, clears, or guards", readers)
	}
	if !slices.Contains(viewKeysOf(t, def), "mode") {
		t.Errorf("the assembled view's keys = %v; `mode` must be IN the "+
			"view for the predicate to read it", viewKeysOf(t, def))
	}
	c := requireCandidate(t, def, "mode-row")
	if pairs := candidateUnknown(t, c); hasUnknown(pairs, "mode", "absent") {
		t.Errorf("`mode-row` reports {mode, absent}: %#v — it is "+
			"MATCH-DECIDED, not undecided. Every row coming back a "+
			"candidate carrying that key absent is the un-extended set's "+
			"answer, the default silently degraded to `--all`", pairs)
	}

	all := runNext(t, model, binds, "--all")
	if a, b := readersOf(t, def), readersOf(t, all); !slices.Equal(a, b) {
		t.Errorf("readers under default = %v, under --all = %v; the demand "+
			"set is mode-independent", a, b)
	}
}

// REQ-94 / MVV 9: "`make check` passes; `go test ./internal/resolve` passes
// with no test changed and `git diff --stat internal/resolve` empty; the
// 0005 `next` suite passes with only C3's five mechanical
// `unresolved`→`unknown` re-homings (ok-bool asserted), and the file header
// names this RDR."
// REQ-119 / `0011:S8`: "Done: S1–S8 green and `make check` passes."
// ADVERSARIAL — the two mechanically checkable halves. `make check` itself
// is the build's gate rather than an oracle, so what is asserted here is
// the kernel's untouchedness in both of the forms MVV 9 names.
func TestReq94And119_MVV9TheKernelIsUntouchedInBothOfTheFormsTheMVVNames(t *testing.T) {
	// (i) `git diff --stat internal/resolve` empty — the ONE mechanically
	// checkable form of "the kernel is untouched".
	if diff := strings.TrimSpace(gitDiffStat(t, "internal/resolve")); diff != "" {
		t.Errorf("`git diff --stat internal/resolve` is NOT empty:\n%s",
			diff)
	}

	// (ii) NO test file under `internal/resolve` changed. A green kernel
	// suite is corroboration, "not a second criterion, since a green suite
	// cannot by itself prove nothing changed" — so this asserts the DIFF
	// of the test files rather than running them.
	root := repoRootFor(t)
	names, err := exec.Command("git", "-C", root, "diff", "--name-only",
		diffBase(t), "--", "internal/resolve").Output()
	if err != nil {
		t.Skipf("`git diff --name-only` failed: %v", err)
	}
	for _, line := range strings.Split(strings.TrimSpace(string(names)), "\n") {
		if line == "" {
			continue
		}
		t.Errorf("`%s` changed; MVV 9 requires `go test ./internal/resolve` "+
			"to pass with NO test changed. `TestReq78_MatchPatternStill"+
			"FoldsAbsenceIntoNonMatch`, the `match_conflicted_test.go` "+
			"oracles, and the reason-set pins `TestReq49`/`TestReq47` stay "+
			"green AS SHIPPED — a kernel diff in the build is a defect, not "+
			"an expected inversion", line)
	}
}
