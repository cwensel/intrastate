package flowbind

// RDR 0028 `0028:C1.3` — the `edit` write binding's execution semantics.
//
// This binding rewrites DECLARED LINES of a caller-bound text artifact.
// It spawns nothing, runs no shell, and performs no word-splitting,
// option parsing or PATH resolution: the only way runtime data could
// become grammar here is a newline inside an interpolated value, which
// C1.2 refuses before mutation, and that is the whole hazard class.
//
// Three properties shape the code more than anything else.
//
// FIRST — every refusal is decided BEFORE any byte is written. Selection,
// the cross-rule collision sweep and the post-rewrite re-anchor pass all
// run against in-memory buffers; the file is opened once for reading and,
// at most, once for the staged write. An implementation that wrote, then
// noticed, then reported would satisfy every error assertion and violate
// the clause.
//
// SECOND — bytes, not strings. A line is the bytes up to and excluding
// "\n", with a preceding "\r" part of the terminator. A bare "\r", NEL
// and U+2028 are line CONTENT. A missing final terminator is preserved,
// never supplied. `strings.Split` plus `strings.Join` gets each of these
// wrong in a different way, so the split keeps every terminator with its
// own line and the join is a concatenation.
//
// THIRD — the reuse of `save`'s discipline stops at stage-and-rename. It
// does NOT drag along `Writer`'s unreachable-locator seal, which keys off
// a declared-path suffix an `edit` entry does not have and would emit
// seal bytes into a reviewed document; and it does not extend `save`
// itself, whose `MkdirAll`, JSON-line store coupling and fixed
// `Chmod(0o600)` are each incompatible with a mode-preserving rewrite of
// an existing file.

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// The six apply-time reason tokens C1.3 and C1.5 mint. They ride a
// refusal's Detail and register NOWHERE in `table.Categories()`: they are
// DISJOINT from C1.4's six load-time wire strings, and registering one
// would claim that lint decides a stale anchor, which `what lint does NOT
// prove:` denies outright.
const (
	tokenAnchorUnmatched = "edit_anchor_unmatched"
	tokenAnchorAmbiguous = "edit_anchor_ambiguous"
	tokenAnchorCollision = "edit_anchor_collision"
	tokenAnchorUnstable  = "edit_anchor_unstable"
	tokenValueMultiline  = "edit_value_multiline"
	tokenClearUndeclared = "edit_clear_undeclared"
)

// EditWriter is the write binding for an `edit`-carried entry.
type EditWriter struct {
	acc  table.Accessor
	name string

	invocations int
}

// NewEditWriter builds the write binding for one `edit`-carried entry.
// The name is the entry id, which a rule-scoped refusal's Detail spells
// as `<id>.edit.<key>` — a line rule's identity is `(entry, key)`, and
// the accessor's own identity triple is unchanged by the carrier.
func NewEditWriter(acc table.Accessor, name string) accessor.WriteBinding {
	return &EditWriter{acc: acc, name: name}
}

// Capability reports write. `edit` is a WRITE carrier only: no
// line-oriented READ carrier is introduced, because the consumer's own
// projector stays the reader (`0028:C1.6` why here:).
func (w *EditWriter) Capability() accessor.Capability { return accessor.CapWrite }

// Invocations counts `Apply` calls exactly as `Writer` does, so a test
// can assert the accessor layer performed no retry or re-derivation
// (`0004:C14`). It counts APPLIES, never spawns — there are none.
func (w *EditWriter) Invocations() int { return w.invocations }

// editRulePlan is one rule carried through the apply pipeline.
type editRulePlan struct {
	key   string
	rule  table.EditRule
	value string
	clear bool

	anchorSegs  []table.EditSegment
	replaceSegs []table.EditSegment
	re          *regexp.Regexp

	// index is the rule's own PRE-EDIT line index, or -1 where the rule
	// selected nothing (a `<clear>` on a zero-match anchor, which is
	// success with no write).
	index int
	// deleted marks a rule that removed its line, so the re-anchor pass
	// can shift a sibling's held index by the deletions of the rules
	// preceding it. The shift is a COUNT, not a diff.
	deleted bool
}

// Apply performs the planned mutation over the caller-bound artifact.
//
// The refusal order is C1.3 `precedence:`, and it is fail-fast: (1) the
// entry-level preconditions, (2) `edit_value_multiline`, (3)
// `edit_clear_undeclared`, (4) per-rule cardinality for EVERY rule, (5)
// the cross-rule collision sweep, (6) the post-rewrite re-anchor pass.
// Steps 1–3 need no file at all, which is why the target is not opened
// until step 4 needs its content.
func (w *EditWriter) Apply(
	_ context.Context, art accessor.Artifact, planned []resolve.Tag,
) error {
	w.invocations++

	plans, err := w.prepare(art, planned)
	if err != nil {
		return err
	}

	// C1.3 `target:` — symlinks are RESOLVED before staging. Renaming
	// onto an unresolved symlink path would replace the link with a
	// regular file, which the clause forbids unconditionally.
	target := art.Path
	if resolved, rerr := filepath.EvalSymlinks(target); rerr == nil {
		target = resolved
	}

	// C1.3 `input:` — the whole file as BYTES. A target that cannot be
	// read refuses BEFORE mutation with the OS error in the Detail. It
	// does not take `load`'s absent-is-empty precedent: that precedent
	// exists so the FIRST write of a flow artifact can create it, and
	// `edit` fences creation out, so a missing target here is a stale
	// model rather than an anchor result.
	buf, err := os.ReadFile(target)
	if err != nil {
		return &accessor.ExecError{Detail: err.Error(), Err: err}
	}
	lines := splitLines(buf)

	if err := w.selectLines(plans, lines); err != nil {
		return err
	}
	if err := w.collisions(plans); err != nil {
		return err
	}

	out := w.rewrite(plans, lines)
	if err := w.reAnchor(plans, out); err != nil {
		return err
	}

	// C1.3 `write:` — a post-edit buffer equal to the input is not
	// written AT ALL: no staging, no rename. Stage-and-rename always
	// replaces the inode, so an unchanged inode is the observable that
	// proves the absence of a write, and file content cannot separate
	// "did not write" from "wrote identical bytes". Evaluated after the
	// re-anchor pass, which is the last refusal step.
	next := joinLines(out)
	if string(buf) == string(next) {
		return nil
	}

	mode := os.FileMode(0o644)
	if fi, serr := os.Stat(target); serr == nil {
		mode = fi.Mode().Perm()
	}
	if err := stageAndRename(target, next, mode); err != nil {
		return &accessor.ExecError{Detail: err.Error(), Err: err}
	}
	return nil
}

// prepare runs C1.3 `precedence:` steps (1) through (3) and parses each
// rule's two templates. None of it touches the artifact: every refusal
// here is decided from the declaration and the plan alone.
func (w *EditWriter) prepare(
	art accessor.Artifact, planned []resolve.Tag,
) ([]*editRulePlan, error) {
	keys := slices.Sorted(editKeys(w.acc.Edit))
	plans := make([]*editRulePlan, 0, len(keys))
	values := map[string]string{}
	for _, t := range planned {
		values[t.Key] = t.Value
	}

	for _, key := range keys {
		rule := w.acc.Edit[key]
		value, planned := values[key]
		if !planned {
			// A key this invocation did not plan has nothing to write.
			// Its rule contributes no selection, no deletion, and no
			// refusal: the plan's Writes are the only tags applied.
			continue
		}
		p := &editRulePlan{
			key:   key,
			rule:  rule,
			value: value,
			clear: accessor.IsClear(value),
			index: -1,
		}

		segs, err := table.ParseEditAnchor(rule.Anchor, nil)
		if err != nil {
			return nil, w.ruleErr(key, "edit_anchor_invalid", err.Error())
		}
		p.anchorSegs = segs

		// (1) The ENTRY-level preconditions. An unbound `{tag.<key>}` is
		// a property of the BINDING, not of a rule, so it names the
		// placeholder rather than a rule — a placeholder is never passed
		// through literally.
		pattern, missing, ok := table.ExpandEditAnchor(segs, art.Context)
		if !ok {
			return nil, &accessor.ExecError{
				Detail: "the placeholder `{tag." + missing +
					"}` is not bound on this invocation's context; a " +
					"placeholder is never passed through literally",
				// An entry-level precondition is about the REQUEST, not
				// the environment, so it takes the exit-2 group like the
				// rule-scoped refusals do (`0028:C1.3` EXIT GROUP:). A
				// missing `--tag` is the one defect re-running the same
				// request unchanged can never repair.
				Err: accessor.ErrDeclaredRequest,
			}
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, w.ruleErr(key, "edit_anchor_invalid", err.Error())
		}
		p.re = re

		rsegs, err := table.ParseEditReplace(rule.Replace, key, re.NumSubexp(), re.SubexpNames())
		if err != nil {
			return nil, w.ruleErr(key, "edit_template_invalid", err.Error())
		}
		p.replaceSegs = rsegs

		plans = append(plans, p)
	}

	// (2) `edit_value_multiline`, over the entry's planned values and the
	// tag values its rules ACTUALLY REFERENCE. The two halves have
	// different scopes and the clause is explicit about it: planned
	// values are scanned per-ENTRY, tag values per-USE-SITE. A tag bound
	// on the context but referenced by no anchor of this entry reaches no
	// line data, so scanning it would fence out a legitimate invocation
	// for a value the entry never touches.
	for _, p := range plans {
		if !p.clear && multiline(p.value) {
			return nil, w.entryErr(tokenValueMultiline,
				"the planned value for `"+p.key+"` carries a newline or "+
					"carriage return, which would turn one line into two")
		}
	}
	for _, p := range plans {
		for _, tag := range table.EditAnchorTagKeys(p.anchorSegs) {
			if multiline(art.Context[tag]) {
				return nil, w.entryErr(tokenValueMultiline,
					"the bound value of `{tag."+tag+"}`, referenced by the "+
						"anchor of `"+p.key+"`, carries a newline or carriage return")
			}
		}
	}

	// (3) `edit_clear_undeclared`, decided on the RAW planned value —
	// BEFORE selection and before any `replace` expansion — so the
	// refusal never depends on whether that rule's anchor matched. An
	// author whose anchor is also stale must be told about the missing
	// `clear` declaration, not sent chasing the anchor.
	for _, p := range plans {
		if p.clear && p.rule.Clear != table.EditClearLine {
			return nil, w.ruleErr(p.key, tokenClearUndeclared,
				"a `<clear>` was planned for a rule that declares no `clear` "+
					"disposition; a Status bullet has no meaningful cleared line "+
					"and the model author says so by omission")
		}
	}
	return plans, nil
}

// selectLines runs C1.3 `precedence:` step (4): per-rule cardinality,
// `edit_anchor_unmatched` then `edit_anchor_ambiguous`, resolved for
// EVERY rule of the entry before the cross-rule sweep can run.
//
// Selections are held as PRE-EDIT line indices against the pre-edit
// content, so a deletion never shifts a sibling rule's target.
func (w *EditWriter) selectLines(plans []*editRulePlan, lines []editLine) error {
	for _, p := range plans {
		var hits []int
		for i, ln := range lines {
			if p.re.Match(ln.content) {
				hits = append(hits, i)
			}
		}
		switch {
		case len(hits) == 0:
			// A `<clear>` on a zero-match anchor is SUCCESS with no write:
			// "clearing a key the artifact does not hold MUST succeed"
			// (`0004:C11`). It contributes no deletion to the re-anchor
			// shift, the shift being a count of lines actually removed.
			if p.clear {
				continue
			}
			// Never insert and never append: creation is fenced out, and
			// an unmatched anchor is a stale model, not a missing line.
			return w.ruleErr(p.key, tokenAnchorUnmatched,
				"the anchor selected no line of the bound artifact")
		case len(hits) > 1:
			return w.ruleErr(p.key, tokenAnchorAmbiguous,
				"the anchor selected more than one line of the bound artifact")
		default:
			p.index = hits[0]
		}
	}
	return nil
}

// collisions runs step (5): the cross-rule `edit_anchor_collision` sweep,
// decidable only once every rule holds a selection — which is why step
// (4) resolves every rule's cardinality first.
//
// It has no rule to name: the defect is a PAIR, so the Detail names both.
func (w *EditWriter) collisions(plans []*editRulePlan) error {
	seen := map[int]string{}
	for _, p := range plans {
		if p.index < 0 {
			continue
		}
		if other, taken := seen[p.index]; taken {
			return w.entryErr(tokenAnchorCollision,
				"the rules `"+w.name+".edit."+other+"` and `"+w.name+
					".edit."+p.key+"` both select the same line")
		}
		seen[p.index] = p.key
	}
	return nil
}

// rewrite applies every rule of the entry to ONE buffer, so they land in
// ONE write. A `<clear>` on a rule declaring `clear = "line"` removes the
// line WITH its terminator: leaving the terminator behind would produce a
// blank line where the bullet was, which is a different document.
func (w *EditWriter) rewrite(plans []*editRulePlan, lines []editLine) []editLine {
	out := slices.Clone(lines)
	drop := map[int]bool{}
	for _, p := range plans {
		if p.index < 0 {
			continue
		}
		if p.clear {
			drop[p.index] = true
			p.deleted = true
			continue
		}
		out[p.index] = editLine{
			content: []byte(expandReplace(p.replaceSegs, p.value, p.re, lines[p.index].content)),
			term:    lines[p.index].term,
		}
	}
	if len(drop) == 0 {
		return out
	}
	kept := make([]editLine, 0, len(out))
	for i, ln := range out {
		if !drop[i] {
			kept = append(kept, ln)
		}
	}
	// Deleting the FINAL line of a file that had no final terminator also
	// removes the preceding line's terminator, so the file's
	// final-terminator state is preserved either way.
	if len(kept) != 0 && len(out) != 0 && drop[len(out)-1] && out[len(out)-1].term == nil {
		kept[len(kept)-1].term = nil
	}
	return kept
}

// reAnchor runs step (6), necessarily last: every rule's anchor is run
// again over the POST-EDIT buffer and must select exactly its own
// rewritten line — or, for a deleted line, zero lines.
//
// IDENTITY, not cardinality. A pass that checked only "exactly one match"
// would admit a rule whose anchor drifted onto ANOTHER rule's line after
// a rewrite or a deletion: exactly one match, wrong line. That is how a
// replacement poisons a sibling's anchor on the next run, and it is the
// whole reason this pass exists rather than a re-run of `selectLines`.
func (w *EditWriter) reAnchor(plans []*editRulePlan, out []editLine) error {
	for _, p := range plans {
		var hits []int
		for i, ln := range out {
			if p.re.Match(ln.content) {
				hits = append(hits, i)
			}
		}
		if p.index < 0 || p.deleted {
			// A rule that deleted its line, or a `<clear>` that matched
			// none, must select zero lines afterwards.
			if len(hits) != 0 {
				return w.ruleErr(p.key, tokenAnchorUnstable,
					"the anchor still selects a line after its own line was removed")
			}
			continue
		}
		// The held pre-edit index, shifted by the deletions of the
		// PRECEDING sibling rules. The shift is a count, not a diff.
		want := p.index
		for _, q := range plans {
			if q.deleted && q.index >= 0 && q.index < p.index {
				want--
			}
		}
		if len(hits) != 1 || hits[0] != want {
			return w.ruleErr(p.key, tokenAnchorUnstable,
				"the anchor does not select exactly its own rewritten line on "+
					"the post-edit buffer; the replacement de-anchors the rule or "+
					"poisons a sibling's anchor")
		}
	}
	return nil
}

// --- refusal carriers -----------------------------------------------------

// ruleErr builds a RULE-SCOPED refusal: its Detail names the rule as
// `<id>.edit.<key>` AND the reason token, which together are what let a
// caller trace a refusal to one declaration among an entry's siblings
// (`0028:C1.3` order:).
func (w *EditWriter) ruleErr(key, token, detail string) error {
	return &accessor.ExecError{
		Detail: token + ": " + w.name + ".edit." + key + ": " + detail,
		// These refusals are about the REQUEST — a stale anchor, an
		// ambiguous one, a value carrying a newline — so they take the
		// exit-2 group rather than `execution_failure`'s default exit 3.
		// An OS error reading the target does NOT carry this: that one
		// really is the environment and a re-run may well succeed.
		Err: accessor.ErrDeclaredRequest,
	}
}

// entryErr builds an ENTRY-LEVEL refusal, which names no rule because it
// has none to name: the multiline scan spans the entry's whole plan and
// the collision sweep is a property of a PAIR.
func (w *EditWriter) entryErr(token, detail string) error {
	return &accessor.ExecError{
		Detail: token + ": " + w.name + ": " + detail,
		Err:    accessor.ErrDeclaredRequest,
	}
}

// --- line handling --------------------------------------------------------

// editLine is one line: its content bytes and its own terminator, kept
// SEPARATE so a rewrite replaces the content and leaves the terminator
// exactly as authored. CRLF is preserved per line and a missing final
// terminator is preserved rather than supplied.
type editLine struct {
	content []byte
	// term is the line's terminator — "\n" or "\r\n" — or nil for a final
	// line the file did not terminate.
	term []byte
}

// splitLines splits bytes into lines on "\n", with a preceding "\r"
// treated as part of the terminator. Only "\n" terminates: a bare "\r",
// NEL and U+2028 are line CONTENT, so a run carrying all three is ONE
// line.
func splitLines(buf []byte) []editLine {
	var out []editLine
	start := 0
	for i := range len(buf) {
		if buf[i] != '\n' {
			continue
		}
		end := i
		term := buf[i : i+1]
		if end > start && buf[end-1] == '\r' {
			end--
			term = buf[end : i+1]
		}
		out = append(out, editLine{content: buf[start:end], term: term})
		start = i + 1
	}
	if start < len(buf) {
		out = append(out, editLine{content: buf[start:]})
	}
	return out
}

// joinLines concatenates lines and their own terminators. It is a
// concatenation and never a `strings.Join`, which would supply a
// separator the file may not have carried.
func joinLines(lines []editLine) []byte {
	var b []byte
	for _, ln := range lines {
		b = append(b, ln.content...)
		b = append(b, ln.term...)
	}
	return b
}

// multiline reports the ONE value shape C1.2 refuses.
func multiline(v string) bool {
	return strings.ContainsAny(v, "\n\r")
}

// expandReplace emits the parsed template's segments as bytes. Captured
// text and the substituted value are emitted LITERALLY and never
// re-scanned for `${…}` or `{…}`: if they were, a document line
// containing `${1}` would let the DOCUMENT choose what gets written —
// runtime data becoming grammar, which is exactly what C1.2 forbids.
// This is also why the expansion is not a call to `regexp.Expand`, which
// re-scans its template at expansion time.
func expandReplace(
	segs []table.EditSegment, value string, re *regexp.Regexp, line []byte,
) string {
	match := re.FindSubmatchIndex(line)
	names := re.SubexpNames()

	var b strings.Builder
	for _, s := range segs {
		switch s.Kind {
		case table.EditKey:
			b.WriteString(value)
		case table.EditGroup:
			n := s.Index
			if n < 0 {
				for i, nm := range names {
					if nm == s.Text {
						n = i
						break
					}
				}
			}
			// A group that did not participate in the match expands to
			// the empty string. The rule matches `regexp.Expand`'s, but
			// the parse is this clause's own.
			if n >= 0 && match != nil && 2*n+1 < len(match) &&
				match[2*n] >= 0 && match[2*n+1] >= 0 {
				b.Write(line[match[2*n]:match[2*n+1]])
			}
		default:
			b.WriteString(s.Text)
		}
	}
	return b.String()
}

// --- the staged write -----------------------------------------------------

// stageAndRename stages the replacement BESIDE the resolved target and
// renames it over, so a reader never observes a half-written document and
// a rename that fails leaves the target untouched by construction.
//
// It is a separate helper rather than an extension of `save`, whose
// `MkdirAll(0o700)` creation, `WriteJSONLine` store coupling and fixed
// `Chmod(0o600)` are each incompatible here: mode is PRESERVED from the
// existing target, or a 0755 document would come back 0644 and show up as
// a mode change in `git diff --summary`.
func stageAndRename(target string, content []byte, mode os.FileMode) error {
	dir := filepath.Dir(target)
	tmp, err := os.CreateTemp(dir, ".flow-edit-*")
	if err != nil {
		return err
	}
	staged := tmp.Name()
	defer func() { _ = os.Remove(staged) }()

	if _, err := tmp.Write(content); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(staged, mode); err != nil {
		return err
	}
	return os.Rename(staged, target)
}

func editKeys[V any](m map[string]V) func(func(string) bool) {
	return func(yield func(string) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}
