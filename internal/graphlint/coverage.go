package graphlint

import (
	"fmt"
	"slices"
	"strings"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/guard"
	"github.com/cwensel/intrastate/internal/table"
)

// checkCoverage runs invariant 4 over one reachable group: the guard
// exhaustiveness claim, the three withholding triggers, and the per-class
// coverage arms.
//
// The claim is DEFAULT-ON for every group whose participating dimensions are
// all finitely declared — never an opt-in annotation, and leaving a
// dimension undeclared is not an opt-out either: it makes the proof
// unavailable, which is the loud blocking outcome.
//
// Coverage is a UNIVERSAL claim, so it is computed from the group's authored
// rows and their declared guard domains alone and never from a reachability
// node. Reachability decides only whether the group is proven at all; it
// never enters the computation. Match keys are not product dimensions and
// contribute no assignment to either side.
func (a *analysis) checkCoverage(g guard.Group) {
	dims := guard.Dimensions(a.model, g)
	if len(dims) == 0 {
		if table.IsDecisionTable(a.model) {
			// A decision table IS its guard product: the cells are the
			// authored guard dimensions, and a group ranging over none has
			// nothing for coverage to be proved over. The empty product's
			// single empty assignment would be "covered" by any row at all,
			// so reporting it covered is a vacuous green — a table whose
			// author discriminated with `[rule.match.<key>]` atoms (which
			// scope the group and contribute no dimension, `0006:C7`) would
			// lint clean while proving nothing.
			//
			// The arm binds THIS CLASS ONLY. Over a state machine a
			// zero-dimension group is legitimate — the owned-state
			// traversal is what carries the proof there — and an
			// unconditional arm would redden 151 such groups across every
			// loadable fixture.
			//
			// It RETURNS rather than falling through. Emitting and then
			// running `emitCoverageArms` would yield both this finding and
			// `graph-coverage-closed-by-escape` over a zero-dimension group
			// carrying a bare escape row, which is the double-report
			// `0010:C5` forbids; ordering alone does not deliver the
			// precedence, the return does.
			a.emitNoParticipatingDimension(g)
			return
		}
		// A group over no guard dimension still has a scoped product — the
		// empty one, carrying the single empty assignment — so the per-class
		// arms below are decidable and the claim is not withheld. What is
		// absent is anything for an unprovable DIMENSION to be found in, so
		// the two scans are skipped rather than run over nothing.
		a.emitCoverageArms(g)
		return
	}

	// The bound is tested BEFORE the projection scan and after every
	// participating dimension is known finite. Testing it earlier would
	// report a computed size for a product carrying an unprovable
	// dimension; testing it later would send the projection scan through a
	// product this implementation has already declined to enumerate.
	if card, ok := guard.Cardinality(a.model, g); ok && card > ProductBound() {
		// REQ-94: the over-large refusal "is simply not among the findings"
		// for a group carrying an unprovable dimension, and REQ-93's second
		// clause forbids reporting a computed size for one. A size counted
		// over a product whose atoms do not project is a number lint cannot
		// witness. REQ-58 settles that an unprojectable atom is unprovable
		// "exactly as a dimension with no finite declared domain does", so
		// the unprojectable-ATOM case is in scope, not only the
		// no-finite-DOMAIN case.
		//
		// The scan is STRUCTURAL — the marker arm of unprovableReason, which
		// reads declarations only and never calls Denotation — so REQ-93's
		// cost ordering holds and no enumeration reaches the refusal path
		// (ADV-2). It also corrects the remedy: "add the marker" rather than
		// "narrow a declared domain", which was never the defect here.
		if a.emitStructurallyUnprovable(g) {
			a.emitWithholdings(g)
			return
		}

		a.emit(clierr.Finding{
			Code:      CodeProductTooLarge,
			Rule:      firstRuleID(g),
			Element:   g.Context.String(),
			Dimension: strings.Join(dims, ","),
			Message: fmt.Sprintf("the declared finite product of group %s is "+
				"%d assignments, above the published product bound of %d; "+
				"narrow a declared domain", g.Context.String(), card,
				ProductBound()),
		})
		// The group has no product to prove over, so no coverage arm runs.
		// The withholding findings still do: each is its own defect and one
		// refusal does not suppress another.
		a.emitWithholdings(g)
		return
	}

	unprovable := a.emitUnprovableDimensions(g)
	withheld := a.emitWithholdings(g)
	if unprovable || withheld {
		// The group has no provable product. The coverage-gap finding is
		// scoped to a provable product, so none is emitted — and neither is
		// the bare-escape closure: WITHHOLDING DOMINATES CLOSURE, so a group
		// both closable by a bare escape row and carrying a row that can
		// refuse resolves to the blocking withheld claim rather than to a
		// success carrying `graph-coverage-closed-by-escape`.
		return
	}

	a.emitCoverageArms(g)
}

// emitNoParticipatingDimension reports the GROUP-level unprovable-coverage
// finding a decision-table group with zero participating dimensions takes
// (`0010:C5`).
//
// It supplies its OWN message rather than routing through
// unprovableMessage, whose `default` limb phrases the remedy as "declare
// the domain" — advice for a dimension that does not exist here. The remedy
// this arm states is the one C5's own prose gives: author the
// discriminators as guard atoms.
//
// The finding is group-scoped, so it carries no `Dimension` and no `Key`:
// there is no dimension to name, and borrowing one would misreport which
// declaration is at fault. It names the class nowhere — `0010:BR6` (a lint
// finding announcing the class) was rejected, so the message speaks about
// the GROUP's authoring, not about what the model declared itself to be.
func (a *analysis) emitNoParticipatingDimension(g guard.Group) {
	a.emit(clierr.Finding{
		Code:    CodeUnprovableCoverage,
		Reason:  ReasonNoParticipatingDimension,
		Rule:    firstRuleID(g),
		Element: g.Context.String(),
		Message: fmt.Sprintf("group %s ranges over no participating "+
			"dimension, so its scoped product carries the single empty "+
			"assignment and its coverage cannot be proved; author the "+
			"discriminating keys as `[rule.guard.all.<key>]` atoms, which "+
			"is what makes them product dimensions", g.Context.String()),
	})
}

// emitUnprovableDimensions reports one blocking finding per unprovable
// participating dimension, naming that dimension and carrying the reason
// that tells the author what to do. It returns whether any fired.
//
// The two triggers take the same code and different reasons because their
// remedies differ: a dimension with no finite declared domain is fixed by
// declaring the domain, while an `eq`/`in`/comparison atom over a tag
// lacking the single-valued marker is fixed by adding the marker.
func (a *analysis) emitUnprovableDimensions(g guard.Group) bool {
	var any bool
	for _, key := range guard.Dimensions(a.model, g) {
		reason, atom := a.unprovableReason(g, key)
		if reason == "" {
			continue
		}
		any = true
		f := clierr.Finding{
			Code:      CodeUnprovableCoverage,
			Reason:    reason,
			Dimension: key,
			Key:       key,
			Element:   g.Context.String(),
			Message:   a.unprovableMessage(reason, key, g),
		}
		if atom != nil {
			f.Operator = atom.Operator
			f.Literal = strings.Join(atom.Literal, ",")
			f.Block = string(atom.Block)
		}
		a.emit(f)
	}
	return any
}

// emitStructurallyUnprovable emits a `graph-unprovable-coverage` finding for
// every participating dimension whose unprovability is decidable from
// DECLARATIONS ALONE — a value atom over a tag carrying no single-valued
// marker (`0003::A21`). It reports whether it emitted any.
//
// It is the subset of unprovableReason that is safe on the over-large refusal
// path: the `Denotation` arm would materialize the very product the bound has
// already declined, which is the naive enumeration ADV-2 closed. The
// no-finite-DOMAIN arm is omitted for a different reason — a dimension with no
// finite domain has no cardinality, so a group carrying one never reaches this
// branch at all (`guard.Cardinality` returns ok=false).
func (a *analysis) emitStructurallyUnprovable(g guard.Group) bool {
	var any bool
	for _, key := range guard.Dimensions(a.model, g) {
		decl := a.model.Tags[key]
		if decl.SingleValued {
			continue
		}
		// `guard.SingleValueOperator` is the one authoritative
		// classification: an atom over the tag's single held value has no
		// projection when the tag is NOT declared single-valued, and takes
		// the blocking inability-to-prove outcome (`0003::A21`).
		atoms := a.groupAtomsOver(g, key)
		for i := range atoms {
			if !guard.SingleValueOperator(atoms[i].Operator) {
				continue
			}
			any = true
			a.emit(clierr.Finding{
				Code:      CodeUnprovableCoverage,
				Reason:    ReasonTagNotSingleValued,
				Dimension: key,
				Key:       key,
				Element:   g.Context.String(),
				Message:   a.unprovableMessage(ReasonTagNotSingleValued, key, g),
				Operator:  atoms[i].Operator,
				Literal:   strings.Join(atoms[i].Literal, ","),
				Block:     string(atoms[i].Block),
			})
			break
		}
	}
	return any
}

// unprovableReason reports why a participating dimension cannot be proved,
// or "" when it can, along with the atom that made it unprovable where one
// did.
func (a *analysis) unprovableReason(g guard.Group, key string) (string, *table.Atom) {
	decl := a.model.Tags[key]

	// The DOMAIN arm is decided first. A dimension with no finite declared
	// domain has nothing for a single-valued marker to range over, so
	// reporting the marker there would name a remedy that does not fix it:
	// adding `single_valued` to a `scalar` still leaves lint with no domain
	// to enumerate. The author must declare the domain, and that is what the
	// reason says.
	if _, ok := guard.AssignmentCount(decl); !ok {
		return ReasonDimensionNotFinite, nil
	}

	// `0003::A21`: over a FINITE dimension, a value atom on a tag lacking
	// the single-valued marker is the one that has no projection, and its
	// remedy is the marker.
	atoms := a.groupAtomsOver(g, key)
	if !decl.SingleValued {
		for i := range atoms {
			if guard.SingleValueOperator(atoms[i].Operator) {
				return ReasonTagNotSingleValued, &atoms[i]
			}
		}
	}

	// A dimension whose declaration is finite can still carry an atom lint
	// cannot project — an operator the declared kind admits no denotation
	// for. That is the same unavailable proof as an undeclared domain, and
	// it takes the same reason.
	for i := range atoms {
		if !guard.Denotation(a.model, key, atoms[i]).Projectable() {
			return ReasonDimensionNotFinite, &atoms[i]
		}
	}
	return "", nil
}

// groupAtomsOver returns every guard atom in the group naming key, in
// canonical atom order so the atom a finding names does not depend on row
// iteration.
func (a *analysis) groupAtomsOver(g guard.Group, key string) []table.Atom {
	var out []table.Atom
	for _, row := range g.Rows {
		for _, atom := range guardAtomsOf(row) {
			if atom.Key == key {
				out = append(out, atom)
			}
		}
	}
	slices.SortFunc(out, compareAtoms)
	return out
}

// unprovableMessage phrases the dimension arms. Neither reads as a
// withholding: both name a declaration the author can supply.
func (a *analysis) unprovableMessage(reason, key string, g guard.Group) string {
	switch reason {
	case ReasonTagNotSingleValued:
		return fmt.Sprintf("the tag %q carries no single-valued marker, so a "+
			"value atom over it has no projection and the coverage of group "+
			"%s cannot be proved; add the marker", key, g.Context.String())
	default:
		return fmt.Sprintf("the dimension %q has no finite declared domain "+
			"lint can project, so the coverage of group %s cannot be proved; "+
			"declare the domain", key, g.Context.String())
	}
}

// emitWithholdings reports the withheld exhaustiveness claim for EACH row in
// the group that can refuse `guard_unevaluable`, naming the participating row
// and the refusing atom. It returns whether any fired.
//
// The population is EVERY row in the group, escape rows included: the
// narrowing is over participation, not over the two overlap populations, and
// reading the split into it would certify green exactly the group the runtime
// refuses. The "can refuse" test is SYNTACTIC over the declared optionality
// field — it never consults the reachability relation.
//
// There is no non-blocking tier for this class, and no second code: the
// withheld claim takes `graph-unprovable-coverage` with the `row-can-refuse`
// reason.
func (a *analysis) emitWithholdings(g guard.Group) bool {
	rows := slices.Clone(g.Rows)
	slices.SortFunc(rows, func(x, y table.Row) int { return strings.Compare(x.RuleID, y.RuleID) })

	var any bool
	for _, row := range rows {
		if !guard.CanRefuse(a.model.Tags, row) {
			continue
		}
		atom := a.refusingAtom(row)
		if atom == nil {
			continue
		}
		any = true
		a.emit(clierr.Finding{
			Code:        CodeUnprovableCoverage,
			Reason:      ReasonRowCanRefuse,
			Rule:        rowIdentity(row),
			Span:        row.SourceLocator,
			Dimension:   atom.Key,
			Key:         atom.Key,
			Operator:    atom.Operator,
			Literal:     strings.Join(atom.Literal, ","),
			Block:       string(atom.Block),
			Fingerprint: Fingerprint(row),
			// The model is fine and lint declines to promise, so the message
			// says the claim is WITHHELD rather than reading as an authoring
			// error.
			Message: fmt.Sprintf("coverage claim withheld for row %q: the tag "+
				"%q is declared optional, so the atom can refuse "+
				"`guard_unevaluable` at runtime and lint declines to promise "+
				"exhaustiveness the runtime does not deliver",
				row.RuleID, atom.Key),
		})
	}
	return any
}

// refusingAtom returns the row's first value atom over an optional key, in
// canonical atom order, so the withholding finding names one atom rather
// than whichever the row happened to author first.
func (a *analysis) refusingAtom(row table.Row) *table.Atom {
	atoms := guardAtomsOf(row)
	slices.SortFunc(atoms, compareAtoms)
	for i := range atoms {
		if atoms[i].Operator == "exists" {
			// An existence atom over an absent key is DECIDED from presence
			// alone, never unevaluable, so it cannot make a row refuse.
			continue
		}
		if decl, ok := a.model.Tags[atoms[i].Key]; ok && !decl.Required {
			return &atoms[i]
		}
	}
	return nil
}

// emitCoverageArms proves `union(row_i accepted assignments) == scoped
// product` per (scoped row group × declared rescuable class).
//
// An escape row closes coverage only for the failure classes it DECLARES, so
// a row declaring one class must not close the group's other arms. Escape
// rows contribute their accepted assignments to the union like any other
// row: there is no separate "an escape row exists" disjunct.
//
// The `ambiguous_match` arm is only demanded where it is REACHABLE. Overlap
// is itself blocking, so a group whose ordinary rows pass invariant 3 cannot
// produce `ambiguous_match` at runtime, and requiring an escape row for that
// arm would demand a row `graph-unreachable-rule` would then flag. The arm is
// therefore checked only for a group carrying a `graph-overlap` finding, and
// treated as vacuously closed otherwise.
func (a *analysis) emitCoverageArms(g guard.Group) {
	dims := guard.Dimensions(a.model, g)
	product := guard.Product(a.model, g)
	if !product.Projectable() {
		return
	}

	overlapping := a.groupHasOverlap(g)
	var closedBy string
	for _, class := range guard.RescuableClasses() {
		if class == string(ambiguousMatch) && !overlapping {
			continue
		}
		union := a.coverageUnionFor(g, class)
		if union.Equal(product) {
			// A bare escape row in the closing population is reported
			// whatever else closed the arm: the clause is that a bare green
			// MUST NOT satisfy it, so the reader can tell a group carrying a
			// catch-all from one that does not without inspecting the model.
			if row := bareEscapeFor(g, class); row != "" {
				closedBy = row
			}
			continue
		}
		// The gap is attributed to the group's ROWS when its own rows leave
		// assignments uncovered, and to the selection CONTEXT alone when
		// they do not — the second arm's gap is the absent rescue row, and
		// there is no authored row to name for a row that was never
		// written. REQ-127's fallback is what keeps the second arm
		// actionable: it carries the graph element id instead.
		gap := clierr.Finding{
			Code:      CodeCoverageGap,
			Element:   g.Context.String(),
			Class:     class,
			Dimension: strings.Join(dims, ","),
			Message: fmt.Sprintf("group %s over rows %v leaves %d of %d "+
				"assignments in its scoped product uncovered for the %s arm; "+
				"the coverage union must equal the scoped product",
				g.Context.String(), ruleIDsOf(g), product.Len()-union.Len(),
				product.Len(), class),
		}
		if !a.ordinaryClosesAlone(g, product) {
			gap.Rule = firstRuleID(g)
		}
		a.emit(gap)
	}

	if closedBy == "" {
		return
	}
	// A bare escape row is not a silent opt-out from the guarantee: the
	// closure it performs is reported as an observable result, so a bare
	// green does not satisfy the clause.
	a.emit(clierr.Finding{
		Code:    CodeCoverageClosedByEscape,
		Rule:    closedBy,
		Element: g.Context.String(),
		Message: fmt.Sprintf("the coverage of group %s is closed by the bare "+
			"escape row %q rather than proved over its declared domains",
			g.Context.String(), closedBy),
	})
}

// coverageUnionFor unions the accepted assignments of the rows that close
// one declared rescuable class's arm. Escape rows contribute their
// assignments like any other row — there is no separate "an escape row
// exists" disjunct — but WHICH rows close an arm differs by class, because
// the two refusals arise from opposite conditions.
//
//   - `no_match` arises exactly where NO row accepts the assignment, so an
//     ordinary row accepting it is what keeps the refusal from occurring at
//     all. The ordinary population therefore closes this arm, alongside the
//     escape rows declaring the class.
//   - `ambiguous_match` arises where the group's ORDINARY rows overlap. The
//     rows that caused the ambiguity cannot also rescue it — the kernel
//     reaches `escapeOrRefuse` precisely because none of them was the
//     exact-one match — so only escape rows declaring the class close this
//     arm. Counting ordinary rows here would let the very overlap that mints
//     the refusal certify it rescued.
func (a *analysis) coverageUnionFor(g guard.Group, class string) guard.AssignmentSet {
	if class == ambiguousMatch {
		return a.escapeUnionFor(g, class)
	}
	if len(guard.Dimensions(a.model, g)) == 0 {
		// Over the empty product the union is decided by membership alone: a
		// row in the closing population denotes the single empty assignment,
		// and an empty population denotes nothing.
		for _, row := range g.Rows {
			if row.Kind() == table.KindEscape && !slices.Contains(row.Escape, class) {
				continue
			}
			return guard.Product(a.model, g)
		}
		return guard.AssignmentSet{}
	}
	return guard.CoverageUnionFor(a.model, g, class)
}

// escapeUnionFor unions the accepted assignments of the escape rows
// declaring class, and nothing else.
func (a *analysis) escapeUnionFor(g guard.Group, class string) guard.AssignmentSet {
	product := guard.Product(a.model, g)
	var union guard.AssignmentSet
	var seeded bool
	for _, row := range g.Rows {
		if row.Kind() != table.KindEscape || !slices.Contains(row.Escape, class) {
			continue
		}
		if len(guard.Dimensions(a.model, g)) == 0 {
			return product
		}
		accepted := guard.AcceptedAssignments(a.model, row)
		if !accepted.Projectable() {
			continue
		}
		if !seeded {
			union, seeded = accepted, true
			continue
		}
		union = union.Union(accepted)
	}
	if !seeded {
		return guard.AssignmentSet{}
	}
	return union
}

// ambiguousMatch names the rescuable class whose arm is demanded only where
// an overlap finding makes it reachable.
const ambiguousMatch = "ambiguous_match"

// groupHasOverlap reports whether this run already emitted an overlap
// finding for the group's ordinary population. The coverage arm reads the
// EMITTED verdict rather than recomputing it, so the two clauses can never
// disagree about whether the group overlaps.
func (a *analysis) groupHasOverlap(g guard.Group) bool {
	ids := ruleIDsOf(g)
	for _, f := range a.findings {
		if f.Code != CodeOverlap || f.Class != "" {
			continue
		}
		if slices.Contains(ids, authoredID(f.Rule)) && slices.Contains(ids, authoredID(f.Element)) {
			return true
		}
	}
	return false
}

// ordinaryClosesAlone reports whether the group's non-escape rows already
// close coverage, in which case no escape row closed anything.
func (a *analysis) ordinaryClosesAlone(g guard.Group, product guard.AssignmentSet) bool {
	var union guard.AssignmentSet
	var seeded bool
	for _, row := range g.Rows {
		if row.Kind() == table.KindEscape {
			continue
		}
		accepted := guard.AcceptedAssignments(a.model, row)
		if !accepted.Projectable() {
			continue
		}
		if !seeded {
			union, seeded = accepted, true
			continue
		}
		union = union.Union(accepted)
	}
	return seeded && union.Equal(product)
}

// bareEscapeFor names the escape row declaring class that carries NO guard
// atoms — the one denoting the whole scoped product, which therefore closes
// coverage by itself.
func bareEscapeFor(g guard.Group, class string) string {
	var found []string
	for _, row := range g.Rows {
		if row.Kind() != table.KindEscape || !slices.Contains(row.Escape, class) {
			continue
		}
		if len(guardAtomsOf(row)) == 0 {
			// The finding names this ROW, so an expanded one publishes its
			// suffixed identity (`0030:C3`).
			found = append(found, rowIdentity(row))
		}
	}
	slices.Sort(found)
	if len(found) == 0 {
		return ""
	}
	return found[0]
}

// ruleIDsOf names the group's rows, sorted and duplicate-free.
func ruleIDsOf(g guard.Group) []string {
	var out []string
	for _, row := range g.Rows {
		if !slices.Contains(out, row.RuleID) {
			out = append(out, row.RuleID)
		}
	}
	slices.Sort(out)
	return out
}

// firstRuleID names the group's lowest-sorting rule id, so a group-scoped
// finding still carries a source rule id rather than the code alone.
func firstRuleID(g guard.Group) string {
	ids := ruleIDsOf(g)
	if len(ids) == 0 {
		return ""
	}
	return ids[0]
}
