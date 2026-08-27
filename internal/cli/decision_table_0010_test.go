package cli

// RDR 0010 §D — the `flow resolve` payload over a decision-table model
// (`0010:C4`), plus the negative REQs naming the sites and files that must
// not change.
//
// Every invocation drives the production Cobra path through
// `ExecuteAndEmit`, per the harness rules `flow_harness_0005_test.go`
// fixes. Nothing calls a renderer or the kernel directly.

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
	"github.com/newcoinc/intrastate/internal/table"
)

// --- the decision-table fixture ------------------------------------------

// dtModel0010 is the RDR 0010 decision-table fixture, shaped by `0010:MVV`
// and SC-5:
//
//   - `class = "decision-table"`, no `[initial]`, no accessors, no owned tag;
//   - two observed enum dimensions of two values each, `single_valued` and
//     `required`, DISCRIMINATED BY `[rule.guard.all.<key>]` atoms;
//   - `[rule.match.recognized]` binds the one outcome;
//   - four ordinary rules closing the product;
//   - the fourth row (`cell-yq`) authors NO `[rule.emit]`, which is the
//     only arm on which `"emit": null` can surface (SC-5).
const dtModel0010 = `outcomes = ["decide"]

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

[[rule]]
id = "cell-xp"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "x"
[rule.guard.all.b]
eq = "p"
[rule.emit]
verdict = "alpha"
code = "1"

[[rule]]
id = "cell-xq"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "x"
[rule.guard.all.b]
eq = "q"
[rule.emit]
verdict = "beta"

[[rule]]
id = "cell-yp"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "y"
[rule.guard.all.b]
eq = "p"
[rule.emit]
verdict = "gamma"

[[rule]]
id = "cell-yq"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "y"
[rule.guard.all.b]
eq = "q"
`

// dtEscapeModel0010 replaces the fourth ordinary rule with an escape row
// rescuing `no_match`, carrying its OWN `[rule.emit]`. A plan rescued by
// that row must carry the escape row's block (C4).
//
// It carries TWO escape rows, so both halves of REQ-44 are reachable:
//
//   - `otherwise` (`no_match`, HAS an emit block) — rescues the (y,q)
//     request, which no ordinary row covers;
//   - `bare-otherwise` (`ambiguous_match`, authors NO emit block) — rescues
//     the (x,p) request, where `cell-xp` and `cell-xp-shadow` BOTH select.
//
// `cell-xp-shadow` exists only to make that second rescue reachable: an
// `ambiguous_match` needs two rows selecting simultaneously
// (`internal/resolve/resolve.go`), and with a single ordinary row the
// `bare-otherwise` arm — REQ-44's "an escape row authoring no `[rule.emit]`
// renders `{}`" clause — is structurally dead. It duplicates `cell-xp`'s
// guard atoms exactly and touches no other cell, so the (y,q) `no_match`
// selection above is undisturbed.
const dtEscapeModel0010 = `outcomes = ["decide"]

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

[[rule]]
id = "cell-xp"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "x"
[rule.guard.all.b]
eq = "p"
[rule.emit]
verdict = "alpha"

[[rule]]
id = "cell-xp-shadow"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "x"
[rule.guard.all.b]
eq = "p"
[rule.emit]
verdict = "shadow"

[[rule]]
id = "otherwise"
escape = ["no_match"]
[rule.match.recognized]
eq = "decide"
[rule.emit]
verdict = "fallback"

[[rule]]
id = "bare-otherwise"
escape = ["ambiguous_match"]
[rule.match.recognized]
eq = "decide"
`

// dtGateModel0010 is the gated decision-table fixture. It exists because
// REQ-42's clause — "gate evaluation is unchanged and PRIOR to the emit
// join … a deny is `flow-gate-denied`, so `emit` is never computed on a
// denied selection" — has no witness on an ungated model: `dtModel0010`
// declares no `[gate.*]` and no rule carries `gate = [...]`, so `gates` is
// always `[]` and every loop over it is vacuous.
//
// The gates key the OBSERVED tag `a`, not an owned one. `0010:C2` forbids
// only accessors whose `keys` name an OWNED tag, and
// `internal/table/load.go::accessorTable` imposes no owned requirement on
// gates while `checkAccessorBindings` arity-checks readers and writers
// only — so a gate over an observed tag is authorable on this class. (This
// overturns `verification.md` §Undecidable's REQ-42 entry; see
// deviations.md D17.)
//
// The verdict a gate answers is its DECLARED PATH's suffix
// (`internal/cli/flowbind/flowbind.go::verdictFor`), so the two rows below
// differ only in which gate they carry.
const dtGateModel0010 = `outcomes = ["decide"]

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

[gate.ok]
role = "state"
path = "flow.gate.allow"
keys = ["a"]
timeout = "2s"

[gate.nope]
role = "state"
path = "flow.gate.deny"
keys = ["a"]
timeout = "2s"

[[rule]]
id = "allowed"
gate = ["ok"]
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "x"
[rule.emit]
verdict = "alpha"

[[rule]]
id = "denied"
gate = ["nope"]
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "y"
[rule.emit]
verdict = "never-computed"
`

// dtExpandingModel0010 authors an emit block on an EXPANDING rule: a
// multi-member `in` MATCH atom mints several rows from one rule. SC-3
// requires the byte-for-byte assertion be taken again through `flow
// resolve` on a selection landing on a NON-FIRST expanded row.
const dtExpandingModel0010 = `outcomes = ["decide"]

[model]
id = "dt"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

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
[rule.emit]
verdict = "shared"
note = "same-block"
`

// --- REQs ----------------------------------------------------------------

// REQ-38: "The `flow resolve` success payload MUST carry `emit`: a JSON
// object of string values, keys in byte order, present as `{}` — never
// `null`, never omitted — when the selected row authored none."
// REQ-98 / REQ-99 / SC-5: "The fixture MUST include one selectable row
// authoring **no** `[rule.emit]` block, and the `{}`-never-`null`
// assertion MUST be taken on that row's resolve"
// HAPPY PATH + BOUNDARY.
func TestReq38_ResolveCarriesEmitAsAnObjectNeverNullNeverOmitted(t *testing.T) {
	model := writeFlowModel(t, dtModel0010)

	t.Run("an authored block is returned", func(t *testing.T) {
		data := flowData(t, requireSuccess(t, "flow", "resolve",
			"--model", model, "--outcome", "decide",
			"--tag", "a=x", "--tag", "b=p", "--as=json"))

		if data["rule"] != "cell-xp" {
			t.Fatalf("`rule` = %#v; want %q", data["rule"], "cell-xp")
		}
		emit, ok := data["emit"].(map[string]any)
		if !ok {
			t.Fatalf("`emit` is not a JSON object: %#v (payload keys: %v)",
				data["emit"], keysOf(data))
		}
		want := map[string]any{"verdict": "alpha", "code": "1"}
		if !reflect.DeepEqual(emit, want) {
			t.Errorf("`emit` = %#v; want %#v — the authored block, verbatim",
				emit, want)
		}
	})

	t.Run("an unauthored block is {} and never null", func(t *testing.T) {
		stdout := requireSuccess(t, "flow", "resolve",
			"--model", model, "--outcome", "decide",
			"--tag", "a=y", "--tag", "b=q", "--as=json")
		data := flowData(t, stdout)

		if data["rule"] != "cell-yq" {
			t.Fatalf("`rule` = %#v; want %q — the row authoring NO emit block",
				data["rule"], "cell-yq")
		}
		raw, present := data["emit"]
		if !present {
			t.Fatalf("`emit` is OMITTED for a row authoring none; it must be "+
				"present as `{}` (payload keys: %v)", keysOf(data))
		}
		if raw == nil {
			t.Fatal("`emit` is null for a row authoring none; it must be `{}`")
		}
		emit, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("`emit` is not a JSON object: %#v", raw)
		}
		if len(emit) != 0 {
			t.Errorf("`emit` = %#v; want the empty object", emit)
		}

		// The BYTES, not only the decoded shape: `"emit":null` decodes to a
		// nil `any` above, but so would a missing key on some readers.
		if strings.Contains(stdout, `"emit":null`) {
			t.Errorf("the wire carries `\"emit\":null`:\n%s", stdout)
		}
		if !strings.Contains(stdout, `"emit":{}`) {
			t.Errorf("the wire does not carry `\"emit\":{}`:\n%s", stdout)
		}
	})
}

// REQ-35 / REQ-36: "`emit` on the wire is a JSON object of string→string,
// keys in byte order, HTML escaping disabled like every other payload map
// (`0005:C1`)"
// BOUNDARY — determinism and the one-encoder rule.
func TestReq35_EmitOnTheWireIsKeyOrderedWithHTMLEscapingDisabled(t *testing.T) {
	src := strings.Replace(dtModel0010,
		"[rule.emit]\nverdict = \"alpha\"\ncode = \"1\"\n",
		"[rule.emit]\nzeta = \"a<b\"\nalpha = \"x&y\"\nmu = \"m\"\n", 1)
	if src == dtModel0010 {
		t.Fatal("the emit substitution did not apply")
	}
	model := writeFlowModel(t, src)

	stdout := requireSuccess(t, "flow", "resolve",
		"--model", model, "--outcome", "decide",
		"--tag", "a=x", "--tag", "b=p", "--as=json")

	t.Run("keys in byte order", func(t *testing.T) {
		want := `"emit":{"alpha":"x&y","mu":"m","zeta":"a<b"}`
		if !strings.Contains(stdout, want) {
			t.Errorf("wire does not carry %s:\n%s", want, stdout)
		}
	})

	t.Run("HTML escaping disabled", func(t *testing.T) {
		for _, bad := range []string{`\u003c`, `\u0026`} {
			if strings.Contains(stdout, bad) {
				t.Errorf("wire carries the HTML-escaped sequence %s; every "+
					"payload map disables HTML escaping (`0005:C1`):\n%s",
					bad, stdout)
			}
		}
	})

	t.Run("repeated invocations are byte-identical", func(t *testing.T) {
		again := requireSuccess(t, "flow", "resolve",
			"--model", model, "--outcome", "decide",
			"--tag", "a=x", "--tag", "b=p", "--as=json")
		if again != stdout {
			t.Errorf("two identical invocations differ:\n%s\n---\n%s",
				stdout, again)
		}
	})
}

// REQ-39 / REQ-40: "Its **position is fixed**: `emit` is declared
// immediately after `Gates` in
// `internal/cli/flow_resolve.go::resolvePayload`, taking that struct from
// thirteen fields to fourteen; the declaration order is the JSON key order
// Go's encoder emits." / "`Gates`, not `Rule`, is `emit`'s predecessor"
// BOUNDARY — the assertion is on the WIRE key order the position produces.
func TestReq39_EmitSitsImmediatelyAfterGatesOnTheWire(t *testing.T) {
	model := writeFlowModel(t, dtModel0010)
	stdout := requireSuccess(t, "flow", "resolve",
		"--model", model, "--outcome", "decide",
		"--tag", "a=x", "--tag", "b=p", "--as=json")

	got := wireKeyOrder(t, stdout)
	want := []string{
		"model", "revision", "observed", "owned", "readers", "outcome",
		"rule", "gates", "emit", "next", "writes", "clear", "escaped",
	}
	if !slices.Equal(got, want) {
		t.Errorf("payload key order = %v;\nwant                 %v\n"+
			"— `emit` is declared immediately after `Gates`, and the "+
			"pre-edit order is not otherwise disturbed", got, want)
	}

	t.Run("the struct declares fourteen fields", func(t *testing.T) {
		rt := reflect.TypeOf(resolvePayload{})
		if rt.NumField() != 14 {
			t.Errorf("resolvePayload declares %d fields; C4 takes it from "+
				"thirteen to fourteen", rt.NumField())
		}
		gates, ok := rt.FieldByName("Gates")
		if !ok {
			t.Fatal("resolvePayload declares no `Gates` field")
		}
		emit, ok := rt.FieldByName("Emit")
		if !ok {
			t.Fatalf("resolvePayload declares no `Emit` field; fields = %v",
				structFieldNames(rt))
		}
		if emit.Index[0] != gates.Index[0]+1 {
			t.Errorf("`Emit` is at index %d and `Gates` at %d; `emit` is "+
				"declared IMMEDIATELY after `Gates`",
				emit.Index[0], gates.Index[0])
		}
	})
}

// REQ-43: "`escaped` is an existing payload field this RDR does not move:
// it stays at its declared position between `clear` and `escape_class`."
// DOMAIN EDGE — negative REQ.
func TestReq43_EscapedStaysBetweenClearAndEscapeClass(t *testing.T) {
	model := writeFlowModel(t, dtEscapeModel0010)
	stdout := requireSuccess(t, "flow", "resolve",
		"--model", model, "--outcome", "decide",
		"--tag", "a=y", "--tag", "b=q", "--as=json")

	got := wireKeyOrder(t, stdout)
	clearAt := slices.Index(got, "clear")
	escapedAt := slices.Index(got, "escaped")
	classAt := slices.Index(got, "escape_class")
	if clearAt < 0 || escapedAt < 0 || classAt < 0 {
		t.Fatalf("payload key order = %v; the escaped plan must carry all "+
			"three of clear/escaped/escape_class", got)
	}
	if escapedAt != clearAt+1 || classAt != escapedAt+1 {
		t.Errorf("payload key order = %v; `escaped` stays between `clear` "+
			"and `escape_class`", got)
	}
}

// REQ-41: "`next`, `writes`, `clear`, `owned`, and `readers` keep their
// `0005:C1` shapes and are empty over a decision-table model."
// HAPPY PATH
func TestReq41_TheStateFieldsKeepTheirShapesAndAreEmpty(t *testing.T) {
	model := writeFlowModel(t, dtModel0010)
	data := flowData(t, requireSuccess(t, "flow", "resolve",
		"--model", model, "--outcome", "decide",
		"--tag", "a=x", "--tag", "b=p", "--as=json"))

	for _, key := range []string{"next", "writes", "owned"} {
		obj, ok := data[key].(map[string]any)
		if !ok {
			t.Errorf("`%s` = %#v; it keeps its `0005:C1` object shape",
				key, data[key])
			continue
		}
		if len(obj) != 0 {
			t.Errorf("`%s` = %#v; want empty over a decision table", key, obj)
		}
	}

	for _, key := range []string{"clear", "readers"} {
		arr, ok := data[key].([]any)
		if !ok {
			t.Errorf("`%s` = %#v; it keeps its `0005:C1` array shape",
				key, data[key])
			continue
		}
		if len(arr) != 0 {
			t.Errorf("`%s` = %#v; want `[]` over a decision table", key, arr)
		}
	}
}

// REQ-42: "Gate evaluation is **unchanged and prior to the emit join** …
// a deny is `flow-gate-denied` — a refusal, never a payload, so `emit` is
// never computed on a denied selection."
// DOMAIN EDGE — negative REQ: no gate behaviour is added.
//
// Both halves are taken on the DECISION-TABLE class, over `dtGateModel0010`
// (see that fixture on why a gate is authorable there). The allow half also
// pins the wire spelling: `internal/cli/flow_exec.go::gateResult` serializes
// `{id, result, reason?}` — there is no `verdict` key on the wire — and it
// `t.Fatal`s on an empty `gates` list so this oracle can never again pass
// vacuously.
func TestReq42_GatesAreEvaluatedBeforeTheEmitJoin(t *testing.T) {
	model := writeFlowModel(t, dtGateModel0010)
	bind := artifactBinding(flowStateRole, newFlowArtifact(t, "state.artifact"))

	t.Run("an allowing gate is reported and the emit join still runs", func(t *testing.T) {
		data := flowData(t, requireSuccess(t, "flow", "resolve",
			"--model", model, "--artifact", bind, "--outcome", "decide",
			"--tag", "a=x", "--as=json"))

		if data["rule"] != "allowed" {
			t.Fatalf("`rule` = %#v; want %q", data["rule"], "allowed")
		}
		gates, ok := data["gates"].([]any)
		if !ok {
			t.Fatalf("`gates` = %#v; it keeps its array shape", data["gates"])
		}
		if len(gates) == 0 {
			t.Fatal("`gates` is empty on a row carrying `gate = [\"ok\"]`; " +
				"an empty list makes every assertion below vacuous")
		}
		if len(gates) != 1 {
			t.Fatalf("`gates` = %#v; the selected row carries exactly one gate",
				gates)
		}
		obj, ok := gates[0].(map[string]any)
		if !ok {
			t.Fatalf("gate entry = %#v; want a JSON object", gates[0])
		}
		want := map[string]any{"id": "ok", "result": "allow"}
		if !reflect.DeepEqual(obj, want) {
			t.Errorf("gate entry = %#v; want %#v — `gateResult` serializes "+
				"`{id, result, reason?}`; there is no `verdict` key on the "+
				"wire, and a success payload's gates all allow", obj, want)
		}

		emit, ok := data["emit"].(map[string]any)
		if !ok {
			t.Fatalf("`emit` is not a JSON object: %#v", data["emit"])
		}
		if !reflect.DeepEqual(emit, map[string]any{"verdict": "alpha"}) {
			t.Errorf("`emit` = %#v; the emit join runs on an ALLOWED "+
				"selection", emit)
		}
	})

	t.Run("a denying gate refuses and no emit is computed", func(t *testing.T) {
		stdout, _, err := runCmd(t, "flow", "resolve",
			"--model", model, "--artifact", bind, "--outcome", "decide",
			"--tag", "a=y", "--as=json")
		if err == nil {
			t.Fatalf("the denied selection succeeded:\n%s", stdout)
		}
		var ce *clierr.CLIError
		if !asCLIError(err, &ce) {
			t.Fatalf("refusal is not a structured CLIError: %v", err)
		}
		if ce.Code != "flow-gate-denied" {
			t.Fatalf("refusal code = %q; want %q — a deny is a refusal, "+
				"never a payload", ce.Code, "flow-gate-denied")
		}
		if got := clierr.ExitCodeFor(err); got != 2 {
			t.Errorf("exit code = %d; want 2", got)
		}
		// `emit` is never COMPUTED on a denied selection: the refusal
		// envelope carries no payload at all, so the emit block the denied
		// row authors cannot appear anywhere on the wire.
		for _, leak := range []string{"never-computed", `"emit"`} {
			if strings.Contains(stdout, leak) {
				t.Errorf("the refusal wire carries %q; `emit` is never "+
					"computed on a denied selection:\n%s", leak, stdout)
			}
		}
	})
}

// REQ-44: "A plan rescued by an escape row carries that escape row's
// **own** `emit`, joined on the same `Plan.RuleID` path as any other
// selection … An escape row authoring no `[rule.emit]` renders `{}` like
// any other."
// DOMAIN EDGE — negative REQ: no escape-specific join arm.
func TestReq44_AnEscapedPlanCarriesTheEscapeRowsOwnEmit(t *testing.T) {
	model := writeFlowModel(t, dtEscapeModel0010)

	t.Run("an escape row's own authored block", func(t *testing.T) {
		data := flowData(t, requireSuccess(t, "flow", "resolve",
			"--model", model, "--outcome", "decide",
			"--tag", "a=y", "--tag", "b=q", "--as=json"))

		if data["escaped"] != true {
			t.Fatalf("`escaped` = %#v; the (y,q) request has no ordinary row "+
				"and is rescued by the `no_match` escape row", data["escaped"])
		}
		if data["rule"] != "otherwise" {
			t.Fatalf("`rule` = %#v; want the rescuing escape row `otherwise`",
				data["rule"])
		}
		emit, ok := data["emit"].(map[string]any)
		if !ok {
			t.Fatalf("`emit` is not a JSON object: %#v", data["emit"])
		}
		want := map[string]any{"verdict": "fallback"}
		if !reflect.DeepEqual(emit, want) {
			t.Errorf("`emit` = %#v; want the ESCAPE row's own block %#v — "+
				"joined on the same `Plan.RuleID` path as any other "+
				"selection", emit, want)
		}
	})

	// "An escape row authoring no `[rule.emit]` renders `{}` like any
	// other." The (x,p) request selects `cell-xp` AND `cell-xp-shadow`, so
	// the `ambiguous_match` escape row `bare-otherwise` rescues — and it
	// authors no block. The assertion is that the join takes the SELECTED
	// escape row's (absent) block, not some other row's: `cell-xp`,
	// `cell-xp-shadow` and `otherwise` all author one, so a join reaching
	// for any of them renders a populated object here.
	t.Run("an escape row authoring no block renders {}", func(t *testing.T) {
		stdout := requireSuccess(t, "flow", "resolve",
			"--model", model, "--outcome", "decide",
			"--tag", "a=x", "--tag", "b=p", "--as=json")
		data := flowData(t, stdout)

		if data["escaped"] != true {
			t.Fatalf("`escaped` = %#v; the (x,p) request selects both "+
				"`cell-xp` and `cell-xp-shadow` and is rescued by the "+
				"`ambiguous_match` escape row", data["escaped"])
		}
		if data["rule"] != "bare-otherwise" {
			t.Fatalf("`rule` = %#v; want the rescuing escape row "+
				"`bare-otherwise`", data["rule"])
		}
		raw, present := data["emit"]
		if !present {
			t.Fatalf("`emit` is OMITTED on an escape row authoring none; it "+
				"must be present as `{}` (payload keys: %v)", keysOf(data))
		}
		if raw == nil {
			t.Fatal("`emit` is null on an escape row authoring none; it must " +
				"be `{}`")
		}
		emit, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("`emit` is not a JSON object: %#v", raw)
		}
		if !reflect.DeepEqual(emit, map[string]any{}) {
			t.Errorf("`emit` = %#v; want the EMPTY object — the rescuing "+
				"escape row authors no block, and the join takes the "+
				"SELECTED row's, not another row's", emit)
		}
		if !strings.Contains(stdout, `"emit":{}`) {
			t.Errorf("the wire does not carry `\"emit\":{}`:\n%s", stdout)
		}
	})
}

// REQ-45 / REQ-99 / SC-5: "Text mode is **not** specified by this RDR:
// `emit` … renders through the generic payload renderer 0005 owns …, which
// emits one path-qualified leaf per line — `emit.<key>: <value>` — and
// renders an empty block as `emit: (none)` via `flatten`'s empty-container
// arm."
// HAPPY PATH — asserted AS BUILT, per ASSUMPTION-10.
func TestReq45_TextModeRendersEmitThroughTheGenericRenderer(t *testing.T) {
	model := writeFlowModel(t, dtModel0010)

	t.Run("one path-qualified leaf per pair", func(t *testing.T) {
		stdout := requireSuccess(t, "flow", "resolve",
			"--model", model, "--outcome", "decide",
			"--tag", "a=x", "--tag", "b=p")

		for _, want := range []string{"emit.verdict: alpha", "emit.code: 1"} {
			if !containsLine(stdout, want) {
				t.Errorf("text output carries no line %q:\n%s", want, stdout)
			}
		}
	})

	t.Run("an empty block renders emit: (none)", func(t *testing.T) {
		stdout := requireSuccess(t, "flow", "resolve",
			"--model", model, "--outcome", "decide",
			"--tag", "a=y", "--tag", "b=q")

		if !containsLine(stdout, "emit: (none)") {
			t.Errorf("text output carries no `emit: (none)` line for a row "+
				"authoring no block; that arm is `flatten`'s empty-container "+
				"arm and is load-bearing for 0005:\n%s", stdout)
		}
	})
}

// REQ-46: "this RDR MUST NOT introduce a per-verb text special case for
// `emit`; no `Overrides` entry against 0005's text rendering is taken, and
// Phase 3 does not touch that file."
// ADVERSARIAL — negative REQ, naming a file that must not be edited.
func TestReq46_NoPerVerbTextSpecialCaseIsAddedForEmit(t *testing.T) {
	fset := token.NewFileSet()
	// Comments are excluded deliberately: `text.go` already says the
	// renderer "emits one path-qualified leaf per line", and a prose match
	// is not a special case. What the clause forbids is CODE in that file
	// keying on the field name.
	f, err := parser.ParseFile(fset, "respond/text.go", nil, 0)
	if err != nil {
		t.Fatalf("parse respond/text.go: %v", err)
	}

	ast.Inspect(f, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}
		if strings.Contains(strings.ToLower(lit.Value), "emit") {
			t.Errorf("internal/cli/respond/text.go carries the string "+
				"literal %s; the generic payload renderer takes no per-verb "+
				"special case for `emit`, and Phase 3 does not touch that "+
				"file", lit.Value)
		}
		return true
	})

	ast.Inspect(f, func(n ast.Node) bool {
		id, ok := n.(*ast.Ident)
		if ok && strings.EqualFold(id.Name, "emit") {
			t.Errorf("internal/cli/respond/text.go declares or references the "+
				"identifier %q; no per-verb special case is added", id.Name)
		}
		return true
	})
}

// REQ-47 / A2: "`--artifact` is not required when the invoked reader set is
// empty … and an `--artifact` binding for a role no invoked reader needs
// MUST be ignored."
// INPUT EDGE — ignored, not refused.
func TestReq47_ArtifactIsUnrequiredAndAStrayBindingIsIgnored(t *testing.T) {
	model := writeFlowModel(t, dtModel0010)

	t.Run("no --artifact succeeds", func(t *testing.T) {
		requireSuccess(t, "flow", "resolve",
			"--model", model, "--outcome", "decide",
			"--tag", "a=x", "--tag", "b=p", "--as=json")
	})

	t.Run("a stray binding is ignored, not refused", func(t *testing.T) {
		stray := artifactBinding("nobody", newFlowArtifact(t, "stray.artifact"))
		with := flowData(t, requireSuccess(t, "flow", "resolve",
			"--model", model, "--artifact", stray, "--outcome", "decide",
			"--tag", "a=x", "--tag", "b=p", "--as=json"))
		without := flowData(t, requireSuccess(t, "flow", "resolve",
			"--model", model, "--outcome", "decide",
			"--tag", "a=x", "--tag", "b=p", "--as=json"))

		if !reflect.DeepEqual(with, without) {
			t.Errorf("binding an artifact for a role no invoked reader needs "+
				"changed the payload:\n%#v\n---\n%#v", with, without)
		}
	})
}

// REQ-48: "`--outcome` remains required in both classes."
// ADVERSARIAL
func TestReq48_OutcomeRemainsRequiredOverADecisionTable(t *testing.T) {
	model := writeFlowModel(t, dtModel0010)

	_, _, err := runCmd(t, "flow", "resolve",
		"--model", model, "--tag", "a=x", "--tag", "b=p", "--as=json")
	if err == nil {
		t.Fatal("`flow resolve` succeeded with no --outcome over a decision " +
			"table; --outcome remains required in BOTH classes")
	}
}

// REQ-49 / REQ-50 / REQ-51 / `0010:BR4`: "`flow next`, `flow read-state`,
// and `flow set-state` are unchanged by this RDR." / "`flow next` carrying
// no `emit` is **deliberate, not an oversight**"
// ADVERSARIAL — negative REQ.
func TestReq50_FlowNextCarriesNoEmit(t *testing.T) {
	model := writeFlowModel(t, dtModel0010)
	data := flowData(t, requireSuccess(t, "flow", "next",
		"--model", model, "--as=json"))

	if _, present := data["emit"]; present {
		t.Errorf("`flow next`'s payload carries `emit`; adding it to the "+
			"candidate preview is `0010:BR4`, rejected (payload keys: %v)",
			keysOf(data))
	}

	candidates, ok := objectsAt(data, "candidates")
	if !ok || len(candidates) == 0 {
		t.Fatalf("`candidates` is empty or malformed: %#v", data["candidates"])
	}
	for _, c := range candidates {
		if _, present := c["emit"]; present {
			t.Errorf("candidate %#v carries `emit`; the candidate preview "+
				"reports what a row would REQUIRE and WRITE, and over a "+
				"decision table that trio is empty", c["rule"])
		}
		req, ok := c["required"].([]any)
		if !ok {
			t.Errorf("candidate %#v carries no `required` array: %#v",
				c["rule"], c["required"])
			continue
		}
		if len(req) != 0 {
			t.Errorf("candidate %#v: `required` = %#v; want empty over a "+
				"decision table (A11)", c["rule"], req)
		}
	}
}

// REQ-51: "`invokedReaders` and `runReaders` are unchanged — the empty
// demand set already makes `--artifact` unrequired"
// DOMAIN EDGE — negative REQ naming the sites that must not change.
//
// Assertable form: over a zero-owned model the demand set is empty, which
// is what `invokedReaders` already computes.
func TestReq51_TheInvokedReaderSetIsEmptyOverADecisionTable(t *testing.T) {
	m := loadFlowModel(t, dtModel0010)

	if got := invokedReaders(m, "decide"); len(got) != 0 {
		t.Errorf("invokedReaders(resolve) = %v; the demand set is keyed on "+
			"OWNED provenance, so zero owned tags yields the empty set", got)
	}
	if got := invokedReaders(m, ""); len(got) != 0 {
		t.Errorf("invokedReaders(next) = %v; want the empty set", got)
	}
}

// REQ-92 / SC-3: "asserted on the expanded rows and again through `flow
// resolve` on a selection that lands on a NON-FIRST expanded row"
// ADVERSARIAL — the mandatory second half of SC-3's expansion case.
func TestReq92_ANonFirstExpandedRowResolvesToTheAuthoredEmitBlock(t *testing.T) {
	model := writeFlowModel(t, dtExpandingModel0010)

	// `n` is the SECOND member of the `in` atom, so the selection lands on
	// the non-first expanded row.
	data := flowData(t, requireSuccess(t, "flow", "resolve",
		"--model", model, "--outcome", "decide", "--tag", "c=n", "--as=json"))

	if data["rule"] != "spread" {
		t.Fatalf("`rule` = %#v; want %q", data["rule"], "spread")
	}
	emit, ok := data["emit"].(map[string]any)
	if !ok {
		t.Fatalf("`emit` is not a JSON object: %#v", data["emit"])
	}
	want := map[string]any{"verdict": "shared", "note": "same-block"}
	if !reflect.DeepEqual(emit, want) {
		t.Errorf("`emit` = %#v; want %#v byte-for-byte on the NON-FIRST "+
			"expanded row", emit, want)
	}
}

// REQ-34: "the existing first-match
// `internal/cli/flow_resolve.go::rowByID` is correct as it stands"
// DOMAIN EDGE — negative REQ: no per-row uniqueness is claimed or needed.
//
// GREEN BY CONSTRUCTION for the "unchanged" half; the discriminating half
// is that the join is sound over an expanding rule — every row a rule
// expands to carries the same block, so first-match returns the right one.
func TestReq34_TheRuleIDJoinIsSoundOverAnExpandingRule(t *testing.T) {
	m := loadFlowModel(t, dtExpandingModel0010)

	rows := 0
	var first []string
	for _, r := range m.Rows {
		if r.RuleID != "spread" {
			continue
		}
		rows++
		got := make([]string, 0, len(r.Emit))
		for _, e := range r.Emit {
			got = append(got, e.Key+"="+e.Value)
		}
		if first == nil {
			first = got
			continue
		}
		if !slices.Equal(first, got) {
			t.Errorf("%s carries emit %v; the first expanded row carries %v — "+
				"the rule-id join is only sound because every row a rule "+
				"expands to carries the SAME block",
				r.Identity(), got, first)
		}
	}
	if rows < 2 {
		t.Fatalf("the `in`-atom rule expanded to %d rows; want at least 2", rows)
	}
}

// REQ-74 / FM: "class/owned-set disagreement → `flow-model-invalid` with a
// `malformed model declaration` finding naming the class and count"
// ADVERSARIAL — the visible failure mode, at the CLI surface.
func TestReq74_AClassOwnedDisagreementSurfacesAsFlowModelInvalid(t *testing.T) {
	src := strings.Replace(dtModel0010, "[tags.a]", `[tags.own]
provenance = "owned"
kind = "scalar"

[tags.a]`, 1)
	model := writeFlowModel(t, src)

	ce := requireRefusal(t, "flow-model-invalid", 2,
		"flow", "resolve", "--model", model, "--outcome", "decide",
		"--tag", "a=x", "--tag", "b=p", "--as=json")

	var named bool
	for _, f := range ce.Findings {
		if strings.Contains(f.Message, "owned=1") ||
			strings.Contains(f.Hint, "owned=1") {
			named = true
		}
	}
	if !named {
		t.Errorf("no finding renders the owned count as the literal token "+
			"`owned=1`; findings = %#v", ce.Findings)
	}
}

// REQ-104 / IP Phase 3: "`emit` on `resolvePayload`, joined by rule id
// after selection, rendered in text mode"
// BOUNDARY — the join is by rule id AFTER selection, never through the
// kernel.
func TestReq104_EmitIsJoinedByRuleIDAfterSelection(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "flow_resolve.go", nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse flow_resolve.go: %v", err)
	}

	// `rowByID` must still be the join, and it must still be first-match.
	var sawRowByID bool
	ast.Inspect(f, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if ok && fn.Name.Name == "rowByID" {
			sawRowByID = true
		}
		return true
	})
	if !sawRowByID {
		t.Error("flow_resolve.go declares no `rowByID`; C3 fixes it as the " +
			"join and says it is correct as it stands")
	}
}

// --- helpers -------------------------------------------------------------

// wireKeyOrder returns the `data` object's key order AS THE ENCODER EMITTED
// IT. Decoding into a map would lose the order, which is the property C4
// fixes, so the raw JSON is scanned instead.
func wireKeyOrder(t *testing.T, stdout string) []string {
	t.Helper()

	var env struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &env); err != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", err, stdout)
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
				if depth < 0 {
					return keys
				}
			}
		case string:
			if depth == 0 {
				keys = append(keys, v)
				// Consume the value.
				if err := skipValue(dec); err != nil {
					return keys
				}
			}
		}
	}
	return keys
}

// skipValue consumes one complete JSON value from dec.
func skipValue(dec *json.Decoder) error {
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	delim, ok := tok.(json.Delim)
	if !ok {
		return nil
	}
	depth := 1
	for depth > 0 {
		tok, err := dec.Token()
		if err != nil {
			return err
		}
		if d, ok := tok.(json.Delim); ok {
			switch d {
			case '{', '[':
				depth++
			case '}', ']':
				depth--
			}
		}
	}
	_ = delim
	return nil
}

// containsLine reports whether stdout carries want as a complete line.
func containsLine(stdout, want string) bool {
	for _, line := range strings.Split(stdout, "\n") {
		if strings.TrimRight(line, "\r") == want {
			return true
		}
	}
	return false
}

func structFieldNames(rt reflect.Type) []string {
	out := make([]string, 0, rt.NumField())
	for i := range rt.NumField() {
		out = append(out, rt.Field(i).Name)
	}
	return out
}

// loadFlowModel loads authored TOML through the real `table.Load` — the
// same loader the CLI hands bytes to — for the two REQs whose oracle is a
// package-internal helper (`invokedReaders`) rather than the CLI surface.
func loadFlowModel(t *testing.T, src string) *table.Model {
	t.Helper()

	m, err := table.Load([]byte(src), "dt.toml")
	if err != nil {
		t.Fatalf("the fixture must load clean; refused: %v", err)
	}
	return m
}
