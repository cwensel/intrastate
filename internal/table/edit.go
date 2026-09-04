package table

// RDR 0028 `0028:C1.2` — the two template vocabularies of the `edit`
// carrier, and the ONE parser that decides both.
//
// The vocabularies are DIFFERENT because one of the templates is a
// regexp, and the asymmetry is the clause's central point:
//
//   - `replace` has a CLOSED vocabulary — `{<key>}`, `${N}`/`${name}`,
//     `$$`, `{{`/`}}` — and any other `{…}` or `$…` form is a defect,
//     including a bare `$1` that `regexp.Expand` would accept.
//   - `anchor` opens a placeholder on the fixed prefix `{tag.` ONLY.
//     Every other `{`, `}` and `$` is passed to RE2 untouched, so
//     `\d{4}`, `a{2,3}` and `^…$` are ordinary pattern text and
//     `{{`/`}}` is NOT an escape there.
//
// The parse is this clause's own and never a call to `regexp.Expand`,
// which re-scans its template at expansion time: captured text and
// substituted values must emit as literal bytes, never as grammar
// (`0028:C1.2` parse once:).
//
// It lives in this package rather than beside the write binding because
// BOTH callers need it and only one of them is in `flowbind`: the loader
// validates the templates at load, and the binding expands them at apply.
// A second parser would be a second dialect.

import (
	"regexp"
	"strconv"
	"strings"
)

// EditSegmentKind discriminates one parsed template segment.
type EditSegmentKind int

const (
	// EditLiteral emits its own bytes.
	EditLiteral EditSegmentKind = iota
	// EditGroup emits the anchor capture group's matched text, or the
	// empty string where the group did not participate.
	EditGroup
	// EditKey emits the planned value of the key the rule table is named
	// for (in `replace`), or the bound tag value (in `anchor`, where it
	// is emitted regexp-QUOTED).
	EditKey
)

// EditSegment is one parsed template segment.
type EditSegment struct {
	Kind EditSegmentKind
	// Text is the literal bytes for EditLiteral, the key name for
	// EditKey, and the named group's name for a named EditGroup.
	Text string
	// Index is the numbered group for EditGroup, or -1 for a named one.
	Index int
}

// editTagPrefix is the ONE prefix that opens a placeholder in an anchor.
const editTagPrefix = "{tag."

// ParseEditReplace parses a `replace` template over its closed
// vocabulary. `owned` is the one key the rule's own `{<key>}` may name;
// `groups` and `names` are the anchor's capture-group arity and names, so
// a reference the anchor does not define is refused here rather than
// expanding empty at apply.
//
// A nil `names` and a negative `groups` skip the group-arity check, which
// is what an apply-time re-parse of an already-validated template wants:
// the anchor is compiled by then and `regexp` answers the arity itself.
func ParseEditReplace(tmpl, owned string, groups int, names []string) ([]EditSegment, error) {
	var out []EditSegment
	var lit strings.Builder
	flush := func() {
		if lit.Len() != 0 {
			out = append(out, EditSegment{Kind: EditLiteral, Text: lit.String()})
			lit.Reset()
		}
	}

	for i := 0; i < len(tmpl); {
		switch tmpl[i] {
		case '$':
			// `$$` is the escape; `${…}` is a group; a bare `$1`/`$name`
			// is refused, which is exactly where this parser and
			// `regexp.Expand` part company.
			if i+1 < len(tmpl) && tmpl[i+1] == '$' {
				lit.WriteByte('$')
				i += 2
				continue
			}
			if i+1 >= len(tmpl) || tmpl[i+1] != '{' {
				return nil, editTemplateErr("the token `$` at offset " +
					strconv.Itoa(i) + " is not `$$` and does not open `${…}`; " +
					"a bare `$1`/`$name` is not in this template's vocabulary")
			}
			end := strings.IndexByte(tmpl[i+2:], '}')
			if end < 0 {
				return nil, editTemplateErr("the `${` at offset " +
					strconv.Itoa(i) + " is never closed")
			}
			ref := tmpl[i+2 : i+2+end]
			seg, err := editGroupSegment(ref, groups, names)
			if err != nil {
				return nil, err
			}
			flush()
			out = append(out, seg)
			i += 2 + end + 1

		case '{':
			// `{{` emits a literal brace; otherwise the token must be
			// exactly the rule's own key.
			if i+1 < len(tmpl) && tmpl[i+1] == '{' {
				lit.WriteByte('{')
				i += 2
				continue
			}
			end := strings.IndexByte(tmpl[i+1:], '}')
			if end < 0 {
				return nil, editTemplateErr("the `{` at offset " +
					strconv.Itoa(i) + " is never closed")
			}
			key := tmpl[i+1 : i+1+end]
			if key != owned {
				return nil, editTemplateErr("the placeholder `{" + key +
					"}` is not this rule's own key `{" + owned + "}`; a `replace` " +
					"admits the planned value of THE key its table is named for " +
					"and no other")
			}
			flush()
			out = append(out, EditSegment{Kind: EditKey, Text: key, Index: -1})
			i += 1 + end + 1

		case '}':
			// `}}` emits a literal brace. A lone `}` is ordinary text —
			// it opens nothing, so there is no malformed token to name.
			if i+1 < len(tmpl) && tmpl[i+1] == '}' {
				lit.WriteByte('}')
				i += 2
				continue
			}
			lit.WriteByte('}')
			i++

		default:
			lit.WriteByte(tmpl[i])
			i++
		}
	}
	flush()
	return out, nil
}

// editGroupSegment resolves one `${…}` reference against the anchor's
// capture groups. A reference the anchor does not define is a defect at
// LOAD, so it never reaches apply to expand silently empty — which is the
// non-participating-group rule, a different case entirely.
func editGroupSegment(ref string, groups int, names []string) (EditSegment, error) {
	if ref == "" {
		return EditSegment{}, editTemplateErr("`${}` names no group")
	}
	if n, err := strconv.Atoi(ref); err == nil {
		if n < 1 {
			return EditSegment{}, editTemplateErr("`${" + ref +
				"}` is not a capture group; groups are 1-based")
		}
		if groups >= 0 && n > groups {
			return EditSegment{}, editTemplateErr("`${" + ref +
				"}` references a group the anchor does not define; it defines " +
				strconv.Itoa(groups))
		}
		return EditSegment{Kind: EditGroup, Text: ref, Index: n}, nil
	}
	if names != nil {
		found := false
		for _, n := range names {
			if n == ref {
				found = true
				break
			}
		}
		if !found {
			return EditSegment{}, editTemplateErr("`${" + ref +
				"}` names a group the anchor does not define")
		}
	}
	return EditSegment{Kind: EditGroup, Text: ref, Index: -1}, nil
}

// ParseEditAnchor parses an `anchor` over its own vocabulary: the fixed
// prefix `{tag.` opens a placeholder, scanned to its closing `}`, and
// EVERY other byte — brace, dollar, caret included — is RE2's.
//
// `declared` reports whether a tag key is bindable at invocation. It is
// nil at apply, where the model has already been linted and the keys are
// whatever the invocation bound.
func ParseEditAnchor(pattern string, declared func(string) bool) ([]EditSegment, error) {
	var out []EditSegment
	var lit strings.Builder
	flush := func() {
		if lit.Len() != 0 {
			out = append(out, EditSegment{Kind: EditLiteral, Text: lit.String()})
			lit.Reset()
		}
	}

	for i := 0; i < len(pattern); {
		// The scan is for the UNESCAPED literal prefix. RE2's own escape
		// (`\{tag\.`) therefore stays ordinary pattern text, so a line
		// whose CONTENT carries `{tag.` remains anchorable. A backslash
		// consumes the byte after it for exactly this reason.
		if pattern[i] == '\\' {
			lit.WriteByte(pattern[i])
			if i+1 < len(pattern) {
				lit.WriteByte(pattern[i+1])
				i += 2
				continue
			}
			i++
			continue
		}
		if !strings.HasPrefix(pattern[i:], editTagPrefix) {
			// "carries any other `{…}` form" (`0028:C1.4`), reconciled
			// with "never for a brace RE2 itself accepts" (`0028:C1.2`).
			//
			// RE2 accepts an unescaped `{` either way: `{artifact}` and
			// `\d{4}` both COMPILE, the first as literal text and the
			// second as a bounded quantifier, so compilation cannot tell
			// them apart and the two clauses would be unsatisfiable
			// together if it had to. The discriminator that satisfies
			// both is the brace's SHAPE: a `{` whose contents are a
			// well-formed repeat spec (`{N}`, `{N,}`, `{N,M}`) is
			// quantifier syntax and reaches RE2 untouched; any other
			// unescaped `{…}` is the "other form" C1.4 refuses, which is
			// what stops `{artifact}` or an entry's own `{status}` from
			// silently anchoring on literal braces the document does not
			// carry.
			if pattern[i] == '{' && !editRepeatSpec(pattern[i:]) {
				return nil, editAnchorErr("the `{` at offset " +
					strconv.Itoa(i) + " opens neither a `" + editTagPrefix +
					"<key>}` placeholder nor an RE2 repeat quantifier; escape " +
					"it as `\\{` to anchor a literal brace")
			}
			lit.WriteByte(pattern[i])
			i++
			continue
		}
		end := strings.IndexByte(pattern[i+len(editTagPrefix):], '}')
		if end < 0 {
			return nil, editAnchorErr("the `" + editTagPrefix + "` at offset " +
				strconv.Itoa(i) + " is never closed")
		}
		key := pattern[i+len(editTagPrefix) : i+len(editTagPrefix)+end]
		if key == "" {
			return nil, editAnchorErr("`{tag.}` names no tag key")
		}
		if declared != nil && !declared(key) {
			return nil, editAnchorErr("`{tag." + key +
				"}` names no tag key the model declares with `observed` " +
				"provenance; a key of any other provenance is structurally " +
				"unbindable at invocation")
		}
		flush()
		out = append(out, EditSegment{Kind: EditKey, Text: key, Index: -1})
		i += len(editTagPrefix) + end + 1
	}
	flush()
	return out, nil
}

// editRepeatSpec reports whether s opens an RE2 repeat quantifier —
// `{N}`, `{N,}` or `{N,M}`. It reads only the shape, never the operand it
// applies to: whether the quantifier has something to repeat is RE2's
// question and it answers it at compile.
func editRepeatSpec(s string) bool {
	if len(s) == 0 || s[0] != '{' {
		return false
	}
	end := strings.IndexByte(s, '}')
	if end < 0 {
		return false
	}
	body := s[1:end]
	lo, hi, comma := strings.Cut(body, ",")
	if !editDigits(lo) {
		return false
	}
	return !comma || hi == "" || editDigits(hi)
}

// editDigits reports whether s is one or more ASCII digits.
func editDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := range len(s) {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// EditAnchorProbe renders an anchor's segments with every `{tag.<key>}`
// replaced by a quoted probe, so the pattern can be compiled at LOAD
// without any binding. `regexp.QuoteMeta` emits no group syntax, so the
// probe cannot change the compiled arity the `replace` check reads
// (`0028:C1.4`).
func EditAnchorProbe(segs []EditSegment) string {
	return expandAnchor(segs, func(string) (string, bool) { return "", true })
}

// ExpandEditAnchor renders an anchor's segments with each `{tag.<key>}`
// replaced by its bound value, regexp-QUOTED — so a bound value is a
// literal-match fragment and can never alter the pattern's structure.
// It reports the first key it could not bind.
func ExpandEditAnchor(segs []EditSegment, bound map[string]string) (string, string, bool) {
	var missing string
	out := expandAnchor(segs, func(key string) (string, bool) {
		v, ok := bound[key]
		if !ok && missing == "" {
			missing = key
		}
		return v, ok
	})
	if missing != "" {
		return "", missing, false
	}
	return out, "", true
}

func expandAnchor(segs []EditSegment, bind func(string) (string, bool)) string {
	var b strings.Builder
	for _, s := range segs {
		if s.Kind == EditKey {
			v, _ := bind(s.Text)
			b.WriteString(regexp.QuoteMeta(v))
			continue
		}
		b.WriteString(s.Text)
	}
	return b.String()
}

// EditAnchorTagKeys returns the tag keys an anchor's segments reference,
// which is the per-USE-SITE scope both the multiline scan and C1.6's
// argv rule are held to: a tag bound on the context but referenced by no
// anchor of this entry is never scanned (`0028:C1.3` precedence: (2)).
func EditAnchorTagKeys(segs []EditSegment) []string {
	var out []string
	for _, s := range segs {
		if s.Kind == EditKey {
			out = append(out, s.Text)
		}
	}
	return out
}

// EditClearLine is the ONE member of `clear`'s closed set (`0028:C1.5`).
const EditClearLine = "line"

// --- the two template error carriers -------------------------------------

// editTemplateError and editAnchorError let the loader map a parse defect
// onto its own category while the apply side, which has no categories,
// reads them as ordinary errors.
type editTemplateError struct{ detail string }

func (e *editTemplateError) Error() string { return e.detail }

type editAnchorError struct{ detail string }

func (e *editAnchorError) Error() string { return e.detail }

func editTemplateErr(detail string) error { return &editTemplateError{detail: detail} }

func editAnchorErr(detail string) error { return &editAnchorError{detail: detail} }
