package table_test

// RDR 0002 sections A–C: the format and its carrier, the
// `[model.metadata]` extension namespace, and the model header and
// version gate.

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/table"
)

// REQ-1: "The transition model MUST be authored as sparse TOML data, not
// generated code and not a fully expanded Cartesian-product table."
// HAPPY PATH
//
// The checkable content of "sparse, not fully expanded" is that one
// authored rule yields more than one candidate row: a Cartesian-product
// source would have to author each row, so the row count would equal the
// rule count. The RDR fixture authors five rules and normalizes to nine
// rows, so the sparseness is observable in the value.
func TestReq1_SparseSourceExpandsToMoreRowsThanAuthoredRules(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	if got := len(m.Rows); got != 9 {
		t.Fatalf("RDR fixture normalized to %d candidate rows; want 9 (%v)",
			got, rowIdentities(m))
	}

	// Three of the five rules expand; two do not. If the source were an
	// expanded table each rule would contribute exactly one row.
	expanding := map[string]int{
		"continue-prelock":         2,
		"continue-prelock-cluster": 2,
		"draft-no-match-escape":    2,
		"terminal-archive":         2,
		"reconcile-rewind":         1,
	}
	for ruleID, want := range expanding {
		if got := len(rowsByRuleID(m, ruleID)); got != want {
			t.Errorf("rule %q yielded %d rows; want %d — the source is sparse "+
				"and normalization expands it", ruleID, got, want)
		}
	}
}

// REQ-2: "The source schema is the closed layout JDR 0001 §D7 fixes: root
// `outcomes`, root `terminal`, `[model]` (with the free-form sub-table
// `[model.metadata]`), `[initial]`, `[tags.<tag>]`, `[read.<id>]`,
// `[write.<id>]`, `[gate.<id>]`, `[context.<id>]`, `[[rule]]`, and
// `[dump]`."
// HAPPY PATH
//
// Every element of the closed layout is authored in the RDR fixture, so
// the layout is asserted by the fixture decoding into a model that
// carries each one. A loader admitting a narrower layout refuses the
// fixture; one admitting a wider layout is caught by REQ-5/REQ-6.
func TestReq2_ClosedLayoutDecodesEveryElement(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	if len(m.Outcomes) == 0 {
		t.Error("root `outcomes` did not decode")
	}
	if len(m.Terminal) == 0 {
		t.Error("root `terminal` did not decode")
	}
	if m.ID == "" || m.Version == 0 {
		t.Errorf("[model] did not decode: id=%q version=%d", m.ID, m.Version)
	}
	if len(m.Metadata) == 0 {
		t.Error("[model.metadata] did not decode")
	}
	if len(m.Initial) == 0 {
		t.Error("[initial] did not decode")
	}
	if len(m.Tags) == 0 {
		t.Error("[tags.<tag>] did not decode")
	}
	if len(m.Readers) == 0 {
		t.Error("[read.<id>] did not decode")
	}
	if len(m.Writers) == 0 {
		t.Error("[write.<id>] did not decode")
	}
	if len(m.Gates) == 0 {
		t.Error("[gate.<id>] did not decode")
	}
	if len(m.Rows) == 0 {
		t.Error("[[rule]] did not decode")
	}
	if len(m.DumpOrder) == 0 {
		t.Error("[dump] did not decode")
	}
}

// REQ-3: "Context predicates live under `[context.<id>.match.<tag>]`; rule
// predicates live under `[rule.match.<tag>]`, `[rule.guard.all.<tag>]`,
// and `[rule.guard.unless.<tag>]`; writes live under `[rule.write]`;
// explicit clears live in a rule-level `clear` list; gate accessor
// references live in a rule-level `gate` list of `[gate.<id>]` ids;
// modeled escape rows live in a rule-level `escape` list."
// HAPPY PATH
//
// Each authoring site is witnessed by the normalized artifact it produces
// on `continue-prelock` (all three atom blocks, a write, a gate list),
// `reconcile-rewind` (a clear), and `draft-no-match-escape` (an escape
// list). The context site is witnessed by an inherited match atom
// reaching the row.
func TestReq3_EveryAuthoringSiteReachesTheNormalizedRow(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	prelock := rowByID(t, m, "rdr.continue-prelock#large")

	// [context.<id>.match.<tag>]: `status.eq=Draft` comes from context
	// `draft` via `prelock` via `large-prelock`, never authored on the rule.
	if got := atomsOn(prelock, "status"); len(got) != 1 || got[0].Block != table.BlockMatch {
		t.Errorf("context match predicate did not reach the row: %+v", got)
	}
	// [rule.match.<tag>]: the outcome atom is lifted, but the inherited
	// `profile.in` expansion lands as a match atom.
	if got := atomsOn(prelock, "profile"); len(got) != 2 {
		t.Errorf("profile atoms = %d; want 2 (one match from expansion, one unless)", len(got))
	}
	// [rule.guard.all.<tag>] and [rule.guard.unless.<tag>].
	if len(atomsInBlock(prelock, table.BlockAll)) != 2 {
		t.Errorf("guard.all atoms = %+v; want iter.lt and finalized_at.exists",
			atomsInBlock(prelock, table.BlockAll))
	}
	if len(atomsInBlock(prelock, table.BlockUnless)) != 1 {
		t.Errorf("guard.unless atoms = %+v; want profile.eq=small",
			atomsInBlock(prelock, table.BlockUnless))
	}
	// [rule.write].
	if _, ok := tagValue(prelock.Writes, "stage"); !ok {
		t.Error("[rule.write] did not reach the row's writes")
	}
	// rule-level `gate`.
	if !reflect.DeepEqual(prelock.Gate, []string{"rdr-lock"}) {
		t.Errorf("rule-level gate list = %v; want [rdr-lock]", prelock.Gate)
	}

	// rule-level `clear`.
	rewind := rowByID(t, m, "rdr.reconcile-rewind")
	if v, ok := tagValue(rewind.Writes, "prelock_lens"); !ok || len(v) != 1 || v[0] != table.ClearSentinel {
		t.Errorf("rule-level clear rendered as %v (present=%v); want the <clear> sentinel", v, ok)
	}

	// rule-level `escape`.
	esc := rowByID(t, m, "rdr.draft-no-match-escape#round-clean")
	if len(esc.Escape) != 1 || esc.Escape[0] != "no_match" {
		t.Errorf("rule-level escape list = %v; want [no_match]", esc.Escape)
	}
}

// REQ-4: "`[model]` additionally carries an optional human `description`,
// and each `[[rule]]` an optional `source` — a provenance *annotation*
// carried alongside the locator, which is derived from `(model id, rule
// id)` and never from `source` (below). Both are admitted keys"
// HAPPY PATH
//
// Both keys are admitted (they do not refuse as unknown schema fields)
// and both survive to the normalized value. The second half — that the
// locator is not derived from `source` — is REQ-92's.
func TestReq4_DescriptionAndSourceAreAdmittedAndCarried(t *testing.T) {
	m := mustLoad(t, rdrFixture)

	if m.Description != "Representative RDR prelock and rewind slice." {
		t.Errorf("[model].description = %q; want the authored text", m.Description)
	}
	row := rowByID(t, m, "rdr.continue-prelock#large")
	if row.Source != "rdr:prelock" {
		t.Errorf("rule `source` annotation = %q; want %q", row.Source, "rdr:prelock")
	}
}

// REQ-5: "No other root key or table is admitted (strict decoding,
// below)."
// ADVERSARIAL
func TestReq5_NoOtherRootKeyOrTableIsAdmitted(t *testing.T) {
	base := readFixture(t, rdrFixture)

	cases := map[string]string{
		"unknown root key":       "\nflavor = \"vanilla\"\n",
		"unknown root table":     "\n[accessors.rdr-status]\nrole = \"rdr\"\n",
		"unknown model sub-key":  "", // handled below by the promoted fixture
		"unknown tags sub-table": "\n[tags.status.extra]\nnope = 1\n",
	}
	for name, suffix := range cases {
		if suffix == "" {
			continue
		}
		t.Run(name, func(t *testing.T) {
			src := append(append([]byte{}, base...), []byte(suffix)...)
			_, err := table.Load(src, rdrFixture)
			if err == nil {
				t.Fatalf("a document carrying an %s loaded clean", name)
			}
			if cat, ok := table.CategoryOf(err); !ok || cat != table.CatUnknownSchemaField {
				t.Errorf("category = %q (categorized=%v); want %q",
					cat, ok, table.CatUnknownSchemaField)
			}
		})
	}

	// The promoted one-mutation fixture: `model.flavor`.
	t.Run("promoted neg-unknown-field", func(t *testing.T) {
		if got := loadCategory(t, "neg/neg-unknown-field.toml"); got != table.CatUnknownSchemaField {
			t.Errorf("category = %q; want %q", got, table.CatUnknownSchemaField)
		}
	})
}

// REQ-6: "**Strict decoding is an obligation on this format, not a
// property of a library.** The decoder MUST reject unmapped keys so an
// unknown schema field is a stable refusal rather than a silent no-op."
// ADVERSARIAL
//
// "Stable refusal rather than a silent no-op" is the whole assertion: a
// permissive decoder returns a model with the unmapped key discarded and
// no error. The control therefore proves a refusal, not merely that the
// key is absent from the model.
func TestReq6_UnmappedKeyIsAStableRefusalNotASilentNoOp(t *testing.T) {
	src := append(readFixture(t, rdrFixture), []byte("\n[tags.status.unmapped]\nx = 1\n")...)

	m, err := table.Load(src, rdrFixture)
	if err == nil {
		t.Fatalf("unmapped key loaded clean with model %+v — a permissive "+
			"decoder passes this by discarding the key", m)
	}
	if cat, ok := table.CategoryOf(err); !ok || cat != table.CatUnknownSchemaField {
		t.Errorf("category = %q (categorized=%v); want %q", cat, ok, table.CatUnknownSchemaField)
	}
}

// REQ-7: "Use `github.com/pelletier/go-toml/v2` as the TOML parser
// candidate."
// BOUNDARY
func TestReq7_TOMLParserIsGoTOMLV2(t *testing.T) {
	mod, err := os.ReadFile(filepath.Join("..", "..", "go.mod"))
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	if !strings.Contains(string(mod), "github.com/pelletier/go-toml/v2") {
		t.Errorf("go.mod does not require github.com/pelletier/go-toml/v2:\n%s", mod)
	}
}

// REQ-8: "The load entry therefore takes **already-read bytes plus a
// source id** for the locator, not a filesystem path — this package
// performs no file I/O and no path resolution."
// BOUNDARY
//
// Two halves. The signature half is structural: Load takes ([]byte,
// string). The no-I/O half is asserted by parsing this package's own
// non-test sources and refusing any import of os, io/fs, or path/filepath
// — the packages a file read or a path resolution must travel.
func TestReq8_LoadTakesBytesAndSourceIDAndDoesNoFileIO(t *testing.T) {
	t.Run("signature is bytes plus source id", func(t *testing.T) {
		fn := reflect.TypeOf(table.Load)
		if fn.Kind() != reflect.Func {
			t.Fatalf("table.Load is not a func: %v", fn)
		}
		if fn.NumIn() != 2 {
			t.Fatalf("Load takes %d parameters; want 2 (bytes, source id)", fn.NumIn())
		}
		if got := fn.In(0); got != reflect.TypeOf([]byte(nil)) {
			t.Errorf("Load's first parameter is %v; want []byte — not a path", got)
		}
		if got := fn.In(1); got.Kind() != reflect.String {
			t.Errorf("Load's second parameter is %v; want a string source id", got)
		}
	})

	t.Run("package performs no file IO and no path resolution", func(t *testing.T) {
		banned := map[string]bool{
			`"os"`:            true,
			`"io/fs"`:         true,
			`"path/filepath"`: true,
			`"io/ioutil"`:     true,
		}
		for _, file := range packageGoFiles(t) {
			f, err := parser.ParseFile(token.NewFileSet(), file, nil, parser.ImportsOnly)
			if err != nil {
				t.Fatalf("parse %s: %v", file, err)
			}
			for _, imp := range f.Imports {
				if banned[imp.Path.Value] {
					t.Errorf("%s imports %s; this package performs no file I/O "+
						"and no path resolution", filepath.Base(file), imp.Path.Value)
				}
			}
		}
	})
}

// packageGoFiles lists the package's non-test Go sources.
func packageGoFiles(t *testing.T) []string {
	t.Helper()

	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read package dir: %v", err)
	}
	var out []string
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		out = append(out, name)
	}
	if len(out) == 0 {
		t.Fatal("no non-test Go sources in the package")
	}
	return out
}

// REQ-9: "**`[model.metadata]` is the one sanctioned extension
// namespace.** The loader MUST decode it as a free-form table, carry it
// through to the normalized model untouched, and MUST NOT interpret any
// key in it; strictness applies everywhere else."
// HAPPY PATH
func TestReq9_MetadataIsFreeFormAndUninterpreted(t *testing.T) {
	// A metadata table carrying keys that collide with the schema's own
	// root vocabulary. Interpreting any of them would change the model;
	// strictness applied inside would refuse the document.
	src := []byte(`outcomes = ["done"]

[model]
id = "m"
version = 1

[model.metadata]
outcomes = ["not-the-alphabet"]
version = 99
flavor = "vanilla"

[tags.recognized]
provenance = "recognized"
kind = "enum"
required = true

[[rule]]
id = "r"
[rule.match.recognized]
eq = "done"
[rule.write]
recognized = "done"
`)
	// The write above is illegal (recognized is not owned), so this
	// document is used only through the metadata leg: load the promoted
	// fixture for the positive case and assert non-interpretation here by
	// the categories that must NOT appear.
	_, err := table.Load(src, "metadata-collision.toml")
	if err != nil {
		if cat, ok := table.CategoryOf(err); ok {
			switch cat {
			case table.CatUnknownSchemaField, table.CatUnsupportedVersion,
				table.CatMalformedRecognizedOutcomeAlphabet:
				t.Errorf("category = %q: a key inside [model.metadata] was "+
					"interpreted or strictness was applied inside it", cat)
			}
		}
	}

	m := mustLoad(t, rdrFixture)
	if m.Version != 1 {
		t.Errorf("model version = %d; want 1", m.Version)
	}
	if _, isMap := any(m.Metadata).(map[string]any); !isMap {
		t.Errorf("metadata field is %T; want a free-form map[string]any", m.Metadata)
	}
}

// REQ-10: "It is a **model-level field and reaches no candidate row**, so
// it is outside the dump's field list and outside the Round-Trip
// invariant"
// BOUNDARY
//
// Two documents differing only in [model.metadata] normalize to identical
// candidate-row sets while carrying different normalized models.
func TestReq10_MetadataReachesNoCandidateRow(t *testing.T) {
	base := readFixture(t, rdrFixture)
	altered := []byte(strings.Replace(string(base),
		`review_cadence = "per-lock"`, `review_cadence = "per-round"`, 1))
	if string(altered) == string(base) {
		t.Fatal("fixture mutation did not apply")
	}

	a, err := table.Load(base, rdrFixture)
	if err != nil {
		t.Fatalf("baseline refused: %v", err)
	}
	b, err := table.Load(altered, rdrFixture)
	if err != nil {
		t.Fatalf("metadata-only variant refused: %v", err)
	}

	if !reflect.DeepEqual(a.Rows, b.Rows) {
		t.Error("candidate-row sets differ across a metadata-only edit; " +
			"metadata reaches no candidate row")
	}
	if reflect.DeepEqual(a.Metadata, b.Metadata) {
		t.Error("normalized models carry identical metadata across the edit; " +
			"the two models must differ")
	}
	// Outside the dump's field list: the rendered dump is unchanged.
	if table.Dump(a) != table.Dump(b) {
		t.Error("dump differs across a metadata-only edit; metadata is outside " +
			"the dump's field list")
	}
}

// REQ-11: "the normalized model MUST expose the decoded table as a field
// (Phase 2 deliverable), and scenario 1 MUST assert **value and nesting
// equality** against the authored table, not merely that its top-level
// key names survive"
// ADVERSARIAL
//
// The failing control this replaces is a top-level-key-name oracle, which
// a loader that decodes the table and discards every value passes. This
// asserts values and nesting depth.
func TestReq11_MetadataAssertedByValueAndNestingEquality(t *testing.T) {
	src := []byte(`outcomes = ["done"]
terminal = []

[model]
id = "m"
version = 1

[model.metadata]
owner = "flow-team"
count = 3
enabled = true

[model.metadata.nested]
depth = "one"

[model.metadata.nested.deeper]
depth = "two"
list = ["a", "b"]

[tags.recognized]
provenance = "recognized"
kind = "enum"
required = true
`)
	m, err := table.Load(src, "metadata-nesting.toml")
	if err != nil {
		t.Fatalf("metadata-nesting document refused: %v", err)
	}

	want := map[string]any{
		"owner":   "flow-team",
		"count":   int64(3),
		"enabled": true,
		"nested": map[string]any{
			"depth": "one",
			"deeper": map[string]any{
				"depth": "two",
				"list":  []any{"a", "b"},
			},
		},
	}
	if !reflect.DeepEqual(m.Metadata, want) {
		t.Errorf("metadata is not value-and-nesting equal to the authored table:\n"+
			" got = %#v\nwant = %#v", m.Metadata, want)
	}
}

// REQ-12: "Its internal shape is deliberately unconstrained (arbitrary
// keys, values, and nesting)"
// INPUT EDGE
func TestReq12_MetadataShapeIsUnconstrained(t *testing.T) {
	src := []byte(`outcomes = ["done"]

[model]
id = "m"
version = 1

[model.metadata]
"a key with spaces" = "ok"
"" = "empty key"
mixed = [1, "two", true]
"<clear>" = "not the sentinel here"
"has#hash" = "no expansion suffix ban inside metadata"
recognized = "not a tag"

[tags.recognized]
provenance = "recognized"
kind = "enum"
required = true
`)
	m, err := table.Load(src, "metadata-arbitrary.toml")
	if err != nil {
		t.Fatalf("arbitrary metadata refused: %v — its internal shape is "+
			"deliberately unconstrained", err)
	}
	if len(m.Metadata) != 6 {
		t.Errorf("metadata carries %d keys; want 6 — %#v", len(m.Metadata), m.Metadata)
	}
	if got := m.Metadata["<clear>"]; got != "not the sentinel here" {
		t.Errorf("the reserved-value rule reached inside metadata: %v", got)
	}
	if got := m.Metadata["has#hash"]; got != "no expansion suffix ban inside metadata" {
		t.Errorf("the `#` reservation reached inside metadata: %v", got)
	}
}

// REQ-13: "`[model]` MUST contain `id` and `version`. Version `1` is the
// only version this RDR accepts; any other version MUST be refused before
// normalization."
// ADVERSARIAL
func TestReq13_ModelHeaderAndVersionOne(t *testing.T) {
	t.Run("version 3 refuses unsupported_version", func(t *testing.T) {
		if got := loadCategory(t, "neg/neg-version.toml"); got != table.CatUnsupportedVersion {
			t.Errorf("category = %q; want %q", got, table.CatUnsupportedVersion)
		}
	})
	t.Run("absent id refuses malformed_model_declaration", func(t *testing.T) {
		if got := loadCategory(t, "neg/neg-no-model-id.toml"); got != table.CatMalformedModelDeclaration {
			t.Errorf("category = %q; want %q", got, table.CatMalformedModelDeclaration)
		}
	})
	t.Run("absent [model] refuses malformed_model_declaration", func(t *testing.T) {
		if got := loadCategory(t, "neg/neg-no-model-table.toml"); got != table.CatMalformedModelDeclaration {
			t.Errorf("category = %q; want %q", got, table.CatMalformedModelDeclaration)
		}
	})
	t.Run("version 1 is accepted", func(t *testing.T) {
		if m := mustLoad(t, rdrFixture); m.Version != 1 {
			t.Errorf("accepted version = %d; want 1", m.Version)
		}
	})
}

// REQ-14: "**The version check MUST run before strict field validation,
// not merely before normalization.** Loading MUST therefore proceed in
// two passes: read `[model]` permissively enough to obtain `version`,
// refuse on any value but `1`, and only then decode the document
// strictly."
// ADVERSARIAL
//
// A14's precedence control: a v2-SHAPED document — version = 2 plus a
// v2-only key the strict decoder would reject as an unknown schema field.
// A v1-shaped document with a bad version value trips unsupported_version
// under either ordering and therefore witnesses nothing about precedence.
func TestReq14_VersionGateRunsBeforeStrictFieldValidation(t *testing.T) {
	got := loadCategory(t, "neg/neg-v2-shaped.toml")

	if got == table.CatUnknownSchemaField {
		t.Fatalf("category = %q: strict decoding ran before the version gate", got)
	}
	if got != table.CatUnsupportedVersion {
		t.Errorf("category = %q; want %q", got, table.CatUnsupportedVersion)
	}
}

// REQ-15: "the ordering is `malformed TOML` → version gate → strict
// decoding → the remaining categories."
// BOUNDARY
//
// The two fixed precedences, each witnessed by a document that trips both
// arms and must report the earlier one.
func TestReq15_FixedPrecedenceOrdering(t *testing.T) {
	t.Run("malformed TOML precedes the version gate", func(t *testing.T) {
		// version = 2 AND a syntax error: the syntax error wins, because a
		// document that does not parse has no readable version.
		src := []byte("version = 2\nid = \"m\nnope")
		_, err := table.Load(src, "toml-before-version.toml")
		if err == nil {
			t.Fatal("a document that is not TOML loaded clean")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatMalformedTOML {
			t.Errorf("category = %q; want %q", cat, table.CatMalformedTOML)
		}
	})

	t.Run("version gate precedes strict decoding", func(t *testing.T) {
		if got := loadCategory(t, "neg/neg-v2-shaped.toml"); got != table.CatUnsupportedVersion {
			t.Errorf("category = %q; want %q", got, table.CatUnsupportedVersion)
		}
	})

	t.Run("promoted malformed TOML fixture", func(t *testing.T) {
		if got := loadCategory(t, "neg/neg-malformed-toml.toml"); got != table.CatMalformedTOML {
			t.Errorf("category = %q; want %q", got, table.CatMalformedTOML)
		}
	})
}

// REQ-16: "**Load is fail-fast: the first category a document trips is
// the refusal.** Beyond the two fixed precedences above, the order in
// which independent defects are checked is deliberately **unspecified**"
// BOUNDARY
//
// The assertable content is fail-fast, not the unspecified order: a
// document carrying two independent defects yields ONE refusal carrying
// ONE category, and that category is one of the two the defects name.
// Asserting a particular one would assert the order the clause leaves
// unspecified.
func TestReq16_LoadIsFailFastWithOneCategory(t *testing.T) {
	base := string(readFixture(t, rdrFixture))
	// Two independent defects: an unknown context on a rule, and a rule id
	// carrying the expansion separator.
	two := strings.Replace(base, `use = ["draft"]`, `use = ["nope"]`, 1)
	two = strings.Replace(two, `id = "reconcile-rewind"`, `id = "reconcile#rewind"`, 1)

	_, err := table.Load([]byte(two), "two-defects.toml")
	if err == nil {
		t.Fatal("a document carrying two independent defects loaded clean")
	}
	cat, ok := table.CategoryOf(err)
	if !ok {
		t.Fatalf("refusal carries no category: %v", err)
	}
	switch cat {
	case table.CatUnknownContext, table.CatMalformedRuleID:
	default:
		t.Errorf("category = %q; want one of %q / %q — the first category the "+
			"document trips", cat, table.CatUnknownContext, table.CatMalformedRuleID)
	}
}

// REQ-17: "What is forbidden is the single `load` entry point returning a
// list instead of a refusal"
// BOUNDARY
func TestReq17_LoadReturnsOneRefusalNotAList(t *testing.T) {
	fn := reflect.TypeOf(table.Load)
	if fn.NumOut() != 2 {
		t.Fatalf("Load returns %d values; want 2 (model, error)", fn.NumOut())
	}
	errType := reflect.TypeOf((*error)(nil)).Elem()
	if got := fn.Out(1); got != errType {
		t.Fatalf("Load's second return is %v; want error", got)
	}
	if got := fn.Out(1); got.Kind() == reflect.Slice {
		t.Errorf("Load returns a list of failures (%v) instead of a refusal", got)
	}

	// And the refusal itself is singular: one category, not an aggregate.
	_, err := table.Load(readFixture(t, "neg/neg-version.toml"), "neg-version.toml")
	if err == nil {
		t.Fatal("neg-version loaded clean")
	}
	if _, ok := table.CategoryOf(err); !ok {
		t.Errorf("the refusal is not a single categorized failure: %v", err)
	}
}

// REQ-18: "Because this class trips no category, it is **asserted
// positively, not by refusal**: the oracle is a count over the normalized
// value — the atom count for the merge case (Testing Strategy 2's
// distinct-literal control) and the write value's member sequence for the
// write case — each asserted to survive normalization intact."
// BOUNDARY
//
// The clause's own claim: neither absorbed-defect shape refuses. The
// counts themselves are REQ-126 and REQ-127; here the assertion is that
// both documents LOAD, so a refusal-shaped oracle would be wrong.
//
// SURVIVAL, not difference. Both legs run the REQ-126/127 five-delimiter
// parameterization (`,` `;` `|` space, empty), mirroring REQ-59's loop.
// The merge leg asserts the count over EVERY expanded row, not just the
// first: `reconcile-rewind` inherits two match-block `in` lists of two
// members each and expands to four rows, so indexing `rows[0]` would let a
// normalizer that merged atoms on rows 2-4 only pass. The write leg
// asserts the member sequence by exact content, sorted per REQ-57: mere
// inequality survives a `members[0]` truncation, which yields ["x,y"] vs
// ["x"] — distinct, yet the sequence did not survive.
func TestReq18_AbsorbedDefectsTripNoCategory(t *testing.T) {
	// Each fixture family spells the same two-member sets with one
	// candidate in-member separator embedded in a member. The expected
	// sequences are derived from that separator, not transcribed, so the
	// oracle stays tied to the fixture's own delimiter.
	seps := map[string]string{
		"comma": ",",
		"semi":  ";",
		"pipe":  "|",
		"space": " ",
		"empty": "",
	}

	t.Run("the merge case loads and every expanded row's atoms survive", func(t *testing.T) {
		for _, d := range []string{"comma", "semi", "pipe", "space", "empty"} {
			t.Run(d, func(t *testing.T) {
				rel := "delim/merge-delim-" + d + ".toml"
				m, err := table.Load(readFixture(t, rel), rel)
				if err != nil {
					t.Fatalf("%s refused with %v; this class trips no category and must "+
						"be asserted positively over the normalized value", rel, err)
				}
				// The positive oracle: a COUNT over the normalized value,
				// asserted on every row the contributing rule expands to.
				rows := rowsByRuleID(m, "reconcile-rewind")
				// Two inherited match-block `in` lists of two members each:
				// the §D6 expansion is the 2x2 product.
				if len(rows) != 4 {
					t.Fatalf("the contributing rule normalized to %d rows; want 4 — "+
						"a count asserted on one row of four is not a count over the "+
						"normalized value: %v", len(rows), rowIdentities(m))
				}
				for _, r := range rows {
					if got := len(atomsOn(r, "finalized_at")); got != 2 {
						t.Errorf("row %s carries %d atoms; want 2 — the two spellings "+
							"stay distinct on EVERY expanded row, never a refusal",
							r.Identity(), got)
					}
				}
			})
		}
	})

	t.Run("the write case loads and its member sequence survives intact", func(t *testing.T) {
		for _, d := range []string{"comma", "semi", "pipe", "space", "empty"} {
			t.Run(d, func(t *testing.T) {
				sep := seps[d]
				// Arm a authors ["x<sep>y","z"], arm b ["x","y<sep>z"]: the
				// same six characters, split two ways. Both are already in
				// byte-lexicographic order (REQ-57), since 'x' < 'y' < 'z'
				// and every separator sorts below 'y'.
				cases := []struct {
					rel  string
					want []string
				}{
					{"delim/write-delim-" + d + "-a.toml", []string{"x" + sep + "y", "z"}},
					{"delim/write-delim-" + d + "-b.toml", []string{"x", "y" + sep + "z"}},
				}
				var got [][]string
				for _, c := range cases {
					m, err := table.Load(readFixture(t, c.rel), c.rel)
					if err != nil {
						t.Fatalf("%s refused with %v; this class trips no category", c.rel, err)
					}
					// The positive oracle: the write value's MEMBER SEQUENCE,
					// asserted by content. Inequality alone would survive a
					// truncation to members[0].
					v := setWriteValue(t, m, "labels")
					if !reflect.DeepEqual(v, c.want) {
						t.Errorf("%s normalized labels to %v; want %v — the member "+
							"sequence must survive normalization intact", c.rel, v, c.want)
					}
					got = append(got, v)
				}
				// And, surviving intact, the two spellings stay distinct.
				if reflect.DeepEqual(got[0], got[1]) {
					t.Errorf("the two spellings normalized to %v alike under the %s "+
						"delimiter; each member is delimited unambiguously", got[0], d)
				}
			})
		}
	})
}
