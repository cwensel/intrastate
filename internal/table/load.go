package table

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	toml "github.com/pelletier/go-toml/v2"

	"github.com/cwensel/intrastate/internal/resolve"
)

// Load parses, validates, and normalizes one source document into its
// candidate-row value.
//
// It takes already-read BYTES plus a source id, never a filesystem path:
// this package performs no file I/O and no path resolution (`0002:EIA`).
//
// "Load" names the whole source-to-candidate-rows pipeline. Every category
// is refused before the pipeline yields rows, and the refusal is singular:
// load is fail-fast and returns one categorized error, never a list
// (`0002:C24`, `0002:C3`).
//
// The two fixed precedences are `malformed TOML` -> version gate -> strict
// decoding -> the remaining categories. Beyond those, the order in which
// independent defects are checked is deliberately unspecified.
func Load(src []byte, sourceID string) (*Model, error) {
	// Pass 1: the version gate, permissive enough to read `[model]` before
	// strict field validation ever runs. A v2-shaped document carrying a
	// v2-only key must refuse `unsupported_version`, not
	// `unknown_schema_field`.
	var probe versionProbe
	if err := toml.Unmarshal(src, &probe); err != nil {
		return nil, fail(CatMalformedTOML, err.Error())
	}
	if probe.Model == nil {
		return nil, fail(CatMalformedModelDeclaration, "no [model] table")
	}
	if probe.Model.Version == nil {
		// An absent version refuses here and is never folded into
		// `unsupported_version`, which would name a version the author
		// never wrote (`0002:C24`).
		return nil, fail(CatMalformedModelDeclaration, "[model] carries no version")
	}
	if *probe.Model.Version != 1 {
		return nil, fail(CatUnsupportedVersion,
			fmt.Sprintf("version %d; this RDR accepts version 1 only", *probe.Model.Version))
	}

	// Pass 2: strict decoding of the closed layout.
	var doc sourceDoc
	if err := decodeStrict(src, &doc); err != nil {
		return nil, err
	}

	ld := &loader{doc: &doc, src: src, sourceID: sourceID}
	return ld.run()
}

// loader carries one document through validation and normalization.
type loader struct {
	doc *sourceDoc
	// src is the document's own bytes, kept for the ONE thing the decoded
	// structs cannot supply: where in the source a declaration was written.
	src      []byte
	sourceID string

	model *Model
	// contexts holds each context's own authored atoms, before inheritance.
	contexts map[string][]Atom
	// resolved memoizes the transitive atom set per context id.
	resolved map[string][]Atom
}

func (l *loader) run() (*Model, error) {
	l.model = &Model{}

	for _, step := range []func() error{
		l.loadModelHeader,
		l.loadOutcomes,
		l.loadTags,
		// RDR 0024 `0024:C2` — the two emit steps, in this order, at this
		// position: immediately after loadTags and ahead of normalizeRules,
		// so both read the SOURCE rules rather than normalized rows and
		// refuse before any row is yielded. loadEmitDecls builds the
		// carrier checkRuleEmit reads, so the order is load-bearing for
		// correctness, not only for reporting: the value check is
		// safe-by-omission and depends on C1's arms having already fired.
		l.loadEmitDecls,
		l.checkRuleEmit,
		l.loadAccessors,
		l.loadDump,
		l.loadContexts,
		l.loadInitial,
		l.loadTerminal,
		l.normalizeRules,
		l.checkClassAgreement,
		l.checkAccessorBindings,
	} {
		if err := step(); err != nil {
			return nil, err
		}
	}
	return l.model, nil
}

// ---------------------------------------------------------------- header

func (l *loader) loadModelHeader() error {
	m := l.doc.Model
	if m == nil {
		return fail(CatMalformedModelDeclaration, "no [model] table")
	}
	if m.ID == nil || *m.ID == "" {
		return fail(CatMalformedModelDeclaration, "[model] carries no id")
	}
	if m.Version == nil {
		return fail(CatMalformedModelDeclaration, "[model] carries no version")
	}

	// The class is `[model]` data, so it is READ here — but the agreement
	// with the owned set is checked later, in a step at or after loadTags,
	// where the tag table exists (`0010:C1`).
	if m.Class != nil {
		switch *m.Class {
		case ClassStateMachine, ClassDecisionTable:
		default:
			return fail(CatMalformedModelDeclaration,
				"[model] class "+strconv.Quote(*m.Class)+
					" is not "+strconv.Quote(ClassStateMachine)+
					" or "+strconv.Quote(ClassDecisionTable))
		}
		l.model.Class = *m.Class
	}

	l.model.ID = *m.ID
	l.model.Version = *m.Version
	l.model.Description = m.Description
	// Carried through untouched: no key inside is interpreted, and its
	// internal shape is deliberately unconstrained (`0002:C2`).
	l.model.Metadata = m.Metadata
	return nil
}

func (l *loader) loadOutcomes() error {
	if l.doc.Outcomes == nil {
		return fail(CatMissingRecognizedOutcomeAlphabet, "no root `outcomes` alphabet")
	}
	alphabet := *l.doc.Outcomes
	if len(alphabet) == 0 {
		return fail(CatMalformedRecognizedOutcomeAlphabet, "the alphabet is empty")
	}
	seen := make(map[string]bool, len(alphabet))
	for _, o := range alphabet {
		switch {
		case o == "":
			return fail(CatMalformedRecognizedOutcomeAlphabet, "the empty string is a member")
		case strings.Contains(o, suffixSep):
			return fail(CatMalformedRecognizedOutcomeAlphabet,
				"member "+o+" contains the expansion-suffix separator "+suffixSep)
		case seen[o]:
			return fail(CatMalformedRecognizedOutcomeAlphabet, "member "+o+" is duplicated")
		}
		seen[o] = true
	}
	l.model.Outcomes = alphabet
	return nil
}

// ---------------------------------------------------- tags and accessors

// tagHeaderLine reports the 1-based line of the `[tags.<key>]` header
// declaring `key`, or ZERO where it cannot be identified.
//
// The scan is textual and deliberately narrow. It recognizes exactly the
// standard-table authoring, `[tags.<key>]` on a line of its own, with the
// key either bare or quoted. Every other authoring TOML permits — an inline
// table under `[tags]`, a dotted key, a key whose own text contains a
// bracket — falls outside it and yields zero.
//
// Zero on anything ambiguous is the point, not a gap. A line number is a
// pointer the author follows to source text; pointing at the wrong text
// costs more than pointing at no text, and the caller renders zero as a
// file-only locator that is honest about what it does not know. So a key
// matched zero times OR more than once declines: a second match means the
// scan cannot tell which header the refusal belongs to.
func tagHeaderLine(src []byte, key string) int {
	bare := "[tags." + key + "]"
	quoted := "[tags." + strconv.Quote(key) + "]"

	line, matches := 0, 0
	for i, raw := range strings.Split(string(src), "\n") {
		// Trailing comments and surrounding space are the only decoration a
		// header line may carry; anything else is not this header.
		text := raw
		if hash := strings.IndexByte(text, '#'); hash >= 0 {
			text = text[:hash]
		}
		text = strings.TrimSpace(text)
		if text != bare && text != quoted {
			continue
		}
		matches++
		line = i + 1
	}
	if matches != 1 {
		return 0
	}
	return line
}

func (l *loader) loadTags() error {
	decls := make(map[string]TagDecl, len(l.doc.Tags))
	for key, src := range l.doc.Tags {
		decl, err := tagDecl(key, src)
		if err != nil {
			// The declaration's own `[tags.<key>]` header is the only
			// position the decoded struct leaves recoverable, and this is
			// the only place that still holds both the key and the source.
			return atLine(err, tagHeaderLine(l.src, key))
		}
		decls[key] = decl
	}

	// A `recognized` declaration MUST be named `recognized`, and no owned
	// or observed declaration may take that name (`0002:C12`).
	//
	// RDR 0008 `0008:C3` — each violation carries the three-field payload:
	// the offending name as authored, the direction's remedy name, and the
	// direction's stable rule identifier.
	//
	// The scan runs over SORTED keys, not Go map order. The two directions
	// prescribe opposite remedies ("rename to `recognized`" vs "rename away
	// from it"), so a model breaching both at once must not report a payload
	// chosen by the map seed: identical bytes would yield contradictory
	// advice across runs, and `0008:C3`'s rule identifier is a token a golden
	// test asserts byte-for-byte and a consumer uses for remediation lookup.
	// Which of two SAME-direction declarations is named stays unspecified
	// (`0008:TS-5`); sorting only makes that choice reproducible.
	for _, key := range slices.Sorted(maps.Keys(decls)) {
		decl := decls[key]
		if decl.Provenance == ProvenanceRecognized && key != RecognizedTagKey {
			// Rename TO the reserved key: the remedy names it.
			return &Failure{
				Category:  CatReservedTagKey,
				Detail:    "tag " + key + " declares provenance recognized under another name",
				Offending: key,
				Remedy:    RecognizedTagKey,
				Rule:      RuleKernelOwned,
			}
		}
		if key == RecognizedTagKey && decl.Provenance != ProvenanceRecognized {
			// Rename AWAY FROM the reserved key: the remedy is the empty
			// string, because "choose any other name" has no one answer and a
			// renderer must not present the reserved key as the required one.
			return &Failure{
				Category:  CatReservedTagKey,
				Detail:    "tag " + RecognizedTagKey + " is reserved for the recognized declaration",
				Offending: key,
				Remedy:    "",
				Rule:      RuleAuthorMustRename,
			}
		}
	}
	if _, ok := decls[RecognizedTagKey]; !ok {
		return fail(CatMalformedModelDeclaration, "no [tags.recognized] declaration")
	}

	l.model.Tags = decls
	return nil
}

// ------------------------------------------- the declared emit vocabulary

// emitKinds is RDR 0024's closed emit-kind vocabulary (`0024:C1`): RDR
// 0003's token spellings reused verbatim, minus `set`. `set` is excluded
// because an emit value is ONE authored string, never a member sequence
// (`0010:C3`).
var emitKinds = []string{"enum", "bool", "int", "scalar"}

// loadEmitDecls runs C1's grammar checks over `[emit]` and builds C3's
// carrier (`0024:C2`).
//
// The opt-in gate is evaluated HERE rather than by the caller: the step
// slice is uniform and holds bound method values, so there is no room for a
// conditional call. The trigger is the COUNT of declared keys, never the
// presence of the table — a bare `[emit]` with zero sub-tables is a
// zero-declaration model, identical in every observable to omitting it.
//
// The carrier is set unconditionally, so `Model.EmitDecls` is non-nil after
// every successful load and empty rather than nil for a zero-declaration
// model, matching the `Model.Tags` convention.
//
// The walk is over SORTED keys for the same reason `loadTags`'s
// reserved-key scan is: load is fail-fast and reports one refusal, so which
// of several malformed declarations is named must not be chosen by the map
// seed.
func (l *loader) loadEmitDecls() error {
	l.model.EmitDecls = make(map[string]EmitDecl, len(l.doc.Emit))
	if len(l.doc.Emit) == 0 {
		return nil
	}

	for _, key := range slices.Sorted(maps.Keys(l.doc.Emit)) {
		decl, err := emitDecl(key, l.doc.Emit[key])
		if err != nil {
			// A declaration defect keys on the declaration's own
			// `[emit.<key>]` header, which `tagHeaderLine`'s technique
			// reaches unchanged. This is the only place still holding both
			// the key and the source bytes.
			return atLine(err, emitHeaderLine(l.src, key))
		}
		l.model.EmitDecls[key] = decl
	}
	return nil
}

// emitDecl validates one `[emit.<key>]` declaration and returns its carried
// form (`0024:C1`).
func emitDecl(key string, src sourceEmitDecl) (EmitDecl, error) {
	bad := func(detail string) (EmitDecl, error) {
		return EmitDecl{}, fail(CatMalformedEmitDeclaration, "emit "+key+": "+detail)
	}

	if !slices.Contains(emitKinds, src.Kind) {
		return bad("kind " + strconv.Quote(src.Kind) +
			" is not one of " + strings.Join(emitKinds, ", "))
	}
	if src.Kind != "enum" {
		// None of bool, int, and scalar takes a domain, in EITHER spelling:
		// the arm refuses the PRESENCE of the key, since the narrower
		// reading would leave `[emit.<key>.domain]` under `kind = "bool"`
		// silently admitted.
		if src.Domain != nil {
			return bad("kind " + src.Kind + " admits no domain; only an enum declares one")
		}
		// Domain stays nil, never an empty non-nil slice: the distinction
		// is observable, because the carrier is asserted by value-equality.
		return EmitDecl{Kind: src.Kind}, nil
	}

	members, dispositions, err := emitDomain(src.Domain)
	if err != nil {
		return bad(err.Error())
	}
	// "No usable domain" is ONE arm, not two: an enum declaring no `domain`
	// key at all and one carrying an empty `[emit.<key>.domain]` sub-table
	// BOTH decode to a nil domain, with nothing to discriminate on. An
	// empty flat `domain = []` is distinguishable but takes the same arm.
	if len(members) == 0 {
		return bad("kind enum carries no usable domain")
	}

	seen := make(map[string]bool, len(members))
	for _, m := range members {
		if m == "" {
			return bad("the empty string is a domain member")
		}
		if seen[m] {
			// Duplicates ACROSS disposition lists included: one member, one
			// disposition.
			return bad("domain member " + strconv.Quote(m) + " is duplicated")
		}
		seen[m] = true
	}

	// The carry is value-preserving, not order-preserving: `Domain` is the
	// union sorted bytewise, so a partitioned domain yields one value
	// whatever order the decoder's map iteration hands the dispositions
	// over. The partition grouping is not carried separately — it is fully
	// recoverable from `Dispositions`.
	slices.Sort(members)
	return EmitDecl{Kind: "enum", Domain: members, Dispositions: dispositions}, nil
}

// emitDomain reads either authored spelling of the one `domain` key: a flat
// member array, or a sub-table whose keys are model-authored disposition
// tokens and whose values are member arrays (`0024:C1`).
//
// It returns the members in AUTHORED order — the caller sorts — plus the
// member→disposition map, which is nil for the flat spelling.
//
// The shape arms here exist only because strictness cannot reach them: the
// `any` field type is C1's carve-out for the two spellings, and it is what
// lets a non-array domain, a non-array disposition value, nesting below the
// disposition level, and a non-string member arrive as decoded values
// rather than as decoder errors.
func emitDomain(domain any) ([]string, map[string]string, error) {
	switch d := domain.(type) {
	case nil:
		return nil, nil, nil
	case []any:
		members, err := emitMembers(d)
		if err != nil {
			return nil, nil, err
		}
		return members, nil, nil
	case map[string]any:
		var members []string
		dispositions := map[string]string{}
		// Sorted, so a malformed disposition list is named reproducibly.
		for _, token := range slices.Sorted(maps.Keys(d)) {
			if token == "" {
				return nil, nil, fmt.Errorf("the empty string is a disposition token")
			}
			list, ok := d[token].([]any)
			if !ok {
				return nil, nil, fmt.Errorf("disposition %s is not an array of members",
					strconv.Quote(token))
			}
			listed, err := emitMembers(list)
			if err != nil {
				return nil, nil, err
			}
			for _, m := range listed {
				dispositions[m] = token
			}
			members = append(members, listed...)
		}
		if len(members) == 0 {
			return nil, nil, nil
		}
		return members, dispositions, nil
	default:
		return nil, nil, fmt.Errorf(
			"domain is neither a flat array of members nor a table of dispositions")
	}
}

// emitMembers reads one authored member array. Every member is a STRING: a
// non-string member is refused on its type, and any nesting below the
// disposition level arrives here as a non-string too.
func emitMembers(list []any) ([]string, error) {
	out := make([]string, 0, len(list))
	for _, item := range list {
		m, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("domain member %v is not a string", item)
		}
		out = append(out, m)
	}
	return out, nil
}

// checkRuleEmit holds every rule's authored `[rule.emit]` block to the
// declared vocabulary (`0024:C2`).
//
// Two cross-checks, over the SOURCE rules rather than normalized rows,
// because this step runs ahead of `normalizeRules`:
//
//   - `unknown_emit_key` — any key of any rule, ordinary or escape, not
//     declared under `[emit]`. Strictness is whole-model, not per-key.
//   - `emit_value_out_of_domain` — an authored value outside its key's
//     declared enum domain, not a `bool` token, or not an `int` literal.
//     A `scalar` value is never refused.
//
// The opt-in gate is evaluated here for the same reason it is in
// loadEmitDecls: with zero declarations neither refusal is reachable and
// the pipeline is byte-for-byte today's.
func (l *loader) checkRuleEmit() error {
	if len(l.model.EmitDecls) == 0 {
		return nil
	}

	for ordinal, rule := range l.doc.Rule {
		ruleID := ""
		if rule.ID != nil {
			ruleID = *rule.ID
		}
		// This step runs AHEAD of `normalizeRules` because `0024:C2` fixes
		// it there ("before yielding rows"), and `normalizeRules` is where
		// the missing-id, malformed-id and duplicate-id categories are
		// minted. So a rule id read here is not yet proven present, and the
		// refusal names the rule POSITIONALLY when it is absent rather than
		// rendering "rule  emits …", which attributes the defect to nothing.
		// The ordering C2 fixes is untouched; only the naming degrades.
		named := "rule " + ruleID
		if ruleID == "" {
			named = "rule #" + strconv.Itoa(ordinal+1)
		}
		// Sorted, so which of several defects in one block is reported does
		// not ride the map seed.
		for _, key := range slices.Sorted(maps.Keys(rule.Emit)) {
			decl, declared := l.model.EmitDecls[key]
			if !declared {
				return atLine(fail(CatUnknownEmitKey,
					named+" emits the undeclared key "+key),
					emitRuleLine(l.src, ruleID, ordinal))
			}
			// A scalar key is short-circuited BEFORE the call rather than
			// relying on conformKind's fall-through: the escape hatch is
			// total, and the skip says so at the call site.
			if decl.Kind == "scalar" {
				continue
			}
			// The throwaway TagDecl literal IS the adapter (`0024:C3`): no
			// extraction, no new exported helper. Min, Max and Elements stay
			// nil, so the int arm's bounds and the set arm never fire, and
			// the reuse is safe-by-omission — which is why C1's arms must
			// have fired first.
			if err := ConformValue(
				TagDecl{Kind: decl.Kind, Domain: decl.Domain}, rule.Emit[key]); err != nil {
				return atLine(fail(CatEmitValueOutOfDomain,
					named+" emits "+key+" = "+
						strconv.Quote(rule.Emit[key])+": "+err.Error()),
					emitRuleLine(l.src, ruleID, ordinal))
			}
		}
	}
	return nil
}

// emitHeaderLine reports the 1-based line of the `[emit.<key>]` header
// declaring key, or ZERO where it cannot be identified.
//
// It is `tagHeaderLine`'s technique against a different prefix, and it
// declines on anything ambiguous for the same reason: pointing at the wrong
// text costs more than pointing at no text.
func emitHeaderLine(src []byte, key string) int {
	return headerLine(src, "[emit."+key+"]", "[emit."+strconv.Quote(key)+"]")
}

// emitRuleLine reports the 1-based line locating the rule that authored the
// offending emit block: its `id` assignment where one is authored, and
// otherwise its `[[rule]]` header.
//
// A rule-side defect keys on the offending RULE rather than on its
// `[rule.emit]` header, because a rule's emit block is not identifiable on
// its own: `[rule.emit]` is spelled identically under every `[[rule]]`, so
// a scan for it cannot tell which rule it belongs to (`0024:C2`).
//
// REQ-32/REQ-76 fix the technique as a rule-id-anchored, BLOCK-BOUNDED
// FORWARD SCAN, and both halves carry weight:
//
//   - Block-bounded. The scan starts at the ordinal's `[[rule]]` header and
//     stops at the next top-level table header, so it can only ever read
//     the offending rule's own text. That is what closes the `[model]`
//     id-collision hazard C2 raises: `[model]`'s `id` is outside every rule
//     block, so it is never a candidate, and a rule id that happens to
//     equal the model id still recovers its own line. It also closes the
//     duplicate-id collision, which C2 assumed away on the premise that
//     `CatDuplicateRuleID` runs first — false as built, since C2 itself
//     places both emit steps AHEAD of `normalizeRules`, where the
//     missing-id, malformed-id and duplicate-id categories are enforced.
//     Ordinal selection makes the premise unnecessary rather than true: the
//     step ordering C2 fixes is normative and is left exactly as specified.
//   - Forward scan, not text equality. The anchor TOKENIZES the assignment
//     rather than matching one canonical spelling. TOML fixes no spelling
//     for `id = "r"`: `id="r"`, `id  =  "r"` and the literal-string form
//     `id = 'r'` are the same document, and an id carrying `#` must not be
//     truncated by a comment strip that cuts inside a quoted string.
//
// `ordinal` is the rule's index in `l.doc.Rule`, which the decoder fills in
// document order, so it maps to the ordinal `[[rule]]` header. `ruleID` is
// used only to CONFIRM the match; a rule that authors no id, or whose id
// does not round-trip through the scan, still yields its header line rather
// than zero, because `0024:C2` requires a source line unconditionally.
func emitRuleLine(src []byte, ruleID string, ordinal int) int {
	lines := strings.Split(string(src), "\n")

	starts := make([]int, 0, 8)
	for i, raw := range lines {
		if tableHeader(raw) == "[[rule]]" {
			starts = append(starts, i)
		}
	}
	// `ordinal` indexes `l.doc.Rule`, so it is only meaningful while the
	// scanned census aligns with what the decoder found. A spelling this
	// scan fails to recognize would shift later ordinals onto the wrong
	// block, so an out-of-range ordinal reports NO line rather than a
	// confidently wrong one.
	if ordinal < 0 || ordinal >= len(starts) {
		return 0
	}

	start := starts[ordinal]
	end := len(lines)
	for i := start + 1; i < len(lines); i++ {
		if tableHeader(lines[i]) != "" {
			end = i
			break
		}
	}

	// Forward scan from the block's header to its close. A rule's own
	// sub-tables (`[rule.match…]`, `[rule.emit]`) are not top-level headers,
	// so they do not end the block; the next `[[rule]]` or `[…]` does.
	for i := start + 1; i < end; i++ {
		value, ok := scalarAssignment(lines[i], "id")
		if !ok {
			continue
		}
		if ruleID == "" || value == ruleID {
			return i + 1
		}
	}
	// No usable `id` assignment inside the block. The `[[rule]]` header is
	// still an honest, unambiguous pointer at the offending rule, and a
	// header line beats no line at all.
	return start + 1
}

// tableHeader reports the bracketed TOML table header a line declares, or
// the empty string where the line declares none. Only TOP-LEVEL headers
// count: a dotted header such as `[rule.emit]` belongs to the block it sits
// in and does not close it, so it reports empty.
//
// The reported header is CANONICAL: TOML fixes no single spelling for a
// key. Whitespace inside the brackets is free, and a bare key, a basic
// string key and a literal string key are interchangeable — `[[rule]]`,
// `[[ rule ]]`, `[["rule"]]` and `[['rule']]` all name the same array
// table and must compare equal. Reporting the authored spelling instead
// would drop the other forms from a caller's `[[rule]]` census and shift
// every later ordinal onto the wrong block — the "pointing at the wrong
// text costs more than pointing at no text" failure `emitRuleLine`
// disclaims, and worse than the bare `:1` `0024:MVV` step 2 forbids,
// because it sends an author to edit an innocent rule.
//
// Only a MATCHED outer quote pair is stripped, and only from an otherwise
// bare key. Folding anything looser would fuse a genuinely distinct table
// name into `rule` and add a PHANTOM census entry, shifting ordinals the
// other way.
func tableHeader(raw string) string {
	text := strings.TrimSpace(stripComment(raw))
	if !strings.HasPrefix(text, "[") || !strings.HasSuffix(text, "]") {
		return ""
	}
	name := unquoteKey(strings.TrimSpace(strings.Trim(text, "[]")))
	// The dot test runs AFTER unquoting, so it is applied to the key's
	// CONTENT rather than to its punctuation. `[["rule.other"]]` is
	// therefore excluded exactly as `[rule.other]` is. That is the
	// conservative direction: TOML reads the quoted dot as one ordinary
	// key, but this scan only needs to know which lines end a `[[rule]]`
	// block, and declining a name it cannot confidently reduce to `rule`
	// keeps it out of the census rather than adding a phantom entry.
	if strings.Contains(name, ".") {
		return ""
	}
	if strings.HasPrefix(text, "[[") && strings.HasSuffix(text, "]]") {
		return "[[" + name + "]]"
	}
	return "[" + name + "]"
}

// unquoteKey strips a MATCHED outer quote pair off a single TOML key,
// reporting the key's content. TOML admits three interchangeable spellings
// of the same key — bare (`rule`), basic string (`"rule"`) and literal
// string (`'rule'`) — so canonicalizing them to the bare form is what lets
// one string comparison recognize a table however it was authored.
//
// It is deliberately narrow. A lone or mismatched quote is NOT a quoted
// key and is reported unchanged, so a table genuinely named `"rule` can
// never be folded into the `rule` census.
//
// A basic-string key is DECODED, not merely unwrapped: `"rule"` names
// the array table `rule` exactly as `"rule"` does, so reporting the raw
// body would drop the block from the caller's census and shift every later
// ordinal — the same failure the plain quoted spelling caused before it was
// canonicalized. Decoding runs through `decodeScalarString`, the shared
// helper `scalarAssignment` uses for the value side, so the key half and
// the id half of the anchor agree by construction rather than by two
// parallel implementations.
func unquoteKey(name string) string {
	if decoded, ok := decodeScalarString(name); ok {
		return decoded
	}
	return name
}

// scalarAssignment reads a bare `key = <string>` assignment off one line and
// reports the value with its quoting removed.
//
// It tolerates every spelling TOML admits for the same document: arbitrary
// horizontal whitespace around the key and the `=`, and all three string
// syntaxes — basic (`"v"`), literal ('v'), and their multi-line forms'
// single-line use. It reports false for a dotted key, an inline table, or
// any non-string value, none of which name a rule the way an `id` does.
//
// The KEY and `=` are tokenized here; the VALUE is handed to
// `decodeScalarString`, so the comparison happens in the DECODED domain and
// an id authored with an escape matches the string the decoder actually put
// in `l.doc.Rule`.
func scalarAssignment(raw, key string) (string, bool) {
	text := strings.TrimSpace(stripComment(raw))
	rest, ok := strings.CutPrefix(text, key)
	if !ok {
		return "", false
	}
	rest = strings.TrimLeft(rest, " \t")
	rest, ok = strings.CutPrefix(rest, "=")
	if !ok {
		return "", false
	}
	return decodeScalarString(strings.TrimSpace(rest))
}

// decodeScalarString reports the value a TOML string literal denotes, or
// false where text is not a single, self-contained string literal.
//
// It does NOT hand-roll an unescaper. It hands one candidate to the SAME
// decoder that filled `l.doc`, as the one-line document `id = <text>`, and
// reports what came back. `"reconcile-rewind"`, `"reconcile-rewind"` and
// `'reconcile-rewind'` all denote one string and must all confirm the rule
// they name; reimplementing TOML's escape table to learn that would be a
// second, divergable parser for a grammar the module already carries.
//
// Only a quoted literal is offered to the decoder. Bare text is rejected
// before the round trip, so a bare key or a non-string value (a number, a
// boolean, an inline table) reports false rather than decoding to something
// that never appeared in the document. Anything the decoder refuses — an
// unterminated quote, a trailing second value, an invalid escape — is also
// false, which leaves the caller on its existing degrade rather than on a
// guess.
func decodeScalarString(text string) (string, bool) {
	if len(text) < 2 {
		return "", false
	}
	if q := text[0]; q != '"' && q != '\'' {
		return "", false
	}
	var probe struct {
		ID string `toml:"id"`
	}
	if err := toml.Unmarshal([]byte("id = "+text), &probe); err != nil {
		return "", false
	}
	return probe.ID, true
}

// stripComment removes a trailing TOML comment from one line, leaving a `#`
// that sits INSIDE a quoted string alone.
//
// `headerLine`'s cruder strip — cut at the first `#` — is right for a
// bracketed header, whose text cannot carry a quoted `#` in the authorings
// that scan admits. It is wrong for a value assignment: a rule id may
// legally contain `#`, and cutting inside its quotes truncates the id so it
// can never match its own anchor.
//
// Inside a BASIC string a backslash escapes the next character, so `\"`
// does not close the string and must not flip the scan back out of it —
// otherwise `id = "a\"b" # c` is read as having closed at the escaped quote
// and its comment is left attached. A LITERAL string defines no escapes at
// all: a backslash in `'a\'` is an ordinary character and the first `'`
// closes it, so the escape rule is applied only to the basic form.
func stripComment(raw string) string {
	var quote rune
	escaped := false
	for i, r := range raw {
		switch {
		case escaped:
			escaped = false
		case quote != 0:
			switch {
			case quote == '"' && r == '\\':
				escaped = true
			case r == quote:
				quote = 0
			}
		case r == '"' || r == '\'':
			quote = r
		case r == '#':
			return raw[:i]
		}
	}
	return raw
}

// headerLine reports the 1-based line whose decoration-stripped text equals
// one of the accepted spellings, or ZERO when it matches zero times or more
// than once — a second match means the scan cannot tell which line the
// refusal belongs to.
func headerLine(src []byte, spellings ...string) int {
	line, matches := 0, 0
	for i, raw := range strings.Split(string(src), "\n") {
		// Trailing comments and surrounding space are the only decoration
		// such a line may carry; anything else is not this line.
		text := raw
		if hash := strings.IndexByte(text, '#'); hash >= 0 {
			text = text[:hash]
		}
		text = strings.TrimSpace(text)
		if !slices.Contains(spellings, text) {
			continue
		}
		matches++
		line = i + 1
	}
	if matches != 1 {
		return 0
	}
	return line
}

// checkClassAgreement holds the declared class and the owned set to the
// ONE-DIRECTIONAL agreement `0010:C1` fixes: a `decision-table` model MUST
// declare zero owned tags.
//
// The second direction is deliberately absent. A `state-machine` model
// declaring zero owned tags is rootless, not malformed, and `0006:C18`
// already reports it as `graph-dangling-edge` at lint — growing an arm here
// would make that model unauthorable and change behaviour this RDR leaves
// untouched.
//
// POSITION (`0010:C1`): the check needs the tag table, so it cannot live in
// `loadModelHeader`; and it must precede `checkAccessorBindings`, `run`'s
// last step, so a decision table that declares both an owned tag and
// `[initial]` refuses on the CLASS rather than on the writer-arity
// diagnostic. C1 draws a second consequence from the window's floor — "an
// undeclared-tag refusal precedes a class-disagreement refusal under
// `run`'s fail-fast order" — and that refusal is minted by
// `normalizeRules`. The two clauses are jointly satisfiable at exactly one
// place in the window: AFTER `normalizeRules` and before
// `checkAccessorBindings`, which is where this step runs.
func (l *loader) checkClassAgreement() error {
	if !IsDecisionTable(l.model) {
		return nil
	}
	var owned int
	for _, decl := range l.model.Tags {
		if decl.Provenance == ProvenanceOwned {
			owned++
		}
	}
	if owned == 0 {
		return nil
	}
	// The detail names the class and renders the count as the literal token
	// `owned=<n>`, which C1 fixes verbatim so a consumer can key on it.
	return fail(CatMalformedModelDeclaration,
		fmt.Sprintf("[model] class %s declares owned=%d; a %s model declares "+
			"zero tags of provenance %s",
			strconv.Quote(ClassDecisionTable), owned, ClassDecisionTable,
			ProvenanceOwned))
}

// kindFieldMapping names, for every kind that declares a type field, the
// field it takes. It is appended to all three kind/field mismatch refusals
// so the trio stays symmetric and an author who tripped one arm reads the
// mapping whole.
const kindFieldMapping = "only an enum declares domain, a set declares " +
	"elements = [...], an int declares min/max"

func tagDecl(key string, src sourceTagDecl) (TagDecl, error) {
	bad := func(detail string) (TagDecl, error) {
		return TagDecl{}, fail(CatMalformedTagDeclaration, "tag "+key+": "+detail)
	}

	switch Provenance(src.Provenance) {
	case ProvenanceOwned, ProvenanceObserved, ProvenanceRecognized:
	default:
		return bad("provenance " + strconv.Quote(src.Provenance) + " is not owned, observed, or recognized")
	}
	// RDR 0003 fixes exactly five kind tokens; `string` is not one of them
	// (deviations.md D1).
	if !IsDeclaredKind(src.Kind) {
		return bad("kind " + strconv.Quote(src.Kind) + " is outside RDR 0003's five tokens")
	}

	// The type model's fields are per kind: only an enum carries a domain,
	// only an int carries bounds, and only a set carries elements.
	//
	// Each arm carries the WHOLE mapping, not the one field its own kind
	// takes. The three arms are shared across kinds — `admits no domain`
	// fires for `set`, `int`, and `scalar` alike — so a hint answering for
	// a single kind would be wrong advice on the others. Stating "a set
	// declares elements" on the domain arm is also what turns a true but
	// useless refusal into an actionable one: read alone, "kind set admits
	// no domain" says a set cannot be finite, which is the opposite of the
	// truth.
	if len(src.Domain) > 0 && src.Kind != "enum" {
		return bad("kind " + src.Kind + " admits no domain; " + kindFieldMapping)
	}
	if (src.Min != nil || src.Max != nil) && src.Kind != "int" {
		return bad("kind " + src.Kind + " admits no min or max; " + kindFieldMapping)
	}
	if len(src.Elements) > 0 && src.Kind != "set" {
		return bad("kind " + src.Kind + " admits no elements; " + kindFieldMapping)
	}
	// The single-valued marker is meaningful only where the kind has both a
	// single-valued and a non-single-valued assignment count. A `set` holds
	// any subset, so its count is `2^|element universe|` and never
	// `|universe|`; a `scalar` has no finite declared domain to partition
	// (`0003:1273-1283`). `0003:1389` rejects both here, in the declaration
	// loader, under `0002:C22`'s malformed tag declaration.
	if src.SingleValued != nil && (src.Kind == "set" || src.Kind == "scalar") {
		return bad("kind " + src.Kind + " admits no single_valued marker")
	}
	if src.Min != nil && src.Max != nil && *src.Min > *src.Max {
		return bad("min exceeds max")
	}

	return TagDecl{
		Provenance:   Provenance(src.Provenance),
		Kind:         src.Kind,
		Domain:       src.Domain,
		Min:          src.Min,
		Max:          src.Max,
		Elements:     src.Elements,
		SingleValued: src.SingleValued != nil && *src.SingleValued,
		Required:     src.Required,
	}, nil
}

func (l *loader) loadAccessors() error {
	readers, err := l.accessorTable(l.doc.Read, "read", false)
	if err != nil {
		return err
	}
	writers, err := l.accessorTable(l.doc.Write, "write", true)
	if err != nil {
		return err
	}
	gates, err := l.accessorTable(l.doc.Gate, "gate", false)
	if err != nil {
		return err
	}

	// Every key in a `[write.<id>]` entry's keys list MUST be an OWNED tag
	// (`0002:C2`).
	for id, w := range writers {
		for _, key := range w.Keys {
			if l.model.Tags[key].Provenance != ProvenanceOwned {
				return fail(CatWriteToNonOwnedTag,
					"writer "+id+" names the non-owned tag "+key)
			}
		}
	}

	l.model.Readers = readers
	l.model.Writers = writers
	l.model.Gates = gates
	return nil
}

// accessorTable validates one capability table. Role, keys, and timeout
// are required, and any of them absent, empty, or ill-formed is a
// `malformed accessor declaration` — except a keys member naming an
// undeclared tag, which carries `unknown tag` (`0002:C2`).
//
// The CARRIER is exactly one of `path` or `command`, and its defects carry
// RDR 0025's own six categories rather than the inherited one (`0025:C1`,
// `0025:C5`). The other four rules are unchanged: C1 relaxes only the path
// rule.
func (l *loader) accessorTable(src map[string]sourceAcc, capability string, wantReadBack bool) (map[string]Accessor, error) {
	out := make(map[string]Accessor, len(src))

	for id, a := range src {
		bad := func(detail string) error {
			return fail(CatMalformedAccessorDeclaration, capability+" "+id+": "+detail)
		}
		switch {
		case a.Role == nil || *a.Role == "":
			return nil, bad("role is absent or empty")
		case a.Keys == nil || len(*a.Keys) == 0:
			return nil, bad("keys is absent or empty")
		case a.Timeout == nil:
			return nil, bad("timeout is absent")
		}
		d, err := time.ParseDuration(*a.Timeout)
		if err != nil {
			return nil, bad("timeout " + strconv.Quote(*a.Timeout) + " is not a Go duration")
		}
		if d <= 0 {
			return nil, bad("timeout " + *a.Timeout + " is not positive")
		}

		// The CARRIER, before the per-key checks: C5's arms judge the
		// entry's own declaration, and a raw-mode arity defect must not be
		// masked by a key that happens to name the reserved kernel key.
		if err := carrierDefect(a, capability, id, *a.Keys, l.model.Tags); err != nil {
			return nil, err
		}

		for _, key := range *a.Keys {
			if _, ok := l.model.Tags[key]; !ok {
				return nil, fail(CatUnknownTag,
					capability+" "+id+" names the undeclared tag "+key)
			}
			// The kernel supplies the recognized key; no accessor reads or
			// writes it (`0002:C2`).
			if key == RecognizedTagKey {
				return nil, fail(CatMalformedAccessorBinding,
					capability+" "+id+" names the reserved key "+RecognizedTagKey)
			}
		}

		if wantReadBack && (a.ReadBack == nil || !*a.ReadBack) {
			return nil, bad("a write entry requires read_back = true")
		}
		if !wantReadBack && a.ReadBack != nil {
			return nil, bad("read_back is a write-entry key")
		}

		acc := Accessor{
			Role:     *a.Role,
			Keys:     *a.Keys,
			Timeout:  *a.Timeout,
			ReadBack: a.ReadBack != nil && *a.ReadBack,
			Output:   a.Output,
		}
		if a.Path != nil {
			acc.Path = *a.Path
		}
		if a.Command != nil {
			acc.Command = slices.Clone(*a.Command)
		}
		if a.ExitAbsent != nil {
			acc.ExitAbsent = slices.Clone(*a.ExitAbsent)
		}
		if a.ExitVerdicts != nil {
			acc.ExitVerdicts = maps.Clone(*a.ExitVerdicts)
		}
		if a.Env != nil {
			acc.Env = maps.Clone(*a.Env)
		}
		if a.EnvPass != nil {
			acc.EnvPass = slices.Clone(*a.EnvPass)
		}
		if a.Edit != nil {
			// PRESENCE survives: a bare `[write.x.edit]` decodes to a
			// non-nil empty map, which is what the "both" arm keys on.
			acc.Edit = make(map[string]EditRule, len(*a.Edit))
			for key, r := range *a.Edit {
				rule := EditRule{}
				if r.Anchor != nil {
					rule.Anchor = *r.Anchor
				}
				if r.Replace != nil {
					rule.Replace = *r.Replace
				}
				if s, ok := r.Clear.(string); ok {
					rule.Clear = s
				}
				acc.Edit[key] = rule
			}
		}
		acc.Select = selectRules(a.Select)
		acc.Steps = stepsRules(a.Steps)
		out[id] = acc
	}
	return out, nil
}

// commandPlaceholders is the CLOSED v1 placeholder vocabulary (`0025:C2`).
// A placeholder is recognized only as a WHOLE argv element.
var commandPlaceholders = []string{"{artifact}"}

// shellInterpreters is the OPEN deny-list: interpreter words that take
// inline code on a flag, each holding its OWN inline-code flags. It is
// deliberately not closed — an unlisted spelling is admitted, so the check
// raises the cost of an inline-shell carrier without claiming to make one
// impossible. The list grows only by amending the clause (`0027:C1`,
// succeeding `0025:C5`'s interpreter-set line).
var shellInterpreters = map[string][]string{
	"sh":     {"-c"},
	"bash":   {"-c"},
	"dash":   {"-c"},
	"ksh":    {"-c"},
	"zsh":    {"-c"},
	"csh":    {"-c"},
	"tcsh":   {"-c"},
	"python": {"-c"},
	"ruby":   {"-e"},
	"node":   {"-e", "--eval"},
	"php":    {"-r"},
}

// envReservedPrefix is the literal, case-sensitive prefix the C4 overlay
// reserves. An `env` key or `env_pass` name matching it is a defect, so
// the overlay is never shadowed silently (`0025:C5`).
const envReservedPrefix = "INTRASTATE_"

// carrierDefect reports the FIRST C5 defect an entry carries, in the
// clause order C5 declares: conflict, empty, unknown placeholder, shell
// interpreter, output shape, env conflict. Load is fail-fast, so an entry
// carrying several reports the earliest (`0025:C5` precedence).
func carrierDefect(a sourceAcc, capability, id string, keys []string, tags map[string]TagDecl) error {
	where := capability + " " + id
	hasPath := a.Path != nil && *a.Path != ""
	hasCommand := a.Command != nil
	// PRESENCE, not emptiness: `edit_carrier_conflict` follows the "both"
	// arm's discipline, so a bare `[write.x.edit]` beside another carrier
	// is a conflict rather than a silently ignored second carrier
	// (`0028:C1.1`).
	hasEdit := a.Edit != nil
	// The same PRESENCE discipline for the fourth carrier (kata q14r).
	hasSteps := a.Steps != nil

	// 1 — command_and_path_conflict: both or neither carrier.
	//
	// "Both" is keyed on the KEYS being declared, so an empty `path`
	// beside a `command` is still a conflict rather than a silently
	// ignored second carrier; "neither" is keyed on neither being a
	// USABLE carrier, which is where the old "path is absent or empty"
	// arm lands now that a command entry is legal (REQ-7, REQ-14).
	//
	// The "neither" arm's PREDICATE widens to "none of the three"
	// (`0028:C1.1`) — an `edit`-only entry would otherwise be refused as
	// carrier-less — while its wire string stays 0025:C5's
	// `command_and_path_conflict`. Only the message text moves, to name
	// three carriers: the wire string is the contract and the message is
	// not. The `steps` carrier widens it to four the same way (kata q14r).
	switch {
	case a.Path != nil && a.Command != nil:
		return fail(CatCommandAndPathConflict,
			where+" declares both `path` and `command`; an entry carries "+
				"exactly one carrier")
	case !hasPath && !hasCommand && !hasEdit && !hasSteps:
		return fail(CatCommandAndPathConflict,
			where+" declares none of `path`, `command`, `edit` or `steps`; "+
				"an entry carries exactly one carrier")
	}

	// 1b — edit_carrier_conflict: `edit` beside another carrier, or on a
	// read or gate entry. It is evaluated with the carrier arms because
	// it IS a carrier arm; C1.4's five edit-table categories are decided
	// only once the carrier is established as one (`0028:C1.4`
	// precedence:).
	if hasEdit {
		switch {
		case a.Path != nil || a.Command != nil:
			return fail(CatEditCarrierConflict,
				where+" declares `edit` beside another carrier; an entry "+
					"carries exactly one of `path`, `command` or `edit`")
		case capability != "write":
			return fail(CatEditCarrierConflict,
				where+" declares `edit`, which is admissible on write entries only")
		}
	}

	// 1c — steps_carrier_conflict: `steps` beside another carrier, or on
	// a read or gate entry, on `edit_carrier_conflict`'s model (kata
	// q14r). Its table categories are decided once the carrier is
	// established as the one, exactly as C1.4's are.
	if hasSteps {
		switch {
		case a.Path != nil || a.Command != nil || hasEdit:
			return fail(CatStepsCarrierConflict,
				where+" declares `steps` beside another carrier; an entry "+
					"carries exactly one of `path`, `command`, `edit` or `steps`")
		case capability != "write":
			return fail(CatStepsCarrierConflict,
				where+" declares `steps`, which is admissible on write entries only")
		}
	}

	// 2–4 and the argv0 rule, factored into `argvDefect` so a `steps`
	// vector is judged by the same checks (kata q14r).
	if hasCommand {
		declared := func(key string) bool {
			_, ok := tags[key]
			return ok
		}
		if err := argvDefect(*a.Command, where, "`command`", declared); err != nil {
			return err
		}
	}

	// The command keys are properties of the invocation envelope
	// (`0025:C1`, C3/C4), so they mean something only on a carrier that
	// spawns a child: `command`, or `steps`, whose every step composes its
	// environment the same way (kata c6t5). These arms run after the
	// carrier arms above, so a carrier-less entry still reports
	// `command_and_path_conflict`.
	spawns := hasCommand || hasSteps

	// 5 — command_output_shape.
	if err := outputShapeDefect(a, capability, where, keys, spawns); err != nil {
		return err
	}

	// 6 — command_env_conflict. On an entry that never spawns, `env` and
	// `env_pass` would be silently ignored, so declaring either there is
	// this category before the reserved-prefix rule is consulted (kata
	// c6t5, composing `0025:C1` and `0025:C5`).
	if !spawns {
		for _, k := range []struct {
			name     string
			declared bool
		}{{"env", a.Env != nil}, {"env_pass", a.EnvPass != nil}} {
			if k.declared {
				return fail(CatCommandEnvConflict,
					where+" declares `"+k.name+"`, which requires a `command` "+
						"or `steps` carrier")
			}
		}
	}
	if a.Env != nil {
		for _, k := range slices.Sorted(maps.Keys(*a.Env)) {
			if strings.HasPrefix(k, envReservedPrefix) {
				return fail(CatCommandEnvConflict,
					where+" declares the `env` key "+k+", which matches the "+
						"reserved "+envReservedPrefix+" prefix the overlay owns")
			}
		}
	}
	if a.EnvPass != nil {
		for _, v := range *a.EnvPass {
			if strings.HasPrefix(v, envReservedPrefix) {
				return fail(CatCommandEnvConflict,
					where+" names the `env_pass` variable "+v+", which matches "+
						"the reserved "+envReservedPrefix+" prefix the overlay owns")
			}
		}
	}

	// The three `command_select_*` categories (kata wchf), after 0025:C5's
	// clauses 1–6 so no existing entry's reported defect moves. `select`
	// is read-only and `edit` write-only, so the two never both apply.
	if err := selectDefect(a, capability, where, keys, tags); err != nil {
		return err
	}

	// C1.4's remaining four edit-table categories, evaluated after
	// 0025:C5's clauses 1–6 and in C1.4's own registration order
	// (`0028:C1.4` precedence:).
	if hasEdit {
		return editTableDefect(*a.Edit, where, keys, tags)
	}
	if hasSteps {
		return stepsTableDefect(*a.Steps, where, keys, tags)
	}
	return nil
}

// argvDefect reports the FIRST defect one declared argv vector carries
// under `command`'s clauses 2–4 and the argv0 rule, in that order: empty,
// unknown placeholder, shell interpreter, `{tag.<key>}` at argv0
// (`0025:C5`, `0027:C1`, `0028:C1.4`).
//
// It is factored out of `carrierDefect` so a `steps` vector is judged by
// the SAME checks a `command` is and carries the same wire strings (kata
// q14r). `what` names the vector in the empty-vector detail, and
// `admitTag` is the `{tag.<key>}` admission: any DECLARED key for a
// `command` (`0028:C1.6`), an OBSERVED key for a step.
func argvDefect(argv []string, where, what string, admitTag func(string) bool) error {
	// 2 — command_empty: an empty vector or an empty element.
	if len(argv) == 0 {
		return fail(CatCommandEmpty, where+" declares an empty "+what+" vector")
	}
	for i, el := range argv {
		if el == "" {
			return fail(CatCommandEmpty,
				where+" declares an empty "+what+" element at index "+strconv.Itoa(i))
		}
	}

	interpEl, isInterp := interpreterForm(argv)

	// 3 — command_unknown_placeholder: an unknown or non-whole-element
	// `{…}` token. Never silently-literal text.
	for _, el := range argv {
		if slices.Contains(commandPlaceholders, el) {
			continue
		}
		// `{tag.<key>}` joins the vocabulary as a FAMILY, admitted
		// whole-element for a key the caller's `admitTag` admits
		// (`0028:C1.6`). An unadmitted key, and a `{…}` element that is
		// neither `{artifact}` nor an admitted `{tag.<key>}`, keep
		// 0025:C5's unchanged `command_unknown_placeholder` wire string.
		if key, ok := CommandTagKey(el); ok && admitTag(key) {
			continue
		}
		// A brace-bearing element carrying whitespace is a command
		// STRING — the `sh -c "cat {artifact}"` shape — which clause 4
		// owns; reading it as a malformed placeholder would report the
		// wrong defect and mask the interpreter form (deviations D6).
		// That exemption belongs to the INTERPRETER FORM, not to
		// whitespace as such: clause 4 fires only when some argv word
		// names a listed interpreter followed by one of its inline-code
		// flags (`0027:C1`), so under any other argv an exempted element
		// is owned by nothing and would reach the executor as literal,
		// unsubstituted argv — the outcome C2 forbids.
		if strings.ContainsAny(el, "{}") &&
			(!isInterp || strings.IndexFunc(el, unicode.IsSpace) < 0) {
			return fail(CatCommandUnknownPlaceholder,
				where+" declares the element "+el+
					", which carries a `{…}` token that is not exactly a known "+
					"placeholder; the v1 vocabulary is "+
					strings.Join(commandPlaceholders, ", "))
		}
	}

	// 4 — command_shell_interpreter: a listed interpreter word plus one
	// of its inline-code flags at any later argv position, under any
	// prefix (`0027:C1`, succeeding `0025:C5`'s argv0 line).
	if isInterp {
		return fail(CatCommandShellInterpreter,
			where+" declares the interpreter form "+interpEl+
				"; inline shell is not a declared command; put it in a script "+
				"and declare the script as argv0")
	}

	// edit_tag_argv0 (`0028:C1.4`, C1.6) — a `{tag.<key>}` element at
	// argv0 of a read, gate or WRITE entry's command. The executable
	// is the one word a reviewer must be able to read off the model,
	// and a caller-bound argv0 makes the interpreter deny-list
	// unenforceable against a name that does not exist until
	// invocation. It is statically decidable, so it is refused where
	// it is visible rather than at spawn.
	//
	// The rule is on THIS FAMILY only: `{artifact}` at argv0 stays
	// admitted, having no such rule and no such reviewer promise to
	// break.
	if _, ok := CommandTagKey(argv[0]); ok {
		return fail(CatEditTagArgv0,
			where+" declares the placeholder "+argv[0]+" at argv0; the "+
				"executable must be readable off the model, so a "+
				"`{tag.<key>}` element is admitted at any later position only")
	}
	return nil
}

// CommandTagKey reports the tag key a WHOLE argv element names as
// `{tag.<key>}`, and whether the element takes that form at all. It is
// whole-element by construction — `x{tag.nnnn}` is not a placeholder —
// which is 0025:C2's substitution rule that C1.6 joins rather than
// widens.
//
// It is exported because the load-time admission check and the
// apply-time substitution must recognize the SAME form: two spellings
// would let an element lint clean and then cross to a child unsubstituted,
// the outcome 0025:C2 forbids.
func CommandTagKey(el string) (string, bool) {
	if !strings.HasPrefix(el, editTagPrefix) || !strings.HasSuffix(el, "}") {
		return "", false
	}
	key := el[len(editTagPrefix) : len(el)-1]
	if key == "" || strings.ContainsAny(key, "{}") {
		return "", false
	}
	return key, true
}

// editTableDefect reports the FIRST C1.4 edit-table defect an entry
// carries, in the clause's registration order: key mismatch, anchor
// invalid, template invalid, clear invalid.
//
// Sibling `edit.<key>` tables are map-ranged, exactly as `accessorTable`'s
// entries are: which of two equally-defective tables is reported is
// unspecified and no test may assert it. The CATEGORY is not unspecified,
// which is why each step sweeps every rule before the next step runs
// (`0028:C1.4` precedence:).
func editTableDefect(
	rules map[string]sourceEditRule, where string, keys []string,
	tags map[string]TagDecl,
) error {
	// 2 — edit_key_mismatch: `keys` and the rule tables must be in
	// bijection, which is what makes "every planned key has a rule" a
	// lint-time guarantee rather than an apply-time surprise.
	for _, key := range keys {
		if _, ok := rules[key]; !ok {
			return fail(CatEditKeyMismatch,
				where+" declares the key "+key+" with no `edit."+key+"` table")
		}
	}
	for _, key := range slices.Sorted(maps.Keys(rules)) {
		if !slices.Contains(keys, key) {
			return fail(CatEditKeyMismatch,
				where+" declares an `edit."+key+"` table for a key not in `keys`")
		}
	}

	// A tag key is bindable at invocation only with `observed`
	// provenance: `flow_input.go::parseTags` refuses an owned key
	// (`flow-tag-owned`) and the recognized key (`flow-tag-reserved`), so
	// a declared key of either kind is structurally unbindable and naming
	// one is a LINT defect rather than a guaranteed runtime failure
	// (`0028:C1.2`).
	bindable := func(key string) bool {
		return tags[key].Provenance == ProvenanceObserved
	}

	// 3 — edit_anchor_invalid, over every rule. The anchor is COMPILED
	// here, with every `{tag.<key>}` replaced by a quoted probe, because
	// clause 4 needs its capture-group arity and `regexp.QuoteMeta` emits
	// no group syntax to perturb it.
	compiled := make(map[string]*regexp.Regexp, len(rules))
	for _, key := range slices.Sorted(maps.Keys(rules)) {
		r := rules[key]
		// `anchor` is REQUIRED. C1.1 marks only `clear` optional, so an
		// omitted anchor is a grammar defect and not an empty pattern:
		// "" compiles to the everything-matcher, which selects the sole
		// line of a one-line file and rewrites it. Refusing here is what
		// keeps `select:`'s exactly-one rule from being satisfied by
		// accident on an artifact no author aimed at.
		if r.Anchor == nil {
			return fail(CatEditAnchorInvalid,
				where+" `edit."+key+"` declares no `anchor`; only `clear` is "+
					"optional, and an absent anchor would match every line "+
					"rather than none")
		}
		anchor := *r.Anchor
		segs, err := ParseEditAnchor(anchor, bindable)
		if err != nil {
			return fail(CatEditAnchorInvalid,
				where+" `edit."+key+"` anchor: "+err.Error())
		}
		re, err := regexp.Compile(EditAnchorProbe(segs))
		if err != nil {
			return fail(CatEditAnchorInvalid,
				where+" `edit."+key+"` anchor does not compile as RE2: "+err.Error())
		}
		compiled[key] = re
	}

	// 4 — edit_template_invalid, over every rule.
	for _, key := range slices.Sorted(maps.Keys(rules)) {
		r := rules[key]
		// `replace` is REQUIRED for the same reason. An omitted template
		// is not "write nothing" — it is the empty LINE, so a rule whose
		// anchor still matches erases that line's content while
		// reporting success. Deleting a line is `clear`'s job (C1.5),
		// declared explicitly.
		if r.Replace == nil {
			return fail(CatEditTemplateInvalid,
				where+" `edit."+key+"` declares no `replace`; only `clear` is "+
					"optional, and an absent template would blank the selected "+
					"line rather than rewrite it (deletion is `clear`)")
		}
		replace := *r.Replace
		re := compiled[key]
		if _, err := ParseEditReplace(
			replace, key, re.NumSubexp(), re.SubexpNames(),
		); err != nil {
			return fail(CatEditTemplateInvalid,
				where+" `edit."+key+"` replace: "+err.Error())
		}
	}

	// 5 — edit_clear_invalid: `clear` outside the closed set {"line"}.
	// The deferred table form `clear = { replace = … }` and a boolean
	// land here rather than as malformed TOML, which is why the decoded
	// field is `any`.
	for _, key := range slices.Sorted(maps.Keys(rules)) {
		r := rules[key]
		if r.Clear == nil {
			continue
		}
		s, ok := r.Clear.(string)
		if !ok || s != EditClearLine {
			return fail(CatEditClearInvalid,
				where+" `edit."+key+"` declares a `clear` outside the closed "+
					"set {\""+EditClearLine+"\"}")
		}
	}
	return nil
}

// stepsTableDefect reports the FIRST `steps` table defect an entry
// carries (kata q14r), in check order: key mismatch, table invalid, then
// each step vector under `command`'s own argv clauses. Sibling tables,
// arms and members are swept in SORTED order, so the reported defect is
// deterministic as `editTableDefect`'s category is.
//
// The `clear` arms a clearing rule OWES are not judged here: which values
// it can hold is read off the normalized rows, so that arm of
// `steps_table_invalid` is `checkStepsClearArms` (kata t02k).
func stepsTableDefect(
	rules map[string]sourceStepsRule, where string, keys []string,
	tags map[string]TagDecl,
) error {
	// steps_key_mismatch — `keys` and the tables in bijection, mirroring
	// `edit.<key>`: every planned key has arms, so no key falls through
	// to a write that runs nothing.
	for _, key := range keys {
		if _, ok := rules[key]; !ok {
			return fail(CatStepsKeyMismatch,
				where+" declares the key "+key+" with no `steps."+key+"` table")
		}
	}
	for _, key := range slices.Sorted(maps.Keys(rules)) {
		if !slices.Contains(keys, key) {
			return fail(CatStepsKeyMismatch,
				where+" declares a `steps."+key+"` table for a key not in `keys`")
		}
	}

	// steps_table_invalid — the arms against the tag's domain. A step
	// table maps each VALUE to literal argv, so the tag must have finitely
	// many values and every one it can be planned must have a set arm;
	// a value the tag cannot take is a typo lint can see.
	for _, key := range slices.Sorted(maps.Keys(rules)) {
		r := rules[key]
		at := where + " `steps." + key + "`"
		bad := func(detail string) error {
			return fail(CatStepsTableInvalid, at+" "+detail)
		}
		decl := tags[key]
		domain, finite := stepsDomain(decl)
		if !finite {
			return bad("is declared for the " + decl.Kind + "-kind tag " + key +
				"; steps name argv per value, so the tag needs a finite domain " +
				"(an enum with a domain, or a bool) — write it through " +
				"`command`, which carries the value on stdin, or `edit`")
		}
		if r.Set == nil {
			return bad("declares no `set` arms")
		}
		for _, arm := range []struct {
			name  string
			table *map[string][][]string
		}{{"set", r.Set}, {"clear", r.Clear}} {
			if arm.table == nil {
				continue
			}
			for _, member := range slices.Sorted(maps.Keys(*arm.table)) {
				if !slices.Contains(domain, member) {
					return bad("declares the `" + arm.name + "` arm " +
						strconv.Quote(member) + ", which is not a value of " + key)
				}
				if len((*arm.table)[member]) == 0 {
					return bad("declares an empty `" + arm.name + "." + member +
						"` arm; an arm runs at least one step")
				}
			}
		}
		for _, member := range domain {
			if _, ok := (*r.Set)[member]; !ok {
				return bad("declares no `set." + member + "` arm; every value " +
					key + " can be planned needs one")
			}
		}
	}

	// The step vectors themselves, under `command`'s clauses 2–4 and the
	// argv0 rule. A step names only an OBSERVED `{tag.<key>}`: that is the
	// one provenance a caller can bind (`flow_input.go::parseTags` refuses
	// owned and recognized keys), so an owned key here would lint clean
	// and then refuse every invocation as unbound — and the planned value
	// never crosses argv at all.
	observed := func(key string) bool {
		return tags[key].Provenance == ProvenanceObserved
	}
	for _, key := range slices.Sorted(maps.Keys(rules)) {
		r := rules[key]
		for _, arm := range []struct {
			name  string
			table *map[string][][]string
		}{{"set", r.Set}, {"clear", r.Clear}} {
			if arm.table == nil {
				continue
			}
			for _, member := range slices.Sorted(maps.Keys(*arm.table)) {
				for i, argv := range (*arm.table)[member] {
					at := where + " `steps." + key + "." + arm.name + "." +
						member + "[" + strconv.Itoa(i) + "]`"
					if err := argvDefect(argv, at, "step", observed); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// stepsDomain is the finite value set a `steps` table enumerates: an
// enum's declared domain, or a bool's two literals (kata q14r).
func stepsDomain(decl TagDecl) ([]string, bool) {
	switch {
	case decl.Kind == "enum" && len(decl.Domain) > 0:
		return decl.Domain, true
	case decl.Kind == "bool":
		return []string{"false", "true"}, true
	default:
		return nil, false
	}
}

// stepsRules converts the decoded tables into the model's typed form.
// PRESENCE survives: a bare `[write.x.steps]` decodes to a non-nil empty
// map, as `edit` does.
func stepsRules(src *map[string]sourceStepsRule) map[string]StepsRule {
	if src == nil {
		return nil
	}
	clone := func(arms *map[string][][]string) map[string][][]string {
		if arms == nil {
			return nil
		}
		out := make(map[string][][]string, len(*arms))
		for member, steps := range *arms {
			cloned := make([][]string, 0, len(steps))
			for _, argv := range steps {
				cloned = append(cloned, slices.Clone(argv))
			}
			out[member] = cloned
		}
		return out
	}
	out := make(map[string]StepsRule, len(*src))
	for key, r := range *src {
		out[key] = StepsRule{Set: clone(r.Set), Clear: clone(r.Clear)}
	}
	return out
}

// interpreterForm reports the offending pair when argv carries a listed
// interpreter word followed, at ANY later argv position, by one of THAT
// interpreter's inline-code flags — under any prefix (`env` and its
// options, `nice`, `timeout`, `xargs`, `doas`, …). It is a two-index scan:
// refuse iff there exist i < j with `filepathBase(argv[i])` a listed
// interpreter and `argv[j]` one of its listed flags, reporting
// `argv[i] + " " + argv[j]` for the lowest i, then the lowest j.
//
// Nothing before argv[i] is read: no wrapper table exists and none may be
// added, so an unenumerated prefix is caught by the same rule the named
// ones are. A listed word with no qualifying flag after it does not stop
// the scan — the tie-break is on the lowest PAIR, not the lowest listed
// basename. The scan iterates argv positions, never `shellInterpreters`,
// so the reported form cannot depend on map iteration order.
//
// The check reads argv WORDS only: it never splits a word on whitespace,
// and never reads stdin, files, PATH, or the resolved binary. A shell
// string carried in ONE word (`env -S "sh -c …"`) and an interpreter
// reading its script from stdin (`sh -s`, bare `sh`, `python -`) are out
// of scope BY NAME and stay green (`0027:C1`).
func interpreterForm(argv []string) (string, bool) {
	for i, word := range argv {
		flags, known := shellInterpreters[filepathBase(word)]
		if !known {
			continue
		}
		for _, arg := range argv[i+1:] {
			if slices.Contains(flags, arg) {
				return word + " " + arg, true
			}
		}
	}
	return "", false
}

// filepathBase is `path.Base` over a declared argv word, so `/bin/sh -c`
// is recognized as the same form as `sh -c`. The post-slash segment is
// matched byte-exactly — no suffix, alias, or case folding (`0027:C1`). It
// stays a string operation: this package performs no path RESOLUTION
// (`0002:EIA`).
func filepathBase(s string) string {
	if i := strings.LastIndexByte(s, '/'); i >= 0 {
		return s[i+1:]
	}
	return s
}

// outputShapeDefect covers C5's three output-shape arms: `output = "raw"`
// with a declared key set whose ARITY is not one, `exit_absent` on a
// non-read entry, and `exit_verdicts` on a non-gate entry or naming a
// value outside the closed verdict set. `output` itself is a read-entry
// key admitting exactly `"json"` or `"raw"` (`0025:C1`, `0025:C3`).
//
// All three shape a spawned command's stdout or exit status, so on an
// entry whose carrier never spawns (`spawns` false) declaring any of them
// is this category first, ahead of the capability arms (kata c6t5).
func outputShapeDefect(a sourceAcc, capability, where string, keys []string, spawns bool) error {
	if !spawns {
		for _, k := range []struct {
			name     string
			declared bool
		}{
			{"output", a.Output != nil},
			{"exit_absent", a.ExitAbsent != nil},
			{"exit_verdicts", a.ExitVerdicts != nil},
		} {
			if k.declared {
				return fail(CatCommandOutputShape,
					where+" declares `"+k.name+"`, which requires a `command` carrier")
			}
		}
	}
	if a.Output != nil {
		if capability != "read" {
			return fail(CatCommandOutputShape,
				where+" declares `output`, which is a read-entry key")
		}
		switch *a.Output {
		case "json":
		case "raw":
			// Valid only when the declared key list has arity one. The
			// ARITY is the declared list's length, not its distinct set: a
			// repeated key still declares two, which raw mode cannot carry.
			if len(keys) != 1 {
				return fail(CatCommandOutputShape,
					where+` declares output = "raw" with `+
						strconv.Itoa(len(keys))+" declared keys; raw mode carries "+
						"the single declared key's value and is valid only at arity one")
			}
		default:
			return fail(CatCommandOutputShape,
				where+" declares output "+strconv.Quote(*a.Output)+
					`; the admitted literals are "json" and "raw"`)
		}
	}
	if a.ExitAbsent != nil && capability != "read" {
		return fail(CatCommandOutputShape,
			where+" declares `exit_absent`, which is a read-entry key")
	}
	if a.ExitVerdicts != nil {
		if capability != "gate" {
			return fail(CatCommandOutputShape,
				where+" declares `exit_verdicts`, which is a gate-entry key")
		}
		for _, code := range slices.Sorted(maps.Keys(*a.ExitVerdicts)) {
			v := (*a.ExitVerdicts)[code]
			if !slices.Contains(commandVerdicts, v) {
				return fail(CatCommandOutputShape,
					where+" maps exit "+code+" to "+strconv.Quote(v)+
						", which is not a gate verdict")
			}
		}
	}
	return nil
}

// commandVerdicts mirrors `accessor.Verdicts()`, which this package cannot
// import: `internal/accessor` imports `internal/table`, so naming the
// three here is what keeps the dependency one-way. RDR 0004 owns the
// vocabulary; C3 fixes that the exit map's values ARE those strings.
var commandVerdicts = []string{"allow", "deny", "indeterminate"}

// checkAccessorBindings enforces the provenance-scoped arity rules, which
// need the normalized rows: the written-key set is what the rules' write
// blocks and clear lists name (`0002:C2`).
//
// Every OWNED tag MUST be served by exactly one reader; an OBSERVED tag MAY
// be served by at most one, and zero is legal because an observed key may
// arrive from the caller (JDR 0001 §JD-9).
//
// The writer arity is scoped by USE rather than by provenance: a key some
// rule writes or clears, or `[initial]` assigns, owes exactly one writer,
// while an owned key nothing writes owes at most one. Both walks cover every
// declared tag — a two-writer key outside `written` is malformed even though
// nothing in the model writes it.
func (l *loader) checkAccessorBindings() error {
	readerCount := map[string]int{}
	for _, r := range l.model.Readers {
		for _, key := range r.Keys {
			readerCount[key]++
		}
	}
	for _, key := range slices.Sorted(maps.Keys(l.model.Tags)) {
		decl := l.model.Tags[key]
		switch decl.Provenance {
		case ProvenanceOwned:
			if readerCount[key] != 1 {
				return fail(CatMalformedAccessorBinding,
					fmt.Sprintf("owned tag %s is served by %d readers; want exactly one",
						key, readerCount[key]))
			}
		case ProvenanceObserved:
			if readerCount[key] > 1 {
				return fail(CatMalformedAccessorBinding,
					fmt.Sprintf("observed tag %s is served by %d readers; want at most one",
						key, readerCount[key]))
			}
		}
	}

	writerCount := map[string]int{}
	for _, w := range l.model.Writers {
		for _, key := range w.Keys {
			writerCount[key]++
		}
	}
	// Every key any rule's write block or clear list names, and every key
	// `[initial]` assigns, MUST be served by exactly one writer.
	written := map[string]bool{}
	for _, row := range l.model.Rows {
		for _, key := range row.RequiresOwned {
			written[key] = true
		}
	}
	for _, t := range l.model.Initial {
		written[t.Key] = true
	}
	// The walk is over every DECLARED tag, the way the reader guard above
	// walks them, rather than over `written` alone. A key can be owned,
	// served by two writers, and named by no rule write, no clear list, and
	// no `[initial]` — outside `written` entirely — and such a model used to
	// load clean, leaving the CLI to route the mutation to whichever writer
	// sorted first (kata 8dg3).
	//
	// The arity itself stays provenance- and use-scoped:
	//
	//   - a key in `written` owes EXACTLY one writer, because something in
	//     the model actually writes it;
	//   - an owned key nothing writes owes AT MOST one, because zero is
	//     legal — the RDR 0005 MVV fixture's `note` is owned, read by one
	//     reader, and served by no writer, and `internal/cli::writerFor`
	//     names that case as the reason it checks writer keys rather than
	//     ownership. Tightening zero into a refusal would make such a key
	//     unauthorable.
	keys := map[string]bool{}
	for key := range l.model.Tags {
		keys[key] = true
	}
	for key := range written {
		keys[key] = true
	}
	for _, key := range slices.Sorted(maps.Keys(keys)) {
		switch {
		case written[key]:
			if writerCount[key] != 1 {
				return fail(CatMalformedAccessorBinding,
					fmt.Sprintf("written tag %s is served by %d writers; want exactly one",
						key, writerCount[key]))
			}
		case l.model.Tags[key].Provenance == ProvenanceOwned:
			if writerCount[key] > 1 {
				return fail(CatMalformedAccessorBinding,
					fmt.Sprintf("owned tag %s is served by %d writers; want at most one",
						key, writerCount[key]))
			}
		}
	}
	// The one writer-against-rows check beyond arity. It rides this step
	// rather than adding one to `run`: it needs the same normalized rows,
	// and `0010:C1` fixes this step as `run`'s last.
	return l.checkStepsClearArms()
}

// checkStepsClearArms is the `steps_table_invalid` arm that needs the
// normalized rows (kata t02k). A planned `<clear>` runs `clear.<held>`
// alone, so every value a clearing rule can HOLD needs a clear arm —
// otherwise the removal would run nothing and read-back would be the first
// to notice (`0004:C11`). A value no clearing rule can hold owes none, and
// leaving it arm-less is what keeps a compare-and-set the FIRST step of
// the `set` arm entered from it: a replace runs `clear.<held>` only when
// that arm exists.
//
// What a row can hold is read off its normalized atoms on the key — match
// `eq` (a match `in` has already expanded into one row per member) and
// `guard.all` `eq`/`in`, intersected, inherited atoms included. A row with
// none can hold any domain value. Every other atom is ignored, which only
// over-approximates: an `unless` may leave an arm owed that no row can
// select, never the reverse.
func (l *loader) checkStepsClearArms() error {
	for _, id := range slices.Sorted(maps.Keys(l.model.Writers)) {
		w := l.model.Writers[id]
		for _, key := range slices.Sorted(maps.Keys(w.Steps)) {
			domain, finite := stepsDomain(l.model.Tags[key])
			if !finite {
				continue
			}
			arms := w.Steps[key].Clear
			for _, row := range l.model.Rows {
				if !rowClears(row, key) {
					continue
				}
				for _, held := range rowHolds(row, key, domain) {
					if _, ok := arms[held]; ok {
						continue
					}
					return fail(CatStepsTableInvalid,
						"write "+id+" `steps."+key+"` declares no `clear."+held+
							"` arm, and rule "+row.RuleID+" clears "+key+
							" while it holds "+strconv.Quote(held)+
							"; declare `clear."+held+"`")
				}
			}
		}
	}
	return nil
}

// rowClears reports whether the row's write set renders `<clear>` for key.
func rowClears(row Row, key string) bool {
	return slices.ContainsFunc(row.Writes, func(t TagValue) bool {
		return t.Key == key && isClear(t.Value)
	})
}

// rowHolds is the subset of domain the row's own atoms on key admit, in
// domain order: the intersection of its match and `guard.all` `eq`/`in`
// literals, or the whole domain when it carries none.
func rowHolds(row Row, key string, domain []string) []string {
	out := slices.Clone(domain)
	for _, a := range row.Atoms {
		if a.Key != key || (a.Operator != "eq" && a.Operator != "in") {
			continue
		}
		if a.Block != BlockMatch && a.Block != BlockAll {
			continue
		}
		out = slices.DeleteFunc(out, func(v string) bool {
			return !slices.Contains(a.Literal, v)
		})
	}
	return out
}

// ------------------------------------------------------------------ dump

func (l *loader) loadDump() error {
	if l.doc.Dump == nil {
		l.model.DumpOrder = DumpColumns()
		return nil
	}
	order := l.doc.Dump.Order
	vocabulary := DumpColumns()

	seen := map[string]bool{}
	for _, col := range order {
		if !slices.Contains(vocabulary, col) {
			return fail(CatMalformedDumpDeclaration, "unknown column identifier "+strconv.Quote(col))
		}
		if seen[col] {
			return fail(CatMalformedDumpDeclaration, "column "+col+" is repeated")
		}
		seen[col] = true
	}
	// Settings MAY reorder the rendered columns; they MUST NOT omit a
	// field — a truncated dump is a refusal, never a silent narrowing
	// (`0002:C19`).
	for _, col := range vocabulary {
		if !seen[col] {
			return fail(CatMalformedDumpDeclaration, "column "+col+" is omitted")
		}
	}
	l.model.DumpOrder = order
	return nil
}

// -------------------------------------------------------------- contexts

func (l *loader) loadContexts() error {
	l.contexts = make(map[string][]Atom, len(l.doc.Context))
	l.resolved = make(map[string][]Atom, len(l.doc.Context))

	for id, ctx := range l.doc.Context {
		atoms, err := l.atomsFromBlock(ctx.Match, BlockMatch, "context "+id)
		if err != nil {
			return err
		}
		l.contexts[id] = atoms
	}
	// Inheritance targets must resolve before the transitive closure runs.
	for id, ctx := range l.doc.Context {
		if ctx.Inherits == "" {
			continue
		}
		if _, ok := l.doc.Context[ctx.Inherits]; !ok {
			return fail(CatUnknownContext,
				"context "+id+" inherits the unknown context "+ctx.Inherits)
		}
	}
	for _, id := range slices.Sorted(maps.Keys(l.doc.Context)) {
		if _, err := l.contextAtoms(id, nil); err != nil {
			return err
		}
	}
	return nil
}

// contextAtoms returns a context's transitive atom set, merging its own
// atoms with those it inherits. Inheritance MUST normalize to an explicit
// predicate set before lint or resolution, and it never overrides — it only
// accumulates (`0002:C6`).
func (l *loader) contextAtoms(id string, stack []string) ([]Atom, error) {
	if atoms, ok := l.resolved[id]; ok {
		return atoms, nil
	}
	if slices.Contains(stack, id) {
		return nil, fail(CatCyclicContextInheritance,
			"context inheritance cycle through "+strings.Join(append(stack, id), " -> "))
	}

	own := l.contexts[id]
	parent := l.doc.Context[id].Inherits
	if parent == "" {
		merged := mergeAtoms(nil, own)
		l.resolved[id] = merged
		return merged, nil
	}

	inherited, err := l.contextAtoms(parent, append(stack, id))
	if err != nil {
		return nil, err
	}
	merged := mergeAtoms(inherited, own)
	l.resolved[id] = merged
	return merged, nil
}

// --------------------------------------------------- root and stop set

func (l *loader) loadInitial() error {
	keys := slices.Sorted(maps.Keys(l.doc.Initial))
	assignments := make([]TagValue, 0, len(keys))
	for _, key := range keys {
		decl, ok := l.model.Tags[key]
		if !ok {
			return fail(CatUnknownTag, "[initial] assigns the undeclared tag "+key)
		}
		members, err := valueMembers(l.doc.Initial[key])
		if err != nil {
			return fail(CatMalformedInitialDeclaration, "[initial] "+key+": "+err.Error())
		}
		// The reserved-value rule takes precedence over the value arm: the
		// sentinel is refused wherever a tag value is authored (`0002:C11`).
		if slices.Contains(members, ClearSentinel) {
			return fail(CatReservedTagValue, "[initial] "+key+" authors the reserved value "+ClearSentinel)
		}
		// An [initial] value ill-formed for its declared kind or outside
		// its declared domain is scoped by SITE to this category, while the
		// same defect on a predicate literal is `malformed predicate atom`
		// (req-list Q2).
		if err := conform(decl, "eq", members); err != nil {
			return fail(CatMalformedInitialDeclaration, "[initial] "+key+": "+err.Error())
		}
		// Well-formedness for the declared kind includes ARITY: only a
		// `set` kind holds a member sequence, so a multi-member value on
		// any other kind is ill-formed for that kind. Refusing it here
		// keeps the seam from truncating it — `resolve.Tag.Value` is one
		// string, so the surplus members would vanish unreported.
		if decl.Kind != "set" && len(members) != 1 {
			return fail(CatMalformedInitialDeclaration,
				"[initial] "+key+": kind "+decl.Kind+
					" holds one value, not a member sequence")
		}
		// Every [initial] key MUST be a declared OWNED tag. A non-owned key
		// is caught by the writer binding, which no observed tag can
		// satisfy without tripping `write to non-owned tag` first.
		assignments = append(assignments, TagValue{Key: key, Value: members})
	}
	l.model.Initial = assignments
	return nil
}

func (l *loader) loadTerminal() error {
	sets := make([][]Atom, 0, len(l.doc.Terminal))
	for _, id := range l.doc.Terminal {
		if _, ok := l.doc.Context[id]; !ok {
			return fail(CatUnknownContext, "terminal names the unknown context "+id)
		}
		// A terminal context id is a REFERENCE, not a state name:
		// normalization dereferences each to the context's explicit
		// predicate set, and the normalized value carries the sets — never
		// the bare ids (`0002:C2`).
		atoms, err := l.contextAtoms(id, nil)
		if err != nil {
			return err
		}
		sets = append(sets, slices.Clone(atoms))
	}
	l.model.Terminal = sets
	return nil
}

// -------------------------------------------------------------- helpers

// valueMembers renders one authored TOML value as a member sequence. A
// scalar is a one-member sequence; an array is its members in authored
// order (sorting is the caller's, since only set-valued literals sort).
func valueMembers(v any) ([]string, error) {
	switch t := v.(type) {
	case string:
		return []string{t}, nil
	case bool:
		return []string{strconv.FormatBool(t)}, nil
	case int64:
		return []string{strconv.FormatInt(t, 10)}, nil
	case float64:
		return []string{strconv.FormatFloat(t, 'g', -1, 64)}, nil
	case []any:
		out := make([]string, 0, len(t))
		for _, item := range t {
			// A nested array is not a member sequence and the format
			// admits no such spelling, so it is refused on its SHAPE, not
			// on its arity. Testing arity instead unwrapped a one-element
			// nested array and silently reinterpreted `[["a"], "b"]` as
			// `["a", "b"]`, which defeats `0002:C3`'s governing principle
			// — a malformed authoring is a stable refusal, never a silent
			// no-op — and made the same authoring error refuse or succeed
			// depending on how many elements the inner array happened to
			// hold.
			if isArray(item) {
				return nil, fmt.Errorf("nested array literals are not a member sequence")
			}
			members, err := valueMembers(item)
			if err != nil {
				return nil, err
			}
			out = append(out, members...)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("value %v is not a tag value", v)
	}
}

// isArray reports whether an authored value was spelled as a TOML array.
func isArray(v any) bool {
	_, ok := v.([]any)
	return ok
}

// isBool reports whether an authored value was spelled as a TOML boolean.
func isBool(v any) bool {
	_, ok := v.(bool)
	return ok
}

// ConformValue checks one CALLER-supplied member against a tag's declared
// type model — its kind first, then its domain: an enum's `domain`, an
// int's `min`/`max`, or a set's `elements`.
//
// It is the runtime counterpart of the load-time `conform`, exported so the
// CLI can hold a `--write` or `--tag` value to the same declaration a rule's
// authored literal is already held to. Letting a caller write what a rule
// may not author would make a declaration advisory.
//
// There is no operator here: a caller supplies a VALUE, not a predicate, so
// the `exists` and comparison-bound arms `conform` carries have no analogue.
// A member is one scalar; a set's members are conformed one at a time by the
// caller, so the refusal can name the offending member as it was authored.
//
// A zero TagDecl conforms everything: an undeclared kind has no domain to
// violate, which is what keeps an undeclared `--tag` key shape-only.
func ConformValue(decl TagDecl, member string) error {
	if err := conformKind(decl, member); err != nil {
		return err
	}
	return conformDomain(decl, member)
}

// conform checks one literal against a tag's declared type model.
//
// Conformance is PER OPERATOR, not per literal: `eq`, `in`, and `contains`
// take members and are domain-checked; `exists` takes a bool literal, never
// a member; `lt`/`lte`/`gt`/`gte` take an ordered BOUND, which is
// kind-checked but not domain-checked (`0002:C17`).
func conform(decl TagDecl, operator string, members []string) error {
	switch operator {
	case resolve.OpExists:
		if len(members) != 1 || (members[0] != resolve.LiteralTrue && members[0] != resolve.LiteralFalse) {
			return fmt.Errorf("an existence literal is %s or %s",
				resolve.LiteralTrue, resolve.LiteralFalse)
		}
		return nil
	case "lt", "lte", "gt", "gte":
		if len(members) != 1 {
			return fmt.Errorf("a comparison bound is one value")
		}
		return conformKind(decl, members[0])
	}

	for _, m := range members {
		if err := conformKind(decl, m); err != nil {
			return err
		}
		if err := conformDomain(decl, m); err != nil {
			return err
		}
	}
	return nil
}

// conformKind checks a member against the tag's declared kind alone.
func conformKind(decl TagDecl, member string) error {
	switch decl.Kind {
	case "int":
		if _, err := strconv.Atoi(member); err != nil {
			return fmt.Errorf("%s is not an int", strconv.Quote(member))
		}
	case "bool":
		if member != "true" && member != "false" {
			return fmt.Errorf("%s is not a bool", strconv.Quote(member))
		}
	}
	return nil
}

// conformDomain checks a member against the tag's declared domain: an
// enum's `domain`, an int's `min`/`max`, or a set's `elements`.
func conformDomain(decl TagDecl, member string) error {
	switch decl.Kind {
	case "enum":
		if len(decl.Domain) > 0 && !slices.Contains(decl.Domain, member) {
			return fmt.Errorf("%s is outside the declared domain", strconv.Quote(member))
		}
	case "set":
		if len(decl.Elements) > 0 && !slices.Contains(decl.Elements, member) {
			return fmt.Errorf("%s is outside the declared elements", strconv.Quote(member))
		}
	case "int":
		n, err := strconv.Atoi(member)
		if err != nil {
			return fmt.Errorf("%s is not an int", strconv.Quote(member))
		}
		if decl.Min != nil && n < *decl.Min {
			return fmt.Errorf("%d is below min %d", n, *decl.Min)
		}
		if decl.Max != nil && n > *decl.Max {
			return fmt.Errorf("%d is above max %d", n, *decl.Max)
		}
	}
	return nil
}
