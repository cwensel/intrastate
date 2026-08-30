package cli

// RDR 0025 — the execution gate at the CLI surface (C6; scenario S7) and
// the single `CLIError.Detail` slot the diagnosis lands in (C4).
//
// S7's own framing, which this file follows exactly: the executor callers
// are the WRONG axis. `flowbind.Registry` has ONE production call site, so
// every executor shares one gated registry and testing them separately
// tests one path repeatedly. With C6's ONE registration on the `flow`
// group, the drift surface is CONTAINMENT — so the oracle DERIVES the set
// of verbs reaching `buildRequest` from the command tree and asserts each
// resolves the flag through the group's persistent set. It never restates
// a literal four, and it never text-matches pflag's error.
//
// PHASE-0 READING Q4: S7's containment arm is asserted AS WORDED — a
// structural `Flags().Lookup("allow-commands") != nil` walk over the
// group's subtree, with `lint` at root as the complementary negative.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/spf13/cobra"
)

// allowCommandsFlag is the flag name C6 fixes. Naming it once here keeps
// the walk below from restating a literal per assertion.
const allowCommandsFlag = "allow-commands"

// lintModelInvalid is the code root `lint` carries for a load refusal. It
// is NOT the `flow`-group-scoped `codeModelInvalid`: the two surfaces
// answer under different codes and always have (deviations D10).
const lintModelInvalid = "model-invalid"

// --- C6: registration and containment (S7) --------------------------------

// REQ-89: "registration: ONE registration, on the `flow` GROUP —
// `internal/cli/flow.go::newFlowCmd`'s `PersistentFlags()`, whose subtree is
// exactly the four `buildRequest` call sites across three verbs
// (`flow_next.go`, `flow_resolve.go`, `flow_state.go` ×2) and nothing
// else."
// REQ-131 (S7's containment arm): "the scenario derives the set of verbs
// reaching `buildRequest` and asserts each RESOLVES the flag through the
// group's persistent set — a structural `Flags().Lookup(\"allow-commands\")
// != nil` walk, never a text match on pflag's error. A fifth verb added
// INSIDE the group inherits the flag and must pass; one added OUTSIDE it
// must fail this test rather than silently inherit a lookup miss."
// BOUNDARY
//
// PHASE-0 READING Q4: asserted as worded.
func TestReq89_EveryVerbInTheFlowGroupResolvesTheAllowCommandsFlag(t *testing.T) {
	root := NewRootCmd()

	flow := subcommand(t, root, "flow")

	// The flag is registered ONCE, on the group's persistent set — not on
	// any verb's own set. A per-verb registration would pass the subtree
	// walk below while leaving a fifth verb to be remembered.
	if flow.PersistentFlags().Lookup(allowCommandsFlag) == nil {
		t.Fatalf("`flow` does not register --%s on PersistentFlags(); ONE "+
			"registration on the group is what makes the caller set structural "+
			"rather than hand-kept", allowCommandsFlag)
	}

	// The DERIVED set: every verb in the group's subtree. The oracle
	// derives it rather than restating a literal four, so a fifth verb
	// added inside the group is covered without editing this test.
	verbs := leafVerbs(flow)
	if len(verbs) == 0 {
		t.Fatal("the `flow` group has no verbs; the containment walk would be " +
			"vacuous")
	}
	for _, v := range verbs {
		// The structural resolution check, exactly as S7 words it. Never a
		// text match on pflag's "flag accessed but not defined".
		//
		// `InheritedFlags()` is cobra's own accessor for the parent-pflag
		// MERGE that `Flags()` reports lazily — it performs the merge and
		// returns; the assertion below stays on `Flags()`, as REQ-131
		// words it, and still discriminates: a verb OUTSIDE the group
		// resolves nothing after the same call (deviations D9).
		_ = v.InheritedFlags()
		if v.Flags().Lookup(allowCommandsFlag) == nil {
			t.Errorf("`flow %s` does not RESOLVE --%s; a verb inside the group "+
				"inherits the persistent flag, and one that does not would "+
				"silently lose command execution", v.Name(), allowCommandsFlag)
		}
	}

	t.Run("a verb added INSIDE the group inherits it", func(t *testing.T) {
		// The positive control S7 names: a fifth verb inside the group must
		// pass the same walk, with no edit to the registration.
		fresh := NewRootCmd()
		g := subcommand(t, fresh, "flow")
		fifth := &cobra.Command{Use: "fifth", RunE: func(*cobra.Command, []string) error { return nil }}
		g.AddCommand(fifth)

		_ = fifth.InheritedFlags()
		if fifth.Flags().Lookup(allowCommandsFlag) == nil {
			t.Errorf("a fifth verb added inside the group does not resolve --%s; "+
				"the coupling must be to the GROUP, not to a hand-kept verb list",
				allowCommandsFlag)
		}
	})

	t.Run("a verb added OUTSIDE the group does NOT inherit it", func(t *testing.T) {
		// The complementary negative. A verb outside the group must FAIL
		// this test rather than silently inherit a lookup miss — which is
		// the shape that fails OPEN once someone discards the error.
		fresh := NewRootCmd()
		outside := &cobra.Command{Use: "outside", RunE: func(*cobra.Command, []string) error { return nil }}
		fresh.AddCommand(outside)

		_ = outside.InheritedFlags()
		if outside.Flags().Lookup(allowCommandsFlag) != nil {
			t.Errorf("a verb at ROOT resolves --%s; the flag's subtree is the "+
				"`flow` group and nothing else", allowCommandsFlag)
		}
	})
}

// REQ-88: "--allow-commands        # v1's ONLY opt-in: a per-invocation
// flag; never the model file; `lint` does not carry it"
// REQ-90: "`lint` sits at ROOT, outside the group, so it does not carry the
// flag"
// BOUNDARY — the complementary arm S7 names.
func TestReq90_LintSitsAtRootAndDoesNotResolveTheFlag(t *testing.T) {
	root := NewRootCmd()

	lint := subcommand(t, root, "lint")

	_ = lint.InheritedFlags()
	if lint.Flags().Lookup(allowCommandsFlag) != nil {
		t.Errorf("`lint` resolves --%s; validation is UNGATED and a flag that "+
			"changes nothing is a false affordance — the exclusion is "+
			"structural, not remembered", allowCommandsFlag)
	}

	// And that lint is genuinely at root rather than nested under the group,
	// which is what makes the exclusion structural.
	if lint.Parent() != root {
		t.Errorf("`lint`'s parent is %q; want the root command", lint.Parent().Name())
	}

	// The discriminating half: root itself must not carry it either, or
	// "lint does not resolve it" would be satisfied by a flag nobody
	// registered at all.
	flow := subcommand(t, root, "flow")
	if flow.PersistentFlags().Lookup(allowCommandsFlag) == nil {
		t.Errorf("--%s is registered nowhere; the lint exclusion above proves "+
			"nothing until the flag exists on the group", allowCommandsFlag)
	}
}

// leafVerbs returns the runnable verbs in a command's subtree.
func leafVerbs(cmd *cobra.Command) []*cobra.Command {
	var out []*cobra.Command
	for _, c := range cmd.Commands() {
		if c.Name() == "help" || c.Name() == "completion" {
			continue
		}
		if c.Runnable() {
			out = append(out, c)
		}
		out = append(out, leafVerbs(c)...)
	}
	return out
}

func subcommand(t *testing.T, parent *cobra.Command, name string) *cobra.Command {
	t.Helper()

	for _, c := range parent.Commands() {
		if c.Name() == name {
			return c
		}
	}
	t.Fatalf("no %q subcommand under %q", name, parent.Name())
	return nil
}

// --- C6: the gate's observable behaviour (S7) -----------------------------

// REQ-92: "absent ⇒ every command invocation refuses execution_failure
// before spawn, Detail naming the gate; lint (C1/C5) validates regardless"
// REQ-97: "With the gate off, a command invocation refuses before spawn
// (`execution_failure`, `Detail` naming `allow_commands`), while `intrastate
// lint` still validates the entries — the model is reviewable before it is
// trusted."
// REQ-131 (S7): "asserted by **absence of a spawn** (a sentinel argv0 that
// would leave an observable trace if executed), not merely by a non-zero
// exit; normal execution under the flag; lint passes in both"
// ADVERSARIAL
func TestReq92_WithoutTheFlagACommandModelRefusesAndNoChildSpawns(t *testing.T) {
	dir := t.TempDir()
	trace := filepath.Join(dir, "spawned")
	model := writeCommandModel(t, dir, traceArgv(t, trace))
	art := filepath.Join(dir, "state.json")

	_, _, err := runCmd(t, "flow", "set-state",
		"--model", model,
		"--artifact", "state="+art,
		"--write", "status=final",
		"--as=json")

	if err == nil {
		t.Fatal("a command-backed set-state succeeded with no --allow-commands; " +
			"execution requires an opt-in outside the model")
	}
	if code := clierr.ErrorCode(err); code != codeAccessorFailed {
		t.Errorf("code = %q; want %q — the gate refusal is an execution failure",
			code, codeAccessorFailed)
	}
	// ABSENCE OF A SPAWN, not merely a non-zero exit.
	if _, serr := os.Stat(trace); serr == nil {
		t.Error("a child ran; the gate refuses BEFORE spawn")
	}

	t.Run("the same invocation runs WITH the flag", func(t *testing.T) {
		// The discriminating half. Without it, "it refused" is satisfied by
		// a build where command entries never work at all.
		trace := filepath.Join(t.TempDir(), "spawned")
		dir := t.TempDir()
		model := writeCommandModel(t, dir, traceArgv(t, trace))
		art := filepath.Join(dir, "state.json")

		_, _, _ = runCmd(t, "flow", "set-state",
			"--model", model,
			"--artifact", "state="+art,
			"--write", "status=final",
			"--allow-commands",
			"--as=json")

		if _, serr := os.Stat(trace); serr != nil {
			t.Errorf("no child ran with --allow-commands; the gate must be the "+
				"only thing that stopped the invocation above: %v", serr)
		}
	})
}

// REQ-92 (lint arm) / REQ-97: "lint (C1/C5) validates regardless"
// REQ-118: "Phase 4: Surface and proof — Wire lint reporting for the C5
// defects"
// HAPPY PATH
//
// The pair: the same model lints clean with and without the flag (which
// lint does not carry at all), and a C5 mutant is REPORTED by lint under
// both — so the model is reviewable before it is trusted.
func TestReq92_LintValidatesACommandModelRegardlessOfTheGate(t *testing.T) {
	dir := t.TempDir()
	model := writeCommandModel(t, dir, []string{"true"})

	t.Run("a valid command model lints clean", func(t *testing.T) {
		if _, _, err := runCmd(t, "lint", "--model", model, "--as=json"); err != nil {
			t.Errorf("a valid command-backed model failed lint: %v — validation "+
				"is ungated", err)
		}
	})

	t.Run("a C5 mutant is reported by lint", func(t *testing.T) {
		mutant := filepath.Join(t.TempDir(), "mutant.toml")
		src, rerr := os.ReadFile(model)
		if rerr != nil {
			t.Fatalf("reading the model: %v", rerr)
		}
		// `writeCommandModel` appends `{artifact}` to the supplied argv, so
		// the write entry reads `command = ["true", "{artifact}"]`
		// (deviations D10).
		bad := strings.Replace(string(src),
			`command = ["true", "{artifact}"]`,
			`command = ["sh", "-c", "true", "{artifact}"]`, 1)
		if bad == string(src) {
			t.Fatal("the mutant substitution did not apply")
		}
		if werr := os.WriteFile(mutant, []byte(bad), 0o600); werr != nil {
			t.Fatalf("writing the mutant: %v", werr)
		}

		_, _, err := runCmd(t, "lint", "--model", mutant, "--as=json")
		if err == nil {
			t.Fatal("lint accepted an `sh -c` interpreter form; C5's defects are " +
				"what make the model reviewable before it is trusted")
		}
		// `lint` is at ROOT and carries its OWN model-invalid code; the
		// `flow`-scoped constant belongs to the group (deviations D10).
		if code := clierr.ErrorCode(err); code != lintModelInvalid {
			t.Errorf("code = %q; want %q", code, lintModelInvalid)
		}
	})
}

// REQ-91: "lookup: `buildRequest` reads the flag with the error CHECKED,
// not discarded. Because persistent registration guarantees presence, a
// lookup miss can only be a wiring bug and MUST surface as one"
// ADVERSARIAL
//
// PHASE-0 READING (ASSUMPTION REQ-91): "surface as one" means the request
// build returns an ERROR through its existing error return, not a panic —
// the repo has no panic-on-wiring-bug idiom.
//
// The oracle is behavioural and structural at once: with the registration
// REMOVED from the group, a verb that reaches buildRequest must refuse
// loudly rather than read the zero value and proceed as if the user had
// declined. A default-deny resting on a lookup miss is NOT the gate — it
// would make the safety property a coincidence between two independent
// decisions.
func TestReq91_ALookupMissIsAWiringBugAndSurfacesAsOne(t *testing.T) {
	requireGit(t)

	// The write tool must genuinely APPLY: the discriminating arm below
	// asserts the intact tree SUCCEEDS on the identical invocation, and a
	// no-op `true` would refuse `read_back_mismatch` for a reason that has
	// nothing to do with the lookup (deviations D13).
	dir := t.TempDir()
	model := writeCommandModel(t, dir, []string{writeWrapper(t, dir)})
	art := filepath.Join(dir, "state.cfg")

	// The wiring bug, staged: the SAME verb is re-parented under a group
	// that never registered the flag, so pflag's GetBool at `buildRequest`
	// returns (false, *NotExistError) — a PROGRAMMER error, whose own text
	// is "flag accessed but not defined", never a statement about user
	// intent. This is exactly the shape C6 refuses to rest on.
	root := NewRootCmd()
	orphanGroup := &cobra.Command{Use: "orphan", SilenceErrors: true, SilenceUsage: true}
	orphanGroup.AddCommand(newFlowSetStateCmd())
	root.AddCommand(orphanGroup)

	if orphanGroup.PersistentFlags().Lookup(allowCommandsFlag) != nil {
		t.Fatal("the staged group registered the flag; the arm below would " +
			"prove nothing")
	}

	err := runOn(t, root, "orphan", "set-state",
		"--model", model,
		"--artifact", "state="+art,
		"--write", "status=final",
		"--as=json")

	if err == nil {
		t.Fatal("set-state proceeded with the flag unregistered; a lookup miss " +
			"is a WIRING BUG and must surface as one — reading the zero value " +
			"makes fail-closed a coincidence between two independent decisions, " +
			"and the dangerous edit (a parent registering it, or a verb " +
			"registering a different default) silently converts the miss into a " +
			"hit and fails OPEN")
	}
	if code := clierr.ErrorCode(err); code == "" {
		t.Errorf("the refusal carries no structured code: %v", err)
	}

	// The discriminating half, and the one that carries the clause. With
	// the registration INTACT the identical invocation SUCCEEDS — so the
	// orphan's refusal above is about the missing registration and nothing
	// else. Without this, "the miss refuses" is satisfied by a verb that
	// refuses every set-state, which is precisely the coincidence C6 says
	// must not stand in for the gate.
	t.Run("with the registration intact the identical invocation succeeds", func(t *testing.T) {
		fresh := NewRootCmd()
		ferr := runOn(t, fresh, "flow", "set-state",
			"--model", model,
			"--artifact", "state="+art,
			"--write", "status=final",
			"--allow-commands",
			"--as=json")
		if ferr != nil {
			t.Fatalf("the intact tree refused the identical invocation: %v — the "+
				"orphan's refusal above then proves nothing about the lookup",
				ferr)
		}
	})
}

// runOn drives an already-built root tree, so a test can mutate the tree
// before running it.
func runOn(t *testing.T, root *cobra.Command, args ...string) error {
	t.Helper()

	var out, errBuf strings.Builder
	root.SetOut(&out)
	root.SetErr(&errBuf)
	return ExecuteAndEmit(root, args)
}

// --- C4: the single CLIError.Detail slot ----------------------------------

// REQ-65: "There is exactly **one** slot to land in —
// `internal/cli/clierr/clierr.go::CLIError` has a single `Detail string` and
// `EmitText` renders a single `detail:` line — so \"renders beneath\" is not
// available and no second carrier is added here. ... they occupy that one
// slot in order: the applied-sense text **first** ..., the stderr tail
// appended after it."
// REQ-124 (S3): "the 4 KiB stderr tail in `Refusal.Detail` — asserted
// non-empty on case 3, where no applied-sense text competes for
// `CLIError.Detail`."
// DOMAIN EDGE
//
// Two arms, because the clause fixes an ORDER and an order needs both
// operands: a refusal carrying only the tail, and one carrying both.
func TestReq65_TheAppliedSenseLeadsAndTheStderrTailFollowsInOneSlot(t *testing.T) {
	dir := t.TempDir()

	t.Run("tail only, no applied sense", func(t *testing.T) {
		// A READ whose tool exits non-zero WITH STDERR: an execution
		// failure with no applied sense competing for the slot. `false`
		// writes no stderr and so has no tail to carry, which is the
		// arm's own subject — the fixture takes FX-exit-codes E1c's shape
		// instead: a non-zero exit that says something (deviations D14).
		model := writeCommandReadModel(t, dir, []string{stderrScript(t, dir)})
		art := filepath.Join(dir, "state.json")

		_, _, err := runCmd(t, "flow", "read-state",
			"--model", model,
			"--artifact", "state="+art,
			"--allow-commands",
			"--as=json")
		if err == nil {
			t.Fatal("a read whose tool exited non-zero succeeded")
		}
		ce := asCLIErr(t, err)
		if ce.Detail == "" {
			t.Error("Detail is empty on an execution failure; C4's stderr tail " +
				"lands in the CLI's single Detail slot")
		}
		if strings.Contains(ce.Detail, "may have been applied") {
			t.Errorf("Detail = %q; no write ran, so no applied sense competes "+
				"for the slot", ce.Detail)
		}
	})

	t.Run("both senses: the applied text leads", func(t *testing.T) {
		// A WRITE that applies, then a read-back whose reader fails with a
		// tail. Losing "the write may have applied" to a diagnostic tail
		// would drop the more consequential fact, so it goes FIRST.
		wdir := t.TempDir()
		model := writeAppliedThenFailingReadBackModel(t, wdir)
		art := filepath.Join(wdir, "state.json")

		_, _, err := runCmd(t, "flow", "set-state",
			"--model", model,
			"--artifact", "state="+art,
			"--write", "status=final",
			"--allow-commands",
			"--as=json")
		if err == nil {
			t.Fatal("a write whose read-back failed reported success")
		}
		ce := asCLIErr(t, err)

		applied := strings.Index(ce.Detail, "applied")
		if applied < 0 {
			t.Fatalf("Detail = %q; want the applied-but-unverified sense", ce.Detail)
		}
		// One slot, in order: the applied sense first, the tail after it.
		if len(ce.Detail) <= applied+len("applied") {
			t.Errorf("Detail = %q; want the stderr tail appended AFTER the "+
				"applied-sense text in the one slot", ce.Detail)
		}
		if strings.HasPrefix(strings.TrimSpace(ce.Detail), "wrapper") {
			t.Errorf("Detail = %q; the applied sense must LEAD — losing it to a "+
				"diagnostic tail drops the more consequential fact", ce.Detail)
		}
	})
}

func asCLIErr(t *testing.T, err error) *clierr.CLIError {
	t.Helper()

	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("the refusal is not a *clierr.CLIError: %v", err)
	}
	return ce
}
