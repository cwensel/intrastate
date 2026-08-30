package table_test

// RDR 0010 §B–§C — the class-conditioned rule shape (`0010:C2`) and the
// `[rule.emit]` grammar, normalization, and dump column (`0010:C3`).

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// REQ-16: "`0002:C4`'s \"an ordinary transition rule MUST contain a write
// block\" is conditioned on class: it binds the `\"state-machine\"` class
// only, and the `malformed rule shape` arm that enforces it MUST NOT fire
// for a `\"decision-table\"` model."
// HAPPY PATH — the one edit C2 requires.
func TestReq16_TheNoWriteBlockArmIsConditionedOnClass(t *testing.T) {
	t.Run("state-machine still refuses", func(t *testing.T) {
		f := refuseSource(t, smNoWriteBlock, "sm-no-write.toml")
		if f.Category != table.CatMalformedRuleShape {
			t.Errorf("category = %q; want %q — the arm binds the "+
				"state-machine class and is unchanged there",
				f.Category, table.CatMalformedRuleShape)
		}
	})

	t.Run("decision-table loads", func(t *testing.T) {
		m := loadSource(t, dtComplete, "dt-complete.toml")
		if len(m.Rows) != 4 {
			t.Fatalf("normalized to %d rows; want 4", len(m.Rows))
		}
	})
}

// REQ-90 / SC-2 (`0010:S2`): "rule-shape normalization over an ordinary
// rule with no write block, in each class." Expected: "`malformed rule
// shape` for `state-machine`; loads with empty `Writes`/`NextTags`/
// `RequiresOwned` for `decision-table` (C2)."
// HAPPY PATH
func TestReq90_DecisionTableRowsNormalizeWithEmptyWriteSurface(t *testing.T) {
	m := loadSource(t, dtComplete, "dt-complete.toml")

	for _, r := range m.Rows {
		if len(r.Writes) != 0 {
			t.Errorf("%s: Writes = %v; want empty", r.Identity(), r.Writes)
		}
		if len(r.NextTags) != 0 {
			t.Errorf("%s: NextTags = %v; want empty", r.Identity(), r.NextTags)
		}
		if len(r.RequiresOwned) != 0 {
			t.Errorf("%s: RequiresOwned = %v; want empty",
				r.Identity(), r.RequiresOwned)
		}
	}
}

// REQ-14: "A decision-table model MUST NOT declare `[initial]`, `terminal`,
// or any accessor whose `keys` name an owned tag, and its ordinary rules
// MUST NOT carry a write block or a clear list."
// REQ-15: "None of these is a new refusal: each is already unauthorable
// once the owned set is empty"
// ADVERSARIAL — negative REQ: the EXISTING arms carry the prohibition.
//
// The assertion is on WHICH arm fires: a new dedicated `decision table
// declares [initial]` category would fail here.
func TestReq15_TheProhibitionsAreCarriedByExistingRefusalArms(t *testing.T) {
	t.Run("[initial] over an undeclared key is unknown tag", func(t *testing.T) {
		src := strings.Replace(dtComplete, "[tags.recognized]", `[initial]
ghost = "seed"

[tags.recognized]`, 1)
		f := refuseSource(t, src, "dt-initial-undeclared.toml")
		if f.Category != table.CatUnknownTag &&
			f.Category != table.CatMalformedInitialDeclaration {
			t.Errorf("category = %q; an `[initial]` key that is not a declared "+
				"tag refuses through 0002's existing arms, not a new one",
				f.Category)
		}
	})

	t.Run("[initial] over a declared observed key refuses at binding", func(t *testing.T) {
		src := strings.Replace(dtComplete, "[tags.recognized]", `[initial]
a = "x"

[tags.recognized]`, 1)
		f := refuseSource(t, src, "dt-initial-observed.toml")
		if f.Category != table.CatMalformedAccessorBinding {
			t.Errorf("category = %q; want %q — a declared observed key in "+
				"`[initial]` is a written tag and the writer-arity check "+
				"refuses it (C2)", f.Category, table.CatMalformedAccessorBinding)
		}
	})

	t.Run("a writer naming a non-owned tag refuses", func(t *testing.T) {
		src := dtComplete + `
[write.state]
role = "state"
path = "s.own"
keys = ["a"]
timeout = "2s"
read_back = true
`
		f := refuseSource(t, src, "dt-writer-observed.toml")
		if f.Category != table.CatWriteToNonOwnedTag {
			t.Errorf("category = %q; want %q — `loadAccessors` refuses any "+
				"writer naming a non-owned tag (C2)",
				f.Category, table.CatWriteToNonOwnedTag)
		}
	})

	// `normalize.go::renderWrites` names this arm by name: "a clear on a
	// decision table names a non-owned tag by construction and refuses
	// through 0002's existing `write to non-owned tag` arm, which is what
	// C2 relies on." Every other clear-refusal test in the corpus is taken
	// on a STATE-MACHINE fixture, so the decision-table half of that claim
	// had no assertion.
	t.Run("a clear list on a decision-table rule refuses", func(t *testing.T) {
		// `clear` is a RULE-level key, so it must precede the first
		// sub-table; appending it after `[rule.guard.all.b]` would author a
		// guard key named `clear` instead.
		src := strings.Replace(dtComplete,
			"id = \"cell-xp\"\n", "id = \"cell-xp\"\nclear = [\"a\"]\n", 1)
		if src == dtComplete {
			t.Fatal("the clear-list substitution did not apply")
		}
		f := refuseSource(t, src, "dt-clear-list.toml")
		if f.Category != table.CatWriteToNonOwnedTag {
			t.Errorf("category = %q; want %q — a clear on a decision-table "+
				"rule names a non-owned tag by construction and refuses "+
				"through 0002's existing arm, which is what C2 relies on",
				f.Category, table.CatWriteToNonOwnedTag)
		}
	})

	t.Run("a write block naming a non-owned tag refuses", func(t *testing.T) {
		src := strings.Replace(dtComplete,
			"[rule.emit]\nverdict = \"alpha\"\n",
			"[rule.write]\na = \"x\"\n[rule.emit]\nverdict = \"alpha\"\n", 1)
		f := refuseSource(t, src, "dt-write-block.toml")
		switch f.Category {
		case table.CatWriteToNonOwnedTag, table.CatMalformedAccessorBinding,
			table.CatUnknownTag:
		default:
			t.Errorf("category = %q; a write block on a decision-table row "+
				"refuses through 0002's existing write/binding arms, not a "+
				"new one", f.Category)
		}
	})
}

// REQ-17: "Rule ids, the match-block obligation, shared contexts, guards,
// gates, escape rules, and outcome binding are unchanged in both classes."
// DOMAIN EDGE — negative REQ.
//
// GREEN BY CONSTRUCTION for the state-machine half; the decision-table half
// is new behaviour (the model must load at all).
func TestReq17_TheMatchBlockObligationHoldsInBothClasses(t *testing.T) {
	t.Run("decision-table rule with no match block refuses", func(t *testing.T) {
		src := strings.Replace(dtComplete,
			"[rule.match.recognized]\neq = \"decide\"\n", "", 1)
		f := refuseSource(t, src, "dt-no-match.toml")
		if f.Category != table.CatMalformedRuleShape &&
			f.Category != table.CatMalformedOutcomeBinding {
			t.Errorf("category = %q; the match-block obligation is unchanged "+
				"for the decision-table class", f.Category)
		}
	})

	t.Run("duplicate rule id refuses in the decision-table class", func(t *testing.T) {
		src := strings.Replace(dtComplete, `id = "cell-yq"`, `id = "cell-xp"`, 1)
		f := refuseSource(t, src, "dt-dup-rule-id.toml")
		if f.Category != table.CatDuplicateRuleID {
			t.Errorf("category = %q; want %q — rule ids are unchanged in "+
				"both classes", f.Category, table.CatDuplicateRuleID)
		}
	})
}

// REQ-19: "Any rule, ordinary or escape, in either class MAY carry
// `[rule.emit]`: a flat TOML table whose values MUST be strings."
// HAPPY PATH — both classes, both rule kinds.
func TestReq19_EmitIsAdmittedOnBothRuleKindsInBothClasses(t *testing.T) {
	t.Run("ordinary rule, decision-table class", func(t *testing.T) {
		m := loadSource(t, dtComplete, "dt-complete.toml")
		got, _ := emitOf(rowByRuleID(t, m, "cell-xp"))
		if got["verdict"] != "alpha" {
			t.Errorf("emit = %v; want verdict=alpha", got)
		}
	})

	t.Run("escape rule, decision-table class", func(t *testing.T) {
		m := loadSource(t, dtEscapeVariant, "dt-escape.toml")
		got, _ := emitOf(rowByRuleID(t, m, "otherwise"))
		if got["verdict"] != "fallback" {
			t.Errorf("escape row emit = %v; want verdict=fallback", got)
		}
	})

	t.Run("ordinary rule, state-machine class", func(t *testing.T) {
		src := strings.Replace(smHeader+`
[[rule]]
id = "advance"
[rule.match.recognized]
eq = "decide"
[rule.match.status]
eq = "draft"
[rule.write]
status = "done"
[rule.emit]
verdict = "sm-answer"
`, "\x00", "", 1)
		m := loadSource(t, src, "sm-emit.toml")
		got, _ := emitOf(rowByRuleID(t, m, "advance"))
		if got["verdict"] != "sm-answer" {
			t.Errorf("state-machine row emit = %v; want verdict=sm-answer", got)
		}
	})

	// The fourth quadrant. The three above cover dt/ordinary, dt/escape and
	// sm/ordinary; without this one "either class, either rule kind" is
	// asserted for three of the four combinations, and the arm that admits
	// an emit block on a state-machine ESCAPE row — the row that carries
	// neither a write block nor a clear list — is unexercised.
	t.Run("escape rule, state-machine class", func(t *testing.T) {
		src := smHeader + `
[[rule]]
id = "advance"
[rule.match.recognized]
eq = "decide"
[rule.match.status]
eq = "draft"
[rule.write]
status = "done"

[[rule]]
id = "sm-otherwise"
escape = ["no_match"]
[rule.match.recognized]
eq = "decide"
[rule.emit]
verdict = "sm-fallback"
`
		m := loadSource(t, src, "sm-escape-emit.toml")
		got, _ := emitOf(rowByRuleID(t, m, "sm-otherwise"))
		if got["verdict"] != "sm-fallback" {
			t.Errorf("state-machine ESCAPE row emit = %v; want "+
				"verdict=sm-fallback", got)
		}
	})
}

// REQ-20 / REQ-21: "`sourceRule.Emit` is typed `map[string]string`, so a
// non-string value — a nested table (`[rule.emit.sub]`) included — is
// refused by the decoder as a type error, which
// `internal/table/source.go::decodeStrict` maps to `malformed TOML`" /
// "The assertion is on the category, never on upstream message text
// (`0002:C24`), and no hand-written type check is added."
// ADVERSARIAL
func TestReq20_ANonStringEmitValueRefusesAsMalformedTOML(t *testing.T) {
	cases := map[string]string{
		"integer": "[rule.emit]\nverdict = 42\n",
		"boolean": "[rule.emit]\nverdict = true\n",
		"array":   "[rule.emit]\nverdict = [\"a\", \"b\"]\n",
		"nested table": "[rule.emit]\nverdict = \"alpha\"\n" +
			"[rule.emit.sub]\ninner = \"x\"\n",
	}

	for name, block := range cases {
		t.Run(name, func(t *testing.T) {
			src := strings.Replace(dtComplete,
				"[rule.emit]\nverdict = \"alpha\"\n", block, 1)
			if src == dtComplete {
				t.Fatal("the emit substitution did not apply")
			}
			f := refuseSource(t, src, "dt-bad-emit-value.toml")
			if f.Category != table.CatMalformedTOML {
				t.Errorf("category = %q; want %q — the decoder's type error "+
					"carries this, and no hand-written type check is added "+
					"(CatUnknownSchemaField is the unknown-KEY arm)",
					f.Category, table.CatMalformedTOML)
			}
		})
	}
}

// REQ-22: "The normalized row carries `Emit []EmitValue`, a **new**
// row-level type `{Key, Value string}` — *not* `TagValue`"
// BOUNDARY
func TestReq22_RowCarriesEmitAsANewKeyValueStringType(t *testing.T) {
	rt := reflect.TypeOf(table.Row{})
	f, ok := rt.FieldByName("Emit")
	if !ok {
		t.Fatalf("table.Row carries no `Emit` field; fields = %v", fieldNames(rt))
	}
	if f.Type.Kind() != reflect.Slice {
		t.Fatalf("Row.Emit is %s; C3 fixes `Emit []EmitValue`", f.Type)
	}

	et := f.Type.Elem()
	if et == reflect.TypeOf(table.TagValue{}) {
		t.Fatal("Row.Emit's element type is table.TagValue; C3 fixes a NEW " +
			"row-level type, explicitly *not* TagValue")
	}
	if et.Name() != "EmitValue" {
		t.Errorf("Row.Emit's element type is %s; C3 names it EmitValue", et.Name())
	}

	key, ok := et.FieldByName("Key")
	if !ok || key.Type.Kind() != reflect.String {
		t.Errorf("%s carries no `Key string`; C3 fixes `{Key, Value string}`",
			et.Name())
	}
	val, ok := et.FieldByName("Value")
	if !ok || val.Type.Kind() != reflect.String {
		t.Errorf("%s carries no `Value string`; a TagValue-shaped "+
			"`Value []string` is what C3 rules out", et.Name())
	}
}

// REQ-28: "Normalization MUST carry the block on the row as a key-sorted
// sequence; an absent block normalizes to an empty sequence, and a
// **present but empty** `[rule.emit]` normalizes to the same empty sequence
// rather than refusing — `emit` keys on length, not on key presence"
// REQ-91 / SC-3: unordered keys, an absent block, a present-but-empty block.
// INPUT EDGE — the deliberate divergence from the write/clear/gate rule.
func TestReq28_EmitNormalizesKeySortedAndKeysOnLengthNotPresence(t *testing.T) {
	t.Run("unordered keys normalize sorted", func(t *testing.T) {
		m := loadSource(t, dtUnordered, "dt-unordered.toml")
		_, order := emitOf(rowByRuleID(t, m, "cell-xp"))
		want := []string{"alpha", "mu", "zeta"}
		if !slices.Equal(order, want) {
			t.Errorf("emit key order = %v; want %v — normalization sorts by "+
				"key, and the fixture authors zeta/alpha/mu", order, want)
		}
	})

	t.Run("an absent block normalizes to an empty sequence", func(t *testing.T) {
		m := loadSource(t, dtComplete, "dt-complete.toml")
		r := rowByRuleID(t, m, "cell-yq")
		if len(r.Emit) != 0 {
			t.Errorf("cell-yq authors no [rule.emit]; Emit = %v, want empty",
				r.Emit)
		}
	})

	t.Run("a present but empty block loads and normalizes to empty", func(t *testing.T) {
		m := loadSource(t, dtEmptyEmit, "dt-empty-emit.toml")
		r := rowByRuleID(t, m, "cell-xp")
		if len(r.Emit) != 0 {
			t.Errorf("a present-but-empty [rule.emit] normalized to %v; want "+
				"the same empty sequence an absent block gives", r.Emit)
		}
	})
}

// REQ-23 / REQ-93: "Duplicate keys are unreachable — a TOML table refuses
// them in the decoder before this code runs — so normalization sorts by key
// and asserts no dedup pass" / "No duplicate-key case: … no dedup arm
// exists to test."
// ADVERSARIAL — negative REQ, asserted in the only direction available: the
// DECODER refuses, so the case never reaches normalization.
func TestReq23_DuplicateEmitKeysAreRefusedByTheDecoder(t *testing.T) {
	src := strings.Replace(dtComplete,
		"[rule.emit]\nverdict = \"alpha\"\n",
		"[rule.emit]\nverdict = \"alpha\"\nverdict = \"beta\"\n", 1)
	if src == dtComplete {
		t.Fatal("the duplicate-key substitution did not apply")
	}

	f := refuseSource(t, src, "dt-dup-emit-key.toml")
	if f.Category != table.CatMalformedTOML {
		t.Errorf("category = %q; want %q — a TOML table refuses duplicate "+
			"keys in the decoder, which is why no dedup arm exists",
			f.Category, table.CatMalformedTOML)
	}
}

// REQ-24 / REQ-92 / SC-3 (`0010:S3`): "`Row.Emit` MUST be carried through
// `internal/table/normalize.go::expand` for every expanded row, which
// builds each row from the seed literal rather than copying `base`." /
// "**every** row expanded from the `in`-atom rule carries the authored
// block byte-for-byte"
// ADVERSARIAL — the Trace's unwitnessed step 4′; the mandatory SC-3 case.
func TestReq24_EveryExpandedRowCarriesTheAuthoredEmitBlock(t *testing.T) {
	m := loadSource(t, dtExpanding, "dt-expanding.toml")

	rows := rowsByRuleID(m, "spread")
	if len(rows) != 2 {
		t.Fatalf("the `in`-atom rule expanded to %d rows; want 2 — the "+
			"fixture's match atom carries two members", len(rows))
	}

	want := map[string]string{"verdict": "shared", "note": "same-block"}
	wantOrder := []string{"note", "verdict"}
	for _, r := range rows {
		got, order := emitOf(r)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: emit = %v; want %v byte-for-byte — expand builds "+
				"each row from the seed literal and must carry Emit",
				r.Identity(), got, want)
		}
		if !slices.Equal(order, wantOrder) {
			t.Errorf("%s: emit key order = %v; want %v",
				r.Identity(), order, wantOrder)
		}
	}
}

// REQ-25 / REQ-26: "The carried sequence MAY be **shared** across the rows
// one rule expands to — no clone is required." / "`Emit` is never mutated
// after normalization: it is read by the dump column and by the
// `resolvePayload` join and by nothing else, and both are readers."
// ADVERSARIAL — negative REQ: rendering must not feed back into the value.
//
// Sharing and cloning are indistinguishable to the byte-for-byte assertion
// above; what IS assertable is that no reader mutates the sequence.
func TestReq26_EmitIsNeverMutatedByTheDumpReader(t *testing.T) {
	m := loadSource(t, dtExpanding, "dt-expanding.toml")

	before := make([][]table.EmitValue, len(m.Rows))
	for i, r := range m.Rows {
		before[i] = slices.Clone(r.Emit)
	}

	_ = table.Dump(m)
	_ = table.Dump(m)

	for i, r := range m.Rows {
		if !slices.Equal(r.Emit, before[i]) {
			t.Errorf("%s: rendering mutated Emit: %v -> %v; the dump column "+
				"is a READER", r.Identity(), before[i], r.Emit)
		}
	}
}

// REQ-27: "Emit keys are NOT tag keys: they are undeclared, uninterpreted,
// compared by exact byte equality, and MUST NOT be matched, guarded,
// written, or read by any accessor"
// ADVERSARIAL — a `[rule.match.<emit-key>]` refuses `unknown tag` (SC-3).
func TestReq27_AnEmitKeyIsNotATagKey(t *testing.T) {
	t.Run("[rule.match.<emit-key>] refuses unknown tag", func(t *testing.T) {
		src := strings.Replace(dtComplete,
			"[rule.guard.all.a]\neq = \"x\"\n[rule.guard.all.b]\neq = \"p\"\n[rule.emit]\nverdict = \"alpha\"\n",
			"[rule.match.verdict]\neq = \"alpha\"\n[rule.guard.all.a]\neq = \"x\"\n[rule.guard.all.b]\neq = \"p\"\n[rule.emit]\nverdict = \"alpha\"\n",
			1)
		if src == dtComplete {
			t.Fatal("the match substitution did not apply")
		}
		f := refuseSource(t, src, "dt-match-emit-key.toml")
		if f.Category != table.CatUnknownTag {
			t.Errorf("category = %q; want %q — an emit key is not a tag key, "+
				"declared or not", f.Category, table.CatUnknownTag)
		}
	})

	t.Run("[rule.guard.all.<emit-key>] refuses unknown tag", func(t *testing.T) {
		src := strings.Replace(dtComplete,
			"[rule.guard.all.b]\neq = \"p\"\n[rule.emit]\nverdict = \"alpha\"\n",
			"[rule.guard.all.b]\neq = \"p\"\n[rule.guard.all.verdict]\neq = \"alpha\"\n[rule.emit]\nverdict = \"alpha\"\n",
			1)
		if src == dtComplete {
			t.Fatal("the guard substitution did not apply")
		}
		f := refuseSource(t, src, "dt-guard-emit-key.toml")
		if f.Category != table.CatUnknownTag {
			t.Errorf("category = %q; want %q", f.Category, table.CatUnknownTag)
		}
	})

	t.Run("an emit key does not enter the tag table", func(t *testing.T) {
		m := loadSource(t, dtComplete, "dt-complete.toml")
		if _, declared := m.Tags["verdict"]; declared {
			t.Error("the emit key `verdict` appears in the model's tag table; " +
				"an emit key is not a tag key, declared or not")
		}
	})
}

// REQ-27 (byte equality): emit keys and values are compared by EXACT byte
// equality and are uninterpreted — no case folding, no trimming, no
// value coercion.
// INPUT EDGE
func TestReq27_EmitKeysAndValuesAreCarriedByteForByte(t *testing.T) {
	src := strings.Replace(dtComplete,
		"[rule.emit]\nverdict = \"alpha\"\n",
		"[rule.emit]\n\"Mixed Case\" = \"  spaced  \"\n\"a<b&c\" = \"<clear>\"\n",
		1)
	if src == dtComplete {
		t.Fatal("the emit substitution did not apply")
	}

	m := loadSource(t, src, "dt-emit-bytes.toml")
	got, _ := emitOf(rowByRuleID(t, m, "cell-xp"))

	want := map[string]string{"Mixed Case": "  spaced  ", "a<b&c": "<clear>"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("emit = %#v; want %#v byte-for-byte — emit values are "+
			"uninterpreted, and the `<clear>` sentinel is not a value emit "+
			"keys have", got, want)
	}
}

// REQ-29: "`emit` MUST join the closed dump column vocabulary (`0002:C19`)
// **appended last**, after `escape`"
// BOUNDARY — the only insertion that leaves every existing column at its
// existing index.
func TestReq29_EmitJoinsTheDumpVocabularyAppendedLast(t *testing.T) {
	got := table.DumpColumns()
	want := []string{
		"identity", "source", "kind", "outcome",
		"atoms", "next", "writes", "requires_owned", "gate", "escape", "emit",
	}
	if !slices.Equal(got, want) {
		t.Fatalf("dump column vocabulary = %v; want %v — `emit` is appended "+
			"LAST, after `escape`", got, want)
	}
}

// REQ-30 / REQ-31 / ASSUMPTION-4: "Its cell is rendered as `key=value`
// pairs in key order, bracketed like `writes` (`[a=x; b=y]`), with the
// value emitted as the raw authored string — **not** through
// `renderValue`" / "An empty block renders as the empty bracket the
// `writes` column already uses for an empty sequence."
// HAPPY PATH
func TestReq30_TheEmitCellRendersKeyValuePairsBracketedLikeWrites(t *testing.T) {
	m := loadSource(t, dtUnordered, "dt-unordered.toml")
	dump := table.Dump(m)

	if !strings.Contains(dump, "emit=[alpha=first; mu=middle; zeta=last]") {
		t.Errorf("dump does not render the emit cell as `key=value` pairs in "+
			"KEY order bracketed like `writes`:\n%s", dump)
	}

	t.Run("an empty block renders as the empty bracket", func(t *testing.T) {
		line := dumpLineFor(t, dump, "dt.cell-yq")
		if !strings.Contains(line, "emit=[]") {
			t.Errorf("a row authoring no [rule.emit] renders %q; want the "+
				"empty bracket `emit=[]`", line)
		}
	})
}

// REQ-30 (raw string): the value is the RAW authored string, not routed
// through `renderValue`, whose quoting and bracketing key on a tag's
// declared kind and member count — neither of which an emit key has.
// ADVERSARIAL
func TestReq30_TheEmitValueIsTheRawAuthoredStringNotRenderValue(t *testing.T) {
	src := strings.Replace(dtComplete,
		"[rule.emit]\nverdict = \"alpha\"\n",
		"[rule.emit]\nverdict = \"a,b\"\n", 1)
	if src == dtComplete {
		t.Fatal("the emit substitution did not apply")
	}

	m := loadSource(t, src, "dt-emit-raw.toml")
	line := dumpLineFor(t, table.Dump(m), "dt.cell-xp")

	if !strings.Contains(line, "emit=[verdict=a,b]") {
		t.Errorf("emit cell in %q; want `emit=[verdict=a,b]` — the value is "+
			"the RAW authored string, never quoted or bracketed by "+
			"renderValue", line)
	}
}

// REQ-32: "A `[dump]` column list MUST name it, and the default column set
// used when `[dump]` is absent MUST include it for both classes."
// ADVERSARIAL — a `[dump]` list omitting `emit` refuses `malformed dump
// declaration` (FM, SC-3).
func TestReq32_ADumpListOmittingEmitRefuses(t *testing.T) {
	t.Run("the default set includes emit for both classes", func(t *testing.T) {
		for name, src := range map[string]string{
			"decision-table": dtComplete,
			"state-machine":  smZeroOwned,
		} {
			t.Run(name, func(t *testing.T) {
				m := loadSource(t, src, name+".toml")
				if !slices.Contains(m.DumpOrder, "emit") {
					t.Errorf("DumpOrder = %v; the default column set used when "+
						"`[dump]` is absent MUST include `emit`", m.DumpOrder)
				}
				if !strings.Contains(table.Dump(m), "emit=") {
					t.Error("the rendered dump carries no `emit` column")
				}
			})
		}
	})

	t.Run("an explicit [dump] list omitting emit refuses", func(t *testing.T) {
		src := dtComplete + `
[dump]
order = [
  "identity", "source", "kind", "outcome",
  "atoms", "next", "writes", "requires_owned", "gate", "escape",
]
`
		f := refuseSource(t, src, "dt-dump-omits-emit.toml")
		if f.Category != table.CatMalformedDumpDeclaration {
			t.Errorf("category = %q; want %q",
				f.Category, table.CatMalformedDumpDeclaration)
		}
	})

	t.Run("an explicit [dump] list naming emit loads", func(t *testing.T) {
		src := dtComplete + `
[dump]
order = [
  "identity", "source", "kind", "outcome",
  "atoms", "next", "writes", "requires_owned", "gate", "escape", "emit",
]
`
		m := loadSource(t, src, "dt-dump-names-emit.toml")
		if !slices.Contains(m.DumpOrder, "emit") {
			t.Errorf("DumpOrder = %v; want it to carry `emit`", m.DumpOrder)
		}
	})
}

// REQ-33: "`emit` MUST NOT be added to the kernel row; it is table data
// joined back by rule id after selection"
// DOMAIN EDGE — negative REQ; the kernel row is untouched (TD item 1).
func TestReq33_EmitNeverCrossesTheKernelBoundary(t *testing.T) {
	rt := reflect.TypeOf(table.Row{}.KernelRow())
	for i := range rt.NumField() {
		if strings.EqualFold(rt.Field(i).Name, "emit") {
			t.Fatalf("resolve.Row carries a %q field; `emit` MUST NOT be "+
				"added to the kernel row", rt.Field(i).Name)
		}
	}

	m := loadSource(t, dtComplete, "dt-complete.toml")
	kt := m.KernelTable()
	if len(kt.Rows) != len(m.Rows) {
		t.Errorf("kernel table carries %d rows; the model carries %d",
			len(kt.Rows), len(m.Rows))
	}
}

// REQ-71 / REQ-87 / ASSUMPTION-8: "The one non-silent widening is the dump
// vocabulary: an explicit `[dump]` list must gain `emit` or fail at load
// (A4, 103 fixtures)"
// DOMAIN EDGE — every checked-in fixture carrying an explicit `[dump]` list
// is updated in this change, or it stops loading.
func TestReq71_EveryCheckedInDumpCarryingFixtureNamesEmit(t *testing.T) {
	for _, rel := range []string{rdrFixture, kataFixture} {
		t.Run(rel, func(t *testing.T) {
			m := mustLoad(t, rel)
			if !slices.Contains(m.DumpOrder, "emit") {
				t.Errorf("%s: [dump].order = %v; the explicit list must gain "+
					"`emit` or the fixture fails at load", rel, m.DumpOrder)
			}
		})
	}
}

// --- helpers -------------------------------------------------------------

// rowByRuleID returns the single normalized row a rule id names, failing
// when the rule expanded to more than one.
func rowByRuleID(t *testing.T, m *table.Model, ruleID string) table.Row {
	t.Helper()

	rows := rowsByRuleID(m, ruleID)
	if len(rows) != 1 {
		t.Fatalf("%d rows carry rule id %q; want exactly 1 (identities: %v)",
			len(rows), ruleID, rowIdentities(m))
	}
	return rows[0]
}

// dumpLineFor returns the dump line whose `identity=` cell is want.
func dumpLineFor(t *testing.T, dump, want string) string {
	t.Helper()

	for _, line := range strings.Split(dump, "\n") {
		if strings.HasPrefix(line, "identity="+want+" ") {
			return line
		}
	}
	t.Fatalf("no dump line for identity %q:\n%s", want, dump)
	return ""
}
