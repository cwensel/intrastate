package cli

// RDR 0021 — the `graph` verb (`0021:C1`, `0021:C5`).
//
// A new ROOT verb beside `lint`, outside the `flow` group, under
// `0005:C1`'s carve-out for command groups owned by the RDR that names
// them. It mirrors lint's selection ARM SET and its CODE SPELLINGS
// verbatim and never redefines them — lint owns those codes — while its
// help text is its own, because a verb that described itself as lint does
// would tell a caller the wrong thing about what it produces.
//
// It runs load, normalization, grouping, and the reachability traversal
// ONLY. It never runs the lint invariants and emits no findings: an export
// is never a lint pass, and the lint gate stays the acceptance surface
// (`0021:C4`).

import (
	"bytes"
	"os"
	"slices"

	"github.com/spf13/cobra"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/cli/respond"
	"github.com/cwensel/intrastate/internal/graphlint"
	"github.com/cwensel/intrastate/internal/table"
)

// emitFlagName is the document-format selector. `--emit` selects the
// DOCUMENT; `--as` selects the ENVELOPE (`0021:D-selection-predicate`).
const emitFlagName = "emit"

const (
	emitJSON = "json"
	emitDOT  = "dot"
)

// codeExportTooLarge is the scalar refusal an incomplete traversal takes.
// It names one subject, so it carries no findings list.
const codeExportTooLarge = "graph-export-too-large"

// emitFormats is the `--emit` vocabulary, an `append-only` set
// (`0029:C2`): a third format MAY be added in a minor, none removed or
// renamed within a major. `json` is the default and leads the list.
var emitFormats = []string{emitJSON, emitDOT}

// EmitFormats returns the `--emit` document-format vocabulary.
//
// It is the enumeration seam `0029:C4` obliges of every tier assignment
// (REQ-10). It is sited in `internal/cli` because this package owns the
// set and nothing it imports imports it back, so the accessor builds with
// no new import and no cycle. Without it the `append-only` tier would bind
// only prose, with nothing to check membership over.
func EmitFormats() []string { return slices.Clone(emitFormats) }

// graphPayload is the verb's success payload.
//
// It satisfies `respond.TextLiner`, so under `--as=text` the gateway
// prints the document and nothing else — the DOT stream pipes to `dot` and
// the JSON document diffs raw in CI. The generic renderer would spread it
// over a sorted leaf per field, which is right for a structured result and
// wrong for a document a caller pipes (`0021:C5`, REQ-64).
type graphPayload struct {
	// document is the rendered document for the text branch.
	document string
	// data is what the JSON branch marshals under `data`: the document
	// OBJECT for `--emit json`, and for `--emit dot` an object carrying the
	// DOT text as the single required string member `dot`, so the
	// documented unwrap is exactly `jq -r .data.dot` (REQ-58).
	data any
}

// TextLine returns the bare document. The text stream is the document by
// contract, exactly as `version`'s TextLine is its identity string: a
// consumer that wants an envelope must use `--as=json` (REQ-61).
func (p graphPayload) TextLine() string { return p.document }

// MarshalJSON emits the payload's data value, so both modes derive from
// ONE export value (REQ-59).
func (p graphPayload) MarshalJSON() ([]byte, error) {
	return marshalGraphJSON(p.data)
}

// marshalGraphJSON encodes a value through the ONE shared
// non-HTML-escaping encoder (`0005:C1`'s one-encoder rule, REQ-43), so
// `<`, `>`, and `&` reach the wire as themselves and no fifth ad hoc
// `SetEscapeHTML(false)` site is minted beside the four A5 records.
//
// `json.Encoder.Encode` terminates its value with a newline. That newline
// is the ENCODER's, not the document's: the gateway's own `Fprintln`
// supplies the single trailing newline F1 pins on text-mode stdout, and
// the `dot` member carries the document without it (REQ-78). Trimming it
// here is what keeps the document value newline-free in both modes.
func marshalGraphJSON(v any) ([]byte, error) {
	var buf bytes.Buffer
	if err := clierr.WriteJSONLine(&buf, v); err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), nil
}

// dotPayload carries the DOT text under its normative member spelling.
type dotPayload struct {
	Dot string `json:"dot"`
}

func newGraphCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "graph",
		Short: "Export a model's normalized graph as a document",
		Long: `Export a transition model's normalized graph as a versioned document.

The document is a documentation artifact, not a verdict: this command
runs load, normalization, grouping, and the reachability traversal only.
It never runs the graph invariants and emits no findings, so a model
` + "`lint`" + ` refuses still exports — which is what makes a failing model
inspectable as a graph. Run ` + "`intrastate lint`" + ` for acceptance.

  --emit ` + emitJSON + `   the ` + schemaMarkerValue + ` JSON document (default)
  --emit ` + emitDOT + `    a DOT digraph of the same relation

Under --as=text stdout carries the selected document verbatim, so the
DOT stream pipes to ` + "`dot`" + ` and the JSON document diffs raw in CI.
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE:          runGraphExport,
	}
	cmd.Flags().String("model", "", "path to the transition model to export")
	cmd.Flags().String("flow", "", "flow id to export (reserved; this build resolves none — use --model)")
	cmd.Flags().String(emitFlagName, emitJSON,
		"document format: "+emitJSON+" | "+emitDOT)
	withExtendedHelp(cmd, graphExtendedDesc())
	return cmd
}

// graphExtendedDesc is the --help-all body: what the document IS, how it
// is versioned, and the soundness rule a consumer needs before reasoning
// over the relation.
func graphExtendedDesc() string {
	var b []byte
	add := func(s string) { b = append(b, s...) }

	add(`This command exports the model's normalized graph as a document.

The JSON document is versioned by its leading ` + "`schema`" + ` field, initial
value ` + schemaMarkerValue + `. Evolution within /1 is ADDITIVE: a member may be
added in a minor, none removed or renamed within a major, so a consumer
must tolerate an unrecognized field and must not assert on the field
set's size or on any member's position.

Two version markers coexist and neither substitutes for the other. The
document's own ` + "`schema`" + ` versions the DOCUMENT; the envelope's
` + "`schema_version`" + ` versions the ENVELOPE carrying it and is never
projected into ` + "`data`" + `.

Document formats (` + "`--emit`" + `, an append-only vocabulary):

`)
	for _, f := range emitFormats {
		add("  " + f + "\n")
	}

	add(`
` + "`--emit`" + ` selects the DOCUMENT and ` + "`--as`" + ` selects the ENVELOPE; all four
combinations are defined. Under --as=text stdout is the bare document
plus one newline. Under --as=json stdout is one terminal ok envelope
whose ` + "`data`" + ` embeds the same document — the object itself for
--emit ` + emitJSON + `, and for --emit ` + emitDOT + ` an object whose single required
string member ` + "`dot`" + ` carries the DOT text, so the unwrap is
` + "`jq -r .data.dot`" + `.

The ` + "`reach`" + ` block carries a required ` + "`abstraction`" + ` marker with the token
value ` + abstractionMarkerValue + `. THE SOUNDNESS RULE: the relation is the
DECLARED over-approximation, not the runtime one — nodes are merged and
guard/observed atoms are unpruned — so universal claims ("no path does
X") proved over it hold at runtime, while existence claims ("some path
reaches X") may be spurious.

This command is NOT an acceptance gate. The document carries no verdict
and no finding field; a model that loads exports, whatever lint would
say about it. Run ` + "`intrastate lint`" + ` to accept or refuse a model.

When the traversal does not complete under this build's node ceiling of
` + itoa(graphlint.NodeCeiling()) + `, the command refuses with ` + codeExportTooLarge + ` rather than
emitting a partial document: a partial graph diffs as a graph change.
Narrow a tag's declared domain and re-run.

For one model and one build, emission is byte-for-byte identical across
invocations in every --emit and --as combination. The promise is scoped
to (model, build) and not across builds.
`)
	return string(b)
}

// runGraphExport is the verb's RunE.
//
// The arm ORDER is normative (REQ-11): the mode is validated first, then
// the `--model`/`--flow` SELECTION arms, then `--emit`'s value, and only
// then any file I/O. So `graph --model <unreadable> --emit=xml` refuses
// `flag-invalid-value` naming `emit`, never `model-unreadable` — the
// request is wrong independent of the environment, and reporting the
// environment would send the caller to fix the wrong thing.
func runGraphExport(cmd *cobra.Command, _ []string) error {
	if ce := respond.ValidateMode(cmd); ce != nil {
		return respond.Fail(cmd, ce)
	}

	// The selection arms, mirroring lint's arm set and code spellings
	// VERBATIM (`internal/cli/lint.go::runLint`). C1 mirrors, never
	// redefines: lint owns these codes.
	modelPath, _ := cmd.Flags().GetString("model")
	flowID, _ := cmd.Flags().GetString("flow")
	if modelPath != "" && flowID != "" {
		return respond.Fail(cmd, &clierr.CLIError{
			Code:    "flag-mutually-exclusive",
			Message: "--flow and --model are mutually exclusive",
			Group:   clierr.GroupUserEnv,
		})
	}
	if modelPath == "" && flowID != "" {
		return respond.Fail(cmd, &clierr.CLIError{
			Code:    "flag-invalid-value",
			Param:   "flow",
			Message: "no model is registered for the flow id " + flowID,
			Group:   clierr.GroupUserEnv,
			Hint:    "give --model <path> instead",
		})
	}
	if modelPath == "" {
		return respond.Fail(cmd, &clierr.CLIError{
			Code:    "flag-required",
			Param:   "model",
			Message: "--model <path> names the transition model to export",
			Group:   clierr.GroupUserEnv,
		})
	}

	// `--emit` is the ONE value-checked flag on this verb, and it is
	// checked BEFORE any file I/O.
	emit, _ := cmd.Flags().GetString(emitFlagName)
	if !slices.Contains(emitFormats, emit) {
		return respond.Fail(cmd, &clierr.CLIError{
			Code:    "flag-invalid-value",
			Param:   emitFlagName,
			Message: "unknown document format " + emit,
			Group:   clierr.GroupUserEnv,
			Hint:    "use --" + emitFlagName + " " + emitJSON + " or --" + emitFlagName + " " + emitDOT,
		})
	}

	src, err := os.ReadFile(modelPath)
	if err != nil {
		return respond.Fail(cmd, &clierr.CLIError{
			Code:    "model-unreadable",
			Param:   "model",
			Message: "cannot read " + modelPath,
			Detail:  err.Error(),
			Group:   clierr.GroupUserEnv,
			Cause:   err,
		})
	}
	// The load path is the shared one lint uses, so the load-category
	// findings are identical. The ADVISORIES it also returns are lint's
	// advisory channel and never graph data, so they are dropped here
	// rather than becoming document content (REQ-105).
	m, _, err := table.LoadWithAdvisories(src, modelPath)
	if err != nil {
		return respond.Fail(cmd, &clierr.CLIError{
			Code:     "model-invalid",
			Param:    "model",
			Message:  "model does not conform to the transition-model schema",
			Detail:   err.Error(),
			Group:    clierr.GroupUserEnv,
			Findings: loadFindings(modelPath, err),
			Cause:    err,
		})
	}

	// The traversal, reached through the same `graphlint` surface lint
	// uses: ONE traversal, two call sites, so verb/lint drift is structural
	// rather than disciplinary (`0021:C4`).
	nodes, edges, complete := graphlint.ReachGraph(m)
	if !complete {
		// NEVER a partial document. The refusal names the ceiling and the
		// narrow-a-domain remedy — the same remedy lint gives — because a
		// partial graph would diff as a graph change.
		return respond.Fail(cmd, &clierr.CLIError{
			Code: codeExportTooLarge,
			Message: "the reachable node set exceeds this build's node ceiling of " +
				itoa(graphlint.NodeCeiling()),
			Detail: "the traversal did not complete, so no document is emitted: " +
				"a partial graph diffs as a graph change",
			Group: clierr.GroupUserEnv,
			Hint:  "narrow a tag's declared domain and re-run",
		})
	}

	doc := buildGraphDocument(m, nodes, edges)

	payload := graphPayload{}
	switch emit {
	case emitDOT:
		dot := renderDOT(doc)
		payload.document = dot
		payload.data = dotPayload{Dot: dot}
	default:
		rendered, merr := marshalGraphJSON(doc)
		if merr != nil {
			// Marshalling a document of strings, slices and maps cannot
			// fail; the branch exists so the error is never silently
			// dropped.
			return respond.Fail(cmd, &clierr.CLIError{
				Code:    "internal-error",
				Message: "the exported document could not be encoded",
				Detail:  merr.Error(),
				Group:   clierr.GroupInternal,
				Cause:   merr,
			})
		}
		payload.document = string(rendered)
		payload.data = doc
	}

	return respond.OK(cmd, respond.Success{Data: payload})
}
