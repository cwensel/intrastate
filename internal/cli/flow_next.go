package cli

// RDR 0005 / RDR 0011 — `flow next`: report the candidate rows the supplied
// state can take, and the legal recognized outcomes behind them.
//
// RDR 0011 overrides 0005's candidate clause. By DEFAULT a row's match
// atoms take part in the verdict: a candidate is a non-escape row whose
// match atoms over keys PRESENT in the assembled view all hold and whose
// guard the kernel does not decide false. `--all` restores 0005's
// predicate — every non-escape row the guards do not exclude, with the
// match pattern taking no part.
//
// The placement is load-bearing: the CLI decides PRESENCE, the kernel
// decides EQUALITY. `flow next` never compares a match value; it restricts
// the probe's match pattern to keys the view holds and asks the kernel.
//
// `next` still REPORTS rather than selects. Everything else about the verb
// follows from that one distinction:
//
//   - a gate DENY is reported on its candidate and the command still exits
//     0 (REQ-43). Only `flow resolve` turns a deny into a refusal, because
//     only `resolve` was asked to pick a row;
//   - candidate summaries come from NORMALIZED model data, never from
//     evaluating a missing fact (REQ-118). A fact nothing supplied is
//     reported under `unknown` with its reason named, which is a report,
//     not a decision;
//   - gates run only under `--evaluate-gates`, so the verb is effect-free
//     by default (REQ-41, REQ-57).

import (
	"slices"
	"strings"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/cli/respond"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
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
	// Unknown names every fact the run could not decide for this row, each
	// as a `{key, reason}` pair: a match atom over an absent key, a guard
	// atom the kernel could not decide, an owned key no invoked reader
	// established, and — absent --evaluate-gates — the row's own gate ids.
	// It is what the caller sees in place of a fact the CLI refuses to
	// invent (REQ-44), and the reason names the remedy (`0011:C1`).
	//
	// Present in BOTH modes, as `[]` rather than omitted, so one consumer
	// struct parses either (`0011:C2`).
	Unknown []unknownFact `json:"unknown"`
	// Gates carries evaluated results, present only under
	// --evaluate-gates.
	Gates []gateResult `json:"gates,omitempty"`
	// Next, Writes, and Clear are the preview the normalized row exposes
	// without evaluating anything (REQ-40, A-5).
	Next   map[string]string `json:"next"`
	Writes map[string]string `json:"writes"`
	Clear  []string          `json:"clear"`
}

// unknownFact is one entry of a candidate's `unknown` list: the key (or
// gate id) the run could not decide, and why (`0011:D-undecided-reporting-shape`).
type unknownFact struct {
	Key    string `json:"key"`
	Reason string `json:"reason"`
}

// The closed reason vocabulary. `absent` and `uncomparable` are
// `0007:C8`'s kernel-owned set, reused here rather than minted;
// `not-evaluated` is the ONE token RDR 0011 adds, scoped to un-run gate ids
// on this CLI list and never appearing on a match, guard, or owned-key
// entry (`0011:D-undecided-vocabulary`).
const (
	reasonAbsent       = string(resolve.ReasonAbsent)
	reasonUncomparable = string(resolve.ReasonUncomparable)
	reasonNotEvaluated = "not-evaluated"
)

// UnknownReasons returns the `flow next` unknown-`reason` vocabulary
// (`0029:C4`). The set is append-only: a consumer tolerates an
// unrecognized member and does not assert on the set's cardinality or a
// member's ordinal position.
//
// The union is the two kernel-owned reasons plus RDR 0011's one added
// token, spelled here as the members they are: `absent`, `uncomparable`,
// `not-evaluated`.
func UnknownReasons() []string {
	return []string{"absent", "uncomparable", "not-evaluated"}
}

func newFlowNextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "next",
		Short: "Report the candidate rules the supplied state can take",
		Long: `Report the candidate rules the supplied state can take.

A candidate is a rule whose match and guard both HOLD or are UNDECIDED over
the state you supplied. A match or guard key the state does not carry does
not exclude the row: it leaves the row a candidate with the key listed
under unknown and its reason named, so you can see what to supply.

A candidate is a row the supplied state does not exclude — it is NOT a row
flow resolve will select. A candidate carrying an unknown entry may still
be refused by flow resolve over the same state, and the entry names what to
supply to settle it. Choosing exactly one row is flow resolve's job.

--all reports every row the guards do not exclude, regardless of match. In
that mode a match atom takes no part in the verdict and contributes no
unknown entry.

A gate that denies constrains only the candidate carrying it: the result is
reported and the command still exits 0. Turning a deny into a refusal is
flow resolve's job.

Gates run only under --evaluate-gates. Without it their ids are reported
under unknown as not-evaluated, and the command invokes no gate accessor at
all.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE:          runFlowNext,
	}
	registerSelectionFlags(cmd)
	registerTagFlag(cmd)
	cmd.Flags().Bool("evaluate-gates", false,
		"run the reported candidates' gate accessors and report each result")
	// `--all` is registered on `next` ONLY — never on the shared
	// registerSelectionFlags, which would place it on resolve, read-state,
	// and set-state too (`0011:C2`). No shorthand: `-a` widens the surface
	// C2 pins (`0011:A14`).
	cmd.Flags().Bool("all", false,
		"report every row the guards do not exclude, regardless of match")
	withExtendedHelp(cmd, flowNextExtendedDesc)
	return cmd
}

const flowNextExtendedDesc = `Reading the output

  outcomes[]   the model's declared recognized alphabet — every outcome
               that can be asked of it, not only those with a surviving
               candidate row.
  candidates[] one entry per row the supplied state does not exclude,
               each naming its rule id, its outcome, the owned keys its
               transition requires, and — under --evaluate-gates — its
               gate results.
  unknown[]    per candidate, the keys that left the row undecided and
               the reason each is unsettled. This is the actionable
               field: it names exactly what to supply.
  observed{}   the tags you passed, echoed as parsed.
  owned{}      the owned tags the invoked readers established.
  readers[]    which readers this call actually invoked — the model's
               demand set for these rows, not every declared reader.

  A candidate carrying an unknown entry may still be refused by flow
  resolve over the same state. next reports what is not excluded;
  choosing among what remains is resolve's job, and the two answering
  differently is the contract working, not a discrepancy.

--all

  Without it, a row is a candidate when its match and its guards both
  hold or are undecided. With it, matches take no part in the verdict:
  every row the GUARDS do not exclude is reported, and a match atom
  contributes no unknown entry. Use it to see the whole guard-legal
  surface at a state, independent of any one outcome.

Gates

  Gates are not evaluated unless you pass --evaluate-gates. Without it,
  no gate accessor is invoked at all and each gate id is reported under
  unknown as not-evaluated — so a bare next never touches the
  environment on a gate's behalf.

  With it, a gate that DENIES constrains only the candidate carrying it.
  The result is reported and the command still exits 0. Turning a deny
  into the refusal ` + codeGateDenied + ` is flow resolve's job: next
  reports, it does not adjudicate.

Shared refusals

  The codes above are the ones specific to this verb. Model selection,
  tag and artifact validation, and the accessor and environment failures
  are common to every flow verb and are listed once under
  "intrastate flow --help-all" rather than repeated here.

Exits

  0  candidates reported, including none, and including denied gates.
  2  the request or the model is wrong.
  3  a reader or an evaluated gate could not be consulted.

Worked call

  intrastate flow next --model flow.toml \
      --artifact state=state.json \
      --tag actor=reviewer --evaluate-gates --as json`

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

	// The invoked reader set and the assembled view are fixed ONCE per
	// invocation, before the row loop, and are the same under --all: the
	// demand set is a property of the MODEL, not of the mode, so `unknown`
	// is never a function of row order or of the flag (`0011:C1`).
	view := assembledView(owned, req.observed)
	evaluate, _ := cmd.Flags().GetBool("evaluate-gates")
	all, _ := cmd.Flags().GetBool("all")

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
		result := probeRow(row, view, owned, req.observed, all)
		if excluded(result) {
			// The probe answered `no_match`: the row's guard is DECIDED
			// FALSE, or (by default) a match atom over a PRESENT key is
			// unequal. Either way it is not a candidate, so it is not
			// reported — and under --evaluate-gates its gates must not run
			// (REQ-42, `0011:C1`).
			continue
		}

		c := summarize(row, view, result, evaluate, all)
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

// summarize renders one candidate row from NORMALIZED model data plus the
// row's own probe Result. It compares no value of its own: every fact it
// reports is either a PRESENCE verdict off its own view walk or a reason
// the KERNEL put on the probe's refusal payload (REQ-44, REQ-118,
// `0011:C1`).
//
// The `unknown` list has four sources and three emission sites:
//
//   - the ATOM walk — a match atom over an absent key and a guard atom over
//     an absent key, both `absent`, and (under --all) match atoms filtered
//     out AT EMISSION so the mode reports nothing about a predicate it does
//     not apply (`0011:C2`);
//   - the OWNED-key walk — an owned key no invoked reader established,
//     `absent`. It is INDEPENDENT of the atom walk and the --all filter
//     never touches it;
//   - the GATE-id loop — absent --evaluate-gates, the row's own declared
//     gates, `not-evaluated`, emitted on every reported disposition;
//
// plus the kernel's `guard_unevaluable` payload, which is the only source
// of `uncomparable` and the only one the CLI cannot derive itself.
//
// Dedup on the `{key, reason}` PAIR and the sort are the LAST step, after
// every source has contributed and after the --all filter has run at
// emission (`0011:C1`).
func summarize(
	row table.Row,
	view map[string]string,
	result resolve.Result,
	gatesRan bool,
	all bool,
) candidate {
	c := candidate{
		Rule:     row.RuleID,
		Outcome:  row.Outcome,
		Required: slices.Clone(row.RequiresOwned),
		Unknown:  []unknownFact{},
		Next:     map[string]string{},
		Writes:   map[string]string{},
		Clear:    []string{},
	}
	if c.Required == nil {
		c.Required = []string{}
	}

	var facts []unknownFact

	// (1) the ATOM walk. A match or guard atom over a key nothing
	// established is `absent` — a presence verdict, never a comparison.
	for _, atom := range row.Atoms {
		if all && atom.Block == table.BlockMatch {
			// Under --all the match pattern takes no part in the verdict,
			// so reporting a match fact would be reporting on a predicate
			// this mode does not apply. The filter is per-ATOM and applied
			// HERE, at emission: a GUARD atom on the same key still
			// contributes, and so does the owned-key walk below
			// (`0011:C2`).
			continue
		}
		if _, known := view[atom.Key]; known {
			continue
		}
		facts = append(facts, unknownFact{Key: atom.Key, Reason: reasonAbsent})
	}

	// (2) the OWNED-key walk, independent of the atom walk and of the mode.
	for _, key := range row.RequiresOwned {
		if _, known := view[key]; known {
			continue
		}
		facts = append(facts, unknownFact{Key: key, Reason: reasonAbsent})
	}

	// (3) the kernel's own reasons for the atoms it could not decide. The
	// read is GUARDED on the probe's disposition being `guard_unevaluable`
	// — the only refusal that carries `Refusal.Undecided` — rather than
	// read off whatever Result is in hand and filtered afterwards
	// (`0011:C1`, `0011:A13`).
	facts = append(facts, undecidedFacts(result)...)

	if !gatesRan {
		// (4) Without --evaluate-gates the gate ids are what the caller can
		// act on, and the reason is neither `absent` (the gate is declared)
		// nor `uncomparable` (a gate id is not a view key with a value):
		// it is `not-evaluated`, minted for this class only. Emitted on
		// EVERY reported disposition (REQ-41, REQ-46, `0011:C1`).
		for _, id := range row.Gate {
			facts = append(facts, unknownFact{Key: id, Reason: reasonNotEvaluated})
		}
	}

	c.Unknown = dedupeUnknown(facts)

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

// probeRow asks the KERNEL for the row's verdict, over a one-row probe
// table built from that row alone.
//
// The CLI decides nothing about guards itself: the row verdict rule
// `all ∧ ¬unless` (`0007:C5`/`0007:C6`) is RDR 0007's, and the typed
// operator semantics behind each atom are RDR 0003's, reached through the
// same `guardSeam()` `flow resolve` hands the kernel. That is the whole
// point of the seam — a hand-rolled filter here answered a question it does
// not own, and its answer disagreed with the kernel silently.
//
// The CLI decides EQUALITY nowhere. What it decides is PRESENCE, by the
// same view-presence test `summarize` applies — and it applies that as a
// FILTER over the CONVERTED probe (`row.KernelRow()`, then drop every
// `resolve.Tag` whose `Key` the view lacks), never by rebuilding match tags
// from `row.Atoms`: `KernelRow` renders each literal through `seamValue`'s
// canonicalization, and a second encoder here would fork the seam
// (`0011:C1`, `0011:BR4`).
//
// Filtering after conversion is sound because `resolve.Tag` carries `Key`
// and presence is a property of the KEY: every match tag on one key shares
// one presence verdict, so the filter hands a key's tags to the kernel
// together or omits them together and can never split a key. Whether a
// key's tags hold TOGETHER is then the kernel's conjunction, which is what
// makes a `0002:C13` dead row a candidate while its key is absent and
// excluded once the key is present at any value — with no CLI literal
// comparison anywhere (`0011:C1`, `0011:A12`).
//
// Under `--all` the match pattern is OMITTED entirely rather than
// presence-filtered, which is what makes the mode's predicate exactly the
// one `0005:C1` specified (`0011:C2`, `0011:A6`).
//
// Two probe properties are unchanged from 0005:
//
//   - `Recognized` is bound to the row's OWN outcome, alongside
//     `Outcomes: []string{row.Outcome}` so `Table.models` holds trivially.
//     The binding keeps the probe modelling its own outcome, so
//     `unmodeled_outcome` is unreachable (`0011:A8`);
//   - the escape list is stripped so the probe cannot be rescued into a
//     plan by the row's own escape modeling, which would report an excluded
//     row as viable (`0011:A7`).
func probeRow(
	row table.Row,
	view map[string]string,
	owned []resolve.Tag,
	observed []resolveTag,
	all bool,
) resolve.Result {
	probe := row.KernelRow()
	probe.Escape = nil
	if all {
		probe.Match = nil
	} else {
		probe.Match = presentMatchTags(probe.Match, view)
	}

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
	if err != nil {
		// A probe the kernel could not even accept is not a verdict, so it
		// cannot exclude: the row stays a candidate.
		return resolve.Result{}
	}
	return result
}

// presentMatchTags restricts a converted probe's match pattern to the keys
// the assembled view HOLDS A VALUE for. Presence is map membership, so an
// empty-but-present value is PRESENT — the same test the kernel's own
// assembled view applies (`0011:A11`).
//
// A match atom whose key is absent is omitted and does not exclude; it is
// reported `{key, absent}` instead. A row ALL of whose match atoms are
// omitted is probed with an empty match pattern, which matches
// unconditionally (`0011:C1`).
func presentMatchTags(match []resolve.Tag, view map[string]string) []resolve.Tag {
	var kept []resolve.Tag
	for _, t := range match {
		if _, known := view[t.Key]; !known {
			continue
		}
		kept = append(kept, t)
	}
	return kept
}

// excluded reports whether the probe's disposition EXCLUDES the row.
//
// `no_match` is the only one that does — the kernel's answer for a row it
// PRUNED as decided false, whether by a guard or by a present-and-unequal
// match atom. `guard_unevaluable` and `owned_state_unavailable` are
// UNDECIDED, not false: they leave the row a reported candidate with their
// facts named, because treating "unknown" as "false" would silently drop
// rows the caller is entitled to see (`0011:C1`).
func excluded(result resolve.Result) bool {
	if !result.Refused() {
		return false
	}
	return result.Refusal.Kind == resolve.KindNoMatch
}

// undecidedFacts reads the reasons the KERNEL reported for atoms it could
// not decide, off the probe's `guard_unevaluable` refusal payload.
//
// The read is GUARDED on that disposition: `guard_unevaluable` is the only
// refusal that carries `Refusal.Undecided`, so reading the field off
// whatever Result is in hand and filtering afterwards would invent facts
// for a row that is not on the wire (`0011:C1`, `0011:A13`).
//
// This is the ONLY source of `uncomparable` — a key that is PRESENT at a
// value the operator cannot parse, which the CLI's own view walk cannot
// see. It also re-reports `absent` for a guard atom over an absent key; the
// walk produces the same pair and dedup on the PAIR collapses them to one
// entry.
//
// One fidelity limit is stated rather than hidden: when the probe refuses
// `owned_state_unavailable` the kernel returns before it collects the guard
// payload, so an `uncomparable` guard atom on that row is not on the wire
// and is not reported. Every `absent` fact on the row still is, from the
// walk (`0011:C1`, `0011:S7`).
func undecidedFacts(result resolve.Result) []unknownFact {
	if !result.Refused() || result.Refusal.Kind != resolve.KindGuardUnevaluable {
		return nil
	}
	var out []unknownFact
	for _, undecided := range result.Refusal.Undecided {
		for _, atom := range undecided.Atoms {
			out = append(out, unknownFact{
				Key:    atom.Key,
				Reason: string(atom.Reason),
			})
		}
	}
	return out
}

// dedupeUnknown is the LAST step: every source has contributed and C2's
// --all filter has already run at emission, so the merged entries are now
// deduplicated by the `{key, reason}` PAIR — not by key, since `absent` and
// `uncomparable` name different remedies — and sorted by `(key, reason)`,
// so a consumer diffing a golden payload sees no churn (`0011:C1`).
func dedupeUnknown(facts []unknownFact) []unknownFact {
	out := make([]unknownFact, 0, len(facts))
	seen := make(map[unknownFact]bool, len(facts))
	for _, f := range facts {
		if seen[f] {
			continue
		}
		seen[f] = true
		out = append(out, f)
	}
	slices.SortFunc(out, func(a, b unknownFact) int {
		if c := strings.Compare(a.Key, b.Key); c != 0 {
			return c
		}
		return strings.Compare(a.Reason, b.Reason)
	})
	return out
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
