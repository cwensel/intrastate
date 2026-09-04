package accessor_test

// RDR 0028 — the `edit` write carrier's EXECUTION semantics (C1.2 apply
// half, C1.3, C1.5, C1.6 binding half) and its seam extension (A1).
//
// This file holds the APPLY half of C1. The load half — the declaration
// grammar and the six lint categories — lives in
// `internal/table/edit_carrier_0028_test.go`, because C1.4's
// `what lint does NOT prove:` fences the two apart: a model is linted
// without its artifacts and without the invocation's bindings, so no
// assertion here may be reachable from `table.Load` and none there may
// need a file.
//
// The coverage floors this file exists to hold.
//
// FIRST — BYTE PRESERVATION IS THE ORACLE, not "did not error". C1.3's
// fidelity table states that every byte outside the selected line is
// unchanged and that `diff` shows exactly one changed line per edited
// file. A green `Apply` over a file this suite never re-reads would pass
// a naive test while having rewritten the whole document, so every
// positive arm reconstructs the file and compares it byte-for-byte
// against a hand-written expected buffer.
//
// SECOND — "REFUSED" MEANS THE FILE IS UNTOUCHED. C1.3 `order:` decides
// every refusal BEFORE any byte is written. A refusal assertion that
// checked only the returned error would pass an implementation that
// wrote, then noticed, then reported — the exact failure the clause
// exists to prevent. So every negative arm asserts the target's bytes AND
// its inode are unchanged.
//
// THIRD — the six apply-time reason tokens ride the Detail and register
// NOWHERE in `table.Categories()` (req-list A-12); the load-time file
// holds the disjointness assertion.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/cmdbind"
	"github.com/cwensel/intrastate/internal/cli/flowbind"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// --- the apply-side fixture ----------------------------------------------

const (
	editFlow   = "editflow"
	editWriter = "record.write"
	editReader = "record.read"
	editRole   = "record"
	editKey    = "status"
)

// recordDoc is the consumer's record shape: a Status bullet among lines
// that must survive byte-identical. The decoy line is deliberate — a
// quoted example of the same bullet, which a first-match or last-match
// selector would silently pick (C1.3 `select:` rejects both).
const recordDoc = `# 0028 — declared line edit writer

- **Status**: Draft
- **Owner**: cwensel

## Notes

A record's status bullet is written as ` + "`- **Status**: Draft`" + ` in the
template, which is why the anchor pins ` + "`^…$`" + `.
`

// statusAnchor selects exactly the Status BULLET and captures the value
// plus any bracketed qualifier, so a `replace` can keep the qualifier by
// backreference (MVV step 3).
const statusAnchor = `^- \*\*Status\*\*: (\w+)(.*)$`

// statusReplace rewrites the whole line, interpolating the planned value
// and carrying the qualifier forward through group 2.
const statusReplace = `- **Status**: {status}${2}`

// writeFixture writes body to a fresh temp file and returns its path.
func writeFixture(t *testing.T, name, body string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write fixture %s: %v", name, err)
	}
	return path
}

// inodeOf returns the target's inode. Stage-and-rename ALWAYS replaces the
// inode (A5), so an unchanged inode is the observable that separates "did
// not write" from "wrote identical bytes" — file content cannot (S20).
func inodeOf(t *testing.T, path string) uint64 {
	t.Helper()

	fi, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		t.Skipf("inode is not observable on this platform")
	}
	return uint64(st.Ino)
}

// readBack returns the target's bytes as a string.
func readBack(t *testing.T, path string) string {
	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back %s: %v", path, err)
	}
	return string(b)
}

// changedLines reports the lines that differ between before and after,
// as `diff` would count them. C1.3's fidelity invariant is "exactly one
// line per edited file", so this is the assertion's unit.
func changedLines(before, after string) []string {
	b := strings.Split(before, "\n")
	a := strings.Split(after, "\n")
	var out []string
	for i := range max(len(a), len(b)) {
		var bl, al string
		if i < len(b) {
			bl = b[i]
		}
		if i < len(a) {
			al = a[i]
		}
		if bl != al {
			out = append(out, al)
		}
	}
	return out
}

// editRule is one declared line rule, mirroring `table.EditRule`.
func editRule(anchor, replace, clear string) table.EditRule {
	return table.EditRule{Anchor: anchor, Replace: replace, Clear: clear}
}

// editAccessor builds the `table.Accessor` an `edit` write entry carries:
// role and keys and timeout, NO path, NO command, and one rule per key.
func editAccessor(rules map[string]table.EditRule) table.Accessor {
	keys := slices.Sorted(editKeysOf(rules))
	return table.Accessor{
		Role:     editRole,
		Keys:     keys,
		Timeout:  "2s",
		ReadBack: true,
		Edit:     rules,
	}
}

func editKeysOf[V any](m map[string]V) func(func(string) bool) {
	return func(yield func(string) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}

// applyEdit invokes the edit binding directly over one artifact, which is
// the narrowest seam that exercises C1.3 without dragging the executor's
// read-back in. The executor-level arms have their own tests below.
func applyEdit(
	t *testing.T,
	rules map[string]table.EditRule,
	path string,
	tags map[string]string,
	planned ...resolve.Tag,
) (accessor.WriteBinding, error) {
	t.Helper()

	acc := editAccessor(rules)
	b := flowbind.NewEditWriter(acc, editWriter)
	art := accessor.Artifact{Role: editRole, Path: path, Context: tags}
	return b, b.Apply(context.Background(), art, planned)
}

// detailOf returns the Detail a refusal carries. C1.3 `order:` requires a
// rule-scoped refusal to name `<id>.edit.<key>` AND its reason token, so
// the Detail — not the error string — is the oracle.
func detailOf(t *testing.T, err error) string {
	t.Helper()

	if err == nil {
		t.Fatalf("want a refusal; got nil")
	}
	var ee *accessor.ExecError
	if !errors.As(err, &ee) {
		t.Fatalf("refusal is not an *accessor.ExecError: %v", err)
	}
	return ee.Detail
}

// assertRefusedUntouched asserts the apply refused, that its Detail names
// the reason token and the rule, and that not one byte and not the inode
// of the target moved. This is the shape C1.3 `order:` requires of EVERY
// refusal in the clause.
func assertRefusedUntouched(
	t *testing.T, err error, token, rule, path, before string, ino uint64,
) {
	t.Helper()

	detail := detailOf(t, err)
	if !strings.Contains(detail, token) {
		t.Errorf("Detail = %q; want it to carry the reason token %q", detail, token)
	}
	if rule != "" && !strings.Contains(detail, rule) {
		t.Errorf("Detail = %q; want it to name the rule %q — a rule-scoped "+
			"refusal names `<id>.edit.<key>`", detail, rule)
	}
	if got := readBack(t, path); got != before {
		t.Errorf("the target changed on a REFUSED apply:\n got %q\nwant %q — "+
			"every refusal in C1.3 is decided BEFORE any byte is written", got, before)
	}
	if got := inodeOf(t, path); got != ino {
		t.Errorf("the target's inode moved from %d to %d on a REFUSED apply; "+
			"stage-and-rename ran when nothing should have been written", ino, got)
	}
}

// --- C1.3 `input:` / `write:` — the happy path and its byte invariant ----

// REQ-31: "input:    the whole file is read as bytes; lines are split on
// \"\\n\" and a preceding \"\\r\" stays with the terminator (CRLF preserved
// per line); a missing final terminator is preserved; the file's mode is
// preserved on write (A5)."
// REQ-47: "write:    all rules of one entry rewrite ONE buffer and land in
// ONE write: staged beside the resolved target and renamed over it"
// REQ-3 (apply half): "replace ... the WHOLE replacement line, terminator
// excluded — never the matched span"
// HAPPY PATH
//
// The byte-preservation invariant, asserted as C1.3's fidelity table
// states it: exactly ONE changed line, every other byte identical. A green
// `Apply` is not the oracle — the reconstructed file is.
func TestReq31And47_AnAppliedEditRewritesExactlyTheAnchoredLine(t *testing.T) {
	path := writeFixture(t, "record.md", recordDoc)

	_, err := applyEdit(t,
		map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "")},
		path, nil,
		resolve.Tag{Key: editKey, Value: "Final"})
	if err != nil {
		t.Fatalf("a well-anchored edit refused: %v", err)
	}

	got := readBack(t, path)
	want := strings.Replace(recordDoc,
		"- **Status**: Draft\n", "- **Status**: Final\n", 1)
	if got != want {
		t.Errorf("post-edit file:\n%q\nwant:\n%q", got, want)
	}
	if changed := changedLines(recordDoc, got); len(changed) != 1 {
		t.Errorf("%d lines changed %#v; C1.3's fidelity invariant is exactly ONE "+
			"changed line per edited file", len(changed), changed)
	}
}

// REQ-3 (span half): "a span-style anchor without `^…$` drops the
// unmatched prefix/suffix by design"
// REQ-2: "a line is SELECTED when the pattern matches anywhere in it
// (terminator excluded); authors pin `^…$`"
// DOMAIN EDGE
//
// The clause states this as a DESIGNED consequence, not a defect: an
// unpinned anchor still selects the whole line and `replace` still emits
// the whole line, so the unmatched prefix and suffix are dropped. An
// implementation that spliced the replacement into the matched span
// instead would look more helpful and would violate the clause.
func TestReq2And3_ReplaceEmitsTheWholeLineNotTheMatchedSpan(t *testing.T) {
	path := writeFixture(t, "record.md", recordDoc)

	// An unpinned anchor matching only the middle of the Owner line.
	_, err := applyEdit(t,
		map[string]table.EditRule{editKey: editRule(`Owner`, `{status}`, "")},
		path, nil,
		resolve.Tag{Key: editKey, Value: "REPLACED"})
	if err != nil {
		t.Fatalf("an unpinned anchor refused: %v", err)
	}

	got := readBack(t, path)
	want := strings.Replace(recordDoc, "- **Owner**: cwensel\n", "REPLACED\n", 1)
	if got != want {
		t.Errorf("post-edit file:\n%q\nwant the WHOLE line replaced:\n%q", got, want)
	}
}

// REQ-30: "a line is the bytes up to and excluding `\\n`, with a preceding
// `\\r` treated as part of the terminator; the file is bytes, not a decoded
// string, and RE2 matches over UTF-8 bytes."
// REQ-56: "terminators: only \"\\n\" — optionally preceded by \"\\r\" —
// terminates a line for selection and rewriting; a bare \"\\r\", NEL or
// U+2028 is line content"
// REQ-57: "Deleting the final line of a file that had no final terminator
// also removes the preceding line's terminator, so the file's final-
// terminator state is preserved either way"
// INPUT EDGE
//
// S19. Four separate line-shape hazards, each of which an implementation
// built on `strings.Split(s, "\n")` plus `strings.Join` gets wrong in a
// different way: CRLF becomes LF, the missing final terminator grows one,
// and a bare `\r` becomes a line break.
func TestReq30And56And57_TerminatorHandlingPreservesLineShape(t *testing.T) {
	t.Run("crlf_is_preserved_per_line", func(t *testing.T) {
		body := "alpha\r\n- **Status**: Draft\r\nomega\r\n"
		path := writeFixture(t, "crlf.md", body)

		if _, err := applyEdit(t,
			map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "")},
			path, nil, resolve.Tag{Key: editKey, Value: "Final"}); err != nil {
			t.Fatalf("a CRLF fixture refused: %v", err)
		}

		want := "alpha\r\n- **Status**: Final\r\nomega\r\n"
		if got := readBack(t, path); got != want {
			t.Errorf("post-edit = %q; want %q — a preceding \\r stays with the "+
				"terminator and CRLF is preserved per line", got, want)
		}
	})

	t.Run("a_missing_final_terminator_is_preserved", func(t *testing.T) {
		body := "alpha\n- **Status**: Draft"
		path := writeFixture(t, "noterm.md", body)

		if _, err := applyEdit(t,
			map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "")},
			path, nil, resolve.Tag{Key: editKey, Value: "Final"}); err != nil {
			t.Fatalf("an unterminated fixture refused: %v", err)
		}

		want := "alpha\n- **Status**: Final"
		if got := readBack(t, path); got != want {
			t.Errorf("post-edit = %q; want %q — a missing final terminator is "+
				"preserved, never supplied", got, want)
		}
	})

	t.Run("bare_cr_nel_and_u2028_are_line_content", func(t *testing.T) {
		// One physical line carrying all three. If any were treated as a
		// terminator the anchor — pinned `^…$` — would not select it.
		body := "- **Status**: Draft\rtail\u0085more\u2028end\n"
		path := writeFixture(t, "oddities.md", body)

		if _, err := applyEdit(t,
			map[string]table.EditRule{editKey: editRule(`^- \*\*Status\*\*: Draft.*$`, `{status}`, "")},
			path, nil, resolve.Tag{Key: editKey, Value: "X"}); err != nil {
			t.Fatalf("a fixture carrying bare \\r/NEL/U+2028 refused: %v", err)
		}

		want := "X\n"
		if got := readBack(t, path); got != want {
			t.Errorf("post-edit = %q; want %q — a bare \\r, NEL and U+2028 are "+
				"line CONTENT, so the whole run is one line", got, want)
		}
	})

	t.Run("deleting_an_unterminated_final_line_drops_the_preceding_terminator", func(t *testing.T) {
		body := "alpha\n- **Status**: Draft"
		path := writeFixture(t, "noterm-clear.md", body)

		if _, err := applyEdit(t,
			map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "line")},
			path, nil,
			resolve.Tag{Key: editKey, Value: accessor.ClearSentinel}); err != nil {
			t.Fatalf("a clear on an unterminated final line refused: %v", err)
		}

		want := "alpha"
		if got := readBack(t, path); got != want {
			t.Errorf("post-clear = %q; want %q — deleting the final line of a "+
				"file that had no final terminator also removes the preceding "+
				"line's terminator", got, want)
		}
	})
}

// REQ-108 (S21): "Mode preservation on a 0644 target and on a 0755 target.
// **Expected**: `git diff --summary` empty in both cases."
// BOUNDARY
//
// The reason mode is copied rather than reused from `flowbind.go::save`,
// whose fixed `Chmod(0o600)` would emit ` mode change 100755 => 100644` on
// the 0755 target — A5's fixture, and the mutant this kills.
func TestReq108_TheTargetsModeIsPreservedAcrossTheWrite(t *testing.T) {
	for _, mode := range []os.FileMode{0o644, 0o755, 0o600} {
		t.Run(mode.String(), func(t *testing.T) {
			path := writeFixture(t, "moded.md", recordDoc)
			if err := os.Chmod(path, mode); err != nil {
				t.Fatalf("chmod: %v", err)
			}
			ino := inodeOf(t, path)

			if _, err := applyEdit(t,
				map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "")},
				path, nil, resolve.Tag{Key: editKey, Value: "Final"}); err != nil {
				t.Fatalf("apply refused: %v", err)
			}

			// The write must actually have HAPPENED, or "the mode did
			// not change" is vacuously true of a binding that did
			// nothing. Content plus inode is the witness.
			want := strings.Replace(recordDoc, "Draft", "Final", 1)
			if got := readBack(t, path); got != want {
				t.Fatalf("the edit did not land:\n%q", got)
			}
			if got := inodeOf(t, path); got == ino {
				t.Fatalf("the inode did not move; stage-and-rename did not run, " +
					"so mode preservation is untested here")
			}

			fi, err := os.Stat(path)
			if err != nil {
				t.Fatalf("stat: %v", err)
			}
			if got := fi.Mode().Perm(); got != mode {
				t.Errorf("mode = %v after the write; want the ORIGINAL %v — a "+
					"fixed chmod would show as a mode change in `git diff --summary`",
					got, mode)
			}
		})
	}
}

// REQ-33: "target:   the caller-bound artifact path for the entry's role
// (0004:C3), symlinks resolved; the model names no path"
// REQ-109 (S22): "A symlinked target. **Expected**: the symlink survives
// and points at the new content"
// DOMAIN EDGE
//
// Renaming onto an UNRESOLVED symlink path replaces the symlink with a
// regular file — a property of `os.Rename` A5 demonstrated. C1.3 forbids
// the unresolved path unconditionally, so the symlink must survive AS a
// symlink and its target must carry the new content.
func TestReq33And109_ASymlinkedTargetIsResolvedBeforeStagingAndSurvives(t *testing.T) {
	dir := t.TempDir()
	real := filepath.Join(dir, "real.md")
	if err := os.WriteFile(real, []byte(recordDoc), 0o644); err != nil {
		t.Fatalf("write real: %v", err)
	}
	link := filepath.Join(dir, "link.md")
	if err := os.Symlink(real, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	if _, err := applyEdit(t,
		map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "")},
		link, nil, resolve.Tag{Key: editKey, Value: "Final"}); err != nil {
		t.Fatalf("apply through a symlink refused: %v", err)
	}

	fi, err := os.Lstat(link)
	if err != nil {
		t.Fatalf("lstat link: %v", err)
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		t.Errorf("the symlink was replaced by a regular file; C1.3 `target:` " +
			"resolves symlinks BEFORE staging so the link survives")
	}
	want := strings.Replace(recordDoc, "Draft", "Final", 1)
	if got := readBack(t, real); got != want {
		t.Errorf("the symlink's target:\n%q\nwant the new content:\n%q", got, want)
	}
}

// REQ-49: "A post-edit buffer equal to the input is not written at all (no
// staging, no rename; S20 asserts it, witnessed by an unchanged inode"
// REQ-107 (S20): "The witness is the target's INODE, unchanged across the
// call ... stage-and-rename always replaces the inode (A5), so an
// unchanged inode is the observable that proves absence of a write."
// req-list A-10: the no-op arm is evaluated AFTER the re-anchor pass.
// BOUNDARY
//
// An already-applied edit re-run against its own output. File CONTENT
// cannot separate "did not write" from "wrote identical bytes" — only the
// inode can, which is why the clause names it as the witness.
func TestReq49And107_ANoOpEditIsNotWrittenAtAllWitnessedByTheInode(t *testing.T) {
	rules := map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "")}

	// POSITIVE CONTROL. An unchanged inode only means "declined to write"
	// if this binding demonstrably DOES write when the buffer differs; a
	// binding that never wrote would pass the no-op assertion vacuously.
	{
		ctl := writeFixture(t, "control.md", recordDoc)
		ctlIno := inodeOf(t, ctl)
		if _, err := applyEdit(t, rules, ctl, nil,
			resolve.Tag{Key: editKey, Value: "Final"}); err != nil {
			t.Fatalf("control apply refused: %v", err)
		}
		if inodeOf(t, ctl) == ctlIno {
			t.Fatalf("the control write did not move the inode; stage-and-rename " +
				"never ran, so the no-op assertion below would be vacuous")
		}
	}

	already := strings.Replace(recordDoc, "Draft", "Final", 1)
	path := writeFixture(t, "already.md", already)
	before := inodeOf(t, path)

	if _, err := applyEdit(t, rules, path, nil,
		resolve.Tag{Key: editKey, Value: "Final"}); err != nil {
		t.Fatalf("re-running an applied edit refused: %v", err)
	}

	if got := readBack(t, path); got != already {
		t.Errorf("content changed on a no-op:\n%q\nwant\n%q", got, already)
	}
	if got := inodeOf(t, path); got != before {
		t.Errorf("inode moved from %d to %d; a post-edit buffer equal to the "+
			"input is not written at all — no staging, no rename", before, got)
	}
}

// --- C1.3 `select:` — cardinality and its refusals ------------------------

// REQ-35: "select:   every rule's anchor is resolved against the PRE-EDIT
// content and selections are held as pre-edit line INDICES (a deletion
// never shifts a sibling rule's target)."
// REQ-36: "Each must select exactly one line: 0 ⇒ `edit_anchor_unmatched`;
// ≥2 ⇒ `edit_anchor_ambiguous`; two rules selecting one line ⇒
// `edit_anchor_collision`."
// REQ-37: "Never last-match (Ansible `lineinfile`), never first-match,
// never insert or append (Puppet `append_on_no_match`, Ansible
// `insertafter`) — creation is fenced out"
// ADVERSARIAL
//
// S9. Each refusal is asserted BEFORE any byte is written, with the file
// byte-identical after. The zero-match arm is where the clause is most
// opinionated: every neighbouring tool would APPEND the line, and doing so
// here would convert a stale model into a silently invented one.
func TestReq35And36And37_EachRuleMustSelectExactlyOneLine(t *testing.T) {
	t.Run("zero_matches_is_edit_anchor_unmatched_and_never_an_append", func(t *testing.T) {
		path := writeFixture(t, "record.md", recordDoc)
		ino := inodeOf(t, path)

		_, err := applyEdit(t,
			map[string]table.EditRule{editKey: editRule(`^- \*\*Absent\*\*: (.*)$`, `{status}`, "")},
			path, nil, resolve.Tag{Key: editKey, Value: "Final"})

		assertRefusedUntouched(t, err, "edit_anchor_unmatched",
			editWriter+".edit."+editKey, path, recordDoc, ino)
	})

	t.Run("two_or_more_matches_is_edit_anchor_ambiguous", func(t *testing.T) {
		dup := recordDoc + "- **Status**: Draft\n"
		path := writeFixture(t, "dup.md", dup)
		ino := inodeOf(t, path)

		_, err := applyEdit(t,
			map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "")},
			path, nil, resolve.Tag{Key: editKey, Value: "Final"})

		assertRefusedUntouched(t, err, "edit_anchor_ambiguous",
			editWriter+".edit."+editKey, path, dup, ino)
	})

	t.Run("two_rules_on_one_line_is_edit_anchor_collision", func(t *testing.T) {
		path := writeFixture(t, "record.md", recordDoc)
		ino := inodeOf(t, path)

		// Two rules whose anchors both select the Status bullet.
		_, err := applyEdit(t, map[string]table.EditRule{
			editKey: editRule(statusAnchor, statusReplace, ""),
			"owner": editRule(`^- \*\*Status\*\*: Draft$`, `{owner}`, ""),
		}, path, nil,
			resolve.Tag{Key: editKey, Value: "Final"},
			resolve.Tag{Key: "owner", Value: "someone"})

		assertRefusedUntouched(t, err, "edit_anchor_collision", "", path, recordDoc, ino)
	})
}

// REQ-105 (S13): "A record whose Status VALUE sits on the continuation
// line. **Expected**: zero matches → `edit_anchor_unmatched`, not a
// rewrite."
// DOMAIN EDGE
//
// The consumer's real drift shape. The bullet exists but its VALUE wrapped
// onto the next line, so the pinned anchor — which requires the value on
// the bullet line — matches zero. Rewriting the bullet anyway would
// produce a document with two Status values.
func TestReq105_AStatusValueOnTheContinuationLineIsUnmatchedNotRewritten(t *testing.T) {
	body := "# rec\n\n- **Status**:\n  Draft [wrapped]\n\ntail\n"
	path := writeFixture(t, "wrapped-value.md", body)
	ino := inodeOf(t, path)

	_, err := applyEdit(t,
		map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "")},
		path, nil, resolve.Tag{Key: editKey, Value: "Final"})

	assertRefusedUntouched(t, err, "edit_anchor_unmatched",
		editWriter+".edit."+editKey, path, body, ino)
}

// REQ-104 (S8): "exactly one line selected per in-scope file; no decoy
// (template comment, quoted example) scores ≥2." ... "The invariant is PER
// FILE and the cardinalities are descriptive, never asserted"
// DOMAIN EDGE
//
// The A3 corpus check as a regression test, held to the clause's own
// discipline: the invariant is PER FILE, and no total is pinned, because
// the corpus grows with every RDR seeded and a pinned total is a scheduled
// false positive. The decoy in `recordDoc` — a quoted `- **Status**: Draft`
// inside a backticked sentence — is the shape S8 names.
func TestReq104_ThePinnedAnchorSelectsExactlyOneLinePerInScopeFile(t *testing.T) {
	// In-scope means Status-bearing: a record carrying a `- **Status**:`
	// bullet. A file without one legitimately selects zero and is out of
	// scope, which is why the predicate is part of the assertion.
	for _, tc := range []struct {
		name    string
		body    string
		inScope bool
	}{
		{name: "bullet_plus_a_quoted_decoy", body: recordDoc, inScope: true},
		{
			name:    "bullet_plus_an_indented_template_comment",
			body:    "- **Status**: Draft\n\n    - **Status**: <value>\n",
			inScope: true,
		},
		{name: "no_status_bullet_is_out_of_scope", body: "# postmortem\n\nno bullet\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := writeFixture(t, "corpus.md", tc.body)
			_, err := applyEdit(t,
				map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "")},
				path, nil, resolve.Tag{Key: editKey, Value: "Final"})

			if tc.inScope {
				if err != nil {
					t.Errorf("an in-scope Status-bearing file refused: %v — the "+
						"pinned anchor must select exactly one line and no decoy "+
						"may score with it", err)
				}
				return
			}
			// Out of scope: zero matches, which is the unmatched refusal
			// and NOT a defect of the anchor.
			if err == nil {
				t.Errorf("a file carrying no Status bullet applied cleanly; " +
					"out-of-scope files select zero")
			}
		})
	}
}

// REQ-119 (S31): "A decoy fixture where the anchor's exactly-one match is
// the WRONG line ... **Expected**: the write applies, and the read-back
// through the role's reader then refuses `read_back_mismatch`. `Applied()`
// is FALSE on this arm"
// ADVERSARIAL
//
// The post-mutation net F4 and F5 name as the SOLE defence against silent
// wrong-line corruption. Pre-mutation refusal cannot reach this class by
// construction, which is why it is tested here rather than folded into the
// refusal-surface test.
//
// `Applied()` is FALSE and that is deliberate: `Executor.Write`'s mismatch
// branch leaves the applied sense unset by an explicit 0004:C14 scoping
// comment, which confines applied-but-unverified to `read_back_incomplete`
// and a post-mutation timeout. A test asserting TRUE here would fail for
// the correct reason and must not be relaxed into deleting the check.
func TestReq119_AWrongLineRewriteIsCaughtByReadBackMismatchWithAppliedFalse(t *testing.T) {
	// The decoy outscores the real bullet: the anchor pins a shape only
	// the quoted example carries.
	body := "# rec\n\n- **Status**: Draft\n\n> example: `Status = Draft`\n"
	path := writeFixture(t, "decoy.md", body)

	acc := editAccessor(map[string]table.EditRule{
		editKey: editRule("^> example: .Status = (\\w+).$", "> example: `Status = {status}`", ""),
	})
	writer := flowbind.NewEditWriter(acc, editWriter)

	// The reader reports the REAL bullet, which the edit never touched, so
	// the read-back compares Draft against the planned Final.
	reader := &editStubReader{values: map[string]string{editKey: "Draft"}}

	reg := accessor.Registry{
		Flow:      editFlow,
		OwnedTags: []string{editKey},
		Definitions: []accessor.Definition{
			{
				Identity: accessor.Identity{Flow: editFlow, Name: editWriter, Capability: accessor.CapWrite},
				Accessor: acc,
				Binding:  writer,
			},
			{
				Identity: accessor.Identity{Flow: editFlow, Name: editReader, Capability: accessor.CapRead},
				Accessor: table.Accessor{Role: editRole, Keys: []string{editKey}, Timeout: "2s"},
				Binding:  reader,
			},
		},
	}
	ex := accessor.NewExecutor(reg, accessor.Artifacts{
		editRole: {Role: editRole, Path: path},
	})

	res := ex.Write(context.Background(), editWriter,
		resolve.Plan{Writes: []resolve.Tag{{Key: editKey, Value: "Final"}}})

	if res.Refusal == nil {
		t.Fatalf("a wrong-line rewrite verified clean; read-back is the only " +
			"defence against silent wrong-line corruption")
	}
	if got := res.Refusal.Class; got != accessor.ClassReadBackMismatch {
		t.Errorf("class = %q; want %q", got, accessor.ClassReadBackMismatch)
	}
	if res.Refusal.Applied() {
		t.Errorf("Applied() = true on a read_back_mismatch; 0004:C14 confines " +
			"the applied-but-unverified sense to read_back_incomplete and a " +
			"post-mutation timeout")
	}

	// "the write APPLIES, and the read-back then refuses" — so the DECOY
	// line must actually carry the new bytes. Without this the mismatch
	// could be caused by a binding that wrote nothing at all, which is a
	// different arm and would leave the wrong-line hazard untested.
	after := readBack(t, path)
	if !strings.Contains(after, "> example: `Status = Final`") {
		t.Errorf("the decoy line was not rewritten:\n%q — the write must APPLY "+
			"for the read-back to be the witness that it hit the wrong line", after)
	}
	if !strings.Contains(after, "- **Status**: Draft") {
		t.Errorf("the real bullet was touched:\n%q — the anchor selected the "+
			"decoy, so only the decoy may change", after)
	}
}

// --- C1.3 `re-anchor:` ----------------------------------------------------

// REQ-38: "re-anchor: after the buffer is rewritten in memory, every
// rule's anchor is run again over the POST-EDIT buffer and must select
// exactly its own rewritten line (or, for a deleted line, zero lines);
// otherwise refuse `edit_anchor_unstable` before any write."
// ADVERSARIAL
//
// S11, premortem P-4 and P-13. A `replace` whose output no longer matches
// its own anchor de-anchors the rule: the next run would report
// `edit_anchor_unmatched` against a line this run wrote. Refusing here is
// what makes the carrier idempotent.
func TestReq38_AReplacementThatDeAnchorsItselfRefusesEditAnchorUnstable(t *testing.T) {
	path := writeFixture(t, "record.md", recordDoc)
	ino := inodeOf(t, path)

	// The anchor requires the `- **Status**: ` prefix; the replacement
	// drops it, so the rewritten line cannot be re-selected.
	_, err := applyEdit(t,
		map[string]table.EditRule{editKey: editRule(statusAnchor, `Status is {status}`, "")},
		path, nil, resolve.Tag{Key: editKey, Value: "Final"})

	assertRefusedUntouched(t, err, "edit_anchor_unstable",
		editWriter+".edit."+editKey, path, recordDoc, ino)
}

// REQ-39: "IDENTITY, not cardinality: the selected line must BE the rule's
// own — its held pre-edit index, shifted by the deletions of preceding
// sibling rules — and a rule that selects exactly one line which is a
// DIFFERENT line refuses."
// req-list A-9: a `<clear>` on a zero-match anchor contributes no deletion
// to the shift.
// ADVERSARIAL
//
// The mutant this kills is a re-anchor pass that checks CARDINALITY only.
// Such a pass would pass a rule whose anchor drifted onto ANOTHER rule's
// line after a rewrite — exactly one match, wrong line — which is how a
// replacement poisons a sibling's anchor on the next run.
func TestReq39_ReAnchorAssertsLineIdentityNotMerelyCardinality(t *testing.T) {
	body := "- **A**: one\n- **B**: two\n"
	path := writeFixture(t, "identity.md", body)
	ino := inodeOf(t, path)

	// Rule `status` rewrites the A line into the SHAPE rule `owner`
	// anchors on. After the rewrite each anchor still selects exactly one
	// line, but `owner`'s selection is no longer its own.
	_, err := applyEdit(t, map[string]table.EditRule{
		editKey: editRule(`^- \*\*A\*\*: (.*)$`, `- **B**: {status}`, ""),
		"owner": editRule(`^- \*\*B\*\*: (.*)$`, `- **B**: {owner}`, ""),
	}, path, nil,
		resolve.Tag{Key: editKey, Value: "moved"},
		resolve.Tag{Key: "owner", Value: "two"})

	assertRefusedUntouched(t, err, "edit_anchor_unstable", "", path, body, ino)
}

// --- C1.2 `value shape:` and `parse once:` --------------------------------

// REQ-20: "value shape: a planned value, or a bound tag value THIS ENTRY's
// rules reference, containing \"\\n\" or \"\\r\" refuses BEFORE mutation
// (Detail `edit_value_multiline`) — the single structural hazard of
// interpolation into line data."
// ADVERSARIAL
//
// S10. A newline in an interpolated value turns one line into two, which
// is the only way data can become GRAMMAR on this path — there is no
// shell, no word-splitting and no option parsing here (REQ-21), so this is
// the whole hazard class.
func TestReq20_AValueCarryingANewlineOrCarriageReturnRefusesBeforeMutation(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value string
	}{
		{name: "planned_value_with_lf", value: "Fi\nnal"},
		{name: "planned_value_with_cr", value: "Fi\rnal"},
		{name: "planned_value_with_crlf", value: "Fi\r\nnal"},
		{name: "planned_value_that_is_only_a_newline", value: "\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := writeFixture(t, "record.md", recordDoc)
			ino := inodeOf(t, path)

			_, err := applyEdit(t,
				map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "")},
				path, nil, resolve.Tag{Key: editKey, Value: tc.value})

			assertRefusedUntouched(t, err, "edit_value_multiline", "", path, recordDoc, ino)
		})
	}

	t.Run("a_referenced_tag_value_with_a_newline", func(t *testing.T) {
		body := "| [0028](0028-x.md) | Draft |\n"
		path := writeFixture(t, "readme.md", body)
		ino := inodeOf(t, path)

		_, err := applyEdit(t,
			map[string]table.EditRule{editKey: editRule(`^\| \[{tag.nnnn}\].*$`, `{status}`, "")},
			path, map[string]string{"nnnn": "00\n28"},
			resolve.Tag{Key: editKey, Value: "Final"})

		assertRefusedUntouched(t, err, "edit_value_multiline", "", path, body, ino)
	})
}

// REQ-52 (S24 half), REQ-111: "the same value bound but named by no argv
// element of the entry; the same value reaching an `anchor` only."
// ... "bound-but-unreferenced is never scanned (per-use-site)"
// DOMAIN EDGE
//
// The per-USE-SITE scope. A tag bound on the invocation context but
// referenced by no `anchor` of this entry reaches no line data, so scanning
// it would fence out a legitimate invocation for a value the entry never
// touches. The mutant this kills is a per-INVOCATION scan.
func TestReq52And111_ABoundButUnreferencedTagValueIsNeverScanned(t *testing.T) {
	path := writeFixture(t, "record.md", recordDoc)

	// `other` carries a newline — a hazard if it were scanned — but no
	// anchor of this entry references it.
	if _, err := applyEdit(t,
		map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "")},
		path, map[string]string{"other": "haz\nard"},
		resolve.Tag{Key: editKey, Value: "Final"}); err != nil {
		t.Fatalf("a bound-but-unreferenced tag caused a refusal: %v — the scope "+
			"is per-USE-SITE, not per-invocation", err)
	}

	want := strings.Replace(recordDoc, "Draft", "Final", 1)
	if got := readBack(t, path); got != want {
		t.Errorf("post-edit:\n%q\nwant:\n%q", got, want)
	}
}

// REQ-19 (apply half), REQ-23: "captured text and substituted values are
// never re-scanned for `${…}` or `{…}`" ... "the parse is this clause's own
// and never a call to `Expand`, which re-scans its template at expansion
// time and would break `parse once:`"
// REQ-15 (S15): "Captured text and substituted values that themselves
// contain `${…}` or `{…}`. **Expected**: emitted as literal bytes"
// ADVERSARIAL
//
// The injection boundary. If captured text were re-scanned, a document
// line containing `${1}` would let the DOCUMENT choose what gets written —
// runtime data becoming grammar, which is exactly what C1.2 forbids.
// `regexp.Expand` re-scans, so an implementation that delegated to it
// fails here and nowhere else in this suite.
func TestReq19And23_CapturedTextAndValuesEmitAsLiteralBytesNeverRescanned(t *testing.T) {
	t.Run("captured_text_carrying_a_group_reference", func(t *testing.T) {
		body := "- **Status**: Draft ${1} {status} $$\n"
		path := writeFixture(t, "rescan.md", body)

		if _, err := applyEdit(t,
			map[string]table.EditRule{editKey: editRule(`^- \*\*Status\*\*: \w+ (.*)$`, `- **Status**: {status} ${1}`, "")},
			path, nil, resolve.Tag{Key: editKey, Value: "Final"}); err != nil {
			t.Fatalf("apply refused: %v", err)
		}

		// Group 1 captured the literal bytes `${1} {status} $$`; they emit
		// unchanged.
		want := "- **Status**: Final ${1} {status} $$\n"
		if got := readBack(t, path); got != want {
			t.Errorf("post-edit = %q; want %q — captured text is never "+
				"re-scanned for `${…}` or `{…}`", got, want)
		}
	})

	t.Run("a_planned_value_carrying_a_placeholder", func(t *testing.T) {
		path := writeFixture(t, "record.md", recordDoc)

		if _, err := applyEdit(t,
			map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "")},
			path, nil, resolve.Tag{Key: editKey, Value: "${1}{status}"}); err != nil {
			t.Fatalf("apply refused: %v", err)
		}

		want := strings.Replace(recordDoc, "- **Status**: Draft", "- **Status**: ${1}{status}", 1)
		if got := readBack(t, path); got != want {
			t.Errorf("post-edit:\n%q\nwant:\n%q — a substituted VALUE is never "+
				"re-scanned", got, want)
		}
	})
}

// REQ-22 (apply half): "escapes: `$$` emits a literal `$`; `{{` and `}}`
// emit literal braces in `replace`; a group that did not participate in
// the match expands to the empty string"
// REQ-121 (S14, replace half)
// INPUT EDGE
//
// The empty-group rule MATCHES `regexp.Expand`'s, which is exactly why it
// must be asserted: an implementation could get this arm right by
// delegating to `Expand` while breaking `parse once:` above. Both tests
// have to pass, and only a hand-rolled segment parser passes both.
func TestReq22And121_ReplaceEscapesAndNonParticipatingGroupsExpandAsStated(t *testing.T) {
	for _, tc := range []struct {
		name    string
		anchor  string
		replace string
		body    string
		want    string
	}{
		{
			name:    "doubled_dollar_emits_one_dollar",
			anchor:  `^X$`,
			replace: `$${status}`,
			body:    "X\n",
			want:    "$Final\n",
		},
		{
			name:    "doubled_braces_emit_literal_braces",
			anchor:  `^X$`,
			replace: `{{{status}}}`,
			body:    "X\n",
			want:    "{Final}\n",
		},
		{
			name:    "a_group_that_did_not_participate_expands_empty",
			anchor:  `^X(a)?$`,
			replace: `{status}[${1}]`,
			body:    "X\n",
			want:    "Final[]\n",
		},
		{
			name:    "a_group_that_participated_expands_to_its_text",
			anchor:  `^X(a)?$`,
			replace: `{status}[${1}]`,
			body:    "Xa\n",
			want:    "Final[a]\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := writeFixture(t, "escapes.md", tc.body)

			if _, err := applyEdit(t,
				map[string]table.EditRule{editKey: editRule(tc.anchor, tc.replace, "")},
				path, nil, resolve.Tag{Key: editKey, Value: "Final"}); err != nil {
				t.Fatalf("apply refused: %v", err)
			}
			if got := readBack(t, path); got != tc.want {
				t.Errorf("post-edit = %q; want %q", got, tc.want)
			}
		})
	}
}

// REQ-16: "anchor  admits: ... {tag.<key>} — a tag key the model declares
// with `observed` provenance, bound on the invocation's context (A1, A7);
// substituted regexp-quoted, so a bound value is a literal-match fragment
// and can never alter the pattern's structure."
// REQ-106 (S16): "A bound tag value containing RE2 metacharacters, used in
// an anchor. **Expected**: regexp-quoted — matches literally, cannot alter
// the pattern's structure (C1.2)."
// ADVERSARIAL
//
// The structural guarantee. Without quoting, a bound value of `.*` would
// turn a pinned anchor into a match-anything pattern and the caller would
// choose which line gets rewritten — the same class of defect as an
// unquoted SQL parameter.
func TestReq16And106_ABoundTagValueIsRegexpQuotedAndCannotAlterThePattern(t *testing.T) {
	t.Run("a_metacharacter_value_matches_literally", func(t *testing.T) {
		body := "| [a.c](x.md) | Draft |\n| [abc](y.md) | Draft |\n"
		path := writeFixture(t, "readme.md", body)

		// Unquoted, `a.c` would match BOTH rows and refuse ambiguous.
		// Quoted, it matches the literal `a.c` row alone.
		if _, err := applyEdit(t,
			map[string]table.EditRule{editKey: editRule(
				`^\| \[{tag.nnnn}\]\([^)]*\) \| (\w+) \|$`,
				`| [{tag.nnnn}](x.md) | {status} |`, "")},
			path, map[string]string{"nnnn": "a.c"},
			resolve.Tag{Key: editKey, Value: "Final"}); err != nil {
			t.Fatalf("a metacharacter-bearing tag value refused: %v — it is "+
				"regexp-QUOTED, so it matches literally", err)
		}

		want := "| [a.c](x.md) | Final |\n| [abc](y.md) | Draft |\n"
		if got := readBack(t, path); got != want {
			t.Errorf("post-edit = %q; want %q — a bound value is a literal-match "+
				"fragment and can never alter the pattern's structure", got, want)
		}
	})

	t.Run("a_structural_metacharacter_value_cannot_widen_the_pattern", func(t *testing.T) {
		body := "| [0001](a.md) | Draft |\n| [0002](b.md) | Draft |\n"
		path := writeFixture(t, "readme.md", body)
		ino := inodeOf(t, path)

		// `.*` unquoted selects both rows. Quoted it selects neither, so
		// the refusal must be UNMATCHED, never a two-row rewrite.
		_, err := applyEdit(t,
			map[string]table.EditRule{editKey: editRule(
				`^\| \[{tag.nnnn}\].*$`, `{status}`, "")},
			path, map[string]string{"nnnn": ".*"},
			resolve.Tag{Key: editKey, Value: "Final"})

		assertRefusedUntouched(t, err, "edit_anchor_unmatched",
			editWriter+".edit."+editKey, path, body, ino)
	})
}

// REQ-111 (anchor-only half, S24): "the same value reaching an `anchor`
// only ... anchor-only is ADMITTED — `-` is ordinary regexp-quoted line
// data there, and asserting otherwise would fence out a legitimate anchor
// on a list-item line"
// DOMAIN EDGE
//
// C1.6's `-`-prefix rule is an ARGV rule, because the hazard is argument
// injection into a child process. There is no child on the edit path, so
// applying the rule to an anchor would fence out every anchor pinned to a
// markdown list item — the single most common line shape in the corpus.
func TestReq111_AFlagShapedTagValueIsAdmittedInAnAnchor(t *testing.T) {
	body := "- draft item\n- other\n"
	path := writeFixture(t, "list.md", body)

	if _, err := applyEdit(t,
		map[string]table.EditRule{editKey: editRule(`^{tag.nnnn} item$`, `{status} item`, "")},
		path, map[string]string{"nnnn": "- draft"},
		resolve.Tag{Key: editKey, Value: "- final"}); err != nil {
		t.Fatalf("a `-`-prefixed value in an ANCHOR refused: %v — the argv rule "+
			"is C1.6's and does not reach line data", err)
	}

	want := "- final item\n- other\n"
	if got := readBack(t, path); got != want {
		t.Errorf("post-edit = %q; want %q", got, want)
	}
}

// --- C1.5: `<clear>` on a line rule ---------------------------------------

// REQ-77: "clear = \"line\" ⇒ a planned `<clear>` deletes the anchored
// line, terminator included; read-back asserts the key ABSENT through the
// role's reader (0004:C11 unchanged)."
// HAPPY PATH
//
// S17's first arm. The terminator goes WITH the line: leaving it behind
// would produce a blank line where the bullet was, which is a different
// document.
func TestReq77_ClearLineDeletesTheAnchoredLineWithItsTerminator(t *testing.T) {
	path := writeFixture(t, "record.md", recordDoc)

	if _, err := applyEdit(t,
		map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "line")},
		path, nil,
		resolve.Tag{Key: editKey, Value: accessor.ClearSentinel}); err != nil {
		t.Fatalf("a declared clear refused: %v", err)
	}

	want := strings.Replace(recordDoc, "- **Status**: Draft\n", "", 1)
	if got := readBack(t, path); got != want {
		t.Errorf("post-clear:\n%q\nwant the line AND its terminator gone:\n%q", got, want)
	}
}

// REQ-78: "An anchor matching ZERO lines on a `<clear>` plan is SUCCESS
// with no write — \"clearing a key the artifact does not hold MUST
// succeed\" (0004:C11); ≥2 matches still refuses `edit_anchor_ambiguous`"
// BOUNDARY
//
// S17's second arm. The zero-match asymmetry is 0004:C11's, not this
// clause's invention: clearing a key the artifact does not hold is the
// caller getting what it asked for. The `≥2` half is asserted alongside so
// the asymmetry is not read as "clear never refuses".
func TestReq78_AZeroMatchClearSucceedsWithNoWriteButTwoMatchesStillRefuses(t *testing.T) {
	t.Run("zero_matches_succeeds_and_writes_nothing", func(t *testing.T) {
		path := writeFixture(t, "record.md", recordDoc)
		ino := inodeOf(t, path)

		if _, err := applyEdit(t,
			map[string]table.EditRule{editKey: editRule(`^- \*\*Absent\*\*:.*$`, `{status}`, "line")},
			path, nil,
			resolve.Tag{Key: editKey, Value: accessor.ClearSentinel}); err != nil {
			t.Fatalf("a zero-match clear refused: %v — clearing a key the "+
				"artifact does not hold MUST succeed", err)
		}
		if got := readBack(t, path); got != recordDoc {
			t.Errorf("the file changed on a zero-match clear:\n%q", got)
		}
		if got := inodeOf(t, path); got != ino {
			t.Errorf("inode moved from %d to %d; a zero-match clear is SUCCESS "+
				"with NO WRITE", ino, got)
		}
	})

	t.Run("two_matches_still_refuses_ambiguous", func(t *testing.T) {
		dup := recordDoc + "- **Status**: Draft\n"
		path := writeFixture(t, "dup.md", dup)
		ino := inodeOf(t, path)

		_, err := applyEdit(t,
			map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "line")},
			path, nil,
			resolve.Tag{Key: editKey, Value: accessor.ClearSentinel})

		assertRefusedUntouched(t, err, "edit_anchor_ambiguous",
			editWriter+".edit."+editKey, path, dup, ino)
	})
}

// REQ-79: "clear absent ⇒ a planned `<clear>` refuses BEFORE mutation,
// Detail `edit_clear_undeclared`, decided on the RAW planned value ahead
// of selection (C1.3 `precedence:`) — so an unmatched anchor on such a
// rule reports the undeclared `<clear>`, never `edit_anchor_unmatched`."
// REQ-53: "(3) `edit_clear_undeclared` (C1.5) — ... decided on the RAW
// planned value, BEFORE selection and before any `replace` expansion, so
// the refusal never depends on whether that rule's anchor matched"
// ADVERSARIAL
//
// The second sub-test is the load-bearing one: the SAME undeclared-clear
// plan against an anchor that matches ZERO lines must still report
// `edit_clear_undeclared`. An implementation that ran selection first would
// report `edit_anchor_unmatched` and send the author chasing a stale anchor
// instead of a missing `clear` declaration.
func TestReq53And79_AnUndeclaredClearRefusesAheadOfSelection(t *testing.T) {
	t.Run("with_a_matching_anchor", func(t *testing.T) {
		path := writeFixture(t, "record.md", recordDoc)
		ino := inodeOf(t, path)

		_, err := applyEdit(t,
			map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "")},
			path, nil,
			resolve.Tag{Key: editKey, Value: accessor.ClearSentinel})

		assertRefusedUntouched(t, err, "edit_clear_undeclared",
			editWriter+".edit."+editKey, path, recordDoc, ino)
	})

	t.Run("with_an_anchor_matching_zero_lines", func(t *testing.T) {
		path := writeFixture(t, "record.md", recordDoc)
		ino := inodeOf(t, path)

		_, err := applyEdit(t,
			map[string]table.EditRule{editKey: editRule(`^- \*\*Absent\*\*:.*$`, `{status}`, "")},
			path, nil,
			resolve.Tag{Key: editKey, Value: accessor.ClearSentinel})

		detail := detailOf(t, err)
		if strings.Contains(detail, "edit_anchor_unmatched") {
			t.Errorf("Detail = %q; the undeclared `<clear>` is decided on the RAW "+
				"planned value BEFORE selection, so it never depends on whether "+
				"the anchor matched", detail)
		}
		assertRefusedUntouched(t, err, "edit_clear_undeclared",
			editWriter+".edit."+editKey, path, recordDoc, ino)
	})
}

// REQ-80: "one-way: a deleted line cannot be re-established by `edit` (no
// append, C1.3), so after a `clear` the next non-clear write refuses
// `edit_anchor_unmatched`"
// REQ-117 (S29): "Two invocations against one artifact ... **Expected**:
// the first deletes the line and read-back reports the key ABSENT; the
// second refuses `edit_anchor_unmatched`"
// DOMAIN EDGE
//
// The only clause in C1 whose behaviour spans TWO invocations, and the
// test that stops a caller from treating `clear` as reversible. A single
// invocation cannot witness it.
func TestReq80And117_ClearIsOneWayAndTheNextWriteRefusesUnmatched(t *testing.T) {
	path := writeFixture(t, "record.md", recordDoc)
	rules := map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "line")}

	// Invocation 1 — the clear lands.
	if _, err := applyEdit(t, rules, path, nil,
		resolve.Tag{Key: editKey, Value: accessor.ClearSentinel}); err != nil {
		t.Fatalf("the clear refused: %v", err)
	}
	cleared := strings.Replace(recordDoc, "- **Status**: Draft\n", "", 1)
	if got := readBack(t, path); got != cleared {
		t.Fatalf("post-clear:\n%q\nwant:\n%q", got, cleared)
	}
	ino := inodeOf(t, path)

	// Invocation 2 — a non-clear write through the SAME rule. `edit`
	// never appends, so the line cannot be re-established.
	_, err := applyEdit(t, rules, path, nil,
		resolve.Tag{Key: editKey, Value: "Final"})

	assertRefusedUntouched(t, err, "edit_anchor_unmatched",
		editWriter+".edit."+editKey, path, cleared, ino)
}

// REQ-82: "the literal string `<clear>` is never substituted into
// `replace` (0004:C11) ... The sentinel is a property of the WHOLE planned
// value, tested before expansion: a `<clear>` authored into a `replace`
// template's LITERAL segment is ordinary bytes and emits as such, since
// only the planned value — never a template segment — is read as the
// sentinel"
// ADVERSARIAL
//
// S18. Both directions in one test, because the clause makes a precise
// distinction a one-sided test would blur: the sentinel is read on the
// PLANNED VALUE and nowhere else. A template literal spelling `<clear>` is
// bytes; a planned value spelling it is a removal.
func TestReq82_TheClearSentinelIsReadOnThePlannedValueAndNeverOnATemplate(t *testing.T) {
	t.Run("a_clear_in_a_template_literal_emits_as_bytes", func(t *testing.T) {
		path := writeFixture(t, "record.md", recordDoc)

		if _, err := applyEdit(t,
			map[string]table.EditRule{editKey: editRule(statusAnchor, `- **Status**: <clear> {status}`, "")},
			path, nil, resolve.Tag{Key: editKey, Value: "Final"}); err != nil {
			t.Fatalf("apply refused: %v", err)
		}

		want := strings.Replace(recordDoc,
			"- **Status**: Draft", "- **Status**: <clear> Final", 1)
		if got := readBack(t, path); got != want {
			t.Errorf("post-edit:\n%q\nwant the literal segment emitted as bytes:\n%q",
				got, want)
		}
	})

	t.Run("the_sentinel_is_never_substituted_as_a_value", func(t *testing.T) {
		path := writeFixture(t, "record.md", recordDoc)

		// A planned `<clear>` on a rule that DECLARES clear deletes; it
		// must never write the literal, which a re-read would report as a
		// mismatch (0004:C11).
		if _, err := applyEdit(t,
			map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "line")},
			path, nil,
			resolve.Tag{Key: editKey, Value: accessor.ClearSentinel}); err != nil {
			t.Fatalf("the clear refused: %v", err)
		}
		if got := readBack(t, path); strings.Contains(got, accessor.ClearSentinel) {
			t.Errorf("the file carries the literal sentinel after a clear:\n%q — "+
				"a re-read holding the literal is a mismatch (0004:C11)", got)
		}
	})
}

// --- C1.3 `order:` and `precedence:` --------------------------------------

// REQ-40: "order:    every refusal in this clause and C1.2's
// `edit_value_multiline` is decided BEFORE any byte is written. A refused
// edit is NOT APPLIED and surfaces as 0004's `execution_failure`"
// REQ-41: "a rule-scoped refusal carries a Detail naming the rule
// (`<id>.edit.<key>`) and the reason token (A8)"
// REQ-42: "no new refusal class is introduced, and no pre-write refusal
// THIS BINDING MINTS ever carries 0004:C14's applied-but-unverified sense"
// req-list A-11: the accessor-level class remains `execution_failure`.
// BOUNDARY
//
// S27's accessor-seam half. Three things at once: the CLASS is 0004's
// unchanged `execution_failure` (no new class, REQ-42), `Applied()` is
// FALSE (nothing was written), and the Detail carries both the rule id and
// the reason token.
func TestReq40And41And42_ARefusedEditIsExecutionFailureNotAppliedWithARuleScopedDetail(t *testing.T) {
	path := writeFixture(t, "record.md", recordDoc)

	acc := editAccessor(map[string]table.EditRule{
		editKey: editRule(`^- \*\*Absent\*\*:.*$`, `{status}`, ""),
	})
	reg := accessor.Registry{
		Flow:      editFlow,
		OwnedTags: []string{editKey},
		Definitions: []accessor.Definition{
			{
				Identity: accessor.Identity{Flow: editFlow, Name: editWriter, Capability: accessor.CapWrite},
				Accessor: acc,
				Binding:  flowbind.NewEditWriter(acc, editWriter),
			},
			{
				Identity: accessor.Identity{Flow: editFlow, Name: editReader, Capability: accessor.CapRead},
				Accessor: table.Accessor{Role: editRole, Keys: []string{editKey}, Timeout: "2s"},
				Binding:  &editStubReader{values: map[string]string{editKey: "Draft"}},
			},
		},
	}
	ex := accessor.NewExecutor(reg, accessor.Artifacts{editRole: {Role: editRole, Path: path}})

	res := ex.Write(context.Background(), editWriter,
		resolve.Plan{Writes: []resolve.Tag{{Key: editKey, Value: "Final"}}})

	if res.Refusal == nil {
		t.Fatalf("an unmatched anchor verified clean")
	}
	if got := res.Refusal.Class; got != accessor.ClassExecutionFailure {
		t.Errorf("class = %q; want 0004's UNCHANGED %q — no new refusal class "+
			"is introduced", got, accessor.ClassExecutionFailure)
	}
	if res.Refusal.Applied() {
		t.Errorf("Applied() = true on a PRE-WRITE refusal; no refusal this " +
			"binding mints carries the applied-but-unverified sense")
	}
	detail := res.Refusal.Detail
	if !strings.Contains(detail, "edit_anchor_unmatched") {
		t.Errorf("Detail = %q; want the reason token", detail)
	}
	if !strings.Contains(detail, editWriter+".edit."+editKey) {
		t.Errorf("Detail = %q; want the rule id `<id>.edit.<key>`", detail)
	}
	if got := readBack(t, path); got != recordDoc {
		t.Errorf("the target changed on a refused write")
	}
}

// REQ-51: "precedence: apply-time refusals fail-fast within one entry in
// this order ... (1) the ENTRY-level preconditions"
// REQ-52: "(2) `edit_value_multiline` ..."
// REQ-53: "(3) `edit_clear_undeclared` ..."
// REQ-54: "(4) per-rule cardinality, `edit_anchor_unmatched` then
// `edit_anchor_ambiguous`, resolved for EVERY rule of the entry before (5)
// the cross-rule `edit_anchor_collision` sweep ... (6)
// `edit_anchor_unstable`, necessarily last"
// REQ-120 (S32): "the earlier category in C1.3 `precedence:` is the one
// reported, and only it"
// ADVERSARIAL
//
// S32. Apply-time fail-fast is observable ONLY through a multi-defect
// input — with one defect per fixture every evaluation order passes, which
// is why the single-defect fixtures above cannot witness it. Each row
// asserts the EARLIER token AND that the later one did not win.
func TestReq51To54And120_ApplyTimeRefusalsFailFastInPrecedenceOrder(t *testing.T) {
	for _, tc := range []struct {
		name  string
		rules map[string]table.EditRule
		tags  map[string]string
		plan  []resolve.Tag
		want  string
		notes string
	}{
		{
			name: "multiline_value_beats_unmatched_anchor",
			rules: map[string]table.EditRule{
				editKey: editRule(`^- \*\*Absent\*\*:.*$`, `{status}`, ""),
			},
			plan:  []resolve.Tag{{Key: editKey, Value: "Fi\nnal"}},
			want:  "edit_value_multiline",
			notes: "step (2) precedes step (4)",
		},
		{
			name: "undeclared_clear_beats_unmatched_anchor",
			rules: map[string]table.EditRule{
				editKey: editRule(`^- \*\*Absent\*\*:.*$`, `{status}`, ""),
			},
			plan:  []resolve.Tag{{Key: editKey, Value: accessor.ClearSentinel}},
			want:  "edit_clear_undeclared",
			notes: "step (3) precedes step (4)",
		},
		{
			name: "unmatched_cardinality_beats_the_collision_sweep",
			rules: map[string]table.EditRule{
				editKey: editRule(`^- \*\*Absent\*\*:.*$`, `{status}`, ""),
				"owner": editRule(statusAnchor, `- **Status**: {owner}`, ""),
				"third": editRule(`^- \*\*Status\*\*: Draft$`, `- **Status**: {third}`, ""),
			},
			plan: []resolve.Tag{
				{Key: editKey, Value: "Final"},
				{Key: "owner", Value: "a"},
				{Key: "third", Value: "b"},
			},
			want:  "edit_anchor_unmatched",
			notes: "every rule's cardinality resolves before the cross-rule sweep",
		},
		{
			name: "multiline_value_beats_undeclared_clear",
			rules: map[string]table.EditRule{
				editKey: editRule(statusAnchor, statusReplace, ""),
				"owner": editRule(`^- \*\*Owner\*\*: (.*)$`, `- **Owner**: {owner}`, ""),
			},
			plan: []resolve.Tag{
				{Key: editKey, Value: "Fi\nnal"},
				{Key: "owner", Value: accessor.ClearSentinel},
			},
			want:  "edit_value_multiline",
			notes: "step (2) precedes step (3)",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := writeFixture(t, "record.md", recordDoc)
			ino := inodeOf(t, path)

			_, err := applyEdit(t, tc.rules, path, tc.tags, tc.plan...)

			detail := detailOf(t, err)
			if !strings.Contains(detail, tc.want) {
				t.Errorf("Detail = %q; want the token %q — %s", detail, tc.want, tc.notes)
			}
			// "and only it": no later token may also be reported.
			for _, later := range []string{
				"edit_value_multiline", "edit_clear_undeclared",
				"edit_anchor_unmatched", "edit_anchor_ambiguous",
				"edit_anchor_collision", "edit_anchor_unstable",
			} {
				if later == tc.want {
					continue
				}
				if strings.Contains(detail, later) {
					t.Errorf("Detail = %q also carries %q; the earlier category "+
						"is reported and ONLY it", detail, later)
				}
			}
			if got := readBack(t, path); got != recordDoc {
				t.Errorf("the file is not byte-identical after a refusal")
			}
			if got := inodeOf(t, path); got != ino {
				t.Errorf("inode moved on a refusal")
			}
		})
	}
}

// REQ-55: "Within one step, siblings are map-ranged and inherit C1.4's
// rule unchanged — which of two equally-defective rules is named is
// unspecified and no test may assert it."
// NEGATIVE REQ.
// ADVERSARIAL
//
// The assertion is stability of the TOKEN across repeated applies, never
// of the named rule. Pinning the winning rule would assert a guarantee the
// contract explicitly declines to make; asserting nothing would let a
// randomised token through.
func TestReq55_TwoEquallyDefectiveSiblingsReportAStableTokenAndNoPinnedWinner(t *testing.T) {
	rules := map[string]table.EditRule{
		editKey: editRule(`^- \*\*NopeA\*\*:.*$`, `{status}`, ""),
		"owner": editRule(`^- \*\*NopeB\*\*:.*$`, `{owner}`, ""),
	}

	for i := range 32 {
		path := writeFixture(t, "record.md", recordDoc)
		_, err := applyEdit(t, rules, path, nil,
			resolve.Tag{Key: editKey, Value: "Final"},
			resolve.Tag{Key: "owner", Value: "x"})

		detail := detailOf(t, err)
		if !strings.Contains(detail, "edit_anchor_unmatched") {
			t.Fatalf("apply %d: Detail = %q; want the token %q on every "+
				"iteration — the WINNING RULE is unspecified but the TOKEN is not",
				i, detail, "edit_anchor_unmatched")
		}
	}
}

// REQ-115 (S27b): "A bound target that cannot be read — absent (ENOENT), a
// directory (EISDIR), and unreadable (EACCES). **Expected**: all three
// refuse BEFORE mutation with `execution_failure` carrying the OS error in
// the Detail"
// REQ-32: "A target that cannot be read refuses BEFORE mutation with the
// OS error in the Detail — it does NOT reuse `edit_anchor_unmatched`, and
// it does NOT take `flowbind.go::load`'s absent-is-empty precedent"
// ADVERSARIAL
//
// The one refusal in C1 with NO `edit_*` reason token, so this also pins
// the Detail's shape where there is no rule to name. The absent-is-empty
// mutant is the dangerous one: `flowbind.go::load` treats a missing file
// as empty so a flow artifact can be CREATED on first write, and `edit`
// fences creation out — inheriting that precedent would launder a stale
// binding into `edit_anchor_unmatched`.
func TestReq32And115_AnUnreadableTargetRefusesWithTheOSErrorNotAnAnchorResult(t *testing.T) {
	dir := t.TempDir()

	for _, tc := range []struct {
		name  string
		setup func(t *testing.T) string
	}{
		{
			name: "absent_enoent",
			setup: func(t *testing.T) string {
				return filepath.Join(dir, "does-not-exist.md")
			},
		},
		{
			name: "a_directory_eisdir",
			setup: func(t *testing.T) string {
				p := filepath.Join(dir, "adir")
				if err := os.Mkdir(p, 0o755); err != nil {
					t.Fatalf("mkdir: %v", err)
				}
				return p
			},
		},
		{
			name: "unreadable_eacces",
			setup: func(t *testing.T) string {
				p := filepath.Join(dir, "locked.md")
				if err := os.WriteFile(p, []byte(recordDoc), 0o000); err != nil {
					t.Fatalf("write: %v", err)
				}
				if os.Geteuid() == 0 {
					t.Skip("running as root; mode 0000 is still readable")
				}
				return p
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := tc.setup(t)

			_, err := applyEdit(t,
				map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "")},
				path, nil, resolve.Tag{Key: editKey, Value: "Final"})
			if err == nil {
				t.Fatalf("an unreadable target applied cleanly")
			}

			detail := detailOf(t, err)
			if strings.Contains(detail, "edit_anchor_unmatched") {
				t.Errorf("Detail = %q; an unreadable target must NOT be laundered "+
					"into an anchor result", detail)
			}
			// The OS error rides the Detail; there is no rule to name.
			if detail == "" {
				t.Errorf("Detail is empty; the OS error must reach it")
			}
		})
	}
}

// --- C1.3 `no subprocess:` and `Invocations()` ----------------------------

// REQ-58: "no subprocess: `edit` spawns nothing; `--allow-commands`
// (0025:C6) is not consulted by the write itself. `Invocations()` counts
// `Apply` calls exactly as `flowbind.go::Writer` does (0004:C14)"
// REQ-110 (S23): "`Invocations()` after an `edit` apply, on an entry whose
// role's reader is FILE-backed, run with the gate OFF. **Expected**:
// counts `Apply` calls exactly as `flowbind.go::Writer` does; the write
// succeeds."
// DOMAIN EDGE
//
// GATE-OFF SUCCESS is the witness. `Invocations()` counts applies, not
// spawns, so it cannot prove the negative on its own — but a gate-off run
// that SUCCEEDS separates "not consulted" from "consulted and permitted".
// The reader is FILE-backed here, which is exactly what distinguishes this
// from the gate pre-check test below, where a gate-off run must REFUSE.
func TestReq58And110_AnEditAppliesWithTheGateOffAndInvocationsCountsApplies(t *testing.T) {
	path := writeFixture(t, "record.md", recordDoc)
	rules := map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "")}
	acc := editAccessor(rules)

	// The registry's gate is OFF and the role's reader is FILE-backed.
	writer := flowbind.NewEditWriter(acc, editWriter)
	reg := accessor.Registry{
		Flow:          editFlow,
		OwnedTags:     []string{editKey},
		AllowCommands: false,
		Definitions: []accessor.Definition{
			{
				Identity: accessor.Identity{Flow: editFlow, Name: editWriter, Capability: accessor.CapWrite},
				Accessor: acc,
				Binding:  writer,
			},
			{
				Identity: accessor.Identity{Flow: editFlow, Name: editReader, Capability: accessor.CapRead},
				Accessor: table.Accessor{
					Role: editRole, Path: path, Keys: []string{editKey}, Timeout: "2s",
				},
				Binding: &editStubReader{values: map[string]string{editKey: "Final"}},
			},
		},
	}
	ex := accessor.NewExecutor(reg, accessor.Artifacts{editRole: {Role: editRole, Path: path}})

	if got := writer.Invocations(); got != 0 {
		t.Fatalf("Invocations() = %d before any apply; want 0", got)
	}

	res := ex.Write(context.Background(), editWriter,
		resolve.Plan{Writes: []resolve.Tag{{Key: editKey, Value: "Final"}}})
	if res.Refusal != nil {
		t.Fatalf("a gate-off write with a FILE-backed reader refused: %+v — "+
			"`--allow-commands` is not consulted by the write itself", res.Refusal)
	}
	if got := writer.Invocations(); got != 1 {
		t.Errorf("Invocations() = %d after one write; want 1 — it counts `Apply` "+
			"calls exactly as `flowbind.go::Writer` does", got)
	}

	// The gate-off run must SUCCEED by actually writing, not by doing
	// nothing: "gate-off success is the witness for no subprocess spawned"
	// only holds if the edit landed.
	want := strings.Replace(recordDoc, "Draft", "Final", 1)
	if got := readBack(t, path); got != want {
		t.Errorf("the gate-off edit did not land:\n got %q\nwant %q", got, want)
	}
}

// REQ-59: "read-back: unchanged — 0004:C12/0004:C13 through the role's
// declared reader. When that reader is command-backed and the gate is off,
// the write MUST refuse BEFORE mutation (Detail naming the gate) ... a
// forgotten flag must not produce `read_back_incomplete` for a write that
// ran no command."
// REQ-60: "AUTHORITATIVE DETECTOR: this pre-check, not `Executor.Write`'s
// existing pre-`Apply` baseline read."
// REQ-112 (S25): "Run it BOTH ways over `protectedKeys`: once with the
// reader declaring only the planned key (`protected` empty ...) and once
// with the reader declaring a second key (`protected` non-empty ...).
// C1.3's pre-check must refuse identically in both"
// ADVERSARIAL
//
// S25. Running it BOTH ways over `protectedKeys` is what makes the
// pre-check the AUTHORITATIVE detector rather than a duplicate of the
// existing baseline arm: the baseline read only runs when `protected` is
// non-empty, and it SWALLOWS its failure into `baselineUnread` and
// proceeds. If the pre-check were absent, the `protected`-empty case would
// mutate and then report `read_back_incomplete` — the applied-but-
// unverified sense for a write that ran no command.
func TestReq59And60And112_AGateOffCommandReaderRefusesTheWriteBeforeMutation(t *testing.T) {
	for _, tc := range []struct {
		name        string
		readerKeys  []string
		description string
	}{
		{
			name:        "protected_empty_baseline_read_never_runs",
			readerKeys:  []string{editKey},
			description: "the reader declares only the planned key",
		},
		{
			name:        "protected_non_empty_baseline_read_runs_and_is_swallowed",
			readerKeys:  []string{editKey, "owner"},
			description: "the reader declares a second key",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := writeFixture(t, "record.md", recordDoc)
			ino := inodeOf(t, path)

			rules := map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "")}
			acc := editAccessor(rules)

			reg := accessor.Registry{
				Flow:          editFlow,
				OwnedTags:     []string{editKey},
				AllowCommands: false, // the gate is OFF
				Definitions: []accessor.Definition{
					{
						Identity: accessor.Identity{Flow: editFlow, Name: editWriter, Capability: accessor.CapWrite},
						Accessor: acc,
						Binding:  flowbind.NewEditWriter(acc, editWriter),
					},
					{
						// COMMAND-backed reader.
						Identity: accessor.Identity{Flow: editFlow, Name: editReader, Capability: accessor.CapRead},
						Accessor: table.Accessor{
							Role:    editRole,
							Command: []string{"rdr", "status", "-json", "{artifact}"},
							Keys:    tc.readerKeys,
							Timeout: "2s",
						},
						Binding: &editStubReader{values: map[string]string{editKey: "Final"}},
					},
				},
			}
			ex := accessor.NewExecutor(reg, accessor.Artifacts{editRole: {Role: editRole, Path: path}})

			res := ex.Write(context.Background(), editWriter,
				resolve.Plan{Writes: []resolve.Tag{{Key: editKey, Value: "Final"}}})

			if res.Refusal == nil {
				t.Fatalf("%s: a gate-off command-backed reader did not refuse the write",
					tc.description)
			}
			if got := res.Refusal.Class; got == accessor.ClassReadBackIncomplete {
				t.Errorf("%s: class = %q; a forgotten flag must NOT produce "+
					"read_back_incomplete for a write that ran no command",
					tc.description, got)
			}
			if got := res.Refusal.Class; got != accessor.ClassExecutionFailure {
				t.Errorf("%s: class = %q; want %q", tc.description, got,
					accessor.ClassExecutionFailure)
			}
			if res.Refusal.Applied() {
				t.Errorf("%s: Applied() = true on a refusal decided BEFORE "+
					"mutation", tc.description)
			}
			if got := readBack(t, path); got != recordDoc {
				t.Errorf("%s: the target changed; the refusal is BEFORE mutation",
					tc.description)
			}
			if got := inodeOf(t, path); got != ino {
				t.Errorf("%s: the inode moved; nothing should have been written",
					tc.description)
			}
			// The Detail names the GATE, having no rule to name.
			if d := res.Refusal.Detail; d == "" || !strings.Contains(d, "allow-commands") {
				t.Errorf("%s: Detail = %q; the entry-level precondition names the "+
					"gate instead of a rule", tc.description, d)
			}
		})
	}
}

// REQ-62: "`internal/accessor` does NOT import `cmdbind` (`cmdbind`
// imports `accessor`; the reverse is a cycle) and `AllowCommands` on
// `cmdbind.Config` stays where it is — the gate is carried to the accessor
// as state, never read across the seam (A10)"
// REQ-61: "SITE (A10, `Verified`): ... the gate from a field on
// `accessor.Registry`, set at `flowbind.go::Registry`"
// NEGATIVE REQ.
// DOMAIN EDGE
//
// The gate is STATE on the registry, not a call across the seam. The
// architectural constraint is load-bearing: reading the gate from
// `cmdbind` would be an import cycle, so an implementation that "just
// asked cmdbind" cannot compile — but one that duplicated the flag
// somewhere else would compile and drift. Asserting the field's presence
// on `Registry` pins the single site.
func TestReq61And62_TheGateIsCarriedAsStateOnTheAccessorRegistry(t *testing.T) {
	reg := accessor.Registry{Flow: editFlow, AllowCommands: true}
	if !reg.AllowCommands {
		t.Errorf("Registry.AllowCommands did not carry; the gate reaches the " +
			"accessor as STATE, never read across the cmdbind seam")
	}
}

// --- C1.1 in-memory residue -----------------------------------------------

// REQ-11: "in-memory: the registry's residue rule (0025:C1 runtime arm) is
// unchanged — an entry with no carrier builds a refusing binding, never a
// file binding"
// REQ-12: "`registry.go::commandBacked` becomes a three-way carrier
// discriminator: `command` → command binding; `edit` → edit binding;
// `path` → file binding; residue → refusing binding"
// REQ-118 (S30): "An entry with NO carrier at all, exercised at runtime
// rather than at load ... `commandBacked` is `len(acc.Command) != 0 ||
// acc.Path == \"\"` today, and its `acc.Path == \"\"` arm IS the residue
// path, so adding the `edit` arm could silently reroute residue to a file
// binding."
// ADVERSARIAL
//
// The dangerous mutant, stated by the scenario itself: the existing
// discriminator's `acc.Path == ""` arm IS the residue path, so an
// implementation that inserted the `edit` arm carelessly reroutes a
// carrier-less entry to a `Path: ""` FILE binding. That binding FAILS OPEN
// — it maps a non-existent file to an empty store — so every declared key
// reads ABSENT and an unapplied write verifies as clean. Silent state
// loss, never a crash.
func TestReq11And12And118_ACarrierLessEntryBuildsARefusingBindingNeverAFileBinding(t *testing.T) {
	m := &table.Model{
		ID: editFlow,
		Writers: map[string]table.Accessor{
			// No path, no command, no edit — the residue.
			"residue": {Role: editRole, Keys: []string{editKey}, Timeout: "2s", ReadBack: true},
		},
	}

	reg := flowbind.Registry(m, t.TempDir(), false)
	def, ok := reg.Lookup("residue", accessor.CapWrite)
	if !ok {
		t.Fatalf("the residue entry produced no definition")
	}

	wb, ok := def.Binding.(accessor.WriteBinding)
	if !ok {
		t.Fatalf("the residue binding is not a WriteBinding: %T", def.Binding)
	}

	path := writeFixture(t, "record.md", recordDoc)
	err := wb.Apply(context.Background(), accessor.Artifact{Role: editRole, Path: path},
		[]resolve.Tag{{Key: editKey, Value: "Final"}})
	if err == nil {
		t.Errorf("the residue binding APPLIED; a carrier-less entry must build a " +
			"REFUSING binding — a `Path: \"\"` file binding fails open and would " +
			"confirm an unapplied write as verified")
	}

	t.Run("an_edit_entry_gets_the_edit_binding_not_a_file_binding", func(t *testing.T) {
		em := &table.Model{
			ID: editFlow,
			Writers: map[string]table.Accessor{
				editWriter: editAccessor(map[string]table.EditRule{
					editKey: editRule(statusAnchor, statusReplace, ""),
				}),
			},
		}
		ereg := flowbind.Registry(em, t.TempDir(), false)
		edef, ok := ereg.Lookup(editWriter, accessor.CapWrite)
		if !ok {
			t.Fatalf("the edit entry produced no definition")
		}
		ewb, ok := edef.Binding.(accessor.WriteBinding)
		if !ok {
			t.Fatalf("the edit binding is not a WriteBinding: %T", edef.Binding)
		}

		p := writeFixture(t, "record.md", recordDoc)
		if err := ewb.Apply(context.Background(),
			accessor.Artifact{Role: editRole, Path: p},
			[]resolve.Tag{{Key: editKey, Value: "Final"}}); err != nil {
			t.Fatalf("the edit binding refused a well-anchored apply: %v", err)
		}
		want := strings.Replace(recordDoc, "Draft", "Final", 1)
		if got := readBack(t, p); got != want {
			t.Errorf("the `edit` carrier did not route to the edit binding:\n%q", got)
		}
	})
}

// REQ-13: "the accessor identity is 0004's `(flow, name, capability)`; the
// carrier does not enter it (0025 Identity LBD). A line rule's identity is
// `(entry, key)`."
// BOUNDARY
//
// GREEN-BY-DESIGN on the accessor half — 0004's identity triple is
// unchanged and this asserts the `edit` carrier did not widen it. The
// load-bearing half is the RULE's identity: `(entry, key)`, which is what
// makes a Detail of `<id>.edit.<key>` sufficient to name one rule
// unambiguously across an entry's siblings.
func TestReq13_TheCarrierDoesNotEnterTheAccessorIdentityAndARuleIsEntryPlusKey(t *testing.T) {
	acc := editAccessor(map[string]table.EditRule{
		editKey: editRule(statusAnchor, statusReplace, ""),
	})
	m := &table.Model{ID: editFlow, Writers: map[string]table.Accessor{editWriter: acc}}
	reg := flowbind.Registry(m, t.TempDir(), false)

	def, ok := reg.Lookup(editWriter, accessor.CapWrite)
	if !ok {
		t.Fatalf("no definition for the edit entry")
	}
	want := accessor.Identity{Flow: editFlow, Name: editWriter, Capability: accessor.CapWrite}
	if def.Identity != want {
		t.Errorf("identity = %+v; want %+v — the carrier does not enter it", def.Identity, want)
	}

	// The rule's identity is `(entry, key)`, witnessed by two sibling
	// rules of one entry being distinguishable in a refusal Detail.
	path := writeFixture(t, "record.md", recordDoc)
	_, err := applyEdit(t, map[string]table.EditRule{
		editKey: editRule(`^- \*\*NopeStatus\*\*:.*$`, `{status}`, ""),
	}, path, nil, resolve.Tag{Key: editKey, Value: "Final"})

	if detail := detailOf(t, err); !strings.Contains(detail, editWriter+".edit."+editKey) {
		t.Errorf("Detail = %q; a rule is named by `(entry, key)` as "+
			"`<id>.edit.<key>`", detail)
	}
}

// --- A1: the seam extension -----------------------------------------------

// REQ-95: "Context tags crossing the write seam | This RDR | Introduced |
// A1 — the one seam extension"
// REQ-96: "context tags must reach `Apply` for `{tag.<key>}` anchors,
// which today's `Artifact{Role, Path}` cannot carry ⇒ this RDR owns that
// field"
// REQ-97: "the SAME `art` value reaches both carriers, so one context map
// on `Artifact` serves the writer's anchor tags and C1.6's reader"
// REQ-98: "Context tags cross to `Apply` and `Read` (A1)"
// BOUNDARY
//
// ONE channel, not two. The scenario's phrase is "the SAME `art` value
// reaches both carriers": an implementation that added a separate
// `Write`-argument for the writer's tags would work for the anchor half
// and leave C1.6's reader unable to identify the row it read, which is the
// whole reason C1.6 exists.
func TestReq95To98_OneContextMapOnArtifactServesBothTheAnchorAndTheReader(t *testing.T) {
	body := "| [0028](0028-x.md) | Draft |\n| [0029](0029-y.md) | Draft |\n"
	path := writeFixture(t, "readme.md", body)

	art := accessor.Artifact{
		Role:    editRole,
		Path:    path,
		Context: map[string]string{"nnnn": "0028"},
	}

	// The writer's anchor consumes the context tag.
	acc := editAccessor(map[string]table.EditRule{
		editKey: editRule(
			`^\| \[{tag.nnnn}\]\([^)]*\) \| (\w+) \|$`,
			`| [{tag.nnnn}](0028-x.md) | {status} |`, ""),
	})
	w := flowbind.NewEditWriter(acc, editWriter)
	if err := w.Apply(context.Background(), art,
		[]resolve.Tag{{Key: editKey, Value: "Final"}}); err != nil {
		t.Fatalf("the anchor's context tag did not reach Apply: %v", err)
	}

	want := "| [0028](0028-x.md) | Final |\n| [0029](0029-y.md) | Draft |\n"
	if got := readBack(t, path); got != want {
		t.Errorf("post-edit = %q; want %q — only the addressed row flips", got, want)
	}

	// The SAME art value reaches a reader: one channel serves both.
	r := &editStubReader{values: map[string]string{editKey: "Final"}}
	if _, _, err := r.Read(context.Background(), art, []string{editKey}); err != nil {
		t.Fatalf("the same art did not reach Read: %v", err)
	}
	if got := r.sawContext["nnnn"]; got != "0028" {
		t.Errorf("the reader saw context %#v; want the same `nnnn=0028` the "+
			"writer's anchor consumed — one context map serves both carriers",
			r.sawContext)
	}
}

// --- C1.6: the argv placeholder at apply ----------------------------------

// REQ-88: "binding: the value comes from the invocation's context ...; an
// unbound tag at invocation refuses `execution_failure` BEFORE spawn,
// Detail naming the placeholder — a placeholder is never passed through
// literally (0025:C2)."
// REQ-89: "VALUE: a bound value beginning with `-` refuses
// `execution_failure` BEFORE spawn, Detail naming the placeholder,
// mirroring `cmdbind.go::substitute`'s existing `{artifact}` rule
// verbatim"
// REQ-90: "No `--` is inserted"
// ADVERSARIAL
//
// S24's argv arms. The class is argument injection (CWE-88), not shell
// injection — `os/exec` runs no shell — so the hazard is a flag-shaped
// word the CHILD reads as a flag. Both refusals are BEFORE spawn: a
// refusal after the process starts would have already run the command.
//
// The pass-through mutant is the one REQ-88 names: an unbound placeholder
// forwarded literally as `{tag.nnnn}` would reach the child as an argument
// and could be read as a filename.
func TestReq88And89And90_UnboundAndFlagShapedTagValuesRefuseBeforeSpawn(t *testing.T) {
	for _, tc := range []struct {
		name    string
		context map[string]string
		reason  string
	}{
		{
			name:    "unbound_tag",
			context: nil,
			reason:  "an unbound tag refuses before spawn, Detail naming the placeholder",
		},
		{
			name:    "flag_shaped_value_long",
			context: map[string]string{"nnnn": "--version"},
			reason:  "a bound value beginning with `-` refuses before spawn",
		},
		{
			name:    "flag_shaped_value_short",
			context: map[string]string{"nnnn": "-rf"},
			reason:  "a bound value beginning with `-` refuses before spawn",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// A command reader whose argv carries `{tag.nnnn}` at a
			// non-leading position. `false` is chosen as argv0 because it
			// exits non-zero: if the refusal did NOT happen before spawn,
			// the observable error would be the child's, not the
			// placeholder's.
			acc := table.Accessor{
				Role:    editRole,
				Command: []string{"false", "{tag.nnnn}"},
				Keys:    []string{editKey},
				Timeout: "2s",
			}
			r := cmdbind.Reader{Accessor: acc, Name: editReader, Config: cmdbind.Config{AllowCommands: true}}

			art := accessor.Artifact{Role: editRole, Path: "/dev/null", Context: tc.context}
			_, _, err := r.Read(context.Background(), art, []string{editKey})
			if err == nil {
				t.Fatalf("%s: the read succeeded", tc.reason)
			}
			detail := detailOf(t, err)
			if !strings.Contains(detail, "{tag.nnnn}") {
				t.Errorf("Detail = %q; want it to NAME the placeholder — %s",
					detail, tc.reason)
			}
		})
	}

	t.Run("no_double_dash_separator_is_inserted", func(t *testing.T) {
		// The admitted arm: a bound, non-flag value substitutes whole and
		// the argv is NOT rewritten with a `--` separator, because only a
		// child that honours the separator would be helped and the model
		// cannot know which do.
		acc := table.Accessor{
			Role:    editRole,
			Command: []string{"echo", "{tag.nnnn}"},
			Output:  ptr("raw"),
			Keys:    []string{editKey},
			Timeout: "2s",
		}
		r := cmdbind.Reader{Accessor: acc, Name: editReader, Config: cmdbind.Config{AllowCommands: true}}
		art := accessor.Artifact{
			Role: editRole, Path: "/dev/null",
			Context: map[string]string{"nnnn": "0028"},
		}
		vals, _, err := r.Read(context.Background(), art, []string{editKey})
		if err != nil {
			t.Fatalf("a bound non-flag value refused: %v", err)
		}
		for _, v := range vals {
			if strings.Contains(v.Value, "--") {
				t.Errorf("the child saw a `--` separator in %q; none is inserted", v.Value)
			}
		}
	})
}

// REQ-91: "The refusal is per-USE-SITE like C1.2's `value shape:` — a tag
// bound on the context but named by no argv element of this entry is never
// scanned."
// DOMAIN EDGE
//
// The mirror of the anchor-side per-use-site scope. A flag-shaped tag
// bound on the invocation but named by NO argv element of this entry
// reaches no child, so scanning it would refuse an invocation for a value
// the entry never passes.
func TestReq91_ABoundButUnreferencedFlagShapedTagIsNeverScannedAtArgv(t *testing.T) {
	acc := table.Accessor{
		Role:    editRole,
		Command: []string{"echo", "ok"},
		Output:  ptr("raw"),
		Keys:    []string{editKey},
		Timeout: "2s",
	}
	r := cmdbind.Reader{Accessor: acc, Name: editReader, Config: cmdbind.Config{AllowCommands: true}}
	art := accessor.Artifact{
		Role: editRole, Path: "/dev/null",
		Context: map[string]string{"unused": "--dangerous"},
	}

	if _, _, err := r.Read(context.Background(), art, []string{editKey}); err != nil {
		t.Errorf("a bound-but-unreferenced flag-shaped tag caused a refusal: %v — "+
			"the refusal is per-USE-SITE", err)
	}

	// POSITIVE CONTROL. The same flag-shaped value REFERENCED by an argv
	// element must refuse — otherwise "never scanned" passes vacuously on
	// an implementation that scans nothing at all.
	refAcc := table.Accessor{
		Role:    editRole,
		Command: []string{"echo", "{tag.unused}"},
		Output:  ptr("raw"),
		Keys:    []string{editKey},
		Timeout: "2s",
	}
	rr := cmdbind.Reader{
		Accessor: refAcc, Name: editReader,
		Config: cmdbind.Config{AllowCommands: true},
	}
	if _, _, err := rr.Read(context.Background(), art, []string{editKey}); err == nil {
		t.Errorf("the SAME flag-shaped value REFERENCED at argv did not refuse; " +
			"the per-use-site scope is a discrimination, not a blanket exemption")
	}
}

// REQ-92: "why here: 0004:C12 read-back re-reads the SAME role; a shared
// artifact (an index) has one reader for many records, and without an
// identity in its argv that reader cannot say which row it read. The
// consumer's own projector stays the reader (no line-oriented READ carrier
// is introduced — Briefly Rejected)"
// NEGATIVE REQ: no `edit` read carrier.
// REQ-93: "no other change to 0025:C1–C1.6: stdin envelopes, exit maps,
// env overlay and the gate are untouched."
// REQ-94: "The shell-interpreter deny-list stays STATIC"
// NEGATIVE REQ.
// DOMAIN EDGE
//
// `edit` is a WRITE carrier only. A read-side `edit` would make the reader
// a second line parser and split the read-back's authority between the
// consumer's projector and this binding — the Briefly Rejected shape.
func TestReq92And93_ThereIsNoEditReadCarrierAndTheReadSideIsUnchanged(t *testing.T) {
	m := &table.Model{
		ID: editFlow,
		Readers: map[string]table.Accessor{
			editReader: {
				Role:    editRole,
				Keys:    []string{editKey},
				Timeout: "2s",
				Edit: map[string]table.EditRule{
					editKey: editRule(statusAnchor, statusReplace, ""),
				},
			},
		},
	}
	reg := flowbind.Registry(m, t.TempDir(), false)
	def, ok := reg.Lookup(editReader, accessor.CapRead)
	if !ok {
		t.Fatalf("no definition for the read entry")
	}
	if _, isWrite := def.Binding.(accessor.WriteBinding); isWrite {
		t.Errorf("a read entry carrying `edit` produced a WriteBinding; `edit` is " +
			"a WRITE carrier only and no line-oriented READ carrier is introduced")
	}

	// The CONTROL that makes the negative above a real discrimination
	// rather than a blanket "nothing is ever an edit binding": the SAME
	// rules on a WRITE entry must route to the edit binding and apply.
	wm := &table.Model{
		ID: editFlow,
		Writers: map[string]table.Accessor{
			editWriter: editAccessor(map[string]table.EditRule{
				editKey: editRule(statusAnchor, statusReplace, ""),
			}),
		},
	}
	wreg := flowbind.Registry(wm, t.TempDir(), false)
	wdef, ok := wreg.Lookup(editWriter, accessor.CapWrite)
	if !ok {
		t.Fatalf("no definition for the write entry")
	}
	wb, ok := wdef.Binding.(accessor.WriteBinding)
	if !ok {
		t.Fatalf("the write entry's `edit` did not produce a WriteBinding: %T",
			wdef.Binding)
	}
	p := writeFixture(t, "record.md", recordDoc)
	if err := wb.Apply(context.Background(),
		accessor.Artifact{Role: editRole, Path: p},
		[]resolve.Tag{{Key: editKey, Value: "Final"}}); err != nil {
		t.Fatalf("the write-side edit binding refused: %v", err)
	}
	if got, want := readBack(t, p), strings.Replace(recordDoc, "Draft", "Final", 1); got != want {
		t.Errorf("the write-side edit did not land:\n%q", got)
	}
}

// --- helper binding -------------------------------------------------------

// editStubReader is a role reader the test controls. It is NOT a mock of
// the unit under test: the units here are the edit WRITE binding and the
// executor's write path, and 0004's own fixtures establish that the write
// binding and the read-back reader must be SEPARATE objects so the re-read
// can fail independently of the write.
type editStubReader struct {
	values     map[string]string
	sawContext map[string]string
}

func (r *editStubReader) Capability() accessor.Capability { return accessor.CapRead }

func (r *editStubReader) Read(
	_ context.Context, art accessor.Artifact, requested []string,
) ([]accessor.KeyValue, []string, error) {
	r.sawContext = art.Context

	var out []accessor.KeyValue
	for _, k := range requested {
		v, ok := r.values[k]
		out = append(out, accessor.KeyValue{Key: k, Value: v, Absent: !ok})
	}
	return out, nil, nil
}

func ptr[T any](v T) *T { return &v }

var _ = time.Second

// --- clauses whose whole content is a NEGATIVE -----------------------------

// REQ-21: "No word-splitting, option parsing, PATH resolution or shell
// exists on this path"
// NEGATIVE REQ.
// REQ-29: "out of scope for admission: a value that is well-formed line
// data but wrong for the DOCUMENT (a `|` inside a markdown table cell) is
// not a hazard this clause can name without document knowledge — the
// reader's read-back is the check (0004:C12), post-mutation by contract"
// NEGATIVE REQ: no per-rule value guard in v1.
// DOMAIN EDGE
//
// The clause's own scoping argument, asserted as behaviour. C1.2 refuses
// exactly ONE value shape — a newline or carriage return — because that is
// the only way line data becomes grammar here. Everything else a shell
// would mangle (spaces, quotes, globs, `$(…)`, a leading `-`) is ordinary
// line content, and a `|` inside a markdown cell is wrong for the DOCUMENT
// but not a hazard this clause can name without document knowledge.
//
// The mutant this kills is a defensive implementation that added a shell
// metacharacter deny-list "to be safe": it would refuse legitimate line
// content and, worse, imply a shell exists on a path that has none.
func TestReq21And29_OnlyNewlinesAreRefusedAndNoShellHazardIsScannedFor(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value string
	}{
		{name: "spaces_are_not_word_split", value: "Final  two   spaces"},
		{name: "a_leading_dash_is_line_data", value: "-rf"},
		{name: "shell_metacharacters_are_line_data", value: "a; b | c && d"},
		{name: "command_substitution_is_line_data", value: "$(rm -rf /)"},
		{name: "a_glob_is_line_data", value: "*.md"},
		{name: "quotes_are_line_data", value: `he said "hi" and 'bye'`},
		{name: "a_backtick_is_line_data", value: "`whoami`"},
		{name: "a_pipe_wrong_for_a_markdown_cell_is_still_admitted", value: "Fi|nal"},
		{name: "a_path_like_value_is_not_resolved", value: "/usr/bin/env"},
		{name: "a_tab_is_line_data", value: "Fi\tnal"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := writeFixture(t, "record.md", recordDoc)

			if _, err := applyEdit(t,
				map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "")},
				path, nil, resolve.Tag{Key: editKey, Value: tc.value}); err != nil {
				t.Fatalf("value %q refused: %v — the single admission hazard is a "+
					"newline or carriage return; there is no shell, no word-"+
					"splitting and no PATH resolution on this path", tc.value, err)
			}

			want := strings.Replace(recordDoc,
				"- **Status**: Draft", "- **Status**: "+tc.value, 1)
			if got := readBack(t, path); got != want {
				t.Errorf("post-edit:\n got %q\nwant %q — the value emits as literal "+
					"bytes", got, want)
			}
		})
	}
}

// REQ-34: "The reuse of `flowbind.go::save`'s discipline stops at
// stage-and-rename and does NOT drag along `Writer`'s unreachable-locator
// seal (`unreachable(path)`/`sealedKey`)" ... "`edit` therefore has no
// seal affordance, by construction rather than by omission"
// NEGATIVE REQ.
// REQ-100: "Extract `stageAndRename(dir, write, mode)` for both callers,
// or copy the discipline; do not extend `save`"
// NEGATIVE REQ: `flowbind.go::save` is not extended.
// ADVERSARIAL
//
// The seal keys off a declared-path SUFFIX an `edit` entry does not have,
// and its purpose is to make the applied-but-unverified sense atomic in a
// key/value artifact — a markdown target has nowhere to hold it. An
// implementation that reused `Writer` wholesale would emit seal bytes into
// the document, which is silent corruption of a reviewed artifact.
//
// The oracle is the FILE: after a successful edit, no seal key, no
// sentinel and no sidecar may appear anywhere in or beside the target.
func TestReq34And100_TheEditWriteEmitsNoSealAndCreatesNoSidecar(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "record.md")
	if err := os.WriteFile(path, []byte(recordDoc), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	if _, err := applyEdit(t,
		map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "")},
		path, nil, resolve.Tag{Key: editKey, Value: "Final"}); err != nil {
		t.Fatalf("apply refused: %v", err)
	}

	got := readBack(t, path)
	// The write must have LANDED, or "no seal bytes appear" is vacuously
	// true of a binding that wrote nothing.
	if want := strings.Replace(recordDoc, "Draft", "Final", 1); got != want {
		t.Fatalf("the edit did not land:\n%q", got)
	}
	for _, seal := range []string{"unreachable", "sealed", "sealedKey"} {
		if strings.Contains(got, seal) {
			t.Errorf("the target carries seal bytes %q:\n%s — `edit` has no seal "+
				"affordance, by construction rather than by omission", seal, got)
		}
	}

	// stage-and-rename leaves NOTHING behind: no temp file, no sidecar,
	// no directory `save`'s `MkdirAll(0o700)` would have created.
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("readdir: %v", err)
	}
	if len(entries) != 1 || entries[0].Name() != "record.md" {
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("the target's directory holds %#v; want only `record.md` — "+
			"stage-and-rename leaves no temp file and no sidecar", names)
	}
}

// REQ-46: "A rename that fails after a good staged write is also NOT
// APPLIED — the target is untouched by construction"
// REQ-48: "The atomicity boundary is ONE ENTRY over ONE file, and it is
// not transactional across entries" ... "a cross-entry transaction is not
// introduced"
// NEGATIVE REQ.
// BOUNDARY
//
// Two halves of the atomicity claim.
//
// WITHIN one entry: every rule of the entry rewrites ONE buffer and lands
// in ONE write, so a multi-rule entry moves the inode exactly once. An
// implementation applying rules one-at-a-time would move it per rule and
// leave the file observable in an intermediate state.
//
// ACROSS entries: entry 1 landing and entry 2 refusing leaves the two
// artifacts DISAGREEING, with the refusal reported for entry 2 only. That
// is 0004's per-entry apply model, and the record states outright that no
// cross-artifact rollback is introduced — asserting a rollback here would
// pin a guarantee the contract declines to make.
func TestReq46And48_AtomicityIsOneEntryOverOneFileAndNeverCrossEntry(t *testing.T) {
	t.Run("all_rules_of_one_entry_land_in_one_write", func(t *testing.T) {
		body := "- **A**: one\n- **B**: two\n- **C**: three\n"
		path := writeFixture(t, "multi.md", body)
		ino := inodeOf(t, path)

		if _, err := applyEdit(t, map[string]table.EditRule{
			editKey: editRule(`^- \*\*A\*\*: (.*)$`, `- **A**: {status}`, ""),
			"owner": editRule(`^- \*\*B\*\*: (.*)$`, `- **B**: {owner}`, ""),
			"third": editRule(`^- \*\*C\*\*: (.*)$`, `- **C**: {third}`, ""),
		}, path, nil,
			resolve.Tag{Key: editKey, Value: "1"},
			resolve.Tag{Key: "owner", Value: "2"},
			resolve.Tag{Key: "third", Value: "3"}); err != nil {
			t.Fatalf("a three-rule entry refused: %v", err)
		}

		want := "- **A**: 1\n- **B**: 2\n- **C**: 3\n"
		if got := readBack(t, path); got != want {
			t.Errorf("post-edit = %q; want %q", got, want)
		}
		// ONE write: the inode moved, and it moved once. A per-rule write
		// would have staged and renamed three times.
		if got := inodeOf(t, path); got == ino {
			t.Errorf("the inode did not move; the entry's rules must land in ONE write")
		}
	})

	t.Run("a_refusing_second_entry_does_not_roll_back_the_first", func(t *testing.T) {
		// Two entries over two files. Entry 1's edit is well-anchored;
		// entry 2's anchor matches nothing.
		one := writeFixture(t, "one.md", recordDoc)
		two := writeFixture(t, "two.md", recordDoc)

		if _, err := applyEdit(t,
			map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "")},
			one, nil, resolve.Tag{Key: editKey, Value: "Final"}); err != nil {
			t.Fatalf("entry 1 refused: %v", err)
		}
		if _, err := applyEdit(t,
			map[string]table.EditRule{editKey: editRule(`^- \*\*Absent\*\*:.*$`, `{status}`, "")},
			two, nil, resolve.Tag{Key: editKey, Value: "Final"}); err == nil {
			t.Fatalf("entry 2 applied; its anchor matches nothing")
		}

		// Entry 1 STAYS applied. The two artifacts disagree and the
		// refusal is reported for entry 2 only — the caller reconciles by
		// re-running, which the no-op arm makes safe.
		wantOne := strings.Replace(recordDoc, "Draft", "Final", 1)
		if got := readBack(t, one); got != wantOne {
			t.Errorf("entry 1's file was rolled back:\n%q\nwant\n%q — the "+
				"atomicity boundary is ONE ENTRY over ONE file and no "+
				"cross-entry transaction is introduced", got, wantOne)
		}
		if got := readBack(t, two); got != recordDoc {
			t.Errorf("entry 2's file changed on a refusal:\n%q", got)
		}
	})
}

// REQ-50: "No lock and no compare-before-rename: a concurrent writer is
// out of scope, as it is for `path`"
// NEGATIVE REQ.
// DOMAIN EDGE
//
// The clause declines a guarantee, and the assertion is that the decline
// is REAL: a second apply against a target whose bytes changed under the
// binding since it was read still lands. An implementation that added a
// compare-before-rename would refuse here and would be making a promise
// the contract explicitly does not make — and one callers would then rely
// on. The witness is that the concurrent writer's line survives (it is
// outside the anchored line) while the edit still applies.
func TestReq50_ThereIsNoLockAndNoCompareBeforeRename(t *testing.T) {
	path := writeFixture(t, "record.md", recordDoc)

	// A "concurrent writer" appends a line between the fixture being laid
	// down and the apply. No lock exists, so the apply proceeds.
	concurrent := recordDoc + "- **Added**: by someone else\n"
	if err := os.WriteFile(path, []byte(concurrent), 0o644); err != nil {
		t.Fatalf("concurrent write: %v", err)
	}

	if _, err := applyEdit(t,
		map[string]table.EditRule{editKey: editRule(statusAnchor, statusReplace, "")},
		path, nil, resolve.Tag{Key: editKey, Value: "Final"}); err != nil {
		t.Fatalf("apply refused: %v — no lock and no compare-before-rename "+
			"exists; a concurrent writer is out of scope, as it is for `path`", err)
	}

	want := strings.Replace(concurrent, "- **Status**: Draft", "- **Status**: Final", 1)
	if got := readBack(t, path); got != want {
		t.Errorf("post-edit:\n got %q\nwant %q", got, want)
	}
}

// REQ-94: "The shell-interpreter deny-list stays STATIC — the two rules
// above are rules on this placeholder family, not entries on that list and
// not a widening of it, which is what keeps 0027:C1's promise (\"argv
// WORDS only\", two named admitted forms) true with no third form to
// disclose (JDR 0003 §D2 (a); 0027 unchanged)"
// NEGATIVE REQ.
// DOMAIN EDGE
//
// C1.6's argv0 and `-`-prefix rules are rules on the `{tag.<key>}` FAMILY.
// They are not entries on 0027's interpreter deny-list and do not widen
// it. The witness is that the deny-list's membership is unchanged: a
// listed interpreter still refuses and an UNLISTED one still LOADS, both
// exactly as 0027 leaves them, with a `{tag.<key>}` element present.
func TestReq94_TheShellInterpreterDenyListIsNeitherWidenedNorNarrowed(t *testing.T) {
	// This is a load-time property, asserted here because it is C1.6's
	// claim about a SIBLING record's surface rather than about this
	// record's categories. `table.Categories()` must carry 0027's
	// category, unchanged and un-renamed, with no `{tag.…}` variant of it.
	registered := map[string]bool{}
	for _, c := range table.Categories() {
		registered[string(c)] = true
	}
	if !registered["command_shell_interpreter"] {
		t.Errorf("`command_shell_interpreter` is no longer registered; 0027:C1 " +
			"is unchanged by this record")
	}
	for _, minted := range []string{
		"edit_shell_interpreter",
		"command_tag_interpreter",
		"edit_tag_shell_interpreter",
	} {
		if registered[minted] {
			t.Errorf("this record minted %q; the deny-list stays STATIC and the "+
				"two `{tag.<key>}` rules are not entries on it", minted)
		}
	}
}

// REQ-99: "Data flow at apply: `Executor.Write` → non-owned check
// (unchanged) → resolve the role's reader (unchanged; C1.3's gate
// pre-check hangs here, A2) → `Apply`: read file → split lines
// (terminators kept) → for each rule: substitute `{tag.*}` into the anchor
// (quoted), compile, select exactly one line → for each rule: expand the
// parsed template segments (group text, planned value) → refuse or rewrite
// → stage + rename → executor read-back through the role's reader."
// BOUNDARY
//
// The phase ORDER, asserted where it is observable: the non-owned check
// runs BEFORE `Apply`, so a plan naming a non-owned tag must not reach the
// binding at all and the target must be untouched. `Invocations()` counts
// `Apply` calls, so a count of ZERO is the witness that the phase before
// it refused.
func TestReq99_TheNonOwnedCheckRunsBeforeApplyAndTheBindingIsNeverReached(t *testing.T) {
	path := writeFixture(t, "record.md", recordDoc)
	ino := inodeOf(t, path)

	acc := editAccessor(map[string]table.EditRule{
		editKey: editRule(statusAnchor, statusReplace, ""),
	})
	writer := flowbind.NewEditWriter(acc, editWriter)

	reg := accessor.Registry{
		Flow: editFlow,
		// `status` is the only OWNED tag; the plan below names `observed`.
		OwnedTags: []string{editKey},
		Definitions: []accessor.Definition{
			{
				Identity: accessor.Identity{Flow: editFlow, Name: editWriter, Capability: accessor.CapWrite},
				Accessor: acc,
				Binding:  writer,
			},
			{
				Identity: accessor.Identity{Flow: editFlow, Name: editReader, Capability: accessor.CapRead},
				Accessor: table.Accessor{Role: editRole, Path: path, Keys: []string{editKey}, Timeout: "2s"},
				Binding:  &editStubReader{values: map[string]string{editKey: "Draft"}},
			},
		},
	}
	ex := accessor.NewExecutor(reg, accessor.Artifacts{editRole: {Role: editRole, Path: path}})

	res := ex.Write(context.Background(), editWriter,
		resolve.Plan{Writes: []resolve.Tag{{Key: "observed", Value: "x"}}})

	if res.Refusal == nil {
		t.Fatalf("a plan naming a non-owned tag was applied")
	}
	if got := writer.Invocations(); got != 0 {
		t.Errorf("Invocations() = %d; want 0 — the non-owned check runs BEFORE "+
			"`Apply`, so the binding is never reached", got)
	}
	if got := readBack(t, path); got != recordDoc {
		t.Errorf("the target changed; the write must not reach the artifact at all")
	}
	if got := inodeOf(t, path); got != ino {
		t.Errorf("the inode moved on a refusal decided before `Apply`")
	}

	// POSITIVE CONTROL for the phase order: the SAME binding, reached
	// through the SAME executor with an OWNED plan, must apply and
	// increment `Invocations()`. Without it, "the binding is never
	// reached" passes on an implementation that never reaches it at all.
	res = ex.Write(context.Background(), editWriter,
		resolve.Plan{Writes: []resolve.Tag{{Key: editKey, Value: "Draft"}}})
	if res.Refusal != nil {
		t.Fatalf("the owned-plan control refused: %+v", res.Refusal)
	}
	if got := writer.Invocations(); got != 1 {
		t.Errorf("Invocations() = %d after the owned-plan control; want 1 — the "+
			"binding IS reached once the non-owned check passes", got)
	}
}
