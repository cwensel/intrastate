package accessor_test

// RDR 0028 — the Minimum Viable Validation, all EIGHT steps end to end.
//
// req-list ASSUMPTION A-7: the MVV has EIGHT steps. The Implementation
// Plan enumerates 1–8 and marks step 8 "required rather than optional";
// Testing Strategy's Done line says "seven". The stale count is the Done
// line's — a step the record argues is load-bearing ("Without this step
// the MVV proves the composed scenario only against a fixture curated to
// exclude the one class of drift the real corpus has") cannot be the one
// the count drops. All eight run here.
//
// WHAT THIS TEST IS NOT. A green exit code and "did not error" are NOT
// the acceptance oracle. The record declares a byte-preservation
// invariant in its fidelity table — "every byte outside the selected
// line(s) is unchanged; `diff` shows exactly one line per edited file" —
// and MVV step 3 restates it as "every other byte of both files is
// unchanged". So every positive step here:
//
//  1. reconstructs each edited file and compares it BYTE-FOR-BYTE against
//     a hand-written expected buffer, and
//  2. counts the changed lines and asserts exactly ONE per edited file
//     PER INVOCATION.
//
// The per-INVOCATION scoping is the record's, not a convenience: across
// the three fixture records the README accumulates three changed lines
// and each record one, so pinning a whole-run total would assert the
// fixture's SIZE instead of the contract.
//
// SCOPE OF THIS FILE vs Phase 4. `REQ-MVV-BLOCK` records that the
// consumer's row-addressed projector verb the C1.6 command reader invokes
// does not ship, so the CLI-level pipeline (`flow resolve | flow
// set-state`) cannot run end to end yet. This file therefore drives the
// same eight steps through the ACCESSOR seam — the executor, the edit
// write bindings, and a role reader per artifact — which is the deepest
// layer that does not depend on the missing verb. The dispositions,
// byte-invariants and refusal tokens asserted are the record's; what is
// deferred to Phase 4 is the envelope, not the behaviour.

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/flowbind"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// --- the MVV fixture ------------------------------------------------------

const (
	mvvFlow = "rdrwrite"

	mvvRecordRole = "record"
	mvvReadmeRole = "readme"

	mvvRecordWriter = "record.status"
	mvvRecordReader = "record.read"
	mvvReadmeWriter = "readme.row"
	mvvReadmeReader = "readme.read"

	mvvKey = "status"
)

// The three fixture records the MVV names: one plain `Draft`, one whose
// bracketed qualifier sits on ONE line, and one whose qualifier WRAPS
// onto a continuation line while still reading `Draft`. The wrapped one
// is A3's normative case — `0022-cache-metrics-surface.md`'s shape, NOT
// the already-`Final` `0021`, so the flip and the preservation are
// exercised together rather than one at a time.
const (
	mvvRecordPlain = `# 0027 — inline shell scope

- **Status**: Draft
- **Owner**: cwensel

Body text that must not move.
`

	mvvRecordJoint = `# 0026 — joint decision record

- **Status**: Draft [joint decision → JDR 0003]
- **Owner**: cwensel

Body text that must not move.
`

	// The qualifier wraps. The CONTINUATION line must come through
	// byte-identical, and the reader must report the qualifier whole.
	mvvRecordWrapped = `# 0022 — cache metrics surface

- **Status**: Draft [revised from Final 2026-08-10; re-verify A1,A2 — the
  qualifier continues onto this line and must not be touched]
- **Owner**: cwensel

Body text that must not move.
`
)

// The README index in the LINKED `| [NNNN](NNNN-slug.md) |` form that
// `rdr-write.toml`'s own `readme-add` rule emits. The unlinked `| NNNN |`
// form two live rows carry is out of this fixture by design and is what
// step 8 exercises.
const mvvReadme = `# Records

| id | title | status |
|---|---|---|
| [0022](0022-cache-metrics-surface.md) | cache metrics surface | Draft |
| [0026](0026-joint-decision.md) | joint decision | Draft |
| [0027](0027-inline-shell-scope.md) | inline shell scope | Draft |
`

// The record's Status anchor keeps any bracketed qualifier by
// backreference: group 1 is the bare value, group 2 the remainder of the
// line. That is what makes `Draft [joint decision → …]` become
// `Final [joint decision → …]` rather than a bare `Final`.
const (
	mvvRecordAnchor  = `^- \*\*Status\*\*: (\w+)(.*)$`
	mvvRecordReplace = `- **Status**: {status}${2}`

	// The README row is addressed by the context tag: `{tag.nnnn}` binds
	// ONE record identity, so it selects one row of the shared index.
	mvvReadmeAnchor  = `^\| \[{tag.nnnn}\]\(([^)]*)\) \| ([^|]*) \| (\w+) \|$`
	mvvReadmeReplace = `| [{tag.nnnn}](${1}) |${2}| {status} |`
)

// mvvWorld is one invocation's artifact pair plus its bindings.
type mvvWorld struct {
	recordPath string
	readmePath string
	recordBody string
	readmeBody string

	exec       *accessor.Executor
	recordEdit accessor.WriteBinding
	readmeEdit accessor.WriteBinding
}

// newMVVWorld lays down one record and one README on disk and builds the
// executor over them. `nnnn` is the record identity the invocation binds
// — the MVV's ONE INVOCATION PER RECORD discipline.
func newMVVWorld(
	t *testing.T, recordBody, readmeBody, nnnn string, allowCommands bool,
	recordRules, readmeRules map[string]table.EditRule,
	recordReads, readmeReads map[string]string,
) *mvvWorld {
	t.Helper()

	w := &mvvWorld{
		recordPath: writeFixture(t, "record.md", recordBody),
		readmePath: writeFixture(t, "README.md", readmeBody),
		recordBody: recordBody,
		readmeBody: readmeBody,
	}

	recAcc := table.Accessor{
		Role: mvvRecordRole, Keys: []string{mvvKey}, Timeout: "2s",
		ReadBack: true, Edit: recordRules,
	}
	rdmAcc := table.Accessor{
		Role: mvvReadmeRole, Keys: []string{mvvKey}, Timeout: "2s",
		ReadBack: true, Edit: readmeRules,
	}
	w.recordEdit = flowbind.NewEditWriter(recAcc, mvvRecordWriter)
	w.readmeEdit = flowbind.NewEditWriter(rdmAcc, mvvReadmeWriter)

	// Both readers are COMMAND-backed, exactly as the MVV fixture
	// declares: (a) `["rdr","status","-json","-filter","status",
	// "{artifact}"]` on role `record`, and (c) a C1.6 reader over the
	// README row carrying `{tag.nnnn}` in its argv. The stub stands in
	// for the consumer's projector verb, which REQ-MVV-BLOCK records as
	// not yet shipping; what it must NOT do is stand in for the write
	// binding, which is the unit under test.
	reg := accessor.Registry{
		Flow:          mvvFlow,
		OwnedTags:     []string{mvvKey},
		AllowCommands: allowCommands,
		Definitions: []accessor.Definition{
			{
				Identity: accessor.Identity{Flow: mvvFlow, Name: mvvRecordWriter, Capability: accessor.CapWrite},
				Accessor: recAcc,
				Binding:  w.recordEdit,
			},
			{
				Identity: accessor.Identity{Flow: mvvFlow, Name: mvvRecordReader, Capability: accessor.CapRead},
				Accessor: table.Accessor{
					Role:    mvvRecordRole,
					Command: []string{"rdr", "status", "-json", "-filter", "status", "{artifact}"},
					Keys:    []string{mvvKey},
					Timeout: "2s",
				},
				Binding: &mvvProjector{path: &w.recordPath, anchor: mvvRecordAnchor, nnnn: nnnn, canned: recordReads},
			},
			{
				Identity: accessor.Identity{Flow: mvvFlow, Name: mvvReadmeWriter, Capability: accessor.CapWrite},
				Accessor: rdmAcc,
				Binding:  w.readmeEdit,
			},
			{
				Identity: accessor.Identity{Flow: mvvFlow, Name: mvvReadmeReader, Capability: accessor.CapRead},
				Accessor: table.Accessor{
					Role:    mvvReadmeRole,
					Command: []string{"rdr", "readme-row", "{tag.nnnn}", "{artifact}"},
					Keys:    []string{mvvKey},
					Timeout: "2s",
				},
				Binding: &mvvProjector{path: &w.readmePath, nnnn: nnnn, canned: readmeReads},
			},
		},
	}

	w.exec = accessor.NewExecutor(reg, accessor.Artifacts{
		mvvRecordRole: {
			Role: mvvRecordRole, Path: w.recordPath,
			Context: map[string]string{"nnnn": nnnn},
		},
		mvvReadmeRole: {
			Role: mvvReadmeRole, Path: w.readmePath,
			Context: map[string]string{"nnnn": nnnn},
		},
	})
	return w
}

// assertUntouched asserts BOTH artifacts are byte-identical to what they
// were laid down as. Steps 5–8 are negative and this is their oracle.
func (w *mvvWorld) assertUntouched(t *testing.T, step string) {
	t.Helper()

	if got := readBack(t, w.recordPath); got != w.recordBody {
		t.Errorf("%s: the record changed:\n got %q\nwant %q", step, got, w.recordBody)
	}
	if got := readBack(t, w.readmePath); got != w.readmeBody {
		t.Errorf("%s: the README changed:\n got %q\nwant %q", step, got, w.readmeBody)
	}
}

// assertOneChangedLine asserts the file's post-state equals want
// byte-for-byte AND that exactly one line differs from its pre-state.
// Both halves are needed: byte equality alone would pass a two-line
// rewrite that happened to be expected, and a line count alone would pass
// a rewrite of the wrong line.
func assertOneChangedLine(t *testing.T, what, before, want, got string) {
	t.Helper()

	if got != want {
		t.Errorf("%s post-state:\n got %q\nwant %q", what, got, want)
		return
	}
	changed := changedLines(before, got)
	if len(changed) != 1 {
		t.Errorf("%s: %d lines changed %#v; the invariant is exactly ONE line "+
			"per edited file PER INVOCATION", what, len(changed), changed)
	}
}

// mvvProjector is the consumer's row-addressed projector standing in for
// the verb REQ-MVV-BLOCK records as not yet shipping. It reads the FILE
// back — it never mirrors what the writer intended — so the read-back is
// a genuine re-read of the artifact and can fail independently of the
// write, which is the property 0004's own fixtures insist on.
type mvvProjector struct {
	path   *string
	anchor string
	nnnn   string
	// canned overrides the file read where a step needs the reader to
	// report a specific value (the qualifier-split assertion).
	canned map[string]string
}

func (p *mvvProjector) Capability() accessor.Capability { return accessor.CapRead }

func (p *mvvProjector) Read(
	_ context.Context, art accessor.Artifact, requested []string,
) ([]accessor.KeyValue, []string, error) {
	// C1.6's whole reason for existing: without the record identity in
	// its argv, a reader over a SHARED artifact cannot say which row it
	// read. The context tag is that identity.
	id := art.Context["nnnn"]
	if id == "" {
		id = p.nnnn
	}

	body, err := readFile(art.Path)
	if err != nil {
		return nil, requested, &accessor.ExecError{Detail: err.Error(), Err: err}
	}

	var out []accessor.KeyValue
	for _, k := range requested {
		if v, ok := p.canned[k]; ok {
			out = append(out, accessor.KeyValue{Key: k, Value: v})
			continue
		}
		v, ok := projectStatus(body, id)
		out = append(out, accessor.KeyValue{Key: k, Value: v, Absent: !ok})
	}
	return out, nil, nil
}

// projectStatus reads the status the artifact holds for one record
// identity: the Status bullet's bare value in a record, or the status
// cell of the identity's row in the README.
func projectStatus(body, nnnn string) (string, bool) {
	for _, line := range strings.Split(body, "\n") {
		if after, ok := strings.CutPrefix(line, "- **Status**: "); ok {
			f := strings.Fields(after)
			if len(f) == 0 {
				return "", false
			}
			return f[0], true
		}
		if strings.HasPrefix(line, "| ["+nnnn+"](") {
			cells := strings.Split(line, "|")
			if len(cells) < 4 {
				return "", false
			}
			return strings.TrimSpace(cells[3]), true
		}
	}
	return "", false
}

func readFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	return string(b), err
}

// mvvProjectorVerbMissing reports whether the consumer's row-addressed
// projector verb is still absent. The probe is the verb itself: if `rdr`
// is not on PATH, or it does not accept `readme-row`, the verb does not
// ship. REQ-MVV-BLOCK names this exact dependency.
func mvvProjectorVerbMissing() bool {
	bin, err := exec.LookPath("rdr")
	if err != nil {
		return true
	}
	cmd := exec.Command(bin, "readme-row", "--help")
	return cmd.Run() != nil
}

// REQ-MVV (`0028:MVV`): "The consumer's acceptance scenario, carried from
// intrastate#zyh0." Fixture: "a `rdr-write.toml` re-authored as a
// state-machine over (a) one 0025 command reader
// `[\"rdr\",\"status\",\"-json\",\"-filter\",\"status\",\"{artifact}\"]` on
// role `record`, (b) an `edit` writer for the record's Status line on role
// `record`, (c) a C1.6 command reader over the README row on role `readme`
// (`{tag.nnnn}` in its argv) and an `edit` writer for that row anchored by
// `{tag.nnnn}`"
// HAPPY PATH
//
// The gating validation for RDR 0028. Eight steps, each a sub-test named
// for its REQ so an end-to-end failure traces to one step.
func TestMVV_DeclaredLineEditWriter(t *testing.T) {
	recordRules := map[string]table.EditRule{
		mvvKey: {Anchor: mvvRecordAnchor, Replace: mvvRecordReplace},
	}
	readmeRules := map[string]table.EditRule{
		mvvKey: {Anchor: mvvReadmeAnchor, Replace: mvvReadmeReplace},
	}

	// REQ-MVV.1: "`intrastate lint --model rdr-write.toml` passes; no
	// wrapper script exists anywhere in the fixture."
	t.Run("step1_the_model_lints_clean_with_no_wrapper_script", func(t *testing.T) {
		// The load-time gate. The `edit` writers carry NO `command` and
		// NO `path` — that absence IS the claim: the consumer's Status
		// flip needs no wrapper script anywhere in the fixture, which is
		// the whole problem this record exists to solve.
		if _, err := table.Load([]byte(mvvModelTOML), "rdr-write.toml"); err != nil {
			t.Fatalf("the MVV model did not lint clean: %v", err)
		}

		m, err := table.Load([]byte(mvvModelTOML), "rdr-write.toml")
		if err != nil {
			t.Fatalf("load: %v", err)
		}
		for name, acc := range m.Writers {
			if len(acc.Command) != 0 {
				t.Errorf("write.%s declares a command %#v; the MVV fixture "+
					"carries no wrapper script anywhere", name, acc.Command)
			}
			if acc.Path != "" {
				t.Errorf("write.%s declares a path %q; the writers are `edit`-carried",
					name, acc.Path)
			}
			if len(acc.Edit) == 0 {
				t.Errorf("write.%s carries no `edit` table", name)
			}
		}
	})

	// REQ-MVV.2: "ONE INVOCATION PER RECORD — `{tag.nnnn}` binds one
	// record identity, so it addresses one record file and one README
	// row; the three fixture records are three sequential runs of the
	// pipeline below, not one."
	// REQ-MVV.3: "End state: the record's Status line reads `Final` (the
	// joint-decision record reads `Final [joint decision → …]`, qualifier
	// kept by backreference; the wrapped record keeps its continuation
	// line byte-identical and its reader reports the qualifier whole); the
	// README row's status cell reads `Final`; both confirmed by read-back
	// through the command reader; every other byte of both files is
	// unchanged."
	t.Run("step2and3_three_sequential_invocations_each_flip_one_line_per_file", func(t *testing.T) {
		for _, tc := range []struct {
			nnnn       string
			record     string
			wantRecord string
		}{
			{
				nnnn:       "0027",
				record:     mvvRecordPlain,
				wantRecord: strings.Replace(mvvRecordPlain, "- **Status**: Draft", "- **Status**: Final", 1),
			},
			{
				// The qualifier is kept by BACKREFERENCE — group 2 of the
				// anchor carries ` [joint decision → JDR 0003]` forward.
				nnnn:   "0026",
				record: mvvRecordJoint,
				wantRecord: strings.Replace(mvvRecordJoint,
					"- **Status**: Draft [joint decision → JDR 0003]",
					"- **Status**: Final [joint decision → JDR 0003]", 1),
			},
			{
				// The wrapped record: only the BULLET line is rewritten
				// and the continuation line must come through
				// byte-identical.
				nnnn:   "0022",
				record: mvvRecordWrapped,
				wantRecord: strings.Replace(mvvRecordWrapped,
					"- **Status**: Draft [revised from Final 2026-08-10; re-verify A1,A2 — the",
					"- **Status**: Final [revised from Final 2026-08-10; re-verify A1,A2 — the", 1),
			},
		} {
			t.Run("record_"+tc.nnnn, func(t *testing.T) {
				w := newMVVWorld(t, tc.record, mvvReadme, tc.nnnn, true,
					recordRules, readmeRules, nil, nil)

				plan := resolve.Plan{Writes: []resolve.Tag{{Key: mvvKey, Value: "Final"}}}

				if res := w.exec.Write(context.Background(), mvvRecordWriter, plan); res.Refusal != nil {
					t.Fatalf("the record write refused: %+v", res.Refusal)
				}
				if res := w.exec.Write(context.Background(), mvvReadmeWriter, plan); res.Refusal != nil {
					t.Fatalf("the README write refused: %+v", res.Refusal)
				}

				// The record: exactly one changed line, every other byte
				// identical — including the wrapped record's continuation.
				assertOneChangedLine(t, "record "+tc.nnnn,
					tc.record, tc.wantRecord, readBack(t, w.recordPath))

				// The README: exactly one changed line — the addressed
				// row's status cell — and the two sibling rows untouched.
				wantReadme := strings.Replace(mvvReadme,
					"("+tc.nnnn+"-", "("+tc.nnnn+"-", 1)
				wantReadme = replaceRowStatus(mvvReadme, tc.nnnn, "Final")
				assertOneChangedLine(t, "README row "+tc.nnnn,
					mvvReadme, wantReadme, readBack(t, w.readmePath))

				// "both confirmed by read-back through the command
				// reader" — the executor's read-back already ran above and
				// did not refuse, which IS the confirmation; assert the
				// projected value too so a reader that verified nothing is
				// caught.
				body := readBack(t, w.recordPath)
				if got, _ := projectStatus(body, tc.nnnn); got != "Final" {
					t.Errorf("the record's reader projects %q; want Final", got)
				}
			})
		}
	})

	// REQ-MVV.3 (qualifier-whole half): "the wrapped record ... its reader
	// reports the qualifier whole"
	// REQ-122 (S26): "Read-back of a planned `Final` against a record
	// carrying a bracketed qualifier. **Expected**: byte-equal — the
	// reader reports `status=Final` with the qualifier on
	// `status_form`/`status.qualifier` (A4's normative fixture)."
	t.Run("step3_the_reader_splits_the_qualifier_and_reports_status_final", func(t *testing.T) {
		w := newMVVWorld(t, mvvRecordJoint, mvvReadme, "0026", true,
			recordRules, readmeRules, nil, nil)

		plan := resolve.Plan{Writes: []resolve.Tag{{Key: mvvKey, Value: "Final"}}}
		if res := w.exec.Write(context.Background(), mvvRecordWriter, plan); res.Refusal != nil {
			t.Fatalf("the write refused: %+v", res.Refusal)
		}

		body := readBack(t, w.recordPath)
		got, ok := projectStatus(body, "0026")
		if !ok || got != "Final" {
			t.Errorf("the reader reports status=%q (found=%v); want Final — the "+
				"qualifier rides `status_form`/`status.qualifier`, never the value",
				got, ok)
		}
		if !strings.Contains(body, "[joint decision → JDR 0003]") {
			t.Errorf("the qualifier was dropped from the line:\n%s", body)
		}
	})

	// REQ-MVV.4: "A `none` or `stopped:*` resolve row applies nothing and
	// exits 0 carrying `dispositions`."
	t.Run("step4_a_none_or_stopped_row_applies_nothing_and_does_not_refuse", func(t *testing.T) {
		w := newMVVWorld(t, mvvRecordPlain, mvvReadme, "0027", true,
			recordRules, readmeRules, nil, nil)
		recIno := inodeOf(t, w.recordPath)
		rdmIno := inodeOf(t, w.readmePath)

		// A resolve row that applies nothing reaches the write accessor
		// with an EMPTY write set. Performing no write is not a failure —
		// the disposition is carried, not a refusal.
		if res := w.exec.Write(context.Background(), mvvRecordWriter,
			resolve.Plan{}); res.Refusal != nil {
			t.Errorf("an empty plan refused: %+v; a `none`/`stopped:*` row applies "+
				"nothing and does not fail", res.Refusal)
		}
		if res := w.exec.Write(context.Background(), mvvReadmeWriter,
			resolve.Plan{}); res.Refusal != nil {
			t.Errorf("an empty plan refused: %+v", res.Refusal)
		}

		w.assertUntouched(t, "step4")
		if got := inodeOf(t, w.recordPath); got != recIno {
			t.Errorf("the record's inode moved on a no-apply row")
		}
		if got := inodeOf(t, w.readmePath); got != rdmIno {
			t.Errorf("the README's inode moved on a no-apply row")
		}
	})

	// REQ-MVV.5: "Negative: the same pipeline without `--allow-commands`
	// on `set-state` refuses before mutation and both files are
	// byte-identical to before."
	t.Run("step5_without_allow_commands_the_write_refuses_before_mutation", func(t *testing.T) {
		// allowCommands = false while BOTH readers are command-backed.
		w := newMVVWorld(t, mvvRecordPlain, mvvReadme, "0027", false,
			recordRules, readmeRules, nil, nil)
		recIno := inodeOf(t, w.recordPath)
		rdmIno := inodeOf(t, w.readmePath)

		plan := resolve.Plan{Writes: []resolve.Tag{{Key: mvvKey, Value: "Final"}}}

		for _, name := range []string{mvvRecordWriter, mvvReadmeWriter} {
			res := w.exec.Write(context.Background(), name, plan)
			if res.Refusal == nil {
				t.Fatalf("%s: the write ran with the gate OFF and a command-backed "+
					"reader", name)
			}
			if res.Refusal.Class == accessor.ClassReadBackIncomplete {
				t.Errorf("%s: class = read_back_incomplete; a forgotten flag must "+
					"not report an applied-but-unverified sense for a write that "+
					"ran no command", name)
			}
			if res.Refusal.Applied() {
				t.Errorf("%s: Applied() = true on a refusal decided BEFORE mutation",
					name)
			}
		}

		w.assertUntouched(t, "step5")
		if got := inodeOf(t, w.recordPath); got != recIno {
			t.Errorf("step5: the record's inode moved; the refusal is pre-mutation")
		}
		if got := inodeOf(t, w.readmePath); got != rdmIno {
			t.Errorf("step5: the README's inode moved; the refusal is pre-mutation")
		}
	})

	// REQ-MVV.6: "Negative: a README with the record's row duplicated
	// refuses `edit_anchor_ambiguous`; files untouched."
	t.Run("step6_a_duplicated_readme_row_refuses_edit_anchor_ambiguous", func(t *testing.T) {
		dup := mvvReadme + "| [0027](0027-inline-shell-scope.md) | inline shell scope | Draft |\n"
		w := newMVVWorld(t, mvvRecordPlain, dup, "0027", true,
			recordRules, readmeRules, nil, nil)
		rdmIno := inodeOf(t, w.readmePath)

		plan := resolve.Plan{Writes: []resolve.Tag{{Key: mvvKey, Value: "Final"}}}
		res := w.exec.Write(context.Background(), mvvReadmeWriter, plan)
		if res.Refusal == nil {
			t.Fatalf("a duplicated row applied; the anchor must select exactly one")
		}
		if d := res.Refusal.Detail; !strings.Contains(d, "edit_anchor_ambiguous") {
			t.Errorf("Detail = %q; want the token edit_anchor_ambiguous", d)
		}
		if got := readBack(t, w.readmePath); got != dup {
			t.Errorf("step6: the README changed on a refusal")
		}
		if got := inodeOf(t, w.readmePath); got != rdmIno {
			t.Errorf("step6: the README's inode moved on a refusal")
		}
	})

	// REQ-MVV.7: "Negative: a `replace` whose output no longer matches its
	// own anchor refuses `edit_anchor_unstable`; files untouched."
	t.Run("step7_a_self_de_anchoring_replace_refuses_edit_anchor_unstable", func(t *testing.T) {
		unstable := map[string]table.EditRule{
			// The anchor requires the `- **Status**: ` prefix; the
			// replacement drops it, so the rewritten line cannot be
			// re-selected on the post-edit buffer.
			mvvKey: {Anchor: mvvRecordAnchor, Replace: `Status is now {status}`},
		}
		w := newMVVWorld(t, mvvRecordPlain, mvvReadme, "0027", true,
			unstable, readmeRules, nil, nil)
		recIno := inodeOf(t, w.recordPath)

		plan := resolve.Plan{Writes: []resolve.Tag{{Key: mvvKey, Value: "Final"}}}
		res := w.exec.Write(context.Background(), mvvRecordWriter, plan)
		if res.Refusal == nil {
			t.Fatalf("a self-de-anchoring replace applied; the re-anchor pass " +
				"must refuse it BEFORE any write")
		}
		if d := res.Refusal.Detail; !strings.Contains(d, "edit_anchor_unstable") {
			t.Errorf("Detail = %q; want the token edit_anchor_unstable", d)
		}
		w.assertUntouched(t, "step7")
		if got := inodeOf(t, w.recordPath); got != recIno {
			t.Errorf("step7: the record's inode moved; the refusal is pre-write")
		}
	})

	// REQ-MVV.8: "Negative, and required rather than optional: a README
	// carrying one row in the live corpus's unlinked `| NNNN |` form (F6's
	// 0001/0006 shape) refuses `edit_anchor_unmatched` for that record,
	// names the rule in the Detail, and leaves both files byte-identical.
	// The other rows in the same README still flip."
	t.Run("step8_an_unlinked_row_refuses_unmatched_while_siblings_still_flip", func(t *testing.T) {
		// One row in the live corpus's unlinked form; the rest linked.
		drifted := `# Records

| id | title | status |
|---|---|---|
| 0001 | first record | Draft |
| [0026](0026-joint-decision.md) | joint decision | Draft |
| [0027](0027-inline-shell-scope.md) | inline shell scope | Draft |
`
		plan := resolve.Plan{Writes: []resolve.Tag{{Key: mvvKey, Value: "Final"}}}

		t.Run("the_unlinked_row_refuses_and_names_the_rule", func(t *testing.T) {
			w := newMVVWorld(t, mvvRecordPlain, drifted, "0001", true,
				recordRules, readmeRules, nil, nil)
			rdmIno := inodeOf(t, w.readmePath)

			res := w.exec.Write(context.Background(), mvvReadmeWriter, plan)
			if res.Refusal == nil {
				t.Fatalf("the unlinked `| NNNN |` row applied; the linked-form " +
					"anchor must not select it")
			}
			d := res.Refusal.Detail
			if !strings.Contains(d, "edit_anchor_unmatched") {
				t.Errorf("Detail = %q; want the token edit_anchor_unmatched", d)
			}
			if !strings.Contains(d, mvvReadmeWriter+".edit."+mvvKey) {
				t.Errorf("Detail = %q; want it to NAME the rule "+
					"`<id>.edit.<key>`", d)
			}
			if got := readBack(t, w.readmePath); got != drifted {
				t.Errorf("step8: the README changed on a refusal")
			}
			if got := inodeOf(t, w.readmePath); got != rdmIno {
				t.Errorf("step8: the README's inode moved on a refusal")
			}
		})

		t.Run("the_other_rows_in_the_same_readme_still_flip", func(t *testing.T) {
			// The drift is PER ROW: a README carrying one unlinked row
			// still serves every record whose row IS linked. Asserting
			// only the refusal would leave "one bad row poisons the whole
			// index" untested, which is the partial-adoption claim the
			// Prerequisites make.
			w := newMVVWorld(t, mvvRecordPlain, drifted, "0027", true,
				recordRules, readmeRules, nil, nil)

			if res := w.exec.Write(context.Background(), mvvReadmeWriter, plan); res.Refusal != nil {
				t.Fatalf("a LINKED row in a README carrying an unlinked sibling "+
					"refused: %+v", res.Refusal)
			}
			want := replaceRowStatus(drifted, "0027", "Final")
			assertOneChangedLine(t, "README row 0027",
				drifted, want, readBack(t, w.readmePath))
		})
	})
}

// replaceRowStatus rewrites the status cell of the linked row for nnnn,
// building the expected buffer by hand so the assertion never reuses the
// implementation's own substitution logic.
func replaceRowStatus(readme, nnnn, status string) string {
	lines := strings.Split(readme, "\n")
	for i, line := range lines {
		if !strings.HasPrefix(line, "| ["+nnnn+"](") {
			continue
		}
		cells := strings.Split(line, "|")
		if len(cells) < 4 {
			continue
		}
		cells[3] = " " + status + " "
		lines[i] = strings.Join(cells, "|")
	}
	return strings.Join(lines, "\n")
}

// mvvModelTOML is the MVV's `rdr-write.toml`, carrying every element the
// fixture paragraph names: a 0025 command reader on `record`, an `edit`
// writer for the record's Status line, a C1.6 command reader over the
// README row with `{tag.nnnn}` in its argv, and an `edit` writer for that
// row anchored by `{tag.nnnn}`. NO wrapper script appears anywhere.
const mvvModelTOML = `outcomes = ["lock"]
terminal = ["locked"]

[model]
id = "rdrwrite"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.nnnn]
provenance = "observed"
kind = "scalar"

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["Draft", "Final"]
single_valued = true
required = true

[read.record]
role = "record"
command = ["rdr", "status", "-json", "-filter", "status", "{artifact}"]
keys = ["status"]
timeout = "5s"

[read.readme]
role = "readme"
command = ["rdr", "readme-row", "{tag.nnnn}", "{artifact}"]
keys = ["status"]
timeout = "5s"

[write.record]
role = "record"
keys = ["status"]
timeout = "5s"
read_back = true

[write.record.edit.status]
anchor  = "^- \\*\\*Status\\*\\*: (\\w+)(.*)$"
replace = "- **Status**: {status}${2}"

[write.readme]
role = "readme"
keys = ["status"]
timeout = "5s"
read_back = true

[write.readme.edit.status]
anchor  = "^\\| \\[{tag.nnnn}\\]\\(([^)]*)\\) \\| ([^|]*) \\| (\\w+) \\|$"
replace = "| [{tag.nnnn}](${1}) |${2}| {status} |"

[context.done]
[context.done.match.status]
eq = "Final"

[initial]
status = "Draft"
`

// REQ-MVV-BLOCK: "BLOCKING the MVV (not the unit work): the consumer's
// row-addressed projector verb the C1.6 reader invokes must exist. The
// Illustrative Code marks it \"verb illustrative; the consumer's to add\"
// and no such verb ships today, so Phase 4's end-to-end proof cannot run
// until it does. Phases 1–3 do not depend on it"
// DOMAIN EDGE
//
// The blocker is recorded as a TEST so it cannot be forgotten: the CLI
// pipeline half of the MVV — `flow resolve | flow set-state` over the real
// `rdr readme-row` verb — is deferred, and this asserts the verb is still
// missing so the deferral stays honest. When the consumer ships the verb
// this test flips red and Phase 4 runs the pipeline end to end.
//
// The eight steps above are UNBLOCKED and run at the accessor seam, which
// is exactly what "Phases 1–3 do not depend on it" means.
func TestReqMVVBlock_ThePhase4PipelineAwaitsTheConsumersRowAddressedVerb(t *testing.T) {
	if !mvvProjectorVerbMissing() {
		t.Errorf("the consumer's row-addressed projector verb now ships; " +
			"MVV Phase 4 is unblocked and the CLI pipeline half of the MVV " +
			"(`flow resolve | flow set-state`) must now run end to end")
	}
}
