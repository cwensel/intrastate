package cli

// Kata 9xjf — the acceptance, end to end: an issue carrying no
// `lifecycle:*` label is claimed through `flow resolve | flow set-state`
// with no seed step and no escape row. The reader declares
// `unmatched = "filed"`, so the unlabelled issue reads `filed` — the
// `[initial]` value — and the key is always present, as `required = true`
// declares. Release WRITES `filed` rather than clearing the key.
//
// The tool is q14r's fake `kata` (`fakeKataScript`).

import (
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
)

// unmatchedModel9xjf is the kata-flight lifecycle with `filed` as the
// projection of "no lifecycle label". `filed` declares no `clear` arm —
// there is no label to remove — so the claim stays the first step.
const unmatchedModel9xjf = `outcomes = ["claim", "advance", "release"]
terminal = ["done"]

[model]
id = "kataflight"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.id]
provenance = "observed"
kind = "scalar"

[tags.session]
provenance = "observed"
kind = "scalar"

[tags.lifecycle]
provenance = "owned"
kind = "enum"
domain = ["filed", "resolving", "refining", "closed"]
single_valued = true
required = true

[read.kata]
role = "kata"
command = [KATA, "show", "{tag.id}", "--json"]
keys = ["lifecycle"]
timeout = "5s"

[read.kata.select.lifecycle]
pointer = "/labels"
element = "/label"
prefix = "lifecycle:"
unmatched = "filed"

[write.kata]
role = "kata"
keys = ["lifecycle"]
timeout = "5s"
read_back = true

[write.kata.steps.lifecycle.set]
filed = [[KATA, "unassign", "{tag.id}"]]
resolving = [[KATA, "claim", "{tag.id}", "--as", "{tag.session}"], [KATA, "label", "add", "{tag.id}", "lifecycle:resolving"]]
refining = [[KATA, "label", "add", "{tag.id}", "lifecycle:refining"]]
closed = [[KATA, "unassign", "{tag.id}"], [KATA, "label", "add", "{tag.id}", "lifecycle:closed"]]

[write.kata.steps.lifecycle.clear]
resolving = [[KATA, "label", "rm", "{tag.id}", "lifecycle:resolving"]]
refining = [[KATA, "label", "rm", "{tag.id}", "lifecycle:refining"]]

[initial]
lifecycle = "filed"

[context.done]
[context.done.match.lifecycle]
eq = "closed"

[[rule]]
id = "claim"
[rule.match.lifecycle]
eq = "filed"
[rule.match.recognized]
eq = "claim"
[rule.write]
lifecycle = "resolving"

[[rule]]
id = "refine"
[rule.match.lifecycle]
eq = "resolving"
[rule.match.recognized]
eq = "advance"
[rule.write]
lifecycle = "refining"

[[rule]]
id = "close"
[rule.match.lifecycle]
eq = "refining"
[rule.match.recognized]
eq = "advance"
[rule.write]
lifecycle = "closed"

[[rule]]
id = "release"
[rule.match.lifecycle]
in = ["resolving", "refining"]
[rule.match.recognized]
eq = "release"
[rule.write]
lifecycle = "filed"
`

// newUnmatchedKata lays down the fake tool beside src, seeded with the
// given labels plus the unrelated `area:cli`.
func newUnmatchedKata(t *testing.T, src string, labels ...string) fakeKata {
	t.Helper()

	dir := t.TempDir()
	script := filepath.Join(dir, "kata")
	if err := os.WriteFile(script, []byte(fakeKataScript), 0o755); err != nil {
		t.Fatalf("writing the fake tool: %v", err)
	}
	seed := strings.Join(append([]string{"area:cli"}, labels...), "\n") + "\n"
	if err := os.WriteFile(filepath.Join(dir, "labels"), []byte(seed), 0o600); err != nil {
		t.Fatalf("seeding labels: %v", err)
	}
	src = strings.ReplaceAll(src, "KATA", strconv.Quote(script))
	return fakeKata{
		dir:   dir,
		model: writeModelFile(t, dir, "kataflight.toml", src),
		bind: []string{"--artifact", artifactBinding("kata", dir),
			"--tag", "id=9xjf", "--tag", "session=kata-ship/1"},
	}
}

func TestSelectUnmatched9xjf_TheModelLintsCleanWithNoEscapeRow(t *testing.T) {
	f := newUnmatchedKata(t, unmatchedModel9xjf)

	stdout, _, err := runCmd(t, "lint", "--model", f.model, "--as=json")
	if err != nil {
		t.Fatalf("the unmatched lifecycle model does not lint clean: %v\n%s",
			err, stdout)
	}
}

func TestSelectUnmatched9xjf_AnUnlabelledIssueIsClaimedAndReleased(t *testing.T) {
	f := newUnmatchedKata(t, unmatchedModel9xjf)

	stdout, err := f.move(t, "claim")
	if err != nil {
		t.Fatalf("claiming an unlabelled issue refused: %v", err)
	}
	if got := flowData(t, stdout)["owned"].(map[string]any)["lifecycle"]; got != "resolving" {
		t.Fatalf("the read-back confirmed lifecycle = %v; want \"resolving\"", got)
	}
	want := []string{
		"claim 9xjf --as kata-ship/1",
		"label add 9xjf lifecycle:resolving",
	}
	if got := f.mutations(t); !slices.Equal(got, want) {
		t.Errorf("the tool saw %q; want %q — no seed step, the claim first", got, want)
	}

	stdout, err = f.move(t, "release")
	if err != nil {
		t.Fatalf("release refused: %v", err)
	}
	if got := flowData(t, stdout)["owned"].(map[string]any)["lifecycle"]; got != "filed" {
		t.Fatalf("the read-back confirmed lifecycle = %v; want \"filed\"", got)
	}
	if got := f.read(t, "labels"); !slices.Equal(got, []string{"area:cli"}) {
		t.Errorf("labels after release = %q; want no lifecycle label", got)
	}
	if got := f.read(t, "owner"); len(got) != 0 {
		t.Errorf("owner after release = %q; want unassigned", got)
	}

	// Released, the issue is claimable again.
	if _, err := f.move(t, "claim"); err != nil {
		t.Fatalf("re-claiming a released issue refused: %v", err)
	}
}

// A key the reader projects can never read absent, so a rule clearing it
// dead-ends the machine; `required = true` makes lint say so.
func TestSelectUnmatched9xjf_ClearingTheKeyFailsLint(t *testing.T) {
	src := strings.Replace(unmatchedModel9xjf, `id = "release"
[rule.match.lifecycle]
in = ["resolving", "refining"]
[rule.match.recognized]
eq = "release"
[rule.write]
lifecycle = "filed"
`, `id = "release"
clear = ["lifecycle"]
[rule.match.lifecycle]
in = ["resolving", "refining"]
[rule.match.recognized]
eq = "release"
[rule.write]
`, 1)
	if src == unmatchedModel9xjf {
		t.Fatal("the clearing-release substitution did not apply")
	}
	f := newUnmatchedKata(t, src)

	stdout, _, err := runCmd(t, "lint", "--model", f.model, "--as=json")
	if err == nil {
		t.Fatalf("a model clearing an `unmatched`-projected key linted clean:\n%s", stdout)
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
