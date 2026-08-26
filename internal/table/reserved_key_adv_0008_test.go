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
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/table"
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
// FAILURE MODE: `internal/table/load.go:165` iterates `decls`, a Go
// `map[string]TagDecl`, and returns on the first match. Go randomizes map
// iteration order per range, so which of the two branches fires is decided
// by the runtime hash seed. The category is stable — both branches carry
// `reserved_tag_key` — but the THREE FIELDS `0008:C3` adds are not: one
// run yields
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
// remedies. Sorting the declaration keys before the scan fixes it.
func TestAdv0008_DoublyBreachingModelReportsOneStableDirection(t *testing.T) {
	const runs = 200

	type payload struct{ offending, remedy, rule string }
	seen := map[payload]int{}

	for range runs {
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
		seen[payload{f.Offending, f.Remedy, f.Rule}]++
	}

	if len(seen) > 1 {
		var lines []string
		for p, n := range seen {
			lines = append(lines, "offending="+p.offending+
				" remedy="+p.remedy+" rule="+p.rule+
				" hits="+itoaAdv(n))
		}
		t.Errorf("the SAME bytes yielded %d different reserved_tag_key "+
			"payloads across %d loads; `0008:C3`'s rule identifier is a "+
			"token \"for golden assertions and remediation lookup\" "+
			"(REQ-27, REQ-33) and the two directions prescribe OPPOSITE "+
			"edits — REQ-82 licenses an unspecified choice between two "+
			"same-direction declarations, not a coin-flip between "+
			"contradictory remedies. load.go:165 ranges a Go map.\n  %s",
			len(seen), runs, strings.Join(lines, "\n  "))
	}
}

func itoaAdv(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
