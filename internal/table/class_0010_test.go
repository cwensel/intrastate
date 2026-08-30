package table_test

// RDR 0010 §A — the model class declaration (`0010:C1`) and its
// one-directional agreement check.
//
// Every oracle here goes through the real `table.Load` and asserts on the
// stable load CATEGORY, never on message text (`0002:C24`) — with the one
// exception C1 itself makes normative: the disagreement detail must carry
// the LITERAL token `owned=<n>`, which the fence names verbatim.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// REQ-1: "`[model]` MAY carry `class`, a string whose only admitted values
// are `\"state-machine\"` and `\"decision-table\"`"
// HAPPY PATH
func TestReq1_ClassAdmitsExactlyTwoValues(t *testing.T) {
	t.Run("decision-table", func(t *testing.T) {
		m := loadSource(t, dtComplete, "dt-complete.toml")
		if m.Class != "decision-table" {
			t.Errorf("Class = %q; want %q", m.Class, "decision-table")
		}
	})

	t.Run("state-machine", func(t *testing.T) {
		m := loadSource(t, smZeroOwned, "sm-zero-owned.toml")
		if m.Class != "state-machine" {
			t.Errorf("Class = %q; want %q", m.Class, "state-machine")
		}
	})
}

// REQ-2: "an absent `class` MUST be read as `\"state-machine\"`"
// INPUT EDGE
//
// The discriminating half is BEHAVIOURAL, not just the field value: a model
// with `class` omitted must still take the state-machine arm of every
// class-keyed decision. `dtClassOmitted` declares zero owned tags, so a
// loader that read the absent key as "decision-table" would load it as one.
func TestReq2_AbsentClassReadsAsStateMachine(t *testing.T) {
	m := loadSource(t, dtClassOmitted, "dt-class-omitted.toml")

	if m.Class != "" && m.Class != "state-machine" {
		t.Errorf("Class = %q; an absent `class` reads as %q (or its zero value)",
			m.Class, "state-machine")
	}
	if table.IsDecisionTable(m) {
		t.Error("a model with `class` omitted reads as a decision table; " +
			"an absent `class` MUST be read as \"state-machine\"")
	}
}

// REQ-3: "A `class` outside that set is `malformed model declaration`."
// ADVERSARIAL
func TestReq3_UnknownClassIsMalformedModelDeclaration(t *testing.T) {
	for _, bad := range []string{"lookup-table", "State-Machine", "decision_table", ""} {
		t.Run("class="+bad, func(t *testing.T) {
			src := strings.Replace(dtComplete,
				`class = "decision-table"`, `class = "`+bad+`"`, 1)
			f := refuseSource(t, src, "bad-class.toml")
			if f.Category != table.CatMalformedModelDeclaration {
				t.Errorf("category = %q; want %q",
					f.Category, table.CatMalformedModelDeclaration)
			}
		})
	}
}

// REQ-4: "The agreement check is **one-directional**: a `\"decision-table\"`
// model MUST declare zero tags of provenance `owned`, and a non-zero owned
// count under that class is a `malformed model declaration` load failure
// whose detail names the class and the owned count, rendered as the token
// `owned=<n>`."
// ADVERSARIAL
//
// The detail assertion is on the LITERAL token the fence names, not on an
// invented substring (SC-1).
func TestReq4_DecisionTableWithOwnedTagRefusesNamingClassAndOwnedCount(t *testing.T) {
	f := refuseSource(t, dtWithOwnedTag, "dt-owned-tag.toml")

	if f.Category != table.CatMalformedModelDeclaration {
		t.Fatalf("category = %q; want %q",
			f.Category, table.CatMalformedModelDeclaration)
	}
	if !strings.Contains(f.Detail, "owned=1") {
		t.Errorf("detail = %q; it must render the owned count as the literal "+
			"token `owned=1`", f.Detail)
	}
	if !strings.Contains(f.Detail, "decision-table") {
		t.Errorf("detail = %q; it must name the declared class", f.Detail)
	}
}

// REQ-4 (count arm): the rendered count is the ACTUAL owned count, not a
// hard-coded 1.
// BOUNDARY
func TestReq4_OwnedTokenCarriesTheActualCount(t *testing.T) {
	src := strings.Replace(dtComplete, "[tags.a]", `[tags.own1]
provenance = "owned"
kind = "scalar"

[tags.own2]
provenance = "owned"
kind = "scalar"

[tags.a]`, 1)

	f := refuseSource(t, src, "dt-two-owned.toml")
	if f.Category != table.CatMalformedModelDeclaration {
		t.Fatalf("category = %q; want %q",
			f.Category, table.CatMalformedModelDeclaration)
	}
	if !strings.Contains(f.Detail, "owned=2") {
		t.Errorf("detail = %q; two owned tags must render `owned=2`", f.Detail)
	}
}

// REQ-5: "A `\"state-machine\"` model MUST NOT be refused for declaring
// zero owned tags — that model is rootless, and `0006:C18` reports it as
// `graph-dangling-edge` at lint, which is the pre-existing behaviour this
// RDR leaves untouched"
// DOMAIN EDGE — negative REQ: no second direction of the agreement check.
//
// GREEN BY CONSTRUCTION against main: nothing refuses this today. It is a
// characterization test pinning that the implementation does not GROW a
// second arm.
func TestReq5_ZeroOwnedStateMachineIsNotRefused(t *testing.T) {
	t.Run("class declared explicitly", func(t *testing.T) {
		m := loadSource(t, smZeroOwned, "sm-zero-owned.toml")
		if len(m.Rows) == 0 {
			t.Fatal("normalized to no candidate rows")
		}
	})

	t.Run("class omitted", func(t *testing.T) {
		src := strings.Replace(smZeroOwned, "class = \"state-machine\"\n", "", 1)
		loadSource(t, src, "sm-zero-owned-implicit.toml")
	})
}

// REQ-6 / REQ-13: "The class is model data and is carried on the loaded
// model; nothing downstream infers it from the owned set." / "The class is
// **declared, not inferred** from the empty owned set."
// ADVERSARIAL — negative REQ.
//
// The discriminating pair: two documents with the SAME (empty) owned set
// and DIFFERENT declared classes must carry different classes. An
// implementation inferring the class from `len(owned) == 0` gives both the
// same answer and fails here.
func TestReq6_ClassIsDeclaredNeverInferredFromTheOwnedSet(t *testing.T) {
	declared := loadSource(t, dtComplete, "dt-complete.toml")
	// The class-omitted half is `smZeroOwned` with its `class` line
	// stripped, NOT `dtClassOmitted`. deviations.md D1 gives
	// `dtClassOmitted` owned-tag scaffolding — a zero-owned model's
	// ordinary rule cannot carry the write block `0002:C4` still demands
	// of a state machine — so that document does not share the DECLARED
	// half's empty owned set and the pair below could not discriminate an
	// implementation inferring the class from `len(owned) == 0`. The
	// zero-owned escape-row document does share it, which makes the
	// invariant real rather than merely asserted.
	omitted := loadSource(t,
		strings.Replace(smZeroOwned, "class = \"state-machine\"\n", "", 1),
		"sm-zero-owned-implicit.toml")

	if !table.IsDecisionTable(declared) {
		t.Error("a model declaring `class = \"decision-table\"` does not read " +
			"as a decision table")
	}
	if table.IsDecisionTable(omitted) {
		t.Error("a zero-owned model with `class` OMITTED reads as a decision " +
			"table; the class is declared, never inferred from the owned set")
	}

	// Both carry the same owned set — the inference an implementation must
	// not make.
	for _, m := range []*table.Model{declared, omitted} {
		for key, decl := range m.Tags {
			if decl.Provenance == table.ProvenanceOwned {
				t.Fatalf("fixture invariant broken: %s declares the owned tag %q",
					m.ID, key)
			}
		}
	}
}

// REQ-7: "The class is carried as an **exported `Class string` field on
// `internal/table/model.go::Model`**, whose zero value — the empty string —
// MUST be read everywhere as `\"state-machine\"`, so a hand-constructed
// `table.Model` is a state machine without naming a class."
// BOUNDARY — an exported field, not an accessor.
func TestReq7_ClassIsAnExportedStringFieldWhoseZeroValueIsStateMachine(t *testing.T) {
	rt := reflect.TypeOf(table.Model{})
	f, ok := rt.FieldByName("Class")
	if !ok {
		t.Fatalf("table.Model carries no exported `Class` field; fields = %v",
			fieldNames(rt))
	}
	if f.Type.Kind() != reflect.String {
		t.Errorf("Model.Class is %s; C1 fixes an exported `Class string` field",
			f.Type)
	}

	// The zero value: a hand-constructed model that never passed the loader
	// is a STATE MACHINE.
	var zero table.Model
	if zero.Class != "" {
		t.Errorf("the zero value of Model.Class = %q; want the empty string",
			zero.Class)
	}
	if table.IsDecisionTable(&zero) {
		t.Error("a hand-constructed table.Model reads as a decision table; " +
			"the zero value MUST be read as \"state-machine\"")
	}
}

// REQ-8: "A reader MAY spell the zero value through a small helper rather
// than repeating the empty-string comparison, but the field is the storage
// and no reader re-derives the class from the owned set."
// BOUNDARY
//
// The FIELD is the storage: mutating it must change the answer, which an
// implementation caching the class elsewhere (or deriving it) fails.
func TestReq8_TheFieldIsTheStorageForTheClass(t *testing.T) {
	m := loadSource(t, dtComplete, "dt-complete.toml")
	if !table.IsDecisionTable(m) {
		t.Fatal("loaded decision table does not read as one")
	}

	m.Class = "state-machine"
	if table.IsDecisionTable(m) {
		t.Error("setting Model.Class = \"state-machine\" did not change the " +
			"answer; the field is the storage, not a cache of a derivation")
	}

	m.Class = ""
	if table.IsDecisionTable(m) {
		t.Error("the empty Class does not read as \"state-machine\"")
	}
}

// REQ-9: "The check MUST run after the tag table is loaded, not in the
// `[model]` header step" … "The class is read in `loadModelHeader` (it is
// `[model]` data) and the agreement is checked in a step at or after
// `loadTags`"
// DOMAIN EDGE
//
// Observable consequence under `run`'s fail-fast order: an UNDECLARED-tag
// refusal precedes a class-disagreement refusal. A check wired into
// `loadModelHeader` would report the class disagreement first.
func TestReq9_AnUndeclaredTagRefusalPrecedesTheClassDisagreement(t *testing.T) {
	// A decision table that BOTH declares an owned tag and names an
	// undeclared tag in a rule guard.
	src := strings.Replace(dtWithOwnedTag,
		"[rule.guard.all.a]\neq = \"x\"\n[rule.guard.all.b]\neq = \"p\"\n[rule.emit]\nverdict = \"alpha\"\n",
		"[rule.guard.all.nope]\neq = \"x\"\n[rule.guard.all.b]\neq = \"p\"\n[rule.emit]\nverdict = \"alpha\"\n",
		1)
	if src == dtWithOwnedTag {
		t.Fatal("the undeclared-tag substitution did not apply")
	}

	f := refuseSource(t, src, "dt-owned-and-unknown-tag.toml")
	if f.Category != table.CatUnknownTag {
		t.Errorf("category = %q; want %q — C1 draws this consequence from "+
			"the window's floor verbatim: \"an undeclared-tag refusal "+
			"precedes a class-disagreement refusal under `run`'s fail-fast "+
			"order\"", f.Category, table.CatUnknownTag)
	}
}

// REQ-10 / REQ-11: "The position is fixed at **both** ends: the check MUST
// run at or after `loadTags` and **before `checkAccessorBindings`**, `run`'s
// last step." / "a decision-table model that declares both an owned tag and
// `[initial]` refuses as `malformed model declaration` naming the class and
// `owned=<n>`, never with the accessor-binding writer-arity diagnostic"
// ADVERSARIAL — the doubly-malformed determinacy case.
func TestReq10_DoublyMalformedRefusesOnTheClassNotOnWriterArity(t *testing.T) {
	src := strings.Replace(dtWithOwnedTag, "[tags.own]", `[initial]
own = "seed"

[tags.own]`, 1)
	if src == dtWithOwnedTag {
		t.Fatal("the [initial] substitution did not apply")
	}

	f := refuseSource(t, src, "dt-owned-and-initial.toml")
	if f.Category != table.CatMalformedModelDeclaration {
		t.Fatalf("category = %q; want %q — the agreement check runs BEFORE "+
			"checkAccessorBindings, so the writer-arity diagnostic must not "+
			"be what a doubly-malformed decision table sees",
			f.Category, table.CatMalformedModelDeclaration)
	}
	if !strings.Contains(f.Detail, "owned=1") {
		t.Errorf("detail = %q; the class refusal must carry the `owned=<n>` "+
			"token rather than the accessor-binding diagnostic", f.Detail)
	}
	if f.Category == table.CatMalformedAccessorBinding {
		t.Error("the refusal is the accessor-binding writer-arity diagnostic")
	}
}

// REQ-12: "Any step in that window satisfies this clause; the RDR fixes the
// window, not the step's name."
// BOUNDARY
//
// The WINDOW is asserted structurally against `run`'s step slice: whatever
// the step is called, it appears at or after `loadTags` and before
// `checkAccessorBindings`. Asserting the name would over-fix what the fence
// deliberately leaves free (ASSUMPTION-1).
func TestReq12_TheAgreementCheckSitsInsideTheLicensedWindow(t *testing.T) {
	steps := loaderRunSteps(t)

	tagsAt := indexOfStep(steps, "loadTags")
	bindingsAt := indexOfStep(steps, "checkAccessorBindings")
	if tagsAt < 0 || bindingsAt < 0 {
		t.Fatalf("run's step slice does not name loadTags/checkAccessorBindings: %v",
			steps)
	}
	if bindingsAt != len(steps)-1 {
		t.Errorf("checkAccessorBindings is at %d of %d; C1 calls it run's last step",
			bindingsAt, len(steps))
	}

	// The class check is whichever step is NEW relative to the steps this
	// RDR did not add. Rather than guess a name, assert the loader grew
	// exactly one step beyond that baseline and that it landed inside the
	// window.
	//
	// RDR 0024 `0024:C2` normatively inserts `loadEmitDecls` and
	// `checkRuleEmit` immediately after `loadTags`, so they belong to the
	// baseline here for the same reason 0002's ten do: they are not the
	// step this test is identifying. Leaving them out would make the
	// technique — "the added step is the one that is new" — fail on any
	// successor that adds a loader step, which is a property of the
	// technique, not of `0010:C1`'s window.
	base := []string{
		"loadModelHeader", "loadOutcomes", "loadTags", "loadEmitDecls",
		"checkRuleEmit", "loadAccessors", "loadDump", "loadContexts",
		"loadInitial", "loadTerminal", "normalizeRules",
		"checkAccessorBindings",
	}
	var added []int
	for i, s := range steps {
		if !contains(base, s) {
			added = append(added, i)
		}
	}
	if len(added) != 1 {
		t.Fatalf("run's steps = %v; expected exactly one step added to 0002's "+
			"fixed order %v (the class/owned-set agreement check)", steps, base)
	}
	at := added[0]
	if at <= tagsAt || at >= bindingsAt {
		t.Errorf("the added step sits at index %d; C1 fixes the window at or "+
			"after loadTags (%d) and before checkAccessorBindings (%d)",
			at, tagsAt, bindingsAt)
	}
}

// REQ-89 / SC-1 (`0010:S1`): "loader table tests over `class` — absent, each
// admitted value, an unknown value, the one disagreement direction
// (`decision-table` declaring an owned tag), and its counterpart control —
// a `state-machine` model (or one with `class` omitted) declaring zero owned
// tags, which MUST load."
// HAPPY PATH — the scenario as a single table.
func TestReq89_LoaderTableOverTheClassKey(t *testing.T) {
	cases := []struct {
		name string
		src  string
		// want is the expected category; empty means "must load clean".
		want  table.Category
		class string
	}{
		{name: "absent", src: dtClassOmitted, class: "state-machine"},
		{name: "decision-table", src: dtComplete, class: "decision-table"},
		{name: "state-machine", src: smZeroOwned, class: "state-machine"},
		{
			name: "unknown value",
			src: strings.Replace(dtComplete,
				`class = "decision-table"`, `class = "table"`, 1),
			want: table.CatMalformedModelDeclaration,
		},
		{
			name: "decision-table declaring an owned tag",
			src:  dtWithOwnedTag,
			want: table.CatMalformedModelDeclaration,
		},
		{name: "zero-owned state machine loads", src: smZeroOwned, class: "state-machine"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.want == "" {
				m := loadSource(t, c.src, c.name+".toml")
				got := m.Class
				if got == "" {
					got = "state-machine"
				}
				if got != c.class {
					t.Errorf("Class = %q; want %q", got, c.class)
				}
				return
			}
			f := refuseSource(t, c.src, c.name+".toml")
			if f.Category != c.want {
				t.Errorf("category = %q; want %q", f.Category, c.want)
			}
			if c.name == "decision-table declaring an owned tag" &&
				!strings.Contains(f.Detail, "owned=") {
				t.Errorf("detail = %q; the disagreement detail must contain "+
					"the literal token `owned=<n>`", f.Detail)
			}
		})
	}
}

// --- structural helpers --------------------------------------------------

// loaderRunSteps parses `internal/table/load.go` and returns the method
// names in `run`'s step slice, in order.
//
// This is a STRUCTURAL oracle for the ONE clause C1 states about position
// rather than behaviour (REQ-10, REQ-12). Every other assertion in this
// file is behavioural.
func loaderRunSteps(t *testing.T) []string {
	t.Helper()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "load.go", nil, 0)
	if err != nil {
		t.Fatalf("parse load.go: %v", err)
	}

	var steps []string
	ast.Inspect(f, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "run" {
			return true
		}
		ast.Inspect(fn.Body, func(inner ast.Node) bool {
			lit, ok := inner.(*ast.CompositeLit)
			if !ok {
				return true
			}
			for _, elt := range lit.Elts {
				sel, ok := elt.(*ast.SelectorExpr)
				if !ok {
					continue
				}
				steps = append(steps, sel.Sel.Name)
			}
			return true
		})
		return false
	})
	if len(steps) == 0 {
		t.Fatal("could not read run's step slice from load.go")
	}
	return steps
}

func indexOfStep(steps []string, want string) int {
	for i, s := range steps {
		if s == want {
			return i
		}
	}
	return -1
}

func contains(hay []string, want string) bool {
	return indexOfStep(hay, want) >= 0
}

func fieldNames(rt reflect.Type) []string {
	out := make([]string, 0, rt.NumField())
	for i := range rt.NumField() {
		out = append(out, rt.Field(i).Name)
	}
	return out
}
