// Command rdr0002spike2 parses the sparse TOML transition model authored to
// the JDR 0001 §D7 closed layout and normalizes it into expanded candidate
// rows, then dumps them.
//
// It is the Stage 4 (iteration 2) witness for RDR 0002's A1: that the §D7
// layout decodes unambiguously under STRICT decoding, and that the §D6
// block-keyed routing, general match-block expansion, and sequence-valued
// expansion suffix normalize as the Normative Contracts state.
//
// Usage: rdr0002spike2 [-strict-only] <fixture.toml>...
package main

import (
	"bytes"
	"cmp"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// Kernel constants RDR 0007 exports and this normalizer must emit verbatim
// for existence atoms. Mirrored here because the reshape has not landed in
// internal/resolve yet (RDR 0007 is Final, unimplemented) — confirmed this
// stage: no OpExists / LiteralTrue / LiteralFalse anywhere under internal/.
const (
	opExists     = "exists"
	literalTrue  = "true"
	literalFalse = "false"
)

// recognizedKey is the reserved tag key RDR 0008 owns. Every rule binds its
// outcome through exactly one atom in its MATCH blocks on this key.
const recognizedKey = "recognized"

// clearSentinel is reserved (JDR 0001 §D5): refused wherever a tag value is
// authored, and the rendering of an authored clear.
const clearSentinel = "<clear>"

// suffixSep joins the rendered identity's rule id and suffix elements. It is
// banned from rule ids, outcome literals, and match-block `in` members.
const suffixSep = "#"

// ---------------------------------------------------------------- source

// Model is the §D7 closed layout. Every root key and table below is admitted;
// strict decoding refuses any other.
type Model struct {
	Outcomes []string       `toml:"outcomes"`
	Terminal []string       `toml:"terminal"`
	Model    ModelMeta      `toml:"model"`
	Initial  map[string]any `toml:"initial"`
	Tags     map[string]Tag `toml:"tags"`
	// Three capability tables — §D7(ii). Each decodes to its own shape, so
	// "exactly one capability per accessor" holds structurally.
	Read    map[string]Accessor `toml:"read"`
	Write   map[string]Accessor `toml:"write"`
	Gate    map[string]Accessor `toml:"gate"`
	Context map[string]Context  `toml:"context"`
	Rule    []Rule              `toml:"rule"`
	Dump    Dump                `toml:"dump"`
}

type ModelMeta struct {
	ID          string `toml:"id"`
	Version     int    `toml:"version"`
	Description string `toml:"description"`
	// Metadata is the one sanctioned extension namespace (§D7(v)): decoded
	// free-form, carried untouched, never interpreted. Being a map, the
	// strict decoder does not descend into it.
	Metadata map[string]any `toml:"metadata"`
}

// Tag is a declaration: provenance (this RDR's) plus the seven type-model
// wire keys (§D7(iii); RDR 0003 owns what they mean).
type Tag struct {
	Provenance   string `toml:"provenance"`
	Kind         string `toml:"kind"`
	Domain       []any  `toml:"domain"`
	Min          *int   `toml:"min"`
	Max          *int   `toml:"max"`
	Elements     []any  `toml:"elements"`
	SingleValued *bool  `toml:"single_valued"`
	Required     *bool  `toml:"required"`
}

// Accessor is one capability-table entry (§D7(ii)). ReadBack is meaningful
// only under [write.<id>]; keys is the sole tag<->accessor binding.
type Accessor struct {
	Role     string   `toml:"role"`
	Path     string   `toml:"path"`
	Keys     []string `toml:"keys"`
	Timeout  string   `toml:"timeout"`
	ReadBack bool     `toml:"read_back"`
}

type Context struct {
	Inherits string                    `toml:"inherits"`
	Match    map[string]map[string]any `toml:"match"`
}

type Rule struct {
	ID     string                    `toml:"id"`
	Use    []string                  `toml:"use"`
	Source string                    `toml:"source"`
	Gate   []string                  `toml:"gate"`
	Clear  []string                  `toml:"clear"`
	Escape []string                  `toml:"escape"`
	Match  map[string]map[string]any `toml:"match"`
	Guard  Guard                     `toml:"guard"`
	Write  map[string]any            `toml:"write"`
}

type Guard struct {
	All    map[string]map[string]any `toml:"all"`
	Unless map[string]map[string]any `toml:"unless"`
}

type Dump struct {
	Order []string `toml:"order"`
}

// versionProbe is the permissive first pass: it reads [model].version and
// nothing else, so the version gate can precede strict field decoding.
type versionProbe struct {
	Model struct {
		Version int `toml:"version"`
	} `toml:"model"`
}

// --------------------------------------------------------------- normalized

// Atom is one predicate atom. Block is the three-valued authored block —
// "match", "all", or "unless" — and is the key the kernel handoff routes on
// (JDR 0001 §D6). Literal holds a set literal as a member SEQUENCE, never a
// joined string.
type Atom struct {
	Key      string
	Operator string
	Literal  []string
	Block    string
}

// literalString renders a literal for DISPLAY only. It joins members, so it
// MUST NOT reach any identity, merge key, dedup, or sort — see identity().
func (a Atom) literalString() string { return strings.Join(a.Literal, ",") }

// identity is the full atom identity tuple the merge is a set over and the
// dump sorts on: (key, block, operator token, literal). The literal enters as
// its member SEQUENCE, length-prefixed per member, never as a joined string:
// any delimiter is authorable inside a tag value, so a joined rendering would
// collapse ["a,b","c"] and ["a","b,c"] into one atom and turn a dead rule into
// a live one. Length prefixing makes the encoding injective for any member.
func (a Atom) identity() string {
	var b strings.Builder
	b.WriteString(a.Key)
	b.WriteString("\x00")
	b.WriteString(a.Block)
	b.WriteString("\x00")
	b.WriteString(a.Operator)
	b.WriteString("\x00")
	b.WriteString(memberKey(a.Literal))
	return b.String()
}

// memberKey encodes a member sequence injectively: each member is prefixed by
// its byte length, so no member content can forge a boundary.
func memberKey(members []string) string {
	var b strings.Builder
	for _, m := range members {
		fmt.Fprintf(&b, "%d:%s", len(m), m)
	}
	return b.String()
}

type Write struct {
	Key   string
	Value []string
	Clear bool
}

// identity compares a write by its member SEQUENCE, for the same reason
// Atom.identity does: RDR 0004 reads a held value back by equality, so a
// joined rendering would decide whether a write reports as applied.
func (w Write) identity() string {
	if w.Clear {
		return w.Key + "\x00" + clearSentinel
	}
	return w.Key + "\x00" + memberKey(w.Value)
}

func (w Write) render() string {
	if w.Clear {
		return w.Key + "=" + clearSentinel
	}
	// Display only -- never an identity (see Write.identity). A set is spelled
	// visibly as a set, and each member is QUOTED: an unquoted separator is
	// forgeable by a member that contains it, so ["x|y","z"] and ["x","y|z"]
	// would render alike and a test asserting on the rendered value would pass
	// the very collision this RDR forbids. The delimiter-parameterized control
	// caught exactly that.
	if len(w.Value) == 1 {
		return w.Key + "=" + w.Value[0]
	}
	return w.Key + "=" + renderMembers(w.Value)
}

// renderMembers spells a member sequence unambiguously for display: each
// member quoted with %q, so no member content can forge a boundary.
func renderMembers(members []string) string {
	parts := make([]string, 0, len(members))
	for _, m := range members {
		parts = append(parts, fmt.Sprintf("%q", m))
	}
	return "[" + strings.Join(parts, " ") + "]"
}

// Row is the normalized candidate row. Identity is (ModelID, RuleID, Suffix),
// Suffix being a SEQUENCE of chosen members (one per multi-member match-block
// `in`), in atom sort order. Locator is diagnostic and never orders.
type Row struct {
	ModelID       string
	RuleID        string
	Suffix        []string
	Locator       string
	Outcome       string
	Atoms         []Atom
	NextTags      []Write
	Writes        []Write
	RequiresOwned []string
	GateIDs       []string
	EscapeClasses []string
}

// kind is a DERIVED view property, never a stored field: escape identity is
// discriminated solely by a non-empty escape class list (RDR 0009).
func (r Row) kind() string {
	if len(r.EscapeClasses) > 0 {
		return "escape"
	}
	return "transition"
}

func main() {
	strictOnly := flag.Bool("strict-only", false, "decode strictly and report, without normalizing")
	flag.Parse()

	for _, path := range flag.Args() {
		data, err := os.ReadFile(path)
		must(err)
		model, err := load(data)
		must(err)
		if *strictOnly {
			fmt.Printf("STRICT-OK %s model=%s metadata=%d\n",
				path, model.Model.ID, len(model.Model.Metadata))
			continue
		}
		rows, err := normalize(model)
		must(err)
		fmt.Printf("MODEL %s rows=%d outcomes=%s terminal=%s metadata_keys=%s\n",
			model.Model.ID, len(rows), strings.Join(model.Outcomes, ","),
			strings.Join(model.Terminal, ","), strings.Join(sortedKeys(model.Model.Metadata), ","))
		for _, row := range rows {
			fmt.Println(render(row))
		}
	}
}

// load runs the two-pass version gate then strict decoding (§D7 / the version
// clause): read [model].version permissively, refuse any value but 1, and
// only then decode the document strictly.
func load(data []byte) (Model, error) {
	var probe versionProbe
	if err := toml.Unmarshal(data, &probe); err != nil {
		return Model{}, fmt.Errorf("malformed TOML: %w", err)
	}
	if probe.Model.Version != 1 {
		return Model{}, fmt.Errorf("unsupported version %d", probe.Model.Version)
	}

	var model Model
	dec := toml.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&model); err != nil {
		var strictErr *toml.StrictMissingError
		if ok := asStrict(err, &strictErr); ok {
			return Model{}, fmt.Errorf("unknown schema field: %s", strictKeys(strictErr))
		}
		return Model{}, fmt.Errorf("malformed TOML: %w", err)
	}
	if err := validate(model); err != nil {
		return Model{}, err
	}
	return model, nil
}

func asStrict(err error, target **toml.StrictMissingError) bool {
	if e, ok := err.(*toml.StrictMissingError); ok {
		*target = e
		return true
	}
	return false
}

func strictKeys(e *toml.StrictMissingError) string {
	parts := make([]string, 0, len(e.Errors))
	for i := range e.Errors {
		parts = append(parts, strings.Join(e.Errors[i].Key(), "."))
	}
	slices.Sort(parts)
	return strings.Join(parts, ", ")
}

// ---------------------------------------------------------------- validate

func validate(m Model) error {
	if m.Model.ID == "" {
		return fmt.Errorf("malformed model: missing id")
	}
	if err := validateAlphabet(m); err != nil {
		return err
	}
	if err := validateTags(m); err != nil {
		return err
	}
	if err := validateAccessors(m); err != nil {
		return err
	}
	if err := validateInitial(m); err != nil {
		return err
	}
	for _, id := range m.Terminal {
		if _, ok := m.Context[id]; !ok {
			return fmt.Errorf("unknown context %q in terminal", id)
		}
	}
	return validateRules(m)
}

func validateAlphabet(m Model) error {
	if len(m.Outcomes) == 0 {
		return fmt.Errorf("malformed recognized outcome alphabet: empty")
	}
	seen := map[string]bool{}
	for _, o := range m.Outcomes {
		switch {
		case o == "":
			return fmt.Errorf("malformed recognized outcome alphabet: contains the empty string")
		case seen[o]:
			return fmt.Errorf("malformed recognized outcome alphabet: duplicate member %q", o)
		case strings.Contains(o, suffixSep):
			return fmt.Errorf("malformed recognized outcome alphabet: member %q contains %q", o, suffixSep)
		}
		seen[o] = true
	}
	return nil
}

func validateTags(m Model) error {
	recognized := 0
	for name, tag := range m.Tags {
		switch {
		case tag.Provenance == "recognized":
			recognized++
			if name != recognizedKey {
				return fmt.Errorf("reserved_tag_key: recognized declaration named %q", name)
			}
		case name == recognizedKey:
			return fmt.Errorf("reserved_tag_key: %s declaration named %q", tag.Provenance, name)
		case tag.Provenance != "owned" && tag.Provenance != "observed":
			return fmt.Errorf("malformed tag declaration: %q has provenance %q", name, tag.Provenance)
		}
		if tag.Kind == "" {
			return fmt.Errorf("malformed tag declaration: %q has no kind", name)
		}
		for _, d := range tag.Domain {
			if fmt.Sprint(d) == clearSentinel {
				return fmt.Errorf("reserved tag value: %q domain contains %s", name, clearSentinel)
			}
		}
	}
	if recognized != 1 {
		return fmt.Errorf("expected exactly one recognized declaration, got %d", recognized)
	}
	return nil
}

// validateAccessors runs the PROVENANCE-SCOPED binding validations this RDR
// settles for §D7(ii): every owned tag exactly one reader; an observed tag at
// most one (zero legal, JD-9); every written/cleared/initial key exactly one
// writer and owned; recognized served by neither.
func validateAccessors(m Model) error {
	readers := map[string]int{}
	writers := map[string]int{}
	collect := func(tables map[string]Accessor, counts map[string]int, what string) error {
		for id, acc := range tables {
			if acc.Timeout == "" {
				return fmt.Errorf("malformed accessor declaration: %s %q has no timeout", what, id)
			}
			if len(acc.Keys) == 0 {
				return fmt.Errorf("malformed accessor declaration: %s %q has an empty keys list", what, id)
			}
			for _, key := range acc.Keys {
				tag, ok := m.Tags[key]
				if !ok {
					return fmt.Errorf("unknown tag %q in %s %q keys", key, what, id)
				}
				if key == recognizedKey {
					return fmt.Errorf("malformed accessor binding: %s %q serves the reserved key %q", what, id, recognizedKey)
				}
				if what == "write" && tag.Provenance != "owned" {
					return fmt.Errorf("write to non-owned tag: writer %q serves %q (%s)", id, key, tag.Provenance)
				}
				counts[key]++
			}
		}
		return nil
	}
	if err := collect(m.Read, readers, "read"); err != nil {
		return err
	}
	if err := collect(m.Write, writers, "write"); err != nil {
		return err
	}
	for _, id := range sortedKeys(m.Gate) {
		acc := m.Gate[id]
		if acc.Timeout == "" {
			return fmt.Errorf("malformed accessor declaration: gate %q has no timeout", id)
		}
		for _, key := range acc.Keys {
			if _, ok := m.Tags[key]; !ok {
				return fmt.Errorf("unknown tag %q in gate %q keys", key, id)
			}
			if key == recognizedKey {
				return fmt.Errorf("malformed accessor binding: gate %q serves the reserved key %q", id, recognizedKey)
			}
		}
	}
	for _, name := range sortedKeys(m.Tags) {
		tag := m.Tags[name]
		switch tag.Provenance {
		case "owned":
			if readers[name] != 1 {
				return fmt.Errorf("malformed accessor binding: owned tag %q served by %d readers, want exactly 1", name, readers[name])
			}
		case "observed":
			if readers[name] > 1 {
				return fmt.Errorf("malformed accessor binding: observed tag %q served by %d readers, want at most 1", name, readers[name])
			}
		case "recognized":
			if readers[name]+writers[name] > 0 {
				return fmt.Errorf("malformed accessor binding: %q is served by an accessor", name)
			}
		}
	}
	// Every key a rule writes or clears, and every [initial] key, needs
	// exactly one writer.
	need := map[string]bool{}
	for key := range m.Initial {
		need[key] = true
	}
	for _, rule := range m.Rule {
		for key := range rule.Write {
			need[key] = true
		}
		for _, key := range rule.Clear {
			need[key] = true
		}
	}
	for _, key := range sortedBoolKeys(need) {
		if _, ok := m.Tags[key]; !ok {
			return fmt.Errorf("unknown tag %q written", key)
		}
		if writers[key] != 1 {
			return fmt.Errorf("malformed accessor binding: written key %q served by %d writers, want exactly 1", key, writers[key])
		}
	}
	return nil
}

func validateInitial(m Model) error {
	for _, key := range sortedKeys(m.Initial) {
		tag, ok := m.Tags[key]
		if !ok {
			return fmt.Errorf("malformed initial declaration: unknown tag %q", key)
		}
		if tag.Provenance != "owned" {
			return fmt.Errorf("malformed initial declaration: %q is %s, want owned", key, tag.Provenance)
		}
		if fmt.Sprint(m.Initial[key]) == clearSentinel {
			return fmt.Errorf("reserved tag value: [initial] %q is %s", key, clearSentinel)
		}
		// The VALUE arm: well-formed for the tag's declared kind and domain.
		// RDR 0003 owns what the keys mean; this is the rejection site.
		if err := wellFormed(tag, m.Initial[key]); err != nil {
			return fmt.Errorf("malformed initial declaration: %q %w", key, err)
		}
	}
	return nil
}

// wellFormed applies RDR 0003's declaration rules at this RDR's two rejection
// sites: a literal outside a declared finite domain, and an int outside
// min/max. Deliberately narrow — the full grammar is RDR 0003's.
func wellFormed(tag Tag, value any) error {
	lit := renderValue(value)
	if len(tag.Domain) > 0 {
		domain := make([]string, 0, len(tag.Domain))
		for _, d := range tag.Domain {
			domain = append(domain, renderValue(d))
		}
		if !slices.Contains(domain, lit) {
			return fmt.Errorf("value %q is outside declared domain [%s]", lit, strings.Join(domain, ","))
		}
	}
	if tag.Kind == "int" {
		n, ok := value.(int64)
		if !ok {
			return fmt.Errorf("value %q is not an int", lit)
		}
		if tag.Min != nil && n < int64(*tag.Min) {
			return fmt.Errorf("value %d is below min %d", n, *tag.Min)
		}
		if tag.Max != nil && n > int64(*tag.Max) {
			return fmt.Errorf("value %d is above max %d", n, *tag.Max)
		}
	}
	return nil
}

func validateRules(m Model) error {
	ids := map[string]bool{}
	for _, rule := range m.Rule {
		if strings.Contains(rule.ID, suffixSep) {
			return fmt.Errorf("malformed rule id: %q contains %q", rule.ID, suffixSep)
		}
		if ids[rule.ID] {
			return fmt.Errorf("duplicate rule id %q", rule.ID)
		}
		ids[rule.ID] = true

		isEscape := len(rule.Escape) > 0
		if isEscape {
			if len(rule.Write) > 0 || len(rule.Clear) > 0 || len(rule.Gate) > 0 {
				return fmt.Errorf("malformed escape declaration: escape rule %q carries a write block, clear list, or gate list", rule.ID)
			}
			for _, class := range rule.Escape {
				if class != "no_match" && class != "ambiguous_match" {
					return fmt.Errorf("malformed escape declaration: rule %q models %q", rule.ID, class)
				}
			}
		} else if len(rule.Write) == 0 {
			return fmt.Errorf("malformed rule shape: ordinary rule %q carries no write block", rule.ID)
		}

		for _, id := range rule.Gate {
			if _, ok := m.Gate[id]; !ok {
				return fmt.Errorf("unknown accessor: rule %q gates on %q", rule.ID, id)
			}
		}
		for _, key := range sortedKeys(rule.Write) {
			if err := ownedTag(m, key, rule.ID); err != nil {
				return err
			}
			if fmt.Sprint(rule.Write[key]) == clearSentinel {
				return fmt.Errorf("reserved tag value: rule %q writes %s to %q", rule.ID, clearSentinel, key)
			}
		}
		for _, key := range rule.Clear {
			if err := ownedTag(m, key, rule.ID); err != nil {
				return err
			}
		}
		// Match blocks admit only eq and in (§D6); a recognized atom in a
		// guard block is a malformed outcome binding.
		if err := checkMatchBlock(m, rule.ID, rule.Match); err != nil {
			return err
		}
		guardBlocks := []struct {
			name  string
			block map[string]map[string]any
		}{{"guard.all", rule.Guard.All}, {"guard.unless", rule.Guard.Unless}}
		for _, gb := range guardBlocks {
			for key := range gb.block {
				if key == recognizedKey {
					return fmt.Errorf("malformed outcome binding: rule %q authors a %q atom in a guard block", rule.ID, recognizedKey)
				}
			}
			// A15: every atom rule applies here exactly as in a match block.
			if err := checkAtomBlock(m, "rule "+rule.ID, gb.name, gb.block); err != nil {
				return err
			}
		}
		for _, id := range rule.Use {
			if _, ok := m.Context[id]; !ok {
				return fmt.Errorf("unknown context %q used by rule %q", id, rule.ID)
			}
		}
	}
	for _, id := range sortedKeys(m.Context) {
		if err := checkMatchBlock(m, "context "+id, m.Context[id].Match); err != nil {
			return err
		}
	}
	return nil
}

func checkMatchBlock(m Model, owner string, block map[string]map[string]any) error {
	return checkAtomBlock(m, owner, "match", block)
}

// guardOperators is RDR 0003's closed set. This RDR does not mint or widen it;
// it cites it by member list so `unknown operator` is decidable at load.
var guardOperators = []string{"eq", "in", "lt", "lte", "gt", "gte", "exists", "contains"}

// checkAtomBlock enforces every atom-level rule this RDR states, in EVERY atom
// block — `match`, `guard.all`, and `guard.unless` (A15, the block-agnostic
// clause). Only two rules are genuinely match-only and stay gated on isMatch:
// the eq/in operator restriction (§D6 routing) and the `#` reservation in `in`
// members (the expansion suffix separator, and only match blocks expand).
func checkAtomBlock(m Model, owner, blockName string, block map[string]map[string]any) error {
	isMatch := blockName == "match"
	for _, key := range sortedKeys(block) {
		if _, ok := m.Tags[key]; !ok {
			return fmt.Errorf("unknown tag %q matched by %s", key, owner)
		}
		for _, op := range sortedKeys(block[key]) {
			if isMatch {
				if op != "eq" && op != "in" {
					return fmt.Errorf("malformed predicate atom: %s matches %s with operator %q under a match block", owner, key, op)
				}
			} else if !slices.Contains(guardOperators, op) {
				return fmt.Errorf("malformed predicate atom: %s constrains %s with unknown operator %q under a %s block", owner, key, op, blockName)
			}
			tag := m.Tags[key]
			raw := block[key][op]
			// The <clear> and `#` rules bind every member regardless of
			// operator; the domain/kind rule is operator-scoped (below).
			for _, lit := range literalMembers(raw) {
				if lit == clearSentinel {
					return fmt.Errorf("reserved tag value: %s matches %s against %s", owner, key, clearSentinel)
				}
				if isMatch && op == "in" && strings.Contains(lit, suffixSep) {
					return fmt.Errorf("malformed predicate atom: %s has match-block `in` member %q containing %q", owner, lit, suffixSep)
				}
			}
			// Domain/kind conformance is per-operator, because operators do
			// not all take a domain MEMBER as their right-hand side:
			//   eq / in / contains -> a member (or set of members)
			//   exists             -> a bool literal, never a member
			//   lt/lte/gt/gte      -> an ordered BOUND, kind-checked but not
			//                         domain-checked (a bound need not itself
			//                         be an authorable value)
			switch op {
			case "exists":
				lit := renderValue(raw)
				if lit != literalTrue && lit != literalFalse {
					return fmt.Errorf("malformed predicate atom: %s.exists has literal %q", key, lit)
				}
			case "lt", "lte", "gt", "gte":
				if tag.Kind == "int" {
					if _, ok := raw.(int64); !ok {
						return fmt.Errorf("malformed predicate atom: %s constrains %s with non-int bound %q", owner, key, renderValue(raw))
					}
				}
			default:
				for _, member := range literalValues(raw) {
					if err := wellFormed(tag, member); err != nil {
						return fmt.Errorf("malformed predicate atom: %s matches %s: %w", owner, key, err)
					}
				}
			}
		}
	}
	return nil
}

func ownedTag(m Model, key, ruleID string) error {
	tag, ok := m.Tags[key]
	if !ok {
		return fmt.Errorf("unknown tag %q written by rule %q", key, ruleID)
	}
	if tag.Provenance != "owned" {
		return fmt.Errorf("write to non-owned tag: rule %q writes %q (%s)", ruleID, key, tag.Provenance)
	}
	return nil
}

// --------------------------------------------------------------- normalize

func normalize(model Model) ([]Row, error) {
	var rows []Row
	for _, rule := range model.Rule {
		// One unified atom set, keyed on the FULL atom identity, so two
		// atoms agreeing on (key, block, operator) but differing in literal
		// both survive, and a byte-identical atom collapses to one (A13).
		atoms := map[string]Atom{}
		for _, contextID := range rule.Use {
			if err := mergeContext(model.Context, contextID, atoms, map[string]bool{}); err != nil {
				return nil, err
			}
		}
		if err := mergeAtoms(atoms, "match", rule.Match); err != nil {
			return nil, err
		}
		if err := mergeAtoms(atoms, "all", rule.Guard.All); err != nil {
			return nil, err
		}
		if err := mergeAtoms(atoms, "unless", rule.Guard.Unless); err != nil {
			return nil, err
		}

		set := make([]Atom, 0, len(atoms))
		for _, a := range atoms {
			set = append(set, a)
		}
		slices.SortFunc(set, compareAtoms)

		// Outcome binding reads the MATCH blocks only.
		outcomeAtom, err := findOutcomeAtom(set, model, rule.ID)
		if err != nil {
			return nil, err
		}

		// General match-block expansion (§D6): the rule's rows are the
		// Cartesian product of every match-block `in` atom, the recognized
		// one included. Only atoms with MORE THAN ONE member contribute a
		// suffix element; the suffix is the sequence of chosen members, one
		// per such atom, in atom sort order.
		combos := product(set)

		writes, next := renderWrites(rule)
		requires := requiresOwned(rule)
		gates := slices.Sorted(slices.Values(rule.Gate))
		escapes := slices.Sorted(slices.Values(rule.Escape))

		_ = outcomeAtom
		for _, combo := range combos {
			rowAtoms := make([]Atom, 0, len(set))
			var suffix []string
			outcome := ""
			for i, a := range set {
				if chosen, expanding := combo[i]; expanding {
					// An expanded `in` becomes an `eq` on the chosen member,
					// still carrying block "match", so every match-block atom
					// the kernel receives is an equality Match can test.
					if len(a.Literal) > 1 {
						suffix = append(suffix, chosen)
					}
					a = Atom{Key: a.Key, Operator: "eq", Literal: []string{chosen}, Block: a.Block}
				}
				if a.Key == recognizedKey && a.Block == "match" {
					// Lifted out of the predicate set into the outcome field.
					outcome = a.Literal[0]
					continue
				}
				rowAtoms = append(rowAtoms, a)
			}
			if outcome == "" {
				return nil, fmt.Errorf("malformed outcome binding: rule %q bound no outcome", rule.ID)
			}
			slices.SortFunc(rowAtoms, compareAtoms)
			rows = append(rows, Row{
				ModelID:       model.Model.ID,
				RuleID:        rule.ID,
				Suffix:        suffix,
				Locator:       rule.Source,
				Outcome:       outcome,
				Atoms:         rowAtoms,
				NextTags:      next,
				Writes:        writes,
				RequiresOwned: requires,
				GateIDs:       gates,
				EscapeClasses: escapes,
			})
		}
	}
	slices.SortFunc(rows, compareRows)
	return rows, nil
}

// findOutcomeAtom enforces "exactly one recognized atom in the match blocks,
// eq or in, every literal an alphabet member".
func findOutcomeAtom(set []Atom, model Model, ruleID string) (Atom, error) {
	var found []Atom
	for _, a := range set {
		if a.Key == recognizedKey {
			if a.Block != "match" {
				return Atom{}, fmt.Errorf("malformed outcome binding: rule %q authors a %q atom in block %q", ruleID, recognizedKey, a.Block)
			}
			found = append(found, a)
		}
	}
	if len(found) != 1 {
		return Atom{}, fmt.Errorf("malformed outcome binding: rule %q binds %d recognized atoms, want exactly 1", ruleID, len(found))
	}
	atom := found[0]
	if atom.Operator != "eq" && atom.Operator != "in" {
		return Atom{}, fmt.Errorf("malformed outcome binding: rule %q binds outcome with operator %q", ruleID, atom.Operator)
	}
	for _, lit := range atom.Literal {
		if !slices.Contains(model.Outcomes, lit) {
			return Atom{}, fmt.Errorf("malformed outcome binding: rule %q binds outcome %q outside the alphabet", ruleID, lit)
		}
	}
	return atom, nil
}

// product yields the Cartesian product of every match-block `in` atom, as a
// slice of maps from atom index to the chosen member. A SINGLE-member `in`
// participates (so its atom still becomes an `eq`) but multiplies the product
// by one and contributes no suffix element — which is what makes `eq = "x"`
// and `in = ["x"]` one spelling of one edge. Atoms are visited in the set's
// sort order, so the suffix sequence is in atom sort order by construction.
func product(set []Atom) []map[int]string {
	combos := []map[int]string{{}}
	for i, a := range set {
		if a.Block != "match" || a.Operator != "in" {
			continue
		}
		grown := make([]map[int]string, 0, len(combos)*len(a.Literal))
		for _, base := range combos {
			for _, member := range a.Literal {
				next := make(map[int]string, len(base)+1)
				for k, v := range base {
					next[k] = v
				}
				next[i] = member
				grown = append(grown, next)
			}
		}
		combos = grown
	}
	return combos
}

// --------------------------------------------------------------- rendering

func render(r Row) string {
	atoms := make([]string, 0, len(r.Atoms))
	for _, a := range r.Atoms {
		atoms = append(atoms, fmt.Sprintf("%s.%s=%s@%s", a.Key, a.Operator, a.literalString(), a.Block))
	}
	writes := make([]string, 0, len(r.Writes))
	for _, w := range r.Writes {
		writes = append(writes, w.render())
	}
	next := make([]string, 0, len(r.NextTags))
	for _, w := range r.NextTags {
		next = append(next, w.render())
	}
	return fmt.Sprintf(
		"%s kind=%s source=%s outcome=%s atoms=[%s] next=[%s] write=[%s] requires_owned=[%s] gate=[%s] escape=[%s]",
		identity(r), r.kind(), r.Locator, r.Outcome,
		strings.Join(atoms, "; "), strings.Join(next, "; "), strings.Join(writes, "; "),
		strings.Join(r.RequiresOwned, ","), strings.Join(r.GateIDs, ","),
		strings.Join(r.EscapeClasses, ","))
}

// identity renders (model id, rule id, expansion suffix), joining the rule id
// and each suffix element with `#`. `#` is banned from every string that can
// reach here, so splitting on it is unambiguous.
func identity(r Row) string {
	id := r.ModelID + "." + r.RuleID
	for _, s := range r.Suffix {
		id += suffixSep + s
	}
	return id
}

// compareRows is the dump's total row ordering: (model id, rule id, expansion
// suffix), field by field, byte-lexicographically. The suffix compares as a
// SEQUENCE — element by element, a shorter prefix sorting first, so the empty
// suffix sorts before any non-empty one. The locator does NOT participate.
func compareRows(a, b Row) int {
	if c := cmp.Compare(a.ModelID, b.ModelID); c != 0 {
		return c
	}
	if c := cmp.Compare(a.RuleID, b.RuleID); c != 0 {
		return c
	}
	return slices.Compare(a.Suffix, b.Suffix)
}

func compareAtoms(a, b Atom) int {
	if c := cmp.Compare(a.Key, b.Key); c != 0 {
		return c
	}
	if c := cmp.Compare(a.Block, b.Block); c != 0 {
		return c
	}
	if c := cmp.Compare(a.Operator, b.Operator); c != 0 {
		return c
	}
	return slices.Compare(a.Literal, b.Literal)
}

func renderWrites(rule Rule) (writes, next []Write) {
	if len(rule.Escape) > 0 {
		return nil, nil
	}
	for _, key := range sortedKeys(rule.Write) {
		writes = append(writes, Write{Key: key, Value: literalMembers(rule.Write[key])})
	}
	for _, key := range rule.Clear {
		writes = append(writes, Write{Key: key, Clear: true})
	}
	slices.SortFunc(writes, func(a, b Write) int { return cmp.Compare(a.Key, b.Key) })
	// Populated explicitly, NOT by aliasing: the two fields are distinct
	// (RDR 0009 A4) and an alias would make a future divergence silent.
	next = make([]Write, len(writes))
	copy(next, writes)
	return writes, next
}

func requiresOwned(rule Rule) []string {
	if len(rule.Escape) > 0 {
		return nil
	}
	seen := map[string]bool{}
	for key := range rule.Write {
		seen[key] = true
	}
	for _, key := range rule.Clear {
		seen[key] = true
	}
	return sortedBoolKeys(seen)
}

func mergeContext(contexts map[string]Context, id string, out map[string]Atom, seen map[string]bool) error {
	ctx, ok := contexts[id]
	if !ok {
		return fmt.Errorf("unknown context %q", id)
	}
	if seen[id] {
		return fmt.Errorf("cyclic context inheritance at %q", id)
	}
	seen[id] = true
	if ctx.Inherits != "" {
		if err := mergeContext(contexts, ctx.Inherits, out, seen); err != nil {
			return err
		}
	}
	return mergeAtoms(out, "match", ctx.Match)
}

// mergeAtoms folds one authored block into the unified atom set, keyed on the
// FULL identity tuple so the merge is a set over (key, block, operator,
// literal) and never on a proper prefix of it.
func mergeAtoms(out map[string]Atom, block string, source map[string]map[string]any) error {
	for _, key := range sortedKeys(source) {
		ops := source[key]
		for _, op := range sortedKeys(ops) {
			members := literalMembers(ops[op])
			if op == opExists {
				if len(members) != 1 || (members[0] != literalTrue && members[0] != literalFalse) {
					return fmt.Errorf("malformed predicate atom: %s.%s has literal %q", key, op, strings.Join(members, ","))
				}
			}
			atom := Atom{Key: key, Operator: op, Literal: members, Block: block}
			out[atom.identity()] = atom
		}
	}
	return nil
}

// literalMembers gives every literal a member SEQUENCE — a set literal is
// never rendered into a delimiter-joined string. Members sort
// byte-lexicographically so two authored orderings of one set are one literal.
func literalMembers(value any) []string {
	if list, ok := value.([]any); ok {
		parts := make([]string, 0, len(list))
		for _, item := range list {
			parts = append(parts, renderValue(item))
		}
		slices.Sort(parts)
		return parts
	}
	return []string{renderValue(value)}
}

// literalValues gives the raw member values of a literal (never rendered), so
// kind-sensitive checks see the authored type rather than its rendering.
func literalValues(value any) []any {
	if list, ok := value.([]any); ok {
		return list
	}
	return []any{value}
}

func renderValue(value any) string {
	if b, ok := value.(bool); ok {
		if b {
			return literalTrue
		}
		return literalFalse
	}
	if list, ok := value.([]any); ok {
		parts := make([]string, 0, len(list))
		for _, item := range list {
			parts = append(parts, renderValue(item))
		}
		slices.Sort(parts)
		return strings.Join(parts, ",")
	}
	return fmt.Sprint(value)
}

func sortedKeys[V any](m map[string]V) []string {
	return slices.Sorted(maps_Keys(m))
}

func sortedBoolKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k, v := range m {
		if v {
			out = append(out, k)
		}
	}
	slices.Sort(out)
	return out
}

func maps_Keys[V any](m map[string]V) func(func(string) bool) {
	return func(yield func(string) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "refused:", err)
		os.Exit(1)
	}
}
