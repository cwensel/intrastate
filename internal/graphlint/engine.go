package graphlint

import (
	"slices"
	"strings"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
	"github.com/newcoinc/intrastate/internal/table"
)

// Request is the normalized graph-lint request. It carries a normalized
// model value, never a path and never command state.
type Request struct {
	Model *table.Model
}

// NewRequest is THE request builder (`0006:AP`, SC-6). Every path that
// lints — the root command, any later alias, any CI invocation — builds
// its request here, so no second constructor can diverge.
func NewRequest(m *table.Model) Request { return Request{Model: m} }

// Report is one lint run's result.
type Report struct {
	// ModelID is the `[model].id` from the model itself, never a path.
	ModelID string
	// Findings is every finding the run decided, blocking and advisory,
	// in deterministic finding-identity order.
	Findings []clierr.Finding
}

// Blocking returns the report's blocking findings, in the same order.
func (r Report) Blocking() []clierr.Finding { return r.partition(true) }

// Advisory returns the report's non-blocking findings, in the same order.
func (r Report) Advisory() []clierr.Finding { return r.partition(false) }

// partition splits the findings by severity. The two halves are exhaustive
// and disjoint: every finding carries a code from exactly one tier.
func (r Report) partition(blocking bool) []clierr.Finding {
	var out []clierr.Finding
	for _, f := range r.Findings {
		if IsBlocking(f.Code) == blocking {
			out = append(out, f)
		}
	}
	return out
}

// Run is THE engine entry point (`0006:AP`, SC-6). It runs every mandatory
// invariant over the request's normalized graph and returns every defect it
// can decide in ONE pass — never the first it encounters.
//
// Emission is complete, not first-failure: withholding a group's
// exhaustiveness claim does not suppress overlap, coverage, or further
// withholding findings for that group, and each unprovable dimension, each
// refusing row, each overlapping pair, and any coverage gap over a provable
// product is its own finding. The emitted set therefore never depends on
// row or dimension iteration order.
func Run(req Request) Report {
	m := req.Model
	if m == nil {
		return Report{}
	}

	a := newAnalysis(m)

	// Every check runs; none short-circuits another. The order below is
	// the taxonomy's, not a precedence: the emitted set is sorted by
	// finding identity at the end, so this order is invisible downstream.
	a.checkDanglingEdge()
	a.checkNodeCeiling()
	a.checkSingleValuedState()
	a.checkAlwaysPresentOwned()
	a.checkDeadEnd()
	a.checkTerminalEscape()
	a.checkGroups()
	a.checkUnreachableRules()

	sortFindings(a.findings)
	return Report{ModelID: m.ID, Findings: a.findings}
}

// sortFindings orders findings by the finding-identity tuple: model id,
// invariant code, source rule/context id or graph element id, then the
// normalized predicate/write fingerprint. A source rule/context id sorts
// before any graph element id, so the namespace is a leading discriminator
// within the (model, code) bucket.
func sortFindings(findings []clierr.Finding) {
	slices.SortFunc(findings, func(a, b clierr.Finding) int {
		return strings.Compare(identityKey(a), identityKey(b))
	})
}

// identityKey renders the finding-identity tuple as one sortable string.
// The trailing fields are not part of the normative tuple; they are the
// tie-break that makes the order TOTAL, so two findings agreeing on the
// tuple still emit in a fixed order rather than a map-iteration one.
func identityKey(f clierr.Finding) string {
	id, namespace := f.Rule, "0"
	if id == "" {
		namespace = "1"
		id = f.Element
		if id == "" {
			id = f.Span
		}
	}
	return strings.Join([]string{
		f.Model, f.Code, namespace, id, f.Fingerprint,
		f.Class, f.Reason, f.Dimension, f.Key, f.Operator, f.Literal, f.Block,
		f.Message,
	}, "\x00")
}

// Fingerprint renders the canonical SORTABLE predicate/write serialization
// the finding-identity tuple closes on — never a hash (`0006:C15`).
//
// It is RDR 0002's canonical atom sort — `(key, block, operator token,
// literal)` — over the row's predicate atoms, followed by its next-state
// tags in the same canonical set-literal form. The atom content stays
// READABLE, so a reviewer can tell two fingerprints apart by eye and a sort
// over them is meaningful rather than an arbitrary permutation.
func Fingerprint(row table.Row) string {
	atoms := slices.Clone(row.Atoms)
	slices.SortFunc(atoms, compareAtoms)

	var b strings.Builder
	for _, a := range atoms {
		b.WriteString(a.Key)
		b.WriteString("|")
		b.WriteString(string(a.Block))
		b.WriteString("|")
		b.WriteString(a.Operator)
		b.WriteString("|")
		b.WriteString(strings.Join(canonicalValues(a.Literal), ","))
		b.WriteString(";")
	}
	b.WriteString("#")
	for _, t := range canonicalTags(row.NextTags) {
		b.WriteString(t.Key)
		b.WriteString("=")
		b.WriteString(strings.Join(canonicalValues(t.Value), ","))
		b.WriteString(";")
	}
	return b.String()
}

// compareAtoms is RDR 0002's canonical atom sort: key, then block, then
// operator token, then literal.
func compareAtoms(a, b table.Atom) int {
	if c := strings.Compare(a.Key, b.Key); c != 0 {
		return c
	}
	if c := strings.Compare(string(a.Block), string(b.Block)); c != 0 {
		return c
	}
	if c := strings.Compare(a.Operator, b.Operator); c != 0 {
		return c
	}
	return strings.Compare(
		strings.Join(canonicalValues(a.Literal), ","),
		strings.Join(canonicalValues(b.Literal), ","))
}

// canonicalTags sorts a tag-value sequence by key so the fingerprint's
// next-state half is order-independent too.
func canonicalTags(in []table.TagValue) []table.TagValue {
	out := slices.Clone(in)
	slices.SortFunc(out, func(a, b table.TagValue) int { return strings.Compare(a.Key, b.Key) })
	return out
}
