package cli

// Kata q14r — the APPLY half of the `steps` write carrier, end to end:
// `flow resolve | flow set-state --plan -` against an argv-only tool, with
// the tool's argv declared in the model and no wrapper script. The load
// half is `internal/table/steps_carrier_q14r_test.go`.
//
// The tool is a fake `kata`: a shell script holding an issue's labels and
// owner in files beside itself, logging every argv it receives, and
// answering `show --json` in the nested shape a `select` reader projects.
// Three control files make it misbehave the ways the real tool can: a
// claim that loses its compare-and-set race, a label add that exits 0
// without persisting, and a label add that fails outright.

import (
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/cli/cmdbind"
	"github.com/cwensel/intrastate/internal/cli/flowbind"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// fakeKataScript is the stand-in tool. `claim-exit` makes `claim` exit
// with its content; `add-exit` does the same for `label add`; `nopersist`
// makes `label add` exit 0 and change nothing.
const fakeKataScript = `#!/bin/sh
dir=$(dirname "$0")
echo "$*" >> "$dir/log"
case "$1" in
show)
  owner=null
  if [ -s "$dir/owner" ]; then owner="\"$(cat "$dir/owner")\""; fi
  printf '{"issue":{"owner":%s},"labels":[' "$owner"
  sep=""
  while IFS= read -r l; do
    [ -n "$l" ] || continue
    printf '%s{"label":"%s"}' "$sep" "$l"
    sep=","
  done < "$dir/labels"
  printf ']}\n'
  ;;
claim)
  if [ -f "$dir/claim-exit" ]; then exit "$(cat "$dir/claim-exit")"; fi
  printf '%s\n' "$4" > "$dir/owner"
  ;;
unassign)
  : > "$dir/owner"
  ;;
label)
  case "$2" in
  add)
    if [ -f "$dir/add-exit" ]; then exit "$(cat "$dir/add-exit")"; fi
    if [ -f "$dir/nopersist" ]; then exit 0; fi
    grep -qxF "$4" "$dir/labels" || printf '%s\n' "$4" >> "$dir/labels"
    ;;
  rm)
    grep -vxF "$4" "$dir/labels" > "$dir/labels.tmp"
    mv "$dir/labels.tmp" "$dir/labels"
    ;;
  esac
  ;;
*)
  exit 64
  ;;
esac
`

// stepsModelQ14r is the kata-flight lifecycle. The claim — a
// compare-and-set — is the FIRST step of `set.resolving`, and `queued`
// declares no clear arm, so entering `resolving` runs the claim before any
// label changes. Every later move removes the old label, then adds the
// new one; `closed` unassigns first.
const stepsModelQ14r = `outcomes = ["claim", "advance"]
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
domain = ["queued", "resolving", "refining", "shipping", "closed"]
single_valued = true

[read.kata]
role = "kata"
command = [KATA, "show", "{tag.id}", "--json"]
keys = ["lifecycle"]
timeout = "5s"

[read.kata.select.lifecycle]
pointer = "/labels"
element = "/label"
prefix = "lifecycle:"

[write.kata]
role = "kata"
keys = ["lifecycle"]
timeout = "5s"
read_back = true

[write.kata.steps.lifecycle.set]
queued = [[KATA, "label", "add", "{tag.id}", "lifecycle:queued"]]
resolving = [[KATA, "claim", "{tag.id}", "--as", "{tag.session}"], [KATA, "label", "rm", "{tag.id}", "lifecycle:queued"], [KATA, "label", "add", "{tag.id}", "lifecycle:resolving"]]
refining = [[KATA, "label", "add", "{tag.id}", "lifecycle:refining"]]
shipping = [[KATA, "label", "add", "{tag.id}", "lifecycle:shipping"]]
closed = [[KATA, "unassign", "{tag.id}"], [KATA, "label", "add", "{tag.id}", "lifecycle:closed"]]

[write.kata.steps.lifecycle.clear]
resolving = [[KATA, "label", "rm", "{tag.id}", "lifecycle:resolving"]]
refining = [[KATA, "label", "rm", "{tag.id}", "lifecycle:refining"]]
shipping = [[KATA, "label", "rm", "{tag.id}", "lifecycle:shipping"]]

[context.done]
[context.done.match.lifecycle]
eq = "closed"

[[rule]]
id = "claim"
[rule.match.lifecycle]
eq = "queued"
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
id = "ship"
[rule.match.lifecycle]
eq = "refining"
[rule.match.recognized]
eq = "advance"
[rule.write]
lifecycle = "shipping"

[[rule]]
id = "close"
[rule.match.lifecycle]
eq = "shipping"
[rule.match.recognized]
eq = "advance"
[rule.write]
lifecycle = "closed"
`

// fakeKata is one fixture: the model, the tool's state directory, and the
// argv every invocation shares.
type fakeKata struct {
	dir   string
	model string
	bind  []string
}

// newFakeKata lays down the script and a model naming it, seeded with
// the given labels. `area:cli` rides along in every seed as the label no
// step may touch.
func newFakeKata(t *testing.T, labels ...string) fakeKata {
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
	src := strings.ReplaceAll(stepsModelQ14r, "KATA", strconv.Quote(script))
	return fakeKata{
		dir:   dir,
		model: writeModelFile(t, dir, "kataflight.toml", src),
		bind: []string{"--artifact", artifactBinding("kata", dir),
			"--tag", "id=q14r", "--tag", "session=kata-ship/1"},
	}
}

// control drops a misbehaviour file beside the tool.
func (f fakeKata) control(t *testing.T, name, body string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(f.dir, name), []byte(body), 0o600); err != nil {
		t.Fatalf("writing %s: %v", name, err)
	}
}

// read returns one state file's lines, blank lines dropped.
func (f fakeKata) read(t *testing.T, name string) []string {
	t.Helper()

	b, err := os.ReadFile(filepath.Join(f.dir, name))
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("reading %s: %v", name, err)
	}
	var out []string
	for line := range strings.SplitSeq(string(b), "\n") {
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

// mutations is the tool's argv log minus its reads: every state change
// the steps made, in the order they ran.
func (f fakeKata) mutations(t *testing.T) []string {
	t.Helper()

	var out []string
	for _, line := range f.read(t, "log") {
		if !strings.HasPrefix(line, "show ") {
			out = append(out, line)
		}
	}
	return out
}

// move is the kata's acceptance pipeline: resolve the outcome against the
// tool's live state, then pipe the envelope straight into set-state.
func (f fakeKata) move(t *testing.T, outcome string) (string, error) {
	t.Helper()

	args := append([]string{"flow", "resolve", "--model", f.model,
		"--allow-commands", "--outcome", outcome, "--as=json"}, f.bind...)
	envelope := strings.TrimSpace(requireSuccess(t, args...))

	set := append([]string{"flow", "set-state", "--model", f.model,
		"--allow-commands", "--plan", "-", "--as=json"}, f.bind...)
	stdout, _, err := runCmdStdin(t, envelope, set...)
	return stdout, err
}

func TestStepsQ14r_TheLifecycleMovesThroughEveryStateWithNoWrapper(t *testing.T) {
	f := newFakeKata(t, "lifecycle:queued")

	for _, step := range []struct {
		outcome, want string
	}{
		{"claim", "resolving"},
		{"advance", "refining"},
		{"advance", "shipping"},
		{"advance", "closed"},
	} {
		stdout, err := f.move(t, step.outcome)
		if err != nil {
			t.Fatalf("%s → %s refused: %v", step.outcome, step.want, err)
		}
		if got := flowData(t, stdout)["owned"].(map[string]any)["lifecycle"]; got != step.want {
			t.Fatalf("the read-back confirmed lifecycle = %v; want %q", got, step.want)
		}
		want := []string{"area:cli", "lifecycle:" + step.want}
		if got := f.read(t, "labels"); !slices.Equal(got, want) {
			t.Fatalf("labels after %s = %q; want %q — exactly one lifecycle "+
				"label, and the unrelated label untouched", step.want, got, want)
		}
		if step.want == "resolving" {
			if got := f.read(t, "owner"); !slices.Equal(got, []string{"kata-ship/1"}) {
				t.Errorf("owner after the claim = %q; want the bound session", got)
			}
		}
	}

	want := []string{
		"claim q14r --as kata-ship/1",
		"label rm q14r lifecycle:queued",
		"label add q14r lifecycle:resolving",
		"label rm q14r lifecycle:resolving",
		"label add q14r lifecycle:refining",
		"label rm q14r lifecycle:refining",
		"label add q14r lifecycle:shipping",
		"label rm q14r lifecycle:shipping",
		"unassign q14r",
		"label add q14r lifecycle:closed",
	}
	if got := f.mutations(t); !slices.Equal(got, want) {
		t.Errorf("the tool saw\n%q\nwant, in order,\n%q", got, want)
	}
	if got := f.read(t, "owner"); len(got) != 0 {
		t.Errorf("owner after close = %q; want unassigned", got)
	}
}

func TestStepsQ14r_ALostClaimRaceRefusesWithNothingApplied(t *testing.T) {
	f := newFakeKata(t, "lifecycle:queued")
	f.control(t, "claim-exit", "5")

	_, err := f.move(t, "claim")
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("a lost claim did not refuse: %v", err)
	}
	if ce.Code != codeAccessorFailed {
		t.Errorf("code = %q; want %q — the claim was the FIRST step, so "+
			"nothing ran and the refusal keeps the safe-to-retry sense", ce.Code, codeAccessorFailed)
	}
	if strings.Contains(ce.Detail, detailMayHaveApplied) {
		t.Errorf("Detail carries the applied sense for a write that applied nothing: %q", ce.Detail)
	}
	if got := f.mutations(t); !slices.Equal(got, []string{"claim q14r --as kata-ship/1"}) {
		t.Errorf("the tool saw %q; want the claim alone — fail-fast stops "+
			"before any label step", got)
	}
	if got := f.read(t, "labels"); !slices.Equal(got, []string{"area:cli", "lifecycle:queued"}) {
		t.Errorf("labels = %q; want them unchanged", got)
	}
}

func TestStepsQ14r_ALabelThatDidNotPersistRefusesAtReadBack(t *testing.T) {
	f := newFakeKata(t, "lifecycle:resolving")
	f.control(t, "nopersist", "")

	_, err := f.move(t, "advance")
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("a label add that exited 0 without persisting verified: %v", err)
	}
	if ce.Code != codeReadBackMismatch {
		t.Errorf("code = %q; want %q — every step exited 0, and only the "+
			"read-back can see the label never landed", ce.Code, codeReadBackMismatch)
	}
}

func TestStepsQ14r_ALaterStepFailureRefusesInTheAppliedSense(t *testing.T) {
	f := newFakeKata(t, "lifecycle:resolving")
	f.control(t, "add-exit", "1")

	_, err := f.move(t, "advance")
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("a failing second step did not refuse: %v", err)
	}
	if ce.Code != codeWriteFailedApplied || !strings.Contains(ce.Detail, detailMayHaveApplied) {
		t.Errorf("code = %q, Detail %q; want %q with the applied sense — "+
			"the label remove already ran", ce.Code, ce.Detail, codeWriteFailedApplied)
	}
	if !strings.Contains(ce.Detail, "`steps.lifecycle.set.refining[0]`") {
		t.Errorf("Detail %q does not name the failing step", ce.Detail)
	}
	if got := clierr.ExitCodeFor(err); got != 3 {
		t.Errorf("exit = %d; want 3", got)
	}
	// No undo: the remove stands (`0004:C14`).
	if got := f.read(t, "labels"); !slices.Equal(got, []string{"area:cli"}) {
		t.Errorf("labels = %q; want the remove standing and nothing undone", got)
	}
}

func TestStepsQ14r_ClearingAnAbsentKeyRunsNothingAndVerifies(t *testing.T) {
	f := newFakeKata(t)

	args := append([]string{"flow", "set-state", "--model", f.model,
		"--allow-commands", "--clear", "lifecycle", "--as=json"}, f.bind...)
	requireSuccess(t, args...)
	if got := f.mutations(t); len(got) != 0 {
		t.Errorf("the tool saw %q; clearing a key the artifact does not "+
			"hold must run no step (`0004:C11`)", got)
	}

	held := newFakeKata(t, "lifecycle:shipping")
	args = append([]string{"flow", "set-state", "--model", held.model,
		"--allow-commands", "--clear", "lifecycle", "--as=json"}, held.bind...)
	requireSuccess(t, args...)
	if got := held.mutations(t); !slices.Equal(got, []string{"label rm q14r lifecycle:shipping"}) {
		t.Errorf("the tool saw %q; want `clear.shipping` alone", got)
	}
}

func TestStepsQ14r_AnUndeclaredClearRefusesBeforeAnyStep(t *testing.T) {
	// `queued` declares no clear arm, so a removal of it has no argv.
	f := newFakeKata(t, "lifecycle:queued")

	args := append([]string{"flow", "set-state", "--model", f.model,
		"--allow-commands", "--clear", "lifecycle"}, f.bind...)
	_, _, err := runCmd(t, args...)
	if got := clierr.ExitCodeFor(err); err == nil || got != 2 {
		t.Fatalf("exit = %d (%v); want 2 — a removal the model declares no "+
			"argv for is about the request", got, err)
	}
	if got := f.mutations(t); len(got) != 0 {
		t.Errorf("the tool saw %q; want no step", got)
	}
}

func TestStepsQ14r_ReassertingTheHeldValueRunsNoClear(t *testing.T) {
	f := newFakeKata(t, "lifecycle:refining")

	args := append([]string{"flow", "set-state", "--model", f.model,
		"--allow-commands", "--write", "lifecycle=refining", "--as=json"}, f.bind...)
	requireSuccess(t, args...)
	want := []string{"label add q14r lifecycle:refining"}
	if got := f.mutations(t); !slices.Equal(got, want) {
		t.Errorf("the tool saw %q; want %q — a clear before re-adding the "+
			"same label would remove what the set then restores", got, want)
	}
}

func TestStepsQ14r_WithoutTheGateNothingRuns(t *testing.T) {
	f := newFakeKata(t, "lifecycle:resolving")

	args := append([]string{"flow", "set-state", "--model", f.model,
		"--write", "lifecycle=refining"}, f.bind...)
	_, _, err := runCmd(t, args...)
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("a steps write ran without --allow-commands: %v", err)
	}
	if !strings.Contains(ce.Detail, "allow_commands") {
		t.Errorf("Detail %q does not name the gate", ce.Detail)
	}
	// A missing opt-in is a refusal of the REQUEST: the exit-2 group, so a
	// caller's retry loop does not re-run what can never succeed.
	if ce.Code != codeRequestRefused {
		t.Errorf("code = %q; want %q", ce.Code, codeRequestRefused)
	}
	if got := clierr.ExitCodeFor(err); got != 2 {
		t.Errorf("exit = %d; want 2", got)
	}
	if strings.Contains(ce.Detail, detailMayHaveApplied) {
		t.Errorf("Detail carries the applied sense: %q", ce.Detail)
	}
	if got := f.read(t, "log"); len(got) != 0 {
		t.Errorf("the tool ran %q without the gate", got)
	}
}

// One Apply is one entry, however many steps it runs (`0025:C7`).
func TestStepsQ14r_AThreeStepWriteIsOneInvocation(t *testing.T) {
	f := newFakeKata(t, "lifecycle:queued")

	src, err := os.ReadFile(f.model)
	if err != nil {
		t.Fatal(err)
	}
	m, err := table.Load(src, f.model)
	if err != nil {
		t.Fatalf("the model refused: %v", err)
	}
	reg := flowbind.Registry(m, f.dir, true)
	def, ok := reg.Lookup("kata", accessor.CapWrite)
	if !ok {
		t.Fatal("no write definition `kata`")
	}
	w, ok := def.Binding.(*cmdbind.StepsWriter)
	if !ok {
		t.Fatalf("the registry bound %T; want *cmdbind.StepsWriter", def.Binding)
	}

	art := accessor.Artifact{Role: "kata", Path: f.dir,
		Context: map[string]string{"id": "q14r", "session": "kata-ship/1"}}
	got := accessor.NewExecutor(reg, accessor.Artifacts{"kata": art}).Write(
		t.Context(), "kata", resolve.Plan{
			Writes: []resolve.Tag{{Key: "lifecycle", Value: "resolving"}},
		})
	if got.Refusal != nil {
		t.Fatalf("the write refused %q: %s", got.Refusal.Class, got.Refusal.Detail)
	}
	if n := w.Invocations(); n != 1 {
		t.Errorf("Invocations() = %d after a three-step write; want 1", n)
	}
	if n := len(f.mutations(t)); n != 3 {
		t.Errorf("the tool saw %d mutations; want the three steps", n)
	}
}

func TestStepsQ14r_AnUnreadablePriorRefusesBeforeAnyStep(t *testing.T) {
	// Two lifecycle labels: the prefix projection cannot say which one the
	// key holds, so it is unreadable — and running `set` without knowing
	// which `clear` it owed could leave three.
	f := newFakeKata(t, "lifecycle:resolving", "lifecycle:refining")

	args := append([]string{"flow", "set-state", "--model", f.model,
		"--allow-commands", "--write", "lifecycle=shipping"}, f.bind...)
	_, _, err := runCmd(t, args...)
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("a write over an unreadable prior ran: %v", err)
	}
	if ce.Code != codeAccessorFailed || strings.Contains(ce.Detail, detailMayHaveApplied) {
		t.Errorf("code = %q, Detail %q; want %q without the applied sense",
			ce.Code, ce.Detail, codeAccessorFailed)
	}
	if got := f.mutations(t); len(got) != 0 {
		t.Errorf("the tool saw %q; want no step", got)
	}
}
