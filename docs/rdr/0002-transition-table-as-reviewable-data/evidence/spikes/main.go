// Command rdr0002spike parses the sparse TOML transition model and
// normalizes it into expanded candidate rows, then dumps them.
//
// It is the Stage 4 witness for RDR 0002's normalization and dump-ordering
// contracts: outcome lifting, RequiresOwned derivation, write-free escape
// rows, per-atom block retention, and the total row ordering.
package main

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// Kernel constants RDR 0007 exports and this normalizer must emit verbatim
// for existence atoms. Mirrored here because the reshape has not landed in
// internal/resolve yet (RDR 0007 is Final, unimplemented).
const (
	opExists     = "exists"
	literalTrue  = "true"
	literalFalse = "false"
)

// recognizedKey is the reserved tag key RDR 0008 owns. Every rule binds its
// outcome through exactly one atom on this key.
const recognizedKey = "recognized"

type Model struct {
	Model struct {
		ID          string
		Version     int
		Description string
	}
	Tags      map[string]Tag
	Outcomes  []string
	Accessors map[string]Accessor
	Context   map[string]Context
	Rule      []Rule
	Dump      Dump
}

type Tag struct {
	Provenance string
	Kind       string
	Accessor   string
}

type Accessor struct {
	Mode string
	Path string
}

type Context struct {
	Inherits string
	Match    map[string]map[string]any
}

type Rule struct {
	ID     string
	Use    []string
	Match  map[string]map[string]any
	Guard  Guard
	Write  map[string]any
	Clear  []string
	Escape []string
	Self   bool
	Source string
}

type Guard struct {
	All    map[string]map[string]any
	Unless map[string]map[string]any
}

type Dump struct {
	Order []string
}

// Atom is one predicate atom in a normalized row's predicate set. It retains
// the block it was authored in; normalization never folds unless into all.
type Atom struct {
	Key      string
	Operator string
	Literal  string
	Block    string // "match", "all", or "unless"
}

// Write is one rendered tag write. Clear is true for an authored clear,
// which renders as a <clear> write.
type Write struct {
	Key   string
	Value string
	Clear bool
}

func (w Write) render() string {
	if w.Clear {
		return w.Key + "=<clear>"
	}
	return w.Key + "=" + w.Value
}

// Row is the normalized candidate row: the full field set the dump contract
// requires. ModelID, RuleID and Suffix are the row identity; Locator is
// diagnostic and does NOT participate in ordering.
type Row struct {
	ModelID       string
	RuleID        string
	Suffix        string // outcome literal when an `in` atom expanded the rule
	Locator       string
	Kind          string // "transition" or "escape"
	Outcome       string
	Atoms         []Atom
	NextTags      []Write
	Writes        []Write
	RequiresOwned []string
	EscapeClasses []string
}

func main() {
	for _, path := range os.Args[1:] {
		data, err := os.ReadFile(path)
		must(err)
		var model Model
		must(toml.Unmarshal(data, &model))
		must(validate(model))
		rows, err := normalize(model)
		must(err)
		fmt.Printf("MODEL %s rows=%d outcomes=%s\n",
			model.Model.ID, len(rows), strings.Join(model.Outcomes, ","))
		for _, row := range rows {
			fmt.Println(render(row))
		}
	}
}

// render emits every field of the normalized candidate-row value, in the
// order the dump contract lists them.
func render(r Row) string {
	atoms := make([]string, 0, len(r.Atoms))
	for _, a := range r.Atoms {
		atoms = append(atoms, fmt.Sprintf("%s.%s=%s@%s", a.Key, a.Operator, a.Literal, a.Block))
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
		"%s kind=%s source=%s outcome=%s atoms=[%s] next=[%s] write=[%s] requires_owned=[%s] escape=[%s]",
		identity(r), r.Kind, r.Locator, r.Outcome,
		strings.Join(atoms, "; "), strings.Join(next, "; "), strings.Join(writes, "; "),
		strings.Join(r.RequiresOwned, ","), strings.Join(r.EscapeClasses, ","))
}

// identity renders the row identity tuple: (model id, rule id, suffix).
func identity(r Row) string {
	id := r.ModelID + "." + r.RuleID
	if r.Suffix != "" {
		id += "#" + r.Suffix
	}
	return id
}

// validate applies the load-time checks the normative contracts require, so
// the spike refuses a malformed model rather than normalizing it.
func validate(m Model) error {
	if m.Model.Version != 1 {
		return fmt.Errorf("unsupported version %d", m.Model.Version)
	}
	if len(m.Outcomes) == 0 {
		return fmt.Errorf("missing recognized outcome alphabet")
	}
	seen := map[string]bool{}
	for _, o := range m.Outcomes {
		if o == "" {
			return fmt.Errorf("outcome alphabet contains the empty string")
		}
		if seen[o] {
			return fmt.Errorf("duplicate outcome %q", o)
		}
		seen[o] = true
	}
	recognized := 0
	for name, tag := range m.Tags {
		if tag.Provenance == "recognized" {
			recognized++
			if name != recognizedKey {
				return fmt.Errorf("reserved_tag_key: recognized declaration named %q", name)
			}
		} else if name == recognizedKey {
			return fmt.Errorf("reserved_tag_key: %s declaration named %q", tag.Provenance, name)
		}
		if tag.Accessor != "" {
			if _, ok := m.Accessors[tag.Accessor]; !ok {
				return fmt.Errorf("unknown accessor %q on tag %q", tag.Accessor, name)
			}
		}
	}
	if recognized != 1 {
		return fmt.Errorf("expected exactly one recognized declaration, got %d", recognized)
	}
	ruleIDs := map[string]bool{}
	for _, rule := range m.Rule {
		if ruleIDs[rule.ID] {
			return fmt.Errorf("duplicate rule id %q", rule.ID)
		}
		ruleIDs[rule.ID] = true
		if len(rule.Escape) > 0 && (len(rule.Write) > 0 || len(rule.Clear) > 0) {
			return fmt.Errorf("escape rule %q carries a write block or clear list", rule.ID)
		}
		for key := range rule.Write {
			if err := ownedTag(m, key, rule.ID); err != nil {
				return err
			}
		}
		for _, key := range rule.Clear {
			if err := ownedTag(m, key, rule.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

// ownedTag enforces declaration completeness and the write-to-owned rule,
// under exact byte equality on the key (no folding).
func ownedTag(m Model, key, ruleID string) error {
	tag, ok := m.Tags[key]
	if !ok {
		return fmt.Errorf("unknown tag %q written by rule %q", key, ruleID)
	}
	if tag.Provenance != "owned" {
		return fmt.Errorf("rule %q writes non-owned tag %q (%s)", ruleID, key, tag.Provenance)
	}
	return nil
}

func normalize(model Model) ([]Row, error) {
	var rows []Row
	for _, rule := range model.Rule {
		atoms := map[string]Atom{}
		for _, contextID := range rule.Use {
			if err := mergeContext(model.Context, contextID, atoms, map[string]bool{}); err != nil {
				return nil, err
			}
		}
		mergeAtoms(atoms, "match", rule.Match)
		mergeAtoms(atoms, "all", rule.Guard.All)
		mergeAtoms(atoms, "unless", rule.Guard.Unless)

		// Lift the single `recognized` atom out of the predicate set into
		// the row's outcome field. An `in` atom expands one row per member.
		outcomes, err := liftOutcome(atoms, model, rule.ID)
		if err != nil {
			return nil, err
		}

		set := make([]Atom, 0, len(atoms))
		for _, a := range atoms {
			set = append(set, a)
		}
		slices.SortFunc(set, compareAtoms)

		writes, next := renderWrites(rule)
		requires := requiresOwned(rule)

		kind := "transition"
		if len(rule.Escape) > 0 {
			kind = "escape"
		}

		for _, outcome := range outcomes {
			suffix := ""
			if len(outcomes) > 1 {
				suffix = outcome
			}
			rows = append(rows, Row{
				ModelID:       model.Model.ID,
				RuleID:        rule.ID,
				Suffix:        suffix,
				Locator:       rule.Source,
				Kind:          kind,
				Outcome:       outcome,
				Atoms:         set,
				NextTags:      next,
				Writes:        writes,
				RequiresOwned: requires,
				EscapeClasses: slices.Sorted(slices.Values(rule.Escape)),
			})
		}
	}

	slices.SortFunc(rows, compareRows)
	return rows, nil
}

// compareRows is the dump's row ordering: the identity tuple
// (model id, rule id, expansion suffix), compared field by field,
// byte-lexicographically. The source locator does NOT participate.
func compareRows(a, b Row) int {
	if c := cmp.Compare(a.ModelID, b.ModelID); c != 0 {
		return c
	}
	if c := cmp.Compare(a.RuleID, b.RuleID); c != 0 {
		return c
	}
	return cmp.Compare(a.Suffix, b.Suffix)
}

// compareAtoms sorts a row's predicate atoms by key, block, operator token,
// then literal — the order the dump contract states.
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
	return cmp.Compare(a.Literal, b.Literal)
}

// liftOutcome removes the rule's single `recognized` atom from the predicate
// set and returns the outcome literals it binds.
func liftOutcome(atoms map[string]Atom, model Model, ruleID string) ([]string, error) {
	var found []Atom
	for id, a := range atoms {
		if a.Key == recognizedKey {
			found = append(found, a)
			delete(atoms, id)
		}
	}
	if len(found) != 1 {
		return nil, fmt.Errorf("rule %q binds %d recognized atoms, want exactly 1", ruleID, len(found))
	}
	atom := found[0]
	var literals []string
	switch atom.Operator {
	case "eq":
		literals = []string{atom.Literal}
	case "in":
		literals = strings.Split(atom.Literal, " ")
	default:
		return nil, fmt.Errorf("rule %q binds outcome with operator %q, want eq or in", ruleID, atom.Operator)
	}
	for _, lit := range literals {
		if !slices.Contains(model.Outcomes, lit) {
			return nil, fmt.Errorf("rule %q binds outcome %q outside the alphabet", ruleID, lit)
		}
	}
	slices.Sort(literals)
	return literals, nil
}

// renderWrites returns the rule's writes (owned-tag writes for the accessor
// layer, including rendered <clear> entries) and its next-state tags. An
// escape rule carries neither: RDR 0009 requires a write-free escape row.
func renderWrites(rule Rule) (writes, next []Write) {
	if len(rule.Escape) > 0 {
		return nil, nil
	}
	for key, value := range rule.Write {
		writes = append(writes, Write{Key: key, Value: fmt.Sprint(value)})
	}
	for _, key := range rule.Clear {
		writes = append(writes, Write{Key: key, Clear: true})
	}
	slices.SortFunc(writes, compareWrites)
	next = slices.Clone(writes)
	return writes, next
}

func compareWrites(a, b Write) int {
	return cmp.Compare(a.Key, b.Key)
}

// requiresOwned derives Row.RequiresOwned as the sorted, duplicate-free set
// of tag keys named by the write block and clear list. An escape row carries
// an empty set (RDR 0007 A21).
func requiresOwned(rule Rule) []string {
	if len(rule.Escape) > 0 {
		return nil
	}
	seen := map[string]bool{}
	var keys []string
	for key := range rule.Write {
		if !seen[key] {
			seen[key] = true
			keys = append(keys, key)
		}
	}
	for _, key := range rule.Clear {
		if !seen[key] {
			seen[key] = true
			keys = append(keys, key)
		}
	}
	slices.Sort(keys)
	return keys
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
	mergeAtoms(out, "match", ctx.Match)
	return nil
}

// mergeAtoms folds one authored block into the predicate set, retaining the
// block each atom was authored in. Existence atoms are emitted with the
// kernel's exported constants; any other existence literal is rejected.
func mergeAtoms(out map[string]Atom, block string, source map[string]map[string]any) {
	for key, ops := range source {
		for op, value := range ops {
			literal := renderLiteral(value)
			if op == opExists && literal != literalTrue && literal != literalFalse {
				must(fmt.Errorf("malformed predicate atom: %s.%s has literal %q", key, op, literal))
			}
			out[block+"\x00"+key+"\x00"+op] = Atom{
				Key: key, Operator: op, Literal: literal, Block: block,
			}
		}
	}
}

// renderLiteral gives every literal a byte-comparable spelling. A set
// literal renders as its members joined by a space, sorted, so two authored
// orderings of the same set are one literal.
func renderLiteral(value any) string {
	if list, ok := value.([]any); ok {
		parts := make([]string, 0, len(list))
		for _, item := range list {
			parts = append(parts, fmt.Sprint(item))
		}
		slices.Sort(parts)
		return strings.Join(parts, " ")
	}
	if b, ok := value.(bool); ok {
		if b {
			return literalTrue
		}
		return literalFalse
	}
	return fmt.Sprint(value)
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "refused:", err)
		os.Exit(1)
	}
}
