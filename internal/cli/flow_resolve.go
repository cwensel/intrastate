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

	"github.com/newcoinc/intrastate/internal/cli/clierr"
	"github.com/newcoinc/intrastate/internal/cli/respond"
	"github.com/newcoinc/intrastate/internal/guard"
	"github.com/newcoinc/intrastate/internal/resolve"
	"github.com/newcoinc/intrastate/internal/table"
	"github.com/spf13/cobra"
)

// resolvePayload is `flow resolve`'s verb-specific `data` (REQ-76).
type resolvePayload struct {
	Model    string            `json:"model"`
	Revision string            `json:"revision"`
	Observed map[string]string `json:"observed"`
	Owned    map[string]string `json:"owned"`
	Readers  []string          `json:"readers"`
	Outcome  string            `json:"outcome"`
	Rule     string            `json:"rule"`
	Gates    []gateResult      `json:"gates"`
	Next     map[string]string `json:"next"`
	Writes   map[string]string `json:"writes"`
	Clear    []string          `json:"clear"`
	Escaped  bool              `json:"escaped"`
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
	return cmd
}

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

	payload := resolvePayload{
		Model:    req.modelRef,
		Revision: req.revision(),
		Observed: observedTagMap(req.observed),
		Owned:    tagMap(owned),
		Readers:  readerIDs(readers),
		Outcome:  outcome,
		Rule:     plan.RuleID,
		Gates:    gates,
		Next:     tagMap(plan.NextTags),
		Writes:   map[string]string{},
		Clear:    []string{},
		Escaped:  plan.Escaped,
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

	return respond.OK(cmd, respond.Success{Data: payload})
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
func escapeShapeBreaches(err error) []*resolve.EscapeShapeBreachError {
	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		var out []*resolve.EscapeShapeBreachError
		for _, elem := range joined.Unwrap() {
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

// guardSeam returns RDR 0003's value-comparison evaluator, the delegated
// seam the kernel calls for typed operator semantics. This RDR does not
// own guard semantics and implements none of its own (REQ-113).
func guardSeam() resolve.GuardEvaluator { return guard.Evaluator{} }
