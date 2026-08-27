package table_test

// RDR 0008 Phase 3b — adversarial coverage for the failure payload's
// determinism.
//
// The RDR's Failure Modes section promises a "Visible" failure whose data
// "carries the direction-specific rule identifier, the offending name, and
// the remedy name". `0008:C3` (REQ-27) makes the rule identifier "a
// comparable token, not prose: a golden test asserts it byte-for-byte",
// and REQ-33 makes it the value "for golden assertions and remediation
// lookup". Both uses presuppose that identical bytes yield an identical
// payload — a remediation lookup that returns "rename TO `recognized`" on
// one run and "rename AWAY FROM `recognized`" on the next is worse than
// no guidance, because the two directions prescribe opposite edits.
//
// REQ-82 (TS-5) licenses exactly one nondeterminism and no more: with two
// same-direction declarations, "which of the two is reported is
// unspecified". It says nothing about a model that breaches BOTH
// directions at once, where the choice is not between two offending names
// but between two contradictory REMEDIES.

import (
	"errors"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// advBothDirections is a model that breaches both reserved-key naming
// rules in one document: `[tags.recognized]` is declared `owned` (the
// author-must-rename direction) while `[tags.outcome]` declares
// provenance `recognized` (the kernel-owned direction). Both are legal
// TOML and distinct keys, so neither the duplicate-key rule (REQ-13) nor
// the lower bound (REQ-15) preempts the reserved-key check.
const advBothDirections = `
outcomes = ["go", "stop"]
terminal = ["done"]

[model]
id = "advboth"
version = 1

[tags.recognized]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[tags.outcome]
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
[rule.match.outcome]
eq = "go"
[rule.write]
status = "b"
`

// TestAdv0008_DoublyBreachingModelReportsOneStableDirection pins that a
// model breaching BOTH naming directions reports the same rule identifier
// and the same remedy on every load of the identical bytes.
//
// FAILURE MODE: a scan that iterates `decls`, a Go `map[string]TagDecl`,
// and returns on the first match leaves the payload to the runtime hash
// seed, because Go randomizes map iteration order per range. The category
// is stable — both branches carry `reserved_tag_key` — but the THREE
// FIELDS `0008:C3` adds are not: one run yields
//
//	Offending "outcome",    Remedy "recognized", Rule reserved-tag-key/kernel-owned
//
// and another yields
//
//	Offending "recognized",  Remedy "",           Rule reserved-tag-key/author-must-rename
//
// REQ-27's byte-for-byte golden assertion and REQ-33's "remediation
// lookup" are both undefined against a payload that flips. REQ-82's
// license covers two same-direction declarations, not two contradictory
// remedies. Sorting the declaration keys before the scan fixes it, which
// is why this test asserts the sorted-first payload exactly rather than
// sampling repeated loads: repeated map ranges are not guaranteed to
// expose every order, so a histogram could pass against the broken scan.
func TestAdv0008_DoublyBreachingModelReportsOneStableDirection(t *testing.T) {
	_, err := table.Load([]byte(advBothDirections), "adv-both.toml")
	if err == nil {
		t.Fatalf("the doubly-breaching model loaded clean; it breaches " +
			"both directions of `0008:C2`")
	}
	cat, ok := table.CategoryOf(err)
	if !ok || cat != table.CatReservedTagKey {
		t.Fatalf("category = %q (ok=%v); want %q — REQ-31 puts both "+
			"rule identifiers inside the one category",
			cat, ok, table.CatReservedTagKey)
	}
	var f *table.Failure
	if !errors.As(err, &f) {
		t.Fatalf("refusal is not a *table.Failure: %v", err)
	}

	// The scan runs over sorted declaration keys, so the doubly-breaching
	// model reports the lexicographically first offender: `outcome` (which
	// declares provenance `recognized` under another name) sorts before
	// `recognized` (declared `owned`). That is an exact payload, not a
	// sampled one — `0008:C3`'s rule identifier is a token "for golden
	// assertions and remediation lookup" (REQ-27, REQ-33), and the two
	// directions prescribe OPPOSITE edits, so a histogram over repeated
	// loads is the wrong oracle: Go never promises a map range exposes every
	// order, and a map-ranging implementation could pass it by luck.
	// REQ-82 licenses an unspecified choice between two SAME-direction
	// declarations; this fixture is cross-direction and is pinned exactly.
	if f.Offending != "outcome" {
		t.Errorf("Offending = %q; want %q — the sorted scan reports the "+
			"lexicographically first breaching declaration, and `0008:C3` "+
			"carries the offending name as authored (REQ-27, REQ-33)",
			f.Offending, "outcome")
	}
	if f.Remedy != "recognized" {
		t.Errorf("Remedy = %q; want %q — the kernel-owned direction's "+
			"remedy names the reserved key to rename TO; the opposite "+
			"direction's remedy is the empty string, and REQ-33's "+
			"remediation lookup is undefined if the two can swap "+
			"(REQ-27, REQ-33, REQ-82)", f.Remedy, "recognized")
	}
	if f.Rule != table.RuleKernelOwned {
		t.Errorf("Rule = %q; want %q — REQ-27 makes the rule identifier a "+
			"comparable token a golden test asserts byte-for-byte, and "+
			"REQ-82 licenses no coin-flip between contradictory remedies",
			f.Rule, table.RuleKernelOwned)
	}
}
