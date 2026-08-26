package guard

import (
	"slices"

	"github.com/newcoinc/intrastate/internal/resolve"
	"github.com/newcoinc/intrastate/internal/table"
)

// Code is one of RDR 0006's lint finding codes this RDR's semantics emit.
// This RDR mints none of its own: it states the semantics and RDR 0006
// owns the codes. The set below is the req-list ASSUMPTION's enumeration
// verbatim, plus `graph-vacuous-atom` — RDR 0006's advisory-tier name for
// the vacuity case RDR 0003 requires be reported rather than rejected
// (`0006:722`, `0006:739`). A code outside RDR 0006's taxonomy would be one
// this RDR minted in a space it disclaimed owning.
type Code string

// The finding codes this RDR's clauses name.
const (
	CodeUnprovableCoverage     Code = "graph-unprovable-coverage"
	CodeProductTooLarge        Code = "graph-product-too-large"
	CodeCoverageGap            Code = "graph-coverage-gap"
	CodeOverlap                Code = "graph-overlap"
	CodeCoverageClosedByEscape Code = "graph-coverage-closed-by-escape"
	CodeOwnedBeforeWrite       Code = "graph-owned-before-write"
	CodeVacuousAtom            Code = "graph-vacuous-atom"
)

var codes = []Code{
	CodeUnprovableCoverage,
	CodeProductTooLarge,
	CodeCoverageGap,
	CodeOverlap,
	CodeCoverageClosedByEscape,
	CodeOwnedBeforeWrite,
	CodeVacuousAtom,
}

// Codes returns the finding codes this RDR's semantics emit.
func Codes() []Code { return slices.Clone(codes) }

// IsBlocking reports whether a code carries THIS RDR's inability-to-prove
// outcome. The two inability-to-prove carriers are blocking; the rest report
// a defect or an observation without refusing a claim that was never made.
//
// This is DELIBERATELY NARROWER than RDR 0006's blocking tier, and is not a
// mirror of it. RDR 0006's taxonomy also classifies `graph-coverage-gap`,
// `graph-overlap`, and `graph-owned-before-write` as blocking; the authority
// for that classification is `graphlint`'s own table (`0006:C8`, `0006:C17`),
// which every consumer of the lint's machine-readable output reads. This
// predicate answers a different question: which codes carry REQ-65's
// refuse-or-downgrade outcome, the one this RDR states is a SINGLE blocking
// outcome with no non-blocking warning category (`0003:C20`). Widening it to
// RDR 0006's tier would erase that distinction and contradict REQ-66, whose
// test asserts this package's severity vocabulary has no advisory tier at
// all — the opposite of RDR 0006's two-tier taxonomy.
//
// The two predicates share a name and answer different questions. Consumers
// wanting RDR 0006's tier MUST read `graphlint`, not this.
func IsBlocking(c Code) bool {
	return c == CodeUnprovableCoverage || c == CodeProductTooLarge
}

// IsRefuseOrDowngrade reports whether a code is a refuse-or-downgrade
// carrier. "Refuse" and "downgrade" are ONE outcome, not an author's
// choice: every case this RDR sends there produces a blocking
// inability-to-prove finding, and no non-blocking warning category is
// minted for any of them.
func IsRefuseOrDowngrade(c Code) bool { return IsBlocking(c) }

// severities is this RDR's severity vocabulary. There is no advisory tier:
// "downgrade" names the same blocking outcome as "refuse" on both sides of
// the seam, and never a silent or advisory one.
var severities = []string{"blocking", "report"}

// Severities returns this RDR's severity vocabulary.
func Severities() []string { return slices.Clone(severities) }

// Finding is one lint finding.
type Finding struct {
	Code     Code
	Context  string
	RuleIDs  []string
	Locators []string
	// Dimension names the guard dimension a per-dimension finding is about.
	Dimension string
	// Atom names the atom a per-atom finding is about — the refusing atom
	// on a withheld claim.
	Atom *table.Atom
	// Class names the declared rescuable class a per-class finding is about.
	Class string
	// Witness carries a concrete uncovered assignment from the scoped
	// product. Naming the group without one does not satisfy the gap clause.
	Witness Assignment
	// ComputedSize and Bound are the two figures the over-large refusal
	// reports, so an author can tell an over-large product from an
	// undeclared dimension. They are unavailable — and stay zero — for a
	// group whose product has no cardinality.
	ComputedSize int
	Bound        int
	// Blocking marks THIS RDR's inability-to-prove outcome, matching
	// IsBlocking — not RDR 0006's blocking tier, which is wider and whose
	// authority is `graphlint`'s taxonomy. A gap or overlap finding here
	// leaves it false while RDR 0006 classifies those codes blocking; that
	// is the deliberate narrowing REQ-65/REQ-66 state, not an unset flag.
	Blocking bool
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
	Context string
	RuleIDs []string
	Verdict Verdict
	Green   bool
	// Provable reports whether the group has a scoped product to prove
	// over. A group containing a can-refuse row has none.
	Provable bool
	// ConformingViewsOnly records that a green claim asserts coverage over
	// conforming views only — the premise every lint claim here is
	// conditional on.
	ConformingViewsOnly bool
	// ClosedByEscape names the bare escape row that closed the group's
	// coverage, so a reviewer can tell a group proved over its declared
	// domains from one closed by a catch-all.
	ClosedByEscape string
	Cardinality    int
	CoverageUnion  AssignmentSet
	Findings       []Finding
}

// Covers reports whether the group's green claim covers v. A non-conforming
// view is outside the claim: this RDR promises nothing about a view that
// violates the declarations.
func (r GroupReport) Covers(v View) bool {
	if !r.Green {
		return false
	}
	return r.CoverageUnion.Contains(Assignment(v))
}

// Lint runs the finite-domain lint over every scoped row group.
//
// The claim is DEFAULT-ON for every group whose participating dimensions
// are all finitely declared: an opt-in flag would let the guarantee be
// silently skipped exactly where it matters, and leaving a dimension
// undeclared is not an opt-out either — it makes the proof unavailable,
// which is the loud blocking inability-to-prove outcome. Every group is
// therefore either proved or carrying a blocking finding; there is no
// third, silently-unchecked state.
//
// Every defect the pass can decide is reported in ONE pass over the group,
// not the first one encountered.
func Lint(m *table.Model) []GroupReport {
	written := writtenKeys(m)
	out := make([]GroupReport, 0, len(Groups(m)))
	for _, g := range Groups(m) {
		out = append(out, lintGroup(m, g, written))
	}
	return out
}

// writtenKeys names every owned tag some rule's write block assigns. It is
// RDR 0006's owned-state reachability relation as this model exposes it;
// this RDR cites that relation and defines no second one.
func writtenKeys(m *table.Model) map[string]bool {
	out := map[string]bool{}
	for _, row := range m.Rows {
		for _, w := range row.Writes {
			out[w.Key] = true
		}
		for _, n := range row.NextTags {
			out[n.Key] = true
		}
	}
	return out
}

func lintGroup(m *table.Model, g Group, written map[string]bool) GroupReport {
	r := GroupReport{
		Context: g.Context.String(),
		RuleIDs: ruleIDs(g.Rows),
	}

	r.Findings = append(r.Findings, ownedBeforeWriteFindings(m, g, written)...)
	r.Findings = append(r.Findings, vacuousExistsFindings(m, g)...)

	if len(Dimensions(m, g)) == 0 {
		// A group over no guard dimension at all makes no exhaustiveness
		// claim: the claim is default-on for every group "whose
		// participating dimensions are all finitely declared", and there
		// is nothing here for a claim to range over. Certifying it green
		// would put a proof beside groups that carry one, and the reader
		// could not tell the two apart.
		//
		// OVERLAP still runs. Withholding the CLAIM is not withholding
		// every check: the exhaustiveness claim carries the
		// finitely-declared precondition, while the overlap invariant
		// ("no finite-domain input assignment may enable two ordinary
		// rows") carries none — it reads declarations, never a product.
		// Two guard-atom-free rows in one selection context are enabled by
		// EVERY assignment, so returning here emitted no verdict and no
		// finding at all: lint exited clean on a model the kernel always
		// refuses `ambiguous_match`. That is the false green this record
		// exists to prevent, reached by skipping a claim-independent check.
		r.Findings = append(r.Findings, overlapFindings(m, g)...)
		if countCodeIn(r.Findings, CodeOverlap) > 0 {
			r.Verdict = VerdictOverlap
		}
		r.CoverageUnion = newSet(nil)
		return r
	}

	// The bound is tested BEFORE the projection scan, and after every
	// participating dimension is known finite — which is exactly the order
	// REQ-93 states ("lint MUST compute the cardinality and test the bound
	// after every participating dimension is known finite"). Both halves of
	// the order matter: testing it earlier would report a computed size for
	// a product carrying an unprovable dimension, and testing it later
	// would send the projection scan through a product this implementation
	// has already declined to enumerate — the naive enumeration on the
	// refusal path that the record's mitigation forbids.
	if card, ok := Cardinality(m, g); ok && card > Bound() {
		r.Verdict = VerdictWithheld
		r.Cardinality = card
		r.CoverageUnion = newSet(nil)
		r.Findings = append(r.Findings, Finding{
			Code:         CodeProductTooLarge,
			Context:      r.Context,
			RuleIDs:      r.RuleIDs,
			Locators:     locators(g.Rows),
			ComputedSize: card,
			Bound:        Bound(),
			Blocking:     true,
		})
		r.Findings = append(r.Findings, withholdingFindings(m, g)...)
		return r
	}

	unprovable := unprovableFindings(m, g)
	refusing := withholdingFindings(m, g)
	r.Findings = append(r.Findings, unprovable...)
	r.Findings = append(r.Findings, refusing...)

	if len(unprovable) > 0 || len(refusing) > 0 {
		// The group has no provable product. The coverage-gap finding is
		// scoped to a provable product, so none is emitted — but overlap
		// among the group's decidable rows is unaffected and still MUST be.
		r.Verdict = VerdictWithheld
		r.CoverageUnion = newSet(nil)
		r.Findings = append(r.Findings, overlapFindings(m, g)...)
		return r
	}

	r.Provable = true
	card, ok := Cardinality(m, g)
	if !ok {
		// Unreachable: every dimension is known finite above. Kept so the
		// size-bearing diagnostic can never report a quantity the table
		// defines no value for.
		r.Verdict = VerdictWithheld
		r.CoverageUnion = newSet(nil)
		return r
	}
	r.Cardinality = card

	// The over-large refusal is already decided above, before the
	// projection scan ran. Reaching here means the product is finite, under
	// the bound, and fully projectable.
	r.CoverageUnion = CoverageUnion(m, g)
	overlaps := overlapFindings(m, g)
	r.Findings = append(r.Findings, overlaps...)
	gaps, closedBy := coverageFindings(m, g, ordinaryOverlaps(overlaps))
	r.Findings = append(r.Findings, gaps...)
	r.ClosedByEscape = closedBy

	switch {
	case countCodeIn(r.Findings, CodeCoverageGap) > 0:
		r.Verdict = VerdictGap
	case countCodeIn(r.Findings, CodeOverlap) > 0:
		r.Verdict = VerdictOverlap
	default:
		r.Verdict = VerdictExhaustive
		r.Green = true
		r.ConformingViewsOnly = true
	}
	return r
}

// unprovableFindings emits one blocking finding per unprovable dimension,
// naming that dimension. Each is its own finding: lint never stops at the
// first, and never treats the covered examples as complete.
func unprovableFindings(m *table.Model, g Group) []Finding {
	var out []Finding
	for _, key := range Dimensions(m, g) {
		if reason := unprovableReason(m, g, key); reason != "" {
			out = append(out, Finding{
				Code:      CodeUnprovableCoverage,
				Context:   g.Context.String(),
				RuleIDs:   ruleIDs(g.Rows),
				Locators:  locators(g.Rows),
				Dimension: key,
				Blocking:  true,
			})
		}
	}
	return out
}

// unprovableReason reports why a dimension cannot be proved, or "" when it
// can. The two cases take the SAME outcome: a dimension with no finite
// declared domain, and one carrying a value atom lint cannot project.
func unprovableReason(m *table.Model, g Group, key string) string {
	if _, ok := AssignmentCount(m.Tags[key]); !ok {
		return "no finite declared domain"
	}
	if unprovableDimension(m, g, key) {
		return "unprojectable atom"
	}
	return ""
}

// withholdingFindings emits the withholding finding for EACH refusing row,
// naming the row and the atom that can refuse. A withheld claim must be
// observable, not silent: emitting nothing does not satisfy the clause,
// since an exit code alone cannot distinguish a withheld claim from a
// proved one.
//
// The population is EVERY row in the group, escape rows included — the
// participation clause's population, not the ordinary-row population the
// overlap check uses. Reading the two-population split into the narrowing
// would certify green exactly the group the runtime refuses.
func withholdingFindings(m *table.Model, g Group) []Finding {
	rows := slices.Clone(g.Rows)
	slices.SortFunc(rows, func(a, b table.Row) int {
		return slices.Compare([]string{a.RuleID, a.SourceLocator}, []string{b.RuleID, b.SourceLocator})
	})

	var out []Finding
	for _, row := range rows {
		if !CanRefuse(m.Tags, row) {
			continue
		}
		out = append(out, Finding{
			Code:      CodeUnprovableCoverage,
			Context:   g.Context.String(),
			RuleIDs:   []string{row.RuleID},
			Locators:  []string{row.SourceLocator},
			Dimension: refusingAtom(m.Tags, row).Key,
			Atom:      refusingAtom(m.Tags, row),
			Blocking:  true,
		})
	}
	return out
}

// overlapFindings emits the pairwise overlap findings, in TWO separate
// populations because the runtime never mixes them.
//
// The ordinary population excludes escape rows: an escape row overlapping a
// guarded row is not a runtime ambiguity and MUST NOT be reported as one.
// The escape population is partitioned per DECLARED failure class, with a
// row declaring several classes placed in EACH of those populations, and
// the pairwise check run within each independently — a flat pairing
// over-reports two rows sharing no class, and a per-row partition
// under-reports two rows colliding on one shared class of several. A pair
// overlapping in more than one class is one finding per class.
// ordinaryOverlaps reports whether the group's ORDINARY population carries an
// overlap — the precondition that makes the `ambiguous_match` arm reachable.
// Escape-population overlaps carry a class and do not count: the kernel
// reaches ambiguity from the ordinary rows, so only their pairwise overlap
// makes an `ambiguous_match` rescue row reachable.
func ordinaryOverlaps(findings []Finding) bool {
	return slices.ContainsFunc(findings, func(f Finding) bool {
		return f.Code == CodeOverlap && f.Class == ""
	})
}

func overlapFindings(m *table.Model, g Group) []Finding {
	var ordinary, escape []table.Row
	for _, row := range g.Rows {
		if row.Kind() == table.KindEscape {
			escape = append(escape, row)
			continue
		}
		ordinary = append(ordinary, row)
	}

	out := pairwiseOverlaps(m, g, ordinary, "")
	for _, class := range RescuableClasses() {
		var population []table.Row
		for _, row := range escape {
			if slices.Contains(row.Escape, class) {
				population = append(population, row)
			}
		}
		out = append(out, pairwiseOverlaps(m, g, population, class)...)
	}
	return out
}

// pairwiseOverlaps reports each pair of rows whose accepted assignments
// intersect non-emptily. Overlap is not an equality test at all — it
// intersects the rows' accepted assignment SETS — so two byte-identical
// guards in two rule ids correctly produce a finding naming both rather
// than being deduplicated into one row.
func pairwiseOverlaps(m *table.Model, g Group, rows []table.Row, class string) []Finding {
	rows = slices.Clone(rows)
	slices.SortFunc(rows, func(a, b table.Row) int {
		return slices.Compare([]string{a.RuleID, a.SourceLocator}, []string{b.RuleID, b.SourceLocator})
	})

	var out []Finding
	for i := range rows {
		a := acceptedIn(m, g, rows[i])
		if !a.Projectable() {
			continue
		}
		for _, other := range rows[i+1:] {
			b := acceptedIn(m, g, other)
			if !b.Projectable() || a.Intersect(b).Len() == 0 {
				continue
			}
			out = append(out, Finding{
				Code:     CodeOverlap,
				Context:  g.Context.String(),
				RuleIDs:  []string{rows[i].RuleID, other.RuleID},
				Locators: []string{rows[i].SourceLocator, other.SourceLocator},
				Class:    class,
			})
		}
	}
	return out
}

// coverageFindings emits one gap finding per unclosed rescuable class, and
// reports the bare escape row that closed a class, if any.
//
// Coverage is `union(row_i accepted assignments) == scoped product`. The
// union is computed per (group × declared rescuable class), since an escape
// row closes coverage only for the classes it can actually rescue. A gap
// has no contributing predicate — it is an ABSENCE — so it is attributed by
// the selection context, EVERY source rule id in the group, and at least
// one concrete uncovered assignment from the product.
// The `ambiguous_match` arm is checked only where it is REACHABLE — for a
// group whose ordinary population carries an overlap — and is treated as
// vacuously closed otherwise. RDR 0006 states the arm mechanics
// (`0006::Normative Contracts`, the per-class coverage clause) and REQ-77
// (`0003:C15`) cites them rather than restating them, so 0006's reachability
// precondition governs emission here. Demanding a rescue row for an arm the
// kernel can never reach would mint a row `graph-unreachable-rule` then
// flags. REQ-77's own clause — that the union is per (group x declared
// class), and that a row declaring one class does not close another's arm —
// is unaffected: the two unions stay distinct either way.
func coverageFindings(m *table.Model, g Group, ordinaryOverlap bool) ([]Finding, string) {
	product := Product(m, g)
	var out []Finding
	var closedBy string

	for _, class := range RescuableClasses() {
		if class == string(resolve.KindAmbiguousMatch) && !ordinaryOverlap {
			continue
		}
		union := CoverageUnionFor(m, g, class)
		if union.Equal(product) {
			if row := bareEscapeFor(g, class); row != "" && !ordinaryClosesAlone(m, g, product) {
				closedBy = row
			}
			continue
		}
		witness, ok := uncovered(product, union)
		if !ok {
			continue
		}
		out = append(out, Finding{
			Code:     CodeCoverageGap,
			Context:  g.Context.String(),
			RuleIDs:  ruleIDs(g.Rows),
			Locators: locators(g.Rows),
			Class:    class,
			Witness:  witness,
		})
	}

	if closedBy != "" {
		// A bare escape row is not a silent opt-out from the guarantee: the
		// closure it performs is reported as an OBSERVABLE result, and a
		// bare green for such a group does not satisfy the clause.
		out = append(out, Finding{
			Code:     CodeCoverageClosedByEscape,
			Context:  g.Context.String(),
			RuleIDs:  []string{closedBy},
			Locators: locatorOf(g, closedBy),
		})
	}
	return out, closedBy
}

// ordinaryClosesAlone reports whether the group's non-escape rows already
// close coverage, in which case no escape row closed anything.
func ordinaryClosesAlone(m *table.Model, g Group, product AssignmentSet) bool {
	union := newSet(productDims(m, g))
	for _, row := range g.Rows {
		if row.Kind() == table.KindEscape {
			continue
		}
		accepted := acceptedIn(m, g, row)
		if accepted.Projectable() {
			union = union.Union(accepted)
		}
	}
	return union.Equal(product)
}

// bareEscapeFor names the escape row declaring class that carries NO guard
// atoms — the one that denotes the whole scoped product and therefore
// closes coverage by itself.
func bareEscapeFor(g Group, class string) string {
	var found []string
	for _, row := range g.Rows {
		if row.Kind() != table.KindEscape || !slices.Contains(row.Escape, class) {
			continue
		}
		if len(guardAtoms(row)) == 0 {
			found = append(found, row.RuleID)
		}
	}
	slices.Sort(found)
	if len(found) == 0 {
		return ""
	}
	return found[0]
}

// uncovered returns one concrete assignment in the product the union does
// not cover, chosen deterministically so the witness does not depend on map
// iteration order.
func uncovered(product, union AssignmentSet) (Assignment, bool) {
	keys := make([]string, 0, len(product.members))
	for k := range product.members {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		a := product.members[k]
		if !union.Contains(a) {
			return a, true
		}
	}
	return nil, false
}

// vacuousExistsFindings reports an `exists` atom over a key declared
// always-present. Such a key contributes no `{absent}` assignment, so the
// atom is well-formed but VACUOUS — lint reports it as such rather than
// rejecting it.
func vacuousExistsFindings(m *table.Model, g Group) []Finding {
	rows := slices.Clone(g.Rows)
	slices.SortFunc(rows, func(a, b table.Row) int {
		return slices.Compare([]string{a.RuleID}, []string{b.RuleID})
	})

	var out []Finding
	for _, row := range rows {
		atoms := guardAtoms(row)
		slices.SortFunc(atoms, compareAtoms)
		for i := range atoms {
			if atoms[i].Operator != resolve.OpExists {
				continue
			}
			if !m.Tags[atoms[i].Key].Required {
				continue
			}
			out = append(out, Finding{
				Code:      CodeVacuousAtom,
				Context:   g.Context.String(),
				RuleIDs:   []string{row.RuleID},
				Locators:  []string{row.SourceLocator},
				Dimension: atoms[i].Key,
				Atom:      &atoms[i],
			})
		}
	}
	return out
}

// ownedBeforeWriteFindings reports a row that MATCHES an owned tag no
// reachable predecessor sets or preserves before the match. Predicate lint
// distinguishes the three provenances: recognized tags are fresh event
// inputs and observed tags are re-read before matching, so neither carries
// the obligation.
func ownedBeforeWriteFindings(m *table.Model, g Group, written map[string]bool) []Finding {
	rows := slices.Clone(g.Rows)
	slices.SortFunc(rows, func(a, b table.Row) int {
		return slices.Compare([]string{a.RuleID}, []string{b.RuleID})
	})

	var out []Finding
	for _, row := range rows {
		var keys []string
		for _, atom := range row.Atoms {
			if atom.Block != table.BlockMatch {
				continue
			}
			if !RequiresPredecessorWrite(m.Tags[atom.Key].Provenance) {
				continue
			}
			if written[atom.Key] || slices.Contains(keys, atom.Key) {
				continue
			}
			keys = append(keys, atom.Key)
		}
		slices.Sort(keys)
		for _, key := range keys {
			out = append(out, Finding{
				Code:      CodeOwnedBeforeWrite,
				Context:   g.Context.String(),
				RuleIDs:   []string{row.RuleID},
				Locators:  []string{row.SourceLocator},
				Dimension: key,
			})
		}
	}
	return out
}

// --- small helpers -------------------------------------------------------

func ruleIDs(rows []table.Row) []string {
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		if !slices.Contains(out, r.RuleID) {
			out = append(out, r.RuleID)
		}
	}
	slices.Sort(out)
	return out
}

func locators(rows []table.Row) []string {
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		if r.SourceLocator != "" && !slices.Contains(out, r.SourceLocator) {
			out = append(out, r.SourceLocator)
		}
	}
	slices.Sort(out)
	return out
}

func locatorOf(g Group, ruleID string) []string {
	for _, r := range g.Rows {
		if r.RuleID == ruleID && r.SourceLocator != "" {
			return []string{r.SourceLocator}
		}
	}
	return nil
}

func countCodeIn(findings []Finding, code Code) int {
	var n int
	for _, f := range findings {
		if f.Code == code {
			n++
		}
	}
	return n
}
