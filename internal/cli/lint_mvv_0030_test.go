package cli

// RDR 0030 — the Minimum Viable Validation (`0030:MVV`) and the lint-,
// load- and export-facing Testing Strategy scenarios (`0030:S1`, S3–S6,
// S11), end to end through the command tree.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/graphlint"
	"github.com/cwensel/intrastate/internal/guard"
	"github.com/cwensel/intrastate/internal/table"
)

// --- fixtures derived from the committed ladders -----------------------------

// c30S11 is the S11 attribution fixture over either ladder: `attempt` starts
// at 3, and the retry rows (one per cell, pinned by a match `eq`) read an
// owned `note` nobody initialises. Every transition out of the root writes
// `note`, so the ONLY reachable owned-state where a retry row reads it absent
// is the root — cell 3. Exactly one admitted cell carries the per-row
// blocking defect.
func c30S11(t *testing.T, stepped bool) string {
	t.Helper()

	name := c30Literal
	if stepped {
		name = c30Step
	}
	src := c30Source(t, name)
	src = c30Edit(t, src, "attempt = 0\ntier = \"small\"", "attempt = 3\ntier = \"small\"")
	src = c30Edit(t, src, "[read.state]", "[tags.note]\nprovenance = \"owned\"\nkind = \"enum\"\n"+
		"domain = [\"x\"]\nsingle_valued = true\n\n[read.state]")
	src = c30EditAll(t, src, `keys = ["status", "attempt", "tier"]`, `keys = ["status", "attempt", "tier", "note"]`)
	src = c30EditAll(t, src, "attempt = 0\n", "attempt = 0\nnote = \"x\"\n")
	if stepped {
		return c30Edit(t, src,
			"[rule.guard.all.attempt]\ngte = 0\nlt = 5\n[rule.write]\nattempt = { step = 1 }\nstatus = \"open\"\n",
			"[rule.match.attempt]\nin = [0, 1, 2, 3, 4]\n[rule.guard.all.attempt]\ngte = 0\nlt = 5\n"+
				"[rule.guard.all.note]\neq = \"x\"\n[rule.write]\nattempt = { step = 1 }\nstatus = \"open\"\nnote = \"x\"\n")
	}
	for n := range 5 {
		c, next := itoa(n), itoa(n+1)
		src = c30Edit(t, src,
			"[rule.guard.all.attempt]\neq = "+c+"\ngte = 0\nlt = 5\n[rule.write]\nattempt = "+next+"\nstatus = \"open\"\n",
			"[rule.match.attempt]\neq = "+c+"\n[rule.guard.all.attempt]\neq = "+c+"\ngte = 0\nlt = 5\n"+
				"[rule.guard.all.note]\neq = \"x\"\n[rule.write]\nattempt = "+next+"\nstatus = \"open\"\nnote = \"x\"\n")
	}
	return src
}

// c30Pin appends a literal row `pin` claiming cell 3 of the retry group, so
// the group's ordinary rows overlap there (S11 b/c).
func c30Pin(t *testing.T, name string) string {
	t.Helper()

	return c30Source(t, name) + `
[[rule]]
id = "pin"
use = ["open"]
[rule.match.recognized]
eq = "retry"
[rule.guard.all.attempt]
eq = 3
[rule.write]
status = "done"
`
}

// c30Uncapped removes `retry-cap`, leaving cells 5..9 of the retry outcome
// claimed by no row (the F4 shape).
func c30Uncapped(t *testing.T, name string) string {
	t.Helper()

	return c30Edit(t, c30Source(t, name), `[[rule]]
id = "retry-cap"
use = ["open"]
[rule.match.recognized]
eq = "retry"
[rule.guard.all.attempt]
gte = 5
[rule.write]
status = "done"
`, "")
}

// c30Dead is S5b: the stepped `retry` excludes the INTERIOR cell 2 with an
// `unless`, and a literal row in the same group claims 2. The literal twin
// carries the same `unless` on its `retry-2` row.
func c30Dead(t *testing.T, stepped bool) string {
	t.Helper()

	claim := `
[[rule]]
id = "retry-two"
use = ["open"]
[rule.match.recognized]
eq = "retry"
[rule.guard.all.attempt]
eq = 2
[rule.write]
attempt = 3
status = "open"
`
	if stepped {
		return c30Edit(t, c30Source(t, c30Step), "lt = 5\n[rule.write]\nattempt = { step = 1 }",
			"lt = 5\n[rule.guard.unless.attempt]\neq = 2\n[rule.write]\nattempt = { step = 1 }") + claim
	}
	return c30Edit(t, c30Source(t, c30Literal), "eq = 2\ngte = 0\nlt = 5\n[rule.write]",
		"eq = 2\ngte = 0\nlt = 5\n[rule.guard.unless.attempt]\neq = 2\n[rule.write]") + claim
}

// c30Split is S5c: the exclusion of cell 2 written as POSITIVE atoms — two
// step rules either side of it — with a literal row claiming 2.
func c30Split(t *testing.T, stepped bool) string {
	t.Helper()

	claim := `
[[rule]]
id = "retry-two"
use = ["open"]
[rule.match.recognized]
eq = "retry"
[rule.guard.all.attempt]
eq = 2
[rule.write]
attempt = 3
status = "open"
`
	if stepped {
		return c30Edit(t, c30Source(t, c30Step), `[[rule]]
id = "retry"
use = ["open"]
[rule.match.recognized]
eq = "retry"
[rule.guard.all.attempt]
gte = 0
lt = 5
[rule.write]
attempt = { step = 1 }
status = "open"
`, `[[rule]]
id = "retry-low"
use = ["open"]
[rule.match.recognized]
eq = "retry"
[rule.guard.all.attempt]
gte = 0
lt = 2
[rule.write]
attempt = { step = 1 }
status = "open"

[[rule]]
id = "retry-high"
use = ["open"]
[rule.match.recognized]
eq = "retry"
[rule.guard.all.attempt]
gt = 2
lt = 5
[rule.write]
attempt = { step = 1 }
status = "open"
`) + claim
	}
	src := c30Source(t, c30Literal)
	for n := range 5 {
		c := itoa(n)
		old := "id = \"retry-" + c + "\"\nuse = [\"open\"]\n[rule.match.recognized]\neq = \"retry\"\n" +
			"[rule.guard.all.attempt]\neq = " + c + "\ngte = 0\nlt = 5\n"
		switch {
		case n < 2:
			src = c30Edit(t, src, old, strings.Replace(old, "lt = 5", "lt = 2", 1))
		case n == 2:
			src = c30Edit(t, src, "[[rule]]\n"+old+"[rule.write]\nattempt = 3\nstatus = \"open\"\n", "")
		default:
			src = c30Edit(t, src, old, strings.Replace(old, "gte = 0", "gt = 2", 1))
		}
	}
	return src + claim
}

// --- REQ-MVV ------------------------------------------------------------------

// REQ-MVV (`0030:MVV`): "Author two in-repo fixtures of the Background's
// ladder …" … "2. `intrastate lint` on both: identical finding sets under
// S1's `(code, key, dimension, class, reason)` projection, zero
// `graph-overlap`, zero `graph-coverage-gap` (A5). Plus … a per-row finding
// on the stepped ladder names the offending CELL (`retry#3`)" … "3. `flow
// resolve` over every (state, outcome) cell of both: identical
// `writes`/`next`/`clear` (A6)." … "4. `ladder-step-unguarded.toml` …: load
// refuses `malformed_tag_declaration` naming `retry`, `attempt`, cell `9`,
// value `10`; the enum sibling names `tier`, `large`. Two more refusal
// fixtures …" … "5. The dead-row outcome and A5's negative control …" …
// "6. `graph` export of `ladder-step.toml` decodes under the shipped
// `intrastate.graph/1` document type with no new member (C3)." … "7. …
// the `[[rules]]` count of `ladder-step.toml` is ONE PER INTENT … and is
// strictly less than `ladder-literal.toml`'s." … "End-state: the step ladder
// is the literal ladder to every consumer, it is one row per intent in
// source, and the unguarded cap is a named load refusal."
// HAPPY PATH
func TestMVV_0030_TheStepLadderIsTheLiteralLadderToEveryConsumer(t *testing.T) {
	t.Run("item 1 — both fixtures load", func(t *testing.T) {
		lit := c30MustLoad(t, c30Literal, c30Source(t, c30Literal))
		step := c30MustLoad(t, c30Step, c30Source(t, c30Step))
		if len(step.Rows) != len(lit.Rows) {
			t.Errorf("the stepped ladder normalizes to %d rows, the literal one to %d",
				len(step.Rows), len(lit.Rows))
		}
		for _, key := range []string{"attempt", "tier"} {
			if step.Tags[key].Kind == "" {
				t.Errorf("the ladder declares no %s tag", key)
			}
		}
		if got := step.Tags["tier"].Domain; !slices.Equal(got, []string{"small", "mid", "large"}) {
			t.Errorf("tier's domain = %v; the MVV fixes [small mid large]", got)
		}
	})

	t.Run("item 2 — lint agrees, overlaps and gaps are zero, the cell is named", func(t *testing.T) {
		lit := c30RunLint(t, c30Path(t, c30Literal))
		step := c30RunLint(t, c30Path(t, c30Step))
		c30Loaded(t, c30Step, step)
		litSet, stepSet := c30ProjSet(lit.findings), c30ProjSet(step.findings)
		if len(litSet) == 0 {
			t.Fatal("the literal ladder lints no finding; the set comparison would be vacuous")
		}
		if !slices.Equal(litSet, stepSet) {
			t.Errorf("projected finding sets differ\n  literal: %v\n  step:    %v", litSet, stepSet)
		}
		for _, r := range []c30Lint{lit, step} {
			if n := c30Count(r.findings, graphlint.CodeOverlap) + c30Count(r.findings, graphlint.CodeCoverageGap); n != 0 {
				t.Errorf("%d overlap/coverage-gap findings; want zero:\n%s", n, r.stdout)
			}
		}
		c30RequireCellAttribution(t)
	})

	t.Run("item 3 — flow resolve agrees over every cell", func(t *testing.T) {
		c30RequireSweepAgrees(t, c30Source(t, c30Literal), c30Source(t, c30Step))
	})

	t.Run("item 4 — the unguarded cap is a named load refusal", func(t *testing.T) {
		f := c30LoadRefusal(t, c30Unguarded, c30Source(t, c30Unguarded), table.CatMalformedTagDeclaration)
		for _, n := range []string{"retry", "attempt", "9", "10"} {
			if !strings.Contains(f.Message, n) {
				t.Errorf("the refusal does not name %q: %s", n, f.Message)
			}
		}
		e := c30LoadRefusal(t, c30UnguardedEnum, c30Source(t, c30UnguardedEnum), table.CatMalformedTagDeclaration)
		for _, n := range []string{"tier", "large"} {
			if !strings.Contains(e.Message, n) {
				t.Errorf("the enum refusal does not name %q: %s", n, e.Message)
			}
		}
		c30LoadRefusal(t, "[initial] attempt = { step = 1 }",
			c30Edit(t, c30Source(t, c30Step), "attempt = 0\ntier", "attempt = { step = 1 }\ntier"),
			table.CatMalformedInitialDeclaration)
		c30LoadRefusal(t, "a predicate literal { step = 1 }",
			c30Edit(t, c30Source(t, c30Step), "gte = 5", "eq = { step = 1 }"),
			table.CatMalformedPredicateAtom)
	})

	t.Run("item 5 — the dead row is advisory and the positive exclusion does not overlap", func(t *testing.T) {
		c30RequireDeadRowAndNoOverlap(t)
	})

	t.Run("item 6 — the graph export decodes strictly", func(t *testing.T) {
		c30RequireStrictGraph(t, c30Path(t, c30Step))
	})

	t.Run("item 7 — one rule per intent, fewer than the literal ladder", func(t *testing.T) {
		var lit, step map[string]any
		if err := toml.Unmarshal([]byte(c30Source(t, c30Literal)), &lit); err != nil {
			t.Fatalf("parsing %s: %v", c30Literal, err)
		}
		if err := toml.Unmarshal([]byte(c30Source(t, c30Step)), &step); err != nil {
			t.Fatalf("parsing %s: %v", c30Step, err)
		}
		perIntent, stepCount := c30StepRules(t, step)
		_, litCount := c30StepRules(t, lit)
		want := []string{"attempt@retry", "tier@escalate"}
		var got []string
		for k, ids := range perIntent {
			got = append(got, k)
			if len(ids) != 1 {
				t.Errorf("intent %s is written as %d step rules %v; want one", k, len(ids), ids)
			}
		}
		slices.Sort(got)
		if !slices.Equal(got, want) {
			t.Errorf("step intents = %v; want %v", got, want)
		}
		if stepCount >= litCount {
			t.Errorf("ladder-step.toml carries %d rules, ladder-literal.toml %d; the step ladder must be "+
				"strictly smaller", stepCount, litCount)
		}
		// The count is only meaningful if the stepped source IS the literal
		// table: it loads.
		c30MustLoad(t, c30Step, c30Source(t, c30Step))
	})
}

// c30RequireCellAttribution is MVV item 2's second half and S11(a): the per-
// row blocking finding names the offending cell on the stepped side and the
// literal row on the twin.
func c30RequireCellAttribution(t *testing.T) {
	t.Helper()

	lit := c30LintSrc(t, c30S11(t, false))
	obw := c30WithCode(lit.findings, graphlint.CodeOwnedBeforeWrite)
	if len(obw) != 1 || obw[0].Rule != "retry-3" {
		t.Fatalf("the literal S11 twin must carry exactly one owned-before-write, on retry-3; got %s",
			c30JSON(obw))
	}

	step := c30LintSrc(t, c30S11(t, true))
	c30Loaded(t, "the stepped S11 fixture", step)
	got := c30WithCode(step.findings, graphlint.CodeOwnedBeforeWrite)
	if len(got) != 1 {
		t.Fatalf("the stepped S11 fixture carries %d owned-before-write findings; want exactly one: %s",
			len(got), c30JSON(got))
	}
	if got[0].Rule != "retry#3" {
		t.Errorf("the per-row finding's rule = %q; want %q — the offending cell, not the bare rule",
			got[0].Rule, "retry#3")
	}
}

// c30RequireDeadRowAndNoOverlap is MVV item 5 / S5b / S5c.
func c30RequireDeadRowAndNoOverlap(t *testing.T) {
	t.Helper()

	// The literal twins pin the shape today.
	for _, stepped := range []bool{false, true} {
		src := c30Dead(t, stepped)
		m := c30MustLoad(t, "S5b", src)
		for _, r := range m.Rows {
			if r.RuleID != "retry" && r.RuleID != "retry-2" {
				continue
			}
			if len(r.Suffix) == 1 && r.Suffix[0] != "2" {
				continue
			}
			// The projectability precondition: the dead row decides.
			acc := guard.AcceptedAssignments(m, r)
			if !acc.Projectable() || acc.Len() != 0 {
				t.Errorf("%s: the dead row must project whole to the empty accepted set (projectable %v, "+
					"len %d)", r.Identity(), acc.Projectable(), acc.Len())
			}
		}
		lint := c30LintSrc(t, src)
		c30Loaded(t, "S5b", lint)
		if n := c30Count(lint.findings, graphlint.CodeRedundantRow); n != 1 {
			t.Errorf("S5b (stepped=%v) lints %d graph-redundant-row; want exactly one:\n%s", stepped, n, lint.stdout)
		}

		split := c30Split(t, stepped)
		sm := c30MustLoad(t, "S5c", split)
		for _, r := range sm.Rows {
			if r.Outcome == "retry" && !guard.AcceptedAssignments(sm, r).Projectable() {
				t.Errorf("S5c: %s does not project; the overlap oracle would skip it", r.Identity())
			}
		}
		sl := c30LintSrc(t, split)
		c30Loaded(t, "S5c", sl)
		if n := c30Count(sl.findings, graphlint.CodeOverlap); n != 0 {
			t.Errorf("S5c (stepped=%v) lints %d graph-overlap; want zero:\n%s", stepped, n, sl.stdout)
		}
	}
}

// c30RequireStrictGraph decodes a model's `graph` export into the shipped
// document type with unknown fields rejected.
func c30RequireStrictGraph(t *testing.T, path string) graphDoc {
	t.Helper()

	stdout, _, err := runGraph(t, "--model", path)
	if err != nil {
		t.Fatalf("graph refused %s: %v\n%s", path, err, stdout)
	}
	dec := json.NewDecoder(bytes.NewReader([]byte(strings.TrimSpace(stdout))))
	dec.DisallowUnknownFields()
	var doc graphDoc
	if err := dec.Decode(&doc); err != nil {
		t.Fatalf("the export does not decode strictly into the shipped intrastate.graph/1 type: %v", err)
	}
	if doc.Schema != schemaMarkerValue {
		t.Errorf("schema = %q; want %q", doc.Schema, schemaMarkerValue)
	}
	return doc
}

// --- S1 -----------------------------------------------------------------------

// REQ-48: "identical finding sets compared as sets over the PROJECTION
// `(code, key, dimension, class, reason)` of each `findings[]` entry — not
// over whole finding objects." … "The compared set MUST be non-empty,
// asserted BEFORE the set comparison" … "(a) the CLEAN pair … must exit OK
// and yield a non-empty ADVISORY set, guaranteed by authoring a
// `graph-idempotent-write` on an UNSTEPPED key … mirrored on both sides; and
// (b) a BLOCKING pair, `ladder-literal-blocking.toml` /
// `ladder-step-blocking.toml` … The blocking defect is
// `graph-owned-before-write` on a tag neither ladder initialises" … "Zero
// `graph-overlap` and zero `graph-coverage-gap` on all four."
// HAPPY PATH
func TestReq48_0030_LintFindingSetsAgreeUnderTheProjection(t *testing.T) {
	for _, tc := range []struct {
		name, lit, step string
		blocking        bool
	}{
		{"clean pair", c30Literal, c30Step, false},
		{"blocking pair", c30LitBlocking, c30StepBlocking, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			lit := c30RunLint(t, c30Path(t, tc.lit))
			step := c30RunLint(t, c30Path(t, tc.step))
			c30Loaded(t, tc.step, step)

			for _, r := range []struct {
				name string
				run  c30Lint
			}{{tc.lit, lit}, {tc.step, step}} {
				if r.run.failed != tc.blocking {
					t.Errorf("%s: lint failed=%v; want %v:\n%s", r.name, r.run.failed, tc.blocking, r.run.stdout)
				}
				if n := c30Count(r.run.findings, graphlint.CodeOverlap) +
					c30Count(r.run.findings, graphlint.CodeCoverageGap); n != 0 {
					t.Errorf("%s: %d overlap/coverage-gap findings; want zero", r.name, n)
				}
			}
			if tc.blocking {
				for _, r := range []c30Lint{lit, step} {
					if c30Count(r.findings, graphlint.CodeOwnedBeforeWrite) == 0 {
						t.Errorf("the blocking pair carries no graph-owned-before-write:\n%s", r.stdout)
					}
				}
			} else if c30Count(lit.findings, graphlint.CodeIdempotentWrite) == 0 {
				t.Errorf("the clean literal ladder carries no graph-idempotent-write advisory:\n%s", lit.stdout)
			}

			litSet, stepSet := c30ProjSet(lit.findings), c30ProjSet(step.findings)
			if len(litSet) == 0 || len(stepSet) == 0 {
				t.Fatalf("an empty projected set (literal %d, step %d); the comparison would be vacuous",
					len(litSet), len(stepSet))
			}
			if !slices.Equal(litSet, stepSet) {
				t.Errorf("projected finding sets differ\n  literal: %v\n  step:    %v", litSet, stepSet)
			}
		})
	}
}

// --- S11 ----------------------------------------------------------------------

// REQ-58: "**Expected**: the finding's `rule` reads `retry#3` — the offending
// cell — and the literal twin's reads `retry-3`" … "Assert additionally that
// a GROUP-level finding (`coverage.go`'s `firstRuleID` sites) still reads the
// bare authored id" … "(b) a `graph-overlap` between a literal row and an
// expanded one: BOTH `rule` and `element` carry their own row's identity —
// the expanded side suffixed, the literal side not" … "(c) the same fixture
// authored so the group's `ambiguous_match` coverage arm is reachable: the
// arm must still fire."
// ADVERSARIAL
func TestReq58_0030_PerRowFindingsNameTheCellAndGroupFindingsStayBare(t *testing.T) {
	t.Run("a — the per-row finding names the cell", func(t *testing.T) {
		c30RequireCellAttribution(t)
	})

	t.Run("group-level findings read the bare authored id", func(t *testing.T) {
		src := c30Uncapped(t, c30Step)
		lit := c30LintSrc(t, c30Uncapped(t, c30Literal))
		if len(c30GroupFindings(lit.findings)) == 0 {
			t.Fatalf("the uncapped literal twin carries no group-level finding with a rule:\n%s", lit.stdout)
		}
		step := c30LintSrc(t, src)
		c30Loaded(t, "the uncapped step ladder", step)
		group := c30GroupFindings(step.findings)
		if len(group) == 0 {
			t.Fatalf("the uncapped step ladder carries no group-level finding with a rule:\n%s", step.stdout)
		}
		for _, f := range group {
			if strings.Contains(f.Rule, "#") {
				t.Errorf("a GROUP-level %s names a row identity %q; it names the group's bare rule id",
					f.Code, f.Rule)
			}
		}
	})

	t.Run("b and c — overlap slots and the ambiguous arm", func(t *testing.T) {
		c30RequireOverlapSlots(t)
	})
}

// c30GroupFindings are the findings `coverage.go`'s `firstRuleID` sites
// emit: a coverage gap, a product-too-large, or the no-participating-
// dimension unprovable-coverage finding, carrying a rule.
func c30GroupFindings(fs []clierr.Finding) []clierr.Finding {
	var out []clierr.Finding
	for _, f := range fs {
		if f.Rule == "" {
			continue
		}
		switch {
		case f.Code == graphlint.CodeCoverageGap, f.Code == graphlint.CodeProductTooLarge,
			f.Code == graphlint.CodeUnprovableCoverage && f.Reason == graphlint.ReasonNoParticipatingDimension:
			out = append(out, f)
		}
	}
	return out
}

// c30RequireOverlapSlots is S11(b)/(c): an overlap between the literal `pin`
// and the expanded cell 3 names each side by its own row identity, and the
// group's ambiguous_match coverage arm still fires.
func c30RequireOverlapSlots(t *testing.T) {
	t.Helper()

	slots := func(fs []clierr.Finding) [][]string {
		var out [][]string
		for _, f := range c30WithCode(fs, graphlint.CodeOverlap) {
			pair := []string{f.Rule, f.Element}
			slices.Sort(pair)
			out = append(out, pair)
		}
		return out
	}
	ambiguous := func(fs []clierr.Finding) bool {
		for _, f := range c30WithCode(fs, graphlint.CodeCoverageGap) {
			if f.Class == "ambiguous_match" {
				return true
			}
		}
		return false
	}

	lit := c30LintSrc(t, c30Pin(t, c30Literal))
	if got := slots(lit.findings); !slices.EqualFunc(got, [][]string{{"pin", "retry-3"}}, slices.Equal) {
		t.Fatalf("the literal twin's overlap slots = %v; want [[pin retry-3]]", got)
	}
	if !ambiguous(lit.findings) {
		t.Fatalf("the literal twin's ambiguous_match arm does not fire:\n%s", lit.stdout)
	}

	step := c30LintSrc(t, c30Pin(t, c30Step))
	c30Loaded(t, "the pinned step ladder", step)
	if got := slots(step.findings); !slices.EqualFunc(got, [][]string{{"pin", "retry#3"}}, slices.Equal) {
		t.Errorf("the overlap's rule/element slots = %v; want [[pin retry#3]] — each side names its "+
			"own row, the expanded one suffixed", got)
	}
	if !ambiguous(step.findings) {
		t.Errorf("the ambiguous_match coverage arm no longer fires on the stepped ladder — the "+
			"read-back join lost the suffixed overlap:\n%s", step.stdout)
	}
}

// REQ-37: "A per-row graph-lint finding on an EXPANDED row therefore
// publishes the suffixed identity `rule#cell`, not the bare authored id … That
// covers BOTH row-naming slots — `Rule`, and the `Element` slot the two
// findings that name a SECOND row populate" … "A finding may not suffix one
// slot and not the other" … "The group-level findings that name a group
// rather than a row … are unchanged". Scope (Q1 → reading (a)): every row
// carrying a suffix, `in`-minted rows included.
// DOMAIN EDGE
func TestReq37_0030_EveryRowCarryingASuffixPublishesItInBothRowSlots(t *testing.T) {
	// (i) A model that steps NOTHING and lints clean: `fan` expands by a
	// match `in`, and its small row writes tier back to small — a per-row
	// advisory on an `in`-minted row.
	clean := c30LintSrc(t, c30Fan(t, true, false))
	c30Loaded(t, "the in-expanded model", clean)
	if clean.failed {
		t.Fatalf("the in-expanded model must lint clean so its advisories publish:\n%s", clean.stdout)
	}
	var fanIdem []string
	for _, f := range c30WithCode(clean.findings, graphlint.CodeIdempotentWrite) {
		if strings.HasPrefix(f.Rule, "fan") {
			fanIdem = append(fanIdem, f.Rule)
		}
	}
	if !slices.Equal(fanIdem, []string{"fan#small"}) {
		t.Errorf("the in-expanded row's idempotent-write names %v; want [fan#small]", fanIdem)
	}

	// (ii) `pin` sits inside the small row: the redundant-row finding names
	// a SECOND row in its Element slot, and BOTH slots name their row.
	second := c30LintSrc(t, c30Fan(t, true, true))
	c30Loaded(t, "the pinned in-expanded model", second)
	var pairs [][]string
	for _, f := range c30WithCode(second.findings, graphlint.CodeRedundantRow) {
		pair := []string{f.Rule, f.Element}
		slices.Sort(pair)
		pairs = append(pairs, pair)
	}
	if !slices.EqualFunc(pairs, [][]string{{"fan#small", "pin"}}, slices.Equal) {
		t.Errorf("the redundant-row rule/element slots = %v; want [[fan#small pin]] — both slots name "+
			"their row", pairs)
	}
	// (iii) Group-level findings stay bare.
	for _, f := range c30GroupFindings(second.findings) {
		if strings.Contains(f.Rule, "#") {
			t.Errorf("a group-level %s names %q; group findings stay bare", f.Code, f.Rule)
		}
	}
}

// c30Fan is the literal ladder plus a `fan` outcome whose rule steps
// NOTHING: expanded, it matches `tier in [small, mid]`; unrolled, it is two
// literal rules. Its small row writes tier back to small. `pin` adds a row
// the small one strictly contains.
func c30Fan(t *testing.T, expanded, pin bool) string {
	t.Helper()

	src := c30Edit(t, c30Source(t, c30Literal), `outcomes = ["retry", "escalate", "succeed"]`,
		`outcomes = ["retry", "escalate", "succeed", "fan"]`)
	if expanded {
		src += `
[[rule]]
id = "fan"
use = ["open"]
[rule.match.recognized]
eq = "fan"
[rule.match.tier]
in = ["small", "mid"]
[rule.write]
tier = "small"
`
	} else {
		src += `
[[rule]]
id = "fan-small"
use = ["open"]
[rule.match.recognized]
eq = "fan"
[rule.match.tier]
eq = "small"
[rule.write]
tier = "small"

[[rule]]
id = "fan-mid"
use = ["open"]
[rule.match.recognized]
eq = "fan"
[rule.match.tier]
eq = "mid"
[rule.write]
tier = "small"
`
	}
	if pin {
		src += `
[[rule]]
id = "pin"
use = ["open"]
[rule.match.recognized]
eq = "fan"
[rule.match.tier]
eq = "small"
[rule.guard.all.attempt]
eq = 3
[rule.write]
status = "done"
`
	}
	return src
}

// REQ-38: "The suffix is therefore applied at EMISSION and the in-package
// read-back joins on the AUTHORED ID RECOVERED from the suffixed form — the
// published identity truncated at the first `#`" … "No in-package consumer
// may resolve a per-row finding's `Rule` against authored `[[rules]]` ids by
// equality." … "`internal/graphlint/engine.go::identityKey` … therefore
// satisfies the fence unchanged and keeps its ordering guarantee under a
// suffixed value (A15)."
// ADVERSARIAL
func TestReq38_0030_ReadBackJoinsOnTheRecoveredAuthoredID(t *testing.T) {
	src := c30Pin(t, c30Step)
	m := c30MustLoad(t, "the pinned step ladder", src)
	ids := map[string]bool{}
	for _, r := range m.Rows {
		ids[r.RuleID] = true
	}

	first := c30LintSrc(t, src)
	c30Loaded(t, "the pinned step ladder", first)
	var suffixed int
	for _, f := range first.findings {
		for _, slot := range []string{f.Rule, f.Element} {
			if !strings.Contains(slot, "#") {
				continue
			}
			suffixed++
			authored, _, _ := strings.Cut(slot, "#")
			if !ids[authored] {
				t.Errorf("%s names %q, whose recovered id %q is no authored rule", f.Code, slot, authored)
			}
		}
	}
	if suffixed == 0 {
		t.Errorf("no per-row finding publishes a suffixed identity:\n%s", first.stdout)
	}
	// The recovered join still gates the ambiguous_match arm (S11 c).
	c30RequireOverlapSlots(t)

	// Ordering stays deterministic under suffixed values.
	for range 5 {
		if again := c30LintSrc(t, src); again.stdout != first.stdout {
			t.Fatalf("two lint runs over one model differ:\n%s\n---\n%s", first.stdout, again.stdout)
		}
	}
}

// --- derived outcomes -----------------------------------------------------------

// REQ-43: "A cap cell excluded by the step rule and claimed by no other row —
// `graph-coverage-gap` at lint, the existing finding." and "A stepped tag
// never initialised — `graph-owned-before-write` at lint."
// DOMAIN EDGE
func TestReq43_0030_AnUnclaimedCapAndAnUninitialisedSteppedTagSurfaceAtLint(t *testing.T) {
	uninit := func(name string) string {
		return c30Edit(t, c30Source(t, name), "attempt = 0\ntier", "tier")
	}
	for _, stepped := range []bool{false, true} {
		name := c30Literal
		if stepped {
			name = c30Step
		}
		gap := c30LintSrc(t, c30Uncapped(t, name))
		c30Loaded(t, "uncapped "+name, gap)
		if c30Count(gap.findings, graphlint.CodeCoverageGap) == 0 || !gap.failed {
			t.Errorf("uncapped %s: no blocking graph-coverage-gap:\n%s", name, gap.stdout)
		}
		obw := c30LintSrc(t, uninit(name))
		c30Loaded(t, "uninitialised "+name, obw)
		var onAttempt bool
		for _, f := range c30WithCode(obw.findings, graphlint.CodeOwnedBeforeWrite) {
			if f.Key == "attempt" {
				onAttempt = true
			}
		}
		if !onAttempt || !obw.failed {
			t.Errorf("uninitialised %s: no blocking graph-owned-before-write on attempt:\n%s", name, obw.stdout)
		}
	}
}

// REQ-44: "appending `critical` to a domain whose step rule carries `in =
// [\"small\",\"medium\"]` does NOT refuse at load … and does NOT silently widen
// the ladder (no cell steps into `critical`). The ladder keeps stopping at
// `large`, and the new member surfaces as a BLOCKING `graph-coverage-gap` at
// lint"
// DOMAIN EDGE
func TestReq44_0030_AppendingAMemberNeitherRefusesNorWidensTheLadder(t *testing.T) {
	ladder := func(domain string) string {
		src := c30Source(t, c30Step)
		src = c30Edit(t, src, `domain = ["small", "mid", "large"]`, "domain = "+domain)
		return c30Edit(t, src, `in = ["small", "mid"]`, `in = ["small", "medium"]`)
	}
	before, after := ladder(`["small", "medium", "large"]`), ladder(`["small", "medium", "large", "critical"]`)

	pre := c30LintSrc(t, before)
	c30Loaded(t, "the three-member ladder", pre)
	if n := c30Count(pre.findings, graphlint.CodeCoverageGap); n != 0 {
		t.Errorf("the three-member ladder already carries %d coverage gaps:\n%s", n, pre.stdout)
	}

	m := c30MustLoad(t, "the appended ladder", after)
	for _, r := range m.Rows {
		for _, carrier := range [][]table.TagValue{r.Writes, r.NextTags} {
			for _, tv := range carrier {
				if tv.Key == "tier" && slices.Contains(tv.Value, "critical") {
					t.Errorf("%s steps into critical; the ladder widened silently", r.Identity())
				}
			}
		}
	}
	var escalated []string
	for _, r := range m.Rows {
		if r.RuleID == "escalate" {
			escalated = append(escalated, r.Identity())
		}
	}
	slices.Sort(escalated)
	if want := []string{"ladder.escalate#medium", "ladder.escalate#small"}; !slices.Equal(escalated, want) {
		t.Errorf("escalate rows = %v; want %v", escalated, want)
	}

	post := c30LintSrc(t, after)
	c30Loaded(t, "the appended ladder", post)
	if c30Count(post.findings, graphlint.CodeCoverageGap) == 0 || !post.failed {
		t.Errorf("the appended member does not surface as a blocking graph-coverage-gap:\n%s", post.stdout)
	}
}

// REQ-45: "INTERIOR over-admitted cell conforms, so it mints a row — but that
// row is DEAD, not overlapping" … "`graph-overlap` cannot fire; what fires is
// the advisory `graph-redundant-row`"
// DOMAIN EDGE
func TestReq45_0030_AnUnlessExcludedInteriorCellMintsADeadRow(t *testing.T) {
	c30RequireDeadRowAndNoOverlap(t)

	m := c30MustLoad(t, "S5b", c30Dead(t, true))
	var dead bool
	for _, r := range m.Rows {
		if r.Identity() == "ladder.retry#2" {
			dead = true
			var unless bool
			for _, a := range r.Atoms {
				if a.Key == "attempt" && a.Block == table.BlockUnless {
					unless = true
				}
			}
			if !unless {
				t.Errorf("the dead row does not retain the authored unless atom: %+v", r.Atoms)
			}
		}
	}
	if !dead {
		t.Errorf("the unless-excluded interior cell minted no row; rows %v", func() []string {
			var ids []string
			for _, r := range m.Rows {
				ids = append(ids, r.Identity())
			}
			return ids
		}())
	}
}

// --- S4 / S5 --------------------------------------------------------------------

// REQ-51: "`[initial] attempt = { step = 1 }` and a predicate literal `{ step
// = 1 }`, each loaded. **Expected**: both refuse, and under the categories C1
// names — the `[initial]` one `malformed_initial_declaration`, the predicate
// one `malformed_predicate_atom`." … "The predicate-literal fixture covers all
// three atom positions — `match … eq`, `guard.all … eq`, and as a member of
// an `in` list"
// ADVERSARIAL
func TestReq51_0030_TableShapesOffTheWritePathRefuseUnderTheirOwnCategories(t *testing.T) {
	step := c30Source(t, c30Step)
	c30LoadRefusal(t, "[initial]", c30Edit(t, step, "attempt = 0\ntier", "attempt = { step = 1 }\ntier"),
		table.CatMalformedInitialDeclaration)
	for _, tc := range []struct{ name, old, repl string }{
		{"match eq", "[rule.match.recognized]\neq = \"succeed\"\n", "[rule.match.recognized]\neq = \"succeed\"\n[rule.match.tier]\neq = { step = 1 }\n"},
		{"guard.all eq", "gte = 5", "eq = { step = 1 }"},
		{"in member", "gte = 5", "in = [5, { step = 1 }]"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c30LoadRefusal(t, tc.name, c30Edit(t, step, tc.old, tc.repl), table.CatMalformedPredicateAtom)
		})
	}
	// The write path, where the same shape is admitted, lints.
	c30Loaded(t, c30Step, c30RunLint(t, c30Path(t, c30Step)))
}

// REQ-52: "zero-cell (from an ordered bound …), zero-step, float-step,
// `#`-in-domain, and duplicate-member-in-domain step models loaded; a step
// write on `bool`, `set`, `scalar`, an unbounded `int`, and an `int` whose
// width is not representable; a rule that both steps a key and names it in
// its `clear` list; and a REPRESENTABLE-width `int` whose step magnitude
// carries an admitted cell past `math.MaxInt`" … "**Expected**: each a load
// refusal under `malformed_tag_declaration` whose detail names the admitted
// form or kinds." … "Each carries its paired negative control" … "The A13
// overflow fixture's expected detail names the cell and the bound it passed,
// NOT \"is not an int\""
// INPUT EDGE
func TestReq52_0030_EveryStepGrammarAndKindDefectIsANamedLoadRefusal(t *testing.T) {
	step := c30Source(t, c30Step)
	retryGuard := "[rule.guard.all.attempt]\ngte = 0\nlt = 5\n[rule.write]\nattempt = { step = 1 }"
	withX := func(decl, write string) string {
		src := c30Edit(t, step, "[read.state]", "[tags.x]\nprovenance = \"owned\"\n"+decl+"\n[read.state]")
		src = c30EditAll(t, src, `keys = ["status", "attempt", "tier"]`, `keys = ["status", "attempt", "tier", "x"]`)
		return c30Edit(t, src, "[rule.match.recognized]\neq = \"succeed\"\n[rule.write]\nstatus = \"done\"",
			"[rule.match.recognized]\neq = \"succeed\"\n[rule.write]\nstatus = \"done\"\nx = "+write)
	}
	tierDomain := func(domain string) string {
		return c30Edit(t, step, `domain = ["small", "mid", "large"]`, "domain = "+domain)
	}
	unstepTier := func(src string) string {
		return c30Edit(t, src, "tier = { step = 1 }", `tier = "mid"`)
	}

	for _, tc := range []struct {
		name, src string
		names     []string
		control   string // loads: the paired negative control
	}{
		{"zero cells", c30Edit(t, step, retryGuard, "[rule.guard.all.attempt]\nlt = 0\n[rule.write]\nattempt = { step = 1 }"),
			nil, c30Edit(t, step, retryGuard, "[rule.guard.all.attempt]\nlt = 0\n[rule.write]\nattempt = 1")},
		{"zero step", c30Edit(t, step, "attempt = { step = 1 }", "attempt = { step = 0 }"), []string{"step"}, ""},
		{"float step", c30Edit(t, step, "attempt = { step = 1 }", "attempt = { step = 1.0 }"), []string{"step"}, ""},
		{"# in the stepped domain", tierDomain(`["small", "mid", "large", "la#rge"]`), nil,
			unstepTier(tierDomain(`["small", "mid", "large", "la#rge"]`))},
		{"repeated member in the stepped domain", tierDomain(`["small", "mid", "large", "mid"]`), nil,
			unstepTier(tierDomain(`["small", "mid", "large", "mid"]`))},
		{"bool", withX("kind = \"bool\"\n", "{ step = 1 }"), []string{"int", "enum"}, ""},
		{"set", withX("kind = \"set\"\nelements = [\"a\", \"b\"]\n", "{ step = 1 }"), []string{"int", "enum"}, ""},
		{"scalar", withX("kind = \"scalar\"\n", "{ step = 1 }"), []string{"int", "enum"}, ""},
		{"unbounded int", withX("kind = \"int\"\nmin = 0\nsingle_valued = true\n", "{ step = 1 }"),
			[]string{"int", "enum"}, ""},
		{"unrepresentable width", withX("kind = \"int\"\nmin = 0\nmax = 9223372036854775807\nsingle_valued = true\n",
			"{ step = 1 }"), nil,
			withX("kind = \"int\"\nmin = 0\nmax = 9223372036854775807\nsingle_valued = true\n", "0")},
		{"stepped and cleared", c30Edit(t, step, "id = \"retry\"\nuse = [\"open\"]\n",
			"id = \"retry\"\nuse = [\"open\"]\nclear = [\"attempt\"]\n"), []string{"clear"}, ""},
		{"A13 overflow", c30Edit(t, step, retryGuard,
			"[rule.guard.all.attempt]\ngte = 1\nlt = 2\n[rule.write]\nattempt = { step = 9223372036854775807 }"),
			[]string{"1", "max"}, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.control != "" {
				r := c30LintSrc(t, tc.control)
				c30Loaded(t, "the negative control", r)
			}
			f := c30LoadRefusal(t, tc.name, tc.src, table.CatMalformedTagDeclaration)
			for _, n := range tc.names {
				if !strings.Contains(f.Message, n) {
					t.Errorf("the refusal does not name %q: %s", n, f.Message)
				}
			}
			if strings.Contains(f.Message, "is not an int") {
				t.Errorf("the refusal reads as a rendered-literal conformance failure: %s", f.Message)
			}
		})
	}

	// The negative-width arm is a PRE-EXISTING refusal, restated.
	c30LoadRefusal(t, "min exceeds max", withX("kind = \"int\"\nmin = 5\nmax = 0\nsingle_valued = true\n", "0"),
		table.CatMalformedTagDeclaration)

	// Every defect above is carved out of a model that loads.
	c30Loaded(t, c30Step, c30RunLint(t, c30Path(t, c30Step)))
}

// --- vocabulary, export, docs ------------------------------------------------------

// c30CLICodes collects the CLIError code strings the CLI source declares: the
// `code…` string constants and literal `Code: "…"` fields.
func c30CLICodes(t *testing.T) []string {
	t.Helper()

	root := repoRootFor(t)
	var out []string
	constRe := regexp.MustCompile(`(?m)^\s*code[A-Za-z0-9]*\s*(?:string\s*)?=\s*"([a-z0-9-]+)"`)
	litRe := regexp.MustCompile(`Code:\s*"([a-z0-9-]+)"`)
	err := filepath.WalkDir(filepath.Join(root, "internal", "cli"), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, re := range []*regexp.Regexp{constRe, litRe} {
			for _, m := range re.FindAllStringSubmatch(string(b), -1) {
				if !slices.Contains(out, m[1]) {
					out = append(out, m[1])
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking internal/cli: %v", err)
	}
	slices.Sort(out)
	return out
}

// REQ-34: "No emitted vocabulary — `internal/table::Categories()`, the
// CLIError codes, the finding codes, the graph document's members — gains a
// member, so no `0029:C4` stability tier is owed by this record." With "Done
// = every scenario below green with no new finding code and no new envelope
// member."
// BOUNDARY
func TestReq34_0030_NoEmittedVocabularyGainsAMember(t *testing.T) {
	// The golden holds on a tree where the step form EXISTS — a vocabulary
	// pinned before the feature lands proves nothing about the feature.
	c30Loaded(t, c30Step, c30RunLint(t, c30Path(t, c30Step)))

	var cats []string
	for _, c := range table.Categories() {
		cats = append(cats, string(c))
	}
	wantCats := []string{
		"malformed_toml", "unknown_schema_field", "missing_recognized_outcome_alphabet",
		"malformed_recognized_outcome_alphabet", "unknown_tag", "unknown_context",
		"cyclic_context_inheritance", "write_to_non_owned_tag", "unknown_accessor",
		"malformed_accessor_declaration", "malformed_accessor_binding", "malformed_tag_declaration",
		"malformed_initial_declaration", "reserved_tag_value", "unsupported_version",
		"malformed_predicate_atom", "malformed_escape_declaration", "malformed_model_declaration",
		"malformed_dump_declaration", "malformed_outcome_binding", "malformed_rule_shape",
		"malformed_rule_id", "duplicate_rule_id", "duplicate_model_id", "reserved_tag_key",
		"malformed_emit_declaration", "unknown_emit_key", "emit_value_out_of_domain",
		"command_and_path_conflict", "command_empty", "command_unknown_placeholder",
		"command_shell_interpreter", "command_output_shape", "command_env_conflict",
		"edit_carrier_conflict", "edit_key_mismatch", "edit_anchor_invalid", "edit_template_invalid",
		"edit_clear_invalid", "edit_tag_argv0",
		// Appended after this record by kata wchf's `select` tables; the
		// step form still adds none.
		"command_select_placement", "command_select_key_mismatch", "command_select_invalid",
	}
	if !slices.Equal(cats, wantCats) {
		t.Errorf("table.Categories() = %v; want the golden %v", cats, wantCats)
	}

	codes := append(graphlint.BlockingCodes(), graphlint.AdvisoryCodes()...)
	wantCodes := []string{
		"graph-dangling-edge", "graph-dead-end", "graph-overlap", "graph-coverage-gap",
		"graph-unprovable-coverage", "graph-single-valued-state", "graph-always-present-owned",
		"graph-owned-before-write", "graph-terminal-escape", "graph-product-too-large",
		"graph-coverage-closed-by-escape", "graph-redundant-row", "graph-unreachable-rule",
		"graph-vacuous-atom", "graph-idempotent-write",
	}
	if !slices.Equal(codes, wantCodes) {
		t.Errorf("the finding codes = %v; want the golden %v", codes, wantCodes)
	}

	wantCLI := []string{
		"command-error", "docs-write-failed", "escape-row-shape-breach",
		"flag-invalid-value", "flag-mutually-exclusive", "flag-required",
		"flow-accessor-capability-mismatch", "flow-accessor-failed", "flow-accessor-request-refused",
		"flow-accessor-timeout", "flow-accessor-unknown", "flow-argv-upstream-stop",
		"flow-artifact-invalid", "flow-artifact-missing", "flow-clear-unbound", "flow-gate-denied",
		"flow-gate-indeterminate", "flow-init-carrier-unsupported", "flow-init-class-unsupported",
		"flow-model-invalid", "flow-model-not-found", "flow-read-incomplete", "flow-tag-duplicate",
		"flow-tag-invalid", "flow-tag-owned", "flow-tag-reserved", "flow-write-duplicate",
		"flow-write-edit-refused", "flow-write-failed-applied", "flow-write-invalid",
		"flow-write-non-owned", "flow-write-readback-incomplete", "flow-write-readback-mismatch",
		"flow-write-readback-timeout", "flow-write-unbound", "internal-error", "model-invalid",
		"model-unreadable",
	}
	if got := c30CLICodes(t); !slices.Equal(got, wantCLI) {
		t.Errorf("the CLIError codes the CLI declares = %v; want the golden %v", got, wantCLI)
	}
}

// REQ-53: "`graph` export of `ladder-step.toml` decoded against the shipped
// `intrastate.graph/1` document type. **Expected**: decodes with no new
// member — asserted by STRICT decode into the shipped document type (unknown
// fields rejected) …; and `internal/table::Categories()`, the CLIError codes
// and the finding codes each gain none, asserted as golden enumerations"
// HAPPY PATH
func TestReq53_0030_TheStepLaddersGraphExportDecodesStrictly(t *testing.T) {
	doc := c30RequireStrictGraph(t, c30Path(t, c30Step))
	var stepped int
	for _, r := range doc.Rows {
		if strings.Contains(r.Identity, "#") {
			stepped++
		}
	}
	if stepped != 7 {
		t.Errorf("the export carries %d expanded rows; want 7 (retry#0..4, escalate#small, escalate#mid)", stepped)
	}
	// The literal twin decodes to the same row count, so nothing was dropped.
	lit := c30RequireStrictGraph(t, c30Path(t, c30Literal))
	if len(lit.Rows) != len(doc.Rows) {
		t.Errorf("the step export carries %d rows, the literal one %d", len(doc.Rows), len(lit.Rows))
	}
}

// c30Section returns the markdown text from the first line containing anchor
// to the next heading of the same or a shallower level.
func c30Section(doc, anchor string) string {
	lines := strings.Split(doc, "\n")
	start, level := -1, 0
	for i, l := range lines {
		if strings.Contains(l, anchor) {
			start = i
			break
		}
	}
	if start < 0 {
		return ""
	}
	for i := start; i >= 0; i-- {
		if strings.HasPrefix(lines[i], "#") {
			level = len(lines[i]) - len(strings.TrimLeft(lines[i], "#"))
			start = i
			break
		}
	}
	end := len(lines)
	for i := start + 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "#") {
			if l := len(lines[i]) - len(strings.TrimLeft(lines[i], "#")); level > 0 && l <= level {
				end = i
				break
			}
		}
	}
	return strings.Join(lines[start:end], "\n")
}

// REQ-35: "The authoring guide (`docs/model-authoring.md`) documents the
// form, the admitted kinds, the admitted-cell rule, and the bound refusal
// beside the existing write-block section — and … the four hazards this
// record creates for an author: that an authored `enum` `domain`'s ORDER is
// the step order …; that `unless` never excludes a cell from a step …; that
// on an `enum` the only atom that can exclude the terminal member is an `in`
// re-listing the domain minus that member" … "and that the cell count is the
// LINT cost, not the load cost" … "What the author sees is `lint` not
// returning." … "`intrastate --help-all` is regenerated."
// HAPPY PATH
func TestReq35_0030_TheAuthoringGuideDocumentsTheStepFormAndItsHazards(t *testing.T) {
	b, err := os.ReadFile(filepath.Join(repoRootFor(t), "docs", "model-authoring.md"))
	if err != nil {
		t.Fatalf("reading the authoring guide: %v", err)
	}
	sec := c30Section(string(b), "{ step =")
	if sec == "" {
		t.Fatal("REQ-35 owed: docs/model-authoring.md shows no `{ step = <n> }` write")
	}
	for _, want := range []string{
		"int", "enum", "min", "max", "domain", // the admitted kinds
		"malformed_tag_declaration", // the bound refusal
		"unless", "guard.all",       // the unless hazard and its rewrite
		"order",     // the domain order is the step order
		"in",        // the terminal-member exclusion
		"lint",      // the cell count is the lint cost
		"returning", // what the author sees
	} {
		if !strings.Contains(sec, want) {
			t.Errorf("the step section does not mention %q", want)
		}
	}
	// Beside the existing write-block section.
	if i, j := strings.Index(string(b), "### Rules advance the state"), strings.Index(string(b), "{ step ="); i < 0 || j < i {
		t.Error("the step form is not documented beside the write-block section")
	}
}

// REQ-39: "the field's published meaning narrows from \"the authored rule
// id\" to \"the row's identity\", and `docs/cli-output-contract.md` describes
// `rule` as graph-lint attribution, so an out-of-repo consumer joining on it
// is affected in the same way and `span` is the stable key it should use."
// HAPPY PATH
func TestReq39_0030_TheOutputContractNarrowsRuleToTheRowIdentity(t *testing.T) {
	b, err := os.ReadFile(filepath.Join(repoRootFor(t), "docs", "cli-output-contract.md"))
	if err != nil {
		t.Fatalf("reading the output contract: %v", err)
	}
	doc := string(b)
	if !strings.Contains(strings.ToLower(doc), "row's identity") && !strings.Contains(strings.ToLower(doc), "row identity") {
		t.Error("REQ-39 owed: the contract does not describe a finding's `rule` as the row's identity")
	}
	if !regexp.MustCompile("`[a-z][a-z0-9-]*#[a-z0-9<>-]+`").MatchString(doc) {
		t.Error("the contract shows no suffixed `rule#cell` identity")
	}
	var stable bool
	for _, line := range strings.Split(doc, "\n") {
		if strings.Contains(line, "`span`") && strings.Contains(strings.ToLower(line), "stable") {
			stable = true
		}
	}
	if !stable {
		t.Error("the contract does not name `span` as the stable join key")
	}
}

// REQ-47: "a model with no `{ step = <n> }` write value loads, lints,
// resolves and exports exactly as today, because the new arm is intercepted
// before `load.go::valueMembers` on the write path only." Observable: the
// committed corpus's `lint`/`dump`/`graph` outputs are unchanged; see Q1 for
// the one reading under which a per-row lint `rule` on an `in`-expanded row
// changes.
// BOUNDARY
func TestReq47_0030_AModelThatStepsNothingIsUnchangedButForRowAttribution(t *testing.T) {
	// An `in`-expanded rule and its hand-unrolled twin: no step anywhere.
	expanded, unrolled := c30Fan(t, true, false), c30Fan(t, false, false)
	em, um := c30MustLoad(t, "expanded", expanded), c30MustLoad(t, "unrolled", unrolled)
	if len(em.Rows) != len(um.Rows) {
		t.Fatalf("the expanded model has %d rows, the unrolled one %d", len(em.Rows), len(um.Rows))
	}
	for _, id := range []string{"ladder.fan#small", "ladder.fan#mid"} {
		var found bool
		for _, r := range em.Rows {
			found = found || r.Identity() == id
		}
		if !found {
			t.Errorf("the in-expanded row %s is gone", id)
		}
	}

	// Lint: the same verdict and the same projected findings.
	e, u := c30LintSrc(t, expanded), c30LintSrc(t, unrolled)
	if !slices.Equal(c30ProjSet(e.findings), c30ProjSet(u.findings)) || e.failed != u.failed {
		t.Errorf("lint verdicts differ under the projection\n  expanded: %v\n  unrolled: %v",
			c30ProjSet(e.findings), c30ProjSet(u.findings))
	}
	// Export: the same reach relation.
	eg, ug := c30RequireStrictGraph(t, writeModel(t, expanded)), c30RequireStrictGraph(t, writeModel(t, unrolled))
	if !reflect.DeepEqual(eg.Reach.Nodes, ug.Reach.Nodes) {
		t.Error("the expanded and unrolled models export different reach nodes")
	}

	// The one sanctioned change (Q1 → (a)): the per-row finding on the
	// in-expanded row names its row.
	var rules []string
	for _, f := range c30WithCode(e.findings, graphlint.CodeIdempotentWrite) {
		if strings.HasPrefix(f.Rule, "fan") {
			rules = append(rules, f.Rule)
		}
	}
	if !slices.Equal(rules, []string{"fan#small"}) {
		t.Errorf("the in-expanded row's per-row finding names %v; want [fan#small]", rules)
	}
}
