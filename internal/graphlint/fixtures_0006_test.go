package graphlint_test

// Shared fixture helpers for the RDR 0006 graph-lint suite.
//
// Nothing here mocks the unit under test. Models are built by the real
// `internal/table` loader from authored TOML, and every lint assertion
// drives the real `internal/graphlint`. The loader is a collaborator RDR
// 0002 owns; this RDR owns what the graph MEANS.

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
	"github.com/newcoinc/intrastate/internal/graphlint"
	"github.com/newcoinc/intrastate/internal/table"
)

// --- model construction --------------------------------------------------

// modelHeader is the minimum JDR 0001 §D7 preamble every fixture needs: a
// model id, the recognized-outcome alphabet, and the reserved
// `[tags.recognized]` declaration.
const modelHeader = `
[model]
id = "t"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true
`

// source assembles a complete fixture document: the header, the caller's
// tag declarations, the derived accessor bindings RDR 0002 requires of
// every owned tag, and the caller's body (contexts, initial, terminal,
// rules).
//
// The accessor bindings are DERIVED rather than authored per fixture so an
// accessor-arity slip can never masquerade as an RDR 0006 failure: every
// fixture below either loads clean or fails on a clause this RDR owns.
func source(decls, body string) string {
	return rootKeys(body) + modelHeader + decls + accessorsFor(decls) + tables(body)
}

// rootKeys lifts a body's ROOT scalar keys (`outcomes`, `terminal`) above
// every table, because TOML binds a bare key to the most recent table
// header: `terminal = [...]` written after `[write.own]` is a key INSIDE
// that accessor, and strict decoding refuses it. Splitting here keeps the
// fixture bodies readable as one authored document.
func rootKeys(body string) string {
	var b strings.Builder
	b.WriteString("\noutcomes = [\"go\", \"stop\"]\n")
	for line := range strings.Lines(body) {
		if strings.HasPrefix(strings.TrimSpace(line), "terminal =") {
			b.WriteString(line)
		}
	}
	return b.String()
}

// tables returns the body with its root scalar keys removed.
func tables(body string) string {
	var b strings.Builder
	for line := range strings.Lines(body) {
		if strings.HasPrefix(strings.TrimSpace(line), "terminal =") {
			continue
		}
		b.WriteString(line)
	}
	return b.String()
}

// accessorsFor emits the reader/writer accessors covering every owned tag
// the declarations name.
func accessorsFor(decls string) string {
	var owned []string
	var key string
	for line := range strings.Lines(decls) {
		line = strings.TrimSpace(line)
		if name, ok := strings.CutPrefix(line, "[tags."); ok {
			// A tag name TOML must quote — one carrying `.`, `;`, or any
			// other bare-key-illegal rune — arrives here still wearing its
			// quotes. The accessor `keys` list is re-quoted below, so the
			// authored quotes must come off first or the accessor binds a
			// tag literally named `"a;b"` and the loader refuses the
			// fixture for an undeclared tag.
			key = strings.TrimSuffix(name, "]")
			if unquoted, err := strconv.Unquote(key); err == nil {
				key = unquoted
			}
			continue
		}
		if line == `provenance = "owned"` && key != "" {
			owned = append(owned, key)
			key = ""
		}
	}
	if len(owned) == 0 {
		return ""
	}
	slices.Sort(owned)
	quoted := make([]string, 0, len(owned))
	for _, k := range owned {
		quoted = append(quoted, strconv.Quote(k))
	}
	keys := "[" + strings.Join(quoted, ", ") + "]"
	return `
[read.own]
role = "t"
path = "t.own"
keys = ` + keys + `
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ` + keys + `
timeout = "2s"
read_back = true
`
}

// mustLoad loads authored TOML that the fixture says must load clean. A
// refusal here is a fixture bug, never the defect under test: RDR 0002
// refuses non-conforming source BEFORE normalization and lint never runs
// (the disposition table's upstream row).
func mustLoad(t *testing.T, src string) *table.Model {
	t.Helper()

	m, err := table.Load([]byte(src), "fixture.toml")
	if err != nil {
		t.Fatalf("fixture must load clean; refused: %v\n--- source ---\n%s", err, src)
	}
	if m == nil {
		t.Fatal("loader returned a nil model with no error")
	}
	return m
}

// lintSource loads src and runs the real engine over it, through the one
// request builder and the one entry point.
func lintSource(t *testing.T, src string) graphlint.Report {
	t.Helper()
	return graphlint.Run(graphlint.NewRequest(mustLoad(t, src)))
}

// lint is the common shorthand: assemble, load, run.
func lint(t *testing.T, decls, body string) graphlint.Report {
	t.Helper()
	return lintSource(t, source(decls, body))
}

// --- report inspection ---------------------------------------------------

// withCode returns every finding in the report carrying code.
func withCode(r graphlint.Report, code string) []clierr.Finding {
	var out []clierr.Finding
	for _, f := range r.Findings {
		if f.Code == code {
			out = append(out, f)
		}
	}
	return out
}

// countCode counts findings carrying code.
func countCode(r graphlint.Report, code string) int {
	return len(withCode(r, code))
}

// hasCode reports whether the report carries at least one finding of code.
func hasCode(r graphlint.Report, code string) bool {
	return countCode(r, code) > 0
}

// codesIn returns the sorted, duplicate-free code set the report carries.
func codesIn(r graphlint.Report) []string {
	var out []string
	for _, f := range r.Findings {
		if !slices.Contains(out, f.Code) {
			out = append(out, f.Code)
		}
	}
	slices.Sort(out)
	return out
}

// codesOf returns the report's codes in sorted order, PRESERVING
// duplicates. Distinct from codesIn, which dedupes: an oracle that must
// bound a report (SC-4's "exact finding lists") has to see a duplicated
// code, since a regression emitting the same code twice is a real defect
// that a deduped comparison hides.
func codesOf(r graphlint.Report) []string {
	out := make([]string, 0, len(r.Findings))
	for _, f := range r.Findings {
		out = append(out, f.Code)
	}
	slices.Sort(out)
	return out
}

// blockingCodes returns the sorted, duplicate-free code set of the
// report's BLOCKING findings.
func blockingCodes(r graphlint.Report) []string {
	var out []string
	for _, f := range r.Blocking() {
		if !slices.Contains(out, f.Code) {
			out = append(out, f.Code)
		}
	}
	slices.Sort(out)
	return out
}

// render formats a report for a failure message, so an assertion always
// shows what lint DID emit.
func render(r graphlint.Report) string {
	if len(r.Findings) == 0 {
		return " (no findings)"
	}
	var b strings.Builder
	for _, f := range r.Findings {
		fmt.Fprintf(&b, "\n  %s sev=%s model=%s rule=%q reason=%q dim=%q "+
			"atom=(%s %s %s %s) class=%q msg=%q",
			f.Code, f.Severity, f.Model, f.Rule, f.Reason, f.Dimension,
			f.Key, f.Operator, f.Literal, f.Block, f.Class, f.Message)
	}
	return b.String()
}

// requireCode asserts the report carries at least one finding of code and
// returns them.
func requireCode(t *testing.T, r graphlint.Report, code string) []clierr.Finding {
	t.Helper()

	got := withCode(r, code)
	if len(got) == 0 {
		t.Fatalf("no %s finding; report:%s", code, render(r))
	}
	return got
}

// requireOneCode asserts the report carries EXACTLY one finding of code.
func requireOneCode(t *testing.T, r graphlint.Report, code string) clierr.Finding {
	t.Helper()

	got := withCode(r, code)
	if len(got) != 1 {
		t.Fatalf("%d %s findings; want exactly 1; report:%s", len(got), code, render(r))
	}
	return got[0]
}

// requireNoCode asserts the report carries no finding of code.
func requireNoCode(t *testing.T, r graphlint.Report, code string) {
	t.Helper()

	if n := countCode(r, code); n != 0 {
		t.Fatalf("%d %s findings; want none; report:%s", n, code, render(r))
	}
}

// requireClean asserts the report carries no BLOCKING finding — the
// success disposition.
func requireClean(t *testing.T, r graphlint.Report) {
	t.Helper()

	if n := len(r.Blocking()); n != 0 {
		t.Fatalf("%d blocking findings; want a clean model; report:%s", n, render(r))
	}
}

// namesRule reports whether some finding of code names ruleID.
func namesRule(r graphlint.Report, code, ruleID string) bool {
	for _, f := range withCode(r, code) {
		if f.Rule == ruleID {
			return true
		}
	}
	return false
}

// --- shared fixture bodies -----------------------------------------------

// twoStateDecls declares one owned enum `status` over {a, b} plus one
// optional owned enum used as a guard dimension. It is the workhorse
// declaration block for reachability and coverage fixtures.
const twoStateDecls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[tags.flag]
provenance = "owned"
kind = "bool"
single_valued = true
required = true
`

// legalBody is a complete, minimal legal model: a root, a terminal, and
// one group whose two rows partition the `flag` dimension.
const legalBody = `
terminal = ["done"]

[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"

[[rule]]
id = "advance-off"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "false"
[rule.write]
status = "b"
`

// danglingBody declares NO initial owned state. RDR 0002's loader refuses
// every SYNTACTIC dangling reference before normalization (an undeclared
// tag, an out-of-alphabet outcome, an unknown gate accessor), so the arm
// of invariant 1 that reaches lint is the model half of the clause: "the
// model must declare an initial owned state".
const danglingBody = `
terminal = ["done"]

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "dangles"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`

// escapeBody carries an escape row alongside an ordinary one, so the
// input contract's row-kind and failure-class fields are exercised.
const escapeBody = `
terminal = ["done"]

[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"

[[rule]]
id = "rescue"
escape = ["no_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
`

// --- declaration blocks used across the invariant suites -----------------

// statusOnlyDecls is the smallest useful owned lattice: one required
// single-valued enum over {a, b}.
const statusOnlyDecls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true
`

// optionalGuardDecls adds an OPTIONAL owned enum alongside `status`. An
// unmarked key defaults to optional (RDR 0003), so `opt` carries the
// presence dimension and makes its readers able to refuse
// `guard_unevaluable`.
const optionalGuardDecls = statusOnlyDecls + `
[tags.opt]
provenance = "owned"
kind = "enum"
domain = ["p", "q"]
single_valued = true
`

// alwaysPresentOwnedDecls adds an owned key carrying RDR 0003's EXPLICIT
// always-present marker.
const alwaysPresentOwnedDecls = statusOnlyDecls + `
[tags.always]
provenance = "owned"
kind = "enum"
domain = ["p", "q"]
single_valued = true
required = true
`

// observedRequiredDecls adds an OBSERVED key carrying the always-present
// marker — the SC-17 control that must take no finding.
const observedRequiredDecls = statusOnlyDecls + `
[tags.seen]
provenance = "observed"
kind = "enum"
domain = ["y", "n"]
single_valued = true
required = true
`

// nonFiniteDecls adds an owned `scalar` — a kind that declares no finite
// domain, so any dimension over it is not provable.
const nonFiniteDecls = statusOnlyDecls + `
[tags.free]
provenance = "owned"
kind = "scalar"
required = true
`

// wideIntDecls adds an owned int whose declared bound makes the product
// large enough to exceed the published bound.
const wideIntDecls = statusOnlyDecls + `
[tags.n]
provenance = "owned"
kind = "int"
min = 0
max = 1000000
single_valued = true
required = true
`

// --- body fragments ------------------------------------------------------

// rootAndTerminal is the standard root plus a `status = b` stop set.
const rootAndTerminal = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"
`

// --- the multi-defect fixture (REQ-MVV, SC-13) ---------------------------
//
// One model carrying several INDEPENDENT decidable defects across two
// groups, so a first-failure engine emitting only the first cannot pass.
// It pairs with multiDefectBodyReversed, which authors the same rules in
// the opposite order: the emitted finding set must be identical.
//
// The defects: an overlapping ordinary pair (`over-one` x `over-two`), two
// rows that can refuse `guard_unevaluable` over the optional `opt`
// (`refuse-one`, `refuse-two`), and a reachable node with no outgoing row
// and no terminal covering it.

const multiDefectBody = `
terminal = ["done"]

[initial]
status = "a"
opt = "p"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "over-one"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "over-two"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "refuse-one"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"
[rule.guard.all.opt]
eq = "p"
[rule.write]
status = "b"

[[rule]]
id = "refuse-two"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"
[rule.guard.all.opt]
eq = "q"
[rule.write]
status = "b"
`

// multiDefectBodyReversed authors the same four rules in the opposite
// order. The emitted finding set must not depend on it.
const multiDefectBodyReversed = `
terminal = ["done"]

[initial]
status = "a"
opt = "p"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "refuse-two"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"
[rule.guard.all.opt]
eq = "q"
[rule.write]
status = "b"

[[rule]]
id = "refuse-one"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"
[rule.guard.all.opt]
eq = "p"
[rule.write]
status = "b"

[[rule]]
id = "over-two"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "over-one"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`

// advisoryBody is a clean model carrying advisory findings only: a group
// closed by a bare escape row, a redundant row (a proper subset of a
// sibling's assignments), and a rule no reachable node satisfies.
const advisoryBody = `
terminal = ["done"]

[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "wide"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "narrow"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"

[[rule]]
id = "only-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"

[[rule]]
id = "bare-rescue"
escape = ["no_match", "ambiguous_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"
`

// legalBodyReordered authors legalBody's `advance-on` rule with its match
// and guard blocks in the opposite order. Normalization sorts the atom
// set, so the two must be indistinguishable downstream.
const legalBodyReordered = `
terminal = ["done"]

[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance-on"
[rule.guard.all.flag]
eq = "true"
[rule.match.recognized]
eq = "go"
[rule.match.status]
eq = "a"
[rule.write]
status = "b"

[[rule]]
id = "advance-off"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "false"
[rule.write]
status = "b"
`
