package table_test

// RDR 0002 sections N–O: normalization output and ordering, the expanded
// table dump, and the load-category taxonomy.

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/resolve"
	"github.com/newcoinc/intrastate/internal/table"
)

// REQ-91: "The tool MUST normalize the sparse source into deterministic
// candidate rows for lint, resolver lookup, diagnostics, and table dumps.
// Each candidate row MUST retain its source rule id and source locator."
// HAPPY PATH
func TestReq91_EveryCandidateRowRetainsRuleIDAndLocator(t *testing.T) {
	for _, rel := range []string{rdrFixture, kataFixture} {
		t.Run(rel, func(t *testing.T) {
			m := mustLoad(t, rel)
			if len(m.Rows) == 0 {
				t.Fatal("normalized to no rows")
			}
			for _, r := range m.Rows {
				if r.RuleID == "" {
					t.Errorf("%s retains no source rule id", r.Identity())
				}
				if r.SourceLocator == "" {
					t.Errorf("%s retains no source locator", r.Identity())
				}
			}
		})
	}
}

// REQ-92: "The locator must identify at least the model id and rule id,
// and it is **derived, not authored**: the loader composes it from `(model
// id, rule id)` … An authored `[[rule]].source` is a provenance
// *annotation* carried alongside it, not the locator itself"
// ADVERSARIAL
//
// The discriminating case: editing `source` must not change the locator.
// A loader using `source` as the locator passes any test that only checks
// the locator is non-empty.
func TestReq92_LocatorIsDerivedFromModelAndRuleIDNeverFromSource(t *testing.T) {
	m := mustLoad(t, rdrFixture)
	row := rowByID(t, m, "rdr.continue-prelock#large")

	if !strings.Contains(row.SourceLocator, m.ID) {
		t.Errorf("locator %q does not identify the model id %q", row.SourceLocator, m.ID)
	}
	if !strings.Contains(row.SourceLocator, row.RuleID) {
		t.Errorf("locator %q does not identify the rule id %q", row.SourceLocator, row.RuleID)
	}

	base := string(readFixture(t, rdrFixture))
	src := strings.Replace(base, `source = "rdr:prelock"`, `source = "somewhere/else.toml:99"`, 1)
	if src == base {
		t.Fatal("source substitution did not apply")
	}
	edited, err := table.Load([]byte(src), rdrFixture)
	if err != nil {
		t.Fatalf("edited document refused: %v", err)
	}
	got := rowByID(t, edited, "rdr.continue-prelock#large")

	if got.SourceLocator != row.SourceLocator {
		t.Errorf("editing `source` changed the locator: %q -> %q — the locator "+
			"is derived from (model id, rule id) and never from `source`",
			row.SourceLocator, got.SourceLocator)
	}
	if got.Source != "somewhere/else.toml:99" {
		t.Errorf("the `source` annotation = %q; want the edited value", got.Source)
	}
}

// REQ-93: JDR 0001 §D13 (the landing this RDR owes; deviations.md D5):
// "`Tag.Value` stays `string`; a set crosses as its canonical JSON array —
// members sorted, duplicate-free, compact encoding — and 0002 declares
// it."
// BOUNDARY
//
// Two layers, not a contradiction: the NORMALIZED value holds a []string
// member sequence (REQ-38/56/58); the KERNEL SEAM encoding — what lands in
// resolve.Tag.Value, a string — is the canonical JSON array. RDR 0007
// already wrote its `contains` leg against §D13, so this build matches
// that form rather than redefining it.
func TestReq93_SetsCrossTheKernelSeamAsCanonicalJSONArrays(t *testing.T) {
	m := mustLoad(t, kataFixture)
	kr := rowByID(t, m, "kata.review-needs-work").KernelRow()

	var labels *resolve.Tag
	for i, w := range kr.Writes {
		if w.Key == "labels" {
			labels = &kr.Writes[i]
		}
	}
	if labels == nil {
		t.Fatalf("no `labels` write reached the kernel row: %+v", kr.Writes)
	}

	t.Run("compact JSON array, not a joined string", func(t *testing.T) {
		var got []string
		if err := json.Unmarshal([]byte(labels.Value), &got); err != nil {
			t.Fatalf("Tag.Value %q is not a JSON array: %v", labels.Value, err)
		}
		if !reflect.DeepEqual(got, []string{"needs work"}) {
			t.Errorf("decoded members = %v; want [\"needs work\"]", got)
		}
		if labels.Value != `["needs work"]` {
			t.Errorf("Tag.Value = %q; want the compact encoding %q",
				labels.Value, `["needs work"]`)
		}
	})

	t.Run("members sorted and duplicate-free", func(t *testing.T) {
		src := strings.Replace(string(readFixture(t, kataFixture)),
			`labels = ["needs work"]`, `labels = ["chore", "bug", "chore"]`, 1)
		k, err := table.Load([]byte(src), "set-encoding.toml")
		if err != nil {
			t.Fatalf("refused: %v", err)
		}
		enc := rowByID(t, k, "kata.review-needs-work").KernelRow()
		for _, w := range enc.Writes {
			if w.Key != "labels" {
				continue
			}
			if w.Value != `["bug","chore"]` {
				t.Errorf("Tag.Value = %q; want the sorted, duplicate-free, "+
					"compact array %q", w.Value, `["bug","chore"]`)
			}
		}
	})

	t.Run("read-back equality is byte equality over the encoding", func(t *testing.T) {
		// Two authored orderings of one set must cross the seam as the
		// same bytes, which is what makes RDR 0004's read-back a byte
		// comparison.
		build := func(members string) string {
			return strings.Replace(string(readFixture(t, kataFixture)),
				`labels = ["needs work"]`, `labels = `+members, 1)
		}
		a, err := table.Load([]byte(build(`["bug", "chore"]`)), "seam-a.toml")
		if err != nil {
			t.Fatalf("seam-a refused: %v", err)
		}
		b, err := table.Load([]byte(build(`["chore", "bug"]`)), "seam-b.toml")
		if err != nil {
			t.Fatalf("seam-b refused: %v", err)
		}
		va := seamValue(t, a, "kata.review-needs-work", "labels")
		vb := seamValue(t, b, "kata.review-needs-work", "labels")
		if va != vb {
			t.Errorf("two authored orderings crossed the seam as %q and %q; "+
				"read-back equality is byte equality over the canonical array", va, vb)
		}
	})

	t.Run("and the normalized value is NOT the serialized form", func(t *testing.T) {
		row := rowByID(t, m, "kata.review-needs-work")
		got, _ := tagValue(row.Writes, "labels")
		if !reflect.DeepEqual(got, []string{"needs work"}) {
			t.Errorf("normalized write value = %v; the JSON array is the SEAM "+
				"encoding, and the normalized value stays a member sequence", got)
		}
	})
}

func seamValue(t *testing.T, m *table.Model, rowID, key string) string {
	t.Helper()

	for _, w := range rowByID(t, m, rowID).KernelRow().Writes {
		if w.Key == key {
			return w.Value
		}
	}
	t.Fatalf("no %q write on the kernel row for %s", key, rowID)
	return ""
}

// REQ-94: "The expanded table dump MUST be derived from the normalized
// candidate-row value and MUST carry every field of it: row identity,
// source locator, outcome, predicate atoms (with block), gate ids,
// next-state tags, writes including `<clear>` entries, required-owned
// keys, and escape failure classes — plus the derived row kind column."
// HAPPY PATH
func TestReq94_DumpCarriesEveryFieldOfTheNormalizedValue(t *testing.T) {
	m := mustLoad(t, rdrFixture)
	out := table.Dump(m)

	witnesses := map[string]string{
		"row identity":    "rdr.reconcile-rewind",
		"source locator":  rowByID(t, m, "rdr.reconcile-rewind").SourceLocator,
		"outcome":         "reconcile-block",
		"predicate atoms": "status",
		"atom block":      string(table.BlockUnless),
		"gate ids":        "rdr-lock",
		"next-state tags": "rewind_scope",
		"<clear> writes":  table.ClearSentinel,
		"required-owned":  "prelock_lens",
		"escape classes":  "no_match",
		"derived kind":    string(table.KindEscape),
	}
	for field, want := range witnesses {
		if !strings.Contains(out, want) {
			t.Errorf("the dump does not carry %s (%q):\n%s", field, want, out)
		}
	}
}

// REQ-95: "`[dump]` carries exactly one key, `order`: a list of column
// identifiers." with the closed vocabulary "identity, source, kind,
// outcome, atoms, next, writes, requires_owned, gate, escape"
// BOUNDARY
//
// The vocabulary is fixed VERBATIM to the fixtures' spelling — it is the
// one place the record forbids re-derivation.
func TestReq95_DumpColumnVocabularyIsClosedAndVerbatim(t *testing.T) {
	want := []string{
		"identity", "source", "kind", "outcome",
		"atoms", "next", "writes", "requires_owned", "gate", "escape",
	}
	got := table.DumpColumns()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("dump column vocabulary = %v; want %v (the fixtures' spelling)", got, want)
	}

	m := mustLoad(t, rdrFixture)
	if !reflect.DeepEqual(m.DumpOrder, want) {
		t.Errorf("fixture [dump].order decoded as %v; want %v", m.DumpOrder, want)
	}

	t.Run("[dump] carries exactly one key", func(t *testing.T) {
		base := string(readFixture(t, rdrFixture))
		src := strings.Replace(base, "[dump]\norder = [", "[dump]\nwidth = 80\norder = [", 1)
		if src == base {
			t.Fatal("[dump] substitution did not apply")
		}
		_, err := table.Load([]byte(src), "dump-two-keys.toml")
		if err == nil {
			t.Fatal("a [dump] table carrying a second key loaded clean")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatUnknownSchemaField {
			t.Errorf("category = %q; want %q", cat, table.CatUnknownSchemaField)
		}
	})
}

// REQ-96: "`[dump]` settings MAY reorder the rendered columns; they MUST
// NOT omit a field. An `order` naming an unknown identifier, repeating
// one, or omitting one is a `malformed dump declaration` refused at load —
// not a silently truncated dump."
// ADVERSARIAL
func TestReq96_DumpOrderMayReorderButNeverOmit(t *testing.T) {
	t.Run("reordering is legal and takes effect", func(t *testing.T) {
		base := string(readFixture(t, rdrFixture))
		src := strings.Replace(base,
			`  "identity", "source", "kind", "outcome",
  "atoms", "next", "writes", "requires_owned", "gate", "escape",`,
			`  "escape", "gate", "requires_owned", "writes", "next", "atoms",
  "outcome", "kind", "source", "identity",`, 1)
		if src == base {
			t.Fatal("[dump].order substitution did not apply")
		}
		m, err := table.Load([]byte(src), "reordered.toml")
		if err != nil {
			t.Fatalf("a reordered [dump].order refused: %v", err)
		}
		if reflect.DeepEqual(m.DumpOrder, table.DumpColumns()) {
			t.Error("the reordered [dump].order did not take effect")
		}
		if table.Dump(m) == table.Dump(mustLoad(t, rdrFixture)) {
			t.Error("reordering the columns did not change the rendered dump")
		}
	})

	arms := map[string]string{
		"unknown identifier": "neg/neg-dump-unknown-column.toml",
		"repeated column":    "neg/neg-dump-repeated-column.toml",
		"omitted column":     "neg/neg-dump-omitted-column.toml",
	}
	for name, rel := range arms {
		t.Run(name, func(t *testing.T) {
			if got := loadCategory(t, rel); got != table.CatMalformedDumpDeclaration {
				t.Errorf("category = %q; want %q — not a silently truncated dump",
					got, table.CatMalformedDumpDeclaration)
			}
		})
	}
}

// REQ-97: "`[dump]` is **presentation, and MUST NOT reach the normalized
// value** … Normalization MUST ignore it entirely, so two models differing
// only in `[dump]` normalize to identical candidate-row sets"
// BOUNDARY
func TestReq97_DumpSettingsDoNotReachTheNormalizedValue(t *testing.T) {
	base := readFixture(t, rdrFixture)
	reordered := []byte(strings.Replace(string(base),
		`  "identity", "source", "kind", "outcome",
  "atoms", "next", "writes", "requires_owned", "gate", "escape",`,
		`  "escape", "gate", "requires_owned", "writes", "next", "atoms",
  "outcome", "kind", "source", "identity",`, 1))
	if string(reordered) == string(base) {
		t.Fatal("[dump].order substitution did not apply")
	}

	a, err := table.Load(base, rdrFixture)
	if err != nil {
		t.Fatalf("baseline refused: %v", err)
	}
	b, err := table.Load(reordered, rdrFixture)
	if err != nil {
		t.Fatalf("reordered refused: %v", err)
	}

	if len(a.Rows) == 0 || len(b.Rows) == 0 {
		t.Fatalf("one side normalized to no rows (%d / %d); the comparison "+
			"would be vacuous", len(a.Rows), len(b.Rows))
	}
	if !reflect.DeepEqual(a.Rows, b.Rows) {
		t.Error("two models differing only in [dump] normalized to different " +
			"candidate-row sets; normalization must ignore [dump] entirely")
	}
	// And the two models DO differ, so the fixture pair is real.
	if reflect.DeepEqual(a.DumpOrder, b.DumpOrder) {
		t.Error("the [dump].order permutation did not take effect; the control " +
			"compares two identical documents")
	}
}

// REQ-98: "Dump ordering MUST be deterministic across source key order.
// Rows sort by row identity, compared field by field as the tuple `(model
// id, rule id, expansion suffix)`; `model id` and `rule id` compare
// byte-lexicographically on the post-parse string, and the expansion
// suffix compares as a sequence — element by element,
// byte-lexicographically, a shorter sequence that is a prefix of a longer
// one sorting first, so the empty suffix sorts before any non-empty one."
// BOUNDARY
func TestReq98_RowsSortByTheIdentityTuple(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	t.Run("rows are emitted in identity-tuple order", func(t *testing.T) {
		got := rowIdentities(m)
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
		if !reflect.DeepEqual(got, want) {
			t.Errorf("row order =\n %v\nwant\n %v", got, want)
		}
	})

	t.Run("the empty suffix sorts before any non-empty one", func(t *testing.T) {
		// One rule spelled both ways: an unexpanded row and its expanded
		// siblings under the same rule id.
		base := string(readFixture(t, rdrFixture))
		src := strings.Replace(base,
			"source = \"rdr:terminal\"\n[rule.match.recognized]\nin = [\"finalized\", \"verdict-flapping\"]",
			"source = \"rdr:terminal\"\n[rule.match.recognized]\nin = [\"finalized\", \"verdict-flapping\", \"round-clean\"]", 1)
		if src == base {
			t.Fatal("alphabet substitution did not apply")
		}
		p, err := table.Load([]byte(src), "suffix-sort.toml")
		if err != nil {
			t.Fatalf("refused: %v", err)
		}
		ids := rowIdentities(p)
		// `rdr.reconcile-rewind` carries the empty suffix; it must sort
		// before every `rdr.terminal-archive#…`.
		empty := slices.Index(ids, "rdr.reconcile-rewind")
		for i, id := range ids {
			if strings.HasPrefix(id, "rdr.terminal-archive#") && i < empty {
				t.Errorf("%s sorts before the empty-suffix row; the empty suffix "+
					"sorts first (%v)", id, ids)
			}
		}
	})

	t.Run("deterministic across source key order", func(t *testing.T) {
		a := mustLoad(t, rdrFixture)
		b := mustLoad(t, "perm/rdr-keyorder.toml")
		if !reflect.DeepEqual(a.Rows, b.Rows) {
			t.Error("a key-order permutation produced a different row sequence")
		}
	})
}

// REQ-99: "Within each row, atoms sort by (key, block, operator token,
// literal), and next-state tags, writes, required-owned keys, gate ids,
// and escape classes sort by key."
// BOUNDARY
func TestReq99_WithinRowSequencesAreSorted(t *testing.T) {
	for _, rel := range []string{rdrFixture, kataFixture} {
		m := mustLoad(t, rel)
		for _, r := range m.Rows {
			t.Run(rel+" "+r.Identity(), func(t *testing.T) {
				if !sort.SliceIsSorted(r.Atoms, func(i, j int) bool {
					return compareAtoms(r.Atoms[i], r.Atoms[j]) < 0
				}) {
					t.Errorf("atoms are not sorted by (key, block, operator, "+
						"literal): %+v", r.Atoms)
				}
				assertSortedByKey(t, "NextTags", r.NextTags)
				assertSortedByKey(t, "Writes", r.Writes)
				if !slices.IsSorted(r.RequiresOwned) {
					t.Errorf("RequiresOwned is not sorted: %v", r.RequiresOwned)
				}
				if !slices.IsSorted(r.Gate) {
					t.Errorf("Gate is not sorted: %v", r.Gate)
				}
				if !slices.IsSorted(r.Escape) {
					t.Errorf("Escape is not sorted: %v", r.Escape)
				}
			})
		}
	}
}

func compareAtoms(a, b table.Atom) int {
	if c := strings.Compare(a.Key, b.Key); c != 0 {
		return c
	}
	if c := strings.Compare(string(a.Block), string(b.Block)); c != 0 {
		return c
	}
	if c := strings.Compare(a.Operator, b.Operator); c != 0 {
		return c
	}
	return slices.Compare(a.Literal, b.Literal)
}

func assertSortedByKey(t *testing.T, name string, tags []table.TagValue) {
	t.Helper()

	for i := 1; i < len(tags); i++ {
		if tags[i-1].Key > tags[i].Key {
			t.Errorf("%s is not sorted by key: %+v", name, tags)
			return
		}
	}
}

// REQ-100: "The **source locator MUST NOT participate in row ordering**."
// ADVERSARIAL
//
// Two rules whose identity tuples order one way and whose locators would
// order the other. A comparator including the locator flips them.
func TestReq100_LocatorDoesNotParticipateInRowOrdering(t *testing.T) {
	base := string(readFixture(t, rdrFixture))
	// Give the alphabetically-first rule a locator-bearing `source` that
	// would sort last, and vice versa.
	src := strings.Replace(base, `source = "rdr:prelock"`, `source = "zzz-last"`, 1)
	src = strings.Replace(src, `source = "rdr:terminal"`, `source = "aaa-first"`, 1)
	if src == base {
		t.Fatal("source substitutions did not apply")
	}

	m, err := table.Load([]byte(src), rdrFixture)
	if err != nil {
		t.Fatalf("refused: %v", err)
	}
	if !reflect.DeepEqual(rowIdentities(m), rowIdentities(mustLoad(t, rdrFixture))) {
		t.Errorf("row order changed when the `source` annotations were "+
			"reordered:\n got = %v", rowIdentities(m))
	}
}

// REQ-101: "**One source document carries exactly one model** … `duplicate
// model id` is never decidable within a single document. It is scoped to a
// caller that loads **several documents into one dump or lint invocation**
// and MUST be checked there, over the set of loaded models, before rows
// are merged for rendering. A loader handed one document MUST NOT report
// it."
// ADVERSARIAL (A10)
func TestReq101_DuplicateModelIDIsCrossDocumentOnly(t *testing.T) {
	a := mustLoad(t, "dup/dup-model-a.toml")
	b := mustLoad(t, "dup/dup-model-b.toml")
	distinct := mustLoad(t, "dup/dup-model-distinct.toml")

	t.Run("a single-document load must not report it", func(t *testing.T) {
		// Both documents load clean alone, which is the whole point: the
		// category is undecidable within one document.
		if a.ID != b.ID {
			t.Fatalf("the fixture pair does not share a model id: %q / %q", a.ID, b.ID)
		}
	})

	t.Run("the multi-document check reports it", func(t *testing.T) {
		err := table.CheckModelIDs([]*table.Model{a, b})
		if err == nil {
			t.Fatal("two documents sharing a model id passed the cross-document check")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatDuplicateModelID {
			t.Errorf("category = %q; want %q", cat, table.CatDuplicateModelID)
		}
	})

	t.Run("distinct model ids pass", func(t *testing.T) {
		if err := table.CheckModelIDs([]*table.Model{a, distinct}); err != nil {
			t.Errorf("two documents with distinct model ids refused: %v", err)
		}
	})

	t.Run("a single-model set passes", func(t *testing.T) {
		if err := table.CheckModelIDs([]*table.Model{a}); err != nil {
			t.Errorf("a one-model set refused: %v", err)
		}
	})
}

// REQ-102: "The dump MUST be emitted from a pre-sorted sequence of
// normalized rows, never by iterating a map and never by delegating key
// order to an encoder."
// ADVERSARIAL
//
// Go randomizes map iteration, so repeated dumps of one model in one
// process are the oracle: a map-iterating renderer produces different
// bytes across runs.
func TestReq102_DumpIsEmittedFromAPreSortedSequence(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	first := table.Dump(m)
	for i := range 32 {
		if got := table.Dump(m); got != first {
			t.Fatalf("dump %d differs from the first; the renderer iterates a "+
				"map or delegates key order to an encoder:\n got = %s\nwant = %s",
				i, got, first)
		}
	}

	// And a freshly loaded model dumps identically.
	if got := table.Dump(mustLoad(t, rdrFixture)); got != first {
		t.Error("a freshly loaded model dumped differently")
	}
}

// REQ-103: "Source order and rendered-row order MUST NOT decide a
// successful transition. … the normalized row set is unordered as far as
// selection is concerned … and normalization MUST NOT emit any positional
// field a consumer could tiebreak on."
// ADVERSARIAL
//
// Rule-order permutation is the assertion that actually excludes a
// positional tiebreak: sorting and checking no adjacent pair compares
// equal is satisfied by the defect it was meant to catch.
func TestReq103_NoPositionalFieldAndRuleOrderIsNeutralized(t *testing.T) {
	t.Run("rule-order permutation normalizes identically", func(t *testing.T) {
		a := mustLoad(t, rdrFixture)
		b := mustLoad(t, "perm/rdr-ruleorder.toml")
		if !reflect.DeepEqual(a.Rows, b.Rows) {
			t.Error("a rule-order permutation produced a different normalized " +
				"value; source order must not survive normalization")
		}
		if table.Dump(a) != table.Dump(b) {
			t.Error("a rule-order permutation dumped differently")
		}
	})

	t.Run("no positional field on the row", func(t *testing.T) {
		rt := reflect.TypeOf(table.Row{})
		for i := range rt.NumField() {
			switch strings.ToLower(rt.Field(i).Name) {
			case "index", "position", "ordinal", "seq", "sequence", "order", "line":
				t.Errorf("table.Row carries a %q field a consumer could tiebreak on",
					rt.Field(i).Name)
			}
		}
	})

	t.Run("the kernel row carries no positional field either", func(t *testing.T) {
		m := mustLoad(t, rdrFixture)
		kr := rowByID(t, m, "rdr.continue-prelock#large").KernelRow()
		if kr.RuleID == "" {
			t.Error("the kernel row lost its rule id")
		}
	})
}

// REQ-104: "The selection procedure itself — gate-then-count, and which
// refusals are escapable — is **JDR 0001 §D2's and the kernel's, and MUST
// NOT be restated here**"
// BOUNDARY
//
// The checkable form: this package mints no selection surface. It exposes
// no Resolve, Select, or Gate entry point of its own, and it consumes the
// kernel's rather than reimplementing it.
func TestReq104_ThisPackageRestatesNoSelectionProcedure(t *testing.T) {
	banned := []string{"Resolve", "Select", "Gate", "Count", "Escape", "Match"}
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
			if slices.Contains(banned, fn.Name.Name) {
				t.Errorf("%s declares a package-level func %q; the selection "+
					"procedure is the kernel's and must not be restated here",
					file, fn.Name.Name)
			}
		}
	}
}

// REQ-105: "A dump spanning more than one model MUST carry a
// model-unique `model id`"
// BOUNDARY
func TestReq105_MultiModelDumpsCarryTheModelID(t *testing.T) {
	a := mustLoad(t, rdrFixture)
	b := mustLoad(t, kataFixture)

	out := table.DumpAll([]*table.Model{a, b})
	for _, id := range []string{a.ID, b.ID} {
		if !strings.Contains(out, id) {
			t.Errorf("the multi-model dump does not carry the model id %q:\n%s", id, out)
		}
	}
	// Every row identity is model-qualified, so rows from two models never
	// collide in one dump.
	for _, m := range []*table.Model{a, b} {
		for _, r := range m.Rows {
			if !strings.HasPrefix(r.Identity(), m.ID+".") {
				t.Errorf("row identity %q is not qualified by the model id %q",
					r.Identity(), m.ID)
			}
		}
	}
}

// REQ-106: "Load-time validation failures MUST retain stable data-level
// categories before CLI mapping, including at minimum malformed TOML,
// unknown schema field, missing recognized outcome alphabet, malformed
// recognized outcome alphabet …, unknown tag, unknown context, cyclic
// context inheritance, write to non-owned tag, unknown accessor, malformed
// accessor declaration, malformed accessor binding, malformed tag
// declaration, malformed initial declaration, reserved tag value,
// unsupported version, malformed predicate atom …, malformed escape
// declaration …, **malformed model declaration** …, **malformed dump
// declaration** …, malformed outcome binding …, **malformed rule shape**
// …, malformed rule id …, duplicate rule id, duplicate model id, and
// `reserved_tag_key` (RDR 0008)."
// BOUNDARY
//
// The full category floor: every named category exists as a stable
// identifier, and each is witnessed by a fixture that trips it and no
// other (REQ-119, REQ-131).
func TestReq106_EveryNamedCategoryExistsAndIsWitnessed(t *testing.T) {
	witnesses := map[table.Category]string{
		table.CatMalformedTOML:                      "neg/neg-malformed-toml.toml",
		table.CatUnknownSchemaField:                 "neg/neg-unknown-field.toml",
		table.CatMissingRecognizedOutcomeAlphabet:   "neg/neg-no-alphabet.toml",
		table.CatMalformedRecognizedOutcomeAlphabet: "neg/neg-dup-alphabet-member.toml",
		table.CatUnknownTag:                         "neg/neg-unknown-tag-match.toml",
		table.CatUnknownContext:                     "neg/neg-unknown-context.toml",
		table.CatCyclicContextInheritance:           "neg/neg-cyclic-context.toml",
		table.CatWriteToNonOwnedTag:                 "neg/neg-write-observed.toml",
		table.CatUnknownAccessor:                    "neg/neg-gate-is-reader.toml",
		table.CatMalformedAccessorDeclaration:       "neg/neg-accessor-no-timeout.toml",
		table.CatMalformedAccessorBinding:           "neg/neg-owned-no-reader.toml",
		table.CatMalformedTagDeclaration:            "neg/neg-tagdecl-domain-on-bool.toml",
		table.CatMalformedInitialDeclaration:        "neg/neg-initial-bad-value.toml",
		table.CatReservedTagValue:                   "neg/neg-clear-as-write.toml",
		table.CatUnsupportedVersion:                 "neg/neg-version.toml",
		table.CatMalformedPredicateAtom:             "neg/neg-bad-exists.toml",
		table.CatMalformedEscapeDeclaration:         "neg/neg-escape-with-write.toml",
		table.CatMalformedModelDeclaration:          "neg/neg-no-version.toml",
		table.CatMalformedDumpDeclaration:           "neg/neg-dump-unknown-column.toml",
		table.CatMalformedOutcomeBinding:            "neg/neg-outcome-outside.toml",
		table.CatMalformedRuleShape:                 "neg/neg-no-write-block.toml",
		table.CatMalformedRuleID:                    "neg/neg-ruleid-hash.toml",
		table.CatDuplicateRuleID:                    "neg/neg-duplicate-rule-id.toml",
		table.CatReservedTagKey:                     "neg/neg-recognized-misnamed.toml",
	}

	t.Run("every category the floor names is exported", func(t *testing.T) {
		declared := table.Categories()
		for cat := range witnesses {
			if !slices.Contains(declared, cat) {
				t.Errorf("category %q is not in table.Categories()", cat)
			}
		}
		// duplicate model id is cross-document and has no single-document
		// witness, but it is still part of the floor.
		if !slices.Contains(declared, table.CatDuplicateModelID) {
			t.Errorf("category %q is not in table.Categories()", table.CatDuplicateModelID)
		}
	})

	for cat, rel := range witnesses {
		t.Run(string(cat), func(t *testing.T) {
			if got := loadCategory(t, rel); got != cat {
				t.Errorf("%s refused %q; want %q — one mutated fixture per "+
					"category, tripping it and no other", rel, got, cat)
			}
		})
	}
}

// REQ-107: "an absent `version` MUST refuse here and MUST NOT be folded
// into `unsupported version`, which would name a version the author never
// wrote"
// ADVERSARIAL
func TestReq107_AbsentVersionIsNotUnsupportedVersion(t *testing.T) {
	got := loadCategory(t, "neg/neg-no-version.toml")

	if got == table.CatUnsupportedVersion {
		t.Fatal("an absent `version` was folded into `unsupported_version`, " +
			"which names a version the author never wrote")
	}
	if got != table.CatMalformedModelDeclaration {
		t.Errorf("category = %q; want %q", got, table.CatMalformedModelDeclaration)
	}
}

// REQ-108: "**"Load" names the whole source-to-candidate-rows pipeline, not
// one callable.** Every category above MUST be refused by that pipeline
// before it yields candidate rows, whatever internal decomposition an
// implementation chooses."
// BOUNDARY
//
// The checkable form: a refused document yields NO candidate rows. A
// pipeline that returns partially normalized rows alongside a refusal
// fails.
func TestReq108_ARefusedDocumentYieldsNoCandidateRows(t *testing.T) {
	for _, rel := range []string{
		"neg/neg-unknown-context.toml",
		"neg/neg-outcome-outside.toml",
		"neg/neg-escape-with-write.toml",
		"neg/neg-version.toml",
	} {
		t.Run(rel, func(t *testing.T) {
			m, err := table.Load(readFixture(t, rel), rel)
			if err == nil {
				t.Fatal("loaded clean")
			}
			if m != nil && len(m.Rows) > 0 {
				t.Errorf("a refused document yielded %d candidate rows; every "+
					"category is refused BEFORE the pipeline yields rows", len(m.Rows))
			}
		})
	}
}

// REQ-109: "The categories are **data-level and MUST NOT depend on the CLI
// envelope**: this package MUST NOT import `internal/cli`"
// ADVERSARIAL
func TestReq109_PackageDoesNotImportInternalCLI(t *testing.T) {
	for _, file := range packageGoFiles(t) {
		f, err := parser.ParseFile(token.NewFileSet(), file, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", file, err)
		}
		for _, imp := range f.Imports {
			if strings.Contains(imp.Path.Value, "internal/cli") {
				t.Errorf("%s imports %s; load categories are data-level and must "+
					"not depend on the CLI envelope", file, imp.Path.Value)
			}
		}
	}
}

// REQ-110: "Category identifiers are **snake_case renderings of the prose
// names above**, anchored on `reserved_tag_key` … The Testing Strategy
// asserts on the category, so these identifiers are an API surface, not
// message text."
// BOUNDARY
func TestReq110_CategoryIdentifiersAreSnakeCase(t *testing.T) {
	t.Run("the anchor", func(t *testing.T) {
		if table.CatReservedTagKey != "reserved_tag_key" {
			t.Errorf("CatReservedTagKey = %q; want %q — the anchor RDR 0008 fixes",
				table.CatReservedTagKey, "reserved_tag_key")
		}
	})

	want := map[table.Category]string{
		table.CatMalformedTOML:                      "malformed_toml",
		table.CatUnknownSchemaField:                 "unknown_schema_field",
		table.CatMissingRecognizedOutcomeAlphabet:   "missing_recognized_outcome_alphabet",
		table.CatMalformedRecognizedOutcomeAlphabet: "malformed_recognized_outcome_alphabet",
		table.CatUnknownTag:                         "unknown_tag",
		table.CatUnknownContext:                     "unknown_context",
		table.CatCyclicContextInheritance:           "cyclic_context_inheritance",
		table.CatWriteToNonOwnedTag:                 "write_to_non_owned_tag",
		table.CatUnknownAccessor:                    "unknown_accessor",
		table.CatMalformedAccessorDeclaration:       "malformed_accessor_declaration",
		table.CatMalformedAccessorBinding:           "malformed_accessor_binding",
		table.CatMalformedTagDeclaration:            "malformed_tag_declaration",
		table.CatMalformedInitialDeclaration:        "malformed_initial_declaration",
		table.CatReservedTagValue:                   "reserved_tag_value",
		table.CatUnsupportedVersion:                 "unsupported_version",
		table.CatMalformedPredicateAtom:             "malformed_predicate_atom",
		table.CatMalformedEscapeDeclaration:         "malformed_escape_declaration",
		table.CatMalformedModelDeclaration:          "malformed_model_declaration",
		table.CatMalformedDumpDeclaration:           "malformed_dump_declaration",
		table.CatMalformedOutcomeBinding:            "malformed_outcome_binding",
		table.CatMalformedRuleShape:                 "malformed_rule_shape",
		table.CatMalformedRuleID:                    "malformed_rule_id",
		table.CatDuplicateRuleID:                    "duplicate_rule_id",
		table.CatDuplicateModelID:                   "duplicate_model_id",
		table.CatReservedTagKey:                     "reserved_tag_key",
	}
	for got, spelling := range want {
		if string(got) != spelling {
			t.Errorf("category %q; want the snake_case rendering %q", got, spelling)
		}
	}

	t.Run("every exported category is snake_case", func(t *testing.T) {
		for _, c := range table.Categories() {
			s := string(c)
			if s != strings.ToLower(s) || strings.ContainsAny(s, " -.") {
				t.Errorf("category %q is not a snake_case identifier", c)
			}
		}
	})
}

// REQ-111: "Cross-row findings — overlap, gap, dead row, read-before-write
// — carry RDR 0006's lint categories, not these."
// BOUNDARY
//
// The checkable form: no cross-row finding name appears in this package's
// category set, and the deliberately overlapping variant LOADS.
func TestReq111_CrossRowFindingsAreNotLoadCategories(t *testing.T) {
	for _, c := range table.Categories() {
		s := string(c)
		for _, lint := range []string{"overlap", "gap", "dead_row", "read_before_write", "unreachable"} {
			if strings.Contains(s, lint) {
				t.Errorf("category %q names the cross-row finding %q; those carry "+
					"RDR 0006's lint categories", c, lint)
			}
		}
	}
}

// REQ-112: "a check decidable from one rule plus the model's declarations
// is load-time and this RDR's; a check that must compare normalized rows
// against each other — overlap, gap, dead row, read-before-write — is
// graph lint and RDR 0006's. Ambiguous overlap between candidate rows is
// therefore a lint finding, not a load failure."
// INPUT EDGE
//
// Asserted positively: a deliberately overlapping model LOADS. Marking it
// as reported would be asserting against an unimplemented lint (REQ-117).
func TestReq112_AmbiguousOverlapIsNotALoadFailure(t *testing.T) {
	base := string(readFixture(t, rdrFixture))
	// Make `continue-prelock-cluster` predicate-identical to
	// `continue-prelock` on the same outcome: two rows overlap exactly.
	src := strings.Replace(base,
		"[rule.guard.all.cluster_ready]\neq = true\n", "", 1)
	if src == base {
		t.Fatal("overlap substitution did not apply")
	}

	m, err := table.Load([]byte(src), "overlapping.toml")
	if err != nil {
		t.Fatalf("a deliberately overlapping model refused at load with %v; "+
			"ambiguous overlap is cross-row and therefore RDR 0006's", err)
	}

	// What this RDR owes the lint handshake: the overlapping rows are
	// distinguishable by row identity and comparable by predicate set.
	a := rowByID(t, m, "rdr.continue-prelock#large")
	b := rowByID(t, m, "rdr.continue-prelock-cluster#large")
	if a.Identity() == b.Identity() {
		t.Error("the overlapping rows are not distinguishable by row identity")
	}
	if a.Outcome != b.Outcome {
		t.Errorf("the overlapping rows bind different outcomes (%q / %q); the "+
			"fixture is not an overlap", a.Outcome, b.Outcome)
	}
	if len(a.Atoms) == 0 || len(b.Atoms) == 0 {
		t.Error("the overlapping rows are not comparable by predicate set")
	}
}
