package cli

// RDR 0005 — Phase 3b adversarial suite.
//
// Three failure modes the Phase 1 suite's own fixtures cannot reach, each
// anchored on a row of the RDR's `## Trade-offs > ### Failure Modes`
// section or on the gate-disposition mini-check that section carries.
//
// The common shape of all three: the CLI answers a question it does not
// own, and its answer disagrees with the seam that DOES own it, silently.
//
//	ADV-1  the CLI's own row-exclusion filter under-decides relative to the
//	       kernel, so `next --evaluate-gates` reports a row the model has
//	       already ruled out and INVOKES that row's gate accessor.
//	       (FM mini-check, row "gate on a guard-excluded row"; REQ-42)
//
//	ADV-2  the narrowed read-accessor set is computed from `RequiresOwned`
//	       alone, which RDR 0002 derives from the write block and clear
//	       list — NOT from the guard. A row guarding on an owned key it
//	       does not write leaves its reader uninvoked, and `resolve`
//	       refuses `flow-guard-unevaluable` over an artifact that holds the
//	       fact and a role the caller bound. (FM kernel row
//	       `flow-guard-unevaluable`; REQ-35, REQ-106)
//
//	ADV-3  `revision` renders the model ID. No model can declare a revision
//	       — RDR 0002's `[model]` schema has no such field — so REQ-25's
//	       "a model that declares none renders `revision` empty rather than
//	       a CLI-invented value" governs EVERY model, and a stand-in makes
//	       REQ-110's request identity blind to the thing it names.
//
// Every fixture here is a complete document RDR 0002's loader ACCEPTS, so
// a failure under test is the `flow` verb's and never a load refusal.

import (
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
)

// --- ADV-1 --------------------------------------------------------------
//
// `internal/cli::excluded` decides a row's guard FALSE only for an `eq`
// atom under `all` over a present value. Two shapes the loader admits fall
// straight through it:
//
//   - a `guard.unless` block. The kernel's row verdict is
//     `all ∧ ¬unless` (`0007:C5`/`0007:C6`), so an `unless` conjunct
//     decided TRUE excludes the row. `excluded` skips every non-`all`
//     block outright.
//   - an `all` atom on any operator other than `eq` — `in` here, which is
//     an ordinary authoring shape RDR 0003 admits over a single-valued
//     enum.
//
// In both, the FM mini-check's "gate on a guard-excluded row" row says
// "not run, not reported" and REQ-42 says "a row whose guards already
// exclude it is not a candidate, so its gates MUST NOT run". Each fixture
// pairs the excluded row with a gate that DENIES, so the test can prove
// the accessor really ran rather than inferring it: a reported `deny`
// result is an accessor invocation that the contract forbids.
//
// The companion `resolve` assertion is what makes this non-circular. It
// asks the KERNEL — the seam that owns guard semantics (REQ-113) — the
// same question over the same state, and the kernel answers `no_match`.
// So the row is guard-excluded as a matter of the model, not of the test's
// opinion.

// advUnlessExcludedModel: `unless flag eq true` over `flag = true` — the
// kernel excludes the row. Its gate denies.
const advUnlessExcludedModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "advunless"
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

[tags.flag]
provenance = "owned"
kind = "bool"
single_valued = true
required = true

[read.state]
role = "state"
path = "flow.state"
keys = ["status", "flag"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status", "flag"]
timeout = "2s"
read_back = true

[gate.blocked]
role = "state"
path = "flow.gate.deny"
keys = ["status"]
timeout = "2s"

[initial]
status = "draft"
flag = "true"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "unless-excluded"
gate = ["blocked"]
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.guard.unless.flag]
eq = "true"
[rule.write]
status = "final"
`

// advInExcludedModel: `all tier in ["mid","high"]` over `tier = low` — the
// kernel excludes the row. Its gate denies.
const advInExcludedModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "advin"
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

[tags.tier]
provenance = "owned"
kind = "enum"
domain = ["low", "mid", "high"]
single_valued = true
required = true

[read.state]
role = "state"
path = "flow.state"
keys = ["status", "tier"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status", "tier"]
timeout = "2s"
read_back = true

[gate.blocked]
role = "state"
path = "flow.gate.deny"
keys = ["status"]
timeout = "2s"

[initial]
status = "draft"
tier = "low"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "in-excluded"
gate = ["blocked"]
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.guard.all.tier]
in = ["mid", "high"]
[rule.write]
status = "final"
`

// TestAdv1_AGuardExcludedRowIsNeitherReportedNorGated proves the FM
// mini-check's "gate on a guard-excluded row → not run, not reported" over
// the two exclusion shapes `excluded` does not decide.
func TestAdv1_AGuardExcludedRowIsNeitherReportedNorGated(t *testing.T) {
	cases := []struct {
		name  string
		model string
		rule  string
		seed  []string
	}{
		{
			name:  "unless-block",
			model: advUnlessExcludedModel,
			rule:  "unless-excluded",
			seed:  []string{"status=draft", "flag=true"},
		},
		{
			name:  "all-block-in-operator",
			model: advInExcludedModel,
			rule:  "in-excluded",
			seed:  []string{"status=draft", "tier=low"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			model := writeFlowModel(t, tc.model)
			art := seedArtifact(t, model, tc.seed...)
			bind := artifactBinding(flowStateRole, art)

			// The kernel is the authority on guard semantics (REQ-113).
			// Asking it first establishes that the row really is excluded,
			// so the assertions below test the CLI and not the fixture.
			_, _, rerr := runCmd(t, "flow", "resolve", "--model", model,
				"--artifact", bind, "--outcome", "advance", "--as=json")
			if rerr == nil {
				t.Fatalf("fixture is not discriminating: the kernel SELECTED "+
					"rule %q, so it is not guard-excluded", tc.rule)
			}
			if got := codeOfError(rerr); got != "flow-no-match" {
				t.Fatalf("fixture is not discriminating: resolve refused %q; "+
					"a guard-excluded sole row leaves zero matches", got)
			}

			data := flowData(t, requireSuccess(t, "flow", "next",
				"--model", model, "--artifact", bind,
				"--evaluate-gates", "--as=json"))

			candidates, ok := objectsAt(data, "candidates")
			if !ok {
				t.Fatalf("`candidates` is not an array of objects: %#v",
					data["candidates"])
			}
			for _, c := range candidates {
				if c["rule"] != tc.rule {
					continue
				}
				// REQ-42 / FM mini-check: not reported.
				t.Errorf("rule %q is reported as a candidate, but its guard "+
					"already excludes it — the kernel answers no_match over "+
					"the same state (REQ-42, FM mini-check row \"gate on a "+
					"guard-excluded row\")", tc.rule)

				// REQ-42 / FM mini-check: not run. A reported gate result is
				// proof the accessor was invoked.
				gates, hasGates := c["gates"]
				if !hasGates {
					continue
				}
				arr, _ := gates.([]any)
				if len(arr) != 0 {
					t.Errorf("rule %q reports %d gate result(s) %v; a "+
						"guard-excluded row's gates MUST NOT run (REQ-42)",
						tc.rule, len(arr), gates)
				}
			}
		})
	}
}

// --- ADV-2 --------------------------------------------------------------
//
// `internal/cli::invokedReaders` computes the invoked read-accessor set
// from `table.Row.RequiresOwned`. RDR 0002 derives that field from the
// rule's WRITE BLOCK and CLEAR LIST only — `internal/resolve::Row`'s own
// doc says so, and adds that "the guard's key set is not required to be a
// subset of this one".
//
// So a row that GUARDS on an owned key it never writes demands that key
// without naming it in `RequiresOwned`. The reader serving it is skipped,
// the kernel is handed a view missing the fact, and `resolve` refuses
// `flow-guard-unevaluable` — the FM kernel row — for an invocation where
// the reader was declared, the role was bound, and the artifact held the
// value.
//
// REQ-35 fixes the invoked set as "exactly those readers serving an owned
// key some candidate row of the requested model requires". A guard the row
// cannot be decided without is a requirement. REQ-106 names the shape:
// refusal-first is the recovery path for a real failure, not a way to turn
// a decidable transition into one.
//
// The fixture is deliberately minimal: `read.flags` serves exactly one
// owned key, that key appears only in a guard, and its role IS bound on
// every invocation below — so no `flow-artifact-missing` and no
// narrowing-by-design (REQ-36) can explain the outcome.

const advGuardOwnedModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "advguardowned"
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

[tags.flag]
provenance = "owned"
kind = "bool"
single_valued = true
required = true

[read.state]
role = "state"
path = "flow.state"
keys = ["status"]
timeout = "2s"

[read.flags]
role = "flags"
path = "flow.flags"
keys = ["flag"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status"]
timeout = "2s"
read_back = true

[write.flags]
role = "flags"
path = "flow.flags"
keys = ["flag"]
timeout = "2s"
read_back = true

[initial]
status = "draft"
flag = "true"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "advance-draft"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "final"
`

const advFlagsRole = "flags"

// TestAdv2_AReaderServingOnlyAGuardsOwnedKeyIsStillInvoked proves REQ-35
// over a row whose owned demand is a GUARD input rather than a write
// target, and shows the consequence the FM `flow-guard-unevaluable` row
// otherwise hides.
func TestAdv2_AReaderServingOnlyAGuardsOwnedKeyIsStillInvoked(t *testing.T) {
	model := writeFlowModel(t, advGuardOwnedModel)
	stateArt := newFlowArtifact(t, "state.artifact")
	flagsArt := newFlowArtifact(t, "flags.artifact")
	stateBind := artifactBinding(flowStateRole, stateArt)
	flagsBind := artifactBinding(advFlagsRole, flagsArt)

	// Establish both owned facts through the production path, so the
	// artifacts hold what the guard reads.
	if _, _, err := runCmd(t, "flow", "set-state", "--model", model,
		"--artifact", stateBind, "--artifact", flagsBind,
		"--write", "status=draft", "--write", "flag=true",
		"--as=json"); err != nil {
		t.Fatalf("seeding both artifacts failed: %v", err)
	}

	t.Run("resolve", func(t *testing.T) {
		stdout, _, err := runCmd(t, "flow", "resolve", "--model", model,
			"--artifact", stateBind, "--artifact", flagsBind,
			"--outcome", "advance", "--as=json")
		if err != nil {
			t.Fatalf("resolve refused %q over an artifact that HOLDS "+
				"flag=true and a bound `flags` role: the reader serving the "+
				"guard's owned key was never invoked (REQ-35, REQ-106)\n%s",
				codeOfError(err), stdout)
		}

		data := flowData(t, stdout)
		readers, _ := stringsAt(data, "readers")
		if !containsString(readers, "flags") {
			t.Errorf("`readers` = %v; `read.flags` serves the owned key the "+
				"selected row's guard requires, so REQ-35 puts it in the "+
				"invoked set", readers)
		}
		owned, _ := data["owned"].(map[string]any)
		if owned["flag"] != "true" {
			t.Errorf("`owned` = %v; it must carry the guard's owned fact, "+
				"assembled from the declared reader (REQ-30, REQ-35)", owned)
		}
	})

	t.Run("next", func(t *testing.T) {
		data := flowData(t, requireSuccess(t, "flow", "next",
			"--model", model, "--artifact", stateBind,
			"--artifact", flagsBind, "--as=json"))

		readers, _ := stringsAt(data, "readers")
		if !containsString(readers, "flags") {
			t.Errorf("`readers` = %v; `read.flags` serves an owned key a "+
				"candidate row requires (REQ-35)", readers)
		}

		candidates, ok := objectsAt(data, "candidates")
		if !ok || len(candidates) == 0 {
			t.Fatalf("`candidates` is empty: %#v", data["candidates"])
		}
		for _, c := range candidates {
			unresolved, _ := stringsAt(c, "unresolved")
			if containsString(unresolved, "flag") {
				t.Errorf("candidate %v reports `flag` unresolved, but a "+
					"declared reader over a bound role holds it: `next` "+
					"reports a fact unresolved that the invoked set was "+
					"required to establish (REQ-35, REQ-44)", c["rule"])
			}
		}
	})
}

// --- ADV-3 --------------------------------------------------------------
//
// REQ-25 (TD, *Model selection*): "`revision` on every payload is the
// loaded model's own revision identity, carried through verbatim … The CLI
// never derives, hashes, or synthesizes it, so a model that declares none
// renders `revision` empty rather than a CLI-invented value."
//
// RDR 0002's `[model]` block admits `id`, `version`, `description`, and
// `metadata` — there is no revision field, so NO model can declare one and
// the clause's consequent governs every payload this CLI emits.
//
// The implementation renders `table.Model.KernelTable().Revision`, which
// RDR 0002 sets to the model ID. The ID is a different identity: REQ-110
// makes request identity "(… `--model` path AND model revision …)", so a
// revision that is only ever the id cannot distinguish two revisions of
// one model — which is precisely the discrimination the field is named
// for. Rendering the id there is a CLI-invented stand-in.
//
// The oracle is deliberately narrow. It does not demand a revision
// mechanism the RDR does not specify; it demands only the consequent
// REQ-25 states for a model that declares none, and it proves the value is
// the id rather than merely non-empty by asserting the exact equality.
func TestAdv3_AModelDeclaringNoRevisionRendersRevisionEmpty(t *testing.T) {
	// Two documents identical but for their `[model].id`, so a `revision`
	// tracking the id is visibly tracking the WRONG identity.
	const idA = "advrevone"
	const idB = "advrevtwo"

	verbs := [][]string{
		{"flow", "next"},
		{"flow", "resolve", "--outcome", "hold"},
		{"flow", "read-state"},
		{"flow", "set-state", "--write", "status=draft"},
	}

	for _, id := range []string{idA, idB} {
		src := strings.Replace(flowMVVModel,
			`id = "mvvflow"`, `id = "`+id+`"`, 1)
		model := writeFlowModel(t, src)
		art := seedArtifact(t, model, "status=draft")
		stateBind := artifactBinding(flowStateRole, art)

		for _, verb := range verbs {
			args := append([]string{}, verb...)
			args = append(args, "--model", model, "--artifact", stateBind)
			if verb[1] == "read-state" {
				// read-state runs every declared reader (REQ-37).
				args = append(args, "--artifact",
					artifactBinding(flowOrphanRole, art))
			}
			args = append(args, "--as=json")

			t.Run(id+"/"+verb[1], func(t *testing.T) {
				data := flowData(t, requireSuccess(t, args...))

				raw, ok := data["revision"]
				if !ok {
					t.Fatalf("payload carries no `revision`; every payload "+
						"carries it (REQ-74..REQ-78). keys = %v", keysOf(data))
				}
				rev, _ := raw.(string)

				if rev == id {
					t.Errorf("`revision` = %q, which is the `[model].id`. "+
						"RDR 0002's `[model]` schema declares no revision "+
						"field, so this model declares none, and REQ-25 "+
						"requires `revision` to render EMPTY rather than a "+
						"CLI-invented value. Standing in the id also makes "+
						"REQ-110's request identity — which names the model "+
						"revision separately from the model selection — "+
						"unable to tell two revisions of one model apart",
						rev)
					return
				}
				if rev != "" {
					t.Errorf("`revision` = %q; this model declares no "+
						"revision, so REQ-25 renders it empty", rev)
				}
			})
		}
	}
}

// --- shared helpers -----------------------------------------------------

// codeOfError reads the stable `flow-*` code off a refusal returned by the
// production path, without importing the concrete error type's fields into
// every assertion above.
func codeOfError(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	if i := strings.Index(msg, ":"); i > 0 {
		return msg[:i]
	}
	return msg
}

func containsString(list []string, want string) bool {
	for _, s := range list {
		if s == want {
			return true
		}
	}
	return false
}

// --- FAIL-1 (Phase 3a) --------------------------------------------------
//
// REQ-50: "all gates on the row MUST run, deny MUST override allow and
// indeterminate, and every gate result MUST be reported."
// REQ-47: "no gate result is ever dropped silently" — binding in BOTH
// verbs.
// Technical Design, *Gates*: "every gate's result is reported, as `gates[]`
// on a success or as ONE `findings[]` ENTRY PER GATE on `flow-gate-denied`
// / `flow-gate-indeterminate`."
//
// `gateVerdictFailure` partitioned the row's results into `denied` and
// `indeterminate` and handed `gateRefusal` only the matching partition, so
// a gate that ALLOWED was in neither slice and vanished from the envelope
// entirely — absent from `findings[]`, from `detail`, and from text mode.
//
// The oracle below is non-circular because it asks the OTHER verb the same
// question over the SAME row: `flow next --evaluate-gates` reports the
// allowing gate on the candidate's `gates[]`. Two verbs disagreeing about
// what the same gates on the same row answered is the drop, stated without
// reference to either implementation.

// advGateReportModel pairs an ALLOWING gate with a refusing one on each of
// two rows, so both refusal arms — deny and indeterminate — carry a gate
// whose result the partitioning dropped.
const advGateReportModel = `outcomes = ["deny-path", "indeterminate-path"]
terminal = ["done"]

[model]
id = "advgatereport"
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
path = "flow.state"
keys = ["status"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status"]
timeout = "2s"
read_back = true

[gate.permits]
role = "state"
path = "flow.gate.allow"
keys = ["status"]
timeout = "2s"

[gate.refuses]
role = "state"
path = "flow.gate.deny"
keys = ["status"]
timeout = "2s"

[gate.undecided]
role = "state"
path = "flow.gate.indeterminate"
keys = ["status"]
timeout = "2s"

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "deny-row"
gate = ["permits", "refuses"]
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "deny-path"
[rule.write]
status = "final"

[[rule]]
id = "indeterminate-row"
gate = ["permits", "undecided"]
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "indeterminate-path"
[rule.write]
status = "final"
`

// TestFail1_EveryGateOnTheSelectedRowIsReportedOnAGateRefusal proves
// REQ-50's "every gate result MUST be reported" on the REFUSAL path, where
// the success path's `gates[]` is not available to carry it.
func TestFail1_EveryGateOnTheSelectedRowIsReportedOnAGateRefusal(t *testing.T) {
	model := writeFlowModel(t, advGateReportModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	cases := []struct {
		name    string
		outcome string
		code    string
		// gates names every gate on the row, allowing gates included.
		gates []string
	}{
		{
			name:    "deny-arm",
			outcome: "deny-path",
			code:    "flow-gate-denied",
			gates:   []string{"permits", "refuses"},
		},
		{
			name:    "indeterminate-arm",
			outcome: "indeterminate-path",
			code:    "flow-gate-indeterminate",
			gates:   []string{"permits", "undecided"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, _, err := runCmd(t, "flow", "resolve", "--model", model,
				"--artifact", bind, "--outcome", tc.outcome, "--as=json")
			if err == nil {
				t.Fatalf("a refusing gate produced a plan:\n%s", stdout)
			}
			var ce *clierr.CLIError
			if !asCLIError(err, &ce) {
				t.Fatalf("not a CLIError: %v", err)
			}
			// Precedence is unchanged: the ENVELOPE code is still the
			// disposition the row implies (REQ-50, REQ-51).
			if ce.Code != tc.code {
				t.Fatalf("code = %q; want %q", ce.Code, tc.code)
			}

			var reported []string
			for _, f := range flowFailureFindings(t, stdout) {
				reported = append(reported, f.Param)
			}
			for _, id := range tc.gates {
				if !containsString(reported, id) {
					t.Errorf("`findings[]` = %v drops the gate %q. REQ-50 "+
						"requires that EVERY gate result be reported and "+
						"REQ-47 that no gate result be dropped silently; "+
						"the Technical Design fixes the refusal-path "+
						"carrier as ONE findings[] entry per gate on the "+
						"selected row", reported, id)
				}
			}
			if len(reported) != len(tc.gates) {
				t.Errorf("`findings[]` carries %d entries for a row with %d "+
					"gates; the carrier is one entry PER GATE",
					len(reported), len(tc.gates))
			}

			// The two verbs must agree about what the same gates on the
			// same row answered. `next --evaluate-gates` is the independent
			// witness that the allowing gate RAN and returned `allow`.
			data := flowData(t, requireSuccess(t, "flow", "next",
				"--model", model, "--artifact", bind,
				"--evaluate-gates", "--as=json"))
			candidates, ok := objectsAt(data, "candidates")
			if !ok {
				t.Fatalf("`candidates` is not an array of objects: %#v",
					data["candidates"])
			}
			var witnessed []string
			for _, c := range candidates {
				if c["outcome"] != tc.outcome {
					continue
				}
				gates, _ := objectsAt(c, "gates")
				for _, g := range gates {
					id, _ := g["id"].(string)
					witnessed = append(witnessed, id)
				}
			}
			for _, id := range witnessed {
				if !containsString(reported, id) {
					t.Errorf("`flow next --evaluate-gates` reports the gate "+
						"%q on this row and `flow resolve` does not: the "+
						"two verbs disagree about what the same gates on "+
						"the same row answered (REQ-47, REQ-50). next = %v, "+
						"resolve findings = %v", id, witnessed, reported)
				}
			}
		})
	}
}

// REQ-36: "A reader no candidate row needs MUST NOT run, and its unbound
// artifact role MUST NOT raise flow-artifact-missing; flow-artifact-missing
// is scoped to the roles the invoked set needs."
// BOUNDARY — an escape row bound to a DIFFERENT outcome is unrescuable by
// construction: the kernel skips it on `row.Outcome != in.Recognized`
// (`internal/resolve.escapeOrRefuse`) before it can rescue anything. Its
// guard-owned keys are therefore not demands of the requested outcome, so
// the reader serving only them must not be INVOKED and its unbound role
// must not raise `flow-artifact-missing`. Retaining such a row was
// harmless until guard-owned keys joined the demand set (DEV-8).
//
// The kernel may still refuse for its own reasons once it evaluates that
// row's guard; this oracle pins only the CLI's reader-narrowing duty,
// which is the clause REQ-35/REQ-36 assign to this RDR.
func TestFi_AnEscapeRowUnderAnotherOutcomeDoesNotDemandItsReader(t *testing.T) {
	model := writeFlowModel(t, flowEscapeOtherOutcomeModel)
	art := seedArtifact(t, model, "status=draft")

	// The `bail` escape row's `sidecar` reader is deliberately left
	// UNBOUND. `flow-artifact-missing` is the refusal that fires when the
	// INVOKED set needs an unbound role, so its absence is what proves the
	// narrowing (REQ-36).
	_, _, err := runCmd(t, "flow", "resolve", "--model", model,
		"--outcome", "advance",
		"--artifact", artifactBinding(flowStateRole, art),
		"--as=json")

	var ce *clierr.CLIError
	if asCLIError(err, &ce) && ce.Code == "flow-artifact-missing" {
		t.Fatalf("resolve raised %q for the unbound `sidecar` role; that "+
			"reader serves only a guard on an escape row bound to `bail`, "+
			"which the kernel can never rescue `advance` with, so it is "+
			"not in the invoked set and its unbound role MUST NOT refuse "+
			"this request (REQ-35, REQ-36)", ce.Code)
	}
}
