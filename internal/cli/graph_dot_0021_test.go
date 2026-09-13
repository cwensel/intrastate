package cli

// RDR 0021 — the DOT document (`0021:C2`'s DOT paragraph, `0021:S6`).
//
// The assertion boundary is fixed by C2: the node/edge SET and the
// abstraction marker are NORMATIVE; styling and attributes are NOT. Every
// oracle below is therefore set equality or marker presence, never a
// rendered-attribute comparison.
//
// S6 is explicit that `dot -Tsvg` exiting 0 is NOT sufficient: a
// well-formed but MIS-ESCAPED identifier still renders. So each identifier
// is unescaped by inverting A6's recorded escaping ORDER and compared to
// its source id.

import (
	"os/exec"
	"slices"
	"strings"
	"testing"
)

// exportDOT runs the verb with `--emit dot` and returns the bare DOT
// document stdout carries.
func exportDOT(t *testing.T, modelSrc string) string {
	t.Helper()

	requireGraphVerb(t)
	path := writeModel(t, modelSrc)
	stdout, _, err := runGraph(t, "--model", path, "--emit", "dot")
	if err != nil {
		t.Fatalf("--emit dot refused: %v", err)
	}
	return stdout
}

// REQ-39: "The DOT document renders the same value: one node per
// reachability node, one edge per reachability edge labeled with its rule
// id, the initial node marked, and the abstraction marker rendered in the
// graph header comment/label so the diagram carries it too"
// REQ-40: "its node/edge SET and the marker are normative, its
// styling/attributes are not."
// REQ-70: "the DOT node/edge id set equals step 2's `reach` block."
// HAPPY PATH
func TestReq39And40And70_TheDOTNodeAndEdgeSetEqualsTheReachBlock(t *testing.T) {
	requireGraphVerb(t)

	doc := decodeDocument(t, exportDocument(t, legalModel))
	if doc.Reach == nil {
		t.Fatal("the JSON document carries no `reach` block")
	}
	dot := exportDOT(t, legalModel)

	// The NODE set, compared as sets. Styling is non-normative, so the
	// comparison is over identifiers only.
	wantNodes := make([]string, 0, len(doc.Reach.Nodes))
	for _, n := range doc.Reach.Nodes {
		wantNodes = append(wantNodes, n.ID)
	}
	gotNodes := dotNodeIDs(t, dot)
	assertSameSet(t, "node", wantNodes, gotNodes)

	// The EDGE set, as (from, to) pairs each labelled with its rule id.
	wantEdges := make([]string, 0, len(doc.Reach.Edges))
	for _, e := range doc.Reach.Edges {
		wantEdges = append(wantEdges, e.From+" -> "+e.To+" ["+e.Rule+"]")
	}
	gotEdges := dotEdges(t, dot)
	assertSameSet(t, "edge", wantEdges, gotEdges)
}

// REQ-29 / REQ-39 (the marker in the DOT arm).
// HAPPY PATH — the diagram carries the abstraction marker too, so a
// reviewer reading only the rendered graph still learns the relation is the
// declared over-approximation (premortem P-7).
func TestReq39_TheDOTDocumentCarriesTheAbstractionMarker(t *testing.T) {
	dot := exportDOT(t, legalModel)

	if !strings.Contains(dot, abstractionMarker) {
		t.Errorf("the DOT document does not carry the abstraction marker "+
			"%q. C2 renders it in the graph header comment/label so the "+
			"DIAGRAM carries it too — a reviewer reading only the picture "+
			"must still learn this is the declared over-approximation, not "+
			"the runtime relation\n%s", abstractionMarker, dot)
	}
}

// REQ-39 (the initial node is marked).
// HAPPY PATH
func TestReq39_TheDOTDocumentMarksTheInitialNode(t *testing.T) {
	requireGraphVerb(t)

	doc := decodeDocument(t, exportDocument(t, legalModel))
	dot := exportDOT(t, legalModel)

	// The marking mechanism is styling and therefore non-normative; what
	// IS normative is that the initial node is distinguished. The oracle
	// is that the DOT carries some statement naming the initial node
	// beyond its bare node statement.
	if doc.Reach == nil || len(doc.Reach.Nodes) == 0 {
		t.Fatal("the JSON document carries no reach nodes")
	}
	if !strings.Contains(strings.ToLower(dot), "initial") &&
		!strings.Contains(dot, "__start") {
		t.Errorf("the DOT document carries no marking of the initial node; "+
			"C2 requires \"the initial node marked\". The MECHANISM is "+
			"styling and free; the marking itself is not\n%s", dot)
	}
}

// REQ-41: "A pure renderer over the export value; equality test against the
// JSON arm's node/edge set."
// REQ-79: "DOT arm equality and hostile content — tag values carrying
// quotes, newlines, and non-ASCII. Expected: DOT node/edge id set equals
// the JSON document's `reach` block, the abstraction marker is present in
// the header, and every fixture survives `dot -Tsvg`"
// ADVERSARIAL — hostile content is where a renderer coupled to analysis
// internals, or one escaping in the wrong order, diverges.
func TestReq41And79_TheDOTArmSurvivesHostileContentAndStillMatches(t *testing.T) {
	requireGraphVerb(t)

	doc := decodeDocument(t, exportDocument(t, hostileModel))
	if doc.Reach == nil {
		t.Fatal("the JSON document carries no `reach` block")
	}
	dot := exportDOT(t, hostileModel)

	wantNodes := make([]string, 0, len(doc.Reach.Nodes))
	for _, n := range doc.Reach.Nodes {
		wantNodes = append(wantNodes, n.ID)
	}
	assertSameSet(t, "node", wantNodes, dotNodeIDs(t, dot))

	if !strings.Contains(dot, abstractionMarker) {
		t.Errorf("the hostile-content DOT carries no abstraction marker")
	}

	// The rendering must survive graphviz. This is necessary, not
	// sufficient — REQ-80 below supplies the discriminating half.
	requireGraphviz(t)
	cmd := exec.Command("dot", "-Tsvg")
	cmd.Stdin = strings.NewReader(dot)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Errorf("`dot -Tsvg` refused the emitted document: %v\nstderr: %s\n"+
			"--- dot ---\n%s", err, stderr.String(), dot)
	}
}

// REQ-80: "Exit 0 is NOT sufficient: each hostile node and edge
// identifier's emitted quoted string is unescaped by inverting A6's
// recorded escaping ORDER (undo the `\n` line-break escape, then `\"`, then
// `\\`) and asserted EQUAL to its source id, so a well-formed but
// mis-escaped identifier fails — the bug class A6's spike hit. Identifiers
// only; DOT label STYLING stays non-normative (F3)."
// ADVERSARIAL — this is the oracle a `dot -Tsvg` exit code cannot supply.
func TestReq80_EachDOTIdentifierUnescapesBackToItsSourceID(t *testing.T) {
	requireGraphVerb(t)

	doc := decodeDocument(t, exportDocument(t, hostileModel))
	if doc.Reach == nil || len(doc.Reach.Nodes) == 0 {
		t.Fatal("the JSON document carries no reach nodes")
	}
	dot := exportDOT(t, hostileModel)

	sourceIDs := map[string]bool{}
	for _, n := range doc.Reach.Nodes {
		sourceIDs[n.ID] = true
	}

	quoted := dotQuotedStrings(dot)
	if len(quoted) == 0 {
		t.Fatalf("the DOT document carries no quoted strings, so the "+
			"unescape oracle is unexercised\n%s", dot)
	}

	var matched int
	for _, q := range quoted {
		unescaped := dotUnescape(q)
		if sourceIDs[unescaped] {
			matched++
		}
	}
	if matched < len(sourceIDs) {
		var got []string
		for _, q := range quoted {
			got = append(got, dotUnescape(q))
		}
		t.Errorf("only %d of %d source node ids round-trip through the DOT "+
			"quoting. Each identifier's emitted quoted string must unescape "+
			"— inverting A6's order: the `\\n` line break, then `\\\"`, then "+
			"`\\\\` — back to its SOURCE id. A well-formed but mis-escaped "+
			"identifier still renders at exit 0, which is why this oracle "+
			"exists and the exit code does not suffice.\n"+
			"source ids: %v\nunescaped: %v\n--- dot ---\n%s",
			matched, len(sourceIDs), dotSourceIDs(sourceIDs), got, dot)
	}
}

// dotUnescape inverts A6's recorded escaping ORDER: undo the `\n` line-break
// escape, then `\"`, then `\\`. The order matters — undoing `\\` first
// would turn the two bytes of an escaped backslash into an escape for
// whatever follows.
func dotUnescape(s string) string {
	s = strings.ReplaceAll(s, `\n`, "\n")
	s = strings.ReplaceAll(s, `\"`, `"`)
	s = strings.ReplaceAll(s, `\\`, `\`)
	return s
}

// dotQuotedStrings returns every double-quoted string in the document, with
// its surrounding quotes stripped and its escapes left intact.
func dotQuotedStrings(dot string) []string {
	var out []string
	var cur strings.Builder
	var inQuote, escaped bool
	for _, r := range dot {
		switch {
		case escaped:
			cur.WriteRune(r)
			escaped = false
		case r == '\\' && inQuote:
			cur.WriteRune(r)
			escaped = true
		case r == '"':
			if inQuote {
				out = append(out, cur.String())
				cur.Reset()
			}
			inQuote = !inQuote
		case inQuote:
			cur.WriteRune(r)
		}
	}
	return out
}

// dotNodeIDs returns the identifiers of the DOT node statements — a quoted
// identifier that is not part of an edge statement.
func dotNodeIDs(t *testing.T, dot string) []string {
	t.Helper()

	var out []string
	for line := range strings.Lines(dot) {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") ||
			strings.HasPrefix(trimmed, "#") {
			continue
		}
		if strings.Contains(trimmed, "->") {
			continue
		}
		quoted := dotQuotedStrings(trimmed)
		if len(quoted) == 0 {
			continue
		}
		// A node statement opens with its identifier; anything after is an
		// attribute value and therefore styling.
		if strings.HasPrefix(trimmed, `"`) {
			out = append(out, dotUnescape(quoted[0]))
		}
	}
	return out
}

// dotEdges returns each edge statement as "from -> to [rule]", with the
// rule read from the edge's label attribute. Styling beyond the label is
// ignored: C2 makes the SET and the labelling normative, not the attributes.
func dotEdges(t *testing.T, dot string) []string {
	t.Helper()

	var out []string
	for line := range strings.Lines(dot) {
		trimmed := strings.TrimSpace(line)
		if !strings.Contains(trimmed, "->") {
			continue
		}
		quoted := dotQuotedStrings(trimmed)
		if len(quoted) < 2 {
			continue
		}
		from, to := dotUnescape(quoted[0]), dotUnescape(quoted[1])
		rule := ""
		if len(quoted) > 2 {
			rule = dotUnescape(quoted[2])
		}
		out = append(out, from+" -> "+to+" ["+rule+"]")
	}
	return out
}

// assertSameSet compares two identifier sets, reporting the symmetric
// difference. Order is not asserted: C2 makes the SET normative.
func assertSameSet(t *testing.T, what string, want, got []string) {
	t.Helper()

	wantSorted := slices.Clone(want)
	gotSorted := slices.Clone(got)
	slices.Sort(wantSorted)
	slices.Sort(gotSorted)
	wantSorted = slices.Compact(wantSorted)
	gotSorted = slices.Compact(gotSorted)

	if slices.Equal(wantSorted, gotSorted) {
		return
	}
	var missing, extra []string
	for _, w := range wantSorted {
		if !slices.Contains(gotSorted, w) {
			missing = append(missing, w)
		}
	}
	for _, g := range gotSorted {
		if !slices.Contains(wantSorted, g) {
			extra = append(extra, g)
		}
	}
	t.Errorf("the DOT %s set does not equal the JSON document's `reach` "+
		"block.\nmissing from DOT: %v\nnot in the document: %v\n"+
		"C2 makes the node/edge SET normative — the DOT arm is a projection "+
		"of the SAME export value (A6), not a second traversal",
		what, missing, extra)
}

// dotSourceIDs renders a set's members for a failure message.
func dotSourceIDs(set map[string]bool) []string {
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	slices.Sort(out)
	return out
}

// requireGraphviz skips when `dot` is unavailable. The renderer's own
// correctness is asserted without graphviz by REQ-80's unescape oracle;
// this guard only governs the "survives `dot -Tsvg`" half.
func requireGraphviz(t *testing.T) {
	t.Helper()

	if _, err := exec.LookPath("dot"); err != nil {
		t.Skip("graphviz `dot` is not installed; S6's render check needs it")
	}
}

// hostileModel carries tag values with quotes, backslashes, newlines, and
// non-ASCII — the escaping surface S6 names, and the content that broke
// A6's spike renderer.
const hostileModel = `
outcomes = ["go", "stop"]
terminal = ["done"]

[model]
id = "hostile"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["say \"hi\"", "back\\slash", "line\nbreak", "ünïcødé"]
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
status = "say \"hi\""

[context.done]
[context.done.match.status]
eq = "ünïcødé"

[[rule]]
id = "to-backslash"
[rule.match.status]
eq = "say \"hi\""
[rule.match.recognized]
eq = "go"
[rule.write]
status = "back\\slash"

[[rule]]
id = "to-newline"
[rule.match.status]
eq = "back\\slash"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "line\nbreak"

[[rule]]
id = "to-unicode"
[rule.match.status]
eq = "line\nbreak"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "ünïcødé"

[[rule]]
id = "go-otherwise"
escape = ["no_match"]
[rule.match.recognized]
eq = "go"
`
