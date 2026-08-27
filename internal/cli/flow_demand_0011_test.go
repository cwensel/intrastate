package cli

// RDR 0011 — the demand-set term, and the `flow resolve` behaviour-change
// class it carries.
//
// C1's predicate is only real if the assembled view actually carries the
// keys it reads, so `internal/cli/flow_exec.go::invokedReaders` gains a
// third term: each row's MATCH-block owned keys, beside `RequiresOwned` and
// the guard-owned keys it already unions. Without it, an owned key a model
// MATCHES on but never writes, clears, or guards demands no reader, is
// absent from the view, and has its atom omitted from EVERY probe — so
// every row becomes a candidate carrying `{key, absent}` and the default
// silently degrades to `--all` for that model.
//
// The term binds BOTH callers of `invokedReaders`, so it changes
// `flow resolve` too. This RDR does not deny that: it NAMES a behaviour
// change class and ACCEPTS it, and S8 pins the class "so the DEV-8 class is
// a known cost, not a surprise." The oracles below are written so the cost
// is asserted rather than discovered.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
)

// sideBindings returns the state and side artifact bindings for the
// match-only-owned-key fixture, seeding `status` through the CLI and
// leaving `mode`'s reader either bound or unbound as asked.
func sideBindings(t *testing.T, model string, bindSide bool, sideWrites ...string) []string {
	t.Helper()

	state := newFlowArtifact(t, "state.artifact")
	stateBind := artifactBinding(flowMatchRole, state)
	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", stateBind, "--write", "status=draft", "--as=json")

	if !bindSide {
		return []string{stateBind}
	}

	side := newFlowArtifact(t, "side.artifact")
	sideBind := artifactBinding(flowSideRole, side)
	if len(sideWrites) > 0 {
		// The `side` role serves `mode` and `extra`, both owned. The model
		// declares no writer for that role, so the value is established
		// through the SAME artifact under a writer-carrying sibling.
		args := []string{"flow", "set-state", "--model",
			writeFlowModel(t, flowSideWriterModel), "--artifact", sideBind}
		requireSuccess(t, append(args, flattenWrites(sideWrites)...)...)
	}
	return []string{stateBind, sideBind}
}

func flattenWrites(writes []string) []string {
	out := make([]string, 0, len(writes)*2+1)
	for _, w := range writes {
		out = append(out, "--write", w)
	}
	return append(out, "--as=json")
}

// REQ-16: "The assembled view MUST actually carry the keys the predicate
// reads, so `internal/cli/flow_exec.go::invokedReaders` MUST add each row's
// MATCH-block owned keys to its demand set, which today is
// `Row.RequiresOwned` unioned with the row's GUARD-owned keys only"
// REQ-114 / `0011:S8`: "the reader IS invoked (it appears in the payload's
// `readers`), the key IS in the view, and the row is match-decided rather
// than reported `{key, absent}`."
// ADVERSARIAL — the discriminating oracle for the extension. "On the
// un-extended set the reader is never invoked and every row comes back a
// candidate carrying that key `absent` — the default silently degraded to
// `--all`."
func TestReq16And114_AMatchOnlyOwnedKeyInvokesItsReaderAndIsMatchDecided(t *testing.T) {
	model := writeFlowModel(t, flowMatchOnlyOwnedModel)
	binds := sideBindings(t, model, true, "mode=fast")

	data := runNext(t, model, binds)

	// (i) the reader IS invoked.
	if readers := readersOf(t, data); !slices.Contains(readers, "side") {
		t.Errorf("readers = %v; `read.side` serves `mode`, an owned key "+
			"`mode-row` MATCHES on and no row writes, clears, or guards. "+
			"Without the match-owned term it is never invoked", readers)
	}

	// (ii) the key IS in the view.
	if owned, ok := data["owned"].(map[string]any); !ok || owned["mode"] == nil {
		t.Errorf("`owned` does not carry `mode`: %#v\nThe demand set is "+
			"what decides which readers run, and a key no reader "+
			"established cannot be in the view", data["owned"])
	}

	// (iii) the row is MATCH-DECIDED rather than reported `{mode, absent}`.
	c := requireCandidate(t, data, "mode-row")
	if pairs := candidateUnknown(t, c); hasUnknown(pairs, "mode", "absent") {
		t.Errorf("`mode-row` reports {mode, absent}: %#v\nThat is the "+
			"un-extended set's answer — the default degraded to `--all`, "+
			"because the CLI omitted the atom from every probe for want of "+
			"a reader it never invoked", pairs)
	}

	// And it genuinely decides: at the other declared value the row is
	// excluded, which the un-extended set can never produce.
	other := runNext(t, writeFlowModel(t, flowMatchOnlyOwnedModel),
		sideBindings(t, writeFlowModel(t, flowMatchOnlyOwnedModel), true, "mode=slow"))
	if c := candidateNamed(t, other, "mode-row"); c != nil {
		t.Errorf("`mode-row` survives `mode=slow` against `eq \"fast\"`: "+
			"%#v\nWithout this arm the assertion above passes against a "+
			"build that invokes the reader but still ignores the match "+
			"pattern", c)
	}
}

// REQ-18: "The demand set is a property of the MODEL (and, for `resolve`,
// the requested outcome), not of the mode: it is computed once per
// invocation, before the row loop, and is identical under --all"
// REQ-115 / `0011:S8`: "Negative controls: the same model under `--all`
// reports the SAME reader set …; and a model with zero owned tags demands
// nothing and invokes no reader, unchanged"
// REQ-21: "No `flow resolve` run over a model WITHOUT such a key changes,
// and the term is empty over a model with zero owned tags (`0010:C4`)."
// A-13 (Q2): "`--all` restores 0005's PREDICATE, not 0005's reader set".
// BOUNDARY — the two negative controls, without which the extension could
// be a mode-dependent hack.
func TestReq18And21And115_TheDemandSetIsModeIndependentAndEmptyWhereNoRowMatchOwns(t *testing.T) {
	t.Run("same-reader-set-under-all", func(t *testing.T) {
		model := writeFlowModel(t, flowMatchOnlyOwnedModel)
		binds := sideBindings(t, model, true, "mode=fast")

		def := readersOf(t, runNext(t, model, binds))
		all := readersOf(t, runNext(t, model, binds, "--all"))

		if !slices.Equal(def, all) {
			t.Errorf("readers under default = %v, under --all = %v; the "+
				"demand set is mode-INDEPENDENT. `--all` restores 0005's "+
				"PREDICATE, not 0005's reader set: on this class it invokes "+
				"one more reader than 0005 did, in BOTH modes", def, all)
		}
		if !slices.Contains(all, "side") {
			t.Errorf("readers under --all = %v; a mode-dependent demand set "+
				"would make the view mode-dependent and break C1's \"fixed "+
				"ONCE per invocation … and the same under --all\"", all)
		}
	})

	t.Run("no-row-match-owns-a-key", func(t *testing.T) {
		model := writeFlowModel(t, flowZeroOwnedMatchModel)
		bind := seedMatchArtifact(t, model, "status=draft")

		readers := readersOf(t, runNext(t, model, []string{bind}))
		if !slices.Equal(readers, []string{"state"}) {
			t.Errorf("readers = %v; want exactly [state]. The added term "+
				"finds nothing over a model where no row match-owns a key, "+
				"so the invoked set is unchanged — the property that keeps "+
				"the extension a strict widening and never a rewrite",
				readers)
		}
	})
}

// REQ-20 / `0011:F6`: the accepted `flow resolve` change is scoped to runs
// where "the requested outcome has a row on that key" — so the demand set
// is a union over THAT OUTCOME's rows, not over the model's.
// REQ-18: "The demand set is a property of the MODEL (and, for `resolve`,
// the requested outcome)" — the parenthesis is the whole assertion here.
// ADVERSARIAL — the NEGATIVE half of the demand-set scope, and the half
// suite-green cannot give. Every other oracle in this file asserts a reader
// IS invoked; a build that unioned the match-owned keys across every row
// BEFORE applying the outcome filter would satisfy all of them and still be
// wrong, because it would invoke `read.side` for outcome `stop` — whose
// only row (`stop-row`) matches on `status` alone. The cost is not
// cosmetic: with `side` unbound that build refuses `flow-artifact-missing`
// (exit 2) a run that must plan, which is `TestReq20And116`'s accepted cost
// LEAKING onto an outcome F6 never licensed it for.
func TestReq20AndF6_MatchOwnedKeysOfOtherOutcomesAreExcludedFromResolve(t *testing.T) {
	// `flowMatchOnlyOwnedModel` is the discriminating fixture: `mode-row`
	// MATCH-owns `mode` (served solely by `read.side`) under outcome `go`,
	// while outcome `stop`'s only row matches on `status` alone. So `go`
	// and `stop` MUST disagree on the reader set over one model — the
	// property a union-before-filter build erases.
	t.Run("stop-invokes-only-the-reader-its-own-row-needs", func(t *testing.T) {
		model := writeFlowModel(t, flowMatchOnlyOwnedModel)
		binds := sideBindings(t, model, true, "mode=fast")

		args := []string{"flow", "resolve", "--model", model}
		for _, b := range binds {
			args = append(args, "--artifact", b)
		}
		data := flowData(t, requireSuccess(t,
			append(args, "--outcome", "stop", "--as=json")...))

		if readers := readersOf(t, data); !slices.Equal(readers, []string{"state"}) {
			t.Errorf("`resolve --outcome stop` readers = %v; want exactly "+
				"[state]. `stop-row` is outcome `stop`'s only row and it "+
				"matches on `status` alone; `mode` is match-owned by "+
				"`mode-row`, which binds outcome `go`. A build that unioned "+
				"the match-owned term across every row BEFORE the outcome "+
				"filter reports [side state] here and passes every other "+
				"oracle in this file", readers)
		}
		// The same model at the OTHER outcome DOES pull the reader in, so
		// the assertion above is a real exclusion and not a fixture that
		// never demanded `side` in the first place. Seeded `mode=slow` so
		// `mode-row` does not match and `plain-row` is the sole `go`
		// candidate: at `mode=fast` BOTH `go` rows match and the kernel
		// refuses `ambiguous_match` before reporting a reader set, which is
		// the two-row fixture's deliberate shape (DEV-6). `read.side` is
		// demanded either way — the demand set is computed before the row
		// loop and does not depend on which rows end up matching.
		gmodel := writeFlowModel(t, flowMatchOnlyOwnedModel)
		gbinds := sideBindings(t, gmodel, true, "mode=slow")
		gargs := []string{"flow", "resolve", "--model", gmodel}
		for _, b := range gbinds {
			gargs = append(gargs, "--artifact", b)
		}
		gdata := flowData(t, requireSuccess(t,
			append(gargs, "--outcome", "go", "--as=json")...))
		if readers := readersOf(t, gdata); !slices.Contains(readers, "side") {
			t.Errorf("`resolve --outcome go` readers = %v; the same model's "+
				"`go` rows DO match-own `mode`. Without this arm the "+
				"exclusion above would pass against a build that never "+
				"invoked `read.side` for any outcome", readers)
		}
	})

	// The caller-visible harm, and the arm that fails LOUDEST on a
	// union-before-filter build: with `side` UNBOUND, `stop` must still
	// plan. `TestReq20And116` pins the mirror — that `go` under the same
	// unbound role refuses exit 2 — so the pair fixes the accepted cost's
	// exact boundary: it lands on `go` and NOT on `stop`.
	t.Run("stop-still-plans-with-the-other-outcomes-reader-unbound", func(t *testing.T) {
		model := writeFlowModel(t, flowMatchOnlyOwnedModel)
		binds := sideBindings(t, model, false)

		args := []string{"flow", "resolve", "--model", model}
		for _, b := range binds {
			args = append(args, "--artifact", b)
		}
		data := flowData(t, requireSuccess(t,
			append(args, "--outcome", "stop", "--as=json")...))

		if data["rule"] != "stop-row" {
			t.Errorf("the run yielded rule %#v; want `stop-row`. With "+
				"`side` unbound a union-before-filter build refuses "+
				"`flow-artifact-missing` (exit 2) for a role no row of "+
				"outcome `stop` needs — F6 scopes the accepted change to "+
				"runs where \"the requested outcome has a row on that "+
				"key\", and this run has none", data["rule"])
		}
	})
}

// REQ-17: "It therefore binds BOTH callers of `invokedReaders` —
// `flow_next.go` (`invokedReaders(req.model, \"\")`) and `flow_resolve.go`
// (`invokedReaders(req.model, outcome)`) — and MUST NOT be scoped to `next`"
// REQ-117 / `0011:S8`: "Both verbs' `readers` over the same model and
// outcome are EQUAL (one `invokedReaders`, one view), asserted on a fixture
// whose reader serves several keys."
// ADVERSARIAL — "a `next`-only term would have the two verbs assemble
// different views over one model, so `next` would report a row a candidate
// that `resolve` then refuses `flow-no-match` for the key it never read."
func TestReq17And117_BothVerbsAgreeOnTheReaderSetOverOneModelAndOutcome(t *testing.T) {
	// The SOLO fixture, for DEV-4's reason: over the two-row model
	// `mode=fast` makes BOTH `go` rows match and `flow resolve` refuses
	// `ambiguous_match` before it can report a `readers` set to compare.
	// The solo model satisfies S8's stated fixture requirement — its
	// `read.side` serves `mode` AND `extra`, so the agreement is asserted
	// over a reader carrying more than the one key the term found — and
	// the two-row shape is the BREAKING arm's subject, which has its own
	// oracle below. Recorded as DEV-6.
	model := writeFlowModel(t, flowMatchOnlyOwnedSoloModel)
	binds := sideBindings(t, model, true, "mode=fast")

	nextReaders := readersOf(t, runNext(t, model, binds))

	args := []string{"flow", "resolve", "--model", model}
	for _, b := range binds {
		args = append(args, "--artifact", b)
	}
	resolveData := flowData(t, requireSuccess(t,
		append(args, "--outcome", "go", "--as=json")...))
	resolveReaders := readersOf(t, resolveData)

	if !slices.Equal(nextReaders, resolveReaders) {
		t.Errorf("`next` readers = %v, `resolve --outcome go` readers = %v; "+
			"ONE `invokedReaders`, ONE view. A `next`-only term makes the "+
			"two verbs assemble different views over one model, so `next` "+
			"reports a row a candidate that `resolve` then refuses "+
			"`flow-no-match` for the key it never read",
			nextReaders, resolveReaders)
	}
	// `read.side` serves TWO keys, which is the fixture shape S8 names —
	// so the agreement is asserted over a reader that carries more than
	// the one key the term found.
	if !slices.Contains(resolveReaders, "side") {
		t.Errorf("`resolve` readers = %v; the term binds BOTH callers and "+
			"MUST NOT be scoped to `next`", resolveReaders)
	}
}

// REQ-19: "after this term it invokes the reader, so a run that used to
// refuse `flow-no-match` MAY now yield a plan, `flow-artifact-missing`
// (exit 2, the reader's role unbound), or the reader's own refusal
// (exit 3)."
// REQ-93 / MVV 8: "Run `flow resolve --outcome <o>` over the same model with
// the reader bound, unbound, and refusing: assert a plan,
// `flow-artifact-missing` (exit 2), and the reader's refusal (exit 3)
// respectively, and `flow-no-match` for that key in none".
// DOMAIN EDGE — the three arms, "each of those names the remedy the old
// refusal hid." `flow-no-match` in NONE is the half that discriminates: it
// is the un-extended set's answer in all three.
func TestReq19And93_ResolveSeparatesIntoThreeArmsAndRefusesFlowNoMatchInNone(t *testing.T) {
	// The SOLO fixture is what separates the three arms: with the two-row
	// fixture, `mode=fast` makes both `go` rows match and the kernel
	// refuses `ambiguous_match` rather than plan. The two-row shape is the
	// BREAKING arm's subject and has its own oracle below.
	t.Run("bound-and-answering-yields-a-plan", func(t *testing.T) {
		model := writeFlowModel(t, flowMatchOnlyOwnedSoloModel)
		binds := sideBindings(t, model, true, "mode=fast")

		args := []string{"flow", "resolve", "--model", model}
		for _, b := range binds {
			args = append(args, "--artifact", b)
		}
		data := flowData(t, requireSuccess(t,
			append(args, "--outcome", "go", "--as=json")...))

		// PRE-extension this run refuses `flow-no-match`: `read.side` is
		// never invoked, `mode` is absent, and the kernel folds an absent
		// match key into non-match — over an artifact that HOLDS the fact.
		// That is DEV-8's ADV-2 shape with the other refusal, and the term
		// is what turns it into a plan.
		if data["rule"] != "mode-row" {
			t.Errorf("the run yielded rule %#v; want a plan for `mode-row`. "+
				"Pre-extension this refuses `flow-no-match` for want of a "+
				"reader it never invoked, over an artifact that holds the "+
				"fact — the refusal the term removes", data["rule"])
		}
	})

	t.Run("role-unbound-yields-artifact-missing-not-no-match", func(t *testing.T) {
		model := writeFlowModel(t, flowMatchOnlyOwnedSoloModel)
		binds := sideBindings(t, model, false)

		args := []string{"flow", "resolve", "--model", model}
		for _, b := range binds {
			args = append(args, "--artifact", b)
		}
		ce := requireRefusal(t, "flow-artifact-missing", 2,
			append(args, "--outcome", "go", "--as=json")...)
		if ce.Code == "flow-no-match" {
			t.Error("the refusal is `flow-no-match`, the un-extended set's " +
				"answer; the term invokes the reader BEFORE the kernel and " +
				"above the escape phase, so an unbound role is what the " +
				"caller is told about")
		}
	})

	t.Run("reader-refusing-yields-exit-3-not-no-match", func(t *testing.T) {
		model := writeFlowModel(t, flowRefusingSideReaderModel)
		state := newFlowArtifact(t, "state.artifact")
		stateBind := artifactBinding(flowMatchRole, state)
		requireSuccess(t, "flow", "set-state", "--model", model,
			"--artifact", stateBind, "--write", "status=draft", "--as=json")
		sideBind := artifactBinding(flowSideRole,
			newFlowArtifact(t, "side.artifact"))

		_, _, err := runCmd(t, "flow", "resolve", "--model", model,
			"--artifact", stateBind, "--artifact", sideBind,
			"--outcome", "go", "--as=json")
		if err == nil {
			t.Fatal("the run SUCCEEDED with the reader's locator " +
				"unreachable; the reader's own refusal is an exit-3 class")
		}
		// A15's limit is respected: the oracle asserts the EXIT CODE and
		// the not-a-`flow-no-match` property, "not a single code, since
		// another refusal mechanism takes a different exit-3 code".
		if exit := clierr.ExitCodeFor(err); exit != 3 {
			t.Errorf("exit = %d; want 3 — the reader's own refusal is an "+
				"ENVIRONMENT class, not a request error. code = %q",
				exit, clierr.ErrorCode(err))
		}
		if code := clierr.ErrorCode(err); code == "flow-no-match" {
			t.Error("the refusal is `flow-no-match`; that is the " +
				"un-extended set's answer in all three arms, and it is what " +
				"the term replaces with a refusal that names the remedy")
		}
	})
}

// REQ-20: "The demand set is a union over the outcome's rows, so the reader
// is invoked even when the row `resolve` would select does not itself match
// on the key: over an UNBOUND or REFUSING reader such a run turns from a
// plan into exit 2/3, before the kernel and above the escape phase — a
// class this contract NAMES AND ACCEPTS here"
// REQ-116 / `0011:S8`: "The BREAKING arm: add a second ordinary row for the
// same outcome that does not match on the key; with the reader unbound,
// assert the pre-extension build returns that row's plan and the
// post-extension build `flow-artifact-missing` — pinned so the DEV-8 class
// is a known cost, not a surprise."
// ADVERSARIAL — the accepted behaviour-change class. This oracle asserts a
// build gets WORSE for one shape, on purpose, so a future reader finds the
// cost recorded rather than discovering it in the field.
func TestReq20And116_TheBreakingArmIsPinnedAsAKnownCost(t *testing.T) {
	model := writeFlowModel(t, flowMatchOnlyOwnedModel)
	// `plain-row` is a second ordinary row for outcome `go` that does NOT
	// match on `mode`. Pre-extension the run returns its plan, because
	// `read.side` is never invoked and the unbound role never bites.
	binds := sideBindings(t, model, false)

	args := []string{"flow", "resolve", "--model", model}
	for _, b := range binds {
		args = append(args, "--artifact", b)
	}
	args = append(args, "--outcome", "go", "--as=json")

	stdout, _, err := runCmd(t, args...)
	if err == nil {
		t.Fatalf("`flow resolve --outcome go` SUCCEEDED with `side` "+
			"unbound, returning:\n%s\nThe demand set is a UNION over the "+
			"outcome's rows, so `read.side` is invoked even though the row "+
			"`resolve` selects (`plain-row`) does not itself match on "+
			"`mode`. Over an unbound reader that run turns from a plan into "+
			"exit 2 — the class this contract NAMES AND ACCEPTS.\n"+
			"One asymmetry is stated rather than hidden: "+
			"`flow-artifact-missing` is raised BEFORE the kernel and is not "+
			"one of the five RefusalKinds, so NO escape row can rescue it — "+
			"where the `flow-no-match` it replaces IS escapable. The trade "+
			"is a rescuable-but-mute refusal for an unrescuable one that "+
			"names the role and the accessor it needs", stdout)
	}

	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("refusal is not a structured CLIError: %v", err)
	}
	if ce.Code != "flow-artifact-missing" {
		t.Errorf("code = %q; want %q — the reader is pulled in ABOVE the "+
			"kernel, so the caller is told which ROLE to bind, and binding "+
			"it restores the plan", ce.Code, "flow-artifact-missing")
	}
	if ce.Param != flowSideRole {
		t.Errorf("param = %q; want the unbound ROLE %q — the refusal names "+
			"the remedy", ce.Param, flowSideRole)
	}

	// Binding the role RESTORES the plan, which is what makes the class a
	// known cost rather than a dead end.
	restored := sideBindings(t, model, true)
	rargs := []string{"flow", "resolve", "--model", model}
	for _, b := range restored {
		rargs = append(rargs, "--artifact", b)
	}
	if _, _, rerr := runCmd(t, append(rargs, "--outcome", "go", "--as=json")...); rerr != nil {
		t.Errorf("binding the role did not restore a resolution: %v\nThe "+
			"model DECLARES that reader (`checkAccessorBindings`), and "+
			"binding it is the remedy the refusal names", rerr)
	}
}

// REQ-35: "The `flow resolve` VERB changes in exactly the one class the
// demand-set paragraph names, and nowhere else."
// REQ-118 / `0011:S8`: "`flow resolve` over `models/rdr.toml` and every
// shipped fixture is byte-identical before and after (A15 iii)."
// A-11: discharged against a pre-change build of the same tree — the oracle
// may be a golden payload "provided it is captured from the pre-change
// behaviour and asserts the FULL payload rather than a subset — which is
// the gap S8 names".
// ADVERSARIAL — the payload-diff form "the spike could not reach": the
// spike compared only `models/rdr.toml` byte-for-byte and left the shipped
// fixtures to suite-green, "a fixture asserting only a subset of its
// payload being exactly what suite-green cannot catch."
func TestReq35And118_ResolveOverTheCheckedInModelAndShippedFixturesIsUnchanged(t *testing.T) {
	// Every owned match key in `models/rdr.toml` is also WRITTEN, so it is
	// already in `RequiresOwned` and the added term finds nothing: no rule
	// there is in the changed class (A15 iii). The golden is captured from
	// the PRE-change behaviour and asserts the FULL payload.
	model := repoModelPath(t)
	art := newFlowArtifact(t, "rdr.artifact")
	bind := artifactBinding("rdr", art)
	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", bind, "--write", "stage=resolved",
		"--write", "status=draft", "--write", "gate_passed=false",
		"--as=json")

	for _, outcome := range []string{"advance", "revise", "abandon"} {
		t.Run(outcome, func(t *testing.T) {
			stdout, _, err := runCmd(t, "flow", "resolve", "--model", model,
				"--artifact", bind, "--outcome", outcome, "--as=json")
			if err != nil {
				t.Fatalf("`flow resolve --outcome %s` refused: %v\nNo rule "+
					"of the checked-in model is in the changed class",
					outcome, err)
			}
			// The FULL payload, not a subset: a golden captured from the
			// pre-change behaviour. The `model` path varies per run, so
			// it is the one field normalized out.
			got := normalizeModelPath(stdout, model)
			want := resolveGolden[outcome]
			if got != want {
				t.Errorf("`flow resolve --outcome %s` payload changed.\n"+
					"got:  %s\nwant: %s\n`flow resolve` over "+
					"`models/rdr.toml` is byte-identical before and after; "+
					"the changed class is a match-ONLY owned key, and every "+
					"owned match key here is also written",
					outcome, got, want)
			}
		})
	}

	// The shipped fixtures, in the same full-payload form — the half the
	// spike attested only by suite-green, and the half this subtest used
	// to discharge with ONE fixture (`flowMVVModel`) against S8's "every
	// shipped fixture".
	for _, tc := range shippedResolveGoldens {
		t.Run("shipped-fixture/"+tc.name+"/"+tc.outcome, func(t *testing.T) {
			fx := writeFlowModel(t, tc.model)

			var binds []string
			for _, s := range tc.seeds {
				art := newFlowArtifact(t, s.role+".artifact")
				bind := artifactBinding(s.role, art)
				args := []string{"flow", "set-state", "--model", fx,
					"--artifact", bind}
				for _, w := range s.writes {
					args = append(args, "--write", w)
				}
				requireSuccess(t, append(args, "--as=json")...)
				binds = append(binds, bind)
			}
			for _, role := range tc.bindOnly {
				binds = append(binds,
					artifactBinding(role, newFlowArtifact(t, role+".artifact")))
			}

			args := []string{"flow", "resolve", "--model", fx}
			for _, b := range binds {
				args = append(args, "--artifact", b)
			}
			for _, tg := range tc.tags {
				args = append(args, "--tag", tg)
			}
			stdout, _, err := runCmd(t,
				append(args, "--outcome", tc.outcome, "--as=json")...)
			if err != nil {
				t.Fatalf("`flow resolve --outcome %s` over the shipped "+
					"fixture `%s` refused: %v\nThis entry names a run that "+
					"reaches a PLAN; a refusal here means the fixture or "+
					"its seeding drifted", tc.outcome, tc.name, err)
			}
			got := normalizeModelPath(stdout, fx)
			if got != tc.golden {
				t.Errorf("`flow resolve --outcome %s` over `%s` changed.\n"+
					"got:  %s\nwant: %s\nNo shipped flow-reachable fixture "+
					"is in the changed class (A15 iii), so every one of "+
					"these payloads is byte-identical before and after",
					tc.outcome, tc.name, got, tc.golden)
			}
		})
	}

	// REQ-118 says "every shipped fixture", so the sweep's MEMBERSHIP is
	// itself an oracle: a fixture added later must land in the table or in
	// the named excluded set, never escape both silently. Without this the
	// table is a snapshot of today's corpus and S8 decays with every new
	// fixture.
	t.Run("every-shipped-fixture-is-swept-or-named-excluded", func(t *testing.T) {
		swept := map[string]bool{}
		for _, tc := range shippedResolveGoldens {
			swept[tc.name] = true
		}
		corpus := declaredFlowFixtures(t)
		declared := map[string]bool{}
		for _, name := range corpus {
			declared[name] = true
		}
		for _, name := range corpus {
			if swept[name] {
				continue
			}
			if _, ok := unsweptShippedFixtures[name]; !ok {
				t.Errorf("`%s` is a shipped fixture that is neither in "+
					"`shippedResolveGoldens` nor named in "+
					"`unsweptShippedFixtures`.\nS8/REQ-118 states the "+
					"byte-identity guarantee over EVERY shipped fixture; a "+
					"fixture that escapes both sets discharges it by "+
					"omission. Add a golden entry, or record why the "+
					"fixture has no reachable `flow resolve` plan", name)
			}
		}
		// The mirror: a name in the excluded set that is no longer a
		// fixture, or that IS swept, means the two lists drifted apart.
		for name := range unsweptShippedFixtures {
			if !declared[name] {
				t.Errorf("`%s` is named in `unsweptShippedFixtures` but is "+
					"not declared by any test source; the exclusion "+
					"outlived its fixture", name)
			}
			if swept[name] {
				t.Errorf("`%s` is BOTH swept and named excluded; the "+
					"exclusion is stale and its reason is now false", name)
			}
		}
	})
}

// normalizeModelPath replaces the run-specific model path with a stable
// token so a full-payload golden compares.
func normalizeModelPath(stdout, path string) string {
	return strings.TrimSpace(strings.ReplaceAll(stdout, path, "<model>"))
}

// resolveGolden and resolveFixtureGolden are the FULL `flow resolve`
// payloads captured from the PRE-change build of this tree, per A-11.
//
// S8 requires byte-identity "before and after (A15 iii)", and names the gap
// the spike could not close: it compared `models/rdr.toml` binary-to-binary
// with `cmp` but attested the shipped fixtures only by suite-green, "a
// fixture asserting only a subset of its payload being exactly what
// suite-green cannot catch." These goldens close it — they assert the WHOLE
// payload, so a field this contract was never meant to touch cannot change
// unnoticed.
//
// The run-specific model path is the one value normalized out; everything
// else is the wire form verbatim.
var resolveGolden = map[string]string{
	"advance": `{"type":"ok","data":{"model":"<model>","revision":"",` +
		`"observed":{},"owned":{"gate_passed":"false","stage":"resolved",` +
		`"status":"draft"},"readers":["rdr-status"],"outcome":"advance",` +
		`"rule":"prelock","gates":[],"emit":{},"next":{"stage":"prelocked"},` +
		`"writes":{"stage":"prelocked"},"clear":[],"escaped":false}}`,
	"revise": `{"type":"ok","data":{"model":"<model>","revision":"",` +
		`"observed":{},"owned":{"gate_passed":"false","stage":"resolved",` +
		`"status":"draft"},"readers":["rdr-status"],"outcome":"revise",` +
		`"rule":"resolve-route-back","gates":[],"emit":{},"next":{"stage":"refined"},` +
		`"writes":{"stage":"refined"},"clear":[],"escaped":false}}`,
	"abandon": `{"type":"ok","data":{"model":"<model>","revision":"",` +
		`"observed":{},"owned":{"gate_passed":"false","stage":"resolved",` +
		`"status":"draft"},"readers":["rdr-status"],"outcome":"abandon",` +
		`"rule":"resolve-abandon","gates":[],"emit":{},` +
		`"next":{"stage":"dropped","status":"abandoned"},` +
		`"writes":{"stage":"dropped","status":"abandoned"},"clear":[],` +
		`"escaped":false}}`,
}

// --- the shipped-fixture sweep (S8 / REQ-118) ----------------------------

// fixtureSeed establishes one artifact role through the PRODUCTION write
// path before the resolve under test — never by writing a file whose format
// the test claims to understand (the 0005 harness rule).
type fixtureSeed struct {
	role   string
	writes []string
}

// shippedFixtureGolden is one `flow resolve` run over one shipped fixture,
// with the FULL payload it must keep producing.
type shippedFixtureGolden struct {
	name    string
	model   string
	outcome string
	seeds   []fixtureSeed
	// bindOnly names roles the run must BIND but not seed — a reader the
	// model declares whose absence would refuse `flow-artifact-missing`
	// before the payload exists.
	bindOnly []string
	// tags are `--tag k=v` supplies for OBSERVED keys a row matches on.
	// A fixture whose only match atoms are observed refuses
	// `flow-no-match` unsupplied and reaches a plan supplied, so the
	// supply is what makes the run a byte-identity subject at all.
	tags   []string
	golden string
}

// shippedResolveGoldens is S8/REQ-118's sweep: every shipped fixture with a
// reachable `flow resolve` PLAN, each asserted as a FULL payload captured
// from this tree's behaviour per A-11 ("the oracle may be a golden payload
// … provided it … asserts the FULL payload rather than a subset — which is
// the gap S8 names").
//
// This replaces a one-fixture discharge. S8 states the guarantee over
// "every shipped fixture" and A15's limit (b) records why suite-green
// cannot give it: "a fixture asserting a subset of its payload would not
// catch an unasserted-field change". One entry per (fixture, outcome) that
// reaches exit 0 — a refusal has no payload to be byte-identical about, and
// the refusal fixtures are pinned by code and exit elsewhere in this suite.
//
// The goldens are deliberately brittle: any change to the resolve envelope
// re-breaks every entry at once. That is the change-detection S8 asks for,
// and the cost is one mechanical regeneration confined to this variable.
var shippedResolveGoldens = []shippedFixtureGolden{
	{
		name: "flowMVVModel", model: flowMVVModel, outcome: "hold",
		seeds:    []fixtureSeed{{flowStateRole, []string{"status=draft"}}},
		bindOnly: []string{flowOrphanRole},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{},"owned":{"status":"draft"},"readers":["state"],` +
			`"outcome":"hold","rule":"hold-draft","gates":[],"emit":{},` +
			`"next":{"status":"draft"},"writes":{"status":"draft"},` +
			`"clear":[],"escaped":false}}`,
	},
	{
		name: "flowEscapeModel", model: flowEscapeModel, outcome: "hold",
		seeds:    []fixtureSeed{{flowStateRole, []string{"status=draft"}}},
		bindOnly: []string{flowOrphanRole},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{},"owned":{"status":"draft"},"readers":["state"],` +
			`"outcome":"hold","rule":"hold-draft","gates":[],"emit":{},` +
			`"next":{"status":"draft"},"writes":{"status":"draft"},` +
			`"clear":[],"escaped":false}}`,
	},
	{
		// The ESCAPED payload shape — `escaped:true` plus `escape_class`,
		// two fields no other entry in this table carries. Without it the
		// sweep would leave the escape envelope to suite-green.
		name: "flowEscapeModel", model: flowEscapeModel, outcome: "bail",
		seeds:    []fixtureSeed{{flowStateRole, []string{"status=draft"}}},
		bindOnly: []string{flowOrphanRole},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{},"owned":{},"readers":[],"outcome":"bail",` +
			`"rule":"bail-escape","gates":[],"emit":{},"next":{},` +
			`"writes":{},"clear":[],"escaped":true,` +
			`"escape_class":"no_match"}}`,
	},
	{
		name: "flowEscapeOtherOutcomeModel", model: flowEscapeOtherOutcomeModel,
		outcome:  "hold",
		seeds:    []fixtureSeed{{flowStateRole, []string{"status=draft"}}},
		bindOnly: []string{flowOrphanRole, "sidecar"},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{},"owned":{"status":"draft"},"readers":["state"],` +
			`"outcome":"hold","rule":"hold-draft","gates":[],"emit":{},` +
			`"next":{"status":"draft"},"writes":{"status":"draft"},` +
			`"clear":[],"escaped":false}}`,
	},
	{
		name: "flowDomainModel", model: flowDomainModel, outcome: "advance",
		seeds: []fixtureSeed{{flowStateRole, []string{"status=draft"}}},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{},"owned":{"status":"draft"},"readers":["state"],` +
			`"outcome":"advance","rule":"advance","gates":[],"emit":{},` +
			`"next":{"status":"final"},"writes":{"status":"final"},` +
			`"clear":[],"escaped":false}}`,
	},
	{
		// The GATED payload shape — a non-empty `gates` array, which no
		// other entry carries.
		name: "flowGatedNextModel", model: flowGatedNextModel, outcome: "advance",
		seeds: []fixtureSeed{{flowStateRole, []string{"status=draft", "flag=true"}}},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{},"owned":{"flag":"true","status":"draft"},` +
			`"readers":["state"],"outcome":"advance",` +
			`"rule":"gated-reported","gates":[{"id":"reported",` +
			`"result":"allow"}],"emit":{},"next":{"status":"final"},` +
			`"writes":{"status":"final"},"clear":[],"escaped":false}}`,
	},
	{
		name: "flowSetObservedModel", model: flowSetObservedModel, outcome: "advance",
		seeds: []fixtureSeed{{flowStateRole, []string{"status=draft"}}},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{},"owned":{"status":"draft"},"readers":["state"],` +
			`"outcome":"advance","rule":"advance","gates":[],"emit":{},` +
			`"next":{"status":"final"},"writes":{"status":"final"},` +
			`"clear":[],"escaped":false}}`,
	},
	{
		name: "flowMatchClassesModel", model: flowMatchClassesModel, outcome: "gamma",
		seeds: []fixtureSeed{{flowStateRole, []string{"status=draft"}}},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{},"owned":{"status":"draft"},"readers":["state"],` +
			`"outcome":"gamma","rule":"no-match-atoms","gates":[],"emit":{},` +
			`"next":{"status":"final"},"writes":{"status":"final"},` +
			`"clear":[],"escaped":false}}`,
	},
	{
		name: "flowSetKindedMatchModel", model: flowSetKindedMatchModel, outcome: "hold",
		seeds: []fixtureSeed{{flowStateRole, []string{"status=draft"}}},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{},"owned":{"status":"draft"},"readers":["state"],` +
			`"outcome":"hold","rule":"scalar-control","gates":[],"emit":{},` +
			`"next":{"status":"final"},"writes":{"status":"final"},` +
			`"clear":[],"escaped":false}}`,
	},
	{
		name: "flowAbsentMatchKeyModel", model: flowAbsentMatchKeyModel, outcome: "hold",
		seeds: []fixtureSeed{{flowStateRole, []string{"status=draft"}}},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{},"owned":{"status":"draft"},"readers":["state"],` +
			`"outcome":"hold","rule":"recognized-only","gates":[],"emit":{},` +
			`"next":{"status":"final"},"writes":{"status":"final"},` +
			`"clear":[],"escaped":false}}`,
	},
	{
		name: "flowFullyResolvedModel", model: flowFullyResolvedModel, outcome: "advance",
		seeds: []fixtureSeed{{flowStateRole, []string{"status=draft", "flag=true"}}},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{},"owned":{"flag":"true","status":"draft"},` +
			`"readers":["state"],"outcome":"advance","rule":"resolved-row",` +
			`"gates":[],"emit":{},"next":{"status":"final"},` +
			`"writes":{"status":"final"},"clear":[],"escaped":false}}`,
	},
	{
		// A set-valued write whose members carry `<` and `&` — the only
		// entry whose golden pins the set wire form through the resolve
		// envelope (`0005:MVV`, REQ-70, REQ-108). Reachable once BOTH
		// owned keys are seeded; seeding only `status` refuses
		// `flow-owned-state-unavailable`, which is why it once sat in the
		// excluded set.
		name: "flowSetPlanModel", model: flowSetPlanModel, outcome: "advance",
		seeds: []fixtureSeed{{flowStateRole, []string{
			"status=draft", `labels=["plain"]`}}},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{},"owned":{"labels":"[\"plain\"]","status":"draft"},` +
			`"readers":["state"],"outcome":"advance","rule":"advance",` +
			`"gates":[],"emit":{},` +
			`"next":{"labels":"[\"a<b\",\"x&y\"]","status":"final"},` +
			`"writes":{"labels":"[\"a<b\",\"x&y\"]","status":"final"},` +
			`"clear":[],"escaped":false}}`,
	},
	{
		// The `side`-role writer. It is NOT in the changed class: its
		// `side-advance` row matches on `mode` AND writes `mode`, so
		// `mode` is already in `RequiresOwned` and the added term finds
		// nothing (A15 iii, and the RDR's ground-sweep finding that no
		// flow-reachable fixture carries a match-ONLY owned key). Its
		// role adjacency to the changed class is not a reason to skip the
		// byte-identity assertion S8 owes it.
		name: "flowSideWriterModel", model: flowSideWriterModel, outcome: "advance",
		seeds: []fixtureSeed{{flowSideRole, []string{"mode=fast"}}},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{},"owned":{"mode":"fast"},"readers":["side"],` +
			`"outcome":"advance","rule":"side-advance","gates":[],"emit":{},` +
			`"next":{"mode":"slow"},"writes":{"mode":"slow"},` +
			`"clear":[],"escaped":false}}`,
	},
	{
		// The guard fixture's SATISFIED arm. It refuses
		// `flow-guard-unevaluable` only when `opt` is ABSENT; seeded
		// `opt=p` the guard evaluates and the row plans, so the fixture
		// does have a byte-identity subject. Its unevaluable arm stays
		// pinned by the oracle that owns it.
		name: "flowGuardUnevaluableModel", model: flowGuardUnevaluableModel,
		outcome: "advance",
		seeds: []fixtureSeed{{flowStateRole, []string{
			"status=draft", "opt=p"}}},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{},"owned":{"opt":"p","status":"draft"},` +
			`"readers":["state"],"outcome":"advance","rule":"guarded-advance",` +
			`"gates":[],"emit":{},"next":{"status":"final"},` +
			`"writes":{"status":"final"},"clear":[],"escaped":false}}`,
	},
	{
		// The zero-owned-match control's SUPPLIED arm. `hint-row` matches
		// an OBSERVED key, so it refuses `flow-no-match` unsupplied and
		// plans under `--tag hint=x`. This is also the sweep's only entry
		// with a non-empty `observed` map, so the golden pins that field's
		// wire form.
		name: "flowZeroOwnedMatchModel", model: flowZeroOwnedMatchModel,
		outcome: "advance",
		seeds:   []fixtureSeed{{flowStateRole, []string{"status=draft"}}},
		tags:    []string{"hint=x"},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{"hint":"x"},"owned":{"status":"draft"},` +
			`"readers":["state"],"outcome":"advance","rule":"hint-row",` +
			`"gates":[],"emit":{},"next":{"status":"final"},` +
			`"writes":{"status":"final"},"clear":[],"escaped":false}}`,
	},
	{
		// Three observed match atoms, all supplied — the row plans. The
		// `flow-no-match` its exclusion cited is the UNSUPPLIED arm only.
		name: "flowMixedStateRowModel", model: flowMixedStateRowModel,
		outcome: "advance",
		seeds:   []fixtureSeed{{flowStateRole, []string{"status=draft"}}},
		tags:    []string{"holds=yes", "fails=yes", "missing=yes"},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{"fails":"yes","holds":"yes","missing":"yes"},` +
			`"owned":{"status":"draft"},"readers":["state"],` +
			`"outcome":"advance","rule":"mixed-row","gates":[],"emit":{},` +
			`"next":{"status":"final"},"writes":{"status":"final"},` +
			`"clear":[],"escaped":false}}`,
	},
	{
		// The sort-order fixture's supplied arm, and the sweep's ONLY
		// entry whose `gates` array is non-empty — without it the gate
		// result's wire form (`{"id","result"}`) is left to suite-green,
		// exactly the gap A15 limit (b) names.
		name: "flowSortOrderModel", model: flowSortOrderModel,
		outcome: "advance",
		seeds: []fixtureSeed{{flowStateRole, []string{
			"status=draft", "mid=m"}}},
		tags: []string{"zebra=z", "alpha=a"},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{"alpha":"a","zebra":"z"},` +
			`"owned":{"mid":"m","status":"draft"},"readers":["state"],` +
			`"outcome":"advance","rule":"sort-row",` +
			`"gates":[{"id":"beta","result":"allow"}],"emit":{},` +
			`"next":{"mid":"m","status":"final"},` +
			`"writes":{"mid":"m","status":"final"},"clear":[],` +
			`"escaped":false}}`,
	},
	{
		// The PAIR fixture's supplied arm. Its tag key `beta` and its gate
		// id `beta` are the SAME token in two namespaces, and this payload
		// is where that separation is observable on the wire: `observed`
		// carries `beta` the tag while `gates` carries `beta` the accessor
		// id, in one envelope. A build that conflated the namespaces could
		// not produce both.
		name: "flowUnknownPairModel", model: flowUnknownPairModel,
		outcome: "advance",
		seeds:   []fixtureSeed{{flowStateRole, []string{"status=draft"}}},
		tags:    []string{"beta=b", "alpha=a"},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{"alpha":"a","beta":"b"},` +
			`"owned":{"status":"draft"},"readers":["state"],` +
			`"outcome":"advance","rule":"pair-row",` +
			`"gates":[{"id":"beta","result":"allow"}],"emit":{},` +
			`"next":{"status":"final"},"writes":{"status":"final"},` +
			`"clear":[],"escaped":false}}`,
	},
	{
		// The uncomparable-guard fixture's COMPARABLE arm. `size` is an
		// `int` tag and the guard's bound is the string `"3"`; the
		// refusal its exclusion cited is that uncomparability, which the
		// oracle owning it pins. Seeded, the row still plans, so the
		// fixture has a payload S8 owes an assertion.
		name: "flowUncomparableGuardModel", model: flowUncomparableGuardModel,
		outcome: "advance",
		seeds: []fixtureSeed{{flowStateRole, []string{
			"status=draft", "size=9", "gone=ready"}}},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{},"owned":{"gone":"ready","size":"9",` +
			`"status":"draft"},"readers":["state"],"outcome":"advance",` +
			`"rule":"uncomparable-row","gates":[],"emit":{},` +
			`"next":{"status":"final"},"writes":{"status":"final"},` +
			`"clear":[],"escaped":false}}`,
	},
	{
		// The same fixture's second outcome — `absent-guard-row`, whose
		// guard is on `gone`. Seeded it evaluates and the row plans.
		name: "flowUncomparableGuardModel", model: flowUncomparableGuardModel,
		outcome: "hold",
		seeds: []fixtureSeed{{flowStateRole, []string{
			"status=draft", "size=9", "gone=ready"}}},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{},"owned":{"gone":"ready","size":"9",` +
			`"status":"draft"},"readers":["state"],"outcome":"hold",` +
			`"rule":"absent-guard-row","gates":[],"emit":{},` +
			`"next":{"status":"final"},"writes":{"status":"final"},` +
			`"clear":[],"escaped":false}}`,
	},
	{
		// The owned-unavailable fixture's AVAILABLE arm. Its reader
		// serves all three owned keys, so seeding them makes the guard
		// evaluable and the row plans; `flow-owned-state-unavailable` is
		// the unseeded arm only.
		name: "flowOwnedUnavailableModel", model: flowOwnedUnavailableModel,
		outcome: "advance",
		seeds: []fixtureSeed{{flowStateRole, []string{
			"status=draft", "size=9", "missingowned=q"}}},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{},"owned":{"missingowned":"q","size":"9",` +
			`"status":"draft"},"readers":["state"],"outcome":"advance",` +
			`"rule":"unavail-row","gates":[],"emit":{},` +
			`"next":{"missingowned":"x","status":"final"},` +
			`"writes":{"missingowned":"x","status":"final"},` +
			`"clear":[],"escaped":false}}`,
	},
	{
		name: "flowLooseWriterModel", model: flowLooseWriterModel, outcome: "advance",
		seeds: []fixtureSeed{{flowStateRole, []string{"status=draft"}}},
		golden: `{"type":"ok","data":{"model":"<model>","revision":"",` +
			`"observed":{},"owned":{"status":"draft"},"readers":["state"],` +
			`"outcome":"advance","rule":"seed-row","gates":[],"emit":{},` +
			`"next":{"status":"final"},"writes":{"status":"final"},` +
			`"clear":[],"escaped":false}}`,
	},
}

// declaredFlowFixtures returns the name of every package-level
// `flow*Model` identifier this package's test sources DECLARE, read from
// the sources themselves with `go/parser`.
//
// The membership oracle below needs a source of truth that a new fixture
// cannot silently escape. A hand-listed map is not one: a fixture omitted
// from it escapes the sweep, the exclusions, AND the check that is
// supposed to catch that — the oracle passes vacuously and S8 decays with
// every fixture added. Nor is a shared registry variable, since a fixture
// can simply not join it.
//
// The declarations are the one artifact a fixture cannot be used without
// creating, so they are what this reads. A fixture added in any form —
// `const`, `var`, computed by `strings.ReplaceAll` — is found by its
// declaration, and the oracle fails until it is swept or named excluded.
func declaredFlowFixtures(t *testing.T) []string {
	t.Helper()
	paths, err := filepath.Glob("*_test.go")
	if err != nil {
		t.Fatalf("globbing this package's test sources: %v", err)
	}
	if len(paths) == 0 {
		t.Fatal("no `*_test.go` sources found; the membership oracle " +
			"cannot derive the fixture corpus and would pass vacuously")
	}
	// This file's own tables NAME fixtures without declaring them; only
	// declarations count.
	want := regexp.MustCompile(`^flow[A-Za-z0-9]*Model$`)
	fset := token.NewFileSet()
	var names []string
	for _, path := range paths {
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("reading %s: %v", path, err)
		}
		f, err := parser.ParseFile(fset, path, src, 0)
		if err != nil {
			t.Fatalf("parsing %s: %v", path, err)
		}
		for _, decl := range f.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || (gen.Tok != token.CONST && gen.Tok != token.VAR) {
				continue
			}
			for _, spec := range gen.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, id := range vs.Names {
					if want.MatchString(id.Name) {
						names = append(names, id.Name)
					}
				}
			}
		}
	}
	slices.Sort(names)
	return slices.Compact(names)
}

// unsweptShippedFixtures names every shipped fixture the golden sweep does
// NOT cover, with the reason. S8's guarantee is over "every shipped
// fixture"; discharging it honestly means the uncovered ones are NAMED, not
// silently absent — so this map is the second half of the oracle above.
//
// Two classes, and neither is a byte-identity subject:
//
//   - NO REACHABLE PLAN — the fixture refuses on EVERY outcome under
//     every seeding, so it has no success payload to be byte-identical
//     about. Its refusal code and exit are pinned by the oracles that own
//     it. The bar is that no seeding reaches exit 0, not that the
//     fixture's headline oracle is a refusal: a fixture whose refusal is
//     an UNSEEDED arm (an unevaluable guard, an unsupplied observed match
//     key, an unavailable owned key) does plan when seeded, and belongs in
//     the sweep above.
//   - INSIDE THE CHANGED CLASS — 0011's own match-only-owned-key fixtures.
//     A15 iii's warrant is that no shipped flow-reachable fixture is in the
//     changed class; these are the deliberate exceptions the RDR BUILT, so
//     asserting they are unchanged would assert the term did nothing.
var unsweptShippedFixtures = map[string]string{
	// --- no reachable plan: refusal fixtures ---
	"flowInvalidModel":          "refuses `flow-model-invalid` at load",
	"flowTwoReadersOneKeyModel": "refuses `flow-model-invalid` at load",
	"flowAmbiguousModel":        "refuses `flow-ambiguous-match`",
	"flowGateDenyModel":         "refuses `flow-gate-denied` / `flow-gate-indeterminate` on every outcome",
	"flowGateFailModel":         "refuses `flow-accessor-failed`",
	"flowReadBackFailModel":     "refuses `flow-read-incomplete`; its writer cannot even seed",
	"flowExpiredTimeoutModel":   "refuses `flow-accessor-timeout`; its reader's deadline is expired before it returns",

	// --- inside the changed class: 0011's own demand-set fixtures ---
	"flowMatchOnlyOwnedModel":     "match-only owned key — the changed class itself (S8's subject)",
	"flowMatchOnlyOwnedSoloModel": "match-only owned key — the changed class itself (S8's subject)",
	"flowEscapeMatchOwnedModel":   "escape-row match-only owned key — DEV-8/DEV-9's subject",
	"flowRefusingSideReaderModel": "a deliberately refusing reader — exit 3, no payload",
}
