# triage — 0029 version promise on machine-readable output

Window frozen once: `HEAD=bdeea9a`, `BASE=2f7d22c` (24 commits), branch
`worktree-rdr-0029`. Batch label `batch:rdr-0029`.

Spine: of 79 open roborev jobs, exactly two are ancestors of `HEAD` and
outside `BASE` — job 7757 (`14801e9`, the red-suite commit, four findings)
and job 7766 (`4539201`, the Phase 3b adversarial commit, two findings).
The cross-commit range net (job 7775, `roborev review --since BASE`, run
once) added one further finding. Every other open job belongs to `main`,
`via-codex`, `via-claude`, or a foreign `worktree-*` branch and is out of
window. The main feat commit `74f25d3` carries no open job.

The eight `docs(rdr):`/`chore(rdr-0029):` artifact commits on this branch
carry no auto-review: they match roborev's `excluded_commit_patterns` and
hold no code. That is a correct skip, not a hole in the spine.

| # | Source | Location | Verdict | odc-type | odc-trigger | Outcome |
|---|---|---|---|---|---|---|
| 1 | job 7757 (`14801e9`) | `internal/cli/schema_mvv_0029_test.go:208-217` | DROP superseded-at-HEAD | test-oracle | design-conformance | drop:superseded-at-HEAD |
| 2 | job 7757 (`14801e9`) | `internal/cli/schema_closed_wording_0029_test.go:93-100` | DROP superseded-at-HEAD | test-oracle | design-conformance | drop:superseded-at-HEAD |
| 3 | job 7757 (`14801e9`) | `internal/cli/schema_tiers_0029_test.go:61-76` | DROP over-engineering | test-oracle | rare-situation | drop:over-engineering |
| 4 | job 7757 (`14801e9`) | `internal/cli/schema_register_0029_test.go:211-214` | OUT-OF-SCOPE | test-oracle | test-coverage | kata 5h4a |
| 5 | job 7766 (`4539201`) | `internal/cli/schema_adversarial_0029_test.go:117` | DROP superseded-at-HEAD | test-oracle | design-conformance | drop:superseded-at-HEAD |
| 6 | job 7766 (`4539201`) | `internal/graphlint/severity_snapshot_adv_0029_test.go:37` | DROP rdr-adjudicated | test-oracle | test-coverage | drop:rdr-adjudicated |
| 7 | job 7775 (range net) | `internal/graphlint/vocabulary_snapshot_0029_test.go:96` | DROP rdr-adjudicated | test-oracle | test-coverage | drop:rdr-adjudicated |

## Evidence per row

**1 — superseded at HEAD.** The finding names the in-process MVV as
jointly unsatisfiable: step 1 requires `data.findings == []` at `"0.1"`
while step 2 requires a new advisory code to fire on that same model. The
hunk was rewritten before HEAD: the empty-findings assertion is gone and
the file now carries the explicit retirement — "The pre-change
empty-findings baseline is NOT asserted here … `0029:S7`'s matrix carries
no empty-findings row for that reason." That is D2's recorded resolution
(`deviations.md`, RESOLVED against `0029:S7`), applied in source.

**2 — superseded at HEAD.** The census sweep was file-scoped, so it
flagged the two lines `TestReq17And63` requires verbatim. At HEAD
`describesUntieredVocabulary()` exists and is called inside the census
loop, exempting the untiered capability and validation-code sets. That is
D3's recorded resolution (REQ-63, `0029:C2`, `0029:S9`) applied to the
assertion rather than to `model.go`, exactly as D3 directed.

**3 — over-engineering.** The concern is `seamLiteralMembers`' AST
literal-scrape falling back to `pkgVarStrings` and merging distinct
vocabularies. Grounding found all three live call sites
(`schema_tiers_0029_test.go:120,145,205` → `respond.Types`,
`respond.Levels`, `cli.UnknownReasons`) return literals directly in-body,
so the fallback is never invoked and no live seam can be silently wrong.
`verification.md` independently confirms `UnknownReasons=[absent
uncomparable not-evaluated]` by value. Defense against a shape no seam
currently takes.

**4 — OUT-OF-SCOPE → kata `5h4a`.** The "is there a reader" meta-check
scans for any `*_test.go` containing `"testdata"` plus `"vocabular"`/
`"seams"`, and the scanning file itself satisfies that predicate. The
self-satisfaction is real. But REQ-21's obligation is discharged by
`internal/graphlint/vocabulary_snapshot_0029_test.go::TestReq21_TheVocabularySnapshotMatchesTheTree`,
a genuine byte-diff against the committed snapshot — proved red in a
detached probe worktree by mutating `graph-vacuous-atom info` →
`blocking` (the exact undisclosed-promotion scenario), then green on
restore. So this is redundant scaffolding beside a working oracle, not a
breach of `0029:C3`: severity low, priority 3, `release:post-1.0`.

**5 — superseded at HEAD.** ADV-1 checked code membership against
`BlockingCodes() + AdvisoryCodes()` without comparing the emitted
severity. At HEAD the predicate is `snapshotRecordsSeverity(snapshot,
code, severity)` — it asserts the committed snapshot carries a
`<code>\t<severity>` row for any code emitted with a severity, which is
the finding's own requested fix. Recorded as ADV-1 `fixed:fa3e3db` in
`verification.md`, under D5's RESOLVED reading 1.

**6 — rdr-adjudicated.** ADV-3's own assertion is header-only
(`strings.Contains(rendered, "[graphlint.Severities]")`) and unchanged at
HEAD, so the finding describes the bytes accurately. It is nonetheless
not a live hole: `vocabulary_snapshot_0029_test.go:103` renders
`graphlint.Severities()`'s members into the section body, and
`TestReq21_TheVocabularySnapshotMatchesTheTree` golden-diffs the whole
rendered output. Proved in a detached probe: adding a third severity
member turns that golden test red while ADV-3 itself stays green. The
failure mode ADV-3 names is caught; ADV-3 is a redundant weak oracle
beside the golden test. `verification.md`'s "ADV-3 — fixed" holds in
effect.

**7 — rdr-adjudicated.** The graphlint golden snapshot covers only
downward-importable vocabularies, per its own header comment, and `cli`
is outside that allow-list. The C4 row for `cli.UnknownReasons` is
REQ-34 — distinct from REQ-36's genuinely unbuildable CLIError-`code` row
(A5) — and is already pinned by value and for uniqueness in
`internal/cli/schema_tiers_0029_test.go:204-236`, with declaration and
shape checks in `schema_probe_0029_test.go:234-292`. The seam is
reviewably covered in the layer that owns it; pulling `cli` into
graphlint's renderer would buy nothing.

## Out of scope by standing decision (not filed)

- The six `modernize: errorsastype` reports from `golangci-lint`
  (`internal/accessor/executor.go`, `internal/cli/clierr/clierr.go`,
  `internal/cli/root.go`, `internal/table/category.go`) are pre-existing:
  `main` carries the same six, and no `errors.As` line was touched on this
  branch. Not this record's debt. No finding raised them.
- The five REQ demoted to EXCLUDED as unassertable (20, 23, 39, 40, 41)
  are grounded in the record's own bytes. No finding claimed one of them
  was uncovered.

## Result

7 findings · 6 DROP · 1 OUT-OF-SCOPE (`5h4a`) · **0 IN-SCOPE / 0 open**.

No row is `open`, so the launch's Phase 3c has nothing to fix and the
completion gate does not hold on `impl_findings_open`. Every originating
job (7757, 7766, 7775) is responded and closed with its per-finding
verdicts.
