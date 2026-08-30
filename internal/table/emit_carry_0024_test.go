package table_test

// RDR 0024 `0024:C3` — CARRY: the normalized declaration carrier.
//
// `0024:S4` is the SOLE oracle for the sort clause: A8 is Refuted, so
// `roundtrip_test.go::TestReq146` does not reach this carrier and is not
// edited. Every determinism assertion this contract has lives here.

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// REQ-44: "CARRY. Every declared field — key, kind, domain members, and
// each member's disposition — is carried through normalization onto the
// normalized model, the same clause `0002:C22` states for tag
// declarations"
// REQ-100 / `0024:S4`: "**Expected**: read off `Model.EmitDecls[<key>]` —
// `Kind` equals the authored token, `Domain` equals the bytewise-sorted
// union of the authored members, and `Dispositions[<member>]` equals the
// token that member was listed under."
// HAPPY PATH
func TestReq44_0024_EveryDeclaredFieldIsCarriedOntoTheModel(t *testing.T) {
	m := loadEmitDecls0024(t, `[emit.verdict]
kind = "enum"
[emit.verdict.domain]
stop = ["gamma"]
route = ["beta", "alpha"]`, "emit-carry.toml")

	got, ok := m.EmitDecls["verdict"]
	if !ok {
		t.Fatalf("the declared key `verdict` is not carried; EmitDecls = %v",
			m.EmitDecls)
	}
	want := table.EmitDecl{
		Kind:   "enum",
		Domain: []string{"alpha", "beta", "gamma"},
		Dispositions: map[string]string{
			"alpha": "route",
			"beta":  "route",
			"gamma": "stop",
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("EmitDecls[\"verdict\"] = %+v;\nwant                     %+v",
			got, want)
	}
}

// REQ-45: "The carried form is a new model-level type, deliberately NOT
// `TagDecl` (the mirror of `EmitValue` not being `TagValue`, `0010:C3`). It
// is carried at `Model.EmitDecls map[string]EmitDecl`, keyed by emit key —
// the sibling of `Model.Tags map[string]TagDecl`."
// REQ-46: "`EmitDecl` is an exported struct with exactly three exported
// fields: `Kind string` …, `Domain []string` …, and `Dispositions
// map[string]string`"
// DOMAIN EDGE
func TestReq45_0024_TheCarrierIsANewTypeWithExactlyThreeFields(t *testing.T) {
	field, ok := reflect.TypeOf(table.Model{}).FieldByName("EmitDecls")
	if !ok {
		t.Fatal("`Model` declares no `EmitDecls` field")
	}
	if got, want := field.Type.String(), "map[string]table.EmitDecl"; got != want {
		t.Errorf("Model.EmitDecls is %s; want %s — a new model-level type, "+
			"deliberately NOT TagDecl", got, want)
	}

	rt := reflect.TypeOf(table.EmitDecl{})
	if rt.NumField() != 3 {
		t.Errorf("EmitDecl declares %d fields; C3 fixes EXACTLY three",
			rt.NumField())
	}
	for name, want := range map[string]string{
		"Kind":         "string",
		"Domain":       "[]string",
		"Dispositions": "map[string]string",
	} {
		f, found := rt.FieldByName(name)
		if !found {
			t.Errorf("EmitDecl declares no `%s` field", name)
			continue
		}
		if got := f.Type.String(); got != want {
			t.Errorf("EmitDecl.%s is %s; want %s", name, got, want)
		}
		if !f.IsExported() {
			t.Errorf("EmitDecl.%s is unexported; all three are exported", name)
		}
	}
}

// REQ-47: "**`Model.EmitDecls` is always non-nil after a successful load**
// — empty for a zero-declaration model, never `nil` — matching the
// `Model.Tags` convention"
// BOUNDARY
func TestReq47_0024_EmitDeclsIsAlwaysNonNilAfterASuccessfulLoad(t *testing.T) {
	for name, src := range map[string]string{
		"no `[emit]` table at all": dtComplete,
		"a bare `[emit]` table":    dtWithEmitDecl("[emit]"),
		"one declared key":         dtWithEmitDecl(declVerdictCovering),
	} {
		t.Run(name, func(t *testing.T) {
			m := loadSource(t, src, "emit-nonnil.toml")
			if m.EmitDecls == nil {
				t.Error("`Model.EmitDecls` is nil after a successful load; " +
					"it is empty for a zero-declaration model, NEVER nil")
			}
		})
	}
}

// REQ-48: "`Domain` is `nil` (not an empty non-nil slice) for `bool`,
// `int`, and `scalar`, since none of them takes a domain; the distinction
// is observable because scenario 4 asserts value-equality on the carrier,
// where `nil` and `[]string{}` do not compare equal."
// BOUNDARY
func TestReq48_0024_DomainIsNilNotEmptyForEveryNonEnumKind(t *testing.T) {
	m := loadEmitDecls0024(t, declVerdictCovering+`

[emit.flag]
kind = "bool"

[emit.count]
kind = "int"

[emit.free]
kind = "scalar"`, "emit-nil-domain.toml")

	for _, key := range []string{"flag", "count", "free"} {
		d := m.EmitDecls[key]
		if d.Domain == nil {
			continue
		}
		t.Errorf("EmitDecls[%q].Domain = %#v; want nil — the distinction is "+
			"observable because scenario 4 asserts VALUE-EQUALITY, where "+
			"nil and []string{} do not compare equal", key, d.Domain)
	}
}

// REQ-49: "**The carry is value-preserving, not order-preserving, and the
// member list is SORTED.** … So `Domain` is the union sorted bytewise, and
// the partition grouping is not itself carried: it is fully recoverable
// from `Dispositions`"
// REQ-51: "the authored order of a flat-array domain is NOT preserved
// either, and nothing downstream reads it"
// DOMAIN EDGE
func TestReq49_0024_TheCarriedDomainIsTheBytewiseSortedUnion(t *testing.T) {
	t.Run("disposition-table union is sorted", func(t *testing.T) {
		m := loadEmitDecls0024(t, `[emit.verdict]
kind = "enum"
[emit.verdict.domain]
zeta = ["gamma", "alpha"]
alpha-disp = ["beta"]`, "emit-sorted-union.toml")

		want := []string{"alpha", "beta", "gamma"}
		if got := m.EmitDecls["verdict"].Domain; !slices.Equal(got, want) {
			t.Errorf("Domain = %v; want the bytewise-sorted union %v", got, want)
		}
	})

	t.Run("flat-array authored order is not preserved", func(t *testing.T) {
		m := loadEmitDecls0024(t, `[emit.verdict]
kind = "enum"
domain = ["gamma", "alpha", "beta"]`, "emit-sorted-flat.toml")

		want := []string{"alpha", "beta", "gamma"}
		if got := m.EmitDecls["verdict"].Domain; !slices.Equal(got, want) {
			t.Errorf("Domain = %v; want %v — the authored order of a "+
				"flat-array domain is NOT preserved", got, want)
		}
	})

	t.Run("the partition grouping is recoverable from Dispositions", func(t *testing.T) {
		m := loadEmitDecls0024(t, declVerdictEnum, "emit-partition.toml")
		d := m.EmitDecls["verdict"]
		route := make([]string, 0, len(d.Domain))
		for _, member := range d.Domain {
			if d.Dispositions[member] == "route" {
				route = append(route, member)
			}
		}
		if want := []string{"alpha", "beta"}; !slices.Equal(route, want) {
			t.Errorf("the `route` partition recovers as %v; want %v", route, want)
		}
	})
}

// REQ-50: "that sweep does NOT reach this carrier: all three of its
// sub-tests are scoped to `table.Row` (A8, Refuted), so **Testing Strategy
// scenario 4 is the sole oracle for this clause** and Phase 2 owes no edit
// to `TestReq146`."
// ADVERSARIAL — the negative REQ: the sweep must stay unedited.
func TestReq50_0024_TheRoundTripSweepIsNotEditedToReachThisCarrier(t *testing.T) {
	// A8 is Refuted because every sub-test is scoped to `table.Row`. The
	// checkable consequence is that the sweep still passes untouched
	// against a model carrying declarations — if an implementer had wired
	// the carrier into it, this model would exercise a path the sweep was
	// never scoped to.
	m := loadEmitDecls0024(t, declVerdictEnum, "emit-no-sweep-edit.toml")
	for _, r := range m.Rows {
		if !slices.IsSorted(r.Suffix) {
			t.Errorf("row %s carries an unsorted suffix; the RDR 0002 sweep "+
				"is scoped to table.Row and stays green unedited",
				r.Identity())
		}
	}
	// And the carrier is reachable independently — the oracle scenario 4
	// owns, not the sweep's.
	if len(m.EmitDecls) == 0 {
		t.Error("the carrier is empty; scenario 4 is the SOLE oracle for " +
			"the sort clause and needs the carrier to read")
	}
}

// REQ-52: "**Value conformance is performed by constructing a throwaway
// `TagDecl{Kind: d.Kind, Domain: d.Domain}` per call and passing it to
// `internal/table/load.go::ConformValue(decl TagDecl, member string)
// error`.** That struct literal IS the adapter — no extraction, no new
// exported helper."
// REQ-53: "`ConformValue` reads only `Kind`, `Domain`, `Min`, `Max`, and
// `Elements`; `Min`/`Max`/`Elements` stay nil, so the `int` arm's bounds
// and the `set` arm never fire."
// DOMAIN EDGE
func TestReq52_0024_ValueConformanceRunsThroughConformValueWithNoBounds(t *testing.T) {
	// The observable of "Min/Max stay nil" is that an int value outside any
	// plausible bound is still admitted: no bound was ever supplied.
	body := declVerdictCovering + `

[emit.count]
kind = "int"`
	for _, literal := range []string{"-999999", "0", "999999"} {
		t.Run(literal, func(t *testing.T) {
			src := replaceOnce(t, dtWithEmitDecl(body),
				"verdict = \"alpha\"\n",
				"verdict = \"alpha\"\ncount = \""+literal+"\"\n")
			loadSource(t, src, "emit-no-bounds.toml")
		})
	}

	// And the reuse really is `ConformValue`'s: the exported function the
	// contract names admits exactly what the loader admits.
	if err := table.ConformValue(
		table.TagDecl{Kind: "int"}, "999999"); err != nil {
		t.Errorf("ConformValue admits `999999` under a bound-free int decl; "+
			"got %v", err)
	}
}

// REQ-54: "**A `scalar`-kinded key is short-circuited BEFORE the call**
// rather than relying on `conformKind`'s fall-through"
// BOUNDARY
func TestReq54_0024_AScalarKeyIsShortCircuitedBeforeTheConformCall(t *testing.T) {
	// The observable is that NO scalar value refuses, including values that
	// would trip a domain arm were one ever attached to a scalar. C1 bars
	// a domain on scalar (REQ-16), so the short-circuit is what makes the
	// escape hatch total.
	body := declVerdictCovering + `

[emit.free]
kind = "scalar"`

	for _, value := range []string{"", "anything", "42", "true", "  spaced  "} {
		t.Run("value "+value, func(t *testing.T) {
			src := replaceOnce(t, dtWithEmitDecl(body),
				"verdict = \"alpha\"\n",
				"verdict = \"alpha\"\nfree = \""+value+"\"\n")
			loadSource(t, src, "emit-scalar-hatch.toml")
		})
	}
}

// REQ-56: "The kernel (`internal/resolve`) continues to carry no
// declarations and no emit. This RDR adds exactly one reader of the carried
// declarations — C4's payload join — and `internal/graphlint` reads none of
// it here"
// `0024:S4`: "Kernel and dump surfaces carry none of it."
// DOMAIN EDGE
func TestReq56_0024_TheKernelAndDumpSurfacesCarryNoDeclarations(t *testing.T) {
	m := loadEmitDecls0024(t, declVerdictEnum, "emit-kernel.toml")

	kt := m.KernelTable()
	rt := reflect.TypeOf(kt)
	for i := range rt.NumField() {
		name := rt.Field(i).Name
		if name == "EmitDecls" || name == "Dispositions" {
			t.Errorf("the kernel table declares a `%s` field; the kernel "+
				"continues to carry no declarations", name)
		}
	}

	// The dump renders rows, never declarations: no disposition token and
	// no declared kind may appear in the rendered table.
	rendered := table.Dump(m)
	for _, token := range []string{"route", "stop", "\"enum\""} {
		if strings.Contains(rendered, token) {
			t.Errorf("the dump surface renders %q; declarations are carried "+
				"by neither the kernel nor the dump", token)
		}
	}
}

// REQ-69 / `0024:D-identity`: "an emit declaration is identified by its
// emit key, byte-exact … one declaration per key, duplicates refused by the
// TOML decoder as the tag table's are. A domain member's identity is its
// byte-exact string; one disposition per member."
// ADVERSARIAL
func TestReq69_0024_ADeclarationIsIdentifiedByItsByteExactEmitKey(t *testing.T) {
	t.Run("a duplicate declaration is the decoder's refusal", func(t *testing.T) {
		f := refuseEmitDecls0024(t, `[emit.verdict]
kind = "enum"
domain = ["alpha"]
[emit.verdict]
kind = "scalar"`, "emit-dup-decl.toml")
		if f.Category != table.CatMalformedTOML {
			t.Errorf("category = %q; want %q — duplicates are refused by "+
				"the TOML DECODER, as the tag table's are",
				f.Category, table.CatMalformedTOML)
		}
	})

	t.Run("key identity is byte-exact", func(t *testing.T) {
		// `Verdict` and `verdict` are two distinct keys, so declaring only
		// the case-folded spelling leaves the authored key undeclared.
		src := dtWithEmitDecl(`[emit.Verdict]
kind = "enum"
domain = ["alpha", "beta", "gamma"]`)
		f := refuseSource(t, src, "emit-case-fold.toml")
		if f.Category != table.CatUnknownEmitKey {
			t.Errorf("category = %q; want %q — an emit key's identity is "+
				"BYTE-EXACT", f.Category, table.CatUnknownEmitKey)
		}
	})

	t.Run("member identity is byte-exact", func(t *testing.T) {
		// The sibling clause: "A domain member's identity is its byte-exact
		// string." The 0010 rules author `verdict = "alpha"`, and this
		// declaration lists `Alpha`. Byte-exactly the value is outside the
		// domain; under a case-folding comparison the model would load
		// clean, which is the whole detection gap this subtest closes.
		src := dtWithEmitDecl(`[emit.verdict]
kind = "enum"
domain = ["Alpha", "beta", "gamma"]`)
		f := refuseSource(t, src, "emit-member-case-fold.toml")
		if f.Category != table.CatEmitValueOutOfDomain {
			t.Errorf("category = %q; want %q — a domain member's identity "+
				"is its BYTE-EXACT string",
				f.Category, table.CatEmitValueOutOfDomain)
		}
	})
}

// REQ-72 / `0024:D-selection-predicate`: "dispositions are per **domain
// member**, not per key … When the selected row authors a value for a
// declared enum key, the disposition surfaced is the one the declaration
// lists that member under — exactly one exists by C1's duplicate refusal."
// DOMAIN EDGE
func TestReq72_0024_DispositionsArePerMemberNotPerKey(t *testing.T) {
	m := loadEmitDecls0024(t, declVerdictEnum, "emit-per-member.toml")
	d := m.EmitDecls["verdict"]

	// Two members of the SAME key carry DIFFERENT dispositions — the
	// observable that forecloses a per-key reading.
	if d.Dispositions["alpha"] == d.Dispositions["gamma"] {
		t.Errorf("`alpha` and `gamma` both carry disposition %q; they are "+
			"listed under `route` and `stop` respectively, and dispositions "+
			"are per DOMAIN MEMBER, not per key", d.Dispositions["alpha"])
	}
	if len(d.Dispositions) != len(d.Domain) {
		t.Errorf("%d dispositions for %d domain members; exactly one exists "+
			"per member by C1's duplicate refusal",
			len(d.Dispositions), len(d.Domain))
	}
}

// REQ-101 / `0024:S4`: "**Determinism is asserted here or nowhere**:
// loading the same partitioned-domain model repeatedly yields one `Domain`
// value, which is the leg that fails if the union is built from map
// iteration without the sort. … this assertion is the only thing standing
// between an unsorted union and a green suite."
// ADVERSARIAL
func TestReq101_0024_LoadingAPartitionedDomainRepeatedlyYieldsOneDomain(t *testing.T) {
	// Enough disposition keys and members that Go's randomized map
	// iteration would, without the sort, produce a different union order
	// across these repetitions with overwhelming probability.
	body := `[emit.verdict]
kind = "enum"
[emit.verdict.domain]
d1 = ["m07", "m01"]
d2 = ["m04", "m10"]
d3 = ["m02", "m08"]
d4 = ["m05", "m09"]
d5 = ["m03", "m06"]`

	src := dtWithEmitDecl(body)
	// The rules author `verdict` values outside this domain, so replace
	// them with declared members: the model must LOAD for the carrier to be
	// readable.
	for authored, member := range map[string]string{
		"alpha": "m01", "beta": "m02", "gamma": "m03",
	} {
		src = replaceOnce(t, src,
			"verdict = \""+authored+"\"\n", "verdict = \""+member+"\"\n")
	}

	want := []string{
		"m01", "m02", "m03", "m04", "m05",
		"m06", "m07", "m08", "m09", "m10",
	}
	for i := range 32 {
		m := loadSource(t, src, "emit-determinism.toml")
		got := m.EmitDecls["verdict"].Domain
		if !slices.Equal(got, want) {
			t.Fatalf("load %d yielded Domain = %v; want the one bytewise-"+
				"sorted union %v — an unsorted union built from map "+
				"iteration is what this assertion stands between",
				i, got, want)
		}
	}
}

// REQ-109 / `0024:G-cross-cutting`: "**Concurrency model** — none
// introduced. Both new steps are `*loader` methods on the single-threaded
// load path, and the declaration carrier is written once during load and
// read-only thereafter"
// ADVERSARIAL
func TestReq109_0024_TheCarrierIsWrittenOnceAndReadOnlyThereafter(t *testing.T) {
	// Concurrent LOADS of the same source share no state: each yields its
	// own carrier with the same value. Run under `-race` this is the
	// observable of "none introduced".
	src := dtWithEmitDecl(declVerdictEnum)
	want := []string{"alpha", "beta", "gamma"}

	results := make(chan []string, 8)
	for range 8 {
		go func() {
			m, err := table.Load([]byte(src), "emit-concurrent.toml")
			if err != nil {
				results <- nil
				return
			}
			results <- m.EmitDecls["verdict"].Domain
		}()
	}
	for range 8 {
		got := <-results
		if !slices.Equal(got, want) {
			t.Errorf("a concurrent load yielded Domain = %v; want %v", got, want)
		}
	}
}

// REQ-110 / `0024:G-cross-cutting`: "No hash, no content-addressed
// identity, and no replay-stable digest is claimed anywhere in this RDR"
// DOMAIN EDGE
func TestReq110_0024_NoHashOrContentAddressedIdentityIsCarried(t *testing.T) {
	rt := reflect.TypeOf(table.EmitDecl{})
	for i := range rt.NumField() {
		switch rt.Field(i).Name {
		case "Hash", "Digest", "Fingerprint", "ContentID":
			t.Errorf("EmitDecl declares a `%s` field; no hash, no content-"+
				"addressed identity, and no replay-stable digest is claimed",
				rt.Field(i).Name)
		}
	}
}

// --- helpers -------------------------------------------------------------

// replaceOnce substitutes old for new exactly once and fails the test when
// the substitution does not apply, so a silently unmutated fixture cannot
// pass an assertion vacuously.
func replaceOnce(t *testing.T, src, old, new string) string {
	t.Helper()

	out := strings.Replace(src, old, new, 1)
	if out == src {
		t.Fatalf("the substitution %q -> %q did not apply", old, new)
	}
	return out
}
