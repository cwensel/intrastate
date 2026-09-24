package table_test

// RDR 0030 `0030:C1` — the step write GRAMMAR and the kinds it is admitted
// on, asserted at the loader (`table.Load`), where every refusal this record
// mints fires.

import (
	"strconv"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// s30AdmittedRetry is the minimal admitted step model: `retry` steps
// `attempt` by one below the cap 5, so it expands to cells 0..4.
func s30AdmittedRetry() string {
	return s30Src(s30Opts{rules: s30Retry("lt = 5", "attempt = { step = 1 }")})
}

// s30RequireAdmittedRetry loads the admitted twin and asserts its five rows.
func s30RequireAdmittedRetry(t *testing.T) {
	t.Helper()

	m := s30MustLoad(t, "the admitted `{ step = 1 }` twin", s30AdmittedRetry())
	want := []string{"m.retry#0", "m.retry#1", "m.retry#2", "m.retry#3", "m.retry#4"}
	if got := s30RowIDs(m, "retry"); !slices0030Equal(got, want) {
		t.Errorf("the admitted twin expands to %v; want %v", got, want)
	}
}

func slices0030Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// REQ-1: "STEP WRITE. A rule's write block MAY assign an owned tag the inline
// table `{ step = <n> }`, `n` a non-zero TOML integer, and no other table
// shape: any other key set, a zero step, or a non-integer step (a TOML float
// `1.0` included) is refused at load under `malformed_tag_declaration` with a
// detail naming the one admitted form."
// INPUT EDGE
func TestReq1_0030_OnlyTheOneKeyNonZeroIntegerStepTableIsAdmitted(t *testing.T) {
	for _, tc := range []struct{ name, value string }{
		{"reporter's add key", "{ add = 1 }"},
		{"reporter's next key", "{ next = 1 }"},
		{"step plus a second key", "{ step = 1, add = 1 }"},
		{"empty table", "{}"},
		{"zero step", "{ step = 0 }"},
		{"float step 1.0", "{ step = 1.0 }"},
		{"float step 0.5", "{ step = 0.5 }"},
		{"string step", `{ step = "1" }`},
		{"bool step", "{ step = true }"},
		{"array step", "{ step = [1] }"},
		{"nested table step", "{ step = { n = 1 } }"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			src := s30Src(s30Opts{rules: s30Retry("lt = 5", "attempt = "+tc.value)})
			f := s30Refusal(t, tc.value, src)
			s30RequireCategory(t, tc.value, f, table.CatMalformedTagDeclaration)
			// The detail names the ONE admitted form, whose key is `step`.
			s30Contains(t, tc.value, f.Detail, "step")
		})
	}

	// The one admitted form loads — the twin every refusal above differs from.
	s30RequireAdmittedRetry(t)
}

// REQ-2: "The form is admitted only on a tag whose declared kind is `int`
// with both `min` and `max`, or `enum` with a non-empty `domain`; on `bool`,
// `set`, `scalar`, or an `int` missing a bound it is refused under the same
// category." With "load refusal, `malformed_tag_declaration`, detail names
// the admitted kinds."
// DOMAIN EDGE
func TestReq2_0030_StepIsAdmittedOnlyOnBoundedIntAndNonEmptyEnum(t *testing.T) {
	stepX := func(decl, guard string) string {
		rule := "[[rule]]\nid = \"bump\"\nuse = [\"open\"]\n[rule.match.recognized]\neq = \"succeed\"\n"
		if guard != "" {
			rule += "[rule.guard.all.x]\n" + guard + "\n"
		}
		rule += "[rule.write]\nx = { step = 1 }\n"
		return s30Src(s30Opts{decls: "[tags.x]\nprovenance = \"owned\"\n" + decl, keys: []string{"x"}, rules: rule})
	}

	for _, tc := range []struct{ name, decl string }{
		{"bool", "kind = \"bool\"\n"},
		{"set", "kind = \"set\"\nelements = [\"a\", \"b\"]\n"},
		{"scalar", "kind = \"scalar\"\n"},
		{"int missing max", "kind = \"int\"\nmin = 0\nsingle_valued = true\n"},
		{"int missing min", "kind = \"int\"\nmax = 9\nsingle_valued = true\n"},
		{"int missing both bounds", "kind = \"int\"\nsingle_valued = true\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := s30Refusal(t, "a step on "+tc.name, stepX(tc.decl, ""))
			s30RequireCategory(t, tc.name, f, table.CatMalformedTagDeclaration)
			s30Contains(t, tc.name, f.Detail, "int", "enum")
		})
	}

	// The admitted kinds load (the twins that make each refusal specific).
	s30MustLoad(t, "a step on a bounded int",
		stepX("kind = \"int\"\nmin = 0\nmax = 9\nsingle_valued = true\n", "lt = 9"))
	s30MustLoad(t, "a step on a non-empty enum",
		stepX("kind = \"enum\"\ndomain = [\"a\", \"b\"]\nsingle_valued = true\n", "eq = \"a\""))
}

// REQ-3: "An `int` whose declared width is not representable — `max - min`
// negative, or equal to `math.MaxInt`, so the inclusive `+1` wraps — is
// refused too" … "That is `internal/guard/declaration.go::intWidth`'s rule,
// and it must be the SAME rule: a declaration whose width lint cannot compute
// … must not be one the loader enumerates."
// BOUNDARY
func TestReq3_0030_AnUnrepresentableIntWidthIsRefusedForAStep(t *testing.T) {
	const maxInt = "9223372036854775807"
	const minInt = "-9223372036854775808"
	decl := func(minV, maxV string) string {
		return "[tags.x]\nprovenance = \"owned\"\nkind = \"int\"\nmin = " + minV +
			"\nmax = " + maxV + "\nsingle_valued = true\n"
	}
	model := func(minV, maxV, guard, write string) string {
		rule := "[[rule]]\nid = \"bump\"\nuse = [\"open\"]\n[rule.match.recognized]\neq = \"succeed\"\n"
		if guard != "" {
			rule += "[rule.guard.all.x]\n" + guard + "\n"
		}
		rule += "[rule.write]\nx = " + write + "\n"
		return s30Src(s30Opts{decls: decl(minV, maxV), keys: []string{"x"}, rules: rule})
	}

	for _, tc := range []struct{ name, minV, maxV string }{
		// max - min == math.MaxInt: the inclusive +1 wraps. NEW with this record.
		{"max - min == MaxInt from zero", "0", maxInt},
		{"max - min == MaxInt ending at -1", minInt, "-1"},
		// max - min overflows negative.
		{"max - min overflows", "-1", maxInt},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Negative control: the SAME declaration with a literal write, so
			// the refusal is the step's and not the declaration's. The literal
			// is the declared min, which every such domain contains.
			s30MustLoad(t, "the same declaration written literally", model(tc.minV, tc.maxV, "", tc.minV))

			f := s30Refusal(t, "a step over "+tc.name, model(tc.minV, tc.maxV, "", "{ step = 1 }"))
			s30RequireCategory(t, tc.name, f, table.CatMalformedTagDeclaration)
		})
	}

	// The negative-width arm is the pre-existing `min exceeds max` refusal,
	// restated: it fires before any write is read, with or without a step.
	for _, write := range []string{"0", "{ step = 1 }"} {
		f := s30Refusal(t, "min 5 max 0 with write "+write, model("5", "0", "", write))
		s30RequireCategory(t, "negative width", f, table.CatMalformedTagDeclaration)
	}

	// The representable twin: the same rule over max = 9 loads.
	s30MustLoad(t, "a step over the representable width 0..9", model("0", "9", "lt = 9", "{ step = 1 }"))
}

// REQ-6: "The bound is representability, not size — a wide but
// representable domain is admitted and merely expensive (A7), and the record
// takes no load-time row ceiling."
// BOUNDARY
func TestReq6_0030_AWideRepresentableDomainLoadsWithNoRowCeiling(t *testing.T) {
	// Wider than any lint bound in the tree (the product bound is 2048, the
	// node ceiling 4096), so a load-time ceiling borrowed from either would
	// refuse it.
	const width = 10000
	decl := "[tags.x]\nprovenance = \"owned\"\nkind = \"int\"\nmin = 0\nmax = " +
		strconv.Itoa(width-1) + "\nsingle_valued = true\n"
	rule := "[[rule]]\nid = \"bump\"\nuse = [\"open\"]\n[rule.match.recognized]\neq = \"succeed\"\n" +
		"[rule.guard.all.x]\nlt = " + strconv.Itoa(width-1) + "\n[rule.write]\nx = { step = 1 }\n"

	m := s30MustLoad(t, "a representable 10000-cell step",
		s30Src(s30Opts{decls: decl, keys: []string{"x"}, rules: rule}))
	if got := len(rowsByRuleID(m, "bump")); got != width-1 {
		t.Errorf("the wide step expands to %d rows; want %d, one per admitted cell", got, width-1)
	}
}

// REQ-7: "A step write over an `enum` whose `domain` repeats a member is
// refused as well" and "A step write over a domain containing a member with
// the suffix separator `#` is refused at load."
// DOMAIN EDGE
func TestReq7_0030_AStepOverADuplicateOrHashCarryingDomainIsRefused(t *testing.T) {
	tierDecl := func(domain string) string {
		return "[tags.tier]\nprovenance = \"owned\"\nkind = \"enum\"\ndomain = " + domain +
			"\nsingle_valued = true\nrequired = true\n"
	}
	stepTier := func(domain string) string {
		return s30Src(s30Opts{tier: tierDecl(domain),
			rules: s30Escalate(`in = ["small"]`, "tier = { step = 1 }")})
	}

	for _, tc := range []struct{ name, domain string }{
		{"repeated member", `["small", "mid", "small"]`},
		{"repeated member elsewhere in the domain", `["small", "mid", "large", "mid"]`},
		{"# member (not an admitted cell)", `["small", "mid", "la#rge"]`},
		{"# member (the stepped-to member)", `["small", "m#id", "large"]`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := s30Refusal(t, "a step over "+tc.name, stepTier(tc.domain))
			s30RequireCategory(t, tc.name, f, table.CatMalformedTagDeclaration)
		})
	}
	// The clean-domain twin loads.
	s30MustLoad(t, "a step over a clean domain", stepTier(`["small", "mid", "large"]`))
}

// REQ-8: "**Both guards are step-scoped, not new declaration rules**: an
// `enum` whose domain repeats a member or carries a `#` member and that NO
// rule steps continues to load exactly as it does today … The refusal fires
// on the step write, not on the declaration … this record changes no
// behaviour on a model that steps nothing."
// DOMAIN EDGE
func TestReq8_0030_DuplicateAndHashDomainsStillLoadWhenTheyAreNotStepped(t *testing.T) {
	for _, domain := range []string{`["small", "mid", "small"]`, `["small", "mid", "la#rge"]`} {
		t.Run(domain, func(t *testing.T) {
			tier := "[tags.tier]\nprovenance = \"owned\"\nkind = \"enum\"\ndomain = " + domain +
				"\nsingle_valued = true\nrequired = true\n"

			// Nothing steps: loads exactly as today.
			s30MustLoad(t, "an unstepped enum over "+domain,
				s30Src(s30Opts{tier: tier, rules: s30Escalate(`eq = "small"`, `tier = "mid"`)}))

			// The guard keys on the STEPPED tag: a step on a different, clean
			// tag in the same model loads beside the unstepped defect-shaped
			// domain.
			m := s30MustLoad(t, "a step on attempt beside the unstepped "+domain,
				s30Src(s30Opts{tier: tier, rules: s30Escalate(`eq = "small"`, `tier = "mid"`) + "\n" +
					s30Retry("lt = 5", "attempt = { step = 1 }")}))
			if got := len(rowsByRuleID(m, "retry")); got != 5 {
				t.Errorf("the attempt step expands to %d rows; want 5", got)
			}
		})
	}
}

// REQ-9: "`[initial]` values and predicate literals keep the literal-only
// grammar; the table shape is admitted on the write-block path alone. Each
// refusing path keeps its OWN existing category, and this record merges none
// of them: `[initial]` refuses under `malformed_initial_declaration` …, a
// predicate literal under `malformed_predicate_atom` …, and the write-block
// path under `malformed_tag_declaration`."
// ADVERSARIAL
func TestReq9_0030_TheTableShapeIsWriteBlockOnlyAndEachPathKeepsItsCategory(t *testing.T) {
	initial := s30Src(s30Opts{initial: "status = \"open\"\nattempt = { step = 1 }\ntier = \"small\"\n",
		rules: s30Retry("lt = 5", "attempt = 1")})
	f := s30Refusal(t, "[initial] attempt = { step = 1 }", initial)
	s30RequireCategory(t, "[initial]", f, table.CatMalformedInitialDeclaration)

	for _, tc := range []struct{ name, rule string }{
		{"match eq", "[[rule]]\nid = \"retry\"\nuse = [\"open\"]\n[rule.match.recognized]\neq = \"retry\"\n" +
			"[rule.match.attempt]\neq = { step = 1 }\n[rule.write]\nattempt = 1\n"},
		{"guard.all eq", s30Retry("eq = { step = 1 }", "attempt = 1")},
		{"in member", s30Retry("in = [1, { step = 1 }]", "attempt = 1")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := s30Refusal(t, "a predicate literal { step = 1 } as "+tc.name, s30Src(s30Opts{rules: tc.rule}))
			s30RequireCategory(t, tc.name, f, table.CatMalformedPredicateAtom)
		})
	}
	// The write-block path admits it.
	s30RequireAdmittedRetry(t)
}

// s30Minted enumerates one fixture per refusal kind this record mints, each
// beside the admitted twin it differs from in exactly the defect.
func s30Minted() []struct{ name, src string } {
	intX := func(minV, maxV string) string {
		return "[tags.x]\nprovenance = \"owned\"\nkind = \"int\"\nmin = " + minV + "\nmax = " + maxV +
			"\nsingle_valued = true\n"
	}
	bumpX := func(guard, write string) string {
		rule := "[[rule]]\nid = \"bump\"\nuse = [\"open\"]\n[rule.match.recognized]\neq = \"succeed\"\n"
		if guard != "" {
			rule += "[rule.guard.all.x]\n" + guard + "\n"
		}
		return rule + "[rule.write]\nx = " + write + "\n"
	}
	tierDomain := func(domain string) string {
		return "[tags.tier]\nprovenance = \"owned\"\nkind = \"enum\"\ndomain = " + domain +
			"\nsingle_valued = true\nrequired = true\n"
	}
	return []struct{ name, src string }{
		{"grammar: foreign key", s30Src(s30Opts{rules: s30Retry("lt = 5", "attempt = { add = 1 }")})},
		{"grammar: zero step", s30Src(s30Opts{rules: s30Retry("lt = 5", "attempt = { step = 0 }")})},
		{"grammar: float step", s30Src(s30Opts{rules: s30Retry("lt = 5", "attempt = { step = 1.0 }")})},
		{"kind: bool", s30Src(s30Opts{decls: "[tags.x]\nprovenance = \"owned\"\nkind = \"bool\"\n",
			keys: []string{"x"}, rules: bumpX("", "{ step = 1 }")})},
		{"kind: unbounded int", s30Src(s30Opts{decls: "[tags.x]\nprovenance = \"owned\"\nkind = \"int\"\nmin = 0\n",
			keys: []string{"x"}, rules: bumpX("", "{ step = 1 }")})},
		{"kind: unrepresentable width", s30Src(s30Opts{decls: intX("0", "9223372036854775807"),
			keys: []string{"x"}, rules: bumpX("", "{ step = 1 }")})},
		{"zero admitted cells", s30Src(s30Opts{rules: s30Retry("lt = 0", "attempt = { step = 1 }")})},
		{"# member in the stepped domain", s30Src(s30Opts{tier: tierDomain(`["small", "mid", "la#rge"]`),
			rules: s30Escalate(`in = ["small"]`, "tier = { step = 1 }")})},
		{"repeated member in the stepped domain", s30Src(s30Opts{tier: tierDomain(`["small", "mid", "small"]`),
			rules: s30Escalate(`in = ["small"]`, "tier = { step = 1 }")})},
		{"bound: past max", s30Src(s30Opts{rules: s30Retry("gte = 0", "attempt = { step = 1 }")})},
		{"bound: below min", s30Src(s30Opts{rules: s30Retry("lt = 5", "attempt = { step = -1 }")})},
		{"bound: past the last enum member", s30Src(s30Opts{rules: s30Escalate("", "tier = { step = 1 }")})},
		{"bound: checked add overflows", s30Src(s30Opts{rules: s30Retry("gte = 1\nlt = 2",
			"attempt = { step = 9223372036854775807 }")})},
		{"collision: stepped and cleared", s30Src(s30Opts{rules: strings.Replace(
			s30Retry("lt = 5", "attempt = { step = 1 }"), "use = [\"open\"]\n",
			"use = [\"open\"]\nclear = [\"attempt\"]\n", 1)})},
	}
}

// REQ-10: "Every refusal THIS clause and its surfaces mint — the kind and
// grammar refusals, the zero-cell refusal, the `#`-member and
// duplicate-member refusals, and C2's bound — is on the write-block path and
// so takes `malformed_tag_declaration`."
// HAPPY PATH
func TestReq10_0030_EveryRefusalThisRecordMintsIsMalformedTagDeclaration(t *testing.T) {
	for _, tc := range s30Minted() {
		t.Run(tc.name, func(t *testing.T) {
			f := s30Refusal(t, tc.name, tc.src)
			s30RequireCategory(t, tc.name, f, table.CatMalformedTagDeclaration)
		})
	}
	// Every fixture above differs from this admitted step in its defect alone.
	s30RequireAdmittedRetry(t)
}

// REQ-19: "A rule admitting zero cells is refused at load." With "A rule with
// zero admitted cells is refused, and that refusal is NEW. The loader does
// not refuse an unsatisfiable guard today" — (TD) and "A step rule admitting
// no cell — a NEW load refusal".
// BOUNDARY
func TestReq19_0030_AStepRuleAdmittingNoCellIsRefused(t *testing.T) {
	// Negative control: the same unsatisfiable `lt = 0` WITHOUT a step loads,
	// today and after — the refusal is the step's, not 0003's.
	s30MustLoad(t, "lt = 0 with a literal write", s30Src(s30Opts{rules: s30Retry("lt = 0", "attempt = 3")}))

	for _, guard := range []string{"lt = 0", "gt = 9", "gte = 3\nlt = 3", "in = [7, 9]\nlt = 5"} {
		f := s30Refusal(t, "a zero-cell step ("+guard+")",
			s30Src(s30Opts{rules: s30Retry(guard, "attempt = { step = 1 }")}))
		s30RequireCategory(t, guard, f, table.CatMalformedTagDeclaration)
	}
	// One cell over: `lt = 1` admits exactly cell 0 and loads.
	m := s30MustLoad(t, "lt = 1 with a step", s30Src(s30Opts{rules: s30Retry("lt = 1", "attempt = { step = 1 }")}))
	if got := s30RowIDs(m, "retry"); !slices0030Equal(got, []string{"m.retry#0"}) {
		t.Errorf("lt = 1 expands to %v; want [m.retry#0]", got)
	}
}

// REQ-20: "A rule that BOTH steps a key and names it in its `clear` list is
// refused at load under the same category" … "the collision refusal PRECEDES
// the per-cell walk: a key both stepped and cleared is refused before any
// cell is evaluated, so it never reaches C2's bound. Where a rule would trip
// both, the collision is what the author is told"
// ADVERSARIAL
func TestReq20_0030_SteppingAndClearingOneKeyIsRefusedBeforeAnyCellIsWalked(t *testing.T) {
	withClear := func(rule, key string) string {
		return strings.Replace(rule, "use = [\"open\"]\n", "use = [\"open\"]\nclear = [\""+key+"\"]\n", 1)
	}

	f := s30Refusal(t, "attempt stepped and cleared",
		s30Src(s30Opts{rules: withClear(s30Retry("lt = 5", "attempt = { step = 1 }"), "attempt")}))
	s30RequireCategory(t, "collision", f, table.CatMalformedTagDeclaration)
	s30Contains(t, "collision", f.Detail, "clear")

	// Trips both: unguarded (cell 9 steps to 10) AND cleared. The collision
	// is what the author is told; no cell was walked, so no result is named.
	both := s30Refusal(t, "unguarded, stepped and cleared",
		s30Src(s30Opts{rules: withClear(s30Retry("gte = 0", "attempt = { step = 1 }"), "attempt")}))
	s30RequireCategory(t, "collision over bound", both, table.CatMalformedTagDeclaration)
	s30Contains(t, "collision over bound", both.Detail, "clear")
	if strings.Contains(both.Detail, "10") {
		t.Errorf("the collision refusal names the stepped result 10, so the per-cell walk "+
			"ran before the collision check: %s", both.Detail)
	}
	// Twin: clearing a DIFFERENT key beside the step loads.
	s30MustLoad(t, "a step on attempt clearing tier",
		s30Src(s30Opts{rules: withClear(s30Retry("lt = 5", "attempt = { step = 1 }"), "tier")}))
}
