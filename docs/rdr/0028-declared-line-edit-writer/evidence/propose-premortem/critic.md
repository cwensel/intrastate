Model: claude-fable-5

| ID | passage/claim | Failure mode | Symptom user sees | Origin |
|---|---|---|---|---|
| P-1 | Lint: "`anchor` carries no placeholder" + consumer README index table has one row per record | The README Status cell cannot be selected by a single declared rule: the row is identified by the record id, the id is per-artifact, and the anchor may not carry it. Either N entries (one per record, hand-edited: the drift returns) or the write cannot be declared. | Author declares the README write, lint passes, apply refuses "anchor matched 12 lines" (every row) or the author copy-pastes an entry per record | negation |
| P-2 | "two rules select the same line -> refuse" + "one rule per key" | A README index row holds Status AND Updated cells on ONE line; two owned keys must land on the same line; the filter refuses exactly the consumer's artifact | Apply refuses "rules status and updated select the same line 41" on every run | refutation (S3) |
| P-3 | Claim 1: the only structural hazard is a line terminator | A value containing `\|` inserted into a table cell shifts columns; the reader (projector) parses the wrong cell; write succeeds, read-back mismatches AFTER mutation | `read_back_mismatch`, README row corrupt in git, "may have been applied" | negation |
| P-4 | Claim 1: value "cannot change which line is selected" | True for this apply, false for the next: a value written on run N contains text that matches another rule's anchor; run N+1 refuses with 2 matches. The tool poisons its own anchors. | Second state change refuses "anchor matched 2 lines"; nothing in the model changed | hindsight |
| P-5 | "rewrites the selected lines" | Whole-line vs matched-span ambiguity: an unanchored regex (`Status\*\*: (\w+)`) matches a span; whole-line replace drops the `- **` prefix | Bullet becomes `Status**: Final`; projector no longer parses the bullet; `read_back_incomplete` post-mutation | negation |
| P-6 | Claim 3: template parsed at load into literal/backref/placeholder | No escape syntax for literal `{`, `}`, `$` in replace; markdown with Hugo/Jinja `{{ }}` or a `$5` literal either fails lint or is mis-segmented; unmatched optional group `${N}` is undefined | Lint "placeholder `{ site.x }` names no key"; or `$5` becomes empty | negation |
| P-7 | Claim 4: gate-off read-back refused BEFORE mutation | The gate check lives in the read-back path (post-write) today; a "widening" (edit needs no gate) becomes a narrowing: edit written, reader refused, artifact applied-but-unverified | File changed on disk, exit says refused, no read-back ran | refutation (S2) |
| P-8 | Claim 4: "runs no subprocess so needs no gate" | The gate's veto semantics are "a model cannot cause execution", not "cannot spawn now"; `edit` writes any path the model names: `.git/hooks/pre-commit`, `Makefile`, `.claude/settings.json` hooks line | A cloned repo's model edits a hook line without `--allow-commands`; execution deferred to the next commit | refutation (S1) |
| P-9 | Claim 5: temp+rename sufficient, no locking | Rename over a symlink replaces the link with a regular file; mode/xattrs/hard links lost; temp on another filesystem fails rename; no compare-on-write so an editor/agent save between read and rename is lost | README symlink severed; file mode 0600 becomes 0644; concurrent agent edit vanishes | negation |
| P-10 | Claim 5: line endings preserved exactly | Last line without trailing newline, mixed CRLF/LF, and `clear` deleting the last line: which terminator goes? Bare `\r` and U+2028/NEL are not treated as terminators by the newline check but the reader or git may | Diff shows a spurious final-newline change on every apply; or a value with U+2028 splits the line for the reader | negation |
| P-11 | Claim 6: `clear` deletes the anchored line | Deleting the README row deletes the whole record's index row (title, date, link), not the Status cell; and after deletion the anchor matches 0 lines, so the next set refuses (R5: no append). Clear is one-way and destructive. | Row gone from README; re-setting status refuses "anchor matched 0 lines" forever | hindsight |
| P-12 | Claim 6: "read-back asserts absence" borrowed from the JSON carrier | Reader defaults a missing Status bullet to "Draft" (projector fills defaults); absence of the line reads back as presence of a value | `read_back_mismatch` after a successful clear; or worse, equality with the old default and a false pass | refutation (S4) |
| P-13 | Claim 7: lint proves statically all that matters | Lint cannot prove: anchors disjoint, exactly-one, template re-emits a captured span containing the old value (`${0}`), post-edit line still matches its anchor (re-anchorability) | First-run refusal or a permanent `read_back_mismatch` that lint passed | negation |
| P-14 | Claim 2: exactly one match is correct for these artifacts | Records quote the template ("`- **Status**: Draft` appears in a fence"), peers are cited with their status bullet, a Status History table repeats the string; 2+ matches refuse a valid record | Apply refuses on a record that is fine; author "fixes" the record content to satisfy the tool | negation |
| P-15 | Claim 8: read-back through the role's reader verifies the edit | The README write is to a different artifact than the record the role's reader parses; the README key has no reader; it is either unverified or `read_back_incomplete` on every run | README write is never verified while the exit says verified; or every apply is incomplete | refutation (S4) |
| P-16 | Claim 8: byte-for-byte equality | Reader normalizes (lowercase, trim, markdown strip, date typing); "Final" reads back "final" | Permanent `read_back_mismatch` after a correct edit | negation |
| P-17 | Claim 9: every pre-write refusal is `execution_failure` = "not applied" | `execution_failure` for the `command` carrier already covers post-partial-write exits; the class's veto semantics are "unknown", not "not applied"; a failed rename (EXDEV, EACCES on dir) after a good temp write is also here | Agent loop retries an "unknown" as if not applied; or treats a not-applied as maybe-applied and stops | refutation (S1) |
| P-18 | Claim 10: anchors resolved against pre-edit content, one buffer | Line-index shift: `clear` deletes line 5, rule for line 9 now edits line 8 unless applied descending or on a pre-split slice | Wrong line rewritten in the same apply; read-back mismatch post-mutation | negation |
| P-19 | R2: sub-verb "same semantics behind a process boundary, no safety gain" | The gate IS the safety gain for arbitrary caller-owned file writes; in-process removes the only per-invocation consent for writing outside tool-owned files | See P-8 | negation |
| P-20 | R3: adapter registry "imports document knowledge" | P-11 shows `clear` needs shape knowledge (delete line vs blank cell); the registry returns as disposition enums per shape | Every new artifact shape is still a release, now as a `clear` variant | negation |
| P-21 | R5: "an unmatched anchor is a stale model, not a missing line" | After the tool's own `clear`, the unmatched anchor is a missing line the tool made; R5's premise is falsified by claim 6 | See P-11 | negation |
| P-22 | R1/R4: wrapper and emit-only "relocate drift" | The wrapper is in the repo and covered by the command carrier's argv lint; `edit`'s regex is equally opaque to lint about content; the drift argument cuts both ways | None visible; the rejection is rhetorical, not a defect | negation |

---

## 1. Failure narrative (prospective hindsight)

The `edit` carrier shipped in the release after 0028 locked. The consumer bound it to two rules in one write entry: the record's `- **Status**: {status}` bullet and the README index row's Status cell. Lint passed on the first try.

The first apply on a real record refused: "rules `status` and `updated` select the same line 41". The README row is one line and carries both the Status and the Updated cells; the design's one-rule-per-key and same-line-refuse rules had never been tried against the README. The author split the README write into a second entry with only `status`, deferring `updated`. The second apply refused: "anchor matched 14 lines". The anchor for the README row could not name the record id because lint forbids placeholders in `anchor`, so the author had pasted an entry per record, then a colleague's record got added and the anchor written for 0027 matched 0027 and 0027's "superseded-by" back-reference row. The author narrowed the regex by hand.

Three weeks in, a record's Status was set to `Final` and its README cell to `Final | 2026-09-12` because the table had been extended with a merged cell and the planned value now carried a pipe. No line terminator, so the pre-write check passed. The rename landed. The projector read the row with one extra column and reported the Status cell as `2026-09-12`. `read_back_mismatch` fired, post-mutation. The README was committed by the agent loop that treats mismatch as "may have been applied, re-run"; the re-run refused with "anchor matched 0 lines" because the rewritten row no longer matched the original anchor. The model was now permanently stuck on that record, and nobody could tell from the refusal whether the file or the model was wrong.

Separately, a `<clear>` of status on a withdrawn record deleted its entire README row (title, link, date), which the projector treated as "record not indexed" and the docs site dropped the page. Re-setting the status refused: 0 matches, no append by design. The only recovery was a hand edit: exactly the edit no model declared and no lint covered.

The last incident was the quiet one. A model in a cloned repo declared an `edit` write against `.git/hooks/pre-commit`, matching the one `exec` line. No `--allow-commands` was needed because the write ran no subprocess. The hook ran at the next commit.

---

## 2. Obstacle negation

### Claim 1 (value interpolation into file data; only hazard is a line terminator)
Negated by P-3, P-4, P-10. The file is not opaque bytes to the consumer; it is markdown parsed by a projector. Interpolation into a markdown table cell has cell separators (`|`), into a bullet has emphasis and code delimiters, into a heading has `#`. None of these is a line terminator, all of them change what the reader reads back, and the check runs before the write so the damage is post-mutation. Cross-run: a value that contains the text of another rule's anchor makes the next apply ambiguous. Refusing terminators closes one class; it does not close "value changes what the reader reads" or "value changes what the next anchor selects". Terminator set: if the check is `\n` only, bare `\r`, U+2028, U+0085 pass the check and split the line for some readers and for git's diff.
Answer required: a per-rule `refuse_chars` or a mandatory post-substitution re-anchor check (the rewritten buffer must still select exactly the same line for every rule).

### Claim 2 (exactly one match)
Negated by P-1, P-14. Records quote the template, cite peers with their status bullet, and carry history tables; the consumer's own artifacts violate the premise. And the README requires an identity in the anchor that lint forbids. "Safer than lineinfile" is true only when the artifact has one candidate line; when it has two legitimate ones, strict single-match is a refusal on a correct record, which pushes authors to edit content to please the tool.
Answer required: a quoted, escaped identity placeholder permitted in `anchor` (regexp.QuoteMeta applied, never raw), and a documented recipe for scoping (anchor includes the section heading via a two-line window, or the rule carries a `section` pre-anchor).

### Claim 3 (segment-wise substitution)
Negated by P-6, P-13. Injection is closed; correctness is not. No escape for literal `{`, `}`, `$`; undefined `${N}` for an unmatched optional group; a template that re-emits `${0}` or a group spanning the old value silently writes the old value back. Whole-line vs span replace (P-5) is unspecified and lint cannot see it.
Answer required: an escape syntax, a lint rule that every `${N}` names a group index within the anchor's group count, and an apply-time invariant that the rewritten line matches the anchor and that the placeholder segment's bytes are the planned value.

### Claim 4 (no gate; pre-mutation refusal of gate-off read-back)
Negated by P-7, P-8. If the gate check is where it lives today (in the reader path, after the write), the pre-mutation refusal is a promise the control flow does not keep: S2's shape exactly. And "no subprocess" is the wrong partition: the gate is consent to let a model cause execution; a line edit to a hook, a Makefile, a CI file, or a settings file is deferred execution.
Answer required: move the gate evaluation into the plan phase (before any carrier runs) and either scope `edit` targets to the artifact's own path or a declared path allowlist, refusing symlinks and paths under `.git/`.

### Claim 5 (temp+rename; no locking; bytes preserved)
Negated by P-9, P-10, P-18. Rename replaces a symlink with a regular file; mode, xattrs, hard links are lost; temp on another mount fails at rename after a successful write (an execution_failure that is honestly "not applied", but for a different reason); no compare-on-write means a concurrent save is lost. Line endings: last-line-without-newline and `clear`-of-last-line have no stated rule.
Answer required: temp in the target's directory, `Lstat` refusal on symlinks, mode copy, a size+hash compare of the file immediately before rename against the bytes read, and explicit terminator rules for the last line.

### Claim 6 (`<clear>` deletes the line)
Negated by P-11, P-12, P-21. Deleting a table row deletes every other key on that line and the record's index presence; the line cannot come back (R5 forbids append), so clear is a one-way door. And absence-as-read-back was defined for a JSON carrier where "key absent" is unambiguous; a projector that defaults a missing bullet to `Draft` makes absence read back as a value.
Answer required: `clear` dispositions `delete-line` and `replace-with` (a template with no placeholder, e.g. blank the cell), and a stated read-back semantics for clear against readers that default.

### Claim 7 (lint proves everything relevant statically)
Negated by P-13. The correctness of an edit depends on file content by construction; the honest statement is "lint proves the model is well-formed; apply proves the file is well-shaped". Missing from apply: re-anchorability, group-count bounds, disjointness. Missing from lint: `${N}` within group count (statically provable via `NumSubexp`).

### Claim 8 (read-back through the role's reader is meaningful)
Negated by P-15, P-16. The role's reader parses the record. The README is not the record. A write entry whose rules span two files has one reader; the second file is unverified or incomplete on every run. Byte equality also assumes the reader does not normalize, which projectors do.
Answer required: a reader per edited file (rule carries its own read-back path or the entry is split per artifact), and either a declared normalization or a documented "reader must emit the literal cell".

### Claim 9 (no new refusal class)
Negated by P-17. `execution_failure`'s existing semantics for the `command` carrier include non-zero exit after partial mutation; the class is not "not applied", it is "the carrier failed". Putting "anchor matched 0 lines", "value contains terminator", "rename failed after temp write" all into it with a detail string means the agent loop that drives intrastate cannot branch on them. Distinct cases with distinct recovery (fix the model vs fix the value vs retry) collapse into one.
Answer required: either a pre-mutation refusal class (`plan_refused`, "not applied" by contract) or structured detail codes on `execution_failure`.

### Claim 10 (one rule per key, one buffer)
Negated by P-2, P-18. The README row is one line with two owned keys. The rule-per-key model with same-line refusal cannot express it. And a buffer with deletions needs an ordering rule.
Answer required: rules keyed by anchor carrying one or more placeholders; every key covered by exactly one rule; deletions applied last or on a pre-split slice.

### R1 (wrapper script)
The wrapper is in the repo the model is in and its argv is lint-covered by the command carrier. What lint cannot cover is the wrapper's body; what lint cannot cover in `edit` is whether the regex selects the right line. The rejection is symmetric. The real reason to prefer `edit` is the in-process refusal before write, which a wrapper cannot promise. Say that instead.

### R2 (sub-verb via command carrier)
"No safety gain" is wrong: the gate is the only consent step for writing into caller-owned files, and P-8 shows why it matters. The sub-verb also gives a unit-testable CLI surface (`intrastate edit --anchor ... --replace ... < value`) that the in-process carrier lacks. If `edit` is chosen anyway, it needs its own consent (path scoping) to match what the sub-verb would have had for free.

### R3 (adapter registry)
P-11 and P-20: `clear` alone already forks on document shape (delete a bullet line vs blank a table cell). The registry is being rejected in name and accepted as an enum. Either accept that `clear` has dispositions and say so, or keep one disposition and accept the README row cannot be cleared.

### R4 (emit-only)
The caller here is an agent loop, and an emitted structured patch (path, line anchor, replacement) applied by a lint-covered `intrastate apply-patch` is exactly `edit` with a persisted intermediate. The rejection conflates a human retyping a value with a mechanical apply of an emitted patch. The stronger reason is that emit-only cannot do the pre-write read-back gate check either. Say that.

### R5 (no append on unmatched anchor)
P-21: after the design's own `clear`, the unmatched anchor is a line the tool removed. Either `clear` must not delete (blank instead) or R5 needs an exception for a rule with a declared `clear` disposition and a declared `insert-after` anchor. As written, R5 and claim 6 contradict each other.

---

## 3. Consumer artifacts (what would have caught each at review)

- P-1, P-2: **User journey "index two records in one README"**: model with one write entry, two records, README rows with Status and Updated cells. Expected: both rows editable by one declared rule set. Would have failed at the anchor-without-placeholder lint and at the same-line refusal.
- P-3: **Test `TestEdit_ValueWithPipeInTableCell`**: planned `status = "Final | note"`, table-cell rule; expect pre-write refusal, not post-write mismatch.
- P-4: **Test `TestEdit_SecondApplyAfterValueContainsPeerAnchor`**: run 1 writes a value that contains another rule's anchor text; run 2 must still select one line. Expect refusal at run 1, not run 2.
- P-5: **Test `TestEdit_UnanchoredRegexPreservesPrefix`**: anchor `Status\*\*: (\w+)`, replace `Status**: {status}`; assert the line still begins with `- **`. Fails as designed; forces the whole-line/span decision.
- P-6: **Test `TestTemplate_LiteralBracesAndDollar`**: replace containing `{{ x }}` and `$5`; expect a defined escape and a lint error naming the rule.
- P-7: **Test `TestEdit_GateOffCommandReader_NoMutation`**: `--allow-commands` absent, reader is a command, edit rule valid; assert the file's bytes and mtime are unchanged and the refusal names the gate.
- P-8: **Test `TestEdit_RefusesPathUnderDotGitAndSymlink`**: target `.git/hooks/pre-commit` and a symlinked README; expect refusal before any read.
- P-9: **Test `TestEdit_PreservesModeAndRefusesConcurrentChange`**: 0600 file; second writer modifies between read and rename; expect mode preserved and a refusal (not a lost update).
- P-10: **Test `TestEdit_LastLineNoTrailingNewline_ClearAndSet`**: golden bytes before and after for LF, CRLF, mixed, and no-final-newline.
- P-11, P-21: **User journey "withdraw a record, then reinstate it"**: clear status, then set status back to Draft. Expect success on both; fails at the second step by design.
- P-12: **Test `TestReadBack_ClearAgainstDefaultingReader`**: reader emits `Draft` when the bullet is absent; clear must be reported honestly (mismatch or a declared absence semantics), never a false pass.
- P-13: **Lint test `TestLint_BackrefIndexWithinGroupCount`** and **apply test `TestEdit_RewrittenLineStillMatchesAnchor`**.
- P-14: **Fixture: a record that quotes the template bullet inside a code fence** plus a peer citation line; expect a documented scoping recipe to select the real one.
- P-15: **User journey "verify the README write"**: one entry, rules in two files, one reader; expect the README key to be read back by something. Fails: no reader for it.
- P-16: **Test `TestReadBack_ReaderNormalizesCase`**: reader lowercases; expect a declared normalization or a refusal at lint ("reader contract"), not a permanent mismatch.
- P-17: **Test `TestRefusal_PreWriteClassIsNotApplied`**: for each pre-write refusal, assert a machine-readable code distinguishing it from a carrier failure after mutation; assert the agent loop's retry policy branches on it.
- P-18: **Test `TestEdit_ClearShiftsLaterLineIndex`**: rule A clears line 5, rule B sets line 9; assert line 9's original content is the one rewritten.

---

## 4. Refutation targets (seed classes applied)

**S1 (routed into an existing class without checking its veto semantics)**: two places. (a) Claim 9 routes every edit-time refusal into `execution_failure` and asserts "means not applied". The class as it exists for the `command` carrier is reached after a process ran and may have written; its consumer-facing semantics are "carrier failed", not "not applied". The claim names a partition (pre-mutation vs post-mutation inside `execution_failure`) the code does not have. (b) Claim 4 routes `edit` outside the `--allow-commands` gate on the premise the gate partitions "spawns a process" from "does not"; the gate's actual veto is "a model may cause execution"; a declared line edit to a hook or build file is in the vetoed class.

**S2 (widening asserted, narrowing shipped)**: claim 4's "a gate-off read-back can be refused BEFORE mutation" is a widening (edit gains a pre-write refusal) that is a narrowing if the gate is evaluated where it is evaluated today, in the reader path above nothing and after the write. Control flow: write, then reader, then gate check returns. The artifact is applied-but-unverified: the exact state the claim says cannot occur. Second instance: claim 2's "safer than lineinfile" is a narrowing for the consumer's records that legitimately contain two candidate lines.

**S3 (filter rests on a cardinality premise)**: "two rules select the same line -> refuse" rests on "one owned key per line". The README index row carries two owned keys on one line. The filter therefore refuses the row when both keys are owned and accepts it when only one is: the same non-monotonicity as S3, and the RDR's stated consumer is the instance. Second cardinality premise: "the anchor selects exactly one line" rests on "one status bullet per record"; template quotes, peer citations, and history tables make two.

**S4 (borrowed authority)**: three borrowings. (a) "the existing executor rule: read-back through the artifact role's declared reader" is borrowed for a write whose second target (README) the role's reader does not read; the rule's entry says the reader reads the artifact, not every file a write touches. (b) "the existing rule that a clear is a removal whose read-back asserts absence" was written for the `path` JSON carrier where absence is a missing key; in a text line, absence is a missing line to a reader that may default it. The rule names no text-line case. (c) "atomic temp-and-rename as the `path` carrier does" is borrowed from a tool-owned flat file to a caller-owned, git-tracked, human-edited, possibly symlinked file; the `path` carrier's atomicity argument assumed sole ownership and the seam never said that assumption travels.

---

## Verdict reasoning

The approach survives: every failure above has a design answer inside the `edit` carrier's own vocabulary. But two of them (P-1, P-2) mean the design as written cannot serve the one consumer it names, and one (P-8) moves the tool from "the only thing that changes state in a source" to "a thing that can change any source", so the mitigations are blocking, not optional.
