package cli

// RDR 0030 — the flow-, dump- and export-facing Testing Strategy scenarios
// (`0030:S2`, S7–S10) and C3's invariance, end to end through the command
// tree. The oracle throughout is the UNROLLED twin: a stepped ladder is
// correct exactly when every consumer reads it as the literal ladder.

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/table"
)

// c30Plan is one `flow resolve --plan-only` answer: a plan, or a refusal code.
type c30Plan struct {
	code string
	rule string
	data map[string]any
}

func c30Resolve(t *testing.T, model, art, outcome string) c30Plan {
	t.Helper()

	stdout, _, err := runCmd(t, "flow", "resolve", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art), "--outcome", outcome, "--plan-only", "--as=json")
	if err != nil {
		return c30Plan{code: clierr.ErrorCode(err)}
	}
	data := flowData(t, stdout)
	rule, _ := data["rule"].(string)
	return c30Plan{rule: rule, data: data}
}

// c30Cell is one (state, outcome) cell of the sweep.
type c30Cell struct {
	status, attempt, tier, outcome string
}

func (c c30Cell) String() string {
	return "status=" + c.status + " attempt=" + c.attempt + " tier=" + c.tier + " outcome=" + c.outcome
}

// c30Sweep resolves every (state, outcome) cell of the FULL cross product of
// the declared domains — status × attempt × tier × outcomes — over both
// models, seeding each state through the literal model's own write path.
// It returns the per-cell plans, literal then stepped.
func c30Sweep(t *testing.T, litSrc, stepSrc string) (cells []c30Cell, lit, step []c30Plan) {
	t.Helper()

	m := c30MustLoad(t, "the literal twin", litSrc)
	litPath, stepPath := writeFlowModel(t, litSrc), writeFlowModel(t, stepSrc)
	decl := m.Tags["attempt"]
	for _, status := range m.Tags["status"].Domain {
		for a := *decl.Min; a <= *decl.Max; a++ {
			for _, tier := range m.Tags["tier"].Domain {
				art := seedArtifact(t, litPath, "status="+status, "attempt="+itoa(a), "tier="+tier)
				for _, outcome := range m.Outcomes {
					cells = append(cells, c30Cell{status, itoa(a), tier, outcome})
					lit = append(lit, c30Resolve(t, litPath, art, outcome))
					step = append(step, c30Resolve(t, stepPath, art, outcome))
				}
			}
		}
	}
	return cells, lit, step
}

// c30RequireSweepAgrees is S2 / MVV item 3: `writes`, `next` and `clear`
// equal cell by cell, and a refusing cell refuses with the same code.
func c30RequireSweepAgrees(t *testing.T, litSrc, stepSrc string) (cells []c30Cell, lit []c30Plan) {
	t.Helper()

	cells, lit, step := c30Sweep(t, litSrc, stepSrc)
	var diffs int
	for i, c := range cells {
		l, s := lit[i], step[i]
		switch {
		case l.code != s.code:
			diffs++
			if diffs <= 5 {
				t.Errorf("%s: literal answers %q, stepped answers %q", c, c30Or(l.code, "a plan"), c30Or(s.code, "a plan"))
			}
		case l.code == "":
			for _, key := range []string{"writes", "next", "clear"} {
				if !reflect.DeepEqual(l.data[key], s.data[key]) {
					diffs++
					if diffs <= 5 {
						t.Errorf("%s: `%s` differs\n  literal: %s\n  stepped: %s", c, key,
							c30JSON(l.data[key]), c30JSON(s.data[key]))
					}
				}
			}
		}
	}
	if diffs > 5 {
		t.Errorf("… %d cell differences in all, over %d cells", diffs, len(cells))
	}
	return cells, lit
}

func c30Or(s, alt string) string {
	if s == "" {
		return alt
	}
	return s
}

// REQ-49: "`flow resolve --as json --plan-only` swept over every (state,
// outcome) cell of the pair, the state driven by rewriting the bound
// artifact file. The sweep's domain is the FULL cross product of both
// stepped tags' declared domains × the outcome set" … "**Expected**:
// `writes`, `next` and `clear` equal cell by cell." … "2b. **Scenario**:
// unclaimed-cell parity. The F4 shape — a cap cell excluded by the step rule
// and claimed by no other row — authored both ways, swept as in 2.
// **Expected**: BOTH sides refuse at that cell, with the same code."
// HAPPY PATH
func TestReq49_0030_ResolvePlansAgreeOverEveryCell(t *testing.T) {
	t.Run("2 — the declared pair", func(t *testing.T) {
		cells, lit := c30RequireSweepAgrees(t, c30Source(t, c30Literal), c30Source(t, c30Step))
		var planned int
		for _, p := range lit {
			if p.code == "" {
				planned++
			}
		}
		if planned == 0 || len(cells) != 2*10*3*3 {
			t.Errorf("the sweep covered %d cells with %d literal plans; want 180 cells, some planning",
				len(cells), planned)
		}
	})

	t.Run("2b — the unclaimed cap refuses on both sides", func(t *testing.T) {
		cells, lit := c30RequireSweepAgrees(t, c30Uncapped(t, c30Literal), c30Uncapped(t, c30Step))
		var refused int
		for i, c := range cells {
			if c.status == "open" && c.outcome == "retry" && c.attempt >= "5" && len(c.attempt) == 1 {
				if lit[i].code == "" {
					t.Errorf("%s: the literal twin plans an unclaimed cap cell", c)
				}
				refused++
			}
		}
		if refused == 0 {
			t.Error("the sweep reached no unclaimed cap cell")
		}
	})
}

// c30Emitting is the S7 fixture over either ladder: a declared `route` emit
// key with one disposition, authored on every retry row.
func c30Emitting(t *testing.T, name string) string {
	t.Helper()

	src := c30Edit(t, c30Source(t, name), "[context.open.match.status]",
		"[emit.route]\nkind = \"enum\"\n[emit.route.domain]\nagain = [\"retry\"]\n\n[context.open.match.status]")
	return c30EditAll(t, src, "[rule.write]\nattempt = ", "[rule.emit]\nroute = \"retry\"\n[rule.write]\nattempt = ")
}

// REQ-54: "`flow resolve` selecting a NON-FIRST expanded row of a stepped
// rule that authors an `emit` block" … "**Expected**: the payload's `emit`
// and `dispositions` are the authored block, and `rule` is the authored id"
// HAPPY PATH
func TestReq54_0030_ANonFirstExpandedRowJoinsTheAuthoredEmitBlock(t *testing.T) {
	litSrc, stepSrc := c30Emitting(t, c30Literal), c30Emitting(t, c30Step)
	litPath, stepPath := writeFlowModel(t, litSrc), writeFlowModel(t, stepSrc)
	art := seedArtifact(t, litPath, "status=open", "attempt=3", "tier=mid")

	lit := c30Resolve(t, litPath, art, "retry")
	if lit.code != "" || lit.rule != "retry-3" {
		t.Fatalf("the literal twin at attempt=3 answers %q/%q; want a plan from retry-3", lit.code, lit.rule)
	}
	step := c30Resolve(t, stepPath, art, "retry")
	if step.code != "" {
		t.Fatalf("the stepped ladder refuses the non-first cell attempt=3: %s", step.code)
	}
	if step.rule != "retry" {
		t.Errorf("`rule` = %q; want the authored id %q", step.rule, "retry")
	}
	if want := map[string]any{"route": "retry"}; !reflect.DeepEqual(step.data["emit"], want) {
		t.Errorf("`emit` = %s; want the authored block %s", c30JSON(step.data["emit"]), c30JSON(want))
	}
	if !reflect.DeepEqual(step.data["dispositions"], lit.data["dispositions"]) ||
		!reflect.DeepEqual(step.data["dispositions"], map[string]any{"route": "again"}) {
		t.Errorf("`dispositions` = %s; the literal twin's are %s", c30JSON(step.data["dispositions"]),
			c30JSON(lit.data["dispositions"]))
	}
}

// c30DumpLines returns the dump lines of rows expanded from rule.
func c30DumpLines(m *table.Model, prefix string) []string {
	var out []string
	for _, line := range strings.Split(table.Dump(m), "\n") {
		if strings.HasPrefix(line, "identity="+prefix) {
			out = append(out, line)
		}
	}
	return out
}

// c30DumpColumn returns one column of a dump line.
func c30DumpColumn(line, col string) string {
	cols := table.DumpColumns()
	start := strings.Index(" "+line, " "+col+"=")
	if start < 0 {
		return ""
	}
	rest := line[start+len(col)+1:]
	end := len(rest)
	for _, other := range cols {
		if i := strings.Index(rest, " "+other+"="); i >= 0 && i < end {
			end = i
		}
	}
	return rest[:end]
}

// REQ-55: "A step rule on `attempt` that ALSO authors `match.attempt in =
// [1,3]`, and a sibling fixture where that same `in` is inherited from a
// context rather than authored on the rule. **Expected**: both yield exactly
// |in ∩ admitted cells| rows — two, not the |members| × |cells| product —
// each carrying `match eq = <cell>` with NO surviving `in`, and a suffix
// carrying exactly one element for `attempt`. Asserted CLI-level on `dump`'s
// `atoms` column"
// HAPPY PATH
func TestReq55_0030_DumpPublishesTheSubsumedInRewrittenPerRow(t *testing.T) {
	step := c30Source(t, c30Step)
	authored := c30Edit(t, step, "[rule.guard.all.attempt]\ngte = 0\nlt = 5\n[rule.write]\nattempt = { step = 1 }",
		"[rule.match.attempt]\nin = [1, 3]\n[rule.guard.all.attempt]\ngte = 0\nlt = 5\n[rule.write]\nattempt = { step = 1 }")
	inherited := c30Edit(t, c30Edit(t, step, "[context.done.match.status]",
		"[context.odd.match.attempt]\nin = [1, 3]\n\n[context.done.match.status]"),
		"id = \"retry\"\nuse = [\"open\"]", "id = \"retry\"\nuse = [\"open\", \"odd\"]")

	for _, tc := range []struct{ name, src string }{{"authored", authored}, {"inherited", inherited}} {
		t.Run(tc.name, func(t *testing.T) {
			m := c30MustLoad(t, tc.name, tc.src)
			lines := c30DumpLines(m, "ladder.retry#")
			if len(lines) != 2 {
				t.Fatalf("the dump carries %d retry rows; want 2:\n%s", len(lines), strings.Join(lines, "\n"))
			}
			for _, line := range lines {
				id := c30DumpColumn(line, "identity")
				cell := strings.TrimPrefix(id, "ladder.retry#")
				if strings.Contains(cell, "#") || (cell != "1" && cell != "3") {
					t.Errorf("identity %q does not carry exactly one attempt element in {1,3}", id)
				}
				atoms := c30DumpColumn(line, "atoms")
				if !strings.Contains(atoms, "attempt.eq="+cell+"@match") {
					t.Errorf("%s: atoms %s carry no `match eq = %s`", id, atoms, cell)
				}
				if strings.Contains(atoms, "attempt.in=") {
					t.Errorf("%s: atoms %s retain an `in` on the stepped tag", id, atoms)
				}
			}
		})
	}
}

// REQ-56: "(a) publishes `retry#<cell>` in `dump`'s `identity` column and the
// graph document's row `identity` — the suffix is appended even at one cell;
// (b) publishes `retry#0 … retry#4` on those same surfaces; (c) mints NO
// suffix element" … "(d) a rule carrying a step on `attempt` AND an
// unsubsumed `match in` on `mode`, which publishes `retry#<attempt>#<mode>` —
// the suffix elements ordered by KEY, not by point kind."
// HAPPY PATH
func TestReq56_0030_RowSurfacesPublishTheSuffixedIdentity(t *testing.T) {
	step := c30Source(t, c30Step)
	oneCell := c30Edit(t, step, "gte = 0\nlt = 5\n[rule.write]\nattempt = { step = 1 }",
		"gte = 4\nlt = 5\n[rule.write]\nattempt = { step = 1 }")
	noStep := c30Edit(t, c30Source(t, c30Literal), "id = \"succeed\"\nuse = [\"open\"]\n[rule.match.recognized]\neq = \"succeed\"",
		"id = \"succeed\"\nuse = [\"open\"]\n[rule.match.recognized]\nin = [\"succeed\"]")
	mode := c30EditAll(t, c30Edit(t, step, "[read.state]",
		"[tags.mode]\nprovenance = \"owned\"\nkind = \"enum\"\ndomain = [\"fast\", \"slow\"]\nsingle_valued = true\n\n[read.state]"),
		`keys = ["status", "attempt", "tier"]`, `keys = ["status", "attempt", "tier", "mode"]`)
	mode = c30Edit(t, mode, "gte = 0\nlt = 5\n[rule.write]\nattempt = { step = 1 }",
		"gte = 0\nlt = 1\n[rule.write]\nattempt = { step = 1 }")
	mode = c30Edit(t, mode, "id = \"retry\"\nuse = [\"open\"]\n[rule.match.recognized]\neq = \"retry\"\n",
		"id = \"retry\"\nuse = [\"open\"]\n[rule.match.recognized]\neq = \"retry\"\n[rule.match.mode]\nin = [\"fast\", \"slow\"]\n")

	for _, tc := range []struct {
		name, src, prefix string
		want              []string
	}{
		{"a — one admitted cell", oneCell, "ladder.retry", []string{"ladder.retry#4"}},
		{"b — the multi-cell rule", step, "ladder.retry",
			[]string{"ladder.retry#0", "ladder.retry#1", "ladder.retry#2", "ladder.retry#3", "ladder.retry#4"}},
		{"c — a single-member in, no step", noStep, "ladder.succeed", []string{"ladder.succeed"}},
		{"d — a step and an unsubsumed in", mode, "ladder.retry", []string{"ladder.retry#0#fast", "ladder.retry#0#slow"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := c30MustLoad(t, tc.name, tc.src)
			var dump []string
			for _, line := range c30DumpLines(m, tc.prefix) {
				id := c30DumpColumn(line, "identity")
				if id == tc.prefix || strings.HasPrefix(id, tc.prefix+"#") {
					dump = append(dump, id)
				}
			}
			slices.Sort(dump)
			if !slices.Equal(dump, tc.want) {
				t.Errorf("dump identities = %v; want %v", dump, tc.want)
			}
			doc := c30RequireStrictGraph(t, writeModel(t, tc.src))
			var graph []string
			for _, r := range doc.Rows {
				if r.Identity == tc.prefix || strings.HasPrefix(r.Identity, tc.prefix+"#") {
					graph = append(graph, r.Identity)
				}
			}
			slices.Sort(graph)
			if !slices.Equal(graph, tc.want) {
				t.Errorf("graph row identities = %v; want %v", graph, tc.want)
			}
		})
	}
}

// REQ-57: "`graph` export and `dump` of `ladder-step.toml`, read for each
// expanded row's successor. **Expected**: every row's `next` (graph) and next
// column (`dump`) carries THAT row's stepped value, not the authored one."
// ADVERSARIAL
func TestReq57_0030_BothCarriersPublishEachRowsSteppedSuccessor(t *testing.T) {
	want := map[string][2]string{
		"ladder.retry#0": {"attempt", "1"}, "ladder.retry#1": {"attempt", "2"},
		"ladder.retry#2": {"attempt", "3"}, "ladder.retry#3": {"attempt", "4"},
		"ladder.retry#4": {"attempt", "5"}, "ladder.escalate#small": {"tier", "mid"},
		"ladder.escalate#mid": {"tier", "large"},
	}
	doc := c30RequireStrictGraph(t, c30Path(t, c30Step))
	seen := map[string]bool{}
	for _, r := range doc.Rows {
		w, ok := want[r.Identity]
		if !ok {
			continue
		}
		seen[r.Identity] = true
		for _, carrier := range []struct {
			name string
			vals []graphTagValue
		}{{"next", r.Next}, {"writes", r.Writes}} {
			var got []string
			for _, tv := range carrier.vals {
				if tv.Key == w[0] {
					got = tv.Value
				}
			}
			if !slices.Equal(got, []string{w[1]}) {
				t.Errorf("graph %s `%s` %s = %v; want [%s]", r.Identity, carrier.name, w[0], got, w[1])
			}
		}
	}
	if len(seen) != len(want) {
		t.Errorf("the export carries %d of the %d expanded rows", len(seen), len(want))
	}

	m := c30MustLoad(t, c30Step, c30Source(t, c30Step))
	for _, line := range c30DumpLines(m, "ladder.") {
		id := c30DumpColumn(line, "identity")
		w, ok := want[id]
		if !ok {
			continue
		}
		for _, col := range []string{"next", "writes"} {
			if !strings.Contains(c30DumpColumn(line, col), w[0]+"="+w[1]) {
				t.Errorf("dump %s `%s` = %s; want %s=%s", id, col, c30DumpColumn(line, col), w[0], w[1])
			}
		}
	}
}

// REQ-36: "Every surface that names a ROW publishes the suffixed identity
// `rule#cell` — the shape `in` expansion already emits there (`dump`'s
// `identity` column, the graph document's `identity` row field) — and every
// surface that names a RULE publishes the authored rule id (`flow next`'s and
// `flow resolve`'s `rule`, the graph document's edge `rule`)."
// HAPPY PATH
func TestReq36_0030_RowSurfacesNameRowsAndRuleSurfacesNameRules(t *testing.T) {
	doc := c30RequireStrictGraph(t, c30Path(t, c30Step))
	var rowSuffixed, edges int
	for _, r := range doc.Rows {
		if strings.HasPrefix(r.Identity, "ladder.retry#") {
			rowSuffixed++
		}
	}
	for _, e := range doc.Reach.Edges {
		if strings.Contains(e.Rule, "#") {
			t.Errorf("graph edge %s -> %s names the row %q; an edge names the authored rule", e.From, e.To, e.Rule)
		}
		if e.Rule == "retry" {
			edges++
		}
	}
	if rowSuffixed != 5 || edges == 0 {
		t.Errorf("graph: %d suffixed retry rows (want 5), %d edges attributed to `retry` (want some)",
			rowSuffixed, edges)
	}

	stepPath := writeFlowModel(t, c30Source(t, c30Step))
	art := seedArtifact(t, writeFlowModel(t, c30Source(t, c30Literal)), "status=open", "attempt=3", "tier=mid")
	if p := c30Resolve(t, stepPath, art, "retry"); p.code != "" || p.rule != "retry" {
		t.Errorf("flow resolve answers %q/%q; want a plan whose `rule` is %q", p.code, p.rule, "retry")
	}
	stdout, _, err := runCmd(t, "flow", "next", "--model", stepPath,
		"--artifact", artifactBinding(flowStateRole, art), "--as=json")
	if err != nil {
		t.Fatalf("flow next refused: %v\n%s", err, stdout)
	}
	cands, _ := objectsAt(flowData(t, stdout), "candidates")
	var rules []string
	for _, c := range cands {
		r, _ := c["rule"].(string)
		rules = append(rules, r)
	}
	if !slices.Contains(rules, "retry") || slices.ContainsFunc(rules, func(r string) bool { return strings.Contains(r, "#") }) {
		t.Errorf("flow next candidate rules = %v; want the authored ids, `retry` among them", rules)
	}
}

// c30NextShape reduces a `flow next` payload to its rule-free content: each
// candidate's outcome and plan, sorted.
func c30NextShape(t *testing.T, model, art string) []string {
	t.Helper()

	stdout, _, err := runCmd(t, "flow", "next", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art), "--as=json")
	if err != nil {
		return []string{"refused:" + clierr.ErrorCode(err)}
	}
	cands, _ := objectsAt(flowData(t, stdout), "candidates")
	var out []string
	for _, c := range cands {
		delete(c, "rule")
		out = append(out, c30JSON(c))
	}
	slices.Sort(out)
	return out
}

// REQ-33: "INVARIANCE. `flow next`, `flow resolve`, `flow set-state`, `lint`,
// `dump`, and `graph` read expanded rows through `Row.Writes`, `Row.NextTags`
// and `Row.Atoms` alone … and none learns a computed value shape; the graph
// document's `intrastate.graph/1` vocabulary gains no member."
// DOMAIN EDGE
func TestReq33_0030_NoConsumerLearnsAComputedValueShape(t *testing.T) {
	litPath, stepPath := writeFlowModel(t, c30Source(t, c30Literal)), writeFlowModel(t, c30Source(t, c30Step))

	// flow set-state takes the same writes on the stepped model.
	stepArt := seedArtifact(t, stepPath, "status=open", "attempt=2", "tier=small")

	for _, state := range [][]string{
		{"status=open", "attempt=0", "tier=small"},
		{"status=open", "attempt=4", "tier=mid"},
		{"status=open", "attempt=7", "tier=large"},
		{"status=done", "attempt=5", "tier=large"},
	} {
		art := seedArtifact(t, litPath, state...)
		if l, s := c30NextShape(t, litPath, art), c30NextShape(t, stepPath, art); !slices.Equal(l, s) {
			t.Errorf("%v: flow next differs\n  literal: %v\n  stepped: %v", state, l, s)
		}
	}
	if got := c30NextShape(t, stepPath, stepArt); len(got) == 0 {
		t.Error("flow next over the stepped model's own artifact offers no candidate")
	}

	// graph: the reach relation is the literal ladder's, bar edge rule names.
	lit, step := c30RequireStrictGraph(t, litPath), c30RequireStrictGraph(t, stepPath)
	pairs := func(d graphDoc) []string {
		var out []string
		for _, e := range d.Reach.Edges {
			out = append(out, e.From+" -> "+e.To)
		}
		slices.Sort(out)
		return slices.Compact(out)
	}
	if !slices.Equal(pairs(lit), pairs(step)) {
		t.Errorf("reach edges differ\n  literal: %v\n  stepped: %v", pairs(lit), pairs(step))
	}
	if !reflect.DeepEqual(lit.Reach.Nodes, step.Reach.Nodes) {
		t.Errorf("reach nodes differ\n  literal: %s\n  stepped: %s", c30JSON(lit.Reach.Nodes), c30JSON(step.Reach.Nodes))
	}
}
