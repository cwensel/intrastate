package cli

// RDR 0011 — C3: the help text, the shipped user-facing prose, and the
// five-read `unresolved` → `unknown` census.
//
// C3 exists because "the override is visible where a reader meets it — help
// text and the test file — not only in this record." Two of its obligations
// are unusual and both have oracles here:
//
//   - the census is a COUNT with `file:line` provenance from lock time. If
//     the tree has drifted the lines are stale but the count is the
//     contract, so the oracle reconciles by SYMBOL: it counts shipped reads
//     of the key across the named files, never by line number.
//   - the prose sweep covers COMMENT and failure-message text a green suite
//     never flags and that "goes factually false the moment the field is
//     renamed." No runtime assertion can see it, so the oracle reads the
//     files.

import (
	"slices"
	"strconv"
	"strings"
	"testing"
)

// c3ReHomedFiles are the three test files C3's five-read census names, in
// its own order.
var c3ReHomedFiles = []string{
	"internal/cli/flow_next_0005_test.go",
	"internal/cli/flow_adversarial_0005_test.go",
	"internal/cli/flow_mvv_0005_test.go",
}

// REQ-54: "The build MUST mechanically re-home every shipped read of the
// `unresolved` key to `unknown` and its element type from string to
// `{key, reason}` — FIVE TEST reads on this build:
// `internal/cli/flow_next_0005_test.go` (3 — the presence check `:103`, the
// gate-id assertion `:163`, and the adversarial guard-fact assertion
// `:404`), `internal/cli/flow_adversarial_0005_test.go` (1, `:446`), and
// `internal/cli/flow_mvv_0005_test.go` (1, `:66`)."
// REQ-110 / `0011:S6`: "the 5 test reads of `unresolved` are re-homed to
// `unknown` with the ok-bool asserted".
// ADVERSARIAL — the census, reconciled by SYMBOL rather than by line. The
// `file:line` numbers are lock-time provenance; the COUNT is the contract,
// and the lines drift with any unrelated edit above them.
func TestReq54And110_TheFiveShippedReadsOfUnresolvedAreReHomedToUnknown(t *testing.T) {
	// The distribution C3 fixes, per file. A sixth occurrence in
	// `flow_next_0005_test.go` is "the `%#v` argument of the `t.Fatalf`
	// reporting `:163`'s failure — it re-homes with that read and is not a
	// read of its own", so the count is of READS, and the payload-key
	// literal is what a read is spelled with.
	wantReads := map[string]int{
		"internal/cli/flow_next_0005_test.go":        3,
		"internal/cli/flow_adversarial_0005_test.go": 1,
		"internal/cli/flow_mvv_0005_test.go":         1,
	}

	var total int
	for _, file := range c3ReHomedFiles {
		src := readRepoFile(t, file)

		// Every shipped read of the OLD key must be gone: the rename IS
		// the break, so a build that left one behind is reading a field
		// the payload no longer carries.
		if n := strings.Count(src, `"unresolved"`); n != 0 {
			t.Errorf("%s still reads the OLD payload key `unresolved` %d "+
				"time(s); C3 requires every shipped read re-homed to "+
				"`unknown`, and the element type from string to "+
				"{key, reason}", file, n)
		}

		got := strings.Count(src, `"unknown"`)
		total += got
		if want := wantReads[file]; got < want {
			t.Errorf("%s reads the `unknown` key %d time(s); the census "+
				"fixes %d for this file. The `file:line` numbers C3 cites "+
				"are LOCK-TIME provenance and drift with any edit above "+
				"them — the COUNT is the contract, reconciled by symbol",
				file, got, want)
		}
	}
	if total < 5 {
		t.Errorf("the three census files carry %d reads of `unknown` in "+
			"total; the census fixes FIVE. A re-homing that dropped a read "+
			"rather than moving it has deleted an assertion", total)
	}
}

// REQ-58: "THREE of the five discard the comma-ok of a helper
// (`flow_harness_0005_test.go::stringsAt` …) — `flow_next_0005_test.go:404`,
// `flow_mvv_0005_test.go:66`, and `flow_adversarial_0005_test.go:446` — and
// the last of those is a NEGATIVE assertion (`if containsString(unresolved,
// \"flag\")`) that passes vacuously on the empty slice, so a re-homed read
// MUST assert the ok bool, and a green suite MUST NOT be cited as evidence
// that this contract was implemented."
// ADVERSARIAL — the ok-bool is what stops a re-homed read from passing
// against a payload that carries no such field at all.
func TestReq58_TheReHomedReadsAssertTheOkBoolRatherThanDiscardingIt(t *testing.T) {
	// VACUITY GUARD. This clause is about the re-homed reads, so it says
	// nothing until they exist: before the rename the files read the OLD
	// key and the discarding forms below simply are not present, so the
	// scan passes for the wrong reason. That is precisely the failure mode
	// C3 warns about — "a green suite MUST NOT be cited as evidence that
	// this contract was implemented" — so the guard is the assertion's
	// first half, not decoration.
	for _, file := range c3ReHomedFiles {
		if strings.Contains(readRepoFile(t, file), `"unresolved"`) {
			t.Fatalf("%s still reads the OLD key `unresolved`; the ok-bool "+
				"clause is about the RE-HOMED reads and is not yet "+
				"assertable. A pass here before the rename would be exactly "+
				"the vacuous green C3 forbids citing", file)
		}
	}

	for _, file := range c3ReHomedFiles {
		src := readRepoFile(t, file)

		// The discarding form. `stringsAt` returns `nil,false` on a
		// missing key, a non-array value, OR a non-string element — so a
		// discarded ok bool turns a MISSING FIELD into an empty slice, and
		// a negative assertion over an empty slice passes vacuously in
		// every possible build.
		for _, discarded := range []string{
			`, _ := unknownAt(`,
			`, _ := stringsAt(c, "unknown")`,
			`, _ := objectsAt(c, "unknown")`,
		} {
			if strings.Contains(src, discarded) {
				t.Errorf("%s discards the comma-ok of a payload read "+
					"(`%s`); a re-homed read MUST assert the ok bool. "+
					"Without it a missing `unknown` field reads as an empty "+
					"list and the assertion passes vacuously — which is why "+
					"a green suite MUST NOT be cited as evidence this "+
					"contract was implemented", file, strings.TrimSpace(discarded))
			}
		}
	}
}

// REQ-59: "`stringsAt` itself MUST NOT be repointed at the new element type
// … it is shared harness used well beyond these five sites, and changing its
// return shape would touch readers this contract has not counted. The
// re-homing adds a SIBLING helper beside it (a `[]{key, reason}` reader over
// the same payload) and switches the five reads to it."
// ADVERSARIAL — the harness fence. A sibling is additive; a repoint is a
// change to readers the census never counted.
func TestReq59_StringsAtKeepsItsShapeAndTheReHomingAddsASiblingHelper(t *testing.T) {
	harness := readRepoFile(t, "internal/cli/flow_harness_0005_test.go")

	// VACUITY GUARD: the fence is about the re-homing, so it bites only
	// once the re-homing has happened. Before it, `stringsAt` trivially
	// keeps its shape because nothing asked it to change.
	for _, file := range c3ReHomedFiles {
		if strings.Contains(readRepoFile(t, file), `"unresolved"`) {
			t.Fatalf("%s still reads the OLD key `unresolved`; the harness "+
				"fence is about the RE-HOMING and is not yet assertable",
				file)
		}
	}

	if !strings.Contains(harness, "func stringsAt(m map[string]any, key string) ([]string, bool)") {
		t.Errorf("`stringsAt`'s signature changed in " +
			"`flow_harness_0005_test.go`; it is SHARED harness used well " +
			"beyond C3's five sites, and repointing it at the new element " +
			"type would touch readers this contract has not counted")
	}

	// The sibling exists and reads the new element type. It is HARNESS,
	// not a sixth re-homing: "the census counts shipped READS of the
	// `unresolved` key."
	sibling := readRepoFile(t, "internal/cli/flow_fixtures_0011_test.go")
	if !strings.Contains(sibling, "func unknownAt(") {
		t.Error("no sibling `[]{key, reason}` reader is defined beside " +
			"`stringsAt`; the re-homing adds one and switches the five " +
			"reads to it")
	}
}

// REQ-60: "The same rename MUST also sweep the COMMENT and failure-message
// prose that names `unresolved` … On this build that is the doc and header
// comments and the `t.Errorf`/`t.Fatalf` strings at
// `internal/cli/flow_next_0005_test.go:23,66,104,135,137,141,169-171,404-408`,
// `internal/cli/flow_mvv_0005_test.go:46,73`, and
// `internal/cli/flow_adversarial_0005_test.go:448,450`."
// ADVERSARIAL — "text a green suite never flags and that goes factually
// false the moment the field is renamed." The line numbers are lock-time
// provenance and drift; the obligation is that the WORD is gone from those
// files, which is what the oracle asserts.
func TestReq60_TheCommentAndFailureMessageProseNamingUnresolvedIsSwept(t *testing.T) {
	// The census's own scope: prose in the three files C3 names. These
	// "are NOT re-homings under the five-read census above and no ok-bool
	// applies, but leaving them stale ships a suite whose messages name a
	// field the payload no longer carries."
	for _, file := range c3ReHomedFiles {
		src := readRepoFile(t, file)
		if n := strings.Count(src, "unresolved"); n != 0 {
			var lines []string
			for i, line := range strings.Split(src, "\n") {
				if strings.Contains(line, "unresolved") {
					lines = append(lines,
						strings.TrimSpace(line)+"  ("+file+":"+strconv.Itoa(i+1)+")")
				}
			}
			t.Errorf("%s names `unresolved` %d time(s) in prose:\n  %s\n"+
				"The rename sweeps COMMENT and failure-message text too. "+
				"The line numbers C3 cites are lock-time provenance and "+
				"drift with any edit above them; the obligation is that the "+
				"word is gone", file, n, strings.Join(lines, "\n  "))
		}
	}
}

// REQ-53: "The 0005 flow next tests in internal/cli/flow_next_0005_test.go
// MUST keep passing unchanged where they assert the alphabet, gate
// handling, reader narrowing, determinism, and non-mutation; an assertion
// that depended on a match-excluded row being reported MUST be re-homed
// under --all, not deleted, and the file's header comment MUST name this
// RDR as the source of the default."
// REQ-111 / `0011:S6`: "\"None re-homed under `--all`\" is a claim about the
// DIFF, not a runtime assertion — no test can observe it — so it is
// discharged at review by reading the changed test files".
// DOMAIN EDGE — the assertable half of REQ-53 is the header naming this
// RDR; the diff half is a review obligation this oracle names rather than
// pretends to check.
func TestReq53And111_The0005NextFileHeaderNamesThisRDRAsTheSourceOfTheDefault(t *testing.T) {
	src := readRepoFile(t, "internal/cli/flow_next_0005_test.go")

	// The header is the first comment block, before the imports.
	header := src
	if i := strings.Index(src, "\nimport ("); i > 0 {
		header = src[:i]
	}

	if !strings.Contains(header, "0011") {
		t.Errorf("the header comment of `flow_next_0005_test.go` does not "+
			"name RDR 0011 as the source of the default:\n%s\n"+
			"The override must be visible where a READER meets it — the "+
			"help text and this file — not only in the record. The header "+
			"still frames the verb as one that ENUMERATES, which is the "+
			"predicate 0011 overrides", header)
	}

	// A4 predicts ZERO moves under `--all` over the shipped fixtures, and
	// "any move found refutes A4 rather than satisfying this scenario". So
	// the assertable form is the NEGATIVE: no 0005 `next` oracle passes
	// the flag. A move that WAS needed would show up here.
	if strings.Contains(src, `"--all"`) {
		t.Errorf("`flow_next_0005_test.go` passes `--all`; A4 classified " +
			"all 28 `flow next` oracles against C1's predicate and reported " +
			"MOVES-UNDER-`--all` = 0. A move found here REFUTES A4 rather " +
			"than satisfying S6 — record it as a deviation rather than " +
			"absorbing it")
	}
}

// REQ-51: "The command's Short and Long help MUST state that flow next
// reports the candidates the supplied state can take (match and guard both
// holding or undecided), that a match or guard key the state does not carry
// leaves a row a candidate with the key listed under `unknown` and its
// reason named, and that --all reports every row the guards do not exclude
// regardless of match."
// REQ-55: "the production rename additionally touches … the shipped
// USER-FACING string in this file, not a comment — the cobra `Long` help at
// `internal/cli/flow_next.go:70-71` (\"next ENUMERATES rather than selects …
// their ids are reported as unresolved facts\")".
// REQ-121 / PH 2: "rewrite Short/Long help per C3".
// DOMAIN EDGE — the help is shipped user-facing prose, so the oracle drives
// `--help` through the CLI rather than reading the source string.
func TestReq51And55And121_TheHelpStatesTheMatchConditionedDefaultAndTheAllFlag(t *testing.T) {
	stdout, stderr, err := runCmd(t, "flow", "next", "--help")
	if err != nil {
		t.Fatalf("`flow next --help` failed: %v", err)
	}
	help := stdout + stderr

	// The OLD framing must be gone: it states the predicate this RDR
	// overrides, and it is shipped user-facing text, not a comment.
	for _, stale := range []string{"ENUMERATES", "unresolved facts"} {
		if strings.Contains(help, stale) {
			t.Errorf("the help still carries %q:\n%s\nThat text states the "+
				"OVERRIDDEN default — `next` enumerating rather than "+
				"selecting — and the payload field it names no longer "+
				"exists", stale, help)
		}
	}

	// REQ-51 binds the PROSE, not the rendered blob. Cobra generates the
	// `Flags:` listing from each flag's registered usage string, and
	// `--all`'s registration in `flow_next.go` reads "report every row the
	// guards do not exclude, regardless of match" — the clause verbatim.
	// Asserting over the whole blob therefore passes even if the entire
	// `--all` paragraph is deleted from `Long`, which is exactly the
	// vacuity this oracle exists to close. Cut the generated sections off
	// and assert against what the author wrote.
	// Lower-cased: the shipped text shouts HOLD / UNDECIDED for emphasis,
	// and emphasis is styling, not contract. Case must not decide whether
	// a clause counts as stated.
	prose := strings.ToLower(helpProse(t, help))

	// Each clause is pinned by a SET of phrasings, not one golden string:
	// the help was just rewritten and will be edited again, and an oracle
	// that fails on any rewording gets weakened rather than fixed. Every
	// alternative below states the same fact, so the set fails only when
	// the CLAUSE is dropped — not when it is rephrased.
	for _, want := range []struct {
		clause string
		anyOf  []string
		why    string
	}{
		{
			clause: "the match-conditioned predicate",
			anyOf: []string{
				"match and guard both hold or are undecided",
				"match and its guards both hold or are undecided",
				"whose match and guard both hold or are undecided",
				"match and guards both hold or are undecided",
			},
			why: "REQ-51 requires the help to state WHAT a candidate is " +
				"under the new default: a rule whose MATCH and GUARD both " +
				"HOLD or are UNDECIDED. Naming `candidate` alone leaves " +
				"the caller with 0005's enumerate-everything reading",
		},
		{
			clause: "an absent key leaves the row a candidate",
			anyOf: []string{
				"does not carry does not exclude",
				"the state does not carry does not exclude",
				"absent key does not exclude",
				"leaves the row a candidate",
				"leaves it a candidate",
			},
			why: "REQ-51 requires the help to state that a match or guard " +
				"key the state DOES NOT CARRY does not exclude the row. " +
				"This is the clause that separates `next` from a filter: " +
				"absence is undecided, not false",
		},
		{
			clause: "the absent key is reported under `unknown` with its reason",
			anyOf: []string{
				"under unknown and its reason named",
				"listed under unknown and its reason",
				"under unknown with its reason",
				"unknown and its reason named",
			},
			why: "REQ-51 requires BOTH halves: the key is listed under " +
				"`unknown` AND its reason is named. `unknown` alone is a " +
				"field name; without the reason the caller cannot tell " +
				"an absent fact from an unevaluated gate",
		},
		{
			clause: "`--all` reports guard-legal rows regardless of match",
			anyOf: []string{
				"--all reports every row the guards do not exclude, regardless of match",
				"every row the guards do not exclude, regardless of match",
				"regardless of match",
			},
			why: "REQ-51 requires the help to state that `--all` reports " +
				"every row the GUARDS do not exclude REGARDLESS OF MATCH. " +
				"Cobra's flag listing carries these words for free, so " +
				"this is asserted against the authored prose only",
		},
		{
			clause: "the `not-evaluated` reason",
			anyOf: []string{
				"not-evaluated",
			},
			why: "the reason vocabulary is half of the `unknown` contract: " +
				"without `--evaluate-gates` a gate id is unknown for the " +
				"reason `not-evaluated`, not because a fact is missing. " +
				"C1 names the token, so the help must spell it",
		},
	} {
		var stated bool
		for _, phrase := range want.anyOf {
			if strings.Contains(prose, phrase) {
				stated = true
				break
			}
		}
		if !stated {
			t.Errorf("the help's prose does not state %s: %s\n"+
				"None of the accepted phrasings %q appear in:\n%s",
				want.clause, want.why, want.anyOf, prose)
		}
	}
}

// helpProse returns the authored help body — everything cobra's template
// emits BEFORE the generated `Usage:` / `Flags:` / `Global Flags:`
// sections. Those sections are assembled from flag registrations, not
// from `Long`, so a REQ-51 clause found only there is not help the author
// wrote and is not evidence the clause is stated.
func helpProse(t *testing.T, help string) string {
	t.Helper()
	const boundary = "\nUsage:"
	i := strings.Index(help, boundary)
	if i < 0 {
		t.Fatalf("`flow next --help` has no `Usage:` section, so the "+
			"authored prose cannot be separated from cobra's generated "+
			"flag listing:\n%s", help)
	}
	prose := strings.TrimSpace(help[:i])
	if prose == "" {
		t.Fatalf("`flow next --help` renders no prose before `Usage:` — "+
			"the command's `Long` is empty, and every REQ-51 clause is "+
			"unstated:\n%s", help)
	}
	return prose
}

// REQ-52: "The help MUST also state that a candidate is a row the supplied
// state does NOT EXCLUDE, not a row flow resolve will select: a candidate
// carrying an `unknown` entry may still be refused by flow resolve over the
// same state, and the entry names what to supply."
// ADVERSARIAL — this is the clause that keeps a caller from reading the
// narrowed list as a decision. `next` and `resolve` genuinely DISAGREE on a
// row whose match key is absent, and the help is where that is said.
func TestReq52_TheHelpSeparatesCandidateFromWhatResolveWillSelect(t *testing.T) {
	stdout, stderr, err := runCmd(t, "flow", "next", "--help")
	if err != nil {
		t.Fatalf("`flow next --help` failed: %v", err)
	}
	help := stdout + stderr
	lower := strings.ToLower(help)

	// The shipped help already names `flow resolve` — but only to say that
	// turning a gate DENY into a refusal is its job. That is `0005`'s
	// separation, not this clause's, so mentioning the verb is not the
	// assertion: the help must say the candidate list is what the state
	// does not EXCLUDE, and that an `unknown` entry may still be refused.
	if !strings.Contains(lower, "resolve") {
		t.Fatalf("the help never mentions `flow resolve` at all:\n%s", help)
	}

	// (i) `unknown` must be named in the same breath as the refusal — the
	// entry is what tells the caller what to supply.
	if !strings.Contains(help, "unknown") {
		t.Errorf("the help does not name the `unknown` list:\n%s\nA "+
			"candidate carrying an `unknown` entry may still be refused by "+
			"`flow resolve` over the same state, and THE ENTRY NAMES WHAT "+
			"TO SUPPLY — which is the caller's remedy", help)
	}

	// (ii) the separation itself: a candidate is a row the supplied state
	// does not EXCLUDE, not a row `resolve` will SELECT. Without this the
	// caller reads the narrowed list as a decision, which the whole
	// `D-selection-predicate` decision forbids.
	var separated bool
	for _, phrase := range []string{
		"does not exclude", "not exclude", "may still be refused",
		"still be refused", "not a row flow resolve will select",
		"not what flow resolve will select",
	} {
		if strings.Contains(lower, phrase) {
			separated = true
		}
	}
	if !separated {
		t.Errorf("the help does not state that a candidate is a row the "+
			"supplied state does NOT EXCLUDE rather than a row "+
			"`flow resolve` will select:\n%s\nThe shipped text separates "+
			"the two verbs only on GATE DENIALS (0005's clause). This "+
			"clause is about the candidate list itself: `next` and "+
			"`resolve` genuinely DISAGREE on a row whose match key is "+
			"absent — `next` lists it naming the key, `resolve` folds it "+
			"into `no_match`", help)
	}
}

// REQ-56: "Two further shipped user-facing descriptions state the OVERRIDDEN
// default and MUST be corrected with the help text, though neither reads the
// payload field: `docs/cli-output-contract.md:125` (\"Enumerate the legal
// outcomes and their candidate rules\") and `README.md:43` (\"list legal
// next outcomes\")."
// REQ-122 / PH 3: "document `flow next`'s payload — including `unknown` and
// the `--all` invocation — in `docs/cli-output-contract.md`, which currently
// shows only the invocation grammar".
// DOMAIN EDGE — shipped descriptive prose, corrected alongside the help. No
// CONTRACT clause of the output-contract doc is invalidated ("it pins the
// envelope, not per-command payload fields"), so these are prose
// corrections rather than a contract amendment.
func TestReq56And122_TheShippedDescriptionsOfTheOldDefaultAreCorrected(t *testing.T) {
	t.Run("cli-output-contract", func(t *testing.T) {
		doc := readRepoFile(t, "docs/cli-output-contract.md")

		if strings.Contains(doc, "Enumerate the legal outcomes") {
			t.Error("`docs/cli-output-contract.md` still describes " +
				"`flow next` as \"Enumerate the legal outcomes and their " +
				"candidate rules\"; that is the predicate this RDR overrides")
		}
		// Phase 3 documents the PAYLOAD there, which today "shows only the
		// invocation grammar".
		for _, want := range []string{"unknown", "--all"} {
			if !strings.Contains(doc, want) {
				t.Errorf("`docs/cli-output-contract.md` does not document "+
					"%q; Phase 3 documents `flow next`'s payload — "+
					"including `unknown` and the `--all` invocation — where "+
					"the doc currently shows only the invocation grammar",
					want)
			}
		}
	})

	t.Run("readme", func(t *testing.T) {
		readme := readRepoFile(t, "README.md")

		if strings.Contains(readme, "list legal next outcomes") {
			t.Error("`README.md` still describes `flow next` as \"list " +
				"legal next outcomes\"; it is the other shipped description " +
				"of the old default")
		}
	})
}

// REQ-109 / `0011:S5`: "A11's key-set agreement is PINNED here …: an oracle
// MUST compare `assembledView`'s key set against the kernel's
// `internal/resolve/resolve.go::assemble` over the same inputs, modulo the
// kernel's `recognized` key"
// ADVERSARIAL — "Today the agreement holds only because `runFlowNext`
// passes that one slice to both consumers (A11, A16); no other scenario
// fails specifically when it breaks, which is what this oracle fixes." A
// refactor that stops threading one `owned` slice to both fails HERE rather
// than silently.
//
// The kernel's `assemble` is unexported and `internal/resolve` is untouched
// by this contract, so the comparison is made against its SPECIFICATION —
// owned ∪ observed ∪ `recognized`, which is what `assemble` builds from the
// `Input` the CLI hands it — reconstructed from the payload's own `owned`
// and `observed` maps. That keeps the oracle inside `internal/cli` and
// keeps `internal/resolve` unchanged, which S5 requires in the same breath.
func TestReq109_TheCLIViewsKeySetAgreesWithTheKernelsModuloRecognized(t *testing.T) {
	model := writeFlowModel(t, flowMatchClassesModel)
	bind := seedMatchArtifact(t, model, "status=draft")

	data := runNext(t, model, []string{bind}, "--tag", "phase=alpha")

	// `assembledView`'s key set, as the payload reports it: owned ∪
	// observed, last write wins.
	cliKeys := viewKeysOf(t, data)

	// The kernel's `assemble` over the SAME inputs adds exactly one key:
	// the reserved `recognized`, which no atom can name — `normalize.go`
	// lifts the match-block `recognized` atom into `Row.Outcome`, and a
	// `recognized` atom under `guard.all`/`guard.unless` is refused
	// `CatMalformedOutcomeBinding` at load.
	for _, k := range cliKeys {
		if k == "recognized" {
			t.Errorf("`assembledView` carries the reserved key "+
				"`recognized`: %v\nThe CLI's view does NOT inject it, "+
				"unlike the kernel's `assemble` — that asymmetry is "+
				"harmless precisely because it is the ONE key the two "+
				"builders differ on, and it is unnameable by any atom",
				cliKeys)
		}
	}

	// The agreement in the direction that matters: every key a MATCH atom
	// can name is present in the CLI's view exactly when the reader or
	// `--tag` established it, so a match atom the CLI KEEPS is never over
	// a key the kernel's view lacks. The observable form is that a
	// present-and-equal key DECIDES (the row survives) and a present-and-
	// unequal one EXCLUDES — neither of which a key-set disagreement
	// could produce, since a key the kernel lacked would fold into
	// `no_match` and drop the row silently.
	if !slices.Contains(cliKeys, "phase") {
		t.Errorf("the CLI's view = %v; `--tag phase=alpha` was supplied, so "+
			"`phase` must be present. An empty-but-present value yields a "+
			"`Tag` with `Value: \"\"`, which `assembledView` stores as a "+
			"present map entry and the kernel stores as a present "+
			"`taggedValue` — agreement, not divergence", cliKeys)
	}
	requireCandidate(t, data, "match-alpha")
	if c := candidateNamed(t, data, "match-beta"); c != nil {
		t.Errorf("`match-beta` survives `phase=alpha`: %#v\nA key the "+
			"kernel's view LACKED would fold into `no_match` and drop the "+
			"row silently; a key it HAS decides. This arm is what makes the "+
			"key-set agreement observable at all", c)
	}
}

// REQ-108 / `0011:S5`: "`TestReq78_MatchPatternStillFoldsAbsenceIntoNonMatch`,
// the `match_conflicted_test.go` oracles, and the reason-set pins
// `TestReq49`/`TestReq47` … stay green as shipped."
// ADVERSARIAL — the upstream pins A13's merge relies on. They must still
// EXIST and be unedited; running them is `go test ./internal/resolve`'s job
// and their unchangedness is asserted as a diff by MVV 9's oracle.
func TestReq108_TheUpstreamKernelPinsStillExistAsShipped(t *testing.T) {
	// All three ship in one file on this build. Their HOME is provenance;
	// the obligation is that the symbols still exist, unedited — the
	// unchangedness itself is asserted as a diff by MVV 9's oracle.
	src := readRepoFile(t, "internal/resolve/guard_atoms_test.go")

	for _, pin := range []struct{ symbol, why string }{
		{"TestReq78_MatchPatternStillFoldsAbsenceIntoNonMatch",
			"the two-valued match seam this RDR leaves DEFERRED " +
				"(`0007:REQ-78`): the kernel still folds an absent match " +
				"key into non-match, and `flow next` reports the same row " +
				"as a candidate naming the key"},
		{"TestReq49_ReasonIsAClosedNamedStringSetWithAnEnumerator",
			"the reason set's closure, which C1 REUSES rather than mints"},
		{"TestReq47_ReasonSetStaysAtExactlyTwoMembers",
			"`Reasons()` at exactly `absent`/`uncomparable` — the closure " +
				"A13's merge relies on, and the reason `not-evaluated` is " +
				"minted on the CLI list ONLY"},
	} {
		if !strings.Contains(src, pin.symbol) {
			t.Errorf("the kernel pin %s is gone; it pins %s. It stays green "+
				"AS SHIPPED", pin.symbol, pin.why)
		}
	}
}
