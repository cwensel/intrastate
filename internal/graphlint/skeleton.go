// Package graphlint is RDR 0006's graph-lint engine: the blocking
// acceptance gate over the normalized transition model.
//
// SKELETON ONLY (Stage 8, Phase 1). Every declaration below is a
// signature returning a zero value so the RDR 0006 test suite COMPILES
// and each test fails on its own assertion — naming the REQ its header
// quotes — rather than the whole package failing with one
// undefined-symbol error that attributes to no clause. Phase 2 replaces
// each body. Nothing here implements a check.
//
// The boundary this package fixes (`0006:C2`, TD): it receives a
// NORMALIZED graph value — never Cobra command state and never sparse
// TOML — and defines no second sparse-source parser and no parallel
// transition semantics.
package graphlint

import (
	"github.com/newcoinc/intrastate/internal/cli/clierr"
	"github.com/newcoinc/intrastate/internal/table"
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
// carries (Technical Design, reason table).
const (
	ReasonDimensionNotFinite = "dimension-not-finite"
	ReasonTagNotSingleValued = "tag-not-single-valued"
	ReasonRowCanRefuse       = "row-can-refuse"
)

// BlockingCodes returns the blocking code set, in taxonomy order.
func BlockingCodes() []string { return nil }

// AdvisoryCodes returns the closed four-member advisory tier.
func AdvisoryCodes() []string { return nil }

// Reasons returns the closed `reason` discriminator set.
func Reasons() []string { return nil }

// Severities returns the two-valued severity vocabulary.
func Severities() []string { return nil }

// IsBlocking reports whether code carries the blocking severity.
func IsBlocking(code string) bool { return false }

// --- the published implementation constants ------------------------------

// ProductBound returns the model-independent product bound above which
// lint declines to prove a group's coverage (`0006:C12`). It is an
// implementation constant, published in the command's help output — never
// a per-model input.
func ProductBound() int { return 0 }

// NodeCeiling returns the model-independent reachable-node ceiling
// (Performance Expectations). Exceeding it is `graph-product-too-large`
// naming the traversal, never an unbounded run.
func NodeCeiling() int { return 0 }

// --- the single request builder and the single engine entry point --------

// Request is the normalized graph-lint request. It carries a normalized
// model value, never a path and never command state.
type Request struct {
	Model *table.Model
}

// NewRequest is THE request builder (`0006:AP`, SC-6). Every path that
// lints — the root command, any later alias, any CI invocation — builds
// its request here, so no second constructor can diverge.
func NewRequest(m *table.Model) Request { return Request{} }

// Report is one lint run's result.
type Report struct {
	// ModelID is the `[model].id` from the model itself, never a path.
	ModelID string
	// Findings is every finding the run decided, blocking and advisory,
	// in deterministic finding-identity order.
	Findings []clierr.Finding
}

// Blocking returns the report's blocking findings, in the same order.
func (r Report) Blocking() []clierr.Finding { return nil }

// Advisory returns the report's non-blocking findings, in the same order.
func (r Report) Advisory() []clierr.Finding { return nil }

// Run is THE engine entry point (`0006:AP`, SC-6). It runs every
// mandatory invariant over the request's normalized graph and returns
// every defect it can decide in ONE pass — never the first it encounters.
func Run(req Request) Report { return Report{} }

// --- the reachability relation (Load-Bearing Decisions) ------------------

// Node is one abstract owned-state: per owned tag, either absent, or held
// with the set of possible declared values. A tag with no finite domain
// abstracts to held/absent.
type Node struct {
	// Values maps each held owned tag to its set of possible values,
	// sorted and duplicate-free. A key absent from the map is absent in
	// this node.
	Values map[string][]string
}

// Reach runs the owned-state traversal to a fixpoint over MERGED nodes
// and returns the reachable node set. The root is the declared initial
// owned state; an edge is a normalized non-escape row whose match pattern
// over owned tags is satisfiable in the source node; escape rows are
// edges too, and they are self-loops.
func Reach(m *table.Model) []Node { return nil }

// Fingerprint renders the canonical SORTABLE predicate/write
// serialization the finding-identity tuple closes on — never a hash
// (`0006:C15`).
func Fingerprint(row table.Row) string { return "" }
