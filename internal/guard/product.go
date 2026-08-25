package guard

import (
	"slices"
	"strings"

	"github.com/newcoinc/intrastate/internal/resolve"
	"github.com/newcoinc/intrastate/internal/table"
)

// Selection is one selection context: the source state and the recognized
// outcome a row group shares. That is exactly the set RDR 0001 resolves
// exact-one over, which is why the row group is defined by it.
type Selection struct {
	Model   string
	Outcome string
	Match   []table.Atom
}

// String renders the selection context.
func (s Selection) String() string {
	var b strings.Builder
	b.WriteString(s.Model)
	b.WriteString("/")
	b.WriteString(s.Outcome)
	for _, a := range s.Match {
		b.WriteString(" ")
		b.WriteString(a.Key)
		b.WriteString(".")
		b.WriteString(a.Operator)
		b.WriteString("=")
		b.WriteString(strings.Join(a.Literal, ","))
	}
	return b.String()
}

// MatchKeys names the match keys that scope the context. They are NOT
// product dimensions: a row's match pattern selects which group the row
// belongs to, and is a separate field from its guard.
func (s Selection) MatchKeys() []string {
	out := make([]string, 0, len(s.Match))
	for _, a := range s.Match {
		if !slices.Contains(out, a.Key) {
			out = append(out, a.Key)
		}
	}
	slices.Sort(out)
	return out
}

// Group is one scoped row group: the normalized candidate rows sharing one
// selection context.
type Group struct {
	Context Selection
	Rows    []table.Row
}

// Groups partitions a model's rows into scoped row groups.
//
// The grouping predicate is self-contained: it reads the model's own rows
// and nothing else. RDR 0006 supplies the graph traversal enumerating which
// selection contexts are REACHABLE; it does not define the grouping, and
// this RDR reads none back from it. Every row therefore lands in exactly
// one group, which an externally-supplied reachability set could not
// guarantee.
func Groups(m *table.Model) []Group {
	index := map[string]int{}
	var out []Group
	for _, row := range m.Rows {
		ctx := selectionOf(m, row)
		key := ctx.String()
		i, ok := index[key]
		if !ok {
			i = len(out)
			index[key] = i
			out = append(out, Group{Context: ctx})
		}
		out[i].Rows = append(out[i].Rows, row)
	}
	slices.SortFunc(out, func(a, b Group) int {
		return strings.Compare(a.Context.String(), b.Context.String())
	})
	return out
}

// selectionOf builds a row's selection context: its recognized outcome plus
// its match pattern, which is the source state this RDR groups by.
func selectionOf(m *table.Model, row table.Row) Selection {
	var match []table.Atom
	for _, a := range row.Atoms {
		if a.Block == table.BlockMatch {
			match = append(match, a)
		}
	}
	slices.SortFunc(match, compareAtoms)
	return Selection{Model: m.ID, Outcome: row.Outcome, Match: match}
}

func compareAtoms(a, b table.Atom) int {
	if c := strings.Compare(a.Key, b.Key); c != 0 {
		return c
	}
	if c := strings.Compare(a.Operator, b.Operator); c != 0 {
		return c
	}
	return strings.Compare(strings.Join(a.Literal, "\x00"), strings.Join(b.Literal, "\x00"))
}

// guardAtoms returns a row's guard atoms — every atom under `guard.all` or
// `guard.unless`. The authored BLOCK is what makes an atom a guard atom,
// regardless of operator, `eq` included.
func guardAtoms(row table.Row) []table.Atom {
	var out []table.Atom
	for _, a := range row.Atoms {
		if a.Block == table.BlockAll || a.Block == table.BlockUnless {
			out = append(out, a)
		}
	}
	return out
}

// Dimensions names the participating guard value dimensions, sorted.
//
// A dimension participates when ANY row in the group carries a guard atom
// over that key, in either block — not only when the rows constrain it
// differently. A key every row constrains identically still bounds the
// product and MUST NOT be dropped for failing to discriminate. Match keys
// never enter: the split is per atom per group, so the same key may be a
// match key in one group and a guard key in another.
func Dimensions(m *table.Model, g Group) []string {
	var out []string
	for _, row := range g.Rows {
		for _, a := range guardAtoms(row) {
			if !slices.Contains(out, a.Key) {
				out = append(out, a.Key)
			}
		}
	}
	slices.Sort(out)
	return out
}

// PresenceDimensions names the keys contributing a presence dimension: the
// participating keys declared OPTIONAL. A key declared always-present
// contributes none — its `{absent}` assignment is not in the product, which
// is what makes the narrowing's negative control provable rather than
// vacuous.
func PresenceDimensions(m *table.Model, g Group) []string {
	var out []string
	for _, key := range Dimensions(m, g) {
		if !m.Tags[key].Required {
			out = append(out, key)
		}
	}
	return out
}

// productDims names every dimension of the group's scoped product: each
// participating key's value dimension, plus a presence dimension for each
// participating key declared optional.
func productDims(m *table.Model, g Group) []string {
	var out []string
	for _, key := range Dimensions(m, g) {
		out = append(out, key)
		if !m.Tags[key].Required {
			out = append(out, presenceDim(key))
		}
	}
	slices.Sort(out)
	return out
}

// Product returns the group's scoped product: the cross-product of every
// participating dimension. A product carrying an unprovable dimension has
// no members and is not projectable — the dimension has no assignment count
// to range over.
//
// A product larger than the published bound is likewise not projectable:
// enumerating it is exactly the silent cap this RDR forbids, so the
// representation declines rather than capping, and lint reports the
// over-large refusal instead.
func Product(m *table.Model, g Group) AssignmentSet {
	if card, ok := Cardinality(m, g); !ok || card > Bound() {
		return AssignmentSet{}
	}
	dims := productDims(m, g)
	ranges := make([][]string, 0, len(dims))
	for _, dim := range dims {
		key, isPresence := valueKeyOf(dim)
		if isPresence {
			ranges = append(ranges, []string{string(PresencePresent), string(PresenceAbsent)})
			continue
		}
		values, ok := valueAssignments(m.Tags[key])
		if !ok {
			return AssignmentSet{}
		}
		ranges = append(ranges, values)
	}

	out := newSet(dims)
	cross(dims, ranges, Assignment{}, 0, func(a Assignment) { out.add(a) })
	return out
}

// cross enumerates the cross-product of ranges, calling yield once per
// point. It is the deterministic enumeration the set-valued proof rests on:
// the same model gives the same points in the same order every run.
func cross(dims []string, ranges [][]string, acc Assignment, i int, yield func(Assignment)) {
	if i == len(dims) {
		yield(acc)
		return
	}
	for _, v := range ranges[i] {
		acc[dims[i]] = v
		cross(dims, ranges, acc, i+1, yield)
	}
}

// Cardinality returns the scoped product's cardinality, and whether every
// participating dimension is finite.
//
// The quantity is the number of ASSIGNMENTS in the product — the product of
// every participating dimension's assignment count — never a bitset width,
// byte size, or row count. A dimension with no finite declared domain has
// NO assignment count, so a product containing one has no cardinality at
// all, which is why the too-large comparison is defined only over a
// fully-provable product.
func Cardinality(m *table.Model, g Group) (int, bool) {
	card := 1
	for _, key := range Dimensions(m, g) {
		n, ok := AssignmentCount(m.Tags[key])
		if !ok {
			return 0, false
		}
		if card > cardinalityCeiling/max(n, 1) {
			// The product is finite but larger than this implementation
			// counts to. Saturating rather than overflowing keeps the
			// too-large comparison a comparison: a wrapped negative would
			// read as under the bound and certify the very product the
			// clause refuses.
			return cardinalityCeiling, true
		}
		card *= n
	}
	return card, true
}

// cardinalityCeiling is the largest cardinality this implementation
// reports. It is far above the published bound, so every product it
// saturates is already refused; it exists only so an astronomically large
// finite product reports a large number rather than an overflowed one.
const cardinalityCeiling = 1 << 40

// ProofCompletes reports whether the proof representation completes within
// the implementation's own resource budget.
//
// The budget IS the published bound: the representation this implementation
// proves with enumerates the scoped product, so it completes exactly when
// the product is small enough to enumerate, and the bound is chosen to be
// that threshold rather than discovered by exhausting memory or wall-clock.
// A15 is refuted if the two ever disagree, which is what makes stating them
// as one quantity — rather than two that happen to track — the claim.
func ProofCompletes(m *table.Model, g Group) bool {
	card, ok := Cardinality(m, g)
	return ok && card <= Bound()
}

// bound is the single integer cardinality bound this implementation
// publishes. It is a declared, model-independent constant: not per-model,
// per-group, or configurable per run, and never discovered by exhausting a
// resource.
//
// 2048 is the largest product this implementation's enumerating proof
// representation completes over within its budget. Another conforming
// implementation MAY publish a different bound; what both MUST do is return
// the same verdict for the same model and report which bound they applied.
const bound = 2048

// Bound returns the single integer cardinality bound this implementation
// publishes.
func Bound() int { return bound }

// --- denotation ----------------------------------------------------------

// Denotation returns the subset of key's dimensions the atom denotes.
//
// An `exists` atom projects onto the key's PRESENCE dimension, denoting
// `{present}` or `{absent}` on it, and reads presence alone — so it
// projects over a key whose values co-occur exactly as over a single-valued
// one. A value atom denotes a subset of the key's VALUE dimension as the
// operator/kind agreement clause projects it, and — because a value atom
// over an absent key is unevaluable rather than false — implicitly
// `{present}`.
func Denotation(m *table.Model, key string, atom table.Atom) AssignmentSet {
	decl := m.Tags[key]
	optional := !decl.Required

	dims := []string{key}
	if optional {
		dims = append(dims, presenceDim(key))
	}

	if atom.Operator == resolve.OpExists {
		// `exists` reads presence alone and always projects. Over an
		// always-present key there is no presence dimension to select on,
		// so it denotes the whole value dimension — well-formed but vacuous.
		values, ok := valueAssignments(decl)
		if !ok {
			// The value dimension is unprovable, but presence is still
			// two-valued and `exists` still reads it. Project over presence
			// alone rather than refusing: dropping the atom is the
			// false-green this RDR's narrowing forbids.
			if !optional {
				return AssignmentSet{}
			}
			out := newSet([]string{presenceDim(key)})
			out.add(Assignment{presenceDim(key): presenceFor(atom)})
			return out
		}
		out := newSet(dims)
		for _, v := range values {
			a := Assignment{key: v}
			if optional {
				a[presenceDim(key)] = presenceFor(atom)
			}
			out.add(a)
		}
		return out
	}

	// A value atom over a dimension the model does not declare
	// single-valued is one lint cannot project AT ALL, exactly as a
	// dimension with no finite declared domain is.
	values, ok := valueAssignments(decl)
	if !ok {
		return AssignmentSet{}
	}
	if isSingleValueOperator(atom.Operator) && !decl.SingleValued {
		return AssignmentSet{}
	}
	if atom.Operator == "contains" && decl.Kind != "set" {
		return AssignmentSet{}
	}

	out := newSet(dims)
	for _, v := range values {
		if !valueSatisfies(atom, v) {
			continue
		}
		a := Assignment{key: v}
		if optional {
			// Implicitly {present}: a value atom over an absent key is
			// unevaluable rather than false, so it selects no assignment in
			// which the key is absent.
			a[presenceDim(key)] = string(PresencePresent)
		}
		out.add(a)
	}
	return out
}

// presenceFor maps an `exists` literal onto the presence value it selects.
func presenceFor(atom table.Atom) string {
	if len(atom.Literal) == 1 && atom.Literal[0] == resolve.LiteralTrue {
		return string(PresencePresent)
	}
	return string(PresenceAbsent)
}

// valueSatisfies decides one value atom against one rendered held value.
//
// It is the same decision the runtime evaluator makes, over the same
// rendered value form, which is what keeps the lint claim about the runtime
// it describes rather than about a second semantics.
func valueSatisfies(atom table.Atom, held string) bool {
	seam := Evaluator{}
	literal := strings.Join(atom.Literal, "")
	if atom.Operator == "in" || atom.Operator == "contains" {
		literal = renderSet(atom.Literal)
	}
	return seam.Evaluate(resolve.GuardAtom{
		Key:      atom.Key,
		Operator: atom.Operator,
		Literal:  literal,
		Block:    atom.Block,
	}, held) == resolve.GuardTrue
}

// --- row acceptance ------------------------------------------------------

// AcceptedAssignments returns a row's accepted assignments: the
// intersection of all positive `all` atom domains, minus the SINGLE
// conjunctive assignment set matched by the row's full `unless` block.
//
// `unless` is not per-atom negation — subtracting every assignment matching
// EITHER unless atom would remove points the block as a whole does not
// match — and it creates no source-order priority.
//
// A row lint cannot project contributes nothing: a can-refuse row and a row
// carrying an unprojectable atom both return the unprojectable set, and
// lint MUST NOT credit either with its `all`-intersection unsubtracted.
func AcceptedAssignments(m *table.Model, row table.Row) AssignmentSet {
	return acceptedIn(m, groupContaining(m, row), row)
}

// groupContaining returns the scoped row group the row sits in.
func groupContaining(m *table.Model, row table.Row) Group {
	want := selectionOf(m, row).String()
	for _, g := range Groups(m) {
		if g.Context.String() == want {
			return g
		}
	}
	return Group{Context: selectionOf(m, row), Rows: []table.Row{row}}
}

// acceptedIn computes a row's accepted assignments within a known group, so
// the product it is expressed over is the group's rather than the row's.
func acceptedIn(m *table.Model, g Group, row table.Row) AssignmentSet {
	if CanRefuse(m.Tags, row) {
		// The subtraction model is defined over DECIDED atoms. A row that
		// can refuse has no decidable accepted-assignment set to
		// contribute, so it yields the unprojectable set rather than its
		// `all`-intersection unsubtracted — crediting it is the false-green
		// the narrowing forbids, and it is why the surviving overlap check
		// is scoped to the group's DECIDABLE rows.
		return AssignmentSet{}
	}

	product := Product(m, g)
	if !product.Projectable() {
		return AssignmentSet{}
	}

	accepted := product
	var unlessTerms []AssignmentSet
	for _, atom := range guardAtoms(row) {
		d := Denotation(m, atom.Key, atom)
		if !d.Projectable() {
			return AssignmentSet{}
		}
		if atom.Block == table.BlockUnless {
			unlessTerms = append(unlessTerms, d)
			continue
		}
		accepted = accepted.Intersect(d)
	}

	if len(unlessTerms) > 0 {
		excluded := product
		for _, term := range unlessTerms {
			excluded = excluded.Intersect(term)
		}
		accepted = accepted.subtract(excluded)
	}
	return accepted
}

// CanRefuse reports whether a row carries a value atom over a key declared
// optional, in either guard block.
//
// The test is SYNTACTIC over declarations so it is total: every row gets an
// answer, including one carrying no atoms and one naming an undeclared key.
// It reads exactly one declared field — the optionality marker — and no
// second presence property and no graph query. An owned tag's graph-level
// presence condition governs whether a row may MATCH, not whether its guard
// can refuse, so provenance never changes the answer.
func CanRefuse(decls map[string]table.TagDecl, row table.Row) bool {
	for _, atom := range guardAtoms(row) {
		if atom.Operator == resolve.OpExists {
			// An existence atom over an absent key is DECIDED from presence
			// alone, never unevaluable, so it cannot make a row refuse.
			continue
		}
		decl, ok := decls[atom.Key]
		if !ok {
			continue
		}
		if !decl.Required {
			return true
		}
	}
	return false
}

// refusingAtom returns the first value atom over an optional key a row
// carries, in a deterministic order, so the withholding finding can name
// the atom that can refuse.
func refusingAtom(decls map[string]table.TagDecl, row table.Row) *table.Atom {
	atoms := guardAtoms(row)
	slices.SortFunc(atoms, compareAtoms)
	for i := range atoms {
		if atoms[i].Operator == resolve.OpExists {
			continue
		}
		if decl, ok := decls[atoms[i].Key]; ok && !decl.Required {
			return &atoms[i]
		}
	}
	return nil
}

// RescuableClasses names the failure classes an escape row may rescue: of
// the classes it declares, only `no_match` and `ambiguous_match`. The
// kernel consults escape rows solely from the zero-match and multiple-match
// arms, so an escape row cannot rescue either blocking class however bare
// its guard.
func RescuableClasses() []string {
	return []string{string(resolve.KindNoMatch), string(resolve.KindAmbiguousMatch)}
}

// CoverageUnion returns the union of the group's rows' accepted
// assignments. Every row participates, escape rows included: excluding them
// from COVERAGE is what MUST NOT happen.
func CoverageUnion(m *table.Model, g Group) AssignmentSet {
	return coverageUnion(m, g, "")
}

// CoverageUnionFor returns the union for one declared rescuable class. The
// union is computed per (scoped row group × declared rescuable class), so a
// row declaring one class does not close the group's other arm.
func CoverageUnionFor(m *table.Model, g Group, class string) AssignmentSet {
	return coverageUnion(m, g, class)
}

// coverageUnion unions the accepted assignments of the group's rows. An
// empty class unions every row; a named class unions the ordinary rows plus
// only the escape rows declaring that class.
func coverageUnion(m *table.Model, g Group, class string) AssignmentSet {
	out := newSet(productDims(m, g))
	if !Product(m, g).Projectable() {
		return AssignmentSet{}
	}
	for _, row := range g.Rows {
		if class != "" && row.Kind() == table.KindEscape && !slices.Contains(row.Escape, class) {
			continue
		}
		accepted := acceptedIn(m, g, row)
		if !accepted.Projectable() {
			// A can-refuse or unprojectable row contributes NO assignments.
			continue
		}
		out = out.Union(accepted)
	}
	return out
}
