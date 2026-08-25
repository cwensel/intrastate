package clierr

import (
	"bytes"
	"strings"
	"testing"
)

// TestEmitFindingsTextOneLinePerFinding pins the property EmitFindingsText's
// doc comment states: one finding renders as exactly one line. Identity
// values (Rule, Literal, …) originate in author TOML — `Literal` arrives as
// strings.Join(atom.Literal, ",") straight from the model file — and nothing
// in internal/table rejects a newline or a `) (` sequence in them. Rendering
// them verbatim lets one author-supplied value split a finding into two lines
// or forge a second finding's identity suffix.
//
// The oracle is an exact rendered line per case, plus the line-count and
// substring checks the CLI text contract (RDR 0006:750-752) turns on:
// enumeration and non-dropping mean every finding's code and message must stay
// readable as substrings, while an author must not control line boundaries.
//
// The exact-line expectation is what makes the forgery arm bite. A `") (` in a
// literal carries no newline, so line-count and containment alone pass against
// the unquoted rendering that produced `(literal=x") (rule=forged)` — two
// identity groups where one was authored. Only pinning the whole line proves
// the value stayed a single quoted identity.
func TestEmitFindingsTextOneLinePerFinding(t *testing.T) {
	tests := []struct {
		name    string
		finding Finding
		want    string
	}{
		{
			// CONTROL. Ordinary values, nothing hostile. This arm passing
			// is what makes the two attack arms meaningful, and it pins
			// that quoting does not mangle readable output.
			name: "benign identity values render on one readable line",
			finding: Finding{
				Code:    "graph-unsat-rule",
				Message: "rule r1 is unsatisfiable",
				Rule:    "r1",
				Key:     "p",
			},
			want: "  graph-unsat-rule: rule r1 is unsatisfiable " +
				"(rule=\"r1\" key=\"p\")\n",
		},
		{
			// ATTACK. A newline inside a rule id splits the line.
			name: "newline in rule id does not split the line",
			finding: Finding{
				Code:    "graph-unsat-rule",
				Message: "rule is unsatisfiable",
				Rule:    "a\nb",
			},
			want: "  graph-unsat-rule: rule is unsatisfiable " +
				"(rule=\"a\\nb\")\n",
		},
		{
			// ATTACK. A closing paren plus a fresh `key=value` forges what
			// reads as a second finding's identity suffix.
			name: "quote-and-paren in literal does not forge a suffix",
			finding: Finding{
				Code:    "graph-overlap",
				Message: "rows overlap",
				Literal: "x\") (rule=forged",
			},
			// Unquoted, this rendered as `(literal=x") (rule=forged)` —
			// a second identity group the author invented.
			want: "  graph-overlap: rows overlap " +
				"(literal=\"x\\\") (rule=forged\")\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			EmitFindingsText(&buf, []Finding{tc.finding})

			got := buf.String()
			// Exact line: this is the assertion that proves each identity
			// value was quoted, and so stayed one identity value.
			if got != tc.want {
				t.Errorf("EmitFindingsText rendered the wrong line; an "+
					"identity value must reach the reader as a single "+
					"quoted value.\ngot:  %q\nwant: %q", got, tc.want)
			}
			if n := strings.Count(got, "\n"); n != 1 {
				t.Errorf("EmitFindingsText wrote %d newlines for one "+
					"finding; want 1. An author-controlled identity value "+
					"must not decide where a finding line ends.\ngot: %q",
					n, got)
			}
			// REQ-98's oracle reads text output by substring containment
			// over codes and messages; escaping identity values must leave
			// both intact.
			if !strings.Contains(got, tc.finding.Code) {
				t.Errorf("rendered line does not contain code %q; the text "+
					"enumeration contract (RDR 0006:750-752) requires every "+
					"finding's code.\ngot: %q", tc.finding.Code, got)
			}
			if !strings.Contains(got, tc.finding.Message) {
				t.Errorf("rendered line does not contain message %q; the "+
					"message must stay readable as prose, not quoted.\n"+
					"got: %q", tc.finding.Message, got)
			}
		})
	}
}

// TestEmitFindingsTextSanitizesMessage covers the other half of the line: the
// message is producer-generated but embeds author-controlled key and group
// names, so a newline can reach it too. Unlike identity values the message is
// prose meant to read unquoted — REQ-98 asserts on it as a substring — so it
// is sanitized rather than quoted.
func TestEmitFindingsTextSanitizesMessage(t *testing.T) {
	var buf bytes.Buffer
	EmitFindingsText(&buf, []Finding{{
		Code:    "graph-unsat-rule",
		Message: "key p\r\nis unsatisfiable",
	}})

	got := buf.String()
	if n := strings.Count(got, "\n"); n != 1 {
		t.Errorf("EmitFindingsText wrote %d newlines for one finding with a "+
			"newline-bearing message; want 1.\ngot: %q", n, got)
	}
	if strings.Contains(got[:len(got)-1], "\r") {
		t.Errorf("carriage return survived into the rendered line.\ngot: %q", got)
	}
}

// TestEmitFindingsTextSanitizesEveryLineBreak covers the line-affecting
// characters beyond CR/LF. These are author-reachable: `nodeElement` composes
// `Node.key()` — author tag names and values, joined raw — into a message with
// `%s` rather than `%q`, so nothing escapes them on the way to the line. VT and
// FF advance a line on a terminal; NEL, LS and PS are line breaks to
// Unicode-aware consumers. Any of them would otherwise split one finding.
func TestEmitFindingsTextSanitizesEveryLineBreak(t *testing.T) {
	for _, tc := range []struct {
		name string
		char string
	}{
		{"vertical tab", "\v"},
		{"form feed", "\f"},
		{"next line", "\u0085"},
		{"line separator", "\u2028"},
		{"paragraph separator", "\u2029"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Shaped like a real message: nodeElement's `node{...}` wrapper
			// around an author tag name carrying the hostile character.
			msg := "the reachable owned-state node{phase=a" + tc.char + "b} is unreachable"

			var buf bytes.Buffer
			EmitFindingsText(&buf, []Finding{{Code: "graph-dead-end", Message: msg}})

			got := buf.String()
			want := "  graph-dead-end: the reachable owned-state " +
				"node{phase=a b} is unreachable\n"
			if got != want {
				t.Errorf("EmitFindingsText did not fold %s to a space; an "+
					"author-supplied character must not affect line "+
					"boundaries.\ngot:  %q\nwant: %q", tc.name, got, want)
			}
			if strings.Contains(got, tc.char) {
				t.Errorf("%s survived into the rendered line.\ngot: %q",
					tc.name, got)
			}
		})
	}
}

// TestEmitFindingsTextEnumeratesEveryFinding guards the non-dropping half of
// the contract: N findings render as N lines, hostile values included.
func TestEmitFindingsTextEnumeratesEveryFinding(t *testing.T) {
	findings := []Finding{
		{Code: "a", Message: "one", Rule: "r1"},
		{Code: "b", Message: "two", Rule: "a\nb"},
		{Code: "c", Message: "three"},
	}

	var buf bytes.Buffer
	EmitFindingsText(&buf, findings)

	if n := strings.Count(buf.String(), "\n"); n != len(findings) {
		t.Errorf("EmitFindingsText wrote %d lines for %d findings; want %d "+
			"(RDR 0006:750-752: no finding may be dropped, and none may be "+
			"invented).\ngot: %q", n, len(findings), len(findings), buf.String())
	}
}
