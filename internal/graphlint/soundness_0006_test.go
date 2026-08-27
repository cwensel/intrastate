package graphlint_test

// RDR 0006 — the soundness direction (REQ-105, REQ-107, REQ-110..REQ-112)
// and the accepted false positives it deliberately buys (REQ-117 / SC-10).

import (
	"fmt"
	"slices"
	"testing"

	"github.com/cwensel/intrastate/internal/graphlint"
	"github.com/cwensel/intrastate/internal/guard"
	"github.com/cwensel/intrastate/internal/table"
)

// REQ-105: "A selection context is *reachable* when some reachable
// owned-state satisfies its match pattern; a row's *reachable
// predecessors* are the rows on any root-to-source path."
// HAPPY PATH
func TestReq105_ASelectionContextIsReachableWhenSomeNodeSatisfiesIt(t *testing.T) {
	// The clause is "SOME reachable owned-state satisfies its match
	// pattern" — not "the root does". The row under test therefore matches
	// `status = mid`, a NON-root value some other row writes, so an engine
	// that only ever treated root-matching contexts as reachable would
	// report it unreachable and fail here. Letting the sole reachable row
	// match the `[initial]` value (the earlier shape) made the two
	// readings indistinguishable.
	//
	// `status = orphan` is written by no row, so its context is not
	// reachable under either reading — and the advisory names exactly it.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "mid", "b", "orphan"]
single_valued = true
required = true
`
	const body = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "leaves-root"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "mid"

[[rule]]
id = "reachable"
[rule.match.status]
eq = "mid"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "unreachable"
[rule.match.status]
eq = "orphan"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`
	m := mustLoad(t, source(decls, body))

	// The relation's own oracle: `orphan` is in no reachable node.
	var sawOrphan bool
	for _, n := range graphlint.Reach(m) {
		if slices.Contains(n.Values["status"], "orphan") {
			sawOrphan = true
		}
	}
	if sawOrphan {
		t.Errorf("`status = orphan` is reachable, but no row writes it; "+
			"nodes=%v", graphlint.Reach(m))
	}

	// PRECONDITION: the row under test matches no root value, so the
	// assertion below distinguishes "some reachable owned-state" from
	// "the root".
	for _, tv := range m.Initial {
		if tv.Key == "status" && slices.Contains(tv.Value, "mid") {
			t.Fatalf("`status = mid` is a ROOT value, so the reachable row "+
				"below is reachable under either reading and the clause is "+
				"untested; initial=%v", m.Initial)
		}
	}

	r := graphlint.Run(graphlint.NewRequest(m))
	if !namesRule(r, graphlint.CodeUnreachableRule, "unreachable") {
		t.Errorf("no %s finding names the row whose selection context no "+
			"reachable owned-state satisfies; report:%s",
			graphlint.CodeUnreachableRule, render(r))
	}
	for _, f := range withCode(r, graphlint.CodeUnreachableRule) {
		if f.Rule == "reachable" || f.Rule == "leaves-root" {
			t.Errorf("%s names the row %q, whose selection context a "+
				"reachable owned-state DOES satisfy; reachability is over "+
				"every reachable owned-state, not the root alone; report:%s",
				graphlint.CodeUnreachableRule, f.Rule, render(r))
		}
	}
}

// REQ-107: "It is a different predicate from RDR 0003's \"can refuse\"
// test, which is decided over the optionality field alone and never
// consults this relation."
// ADVERSARIAL
func TestReq107_CanRefuseIsDecidedOverOptionalityAndNeverConsultsReachability(t *testing.T) {
	// Two models identical but for the reachability of the row's context.
	// The can-refuse verdict must be the SAME in both, because it reads
	// the declaration and never the relation.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b", "orphan"]
single_valued = true
required = true

[tags.opt]
provenance = "owned"
kind = "enum"
domain = ["p", "q"]
single_valued = true
`
	const reachable = `
terminal = ["done"]

[initial]
status = "a"
opt = "p"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "subject"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.opt]
eq = "p"
[rule.write]
status = "b"
`
	const unreachable = `
terminal = ["done"]

[initial]
status = "a"
opt = "p"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "subject"
[rule.match.status]
eq = "orphan"
[rule.match.recognized]
eq = "go"
[rule.guard.all.opt]
eq = "p"
[rule.write]
status = "b"
`
	rm := mustLoad(t, source(decls, reachable))
	um := mustLoad(t, source(decls, unreachable))

	// RDR 0003's predicate is the authority and it reads declarations only.
	rRow := rowByRuleID(t, rm, "subject")
	uRow := rowByRuleID(t, um, "subject")
	if guard.CanRefuse(rm.Tags, rRow) != guard.CanRefuse(um.Tags, uRow) {
		t.Fatalf("fixture precondition failed: RDR 0003's own can-refuse "+
			"verdict already differs between the two models (%v vs %v)",
			guard.CanRefuse(rm.Tags, rRow), guard.CanRefuse(um.Tags, uRow))
	}
	if !guard.CanRefuse(rm.Tags, rRow) {
		t.Fatal("fixture precondition failed: the subject row cannot refuse, " +
			"so there is no withholding to compare")
	}

	// So this RDR's withholding verdict must agree in both: the reachable
	// model withholds, and the unreachable one's row is still a can-refuse
	// row — the relation filtered the GROUP, it did not change the
	// declaration-level predicate.
	rr := graphlint.Run(graphlint.NewRequest(rm))
	if !namesRule(rr, graphlint.CodeUnprovableCoverage, "subject") {
		t.Errorf("the reachable model did not withhold the claim for the "+
			"can-refuse row; report:%s", render(rr))
	}
}

// REQ-110: "**Existential checks read merged nodes directly.** Overlap …,
// dangling edge, and owned-set-before-match ask whether *some* witness
// exists."
// REQ-112: "invariant 6's \"every reachable owned-state that satisfies the
// row's match pattern\" is the normative form, evaluated against merged
// fixpoint nodes. … where the two readings differ the node form governs."
// DOMAIN EDGE
func TestReq110And112_OwnedSetBeforeMatchReadsMergedFixpointNodes(t *testing.T) {
	// Two paths converge on `status = mid`: one has written `opt`, the
	// other has not. The MERGED node therefore does not hold `opt` on
	// every concrete view it stands for, and invariant 6's node form
	// reports the row that reads it.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["seed", "mid", "end"]
single_valued = true
required = true

[tags.opt]
provenance = "owned"
kind = "enum"
domain = ["p", "q"]
single_valued = true
`
	const body = `
terminal = ["done"]

[initial]
status = "seed"

[context.done]
[context.done.match.status]
eq = "end"

[[rule]]
id = "writes-opt"
[rule.match.status]
eq = "seed"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "mid"
opt = "p"

[[rule]]
id = "skips-opt"
[rule.match.status]
eq = "seed"
[rule.match.recognized]
eq = "stop"
[rule.write]
status = "mid"

[[rule]]
id = "reads-opt"
[rule.match.status]
eq = "mid"
[rule.match.recognized]
eq = "go"
[rule.guard.all.opt]
eq = "p"
[rule.write]
status = "end"
`
	m := mustLoad(t, source(decls, body))

	// The node form, not a per-path enumeration. `skips-opt` reaches the
	// owned-state {status=mid} and `writes-opt` reaches {status=mid, opt=p}:
	// two distinct abstract owned-states, since REQ-101 makes "absent" a
	// value of a tag's dimension, and one node each. What must NOT appear is
	// a second node for the SAME owned-state — that is the path-sensitive
	// enumeration REQ-108's fixpoint rules out — so the oracle counts nodes
	// per owned-state identity rather than per `status` value.
	//
	// Counting `status = mid` nodes instead would demand the two footprints
	// be folded into one, which under absence-dominates erases a key no
	// sibling path establishes. REQ-108 licenses the join only between edges
	// reaching "the same successor", and these two do not: see FAIL-2 in
	// artifacts/verification.md, where folding divergent footprints is the
	// false green.
	seen := map[string]int{}
	for _, n := range graphlint.Reach(m) {
		if !slices.Contains(n.Values["status"], "mid") {
			continue
		}
		seen[fmt.Sprint(n.Values)]++
	}
	for id, count := range seen {
		if count != 1 {
			t.Errorf("%d nodes carry the owned-state %s; the traversal is a "+
				"fixpoint over MERGED nodes, so every path reaching one "+
				"owned-state reaches ONE node; nodes=%v", count, id,
				graphlint.Reach(m))
		}
	}
	// The `opt`-absent path must be represented: it is the node invariant 6
	// reads below, and losing it would make the finding vacuous.
	var sawAbsent bool
	for _, n := range graphlint.Reach(m) {
		if slices.Contains(n.Values["status"], "mid") && len(n.Values["opt"]) == 0 {
			sawAbsent = true
		}
	}
	if !sawAbsent {
		t.Errorf("no reachable owned-state carries `status = mid` without "+
			"`opt`, but `skips-opt` reaches exactly that; nodes=%v",
			graphlint.Reach(m))
	}

	r := graphlint.Run(graphlint.NewRequest(m))
	if !namesRule(r, graphlint.CodeOwnedBeforeWrite, "reads-opt") {
		t.Errorf("no %s finding names the row reading an owned tag that the "+
			"merged node does not hold on every path; report:%s",
			graphlint.CodeOwnedBeforeWrite, render(r))
	}
}

// REQ-111: "**Universal checks must not.**" Coverage's remedy is that "it
// never reads a node at all"; dead end's remedy is to "**split** the node
// on terminal-participating keys first, recovering exactness."
// ADVERSARIAL
func TestReq111_UniversalChecksDoNotReadWidenedNodesDirectly(t *testing.T) {
	// Coverage's remedy, restated as an oracle: two models whose GROUPS
	// are identical and whose reachable node sets differ must produce the
	// same coverage verdict, because no node enters the computation.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["seed", "a", "end"]
single_valued = true
required = true

[tags.flag]
provenance = "owned"
kind = "bool"
single_valued = true
required = true
`
	// One route into `status = a`, so the node's `flag` set is narrow.
	const oneRoute = `
terminal = ["done"]

[initial]
status = "seed"
flag = "true"

[context.done]
[context.done.match.status]
eq = "end"

[[rule]]
id = "seed-on"
[rule.match.status]
eq = "seed"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "a"
flag = "true"

[[rule]]
id = "only-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "end"
`
	// Two routes, so the merged node's `flag` set is {true, false}. The
	// group under test is byte-identical.
	const twoRoutes = `
terminal = ["done"]

[initial]
status = "seed"
flag = "true"

[context.done]
[context.done.match.status]
eq = "end"

[[rule]]
id = "seed-on"
[rule.match.status]
eq = "seed"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "a"
flag = "true"

[[rule]]
id = "seed-off"
[rule.match.status]
eq = "seed"
[rule.match.recognized]
eq = "stop"
[rule.write]
status = "a"
flag = "false"

[[rule]]
id = "only-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "end"
`
	one := lint(t, decls, oneRoute)
	two := lint(t, decls, twoRoutes)

	// The `only-on` group leaves `flag = false` uncovered in BOTH, because
	// coverage never reads a node. A node-narrowing reading would certify
	// the one-route model and produce a false green.
	for name, r := range map[string]graphlint.Report{"one-route": one, "two-routes": two} {
		if !namesRule(r, graphlint.CodeCoverageGap, "only-on") {
			t.Errorf("%s: no %s finding names `only-on`. Coverage is a "+
				"universal claim computed from the group's authored rows and "+
				"declared domains alone — it never reads a node, so the "+
				"reachable set must not change the verdict; report:%s",
				name, graphlint.CodeCoverageGap, render(r))
		}
	}
}

// REQ-117 / SC-10: "Both are accepted false positives whose cure is an
// explicit write, clear, or terminal declaration on the model, never a
// guard-aware lint."
// DOMAIN EDGE
func TestReq117_AcceptedFalsePositivesAreCuredByTheModelNotByAWeakerLint(t *testing.T) {
	// SC-10's owned-before-write half: the row reads an owned tag held on
	// every GUARD-FEASIBLE path but absent on one path lint cannot prune.
	// The finding is expected — it is the accepted false positive.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["seed", "mid", "end"]
single_valued = true
required = true

[tags.opt]
provenance = "owned"
kind = "enum"
domain = ["p", "q"]
single_valued = true
`
	const uncured = `
terminal = ["done"]

[initial]
status = "seed"

[context.done]
[context.done.match.status]
eq = "end"

[[rule]]
id = "writes-opt"
[rule.match.status]
eq = "seed"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "mid"
opt = "p"

[[rule]]
id = "skips-opt"
[rule.match.status]
eq = "seed"
[rule.match.recognized]
eq = "stop"
[rule.write]
status = "mid"

[[rule]]
id = "reads-opt"
[rule.match.status]
eq = "mid"
[rule.match.recognized]
eq = "go"
[rule.guard.all.opt]
eq = "p"
[rule.write]
status = "end"
`
	r := lint(t, decls, uncured)
	if !namesRule(r, graphlint.CodeOwnedBeforeWrite, "reads-opt") {
		t.Errorf("the accepted false positive was not reported; the "+
			"over-approximation is deliberate and its cure is a model edit, "+
			"never a guard-aware lint; report:%s", render(r))
	}

	// THE CURE: an explicit write on the other path. The finding clears
	// because the model changed, not because lint got weaker.
	const cured = `
terminal = ["done"]

[initial]
status = "seed"

[context.done]
[context.done.match.status]
eq = "end"

[[rule]]
id = "writes-opt"
[rule.match.status]
eq = "seed"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "mid"
opt = "p"

[[rule]]
id = "skips-opt"
[rule.match.status]
eq = "seed"
[rule.match.recognized]
eq = "stop"
[rule.write]
status = "mid"
opt = "q"

[[rule]]
id = "reads-opt"
[rule.match.status]
eq = "mid"
[rule.match.recognized]
eq = "go"
[rule.guard.all.opt]
eq = "p"
[rule.write]
status = "end"
`
	c := lint(t, decls, cured)
	for _, f := range withCode(c, graphlint.CodeOwnedBeforeWrite) {
		if f.Rule == "reads-opt" {
			t.Errorf("the explicit write on the other path did not cure the "+
				"finding; the cure is a model edit and it must work; "+
				"report:%s", render(c))
		}
	}
}

// REQ-16 / REQ-45 (declaration-reading half): lint reads optionality from
// RDR 0003's declaration model, and the two packages must agree — a second
// reading is exactly the inference the clause forbids.
// BOUNDARY
func TestReq45_OptionalityAndSingleValuednessComeFromTheDeclarationModel(t *testing.T) {
	m := mustLoad(t, source(optionalGuardDecls, rootAndTerminal+advanceRule))

	for key, want := range map[string]struct {
		optional     bool
		singleValued bool
	}{
		// `status` carries `required = true`, so it is always-present.
		"status": {optional: false, singleValued: true},
		// `opt` carries no marker, so it DEFAULTS to optional.
		"opt": {optional: true, singleValued: true},
	} {
		got := guard.DeclarationOf(m, key)
		if got.Optional != want.optional {
			t.Errorf("tag %q optional = %v; want %v — read from the "+
				"declaration, never inferred", key, got.Optional, want.optional)
		}
		if got.SingleValued != want.singleValued {
			t.Errorf("tag %q single-valued = %v; want %v", key,
				got.SingleValued, want.singleValued)
		}
		if m.Tags[key].Provenance != table.ProvenanceOwned {
			t.Errorf("tag %q provenance = %q; want owned", key,
				m.Tags[key].Provenance)
		}
	}
}
