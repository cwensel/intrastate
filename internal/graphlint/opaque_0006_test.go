package graphlint_test

// REQ-101 / REQ-106 / REQ-109: a tag with no finite declared domain
// abstracts to held/absent. The traversal must carry such a tag as the
// single opaque value rather than as its concrete authored string —
// otherwise a value atom over it is PRUNED, the edge behind it is never
// taken, and the live rows downstream are reported
// `graph-unreachable-rule` (advisory) while their blocking overlap goes
// unchecked. That is the false-green direction REQ-106 forbids, and the
// same mechanism deviation D11 recorded on the presence axis.

import (
	"slices"
	"testing"

	"github.com/cwensel/intrastate/internal/graphlint"
)

// scalarGateDecls declares the owned `scalar` whose domain is not finite
// alongside the ordinary `phase` enum the downstream group lives on.
const scalarGateDecls = `
[tags.phase]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[tags.free]
provenance = "owned"
kind = "scalar"
required = true
`

// unboundedIntGateDecls is the same shape with an owned `int` declared
// WITHOUT `min`/`max`. It exercises `guard.AssignmentCount`'s other
// no-finite-domain arm, so the abstraction is gated on the declared
// domain rather than on the `scalar` kind alone.
const unboundedIntGateDecls = `
[tags.phase]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[tags.free]
provenance = "owned"
kind = "int"
single_valued = true
required = true
`

// scalarGateBody roots at `free = "start"`, gates the only edge out of
// `phase = a` behind `free eq "other"`, and puts two IDENTICAL rows
// behind that gate. The duplicate pair is a blocking `graph-overlap`
// defect; it is visible only if the gated edge is traversed.
const scalarGateBody = `
terminal = ["done"]

[initial]
phase = "a"
free = "start"

[context.done]
[context.done.match.phase]
eq = "b"

[[rule]]
id = "gate"
[rule.match.phase]
eq = "a"
[rule.match.free]
eq = "other"
[rule.match.recognized]
eq = "go"
[rule.write]
phase = "b"

[[rule]]
id = "dup1"
[rule.match.phase]
eq = "b"
[rule.match.recognized]
eq = "stop"
[rule.write]
phase = "b"

[[rule]]
id = "dup2"
[rule.match.phase]
eq = "b"
[rule.match.recognized]
eq = "stop"
[rule.write]
phase = "b"
`

// intGateBody is scalarGateBody with int-valued literals for `free`, so
// the loader accepts the root and the gate against an `int` declaration.
const intGateBody = `
terminal = ["done"]

[initial]
phase = "a"
free = "1"

[context.done]
[context.done.match.phase]
eq = "b"

[[rule]]
id = "gate"
[rule.match.phase]
eq = "a"
[rule.match.free]
eq = "2"
[rule.match.recognized]
eq = "go"
[rule.write]
phase = "b"

[[rule]]
id = "dup1"
[rule.match.phase]
eq = "b"
[rule.match.recognized]
eq = "stop"
[rule.write]
phase = "b"

[[rule]]
id = "dup2"
[rule.match.phase]
eq = "b"
[rule.match.recognized]
eq = "stop"
[rule.write]
phase = "b"
`

// TestReq101_ScalarGatedEdgeIsTraversable is the load-bearing regression:
// the blocking overlap behind a scalar-gated edge must be REPORTED, not
// masked behind advisory unreachability at exit 0.
func TestReq101_ScalarGatedEdgeIsTraversable(t *testing.T) {
	for name, tc := range map[string]struct{ decls, body string }{
		"scalar":        {scalarGateDecls, scalarGateBody},
		"unbounded-int": {unboundedIntGateDecls, intGateBody},
	} {
		t.Run(name, func(t *testing.T) {
			r := lint(t, tc.decls, tc.body)

			if !namesRule(r, graphlint.CodeOverlap, "dup1") &&
				!namesRule(r, graphlint.CodeOverlap, "dup2") {
				t.Fatalf("no graph-overlap over the duplicate pair behind the "+
					"gated edge: the gate's value atom over a tag with no "+
					"finite domain pruned the edge, so a BLOCKING defect is "+
					"masked; report:%s", render(r))
			}
			for _, id := range []string{"gate", "dup1", "dup2"} {
				if namesRule(r, graphlint.CodeUnreachableRule, id) {
					t.Errorf("rule %q reported graph-unreachable-rule: it is "+
						"live while `free` is held, and a tag with no finite "+
						"domain abstracts to held/absent; report:%s",
						id, render(r))
				}
			}
		})
	}
}

// TestReq101_NonFiniteTagAbstractsToOpaqueInTheRoot pins the normalization
// itself: it is gated on `guard.AssignmentCount`, not applied blanket. A
// finitely-declared tag keeps its concrete authored member.
func TestReq101_NonFiniteTagAbstractsToOpaqueInTheRoot(t *testing.T) {
	m := mustLoad(t, source(scalarGateDecls, scalarGateBody))

	nodes := graphlint.Reach(m)
	if len(nodes) == 0 {
		t.Fatal("the traversal produced no nodes")
	}
	root := nodes[0]

	if got := root.Values["free"]; !slices.Equal(got, []string{graphlint.OpaqueValue}) {
		t.Errorf("root holds `free` as %v; a tag with no finite declared "+
			"domain must abstract to %v", got, []string{graphlint.OpaqueValue})
	}
	if got := root.Values["phase"]; !slices.Equal(got, []string{"a"}) {
		t.Errorf("root holds `phase` as %v; a finitely-declared enum must "+
			"keep its concrete authored member [a]", got)
	}
}

// TestReq101_UnboundedIntAbstractsToOpaqueInTheRoot is the same pin over
// `AssignmentCount`'s other no-finite-domain arm: an `int` declared
// without `min`/`max`.
func TestReq101_UnboundedIntAbstractsToOpaqueInTheRoot(t *testing.T) {
	m := mustLoad(t, source(unboundedIntGateDecls, intGateBody))

	nodes := graphlint.Reach(m)
	if len(nodes) == 0 {
		t.Fatal("the traversal produced no nodes")
	}
	if got := nodes[0].Values["free"]; !slices.Equal(got, []string{graphlint.OpaqueValue}) {
		t.Errorf("root holds an unbounded `int` as %v; it declares no finite "+
			"domain and must abstract to %v",
			got, []string{graphlint.OpaqueValue})
	}
}

// TestReq101_SuccessorAbstractsANonFiniteWrite pins the second call site:
// a row WRITING a tag with no finite domain must put the opaque value in
// the successor, so an atom over that tag downstream stays satisfiable.
func TestReq101_SuccessorAbstractsANonFiniteWrite(t *testing.T) {
	const body = `
terminal = ["done"]

[initial]
phase = "a"
free = "start"

[context.done]
[context.done.match.phase]
eq = "b"

[[rule]]
id = "advance"
[rule.match.phase]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
phase = "b"
free = "written"
`
	m := mustLoad(t, source(scalarGateDecls, body))

	var saw bool
	for _, n := range graphlint.Reach(m) {
		if slices.Contains(n.Values["phase"], "b") {
			saw = true
			if got := n.Values["free"]; !slices.Equal(got, []string{graphlint.OpaqueValue}) {
				t.Errorf("successor holds the written `free` as %v; the write "+
					"lands on a tag with no finite domain and must abstract "+
					"to %v", got, []string{graphlint.OpaqueValue})
			}
		}
	}
	if !saw {
		t.Fatalf("no successor holding `phase = b`; nodes:%s",
			renderNodes(graphlint.Reach(m)))
	}
}

// boundedIntDecls declares `free` as an `int` WITH `min`/`max`. It is the
// finite control for the two unbounded-int cases above: the same declared
// KIND, differing only in whether a bound is carried.
const boundedIntDecls = `
[tags.phase]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[tags.free]
provenance = "owned"
kind = "int"
min = 0
max = 9
single_valued = true
required = true
`

// TestReq101_BoundedIntKeepsItsConcreteValue is the control that makes the
// two unbounded-int pins discriminating. The abstraction is gated on
// `guard.AssignmentCount`, not on the declared kind, so an `int` CARRYING a
// bound declares a finite domain and keeps its concrete authored value —
// at the root and across a write — exactly as an enum does. Without this,
// a `heldValues` keyed on `kind == "int"`, or one abstracting every tag
// blanket, would satisfy every other case in this file.
func TestReq101_BoundedIntKeepsItsConcreteValue(t *testing.T) {
	const body = `
terminal = ["done"]

[initial]
phase = "a"
free = "1"

[context.done]
[context.done.match.phase]
eq = "b"

[[rule]]
id = "advance"
[rule.match.phase]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
phase = "b"
free = "7"
`
	m := mustLoad(t, source(boundedIntDecls, body))

	nodes := graphlint.Reach(m)
	if len(nodes) == 0 {
		t.Fatal("the traversal produced no nodes")
	}
	if got := nodes[0].Values["free"]; !slices.Equal(got, []string{"1"}) {
		t.Errorf("root holds the bounded `int` as %v; it declares a finite "+
			"domain and must keep its concrete authored value [1]", got)
	}

	var saw bool
	for _, n := range nodes {
		if slices.Contains(n.Values["phase"], "b") {
			saw = true
			if got := n.Values["free"]; !slices.Equal(got, []string{"7"}) {
				t.Errorf("successor holds the written bounded `int` as %v; "+
					"the write lands on a finite domain and must keep its "+
					"concrete value [7]", got)
			}
		}
	}
	if !saw {
		t.Fatalf("no successor holding `phase = b`; nodes:%s",
			renderNodes(nodes))
	}
}
