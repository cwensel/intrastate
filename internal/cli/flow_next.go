package cli

// RDR 0005 — `flow next`: enumerate the legal recognized outcomes and the
// candidate rows behind them.
//
// `next` ENUMERATES; it does not select. Everything else about the verb
// follows from that one distinction:
//
//   - a gate DENY is reported on its candidate and the command still exits
//     0 (REQ-43). Only `flow resolve` turns a deny into a refusal, because
//     only `resolve` was asked to pick a row;
//   - candidate summaries come from NORMALIZED model data, never from
//     evaluating a missing fact (REQ-118). A fact nothing supplied stays
//     `unresolved`, which is a report, not a decision;
//   - gates run only under `--evaluate-gates`, so the verb is effect-free
//     by default (REQ-41, REQ-57).

import (
	"slices"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
	"github.com/newcoinc/intrastate/internal/cli/respond"
	"github.com/newcoinc/intrastate/internal/resolve"
	"github.com/newcoinc/intrastate/internal/table"
	"github.com/spf13/cobra"
)

// nextPayload is `flow next`'s verb-specific `data` (REQ-74).
type nextPayload struct {
	Model     string            `json:"model"`
	Revision  string            `json:"revision"`
	Observed  map[string]string `json:"observed"`
	Owned     map[string]string `json:"owned"`
	Readers   []string          `json:"readers"`
	Outcomes  []string          `json:"outcomes"`
	Candidate []candidate       `json:"candidates"`
}

// candidate is one reported candidate row (REQ-75).
type candidate struct {
	Rule    string `json:"rule"`
	Outcome string `json:"outcome"`
	// Required names the owned facts the row's transition depends on, read
	// verbatim from the normalized row.
	Required []string `json:"required"`
	// Unresolved names the guard facts nothing established plus, when
	// gates did not run, the row's gate ids. It is what the caller sees in
	// place of a fact the CLI refuses to invent (REQ-44).
	Unresolved []string `json:"unresolved"`
	// Gates carries evaluated results, present only under
	// --evaluate-gates.
	Gates []gateResult `json:"gates,omitempty"`
	// Next, Writes, and Clear are the preview the normalized row exposes
	// without evaluating anything (REQ-40, A-5).
	Next   map[string]string `json:"next"`
	Writes map[string]string `json:"writes"`
	Clear  []string          `json:"clear"`
}

func newFlowNextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "next",
		Short: "List the legal recognized outcomes and their candidate rules",
		Long: `List the legal recognized outcomes and their candidate rules.

next ENUMERATES rather than selects. A gate that denies constrains only the
candidate carrying it: the result is reported and the command still exits
0. Turning a deny into a refusal is flow resolve's job.

Gates run only under --evaluate-gates. Without it their ids are reported as
unresolved facts, and the command invokes no gate accessor at all.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE:          runFlowNext,
	}
	registerSelectionFlags(cmd)
	registerTagFlag(cmd)
	cmd.Flags().Bool("evaluate-gates", false,
		"run the reported candidates' gate accessors and report each result")
	return cmd
}

func runFlowNext(cmd *cobra.Command, _ []string) error {
	if ce := respond.ValidateMode(cmd); ce != nil {
		return respond.Fail(cmd, ce)
	}

	req, ce := buildRequest(cmd, true)
	if ce != nil {
		return respond.Fail(cmd, ce)
	}

	// `next` narrows the reader set by the union over ALL model rows: it
	// takes no --outcome, so every row is potentially a candidate (DEV-1,
	// reading (a)).
	readers, owned, ce := req.runReaders(cmd.Context(), invokedReaders(req.model, ""))
	if ce != nil {
		return respond.Fail(cmd, ce)
	}

	view := assembledView(owned, req.observed)
	evaluate, _ := cmd.Flags().GetBool("evaluate-gates")

	payload := nextPayload{
		Model:     req.modelRef,
		Revision:  req.revision(),
		Observed:  observedTagMap(req.observed),
		Owned:     tagMap(owned),
		Readers:   readerIDs(readers),
		Candidate: []candidate{},
	}

	for _, row := range req.model.Rows {
		if len(row.Escape) != 0 {
			// An escape row is not an ordinary candidate: it participates
			// only in the rescue phase for its declared class, so
			// enumerating it here would advertise a transition the caller
			// cannot request.
			continue
		}
		if excluded(row, owned, req.observed) {
			// The row's guard is DECIDED FALSE over what is known. It is
			// not a candidate, so it is not reported — and under
			// --evaluate-gates its gates must not run (REQ-42).
			continue
		}

		c := summarize(row, view, evaluate)
		if evaluate {
			results, gce := req.runGates(cmd.Context(), row.Gate)
			if gce != nil {
				// A gate that could not be CONSULTED is exit 3 in BOTH
				// verbs — it is not a verdict, so `next` cannot report it
				// as one (REQ-46, REQ-47).
				return respond.Fail(cmd, gce)
			}
			// A deny rides the candidate and is NOT a refusal here.
			c.Gates = results
		}

		payload.Candidate = append(payload.Candidate, c)
	}

	// The recognized-outcome alphabet, as the model declares it. Each entry
	// is a bare tag — RDR 0002's layout carries no per-outcome recognizer
	// text and `[model.metadata]` is not its carrier (REQ-74).
	//
	// ASSUMPTION A-12 read REQ-40's "for the supplied state" as FILTERING
	// the alphabet down to outcomes with a surviving candidate row. The MVV
	// fixture refutes that reading: it declares `bail` with no responding
	// rule and requires `next` to report all three declared outcomes
	// (`0005:MVV`). The qualifier governs the CANDIDATES, which are what
	// varies with state; the alphabet is what the caller may legally
	// request, and a caller cannot discover that an outcome is currently
	// unreachable if the verb hides it. Recorded as DEV-4.
	payload.Outcomes = slices.Clone(req.model.Outcomes)
	if payload.Outcomes == nil {
		payload.Outcomes = []string{}
	}

	return respond.OK(cmd, respond.Success{Data: payload})
}

// summarize renders one candidate row from NORMALIZED model data alone. It
// evaluates nothing: a fact it cannot see is reported unresolved, never
// decided (REQ-44, REQ-118).
func summarize(row table.Row, view map[string]string, gatesRan bool) candidate {
	c := candidate{
		Rule:       row.RuleID,
		Outcome:    row.Outcome,
		Required:   slices.Clone(row.RequiresOwned),
		Unresolved: []string{},
		Next:       map[string]string{},
		Writes:     map[string]string{},
		Clear:      []string{},
	}
	if c.Required == nil {
		c.Required = []string{}
	}

	// A guard atom over a key nothing established is an UNRESOLVED FACT.
	for _, atom := range row.Atoms {
		if _, known := view[atom.Key]; known {
			continue
		}
		if !slices.Contains(c.Unresolved, atom.Key) {
			c.Unresolved = append(c.Unresolved, atom.Key)
		}
	}
	// An owned key no reader established is unresolved for the same reason.
	for _, key := range row.RequiresOwned {
		if _, known := view[key]; known {
			continue
		}
		if !slices.Contains(c.Unresolved, key) {
			c.Unresolved = append(c.Unresolved, key)
		}
	}
	if !gatesRan {
		// Without --evaluate-gates the gate ids ARE the unresolved facts:
		// the gate did not run, and its id is what the caller can act on
		// (REQ-41, REQ-46).
		for _, id := range row.Gate {
			if !slices.Contains(c.Unresolved, id) {
				c.Unresolved = append(c.Unresolved, id)
			}
		}
	}

	for _, t := range row.KernelRow().NextTags {
		c.Next[t.Key] = t.Value
	}
	// Clear keys are SPLIT OUT of the preview writes. RDR 0002 normalizes an
	// authored clear into a `<clear>` write, and leaving the sentinel among
	// the write targets would advertise it as a value a caller could
	// transcribe — which `--write k=<clear>` then refuses (A-5, REQ-62).
	for _, t := range row.KernelRow().Writes {
		if t.Value == table.ClearSentinel {
			c.Clear = append(c.Clear, t.Key)
			continue
		}
		c.Writes[t.Key] = t.Value
	}
	return c
}

// assembledView is the flat key/value view the CLI uses for its OWN
// reporting decisions — which facts are known, which rows are excluded.
// It is not the kernel's evaluation view: the kernel assembles its own from
// the same inputs, with provenance preserved.
func assembledView(owned []resolve.Tag, observed []resolveTag) map[string]string {
	view := make(map[string]string, len(owned)+len(observed))
	for _, t := range owned {
		view[t.Key] = t.Value
	}
	for _, t := range observed {
		view[t.Key] = t.Value
	}
	return view
}

// excluded reports whether a row's guard is DECIDED FALSE over what is
// known. Only a decided-false guard excludes: a guard the CLI cannot decide
// leaves the row a candidate with the fact reported unresolved, because
// treating "unknown" as "false" would silently drop rows the caller is
// entitled to see.
//
// The verdict is the KERNEL's, asked over a one-row probe table (REQ-113).
// The CLI decides nothing about guards itself: the row verdict rule
// `all ∧ ¬unless` (`0007:C5`/`0007:C6`) is RDR 0007's, and the typed
// operator semantics behind each atom are RDR 0003's, reached through the
// same `guardSeam()` `flow resolve` hands the kernel. That is the whole
// point of the seam — a hand-rolled filter here answered a question it does
// not own, and its answer disagreed with the kernel silently: it decided
// only `eq` under `all` over a present value, so an `unless` block or any
// other operator fell straight through, leaving the row REPORTED as a
// candidate and, under `--evaluate-gates`, its gate accessor INVOKED —
// which REQ-42 and the FM mini-check row "gate on a guard-excluded row"
// forbid.
//
// The probe strips the row's MATCH pattern and the model's other rows so
// the answer is that row's GUARD verdict alone:
//
//   - `Match` is the kernel's selection pattern, not a guard, and `next`
//     enumerates rather than selects — dropping a row whose match pattern
//     the supplied facts do not satisfy would be a selection this verb was
//     not asked to make;
//   - the escape list is stripped so the probe cannot be rescued into a
//     plan by the row's own escape modeling, which would report an
//     excluded row as viable.
//
// Only `no_match` — the kernel's answer for a row it PRUNED as decided
// false — excludes. `guard_unevaluable` and `owned_state_unavailable` are
// undecided, not false, and leave the row a reported candidate.
func excluded(row table.Row, owned []resolve.Tag, observed []resolveTag) bool {
	probe := row.KernelRow()
	probe.Match = nil
	probe.Escape = nil

	result, err := resolve.Resolve(resolve.Input{
		Table: resolve.Table{
			Outcomes: []string{row.Outcome},
			Rows:     []resolve.Row{probe},
		},
		Owned:      owned,
		Observed:   kernelTags(observed),
		Recognized: row.Outcome,
		Guards:     guardSeam(),
	})
	if err != nil || !result.Refused() {
		return false
	}
	return result.Refusal.Kind == resolve.KindNoMatch
}

func readerIDs(outputs []readerOutput) []string {
	out := make([]string, 0, len(outputs))
	for _, r := range outputs {
		out = append(out, r.ID)
	}
	return out
}

// LintFindings lets a refusal payload satisfy respond.FindingCarrier. The
// success payloads here carry none, so this file declares no carrier.
var _ = clierr.Finding{}
