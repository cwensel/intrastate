package table_test

// RDR 0024 — shared fixture corpus for the declared-emit-vocabulary suite.
//
// Every model here is authored as a complete TOML document and handed to
// the REAL `table.Load`. Nothing mocks the loader, the grammar checks, or
// the carrier.
//
// The corpus builds on RDR 0010's decision-table header (`dtHeader` in
// `fixtures_0010_test.go`) because `0024:C1`'s `[emit]` table is a
// TOP-LEVEL sibling of `[tags]`, and the 0010 fixtures already carry the
// two discriminating observed dimensions plus `[rule.emit]` blocks the
// cross-checks read.

import (
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// emitDeclTable renders a top-level `[emit]` block from already-authored
// TOML body text. The body is authored verbatim so a fixture can express
// shapes a typed builder could not — a non-array `domain`, nesting below
// the disposition level, a non-string member.
func emitDeclTable(body string) string {
	if body == "" {
		return "\n[emit]\n"
	}
	return "\n" + strings.TrimRight(body, "\n") + "\n"
}

// dtWithEmitDecl prefixes RDR 0010's complete decision table with an
// `[emit]` declaration body.
//
// The declaration is authored BEFORE the rules so the `[[rule]]` array
// items keep their document positions: a top-level table appearing after
// an array-of-tables would be parsed into the last rule.
func dtWithEmitDecl(body string) string {
	return dtHeader + emitDeclTable(body) + dtRulesOnly()
}

// dtRulesOnly is `dtComplete` minus the header — the four ordinary rules
// of the 0010 decision table, three emitting `verdict` and the fourth
// authoring no emit block.
func dtRulesOnly() string {
	return strings.TrimPrefix(dtComplete, dtHeader)
}

// declVerdictEnum is the well-formed declaration the positive fixtures
// share: the `verdict` key declared `enum` over the three values the
// 0010 rules actually author, partitioned into two model-authored
// dispositions.
//
// `0024:C1` fixes that intrastate owns no disposition vocabulary, so the
// tokens here (`route`, `stop`) are the record's own examples and carry no
// meaning to the loader.
const declVerdictEnum = `[emit.verdict]
kind = "enum"
[emit.verdict.domain]
route = ["alpha", "beta"]
stop = ["gamma"]`

// declVerdictFlat is the same key declared with the FLAT array spelling of
// the same `domain` key — no dispositions.
const declVerdictFlat = `[emit.verdict]
kind = "enum"
domain = ["alpha", "beta", "gamma"]`

// loadEmitDecls0024 loads a model carrying an `[emit]` declaration body
// and returns the normalized model.
func loadEmitDecls0024(t *testing.T, body, id string) *table.Model {
	t.Helper()

	return loadSource(t, dtWithEmitDecl(body), id)
}

// refuseEmitDecls0024 loads a model carrying an `[emit]` declaration body
// that must refuse, and returns the categorized failure.
func refuseEmitDecls0024(t *testing.T, body, id string) *table.Failure {
	t.Helper()

	return refuseSource(t, dtWithEmitDecl(body), id)
}
