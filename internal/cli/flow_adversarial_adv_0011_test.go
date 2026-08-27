package cli

// RDR 0011 — Phase 3b adversarial failure-mode oracles, resolved in
// Phase 3c.
//
// Phase 1's suite pins what C1/C2/C3 SAY. This file asks where the shipped
// implementation of those clauses breaks on a model class the record did
// not survey. Every oracle here is anchored in a `Trade-offs > Failure
// Modes` element and was written to FAIL against the Phase 2 build.
//
// The common thread across all three is the demand-set term. C1 adds each
// row's MATCH-block owned keys to `invokedReaders`, and the Phase 2 build
// applied that term to EVERY row of the model, ESCAPE ROWS INCLUDED —
// while the two verbs consult an escape row very differently:
//
//   - `flow next` SKIPS escape rows structurally (`flow_next.go`'s row loop
//     `if len(row.Escape) != 0 { continue }`), so no probe is ever built
//     from one and no candidate is ever reported for one. C1 fixes `next`'s
//     predicate over "each NON-ESCAPE row" and justifies the demand-set
//     term by "the assembled view MUST actually carry the keys THE
//     PREDICATE READS"; the RDR's joint-check states it flatly —
//     "`runFlowNext` never lists escape rows". Demanding a reader there
//     refused the WHOLE invocation exit 2 above the row loop, for a fact no
//     reported row could consume, with no remedy (the key is owned, so
//     `--tag` is refused `flow-tag-owned`) and with `--all` inheriting the
//     same refusal — denying F1 its two-run diagnostic. ADV-1 and ADV-3
//     pinned that; Phase 3c scoped the match term to non-escape rows under
//     `next` (DEV-9).
//   - `flow resolve` DOES consult escape rows:
//     `internal/resolve.escapeOrRefuse` evaluates `view.matches(row.Match)`
//     on every escape row binding the requested outcome, and that phase's
//     reachability is not knowable before the kernel runs. C1 names and
//     ACCEPTS the resulting cost verbatim — such a run "turns from a plan
//     into exit 2/3, before the kernel and above the escape phase". ADV-2
//     originally denied that class; re-anchored in Phase 3c to pin its
//     BOUNDARY instead (DEV-9).
//
// `0002:C4` makes the shape unavoidable rather than opt-in: EVERY rule,
// escape rules included, MUST carry a local match block, and
// `normalizeRule` refuses one that does not. DEV-8's guard term had the
// same shape but a guard block is optional on an escape rule; a match block
// is not.

import (
	"slices"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
)

// --- fixtures ------------------------------------------------------------

// flowEscapeMatchOwnedModel authors the minimal shape the three oracles
// below share: ONE ordinary row that needs only `read.state`, and ONE
// ESCAPE row whose match block names `mode`, an owned key served solely by
// `read.side`.
//
// Nothing else in the model references `mode`: no row writes it, clears it,
// guards on it, or matches on it. Before C1's demand-set term `read.side`
// was therefore never invoked and the `side` role never had to be bound.
//
// The escape row is conformant on every axis `0002:C4`/`0009:C1` fence: it
// carries an `escape` list naming a modelable class, no write block, no
// clear list, and no gate list. Its match block is REQUIRED — `0002:C4`
// binds escape rules too ("EVERY transition rule MUST carry a local match
// block") — so a model author cannot avoid this shape by omitting it.
const flowEscapeMatchOwnedModel = `outcomes = ["go"]
terminal = ["done"]

[model]
id = "escapematchowned"
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

[tags.mode]
provenance = "owned"
kind = "enum"
domain = ["fast", "slow"]
single_valued = true

[read.state]
role = "state"
path = "flow.state"
keys = ["status"]
timeout = "2s"

[read.side]
role = "side"
path = "flow.side"
keys = ["mode"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status"]
timeout = "2s"
read_back = true

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "plain-row"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "final"

[[rule]]
id = "rescue-row"
escape = ["no_match"]
[rule.match.mode]
eq = "fast"
[rule.match.recognized]
eq = "go"
`

// escapeOwnedStateOnly seeds `status` on the `state` role and returns ONLY
// that binding. The `side` role is deliberately left unbound: that is the
// condition every oracle here turns on, and it is the condition under which
// the model ran clean before this RDR.
func escapeOwnedStateOnly(t *testing.T, model string) string {
	t.Helper()

	art := newFlowArtifact(t, "state.artifact")
	bind := artifactBinding(flowMatchRole, art)
	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", bind, "--write", "status=draft", "--as=json")
	return bind
}

// --- ADV-1 ---------------------------------------------------------------

// FAILURE MODE: `flow next` refuses `flow-artifact-missing` (exit 2) over a
// model whose ONLY unbound-reader demand comes from an ESCAPE row — a row
// `flow next` structurally never reports as a candidate and never builds a
// probe from. The verb refuses to run at all for a fact no reported row
// could ever consume.
//
// RDR ANCHOR: `0011:F6` — "a `flow resolve` that used to return a plan
// exits 2 (`flow-artifact-missing`) or 3 — the model matches on an owned
// key no rule writes, the requested outcome has a row on that key, and its
// declared reader is unbound or refusing; bind the reader (C1's demand-set
// clause, the class C1 names)."
//
// F6 is the RDR's ONE named behaviour-change class for the demand-set term,
// and it is scoped to `flow resolve`. C1 says so in the same voice: "The
// `flow resolve` VERB changes in exactly the one class the demand-set
// paragraph names, and nowhere else", and that paragraph's whole
// justification is "a row cannot be match-decided without the key". An
// escape row is never match-decided by `next`: the row loop skips it before
// `probeRow` is called. No Failure Mode in this record contemplates
// `flow next` REFUSING; F7 says the opposite — "none needed — `next` is
// effect-free in both modes; re-run with the corrected inputs" — and there
// is no correction to re-run with, because C1's own text records that
// `--tag` on an owned key is refused `flow-tag-owned` and "no other flag or
// ambient channel can supply it".
//
// The oracle is the run C1's demand-set paragraph itself promises is
// untouched: "No `flow resolve` run over a model WITHOUT such a key
// changes" — here `flow next` over a model where no CANDIDATE row carries
// such a key at all.
//
// ADVERSARIAL
func TestAdv1_NextDoesNotRefuseForAnEscapeRowsMatchOwnedKey(t *testing.T) {
	model := writeFlowModel(t, flowEscapeMatchOwnedModel)
	bind := escapeOwnedStateOnly(t, model)

	stdout, _, err := runCmd(t, "flow", "next", "--model", model,
		"--artifact", bind, "--as=json")
	if err != nil {
		var ce *clierr.CLIError
		code := "<not a CLIError>"
		if asCLIError(err, &ce) {
			code = ce.Code
		}
		t.Fatalf("`flow next` refused %s over a model whose only unbound "+
			"reader demand comes from an ESCAPE row.\n\n"+
			"`flow next` SKIPS escape rows: the row loop in flow_next.go "+
			"does `if len(row.Escape) != 0 { continue }` before probeRow, "+
			"so `rescue-row` is never probed, never reported, and its "+
			"`mode` atom never reaches the kernel. C1 justifies the "+
			"demand-set term by \"a row cannot be match-decided without "+
			"the key\"; this row is not match-decided at all.\n\n"+
			"F6 names ONE accepted behaviour-change class and scopes it to "+
			"`flow resolve` over \"the requested outcome[']s\" rows. No "+
			"Failure Mode contemplates `flow next` refusing, and F7 "+
			"records \"Recovery: none needed\" — yet there is no recovery "+
			"here: `mode` is owned, so `--tag mode=fast` is refused "+
			"`flow-tag-owned`, and C1 states no other channel supplies it. "+
			"`--all` cannot rescue it either: the refusal precedes the row "+
			"loop and the demand set is mode-independent by C1's own "+
			"requirement.\n\nerr = %v\nstdout = %s", code, err, stdout)
	}

	data := flowData(t, stdout)
	if got := candidateRules(t, data); !slices.Equal(got, []string{"plain-row"}) {
		t.Errorf("candidates = %v; want exactly [plain-row] — the escape "+
			"row is never an ordinary candidate", got)
	}
	if readers := readersOf(t, data); slices.Contains(readers, "side") {
		t.Errorf("readers = %v; `read.side` serves only `mode`, which no "+
			"REPORTED row matches on, writes, clears, or guards. C1's term "+
			"is justified by the candidate predicate reading the key, and "+
			"`next` never reads it for an escape row", readers)
	}
}

// --- ADV-2 ---------------------------------------------------------------

// BOUNDARY (Phase 3c: re-anchored; DEV-9). Phase 3b wrote this oracle to
// assert that `flow resolve --outcome go` KEEPS its plan when only an
// escape row demands the unbound reader. Against the record's own bytes
// that assertion is wrong, and it is wrong against the very clause it
// cited.
//
// C1's demand-set paragraph names this exact shape and ACCEPTS it,
// verbatim: "The demand set is a union over the outcome's rows, so the
// reader is invoked even when the row `resolve` would select does not
// itself match on the key: over an UNBOUND or REFUSING reader such a run
// turns from a plan into exit 2/3, before the kernel and ABOVE THE ESCAPE
// PHASE — a class this contract NAMES AND ACCEPTS here … That class is
// named in Consequences and pinned by S8; it is not denied."
//
// "Above the escape phase" is not incidental wording: it is the record
// contemplating precisely a run whose rescue phase is never reached and
// stating that the exit-2 there is accepted. F6, the anchor Phase 3b cited,
// scopes the class as "the requested outcome has a row on that key" —
// `rescue-row` binds `go` and matches on `mode`, so this fixture is INSIDE
// F6's scope, not outside it.
//
// The implementation reason is independent and points the same way.
// `internal/resolve.escapeOrRefuse` evaluates `view.matches(row.Match)` on
// every escape row binding the requested outcome, and `TagSet.matches` is
// TWO-VALUED: an absent key makes the escape row simply not match, and the
// rescue silently fails into the original `no_match`. Dropping escape rows
// from `resolve`'s demand set would therefore reinstate, in the rescue
// phase, the exact hidden-fact defect C1's term exists to remove — refusing
// over an artifact that HOLDS the fact. And the CLI cannot condition the
// demand on reachability: whether `Resolve` enters `escapeOrRefuse` is not
// knowable until after the kernel has run.
//
// So this oracle now pins the BOUNDARY rather than denying the class. It
// asserts the two halves that actually distinguish a correct build:
//
//  1. `resolve` still demands the escape row's reader — the accepted class
//     holds, and a build that "fixed" ADV-2 by narrowing `resolve` would
//     break the rescue path silently.
//  2. The refusal is `flow-artifact-missing` (exit 2) naming the ROLE, the
//     remedy F6 prescribes — not the mute `flow-no-match` the term
//     replaced, and not the reader's own refusal.
//
// Recorded as DEV-9 (TEST-FIXTURE): the assertion, not the implementation,
// contradicted the record.
//
// ADVERSARIAL
func TestAdv2_ResolveKeepsItsPlanWhenOnlyAnEscapeRowDemandsTheReader(t *testing.T) {
	model := writeFlowModel(t, flowEscapeMatchOwnedModel)
	bind := escapeOwnedStateOnly(t, model)

	_, _, err := runCmd(t, "flow", "resolve", "--model", model,
		"--artifact", bind, "--outcome", "go", "--as=json")
	if err == nil {
		t.Fatalf("`flow resolve --outcome go` returned a plan while the " +
			"role serving `rescue-row`'s match-owned `mode` is UNBOUND.\n\n" +
			"`escapeOrRefuse` evaluates `view.matches(row.Match)` on every " +
			"escape row binding the requested outcome, and `TagSet.matches` " +
			"is two-valued: an absent `mode` makes `rescue-row` silently " +
			"not match, so a rescue fails into the original `no_match` over " +
			"an artifact that HOLDS the fact — the hidden-fact defect C1's " +
			"demand-set term exists to remove. Reachability of the rescue " +
			"phase is not knowable before the kernel runs, so the demand " +
			"cannot be conditioned on it.\n\n" +
			"C1 NAMES AND ACCEPTS this cost verbatim: such a run \"turns " +
			"from a plan into exit 2/3, before the kernel and above the " +
			"escape phase\". Binding the `side` role restores the plan.")
	}

	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("`flow resolve` failed with a non-CLIError: %v", err)
	}
	if ce.Code != "flow-artifact-missing" {
		t.Errorf("`flow resolve` refused %q; want `flow-artifact-missing`. "+
			"F6 prescribes the remedy — \"bind the reader\" — and the "+
			"refusal must NAME the unbound role rather than surfacing as "+
			"the mute `flow-no-match` this term replaced", ce.Code)
	}
	if ce.Param != "side" {
		t.Errorf("refusal param = %q; want `side` — the role whose reader "+
			"serves the escape row's match-owned `mode`. C1 accepts the "+
			"exit only because it names what to bind", ce.Param)
	}
}

// --- ADV-3 ---------------------------------------------------------------

// FAILURE MODE: `--all` is not the escape hatch F1 makes it. F1 prescribes
// `--all` as the SOLE diagnostic for "a candidate the caller expected is
// absent" — and over this model class the caller cannot reach either mode:
// both refuse before the row loop, on a demand no reported row generates.
// The mode-independence C1 requires of the demand set is exactly what
// denies the remedy F1 promises.
//
// RDR ANCHOR: `0011:F1` — "a candidate the caller expected is absent.
// Diagnose with `--all`: if it appears there, the row was match-excluded …
// `--all` localizes the cause to match, and no further … That
// eye-comparison is the whole remedy today: `intrastate` ships no
// model-inspection verb … and this RDR adds none".
//
// F1's remedy is a two-run protocol: run the default, run `--all`, diff. It
// presumes BOTH runs produce a payload. C1 makes that presumption fail
// closed rather than open, and does so deliberately: "The demand set is a
// property of the MODEL … not of the mode: it is computed once per
// invocation, before the row loop, and is identical under --all". So the
// one mode F1 hands the caller for diagnosis inherits the same refusal, and
// the caller has no verb left — `dump` is asserted FOREIGN to the flow
// group (`flow_surface_0005_test.go:105`), and F1 records that this RDR
// adds no model-inspection verb.
//
// The oracle asserts only what F1 needs to be true: that `--all` returns a
// payload the caller can diff. It asserts nothing about which rows that
// payload holds — F1's own scope — so it cannot be satisfied by weakening
// the predicate.
//
// ADVERSARIAL
func TestAdv3_AllRemainsF1sDiagnosticOverAnEscapeRowDemand(t *testing.T) {
	model := writeFlowModel(t, flowEscapeMatchOwnedModel)
	bind := escapeOwnedStateOnly(t, model)

	for _, mode := range []struct {
		name string
		args []string
	}{
		{"default", nil},
		{"all", []string{"--all"}},
	} {
		t.Run(mode.name, func(t *testing.T) {
			args := []string{"flow", "next", "--model", model, "--artifact", bind}
			args = append(args, mode.args...)
			stdout, _, err := runCmd(t, append(args, "--as=json")...)
			if err != nil {
				var ce *clierr.CLIError
				code := "<not a CLIError>"
				if asCLIError(err, &ce) {
					code = ce.Code
				}
				t.Fatalf("`flow next %v` refused %s, so F1's diagnostic "+
					"protocol has no second run to diff against.\n\n"+
					"F1 prescribes `--all` as the WHOLE remedy for an "+
					"absent candidate: \"Diagnose with `--all`: if it "+
					"appears there, the row was match-excluded\". That is "+
					"a two-run protocol and it presumes both runs emit a "+
					"payload. Over this model neither does, because the "+
					"demand set is mode-independent by C1's own "+
					"requirement and the refusal precedes the row loop.\n\n"+
					"F1 also records that the caller has nothing else: "+
					"`intrastate` ships no model-inspection verb, `dump` "+
					"is FOREIGN to the flow group, and \"this RDR adds "+
					"none\".\n\nerr = %v\nstdout = %s",
					mode.args, code, err, stdout)
			}

			data := flowData(t, stdout)
			if _, ok := data["candidates"]; !ok {
				t.Errorf("payload carries no `candidates` under %s; F1's "+
					"diff has nothing to compare. keys = %v",
					mode.name, keysOf(data))
			}
		})
	}
}
