package graphlint

import (
	"fmt"
	"slices"
	"strings"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
	"github.com/newcoinc/intrastate/internal/guard"
	"github.com/newcoinc/intrastate/internal/table"
)

// checkGroups runs the per-group invariants: determinism/overlap
// (invariant 3), guard exhaustiveness and the withholding rules
// (invariant 4), owned-set-before-match (invariant 6), and the redundant-row
// and vacuous-atom advisories.
//
// Invariants 3 and 4 run once per REACHABLE group. Membership is syntactic
// and per-row — the authored match pattern — so the group count is bounded
// by the authored rule count regardless of the lattice; only the
// reachability filter scales with nodes.
//
// Every check runs for every group; none short-circuits another. Withholding
// a group's exhaustiveness claim does not suppress overlap, coverage, or
// further withholding findings for that same group.
func (a *analysis) checkGroups() {
	for _, g := range a.groups {
		// Invariant 6 is per row and reads the reachability nodes, not the
		// group's product, so it runs for every group whether or not the
		// group's context is reachable.
		a.checkOwnedBeforeMatch(g)
		a.checkVacuousAtoms(g)

		if !a.reachable[g.Context.String()] {
			// Reachability is a filter over contexts: an unreachable group
			// is not proven at all, and its rows take the unreachable-rule
			// advisory instead. Its rows are never pooled into a reachable
			// group — membership is untouched by the filter.
			continue
		}

		a.checkOverlap(g)
		a.checkRedundantRows(g)
		a.checkCoverage(g)
	}
}

// --- invariant 6: owned-set-before-match ---------------------------------

// checkOwnedBeforeMatch checks that a row reading an owned tag finds it held
// in every reachable owned-state that satisfies the row's match pattern. The
// declared initial owned state counts as a write.
//
// The read-set is computed HERE from the row's atoms, never read from
// `Row.RequiresOwned`: that field is the sorted key set the write block and
// clear list name, so reading it would check writes-before-writes instead of
// reads-before-writes and silently never fire on the defect this invariant
// exists for. The read-set is the owned-provenance keys named by the row's
// match pattern and its `all`/`unless` guard atoms.
//
// The test is evaluated against MERGED fixpoint nodes. It is existential in
// the direction that matters — some reachable node fails to hold the key —
// so the widened domain yields false positives at worst, never a missed
// defect. Those false positives are accepted: their cure is an explicit
// write, clear, or terminal declaration on the model, never a guard-aware
// lint.
func (a *analysis) checkOwnedBeforeMatch(g guard.Group) {
	rows := slices.Clone(g.Rows)
	slices.SortFunc(rows, func(x, y table.Row) int { return strings.Compare(x.RuleID, y.RuleID) })

	for _, row := range rows {
		for _, key := range a.ownedReadSet(row) {
			node, missing := a.nodeMissingKey(row, key)
			if !missing {
				continue
			}
			a.emit(clierr.Finding{
				Code:        CodeOwnedBeforeWrite,
				Rule:        row.RuleID,
				Span:        row.SourceLocator,
				Key:         key,
				Dimension:   key,
				Element:     nodeElement(node),
				Fingerprint: Fingerprint(row),
				Message: fmt.Sprintf("row %q reads the owned tag %q, but the "+
					"reachable owned-state %s that satisfies its match pattern "+
					"does not hold it; add an explicit write or establish it "+
					"in `[initial]`", row.RuleID, key, nodeElement(node)),
			})
		}
	}
}

// ownedReadSet derives a row's owned read-set: the owned-provenance keys its
// match pattern and its `all`/`unless` guard atoms name, sorted and
// duplicate-free.
//
// An `exists` atom is excluded: it reads PRESENCE, so absence is a verdict
// it decides rather than a value it fails to find, and demanding the key be
// held would make the atom's own negative arm a defect.
func (a *analysis) ownedReadSet(row table.Row) []string {
	var out []string
	for _, atom := range row.Atoms {
		if atom.Operator == "exists" {
			continue
		}
		if a.model.Tags[atom.Key].Provenance != table.ProvenanceOwned {
			continue
		}
		if !slices.Contains(out, atom.Key) {
			out = append(out, atom.Key)
		}
	}
	slices.Sort(out)
	return out
}

// nodeMissingKey returns the first reachable node satisfying the row's match
// pattern that does not hold key, in traversal order.
func (a *analysis) nodeMissingKey(row table.Row, key string) (Node, bool) {
	for _, n := range a.nodes {
		if !matchSatisfiable(a.model, n, row) {
			continue
		}
		if len(n.Values[key]) == 0 {
			return n, true
		}
	}
	return Node{}, false
}

// --- invariant 3: determinism / overlap ----------------------------------

// checkOverlap reports each pair of rows in the group whose accepted
// assignments intersect, in TWO separate populations because the runtime
// never mixes them.
//
// The ordinary population excludes escape rows: the kernel matches an escape
// row only to rescue a `no_match` / `ambiguous_match` refusal and only when
// exactly one matches, so an escape row overlapping an ordinary row is not a
// runtime ambiguity and MUST NOT be reported. The escape population is
// partitioned per DECLARED failure class, with a row declaring several
// classes placed in EACH of those populations — so a pair sharing two
// classes yields one finding per shared class.
//
// A pair one of whose accepted sets is a PROPER SUBSET of the other's is not
// reported here: that is the advisory redundant-row case, which is distinct
// from overlap by the taxonomy's own definition ("overlap … is a partial
// intersection between two rows neither of which subsumes the other").
func (a *analysis) checkOverlap(g guard.Group) {
	var ordinary, escape []table.Row
	for _, row := range g.Rows {
		if row.Kind() == table.KindEscape {
			escape = append(escape, row)
			continue
		}
		ordinary = append(ordinary, row)
	}

	a.emitOverlaps(g, ordinary, "")
	for _, class := range guard.RescuableClasses() {
		var population []table.Row
		for _, row := range escape {
			if slices.Contains(row.Escape, class) {
				population = append(population, row)
			}
		}
		a.emitOverlaps(g, population, class)
	}
}

// emitOverlaps reports each overlapping pair within one population. ONE pair
// is ONE finding naming BOTH rows, never one per row.
func (a *analysis) emitOverlaps(g guard.Group, rows []table.Row, class string) {
	rows = slices.Clone(rows)
	slices.SortFunc(rows, func(x, y table.Row) int { return strings.Compare(x.RuleID, y.RuleID) })

	for i := range rows {
		left := guard.AcceptedAssignments(a.model, rows[i])
		if !left.Projectable() {
			continue
		}
		for _, right := range rows[i+1:] {
			other := guard.AcceptedAssignments(a.model, right)
			if !other.Projectable() || left.Intersect(other).Len() == 0 {
				continue
			}
			if properSubset(left, other) || properSubset(other, left) {
				// Subsumption is the redundant-row advisory, not overlap.
				continue
			}
			a.emit(clierr.Finding{
				Code:        CodeOverlap,
				Rule:        rows[i].RuleID,
				Span:        rows[i].SourceLocator,
				Element:     right.RuleID,
				Class:       class,
				Fingerprint: Fingerprint(rows[i]) + "|" + Fingerprint(right),
				Message: fmt.Sprintf("rows %q and %q are enabled together by "+
					"some assignment in %s%s; lint rejects the ambiguity "+
					"rather than selecting by source order",
					rows[i].RuleID, right.RuleID, g.Context.String(),
					classSuffix(class)),
			})
		}
	}
}

// classSuffix renders the failure-class scope for an escape-population
// message, and nothing at all for the ordinary population.
func classSuffix(class string) string {
	if class == "" {
		return ""
	}
	return " (escape population " + class + ")"
}

// properSubset reports whether left's members are a strict subset of
// right's.
func properSubset(left, right guard.AssignmentSet) bool {
	return left.Subset(right) && left.Len() < right.Len()
}

// --- the redundant-row advisory ------------------------------------------

// checkRedundantRows reports a row whose accepted assignments are a PROPER
// SUBSET of a sibling's in the same group, so it can never be the exact-one
// match. It is advisory and never changes the success disposition.
func (a *analysis) checkRedundantRows(g guard.Group) {
	rows := slices.Clone(g.Rows)
	slices.SortFunc(rows, func(x, y table.Row) int { return strings.Compare(x.RuleID, y.RuleID) })

	for i := range rows {
		mine := guard.AcceptedAssignments(a.model, rows[i])
		if !mine.Projectable() {
			continue
		}
		for j := range rows {
			if i == j || rows[i].Kind() != rows[j].Kind() {
				continue
			}
			sibling := guard.AcceptedAssignments(a.model, rows[j])
			if !sibling.Projectable() || !properSubset(mine, sibling) {
				continue
			}
			a.emit(clierr.Finding{
				Code:        CodeRedundantRow,
				Rule:        rows[i].RuleID,
				Span:        rows[i].SourceLocator,
				Element:     rows[j].RuleID,
				Fingerprint: Fingerprint(rows[i]),
				Message: fmt.Sprintf("row %q accepts a proper subset of row "+
					"%q's assignments, so it can never be the exact-one match",
					rows[i].RuleID, rows[j].RuleID),
			})
			break
		}
	}
}

// --- the vacuous-atom advisory -------------------------------------------

// checkVacuousAtoms reports an `exists` atom over a key declared
// always-present. Such a key contributes no `{absent}` assignment, so the
// atom is well-formed but VACUOUS — RDR 0003 requires it be reported as
// such rather than rejected, which is why it is advisory.
func (a *analysis) checkVacuousAtoms(g guard.Group) {
	rows := slices.Clone(g.Rows)
	slices.SortFunc(rows, func(x, y table.Row) int { return strings.Compare(x.RuleID, y.RuleID) })

	for _, row := range rows {
		atoms := guardAtomsOf(row)
		slices.SortFunc(atoms, compareAtoms)
		for _, atom := range atoms {
			if atom.Operator != "exists" || !a.model.Tags[atom.Key].Required {
				continue
			}
			a.emit(clierr.Finding{
				Code:        CodeVacuousAtom,
				Rule:        row.RuleID,
				Span:        row.SourceLocator,
				Key:         atom.Key,
				Operator:    atom.Operator,
				Literal:     strings.Join(atom.Literal, ","),
				Block:       string(atom.Block),
				Dimension:   atom.Key,
				Fingerprint: Fingerprint(row),
				Message: fmt.Sprintf("row %q carries an `exists` atom over the "+
					"always-present key %q; the atom is well-formed but "+
					"vacuous", row.RuleID, atom.Key),
			})
		}
	}
}

// guardAtomsOf returns a row's guard atoms — every atom authored under
// `guard.all` or `guard.unless`. The authored BLOCK is what makes an atom a
// guard atom, regardless of operator, `eq` included.
func guardAtomsOf(row table.Row) []table.Atom {
	var out []table.Atom
	for _, atom := range row.Atoms {
		if atom.Block == table.BlockAll || atom.Block == table.BlockUnless {
			out = append(out, atom)
		}
	}
	return out
}
