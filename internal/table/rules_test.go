package table_test

// RDR 0002 sections F–H: rule shape, ids, and writes; escape rows; and
// shared contexts and their merge semantics.

import (
	"reflect"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/table"
)

// REQ-34: "Rule ids MUST be unique within a model, compared by exact byte
// equality; a duplicate rule id is a load failure."
// ADVERSARIAL
func TestReq34_RuleIDsAreUniqueByExactByteEquality(t *testing.T) {
	t.Run("duplicate refuses", func(t *testing.T) {
		got := loadCategory(t, "neg/neg-duplicate-rule-id.toml")
		if got != table.CatDuplicateRuleID {
			t.Errorf("category = %q; want %q", got, table.CatDuplicateRuleID)
		}
	})

	// "compared by exact byte equality" — two ids differing only in case
	// are two ids, not a duplicate. A normalizer folding case refuses a
	// legal model.
	t.Run("case-differing ids are not duplicates", func(t *testing.T) {
		base := string(readFixture(t, rdrFixture))
		src := strings.Replace(base, `id = "reconcile-rewind"`, `id = "Continue-Prelock"`, 1)
		if src == base {
			t.Fatal("rule id substitution did not apply")
		}
		m, err := table.Load([]byte(src), "case-differing-ids.toml")
		if err != nil {
			t.Fatalf("`Continue-Prelock` and `continue-prelock` are two ids, "+
				"not a duplicate; the document MUST load clean, refused: %v", err)
		}
		// Both ids survive normalization as distinct rows: neither was
		// folded into the other, and neither was dropped.
		for _, id := range []string{"Continue-Prelock", "continue-prelock"} {
			if len(rowsByRuleID(m, id)) == 0 {
				t.Errorf("no normalized row carries rule id %q; case-differing "+
					"ids were folded together", id)
			}
		}
	})
}

// REQ-35: "Each transition rule MUST contain a stable rule id, zero or
// more shared-context references, and a local match block. An ordinary
// transition rule MUST contain a write block, MAY contain a rule-level
// explicit clear list, and MAY contain a rule-level `gate` list. A write
// block MAY assign more than one tag."
// ADVERSARIAL
func TestReq35_OrdinaryRuleShape(t *testing.T) {
	t.Run("ordinary rule with no write block refuses", func(t *testing.T) {
		got := loadCategory(t, "neg/neg-no-write-block.toml")
		if got != table.CatMalformedRuleShape {
			t.Errorf("category = %q; want %q", got, table.CatMalformedRuleShape)
		}
	})

	// "and a local match block" — the "zero or more" allowance is spent on
	// the shared-context references, so the match-block obligation is
	// unqualified. A rule whose entire selection criterion is inherited
	// declares nothing at the rule site, which `0002:C3` makes a stable
	// refusal rather than a silent admission. Presence-keyed per D2: an
	// absent block and an explicit empty one both decode to nil.
	for name, rel := range map[string]string{
		"no match block":         "neg/neg-rule-no-match-block.toml",
		"empty match block":      "neg/neg-rule-empty-match-block.toml",
		"escape, no match block": "neg/neg-escape-no-match-block.toml",
	} {
		t.Run(name+" refuses", func(t *testing.T) {
			got := loadCategory(t, rel)
			if got != table.CatMalformedRuleShape {
				t.Errorf("category = %q; want %q", got, table.CatMalformedRuleShape)
			}
		})
	}

	t.Run("a write block MAY assign more than one tag", func(t *testing.T) {
		m := mustLoad(t, rdrFixture)
		row := rowByID(t, m, "rdr.reconcile-rewind")
		// Three authored writes plus one rendered clear.
		want := map[string][]string{
			"stage":        {"resolve"},
			"status":       {"Draft"},
			"rewind_scope": {"assumptions"},
			"prelock_lens": {table.ClearSentinel},
		}
		if len(row.Writes) != len(want) {
			t.Fatalf("multi-tag write yielded %d writes; want %d (%+v)",
				len(row.Writes), len(want), row.Writes)
		}
		for k, v := range want {
			got, ok := tagValue(row.Writes, k)
			if !ok {
				t.Errorf("write on %q is missing", k)
				continue
			}
			if !reflect.DeepEqual(got, v) {
				t.Errorf("write %q = %v; want %v", k, got, v)
			}
		}
	})

	t.Run("zero shared-context references is legal", func(t *testing.T) {
		base := string(readFixture(t, kataFixture))
		src := strings.Replace(base, "use = [\"review\"]\nsource = \"kata:review\"\n[rule.match.recognized]\neq = \"needs-work\"",
			"source = \"kata:review\"\n[rule.match.recognized]\neq = \"needs-work\"\n[rule.match.status]\neq = \"open\"", 1)
		if src == base {
			t.Fatal("use-list substitution did not apply")
		}
		if _, err := table.Load([]byte(src), "no-use.toml"); err != nil {
			t.Errorf("a rule with zero context references refused: %v", err)
		}
	})
}

// REQ-36: "An escape rule MUST contain an `escape` list and MUST NOT
// contain a write block, clear list, or `gate` list, even an empty one
// (RDR 0009 binds the write-free obligation at the kernel boundary)."
// ADVERSARIAL
//
// "even an empty one" is the discriminating clause: the loader keys on key
// PRESENCE, not on length (deviations.md D2 — the spike checked
// `len(rule.Write) > 0` and would pass a `write = []` escape rule).
func TestReq36_EscapeRuleCarriesNoPlanBearingField(t *testing.T) {
	cases := map[string]string{
		"write block": "neg/neg-escape-with-write.toml",
		"clear list":  "neg/neg-escape-with-clear.toml",
		"gate list":   "neg/neg-escape-with-gate.toml",
		// D2: the presence-keyed arm.
		"empty write block": "neg/neg-escape-with-empty-write.toml",
	}
	for name, rel := range cases {
		t.Run(name, func(t *testing.T) {
			if got := loadCategory(t, rel); got != table.CatMalformedEscapeDeclaration {
				t.Errorf("category = %q; want %q", got, table.CatMalformedEscapeDeclaration)
			}
		})
	}

	// Each fixture carries exactly one of the three shapes: the category
	// oracle alone cannot tell them apart, so the discriminating assertion
	// is that the other two shapes are absent from each fixture.
	t.Run("one shape per fixture", func(t *testing.T) {
		shapes := map[string][]string{
			"neg/neg-escape-with-write.toml":       {"clear = [", "gate = ["},
			"neg/neg-escape-with-clear.toml":       {"[rule.write]", "gate = ["},
			"neg/neg-escape-with-gate.toml":        {"[rule.write]", "clear = ["},
			"neg/neg-escape-with-empty-write.toml": {"clear = [", "gate = ["},
		}
		for rel, absent := range shapes {
			src := escapeRuleBlock(t, rel)
			for _, shape := range absent {
				if strings.Contains(src, shape) {
					t.Errorf("%s: the escape rule also carries %q; each fixture "+
						"must carry exactly one shape", rel, shape)
				}
			}
		}
	})
}

// escapeRuleBlock returns the text of the escape rule in a fixture, from
// its `id = "draft-no-match-escape"` line to the next `[[rule]]`.
func escapeRuleBlock(t *testing.T, rel string) string {
	t.Helper()

	src := string(readFixture(t, rel))
	i := strings.Index(src, `id = "draft-no-match-escape"`)
	if i < 0 {
		t.Fatalf("%s carries no escape rule", rel)
	}
	rest := src[i:]
	if j := strings.Index(rest, "[[rule]]"); j >= 0 {
		rest = rest[:j]
	}
	return rest
}

// REQ-37: "**A write replaces.** A write assigns a tag's whole value and
// supplants whatever was held; for a `set` kind the array literal is the
// whole new set. There is no accumulate form."
// BOUNDARY
//
// The checkable form of "replaces, no accumulate form" is that a set-kind
// write normalizes to exactly the authored members — the loader never
// unions with `[initial]` or with any other rule's write on the same key.
func TestReq37_AWriteReplaces(t *testing.T) {
	m := mustLoad(t, kataFixture)

	row := rowByID(t, m, "kata.review-needs-work")
	got, ok := tagValue(row.Writes, "labels")
	if !ok {
		t.Fatal("review-needs-work carries no `labels` write")
	}
	want := []string{"needs work"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("set-kind write = %v; want exactly the authored members %v — "+
			"a write replaces and there is no accumulate form", got, want)
	}

	// The sibling rule writes `status` too; neither accumulates the other.
	accepted := rowByID(t, m, "kata.review-accepted")
	if v, _ := tagValue(accepted.Writes, "status"); !reflect.DeepEqual(v, []string{"accepted"}) {
		t.Errorf("review-accepted status write = %v; want [accepted]", v)
	}
	if v, _ := tagValue(row.Writes, "status"); !reflect.DeepEqual(v, []string{"open"}) {
		t.Errorf("review-needs-work status write = %v; want [open]", v)
	}
}

// REQ-38: "a `set`-kind write value MUST normalize to an ordered,
// member-sorted sequence and MUST NOT be joined into a single
// delimiter-separated string, in the stored value or in any comparison
// derived from it."
// ADVERSARIAL
//
// The delimiter-parameterized control (Testing Strategy 2). A control
// using only scalar or space-bearing members does not discriminate: the
// pre-fix spike joined on `,` and collapsed ONLY the comma case, so four
// of five delimiters returned distinct values while the defect was live.
func TestReq38_SetWriteValueIsAMemberSequenceNotAJoinedString(t *testing.T) {
	delims := []string{"comma", "semi", "pipe", "space", "empty"}

	for _, d := range delims {
		t.Run(d, func(t *testing.T) {
			relA := "delim/write-delim-" + d + "-a.toml"
			relB := "delim/write-delim-" + d + "-b.toml"

			a, err := table.Load(readFixture(t, relA), relA)
			if err != nil {
				t.Fatalf("%s refused: %v", relA, err)
			}
			b, err := table.Load(readFixture(t, relB), relB)
			if err != nil {
				t.Fatalf("%s refused: %v", relB, err)
			}

			va := setWriteValue(t, a, "labels")
			vb := setWriteValue(t, b, "labels")

			// Read the NORMALIZED value, never the rendered one: Stage 6's
			// `|` case showed a render can still collide after the identity
			// is correct.
			if reflect.DeepEqual(va, vb) {
				t.Errorf("the two set spellings normalized to the same value %v; "+
					"they must be distinct — a joined rendering collapses them", va)
			}
			if len(va) != 2 || len(vb) != 2 {
				t.Errorf("write values are not two-member sequences: a=%v b=%v — "+
					"a joined value has length 1", va, vb)
			}
		})
	}

	t.Run("members are sorted", func(t *testing.T) {
		src := []byte(`outcomes = ["done"]

[model]
id = "m"
version = 1

[tags.labels]
provenance = "owned"
kind = "set"
elements = ["a", "b", "c"]

[tags.recognized]
provenance = "recognized"
kind = "enum"
required = true

[read.r]
role = "m"
path = "m"
keys = ["labels"]
timeout = "2s"

[write.w]
role = "m"
path = "m"
keys = ["labels"]
timeout = "2s"
read_back = true

[[rule]]
id = "r"
[rule.match.recognized]
eq = "done"
[rule.write]
labels = ["c", "a", "b"]
`)
		m, err := table.Load(src, "sorted-write.toml")
		if err != nil {
			t.Fatalf("refused: %v", err)
		}
		got := setWriteValue(t, m, "labels")
		want := []string{"a", "b", "c"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("set write value = %v; want the member-sorted sequence %v", got, want)
		}
	})
}

// setWriteValue returns the value the model's single labels-writing row
// assigns to key.
func setWriteValue(t *testing.T, m *table.Model, key string) []string {
	t.Helper()

	for _, r := range m.Rows {
		if v, ok := tagValue(r.Writes, key); ok {
			return v
		}
	}
	t.Fatalf("no normalized row writes %q", key)
	return nil
}

// REQ-39: "An `escape` list MUST contain only resolver failure classes
// that RDR 0001 allows the table to model: `no_match` and
// `ambiguous_match`."
// ADVERSARIAL
func TestReq39_EscapeListAdmitsOnlyTwoFailureClasses(t *testing.T) {
	base := string(readFixture(t, rdrFixture))

	t.Run("both admitted classes load", func(t *testing.T) {
		for _, class := range []string{"no_match", "ambiguous_match"} {
			src := strings.Replace(base, `escape = ["no_match"]`,
				`escape = ["`+class+`"]`, 1)
			if _, err := table.Load([]byte(src), "escape-"+class+".toml"); err != nil {
				t.Errorf("escape class %q refused: %v", class, err)
			}
		}
	})

	t.Run("a non-modelable class refuses", func(t *testing.T) {
		// The kernel's other three refusal kinds are never escapable.
		for _, class := range []string{"owned_state_unavailable", "guard_unevaluable", "unmodeled_outcome"} {
			src := strings.Replace(base, `escape = ["no_match"]`,
				`escape = ["`+class+`"]`, 1)
			_, err := table.Load([]byte(src), "escape-"+class+".toml")
			if err == nil {
				t.Errorf("escape class %q loaded clean; only no_match and "+
					"ambiguous_match may be modeled", class)
				continue
			}
			if cat, _ := table.CategoryOf(err); cat != table.CatMalformedEscapeDeclaration {
				t.Errorf("escape class %q: category = %q; want %q",
					class, cat, table.CatMalformedEscapeDeclaration)
			}
		}
	})
}

// REQ-40: "Normalization MUST render an escape rule as a candidate row
// carrying its normal predicate set, outcome, source rule id, source
// locator, and modeled failure class list."
// HAPPY PATH
func TestReq40_EscapeRuleNormalizesToACandidateRow(t *testing.T) {
	m := mustLoad(t, rdrFixture)
	row := rowByID(t, m, "rdr.draft-no-match-escape#round-clean")

	if row.RuleID != "draft-no-match-escape" {
		t.Errorf("source rule id = %q; want draft-no-match-escape", row.RuleID)
	}
	if row.SourceLocator == "" {
		t.Error("escape row carries no source locator")
	}
	if row.Outcome != "round-clean" {
		t.Errorf("outcome = %q; want round-clean", row.Outcome)
	}
	if !reflect.DeepEqual(row.Escape, []string{"no_match"}) {
		t.Errorf("modeled failure class list = %v; want [no_match]", row.Escape)
	}
	// Its normal predicate set: the inherited `draft` context's atom.
	want := []table.Atom{{
		Key: "status", Block: table.BlockMatch, Operator: "eq", Literal: []string{"Draft"},
	}}
	if !reflect.DeepEqual(row.Atoms, want) {
		t.Errorf("predicate set = %+v; want %+v", row.Atoms, want)
	}
}

// REQ-41: "An escape row rescues only resolves carrying the outcome it
// binds" / "There is no table-wide catch-all: covering an alphabet of N
// outcomes requires N escape rows."
// BOUNDARY
//
// The RDR fixture's alphabet has four members and its one escape rule
// binds two of them, so exactly two escape rows exist and neither carries
// a wildcard outcome. A catch-all implementation would produce one row
// with an empty or wildcard outcome.
func TestReq41_EscapeRowsBindOneOutcomeEachWithNoCatchAll(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	var escapes []table.Row
	for _, r := range m.Rows {
		if len(r.Escape) > 0 {
			escapes = append(escapes, r)
		}
	}
	if len(escapes) != 2 {
		t.Fatalf("%d escape rows; want 2 — one per bound alphabet member", len(escapes))
	}

	bound := map[string]bool{}
	for _, r := range escapes {
		if r.Outcome == "" {
			t.Errorf("escape row %s binds no outcome; there is no table-wide "+
				"catch-all", r.Identity())
		}
		bound[r.Outcome] = true
	}
	want := map[string]bool{"round-clean": true, "reconcile-block": true}
	if !reflect.DeepEqual(bound, want) {
		t.Errorf("escape rows bind %v; want %v", bound, want)
	}
	// The two unbound alphabet members are rescued by nothing.
	for _, o := range []string{"verdict-flapping", "finalized"} {
		if bound[o] {
			t.Errorf("outcome %q is rescued by an escape row the fixture does "+
				"not author", o)
		}
	}
}

// REQ-42: "An escape rule MAY bind its outcome with an `in` atom and
// expand like any other rule … each expansion is a separate row rescuing
// its own outcome, and the expansion suffix distinguishes them."
// HAPPY PATH (A11)
func TestReq42_EscapeRuleExpandsUnderIn(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	rows := rowsByRuleID(m, "draft-no-match-escape")
	if len(rows) != 2 {
		t.Fatalf("draft-no-match-escape yielded %d rows; want 2", len(rows))
	}

	suffixes := map[string]string{}
	for _, r := range rows {
		if len(r.Suffix) != 1 {
			t.Errorf("%s suffix = %v; want a one-element sequence", r.Identity(), r.Suffix)
			continue
		}
		suffixes[r.Suffix[0]] = r.Outcome
		if len(r.Writes) != 0 {
			t.Errorf("%s carries writes %+v; an escape row has an empty write set",
				r.Identity(), r.Writes)
		}
		if !reflect.DeepEqual(r.Escape, []string{"no_match"}) {
			t.Errorf("%s escape classes = %v; each expansion retains the modeled "+
				"failure-class list", r.Identity(), r.Escape)
		}
	}
	want := map[string]string{"round-clean": "round-clean", "reconcile-block": "reconcile-block"}
	if !reflect.DeepEqual(suffixes, want) {
		t.Errorf("expansion suffixes -> outcomes = %v; want %v", suffixes, want)
	}
}

// REQ-43: "Row kind (`transition` / `escape`) is a **derived view
// property, never a row field**: escape identity is discriminated solely
// by a non-empty escape class list"
// BOUNDARY
//
// The structural half: the normalized row type carries no kind field. The
// behavioural half: the derived kind agrees with the escape-list
// predicate on every row of both fixtures.
func TestReq43_RowKindIsDerivedNotAField(t *testing.T) {
	t.Run("the row type carries no kind field", func(t *testing.T) {
		rt := reflect.TypeOf(table.Row{})
		for i := range rt.NumField() {
			switch strings.ToLower(rt.Field(i).Name) {
			case "kind", "rowkind", "iskind":
				t.Errorf("table.Row carries a %q field; row kind is a derived "+
					"view property, never a row field", rt.Field(i).Name)
			}
		}
	})

	t.Run("kind is discriminated solely by a non-empty escape list", func(t *testing.T) {
		for _, rel := range []string{rdrFixture, kataFixture} {
			m := mustLoad(t, rel)
			for _, r := range m.Rows {
				want := table.KindTransition
				if len(r.Escape) > 0 {
					want = table.KindEscape
				}
				if got := r.Kind(); got != want {
					t.Errorf("%s: %s kind = %q; want %q (escape list %v)",
						rel, r.Identity(), got, want, r.Escape)
				}
			}
		}
	})
}

// REQ-44: "a render-time view struct MAY materialize the computed column …
// What it MUST NOT do is let that view feed back into the normalized value
// or the kernel row"
// ADVERSARIAL
//
// The dump renders the kind column, and the normalized value must be
// unchanged by rendering it. A view feeding back is caught by comparing
// the row set across a dump.
func TestReq44_RenderedKindDoesNotFeedBackIntoTheNormalizedValue(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	before := make([]table.Row, len(m.Rows))
	copy(before, m.Rows)

	out := table.Dump(m)
	if !strings.Contains(out, string(table.KindEscape)) {
		t.Errorf("dump does not render the derived kind column:\n%s", out)
	}

	if !reflect.DeepEqual(before, m.Rows) {
		t.Error("rendering the dump mutated the normalized candidate rows; the " +
			"view must not feed back into the normalized value")
	}

	// And the kernel row carries no kind either.
	kr := rowByID(t, m, "rdr.draft-no-match-escape#round-clean").KernelRow()
	if len(kr.Escape) != 1 {
		t.Errorf("kernel row escape list = %v; want one class", kr.Escape)
	}
}

// REQ-45: "Shared contexts MAY inherit from other contexts, but
// inheritance MUST normalize to an explicit predicate set before lint or
// resolution."
// HAPPY PATH
//
// `large-prelock` inherits `prelock` inherits `draft`, so a row using
// `large-prelock` carries all three contexts' atoms explicitly and no
// context reference survives on the row.
func TestReq45_InheritanceNormalizesToAnExplicitPredicateSet(t *testing.T) {
	m := mustLoad(t, rdrFixture)
	row := rowByID(t, m, "rdr.continue-prelock#large")

	// draft -> status.eq=Draft; prelock -> stage.eq=prelock;
	// large-prelock -> profile.in expanded to profile.eq=large.
	for _, want := range []table.Atom{
		{Key: "status", Block: table.BlockMatch, Operator: "eq", Literal: []string{"Draft"}},
		{Key: "stage", Block: table.BlockMatch, Operator: "eq", Literal: []string{"prelock"}},
		{Key: "profile", Block: table.BlockMatch, Operator: "eq", Literal: []string{"large"}},
	} {
		if !containsAtom(row.Atoms, want) {
			t.Errorf("inherited atom %+v is absent from the normalized set %+v",
				want, row.Atoms)
		}
	}

	// No bare context reference survives on the normalized row.
	rt := reflect.TypeOf(table.Row{})
	for i := range rt.NumField() {
		switch strings.ToLower(rt.Field(i).Name) {
		case "use", "contexts", "inherits":
			t.Errorf("table.Row carries a %q field; inheritance normalizes to an "+
				"explicit predicate set", rt.Field(i).Name)
		}
	}
}

func containsAtom(atoms []table.Atom, want table.Atom) bool {
	for _, a := range atoms {
		if reflect.DeepEqual(a, want) {
			return true
		}
	}
	return false
}

// REQ-46: "The combined predicate set is a **set over the full atom
// identity** — the tuple `(key, block, operator token, literal)`, the same
// tuple the dump sorts on. Merging MUST NOT key on any proper prefix of
// it: two atoms agreeing on `(key, block, operator)` but differing in
// literal are **distinct atoms** and both survive the merge."
// ADVERSARIAL
//
// The scalar witnesses answer the (key, block, operator) question only.
// The delimiter-bearing controls close the FULL-identity question, and
// they parameterize over five plausible join delimiters because a
// black-box test cannot read the delimiter a wrong implementation chose —
// the pre-fix spike collapsed only the comma case.
func TestReq46_MergeKeysOnTheFullAtomIdentity(t *testing.T) {
	t.Run("scalar witnesses: distinct literals both survive", func(t *testing.T) {
		// Both fixtures give `reconcile-rewind` two inherited contexts
		// contributing status.eq=Draft and status.eq=Final.
		for _, rel := range []string{"merge-distinct.toml", "merge-distinct-rev.toml"} {
			m := mustLoad(t, rel)
			row := rowByID(t, m, "rdr.reconcile-rewind")
			got := atomsOn(row, "status")
			if len(got) != 2 {
				t.Errorf("%s: %d atoms on `status`; want 2 — a merge keyed on "+
					"(key, block, operator) drops one: %+v", rel, len(got), got)
			}
		}
	})

	t.Run("use-order independence", func(t *testing.T) {
		a := mustLoad(t, "merge-distinct.toml")
		b := mustLoad(t, "merge-distinct-rev.toml")
		if !reflect.DeepEqual(a.Rows, b.Rows) {
			t.Error("the two `use` orders normalized to different row sets; the " +
				"merge result must not depend on `use` iteration order")
		}
	})

	t.Run("delimiter-parameterized full-identity control", func(t *testing.T) {
		for _, d := range []string{"comma", "semi", "pipe", "space", "empty"} {
			rel := "delim/merge-delim-" + d + ".toml"
			m, err := table.Load(readFixture(t, rel), rel)
			if err != nil {
				t.Fatalf("%s refused: %v", rel, err)
			}
			// Two contexts contribute `in = ["a,b", "c"]` and
			// `in = ["a", "b,c"]` on one key and block, so the rule expands
			// to four rows; each carries one atom per contributing context.
			if len(rowsByRuleID(m, "reconcile-rewind")) != 4 {
				t.Errorf("%s: reconcile-rewind yielded %d rows; want 4 (the "+
					"product of two two-member `in` atoms)",
					rel, len(rowsByRuleID(m, "reconcile-rewind")))
			}
			got := atomsOn(rowsByRuleID(m, "reconcile-rewind")[0], "finalized_at")
			if len(got) != 2 {
				t.Errorf("%s: %d atoms on `finalized_at`; want 2 — a merge keyed "+
					"on a %s-joined rendering collapses them: %+v", rel, len(got), d, got)
			}
		}
	})
}

// REQ-47: "Merging is idempotent on identical atoms: the same atom
// contributed by a rule and by one or more inherited contexts collapses to
// one. Inheritance therefore never overrides — it only accumulates."
// BOUNDARY (A13's idempotence mirror)
//
// Without this a normalizer that de-duplicates nothing passes every
// distinct-literal control above, since they only ever count the
// two-distinct-literal case.
func TestReq47_MergeIsIdempotentOnIdenticalAtoms(t *testing.T) {
	// `reconcile-rewind` authors status.eq=Draft locally AND inherits the
	// byte-identical atom from context `draft`.
	m := mustLoad(t, "merge-idempotent.toml")

	got := atomsOn(rowByID(t, m, "rdr.reconcile-rewind"), "status")
	if len(got) != 1 {
		t.Fatalf("%d atoms on `status`; want exactly 1 — a byte-identical atom "+
			"from a rule and an inherited context collapses to one: %+v", len(got), got)
	}
	want := table.Atom{Key: "status", Block: table.BlockMatch, Operator: "eq", Literal: []string{"Draft"}}
	if !reflect.DeepEqual(got[0], want) {
		t.Errorf("collapsed atom = %+v; want %+v", got[0], want)
	}
}

// REQ-48: "authoring two atoms on one key that no view can satisfy
// together yields a dead rule, which is RDR 0006's unreachable-rule
// finding, not a load failure here."
// INPUT EDGE
//
// Asserted positively: the document LOADS and both atoms survive. Stating
// it as "the rule is a dead rule for lint" would be an RDR 0006 verdict,
// which the MVV refuses to count against an unimplemented lint.
func TestReq48_UnsatisfiableAtomPairIsNotALoadFailure(t *testing.T) {
	m := mustLoad(t, "merge-distinct.toml")

	got := atomsOn(rowByID(t, m, "rdr.reconcile-rewind"), "status")
	if len(got) != 2 {
		t.Fatalf("%d atoms on `status`; want 2 surviving atoms: %+v", len(got), got)
	}
	literals := map[string]bool{}
	for _, a := range got {
		if len(a.Literal) != 1 {
			t.Errorf("atom %+v carries a non-single literal", a)
			continue
		}
		literals[a.Literal[0]] = true
	}
	if !reflect.DeepEqual(literals, map[string]bool{"Draft": true, "Final": true}) {
		t.Errorf("surviving literals = %v; want Draft and Final", literals)
	}
}
