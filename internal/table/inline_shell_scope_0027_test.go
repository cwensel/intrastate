package table_test

// RDR 0027 — what the inline-shell check on command bindings actually
// promises. The PREDICATE half.
//
// `0027:C1` succeeds two lines of `0025:C5`: the `command_shell_interpreter`
// line and the `interpreter set` line. Nothing else in C5 moves, so this file
// adds probes beside the shipped 0025 ones and edits none of them
// (`0027:S1` "the four REQ-74 probes and `neg/neg-command-shell-interpreter.toml`
// unedited and still green").
//
// The coverage floor this file exists to hold has two mutants.
//
// The first is an ARGV0-ANCHORED predicate. Today `interpreterForm` walks a
// leading `env` chain and then requires the interpreter at that one position,
// so every wrapper — `nice`, `timeout 5`, `xargs`, `doas`, and `env` with any
// OPTION flag — lints green while executing a shell. A suite that probes only
// argv0 forms passes against that code, which is why every wrapper scenario
// below puts the interpreter at argv >= 1.
//
// The second is a SPELLINGS-SCOPED reading of C1's out-of-scope line. C1
// scopes the admitted stdin forms by CHANNEL — "any listed interpreter taking
// its code on stdin rather than as a later argv word is admitted, however
// spelled" — so `["python","-"]` is green even though `python` IS on the
// deny-list. An implementation that read the line as a list of shell
// spellings refuses it. `["python","-"]` and `["sh","-es"]` are the
// discriminating rows (`0027:S4`).
//
// Every scenario is driven through `table.Load` so the oracle is the CATEGORY
// a model author actually sees, never a direct call on the unexported
// predicate (`0027:S1` preamble). The fixture, `swapEntry`, `loadCategoryOf`
// and `cmdFailureOf` are the 0025 file's; reusing them is what keeps the two
// records' probes comparable.

import (
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// --- helpers -------------------------------------------------------------

// interpEntry renders a `[read.state]` block carrying the given argv. It is
// the one mutation every scenario below makes to the 0025 fixture, so a
// scenario's argv is the only thing that varies between a refusal and a green
// load.
func interpEntry(argv []string) string {
	quoted := make([]string, len(argv))
	for i, w := range argv {
		quoted[i] = `"` + strings.ReplaceAll(w, `"`, `\"`) + `"`
	}
	return `[read.state]
role = "state"
command = [` + strings.Join(quoted, ", ") + `]
keys = ["status"]
timeout = "2s"`
}

// loadWithCommand swaps argv into the 0025 command-carrier fixture's read
// entry and loads it, returning the model and the load error unexamined.
func loadWithCommand(t *testing.T, argv []string) (*table.Model, error) {
	t.Helper()

	src := swapEntry(t, cmdCarrierModel, cmdReadBlock, interpEntry(argv))
	return table.Load([]byte(src), "cmd-0027.toml")
}

// refusalFor loads argv and requires a categorized refusal, returning the
// failure DATA. It fails the test rather than returning an error, so a
// scenario reads as one assertion.
func refusalFor(t *testing.T, argv []string) *table.Failure {
	t.Helper()

	_, err := loadWithCommand(t, argv)
	if err == nil {
		t.Fatalf("%v loaded clean; want a categorized refusal", argv)
	}
	return cmdFailureOf(t, err)
}

// wantInterpreterRefusal is S1/S3's full oracle in one place: the category
// string (not "load failed"), the two matched words in the detail, and the
// unchanged script remediation.
func wantInterpreterRefusal(t *testing.T, argv []string, wantForm string) {
	t.Helper()

	f := refusalFor(t, argv)
	if f.Category != table.CatCommandShellInterpreter {
		t.Fatalf("%v refused as %q; want %q — the wrapper prefix must not "+
			"change which defect a wrapped interpreter reports",
			argv, f.Category, table.CatCommandShellInterpreter)
	}
	if !strings.Contains(f.Detail, wantForm) {
		t.Errorf("%v detail = %q; want it to name the two matched words %q — "+
			"a bare category assertion passes even on today's message "+
			"(`0027:§mini-check-oracle`, S5 row)", argv, f.Detail, wantForm)
	}
	if !strings.Contains(f.Detail, "put it in a script and declare the script as argv0") {
		t.Errorf("%v detail = %q; want the UNCHANGED script remediation "+
			"(`0027:C1` report: \"remediation … unchanged\")", argv, f.Detail)
	}
}

// wantAdmitted is S4's oracle: green AND the binding loads with its argv
// unchanged. Green alone is an absence-of-error assertion a no-op predicate
// also satisfies (`0027:§mini-check-oracle`, S4 row).
func wantAdmitted(t *testing.T, argv []string) {
	t.Helper()

	m, err := loadWithCommand(t, argv)
	if err != nil {
		t.Fatalf("%v must lint GREEN — C1 names it out of scope BY NAME, and a "+
			"predicate change that starts refusing it has broken the promise, "+
			"not tightened it (`0027:S4`); refused: %v", argv, err)
	}
	acc, ok := m.Readers["state"]
	if !ok {
		t.Fatalf("the loaded model carries no `read.state` accessor")
	}
	if !slices.Equal(acc.Command, argv) {
		t.Errorf("loaded argv = %#v; want %#v unchanged — \"loads with its argv "+
			"unchanged\", not merely that no error was returned", acc.Command, argv)
	}
}

// --- C1: the position-free predicate -------------------------------------

// REQ-1: "command_shell_interpreter     # a listed interpreter word followed,
// at ANY later argv position, by one of that interpreter's inline-code flags —
// under any prefix (env and its options, nice, timeout, xargs, doas, …); no
// opt-in in v1"
// REQ-42 / `0027:S1`: "One mutant per wrapper the Problem Statement names —
// `env -i`, `env -u FOO`, `nice`, `timeout 5`, `xargs`, `nohup`, `setsid`,
// `stdbuf -o0`, `chpst`, `doas` — each wrapping `sh -c`. **Expected**: each
// refuses `command_shell_interpreter`, detail naming the matched words `sh -c`
// and the script remediation … the wrapper binary need not exist on the test
// host, since the predicate reads argv words only."
// REQ-51: "interpreter word + its inline-code flag at any later argv position
// (`[\"nice\",\"sh\",\"-c\",\"…\"]`) | refuse — load-time defect |
// `command_shell_interpreter`, detail names the two matched words + the script
// remediation | none (load fails) | **loud**"
// REQ-61: "Visible: a refused wrapper form names `sh -c` (the matched words)
// and the wrapper-file remediation"
// ADVERSARIAL
//
// Ten wrappers, each defeating today's argv0 anchor by one leading word. Four
// of them (`setsid`, `timeout`, `chpst`, `doas`) do not exist on stock darwin;
// that is deliberate and costs nothing, because the predicate reads argv WORDS
// and never the resolved binary (`0027:C1` reads:).
func TestReq42_EveryNamedWrapperStillRefusesTheInterpreterForm(t *testing.T) {
	cases := []struct {
		name string
		argv []string
	}{
		{"env -i", []string{"env", "-i", "sh", "-c", "echo hi"}},
		{"env -u FOO", []string{"env", "-u", "FOO", "sh", "-c", "echo hi"}},
		{"nice", []string{"nice", "sh", "-c", "echo hi"}},
		{"timeout 5", []string{"timeout", "5", "sh", "-c", "echo hi"}},
		{"xargs", []string{"xargs", "sh", "-c", "echo hi"}},
		{"nohup", []string{"nohup", "sh", "-c", "echo hi"}},
		{"setsid", []string{"setsid", "sh", "-c", "echo hi"}},
		{"stdbuf -o0", []string{"stdbuf", "-o0", "sh", "-c", "echo hi"}},
		{"chpst", []string{"chpst", "sh", "-c", "echo hi"}},
		{"doas", []string{"doas", "sh", "-c", "echo hi"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wantInterpreterRefusal(t, tc.argv, "sh -c")
		})
	}
}

// REQ-4: "Nothing before argv[i] is read"
// REQ-5: "no wrapper table exists and none may be added"
// REQ-68 / `0027:ALT2` rejected: "the position-free scan achieves the class
// closure with zero grammar and zero table"
// ADVERSARIAL — NEGATIVE REQ.
//
// The forward-binding arm of REQ-5: a wrapper nobody enumerated must refuse on
// the same predicate, because the predicate never consults a list of wrappers.
// An implementation that shipped a wrapper table would pass the ten named rows
// above and fail every row here — which is exactly the failure REQ-5 forbids.
func TestReq5_AnUnenumeratedPrefixIsRefusedBecauseNothingBeforeTheInterpreterIsRead(t *testing.T) {
	cases := []struct {
		name string
		argv []string
	}{
		{"never-named wrapper", []string{"some-wrapper-nobody-listed", "sh", "-c", "echo hi"}},
		{"two stacked wrappers", []string{"nice", "timeout", "5", "sh", "-c", "echo hi"}},
		{"wrapper with its own long flags", []string{"runner", "--jobs=4", "--verbose", "sh", "-c", "echo hi"}},
		{"deep prefix", []string{"a", "b", "c", "d", "e", "f", "sh", "-c", "echo hi"}},
		{"absolute-path wrapper", []string{"/opt/vendor/bin/wrap", "sh", "-c", "echo hi"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wantInterpreterRefusal(t, tc.argv, "sh -c")
		})
	}
}

// REQ-44 / `0027:S3`: "The `env`-option forms intrastate#q2q1 enumerates —
// `-i`, `-u FOO`, `--unset=FOO`, `-0`, `-C DIR`, `--chdir=DIR`, `--`, and a
// mixed `-i A=1 -- sh -c` chain. **Expected**: all refuse
// `command_shell_interpreter`, **detail naming the matched words `sh -c`** as
// in scenario 1 … Proves the deleted `env` walk is subsumed rather than merely
// removed"
// REQ-31: "The `env`-chain walk in `internal/table/load.go::interpreterForm`
// is deleted: this record DECIDES that C1 subsumes it"
// ADVERSARIAL
//
// Subsumption, not removal. Today's walk stops at the first non-`NAME=VALUE`
// token, so every option form below lints green; under C1 they refuse for the
// same reason `nice` does — nothing before the interpreter word is read.
func TestReq44_TheDeletedEnvWalkIsSubsumedByThePositionFreeScan(t *testing.T) {
	cases := []struct {
		name string
		argv []string
	}{
		{"-i", []string{"env", "-i", "sh", "-c", "echo hi"}},
		{"-u FOO", []string{"env", "-u", "FOO", "sh", "-c", "echo hi"}},
		{"--unset=FOO", []string{"env", "--unset=FOO", "sh", "-c", "echo hi"}},
		{"-0", []string{"env", "-0", "sh", "-c", "echo hi"}},
		{"-C DIR", []string{"env", "-C", "/tmp", "sh", "-c", "echo hi"}},
		{"--chdir=DIR", []string{"env", "--chdir=/tmp", "sh", "-c", "echo hi"}},
		{"--", []string{"env", "--", "sh", "-c", "echo hi"}},
		{"mixed -i A=1 -- sh -c", []string{"env", "-i", "A=1", "--", "sh", "-c", "echo hi"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wantInterpreterRefusal(t, tc.argv, "sh -c")
		})
	}
}

// REQ-2: "predicate:  refuse iff there exist i < j with base(argv[i]) a listed
// interpreter and argv[j] one of its listed flags"
// REQ-25: "when several pairs qualify, the lowest interpreter index and then
// the lowest flag index are reported; order is part of the predicate (a flag
// word before the interpreter word never matches)"
// BOUNDARY
//
// The `i < j` ordering is the SOLE bound on the false-refusal class
// position-freedom opens (`0027:§key-discoveries`), so it is asserted in both
// directions: a flag BEFORE the interpreter word is green, the same two words
// in the other order refuse.
func TestReq25_TheFlagMustFollowTheInterpreterWordNotPrecedeIt(t *testing.T) {
	t.Run("flag before interpreter is green", func(t *testing.T) {
		for _, argv := range [][]string{
			{"grep", "-c", "pattern", "sh"},
			{"tool", "-e", "expr", "ruby"},
			{"setup-tests", "bash"},
		} {
			wantAdmitted(t, argv)
		}
	})

	t.Run("the same two words in i<j order refuse", func(t *testing.T) {
		wantInterpreterRefusal(t, []string{"grep", "sh", "pattern", "-c"}, "sh -c")
	})
}

// REQ-6: "The reported form is argv[i] + \" \" + argv[j] for the lowest i,
// then the lowest j"
// REQ-27: "Two listed interpreters in one argv resolve by the same rule and
// the outer word wins: `[\"python\",\"sh\",\"-c\",\"echo\"]` reports `python
// -c`, not the `sh -c` a reader's eye goes to, because `python` and `sh` both
// carry `-c` (`shellInterpreters`) and `python` holds the lower i. Refusal is
// unaffected — either pair refuses — so this is message attribution only"
// REQ-57: "the reported form cannot depend on map iteration order over
// `internal/table/load.go::shellInterpreters`"
// DOMAIN EDGE
//
// The record's own worked case, asserted as written. The counter-assertion is
// what makes it discriminating: an implementation that iterated the
// interpreter MAP rather than argv POSITIONS would report `sh -c` here on some
// runs and `python -c` on others, and a test that only asked "does it refuse"
// would never see it.
func TestReq27_TheOuterInterpreterWordWinsTheReportedForm(t *testing.T) {
	f := refusalFor(t, []string{"python", "sh", "-c", "echo"})
	if f.Category != table.CatCommandShellInterpreter {
		t.Fatalf("category = %q; want %q", f.Category, table.CatCommandShellInterpreter)
	}
	if !strings.Contains(f.Detail, "python -c") {
		t.Errorf("detail = %q; want it to report `python -c` — the LOWEST i wins, "+
			"and `python` carries `-c` too (`0027:D-selection-predicate`)", f.Detail)
	}
	if strings.Contains(f.Detail, "sh -c") {
		t.Errorf("detail = %q; reports `sh -c`, the pair a reader's eye goes to. "+
			"The tie-break is argv POSITION, not scanning order, precisely so "+
			"the reported form cannot depend on map iteration order (REQ-57)",
			f.Detail)
	}
}

// REQ-6 tie-break, second axis: "for the lowest i, then the lowest j"
// REQ-23: "the tie-break (lowest `i`, then lowest `j`)"
// BOUNDARY
//
// With one interpreter and two qualifying flags, the LOWEST j is reported.
// `node` is the only listed name carrying two inline-code flags (`-e`,
// `--eval`), so it is the only vector that can exercise this axis at all.
func TestReq6_WithSeveralQualifyingFlagsTheLowestJIsReported(t *testing.T) {
	f := refusalFor(t, []string{"nice", "node", "--eval", "x", "-e", "y"})
	if !strings.Contains(f.Detail, "node --eval") {
		t.Errorf("detail = %q; want `node --eval` — the LOWEST j after the "+
			"interpreter word, not whichever flag the scan happened to reach",
			f.Detail)
	}
}

// REQ-3: "base = the text after the last `/`, matched exactly (no suffix or
// alias folding: `python3`, `nodejs`, `busybox` stay unlisted spellings under
// the OPEN rule)"
// REQ-56: "Basename matching is `base = the text after the last /`, matched
// **exactly** — no suffix folding, no alias folding, no case folding … the
// check does no Unicode normalization, no case mapping, and no
// locale-dependent comparison; a listed name is a byte-exact match on the
// post-slash segment."
// INPUT EDGE
//
// The basename rule is asserted on both sides: a path-spelled interpreter
// under a wrapper still matches, and every folding a future "helpful" change
// might add stays UNMATCHED. A folding rule would silently widen the deny-list
// without amending C1, which the OPEN-list clause reserves to the clause.
func TestReq3_BasenameIsThePostSlashSegmentMatchedByteExactly(t *testing.T) {
	t.Run("path-spelled interpreter under a wrapper matches on its basename", func(t *testing.T) {
		for _, tc := range []struct {
			argv []string
			form string
		}{
			{[]string{"nice", "/bin/sh", "-c", "echo hi"}, "/bin/sh -c"},
			{[]string{"env", "-i", "/usr/local/bin/bash", "-c", "echo hi"}, "/usr/local/bin/bash -c"},
			{[]string{"timeout", "5", "./sh", "-c", "echo hi"}, "./sh -c"},
		} {
			wantInterpreterRefusal(t, tc.argv, tc.form)
		}
	})

	t.Run("no suffix, alias, or case folding", func(t *testing.T) {
		for _, argv := range [][]string{
			{"nice", "python3", "-c", "print(1)"},
			{"nice", "nodejs", "-e", "1"},
			{"nice", "busybox", "sh", "-x"},
			{"nice", "SH", "-c", "echo hi"},
			{"nice", "Bash", "-c", "echo hi"},
			{"nice", "shell", "-c", "echo hi"},
			{"nice", "ssh", "-c", "echo hi"},
			{"nice", "/bin/sh-wrapper", "-c", "echo hi"},
		} {
			wantAdmitted(t, argv)
		}
	})
}

// REQ-7: "reads:      argv WORDS only. The check never splits a word on
// whitespace, never reads stdin, files, PATH, or the resolved binary"
// REQ-8: "out of scope, BY NAME (admitted by lint; an interpreter may run): a
// shell string carried in ONE word (`env -S \"sh -c …\"`, or a single `\"sh -c
// …\"` element handed to a tool that re-splits it)"
// REQ-52: "shell string inside ONE word (`[\"env\",\"-S\",\"sh -c echo\"]`) |
// green | none | binding loads; a shell may run | **silent by contract** … no
// lint event is minted"
// REQ-30: "Not modelled, on purpose: wrapper option grammars, `--`
// conventions, `env -S` splitting, PATH or binary identity."
// ADVERSARIAL — NEGATIVE REQ.
//
// The one-word form is where a "helpful" implementation is most tempted to
// split on whitespace and catch what it can see. C1 forbids it: a word is
// never split, so the form is admitted knowingly and the shipped description
// says so (`0027:S7`).
func TestReq8_AShellStringCarriedInOneWordIsAdmittedAndNeverSplit(t *testing.T) {
	for _, argv := range [][]string{
		{"env", "-S", "sh -c echo"},
		{"env", "-S", "sh -c 'cat file'"},
		{"some-tool", "sh -c echo hi"},
		{"sh -c echo hi"},
		{"nice", "bash -c echo"},
	} {
		wantAdmitted(t, argv)
	}
}

// REQ-9: "an interpreter that reads its script from STDIN (`sh -s`, bare `sh`,
// `sh -es`, `python -`, `node -`) — owned by the charted `stdin = \"none\" |
// \"envelope\"` successor, whose remit is what a command entry CONSUMES and so
// covers the non-shell spellings too"
// REQ-10: "The channel is the scope: any listed interpreter taking its code on
// stdin rather than as a later argv word is admitted, however spelled."
// REQ-45 / `0027:S4`: "The admitted forms — `[\"sh\",\"./gate.sh\"]`,
// `[\"env\",\"-S\",\"sh -c echo\"]`, `[\"sh\",\"-s\"]`, `[\"sh\"]`,
// `[\"sh\",\"-es\"]`, `[\"python\",\"-\"]`, `[\"node\",\"-\"]`. **Expected**:
// every one lints GREEN **and the binding loads with its argv unchanged** —
// not merely that no error was returned, which a no-op predicate also
// satisfies … `[\"python\",\"-\"]` and `[\"sh\",\"-es\"]` are the
// discriminating cases"
// REQ-46: "These are C1's out-of-scope line as a test: a future predicate
// change that starts refusing one of them has broken the promise, not
// tightened it."
// REQ-53: "interpreter reading its script from stdin … | green … **silent by
// contract** — named by CHANNEL, not by shell spelling"
// REQ-62: "Silent: a one-word shell string or a stdin-fed interpreter (`sh
// -s`, bare `sh`, `sh -es`, `python -`, `node -`) lints green — by contract,
// named in C1"
// DOMAIN EDGE — NEGATIVE REQ.
//
// S4's seven rows verbatim. `["python","-"]` is the CHANNEL probe: `python` is
// a `shellInterpreters` member, so any implementation reading C1's
// out-of-scope line as a list of SHELL SPELLINGS refuses it and fails here.
// `["sh","-es"]` is the second discriminator — a bundled short-option cluster
// that is not the listed `-c` flag.
func TestReq45_TheSevenAdmittedFormsLintGreenAndLoadWithArgvUnchanged(t *testing.T) {
	cases := []struct {
		name string
		argv []string
	}{
		{"sanctioned wrapper file", []string{"sh", "./gate.sh"}},
		{"one-word shell string", []string{"env", "-S", "sh -c echo"}},
		{"stdin -s", []string{"sh", "-s"}},
		{"bare sh", []string{"sh"}},
		{"bundled -es", []string{"sh", "-es"}},
		{"non-shell stdin spelling python -", []string{"python", "-"}},
		{"non-shell stdin spelling node -", []string{"node", "-"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wantAdmitted(t, tc.argv)
		})
	}
}

// REQ-10 (channel scope), widened: "any listed interpreter taking its code on
// stdin rather than as a later argv word is admitted, however spelled"
// DOMAIN EDGE — NEGATIVE REQ.
//
// The channel reading as a forward-binding assertion, over listed names the
// record does not itself enumerate. A spellings-scoped implementation that
// special-cased only `python -` and `node -` passes S4's seven rows and fails
// here.
func TestReq10_TheStdinChannelIsAdmittedForEveryListedSpelling(t *testing.T) {
	for _, argv := range [][]string{
		{"ruby", "-"},
		{"php", "-"},
		{"bash", "-s"},
		{"zsh", "-s"},
		{"nice", "python", "-"},
		{"env", "-i", "sh", "-s"},
	} {
		wantAdmitted(t, argv)
	}
}

// REQ-11: "`sh script.sh` is the sanctioned wrapper-file form and never a
// defect"
// HAPPY PATH — NEGATIVE REQ.
//
// The remediation the refusal MESSAGE prescribes ("put it in a script and
// declare the script as argv0") must itself be green, under a wrapper too, or
// the refusal sends an author in circles — the premortem's first failure.
func TestReq11_TheSanctionedWrapperFileFormIsNeverADefect(t *testing.T) {
	for _, argv := range [][]string{
		{"sh", "./gate.sh"},
		{"sh", "scripts/gate.sh"},
		{"/bin/sh", "/opt/gate.sh"},
		{"nice", "sh", "./gate.sh"},
		{"env", "-i", "sh", "./gate.sh"},
		{"python", "tools/check.py"},
	} {
		wantAdmitted(t, argv)
	}
}

// REQ-12: "interpreter set: OPEN (deny-listed, not closed), unchanged from
// 0025:C5 — an unlisted spelling is admitted; the list grows by amendment of
// THIS clause."
// REQ-54: "unlisted spelling (`[\"python3\",\"-c\",\"…\"]`,
// `[\"perl\",\"-e\",\"…\"]`) | green … **silent by contract** — OPEN
// deny-list; the list grows only by amending C1"
// REQ-70 (briefly rejected): "no spelling folding (`python3` → `python`) and
// no closing of the interpreter set"
// DOMAIN EDGE — NEGATIVE REQ.
//
// Position-freedom widens over argv POSITION, never over which basenames are
// listed. Without this arm the suite reads as "inline shell is impossible",
// the guarantee C5 explicitly declines to make — and a position-free scan is
// exactly the change most likely to be over-applied to the name list too.
func TestReq12_TheInterpreterSetStaysOpenUnderThePositionFreeScan(t *testing.T) {
	for _, argv := range [][]string{
		{"perl", "-e", "print 1"},
		{"nice", "perl", "-e", "print 1"},
		{"python3", "-c", "print(1)"},
		{"env", "-i", "python3", "-c", "print(1)"},
		{"timeout", "5", "nodejs", "-e", "1"},
		// `busybox` as an unlisted spelling taking an inline-code flag.
		// DEVIATION D4 (TEST-FIXTURE): this row was authored as
		// `["xargs","busybox","sh","-c","echo"]`, which carries a real `sh`
		// at i=2 and its own `-c` at j=3 and so refuses under C1's `exists
		// i < j` quantifier — the same shape `Req28_…` requires to refuse
		// (`["wc","-l","python","-c"]`). The discriminating property the row
		// exists for is that `busybox` is UNLISTED, which needs no listed
		// pair beside it.
		{"xargs", "busybox", "-c", "echo"},
	} {
		wantAdmitted(t, argv)
	}
}

// REQ-13: "The membership itself is not restated here: it is
// `internal/table/load.go::shellInterpreters`, one map holding each listed
// name with its own inline-code flags (`sh`/`python` → `-c`, `ruby` → `-e`)"
// REQ-2 ASSUMPTION: "\"one of its listed flags\" means the flag list keyed by
// the interpreter found at `i`, not the union over all listed interpreters …
// so `[\"ruby\",\"-c\"]` does NOT refuse (`ruby` carries only `-e`)"
// BOUNDARY
//
// The per-interpreter flag list is what separates C1 from a union-of-all-flags
// reading. `ruby -c` is Ruby's syntax-CHECK flag, not an inline-code flag, and
// the union reading would refuse it. Asserted under a wrapper too, so the
// widened scan cannot quietly become a union.
func TestReq13_TheFlagListIsKeyedByTheInterpreterFoundNotTheUnion(t *testing.T) {
	t.Run("a flag belonging to another interpreter does not match", func(t *testing.T) {
		for _, argv := range [][]string{
			{"ruby", "-c", "tool.rb"},
			{"nice", "ruby", "-c", "tool.rb"},
			{"sh", "-e", "x"},
			{"python", "-e", "x"},
			{"php", "-c", "x"},
			{"node", "-c", "x"},
		} {
			wantAdmitted(t, argv)
		}
	})

	t.Run("the interpreter's OWN flag matches under a wrapper", func(t *testing.T) {
		for _, tc := range []struct {
			argv []string
			form string
		}{
			{[]string{"nice", "ruby", "-e", "puts 1"}, "ruby -e"},
			{[]string{"nice", "node", "-e", "1"}, "node -e"},
			{[]string{"nice", "node", "--eval", "1"}, "node --eval"},
			{[]string{"nice", "php", "-r", "echo 1;"}, "php -r"},
			{[]string{"nice", "python", "-c", "print(1)"}, "python -c"},
		} {
			wantInterpreterRefusal(t, tc.argv, tc.form)
		}
	})
}

// REQ-2 ASSUMPTION: "`i` and `j` range over the whole argv INCLUDING index 0,
// so the pre-existing argv0 forms (`[\"sh\",\"-c\",…]`) still refuse under the
// widened predicate; \"position-free\" widens the set, it never excludes argv0."
// REQ-41 / `0027:S1`: "the four REQ-74 probes and
// `neg/neg-command-shell-interpreter.toml` unedited and still green"
// REGRESSION / HAPPY PATH
//
// The widen is a SUPERSET claim. A scan written as `for i := 1` would pass
// every wrapper scenario in this file and silently un-refuse the four forms
// 0025 shipped, so the argv0 arm is asserted here as well as by the untouched
// 0025 probes.
func TestReq2_PositionFreedomIncludesArgv0AndDoesNotExcludeIt(t *testing.T) {
	for _, tc := range []struct {
		argv []string
		form string
	}{
		{[]string{"sh", "-c", "cat {artifact}"}, "sh -c"},
		{[]string{"bash", "-c", "echo hi"}, "bash -c"},
		{[]string{"python", "-c", "print(1)"}, "python -c"},
		{[]string{"env", "sh", "-c", "echo hi"}, "sh -c"},
		{[]string{"env", "A=1", "sh", "-c", "echo hi"}, "sh -c"},
	} {
		wantInterpreterRefusal(t, tc.argv, tc.form)
	}
}

// REQ-6 ASSUMPTION: "the tie-break is \"lowest `i` that is a listed
// interpreter AND has some qualifying `j > i`\". A listed interpreter at a
// lower `i` with no following flag must not short-circuit the scan"
// BOUNDARY
//
// The short-circuit bug this pins is invisible to every other scenario: an
// implementation that returned on the first listed BASENAME rather than on the
// first qualifying PAIR admits a real `sh -c` whenever a bare interpreter word
// precedes it. `["setup-tests","bash"]` (A7's own row) shows the bare trailing
// interpreter is green on its own.
func TestReq6_ALowerListedInterpreterWithNoFollowingFlagDoesNotStopTheScan(t *testing.T) {
	t.Run("bare listed word alone is green", func(t *testing.T) {
		for _, argv := range [][]string{
			{"setup-tests", "bash"},
			{"docker", "run", "sh"},
		} {
			wantAdmitted(t, argv)
		}
	})

	t.Run("a later qualifying pair is still found", func(t *testing.T) {
		wantInterpreterRefusal(t, []string{"ruby", "tool.rb", "sh", "-c", "echo"}, "sh -c")
		wantInterpreterRefusal(t, []string{"bash", "node", "-e", "1"}, "node -e")
	})
}

// --- C1: clause-3 coupling ------------------------------------------------

// REQ-14: "clause 3 coupling: 0025:C5 clause 3's whitespace exemption keys on
// this predicate, so a brace-bearing string under a wrapper reports the
// interpreter form, never command_unknown_placeholder"
// REQ-43 / `0027:S2`: "`[\"nice\",\"sh\",\"-c\",\"cat {artifact}\"]` — a
// brace-bearing command string under a wrapper. **Expected**:
// `command_shell_interpreter`, NOT `command_unknown_placeholder`."
// ADVERSARIAL
//
// The wrong-defect masking 0025 deviation D18 left open under a wrapper. The
// oracle is the category STRING, not "load failed": today this vector refuses
// too — as `command_unknown_placeholder` — so a test asserting only that the
// load failed passes against the unchanged code and proves nothing.
func TestReq43_ABraceBearingStringUnderAWrapperReportsTheInterpreterDefect(t *testing.T) {
	cases := []struct {
		name string
		argv []string
	}{
		{"nice", []string{"nice", "sh", "-c", "cat {artifact}"}},
		{"timeout 5", []string{"timeout", "5", "sh", "-c", "cat {artifact}"}},
		{"env -i", []string{"env", "-i", "bash", "-c", "echo {artifact} | tee out"}},
		{"unknown brace token under a wrapper", []string{"nice", "sh", "-c", "cat {nope}"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := refusalFor(t, tc.argv)
			if f.Category == table.CatCommandUnknownPlaceholder {
				t.Fatalf("%v reported %q — the WRONG defect. Clause 3's whitespace "+
					"exemption keys on the same `isInterp` the widened predicate "+
					"sets, so a wrapper must not mask the interpreter form "+
					"(`0027:C1` clause 3 coupling)", tc.argv, f.Category)
			}
			if f.Category != table.CatCommandShellInterpreter {
				t.Fatalf("%v reported %q; want %q", tc.argv, f.Category,
					table.CatCommandShellInterpreter)
			}
		})
	}
}

// REQ-14 ASSUMPTION: "clause 3's exemption keeps its existing shape
// (`!isInterp || no-whitespace ⇒ defect`); only the WRITER's predicate changes."
// REQ-15: "within-entry precedence is 0025:C5's, unchanged"
// BOUNDARY — NEGATIVE REQ.
//
// The exemption belongs to the INTERPRETER FORM, not to whitespace as such. If
// the widen were mis-implemented as "any brace-bearing whitespace element is
// exempt", a malformed placeholder under a non-interpreter argv0 would reach
// the executor as literal argv — the outcome 0025:C2 forbids.
func TestReq14_TheClause3ExemptionStillRequiresTheInterpreterForm(t *testing.T) {
	for _, argv := range [][]string{
		{"tool", "cat {nope}"},
		{"nice", "tool", "run {artifact} now"},
		{"sh", "./gate.sh", "cat {nope}"},
	} {
		f := refusalFor(t, argv)
		if f.Category != table.CatCommandUnknownPlaceholder {
			t.Errorf("%v reported %q; want %q — the exemption belongs to the "+
				"interpreter form, and under any other argv the element is "+
				"owned by nothing", argv, f.Category, table.CatCommandUnknownPlaceholder)
		}
	}
}

// --- C1: the report line --------------------------------------------------

// REQ-16: "report:     category string, remediation (\"inline shell is not a
// declared command; put it in a script and declare the script as argv0\") and
// position in `table.Categories()` are unchanged; the detail additionally
// names the two matched words"
// REQ-47 / `0027:S5`: "The regression set for the accepted false-refusal class
// — `[\"ruby\",\"tool.rb\",\"-e\",\"prod\"]` (a script whose own argument is
// `-e`). **Expected**: refuses, and the detail names `ruby -e` — the two
// matched words — so the author can see which pair collided."
// REQ-26: "the flag is matched at any later position, as today, so the
// already-accepted false-refusal class (a script argument that is itself `-c` /
// `-e`) is unchanged and is diagnosable from the named words"
// REQ-64: "the flag-anywhere laxity stays, so a script argument that is itself
// `-c`/`-e` after an interpreter word is refused (already true today); the
// message now names the words so the fix is visible."
// DOMAIN EDGE
//
// The premortem's first failure: the refusal is correct-by-contract but the
// MESSAGE was the failure, because "put it in a script" sends an author in
// circles when there is no inline code. Asserting only the category passes on
// today's message, which is why the detail text is the oracle here.
func TestReq47_TheAcceptedFalseRefusalClassNamesTheTwoWordsThatCollided(t *testing.T) {
	cases := []struct {
		name string
		argv []string
		form string
	}{
		{"ruby script with its own -e", []string{"ruby", "tool.rb", "-e", "prod"}, "ruby -e"},
		{"python script with its own -c", []string{"python", "run.py", "-c", "config.ini"}, "python -c"},
		{"node script with its own --eval", []string{"node", "cli.js", "--eval", "x"}, "node --eval"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wantInterpreterRefusal(t, tc.argv, tc.form)
		})
	}
}

// REQ-28: "Position-freedom additionally admits a NEW false-refusal class — a
// listed basename as a data argument under an unlisted argv0 — measured at A7
// … so the class is accepted as a stated cost"
// DOMAIN EDGE — NEGATIVE REQ.
//
// The class is ACCEPTED, not bounded. This test asserts the cost is actually
// paid: an implementation that special-cased it away — "only fire when the
// interpreter follows a known wrapper", or "only when argv0 is unlisted AND
// the flag adjoins" — would reintroduce the wrapper table REQ-5 forbids.
// A7's own constructed vectors, so the class is pinned as the record measured
// it and any future narrowing shows up here rather than as silent drift.
func TestReq28_TheNewFalseRefusalClassIsAcceptedNotSuppressed(t *testing.T) {
	cases := []struct {
		argv []string
		form string
	}{
		{[]string{"ls", "python", "-c"}, "python -c"},
		{[]string{"cat", "php", "-r"}, "php -r"},
		{[]string{"wc", "-l", "python", "-c"}, "python -c"},
		{[]string{"git", "log", "--grep", "node", "--eval"}, "node --eval"},
	}

	for _, tc := range cases {
		t.Run(strings.Join(tc.argv, "_"), func(t *testing.T) {
			wantInterpreterRefusal(t, tc.argv, tc.form)
		})
	}
}

// REQ-48 / `0027:S6`: "`table.Categories()` membership and order, and the
// `command_shell_interpreter` wire string. **Expected**: unchanged. The edit
// touches `load.go` only; the category constant and its position live in
// `category.go:60`/`:103` and are not edited"
// REQ-24 / `0027:D-naming`: "the category stays `command_shell_interpreter`;
// rejected: `command_inline_shell`"
// REQ-33: "the edit is confined to `load.go`"
// BOUNDARY — NEGATIVE REQ.
//
// DEVIATION D1 (pre-seeded, 7.1 pairwise 0027×0028 PW1). The golden is
// RELATIVE order, never tail- or `len`-equality: `0025:REQ-79` makes the
// list's total size a non-contract at any point, and a peer record appends
// further categories at the tail. A tail golden would break on that append
// while the position this clause protects — `command_shell_interpreter`'s
// place among its C5 neighbours — never moved.
func TestReq48_TheCategoryWireStringAndItsRelativePositionAreUnchanged(t *testing.T) {
	all := table.Categories()

	t.Run("the wire string is registered", func(t *testing.T) {
		if !slices.Contains(all, table.Category("command_shell_interpreter")) {
			t.Fatalf("Categories() carries no member with the wire string " +
				"`command_shell_interpreter`; the STRING is the contract, and " +
				"D-naming rejects renaming it")
		}
		if string(table.CatCommandShellInterpreter) != "command_shell_interpreter" {
			t.Errorf("CatCommandShellInterpreter = %q; want "+
				"`command_shell_interpreter` unchanged",
				table.CatCommandShellInterpreter)
		}
	})

	t.Run("the six C5 categories keep their RELATIVE clause order", func(t *testing.T) {
		// Relative, not absolute: index each member and require the
		// sequence to ascend. This is the assertion that survives a peer
		// record appending members at the tail (deviation D1) while still
		// failing if a category is moved, dropped, or reordered.
		clauseOrder := []table.Category{
			table.CatCommandAndPathConflict,
			table.CatCommandEmpty,
			table.CatCommandUnknownPlaceholder,
			table.CatCommandShellInterpreter,
			table.CatCommandOutputShape,
			table.CatCommandEnvConflict,
		}
		prev := -1
		for _, c := range clauseOrder {
			at := slices.Index(all, c)
			if at < 0 {
				t.Fatalf("Categories() does not carry %q", c)
			}
			if at <= prev {
				t.Errorf("%q is registered at index %d, at or before its "+
					"predecessor at %d; the six C5 categories keep their "+
					"relative clause order (`0027:S6`)", c, at, prev)
			}
			prev = at
		}
	})

	t.Run("no duplicate registration", func(t *testing.T) {
		seen := map[table.Category]bool{}
		for _, c := range all {
			if seen[c] {
				t.Errorf("category %q is registered twice", c)
			}
			seen[c] = true
		}
	})
}

// --- C1: determinacy and purity ------------------------------------------

// REQ-22: "Determinacy: n/a — no new signature, type, or API surface. C1
// replaces `interpreterForm`'s body at its existing name and signature"
// REQ-58: "`interpreterForm` is a pure function over a fixed argv slice with
// no shared state, invoked from the single call site … and this record adds no
// goroutine, no cache and no mutable package-level state."
// REQ-63: "Recovery: none needed for the predicate (no state)"
// REQ-59: "This record does not claim byte-identical output, content-addressed
// identity, or replay-stable hashes"
// ADVERSARIAL — NEGATIVE REQ.
//
// Purity, asserted as observable behaviour rather than by reading the source:
// the same vector must report the same form every time, and interleaving other
// vectors between two loads of the same one must not change either verdict.
// A cache or any mutable package-level state added under the widen shows up
// here as a divergent detail on a repeated load.
func TestReq58_ThePredicateIsPureAndCarriesNoStateAcrossLoads(t *testing.T) {
	const runs = 8
	probe := []string{"nice", "python", "sh", "-c", "echo"}

	var first string
	for i := range runs {
		// Interleave unrelated vectors so a cache keyed on anything but the
		// argv itself would be poisoned between the repeated probes.
		wantAdmitted(t, []string{"sh", "./gate.sh"})
		wantInterpreterRefusal(t, []string{"xargs", "sh", "-c", "echo"}, "sh -c")
		wantAdmitted(t, []string{"python", "-"})

		f := refusalFor(t, probe)
		if i == 0 {
			first = f.Detail
			continue
		}
		if f.Detail != first {
			t.Fatalf("run %d reported %q; run 0 reported %q — the predicate is a "+
				"pure function over a fixed argv slice with no shared state",
				i, f.Detail, first)
		}
	}
}

// REQ-66: "the predicate is O(len(argv)²) worst case against O(len(argv))
// today … the scan is not on any hot path and no alternative was rejected for
// cost."
// INPUT EDGE — NEGATIVE REQ.
//
// No performance assertion is owed and none is made: the quadratic scan is
// accepted by name. What IS asserted is that a long argv is answered
// correctly rather than truncated, capped, or short-circuited by any bound an
// implementation might add to "protect" the scan.
func TestReq66_ALongArgvIsScannedToTheEndWithNoBound(t *testing.T) {
	argv := make([]string, 0, 130)
	for range 128 {
		argv = append(argv, "filler")
	}
	argv = append(argv, "sh", "-c")

	wantInterpreterRefusal(t, argv, "sh -c")
}

// REQ-55: "`C1` states `no opt-in in v1`, which is deliberate: an opt-in flag
// would make the promise conditional and defeat the visibility the check
// exists for" … "No deprecation window, no compatibility flag, no staged
// rollout — the refusal set moves once, at the version that ships `C1`."
// REQ-70 (briefly rejected): "no declared inline-shell opt-in (`shell` field)"
// ADVERSARIAL — NEGATIVE REQ.
//
// The opt-in this forbids would land as an entry-level field, and the loader's
// unknown-schema-field arm is the surface that proves none was added: a
// `shell = true` escape hatch must refuse as an unknown field, not silently
// admit the interpreter form.
func TestReq55_NoInlineShellOptInFieldExistsOnAnEntry(t *testing.T) {
	src := swapEntry(t, cmdCarrierModel, cmdReadBlock, `[read.state]
role = "state"
command = ["nice", "sh", "-c", "echo hi"]
shell = true
keys = ["status"]
timeout = "2s"`)

	_, err := table.Load([]byte(src), "cmd-optin.toml")
	if err == nil {
		t.Fatalf("a `shell = true` opt-in loaded clean; C1 fixes `no opt-in in " +
			"v1` and an opt-in flag would make the promise conditional")
	}
	cat, ok := table.CategoryOf(err)
	if !ok {
		t.Fatalf("refusal carried no category: %v", err)
	}
	if cat != table.CatUnknownSchemaField && cat != table.CatCommandShellInterpreter {
		t.Errorf("category = %q; want the field refused as unknown, or the "+
			"interpreter form refused regardless — never an opt-in that "+
			"admits it", cat)
	}
}

// REQ-69 / `0027:ALT3` rejected: "only the `sh -s` axis belongs to the
// successor, and C1 routes exactly that" — the category stays a load-time
// REFUSAL, not an advisory.
// REQ-60: "That routing is an **ownership claim, not a closure** … until one
// does, the forms are admitted and documented as such" — NEGATIVE REQ: no code
// lands for the stdin axis in 0027.
// BOUNDARY — NEGATIVE REQ.
//
// ALT3 would have demoted the refusal to a hint whose text promises nothing.
// The oracle is that the wrapped form is a LOAD-TIME refusal — `table.Load`
// returns an error and no model — rather than a model that loaded carrying an
// advisory. Paired with the stdin arm, which must stay green because 0027
// ships no code for that axis.
func TestReq69_TheCategoryStaysALoadTimeRefusalAndTheStdinAxisShipsNoCode(t *testing.T) {
	t.Run("the wrapped interpreter form refuses at LOAD, not as an advisory", func(t *testing.T) {
		m, err := loadWithCommand(t, []string{"nice", "sh", "-c", "echo hi"})
		if err == nil {
			t.Fatalf("loaded clean; ALT3's advisory demotion is rejected — the " +
				"category is a load-time refusal")
		}
		if m != nil {
			t.Errorf("a refused load returned a model; load is fail-fast and " +
				"yields no model beside its refusal")
		}
	})

	t.Run("no stdin-axis refusal was added", func(t *testing.T) {
		// A5 is Pending and downgraded; the withholding point is the
		// successor's to add. If 0027 had shipped code for the axis, these
		// would refuse.
		for _, argv := range [][]string{
			{"sh", "-s"},
			{"python", "-"},
		} {
			wantAdmitted(t, argv)
		}
	})
}
