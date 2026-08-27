package cli

// RDR 0005 — the Minimum Viable Validation (`0005:MVV`), run end to end
// through the production Cobra path.
//
// The MVV fixes six obligations this file discharges as named subtests:
//
//  1. all four verbs prove out over ONE fixture-backed flow;
//  2. `set-state` persists one scalar write, one set write as a JSON array,
//     and one `--clear`, then read-back-verifies them;
//  3. the set write's members include one containing `<` and one containing
//     `&`, asserted BYTE-IDENTICAL through plan -> request -> read-back —
//     "so the encoder rule is covered by the validation rather than only by
//     review". A green exit code is explicitly NOT sufficient here;
//  4. the happy path, one escaped plan, one gate deny, one kernel refusal,
//     and one exit-3 accessor failure each run in `--as=text` AND
//     `--as=json`, "asserting `findings[]` where the table names it";
//  5. the unbound-reader pair: `next` / `resolve` succeed without invoking
//     it, `read-state` invokes it and refuses the missing binding;
//  6. `next --evaluate-gates` over one gated candidate plus one
//     guard-excluded row reports the gate on the reported candidate, runs
//     no gate for the excluded row, and exits 0 even on a deny.

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
)

// REQ-MVV / `0005:MVV`: "Implement one fixture-backed flow and prove all
// four verbs through the production Cobra path".
// REQ-122: "Command tests exercise the `flow` command group through the
// production Cobra path rather than calling renderers or kernel functions
// directly."
// HAPPY PATH
func TestMVV_AllFourVerbsProveOutOverOneFixtureBackedFlow(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	stateBind := artifactBinding(flowStateRole, art)
	orphanBind := artifactBinding(flowOrphanRole, art)

	// `flow next` returns the legal outcome alphabet with gate ids as
	// unknown facts.
	t.Run("next", func(t *testing.T) {
		data := flowData(t, requireSuccess(t, "flow", "next",
			"--model", model, "--artifact", stateBind, "--as=json"))

		outcomes, ok := stringsAt(data, "outcomes")
		if !ok {
			t.Fatalf("`outcomes` is not an array: %#v", data["outcomes"])
		}
		for _, want := range []string{"advance", "hold", "bail"} {
			if !slices.Contains(outcomes, want) {
				t.Errorf("`outcomes` = %v omits the declared outcome %q",
					outcomes, want)
			}
		}

		// Gate ids appear as UNKNOWN facts, because no gate ran.
		candidates, _ := objectsAt(data, "candidates")
		var gateListed bool
		for _, c := range candidates {
			unknown, ok := unknownAt(c, "unknown")
			if !ok {
				t.Fatalf("candidate %v carries no `unknown` list of "+
					"{key, reason} pairs: %#v", c["rule"], c["unknown"])
			}
			if hasUnknown(unknown, "approval", "not-evaluated") {
				gateListed = true
			}
		}
		if !gateListed {
			t.Error("no candidate lists the gate id `approval` as an " +
				"unknown fact {approval, not-evaluated}; without " +
				"--evaluate-gates the gate does not run and its id is what " +
				"the caller sees")
		}
	})

	// `flow resolve` reads owned state from the fixture artifact, maps one
	// outcome to one plan, and runs one gate on the selected row.
	t.Run("resolve", func(t *testing.T) {
		gated := seedArtifact(t, model, "status=draft", "stale=obsolete")
		data := flowData(t, requireSuccess(t, "flow", "resolve",
			"--model", model,
			"--artifact", artifactBinding(flowStateRole, gated),
			"--outcome", "advance", "--as=json"))

		if data["rule"] != "advance-draft" {
			t.Errorf("`rule` = %#v; one outcome maps to ONE plan on the row "+
				"`advance-draft`", data["rule"])
		}
		owned, ok := data["owned"].(map[string]any)
		if !ok || owned["status"] != "draft" {
			t.Errorf("`owned` = %#v; owned state is read from the fixture "+
				"artifact", data["owned"])
		}
		gates, ok := objectsAt(data, "gates")
		if !ok || len(gates) != 1 {
			t.Fatalf("`gates[]` = %#v; exactly ONE gate rides the selected "+
				"row", data["gates"])
		}
		if gates[0]["id"] != "approval" {
			t.Errorf("gate id = %#v; want %q", gates[0]["id"], "approval")
		}
	})

	// `flow read-state` reads the fixture artifact tags per reader.
	t.Run("read-state", func(t *testing.T) {
		data := flowData(t, requireSuccess(t, "flow", "read-state",
			"--model", model, "--artifact", stateBind,
			"--artifact", orphanBind, "--as=json"))

		readers, ok := objectsAt(data, "readers")
		if !ok || len(readers) == 0 {
			t.Fatalf("`readers` is empty or malformed: %#v", data["readers"])
		}
		if got := readerTagValue(t, data, "state", "status"); got != "draft" {
			t.Errorf("reader `state` reports status = %#v; want %q",
				got, "draft")
		}
	})

	// `flow set-state` persists one scalar write, one set write as a JSON
	// array, and one `--clear`, then read-back-verifies them.
	t.Run("set-state", func(t *testing.T) {
		target := seedArtifact(t, model, "status=draft", "stale=obsolete")
		bind := artifactBinding(flowStateRole, target)

		data := flowData(t, requireSuccess(t, "flow", "set-state",
			"--model", model, "--artifact", bind,
			"--write", "status=final",
			"--write", "labels="+canonicalSetLiteral,
			"--clear", "stale", "--as=json"))

		owned, ok := data["owned"].(map[string]any)
		if !ok {
			t.Fatalf("`owned` is not an object: %#v", data["owned"])
		}
		if owned["status"] != "final" {
			t.Errorf("the scalar write was not read-back-verified: "+
				"owned[status] = %#v", owned["status"])
		}
		assertCanonicalSetValue(t, "owned[labels]", owned["labels"])
		if v, held := owned["stale"]; held {
			t.Errorf("the cleared key reads back as %#v; a cleared key reads "+
				"back ABSENT", v)
		}
	})
}

// REQ-MVV: "The set write's members MUST include one containing `<` and one
// containing `&`, asserted byte-identical through plan → request →
// read-back, so the encoder rule is covered by the validation rather than
// only by review."
// REQ-70 / REQ-71 / REQ-108: the canonical form and its byte equality.
// ADVERSARIAL — the three hops are compared as BYTES, not as decoded
// values: a green exit code is explicitly not sufficient, because fidelity
// loss hides behind a passing run.
func TestMVV_SetMembersWithAngleAndAmpersandSurviveAllThreeHopsByteIdentical(t *testing.T) {
	model := writeFlowModel(t, flowSetPlanModel)
	art := seedArtifact(t, model, "status=draft", "labels="+canonicalSetLiteral)
	bind := artifactBinding(flowStateRole, art)

	// --- HOP 1: the PLAN emits the set value ---------------------------
	planOut := requireSuccess(t, "flow", "resolve", "--model", model,
		"--artifact", bind, "--outcome", "advance", "--as=json")
	assertNoHTMLEscapes(t, "plan (hop 1)", planOut)

	plan := flowData(t, planOut)
	writes, ok := plan["writes"].(map[string]any)
	if !ok {
		t.Fatalf("the plan's `writes` is not an object: %#v", plan["writes"])
	}
	planLiteral := rawLiteralFor(t, planOut, "labels")
	if planLiteral == "" {
		t.Fatalf("the plan does not carry `labels` as a set value:\n%s",
			planOut)
	}
	assertCanonicalSetValue(t, "plan writes[labels]", writes["labels"])

	// --- HOP 2: the REQUEST carries it back verbatim -------------------
	// Copy-through: the plan's own emitted value becomes the `--write`
	// argument with no re-encoding by the test (REQ-108).
	requestLiteral := planValueLiteral(t, writes["labels"])
	if requestLiteral != canonicalSetLiteral {
		t.Errorf("the request literal reconstructed from the plan is %q; "+
			"want %q — the plan's value is accepted verbatim by "+
			"`set-state`'s grammar", requestLiteral, canonicalSetLiteral)
	}

	target := seedArtifact(t, model, "status=draft")
	targetBind := artifactBinding(flowStateRole, target)
	requestOut := requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", targetBind, "--write", "labels="+requestLiteral,
		"--as=json")
	assertNoHTMLEscapes(t, "request echo (hop 2)", requestOut)
	requestEcho := rawLiteralFor(t, requestOut, "labels")

	// --- HOP 3: the READ-BACK returns it -------------------------------
	readOut := requireSuccess(t, "flow", "read-state", "--model", model,
		"--artifact", targetBind, "--as=json")
	assertNoHTMLEscapes(t, "read-back (hop 3)", readOut)
	readBack := rawLiteralFor(t, readOut, "labels")

	// The three hops agree BYTE-FOR-BYTE.
	if planLiteral != requestEcho {
		t.Errorf("plan -> request lost fidelity:\n  plan:    %s\n  request: %s",
			planLiteral, requestEcho)
	}
	if requestEcho != readBack {
		t.Errorf("request -> read-back lost fidelity:\n  request:   %s\n"+
			"  read-back: %s", requestEcho, readBack)
	}
	if readBack != canonicalSetLiteral {
		t.Errorf("the value read back is %s; want the canonical %s — "+
			"members sorted, duplicate-free, compact, and `<` / `&` "+
			"rendered as THEMSELVES", readBack, canonicalSetLiteral)
	}

	// And the two required members are literally present, unescaped, at
	// every hop — the assertion the MVV says review alone cannot give.
	for _, hop := range []struct{ name, out string }{
		{"plan", planOut}, {"request", requestOut}, {"read-back", readOut},
	} {
		for _, member := range []string{setMemberLT, setMemberAmp} {
			if !strings.Contains(hop.out, member) {
				t.Errorf("hop %q does not carry the member %q as itself:\n%s",
					hop.name, member, hop.out)
			}
		}
	}
}

// REQ-MVV: "Run the happy path, one escaped plan, one gate deny, one kernel
// refusal, and one exit-3 accessor failure in `--as=text` and `--as=json`,
// asserting `findings[]` where the table names it."
// REQ-120: "Derive both renderings from one typed result per verb and test
// both modes."
// DOMAIN EDGE — five scenarios × two modes, with the exit and the refusal
// identity required to AGREE across modes.
func TestMVV_FiveScenariosInBothModesAgreeOnExitAndIdentity(t *testing.T) {
	mvv := writeFlowModel(t, flowMVVModel)
	escape := writeFlowModel(t, flowEscapeModel)
	deny := writeFlowModel(t, flowGateDenyModel)
	fail := writeFlowModel(t, flowGateFailModel)

	mvvArt := seedArtifact(t, mvv, "status=draft")
	escapeArt := seedArtifact(t, escape, "status=draft")
	denyArt := seedArtifact(t, deny, "status=draft")
	failArt := seedArtifact(t, fail, "status=draft")

	scenarios := []struct {
		name string
		args []string
		// code is "" for the two success scenarios.
		code string
		exit int
		// findings records whether the code table names `findings[]` as
		// this refusal's carrier.
		findings bool
	}{
		{
			name: "happy-path",
			args: []string{"flow", "resolve", "--model", mvv,
				"--artifact", artifactBinding(flowStateRole, mvvArt),
				"--outcome", "hold"},
			exit: 0,
		},
		{
			name: "escaped-plan",
			args: []string{"flow", "resolve", "--model", escape,
				"--artifact", artifactBinding(flowStateRole, escapeArt),
				"--outcome", "bail"},
			exit: 0,
		},
		{
			name: "gate-deny",
			args: []string{"flow", "resolve", "--model", deny,
				"--artifact", artifactBinding(flowStateRole, denyArt),
				"--outcome", "deny-path"},
			code: "flow-gate-denied", exit: 2, findings: true,
		},
		{
			name: "kernel-refusal",
			args: []string{"flow", "resolve", "--model", mvv,
				"--artifact", artifactBinding(flowStateRole, mvvArt),
				"--outcome", "bail"},
			code: "flow-no-match", exit: 2, findings: true,
		},
		{
			name: "exit-3-accessor-failure",
			args: []string{"flow", "resolve", "--model", fail,
				"--artifact", artifactBinding(flowStateRole, failArt),
				"--outcome", "advance"},
			code: "flow-accessor-failed", exit: 3,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			jsonOut, _, jerr := runCmd(t, append(slices.Clone(sc.args), "--as=json")...)
			textOut, textErr, terr := runCmd(t, slices.Clone(sc.args)...)

			// Both modes agree on the EXIT.
			jExit, tExit := clierr.ExitCodeFor(jerr), clierr.ExitCodeFor(terr)
			if jExit != sc.exit {
				t.Errorf("--as=json exit = %d; want %d\nstdout:\n%s",
					jExit, sc.exit, jsonOut)
			}
			if tExit != sc.exit {
				t.Errorf("--as=text exit = %d; want %d\nstdout:\n%s\nstderr:\n%s",
					tExit, sc.exit, textOut, textErr)
			}

			if sc.code == "" {
				// A success in both modes emits a payload in both modes.
				if strings.TrimSpace(jsonOut) == "" {
					t.Error("--as=json emitted no success envelope")
				}
				if strings.TrimSpace(textOut) == "" {
					t.Error("--as=text emitted no success payload")
				}
				if sc.name == "escaped-plan" {
					data := flowData(t, jsonOut)
					if data["escaped"] != true {
						t.Errorf("`escaped` = %#v; an escaped plan is a "+
							"SUCCESS carrying escaped=true", data["escaped"])
					}
					if _, ok := data["escape_class"].(string); !ok {
						t.Errorf("an escaped plan carries no `escape_class`: "+
							"%#v", data["escape_class"])
					}
				}
				return
			}

			// Both modes agree on the refusal IDENTITY.
			var jce, tce *clierr.CLIError
			if !asCLIError(jerr, &jce) {
				t.Fatalf("--as=json refusal is not a CLIError: %v", jerr)
			}
			if !asCLIError(terr, &tce) {
				t.Fatalf("--as=text refusal is not a CLIError: %v", terr)
			}
			if jce.Code != sc.code {
				t.Errorf("--as=json code = %q; want %q", jce.Code, sc.code)
			}
			if tce.Code != jce.Code {
				t.Errorf("--as=text code = %q but --as=json code = %q; the "+
					"two renderings derive from ONE typed result",
					tce.Code, jce.Code)
			}

			if !sc.findings {
				return
			}
			// `findings[]` where the table names it — in BOTH modes.
			findings := flowFailureFindings(t, jsonOut)
			if len(findings) == 0 {
				t.Fatalf("%s carries no `findings[]`, which the code table "+
					"names as its carrier", sc.code)
			}
			for _, f := range findings {
				if !strings.Contains(textErr, f.Code) {
					t.Errorf("--as=text omits the finding code %q that "+
						"--as=json reports\nstderr:\n%s", f.Code, textErr)
				}
			}
		})
	}
}

// REQ-MVV: "The fixture model MUST also declare one reader no candidate row
// needs, with its artifact role left unbound: `flow next` and `flow resolve`
// MUST succeed without invoking it or raising `flow-artifact-missing`, while
// `flow read-state` on the same model MUST invoke it and refuse the missing
// binding. That pair is the only assertion that distinguishes the narrowed
// invoked set from \"every declared reader\"."
// REQ-35 / REQ-36 / REQ-37, DEV-1 reading (a).
// ADVERSARIAL — the SAME argv succeeds on two verbs and refuses on the
// third, which no "run every declared reader" implementation can satisfy.
func TestMVV_TheUnboundReaderPairDistinguishesTheNarrowedInvokedSet(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	// ONLY the `state` role is bound throughout.
	bind := artifactBinding(flowStateRole, art)

	t.Run("next-succeeds-without-invoking-it", func(t *testing.T) {
		stdout, _, err := runCmd(t, "flow", "next", "--model", model,
			"--artifact", bind, "--as=json")
		if err != nil {
			var ce *clierr.CLIError
			if asCLIError(err, &ce) && ce.Code == "flow-artifact-missing" {
				t.Fatalf("`next` raised flow-artifact-missing for the role "+
					"%q of a reader NO candidate row needs; "+
					"flow-artifact-missing is scoped to the roles the "+
					"INVOKED set needs", ce.Param)
			}
			t.Fatalf("`next` failed: %v", err)
		}
		assertOrphanReaderNotInvoked(t, flowData(t, stdout))
	})

	t.Run("resolve-succeeds-without-invoking-it", func(t *testing.T) {
		stdout, _, err := runCmd(t, "flow", "resolve", "--model", model,
			"--artifact", bind, "--outcome", "hold", "--as=json")
		if err != nil {
			var ce *clierr.CLIError
			if asCLIError(err, &ce) && ce.Code == "flow-artifact-missing" {
				t.Fatalf("`resolve` raised flow-artifact-missing for the "+
					"role %q of a reader no candidate row needs", ce.Param)
			}
			t.Fatalf("`resolve` failed: %v", err)
		}
		assertOrphanReaderNotInvoked(t, flowData(t, stdout))
	})

	t.Run("read-state-invokes-it-and-refuses-the-missing-binding", func(t *testing.T) {
		ce := requireRefusal(t, "flow-artifact-missing", 2,
			"flow", "read-state", "--model", model,
			"--artifact", bind, "--as=json")
		if ce.Param != flowOrphanRole {
			t.Errorf("param = %q; want the unbound role %q — `read-state` "+
				"has no candidate set to narrow by, so it runs EVERY "+
				"declared reader", ce.Param, flowOrphanRole)
		}
	})
}

// REQ-MVV: "`flow next --evaluate-gates` over a model with one gated
// candidate and one guard-excluded row MUST report the gate result on the
// reported candidate, run no gate for the excluded row, and exit 0 even
// when that gate denies."
// REQ-42 / REQ-43 / REQ-46.
// ADVERSARIAL
func TestMVV_NextEvaluateGatesReportsOneGateAndExitsZeroEvenOnADeny(t *testing.T) {
	model := writeFlowModel(t, flowGatedNextModel)
	art := seedArtifact(t, model, "status=draft", "flag=true")

	stdout, _, err := runCmd(t, "flow", "next", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--evaluate-gates", "--as=json")
	if err != nil {
		t.Fatalf("`next --evaluate-gates` exited non-zero: %v\n"+
			"`next` enumerates rather than selects: a deny constrains only "+
			"that candidate and is reported, never refused\nstdout:\n%s",
			err, stdout)
	}

	data := flowData(t, stdout)
	candidates, ok := objectsAt(data, "candidates")
	if !ok {
		t.Fatalf("`candidates` is not an array of objects: %#v",
			data["candidates"])
	}

	var totalGates int
	var reportedHasItsGate bool
	for _, c := range candidates {
		if c["rule"] == "gated-excluded" {
			t.Error("the guard-excluded row `gated-excluded` was reported " +
				"as a candidate; its guard excludes it, so it is not a " +
				"candidate and its gate MUST NOT run")
		}
		gates, _ := objectsAt(c, "gates")
		totalGates += len(gates)
		for _, g := range gates {
			if g["id"] == "excluded" {
				t.Error("a gate result for the EXCLUDED row's gate " +
					"`excluded` was reported; that gate must not run at all")
			}
			if c["rule"] == "gated-reported" && g["id"] == "reported" {
				reportedHasItsGate = true
			}
		}
	}

	if !reportedHasItsGate {
		t.Error("the reported candidate `gated-reported` does not carry its " +
			"own gate `reported`'s result; each gate result MUST be " +
			"reported on the candidate that carries it")
	}
	if totalGates != 1 {
		t.Errorf("%d gate results were reported across all candidates; "+
			"EXACTLY one gate runs — the reported candidate's — and no "+
			"others", totalGates)
	}
}

// REQ-132: "The four illustrative invocation shapes in `Illustrative Code`
// must parse under the shipped flag grammar."
// BOUNDARY — the shapes exercise `--model`, `--flow`, `--artifact`,
// `--tag`, `--outcome`, `--write` scalar, `--write` set literal, `--clear`,
// and `--as=json`. They are non-normative as OUTPUT, so only the grammar is
// asserted: no shape may fail with a flag-parse error.
func TestReq132_TheIllustrativeInvocationShapesParseUnderTheShippedGrammar(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	shapes := [][]string{
		{"flow", "next", "--model", model, "--artifact", bind,
			"--tag", "profile=mid", "--as=json"},
		{"flow", "resolve", "--model", model, "--artifact", bind,
			"--tag", "profile=mid", "--outcome", "hold"},
		{"flow", "read-state", "--flow", "mvvflow", "--artifact", bind},
		{"flow", "set-state", "--model", model, "--artifact", bind,
			"--write", "status=final",
			"--write", `labels=["cli","final"]`,
			"--clear", "stale"},
	}

	for i, shape := range shapes {
		t.Run(strings.Join(shape[:2], "-"), func(t *testing.T) {
			_, _, err := runCmd(t, shape...)
			if err == nil {
				return
			}
			var ce *clierr.CLIError
			if !asCLIError(err, &ce) {
				t.Fatalf("shape %d produced a non-CLIError: %v", i, err)
			}
			// A grammar failure is what this REQ forbids. A DOMAIN refusal
			// (an unresolvable `--flow` id, say) is fine — the shapes are
			// non-normative as output.
			switch ce.Code {
			case "command-error", "flag-invalid-value":
				t.Errorf("shape %d does not PARSE under the shipped flag "+
					"grammar: %s: %s\nargv: %v", i, ce.Code, ce.Message, shape)
			}
		})
	}
}

// rawLiteralFor extracts the RAW JSON bytes a payload carries for key,
// searched anywhere in the emitted line. It compares hops as BYTES rather
// than as decoded values, which is what "byte-identical" requires: decoding
// both sides would launder exactly the escaping defect under test.
func rawLiteralFor(t *testing.T, out, key string) string {
	t.Helper()

	var walk func(raw json.RawMessage) string
	walk = func(raw json.RawMessage) string {
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(raw, &obj); err == nil {
			if v, held := obj[key]; held {
				s := strings.TrimSpace(string(v))
				// A set value rides either as a JSON array or as a string
				// holding the canonical literal. Unwrap the string case so
				// both render the same comparable bytes.
				if strings.HasPrefix(s, `"`) {
					var unquoted string
					if err := json.Unmarshal(v, &unquoted); err == nil {
						return unquoted
					}
				}
				return s
			}
			for _, v := range obj {
				if found := walk(v); found != "" {
					return found
				}
			}
			return ""
		}
		var arr []json.RawMessage
		if err := json.Unmarshal(raw, &arr); err == nil {
			for _, e := range arr {
				if found := walk(e); found != "" {
					return found
				}
			}
		}
		return ""
	}

	for _, line := range nonEmptyLines(out) {
		if found := walk(json.RawMessage(line)); found != "" {
			return found
		}
	}
	return ""
}
