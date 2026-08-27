package cli

// RDR 0005 — `flow resolve` (REQ-48..REQ-57, REQ-76), the kernel refusal
// mapping (REQ-4, REQ-49, REQ-92..REQ-98), and the post-selection gate site
// (REQ-50, REQ-51, REQ-55, REQ-116, REQ-119).
//
// The load-bearing distinction across this file: `resolve` SELECTS. Where
// `flow next` reports a gate result and exits 0, `resolve` turns a deny into
// a typed refusal — and never into a plan and never into an escape class.

import (
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/resolve"
)

// REQ-48: "flow resolve MUST return exactly one plan or exactly one CLIError
// refusal."
// REQ-53: "flow resolve MUST NOT print directly, initiate skill work, or
// choose among multiple matching rows."
// HAPPY PATH — exactly one terminal record, never both dispositions.
func TestReq48And53_ResolveReturnsExactlyOnePlanOrExactlyOneRefusal(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	t.Run("plan", func(t *testing.T) {
		stdout := requireSuccess(t, "flow", "resolve", "--model", model,
			"--artifact", bind, "--outcome", "hold", "--as=json")
		if n := len(nonEmptyLines(stdout)); n != 1 {
			t.Fatalf("stdout carries %d lines; resolve returns EXACTLY one "+
				"terminal record — a direct print would add another:\n%s",
				n, stdout)
		}
		data := flowData(t, stdout)
		if data["rule"] == nil || data["rule"] == "" {
			t.Error("the plan names no matched `rule`")
		}
	})

	t.Run("refusal", func(t *testing.T) {
		stdout, _, err := runCmd(t, "flow", "resolve", "--model", model,
			"--artifact", bind, "--outcome", "not-an-outcome", "--as=json")
		if err == nil {
			t.Fatalf("an unmodeled outcome produced a plan:\n%s", stdout)
		}
		if n := len(nonEmptyLines(stdout)); n != 1 {
			t.Fatalf("stdout carries %d lines; resolve returns EXACTLY one "+
				"terminal record:\n%s", n, stdout)
		}
	})
}

// REQ-4: "Kernel-mapped codes mirror the kernel refusal kinds one-to-one
// (`flow-<kind>`, underscores to hyphens)."
// REQ-49: "Kernel refusals MUST map one-to-one onto flow-unmodeled-outcome,
// flow-no-match, flow-ambiguous-match, flow-owned-state-unavailable, and
// flow-guard-unevaluable."
// REQ-92 / REQ-93 / REQ-94 / REQ-95 / REQ-96: each row's group, exit, and
// carrier.
// DOMAIN EDGE
func TestReq4And49_EveryKernelKindMapsToItsMirroredCodeOneToOne(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	ambiguous := writeFlowModel(t, flowAmbiguousModel)
	art := seedArtifact(t, model, "status=draft")
	ambigArt := seedArtifact(t, ambiguous, "status=draft")

	t.Run("unmodeled_outcome", func(t *testing.T) {
		ce := requireRefusal(t, "flow-unmodeled-outcome", 2,
			"flow", "resolve", "--model", model,
			"--artifact", artifactBinding(flowStateRole, art),
			"--outcome", "not-an-outcome", "--as=json")
		// REQ-92 fixes `param` as the carrier: the offending outcome.
		if ce.Param == "" {
			t.Error("`flow-unmodeled-outcome` carries no `param`; the code " +
				"table fixes `param` as its carrier")
		}
	})

	t.Run("no_match", func(t *testing.T) {
		// `bail` is a declared outcome with no matching row in the MVV model.
		stdout, _, err := runCmd(t, "flow", "resolve", "--model", model,
			"--artifact", artifactBinding(flowStateRole, art),
			"--outcome", "bail", "--as=json")
		if err == nil {
			t.Fatalf("`bail` produced a plan; no row matches it:\n%s", stdout)
		}
		var ce *clierr.CLIError
		if !asCLIError(err, &ce) {
			t.Fatalf("not a CLIError: %v", err)
		}
		if ce.Code != "flow-no-match" {
			t.Fatalf("code = %q; want %q", ce.Code, "flow-no-match")
		}
		// REQ-93 fixes `findings[]` (rows considered) as the carrier.
		if len(flowFailureFindings(t, stdout)) == 0 {
			t.Error("`flow-no-match` carried no findings; the code table " +
				"fixes `findings[]` (rows considered) as its carrier")
		}
	})

	t.Run("ambiguous_match", func(t *testing.T) {
		stdout, _, err := runCmd(t, "flow", "resolve", "--model", ambiguous,
			"--artifact", artifactBinding(flowStateRole, ambigArt),
			"--outcome", "advance", "--as=json")
		if err == nil {
			t.Fatalf("two matching rows produced a plan; resolve MUST NOT "+
				"choose among them:\n%s", stdout)
		}
		var ce *clierr.CLIError
		if !asCLIError(err, &ce) {
			t.Fatalf("not a CLIError: %v", err)
		}
		if ce.Code != "flow-ambiguous-match" {
			t.Fatalf("code = %q; want %q", ce.Code, "flow-ambiguous-match")
		}
		// REQ-94: `findings[]` (matching rows) — BOTH conflicting rows are
		// named, because a tie-break would be exactly the forbidden choice.
		findings := flowFailureFindings(t, stdout)
		if len(findings) < 2 {
			t.Fatalf("`flow-ambiguous-match` carried %d findings; both "+
				"matching rows must be named — reporting one would be the "+
				"tie-break REQ-53 forbids", len(findings))
		}
		var named []string
		for _, f := range findings {
			named = append(named, f.Rule)
		}
		for _, want := range []string{"first", "second"} {
			if !slices.Contains(named, want) {
				t.Errorf("findings name rows %v; the conflicting row %q is "+
					"missing", named, want)
			}
		}
	})

	t.Run("owned_state_unavailable", func(t *testing.T) {
		// `advance` requires owned `stale` and `status`. With only `status`
		// established, the kernel refuses on the absent owned key.
		stdout, _, err := runCmd(t, "flow", "resolve", "--model", model,
			"--artifact", artifactBinding(flowStateRole, art),
			"--outcome", "advance", "--as=json")
		if err == nil {
			t.Fatalf("`advance` produced a plan though an owned key is "+
				"unavailable:\n%s", stdout)
		}
		var ce *clierr.CLIError
		if !asCLIError(err, &ce) {
			t.Fatalf("not a CLIError: %v", err)
		}
		if ce.Code != "flow-owned-state-unavailable" {
			t.Fatalf("code = %q; want %q", ce.Code,
				"flow-owned-state-unavailable")
		}
		// REQ-95: `findings[]` (keys).
		findings := flowFailureFindings(t, stdout)
		if len(findings) == 0 {
			t.Fatal("`flow-owned-state-unavailable` carried no findings; " +
				"the code table fixes `findings[]` (keys) as its carrier")
		}
		var keys []string
		for _, f := range findings {
			keys = append(keys, f.Key)
		}
		if !slices.Contains(keys, "stale") {
			t.Errorf("findings name keys %v; the unavailable owned key "+
				"`stale` is missing", keys)
		}
	})
}

// REQ-52: "An escaped plan MUST be a success carrying escaped=true and
// escape_class."
// REQ-116: "A denied gate is a refusal, not a plan and not an escape class."
// HAPPY PATH — the escape arm is a SUCCESS, not a laundered refusal.
func TestReq52_AnEscapedPlanIsASuccessCarryingEscapedTrueAndItsClass(t *testing.T) {
	model := writeFlowModel(t, flowEscapeModel)
	art := seedArtifact(t, model, "status=draft")

	stdout := requireSuccess(t, "flow", "resolve", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--outcome", "bail", "--as=json")

	data := flowData(t, stdout)
	if data["escaped"] != true {
		t.Errorf("`escaped` = %#v; an escaped plan carries escaped=true",
			data["escaped"])
	}
	// ASSUMPTION A-3: the class reported is the one that was RESCUED — the
	// request's own refusal kind, here `no_match`.
	class, ok := data["escape_class"].(string)
	if !ok || class == "" {
		t.Fatalf("`escape_class` = %#v; an escaped plan carries the class "+
			"the plan came from", data["escape_class"])
	}
	if class != "no_match" {
		t.Errorf("`escape_class` = %q; the escape row models `no_match`, "+
			"and the reported class is the one that was rescued — the "+
			"request's own refusal kind", class)
	}
	if data["rule"] != "bail-escape" {
		t.Errorf("`rule` = %#v; want the escape row `bail-escape`",
			data["rule"])
	}
}

// REQ-52 / REQ-116 (negative arm): an ORDINARY plan is not marked escaped,
// so `escaped` genuinely discriminates.
// BOUNDARY
func TestReq52_AnOrdinaryPlanIsNotMarkedEscaped(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")

	data := flowData(t, requireSuccess(t, "flow", "resolve", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--outcome", "hold", "--as=json"))

	if data["escaped"] == true {
		t.Error("an ordinary exact-one match reported escaped=true; the " +
			"flag would then discriminate nothing")
	}
	if class, ok := data["escape_class"].(string); ok && class != "" {
		t.Errorf("an ordinary plan carries escape_class = %q; the class is "+
			"carried only when the plan came from an escape row", class)
	}
}

// REQ-54: "`flow resolve` requires `--outcome <tag>`; an absent or empty
// value is `flow-tag-invalid` at the CLI, and a value outside the model's
// alphabet is the kernel's `flow-unmodeled-outcome`."
// REQ-125 / `0005:S3`: "`flow resolve` with an unmodeled outcome, an empty
// `--outcome`, a zero-match row, a multi-match row, and an escape row for
// `no_match`" — "each refusal maps to its one-to-one code with non-zero exit
// under both modes".
// INPUT EDGE — the two arms carry DIFFERENT codes, which is the point.
func TestReq54_AbsentOutcomeIsCLISideAndOutOfAlphabetIsKernelSide(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	t.Run("absent", func(t *testing.T) {
		requireRefusal(t, "flow-tag-invalid", 2,
			"flow", "resolve", "--model", model, "--artifact", bind, "--as=json")
	})
	t.Run("empty", func(t *testing.T) {
		requireRefusal(t, "flow-tag-invalid", 2,
			"flow", "resolve", "--model", model, "--artifact", bind,
			"--outcome", "", "--as=json")
	})
	t.Run("outside-alphabet", func(t *testing.T) {
		requireRefusal(t, "flow-unmodeled-outcome", 2,
			"flow", "resolve", "--model", model, "--artifact", bind,
			"--outcome", "not-an-outcome", "--as=json")
	})
}

// REQ-56: "`flow resolve` succeeds only when the kernel reports exactly one
// matching row (or exactly one escape row for an escapable class) and every
// gate on that row allows. Zero matches, multiple matches, unmodeled
// outcomes, unavailable owned state, unevaluable guards, and gate
// deny/indeterminate become typed refusals rather than tie-breaks."
// REQ-119: "Kernel codes mirror the kernel kinds one-to-one; gates run after
// selection and deny is its own refusal".
// another, and none collapses into `command-error`.
// ADVERSARIAL — the closure property: no refusal family collapses into
func TestReq56And119_EveryRefusalFamilyKeepsItsOwnIdentity(t *testing.T) {
	seen := map[string]string{}
	for _, tc := range flowRefusalInvocations(t) {
		t.Run(tc.name, func(t *testing.T) {
			ce := requireRefusal(t, tc.code, tc.exit,
				append(tc.args, "--as=json")...)
			if ce.Code == "command-error" {
				t.Fatal("the refusal collapsed into the generic " +
					"`command-error`; every family carries its own code")
			}
			if prev, held := seen[tc.name]; held && prev != ce.Code {
				t.Errorf("code drifted between runs: %q then %q", prev, ce.Code)
			}
			seen[tc.name] = ce.Code
		})
	}
}

// REQ-50: "It MUST run the selected row's gates only after exact-one
// selection and before emitting the plan; all gates on the row MUST run,
// deny MUST override allow and indeterminate, and every gate result MUST be
// reported."
// REQ-76: `resolve` data minimum includes "`gates[]` (`id`, `result`)".
// REQ-124 / `0005:S2`: "the reader runs, `owned` is assembled from it, the
// plan carries `rule`, `next`, `writes`, `clear[]`, and `gates[]` with
// `allow`; no write accessor runs".
// HAPPY PATH — the allow arm reports the gate on the plan.
func TestReq50And76_GatesRunAfterSelectionAndEveryResultIsReported(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	// `advance` selects the gated row `advance-draft`; the owned keys it
	// requires are established so selection reaches the gate.
	art := seedArtifact(t, model, "status=draft", "stale=x")

	stdout := requireSuccess(t, "flow", "resolve", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--outcome", "advance", "--as=json")

	data := flowData(t, stdout)
	gates, ok := objectsAt(data, "gates")
	if !ok {
		t.Fatalf("the plan carries no `gates[]`; the data minimum fixes "+
			"`gates[]` (`id`, `result`). keys = %v", keysOf(data))
	}
	if len(gates) == 0 {
		t.Fatal("`gates[]` is empty though the selected row carries the " +
			"gate `approval`; every gate on the row MUST run and every " +
			"result MUST be reported")
	}
	var sawApproval bool
	for _, g := range gates {
		if g["id"] == "approval" {
			sawApproval = true
			if g["result"] == nil || g["result"] == "" {
				t.Error("gate `approval` carries no `result`")
			}
		}
	}
	if !sawApproval {
		t.Errorf("`gates[]` = %v omits the selected row's gate `approval`",
			gates)
	}

	// A gate on a row that was NOT selected never ran: `hold-draft` carries
	// no gate, so a gate result for it would mean gates ran pre-selection.
	if len(gates) > 1 {
		t.Errorf("`gates[]` carries %d entries; only the SELECTED row's "+
			"gates run, and `advance-draft` carries exactly one", len(gates))
	}
}

// REQ-50 / REQ-51: "A denied or indeterminate gate MUST surface as
// flow-gate-denied or flow-gate-indeterminate with one findings[] entry per
// gate; it is never a plan and never an escape class."
// REQ-97 / REQ-98: `flow-gate-denied` and `flow-gate-indeterminate` are
// `GroupUserEnv` / 2 with `findings[]` one per gate. `flow-gate-indeterminate`
// fires only when "no gate denied and at least one indeterminate".
// REQ-126 / `0005:S4`: "`flow-gate-denied` then `flow-gate-indeterminate`,
// each with one `findings[]` entry per gate; deny overrides indeterminate
// when both occur; no plan is emitted".
// ADVERSARIAL — deny OVERRIDES indeterminate; neither becomes an escape.
func TestReq51And97And98_GateDenyAndIndeterminateAreTypedRefusalsNeverPlans(t *testing.T) {
	model := writeFlowModel(t, flowGateDenyModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	t.Run("deny", func(t *testing.T) {
		stdout, _, err := runCmd(t, "flow", "resolve", "--model", model,
			"--artifact", bind, "--outcome", "deny-path", "--as=json")
		if err == nil {
			t.Fatalf("a denied gate produced a plan; a deny is a REFUSAL, "+
				"never a plan and never an escape class:\n%s", stdout)
		}
		var ce *clierr.CLIError
		if !asCLIError(err, &ce) {
			t.Fatalf("not a CLIError: %v", err)
		}
		if ce.Code != "flow-gate-denied" {
			t.Fatalf("code = %q; want %q — a deny must NOT be laundered into "+
				"`flow-no-match` or an escape", ce.Code, "flow-gate-denied")
		}
		if clierr.ExitCodeFor(err) != 2 {
			t.Errorf("exit = %d; want 2", clierr.ExitCodeFor(err))
		}
		findings := flowFailureFindings(t, stdout)
		if len(findings) == 0 {
			t.Fatal("`flow-gate-denied` carried no findings; the code table " +
				"fixes `findings[]` ONE PER GATE as its carrier")
		}
	})

	t.Run("indeterminate", func(t *testing.T) {
		stdout, _, err := runCmd(t, "flow", "resolve", "--model", model,
			"--artifact", bind, "--outcome", "indeterminate-path", "--as=json")
		if err == nil {
			t.Fatalf("an indeterminate gate produced a plan:\n%s", stdout)
		}
		var ce *clierr.CLIError
		if !asCLIError(err, &ce) {
			t.Fatalf("not a CLIError: %v", err)
		}
		if ce.Code != "flow-gate-indeterminate" {
			t.Fatalf("code = %q; want %q", ce.Code, "flow-gate-indeterminate")
		}
		// REQ-21: "an indeterminate gate as an exit-2 refusal" — it is NOT
		// an environment failure, so it never exits 3.
		if got := clierr.ExitCodeFor(err); got != 2 {
			t.Errorf("exit = %d; want 2 — an INDETERMINATE gate is a "+
				"decision the model could not make, not an environment that "+
				"could not be consulted", got)
		}
		if len(flowFailureFindings(t, stdout)) == 0 {
			t.Error("`flow-gate-indeterminate` carried no findings")
		}
	})

	t.Run("deny-overrides-indeterminate", func(t *testing.T) {
		// A row carrying BOTH a denying and an indeterminate gate refuses
		// `flow-gate-denied`: deny overrides allow AND indeterminate.
		requireRefusal(t, "flow-gate-denied", 2,
			"flow", "resolve", "--model", model, "--artifact", bind,
			"--outcome", "both-path", "--as=json")
	})
}

// REQ-47: "no gate result is ever dropped silently, and a gate that could
// not be consulted is exit 3 rather than a deny." — binding in BOTH verbs.
// REQ-55: "A gate's timeout or execution failure is an accessor refusal
// (exit 3), not a gate result."
// ADVERSARIAL — the sharpest silent-failure guard in the contract.
func TestReq47And55_AGateThatCouldNotBeConsultedIsExit3AndNeverADeny(t *testing.T) {
	model := writeFlowModel(t, flowGateFailModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	t.Run("resolve", func(t *testing.T) {
		_, _, err := runCmd(t, "flow", "resolve", "--model", model,
			"--artifact", bind, "--outcome", "advance", "--as=json")
		if err == nil {
			t.Fatal("a gate that could not execute produced a plan")
		}
		var ce *clierr.CLIError
		if !asCLIError(err, &ce) {
			t.Fatalf("not a CLIError: %v", err)
		}
		if ce.Code == "flow-gate-denied" {
			t.Fatal("a gate that could not be CONSULTED was reported as a " +
				"DENY; that is the silent-failure mode the contract names " +
				"— it must be exit 3 instead")
		}
		if !slices.Contains(
			[]string{"flow-accessor-failed", "flow-accessor-timeout"},
			ce.Code) {
			t.Fatalf("code = %q; a gate timeout or execution failure is an "+
				"ACCESSOR refusal: `flow-accessor-timeout` or "+
				"`flow-accessor-failed`", ce.Code)
		}
		if got := clierr.ExitCodeFor(err); got != 3 {
			t.Errorf("exit = %d; want 3 — the environment could not be "+
				"consulted and the same request may be re-run unchanged", got)
		}
	})

	t.Run("next-evaluate-gates", func(t *testing.T) {
		// REQ-46: "gate times out / fails to execute" → "exit 3 … exit 3,
		// same codes" — the rule binds under `next` too.
		_, _, err := runCmd(t, "flow", "next", "--model", model,
			"--artifact", bind, "--evaluate-gates", "--as=json")
		if err == nil {
			t.Fatal("`next --evaluate-gates` succeeded though a gate could " +
				"not be consulted; that is exit 3 in BOTH verbs")
		}
		if got := clierr.ExitCodeFor(err); got != 3 {
			t.Errorf("exit = %d; want 3 — the same codes bind under `next`",
				got)
		}
	})
}

// REQ-76: `resolve` data minimum — "`model`, `revision`, `observed`,
// `owned`, `readers[]`, `outcome`, `rule` …, `gates[]` …, `next` …,
// `writes` …, `clear[]` …, `escaped` … and `escape_class`".
// HAPPY PATH
func TestReq76_ResolvePayloadCarriesItsFullDataMinimum(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft", "stale=x")

	data := flowData(t, requireSuccess(t, "flow", "resolve", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--outcome", "advance", "--as=json"))

	for _, key := range []string{
		"model", "revision", "observed", "owned", "readers", "outcome",
		"rule", "gates", "next", "writes", "clear", "escaped",
	} {
		if _, ok := data[key]; !ok {
			t.Errorf("`resolve` payload carries no %q; the data minimum "+
				"fixes it. keys = %v", key, keysOf(data))
		}
	}

	if data["outcome"] != "advance" {
		t.Errorf("`outcome` = %#v; want %q", data["outcome"], "advance")
	}
	if data["rule"] != "advance-draft" {
		t.Errorf("`rule` = %#v; want the matched rule identity %q",
			data["rule"], "advance-draft")
	}
	// The selected row's clear list rides `clear[]`, and the `<clear>`
	// sentinel never appears as a write VALUE.
	clears, ok := stringsAt(data, "clear")
	if !ok || !slices.Contains(clears, "stale") {
		t.Errorf("`clear[]` = %#v; the selected row clears `stale`",
			data["clear"])
	}
	if writes, ok := data["writes"].(map[string]any); ok {
		for k, v := range writes {
			if s, ok := v.(string); ok && s == "<clear>" {
				t.Errorf("`writes[%s]` carries the `<clear>` SENTINEL as a "+
					"value; planned clears ride `clear[]`", k)
			}
		}
	}
}

// REQ-57: "`flow resolve` performs reads, so it is no longer effect-free".
// REQ-30: "Owned state is never caller-supplied: it is assembled from
// declared read accessors over `--artifact role=path` bindings".
// HAPPY PATH — `owned` came from the reader, and the reader is named.
func TestReq30And57_ResolveAssemblesOwnedStateFromTheDeclaredReaders(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")

	data := flowData(t, requireSuccess(t, "flow", "resolve", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--outcome", "hold", "--as=json"))

	owned, ok := data["owned"].(map[string]any)
	if !ok {
		t.Fatalf("`owned` is not an object: %#v", data["owned"])
	}
	if owned["status"] != "draft" {
		t.Errorf("owned[status] = %#v; the reader established %q",
			owned["status"], "draft")
	}
	// The reader that produced it is named among the invoked identities.
	if ids, ok := stringsAt(data, "readers"); ok {
		if !slices.Contains(ids, "state") {
			t.Errorf("`readers[]` = %v omits `state`, which produced the "+
				"owned snapshot", ids)
		}
		return
	}
	objs, _ := objectsAt(data, "readers")
	var sawState bool
	for _, r := range objs {
		if r["id"] == "state" {
			sawState = true
		}
	}
	if !sawState {
		t.Error("`readers[]` omits `state`, which produced the owned snapshot")
	}
}

// REQ-106: "The main silent-failure risk is treating a failed accessor, a
// denied gate, or an ambiguous row as a successful transition. The recovery
// path is refusal-first: the command exits non-zero, emits the structured
// error envelope, and leaves skill judgment or model repair to the caller."
// a non-zero exit carrying a structured envelope.
// ADVERSARIAL — the three named silent-failure shapes, each asserted to be
func TestReq106_NoneOfTheThreeSilentFailureShapesBecomesASuccess(t *testing.T) {
	ambiguous := writeFlowModel(t, flowAmbiguousModel)
	denyModel := writeFlowModel(t, flowGateDenyModel)
	failModel := writeFlowModel(t, flowGateFailModel)

	ambigArt := seedArtifact(t, ambiguous, "status=draft")
	denyArt := seedArtifact(t, denyModel, "status=draft")
	failArt := seedArtifact(t, failModel, "status=draft")

	cases := []struct {
		name string
		args []string
	}{
		{"ambiguous-row", []string{"flow", "resolve", "--model", ambiguous,
			"--artifact", artifactBinding(flowStateRole, ambigArt),
			"--outcome", "advance"}},
		{"denied-gate", []string{"flow", "resolve", "--model", denyModel,
			"--artifact", artifactBinding(flowStateRole, denyArt),
			"--outcome", "deny-path"}},
		{"failed-accessor", []string{"flow", "resolve", "--model", failModel,
			"--artifact", artifactBinding(flowStateRole, failArt),
			"--outcome", "advance"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, _, err := runCmd(t, append(tc.args, "--as=json")...)
			if err == nil {
				t.Fatalf("%s was treated as a SUCCESSFUL transition; the "+
					"recovery path is refusal-first\nstdout:\n%s",
					tc.name, stdout)
			}
			if clierr.ExitCodeFor(err) == 0 {
				t.Errorf("%s exited 0", tc.name)
			}
			var ce *clierr.CLIError
			if !asCLIError(err, &ce) {
				t.Fatalf("%s did not emit a structured error envelope: %v",
					tc.name, err)
			}
			if !strings.HasPrefix(ce.Code, "flow-") {
				t.Errorf("%s carries code %q; every refusal in this contract "+
					"is a `flow-*` code", tc.name, ce.Code)
			}
		})
	}
}

// REQ-113: "It does not own … skill execution".
// REQ-53: "flow resolve MUST NOT … initiate skill work".
// data the caller acts on, never work the CLI performs.
// BOUNDARY — the kernel is called once and the verb terminates; the plan is
func TestReq53And113_ResolveEmitsAPlanAndInitiatesNoWorkItself(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	before := requireSuccess(t, "flow", "read-state", "--model", model,
		"--artifact", bind, "--artifact", artifactBinding(flowOrphanRole, art),
		"--as=json")

	// `resolve` on the row that PLANS a write to `status`…
	data := flowData(t, requireSuccess(t, "flow", "resolve", "--model", model,
		"--artifact", bind, "--outcome", "hold", "--as=json"))
	if _, ok := data["writes"]; !ok {
		t.Fatal("the plan carries no `writes`")
	}

	// …changes nothing: the plan is DESCRIBED, never applied. Applying it
	// is `set-state`'s job, and nothing links the two calls (REQ-68).
	after := requireSuccess(t, "flow", "read-state", "--model", model,
		"--artifact", bind, "--artifact", artifactBinding(flowOrphanRole, art),
		"--as=json")
	if before != after {
		t.Errorf("`resolve` applied its own plan; it emits a plan and "+
			"initiates no work\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

// REQ-96: `flow-guard-unevaluable` — kernel `guard_unevaluable` —
// `GroupUserEnv` / 2 — "`findings[]` (rows; atoms when the kernel names
// them)".
// DOMAIN EDGE — the atom fields ride FLAT on the finding (REQ-15).
func TestReq96_GuardUnevaluableCarriesRowsAndFlatAtomFields(t *testing.T) {
	model := writeFlowModel(t, flowGuardUnevaluableModel)
	art := seedArtifact(t, model, "status=draft")

	stdout, _, err := runCmd(t, "flow", "resolve", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--outcome", "advance", "--as=json")
	if err == nil {
		t.Fatalf("an unevaluable guard produced a plan:\n%s", stdout)
	}
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("not a CLIError: %v", err)
	}
	if ce.Code != "flow-guard-unevaluable" {
		t.Fatalf("code = %q; want %q", ce.Code, "flow-guard-unevaluable")
	}
	if got := clierr.ExitCodeFor(err); got != 2 {
		t.Errorf("exit = %d; want 2", got)
	}

	findings := flowFailureFindings(t, stdout)
	if len(findings) == 0 {
		t.Fatal("`flow-guard-unevaluable` carried no findings; the code " +
			"table fixes `findings[]` (rows; atoms) as its carrier")
	}
	var sawRow, sawAtom bool
	for _, f := range findings {
		if f.Rule != "" {
			sawRow = true
		}
		// REQ-15: "no producer nests its own fields in a sub-object" — the
		// atom's key/operator/literal/block sit FLAT on the record.
		if f.Key != "" || f.Operator != "" || f.Block != "" {
			sawAtom = true
		}
	}
	if !sawRow {
		t.Error("no finding names a row; the kernel names every undecidable " +
			"row")
	}
	if !sawAtom {
		t.Error("no finding carries the atom's flat `key` / `operator` / " +
			"`block` fields, though the kernel names the unevaluable atoms")
	}
}

// REQ-49 closure: "Kernel refusals MUST map one-to-one onto" the five
// codes — so the CLI raises FIVE DISTINCT kernel-mapped codes over the
// five kinds, and no two kinds collapse onto one code.
// BOUNDARY
func TestReq49_TheFiveKernelKindsRaiseFiveDistinctCodes(t *testing.T) {
	mvv := writeFlowModel(t, flowMVVModel)
	ambiguous := writeFlowModel(t, flowAmbiguousModel)
	guardUnev := writeFlowModel(t, flowGuardUnevaluableModel)
	mvvArt := seedArtifact(t, mvv, "status=draft")
	ambigArt := seedArtifact(t, ambiguous, "status=draft")
	guardArt := seedArtifact(t, guardUnev, "status=draft")

	invocations := map[resolve.RefusalKind][]string{
		resolve.KindUnmodeledOutcome: {"flow", "resolve", "--model", mvv,
			"--artifact", artifactBinding(flowStateRole, mvvArt),
			"--outcome", "not-an-outcome"},
		resolve.KindNoMatch: {"flow", "resolve", "--model", mvv,
			"--artifact", artifactBinding(flowStateRole, mvvArt),
			"--outcome", "bail"},
		resolve.KindAmbiguousMatch: {"flow", "resolve", "--model", ambiguous,
			"--artifact", artifactBinding(flowStateRole, ambigArt),
			"--outcome", "advance"},
		resolve.KindOwnedStateUnavailable: {"flow", "resolve", "--model", mvv,
			"--artifact", artifactBinding(flowStateRole, mvvArt),
			"--outcome", "advance"},
		resolve.KindGuardUnevaluable: {"flow", "resolve", "--model", guardUnev,
			"--artifact", artifactBinding(flowStateRole, guardArt),
			"--outcome", "advance"},
	}

	seen := map[string]resolve.RefusalKind{}
	for _, kind := range resolve.RefusalKinds() {
		args, ok := invocations[kind]
		if !ok {
			t.Errorf("the kernel kind %q has no invocation in this suite; "+
				"the mapping is one-to-one over the CLOSED five-kind set, "+
				"so a new kind needs a code and a test", kind)
			continue
		}
		_, _, err := runCmd(t, append(slices.Clone(args), "--as=json")...)
		if err == nil {
			t.Errorf("%s: the invocation succeeded; want a kernel refusal",
				kind)
			continue
		}
		var ce *clierr.CLIError
		if !asCLIError(err, &ce) {
			t.Errorf("%s: not a CLIError: %v", kind, err)
			continue
		}
		want := "flow-" + strings.ReplaceAll(string(kind), "_", "-")
		if ce.Code != want {
			t.Errorf("kernel kind %q raised %q; the codes mirror the kinds "+
				"ONE-TO-ONE as `flow-<kind>` with underscores hyphenated, "+
				"so it must be %q", kind, ce.Code, want)
			continue
		}
		if prev, collided := seen[ce.Code]; collided {
			t.Errorf("kinds %q and %q both raise %q; the mapping is "+
				"one-to-one and no two kinds may collapse onto one code",
				prev, kind, ce.Code)
		}
		seen[ce.Code] = kind
	}

	if len(seen) != 5 {
		t.Errorf("the CLI raised %d distinct kernel-mapped codes; the "+
			"one-to-one mapping over the closed five-kind set yields five",
			len(seen))
	}
}
