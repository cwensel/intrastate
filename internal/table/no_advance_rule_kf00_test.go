package table_test

// Kata kf00 — a state-machine rule must be able to DECIDE without
// advancing, so one model can both apply and refuse.
//
// `0010:C2` gives a decision table zero owned tags, so it can never write;
// `0002:C4`, conditioned on the state-machine class by `0010`, made every
// ordinary rule carry `[rule.write]`, so it could never merely decide. A
// model was therefore entirely decide-only or entirely write-per-rule,
// while a real flow is mixed: the rows that refuse must emit their reason
// and write NOTHING.
//
// `advance = false` on `[[rule]]` is the opt-out. It is EXPLICIT rather
// than inferred from the missing `[rule.write]`, so an author who merely
// forgot the block still takes the refusal, and it is a rule-level marker
// rather than a reading of the emit's disposition, because RDR 0024 fixes
// no disposition vocabulary and intrastate never interprets a disposition
// token — the loader has no `stop`/`none` spelling to key on.
//
// Fixtures are authored in-test on the `smHeader` state-machine preamble,
// the same shape the RDR 0010 class suite uses. A new file under
// `testdata/neg/` would be wrong here: that directory is the REQ-118
// promoted spike set, pinned per fixture to the category it witnessed.

import (
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// smEmitNoMarker is the RED fixture: a state-machine ordinary rule that
// emits and carries no write block, and does NOT declare `advance`. The
// emit alone must not lift the obligation — that is the forgotten-block
// case the refusal exists to catch.
const smEmitNoMarker = smHeader + `
[emit.op]
kind = "enum"
domain = ["stopped:no-gate"]

[[rule]]
id = "stop-no-marker"
[rule.match.recognized]
eq = "decide"
[rule.match.status]
eq = "draft"
[rule.emit]
op = "stopped:no-gate"
`

// smNoAdvance is the GREEN fixture: the SAME rule carrying the explicit
// marker. It must load clean and normalize to an empty write surface.
const smNoAdvance = smHeader + `
[emit.op]
kind = "enum"
domain = ["stopped:no-gate"]

[[rule]]
id = "stop-no-gate"
advance = false
[rule.match.recognized]
eq = "decide"
[rule.match.status]
eq = "draft"
[rule.emit]
op = "stopped:no-gate"
`

// TestKataKf00_EmitWithoutTheMarkerStillRefuses is the RED half: the
// strictness `0002:C4` buys must not regress. A rule that emits and
// forgets its write block is still malformed.
func TestKataKf00_EmitWithoutTheMarkerStillRefuses(t *testing.T) {
	f := refuseSource(t, smEmitNoMarker, "sm-emit-no-marker.toml")
	if f.Category != table.CatMalformedRuleShape {
		t.Errorf("category = %q; want %q — an emit block is not the opt-out, "+
			"only the explicit `advance = false` is, so an author who merely "+
			"forgot `[rule.write]` still takes this refusal",
			f.Category, table.CatMalformedRuleShape)
	}

	// The discriminating control: the ONLY difference between the two
	// fixtures is the marker. If this stopped loading the RED assertion
	// above would pass for the wrong reason.
	loadSource(t, smNoAdvance, "sm-no-advance.toml")
}

// TestKataKf00_NoAdvanceRowNormalizesToAnEmptyWriteSurface is the GREEN
// half at the loader: the row carries no writes, no next-state tags, and
// no required-owned keys, so nothing downstream can plan a write from it.
func TestKataKf00_NoAdvanceRowNormalizesToAnEmptyWriteSurface(t *testing.T) {
	m := loadSource(t, smNoAdvance, "sm-no-advance.toml")

	row := rowByID(t, m, "sm.stop-no-gate")
	if len(row.Writes) != 0 {
		t.Errorf("Writes = %v; want empty — a non-advancing row must plan no write", row.Writes)
	}
	if len(row.NextTags) != 0 {
		t.Errorf("NextTags = %v; want empty — a non-advancing row names no next state", row.NextTags)
	}
	if len(row.RequiresOwned) != 0 {
		t.Errorf("RequiresOwned = %v; want empty", row.RequiresOwned)
	}
	// The row is an ORDINARY transition row, not an escape row: it is
	// selected by its own predicate, never by a resolver failure class.
	if row.Kind() != table.KindTransition {
		t.Errorf("Kind() = %v; want %v", row.Kind(), table.KindTransition)
	}
	// The emit is what makes declining to advance an ANSWER rather than a
	// silent no-op, so it must survive normalization.
	if len(row.Emit) != 1 || row.Emit[0].Key != "op" {
		t.Errorf("Emit = %v; want the authored `op` pair", row.Emit)
	}
}

// TestKataKf00_TheMarkerIsRefusedWhereItDeclaresNothing pins the arms that
// keep `advance = false` meaningful rather than decorative.
func TestKataKf00_TheMarkerIsRefusedWhereItDeclaresNothing(t *testing.T) {
	// A non-advancing row whose successor would differ from its source is
	// a contradiction. Both carriers are presence-keyed, the same way the
	// escape-row arms are.
	t.Run("with a write block", func(t *testing.T) {
		src := smHeader + `
[[rule]]
id = "contradiction"
advance = false
[rule.match.recognized]
eq = "decide"
[rule.match.status]
eq = "draft"
[rule.write]
status = "done"
`
		if f := refuseSource(t, src, "sm-no-advance-write.toml"); f.Category != table.CatMalformedRuleShape {
			t.Errorf("category = %q; want %q", f.Category, table.CatMalformedRuleShape)
		}
	})

	t.Run("with a clear list", func(t *testing.T) {
		src := smHeader + `
[[rule]]
id = "contradiction"
advance = false
clear = ["status"]
[rule.match.recognized]
eq = "decide"
[rule.match.status]
eq = "draft"
`
		if f := refuseSource(t, src, "sm-no-advance-clear.toml"); f.Category != table.CatMalformedRuleShape {
			t.Errorf("category = %q; want %q", f.Category, table.CatMalformedRuleShape)
		}
	})

	// A row that neither advances nor answers is INERT: it can be selected
	// and contributes nothing the caller can branch on.
	t.Run("with no emit block", func(t *testing.T) {
		src := smHeader + `
[[rule]]
id = "inert"
advance = false
[rule.match.recognized]
eq = "decide"
[rule.match.status]
eq = "draft"
`
		if f := refuseSource(t, src, "sm-no-advance-inert.toml"); f.Category != table.CatMalformedRuleShape {
			t.Errorf("category = %q; want %q", f.Category, table.CatMalformedRuleShape)
		}
	})

	// An escape row already advances nothing — it carries neither a write
	// block nor a clear list — so the marker declares what the escape list
	// already fixes. It joins that arm's presence-keyed list.
	t.Run("on an escape row", func(t *testing.T) {
		src := smHeader + `
[[rule]]
id = "rescue"
advance = false
escape = ["no_match"]
[rule.match.recognized]
eq = "decide"
`
		if f := refuseSource(t, src, "sm-escape-advance.toml"); f.Category != table.CatMalformedEscapeDeclaration {
			t.Errorf("category = %q; want %q", f.Category, table.CatMalformedEscapeDeclaration)
		}
	})

	// A decision table owns no state (`0010:C1`, C2), so NO row in it ever
	// advances and the marker is a per-row property the class does not
	// have.
	t.Run("on a decision-table rule", func(t *testing.T) {
		src := dtHeader + `
[[rule]]
id = "decide-only"
advance = false
[rule.match.recognized]
eq = "decide"
[rule.match.a]
eq = "x"
[rule.emit]
verdict = "yes"
`
		if f := refuseSource(t, src, "dt-advance.toml"); f.Category != table.CatMalformedRuleShape {
			t.Errorf("category = %q; want %q", f.Category, table.CatMalformedRuleShape)
		}
	})
}

// TestKataKf00_ExplicitAdvanceTrueIsTheUnchangedDefault pins that only the
// explicit FALSE lifts the obligation: `advance = true` restates the
// default and leaves `0002:C4` binding.
func TestKataKf00_ExplicitAdvanceTrueIsTheUnchangedDefault(t *testing.T) {
	src := smHeader + `
[[rule]]
id = "advances"
advance = true
[rule.match.recognized]
eq = "decide"
[rule.match.status]
eq = "draft"
`
	if f := refuseSource(t, src, "sm-advance-true.toml"); f.Category != table.CatMalformedRuleShape {
		t.Errorf("category = %q; want %q — `advance = true` is the default, "+
			"not an opt-out", f.Category, table.CatMalformedRuleShape)
	}
}
