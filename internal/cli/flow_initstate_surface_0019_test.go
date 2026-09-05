package cli

// RDR 0019 — the verb surface: `0019:C1`'s CARRIER paragraph, the
// `0005:C1` verb-enumeration override the Metadata field carries, and the
// non-goals the record records so a later phase does not manufacture work.
//
// The discriminating choice throughout: the verb set is read off the
// REGISTERED cobra tree and asserted by SET EQUALITY, never by "contains
// init-state". An additive override that also withdrew a verb would pass a
// containment test and fail this one.

import (
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/spf13/cobra"
)

// initStateCommand returns the registered `flow init-state` cobra command.
// A missing registration is a FATAL setup failure, not a skip: every
// surface assertion below is about a command that must exist.
func initStateCommand(t *testing.T) *cobra.Command {
	t.Helper()

	for _, c := range NewRootCmd().Commands() {
		if c.Name() != "flow" {
			continue
		}
		for _, sub := range c.Commands() {
			if sub.Name() == initStateVerb {
				return sub
			}
		}
	}
	t.Fatalf("the `flow` group registers no `%s` verb", initStateVerb)
	return nil
}

// initStateVerb is the exact spelling `0019:D-naming` fixes.
const initStateVerb = "init-state"

// flowVerbs0019 is the verb set AFTER this record's override: the four
// `flow_surface_0005_test.go` fixes, plus exactly this spelling.
var flowVerbs0019 = []string{
	"next", "resolve", "read-state", "set-state", initStateVerb,
}

// REQ-8: "The `flow` group gains one verb, `init-state` — an override
// extending 0005:C1's verb enumeration and 0005:D-naming's verb list by
// exactly this spelling."
// REQ-13: "`init-state`, completing the `read-state`/`set-state` family."
// — the verb spelling is normative; `init`, `seed`, and a `--from-initial`
// flag on `set-state` are rejected.
// REQ-14: "`internal/cli/flow_surface_0005_test.go` fixes `flowVerbs =
// {next, resolve, read-state, set-state}` and asserts SET EQUALITY against
// the registered group, so it must gain `init-state` here"
// HAPPY PATH — set equality, so a withdrawal or a second new verb fails
// here rather than passing a containment check.
func TestReq8And13And14_0019_FlowGroupGainsExactlyInitState(t *testing.T) {
	got := flowGroupNames(t)
	if got == nil {
		t.Fatalf("no `flow` command group is registered")
	}

	want := slices.Clone(flowVerbs0019)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Errorf("flow group verbs = %v; want exactly %v — this record's "+
			"override extends the enumeration by EXACTLY the `init-state` "+
			"spelling and withdraws nothing", got, want)
	}
}

// REQ-9: "**appended, not conditioned**, and additive on its owner's
// grammar: no existing verb changes and nothing is withdrawn."
// BOUNDARY — the four shipped verbs are asserted STILL PRESENT
// independently of the set-equality test above, so a regression that
// swapped one for `init-state` names the withdrawn verb.
func TestReq9_0019_NoShippedVerbIsWithdrawn(t *testing.T) {
	// The override is ADDITIVE, so the new verb must be present ALONGSIDE
	// the four. Asserting only the four would go green against a tree that
	// never gained `init-state` at all.
	if !flowSubcommand(t, initStateVerb) {
		t.Fatalf("the `flow` group registers no `%s`; the override is "+
			"APPENDED, so the new verb sits beside the four shipped ones",
			initStateVerb)
	}
	for _, verb := range []string{"next", "resolve", "read-state", "set-state"} {
		t.Run(verb, func(t *testing.T) {
			if !flowSubcommand(t, verb) {
				t.Errorf("shipped verb %q is no longer registered — this "+
					"record's override is ADDITIVE: nothing is withdrawn", verb)
			}
		})
	}
}

// REQ-9: "The new verb inherits C1's per-verb MUSTs (`respond.ValidateMode`
// first, `respond.OK`/`Fail`, `SilenceErrors`/`SilenceUsage`) unchanged"
// INPUT EDGE — `ValidateMode` runs FIRST, so an invalid `--as` mode refuses
// through the structured gateway even when every other input is also
// wrong. Driving it with a nonexistent model proves the ORDER: if
// ValidateMode did not run first, the model refusal would surface instead.
func TestReq9_0019_ValidateModeRunsFirstOnInitState(t *testing.T) {
	_, _, err := runCmd(t, "flow", initStateVerb,
		"--model", "/nonexistent/model.toml",
		"--artifact", artifactBinding(initRoleA, "/nonexistent/art.json"),
		"--as=not-a-mode")
	if err == nil {
		t.Fatalf("an invalid --as mode succeeded; want a structured refusal")
	}
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("refusal is not a structured CLIError: %v", err)
	}
	requireVerbRegistered(t, ce)
	if !strings.Contains(strings.ToLower(ce.Message+ce.Code), "mode") {
		t.Errorf("refusal = %s / %s; want the OUTPUT MODE refusal — "+
			"`respond.ValidateMode` runs FIRST, before model selection, so "+
			"the mode fault wins over the missing model", ce.Code, ce.Message)
	}
}

// REQ-9: "`SilenceErrors`/`SilenceUsage`" — a refusal must not print
// cobra's usage block.
// ADVERSARIAL — a verb registered without these flags emits its usage text
// on every refusal, which is exactly the unstructured output the CLI
// contract forbids. Asserted on the captured streams, not on a field.
func TestReq9_0019_RefusalPrintsNoUsageBlock(t *testing.T) {
	model := writeFlowModel(t, initDecisionTableModel)

	stdout, stderr, err := runCmd(t, initStateArgs(model)...)
	if err == nil {
		t.Fatalf("a decision-table model succeeded; want a class refusal")
	}
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("refusal is not a structured CLIError: %v", err)
	}
	requireVerbRegistered(t, ce)

	for name, stream := range map[string]string{"stdout": stdout, "stderr": stderr} {
		if strings.Contains(stream, "Usage:") {
			t.Errorf("%s carries cobra's usage block on a refusal — the verb "+
				"inherits `SilenceErrors`/`SilenceUsage` unchanged:\n%s",
				name, stream)
		}
	}
}

// REQ-10: "It takes the shared selection flags (`--flow`|`--model`),
// explicit `--artifact role=path` bindings, and no write grammar: its
// planned writes are the model's `[initial]` assignments."
// INPUT EDGE — the write grammar's ABSENCE is the claim, so each flag
// `set-state` owns is offered and must be REJECTED as unknown. A verb that
// quietly accepted `--write` would let a caller author a seed, which is
// exactly the plan source this verb exists to replace.
func TestReq10_0019_InitStateCarriesNoWriteGrammar(t *testing.T) {
	model := writeFlowModel(t, initMVVModel)
	art := newFlowArtifact(t, "state.artifact")

	for _, flag := range []struct{ name, argv string }{
		{"write", "--write"},
		{"clear", "--clear"},
		{"plan", "--plan"},
	} {
		t.Run(flag.name, func(t *testing.T) {
			args := initStateArgs(model, artifactBinding(initRoleA, art))
			args = append(args, flag.argv, "stage=seeded")

			ce := requireInitRefusedOnItsOwnTerms(t,
				"`init-state "+flag.argv+"`", args...)
			if !strings.Contains(ce.Message, "unknown flag") {
				t.Errorf("`init-state %s` refused with %s: %s, which is not "+
					"the UNKNOWN-FLAG refusal; the verb carries NO write "+
					"grammar — its planned writes are the model's `[initial]` "+
					"assignments", flag.argv, ce.Code, ce.Message)
			}
		})
	}
}

// REQ-10: "It takes the shared selection flags (`--flow`|`--model`)"
// HAPPY PATH — both selection flags are registered on the verb. Asserted on
// the command's own flag set, since a missing selection flag makes every
// other test in this suite unreachable.
func TestReq10_0019_InitStateTakesTheSharedSelectionFlags(t *testing.T) {
	cmd := initStateCommand(t)

	for _, name := range []string{"flow", "model", "artifact"} {
		if cmd.Flags().Lookup(name) == nil {
			t.Errorf("`init-state` does not register --%s; the verb takes "+
				"the shared selection flags and explicit `--artifact "+
				"role=path` bindings", name)
		}
	}
}

// REQ-11: "It also carries `--allow-commands`, not by declaring it but by
// INHERITING it: the flag is registered once on the `flow` group's
// persistent set (`internal/cli/flow.go::newFlowCmd`), so every child verb
// takes it structurally and this one cannot decline it."
// REQ-11: "The enumeration above is the verb's OWN grammar, not the closed
// set of flags it accepts"
// BOUNDARY — the flag must resolve on the verb (inheritance) AND must NOT
// be declared on the verb's own flag set (which would be declaring it).
// Testing only "the flag works" would pass against a local redeclaration,
// which is the reading C1 rules out.
func TestReq11_0019_AllowCommandsIsInheritedNotDeclared(t *testing.T) {
	cmd := initStateCommand(t)

	if cmd.Flags().Lookup("allow-commands") == nil {
		t.Errorf("--allow-commands does not resolve on `init-state`; the flag " +
			"is registered on the `flow` group's PERSISTENT set, so every " +
			"child verb takes it structurally and this one cannot decline it")
	}
	if cmd.LocalNonPersistentFlags().Lookup("allow-commands") != nil {
		t.Errorf("`init-state` DECLARES --allow-commands on its own flag set; " +
			"the contract fixes that it INHERITS the group's persistent flag, " +
			"never declares one")
	}
}

// REQ-11: the verb's own grammar "is not the closed set of flags it
// accepts" — passing the inherited flag is admitted, not refused as
// unknown.
// INPUT EDGE — a reader who took the CARRIER enumeration as CLOSED would
// build a verb that refuses `--allow-commands` as an unknown flag. This
// asserts the flag PARSES; the carrier gate's preemption of its refusal is
// asserted separately (REQ-36).
func TestReq11_0019_AllowCommandsParsesOnInitState(t *testing.T) {
	model := writeFlowModel(t, initMVVModel)
	art := newFlowArtifact(t, "state.artifact")

	args := append(initStateArgs(model, artifactBinding(initRoleA, art)),
		"--allow-commands")
	_, _, err := runCmd(t, args...)
	if err != nil && strings.Contains(err.Error(), "unknown flag") {
		t.Errorf("`init-state --allow-commands` refused as an UNKNOWN FLAG: "+
			"%v — the verb's own grammar is not the closed set of flags it "+
			"accepts; the group's persistent flag is inherited", err)
	}
}

// REQ-15: "Six code/doc sites hard-code the cardinal and change with it:
// `internal/cli/flow.go`'s package comment, `::newFlowCmd`'s `Long` body,
// `::flowExtendedDesc`, `internal/cli/flow_exec.go`'s header comment,
// `docs/cli-reference.md` (two strings), and `docs/model-authoring.md`."
// ADVERSARIAL — the cardinal is a claim about the verb count that goes
// STALE silently. Asserted as the absence of the OLD cardinal ("four
// verbs") anywhere the record names, which is what a missed site leaves
// behind.
func TestReq15_0019_NoSiteStillClaimsFourVerbs(t *testing.T) {
	sites := []string{
		"internal/cli/flow.go",
		"internal/cli/flow_exec.go",
		"docs/cli-reference.md",
		"docs/model-authoring.md",
	}
	stale := []string{"four verbs", "four flow verbs", "the four verbs"}

	for _, site := range sites {
		t.Run(site, func(t *testing.T) {
			body := readRepoFile(t, site)
			lower := strings.ToLower(body)
			for _, phrase := range stale {
				if strings.Contains(lower, phrase) {
					t.Errorf("%s still hard-codes the OLD cardinal %q; the "+
						"group now carries %d verbs and every site that "+
						"names the count changes with it",
						site, phrase, len(flowVerbs0019))
				}
			}
		})
	}
}

// REQ-15: the same six sites must NAME the new verb where they enumerate
// the group's verbs.
// BOUNDARY — the docs sites are the operator-facing enumeration; a verb
// absent from them is unreachable in practice even when registered.
func TestReq15_0019_DocSitesEnumerateInitState(t *testing.T) {
	for _, site := range []string{
		"docs/cli-reference.md",
		"docs/model-authoring.md",
	} {
		t.Run(site, func(t *testing.T) {
			if !strings.Contains(readRepoFile(t, site), initStateVerb) {
				t.Errorf("%s does not name `%s`; it is one of the sites the "+
					"record fixes as changing with the verb set",
					site, initStateVerb)
			}
		})
	}
}

// REQ-16: "Regenerate and commit `llms.txt` (`internal/cli/docs.go`); it
// enumerates the verbs as bullets and hard-codes no cardinal, as does
// `docs/cli-output-contract.md`, so neither needs a cardinal edit"
// HAPPY PATH — the obligation is that the new verb APPEARS in the
// generated enumeration, not that a cardinal changed. A stale `llms.txt` is
// a committed artifact that silently omits the verb.
func TestReq16_0019_LLMsTxtEnumeratesInitState(t *testing.T) {
	for _, site := range []string{"llms.txt", "docs/cli-output-contract.md"} {
		t.Run(site, func(t *testing.T) {
			body := readRepoFile(t, site)
			if !strings.Contains(body, initStateVerb) {
				t.Errorf("%s does not enumerate `%s` — it lists the verbs as "+
					"bullets and must gain this one", site, initStateVerb)
			}
			if strings.Contains(strings.ToLower(body), "four verbs") {
				t.Errorf("%s hard-codes a cardinal; the record fixes that it "+
					"does NOT and needs no cardinal edit", site)
			}
		})
	}
}

// REQ-95: "`docs/cli-output-contract.md` (verb I/O, refusal codes),
// `docs/cli-reference.md`, `docs/model-authoring.md` (the start-state
// sentence gains its runtime carrier), `--help-all` text."
// BOUNDARY — `--help-all` is the surface a caller reaches for the shared
// refusal families, so the new verb must appear there too.
func TestReq95_0019_HelpAllNamesInitState(t *testing.T) {
	stdout := requireSuccess(t, "flow", "--help-all")
	if !strings.Contains(stdout, initStateVerb) {
		t.Errorf("`flow --help-all` does not name `%s`:\n%s",
			initStateVerb, stdout)
	}
}

// REQ-96: "it states the FILE-BACKED scope (C1 carrier scope), not an
// unqualified \"every state-machine flow\"; and it names the verb from the
// surface where an operator meets the wall — `flow next`'s
// `unknown[].reason: absent` report"
// DOMAIN EDGE — the doc obligation is a SCOPE qualification. A doc that
// promised the verb for "every state-machine flow" would be wrong for
// edit-carried and command-backed models, which the carrier gate refuses.
func TestReq96_0019_ModelAuthoringStatesTheFileBackedScope(t *testing.T) {
	body := strings.ToLower(readRepoFile(t, "docs/model-authoring.md"))

	if !strings.Contains(body, initStateVerb) {
		t.Fatalf("docs/model-authoring.md does not name `%s` at all; the "+
			"start-state sentence gains its runtime carrier here, and the "+
			"scope qualification below has no subject until it does",
			initStateVerb)
	}
	if !strings.Contains(body, "file-backed") {
		t.Errorf("docs/model-authoring.md names `%s` without stating the "+
			"FILE-BACKED scope; the verb refuses edit-carried and "+
			"command-backed carriers, so an unqualified promise is wrong",
			initStateVerb)
	}
	if strings.Contains(body, "every state-machine flow") {
		t.Errorf("docs/model-authoring.md makes the UNQUALIFIED claim the " +
			"record forbids (\"every state-machine flow\")")
	}
}

// REQ-97: "Pointing an operator from `flow next` to `init-state` in-band
// would mean adding a field or a reason-string to another verb's payload,
// which is an override of `0005:C1`'s I/O for `next` that this record does
// not carry — its override is the verb enumeration, nothing more."
// ADVERSARIAL — the tempting change is exactly the one forbidden. Asserted
// on `next`'s payload against a store with an absent `[initial]` key: no
// new top-level field, and no reason string other than the shipped ones.
func TestReq97_0019_NextPayloadGainsNoInitStateField(t *testing.T) {
	// The non-goal only has content once the verb exists: a tree with no
	// `init-state` trivially mentions it nowhere.
	if !flowSubcommand(t, initStateVerb) {
		t.Fatalf("the `flow` group registers no `%s`; this non-goal has no "+
			"subject until it does", initStateVerb)
	}
	model := writeFlowModel(t, initMVVModel)
	art := newFlowArtifact(t, "state.artifact")

	stdout := requireSuccess(t, "flow", "next", "--model", model,
		"--artifact", artifactBinding(initRoleA, art),
		"--as=json")

	for _, s := range payloadValueStrings(t, stdout) {
		if strings.Contains(s, initStateVerb) {
			t.Errorf("`flow next`'s payload carries the string %q; this "+
				"record's override is the VERB ENUMERATION, nothing more — "+
				"it adds no field or reason-string to another verb's payload",
				s)
		}
	}

	data := flowData(t, stdout)
	for _, forbidden := range []string{"init", "init_state", "initState", "remedy"} {
		if _, present := data[forbidden]; present {
			t.Errorf("`flow next`'s payload gained a top-level %q field — "+
				"this record carries no `0005:C1` I/O override for `next`",
				forbidden)
		}
	}
}

// REQ-99: "no existing verb, flag, artifact shape or model schema changes."
// REQ-100: "no lint arm owed; a model that loads already satisfies it" —
// the verb owes no writer-arity lint arm.
// ADVERSARIAL — the non-goal is that `flow lint` gains NOTHING. A model
// that loads and lints clean today must still lint clean, and lint must not
// start reporting an un-seeded store.
func TestReq99And100_0019_LintGainsNoArm(t *testing.T) {
	if !flowSubcommand(t, initStateVerb) {
		t.Fatalf("the `flow` group registers no `%s`; \"lint gains no arm\" "+
			"is vacuous until the verb this record adds exists", initStateVerb)
	}
	model := writeFlowModel(t, initMVVModel)

	stdout := requireSuccess(t, "lint", "--model", model, "--as=json")

	for _, s := range payloadValueStrings(t, stdout) {
		lower := strings.ToLower(s)
		if strings.Contains(lower, initStateVerb) ||
			strings.Contains(lower, "un-seeded") ||
			strings.Contains(lower, "unseeded") {
			t.Errorf("`flow lint` reports %q — this record owes NO lint arm: "+
				"a model that loads already satisfies the writer-arity "+
				"requirement, and `[initial]`'s lint role (0006:C18) is "+
				"unchanged", s)
		}
	}
}

// REQ-99: "no existing verb, flag, artifact shape or model schema changes."
// REQ-13: "`init-state`, completing the `read-state`/`set-state` family." —
// the verb spelling is normative; `init`, `seed`, and a `--from-initial`
// flag on `set-state` are rejected.
// BOUNDARY — `set-state` in particular must not gain the `--from-initial`
// flag `0019:D-naming` names as REJECTED.
func TestReq99And13_0019_SetStateGainsNoFromInitialFlag(t *testing.T) {
	model := writeFlowModel(t, initMVVModel)
	art := newFlowArtifact(t, "state.artifact")

	if !flowSubcommand(t, initStateVerb) {
		t.Fatalf("the `flow` group registers no `%s`; `0019:D-naming` "+
			"rejects `--from-initial` IN FAVOUR OF this verb, so the "+
			"rejection has no subject until the verb exists", initStateVerb)
	}

	_, _, err := runCmd(t, "flow", "set-state", "--model", model,
		"--artifact", artifactBinding(initRoleA, art),
		"--from-initial", "--as=json")
	if err == nil {
		t.Errorf("`set-state --from-initial` was ACCEPTED; `0019:D-naming` " +
			"REJECTS a `--from-initial` flag on `set-state` in favour of the " +
			"`init-state` verb")
	}
}

// REQ-1: "`[initial]` is BOTH the lint-time reachability root (0006:C18
// unchanged — a state-machine model declaring none stays a blocking
// finding) AND the runtime bootstrap source for owned state."
// BOUNDARY — the word is BOTH, so the test asserts the two roles TOGETHER
// over one pair of models. A test of either role alone would pass against
// an implementation that traded one for the other — which is exactly the
// change "0006:C18 unchanged" forbids.
func TestReq1_0019_InitialIsBothTheLintRootAndTheRuntimeBootstrapSource(t *testing.T) {
	// Role 1, unchanged: a state-machine model declaring NO `[initial]`
	// stays a BLOCKING lint finding.
	noRoot := writeFlowModel(t, initNoInitialModel)
	stdout, _, err := runCmd(t, "lint", "--model", noRoot, "--as=json")
	if err == nil {
		data := flowData(t, stdout)
		findings, _ := data["findings"].([]any)
		if len(findings) == 0 {
			t.Errorf("`lint` reports NO finding for a state-machine model " +
				"declaring no `[initial]`; 0006:C18 is UNCHANGED here — such " +
				"a model stays a blocking finding, and this record adds a " +
				"runtime role WITHOUT withdrawing the lint one")
		}
	}

	// Role 2, new: the SAME declaration is the runtime bootstrap source, so
	// a model that DOES declare `[initial]` seeds from it.
	withRoot := writeFlowModel(t, initMVVModel)
	bind := artifactBinding(initRoleA, newFlowArtifact(t, "state.artifact"))
	requireSuccess(t, initStateArgs(withRoot, bind)...)

	owned := ownedFromReadState(t, requireSuccess(t, readStateArgs(withRoot, bind)...))
	requireReadsBack(t, owned, "stage", "seeded")
	requireReadsBack(t, owned, "note", "hello")
}

// initNoInitialModel is a state-machine model declaring NO `[initial]`.
// 0006:C18 makes that a blocking lint finding, and this record leaves that
// arm unchanged.
const initNoInitialModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "initnoroot"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.stage]
provenance = "owned"
kind = "enum"
domain = ["seeded", "final"]
single_valued = true

[read.state]
role = "state"
path = "flow.state"
keys = ["stage"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["stage"]
timeout = "2s"
read_back = true

[context.done]
[context.done.match.stage]
eq = "final"

[[rule]]
id = "advance"
[rule.match.stage]
eq = "seeded"
[rule.match.recognized]
eq = "advance"
[rule.write]
stage = "final"
`

// REQ-12: "It MUST route every seed through the declared write accessors
// with commit-time read-back, on the same writer-routing/no-cross-writer-
// atomicity terms as `set-state` (0005:C1); it MUST NOT write an artifact
// directly."
// ADVERSARIAL — "MUST NOT write an artifact directly" is the claim, and a
// direct write is INVISIBLE to an oracle that only checks the resulting
// bytes. The discriminating construction is a model whose write accessor's
// commit-time READ-BACK cannot complete: routing THROUGH the accessor
// surfaces that as an exit-3 refusal, while writing the artifact directly
// would succeed at exit 0, never having consulted the read-back at all.
func TestReq12_0019_EverySeedIsRoutedThroughTheWriteAccessorWithCommitTimeReadBack(t *testing.T) {
	model := writeFlowModel(t, initSealModel)
	bind := artifactBinding(initRoleA, newFlowArtifact(t, "sealed.artifact"))

	stdout, _, err := runCmd(t, initStateArgs(model, bind)...)
	if err == nil {
		t.Fatalf("the seed SUCCEEDED against a writer whose commit-time "+
			"read-back cannot complete; the verb MUST route every seed "+
			"through the declared write accessors WITH commit-time "+
			"read-back and MUST NOT write an artifact directly — a direct "+
			"write is exactly what exits 0 here\nstdout: %s", stdout)
	}
	// And it is the READ-BACK that refused: exit 3, the environment group
	// the incomplete-read-back arm carries.
	initRefusal(t, 3, initStateArgs(model, bind)...)
}
