package table

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/newcoinc/intrastate/internal/resolve"
)

// ClearSentinel is the reserved tag value an explicit `clear` entry
// normalizes to. It is refused wherever a tag value is authored
// (`0002:C11`).
const ClearSentinel = "<clear>"

// RecognizedTagKey is the reserved tag key the freshly recognized outcome
// takes. A `recognized` declaration must be named this, and no owned or
// observed declaration may take the name (`0002:C12`).
const RecognizedTagKey = "recognized"

// suffixSep joins the rule id and each expansion-suffix element in a
// rendered row identity. It is banned from rule ids, outcome literals, and
// match-block `in` members so the rendered identity stays unambiguous
// (`0002:C11`).
const suffixSep = "#"

// Block names the atom block an atom was authored in — a three-valued
// domain. The values spell the kernel's exported constants: JDR 0001 §D12
// widened resolve.Block to three members precisely so this RDR need not
// fork a second vocabulary (`0002:C7`).
type Block = resolve.Block

// The three authored atom blocks, aliased from the kernel.
const (
	BlockMatch  = resolve.BlockMatch
	BlockAll    = resolve.BlockAll
	BlockUnless = resolve.BlockUnless
)

// Kind is the derived row-kind view column. It is computed from the escape
// class list and is never a row field (`0002:C5`).
type Kind string

const (
	// KindTransition is an ordinary transition row.
	KindTransition Kind = "transition"
	// KindEscape is a modeled escape row: a non-empty escape class list.
	KindEscape Kind = "escape"
)

// The two admitted model classes (`0010:C1`). A `class` outside this set
// is `malformed model declaration`; an ABSENT `class` reads as
// ClassStateMachine, which is what makes the empty Model.Class a state
// machine.
const (
	ClassStateMachine  = "state-machine"
	ClassDecisionTable = "decision-table"
)

// IsDecisionTable reports whether m declares the decision-table class.
//
// It is the small helper `0010:C1` licenses for spelling the zero value —
// "a reader MAY spell the zero value through a small helper rather than
// repeating the empty-string comparison". The FIELD stays the storage: this
// reads Model.Class and nothing else, and in particular never re-derives
// the class from `len(owned) == 0`. It is exported because two of the four
// class readers live in `internal/graphlint`.
func IsDecisionTable(m *Model) bool {
	return m != nil && m.Class == ClassDecisionTable
}

// Provenance is a tag declaration's provenance (`0002:C21`).
type Provenance string

const (
	ProvenanceOwned      Provenance = "owned"
	ProvenanceObserved   Provenance = "observed"
	ProvenanceRecognized Provenance = "recognized"
)

// declaredKinds is RDR 0003's closed five-token type-model vocabulary. It
// is 0003's to fix and this RDR's to enforce; `string` is not one of them
// (deviations.md D1).
var declaredKinds = []string{"enum", "bool", "int", "set", "scalar"}

// IsDeclaredKind reports whether kind is one of RDR 0003's five tokens.
func IsDeclaredKind(kind string) bool {
	return slices.Contains(declaredKinds, kind)
}

// operators is RDR 0003's closed operator set. This RDR does not mint
// operators and does not widen the set (`0002:C16`).
var operators = []string{"eq", "in", "lt", "lte", "gt", "gte", "exists", "contains"}

// Operators returns the closed admitted operator set.
func Operators() []string {
	return slices.Clone(operators)
}

// operatorKinds is RDR 0003's operator/kind matrix, mirrored here for the
// same reason `declaredKinds` and `operators` are: RDR 0003 owns the
// matrix, this RDR's loader is what enforces it at load, and `guard`
// imports this package rather than the other way round. `guard.Accepts` is
// the authority and its REQ-4 test pins the matrix cell by cell.
//
// `exists` accepts every kind — it reads presence, not value. The RDR's
// matrix cell adds "provided the tag is declared optional", but that is a
// vacuity report, not a rejection ("well-formed but vacuous — lint reports
// it as such rather than rejecting it"), so it is lint's, not load's.
var operatorKinds = map[string][]string{
	"eq":       {"enum", "bool", "int", "scalar"},
	"in":       {"enum", "bool", "int", "scalar"},
	"lt":       {"int"},
	"lte":      {"int"},
	"gt":       {"int"},
	"gte":      {"int"},
	"exists":   {"enum", "bool", "int", "set", "scalar"},
	"contains": {"set"},
}

// operatorAcceptsKind reports whether the operator/kind matrix admits the
// pair. An operator outside the closed set accepts no kind.
func operatorAcceptsKind(operator, kind string) bool {
	return slices.Contains(operatorKinds[operator], kind)
}

// TagDecl is one tag's declaration: its provenance plus the type model RDR
// 0003 owns, spelled as the wire keys `kind`, `domain`, `min`, `max`,
// `elements`, `single_valued`, and `required` (`0002:C22`). Every declared
// field is carried through normalization without loss.
type TagDecl struct {
	Provenance   Provenance
	Kind         string
	Domain       []string
	Min          *int
	Max          *int
	Elements     []string
	SingleValued bool
	Required     bool
}

// Accessor is one `[read.<id>]`, `[write.<id>]`, or `[gate.<id>]` entry.
// All four of role, path, keys, and timeout are required (`0002:C2`).
type Accessor struct {
	Role     string
	Path     string
	Keys     []string
	Timeout  string
	ReadBack bool
}

// TagValue is one key bound to a member sequence. A set-valued write holds
// its members as a sequence, never a delimiter-joined string (`0002:C10`).
type TagValue struct {
	Key   string
	Value []string
}

// EmitValue is one `[rule.emit]` pair (`0010:C3`).
//
// It is deliberately NOT a TagValue. An emit key is not a tag key: it is
// undeclared, uninterpreted, and compared by exact byte equality, and its
// value is one authored string rather than a member sequence — so the
// `Value []string` a TagValue carries would invite the set semantics,
// domain conformance, and `renderValue` quoting that none of this applies
// to.
type EmitValue struct {
	Key   string
	Value string
}

// Atom is one normalized predicate atom. Its identity is the full tuple
// (key, block, operator token, literal) — the tuple the merge is a set
// over and the dump sorts on — and the literal enters that identity as its
// member SEQUENCE, never as a joined string (`0002:C6`, `0002:C10`).
type Atom struct {
	Key      string
	Block    Block
	Operator string
	Literal  []string
}

// Row is one normalized candidate row.
//
// The atoms travel as ONE unified set with each atom's authored block
// retained; the Match/guard split is applied at the kernel handoff, never
// stored here (`0002:C16`). There is no row-kind field and no positional
// field: kind is derived from the escape list and selection must not be
// decidable from row order (`0002:C5`, `0002:C20`).
type Row struct {
	// ModelID and RuleID plus Suffix are the identity tuple the dump sorts
	// on. SourceLocator is diagnostic and never participates in ordering.
	ModelID       string
	RuleID        string
	Suffix        []string
	SourceLocator string
	// Source is the authored `[[rule]].source` provenance annotation,
	// carried alongside the locator and never used to derive it.
	Source string

	// Outcome is the single outcome the row binds, lifted out of the
	// predicate set at normalization.
	Outcome string
	// Atoms is the unified predicate set, sorted by the identity tuple.
	Atoms []Atom
	// Gate names the gate accessors the executor requires before applying
	// this row's plan.
	Gate []string
	// NextTags and Writes are populated independently from the rule; one
	// is never an alias of the other (`0002:C15`).
	NextTags []TagValue
	Writes   []TagValue
	// RequiresOwned is derived, never authored: the sorted, duplicate-free
	// key set the write block and clear list name (`0002:C14`).
	RequiresOwned []string
	// Escape lists the modeled resolver failure classes. A non-empty list
	// is what discriminates an escape row.
	Escape []string
	// Emit is the authored `[rule.emit]` block, key-sorted (`0010:C3`).
	//
	// It is table data and never crosses the kernel seam: KernelRow does
	// not carry it, and the CLI joins it back by rule id AFTER selection.
	// An absent block and a present-but-empty one both normalize to the
	// empty sequence — `emit` keys on LENGTH, not on key presence, which is
	// the deliberate divergence from the write/clear/gate blocks.
	//
	// The sequence MAY be shared across the rows one rule expands to: it is
	// read by the dump column and by the resolve payload join, and both are
	// readers, so no defensive copy is taken.
	Emit []EmitValue

	// setKeys names the tag keys whose declared kind is `set`, sorted.
	//
	// It is unexported and carries no authored form: the normalized value
	// holds every literal as a member SEQUENCE regardless of kind
	// (`0002:C10`), and this is only how the KERNEL HANDOFF knows to spell
	// a one-member set as a one-element JSON array rather than as the bare
	// member (JDR 0001 §D13). It is derived from the model's declarations,
	// so it is identical across the semantics-preserving permutations the
	// Round-Trip invariant compares.
	setKeys []string
}

// isSet reports whether key's declared kind is `set`.
func (r Row) isSet(key string) bool {
	return slices.Contains(r.setKeys, key)
}

// Identity renders the identity tuple (model id, rule id, expansion
// suffix), joining the rule id and each suffix element with `#`. Every
// string that can reach here has `#` banned at load, so the rendering is
// unambiguous.
func (r Row) Identity() string {
	var b strings.Builder
	b.WriteString(r.ModelID)
	b.WriteString(".")
	b.WriteString(r.RuleID)
	for _, s := range r.Suffix {
		b.WriteString(suffixSep)
		b.WriteString(s)
	}
	return b.String()
}

// Kind is the derived row-kind view property, discriminated solely by a
// non-empty escape class list (`0002:C5`).
func (r Row) Kind() Kind {
	if len(r.Escape) > 0 {
		return KindEscape
	}
	return KindTransition
}

// KernelRow constructs the resolve.Row for the kernel, applying the
// match/guard split. Routing is by BLOCK and never by operator: an `eq`
// atom under `guard.all` is a guard atom. The routing is exhaustive and
// disjoint by construction (`0002:C16`).
func (r Row) KernelRow() resolve.Row {
	var match []resolve.Tag
	var guard []resolve.GuardAtom
	for _, a := range r.Atoms {
		if a.Block == BlockMatch {
			match = append(match, resolve.Tag{Key: a.Key, Value: seamValue(a.Literal, r.isSet(a.Key))})
			continue
		}
		guard = append(guard, resolve.GuardAtom{
			Key:      a.Key,
			Operator: a.Operator,
			Literal:  seamValue(a.Literal, setValuedLiteral(a.Operator)),
			Block:    a.Block,
		})
	}

	escape := make([]resolve.RefusalKind, 0, len(r.Escape))
	for _, e := range r.Escape {
		escape = append(escape, resolve.RefusalKind(e))
	}

	return resolve.Row{
		RuleID:        r.RuleID,
		SourceLocator: r.SourceLocator,
		Outcome:       r.Outcome,
		Match:         match,
		RequiresOwned: slices.Clone(r.RequiresOwned),
		Guard:         guard,
		NextTags:      r.seamTags(r.NextTags),
		Writes:        r.seamTags(r.Writes),
		Escape:        escape,
	}
}

// setValuedLiteral reports whether an operator's right-hand side is a
// member SET rather than a single value, which decides the literal's
// carriage form independently of the tag's declared kind.
//
// `0002:C17` fixes the split: "`eq`, `in`, and `contains` take members and
// are domain-checked" — but RDR 0003 fences `eq` and the comparisons as
// SINGLE-VALUE operators ("the tag's single held value IS the literal"),
// leaving `in` (membership in a literal set) and `contains` (set
// containment) as the two whose literal is a set. The shipped cross-RDR
// contract test agrees byte for byte:
// `internal/resolve/guardcontract.go` drives `{"in", ["alpha","beta"],
// "alpha"}` and `{"contains", ["alpha"], ["alpha","beta"]}`, so an `in` or
// `contains` literal MUST cross as the §D13 array whatever the tag's kind.
//
// Getting this from the declared kind alone truncated the literal to
// `members[0]`, and on the `in` arm that DEADLOCKED rather than refused:
// RDR 0007's evaluator parses the literal as a JSON array, must answer
// `GuardUnevaluable` on a bare string, and the unevaluable-candidate veto
// then blocks decidable siblings.
func setValuedLiteral(operator string) bool {
	return operator == "in" || operator == "contains"
}

// seamValue encodes a member sequence for the kernel seam, where
// resolve.Tag.Value and GuardAtom.Literal are both `string`.
//
// JDR 0001 §D13 fixes the byte form: a SET crosses as its canonical JSON
// array — members sorted, duplicate-free, compact encoding — and read-back
// equality (RDR 0004) is byte equality over that array. RDR 0007 already
// wrote its `contains` leg against §D13, so this is the form it matches
// rather than a second encoding.
//
// Two things make a value a set here, and they answer different questions.
// For a tag VALUE — `Tag.Value` on a match tag, a next-state tag, or a
// write — set-ness is the DECLARED kind, not the member count: a one-member
// set still crosses as a one-element array, or read-back could not tell
// `["a"]` from the scalar `a`. For a guard atom LITERAL it is the OPERATOR
// ALONE, per setValuedLiteral above: the kind says what the TAG holds, and
// the literal is the right-hand side, not the tag. An `exists` literal is a
// bare bool constant and a comparison bound is a bare single value however
// the tag is declared, and the kernel compares both verbatim.
//
// A non-set value crosses as its single member verbatim, which is what
// keeps a `<clear>` write the bare sentinel the kernel and RDR 0004 expect.
//
// The normalized value never holds this form: it stays a member sequence,
// and no identity is derived from the serialization (`0002:C10`).
func seamValue(members []string, isSet bool) string {
	if !isSet {
		if len(members) == 0 {
			return ""
		}
		// A non-set value is one member by construction: load refuses a
		// multi-member literal on a single-value operator and a
		// multi-member value on a non-`set` kind, at every authoring site.
		// The seam is where a surplus member would vanish unreported, so
		// the refusals are upstream of it rather than a truncation here.
		return members[0]
	}
	canonical := slices.Compact(slices.Sorted(slices.Values(members)))

	// HTML escaping is DISABLED, so `<`, `>`, and `&` serialize as
	// themselves. Bare `json.Marshal` escapes them, which would make this
	// seam spell a set value differently from every other site that emits or
	// compares one — and §D13 read-back equality is BYTE equality, so a
	// member carrying `&` would fail to match the request that wrote it
	// (JDR 0001 §D13; RDR 0005 REQ-70/REQ-71/REQ-72 requires one encoder on
	// every path a set value crosses).
	var buf strings.Builder
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(canonical); err != nil {
		// Encoding a []string cannot fail.
		return ""
	}
	// Encode appends a newline; the seam value carries none.
	return strings.TrimRight(buf.String(), "\n")
}

func (r Row) seamTags(in []TagValue) []resolve.Tag {
	if len(in) == 0 {
		return nil
	}
	out := make([]resolve.Tag, 0, len(in))
	for _, t := range in {
		// A rendered clear is the sentinel itself whatever the tag's kind:
		// it names a deletion, not a set value (`0002:C23`).
		isSet := r.isSet(t.Key) && !isClear(t.Value)
		out = append(out, resolve.Tag{Key: t.Key, Value: seamValue(t.Value, isSet)})
	}
	return out
}

// isClear reports whether a rendered value is the `<clear>` sentinel.
func isClear(value []string) bool {
	return len(value) == 1 && value[0] == ClearSentinel
}

// Model is the normalized transition model: the decoded declarations plus
// the deterministic candidate-row set normalization produced.
type Model struct {
	ID          string
	Version     int
	Description string
	// Class is the declared model class (`0010:C1`). Its zero value — the
	// empty string — reads everywhere as ClassStateMachine, so a
	// hand-constructed Model is a state machine without naming a class.
	//
	// The FIELD is the storage: no reader re-derives the class from the
	// owned set, and nothing downstream infers it.
	Class string
	// Metadata is the one sanctioned extension namespace, carried through
	// untouched and never interpreted (`0002:C2`).
	Metadata map[string]any

	Outcomes []string
	// Initial is the root assignment set.
	Initial []TagValue
	// Terminal carries the stop set as DEREFERENCED predicate sets over
	// the contexts the root `terminal` list names — never the bare ids
	// (`0002:C2`).
	Terminal [][]Atom

	Tags    map[string]TagDecl
	Readers map[string]Accessor
	Writers map[string]Accessor
	Gates   map[string]Accessor

	// DumpOrder is presentation only and never reaches a candidate row.
	DumpOrder []string

	// Rows is the normalized candidate-row set, pre-sorted by the identity
	// tuple.
	Rows []Row
}

// KernelTable builds the resolve.Table this model hands the kernel. It
// adds no runtime ordering and no host-code predicate callback: the guard
// travels as parsed atoms (`0002:IP` Phase 4).
func (m *Model) KernelTable() resolve.Table {
	rows := make([]resolve.Row, 0, len(m.Rows))
	for _, r := range m.Rows {
		rows = append(rows, r.KernelRow())
	}
	return resolve.Table{
		Revision: m.ID,
		Outcomes: slices.Clone(m.Outcomes),
		Rows:     rows,
	}
}

// CheckModelIDs reports a duplicate model id over a SET of loaded models.
//
// `duplicate model id` is never decidable within one document — one source
// document carries exactly one model — so the check is scoped to a caller
// loading several documents into one dump or lint invocation, and must run
// before rows are merged for rendering. A loader handed one document never
// reports it (`0002:C19`).
func CheckModelIDs(models []*Model) error {
	seen := make(map[string]bool, len(models))
	for _, m := range models {
		if m == nil {
			continue
		}
		if seen[m.ID] {
			return fail(CatDuplicateModelID, "model id "+m.ID+" is carried by more than one document")
		}
		seen[m.ID] = true
	}
	return nil
}
