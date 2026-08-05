// Package resolver provides the pure transition-resolution kernel.
package resolver

import (
	"maps"
	"slices"
)

// TagSet is a state snapshot keyed by tag name.
type TagSet map[string]string

// Tag is a freshly recognized named outcome.
type Tag struct {
	Name  string
	Value string
}

// Provenance identifies the source of a tag used by a guard.
type Provenance string

const (
	// ProvenanceOwned selects a tag from the accessor-produced owned snapshot.
	ProvenanceOwned Provenance = "owned"
	// ProvenanceObserved selects a caller-supplied observed tag.
	ProvenanceObserved Provenance = "observed"
	// ProvenanceRecognized selects the freshly recognized outcome tag.
	ProvenanceRecognized Provenance = "recognized"
)

// TagPredicate requires one provenance-qualified tag value.
type TagPredicate struct {
	Provenance Provenance
	Name       string
	Value      string
}

// Guard is a conjunction of tag predicates.
type Guard struct {
	All         []TagPredicate
	Unevaluable bool
}

// OwnedSnapshot is the accessor-produced owned state supplied to Resolve.
type OwnedSnapshot struct {
	Available bool
	Tags      TagSet
}

// OwnedTagWrite describes an inert write for the accessor layer to apply.
type OwnedTagWrite struct {
	Role  string
	Name  string
	Value string
}

// RefusalKind is a stable, value-level reason resolution could not select an edge.
type RefusalKind string

const (
	// RefusalNoMatch means no ordinary or corresponding escape edge matched.
	RefusalNoMatch RefusalKind = "no_match"
	// RefusalAmbiguousMatch means more than one candidate edge matched.
	RefusalAmbiguousMatch RefusalKind = "ambiguous_match"
	// RefusalOwnedStateUnavailable means the accessor supplied no owned snapshot.
	RefusalOwnedStateUnavailable RefusalKind = "owned_state_unavailable"
	// RefusalGuardUnevaluable means a relevant symbolic guard could not be evaluated.
	RefusalGuardUnevaluable RefusalKind = "guard_unevaluable"
	// RefusalUnmodeledOutcome means the recognized outcome is outside the declared alphabet.
	RefusalUnmodeledOutcome RefusalKind = "unmodeled_outcome"
)

// Edge is one normalized transition-table candidate.
type Edge struct {
	Outcome   string
	EscapeFor RefusalKind
	Guard     Guard
	NextTags  TagSet
	Writes    []OwnedTagWrite
}

// Input is the complete, explicit resolution tuple.
type Input struct {
	FlowID        string
	TableRevision string
	Owned         OwnedSnapshot
	Observed      TagSet
	Recognized    Tag
	Outcomes      []string
	Table         []Edge
}

// TransitionPlan is the selected next state and its inert owned-tag writes.
type TransitionPlan struct {
	NextTags TagSet
	Writes   []OwnedTagWrite
}

// Refusal is a modeled resolution failure.
type Refusal struct {
	Kind RefusalKind
}

// Disposition contains exactly one transition plan or one refusal.
type Disposition struct {
	Plan    *TransitionPlan
	Refusal *Refusal
}

// Resolve selects exactly one matching transition or returns a modeled refusal.
func Resolve(input Input) (Disposition, error) {
	if !input.Owned.Available {
		return refuse(RefusalOwnedStateUnavailable), nil
	}

	if !modelsOutcome(input.Outcomes, input.Recognized.Value) {
		return refuse(RefusalUnmodeledOutcome), nil
	}

	matches, unevaluable := matchingEdges(input, "")
	if unevaluable {
		return refuse(RefusalGuardUnevaluable), nil
	}
	if len(matches) == 1 {
		return plan(matches[0]), nil
	}

	failedKind := RefusalNoMatch
	if len(matches) > 1 {
		failedKind = RefusalAmbiguousMatch
	}

	escapes, escapeUnevaluable := matchingEdges(input, failedKind)
	if escapeUnevaluable {
		return refuse(RefusalGuardUnevaluable), nil
	}
	if len(escapes) == 1 {
		return plan(escapes[0]), nil
	}
	if len(escapes) > 1 {
		return refuse(RefusalAmbiguousMatch), nil
	}

	return refuse(failedKind), nil
}

func modelsOutcome(outcomes []string, outcome string) bool {
	return slices.Contains(outcomes, outcome)
}

func matchingEdges(input Input, escapeFor RefusalKind) ([]Edge, bool) {
	matches := make([]Edge, 0)
	for _, edge := range input.Table {
		if edge.Outcome != input.Recognized.Value || edge.EscapeFor != escapeFor {
			continue
		}

		matchesGuard, unevaluable := evaluateGuard(input, edge.Guard)
		if unevaluable {
			return nil, true
		}
		if matchesGuard {
			matches = append(matches, edge)
		}
	}
	return matches, false
}

func evaluateGuard(input Input, guard Guard) (bool, bool) {
	if guard.Unevaluable {
		return false, true
	}

	for _, predicate := range guard.All {
		matches, unevaluable := evaluatePredicate(input, predicate)
		if unevaluable {
			return false, true
		}
		if !matches {
			return false, false
		}
	}
	return true, false
}

func evaluatePredicate(input Input, predicate TagPredicate) (bool, bool) {
	switch predicate.Provenance {
	case ProvenanceOwned:
		value, ok := input.Owned.Tags[predicate.Name]
		return ok && value == predicate.Value, false
	case ProvenanceObserved:
		value, ok := input.Observed[predicate.Name]
		return ok && value == predicate.Value, false
	case ProvenanceRecognized:
		matchesName := input.Recognized.Name == predicate.Name
		matchesValue := input.Recognized.Value == predicate.Value
		return matchesName && matchesValue, false
	default:
		return false, true
	}
}

func plan(edge Edge) Disposition {
	return Disposition{
		Plan: &TransitionPlan{
			NextTags: maps.Clone(edge.NextTags),
			Writes:   append([]OwnedTagWrite{}, edge.Writes...),
		},
	}
}

func refuse(kind RefusalKind) Disposition {
	return Disposition{
		Refusal: &Refusal{Kind: kind},
	}
}
