package cli

// RDR 0021 — the exported JSON document (`0021:C2`): the field list and its
// normative spellings, the declared orders, the empty/absent rule, the
// required markers, the stability tier, and the round-trip invariants.
//
// C2's STABILITY clause binds THIS suite as much as any other consumer: a
// test here MUST NOT assert the field set's cardinality, a member's ordinal
// position, or a tail position. Every assertion below is therefore about a
// named member's presence, value, or order-within-its-own-sequence.

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// exportDocument runs the verb in its default (text) mode and returns the
// bare document stdout carries. C5 makes that stream the document itself,
// so this is the document under test throughout this file.
func exportDocument(t *testing.T, modelSrc string) string {
	t.Helper()

	requireGraphVerb(t)
	path := writeModel(t, modelSrc)
	stdout, _, err := runGraph(t, "--model", path)
	if err != nil {
		t.Fatalf("export refused: %v", err)
	}
	return stdout
}

// REQ-16: "The JSON document is a documentation artifact, versioned by a
// required leading `schema` field, initial value `intrastate.graph/1`"
// REQ-19 (`schema` member), ASSUMPTION A-2 (leading = first key on the wire).
// HAPPY PATH
func TestReq16_TheDocumentCarriesALeadingSchemaField(t *testing.T) {
	body := exportDocument(t, legalModel)
	doc := decodeDocument(t, body)

	if doc.Schema != schemaMarker {
		t.Errorf("schema = %q; want %q — C2 fixes the initial value of the "+
			"document's version marker", doc.Schema, schemaMarker)
	}

	// LEADING means first key on the wire. A2 reads C2's "required leading
	// `schema` field" together with XC's "schema first"; the encoder's
	// struct field order is what delivers it.
	trimmed := strings.TrimSpace(body)
	if !strings.HasPrefix(trimmed, `{"schema":`) {
		head := trimmed
		if len(head) > 60 {
			head = head[:60]
		}
		t.Errorf("the document does not open with the `schema` key: %q… — "+
			"C2 requires a LEADING `schema` field (XC restates it as "+
			"\"schema first\")", head)
	}
}

// REQ-19: "`schema`; `model` (`table.Model.ID`, the AUTHORED `[model] id`,
// never the `--model <path>` argument or any path-derived string); `class`;"
// REQ-46: "Identity is therefore invocation-independent: the same model
// exported from two checkouts is byte-identical"
// ADVERSARIAL — the defect is a `model` member derived from the PATH, which
// would make the document differ between two checkouts of one model and
// destroy the CI-diffability the record exists for.
func TestReq19And46_ModelIsTheAuthoredIDNeverThePathArgument(t *testing.T) {
	requireGraphVerb(t)

	// The same model source written to two differently-named files in two
	// directories. A path-derived identity differs between them; the
	// authored `[model] id` does not.
	dirA, dirB := t.TempDir(), t.TempDir()
	pathA := writeModelAt(t, dirA, "alpha.toml", legalModel)
	pathB := writeModelAt(t, dirB, "beta-renamed.toml", legalModel)

	outA, _, err := runGraph(t, "--model", pathA)
	if err != nil {
		t.Fatalf("export A refused: %v", err)
	}
	outB, _, err := runGraph(t, "--model", pathB)
	if err != nil {
		t.Fatalf("export B refused: %v", err)
	}

	docA := decodeDocument(t, outA)
	if docA.Model != "lintfix" {
		t.Errorf("model = %q; want the AUTHORED `[model] id` %q, never the "+
			"`--model <path>` argument or any path-derived string",
			docA.Model, "lintfix")
	}
	if outA != outB {
		t.Errorf("the same model exported from two paths produced different "+
			"bytes; D-identity makes identity invocation-INDEPENDENT, which "+
			"is what C3's (model, build) narrowing assumes and what CI "+
			"diffability rests on\n--- %s ---\n%s\n--- %s ---\n%s",
			pathA, outA, pathB, outB)
	}
	for _, path := range []string{pathA, "alpha", "beta-renamed", ".toml"} {
		if strings.Contains(docA.Model, path) {
			t.Errorf("model %q carries path-derived content %q",
				docA.Model, path)
		}
	}
}

// REQ-19 (`class` member).
// HAPPY PATH — the declared model class travels in the document, for both
// classes.
func TestReq19_ClassCarriesTheDeclaredModelClass(t *testing.T) {
	requireGraphVerb(t)

	for _, tc := range []struct {
		name string
		src  string
		want string
	}{
		{"state machine", legalModel, table.ClassStateMachine},
		{"decision table", decisionTableFixture, table.ClassDecisionTable},
	} {
		t.Run(tc.name, func(t *testing.T) {
			doc := decodeDocument(t, exportDocument(t, tc.src))
			if doc.Class != tc.want {
				t.Errorf("class = %q; want %q", doc.Class, tc.want)
			}
		})
	}
}

// REQ-20: "`tags[{name, provenance, kind, required, single_valued,
// domain}]` (`domain` carries the declaration's AUTHORED members and is
// present exactly when there are any — an `enum` with a non-empty `domain`
// — and ABSENT, its key omitted, otherwise)"
// REQ-28 (the ABSENT half of the empty/absent rule).
// BOUNDARY — every arm is asserted: an `enum` WITH members is PRESENT; the
// finite-but-member-less kinds (`bool`, bounded `int`, member-less `enum`)
// and the non-finite `scalar` are each ABSENT, key omitted, never `null`.
// Presence is NOT `guard.AssignmentCount` finiteness — that predicate
// reports finite for three kinds that author no `decl.Domain`.
func TestReq20And28_TagDomainIsPresentExactlyWhenAuthored(t *testing.T) {
	requireGraphVerb(t)

	// `status` authors enum members; `free` is a `scalar` carrying no finite
	// domain; `flag`, `count` and `recognized` are FINITE (or, for the
	// member-less enum, member-less) while authoring no members at all.
	const mixed = `
outcomes = ["go", "stop"]
terminal = ["done"]

[model]
id = "domains"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[tags.free]
provenance = "observed"
kind = "scalar"
required = true

[tags.flag]
provenance = "observed"
kind = "bool"
single_valued = true
required = true

[tags.count]
provenance = "observed"
kind = "int"
min = 0
max = 3
single_valued = true
required = true

[read.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"
read_back = true

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "advance-otherwise"
escape = ["no_match"]
[rule.match.recognized]
eq = "go"
`
	body := exportDocument(t, mixed)

	// Decode generically so ABSENT is distinguishable from null and from
	// the empty array.
	var generic struct {
		Tags []map[string]json.RawMessage `json:"tags"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(body)), &generic); err != nil {
		t.Fatalf("document is not JSON: %v\n%s", err, body)
	}
	if len(generic.Tags) == 0 {
		t.Fatalf("the document carries no `tags`: %s", body)
	}

	seen := map[string]map[string]json.RawMessage{}
	for _, raw := range generic.Tags {
		var name string
		if err := json.Unmarshal(raw["name"], &name); err != nil {
			t.Fatalf("a tag carries no decodable `name`: %v", err)
		}
		seen[name] = raw
	}

	finite, ok := seen["status"]
	if !ok {
		t.Fatalf("no `status` tag in the document: %s", body)
	}
	if _, has := finite["domain"]; !has {
		t.Errorf("the member-authoring tag `status` carries no `domain` " +
			"member; C2 requires it PRESENT exactly when the declaration " +
			"AUTHORS members — an `enum` with a non-empty domain")
	}

	opaque, ok := seen["free"]
	if !ok {
		t.Fatalf("no `free` tag in the document: %s", body)
	}
	if raw, has := opaque["domain"]; has {
		t.Errorf("the non-finite tag `free` carries `domain`: %s. C2 makes "+
			"`domain` present EXACTLY when the declaration AUTHORS members; "+
			"an inapplicable optional member is ABSENT, its key omitted, "+
			"never null", raw)
	}

	// The FINITE-but-member-less kinds are the discriminating arms: C2's
	// presence rule is over AUTHORED members, not over
	// `guard.AssignmentCount` finiteness, which reports finite for `bool`,
	// for a bounded `int`, and for `set` — none of which populates
	// `decl.Domain`. A member-less `enum` is the third case. Each omits
	// `domain`; `kind` already carries the type.
	for _, name := range []string{"flag", "count", "recognized"} {
		tag, ok := seen[name]
		if !ok {
			t.Fatalf("no %q tag in the document: %s", name, body)
		}
		if raw, has := tag["domain"]; has {
			t.Errorf("the finite-but-member-less tag %q carries `domain`: "+
				"%s. C2 predicates presence on AUTHORED members — an `enum` "+
				"with a non-empty domain — and explicitly NOT on "+
				"`guard.AssignmentCount` finiteness, so no `domain` array is "+
				"owed where no authority spells what it would carry",
				name, raw)
		}
	}

	// The five always-present members are named individually; this is not
	// a cardinality assertion (C2's STABILITY clause forbids that).
	for _, member := range []string{
		"name", "provenance", "kind", "required", "single_valued",
	} {
		if _, has := finite[member]; !has {
			t.Errorf("the `status` tag carries no %q member; C2 spells the "+
				"tag object as {name, provenance, kind, required, "+
				"single_valued, domain}", member)
		}
	}
}

// REQ-21: "`initial`; `terminal` (the declared predicate sets, carried as
// declared);"
// REQ-32: "the declared sets travel in the document for a consumer to
// evaluate."
// HAPPY PATH
func TestReq21And32_InitialAndTerminalTravelAsDeclared(t *testing.T) {
	body := exportDocument(t, legalModel)
	doc := decodeGeneric(t, body)

	for _, member := range []string{"initial", "terminal"} {
		raw, has := doc[member]
		if !has {
			t.Errorf("the document carries no %q member; C2 lists both, and "+
				"REQ-32 has the declared sets travel for a consumer to "+
				"evaluate", member)
			continue
		}
		if raw == nil {
			t.Errorf("%q is null; a declared collection renders `[]` when "+
				"empty, never null (C2)", member)
		}
	}

	// `terminal` carries the DECLARED predicate sets. The fixture declares
	// one terminal context (`done`, status eq b), so the member is a
	// non-empty sequence rather than an empty receipt.
	terminal, _ := doc["terminal"].([]any)
	if len(terminal) == 0 {
		t.Errorf("`terminal` is empty for a model declaring one terminal " +
			"context; C2 carries the declared predicate sets AS DECLARED")
	}
}

// REQ-22: "`rows[{identity, source, kind, outcome, atoms, next, writes,
// requires_owned, gate, escape, emit}]` — RDR 0002's dump field list in the
// row order `internal/table/dump.go::dumpColumns` fixes, atoms (each `{key,
// operator, literal[], block}`) in the atom order 0002 fixes, lowercase
// snake_case"
// REQ-23: "the `emit` member is the one `0010:C3` adds"
// REQ-108: "field names are lowercase snake_case against
// `internal/table/dump.go::dumpColumns`"
// HAPPY PATH
func TestReq22And23And108_RowsCarryTheDumpFieldVocabulary(t *testing.T) {
	body := exportDocument(t, legalModel)

	var generic struct {
		Rows []map[string]json.RawMessage `json:"rows"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(body)), &generic); err != nil {
		t.Fatalf("document is not JSON: %v\n%s", err, body)
	}
	if len(generic.Rows) == 0 {
		t.Fatalf("the document carries no rows: %s", body)
	}

	// The vocabulary is 0002's, read from the shipped accessor rather than
	// restated here — a column renamed there renames itself in this
	// assertion. This checks MEMBERSHIP, never cardinality: C2's STABILITY
	// clause forbids asserting the field set's size.
	row := generic.Rows[0]
	for _, column := range table.DumpColumns() {
		if _, has := row[column]; !has {
			t.Errorf("a row carries no %q member; C2 carries RDR 0002's dump "+
				"field vocabulary (`dumpColumns`), including the `emit` "+
				"member `0010:C3` appends", column)
		}
	}

	// lowercase snake_case: no member of any row may carry an uppercase
	// letter or a hyphen.
	for _, r := range generic.Rows {
		for name := range r {
			if name != strings.ToLower(name) || strings.Contains(name, "-") {
				t.Errorf("row member %q is not lowercase snake_case (XC)",
					name)
			}
		}
	}
}

// REQ-22 (the atom object and 0002's atom order).
// HAPPY PATH — each atom carries {key, operator, literal[], block}, and the
// atom sequence is in 0002's canonical order: (key, block, operator,
// literal), the same tuple `internal/table/normalize.go::compareAtoms` fixes.
func TestReq22_AtomsCarryTheirFieldsInTheCanonicalAtomOrder(t *testing.T) {
	body := exportDocument(t, legalModel)
	doc := decodeDocument(t, body)

	var checked bool
	for _, row := range doc.Rows {
		if len(row.Atoms) == 0 {
			continue
		}
		checked = true
		for _, a := range row.Atoms {
			if a.Key == "" || a.Operator == "" || a.Block == "" {
				t.Errorf("atom %+v is missing one of {key, operator, block}; "+
					"C2 spells the atom object as {key, operator, literal[], "+
					"block}", a)
			}
			if a.Literal == nil {
				t.Errorf("atom %+v carries a null `literal`; a set-valued "+
					"member is a JSON ARRAY, and an empty one renders `[]` "+
					"never null (C2, REQ-33)", a)
			}
		}
		// 0002's canonical atom order, compared field by field.
		if !slices.IsSortedFunc(row.Atoms, compareExportedAtoms) {
			t.Errorf("row %q's atoms are not in RDR 0002's canonical atom "+
				"order (key, block, operator, literal): %+v",
				row.Identity, row.Atoms)
		}
	}
	if !checked {
		t.Fatal("no row in the fixture carries an atom; the clause under " +
			"test is unexercised")
	}
}

// compareExportedAtoms is 0002's canonical atom order over the EXPORTED
// shape — the same tuple `internal/table/normalize.go::compareAtoms`
// compares: (key, block, operator, literal), the literal element by element.
func compareExportedAtoms(a, b graphAtom) int {
	if c := strings.Compare(a.Key, b.Key); c != 0 {
		return c
	}
	if c := strings.Compare(a.Block, b.Block); c != 0 {
		return c
	}
	if c := strings.Compare(a.Operator, b.Operator); c != 0 {
		return c
	}
	return slices.Compare(a.Literal, b.Literal)
}

// REQ-22 (row order).
// HAPPY PATH — rows are in 0002's canonical row order, the identity tuple
// (model id, rule id, expansion suffix) `compareRows` fixes.
func TestReq22_RowsAreInTheCanonicalRowOrder(t *testing.T) {
	doc := decodeDocument(t, exportDocument(t, legalModel))

	identities := make([]string, 0, len(doc.Rows))
	for _, r := range doc.Rows {
		identities = append(identities, r.Identity)
	}
	if !slices.IsSorted(identities) {
		t.Errorf("the exported rows are not in RDR 0002's canonical row "+
			"order; C2 carries \"the row order `dumpColumns` fixes\" and "+
			"D-identity makes a row's identity 0002's row identity: %v",
			identities)
	}
}

// REQ-24: "`groups[{context, rules}]`;"
// HAPPY PATH
func TestReq24_GroupsCarryContextAndRules(t *testing.T) {
	doc := decodeDocument(t, exportDocument(t, legalModel))

	if len(doc.Groups) == 0 {
		t.Fatalf("the document carries no `groups`; C2 lists "+
			"`groups[{context, rules}]` and the fixture declares several "+
			"rules: %+v", doc)
	}
	for _, g := range doc.Groups {
		if g.Context == "" {
			t.Errorf("a group carries an empty `context`: %+v", g)
		}
		if g.Rules == nil {
			t.Errorf("group %q carries a null `rules`; a declared "+
				"collection renders `[]` when empty, never null (C2)",
				g.Context)
		}
	}
}

// REQ-25: "`reach{abstraction, nodes[{id, values}], edges[{from, to,
// rule}]}` … `values` an OBJECT keyed by tag name whose every value is that
// tag's sorted, deduplicated value ARRAY"
// REQ-107 (each tag's `values` array sorted and deduplicated).
// HAPPY PATH — A7 resolved this shape: the object keyed by tag is the
// invention-free projection of `reach.go::Node`'s `Values
// map[string][]string`. Joined `key=value` strings exist only inside
// `(Node).key`'s internal fingerprint and are never an exported value.
func TestReq25And107_ReachNodeValuesIsATagKeyedObjectOfSortedArrays(t *testing.T) {
	doc := decodeDocument(t, exportDocument(t, legalModel))

	if doc.Reach == nil {
		t.Fatalf("the document carries no `reach` block; C2 lists " +
			"`reach{abstraction, nodes, edges}`")
	}
	if len(doc.Reach.Nodes) == 0 {
		t.Fatalf("the `reach` block carries no nodes for a model with a "+
			"declared root: %+v", doc.Reach)
	}

	for _, n := range doc.Reach.Nodes {
		if n.ID == "" {
			t.Errorf("a reach node carries an empty `id`; C2 fixes the id as " +
				"the node key `reach.go::(Node).key` renders")
		}
		if n.Values == nil {
			t.Errorf("node %q carries a null `values`; the member is an "+
				"OBJECT keyed by tag name, and an empty one renders `{}` "+
				"never null (C2)", n.ID)
			continue
		}
		for tag, values := range n.Values {
			if values == nil {
				t.Errorf("node %q tag %q carries a null value array; every "+
					"value is that tag's value ARRAY (C2)", n.ID, tag)
				continue
			}
			if !slices.IsSorted(values) {
				t.Errorf("node %q tag %q's values %v are not sorted; C2 "+
					"requires each tag's value array sorted and "+
					"deduplicated so construction order is unobservable",
					n.ID, tag, values)
			}
			for i := 1; i < len(values); i++ {
				if values[i] == values[i-1] {
					t.Errorf("node %q tag %q's values %v carry a duplicate; "+
						"C2 requires the array deduplicated",
						n.ID, tag, values)
					break
				}
			}
		}
	}
}

// REQ-25 (the edge object).
// HAPPY PATH
func TestReq25_ReachEdgesCarryFromToAndRule(t *testing.T) {
	doc := decodeDocument(t, exportDocument(t, legalModel))

	if doc.Reach == nil || len(doc.Reach.Edges) == 0 {
		t.Fatalf("the `reach` block carries no edges; the fixture declares "+
			"advancing rules, so the relation is non-empty: %+v", doc.Reach)
	}
	ids := map[string]bool{}
	for _, n := range doc.Reach.Nodes {
		ids[n.ID] = true
	}
	for _, e := range doc.Reach.Edges {
		if e.From == "" || e.To == "" || e.Rule == "" {
			t.Errorf("edge %+v is missing one of {from, to, rule}; C2 spells "+
				"the edge object as {from, to, rule}", e)
		}
		// An edge names nodes the document also carries: the relation is
		// self-contained, which is what lets the DOT arm render from the
		// document alone (A6).
		if !ids[e.From] {
			t.Errorf("edge %+v names a `from` node absent from `reach.nodes`", e)
		}
		if !ids[e.To] {
			t.Errorf("edge %+v names a `to` node absent from `reach.nodes`", e)
		}
	}
}

// REQ-26: "nodes sorted by node key and edges by (from, to, rule) so
// construction order is unobservable."
// REQ-107 (the same orders, restated in the cross-cutting policy).
// HAPPY PATH
func TestReq26And107_ReachNodesAndEdgesAreSorted(t *testing.T) {
	doc := decodeDocument(t, exportDocument(t, legalModel))
	if doc.Reach == nil {
		t.Fatal("the document carries no `reach` block")
	}

	ids := make([]string, 0, len(doc.Reach.Nodes))
	for _, n := range doc.Reach.Nodes {
		ids = append(ids, n.ID)
	}
	if !slices.IsSorted(ids) {
		t.Errorf("`reach.nodes` are not sorted by node key: %v — C2 sorts "+
			"them so construction order is unobservable", ids)
	}

	if !slices.IsSortedFunc(doc.Reach.Edges, func(a, b graphEdge) int {
		if c := strings.Compare(a.From, b.From); c != 0 {
			return c
		}
		if c := strings.Compare(a.To, b.To); c != 0 {
			return c
		}
		return strings.Compare(a.Rule, b.Rule)
	}) {
		t.Errorf("`reach.edges` are not sorted by (from, to, rule): %+v",
			doc.Reach.Edges)
	}
}

// REQ-27: "A tag with no finite declared domain carries the single value
// `<opaque>` (`reach.go::OpaqueValue`) in its array, passed through
// verbatim; this record attaches no meaning to that spelling."
// DOMAIN EDGE
func TestReq27_ANonFiniteOwnedTagCarriesTheOpaqueValueVerbatim(t *testing.T) {
	requireGraphVerb(t)

	// `note` is an owned `scalar`: no finite declared domain, so the
	// traversal abstracts it to held/absent and carries the opaque value.
	const opaqueModel = `
outcomes = ["go", "stop"]
terminal = ["done"]

[model]
id = "opaque"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[tags.note]
provenance = "owned"
kind = "scalar"
required = true

[read.own]
role = "t"
path = "t.own"
keys = ["status", "note"]
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ["status", "note"]
timeout = "2s"
read_back = true

[initial]
status = "a"
note = "anything"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "advance-otherwise"
escape = ["no_match"]
[rule.match.recognized]
eq = "go"
`
	doc := decodeDocument(t, exportDocument(t, opaqueModel))
	if doc.Reach == nil {
		t.Fatal("the document carries no `reach` block")
	}

	var found bool
	for _, n := range doc.Reach.Nodes {
		values, held := n.Values["note"]
		if !held {
			continue
		}
		found = true
		if !slices.Contains(values, "<opaque>") {
			t.Errorf("node %q holds `note` as %v; a tag with no finite "+
				"declared domain carries the single value `<opaque>` "+
				"(`reach.go::OpaqueValue`), passed through VERBATIM",
				n.ID, values)
		}
	}
	if !found {
		t.Fatal("no exported node holds the non-finite owned tag `note`; " +
			"the clause under test is unexercised")
	}
}

// REQ-28: "Every declared collection renders as an empty JSON array `[]`
// (or object `{}`) when it has no members — never `null`; an optional
// member that does not apply is ABSENT, its key omitted, never `null`."
// REQ-38: "C2.s schema-versioned JSON document; exact field spellings
// normative per C2.s list; empty declared collections render `[]`/`{}` and
// an inapplicable optional member is absent, never `null`; DOT styling
// explicitly non-normative."
// REQ-94: "producer code choosing `[]T{}`, never nil"
// ADVERSARIAL — the three states must be DISTINGUISHABLE on the wire. A nil
// Go slice marshals to `null`, which is the defect this clause forbids.
func TestReq28And94_NoDeclaredCollectionEverRendersAsNull(t *testing.T) {
	requireGraphVerb(t)

	// A decision table declares no owned tag and no terminal, so several
	// declared collections are genuinely EMPTY — the arm where a nil slice
	// would surface as `null`.
	for _, tc := range []struct {
		name string
		src  string
	}{
		{"state machine", legalModel},
		{"decision table", decisionTableFixture},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := exportDocument(t, tc.src)
			assertNoNullsAnywhere(t, body)
		})
	}
}

// assertNoNullsAnywhere walks the decoded document and fails on any null,
// at any depth. C2 admits exactly two renderings for a member with no
// content — `[]`/`{}` for a declared collection, and ABSENCE for an
// inapplicable optional member — and `null` is neither.
func assertNoNullsAnywhere(t *testing.T, body string) {
	t.Helper()

	var any_ any
	if err := json.Unmarshal([]byte(strings.TrimSpace(body)), &any_); err != nil {
		t.Fatalf("document is not JSON: %v\n%s", err, body)
	}
	var walk func(path string, v any)
	walk = func(path string, v any) {
		switch typed := v.(type) {
		case nil:
			t.Errorf("%s is null; C2 renders an empty declared collection as "+
				"`[]`/`{}` and OMITS an inapplicable optional member — "+
				"`null` is neither, and the three states must stay "+
				"distinguishable on the wire", path)
		case map[string]any:
			for k, sub := range typed {
				walk(path+"."+k, sub)
			}
		case []any:
			for i, sub := range typed {
				walk(path+"["+itoa(i)+"]", sub)
			}
		}
	}
	walk("$", any_)
}

// REQ-29: "The `reach` block carries a REQUIRED marker, `abstraction` with
// the token value `declared-over-approximation`"
// HAPPY PATH — the marker is what stops a consumer reading the merged,
// over-approximate relation as the runtime one.
func TestReq29_ReachCarriesTheRequiredAbstractionMarker(t *testing.T) {
	doc := decodeDocument(t, exportDocument(t, legalModel))

	if doc.Reach == nil {
		t.Fatal("the document carries no `reach` block")
	}
	if doc.Reach.Abstraction != abstractionMarker {
		t.Errorf("reach.abstraction = %q; want the token %q. C2 makes the "+
			"marker REQUIRED: the relation is the DECLARED "+
			"over-approximation (merged nodes, guard/observed atoms "+
			"unpruned), and a consumer reading it as the runtime relation "+
			"will over-count edges", doc.Reach.Abstraction, abstractionMarker)
	}
}

// REQ-31: "The document carries NO verdict or finding field — an export is
// never a lint pass — and NO per-node terminal marking"
// ADVERSARIAL — a clean, schema-stamped export that carried a verdict would
// read as certification and let a PR merge on the diagram while lint failed
// in another job (premortem P-6).
func TestReq31_TheDocumentCarriesNoVerdictFindingOrTerminalMarking(t *testing.T) {
	requireGraphVerb(t)

	// Both a clean model and one lint REFUSES: neither may carry a verdict.
	for _, tc := range []struct {
		name string
		src  string
	}{
		{"clean model", legalModel},
		{"model lint refuses", illegalModel},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := exportDocument(t, tc.src)
			doc := decodeGeneric(t, body)

			for _, forbidden := range []string{
				"verdict", "findings", "finding", "pass", "ok", "valid",
			} {
				if _, has := doc[forbidden]; has {
					t.Errorf("the document carries a top-level %q member; "+
						"C2 forbids a verdict or finding field — an export "+
						"is NEVER a lint pass, and the lint gate stays the "+
						"acceptance surface (`0006:C19`)", forbidden)
				}
			}

			// No PER-NODE terminal marking: which merged nodes satisfy a
			// `terminal` predicate set is the dead-end quantifier RDR 0015
			// owns (JDR 0001 §JD-23), and this record declares no evaluator.
			var probe struct {
				Reach *struct {
					Nodes []map[string]json.RawMessage `json:"nodes"`
				} `json:"reach"`
			}
			if err := json.Unmarshal([]byte(strings.TrimSpace(body)), &probe); err != nil {
				t.Fatalf("document is not JSON: %v", err)
			}
			if probe.Reach == nil {
				return
			}
			for _, n := range probe.Reach.Nodes {
				for _, forbidden := range []string{
					"terminal", "is_terminal", "terminal_satisfying", "dead_end",
				} {
					if _, has := n[forbidden]; has {
						t.Errorf("a reach node carries a %q member; C2 "+
							"declares NO per-node terminal marking — that "+
							"quantifier is RDR 0015's (JDR 0001 §JD-23), and "+
							"this record takes no side of it", forbidden)
					}
				}
			}
		})
	}
}

// REQ-33: "Set-valued members are JSON arrays, closing
// `0002:§round-trip-inverse-invariants`'s lossy set-literal rendering for
// this document"
// REQ-65: "`json-decode ∘ export = value identity on every field C2 lists`"
// REQ-84: "JSON round-trip — decode the exported document. Expected: value
// identity on every C2 field including exact set members (RT1)"
// HAPPY PATH — the recoverability 0002's TEXT dump deliberately declined.
func TestReq33And65And84_SetMembersSurviveDecodeAsExactArrays(t *testing.T) {
	requireGraphVerb(t)

	// A set-valued write, whose members are exactly what 0002's text dump
	// rendered lossily.
	const setModel = `
outcomes = ["go", "stop"]
terminal = ["done"]

[model]
id = "sets"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[tags.flavors]
provenance = "owned"
kind = "set"
elements = ["x", "y", "z"]
required = true

[read.own]
role = "t"
path = "t.own"
keys = ["status", "flavors"]
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ["status", "flavors"]
timeout = "2s"
read_back = true

[initial]
status = "a"
flavors = ["x", "y"]

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
flavors = ["y", "z"]

[[rule]]
id = "advance-otherwise"
escape = ["no_match"]
[rule.match.recognized]
eq = "go"
`
	body := exportDocument(t, setModel)

	// The round trip is VALUE identity, not "it decoded without error": a
	// green decode with a dropped member is exactly the fidelity loss RT1
	// exists to catch, so the re-encoded value is compared to the original.
	var first, second any
	if err := json.Unmarshal([]byte(strings.TrimSpace(body)), &first); err != nil {
		t.Fatalf("the exported document does not decode: %v\n%s", err, body)
	}
	reencoded, err := json.Marshal(first)
	if err != nil {
		t.Fatalf("re-encoding the decoded document: %v", err)
	}
	if err := json.Unmarshal(reencoded, &second); err != nil {
		t.Fatalf("the re-encoded document does not decode: %v", err)
	}
	if !jsonDeepEqual(first, second) {
		t.Errorf("decode ∘ export is not value-preserving:\n--- first ---\n"+
			"%v\n--- second ---\n%v", first, second)
	}

	// Every set-valued member on the wire is a JSON ARRAY, never a rendered
	// literal string: that is what closes 0002's lossy set-literal
	// rendering for this document.
	doc := decodeDocument(t, body)
	var sawSetWrite bool
	for _, row := range doc.Rows {
		for _, a := range row.Atoms {
			if a.Literal == nil {
				t.Errorf("atom %+v carries no `literal` array", a)
			}
		}
		if row.Identity != "" {
			sawSetWrite = true
		}
	}
	if !sawSetWrite {
		t.Fatal("no rows decoded; the clause under test is unexercised")
	}
}

// jsonDeepEqual compares two decoded JSON values structurally.
func jsonDeepEqual(a, b any) bool {
	left, errL := json.Marshal(a)
	right, errR := json.Marshal(b)
	return errL == nil && errR == nil && string(left) == string(right)
}

// REQ-35: "A consumer MUST therefore tolerate an unrecognized field and
// MUST NOT assert on the field set's cardinality, a member's ordinal
// position, or a tail position"
// REQ-34 (the append-only tier this restates), REQ-17 (additive within /1).
// ADVERSARIAL — this binds THIS suite. The oracle is that decoding into a
// view carrying only SOME of C2's members still works, which is what
// "tolerate an unrecognized field" means for a consumer.
func TestReq17And34And35_AConsumerToleratesUnrecognizedFields(t *testing.T) {
	body := exportDocument(t, legalModel)

	// A deliberately PARTIAL view: it names three members and ignores every
	// other. If decoding this required the full field set, the document
	// would not be additively evolvable, which is C2's whole promise.
	var partial struct {
		Schema string `json:"schema"`
		Model  string `json:"model"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(body)), &partial); err != nil {
		t.Fatalf("a consumer reading only {schema, model} could not decode "+
			"the document: %v — evolution within /1 is additive, so a "+
			"consumer ignoring unknown fields keeps working", err)
	}
	if partial.Schema != schemaMarker || partial.Model == "" {
		t.Errorf("the partial view decoded to %+v; the two named members "+
			"must survive a decode that ignores the rest", partial)
	}

	// The complementary half: an ADDED member does not break a consumer.
	// Simulated by injecting one and re-decoding through the same view.
	var full map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(body)), &full); err != nil {
		t.Fatalf("document is not JSON: %v", err)
	}
	full["a_future_member_this_build_does_not_know"] = "x"
	widened, err := json.Marshal(full)
	if err != nil {
		t.Fatalf("re-encoding: %v", err)
	}
	if err := json.Unmarshal(widened, &partial); err != nil {
		t.Errorf("a consumer could not decode a document carrying one added "+
			"member: %v — a member MAY be added in a minor (C2 STABILITY), "+
			"so an unrecognized field must be tolerated", err)
	}
}

// REQ-37: "this `schema` field versions the DOCUMENT, while `0029:C1`'s
// `schema_version` versions the ENVELOPE carrying it and is never projected
// into `data`."
// ADVERSARIAL — the two markers coexist at different levels and neither
// substitutes for the other. `schema_version` inside the document is the
// "one key meaning two things on one wire" defect.
func TestReq37_TheEnvelopeVersionIsNeverProjectedIntoTheDocument(t *testing.T) {
	requireGraphVerb(t)

	// The bare document (text mode) carries no envelope marker at all.
	body := exportDocument(t, legalModel)
	if strings.Contains(body, "schema_version") {
		t.Errorf("the bare document carries `schema_version`; that field "+
			"versions the ENVELOPE (`0029:C1`) and is NEVER projected into "+
			"the document\n%s", body)
	}

	// Under --as=json the envelope carries it at the TOP level, and the
	// embedded document still does not.
	path := writeModel(t, legalModel)
	stdout, _, err := runGraph(t, "--model", path, "--as=json")
	if err != nil {
		t.Fatalf("export --as=json refused: %v", err)
	}
	var env struct {
		SchemaVersion string          `json:"schema_version"`
		Data          json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &env); err != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", err, stdout)
	}
	if env.SchemaVersion == "" {
		t.Errorf("the envelope carries no top-level `schema_version`; "+
			"`0029:C1` makes it non-omitempty on both terminal records: %s",
			stdout)
	}
	var data map[string]any
	if err := json.Unmarshal(env.Data, &data); err != nil {
		t.Fatalf("`data` is not an object: %v\n%s", err, stdout)
	}
	if _, has := data["schema_version"]; has {
		t.Errorf("`data` carries `schema_version`; it rides the TOP LEVEL " +
			"only and is never projected into `data`, where it would " +
			"collide with this document's own `schema` marker (C2)")
	}
	if got, _ := data["schema"].(string); got != schemaMarker {
		t.Errorf("`data.schema` = %q; want %q — the DOCUMENT's own marker "+
			"travels inside the envelope", got, schemaMarker)
	}
}

// REQ-106: "This record claims no content-addressed identity and no
// replay-stable hash; the promise is byte identity of the emitted stream
// (C3, RT2), not a digest."
// BOUNDARY — a negative fence with a direct observable: no hash/digest
// member anywhere in the document.
func TestReq106_TheDocumentClaimsNoHashOrDigest(t *testing.T) {
	body := exportDocument(t, legalModel)
	doc := decodeGeneric(t, body)

	for _, forbidden := range []string{
		"hash", "digest", "sha256", "sha", "checksum", "fingerprint", "etag",
	} {
		if _, has := doc[forbidden]; has {
			t.Errorf("the document carries a %q member; XC records that this "+
				"record claims NO content-addressed identity and no "+
				"replay-stable hash — the promise is byte identity of the "+
				"stream, not a digest", forbidden)
		}
	}
}

// REQ-109: "The shared encoder is non-HTML-escaping, so tag values and rule
// ids reach the wire unescaped"
// DOMAIN EDGE — `<`, `>`, and `&` serialize as themselves. The `<opaque>`
// sentinel alone proves it for `<` and `>`, and a rule id carrying `&`
// proves it for the third.
func TestReq109_AngleBracketsAndAmpersandsReachTheWireUnescaped(t *testing.T) {
	requireGraphVerb(t)

	const htmlish = `
outcomes = ["go", "stop"]
terminal = ["done"]

[model]
id = "htmlish"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a<b", "x&y"]
single_valued = true
required = true

[read.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"
read_back = true

[initial]
status = "a<b"

[context.done]
[context.done.match.status]
eq = "x&y"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a<b"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "x&y"

[[rule]]
id = "advance-otherwise"
escape = ["no_match"]
[rule.match.recognized]
eq = "go"
`
	body := exportDocument(t, htmlish)

	// The needles are the ESCAPE SEQUENCES, not the raw characters: the
	// fixture authors `a<b` and `x&y`, so a correct document necessarily
	// carries `<` and `&` as literal bytes (the loop below requires exactly
	// that). What REQ-109 forbids is the encoder rewriting them as
	// `\u003c`/`\u003e`/`\u0026`.
	for _, escaped := range []string{"\\u003c", "\\u003e", "\\u0026"} {
		if strings.Contains(body, escaped) {
			t.Errorf("the document carries the HTML escape %s; the shared "+
				"encoder `clierr.WriteJSONLine` disables HTML escaping, so "+
				"`<`, `>`, and `&` serialize as THEMSELVES (XC, C3)\n%s",
				escaped, body)
		}
	}
	for _, raw := range []string{"a<b", "x&y"} {
		if !strings.Contains(body, raw) {
			t.Errorf("the document does not carry the authored value %q as "+
				"itself; no case folding is applied to values (XC)\n%s",
				raw, body)
		}
	}
}

// decisionTableFixture is MVV step 1's decision-table half: a stateless
// table declaring no owned tag, whose cells are claimed by guard atoms.
// Its shape is copied from `models/examples/pricing-decision-table.toml`,
// the CI-linted authored example.
const decisionTableFixture = `
outcomes = ["decide"]

[model]
id = "pricing"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.tier]
provenance = "observed"
kind = "enum"
domain = ["free", "paid"]
single_valued = true
required = true

[tags.region]
provenance = "observed"
kind = "enum"
domain = ["eu", "us"]
single_valued = true
required = true

[[rule]]
id = "free-eu"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.tier]
eq = "free"
[rule.guard.all.region]
eq = "eu"
[rule.emit]
plan = "basic"
dpa = "required"

[[rule]]
id = "free-us"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.tier]
eq = "free"
[rule.guard.all.region]
eq = "us"
[rule.emit]
plan = "basic"
dpa = "none"

[[rule]]
id = "paid-eu"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.tier]
eq = "paid"
[rule.guard.all.region]
eq = "eu"
[rule.emit]
plan = "pro"
dpa = "required"

[[rule]]
id = "paid-us"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.tier]
eq = "paid"
[rule.guard.all.region]
eq = "us"
[rule.emit]
plan = "pro"
dpa = "none"
`
