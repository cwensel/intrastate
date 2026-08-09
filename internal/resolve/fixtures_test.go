package resolve_test

import (
	"github.com/newcoinc/intrastate/internal/resolve"
)

// This file holds only fixture builders. Every fixture is a value the
// caller supplies to the kernel: a normalized transition table, an
// accessor-produced owned snapshot, caller-supplied observed tags, and a
// guard-evaluation seam standing in for RDR 0003. Nothing here mocks the
// unit under test — the kernel itself is always the real Resolve.

// fixtureGuards is a table-driven stand-in for the RDR 0003 guard seam.
// Predicates it does not know are reported GuardUnevaluable, which is
// exactly the condition the kernel must answer with guard_unevaluable.
type fixtureGuards struct {
	decided map[string]bool
}

func (g fixtureGuards) Evaluate(guard string, _ resolve.TagSet) resolve.GuardResult {
	if guard == "" {
		return resolve.GuardTrue
	}
	v, ok := g.decided[guard]
	if !ok {
		return resolve.GuardUnevaluable
	}
	if v {
		return resolve.GuardTrue
	}
	return resolve.GuardFalse
}

// allGuardsTrue decides every named guard as holding.
func allGuardsTrue(names ...string) fixtureGuards {
	d := make(map[string]bool, len(names))
	for _, n := range names {
		d[n] = true
	}
	return fixtureGuards{decided: d}
}

// singleMatchTable models one legal edge: recognized outcome "successful"
// on an owned status:Draft artifact moves to status:Final and writes that
// owned tag back.
func singleMatchTable() resolve.Table {
	return resolve.Table{
		Revision: "rev-1",
		Outcomes: []string{"successful", "failed"},
		Rows: []resolve.Row{
			{
				RuleID:        "rdr.draft.successful",
				SourceLocator: "flows/rdr.toml:12",
				Outcome:       "successful",
				Match:         []resolve.Tag{{Key: "status", Value: "Draft"}},
				RequiresOwned: []string{"status"},
				NextTags:      []resolve.Tag{{Key: "status", Value: "Final"}},
				Writes:        []resolve.Tag{{Key: "status", Value: "Final"}},
			},
		},
	}
}

// legalInput is the canonical legal input tuple: it selects exactly one
// row of singleMatchTable.
func legalInput() resolve.Input {
	return resolve.Input{
		Flow:       "rdr",
		Table:      singleMatchTable(),
		Owned:      []resolve.Tag{{Key: "status", Value: "Draft"}},
		Observed:   []resolve.Tag{{Key: "reviews", Value: "2"}},
		Recognized: "successful",
		Guards:     allGuardsTrue(),
	}
}

// noMatchTable declares the outcome in its alphabet but models no row
// that matches a status:Draft artifact — a genuine zero-match, distinct
// from an unmodeled outcome.
func noMatchTable() resolve.Table {
	return resolve.Table{
		Revision: "rev-nomatch",
		Outcomes: []string{"successful"},
		Rows: []resolve.Row{
			{
				RuleID:        "rdr.final.successful",
				SourceLocator: "flows/rdr.toml:20",
				Outcome:       "successful",
				Match:         []resolve.Tag{{Key: "status", Value: "Final"}},
				RequiresOwned: []string{"status"},
				NextTags:      []resolve.Tag{{Key: "status", Value: "Landed"}},
				Writes:        []resolve.Tag{{Key: "status", Value: "Landed"}},
			},
		},
	}
}

// ambiguousTable models two rows that both match the same tag-set for the
// same recognized outcome. Both retain distinct source identity so an
// ambiguous refusal can report the conflict.
func ambiguousTable() resolve.Table {
	return resolve.Table{
		Revision: "rev-ambiguous",
		Outcomes: []string{"successful"},
		Rows: []resolve.Row{
			{
				RuleID:        "rdr.draft.successful.a",
				SourceLocator: "flows/rdr.toml:12",
				Outcome:       "successful",
				Match:         []resolve.Tag{{Key: "status", Value: "Draft"}},
				RequiresOwned: []string{"status"},
				NextTags:      []resolve.Tag{{Key: "status", Value: "Final"}},
				Writes:        []resolve.Tag{{Key: "status", Value: "Final"}},
			},
			{
				RuleID:        "rdr.draft.successful.b",
				SourceLocator: "flows/rdr.toml:30",
				Outcome:       "successful",
				Match:         []resolve.Tag{{Key: "status", Value: "Draft"}},
				RequiresOwned: []string{"status"},
				NextTags:      []resolve.Tag{{Key: "status", Value: "Review"}},
				Writes:        []resolve.Tag{{Key: "status", Value: "Review"}},
			},
		},
	}
}

// missingOwnedTable's only candidate requires an owned tag the accessor
// snapshot will not carry.
func missingOwnedTable() resolve.Table {
	return resolve.Table{
		Revision: "rev-missing-owned",
		Outcomes: []string{"successful"},
		Rows: []resolve.Row{
			{
				RuleID:        "rdr.gate.successful",
				SourceLocator: "flows/rdr.toml:40",
				Outcome:       "successful",
				Match:         []resolve.Tag{{Key: "status", Value: "Draft"}},
				RequiresOwned: []string{"status", "gate"},
				NextTags:      []resolve.Tag{{Key: "status", Value: "Final"}},
				Writes:        []resolve.Tag{{Key: "status", Value: "Final"}},
			},
		},
	}
}

// unevaluableGuardTable's only candidate references a predicate the guard
// seam will report as undecidable.
func unevaluableGuardTable() resolve.Table {
	return resolve.Table{
		Revision: "rev-guard",
		Outcomes: []string{"successful"},
		Rows: []resolve.Row{
			{
				RuleID:        "rdr.guarded.successful",
				SourceLocator: "flows/rdr.toml:50",
				Outcome:       "successful",
				Match:         []resolve.Tag{{Key: "status", Value: "Draft"}},
				RequiresOwned: []string{"status"},
				Guard:         "reviews >= quorum",
				NextTags:      []resolve.Tag{{Key: "status", Value: "Final"}},
				Writes:        []resolve.Tag{{Key: "status", Value: "Final"}},
			},
		},
	}
}

// twoRowsOneGuardFalseTable models two rows matching the same tag-set,
// each behind a distinct guard, so the guard seam's verdict decides how
// many candidates survive.
func twoRowsOneGuardFalseTable() resolve.Table {
	return resolve.Table{
		Revision: "rev-guarded",
		Outcomes: []string{"successful"},
		Rows: []resolve.Row{
			{
				RuleID:        "rdr.draft.successful.fast",
				SourceLocator: "flows/rdr.toml:60",
				Outcome:       "successful",
				Match:         []resolve.Tag{{Key: "status", Value: "Draft"}},
				RequiresOwned: []string{"status"},
				Guard:         "is-fast-lane",
				NextTags:      []resolve.Tag{{Key: "status", Value: "Final"}},
				Writes:        []resolve.Tag{{Key: "status", Value: "Final"}},
			},
			{
				RuleID:        "rdr.draft.successful.slow",
				SourceLocator: "flows/rdr.toml:70",
				Outcome:       "successful",
				Match:         []resolve.Tag{{Key: "status", Value: "Draft"}},
				RequiresOwned: []string{"status"},
				Guard:         "is-slow-lane",
				NextTags:      []resolve.Tag{{Key: "status", Value: "Review"}},
				Writes:        []resolve.Tag{{Key: "status", Value: "Review"}},
			},
		},
	}
}

// observedSensitiveTable's only row can match only if the caller-supplied
// observed tag reaches the evaluation view.
func observedSensitiveTable() resolve.Table {
	return resolve.Table{
		Revision: "rev-observed",
		Outcomes: []string{"successful"},
		Rows: []resolve.Row{
			{
				RuleID:        "rdr.draft.reviewed",
				SourceLocator: "flows/rdr.toml:80",
				Outcome:       "successful",
				Match: []resolve.Tag{
					{Key: "status", Value: "Draft"},
					{Key: "reviews", Value: "2"},
				},
				RequiresOwned: []string{"status"},
				NextTags:      []resolve.Tag{{Key: "status", Value: "Final"}},
				Writes:        []resolve.Tag{{Key: "status", Value: "Final"}},
			},
		},
	}
}

// recognizedTagSensitiveTable's only row can match only if the freshly
// recognized outcome is merged into the evaluation view as a tag.
func recognizedTagSensitiveTable() resolve.Table {
	return resolve.Table{
		Revision: "rev-recognized",
		Outcomes: []string{"successful"},
		Rows: []resolve.Row{
			{
				RuleID:        "rdr.draft.on-recognized-tag",
				SourceLocator: "flows/rdr.toml:85",
				Outcome:       "successful",
				Match: []resolve.Tag{
					{Key: "status", Value: "Draft"},
					{Key: "recognized", Value: "successful"},
				},
				RequiresOwned: []string{"status"},
				NextTags:      []resolve.Tag{{Key: "status", Value: "Final"}},
				Writes:        []resolve.Tag{{Key: "status", Value: "Final"}},
			},
		},
	}
}

// escapeRow builds a row modeled as an escape for the given failure class.
// It matches the base tag-set so it is a genuine escape candidate.
func escapeRow(ruleID, locator string, class resolve.RefusalKind) resolve.Row {
	return resolve.Row{
		RuleID:        ruleID,
		SourceLocator: locator,
		Outcome:       "successful",
		Match:         []resolve.Tag{{Key: "status", Value: "Draft"}},
		RequiresOwned: []string{"status"},
		NextTags:      []resolve.Tag{{Key: "status", Value: "Blocked"}},
		Writes:        []resolve.Tag{{Key: "status", Value: "Blocked"}},
		Escape:        []resolve.RefusalKind{class},
	}
}

// noMatchInput: the outcome is in the alphabet, but no row matches the
// assembled tag-set.
func noMatchInput() resolve.Input {
	in := legalInput()
	in.Table = noMatchTable()
	return in
}

// ambiguousInput: two rows match after guard evaluation.
func ambiguousInput() resolve.Input {
	in := legalInput()
	in.Table = ambiguousTable()
	return in
}

// missingOwnedInput: the candidate requires an owned tag the accessor
// snapshot does not carry.
func missingOwnedInput() resolve.Input {
	in := legalInput()
	in.Table = missingOwnedTable()
	return in
}

// unevaluableGuardInput: the candidate's guard is unknown to the seam.
func unevaluableGuardInput() resolve.Input {
	in := legalInput()
	in.Table = unevaluableGuardTable()
	in.Guards = allGuardsTrue() // knows nothing about "reviews >= quorum"
	return in
}

// unmodeledOutcomeInput: the recognized outcome is outside the table's
// declared recognized-outcome alphabet.
func unmodeledOutcomeInput() resolve.Input {
	in := legalInput()
	in.Table = resolve.Table{
		Revision: "rev-alphabet",
		Outcomes: []string{"successful", "failed"},
		Rows:     singleMatchTable().Rows,
	}
	in.Recognized = "abandoned"
	return in
}

// refusalInputBuilders returns one builder per kernel refusal kind, so
// each case can be constructed fresh for replay comparisons.
func refusalInputBuilders() map[string]func() resolve.Input {
	return map[string]func() resolve.Input{
		"no_match":                noMatchInput,
		"ambiguous_match":         ambiguousInput,
		"owned_state_unavailable": missingOwnedInput,
		"guard_unevaluable":       unevaluableGuardInput,
		"unmodeled_outcome":       unmodeledOutcomeInput,
	}
}

// allRefusalInputs materializes one input per refusal kind.
func allRefusalInputs() map[string]resolve.Input {
	out := map[string]resolve.Input{}
	for name, build := range refusalInputBuilders() {
		out[name] = build()
	}
	return out
}

// allDispositionInputs is every refusal input plus the legal one.
func allDispositionInputs() map[string]resolve.Input {
	out := allRefusalInputs()
	out["legal"] = legalInput()
	return out
}

// deepCopyTags copies a tag slice so a test can prove Resolve did not
// mutate the caller's values.
func deepCopyTags(in []resolve.Tag) []resolve.Tag {
	if in == nil {
		return nil
	}
	out := make([]resolve.Tag, len(in))
	copy(out, in)
	return out
}

// deepCopyInput copies every mutable part of the input tuple.
func deepCopyInput(in resolve.Input) resolve.Input {
	out := in
	out.Owned = deepCopyTags(in.Owned)
	out.Observed = deepCopyTags(in.Observed)
	out.Table.Outcomes = append([]string(nil), in.Table.Outcomes...)
	if in.Table.Rows != nil {
		rows := make([]resolve.Row, len(in.Table.Rows))
		for i, r := range in.Table.Rows {
			r.Match = deepCopyTags(r.Match)
			r.NextTags = deepCopyTags(r.NextTags)
			r.Writes = deepCopyTags(r.Writes)
			r.RequiresOwned = append([]string(nil), r.RequiresOwned...)
			r.Escape = append([]resolve.RefusalKind(nil), r.Escape...)
			rows[i] = r
		}
		out.Table.Rows = rows
	}
	return out
}
