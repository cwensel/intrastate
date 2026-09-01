package cli

// kata c3xz — `flow set-state --plan <file|->`: apply a `flow resolve`
// envelope without transcribing it.
//
// RDR 0005 named the flag, its `<file|->` argument, and the property that
// motivates it — "copy-through from plan to request is byte-identical" — and
// deferred only the deciding (§506-509). This is that seed discharged as a
// flag on an EXISTING verb: `0005:C1`'s four-verb set is untouched, and so is
// the invariant that "nothing links a `set-state` request to a prior
// `resolve`". A carried plan is INPUT, re-derived against the model exactly
// as a `--write` is; it is not a linkage, and no token or session crosses.
//
// The oracles below are written so that each one fails for a DIFFERENT
// implementation defect:
//
//   - the round-trip pair fails if plan-derived keys skip any part of the
//     `--write` path, and the flag-driven COMPARISON is what makes it a
//     claim about equivalence rather than merely about success;
//   - the envelope-shape arms fail if the `"type"` discriminator is guessed
//     rather than probed, and the refusal arm fails if a `{"type":"failed"}`
//     document is treated as a plan carrying no writes;
//   - the merge arms fail if precedence is implemented as anything other
//     than seeding the EXISTING duplicate check;
//   - the canonicalisation arm fails if a plan value is trusted verbatim
//     instead of being re-canonicalised (REQ-71).

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
)

// runCmdStdin is `runCmd` with a stdin body, which `--plan -` needs and the
// shipped harness has no reason to carry. It drives the SAME
// `ExecuteAndEmit` production path — only the input stream is added, so a
// piped plan is exercised through the real flag/stream wiring rather than by
// calling the decoder directly.
func runCmdStdin(t *testing.T, stdin string, args ...string) (stdout, stderr string, err error) {
	t.Helper()

	cmd := NewRootCmd()
	var out, errBuf bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errBuf)
	cmd.SetIn(strings.NewReader(stdin))
	err = ExecuteAndEmit(cmd, args)
	return out.String(), errBuf.String(), err
}

// resolvePlanC3xz drives `flow resolve --as json` and returns the raw
// envelope line — the exact bytes a caller pipes.
func resolvePlanC3xz(t *testing.T, model, bind, outcome string) string {
	t.Helper()

	return strings.TrimSpace(requireSuccess(t, "flow", "resolve",
		"--model", model, "--artifact", bind,
		"--outcome", outcome, "--as=json"))
}

// planDataC3xz extracts the `data` object out of an envelope, as a caller
// running `jq .data` would.
func planDataC3xz(t *testing.T, envelope string) string {
	t.Helper()

	var env struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal([]byte(envelope), &env); err != nil {
		t.Fatalf("the resolve envelope is not one JSON object: %v\n%s",
			err, envelope)
	}
	return string(env.Data)
}

// writePlanFileC3xz writes a plan document to a temp file for `--plan <file>`.
func writePlanFileC3xz(t *testing.T, body string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "plan.json")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("writing the plan fixture failed: %v", err)
	}
	return path
}

// applyResultC3xz is the part of a `set-state` payload two invocations must
// agree on for the plan path to be equivalent to the flag path.
type applyResultC3xz struct {
	writes map[string]any
	clear  []string
	owned  map[string]any
}

func applyResultOfC3xz(t *testing.T, stdout string) applyResultC3xz {
	t.Helper()

	data := flowData(t, stdout)
	writes, ok := data["writes"].(map[string]any)
	if !ok {
		t.Fatalf("`writes` is not an object: %#v", data["writes"])
	}
	owned, ok := data["owned"].(map[string]any)
	if !ok {
		t.Fatalf("`owned` is not an object: %#v", data["owned"])
	}
	clear, ok := stringsAt(data, "clear")
	if !ok {
		t.Fatalf("`clear` is not an array of strings: %#v", data["clear"])
	}
	return applyResultC3xz{writes: writes, clear: clear, owned: owned}
}

func sameApplyResultC3xz(t *testing.T, got, want applyResultC3xz, why string) {
	t.Helper()

	if g, w := jsonOfC3xz(t, got.writes), jsonOfC3xz(t, want.writes); g != w {
		t.Errorf("writes = %s; the flag-driven equivalent produced %s.\n%s",
			g, w, why)
	}
	if g, w := jsonOfC3xz(t, got.clear), jsonOfC3xz(t, want.clear); g != w {
		t.Errorf("clear = %s; the flag-driven equivalent produced %s.\n%s",
			g, w, why)
	}
	if g, w := jsonOfC3xz(t, got.owned), jsonOfC3xz(t, want.owned); g != w {
		t.Errorf("owned = %s; the flag-driven equivalent produced %s.\n"+
			"`owned` is the READ-BACK-CONFIRMED state, so a difference here "+
			"means the two carriers wrote different things to the artifact.\n"+
			"%s", g, w, why)
	}
}

func jsonOfC3xz(t *testing.T, v any) string {
	t.Helper()

	out, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshalling a payload fragment failed: %v", err)
	}
	return string(out)
}

// --- the round trip ------------------------------------------------------

// The kata's acceptance: `resolve --as json | set-state --plan -` applies a
// plan with ZERO `--write` flags, and lands what the flag-driven request
// lands.
//
// The comparison is against a flag-driven run over a SEPARATE artifact
// seeded identically, so the claim is equivalence and not merely "the plan
// run succeeded". A `--plan` that dropped a key, skipped the clear, or wrote
// a differently-canonicalised value all survive a success-only oracle and
// all fail here.
//
// HAPPY PATH
func TestC3xz_AResolveEnvelopePipedAsAPlanAppliesWithZeroWriteFlags(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)

	// The MVV model's `advance` row writes `status=final` and clears
	// `stale`, so one plan exercises BOTH carriers at once.
	planArt := seedArtifact(t, model, "status=draft", "stale=obsolete")
	planBind := artifactBinding(flowStateRole, planArt)
	envelope := resolvePlanC3xz(t, model, planBind, "advance")

	viaPlan := applyResultOfC3xz(t, requireSuccessStdinC3xz(t, envelope,
		"flow", "set-state", "--model", model,
		"--artifact", planBind, "--plan", "-", "--as=json"))

	// The flag-driven equivalent, transcribed by hand the way a caller does
	// today, over its own identically-seeded artifact.
	flagArt := seedArtifact(t, model, "status=draft", "stale=obsolete")
	viaFlags := applyResultOfC3xz(t, requireSuccess(t, "flow", "set-state",
		"--model", model, "--artifact", artifactBinding(flowStateRole, flagArt),
		"--write", "status=final", "--clear", "stale", "--as=json"))

	sameApplyResultC3xz(t, viaPlan, viaFlags,
		"a carried plan must produce the SAME request as transcribing it "+
			"by hand — that equivalence is the whole point of --plan, and "+
			"it holds because plan-derived keys run through the identical "+
			"writerFor + canonicalValue path")

	// And the transition really happened, observed independently.
	if viaPlan.owned["status"] != "final" {
		t.Errorf("owned[status] = %#v after applying the plan; want %q",
			viaPlan.owned["status"], "final")
	}
	if _, held := viaPlan.owned["stale"]; held {
		t.Error("`stale` is present in `owned` after the plan's clear[] " +
			"was applied; a cleared key reads back ABSENT")
	}
}

// requireSuccessStdinC3xz is `requireSuccess` with a stdin body.
func requireSuccessStdinC3xz(t *testing.T, stdin string, args ...string) string {
	t.Helper()

	stdout, stderr, err := runCmdStdin(t, stdin, args...)
	if err != nil {
		t.Fatalf("invocation failed: %v\nstdout:\n%s\nstderr:\n%s",
			err, stdout, stderr)
	}
	return stdout
}

// requireRefusalStdinC3xz is `requireRefusal` with a stdin body: exact code,
// exact exit.
func requireRefusalStdinC3xz(t *testing.T, stdin, code string, exit int, args ...string) *clierr.CLIError {
	t.Helper()

	stdout, _, err := runCmdStdin(t, stdin, args...)
	if err == nil {
		t.Fatalf("invocation succeeded; want refusal %s\nstdout:\n%s",
			code, stdout)
	}
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("refusal is not a structured CLIError: %v", err)
	}
	if ce.Code != code {
		t.Fatalf("refusal code = %q; want %q", ce.Code, code)
	}
	if got := clierr.ExitCodeFor(err); got != exit {
		t.Errorf("exit = %d; want %d — the group fixes the exit", got, exit)
	}
	return ce
}

// `--plan <file>` is the same carrier over a different source. A file and a
// pipe must not diverge: the document is identical, so the applied result
// must be too.
//
// HAPPY PATH
func TestC3xz_APlanFileAppliesIdenticallyToAPipedPlan(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)

	src := seedArtifact(t, model, "status=draft", "stale=obsolete")
	envelope := resolvePlanC3xz(t, model, artifactBinding(flowStateRole, src),
		"advance")
	planFile := writePlanFileC3xz(t, envelope)

	fileArt := seedArtifact(t, model, "status=draft", "stale=obsolete")
	viaFile := applyResultOfC3xz(t, requireSuccess(t, "flow", "set-state",
		"--model", model, "--artifact", artifactBinding(flowStateRole, fileArt),
		"--plan", planFile, "--as=json"))

	pipeArt := seedArtifact(t, model, "status=draft", "stale=obsolete")
	viaPipe := applyResultOfC3xz(t, requireSuccessStdinC3xz(t, envelope,
		"flow", "set-state", "--model", model,
		"--artifact", artifactBinding(flowStateRole, pipeArt),
		"--plan", "-", "--as=json"))

	sameApplyResultC3xz(t, viaFile, viaPipe,
		"`--plan <file>` and `--plan -` read the SAME document; only the "+
			"source differs, so nothing about the applied request may")
}

// --- the accepted envelope shapes ---------------------------------------

// BOTH document shapes are accepted, discriminated on the presence of a
// top-level `"type"` key: the full `{"type":"ok","data":…}` envelope, and
// the bare `data` object a caller extracted from it.
//
// The two arms are compared to EACH OTHER, so an implementation that
// accepted one and silently applied nothing from the other fails — which a
// per-arm success assertion would not catch, since applying no writes is
// itself a success.
//
// DOMAIN EDGE
func TestC3xz_BothTheFullEnvelopeAndABareDataObjectAreAccepted(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)

	src := seedArtifact(t, model, "status=draft", "stale=obsolete")
	envelope := resolvePlanC3xz(t, model, artifactBinding(flowStateRole, src),
		"advance")
	bare := planDataC3xz(t, envelope)

	// The control: the two documents really are different bytes, so the
	// comparison below is about two SHAPES and not about one document read
	// twice.
	if bare == envelope {
		t.Fatal("the bare `data` object is byte-identical to the envelope; " +
			"the two shapes are not distinguishable and this oracle proves " +
			"nothing")
	}
	if strings.Contains(bare, `"type"`) {
		t.Fatalf("the bare `data` object carries a top-level `type` key, so "+
			"the discriminator is ambiguous:\n%s", bare)
	}

	envArt := seedArtifact(t, model, "status=draft", "stale=obsolete")
	viaEnvelope := applyResultOfC3xz(t, requireSuccessStdinC3xz(t, envelope,
		"flow", "set-state", "--model", model,
		"--artifact", artifactBinding(flowStateRole, envArt),
		"--plan", "-", "--as=json"))

	bareArt := seedArtifact(t, model, "status=draft", "stale=obsolete")
	viaBare := applyResultOfC3xz(t, requireSuccessStdinC3xz(t, bare,
		"flow", "set-state", "--model", model,
		"--artifact", artifactBinding(flowStateRole, bareArt),
		"--plan", "-", "--as=json"))

	sameApplyResultC3xz(t, viaBare, viaEnvelope,
		"a bare `data` object and the envelope that wraps it carry the SAME "+
			"plan; the discriminator selects which bytes to decode and "+
			"changes nothing about what is applied")

	// Not vacuous: the plan really carried something.
	if len(viaEnvelope.writes) == 0 && len(viaEnvelope.clear) == 0 {
		t.Fatal("the applied plan carried neither writes nor clears; the " +
			"comparison above would hold between two no-ops")
	}
}

// A REFUSAL envelope piped in place of a plan REFUSES. This is the one
// failure mode a carried plan makes newly reachable: `resolve` writes its
// refusal to the same stream its success goes to under `--as=json`, so a
// pipe whose left side failed delivers `{"type":"failed",…}` — and a decoder
// that merely found no `writes` in it would report SUCCESS for a transition
// the model declined.
//
// ADVERSARIAL
func TestC3xz_ARefusalEnvelopeOnStdinRefuses(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	// A real refusal from the real verb, not a hand-written one: `hold` on a
	// model whose `status` is already `final` matches no row.
	advanced := seedArtifact(t, model, "status=final")
	stdout, _, err := runCmd(t, "flow", "resolve", "--model", model,
		"--artifact", artifactBinding(flowStateRole, advanced),
		"--outcome", "advance", "--as=json")
	if err == nil {
		t.Fatalf("the fixture `resolve` SUCCEEDED; this oracle needs a real "+
			"refusal envelope to pipe\n%s", stdout)
	}
	refusal := strings.TrimSpace(stdout)

	// The shape that makes this oracle sharp, pinned rather than assumed:
	// the refusal envelope is the BARE `CLIError` (REQ-8 — `findings` is a
	// top-level sibling of `code`, with no wrapper), so it carries NO `type`
	// key at all. A decoder that discriminated on `type` alone would take
	// this document for a bare `data` object, find no `writes` in it, and
	// report a NO-OP SUCCESS for a transition the model refused. `code` is
	// what separates them — `resolvePayload` has no such member.
	if strings.Contains(refusal, `"type"`) {
		t.Fatalf("the refusal envelope now carries a `type` key; this oracle "+
			"exists for the spelling that does NOT, which is the one a "+
			"`type`-only discriminator mistakes for a plan:\n%s", refusal)
	}
	if !strings.Contains(refusal, `"code"`) {
		t.Fatalf("the refusal envelope carries no `code`; nothing "+
			"distinguishes it from a bare `data` object:\n%s", refusal)
	}

	ce := requireRefusalStdinC3xz(t, refusal, codeWriteInvalid, 2,
		"flow", "set-state", "--model", model,
		"--artifact", bind, "--plan", "-", "--as=json")
	if ce.Param != planFlagName {
		t.Errorf("param = %q; want %q — the defect is the DOCUMENT handed to "+
			"--plan, and a caller must be told which input to fix",
			ce.Param, planFlagName)
	}

	// The transition did not happen: a refused plan applies nothing.
	state := flowData(t, requireSuccess(t, "flow", "read-state",
		"--model", model, "--artifact", bind,
		"--artifact", artifactBinding(flowOrphanRole, art), "--as=json"))
	if got := readerTagValue(t, state, "state", "status"); got != "draft" {
		t.Errorf("status = %#v after a REFUSED plan; want %q — every input "+
			"refusal precedes the accessor (REQ-61)", got, "draft")
	}
}

// An explicit `{"type":"failed"}` envelope refuses too, and by its TYPE
// rather than by happening to carry no writes. The oracle feeds a document
// that DOES carry a plausible `writes` block, so an implementation that
// checked only for emptiness would apply it.
//
// ADVERSARIAL
func TestC3xz_AFailedTypeEnvelopeRefusesEvenCarryingWrites(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	poisoned := `{"type":"failed","data":{"writes":{"status":"final"},"clear":[]}}`

	ce := requireRefusalStdinC3xz(t, poisoned, codeWriteInvalid, 2,
		"flow", "set-state", "--model", model,
		"--artifact", bind, "--plan", "-", "--as=json")
	if ce.Param != planFlagName {
		t.Errorf("param = %q; want %q", ce.Param, planFlagName)
	}

	state := flowData(t, requireSuccess(t, "flow", "read-state",
		"--model", model, "--artifact", bind,
		"--artifact", artifactBinding(flowOrphanRole, art), "--as=json"))
	if got := readerTagValue(t, state, "state", "status"); got != "draft" {
		t.Errorf("status = %#v; a `failed` envelope's `writes` was APPLIED. "+
			"The type discriminates, not the emptiness of the block", got)
	}
}

// --- merge precedence ----------------------------------------------------

// `--write` alongside `--plan` is admitted for keys the plan leaves UNSET
// and refused for keys the plan SETS — the latter under the EXISTING
// duplicate code, because two mutations for one key have no defined order
// whichever carrier they arrived on.
//
// Both directions are asserted in one test: the refusal alone is satisfied
// by an implementation that refuses every `--write` given beside a `--plan`,
// which would make the composition useless.
//
// INPUT EDGE
func TestC3xz_WriteBesideAPlanRefusesOnAPlanSetKeyAndAppliesOnAnUnsetOne(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)

	src := seedArtifact(t, model, "status=draft", "stale=obsolete")
	envelope := resolvePlanC3xz(t, model, artifactBinding(flowStateRole, src),
		"advance")

	t.Run("plan-set-key-refuses-flow-write-duplicate", func(t *testing.T) {
		art := seedArtifact(t, model, "status=draft", "stale=obsolete")

		// The plan sets `status`; naming it again is ambiguous.
		ce := requireRefusalStdinC3xz(t, envelope, codeWriteDuplicate, 2,
			"flow", "set-state", "--model", model,
			"--artifact", artifactBinding(flowStateRole, art),
			"--plan", "-", "--write", "status=draft", "--as=json")
		if ce.Param != "status" {
			t.Errorf("param = %q; want %q — the duplicate refusal names the "+
				"KEY in collision", ce.Param, "status")
		}
	})

	t.Run("plan-cleared-key-refuses-too", func(t *testing.T) {
		art := seedArtifact(t, model, "status=draft", "stale=obsolete")

		// The plan CLEARS `stale`; setting it is the same collision across
		// the two carriers that `--write` + `--clear` already refuses.
		requireRefusalStdinC3xz(t, envelope, codeWriteDuplicate, 2,
			"flow", "set-state", "--model", model,
			"--artifact", artifactBinding(flowStateRole, art),
			"--plan", "-", "--write", "stale=x", "--as=json")
	})

	t.Run("plan-unset-key-composes", func(t *testing.T) {
		art := seedArtifact(t, model, "status=draft", "stale=obsolete")

		// `labels` is owned, served by the same writer, and the plan leaves
		// it unset — a caller-composed value the table cannot know.
		result := applyResultOfC3xz(t, requireSuccessStdinC3xz(t, envelope,
			"flow", "set-state", "--model", model,
			"--artifact", artifactBinding(flowStateRole, art),
			"--plan", "-", "--write", "labels="+canonicalSetLiteral,
			"--as=json"))

		// The plan's own key still landed…
		if result.owned["status"] != "final" {
			t.Errorf("owned[status] = %#v; the plan's write must still apply "+
				"beside a composed one", result.owned["status"])
		}
		// …and so did the composed one.
		if _, held := result.writes["labels"]; !held {
			t.Errorf("`labels` is absent from writes = %#v; a key the plan "+
				"leaves UNSET is admitted from --write, or --plan and "+
				"--write cannot be composed at all", result.writes)
		}
		assertCanonicalSetValue(t, "labels", result.owned["labels"])
	})
}

// A plan's `clear[]` naming a key no writer serves refuses under
// `flow-clear-unbound`, NOT `flow-write-unbound`. The two codes are what
// tell a caller which half of the plan to fix, and routing a plan's clears
// through the write arm would report the wrong one.
//
// INPUT EDGE
func TestC3xz_APlanClearOnAnUnservedKeyIsClearUnbound(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")

	// `note` is owned but served by NO writer — the fixture's deliberate
	// unbound key, and the same one `--clear note` refuses on today.
	plan := `{"type":"ok","data":{"writes":{},"clear":["note"]}}`

	ce := requireRefusalStdinC3xz(t, plan, codeClearUnbound, 2,
		"flow", "set-state", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--plan", "-", "--as=json")
	if ce.Param != "note" {
		t.Errorf("param = %q; want %q — the binding refusal names the KEY",
			ce.Param, "note")
	}
}

// A plan's `writes` naming a key no writer serves refuses under
// `flow-write-unbound` — the write arm, the mirror of the clause above.
//
// INPUT EDGE
func TestC3xz_APlanWriteOnAnUnservedKeyIsWriteUnbound(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")

	plan := `{"type":"ok","data":{"writes":{"note":"x"},"clear":[]}}`

	requireRefusalStdinC3xz(t, plan, codeWriteUnbound, 2,
		"flow", "set-state", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--plan", "-", "--as=json")
}

// --- the malformed plan --------------------------------------------------

// A `<clear>` sentinel appearing inside `data.writes` is a MALFORMED PLAN,
// not a removal. `resolve` splits removals into `clear[]` and never emits
// the sentinel as a write value, so a document carrying one did not come
// from `resolve` unaltered — and silently treating it as a clear would make
// the sentinel authorable through a back door the `--write` grammar closes
// deliberately (REQ-62).
//
// ADVERSARIAL
func TestC3xz_AClearSentinelInsideAPlansWritesIsAMalformedPlan(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft", "stale=obsolete")
	bind := artifactBinding(flowStateRole, art)

	plan := `{"type":"ok","data":{"writes":{"stale":"<clear>"},"clear":[]}}`

	ce := requireRefusalStdinC3xz(t, plan, codeWriteInvalid, 2,
		"flow", "set-state", "--model", model,
		"--artifact", bind, "--plan", "-", "--as=json")
	if ce.Param != planFlagName {
		t.Errorf("param = %q; want %q — the defect is the plan DOCUMENT's "+
			"shape, not the key's value", ce.Param, planFlagName)
	}

	// And it was not quietly applied as a removal.
	state := flowData(t, requireSuccess(t, "flow", "read-state",
		"--model", model, "--artifact", bind,
		"--artifact", artifactBinding(flowOrphanRole, art), "--as=json"))
	if _, held := readerTag(t, state, "state", "stale"); !held {
		t.Error("`stale` was REMOVED by a `<clear>` sentinel in the plan's " +
			"writes{}; the sentinel is unauthorable as a write value on " +
			"every carrier, and a plan's removals ride clear[]")
	}
}

// Malformed plan documents refuse under `flow-write-invalid` on param
// `plan`, never by panicking, succeeding vacuously, or minting a new code.
// The refusal table RDR 0005 fixes is CLOSED, so a new carrier reuses a
// published code rather than widening it.
//
// ADVERSARIAL
func TestC3xz_AMalformedPlanDocumentRefusesUnderAPublishedCode(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	for _, tc := range []struct {
		name string
		body string
	}{
		{"empty", ""},
		{"whitespace-only", "   \n\t\n"},
		{"not-json", "this is not a plan"},
		{"a-json-array", `["status=final"]`},
		{"envelope-with-no-data", `{"type":"ok"}`},
		// `null` in either position. `json.Unmarshal` takes the literal
		// `null` into a struct as a NO-OP, so both of these decode without
		// error to a zero plan and would report a no-op SUCCESS for a
		// document that carries no plan at all. `jq .data` over a refusal
		// prints exactly `null` — a bare `CLIError` has no `data` member —
		// so this is the piped-refusal hazard reached by a second spelling.
		{"a-bare-null-document", `null`},
		{"an-envelope-whose-data-is-null", `{"type":"ok","data":null}`},
		{"a-bare-json-scalar", `"status=final"`},
		{"data-is-a-json-scalar", `{"type":"ok","data":"status=final"}`},
		// A PRESENT but null/non-string `type`. Decoded as `*string` these
		// are nil — indistinguishable from an ABSENT key — so the document
		// would take the bare-`data` arm. The write-carrying arm below is
		// the sharp one: it must never silently drop a real write.
		{"type-is-null", `{"type":null,"data":null}`},
		{"type-is-not-a-string", `{"type":7,"data":{"writes":{"status":"final"}}}`},
		{"writes-is-not-an-object", `{"type":"ok","data":{"writes":[]}}`},
		{"a-write-value-is-not-a-string", `{"type":"ok","data":{"writes":{"status":7}}}`},
		{"clear-is-not-an-array", `{"type":"ok","data":{"clear":"stale"}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ce := requireRefusalStdinC3xz(t, tc.body, codeWriteInvalid, 2,
				"flow", "set-state", "--model", model,
				"--artifact", bind, "--plan", "-", "--as=json")
			if ce.Param != planFlagName {
				t.Errorf("param = %q; want %q", ce.Param, planFlagName)
			}
			if ce.Message == "" {
				t.Error("the refusal carries no message; `message` must " +
					"stand alone without any structured field being consulted")
			}
		})
	}

	// An unreadable path is the same class: the CLI owns the read, so a
	// missing plan file is this flag's refusal and never a load category.
	requireRefusal(t, codeWriteInvalid, 2,
		"flow", "set-state", "--model", model, "--artifact", bind,
		"--plan", filepath.Join(t.TempDir(), "does-not-exist.json"),
		"--as=json")
}

// A present-but-null `type` never silently DROPS the plan's writes.
//
// This is the sharpest arm of the null family and the reason `Type` is raw
// bytes rather than `*string`: decoded as a pointer, an absent `type` and a
// present `"type":null` are both nil. A document that carries real writes
// under a null discriminator would then take the bare-`data` arm, look for
// `writes` at the TOP level, find none, and report a no-op success — the
// caller's transition silently discarded while the exit code says it
// applied. A refusal is the only safe answer; applying the nested writes
// would be guessing at a shape the document does not have.
//
// ADVERSARIAL
func TestC3xz_ANullEnvelopeTypeNeverSilentlyDropsThePlansWrites(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	const carriesARealWrite = `{"type":null,"data":{"writes":{"status":"final"}}}`

	ce := requireRefusalStdinC3xz(t, carriesARealWrite, codeWriteInvalid, 2,
		"flow", "set-state", "--model", model, "--artifact", bind,
		"--plan", "-", "--as=json")
	if ce.Param != planFlagName {
		t.Errorf("param = %q; want %q", ce.Param, planFlagName)
	}

	// The write must NOT have landed. A refusal that still mutated the
	// artifact would be worse than the silent drop it replaces.
	state := flowData(t, requireSuccess(t, "flow", "read-state",
		"--model", model, "--artifact", bind,
		"--artifact", artifactBinding(flowOrphanRole, art), "--as=json"))
	if got := readerTagValue(t, state, "state", "status"); got != "draft" {
		t.Errorf("status = %#v after a REFUSED null-type plan; want %q — the "+
			"plan's nested write must neither apply nor vanish silently",
			got, "draft")
	}
}

// Another verb's success envelope is not a plan.
//
// `flow read-state --as=json` is a well-formed `{"type":"ok","data":…}`
// document carrying `model` / `artifacts` / `readers` and no `writes`, so a
// carrier that only asks "did this decode?" takes it for an empty plan and
// reports a no-op SUCCESS for a pipeline that never produced a plan. That is
// the wrong-verb sibling of the piped-refusal hole, and it is reachable by
// typing the wrong verb into a real pipeline.
//
// The oracle pins the DISCRIMINATION, not a field allowlist: RDR 0023
// (§1195-1210) rejected a per-verb pickable-field universe, so unknown and
// future keys must keep flowing through. The `{}` arm is the fence that
// keeps this from becoming one.
//
// ADVERSARIAL
func TestC3xz_AnotherVerbsEnvelopeIsNotAPlan(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	// A REAL `read-state` envelope, not a hand-written lookalike, so the
	// oracle tracks whatever that verb actually emits.
	readState := requireSuccess(t, "flow", "read-state", "--model", model,
		"--artifact", bind,
		"--artifact", artifactBinding(flowOrphanRole, art), "--as=json")

	ce := requireRefusalStdinC3xz(t, readState, codeWriteInvalid, 2,
		"flow", "set-state", "--model", model, "--artifact", bind,
		"--plan", "-", "--as=json")
	if ce.Param != planFlagName {
		t.Errorf("param = %q; want %q", ce.Param, planFlagName)
	}

	// The FENCE. An empty plan is still an accepted no-op: that is the
	// decided `none` / `stopped:*` behavior, and a fix that refused it
	// would relax nothing but would break the row this kata must carry.
	requireSuccessStdinC3xz(t, `{"type":"ok","data":{}}`,
		"flow", "set-state", "--model", model, "--artifact", bind,
		"--plan", "-", "--as=json")

	// A plan carrying ONLY `clear` is a plan; the check takes either key.
	requireSuccessStdinC3xz(t, `{"type":"ok","data":{"clear":[]}}`,
		"flow", "set-state", "--model", model, "--artifact", bind,
		"--plan", "-", "--as=json")

	// Append-only evolution keeps working: an UNKNOWN sibling key alongside
	// a real `writes` is tolerated, never refused. This is the arm that
	// fails if the fix is ever tightened into a field allowlist.
	requireSuccessStdinC3xz(t,
		`{"type":"ok","data":{"writes":{"status":"final"},"a_future_key":42}}`,
		"flow", "set-state", "--model", model, "--artifact", bind,
		"--plan", "-", "--as=json")
}

// An explicitly EMPTY `--plan` refuses; only an OMITTED one is the
// flag-driven path.
//
// `--plan=` and `--plan "$UNSET_VAR"` are what a caller writes when the
// path was supposed to come from a variable that never got set. Reading the
// flag's VALUE cannot tell that from "no --plan at all", so both degrade to
// a silent no-op success — the caller believes a plan applied when none was
// ever read. `Changed` separates them, which is the whole fix.
//
// ADVERSARIAL
func TestC3xz_AnExplicitlyEmptyPlanRefusesWhileAnOmittedOneDoesNot(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	for _, empty := range []string{"", "   "} {
		ce := requireRefusal(t, codeWriteInvalid, 2,
			"flow", "set-state", "--model", model, "--artifact", bind,
			"--plan", empty, "--as=json")
		if ce.Param != planFlagName {
			t.Errorf("--plan %q: param = %q; want %q", empty, ce.Param,
				planFlagName)
		}
	}

	// The fence: OMITTING `--plan` is still the ordinary flag-driven path
	// and must NOT be caught by the refusal above. Without this arm the fix
	// could be "always refuse" and the suite would not notice.
	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", bind, "--write", "status=final", "--as=json")
}

// --- the empty plan ------------------------------------------------------

// An empty `writes` with an absent `clear` — a `none` / `stopped:*` row —
// is a NO-OP SUCCESS at exit 0, which is what `set-state` already does for a
// request carrying no mutations. Nothing is relaxed here and no refusal is
// invented: the caller branches on the plan's own content, and a refusal
// would force them to pre-inspect the very document they handed over.
//
// `clear[]` renders as `[]` and never `null`, so a consumer can iterate it
// without a nil check.
//
// BOUNDARY
func TestC3xz_AnEmptyPlanIsANoOpSuccessAtExitZero(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	for _, tc := range []struct {
		name string
		body string
	}{
		{"empty-writes-and-absent-clear", `{"type":"ok","data":{"writes":{}}}`},
		{"absent-both", `{"type":"ok","data":{}}`},
		{"bare-data-empty-object", `{}`},
		{"explicitly-null-members", `{"type":"ok","data":{"writes":null,"clear":null}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stdout := requireSuccessStdinC3xz(t, tc.body,
				"flow", "set-state", "--model", model,
				"--artifact", bind, "--plan", "-", "--as=json")

			data := flowData(t, stdout)
			writes, ok := data["writes"].(map[string]any)
			if !ok {
				t.Fatalf("`writes` is not an object: %#v", data["writes"])
			}
			if len(writes) != 0 {
				t.Errorf("writes = %#v; an empty plan plans nothing", writes)
			}
			// `[]` and not `null`: the existing normalization must survive
			// the plan path.
			raw, held := data["clear"]
			if !held {
				t.Fatal("`clear` is absent from the payload; the data " +
					"minimum fixes it")
			}
			if raw == nil {
				t.Error("`clear` is null; it renders as the empty array `[]` " +
					"so a consumer can iterate it without a nil check")
			}
			if arr, ok := raw.([]any); !ok || len(arr) != 0 {
				t.Errorf("clear = %#v; want the empty array", raw)
			}
		})
	}

	// The state is untouched by a no-op, which is what makes it a no-op
	// rather than a write of empty values.
	state := flowData(t, requireSuccess(t, "flow", "read-state",
		"--model", model, "--artifact", bind,
		"--artifact", artifactBinding(flowOrphanRole, art), "--as=json"))
	if got := readerTagValue(t, state, "state", "status"); got != "draft" {
		t.Errorf("status = %#v after a no-op plan; want %q", got, "draft")
	}
}

// --- canonicalisation (REQ-71) -------------------------------------------

// A set-valued plan write handed in UNSORTED and DUPLICATED form is
// re-canonicalised, not trusted verbatim, and reads back byte-equal to the
// canonical literal.
//
// This is the property `canonicalSet`'s doc comment names as what makes
// plan-to-request copy-through byte equality. A carried plan is exactly the
// path on which a caller might hand back a literal they edited, so the
// re-canonicalisation must bite HERE and not only on `--write`.
//
// DOMAIN EDGE
func TestC3xz_ASetValuedPlanWriteIsRecanonicalisedNotTrusted(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	// Unsorted AND duplicated, over the two members the MVV fixes — one
	// carrying `<`, one carrying `&`.
	unsorted := `["x&y","a<b","x&y"]`
	if unsorted == canonicalSetLiteral {
		t.Fatal("the fixture literal is already canonical; this oracle " +
			"needs a non-canonical spelling to have anything to fix")
	}
	plan := `{"type":"ok","data":{"writes":{"labels":` +
		jsonOfC3xz(t, unsorted) + `}}}`

	result := applyResultOfC3xz(t, requireSuccessStdinC3xz(t, plan,
		"flow", "set-state", "--model", model,
		"--artifact", bind, "--plan", "-", "--as=json"))

	if got := result.writes["labels"]; got != canonicalSetLiteral {
		t.Errorf("writes[labels] = %#v; want the canonical %q — a plan's set "+
			"value is RE-canonicalised (members sorted, deduplicated, "+
			"compact) rather than echoed, which is what makes read-back "+
			"equality byte equality (REQ-71)", got, canonicalSetLiteral)
	}

	// And read-back agrees, through an independent call.
	state := flowData(t, requireSuccess(t, "flow", "read-state",
		"--model", model, "--artifact", bind,
		"--artifact", artifactBinding(flowOrphanRole, art), "--as=json"))
	assertCanonicalSetValue(t, "labels",
		readerTagRaw(t, state, "state", "labels"))

	// The HTML-escaped spelling must appear nowhere: `<` and `&` serialize
	// as themselves on every emit site (REQ-70, REQ-72).
	if strings.Contains(jsonOfC3xz(t, result.writes), escapedSetLiteral) {
		t.Errorf("the applied write carries the HTML-ESCAPED set literal %s; "+
			"the canonical encoder disables HTML escaping", escapedSetLiteral)
	}
}

// A plan value that does not conform to its declaration refuses, exactly as
// the same value would through `--write`. Nothing about having come from a
// plan makes a value legal: legality is re-derived from the REQUEST, which
// is what keeps 0005's "nothing links a set-state request to a prior
// resolve" true of a carried plan.
//
// ADVERSARIAL
func TestC3xz_APlanValueIsStillHeldToItsDeclaration(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	// `status` is an enum over {draft, final}. A plan asserting anything
	// else is refused — a stale plan, or one for another model, cannot
	// smuggle a value past the declaration.
	plan := `{"type":"ok","data":{"writes":{"status":"archived"}}}`

	requireRefusalStdinC3xz(t, plan, codeWriteInvalid, 2,
		"flow", "set-state", "--model", model,
		"--artifact", bind, "--plan", "-", "--as=json")

	// A set-valued key handed a bare scalar is the wrong-KIND refusal, the
	// other half of "well-formed for its declared kind".
	requireRefusalStdinC3xz(t, `{"type":"ok","data":{"writes":{"labels":"plain"}}}`,
		codeWriteInvalid, 2,
		"flow", "set-state", "--model", model,
		"--artifact", bind, "--plan", "-", "--as=json")
}
