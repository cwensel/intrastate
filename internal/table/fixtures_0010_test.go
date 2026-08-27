package table_test

// RDR 0010 — shared fixture corpus for the stateless-decision-table suite.
//
// Every model here is authored as a complete TOML document and handed to
// the REAL `table.Load`. Nothing mocks the loader, the normalizer, or the
// dump renderer.
//
// The corpus is shaped by the MVV (`0010:MVV`) and by SC-1/SC-2/SC-3:
//
//   - two observed enum dimensions of two values each, each `single_valued`
//     and `required`, DISCRIMINATED BY `[rule.guard.all.<key>]` atoms — the
//     authoring that makes them participate in the scoped product
//     (`0006:C7`, `0003:C13`);
//   - one recognized outcome bound by `[rule.match.recognized]`;
//   - no `[initial]`, no accessors, no owned tags;
//   - `[rule.emit]` on each ordinary rule.
//
// ASSUMPTION-9: this fixture is a SUPERSET of the MVV's three-rule table —
// SC-5 additionally requires one selectable row authoring NO emit block, so
// the complete four-rule table below leaves the fourth rule's emit absent.

import (
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// dtHeader is the decision-table preamble: the class declaration, the
// recognized-outcome alphabet, and the two observed enum dimensions.
//
// No `[initial]`, no `[read.*]`/`[write.*]`/`[gate.*]`, no owned tag.
const dtHeader = `outcomes = ["decide"]

[model]
id = "dt"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.a]
provenance = "observed"
kind = "enum"
domain = ["x", "y"]
single_valued = true
required = true

[tags.b]
provenance = "observed"
kind = "enum"
domain = ["p", "q"]
single_valued = true
required = true
`

// dtRule renders one ordinary decision-table rule: the outcome binding, the
// two guard atoms that discriminate the cell, and an optional emit block.
//
// The rule carries NO write block — the shape C2 conditions `0002:C4` on
// the state-machine class to admit.
func dtRule(id, a, b string, emit map[string]string) string {
	var sb strings.Builder
	sb.WriteString("\n[[rule]]\nid = \"" + id + "\"\n")
	sb.WriteString("[rule.match.recognized]\neq = \"decide\"\n")
	sb.WriteString("[rule.guard.all.a]\neq = \"" + a + "\"\n")
	sb.WriteString("[rule.guard.all.b]\neq = \"" + b + "\"\n")
	if emit != nil {
		sb.WriteString("[rule.emit]\n")
		for _, k := range sortedKeys(emit) {
			sb.WriteString(k + " = \"" + emit[k] + "\"\n")
		}
	}
	return sb.String()
}

// sortedKeys renders a map's keys in byte order so a fixture's authored
// text is stable across runs. It deliberately does NOT sort the way the
// implementation must — the emit ORDER assertion is taken on the loaded
// row, from an AUTHORED order that is not already sorted (see dtUnordered).
func sortedKeys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// dtPartial is the MVV step-1/step-2 table: three of the four (a,b) cells
// covered, each ordinary rule carrying `[rule.emit]`.
var dtPartial = dtHeader +
	dtRule("cell-xp", "x", "p", map[string]string{"verdict": "alpha"}) +
	dtRule("cell-xq", "x", "q", map[string]string{"verdict": "beta"}) +
	dtRule("cell-yp", "y", "p", map[string]string{"verdict": "gamma"})

// dtComplete is the MVV step-3 table: the fourth ordinary rule closes the
// product. Per SC-5 that fourth rule authors NO `[rule.emit]`, so the
// `{}`-never-`null` assertion has a selectable row to take.
var dtComplete = dtPartial + dtRule("cell-yq", "y", "q", nil)

// dtCompleteEmitting is dtComplete with the fourth rule AUTHORING an emit
// block — the MVV step-4 selection whose `data.emit` equals its authored
// block.
var dtCompleteEmitting = dtPartial +
	dtRule("cell-yq", "y", "q", map[string]string{"verdict": "delta", "code": "4"})

// dtEscapeVariant is the MVV step-3 variant: the fourth rule is replaced by
// an escape row rescuing `no_match`, itself carrying `[rule.emit]`.
var dtEscapeVariant = dtPartial + `
[[rule]]
id = "otherwise"
escape = ["no_match"]
[rule.match.recognized]
eq = "decide"
[rule.emit]
verdict = "fallback"
`

// dtClassOmitted is the MVV step-6 negative control: the decision table
// with `class` omitted. It must LOAD — the agreement check is
// one-directional (C1) — and lint as a rootless state machine
// (`graph-dangling-edge` at `element = model`, `0006:C18` unchanged).
//
// DEVIATION (see deviations.md D1): the rules carry a write block and the
// model declares one owned tag with its accessors. Without them the
// document refuses at LOAD with `malformed rule shape` — `0002:C4`'s
// write-block arm binds the state-machine class, which is exactly what
// REQ-16 conditions — so a byte-for-byte `class`-strip of `dtComplete`
// could not exhibit the lint finding MVV step 6 names. What the control
// pins is the CLASS-KEYED behaviour: zero `[initial]` under the
// state-machine class still traverses nothing and takes `0006:C18`'s
// missing-root finding (REQ-53, REQ-83).
const dtClassOmitted = `outcomes = ["decide"]

[model]
id = "dt"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.a]
provenance = "observed"
kind = "enum"
domain = ["x", "y"]
single_valued = true
required = true

[tags.b]
provenance = "observed"
kind = "enum"
domain = ["p", "q"]
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["seen"]
single_valued = true
required = true

[read.state]
role = "state"
path = "s.own"
keys = ["status"]
timeout = "2s"

[write.state]
role = "state"
path = "s.own"
keys = ["status"]
timeout = "2s"
read_back = true

[[rule]]
id = "cell-xp"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "x"
[rule.guard.all.b]
eq = "p"
[rule.write]
status = "seen"
[rule.emit]
verdict = "alpha"
`

// dtWithOwnedTag is the MVV step-6 negative control: a declared decision
// table that also declares one `provenance = "owned"` tag. It must refuse
// at load with `malformed model declaration`.
var dtWithOwnedTag = strings.Replace(dtComplete, "[tags.a]", `[tags.own]
provenance = "owned"
kind = "scalar"

[tags.a]`, 1)

// dtUnordered authors one rule's emit keys in NON-sorted order, so the
// key-sorted normalization (C3) is asserted against an authored order that
// differs from the normalized one.
var dtUnordered = dtHeader + `
[[rule]]
id = "cell-xp"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "x"
[rule.guard.all.b]
eq = "p"
[rule.emit]
zeta = "last"
alpha = "first"
mu = "middle"
` + dtRule("cell-xq", "x", "q", nil) +
	dtRule("cell-yp", "y", "p", nil) +
	dtRule("cell-yq", "y", "q", nil)

// dtEmptyEmit authors a PRESENT BUT EMPTY `[rule.emit]` — the deliberate
// divergence from the write/clear/gate "EVEN AN EMPTY ONE" rule (C3).
var dtEmptyEmit = dtHeader + `
[[rule]]
id = "cell-xp"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "x"
[rule.guard.all.b]
eq = "p"
[rule.emit]
` + dtRule("cell-xq", "x", "q", nil) +
	dtRule("cell-yp", "y", "p", nil) +
	dtRule("cell-yq", "y", "q", nil)

// dtExpanding authors an emit block on an EXPANDING rule: a multi-member
// `in` MATCH atom mints several rows from one rule (SC-3, mandatory).
//
// The `in` atom sits in the MATCH block because `normalize.go::expand` is a
// match-block mechanism — a guard `in` is one predicate and expands to no
// second row.
var dtExpanding = dtHeader + `
[tags.c]
provenance = "observed"
kind = "enum"
domain = ["m", "n"]
single_valued = true
required = true

[[rule]]
id = "spread"
[rule.match.recognized]
eq = "decide"
[rule.match.c]
in = ["m", "n"]
[rule.guard.all.a]
eq = "x"
[rule.guard.all.b]
eq = "p"
[rule.emit]
verdict = "shared"
note = "same-block"
` + dtRule("cell-xq", "x", "q", nil) +
	dtRule("cell-yp", "y", "p", nil) +
	dtRule("cell-yq", "y", "q", nil)

// smHeader is a minimal STATE-MACHINE preamble carrying one owned tag, its
// accessors, and an `[initial]`. It is the control every class-conditioned
// assertion is taken against.
const smHeader = `outcomes = ["decide"]

[model]
id = "sm"
version = 1
class = "state-machine"

[initial]
status = "draft"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "done"]
single_valued = true
required = true

[read.state]
role = "state"
path = "s.own"
keys = ["status"]
timeout = "2s"

[write.state]
role = "state"
path = "s.own"
keys = ["status"]
timeout = "2s"
read_back = true
`

// smNoWriteBlock is a STATE-MACHINE model whose ordinary rule carries no
// write block. It must refuse `malformed rule shape` — the arm C2 leaves
// live for this class (SC-2).
const smNoWriteBlock = smHeader + `
[[rule]]
id = "no-write"
[rule.match.recognized]
eq = "decide"
[rule.match.status]
eq = "draft"
`

// smZeroOwned is the SC-1 counterpart control: a `state-machine` model
// declaring ZERO owned tags. It MUST LOAD — the agreement check is
// one-directional (C1, REQ-5, REQ-89) — and is left to lint, which reports
// it as a rootless machine.
//
// DEVIATION (deviations.md D1): its one rule is an ESCAPE row. A
// zero-owned model's ORDINARY rule cannot carry a write block (there is no
// owned tag to write), and `0002:C4`'s write-block arm still binds this
// class (REQ-16), so an escape row is the only rule shape a zero-owned
// state machine can author. That is what makes the control writable at all
// — and the write-block arm being live here is itself the REQ-16
// discrimination, taken separately on `smNoWriteBlock`.
const smZeroOwned = `outcomes = ["decide"]

[model]
id = "sm0"
version = 1
class = "state-machine"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.a]
provenance = "observed"
kind = "enum"
domain = ["x", "y"]
single_valued = true
required = true

[[rule]]
id = "only"
escape = ["no_match"]
[rule.match.recognized]
eq = "decide"
[rule.emit]
verdict = "zero-owned"
`

// --- oracles -------------------------------------------------------------

// loadSource loads authored TOML that must load clean.
func loadSource(t *testing.T, src, id string) *table.Model {
	t.Helper()

	m, err := table.Load([]byte(src), id)
	if err != nil {
		t.Fatalf("%s must load clean; refused: %v\n--- source ---\n%s", id, err, src)
	}
	if m == nil {
		t.Fatalf("%s loaded a nil model with no error", id)
	}
	return m
}

// refuseSource loads authored TOML expected to refuse and returns the
// categorized failure. The oracle is the CATEGORY, never message text
// (`0002:C24`).
func refuseSource(t *testing.T, src, id string) *table.Failure {
	t.Helper()

	_, err := table.Load([]byte(src), id)
	if err == nil {
		t.Fatalf("%s loaded clean; want a refusal\n--- source ---\n%s", id, src)
	}
	var f *table.Failure
	if !asFailure(err, &f) {
		t.Fatalf("%s refused with a non-categorized error: %v", id, err)
	}
	return f
}

// asFailure unwraps a categorized load refusal.
func asFailure(err error, target **table.Failure) bool {
	f, ok := err.(*table.Failure)
	if ok {
		*target = f
	}
	return ok
}

// emitOf returns a row's emit sequence as a key->value map plus the key
// order it was carried in, so both content and ORDER can be asserted.
func emitOf(r table.Row) (map[string]string, []string) {
	vals := make(map[string]string, len(r.Emit))
	order := make([]string, 0, len(r.Emit))
	for _, e := range r.Emit {
		vals[e.Key] = e.Value
		order = append(order, e.Key)
	}
	return vals, order
}
