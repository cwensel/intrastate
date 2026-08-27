package graphlint_test

// RDR 0010 Phase 3b — adversarial coverage for the ∅-rooted graph lint.
//
// FAILURE MODES ANCHOR. The record's `Trade-offs / Failure Modes` section
// files exactly one silent class, and both tests here sit inside it:
//
//	**Silent (guarded)**: a decision table with no coverage findings when it
//	should have them. Two distinct paths, both now guarded: the ∅-root
//	regression (the traversal seeds nothing, `checkGroups` skips every
//	group) …
//
// The record guards the ∅-root regression in ONE direction — the traversal
// seeding NOTHING. The opposite direction is unguarded and is what these
// tests attack: seeding the ∅ node makes `len(a.nodes) == 1` over a
// decision table, and every machine-only invariant whose ONLY gate was
// "the traversal found no node" now runs over a class `0010:C5` declares it
// vacuous for.
//
// `0010:C5` is unambiguous about the intended silence:
//
//	The missing-root arm of invariant 1 (`0006:C18`), dead end (2),
//	always-present-owned (5), owned-set-before-match (6), and single-valued
//	state are vacuous by construction over this class and MUST NOT emit.
//	… Declared-terminal handling (7) proper, and the escape arm of it, stay
//	silent by construction: they key on a declared terminal, which C2
//	forbids.
//
// That last sentence is the load-bearing mistake. `checkTerminalEscape`
// (invariant 7) does NOT key on a declared terminal — it gates on
// `len(a.nodes) == 0` alone (`internal/graphlint/analysis.go`), which was
// true for every rootless model before this RDR and is false for every
// decision table after it. `checkDeadEnd` (invariant 2) gates on
// `len(a.model.Terminal) == 0 || len(a.nodes) == 0`, so it too is now
// reachable for a decision table the moment `terminal` is authored — and
// `terminal` on a decision table is a LINT finding under C2, not a load
// refusal, so such a model reaches lint.
//
// Both leaks are worse than noise. The message invariant 7 emits reads
// "declare it in the root `terminal` list" — a remedy `0010:C2` makes an
// authoring error for this class, so following the finding's own advice
// yields `graph-dangling-edge` instead. The record's `Recovery` bullet
// names the only two real remedies ("drop `class`" or "add
// `provenance = \"owned\"` state and an `[initial]`"), neither of which the
// finding states.

import (
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/graphlint"
	"github.com/newcoinc/intrastate/internal/table"
)

// advLint loads src through the real loader and runs the real engine.
// Nothing is mocked: a finding here is a finding a user sees.
func advLint(t *testing.T, src string) []string {
	t.Helper()
	m, err := table.Load([]byte(src), "adv-0010.toml")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if !table.IsDecisionTable(m) {
		t.Fatalf("fixture is not a decision table: class = %q", m.Class)
	}
	rep := graphlint.Run(graphlint.NewRequest(m))
	codes := make([]string, 0, len(rep.Findings))
	for _, f := range rep.Findings {
		codes = append(codes, f.Code)
	}
	return codes
}

// advEscapeOnly is a decision table whose whole population is one escape
// row. It is a legal model: `0010:C2` strips a decision table of
// `[initial]`, `terminal`, owned accessors, and write blocks, and
// `0002:C4`'s "an ordinary transition rule MUST contain a write block" is
// conditioned away by `0010:C2` — so nothing in the grammar demands an
// ordinary row. The loader accepts it; only lint has an opinion.
const advEscapeOnly = `
outcomes = ["go"]

[model]
id = "adv-escape-only"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[[rule]]
id = "r-esc"
escape = ["no_match"]
[rule.match.recognized]
eq = "go"
[rule.emit]
answer = "FALLBACK"
`

// TestAdv0010_TerminalEscapeStaysSilentOverADecisionTable pins `0010:C5`'s
// "Declared-terminal handling (7) proper, and the escape arm of it, stay
// silent by construction" against the ∅-rooted traversal that makes them
// reachable.
//
// FAILURE MODE GUARDED: the ∅-root seeding admits a machine-only invariant
// it should refuse. `checkTerminalEscape`'s only gate is `len(a.nodes) ==
// 0`; seeding the ∅ node makes that gate false for every decision table, so
// invariant 7 walks the single ∅ node, finds no outgoing ordinary row, and
// accuses a stateless table of "relying on an inferred terminal". Its
// remedy — "declare it in the root `terminal` list" — is the very
// declaration `0010:C2` forbids this class, so the finding is not merely
// noise: it routes the author into a second, different finding.
func TestAdv0010_TerminalEscapeStaysSilentOverADecisionTable(t *testing.T) {
	for _, code := range advLint(t, advEscapeOnly) {
		if code == graphlint.CodeTerminalEscape {
			t.Fatalf("lint emitted %s over a decision table; `0010:C5` "+
				"requires declared-terminal handling (invariant 7) to stay "+
				"silent for this class, and its remedy names the `terminal` "+
				"declaration `0010:C2` forbids", code)
		}
	}
}

// TestAdv0010_TerminalEscapeSilentOnARuleFreeDecisionTable is the same leak
// at its smallest reproducer: a decision table with NO rules at all. There
// is no escape row to blame and no terminal to point at, yet invariant 7
// still fires on the ∅ node.
//
// FAILURE MODE GUARDED: the same ∅-root over-admission, minted from a model
// carrying nothing but a class declaration and a tag table. If the guard
// were "the class is vacuous for this invariant" rather than "the traversal
// found no node", this model could not produce a finding at all.
func TestAdv0010_TerminalEscapeSilentOnARuleFreeDecisionTable(t *testing.T) {
	src := advEscapeOnly[:strings.Index(advEscapeOnly, "[[rule]]")]
	for _, code := range advLint(t, src) {
		if code == graphlint.CodeTerminalEscape {
			t.Fatalf("lint emitted %s over a rule-free decision table; "+
				"`0010:C5` requires invariant 7 to stay silent for this class",
				code)
		}
	}
}

// TestAdv0010_DeadEndStaysSilentOverADecisionTable pins `0010:C5`'s "dead
// end (2) … vacuous by construction over this class and MUST NOT emit".
//
// FAILURE MODE GUARDED: the second half of the same ∅-root over-admission.
// `checkDeadEnd` gates on `len(a.model.Terminal) == 0 || len(a.nodes) == 0`.
// `0010:C2` forbids a decision table from declaring `terminal`, but it
// enforces that prohibition AT LINT — "a terminal predicate over a non-owned
// tag is 0006's `graph-dangling-edge` terminal arm, which C5 keeps live for
// exactly this reason" — so the model LOADS and reaches the engine with a
// non-empty `Terminal`. With the ∅ node seeded, the second disjunct is false
// too, and invariant 2 accuses the ∅ owned-state of being a dead end. The
// author now reads two findings for one authoring mistake, only one of which
// (`graph-dangling-edge`) `0010:C2` licenses.
func TestAdv0010_DeadEndStaysSilentOverADecisionTable(t *testing.T) {
	src := "terminal = [\"done\"]\n" + advEscapeOnly + `
[context.done]
[context.done.match.recognized]
eq = "go"
`
	for _, code := range advLint(t, src) {
		if code == graphlint.CodeDeadEnd {
			t.Fatalf("lint emitted %s over a decision table; `0010:C5` "+
				"lists dead end (invariant 2) among the machine-only "+
				"invariants that are vacuous by construction for this class "+
				"and MUST NOT emit", code)
		}
	}
}
