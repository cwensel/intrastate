package cli

// RDR 0010 §F–§L — cross-cutting compatibility (REQ-70..REQ-73), the
// failure modes in assertable form (REQ-74..REQ-76), and the naming /
// rejected-shape negative REQs (REQ-107..REQ-110).
//
// These are the clauses whose oracle is the CLI surface rather than a
// package API, plus the negative REQs whose only assertable form is
// "the rejected shape is absent".

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
	"github.com/newcoinc/intrastate/internal/table"
)

// REQ-70: "`class` is optional and its absence reads as `\"state-machine\"`
// (C1), so every checked-in model loads, lints, and resolves unchanged; the
// agreement check is one-directional, so no existing model can be newly
// refused."
// DOMAIN EDGE — negative REQ.
//
// GREEN BY CONSTRUCTION against `main` for the "loads and lints" half; the
// discriminating half is that it stays green after the class check lands.
func TestReq70_EveryCheckedInModelLoadsLintsAndResolvesUnchanged(t *testing.T) {
	path := filepath.Join(repoRootFor(t), checkedInModelPath)

	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", checkedInModelPath, err)
	}

	m, lerr := table.Load(src, checkedInModelPath)
	if lerr != nil {
		t.Fatalf("the checked-in model no longer loads: %v", lerr)
	}
	if m.Class != "" && m.Class != "state-machine" {
		t.Errorf("Class = %q; the checked-in model declares no `class`, and "+
			"an absent `class` reads as `state-machine`", m.Class)
	}

	stdout, _, cerr := runCmd(t, "lint", "--model", path, "--as=json")
	if cerr != nil {
		t.Errorf("the checked-in model no longer lints clean: %v\n%s",
			cerr, stdout)
	}
}

// REQ-73: "`[model] version` stays the model's own version, per `0002:C3`"
// and "The model grammar carries no version negotiation and this RDR adds
// none"
// DOMAIN EDGE — negative REQ.
//
// GREEN BY CONSTRUCTION: nothing today negotiates a version, and this test
// pins that `class` does not become a second version axis.
func TestReq73_TheClassIsNotAVersionAndNoNegotiationIsAdded(t *testing.T) {
	src := strings.Replace(dtModel0010, "version = 1", "version = 2", 1)

	_, err := table.Load([]byte(src), "dt-v2.toml")
	if err == nil {
		t.Fatal("a decision table declaring `version = 2` loaded clean; " +
			"`[model] version` stays the model's own version and this RDR " +
			"adds no negotiation")
	}
	if cat, ok := table.CategoryOf(err); !ok ||
		cat != table.CatUnsupportedVersion {
		t.Errorf("category = %q; want %q — the class is not a version axis",
			cat, table.CatUnsupportedVersion)
	}
}

// REQ-75 / FM (Silent, guarded): "`intrastate lint --as=json` over a
// deliberately partial table must report a coverage finding, and over a
// match-discriminated table must report `graph-unprovable-coverage` (C5) —
// neither may exit 0 with `[]`."
// ADVERSARIAL — the two silent failure modes, at the CLI surface.
func TestReq75_NeitherGuardedSilentModeExitsZeroWithAnEmptyFindingList(t *testing.T) {
	cases := map[string]struct {
		src  string
		code string
	}{
		"partial table":             {mvvPartial0010, "graph-coverage-gap"},
		"match-discriminated table": {dtMatchOnlyModel0010, "graph-unprovable-coverage"},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			path := writeModel(t, c.src)

			stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
			if err == nil {
				t.Fatalf("lint exited 0 over the %s; that is the silent "+
					"failure mode this clause guards:\n%s", name, stdout)
			}
			if got := clierr.ExitCodeFor(err); got != 2 {
				t.Errorf("exit code = %d; want 2", got)
			}

			env := parseFailure(t, stdout)
			if env.Findings == nil || len(*env.Findings) == 0 {
				t.Fatalf("the failure carries an EMPTY finding list; neither "+
					"mode may exit with `[]`:\n%s", stdout)
			}
			var codes []string
			for _, f := range *env.Findings {
				codes = append(codes, f.Code)
			}
			if !slices.Contains(codes, c.code) {
				t.Errorf("findings = %v; want one carrying %q", codes, c.code)
			}
		})
	}
}

// REQ-76 / FM (Recovery): "drop `class` (the model still loads and lint
// reports it as a rootless machine — `graph-dangling-edge`, exactly as
// today)" and "nothing persistent is written by this RDR."
// DOMAIN EDGE
func TestReq76_DroppingTheClassRecoversAndNothingPersistentIsWritten(t *testing.T) {
	t.Run("dropping the class recovers to today's behaviour", func(t *testing.T) {
		path := writeModel(t, mvvClassOmitted0010)

		stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
		if err == nil {
			t.Fatalf("the class-omitted model linted clean:\n%s", stdout)
		}
		env := parseFailure(t, stdout)
		if env.Findings == nil {
			t.Fatalf("no findings:\n%s", stdout)
		}
		var rooted bool
		for _, f := range *env.Findings {
			if f.Code == "graph-dangling-edge" && f.Element == "model" {
				rooted = true
			}
		}
		if !rooted {
			t.Errorf("no rootless-machine finding; recovery is `drop the "+
				"class` and get exactly today's behaviour:\n%s", stdout)
		}
	})

	// The oracle is a CONTENT-AND-IDENTITY snapshot, not an entry count: an
	// overwrite of `dt.toml` in place, a rename, or any one-for-one
	// replacement leaves the count unchanged and would pass a `len(after)
	// != len(before)` compare while having written something persistent.
	t.Run("resolve writes nothing persistent", func(t *testing.T) {
		dir := t.TempDir()
		model := filepath.Join(dir, "dt.toml")
		if err := os.WriteFile(model, []byte(dtModel0010), 0o600); err != nil {
			t.Fatalf("write model: %v", err)
		}

		before := snapshotDir(t, dir)

		requireSuccess(t, "flow", "resolve", "--model", model,
			"--outcome", "decide", "--tag", "a=x", "--tag", "b=p", "--as=json")

		after := snapshotDir(t, dir)
		if !reflect.DeepEqual(after, before) {
			t.Errorf("`flow resolve` over a decision table changed the "+
				"directory; nothing persistent is written by this RDR\n"+
				"before: %#v\nafter:  %#v", before, after)
		}

		// The model's BYTES specifically: an in-place rewrite of identical
		// length and mode would survive the walk above on a filesystem with
		// coarse mtimes.
		got, err := os.ReadFile(model)
		if err != nil {
			t.Fatalf("re-read model: %v", err)
		}
		if string(got) != dtModel0010 {
			t.Errorf("`dt.toml` was rewritten by `flow resolve`; the model " +
				"file is read-only input")
		}
	})
}

// REQ-107 / `0010:D-naming`: "**Naming** — `class` with values
// `state-machine` / `decision-table`. Rejected: `kind` …, `type` …,
// `stateless = true` …. `emit` for the answer block. Rejected: `output` …,
// `result`, `answer`."
// ADVERSARIAL — negative REQ: the rejected spellings are unauthorable.
func TestReq107_TheRejectedNamingSpellingsAreUnauthorable(t *testing.T) {
	t.Run("the model-key spellings", func(t *testing.T) {
		for _, rejected := range []string{
			`kind = "decision-table"`,
			`type = "decision-table"`,
			`stateless = true`,
		} {
			t.Run(rejected, func(t *testing.T) {
				src := strings.Replace(dtModel0010,
					`class = "decision-table"`, rejected, 1)
				_, err := table.Load([]byte(src), "rejected.toml")
				if err == nil {
					t.Fatalf("`[model] %s` loaded clean; the closed `[model]` "+
						"layout admits `class` and no rejected alias", rejected)
				}
				if cat, _ := table.CategoryOf(err); cat != table.CatUnknownSchemaField {
					t.Errorf("category = %q; want %q — a rejected alias is an "+
						"unknown schema field", cat, table.CatUnknownSchemaField)
				}
			})
		}
	})

	t.Run("the rule-block spellings", func(t *testing.T) {
		for _, rejected := range []string{"output", "result", "answer"} {
			t.Run(rejected, func(t *testing.T) {
				src := strings.Replace(dtModel0010,
					"[rule.emit]\nverdict = \"alpha\"\ncode = \"1\"\n",
					"[rule."+rejected+"]\nverdict = \"alpha\"\n", 1)
				if src == dtModel0010 {
					t.Fatal("the substitution did not apply")
				}
				_, err := table.Load([]byte(src), "rejected.toml")
				if err == nil {
					t.Fatalf("`[rule.%s]` loaded clean; `emit` is the answer "+
						"block and the aliases are rejected", rejected)
				}
				if cat, _ := table.CategoryOf(err); cat != table.CatUnknownSchemaField {
					t.Errorf("category = %q; want %q",
						cat, table.CatUnknownSchemaField)
				}
			})
		}
	})
}

// REQ-108 / `0010:D-selection-predicate` / A5: "**Selection / predicate** —
// unchanged: the kernel's exact-one match is the decision table's hit
// policy (A5); two matching rules is `flow-ambiguous-match` in both
// classes, and an escape row rescues it the same way."
// ADVERSARIAL — negative REQ: no first-hit, priority, or collect policy.
func TestReq108_TwoMatchingRulesIsAmbiguousMatchInTheDecisionTableClass(t *testing.T) {
	// Two rules accepting the SAME (a,b) cell. Under a first-hit or
	// priority policy the first would win; the kernel's exact-one match
	// refuses.
	src := dtModel0010 + `
[[rule]]
id = "cell-xp-again"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "x"
[rule.guard.all.b]
eq = "p"
[rule.emit]
verdict = "shadow"
`
	model := writeFlowModel(t, src)

	requireRefusal(t, "flow-ambiguous-match", 2,
		"flow", "resolve", "--model", model, "--outcome", "decide",
		"--tag", "a=x", "--tag", "b=p", "--as=json")
}

// REQ-109 / `0010:BR1`..`0010:BR5`: the briefly-rejected shapes.
// ADVERSARIAL — negative REQs, each in its assertable form.
func TestReq109_TheBrieflyRejectedShapesAreAbsent(t *testing.T) {
	model := writeFlowModel(t, dtModel0010)

	t.Run("BR1 — the answer is not row identity ALONE", func(t *testing.T) {
		data := flowData(t, requireSuccess(t, "flow", "resolve",
			"--model", model, "--outcome", "decide",
			"--tag", "a=x", "--tag", "b=p", "--as=json"))
		if _, ok := data["emit"]; !ok {
			t.Error("the payload carries only `rule`; `0010:BR1` (answer = " +
				"row identity only) was rejected in favour of `rule` + `emit`")
		}
	})

	t.Run("BR2 — the answer is not a pseudo-owned `next` tag", func(t *testing.T) {
		data := flowData(t, requireSuccess(t, "flow", "resolve",
			"--model", model, "--outcome", "decide",
			"--tag", "a=x", "--tag", "b=p", "--as=json"))
		next, ok := data["next"].(map[string]any)
		if !ok {
			t.Fatalf("`next` = %#v", data["next"])
		}
		if len(next) != 0 {
			t.Errorf("`next` = %#v; the answer does not ride a pseudo-owned "+
				"next tag (`0010:BR2`, the PoC's shape, rejected)", next)
		}
	})

	t.Run("BR3 — emit values are strings, never `any`", func(t *testing.T) {
		src := strings.Replace(dtModel0010,
			"[rule.emit]\nverdict = \"alpha\"\ncode = \"1\"\n",
			"[rule.emit]\nverdict = 1\n", 1)
		if _, err := table.Load([]byte(src), "any-emit.toml"); err == nil {
			t.Error("an `any`-typed emit value loaded clean; `0010:BR3` was " +
				"rejected and the values MUST be strings")
		}
	})

	t.Run("BR5 — --outcome is not optional for a one-outcome model", func(t *testing.T) {
		if _, _, err := runCmd(t, "flow", "resolve",
			"--model", model, "--tag", "a=x", "--tag", "b=p", "--as=json"); err == nil {
			t.Error("`--outcome` was optional over a model declaring one " +
				"outcome; `0010:BR5` was rejected")
		}
	})

	t.Run("BR6 — no lint finding announces the class", func(t *testing.T) {
		path := writeModel(t, mvvComplete0010)
		stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
		if err != nil {
			t.Fatalf("the complete decision table did not lint clean: %v\n%s",
				err, stdout)
		}
		if strings.Contains(stdout, "decision-table") {
			t.Errorf("the lint envelope names the class; `0010:BR6` (a new "+
				"lint finding announcing the class) was rejected:\n%s", stdout)
		}
	})
}

// dtMatchOnlyModel0010 is the match-discriminated table REQ-75's second arm
// needs: the same shape discriminating by `[rule.match.<key>]`, so the
// scoped product has zero participating dimensions.
const dtMatchOnlyModel0010 = `outcomes = ["decide"]

[model]
id = "dt"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.a]
provenance = "observed"
kind = "enum"
domain = ["x", "y"]
single_valued = true
required = true

[[rule]]
id = "cell-x"
[rule.match.recognized]
eq = "decide"
[rule.match.a]
eq = "x"
[rule.emit]
verdict = "x"
`

// dirEntrySnapshot is one filesystem entry's identity: the fields a
// persistent write would have to disturb. Name alone is not enough — an
// in-place overwrite keeps it — so size, mode, and mtime ride along.
type dirEntrySnapshot struct {
	Size  int64
	Mode  fs.FileMode
	MTime int64
	IsDir bool
}

// snapshotDir walks dir and records every entry's relative path and
// identity, so "nothing persistent is written" can be asserted as an
// equality on the whole tree rather than as an entry COUNT.
func snapshotDir(t *testing.T, dir string) map[string]dirEntrySnapshot {
	t.Helper()

	out := map[string]dirEntrySnapshot{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		out[rel] = dirEntrySnapshot{
			Size:  info.Size(),
			Mode:  info.Mode(),
			MTime: info.ModTime().UnixNano(),
			IsDir: d.IsDir(),
		}
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot %s: %v", dir, err)
	}
	return out
}
