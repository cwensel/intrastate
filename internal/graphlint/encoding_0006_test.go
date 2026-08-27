package graphlint_test

// RDR 0006 — the composite keys the traversal and the finding identity
// close on must be INJECTIVE (REQ-87's sortable serialization, REQ-106's
// over-approximation contract, and RDR 0002 REQ-126/127's two-atom
// requirement over delimiter-bearing set members).
//
// Every key here is a delimiter-joined rendering of authored strings, and
// the loader accepts authored strings carrying those delimiters — no
// charset check exists anywhere in `internal/table`. A collision is
// therefore reachable from conforming source, and each collision below is
// false-green: it erases owned keys from the reachable set, terminates the
// fixpoint early, or collapses two findings into one identity.

import (
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/graphlint"
	"github.com/cwensel/intrastate/internal/table"
)

// delimTagDecls declares three owned tags whose names collide under an
// unescaped `;`-join: the footprint {"a;b"} and the footprint {"a","b"}
// both render `a;b;`.
const delimTagDecls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[tags."a;b"]
provenance = "owned"
kind = "enum"
domain = ["l"]
single_valued = true

[tags.a]
provenance = "owned"
kind = "enum"
domain = ["l"]
single_valued = true

[tags.b]
provenance = "owned"
kind = "enum"
domain = ["l"]
single_valued = true
`

// delimTagBody sends two rows out of the root establishing DIFFERENT owned
// tag sets: one writes the single tag `a;b`, the other writes `a` and `b`.
// The two successors are distinct owned-states and REQ-108 licenses no
// join between them.
const delimTagBody = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "writes-joined"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
"a;b" = "l"

[[rule]]
id = "writes-split"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
a = "l"
b = "l"
`

// RED-1 / REQ-106, REQ-108: `presenceKey` is the successor-join group id,
// so a collision between two DIFFERENT presence footprints sends two
// distinct successors through `joinNodes`, whose absence-dominates arm
// then drops every key held on only one side. That erases owned keys from
// the reachable set — the false-green direction REQ-106 forbids, since
// invariant 6 will certify a read of a tag the runtime does hold.
// ADVERSARIAL
func TestPresenceKeyIsInjectiveOverTagNamesCarryingTheJoinDelimiter(t *testing.T) {
	m := mustLoad(t, source(delimTagDecls, delimTagBody))

	// PRECONDITION: the loader really accepted a tag named `a;b`. Nothing
	// in `internal/table` charset-checks a tag name, so the collision is
	// reachable from conforming source rather than hypothetical.
	if _, declared := m.Tags["a;b"]; !declared {
		t.Fatalf("the loader did not carry a tag literally named `a;b`, so "+
			"this fixture does not exercise the collision; tags=%v",
			slices.Sorted(tagNames(m)))
	}

	nodes := graphlint.Reach(m)

	// Each row's successor must RETAIN the keys it wrote. Under the
	// unescaped join both footprints render `a;b;status;`, the two
	// successors land in one group, and `joinNodes` drops all three
	// owned keys — the observed reachable set is {status:[a]}, {status:[b]}.
	var sawJoined, sawSplit bool
	for _, n := range nodes {
		if _, held := n.Values["a;b"]; held {
			sawJoined = true
		}
		_, heldA := n.Values["a"]
		_, heldB := n.Values["b"]
		if heldA && heldB {
			sawSplit = true
		}
	}
	if !sawJoined {
		t.Errorf("no reachable node holds the owned tag `a;b`, though row "+
			"`writes-joined` leaves the root writing it. The successor "+
			"footprint collided with `{a, b}` under the unescaped `;`-join, "+
			"so `joinNodes` dropped the key — a reachable owned-state that "+
			"lint cannot see, which is the false-green direction REQ-106 "+
			"forbids.\nnodes:%s", renderNodes(nodes))
	}
	if !sawSplit {
		t.Errorf("no reachable node holds both owned tags `a` and `b`, "+
			"though row `writes-split` leaves the root writing both. The "+
			"successor footprint collided with `{\"a;b\"}` under the "+
			"unescaped `;`-join and the keys were erased.\nnodes:%s",
			renderNodes(nodes))
	}

	// The two successors are DIFFERENT owned-states, so they must be two
	// nodes. A single node standing for both is the merged reading REQ-108
	// does not license: it names an abstract state no path produces.
	if sawJoined && sawSplit {
		for _, n := range nodes {
			_, heldJoined := n.Values["a;b"]
			_, heldA := n.Values["a"]
			if heldJoined && heldA {
				t.Errorf("one node holds both `a;b` and `a`; the two rows "+
					"establish different owned tag sets and reach different "+
					"successors, so REQ-108 licenses no join between "+
					"them.\nnodes:%s", renderNodes(nodes))
			}
		}
	}
}

// tagNames yields the model's declared tag names.
func tagNames(m *table.Model) func(func(string) bool) {
	return func(yield func(string) bool) {
		for k := range m.Tags {
			if !yield(k) {
				return
			}
		}
	}
}

// RED-2 / REQ-106, REQ-83: `Node.key` renders `k=v1,v2;` per key, so the
// owned-states {a: "b;c=d"} and {a: "b", c: "d"} both render `a=b;c=d;`.
// That rendering is the `seen` de-duplication key BOTH terminal walks in
// `analysis.go` close on AND the fixpoint equality in `reach.go`, so a
// collision SUPPRESSES the second owned-state's finding entirely — one
// defect emitted where two are owed, and the one that vanished is the
// false-green direction REQ-106 forbids. It is also the node ELEMENT id
// the finding carries, so two distinct owned-states could not be told
// apart in the emitted set even were both emitted.
// ADVERSARIAL
func TestNodeKeyIsInjectiveOverValuesCarryingTheJoinDelimiters(t *testing.T) {
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["start", "stop"]
single_valued = true
required = true

[tags.a]
provenance = "owned"
kind = "enum"
domain = ["b;c=d", "b"]
single_valued = true

[tags.c]
provenance = "owned"
kind = "enum"
domain = ["d"]
single_valued = true
`
	// Two rows leave the root establishing owned-states that render
	// identically under the unescaped `=`/`,`/`;` join: {status: stop,
	// a: "b;c=d"} and {status: stop, a: "b", c: "d"} both render
	// `a=b;c=d;status=stop;`. Neither satisfies the declared terminal and
	// neither is the source of a non-escape row, so invariant 2 owes a
	// `graph-dead-end` finding against EACH.
	const body = `
terminal = ["fin"]

[initial]
status = "start"

[context.fin.match.status]
eq = "start"

[[rule]]
id = "one-key"
[rule.match.status]
eq = "start"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "stop"
a = "b;c=d"

[[rule]]
id = "two-keys"
[rule.match.status]
eq = "start"
[rule.match.recognized]
eq = "stop"
[rule.write]
status = "stop"
a = "b"
c = "d"
`
	m := mustLoad(t, source(decls, body))
	nodes := graphlint.Reach(m)

	// PRECONDITION: both owned-states really are reachable and the
	// delimiter-bearing value really did survive normalization, so the
	// assertion below is about the RENDERING and not about reachability.
	var sawOneKey, sawTwoKeys bool
	for _, n := range nodes {
		if slices.Contains(n.Values["a"], "b;c=d") {
			sawOneKey = true
		}
		if _, held := n.Values["c"]; held {
			sawTwoKeys = true
		}
	}
	if !sawOneKey {
		t.Fatalf("no reachable node holds the authored value `b;c=d` for "+
			"`a`, so this fixture does not exercise the collision.\nnodes:%s",
			renderNodes(nodes))
	}
	if !sawTwoKeys {
		t.Fatalf("no reachable node holds `c`, so the second owned-state is "+
			"not reachable and the collision is untested.\nnodes:%s",
			renderNodes(nodes))
	}

	// The defect, observed through the engine. Both owned-states are dead
	// ends, so invariant 2 owes two findings at two DISTINCT element ids.
	r := graphlint.Run(graphlint.NewRequest(m))
	elements := map[string]bool{}
	for _, f := range withCode(r, graphlint.CodeDeadEnd) {
		elements[f.Element] = true
	}
	if len(elements) != 2 {
		t.Errorf("invariant 2 named %d distinct owned-state(s) as dead ends, "+
			"want exactly 2. {a: \"b;c=d\"} and {a: \"b\", c: \"d\"} are different "+
			"owned-states that render the same `Node.key`, so the `seen` "+
			"de-duplication swallowed one finding outright and the node "+
			"element id could not tell them apart. That same rendering is "+
			"the fixpoint equality in `reach`, so the collision can "+
			"terminate the widening early too. More than two would mean the "+
			"fixture reaches an owned-state it does not intend, which would "+
			"satisfy a floor check without exercising the "+
			"collision.\nelements=%v\nnodes:%s"+
			"\nreport:%s",
			len(elements), slices.Sorted(maps.Keys(elements)),
			renderNodes(nodes), render(r))
	}

	// Pin the two element identities, not just their count. The node
	// element id IS the colliding `Node.key`, so asserting the exact pair
	// is what distinguishes "two findings at two distinct ids" from "two
	// findings that happen to differ somewhere else".
	wantElements := []string{
		`node{a=b,;c=d,;status=stop,}`,
		`node{a=b\;c\=d,;status=stop,}`,
	}
	if got := slices.Sorted(maps.Keys(elements)); !slices.Equal(got, wantElements) {
		t.Errorf("invariant 2's dead-end element ids are %v, want %v. The "+
			"element id is the rendered `Node.key`, so the escaped and "+
			"unescaped owned-states must appear as these two distinct "+
			"identities.\nreport:%s", got, wantElements, render(r))
	}
}

// RED-3 / RDR 0002 REQ-126/127 and RDR 0006 REQ-87: `in = ["a,b", "c"]`
// and `in = ["a", "b,c"]` MUST yield two atoms in the normalized set, and
// the fingerprint that closes the finding-identity tuple must tell them
// apart. Under the unescaped `,`-join both render `a,b,c`, so two rows
// expanded from the two contexts carry the SAME finding identity — the
// emitted set collapses and its order stops being determined by the tuple.
// ADVERSARIAL
func TestFingerprintIsInjectiveOverSetMembersCarryingTheJoinDelimiter(t *testing.T) {
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["start", "done"]
single_valued = true
required = true

[tags.mark]
provenance = "observed"
kind = "scalar"
`
	const body = `
terminal = ["fin"]

[initial]
status = "start"

[context.fin]
[context.fin.match.status]
eq = "done"

[[rule]]
id = "delim-a"
[rule.match.status]
eq = "start"
[rule.match.recognized]
eq = "go"
[rule.guard.all.mark]
in = ["a,b", "c"]
[rule.write]
status = "done"

[[rule]]
id = "delim-b"
[rule.match.status]
eq = "start"
[rule.match.recognized]
eq = "go"
[rule.guard.all.mark]
in = ["a", "b,c"]
[rule.write]
status = "done"
`
	m := mustLoad(t, source(decls, body))
	a := rowByRuleID(t, m, "delim-a")
	b := rowByRuleID(t, m, "delim-b")

	// PRECONDITION (RDR 0002 REQ-126/127): the loader carried the two
	// member sequences through normalization DISTINCTLY — each is two
	// atoms' worth of literal, not one flattened three-member set.
	litA, litB := markLiteral(t, a), markLiteral(t, b)
	wantA := []string{"a,b", "c"}
	wantB := []string{"a", "b,c"}
	if !slices.Equal(litA, wantA) || !slices.Equal(litB, wantB) {
		t.Fatalf("the loader normalized the two `in` literals to %v and %v, "+
			"want %v and %v. Asserting only that they DIFFER would let an "+
			"unrelated malformed value satisfy the precondition without "+
			"exercising the delimiter collision this test is about; the "+
			"exact canonical member sets are what RDR 0002 REQ-126/127 "+
			"require.", litA, litB, wantA, wantB)
	}

	fpA, fpB := graphlint.Fingerprint(a), graphlint.Fingerprint(b)
	if fpA == fpB {
		t.Errorf("rows whose `in` members are %v and %v share the "+
			"fingerprint %q. RDR 0002 REQ-126/127 require the two to yield "+
			"TWO atoms in the normalized set, and REQ-87's fingerprint "+
			"closes the finding-identity tuple — so the collision gives two "+
			"distinct rows one identity, collapsing the emitted set and "+
			"leaving its order undetermined by the tuple.", litA, litB, fpA)
	}
}

// markLiteral returns the row's `mark` `in`-atom literal. The atom lives in
// a GUARD block deliberately: a match-block `in` EXPANDS to one row per
// member under RDR 0002 §D6, so its multi-member literal never survives
// normalization and could not collide. A guard `in` keeps its whole member
// sequence in one atom, which is the literal REQ-126/127 speak about and
// the one the fingerprint must render injectively.
func markLiteral(t *testing.T, row table.Row) []string {
	t.Helper()

	for _, a := range row.Atoms {
		if a.Key == "mark" && a.Operator == "in" {
			return a.Literal
		}
	}
	t.Fatalf("row %q carries no `mark` `in` atom; atoms=%v", row.RuleID, row.Atoms)
	return nil
}

// The loader's OWN delimiter fixtures, linted here rather than
// re-authored: RDR 0002 REQ-126/127 name `["a,b", "c"]` vs `["a", "b,c"]`
// verbatim and `internal/table/testdata/delim/merge-delim-comma.toml`
// already ships both as contexts over `finalized_at`. A match-block `in`
// EXPANDS per member (§D6), so those four rows are already distinct — this
// pins that they STAY distinct and STAY readable once the fields are
// escaped, so the fix cannot regress the requirement it was written for.
// REGRESSION
func TestFingerprintSeparatesTheLoaderDelimFixtureArms(t *testing.T) {
	m := loadTestdata(t, "../table/testdata/delim/merge-delim-comma.toml")

	// `reconcile-rewind` uses BOTH delim contexts, so its expansion is the
	// cross product of the two `in` member sequences — four rows, each
	// carrying one member from each arm.
	seen := map[string]bool{}
	var rows int
	for _, row := range m.Rows {
		if row.RuleID != "reconcile-rewind" {
			continue
		}
		rows++
		fp := graphlint.Fingerprint(row)
		if seen[fp] {
			t.Errorf("two expansions of `reconcile-rewind` share the "+
				"fingerprint %q; RDR 0002 REQ-126/127 require the "+
				"delimiter-bearing members to stay distinct, and the "+
				"finding-identity tuple closes on this rendering", fp)
		}
		seen[fp] = true

		// Never a hash (REQ-87): the authored MEMBERS are still legible,
		// not just the constant atom key. Asserting only on
		// `finalized_at` would pass against an implementation that
		// replaced the delimiter-bearing literals with opaque tokens,
		// which is the failure mode REQ-87 names.
		//
		// §D6 expands a match-block `in` per member, so each expansion
		// carries one member from each arm: `a,b` or `c` from `delim-a`,
		// `a` or `b,c` from `delim-b`. Escaped, the comma-bearing members
		// read `a\,b` and `b\,c` — still eyeball-distinguishable from the
		// plain `a` and `c`, which is what "readable" asks.
		members := []string{`a\,b`, `c`, `a`, `b\,c`}
		var carried int
		for _, want := range members {
			if strings.Contains(fp, want) {
				carried++
			}
		}
		if carried == 0 {
			t.Errorf("the fingerprint %q carries none of the authored "+
				"members %v; a canonical sortable serialization is not a "+
				"hash", fp, members)
		}
		if !strings.Contains(fp, "finalized_at") {
			t.Errorf("the fingerprint %q does not carry the atom key "+
				"`finalized_at`; a canonical sortable serialization is not "+
				"a hash", fp)
		}
	}

	// Across the four expansions, BOTH delimiter-bearing members must
	// appear somewhere. A fingerprint that dropped or opaqued them would
	// still satisfy the per-row check above via the plain members.
	var sawAB, sawBC bool
	for fp := range seen {
		if strings.Contains(fp, `a\,b`) {
			sawAB = true
		}
		if strings.Contains(fp, `b\,c`) {
			sawBC = true
		}
	}
	if !sawAB || !sawBC {
		t.Errorf("the four expansions do not carry both delimiter-bearing "+
			"members in escaped form (a\\,b=%v b\\,c=%v). RDR 0002 "+
			"REQ-126/127 name this pair verbatim, and REQ-87 requires the "+
			"rendering stay readable rather than opaque.\nfingerprints=%v",
			sawAB, sawBC, slices.Sorted(maps.Keys(seen)))
	}
	if rows != 4 {
		t.Fatalf("`reconcile-rewind` expanded to %d rows, want 4 (the cross "+
			"product of two two-member `in` arms); the fixture no longer "+
			"exercises REQ-126/127 and this regression is vacuous", rows)
	}
}

// loadTestdata loads a fixture authored under another package's testdata,
// so the loader's OWN delimiter fixtures — the ones RDR 0002 REQ-126/127
// ships — are linted here rather than re-authored and drifting from them.
func loadTestdata(t *testing.T, path string) *table.Model {
	t.Helper()

	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading fixture %s: %v", path, err)
	}
	m, err := table.Load(src, filepath.Base(path))
	if err != nil {
		t.Fatalf("fixture %s must load clean; refused: %v", path, err)
	}
	return m
}

// REQ-87 requires the fingerprint carry RDR 0002's canonical atom order,
// and the escaping must not take that away.
//
// The escape itself is NOT order-preserving in general — it inserts `\`
// (U+005C) at the delimiter's position, so `["a,b"]` and `["a0"]` transpose
// under it. Sortability is therefore carried by `compareAtoms`, which
// orders atoms over the canonical AUTHORED members BEFORE rendering, per
// `0002:C19`. This test pins that the order the fingerprint comes out in is
// the authored one; `TestCanonicalAtomOrderSurvivesADelimiterBearingLiteral`
// pins the transposing pair the escape alone would get wrong.
// BOUNDARY
func TestEscapingPreservesTheFingerprintSortOrder(t *testing.T) {
	// Rows whose only difference is a guard `in` literal, so the sort over
	// their fingerprints is decided by the escaped literal alone.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["start", "done"]
single_valued = true
required = true

[tags.mark]
provenance = "observed"
kind = "scalar"
`
	body := func(id, literal string) string {
		return `
[[rule]]
id = "` + id + `"
[rule.match.status]
eq = "start"
[rule.match.recognized]
eq = "go"
[rule.guard.all.mark]
in = ` + literal + `
[rule.write]
status = "done"
`
	}
	const preamble = `
terminal = ["fin"]

[initial]
status = "start"

[context.fin]
[context.fin.match.status]
eq = "done"
`
	// Two literals ordered `["a"] < ["a,b"]` by their AUTHORED members: `a`
	// is a prefix of `a,b`, so the shorter sorts first, and the escaping
	// must not flip that.
	m := mustLoad(t, source(decls, preamble+
		body("plain", `["a"]`)+body("delim", `["a,b"]`)))

	fpPlain := graphlint.Fingerprint(rowByRuleID(t, m, "plain"))
	fpDelim := graphlint.Fingerprint(rowByRuleID(t, m, "delim"))

	if fpPlain >= fpDelim {
		t.Errorf("the fingerprints are out of canonical order: %q must sort "+
			"before %q, because the authored member `a` precedes `a,b` "+
			"element-by-element under `0002:C19`. REQ-87's canonical atom "+
			"order is not being carried into the rendering.",
			fpPlain, fpDelim)
	}

	// Never a hash (REQ-87): the escaped member is still legible — the
	// authored `a,b` reads as `a\,b`, not as an opaque token.
	if !strings.Contains(fpDelim, `a\,b`) {
		t.Errorf("the fingerprint %q does not carry the escaped member "+
			"`a\\,b`; a canonical sortable serialization is not a hash",
			fpDelim)
	}
}

// REQ-87 / `0002:C19`: the canonical atom order is over the AUTHORED member
// sequence, compared element by element. The escape is an encoding applied
// AFTER that order is fixed, and it is not order-preserving against an
// arbitrary neighbour: it inserts `\` (U+005C) where the delimiter stood,
// so `["a,b"]` and `["a0"]` compare `,` (U+002C) < `0` (U+0030) authored
// but `\` > `0` escaped.
//
// Sorting the atoms on the escaped rendering therefore emits a row's atoms
// in an order that is an artifact of the encoding rather than the canonical
// one REQ-87 names, and REQ-89's deterministic emission drifts with it.
// ADVERSARIAL
func TestCanonicalAtomOrderSurvivesADelimiterBearingLiteral(t *testing.T) {
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["start", "done"]
single_valued = true
required = true

[tags.mark]
provenance = "observed"
kind = "scalar"
`
	const preamble = `
terminal = ["fin"]

[initial]
status = "start"

[context.fin]
[context.fin.match.status]
eq = "done"
`
	// A row carrying two atoms that tie on key, block, AND operator, so
	// `compareAtoms` is decided by the literal alone — the only position
	// where the literal comparison is reachable, since a block is a table
	// keyed by tag and cannot hold two atoms for one key.
	//
	// `["a,b"]` sorts BEFORE `["a0"]` as an authored member sequence
	// (`,` U+002C < `0` U+0030) and AFTER it once escaped (`\` U+005C >
	// `0`), so this pair is exactly the transposition the escape gets
	// wrong. RDR 0002 `0002:C19` compares the literal element by element
	// over the AUTHORED members, and REQ-87 names that order as the one
	// the fingerprint must be in.
	row := table.Row{
		RuleID: "two-atoms",
		Atoms: []table.Atom{
			{Key: "mark", Block: table.BlockAll, Operator: "in", Literal: []string{"a0"}},
			{Key: "mark", Block: table.BlockAll, Operator: "in", Literal: []string{"a,b"}},
		},
	}

	// PRECONDITION: the escape really does transpose this pair, so a sort
	// over the encoding and a sort over the authored members disagree.
	// Without that the fixture would pass for the wrong reason.
	if !(`a\,b` > `a0`) {
		t.Fatalf("fixture no longer exercises the transposition: escaped " +
			"`a\\,b` must sort after `a0` for this pair to discriminate")
	}

	fp := graphlint.Fingerprint(row)

	iDelim := strings.Index(fp, `a\,b`)
	iDigit := strings.Index(fp, `a0`)
	if iDelim < 0 || iDigit < 0 {
		t.Fatalf("the fingerprint %q lost one of the two literals; it must "+
			"carry both to order them at all", fp)
	}

	// The decisive assertion: `["a,b"]` precedes `["a0"]` in the rendering
	// because that is their AUTHORED order, even though the escaped forms
	// sort the other way.
	if iDelim > iDigit {
		t.Errorf("the row's atoms were ordered by the ESCAPED rendering "+
			"rather than by RDR 0002's canonical order over authored "+
			"members: in %q the atom carrying `a,b` must precede the one "+
			"carrying `a0`, since `,` (U+002C) < `0` (U+0030) as authored. "+
			"Sorting on the encoding orders a row's atoms by an artifact "+
			"of the escape, which is not the canonical sortable "+
			"serialization REQ-87 requires.", fp)
	}
}

// REQ-106 / REQ-108: the composite key must be injective over member
// sequence CARDINALITY, not only over member content.
//
// An infix join renders `[]` and `[""]` identically, so a node holding a
// key with no members and one holding the empty string collapse to the same
// `Node.key` and the same `presenceKey`. Both are authorable — a `set` tag
// declaring no `elements` constrains no member, and `conform` loops zero
// times over an empty sequence — and they are NOT the same owned-state:
// `[]` satisfies no value atom while `[""]` satisfies `eq ""`. Collapsing
// them lets the `seen` maps and the fixpoint equality treat a reachable
// owned-state as already visited, which is the false-green direction.
// ADVERSARIAL
func TestNodeKeyIsInjectiveOverAnEmptyVersusEmptyStringMemberSequence(t *testing.T) {
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["start", "empty", "blank"]
single_valued = true
required = true

[tags.freeset]
provenance = "owned"
kind = "set"
`
	const preamble = `
terminal = ["fin"]

[initial]
status = "start"

[context.fin]
[context.fin.match.status]
eq = "start"
`
	// Two rows leaving the root with the SAME key held but different
	// member sequences: the empty sequence and the sequence holding "".
	m := mustLoad(t, source(decls, preamble+`
[[rule]]
id = "writes-empty"
[rule.match.status]
eq = "start"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "empty"
freeset = []

[[rule]]
id = "writes-blank"
[rule.match.status]
eq = "start"
[rule.match.recognized]
eq = "stop"
[rule.write]
status = "blank"
freeset = [""]
`))

	// PRECONDITION: the loader really accepted both member sequences. A
	// `set` with no declared `elements` constrains no member, so the pair
	// is reachable from conforming source rather than hypothetical.
	empty := rowByRuleID(t, m, "writes-empty")
	blank := rowByRuleID(t, m, "writes-blank")
	if got := valueForKey(t, empty, "freeset"); len(got) != 0 {
		t.Fatalf("row `writes-empty` did not carry an EMPTY member "+
			"sequence for `freeset`; got %#v", got)
	}
	if got := valueForKey(t, blank, "freeset"); !slices.Equal(got, []string{""}) {
		t.Fatalf("row `writes-blank` did not carry the one-element member "+
			"sequence [\"\"] for `freeset`; got %#v", got)
	}

	// Compare the two rows reduced to the `freeset` write ALONE. The rows
	// necessarily differ elsewhere (each needs its own outcome binding and
	// status write to be authorable at all), and that difference would mask
	// a `freeset` collision in a whole-row comparison.
	onlyFreeset := func(r table.Row) table.Row {
		return table.Row{
			RuleID:   r.RuleID,
			NextTags: []table.TagValue{{Key: "freeset", Value: valueForKey(t, r, "freeset")}},
		}
	}

	// REQ-114 closes the finding identity on the fingerprint, so a
	// collision here gives two distinct rows one identity.
	if fpEmpty, fpBlank := graphlint.Fingerprint(onlyFreeset(empty)), graphlint.Fingerprint(onlyFreeset(blank)); fpEmpty == fpBlank {
		t.Errorf("the empty member sequence and [\"\"] render the same "+
			"fingerprint %q. An infix join cannot tell them apart, so two "+
			"different owned-states share one identity — the false-green "+
			"direction REQ-106 forbids.", fpEmpty)
	}

	// `Reach` itself abstracts an unconstrained `set` to `OpaqueValue`
	// (a set declaring no `elements` carries no finite domain), so the
	// traversal never renders these two member sequences. The fingerprint
	// does: it closes over the row's raw `NextTags`, with no such
	// abstraction between the authored value and the rendering. That is
	// the path this pins.
}

// valueForKey returns the member sequence the row writes for key.
func valueForKey(t *testing.T, row table.Row, key string) []string {
	t.Helper()
	for _, tv := range row.NextTags {
		if tv.Key == key {
			return tv.Value
		}
	}
	t.Fatalf("row %s writes no tag %q", row.RuleID, key)
	return nil
}
