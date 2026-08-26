package cli

// RDR 0008 — the CLI-side citations JDR 0001 §D8 and §D10 owe this Stage
// (deviations D3).
//
// Neither is new behavior. `internal/cli/flow_input.go:196-200` already
// refuses `--tag recognized=`, and `clierr.Finding` already carries
// `Code`/`Param`/`Hint`. What this file adds is the PIN: the refusal is a
// `GroupUserEnv` CLI refusal at the CLI site, NOT block 4's Go error path,
// and `--tag` never routes a value into `Input.Recognized`.
//
// These tests are green against the predecessor by design. Their job is to
// fail if a future change re-routes the refusal onto the kernel's
// programmer-mistake path or lets `--tag` reach the recognized channel.

import (
	"testing"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
	"github.com/newcoinc/intrastate/internal/table"
)

// REQ-104: JDR 0001 §D8 "A `--tag` naming an **owned** key or `recognized` is
// refused at the CLI before any accessor runs (P2 — refuse rather than
// silently shadow under owned-over-observed precedence); both are
// `GroupUserEnv` with their own codes (§D11). 0008's \"programmer mistake\" is
// the CLI caller's, which from the CLI's seat is the user; `GroupInternal` is
// for the CLI's own invariants."
// ADVERSARIAL
//
// The site and the group are the whole claim. A `--tag recognized=` is refused
// at the CLI with a `GroupUserEnv` code — never routed to `GroupInternal`, and
// never allowed through to become the kernel predicate's Go error, which would
// present a user mistake as a programmer one.
func TestReq104_TagOnTheReservedKeyIsAGroupUserEnvRefusalAtTheCLI(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")

	ce := requireRefusal(t, codeTagReserved, 2,
		"flow", "next", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--tag", table.RecognizedTagKey+"=whatever", "--as=json")

	if ce.Group != clierr.GroupUserEnv {
		t.Errorf("group = %v; want GroupUserEnv — from the CLI's seat the "+
			"reserved-key mistake is the user's, and GroupInternal is reserved "+
			"for the CLI's own invariants", ce.Group)
	}
	if ce.Param != table.RecognizedTagKey {
		t.Errorf("param = %q; want the offending key %q",
			ce.Param, table.RecognizedTagKey)
	}
}

// REQ-105: JDR 0001 §D8 "`--tag` stays `Observed`."
// BOUNDARY
//
// The CLI must not route a `--tag` value into `Input.Recognized` or
// `Input.Owned`. Asserted on the CLI's own tag-collection seam: every tag it
// collects is observed context, and the reserved key never survives it — which
// is precisely what keeps block 4's producer obligation satisfied by the CLI.
func TestReq105_TagValuesStayObservedAndNeverReachTheRecognizedChannel(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	bind := artifactBinding(flowStateRole, art)

	// An ordinary undeclared observed key passes through as context.
	requireSuccess(t, "flow", "next", "--model", model, "--artifact", bind,
		"--tag", "sidechannel=ctx", "--as=json")

	// The reserved key does not, at any value, including the empty one.
	for _, value := range []string{"", "advance", "anything"} {
		requireRefusal(t, codeTagReserved, 2,
			"flow", "next", "--model", model, "--artifact", bind,
			"--tag", table.RecognizedTagKey+"="+value, "--as=json")
	}
}

// REQ-103: JDR 0001 §D10 "`flow-model-invalid` (every 0002 load category,
// 0003's two, 0008's `reserved_tag_key`) | UserEnv / 2 | `findings[]`"
// DOMAIN EDGE
//
// `reserved_tag_key` maps to the shared model-load CLI code, not to a code of
// its own: it is one of "every 0002 load category" in that row. A model whose
// declaration takes the reserved key under a wrong provenance must refuse
// `flow-model-invalid` at exit 2.
func TestReq103_ReservedTagKeyMapsToTheSharedModelInvalidCodeAtExit2(t *testing.T) {
	model := writeFlowModel(t, reservedKeyModel)

	ce := requireRefusal(t, "flow-model-invalid", 2,
		"flow", "next", "--model", model, "--as=json")

	if ce.Group != clierr.GroupUserEnv {
		t.Errorf("group = %v; want GroupUserEnv / exit 2", ce.Group)
	}
}

// REQ-102: JDR 0001 §D10 "**0007's reopening is accepted: `CLIError` gains
// exactly one `omitempty` structured field, and it is the `Finding` record
// 0006 already mandates** — non-normative `Findings []clierr.Finding`
// (`json:\"findings,omitempty\"`), `Finding{Code, Message, Param, Locator,
// Hint}`, subsystem-agnostic. One field serves all four carriers: … 0008's
// per-key `reserved_tag_key` with the near-miss advisory in `hint`"
// BOUNDARY
//
// The CARRIER's shape is the claim: one `Findings []clierr.Finding` field, and
// a `Finding` carrying the five named fields including `Hint`, which is where
// §D10 places the near-miss advisory. This RDR does not add the CLI delivery
// (deviation D5 defers that to RDR 0005/0006) — it pins that the carrier §D10
// names exists and has room for the advisory.
func TestReq102_TheFindingCarrierHasTheFiveFieldsIncludingHint(t *testing.T) {
	f := clierr.Finding{
		Code:    string(table.CatReservedTagKey),
		Message: "a tag declaration takes the reserved key",
		Param:   table.RecognizedTagKey,
		Locator: "model.toml:67",
		Hint:    "did you mean `recognized`?",
	}

	ce := &clierr.CLIError{
		Code:     "flow-model-invalid",
		Group:    clierr.GroupUserEnv,
		Findings: []clierr.Finding{f},
	}

	if len(ce.Findings) != 1 {
		t.Fatalf("Findings carries %d records; want the one §D10 names",
			len(ce.Findings))
	}
	got := ce.Findings[0]
	if got.Code != string(table.CatReservedTagKey) {
		t.Errorf("Code = %q; want %q", got.Code, table.CatReservedTagKey)
	}
	if got.Hint == "" {
		t.Error("Hint is empty; §D10 places 0008's near-miss advisory there")
	}
	if got.Param != table.RecognizedTagKey {
		t.Errorf("Param = %q; want %q", got.Param, table.RecognizedTagKey)
	}
}

// reservedKeyModel declares an OWNED tag under the reserved key, so RDR 0002's
// loader refuses it in the `reserved_tag_key` category. It is the smallest
// model that reaches that category through the CLI.
const reservedKeyModel = `outcomes = ["advance"]

[model]
id = "reserved"
version = 1

[initial]
status = "draft"

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true
required = true

[tags.recognized]
provenance = "owned"
kind = "enum"
domain = ["advance"]
single_valued = true
required = true

[read.state]
role = "state"
path = "state.status"
keys = ["status"]
timeout = "2s"

[write.state]
role = "state"
path = "state.status"
keys = ["status"]
timeout = "2s"
read_back = true

[[rule]]
id = "advance"
source = "reserved:advance"
[rule.match.recognized]
eq = "advance"
[rule.match.status]
eq = "draft"
[rule.write]
status = "final"
`
