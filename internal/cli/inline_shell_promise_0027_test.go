package cli

// RDR 0027 — the PROMISE half of `0027:C1`.
//
// `0027:S7` is the sole oracle for this half, and the record says why in as
// many words: "scenarios 1–5 assert what the predicate does, and none of them
// would fail if the description were missing, stale, or silent about the
// admitted forms" (REQ-50). So this file cannot be folded into the predicate
// suite in `internal/table` — every test there passes with no description
// shipped at all, and the gap that leaves is the premortem's second failure
// ("a reviewer who read 'closes the wrapper class' stops reading argv")
// reaching production with everything else green.
//
// The oracle is the RENDERED `--help-all` body, reached through the same
// production `ExecuteAndEmit` path a reviewer uses, with NO model loaded and
// NO defect present (`0027:S7`, A6). It is deliberately not a call on the
// accessor: an accessor-only assertion passes even if the body never renders
// the text, and reachability-without-a-refusal is the whole point of the
// clause. The accessor's Go name and signature are unconstrained
// implementation choice (`0027:REQ-36/REQ-38` ASSUMPTION), so nothing here
// names a symbol.
//
// PHASE-0 READING (Q2, recorded in coverage.md): "THESE words are the text it
// ships" is read as SUBSTANTIVE FIDELITY, not byte equality. A byte-equal copy
// of C1's fence would ship contract syntax (`base(argv[i])`, `i < j`,
// peer-record ids) to a CLI reader. S7's own oracle is "its text names both
// out-of-scope forms … in C1's channel-scoped words" — a naming assertion.
// Every check below is therefore a naming or scoping assertion, never a golden
// string.
//
// PHASE-0 READING (Q1, recorded in coverage.md): the description surface is
// PER-CATEGORY OPT-IN, not total over `table.Categories()`. A totality
// assertion would couple this record to every future category a peer record
// appends, contradicting `0025:REQ-79`'s size-is-not-a-contract rule. C1 ships
// text for exactly ONE category and S7 asserts exactly that one, so nothing
// below requires every category to carry text.

import (
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// lintHelpAllBody renders the `--help-all` body for `lint` the way a reviewer
// reads it: through ExecuteAndEmit, with no `--model`, so no model is loaded
// and no refusal is provoked. `internal/cli/help_all.go` keeps this stream off
// the respond gateway (REQ-39), which is what makes the read possible at all.
func lintHelpAllBody(t *testing.T) string {
	t.Helper()

	stdout, _ := runHelp(t, "lint", "--help-all")
	return stdout
}

// REQ-17: "promise:    what a reviewer may rely on is the predicate line and
// the two out-of-scope forms named above, carried in a description Phase 3
// ships on the `--help-all` surface"
// REQ-37: "No description surface exists for `table.Category` today: it is a
// bare string and `Categories()` returns identifiers only … so this phase adds
// the surface as well as the text — the description must be reachable without
// provoking a refusal, which the detail string (`load.go::carrierDefect`) is
// not."
// REQ-49 / `0027:S7`: "The description Phase 3 ships for
// `command_shell_interpreter`, read from the `--help-all` extended-help body
// (A6) with no model loaded and no defect present, so nothing routes through
// the refusal gateway. **Expected**: it exists, and its text names both
// out-of-scope forms — the one-word shell string and the stdin-fed
// interpreter — in C1's channel-scoped words."
// HAPPY PATH
//
// The category-under-description must be named on the surface at all. Today
// `Categories()` returns identifiers only and nothing renders them, so this
// starts red by construction — which is the control the `oracle` mini-check's
// S7 row calls for.
func TestReq49_TheHelpAllSurfaceDescribesTheShellInterpreterCategory(t *testing.T) {
	body := lintHelpAllBody(t)

	if !strings.Contains(body, string(table.CatCommandShellInterpreter)) {
		t.Fatalf("`lint --help-all` never names %q. C1's promise: clause says a "+
			"reviewer may rely on the predicate line and the two out-of-scope "+
			"forms, \"carried in a description Phase 3 ships on the --help-all "+
			"surface\"; with no description shipped, the predicate line alone is "+
			"what holds and the promise overstates it.\n\n%s",
			table.CatCommandShellInterpreter, body)
	}
}

// REQ-49 oracle row (`0027:§mini-check-oracle`, S7): "assert the description
// text contains BOTH admitted forms (the one-word string and the stdin
// channel), not that a description is non-empty: a placeholder string passes a
// non-empty check"
// REQ-21: "the lint's description carries it [C1's out-of-scope line]
// verbatim; the docs never say \"closes inline shell\""
// REQ-65: "`env -S \"sh -c …\"` and the stdin-fed forms (`sh -s`, bare `sh`,
// `sh -es`, `python -`, `node -`) are admitted and documented as such"
// ADVERSARIAL
//
// The discriminating assertion of the whole promise half. A placeholder
// description ("the command declares an inline shell interpreter") satisfies
// existence and says nothing a reviewer can act on, so the oracle is that BOTH
// out-of-scope forms are named.
func TestReq49_TheDescriptionNamesBothOutOfScopeFormsByName(t *testing.T) {
	body := lintHelpAllBody(t)

	t.Run("form one: a shell string carried in ONE word", func(t *testing.T) {
		// C1 names the form and gives `env -S "sh -c …"` as its spelling.
		// Either the spelling or the one-word wording must appear, or a
		// reviewer cannot tell that `["env","-S","sh -c echo"]` is admitted.
		named := strings.Contains(body, "env -S") ||
			(containsFold(body, "one word") || containsFold(body, "single word") ||
				containsFold(body, "ONE word"))
		if !named {
			t.Errorf("the description never names C1's FIRST out-of-scope form — "+
				"\"a shell string carried in ONE word (`env -S \\\"sh -c …\\\"`)\". "+
				"A reviewer who reads only \"closes the wrapper class\" ships an "+
				"`env -S` string; naming the form is the mitigation C1 fixes.\n\n%s",
				body)
		}
	})

	t.Run("form two: an interpreter reading its script from STDIN", func(t *testing.T) {
		if !containsFold(body, "stdin") {
			t.Errorf("the description never names C1's SECOND out-of-scope form — "+
				"\"an interpreter that reads its script from STDIN (`sh -s`, bare "+
				"`sh`, `sh -es`, `python -`, `node -`)\".\n\n%s", body)
		}
	})
}

// REQ-10: "The channel is the scope: any listed interpreter taking its code on
// stdin rather than as a later argv word is admitted, however spelled."
// REQ-9: "an interpreter that reads its script from STDIN … whose remit is
// what a command entry CONSUMES and so covers the non-shell spellings too"
// REQ-49 / `0027:S7`: "in C1's channel-scoped words"
// DOMAIN EDGE
//
// The scoping, not merely the mention. C1 widened this line from "a stdin-fed
// shell" to the CHANNEL precisely so a reviewer is not left to discover
// `["python","-"]` themselves — "the one thing this record exists to stop"
// (§approach). A description that named only shell spellings would pass the
// previous test and re-open exactly that gap, so the text must show the scope
// reaches beyond `sh`.
func TestReq10_TheDescriptionScopesTheStdinFormByChannelNotByShellSpelling(t *testing.T) {
	body := lintHelpAllBody(t)

	// Either the channel wording ("however spelled" / "any listed
	// interpreter") or a non-shell spelling C1 names (`python -`, `node -`)
	// demonstrates the scope. A description carrying neither is
	// spellings-scoped and misleads on the very case C1 calls out.
	channelScoped := containsFold(body, "however spelled") ||
		containsFold(body, "any listed interpreter") ||
		strings.Contains(body, "python -") ||
		strings.Contains(body, "node -")
	if !channelScoped {
		t.Errorf("the description names stdin but not its CHANNEL scope. C1: "+
			"\"any listed interpreter taking its code on stdin rather than as a "+
			"later argv word is admitted, however spelled\" — `python` is on the "+
			"deny-list, so a spellings-scoped sentence leaves a reviewer to "+
			"discover `[\"python\",\"-\"]` themselves.\n\n%s", body)
	}
}

// REQ-19: "THESE words are the text it ships."
// REQ-1 / REQ-2: the predicate line — "a listed interpreter word followed, at
// ANY later argv position, by one of that interpreter's inline-code flags —
// under any prefix"
// REQ-17: "what a reviewer may rely on is the predicate LINE and the two
// out-of-scope forms"
// HAPPY PATH
//
// The promise is two things, not one: the predicate line AND the out-of-scope
// forms. A description that named only what is admitted would leave the
// reviewer without the sentence they are meant to be able to hold — "if an
// interpreter and its inline-code flag appear as two separate words in that
// order, anywhere in argv, the binding is refused" (§approach).
func TestReq17_TheDescriptionAlsoCarriesThePredicateLine(t *testing.T) {
	body := lintHelpAllBody(t)

	if !containsFold(body, "argv") {
		t.Errorf("the description never says the check reads argv WORDS; the "+
			"promise is the predicate line AND the out-of-scope forms\n\n%s", body)
	}
	// Position-freedom is the delta this record ships. A description written
	// against the old argv0 predicate is STALE text beside a widened
	// predicate — the drift S7 exists to catch.
	positionFree := containsFold(body, "any position") ||
		containsFold(body, "any later") ||
		containsFold(body, "anywhere in argv") ||
		containsFold(body, "any prefix") ||
		containsFold(body, "position-free")
	if !positionFree {
		t.Errorf("the description does not state the POSITION-FREE predicate. "+
			"C1: \"a listed interpreter word followed, at ANY later argv "+
			"position, by one of that interpreter's inline-code flags — under "+
			"any prefix\". Text written against the old argv0 rule is stale "+
			"beside the widened predicate.\n\n%s", body)
	}
}

// REQ-21: "the docs never say \"closes inline shell\"" — NEGATIVE REQ on the
// shipped wording.
// REQ-55 / `0027:G-cross-cutting`: "an opt-in flag would make the promise
// conditional and defeat the visibility the check exists for"
// ADVERSARIAL — NEGATIVE REQ.
//
// The premortem's second failure as a test on the wording itself. An
// overclaiming sentence is worse than no description: a reviewer who believes
// inline shell is closed stops reading argv, which is precisely when the
// admitted forms ship.
func TestReq21_TheShippedWordingNeverClaimsInlineShellIsClosed(t *testing.T) {
	body := lintHelpAllBody(t)
	lower := strings.ToLower(body)

	for _, overclaim := range []string{
		"closes inline shell",
		"prevents inline shell",
		"makes inline shell impossible",
		"inline shell is impossible",
		"all inline shell",
		"blocks all shell",
	} {
		if strings.Contains(lower, overclaim) {
			t.Errorf("the shipped wording contains %q. The deny-list is "+
				"\"defense in depth over that, not the barrier itself\" "+
				"(0025:C5), and C1's mitigation is that the docs never say "+
				"\"closes inline shell\".", overclaim)
		}
	}
}

// REQ-39: "`internal/cli/help_all.go` keeps that stream off the respond
// gateway, so it renders with no defect present."
// REQ-49 / `0027:S7`: "with no model loaded and no defect present, so nothing
// routes through the refusal gateway"
// BOUNDARY — NEGATIVE REQ.
//
// The description must be reachable WITHOUT provoking a refusal — that is the
// property the refusal detail string (`load.go::carrierDefect`) lacks and the
// reason Phase 3 adds a surface at all (REQ-37). Asserted under `--as=json`,
// where the never-silent contract reserves stdout for the single terminal
// record: help carries no envelope, so a description routed through the
// gateway shows up here as one.
func TestReq39_TheDescriptionRendersOffTheRespondGatewayWithNoDefectPresent(t *testing.T) {
	stdout, _ := runHelp(t, "lint", "--help-all", "--as=json")

	if strings.Contains(stdout, `"type":"ok"`) || strings.Contains(stdout, `"code":`) {
		t.Errorf("the extended-help stream emitted an envelope; help is not a "+
			"terminal result and the description must render with no defect "+
			"present\n\n%s", stdout)
	}
	if !strings.Contains(stdout, string(table.CatCommandShellInterpreter)) {
		t.Errorf("the description is absent from the help stream under "+
			"--as=json; S7 reads it with no model loaded and no defect "+
			"present\n\n%s", stdout)
	}
}

// REQ-38: "Where it lands is settled (A6, Verified): the `--help-all`
// extended-help body … registered via `withExtendedHelp` and mirrored into
// `docs/cli-reference.md` by `internal/cli/docs.go::runDocs`."
// REQ-18: "the landing site is settled, not open."
// BOUNDARY
//
// The mirror is half the landing site: `docs/cli-reference.md` is generated
// wholesale from the live command tree and staleness-gated in `make check`, so
// a description that renders on `--help-all` reaches the committed reference
// by the same derivation. Asserted through the renderer rather than by reading
// the committed file, so the test does not depend on `make docs` having been
// run in the working tree.
func TestReq38_TheDescriptionReachesTheCLIReferenceMirrorByTheSameDerivation(t *testing.T) {
	var b strings.Builder
	writeCLIReference(&b, NewRootCmd())
	ref := b.String()

	if !strings.Contains(ref, string(table.CatCommandShellInterpreter)) {
		t.Errorf("the generated docs/cli-reference.md body never names %q; the "+
			"landing site is the --help-all body MIRRORED into the committed "+
			"reference, so a description reachable in one is reachable in the "+
			"other by construction", table.CatCommandShellInterpreter)
	}
	if !containsFold(ref, "stdin") {
		t.Errorf("the generated docs/cli-reference.md body never names the " +
			"stdin out-of-scope form; the mirror carries the same words the " +
			"help surface does")
	}
}

// REQ-20: "The surface carries no admitted-form disclosure until Phase 3 lands
// it — S7 is that assertion, and until it passes the predicate line alone is
// what holds"
// REQ-50: "scenarios 1–5 assert what the predicate does, and none of them
// would fail if the description were missing, stale, or silent about the
// admitted forms."
// DOMAIN EDGE
//
// The structural claim that S7 is the SOLE oracle for the promise half, made
// checkable: the description text lives on the help surface, and no refusal
// detail carries it. If a future change moved the admitted-form disclosure
// into the refusal message instead, the promise would again be readable only
// by provoking a defect — the failure REQ-37 names.
func TestReq20_TheAdmittedFormDisclosureLivesOnTheHelpSurfaceNotInARefusal(t *testing.T) {
	body := lintHelpAllBody(t)
	if !containsFold(body, "stdin") {
		t.Fatalf("the help surface carries no admitted-form disclosure; until " +
			"it does, the predicate line alone is what holds")
	}

	// The refusal detail is not the disclosure site: it fires only ON a
	// defect and so by construction cannot tell a reviewer what is ADMITTED
	// (`0027:A6`). Nothing requires it to stay silent about the forms, but a
	// description reachable ONLY there fails the reachability half of S7,
	// which the help-surface assertion above is what pins.
	// DEVIATION D7 (TEST-FIXTURE): `[tags.recognized]` is declared because
	// `table.Load` refuses a model without one as
	// `malformed_model_declaration` (load.go:270) BEFORE reaching C5's
	// command clauses — so without it this fixture asserts nothing about
	// the interpreter form. The sibling MVV fixture declares it too.
	src := `outcomes = ["advance"]
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
command = ["nice", "sh", "-c", "echo hi"]
keys = ["status"]
timeout = "2s"

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"
`
	_, err := table.Load([]byte(src), "cmd-0027-promise.toml")
	if err == nil {
		t.Fatal("the wrapper form loaded clean; the predicate half must be red " +
			"before the promise half means anything")
	}
	if cat, ok := table.CategoryOf(err); ok && cat != table.CatCommandShellInterpreter {
		t.Errorf("the wrapped form refused as %q; want %q", cat,
			table.CatCommandShellInterpreter)
	}
}

// containsFold is a case-insensitive substring test. The description's exact
// casing is an authoring choice C1 does not fix (Q2: substantive fidelity, not
// byte equality), so the assertions above must not fail on capitalization.
func containsFold(haystack, needle string) bool {
	return strings.Contains(strings.ToLower(haystack), strings.ToLower(needle))
}
