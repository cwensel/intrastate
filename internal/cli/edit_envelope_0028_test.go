package cli

// RDR 0028 — the CLI ENVELOPE half of C1.3's exit-group clause.
//
// The accessor-seam half is held in `internal/accessor`: the refusal's
// CLASS is 0004's unchanged `execution_failure`, `Applied()` is false, and
// the Detail carries the rule id and the reason token. None of that is
// observable at the envelope, and none of what this file pins is
// observable at that seam — which is why the two are separate REQs and why
// coverage listed these three as orphans reachable only here.
//
// WHAT THIS FILE PINS, and what it must NOT.
//
// C1.3's EXIT GROUP clause fixes three things and explicitly declines to
// fix a fourth:
//
//  1. the exit GROUP — 2, not `execution_failure`'s default 3;
//  2. the DISCRIMINATOR at `flow_exec.go::accessorFailureOf`, which is
//     what routes one class to two groups;
//  3. the `findings[]` CARRIAGE of the rule id and the reason token, per
//     `docs/cli-output-contract.md`.
//
// The fourth is the CLI code's SPELLING, which the clause states outright
// is "Stage 8's, non-normative here" (REQ-44, restated at S27 as REQ-113).
// So every assertion below reads `clierr.ExitCodeFor` and the findings
// shape, and NONE compares `ce.Code` to a literal. A test that pinned the
// string would convert an implementation's free choice into a contract and
// would go red on a rename that breaks nothing.
//
// THE REGRESSION THIS EXISTS TO CATCH is an exit-3 assertion. Exit 3
// promises a caller "repair the environment and re-run the same request
// unchanged". A stale anchor is not an environment fault: re-running the
// same request spins the caller forever on an input they must instead fix.
// `execution_failure` defaults to exit 3 at `accessorFailureOf`, so the
// edit refusals reach the right group only through a live discriminator —
// and the negative control below is what proves the discriminator
// discriminates rather than moving the whole class.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/table"
)

// --- the fixture ----------------------------------------------------------

// editEnvelopeRole is the artifact role the edit writer and its reader
// both bind.
const editEnvelopeRole = "record"

// editEnvelopeModel is an `edit`-carried flow whose reader is FILE-backed,
// so the gate is out of the picture and the only refusals reaching the
// envelope are the ones C1.3 mints.
//
// The anchor is pinned `^…$` and also matches its own output, so a
// well-anchored write applies and only a fixture that earns it refuses.
const editEnvelopeModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "editenv"
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

[read.record]
role = "record"
path = "record.state"
keys = ["status"]
timeout = "2s"

[write.record]
role = "record"
keys = ["status"]
timeout = "2s"
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

// editEnvelopeDoc is the record shape the writer edits.
const editEnvelopeDoc = `# a record

- **Status**: Draft
- **Owner**: cwensel
`

// editEnvelopeFixture lays down the model and a target document, and
// returns the two argv fragments an invocation needs.
func editEnvelopeFixture(t *testing.T, doc string) (model, binding string) {
	t.Helper()

	dir := t.TempDir()
	model = writeModelFile(t, dir, "edit-env.toml", editEnvelopeModel)
	target := filepath.Join(dir, "record.md")
	if err := os.WriteFile(target, []byte(doc), 0o644); err != nil {
		t.Fatalf("writing the target document: %v", err)
	}
	return model, artifactBinding(editEnvelopeRole, target)
}

// --- REQ-43, REQ-113: the exit GROUP -------------------------------------

// REQ-43: "EXIT GROUP: these refusals are about the REQUEST, not the
// environment, so they take the exit-2 group, not `execution_failure`'s
// default exit 3 — the executor-facing typed `Err` discriminates at
// `flow_exec.go::accessorFailureOf`"
// REQ-113 (S27): "At the envelope the refusal exits **2**, not 3 ... The
// code's SPELLING is Stage 8's to choose and no test may pin the string;
// what this scenario pins is the exit group, the discriminator and the
// `findings[]` carriage — an exit-3 assertion here is the regression this
// scenario exists to catch"
// ADVERSARIAL
//
// Three fixtures, each earning a different reason token, all reaching the
// same group. One fixture would leave "the group is right for THIS token"
// untested — and the tokens are minted at three different steps of C1.3's
// precedence ladder, so a discriminator wired at only one of them passes a
// single-fixture test.
func TestReq43And113_ADeclaredLineEditRefusalTakesTheExit2Group(t *testing.T) {
	for _, tc := range []struct {
		name  string
		doc   string
		token string
	}{
		{
			name:  "unmatched_anchor",
			doc:   "# a record\n\nno status bullet here\n",
			token: "edit_anchor_unmatched",
		},
		{
			name:  "ambiguous_anchor",
			doc:   editEnvelopeDoc + "- **Status**: Draft\n",
			token: "edit_anchor_ambiguous",
		},
		{
			name:  "undeclared_clear",
			doc:   editEnvelopeDoc,
			token: "edit_clear_undeclared",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			model, binding := editEnvelopeFixture(t, tc.doc)

			args := []string{"flow", "set-state", "--model", model,
				"--artifact", binding, "--write", "status=Final"}
			if tc.token == "edit_clear_undeclared" {
				args = []string{"flow", "set-state", "--model", model,
					"--artifact", binding, "--clear", "status"}
			}

			_, _, err := runCmd(t, args...)
			if err == nil {
				t.Fatalf("the invocation succeeded; want a refusal carrying %s",
					tc.token)
			}

			// THE assertion. Exit 2 is the request group; exit 3 would
			// tell the caller to re-run unchanged, which cannot help.
			if got := clierr.ExitCodeFor(err); got != 2 {
				t.Errorf("exit code = %d; want 2 — these refusals are about the "+
					"REQUEST, and exit 3 promises the caller that re-running the "+
					"same request unchanged may succeed", got)
			}

			var ce *clierr.CLIError
			if !asCLIError(err, &ce) {
				t.Fatalf("the refusal is not a structured CLIError: %v", err)
			}
			if ce.Group != clierr.GroupUserEnv {
				t.Errorf("group = %v; want GroupUserEnv — the group is what "+
					"fixes the exit, and the code's spelling is not the contract",
					ce.Group)
			}
		})
	}
}

// REQ-43 (discriminator half): "the executor-facing typed `Err`
// discriminates at `flow_exec.go::accessorFailureOf`"
// ADVERSARIAL
//
// The negative control, and the reason this file is not one test.
// `execution_failure` maps to exit 3 for every OTHER refusal in that class,
// and an implementation that simply moved the whole class to exit 2 would
// pass every assertion above while breaking 0005's contract for a genuine
// environment failure. So the SAME class, reached through a refusal the
// binding did NOT mint about the request — an unreadable target — must
// still exit 3.
//
// C1.3 states this scoping by hand: a target that cannot be read "refuses
// BEFORE mutation with the OS error in the Detail", and that one really is
// the environment, where re-running after fixing permissions succeeds.
func TestReq43_TheDiscriminatorLeavesGenuineEnvironmentFailuresAtExit3(t *testing.T) {
	dir := t.TempDir()
	model := writeModelFile(t, dir, "edit-env.toml", editEnvelopeModel)
	// The bound target does not exist. `edit` fences creation out, so this
	// is an OS read error and not an anchor result.
	binding := artifactBinding(editEnvelopeRole, filepath.Join(dir, "absent.md"))

	_, _, err := runCmd(t, "flow", "set-state", "--model", model,
		"--artifact", binding, "--write", "status=Final")
	if err == nil {
		t.Fatalf("a write against an absent target succeeded")
	}
	if got := clierr.ExitCodeFor(err); got != 3 {
		t.Errorf("exit code = %d; want 3 — an unreadable target IS an "+
			"environment failure, and a discriminator that moved the whole "+
			"`execution_failure` class to exit 2 would break 0005's contract "+
			"for every other member of it", got)
	}
}

// --- REQ-114: the `findings[]` carriage ----------------------------------

// REQ-114 (S27): "Refusal reporting follows
// `docs/cli-output-contract.md`: a scalar failure carries `param`, an
// aggregate carries `findings[]`, regardless of runtime cardinality."
// REQ-41 (envelope half): the rule id `<id>.edit.<key>` and the reason
// token reach the caller.
// BOUNDARY
//
// The contract splits the two carriers by FAMILY and says so explicitly —
// "regardless of runtime cardinality". This family's subject is a RULE,
// not a flag or an argument, and one entry can carry several rules, so it
// is an aggregate family: `findings[]` even at one entry, and no top-level
// `param`.
//
// The oracle reads the marshalled ENVELOPE rather than the Go struct,
// because `findings` is what a consumer parses and `omitempty` means a
// struct field can be populated while the wire key elides.
func TestReq114_TheRefusalCarriesTheRuleAndTokenInFindings(t *testing.T) {
	model, binding := editEnvelopeFixture(t, "# a record\n\nno status bullet\n")

	_, _, err := runCmd(t, "flow", "set-state", "--model", model,
		"--artifact", binding, "--write", "status=Final")
	if err == nil {
		t.Fatalf("the invocation succeeded; want a refusal")
	}
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("the refusal is not a structured CLIError: %v", err)
	}

	raw, merr := json.Marshal(ce)
	if merr != nil {
		t.Fatalf("marshalling the envelope: %v", merr)
	}
	var env struct {
		Code     string `json:"code"`
		Param    string `json:"param"`
		Findings []struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"findings"`
	}
	if uerr := json.Unmarshal(raw, &env); uerr != nil {
		t.Fatalf("decoding the envelope: %v", uerr)
	}

	if len(env.Findings) == 0 {
		t.Fatalf("the envelope carries no `findings[]`:\n%s\n— an aggregate "+
			"family reports in `findings[]` regardless of runtime cardinality", raw)
	}
	if env.Param != "" {
		t.Errorf("the envelope carries a top-level `param` of %q; an aggregate "+
			"family names its subjects in `findings[]` and not in `param`", env.Param)
	}

	// The rule id and the reason token both reach the caller, which is
	// what makes the refusal diagnosable without re-running under a
	// debugger. Asserted over the findings' own text, never over `detail`
	// alone: `detail` is the ONE free-text slot and a consumer that
	// branches on structure reads `findings`.
	joined := env.Findings[0].Message
	if !strings.Contains(joined, "edit_anchor_unmatched") {
		t.Errorf("the finding %q does not carry the reason token", joined)
	}
	if !strings.Contains(joined, "record.edit.status") {
		t.Errorf("the finding %q does not name the rule `<id>.edit.<key>`", joined)
	}

	// REQ-44 / REQ-113: the code's SPELLING is not asserted anywhere in
	// this file. What IS asserted is that the finding carries the SAME
	// code as the envelope, so a caller correlating the two is not left
	// matching one string against another.
	if env.Findings[0].Code != env.Code {
		t.Errorf("finding code = %q but envelope code = %q; a caller "+
			"correlating the two reads one code, not two",
			env.Findings[0].Code, env.Code)
	}
}

// --- the positive control -------------------------------------------------

// The `edit` binding IS REACHED through the CLI and DOES rewrite the
// document, which is what stops every refusal assertion above from passing
// vacuously against an envelope that refuses whatever it is handed.
//
// The oracle is the DOCUMENT, not the exit code. This model's reader is
// file-backed over the same caller-bound path the writer edits, so the
// post-mutation read-back cannot parse a markdown target and the
// invocation refuses at the read-back — a 0004:C12/C13 disposition that is
// downstream of everything C1.3 owns. C1.6 names exactly this: the
// consumer's own PROJECTOR is the reader for a document target, and
// REQ-MVV-BLOCK records that the projector verb does not ship, so a
// green-exit end-to-end run is Phase 4's and not reachable here.
//
// What is reachable, and what this asserts, is that the write LANDED
// before that read-back ran: the anchored line carries the planned value
// and every other byte is unchanged. A binding the envelope never reached
// would leave the document untouched, and the refusal assertions above
// would then be proving nothing about the `edit` carrier at all.
func TestReq43_TheEditBindingIsReachedAndRewritesTheAnchoredLine(t *testing.T) {
	dir := t.TempDir()
	model := writeModelFile(t, dir, "edit-env.toml", editEnvelopeModel)
	target := filepath.Join(dir, "record.md")
	if err := os.WriteFile(target, []byte(editEnvelopeDoc), 0o644); err != nil {
		t.Fatalf("writing the target: %v", err)
	}

	_, _, err := runCmd(t, "flow", "set-state", "--model", model,
		"--artifact", artifactBinding(editEnvelopeRole, target),
		"--write", "status=Final", "--as=json")

	// The read-back's own disposition is 0004's and not this record's, so
	// it is not the oracle. What must NOT happen is a refusal C1.3 mints:
	// those are decided BEFORE mutation and the document below proves the
	// mutation ran.
	if err != nil {
		var ce *clierr.CLIError
		if asCLIError(err, &ce) && strings.Contains(ce.Detail, "edit_") {
			t.Fatalf("a well-anchored edit was refused by C1.3: %v", ce.Detail)
		}
	}

	got, rerr := os.ReadFile(target)
	if rerr != nil {
		t.Fatalf("reading the target back: %v", rerr)
	}
	want := strings.Replace(editEnvelopeDoc,
		"- **Status**: Draft", "- **Status**: Final", 1)
	if string(got) != want {
		t.Errorf("post-edit document:\n got %q\nwant %q — exactly the anchored "+
			"line is rewritten and every other byte is unchanged", got, want)
	}
}

// REQ-102: "`docs/cli-output-contract.md` names the new categories and
// Detail tokens" — (0028 Phase 4, §implementation-plan).
//
// DOMAIN EDGE. The doc is the consumer-facing contract, and C1.4's six
// load-time wire strings and C1.3/C1.5's six apply-time reason tokens are
// DISJOINT sets: only the former register in `table.Categories()`. A
// consumer matching an apply-time token against the load-category list
// finds nothing, so the doc has to name both sets and say they are
// separate. This test is the drift gate on that paragraph — it pins the
// NAMES, which are wire strings, not the prose around them.
func TestReq102_TheOutputContractNamesBothEditNameSets(t *testing.T) {
	body, err := os.ReadFile(filepath.Join("..", "..", "docs", "cli-output-contract.md"))
	if err != nil {
		t.Fatalf("read cli-output-contract.md: %v", err)
	}
	doc := string(body)

	// C1.4's six load-time categories, in the order Categories() registers
	// them. Sourced from the table package so a rename cannot drift the
	// doc and the code apart silently.
	loadTime := []table.Category{
		table.CatEditCarrierConflict,
		table.CatEditKeyMismatch,
		table.CatEditAnchorInvalid,
		table.CatEditTemplateInvalid,
		table.CatEditClearInvalid,
		table.CatEditTagArgv0,
	}
	for _, c := range loadTime {
		if !strings.Contains(doc, string(c)) {
			t.Errorf("cli-output-contract.md does not name load-time category %q", c)
		}
	}

	// C1.3 and C1.5's six apply-time reason tokens. These are unexported
	// in flowbind, so they are spelled here as the wire strings the doc
	// must carry; the disjointness assertion below is what keeps that
	// honest.
	applyTime := []string{
		"edit_anchor_unmatched",
		"edit_anchor_ambiguous",
		"edit_anchor_collision",
		"edit_anchor_unstable",
		"edit_value_multiline",
		"edit_clear_undeclared",
	}
	for _, tok := range applyTime {
		if !strings.Contains(doc, tok) {
			t.Errorf("cli-output-contract.md does not name apply-time reason token %q", tok)
		}
	}

	// The disjointness the doc asserts must be TRUE of the shipped code,
	// or the doc documents a distinction that does not exist.
	registered := make(map[string]bool)
	for _, c := range table.Categories() {
		registered[string(c)] = true
	}
	for _, tok := range applyTime {
		if registered[tok] {
			t.Errorf("apply-time reason token %q is registered in table.Categories(); "+
				"C1.4's load-time set and C1.3/C1.5's apply-time set must stay disjoint", tok)
		}
	}
	for _, c := range loadTime {
		if !registered[string(c)] {
			t.Errorf("load-time category %q is NOT registered in table.Categories()", c)
		}
	}
}
