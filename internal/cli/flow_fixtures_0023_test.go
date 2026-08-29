package cli

// RDR 0023 — shared fixtures and helpers for the resolve-envelope
// projection suite.
//
// Three design rules hold throughout, and each is load-bearing for whether
// the suite actually discriminates `0023:C1` / `0023:C2` rather than
// restating the implementation:
//
//  1. NOTHING here reads the projection's own code to learn what the
//     projection does. C2 forbids the tautological oracle by name — "an
//     oracle that derives the ECHO set by observing what the projection
//     drops is tautological — it restates the implementation and cannot
//     fail" — so every expectation below is an AUTHORED literal or the
//     SAME REQUEST'S DEFAULT-MODE output, never the projected run's own
//     shape.
//
//  2. Every symbol the implementation still owes (the declared assignment
//     table, the reader-execution seam, the total walker) is reached
//     REFLECTIVELY or BEHAVIOURALLY, never by a direct Go reference. A
//     direct reference would fail to COMPILE the package, taking all 128
//     shipped test files red with it — which would destroy the very
//     "no predecessor oracle moves" property `0023:A3` claims and make the
//     red gate unreadable. Each test below fails on its own, for its own
//     clause.
//
//  3. The measurand is fixed where the contract fixes it. C1 measures
//     STRICTLY SHORTER on the FULL EMITTED LINE — "the complete
//     {"type":"ok","data":{…}} NDJSON record as written, excluding the
//     trailing newline — not on the .data payload alone". `emittedLine`
//     below is that unit and nothing in this suite measures `.data`.

import (
	"encoding/json"
	"reflect"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// --- the normative partition ---------------------------------------------

// planOnlyFlag is the normative flag spelling (`0023:C1`, `0023:D-naming`).
// It is a const rather than a literal at each site so a rename shows up as
// one edit rather than as thirty silently-diverging strings.
const planOnlyFlag = "plan-only"

// echoGroup0023 is `0023:C2`'s ECHO group, by WIRE key: "the ECHO group is
// the model reference, the observed tags (--tag, echoed unchanged), the
// assembled owned view and the invoked reader identities, and the requested
// outcome". The census (`0023:CEN`) fixes the five names.
//
// AUTHORED, never derived: this list is what makes the S2/S3 oracles able
// to fail. A list computed from the projection would agree with any
// implementation, correct or not.
var echoGroup0023 = []string{
	"model", "observed", "owned", "readers", "outcome",
}

// planGroup0023 is `0023:C2`'s PLAN group by wire key, and `0023:C1`'s
// "exactly the plan group — rule, gates, emit, next, writes, clear,
// escaped, escape_class, revision".
//
// `escape_class` is here because it is a PLAN-group member and always-keep
// core; whether it is EMITTED on a given run is its producer's presence
// rule (`0005:A-3`), which REQ-24/REQ-25 require asserting against the same
// run's default output rather than unconditionally.
var planGroup0023 = []string{
	"revision", "rule", "gates", "emit", "next", "writes", "clear",
	"escaped", "escape_class",
}

// projectedKeyLiteral0023 is `0023:S2`'s explicit literal list, in
// `0023:A-1`'s order — "the projected order is the default order minus the
// deleted keys, never a re-sort", which puts `revision` FIRST.
//
// It carries EIGHT names; `escape_class` is the conditional ninth
// (`0023:A-2`) and is appended by the oracle only when the same request's
// DEFAULT output carried it.
var projectedKeyLiteral0023 = []string{
	"revision", "rule", "gates", "emit", "next", "writes", "clear", "escaped",
}

// alwaysKeepCore0023 is this verb's always-keep core (`0023:C2`,
// JDR 0002 §D1): "rule, escaped, escape_class, revision".
var alwaysKeepCore0023 = []string{"rule", "escaped", "escape_class", "revision"}

// --- the emitted-line measurand ------------------------------------------

// emittedLine returns C1's measurand: "the complete
// {"type":"ok","data":{…}} NDJSON record as written, excluding the trailing
// newline". Not the `.data` payload — the envelope is the unit the caller
// actually pays for, and it is the stricter reading.
//
// It is FATAL on a multi-record stdout: `0005:C1` fixes one terminal
// envelope, and a width measured over two records would measure the wrong
// thing while still producing a number.
func emittedLine(t *testing.T, stdout string) string {
	t.Helper()

	trimmed := strings.TrimRight(stdout, "\n")
	if trimmed == "" {
		t.Fatal("stdout carried no terminal record; `0005:C1` fixes exactly " +
			"one, and a width measured over none is not a measurement")
	}
	if strings.Contains(trimmed, "\n") {
		t.Fatalf("stdout carries MORE THAN ONE line under --as=json; the "+
			"width measurand is the single terminal record\n%s", stdout)
	}
	return trimmed
}

// textLines returns the emitted text-mode lines, trailing blank dropped.
// `0023:S5`'s subset is over these, per line and byte-identically.
func textLines(stdout string) []string {
	trimmed := strings.TrimRight(stdout, "\n")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "\n")
}

// textWidth is `0023:C1`'s text measurand: "the total rendered byte count
// of the emitted lines". Measured on the stream as written, so a projection
// that drops no bytes cannot pass by a line-count technicality.
func textWidth(stdout string) int { return len(strings.TrimRight(stdout, "\n")) }

// --- payload key access, order-preserving --------------------------------

// payloadKeyOrder returns the `data` object's TOP-LEVEL keys in WIRE order.
//
// It decodes with a streaming token reader rather than into a map, because
// `0023:C1` requires "the carried keys keep their default-mode relative
// order" and a Go map has no order to assert. Nested objects are skipped:
// the clause is about the TOP-LEVEL key set (`0023:REQ-22`).
func payloadKeyOrder(t *testing.T, stdout string) []string {
	t.Helper()

	var env struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &env); err != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", err, stdout)
	}
	if env.Type != "ok" {
		t.Fatalf("envelope type = %q; want %q — the projection rides the "+
			"SUCCESS path only (`0023:C1`)", env.Type, "ok")
	}

	dec := json.NewDecoder(strings.NewReader(string(env.Data)))
	tok, err := dec.Token()
	if err != nil {
		t.Fatalf("read `data`: %v", err)
	}
	if delim, ok := tok.(json.Delim); !ok || delim != '{' {
		t.Fatalf("`data` is not a JSON object: %s", env.Data)
	}

	var keys []string
	depth := 0
	for dec.More() || depth > 0 {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch v := tok.(type) {
		case json.Delim:
			switch v {
			case '{', '[':
				depth++
			case '}', ']':
				depth--
			}
		case string:
			if depth == 0 {
				keys = append(keys, v)
				// Consume the value; only its DELIMITERS matter for depth.
				if !dec.More() {
					continue
				}
			}
		}
	}
	return keys
}

// rawFieldsOf returns the `data` object's top-level members as UNDECODED
// JSON bytes, keyed by wire key.
//
// Raw is the point: `0023:C1` requires "every carried field's rendering
// MUST be byte-identical to its default-mode rendering (same encoder, HTML
// escaping disabled)", and a decode-then-compare would launder exactly the
// encoder differences that clause exists to catch — `<` re-escaped to
// `<` decodes equal and renders differently.
func rawFieldsOf(t *testing.T, stdout string) map[string]json.RawMessage {
	t.Helper()

	var env struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &env); err != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", err, stdout)
	}
	if env.Type != "ok" {
		t.Fatalf("envelope type = %q; want %q", env.Type, "ok")
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(env.Data, &fields); err != nil {
		t.Fatalf("`data` is not a JSON object: %v\n%s", err, string(env.Data))
	}
	return fields
}

// sortedKeys returns a raw-field map's keys, sorted, for set comparisons
// and for readable failure messages.
func sortedKeys(m map[string]json.RawMessage) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// --- invocation shapes ---------------------------------------------------

// resolveArgs builds a `flow resolve` argv. Extra args ride LAST so a
// caller appends `--plan-only` / `--as=text` without rebuilding the base.
func resolveArgs(model, bind, outcome string, extra ...string) []string {
	args := []string{"flow", "resolve", "--model", model, "--outcome", outcome}
	if bind != "" {
		args = append(args, "--artifact", bind)
	}
	return append(args, extra...)
}

// pricingModelPath returns the CHECKED-IN pricing decision table
// (`models/examples/pricing-decision-table.toml`), the normative fixture's
// model: `0023:S1` / `0023:S2` pin "the projected wire record of the
// pricing 2×2 call", and `0023:MVV` runs "over the checked-in
// decision-table fixture".
//
// Checked-in matters for `0023:REQ-104`: the 40% MVV bar is stated on the
// three checked-in fixtures precisely because a spike's scratch model is
// not reproducible from the repo.
func pricingModelPath(t *testing.T) string {
	t.Helper()

	return repoRootFor(t) + "/models/examples/pricing-decision-table.toml"
}

// releaseModelPath returns the checked-in release grammar model — the
// gate/write-heavy shape `0023:MVV` step 6 names beside the pricing table.
func releaseModelPath(t *testing.T) string {
	t.Helper()

	return repoRootFor(t) + "/models/examples/release-grammar.toml"
}

// reviewModelPath returns the checked-in review state machine, the third
// fixture `0023:REQ-104`'s 40% bar names.
func reviewModelPath(t *testing.T) string {
	t.Helper()

	return repoRootFor(t) + "/models/examples/review-state-machine.toml"
}

// pricingCall is the NORMATIVE FIXTURE invocation: the pricing 2×2 call
// resolving `paid-eu` (`0023:S1`, `0023:DT` step 2's witness
// `{"revision":"","rule":"paid-eu",…}`).
//
// The model declares no owned tags and no accessors, so no artifact is
// bound — which is also why it is the shape whose echo containers are
// EMPTY, the arm `0023:A9`'s Phase-1 oracle needs.
func pricingCall(t *testing.T, extra ...string) []string {
	t.Helper()

	return resolveArgs(pricingModelPath(t), "", "decide",
		append([]string{"--tag", "tier=paid", "--tag", "region=eu"}, extra...)...)
}

// --- the flag surface, probed without a compile dependency ---------------

// planOnlyRegistered reports whether cmd's OWN flag set registers
// `plan-only`. Exact-name `Lookup`, the same probe shape the `--all`
// oracles use (`0023:A5`: "the `--all` probes are exact-name
// `Lookup("all")`").
//
// Persistent flags are checked too: a persistent registration on a parent
// is exactly the silent surface-widening C1's non-registration clause
// fences.
func planOnlyRegistered(cmd *cobra.Command) bool {
	return cmd.Flags().Lookup(planOnlyFlag) != nil ||
		cmd.PersistentFlags().Lookup(planOnlyFlag) != nil
}

// resolveCmd returns the `flow resolve` command off a freshly built root,
// or nil. Tests use it to assert the flag's registration SHAPE without
// reaching into the constructor.
func resolveCmd(t *testing.T) *cobra.Command {
	t.Helper()

	for _, c := range NewRootCmd().Commands() {
		if c.Name() != "flow" {
			continue
		}
		for _, sub := range c.Commands() {
			if sub.Name() == "resolve" {
				return sub
			}
		}
	}
	return nil
}

// commandPath renders a command's full path from the root, so a failure
// message names `intrastate flow resolve` rather than the ambiguous
// `resolve`.
func commandPath(cmd *cobra.Command) string {
	var parts []string
	for c := cmd; c != nil; c = c.Parent() {
		parts = append(parts, c.Name())
	}
	slices.Reverse(parts)
	return strings.Join(parts, " ")
}

// --- reflective access to the payload struct -----------------------------

// resolvePayloadType is the reflected `resolvePayload` struct type.
//
// The reflection surface is deliberately the one `0023:A3` audited —
// `NumField`, `FieldByName`, and the JSON tag reachable from a field — so
// the partition oracle can be written over field NAMES (`0023:A-4`) without
// depending on any type the conversion changes.
func resolvePayloadType() reflect.Type { return reflect.TypeOf(resolvePayload{}) }

// wireKeyForField returns the JSON wire key a struct field renders under,
// stripping `,omitempty` and friends. `0023:A-10` keeps the key STRINGS
// unchanged under the pointer conversion, which is what keeps the three
// shipped wire oracles green — so this mapping is stable across the change.
func wireKeyForField(f reflect.StructField) string {
	tag := f.Tag.Get("json")
	if tag == "" {
		return f.Name
	}
	if i := strings.IndexByte(tag, ','); i >= 0 {
		tag = tag[:i]
	}
	if tag == "" || tag == "-" {
		return f.Name
	}
	return tag
}

// payloadWireKeys returns every wire key `resolvePayload` declares, in
// DECLARATION order — which is the JSON key order Go's encoder emits and
// the order `0023:C1`'s "carried keys keep their default-mode relative
// order" is stated against.
func payloadWireKeys() []string {
	rt := resolvePayloadType()
	out := make([]string, 0, rt.NumField())
	for i := range rt.NumField() {
		out = append(out, wireKeyForField(rt.Field(i)))
	}
	return out
}
