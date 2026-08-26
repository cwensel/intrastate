package resolve

// RDR 0008 — the package-internal path.
//
// Blocks 5 and 6 keep two behaviors alive as DOCUMENTED RESIDUALS reachable
// only by a caller that reaches `assemble` without passing the entry
// precondition (`0008:C5`, `0008:C6`): the `owned_state_unavailable` refusal
// for a row naming the reserved key in `RequiresOwned` (scenario 9's first
// half), and RDR 0001 D3's precedence resolving a reserved-key collision
// between owned and observed tags (scenario 8).
//
// The residuals are unreachable through `Resolve` once the precondition
// lands, so the tests that pin them need this seam. It exports nothing
// outside the test binary and adds no production surface — it exists
// precisely so that "documented residual" means something executable rather
// than something asserted.

// ResolveBypassingPreconditionForTest runs the kernel's disposition pipeline
// starting at `assemble`, skipping whatever entry precondition `Resolve`
// applies. It is the package-internal path REQ-51/REQ-96 name.
func ResolveBypassingPreconditionForTest(in Input) Result {
	view := assemble(in)

	if !in.Table.models(in.Recognized) {
		return refuse(in, Refusal{Kind: KindUnmodeledOutcome})
	}

	var candidates []Row
	for _, row := range in.Table.Rows {
		if len(row.Escape) != 0 || row.Outcome != in.Recognized {
			continue
		}
		if !view.matches(row.Match) {
			continue
		}
		candidates = append(candidates, row)
	}

	selected, blocked := gate(candidates, in.Guards, view)
	if blocked != nil {
		return refuse(in, *blocked)
	}
	switch len(selected) {
	case 1:
		return Result{Plan: planOf(in, selected[0], false)}
	case 0:
		return escapeOrRefuse(in, view, Refusal{Kind: KindNoMatch})
	default:
		return escapeOrRefuse(in, view, Refusal{
			Kind: KindAmbiguousMatch,
			Rows: rowRefs(selected),
		})
	}
}

// AssembledBindingForTest reports what the assembled view binds at key,
// bypassing the entry precondition. Scenario 8 needs it to observe D3's
// precedence over a collision no conforming caller can construct.
func AssembledBindingForTest(in Input, key string) (value string, prov Provenance, ok bool) {
	return assemble(in).Lookup(key)
}

// AssembledLenForTest reports the assembled view's distinct key count on the
// package-internal path.
func AssembledLenForTest(in Input) int { return assemble(in).Len() }
