package guard

// Phase 1 skeleton — SIGNATURES ONLY.
//
// This file exists so the RDR 0003 test suite compiles and each test fails
// on its own ASSERTION, naming the REQ its header quotes, rather than the
// whole package failing with one undefined-symbol error that attributes to
// no clause. Every function here returns a zero value and decides nothing.
//
// Phase 2 replaces this file. Nothing in it encodes a spec decision: the
// zero values are deliberately the WRONG answer for every REQ, which is
// what makes the suite red for the right reason.

import (
	"github.com/newcoinc/intrastate/internal/resolve"
	"github.com/newcoinc/intrastate/internal/table"
)

// --- the value seam ------------------------------------------------------

// Evaluator is RDR 0003's value-comparison seam: it decides value
// semantics over a PRESENT value and never reads the tag view.
type Evaluator struct{}

// Evaluate decides one atom against the present value its key holds.
func (Evaluator) Evaluate(atom resolve.GuardAtom, value string) resolve.GuardResult {
	return resolve.GuardUnevaluable
}

// --- the operator/kind matrix --------------------------------------------

// Shape names an operator's literal shape.
type Shape string

// The literal shapes the operator/kind matrix publishes.
const (
	ShapeScalar     Shape = "one typed scalar"
	ShapeScalarSet  Shape = "non-empty typed scalar set"
	ShapeInteger    Shape = "one typed integer"
	ShapeBoolean    Shape = "boolean"
	ShapeElementSet Shape = "non-empty typed element set"
)

// Operators returns this RDR's closed operator vocabulary.
func Operators() []string { return nil }

// KnownOperator reports whether token is in the closed vocabulary.
func KnownOperator(token string) bool { return false }

// Kinds returns the five value-kind tokens.
func Kinds() []string { return nil }

// Accepts reports whether the operator/kind matrix admits the pair.
func Accepts(operator, kind string) bool { return false }

// LiteralShape returns the operator's published literal shape.
func LiteralShape(operator string) Shape { return "" }

// --- the tag declaration model -------------------------------------------

// Declaration is one tag's declaration as this RDR reads it: the five
// fields the declaration model carries.
type Declaration struct {
	Provenance   table.Provenance
	Kind         string
	Domain       []string
	Elements     []string
	SingleValued bool
	Optional     bool
}

// DeclarationFields names the five fields the declaration model carries.
func DeclarationFields() []string { return nil }

// DeclarationOf reads one tag's declaration out of a loaded model.
func DeclarationOf(m *table.Model, key string) Declaration { return Declaration{} }

// AssignmentCount returns a dimension's assignment count under the
// per-kind table, and whether the declaration carries a finite domain.
func AssignmentCount(d table.TagDecl) (int, bool) { return 0, false }

// IntDomain enumerates a bounded int declaration's inclusive values.
func IntDomain(d table.TagDecl) []int { return nil }

// RequiresPredecessorWrite reports whether a provenance carries the
// reachable-predecessor-write obligation.
func RequiresPredecessorWrite(p table.Provenance) bool { return false }

// --- conformance ---------------------------------------------------------

// View is one evaluation view: tag key to rendered held value.
type View map[string]string

// Conforms reports whether v satisfies the declared model's two
// conjuncts: every always-present key is present, and every single-valued
// tag holds at most one declared domain value.
func Conforms(m *table.Model, v View) error { return nil }

// --- predicate semantic kinds --------------------------------------------

// SemanticKind is one of this RDR's predicate semantic kinds.
type SemanticKind string

// The predicate semantic kinds this RDR owns.
const (
	SemanticKindUnknownOperator         SemanticKind = "unknown-operator"
	SemanticKindOperatorKindMismatch    SemanticKind = "operator-kind-mismatch"
	SemanticKindLiteralParseFailure     SemanticKind = "literal-parse-failure"
	SemanticKindLiteralOutsideDomain    SemanticKind = "literal-outside-declared-domain"
	SemanticKindDeclarationDisagreement SemanticKind = "declaration-kind-disagreement"
	SemanticKindUnknownTag              SemanticKind = "unknown-tag"
)

// SemanticKinds returns the closed predicate-semantic-kind set.
func SemanticKinds() []SemanticKind { return nil }

// LoadCategoryFor maps a predicate semantic kind onto RDR 0002's load
// category that carries it.
func LoadCategoryFor(k SemanticKind) table.Category { return "" }

// --- identity ------------------------------------------------------------

// Identity is one guard atom's identity tuple: the enclosing row's source
// identity joined to the atom's own four fields.
type Identity string

// String renders the identity tuple.
func (i Identity) String() string { return string(i) }

// AtomIdentity builds the six-field identity tuple.
func AtomIdentity(row table.Row, atom table.Atom) Identity { return "" }

// SemanticKey renders an atom's semantic equality key: (tag, operator,
// literal), with neither block nor source identity.
func SemanticKey(atom table.Atom) string { return "" }

// --- assignments and denotations -----------------------------------------

// Assignment is one point of the scoped product: dimension to rendered
// held value, plus a presence marker for a key's presence dimension.
type Assignment map[string]string

// Presence names one value of a key's presence dimension.
type Presence string

// The two values of a presence dimension.
const (
	PresencePresent Presence = "present"
	PresenceAbsent  Presence = "absent"
	PresenceEither  Presence = "either"
)

// AssignmentSet is a set of assignments over a scoped product.
type AssignmentSet struct{}

// Len reports how many assignments the set holds.
func (AssignmentSet) Len() int { return 0 }

// Contains reports whether a is in the set.
func (AssignmentSet) Contains(a Assignment) bool { return false }

// Equal reports set equality.
func (AssignmentSet) Equal(other AssignmentSet) bool { return false }

// Subset reports whether the receiver is contained in other.
func (AssignmentSet) Subset(other AssignmentSet) bool { return false }

// Intersect returns the intersection.
func (AssignmentSet) Intersect(other AssignmentSet) AssignmentSet { return AssignmentSet{} }

// Union returns the union.
func (AssignmentSet) Union(other AssignmentSet) AssignmentSet { return AssignmentSet{} }

// Presence reports which presence value the set selects.
func (AssignmentSet) Presence() Presence { return "" }

// Values reports the value subset the set selects on its key's value
// dimension.
func (AssignmentSet) Values() []string { return nil }

// Projectable reports whether lint can project the atom at all.
func (AssignmentSet) Projectable() bool { return false }

// Denotation returns the subset of key's dimensions the atom denotes.
func Denotation(m *table.Model, key string, atom table.Atom) AssignmentSet {
	return AssignmentSet{}
}

// --- scoped row groups ---------------------------------------------------

// Selection is one selection context: the source state and the recognized
// outcome a row group shares.
type Selection struct {
	Model   string
	Outcome string
	Match   []table.Atom
}

// String renders the selection context.
func (Selection) String() string { return "" }

// MatchKeys names the match keys that scope the context.
func (Selection) MatchKeys() []string { return nil }

// Group is one scoped row group.
type Group struct {
	Context Selection
	Rows    []table.Row
}

// Groups partitions a model's rows into scoped row groups.
func Groups(m *table.Model) []Group { return nil }

// Dimensions names the participating guard value dimensions.
func Dimensions(m *table.Model, g Group) []string { return nil }

// PresenceDimensions names the keys contributing a presence dimension.
func PresenceDimensions(m *table.Model, g Group) []string { return nil }

// Product returns the group's scoped product.
func Product(m *table.Model, g Group) AssignmentSet { return AssignmentSet{} }

// Cardinality returns the scoped product's cardinality, and whether every
// participating dimension is finite.
func Cardinality(m *table.Model, g Group) (int, bool) { return 0, false }

// ProofCompletes reports whether the proof representation completes within
// the implementation's own resource budget.
func ProofCompletes(m *table.Model, g Group) bool { return false }

// AcceptedAssignments returns a row's accepted assignments.
func AcceptedAssignments(m *table.Model, row table.Row) AssignmentSet {
	return AssignmentSet{}
}

// CoverageUnion returns the union of the group's rows' accepted
// assignments.
func CoverageUnion(m *table.Model, g Group) AssignmentSet { return AssignmentSet{} }

// CoverageUnionFor returns the union for one declared rescuable class.
func CoverageUnionFor(m *table.Model, g Group, class string) AssignmentSet {
	return AssignmentSet{}
}

// CanRefuse reports whether a row carries a value atom over a key declared
// optional, in either guard block.
func CanRefuse(decls map[string]table.TagDecl, row table.Row) bool { return false }

// RescuableClasses names the failure classes an escape row may rescue.
func RescuableClasses() []string { return nil }

// --- the published bound -------------------------------------------------

// Bound returns the single integer cardinality bound this implementation
// publishes.
func Bound() int { return 0 }

// --- lint ----------------------------------------------------------------

// Code is one of RDR 0006's lint finding codes this RDR's semantics emit.
type Code string

// The finding codes this RDR's clauses name.
const (
	CodeUnprovableCoverage     Code = "graph-unprovable-coverage"
	CodeProductTooLarge        Code = "graph-product-too-large"
	CodeCoverageGap            Code = "graph-coverage-gap"
	CodeOverlap                Code = "graph-overlap"
	CodeCoverageClosedByEscape Code = "graph-coverage-closed-by-escape"
	CodeOwnedBeforeWrite       Code = "graph-owned-before-write"
	CodeVacuousExists          Code = "graph-vacuous-exists"
)

// Codes returns the finding codes this RDR's semantics emit.
func Codes() []Code { return nil }

// IsBlocking reports whether a code carries the blocking outcome.
func IsBlocking(c Code) bool { return false }

// IsRefuseOrDowngrade reports whether a code is a refuse-or-downgrade
// carrier.
func IsRefuseOrDowngrade(c Code) bool { return false }

// Severities returns this RDR's severity vocabulary.
func Severities() []string { return nil }

// Finding is one lint finding.
type Finding struct {
	Code         Code
	Context      string
	RuleIDs      []string
	Locators     []string
	Dimension    string
	Atom         *table.Atom
	Class        string
	Witness      Assignment
	ComputedSize int
	Bound        int
	Blocking     bool
}

// Verdict is a group's exhaustiveness verdict.
type Verdict string

// The verdicts a group report carries.
const (
	VerdictExhaustive Verdict = "exhaustive"
	VerdictWithheld   Verdict = "withheld"
	VerdictGap        Verdict = "gap"
	VerdictOverlap    Verdict = "overlap"
)

// GroupReport is one scoped row group's lint result.
type GroupReport struct {
	Context             string
	RuleIDs             []string
	Verdict             Verdict
	Green               bool
	Provable            bool
	ConformingViewsOnly bool
	ClosedByEscape      string
	Cardinality         int
	CoverageUnion       AssignmentSet
	Findings            []Finding
}

// Covers reports whether the group's green claim covers v.
func (GroupReport) Covers(v View) bool { return false }

// Lint runs the finite-domain lint over every scoped row group.
func Lint(m *table.Model) []GroupReport { return nil }
