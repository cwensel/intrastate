package cmdbind_test

// Kata zwdn — a refusal whose child wrote NO stderr must still name its
// cause (`0004:FM`, `0025:C4` `detail:`).
//
// The defect: `executor.go::refusalOf` copied only `ExecError.Detail` and
// never read `ExecError.Err`, so every no-tail invocation failure — an
// unparseable stdout, a command that does not exist, an exit the entry's
// map does not list — arrived at the caller as a bare
// `flow-accessor-failed: the accessor <id> could not be executed` with no
// `detail` at all. Three different causes, one indistinguishable envelope.
//
// Nothing here mocks the binding. Every case spawns a REAL child through
// the same `Reader`/`Gate` the CLI binds and reads the `Refusal` the real
// executor mints, because "the cause reaches the refusal site" is a
// property of that whole trip or it is nothing — asserting it on a
// hand-built `ExecError` would pass against the very code this file was
// written to kill.
//
// The mutant every case below kills is the original line,
// `r.Detail = ee.Detail`: restoring it empties `Detail` on all seven and
// each `Detail == ""` arm goes red.

import (
	"strconv"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/cmdbind"
	"github.com/cwensel/intrastate/internal/table"
)

// --- harness --------------------------------------------------------------

// readRefusal drives one read through the REAL executor over a
// command-backed reader and returns the refusal it minted.
func readRefusal(t *testing.T, acc table.Accessor, art accessor.Artifact) *accessor.Refusal {
	t.Helper()

	reg := accessor.Registry{
		Flow: "cmdflow",
		Definitions: []accessor.Definition{{
			Identity: accessor.Identity{
				Flow: "cmdflow", Name: fxName, Capability: accessor.CapRead,
			},
			Accessor: acc,
			Binding:  cmdbind.Reader{Accessor: acc, Name: fxName, Config: allowed(t)},
		}},
		OwnedTags: acc.Keys,
	}
	e := accessor.NewExecutor(reg, accessor.Artifacts{fxRole: art})

	got := e.Read(ctxOf(t), fxName)
	if got.Refusal == nil {
		t.Fatal("the read did not refuse; every case in this file is an " +
			"invocation failure and a resolved read proves nothing about Detail")
	}
	return got.Refusal
}

// gateRefusal is the same for a gate.
func gateRefusal(t *testing.T, acc table.Accessor, art accessor.Artifact) *accessor.Refusal {
	t.Helper()

	reg := accessor.Registry{
		Flow: "cmdflow",
		Definitions: []accessor.Definition{{
			Identity: accessor.Identity{
				Flow: "cmdflow", Name: fxName, Capability: accessor.CapGate,
			},
			Accessor: acc,
			Binding:  cmdbind.Gate{Accessor: acc, Name: fxName, Config: allowed(t)},
		}},
		OwnedTags: acc.Keys,
	}
	e := accessor.NewExecutor(reg, accessor.Artifacts{fxRole: art})

	got := e.Gate(ctxOf(t), fxName)
	if got.Refusal == nil {
		t.Fatal("the gate did not refuse")
	}
	return got.Refusal
}

// wantsCause asserts the refusal carries a NON-EMPTY Detail naming every
// fragment. The non-empty arm is stated separately from the naming arms
// because it is the one the mutant trips, and its message is what tells a
// reader which defect came back.
func wantsCause(t *testing.T, r *accessor.Refusal, fragments ...string) {
	t.Helper()

	if r.Class != accessor.ClassExecutionFailure {
		t.Fatalf("Class = %q; want %q — these are invocation failures and the "+
			"exit-3 population is fixed, so this record changes no class",
			r.Class, accessor.ClassExecutionFailure)
	}
	if r.Detail == "" {
		t.Fatal("Detail is EMPTY; the child wrote no stderr, so the wrapped " +
			"error's own text is the only thing naming the cause and dropping " +
			"it makes this refusal read identically to every other one")
	}
	for _, f := range fragments {
		if !strings.Contains(r.Detail, f) {
			t.Errorf("Detail = %q; want it to name %q", r.Detail, f)
		}
	}
	// `0004:C7` reserves Reason for a gate DENY. The cause must not be
	// routed through it on the way past.
	if r.Reason != "" {
		t.Errorf("Reason = %q; `0004:C7` reserves it for a gate deny and the "+
			"cause travels on Detail", r.Reason)
	}
}

// --- the reproduction table ----------------------------------------------

// A stdout that is valid JSON but not `0025:C3`'s flat object of strings —
// the case that cost the RDR 0028 consumer a debugging cycle, because every
// `rdr` projector emits `{"facts":[…],"record":…}` and the refusal did not
// say so.
func TestNoTail_NestedObjectStdoutNamesTheWireShape(t *testing.T) {
	bin := helperBin(t)
	acc := entry(fxRole, []string{bin, "emit",
		`{"facts":[{"name":"status","value":"Draft"}],"record":"0099"}`, "0"})

	r := readRefusal(t, acc, artifactAt(t, "state.json"))

	// The SHAPE is the remedy — "project it" — and `encoding/json`'s own
	// text names a Go target type the author never wrote. The clause id is
	// NOT asserted: it is provenance the author cannot act on, so it stays
	// in the source comment and out of the CLI's `detail`.
	wantsCause(t, r, "flat JSON object of", "string values")
}

// A stdout that is not JSON at all. It must read differently from the
// nested-object case above, or naming the shape bought nothing.
func TestNoTail_NonJSONStdoutNamesTheWireShape(t *testing.T) {
	bin := helperBin(t)

	bad := readRefusal(t, entry(fxRole, []string{bin, "emit", "hello\n", "0"}),
		artifactAt(t, "state.json"))
	wantsCause(t, bad, "flat JSON object of", "string values")

	nested := readRefusal(t, entry(fxRole, []string{bin, "emit", `{"a":{"b":"c"}}`, "0"}),
		artifactAt(t, "state.json"))

	// Both name the shape; the decoder's complaint is what separates them.
	if bad.Detail == nested.Detail {
		t.Errorf("a non-JSON stdout and a wrongly-shaped one both read %q; the "+
			"two must stay distinguishable", bad.Detail)
	}
}

// An argv0 that does not resolve. The refusal must name the tool the author
// declared — it is the one thing they can act on, and `exec.ErrNotFound`
// carries it.
func TestNoTail_CommandNotFoundNamesTheDeclaredArgv0(t *testing.T) {
	const missing = "no-such-tool-zwdn"
	acc := entry(fxRole, []string{missing, "{artifact}"})

	r := readRefusal(t, acc, artifactAt(t, "state.json"))

	wantsCause(t, r, missing, "not found")
}

// Raw mode carries exactly one declared key. A read that requested more is
// a declaration defect, and the count is what says so.
func TestNoTail_RawModeWrongKeyCountNamesTheCount(t *testing.T) {
	bin := helperBin(t)
	acc := entry(fxRole, []string{bin, "emit", "Draft\n", "0"}, fxKey, fxKeyB)
	acc.Output = strptr("raw")

	r := readRefusal(t, acc, artifactAt(t, "state.json"))

	wantsCause(t, r, "raw mode carries one declared key", strconv.Itoa(2))
}

// An empty stdout and a signal before any exit: the entry's `exit_absent`
// map cannot apply, because there is no exit code for it to claim.
func TestNoTail_SignaledBeforeExitSaysTheExitMapCannotApply(t *testing.T) {
	bin := helperBin(t)
	// The helper kills its own process before exiting, so `Wait` reports a
	// signal and no exit code at all.
	acc := entry(fxRole, []string{bin, "self-signal"})
	acc.ExitAbsent = []int{0}

	r := readRefusal(t, acc, artifactAt(t, "state.json"))

	wantsCause(t, r, "killed by a signal", "exit_absent")
}

// An empty stdout and an exit the entry's `exit_absent` does not list. A
// silent tool is a broken tool, and the refusal must say which code it saw.
func TestNoTail_UnlistedExitNamesTheCodeAndTheMap(t *testing.T) {
	bin := helperBin(t)
	acc := entry(fxRole, []string{bin, "silent", "7"})
	acc.ExitAbsent = []int{1}

	r := readRefusal(t, acc, artifactAt(t, "state.json"))

	wantsCause(t, r, "exited 7", "exit_absent")
}

// A gate whose stdout is not the verdict envelope. Execution failure is
// never laundered into a verdict (`0025:C3`), and the refusal it takes
// instead must still say what was wrong with the stdout.
func TestNoTail_MalformedGateVerdictEnvelopeNamesTheShape(t *testing.T) {
	bin := helperBin(t)
	acc := entry(fxRole, []string{bin, "emit", "not-an-envelope\n", "0"})
	acc.ExitVerdicts = map[string]string{"0": "allow"}

	r := gateRefusal(t, acc, artifactAt(t, "state.json"))

	wantsCause(t, r, "verdict", "`reason`")
}

// --- the discriminating property -----------------------------------------

// The whole point of the kata: the three causes the reproduction table
// found indistinguishable must arrive as three DIFFERENT details. A build
// that stuffs one constant string into Detail everywhere would satisfy
// every "non-empty" arm above and still leave the author bisecting by hand.
func TestNoTail_TheThreeCausesAreDistinguishableFromEachOther(t *testing.T) {
	bin := helperBin(t)

	details := map[string]string{
		"nested-object stdout": readRefusal(t,
			entry(fxRole, []string{bin, "emit", `{"facts":[]}`, "0"}),
			artifactAt(t, "a.json")).Detail,
		"non-JSON stdout": readRefusal(t,
			entry(fxRole, []string{bin, "emit", "hello\n", "0"}),
			artifactAt(t, "b.json")).Detail,
		"command not found": readRefusal(t,
			entry(fxRole, []string{"no-such-tool-zwdn"}),
			artifactAt(t, "c.json")).Detail,
		"unlisted exit": readRefusal(t,
			entry(fxRole, []string{bin, "silent", "7"}),
			artifactAt(t, "d.json")).Detail,
	}

	seen := map[string]string{}
	for name, d := range details {
		if d == "" {
			t.Errorf("%s produced an empty Detail", name)
			continue
		}
		if prior, dup := seen[d]; dup {
			t.Errorf("%s and %s both read %q; the reproduction table's whole "+
				"finding is that these must not be the same envelope",
				prior, name, d)
			continue
		}
		seen[d] = name
	}
}

// A child that DID write stderr keeps the tail as its Detail, unprefixed.
// The fallback is exclusive on purpose: a present Detail is already the
// binding's composed diagnosis (`invocation.detail` leads with the
// held-pipe reason), and the `Err` beside it is either the bare exit status
// the tail already explains or RDR 0028's routing sentinel, whose text is a
// retry instruction and not a cause. This is the guard the premortem named.
func TestNoTail_AGenuineStderrTailIsNotPrefixedByTheErrorText(t *testing.T) {
	bin := helperBin(t)
	const tail = "cmdhelper: the tool refused\n"
	acc := entry(fxRole, []string{bin, "emit-stderr", tail, "9"})

	r := readRefusal(t, acc, artifactAt(t, "state.json"))

	if r.Detail != tail {
		t.Errorf("Detail = %q; want exactly the stderr tail %q — where the "+
			"child wrote one it is the diagnosis, and the wrapped error's "+
			"text must not be prepended to it (`0025:C4`)", r.Detail, tail)
	}
}
