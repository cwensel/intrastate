package table_test

// RDR 0002 Phase 3b — adversarial failure-mode tests.
//
// Every test here is anchored in the record's
// `## Trade-offs > ### Failure Modes` section, which enumerates what MUST
// fail at load "with stable CLI errors — each is decidable from one rule
// plus the declarations", and which cross-row conditions ("a row gap,
// overlap, dead row, read-before-write condition, or ambiguous expansion")
// are instead RDR 0006's. The three modes below are the places where the
// implementation quietly lands on the wrong side of that line.
//
// These are ADDED tests. Nothing existing is weakened.

import (
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/table"
)

// ---------------------------------------------------------------- ADV-1

// FAILURE MODE 1 — expansion escapes the match blocks.
//
// Failure Modes: "malformed predicate atoms (including a non-`eq`/`in`
// operator under a match block or a literal outside its domain)" fail at
// load, while "ambiguous expansion requires comparing normalized rows and
// is therefore an RDR 0006 lint failure". Expansion is the record's
// match-block-only mechanism:
//
//	REQ-72 (`0002:C13`): "Every `in` atom in a rule's MATCH BLOCKS — local
//	`match` and inherited context `match` … expands into one candidate row
//	per member, and the rule's rows are the Cartesian product of its
//	expanding atoms."
//	REQ-73 (`0002:C13`): "In each expanded row the `in` atom becomes an
//	`eq` atom on the chosen member, STILL CARRYING BLOCK `match`".
//	REQ-51 (`0002:C7`): "Normalization MUST NOT fold `unless` atoms into
//	`all` or `match` atoms into either guard block."
//	REQ-74 (`0002:C13`): the expansion suffix is "the sequence of chosen
//	members, one per atom with more than one member".
//
// A guard-block `in` is a single predicate over a member set: `unless x in
// {a, b}` means "not a AND not b" — one row. Expanding it mints one row per
// member, each excluding only its own member, which turns the author's
// conjunction into a disjunction: every minted row now MATCHES a view the
// author excluded. The defect is silent — the extra rows are well-formed,
// carry a plausible expansion suffix, and the dump certifies them as
// canonical — and it is not recoverable downstream, because RDR 0006 sees
// rows, not the authored guard.
//
// The same over-reach bleeds `#` into the rendered row identity: the `#`
// ban is deliberately match-only (REQ-90, "only match blocks expand"), so a
// guard `in` member carrying `#` is legally authorable and, once wrongly
// expanded, forges an ambiguous identity — the exact collision `0002:C11`
// exists to close.
func TestAdv1_GuardBlockInMustNotExpand(t *testing.T) {
	base := string(readFixture(t, rdrFixture))

	// `continue-prelock` inherits `profile in ["large", "foundational"]`
	// from `large-prelock` (two match rows) and authors one `unless` atom
	// on `profile`. Replacing that `unless eq` with an `unless in` over two
	// members must leave the rule at TWO rows, not four.
	src := strings.Replace(base,
		"[rule.guard.unless.profile]\neq = \"small\"",
		"[rule.guard.unless.profile]\nin = [\"small\", \"mid\"]", 1)
	if src == base {
		t.Fatal("guard.unless substitution did not apply")
	}

	m, err := table.Load([]byte(src), "adv1-guard-in.toml")
	if err != nil {
		t.Fatalf("a guard-block `in` refused at load: %v — `in` is an "+
			"admitted operator in every block (REQ-87)", err)
	}

	rows := rowsByRuleID(m, "continue-prelock")
	if len(rows) != 2 {
		t.Errorf("rule expanded to %d rows; want 2 — only the inherited "+
			"MATCH-block `in` on `profile` expands (REQ-72). Identities: %v",
			len(rows), func() []string {
				out := make([]string, 0, len(rows))
				for _, r := range rows {
					out = append(out, r.Identity())
				}
				return out
			}())
	}

	for _, r := range rows {
		// REQ-74: one suffix element per EXPANDING atom. Only `profile`
		// in the match block expands, so the suffix is exactly one member.
		if len(r.Suffix) != 1 {
			t.Errorf("row %s carries a %d-element expansion suffix %v; want 1 — "+
				"a guard atom contributes no suffix element (REQ-74)",
				r.Identity(), len(r.Suffix), r.Suffix)
		}

		// REQ-51/REQ-73: the `unless` atom must survive verbatim as the
		// single `in` predicate it was authored as, never folded into a
		// per-member `eq`.
		var unless []table.Atom
		for _, a := range atomsOn(r, "profile") {
			if a.Block == table.BlockUnless {
				unless = append(unless, a)
			}
		}
		if len(unless) != 1 {
			t.Errorf("row %s carries %d `unless` atoms on `profile`; want 1 — "+
				"the authored `in` is one predicate, not one per member "+
				"(REQ-51): %+v", r.Identity(), len(unless), unless)
			continue
		}
		if unless[0].Operator != "in" {
			t.Errorf("row %s: `unless` atom operator = %q; want \"in\" — only a "+
				"MATCH-block `in` becomes an `eq` on a chosen member (REQ-73)",
				r.Identity(), unless[0].Operator)
		}
		if got, want := strings.Join(unless[0].Literal, ","), "mid,small"; got != want {
			t.Errorf("row %s: `unless` literal = %v; want the whole member set "+
				"[mid small] (REQ-56/REQ-57), member-sorted",
				r.Identity(), unless[0].Literal)
		}
	}

	// The semantic consequence, stated as its own oracle: no normalized row
	// may carry an `unless` atom that admits a member the author excluded.
	// With the guard `in` expanded, `#large#mid` excludes only `mid` and so
	// MATCHES profile=small — a view the authored rule refused.
	for _, r := range rows {
		for _, a := range atomsOn(r, "profile") {
			if a.Block != table.BlockUnless {
				continue
			}
			if a.Operator == "eq" && len(a.Literal) == 1 {
				t.Errorf("row %s excludes only %q; the authored guard excluded "+
					"both `small` and `mid`. Expanding a guard `in` turns the "+
					"author's conjunction into a disjunction and admits views "+
					"the rule refused",
					r.Identity(), a.Literal[0])
			}
		}
	}
}

// TestAdv1b_GuardInMustNotBleedHashIntoRowIdentity is ADV-1's identity
// half. `0002:C11` (REQ-60) bans `#` only in rule ids, outcome literals,
// and MATCH-block `in` members, and REQ-90 states the narrowing on purpose:
// "the `#` reservation in `in` members, since ONLY MATCH BLOCKS EXPAND". A
// guard `in` member holding `#` is therefore legal — and must stay out of
// the rendered identity, whose recoverability rests on splitting on `#`.
func TestAdv1b_GuardInMustNotBleedHashIntoRowIdentity(t *testing.T) {
	base := string(readFixture(t, rdrFixture))

	// Widen `prelock_lens`'s domain to admit a `#`-bearing member, then
	// author it as a guard-block `in`.
	src := strings.Replace(base,
		`domain = ["grounding", "cove", "3amigo", "critique", "repeatability"]`,
		`domain = ["grounding", "cove", "3amigo", "critique", "repeatability", "a#b"]`, 1)
	if src == base {
		t.Fatal("prelock_lens domain substitution did not apply")
	}
	src2 := strings.Replace(src,
		"[rule.guard.unless.profile]\neq = \"small\"",
		"[rule.guard.unless.prelock_lens]\nin = [\"a#b\", \"cove\"]", 1)
	if src2 == src {
		t.Fatal("guard.unless substitution did not apply")
	}

	m, err := table.Load([]byte(src2), "adv1b-guard-in-hash.toml")
	if err != nil {
		t.Fatalf("a `#`-bearing GUARD `in` member refused at load: %v — the "+
			"`#` ban is match-only (REQ-60, REQ-90)", err)
	}

	for _, r := range rowsByRuleID(m, "continue-prelock") {
		id := r.Identity()
		// The identity is `<model>.<ruleID>` plus one `#`-joined element
		// per suffix member. With only the inherited match `in` expanding,
		// exactly one `#` may appear.
		if n := strings.Count(id, "#"); n != 1 {
			t.Errorf("row identity %q carries %d `#` separators; want 1 — a "+
				"guard-block member reached the rendered identity, which "+
				"`0002:C11` closes by banning `#` only where expansion can "+
				"carry it", id, n)
		}
		for _, s := range r.Suffix {
			if strings.Contains(s, "#") {
				t.Errorf("row %q carries the suffix element %q containing `#`; "+
					"splitting the rendered identity on `#` is no longer "+
					"unambiguous (REQ-60)", id, s)
			}
		}
	}
}

// ---------------------------------------------------------------- ADV-2

// FAILURE MODE 2 — a rule writing a non-owned tag is refused by the wrong
// authority, under the wrong category.
//
// Failure Modes: "... or A WRITE TO AN UNDECLARED OR NON-OWNED TAG fail at
// load/validation time with stable CLI errors — EACH IS DECIDABLE FROM ONE
// RULE PLUS THE DECLARATIONS."
//
// The record leans on this being a rule-level check, not an accessor-table
// one. `0002:C14` (REQ-77) derives `RequiresOwned` from the write block and
// clear list and then asserts: "BY THE WRITE-TO-NON-OWNED-TAG RULE every
// such key is an owned tag". That guarantee holds only if the rule's own
// write block is checked against the declarations.
//
// The implementation instead checks provenance only on `[write.<id>].keys`,
// and catches a rule writing an observed tag downstream, via the
// writer-arity count, as `malformed_accessor_binding`. Two things go wrong:
//
//   - REQ-110 (`0002:C24`): "The Testing Strategy asserts on the category,
//     so these identifiers are an API surface, not message text." The
//     caller is told the model's ACCESSOR WIRING is wrong when the defect
//     is the RULE. `TestReq27_...` freezes that the arity and provenance
//     obligations "do NOT collapse into one category" — they collapse here.
//   - The refusal is not decidable from one rule plus the declarations: it
//     needs the whole writer table and every rule's write set. Author a
//     writer serving the observed key and the arity check is satisfied,
//     leaving the non-owned write to be caught, if at all, by a different
//     rule than the one that is wrong.
//
// `RequiresOwned` is the kernel's owned-state gate (`resolve.go::
// missingOwned`), so a non-owned key reaching it is a runtime refusal the
// load was supposed to pre-empt.
func TestAdv2_RuleWriteToNonOwnedTagIsRefusedAsSuch(t *testing.T) {
	base := string(readFixture(t, rdrFixture))

	// `finalized_at` is declared `provenance = "observed"`. A rule write
	// block naming it is a write to a non-owned tag, decidable from this
	// rule plus the declarations alone.
	t.Run("write block names an observed tag", func(t *testing.T) {
		src := strings.Replace(base,
			"[rule.write]\nstage = \"resolve\"",
			"[rule.write]\nfinalized_at = \"2026-01-01\"\nstage = \"resolve\"", 1)
		if src == base {
			t.Fatal("write block substitution did not apply")
		}

		_, err := table.Load([]byte(src), "adv2-write-observed.toml")
		if err == nil {
			t.Fatal("a rule writing an observed tag loaded clean")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatWriteToNonOwnedTag {
			t.Errorf("category = %q; want %q — the Failure Modes section names "+
				"\"a write to an undeclared or non-owned tag\" as its own stable "+
				"refusal, and `0002:C14` relies on it to guarantee every "+
				"`RequiresOwned` key is owned (got: %v)",
				cat, table.CatWriteToNonOwnedTag, err)
		}
	})

	// The clear list is the write block's other half: `0002:C14` derives
	// `RequiresOwned` from "the rule's write block AND CLEAR LIST", so the
	// same provenance rule binds it.
	t.Run("clear list names an observed tag", func(t *testing.T) {
		src := strings.Replace(base,
			`clear = ["prelock_lens"]`,
			`clear = ["finalized_at"]`, 1)
		if src == base {
			t.Fatal("clear list substitution did not apply")
		}

		_, err := table.Load([]byte(src), "adv2-clear-observed.toml")
		if err == nil {
			t.Fatal("a rule clearing an observed tag loaded clean")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatWriteToNonOwnedTag {
			t.Errorf("category = %q; want %q — a clear is a `<clear>` write "+
				"(`0002:C23`) and carries the same provenance obligation (got: %v)",
				cat, table.CatWriteToNonOwnedTag, err)
		}
	})

	// The discriminating oracle, mirroring
	// `TestReq27_ReaderArityAndWriterProvenanceCarryDistinctCategories`:
	// the writer-ARITY obligation and the write-PROVENANCE obligation are
	// stated separately in the record and MUST NOT collapse into one
	// category. `neg-owned-no-reader.toml` is the arity defect; a rule
	// writing an observed tag is the provenance defect. They currently
	// report the same category, so a caller cannot tell a mis-wired
	// accessor table from a mis-authored rule.
	t.Run("rule provenance does not collapse into accessor arity", func(t *testing.T) {
		arity := loadCategory(t, "neg/neg-owned-no-reader.toml")

		src := strings.Replace(base,
			"[rule.write]\nstage = \"resolve\"",
			"[rule.write]\nfinalized_at = \"2026-01-01\"\nstage = \"resolve\"", 1)
		if src == base {
			t.Fatal("write block substitution did not apply")
		}
		_, err := table.Load([]byte(src), "adv2-collapse.toml")
		if err == nil {
			t.Fatal("a rule writing an observed tag loaded clean")
		}
		provenance, _ := table.CategoryOf(err)

		if arity == provenance {
			t.Errorf("both defects refuse with %q; the accessor-arity "+
				"obligation and the write-to-non-owned-tag obligation are "+
				"separate categories (`0002:C24`, REQ-110), and the Failure "+
				"Modes section requires the rule-level defect to be "+
				"\"decidable from one rule plus the declarations\"", arity)
		}
	})
}

// ---------------------------------------------------------------- ADV-3

// FAILURE MODE 3 — a repeated `in` member mints two rows sharing one
// identity, so the identity tuple stops being total.
//
// Failure Modes: "malformed predicate atoms" fail at load, "each …
// decidable from one rule plus the declarations". `0002:C13` (REQ-72/76)
// makes the dependency explicit:
//
//	"Because RDR 0003 rejects a repeated `in` element at parse, NO TWO ROWS
//	OF ONE RULE SHARE A SUFFIX, WHICH IS WHAT KEEPS THE IDENTITY TUPLE
//	TOTAL."
//
// REQ-64 (`0002:C22`) lands the enforcement here: this RDR owns "the two
// load categories that carry RDR 0003's rejection rules", of which
// `malformed predicate atom` is one. RDR 0003's parser is not what loads
// this document — this package is — so the premise is this package's to
// discharge.
//
// Unmet, the consequences compound and every one of them is silent:
//
//   - REQ-98 (`0002:C19`): rows sort by "(model id, rule id, expansion
//     suffix)". Two rows with equal tuples have no defined order, so the
//     dump's "deterministic across source key order" guarantee is void for
//     that model.
//   - REQ-91: each row retains its source rule id, and expansion count per
//     rule is derived from the dump (Failure Modes). A phantom duplicate
//     inflates that count against a rule the author wrote once.
//   - On an ESCAPE rule the duplicate is worst: `0002:C5` (REQ-41) has each
//     escape row rescue the outcome it binds, and the kernel counts
//     survivors. Two identical escape rows binding one outcome are two
//     survivors — a self-inflicted `ambiguous_match` on the very row that
//     exists to rescue one.
func TestAdv3_RepeatedInMemberIsRefused(t *testing.T) {
	base := string(readFixture(t, rdrFixture))

	t.Run("repeated member on the outcome binding", func(t *testing.T) {
		src := strings.Replace(base,
			`in = ["round-clean", "reconcile-block"]`,
			`in = ["round-clean", "round-clean"]`, 1)
		if src == base {
			t.Fatal("escape outcome substitution did not apply")
		}

		_, err := table.Load([]byte(src), "adv3-dup-outcome-member.toml")
		if err == nil {
			t.Fatal("a repeated `in` member loaded clean; the identity tuple's " +
				"totality rests on this being rejected (`0002:C13`)")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatMalformedPredicateAtom {
			t.Errorf("category = %q; want %q (got: %v)",
				cat, table.CatMalformedPredicateAtom, err)
		}
	})

	t.Run("repeated member on an inherited context match atom", func(t *testing.T) {
		src := strings.Replace(base,
			"[context.large-prelock.match.profile]\nin = [\"large\", \"foundational\"]",
			"[context.large-prelock.match.profile]\nin = [\"large\", \"large\"]", 1)
		if src == base {
			t.Fatal("context substitution did not apply")
		}

		_, err := table.Load([]byte(src), "adv3-dup-context-member.toml")
		if err == nil {
			t.Fatal("a repeated `in` member in a context loaded clean")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatMalformedPredicateAtom {
			t.Errorf("category = %q; want %q (got: %v)",
				cat, table.CatMalformedPredicateAtom, err)
		}
	})
}

// TestAdv3b_RowIdentityIsTotalAcrossTheModel is ADV-3's consequence
// oracle, written so it bites whatever the loader decides to do about the
// repeated member: either the document is refused, or every normalized row
// carries a unique identity. Nothing else satisfies REQ-98's total order.
//
// It is stated over the model rather than the rule so it also covers the
// only other way two rows could tie — `0002:C4`'s duplicate rule id, which
// the loader does refuse.
func TestAdv3b_RowIdentityIsTotalAcrossTheModel(t *testing.T) {
	base := string(readFixture(t, rdrFixture))

	cases := map[string]string{
		"repeated outcome member": strings.Replace(base,
			`in = ["round-clean", "reconcile-block"]`,
			`in = ["round-clean", "round-clean"]`, 1),
		"repeated context member": strings.Replace(base,
			"[context.large-prelock.match.profile]\nin = [\"large\", \"foundational\"]",
			"[context.large-prelock.match.profile]\nin = [\"large\", \"large\"]", 1),
	}

	for name, src := range cases {
		t.Run(name, func(t *testing.T) {
			if src == base {
				t.Fatal("substitution did not apply")
			}
			m, err := table.Load([]byte(src), "adv3b.toml")
			if err != nil {
				// Refusal is the correct outcome; ADV-3 asserts the category.
				return
			}
			seen := map[string]int{}
			for _, r := range m.Rows {
				seen[r.Identity()]++
			}
			for id, n := range seen {
				if n > 1 {
					t.Errorf("identity %q is carried by %d normalized rows; the "+
						"identity tuple (model id, rule id, expansion suffix) "+
						"MUST be total, and the dump's row order is undefined "+
						"while it is not (REQ-98)", id, n)
				}
			}
		})
	}
}
