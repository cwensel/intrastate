package cli

// RDR 0005 — `flow resolve`: map one recognized outcome to exactly one
// plan, or to exactly one typed refusal.
//
// `resolve` SELECTS, which is what separates it from `next`. The order it
// works in is fixed by JDR 0001 §D9 and is the whole safety argument:
//
//	parse input -> load model -> run the narrowed readers -> ONE kernel
//	call -> run the SELECTED row's gates -> render one terminal result
//
// Gates run after exact-one selection and before the plan is emitted. Any
// earlier and the CLI would gate rows the model never chose; any later and
// it would emit a plan a gate denies. A deny is its own refusal — never a
// no-match, never an escape class — so a denied transition can never be
// laundered into a rescue the model did not model (REQ-116, REQ-119).

import (
	"errors"
	"slices"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/cli/respond"
	"github.com/cwensel/intrastate/internal/guard"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
	"github.com/spf13/cobra"
)

// resolvePayload is `flow resolve`'s verb-specific `data` (REQ-76).
//
// RDR 0023 `0023:A2` retypes the five ECHO-group fields — `Model`,
// `Observed`, `Owned`, `Readers`, `Outcome` — from `T` to `*T` with
// `,omitempty`, leaving `T` itself and every JSON key string unchanged.
// That is the mechanism `--plan-only`'s projection deletes a key with: a
// nilled pointer under `omitempty` drops the key entirely, which is what
// "an omitted key is ABSENT, never null and never an empty placeholder"
// requires. A bare non-pointer `omitempty` cannot serve, because it would
// also drop an EMPTY `{}`/`[]` in DEFAULT mode and silently move bytes for
// callers who never opted in (`0023:A9`).
//
// The conversion changes TYPES only: the field COUNT stays 14 and the
// declaration ORDER — which is the JSON key order Go's encoder emits — is
// untouched, so `decision_table_0010_test.go`'s pinned count, key list and
// `Emit`-after-`Gates` adjacency all pass verbatim (`0023:A3`).
//
// The five pointers are non-nil on every emitted SUCCESS payload; that is
// now an implementation invariant rather than an incidental property
// (`0023:A9`). `runFlowResolve` below allocates each one unconditionally.
type resolvePayload struct {
	Model    *string            `json:"model,omitempty"`
	Revision string             `json:"revision"`
	Observed *map[string]string `json:"observed,omitempty"`
	Owned    *map[string]string `json:"owned,omitempty"`
	Readers  *[]string          `json:"readers,omitempty"`
	Outcome  *string            `json:"outcome,omitempty"`
	Rule     string             `json:"rule"`
	Gates    []gateResult       `json:"gates"`
	// Emit is the selected row's `[rule.emit]` block (`0010:C4`): a JSON
	// object of string values, keys in byte order, present as `{}` — never
	// `null`, never omitted — when the selected row authored none.
	//
	// Its POSITION is fixed immediately after Gates, because the
	// declaration order IS the JSON key order Go's encoder emits and
	// displacing Next..EscapeClass would rewrite a key order this RDR does
	// not own. Gates, not Rule, is the predecessor.
	Emit    map[string]string `json:"emit"`
	Next    map[string]string `json:"next"`
	Writes  map[string]string `json:"writes"`
	Clear   []string          `json:"clear"`
	Escaped bool              `json:"escaped"`
	// EscapeClass is carried only when the plan came from an escape row.
	EscapeClass string `json:"escape_class,omitempty"`
}

func newFlowResolveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "Map one recognized outcome to exactly one plan",
		Long: `Map one recognized outcome to exactly one plan, or refuse.

resolve returns exactly one plan or exactly one structured refusal. It
never chooses among several matching rules: two matches is
flow-ambiguous-match, because picking one would be a tie-break the model
did not author.

Gates on the selected rule run after selection and before the plan is
emitted. A deny is flow-gate-denied — a refusal, never a plan and never an
escape class.

The plan is data the caller acts on. resolve applies nothing: writing the
plan back is flow set-state's job, and nothing links the two calls.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE:          runFlowResolve,
	}
	registerSelectionFlags(cmd)
	registerTagFlag(cmd)
	cmd.Flags().String("outcome", "", "the recognized outcome tag to resolve")
	// RDR 0023 `0023:C1` — VERB-LOCAL, deliberately. The projection axis
	// rides `resolve` alone: it never enters `registerSelectionFlags` (the
	// shared registrar), the `flow` group's own or persistent set, or the
	// root's persistent set. That non-registration IS the fence — the
	// sibling verbs refuse `--plan-only` at the parse, through the shared
	// `command-error` bucket, and this contract mints no code of its own.
	cmd.Flags().Bool(planOnlyFlagName, false, planOnlyFlagUsage)
	withExtendedHelp(cmd, flowResolveExtendedDesc)
	return cmd
}

var flowResolveExtendedDesc = `Reading a successful plan

  rule          the id of the one row that was selected.
  next{}        the owned state that row's writes call for.
  writes{}      the planned writes, as name=value.
  clear[]       the owned keys the plan removes. The <clear> sentinel
                never appears as a write VALUE, so transcribe these with
                --clear <key> on set-state, never --write key=<clear>.
  emit{}        a decision table's answer block for the selected row.
                Its keys are not tags: nothing declares them, nothing
                writes them, and they take no part in selection.
  escaped       true when the plan came from an escape row.
  escape_class  the failure class that escape row RESCUED — the request's
                own refusal kind, not the whole list the row declares. A
                row modeling several classes rescued this request under
                exactly one.
  gates[]       every gate on the selected row, with its result.

  The plan is data you act on. resolve applies nothing: writing it back
  is flow set-state's job, and nothing links the two calls.

  The success payload has two halves. The PLAN half above is what this
  call DECIDED. The other half is the request read back to you:

  model         the model reference you passed.
  observed      the --tag values you passed, echoed unchanged.
  owned         the owned state the readers assembled.
  readers       the read accessors that were invoked.
  outcome       the --outcome you passed, echoed unchanged.

  --plan-only omits that second half. Nothing else changes: the same
  rule is selected, the same gates run, refusals are byte-identical, and
  every field that is carried is carried byte-for-byte. Pass it when you
  already hold the request and only want the decision — and do not pass
  it if you were relying on the payload to tell you which model answered,
  since model is exactly what it drops.

Refusals

  ` + codeTagInvalid + `
      --outcome was absent or empty. A malformed request.
  ` + codeUnmodeledOutcome + `
      the outcome is outside the model's declared alphabet. A request
      the model does not recognise — a different mistake, so a different
      code. The offending outcome is in param.
  ` + codeNoMatch + `
      no rule matches over the assembled state. Every rule considered is
      named in findings.
  ` + codeAmbiguousMatch + `
      more than one rule matches. EVERY conflicting row is named in
      findings: reporting one would be exactly the tie-break the model
      did not author. resolve never picks among them.
  ` + codeOwnedStateUnavailable + `
      a candidate rule requires owned state no reader established. One
      finding per missing key.
  ` + codeGateDenied + `
      a gate on the selected row denied. A refusal — never a plan, and
      never an escape class: a denied transition cannot be rescued into
      one.
  ` + codeGuardUnevaluable + `
      a guard on a matching rule could not be evaluated over the
      assembled state.
  ` + codeGateIndeterminate + `
      a gate answered neither allow nor deny.

  Gates on the selected row all run; every result is reported; deny
  overrides both allow and indeterminate.

Shared refusals

  The codes above are the ones specific to this verb. Model selection,
  tag and artifact validation, and the accessor and environment failures
  are common to every flow verb and are listed once under
  "intrastate flow --help-all" rather than repeated here.

Exits

  0  exactly one plan.
  2  any refusal above — the request or the model is wrong, or the model
     said no.
  3  a reader or gate could not be consulted; re-run unchanged once the
     environment is repaired.

Worked call

  intrastate flow resolve --model flow.toml \
      --artifact state=state.json \
      --tag actor=reviewer --outcome approved --as json`

func runFlowResolve(cmd *cobra.Command, _ []string) error {
	if ce := respond.ValidateMode(cmd); ce != nil {
		return respond.Fail(cmd, ce)
	}

	req, ce := buildRequest(cmd, true)
	if ce != nil {
		return respond.Fail(cmd, ce)
	}

	// An absent or empty `--outcome` is a CLI-side refusal; a value outside
	// the model's alphabet is the KERNEL's unmodeled_outcome. The two carry
	// different codes because they are different mistakes: one is a
	// malformed request, the other a request the model does not recognise
	// (REQ-54).
	outcome, _ := cmd.Flags().GetString("outcome")
	if outcome == "" {
		return respond.Fail(cmd, userErr(codeTagInvalid, "outcome",
			"flow resolve requires --outcome <tag>"))
	}

	readers, owned, ce := req.runReaders(cmd.Context(),
		invokedReaders(req.model, outcome))
	if ce != nil {
		return respond.Fail(cmd, ce)
	}

	// ONE kernel call. The kernel owns selection, guard evaluation, and the
	// escape phase; the CLI neither pre-filters candidates nor re-decides
	// anything it returns (REQ-111, REQ-113).
	result, err := resolve.Resolve(resolve.Input{
		Flow:       req.model.ID,
		Table:      req.model.KernelTable(),
		Owned:      owned,
		Observed:   kernelTags(req.observed),
		Recognized: outcome,
		Guards:     guardSeam(),
	})
	if err != nil {
		return respond.Fail(cmd, kernelResolveFailure(err))
	}
	if result.Refused() {
		return respond.Fail(cmd, kernelFailure(*result.Refusal))
	}

	plan := result.Plan
	row, found := rowByID(req.model, plan.RuleID)
	if !found {
		return respond.Fail(cmd, internalErr(codeAccessorFailed,
			"the kernel selected a rule the loaded model does not carry: "+
				plan.RuleID))
	}

	// --- the gate site (JDR 0001 §D9) -----------------------------------
	// All gates on the SELECTED row run; every result is reported; deny
	// overrides allow and indeterminate.
	gates, ce := req.runGates(cmd.Context(), row.Gate)
	if ce != nil {
		return respond.Fail(cmd, ce)
	}
	if ce := gateVerdictFailure(gates); ce != nil {
		return respond.Fail(cmd, ce)
	}

	// Payload assembly is FLAG-BLIND (`0023:C2`): it never reads
	// `--plan-only` and never takes it as a parameter. Every echo field is
	// assembled and every pointer allocated here, whatever the caller asked
	// for — projection is the only flag-aware step, and it runs on the
	// FULLY ASSEMBLED value below.
	//
	// The five echo pointers are allocated unconditionally, which is what
	// keeps default-mode bytes identical to the pre-change binary: an empty
	// container still renders `{}`/`[]` rather than `null` or nothing
	// (`0023:A9`).
	observedView := observedTagMap(req.observed)
	ownedView := tagMap(owned)
	readerView := readerIDs(readers)
	payload := resolvePayload{
		Model:    &req.modelRef,
		Revision: req.revision(),
		Observed: &observedView,
		Owned:    &ownedView,
		Readers:  &readerView,
		Outcome:  &outcome,
		Rule:     plan.RuleID,
		Gates:    gates,
		// The join is by rule id, AFTER selection and after the gates: `row`
		// is what `rowByID` already returned for `plan.RuleID`, so an
		// escaped plan carries the ESCAPE ROW's own block with no
		// escape-specific arm, and a row authoring none renders `{}`.
		// `emit` never crosses the kernel seam, and a gate deny is a
		// refusal rather than a payload, so this is never computed on a
		// denied selection.
		Emit:    emitMap(row.Emit),
		Next:    tagMap(plan.NextTags),
		Writes:  map[string]string{},
		Clear:   []string{},
		Escaped: plan.Escaped,
	}
	if payload.Gates == nil {
		payload.Gates = []gateResult{}
	}
	if payload.Next == nil {
		payload.Next = map[string]string{}
	}

	// Planned clears ride `clear[]`; the `<clear>` sentinel never appears
	// as a write VALUE, so a caller transcribing the plan reaches for
	// `--clear` rather than for a literal `--write` refuses (REQ-62,
	// REQ-76, REQ-108).
	for _, t := range plan.Writes {
		if t.Value == table.ClearSentinel {
			payload.Clear = append(payload.Clear, t.Key)
			continue
		}
		payload.Writes[t.Key] = t.Value
	}

	if plan.Escaped {
		// The class reported is the one that was RESCUED — the request's
		// own refusal kind — not the whole `Escape` list a row may declare
		// (A-3). A row modeling several classes rescued this request under
		// exactly one of them, and naming the others would misreport why
		// the plan exists.
		payload.EscapeClass = escapeClassOf(row, req, owned, outcome)
	}

	// RDR 0023 `0023:C1` — the projection site, and the only place this verb
	// reads `--plan-only`. It is AFTER the last `respond.Fail` return above,
	// so every refusal path is flag-blind structurally rather than by a
	// guard, and it is BEFORE the gateway, so both output modes render one
	// projected result through the one renderer.
	return respond.OK(cmd, respond.Success{
		Data: projectResolvePayload(cmd, payload),
	})
}

// emitMap renders the selected row's emit sequence as the payload's
// string→string object (`0010:C4`).
//
// It allocates UNCONDITIONALLY, so a row authoring no block yields the
// empty map rather than nil — the difference between `"emit":{}` and
// `"emit":null` on the wire, and C4 fixes the former. Key ORDER is not this
// function's to fix: the payload encoder emits a map's keys in byte order
// with HTML escaping disabled, which is the one policy every other payload
// map already crosses on, and the row's sequence is key-sorted anyway.
func emitMap(emit []table.EmitValue) map[string]string {
	out := make(map[string]string, len(emit))
	for _, e := range emit {
		out[e.Key] = e.Value
	}
	return out
}

// rowByID finds the normalized row the kernel selected.
func rowByID(m *table.Model, ruleID string) (table.Row, bool) {
	for _, row := range m.Rows {
		if row.RuleID == ruleID {
			return row, true
		}
	}
	return table.Row{}, false
}

// escapeClassOf recovers the failure class an escape row rescued.
//
// `resolve.Plan` carries `Escaped bool` but no class, so the CLI recovers
// it by re-asking the kernel what the request would have refused WITHOUT
// the escape rows present. That answer is the request's own refusal kind —
// precisely the class the escape row rescued — and it needs no kernel
// change (A-3).
//
// A row declaring several classes therefore reports the one that fired, not
// the list. If the probe cannot name a kind the row actually rescues, the
// row's single declared class is the answer; a row declaring several and
// rescuing an unprobeable kind reports nothing rather than guessing.
func escapeClassOf(row table.Row, req flowRequest, owned []resolve.Tag, outcome string) string {
	probe := req.model.KernelTable()
	var ordinary []resolve.Row
	for _, r := range probe.Rows {
		if len(r.Escape) == 0 {
			ordinary = append(ordinary, r)
		}
	}
	probe.Rows = ordinary

	result, err := resolve.Resolve(resolve.Input{
		Flow:       req.model.ID,
		Table:      probe,
		Owned:      owned,
		Observed:   kernelTags(req.observed),
		Recognized: outcome,
		Guards:     guardSeam(),
	})
	if err == nil && result.Refused() {
		kind := string(result.Refusal.Kind)
		if slices.Contains(row.Escape, kind) {
			return kind
		}
	}
	if len(row.Escape) == 1 {
		return row.Escape[0]
	}
	return ""
}

// LintFindings keeps the findings type referenced from this file's imports.
var _ = clierr.Finding{}

// --- kernel Go-error mapping (RDR 0009 `0009:C7`) ------------------------

// kernelResolveFailure wraps one kernel Go error — the channel RDR 0001
// reserves for producer/programmer mistakes, never for modeled refusals —
// into the CLI envelope. It is the sibling of kernelFailure, which maps
// modeled refusals.
//
// It DISCRIMINATES: an escape-row shape breach gets its own stable Code and
// carries every offending row identity structurally; every other kernel
// error still falls through to the generic accessor-failed branch. A
// blanket recode would mislabel RDR 0008's reserved-key breach.
func kernelResolveFailure(err error) *clierr.CLIError {
	if err == nil {
		return nil
	}
	if !errors.Is(err, resolve.ErrEscapeShapeBreach) {
		return internalErr(codeAccessorFailed,
			"the resolution kernel reported a programmer error: "+err.Error())
	}

	// The kernel's aggregate is exactly one level deep, every element a
	// *EscapeShapeBreachError, already sorted by RowRef identity with equal
	// identities collapsed (`0009:C5`). The verb traverses it rather than
	// calling errors.As once, which would report the first breach alone.
	ce := &clierr.CLIError{
		Code:  codeEscapeRowShapeBreach,
		Group: clierr.GroupInternal,
		Message: "the transition table carries an escape row that also " +
			"carries writes; an escape row describes no owned-state mutation",
		Hint:  hintEscapeRowShapeBreach,
		Cause: err,
	}

	// The identities reach the wire on `findings` — `Cause` is json:"-", so
	// the Go chain that holds the RowRef values is not wire-visible, and
	// folding them into `Detail` would make a caller re-parse prose.
	// clierr stays a leaf: the conversion from resolve.RowRef to plain
	// strings happens here, in the layer that already imports both.
	for _, breach := range escapeShapeBreaches(err) {
		ce.Findings = append(ce.Findings, clierr.Finding{
			Code:    codeEscapeRowShapeBreach,
			Rule:    breach.Ref.RuleID,
			Locator: breach.Ref.SourceLocator,
			Count:   breach.Count,
			Message: "this escape row carries writes",
			Hint:    hintEscapeRowShapeBreach,
		})
	}
	return ce
}

// escapeShapeBreaches returns the per-identity breach elements err carries,
// in the kernel's order. The traversal is one level deep, matching the flat
// shape `0009:C4` guarantees; a bare (unjoined) breach is handled too, so
// the verb does not depend on the aggregate's cardinality.
//
// Reaching the aggregate walks the single-unwrap chain the way errors.As
// does, rather than asserting on the outermost error alone: REQ-36 forbids
// the kernel from wrapping its own join, but a bare assertion would silently
// reduce an N-breach report to its first breach under any wrap a caller adds.
func escapeShapeBreaches(err error) []*resolve.EscapeShapeBreachError {
	if joined := joinedBreaches(err); joined != nil {
		var out []*resolve.EscapeShapeBreachError
		for _, elem := range joined {
			var breach *resolve.EscapeShapeBreachError
			if errors.As(elem, &breach) {
				out = append(out, breach)
			}
		}
		return out
	}
	var breach *resolve.EscapeShapeBreachError
	if errors.As(err, &breach) {
		return []*resolve.EscapeShapeBreachError{breach}
	}
	return nil
}

// joinedBreaches returns the elements of the first `Unwrap() []error`
// aggregate on err's single-unwrap chain, or nil if the chain holds none.
// The chain walk mirrors errors.As; the aggregate itself is not descended
// into, keeping the element traversal exactly one level deep (REQ-35).
func joinedBreaches(err error) []error {
	for err != nil {
		if joined, ok := err.(interface{ Unwrap() []error }); ok {
			return joined.Unwrap()
		}
		unwrapped, ok := err.(interface{ Unwrap() error })
		if !ok {
			return nil
		}
		err = unwrapped.Unwrap()
	}
	return nil
}

// guardSeam returns RDR 0003's value-comparison evaluator, the delegated
// seam the kernel calls for typed operator semantics. This RDR does not
// own guard semantics and implements none of its own (REQ-113).
func guardSeam() resolve.GuardEvaluator { return guard.Evaluator{} }
