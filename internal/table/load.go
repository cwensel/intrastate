package table

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"

	toml "github.com/pelletier/go-toml/v2"

	"github.com/newcoinc/intrastate/internal/resolve"
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

	ld := &loader{doc: &doc, sourceID: sourceID}
	return ld.run()
}

// loader carries one document through validation and normalization.
type loader struct {
	doc      *sourceDoc
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
		l.loadAccessors,
		l.loadDump,
		l.loadContexts,
		l.loadInitial,
		l.loadTerminal,
		l.normalizeRules,
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

func (l *loader) loadTags() error {
	decls := make(map[string]TagDecl, len(l.doc.Tags))
	for key, src := range l.doc.Tags {
		decl, err := tagDecl(key, src)
		if err != nil {
			return err
		}
		decls[key] = decl
	}

	// A `recognized` declaration MUST be named `recognized`, and no owned
	// or observed declaration may take that name (`0002:C12`).
	for key, decl := range decls {
		if decl.Provenance == ProvenanceRecognized && key != RecognizedTagKey {
			return fail(CatReservedTagKey,
				"tag "+key+" declares provenance recognized under another name")
		}
		if key == RecognizedTagKey && decl.Provenance != ProvenanceRecognized {
			return fail(CatReservedTagKey,
				"tag "+RecognizedTagKey+" is reserved for the recognized declaration")
		}
	}
	if _, ok := decls[RecognizedTagKey]; !ok {
		return fail(CatMalformedModelDeclaration, "no [tags.recognized] declaration")
	}

	l.model.Tags = decls
	return nil
}

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
	if len(src.Domain) > 0 && src.Kind != "enum" {
		return bad("kind " + src.Kind + " admits no domain")
	}
	if (src.Min != nil || src.Max != nil) && src.Kind != "int" {
		return bad("kind " + src.Kind + " admits no min or max")
	}
	if len(src.Elements) > 0 && src.Kind != "set" {
		return bad("kind " + src.Kind + " admits no elements")
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
	for _, key := range slices.Sorted(maps.Keys(written)) {
		if writerCount[key] != 1 {
			return fail(CatMalformedAccessorBinding,
				fmt.Sprintf("written tag %s is served by %d writers; want exactly one",
					key, writerCount[key]))
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
