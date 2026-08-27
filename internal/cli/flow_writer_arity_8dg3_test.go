package cli

// Kata 8dg3 — `writerFor` must require EXACTLY one `[write.<id>]` serving a
// key, not merely at least one.
//
// The original implementation returned on the FIRST writer naming the key,
// so a key served by two writers passed the guard and `groupByWriter` then
// routed the mutation to whichever writer sorted first. Read-back did not
// surface the wrong destination, because it re-reads through that same
// writer.
//
// The model here is built in memory rather than loaded from TOML on
// purpose: RDR 0002's loader now refuses a two-writers-one-key document, so
// a fixture-driven test would never reach `writerFor` and would assert the
// LOADER's guard a second time instead of the CLI's. Bypassing the loader
// is what keeps this test on the CLI guard's own terms — the guard must
// stand alone, because it is the only thing between a hand-assembled model
// and an arbitrary write destination.

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/table"
)

// twoWriterModel returns a model whose `status` key is served by TWO write
// accessors and whose `solo` key is served by exactly one. `orphan` is a
// declared owned tag served by NO writer — the `note` shape from the RDR
// 0005 MVV fixture, which must keep reporting the original zero-writer
// message.
func twoWriterModel() *table.Model {
	return &table.Model{
		ID:      "arity",
		Version: 1,
		Tags: map[string]table.TagDecl{
			"status": {Provenance: table.ProvenanceOwned, Kind: "scalar"},
			"solo":   {Provenance: table.ProvenanceOwned, Kind: "scalar"},
			"orphan": {Provenance: table.ProvenanceOwned, Kind: "scalar"},
		},
		Writers: map[string]table.Accessor{
			// `alpha` sorts first, so an implementation that short-circuits
			// on the first match silently routes `status` here.
			"alpha": {Role: "a", Path: "a.state", Keys: []string{"status"},
				Timeout: "2s", ReadBack: true},
			"beta": {Role: "b", Path: "b.state", Keys: []string{"status", "solo"},
				Timeout: "2s", ReadBack: true},
		},
	}
}

// setStateFlags returns the real `flow set-state` command with `--write` and
// `--clear` set to the given entries, so `parseWrites` is driven through the
// flag surface a caller actually uses.
func setStateFlags(t *testing.T, writes, clears []string) *cobra.Command {
	t.Helper()

	cmd := newFlowSetStateCmd()
	for _, w := range writes {
		if err := cmd.Flags().Set("write", w); err != nil {
			t.Fatalf("set --write %q: %v", w, err)
		}
	}
	for _, c := range clears {
		if err := cmd.Flags().Set("clear", c); err != nil {
			t.Fatalf("set --clear %q: %v", c, err)
		}
	}
	return cmd
}

// TestKata8dg3_KeyServedByTwoWritersIsUnbound is the RED test: before the
// fix both arms returned nil from `writerFor` and the request proceeded.
func TestKata8dg3_KeyServedByTwoWritersIsUnbound(t *testing.T) {
	m := twoWriterModel()

	t.Run("write", func(t *testing.T) {
		cmd := setStateFlags(t, []string{"status=final"}, nil)
		_, _, ce := parseWrites(cmd, m)
		if ce == nil {
			t.Fatal("a key served by TWO write accessors passed the guard; " +
				"REQ-61 requires exactly one, and `groupByWriter` would route " +
				"the mutation to the lexicographically first writer")
		}
		requireUnbound(t, ce, codeWriteUnbound, "status")
		if !strings.Contains(ce.Message, "2 write accessors") {
			t.Errorf("message = %q; the two-writer arm must say how many "+
				"accessors name the key, not reuse the zero-writer wording",
				ce.Message)
		}
	})

	t.Run("clear", func(t *testing.T) {
		cmd := setStateFlags(t, nil, []string{"status"})
		_, _, ce := parseWrites(cmd, m)
		if ce == nil {
			t.Fatal("a --clear key served by TWO write accessors passed the guard")
		}
		requireUnbound(t, ce, codeClearUnbound, "status")
		if !strings.Contains(ce.Message, "2 write accessors") {
			t.Errorf("message = %q; want the two-writer wording", ce.Message)
		}
	})
}

// TestKata8dg3_SingleAndZeroWriterBehaviourIsUnchanged pins the two arms the
// count must not disturb.
func TestKata8dg3_SingleAndZeroWriterBehaviourIsUnchanged(t *testing.T) {
	m := twoWriterModel()

	t.Run("exactly one writer passes", func(t *testing.T) {
		cmd := setStateFlags(t, []string{"solo=x"}, nil)
		planned, _, ce := parseWrites(cmd, m)
		if ce != nil {
			t.Fatalf("a key served by exactly one writer was refused %q: %s",
				ce.Code, ce.Message)
		}
		if len(planned) != 1 || planned[0].Key != "solo" {
			t.Fatalf("planned = %+v; want the single `solo` write", planned)
		}
	})

	t.Run("zero writers keeps the original message", func(t *testing.T) {
		// `orphan` is owned and served by no writer — the MVV `note` shape.
		// Its refusal text is the one REQ-86 published; the count must not
		// reword it.
		cmd := setStateFlags(t, []string{"orphan=x"}, nil)
		_, _, ce := parseWrites(cmd, m)
		if ce == nil {
			t.Fatal("a key served by zero writers passed the guard")
		}
		requireUnbound(t, ce, codeWriteUnbound, "orphan")
		want := "no declared write accessor's `keys` list names the tag `orphan`"
		if ce.Message != want {
			t.Errorf("message = %q; want %q — the zero-writer wording is "+
				"unchanged", ce.Message, want)
		}
	})

	t.Run("zero writers on --clear keeps the original message", func(t *testing.T) {
		cmd := setStateFlags(t, nil, []string{"orphan"})
		_, _, ce := parseWrites(cmd, m)
		if ce == nil {
			t.Fatal("a --clear key served by zero writers passed the guard")
		}
		requireUnbound(t, ce, codeClearUnbound, "orphan")
		want := "no declared write accessor's `keys` list names the tag `orphan`"
		if ce.Message != want {
			t.Errorf("message = %q; want %q", ce.Message, want)
		}
	})
}

// requireUnbound asserts the refusal carries the published code, the key as
// its param, and the group that fixes exit 2 (REQ-86, REQ-87). The code
// spellings are normative and RDR 0005's table is closed, so a fix that
// minted a new code fails here.
func requireUnbound(t *testing.T, ce *clierr.CLIError, code, param string) {
	t.Helper()

	if ce.Code != code {
		t.Errorf("code = %q; want %q — RDR 0005's refusal table is closed at "+
			"25 rows and the two-writer case reuses the published code",
			ce.Code, code)
	}
	if ce.Param != param {
		t.Errorf("param = %q; want %q", ce.Param, param)
	}
	if got := clierr.ExitCodeFor(ce); got != 2 {
		t.Errorf("exit code = %d; want 2", got)
	}
}
