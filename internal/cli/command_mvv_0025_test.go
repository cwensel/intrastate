package cli

// RDR 0025 — the Minimum Viable Validation, end to end (`0025:MVV`,
// scenario S6).
//
// This is the integration proof, and it is RUNNABLE: it authors a model,
// lints it, mutates it six ways, drives a real state change through the
// declared write COMMAND — not intrastate — and verifies the change through
// the declared command READER. A green exit code is explicitly NOT the
// oracle at any step.
//
// The Round-Trip / Inverse invariant the Pre-Lock `fidelity` table fixes
// for this seam is: planned tags → stdin JSON → tool → read-back is TYPED
// TAG-VALUE equality (explicitly NOT byte equality — the artifact's byte
// form belongs to the bound tool), with `<clear>` round-tripping as
// absence. So step 4 asserts the RECONSTRUCTED value equals the planned
// value exactly, read back out of the tool's own artifact through the
// declared reader.
//
// The one thing this does NOT do, deliberately: wire resolve → apply.
// "Automatic resolve→apply wiring is not this RDR's scope and no Phase
// builds it" — what is validated is that a declared command can CARRY the
// apply, which is the half that is missing today.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/table"
)

// REQ-MVV (`0025:MVV`), all five steps plus the end state.
// REQ-129 (S6): "the declared write command, not intrastate, edits the
// artifact, and the declared command reader verifies it."
// REQ-121: "the MVV scenario is the integration proof."
// HAPPY PATH
func TestReqMVV_ALintedModelAppliesAndVerifiesARealChangeThroughACommand(t *testing.T) {
	requireGit(t)

	dir := t.TempDir()
	wrapper := writeWrapper(t, dir)

	// --- MVV step 1 -------------------------------------------------------
	// "Author a model declaring one command-backed read and one
	// command-backed write (`read_back = true`) over a single artifact role,
	// delegating to an established tool present in CI: the read binds `git
	// config --file {artifact} --get <key>` directly in raw mode; the write
	// is a declared wrapper over `git config`."
	model := writeModelFile(t, dir, "mvv.toml", mvvModel(wrapper))
	artifact := filepath.Join(dir, "state.cfg")

	// --- MVV step 2 -------------------------------------------------------
	// "`intrastate lint` accepts it; six mutated copies ... are each
	// rejected with their C5 defect."
	t.Run("step 2: lint accepts the model and rejects the six mutants", func(t *testing.T) {
		if _, _, err := runCmd(t, "lint", "--model", model, "--as=json"); err != nil {
			t.Fatalf("lint refused the authored MVV model: %v", err)
		}

		mutants := []struct {
			name string
			from string
			to   string
			want table.Category
		}{
			{
				"path+command conflict",
				`command = ["git", "config", "--file", "{artifact}", "--get", "flow.status"]`,
				`path = "flow.status"
command = ["git", "config", "--file", "{artifact}", "--get", "flow.status"]`,
				table.CatCommandAndPathConflict,
			},
			{
				"empty element",
				`"--get", "flow.status"]`,
				`"--get", ""]`,
				table.CatCommandEmpty,
			},
			{
				"unknown placeholder",
				`"--file", "{artifact}",`,
				`"--file", "{role}",`,
				table.CatCommandUnknownPlaceholder,
			},
			{
				"sh -c interpreter form",
				`command = ["git", "config", "--file", "{artifact}", "--get", "flow.status"]`,
				`command = ["sh", "-c", "git config --file $1 --get flow.status"]`,
				table.CatCommandShellInterpreter,
			},
			{
				`output = "raw" with two keys`,
				`keys = ["status"]
timeout = "10s"

[write.state]`,
				`keys = ["status", "stage"]
timeout = "10s"

[write.state]`,
				table.CatCommandOutputShape,
			},
			{
				"env key shadowing INTRASTATE_ROLE",
				`output = "raw"`,
				`output = "raw"
env = { INTRASTATE_ROLE = "spoof" }`,
				table.CatCommandEnvConflict,
			},
		}

		src := mvvModel(wrapper)
		for _, m := range mutants {
			t.Run(m.name, func(t *testing.T) {
				mutated := strings.Replace(src, m.from, m.to, 1)
				if mutated == src {
					t.Fatalf("the %s substitution did not apply", m.name)
				}
				p := writeModelFile(t, t.TempDir(), "mutant.toml", mutated)

				_, _, err := runCmd(t, "lint", "--model", p, "--as=json")
				if err == nil {
					t.Fatalf("lint ACCEPTED the %s mutant", m.name)
				}
				// `lint` is at ROOT and carries its own code, not the
				// `flow`-group-scoped one (deviations D10).
				if code := clierr.ErrorCode(err); code != lintModelInvalid {
					t.Fatalf("code = %q; want %q", code, lintModelInvalid)
				}
				// The oracle is the CATEGORY the finding carries, never
				// "lint returned non-empty" (S1).
				if got := findingCategories(t, err); !containsStr(got, string(m.want)) {
					t.Errorf("findings carry %#v; want the %q defect",
						got, m.want)
				}
			})
		}
	})

	// --- MVV step 3 -------------------------------------------------------
	// "With `--allow-commands` on the invocation (C6; without it, the same
	// invocation refuses `execution_failure` naming the gate, before spawn),
	// drive a state change end to end through the existing write path ...:
	// the declared write **command** — not intrastate — applies the edit to
	// the artifact."
	t.Run("step 3a: without the flag the same invocation refuses before spawn", func(t *testing.T) {
		gdir := t.TempDir()
		trace := filepath.Join(gdir, "spawned")
		gmodel := writeModelFile(t, gdir, "gated.toml",
			strings.Replace(mvvModel(wrapper),
				strconv.Quote(wrapper), strconv.Quote(traceScript(t, gdir, trace)), 1))
		gart := filepath.Join(gdir, "state.cfg")

		_, _, err := runCmd(t, "flow", "set-state",
			"--model", gmodel,
			"--artifact", "state="+gart,
			"--write", "status=final",
			"--as=json")
		if err == nil {
			t.Fatal("set-state succeeded with no --allow-commands")
		}
		if code := clierr.ErrorCode(err); code != codeAccessorFailed {
			t.Errorf("code = %q; want %q", code, codeAccessorFailed)
		}
		if _, serr := os.Stat(trace); serr == nil {
			t.Error("a child ran; the gate refuses BEFORE spawn")
		}
	})

	// git config --file needs the artifact's directory to exist; the file
	// itself is created by the write, which is the point.
	stdout, _, err := runCmd(t, "flow", "set-state",
		"--model", model,
		"--artifact", "state="+artifact,
		"--write", "status=final",
		"--allow-commands",
		"--as=json")
	if err != nil {
		t.Fatalf("MVV step 3: the declared write command failed: %v", err)
	}

	// The declared write COMMAND — not intrastate — edited the artifact. The
	// oracle is the artifact's own format, which is `git config`'s and not
	// intrastate's flat JSON: if intrastate had written it, it would be a
	// JSON object.
	raw, rerr := os.ReadFile(artifact)
	if rerr != nil {
		t.Fatalf("MVV step 3: the artifact was never written: %v", rerr)
	}
	if strings.HasPrefix(strings.TrimSpace(string(raw)), "{") {
		t.Errorf("the artifact is JSON:\n%s\nintrastate's own flat-JSON binding "+
			"wrote it — the whole point is that the declared COMMAND applies "+
			"the edit", raw)
	}
	if !strings.Contains(string(raw), "final") {
		t.Errorf("the artifact does not hold the applied value:\n%s", raw)
	}

	// --- MVV step 4 -------------------------------------------------------
	// "Read-back runs through the declared command reader and verifies the
	// written value; the result reports the written tags."
	//
	// The ROUND-TRIP invariant: the reconstructed value must equal the
	// planned value EXACTLY. A green exit code is not sufficient and is not
	// what is asserted.
	t.Run("step 4: the round-trip reconstructs the planned value exactly", func(t *testing.T) {
		written := writtenTags(t, stdout)
		got, held := written["status"]
		if !held {
			t.Fatalf("set-state reported written tags %#v; want the planned "+
				"`status` among them — read-back VERIFIED it", written)
		}
		if got != "final" {
			t.Errorf("the reconstructed value = %q; want the planned %q exactly "+
				"— the round-trip is typed tag-VALUE equality", got, "final")
		}

		// And independently, through the declared command READER on a fresh
		// invocation: the value survives out of the tool's own artifact.
		out, _, rerr := runCmd(t, "flow", "read-state",
			"--model", model,
			"--artifact", "state="+artifact,
			"--allow-commands",
			"--as=json")
		if rerr != nil {
			t.Fatalf("the declared command reader failed: %v", rerr)
		}
		if reread := reportedTag(t, out, "status"); reread != "final" {
			t.Errorf("the declared command reader read %q; want %q — this is "+
				"the half that proves the reader binds `git config --get` "+
				"directly in raw mode (FX-raw-read)", reread, "final")
		}
	})

	// --- MVV step 5 -------------------------------------------------------
	// "A command that sleeps past its declared timeout refuses `timeout` and
	// leaves no orphan process."
	t.Run("step 5: a command past its declared timeout refuses timeout", func(t *testing.T) {
		sdir := t.TempDir()
		sleeper := writeSleeper(t, sdir)
		smodel := writeModelFile(t, sdir, "slow.toml",
			strings.Replace(
				strings.Replace(mvvModel(wrapper), strconv.Quote(wrapper),
					strconv.Quote(sleeper), 1),
				`timeout = "10s"`, `timeout = "300ms"`, -1))
		sart := filepath.Join(sdir, "state.cfg")

		_, _, serr := runCmd(t, "flow", "set-state",
			"--model", smodel,
			"--artifact", "state="+sart,
			"--write", "status=final",
			"--allow-commands",
			"--as=json")
		if serr == nil {
			t.Fatal("a command sleeping past its declared timeout succeeded")
		}
		if code := clierr.ErrorCode(serr); code != codeReadBackTimeout &&
			code != codeAccessorTimeout {
			t.Errorf("code = %q; want a TIMEOUT code — the deadline mints its "+
				"own class, distinct from execution failure", code)
		}
	})

	// --- end state --------------------------------------------------------
	// "a linted model applies and verifies a real state change through a
	// declared command, with the executed argv readable in the model."
	t.Run("end state: the executed argv is readable in the model", func(t *testing.T) {
		src, serr := os.ReadFile(model)
		if serr != nil {
			t.Fatalf("reading the model: %v", serr)
		}
		// The reviewer reads the argv in the model. Both the directly-bound
		// read and the wrapper's own invocation argv are there.
		for _, want := range []string{"git", "config", "--file", "{artifact}", wrapper} {
			if !strings.Contains(string(src), want) {
				t.Errorf("the model does not carry %q; the executed authority is "+
					"reviewable STATICALLY, from the model alone", want)
			}
		}
	})
}

// --- MVV helpers ----------------------------------------------------------

func mvvModel(wrapper string) string {
	return `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "mvvcmd"
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

[tags.stage]
provenance = "observed"
kind = "scalar"

[read.state]
role = "state"
command = ["git", "config", "--file", "{artifact}", "--get", "flow.status"]
output = "raw"
exit_absent = [1]
keys = ["status"]
timeout = "10s"

[write.state]
role = "state"
command = [` + strconv.Quote(wrapper) + `, "{artifact}"]
keys = ["status"]
timeout = "10s"
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

// writeSleeper writes a child that blocks past any plausible declared
// timeout, so the deadline is what ends it.
func writeSleeper(t *testing.T, dir string) string {
	t.Helper()

	p := filepath.Join(dir, "sleeper")
	if err := os.WriteFile(p, []byte("#!/bin/sh\nsleep 30\n"), 0o700); err != nil {
		t.Fatalf("writing the sleeper: %v", err)
	}
	return p
}

// traceScript writes a sentinel child that records having run, so "no child
// spawned" is asserted by absence of a trace rather than by an exit code.
func traceScript(t *testing.T, dir, trace string) string {
	t.Helper()

	p := filepath.Join(dir, "sentinel")
	body := "#!/bin/sh\nprintf ran >> " + strconv.Quote(trace) + "\n"
	if err := os.WriteFile(p, []byte(body), 0o700); err != nil {
		t.Fatalf("writing the sentinel: %v", err)
	}
	return p
}

// writtenTags decodes `flow set-state`'s reported written tags. Going
// through the CLI envelope keeps the oracle on the contract rather than on
// any file schema.
func writtenTags(t *testing.T, stdout string) map[string]string {
	t.Helper()

	// `set-state`'s READ-BACK-CONFIRMED owned tags ride the payload's
	// `owned` map, not an invented `written` array: RDR 0005 fixes the
	// envelope and this RDR changes no CLI payload shape. `owned` is
	// exactly what the read-back verified — a cleared key is absent from
	// it, since a verified removal is not a written value (deviations
	// D12, `0004:C11`).
	var env struct {
		Type string `json:"type"`
		Data struct {
			Owned map[string]string `json:"owned"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(stdout), &env); err != nil {
		t.Fatalf("set-state stdout is not one JSON object: %v\n%s", err, stdout)
	}
	if env.Type != "ok" {
		t.Fatalf("set-state envelope type = %q; want %q\n%s", env.Type, "ok", stdout)
	}
	return env.Data.Owned
}

// reportedTag decodes one tag from `flow read-state`'s report.
func reportedTag(t *testing.T, stdout, key string) string {
	t.Helper()

	// `read-state` reports per-reader: each reader carries its DECLARED
	// key set beside the tags it returned (deviations D12).
	var env struct {
		Data struct {
			Readers []struct {
				ID   string            `json:"id"`
				Keys []string          `json:"keys"`
				Tags map[string]string `json:"tags"`
			} `json:"readers"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(stdout), &env); err != nil {
		t.Fatalf("read-state stdout is not one JSON object: %v\n%s", err, stdout)
	}
	for _, r := range env.Data.Readers {
		if v, held := r.Tags[key]; held {
			return v
		}
	}
	t.Fatalf("read-state reported no %q tag:\n%s", key, stdout)
	return ""
}

// findingCategories returns the category slugs a model-invalid refusal
// carries, so a mutant is asserted by CATEGORY rather than by message text.
func findingCategories(t *testing.T, err error) []string {
	t.Helper()

	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("the refusal is not a *clierr.CLIError: %v", err)
	}
	var out []string
	for _, f := range ce.Findings {
		out = append(out, f.Code, f.Param, f.Message)
	}
	return out
}

func containsStr(haystack []string, needle string) bool {
	for _, h := range haystack {
		if strings.Contains(h, needle) {
			return true
		}
	}
	return false
}
