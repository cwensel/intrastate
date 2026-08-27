package cli

// RDR 0011 — Phase 3b adversarial failure-mode oracles.
//
// Phase 1's suite pins what C1/C2/C3 SAY. This file asks where the shipped
// implementation of those clauses breaks on a model class the record did
// not survey. Every oracle here is anchored in a `Trade-offs > Failure
// Modes` element and is written to FAIL against the current build.
//
// The common thread across all three is the demand-set term. C1 adds each
// row's MATCH-block owned keys to `invokedReaders`, and F6 names the
// behaviour-change class that term carries — but it names it for
// `flow resolve`, over "a row on that key" of "the requested outcome". The
// term as shipped applies to EVERY row of the model, ESCAPE ROWS INCLUDED,
// and an escape row is a row neither verb evaluates the way F6's remedy
// assumes:
//
//   - `flow next` SKIPS escape rows structurally (`flow_next.go`'s row loop
//     `if len(row.Escape) != 0 { continue }`), so no probe is ever built
//     from one and no candidate is ever reported for one. Yet its
//     match-owned keys now demand a reader, and an unbound one refuses the
//     WHOLE invocation exit 2 before the row loop runs.
//   - `flow resolve` consults an escape row only in the RESCUE phase, which
//     `internal/resolve.escapeOrRefuse` reaches only after the ordinary
//     rows produce `no_match` or `ambiguous_match`. A run whose ordinary
//     row yields a plan never reaches it — yet that run now refuses exit 2
//     for the escape row's reader.
//
// `0002:C4` makes this unavoidable rather than opt-in: EVERY rule, escape
// rules included, MUST carry a local match block, and `normalizeRule`
// refuses one that does not. DEV-8's guard term had the same shape but a
// guard block is optional on an escape rule; a match block is not. The
// `invokedReaders` doc comment already reasons about exactly this masking
// hazard for the guard term and fences it with the OUTCOME filter alone —
// a fence that is empty for `next`, which passes `outcome == ""`.

import (
	"slices"
	"testing"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
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

// FAILURE MODE: `flow resolve --outcome go` that returned a PLAN now exits
// 2, because an ESCAPE row binding the same outcome matches on an owned key
// whose reader is unbound. The rescue phase that would consult that row is
// unreachable on this run — `escapeOrRefuse` is called only from the
// `len(selected) == 0` and `default:` arms of `Resolve`'s selection switch,
// and this run's ordinary row selects exactly one — so the reader is
// demanded for a phase the run provably never enters.
//
// RDR ANCHOR: `0011:F6`, at its boundary. F6's accepted class is stated
// with a specific remedy and a specific reason: "the requested outcome has
// a row on that key … bind the reader … Each of those names the remedy the
// old refusal hid." C1 expands it: the class exists because "over a model
// with a match-only owned key, `resolve` today never invokes the serving
// reader and refuses `flow-no-match` for every row matching on that key,
// over an artifact that HOLDS the fact" — the old refusal HID a fact the
// caller could act on.
//
// That reasoning does not reach an escape row on a run that selects a plan.
// There is no hidden `flow-no-match` to redeem: the run succeeded, and it
// succeeded through `plain-row`, which needs nothing from `read.side`. C1's
// own acceptance sentence scopes the cost to "a run whose selected row does
// not itself match on the key" turning "from a plan into exit 2/3" — but it
// reaches that scope through "The demand set is a union over the OUTCOME's
// rows", reasoning about rows the SELECTION phase considers. `invokedReaders`
// already fences exactly this hazard for escape rows one comment above the
// new term ("an unbound or failing reader serving only that unrescuable row
// now refuses the whole request at exit 3 — masking a valid plan"), and
// fences it by OUTCOME alone — which does not exclude an escape row that
// happens to bind the requested outcome.
//
// ADVERSARIAL
func TestAdv2_ResolveKeepsItsPlanWhenOnlyAnEscapeRowDemandsTheReader(t *testing.T) {
	model := writeFlowModel(t, flowEscapeMatchOwnedModel)
	bind := escapeOwnedStateOnly(t, model)

	stdout, _, err := runCmd(t, "flow", "resolve", "--model", model,
		"--artifact", bind, "--outcome", "go", "--as=json")
	if err != nil {
		var ce *clierr.CLIError
		code := "<not a CLIError>"
		if asCLIError(err, &ce) {
			code = ce.Code
		}
		t.Fatalf("`flow resolve --outcome go` refused %s where it returns a "+
			"PLAN through `plain-row`.\n\n"+
			"`rescue-row` is an ESCAPE row. The kernel consults it only in "+
			"`escapeOrRefuse`, reached from the `len(selected) == 0` and "+
			"`default:` arms of Resolve's switch — never on a run whose "+
			"ordinary rows select exactly one. This run selects "+
			"`plain-row`, which matches on `status` alone.\n\n"+
			"F6 accepts a class where \"a `flow resolve` that used to "+
			"return a plan exits 2\" — but its stated warrant is that the "+
			"old refusal HID a fact: `resolve` \"refuses `flow-no-match` "+
			"for every row matching on that key, over an artifact that "+
			"HOLDS the fact\". Nothing is hidden here; the plan is "+
			"correct and complete without `mode`. `invokedReaders` already "+
			"fences this masking hazard for escape rows — \"masking a "+
			"valid plan\" is its own words — but fences it by OUTCOME "+
			"only, which an escape row binding `go` passes.\n\n"+
			"err = %v\nstdout = %s", code, err, stdout)
	}

	data := flowData(t, stdout)
	if rule, _ := data["rule"].(string); rule != "" && rule != "plain-row" {
		t.Errorf("plan rule = %q; want `plain-row` — the ordinary row that "+
			"needs nothing from `read.side`", rule)
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
