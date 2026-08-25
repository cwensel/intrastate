// Package clierr is the leaf home of the CLI's structured-error type.
// It lives in its own package so other internal packages (config, …)
// can construct CLIErrors without importing internal/cli and forming an
// import cycle.
//
// The design intent carried over into the scaffold: every user-facing
// failure is a *CLIError carrying a stable machine-readable Code, a
// human Message, and an ErrorGroup that maps to a process exit code.
// Verbs route every error through the respond gateway so the CLI is
// never silent — a non-zero exit is always accompanied by a structured
// envelope on the wire.
package clierr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ErrorGroup classifies a CLIError for exit-code mapping. The numeric
// exit contract a script can branch on lives in ExitCodeFor; the groups
// below are the named buckets that feed it. Add groups as the exit
// taxonomy grows — keep ExitCodeFor in sync.
type ErrorGroup int

const (
	// GroupSuccess and GroupWarning exit 0 — the command completed.
	GroupSuccess ErrorGroup = iota
	GroupWarning
	// GroupUserEnv maps to exit 2: bad input, bad flags, or a
	// not-found in the user's environment. The common refusal bucket.
	GroupUserEnv
	// GroupEnvUnavailable maps to exit 3: a required external facility
	// (editor, network, …) was unavailable through no fault of input.
	GroupEnvUnavailable
	// GroupInternal maps to exit 2: an internal/IO error that is not a
	// user-input refusal but also not an environment unavailability.
	GroupInternal
	// GroupSignalCancel maps to exit 130: interrupted (SIGINT).
	GroupSignalCancel
)

// CLIError is the single CLI-side structured error type. The JSON form
// is the on-the-wire envelope; Group drives the exit code and is not
// serialized. Extend with new optional fields as needed — keep them
// `omitempty` so the envelope stays append-only and stable for tools.
type CLIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	// Param names the offending flag/argument, when the failure is
	// attributable to one.
	Param string `json:"param,omitempty"`
	// Detail carries hard facts about the failure (the underlying
	// syscall reason, a parser diagnostic). May be multi-line.
	Detail string `json:"detail,omitempty"`
	// Hint is an optional one-line remedy.
	Hint string `json:"hint,omitempty"`

	// Findings carries the individual structured findings an aggregate
	// failure reports. It is the ONE exception to this envelope's
	// `omitempty` habit: the key must serialize even when the list is
	// empty, because an empty list is the proof's receipt (`0006:C14`).
	// It is a top-level SIBLING of `code` on the marshalled CLIError —
	// EmitJSON marshals the *CLIError itself, so there is no `error`
	// wrapper to nest under.
	Findings []Finding `json:"findings"`

	// Group selects the exit code; not serialized.
	Group ErrorGroup `json:"-"`

	// Cause preserves the underlying Go error for errors.Is/errors.As
	// traversal. Not serialized — the wire-visible cause surface is
	// Detail.
	Cause error `json:"-"`
}

func (e *CLIError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Message == "" {
		return e.Code
	}
	return e.Code + ": " + e.Message
}

// Unwrap exposes the wrapped Cause so errors.Is / errors.As can
// traverse the chain.
func (e *CLIError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// ErrorCode returns the structured code carried by err, or "" if err is
// not a recognized CLIError.
func ErrorCode(err error) string {
	if err == nil {
		return ""
	}
	var ce *CLIError
	if errors.As(err, &ce) {
		return ce.Code
	}
	return ""
}

// ExitCodeFor returns the process exit code for err. The contract a
// caller scripts against:
//
//	0  success / warning
//	2  user or internal error (bad input, IO failure)
//	3  environment unavailable
//	130 interrupted
//
// A non-CLIError defaults to exit 1 (unexpected/unclassified).
func ExitCodeFor(err error) int {
	if err == nil {
		return 0
	}
	var ce *CLIError
	if errors.As(err, &ce) {
		switch ce.Group {
		case GroupSuccess, GroupWarning:
			return 0
		case GroupUserEnv, GroupInternal:
			return 2
		case GroupEnvUnavailable:
			return 3
		case GroupSignalCancel:
			return 130
		}
	}
	return 1
}

// EmitJSON writes the structured CLIError envelope as one NDJSON line.
func EmitJSON(out io.Writer, e *CLIError) {
	if e == nil {
		return
	}
	if buf, err := json.Marshal(e); err == nil {
		_, _ = fmt.Fprintln(out, string(buf))
	}
}

// EmitText writes a human-readable "error: <code>: <message>" line,
// followed by indented detail and hint lines when present, and then EVERY
// carried finding's code and message.
//
// Text mode enumerates every finding rather than summarising: the set of
// finding codes in text output must equal the set in JSON output for the
// same input, so a renderer that drops one has dropped a defect the author
// needs to see. Layout and wrapping are free; completeness is not
// (`0006:C14`).
func EmitText(out io.Writer, e *CLIError) {
	if e == nil {
		return
	}
	_, _ = fmt.Fprintf(out, "error: %s\n", e.Error())
	if e.Detail != "" {
		_, _ = fmt.Fprintf(out, "  detail: %s\n", e.Detail)
	}
	if e.Hint != "" {
		_, _ = fmt.Fprintf(out, "  hint: %s\n", e.Hint)
	}
	EmitFindingsText(out, e.Findings)
}

// EmitFindingsText renders one finding per line as "  <code>: <message>",
// with the identity fields the finding carries appended so a reader can act
// on it without re-running in JSON mode.
//
// Identity values reach here from author TOML (`Literal` is the joined atom
// literal), and nothing on the load path rejects a newline or a `) (` in one.
// They are therefore quoted, so no authored value can split a finding across
// lines or forge a second finding's identity suffix — one finding is one line.
func EmitFindingsText(out io.Writer, findings []Finding) {
	for _, f := range findings {
		_, _ = fmt.Fprintf(out, "  %s: %s%s\n",
			f.Code, sanitizeLine(f.Message), f.identitySuffix())
	}
}

// sanitizeLine folds every line-affecting character in producer-generated
// prose to a space so one finding stays one line. Messages embed
// author-controlled tag names and values — `nodeElement` composes `Node.key()`
// into a message with `%s`, not `%q`, so such a character reaches the line
// unescaped. Unlike identity values a message is prose meant to read unquoted,
// and it is asserted on as a plain substring, so it is flattened rather than
// quoted.
//
// CR and LF are folded as a pair first so a CRLF becomes one space rather than
// two. The rest are folded individually: VT and FF advance a line on a
// terminal, and NEL (U+0085), LS (U+2028) and PS (U+2029) are line breaks to
// Unicode-aware consumers.
func sanitizeLine(s string) string {
	if !strings.ContainsAny(s, "\r\n\v\f\u0085\u2028\u2029") {
		return s
	}
	return strings.NewReplacer(
		"\r\n", " ",
		"\r", " ",
		"\n", " ",
		"\v", " ",
		"\f", " ",
		"\u0085", " ",
		"\u2028", " ",
		"\u2029", " ",
	).Replace(s)
}

// identitySuffix renders the finding's identity fields as a trailing
// parenthesised list of quoted values, or "" when it carries none. Quoting is
// what keeps the list unforgeable by an authored value; see EmitFindingsText.
func (f Finding) identitySuffix() string {
	var parts []string
	for _, p := range [][2]string{
		{"rule", f.Rule},
		{"span", f.Span},
		{"element", f.Element},
		{"reason", f.Reason},
		{"dimension", f.Dimension},
		{"key", f.Key},
		{"operator", f.Operator},
		{"literal", f.Literal},
		{"block", f.Block},
		{"class", f.Class},
	} {
		if p[1] != "" {
			parts = append(parts, p[0]+"="+strconv.Quote(p[1]))
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return " (" + strings.Join(parts, " ") + ")"
}

// --- graph-lint findings (RDR 0006) --------------------------------------
//
// SKELETON (RDR 0006 Stage 8, Phase 1). The Finding record is defined
// HERE, in the leaf error package, so `clierr` gains no dependency on the
// graph-lint package (`0006:C14`). The atom and class fields are declared
// STRING-typed, not enums: `Block` is RDR 0002 authoring vocabulary and
// `Class` is RDR 0001's RefusalKind, and importing either type would give
// `clierr` exactly the dependency this clause prevents. `clierr` ascribes
// these fields no meaning — it only transports and serializes them; the
// producing package owns the vocabulary.

// Finding is one graph-lint finding as it crosses the wire. It is
// subsystem-agnostic by construction.
type Finding struct {
	Code     string `json:"code"`
	Model    string `json:"model"`
	Severity string `json:"severity"`
	Message  string `json:"message"`

	// Rule names the source rule/context id when the normalized model can
	// provide one; Span and Element carry the fallbacks.
	Rule    string `json:"rule,omitempty"`
	Span    string `json:"span,omitempty"`
	Element string `json:"element,omitempty"`

	// Reason is the discriminator `graph-unprovable-coverage` carries.
	Reason string `json:"reason,omitempty"`
	// Dimension names the dimension a per-dimension finding is about.
	Dimension string `json:"dimension,omitempty"`

	// Key/Operator/Literal/Block are the atom fields an atom-attributed
	// finding carries. Declared string, never an enum.
	Key      string `json:"key,omitempty"`
	Operator string `json:"operator,omitempty"`
	Literal  string `json:"literal,omitempty"`
	Block    string `json:"block,omitempty"`
	// Class is the declared failure class an escape-scoped finding carries.
	Class string `json:"class,omitempty"`

	// Fingerprint is the canonical sortable predicate/write serialization
	// the finding-identity tuple closes on.
	Fingerprint string `json:"fingerprint,omitempty"`
}
