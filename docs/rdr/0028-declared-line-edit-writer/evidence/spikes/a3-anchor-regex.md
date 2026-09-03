Model: claude-sonnet-5

# A3 spike — RE2 anchor expressiveness for the Status bullet and README row

## Assumption under test (verbatim, `rdr inspect --select 0028:A3 0028`)

- **A3 [Go's RE2 `regexp` is expressive enough for the consumer's two anchors —
  the `- **Status**:` bullet with an optional bracketed qualifier as a capture
  group, and the README `| [NNNN](` row with the status cell as a group — with
  no lookaround; AND each anchor selects exactly one line in every existing
  consumer record and README, template comments and quoted examples included
  (the premortem's decoy-line case)]**
  - Status: Pending
  - Method: Spike

## Patterns tested

```
statusAnchor      = ^- \*\*Status\*\*: (.+)$
readmeRowGeneric  = ^\| \[\d{4}\]\([^)]+\)[^|]*\| [^|]*\| ([^|]+) \|
readmeRowForID(id) = ^\| \[<id-quoted>\]\([^)]+\)[^|]*\| [^|]*\| ([^|]+) \|
qualifierAnchor   = ^- \*\*Status\*\*: (\S+)( \[.*)?$
```

`statusAnchor` is pinned at column 0 (`^- `) specifically to exclude the
nested `  - **Status**:` lines that appear inside Assumption/Detail blocks
(2-space indent) — those are a real decoy source in this corpus (0001 and
0006 alone carry 11 nested Status lines between them) and the zero-indent
pin defeats them cleanly with no lookaround.

## Commands run

```
cd /tmp/a3spike && go run main.go \
  /Users/cwensel/sandbox/newcoinc/intrastate/docs/rdr \
  /Users/cwensel/sandbox/newcoinc/rdr/tools/rdr/testdata/status/records/0021-cache-warmup-order.md
```

(go1.26.6 darwin/arm64; module a3spike, no external deps — stdlib
`regexp`, `regexp/syntax` only.)

## Lookaround check

`regexp/syntax.Parse` under `syntax.Perl` succeeded for both patterns; no
`(?=`, `(?!`, `(?<=`, `(?<!` token appears in either. **No lookaround was
needed or used.**

## Status-anchor counts — every `*.md` under `docs/rdr/` (excluding README.md)

33 files scanned (32 records + BUILD-ORDER.md, README.md scanned separately).

| Result | Count of files |
| --- | --- |
| count == 1 | 28 |
| count == 0 | 5 |

Files at count == 0 (all are non-record documents with no Metadata Status
bullet at all — postmortems and the build-order doc, not decoys the anchor
mis-selected):

- `0003-guard-predicate-exhaustiveness-postmortem.md`
- `0007-guard-predicate-totality-postmortem.md`
- `0008-recognized-tag-key-ownership-postmortem.md`
- `0011-flow-next-match-conditioned-candidates-postmortem.md`
- `BUILD-ORDER.md` (has a bare `**Status**: authored ...` line, no leading
  `- `, so it correctly does not match; it is not a record and was never in
  scope for this anchor)

**No file anywhere in the corpus scored ≥2** — including `0001-resolution-kernel.md`
(6 nested Status lines under Assumptions/Details: `Verified`×4, `Accepted`×1,
plus the top-level `Implemented`) and `0006-graph-lint-authority-and-guarantees.md`
(11 nested Status lines: `Verified`×4, `Pending`×6, plus the top-level
`Implemented`). The zero-indent pin correctly selects exactly the one
Metadata bullet in both, ignoring all nested decoys. This is the
strongest confirmation of the "no ≥2" half of A3 in this corpus.

Verdict on this half: the anchor is safe for every file that actually
has a Metadata Status bullet. The 5 zero-count files are out of the
anchor's intended domain (not RDR records in the Metadata-Status sense),
not anchor failures — worth flagging to the author as a scope note, not
a refutation.

## README anchor counts

Generic row anchor (`readmeRowGeneric`): 26 matches, one per linked
(`[NNNN](...)`) row in the Index table. Correct — matches the number of
rows that use the link form.

Per-record anchor (`readmeRowForID(id)`), tested against `README.md`:

| id | count | note |
| --- | --- | --- |
| 0001 | 0 | row is `\| 0001 \| ... \|` — NOT a `[0001](...)` link; anchor requires the link form |
| 0002 | 1 | ordinary linked row |
| 0006 | 0 | same non-link-row case as 0001 |
| 0026 | 1 | ordinary linked row |
| 0027 | 1 | ordinary linked row |
| 0028 | 1 | ordinary linked row |
| 9999 | 0 | absent id, correctly zero (not a false positive) |

**Finding**: two README rows (0001, 0006) are written without the
`[NNNN](file.md)` link — plain `| 0001 | ... |`. The RDR's own anchor
grammar (`| [NNNN](`, per the task brief) assumes every row is linked.
For these two ids the per-record anchor selects **zero** lines, which
under C3 is `edit_anchor_unmatched` — a refusal, not a wrong-line write.
This is not a regex-expressiveness gap (RE2 handles the link-row shape
fine) but a real corpus fact the author should know: 0001 and 0006 would
refuse a status-anchor `edit` against the README today, until either the
rows are given the link form or the anchor is relaxed to admit both
row shapes.

No README row (linked or not) ever scored ≥2 for a given id. No decoy
row collision found.

## Qualifier capture + rewrite demo

No live `docs/rdr/*.md` record in this consumer tree currently carries a
top-level bracket-qualified Status line (`Draft [joint decision → …]`) —
only prose/backticked mentions mid-sentence (e.g. `0008-recognized-tag-key-ownership.md:665`,
`0028-declared-line-edit-writer.md:685`), which the zero-indent anchor
correctly does not match. To exercise a real qualified Status line the
spike used the RDR engine's own test fixture, which has one:
`/Users/cwensel/sandbox/newcoinc/rdr/tools/rdr/testdata/status/records/0021-cache-warmup-order.md:9`.

### Case 1 — real fixture, wrapped qualifier (found, not synthetic)

```
BEFORE (line 9): "- **Status**: Final [joint decision → 0022-cache-metrics-surface §"
```

The bracketed qualifier in this fixture wraps onto the next physical line
(`  A1: whether warm-up counts toward the hit-rate metric]`, line 10).
Under C3's line-split semantics ("lines are split on `\n`"), the anchor
only ever sees line 9. The qualifier-capture group therefore captures
`" [joint decision → 0022-cache-metrics-surface §"` — truncated, missing
the closing `]` and the `A1: ...` continuation, which live on line 10.
`regexp.Expand`/`ReplaceAllString` on line 9 alone reproduces the line
unchanged (word `Final` is already the target value in this demo, so no
value swap was needed to show the truncation):

```
AFTER (regexp.Expand via ReplaceAllString): "- **Status**: Final [joint decision → 0022-cache-metrics-surface §"
```

**This is a genuine boundary finding**: a `replace` template using
`${2}` to preserve a wrapped qualifier will silently drop everything
after the line break — no error, no refusal, just a truncated qualifier
in the rewritten line, and the orphaned continuation line (`  A1: ...]`)
left dangling below it with no clause tying it back. No live intrastate
record hits this today (see above), but the RDR engine's own testdata
proves the shape occurs in this document family, so C3 or the anchor
design should account for it (e.g. author guidance to keep qualifiers
on one physical line, or a documented refusal/warning when a
Status-anchor capture group's tail looks unterminated).

### Case 2 — synthetic single-line qualifier (control, shows the intended happy path)

```
BEFORE: "- **Status**: Draft [joint decision → JDR 0003 §D1]"
AFTER  (regexp.Expand, ${2} preserves qualifier): "- **Status**: Final [joint decision → JDR 0003 §D1]"
```

When the qualifier stays on one physical line, `${2}` preserves it
exactly across a `Draft`→`Final` swap, via both `ReplaceAllString` and
the lower-level `regexp.Expand` (submatch-index form). This is the
happy path C2/C3 assume.

## Verdict

**PASS with one qualified finding.** RE2 with no lookaround suffices for
both anchors; the zero-indent pin defeats every nested-Status decoy in
the corpus (0001, 0006's Assumption blocks); every Status-bearing record
scores exactly 1; every README link-row scores exactly 1 for its own id
and 0 for others. Two real, non-hypothetical gaps surfaced that the
author should see before lock:

1. README rows 0001 and 0006 lack the `[NNNN](...)` link form the anchor
   assumes — a real `edit_anchor_unmatched` refusal risk against the
   README today, not a regex expressiveness problem.
2. A bracketed Status qualifier that wraps across a markdown continuation
   line (proven to occur in this document family via the RDR engine's own
   testdata fixture) truncates silently under `${N}`-preserving `replace`
   templates, because C3's anchor is strictly per-line. No live
   intrastate record hits this yet, but it is a live shape in the family
   and worth a design note.

Neither finding requires lookaround or a second regex dialect — both are
either author-corpus fixes (link the two rows) or a C2/C3 clarification
(qualifiers must stay single-line, or the anchor design should detect
and refuse an apparently-unterminated bracket rather than silently
truncate).
