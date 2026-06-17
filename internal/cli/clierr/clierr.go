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
// followed by indented detail and hint lines when present.
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
}
