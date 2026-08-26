package cli

// RDR 0005 — the exit-code contract (REQ-19..REQ-21, REQ-99..REQ-105), the
// internal-group codes (REQ-102), the ownership fences (REQ-113..REQ-117),
// and the documentation obligation (REQ-131).
//
// The load-bearing property here is that exit 3 is EARNED, not defaulted:
// "the environment could not be consulted and the same request may be
// re-run unchanged". A refusal that exits 3 must be genuinely retriable,
// and a refusal about the REQUEST must never exit 3 — otherwise a caller's
// retry loop spins on an input error forever.

import (
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/accessor"
	"github.com/newcoinc/intrastate/internal/cli/clierr"
)

// REQ-19: "Exit 3 MUST mean the environment could not be consulted and the
// same request may be re-run unchanged; every other failure MUST exit 2. No
// new exit group."
// REQ-20: "the exit-3 population is exactly those five classes" — accessor
// timeout, execution failure, incomplete read, read-back incomplete, and
// post-mutation timeout.
// BOUNDARY — the exit-3 code set is CLOSED at exactly five spellings.
func TestReq19And20_TheExit3PopulationIsExactlyTheFiveNamedClasses(t *testing.T) {
	// The five codes the code table places at `GroupEnvUnavailable` / 3.
	exit3 := []string{
		"flow-accessor-timeout",
		"flow-accessor-failed",
		"flow-read-incomplete",
		"flow-write-readback-incomplete",
		"flow-write-readback-timeout",
	}

	// Every OTHER refusal this contract raises exits 2. The refusal table
	// drives the assertion, so a code that silently migrates to exit 3
	// fails here.
	for _, tc := range flowRefusalInvocations(t) {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := runCmd(t, append(tc.args, "--as=json")...)
			if err == nil {
				t.Fatalf("%s succeeded; want refusal %s", tc.name, tc.code)
			}
			var ce *clierr.CLIError
			if !asCLIError(err, &ce) {
				t.Fatalf("not a CLIError: %v", err)
			}
			got := clierr.ExitCodeFor(err)

			var wantEnvUnavailable bool
			for _, code := range exit3 {
				if ce.Code == code {
					wantEnvUnavailable = true
				}
			}
			switch {
			case wantEnvUnavailable && got != 3:
				t.Errorf("%s exits %d; it names an environment that could "+
					"not be consulted and MUST exit 3", ce.Code, got)
			case !wantEnvUnavailable && got != 2:
				t.Errorf("%s exits %d; every failure that is not one of the "+
					"five environment classes MUST exit 2 — a request error "+
					"that exits 3 tells a caller to retry an input it must "+
					"instead fix", ce.Code, got)
			}
		})
	}
}

// REQ-19: "the same request may be re-run unchanged" — the retriability
// promise exit 3 makes.
// REQ-110: refusal identity is deterministic.
// ADVERSARIAL — an exit-3 refusal re-run unchanged must reach the SAME
// disposition, because a caller's remedy is to repair the environment and
// re-issue the identical request.
func TestReq19_AnExit3RefusalIsReRunnableUnchangedToTheSameDisposition(t *testing.T) {
	model := writeFlowModel(t, flowGateFailModel)
	art := seedArtifact(t, model, "status=draft")

	args := []string{"flow", "resolve", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--outcome", "advance", "--as=json"}

	first, _, ferr := runCmd(t, args...)
	if clierr.ExitCodeFor(ferr) != 3 {
		t.Fatalf("the fixture did not produce an exit-3 refusal: exit=%d "+
			"err=%v", clierr.ExitCodeFor(ferr), ferr)
	}
	for i := range 3 {
		again, _, aerr := runCmd(t, args...)
		if clierr.ExitCodeFor(aerr) != 3 {
			t.Errorf("re-run %d exits %d; the SAME request re-run unchanged "+
				"reaches the same disposition", i+1, clierr.ExitCodeFor(aerr))
		}
		if again != first {
			t.Errorf("re-run %d differs:\nfirst:\n%s\nagain:\n%s",
				i+1, first, again)
		}
	}
}

// REQ-21: "an unavailable accessor surfaces as an exit-3 refusal and an
// indeterminate gate as an exit-2 refusal."
// DOMAIN EDGE — the pair is the point: both involve a gate, and they exit
// DIFFERENTLY, because one is an environment fault and the other a modeled
// verdict.
func TestReq21_UnavailableAccessorExits3ButIndeterminateGateExits2(t *testing.T) {
	failModel := writeFlowModel(t, flowGateFailModel)
	denyModel := writeFlowModel(t, flowGateDenyModel)
	failArt := seedArtifact(t, failModel, "status=draft")
	denyArt := seedArtifact(t, denyModel, "status=draft")

	_, _, unavailable := runCmd(t, "flow", "resolve", "--model", failModel,
		"--artifact", artifactBinding(flowStateRole, failArt),
		"--outcome", "advance", "--as=json")
	if got := clierr.ExitCodeFor(unavailable); got != 3 {
		t.Errorf("an unavailable accessor exits %d; want 3", got)
	}

	_, _, indeterminate := runCmd(t, "flow", "resolve", "--model", denyModel,
		"--artifact", artifactBinding(flowStateRole, denyArt),
		"--outcome", "indeterminate-path", "--as=json")
	if got := clierr.ExitCodeFor(indeterminate); got != 2 {
		t.Errorf("an indeterminate gate exits %d; want 2 — the model could "+
			"not decide, which is not an environment that could not be "+
			"consulted, and retrying it unchanged cannot help", got)
	}
}

// REQ-99: `flow-accessor-timeout` — "accessor timed out" —
// `GroupEnvUnavailable` / 3 — "`param` = accessor id".
// REQ-100: `flow-accessor-failed` — "accessor execution failed" — same
// group — "`param` = accessor id".
// REQ-101: `flow-read-incomplete` — "read returned an incomplete key set" —
// same group — "`param` = accessor id".
// DOMAIN EDGE — the carrier is the ACCESSOR ID, which is what a caller
// needs to repair the environment.
func TestReq99And100And101_AccessorRefusalsCarryTheAccessorIdAsParam(t *testing.T) {
	model := writeFlowModel(t, flowGateFailModel)
	art := seedArtifact(t, model, "status=draft")

	_, _, err := runCmd(t, "flow", "resolve", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--outcome", "advance", "--as=json")
	if err == nil {
		t.Fatal("the unconsultable gate produced a plan")
	}
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("not a CLIError: %v", err)
	}
	if clierr.ExitCodeFor(err) != 3 {
		t.Fatalf("exit = %d; want 3", clierr.ExitCodeFor(err))
	}
	if ce.Param != "unreachable" {
		t.Errorf("param = %q; the code table fixes `param` = ACCESSOR ID, "+
			"and the unconsultable accessor here is %q — without it the "+
			"caller cannot tell which part of the environment to repair",
			ce.Param, "unreachable")
	}
}

// REQ-102: `flow-accessor-unknown`, `flow-accessor-capability-mismatch`,
// `flow-write-non-owned` — "unknown accessor, capability mismatch, write to
// a non-owned tag (unreachable past load and CLI checks)" — `GroupInternal`
// / 2 — no carrier.
// A-10: `flow-write-non-owned` is "a defensive arm, not a reachable user
// path; REQ-61's `flow-write-unbound` is what a caller actually hits."
// ADVERSARIAL — the observable obligation is that these three are NOT what
// a caller reaches: a user-path mistake must surface as its user-facing
// code, never as an internal one.
func TestReq102_TheThreeInternalCodesAreUnreachableFromTheUserPath(t *testing.T) {
	internal := []string{
		"flow-accessor-unknown",
		"flow-accessor-capability-mismatch",
		"flow-write-non-owned",
	}

	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	// Every user-path mistake the CLI surface admits.
	attempts := [][]string{
		{"flow", "set-state", "--model", model, "--artifact", bind,
			"--write", "note=x"},
		{"flow", "set-state", "--model", model, "--artifact", bind,
			"--write", "no-such-tag=x"},
		{"flow", "set-state", "--model", model, "--artifact", bind,
			"--clear", "note"},
		{"flow", "set-state", "--model", model, "--artifact", bind,
			"--write", "profile=mid"},
		{"flow", "next", "--model", model, "--artifact", "orphan=/nonexistent"},
	}

	for _, args := range attempts {
		_, _, err := runCmd(t, append(args, "--as=json")...)
		if err == nil {
			continue
		}
		var ce *clierr.CLIError
		if !asCLIError(err, &ce) {
			continue
		}
		for _, code := range internal {
			if ce.Code == code {
				t.Errorf("the user-path invocation %v reached the INTERNAL "+
					"code %q; those three are unreachable past load and the "+
					"CLI checks — a caller mistake surfaces as its own "+
					"user-facing code", args, code)
			}
		}
	}

	// And when raised, they carry GroupInternal / exit 2 — never exit 3,
	// which would tell a caller to retry an internal invariant violation.
	for _, code := range internal {
		ce := &clierr.CLIError{Code: code, Group: clierr.GroupInternal}
		if got := clierr.ExitCodeFor(ce); got != 2 {
			t.Errorf("%s at GroupInternal exits %d; want 2", code, got)
		}
	}
}

// REQ-20 (the mapping leg) / A-9: the accessor layer's refusal classes are
// what the exit-3 codes are derived FROM, and the class set is closed at
// eight.
// BOUNDARY — a new accessor class appearing without a CLI mapping would be
// an unmapped refusal, which the never-silent contract forbids.
func TestReq20_EveryAccessorRefusalClassHasACLIMapping(t *testing.T) {
	// The mapping A-9 records, class -> the `flow-*` code(s) it may take.
	// Two classes are verb-phase split: only the invoking phase
	// distinguishes a read from a post-mutation read-back.
	mapping := map[accessor.RefusalClass][]string{
		accessor.ClassTimeout: {
			"flow-accessor-timeout", "flow-write-readback-timeout",
		},
		accessor.ClassExecutionFailure:   {"flow-accessor-failed"},
		accessor.ClassIncompleteRead:     {"flow-read-incomplete"},
		accessor.ClassReadBackIncomplete: {"flow-write-readback-incomplete"},
		accessor.ClassReadBackMismatch:   {"flow-write-readback-mismatch"},
		accessor.ClassGateIndeterminate:  {"flow-gate-indeterminate"},
		accessor.ClassUnknownAccessor:    {"flow-accessor-unknown"},
		accessor.ClassCapabilityMismatch: {"flow-accessor-capability-mismatch"},
	}

	// The clause is about THIS contract's mapping, so the surface that
	// carries it must exist before the check discriminates anything.
	if names := flowGroupNames(t); len(names) != len(flowVerbs) {
		t.Fatalf("the `flow` group registers %v; the accessor-class mapping "+
			"is this RDR's, so its verbs must exist before the mapping is "+
			"assertable", names)
	}

	classes := accessor.RefusalClasses()
	if len(classes) != len(mapping) {
		t.Errorf("the accessor layer exposes %d refusal classes but this "+
			"contract maps %d; every class MUST reach a `flow-*` code, or a "+
			"refusal escapes the never-silent envelope",
			len(classes), len(mapping))
	}
	for _, class := range classes {
		codes, mapped := mapping[class]
		if !mapped {
			t.Errorf("accessor refusal class %q has no `flow-*` mapping",
				class)
			continue
		}
		for _, code := range codes {
			if !strings.HasPrefix(code, "flow-") {
				t.Errorf("class %q maps to %q, which is not a `flow-*` code",
					class, code)
			}
		}
	}
}

// REQ-117: "the read-back is the commit-time check, and the window between
// `resolve` and `set-state` is a stated non-guarantee of the two-verb
// design (§D9)."
// REQ-68: "nothing links a `set-state` request to a prior `resolve`."
// DOMAIN EDGE — state that CHANGED between the two calls does not
// invalidate the `set-state`: the read-back, not the plan's freshness, is
// what gates success. This is a stated non-guarantee, so the test pins the
// behaviour rather than treating it as a defect.
func TestReq117_TheResolveToSetStateWindowIsAStatedNonGuarantee(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	// A plan is minted against `status = draft` …
	plan := flowData(t, requireSuccess(t, "flow", "resolve", "--model", model,
		"--artifact", bind, "--outcome", "hold", "--as=json"))
	if plan["rule"] != "hold-draft" {
		t.Fatalf("`rule` = %#v; want %q", plan["rule"], "hold-draft")
	}

	// … the world moves underneath it …
	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", bind, "--write", "status=final", "--as=json")

	// … and the transcribed plan is still accepted, because nothing links
	// the two calls and the read-back is the ONLY commit-time check.
	data := flowData(t, requireSuccess(t, "flow", "set-state",
		"--model", model, "--artifact", bind,
		"--write", "status=draft", "--as=json"))

	if owned, ok := data["owned"].(map[string]any); ok {
		if owned["status"] != "draft" {
			t.Errorf("owned[status] = %#v; the read-back verifies what the "+
				"REQUEST planned, and the request planned %q",
				owned["status"], "draft")
		}
	}
}

// REQ-131: "Land the full code table, update `docs/cli-output-contract.md`
// for the `findings` field and the exit-3 rule, and document illustrative
// invocations."
// BOUNDARY — the documentation obligation is checked against the shipped
// doc, which the RDR names as the authoritative output contract.
func TestReq131_TheOutputContractDocumentsFindingsAndTheExit3Rule(t *testing.T) {
	doc := readRepoFile(t, "docs/cli-output-contract.md")

	if !strings.Contains(doc, "findings") {
		t.Error("docs/cli-output-contract.md does not mention `findings`; " +
			"REQ-131 requires the doc be updated for the field this RDR adds")
	}

	// The exit-3 rule: the doc must say what exit 3 MEANS, not merely that
	// the number exists — a caller branches on the promise, not the digit.
	lower := strings.ToLower(doc)
	if !strings.Contains(lower, "exit 3") && !strings.Contains(doc, "`3`") &&
		!strings.Contains(lower, "exit code 3") {
		t.Error("docs/cli-output-contract.md does not document the exit-3 " +
			"rule")
	}

	// The `flow` group's verbs are documented as illustrative invocations.
	var missing []string
	for _, verb := range flowVerbs {
		if !strings.Contains(doc, "flow "+verb) {
			missing = append(missing, verb)
		}
	}
	if len(missing) > 0 {
		t.Errorf("docs/cli-output-contract.md documents no illustrative "+
			"invocation for %v; REQ-131 requires them", missing)
	}
}
