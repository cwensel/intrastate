package table_test

// RDR 0002 sections K–M: outcome binding and expansion, derived row
// fields, and the kernel handoff's match/guard routing.

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/resolve"
	"github.com/newcoinc/intrastate/internal/table"
)

// REQ-68: "Every rule — ordinary or escape — MUST bind exactly one
// outcome. **Outcome binding reads the match blocks only** — the rule's
// local `match` block plus the `match` blocks of its inherited contexts.
// That set MUST contain exactly one atom on `recognized`, using `eq` or
// `in`, whose literal(s) are members of the `outcomes` alphabet."
// HAPPY PATH
func TestReq68_EveryRuleBindsExactlyOneOutcomeFromItsMatchBlocks(t *testing.T) {
	for _, rel := range []string{rdrFixture, kataFixture} {
		t.Run(rel, func(t *testing.T) {
			m := mustLoad(t, rel)
			for _, r := range m.Rows {
				if r.Outcome == "" {
					t.Errorf("%s binds no outcome", r.Identity())
					continue
				}
				if !slices.Contains(m.Outcomes, r.Outcome) {
					t.Errorf("%s binds %q, outside the alphabet %v",
						r.Identity(), r.Outcome, m.Outcomes)
				}
			}
		})
	}

	t.Run("binding via an inherited context match block", func(t *testing.T) {
		base := string(readFixture(t, rdrFixture))
		// Move the outcome atom off the rule and onto an inherited context.
		src := strings.Replace(base,
			"[context.archived.match.stage]\neq = \"archive\"",
			"[context.archived.match.stage]\neq = \"archive\"\n[context.archived.match.recognized]\neq = \"finalized\"", 1)
		src = strings.Replace(src,
			"source = \"rdr:terminal\"\n[rule.match.recognized]\nin = [\"finalized\", \"verdict-flapping\"]",
			"source = \"rdr:terminal\"", 1)
		m, err := table.Load([]byte(src), "inherited-outcome.toml")
		if err != nil {
			t.Fatalf("an outcome bound by an inherited context refused: %v", err)
		}
		rows := rowsByRuleID(m, "terminal-archive")
		if len(rows) != 1 || rows[0].Outcome != "finalized" {
			t.Errorf("inherited outcome binding yielded %+v; want one row on "+
				"`finalized`", rows)
		}
	})
}

// REQ-69: "Normalization lifts that atom out of the predicate set into the
// row's outcome field"
// BOUNDARY
//
// "out of" is the discriminating half: a normalizer that copies the atom
// into the outcome field and leaves it in the set passes a
// row.Outcome-only assertion.
func TestReq69_TheRecognizedAtomIsLiftedOutOfThePredicateSet(t *testing.T) {
	for _, rel := range []string{rdrFixture, kataFixture} {
		t.Run(rel, func(t *testing.T) {
			m := mustLoad(t, rel)
			for _, r := range m.Rows {
				if got := atomsOn(r, table.RecognizedTagKey); len(got) != 0 {
					t.Errorf("%s still carries %d `recognized` atoms in its "+
						"predicate set: %+v — the atom is lifted OUT",
						r.Identity(), len(got), got)
				}
			}
		})
	}
}

// REQ-70: "A rule binding zero outcomes, more than one `recognized` atom,
// or a literal outside the alphabet is a load failure."
// ADVERSARIAL
func TestReq70_ZeroTwoOrOutOfAlphabetOutcomeBindingsRefuse(t *testing.T) {
	base := string(readFixture(t, rdrFixture))

	t.Run("literal outside the alphabet", func(t *testing.T) {
		got := loadCategory(t, "neg/neg-outcome-outside.toml")
		if got != table.CatMalformedOutcomeBinding {
			t.Errorf("category = %q; want %q", got, table.CatMalformedOutcomeBinding)
		}
	})

	t.Run("zero recognized atoms", func(t *testing.T) {
		src := strings.Replace(base,
			"[rule.match.recognized]\neq = \"reconcile-block\"",
			"[rule.match.status]\neq = \"Draft\"", 1)
		if src == base {
			t.Fatal("outcome atom removal did not apply")
		}
		_, err := table.Load([]byte(src), "zero-outcomes.toml")
		if err == nil {
			t.Fatal("a rule binding zero outcomes loaded clean")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatMalformedOutcomeBinding {
			t.Errorf("category = %q; want %q", cat, table.CatMalformedOutcomeBinding)
		}
	})

	t.Run("two recognized atoms", func(t *testing.T) {
		// One on the rule and one on an inherited context: both are match
		// blocks, so both are read for binding.
		src := strings.Replace(base,
			"[context.draft.match.status]\neq = \"Draft\"",
			"[context.draft.match.status]\neq = \"Draft\"\n[context.draft.match.recognized]\neq = \"finalized\"", 1)
		if src == base {
			t.Fatal("second outcome atom insertion did not apply")
		}
		_, err := table.Load([]byte(src), "two-outcomes.toml")
		if err == nil {
			t.Fatal("a rule binding two `recognized` atoms loaded clean")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatMalformedOutcomeBinding {
			t.Errorf("category = %q; want %q", cat, table.CatMalformedOutcomeBinding)
		}
	})
}

// REQ-71: "A `recognized` atom authored under `guard.all` or
// `guard.unless` MUST be refused at load as a malformed outcome binding —
// never lifted."
// ADVERSARIAL
func TestReq71_RecognizedAtomInAGuardBlockRefuses(t *testing.T) {
	t.Run("promoted neg-recognized-in-guard", func(t *testing.T) {
		got := loadCategory(t, "neg/neg-recognized-in-guard.toml")
		if got != table.CatMalformedOutcomeBinding {
			t.Errorf("category = %q; want %q", got, table.CatMalformedOutcomeBinding)
		}
	})

	base := string(readFixture(t, rdrFixture))
	for name, block := range map[string]string{
		"guard.all":    "[rule.guard.all.recognized]\neq = \"finalized\"",
		"guard.unless": "[rule.guard.unless.recognized]\neq = \"finalized\"",
	} {
		t.Run(name, func(t *testing.T) {
			src := strings.Replace(base,
				"[rule.guard.unless.profile]\neq = \"small\"",
				"[rule.guard.unless.profile]\neq = \"small\"\n"+block, 1)
			if src == base {
				t.Fatal("guard block insertion did not apply")
			}
			_, err := table.Load([]byte(src), "recognized-in-"+name+".toml")
			if err == nil {
				t.Fatalf("a `recognized` atom under %s loaded clean", name)
			}
			if cat, _ := table.CategoryOf(err); cat != table.CatMalformedOutcomeBinding {
				t.Errorf("category = %q; want %q", cat, table.CatMalformedOutcomeBinding)
			}
		})
	}
}

// REQ-72: "Every `in` atom in a rule's match blocks — local `match` and
// inherited context `match`, on `recognized` or on any other declared tag
// — expands into one candidate row per member, and the rule's rows are the
// Cartesian product of its expanding atoms."
// HAPPY PATH
//
// The expansion is GENERAL, not outcome-only: `continue-prelock` expands
// on the INHERITED non-`recognized` `profile.in`. The product arm is
// witnessed by a rule carrying two multi-member match `in` atoms.
func TestReq72_EveryMatchBlockInExpandsAndTheyProduct(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	t.Run("recognized in expands", func(t *testing.T) {
		if got := len(rowsByRuleID(m, "terminal-archive")); got != 2 {
			t.Errorf("terminal-archive yielded %d rows; want 2", got)
		}
	})
	t.Run("inherited non-recognized in expands", func(t *testing.T) {
		if got := len(rowsByRuleID(m, "continue-prelock")); got != 2 {
			t.Errorf("continue-prelock yielded %d rows; want 2 — the inherited "+
				"`profile.in` expands like any other match-block `in`", got)
		}
	})

	t.Run("two expanding atoms yield their Cartesian product", func(t *testing.T) {
		base := string(readFixture(t, rdrFixture))
		// terminal-archive already has recognized.in over two members;
		// give it a second two-member match `in`.
		src := strings.Replace(base,
			"source = \"rdr:terminal\"\n[rule.match.recognized]\nin = [\"finalized\", \"verdict-flapping\"]",
			"source = \"rdr:terminal\"\n[rule.match.profile]\nin = [\"large\", \"mid\"]\n[rule.match.recognized]\nin = [\"finalized\", \"verdict-flapping\"]", 1)
		if src == base {
			t.Fatal("product substitution did not apply")
		}
		p, err := table.Load([]byte(src), "product.toml")
		if err != nil {
			t.Fatalf("product document refused: %v", err)
		}
		rows := rowsByRuleID(p, "terminal-archive")
		if len(rows) != 4 {
			t.Fatalf("%d rows; want 4 — the Cartesian product of two two-member "+
				"`in` atoms (%v)", len(rows), rowIdentities(p))
		}
		for _, r := range rows {
			if len(r.Suffix) != 2 {
				t.Errorf("%s suffix = %v; want a two-element sequence",
					r.Identity(), r.Suffix)
			}
		}
	})
}

// REQ-73: "In each expanded row the `in` atom becomes an `eq` atom on the
// chosen member, still carrying block `match`"
// BOUNDARY
func TestReq73_ExpandedInBecomesEqStillCarryingBlockMatch(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	for _, id := range []string{"rdr.continue-prelock#large", "rdr.continue-prelock#foundational"} {
		row := rowByID(t, m, id)
		var got *table.Atom
		for _, a := range atomsOn(row, "profile") {
			if a.Block == table.BlockMatch {
				got = &a
			}
		}
		if got == nil {
			t.Errorf("%s carries no match-block `profile` atom", id)
			continue
		}
		if got.Operator != "eq" {
			t.Errorf("%s: expanded operator = %q; want eq", id, got.Operator)
		}
		if got.Block != table.BlockMatch {
			t.Errorf("%s: expanded block = %q; want match", id, got.Block)
		}
		if len(got.Literal) != 1 {
			t.Errorf("%s: expanded literal = %v; want the one chosen member",
				id, got.Literal)
		}
	}

	// No `in` operator survives expansion in any match block.
	for _, r := range m.Rows {
		for _, a := range atomsInBlock(r, table.BlockMatch) {
			if a.Operator == "in" {
				t.Errorf("%s retains a match-block `in` atom %+v; every one "+
					"expands", r.Identity(), a)
			}
		}
	}
}

// REQ-74: "The **expansion suffix** is the sequence of chosen members, one
// per atom with **more than one member**, taken in the atoms' sort order
// `(key, block, operator token, literal)`; it is empty when no such atom
// exists."
// BOUNDARY
func TestReq74_ExpansionSuffixIsASequenceInAtomSortOrder(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	t.Run("empty when no atom expands", func(t *testing.T) {
		row := rowByID(t, m, "rdr.reconcile-rewind")
		if len(row.Suffix) != 0 {
			t.Errorf("suffix = %v; want empty — reconcile-rewind expands nothing",
				row.Suffix)
		}
	})

	t.Run("one element per expanding atom", func(t *testing.T) {
		for _, r := range rowsByRuleID(m, "continue-prelock") {
			if len(r.Suffix) != 1 {
				t.Errorf("%s suffix = %v; want one element", r.Identity(), r.Suffix)
			}
		}
	})

	t.Run("sort order (key, block, operator, literal)", func(t *testing.T) {
		base := string(readFixture(t, rdrFixture))
		// Two expanding atoms: `profile` (key p) and `recognized` (key r).
		// `profile` sorts first, so it is the FIRST suffix element even
		// though `recognized` is authored first.
		src := strings.Replace(base,
			"source = \"rdr:terminal\"\n[rule.match.recognized]\nin = [\"finalized\", \"verdict-flapping\"]",
			"source = \"rdr:terminal\"\n[rule.match.recognized]\nin = [\"finalized\", \"verdict-flapping\"]\n[rule.match.profile]\nin = [\"large\", \"mid\"]", 1)
		if src == base {
			t.Fatal("second expanding atom insertion did not apply")
		}
		p, err := table.Load([]byte(src), "suffix-order.toml")
		if err != nil {
			t.Fatalf("refused: %v", err)
		}
		for _, r := range rowsByRuleID(p, "terminal-archive") {
			if len(r.Suffix) != 2 {
				t.Fatalf("%s suffix = %v; want two elements", r.Identity(), r.Suffix)
			}
			// `profile`'s member first (key "profile" < "recognized"),
			// then `recognized`'s.
			if !slices.Contains([]string{"large", "mid"}, r.Suffix[0]) {
				t.Errorf("%s suffix[0] = %q; want a `profile` member — atoms sort "+
					"by key and `profile` precedes `recognized`", r.Identity(), r.Suffix[0])
			}
			if !slices.Contains([]string{"finalized", "verdict-flapping"}, r.Suffix[1]) {
				t.Errorf("%s suffix[1] = %q; want a `recognized` member",
					r.Identity(), r.Suffix[1])
			}
		}
	})
}

// REQ-75: "A product row whose expanded atoms cannot be satisfied together
// (two `eq` on one key, different literals) is a dead row for RDR 0006,
// not a load failure."
// INPUT EDGE
func TestReq75_UnsatisfiableProductRowIsNotALoadFailure(t *testing.T) {
	base := string(readFixture(t, rdrFixture))
	// A second match `in` on `profile`, disjoint from the inherited one:
	// some product rows carry two conflicting `profile.eq` atoms.
	src := strings.Replace(base,
		"source = \"rdr:prelock\"\ngate = [\"rdr-lock\"]",
		"source = \"rdr:prelock\"\ngate = [\"rdr-lock\"]\n[rule.match.profile]\nin = [\"small\", \"mid\"]", 1)
	if src == base {
		t.Fatal("conflicting atom insertion did not apply")
	}

	m, err := table.Load([]byte(src), "dead-product.toml")
	if err != nil {
		t.Fatalf("a product row that no view satisfies refused: %v — it is RDR "+
			"0006's dead-row finding, not a load failure", err)
	}
	if got := len(rowsByRuleID(m, "continue-prelock")); got != 4 {
		t.Errorf("continue-prelock yielded %d rows; want 4 — the product is "+
			"formed even where a row is dead", got)
	}
}

// REQ-76: "A single-member `in` expands to one row whose suffix is empty,
// so `eq = "x"` and `in = ["x"]` are one spelling of one edge and cannot
// mint two identities. A suffix is non-empty exactly when the rule
// produced more than one row."
// BOUNDARY
func TestReq76_SingleMemberInIsOneSpellingOfOneEdge(t *testing.T) {
	t.Run("eq and single-member in normalize identically", func(t *testing.T) {
		base := mustLoad(t, rdrFixture)
		perm := mustLoad(t, "perm/rdr-eq-as-in.toml")
		if !reflect.DeepEqual(base.Rows, perm.Rows) {
			t.Error("`eq = \"x\"` and `in = [\"x\"]` normalized to different row " +
				"sets; they are one spelling of one edge")
		}
	})

	t.Run("suffix non-empty exactly when the rule produced more than one row", func(t *testing.T) {
		for _, rel := range []string{rdrFixture, kataFixture, "perm/rdr-eq-as-in.toml"} {
			m := mustLoad(t, rel)
			counts := map[string]int{}
			for _, r := range m.Rows {
				counts[r.RuleID]++
			}
			for _, r := range m.Rows {
				multi := counts[r.RuleID] > 1
				if multi != (len(r.Suffix) > 0) {
					t.Errorf("%s: %s produced %d rows but carries suffix %v",
						rel, r.Identity(), counts[r.RuleID], r.Suffix)
				}
			}
		}
	})
}

// REQ-77: "`Row.RequiresOwned` has no authored form. The normalizer MUST
// derive it for every row as the sorted, duplicate-free set of tag keys
// named by the rule's write block and clear list"
// HAPPY PATH
func TestReq77_RequiresOwnedIsDerivedFromWritesAndClears(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	cases := map[string][]string{
		// Writes stage and iter.
		"rdr.continue-prelock#large": {"iter", "stage"},
		// Writes stage and prelock_lens.
		"rdr.continue-prelock-cluster#large": {"prelock_lens", "stage"},
		// Writes stage, status, rewind_scope; clears prelock_lens.
		"rdr.reconcile-rewind": {"prelock_lens", "rewind_scope", "stage", "status"},
	}
	for id, want := range cases {
		row := rowByID(t, m, id)
		if !reflect.DeepEqual(row.RequiresOwned, want) {
			t.Errorf("%s RequiresOwned = %v; want the sorted, duplicate-free "+
				"write-plus-clear key set %v", id, row.RequiresOwned, want)
		}
	}

	t.Run("no authored form", func(t *testing.T) {
		base := string(readFixture(t, rdrFixture))
		src := strings.Replace(base, `gate = ["rdr-lock"]`,
			"gate = [\"rdr-lock\"]\nrequires_owned = [\"status\"]", 1)
		if src == base {
			t.Fatal("requires_owned insertion did not apply")
		}
		_, err := table.Load([]byte(src), "authored-requires-owned.toml")
		if err == nil {
			t.Fatal("an authored `requires_owned` loaded clean; the field has no " +
				"authored form")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatUnknownSchemaField {
			t.Errorf("category = %q; want %q", cat, table.CatUnknownSchemaField)
		}
	})
}

// REQ-78: "An escape rule carries neither a write block nor a clear list,
// so a normalized escape row MUST carry an empty set (RDR 0007 A21)."
// BOUNDARY
func TestReq78_EscapeRowsCarryAnEmptyRequiresOwnedSet(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	for _, r := range m.Rows {
		if len(r.Escape) == 0 {
			continue
		}
		if len(r.RequiresOwned) != 0 {
			t.Errorf("escape row %s RequiresOwned = %v; want empty",
				r.Identity(), r.RequiresOwned)
		}
	}
}

// REQ-79: "this RDR is its producer (JDR 0001 §JD-3) and does not add
// guard-read keys to it."
// ADVERSARIAL
//
// `continue-prelock` guards on `iter` (also written, so legitimately in
// the set) and on `finalized_at` (observed, guard-read only). If the
// normalizer added guard-read keys, `finalized_at` would appear.
// `continue-prelock-cluster` guards on the owned `cluster_ready`, which it
// does NOT write: that is the discriminating case.
func TestReq79_GuardReadKeysAreNotAddedToRequiresOwned(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	prelock := rowByID(t, m, "rdr.continue-prelock#large")
	if slices.Contains(prelock.RequiresOwned, "finalized_at") {
		t.Errorf("RequiresOwned = %v; a guard-read observed key was added",
			prelock.RequiresOwned)
	}

	cluster := rowByID(t, m, "rdr.continue-prelock-cluster#large")
	if slices.Contains(cluster.RequiresOwned, "cluster_ready") {
		t.Errorf("RequiresOwned = %v; `cluster_ready` is guard-read and not "+
			"written, so this RDR does not add it", cluster.RequiresOwned)
	}
}

// REQ-80: "Normalization MUST populate the next-state tags with the tag
// values the rule's write block and clear list produce, and the writes
// with the owned-tag writes the accessor layer applies — the same rendered
// set, including `<clear>` entries."
// HAPPY PATH
func TestReq80_NextTagsAndWritesAreBothPopulatedIncludingClears(t *testing.T) {
	m := mustLoad(t, rdrFixture)
	row := rowByID(t, m, "rdr.reconcile-rewind")

	want := []table.TagValue{
		{Key: "prelock_lens", Value: []string{table.ClearSentinel}},
		{Key: "rewind_scope", Value: []string{"assumptions"}},
		{Key: "stage", Value: []string{"resolve"}},
		{Key: "status", Value: []string{"Draft"}},
	}
	if !reflect.DeepEqual(row.NextTags, want) {
		t.Errorf("NextTags =\n %+v\nwant\n %+v", row.NextTags, want)
	}
	if !reflect.DeepEqual(row.Writes, want) {
		t.Errorf("Writes =\n %+v\nwant\n %+v", row.Writes, want)
	}
}

// REQ-81: "An escape row carries neither a write block nor a clear list,
// so it normalizes to a row with both empty"
// BOUNDARY
func TestReq81_EscapeRowsCarryEmptyNextTagsAndWrites(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	for _, r := range m.Rows {
		if len(r.Escape) == 0 {
			continue
		}
		if len(r.NextTags) != 0 {
			t.Errorf("escape row %s NextTags = %+v; want empty", r.Identity(), r.NextTags)
		}
		if len(r.Writes) != 0 {
			t.Errorf("escape row %s Writes = %+v; want empty", r.Identity(), r.Writes)
		}
	}
}

// REQ-82: "Normalization MUST populate both explicitly from the rule and
// MUST NOT populate one by aliasing the other"
// ADVERSARIAL (TS scenario 2's mutation test)
//
// Value comparison alone is insufficient: the two fields hold equal sets
// under this RDR's authoring surface, so an aliased pair and an
// independently built pair compare equal. Both are slices, so an alias
// shares a backing array — the control MUTATES one field in place and
// asserts the other is unchanged.
func TestReq82_NextTagsAndWritesAreNotAliased(t *testing.T) {
	m := mustLoad(t, rdrFixture)
	row := rowByID(t, m, "rdr.reconcile-rewind")

	if len(row.NextTags) == 0 || len(row.Writes) == 0 {
		t.Fatal("the row under test carries no writes")
	}

	before := make([]table.TagValue, len(row.Writes))
	copy(before, row.Writes)

	// Mutate one field in place. An aliased pair shares the backing array,
	// so the sibling changes too.
	row.NextTags[0].Key = "MUTATED"
	row.NextTags[0].Value = []string{"MUTATED"}

	if !reflect.DeepEqual(row.Writes, before) {
		t.Errorf("mutating NextTags[0] changed Writes:\n got = %+v\nwant = %+v — "+
			"the two fields must be populated independently, never by aliasing",
			row.Writes, before)
	}

	// And the mirror direction.
	m2 := mustLoad(t, rdrFixture)
	row2 := rowByID(t, m2, "rdr.reconcile-rewind")
	beforeNext := make([]table.TagValue, len(row2.NextTags))
	copy(beforeNext, row2.NextTags)

	row2.Writes[0].Key = "MUTATED"
	row2.Writes[0].Value = []string{"MUTATED"}

	if !reflect.DeepEqual(row2.NextTags, beforeNext) {
		t.Errorf("mutating Writes[0] changed NextTags:\n got = %+v\nwant = %+v",
			row2.NextTags, beforeNext)
	}
}

// REQ-83: "Clearing a tag MUST be represented by an explicit rule-level
// `clear` entry that normalization renders as a `<clear>` write. Absence
// from both the write block and the clear list MUST NOT imply deletion."
// BOUNDARY
func TestReq83_ClearingIsExplicitAndAbsenceIsNotDeletion(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	t.Run("an explicit clear renders as a <clear> write", func(t *testing.T) {
		row := rowByID(t, m, "rdr.reconcile-rewind")
		got, ok := tagValue(row.Writes, "prelock_lens")
		if !ok {
			t.Fatal("the cleared key is absent from the writes")
		}
		if !reflect.DeepEqual(got, []string{table.ClearSentinel}) {
			t.Errorf("cleared value = %v; want [%q]", got, table.ClearSentinel)
		}
	})

	t.Run("absence implies no deletion", func(t *testing.T) {
		// `continue-prelock` writes stage and iter and clears nothing. No
		// other owned tag may appear as a <clear> write.
		row := rowByID(t, m, "rdr.continue-prelock#large")
		for _, w := range row.Writes {
			if len(w.Value) == 1 && w.Value[0] == table.ClearSentinel {
				t.Errorf("%q is rendered as a <clear> write on a rule with no "+
					"clear list; absence must not imply deletion", w.Key)
			}
		}
		if len(row.Writes) != 2 {
			t.Errorf("Writes = %+v; want exactly the two authored assignments",
				row.Writes)
		}
	})
}

// REQ-84: "The normalized candidate row carries the atoms as **one unified
// set with each atom's authored block retained**; that single field is
// what the dump contract lists, what the Round-Trip invariant compares,
// and what RDR 0003 reads downstream. The split is applied when a
// `resolve.Row` is constructed for the kernel."
// BOUNDARY
//
// The structural half: the normalized row carries ONE atom field, not a
// pre-split Match/Guard pair. A row already split could not satisfy the
// block-retention obligation.
func TestReq84_NormalizedRowCarriesOneUnifiedAtomSet(t *testing.T) {
	rt := reflect.TypeOf(table.Row{})
	for i := range rt.NumField() {
		switch rt.Field(i).Name {
		case "Match", "Guard":
			t.Errorf("table.Row carries a %q field; the normalized value holds "+
				"one unified atom set and the split happens at the kernel handoff",
				rt.Field(i).Name)
		}
	}

	f, ok := rt.FieldByName("Atoms")
	if !ok {
		t.Fatal("table.Row has no Atoms field")
	}
	if f.Type.Elem() != reflect.TypeOf(table.Atom{}) {
		t.Errorf("Atoms is %v; want []table.Atom", f.Type)
	}
}

// REQ-85: "the handoff MUST route each atom to exactly one of them by its
// block: every atom carrying block `match` populates `Match`; every atom
// carrying block `all` or `unless` belongs to the guard predicate,
// **regardless of operator** — an `eq` atom under `guard.all` is a guard
// atom. The routing is exhaustive and disjoint by construction"
// ADVERSARIAL (A12)
//
// The discriminating cases are equality operators that route to the guard
// BECAUSE of their block: `continue-prelock-cluster`'s
// `[rule.guard.all.cluster_ready] eq = true` and the kata fixture's
// `status.eq=closed@unless`. Operator-keyed routing sends both to Match.
func TestReq85_HandoffRoutesByBlockNeverByOperator(t *testing.T) {
	t.Run("eq under guard.all is a guard atom", func(t *testing.T) {
		m := mustLoad(t, rdrFixture)
		kr := rowByID(t, m, "rdr.continue-prelock-cluster#large").KernelRow()

		if !hasGuardAtom(kr.Guard, "cluster_ready", "eq") {
			t.Errorf("`cluster_ready.eq` is absent from the kernel guard %+v; an "+
				"`eq` atom under guard.all is a guard atom", kr.Guard)
		}
		for _, mt := range kr.Match {
			if mt.Key == "cluster_ready" {
				t.Errorf("`cluster_ready` reached Match; routing is by block, "+
					"never by operator (Match = %+v)", kr.Match)
			}
		}
	})

	t.Run("eq under guard.unless is a guard atom", func(t *testing.T) {
		m := mustLoad(t, kataFixture)
		kr := rowByID(t, m, "kata.review-accepted").KernelRow()

		var unless *resolve.GuardAtom
		for i, a := range kr.Guard {
			if a.Key == "status" && a.Block == resolve.BlockUnless {
				unless = &kr.Guard[i]
			}
		}
		if unless == nil {
			t.Fatalf("`status.eq=closed@unless` is absent from the guard %+v", kr.Guard)
		}
		if unless.Operator != "eq" {
			t.Errorf("operator = %q; want eq — the operator is unchanged by routing",
				unless.Operator)
		}
	})

	t.Run("routing is exhaustive and disjoint over both fixtures", func(t *testing.T) {
		for _, rel := range []string{rdrFixture, kataFixture} {
			m := mustLoad(t, rel)
			for _, row := range m.Rows {
				kr := row.KernelRow()

				if got := len(kr.Match) + len(kr.Guard); got != len(row.Atoms) {
					t.Errorf("%s: |Match| + |Guard| = %d; want %d — the routing is "+
						"exhaustive and disjoint (atoms = %+v)",
						row.Identity(), got, len(row.Atoms), row.Atoms)
				}
				matchKeys := map[string]bool{}
				for _, mt := range kr.Match {
					matchKeys[mt.Key] = true
				}
				for _, g := range kr.Guard {
					if g.Block == resolve.BlockMatch {
						t.Errorf("%s: a BlockMatch atom %+v reached the guard",
							row.Identity(), g)
					}
				}
				// Every Match entry traces to a match-block atom.
				for _, mt := range kr.Match {
					if !containsAtom(atomsInBlock(row, table.BlockMatch), table.Atom{
						Key: mt.Key, Block: table.BlockMatch, Operator: "eq",
						Literal: []string{mt.Value},
					}) {
						t.Errorf("%s: Match entry %+v has no match-block atom",
							row.Identity(), mt)
					}
				}
			}
		}
	})
}

func hasGuardAtom(atoms []resolve.GuardAtom, key, op string) bool {
	for _, a := range atoms {
		if a.Key == key && a.Operator == op {
			return true
		}
	}
	return false
}

// REQ-86: "a **match block admits only `eq` and `in`** (`in` by expansion
// into per-member `eq` rows): a comparison, existence, `contains`, or any
// other operator under `[rule.match]` or `[context.<id>.match]` is refused
// at load as a malformed predicate atom"
// ADVERSARIAL
func TestReq86_MatchBlocksAdmitOnlyEqAndIn(t *testing.T) {
	t.Run("promoted controls", func(t *testing.T) {
		for _, rel := range []string{"neg/neg-match-lt.toml", "neg/neg-match-exists.toml"} {
			if got := loadCategory(t, rel); got != table.CatMalformedPredicateAtom {
				t.Errorf("%s: category = %q; want %q", rel, got, table.CatMalformedPredicateAtom)
			}
		}
	})

	base := string(readFixture(t, rdrFixture))
	for _, op := range []string{"lt", "lte", "gt", "gte", "exists", "contains"} {
		t.Run("match-block "+op, func(t *testing.T) {
			lit := `"x"`
			switch op {
			case "lt", "lte", "gt", "gte":
				lit = "3"
			case "exists":
				lit = "true"
			case "contains":
				lit = `["x"]`
			}
			src := strings.Replace(base,
				"[context.draft.match.status]\neq = \"Draft\"",
				"[context.draft.match.status]\n"+op+" = "+lit, 1)
			if src == base {
				t.Fatal("match operator substitution did not apply")
			}
			_, err := table.Load([]byte(src), "match-"+op+".toml")
			if err == nil {
				t.Fatalf("operator %q under a match block loaded clean", op)
			}
			if cat, _ := table.CategoryOf(err); cat != table.CatMalformedPredicateAtom {
				t.Errorf("category = %q; want %q", cat, table.CatMalformedPredicateAtom)
			}
		})
	}

	t.Run("eq and in are admitted", func(t *testing.T) {
		if _, err := table.Load(readFixture(t, rdrFixture), rdrFixture); err != nil {
			t.Errorf("the fixture's match-block eq and in refused: %v", err)
		}
	})
}

// REQ-87: "**The admitted operator set, guard blocks included, is RDR
// 0003's closed set `{eq, in, lt, lte, gt, gte, exists, contains}`.** This
// RDR does not mint operators and does not widen that set … A token under
// `[rule.guard.all.<tag>]` or `[rule.guard.unless.<tag>]` outside the set
// is a `malformed predicate atom` (`unknown operator`) refused at load"
// ADVERSARIAL
func TestReq87_AdmittedOperatorSetIsClosedAtEight(t *testing.T) {
	t.Run("the exported set is exactly RDR 0003's eight tokens", func(t *testing.T) {
		want := []string{"contains", "eq", "exists", "gt", "gte", "in", "lt", "lte"}
		got := slices.Clone(table.Operators())
		slices.Sort(got)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("operator set = %v; want %v — this RDR does not mint "+
				"operators and does not widen the set", got, want)
		}
	})

	t.Run("promoted guard-block unknown-operator controls", func(t *testing.T) {
		for _, rel := range []string{
			"neg/neg-guard-all-bad-operator.toml",
			"neg/neg-guard-unless-bad-operator.toml",
		} {
			if got := loadCategory(t, rel); got != table.CatMalformedPredicateAtom {
				t.Errorf("%s: category = %q; want %q", rel, got, table.CatMalformedPredicateAtom)
			}
		}
	})
}

// REQ-88: "**Atom-level validation is block-agnostic.** Every rule this
// RDR states about an authored atom — operator membership (above), the
// `<clear>` reserved value, domain and kind conformance of a literal, `#`
// reservation in members, and tag-key declaration — MUST be enforced
// identically in all three atom blocks: `match`, `guard.all`, and
// `guard.unless`."
// ADVERSARIAL (A15)
//
// `gen-cases.py` originally mutated match blocks only, so a rule "witnessed"
// in one block proved nothing about the other two. Each block-agnostic rule
// gets a control per block.
func TestReq88_AtomLevelValidationIsBlockAgnostic(t *testing.T) {
	base := string(readFixture(t, rdrFixture))

	// Anchors: one authored atom per block that a mutation can replace.
	anchors := map[string]string{
		"match":        "[context.draft.match.status]\neq = \"Draft\"",
		"guard.all":    "[rule.guard.all.iter]\nlt = 3",
		"guard.unless": "[rule.guard.unless.profile]\neq = \"small\"",
	}
	// Each rule, spelled per block against that block's tag.
	rules := map[string]struct {
		mutate map[string]string
		want   table.Category
	}{
		"operator membership": {
			mutate: map[string]string{
				"match":        "[context.draft.match.status]\nfrobnicate = \"Draft\"",
				"guard.all":    "[rule.guard.all.iter]\nfrobnicate = 3",
				"guard.unless": "[rule.guard.unless.profile]\nfrobnicate = \"small\"",
			},
			want: table.CatMalformedPredicateAtom,
		},
		"<clear> reserved value": {
			mutate: map[string]string{
				"match":        "[context.draft.match.status]\neq = \"<clear>\"",
				"guard.all":    "[rule.guard.all.iter]\neq = \"<clear>\"",
				"guard.unless": "[rule.guard.unless.profile]\neq = \"<clear>\"",
			},
			want: table.CatReservedTagValue,
		},
		"domain conformance": {
			mutate: map[string]string{
				"match":        "[context.draft.match.status]\neq = \"Archived\"",
				"guard.all":    "[rule.guard.all.iter]\neq = \"NOT_AN_INT\"",
				"guard.unless": "[rule.guard.unless.profile]\neq = \"NOT_A_PROFILE\"",
			},
			want: table.CatMalformedPredicateAtom,
		},
		"tag-key declaration": {
			mutate: map[string]string{
				"match":        "[context.draft.match.frobnitz]\neq = \"Draft\"",
				"guard.all":    "[rule.guard.all.frobnitz]\nlt = 3",
				"guard.unless": "[rule.guard.unless.frobnitz]\neq = \"small\"",
			},
			want: table.CatUnknownTag,
		},
	}

	for rule, spec := range rules {
		for block, anchor := range anchors {
			t.Run(rule+" in "+block, func(t *testing.T) {
				src := strings.Replace(base, anchor, spec.mutate[block], 1)
				if src == base {
					t.Fatalf("anchor %q did not apply", anchor)
				}
				_, err := table.Load([]byte(src), "block-agnostic.toml")
				if err == nil {
					t.Fatalf("%s in %s loaded clean; atom-level validation is "+
						"block-agnostic", rule, block)
				}
				if cat, _ := table.CategoryOf(err); cat != spec.want {
					t.Errorf("category = %q; want %q", cat, spec.want)
				}
			})
		}
	}

	t.Run("promoted guard-block controls", func(t *testing.T) {
		promoted := map[string]table.Category{
			"neg/neg-guard-all-clear.toml":            table.CatReservedTagValue,
			"neg/neg-guard-unless-clear.toml":         table.CatReservedTagValue,
			"neg/neg-guard-unless-out-of-domain.toml": table.CatMalformedPredicateAtom,
			"neg/neg-guard-unless-unknown-tag.toml":   table.CatUnknownTag,
		}
		for rel, want := range promoted {
			if got := loadCategory(t, rel); got != want {
				t.Errorf("%s: category = %q; want %q", rel, got, want)
			}
		}
	})
}

// REQ-89: "**Conformance is per operator, not per literal.** … `eq`, `in`,
// and `contains` take members and are domain-checked; `exists` takes a
// bool literal (`<true>`/`<false>`), never a member; `lt`/`lte`/`gt`/`gte`
// take an ordered **bound**, which is kind-checked but not domain-checked"
// BOUNDARY
//
// The RDR fixture's own `iter.lt = 3` against `min = 0, max = 9` is the
// witness for the bound arm: a bound need not itself be an authorable
// value, so `lt = 99` must LOAD while `eq = 99` on the same tag refuses.
func TestReq89_ConformanceIsPerOperator(t *testing.T) {
	base := string(readFixture(t, rdrFixture))

	t.Run("a comparison bound is kind-checked but not domain-checked", func(t *testing.T) {
		src := strings.Replace(base, "[rule.guard.all.iter]\nlt = 3",
			"[rule.guard.all.iter]\nlt = 99", 1)
		if src == base {
			t.Fatal("bound substitution did not apply")
		}
		if _, err := table.Load([]byte(src), "bound-out-of-range.toml"); err != nil {
			t.Errorf("`lt = 99` against min=0 max=9 refused with %v; a bound is "+
				"kind-checked but not domain-checked", err)
		}
	})

	t.Run("a comparison bound of the wrong kind refuses", func(t *testing.T) {
		src := strings.Replace(base, "[rule.guard.all.iter]\nlt = 3",
			"[rule.guard.all.iter]\nlt = \"three\"", 1)
		if src == base {
			t.Fatal("bound substitution did not apply")
		}
		_, err := table.Load([]byte(src), "bound-wrong-kind.toml")
		if err == nil {
			t.Fatal("`lt = \"three\"` on an int tag loaded clean")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatMalformedPredicateAtom {
			t.Errorf("category = %q; want %q", cat, table.CatMalformedPredicateAtom)
		}
	})

	t.Run("eq is domain-checked", func(t *testing.T) {
		src := strings.Replace(base, "[rule.guard.all.iter]\nlt = 3",
			"[rule.guard.all.iter]\neq = 99", 1)
		if src == base {
			t.Fatal("eq substitution did not apply")
		}
		_, err := table.Load([]byte(src), "eq-out-of-range.toml")
		if err == nil {
			t.Fatal("`eq = 99` against min=0 max=9 loaded clean; eq takes a " +
				"member and is domain-checked")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatMalformedPredicateAtom {
			t.Errorf("category = %q; want %q", cat, table.CatMalformedPredicateAtom)
		}
	})

	t.Run("exists takes a bool literal, never a member", func(t *testing.T) {
		if got := loadCategory(t, "neg/neg-bad-exists.toml"); got != table.CatMalformedPredicateAtom {
			t.Errorf("category = %q; want %q", got, table.CatMalformedPredicateAtom)
		}
	})
}

// REQ-90: "The `<clear>` ban and the tag-key declaration rule bind every
// atom regardless of operator. Two rules are legitimately match-only and
// are **not** block-agnostic: the `eq`/`in` operator restriction (§D6
// routing) and the `#` reservation in `in` members, since only match
// blocks expand."
// BOUNDARY
//
// The discriminating half: the two match-only rules must NOT fire in a
// guard block. A `contains` under guard.all is legal, and a `#` inside a
// guard-block `in` member is legal, because only match blocks expand.
func TestReq90_TwoRulesAreLegitimatelyMatchOnly(t *testing.T) {
	base := string(readFixture(t, rdrFixture))

	t.Run("a non-eq/in operator is legal in a guard block", func(t *testing.T) {
		// `lt` under guard.all is the fixture's own authoring, and
		// `contains` on a set tag is the other member-taking operator.
		if _, err := table.Load(readFixture(t, rdrFixture), rdrFixture); err != nil {
			t.Errorf("the fixture's guard-block `lt` refused: %v", err)
		}
	})

	t.Run("`#` inside a guard-block `in` member is legal", func(t *testing.T) {
		src := strings.Replace(base, "[rule.guard.unless.profile]\neq = \"small\"",
			"[rule.guard.unless.finalized_at]\nin = [\"a#b\"]", 1)
		if src == base {
			t.Fatal("guard `in` substitution did not apply")
		}
		if _, err := table.Load([]byte(src), "guard-in-hash.toml"); err != nil {
			t.Errorf("a `#` inside a guard-block `in` member refused with %v; only "+
				"match blocks expand, so the reservation is match-only", err)
		}
	})

	t.Run("`#` inside a match-block `in` member refuses", func(t *testing.T) {
		if got := loadCategory(t, "neg/neg-in-member-hash.toml"); got != table.CatMalformedPredicateAtom {
			t.Errorf("category = %q; want %q", got, table.CatMalformedPredicateAtom)
		}
	})
}
