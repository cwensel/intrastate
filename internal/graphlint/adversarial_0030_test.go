package graphlint_test

// RDR 0030 Phase 3b — adversarial coverage anchored in the record's
// `Trade-offs / Failure Modes` section, fourth bullet:
//
//	A cap cell excluded by the step rule and claimed by no other row —
//	`graph-coverage-gap` at lint, the existing finding.
//
// read with A5 and the MVV ("zero `graph-overlap`, zero
// `graph-coverage-gap`" on a fully claimed ladder): the gap is the
// UNCLAIMED cell's signal, so a ladder claiming every cell lints no gap,
// exactly as its hand-unrolled twin does.
//
// ADV-2 — the S8 shape. C1 has a step SUBSUME a match `in` on its tag and
// emit, per row, the rewritten `match eq = <cell>` AND the expansion's own
// `guard.all eq = <cell>`. The match `eq` splits the rule into one group
// per cell; the `guard.all eq` then makes the stepped tag a dimension of
// that group's scoped product, which is NOT narrowed by the group's own
// match `eq`, so every per-cell group reports "9 of 10 assignments
// uncovered" — a BLOCKING `graph-coverage-gap` on every CLAIMED cell. The
// hand-unrolled twin (one `match eq` per cell) lints clean. So the gap
// fires where F4 says it must not, and the stepped ladder is not the
// literal ladder to lint.

import (
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/graphlint"
	"github.com/cwensel/intrastate/internal/table"
)

func adv2Model(rules string) string {
	return `outcomes = ["retry", "succeed"]
terminal = ["done"]

[model]
id = "adv2"
version = 1

[initial]
status = "open"
attempt = 0

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["open", "done"]
single_valued = true
required = true

[tags.attempt]
provenance = "owned"
kind = "int"
min = 0
max = 9
single_valued = true
required = true

[read.state]
role = "state"
path = "flow.state"
keys = ["status", "attempt"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["status", "attempt"]
timeout = "2s"
read_back = true

[context.open.match.status]
eq = "open"

[context.done.match.status]
eq = "done"

` + rules + `

# The cap claims every cell the ladder does not.
[[rule]]
id = "retry-cap"
use = ["open"]
[rule.match.recognized]
eq = "retry"
[rule.match.attempt]
in = [5, 6, 7, 8, 9]
[rule.write]
status = "done"

[[rule]]
id = "succeed"
use = ["open"]
[rule.match.recognized]
eq = "succeed"
[rule.write]
status = "done"
`
}

// adv2Stepped is the S8 shape: one rule stepping `attempt` over the cells a
// match `in` on it names.
const adv2Stepped = `[[rule]]
id = "retry"
use = ["open"]
[rule.match.recognized]
eq = "retry"
[rule.match.attempt]
in = [0, 1, 2, 3, 4]
[rule.write]
attempt = { step = 1 }
`

// adv2Literal is the same ladder written out cell by cell.
func adv2Literal() string {
	out := ""
	for i := 0; i < 5; i++ {
		c, v := string(rune('0'+i)), string(rune('1'+i))
		out += `[[rule]]
id = "retry-` + c + `"
use = ["open"]
[rule.match.recognized]
eq = "retry"
[rule.match.attempt]
eq = ` + c + `
[rule.write]
attempt = ` + v + `

`
	}
	return out
}

func adv2Gaps(t *testing.T, rules string) []string {
	t.Helper()
	m, err := table.Load([]byte(adv2Model(rules)), "adv2.toml")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	var gaps []string
	for _, f := range graphlint.Run(graphlint.NewRequest(m)).Findings {
		if f.Code == graphlint.CodeCoverageGap {
			gaps = append(gaps, f.Element+": "+f.Message)
		}
	}
	return gaps
}

// TestAdv2_0030_SubsumedInLadderLintsNoCoverageGap: a ladder whose every
// cell is claimed lints zero `graph-coverage-gap` whether written out or
// stepped over a subsumed match `in`.
func TestAdv2_0030_SubsumedInLadderLintsNoCoverageGap(t *testing.T) {
	if gaps := adv2Gaps(t, adv2Literal()); len(gaps) != 0 {
		t.Fatalf("control drifted: the hand-unrolled ladder lints coverage gaps %v", gaps)
	}
	if gaps := adv2Gaps(t, adv2Stepped); len(gaps) != 0 {
		t.Errorf("the stepped ladder claims every cell its literal twin does, yet lints "+
			"%d blocking graph-coverage-gap (F4 reserves the gap for an UNCLAIMED cell):", len(gaps))
		for _, g := range gaps {
			t.Errorf("  %s", g)
		}
	}
}

// adv2PinnedByEq steps `attempt` on a rule whose AUTHORED match `eq` pins the
// one admitted cell; the rest of the ladder is written out cell by cell.
const adv2PinnedByEq = `[[rule]]
id = "retry"
use = ["open"]
[rule.match.recognized]
eq = "retry"
[rule.match.attempt]
eq = 0
[rule.write]
attempt = { step = 1 }
`

// TestAdv2_0030_AuthoredMatchEqStepLintsNoCoverageGap: the D13 resolution
// covers every match `eq` that pins the stepped tag, not only a subsumed
// `in`'s rewrite — a step beside an authored match `eq` carries no
// `guard.all eq` dimension, so its fully claimed ladder lints zero gaps.
func TestAdv2_0030_AuthoredMatchEqStepLintsNoCoverageGap(t *testing.T) {
	rest := strings.SplitAfterN(adv2Literal(), "\n\n", 2)[1]
	if gaps := adv2Gaps(t, adv2PinnedByEq+"\n"+rest); len(gaps) != 0 {
		t.Errorf("a step pinned by an authored match eq lints %d graph-coverage-gap on a fully "+
			"claimed ladder:", len(gaps))
		for _, g := range gaps {
			t.Errorf("  %s", g)
		}
	}
}

// adv2Escape is a bare escape row whose match `in` expands it into one row
// per outcome, `esc#retry` and `esc#succeed`.
const adv2Escape = `[[rule]]
id = "esc"
use = ["open"]
escape = ["no_match"]
[rule.match.recognized]
in = ["retry", "succeed"]
`

// TestReq37_0030_ClosedByEscapeNamesTheExpandedEscapeRow is Phase 3a
// FAIL-2's regression. `graph-coverage-closed-by-escape` names the bare
// escape ROW that closes the group, so under C3 ("every surface that names a
// ROW publishes the suffixed identity") an expanded escape row publishes
// `esc#<cell>` in both the Rule slot and the message, never the bare `esc`.
func TestReq37_0030_ClosedByEscapeNamesTheExpandedEscapeRow(t *testing.T) {
	m, err := table.Load([]byte(adv2Model(adv2Literal()+adv2Escape)), "adv2-esc.toml")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	got := map[string]string{}
	for _, f := range graphlint.Run(graphlint.NewRequest(m)).Findings {
		if f.Code == graphlint.CodeCoverageClosedByEscape {
			got[f.Rule] = f.Message
		}
	}
	for _, want := range []string{"esc#retry", "esc#succeed"} {
		msg, ok := got[want]
		if !ok {
			t.Errorf("no %s finding names the row %q; got rules %v",
				graphlint.CodeCoverageClosedByEscape, want, got)
			continue
		}
		if !strings.Contains(msg, `"`+want+`"`) {
			t.Errorf("finding for %q names another row in its message: %s", want, msg)
		}
	}
	if msg, ok := got["esc"]; ok {
		t.Errorf("a closed-by-escape finding names the bare authored id \"esc\": %s", msg)
	}
}
