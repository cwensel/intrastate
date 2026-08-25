// Package graphlint is RDR 0006's graph-lint engine: the blocking
// acceptance gate over the normalized transition model.
//
// The boundary this package fixes (`0006:C2`, TD): it receives a
// NORMALIZED graph value — never Cobra command state and never sparse
// TOML — and defines no second sparse-source parser and no parallel
// transition semantics. RDR 0002's loader produces the model; RDR 0003's
// `internal/guard` owns the scoped row group, the tag declaration model,
// the finite-domain product, and the can-refuse predicate, all of which
// this package CITES rather than restates. What this package adds is the
// owned-state reachability relation, the invariant taxonomy, and the
// blocking acceptance verdict.
package graphlint

import (
	"slices"

	"github.com/newcoinc/intrastate/internal/guard"
)

// --- the finding taxonomy ------------------------------------------------

// The ten blocking finding codes (Technical Design, finding-code table).
const (
	CodeDanglingEdge       = "graph-dangling-edge"
	CodeDeadEnd            = "graph-dead-end"
	CodeOverlap            = "graph-overlap"
	CodeCoverageGap        = "graph-coverage-gap"
	CodeUnprovableCoverage = "graph-unprovable-coverage"
	CodeSingleValuedState  = "graph-single-valued-state"
	CodeAlwaysPresentOwned = "graph-always-present-owned"
	CodeOwnedBeforeWrite   = "graph-owned-before-write"
	CodeTerminalEscape     = "graph-terminal-escape"
	CodeProductTooLarge    = "graph-product-too-large"
)

// The four advisory finding codes. The tier is CLOSED at these four
// (`0006:C17`).
const (
	CodeCoverageClosedByEscape = "graph-coverage-closed-by-escape"
	CodeRedundantRow           = "graph-redundant-row"
	CodeUnreachableRule        = "graph-unreachable-rule"
	CodeVacuousAtom            = "graph-vacuous-atom"
)

// AggregateCode is the one aggregate CLIError code a blocking run returns
// (`0006:C14`).
const AggregateCode = "graph-lint-failed"

// The two severities. There is no third tier (`0006:C8`, `0006:C17`).
const (
	SeverityBlocking = "blocking"
	SeverityInfo     = "info"
)

// The closed, append-only `reason` set `graph-unprovable-coverage`
// carries (Technical Design, reason table). It is carried on that code
// alone.
const (
	ReasonDimensionNotFinite = "dimension-not-finite"
	ReasonTagNotSingleValued = "tag-not-single-valued"
	ReasonRowCanRefuse       = "row-can-refuse"
)

// blockingCodes is the blocking tier, in taxonomy order.
var blockingCodes = []string{
	CodeDanglingEdge,
	CodeDeadEnd,
	CodeOverlap,
	CodeCoverageGap,
	CodeUnprovableCoverage,
	CodeSingleValuedState,
	CodeAlwaysPresentOwned,
	CodeOwnedBeforeWrite,
	CodeTerminalEscape,
	CodeProductTooLarge,
}

// advisoryCodes is the closed four-member advisory tier.
var advisoryCodes = []string{
	CodeCoverageClosedByEscape,
	CodeRedundantRow,
	CodeUnreachableRule,
	CodeVacuousAtom,
}

// reasons is the closed `reason` discriminator set.
var reasons = []string{
	ReasonDimensionNotFinite,
	ReasonTagNotSingleValued,
	ReasonRowCanRefuse,
}

// severities is the two-valued severity vocabulary.
var severities = []string{SeverityBlocking, SeverityInfo}

// BlockingCodes returns the blocking code set, in taxonomy order.
func BlockingCodes() []string { return slices.Clone(blockingCodes) }

// AdvisoryCodes returns the closed four-member advisory tier.
func AdvisoryCodes() []string { return slices.Clone(advisoryCodes) }

// Reasons returns the closed `reason` discriminator set.
func Reasons() []string { return slices.Clone(reasons) }

// Severities returns the two-valued severity vocabulary.
func Severities() []string { return slices.Clone(severities) }

// IsBlocking reports whether code carries the blocking severity. A code
// outside the taxonomy is not blocking: the two tiers are the whole
// vocabulary, and a stranger cannot refuse a model.
func IsBlocking(code string) bool { return slices.Contains(blockingCodes, code) }

// severityFor maps a code onto its declared severity.
func severityFor(code string) string {
	if IsBlocking(code) {
		return SeverityBlocking
	}
	return SeverityInfo
}

// --- the published implementation constants ------------------------------

// nodeCeiling is the model-independent reachable-node ceiling. The
// owned-state lattice is finite but can be astronomically wide, so the
// traversal stops at this many merged nodes and reports rather than
// running unboundedly.
const nodeCeiling = 4096

// ProductBound returns the model-independent product bound above which
// lint declines to prove a group's coverage (`0006:C12`). It is an
// implementation constant, published in the command's help output — never
// a per-model input.
//
// It IS RDR 0003's published bound rather than a second number: the proof
// representation this engine declines over is `internal/guard`'s
// enumerating one, so a bound of its own would name a threshold nothing
// applies and the two would drift.
func ProductBound() int { return guard.Bound() }

// NodeCeiling returns the model-independent reachable-node ceiling
// (Performance Expectations). Exceeding it is `graph-product-too-large`
// naming the traversal, never an unbounded run.
func NodeCeiling() int { return nodeCeiling }
