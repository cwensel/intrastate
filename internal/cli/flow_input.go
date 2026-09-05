package cli

// RDR 0005 — model selection and the input grammar every `flow` verb
// shares: `--tag`, `--artifact`, `--write`, and `--clear`.
//
// Everything in this file runs BEFORE any accessor does. That is a
// contract obligation, not an optimisation: REQ-27 and REQ-61 both say
// "before any accessor runs", and the Phase 1 suite proves the ordering by
// leaving an artifact role unbound on invocations whose input is also
// wrong — if an accessor ran first, the refusal would name the wrong
// subject.

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/cli/flowbind"
	"github.com/cwensel/intrastate/internal/table"
	"github.com/spf13/cobra"
)

// --- the stable code table (`0005:FM`) -----------------------------------
//
// The spellings are normative and the group fixes the exit. They are
// declared as constants in one place so a producer cannot spell one
// slightly differently at a second site.

const (
	codeTagInvalid      = "flow-tag-invalid"
	codeTagDuplicate    = "flow-tag-duplicate"
	codeTagReserved     = "flow-tag-reserved"
	codeTagOwned        = "flow-tag-owned"
	codeWriteInvalid    = "flow-write-invalid"
	codeWriteDuplicate  = "flow-write-duplicate"
	codeWriteUnbound    = "flow-write-unbound"
	codeClearUnbound    = "flow-clear-unbound"
	codeArtifactInvalid = "flow-artifact-invalid"
	codeArtifactMissing = "flow-artifact-missing"
	codeModelNotFound   = "flow-model-not-found"
	codeModelInvalid    = "flow-model-invalid"

	codeGateDenied        = "flow-gate-denied"
	codeGateIndeterminate = "flow-gate-indeterminate"

	codeAccessorTimeout = "flow-accessor-timeout"
	codeAccessorFailed  = "flow-accessor-failed"
	// codeWriteFailedApplied is the applied WRITE refusal's own exit-3 code
	// (JDR 0003 §D4 (b), `0026:C1` `refusal:`). A write that failed only
	// AFTER its command already ran is a different fact from one that could
	// not run at all, and an agent's documented branch on
	// `flow-accessor-failed` is retry — so the two must not share a code.
	// Not-applied execution failures keep `flow-accessor-failed`.
	codeWriteFailedApplied = "flow-write-failed-applied"
	codeReadIncomplete     = "flow-read-incomplete"
	codeAccessorUnknown    = "flow-accessor-unknown"
	codeCapabilityMismat   = "flow-accessor-capability-mismatch"
	codeWriteNonOwned      = "flow-write-non-owned"

	// codeEscapeRowShapeBreach is RDR 0009 `0009:C7`'s stable code for a
	// kernel escape-row shape breach. Note it is deliberately NOT
	// `flow`-prefixed like its neighbours: the contract fixes the literal
	// and REQ-67 makes the Code the conformance oracle, so a prefixed
	// variant would fail the contract. It is not a typo.
	codeEscapeRowShapeBreach = "escape-row-shape-breach"
	// hintEscapeRowShapeBreach is the remedy the record fixes verbatim.
	hintEscapeRowShapeBreach = "fix the table producer: an escape row must " +
		"carry no writes"

	// codeWriteEditRefused is RDR 0028 `0028:C1.3` EXIT GROUP:'s distinct
	// code for a declared-line-edit refusal decided before mutation. Its
	// SPELLING is this stage's to choose and is explicitly non-normative
	// — no test may pin the string; what the contract fixes is the exit
	// GROUP and the `findings[]` carriage.
	codeWriteEditRefused = "flow-write-edit-refused"

	// codeRequestRefused is the sibling of `codeWriteEditRefused` for every
	// declared-request refusal a line edit did NOT mint.
	//
	// The two are split by ORIGIN, not by phase or carrier. C1.6's argv
	// preconditions — an unbound `{tag.<key>}`, a `-`-prefixed bound value
	// — are checked during argv substitution, which runs for readers,
	// gates and command-backed WRITERS alike; the read-back preconditions
	// come from an edit writer and are still about a reader's argv. None
	// of those has a rule `<id>.edit.<key>` to name, so reporting them
	// under the line-edit code invented a subject that does not exist.
	//
	// Its SPELLING is this stage's to choose on the same terms as its
	// sibling's: what C1.3 fixes is the exit GROUP, which is 2 on both.
	codeRequestRefused = "flow-accessor-request-refused"

	codeReadBackMismatch   = "flow-write-readback-mismatch"
	codeReadBackIncomplete = "flow-write-readback-incomplete"
	codeReadBackTimeout    = "flow-write-readback-timeout"

	// An argv token beginning `stopped:` is the shared stop-packet
	// prefix of an upstream tool whose output was command-substituted
	// into this invocation. The check fires at the root gateway
	// (`root.go::argvCarriesUpstreamStop`) rather than in a verb, but
	// the spelling belongs to this one table.
	codeArgvUpstreamStop = "flow-argv-upstream-stop"
)

// userErr builds a `GroupUserEnv` refusal: the request or the model is
// wrong, or the model said no. Exit 2 (REQ-19).
func userErr(code, param, message string) *clierr.CLIError {
	return &clierr.CLIError{
		Code: code, Param: param, Message: message,
		Group: clierr.GroupUserEnv,
	}
}

// envErr builds a `GroupEnvUnavailable` refusal: the environment could not
// be consulted and the same request may be re-run unchanged. Exit 3
// (REQ-19, REQ-20).
func envErr(code, param, message string) *clierr.CLIError {
	return &clierr.CLIError{
		Code: code, Param: param, Message: message,
		Group: clierr.GroupEnvUnavailable,
	}
}

// --- model selection -----------------------------------------------------

// selectModel resolves `--model <path>` or `--flow <id>` to a loaded model.
//
// The three arms of `flow-model-not-found` (REQ-22, REQ-90): neither flag,
// both flags, and a selection that does not resolve. Carrying the selection
// ARITY error under the same code as non-resolution is RDR-owned (CA A6) —
// from the caller's side both mean "this invocation names no one model".
//
// The CLI performs the file I/O and hands the loader BYTES plus a source id
// (REQ-23): `table.Load` does no file I/O, so an unreadable path is this
// command's failure and never a load category.
// selectModelPath resolves the `--model` / `--flow` pair to a path, or to
// the ONE refusal REQ-90 fixes for every arm of the selection.
//
// It is shared by the `flow` verbs and by root `lint` so the two cannot
// disagree about what selection means. `lint` previously carried its own
// copy that answered `flag-required` / `flag-mutually-exclusive`, which
// left a caller branching on `flow-model-not-found` — the code the
// contract names for "neither / both … or the selection does not
// resolve" — with an unhandled refusal from one model-taking command.
func selectModelPath(cmd *cobra.Command) (string, *clierr.CLIError) {
	path, _ := cmd.Flags().GetString("model")
	flowID, _ := cmd.Flags().GetString("flow")

	switch {
	case path != "" && flowID != "":
		return "", userErr(codeModelNotFound, "model",
			"--model and --flow are mutually exclusive; give exactly one")
	case path == "" && flowID == "":
		return "", userErr(codeModelNotFound, "model",
			"model selection requires exactly one of --model <path> or --flow <id>")
	case path == "":
		// `--flow <id>` config discovery. `--model <path>` ships first
		// (`0005:EIA`); an id this build cannot resolve is a REFUSAL under
		// the same code, not a missing feature (A-1).
		return "", userErr(codeModelNotFound, "flow",
			"no model is registered for the flow id "+flowID+
				"; give --model <path> instead")
	}
	return path, nil
}

func selectModel(cmd *cobra.Command) (*table.Model, string, *clierr.CLIError) {
	path, ce := selectModelPath(cmd)
	if ce != nil {
		return nil, "", ce
	}

	src, err := os.ReadFile(path)
	if err != nil {
		// The CLI owns the read. This is a selection failure — the loader
		// never saw these bytes, so it is emphatically not a load category
		// (REQ-23).
		return nil, "", userErr(codeModelNotFound, "model",
			"the model at "+path+" could not be read: "+err.Error())
	}

	model, err := table.Load(src, path)
	if err != nil {
		return nil, "", loadFailure(path, err)
	}
	return model, path, nil
}

// loadFailure maps a loader refusal onto `flow-model-invalid` with one
// findings entry per category hit (REQ-24, REQ-91).
//
// Every load-time category rides ONE CLI code, so the category slug is what
// a caller branches on and it travels as the finding's `code`. `locator` is
// the finding's file:line — the position the loader attributes the defect
// to — which is why REQ-14 needed a field the shipped record did not carry.
func loadFailure(path string, err error) *clierr.CLIError {
	return &clierr.CLIError{
		Code:     codeModelInvalid,
		Message:  "the selected model could not be loaded",
		Group:    clierr.GroupUserEnv,
		Findings: loadFindings(path, err),
	}
}

// loadFindings renders a loader refusal as the findings list every CLI verb
// that loads a model carries (REQ-91, JDR 0001 §D10 item 3). It is shared by
// `flow`'s selectModel and by `lint` so the two surfaces cannot drift: one
// structured field, the category slug as the inner discriminator, and the
// per-category payload the loader populated.
//
// `locator` is `file:<line>` where the loader grounded a source position,
// and `file:1` where it did not. The loader attributes a line only to the
// refusals it can honestly place — a malformed `[tags.<key>]` declaration
// locates at its own header — and reports zero, meaning unknown, for the
// rest. `:1` is the documented fallback for that zero, not a claim about
// the defect: a reader still needs the file to act on, and inventing a line
// would point them at innocent source text.
func loadFindings(path string, err error) []clierr.Finding {
	category, ok := table.CategoryOf(err)
	if !ok {
		// A loader error carrying no category still refuses under this
		// code: it is a defect in the selected model, and inventing a
		// second code for it would break the one-code rule.
		category = "unknown"
	}
	// The line is read before the `Rule` gate below, because a position and
	// the reserved-key payload are independent: any load category may carry
	// a line, only one carries a rename remedy.
	line := 1
	var located *table.Failure
	if errors.As(err, &located) && located.Line > 0 {
		line = located.Line
	}
	finding := clierr.Finding{
		Code:    string(category),
		Message: err.Error(),
		Locator: path + ":" + strconv.Itoa(line),
	}

	// RDR 0008 `0008:C3` — the three-field `reserved_tag_key` payload rides
	// the same record rather than a category-specific envelope: `Rule` is the
	// direction identifier a consumer branches on, `Param` the offending name
	// as authored, `Hint` the remedy.
	//
	// The rename hint is gated on `Rule`, not on the errors.As alone: every
	// load category is a `*table.Failure`, but only this one carries a
	// per-key payload, and "rename the declaration" is nonsense advice on a
	// malformed-TOML refusal.
	var f *table.Failure
	if errors.As(err, &f) && f.Rule != "" {
		finding.Rule = f.Rule
		finding.Param = f.Offending
		// The empty `Remedy` is documented, not missing: for
		// RuleAuthorMustRename the remedy is "choose any other name", and
		// REQ-30 forbids a renderer presenting the reserved key as the
		// required one. The offending name already travels on `param`, so
		// the hint names NO name at all rather than echoing the reserved key
		// back as if it were the answer.
		if f.Remedy != "" {
			finding.Hint = "rename the declaration to `" + f.Remedy + "`"
		} else {
			finding.Hint = "rename the declaration to any other name"
		}
	}
	return []clierr.Finding{finding}
}

// --- the carried plan (`--plan <file|->`) --------------------------------

// planFlagName is the flag `set-state` reads a carried plan from. RDR 0005
// deferred it as "a carried plan artifact (`--plan <file|->`) … a seed, not
// part of this contract"; this is that seed, landed as a flag on an EXISTING
// verb rather than a fifth one (`0005:C1` is untouched).
//
// It is registered on `set-state` ALONE, so it does not disturb RDR 0023's
// `--plan-only`, which rides `resolve` alone. The two never meet on one
// command, and pflag does no long-flag prefix matching — `--plan-only` on
// `set-state` still fails the parse as an unknown flag through the shared
// `command-error` bucket, exactly as `0023:C1`'s fence requires.
const planFlagName = "plan"

// planFlagUsage backticks ONLY the metavariable. pflag reads the first
// backticked span in a usage string as the flag's metavariable name, so
// backticking the prose `flow resolve` rendered the synopsis as
// `--plan flow resolve` in the generated reference.
const planFlagUsage = "apply a flow resolve plan from a `file|-` (`-` is stdin)"

// planStdinSentinel is the conventional spelling for "read the plan from
// stdin", which is what makes `flow resolve … --as json | flow set-state …
// --plan -` a pipe rather than a temporary file.
const planStdinSentinel = "-"

// carriedPlan is the writes and clears a `--plan` document supplies.
//
// The field names and JSON tags mirror `resolvePayload`'s `Writes` / `Clear`
// exactly (`flow_resolve.go`), which is the whole point: `resolve` emits the
// document this decodes, so the copy-through from plan to request is byte
// identical (REQ-71, REQ-108) rather than a second spelling that has to be
// kept in step.
type carriedPlan struct {
	Writes map[string]string `json:"writes"`
	Clear  []string          `json:"clear"`
}

// planEnvelope is the discriminating probe over a carried plan document.
//
// A plan arrives in one of two shapes and BOTH are accepted: the full
// `{"type":"ok","data":…}` envelope `respond.OK` emits, or the bare `data`
// object a caller extracted from it (`jq .data`, say). The discriminator is
// the presence of a top-level `"type"` key, which is unambiguous because a
// bare `data` payload never carries one — `resolvePayload`'s fourteen keys
// are fixed and `type` is not among them.
//
// `Type` and `Data` are `json.RawMessage` rather than decoded values so the
// bare-`data` arm can re-decode the WHOLE document without a second read:
// presence of `Type` selects which bytes are the payload, and nothing is
// parsed twice under a shape it does not have.
//
// `Type` must be RAW, not `*string`. A `*string` is nil for an ABSENT key
// and also for a present `"type":null`, so the two are indistinguishable —
// and `{"type":null,"data":{"writes":…}}` would take the bare-`data` arm,
// find no `writes` at the top level, and SILENTLY DROP a real write while
// reporting success. Raw bytes make presence a byte-level fact, which is
// what the "presence of a top-level `type` key" rule above actually says.
//
// `Code` and `Message` are the REFUSAL envelope's discriminators. A failed
// run under `--as=json` emits the bare `CLIError` — `{"code":…,"message":…}`
// with `findings` as a TOP-LEVEL sibling and no wrapper at all (REQ-8) — so
// a refusal carries NO `type` key and would otherwise land in the bare-`data`
// arm, decode to zero writes, and report a no-op SUCCESS for a transition the
// model declined. `resolvePayload` has no `code` member, so its presence
// separates the two without ambiguity.
type planEnvelope struct {
	Type    json.RawMessage `json:"type"`
	Data    json.RawMessage `json:"data"`
	Code    string          `json:"code"`
	Message string          `json:"message"`
}

// readPlan loads and decodes the `--plan` document, returning the writes and
// clears it carries.
//
// Every refusal here is `flow-write-invalid` on param `plan`. That reuses a
// PUBLISHED code deliberately: RDR 0005's refusal table is closed (see
// `writerFor`'s doc comment), and a malformed carried plan is a malformed
// write request — the same class as `--write status=<clear>` — reaching the
// CLI through a different carrier. Minting a `flow-plan-*` code would widen
// a table the contract fixes, for a caller who already branches on
// `flow-write-invalid` plus `param`.
//
// A REFUSAL envelope on stdin refuses. `flow resolve` writes its refusal to
// the same stream its success goes to under `--as=json`, so a pipe that lost
// its plan carries `{"type":"failed",…}` instead — and applying an empty
// write set from it would report success for a transition the model refused.
// That is the one failure mode a carried plan makes newly reachable, so it
// is refused by name rather than by falling through to "no writes".
func readPlan(cmd *cobra.Command, path string) (*carriedPlan, *clierr.CLIError) {
	var src []byte
	if path == planStdinSentinel {
		read, err := io.ReadAll(cmd.InOrStdin())
		if err != nil {
			return nil, userErr(codeWriteInvalid, planFlagName,
				"the plan could not be read from stdin: "+err.Error())
		}
		src = read
	} else {
		read, err := os.ReadFile(path)
		if err != nil {
			return nil, userErr(codeWriteInvalid, planFlagName,
				"the plan at "+path+" could not be read: "+err.Error())
		}
		src = read
	}

	if len(bytes.TrimSpace(src)) == 0 {
		return nil, userErr(codeWriteInvalid, planFlagName,
			"the plan is empty; --plan takes one `flow resolve --as json` "+
				"envelope or its `data` object")
	}

	// The document must be a JSON OBJECT. `json.Unmarshal` accepts the
	// literal `null` into a struct as a no-op, leaving a zero `planEnvelope`
	// that carries no `type` and no `code` — so a bare `null` would fall
	// through the refusal arm, decode to zero writes, and report a no-op
	// SUCCESS. That is the same failure mode the refusal-envelope arm exists
	// to close, reached by a different spelling: `jq .data` over a refusal
	// prints `null`, since a bare `CLIError` has no `data` member. An empty,
	// malformed, or wrong-kind carrier is `flow-write-invalid` (0005's input
	// family names "malformed, empty, or wrong-kind" for exactly this).
	if !isJSONObject(src) {
		return nil, userErr(codeWriteInvalid, planFlagName,
			"the plan is not one JSON object; --plan takes one "+
				"`flow resolve --as json` envelope or its `data` object")
	}

	var env planEnvelope
	if err := json.Unmarshal(src, &env); err != nil {
		return nil, userErr(codeWriteInvalid, planFlagName,
			"the plan is not one JSON object: "+err.Error())
	}

	// A REFUSAL, in either of its two spellings, before anything is decoded
	// as a plan. The bare `CLIError` arm is the load-bearing one: it carries
	// no `type`, so without this check it would fall through to the
	// bare-`data` arm, decode to zero writes, and report a NO-OP SUCCESS for
	// a transition the model declined — the one failure mode a piped plan
	// makes newly reachable, since `resolve` writes success and refusal to
	// the same stream.
	typePresent := len(bytes.TrimSpace(env.Type)) > 0
	if !typePresent && env.Code != "" {
		ce := userErr(codeWriteInvalid, planFlagName,
			"the plan is a refusal envelope (`"+env.Code+"`), not a plan; a "+
				"refusal carries no writes to apply")
		ce.Hint = "the `flow resolve` that produced this refused; resolve " +
			"that refusal first"
		return nil, ce
	}

	// The full envelope. A non-success type is a REFUSAL that was piped in
	// place of a plan, and there is nothing in it to apply.
	payload := src
	if typePresent {
		// A PRESENT `type` must be the success string. Anything else — a
		// refusal spelling, `null`, a number, an object — is not a plan this
		// applies. Decoding into a string first means a non-string `type`
		// refuses by name rather than being mistaken for an absent key.
		//
		// The target is `*string`, not `string`: `json.Unmarshal` takes the
		// literal `null` into a `string` as a no-op, leaving `""`, which
		// would report `"type":null` as an empty-named ENVELOPE and hand
		// back the resolve-the-refusal hint for a document that is not a
		// refusal. A pointer keeps "absent value" distinct from "the empty
		// string", so the null lands on the non-string arm where it belongs.
		var envType *string
		if err := json.Unmarshal(env.Type, &envType); err != nil || envType == nil {
			return nil, userErr(codeWriteInvalid, planFlagName,
				"the plan's `type` is not a string; --plan takes one "+
					"`flow resolve --as json` envelope or its `data` object")
		}
		if *envType != planEnvelopeOK {
			ce := userErr(codeWriteInvalid, planFlagName,
				"the plan is a `"+*envType+"` envelope, not a plan; a "+
					"refusal carries no writes to apply")
			ce.Hint = "resolve the refusal first; `--plan` applies only a " +
				"successful plan"
			return nil, ce
		}
		// `data` must itself be an object. A literal `{"data":null}` passes
		// the emptiness check above (`null` is four bytes) and unmarshals
		// into `carriedPlan` as a no-op, which is the same silent-success
		// hole the top-level check closes.
		if len(bytes.TrimSpace(env.Data)) == 0 || !isJSONObject(env.Data) {
			return nil, userErr(codeWriteInvalid, planFlagName,
				"the plan envelope carries no `data` object")
		}
		payload = env.Data
	}

	var plan carriedPlan
	if err := json.Unmarshal(payload, &plan); err != nil {
		return nil, userErr(codeWriteInvalid, planFlagName,
			"the plan's `data` is not a plan object: "+err.Error())
	}

	// A NON-EMPTY payload must carry at least one of the two plan keys.
	//
	// Without this, ANOTHER verb's success envelope is a valid empty plan:
	// `flow read-state --as=json` carries `model` / `artifacts` / `readers`
	// and no `writes`, so piping it here reported a no-op SUCCESS for a
	// pipeline that never produced a plan — the wrong-verb sibling of the
	// piped-refusal hole.
	//
	// The check is deliberately a PRESENCE test on two keys, not a field
	// allowlist. RDR 0023 (§1195-1210) rejected "a declared pickable-field
	// list per verb, validation, and its own refusal for an unknown field"
	// as cost without need, and 0005 keeps the envelope append-only, so
	// unknown and future fields must stay tolerated. This adds no opinion
	// about any key other than the two the carrier actually reads.
	//
	// Safe under both projections: `writes` and `clear` carry NO `omitempty`
	// on `resolvePayload` (unlike the echo group) and both sit in the PLAN
	// group `--plan-only` keeps, so a real plan always emits them. An empty
	// `{}` stays an accepted no-op plan — that is the decided `none` /
	// `stopped:*` behavior and is not what this refuses.
	if !isEmptyJSONObject(payload) && !objectHasAnyKey(payload, "writes", "clear") {
		ce := userErr(codeWriteInvalid, planFlagName,
			"the plan carries neither `writes` nor `clear`; --plan takes a "+
				"`flow resolve` plan, not another command's envelope")
		ce.Hint = "pipe `flow resolve … --as json`; an empty plan is `{}`"
		return nil, ce
	}

	// `set-state`'s OWN success envelope carries `writes` and `clear` too,
	// so the presence test above admits it and a caller who piped one
	// set-state into another would silently REPLAY the first one's
	// mutations.
	//
	// `writers` is the discriminator, and it is the ONLY safe one:
	// `setStatePayload` alone declares it (`flow_state.go`), while `owned`
	// — the other tempting candidate — is ALSO a `resolvePayload` field,
	// carried in the ECHO group as `owned,omitempty` and present on every
	// plan that was not projected with `--plan-only`. Discriminating on
	// `owned` refuses real plans.
	//
	// This stays a single-key presence test rather than a field allowlist,
	// so it keeps RDR 0023 §1195-1210's rejection of a per-verb field
	// universe intact and holds no opinion about keys `resolve` may add.
	//
	// Refused rather than tolerated as idempotent: re-applying happens to
	// converge here, but "this envelope came from the wrong verb" is a
	// caller mistake worth naming, and a plan carrier that accepts its own
	// output invites a pipeline that looks like it re-derived a decision
	// when nothing consulted the table.
	if objectHasAnyKey(payload, "writers") {
		ce := userErr(codeWriteInvalid, planFlagName,
			"the plan carries `writers`, which only a `flow set-state` "+
				"envelope has; --plan takes a `flow resolve` plan")
		ce.Hint = "pipe `flow resolve … --as json`, not another `set-state` " +
			"result"
		return nil, ce
	}

	// A `flow next` CANDIDATE also carries `writes` and `clear` and no
	// `writers`, so `jq '.data.candidates[0]'` would otherwise apply here.
	// It must not: a candidate's writes are the PREVIEW a normalized row
	// exposes "without evaluating anything" (`flow_next.go`, REQ-40/A-5),
	// its `unknown` lists facts the run could not decide, and absent
	// `--evaluate-gates` its gates were never run. Applying one would
	// commit a transition that skipped exact-one selection AND gate
	// approval — the two things `resolve` exists to perform, and which
	// 0005 places between a decision and a write.
	//
	// `required` and `unknown` are `candidate`'s alone; neither
	// `resolvePayload` nor `setStatePayload` declares either, so this stays
	// a presence test on foreign keys rather than a field allowlist.
	if foreign := firstPresentKey(payload, "required", "unknown"); foreign != "" {
		ce := userErr(codeWriteInvalid, planFlagName,
			"the plan carries `"+foreign+"`, which only a `flow next` "+
				"candidate has; a candidate is a PREVIEW that skipped "+
				"selection and gates, not a plan")
		ce.Hint = "resolve the outcome first: `flow resolve … --as json`"
		return nil, ce
	}
	return &plan, nil
}

// firstPresentKey returns the first of names present as a member of the
// JSON object src, or "" when none is.
func firstPresentKey(src []byte, names ...string) string {
	var members map[string]json.RawMessage
	if err := json.Unmarshal(src, &members); err != nil {
		return ""
	}
	for _, name := range names {
		if _, ok := members[name]; ok {
			return name
		}
	}
	return ""
}

// isEmptyJSONObject reports whether src is the object `{}`, the one shape
// that legitimately carries no plan keys and still applies as a no-op.
func isEmptyJSONObject(src []byte) bool {
	var members map[string]json.RawMessage
	if err := json.Unmarshal(src, &members); err != nil {
		return false
	}
	return len(members) == 0
}

// objectHasAnyKey reports whether the JSON object src carries any of names
// as a member, regardless of that member's value.
//
// Presence is what matters, not the decoded value: `{"writes":null}` is a
// plan that carries the key and resolves to no writes, while a document
// with no `writes` member at all did not come from `flow resolve`.
func objectHasAnyKey(src []byte, names ...string) bool {
	var members map[string]json.RawMessage
	if err := json.Unmarshal(src, &members); err != nil {
		return false
	}
	for _, name := range names {
		if _, ok := members[name]; ok {
			return true
		}
	}
	return false
}

// isJSONObject reports whether src is a JSON object rather than a scalar,
// an array, or the literal `null`.
//
// It exists because `json.Unmarshal` into a struct silently accepts `null`
// as a no-op, so a struct decode alone cannot tell "an object with no
// matching keys" from "not an object at all" — and the two must refuse
// differently from a plan that legitimately carries no writes.
func isJSONObject(src []byte) bool {
	trimmed := bytes.TrimSpace(src)
	return len(trimmed) > 0 && trimmed[0] == '{'
}

// planEnvelopeOK is the success discriminator `respond.OK` stamps. It is
// spelled here rather than imported so this decoder cannot be made to accept
// a value the emitter does not produce by a change on the other side.
const planEnvelopeOK = "ok"

// --- artifact bindings ---------------------------------------------------

// parseArtifacts parses `--artifact role=path` bindings. Nothing discovers
// artifacts ambiently; location comes ONLY from these (REQ-31, `0004:C3`).
func parseArtifacts(cmd *cobra.Command) (map[string]string, *clierr.CLIError) {
	raw, _ := cmd.Flags().GetStringArray("artifact")
	out := make(map[string]string, len(raw))
	for _, binding := range raw {
		role, path, ok := strings.Cut(binding, "=")
		if !ok || role == "" || path == "" {
			return nil, userErr(codeArtifactInvalid, "artifact",
				"--artifact takes role=path; got "+binding)
		}
		out[role] = path
	}
	return out, nil
}

// --- observed tags -------------------------------------------------------

// parseTags parses `--tag name=value` into observed context.
//
// Three refusals precede any accessor (REQ-26, REQ-27, REQ-28):
// `flow-tag-reserved` for the reserved `recognized` key, `flow-tag-owned`
// for a key the model declares owned, and `flow-tag-duplicate` for a
// repeated name. The first two protect provenance: owned state is assembled
// from readers, never from the caller, so a `--tag` that could satisfy an
// owned dependency would let a caller assert state the artifact does not
// hold (REQ-114).
func parseTags(cmd *cobra.Command, m *table.Model) ([]resolveTag, *clierr.CLIError) {
	raw, _ := cmd.Flags().GetStringArray("tag")
	owned := flowbind.OwnedTags(m)

	var out []resolveTag
	seen := map[string]bool{}
	for _, entry := range raw {
		key, value, ok := strings.Cut(entry, "=")
		if !ok || key == "" {
			return nil, userErr(codeTagInvalid, "tag",
				"--tag takes name=value; got "+entry)
		}
		switch {
		case key == table.RecognizedTagKey:
			return nil, userErr(codeTagReserved, key,
				"the tag key `"+key+"` is reserved; the recognized outcome "+
					"enters through --outcome")
		case slices.Contains(owned, key):
			return nil, userErr(codeTagOwned, key,
				"the tag key `"+key+"` is owned; owned state is read from "+
					"the declared read accessors, never supplied by the caller")
		case seen[key]:
			return nil, userErr(codeTagDuplicate, key,
				"the tag `"+key+"` is given more than once")
		}
		seen[key] = true

		// A `--tag` legitimately carries an observed key the model does not
		// declare, and `0020:C1` decides what that absence MEANS here: a key
		// absent from the normalized tag table is a PURE CARRIER — the value
		// crosses VERBATIM, byte-preserved, JSON array literals included, and
		// is never canonicalised, conformed, kind-checked, or compared. The
		// lookup is therefore the two-value form: `declared` selects the
		// arm, and the kind/shape/domain arms below run only under a real
		// declaration, which is what makes their messages truthful.
		decl, declared := m.Tags[key]

		// The empty-value arm is HOISTED out of `canonicalValue` (`0020:C1`):
		// an empty observed value is indistinguishable from unset, a
		// grammar-level fact about the value that holds with or without a
		// declaration, so it binds the carrier too. It sits AFTER the
		// duplicate arm and the `seen[key]` mark, never before, so a
		// repeated key still refuses `flow-tag-duplicate`.
		//
		// The arm is GUARDED, not unconditional on `value == ""`: a declared
		// SET key keeps the set-specific CONFORMANCE message below, which
		// presupposes a declared kind a carrier by definition has none of.
		if value == "" && (!declared || decl.Kind != "set") {
			return nil, userErr(codeTagInvalid, key,
				"the tag `"+key+"` was given an empty value")
		}

		// The carrier branch sits where the conformance arms would have run.
		if !declared {
			out = append(out, resolveTag{Key: key, Value: value})
			continue
		}

		canonical, ce := canonicalValue(key, value, decl, "tag")
		if ce != nil {
			return nil, ce
		}
		out = append(out, resolveTag{Key: key, Value: canonical})
	}
	return out, nil
}

// resolveTag is one parsed key/value pair. It mirrors `resolve.Tag` without
// importing the kernel into the parsing layer.
type resolveTag struct {
	Key   string
	Value string
}

// canonicalValue validates a value against its declared kind and renders it
// in the form the seam carries (REQ-29, REQ-70).
//
// A set value crosses as JDR 0001 §D13's canonical JSON array — members
// sorted, duplicate-free, compact — and is RE-CANONICALISED here rather
// than passed through: the caller may legitimately hand back a plan's own
// literal, but may equally hand an unsorted or duplicated one, and read-back
// equality is byte equality over the canonical form. A scalar handed an
// array, or a set handed a bare scalar, is the wrong-kind refusal.
//
// "Well-formed for its declared kind" (REQ-61) reaches the DOMAIN, not only
// the set-vs-scalar shape: the declaration carries an enum's `domain`, an
// int's `min`/`max`, a bool's two literals, and a set's `elements`, and
// `table.ConformValue` holds the value to all of them. RDR 0002's loader
// already hard-refuses such a literal when a RULE authors it, so admitting
// it from `--write` would make the declaration advisory on exactly the path
// that persists.
//
// Set members are conformed BEFORE `canonicalSet`, so the refusal names the
// offending member as the caller spelled it rather than after sorting and
// de-duplication have moved it.
func canonicalValue(key, value string, decl table.TagDecl, flag string) (string, *clierr.CLIError) {
	looksArray := strings.HasPrefix(strings.TrimSpace(value), "[")
	isSet := decl.Kind == "set"

	if !isSet {
		if looksArray {
			return "", userErr(codeInvalidFor(flag), key,
				"the tag `"+key+"` is not set-valued; got the array literal "+value)
		}
		if value == "" {
			return "", userErr(codeInvalidFor(flag), key,
				"the tag `"+key+"` was given an empty value")
		}
		if err := table.ConformValue(decl, value); err != nil {
			return "", userErr(codeInvalidFor(flag), key,
				"the value for `"+key+"` does not conform to its declaration: "+
					err.Error())
		}
		return value, nil
	}

	if !looksArray {
		return "", userErr(codeInvalidFor(flag), key,
			"the tag `"+key+"` is set-valued and takes a JSON array literal; got "+value)
	}
	var members []string
	if err := json.Unmarshal([]byte(value), &members); err != nil {
		return "", userErr(codeInvalidFor(flag), key,
			"the set literal for `"+key+"` is not a JSON array of strings: "+err.Error())
	}
	for _, m := range members {
		if err := table.ConformValue(decl, m); err != nil {
			return "", userErr(codeInvalidFor(flag), key,
				"the set member "+strconv.Quote(m)+" for `"+key+"` does not "+
					"conform to its declaration: "+err.Error())
		}
	}
	return canonicalSet(members), nil
}

// codeInvalidFor selects the wrong-kind code for the flag the value arrived
// on. `--tag` and `--outcome` take `flow-tag-invalid`; `--write` takes
// `flow-write-invalid` (REQ-80, REQ-84).
func codeInvalidFor(flag string) string {
	if flag == "write" {
		return codeWriteInvalid
	}
	return codeTagInvalid
}

// canonicalSet renders members in JDR 0001 §D13's canonical form: sorted,
// duplicate-free, compact, HTML escaping disabled.
//
// This is THE encoder for a set literal at the CLI. Every site that emits
// or compares one calls it, which is what makes plan-to-request
// copy-through and read-back equality byte equality (REQ-71).
func canonicalSet(members []string) string {
	canonical := slices.Compact(slices.Sorted(slices.Values(members)))
	if canonical == nil {
		canonical = []string{}
	}
	var buf strings.Builder
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(canonical); err != nil {
		// Encoding a []string cannot fail.
		return "[]"
	}
	return strings.TrimRight(buf.String(), "\n")
}
