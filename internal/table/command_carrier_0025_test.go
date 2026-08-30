package table_test

// RDR 0025 — the command carrier and its six load-time defect classes
// (C1, C5; REQ-1..REQ-15, REQ-71..REQ-87, REQ-115, REQ-122).
//
// The coverage floor this file exists to hold: an implementation that
// accepts `command` but never registers its defect categories refuses
// correctly at the call site while staying invisible to every consumer
// that enumerates `table.Categories()` (C5, REQ-85). So every mutant is
// asserted twice — by the category the refusal carries AND by that
// category's membership in the registered set — and the registration
// order is asserted as a tail append, because `Categories()` returns a
// literal slice and a golden keyed on its order moves under an insertion.
//
// The second mutant this file kills is the closed-set reading of C5's
// interpreter deny-list. `perl -e` is an unlisted interpreter and MUST
// LOAD (REQ-83, REQ-122): a suite that only asserts refusals reads as
// "inline shell is impossible", which is exactly the guarantee C5
// declines to make.

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/table"
)

// --- the command-carrier fixture -----------------------------------------

// cmdCarrierModel is a minimal model whose read/write/gate entries are
// command-backed. It is authored, not derived from an existing fixture,
// because the whole point is that `path` is ABSENT on these entries — a
// substitution into a path-backed fixture would leave the old locator.
const cmdCarrierModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "cmdflow"
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

[read.state]
role = "state"
command = ["git", "config", "--file", "{artifact}", "--get", "flow.status"]
output = "raw"
exit_absent = [1]
keys = ["status"]
timeout = "2s"

[write.state]
role = "state"
command = ["tools/flowstate-write", "{artifact}"]
keys = ["status"]
timeout = "2s"
read_back = true

[gate.approval]
role = "state"
command = ["git", "-C", "{artifact}", "diff", "--quiet"]
exit_verdicts = { "0" = "allow", "1" = "deny" }
keys = ["status"]
timeout = "2s"

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "advance-draft"
gate = ["approval"]
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"
`

// swapEntry replaces one whole capability-table block in the fixture. It
// fails loudly if the substitution does not apply, so a mutant that
// silently equals the base can never pass as a refusal test.
func swapEntry(t *testing.T, base, old, replacement string) string {
	t.Helper()

	out := strings.Replace(base, old, replacement, 1)
	if out == base {
		t.Fatalf("mutant substitution did not apply; the block\n%s\nis not in the fixture", old)
	}
	return out
}

const cmdReadBlock = `[read.state]
role = "state"
command = ["git", "config", "--file", "{artifact}", "--get", "flow.status"]
output = "raw"
exit_absent = [1]
keys = ["status"]
timeout = "2s"`

// loadCategoryOf loads src and returns the category its refusal carries.
// The oracle is the CATEGORY, never "an error occurred" and never the
// message text (C5, S1: "asserted by category, never by 'validation
// returned non-empty'").
func loadCategoryOf(t *testing.T, src, id string) table.Category {
	t.Helper()

	_, err := table.Load([]byte(src), id)
	if err == nil {
		t.Fatalf("%s loaded clean; want a categorized refusal", id)
	}
	cat, ok := table.CategoryOf(err)
	if !ok {
		t.Fatalf("%s refused with a non-categorized error: %v", id, err)
	}
	return cat
}

// cmdFailureOf returns the *table.Failure a load refusal carries, so a
// test can read the failure DATA rather than a rendered message.
func cmdFailureOf(t *testing.T, err error) *table.Failure {
	t.Helper()

	var f *table.Failure
	if !errors.As(err, &f) {
		t.Fatalf("refusal is not a *table.Failure: %v", err)
	}
	return f
}

// --- C1: the carrier fields ----------------------------------------------

// REQ-1: "command       = [\"<argv0>\", \"<arg>\", ...]   # []string"
// HAPPY PATH
//
// The oracle is the DECODED vector on the loaded model, not merely that
// the document loaded: a loader that accepts the field and discards it
// would pass a load-clean assertion and lose the authority bound.
func TestReq1_CommandDecodesAsTheEntrysArgvVector(t *testing.T) {
	m, err := table.Load([]byte(cmdCarrierModel), "cmd-carrier.toml")
	if err != nil {
		t.Fatalf("a command-backed model must load clean; refused: %v", err)
	}

	want := []string{"git", "config", "--file", "{artifact}", "--get", "flow.status"}
	got := m.Readers["state"].Command
	if !slices.Equal(got, want) {
		t.Errorf("read.state command = %#v; want the declared vector %#v", got, want)
	}
}

// REQ-2: "output        = \"json\" | \"raw\"              # *string, read
// entries only; default \"json\" when omitted"
// REQ-13: "`output` is a pointer so that omitted and explicit-`\"json\"`
// are distinguishable at load"
// INPUT EDGE
//
// PHASE-0 READING (ASSUMPTION REQ-2): `output` admits exactly the two
// literals; any other string is `command_output_shape`.
func TestReq2_OutputCarriesJSONOrRawAndDistinguishesOmission(t *testing.T) {
	t.Run("raw_decodes", func(t *testing.T) {
		m, err := table.Load([]byte(cmdCarrierModel), "cmd-carrier.toml")
		if err != nil {
			t.Fatalf("fixture must load clean: %v", err)
		}
		out := m.Readers["state"].Output
		if out == nil {
			t.Fatal("read.state Output is nil; the entry declares output = \"raw\"")
		}
		if *out != "raw" {
			t.Errorf("read.state output = %q; want %q", *out, "raw")
		}
	})

	t.Run("omitted_is_distinguishable_from_explicit_json", func(t *testing.T) {
		omitted := swapEntry(t, cmdCarrierModel, cmdReadBlock, `[read.state]
role = "state"
command = ["reader"]
keys = ["status"]
timeout = "2s"`)
		m, err := table.Load([]byte(omitted), "cmd-output-omitted.toml")
		if err != nil {
			t.Fatalf("an entry omitting output must load clean: %v", err)
		}
		if got := m.Readers["state"].Output; got != nil {
			t.Errorf("omitted output decoded as %q; want nil so omission is "+
				"distinguishable from an explicit \"json\" (REQ-13)", *got)
		}

		explicit := swapEntry(t, cmdCarrierModel, cmdReadBlock, `[read.state]
role = "state"
command = ["reader"]
output = "json"
keys = ["status"]
timeout = "2s"`)
		m2, err := table.Load([]byte(explicit), "cmd-output-json.toml")
		if err != nil {
			t.Fatalf("an entry declaring output = \"json\" must load clean: %v", err)
		}
		got := m2.Readers["state"].Output
		if got == nil || *got != "json" {
			t.Errorf("explicit output decoded as %v; want a non-nil %q", got, "json")
		}
	})

	t.Run("an_unknown_output_literal_is_command_output_shape", func(t *testing.T) {
		src := swapEntry(t, cmdCarrierModel, cmdReadBlock, `[read.state]
role = "state"
command = ["reader"]
output = "yaml"
keys = ["status"]
timeout = "2s"`)
		if cat := loadCategoryOf(t, src, "cmd-output-yaml.toml"); cat != table.CatCommandOutputShape {
			t.Errorf("category = %q; want %q", cat, table.CatCommandOutputShape)
		}
	})
}

// REQ-3: "exit_absent   = [<code>, ...]               # []int,
// read entries only"
// REQ-4: "exit_verdicts = { \"<code>\" = \"<verdict>\" }  # map[string]string,
// gate entries only"
// HAPPY PATH
func TestReq3And4_ExitMapsDecodeOnTheirOwnCapability(t *testing.T) {
	m, err := table.Load([]byte(cmdCarrierModel), "cmd-carrier.toml")
	if err != nil {
		t.Fatalf("fixture must load clean: %v", err)
	}

	if got := m.Readers["state"].ExitAbsent; !slices.Equal(got, []int{1}) {
		t.Errorf("read.state exit_absent = %#v; want []int{1}", got)
	}
	want := map[string]string{"0": "allow", "1": "deny"}
	got := m.Gates["approval"].ExitVerdicts
	if len(got) != len(want) {
		t.Fatalf("gate.approval exit_verdicts = %#v; want %#v", got, want)
	}
	for code, verdict := range want {
		if got[code] != verdict {
			t.Errorf("exit_verdicts[%q] = %q; want %q", code, got[code], verdict)
		}
	}
}

// REQ-5: "env           = { \"<KEY>\" = \"<value>\" }     # map[string]string"
// REQ-6: "env_pass      = [\"<VAR>\", ...]              # []string"
// HAPPY PATH
func TestReq5And6_EnvAndEnvPassDecodeOnAnEntry(t *testing.T) {
	src := swapEntry(t, cmdCarrierModel, `[write.state]
role = "state"
command = ["tools/flowstate-write", "{artifact}"]
keys = ["status"]
timeout = "2s"
read_back = true`, `[write.state]
role = "state"
command = ["tools/flowstate-write", "{artifact}"]
env = { GIT_CONFIG_NOSYSTEM = "1" }
env_pass = ["SSH_AUTH_SOCK"]
keys = ["status"]
timeout = "2s"
read_back = true`)

	m, err := table.Load([]byte(src), "cmd-env.toml")
	if err != nil {
		t.Fatalf("an entry declaring env and env_pass must load clean: %v", err)
	}
	w := m.Writers["state"]
	if w.Env["GIT_CONFIG_NOSYSTEM"] != "1" {
		t.Errorf("write.state env = %#v; want GIT_CONFIG_NOSYSTEM=1", w.Env)
	}
	if !slices.Equal(w.EnvPass, []string{"SSH_AUTH_SOCK"}) {
		t.Errorf("write.state env_pass = %#v; want [SSH_AUTH_SOCK]", w.EnvPass)
	}
}

// REQ-7: "exactly one of `path` / `command` per entry. Load-time this is
// C5 `command_and_path_conflict`"
// REQ-71: "command_and_path_conflict     # both or neither of path/command
// declared"
// REQ-14: "the path-absent-and-command-absent case moves to
// `command_and_path_conflict`. Existing fixtures asserting the old
// category on a path-less entry change category, not verdict"
// ADVERSARIAL
func TestReq7_ExactlyOneOfPathOrCommandPerEntry(t *testing.T) {
	cases := []struct {
		name  string
		entry string
	}{
		{"both carriers declared", `[read.state]
role = "state"
path = "flow.status"
command = ["reader"]
keys = ["status"]
timeout = "2s"`},
		{"neither carrier declared", `[read.state]
role = "state"
keys = ["status"]
timeout = "2s"`},
		{"empty path beside a command", `[read.state]
role = "state"
path = ""
command = ["reader"]
keys = ["status"]
timeout = "2s"`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := swapEntry(t, cmdCarrierModel, cmdReadBlock, tc.entry)
			cat := loadCategoryOf(t, src, "cmd-conflict.toml")
			if cat != table.CatCommandAndPathConflict {
				t.Errorf("category = %q; want %q — the exactly-one rule is C5's "+
					"own arm, not the inherited malformed-declaration one",
					cat, table.CatCommandAndPathConflict)
			}
		})
	}
}

// REQ-9: "`command` must be non-empty and contain no empty element."
// REQ-72: "command_empty                 # empty vector or empty argv
// element"
// REQ-24: "an empty argv0 is a C5 defect"
// BOUNDARY
func TestReq9_CommandIsNonEmptyAndCarriesNoEmptyElement(t *testing.T) {
	cases := []struct {
		name  string
		entry string
	}{
		{"empty vector", `[read.state]
role = "state"
command = []
keys = ["status"]
timeout = "2s"`},
		{"empty argv0", `[read.state]
role = "state"
command = ["", "--get"]
keys = ["status"]
timeout = "2s"`},
		{"empty trailing element", `[read.state]
role = "state"
command = ["reader", ""]
keys = ["status"]
timeout = "2s"`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := swapEntry(t, cmdCarrierModel, cmdReadBlock, tc.entry)
			if cat := loadCategoryOf(t, src, "cmd-empty.toml"); cat != table.CatCommandEmpty {
				t.Errorf("category = %q; want %q", cat, table.CatCommandEmpty)
			}
		})
	}
}

// REQ-10: "`role`, `keys`, `timeout`, and `read_back` rules are unchanged
// from RDR 0002/0004."
// DOMAIN EDGE
//
// NEGATIVE-adjacent: C1 relaxes ONLY the path rule. A loader that relaxed
// the whole switch would let a command entry ship with no role or no
// timeout, which every downstream bound depends on.
func TestReq10_TheOtherFourRulesAreUnchangedForACommandEntry(t *testing.T) {
	cases := []struct {
		name  string
		entry string
	}{
		{"absent role", `[read.state]
command = ["reader"]
keys = ["status"]
timeout = "2s"`},
		{"absent keys", `[read.state]
role = "state"
command = ["reader"]
timeout = "2s"`},
		{"absent timeout", `[read.state]
role = "state"
command = ["reader"]
keys = ["status"]`},
		{"non-positive timeout", `[read.state]
role = "state"
command = ["reader"]
keys = ["status"]
timeout = "0s"`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := swapEntry(t, cmdCarrierModel, cmdReadBlock, tc.entry)
			cat := loadCategoryOf(t, src, "cmd-inherited.toml")
			if cat != table.CatMalformedAccessorDeclaration {
				t.Errorf("category = %q; want the INHERITED %q — C1 relaxes only "+
					"the path rule", cat, table.CatMalformedAccessorDeclaration)
			}
		})
	}

	t.Run("a command write still requires read_back", func(t *testing.T) {
		src := swapEntry(t, cmdCarrierModel, `[write.state]
role = "state"
command = ["tools/flowstate-write", "{artifact}"]
keys = ["status"]
timeout = "2s"
read_back = true`, `[write.state]
role = "state"
command = ["tools/flowstate-write", "{artifact}"]
keys = ["status"]
timeout = "2s"`)
		cat := loadCategoryOf(t, src, "cmd-no-readback.toml")
		if cat != table.CatMalformedAccessorDeclaration {
			t.Errorf("category = %q; want %q", cat, table.CatMalformedAccessorDeclaration)
		}
	})
}

// REQ-87: "the loader decodes strictly (`source.go::decodeStrict` sets
// `DisallowUnknownFields`), so a field absent from the struct is a
// **document-level** `CatUnknownSchemaField`, not one of C5's per-entry
// categories."
// ADVERSARIAL
//
// This is the precondition of C5 reporting anything: if the six fields do
// NOT decode, every mutant below refuses as an unknown schema field and
// the per-entry categories are never reached.
func TestReq87_AnUndeclaredEntryFieldIsADocumentLevelUnknownSchemaField(t *testing.T) {
	// The discriminating PAIR. "`stdin` refuses as an unknown field" is
	// satisfied vacuously while NO carrier field decodes, so the arm that
	// carries the contract is the other half: the six declared fields must
	// NOT refuse that way. Asserted together, in one test, so neither can
	// pass alone.
	t.Run("a charted-but-unbuilt field is an unknown schema field", func(t *testing.T) {
		src := swapEntry(t, cmdCarrierModel, cmdReadBlock, `[read.state]
role = "state"
command = ["reader"]
stdin = "none"
keys = ["status"]
timeout = "2s"`)

		cat := loadCategoryOf(t, src, "cmd-unknown-field.toml")
		if cat != table.CatUnknownSchemaField {
			t.Errorf("category = %q; want %q — `stdin` is charted, not built "+
				"(REQ-136), so it decodes as an unknown field", cat,
				table.CatUnknownSchemaField)
		}
	})

	t.Run("the six declared carrier fields are not unknown fields", func(t *testing.T) {
		// This is C5's precondition: until the six DECODE, every mutant
		// below refuses `unknown_schema_field` and the per-entry
		// categories are unreachable.
		_, err := table.Load([]byte(cmdCarrierModel), "cmd-carrier.toml")
		if err == nil {
			return
		}
		if cat, _ := table.CategoryOf(err); cat == table.CatUnknownSchemaField {
			t.Fatalf("a model declaring the six carrier fields refused %q; the "+
				"fields must decode onto the entry struct before any C5 "+
				"category is reachable: %v", table.CatUnknownSchemaField, err)
		}
		t.Fatalf("the command-carrier fixture must load clean: %v", err)
	})
}

// --- C2's load-time arm --------------------------------------------------

// REQ-18: "A placeholder is recognized only as a whole argv element. An
// element that contains a `{...}` token without being exactly a known
// placeholder is a load-time defect (C5) — never silently-literal text."
// REQ-73: "command_unknown_placeholder   # unknown or non-whole-element
// {…} token"
// REQ-142: "The placeholder vocabulary is versioned as a closed set: v1 is
// `{artifact}` and is complete (C2)."
// ADVERSARIAL
func TestReq18_PlaceholdersAreWholeElementAndDrawnFromTheClosedV1Vocabulary(t *testing.T) {
	cases := []struct {
		name  string
		entry string
	}{
		{"unknown placeholder token", `[read.state]
role = "state"
command = ["reader", "{role}"]
keys = ["status"]
timeout = "2s"`},
		{"known token embedded mid-string", `[read.state]
role = "state"
command = ["reader", "--file={artifact}"]
keys = ["status"]
timeout = "2s"`},
		{"known token with a suffix", `[read.state]
role = "state"
command = ["reader", "{artifact}.bak"]
keys = ["status"]
timeout = "2s"`},
		{"partial brace token", `[read.state]
role = "state"
command = ["reader", "{artifact"]
keys = ["status"]
timeout = "2s"`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := swapEntry(t, cmdCarrierModel, cmdReadBlock, tc.entry)
			cat := loadCategoryOf(t, src, "cmd-placeholder.toml")
			if cat != table.CatCommandUnknownPlaceholder {
				t.Errorf("category = %q; want %q — a non-whole-element token is "+
					"never silently-literal text", cat,
					table.CatCommandUnknownPlaceholder)
			}
		})
	}

	t.Run("the v1 vocabulary member loads whole-element", func(t *testing.T) {
		if _, err := table.Load([]byte(cmdCarrierModel), "cmd-carrier.toml"); err != nil {
			t.Fatalf("a whole-element {artifact} must load clean: %v", err)
		}
	})
}

// --- C5: the interpreter deny-list ---------------------------------------

// REQ-74: "command_shell_interpreter     # argv0 + inline-code flag (sh -c,
// bash -c, python -c, env chains); no opt-in in v1"
// REQ-86: "the defect's report says so (\"inline shell is not a declared
// command; put it in a script and declare the script as argv0\")"
// ADVERSARIAL
func TestReq74_KnownInterpreterWithInlineCodeIsALoadTimeDefect(t *testing.T) {
	cases := []struct {
		name  string
		entry string
	}{
		{"sh -c", `[read.state]
role = "state"
command = ["sh", "-c", "cat {artifact}"]
keys = ["status"]
timeout = "2s"`},
		{"bash -c", `[read.state]
role = "state"
command = ["bash", "-c", "echo hi"]
keys = ["status"]
timeout = "2s"`},
		{"python -c", `[read.state]
role = "state"
command = ["python", "-c", "print(1)"]
keys = ["status"]
timeout = "2s"`},
		{"env chain to sh -c", `[read.state]
role = "state"
command = ["env", "sh", "-c", "echo hi"]
keys = ["status"]
timeout = "2s"`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := swapEntry(t, cmdCarrierModel, cmdReadBlock, tc.entry)
			_, err := table.Load([]byte(src), "cmd-interpreter.toml")
			if err == nil {
				t.Fatalf("%s loaded clean; a fixed interpreter vector is still a shell runner", tc.name)
			}
			cat, ok := table.CategoryOf(err)
			if !ok {
				t.Fatalf("refusal carried no category: %v", err)
			}
			if cat != table.CatCommandShellInterpreter {
				t.Fatalf("category = %q; want %q", cat, table.CatCommandShellInterpreter)
			}
			// REQ-86: the remediation travels IN the refusal, the
			// remediation-in-the-refusal pattern C5 names. Asserted on the
			// carried detail, which is data, not a rendered log line.
			f := cmdFailureOf(t, err)
			if !strings.Contains(f.Detail, "script") {
				t.Errorf("detail = %q; want it to name the remediation "+
					"(\"put it in a script and declare the script as argv0\")", f.Detail)
			}
		})
	}
}

// REQ-83: "interpreter set: OPEN (deny-listed, not closed) — an unlisted
// interpreter is admitted, so the list grows by amendment"
// REQ-122 (S1's open-deny-list negative): "an entry whose argv0 is an
// unlisted interpreter (`perl -e`) **loads**"
// DOMAIN EDGE
//
// NEGATIVE REQ. Without this arm the suite reads as "inline shell is
// impossible", which is the guarantee C5 explicitly declines to make.
func TestReq83_TheInterpreterSetIsDenyListedNotClosed(t *testing.T) {
	src := swapEntry(t, cmdCarrierModel, cmdReadBlock, `[read.state]
role = "state"
command = ["perl", "-e", "print 1"]
keys = ["status"]
timeout = "2s"`)

	if _, err := table.Load([]byte(src), "cmd-perl.toml"); err != nil {
		t.Errorf("an UNLISTED interpreter must load — C5's deny-list is open and "+
			"raises the cost of an inline-shell carrier without claiming to make "+
			"one impossible; refused: %v", err)
	}
}

// --- C5: output shape ----------------------------------------------------

// REQ-75: "command_output_shape          # output = \"raw\" with keys ≠ 1;
// exit_absent on a non-read; exit_verdicts on a non-gate or naming a
// non-verdict"
// REQ-32: "valid only when `keys` has exactly one entry"
// ADVERSARIAL
//
// PHASE-0 READING (ASSUMPTION REQ-3/REQ-4): the capability test is the
// discriminator, so `exit_absent` on a gate and `exit_verdicts` on a read
// are both this category.
func TestReq75_OutputShapeArmsCoverRawArityAndOffCapabilityExitMaps(t *testing.T) {
	writeBlock := `[write.state]
role = "state"
command = ["tools/flowstate-write", "{artifact}"]
keys = ["status"]
timeout = "2s"
read_back = true`
	gateBlock := `[gate.approval]
role = "state"
command = ["git", "-C", "{artifact}", "diff", "--quiet"]
exit_verdicts = { "0" = "allow", "1" = "deny" }
keys = ["status"]
timeout = "2s"`

	cases := []struct {
		name    string
		old     string
		replace string
	}{
		{"raw output with two declared keys", cmdReadBlock, `[read.state]
role = "state"
command = ["reader"]
output = "raw"
keys = ["status", "recognized"]
timeout = "2s"`},
		{"raw output with zero-arity is caught by arity not by keys", cmdReadBlock, `[read.state]
role = "state"
command = ["reader"]
output = "raw"
keys = ["status", "status"]
timeout = "2s"`},
		{"exit_absent on a gate entry", gateBlock, `[gate.approval]
role = "state"
command = ["gater"]
exit_absent = [1]
keys = ["status"]
timeout = "2s"`},
		{"exit_absent on a write entry", writeBlock, `[write.state]
role = "state"
command = ["writer"]
exit_absent = [1]
keys = ["status"]
timeout = "2s"
read_back = true`},
		{"exit_verdicts on a read entry", cmdReadBlock, `[read.state]
role = "state"
command = ["reader"]
exit_verdicts = { "0" = "allow" }
keys = ["status"]
timeout = "2s"`},
		{"exit_verdicts naming a non-verdict", gateBlock, `[gate.approval]
role = "state"
command = ["gater"]
exit_verdicts = { "0" = "maybe" }
keys = ["status"]
timeout = "2s"`},
		{"output on a gate entry", gateBlock, `[gate.approval]
role = "state"
command = ["gater"]
output = "raw"
keys = ["status"]
timeout = "2s"`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := swapEntry(t, cmdCarrierModel, tc.old, tc.replace)
			cat := loadCategoryOf(t, src, "cmd-output-shape.toml")
			if cat != table.CatCommandOutputShape {
				t.Errorf("category = %q; want %q", cat, table.CatCommandOutputShape)
			}
		})
	}
}

// REQ-39: "the gate `verdict` strings are `internal/accessor/model.go::
// Verdict`'s (`0004:C9`)"
// DOMAIN EDGE
//
// The oracle is the ACCESSOR package's closed verdict set, so a loader
// that grew a private fourth verdict string would fail here even while
// every literal above still passed.
func TestReq39_ExitVerdictsAdmitExactlyTheThreeAccessorVerdicts(t *testing.T) {
	gateBlock := `[gate.approval]
role = "state"
command = ["git", "-C", "{artifact}", "diff", "--quiet"]
exit_verdicts = { "0" = "allow", "1" = "deny" }
keys = ["status"]
timeout = "2s"`

	for _, v := range accessor.Verdicts() {
		t.Run(string(v), func(t *testing.T) {
			src := swapEntry(t, cmdCarrierModel, gateBlock, `[gate.approval]
role = "state"
command = ["gater"]
exit_verdicts = { "7" = "`+string(v)+`" }
keys = ["status"]
timeout = "2s"`)
			if _, err := table.Load([]byte(src), "cmd-verdict.toml"); err != nil {
				t.Errorf("verdict %q must be admissible in exit_verdicts; refused: %v", v, err)
			}
		})
	}
}

// --- C5: env conflict ----------------------------------------------------

// REQ-54: "An `env` key or `env_pass` name matching the `INTRASTATE_`
// prefix is a C5 `command_env_conflict` defect at load, so the overlay is
// never shadowed silently"
// REQ-76: "command_env_conflict          # an `env` key or `env_pass` name
// matching the reserved `INTRASTATE_` prefix (C4 overlay)"
// ADVERSARIAL
//
// PHASE-0 READING (ASSUMPTION REQ-54): the prefix check is case-sensitive
// and literal, matching the `LC_` rule's stated literalness — so
// `intrastate_role` is NOT a defect.
func TestReq54_TheIntrastatePrefixIsReservedFromEnvAndEnvPass(t *testing.T) {
	base := `[write.state]
role = "state"
command = ["tools/flowstate-write", "{artifact}"]
keys = ["status"]
timeout = "2s"
read_back = true`

	refused := []struct {
		name  string
		entry string
	}{
		{"env key shadowing INTRASTATE_ROLE", `[write.state]
role = "state"
command = ["writer"]
env = { INTRASTATE_ROLE = "spoof" }
keys = ["status"]
timeout = "2s"
read_back = true`},
		{"env key with the bare prefix", `[write.state]
role = "state"
command = ["writer"]
env = { INTRASTATE_ANYTHING = "x" }
keys = ["status"]
timeout = "2s"
read_back = true`},
		{"env_pass naming a reserved var", `[write.state]
role = "state"
command = ["writer"]
env_pass = ["INTRASTATE_PROTOCOL"]
keys = ["status"]
timeout = "2s"
read_back = true`},
	}

	for _, tc := range refused {
		t.Run(tc.name, func(t *testing.T) {
			src := swapEntry(t, cmdCarrierModel, base, tc.entry)
			cat := loadCategoryOf(t, src, "cmd-env-conflict.toml")
			if cat != table.CatCommandEnvConflict {
				t.Errorf("category = %q; want %q", cat, table.CatCommandEnvConflict)
			}
		})
	}

	t.Run("a lowercase near-miss is not the reserved prefix", func(t *testing.T) {
		// PHASE-0 READING (ASSUMPTION REQ-54): literal, case-sensitive.
		src := swapEntry(t, cmdCarrierModel, base, `[write.state]
role = "state"
command = ["writer"]
env = { intrastate_role = "fine" }
keys = ["status"]
timeout = "2s"
read_back = true`)
		if _, err := table.Load([]byte(src), "cmd-env-lower.toml"); err != nil {
			t.Errorf("a lowercase key is not the reserved INTRASTATE_ prefix; refused: %v", err)
		}
	})
}

// --- C5: registration ----------------------------------------------------

// REQ-77: "registration: all six are appended to `table.Categories()`,
// whose hand-maintained list IS the closed set; appended in the order
// declared above, after the existing members, so a consumer enumerating
// the list sees additions only at the tail."
// REQ-85: "A `Category` constant is not in the closed set until it is
// appended to `table.Categories()` ... so S1 must assert membership, not
// merely the refusal."
// BOUNDARY
//
// PHASE-0 READING (ASSUMPTION REQ-77): "in clause order" is the C5 block's
// declared order — conflict, empty, unknown_placeholder, shell_interpreter,
// output_shape, env_conflict.
func TestReq77_TheSixCommandCategoriesAreRegisteredAtTheTailInClauseOrder(t *testing.T) {
	want := []table.Category{
		table.CatCommandAndPathConflict,
		table.CatCommandEmpty,
		table.CatCommandUnknownPlaceholder,
		table.CatCommandShellInterpreter,
		table.CatCommandOutputShape,
		table.CatCommandEnvConflict,
	}

	all := table.Categories()
	if len(all) < len(want) {
		t.Fatalf("Categories() has %d members; cannot carry the six additions", len(all))
	}
	tail := all[len(all)-len(want):]
	if !slices.Equal(tail, want) {
		t.Errorf("Categories() tail = %#v; want the six C5 categories appended in "+
			"clause order %#v — a consumer enumerating the list must see "+
			"additions only at the tail", tail, want)
	}
}

// REQ-78: "Each gets a typed `Category` constant of the existing `Cat…`
// form beside the others; the wire STRINGS above are the contract and the
// constant identifiers are not."
// BOUNDARY
func TestReq78_TheSixCategoriesCarryTheirContractWireStrings(t *testing.T) {
	// The oracle reads the REGISTERED set by wire string, not the Go
	// constants: C5 fixes the strings as the contract and explicitly
	// leaves the identifiers unconstrained, so a test that compared
	// constants to literals would pin the half the clause does not fix
	// and pass while the category was invisible to every consumer.
	registered := map[string]bool{}
	for _, c := range table.Categories() {
		registered[string(c)] = true
	}

	for _, want := range []string{
		"command_and_path_conflict",
		"command_empty",
		"command_unknown_placeholder",
		"command_shell_interpreter",
		"command_output_shape",
		"command_env_conflict",
	} {
		if !registered[want] {
			t.Errorf("Categories() carries no member with the wire string %q; "+
				"the STRING is the contract and registration is what puts it "+
				"in the closed set", want)
		}
	}
}

// REQ-79: "The list's total size is not a contract at any point — it is
// append-only, and no clause or test may assert a count"
// REQ-81: "ACROSS entries in the SAME table there is no order ... and no
// test may assert it."
// DOMAIN EDGE — NEGATIVE REQ.
//
// These two forbid assertions rather than requiring them, so the test that
// would catch a violation is a test that the PROPERTY they protect holds:
// the set is duplicate-free and every registered category is reachable by
// value, which is what a consumer needs from an append-only list. A count
// assertion anywhere else in the suite is caught by review, not by code —
// recorded in coverage.md.
func TestReq79_TheCategorySetIsAppendOnlyAndDuplicateFree(t *testing.T) {
	all := table.Categories()

	// Positive precondition: the six additions must be present, or the
	// append-only property below is asserted over a set this RDR never
	// touched and the test proves nothing about the change.
	six := []table.Category{
		table.CatCommandAndPathConflict,
		table.CatCommandEmpty,
		table.CatCommandUnknownPlaceholder,
		table.CatCommandShellInterpreter,
		table.CatCommandOutputShape,
		table.CatCommandEnvConflict,
	}
	for _, c := range six {
		if !slices.Contains(all, c) {
			t.Fatalf("Categories() does not carry %q; the append-only assertion "+
				"below would range a set this RDR did not extend", c)
		}
	}

	// The prefix before the six additions must be UNCHANGED — that is what
	// "append-only" means for a consumer enumerating the list, and it is
	// the property a count assertion would falsely stand in for (REQ-79
	// forbids the count).
	head := all[:len(all)-len(six)]
	if slices.Contains(head, table.CatCommandAndPathConflict) {
		t.Error("a command category is registered before the tail; additions " +
			"must be visible only at the tail")
	}

	seen := map[table.Category]bool{}
	for _, c := range all {
		if seen[c] {
			t.Errorf("category %q is registered twice; the list is append-only "+
				"and a duplicate makes membership ambiguous", c)
		}
		seen[c] = true
		if c == "" {
			t.Error("Categories() carries an empty category")
		}
	}
}

// REQ-80: "precedence: WITHIN one entry — load is fail-fast (`0002:C24` —
// one categorized error for the whole document, never a list), so an entry
// carrying several of these defects reports the first in the order declared
// above."
// BOUNDARY
func TestReq80_WithinOneEntryTheEarlierClauseWins(t *testing.T) {
	cases := []struct {
		name  string
		entry string
		want  table.Category
	}{
		{
			// both carriers AND an empty element: conflict is clause 1.
			name: "conflict outranks empty",
			entry: `[read.state]
role = "state"
path = "flow.status"
command = ["reader", ""]
keys = ["status"]
timeout = "2s"`,
			want: table.CatCommandAndPathConflict,
		},
		{
			// empty element AND an unknown placeholder: empty is clause 2.
			name: "empty outranks unknown placeholder",
			entry: `[read.state]
role = "state"
command = ["reader", "", "{role}"]
keys = ["status"]
timeout = "2s"`,
			want: table.CatCommandEmpty,
		},
		{
			// unknown placeholder AND sh -c: placeholder is clause 3.
			name: "unknown placeholder outranks shell interpreter",
			entry: `[read.state]
role = "state"
command = ["sh", "-c", "{role}"]
keys = ["status"]
timeout = "2s"`,
			want: table.CatCommandUnknownPlaceholder,
		},
		{
			// sh -c AND raw-with-two-keys: interpreter is clause 4.
			name: "shell interpreter outranks output shape",
			entry: `[read.state]
role = "state"
command = ["sh", "-c", "cat"]
output = "raw"
keys = ["status", "recognized"]
timeout = "2s"`,
			want: table.CatCommandShellInterpreter,
		},
		{
			// raw-with-two-keys AND a reserved env key: shape is clause 5.
			name: "output shape outranks env conflict",
			entry: `[read.state]
role = "state"
command = ["reader"]
output = "raw"
env = { INTRASTATE_ROLE = "x" }
keys = ["status", "recognized"]
timeout = "2s"`,
			want: table.CatCommandOutputShape,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := swapEntry(t, cmdCarrierModel, cmdReadBlock, tc.entry)
			if cat := loadCategoryOf(t, src, "cmd-precedence.toml"); cat != tc.want {
				t.Errorf("category = %q; want %q — fail-fast reports the FIRST "+
					"clause in declared order", cat, tc.want)
			}
		})
	}
}

// REQ-82: "ACROSS tables the order IS fixed — `load.go::loadAccessors` runs
// read, then write, then gate, returning on the first error — so a
// defective read entry always masks a defective gate entry; that ordering
// is `0002`'s, inherited not established here"
// DOMAIN EDGE
func TestReq82_ADefectiveReadEntryMasksADefectiveGateEntry(t *testing.T) {
	src := swapEntry(t, cmdCarrierModel, cmdReadBlock, `[read.state]
role = "state"
command = []
keys = ["status"]
timeout = "2s"`)
	src = swapEntry(t, src, `[gate.approval]
role = "state"
command = ["git", "-C", "{artifact}", "diff", "--quiet"]
exit_verdicts = { "0" = "allow", "1" = "deny" }
keys = ["status"]
timeout = "2s"`, `[gate.approval]
role = "state"
command = ["gater", "{role}"]
keys = ["status"]
timeout = "2s"`)

	cat := loadCategoryOf(t, src, "cmd-table-order.toml")
	if cat != table.CatCommandEmpty {
		t.Errorf("category = %q; want the READ table's %q — loadAccessors runs "+
			"read, write, gate and returns on the first error (0002's contract)",
			cat, table.CatCommandEmpty)
	}
}

// REQ-84: "these are `internal/table/category.go::Category` constants ...
// not `accessor.ValidationCode`, whose eight-member closure is RDR 0004's
// own test contract."
// REQ-122 (S1): "`accessor.ValidationCodes` stays at eight (REQ-84)."
// ADVERSARIAL — NEGATIVE REQ.
//
// The mutant this kills is an implementation that homes the six defects on
// `accessor.ValidationCode` instead, which would break RDR 0004's REQ-84
// closure. Asserting HERE, in this RDR's own suite, makes that breakage
// this record's failure rather than a mysterious 0004 regression.
func TestReq84_TheAccessorValidationCodeSetIsNotExtendedByThisRDR(t *testing.T) {
	codes := accessor.ValidationCodes()
	for _, c := range codes {
		if strings.HasPrefix(string(c), "command_") {
			t.Errorf("accessor.ValidationCodes() carries %q; C5's defects home on "+
				"table.Category, and extending this set breaks 0004's REQ-84 "+
				"closure", c)
		}
	}
	// Positive precondition: the set must still be the 0004 eight, or the
	// loop above is vacuous against a set someone emptied.
	if len(codes) == 0 {
		t.Fatal("accessor.ValidationCodes() is empty; the closure assertion is vacuous")
	}

	// And the discriminating half: the six defects must live SOMEWHERE.
	// Without this, "not on ValidationCode" is satisfied by an
	// implementation that never built them at all, and the clause's
	// content — which family owns them — goes unasserted.
	registered := table.Categories()
	for _, c := range []table.Category{
		table.CatCommandAndPathConflict,
		table.CatCommandEmpty,
		table.CatCommandUnknownPlaceholder,
		table.CatCommandShellInterpreter,
		table.CatCommandOutputShape,
		table.CatCommandEnvConflict,
	} {
		if !slices.Contains(registered, c) {
			t.Errorf("%q is on neither accessor.ValidationCodes() nor "+
				"table.Categories(); C5 homes the six on table.Category", c)
		}
	}
}

// REQ-132: "load-time C5 defects name the entry and the offending element"
// DOMAIN EDGE
func TestReq132_ALoadTimeCommandDefectNamesTheEntryAndTheOffendingElement(t *testing.T) {
	src := swapEntry(t, cmdCarrierModel, cmdReadBlock, `[read.state]
role = "state"
command = ["reader", "{role}"]
keys = ["status"]
timeout = "2s"`)

	_, err := table.Load([]byte(src), "cmd-detail.toml")
	if err == nil {
		t.Fatal("the mutant loaded clean")
	}
	f := cmdFailureOf(t, err)
	if !strings.Contains(f.Detail, "state") {
		t.Errorf("detail = %q; want it to name the entry `state`", f.Detail)
	}
	if !strings.Contains(f.Detail, "{role}") {
		t.Errorf("detail = %q; want it to name the offending element `{role}`", f.Detail)
	}
}

// REQ-115: "Phase 1: Carrier — Extend `internal/table` (`sourceAcc`,
// `Accessor`, `accessorTable`) with `command` and the C1/C5 load rules."
// REQ-12: "The six fields above are the complete set this RDR adds"
// HAPPY PATH
//
// The oracle is that all six fields survive the loader onto the normalized
// `table.Accessor`, which is what every downstream binding reads. A carrier
// that decoded them and dropped them on normalization would pass every
// refusal test above and carry no authority at runtime.
func TestReq12_AllSixCarrierFieldsSurviveOntoTheNormalizedAccessor(t *testing.T) {
	src := swapEntry(t, cmdCarrierModel, cmdReadBlock, `[read.state]
role = "state"
command = ["reader", "{artifact}"]
output = "raw"
exit_absent = [1, 7]
env = { TOOL_MODE = "strict" }
env_pass = ["SSH_AUTH_SOCK"]
keys = ["status"]
timeout = "2s"`)

	m, err := table.Load([]byte(src), "cmd-all-fields.toml")
	if err != nil {
		t.Fatalf("an entry declaring every carrier field must load clean: %v", err)
	}
	acc := m.Readers["state"]

	if !slices.Equal(acc.Command, []string{"reader", "{artifact}"}) {
		t.Errorf("Command = %#v; want the declared vector", acc.Command)
	}
	if acc.Output == nil || *acc.Output != "raw" {
		t.Errorf("Output = %v; want a non-nil \"raw\"", acc.Output)
	}
	if !slices.Equal(acc.ExitAbsent, []int{1, 7}) {
		t.Errorf("ExitAbsent = %#v; want []int{1, 7}", acc.ExitAbsent)
	}
	if acc.Env["TOOL_MODE"] != "strict" {
		t.Errorf("Env = %#v; want TOOL_MODE=strict", acc.Env)
	}
	if !slices.Equal(acc.EnvPass, []string{"SSH_AUTH_SOCK"}) {
		t.Errorf("EnvPass = %#v; want [SSH_AUTH_SOCK]", acc.EnvPass)
	}
	// A command entry carries NO path: C1's exactly-one rule.
	if acc.Path != "" {
		t.Errorf("Path = %q on a command entry; want empty — exactly one carrier", acc.Path)
	}
}

// REQ-139: "Every existing model keeps working untouched, and no
// path-backed entry changes meaning."
// DOMAIN EDGE — NEGATIVE REQ.
//
// The mutant this kills is a carrier change that makes `path` optional in
// the wrong direction, or that alters what a path-backed entry normalizes
// to. The oracle is the shipped RDR 0002 fixture, unmodified.
func TestReq139_APathBackedModelIsUnchangedByTheCommandCarrier(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	acc, ok := m.Readers["rdr-status"]
	if !ok {
		t.Fatal("the promoted fixture no longer carries read.rdr-status")
	}
	if acc.Path != "rdr.status" {
		t.Errorf("Path = %q; want the declared locator %q unchanged", acc.Path, "rdr.status")
	}
	if len(acc.Command) != 0 {
		t.Errorf("Command = %#v on a path-backed entry; want empty", acc.Command)
	}
	if acc.Output != nil {
		t.Errorf("Output = %v on a path-backed entry; want nil", acc.Output)
	}

	// Positive precondition. "Nothing changed for path-backed entries" is
	// satisfied vacuously while the command carrier does not exist at all,
	// so the no-regression claim is only meaningful once a command entry
	// loads. Both halves in one test: the co-existence is the contract.
	if _, err := table.Load([]byte(cmdCarrierModel), "cmd-carrier.toml"); err != nil {
		t.Fatalf("the command carrier must be live for the no-regression claim "+
			"to say anything; a command-backed model refused: %v", err)
	}
}
