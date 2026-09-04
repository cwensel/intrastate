package cli

// The declared-request refusal envelope, split by the refusal's ORIGIN.
//
// RDR 0028 `0028:C1.3` EXIT GROUP: routes an `execution_failure` refusal
// carrying the declared-request marker to the exit-2 REQUEST group. C1.6
// then widened where that marker is minted, and it is now minted at three
// different KINDS of site of which only ONE is a declared line edit:
//
//   - `flowbind.EditWriter` — the rule-scoped and entry-level refusals
//     C1.3 was written about. These name a rule `<id>.edit.<key>`.
//   - `cmdbind.substitute` — C1.6's argv preconditions, which run for
//     EVERY command-backed accessor: a reader, a gate, or a WRITER that
//     may carry no `edit` table at all.
//   - `accessor.Executor` — the read-back argv preconditions, which come
//     from an edit writer and are still about the READER's argv.
//
// `accessorFailureOf` originally answered the marker with one envelope,
// C1.3's own: a code naming a line edit and a message naming "the write
// accessor". Every argv precondition was therefore reported as a line-edit
// refusal, with `editRefusalFindings` coding the argv complaint as one —
// inventing a rule subject that does not exist.
//
// WHY ORIGIN, AND NOT PHASE OR CARRIER. Neither proxy holds. The invoking
// PHASE fails because a command-backed writer's argv refusal arrives at
// `phaseWrite` with no edit table in the model. The writer's CARRIER fails
// because an edit writer's read-back precondition has `Edit != nil` and is
// still not a line edit's refusal. Only the site that minted the refusal
// knows, so it says so, on the existing `Err` slot.
//
// WHAT THIS FILE PINS. The envelope names the role the accessor ACTUALLY
// ran under, and only a refusal a declared line edit minted carries the
// line-edit code and its `findings[]`.
//
// WHAT IT MUST NOT PIN. The CLI codes' SPELLINGS. C1.3 states outright
// that the code is "Stage 8's, non-normative here" (REQ-44, REQ-113), and
// deviation D7 records that no test pins it. So the assertions below
// compare the two codes to EACH OTHER — they must differ, which is the
// whole defect — and never to a literal.
//
// THE REGRESSION THIS EXISTS TO CATCH runs both ways. A fix that moved the
// non-edit arms off the line-edit envelope by dropping them back to
// `codeAccessorFailed` would regress the exit group to 3, telling a caller
// to re-run unchanged an invocation that can only ever fail again; and one
// that discriminated on phase or carrier would leave a command-backed
// writer still reported as a line edit. Both are asserted here.

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
)

// --- the fixture ----------------------------------------------------------

// requestPhaseRole is the artifact role every accessor in these models
// binds.
const requestPhaseRole = "state"

// requestPhaseModel builds a flow whose COMMAND-backed accessor at the
// named table carries a `{tag.<key>}` element naming an undeclarable-by-
// the-caller key, so substitution refuses before any child spawns.
//
// The `where` argument selects which accessor is command-backed; every
// other accessor in the model stays file-backed, so exactly one phase can
// reach the refusal and the test cannot be passing for the wrong reason.
//
// `{tag.branch}` sits at a LATER argv position, never argv0: C1.4's
// `edit_tag_argv0` refuses a placeholder there at LOAD, which would take
// the invocation out before any accessor ran.
func requestPhaseModel(where, argv0 string) string {
	cmd := `command = [` + strconv.Quote(argv0) +
		`, "{artifact}", "{tag.branch}"]`

	read := `[read.state]
role = "state"
path = "flow.status"
keys = ["status"]
timeout = "10s"
`
	if where == "read" {
		read = `[read.state]
role = "state"
` + cmd + `
output = "raw"
exit_absent = [1]
keys = ["status"]
timeout = "10s"
`
	}

	gate, gateRef := "", ""
	if where == "gate" {
		gate = `[gate.approval]
role = "state"
` + cmd + `
exit_verdicts = { "0" = "allow", "1" = "deny" }
keys = ["status"]
timeout = "10s"
`
		gateRef = "gate = [\"approval\"]\n"
	}

	// The write arm is a COMMAND-backed writer with NO `edit` table at
	// all, which is what makes it the sharp case: the invoking phase is
	// `phaseWrite` and there is no declared line edit anywhere in the
	// model, so a discriminator reading the phase — or the writer's
	// carrier — reports a rule `<id>.edit.<key>` that does not exist.
	write := `[write.state]
role = "state"
path = "flow.status"
keys = ["status"]
timeout = "10s"
read_back = true
`
	if where == "write" {
		write = `[write.state]
role = "state"
` + cmd + `
keys = ["status"]
timeout = "10s"
read_back = true
`
	}

	return `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "reqphase"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true
required = true

[tags.branch]
provenance = "observed"
kind = "scalar"

` + read + `
` + write + `
` + gate + `
[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "advance-draft"
` + gateRef + `[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"
`
}

// requestPhaseState is the seeded owned state the FILE-backed reader
// serves. The one rule matches on `status = "draft"`, so without it the
// invocation refuses at `flow-no-match` before any accessor under test
// runs — `[initial]` is the model's starting POINT, not a default the
// matcher falls back to for an absent key.
const requestPhaseState = `{"status":"draft"}` + "\n"

// requestPhaseFixture lays the model and a seeded state file down and
// returns the two argv fragments an invocation needs.
//
// The command-backed accessor never opens that file: its refusal is
// decided during argv substitution, before any child is spawned. The file
// is there so the FILE-backed accessors succeed and selection reaches the
// phase under test.
func requestPhaseFixture(
	t *testing.T, where, argv0 string,
) (model, binding string) {
	t.Helper()

	dir := t.TempDir()
	model = writeModelFile(t, dir, "req-phase.toml",
		requestPhaseModel(where, argv0))
	target := filepath.Join(dir, "state.json")
	if err := os.WriteFile(target, []byte(requestPhaseState), 0o644); err != nil {
		t.Fatalf("seeding the state file: %v", err)
	}
	return model, artifactBinding(requestPhaseRole, target)
}

// requestPhaseChild writes a child that records its own argv and answers
// the reader's declared `status` key, so a run that gets PAST substitution
// leaves proof it did.
func requestPhaseChild(t *testing.T, dir, trace string) string {
	t.Helper()

	p := filepath.Join(dir, "child")
	body := "#!/bin/sh\nprintf '%s\\n' \"$2\" >> " + strconv.Quote(trace) +
		"\nprintf draft\n"
	if err := os.WriteFile(p, []byte(body), 0o700); err != nil {
		t.Fatalf("writing the child: %v", err)
	}
	return p
}

// --- the read and gate arms ----------------------------------------------

// A command-backed READER or GATE whose `{tag.<key>}` is unbound is
// reported as what it is. The envelope names the role that ran, carries no
// edit findings, and stays in the exit-2 REQUEST group.
func TestDeclaredRequestRefusalNamesTheAccessorRoleThatActuallyRan(t *testing.T) {
	for _, tc := range []struct {
		name  string
		where string
		// subject is the phrase the message must carry: the role the
		// accessor ran under, spelled as that phase's other refusals
		// already spell it.
		subject string
	}{
		{name: "read_phase", where: "read", subject: "read accessor `state`"},
		{name: "gate_phase", where: "gate", subject: "gate `approval`"},
		// The sharp one. `phaseWrite` is NOT a proxy for "a declared line
		// edit refused": this model has no `edit` table anywhere, so a
		// discriminator reading the phase reports a rule that does not
		// exist. The writer's CARRIER is no better an axis — an edit
		// writer's read-back precondition is `Edit != nil` and is still
		// about the reader's argv — which is why the split is on the
		// refusal's ORIGIN.
		{name: "write_phase_command_writer", where: "write",
			subject: "write accessor `state`"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// The child is a script that never runs: every refusal under
			// test is decided before spawn. Naming a real path keeps the
			// model free of a load-time argv0 objection.
			model, binding := requestPhaseFixture(t, tc.where,
				requestPhaseChild(t, t.TempDir(),
					filepath.Join(t.TempDir(), "trace")))

			// The write arm needs the verb that WRITES; `resolve` applies
			// nothing and would never reach a write accessor.
			args := []string{"flow", "resolve", "--model", model,
				"--artifact", binding, "--outcome", "advance",
				"--allow-commands"}
			if tc.where == "write" {
				args = []string{"flow", "set-state", "--model", model,
					"--artifact", binding, "--write", "status=final",
					"--allow-commands"}
			}

			_, _, err := runCmd(t, args...)
			if err == nil {
				t.Fatalf("the invocation succeeded; want a refusal from the " +
					"unbound `{tag.branch}` placeholder")
			}
			var ce *clierr.CLIError
			if !asCLIError(err, &ce) {
				t.Fatalf("the refusal is not a structured CLIError: %v", err)
			}

			// The exit group, unchanged by the split. Exit 3 would promise
			// the caller that re-running the same request unchanged may
			// succeed; an unbound placeholder is a request they must fix.
			if got := clierr.ExitCodeFor(err); got != 2 {
				t.Errorf("exit code = %d; want 2 — an unbound `{tag.<key>}` is "+
					"about the REQUEST, and dropping this arm back to the "+
					"generic execution-failure code would regress it to 3", got)
			}
			if ce.Group != clierr.GroupUserEnv {
				t.Errorf("group = %v; want GroupUserEnv", ce.Group)
			}

			// THE defect. The message named a LINE EDIT for a refusal no
			// declared line edit produced, and on the read and gate arms
			// it named a write accessor that never ran at all.
			if !strings.Contains(ce.Message, tc.subject) {
				t.Errorf("message = %q; want it to name %q — naming the role "+
					"that actually ran is the whole point", ce.Message,
					tc.subject)
			}
			// Asserted on EVERY arm including the write one, where a write
			// accessor genuinely did run: none of these models declares an
			// `edit` table, so no invocation here has a rule
			// `<id>.edit.<key>` to name and none may claim one.
			if strings.Contains(ce.Message, "line edit") {
				t.Errorf("message = %q; this model declares no `edit` table, "+
					"so there is no declared line edit to have refused",
					ce.Message)
			}

			// No `findings[]`: the edit-findings builder reports a RULE id
			// `<id>.edit.<key>`, and substitution refused before any edit
			// rule was consulted. A finding naming a rule that does not
			// exist is worse than none.
			if len(ce.Findings) != 0 {
				t.Errorf("the refusal carries %d finding(s): %+v — "+
					"`editRefusalFindings` reports a rule id, and no edit rule "+
					"is in play when argv substitution refuses",
					len(ce.Findings), ce.Findings)
			}

			// The offending placeholder still reaches the caller. It rides
			// `Detail`, which `withStderrTail` appends, so the split needs
			// no new carriage — but the diagnosis must not be lost with the
			// findings.
			if !strings.Contains(ce.Detail, "{tag.branch}") {
				t.Errorf("detail = %q; want it to name the offending "+
					"placeholder, which is the one fact that makes the refusal "+
					"actionable", ce.Detail)
			}
		})
	}
}

// --- the write arm, unmoved ----------------------------------------------

// The other direction. C1.3's own subject — a declared line edit refused
// before mutation — keeps its envelope and its `findings[]` carriage, and
// the two arms answer under DIFFERENT codes.
//
// Neither code is compared to a literal: C1.3 makes both spellings
// non-normative. What is asserted is that they DIFFER, which is exactly
// what a caller distinguishing a read refusal from a write refusal needs
// and exactly what the defect denied them.
func TestTheWritePhaseKeepsTheEditRefusalEnvelopeAndFindings(t *testing.T) {
	// The write arm reuses the `edit` fixture, whose refusal is minted
	// inside the edit writer rather than at substitution — the one site
	// C1.3 was written about.
	model, binding := editEnvelopeFixture(t, "# a record\n\nno status bullet\n")

	_, _, werr := runCmd(t, "flow", "set-state", "--model", model,
		"--artifact", binding, "--write", "status=Final")
	if werr == nil {
		t.Fatalf("the write invocation succeeded; want an edit refusal")
	}
	var write *clierr.CLIError
	if !asCLIError(werr, &write) {
		t.Fatalf("the write refusal is not a structured CLIError: %v", werr)
	}
	if got := clierr.ExitCodeFor(werr); got != 2 {
		t.Errorf("write-phase exit code = %d; want 2", got)
	}
	if len(write.Findings) == 0 {
		t.Errorf("the write-phase refusal carries no `findings[]`; C1.3's "+
			"carriage of the rule id and the reason token is unmoved by the "+
			"phase split: %+v", write)
	}
	if !strings.Contains(write.Message, "write accessor") {
		t.Errorf("write-phase message = %q; want it to name the write "+
			"accessor", write.Message)
	}

	// The read arm, for the code comparison.
	readModel, readBinding := requestPhaseFixture(t, "read",
		requestPhaseChild(t, t.TempDir(),
			filepath.Join(t.TempDir(), "trace")))
	_, _, rerr := runCmd(t, "flow", "resolve", "--model", readModel,
		"--artifact", readBinding, "--outcome", "advance", "--allow-commands")
	if rerr == nil {
		t.Fatalf("the read invocation succeeded; want a refusal")
	}
	var read *clierr.CLIError
	if !asCLIError(rerr, &read) {
		t.Fatalf("the read refusal is not a structured CLIError: %v", rerr)
	}

	if read.Code == write.Code {
		t.Errorf("both phases answer under the code %q; a caller counting "+
			"write refusals cannot tell a reader's unbound placeholder from an "+
			"edit writer's refusal, which is the defect", read.Code)
	}
}

// --- the positive control -------------------------------------------------

// The reader's command IS reached and its placeholder IS substituted when
// the tag is bound, which is what stops the refusal assertions above from
// passing against a model that refuses whatever it is handed.
//
// The oracle is the CHILD'S OWN ARGV, recorded to a trace file. A build
// where command entries never work at all, or one that refused every
// `{tag.<key>}` bound or not, would leave that file absent — and the
// refusals asserted above would then be proving nothing about the UNBOUND
// case in particular.
func TestABoundTagPlaceholderReachesTheChildAndIsSubstituted(t *testing.T) {
	dir := t.TempDir()
	trace := filepath.Join(dir, "trace")
	model, binding := requestPhaseFixture(t, "read",
		requestPhaseChild(t, dir, trace))

	_, _, err := runCmd(t, "flow", "resolve", "--model", model,
		"--artifact", binding, "--outcome", "advance", "--allow-commands",
		"--tag", "branch=main")
	if err != nil {
		t.Fatalf("a BOUND `{tag.branch}` was refused: %v", err)
	}

	got, rerr := os.ReadFile(trace)
	if rerr != nil {
		t.Fatalf("no child ran: %v — substitution is reached only when the "+
			"invocation gets past the placeholder checks, so an absent trace "+
			"means the unbound-case refusals above prove nothing", rerr)
	}
	if strings.TrimSpace(string(got)) != "main" {
		t.Errorf("the child saw %q as its second argument; want the BOUND "+
			"value `main` — the placeholder is replaced whole-element and is "+
			"never passed through literally", strings.TrimSpace(string(got)))
	}
}

// --- the edit writer whose READ-BACK reader is the one at fault ----------

// editReadBackRole is the artifact role the edit writer and its
// command-backed read-back reader both bind.
const editReadBackRole = "record"

// editReadBackModel is an `edit`-carried writer — a REAL rule table, so
// `Edit != nil` — whose read-back goes through a COMMAND-backed reader
// carrying `{tag.branch}`.
//
// This is the case that separates ORIGIN from CARRIER, and the only one
// that can. Every other fixture in this file agrees with a carrier test by
// accident: the reader, gate and command-writer arms all have no `edit`
// table, and `edit_envelope_0028_test.go`'s writer has one and genuinely
// does refuse a line edit. Here the two answers diverge — `Edit != nil`,
// yet the refusal is about the READER's argv and names no rule.
//
// `{tag.branch}` sits at a later argv position, never argv0, which C1.4's
// `edit_tag_argv0` refuses at LOAD.
func editReadBackModel(argv0 string) string {
	return `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "editrb"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["Draft", "Final"]
single_valued = true
required = true

[tags.branch]
provenance = "observed"
kind = "scalar"

[read.record]
role = "record"
command = [` + strconv.Quote(argv0) + `, "{artifact}", "{tag.branch}"]
output = "raw"
exit_absent = [1]
keys = ["status"]
timeout = "10s"

[write.record]
role = "record"
keys = ["status"]
timeout = "10s"
read_back = true

[write.record.edit.status]
anchor  = "^- \\*\\*Status\\*\\*: (.+)$"
replace = "- **Status**: {status}"

[context.done]
[context.done.match.status]
eq = "Final"

[initial]
status = "Draft"
`
}

// editReadBackDoc is the record the edit writer would rewrite if it ran.
const editReadBackDoc = `# a record

- **Status**: Draft
- **Owner**: cwensel
`

// An EDIT writer whose command-backed read-back reader carries a bad
// `{tag.<key>}` refuses under the GENERIC envelope, because the reader's
// argv is what is wrong and no rule refused.
//
// WHY THIS CASE EXISTS. It is the one that pins the discriminator's AXIS.
// The refusal arrives at `phaseWrite` from a definition with `Edit != nil`,
// so both rejected proxies — the invoking phase and the writer's carrier —
// answer "line edit" here, and both are wrong: `editRefusalFindings` would
// code a complaint about the READER's argv as a rule `<id>.edit.<key>` that
// no one declared. Only the refusal's ORIGIN gets it right.
//
// Verified against the mutant it exists to catch: rewriting
// `accessorFailureOf`'s branch to test `len(def.Accessor.Edit) != 0`
// instead of `DeclaredEdit()` leaves every OTHER test in the repository
// green and turns this one red.
func TestAnEditWritersReadBackArgvRefusalIsNotALineEditRefusal(t *testing.T) {
	for _, tc := range []struct {
		name string
		// tag is the `--tag branch=…` binding, or empty for none. The two
		// cases are the two argv preconditions C1.6 names, which C1.3's
		// `precedence:` orders unbound-first — so a fix wired at only one
		// of them is caught.
		tag string
	}{
		{name: "unbound_placeholder", tag: ""},
		{name: "flag_shaped_value", tag: "branch=-rf"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			// A child that would answer the read-back correctly if it ever
			// ran, so the refusal cannot be mistaken for the child failing.
			child := filepath.Join(dir, "child")
			if err := os.WriteFile(child,
				[]byte("#!/bin/sh\nprintf Draft\n"), 0o700); err != nil {
				t.Fatalf("writing the child: %v", err)
			}
			model := writeModelFile(t, dir, "edit-rb.toml",
				editReadBackModel(child))
			target := filepath.Join(dir, "record.md")
			if err := os.WriteFile(target,
				[]byte(editReadBackDoc), 0o644); err != nil {
				t.Fatalf("writing the target: %v", err)
			}

			args := []string{"flow", "set-state", "--model", model,
				"--artifact", artifactBinding(editReadBackRole, target),
				"--write", "status=Final", "--allow-commands"}
			if tc.tag != "" {
				args = append(args, "--tag", tc.tag)
			}

			_, _, err := runCmd(t, args...)
			if err == nil {
				t.Fatalf("the invocation succeeded; want a refusal about the " +
					"read-back reader's argv")
			}
			var ce *clierr.CLIError
			if !asCLIError(err, &ce) {
				t.Fatalf("the refusal is not a structured CLIError: %v", err)
			}

			// Exit 2. The reader's argv is a defect of the REQUEST, and
			// exit 3's "re-run unchanged" is advice that can never work.
			if got := clierr.ExitCodeFor(err); got != 2 {
				t.Errorf("exit code = %d; want 2", got)
			}
			if ce.Group != clierr.GroupUserEnv {
				t.Errorf("group = %v; want GroupUserEnv", ce.Group)
			}

			// THE assertion. No rule refused, so no rule may be named.
			if strings.Contains(ce.Message, "line edit") {
				t.Errorf("message = %q; the refusal is about the read-back "+
					"READER's argv and no `edit` rule was consulted, so naming "+
					"a declared line edit invents a subject — this is what a "+
					"carrier-based discriminator gets wrong, because the "+
					"definition here does carry an `edit` table", ce.Message)
			}
			if len(ce.Findings) != 0 {
				t.Errorf("the refusal carries %d finding(s): %+v — "+
					"`editRefusalFindings` codes its entry as a line-edit "+
					"refusal and its Detail reads as a rule `<id>.edit.<key>`; "+
					"there is no such rule here", len(ce.Findings), ce.Findings)
			}

			// The diagnosis still reaches the caller, naming the reader
			// that is actually at fault.
			if !strings.Contains(ce.Detail, "read-back") {
				t.Errorf("detail = %q; want it to name the read-back as the "+
					"thing that could not be satisfied", ce.Detail)
			}

			// NO MUTATION. C1.3 `order:` requires every refusal in the
			// clause to be decided before any byte is written, and this
			// one is raised ahead of `Apply` precisely so the caller is
			// never told a write "may have been applied" for what is
			// purely a defect of the request.
			got, rerr := os.ReadFile(target)
			if rerr != nil {
				t.Fatalf("reading the target back: %v", rerr)
			}
			if string(got) != editReadBackDoc {
				t.Errorf("the target was mutated:\n got %q\nwant %q — this "+
					"refusal is decided BEFORE any byte is written",
					got, editReadBackDoc)
			}
		})
	}
}
