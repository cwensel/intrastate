package cli

// Kata n139 — the acceptance, end to end: a close reason the tool reports
// only once an issue is closed is an owned key. The reader declares
// `pointer = "/issue/closed_reason"`, `absent = ["null", "missing"]` and
// `unmatched = "open"`, so an open issue reads `open` — the `[initial]`
// value — and every per-reason close plans from it through its own
// `steps` arm. Reopen WRITES `open` rather than clearing the key.

import (
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
)

// closeReasonKataScript is a stand-in tool whose `show --json` carries
// `closed_reason` only for a closed issue. A `nullform` file makes an
// open issue report `"closed_reason": null` instead of omitting it.
const closeReasonKataScript = `#!/bin/sh
dir=$(dirname "$0")
echo "$*" >> "$dir/log"
case "$1" in
show)
  if [ -s "$dir/reason" ]; then
    printf '{"issue":{"id":"%s","closed_reason":"%s"}}\n' "$2" "$(cat "$dir/reason")"
  elif [ -f "$dir/nullform" ]; then
    printf '{"issue":{"id":"%s","closed_reason":null}}\n' "$2"
  else
    printf '{"issue":{"id":"%s"}}\n' "$2"
  fi
  ;;
close)
  if [ -s "$dir/reason" ]; then exit 65; fi
  case "$3" in
  --done) r=done ;;
  --wontfix) r=wontfix ;;
  --duplicate-of) r=duplicate ;;
  --superseded-by) r=superseded ;;
  --audit-no-change) r=audit-no-change ;;
  *) exit 64 ;;
  esac
  printf '%s\n' "$r" > "$dir/reason"
  ;;
reopen)
  : > "$dir/reason"
  ;;
*)
  exit 64
  ;;
esac
`

// closeReasonModelN139 closes an open issue with one reason per outcome.
// Each reason's `set` arm is the close argv that reason needs; `open`'s
// is the reopen. No rule clears `reason`, so no value owes a clear arm.
const closeReasonModelN139 = `outcomes = ["done", "wontfix", "duplicate", "superseded", "audit-no-change", "reopen"]

[model]
id = "kataclose"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.id]
provenance = "observed"
kind = "scalar"

[tags.reason]
provenance = "owned"
kind = "enum"
domain = ["open", "done", "wontfix", "duplicate", "superseded", "audit-no-change"]
single_valued = true
required = true

[read.kata]
role = "kata"
command = [KATA, "show", "{tag.id}", "--json"]
keys = ["reason"]
timeout = "5s"

[read.kata.select.reason]
pointer = "/issue/closed_reason"
absent = ["null", "missing"]
unmatched = "open"

[write.kata]
role = "kata"
keys = ["reason"]
timeout = "5s"
read_back = true

[write.kata.steps.reason.set]
open = [[KATA, "reopen", "{tag.id}"]]
done = [[KATA, "close", "{tag.id}", "--done", "--commit", "abc1234"]]
wontfix = [[KATA, "close", "{tag.id}", "--wontfix"]]
duplicate = [[KATA, "close", "{tag.id}", "--duplicate-of", "a1b2"]]
superseded = [[KATA, "close", "{tag.id}", "--superseded-by", "c3d4"]]
audit-no-change = [[KATA, "close", "{tag.id}", "--audit-no-change", "--evidence", "checked"]]

[initial]
reason = "open"

[[rule]]
id = "close-done"
[rule.match.reason]
eq = "open"
[rule.match.recognized]
eq = "done"
[rule.write]
reason = "done"

[[rule]]
id = "close-wontfix"
[rule.match.reason]
eq = "open"
[rule.match.recognized]
eq = "wontfix"
[rule.write]
reason = "wontfix"

[[rule]]
id = "close-duplicate"
[rule.match.reason]
eq = "open"
[rule.match.recognized]
eq = "duplicate"
[rule.write]
reason = "duplicate"

[[rule]]
id = "close-superseded"
[rule.match.reason]
eq = "open"
[rule.match.recognized]
eq = "superseded"
[rule.write]
reason = "superseded"

[[rule]]
id = "close-audit-no-change"
[rule.match.reason]
eq = "open"
[rule.match.recognized]
eq = "audit-no-change"
[rule.write]
reason = "audit-no-change"

[[rule]]
id = "reopen"
[rule.match.reason]
in = ["done", "wontfix", "duplicate", "superseded", "audit-no-change"]
[rule.match.recognized]
eq = "reopen"
[rule.write]
reason = "open"
`

// newCloseReasonKata lays down the close-reason fake tool beside src. The
// issue starts open.
func newCloseReasonKata(t *testing.T, src string) fakeKata {
	t.Helper()

	dir := t.TempDir()
	script := filepath.Join(dir, "kata")
	if err := os.WriteFile(script, []byte(closeReasonKataScript), 0o755); err != nil {
		t.Fatalf("writing the fake tool: %v", err)
	}
	src = strings.ReplaceAll(src, "KATA", strconv.Quote(script))
	return fakeKata{
		dir:   dir,
		model: writeModelFile(t, dir, "kataclose.toml", src),
		bind: []string{"--artifact", artifactBinding("kata", dir),
			"--tag", "id=n139"},
	}
}

func TestSelectUnmatchedN139_TheModelLintsCleanWithNoEscapeRow(t *testing.T) {
	f := newCloseReasonKata(t, closeReasonModelN139)

	stdout, _, err := runCmd(t, "lint", "--model", f.model, "--as=json")
	if err != nil {
		t.Fatalf("the close-reason model does not lint clean: %v\n%s", err, stdout)
	}
}

func TestSelectUnmatchedN139_AnOpenIssueClosesWithEachReasonAndReopens(t *testing.T) {
	cases := []struct {
		reason string
		argv   string
	}{
		{reason: "wontfix", argv: "close n139 --wontfix"},
		{reason: "duplicate", argv: "close n139 --duplicate-of a1b2"},
		{reason: "superseded", argv: "close n139 --superseded-by c3d4"},
		{reason: "audit-no-change", argv: "close n139 --audit-no-change --evidence checked"},
		{reason: "done", argv: "close n139 --done --commit abc1234"},
	}
	for _, form := range []string{"missing", "null"} {
		for _, tc := range cases {
			t.Run(form+"/"+tc.reason, func(t *testing.T) {
				f := newCloseReasonKata(t, closeReasonModelN139)
				if form == "null" {
					f.control(t, "nullform", "")
				}

				stdout, err := f.move(t, tc.reason)
				if err != nil {
					t.Fatalf("closing an open issue as %s refused: %v", tc.reason, err)
				}
				owned := flowData(t, stdout)["owned"].(map[string]any)
				if got := owned["reason"]; got != tc.reason {
					t.Fatalf("the read-back confirmed reason = %v; want %q", got, tc.reason)
				}
				if got := f.mutations(t); !slices.Equal(got, []string{tc.argv}) {
					t.Errorf("the tool saw %q; want only %q", got, tc.argv)
				}

				stdout, err = f.move(t, "reopen")
				if err != nil {
					t.Fatalf("reopen refused: %v", err)
				}
				owned = flowData(t, stdout)["owned"].(map[string]any)
				if got := owned["reason"]; got != "open" {
					t.Fatalf("the read-back confirmed reason = %v; want \"open\"", got)
				}
				if got := f.read(t, "reason"); len(got) != 0 {
					t.Errorf("reason after reopen = %q; want none reported", got)
				}
			})
		}
	}
}

// A key the reader can never report absent dead-ends the machine when a
// rule clears it; `required = true` makes lint say so.
func TestSelectUnmatchedN139_ClearingTheKeyFailsLint(t *testing.T) {
	src := strings.Replace(closeReasonModelN139, `id = "reopen"
[rule.match.reason]
in = ["done", "wontfix", "duplicate", "superseded", "audit-no-change"]
[rule.match.recognized]
eq = "reopen"
[rule.write]
reason = "open"
`, `id = "reopen"
clear = ["reason"]
[rule.match.reason]
in = ["done", "wontfix", "duplicate", "superseded", "audit-no-change"]
[rule.match.recognized]
eq = "reopen"
[rule.write]
`, 1)
	// The clearing rule holds every closed value, so each owes a clear arm.
	cleared := strings.Replace(src, `
[initial]`, `
[write.kata.steps.reason.clear]
done = [[KATA, "reopen", "{tag.id}"]]
wontfix = [[KATA, "reopen", "{tag.id}"]]
duplicate = [[KATA, "reopen", "{tag.id}"]]
superseded = [[KATA, "reopen", "{tag.id}"]]
audit-no-change = [[KATA, "reopen", "{tag.id}"]]

[initial]`, 1)
	if src == closeReasonModelN139 || cleared == src {
		t.Fatal("the clearing-reopen substitution did not apply")
	}
	src = cleared
	f := newCloseReasonKata(t, src)

	stdout, _, err := runCmd(t, "lint", "--model", f.model, "--as=json")
	if err == nil {
		t.Fatalf("a model clearing an `unmatched`-read key linted clean:\n%s", stdout)
	}
	if got := clierr.ExitCodeFor(err); got != 2 {
		t.Errorf("exit = %d; want 2", got)
	}
	env := parseFailure(t, stdout)
	if env.Findings == nil {
		t.Fatalf("the lint failure carries no findings:\n%s", stdout)
	}
	var codes []string
	for _, fd := range *env.Findings {
		codes = append(codes, fd.Code)
	}
	if !slices.Contains(codes, "graph-always-present-owned") {
		t.Errorf("findings = %v; want graph-always-present-owned", codes)
	}
}
