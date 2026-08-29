package cli

// RDR 0023 — `0023:C1`: the `--plan-only` flag surface, the projected
// payload, the report-only differential, the projection site, and text
// mode.
//
// The whole contract turns on one asymmetry, and every oracle here is
// shaped by it: the flag may narrow the REPORT and may change NOTHING
// else. So each oracle below compares the projected run against THE SAME
// REQUEST'S DEFAULT RUN rather than against a literal wherever the clause
// permits, and against an AUTHORED literal exactly where the clause
// demands one ("The projected TOP-LEVEL KEY SET is itself normative: an
// oracle MUST assert it as an explicit literal list").
//
// Two traps the RDR names by hand, avoided by construction here:
//
//  1. The reader-set assertion is NOT written by re-calling
//     `invokedReaders(model, outcome)` — C1 rejects that form as "green by
//     construction … worse than absent". It is written as OBSERVED READER
//     EXECUTION, through the run's own side effects.
//  2. The width assertion is NOT measured on `.data` — C1 fixes the FULL
//     EMITTED LINE, envelope included.

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/spf13/cobra"
)

// --- flag surface and registration (`0023:C1`) ---------------------------

// REQ-1: "flow resolve MUST accept a boolean flag --plan-only, default
// false, long form only (no shorthand)."
// REQ-15: "This scenario also carries the POSITIVE registration-shape
// assertion on the single registrant — `plan-only` is boolean, defaults
// false, no shorthand (C1, `D-naming`) — … the shape is otherwise stated
// twice in this RDR and asserted nowhere."
// REQ-16: "the flag is `--plan-only`, boolean, long form only." — the
// spelling is normative; `--select`, `--json <fields>`/`--fields`,
// `--quiet`, `--no-echo`, `--short`/`--brief` and an inverted `--explain`
// default are all NAMED rejected spellings.
// HAPPY PATH — the positive half, on the `--all` precedent
// (`flow_all_0011_test.go::TestReq36And50_...`).
func TestReq1And15And16_ResolveRegistersABooleanPlanOnlyFlagDefaultFalseNoShorthand(t *testing.T) {
	// Vacuously true of an unregistered verb, so `resolve` must exist
	// before anything about its flag set means anything.
	if !flowSubcommand(t, "resolve") {
		t.Fatal("`flow resolve` is not registered; the flag surface this " +
			"contract fixes is not yet assertable")
	}

	sub := resolveCmd(t)
	if sub == nil {
		t.Fatal("the `flow` group registers no `resolve` verb")
	}

	f := sub.Flags().Lookup(planOnlyFlag)
	if f == nil {
		t.Fatalf("`flow resolve` does not register --%s; the spelling is "+
			"normative (`0023:D-naming`), and %v are NAMED rejected "+
			"alternatives", planOnlyFlag,
			[]string{"--select", "--fields", "--quiet", "--no-echo",
				"--short", "--brief", "--explain"})
	}
	if f.Value.Type() != "bool" {
		t.Errorf("--%s is of type %q; want `bool` — the axis is a FIXED "+
			"partition with exactly two widths, which is why the "+
			"named-field form (`--fields <list>`) was rejected as "+
			"Alternative 1", planOnlyFlag, f.Value.Type())
	}
	if f.DefValue != "false" {
		t.Errorf("--%s defaults to %q; want \"false\" — the echo default is "+
			"AGREED CORRECT and stays, so an inverted default (`--explain` "+
			"to get the echo) is the named Briefly-Rejected alternative: it "+
			"would put a flag on every one-call consumer to preserve "+
			"behaviour that is already right", planOnlyFlag, f.DefValue)
	}
	if f.Shorthand != "" {
		t.Errorf("--%s registers the shorthand %q; `0023:D-naming` reads "+
			"the spelling as the LONG FORM ONLY", planOnlyFlag, f.Shorthand)
	}

	// The rejected spellings must not ship as aliases on this verb. Reusing
	// `--select` in particular would make `0011:D-naming`'s rejection and
	// this RDR's ambiguous.
	for _, rejected := range []string{
		"select", "fields", "json", "quiet", "no-echo", "short", "brief",
		"explain",
	} {
		if sub.Flags().Lookup(rejected) != nil {
			t.Errorf("`flow resolve` registers --%s; it is a NAMED rejected "+
				"spelling for this axis, not an alias", rejected)
		}
	}
}

// REQ-2: "--plan-only MUST NOT be accepted by any other command —
// satisfied by NON-REGISTRATION (registered on resolve only, never the
// shared registrars, the flow group's own or persistent set, or the root's
// persistent set)"
// REQ-3: the negative is "pinned by a STRUCTURAL oracle that walks the
// WHOLE command tree from the root and asserts the set of commands
// registering plan-only is exactly {flow resolve} — closed over verbs added
// later"
// REQ-4: "vacuity-guarded by requiring the four flow verbs to exist."
// REQ-5: "WHOLE means total: the walk descends every child of the root with
// no name-based skip and no Hidden gate, INCLUDING the auto-generated help
// and completion commands."
// REQ-7: "The oracle MUST therefore materialize both auto-generated
// commands on the root before walking (calling cobra's own initializers,
// never hand-constructing a stand-in), and MUST assert their presence in
// the walked set as a precondition."
// REQ-8: "help is present by default, so it cannot distinguish a total
// walker from an untested tree — completion is the discriminating case and
// the control MUST name it."
// REQ-9: "if a cobra upgrade renames, relocates or stops materializing the
// command, the precondition fails LOUDLY and the oracle goes red, instead
// of the walk quietly covering a smaller tree and passing vacuously."
// REQ-13: "The walker yields the WALKED COMMAND SET — every command
// reached, unfiltered — and the registrant filter is applied by the caller,
// not inside the walk."
// REQ-93: "Extend verb-local only | The flag never enters the shared
// registrar (C1)"
// ADVERSARIAL — the one structural-absence oracle in this RDR, deliberately
// stronger than the enumerated sweep `--all` shipped.
func TestReq2And3And4And5And7And8And9And13And93_PlanOnlyIsRegisteredOnExactlyFlowResolveAcrossTheWholeTree(t *testing.T) {
	// --- vacuity guard 1: the four flow verbs (REQ-4) -----------------
	//
	// An absence oracle over an empty tree passes for the wrong reason,
	// and `0023:OD` names this guard as "what stops it passing on an empty
	// command tree".
	if names := flowGroupNames(t); len(names) != len(flowVerbs) {
		t.Fatalf("the `flow` group registers %v; all four verbs must exist "+
			"before the ABSENCE of --%s on the rest of the tree is "+
			"assertable", names, planOnlyFlag)
	}

	root := NewRootCmd()

	// --- materialize the auto-generated commands (REQ-7) --------------
	//
	// Cobra's OWN initializers, called directly on the root — never a
	// hand-constructed stand-in, which would assert against a tree this
	// binary never builds. `0023:A-6` records both as exported and
	// callable from a test.
	root.InitDefaultHelpCmd()
	root.InitDefaultCompletionCmd()

	// --- the TOTAL walk (REQ-5, REQ-13) -------------------------------
	//
	// Yields every command reached, UNFILTERED, with no name-based skip
	// and no Hidden gate. The registrant filter is the caller's, below —
	// which is exactly what makes the preconditions assertable: a walker
	// that returned only registrants could never contain `completion` in a
	// passing build.
	var walked []*cobra.Command
	var walk func(*cobra.Command)
	walk = func(c *cobra.Command) {
		walked = append(walked, c)
		for _, sub := range c.Commands() {
			walk(sub)
		}
	}
	walk(root)

	names := make([]string, 0, len(walked))
	for _, c := range walked {
		names = append(names, c.Name())
	}

	// --- precondition, asserted not assumed (REQ-7, REQ-8, REQ-9) -----
	//
	// `completion` is the DISCRIMINATING member. `help` is force-inited
	// already (`help_all.go::registerHelpAllOnTree`), so its presence
	// proves nothing about totality; a bare `NewRootCmd()` tree carries no
	// `completion` at all. If a cobra upgrade stops materializing it this
	// fails LOUDLY rather than letting the walk cover a smaller tree.
	for _, required := range []string{"help", "completion"} {
		if !slices.Contains(names, required) {
			t.Fatalf("the walked set does not contain %q; walked = %v\n"+
				"C1 widens the walk to the AUTO-GENERATED commands "+
				"precisely because a registration reachable only there is "+
				"what an enumerated sweep misses — and the repo already "+
				"contains one (`help --all`). Without this precondition the "+
				"walk passes VACUOUSLY on the command class the scope "+
				"exists to reach", required, names)
		}
	}

	// --- the assertion: registrants == exactly {flow resolve} ---------
	var registrants []string
	for _, c := range walked {
		if planOnlyRegistered(c) {
			registrants = append(registrants, commandPath(c))
		}
	}
	want := []string{"intrastate flow resolve"}
	slices.Sort(registrants)
	if !slices.Equal(registrants, want) {
		t.Errorf("commands registering --%s = %v; want exactly %v\n"+
			"The negative is satisfied by NON-REGISTRATION: the flag rides "+
			"`resolve` alone, never `registerSelectionFlags`, never the "+
			"`flow` group's own or persistent set, never the root's "+
			"persistent set. This walk is CLOSED OVER VERBS ADDED LATER, "+
			"which is what makes it stronger than the enumerated sibling "+
			"sweep `--all` shipped.\nwalked (%d commands) = %v",
			planOnlyFlag, registrants, want, len(walked), names)
	}
}

// REQ-6: "It therefore MUST NOT reuse
// internal/cli/help_all.go::walkCommandTree, which skips children named
// help or completion before recursing and whose docs.go callers further
// gate on Hidden — that walker asserts a strictly weaker negative than this
// clause requires."
// REQ-14: the new walker "is also a distinct symbol from
// `internal/cli/help_all.go::walkCommandTree` … that one is a visitor
// (`func(root *cobra.Command, fn func(*cobra.Command))`, no return) and
// skips `help`/`completion`, so neither its name nor its shape carries
// over."
// ADVERSARIAL — the prohibition is load-bearing rather than stylistic, and
// this asserts WHY: the shipped walker demonstrably cannot reach the scope
// C1 mandates, so a build that reused it would assert a weaker negative
// while looking green.
func TestReq6And14_TheShippedPartialWalkerCannotReachC1sScope(t *testing.T) {
	root := NewRootCmd()
	root.InitDefaultHelpCmd()
	root.InitDefaultCompletionCmd()

	// The TOTAL walk C1 mandates.
	var total []string
	var walk func(*cobra.Command)
	walk = func(c *cobra.Command) {
		total = append(total, commandPath(c))
		for _, sub := range c.Commands() {
			walk(sub)
		}
	}
	walk(root)

	// The SHIPPED partial walker, over the SAME materialized root.
	var partial []string
	walkCommandTree(root, func(c *cobra.Command) {
		partial = append(partial, commandPath(c))
	})

	if len(partial) >= len(total) {
		t.Fatalf("the shipped `walkCommandTree` reached %d commands and the "+
			"total walk %d; the prohibition on reuse rests on the shipped "+
			"walker reaching STRICTLY FEWER (it skips `help`/`completion` "+
			"before recursing). If they now agree, C1's reuse prohibition "+
			"needs re-deriving rather than silently satisfying.\n"+
			"total   = %v\npartial = %v", len(partial), len(total),
			total, partial)
	}

	// The dropped set must include the `completion` subtree — the six
	// commands a total sweep exists to cover (`0023:A5`: 15 vs 9).
	var dropped []string
	for _, name := range total {
		if !slices.Contains(partial, name) {
			dropped = append(dropped, name)
		}
	}
	var droppedCompletion bool
	for _, name := range dropped {
		if strings.Contains(name, "completion") {
			droppedCompletion = true
		}
	}
	if !droppedCompletion {
		t.Errorf("the shipped walker drops %v, none of them under "+
			"`completion`; the reuse prohibition names the completion "+
			"subtree as the concrete loss", dropped)
	}

	// And the new oracle must not BE that walker: a total walk reaches
	// `completion`, so any oracle whose walked set omits it is the weaker
	// negative wearing the stronger clause's name.
	if slices.ContainsFunc(partial, func(s string) bool {
		return strings.Contains(s, "completion")
	}) {
		t.Errorf("`walkCommandTree` reached a `completion` command: %v\n"+
			"Its documented behaviour is to SKIP children named `help` or "+
			"`completion` before recursing", partial)
	}
}

// REQ-81, S4's negative controls: "register `plan-only` on a second
// command; S4 must go red. Second control: register it on the
// auto-generated `completion` command — S4 must go red there too, which is
// the control a bare `NewRootCmd()` tree cannot run" — (OD, `0023:S4` row)
// ADVERSARIAL — the discriminability of the whole-tree walk, asserted by
// RUNNING both controls rather than by describing them. A registration
// oracle that cannot fail is worse than absent, and the second control is
// the one that separates a total walker from a partial one: it plants the
// flag on a command a bare `NewRootCmd()` tree does not even contain.
func TestReq81_TheWholeTreeWalkGoesRedWhenASecondCommandOrCompletionRegistersTheFlag(t *testing.T) {
	// The registrant-set computation under test, factored so both controls
	// drive exactly the walk the S4 oracle performs.
	registrantsOf := func(root *cobra.Command) []string {
		var out []string
		var walk func(*cobra.Command)
		walk = func(c *cobra.Command) {
			if planOnlyRegistered(c) {
				out = append(out, commandPath(c))
			}
			for _, sub := range c.Commands() {
				walk(sub)
			}
		}
		walk(root)
		slices.Sort(out)
		return out
	}

	want := []string{"intrastate flow resolve"}

	// BASELINE: the untouched tree must be at the contract's value.
	// Without this the two controls below "go red" trivially, proving
	// nothing about discrimination.
	base := NewRootCmd()
	base.InitDefaultHelpCmd()
	base.InitDefaultCompletionCmd()
	if got := registrantsOf(base); !slices.Equal(got, want) {
		t.Fatalf("the untouched tree's registrant set = %v; want %v. The "+
			"controls below discriminate only against a correct baseline",
			got, want)
	}

	t.Run("control_1_a_second_ordinary_command_registers_the_flag", func(t *testing.T) {
		root := NewRootCmd()
		root.InitDefaultHelpCmd()
		root.InitDefaultCompletionCmd()
		for _, c := range root.Commands() {
			if c.Name() == "version" {
				c.Flags().Bool(planOnlyFlag, false, "planted by the S4 control")
			}
		}
		if got := registrantsOf(root); slices.Equal(got, want) {
			t.Errorf("a SECOND command registers --%s and the registrant "+
				"set is still %v; the oracle cannot see a widened surface "+
				"and is green by construction", planOnlyFlag, got)
		}
	})

	t.Run("control_2_the_auto_generated_completion_command_registers_the_flag", func(t *testing.T) {
		root := NewRootCmd()
		root.InitDefaultHelpCmd()
		root.InitDefaultCompletionCmd()

		var planted bool
		for _, c := range root.Commands() {
			if c.Name() == "completion" {
				c.Flags().Bool(planOnlyFlag, false, "planted by the S4 control")
				planted = true
			}
		}
		if !planted {
			t.Fatal("no `completion` command to plant the flag on; this is " +
				"THE control a bare `NewRootCmd()` tree cannot run, which " +
				"is exactly why the oracle materializes cobra's own " +
				"initializers before walking")
		}
		if got := registrantsOf(root); slices.Equal(got, want) {
			t.Errorf("the AUTO-GENERATED `completion` command registers "+
				"--%s and the registrant set is still %v; a registration "+
				"reachable only on an auto-generated command is precisely "+
				"what an enumerated sweep misses — the repo already contains "+
				"one (`help --all`) that the shipped 0011 oracle does not "+
				"see because it descends the `flow` group only",
				planOnlyFlag, got)
		}
	})
}

// REQ-92: "Register `--plan-only` on `resolve` only and hand `respond.OK`
// the projected verb result when set — the mechanism A2 verified, applied
// after payload assembly, before the gateway." — (PH, Phase 1)
// HAPPY PATH — Phase 1's whole deliverable in one assertion: the flag
// exists on exactly one verb, and setting it changes the reported width
// without changing anything else. The two halves are stated together
// because either alone is satisfiable by a build that does not work.
func TestReq92_TheFlagRidesResolveAloneAndHandsRespondOKTheProjectedResult(t *testing.T) {
	sub := resolveCmd(t)
	if sub == nil || sub.Flags().Lookup(planOnlyFlag) == nil {
		t.Fatalf("`flow resolve` does not register --%s", planOnlyFlag)
	}

	def := requireSuccess(t, append(pricingCall(t), "--as=json")...)
	proj := requireSuccess(t,
		append(pricingCall(t, "--"+planOnlyFlag), "--as=json")...)

	// The flag was HANDED to `respond.OK` — observable as a narrower
	// envelope on the same success path, with the envelope's own shape
	// untouched (`{"type":"ok","data":{…}}`).
	if def == proj {
		t.Fatalf("the payload is identical ± the flag; setting it must hand "+
			"`respond.OK` the PROJECTED verb result\n%s", def)
	}
	for _, out := range []string{def, proj} {
		if !strings.HasPrefix(strings.TrimSpace(out), `{"type":"ok","data":{`) {
			t.Errorf("the emitted record does not carry the shipped success "+
				"envelope; the projection lands on the verb result BEFORE "+
				"`respond.OK`, so the gateway's own shape is untouched:\n%s",
				out)
		}
	}
	if len(emittedLine(t, proj)) >= len(emittedLine(t, def)) {
		t.Errorf("the projected record is not narrower than the default; " +
			"the mechanism is applied AFTER payload assembly and BEFORE the " +
			"gateway, and narrowing is its only observable effect")
	}
}

// REQ-12: "`--plan-only` on `next`/`read-state`/`set-state` fails the parse
// — `command-error`, exit 2, flag name only in pflag's free text (the
// `0011:A14` shared-bucket limit, accepted again here)." — [0011-carried]
// REQ-11: "This contract mints no new refusal code and no new exit group."
// REQ-10: "The behavioural command-error run is corroboration, never the
// assertion."
// REQ-132: "a binary predating this RDR refuses `--plan-only` at the parse
// (`command-error`, exit 2), so a caller can never hold a full-width
// payload believing it was projected" — the SAME refusal is what makes
// version skew loud in both directions.
// ADVERSARIAL — written so it cannot be mistaken for the structural
// assertion: it asserts the CLASS and the exit, never pflag's free text,
// which `0011:A14` records as text this repo does not own.
func TestReq10And11And12And132_PlanOnlyOnASiblingVerbIsTheSharedUsageBucketAndMintsNoCode(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	// CONTROL — without this the whole test is tautological. A tree where
	// the flag exists NOWHERE refuses it on every verb for the wrong
	// reason: cobra's unknown-flag path, not this contract's
	// non-registration. `0011:A14` makes the code a SHARED bucket, so the
	// code alone cannot attribute the refusal to this flag — only the
	// contrast with a verb that ACCEPTS it can.
	//
	// So the sibling refusals below mean something only once `resolve`
	// accepts the flag. This is the same discipline
	// `flow_all_0011_test.go` applies by asserting the positive
	// registration in a separate oracle.
	sub := resolveCmd(t)
	if sub == nil || sub.Flags().Lookup(planOnlyFlag) == nil {
		t.Fatalf("`flow resolve` does not register --%s; until it does, a "+
			"sibling verb's refusal is cobra's unknown-flag path rather "+
			"than this contract's NON-REGISTRATION, and the shared "+
			"`command-error` bucket cannot tell the two apart. The "+
			"asymmetry the clause fences — accepted on `resolve`, refused "+
			"everywhere else — is not yet observable", planOnlyFlag)
	}
	// And the positive half of the same asymmetry: `resolve` ACCEPTS it.
	if _, _, err := runCmd(t, append(pricingCall(t, "--"+planOnlyFlag),
		"--as=json")...); err != nil {
		t.Fatalf("`flow resolve --%s` refused: %v; the flag is accepted "+
			"there and refused everywhere else, and both halves are needed "+
			"for the refusals below to be about THIS flag", planOnlyFlag, err)
	}

	for _, verb := range []struct {
		name string
		args []string
	}{
		{"next", []string{"flow", "next", "--model", model,
			"--artifact", bind, "--" + planOnlyFlag, "--as=json"}},
		{"read-state", []string{"flow", "read-state", "--model", model,
			"--artifact", bind, "--" + planOnlyFlag, "--as=json"}},
		{"set-state", []string{"flow", "set-state", "--model", model,
			"--artifact", bind, "--write", "status=draft",
			"--" + planOnlyFlag, "--as=json"}},
	} {
		t.Run(verb.name, func(t *testing.T) {
			stdout, _, err := runCmd(t, verb.args...)
			if err == nil {
				t.Fatalf("`flow %s --%s` SUCCEEDED; the flag is registered "+
					"on `resolve` only\nstdout:\n%s",
					verb.name, planOnlyFlag, stdout)
			}
			if code := clierr.ErrorCode(err); code != "command-error" {
				t.Errorf("code = %q; want %q — pflag fails the parse before "+
					"any RunE, so no `flow-*` code is reachable and this "+
					"contract mints NONE", code, "command-error")
			}
			if exit := clierr.ExitCodeFor(err); exit != 2 {
				t.Errorf("exit = %d; want 2 — this contract mints no new "+
					"exit group either", exit)
			}
			// No bespoke code may appear on the wire. `0023:BR5` names
			// "registering the flag but having other verbs refuse it with a
			// typed code" as REJECTED.
			for _, minted := range []string{
				"flow-plan-only", "flow-projection", "plan-only-unsupported",
			} {
				if strings.Contains(stdout, minted) {
					t.Errorf("a bespoke %q code appears on the wire:\n%s\n"+
						"C1 follows `0011:C2`'s settled disposition — "+
						"non-registration with the shared `command-error` "+
						"bucket — and mints no new refusal code", minted, stdout)
				}
			}
		})
	}
}

// REQ-29: "No clause of 0011 is overridden: `0011:C2`'s fence is `--all`'s
// non-registration on the other three verbs, which this RDR leaves intact,
// and its `--select` rejected-spelling pin is on `flow next`" —
// [0011-carried]
// REQ-30: "`--all` probes are exact-name `Lookup(\"all\")`; a resolve-local
// `plan-only` is invisible to them." — [0011-carried]
// DOMAIN EDGE — the boundary of the override in the OTHER direction: adding
// this flag must not disturb the sibling fence, and the sibling fence must
// not have been widened to cover this flag.
func TestReq29And30_The0011AllFenceIsUntouchedAndInvisibleToPlanOnly(t *testing.T) {
	root := NewRootCmd()

	// `--all` still rides `next` alone, exactly as `0011:C2` fences it.
	var checked bool
	for _, c := range root.Commands() {
		if c.Name() != "flow" {
			continue
		}
		if c.Flags().Lookup("all") != nil || c.PersistentFlags().Lookup("all") != nil {
			t.Error("the `flow` group now registers --all of its own; " +
				"`0011:C2` places it on `next` ONLY and this RDR overrides " +
				"no clause of 0011")
		}
		for _, sub := range c.Commands() {
			if sub.Name() == "next" {
				checked = true
				if sub.Flags().Lookup("all") == nil {
					t.Error("`flow next` no longer registers --all; " +
						"`0011:C2`'s positive is untouched by this RDR")
				}
				// `flow next` must NOT acquire the projection flag: the
				// registrant set is exactly {flow resolve}.
				if sub.Flags().Lookup(planOnlyFlag) != nil {
					t.Errorf("`flow next` registers --%s; the projection "+
						"axis rides `resolve` alone", planOnlyFlag)
				}
				continue
			}
			if sub.Flags().Lookup("all") != nil {
				t.Errorf("`flow %s` registers --all; `0011:C2`'s structural "+
					"negative is left INTACT by this RDR", sub.Name())
			}
		}
	}
	if !checked {
		t.Fatal("the `flow` group registers no `next` verb")
	}

	// The `--select` rejected-spelling pin is on `flow next` and stays
	// there; this RDR reuses neither the spelling nor the pin.
	sub := resolveCmd(t)
	if sub == nil {
		t.Fatal("the `flow` group registers no `resolve` verb")
	}
	if sub.Flags().Lookup("select") != nil {
		t.Error("`flow resolve` registers --select; even though `0011`'s " +
			"pin test is scoped to `next`, reusing the spelling for a " +
			"different axis would make the two rejections ambiguous")
	}

	// And the exact-name probe direction: a resolve-local `plan-only` is
	// invisible to `Lookup("all")`, so the two fences cannot collide.
	if sub.Flags().Lookup("all") != nil {
		t.Errorf("`flow resolve` registers --all; `0011:C2` fences its " +
			"absence there and this RDR does not reopen it")
	}
}

// --- the projected payload (`0023:C1`, `0023:S2`) ------------------------

// REQ-17: "Under --plan-only a SUCCESS payload MUST carry exactly the plan
// group — rule, gates, emit, next, writes, clear, escaped, escape_class,
// revision, plus every field a successor assigns to the plan group under C2"
// REQ-18: "MUST NOT carry the echo group: model, observed, owned, readers,
// outcome."
// REQ-22: "The projected TOP-LEVEL KEY SET is itself normative: an oracle
// MUST assert it as an explicit literal list, so a later payload field
// reaches the projected width only by a conscious C2 assignment, never by
// falling through an omit-list."
// REQ-23: the literal is "exactly `revision`, `rule`, `gates`, `emit`,
// `next`, `writes`, `clear`, `escaped`, plus `escape_class` exactly when
// the same request's DEFAULT output carries it"
// REQ-98 / MVV step 2: "exactly the normative projected key set … every
// carried field byte-identical to step 1's;
// `model`/`observed`/`owned`/`readers`/`outcome` absent (not null, not
// empty)."
// REQ-108: "Normative fixture: the projected wire record of the pricing 2×2
// call, whose exact bytes are the \"Projected reference\" line in
// `evidence/spikes/a2-encoder-mechanism.md` §Reference output"; S2 pins the
// same fixture — which is why this oracle runs `pricingCall` rather than a
// shape authored here.
// HAPPY PATH — the explicit-literal oracle, on the NORMATIVE FIXTURE (the
// pricing 2×2 call, `0023:S2`).
func TestReq17And18And22And23And98And108_TheProjectedKeySetIsExactlyTheAuthoredLiteral(t *testing.T) {
	def := requireSuccess(t, append(pricingCall(t), "--as=json")...)
	proj := requireSuccess(t, append(pricingCall(t, "--"+planOnlyFlag), "--as=json")...)

	defFields := rawFieldsOf(t, def)
	projFields := rawFieldsOf(t, proj)

	// The literal, per `0023:A-2`: eight names, plus `escape_class` exactly
	// when the SAME REQUEST'S DEFAULT output carried it (REQ-25 — never an
	// unconditional "when escaped").
	want := slices.Clone(projectedKeyLiteral0023)
	if _, carried := defFields["escape_class"]; carried {
		want = append(want, "escape_class")
	}
	slices.Sort(want)

	got := sortedKeys(projFields)
	if !slices.Equal(got, want) {
		t.Errorf("projected top-level key set = %v;\nwant                       "+
			"%v\nThe list is an AUTHORED LITERAL, not a complement: a later "+
			"payload field reaches the projected width only by a conscious "+
			"C2 assignment, never by falling through an omit-list.\n"+
			"default keys = %v", got, want, sortedKeys(defFields))
	}

	// The echo group, named individually so a failure says WHICH field
	// leaked rather than only that a set differed.
	for _, key := range echoGroup0023 {
		if _, held := projFields[key]; held {
			t.Errorf("the projected payload carries %q; it is an ECHO-group "+
				"field (`0023:C2`) and the flag's entire purpose is to omit "+
				"it — the caller already holds every one of these",
				key)
		}
	}
}

// REQ-19: "An omitted key is ABSENT, never null and never an empty
// placeholder." — restated at `0023:D-wire-byte-format`: "Absent means
// absent — never `null`, `{}`, or `\"\"` stand-ins for a projected key."
// REQ-113: "the explicit projected-key-set oracle (S2), which CARRIES the
// absent-not-null assertion — a projected key must be absent, never `null`
// and never an empty placeholder, the live hazard A2's spike reproduced
// with bare non-pointer `omitempty`"
// REQ-79 (S2's absence control): "render one echo key as `null` … S2 must
// go red where a key-presence-only check would stay green"
// ADVERSARIAL — the absence control's own arm. A key-presence-only check
// would stay GREEN against a build that renders `"owned":null`, which is
// the exact hazard `0023:A2`'s spike reproduced with bare non-pointer
// `omitempty`. This asserts on the RAW BYTES, so `null` is caught.
func TestReq19And79And113_AProjectedKeyIsAbsentNotNullNotAnEmptyPlaceholder(t *testing.T) {
	proj := requireSuccess(t, append(pricingCall(t, "--"+planOnlyFlag), "--as=json")...)
	line := emittedLine(t, proj)
	fields := rawFieldsOf(t, proj)

	for _, key := range echoGroup0023 {
		raw, held := fields[key]
		if !held {
			continue // absent — correct
		}
		// Present. Name the placeholder so the message discriminates the
		// three failure shapes rather than reporting one generic error.
		shape := string(raw)
		switch strings.TrimSpace(shape) {
		case "null":
			t.Errorf("projected key %q renders as `null`; absent means "+
				"ABSENT. This is the live hazard `0023:A2`'s spike "+
				"reproduced — a nil value under bare `omitempty` — and a "+
				"key-presence-only check would pass right through it",
				key)
		case "{}", "[]", `""`:
			t.Errorf("projected key %q renders as the empty placeholder %s; "+
				"`0023:D-wire-byte-format` forbids `null`, `{}` and `\"\"` "+
				"stand-ins for a projected key", key, shape)
		default:
			t.Errorf("projected key %q is PRESENT as %s; the echo group is "+
				"omitted under the flag", key, shape)
		}
	}

	// The wire bytes carry no `null` at all on this fixture: every carried
	// value is a string, a bool, or a non-nil container.
	if strings.Contains(line, ":null") {
		t.Errorf("the projected emitted line carries a `null` value:\n%s\n"+
			"Nothing in the plan group renders `null` — `emit` is `{}` "+
			"(`0010:C4`), the containers are allocated unconditionally "+
			"(`0023:A9`) — so a `null` here is a projection defect", line)
	}
}

// REQ-20: "Every carried field's rendering MUST be byte-identical to its
// default-mode rendering (same encoder, HTML escaping disabled)"
// REQ-26: "projection is key deletion, never re-encoding: the projected
// payload is the default payload minus the echo keys, byte-for-byte on
// every surviving field, keys in surviving declaration order, same
// non-HTML-escaping encoder."
// REQ-124: "the projection deletes whole keys and never rewrites a carried
// value."
// REQ-127: "one wire encoder
// (`internal/cli/clierr/clierr.go::WriteJSONLine` — `json.Encoder`,
// `SetEscapeHTML(false)`, compact one-line NDJSON)"
// REQ-78 (S1's mutation control): "mutate one carried value under the flag"
// — this oracle is what goes red.
// ADVERSARIAL — compared on RAW BYTES, not decoded values: a re-encode that
// HTML-escapes `<` to `<` decodes EQUAL and renders differently, and a
// decode-then-compare would launder exactly the defect this clause names.
func TestReq20And26And124And127_EveryCarriedFieldIsByteIdenticalToItsDefaultRendering(t *testing.T) {
	// The release grammar carries a set-valued tag and a `<clear>` plan, so
	// the carried values are richer than the pricing table's — including
	// the canonical set literal whose `<`/`&` members are exactly what a
	// re-encode would corrupt.
	for _, shape := range projectionShapes0023(t) {
		t.Run(shape.name, func(t *testing.T) {
			def := requireSuccess(t, append(shape.args, "--as=json")...)
			proj := requireSuccess(t,
				append(slices.Clone(shape.args), "--"+planOnlyFlag, "--as=json")...)

			defFields := rawFieldsOf(t, def)
			projFields := rawFieldsOf(t, proj)

			for key, projRaw := range projFields {
				defRaw, held := defFields[key]
				if !held {
					t.Errorf("the projected payload carries %q, which the "+
						"SAME REQUEST's default payload does not; the "+
						"projected key set is a strict SUBSET of the "+
						"default's — projection is key DELETION", key)
					continue
				}
				if string(projRaw) != string(defRaw) {
					t.Errorf("carried field %q differs between widths:\n"+
						"  default   = %s\n  projected = %s\n"+
						"Projection is key deletion, NEVER re-encoding: the "+
						"projected payload is the default payload minus the "+
						"echo keys, byte-for-byte on every surviving field, "+
						"through the same non-HTML-escaping encoder",
						key, defRaw, projRaw)
				}
			}
		})
	}
}

// REQ-21: "the carried keys keep their default-mode relative order."
// REQ-121: "Carried fields keep their default-mode relative order."
// REQ-122: "JSON key order is struct declaration order, preserved under the
// pointer conversion."
// A-1: "the projected order is the default order minus the deleted keys,
// never a re-sort" — which puts `revision` FIRST.
// BOUNDARY — order is the property a map-based or re-sorting implementation
// silently loses while every key-SET assertion stays green.
func TestReq21And121And122_TheCarriedKeysKeepTheirDefaultModeRelativeOrder(t *testing.T) {
	def := requireSuccess(t, append(pricingCall(t), "--as=json")...)
	proj := requireSuccess(t, append(pricingCall(t, "--"+planOnlyFlag), "--as=json")...)

	defOrder := payloadKeyOrder(t, def)
	projOrder := payloadKeyOrder(t, proj)

	// The default order minus the echo keys IS the expected projected
	// order. Derived from the same run's DEFAULT output, so this does not
	// restate the implementation.
	var want []string
	for _, k := range defOrder {
		if !slices.Contains(echoGroup0023, k) {
			want = append(want, k)
		}
	}
	if !slices.Equal(projOrder, want) {
		t.Errorf("projected key ORDER = %v;\nwant                 %v\n"+
			"The projected order is the DEFAULT order minus the deleted "+
			"keys, never a re-sort. A map-valued or re-serialising "+
			"projection loses this while every key-SET assertion stays "+
			"green.\ndefault order = %v", projOrder, want, defOrder)
	}

	// `revision` FIRST is the concrete witness `0023:A-1` records and the
	// desk trace pins byte-for-byte
	// (`{"revision":"","rule":"paid-eu","gates":[],…}`).
	if len(projOrder) == 0 || projOrder[0] != "revision" {
		t.Errorf("the projected payload's FIRST key is %v; want `revision` "+
			"— deleting the five echo keys from the shipped declaration "+
			"order leaves `revision` at the head", projOrder)
	}
}

// REQ-24: "Presence-rule fields keep their own contracts inside the plan
// group, unchanged and not restated here: emit stays present as {}
// (0010:C4), and escape_class keeps whatever presence rule its producer
// already has (0005:A-3 — omitted when unescaped, and also when an escaped
// row's class is unprobeable)." — [0010-carried], [0005-carried]
// REQ-25: "This projection neither widens nor narrows those rules: a field
// the default mode omits is omitted under the flag for the same reason, so
// an oracle MUST assert presence-rule fields by comparison against the SAME
// run's default output, never against an unconditional literal."
// REQ-123: "`emit` stays present as `{}` per `0010:C4`" — [0010-carried]
// DOMAIN EDGE — the two presence-rule fields, on BOTH arms: an escaped plan
// (class present) and an ordinary one (class absent).
func TestReq24And25And123_PresenceRuleFieldsAreAssertedAgainstTheSameRunsDefaultNotALiteral(t *testing.T) {
	escapeModel := writeFlowModel(t, dtEscapeModel0010)
	plainModel := writeFlowModel(t, dtModel0010)

	for _, arm := range []struct {
		name        string
		args        []string
		wantEscaped bool
	}{
		{
			// (y,q) is covered by no ordinary row, so `otherwise` rescues
			// it and the plan IS escaped — `escape_class` present.
			name: "escaped/class-present",
			args: resolveArgs(escapeModel, "", "decide",
				"--tag", "a=y", "--tag", "b=q"),
			wantEscaped: true,
		},
		{
			// An ordinary selection: `escape_class` is OMITTED by its
			// producer's own rule (`0005:A-3`), in BOTH widths.
			name: "ordinary/class-absent",
			args: resolveArgs(plainModel, "", "decide",
				"--tag", "a=x", "--tag", "b=p"),
			wantEscaped: false,
		},
	} {
		t.Run(arm.name, func(t *testing.T) {
			def := requireSuccess(t, append(slices.Clone(arm.args), "--as=json")...)
			proj := requireSuccess(t,
				append(slices.Clone(arm.args), "--"+planOnlyFlag, "--as=json")...)

			defFields := rawFieldsOf(t, def)
			projFields := rawFieldsOf(t, proj)

			// Control: the arm is the shape it claims to be. Without this
			// the assertion below passes on a build that never escapes.
			_, defHasClass := defFields["escape_class"]
			if defHasClass != arm.wantEscaped {
				t.Fatalf("the DEFAULT run carries escape_class = %v; this "+
					"arm exists to exercise %v. Without the control the "+
					"projected assertion proves nothing",
					defHasClass, arm.wantEscaped)
			}

			// The clause: presence under the flag == presence by default,
			// for the SAME request. Never an unconditional literal.
			_, projHasClass := projFields["escape_class"]
			if projHasClass != defHasClass {
				t.Errorf("escape_class is %s under the flag and %s by "+
					"default on the SAME request; the projection neither "+
					"widens nor narrows a producer's presence rule — a "+
					"field the default mode omits is omitted under the flag "+
					"FOR THE SAME REASON",
					presence(projHasClass), presence(defHasClass))
			}

			// `emit` is the other presence-rule field, and its rule is
			// UNCONDITIONAL presence as `{}` (`0010:C4`) — so it rides both
			// widths always.
			for width, fields := range map[string]map[string]json.RawMessage{
				"default": defFields, "projected": projFields,
			} {
				raw, held := fields["emit"]
				if !held {
					t.Errorf("the %s payload carries no `emit`; `0010:C4` "+
						"fixes it PRESENT — as `{}` when the selected row "+
						"authored none, never omitted and never `null`",
						width)
					continue
				}
				if strings.TrimSpace(string(raw)) == "null" {
					t.Errorf("`emit` renders `null` in the %s payload; "+
						"`0010:C4` fixes `{}`", width)
				}
			}
		})
	}
}

func presence(b bool) string {
	if b {
		return "PRESENT"
	}
	return "ABSENT"
}

// --- the report-only differential (`0023:C1`, `0023:S1`) -----------------

// REQ-31: "--plan-only is REPORT-ONLY, and report-only is ORACLE-ENFORCED,
// not aspirational: it MUST NOT change what is decided or how the run
// fails"
// REQ-32: the differential asserts "identical exit codes"
// REQ-34: "the projected success payload a strict key-subset of the default
// payload with byte-identical values on every carried key"
// REQ-45: "the same invocation ± the flag selects the same rule, runs the
// same gates, and refuses identically — only report width differs … an
// oracle asserts the ±-flag plan-group equality."
// REQ-130: "a projected payload missing a plan-group field (mechanism bug)
// fails the ±-flag equality oracle"
// REQ-133: "`--as=json` ± `--plan-only` over the same request diffs to
// exactly the echo keys; any other diff is a defect in the projection."
// HAPPY PATH — the differential's core: the diff is EXACTLY the echo keys,
// in both directions.
func TestReq31And32And34And45And130And133_TheDifferentialDiffsToExactlyTheEchoKeys(t *testing.T) {
	for _, shape := range projectionShapes0023(t) {
		t.Run(shape.name, func(t *testing.T) {
			defOut, _, defErr := runCmd(t, append(slices.Clone(shape.args), "--as=json")...)
			projOut, _, projErr := runCmd(t,
				append(slices.Clone(shape.args), "--"+planOnlyFlag, "--as=json")...)

			// Identical exit codes (REQ-32). Asserted before anything else:
			// a flag that changed the disposition would make every payload
			// comparison below meaningless.
			if got, want := clierr.ExitCodeFor(projErr), clierr.ExitCodeFor(defErr); got != want {
				t.Fatalf("exit code = %d under the flag and %d by default; "+
					"--%s is REPORT-ONLY and MUST NOT change what is decided "+
					"or how the run fails\ndefault:\n%s\nprojected:\n%s",
					got, want, planOnlyFlag, defOut, projOut)
			}
			if defErr != nil {
				t.Fatalf("the shape %q did not succeed: %v\n%s",
					shape.name, defErr, defOut)
			}

			defFields := rawFieldsOf(t, defOut)
			projFields := rawFieldsOf(t, projOut)

			// Strict key-SUBSET (REQ-34): nothing appears under the flag
			// that the default width does not carry.
			for key := range projFields {
				if _, held := defFields[key]; !held {
					t.Errorf("the projected payload carries %q, absent from "+
						"the default payload; the projected key set is a "+
						"STRICT SUBSET — a projection can only delete", key)
				}
			}

			// The diff is EXACTLY the echo keys (REQ-133). Both directions:
			// a plan-group field silently dropped is as much a defect as an
			// echo field leaking through.
			var missing []string
			for key := range defFields {
				if _, held := projFields[key]; !held {
					missing = append(missing, key)
				}
			}
			slices.Sort(missing)
			want := slices.Clone(echoGroup0023)
			slices.Sort(want)
			if !slices.Equal(missing, want) {
				t.Errorf("the keys the flag DROPS = %v;\nwant             "+
					"      %v\n`--as=json` ± the flag over the same request "+
					"diffs to EXACTLY the echo keys; any other diff is a "+
					"defect in the projection — a plan-group field missing "+
					"here is the mechanism bug `0023:F2` names",
					missing, want)
			}
		})
	}
}

// REQ-33: the differential asserts "byte-identical refusal envelopes
// (CLIError, findings — never projected, and carrying no echo group to
// project)"
// REQ-100 / MVV step 4: "Re-run with an unrecognized outcome, ± the flag →
// byte-identical refusal envelopes and exit codes."
// REQ-131: "a projection accidentally applied to a refusal path would
// change refusal bytes — A6 plus a refusal-unchanged oracle make it a test
// failure, not a field report."
// A-7: "byte-identical refusal envelopes" is asserted on the full emitted
// refusal LINE, the same unit as the success-width clause.
// ADVERSARIAL — refusals are where a defensive guard or a misplaced
// projection site would show, and `0023:A6` says by construction that not
// one byte may move.
func TestReq33And100And131_RefusalEnvelopesAreByteIdenticalPlusOrMinusTheFlag(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)
	ambiguous := writeFlowModel(t, flowAmbiguousModel)
	ambigArt := seedArtifact(t, ambiguous, "status=draft")

	for _, refusal := range []struct {
		name string
		args []string
		code string
	}{
		{
			// MVV step 4's own shape.
			name: "unmodeled-outcome",
			args: resolveArgs(model, bind, "not-an-outcome"),
			code: "flow-unmodeled-outcome",
		},
		{
			// A refusal reached AFTER the readers have already run — the
			// census notes six of nine Fail sites are downstream of the
			// reader pass, so this exercises the path where a projection
			// site placed too early would be reachable.
			name: "ambiguous-match",
			args: resolveArgs(ambiguous,
				artifactBinding(flowStateRole, ambigArt), "advance"),
			code: "flow-ambiguous-match",
		},
		{
			// A CLI-side refusal reached BEFORE the readers: the earliest
			// Fail site, where a flag consulted in refusal construction
			// would be most visible.
			name: "tag-invalid/empty-outcome",
			args: resolveArgs(model, bind, ""),
			code: "flow-tag-invalid",
		},
	} {
		t.Run(refusal.name, func(t *testing.T) {
			defOut, _, defErr := runCmd(t, append(slices.Clone(refusal.args), "--as=json")...)
			projOut, _, projErr := runCmd(t,
				append(slices.Clone(refusal.args), "--"+planOnlyFlag, "--as=json")...)

			// Control: this really is the refusal it claims to be. Without
			// it the byte comparison below would pass on two identical
			// SUCCESSES, or on two identical parse errors from an
			// unregistered flag.
			if defErr == nil {
				t.Fatalf("the default run SUCCEEDED; this entry exists to "+
					"exercise the %s refusal\n%s", refusal.code, defOut)
			}
			if code := clierr.ErrorCode(defErr); code != refusal.code {
				t.Fatalf("default refusal code = %q; want %q", code, refusal.code)
			}
			if code := clierr.ErrorCode(projErr); code != refusal.code {
				t.Errorf("refusal code under the flag = %q; want %q — a "+
					"refusal is FLAG-BLIND by construction (`0023:A6`): "+
					"`respond.Fail` renders `clierr.CLIError` with no "+
					"request-echo member, and C1 forbids the flag from being "+
					"readable on that path at all",
					clierr.ErrorCode(projErr), refusal.code)
			}

			// Byte identity on the FULL EMITTED LINE (`0023:A-7`), the same
			// unit the success-width clause fixes.
			if got, want := strings.TrimRight(projOut, "\n"),
				strings.TrimRight(defOut, "\n"); got != want {
				t.Errorf("the refusal envelope differs ± the flag:\n"+
					"  default   = %s\n  projected = %s\n"+
					"--%s cannot change a single refusal byte: the refusal "+
					"envelope's fields are code/message/param/detail/hint/"+
					"findings, with NO echo-group member to project",
					want, got, planOnlyFlag)
			}
			if got, want := clierr.ExitCodeFor(projErr), clierr.ExitCodeFor(defErr); got != want {
				t.Errorf("refusal exit = %d under the flag and %d by "+
					"default; the flag MUST NOT change how the run fails",
					got, want)
			}
		})
	}
}

// REQ-36: the differential asserts "the projected encoding STRICTLY SHORTER
// than the default. The width assertion is not implied by the subset one —
// a later plan-group field can grow the projected payload past today's full
// width with the key-set and partition oracles still green — so the flag's
// whole reason for existing is itself oracle-enforced."
// REQ-38: "STRICTLY SHORTER is measured on the FULL EMITTED LINE — the
// complete {\"type\":\"ok\",\"data\":{…}} NDJSON record as written,
// excluding the trailing newline — not on the .data payload alone, and it
// is asserted in BOTH output modes."
// REQ-111: the width assertion is "unconditional because the five echo keys
// always render, `model`/`outcome` being non-`omitempty` and the three
// containers rendering `{}`/`[]`"
// A-5: "S1 asserts `len(projectedLine) < len(defaultLine)` with no
// threshold; the 40% bar is MVV-only and is not an oracle."
// BOUNDARY — strict inequality, on the envelope-inclusive unit, in BOTH
// modes. This is the clause that fails on `0023:P-12`'s "bigger output,
// every test green" scenario.
func TestReq36And38And111_TheProjectedEncodingIsStrictlyShorterOnTheFullEmittedLineInBothModes(t *testing.T) {
	for _, shape := range projectionShapes0023(t) {
		t.Run(shape.name, func(t *testing.T) {
			t.Run("json", func(t *testing.T) {
				def := emittedLine(t,
					requireSuccess(t, append(slices.Clone(shape.args), "--as=json")...))
				proj := emittedLine(t, requireSuccess(t,
					append(slices.Clone(shape.args), "--"+planOnlyFlag, "--as=json")...))

				if len(proj) >= len(def) {
					t.Errorf("projected line is %d B and default %d B; "+
						"STRICTLY SHORTER is the flag's whole reason for "+
						"existing and is oracle-enforced. The unit is the "+
						"FULL EMITTED LINE including the "+
						"{\"type\":\"ok\",\"data\":{…}} envelope — the unit "+
						"the caller actually pays for.\n"+
						"  default   = %s\n  projected = %s",
						len(proj), len(def), def, proj)
				}
			})

			t.Run("text", func(t *testing.T) {
				// REQ-39: "In text mode the measurand is the total rendered
				// byte count of the emitted lines; a text projection that
				// drops no bytes fails this clause even though the S5
				// subset assertion would still pass, a subset being
				// satisfied by the equal set."
				def := requireSuccess(t, append(slices.Clone(shape.args), "--as=text")...)
				proj := requireSuccess(t,
					append(slices.Clone(shape.args), "--"+planOnlyFlag, "--as=text")...)

				if textWidth(proj) >= textWidth(def) {
					t.Errorf("projected text is %d B and default %d B; C1's "+
						"width clause binds BOTH modes unconditionally — "+
						"`0023:A8` measured 43.7-54.4%% on the same three "+
						"shapes. A text projection that drops no bytes fails "+
						"this clause even though the subset assertion would "+
						"still pass, a subset being satisfied by the EQUAL "+
						"SET.\ndefault:\n%s\nprojected:\n%s",
						textWidth(proj), textWidth(def), def, proj)
				}
			})
		})
	}
}

// REQ-39: "In text mode the measurand is the total rendered byte count of
// the emitted lines; a text projection that drops no bytes fails this
// clause even though the S5 subset assertion would still pass"
// REQ-57: "Additionally the projected text is STRICTLY SHORTER in total
// rendered bytes (C1's width clause binds both modes); a subset assertion
// alone is satisfied by the equal set, so a zero-saving text projection
// would otherwise pass."
// REQ-82 (S5's second control): "make the text projection a no-op (project
// nothing) — the subset assertion stays green and only the width assertion
// goes red"
// ADVERSARIAL — the two measurands are SEPARABLE and this asserts the
// separation directly: a no-op text projection satisfies subset and fails
// width, so the pair is jointly necessary.
func TestReq39And57And82_TheTextSubsetAndTextWidthAssertionsAreSeparablyNecessary(t *testing.T) {
	args := pricingCall(t)
	def := requireSuccess(t, append(slices.Clone(args), "--as=text")...)
	proj := requireSuccess(t,
		append(slices.Clone(args), "--"+planOnlyFlag, "--as=text")...)

	defLines, projLines := textLines(def), textLines(proj)

	// The no-op control, stated as an assertion rather than as prose: if
	// the projected line SET equals the default line set, the subset
	// assertion is satisfied by the equal set and width is the only clause
	// left standing. That is precisely the state C1 forbids.
	equalSet := len(projLines) == len(defLines)
	if equalSet {
		for _, line := range projLines {
			if !slices.Contains(defLines, line) {
				equalSet = false
				break
			}
		}
	}
	if equalSet {
		t.Errorf("the projected text line set EQUALS the default's (%d "+
			"lines); the subset assertion is satisfied by the equal set, so "+
			"only the width clause discriminates a no-op text projection — "+
			"and a no-op is what this is.\n%s", len(projLines), proj)
	}
	if textWidth(proj) >= textWidth(def) {
		t.Errorf("the projected text saves no bytes (%d vs %d)",
			textWidth(proj), textWidth(def))
	}

	// And the reduction can come ONLY from dropped lines (REQ-58): every
	// projected line is a default line byte-for-byte, so no carried value
	// was shortened to manufacture the saving.
	for _, line := range projLines {
		if !slices.Contains(defLines, line) {
			t.Errorf("projected text line %q appears in no default line; "+
				"the two measurands are JOINTLY SUFFICIENT only because the "+
				"subset is asserted PER LINE and BYTE-IDENTICALLY — "+
				"otherwise a projection that kept every line and shortened "+
				"one carried value would satisfy both loosely", line)
		}
	}
}

// REQ-55: "The projected text lines MUST be a subset of the default-mode
// text lines, byte-identical per line."
// REQ-56: "projected text lines a byte-identical, stable subset of the
// default-mode lines — set membership, not a subsequence: the wire keeps
// struct declaration order while text is sorted by the shipped `flatten`,
// and no clause requires the two orders to agree."
// REQ-58: "every projected line must appear in the default set unchanged,
// so the width reduction can only come from DROPPED lines, never from a
// shortened carried value."
// REQ-59 / MVV step 3: "the projected lines are a byte-identical subset of
// step 1's text lines, stable across repeated runs."
// REQ-110: "Witness: 15 default lines → 9 projected on the normative
// fixture"; the desk trace names the six dropped lines: "`model`,
// `observed.region`, `observed.tier`, `outcome`, `owned`, `readers`".
// REQ-82 (S5's first control): "drop a plan-group line under the flag; S5
// must go red."
// HAPPY PATH — the text-subset oracle on the normative fixture, with the
// dropped set named rather than counted.
func TestReq55And56And58And59And110_ProjectedTextLinesAreAByteIdenticalStableSubset(t *testing.T) {
	args := pricingCall(t)
	def := requireSuccess(t, append(slices.Clone(args), "--as=text")...)

	// Stability across repeated runs (REQ-59): the flattener sorts keys
	// from the decoded map, so no map-iteration-order path exists — but
	// that is a claim, and this asserts it.
	var proj string
	for i := range 3 {
		out := requireSuccess(t,
			append(slices.Clone(args), "--"+planOnlyFlag, "--as=text")...)
		if i == 0 {
			proj = out
			continue
		}
		if out != proj {
			t.Fatalf("projected text run %d differs from run 0; the "+
				"rendering is deterministic — `respond/text.go::flatten` "+
				"sorts keys from the decoded map\nfirst:\n%s\nagain:\n%s",
				i, proj, out)
		}
	}

	defLines, projLines := textLines(def), textLines(proj)

	// SET MEMBERSHIP, byte-identical per line — not a subsequence. The two
	// orderings (wire declaration order, text sorted order) are different
	// and no clause requires them to agree.
	for _, line := range projLines {
		if !slices.Contains(defLines, line) {
			t.Errorf("projected text line %q is not a default line "+
				"byte-for-byte; the subset is asserted PER LINE, so the "+
				"width reduction can only come from DROPPED lines, never "+
				"from a shortened carried value\ndefault:\n%s\nprojected:\n%s",
				line, def, proj)
		}
	}

	// The dropped lines are exactly the echo group's, NAMED — the desk
	// trace's witness. Counting alone would pass on a build that dropped a
	// plan line and kept an echo one.
	var dropped []string
	for _, line := range defLines {
		if !slices.Contains(projLines, line) {
			dropped = append(dropped, line)
		}
	}
	for _, line := range dropped {
		key := line
		if i := strings.IndexByte(line, ':'); i >= 0 {
			key = line[:i]
		}
		root := key
		if i := strings.IndexByte(key, '.'); i >= 0 {
			root = key[:i]
		}
		if !slices.Contains(echoGroup0023, root) {
			t.Errorf("the flag dropped the text line %q, whose root key %q "+
				"is NOT in the echo group %v; dropping a plan-group line is "+
				"the S5 control that must go red", line, root, echoGroup0023)
		}
	}
	if len(dropped) == 0 {
		t.Errorf("the flag dropped no text line at all; the desk trace's "+
			"witness on this fixture is 15 default lines → 9 projected, the "+
			"six dropped being `model`, `observed.region`, `observed.tier`, "+
			"`outcome`, `owned`, `readers`\ndefault:\n%s", def)
	}
}

// REQ-54: "Text mode renders the projected result through the same generic
// payload renderer as every other success, so the two modes cannot disagree
// (0005:C1 held by construction)." — [0005-carried]
// REQ-94: "Both-modes rendering | `internal/cli/respond` (`OK`,
// `text.go::flatten`) … Reuse unchanged | Projection lands before
// `respond.OK`, so no gateway change" — negative REQ: no edit to
// `internal/cli/respond`.
// DOMAIN EDGE — both-modes agreement asserted as CONTENT equality: the text
// rendering of the projected run must carry exactly the leaves the
// projected JSON payload carries, which is what "derived from the same
// verb-specific result" means and what a per-mode projection would break.
func TestReq54And94_BothModesRenderTheSameProjectedResultThroughTheOneGateway(t *testing.T) {
	args := pricingCall(t, "--"+planOnlyFlag)

	jsonOut := requireSuccess(t, append(slices.Clone(args), "--as=json")...)
	textOut := requireSuccess(t, append(slices.Clone(args), "--as=text")...)

	projFields := rawFieldsOf(t, jsonOut)
	lines := textLines(textOut)

	// Every top-level key the projected JSON carries must have at least one
	// text line rooted at it, and no text line may be rooted at a key the
	// projected JSON omits. A gateway-level or per-mode projection makes
	// exactly this disagree.
	for key := range projFields {
		var found bool
		for _, line := range lines {
			if line == key || strings.HasPrefix(line, key+":") ||
				strings.HasPrefix(line, key+".") ||
				strings.HasPrefix(line, key+"[") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("the projected JSON carries %q but no text line is "+
				"rooted at it; text is DERIVED FROM THE SAME verb-specific "+
				"result through one `flatten`, so the two modes cannot "+
				"disagree about what the run reported\ntext:\n%s",
				key, textOut)
		}
	}
	for _, line := range lines {
		root := line
		if i := strings.IndexAny(line, ":.["); i >= 0 {
			root = line[:i]
		}
		root = strings.TrimSpace(root)
		if root == "" {
			continue
		}
		if _, held := projFields[root]; !held {
			t.Errorf("the projected TEXT carries a line rooted at %q, which "+
				"the projected JSON payload does not carry; a mode-dependent "+
				"omission is structurally impossible through the one "+
				"`flatten` gateway, so this means the projection is being "+
				"applied per output mode rather than BEFORE `respond.OK`\n"+
				"json keys = %v", root, sortedKeys(projFields))
		}
	}
}

// --- projection site and flag-blindness (`0023:C1`, `0023:C2`) ----------

// REQ-47: "The projection SITE is normative … applied to the verb-specific
// result BEFORE respond.OK, never in the respond gateway and never per
// output mode"
// REQ-48: "it sits on the SUCCESS path only, AFTER the last respond.Fail
// return"
// REQ-49: "the flag is read at exactly ONE lexical site, that projection
// branch — never passed as a parameter into payload assembly, and never
// consulted in gate evaluation, rule selection, or refusal construction."
// REQ-50: "A defensive \"if refusing, skip projection\" guard is FORBIDDEN:
// refusal flag-blindness (A6) must hold structurally, because a refusal
// returns before the projection is reachable"
// REQ-53: "Nothing upstream of payload assembly reads the flag: readers,
// kernel call, gate run, and every refusal path are flag-blind. The respond
// gateway is untouched"
// ADVERSARIAL — the site is asserted BEHAVIOURALLY, through what a
// flag-aware upstream would change: the selected rule, the gate results,
// and every refusal disposition are identical ± the flag across every
// decision shape the verb has.
func TestReq47And48And49And50And53_NothingUpstreamOfPayloadAssemblyReadsTheFlag(t *testing.T) {
	gateModel := writeFlowModel(t, dtGateModel0010)
	escapeModel := writeFlowModel(t, dtEscapeModel0010)
	mvv := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, mvv, "status=draft", "stale=x")
	bind := artifactBinding(flowStateRole, art)

	for _, decision := range []struct {
		name string
		args []string
		// key is the plan-group field whose equality across widths is the
		// discriminating witness for THIS decision path.
		key string
	}{
		// Rule selection: a flag consulted in selection would move `rule`.
		{"selection", resolveArgs(pricingModelPath(t), "", "decide",
			"--tag", "tier=paid", "--tag", "region=eu"), "rule"},
		// Gate evaluation: a flag consulted in the gate run would move
		// `gates`, whose results ride the payload.
		{"gates", resolveArgs(gateModel, "", "decide", "--tag", "a=x"), "gates"},
		// The escape phase: a flag consulted there would move `escaped`.
		{"escape", resolveArgs(escapeModel, "", "decide",
			"--tag", "a=y", "--tag", "b=q"), "escaped"},
		// The write/clear plan: reached through the reader pass, the
		// upstream stage furthest from the projection site.
		{"writes-and-clear", resolveArgs(mvv, bind, "advance"), "clear"},
	} {
		t.Run(decision.name, func(t *testing.T) {
			def := rawFieldsOf(t,
				requireSuccess(t, append(slices.Clone(decision.args), "--as=json")...))
			proj := rawFieldsOf(t, requireSuccess(t,
				append(slices.Clone(decision.args), "--"+planOnlyFlag, "--as=json")...))

			defRaw, held := def[decision.key]
			if !held {
				t.Fatalf("the default payload carries no %q; this arm exists "+
					"to witness that decision path and cannot without it",
					decision.key)
			}
			projRaw, held := proj[decision.key]
			if !held {
				t.Fatalf("the projected payload carries no %q; it is a "+
					"PLAN-group field and rides both widths", decision.key)
			}
			if string(defRaw) != string(projRaw) {
				t.Errorf("%q = %s by default and %s under the flag; the same "+
					"invocation ± the flag selects the same rule, runs the "+
					"same gates, and refuses identically — ONLY report width "+
					"differs. The flag is read at exactly ONE lexical site, "+
					"the projection branch, and is never consulted in gate "+
					"evaluation, rule selection, or refusal construction",
					decision.key, defRaw, projRaw)
			}
		})
	}
}

// REQ-51: "The projection function is likewise SEPARATE from payload
// assembly: it takes the fully assembled payload and returns the projected
// one, and the flag is never a parameter to the assembly function."
// REQ-52: "Fusing the two — assembling conditionally under the flag —
// destroys both checkable properties this contract rests on … Assembly is
// flag-blind; projection is the only flag-aware step."
// REQ-67: "the projection is implemented FROM that declaration as a
// DISTINCT function taking the assembled payload and returning the
// projected one"
// ADVERSARIAL — separateness has an observable consequence, and this
// asserts it rather than reading the source: a projection applied to the
// FULLY ASSEMBLED payload cannot produce a value the default assembly never
// produced, so every carried field must be reachable in default mode with
// identical bytes on EVERY shape, including ones where conditional assembly
// would take a cheaper path.
func TestReq51And52And67_ProjectionTakesTheFullyAssembledPayloadAndOnlyDeletes(t *testing.T) {
	for _, shape := range projectionShapes0023(t) {
		t.Run(shape.name, func(t *testing.T) {
			def := rawFieldsOf(t,
				requireSuccess(t, append(slices.Clone(shape.args), "--as=json")...))
			proj := rawFieldsOf(t, requireSuccess(t,
				append(slices.Clone(shape.args), "--"+planOnlyFlag, "--as=json")...))

			// Conditional assembly shows up as a carried field that differs
			// from — or is missing relative to — the assembled default. A
			// projection that takes the fully assembled payload and DELETES
			// cannot do either.
			for _, key := range planGroup0023 {
				defRaw, defHeld := def[key]
				projRaw, projHeld := proj[key]
				if defHeld != projHeld {
					t.Errorf("plan-group field %q is %s by default and %s "+
						"under the flag; assembly is FLAG-BLIND and "+
						"projection only DELETES the echo group — a "+
						"plan-group field whose presence tracks the flag "+
						"means the two were fused",
						key, presence(defHeld), presence(projHeld))
					continue
				}
				if defHeld && string(defRaw) != string(projRaw) {
					t.Errorf("plan-group field %q = %s by default and %s "+
						"under the flag; the projection takes the FULLY "+
						"ASSEMBLED payload and returns it minus the echo "+
						"keys — it cannot compute a different value",
						key, defRaw, projRaw)
				}
			}
		})
	}
}

// --- shared shapes -------------------------------------------------------

// projectionShape0023 is one request the differential runs over.
type projectionShape0023 struct {
	name string
	args []string
}

// projectionShapes0023 returns the request shapes the differential oracles
// sweep. They are chosen to span the DIFFERENT FACT/GATE RATIOS `0023:A1`
// measured — a single shape would let a projection that happens to work on
// a stateless table pass while breaking on a reader-backed one.
//
//   - pricing 2×2: the NORMATIVE fixture; stateless, no accessors, so its
//     echo containers are EMPTY — the arm `0023:A9` names;
//   - MVV model: reader-backed, so `owned` and `readers` are POPULATED and
//     the reader pass actually runs;
//   - escaped plan: `escape_class` present, so the conditional ninth key
//     rides the projected width.
func projectionShapes0023(t *testing.T) []projectionShape0023 {
	t.Helper()

	mvv := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, mvv, "status=draft", "stale=x", "labels=[\"a<b\",\"x&y\"]")
	escapeModel := writeFlowModel(t, dtEscapeModel0010)

	return []projectionShape0023{
		{
			name: "pricing-2x2/empty-containers",
			args: pricingCall(t),
		},
		{
			name: "mvv/populated-owned-and-readers",
			args: resolveArgs(mvv, artifactBinding(flowStateRole, art),
				"advance", "--tag", "profile=mid"),
		},
		{
			name: "escaped/class-carried",
			args: resolveArgs(escapeModel, "", "decide",
				"--tag", "a=y", "--tag", "b=q"),
		},
	}
}
