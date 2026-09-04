package cli

// RDR 0027 Stage 8 Phase 3b — ADVERSARIAL.
//
// Written against `0027:§failure-modes` by a reviewer whose working hypothesis
// is that the implementation is wrong. The predicate suite
// (`internal/table/inline_shell_scope_0027_test.go`) and the promise suite
// (`internal/cli/inline_shell_promise_0027_test.go`) are each internally
// thorough, and that is precisely the shape of the gap attacked here: they are
// two INDEPENDENT sources of truth for one contract, and nothing joins them.
//
// `0027:C1` promise: makes the shipped words a contract, not prose — "what a
// reviewer may rely on is the predicate line and the two out-of-scope forms
// named above … THESE words are the text it ships". The record's Risks line
// spells out the harm when the two drift: "'closes the wrapper class' is read
// as 'closes inline shell' and reviewers stop reading argv". §failure-modes
// then splits the outcomes by CHANNEL — a Visible refusal that names the two
// matched words, and a Silent admission "by contract, named in C1; diagnosis
// is reading argv, which C1's promise TELLS the reviewer to do for exactly
// those shapes".
//
// The words "exactly those shapes" are load-bearing and are what these tests
// hold. A description that names a shape the predicate refuses sends a
// reviewer to read argv for a binding that could never have loaded; a
// description that names a shape as admitted which the predicate later starts
// refusing has broken the promise rather than tightened it (`0027:S4`). Either
// way the reviewer's contract is false while every existing test stays green,
// because no existing test evaluates the description's own claims against the
// predicate.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// unwrapped collapses the rendered body's runs of whitespace to single
// spaces. The description is hard-wrapped for the terminal, so an anchor
// phrase can straddle a newline; matching on the unwrapped text asserts the
// SENTENCE is present without pinning the column the renderer wraps at.
func unwrapped(body string) string { return strings.Join(strings.Fields(body), " ") }

// admittedInDescription is every argv form the SHIPPED description names as
// admitted, transcribed from the rendered text rather than from `0027:C1`.
//
// Transcribing from the record would test the record against itself; the
// failure mode under attack is the SHIPPED text making a claim the SHIPPED
// predicate does not honour, so the left-hand side has to be what a reviewer
// actually reads. Each row carries the exact substring that anchors it in the
// body, so a test failure names the sentence that became false rather than
// only the argv that broke.
var admittedInDescription = []struct {
	anchor string   // a substring the description must contain to make this claim
	argv   []string // the form that claim tells a reviewer is admitted
}{
	{`env -S "sh -c …"`, []string{"env", "-S", "sh -c echo"}},
	{"sh -s", []string{"sh", "-s"}},
	{"bare sh", []string{"sh"}},
	{"sh -es", []string{"sh", "-es"}},
	{"python -", []string{"python", "-"}},
	{"node -", []string{"node", "-"}},
	{"sh script.sh", []string{"sh", "script.sh"}},
	{"python3", []string{"python3", "-c", "print(1)"}},
	{"nodejs", []string{"nodejs", "-e", "1"}},
	{"busybox", []string{"busybox", "sh", "-x"}},
}

// ADV-1 — the description's ADMITTED claims, evaluated against the predicate.
//
// `0027:C1` promise: "THESE words are the text it ships"; `0027:§failure-modes`
// Silent: the admitted forms lint green "by contract, named in C1".
// `0027:S4`: "a future predicate change that starts refusing one of them has
// broken the promise, not tightened it".
// ADVERSARIAL — CROSS-SURFACE.
//
// The description enumerates ten spellings by name and tells a reviewer each
// one loads. Nothing in either shipped suite reads those spellings OUT of the
// description and asks the loader. The table suite asserts admission for a
// list transcribed from the RECORD; the cli suite asserts the text mentions
// "stdin" and "one word". So the two lists can diverge silently: add `-s` to
// `shellInterpreters["sh"]`, or fold `python3` onto `python`, and the shipped
// sentence "sh -es … is admitted, however spelled" becomes a lie on the
// surface a reviewer is told to rely on, with every existing test green.
//
// The anchor assertion is half the test: it fails if the claim is dropped from
// the text as well as if the predicate stops honouring it, so the pair cannot
// be brought back into agreement by DELETING the promise.
func TestADV1_EveryFormTheDescriptionNamesAsAdmittedActuallyLoads(t *testing.T) {
	body := lintHelpAllBody(t)
	flat := unwrapped(body)

	for _, tc := range admittedInDescription {
		t.Run(strings.Join(tc.argv, "_"), func(t *testing.T) {
			if !strings.Contains(flat, tc.anchor) {
				t.Fatalf("the shipped description no longer names %q. C1's "+
					"promise: clause fixes the admitted forms IN THE TEXT, so "+
					"dropping the claim is not a way to keep text and predicate "+
					"in agreement — the reviewer loses the disclosure either "+
					"way\n\n%s", tc.anchor, body)
			}
			if _, err := table.Load([]byte(mvv0027Model(tc.argv)), "adv1.toml"); err != nil {
				t.Errorf("the description tells a reviewer %q is ADMITTED "+
					"(anchor %q), but the loader refuses it: %v\n\nA promise "+
					"the predicate does not honour is worse than none: the "+
					"reviewer is told to read argv for exactly these shapes "+
					"(`0027:§failure-modes` Silent) and the shape can never "+
					"reach them", tc.argv, tc.anchor, err)
			}
		})
	}
}

// refusedInDescription is the converse: forms the shipped description's
// REFUSED line commits to catching. The line names its wrappers explicitly
// ("env and its options, nice, timeout, xargs, doas") and then makes the
// stronger structural claim — "Nothing before the interpreter word is read, so
// no wrapper table exists and none is consulted", which promises the
// unenumerated prefix too.
var refusedInDescription = []struct {
	anchor string
	argv   []string
	form   string
}{
	{"env and its options", []string{"env", "-i", "sh", "-c", "x"}, "sh -c"},
	{"nice", []string{"nice", "sh", "-c", "x"}, "sh -c"},
	{"timeout", []string{"timeout", "5", "sh", "-c", "x"}, "sh -c"},
	{"xargs", []string{"xargs", "sh", "-c", "x"}, "sh -c"},
	{"doas", []string{"doas", "sh", "-c", "x"}, "sh -c"},
	{"wrappers nobody enumerated", []string{"runitwrap", "-q", "bash", "-c", "x"}, "bash -c"},
	{"at ANY later argv", []string{"sh", "./gate.sh", "--", "-c"}, "sh -c"},
}

// ADV-2 — the description's REFUSED claims, evaluated against the predicate,
// with the DETAIL as the oracle.
//
// `0027:§failure-modes` Visible: "a refused wrapper form names `sh -c` (the
// matched words) and the wrapper-file remediation".
// `0027:C1` report: "the detail additionally names the two matched words".
// ADVERSARIAL — CROSS-SURFACE.
//
// The mirror of ADV-1, and the half that carries the record's Risks line: a
// reviewer who reads "under any prefix … wrappers nobody enumerated" and finds
// an unenumerated prefix loading green has been told the class is closed when
// it is not. Asserting refusal alone is too weak — §failure-modes makes the
// Visible outcome the two matched words plus the remediation, so a refusal
// under the right category with a bare message still fails the reviewer the
// description addressed.
func TestADV2_EveryFormTheDescriptionClaimsToRefuseIsRefusedWithTheMatchedWords(t *testing.T) {
	body := lintHelpAllBody(t)
	flat := unwrapped(body)

	for _, tc := range refusedInDescription {
		t.Run(strings.Join(tc.argv, "_"), func(t *testing.T) {
			if !strings.Contains(flat, tc.anchor) {
				t.Fatalf("the shipped description no longer names %q; the "+
					"REFUSED line is the half of the promise a reviewer leans "+
					"on hardest\n\n%s", tc.anchor, body)
			}
			_, err := table.Load([]byte(mvv0027Model(tc.argv)), "adv2.toml")
			if err == nil {
				t.Fatalf("the description claims %q is refused (anchor %q), but "+
					"%v loads clean. This is the record's own Risks line — "+
					"\"'closes the wrapper class' is read as 'closes inline "+
					"shell' and reviewers stop reading argv\"", tc.anchor,
					tc.anchor, tc.argv)
			}
			f, ok := err.(*table.Failure)
			if !ok {
				t.Fatalf("%v refused with a non-categorized error: %v", tc.argv, err)
			}
			if f.Category != table.CatCommandShellInterpreter {
				t.Fatalf("%v refused as %q; the description promises %q",
					tc.argv, f.Category, table.CatCommandShellInterpreter)
			}
			if !strings.Contains(f.Detail, tc.form) {
				t.Errorf("%v detail = %q; `0027:§failure-modes` Visible makes the "+
					"matched words %q part of the outcome, not a nicety",
					tc.argv, f.Detail, tc.form)
			}
			if !strings.Contains(f.Detail, "put it in a script and declare the script as argv0") {
				t.Errorf("%v detail = %q; §failure-modes Visible pairs the "+
					"matched words with the WRAPPER-FILE REMEDIATION, and the "+
					"description's closing line (`sh script.sh` is sanctioned) "+
					"is what that remediation points at", tc.argv, f.Detail)
			}
		})
	}
}

// ADV-3 — the description must not name a form on BOTH sides.
//
// `0027:C1` out of scope, BY NAME; `0027:§failure-modes` Silent ("diagnosis is
// reading argv, which C1's promise tells the reviewer to do for exactly those
// shapes").
// ADVERSARIAL — NEGATIVE.
//
// ADV-1 and ADV-2 each hold one side against the predicate. Neither catches a
// description that is internally incoherent — one that lists a spelling under
// "Out of scope, BY NAME" while the REFUSED paragraph above it also covers
// that spelling. `sh -es` is the live hazard: `-es` is a single word today, so
// the predicate admits it, but the text sits three lines below a paragraph
// promising "one of that interpreter's own inline-code flags" at any later
// position. A reader who takes the refusal line at face value concludes
// `["sh","-es"]` is caught. The invariant that keeps the two paragraphs honest
// is exactly `0027:C1` reads: — the check "never splits a word on whitespace"
// — so the text must carry that sentence, and every admitted stdin spelling
// must in fact be a single word carrying no listed flag as a WHOLE element.
func TestADV3_TheAdmittedStdinSpellingsAreNotAlsoCoveredByTheRefusedLine(t *testing.T) {
	body := lintHelpAllBody(t)

	if !containsFold(unwrapped(body), "never splits a word on whitespace") {
		t.Errorf("the description drops C1's reads: sentence. Without it the "+
			"admitted `env -S \"sh -c …\"` row reads as a contradiction of the "+
			"refusal paragraph directly above it, and a reviewer resolves the "+
			"contradiction in favour of the stronger claim\n\n%s", body)
	}

	// Every spelling the text lists as stdin-admitted must be admitted for
	// the reason the text gives (the code is on the CHANNEL, so no listed
	// flag appears as a whole argv word) — not by accident of the flag map.
	for _, tc := range admittedInDescription {
		if !strings.Contains(tc.anchor, "-") && tc.anchor != "bare sh" {
			continue
		}
		t.Run(strings.Join(tc.argv, "_"), func(t *testing.T) {
			// A whole-element listed flag anywhere after a listed
			// interpreter word is the REFUSED shape by the description's own
			// first paragraph. If an admitted row ever grew one, the two
			// paragraphs would disagree on the same argv.
			m, err := table.Load([]byte(mvv0027Model(tc.argv)), "adv3.toml")
			if err != nil {
				t.Fatalf("%v is named admitted but refused: %v", tc.argv, err)
			}
			acc, ok := m.Readers["state"]
			if !ok {
				t.Fatalf("loaded model carries no `read.state` accessor")
			}
			if len(acc.Command) != len(tc.argv) {
				t.Errorf("loaded argv %#v is not the declared %#v; the check "+
					"reads argv WORDS and must not re-split the declaration",
					acc.Command, tc.argv)
			}
		})
	}
}

// ADV-4 — the promise as a reviewer reads it OUTSIDE the terminal: the
// COMMITTED `docs/cli-reference.md`, not a fresh in-memory render.
//
// `0027:C1` promise: the text is "mirrored into `docs/cli-reference.md`";
// `0027:REQ-38`/A6 settle that mirror as half the landing site.
// ADVERSARIAL — ORACLE STRENGTH.
//
// The shipped promise suite asserts the mirror by calling `writeCLIReference`
// on a live command tree and grepping the RESULT, and says why in a comment:
// "so the test does not depend on `make docs` having been run in the working
// tree". That is a defensible convenience and a strictly weaker oracle — it
// asserts the generator would produce the text, never that the file a reviewer
// opens carries it. The gap is not hypothetical: `docs/cli-reference.md` is a
// committed artifact, and every other assertion in this record's suites stays
// green against a stale one. §failure-modes' Silent row makes the reviewer's
// recourse "reading argv, which C1's promise tells the reviewer to do" — a
// promise they never received if the committed mirror predates Phase 3.
//
// Read against the checked-in bytes, so a build that regenerates the help body
// without regenerating the reference fails HERE rather than in a `make check`
// staleness gate outside the Go suite.
func TestADV4_TheCommittedCLIReferenceCarriesThePromiseNotJustTheGenerator(t *testing.T) {
	path := filepath.Join(repoRootFor(t), "docs", "cli-reference.md")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	ref := unwrapped(string(raw))

	if !strings.Contains(ref, string(table.CatCommandShellInterpreter)) {
		t.Fatalf("the COMMITTED docs/cli-reference.md never names %q. The "+
			"generator producing it is not the same claim as the artifact "+
			"carrying it, and the committed file is what a reviewer who is not "+
			"at a terminal reads (`0027:C1` promise:, REQ-38)",
			table.CatCommandShellInterpreter)
	}

	// The same two out-of-scope forms S7 requires of the terminal surface.
	// Naming the category alone would pass on a mirror that carried the
	// heading and dropped the body.
	for _, want := range []string{"STDIN", `env -S "sh -c …"`, "sh script.sh"} {
		if !strings.Contains(ref, want) {
			t.Errorf("the committed docs/cli-reference.md omits %q; the mirror "+
				"carries the same words the help surface does, or the reviewer "+
				"who reads only the docs is told the wrapper class is closed "+
				"and stops reading argv (`0027:§risks`)", want)
		}
	}

	// And the mirror must agree with the LIVE render, not merely contain
	// some description: a hand-edited reference that says something kinder
	// than the shipped text is the same defect wearing a different hat.
	var b strings.Builder
	writeCLIReference(&b, NewRootCmd())
	desc, ok := table.CategoryDescription(table.CatCommandShellInterpreter)
	if !ok {
		t.Fatal("the shell-interpreter category ships no description at all")
	}
	if !strings.Contains(unwrapped(b.String()), unwrapped(desc)) {
		t.Errorf("the freshly generated reference does not carry the accessor's " +
			"own text; the mirror and the accessor have diverged")
	}
	if !strings.Contains(ref, unwrapped(desc)) {
		t.Errorf("the COMMITTED docs/cli-reference.md does not carry the "+
			"accessor's current text verbatim — it is stale or hand-edited. "+
			"Regenerate it (`make docs`); C1 makes THESE words the text it "+
			"ships, on both halves of the landing site.\n\nwant:\n%s", desc)
	}
}
