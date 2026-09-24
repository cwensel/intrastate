package table_test

// RDR 0030 `0030:C1`/`0030:C2`/`0030:D-identity` — the EXPANSION a step
// write normalizes to, and the BOUND refusal, asserted on the normalized
// model and the refusal the loader returns.

import (
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

const (
	s30LadderStep       = "ladder-step.toml"
	s30LadderLiteral    = "ladder-literal.toml"
	s30LadderUnguarded  = "ladder-step-unguarded.toml"
	s30LadderUnguardedE = "ladder-step-unguarded-enum.toml"
)

// s30Ladder loads a committed ladder fixture.
func s30Ladder(t *testing.T, rel string) *table.Model {
	t.Helper()

	m, err := table.Load(readFixture(t, rel), rel)
	if err != nil {
		t.Fatalf("%s must load clean; refused: %v", rel, err)
	}
	return m
}

// s30Twin names the literal-ladder row a stepped-ladder row corresponds to.
func s30Twin(r table.Row) string {
	switch {
	case r.RuleID == "retry" && len(r.Suffix) == 1:
		return "retry-" + r.Suffix[0]
	case r.RuleID == "escalate" && len(r.Suffix) == 1:
		return "escalate-" + r.Suffix[0]
	default:
		return r.RuleID
	}
}

// s30Atom reports whether row carries the atom (key, block, op, literal…).
func s30Atom(r table.Row, key string, block table.Block, op string, literal ...string) bool {
	for _, a := range r.Atoms {
		if a.Key == key && a.Block == block && a.Operator == op && slices.Equal(a.Literal, literal) {
			return true
		}
	}
	return false
}

// REQ-11: "The loader expands a rule carrying a step write into literal
// rows, one per ADMITTED CELL of the stepped tag, composed into the rule's
// existing expansion product (`0002:C13`) as one more choice point. An
// admitted cell is a domain member … that satisfies the CONJUNCTION of every
// positive atom the rule AUTHORS on that tag, however many per block …, each
// evaluated per member as the runtime evaluator would, against the AUTHORED
// literal and not against any member a sibling choice point has since
// chosen; `unless` atoms are not consulted."
// HAPPY PATH
func TestReq11_0030_AdmittedCellsAreTheConjunctionOfEveryPositiveAtomOnTheTag(t *testing.T) {
	rule := `[[rule]]
id = "retry"
use = ["open"]
[rule.match.recognized]
eq = "retry"
[rule.match.attempt]
in = [1, 2, 3, 4, 7]
[rule.match.tier]
in = ["small", "mid"]
[rule.guard.all.attempt]
gte = 2
lt = 6
in = [2, 3, 4, 5]
[rule.guard.unless.attempt]
eq = 3
[rule.write]
attempt = { step = 1 }
`
	m := s30MustLoad(t, "a four-atom conjunction on attempt beside an in on tier", s30Src(s30Opts{rules: rule}))

	// {1,2,3,4,7} ∩ [2,∞) ∩ (−∞,6) ∩ {2,3,4,5} = {2,3,4}; the unless on 3 is
	// NOT consulted, and the tier `in` is its own choice point.
	var want []string
	for _, cell := range []string{"2", "3", "4"} {
		for _, tier := range []string{"mid", "small"} {
			want = append(want, "m.retry#"+cell+"#"+tier)
		}
	}
	slices.Sort(want)
	if got := s30RowIDs(m, "retry"); !slices.Equal(got, want) {
		t.Errorf("admitted rows = %v; want %v", got, want)
	}
}

// REQ-12: "on a stepped tag the step point SUBSUMES a match `in` on that tag
// rather than composing with it — whether the `in` is authored locally or
// inherited from a context … One choice point is emitted per stepped tag,
// never two, and the `in` contributes its members to the conjunction instead
// of expanding separately."
// HAPPY PATH
func TestReq12_0030_AStepSubsumesAMatchInOnItsTagAuthoredOrInherited(t *testing.T) {
	local := s30Src(s30Opts{rules: `[[rule]]
id = "retry"
use = ["open"]
[rule.match.recognized]
eq = "retry"
[rule.match.attempt]
in = [1, 3]
[rule.guard.all.attempt]
lt = 5
[rule.write]
attempt = { step = 1 }
`})
	inherited := s30Src(s30Opts{contexts: "\n[context.odd.match.attempt]\nin = [1, 3]\n", rules: `[[rule]]
id = "retry"
use = ["open", "odd"]
[rule.match.recognized]
eq = "retry"
[rule.guard.all.attempt]
lt = 5
[rule.write]
attempt = { step = 1 }
`})
	for _, tc := range []struct{ name, src string }{{"authored", local}, {"inherited", inherited}} {
		t.Run(tc.name, func(t *testing.T) {
			m := s30MustLoad(t, "a step with an "+tc.name+" match in on attempt", tc.src)
			want := []string{"m.retry#1", "m.retry#3"}
			if got := s30RowIDs(m, "retry"); !slices.Equal(got, want) {
				t.Errorf("rows = %v; want %v — |in ∩ admitted| rows, one choice point", got, want)
			}
			for _, r := range rowsByRuleID(m, "retry") {
				if len(r.Suffix) != 1 {
					t.Errorf("%s: suffix %v carries %d elements; one choice point mints one",
						r.Identity(), r.Suffix, len(r.Suffix))
				}
			}
		})
	}
}

// REQ-13: "A subsumed `in` is REWRITTEN per row to `match eq = <cell>`,
// exactly as `expand`'s existing `case expanding:` arm rewrites an `in` it
// expands — never retained in its `in` form." … "So this mechanism
// contributes exactly two atoms to the stepped tag per row — the rewritten
// `match eq = <cell>` where an `in` was subsumed, and the expansion's own
// `guard.all eq = <cell>` — and no more." … "every OTHER authored atom on the
// tag, the `guard.all` bounds among them, is retained unchanged beside these
// two"
// Read as resolved (deviations D13): the per-cell atom lands in the block the
// author scoped the tag with, so a row whose subsumed `in` is rewritten to
// `match eq = <cell>` carries no second `guard.all eq = <cell>` (a guard
// dimension beside the group's own match `eq` reads every other cell as a
// coverage gap, `0006:C7`); with nothing subsuming, the `guard.all eq` stands.
// HAPPY PATH
func TestReq13_0030_ASubsumedInIsRewrittenAndEveryOtherAtomIsRetained(t *testing.T) {
	m := s30MustLoad(t, "a step subsuming a match in", s30Src(s30Opts{rules: `[[rule]]
id = "retry"
use = ["open"]
[rule.match.recognized]
eq = "retry"
[rule.match.attempt]
in = [1, 3]
[rule.guard.all.attempt]
lt = 5
[rule.write]
attempt = { step = 1 }
`}))
	rows := rowsByRuleID(m, "retry")
	if len(rows) != 2 {
		t.Fatalf("rows = %v; want two", s30RowIDs(m, "retry"))
	}
	for _, r := range rows {
		cell := r.Suffix[len(r.Suffix)-1]
		on := atomsOn(r, "attempt")
		if len(on) != 2 {
			t.Errorf("%s carries %d atoms on attempt %+v; want exactly 2 (match eq, all lt)",
				r.Identity(), len(on), on)
		}
		if !s30Atom(r, "attempt", table.BlockMatch, "eq", cell) {
			t.Errorf("%s: no rewritten `match eq = %s`: %+v", r.Identity(), cell, on)
		}
		if s30Atom(r, "attempt", table.BlockAll, "eq", cell) {
			t.Errorf("%s: the match `eq` pins the cell, yet a `guard.all eq = %s` rides beside it: %+v",
				r.Identity(), cell, on)
		}
		if !s30Atom(r, "attempt", table.BlockAll, "lt", "5") {
			t.Errorf("%s: the authored `guard.all lt = 5` was not retained: %+v", r.Identity(), on)
		}
		for _, a := range on {
			if a.Operator == "in" {
				t.Errorf("%s retains an `in` atom on the stepped tag: %+v", r.Identity(), a)
			}
		}
	}

	// The committed ladder: with nothing to subsume, escalate's rows carry
	// their own `guard.all eq = <cell>` beside the RETAINED guard `in`, and
	// no match atom on the stepped tag appears from nowhere.
	ladder := s30Ladder(t, s30LadderStep)
	for _, cell := range []string{"small", "mid"} {
		r := s30Row(t, ladder, "ladder.escalate#"+cell)
		if !s30Atom(r, "tier", table.BlockAll, "eq", cell) {
			t.Errorf("%s: no `guard.all eq = %s`", r.Identity(), cell)
		}
		if !s30Atom(r, "tier", table.BlockAll, "in", "mid", "small") &&
			!s30Atom(r, "tier", table.BlockAll, "in", "small", "mid") {
			t.Errorf("%s: the authored guard `in` was not retained: %+v", r.Identity(), atomsOn(r, "tier"))
		}
		if len(atomsOn(r, "tier")) != 2 {
			t.Errorf("%s carries %d atoms on tier %+v; want exactly 2", r.Identity(), len(atomsOn(r, "tier")),
				atomsOn(r, "tier"))
		}
	}
}

// REQ-16: "Each expanded row carries a `guard.all` atom `eq = <cell>` on the
// tag, the literal write `cell + n` (`int`) or the domain member `n`
// positions from the cell (`enum`), and the cell appended to its expansion
// suffix."
// HAPPY PATH
func TestReq16_0030_EachRowCarriesItsCellGuardItsLiteralWriteAndItsSuffix(t *testing.T) {
	m := s30Ladder(t, s30LadderStep)
	for cell := range 5 {
		c := strconv.Itoa(cell)
		r := s30Row(t, m, "ladder.retry#"+c)
		if !s30Atom(r, "attempt", table.BlockAll, "eq", c) {
			t.Errorf("%s: no `guard.all eq = %s`", r.Identity(), c)
		}
		if got := s30Value(t, r.Identity(), r.Writes, "attempt"); got != strconv.Itoa(cell+1) {
			t.Errorf("%s writes attempt = %s; want %d", r.Identity(), got, cell+1)
		}
		if !slices.Equal(r.Suffix, []string{c}) {
			t.Errorf("%s: suffix %v; want [%s]", r.Identity(), r.Suffix, c)
		}
	}
	for cell, next := range map[string]string{"small": "mid", "mid": "large"} {
		r := s30Row(t, m, "ladder.escalate#"+cell)
		if got := s30Value(t, r.Identity(), r.Writes, "tier"); got != next {
			t.Errorf("%s writes tier = %s; want %s", r.Identity(), got, next)
		}
	}

	// n > 1 over an enum: the member n positions on.
	tier := "[tags.tier]\nprovenance = \"owned\"\nkind = \"enum\"\ndomain = [\"small\", \"a\", \"b\", \"c\"]\n" +
		"single_valued = true\nrequired = true\n"
	two := s30MustLoad(t, "an enum step of 2", s30Src(s30Opts{tier: tier,
		rules: s30Escalate(`in = ["small", "a"]`, "tier = { step = 2 }")}))
	for cell, next := range map[string]string{"small": "b", "a": "c"} {
		r := s30Row(t, two, "m.escalate#"+cell)
		if got := s30Value(t, r.Identity(), r.Writes, "tier"); got != next {
			t.Errorf("%s writes tier = %s; want %s (two positions on)", r.Identity(), got, next)
		}
	}
}

// REQ-17: "The suffix element is appended at EVERY admitted cell, including
// when exactly one is admitted: a stepped row is never the authored row, so
// it always carries its cell." … "the two arms keep their own rules, and the
// `in` expansion's behaviour is unchanged by this record."
// BOUNDARY
func TestReq17_0030_TheCellIsAppendedEvenWhenExactlyOneIsAdmitted(t *testing.T) {
	m := s30MustLoad(t, "a one-cell step", s30Src(s30Opts{rules: s30Retry("eq = 3", "attempt = { step = 1 }")}))
	if got := s30RowIDs(m, "retry"); !slices.Equal(got, []string{"m.retry#3"}) {
		t.Errorf("a one-cell step publishes %v; want [m.retry#3]", got)
	}

	// Control: a single-member `in` with no step still mints no suffix.
	in := s30MustLoad(t, "a single-member in", s30Src(s30Opts{rules: `[[rule]]
id = "retry"
use = ["open"]
[rule.match.recognized]
in = ["retry"]
[rule.write]
attempt = 1
`}))
	if got := s30RowIDs(in, "retry"); !slices.Equal(got, []string{"m.retry"}) {
		t.Errorf("a single-member in publishes %v; want [m.retry] (no suffix)", got)
	}
}

// REQ-18: "The stepped value is computed without overflow at every admitted
// cell, whatever the step's magnitude; whether it lands inside the domain is
// C2's bound, not an admission question."
// BOUNDARY
func TestReq18_0030_TheSteppedValueIsComputedWithoutOverflowInEitherDirection(t *testing.T) {
	attempt := "[tags.attempt]\nprovenance = \"owned\"\nkind = \"int\"\nmin = -5\nmax = 5\n" +
		"single_valued = true\nrequired = true\n"

	// MinInt64 from cell -5 overflows downward: named as the cell and the
	// bound, never as a wrapped positive literal.
	f := s30Refusal(t, "step MinInt64 over cells -5..-1", s30Src(s30Opts{attempt: attempt,
		rules: s30Retry("lt = 0", "attempt = { step = -9223372036854775808 }")}))
	s30RequireCategory(t, "downward overflow", f, table.CatMalformedTagDeclaration)
	if got := s30FirstIntAfter(f.Detail, "attempt"); got != "-5" {
		t.Errorf("the first number after the key is %q; want the cell -5 (domain order first): %s", got, f.Detail)
	}
	s30Contains(t, "downward overflow", f.Detail, "min")
	if strings.Contains(f.Detail, "9223372036854775803") {
		t.Errorf("the detail names a WRAPPED value: %s", f.Detail)
	}

	// MaxInt from cell 1 overflows upward.
	g := s30Refusal(t, "step MaxInt64 from cell 1", s30Src(s30Opts{attempt: attempt,
		rules: s30Retry("gte = 1\nlt = 2", "attempt = { step = 9223372036854775807 }")}))
	s30RequireCategory(t, "upward overflow", g, table.CatMalformedTagDeclaration)
	if strings.Contains(g.Detail, "-9223372036854775808") {
		t.Errorf("the detail names a WRAPPED value: %s", g.Detail)
	}
	// Twin: a small negative step over the same cells loads.
	s30MustLoad(t, "step -1 over cells -4..-1", s30Src(s30Opts{attempt: attempt,
		rules: s30Retry("gte = -4\nlt = 0", "attempt = { step = -1 }")}))
}

// REQ-21: "**The stepped key is placed in `assignments` like any other
// written key**, carrying a placeholder the expansion replaces per cell:
// `0002:C14` derives `RequiresOwned` from `maps.Keys(assignments)`, so a spec
// tracked beside the map instead would silently drop the stepped key from
// the kernel's owned-state gate."
// HAPPY PATH
func TestReq21_0030_EveryExpandedRowRequiresItsSteppedKeyOwned(t *testing.T) {
	step := s30Ladder(t, s30LadderStep)
	lit := s30Ladder(t, s30LadderLiteral)
	for _, r := range step.Rows {
		if r.RuleID != "retry" && r.RuleID != "escalate" {
			continue
		}
		key := map[string]string{"retry": "attempt", "escalate": "tier"}[r.RuleID]
		if !slices.Contains(r.RequiresOwned, key) {
			t.Errorf("%s: RequiresOwned %v drops the stepped key %q", r.Identity(), r.RequiresOwned, key)
		}
		twin := rowsByRuleID(lit, s30Twin(r))
		if len(twin) != 1 {
			t.Fatalf("no single literal twin %q for %s", s30Twin(r), r.Identity())
		}
		if !slices.Equal(r.RequiresOwned, twin[0].RequiresOwned) {
			t.Errorf("%s: RequiresOwned %v; its literal twin's is %v", r.Identity(), r.RequiresOwned,
				twin[0].RequiresOwned)
		}
	}
	if len(rowsByRuleID(step, "retry")) == 0 {
		t.Error("the stepped ladder expanded no retry rows")
	}
}

// REQ-23: "The per-cell stepped literal is the row's rendered assignment for
// that key, so it lands in BOTH carriers the assignment feeds: `Row.Writes`
// and `Row.NextTags` are populated per row from the same stepped value, never
// by aliasing one to the other (`0002:C15`)."
// ADVERSARIAL
func TestReq23_0030_WritesAndNextTagsBothCarryTheSteppedValueWithoutAliasing(t *testing.T) {
	m := s30Ladder(t, s30LadderStep)
	rows := rowsByRuleID(m, "retry")
	if len(rows) != 5 {
		t.Fatalf("retry expands to %d rows; want 5", len(rows))
	}
	for _, r := range rows {
		cell, _ := strconv.Atoi(r.Suffix[0])
		want := strconv.Itoa(cell + 1)
		if got := s30Value(t, r.Identity()+" writes", r.Writes, "attempt"); got != want {
			t.Errorf("%s: Writes attempt = %s; want %s", r.Identity(), got, want)
		}
		if got := s30Value(t, r.Identity()+" next", r.NextTags, "attempt"); got != want {
			t.Errorf("%s: NextTags attempt = %s; want %s", r.Identity(), got, want)
		}
	}

	// Mutation control: writing through one carrier of one row moves neither
	// the other carrier nor any sibling row.
	clone := cloneRows(m.Rows)
	for i := range m.Rows {
		if m.Rows[i].RuleID != "retry" {
			continue
		}
		for j := range m.Rows[i].Writes {
			if m.Rows[i].Writes[j].Key == "attempt" {
				m.Rows[i].Writes[j].Value[0] = "mutated"
			}
		}
		break
	}
	mutated := 0
	for i := range m.Rows {
		for j := range m.Rows[i].Writes {
			if !reflect.DeepEqual(m.Rows[i].Writes[j], clone[i].Writes[j]) {
				mutated++
			}
		}
		if !reflect.DeepEqual(m.Rows[i].NextTags, clone[i].NextTags) {
			t.Errorf("%s: NextTags moved when one row's Writes was mutated — the carriers alias",
				m.Rows[i].Identity())
		}
	}
	if mutated != 1 {
		t.Errorf("one Writes mutation changed %d write entries; the expanded rows share storage", mutated)
	}
}

// REQ-24: "the step expansion rides the same `expand` loop and copies
// `Emit`, `Gate`, `RequiresOwned` and `Escape` identically to every row it
// mints, so every row a stepped rule expands to carries the same block"
// HAPPY PATH
func TestReq24_0030_EveryExpandedRowCarriesTheRulesBlockIdentically(t *testing.T) {
	rule := `[[rule]]
id = "retry"
use = ["open"]
gate = ["lock"]
[rule.match.recognized]
eq = "retry"
[rule.guard.all.attempt]
lt = 5
[rule.write]
attempt = { step = 1 }
status = "open"
[rule.emit]
route = "again"
note = "retrying"
`
	gate := "\n[gate.lock]\nrole = \"state\"\npath = \"flow.lock\"\nkeys = [\"status\"]\ntimeout = \"2s\"\n"
	m := s30MustLoad(t, "a gated, emitting step rule", s30Src(s30Opts{contexts: gate, rules: rule}))
	rows := rowsByRuleID(m, "retry")
	if len(rows) != 5 {
		t.Fatalf("retry expands to %d rows; want 5", len(rows))
	}
	first := rows[0]
	if !slices.Equal(first.Gate, []string{"lock"}) {
		t.Errorf("%s: Gate = %v; want [lock]", first.Identity(), first.Gate)
	}
	wantEmit := []table.EmitValue{{Key: "note", Value: "retrying"}, {Key: "route", Value: "again"}}
	if !reflect.DeepEqual(first.Emit, wantEmit) {
		t.Errorf("%s: Emit = %+v; want %+v", first.Identity(), first.Emit, wantEmit)
	}
	for _, r := range rows[1:] {
		if !reflect.DeepEqual(r.Emit, first.Emit) || !slices.Equal(r.Gate, first.Gate) ||
			!slices.Equal(r.RequiresOwned, first.RequiresOwned) || !slices.Equal(r.Escape, first.Escape) {
			t.Errorf("%s carries a different block than %s:\n  %+v\n  %+v", r.Identity(), first.Identity(), r, first)
		}
	}
}

// REQ-25: "After normalization no surface distinguishes an expanded row from
// an authored literal row: both carriers hold literals only, and `0002:C4`'s
// write-replaces clause applies per row unchanged."
// HAPPY PATH
func TestReq25_0030_AnExpandedRowIsIndistinguishableFromItsLiteralTwin(t *testing.T) {
	step := s30Ladder(t, s30LadderStep)
	lit := s30Ladder(t, s30LadderLiteral)
	if len(step.Rows) != len(lit.Rows) {
		t.Fatalf("the stepped ladder normalizes to %d rows, the literal one to %d: %v vs %v",
			len(step.Rows), len(lit.Rows), rowIdentities(step), rowIdentities(lit))
	}
	for _, r := range step.Rows {
		twins := rowsByRuleID(lit, s30Twin(r))
		if len(twins) != 1 {
			t.Errorf("%s has no single literal twin %q", r.Identity(), s30Twin(r))
			continue
		}
		w := twins[0]
		for _, f := range []struct {
			name      string
			got, want any
		}{
			{"Outcome", r.Outcome, w.Outcome},
			{"Atoms", r.Atoms, w.Atoms},
			{"Gate", r.Gate, w.Gate},
			{"NextTags", r.NextTags, w.NextTags},
			{"Writes", r.Writes, w.Writes},
			{"RequiresOwned", r.RequiresOwned, w.RequiresOwned},
			{"Escape", r.Escape, w.Escape},
			{"Emit", r.Emit, w.Emit},
		} {
			if !reflect.DeepEqual(f.got, f.want) {
				t.Errorf("%s vs %s: %s differs\n  step:    %+v\n  literal: %+v",
					r.Identity(), w.Identity(), f.name, f.got, f.want)
			}
		}
		// The kernel sees the same row, bar the rule id.
		kr, kw := r.KernelRow(), w.KernelRow()
		kr.RuleID, kw.RuleID = "", ""
		kr.SourceLocator, kw.SourceLocator = "", ""
		if !reflect.DeepEqual(kr, kw) {
			t.Errorf("%s vs %s: the kernel rows differ\n  step:    %+v\n  literal: %+v",
				r.Identity(), w.Identity(), kr, kw)
		}
	}
}

// REQ-26: "**Step points order by KEY alone.**" … "So `attempt = { step = 1 }`
// with `tier = { step = 1 }` yields `retry#<attempt>#<tier>`, ordered by the
// key names. **The same key comparison places a step point among the OTHER
// candidates** — unsubsumed `in` points on other keys, and the outcome point
// … a rule with a step on `attempt` and an `in` on `mode` yields
// `retry#<attempt>#<mode>`." … "The element is the bare cell, not
// tag-qualified"
// HAPPY PATH
func TestReq26_0030_SuffixElementsOrderByKeyAlone(t *testing.T) {
	mode := "[tags.mode]\nprovenance = \"owned\"\nkind = \"enum\"\ndomain = [\"fast\", \"slow\"]\n" +
		"single_valued = true\n"
	head := func(match string) string {
		return "[[rule]]\nid = \"retry\"\nuse = [\"open\"]\n[rule.match.recognized]\n" + match + "\n"
	}
	for _, tc := range []struct {
		name string
		rule string
		want []string
	}{
		{"two steps: attempt before tier",
			head(`eq = "retry"`) + "[rule.guard.all.attempt]\nlt = 2\n[rule.guard.all.tier]\nin = [\"small\"]\n" +
				"[rule.write]\nattempt = { step = 1 }\ntier = { step = 1 }\n",
			[]string{"m.retry#0#small", "m.retry#1#small"}},
		{"step on attempt, in on mode (in sorts AFTER the step)",
			head(`eq = "retry"`) + "[rule.match.mode]\nin = [\"fast\", \"slow\"]\n[rule.guard.all.attempt]\nlt = 1\n" +
				"[rule.write]\nattempt = { step = 1 }\n",
			[]string{"m.retry#0#fast", "m.retry#0#slow"}},
		{"step on tier, in on mode (in sorts BEFORE the step)",
			head(`eq = "retry"`) + "[rule.match.mode]\nin = [\"fast\", \"slow\"]\n[rule.guard.all.tier]\nin = [\"small\"]\n" +
				"[rule.write]\ntier = { step = 1 }\n",
			[]string{"m.retry#fast#small", "m.retry#slow#small"}},
		{"step on attempt, outcome in (outcome sorts AFTER)",
			head(`in = ["retry", "escalate"]`) + "[rule.guard.all.attempt]\nlt = 1\n" +
				"[rule.write]\nattempt = { step = 1 }\n",
			[]string{"m.retry#0#escalate", "m.retry#0#retry"}},
		{"step on tier, outcome in (outcome sorts BEFORE)",
			head(`in = ["retry", "escalate"]`) + "[rule.guard.all.tier]\nin = [\"small\"]\n" +
				"[rule.write]\ntier = { step = 1 }\n",
			[]string{"m.retry#escalate#small", "m.retry#retry#small"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := s30MustLoad(t, tc.name, s30Src(s30Opts{decls: mode, keys: []string{"mode"}, rules: tc.rule}))
			if got := s30RowIDs(m, "retry"); !slices.Equal(got, tc.want) {
				t.Errorf("identities = %v; want %v (suffix elements ordered by key)", got, tc.want)
			}
		})
	}
}

// REQ-27: "BOUND. At every admitted cell the stepped literal MUST conform to
// the tag's declaration exactly as an authored literal write does … A cell
// whose stepped value leaves the domain in either direction — past `max` or
// below `min`, past the last member or before the first — is a load refusal
// under `malformed_tag_declaration` … whose detail names the rule, the tag,
// the cell, and the result."
// BOUNDARY
func TestReq27_0030_ASteppedValueLeavingTheDomainEitherWayIsRefused(t *testing.T) {
	for _, tc := range []struct {
		name, refused, admitted string
		names                   []string
	}{
		{"past max", s30Retry("gte = 0", "attempt = { step = 1 }"),
			s30Retry("gte = 0\nlt = 9", "attempt = { step = 1 }"), []string{"retry", "attempt", "9", "10"}},
		{"below min", s30Retry("lt = 5", "attempt = { step = -1 }"),
			s30Retry("gte = 1\nlt = 5", "attempt = { step = -1 }"), []string{"retry", "attempt", "0", "-1"}},
		{"past the last member", s30Escalate("", "tier = { step = 1 }"),
			s30Escalate(`in = ["small", "mid"]`, "tier = { step = 1 }"), []string{"escalate", "tier", "large"}},
		{"before the first member", s30Escalate("", "tier = { step = -1 }"),
			s30Escalate(`in = ["mid", "large"]`, "tier = { step = -1 }"), []string{"escalate", "tier", "small"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := s30Refusal(t, tc.name, s30Src(s30Opts{rules: tc.refused}))
			s30RequireCategory(t, tc.name, f, table.CatMalformedTagDeclaration)
			s30Contains(t, tc.name, f.Detail, tc.names...)
			s30MustLoad(t, "the in-domain twin", s30Src(s30Opts{rules: tc.admitted}))
		})
	}
}

// REQ-27 (0030:C2 with 0002:C11): "At every admitted cell the stepped literal
// MUST conform to the tag's declaration exactly as an authored literal write
// does", and an authored literal write of `<clear>` is refused as a
// `reserved tag value`. A domain may declare the sentinel as a member, so a
// step landing on it refuses under the literal write's category, naming the
// cell, rather than minting a write every consumer reads as removal.
// ADVERSARIAL
func TestReq27_0030_AStepLandingOnTheClearSentinelIsRefusedAsTheLiteralIs(t *testing.T) {
	tier := func(domain string) string {
		return "[tags.tier]\nprovenance = \"owned\"\nkind = \"enum\"\ndomain = " + domain +
			"\nsingle_valued = true\nrequired = true\n"
	}
	sentinel := tier(`["small", "<clear>"]`)

	f := s30Refusal(t, "a step onto <clear>",
		s30Src(s30Opts{tier: sentinel, rules: s30Escalate(`eq = "small"`, "tier = { step = 1 }")}))
	s30RequireCategory(t, "a step onto <clear>", f, table.CatReservedTagValue)
	s30Contains(t, "a step onto <clear>", f.Detail, "escalate", "tier", "cell small", table.ClearSentinel)

	// The literal twin refuses under the same category.
	lit := s30Refusal(t, "a literal <clear> write",
		s30Src(s30Opts{tier: sentinel, rules: s30Escalate(`eq = "small"`, `tier = "<clear>"`)}))
	s30RequireCategory(t, "a literal <clear> write", lit, table.CatReservedTagValue)

	// A step over a sentinel-carrying domain that never lands on it loads.
	s30MustLoad(t, "a step onto mid beside a <clear> member",
		s30Src(s30Opts{tier: tier(`["small", "mid", "<clear>"]`),
			rules: s30Escalate(`eq = "small"`, "tier = { step = 1 }")}))
}

// REQ-28: "for `int` it is the computed literal (`10`), which
// `internal/table/load.go::conformDomain` already reports as \"is above max
// 9\" / \"is below min 0\"; for `enum` there IS no stepped value — no member
// sits `n` positions from the cell — so the slot names the signed offset off
// the end instead (`1 past the last member`, `2 before the first`), and on
// the checked-add arm (A13) there is no representable literal either, so the
// slot names the bound it passed."
// BOUNDARY
func TestReq28_0030_TheResultSlotIsTheLiteralForIntAndTheOffsetForEnum(t *testing.T) {
	for _, tc := range []struct {
		name, rule string
		names      []string
	}{
		{"int above max", s30Retry("gte = 0", "attempt = { step = 1 }"), []string{"10", "is above max 9"}},
		{"int below min", s30Retry("lt = 5", "attempt = { step = -1 }"), []string{"-1", "is below min 0"}},
		{"enum one past", s30Escalate("", "tier = { step = 1 }"), []string{"large", "1 past the last member"}},
		{"enum two past", s30Escalate(`in = ["mid", "large"]`, "tier = { step = 2 }"),
			[]string{"mid", "1 past the last member"}},
		{"enum two before", s30Escalate(`in = ["small"]`, "tier = { step = -2 }"),
			[]string{"small", "2 before the first member"}},
		{"checked add names the bound", s30Retry("gte = 1\nlt = 2", "attempt = { step = 9223372036854775807 }"),
			[]string{"max"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := s30Refusal(t, tc.name, s30Src(s30Opts{rules: tc.rule}))
			s30RequireCategory(t, tc.name, f, table.CatMalformedTagDeclaration)
			s30Contains(t, tc.name, f.Detail, tc.names...)
		})
	}
}

// REQ-29: "All four parts ride the refusal's `Detail` STRING" … "This record
// adds none" … "Tests therefore assert the `Category` plus the detail text,
// and the detail names the cell FIRST so the assertion is anchored at a
// stable position rather than mid-sentence." Includes: `Failure.Rule`,
// `Offending` and `Remedy` stay unpopulated on this refusal.
// ADVERSARIAL
func TestReq29_0030_TheBoundRefusalRidesDetailAloneAndNamesTheCellFirst(t *testing.T) {
	f := s30Refusal(t, "the unguarded retry", s30Src(s30Opts{rules: s30Retry("gte = 0", "attempt = { step = 1 }")}))
	s30RequireCategory(t, "bound", f, table.CatMalformedTagDeclaration)
	if f.Rule != "" || f.Offending != "" || f.Remedy != "" {
		t.Errorf("the bound refusal populates a structured slot (Rule %q, Offending %q, Remedy %q); "+
			"all four parts ride Detail", f.Rule, f.Offending, f.Remedy)
	}
	if got := s30FirstIntAfter(f.Detail, "attempt"); got != "9" {
		t.Errorf("the first number after the tag key is %q; want the cell 9 named FIRST: %s", got, f.Detail)
	}
	if i, j := strings.Index(f.Detail, "attempt"), strings.LastIndex(f.Detail, "10"); i < 0 || j < i {
		t.Errorf("the result 10 is not named after the key and cell: %s", f.Detail)
	}
}

// REQ-30: "Where more than one admitted cell fails, the detail names the
// FIRST in the tag's own domain order — ascending `min..max` for `int`, the
// authored `domain` order for `enum` — not the row's suffix order"
// BOUNDARY
func TestReq30_0030_TheFirstFailingCellInDomainOrderIsNamed(t *testing.T) {
	attempt := "[tags.attempt]\nprovenance = \"owned\"\nkind = \"int\"\nmin = 0\nmax = 11\n" +
		"single_valued = true\nrequired = true\n"
	// Cells 9, 10, 11 all fail; "10" sorts first lexically, 9 in the domain.
	f := s30Refusal(t, "three failing int cells", s30Src(s30Opts{attempt: attempt,
		rules: s30Retry("gte = 9", "attempt = { step = 3 }")}))
	s30RequireCategory(t, "int multi-failure", f, table.CatMalformedTagDeclaration)
	if got := s30FirstIntAfter(f.Detail, "attempt"); got != "9" {
		t.Errorf("the named cell is %q; want 9 (domain-first), not 10 (suffix-lexical-first): %s", got, f.Detail)
	}

	// mid and large both fail at step 2; "large" sorts first lexically.
	g := s30Refusal(t, "two failing enum cells", s30Src(s30Opts{
		rules: s30Escalate(`in = ["mid", "large"]`, "tier = { step = 2 }")}))
	s30RequireCategory(t, "enum multi-failure", g, table.CatMalformedTagDeclaration)
	mid, large := strings.Index(g.Detail, "mid"), strings.Index(g.Detail, "large")
	if mid < 0 || (large >= 0 && large < mid) {
		t.Errorf("the detail does not name mid (authored-domain-first) ahead of large: %s", g.Detail)
	}
}

// REQ-31: "There is no saturating and no wrapping form. The stepped value is
// computed as a CHECKED `int` add, and that check is ordered BEFORE the
// render: a step whose magnitude carries a cell outside `int` range is
// refused as this bound failure, naming the cell, rather than being rendered
// and handed to conformance."
// ADVERSARIAL
func TestReq31_0030_TheCheckedAddRefusesAnOverflowingCellBeforeRender(t *testing.T) {
	f := s30Refusal(t, "step MaxInt from cell 1", s30Src(s30Opts{
		rules: s30Retry("gte = 1\nlt = 2", "attempt = { step = 9223372036854775807 }")}))
	s30RequireCategory(t, "overflow", f, table.CatMalformedTagDeclaration)
	if got := s30FirstIntAfter(f.Detail, "attempt"); got != "1" {
		t.Errorf("the named cell is %q; want 1: %s", got, f.Detail)
	}
	s30Contains(t, "overflow", f.Detail, "max")
	for _, bad := range []string{"is not an int", "-9223372036854775808"} {
		if strings.Contains(f.Detail, bad) {
			t.Errorf("the detail carries %q — the value was rendered (or wrapped) and handed to "+
				"conformance before the checked add: %s", bad, f.Detail)
		}
	}
	// Twin: the same admitted cell with a unit step loads.
	s30MustLoad(t, "the unit-step twin", s30Src(s30Opts{rules: s30Retry("gte = 1\nlt = 2", "attempt = { step = 1 }")}))
}

// REQ-32: "Because `unless` atoms are not consulted when cells are admitted
// (C1), the detail also names any `unless` atom the rule authors on the
// stepped tag and says it did not exclude the cell" With "surfaces as C2's
// bound refusal, whose detail names the unconsulted `unless` atom and directs
// the author to the positive-atom (`guard.all`) rewrite"
// ADVERSARIAL
func TestReq32_0030_TheBoundRefusalNamesAnUnconsultedUnlessAtom(t *testing.T) {
	unless := func(guard string) string {
		return strings.Replace(s30Retry(guard, "attempt = { step = 1 }"), "[rule.write]",
			"[rule.guard.unless.attempt]\neq = 9\n[rule.write]", 1)
	}
	f := s30Refusal(t, "an unless meant to cap the ladder", s30Src(s30Opts{rules: unless("gte = 0")}))
	s30RequireCategory(t, "unless", f, table.CatMalformedTagDeclaration)
	s30Contains(t, "unless", f.Detail, "unless", "guard.all")

	// Twin: the positive rewrite of the same exclusion loads.
	s30MustLoad(t, "the positive guard.all rewrite", s30Src(s30Opts{rules: unless("lt = 9")}))
}

// REQ-42: "the loader takes the CANONICALIZING form, since a set literal's two
// spellings are one literal"
// INPUT EDGE
func TestReq42_0030_ASetLiteralsTwoSpellingsAdmitTheSameCells(t *testing.T) {
	a := s30MustLoad(t, "in = [3, 1]", s30Src(s30Opts{rules: s30Retry("in = [3, 1]", "attempt = { step = 1 }")}))
	b := s30MustLoad(t, "in = [1, 3]", s30Src(s30Opts{rules: s30Retry("in = [1, 3]", "attempt = { step = 1 }")}))
	want := []string{"m.retry#1", "m.retry#3"}
	for _, m := range []*table.Model{a, b} {
		if got := s30RowIDs(m, "retry"); !slices.Equal(got, want) {
			t.Errorf("rows = %v; want %v", got, want)
		}
	}
	if !reflect.DeepEqual(rowsByRuleID(a, "retry"), rowsByRuleID(b, "retry")) {
		t.Errorf("two spellings of one set literal normalize to different rows:\n  %+v\n  %+v",
			rowsByRuleID(a, "retry"), rowsByRuleID(b, "retry"))
	}
}

// REQ-46: "No map iteration order reaches the row order." Observable: loading
// the same stepped model repeatedly yields byte-identical `dump` output and
// row order.
// ADVERSARIAL
func TestReq46_0030_RepeatedLoadsOfASteppedModelAreByteIdentical(t *testing.T) {
	two := s30Src(s30Opts{rules: `[[rule]]
id = "both"
use = ["open"]
[rule.match.recognized]
eq = "retry"
[rule.guard.all.attempt]
lt = 4
[rule.guard.all.tier]
in = ["small", "mid"]
[rule.write]
tier = { step = 1 }
attempt = { step = 1 }
`})
	for _, tc := range []struct{ name, src string }{
		{"ladder-step.toml", string(readFixture(t, s30LadderStep))},
		{"two stepped keys", two},
	} {
		t.Run(tc.name, func(t *testing.T) {
			first := s30MustLoad(t, tc.name, tc.src)
			dump, ids := table.Dump(first), rowIdentities(first)
			for range 25 {
				m := s30MustLoad(t, tc.name, tc.src)
				if got := table.Dump(m); got != dump {
					t.Fatalf("two loads dump differently:\n%s\n---\n%s", dump, got)
				}
				if got := rowIdentities(m); !slices.Equal(got, ids) {
					t.Fatalf("two loads order rows differently: %v vs %v", ids, got)
				}
			}
		})
	}
}

// REQ-50: "`ladder-step-unguarded.toml` — the `lt 5` atom removed — loaded.
// **Expected**: refusal under `malformed_tag_declaration` naming rule
// `retry`, tag `attempt`, cell `9`, result `10`; the enum sibling names
// `tier`, `large`, and `1 past the last member`" … "3c. … (i) refuses, naming
// cell `0` and `min` on the `int` arm and `1 before the first member` on the
// `enum` arm … (ii) loads, and its rows are `retry#1 … retry#9` each writing
// `cell - 1`."
// BOUNDARY
func TestReq50_0030_TheUnguardedLadderRefusesNamingRuleTagCellAndResult(t *testing.T) {
	for _, tc := range []struct {
		rel   string
		names []string
		cell  string
	}{
		{s30LadderUnguarded, []string{"retry", "attempt", "10"}, "9"},
		{s30LadderUnguardedE, []string{"tier", "large", "1 past the last member"}, ""},
	} {
		t.Run(tc.rel, func(t *testing.T) {
			f := s30Refusal(t, tc.rel, string(readFixture(t, tc.rel)))
			s30RequireCategory(t, tc.rel, f, table.CatMalformedTagDeclaration)
			s30Contains(t, tc.rel, f.Detail, tc.names...)
			if tc.cell != "" {
				if got := s30FirstIntAfter(f.Detail, "attempt"); got != tc.cell {
					t.Errorf("the first number after the key is %q; want the cell %s: %s", got, tc.cell, f.Detail)
				}
			}
		})
	}

	// 3c (i): a negative step with no excluding atom.
	i := s30Refusal(t, "step -1 over 0..9", s30Src(s30Opts{rules: s30Retry("", "attempt = { step = -1 }")}))
	s30RequireCategory(t, "3c(i) int", i, table.CatMalformedTagDeclaration)
	if got := s30FirstIntAfter(i.Detail, "attempt"); got != "0" {
		t.Errorf("3c(i): the named cell is %q; want 0: %s", got, i.Detail)
	}
	s30Contains(t, "3c(i) int", i.Detail, "min")
	ie := s30Refusal(t, "enum step -1", s30Src(s30Opts{rules: s30Escalate("", "tier = { step = -1 }")}))
	s30RequireCategory(t, "3c(i) enum", ie, table.CatMalformedTagDeclaration)
	s30Contains(t, "3c(i) enum", ie.Detail, "1 before the first member")

	// 3c (ii): the first cell excluded — loads, writing cell - 1.
	m := s30MustLoad(t, "step -1 with gte = 1", s30Src(s30Opts{rules: s30Retry("gte = 1", "attempt = { step = -1 }")}))
	var want []string
	for c := 1; c <= 9; c++ {
		want = append(want, "m.retry#"+strconv.Itoa(c))
	}
	slices.Sort(want)
	if got := s30RowIDs(m, "retry"); !slices.Equal(got, want) {
		t.Errorf("3c(ii) rows = %v; want %v", got, want)
	}
	for _, r := range rowsByRuleID(m, "retry") {
		c, _ := strconv.Atoi(r.Suffix[0])
		if got := s30Value(t, r.Identity(), r.Writes, "attempt"); got != strconv.Itoa(c-1) {
			t.Errorf("%s writes %s; want %d", r.Identity(), got, c-1)
		}
	}
	me := s30MustLoad(t, "enum step -1 over mid, large",
		s30Src(s30Opts{rules: s30Escalate(`in = ["mid", "large"]`, "tier = { step = -1 }")}))
	for cell, prev := range map[string]string{"mid": "small", "large": "mid"} {
		r := s30Row(t, me, "m.escalate#"+cell)
		if got := s30Value(t, r.Identity(), r.Writes, "tier"); got != prev {
			t.Errorf("%s writes %s; want %s", r.Identity(), got, prev)
		}
	}
}
