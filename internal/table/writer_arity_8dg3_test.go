package table_test

// Kata 8dg3 — the writer-arity guard must cover EVERY declared owned tag,
// not only the keys some rule writes or `[initial]` assigns.
//
// `checkAccessorBindings` originally counted writers over a narrowed set:
// `written` = the union of every row's `RequiresOwned` and `[initial]`.
// An owned key declared in `[tags]`, served by TWO `[write.<id>]` entries,
// and named by no rule write, no clear list, and no `[initial]` fell
// outside that set and LOADED CLEAN — verified empirically before the fix.
// The reader guard already iterates `model.Tags` and has no such hole,
// which is what made the missing symmetry look sound.
//
// These fixtures are authored in-test by mutating the promoted RDR
// fixture, the same shape `TestReq24_EveryWrittenKeyIsServedByExactlyOneWriter`
// uses. A new file under `testdata/neg/` would be wrong here: that
// directory is the REQ-118 promoted spike set, pinned per fixture to the
// category it was promoted witnessing.

import (
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/table"
)

// secondClusterWriter is a second `[write.<id>]` naming `cluster_ready` —
// an owned tag in the RDR fixture that NO rule writes or clears and that
// `[initial]` does not assign. Appending it to the fixture produces the
// exact model class the narrowed guard could not see.
const secondClusterWriter = `
[write.rdr-cluster]
role = "rdr"
path = "rdr.cluster"
keys = ["cluster_ready"]
timeout = "2s"
read_back = true
`

// TestKata8dg3_OwnedKeyServedByTwoWritersRefusesEvenWhenNoRuleNamesIt is
// the RED test: before the fix this model loaded clean.
func TestKata8dg3_OwnedKeyServedByTwoWritersRefusesEvenWhenNoRuleNamesIt(t *testing.T) {
	base := string(readFixture(t, rdrFixture))

	// Guard the premise: if a future fixture edit makes some rule write or
	// clear `cluster_ready`, or `[initial]` assign it, this model would be
	// caught by the pre-existing `written` guard and the test would stop
	// discriminating the widened one.
	pre := mustLoad(t, rdrFixture)
	for _, row := range pre.Rows {
		if strings.Contains(strings.Join(row.RequiresOwned, ","), "cluster_ready") {
			t.Fatalf("row %s requires `cluster_ready`; the fixture no longer "+
				"witnesses the key class this test exists for — it would be "+
				"caught by the narrowed `written` guard instead", row.Identity())
		}
	}
	for _, tv := range pre.Initial {
		if tv.Key == "cluster_ready" {
			t.Fatal("`[initial]` assigns `cluster_ready`; the fixture no longer " +
				"witnesses the key class this test exists for")
		}
	}
	if decl, ok := pre.Tags["cluster_ready"]; !ok || decl.Provenance != table.ProvenanceOwned {
		t.Fatalf("`cluster_ready` provenance = %v (declared: %v); want owned",
			decl.Provenance, ok)
	}

	src := base + secondClusterWriter
	_, err := table.Load([]byte(src), "two-writers-unwritten-key.toml")
	if err == nil {
		t.Fatal("an owned key served by TWO writers loaded clean because no " +
			"rule names it; the writer guard must cover every declared owned " +
			"tag, as the reader guard already does")
	}
	if cat, _ := table.CategoryOf(err); cat != table.CatMalformedAccessorBinding {
		t.Errorf("category = %q; want %q", cat, table.CatMalformedAccessorBinding)
	}
}

// TestKata8dg3_OwnedKeyServedByZeroWritersStillLoads pins the allowance the
// widened guard must NOT swallow.
//
// This is the trap the fix had to avoid: tightening "owned key, wrong
// writer count" into `writerCount != 1` for every owned tag would refuse
// the RDR 0005 MVV fixture, whose `note` key is owned, read by one reader,
// and served by no writer at all. `internal/cli::writerFor`'s doc comment
// names that exact case as the reason it checks writer keys rather than
// ownership.
func TestKata8dg3_OwnedKeyServedByZeroWritersStillLoads(t *testing.T) {
	base := string(readFixture(t, rdrFixture))

	// Drop `cluster_ready` from the sole writer's key list. It stays owned,
	// stays served by one reader, and is still named by no rule write, no
	// clear list, and no `[initial]`.
	src := strings.Replace(base,
		`keys = ["status", "stage", "profile", "iter", "rewind_scope", "prelock_lens", "cluster_ready"]
timeout = "2s"
read_back = true`,
		`keys = ["status", "stage", "profile", "iter", "rewind_scope", "prelock_lens"]
timeout = "2s"
read_back = true`, 1)
	if src == base {
		t.Fatal("writer keys substitution did not apply")
	}

	m, err := table.Load([]byte(src), "owned-zero-writers.toml")
	if err != nil {
		t.Fatalf("an owned key served by ZERO writers and named by no rule must "+
			"still load — the MVV fixture's `note` key depends on it; refused: %v",
			err)
	}
	for id, w := range m.Writers {
		for _, k := range w.Keys {
			if k == "cluster_ready" {
				t.Fatalf("writer %q still serves `cluster_ready`; the substitution "+
					"did not produce the zero-writer model this test asserts", id)
			}
		}
	}
}

// TestKata8dg3_WrittenKeyArityMessagesAreUnchanged pins that widening the
// guard did not reword or recategorize the two cases the narrowed guard
// already covered: a `written` key served by zero writers and by two.
func TestKata8dg3_WrittenKeyArityMessagesAreUnchanged(t *testing.T) {
	base := string(readFixture(t, rdrFixture))

	t.Run("written key with zero writers", func(t *testing.T) {
		// `stage` is written by three rules.
		src := strings.Replace(base,
			`[write.rdr-status]
role = "rdr"
path = "rdr.status"
keys = ["status", "stage", "profile", "iter", "rewind_scope", "prelock_lens", "cluster_ready"]`,
			`[write.rdr-status]
role = "rdr"
path = "rdr.status"
keys = ["status", "profile", "iter", "rewind_scope", "prelock_lens", "cluster_ready"]`, 1)
		if src == base {
			t.Fatal("writer keys substitution did not apply")
		}
		_, err := table.Load([]byte(src), "written-zero-writers.toml")
		if err == nil {
			t.Fatal("a written key served by zero writers loaded clean")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatMalformedAccessorBinding {
			t.Errorf("category = %q; want %q", cat, table.CatMalformedAccessorBinding)
		}
		if got := err.Error(); !strings.Contains(got, "written tag stage is served by 0 writers; want exactly one") {
			t.Errorf("message = %q; the `written`-key wording is unchanged by "+
				"the widened guard", got)
		}
	})

	t.Run("written key with two writers", func(t *testing.T) {
		// `stage` IS written by rules, so it is in `written` as well as in
		// `[tags]`. The `want exactly one` arm must win over the owned-tag
		// arm's `at most one`.
		src := base + `
[write.rdr-stage]
role = "rdr"
path = "rdr.stage"
keys = ["stage"]
timeout = "2s"
read_back = true
`
		_, err := table.Load([]byte(src), "written-two-writers.toml")
		if err == nil {
			t.Fatal("a written key served by two writers loaded clean")
		}
		if cat, _ := table.CategoryOf(err); cat != table.CatMalformedAccessorBinding {
			t.Errorf("category = %q; want %q", cat, table.CatMalformedAccessorBinding)
		}
		if got := err.Error(); !strings.Contains(got, "written tag stage is served by 2 writers; want exactly one") {
			t.Errorf("message = %q; a `written` key keeps the exactly-one "+
				"wording rather than falling to the owned-tag arm", got)
		}
	})
}
