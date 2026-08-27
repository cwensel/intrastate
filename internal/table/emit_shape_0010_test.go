package table_test

// RDR 0010 — the remaining C3 shape clauses and the SC-3 scenario as one
// table, plus the negative REQs whose only assertable form is "the
// rejected machinery is absent".

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/table"
)

// REQ-21: "The assertion is on the category, never on upstream message
// text (`0002:C24`), and no hand-written type check is added."
// DOMAIN EDGE — negative REQ.
//
// The assertable form: the refusal for a wrong-typed emit value is the
// DECODER's category, and it is not distinguishable from any other decoder
// type error — which a hand-written check would make it.
func TestReq21_NoHandWrittenTypeCheckDiscriminatesTheEmitRefusal(t *testing.T) {
	emitTyped := refuseSource(t, strings.Replace(dtComplete,
		"[rule.emit]\nverdict = \"alpha\"\n",
		"[rule.emit]\nverdict = 42\n", 1), "emit-typed.toml")

	// An unrelated wrong-typed scalar the decoder also refuses.
	otherTyped := refuseSource(t, strings.Replace(dtComplete,
		`id = "cell-xp"`, `id = 42`, 1), "other-typed.toml")

	if emitTyped.Category != otherTyped.Category {
		t.Errorf("a wrong-typed emit value refuses %q while another "+
			"wrong-typed scalar refuses %q; the emit arm is a hand-written "+
			"type check rather than the decoder's",
			emitTyped.Category, otherTyped.Category)
	}
	if emitTyped.Category != table.CatMalformedTOML {
		t.Errorf("category = %q; want %q",
			emitTyped.Category, table.CatMalformedTOML)
	}
}

// REQ-25: "The carried sequence MAY be **shared** across the rows one rule
// expands to — no clone is required."
// DOMAIN EDGE — negative REQ: do not add a defensive copy.
//
// Sharing and cloning are indistinguishable to a byte-for-byte assertion,
// so the assertable form is that the LICENCE holds: every expanded row's
// emit is EQUAL, and nothing downstream depends on them being distinct
// backing arrays.
func TestReq25_SharingTheEmitSequenceAcrossExpandedRowsIsLicensed(t *testing.T) {
	m := loadSource(t, dtExpanding, "dt-expanding.toml")

	rows := rowsByRuleID(m, "spread")
	if len(rows) < 2 {
		t.Fatalf("the `in`-atom rule expanded to %d rows; want at least 2",
			len(rows))
	}
	for i := 1; i < len(rows); i++ {
		if !slices.Equal(rows[0].Emit, rows[i].Emit) {
			t.Errorf("%s carries emit %v while %s carries %v; every row a "+
				"rule expands to carries the SAME block, shared or cloned",
				rows[0].Identity(), rows[0].Emit,
				rows[i].Identity(), rows[i].Emit)
		}
	}

	// Nothing downstream mutates it, so sharing is safe — the dump is a
	// reader (asserted in TestReq26) and the resolvePayload join is a
	// reader. Re-rendering must not make the rows diverge.
	_ = table.Dump(m)
	for i := 1; i < len(rows); i++ {
		if !slices.Equal(rows[0].Emit, rows[i].Emit) {
			t.Errorf("rendering made expanded rows' emit diverge: %v vs %v",
				rows[0].Emit, rows[i].Emit)
		}
	}
}

// REQ-31: "An empty block renders as the empty bracket the `writes` column
// already uses for an empty sequence."
// BOUNDARY — the SAME rendering, not a lookalike.
func TestReq31_AnEmptyEmitRendersTheSameEmptyBracketWritesUses(t *testing.T) {
	m := loadSource(t, dtEmptyEmit, "dt-empty-emit.toml")
	line := dumpLineFor(t, table.Dump(m), "dt.cell-xp")

	fields := strings.Fields(line)
	var writesCell, emitCell string
	for _, f := range fields {
		if after, ok := strings.CutPrefix(f, "writes="); ok {
			writesCell = after
		}
		if after, ok := strings.CutPrefix(f, "emit="); ok {
			emitCell = after
		}
	}
	if writesCell == "" || emitCell == "" {
		t.Fatalf("dump line carries no writes/emit cell:\n%s", line)
	}
	if emitCell != writesCell {
		t.Errorf("an empty emit renders %q while an empty writes renders "+
			"%q; C3 fixes them as the SAME empty bracket",
			emitCell, writesCell)
	}
}

// REQ-91 / SC-3 (`0010:S3`): the scenario as one table — "unordered keys, a
// non-string value, a nested table (`[rule.emit.sub]`), an absent block, a
// **present but empty** block; dump with and without an explicit `[dump]`
// naming `emit`; `[rule.match.<emit-key>]`"
// HAPPY PATH
func TestReq91_TheEmitNormalizationScenarioAsOneTable(t *testing.T) {
	cases := []struct {
		name string
		src  string
		// want is the expected refusal category; empty means "must load".
		want table.Category
	}{
		{name: "unordered keys", src: dtUnordered},
		{name: "absent block", src: dtComplete},
		{name: "present but empty block", src: dtEmptyEmit},
		{
			name: "non-string value",
			src: strings.Replace(dtComplete,
				"[rule.emit]\nverdict = \"alpha\"\n",
				"[rule.emit]\nverdict = 42\n", 1),
			want: table.CatMalformedTOML,
		},
		{
			name: "nested table",
			src: strings.Replace(dtComplete,
				"[rule.emit]\nverdict = \"alpha\"\n",
				"[rule.emit]\nverdict = \"alpha\"\n[rule.emit.sub]\ninner = \"x\"\n", 1),
			want: table.CatMalformedTOML,
		},
		{
			name: "[dump] omitting emit",
			src: dtComplete + `
[dump]
order = [
  "identity", "source", "kind", "outcome",
  "atoms", "next", "writes", "requires_owned", "gate", "escape",
]
`,
			want: table.CatMalformedDumpDeclaration,
		},
		{
			name: "[dump] naming emit",
			src: dtComplete + `
[dump]
order = [
  "identity", "source", "kind", "outcome",
  "atoms", "next", "writes", "requires_owned", "gate", "escape", "emit",
]
`,
		},
		{
			name: "[rule.match.<emit-key>]",
			src: strings.Replace(dtComplete,
				"[rule.guard.all.a]\neq = \"x\"\n[rule.guard.all.b]\neq = \"p\"\n[rule.emit]\nverdict = \"alpha\"\n",
				"[rule.match.verdict]\neq = \"alpha\"\n[rule.guard.all.a]\neq = \"x\"\n[rule.guard.all.b]\neq = \"p\"\n[rule.emit]\nverdict = \"alpha\"\n",
				1),
			want: table.CatUnknownTag,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.want == "" {
				loadSource(t, c.src, c.name+".toml")
				return
			}
			f := refuseSource(t, c.src, c.name+".toml")
			if f.Category != c.want {
				t.Errorf("category = %q; want %q", f.Category, c.want)
			}
		})
	}
}

// REQ-11 (loader arm): the doubly-malformed determinacy case, taken on the
// LOADER as C1 states it — a decision-table model declaring BOTH an owned
// tag and `[initial]` refuses naming the class, never with the
// accessor-binding writer-arity diagnostic.
// BOUNDARY
//
// Companion to TestReq10 in class_0010_test.go: this arm gives the owned
// tag its reader AND writer, so the arity check would PASS and the only
// refusal left is the class disagreement. Without a writer TestReq10's
// document could refuse on arity for the wrong reason.
func TestReq11_TheClassRefusalWinsEvenWhenTheArityCheckWouldPass(t *testing.T) {
	src := strings.Replace(dtComplete, "[tags.a]", `[initial]
own = "seed"

[tags.own]
provenance = "owned"
kind = "enum"
domain = ["seed"]
single_valued = true
required = true

[read.own]
role = "own"
path = "o"
keys = ["own"]
timeout = "2s"

[write.own]
role = "own"
path = "o"
keys = ["own"]
timeout = "2s"
read_back = true

[tags.a]`, 1)

	f := refuseSource(t, src, "dt-owned-initial-bound.toml")
	if f.Category != table.CatMalformedModelDeclaration {
		t.Fatalf("category = %q; want %q — with the arity check satisfied, "+
			"the class disagreement is the only refusal left, and it must "+
			"fire", f.Category, table.CatMalformedModelDeclaration)
	}
	if !strings.Contains(f.Detail, "owned=1") {
		t.Errorf("detail = %q; want the `owned=1` token", f.Detail)
	}
}

// REQ-18: "a terminal predicate over a non-owned tag is 0006's
// `graph-dangling-edge` terminal arm, which C5 keeps live for exactly this
// reason"
// DOMAIN EDGE — the `terminal` prohibition is enforced at LINT, not by a
// load refusal.
func TestReq18_ATerminalOverANonOwnedTagIsNotALoadRefusal(t *testing.T) {
	src := "terminal = [\"stop\"]\n" + dtComplete + `
[context.stop]
[context.stop.match.a]
eq = "x"
`

	m := loadSource(t, src, "dt-stray-terminal.toml")
	if len(m.Terminal) != 1 {
		t.Errorf("Terminal = %v; the declaration loads and is left to lint — "+
			"C2's prohibition is carried by 0006's terminal arm, not by a "+
			"new load refusal", m.Terminal)
	}
}

// REQ-13: "The class is **declared, not inferred** from the empty owned
// set." (restating REQ-6 at the design level)
// ADVERSARIAL — negative REQ, taken on the SIBLING discriminator the
// Approach names.
//
// `Row.Kind` IS inferred from presence (`len(r.Escape) > 0`) and is
// per-row. The class must be neither: it is per-MODEL and declared.
func TestReq13_TheClassIsPerModelAndDeclaredUnlikeRowKind(t *testing.T) {
	m := loadSource(t, dtEscapeVariant, "dt-escape.toml")

	// Row.Kind is inferred and per-row: this model carries both kinds.
	var kinds []table.Kind
	for _, r := range m.Rows {
		if !slices.Contains(kinds, r.Kind()) {
			kinds = append(kinds, r.Kind())
		}
	}
	if len(kinds) != 2 {
		t.Fatalf("the fixture carries row kinds %v; it must carry both so "+
			"the contrast with the per-model class is real", kinds)
	}

	// The class is one value for the whole model, and it is what was
	// DECLARED — not something varying with the rows.
	if !table.IsDecisionTable(m) {
		t.Error("the model does not read as a decision table despite " +
			"declaring the class")
	}
	rt := reflect.TypeOf(table.Row{})
	if _, ok := rt.FieldByName("Class"); ok {
		t.Error("table.Row carries a `Class` field; the class is per-MODEL, " +
			"never per-row")
	}
}

// REQ-110 / IC: "Illustrative — a two-dimension decision table and the
// resolve envelope it yields; tests must not assert these literally."
// DOMAIN EDGE — negative REQ.
//
// Discharged by construction: no test in this suite reads the record's
// Illustrative Code or its mermaid figure. This test states that as an
// executable claim about the fixture corpus — every fixture is authored
// here, and its identifiers are this suite's own.
func TestReq110_NoFixtureIsLiftedFromTheIllustrativeCode(t *testing.T) {
	// The RDR's illustrative table names its dimensions and outcome
	// differently from this suite's fixtures. If a future edit copied the
	// IC verbatim, these identifiers would appear.
	m := loadSource(t, dtComplete, "dt-complete.toml")

	if m.ID != "dt" {
		t.Errorf("model id = %q; the fixture corpus is authored by this "+
			"suite, not lifted from the record's non-normative Illustrative "+
			"Code", m.ID)
	}
	for _, r := range m.Rows {
		if !strings.HasPrefix(r.RuleID, "cell-") {
			t.Errorf("rule id %q does not follow this suite's own naming; "+
				"the Illustrative Code is not normative and must not be "+
				"asserted literally", r.RuleID)
		}
	}
}
