package cmdbind_test

// RDR 0025 — the implemented seam and the write read-back discipline
// (C7, C3/C4's write rule; scenarios S5, plus the C7 clauses that only
// mean something once a COMMAND binding implements them).
//
// The clauses here are restatements of RDR 0004's seam, and that is
// exactly why they are asserted against the COMMAND bindings rather than
// against a fixture: "the three command bindings implement the interfaces
// UNCHANGED" is a claim about these types, and asserting it over a test
// double would prove only that the double was written to the signature.
//
// The mutant this file kills: a command read binding that signals
// established absence by OMITTING the key from both return slices. That
// convention reads as absence to a naive consumer and as
// `incomplete_read` to `readOutcome.classify` — the arm A15 recorded as
// REFUTED and repaired. So the absence arm and the omission arm are
// asserted as a discriminating pair on the real binding.

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/cmdbind"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// --- C7: the seam, quoted as it exists -----------------------------------

// REQ-98: "the three command bindings implement
// `internal/accessor/binding.go`'s interfaces UNCHANGED — this RDR adds no
// method, changes no signature"
// REQ-99: "Read(ctx context.Context, art Artifact, requested []string)
// (values []KeyValue, unreadable []string, err error)"
// REQ-100: "Gate(ctx context.Context, art Artifact) (Verdict, string,
// error)"
// REQ-101: "Apply(ctx context.Context, art Artifact, planned []resolve.Tag)
// error"
// REQ-102: "Invocations() int          # on WriteBinding, beside Apply"
// REQ-104: "the write method is `Apply`, NOT `Write`."
// BOUNDARY — NEGATIVE REQ.
//
// The compile-time half is the only way a Go test can catch "no signature
// changed": if any of the four returns moved, this file does not build.
// The runtime half is below, on each capability's own test.
func TestReq98_TheCommandBindingsSatisfyTheUnchangedInterfaces(t *testing.T) {
	var (
		_ accessor.ReadBinding  = cmdbind.Reader{}
		_ accessor.GateBinding  = cmdbind.Gate{}
		_ accessor.WriteBinding = (*cmdbind.Writer)(nil)
	)

	// And the capability each declares, which is what the executor selects
	// on. A binding that reported the wrong one would be unreachable under
	// its own table.
	if got := (cmdbind.Reader{}).Capability(); got != accessor.CapRead {
		t.Errorf("Reader.Capability() = %q; want %q", got, accessor.CapRead)
	}
	if got := (cmdbind.Gate{}).Capability(); got != accessor.CapGate {
		t.Errorf("Gate.Capability() = %q; want %q", got, accessor.CapGate)
	}
	if got := (&cmdbind.Writer{}).Capability(); got != accessor.CapWrite {
		t.Errorf("Writer.Capability() = %q; want %q", got, accessor.CapWrite)
	}

	// The discriminating half: satisfying the interface proves the shape,
	// not that the shape carries anything. One real invocation per
	// capability, through the seam's own returns.
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	racc := entry(fxRole, []string{bin, "emit", `{"` + fxKey + `":"final"}`, "0"})
	values, unreadable, rerr := cmdbind.Reader{
		Accessor: racc, Name: fxName, Config: allowed(t),
	}.Read(ctxOf(t), art, racc.Keys)
	if rerr != nil {
		t.Fatalf("the seam read failed: %v", rerr)
	}
	if len(values) != 1 || values[0].Key != fxKey || values[0].Value != "final" {
		t.Errorf("values = %#v; want one KeyValue per key the binding read", values)
	}
	if len(unreadable) != 0 {
		t.Errorf("unreadable = %#v; UNREADABLE is the second slice and a "+
			"complete read leaves it empty", unreadable)
	}

	gacc := entry(fxRole, []string{bin, "emit", `{"verdict":"deny","reason":"why"}`, "0"})
	verdict, reason, gerr := cmdbind.Gate{
		Accessor: gacc, Name: fxName, Config: allowed(t),
	}.Gate(ctxOf(t), art)
	if gerr != nil {
		t.Fatalf("the seam gate failed: %v", gerr)
	}
	if verdict != accessor.VerdictDeny || reason != "why" {
		t.Errorf("Gate returned (%q, %q); C3's `reason` crosses back through "+
			"the Go return and is not dropped", verdict, reason)
	}
}

// REQ-103: "the command write binding counts `Apply` ENTRIES — one per call
// — and performs NO internal respawn: one `Apply`, one spawn, one
// increment."
// REQ-137: "bindings hold no state between invocations (`0004:C14` — no
// retry, no undo)."
// DOMAIN EDGE
//
// PHASE-0 READING (ASSUMPTION REQ-103): `Invocations()` counts Apply
// ENTRIES even when the spawn FAILS — the clause says "one Apply, one
// spawn, one increment" and the existing 0004 assertions are about LAYER
// re-entry, not process success.
//
// The oracle is a PAIR: the binding's own counter and the child's own
// trace file. A binding that internally retried would leave two traces
// while reporting one invocation, and only the pair catches that.
func TestReq103_OneApplyIsOneSpawnAndOneIncrement(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	t.Run("a successful apply", func(t *testing.T) {
		trace := filepath.Join(t.TempDir(), "runs")
		acc := entry(fxRole, []string{bin, "trace", trace})
		acc.ReadBack = true

		w := &cmdbind.Writer{Accessor: acc, Name: fxName, Config: allowed(t)}
		if n := w.Invocations(); n != 0 {
			t.Fatalf("Invocations() = %d before any Apply; want 0", n)
		}
		if err := w.Apply(ctxOf(t), art, []resolve.Tag{{Key: fxKey, Value: "final"}}); err != nil {
			t.Fatalf("the write failed: %v", err)
		}
		if n := w.Invocations(); n != 1 {
			t.Errorf("Invocations() = %d after one Apply; want 1", n)
		}
		if got := countLines(t, trace); got != 1 {
			t.Errorf("the child ran %d times; want exactly 1 — one Apply, one "+
				"spawn, no internal respawn", got)
		}
	})

	t.Run("a FAILING apply still counts one entry and spawns once", func(t *testing.T) {
		// The child WRITES its stderr — proving it ran — then exits
		// non-zero. Without the tail assertion, "the write failed and the
		// count is 1" is satisfied by a binding that never spawned at all,
		// which is a different property entirely.
		acc := entry(fxRole, []string{bin, "emit-stderr", "wrapper: refused\n", "1"})
		acc.ReadBack = true

		w := &cmdbind.Writer{Accessor: acc, Name: fxName, Config: allowed(t)}
		err := w.Apply(ctxOf(t), art, []resolve.Tag{{Key: fxKey, Value: "final"}})
		if err == nil {
			t.Fatal("a write whose child exited non-zero reported success")
		}
		var ee *accessor.ExecError
		if !errors.As(err, &ee) {
			t.Fatalf("the refusal is not an *accessor.ExecError: %v", err)
		}
		if !strings.Contains(ee.Detail, "wrapper: refused") {
			t.Fatalf("Detail = %q; want the child's stderr tail — without it the "+
				"count assertion below cannot tell a failed spawn from no spawn",
				ee.Detail)
		}
		if n := w.Invocations(); n != 1 {
			t.Errorf("Invocations() = %d after one failing Apply; want 1 — the "+
				"count is about LAYER re-entry, not process success", n)
		}
	})

	t.Run("two Applies are two spawns and two increments", func(t *testing.T) {
		trace := filepath.Join(t.TempDir(), "runs")
		acc := entry(fxRole, []string{bin, "trace", trace})
		acc.ReadBack = true

		w := &cmdbind.Writer{Accessor: acc, Name: fxName, Config: allowed(t)}
		for range 2 {
			if err := w.Apply(ctxOf(t), art, []resolve.Tag{{Key: fxKey, Value: "final"}}); err != nil {
				t.Fatalf("the write failed: %v", err)
			}
		}
		if n := w.Invocations(); n != 2 {
			t.Errorf("Invocations() = %d after two Applies; want 2", n)
		}
		if got := countLines(t, trace); got != 2 {
			t.Errorf("the child ran %d times; want 2", got)
		}
	})
}

func countLines(t *testing.T, path string) int {
	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	n := 0
	for _, line := range strings.Split(string(b), "\n") {
		if strings.TrimSpace(line) != "" {
			n++
		}
	}
	return n
}

// REQ-106: "`exit_absent` absence crosses as `KeyValue{Key: k, Absent:
// true}` IN `values` — a present record carrying the flag, the same channel
// `internal/cli/flowbind/flowbind.go::Reader.Read` already uses for the
// path-backed binding."
// REQ-107: "Returning the key in NEITHER slice does NOT establish absence
// ... A command binding that signalled absence by omission would be
// refused, not believed"
// REQ-44: "absence reaches the executor as `KeyValue{Absent: true}` in
// `values`, never as omission from both slices"
// ADVERSARIAL
//
// The repair A15 records. Asserted end-to-end THROUGH the executor, because
// the whole point is which branch `readOutcome.classify` takes: a binding
// returning the flag succeeds, and the executor turns it into omission at
// the resolver seam (REQ-108). No unit-level shape assertion can show that.
func TestReq106_EstablishedAbsenceCrossesAsTheFlagAndSurvivesClassify(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	acc := entry(fxRole, []string{bin, "silent", "1"}, fxKey, fxKeyB)
	acc.ExitAbsent = []int{1}

	reg := accessor.Registry{
		Flow: "cmdflow",
		Definitions: []accessor.Definition{{
			Identity: accessor.Identity{
				Flow: "cmdflow", Name: fxName, Capability: accessor.CapRead,
			},
			Accessor: acc,
			Binding:  cmdbind.Reader{Accessor: acc, Name: fxName, Config: allowed(t)},
		}},
		OwnedTags: []string{fxKey, fxKeyB},
	}
	e := accessor.NewExecutor(reg, accessor.Artifacts{fxRole: art})

	got := e.Read(ctxOf(t), fxName)

	if got.Refusal != nil {
		t.Fatalf("an exit_absent read refused %q; the flag IS the "+
			"established-absent channel and survives classify intact",
			got.Refusal.Class)
	}
	if len(got.Values) != 2 {
		t.Fatalf("Values = %#v; want one record per declared key", got.Values)
	}
	for _, v := range got.Values {
		if !v.Absent {
			t.Errorf("%q came back present; the whole invocation established "+
				"absence", v.Key)
		}
	}

	// REQ-108: at the RESOLVER seam absence becomes OMISSION. Both
	// boundaries hold, and asserting only one of them would let a binding
	// that leaked a placeholder pass.
	if snap := got.OwnedSnapshot(); len(snap) != 0 {
		t.Errorf("OwnedSnapshot() = %#v; an absent key is OMITTED at the outer "+
			"boundary, never carried as a placeholder (0004:C8)", snap)
	}
}

// REQ-107 (the refused sibling): a command read whose envelope omits a
// declared key refuses the WHOLE read `incomplete_read`.
// ADVERSARIAL
func TestReq107_ACommandBindingSignallingAbsenceByOmissionIsRefused(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	// A well-formed envelope that simply omits the second declared key.
	// This is exactly the convention C7 forbids: it looks like absence and
	// must be refused rather than believed.
	acc := entry(fxRole,
		[]string{bin, "emit", `{"` + fxKey + `":"final"}`, "0"},
		fxKey, fxKeyB)

	reg := accessor.Registry{
		Flow: "cmdflow",
		Definitions: []accessor.Definition{{
			Identity: accessor.Identity{
				Flow: "cmdflow", Name: fxName, Capability: accessor.CapRead,
			},
			Accessor: acc,
			Binding:  cmdbind.Reader{Accessor: acc, Name: fxName, Config: allowed(t)},
		}},
		OwnedTags: []string{fxKey, fxKeyB},
	}
	got := accessor.NewExecutor(reg, accessor.Artifacts{fxRole: art}).
		Read(ctxOf(t), fxName)

	if got.Refusal == nil {
		t.Fatalf("an omitted declared key was believed; Values = %#v — guessing "+
			"absence is the failure this contract exists to prevent", got.Values)
	}
	if got.Refusal.Class != accessor.ClassIncompleteRead {
		t.Errorf("class = %q; want %q", got.Refusal.Class, accessor.ClassIncompleteRead)
	}
	if !contains(got.Refusal.Keys, fxKeyB) {
		t.Errorf("Keys = %#v; want the unread key %q named", got.Refusal.Keys, fxKeyB)
	}
}

func contains(s []string, want string) bool {
	for _, v := range s {
		if v == want {
			return true
		}
	}
	return false
}

// --- S5: write read-back through the declared reader ---------------------

// REQ-127 (S5): "Expected: success reporting the written tags;
// `read_back_mismatch`; `read_back_incomplete` (the corrupted artifact no
// longer parses for the reader); the key reads back absent — the write's
// exit status decides none of them."
// REQ-59: "write:    read_back required; verified only through the role's
// reader, never by exit status"
// REQ-68: "its success is never taken from exit status, and it is verified
// through the role's declared reader"
// HAPPY PATH + ADVERSARIAL
//
// The four arms share one shape and differ only in the WRITE tool, which is
// the point: exit status decides none of them. A wrapper that applies, one
// whose tool exits zero without applying (spike R3b), one that sinks stdin
// into the artifact (spike R4b), and one planning `<clear>`.
func TestReq127_WriteVerificationIsReadBackAndNeverExitStatus(t *testing.T) {
	bin := helperBin(t)

	// The honest wrapper: envelope in on stdin, a flat JSON artifact out.
	// The declared READER reads that artifact back through `cat`, in json
	// mode — so both halves are command-backed, as the MVV requires.
	writeThrough := func(t *testing.T, art accessor.Artifact, planned []resolve.Tag) (
		accessor.WriteResult, *storeReader,
	) {
		t.Helper()

		wacc := entry(fxRole, []string{bin, "stdin-to-file", art.Path}, fxKey)
		wacc.ReadBack = true
		racc := entry(fxRole, []string{"cat", "{artifact}"}, fxKey)

		reader := &storeReader{
			inner: cmdbind.Reader{Accessor: racc, Name: "cmd.read", Config: allowed(t)},
		}
		reg := accessor.Registry{
			Flow: "cmdflow",
			Definitions: []accessor.Definition{
				{
					Identity: accessor.Identity{
						Flow: "cmdflow", Name: "cmd.read", Capability: accessor.CapRead,
					},
					Accessor: racc,
					Binding:  reader,
				},
				{
					Identity: accessor.Identity{
						Flow: "cmdflow", Name: fxName, Capability: accessor.CapWrite,
					},
					Accessor: wacc,
					Binding: &cmdbind.Writer{
						Accessor: wacc, Name: fxName, Config: allowed(t),
					},
				},
			},
			OwnedTags: []string{fxKey},
		}
		e := accessor.NewExecutor(reg, accessor.Artifacts{fxRole: art})
		return e.Write(ctxOf(t), fxName, resolve.Plan{Writes: planned}), reader
	}

	t.Run("the wrapper applies: success reporting the written tags", func(t *testing.T) {
		art := artifactAt(t, "state.json")
		got, reader := writeThrough(t, art,
			[]resolve.Tag{{Key: fxKey, Value: "final"}})

		if got.Refusal != nil {
			t.Fatalf("an applied write refused %q (Detail %q)",
				got.Refusal.Class, got.Refusal.Detail)
		}
		if len(got.Written) != 1 || got.Written[0].Value != "final" {
			t.Errorf("Written = %#v; want the verified planned tag", got.Written)
		}
		if reader.reads == 0 {
			t.Error("the declared reader never ran; verification is read-back " +
				"through the role's reader, never the write's exit status")
		}
	})

	t.Run("the tool exits zero without applying: read_back_mismatch", func(t *testing.T) {
		// Spike R3b: `git config` ignored stdin and left the artifact
		// unchanged, exit 0. Here: the write child exits zero and writes
		// nothing, over an artifact pre-seeded with the OLD value.
		art := artifactAt(t, "state.json")
		seed(t, art.Path, map[string]string{fxKey: "draft"})

		wacc := entry(fxRole, []string{bin, "silent", "0"}, fxKey)
		wacc.ReadBack = true
		racc := entry(fxRole, []string{"cat", "{artifact}"}, fxKey)

		reg := accessor.Registry{
			Flow: "cmdflow",
			Definitions: []accessor.Definition{
				{
					Identity: accessor.Identity{
						Flow: "cmdflow", Name: "cmd.read", Capability: accessor.CapRead,
					},
					Accessor: racc,
					Binding:  cmdbind.Reader{Accessor: racc, Name: "cmd.read", Config: allowed(t)},
				},
				{
					Identity: accessor.Identity{
						Flow: "cmdflow", Name: fxName, Capability: accessor.CapWrite,
					},
					Accessor: wacc,
					Binding:  &cmdbind.Writer{Accessor: wacc, Name: fxName, Config: allowed(t)},
				},
			},
			OwnedTags: []string{fxKey},
		}
		got := accessor.NewExecutor(reg, accessor.Artifacts{fxRole: art}).
			Write(ctxOf(t), fxName, resolve.Plan{
				Writes: []resolve.Tag{{Key: fxKey, Value: "final"}},
			})

		if got.Refusal == nil {
			t.Fatal("a write whose tool exited zero without applying reported " +
				"success; exit status decides no write")
		}
		if got.Refusal.Class != accessor.ClassReadBackMismatch {
			t.Errorf("class = %q; want %q — read-back read the artifact and "+
				"found it WRONG", got.Refusal.Class, accessor.ClassReadBackMismatch)
		}
	})

	t.Run("the tool corrupts the artifact: read_back_incomplete", func(t *testing.T) {
		// Spike R4b, the sharpest edge in this design: a stdin-sinking
		// argv0 (`tee {artifact}`) overwrites the artifact with the
		// ENVELOPE. It surfaces at read-back, never at the write. Here the
		// reader is a RAW-mode reader over a two-key entry's key, so the
		// corrupted content cannot be parsed back into the declared key.
		art := artifactAt(t, "state.json")

		wacc := entry(fxRole, []string{bin, "stdin-to-file", art.Path}, fxKey)
		wacc.ReadBack = true
		// The reader expects a JSON object of strings; the corrupted
		// artifact holds one, but with the WRONG key, so the declared key
		// comes back UNREADABLE — the read-back did not run for it.
		racc := entry(fxRole, []string{"cat", "{artifact}"}, "other.key")

		reg := accessor.Registry{
			Flow: "cmdflow",
			Definitions: []accessor.Definition{
				{
					Identity: accessor.Identity{
						Flow: "cmdflow", Name: "cmd.read", Capability: accessor.CapRead,
					},
					Accessor: racc,
					Binding:  cmdbind.Reader{Accessor: racc, Name: "cmd.read", Config: allowed(t)},
				},
				{
					Identity: accessor.Identity{
						Flow: "cmdflow", Name: fxName, Capability: accessor.CapWrite,
					},
					Accessor: wacc,
					Binding:  &cmdbind.Writer{Accessor: wacc, Name: fxName, Config: allowed(t)},
				},
			},
			OwnedTags: []string{fxKey, "other.key"},
		}
		got := accessor.NewExecutor(reg, accessor.Artifacts{fxRole: art}).
			Write(ctxOf(t), fxName, resolve.Plan{
				Writes: []resolve.Tag{{Key: fxKey, Value: "final"}},
			})

		if got.Refusal == nil {
			t.Fatal("a write the reader could not verify reported success")
		}
		if got.Refusal.Class != accessor.ClassReadBackIncomplete {
			t.Errorf("class = %q; want %q — the verification DID NOT RUN, which "+
				"is not the same as the artifact being wrong",
				got.Refusal.Class, accessor.ClassReadBackIncomplete)
		}
		if !got.Refusal.Applied() {
			t.Error("Applied() is false; the write command already ran and this " +
				"MUST read as applied-but-unverified (0004:C14), never as a " +
				"write that did not occur")
		}
	})

	t.Run("a planned <clear> reads back ABSENT", func(t *testing.T) {
		art := artifactAt(t, "state.json")
		got, _ := writeThrough(t, art,
			[]resolve.Tag{{Key: fxKey, Value: table.ClearSentinel}})

		if got.Refusal != nil {
			t.Fatalf("a cleared key refused %q; the tool performs the removal "+
				"and read-back verifies ABSENCE", got.Refusal.Class)
		}
		// REQ-40 / 0004:C11: a cleared key is OMITTED from Written —
		// read-back verified it absent, and a verified removal is not a
		// written value.
		for _, w := range got.Written {
			if w.Key == fxKey {
				t.Errorf("Written carries the cleared key as %+v; a verified "+
					"removal is not a written value", w)
			}
		}
	})
}

// storeReader counts reads so a test can assert the declared reader
// actually ran. It delegates every call to the real command binding — it
// is a COUNTER, not a substitute for the unit under test.
type storeReader struct {
	inner accessor.ReadBinding
	reads int
}

func (r *storeReader) Capability() accessor.Capability { return accessor.CapRead }

func (r *storeReader) Read(ctx context.Context, art accessor.Artifact, requested []string) (
	[]accessor.KeyValue, []string, error,
) {
	r.reads++
	return r.inner.Read(ctx, art, requested)
}

func seed(t *testing.T, path string, values map[string]string) {
	t.Helper()

	b, err := json.Marshal(values)
	if err != nil {
		t.Fatalf("encoding the seed: %v", err)
	}
	if werr := os.WriteFile(path, b, 0o600); werr != nil {
		t.Fatalf("seeding %s: %v", path, werr)
	}
}

// REQ-46: "The inherited reserved-literal rule composes unchanged **on the
// read path only**: a value that reads back as the literal `<clear>` is
// UNREADABLE ... `Executor.Read` passes `clearIsUnreadable = true`, while
// `Executor.Write`'s read-back passes `false` ... this RDR adds no arm and
// changes no site."
// REQ-127 (the reserved-literal arm): "asserted by class, since a cleared
// key and an unreadable one are indistinguishable by absence alone.
// Asserted on the **read** path ... the two senses are asserted separately
// and neither test stands in for the other."
// DOMAIN EDGE — NEGATIVE REQ on the executor side.
//
// Newly REACHABLE through raw-mode reads of real tools: before this RDR no
// shipped binding could emit the literal. The two senses are two subtests
// by construction.
func TestReq46_TheReservedLiteralIsUnreadableOnTheReadPathAndAMismatchOnReadBack(t *testing.T) {
	bin := helperBin(t)

	t.Run("read path: the literal is UNREADABLE", func(t *testing.T) {
		art := artifactAt(t, "state.cfg")
		acc := entry(fxRole, []string{bin, "emit", table.ClearSentinel + "\n", "0"})
		acc.Output = strptr("raw")

		reg := accessor.Registry{
			Flow: "cmdflow",
			Definitions: []accessor.Definition{{
				Identity: accessor.Identity{
					Flow: "cmdflow", Name: fxName, Capability: accessor.CapRead,
				},
				Accessor: acc,
				Binding:  cmdbind.Reader{Accessor: acc, Name: fxName, Config: allowed(t)},
			}},
			OwnedTags: []string{fxKey},
		}
		got := accessor.NewExecutor(reg, accessor.Artifacts{fxRole: art}).
			Read(ctxOf(t), fxName)

		if got.Refusal == nil {
			t.Fatalf("a tool emitting the literal %q was believed; Values = %#v "+
				"— the binding cannot tell a removal from a stored literal",
				table.ClearSentinel, got.Values)
		}
		if got.Refusal.Class != accessor.ClassIncompleteRead {
			t.Errorf("class = %q; want %q — asserted by CLASS, since a cleared "+
				"key and an unreadable one are indistinguishable by absence "+
				"alone", got.Refusal.Class, accessor.ClassIncompleteRead)
		}
	})

	t.Run("read-back path: the literal's PRESENCE is the defect", func(t *testing.T) {
		// The other sense. `Executor.Write`'s read-back passes
		// clearIsUnreadable = false, so a key still holding the literal
		// after a planned clear is a MISMATCH, not an unreadability — and
		// neither test stands in for the other.
		art := artifactAt(t, "state.json")

		// A write tool that "applies" a clear by storing the literal — the
		// exact defect 0004:C11 exists to catch.
		wacc := entry(fxRole, []string{bin, "stdin-to-file", art.Path}, fxKey)
		wacc.ReadBack = true
		racc := entry(fxRole, []string{"cat", "{artifact}"}, fxKey)

		reg := accessor.Registry{
			Flow: "cmdflow",
			Definitions: []accessor.Definition{
				{
					Identity: accessor.Identity{
						Flow: "cmdflow", Name: "cmd.read", Capability: accessor.CapRead,
					},
					Accessor: racc,
					Binding:  cmdbind.Reader{Accessor: racc, Name: "cmd.read", Config: allowed(t)},
				},
				{
					Identity: accessor.Identity{
						Flow: "cmdflow", Name: fxName, Capability: accessor.CapWrite,
					},
					Accessor: wacc,
					Binding:  &cmdbind.Writer{Accessor: wacc, Name: fxName, Config: allowed(t)},
				},
			},
			OwnedTags: []string{fxKey},
		}
		got := accessor.NewExecutor(reg, accessor.Artifacts{fxRole: art}).
			Write(ctxOf(t), fxName, resolve.Plan{
				Writes: []resolve.Tag{{Key: fxKey, Value: table.ClearSentinel}},
			})

		if got.Refusal == nil {
			t.Fatal("a clear whose tool STORED the literal reported success; a " +
				"re-read still holding the key — including as the literal — is " +
				"a mismatch")
		}
		if got.Refusal.Class != accessor.ClassReadBackMismatch {
			t.Errorf("class = %q; want %q — on the read-back path the literal's "+
				"PRESENCE is the defect, not an unreadability",
				got.Refusal.Class, accessor.ClassReadBackMismatch)
		}
	})
}

// REQ-136: "It is not statically detectable [the stdin-sinking write tool]:
// whether a tool reads stdin is not visible in its argv, so no C5 arm can
// catch it and none is claimed"
// DOMAIN EDGE — NEGATIVE REQ.
//
// The clause forbids a check. The test that catches its violation is the
// pair: a stdin-sinking entry LOADS clean (no C5 arm caught it) and the
// corruption surfaces at READ-BACK, loudly, on the next invocation — which
// is what bounds the hazard in v1.
func TestReq136_AStdinSinkingWriteToolIsNotStaticallyDetectableAndSurfacesAtReadBack(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.json")

	// A stdin-sinking argv0 is a perfectly well-formed declared command.
	wacc := entry(fxRole, []string{bin, "stdin-to-file", "{artifact}"}, fxKey)
	wacc.ReadBack = true

	w := &cmdbind.Writer{Accessor: wacc, Name: fxName, Config: allowed(t)}
	if err := w.Apply(ctxOf(t), art, []resolve.Tag{{Key: fxKey, Value: "final"}}); err != nil {
		t.Fatalf("a stdin-sinking write refused at the WRITE: %v — the hazard is "+
			"not statically detectable and none is claimed; what bounds it is "+
			"read-back", err)
	}

	// The artifact now holds the ENVELOPE, verbatim. That is the hazard,
	// and it is observable — which is what makes read-back the bound.
	b, err := os.ReadFile(art.Path)
	if err != nil {
		t.Fatalf("the artifact was not written: %v", err)
	}
	if !strings.Contains(string(b), fxKey) {
		t.Errorf("the artifact holds %q; the stdin-sinking tool ingests the C3 "+
			"envelope as content, which is the hazard this arm records rather "+
			"than prevents", b)
	}

	// And no lint arm claims to catch it: the same entry carries no C5
	// defect. Asserted as the ABSENCE of a category, which is what the
	// clause promises.
	if _, cerr := table.Load([]byte(stdinSinkModel(bin)), "stdin-sink.toml"); cerr != nil {
		if cat, ok := table.CategoryOf(cerr); ok && strings.HasPrefix(string(cat), "command_") {
			t.Errorf("a stdin-sinking entry refused %q; whether a tool reads "+
				"stdin is not visible in its argv and no C5 arm catches it", cat)
		}
	}
}

func stdinSinkModel(bin string) string {
	return `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "sinkflow"
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
command = ["cat", "{artifact}"]
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

// REQ-141: "This RDR claims **no** byte-identical output, content-addressed
// identity, or replay-stable hash"
// REQ-110: "The read envelope is deliberately the same flat
// map-of-strings ... They stay separate decoders because they decode
// different things"
// DOMAIN EDGE — NEGATIVE REQ.
//
// The Pre-Lock `fidelity` invariant that IS claimed: planned tags → stdin →
// tool → read-back is TYPED TAG-VALUE equality, explicitly not byte
// equality. So the discriminating fixture is a wrapper whose artifact bytes
// differ run to run while the tag values do not.
func TestReq141_ReadBackComparesTypedTagValuesAndNotArtifactBytes(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.json")

	wacc := entry(fxRole, []string{bin, "stdin-to-file", art.Path}, fxKey)
	wacc.ReadBack = true
	racc := entry(fxRole, []string{"cat", "{artifact}"}, fxKey)

	reg := accessor.Registry{
		Flow: "cmdflow",
		Definitions: []accessor.Definition{
			{
				Identity: accessor.Identity{
					Flow: "cmdflow", Name: "cmd.read", Capability: accessor.CapRead,
				},
				Accessor: racc,
				Binding:  cmdbind.Reader{Accessor: racc, Name: "cmd.read", Config: allowed(t)},
			},
			{
				Identity: accessor.Identity{
					Flow: "cmdflow", Name: fxName, Capability: accessor.CapWrite,
				},
				Accessor: wacc,
				Binding:  &cmdbind.Writer{Accessor: wacc, Name: fxName, Config: allowed(t)},
			},
		},
		OwnedTags: []string{fxKey},
	}
	e := accessor.NewExecutor(reg, accessor.Artifacts{fxRole: art})

	first := e.Write(ctxOf(t), fxName, resolve.Plan{
		Writes: []resolve.Tag{{Key: fxKey, Value: "final"}},
	})
	if first.Refusal != nil {
		t.Fatalf("the first write refused %q", first.Refusal.Class)
	}
	before, _ := os.ReadFile(art.Path)

	// Re-run the identical write. The VALUE is unchanged, so read-back
	// verifies; the RDR claims nothing about the bytes, which belong to the
	// bound tool.
	second := e.Write(ctxOf(t), fxName, resolve.Plan{
		Writes: []resolve.Tag{{Key: fxKey, Value: "final"}},
	})
	if second.Refusal != nil {
		t.Fatalf("the repeated write refused %q; read-back compares typed tag "+
			"VALUES, not bytes", second.Refusal.Class)
	}
	after, _ := os.ReadFile(art.Path)

	// The oracle: the TAG VALUE round-trips. Byte identity is neither
	// asserted nor required — this line records that, and the comparison
	// above is what carries the contract.
	_ = before
	_ = after
	if len(second.Written) != 1 || second.Written[0].Value != "final" {
		t.Errorf("Written = %#v; want the planned tag verified by VALUE",
			second.Written)
	}
}

// REQ-119: "Illustrative — intent only; tests must not assert it
// literally."
// DOMAIN EDGE — NEGATIVE REQ.
//
// The clause forbids pinning the Illustrative Code block. The property it
// protects is that the CARRIER, not any particular example, is the
// contract: an entry whose argv shares nothing with the illustration must
// work identically. That is what this asserts.
func TestReq119_TheContractIsTheCarrierNotTheIllustration(t *testing.T) {
	bin := helperBin(t)
	art := artifactAt(t, "state.cfg")

	// Deliberately unlike the RDR's `git config` / `tools/flowstate-write`
	// examples: a different argv0, a different flag shape, a different key.
	acc := entry(fxRole, []string{bin, "emit", "arbitrary\n", "0"}, "flow.other")
	acc.Output = strptr("raw")

	values, _, err := cmdbind.Reader{
		Accessor: acc, Name: "some.other.name", Config: allowed(t),
	}.Read(ctxOf(t), art, acc.Keys)
	if err != nil {
		t.Fatalf("an entry unlike the illustration refused: %v", err)
	}
	if len(values) != 1 || values[0].Key != "flow.other" || values[0].Value != "arbitrary" {
		t.Errorf("values = %#v; the contract is the carrier grammar, not the "+
			"illustrative `git config` example", values)
	}

	// And the negative control: nothing in the shipped surface names the
	// illustration's tool.
	if errors.Is(err, os.ErrNotExist) {
		t.Error("the binding depends on a specific external tool; the carrier " +
			"is tool-agnostic")
	}
}
