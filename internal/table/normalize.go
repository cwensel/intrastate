package table

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/newcoinc/intrastate/internal/resolve"
)

// atomsFromBlock validates and renders one authored atom block —
// `[<owner>.match.<tag>]`, `[rule.guard.all.<tag>]`, or
// `[rule.guard.unless.<tag>]` — into normalized atoms.
//
// Atom-level validation is BLOCK-AGNOSTIC: operator membership, the
// `<clear>` reserved value, domain and kind conformance, and the tag-key
// declaration rule are enforced identically in all three blocks. Two rules
// are legitimately match-only and are NOT block-agnostic — the `eq`/`in`
// operator restriction and the `#` reservation in `in` members, since only
// match blocks expand (`0002:C17`).
func (l *loader) atomsFromBlock(block map[string]map[string]any, b Block, owner string) ([]Atom, error) {
	var out []Atom
	for _, key := range slices.Sorted(maps.Keys(block)) {
		decl, ok := l.model.Tags[key]
		if !ok {
			return nil, fail(CatUnknownTag, owner+" references the undeclared tag "+key)
		}
		predicates := block[key]
		for _, op := range slices.Sorted(maps.Keys(predicates)) {
			atom, err := l.atom(decl, key, op, predicates[op], b, owner)
			if err != nil {
				return nil, err
			}
			out = append(out, atom)
		}
	}
	return out, nil
}

// atom validates one authored predicate and renders it.
func (l *loader) atom(decl TagDecl, key, operator string, raw any, b Block, owner string) (Atom, error) {
	where := fmt.Sprintf("%s: %s.%s", owner, key, operator)
	badAtom := func(detail string) (Atom, error) {
		return Atom{}, fail(CatMalformedPredicateAtom, where+": "+detail)
	}

	// The admitted operator set, guard blocks included, is RDR 0003's
	// closed eight. This RDR mints none and widens it nowhere (`0002:C16`).
	if !slices.Contains(operators, operator) {
		return badAtom("unknown operator")
	}
	// A match block admits only `eq` and `in`; `in` expands into per-member
	// `eq` rows. Any other operator under a match block is refused
	// (`0002:C16`).
	if b == BlockMatch && operator != "eq" && operator != "in" {
		return badAtom("a match block admits only eq and in")
	}

	members, err := valueMembers(raw)
	if err != nil {
		return badAtom(err.Error())
	}

	// The `<clear>` ban binds every atom regardless of operator or block
	// (`0002:C17`), and outranks kind conformance so a sentinel authored on
	// a bool or int tag still reports as the reserved value it is.
	for _, m := range members {
		if m == ClearSentinel {
			return Atom{}, fail(CatReservedTagValue, where+" authors the reserved value "+ClearSentinel)
		}
	}

	switch operator {
	case "in":
		if !isArray(raw) {
			return badAtom("`in` takes an array of members")
		}
		// A repeated `in` element is refused here. `0002:C13` rests the
		// identity tuple's totality on it — "Because RDR 0003 rejects a
		// repeated `in` element at parse, no two rows of one rule share a
		// suffix, which is what keeps the identity tuple total" — and RDR
		// 0003's parser is not what loads this document, so `0002:C22`
		// lands the rejection in this package's `malformed predicate atom`
		// category. Deduplicating silently would mint the row set the
		// author meant while leaving the authoring error unreported; on an
		// escape rule the duplicate is a self-inflicted `ambiguous_match`
		// (`0002:C5`).
		if dup, ok := firstDuplicate(members); ok {
			return badAtom("member " + strconv.Quote(dup) + " is repeated")
		}
		if b == BlockMatch {
			// The `#` reservation is match-only: only match blocks expand,
			// and a `#` inside a chosen member would make the rendered row
			// identity ambiguous (`0002:C11`, `0002:C17`).
			for _, m := range members {
				if strings.Contains(m, suffixSep) {
					return badAtom("member " + strconv.Quote(m) +
						" contains the expansion-suffix separator " + suffixSep)
				}
			}
		}
	case "contains":
		if !isArray(raw) {
			return badAtom("`contains` takes an array of members")
		}
	case resolve.OpExists:
		// An existence atom emits the kernel's exported constants verbatim,
		// and any literal that is not one of the two boolean forms is
		// refused at load (`0002:C8`).
		if !isBool(raw) {
			return badAtom("an existence literal is a boolean")
		}
	case "lt", "lte", "gt", "gte":
		if isArray(raw) {
			return badAtom("a comparison bound is one value")
		}
	}

	if err := conform(decl, operator, members); err != nil {
		return badAtom(err.Error())
	}

	if operator == resolve.OpExists {
		// Verbatim: OpExists plus LiteralTrue or LiteralFalse (`0002:C8`).
		literal := resolve.LiteralFalse
		if members[0] == resolve.LiteralTrue {
			literal = resolve.LiteralTrue
		}
		return Atom{Key: key, Block: b, Operator: resolve.OpExists, Literal: []string{literal}}, nil
	}

	// A set-valued literal normalizes to an ORDERED member sequence, each
	// member compared byte-exactly, and members sort byte-lexicographically
	// so two authored orderings of one set are one literal (`0002:C10`).
	// This binds `in` too: its literal is a member SET, and a repeated
	// element was refused above, so the sort needs no dedup pass to be
	// faithful. Row identity is unaffected — a match-block `in` still
	// yields one row per member and rows re-sort by suffix (`0002:C19`).
	// A single-value operator keeps its one member as authored.
	if isArray(raw) {
		members = slices.Compact(slices.Sorted(slices.Values(members)))
	}
	return Atom{Key: key, Block: b, Operator: operator, Literal: members}, nil
}

// firstDuplicate reports the first member that repeats an earlier one.
func firstDuplicate(members []string) (string, bool) {
	seen := make(map[string]bool, len(members))
	for _, m := range members {
		if seen[m] {
			return m, true
		}
		seen[m] = true
	}
	return "", false
}

// atomIdentity encodes the full atom identity tuple (key, block, operator
// token, literal) injectively.
//
// The literal enters as its member SEQUENCE, each member length-prefixed —
// never as a joined string. Any delimiter is authorable inside a tag value,
// so a joined rendering would collapse ["a,b","c"] and ["a","b,c"] into one
// atom and turn a dead rule into a live one (`0002:C10`). Length prefixing
// makes the encoding injective for any member content.
func atomIdentity(a Atom) string {
	var b strings.Builder
	b.WriteString(a.Key)
	b.WriteString("\x00")
	b.WriteString(string(a.Block))
	b.WriteString("\x00")
	b.WriteString(a.Operator)
	b.WriteString("\x00")
	writeMemberKey(&b, a.Literal)
	return b.String()
}

// writeMemberKey encodes a member sequence injectively: each member is
// prefixed by its byte length, so no member content can forge a boundary.
func writeMemberKey(b *strings.Builder, members []string) {
	for _, m := range members {
		fmt.Fprintf(b, "%d:%s", len(m), m)
	}
}

// mergeAtoms is the set union over the FULL atom identity.
//
// Merging never keys on a proper prefix of the tuple: two atoms agreeing on
// (key, block, operator) but differing in literal are distinct atoms and
// both survive. It is idempotent on identical atoms — the same atom
// contributed by a rule and by one or more inherited contexts collapses to
// one — so inheritance never overrides, it only accumulates (`0002:C6`).
func mergeAtoms(sets ...[]Atom) []Atom {
	seen := map[string]bool{}
	var out []Atom
	for _, set := range sets {
		for _, a := range set {
			id := atomIdentity(a)
			if seen[id] {
				continue
			}
			seen[id] = true
			out = append(out, a)
		}
	}
	slices.SortFunc(out, compareAtoms)
	return out
}

// compareAtoms is the within-row atom order: (key, block, operator token,
// literal), the literal compared as a member sequence element by element
// (`0002:C19`).
func compareAtoms(a, b Atom) int {
	if c := strings.Compare(a.Key, b.Key); c != 0 {
		return c
	}
	if c := strings.Compare(string(a.Block), string(b.Block)); c != 0 {
		return c
	}
	if c := strings.Compare(a.Operator, b.Operator); c != 0 {
		return c
	}
	return slices.Compare(a.Literal, b.Literal)
}

// ------------------------------------------------------------ rule shape

// normalizeRules validates each rule and expands it into candidate rows.
func (l *loader) normalizeRules() error {
	seenID := map[string]bool{}
	var rows []Row

	setKeys := l.setValuedKeys()

	for i := range l.doc.Rule {
		rule := &l.doc.Rule[i]
		if rule.ID == nil || *rule.ID == "" {
			return fail(CatMalformedRuleShape, "a [[rule]] carries no id")
		}
		id := *rule.ID
		// Rule ids MUST NOT contain the expansion-suffix separator, or the
		// rendered identity of rule `a#x` and rule `a` expanding on member
		// `x` would collide (`0002:C11`).
		if strings.Contains(id, suffixSep) {
			return fail(CatMalformedRuleID,
				"rule id "+id+" contains the expansion-suffix separator "+suffixSep)
		}
		// Compared by exact byte equality: two ids differing only in case
		// are two ids (`0002:C4`).
		if seenID[id] {
			return fail(CatDuplicateRuleID, "rule id "+id+" is declared more than once")
		}
		seenID[id] = true

		expanded, err := l.normalizeRule(rule, id, setKeys)
		if err != nil {
			return err
		}
		rows = append(rows, expanded...)
	}

	// The dump is emitted from a PRE-SORTED sequence, never by iterating a
	// map and never by delegating key order to an encoder (`0002:C19`).
	slices.SortFunc(rows, compareRows)
	l.model.Rows = rows
	return nil
}

// setValuedKeys is the sorted set of tag keys whose declared kind is `set`.
func (l *loader) setValuedKeys() []string {
	var out []string
	for key, decl := range l.model.Tags {
		if decl.Kind == "set" {
			out = append(out, key)
		}
	}
	slices.Sort(out)
	return out
}

// normalizeRule validates one rule and returns its candidate rows.
func (l *loader) normalizeRule(rule *sourceRule, id string, setKeys []string) ([]Row, error) {
	isEscape := rule.Escape != nil

	// An escape rule MUST contain an `escape` list and MUST NOT contain a
	// write block, clear list, or gate list, EVEN AN EMPTY ONE: the loader
	// keys on key PRESENCE, not on length (`0002:C4`, deviations.md D2).
	if isEscape {
		switch {
		case rule.Write != nil:
			return nil, fail(CatMalformedEscapeDeclaration, "escape rule "+id+" carries a write block")
		case rule.Clear != nil:
			return nil, fail(CatMalformedEscapeDeclaration, "escape rule "+id+" carries a clear list")
		case rule.Gate != nil:
			return nil, fail(CatMalformedEscapeDeclaration, "escape rule "+id+" carries a gate list")
		}
		// The list admits only the two resolver failure classes RDR 0001
		// allows the table to model (`0002:C5`).
		for _, class := range *rule.Escape {
			switch resolve.RefusalKind(class) {
			case resolve.KindNoMatch, resolve.KindAmbiguousMatch:
			default:
				return nil, fail(CatMalformedEscapeDeclaration,
					"escape rule "+id+" models the unmodelable class "+class)
			}
		}
		if len(*rule.Escape) == 0 {
			return nil, fail(CatMalformedEscapeDeclaration, "escape rule "+id+" models no failure class")
		}
	} else if rule.Write == nil {
		// An ordinary transition rule MUST contain a write block
		// (`0002:C4`).
		return nil, fail(CatMalformedRuleShape, "ordinary rule "+id+" carries no write block")
	}

	// Shared-context references must resolve; inheritance normalizes to an
	// explicit predicate set before anything downstream reads it.
	var inherited []Atom
	for _, ref := range rule.Use {
		if _, ok := l.doc.Context[ref]; !ok {
			return nil, fail(CatUnknownContext, "rule "+id+" uses the unknown context "+ref)
		}
		atoms, err := l.contextAtoms(ref, nil)
		if err != nil {
			return nil, err
		}
		inherited = append(inherited, atoms...)
	}

	local, err := l.atomsFromBlock(rule.Match, BlockMatch, "rule "+id+" match")
	if err != nil {
		return nil, err
	}

	var guardAll, guardUnless []Atom
	if rule.Guard != nil {
		if guardAll, err = l.atomsFromBlock(rule.Guard.All, BlockAll, "rule "+id+" guard.all"); err != nil {
			return nil, err
		}
		if guardUnless, err = l.atomsFromBlock(rule.Guard.Unless, BlockUnless, "rule "+id+" guard.unless"); err != nil {
			return nil, err
		}
	}

	// A `recognized` atom authored under guard.all or guard.unless is
	// refused as a malformed outcome binding — never lifted (`0002:C13`).
	for _, a := range slices.Concat(guardAll, guardUnless) {
		if a.Key == RecognizedTagKey {
			return nil, fail(CatMalformedOutcomeBinding,
				"rule "+id+" authors a "+RecognizedTagKey+" atom under guard."+string(a.Block))
		}
	}

	// Outcome binding reads the MATCH blocks only — the rule's local match
	// plus the match blocks of its inherited contexts (`0002:C13`).
	matchAtoms := mergeAtoms(inherited, local)
	outcomeAtom, err := l.outcomeBinding(id, matchAtoms)
	if err != nil {
		return nil, err
	}

	// Normalization lifts that atom OUT of the predicate set into the row's
	// outcome field.
	predicates := make([]Atom, 0, len(matchAtoms))
	for _, a := range matchAtoms {
		if a.Key != RecognizedTagKey {
			predicates = append(predicates, a)
		}
	}
	predicates = mergeAtoms(predicates, guardAll, guardUnless)

	writes, requiresOwned, err := l.renderWrites(rule, id, isEscape)
	if err != nil {
		return nil, err
	}
	if rule.Gate != nil {
		for _, gid := range *rule.Gate {
			if _, ok := l.model.Gates[gid]; !ok {
				return nil, fail(CatUnknownAccessor,
					"rule "+id+" gates on the unknown gate accessor "+gid)
			}
		}
	}

	base := Row{
		ModelID: l.model.ID,
		RuleID:  id,
		// The locator is DERIVED from (model id, rule id), never from the
		// authored `source` annotation (`0002:TD`).
		SourceLocator: l.model.ID + ":" + id,
		Source:        rule.Source,
		Gate:          sortedList(rule.Gate),
		RequiresOwned: requiresOwned,
		Escape:        sortedList(rule.Escape),
		setKeys:       setKeys,
	}

	return expand(base, predicates, outcomeAtom, writes), nil
}

// outcomeBinding finds the single `recognized` atom the rule's match blocks
// carry. Every rule — ordinary or escape — MUST bind exactly one outcome,
// using `eq` or `in`, whose literals are alphabet members (`0002:C13`).
func (l *loader) outcomeBinding(id string, matchAtoms []Atom) (Atom, error) {
	var found []Atom
	for _, a := range matchAtoms {
		if a.Key == RecognizedTagKey {
			found = append(found, a)
		}
	}
	switch len(found) {
	case 1:
	case 0:
		return Atom{}, fail(CatMalformedOutcomeBinding, "rule "+id+" binds no outcome")
	default:
		return Atom{}, fail(CatMalformedOutcomeBinding,
			fmt.Sprintf("rule %s carries %d `%s` atoms; want exactly one",
				id, len(found), RecognizedTagKey))
	}

	atom := found[0]
	if atom.Operator != "eq" && atom.Operator != "in" {
		return Atom{}, fail(CatMalformedOutcomeBinding,
			"rule "+id+" binds its outcome with "+atom.Operator+"; want eq or in")
	}
	for _, member := range atom.Literal {
		if !slices.Contains(l.model.Outcomes, member) {
			return Atom{}, fail(CatMalformedOutcomeBinding,
				"rule "+id+" binds "+strconv.Quote(member)+", outside the declared alphabet")
		}
	}
	return atom, nil
}

// renderWrites renders a rule's write block and clear list into the
// next-state tags and the writes, and derives the required-owned key set.
//
// An escape rule carries neither, so an escape row normalizes to both empty
// and an empty required-owned set (`0002:C14`, `0002:C15`).
func (l *loader) renderWrites(rule *sourceRule, id string, isEscape bool) ([]TagValue, []string, error) {
	if isEscape {
		return nil, nil, nil
	}

	assignments := map[string][]string{}
	for _, key := range slices.Sorted(maps.Keys(*rule.Write)) {
		decl, ok := l.model.Tags[key]
		if !ok {
			return nil, nil, fail(CatUnknownTag, "rule "+id+" writes the undeclared tag "+key)
		}
		raw := (*rule.Write)[key]
		members, err := valueMembers(raw)
		if err != nil {
			return nil, nil, fail(CatMalformedTagDeclaration, "rule "+id+" write "+key+": "+err.Error())
		}
		for _, m := range members {
			if m == ClearSentinel {
				return nil, nil, fail(CatReservedTagValue,
					"rule "+id+" write "+key+" authors the reserved value "+ClearSentinel)
			}
		}
		// A write-block value's kind and domain conformance is filed under
		// the declaration category (deviations.md D3): `0002:C24` names no
		// dedicated category, and D3's check names this one.
		if err := conform(decl, "eq", members); err != nil {
			return nil, nil, fail(CatMalformedTagDeclaration, "rule "+id+" write "+key+": "+err.Error())
		}
		// A write REPLACES: for a set kind the array literal is the whole
		// new set, and it normalizes to a member-sorted sequence rather
		// than a delimiter-joined string (`0002:C4`).
		if decl.Kind == "set" {
			members = slices.Compact(slices.Sorted(slices.Values(members)))
		}
		assignments[key] = members
	}

	// Clearing a tag is represented by an explicit rule-level `clear`
	// entry that normalization renders as a `<clear>` write. Absence from
	// both the write block and the clear list never implies deletion
	// (`0002:C23`).
	if rule.Clear != nil {
		for _, key := range *rule.Clear {
			if _, ok := l.model.Tags[key]; !ok {
				return nil, nil, fail(CatUnknownTag, "rule "+id+" clears the undeclared tag "+key)
			}
			assignments[key] = []string{ClearSentinel}
		}
	}

	keys := slices.Sorted(maps.Keys(assignments))
	// NextTags and Writes are populated INDEPENDENTLY and neither is an
	// alias of the other: the two slices are built separately so a mutation
	// of one is not observable through the other (`0002:C15`).
	next := make([]TagValue, 0, len(keys))
	writes := make([]TagValue, 0, len(keys))
	for _, key := range keys {
		next = append(next, TagValue{Key: key, Value: slices.Clone(assignments[key])})
		writes = append(writes, TagValue{Key: key, Value: slices.Clone(assignments[key])})
	}
	return writes, slices.Clone(keys), nil
}

// sortedList renders an optional authored list as a sorted slice, so map
// and authoring order never reach the normalized value.
func sortedList(in *[]string) []string {
	if in == nil || len(*in) == 0 {
		return nil
	}
	out := slices.Clone(*in)
	slices.Sort(out)
	return out
}

// ------------------------------------------------------------- expansion

// expand renders one rule's candidate rows.
//
// Every `in` atom in the rule's match blocks — on `recognized` or on any
// other declared tag — expands into one row per member, and the rule's rows
// are the CARTESIAN PRODUCT of its expanding atoms. In each expanded row
// the `in` atom becomes an `eq` atom on the chosen member, still carrying
// block `match` (`0002:C13`).
//
// The expansion suffix is the sequence of chosen members, one per atom with
// MORE THAN ONE member, taken in the atoms' sort order. A single-member
// `in` therefore expands to one row whose suffix is empty, so `eq = "x"`
// and `in = ["x"]` are one spelling of one edge and cannot mint two
// identities.
func expand(base Row, predicates []Atom, outcome Atom, writes []TagValue) []Row {
	// The outcome atom participates in the product like any other
	// expanding match atom, in the same sort order, but its chosen member
	// lands in the row's outcome field rather than its predicate set.
	type choicePoint struct {
		// index is the atom's position in the sorted expansion order; the
		// outcome atom is identified by outcomeIndex.
		atom    Atom
		outcome bool
	}

	candidates := make([]choicePoint, 0, len(predicates)+1)
	for _, a := range predicates {
		candidates = append(candidates, choicePoint{atom: a})
	}
	candidates = append(candidates, choicePoint{atom: outcome, outcome: true})
	slices.SortFunc(candidates, func(x, y choicePoint) int {
		return compareAtoms(x.atom, y.atom)
	})

	// combos enumerates one chosen member per expanding atom, in sort
	// order. A non-expanding atom contributes its own single member and no
	// suffix element.
	rows := []Row{{
		ModelID:       base.ModelID,
		RuleID:        base.RuleID,
		SourceLocator: base.SourceLocator,
		Source:        base.Source,
		Gate:          base.Gate,
		RequiresOwned: base.RequiresOwned,
		Escape:        base.Escape,
		setKeys:       base.setKeys,
	}}
	atomSets := [][]Atom{nil}

	for _, c := range candidates {
		expanding := c.atom.Operator == "in"
		members := c.atom.Literal
		if !expanding {
			members = []string{""}
		}
		// A suffix element is emitted only for an atom with MORE THAN one
		// member, so a single-member `in` mints no second identity.
		suffixed := expanding && len(members) > 1

		grown := make([]Row, 0, len(rows)*len(members))
		grownAtoms := make([][]Atom, 0, len(rows)*len(members))
		for i, row := range rows {
			for _, member := range members {
				next := row
				next.Suffix = slices.Clone(row.Suffix)
				if suffixed {
					next.Suffix = append(next.Suffix, member)
				}

				atoms := slices.Clone(atomSets[i])
				switch {
				case c.outcome && expanding:
					next.Outcome = member
				case c.outcome:
					next.Outcome = c.atom.Literal[0]
				case expanding:
					// The `in` atom becomes an `eq` atom on the chosen
					// member, still carrying block match.
					atoms = append(atoms, Atom{
						Key:      c.atom.Key,
						Block:    c.atom.Block,
						Operator: "eq",
						Literal:  []string{member},
					})
				default:
					atoms = append(atoms, c.atom)
				}

				grown = append(grown, next)
				grownAtoms = append(grownAtoms, atoms)
			}
		}
		rows = grown
		atomSets = grownAtoms
	}

	for i := range rows {
		// The product can re-derive an atom an inherited context already
		// contributed, so the row's set is merged rather than concatenated.
		rows[i].Atoms = mergeAtoms(atomSets[i])
		// NextTags and Writes are built per row from the same rendered
		// assignments, never by aliasing one to the other (`0002:C15`).
		rows[i].NextTags = cloneTagValues(writes)
		rows[i].Writes = cloneTagValues(writes)
	}
	return rows
}

func cloneTagValues(in []TagValue) []TagValue {
	if len(in) == 0 {
		return nil
	}
	out := make([]TagValue, 0, len(in))
	for _, t := range in {
		out = append(out, TagValue{Key: t.Key, Value: slices.Clone(t.Value)})
	}
	return out
}

// compareRows is the dump's total row ordering: the identity tuple (model
// id, rule id, expansion suffix), field by field, byte-lexicographically.
// The suffix compares as a SEQUENCE — element by element, a shorter
// sequence that is a prefix of a longer one sorting first, so the empty
// suffix sorts before any non-empty one. The source locator does NOT
// participate (`0002:C19`).
func compareRows(a, b Row) int {
	if c := strings.Compare(a.ModelID, b.ModelID); c != 0 {
		return c
	}
	if c := strings.Compare(a.RuleID, b.RuleID); c != 0 {
		return c
	}
	return slices.Compare(a.Suffix, b.Suffix)
}
