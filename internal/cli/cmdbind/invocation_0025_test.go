package cmdbind_test

// RDR 0025 — the invocation envelope, the authority bound, and execution
// safety (C2, C3, C4; scenarios S2, S3, S4, S4b, S6b).
//
// The coverage floor this file holds, and the two mutants it kills:
//
//  1. A binding that PARSES the exit code before the stdout envelope turns
//     a well-formed deny with a non-zero exit into an execution failure —
//     a verdict the model DID decide, dropped. C4's ordering is normative
//     for exactly that reason, and TestReq48_ asserts the pair (deny with
//     a non-zero exit; malformed envelope with exit zero) so neither arm
//     can pass alone.
//  2. A binding that maps an exit code on a NON-empty stdout laundres a
//     truncated 1 MiB read into an established absence. C3 restricts the
//     maps to empty stdout, which is what makes `exit_absent` sound at
//     whole-invocation granularity.
//
// Every assertion is on what the CHILD observed or on what the binding
// RETURNED through the seam. Nothing reads a private field, no call count
// is asserted that the spec does not name, and the deadline arms use a
// real process group.

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/cmdbind"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

func ctxOf(t *testing.T) context.Context {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	return ctx
}

// readOnce drives one read through the binding and returns the seam's
// three-valued answer.
func readOnce(t *testing.T, acc table.Accessor, art accessor.Artifact) (
	[]accessor.KeyValue, []string, error,
) {
	t.Helper()

	r := cmdbind.Reader{Accessor: acc, Name: fxName, Config: allowed(t)}
	return r.Read(ctxOf(t), art, acc.Keys)
}

// rawValue drives a raw-mode read and returns the single established value.
func rawValue(t *testing.T, acc table.Accessor, art accessor.Artifact) accessor.KeyValue {
	t.Helper()

	values, unreadable, err := readOnce(t, acc, art)
	if err != nil {
		t.Fatalf("the read failed: %v", err)
	}
	if len(unreadable) != 0 {
		t.Fatalf("unreadable = %#v; want none", unreadable)
	}
	if len(values) != 1 {
		t.Fatalf("values = %#v; want exactly one", values)
	}
	return values[0]
}

// --- C2: the authority bound ---------------------------------------------

// REQ-16: "{artifact}   # the closed placeholder vocabulary, v1 complete:
// replaced whole-element by the caller-bound artifact path for the entry's
// declared role"
// REQ-17: "The executed argv is exactly the declared vector after
// whole-element placeholder substitution."
// REQ-19: "Execution invokes no shell and performs no other rewriting of
// the vector."
// REQ-123 (S2): "the absolute path is substituted whole-element and the
// argv the child observes equals the declared vector otherwise
// byte-for-byte"
// HAPPY PATH
//
// The Pre-Lock `fidelity` invariant: declared argv → executed argv is
// BYTE-FOR-BYTE equality except the one substitution. The oracle is the
// argv the CHILD observed, so a binding that re-quoted, word-split, or
// glob-expanded any element fails here.
func TestReq17_TheExecutedArgvIsTheDeclaredVectorAfterOneSubstitution(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	// Deliberately hostile elements: a glob, a space, a shell
	// metacharacter, and a quoted word. Under `os/exec` they cross
	// verbatim; under any shell they would not.
	declared := []string{
		bin, "echo-argv", "*", "two words", "a;b", "$HOME", "{artifact}",
	}
	acc := entry(fxRole, declared)
	acc.Output = strptr("raw")

	got := rawValue(t, acc, art)
	argv := observedArgv(t, got.Value)

	want := append([]string{}, declared...)
	want[len(want)-1] = art.Path
	if len(argv) != len(want) {
		t.Fatalf("the child observed %d args %#v; want %d %#v",
			len(argv), argv, len(want), want)
	}
	for i := range want {
		if argv[i] != want[i] {
			t.Errorf("argv[%d] = %q; want %q — the executed argv is the declared "+
				"vector byte-for-byte except the one whole-element substitution",
				i, argv[i], want[i])
		}
	}
}

// REQ-20: "the **command binding** refuses the invocation
// (`execution_failure`, before spawn) when the caller-bound artifact path
// is relative or begins with `-`"
// REQ-21: "The check sites in the binding, not in `executor.go`" and it
// "**refuses rather than absolutizes**"
// REQ-22: "the refusal `Detail` says so (\"artifact path must be absolute
// for a command entry\")"
// REQ-135: the same, as a Failure Mode.
// ADVERSARIAL
//
// PHASE-0 READING (ASSUMPTION REQ-20): "begins with `-`" is tested on the
// SUBSTITUTED artifact path only, not on other declared argv elements — a
// declared `--flag` is legitimate and must still run (asserted below).
func TestReq20_ARelativeOrDashLedArtifactPathRefusesBeforeSpawn(t *testing.T) {
	bin := helperBin(t)
	trace := filepath.Join(t.TempDir(), "spawned")

	cases := []struct {
		name string
		path string
	}{
		{"relative path", filepath.Join("relative", "state.cfg")},
		{"bare relative name", "state.cfg"},
		{"dash-prefixed path", "-oProxyCommand=evil"},
		{"dash-prefixed absolute-looking", "--file=/tmp/x"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			acc := entry(fxRole, []string{bin, "trace", trace, "{artifact}"})
			acc.Output = strptr("raw")

			_, _, err := readOnce(t, acc, accessor.Artifact{Role: fxRole, Path: tc.path})

			if err == nil {
				t.Fatal("the read did not refuse; a path that is not absolute can " +
					"be parsed as a flag by the invoked tool")
			}
			var ee *accessor.ExecError
			if !errors.As(err, &ee) {
				t.Fatalf("the refusal is not an *accessor.ExecError: %v", err)
			}
			if !strings.Contains(ee.Detail, "absolute") {
				t.Errorf("Detail = %q; want it to say the artifact path must be "+
					"absolute for a command entry", ee.Detail)
			}
			// BEFORE spawn: absence of a spawn, asserted by the trace file
			// the child would have written had it run.
			if _, serr := os.Stat(trace); serr == nil {
				t.Error("a child ran; the refusal is before spawn")
			}
		})
	}

	t.Run("a declared dash-led ELEMENT is unaffected", func(t *testing.T) {
		// The guard is on the substituted value, not on the vector.
		art := artifactAt(t, "state.cfg")
		acc := entry(fxRole, []string{bin, "echo-argv", "--file", "{artifact}"})
		acc.Output = strptr("raw")

		got := rawValue(t, acc, art)
		argv := observedArgv(t, got.Value)
		if len(argv) < 4 || argv[2] != "--file" {
			t.Errorf("argv = %#v; a declared --flag is legitimate and must cross "+
				"unchanged", argv)
		}
	})
}

// REQ-23: "argv0 resolution is fixed at load, never cwd-relative: a bare
// name resolves through the parent's `PATH` at spawn (A2); a name
// containing a path separator resolves against the **model file's
// directory**"
// REQ-25: "the command-binding constructor receives `filepath.Dir` of the
// **absolutized** model path from the caller that opened the file"
// REQ-27: "a relative `--model` must not make argv0 resolution
// cwd-dependent."
// DOMAIN EDGE
func TestReq23_ASeparatorBearingArgv0ResolvesAgainstTheModelDirNotTheCwd(t *testing.T) {
	bin := helperBin(t)
	baseDir := filepath.Dir(bin)
	art := artifactAt(t, "state.cfg")

	// A separator-bearing argv0 naming the helper RELATIVE to the base dir.
	rel := "." + string(filepath.Separator) + filepath.Base(bin)

	acc := entry(fxRole, []string{rel, "echo-argv"})
	acc.Output = strptr("raw")

	r := cmdbind.Reader{
		Accessor: acc,
		Name:     fxName,
		Config:   cmdbind.Config{BaseDir: baseDir, AllowCommands: true},
	}
	values, _, err := r.Read(ctxOf(t), art, acc.Keys)
	if err != nil {
		t.Fatalf("a separator-bearing argv0 must resolve against the model dir; "+
			"refused: %v", err)
	}
	if len(values) != 1 {
		t.Fatalf("values = %#v; want one", values)
	}
	argv := observedArgv(t, values[0].Value)
	if len(argv) == 0 || filepath.Base(argv[0]) != filepath.Base(bin) {
		t.Errorf("the child's argv0 = %#v; want the helper resolved against the "+
			"model file's directory", argv)
	}

	t.Run("the process cwd is never the fallback", func(t *testing.T) {
		// Same relative argv0, but the base dir is a directory that does
		// NOT hold the helper. Go reversed cwd-relative PATH resolution
		// (exec.ErrDot, 1.19) and this binding must not restore it, even
		// when the cwd happens to hold a matching name.
		other := t.TempDir()
		r := cmdbind.Reader{
			Accessor: acc,
			Name:     fxName,
			Config:   cmdbind.Config{BaseDir: other, AllowCommands: true},
		}
		if _, _, rerr := r.Read(ctxOf(t), art, acc.Keys); rerr == nil {
			t.Error("a separator-bearing argv0 resolved outside the declared " +
				"base dir; implicit current-directory lookup is the reversal " +
				"this clause exists to honour")
		}
	})
}

// REQ-28: "A separator-bearing argv0 in [a model with no source file] is
// refused `execution_failure` at **binding construction**, before any spawn
// — not a load-time defect ... What it must never do is fall back to the
// process cwd"
// BOUNDARY
func TestReq28_ASeparatorBearingArgv0WithNoBaseDirRefusesBeforeSpawn(t *testing.T) {
	art := artifactAt(t, "state.cfg")
	trace := filepath.Join(t.TempDir(), "spawned")

	acc := entry(fxRole, []string{"./tools/writer", "trace", trace})
	acc.Output = strptr("raw")

	r := cmdbind.Reader{
		Accessor: acc,
		Name:     fxName,
		Config:   cmdbind.Config{BaseDir: "", AllowCommands: true},
	}
	_, _, err := r.Read(ctxOf(t), art, acc.Keys)

	if err == nil {
		t.Fatal("a separator-bearing argv0 with no base dir did not refuse")
	}
	if _, serr := os.Stat(trace); serr == nil {
		t.Error("a child ran; the refusal is at construction, before any spawn")
	}

	t.Run("the SAME argv0 resolves once a base dir is supplied", func(t *testing.T) {
		// Positive control. Without it, "a separator-bearing argv0 with no
		// base dir refuses" is satisfied by a binding that refuses every
		// separator-bearing argv0, or by one that never runs anything at
		// all — and the clause's content, that the BASE DIR is what makes
		// the difference, goes unasserted.
		bin := helperBin(t)
		rel := "." + string(filepath.Separator) + filepath.Base(bin)
		acc := entry(fxRole, []string{rel, "emit", "draft\n", "0"})
		acc.Output = strptr("raw")

		r := cmdbind.Reader{
			Accessor: acc,
			Name:     fxName,
			Config:   cmdbind.Config{BaseDir: filepath.Dir(bin), AllowCommands: true},
		}
		values, _, rerr := r.Read(ctxOf(t), art, acc.Keys)
		if rerr != nil {
			t.Fatalf("the same argv0 refused WITH a base dir: %v", rerr)
		}
		if len(values) != 1 || values[0].Value != "draft" {
			t.Errorf("values = %#v; want the child's value", values)
		}
	})

	t.Run("a BARE argv0 is unaffected by a missing base dir", func(t *testing.T) {
		// A bare name resolves through the parent's PATH (A2), which the C4
		// allowlist passes on unchanged — the base dir is irrelevant to it.
		// The oracle is the VALUE, not a message: a bare `echo` under an
		// empty base dir must produce its output, so a build that refused
		// every entry for a missing base dir fails here.
		acc := entry(fxRole, []string{"echo", "draft"})
		acc.Output = strptr("raw")
		r := cmdbind.Reader{
			Accessor: acc,
			Name:     fxName,
			Config:   cmdbind.Config{BaseDir: "", AllowCommands: true},
		}
		values, _, rerr := r.Read(ctxOf(t), art, acc.Keys)
		if rerr != nil {
			t.Fatalf("a BARE argv0 was refused with no base dir; it resolves "+
				"through the parent's PATH at spawn (A2): %v", rerr)
		}
		if len(values) != 1 || values[0].Value != "draft" {
			t.Errorf("values = %#v; want the child's single value %q",
				values, "draft")
		}
	})
}

// --- C3: the invocation envelope -----------------------------------------

// REQ-30: "stdin (read/gate): {}                       # nothing
// per-invocation beyond {artifact}; the object is always sent"
// REQ-37: "Per-invocation data beyond the artifact path crosses only on
// stdin; read and gate results return only on stdout"
// REQ-111: "Rejected: carrying the requested key set on stdin (a list is
// not a string value, and the keys are already declared in the model)."
// BOUNDARY — REQ-111 is a NEGATIVE REQ.
func TestReq30_AReadSendsAnEmptyObjectOnStdinAndNeverTheKeySet(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	acc := entry(fxRole, []string{bin, "echo-stdin"}, fxKey, fxKeyB)
	acc.Output = strptr("raw")
	acc.Keys = []string{fxKey}

	got := rawValue(t, acc, art)

	if got.Value != "{}" {
		t.Errorf("the child observed stdin %q; want the empty object %q — the "+
			"object is ALWAYS sent, and the requested key set is not on it "+
			"(REQ-111)", got.Value, "{}")
	}
}

// REQ-29: "stdin (write): {\"<key>\": \"<value>\", ...}   # the planned
// tags; a planned `<clear>` crosses as the literal reserved value"
// REQ-40: "`<clear>` is carried unchanged as the planned value — the tool
// or wrapper performs the removal"
// HAPPY PATH
func TestReq29_AWriteSendsThePlannedTagsAsAFlatObjectOfStrings(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")
	sink := filepath.Join(t.TempDir(), "stdin.json")

	acc := entry(fxRole, []string{bin, "stdin-to-file", sink}, fxKey, fxKeyB)
	acc.ReadBack = true

	w := &cmdbind.Writer{Accessor: acc, Name: fxName, Config: allowed(t)}
	err := w.Apply(ctxOf(t), art, []resolve.Tag{
		{Key: fxKey, Value: "final"},
		{Key: fxKeyB, Value: table.ClearSentinel},
	})
	if err != nil {
		t.Fatalf("the write failed: %v", err)
	}

	b, rerr := os.ReadFile(sink)
	if rerr != nil {
		t.Fatalf("the child recorded no stdin: %v", rerr)
	}
	got := strings.TrimSpace(string(b))

	// Both keys, both as STRING values, and the clear as its literal.
	if !strings.Contains(got, `"`+fxKey+`":"final"`) &&
		!strings.Contains(got, `"`+fxKey+`": "final"`) {
		t.Errorf("stdin = %s; want the planned value carried as a string", got)
	}
	if !strings.Contains(got, table.ClearSentinel) {
		t.Errorf("stdin = %s; want the planned `%s` carried UNCHANGED as the "+
			"literal reserved value — the tool performs the removal", got,
			table.ClearSentinel)
	}
}

// REQ-31: "stdout (read, output = \"json\", default): flat JSON object of
// strings; a declared key the object omits is UNREADABLE, never
// established-absent"
// REQ-105: "a parsed key becomes a `KeyValue` in `values`, an omitted
// declared key a name in `unreadable`"
// ADVERSARIAL
//
// The discriminating PAIR: the same envelope shape, one carrying both
// declared keys and one omitting the second. Omission must land in
// `unreadable`, NOT as an Absent value — the two are different answers and
// only one is safe to guess.
func TestReq31_AJSONEnvelopeOmissionIsUnreadableNeverEstablishedAbsent(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	t.Run("both keys carried", func(t *testing.T) {
		acc := entry(fxRole,
			[]string{bin, "emit", `{"` + fxKey + `":"final","` + fxKeyB + `":"prelock"}`, "0"},
			fxKey, fxKeyB)

		values, unreadable, err := readOnce(t, acc, art)
		if err != nil {
			t.Fatalf("the read failed: %v", err)
		}
		if len(unreadable) != 0 {
			t.Errorf("unreadable = %#v; want none", unreadable)
		}
		if len(values) != 2 {
			t.Fatalf("values = %#v; want both declared keys", values)
		}
		for _, v := range values {
			if v.Absent {
				t.Errorf("%q came back Absent; the envelope carried it", v.Key)
			}
		}
	})

	t.Run("an omitted declared key is UNREADABLE", func(t *testing.T) {
		acc := entry(fxRole,
			[]string{bin, "emit", `{"` + fxKey + `":"final"}`, "0"},
			fxKey, fxKeyB)

		values, unreadable, err := readOnce(t, acc, art)
		if err != nil {
			t.Fatalf("the read failed: %v", err)
		}
		if len(unreadable) != 1 || unreadable[0] != fxKeyB {
			t.Errorf("unreadable = %#v; want [%s] — a tool that fails to report "+
				"a key has not established that the key has no value",
				unreadable, fxKeyB)
		}
		for _, v := range values {
			if v.Key == fxKeyB {
				t.Errorf("%q came back in values as %+v; omission is UNREADABLE, "+
					"never established-absent", fxKeyB, v)
			}
		}
	})
}

// REQ-32: "stdout (read, output = \"raw\"):            the single declared
// key's value = stdout minus one trailing \"\\n\"; valid only when `keys`
// has exactly one entry"
// REQ-140: "The `raw` mode is defined bytewise: the value is the child's
// stdout minus **exactly one** trailing `\n` — no trimming, no case
// folding, no whitespace normalization"
// BOUNDARY — normative fixture FX-raw-read (spike R1: stdout "draft\n" →
// "draft").
//
// PHASE-0 READING (ASSUMPTION REQ-32): "minus one trailing `\n`" removes at
// most one, and only if present.
func TestReq32_RawModeStripsExactlyOneTrailingNewlineAndNothingElse(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	cases := []struct {
		name   string
		stdout string
		want   string
	}{
		// FX-raw-read: the normative fixture line.
		{"fixture FX-raw-read", "draft\n", "draft"},
		{"no trailing newline is used verbatim", "draft", "draft"},
		{"a doubled newline keeps one", "draft\n\n", "draft\n"},
		{"leading whitespace survives", "  draft\n", "  draft"},
		{"internal whitespace survives", "a b\tc\n", "a b\tc"},
		{"case is not folded", "DrAfT\n", "DrAfT"},
		{"a trailing CR is not stripped", "draft\r\n", "draft\r"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			acc := entry(fxRole, []string{bin, "emit", tc.stdout, "0"})
			acc.Output = strptr("raw")

			got := rawValue(t, acc, art)
			if got.Value != tc.want {
				t.Errorf("value = %q; want %q — no trimming, no case folding, no "+
					"whitespace normalization", got.Value, tc.want)
			}
			if got.Key != fxKey {
				t.Errorf("key = %q; want the single declared key %q", got.Key, fxKey)
			}
			if got.Absent {
				t.Error("the value came back Absent; raw mode ESTABLISHES the key " +
					"from stdout")
			}
		})
	}
}

// REQ-33: "stdout (gate): {\"verdict\": \"allow\" | \"deny\" |
// \"indeterminate\", \"reason\": \"<text>\"}"
// REQ-104 (gate arm): "`Gate` returns verdict + reason + error, so C3's
// `reason` field crosses back through the Go return and is not dropped."
// HAPPY PATH
func TestReq33_TheGateEnvelopeCarriesVerdictAndReasonBackThroughTheSeam(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	for _, v := range accessor.Verdicts() {
		t.Run(string(v), func(t *testing.T) {
			acc := entry(fxRole, []string{bin, "emit",
				`{"verdict":"` + string(v) + `","reason":"the tree is dirty"}`, "0"})

			g := cmdbind.Gate{Accessor: acc, Name: fxName, Config: allowed(t)}
			verdict, reason, err := g.Gate(ctxOf(t), art)
			if err != nil {
				t.Fatalf("the gate failed: %v", err)
			}
			if verdict != v {
				t.Errorf("verdict = %q; want %q", verdict, v)
			}
			if reason != "the tree is dirty" {
				t.Errorf("reason = %q; want the envelope's reason carried back "+
					"through the Go return", reason)
			}
		})
	}
}

// --- C3/C4: the ordering and the exit maps (S3's nine cases) -------------

// REQ-43: "A non-empty stdout is always parsed first (C4 ordering); the
// exit maps apply only to an empty stdout — including a stdout truncated at
// the 1 MiB cap, which is *non-empty* and therefore an `execution_failure`
// from the bound (C4), never an exit-map consultation."
// REQ-48: "order:    parse the stdout envelope, THEN classify the exit
// code"
// REQ-67: "a well-formed deny envelope with a non-zero exit is a deny, not
// a failure."
// REQ-133: "a gate command with a malformed envelope refuses
// `execution_failure` — never a deny the model didn't decide."
// ADVERSARIAL
//
// S3 cases 1 and 2, the CONTRACT-DERIVED pair: they have no fixture oracle
// and need none — their expected values follow from C3/C4's ordering rule
// directly. Asserted together so neither passes alone.
func TestReq48_TheStdoutEnvelopeIsParsedBeforeTheExitCodeIsClassified(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	t.Run("case 1: a well-formed deny envelope with a non-zero exit is a DENY", func(t *testing.T) {
		acc := entry(fxRole, []string{bin, "emit",
			`{"verdict":"deny","reason":"unclean"}`, "1"})

		g := cmdbind.Gate{Accessor: acc, Name: fxName, Config: allowed(t)}
		verdict, reason, err := g.Gate(ctxOf(t), art)
		if err != nil {
			t.Fatalf("a well-formed deny with a non-zero exit refused: %v — "+
				"parse-first means the envelope wins over the exit", err)
		}
		if verdict != accessor.VerdictDeny {
			t.Errorf("verdict = %q; want %q", verdict, accessor.VerdictDeny)
		}
		if reason != "unclean" {
			t.Errorf("reason = %q; want %q", reason, "unclean")
		}
	})

	t.Run("case 2: a malformed envelope with exit ZERO is execution_failure", func(t *testing.T) {
		acc := entry(fxRole, []string{bin, "emit", "not json at all", "0"})

		g := cmdbind.Gate{Accessor: acc, Name: fxName, Config: allowed(t)}
		verdict, _, err := g.Gate(ctxOf(t), art)
		if err == nil {
			t.Fatalf("a malformed envelope answered verdict %q; a malformed "+
				"envelope is a failure, never a deny the model didn't decide",
				verdict)
		}

		// The discriminating sibling: the SAME exit code with a WELL-FORMED
		// envelope must answer. Without it, "malformed refuses" is
		// satisfied by a gate that refuses everything.
		ok := entry(fxRole, []string{bin, "emit", `{"verdict":"allow"}`, "0"})
		g2 := cmdbind.Gate{Accessor: ok, Name: fxName, Config: allowed(t)}
		v2, _, err2 := g2.Gate(ctxOf(t), art)
		if err2 != nil {
			t.Fatalf("a well-formed envelope with exit 0 refused: %v", err2)
		}
		if v2 != accessor.VerdictAllow {
			t.Errorf("verdict = %q; want %q", v2, accessor.VerdictAllow)
		}
	})
}

// REQ-35: "empty stdout, no exit-map match: execution_failure, on BOTH read
// modes and on gate — never UNREADABLE and never established-absent. ... a
// read whose exit is unlisted (exit 0 included) refuses"
// REQ-45: "Under `output = \"raw\"` ... an empty stdout with an unlisted
// exit is `execution_failure`, not an empty string."
// ADVERSARIAL
//
// S3 case 3 (unmapped non-zero exit) plus the exit-ZERO sibling the clause
// makes explicit. Empty-stdout arms follow normative fixture FX-exit-codes
// (E1c: exit 2, empty stdout, stderr set).
func TestReq35_AnEmptyStdoutWithNoExitMapMatchIsAlwaysExecutionFailure(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	t.Run("raw read, unlisted non-zero exit", func(t *testing.T) {
		// FX-exit-codes E1c: `test --bogus` → exit 2, empty stdout.
		acc := entry(fxRole, []string{bin, "emit-stderr",
			"test: --bogus: unexpected operator\n", "2"})
		acc.Output = strptr("raw")
		acc.ExitAbsent = []int{1}

		_, _, err := readOnce(t, acc, art)
		if err == nil {
			t.Fatal("an unlisted exit with empty stdout did not refuse; a silent " +
				"tool is a broken tool, not an answer")
		}
		// REQ-124 (S3): the tail is asserted non-empty on THIS case, where
		// no applied-sense text competes for the CLI's single Detail slot.
		var ee *accessor.ExecError
		if !errors.As(err, &ee) {
			t.Fatalf("the refusal is not an *accessor.ExecError: %v", err)
		}
		if ee.Detail == "" {
			t.Error("Detail is empty; C4 carries the bounded stderr tail on an " +
				"execution failure")
		}
		if !strings.Contains(ee.Detail, "unexpected operator") {
			t.Errorf("Detail = %q; want the child's stderr tail", ee.Detail)
		}
	})

	t.Run("raw read, exit ZERO with empty stdout still refuses", func(t *testing.T) {
		acc := entry(fxRole, []string{bin, "silent", "0"})
		acc.Output = strptr("raw")

		values, unreadable, err := readOnce(t, acc, art)
		if err == nil {
			t.Fatalf("exit 0 with empty stdout answered values=%#v unreadable=%#v; "+
				"the clause says exit 0 INCLUDED refuses — never an empty string",
				values, unreadable)
		}

		// The discriminating sibling: the same exit 0 with NON-empty stdout
		// establishes the value. Without it, "empty stdout refuses" is
		// satisfied by a read that refuses everything.
		ok := entry(fxRole, []string{bin, "emit", "draft\n", "0"})
		ok.Output = strptr("raw")
		if got := rawValue(t, ok, art); got.Value != "draft" {
			t.Errorf("value = %q; want %q", got.Value, "draft")
		}
	})

	t.Run("json read, exit ZERO with empty stdout refuses", func(t *testing.T) {
		acc := entry(fxRole, []string{bin, "silent", "0"})

		if _, _, err := readOnce(t, acc, art); err == nil {
			t.Fatal("exit 0 with empty stdout did not refuse on the json mode; " +
				"the rule is BOTH read modes")
		}

		// The discriminating sibling on the json mode.
		ok := entry(fxRole, []string{bin, "emit", `{"` + fxKey + `":"final"}`, "0"})
		values, unreadable, err := readOnce(t, ok, art)
		if err != nil {
			t.Fatalf("a well-formed json envelope refused: %v", err)
		}
		if len(values) != 1 || len(unreadable) != 0 {
			t.Errorf("values=%#v unreadable=%#v; want the one declared key",
				values, unreadable)
		}
	})

	t.Run("gate, unlisted non-zero exit with empty stdout", func(t *testing.T) {
		// FX-exit-codes E1f: `git -C norepo diff --quiet` → exit 128.
		acc := entry(fxRole, []string{bin, "silent", "128"})
		acc.ExitVerdicts = map[string]string{"0": "allow", "1": "deny"}

		g := cmdbind.Gate{Accessor: acc, Name: fxName, Config: allowed(t)}
		verdict, _, err := g.Gate(ctxOf(t), art)
		if err == nil {
			t.Fatalf("exit 128 under a map listing only 0 and 1 answered %q; an "+
				"unlisted exit is an execution failure, never laundered into a "+
				"verdict", verdict)
		}

		// The discriminating sibling: a LISTED exit under the same map is
		// that verdict, so the refusal above is about the map and not about
		// a gate that refuses everything.
		ok := entry(fxRole, []string{bin, "silent", "1"})
		ok.ExitVerdicts = map[string]string{"0": "allow", "1": "deny"}
		g2 := cmdbind.Gate{Accessor: ok, Name: fxName, Config: allowed(t)}
		v2, _, err2 := g2.Gate(ctxOf(t), art)
		if err2 != nil {
			t.Fatalf("a listed exit refused: %v", err2)
		}
		if v2 != accessor.VerdictDeny {
			t.Errorf("verdict = %q; want %q", v2, accessor.VerdictDeny)
		}
	})
}

// REQ-36: "exit_verdicts = { \"<code>\" = \"allow\" | \"deny\" |
// \"indeterminate\", ... }   # gate entries: a listed exit with empty
// stdout is that verdict"
// REQ-41: "the maps list **verdict** codes only — any unlisted exit, spawn
// failure, or malformed stdout is `execution_failure`"
// DOMAIN EDGE — normative fixture FX-exit-codes (E1d/E1e: `git diff
// --quiet` 0 = clean, 1 = dirty).
func TestReq36_AListedExitWithEmptyStdoutIsThatVerdict(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	cases := []struct {
		name string
		code string
		want accessor.Verdict
	}{
		{"FX-exit-codes E1d: clean tree exits 0", "0", accessor.VerdictAllow},
		{"FX-exit-codes E1e: dirty tree exits 1", "1", accessor.VerdictDeny},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			acc := entry(fxRole, []string{bin, "silent", tc.code})
			acc.ExitVerdicts = map[string]string{"0": "allow", "1": "deny"}

			g := cmdbind.Gate{Accessor: acc, Name: fxName, Config: allowed(t)}
			verdict, _, err := g.Gate(ctxOf(t), art)
			if err != nil {
				t.Fatalf("a listed exit with empty stdout refused: %v", err)
			}
			if verdict != tc.want {
				t.Errorf("verdict = %q; want %q", verdict, tc.want)
			}
		})
	}
}

// REQ-34: "exit_absent   = [<code>, ...]                # read entries: a
// listed exit with empty stdout establishes every declared key absent"
// REQ-44: "absence reaches the executor as `KeyValue{Absent: true}` in
// `values`, never as omission from both slices"
// REQ-106: the same, quoted on C7.
// DOMAIN EDGE — normative fixture FX-exit-codes (R6/E1g: `git config --get`
// of a missing key → exit 1, empty stdout).
//
// PHASE-0 READING (ASSUMPTION REQ-34/REQ-106): EVERY key in the entry's
// declared `keys` is returned Absent — the clause says "every declared key
// absent" and the signal is whole-invocation.
func TestReq34_AListedExitAbsentEstablishesEveryDeclaredKeyAbsent(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	// Two declared keys, so "every declared key" is discriminating: a
	// binding that flagged only the first would pass a one-key fixture.
	acc := entry(fxRole, []string{bin, "silent", "1"}, fxKey, fxKeyB)
	acc.ExitAbsent = []int{1}

	values, unreadable, err := readOnce(t, acc, art)
	if err != nil {
		t.Fatalf("a listed exit_absent code refused: %v — FX-exit-codes R6 is "+
			"the established-absent read", err)
	}
	if len(unreadable) != 0 {
		t.Errorf("unreadable = %#v; established absence is a VALUE, not an "+
			"unreadability", unreadable)
	}
	if len(values) != 2 {
		t.Fatalf("values = %#v; want one record per DECLARED key", values)
	}
	for _, v := range values {
		if !v.Absent {
			t.Errorf("%q came back present with %q; a listed exit with empty "+
				"stdout establishes EVERY declared key absent", v.Key, v.Value)
		}
	}
}

// REQ-42: "A spawn failure has no exit code at all (`exec.ErrNotFound`,
// fixture E1h), so no map entry can match it however the map is written —
// the maps are consulted only for a process that ran and exited."
// ADVERSARIAL — normative fixture FX-exit-codes E1h, the NINTH case.
//
// This is the arm that proves spawn failure cannot be laundered into a
// verdict by an over-broad map. The map below lists EVERY plausible code
// and still must not claim the spawn failure.
func TestReq42_ASpawnFailureIsNeverClaimedByAnOverBroadExitMap(t *testing.T) {
	art := artifactAt(t, "state.cfg")

	broad := map[string]string{}
	for i := range 256 {
		broad[itoa(i)] = "allow"
	}

	acc := entry(fxRole, []string{"definitely-not-a-binary-xyz", "{artifact}"})
	acc.ExitVerdicts = broad

	g := cmdbind.Gate{Accessor: acc, Name: fxName, Config: allowed(t)}
	verdict, _, err := g.Gate(ctxOf(t), art)

	if err == nil {
		t.Fatalf("a spawn failure answered verdict %q under a map listing every "+
			"code 0-255; a spawn failure has NO exit code, so no map entry can "+
			"match it", verdict)
	}
	// And it stays distinguishable from a non-zero exit (REQ-58's claim,
	// witnessed at the site that produces it).
	if !errors.Is(err, exec.ErrNotFound) {
		t.Errorf("errors.Is(err, exec.ErrNotFound) is false for a missing "+
			"executable; the wrap must survive: %v", err)
	}

	t.Run("the read side too", func(t *testing.T) {
		racc := entry(fxRole, []string{"definitely-not-a-binary-xyz"})
		racc.ExitAbsent = make([]int, 256)
		for i := range racc.ExitAbsent {
			racc.ExitAbsent[i] = i
		}
		values, _, rerr := readOnce(t, racc, art)
		if rerr == nil {
			t.Fatalf("a spawn failure established absence %#v under an "+
				"all-codes exit_absent map", values)
		}

		// The discriminating sibling: a real child that EXITS with a listed
		// code does establish absence, so the refusal above is about the
		// spawn failure having no exit code at all.
		bin := helperBin(t)
		ok := entry(fxRole, []string{bin, "silent", "1"})
		ok.ExitAbsent = []int{1}
		okValues, _, oerr := readOnce(t, ok, art)
		if oerr != nil {
			t.Fatalf("a listed exit_absent code refused: %v", oerr)
		}
		if len(okValues) != 1 || !okValues[0].Absent {
			t.Errorf("values = %#v; want one established-absent record", okValues)
		}
	})
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var b []byte
	for i > 0 {
		b = append([]byte{byte('0' + i%10)}, b...)
		i /= 10
	}
	return string(b)
}

// REQ-50 (bounds arm): "stdout capped at 1 MiB (overflow =
// execution_failure)"
// REQ-43 (cap arm): "a stdout truncated at the 1 MiB cap ... is *non-empty*
// and therefore an `execution_failure` from the bound (C4), never an
// exit-map consultation."
// BOUNDARY
//
// PHASE-0 READING (ASSUMPTION REQ-50): the cap is enforced by BOUNDING the
// read, so an unbounded child cannot exhaust memory; overflow is DETECTED,
// not silently truncated.
func TestReq50_AStdoutOverTheOneMiBCapIsAnExecutionFailureNotATruncation(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	t.Run("one byte over the cap refuses", func(t *testing.T) {
		acc := entry(fxRole, []string{bin, "flood", itoa(cmdbind.StdoutCap + 1)})
		acc.Output = strptr("raw")
		// An exit map that WOULD claim exit 0, to prove the cap is not
		// exit-mapped: the over-cap stdout is non-empty.
		acc.ExitAbsent = []int{0}

		values, _, err := readOnce(t, acc, art)
		if err == nil {
			t.Fatalf("a stdout of cap+1 bytes answered %#v; the overflow is an "+
				"execution failure from the bound, never an exit-map "+
				"consultation", values)
		}

		// The discriminating sibling: a stdout comfortably UNDER the cap
		// answers, under the same exit-map entry. Without it, "over the cap
		// refuses" is satisfied by a read that refuses every flood.
		ok := entry(fxRole, []string{bin, "flood", itoa(1024)})
		ok.Output = strptr("raw")
		ok.ExitAbsent = []int{0}
		if got := rawValue(t, ok, art); len(got.Value) != 1024 {
			t.Errorf("value length = %d; want 1024 — an under-cap stdout is a "+
				"value, and a non-empty stdout is never exit-mapped",
				len(got.Value))
		}
	})

	t.Run("exactly the cap is admitted", func(t *testing.T) {
		// The boundary's other side: the cap is inclusive, so a stdout of
		// exactly 1 MiB is a value, not a failure. Without this arm a
		// binding that refused at cap-1 would pass.
		acc := entry(fxRole, []string{bin, "flood", itoa(cmdbind.StdoutCap)})
		acc.Output = strptr("raw")

		got := rawValue(t, acc, art)
		if len(got.Value) != cmdbind.StdoutCap {
			t.Errorf("value length = %d; want exactly the cap %d",
				len(got.Value), cmdbind.StdoutCap)
		}
	})
}

// --- C4: the env allowlist (S4b) -----------------------------------------

// REQ-51: "env:      child env = allowlisted parent vars {PATH, HOME,
// TMPDIR, LANG, LC_*} + named `env_pass = [\"VAR\", ...]` vars + the
// entry's literal `env = { KEY = \"value\" }` + the overlay {INTRASTATE_ROLE,
// INTRASTATE_CAPABILITY, INTRASTATE_ACCESSOR, INTRASTATE_PROTOCOL=1};
// nothing else is inherited."
// REQ-47: "The protocol's version rides out-of-band as `INTRASTATE_PROTOCOL`
// in the child env"
// REQ-53: "`LC_*` is a literal prefix match on `LC_` (the one prefix rule;
// `env_pass` names whole variables and admits no pattern)."
// REQ-126 (S4b): "the injection never reaches the child (its env holds only
// the allowlist) ... Asserted on the child's observed environment, not on
// the read result, since an env defect can leave the value correct."
// ADVERSARIAL
//
// Asserted on the CHILD's observed environment, not on the read result,
// since an env defect can leave the value correct (S4b's own wording).
func TestReq51_TheChildEnvIsComposedFromTheAllowlistNotInherited(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	// The A7 hazard, set in the PARENT: the same binary and argv answer
	// differently under an inherited GIT_CONFIG_COUNT.
	t.Setenv("GIT_CONFIG_COUNT", "1")
	t.Setenv("GIT_CONFIG_KEY_0", "flow.status")
	t.Setenv("SSH_AUTH_SOCK", "/tmp/agent.sock")
	t.Setenv("LC_ALL", "C")
	t.Setenv("LCFOO", "nope")
	t.Setenv("HOME", "/home/fixture")

	acc := entry(fxRole, []string{bin, "echo-env"})
	acc.Output = strptr("raw")
	acc.Env = map[string]string{"TOOL_MODE": "strict"}
	acc.EnvPass = []string{"SSH_AUTH_SOCK"}

	got := rawValue(t, acc, art)
	env := observedEnv(t, got.Value)

	t.Run("the ambient injection never reaches the child", func(t *testing.T) {
		for _, k := range []string{"GIT_CONFIG_COUNT", "GIT_CONFIG_KEY_0"} {
			if v, held := env[k]; held {
				t.Errorf("the child observed %s=%q; nothing outside the "+
					"allowlist is inherited (A7)", k, v)
			}
		}
	})

	t.Run("the allowlist passes", func(t *testing.T) {
		if env["PATH"] == "" {
			t.Error("the child observed no PATH; argv0 resolves through the " +
				"parent's PATH (A2) and the allowlist passes it on")
		}
		if env["HOME"] != "/home/fixture" {
			t.Errorf("HOME = %q; want the parent's value", env["HOME"])
		}
	})

	t.Run("LC_ is a literal prefix and LCFOO is not covered", func(t *testing.T) {
		// REQ-53: the ONE prefix rule.
		if env["LC_ALL"] != "C" {
			t.Errorf("LC_ALL = %q; want the parent's %q — LC_* is allowlisted",
				env["LC_ALL"], "C")
		}
		if v, held := env["LCFOO"]; held {
			t.Errorf("the child observed LCFOO=%q; the prefix is the literal "+
				"`LC_` and admits no other match", v)
		}
	})

	t.Run("env_pass forwards a named variable", func(t *testing.T) {
		// REQ-6: the named, no-glob escape hatch.
		if env["SSH_AUTH_SOCK"] != "/tmp/agent.sock" {
			t.Errorf("SSH_AUTH_SOCK = %q; want the parent's value forwarded by "+
				"env_pass", env["SSH_AUTH_SOCK"])
		}
	})

	t.Run("the entry env literal reaches the child", func(t *testing.T) {
		if env["TOOL_MODE"] != "strict" {
			t.Errorf("TOOL_MODE = %q; want the entry's literal %q",
				env["TOOL_MODE"], "strict")
		}
	})

	t.Run("the INTRASTATE overlay is present and complete", func(t *testing.T) {
		want := map[string]string{
			cmdbind.EnvRole:       fxRole,
			cmdbind.EnvCapability: string(accessor.CapRead),
			cmdbind.EnvAccessor:   fxName,
			cmdbind.EnvProtocol:   cmdbind.ProtocolVersion,
		}
		for k, v := range want {
			if env[k] != v {
				t.Errorf("%s = %q; want %q — the overlay gives wrappers their "+
					"context without new placeholders", k, env[k], v)
			}
		}
	})
}

// REQ-52: "On a key collision the LATER layer wins, in exactly that order:
// overlay > entry `env` > `env_pass` > parent allowlist."
// BOUNDARY
//
// Each adjacent pair is asserted on its own, so a build that got one
// boundary right and another wrong cannot pass on the composed outcome.
func TestReq52_OnACollisionTheLaterLayerWins(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	t.Setenv("HOME", "/from/parent")
	t.Setenv("TOOL_MODE", "from-parent")

	t.Run("entry env beats the parent allowlist", func(t *testing.T) {
		acc := entry(fxRole, []string{bin, "echo-env"})
		acc.Output = strptr("raw")
		acc.Env = map[string]string{"HOME": "/from/entry"}

		env := observedEnv(t, rawValue(t, acc, art).Value)
		if env["HOME"] != "/from/entry" {
			t.Errorf("HOME = %q; want the ENTRY's value — entry env outranks the "+
				"parent allowlist", env["HOME"])
		}
	})

	t.Run("entry env beats env_pass", func(t *testing.T) {
		acc := entry(fxRole, []string{bin, "echo-env"})
		acc.Output = strptr("raw")
		acc.EnvPass = []string{"TOOL_MODE"}
		acc.Env = map[string]string{"TOOL_MODE": "from-entry"}

		env := observedEnv(t, rawValue(t, acc, art).Value)
		if env["TOOL_MODE"] != "from-entry" {
			t.Errorf("TOOL_MODE = %q; want the ENTRY's value — entry env "+
				"outranks env_pass", env["TOOL_MODE"])
		}
	})

	t.Run("env_pass beats the parent allowlist", func(t *testing.T) {
		// HOME is in the allowlist AND named in env_pass. Both resolve to
		// the same parent value, so the discriminating shape is that the
		// variable arrives exactly once with the parent's value.
		acc := entry(fxRole, []string{bin, "echo-env"})
		acc.Output = strptr("raw")
		acc.EnvPass = []string{"HOME"}

		env := observedEnv(t, rawValue(t, acc, art).Value)
		if env["HOME"] != "/from/parent" {
			t.Errorf("HOME = %q; want the parent's value", env["HOME"])
		}
	})
}

// REQ-51 (unset arms) — PHASE-0 READINGS.
// INPUT EDGE
//
// PHASE-0 READING (ASSUMPTION REQ-51): a var named in the allowlist but
// UNSET in the parent is simply not passed (no empty-string synthesis), and
// an `env_pass` naming an unset variable passes nothing and is NOT a load
// defect.
func TestReq51_AnUnsetAllowlistedOrPassedVariableIsNotSynthesized(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	t.Setenv("TMPDIR", "")
	_ = os.Unsetenv("TMPDIR")

	acc := entry(fxRole, []string{bin, "echo-env"})
	acc.Output = strptr("raw")
	acc.EnvPass = []string{"DEFINITELY_UNSET_FIXTURE_VAR"}

	env := observedEnv(t, rawValue(t, acc, art).Value)

	if v, held := env["DEFINITELY_UNSET_FIXTURE_VAR"]; held {
		t.Errorf("an env_pass name that is unset in the parent arrived as %q; "+
			"passing nothing is not synthesizing an empty string", v)
	}
	if v, held := env["TMPDIR"]; held && v == "" {
		t.Error("an unset allowlisted var arrived as an empty string; absence " +
			"is not an empty value")
	}
}

// --- C4: the deadline (S4) -----------------------------------------------

// REQ-49: "deadline: Setpgid; at ctx deadline Cancel = SIGKILL to -pgid;
// WaitDelay = 500ms bounds the stdin write and the pipe drain; timeout is
// classified from ctx.Err()"
// REQ-70: "the binding builds on `exec.CommandContext` and carries no timer
// of its own (A1)."
// REQ-125 (S4): "`timeout` refusal within `timeout + WaitDelay` in all
// three, no hang, and no process from the child's group surviving the
// refusal"
// ADVERSARIAL — normative fixture FX-deadline (S2/S4/S5 ≈1.0s, 0 survivors).
//
// The ablations (no process group, no WaitDelay) are deliberately NOT test
// arms: C4 states the triple as fixed behaviour with no injection seam, and
// adding one purely to disable it would put a way to weaken the deadline
// into the shipping binding. The suite asserts the ASSEMBLED triple.
func TestReq49_TheDeadlineBoundsTheWholeProcessGroupAndLeavesNoSurvivor(t *testing.T) {
	if testing.Short() {
		t.Skip("the deadline arms spawn real children")
	}
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	const bound = 300 * time.Millisecond
	// The bound the clause fixes: timeout + WaitDelay, plus slack for the
	// spawn itself. The executor owns the deadline, so the test supplies it.
	limit := bound + time.Duration(cmdbind.WaitDelay)*time.Millisecond + 2*time.Second

	t.Run("FX-deadline S5: a direct sleeper", func(t *testing.T) {
		acc := entry(fxRole, []string{bin, "sleep", "30s"})
		acc.Output = strptr("raw")

		ctx, cancel := context.WithTimeout(ctxOf(t), bound)
		defer cancel()

		r := cmdbind.Reader{Accessor: acc, Name: fxName, Config: allowed(t)}
		start := time.Now()
		_, _, err := r.Read(ctx, art, acc.Keys)
		elapsed := time.Since(start)

		if err == nil {
			t.Fatal("a child sleeping past its deadline answered a value")
		}
		if elapsed > limit {
			t.Errorf("the read returned after %v; want within timeout + "+
				"WaitDelay (%v) — a naive CommandContext blocks Wait past the "+
				"bound (FX-deadline S1)", elapsed, limit)
		}
		if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
			t.Error("the deadline did not fire; the timeout must be classified " +
				"from ctx.Err(), never from the wait error")
		}
	})

	t.Run("FX-deadline S2: a grandchild holds the stdout pipe", func(t *testing.T) {
		pidFile := filepath.Join(t.TempDir(), "grandchild.pid")
		acc := entry(fxRole, []string{bin, "orphan-holds-pipe", "30s", pidFile})
		acc.Output = strptr("raw")

		ctx, cancel := context.WithTimeout(ctxOf(t), bound)
		defer cancel()

		r := cmdbind.Reader{Accessor: acc, Name: fxName, Config: allowed(t)}
		start := time.Now()
		_, _, err := r.Read(ctx, art, acc.Keys)
		elapsed := time.Since(start)

		if err == nil {
			t.Fatal("the invocation answered a value despite the deadline")
		}
		if elapsed > limit {
			t.Errorf("the read returned after %v; want within %v — dropping "+
				"WaitDelay hangs Wait on the held pipe (FX-deadline S1)",
				elapsed, limit)
		}
		// The load-bearing half: NO process from the child's GROUP survives.
		// Dropping the group signal orphans this grandchild (FX-deadline S3).
		pid := pidFrom(t, pidFile)
		time.Sleep(100 * time.Millisecond)
		if alive(pid) {
			t.Errorf("the grandchild pid %d survived the refusal; the deadline "+
				"must terminate the child's process GROUP, not only the direct "+
				"child (FX-deadline S3's orphan)", pid)
		}
	})

	t.Run("FX-deadline S4: a child that never reads stdin", func(t *testing.T) {
		// A write, so a real stdin payload is offered. WaitDelay bounds the
		// stdin write, so a non-reading child cannot block the parent.
		big := make([]resolve.Tag, 0, 4096)
		filler := strings.Repeat("x", 256)
		for i := range 4096 {
			big = append(big, resolve.Tag{Key: fxKey + itoa(i), Value: filler})
		}
		acc := entry(fxRole, []string{bin, "never-reads-stdin", "30s"})
		acc.ReadBack = true

		ctx, cancel := context.WithTimeout(ctxOf(t), bound)
		defer cancel()

		w := &cmdbind.Writer{Accessor: acc, Name: fxName, Config: allowed(t)}
		start := time.Now()
		err := w.Apply(ctx, art, big)
		elapsed := time.Since(start)

		if err == nil {
			t.Fatal("the write answered success despite the deadline")
		}
		if elapsed > limit {
			t.Errorf("Apply returned after %v; want within %v — WaitDelay bounds "+
				"the stdin write so a non-reading child cannot hang the CLI",
				elapsed, limit)
		}

		// The discriminating sibling: a child that DOES read the same large
		// stdin completes well inside the bound. Without it, "the write
		// refuses" is satisfied by a binding that refuses every write.
		sink := filepath.Join(t.TempDir(), "stdin.json")
		ok := entry(fxRole, []string{bin, "stdin-to-file", sink})
		ok.ReadBack = true
		okCtx, okCancel := context.WithTimeout(ctxOf(t), 10*time.Second)
		defer okCancel()
		w2 := &cmdbind.Writer{Accessor: ok, Name: fxName, Config: allowed(t)}
		if aerr := w2.Apply(okCtx, art, big); aerr != nil {
			t.Fatalf("a child that reads its stdin refused: %v", aerr)
		}
		if _, serr := os.Stat(sink); serr != nil {
			t.Errorf("the reading child recorded no stdin: %v", serr)
		}
	})
}

// --- C4: the platform refusal (S6b) --------------------------------------

// REQ-60: "platform: a runtime `runtime.GOOS` check in the command binding
// refuses `execution_failure` before spawn on non-Unix — the predicate is
// `goos == \"windows\" || goos == \"js\" || goos == \"plan9\"` (refuse-listed,
// not allow-listed)"
// REQ-61: "The check reads an injectable package-level `goos` var
// (defaulting to `runtime.GOOS`) so the refusal is testable on Unix CI"
// REQ-134: the same, as a Failure Mode.
// REQ-130 (S6b): "`execution_failure` naming the unsupported platform, with
// no child spawned (asserted by absence of a spawn, as in scenario 7); the
// same model passes `intrastate lint` under that setting, since lint is
// platform-neutral."
// DOMAIN EDGE
//
// PHASE-0 READING (ASSUMPTION REQ-60/REQ-61): the `goos` var is unexported
// and lives in this package, so this test sets it through the package's own
// test hook. Exporting it would put a runtime weakening seam in the public
// surface.
func TestReq60_ACommandEntryRefusesOnARefuseListedPlatformBeforeSpawn(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	for _, os := range []string{"windows", "js", "plan9"} {
		t.Run(os, func(t *testing.T) {
			cmdbind.SetGOOSForTest(t, os)

			trace := filepath.Join(t.TempDir(), "spawned")
			acc := entry(fxRole, []string{bin, "trace", trace})
			acc.Output = strptr("raw")

			_, _, err := readOnce(t, acc, art)
			if err == nil {
				t.Fatal("a command entry ran on a refuse-listed platform")
			}
			var ee *accessor.ExecError
			if !errors.As(err, &ee) {
				t.Fatalf("the refusal is not an *accessor.ExecError: %v", err)
			}
			if !strings.Contains(ee.Detail, os) {
				t.Errorf("Detail = %q; want it to name the unsupported platform "+
					"%q", ee.Detail, os)
			}
			if _, serr := stat(trace); serr == nil {
				t.Error("a child ran; the platform refusal is BEFORE spawn")
			}
		})
	}

	t.Run("lint stays platform-neutral in the same binary", func(t *testing.T) {
		// REQ-130 / REQ-134: "the same model passes `intrastate lint` under
		// that setting, since lint is platform-neutral" — a model authored
		// on Windows still validates. A BUILD CONSTRAINT would have made
		// the refusal unbuildable-on-Windows rather than observable, and it
		// would have taken lint with it; the runtime check is what keeps
		// the two separable, so the separation is what this asserts.
		cmdbind.SetGOOSForTest(t, "windows")

		if _, err := table.Load([]byte(platformNeutralModel(bin)), "win.toml"); err != nil {
			t.Errorf("a command model failed to LOAD under GOOS=windows: %v — "+
				"validation is platform-neutral and the refusal is a runtime "+
				"check in the binding, not a build constraint", err)
		}
	})

	t.Run("an unlisted Unix is admitted without enumeration", func(t *testing.T) {
		// The predicate is refuse-listed, so a BSD — which nothing
		// enumerates — must run. An allow-list implementation fails here.
		cmdbind.SetGOOSForTest(t, "freebsd")

		acc := entry(fxRole, []string{bin, "emit", "draft\n", "0"})
		acc.Output = strptr("raw")

		if got := rawValue(t, acc, art); got.Value != "draft" {
			t.Errorf("value = %q on GOOS=freebsd; every Unix that supports the "+
				"process-group mechanism is admitted without being enumerated",
				got.Value)
		}
	})
}

func stat(p string) (any, error) { return os.Stat(p) }

// --- C6: the gate, at the binding ----------------------------------------

// REQ-92: "absent ⇒ every command invocation refuses execution_failure
// before spawn, Detail naming the gate"
// REQ-96: "The gate sites **in the command binding's constructor**, not in
// the executor"
// REQ-97: "With the gate off, a command invocation refuses before spawn
// (`execution_failure`, `Detail` naming `allow_commands`)"
// ADVERSARIAL
//
// Asserted by ABSENCE OF A SPAWN — a sentinel argv0 that would leave an
// observable trace if executed — not merely by a non-zero exit (S7's own
// wording).
func TestReq92_WithTheGateOffEveryCapabilityRefusesBeforeSpawn(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	off := cmdbind.Config{BaseDir: filepath.Dir(bin), AllowCommands: false}

	t.Run("read", func(t *testing.T) {
		trace := filepath.Join(t.TempDir(), "spawned")
		acc := entry(fxRole, []string{bin, "trace", trace})
		acc.Output = strptr("raw")

		r := cmdbind.Reader{Accessor: acc, Name: fxName, Config: off}
		_, _, err := r.Read(ctxOf(t), art, acc.Keys)

		requireGateRefusal(t, err, trace)
	})

	t.Run("gate", func(t *testing.T) {
		trace := filepath.Join(t.TempDir(), "spawned")
		acc := entry(fxRole, []string{bin, "trace", trace})

		g := cmdbind.Gate{Accessor: acc, Name: fxName, Config: off}
		_, _, err := g.Gate(ctxOf(t), art)

		requireGateRefusal(t, err, trace)
	})

	t.Run("write", func(t *testing.T) {
		trace := filepath.Join(t.TempDir(), "spawned")
		acc := entry(fxRole, []string{bin, "trace", trace})
		acc.ReadBack = true

		w := &cmdbind.Writer{Accessor: acc, Name: fxName, Config: off}
		err := w.Apply(ctxOf(t), art, []resolve.Tag{{Key: fxKey, Value: "final"}})

		requireGateRefusal(t, err, trace)
	})

	t.Run("the same entry RUNS with the gate on", func(t *testing.T) {
		// The discriminating half: without it, "refuses" is satisfied by a
		// binding that never runs anything.
		trace := filepath.Join(t.TempDir(), "spawned")
		acc := entry(fxRole, []string{bin, "trace", trace})
		acc.Output = strptr("raw")

		r := cmdbind.Reader{Accessor: acc, Name: fxName, Config: allowed(t)}
		_, _, _ = r.Read(ctxOf(t), art, acc.Keys)

		if _, serr := os.Stat(trace); serr != nil {
			t.Errorf("no child ran with --allow-commands set; the gate must be "+
				"the only thing that stopped it: %v", serr)
		}
	})
}

func requireGateRefusal(t *testing.T, err error, trace string) {
	t.Helper()

	if err == nil {
		t.Fatal("the invocation did not refuse with the gate off")
	}
	var ee *accessor.ExecError
	if !errors.As(err, &ee) {
		t.Fatalf("the refusal is not an *accessor.ExecError: %v", err)
	}
	if !strings.Contains(ee.Detail, "allow") {
		t.Errorf("Detail = %q; want it to name the `allow_commands` gate",
			ee.Detail)
	}
	if _, serr := os.Stat(trace); serr == nil {
		t.Error("a child ran; the gate refuses BEFORE spawn — asserted by the " +
			"absence of a spawn, not by a non-zero exit")
	}
}

// platformNeutralModel is a valid command-backed model, used to assert that
// LOADING is unaffected by the platform the binding refuses on.
func platformNeutralModel(bin string) string {
	return `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "winflow"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true
required = true

[read.state]
role = "state"
command = ["` + bin + `", "emit", "draft", "0"]
output = "raw"
keys = ["status"]
timeout = "2s"

[write.state]
role = "state"
command = ["` + bin + `", "stdin-to-file", "{artifact}"]
keys = ["status"]
timeout = "2s"
read_back = true

[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "advance-draft"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"
`
}
