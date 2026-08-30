package table

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"

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
func tableHeader(raw string) string {
	text := strings.TrimSpace(stripComment(raw))
	if !strings.HasPrefix(text, "[") || !strings.HasSuffix(text, "]") {
		return ""
	}
	if strings.Contains(strings.Trim(text, "[]"), ".") {
		return ""
	}
	return text
}

// scalarAssignment reads a bare `key = <string>` assignment off one line and
// reports the value with its quoting removed.
//
// It tolerates every spelling TOML admits for the same document: arbitrary
// horizontal whitespace around the key and the `=`, and all three string
// syntaxes — basic (`"v"`), literal ('v'), and their multi-line forms'
// single-line use. It reports false for a dotted key, an inline table, or
// any non-string value, none of which name a rule the way an `id` does.
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
	rest = strings.TrimSpace(rest)
	for _, q := range []string{`"""`, "'''", `"`, "'"} {
		body, ok := strings.CutPrefix(rest, q)
		if !ok {
			continue
		}
		body, ok = strings.CutSuffix(body, q)
		if !ok {
			return "", false
		}
		return body, true
	}
	return "", false
}

// stripComment removes a trailing TOML comment from one line, leaving a `#`
// that sits INSIDE a quoted string alone.
//
// `headerLine`'s cruder strip — cut at the first `#` — is right for a
// bracketed header, whose text cannot carry a quoted `#` in the authorings
// that scan admits. It is wrong for a value assignment: a rule id may
// legally contain `#`, and cutting inside its quotes truncates the id so it
// can never match its own anchor.
func stripComment(raw string) string {
	var quote rune
	for i, r := range raw {
		switch {
		case quote != 0:
			if r == quote {
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

// accessorTable validates one capability table. Each of role, path, keys,
// and timeout is required, and any of them absent, empty, or ill-formed is
// a `malformed accessor declaration` — except a keys member naming an
// undeclared tag, which carries `unknown tag` (`0002:C2`).
func (l *loader) accessorTable(src map[string]sourceAcc, capability string, wantReadBack bool) (map[string]Accessor, error) {
	out := make(map[string]Accessor, len(src))
	for id, a := range src {
		bad := func(detail string) error {
			return fail(CatMalformedAccessorDeclaration, capability+" "+id+": "+detail)
		}
		switch {
		case a.Role == nil || *a.Role == "":
			return nil, bad("role is absent or empty")
		case a.Path == nil || *a.Path == "":
			return nil, bad("path is absent or empty")
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

		out[id] = Accessor{
			Role:     *a.Role,
			Path:     *a.Path,
			Keys:     *a.Keys,
			Timeout:  *a.Timeout,
			ReadBack: a.ReadBack != nil && *a.ReadBack,
		}
	}
	return out, nil
}

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
	return nil
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
