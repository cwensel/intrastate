Model: claude-sonnet-5

# A2 verification spike: widened interpreter predicate causes zero new refusals

## Claim under test

"No binding that lints green today carries a listed interpreter basename as
a non-argv0 word followed later by one of that interpreter's flags — i.e.
widening the interpreter predicate from argv0-after-an-`env`-walk to ANY
argv position refuses NOTHING that is currently accepted in the repo's
fixtures or the 0025 verification models."

Expected result: ZERO new refusals.

## Method

Population sweep of every `command = [...]` argv vector in:
- `internal/table/testdata/neg/*.toml` (6 fixtures)
- `internal/table/command_carrier_0025_test.go` (all inline argv literals
  used across its test cases)

No fixtures under `internal/table/testdata/pos/` declare a `command` field
(none exist — `grep -rn 'command' internal/table/testdata/pos` returned
nothing).

For each argv vector, computed both the CURRENT predicate
(`interpreterForm`, argv0-after-env-walk) and the PROPOSED widened
predicate (`widenedInterpreterForm`, a plain two-index scan: lowest i where
`filepathBase(argv[i])` is a `shellInterpreters` key, lowest j>i where
`argv[j]` is one of that key's flags — nothing before i read) via a
throwaway test file `internal/table/zz_a2_spike_test.go` in package
`table`. The file was deleted after the run; no committed source was
modified.

## Commands run

```
go test ./internal/table/ -run TestA2Spike -v
go test ./internal/table/
git status --porcelain internal/table/
```

## Spike test source (deleted after run)

```go
package table

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

// widenedInterpreterForm is the PROPOSED (A2) predicate: a two-index scan
// for the lowest i where filepathBase(argv[i]) is a shellInterpreters key,
// and the lowest j>i where argv[j] is one of that key's flags. Nothing
// before i is read. This is a throwaway spike copy; it is NOT wired into
// production code.
func widenedInterpreterForm(argv []string) (string, bool) {
	for i := 0; i < len(argv); i++ {
		flags, known := shellInterpreters[filepathBase(argv[i])]
		if !known {
			continue
		}
		for j := i + 1; j < len(argv); j++ {
			if slices.Contains(flags, argv[j]) {
				return argv[i] + " " + argv[j], true
			}
		}
	}
	return "", false
}

func TestA2Spike(t *testing.T) {
	type row struct {
		label string
		argv  []string
		green bool // true if this fixture is expected to lint green (accepted) today
	}

	cases := []row{
		// --- neg-command-shell-interpreter.toml (already-negative) ---
		{"neg-command-shell-interpreter.toml", []string{"sh", "-c", "rdr-gate"}, false},

		// --- command_carrier_0025_test.go: TestReq74 (already-negative) ---
		{"Req74 sh -c", []string{"sh", "-c", "cat {artifact}"}, false},
		{"Req74 bash -c", []string{"bash", "-c", "echo hi"}, false},
		{"Req74 python -c", []string{"python", "-c", "print(1)"}, false},
		{"Req74 env chain to sh -c", []string{"env", "sh", "-c", "echo hi"}, false},

		// --- command_carrier_0025_test.go: TestReq83 (GREEN: unlisted interpreter) ---
		{"Req83 perl -e (unlisted, green)", []string{"perl", "-e", "print 1"}, true},

		// --- command_carrier_0025_test.go: interpreter-outranks-others fixtures (already-negative) ---
		{"shell interpreter outranks placeholder", []string{"sh", "-c", "{role}"}, false},
		{"shell interpreter outranks output shape", []string{"sh", "-c", "cat"}, false},

		// --- command_carrier_0025_test.go: placeholder-clause non-interpreter-argv0 cases (GREEN except where noted) ---
		{"known token whole element", []string{"reader", "{artifact}"}, true},
		{"unknown token whole element", []string{"reader", "{role}"}, false}, // refused: unknown placeholder, not interpreter-related
		{"known token embedded mid-string", []string{"reader", "--file={artifact}"}, false},
		{"known token with a suffix", []string{"reader", "{artifact}.bak"}, false},
		{"partial brace token", []string{"reader", "{artifact"}, false},
		{"known token with a prefix word under non-interpreter argv0", []string{"reader", "prefix {artifact}"}, false},
		{"unknown token with a prefix word under non-interpreter argv0", []string{"echo", "hello {role}"}, false},
		{"known token with a trailing word under non-interpreter argv0", []string{"reader", "{artifact} tail"}, false},

		// --- command_carrier_0025_test.go: other GREEN argv across the suite ---
		{"git config carrier (green)", []string{"git", "config", "--file", "{artifact}", "--get", "flow.status"}, true},
		{"tools/flowstate-write (green)", []string{"tools/flowstate-write", "{artifact}"}, true},
		{"git -C diff --quiet (green)", []string{"git", "-C", "{artifact}", "diff", "--quiet"}, true},
		{"reader bare argv0 (green)", []string{"reader"}, true},
		{"gater bare argv0 (green)", []string{"gater"}, true},
		{"writer bare argv0 (green)", []string{"writer"}, true},
		{"gater with role token (unknown placeholder, negative)", []string{"gater", "{role}"}, false},
		{"reader with token flag (green)", []string{"reader", "--token", "ghp_0123456789abcdefghijklmnopqrstuvwxyz"}, true},

		// --- neg-command-empty.toml / other neg fixtures (already-negative, non-interpreter) ---
		{"neg-command-empty.toml", []string{}, false},
		{"neg-command-output-shape.toml", []string{"rdr-read", "{artifact}"}, false},
		{"neg-command-and-path-conflict.toml", []string{"rdr-gate", "{artifact}"}, false},
		{"neg-command-env-conflict.toml", []string{"rdr-gate", "{artifact}"}, false},
		{"neg-command-unknown-placeholder.toml", []string{"rdr-gate", "{role}"}, false},

		// --- empty / malformed element edge fixtures (already-negative) ---
		{"empty element with flag-like word", []string{"", "--get"}, false},
		{"reader with empty element", []string{"reader", ""}, false},
		{"reader with empty then role", []string{"reader", "", "{role}"}, false},
	}

	fmt.Println("argv | green | old(bool,form) | new(bool,form) | DIVERGES?")
	fmt.Println("---")

	divergences := 0
	for _, c := range cases {
		oldForm, oldOK := interpreterForm(c.argv)
		newForm, newOK := widenedInterpreterForm(c.argv)
		diverges := oldOK != newOK
		flag := ""
		if diverges {
			flag = "<<< DIVERGES"
			if !oldOK && newOK {
				flag += " (NEW REFUSAL)"
			}
			divergences++
		}
		fmt.Printf("%-70q | green=%-5v | old=(%v,%q) | new=(%v,%q) | %s  [%s]\n",
			c.argv, c.green, oldOK, oldForm, newOK, newForm, flag, c.label)

		// The finding that matters: a GREEN fixture that the widened
		// predicate would newly refuse.
		if c.green && !oldOK && newOK {
			t.Errorf("NEW REFUSAL on a green fixture %q: argv=%#v old=(%v,%q) new=(%v,%q)",
				c.label, c.argv, oldOK, oldForm, newOK, newForm)
		}
	}

	fmt.Println("---")
	fmt.Printf("total cases: %d, divergences: %d\n", len(cases), divergences)

	if strings.Contains("", "unused") { // keep strings import tidy if trimmed later
		return
	}
}
```

## Full output

```
=== RUN   TestA2Spike
argv | green | old(bool,form) | new(bool,form) | DIVERGES?
---
["sh" "-c" "rdr-gate"] | green=false | old=(true,"sh -c") | new=(true,"sh -c") |   [neg-command-shell-interpreter.toml]
["sh" "-c" "cat {artifact}"] | green=false | old=(true,"sh -c") | new=(true,"sh -c") |   [Req74 sh -c]
["bash" "-c" "echo hi"] | green=false | old=(true,"bash -c") | new=(true,"bash -c") |   [Req74 bash -c]
["python" "-c" "print(1)"] | green=false | old=(true,"python -c") | new=(true,"python -c") |   [Req74 python -c]
["env" "sh" "-c" "echo hi"] | green=false | old=(true,"sh -c") | new=(true,"sh -c") |   [Req74 env chain to sh -c]
["perl" "-e" "print 1"] | green=true  | old=(false,"") | new=(false,"") |   [Req83 perl -e (unlisted, green)]
["sh" "-c" "{role}"] | green=false | old=(true,"sh -c") | new=(true,"sh -c") |   [shell interpreter outranks placeholder]
["sh" "-c" "cat"] | green=false | old=(true,"sh -c") | new=(true,"sh -c") |   [shell interpreter outranks output shape]
["reader" "{artifact}"] | green=true  | old=(false,"") | new=(false,"") |   [known token whole element]
["reader" "{role}"] | green=false | old=(false,"") | new=(false,"") |   [unknown token whole element]
["reader" "--file={artifact}"] | green=false | old=(false,"") | new=(false,"") |   [known token embedded mid-string]
["reader" "{artifact}.bak"] | green=false | old=(false,"") | new=(false,"") |   [known token with a suffix]
["reader" "{artifact"] | green=false | old=(false,"") | new=(false,"") |   [partial brace token]
["reader" "prefix {artifact}"] | green=false | old=(false,"") | new=(false,"") |   [known token with a prefix word under non-interpreter argv0]
["echo" "hello {role}"] | green=false | old=(false,"") | new=(false,"") |   [unknown token with a prefix word under non-interpreter argv0]
["reader" "{artifact} tail"] | green=false | old=(false,"") | new=(false,"") |   [known token with a trailing word under non-interpreter argv0]
["git" "config" "--file" "{artifact}" "--get" "flow.status"] | green=true  | old=(false,"") | new=(false,"") |   [git config carrier (green)]
["tools/flowstate-write" "{artifact}"] | green=true  | old=(false,"") | new=(false,"") |   [tools/flowstate-write (green)]
["git" "-C" "{artifact}" "diff" "--quiet"] | green=true  | old=(false,"") | new=(false,"") |   [git -C diff --quiet (green)]
["reader"] | green=true  | old=(false,"") | new=(false,"") |   [reader bare argv0 (green)]
["gater"] | green=true  | old=(false,"") | new=(false,"") |   [gater bare argv0 (green)]
["writer"] | green=true  | old=(false,"") | new=(false,"") |   [writer bare argv0 (green)]
["gater" "{role}"] | green=false | old=(false,"") | new=(false,"") |   [gater with role token (unknown placeholder, negative)]
["reader" "--token" "ghp_0123456789abcdefghijklmnopqrstuvwxyz"] | green=true  | old=(false,"") | new=(false,"") |   [reader with token flag (green)]
[] | green=false | old=(false,"") | new=(false,"") |   [neg-command-empty.toml]
["rdr-read" "{artifact}"] | green=false | old=(false,"") | new=(false,"") |   [neg-command-output-shape.toml]
["rdr-gate" "{artifact}"] | green=false | old=(false,"") | new=(false,"") |   [neg-command-and-path-conflict.toml]
["rdr-gate" "{artifact}"] | green=false | old=(false,"") | new=(false,"") |   [neg-command-env-conflict.toml]
["rdr-gate" "{role}"] | green=false | old=(false,"") | new=(false,"") |   [neg-command-unknown-placeholder.toml]
["" "--get"] | green=false | old=(false,"") | new=(false,"") |   [empty element with flag-like word]
["reader" ""] | green=false | old=(false,"") | new=(false,"") |   [reader with empty element]
["reader" "" "{role}"] | green=false | old=(false,"") | new=(false,"") |   [reader with empty then role]
---
total cases: 32, divergences: 0
--- PASS: TestA2Spike (0.00s)
PASS
ok  	github.com/cwensel/intrastate/internal/table	0.297s
```

(Argv columns above are reflowed to single-line per row for readability;
the raw terminal output padded each `%q` element to a fixed column width,
content is otherwise identical.)

## Baseline: full `internal/table` suite before and after the spike

```
$ go test ./internal/table/
ok  	github.com/cwensel/intrastate/internal/table	0.542s
```

Run both before adding `zz_a2_spike_test.go` and again after deleting it —
green both times.

## Population swept

32 argv vectors total:
- 6 from `internal/table/testdata/neg/*.toml` (all already-negative
  fixtures; none are green).
- 26 argv literals drawn from `internal/table/command_carrier_0025_test.go`
  inline test cases, spanning: the C5 interpreter deny-list tests
  (`TestReq74`, `TestReq83`), the C2 placeholder-clause non-interpreter-argv0
  cases (the ones most structurally adjacent to A2, since they already probe
  "known token as a non-whole-element word"), the C1 carrier decode tests,
  the C4 precedence-ordering tests, and assorted carrier/env/rejected-spelling
  tests.
- No fixtures exist under `internal/table/testdata/pos/` with a `command`
  field — confirmed via `grep -rn command internal/table/testdata/pos`
  returning empty. All green `command` cases live in the Go test file, not
  TOML positives.

Of the 32, 11 are green/accepted today: `perl -e` (unlisted interpreter),
`reader {artifact}`, `git config ... {artifact} ... flow.status`,
`tools/flowstate-write {artifact}`, `git -C {artifact} diff --quiet`, bare
`reader`/`gater`/`writer`, and `reader --token ghp_...`.

## Finding

Zero divergences (old != new) across all 32 argv vectors, and in
particular zero cases where old=false and new=true among the 11 green
cases. The widened predicate refuses nothing that lints green today.

The one case structurally closest to a divergence, `perl -e` (green,
unlisted interpreter), stays `(false, "")` under both predicates because
`perl` is not a `shellInterpreters` key — the widen is over argv
*position*, not over which basenames are listed, so an unlisted
interpreter is unaffected regardless of scan strategy.

No fixture in the swept population places a LISTED interpreter basename
(sh/bash/dash/ksh/zsh/csh/tcsh/python/ruby/node/php) as a non-argv0 word
with one of its own flags appearing later in the same vector while still
being green — the only argv vectors containing a listed interpreter
basename anywhere are the ones already refused today (`TestReq74`, the
`neg-command-shell-interpreter.toml` fixture, and the two
precedence-ordering fixtures), all of which stay refused (old=true,
new=true) under the widened predicate.

## Verdict

PASS — assumption A2 holds against the full argv population found in this
repo's fixtures and 0025 verification test file. No committed source was
modified; `internal/table/zz_a2_spike_test.go` was created, run, and
deleted as part of this spike.
