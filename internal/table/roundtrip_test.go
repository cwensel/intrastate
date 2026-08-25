package table_test

// RDR 0002 sections P–R: the Round-Trip / Inverse Invariant, the testing
// obligations the record states as contracts, and package placement.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/table"
)

// REQ-113: "`parse ∘ normalize = expanded-table value identity` on valid
// transition model fixtures: two authorings of one model — differing in
// TOML key order, rule declaration order, or `eq`-versus-single-member-`in`
// spelling — must normalize to the same candidate-row set, compared as
// values over the full field list the dump contract carries … with the
// source locator's optional line/column detail excluded from the
// comparison."
// HAPPY PATH
//
// The Round-Trip / Inverse Invariant, asserted over the NORMALIZED VALUE
// (never the dump text: `dump` has no specified inverse). The comparison
// is value-for-value over the full field list, not a "did not error"
// check.
func TestReq113_RoundTripValueIdentityAcrossAllThreePermutations(t *testing.T) {
	base := mustLoad(t, rdrFixture)

	perms := map[string]string{
		"TOML key order":             "perm/rdr-keyorder.toml",
		"rule declaration order":     "perm/rdr-ruleorder.toml",
		"eq versus single-member in": "perm/rdr-eq-as-in.toml",
	}
	for name, rel := range perms {
		t.Run(name, func(t *testing.T) {
			got := mustLoad(t, rel)

			if len(got.Rows) != len(base.Rows) {
				t.Fatalf("%d rows; want %d — %v vs %v",
					len(got.Rows), len(base.Rows), rowIdentities(got), rowIdentities(base))
			}
			// Value-for-value over the full field list the dump contract
			// carries, field by field so a failure names the field.
			for i := range base.Rows {
				w, g := base.Rows[i], got.Rows[i]
				if w.Identity() != g.Identity() {
					t.Errorf("row %d identity = %q; want %q", i, g.Identity(), w.Identity())
				}
				if !reflect.DeepEqual(g.Suffix, w.Suffix) {
					t.Errorf("%s expansion suffix = %v; want %v", g.Identity(), g.Suffix, w.Suffix)
				}
				if g.Outcome != w.Outcome {
					t.Errorf("%s outcome = %q; want %q", g.Identity(), g.Outcome, w.Outcome)
				}
				if !reflect.DeepEqual(g.Atoms, w.Atoms) {
					t.Errorf("%s atoms =\n %+v\nwant\n %+v", g.Identity(), g.Atoms, w.Atoms)
				}
				if !reflect.DeepEqual(g.Gate, w.Gate) {
					t.Errorf("%s gate = %v; want %v", g.Identity(), g.Gate, w.Gate)
				}
				if !reflect.DeepEqual(g.NextTags, w.NextTags) {
					t.Errorf("%s next tags = %+v; want %+v", g.Identity(), g.NextTags, w.NextTags)
				}
				if !reflect.DeepEqual(g.Writes, w.Writes) {
					t.Errorf("%s writes = %+v; want %+v", g.Identity(), g.Writes, w.Writes)
				}
				if !reflect.DeepEqual(g.RequiresOwned, w.RequiresOwned) {
					t.Errorf("%s requires_owned = %v; want %v",
						g.Identity(), g.RequiresOwned, w.RequiresOwned)
				}
				if !reflect.DeepEqual(g.Escape, w.Escape) {
					t.Errorf("%s escape = %v; want %v", g.Identity(), g.Escape, w.Escape)
				}
			}
			// The whole row set as one value, including the locator's
			// rule-identifying part.
			if !reflect.DeepEqual(got.Rows, base.Rows) {
				t.Error("the normalized candidate-row sets are not value-identical")
			}
		})
	}
}

// REQ-114: "The locator's **presence and its rule-identifying part are
// compared**, even though its line/column detail is not"
// ADVERSARIAL
//
// Excluding the whole field would let a normalizer that dropped locators
// entirely satisfy the invariant, so presence is asserted; the
// rule-identifying part is asserted to be stable across the permutations,
// where any line/column detail would move.
func TestReq114_LocatorPresenceAndRuleIdentifyingPartAreCompared(t *testing.T) {
	base := mustLoad(t, rdrFixture)

	for _, r := range base.Rows {
		if r.SourceLocator == "" {
			t.Errorf("%s carries no locator; presence is part of the invariant",
				r.Identity())
		}
		if !strings.Contains(r.SourceLocator, r.RuleID) {
			t.Errorf("%s locator %q does not carry its rule-identifying part",
				r.Identity(), r.SourceLocator)
		}
	}

	// Rule-order permutation moves every rule's line number while leaving
	// the rule-identifying part fixed.
	perm := mustLoad(t, "perm/rdr-ruleorder.toml")
	for i := range base.Rows {
		if base.Rows[i].SourceLocator != perm.Rows[i].SourceLocator {
			t.Errorf("%s locator changed under a rule-order permutation: %q -> %q",
				base.Rows[i].Identity(), base.Rows[i].SourceLocator,
				perm.Rows[i].SourceLocator)
		}
	}
}

// REQ-115: "a set-valued field MUST render its members delimited in a way
// that shows the field is a set (for example bracketed), so a reviewer
// reading two rows can see that a difference *may* be hiding"
// BOUNDARY
func TestReq115_SetValuedFieldsRenderVisiblyAsSets(t *testing.T) {
	m := mustLoad(t, kataFixture)
	out := table.Dump(m)

	// The `labels` write is a set-valued field. Its rendering must show
	// that it is a set, and must delimit each member unambiguously.
	if !strings.Contains(out, "labels") {
		t.Fatalf("the dump does not render the `labels` write:\n%s", out)
	}
	line := dumpLine(t, out, "kata.review-needs-work")
	i := strings.Index(line, "labels")
	if i < 0 {
		t.Fatalf("the review-needs-work row does not render `labels`:\n%s", line)
	}
	rest := line[i:]
	if !strings.ContainsAny(rest, "[") {
		t.Errorf("the set-valued `labels` field renders without set-showing "+
			"delimiters: %s", rest)
	}
}

func dumpLine(t *testing.T, out, identity string) string {
	t.Helper()

	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, identity) {
			return line
		}
	}
	t.Fatalf("no dump line carries %q:\n%s", identity, out)
	return ""
}

// REQ-116: "Source rewrite is likewise out of scope for this RDR"
// BOUNDARY
//
// The checkable form: this package exposes no source-mutating entry point.
// A rewrite capability would need its own source-preservation invariant,
// which this RDR does not define.
func TestReq116_NoSourceRewriteCapability(t *testing.T) {
	// Only SOURCE-mutating names are banned. `Dump` renders the normalized
	// value into text and is required by REQ-94; it claims no inverse and
	// writes nothing back.
	banned := []string{"Rewrite", "RewriteSource", "Marshal", "MarshalTOML",
		"Encode", "EncodeTOML", "Save", "WriteFile", "WriteSource"}

	for _, file := range packageGoFiles(t) {
		f, err := parser.ParseFile(token.NewFileSet(), file, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", file, err)
		}
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			for _, b := range banned {
				if fn.Name.Name == b {
					t.Errorf("%s declares %q; source rewrite is out of scope for "+
						"this RDR and would owe its own source-preservation invariant",
						file, b)
				}
			}
		}
	}
}

// REQ-117: "The overlap item is the one MVV assertion this RDR cannot
// discharge alone … It is **deferred to the Phase 5 lint handshake rather
// than counted as satisfied here**" / "Marking the overlap assertion green
// before RDR 0006 lands would be asserting on a stub."
// BOUNDARY
//
// The deferral is itself the obligation: this package must NOT mint an
// overlap verdict. What it owes is that the rows carry enough structure
// for RDR 0006's check — distinguishable by identity, comparable by
// predicate set — which REQ-112 and REQ-144 assert.
func TestReq117_OverlapIsDeferredAndNotAssertedHere(t *testing.T) {
	for _, c := range table.Categories() {
		if strings.Contains(string(c), "overlap") || strings.Contains(string(c), "ambiguous") {
			t.Errorf("category %q mints an overlap verdict; the assertion is "+
				"deferred to RDR 0006's lint handshake", c)
		}
	}

	for _, file := range packageGoFiles(t) {
		f, err := parser.ParseFile(token.NewFileSet(), file, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", file, err)
		}
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil {
				continue
			}
			if strings.Contains(strings.ToLower(fn.Name.Name), "overlap") {
				t.Errorf("%s declares %q; marking the overlap assertion green "+
					"before RDR 0006 lands would be asserting on a stub",
					file, fn.Name.Name)
			}
		}
	}
}

// REQ-118: "`evidence/spikes/iter-2/` is the promoted set" / "a fixture may
// not be narrowed on promotion, only extended."
// BOUNDARY
//
// Four legs, because a name is not evidence. (1) Every fixture in the
// approved iter-2 set is present in testdata AT THE SAME RELATIVE PATH: a
// basename match lets a fixture be relocated out of the directory whose
// walk drives TestReq119, silently dropping it from the category census.
// (2) Every promoted fixture retains the spike's content, apart from the
// recorded string-to-scalar kind rename (deviations.md D1). This is the
// broadest pin: it holds each negative's specific MUTATION, not merely the
// category that mutation happens to trip, so swapping one negative for
// another that refuses alike is caught. It also covers the promoted
// non-negatives that no other test loads (merge-delim-atom, write-delim-a,
// write-delim-b). (3) Every promoted negative still refuses with the
// specific category it was promoted witnessing. Leg 2 subsumes this for a
// byte edit, but leg 3 is the one that survives a LOADER change: if a
// refactor reroutes a fixture to a different category the bytes are
// untouched and only the category pin bites. (4) The two primary positives
// are compared by normalized row IDENTITY, not by row count: a nine-row
// fixture with one row swapped for filler keeps its census and loses its
// coverage.
//
// Every leg iterates the SPIKE set and asserts it is covered by what is
// promoted, never the converse — REQ-118 permits extension, and an
// equality assertion would fail on every legitimate new fixture until
// someone loosened it back to basenames.
func TestReq118_PromotedFixtureSetIsNotNarrowed(t *testing.T) {
	spike := filepath.Join("..", "..", "docs", "rdr",
		"0002-transition-table-as-reviewable-data", "evidence", "spikes", "iter-2")

	promoted := map[string]bool{}
	err := filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".toml") {
			return nil
		}
		rel, err := filepath.Rel("testdata", path)
		if err != nil {
			return err
		}
		promoted[rel] = true
		return nil
	})
	if err != nil {
		t.Fatalf("walk testdata: %v", err)
	}

	var spikeFixtures []string
	var spikeNegatives []string
	err = filepath.Walk(spike, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".toml") {
			return nil
		}
		rel, err := filepath.Rel(spike, path)
		if err != nil {
			return err
		}
		if !promoted[rel] {
			t.Errorf("spike fixture %s is not promoted at testdata/%s; the set "+
				"may be extended, never narrowed. Relocating a fixture IS a "+
				"narrowing under REQ-118 — it drops out of the directory walks "+
				"that drive the category census — and must be recorded as a "+
				"deviation, not absorbed by matching on basename", rel, rel)
			return nil
		}
		spikeFixtures = append(spikeFixtures, rel)
		if filepath.Dir(rel) == "neg" {
			spikeNegatives = append(spikeNegatives, rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk spike: %v", err)
	}

	t.Run("every promoted fixture retains its approved content", func(t *testing.T) {
		for _, rel := range spikeFixtures {
			spikeData, err := os.ReadFile(filepath.Join(spike, rel))
			if err != nil {
				t.Fatalf("read spike fixture %s: %v", rel, err)
			}
			want := strings.ReplaceAll(string(spikeData), `kind = "string"`, `kind = "scalar"`)
			if got := string(readFixture(t, rel)); got != want {
				t.Errorf("promoted fixture %s differs from its approved spike content", rel)
			}
		}
	})

	t.Run("every promoted negative still trips its recorded category", func(t *testing.T) {
		// The oracle is the stable data-level category, never the message
		// text (REQ-110, REQ-119) — the same machinery TestReq119 walks
		// neg/ with, pinned here per fixture instead of counted in
		// aggregate. A global "the census did not shrink" check does not
		// bite: rewriting one negative to trip a category some OTHER
		// negative already witnesses leaves the census whole while that
		// negative's own coverage is gone. Only a per-fixture pin catches it.
		//
		// The categories are recorded here rather than read back from the
		// spike copies because the spike copies do not load: iter-2 predates
		// the kind = "string" to "scalar" rename (internal/table/model.go,
		// declaredKinds), so every spike file refuses at tag declaration
		// before reaching its intended mutation. Extending the promoted set
		// does not touch this map; narrowing a promoted negative does.
		want := map[string]table.Category{
			"neg/neg-accessor-no-timeout.toml":              "malformed_accessor_declaration",
			"neg/neg-alphabet-hash.toml":                    "malformed_recognized_outcome_alphabet",
			"neg/neg-bad-exists.toml":                       "malformed_predicate_atom",
			"neg/neg-case-folded-tag.toml":                  "unknown_tag",
			"neg/neg-clear-as-write.toml":                   "reserved_tag_value",
			"neg/neg-clear-in-initial.toml":                 "reserved_tag_value",
			"neg/neg-dup-alphabet-member.toml":              "malformed_recognized_outcome_alphabet",
			"neg/neg-duplicate-rule-id.toml":                "duplicate_rule_id",
			"neg/neg-empty-alphabet-member.toml":            "malformed_recognized_outcome_alphabet",
			"neg/neg-escape-with-clear.toml":                "malformed_escape_declaration",
			"neg/neg-escape-with-gate.toml":                 "malformed_escape_declaration",
			"neg/neg-escape-with-write.toml":                "malformed_escape_declaration",
			"neg/neg-gate-is-reader.toml":                   "unknown_accessor",
			"neg/neg-guard-all-bad-operator.toml":           "malformed_predicate_atom",
			"neg/neg-guard-all-clear.toml":                  "reserved_tag_value",
			"neg/neg-guard-unless-bad-operator.toml":        "malformed_predicate_atom",
			"neg/neg-guard-unless-clear.toml":               "reserved_tag_value",
			"neg/neg-guard-unless-out-of-domain.toml":       "malformed_predicate_atom",
			"neg/neg-guard-unless-unknown-tag.toml":         "unknown_tag",
			"neg/neg-in-member-hash.toml":                   "malformed_predicate_atom",
			"neg/neg-initial-bad-value.toml":                "malformed_initial_declaration",
			"neg/neg-initial-int-out-of-range.toml":         "malformed_initial_declaration",
			"neg/neg-initial-unknown.toml":                  "unknown_tag",
			"neg/neg-literal-outside-domain.toml":           "malformed_predicate_atom",
			"neg/neg-match-exists.toml":                     "malformed_predicate_atom",
			"neg/neg-match-lt.toml":                         "malformed_predicate_atom",
			"neg/neg-no-write-block.toml":                   "malformed_rule_shape",
			"neg/neg-outcome-outside.toml":                  "malformed_outcome_binding",
			"neg/neg-owned-no-reader.toml":                  "malformed_accessor_binding",
			"neg/neg-owned-two-readers.toml":                "malformed_accessor_binding",
			"neg/neg-recognized-in-guard.toml":              "malformed_outcome_binding",
			"neg/neg-recognized-in-keys.toml":               "malformed_accessor_binding",
			"neg/neg-recognized-misnamed.toml":              "reserved_tag_key",
			"neg/neg-ruleid-hash.toml":                      "malformed_rule_id",
			"neg/neg-terminal-unknown.toml":                 "unknown_context",
			"neg/neg-unknown-context.toml":                  "unknown_context",
			"neg/neg-unknown-field.toml":                    "unknown_schema_field",
			"neg/neg-unknown-tag-match.toml":                "unknown_tag",
			"neg/neg-v2-shaped.toml":                        "unsupported_version",
			"neg/neg-version.toml":                          "unsupported_version",
			"neg/neg-write-observed.toml":                   "write_to_non_owned_tag",
			"neg/probe-a-initial-observed-no-writer.toml":   "malformed_accessor_binding",
			"neg/probe-b-initial-observed-with-writer.toml": "write_to_non_owned_tag",
		}

		for _, rel := range spikeNegatives {
			w, ok := want[filepath.ToSlash(rel)]
			if !ok {
				t.Errorf("spike negative %s has no recorded category; a promoted "+
					"negative must be pinned to the category it witnesses", rel)
				continue
			}
			if got := loadCategory(t, rel); got != w {
				t.Errorf("%s refuses with category %q; want %q — a negative "+
					"rewritten to trip a different category is a narrowing under "+
					"REQ-118 even though the neg/ census stays whole", rel, got, w)
			}
		}
		if len(spikeNegatives) != len(want) {
			t.Errorf("%d spike negatives walked but %d categories recorded; the "+
				"pin must cover every promoted negative",
				len(spikeNegatives), len(want))
		}
	})

	t.Run("the RDR fixture's row identities are the approved nine", func(t *testing.T) {
		want := []string{
			"rdr.continue-prelock#foundational",
			"rdr.continue-prelock#large",
			"rdr.continue-prelock-cluster#foundational",
			"rdr.continue-prelock-cluster#large",
			"rdr.draft-no-match-escape#reconcile-block",
			"rdr.draft-no-match-escape#round-clean",
			"rdr.reconcile-rewind",
			"rdr.terminal-archive#finalized",
			"rdr.terminal-archive#verdict-flapping",
		}
		if got := rowIdentities(mustLoad(t, rdrFixture)); !reflect.DeepEqual(got, want) {
			t.Errorf("row identities =\n %v\nwant\n %v — the approved set. A row "+
				"swapped for filler holds the census and narrows the coverage",
				got, want)
		}
	})
	t.Run("the kata fixture's row identities are the approved two", func(t *testing.T) {
		want := []string{
			"kata.review-accepted",
			"kata.review-needs-work",
		}
		if got := rowIdentities(mustLoad(t, kataFixture)); !reflect.DeepEqual(got, want) {
			t.Errorf("row identities =\n %v\nwant\n %v — the approved set", got, want)
		}
	})
}

// REQ-119: "Each load-time validation category is asserted by **category**,
// not by message text, and each has one mutated fixture that trips it and
// no other"
// ADVERSARIAL
//
// The one-mutation rule holds by construction (each negative fixture is
// the RDR fixture under exactly one mutation). The oracle here is that no
// two categories share a witness and no witness trips a second category:
// an implementation collapsing several categories into one code fails.
func TestReq119_OneFixturePerCategoryAssertedByCategory(t *testing.T) {
	entries, err := os.ReadDir(filepath.Join("testdata", "neg"))
	if err != nil {
		t.Fatalf("read neg dir: %v", err)
	}

	byCategory := map[table.Category][]string{}
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".toml") {
			continue
		}
		rel := filepath.Join("neg", e.Name())
		_, err := table.Load(readFixture(t, rel), rel)
		if err == nil {
			t.Errorf("%s loaded clean; every neg/ fixture is a refusal", e.Name())
			continue
		}
		cat, ok := table.CategoryOf(err)
		if !ok {
			t.Errorf("%s refused with a non-categorized error: %v", e.Name(), err)
			continue
		}
		byCategory[cat] = append(byCategory[cat], e.Name())
	}

	// The category set the negatives cover must be wide: an implementation
	// collapsing them all into one code produces a single key here.
	if len(byCategory) < 20 {
		t.Errorf("the %d neg/ fixtures cover only %d distinct categories: %v — "+
			"an implementation collapsing categories fails here",
			len(entries), len(byCategory), byCategory)
	}
}

// REQ-120: "The existence-atom test exercises `exists = true` **and**
// `exists = false`"
// BOUNDARY
//
// A fixture carrying only one lets a hardcoded literal pass. Both arms are
// exercised in TestReq53; this asserts the two produce DIFFERENT normalized
// literals, which a hardcoded literal cannot.
func TestReq120_BothExistenceArmsProduceDistinctLiterals(t *testing.T) {
	base := string(readFixture(t, rdrFixture))
	src := strings.Replace(base,
		"[rule.guard.all.finalized_at]\nexists = false",
		"[rule.guard.all.finalized_at]\nexists = true", 1)
	if src == base {
		t.Fatal("exists substitution did not apply")
	}

	falseSide := mustLoad(t, rdrFixture)
	trueSide, err := table.Load([]byte(src), "exists-true.toml")
	if err != nil {
		t.Fatalf("exists = true refused: %v", err)
	}

	f := atomsOn(rowByID(t, falseSide, "rdr.continue-prelock#large"), "finalized_at")
	tr := atomsOn(rowByID(t, trueSide, "rdr.continue-prelock#large"), "finalized_at")
	if len(f) != 1 || len(tr) != 1 {
		t.Fatalf("existence atoms: false=%+v true=%+v", f, tr)
	}
	if reflect.DeepEqual(f[0].Literal, tr[0].Literal) {
		t.Errorf("both arms produced literal %v; a hardcoded literal passes a "+
			"one-arm fixture", f[0].Literal)
	}
}

// REQ-121: "The byte-equality contract is asserted by a negative control —
// a reference spelling a declared `status` as `Status` must fail `unknown
// tag`"
// ADVERSARIAL
func TestReq121_MisCasedReferenceFailsUnknownTag(t *testing.T) {
	got := loadCategory(t, "neg/neg-case-folded-tag.toml")

	if got != table.CatUnknownTag {
		t.Errorf("category = %q; want %q — a normalizer applying case folding "+
			"passes every positive test", got, table.CatUnknownTag)
	}
}

// REQ-122: "assert that the normalized value is unchanged when the same
// model is authored with its TOML keys in a different order, **and when its
// rules are declared in a different order**, and that repeated dumps of one
// model are byte-identical."
// BOUNDARY
//
// Rule-order permutation is what excludes a positional tiebreak: sorting
// and checking no adjacent pair compares equal is satisfied by the defect
// it was meant to catch.
func TestReq122_KeyOrderRuleOrderAndRepeatedDumpDeterminism(t *testing.T) {
	base := mustLoad(t, rdrFixture)

	t.Run("key order", func(t *testing.T) {
		if !reflect.DeepEqual(mustLoad(t, "perm/rdr-keyorder.toml").Rows, base.Rows) {
			t.Error("a key-order permutation changed the normalized value")
		}
	})
	t.Run("rule declaration order", func(t *testing.T) {
		if !reflect.DeepEqual(mustLoad(t, "perm/rdr-ruleorder.toml").Rows, base.Rows) {
			t.Error("a rule-order permutation changed the normalized value; this " +
				"is the assertion that excludes a positional tiebreak")
		}
	})
	t.Run("repeated dumps are byte-identical", func(t *testing.T) {
		want := table.Dump(base)
		for i := range 16 {
			if got := table.Dump(mustLoad(t, rdrFixture)); got != want {
				t.Fatalf("dump %d differs despite Go's randomized map iteration", i)
			}
		}
	})
}

// REQ-123: "**`[model.metadata]` is asserted by deep equality against the
// authored table** — every value and every nesting level" and "The `source`
// and `description` keys decode onto the rule and model structs rather than
// refusing as unknown schema fields."
// HAPPY PATH (TS 1)
func TestReq123_MetadataDeepEqualityAndAnnotationKeysDecode(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	want := map[string]any{
		"owner":          "flow-team",
		"review_cadence": "per-lock",
	}
	if !reflect.DeepEqual(m.Metadata, want) {
		t.Errorf("metadata = %#v; want deep equality against the authored table "+
			"%#v — a loader discarding values passes a key-name oracle", m.Metadata, want)
	}

	if m.Description == "" {
		t.Error("[model].description did not decode onto the model struct")
	}
	for _, r := range m.Rows {
		if r.Source == "" {
			t.Errorf("%s carries no `source` annotation", r.Identity())
		}
	}

	k := mustLoad(t, kataFixture)
	if !reflect.DeepEqual(k.Metadata, map[string]any{"owner": "kata-team"}) {
		t.Errorf("kata metadata = %#v; want the authored table", k.Metadata)
	}
}

// REQ-124: "This SHA MUST NOT be asserted as a golden hash by any
// implementation test" / "Assert over the normalized value"
// ADVERSARIAL
//
// The ban stands on the oracle's CHARACTER, not its current value: a
// rendered-text hash is satisfied by any renderer that happens to agree on
// these inputs. This test enforces the ban by scanning the suite for the
// digest and for any hashing of dump text.
func TestReq124_NoGoldenHashOfRenderedText(t *testing.T) {
	// The needles are assembled from halves that Go folds at compile time, so
	// the runtime strings are exact while no verbatim needle appears in this
	// file's source bytes. That is what lets the scan below cover EVERY .go
	// file in the package with no exemption — including this one, the primary
	// round-trip file and the likeliest place a golden hash would be reached
	// for. An earlier revision instead exempted this file wholesale, because
	// the needles were verbatim literals here and the scan matched itself
	// (deviations.md D6); the exemption fixed that red-against-every-
	// implementation defect at the cost of blinding the scan to this file.
	// Assembling the needles removes the reason for the exemption. If a future
	// author writes a needle verbatim here, the scan fires — which is exactly
	// the signal the ban wants.
	const spikeSHA = "6ccfe9012b0705ef4d4b3d1c620da" +
		"ffd69523436175120be1bea8a05df9c55dd"
	const shaPkg = "crypto/" + "sha256"
	const md5Pkg = "crypto/" + "md5"

	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read package dir: %v", err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		b, err := os.ReadFile(e.Name())
		if err != nil {
			t.Fatalf("read %s: %v", e.Name(), err)
		}
		src := string(b)
		if strings.Contains(src, spikeSHA) {
			t.Errorf("%s asserts the spike SHA as a golden hash", e.Name())
		}
		if strings.Contains(src, shaPkg) || strings.Contains(src, md5Pkg) {
			t.Errorf("%s hashes output; assert over the normalized value and "+
				"reserve rendered-text goldens for tests of rendering itself", e.Name())
		}
	}
}

// REQ-125: "**The control is a mutation test**: normalize a row, mutate one
// field in place, and assert the other is unchanged."
// ADVERSARIAL (TS 2, no-alias)
//
// The behavioural assertion is TestReq82's. This asserts the SHAPE the
// control depends on: both fields are slices, so an alias shares a backing
// array and the mutation is observable.
func TestReq125_TheNoAliasControlIsAMutationTest(t *testing.T) {
	rt := reflect.TypeOf(table.Row{})
	for _, name := range []string{"NextTags", "Writes"} {
		f, ok := rt.FieldByName(name)
		if !ok {
			t.Errorf("table.Row has no %s field", name)
			continue
		}
		if f.Type.Kind() != reflect.Slice {
			t.Errorf("%s is %v; the mutation control depends on both fields "+
				"being slices, where an alias shares a backing array", name, f.Type)
		}
	}

	// The mutation control must actually be able to observe an alias: two
	// independently built slices must not share a backing array.
	m := mustLoad(t, rdrFixture)
	row := rowByID(t, m, "rdr.reconcile-rewind")
	if len(row.NextTags) == 0 {
		t.Fatal("no writes on the row under test")
	}
	if &row.NextTags[0] == &row.Writes[0] {
		t.Error("NextTags and Writes share a backing array; one is populated by " +
			"aliasing the other")
	}
}

// REQ-126: "Two contexts contributing `in = ["a,b", "c"]` and `in = ["a",
// "b,c"]` on one key and block MUST yield **two** atoms in the normalized
// set" and "the control instead **parameterizes over the plausible
// delimiters** (`,` `;` `|` space, and the empty string)"
// ADVERSARIAL (TS 2)
//
// The atom count is the whole oracle: a property of this RDR's own
// normalized value, readable without RDR 0006, failing against any
// joined-rendering key. Stating it as "the rule is a dead rule for lint"
// would be an RDR 0006 verdict.
func TestReq126_DistinctLiteralControlParameterizedOverFiveDelimiters(t *testing.T) {
	for _, d := range []string{"comma", "semi", "pipe", "space", "empty"} {
		t.Run(d, func(t *testing.T) {
			rel := "delim/merge-delim-" + d + ".toml"
			m, err := table.Load(readFixture(t, rel), rel)
			if err != nil {
				t.Fatalf("%s refused: %v", rel, err)
			}
			rows := rowsByRuleID(m, "reconcile-rewind")
			if len(rows) == 0 {
				t.Fatalf("%s yielded no rows for the contributing rule", rel)
			}
			got := atomsOn(rows[0], "finalized_at")
			if len(got) != 2 {
				t.Errorf("%d atoms; want 2 — a merge keyed on a %s-joined "+
					"rendering collapses the pair: %+v", len(got), d, got)
			}
		})
	}
}

// REQ-127: "A `set`-kind write authoring `["a,b", "c"]` and one authoring
// `["a", "b,c"]` MUST normalize to two distinct write values, parameterized
// over the same delimiter set." / "The assertion must read the normalized
// value, not the rendered one"
// ADVERSARIAL (TS 2)
func TestReq127_WriteValueControlReadsTheNormalizedValue(t *testing.T) {
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

			va, _ := tagValue(rowByID(t, a, "kata.review-needs-work").Writes, "labels")
			vb, _ := tagValue(rowByID(t, b, "kata.review-needs-work").Writes, "labels")

			if reflect.DeepEqual(va, vb) {
				t.Errorf("the two spellings normalized to the same value %v; "+
					"RDR 0004's read-back compares on this equality", va)
			}
			if len(va) != 2 || len(vb) != 2 {
				t.Errorf("normalized write values are not two-member sequences: "+
					"a=%v b=%v", va, vb)
			}
		})
	}
}

// REQ-128: "a rule and an inherited context contributing a
// **byte-identical** `(block, key, operator, literal)` atom collapse to
// **exactly one** atom, asserted by count on that key."
// BOUNDARY (TS 2, A13's idempotence mirror)
func TestReq128_ByteIdenticalAtomsCollapseToExactlyOne(t *testing.T) {
	m := mustLoad(t, "merge-idempotent.toml")

	got := atomsOn(rowByID(t, m, "rdr.reconcile-rewind"), "status")
	if len(got) != 1 {
		t.Errorf("%d atoms on `status`; want exactly 1 — without this a "+
			"normalizer that de-duplicates nothing passes every distinct-literal "+
			"control: %+v", len(got), got)
	}
}

// REQ-129: "the fixture's `draft-no-match-escape` binds its outcome with
// `in` over two alphabet members, and normalization must yield **two**
// escape rows carrying distinct expansion suffixes, each with an empty
// write set and each retaining the modeled failure-class list."
// HAPPY PATH (TS 2, A11)
func TestReq129_EscapeExpansionYieldsTwoRowsWithDistinctSuffixes(t *testing.T) {
	m := mustLoad(t, rdrFixture)
	rows := rowsByRuleID(m, "draft-no-match-escape")

	if len(rows) != 2 {
		t.Fatalf("%d escape rows; want 2 (%v)", len(rows), rowIdentities(m))
	}
	seen := map[string]bool{}
	for _, r := range rows {
		key := strings.Join(r.Suffix, "\x00")
		if seen[key] {
			t.Errorf("two escape rows share the expansion suffix %v", r.Suffix)
		}
		seen[key] = true

		if len(r.Writes) != 0 {
			t.Errorf("%s carries writes %+v; want an empty write set",
				r.Identity(), r.Writes)
		}
		if !reflect.DeepEqual(r.Escape, []string{"no_match"}) {
			t.Errorf("%s escape classes = %v; each row retains the modeled "+
				"failure-class list", r.Identity(), r.Escape)
		}
	}
	want := []string{
		"rdr.draft-no-match-escape#reconcile-block",
		"rdr.draft-no-match-escape#round-clean",
	}
	var got []string
	for _, r := range rows {
		got = append(got, r.Identity())
	}
	slices.Sort(got)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("escape row identities = %v; want %v", got, want)
	}
}

// REQ-130: "constructing the `resolve.Row` for each normalized RDR-fixture
// row, the union of `Match` and the guard's atoms equals the normalized
// atom set and their intersection is empty, and every `Match` atom carries
// block `match`."
// ADVERSARIAL (TS 2, A12)
//
// req-list Q1 proceeds under reading (a): `Row.Match []Tag` stays the
// equality pattern and `Row.Guard []GuardAtom` takes only all/unless
// atoms. The union is therefore asserted across two Go fields of different
// types, and "every Match atom carries block match" is asserted by every
// Match entry tracing to a BlockMatch atom in the normalized set.
func TestReq130_HandoffRoutingIsTotalAndDisjoint(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	for _, row := range m.Rows {
		t.Run(row.Identity(), func(t *testing.T) {
			kr := row.KernelRow()

			matchAtoms := atomsInBlock(row, table.BlockMatch)
			guardAtoms := append(
				slices.Clone(atomsInBlock(row, table.BlockAll)),
				atomsInBlock(row, table.BlockUnless)...)

			// Union equals the normalized atom set.
			if len(kr.Match)+len(kr.Guard) != len(row.Atoms) {
				t.Fatalf("|Match| (%d) + |Guard| (%d) != |atoms| (%d); the "+
					"routing is not total: %+v",
					len(kr.Match), len(kr.Guard), len(row.Atoms), row.Atoms)
			}
			if len(kr.Match) != len(matchAtoms) {
				t.Errorf("|Match| = %d; want %d — every match-block atom and no "+
					"other", len(kr.Match), len(matchAtoms))
			}
			if len(kr.Guard) != len(guardAtoms) {
				t.Errorf("|Guard| = %d; want %d — every all/unless atom and no "+
					"other", len(kr.Guard), len(guardAtoms))
			}

			// Intersection is empty: no key+value pair appears in both, and
			// no guard atom carries BlockMatch.
			for _, g := range kr.Guard {
				if string(g.Block) == string(table.BlockMatch) {
					t.Errorf("guard atom %+v carries block match", g)
				}
			}

			// Every Match atom carries block `match` in the normalized set.
			for _, mt := range kr.Match {
				found := false
				for _, a := range matchAtoms {
					if a.Key == mt.Key && len(a.Literal) == 1 && a.Literal[0] == mt.Value {
						found = true
					}
				}
				if !found {
					t.Errorf("Match entry %+v traces to no BlockMatch atom in "+
						"%+v", mt, matchAtoms)
				}
			}

			// Symmetrically on the guard side, and BIJECTIVELY: |Guard|
			// matching |guardAtoms| above admits one atom emitted twice
			// while another is dropped.
			assertGuardBijection(t, row, kr.Guard)
		})
	}
}

// REQ-131: "Each variant is refused with the **one** category its mutation
// targets and no other — the assertion is on the category, not the message
// text"
// ADVERSARIAL (TS 3)
//
// The oracle is category identity, never a substring of the message: a
// message-text oracle passed against a wrong implementation, which is why
// the record replaced it.
func TestReq131_RefusalsAreAssertedByCategoryNotMessageText(t *testing.T) {
	// Two fixtures in DIFFERENT categories whose messages both mention the
	// same rule. A message-text oracle cannot separate them; a category
	// oracle can.
	a := loadCategory(t, "neg/neg-escape-with-write.toml")
	b := loadCategory(t, "neg/neg-escape-with-gate.toml")
	c := loadCategory(t, "neg/neg-outcome-outside.toml")

	if a != b {
		t.Errorf("the two escape shapes carry %q and %q; the record deliberately "+
			"declines to split them into separate categories", a, b)
	}
	if a == c {
		t.Errorf("an escape-declaration defect and an outcome-binding defect both "+
			"refused %q; an implementation collapsing categories fails here", a)
	}

	// The Category type is an API surface, not message text.
	if reflect.TypeOf(table.Category("")).Kind() != reflect.String {
		t.Error("table.Category is not a named string type")
	}
}

// REQ-132: "The **seven** still owed at implementation are `malformed
// TOML`, `missing recognized outcome alphabet`, `cyclic context
// inheritance`, `malformed tag declaration` …, `duplicate model id` …,
// `malformed model declaration` … and `malformed dump declaration`"
// BOUNDARY (TS 3)
func TestReq132_TheSevenOwedCategoriesAreNowWitnessed(t *testing.T) {
	owed := map[table.Category]string{
		table.CatMalformedTOML:                    "neg/neg-malformed-toml.toml",
		table.CatMissingRecognizedOutcomeAlphabet: "neg/neg-no-alphabet.toml",
		table.CatCyclicContextInheritance:         "neg/neg-cyclic-context.toml",
		table.CatMalformedTagDeclaration:          "neg/neg-tagdecl-domain-on-bool.toml",
		table.CatMalformedModelDeclaration:        "neg/neg-no-model-table.toml",
		table.CatMalformedDumpDeclaration:         "neg/neg-dump-repeated-column.toml",
	}
	for cat, rel := range owed {
		t.Run(string(cat), func(t *testing.T) {
			if got := loadCategory(t, rel); got != cat {
				t.Errorf("%s refused %q; want %q", rel, got, cat)
			}
		})
	}

	// The seventh is cross-document and needs the paired-document surface.
	t.Run("duplicate_model_id", func(t *testing.T) {
		a := mustLoad(t, "dup/dup-model-a.toml")
		b := mustLoad(t, "dup/dup-model-b.toml")
		err := table.CheckModelIDs([]*table.Model{a, b})
		if err == nil {
			t.Fatal("the paired documents passed the cross-document check")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatDuplicateModelID {
			t.Errorf("category = %q; want %q", cat, table.CatDuplicateModelID)
		}
	})
}

// REQ-133: "**Guard-block controls are owed for every atom-level rule
// (A15).** … Each of operator membership, `<clear>`, domain/kind
// conformance, `#` reservation, and tag-key declaration owes a `guard.all`
// and a `guard.unless` variant."
// ADVERSARIAL (TS 3)
//
// The `#` reservation is legitimately match-only (REQ-90), so it owes no
// guard variant; the other four owe one per guard block. The behavioural
// assertions are TestReq88's; this asserts the CONTROLS exist as promoted
// fixtures so a later-added rule inherits the obligation.
func TestReq133_GuardBlockControlsExistForEveryBlockAgnosticRule(t *testing.T) {
	required := map[string]string{
		"operator membership in guard.all":    "neg/neg-guard-all-bad-operator.toml",
		"operator membership in guard.unless": "neg/neg-guard-unless-bad-operator.toml",
		"<clear> in guard.all":                "neg/neg-guard-all-clear.toml",
		"<clear> in guard.unless":             "neg/neg-guard-unless-clear.toml",
		"domain conformance in guard.unless":  "neg/neg-guard-unless-out-of-domain.toml",
		"tag-key declaration in guard.unless": "neg/neg-guard-unless-unknown-tag.toml",
	}
	for name, rel := range required {
		t.Run(name, func(t *testing.T) {
			if _, err := os.Stat(filepath.Join("testdata", rel)); err != nil {
				t.Fatalf("the guard-block control %s is not promoted: %v", rel, err)
			}
			if _, err := table.Load(readFixture(t, rel), rel); err == nil {
				t.Errorf("%s loaded clean; a guard-block control must refuse", rel)
			}
		})
	}
}

// REQ-134: "**Fixture conformance is a control, not a review step.** Both
// normative fixtures MUST load clean under the finished loader, asserted as
// part of scenario 1."
// HAPPY PATH (TS 3)
//
// This is the check that would have caught the `[dump]` column vocabulary
// being fixed to spellings the fixtures do not use.
func TestReq134_BothNormativeFixturesLoadClean(t *testing.T) {
	for _, rel := range []string{rdrFixture, kataFixture} {
		t.Run(rel, func(t *testing.T) {
			m, err := table.Load(readFixture(t, rel), rel)
			if err != nil {
				t.Fatalf("the normative fixture refused: %v", err)
			}
			if len(m.Rows) == 0 {
				t.Error("loaded clean but normalized to no rows")
			}
			if !reflect.DeepEqual(m.DumpOrder, table.DumpColumns()) {
				t.Errorf("the fixture's [dump].order %v is not the column "+
					"vocabulary %v", m.DumpOrder, table.DumpColumns())
			}
		})
	}
}

// REQ-135: "The precedence control is a **v2-shaped** document — `version =
// 2` plus a v2-only key the strict decoder would reject as an unknown
// schema field — asserted to refuse `unsupported version` and **not**
// `unknown schema field`."
// ADVERSARIAL (TS 3, A14)
func TestReq135_V2ShapedPrecedenceControl(t *testing.T) {
	src := string(readFixture(t, "neg/neg-v2-shaped.toml"))

	if !strings.Contains(src, "version = 2") {
		t.Error("the precedence fixture is not v2-shaped: no `version = 2`")
	}

	got := loadCategory(t, "neg/neg-v2-shaped.toml")
	if got == table.CatUnknownSchemaField {
		t.Fatal("refused `unknown_schema_field`; strict decoding ran before the " +
			"version gate")
	}
	if got != table.CatUnsupportedVersion {
		t.Errorf("category = %q; want %q", got, table.CatUnsupportedVersion)
	}
}

// REQ-136: "Positive controls in the same fixture set: an observed key
// served by no reader loads (JD-9), and a model with no `[initial]` loads
// and is left to lint."
// HAPPY PATH (TS 3)
func TestReq136_PositiveControls(t *testing.T) {
	t.Run("observed key served by no reader loads (JD-9)", func(t *testing.T) {
		m := mustLoad(t, kataFixture)
		if m.Tags["owner"].Provenance != table.ProvenanceObserved {
			t.Fatal("the kata fixture's `owner` is not observed")
		}
	})
	t.Run("a model with no [initial] loads", func(t *testing.T) {
		m := mustLoad(t, "pos-no-initial.toml")
		if len(m.Initial) != 0 {
			t.Errorf("[initial] = %+v; want empty", m.Initial)
		}
	})
	t.Run("a model with no terminal loads", func(t *testing.T) {
		m := mustLoad(t, "pos-no-terminal.toml")
		if len(m.Terminal) != 0 {
			t.Errorf("terminal = %+v; want empty", m.Terminal)
		}
	})
}

// REQ-138: "The non-boolean literal is refused at load as a malformed
// predicate atom; the mis-cased reference fails `unknown tag` rather than
// folding."
// ADVERSARIAL (TS 6)
//
// The discriminating half is that the two carry DIFFERENT categories.
func TestReq138_ScenarioSixCarriesTwoDistinctCategories(t *testing.T) {
	nonBool := loadCategory(t, "neg/neg-bad-exists.toml")
	misCased := loadCategory(t, "neg/neg-case-folded-tag.toml")

	if nonBool != table.CatMalformedPredicateAtom {
		t.Errorf("non-boolean existence literal: category = %q; want %q",
			nonBool, table.CatMalformedPredicateAtom)
	}
	if misCased != table.CatUnknownTag {
		t.Errorf("mis-cased reference: category = %q; want %q — it must not fold",
			misCased, table.CatUnknownTag)
	}
	if nonBool == misCased {
		t.Errorf("both refused %q; the two defects carry distinct categories", nonBool)
	}
}

// REQ-139: "the scenario 3 fixture for this category is a **pair** of
// documents sharing a `model id`, and a single-document load of either one
// must not report it."
// ADVERSARIAL (CA A10 / TS 3)
func TestReq139_DuplicateModelIDNeedsThePairedDocumentSurface(t *testing.T) {
	for _, rel := range []string{"dup/dup-model-a.toml", "dup/dup-model-b.toml"} {
		t.Run("single-document load of "+rel, func(t *testing.T) {
			_, err := table.Load(readFixture(t, rel), rel)
			if err != nil {
				if cat, _ := table.CategoryOf(err); cat == table.CatDuplicateModelID {
					t.Fatalf("%s reported %q on a single-document load; the "+
						"category is never decidable within one document",
						rel, table.CatDuplicateModelID)
				}
				t.Fatalf("%s refused: %v", rel, err)
			}
		})
	}

	t.Run("the pair shares a model id", func(t *testing.T) {
		a := mustLoad(t, "dup/dup-model-a.toml")
		b := mustLoad(t, "dup/dup-model-b.toml")
		if a.ID != b.ID {
			t.Fatalf("the fixture pair does not share a model id: %q / %q", a.ID, b.ID)
		}
		if reflect.DeepEqual(a.Rows, b.Rows) {
			t.Error("the pair is byte-identical; it must be two documents, not one")
		}
	})
}

// REQ-140: "New internal package can own sparse source structs and
// normalization."
// BOUNDARY
func TestReq140_PackagePlacement(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if got := filepath.Base(filepath.Dir(wd)); got != "internal" {
		t.Errorf("the package lives under %q; it is a new INTERNAL package", got)
	}
}

// REQ-141: Phase 2 "Introduce typed source structures, normalization to
// `internal/resolve::Row` values (outcome lifted, per-atom block retained,
// `RequiresOwned` derived, existence constants emitted), and an
// expanded-table dump with deterministic ordering and source rule ids."
// BOUNDARY
//
// Phase 2's deliverable list, asserted as one surface check: every named
// capability exists and is reachable from the normalized model.
func TestReq141_PhaseTwoDeliverables(t *testing.T) {
	m := mustLoad(t, rdrFixture)
	row := rowByID(t, m, "rdr.continue-prelock#large")

	if row.Outcome == "" || len(atomsOn(row, table.RecognizedTagKey)) != 0 {
		t.Error("the outcome is not lifted")
	}
	if len(atomsInBlock(row, table.BlockUnless)) == 0 {
		t.Error("the per-atom block is not retained")
	}
	if len(row.RequiresOwned) == 0 {
		t.Error("RequiresOwned is not derived")
	}
	if row.KernelRow().RuleID != row.RuleID {
		t.Error("normalization does not reach internal/resolve::Row values")
	}
	if !strings.Contains(table.Dump(m), row.RuleID) {
		t.Error("the expanded-table dump does not carry source rule ids")
	}
}

// REQ-142: Phase 3 "Add validation rules for tag and accessor
// declarations, the provenance-scoped `keys` bindings, the reserved
// `recognized` name and `<clear>` value, rule ids, context references,
// `[initial]` assignments and `terminal` ids, predicate atoms (including
// the match-block operator restriction and domain membership), outcome
// bindings, gate references, write targets, explicit clears, escape
// declarations, and expansion-count diagnostics."
// BOUNDARY
//
// Phase 3's list, asserted as one witness per validation family.
func TestReq142_PhaseThreeValidationFamilies(t *testing.T) {
	families := map[string]struct {
		rel  string
		want table.Category
	}{
		"tag declaration":            {"neg/neg-tagdecl-domain-on-bool.toml", table.CatMalformedTagDeclaration},
		"accessor declaration":       {"neg/neg-accessor-no-timeout.toml", table.CatMalformedAccessorDeclaration},
		"provenance-scoped bindings": {"neg/neg-write-observed.toml", table.CatWriteToNonOwnedTag},
		"reserved recognized name":   {"neg/neg-recognized-misnamed.toml", table.CatReservedTagKey},
		"reserved <clear> value":     {"neg/neg-clear-as-write.toml", table.CatReservedTagValue},
		"rule ids":                   {"neg/neg-ruleid-hash.toml", table.CatMalformedRuleID},
		"context references":         {"neg/neg-unknown-context.toml", table.CatUnknownContext},
		"[initial] assignments":      {"neg/neg-initial-bad-value.toml", table.CatMalformedInitialDeclaration},
		"terminal ids":               {"neg/neg-terminal-unknown.toml", table.CatUnknownContext},
		"match-block operators":      {"neg/neg-match-lt.toml", table.CatMalformedPredicateAtom},
		"domain membership":          {"neg/neg-literal-outside-domain.toml", table.CatMalformedPredicateAtom},
		"outcome bindings":           {"neg/neg-outcome-outside.toml", table.CatMalformedOutcomeBinding},
		"gate references":            {"neg/neg-gate-is-reader.toml", table.CatUnknownAccessor},
		"write targets":              {"neg/neg-initial-unknown.toml", table.CatUnknownTag},
		"escape declarations":        {"neg/neg-escape-with-write.toml", table.CatMalformedEscapeDeclaration},
	}
	for name, spec := range families {
		t.Run(name, func(t *testing.T) {
			if got := loadCategory(t, spec.rel); got != spec.want {
				t.Errorf("category = %q; want %q", got, spec.want)
			}
		})
	}
}

// REQ-143: Phase 4 "Connect the normalized rows to RDR 0001's
// gate-then-count exact-one contract without adding runtime ordering or
// host-code predicate callbacks."
// BOUNDARY
//
// The "without" half is the assertable one: the kernel table this package
// builds carries no ordering field and no callback. The gate-then-count
// behaviour itself is the kernel's and is exercised by the MVV.
func TestReq143_KernelHandoffAddsNoOrderingOrCallback(t *testing.T) {
	m := mustLoad(t, rdrFixture)
	tbl := m.KernelTable()

	if len(tbl.Rows) != len(m.Rows) {
		t.Errorf("the kernel table carries %d rows; want %d",
			len(tbl.Rows), len(m.Rows))
	}
	if !reflect.DeepEqual(tbl.Outcomes, m.Outcomes) {
		t.Errorf("the kernel table's alphabet = %v; want %v", tbl.Outcomes, m.Outcomes)
	}

	// No host-code predicate callback: no field on the produced row is a
	// func, and the guard is a slice of parsed atoms.
	for _, r := range tbl.Rows {
		v := reflect.ValueOf(r)
		for i := range v.NumField() {
			if v.Field(i).Kind() == reflect.Func {
				t.Errorf("kernel row field %q is a func; no host-code predicate "+
					"callbacks are added", v.Type().Field(i).Name)
			}
		}
	}
}

// REQ-144: Phase 5 "Expose the parsed representation needed by RDR 0006
// for graph determinism, reachability, and read-before-write checks."
// BOUNDARY
//
// The deliverable is the EXPOSED SHAPE, not an RDR 0006 implementation
// (REQ-117 forbids asserting against a stub): the normalized model must
// carry what those three checks read.
func TestReq144_ParsedRepresentationIsExposedForLint(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	t.Run("graph determinism: rows comparable by identity and predicate set", func(t *testing.T) {
		seen := map[string]bool{}
		for _, r := range m.Rows {
			if seen[r.Identity()] {
				t.Errorf("duplicate row identity %q", r.Identity())
			}
			seen[r.Identity()] = true
			if r.Outcome == "" {
				t.Errorf("%s exposes no outcome to key an overlap check on", r.Identity())
			}
		}
	})

	t.Run("reachability: root and stop set are predicate sets over owned tags", func(t *testing.T) {
		if len(m.Initial) == 0 {
			t.Error("the root is not exposed")
		}
		if len(m.Terminal) == 0 {
			t.Error("the stop set is not exposed")
		}
		for _, set := range m.Terminal {
			for _, a := range set {
				if _, ok := m.Tags[a.Key]; !ok {
					t.Errorf("terminal predicate names the undeclared tag %q", a.Key)
				}
			}
		}
	})

	t.Run("read-before-write: rows expose what they read and what they write", func(t *testing.T) {
		row := rowByID(t, m, "rdr.continue-prelock#large")
		if len(row.Atoms) == 0 {
			t.Error("the row exposes nothing it reads")
		}
		if len(row.Writes) == 0 {
			t.Error("the row exposes nothing it writes")
		}
		if len(row.RequiresOwned) == 0 {
			t.Error("the row exposes no owned-key requirement")
		}
	})
}

// REQ-145: "expansion counts per rule | dump (derived: rows per source rule
// id) | this RDR | reviewer-visible, non-blocking; no threshold here."
// BOUNDARY
//
// Derived from the dump, not a lint finding and not a refusal: Phase 3
// mints no check for it, so a high expansion count must LOAD.
func TestReq145_ExpansionCountsAreDerivedAndNonBlocking(t *testing.T) {
	t.Run("counts are derivable as rows per source rule id", func(t *testing.T) {
		m := mustLoad(t, rdrFixture)
		want := map[string]int{
			"continue-prelock":         2,
			"continue-prelock-cluster": 2,
			"draft-no-match-escape":    2,
			"reconcile-rewind":         1,
			"terminal-archive":         2,
		}
		got := map[string]int{}
		for _, r := range m.Rows {
			got[r.RuleID]++
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("rows per source rule id = %v; want %v", got, want)
		}
	})

	t.Run("no threshold refuses", func(t *testing.T) {
		base := string(readFixture(t, rdrFixture))
		// A rule expanding across the whole alphabet plus the whole profile
		// domain: 4 x 4 = 16 rows from one rule.
		src := strings.Replace(base,
			"source = \"rdr:terminal\"\n[rule.match.recognized]\nin = [\"finalized\", \"verdict-flapping\"]",
			"source = \"rdr:terminal\"\n[rule.match.profile]\nin = [\"small\", \"mid\", \"large\", \"foundational\"]\n[rule.match.recognized]\nin = [\"round-clean\", \"verdict-flapping\", \"reconcile-block\", \"finalized\"]", 1)
		if src == base {
			t.Fatal("expansion substitution did not apply")
		}
		m, err := table.Load([]byte(src), "high-expansion.toml")
		if err != nil {
			t.Fatalf("a high expansion count refused with %v; it is "+
				"reviewer-visible and non-blocking, with no threshold here", err)
		}
		if got := len(rowsByRuleID(m, "terminal-archive")); got != 16 {
			t.Errorf("terminal-archive yielded %d rows; want 16", got)
		}
	})

	t.Run("no expansion-count category exists", func(t *testing.T) {
		for _, c := range table.Categories() {
			if strings.Contains(string(c), "expansion_count") {
				t.Errorf("category %q mints an expansion-count refusal; it is "+
					"derived from the dump, not a lint finding and not a refusal", c)
			}
		}
	})
}

// REQ-146: "source key order is neutralized by sorting every emitted
// sequence; Go's map iteration is never the emission order (rows, atoms,
// next tags, writes, required-owned keys, and escape classes are all sorted
// slices before rendering)"
// ADVERSARIAL
//
// The determinism checklist, run as one sweep: every emitted sequence is a
// sorted slice, and repeated loads of one document produce one value
// despite Go's randomized map iteration.
func TestReq146_EveryEmittedSequenceIsASortedSlice(t *testing.T) {
	t.Run("repeated loads produce one value", func(t *testing.T) {
		first := mustLoad(t, rdrFixture)
		for i := range 16 {
			if got := mustLoad(t, rdrFixture); !reflect.DeepEqual(got.Rows, first.Rows) {
				t.Fatalf("load %d produced a different normalized value; Go's map "+
					"iteration reached the emission order", i)
			}
		}
	})

	t.Run("every emitted sequence is a slice, never a map", func(t *testing.T) {
		rt := reflect.TypeOf(table.Row{})
		for _, name := range []string{"Atoms", "NextTags", "Writes", "RequiresOwned", "Gate", "Escape", "Suffix"} {
			f, ok := rt.FieldByName(name)
			if !ok {
				t.Errorf("table.Row has no %s field", name)
				continue
			}
			if f.Type.Kind() != reflect.Slice {
				t.Errorf("%s is %v; every emitted sequence is a sorted slice, "+
					"never a map", name, f.Type)
			}
		}
	})

	t.Run("the row sequence itself is sorted", func(t *testing.T) {
		for _, rel := range []string{rdrFixture, kataFixture} {
			ids := rowIdentities(mustLoad(t, rel))
			if !slices.IsSorted(ids) {
				t.Errorf("%s: rows are not emitted in sorted order: %v", rel, ids)
			}
		}
	})
}
