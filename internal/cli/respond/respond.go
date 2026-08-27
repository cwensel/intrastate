// Package respond is the CLI's single output gateway. Verbs call OK,
// Fail, Note, and Warn rather than printing to stdout/stderr directly,
// so machine-friendly JSON output and human-friendly text output are
// both derived from one render path keyed off the persistent
// --as=text|json root flag.
//
// # The contract
//
// Under --as=json, stdout carries the terminal disposition as one
// NDJSON object discriminated by a "type" field:
//
//	{"type": "ok",     ...}   # success
//	{"type": "failed", ...}   # graceful failure
//
// Exactly one terminal record is emitted on every graceful exit; its
// absence means the process was killed. Stderr under --as=json carries
// advisories only, discriminated by "level" ("note" | "warning").
//
// Under --as=text, stdout is verb-defined human output and stderr
// carries "error: …", "note: …", and "warning: …" lines.
//
// As the CLI grows, add an intermediate-record emitter (Stream) for
// verbs that produce zero-or-more records before the terminal line. The
// terminal "ok"/"failed" type names are reserved.
package respond

import (
	"fmt"
	"io"
	"strings"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Mode is the wire format chosen by --as.
type Mode int

const (
	ModeText Mode = iota
	ModeJSON
)

// FlagName is the persistent root flag registered in root.go.
const FlagName = "as"

// Success is the terminal-success envelope under --as=json. Type is
// always "ok" (Fail emits "failed"). Data carries the verb-specific
// payload; Notes and Warnings are advisory channels.
type Success struct {
	Type     string     `json:"type"`
	Notes    []Advisory `json:"notes,omitempty"`
	Warnings []Warning  `json:"warnings,omitempty"`
	Data     any        `json:"data,omitempty"`
}

// TextLiner is the interface a verb's Data payload satisfies when its
// whole content is one canonical line — a build identity, say — rather
// than a set of fields worth enumerating.
//
// Like FindingCarrier it is an interface rather than a field on Success,
// so the gateway still gains no knowledge of any verb's payload: a verb
// opts in by having a TextLine method, and the gateway asks.
type TextLiner interface {
	TextLine() string
}

// FindingCarrier is the interface a verb's Data payload satisfies when it
// carries structured findings the text branch must enumerate.
//
// It is an interface rather than a field on Success so the gateway gains no
// dependency on any producing subsystem: `clierr.Finding` is
// subsystem-agnostic by construction, and a payload that carries none simply
// does not implement this.
type FindingCarrier interface {
	// LintFindings returns the findings the text branch renders. The JSON
	// branch never calls it — the payload marshals itself.
	LintFindings() []clierr.Finding
}

// Advisory is one structured non-fatal note. Code is the stable
// machine-readable identifier; Message is the human phrasing.
type Advisory struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

// Warning is a non-fatal advisory rendered as "warning: …" in text mode
// or carried inside Success.Warnings in JSON mode.
type Warning struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
	Hint    string `json:"hint,omitempty"`
}

// ModeOf reads the persistent --as flag from the root command. Unknown
// values default to ModeText; ValidateMode is the strict counterpart
// that refuses them.
func ModeOf(cmd *cobra.Command) Mode {
	flag := lookupAs(cmd)
	if flag == nil {
		return ModeText
	}
	if strings.ToLower(flag.Value.String()) == "json" {
		return ModeJSON
	}
	return ModeText
}

// ValidateMode returns a CLIError when --as carries an unrecognized
// value. Verbs call this at the top of RunE to fail fast on `--as=yaml`
// and similar typos.
func ValidateMode(cmd *cobra.Command) *clierr.CLIError {
	flag := lookupAs(cmd)
	if flag == nil {
		return nil
	}
	switch strings.ToLower(flag.Value.String()) {
	case "", "text", "json":
		return nil
	}
	return &clierr.CLIError{
		Code:    "flag-invalid-value",
		Param:   FlagName,
		Message: "--as must be one of text|json; got " + flag.Value.String(),
		Group:   clierr.GroupUserEnv,
	}
}

// OK emits the terminal-success envelope on stdout and any advisories
// on stderr (text mode) or folded into the envelope (json mode). It is
// the natural last act of a verb's RunE on the success path.
func OK(cmd *cobra.Command, s Success) error {
	s.Type = "ok"
	switch ModeOf(cmd) {
	case ModeJSON:
		return writeJSONLine(cmd.OutOrStdout(), s)
	default:
		stderr := cmd.ErrOrStderr()
		for _, n := range s.Notes {
			_, _ = fmt.Fprintf(stderr, "note: %s\n", n.Message)
		}
		for _, w := range s.Warnings {
			writeTextWarning(stderr, w)
		}
		// The text branch renders Notes and Warnings, and — for a payload
		// that carries them — every structured finding. Dropping Data
		// wholesale would make the two modes disagree about which findings a
		// run reported, which the output contract forbids (`0006:C14`).
		if carrier, ok := s.Data.(FindingCarrier); ok {
			clierr.EmitFindingsText(cmd.OutOrStdout(), carrier.LintFindings())
			return nil
		}
		// A payload whose content is ONE canonical line renders as that
		// line. The generic renderer below would spread it over a sorted
		// leaf per field, which is right for a structured result and wrong
		// for an identity string a caller pipes or greps.
		//
		// It is the same content either way — TextLine is derived from the
		// same value the JSON branch marshals — so this narrows the
		// FORMAT, which A-11 leaves free, and not the content, which
		// REQ-11 fixes.
		if liner, ok := s.Data.(TextLiner); ok {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), liner.TextLine())
			return nil
		}
		// Any other payload is rendered field-for-field from the SAME value
		// the JSON branch marshals, so the two modes cannot disagree about
		// what the run reported (RDR 0005 REQ-7/REQ-9/REQ-11/REQ-120).
		writeTextPayload(cmd.OutOrStdout(), s.Data)
		return nil
	}
}

// Fail emits the terminal-error envelope and returns ce so a verb's
// RunE can `return respond.Fail(cmd, ce)`. In text mode the envelope
// goes to stderr; in json mode it goes to stdout so a single stream
// carries both ok and failed dispositions.
func Fail(cmd *cobra.Command, ce *clierr.CLIError) *clierr.CLIError {
	if ce == nil {
		return nil
	}
	switch ModeOf(cmd) {
	case ModeJSON:
		clierr.EmitJSON(cmd.OutOrStdout(), ce)
	default:
		clierr.EmitText(cmd.ErrOrStderr(), ce)
	}
	return ce
}

// Note emits a free-form advisory: "note: …" on stderr in text mode, a
// {"level":"note", …} line on stderr in json mode.
//
// No verb emits an advisory yet, so static analysis reports this and Warn
// below as unreachable. They are kept deliberately: the advisory stream is
// a specified tier of the wire contract (docs/cli-output-contract.md,
// RDR 0005), and this is the gateway API that tier is served through.
// Deleting them would delete a documented surface and push the next verb
// that needs an advisory toward a direct print — the exact drift the
// gateway exists to prevent.
func Note(cmd *cobra.Command, message string) {
	stderr := cmd.ErrOrStderr()
	if ModeOf(cmd) == ModeJSON {
		_ = writeJSONLine(stderr, map[string]string{"level": "note", "message": message})
		return
	}
	_, _ = fmt.Fprintf(stderr, "note: %s\n", message)
}

// Warn emits a structured warning on stderr in both modes.
func Warn(cmd *cobra.Command, w Warning) {
	stderr := cmd.ErrOrStderr()
	if ModeOf(cmd) == ModeJSON {
		payload := map[string]string{"level": "warning", "message": w.Message}
		if w.Code != "" {
			payload["code"] = w.Code
		}
		_ = writeJSONLine(stderr, payload)
		return
	}
	writeTextWarning(stderr, w)
}

func lookupAs(cmd *cobra.Command) *pflag.Flag {
	if cmd == nil {
		return nil
	}
	root := cmd.Root()
	if root == nil {
		return nil
	}
	return root.PersistentFlags().Lookup(FlagName)
}

// writeJSONLine routes through the ONE shared non-HTML-escaping encoder
// `clierr` owns, so this success site and the failure site render a
// canonical set literal identically (`0005:C1`, REQ-71/REQ-72).
func writeJSONLine(out io.Writer, v any) error {
	return clierr.WriteJSONLine(out, v)
}

func writeTextWarning(out io.Writer, w Warning) {
	_, _ = fmt.Fprintf(out, "warning: %s\n", w.Message)
	if w.Detail != "" {
		_, _ = fmt.Fprintf(out, "  detail: %s\n", w.Detail)
	}
	if w.Hint != "" {
		_, _ = fmt.Fprintf(out, "  hint: %s\n", w.Hint)
	}
}
