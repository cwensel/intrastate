package respond

// Text rendering of a verb's success payload (RDR 0005, REQ-7/REQ-9/REQ-11).
//
// The renderer is respond-owned and driven by respond.OK, so a verb never
// prints its payload itself — `internal/cli::newVersionCmd`'s `cmd.Println`
// is the named drift this exists to avoid (REQ-10). A verb that printed
// directly would be mode-blind and would put a second record on stdout
// under `--as=json`.
//
// # Why it renders the payload's JSON shape rather than a hand-written form
//
// REQ-11 fixes two properties: text "renders the same result content" and
// "does not invent fields absent from the JSON payload". REQ-120 adds the
// mechanism — "derive both renderings from one typed result per verb".
//
// A hand-written text template satisfies neither for long. It is a SECOND
// description of the payload, and the moment a verb gains a field the
// template does not, the two modes disagree about what the run reported —
// silently, because nothing fails. Rendering the marshalled form of the
// same value makes drift impossible by construction: there is one source
// of content, and text is a projection of it.
//
// Layout is deliberately free (A-11): the contract fixes content
// completeness and the one-finding-per-line rule, not a format.

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// writeTextPayload renders a verb's success payload as human-scannable
// lines: one leaf value per line, path-qualified.
//
// It round-trips through the payload's JSON encoding on purpose — that is
// what makes "the same result content" literally the same content, and it
// is the same encoder the JSON branch uses, so a canonical set literal
// renders identically in both modes.
func writeTextPayload(out io.Writer, data any) {
	if data == nil {
		return
	}
	buf, err := marshalCanonical(data)
	if err != nil {
		return
	}
	var decoded any
	if err := json.Unmarshal(buf, &decoded); err != nil {
		return
	}
	for _, line := range flatten("", decoded) {
		_, _ = fmt.Fprintln(out, line)
	}
}

// flatten renders one decoded JSON value as path-qualified leaf lines.
//
// An empty object or array still emits a line for its path. Absence and
// emptiness are different claims in this contract — a cleared key is absent
// while an empty set reads back `[]` — so a renderer that dropped empty
// containers would erase exactly the distinction REQ-66 turns on.
func flatten(path string, v any) []string {
	switch t := v.(type) {
	case map[string]any:
		if len(t) == 0 {
			return []string{label(path) + "(none)"}
		}
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var out []string
		for _, k := range keys {
			out = append(out, flatten(join(path, k), t[k])...)
		}
		return out

	case []any:
		if len(t) == 0 {
			return []string{label(path) + "(none)"}
		}
		var out []string
		for i, e := range t {
			out = append(out, flatten(fmt.Sprintf("%s[%d]", path, i), e)...)
		}
		return out

	case nil:
		return []string{label(path) + "(none)"}

	case bool:
		return []string{label(path) + fmt.Sprintf("%t", t)}

	case float64:
		// json.Unmarshal renders every number as float64; render integral
		// values without a spurious fractional part.
		if t == float64(int64(t)) {
			return []string{label(path) + fmt.Sprintf("%d", int64(t))}
		}
		return []string{label(path) + fmt.Sprintf("%g", t)}

	case string:
		return []string{label(path) + t}

	default:
		return []string{label(path) + fmt.Sprintf("%v", t)}
	}
}

func join(path, key string) string {
	if path == "" {
		return key
	}
	return path + "." + key
}

func label(path string) string {
	if path == "" {
		return ""
	}
	return path + ": "
}

// marshalCanonical encodes v with the same non-HTML-escaping settings the
// wire sites use, so `<` and `&` render as themselves in text mode too.
func marshalCanonical(v any) ([]byte, error) {
	var buf strings.Builder
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return []byte(strings.TrimRight(buf.String(), "\n")), nil
}
