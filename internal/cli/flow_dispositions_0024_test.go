package cli

// RDR 0024 `0024:C4` / `0024:S5` — ENVELOPE: the `dispositions` payload
// field.
//
// Every assertion drives the ROOT command through the same shape CI uses
// and reads the JSON envelope. The join lives in `internal/cli`, so this
// suite is the only place it is observable.

import (
	"encoding/json"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// dtDecl0024 is the declaration `0024:S5` needs over RDR 0010's decision
// table: `verdict` declared enum, its three authored values partitioned
// into two model-authored dispositions.
//
// The declaration is authored between the header and the rules so the
// top-level table is not parsed into the last `[[rule]]`.
const dtDecl0024 = `
[emit.verdict]
kind = "enum"
[emit.verdict.domain]
route = ["alpha", "beta"]
stop = ["gamma"]
`

// declaredDT0024 is `dtModel0010` with the `[emit]` declaration inserted.
func declaredDT0024(t *testing.T) string {
	t.Helper()

	const anchor = "[[rule]]\nid = \"cell-xp\""
	if !strings.Contains(dtModel0010, anchor) {
		t.Fatalf("the 0010 decision-table fixture no longer carries %q", anchor)
	}
	return strings.Replace(dtModel0010, anchor, dtDecl0024+"\n"+anchor, 1)
}

// resolveDispositions0024 runs `flow resolve` over src with the supplied
// discriminating tags and returns the decoded payload.
func resolveDispositions0024(t *testing.T, src string, tags ...string) map[string]any {
	t.Helper()

	model := writeFlowModel(t, src)
	args := []string{"flow", "resolve", "--model", model, "--outcome", "decide"}
	for _, tag := range tags {
		args = append(args, "--tag", tag)
	}
	return flowData(t, requireSuccess(t, append(args, "--as=json")...))
}

// dispositionsOf reads the payload's `dispositions` object, failing when
// the key is absent or null — the two forms C4 forbids.
func dispositionsOf(t *testing.T, data map[string]any) map[string]string {
	t.Helper()

	raw, present := data["dispositions"]
	if !present {
		t.Fatalf("the payload carries no `dispositions` key; it is present "+
			"as `{}` — never null, never omitted. keys = %v", keysOf(data))
	}
	if raw == nil {
		t.Fatal("`dispositions` decoded to null; it is present as `{}` — " +
			"never `null`, never omitted")
	}
	obj, ok := raw.(map[string]any)
	if !ok {
		t.Fatalf("`dispositions` is %T; want a JSON object", raw)
	}
	out := make(map[string]string, len(obj))
	for k, v := range obj {
		s, ok := v.(string)
		if !ok {
			t.Errorf("dispositions[%q] is %T; the token is a string", k, v)
			continue
		}
		out[k] = s
	}
	return out
}

// REQ-57: "ENVELOPE. The `flow resolve` success payload gains
// `dispositions`: a JSON object mapping emit key → the disposition token
// the declaration assigns the selected row's authored value, keys in byte
// order, present as `{}` — never `null`, never omitted — when no selected
// value carries one (undeclared model, non-enum kind, flat-array domain, or
// empty emit block alike)."
// HAPPY PATH
func TestReq57_0024_ResolveCarriesDispositionsForTheSelectedValue(t *testing.T) {
	// (x,p) selects `cell-xp`, which emits `verdict = "alpha"`, listed
	// under `route`.
	data := resolveDispositions0024(t, declaredDT0024(t), "a=x", "b=p")

	got := dispositionsOf(t, data)
	want := map[string]string{"verdict": "route"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("dispositions = %v; want %v — the token the declaration "+
			"assigns the selected row's authored value", got, want)
	}
}

// REQ-57 tail / REQ-102 / `0024:S5`: "`dispositions` maps only
// member-disposition hits and is `{}` in every other case — never `null`,
// never omitted".
// BOUNDARY
func TestReq57_0024_DispositionsIsEmptyObjectInEveryOtherCase(t *testing.T) {
	cases := []struct {
		name string
		src  func(*testing.T) string
		tags []string
	}{
		{
			name: "undeclared model",
			src:  func(*testing.T) string { return dtModel0010 },
			tags: []string{"a=x", "b=p"},
		},
		{
			name: "empty emit block on the selected row",
			// `cell-yq` authors NO `[rule.emit]`.
			src:  declaredDT0024,
			tags: []string{"a=y", "b=q"},
		},
		{
			name: "flat-array domain carries no dispositions",
			src: func(t *testing.T) string {
				return strings.Replace(declaredDT0024(t), dtDecl0024, `
[emit.verdict]
kind = "enum"
domain = ["alpha", "beta", "gamma"]
`, 1)
			},
			tags: []string{"a=x", "b=p"},
		},
		{
			name: "non-enum kind",
			src: func(t *testing.T) string {
				return strings.Replace(declaredDT0024(t), dtDecl0024, `
[emit.verdict]
kind = "scalar"
`, 1)
			},
			tags: []string{"a=x", "b=p"},
		},
		{
			name: "a bare `[emit]` table",
			src: func(t *testing.T) string {
				return strings.Replace(declaredDT0024(t), dtDecl0024,
					"\n[emit]\n", 1)
			},
			tags: []string{"a=x", "b=p"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data := resolveDispositions0024(t, tc.src(t), tc.tags...)
			if got := dispositionsOf(t, data); len(got) != 0 {
				t.Errorf("dispositions = %v; want `{}` for the %s case",
					got, tc.name)
			}
		})
	}
}

// REQ-58: "The field is `Dispositions map[string]string` with tag
// `json:\"dispositions\"` and no `omitempty`."
// REQ-63: "**`omitempty` is available on this field and is deliberately not
// taken**"
// DOMAIN EDGE
func TestReq58_0024_TheFieldIsAStringMapWithNoOmitempty(t *testing.T) {
	rt := reflect.TypeOf(resolvePayload{})
	f, ok := rt.FieldByName("Dispositions")
	if !ok {
		t.Fatalf("resolvePayload declares no `Dispositions` field; fields = %v",
			structFieldNames(rt))
	}
	if got, want := f.Type.String(), "map[string]string"; got != want {
		t.Errorf("Dispositions is %s; want %s", got, want)
	}
	tag := f.Tag.Get("json")
	if tag != "dispositions" {
		t.Errorf("json tag = %q; want %q — `omitempty` is available and is "+
			"deliberately NOT taken", tag, "dispositions")
	}
}

// REQ-59: "**Byte order is the marshaller's, not the join's**:
// `encoding/json` sorts map keys on marshal, so the join performs no sort
// and the carried value is an ordinary unordered Go map … An implementer
// must not add a sort over this map"
// DOMAIN EDGE
func TestReq59_0024_DispositionKeysArriveInByteOrderFromTheMarshaller(t *testing.T) {
	// Two declared keys whose byte order is the REVERSE of the order the
	// rule authors them in. The marshaller's sort is the only thing that
	// can produce the wire order.
	src := strings.Replace(declaredDT0024(t), dtDecl0024, `
[emit.verdict]
kind = "enum"
[emit.verdict.domain]
route = ["alpha", "beta"]
stop = ["gamma"]

[emit.aaa]
kind = "enum"
[emit.aaa.domain]
terminal = ["zzz"]
`, 1)
	src = strings.Replace(src, "verdict = \"alpha\"\ncode = \"1\"",
		"verdict = \"alpha\"\naaa = \"zzz\"", 1)
	if strings.Contains(src, "code = \"1\"") {
		t.Fatal("the emit substitution did not apply")
	}

	model := writeFlowModel(t, src)
	stdout := requireSuccess(t, "flow", "resolve", "--model", model,
		"--outcome", "decide", "--tag", "a=x", "--tag", "b=p", "--as=json")

	got := dispositionKeyOrder(t, stdout)
	want := []string{"aaa", "verdict"}
	if !slices.Equal(got, want) {
		t.Errorf("`dispositions` key order = %v; want %v — byte order is "+
			"the MARSHALLER's, and the join performs no sort", got, want)
	}
}

// REQ-60: "**The map is keyed off the SELECTED ROW's authored emit, never
// off the declaration set.** An entry exists for key `k` exactly when the
// selected row authors `k` AND `k`'s declaration lists that authored value
// under a disposition. A declared key the selected row does not emit
// contributes NO entry — `dispositions` is never padded with nulls, empty
// strings, or absent-markers to the declared key set."
// REQ-102 tail / `0024:S5`: "a model declaring two keys whose selected row
// emits one yields exactly one `dispositions` entry"
// ADVERSARIAL
func TestReq60_0024_TheMapIsKeyedOffTheSelectedRowNotTheDeclarationSet(t *testing.T) {
	// Two declared keys; `cell-xq` emits only `verdict = "beta"`.
	src := strings.Replace(declaredDT0024(t), dtDecl0024, `
[emit.verdict]
kind = "enum"
[emit.verdict.domain]
route = ["alpha", "beta"]
stop = ["gamma"]

[emit.code]
kind = "enum"
[emit.code.domain]
terminal = ["1"]

[emit.never-emitted]
kind = "enum"
[emit.never-emitted.domain]
ghost = ["nothing"]
`, 1)

	data := resolveDispositions0024(t, src, "a=x", "b=q")
	got := dispositionsOf(t, data)
	want := map[string]string{"verdict": "route"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("dispositions = %v; want %v — a declared key the selected "+
			"row does not emit contributes NO entry, and the map is never "+
			"padded to the declared key set", got, want)
	}
}

// REQ-60 second leg: an entry exists only when the AUTHORED VALUE carries a
// disposition, not merely when the key is declared.
// BOUNDARY
func TestReq60_0024_ADeclaredKeyWhoseValueCarriesNoDispositionContributesNoEntry(t *testing.T) {
	// `verdict` is declared with a disposition table listing only `alpha`
	// and `beta`; `gamma` sits in a FLAT continuation, so the model must
	// list it — instead declare `code` flat (no dispositions) and let
	// `cell-xp` emit both.
	src := strings.Replace(declaredDT0024(t), dtDecl0024, `
[emit.verdict]
kind = "enum"
[emit.verdict.domain]
route = ["alpha", "beta"]
stop = ["gamma"]

[emit.code]
kind = "enum"
domain = ["1"]
`, 1)

	data := resolveDispositions0024(t, src, "a=x", "b=p")
	got := dispositionsOf(t, data)
	if _, present := got["code"]; present {
		t.Errorf("dispositions = %v; `code` is declared with a FLAT domain "+
			"carrying no dispositions, so its authored value contributes "+
			"no entry", got)
	}
	if got["verdict"] != "route" {
		t.Errorf("dispositions[\"verdict\"] = %q; want %q",
			got["verdict"], "route")
	}
}

// REQ-61: "It is inserted at one fixed position in
// `internal/cli/flow_resolve.go::resolvePayload`, immediately after `emit`"
// REQ-68: "The insertion takes `resolvePayload` from fourteen struct fields
// to **fifteen** and the wire from thirteen keys to **fourteen**, with
// `dispositions` at index 9"
// REQ-103 / `0024:S5`: the post-change wire key order is exactly `model,
// revision, observed, owned, readers, outcome, rule, gates, emit,
// dispositions, next, writes, clear, escaped` "and
// `reflect.TypeOf(resolvePayload{}).NumField() == 15`"
// DOMAIN EDGE
func TestReq103_0024_DispositionsSitsImmediatelyAfterEmitOnTheWire(t *testing.T) {
	model := writeFlowModel(t, declaredDT0024(t))
	stdout := requireSuccess(t, "flow", "resolve", "--model", model,
		"--outcome", "decide", "--tag", "a=x", "--tag", "b=p", "--as=json")

	got := wireKeyOrder(t, stdout)
	want := []string{
		"model", "revision", "observed", "owned", "readers", "outcome",
		"rule", "gates", "emit", "dispositions", "next", "writes", "clear",
		"escaped",
	}
	if !slices.Equal(got, want) {
		t.Errorf("payload key order = %v;\nwant                 %v\n"+
			"— `dispositions` is inserted at index 9, immediately after "+
			"`emit`, and the pre-edit order is not otherwise disturbed",
			got, want)
	}
	if len(want) != 14 {
		t.Fatalf("the expected key list carries %d keys; C4 takes the wire "+
			"from thirteen to FOURTEEN", len(want))
	}

	t.Run("the struct declares fifteen fields", func(t *testing.T) {
		rt := reflect.TypeOf(resolvePayload{})
		if rt.NumField() != 15 {
			t.Errorf("resolvePayload declares %d fields; C4 takes it from "+
				"fourteen to fifteen", rt.NumField())
		}
		emit, ok := rt.FieldByName("Emit")
		if !ok {
			t.Fatal("resolvePayload declares no `Emit` field")
		}
		disp, ok := rt.FieldByName("Dispositions")
		if !ok {
			t.Fatalf("resolvePayload declares no `Dispositions` field; "+
				"fields = %v", structFieldNames(rt))
		}
		if disp.Index[0] != emit.Index[0]+1 {
			t.Errorf("`Dispositions` is at index %d and `Emit` at %d; "+
				"`dispositions` is declared IMMEDIATELY after `emit`",
				disp.Index[0], emit.Index[0])
		}
		if disp.Index[0] != 9 {
			t.Errorf("`Dispositions` is at struct index %d; C4 fixes it at "+
				"index 9", disp.Index[0])
		}
	})
}

// REQ-62: "**The join itself is performed in `internal/cli`, reading
// `Model.EmitDecls`; `internal/table` gains no payload-shaped helper.** … a
// method on `*Model` returning the payload's map would shape the table
// package's API around the CLI's wire format"
// ADVERSARIAL — negative REQ.
func TestReq62_0024_TheTablePackageGainsNoPayloadShapedHelper(t *testing.T) {
	// The observable: `*table.Model` declares no method returning the
	// payload's `map[string]string` disposition map. Reflection over the
	// real type is the oracle a rename cannot dodge.
	mt := reflect.TypeOf(loadFlowModel(t, declaredDT0024(t)))
	for i := range mt.NumMethod() {
		m := mt.Method(i)
		if m.Type.NumOut() != 1 {
			continue
		}
		out := m.Type.Out(0)
		if out.Kind() != reflect.Map ||
			out.Key().Kind() != reflect.String ||
			out.Elem().Kind() != reflect.String {
			continue
		}
		t.Errorf("`*table.Model` declares `%s` returning %s; the join is "+
			"performed in `internal/cli` and `internal/table` gains NO "+
			"payload-shaped helper", m.Name, out)
	}
}

// REQ-64: "The token is surfaced verbatim; intrastate never interprets it."
// ADVERSARIAL
func TestReq64_0024_TheDispositionTokenIsSurfacedVerbatim(t *testing.T) {
	// Tokens intrastate could plausibly want to normalize: mixed case,
	// inner spaces, a leading dash, a token that collides with an existing
	// escape class name.
	for _, token := range []string{
		"Route", "  padded  ", "-dash", "no_match", "STOP/hard",
	} {
		t.Run(token, func(t *testing.T) {
			src := strings.Replace(declaredDT0024(t), dtDecl0024, `
[emit.verdict]
kind = "enum"
[emit.verdict.domain]
"`+token+`" = ["alpha", "beta"]
stop = ["gamma"]
`, 1)
			data := resolveDispositions0024(t, src, "a=x", "b=p")
			if got := dispositionsOf(t, data)["verdict"]; got != token {
				t.Errorf("dispositions[\"verdict\"] = %q; want %q verbatim "+
					"— intrastate NEVER interprets a disposition token",
					got, token)
			}
		})
	}
}

// REQ-65: "Text mode renders through the generic payload renderer as
// `dispositions.<key>: <token>` / `dispositions: (none)` with no per-verb
// special case (the `0010:C4` clause)."
// HAPPY PATH
func TestReq65_0024_TextModeRendersThroughTheGenericPayloadRenderer(t *testing.T) {
	t.Run("a joined disposition", func(t *testing.T) {
		model := writeFlowModel(t, declaredDT0024(t))
		stdout := requireSuccess(t, "flow", "resolve", "--model", model,
			"--outcome", "decide", "--tag", "a=x", "--tag", "b=p")
		if !strings.Contains(stdout, "dispositions.verdict: route") {
			t.Errorf("text mode does not render "+
				"`dispositions.verdict: route`:\n%s", stdout)
		}
	})

	t.Run("no joined disposition renders (none)", func(t *testing.T) {
		// `cell-yq` authors no emit block.
		model := writeFlowModel(t, declaredDT0024(t))
		stdout := requireSuccess(t, "flow", "resolve", "--model", model,
			"--outcome", "decide", "--tag", "a=y", "--tag", "b=q")
		if !strings.Contains(stdout, "dispositions: (none)") {
			t.Errorf("text mode does not render `dispositions: (none)`:\n%s",
				stdout)
		}
	})
}

// REQ-66: "A plan rescued by an escape row joins `dispositions` from that
// escape row's OWN authored values — the same single `Plan.RuleID` join
// path as `emit` (`0010:C4`), so no defaulted or merged value can reach the
// payload unjoined"
// DOMAIN EDGE
func TestReq66_0024_AnEscapeRescueJoinsFromTheEscapeRowsOwnValues(t *testing.T) {
	const decl = `
[emit.verdict]
kind = "enum"
[emit.verdict.domain]
route = ["alpha", "shadow"]
stop = ["fallback"]
`
	const anchor = "[[rule]]\nid = \"cell-xp\""
	if !strings.Contains(dtEscapeModel0010, anchor) {
		t.Fatalf("the 0010 escape fixture no longer carries %q", anchor)
	}
	src := strings.Replace(dtEscapeModel0010, anchor, decl+"\n"+anchor, 1)

	// (y,q) is covered by no ordinary row, so the `otherwise` escape row
	// rescues it with its OWN `verdict = "fallback"`.
	data := resolveDispositions0024(t, src, "a=y", "b=q")

	if escaped, _ := data["escaped"].(bool); !escaped {
		t.Fatalf("the (y,q) request was not rescued by an escape row; "+
			"payload = %v", data)
	}
	got := dispositionsOf(t, data)
	want := map[string]string{"verdict": "stop"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("dispositions = %v; want %v — the rescue joins from the "+
			"ESCAPE ROW's own authored values via the single `Plan.RuleID` "+
			"join path", got, want)
	}
}

// REQ-66 second leg: an escape row authoring NO emit block joins `{}`,
// never a merged or defaulted value from the shadowed ordinary rows.
// BOUNDARY
func TestReq66_0024_AnEscapeRowAuthoringNoEmitJoinsAnEmptyMap(t *testing.T) {
	const decl = `
[emit.verdict]
kind = "enum"
[emit.verdict.domain]
route = ["alpha", "shadow"]
stop = ["fallback"]
`
	src := strings.Replace(dtEscapeModel0010,
		"[[rule]]\nid = \"cell-xp\"", decl+"\n[[rule]]\nid = \"cell-xp\"", 1)

	// (x,p) selects both `cell-xp` and `cell-xp-shadow`, so
	// `bare-otherwise` rescues the ambiguity — and it authors no emit.
	data := resolveDispositions0024(t, src, "a=x", "b=p")
	if got := dispositionsOf(t, data); len(got) != 0 {
		t.Errorf("dispositions = %v; want `{}` — no defaulted or merged "+
			"value can reach the payload unjoined", got)
	}
}

// REQ-67: "`flow next` carries no `dispositions`, for `0010:C4`'s reason:
// the answer is what `resolve` selects."
// ADVERSARIAL — negative REQ.
func TestReq67_0024_FlowNextCarriesNoDispositions(t *testing.T) {
	model := writeFlowModel(t, declaredDT0024(t))
	stdout := requireSuccess(t, "flow", "next", "--model", model,
		"--tag", "a=x", "--tag", "b=p", "--as=json")

	if strings.Contains(stdout, "dispositions") {
		t.Errorf("`flow next` mentions `dispositions`; the answer is what "+
			"`resolve` selects:\n%s", stdout)
	}
	if _, ok := reflect.TypeOf(nextPayload{}).FieldByName("Dispositions"); ok {
		t.Error("`nextPayload` declares a `Dispositions` field; `flow next` " +
			"carries none")
	}
}

// --- helpers -------------------------------------------------------------

// dispositionKeyOrder returns the `dispositions` object's key order AS THE
// ENCODER EMITTED IT. Decoding into a map would lose the order, which is
// the property REQ-59 fixes.
func dispositionKeyOrder(t *testing.T, stdout string) []string {
	t.Helper()

	var env struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &env); err != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", err, stdout)
	}
	var data map[string]json.RawMessage
	if err := json.Unmarshal(env.Data, &data); err != nil {
		t.Fatalf("`data` is not a JSON object: %v\n%s", err, string(env.Data))
	}
	raw, ok := data["dispositions"]
	if !ok {
		t.Fatalf("the payload carries no `dispositions` key:\n%s", string(env.Data))
	}

	dec := json.NewDecoder(strings.NewReader(string(raw)))
	if _, err := dec.Token(); err != nil {
		t.Fatalf("read `dispositions`: %v", err)
	}
	var keys []string
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("read `dispositions` key: %v", err)
		}
		key, ok := tok.(string)
		if !ok {
			t.Fatalf("`dispositions` key token is %T; want a string", tok)
		}
		keys = append(keys, key)
		if _, err := dec.Token(); err != nil {
			t.Fatalf("read `dispositions` value: %v", err)
		}
	}
	return keys
}
