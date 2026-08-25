package table_test

// RDR 0002 sections I–J: atoms — blocks, identity, and literals — and tag
// declarations with their provenance and type model.

import (
	"reflect"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/resolve"
	"github.com/newcoinc/intrastate/internal/table"
)

// REQ-49: "Guard predicates MUST be represented as positive `all`
// predicates and negative `unless` predicates. Normalization MUST combine
// the match atoms and both guard blocks into one candidate-row predicate
// set before ambiguity checks, and each atom in that set MUST retain the
// key, operator token, literal, and the block it was authored in — a
// three-valued domain, `match`, `all`, or `unless`."
// HAPPY PATH
func TestReq49_OneUnifiedPredicateSetWithBlockRetained(t *testing.T) {
	m := mustLoad(t, rdrFixture)
	row := rowByID(t, m, "rdr.continue-prelock#large")

	want := []table.Atom{
		{Key: "finalized_at", Block: table.BlockAll, Operator: resolve.OpExists, Literal: []string{resolve.LiteralFalse}},
		{Key: "iter", Block: table.BlockAll, Operator: "lt", Literal: []string{"3"}},
		{Key: "profile", Block: table.BlockMatch, Operator: "eq", Literal: []string{"large"}},
		{Key: "profile", Block: table.BlockUnless, Operator: "eq", Literal: []string{"small"}},
		{Key: "stage", Block: table.BlockMatch, Operator: "eq", Literal: []string{"prelock"}},
		{Key: "status", Block: table.BlockMatch, Operator: "eq", Literal: []string{"Draft"}},
	}
	if !reflect.DeepEqual(row.Atoms, want) {
		t.Errorf("unified predicate set =\n %+v\nwant\n %+v", row.Atoms, want)
	}

	// The block domain is exactly three values across both fixtures.
	seen := map[table.Block]bool{}
	for _, rel := range []string{rdrFixture, kataFixture} {
		for _, r := range mustLoad(t, rel).Rows {
			for _, a := range r.Atoms {
				seen[a.Block] = true
			}
		}
	}
	for _, b := range []table.Block{table.BlockMatch, table.BlockAll, table.BlockUnless} {
		if !seen[b] {
			t.Errorf("no normalized atom carries block %q; the fixtures author all three", b)
		}
		delete(seen, b)
	}
	if len(seen) != 0 {
		t.Errorf("atoms carry blocks outside the three-valued domain: %v", seen)
	}
}

// REQ-50: "An implementer building the atom type from RDR 0007 alone gets
// two members and MUST widen to three here."
// BOUNDARY
//
// Satisfied by the shipped kernel (JDR 0001 §D12, landed in 0007 Phase 1):
// this RDR must NOT mint a second, local block type. The assertion is that
// the table's Block values are the kernel's exported constants.
func TestReq50_BlockDomainIsTheKernelsThreeValues(t *testing.T) {
	cases := map[table.Block]resolve.Block{
		table.BlockMatch:  resolve.BlockMatch,
		table.BlockAll:    resolve.BlockAll,
		table.BlockUnless: resolve.BlockUnless,
	}
	for got, want := range cases {
		if string(got) != string(want) {
			t.Errorf("table block %q does not spell the kernel's %q; the atom "+
				"type widens to three members, it does not fork", got, want)
		}
	}

	// The behavioural half: an authored guard block reaches the kernel row
	// carrying the KERNEL's constant, so the widening is real rather than a
	// pair of parallel vocabularies that happen to agree.
	kr := rowByID(t, mustLoad(t, kataFixture), "kata.review-accepted").KernelRow()

	seen := map[resolve.Block]bool{}
	for _, a := range kr.Guard {
		seen[a.Block] = true
	}
	if !seen[resolve.BlockAll] {
		t.Errorf("no kernel guard atom carries resolve.BlockAll: %+v", kr.Guard)
	}
	if !seen[resolve.BlockUnless] {
		t.Errorf("no kernel guard atom carries resolve.BlockUnless: %+v", kr.Guard)
	}
	// And a match-block atom reaches Match rather than being minted as a
	// third kernel block.
	if len(kr.Match) == 0 {
		t.Error("no match-block atom reached the kernel row's Match pattern")
	}
}

// REQ-51: "Normalization MUST NOT fold `unless` atoms into `all` or
// `match` atoms into either guard block."
// ADVERSARIAL
//
// The failing control this replaces: a normalizer dropping every `unless`
// atom passes a determinism-only assertion. Here the counts per block are
// asserted directly on a row authoring all three.
func TestReq51_BlocksAreNotFolded(t *testing.T) {
	m := mustLoad(t, rdrFixture)
	row := rowByID(t, m, "rdr.continue-prelock#large")

	counts := map[table.Block]int{
		table.BlockMatch:  3,
		table.BlockAll:    2,
		table.BlockUnless: 1,
	}
	for b, want := range counts {
		if got := len(atomsInBlock(row, b)); got != want {
			t.Errorf("block %q carries %d atoms; want %d — folding one block "+
				"into another changes these counts: %+v", b, got, want, row.Atoms)
		}
	}

	// `profile` is authored in `match` (via the inherited `in`) and in
	// `unless`; the two must not collapse or migrate.
	prof := atomsOn(row, "profile")
	if len(prof) != 2 {
		t.Fatalf("%d atoms on `profile`; want 2: %+v", len(prof), prof)
	}
	blocks := map[table.Block]string{}
	for _, a := range prof {
		blocks[a.Block] = a.Literal[0]
	}
	want := map[table.Block]string{table.BlockMatch: "large", table.BlockUnless: "small"}
	if !reflect.DeepEqual(blocks, want) {
		t.Errorf("profile atoms by block = %v; want %v", blocks, want)
	}
}

// REQ-52: "one atom authored in **both** `all` and `unless` is two
// distinct atoms and both survive normalization; this is not a load
// failure." / "It MUST NOT be silently pruned at load"
// INPUT EDGE
//
// Asserted positively — the shape must SURVIVE load, not be refused. The
// deadness that follows is RDR 0006's verdict, not this test's.
func TestReq52_SameAtomInAllAndUnlessSurvivesAsTwo(t *testing.T) {
	base := string(readFixture(t, rdrFixture))
	src := strings.Replace(base,
		"[rule.guard.unless.profile]\neq = \"small\"",
		"[rule.guard.unless.profile]\neq = \"small\"\n[rule.guard.all.profile]\neq = \"small\"", 1)
	if src == base {
		t.Fatal("guard block substitution did not apply")
	}

	m, err := table.Load([]byte(src), "all-and-unless.toml")
	if err != nil {
		t.Fatalf("an atom authored in both all and unless refused: %v — this is "+
			"not a load failure", err)
	}

	row := rowByID(t, m, "rdr.continue-prelock#large")
	var small []table.Atom
	for _, a := range atomsOn(row, "profile") {
		if len(a.Literal) == 1 && a.Literal[0] == "small" {
			small = append(small, a)
		}
	}
	if len(small) != 2 {
		t.Fatalf("%d `profile.eq=small` atoms; want 2 — one per block, neither "+
			"silently pruned: %+v", len(small), small)
	}
	blocks := map[table.Block]bool{small[0].Block: true, small[1].Block: true}
	if !blocks[table.BlockAll] || !blocks[table.BlockUnless] {
		t.Errorf("the two atoms sit in blocks %v; want all and unless", blocks)
	}
}

// REQ-53: "For an existence atom the normalizer MUST emit the kernel's
// exported constants verbatim — operator token `OpExists` and literal
// `LiteralTrue` or `LiteralFalse` (RDR 0007, existence-operator clause;
// JDR 0001 §D4) — and MUST reject at load, as a malformed predicate atom,
// any existence literal that is not one of those two forms."
// ADVERSARIAL
//
// The MVV requires BOTH arms: a fixture carrying only one lets a hardcoded
// literal pass. The constants are referenced through internal/resolve
// directly rather than a local mirror.
func TestReq53_ExistenceAtomsEmitTheKernelConstants(t *testing.T) {
	base := string(readFixture(t, rdrFixture))

	t.Run("exists = false", func(t *testing.T) {
		m := mustLoad(t, rdrFixture)
		row := rowByID(t, m, "rdr.continue-prelock#large")
		got := atomsOn(row, "finalized_at")
		if len(got) != 1 {
			t.Fatalf("%d atoms on finalized_at; want 1: %+v", len(got), got)
		}
		if got[0].Operator != resolve.OpExists {
			t.Errorf("operator = %q; want resolve.OpExists (%q)", got[0].Operator, resolve.OpExists)
		}
		if !reflect.DeepEqual(got[0].Literal, []string{resolve.LiteralFalse}) {
			t.Errorf("literal = %v; want [%q]", got[0].Literal, resolve.LiteralFalse)
		}
	})

	t.Run("exists = true", func(t *testing.T) {
		src := strings.Replace(base,
			"[rule.guard.all.finalized_at]\nexists = false",
			"[rule.guard.all.finalized_at]\nexists = true", 1)
		if src == base {
			t.Fatal("exists substitution did not apply")
		}
		m, err := table.Load([]byte(src), "exists-true.toml")
		if err != nil {
			t.Fatalf("exists = true refused: %v", err)
		}
		got := atomsOn(rowByID(t, m, "rdr.continue-prelock#large"), "finalized_at")
		if len(got) != 1 {
			t.Fatalf("%d atoms on finalized_at; want 1: %+v", len(got), got)
		}
		if !reflect.DeepEqual(got[0].Literal, []string{resolve.LiteralTrue}) {
			t.Errorf("literal = %v; want [%q]", got[0].Literal, resolve.LiteralTrue)
		}
	})

	t.Run("non-boolean existence literal refuses", func(t *testing.T) {
		got := loadCategory(t, "neg/neg-bad-exists.toml")
		if got != table.CatMalformedPredicateAtom {
			t.Errorf("category = %q; want %q", got, table.CatMalformedPredicateAtom)
		}
	})
}

// REQ-54: "Tag-key identity is exact byte equality on the post-parse key
// string at every stage — declaration lookup, context and rule predicate
// references, write and clear targets, and the keys a normalized row
// carries. The normalizer performs no case folding, trimming, or namespace
// rewriting"
// ADVERSARIAL
//
// The MVV's negative control: a normalizer applying case folding passes
// every positive test, so the oracle is a mis-cased reference failing
// `unknown tag`. Each named stage gets its own case.
func TestReq54_TagKeyIdentityIsExactByteEqualityAtEveryStage(t *testing.T) {
	base := string(readFixture(t, rdrFixture))

	t.Run("promoted neg-case-folded-tag (context predicate reference)", func(t *testing.T) {
		if got := loadCategory(t, "neg/neg-case-folded-tag.toml"); got != table.CatUnknownTag {
			t.Errorf("category = %q; want %q", got, table.CatUnknownTag)
		}
	})

	stages := map[string][2]string{
		"rule predicate reference": {
			"[rule.guard.unless.profile]\neq = \"small\"",
			"[rule.guard.unless.Profile]\neq = \"small\"",
		},
		"write target": {
			"[rule.write]\nstage = \"resolve\"",
			"[rule.write]\nStage = \"resolve\"",
		},
		"clear target": {
			`clear = ["prelock_lens"]`,
			`clear = ["Prelock_Lens"]`,
		},
		"accessor keys member": {
			`keys = ["finalized_at"]`,
			`keys = ["Finalized_At"]`,
		},
		"[initial] key": {
			"[initial]\nstatus = \"Draft\"",
			"[initial]\nStatus = \"Draft\"",
		},
	}
	for name, pair := range stages {
		t.Run(name, func(t *testing.T) {
			src := strings.Replace(base, pair[0], pair[1], 1)
			if src == base {
				t.Fatalf("substitution %q did not apply", pair[0])
			}
			_, err := table.Load([]byte(src), "case-folded.toml")
			if err == nil {
				t.Fatalf("a mis-cased %s loaded clean; the normalizer performs "+
					"no case folding", name)
			}
			if cat, _ := table.CategoryOf(err); cat != table.CatUnknownTag {
				t.Errorf("category = %q; want %q", cat, table.CatUnknownTag)
			}
		})
	}

	t.Run("no trimming", func(t *testing.T) {
		src := strings.Replace(base, "[rule.guard.unless.profile]", `[rule.guard.unless." profile "]`, 1)
		if src == base {
			t.Fatal("whitespace substitution did not apply")
		}
		_, err := table.Load([]byte(src), "untrimmed.toml")
		if err == nil {
			t.Fatal("a whitespace-padded key loaded clean; the normalizer performs no trimming")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatUnknownTag {
			t.Errorf("category = %q; want %q", cat, table.CatUnknownTag)
		}
	})
}

// REQ-55: "The canonical spelling of a tag is its `[tags.<tag>]`
// declaration key; every reference resolves to a declaration by exact
// match or fails `unknown tag`"
// ADVERSARIAL
func TestReq55_ReferencesResolveToTheDeclarationKey(t *testing.T) {
	t.Run("undeclared match reference", func(t *testing.T) {
		if got := loadCategory(t, "neg/neg-unknown-tag-match.toml"); got != table.CatUnknownTag {
			t.Errorf("category = %q; want %q", got, table.CatUnknownTag)
		}
	})
	t.Run("undeclared guard.unless reference", func(t *testing.T) {
		if got := loadCategory(t, "neg/neg-guard-unless-unknown-tag.toml"); got != table.CatUnknownTag {
			t.Errorf("category = %q; want %q", got, table.CatUnknownTag)
		}
	})
	t.Run("the declaration key is the canonical spelling carried on rows", func(t *testing.T) {
		m := mustLoad(t, rdrFixture)
		for _, r := range m.Rows {
			for _, a := range r.Atoms {
				if _, ok := m.Tags[a.Key]; !ok {
					t.Errorf("%s carries atom key %q with no [tags.%s] declaration",
						r.Identity(), a.Key, a.Key)
				}
			}
		}
	})
}

// REQ-56: "A set-valued literal MUST normalize to an ordered sequence of
// its members, each compared byte-exactly; it MUST NOT be rendered into a
// single delimiter-joined string as its normalized value."
// BOUNDARY
func TestReq56_SetLiteralIsAnOrderedMemberSequence(t *testing.T) {
	// The literal field's Go type is the structural half: a joined value
	// is a string, a member sequence is a slice.
	f, ok := reflect.TypeOf(table.Atom{}).FieldByName("Literal")
	if !ok {
		t.Fatal("table.Atom has no Literal field")
	}
	if f.Type.Kind() != reflect.Slice || f.Type.Elem().Kind() != reflect.String {
		t.Fatalf("Atom.Literal is %v; want []string — a member sequence, never "+
			"a delimiter-joined string", f.Type)
	}

	// A member containing a space normalizes to ONE member, so
	// `["needs work"]` and `["needs", "work"]` are distinct literals.
	m := mustLoad(t, kataFixture)
	row := rowByID(t, m, "kata.review-needs-work")
	got, _ := tagValue(row.Writes, "labels")
	if !reflect.DeepEqual(got, []string{"needs work"}) {
		t.Errorf("labels = %v; want the one-member sequence [\"needs work\"]", got)
	}
}

// REQ-57: "Members sort byte-lexicographically so two authored orderings
// of one set are one literal"
// BOUNDARY
func TestReq57_MembersSortByteLexicographically(t *testing.T) {
	build := func(members string) string {
		return `outcomes = ["done"]

[model]
id = "m"
version = 1

[tags.labels]
provenance = "owned"
kind = "set"
elements = ["a", "b", "Z", "z"]

[tags.tier]
provenance = "owned"
kind = "enum"
domain = ["one"]
single_valued = true

[tags.recognized]
provenance = "recognized"
kind = "enum"
required = true

[read.r]
role = "m"
path = "m"
keys = ["labels", "tier"]
timeout = "2s"

[write.w]
role = "m"
path = "m"
keys = ["labels", "tier"]
timeout = "2s"
read_back = true

[[rule]]
id = "r"
[rule.match.recognized]
eq = "done"
[rule.guard.all.labels]
contains = ` + members + `
[rule.write]
labels = ["a"]
`
	}

	a, err := table.Load([]byte(build(`["z", "Z", "b", "a"]`)), "order-a.toml")
	if err != nil {
		t.Fatalf("order-a refused: %v", err)
	}
	b, err := table.Load([]byte(build(`["a", "b", "Z", "z"]`)), "order-b.toml")
	if err != nil {
		t.Fatalf("order-b refused: %v", err)
	}

	if len(a.Rows) == 0 || len(b.Rows) == 0 {
		t.Fatalf("the ordering fixtures normalized to no rows: a=%d b=%d",
			len(a.Rows), len(b.Rows))
	}
	labels := atomsOn(a.Rows[0], "labels")
	if len(labels) != 1 {
		t.Fatalf("%d atoms on `labels`; want 1: %+v", len(labels), labels)
	}
	got := labels[0].Literal
	// Byte-lexicographic: uppercase Z (0x5A) sorts before lowercase a (0x61).
	want := []string{"Z", "a", "b", "z"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("member order = %v; want the byte-lexicographic %v", got, want)
	}
	if !reflect.DeepEqual(a.Rows, b.Rows) {
		t.Error("two authored orderings of one set normalized to different " +
			"literals; they are one literal")
	}
}

// REQ-58: "Wherever atoms are compared, keyed, deduplicated, or sorted —
// in particular the **merge key** … the literal field MUST be compared as
// the member sequence, element by element. An implementation MUST NOT
// derive that key, or any other identity, by joining members into a
// string."
// ADVERSARIAL
//
// The merge key is exercised by the delimiter-parameterized control in
// REQ-46; here the SORT is exercised, since a joined sort key orders
// differently from an element-by-element one.
//
// The sort only reaches the literal when (key, block, operator) TIE, since
// the literal is the tuple's LAST element. Two atoms differing in block —
// one under `guard.all`, one under `guard.unless` — are already ordered by
// block, so the literal is never consulted and a joined key passes. This
// control therefore ties all three leading fields: two `[context]` blocks
// each author a match `in` on the SAME tag, which merge into one row's
// match block as two atoms agreeing on (labels, match, in).
//
// The observable is the EXPANSION SUFFIX element order, which `expand`
// takes in the atoms' sort order. Element-by-element, ["a","z"] < ["a!q","b"]
// (first members "a" < "a!q" by prefix). Comma-joined, "a!q,b" < "a,z"
// ('!' 0x21 < ',' 0x2c at the second byte) — the opposite order. So a
// joined sort key transposes every suffix.
func TestReq58_LiteralComparesAsAMemberSequenceNotAJoinedString(t *testing.T) {
	src := `outcomes = ["done"]

[model]
id = "m"
version = 1

[tags.labels]
provenance = "owned"
kind = "set"
elements = ["a", "b", "z", "a!q"]

[tags.tier]
provenance = "owned"
kind = "enum"
domain = ["one"]
single_valued = true

[tags.recognized]
provenance = "recognized"
kind = "enum"
required = true

[read.r]
role = "m"
path = "m"
keys = ["labels", "tier"]
timeout = "2s"

[write.w]
role = "m"
path = "m"
keys = ["labels", "tier"]
timeout = "2s"
read_back = true

[context.c1.match.recognized]
eq = "done"

[context.seq-a.match.labels]
in = ["a", "z"]

[context.seq-b.match.labels]
in = ["a!q", "b"]

[[rule]]
id = "r"
use = ["c1", "seq-a", "seq-b"]
# REQ-35 requires a local match block on every rule. The tier atom is chosen
# so it ties with neither in-atom on labels, leaving the sort control intact.
[rule.match.tier]
eq = "one"
[rule.write]
labels = ["a"]
`
	m, err := table.Load([]byte(src), "seq-identity.toml")
	if err != nil {
		t.Fatalf("refused: %v", err)
	}

	// The two same-key match atoms tie on (key, block, operator) and merge
	// over the FULL identity (REQ-46), so both survive: 2x2 members.
	if len(m.Rows) != 4 {
		t.Fatalf("%d rows; want 4 — the two same-block `in` atoms must both "+
			"survive the merge and expand as a product", len(m.Rows))
	}

	// seq-a's member leads every suffix, because ["a","z"] sorts before
	// ["a!q","b"] element by element. A comma-joined sort key reverses it.
	var got []string
	for _, r := range m.Rows {
		if len(r.Suffix) != 2 {
			t.Fatalf("suffix %v carries %d elements; want 2", r.Suffix, len(r.Suffix))
		}
		got = append(got, strings.Join(r.Suffix, "|"))
	}
	want := []string{"a|a!q", "a|b", "z|a!q", "z|b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expansion suffixes = %v; want %v — the literal must be "+
			"compared as a member sequence, not as a joined string", got, want)
	}
}

// REQ-59: "Producers rendering into that field MUST spell a set visibly as
// a set **with each member delimited unambiguously** — quoting each member
// is the cheap way"
// BOUNDARY
//
// Stage 6's `|` case is the witness: an unquoted separator rendered
// ["x|y","z"] and ["x","y|z"] alike. The oracle is therefore that the two
// spellings RENDER differently, not merely that they normalize
// differently.
func TestReq59_SetValuedFieldsRenderWithMembersUnambiguouslyDelimited(t *testing.T) {
	for _, d := range []string{"comma", "semi", "pipe", "space", "empty"} {
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
			if table.Dump(a) == table.Dump(b) {
				t.Errorf("the two set spellings render identically under the %s "+
					"delimiter; each member must be delimited unambiguously", d)
			}
		})
	}
}

// REQ-60: "**Rule ids, outcome literals, and the members of any
// match-block `in` literal MUST NOT contain the expansion-suffix separator
// `#`.** … Violations are a malformed rule id, a malformed recognized
// outcome alphabet, and a malformed predicate atom respectively."
// ADVERSARIAL
//
// Three sites, three DISTINCT categories: an implementation collapsing
// them into one fails the category oracle.
func TestReq60_ExpansionSeparatorIsBannedAtThreeSitesWithThreeCategories(t *testing.T) {
	cases := map[string]table.Category{
		"neg/neg-ruleid-hash.toml":    table.CatMalformedRuleID,
		"neg/neg-alphabet-hash.toml":  table.CatMalformedRecognizedOutcomeAlphabet,
		"neg/neg-in-member-hash.toml": table.CatMalformedPredicateAtom,
	}
	seen := map[table.Category]string{}
	for rel, want := range cases {
		t.Run(rel, func(t *testing.T) {
			got := loadCategory(t, rel)
			if got != want {
				t.Errorf("category = %q; want %q", got, want)
			}
			if prev, dup := seen[got]; dup {
				t.Errorf("%s and %s both refused %q; the three sites carry three "+
					"distinct categories", prev, rel, got)
			}
			seen[got] = rel
		})
	}
}

// REQ-61: "The string `<clear>` MUST be refused at load wherever a tag
// value is authored — a write-block value, an `[initial]` assignment, or a
// predicate literal — as a `reserved tag value`."
// ADVERSARIAL
func TestReq61_ClearSentinelIsReservedAtEveryAuthoringSite(t *testing.T) {
	sites := map[string]string{
		"write-block value":    "neg/neg-clear-as-write.toml",
		"[initial] value":      "neg/neg-clear-in-initial.toml",
		"guard.all literal":    "neg/neg-guard-all-clear.toml",
		"guard.unless literal": "neg/neg-guard-unless-clear.toml",
	}
	for name, rel := range sites {
		t.Run(name, func(t *testing.T) {
			if got := loadCategory(t, rel); got != table.CatReservedTagValue {
				t.Errorf("category = %q; want %q", got, table.CatReservedTagValue)
			}
		})
	}

	t.Run("match-block literal", func(t *testing.T) {
		base := string(readFixture(t, rdrFixture))
		src := strings.Replace(base, "[context.draft.match.status]\neq = \"Draft\"",
			"[context.draft.match.status]\neq = \"<clear>\"", 1)
		if src == base {
			t.Fatal("match literal substitution did not apply")
		}
		_, err := table.Load([]byte(src), "clear-in-match.toml")
		if err == nil {
			t.Fatal("a <clear> match literal loaded clean")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatReservedTagValue {
			t.Errorf("category = %q; want %q", cat, table.CatReservedTagValue)
		}
	})

	t.Run("the sentinel string is exactly <clear>", func(t *testing.T) {
		if table.ClearSentinel != "<clear>" {
			t.Errorf("ClearSentinel = %q; want %q", table.ClearSentinel, "<clear>")
		}
	})
}

// REQ-62: "The model MUST declare every tag it matches or writes, including
// each tag's provenance: owned, observed, or recognized."
// HAPPY PATH
func TestReq62_EveryMatchedOrWrittenTagIsDeclaredWithProvenance(t *testing.T) {
	for _, rel := range []string{rdrFixture, kataFixture} {
		t.Run(rel, func(t *testing.T) {
			m := mustLoad(t, rel)

			legal := map[table.Provenance]bool{
				table.ProvenanceOwned:      true,
				table.ProvenanceObserved:   true,
				table.ProvenanceRecognized: true,
			}
			for key, decl := range m.Tags {
				if !legal[decl.Provenance] {
					t.Errorf("tag %q provenance = %q; want owned, observed, or "+
						"recognized", key, decl.Provenance)
				}
			}
			for _, r := range m.Rows {
				for _, a := range r.Atoms {
					if _, ok := m.Tags[a.Key]; !ok {
						t.Errorf("matched tag %q is undeclared", a.Key)
					}
				}
				for _, w := range r.Writes {
					if _, ok := m.Tags[w.Key]; !ok {
						t.Errorf("written tag %q is undeclared", w.Key)
					}
				}
			}
		})
	}
}

// REQ-63: "A tag declaration also carries its **type model**, spelled as
// the wire keys `kind`, `domain`, `min`, `max`, `elements`,
// `single_valued`, and `required` (JDR 0001 §D7(iii); `required` is the
// optionality marker and defaults to optional)."
// BOUNDARY
func TestReq63_TypeModelWireKeys(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	t.Run("enum with a domain", func(t *testing.T) {
		d := m.Tags["status"]
		if d.Kind != "enum" {
			t.Errorf("kind = %q; want enum", d.Kind)
		}
		if !reflect.DeepEqual(d.Domain, []string{"Draft", "Final"}) {
			t.Errorf("domain = %v; want [Draft Final]", d.Domain)
		}
		if !d.SingleValued || !d.Required {
			t.Errorf("single_valued=%v required=%v; want both true", d.SingleValued, d.Required)
		}
	})

	t.Run("int with min and max", func(t *testing.T) {
		d := m.Tags["iter"]
		if d.Kind != "int" {
			t.Errorf("kind = %q; want int", d.Kind)
		}
		if d.Min == nil || *d.Min != 0 {
			t.Errorf("min = %v; want 0", d.Min)
		}
		if d.Max == nil || *d.Max != 9 {
			t.Errorf("max = %v; want 9", d.Max)
		}
	})

	t.Run("required defaults to optional", func(t *testing.T) {
		// `cluster_ready` declares no `required` key.
		if d := m.Tags["cluster_ready"]; d.Required {
			t.Error("cluster_ready is required; `required` defaults to optional")
		}
	})

	t.Run("set with elements", func(t *testing.T) {
		k := mustLoad(t, kataFixture)
		d := k.Tags["labels"]
		if d.Kind != "set" {
			t.Errorf("kind = %q; want set", d.Kind)
		}
		if !reflect.DeepEqual(d.Elements, []string{"bug", "chore", "needs work"}) {
			t.Errorf("elements = %v; want the authored universe", d.Elements)
		}
	})
}

// REQ-64: "This RDR owns where a declaration is authored — under
// `[tags.<tag>]`, beside `provenance`, with no accessor reference — and
// the two load categories that carry RDR 0003's rejection rules:
// `malformed tag declaration` … and `malformed predicate atom` (a literal
// outside the tag's declared domain …)."
// ADVERSARIAL
//
// Per deviations.md D1 the kind vocabulary is RDR 0003's five tokens
// (`enum`, `bool`, `int`, `set`, `scalar`); `string` is not one, so the
// promoted fixtures were rewritten and an unknown token refuses. D3 files
// write-block value conformance under the same declaration category.
func TestReq64_DeclarationSiteAndTheTwoRejectionCategories(t *testing.T) {
	base := string(readFixture(t, rdrFixture))

	t.Run("no accessor reference is admitted on a declaration", func(t *testing.T) {
		src := strings.Replace(base,
			"[tags.cluster_ready]\nprovenance = \"owned\"",
			"[tags.cluster_ready]\nprovenance = \"owned\"\naccessor = \"rdr-status\"", 1)
		if src == base {
			t.Fatal("declaration substitution did not apply")
		}
		_, err := table.Load([]byte(src), "tag-accessor.toml")
		if err == nil {
			t.Fatal("a tag-side accessor reference loaded clean")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatUnknownSchemaField {
			t.Errorf("category = %q; want %q", cat, table.CatUnknownSchemaField)
		}
	})

	t.Run("kind outside RDR 0003's five tokens (D1)", func(t *testing.T) {
		src := strings.Replace(base, `kind = "scalar"`, `kind = "string"`, 1)
		if src == base {
			t.Fatal("kind substitution did not apply")
		}
		_, err := table.Load([]byte(src), "kind-string.toml")
		if err == nil {
			t.Fatal("kind = \"string\" loaded clean; RDR 0003's vocabulary is " +
				"exactly enum, bool, int, set, scalar")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatMalformedTagDeclaration {
			t.Errorf("category = %q; want %q", cat, table.CatMalformedTagDeclaration)
		}
	})

	t.Run("all five kind tokens are admitted", func(t *testing.T) {
		for _, kind := range []string{"enum", "bool", "int", "set", "scalar"} {
			if !table.IsDeclaredKind(kind) {
				t.Errorf("kind %q is not admitted; RDR 0003 fixes exactly five tokens", kind)
			}
		}
		if table.IsDeclaredKind("string") {
			t.Error("kind \"string\" is admitted; it is not one of RDR 0003's five tokens")
		}
	})

	t.Run("declaration RDR 0003 rejects: domain on a kind admitting none", func(t *testing.T) {
		got := loadCategory(t, "neg/neg-tagdecl-domain-on-bool.toml")
		if got != table.CatMalformedTagDeclaration {
			t.Errorf("category = %q; want %q", got, table.CatMalformedTagDeclaration)
		}
	})

	// `single_valued` is meaningful only where the kind has both
	// single-valued and non-single-valued assignment counts. `0003:1273-1283`
	// fixes a `set` at `2^|element universe|` — "never `|universe|`" — and
	// lists no row for `scalar`, a kind with no finite declared domain, so on
	// those two the marker asserts nothing. `0003:1389` routes such a
	// declaration to the loader; `0002:C22` (REQ-64) owns the category.
	t.Run("declaration RDR 0003 rejects: single_valued on a kind admitting none",
		func(t *testing.T) {
			for _, tc := range []struct {
				rel    string
				marker string
			}{
				{"neg/neg-tagdecl-single-valued-on-set.toml", "elements = [\"bug\", \"chore\", \"needs work\"]\nsingle_valued = true"},
				{"neg/neg-tagdecl-single-valued-on-scalar.toml", "kind = \"scalar\"\nsingle_valued = true"},
			} {
				if got := loadCategory(t, tc.rel); got != table.CatMalformedTagDeclaration {
					t.Errorf("%s: category = %q; want %q", tc.rel, got, table.CatMalformedTagDeclaration)
				}

				authored := string(readFixture(t, tc.rel))
				src := strings.Replace(authored, tc.marker, strings.TrimSuffix(tc.marker, "true")+"false", 1)
				if src == authored {
					t.Fatalf("%s: target single_valued marker not found", tc.rel)
				}
				_, err := table.Load([]byte(src), tc.rel)
				if got, ok := table.CategoryOf(err); !ok || got != table.CatMalformedTagDeclaration {
					t.Errorf("%s with single_valued = false: category = %q, %v; want %q",
						tc.rel, got, err, table.CatMalformedTagDeclaration)
				}
			}
		})

	t.Run("predicate literal outside the declared domain", func(t *testing.T) {
		got := loadCategory(t, "neg/neg-literal-outside-domain.toml")
		if got != table.CatMalformedPredicateAtom {
			t.Errorf("category = %q; want %q", got, table.CatMalformedPredicateAtom)
		}
	})

	t.Run("D3: write-block value conformance", func(t *testing.T) {
		for _, rel := range []string{
			"neg/neg-write-value-outside-domain.toml",
			"neg/neg-write-value-wrong-kind.toml",
		} {
			if got := loadCategory(t, rel); got != table.CatMalformedTagDeclaration {
				t.Errorf("%s: category = %q; want %q — deviations.md D3 files "+
					"write-value conformance under the declaration category",
					rel, got, table.CatMalformedTagDeclaration)
			}
		}
	})
}

// REQ-65: "Normalization MUST carry every declared field through to the
// normalized model without loss"
// BOUNDARY
//
// A loader dropping `min`, `max`, `elements`, or `single_valued` passes
// every predicate test, so the oracle is field-by-field equality against
// the authored declarations.
func TestReq65_EveryDeclaredFieldSurvivesNormalization(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	zero := 0
	nine := 9
	want := map[string]table.TagDecl{
		"iter": {
			Provenance: table.ProvenanceOwned,
			Kind:       "int",
			Min:        &zero,
			Max:        &nine,
			Required:   true,
		},
		"cluster_ready": {
			Provenance: table.ProvenanceOwned,
			Kind:       "bool",
		},
		"prelock_lens": {
			Provenance:   table.ProvenanceOwned,
			Kind:         "enum",
			Domain:       []string{"grounding", "cove", "3amigo", "critique", "repeatability"},
			SingleValued: true,
		},
		"finalized_at": {
			Provenance: table.ProvenanceObserved,
			Kind:       "scalar",
		},
	}
	for key, w := range want {
		got, ok := m.Tags[key]
		if !ok {
			t.Errorf("declaration for %q is missing", key)
			continue
		}
		if !reflect.DeepEqual(got, w) {
			t.Errorf("declaration %q =\n %+v\nwant\n %+v", key, got, w)
		}
	}
}

// REQ-66: "A tag declaration with provenance `recognized` MUST be named
// `recognized`, and no owned or observed declaration may take that name;
// violations fail in RDR 0008's `reserved_tag_key` category"
// ADVERSARIAL
func TestReq66_RecognizedIsTheReservedName(t *testing.T) {
	base := string(readFixture(t, rdrFixture))

	t.Run("a recognized declaration under another name", func(t *testing.T) {
		got := loadCategory(t, "neg/neg-recognized-misnamed.toml")
		if got != table.CatReservedTagKey {
			t.Errorf("category = %q; want %q", got, table.CatReservedTagKey)
		}
	})

	t.Run("an owned declaration named recognized", func(t *testing.T) {
		src := strings.Replace(base,
			"[tags.recognized]\nprovenance = \"recognized\"",
			"[tags.recognized]\nprovenance = \"owned\"", 1)
		if src == base {
			t.Fatal("provenance substitution did not apply")
		}
		_, err := table.Load([]byte(src), "owned-recognized.toml")
		if err == nil {
			t.Fatal("an owned declaration named `recognized` loaded clean")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatReservedTagKey {
			t.Errorf("category = %q; want %q", cat, table.CatReservedTagKey)
		}
	})

	t.Run("the reserved key spells the kernel's", func(t *testing.T) {
		if table.RecognizedTagKey != "recognized" {
			t.Errorf("RecognizedTagKey = %q; want %q", table.RecognizedTagKey, "recognized")
		}
	})
}

// REQ-67: "the `outcomes` alphabet MUST be non-empty, duplicate-free, and
// MUST NOT contain the empty string."
// ADVERSARIAL
//
// Each arm is distinguishable within the one category (REQ-106), so all
// three are exercised.
func TestReq67_OutcomeAlphabetIsNonEmptyDuplicateFreeAndHasNoEmptyMember(t *testing.T) {
	base := string(readFixture(t, rdrFixture))

	t.Run("promoted duplicate member", func(t *testing.T) {
		got := loadCategory(t, "neg/neg-dup-alphabet-member.toml")
		if got != table.CatMalformedRecognizedOutcomeAlphabet {
			t.Errorf("category = %q; want %q", got, table.CatMalformedRecognizedOutcomeAlphabet)
		}
	})
	t.Run("promoted empty-string member", func(t *testing.T) {
		got := loadCategory(t, "neg/neg-empty-alphabet-member.toml")
		if got != table.CatMalformedRecognizedOutcomeAlphabet {
			t.Errorf("category = %q; want %q", got, table.CatMalformedRecognizedOutcomeAlphabet)
		}
	})
	t.Run("empty alphabet", func(t *testing.T) {
		src := strings.Replace(base,
			`outcomes = ["round-clean", "verdict-flapping", "reconcile-block", "finalized"]`,
			`outcomes = []`, 1)
		if src == base {
			t.Fatal("alphabet substitution did not apply")
		}
		_, err := table.Load([]byte(src), "empty-alphabet.toml")
		if err == nil {
			t.Fatal("an empty alphabet loaded clean")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatMalformedRecognizedOutcomeAlphabet {
			t.Errorf("category = %q; want %q", cat, table.CatMalformedRecognizedOutcomeAlphabet)
		}
	})
	t.Run("absent alphabet is a different category", func(t *testing.T) {
		// REQ-106 names `missing recognized outcome alphabet` separately
		// from `malformed recognized outcome alphabet`.
		got := loadCategory(t, "neg/neg-no-alphabet.toml")
		if got != table.CatMissingRecognizedOutcomeAlphabet {
			t.Errorf("category = %q; want %q", got, table.CatMissingRecognizedOutcomeAlphabet)
		}
	})
}
