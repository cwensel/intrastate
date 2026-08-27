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

	// The shipped 0005 fixtures, in the same full-payload form — the half
	// the spike attested only by suite-green.
	t.Run("shipped-fixture", func(t *testing.T) {
		fx := writeFlowModel(t, flowMVVModel)
		fart := seedArtifact(t, fx, "status=draft")
		fbind := artifactBinding(flowStateRole, fart)

		stdout, _, err := runCmd(t, "flow", "resolve", "--model", fx,
			"--artifact", fbind, "--outcome", "hold", "--as=json")
		if err != nil {
			t.Fatalf("`flow resolve` over the shipped MVV fixture refused: "+
				"%v", err)
		}
		got := normalizeModelPath(stdout, fx)
		if got != resolveFixtureGolden {
			t.Errorf("`flow resolve` over `flowMVVModel` changed.\n"+
				"got:  %s\nwant: %s\nNo shipped flow-reachable fixture is "+
				"in the changed class (A15 iii)", got, resolveFixtureGolden)
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
		`"rule":"prelock","gates":[],"next":{"stage":"prelocked"},` +
		`"writes":{"stage":"prelocked"},"clear":[],"escaped":false}}`,
	"revise": `{"type":"ok","data":{"model":"<model>","revision":"",` +
		`"observed":{},"owned":{"gate_passed":"false","stage":"resolved",` +
		`"status":"draft"},"readers":["rdr-status"],"outcome":"revise",` +
		`"rule":"resolve-route-back","gates":[],"next":{"stage":"refined"},` +
		`"writes":{"stage":"refined"},"clear":[],"escaped":false}}`,
	"abandon": `{"type":"ok","data":{"model":"<model>","revision":"",` +
		`"observed":{},"owned":{"gate_passed":"false","stage":"resolved",` +
		`"status":"draft"},"readers":["rdr-status"],"outcome":"abandon",` +
		`"rule":"resolve-abandon","gates":[],` +
		`"next":{"stage":"dropped","status":"abandoned"},` +
		`"writes":{"stage":"dropped","status":"abandoned"},"clear":[],` +
		`"escaped":false}}`,
}

// resolveFixtureGolden is the same, over the shipped 0005 MVV fixture — the
// half A15's spike attested only by suite-green.
const resolveFixtureGolden = `{"type":"ok","data":{"model":"<model>",` +
	`"revision":"","observed":{},"owned":{"status":"draft"},` +
	`"readers":["state"],"outcome":"hold","rule":"hold-draft","gates":[],` +
	`"next":{"status":"draft"},"writes":{"status":"draft"},"clear":[],` +
	`"escaped":false}}`
