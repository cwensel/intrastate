# Triage — RDR 0027 inline-shell-detection-scope

Window frozen once: `BASE=f85fde7` .. `HEAD=99261f3` (12 commits).
Batch label: `batch:rdr-0027`.

Spine: 3 in-window per-commit auto-reviews (jobs 6733, 6736, 6737) carrying
7 findings. Cross-commit net (job 6738, `--since f85fde7`) surfaced 2
findings, both duplicates of spine findings — deduped by
(source-anchor, problem-gist) and charged to the per-commit job so the
close path stays unambiguous.

Note: the two commits carrying the actual contract change — `17fcc4d`
(position-free predicate) and `2726a44` (promise wording) — reviewed
CLEAN. Every finding lands on test scaffolding, not on the shipped
behaviour.

| # | Job | Location | Problem (gist) | Verdict | odc-type | odc-trigger | Outcome |
|---|-----|----------|----------------|---------|----------|-------------|---------|
| F1 | 6733 | `internal/table/inline_shell_scope_0027_test.go:499` | admitted row `["xargs","busybox","sh","-c","echo"]` carries a real `sh -c` pair and must refuse | DROP `superseded-at-HEAD` | test-oracle | design-conformance | fixed in Phase 2 as deviation D4; vector now `["xargs","busybox","-c","echo"]` |
| F2 | 6733 | `internal/cli/inline_shell_promise_0027_test.go:314` | fixture omits `[tags.recognized]`, refused as `malformed_model_declaration` before reaching C5 | DROP `superseded-at-HEAD` | test-oracle | design-conformance | fixed in Phase 2 as deviation D7; `tags.recognized` present at HEAD |
| F3 | 6733 (+6738 dup) | `internal/cli/inline_shell_mvv_0027_test.go:190` | MVV step-3 exec arm never runs the `env -S` one-word shell form; no capability check exists despite the comment claiming a skip | KATA-BUG | test-oracle | test-coverage | filed |
| F4 | 6736 (+6738 dup) | `internal/cli/inline_shell_adversarial_0027_test.go:211` | ADV-3 never checks the refused paragraph for overlap with admitted forms | DROP `unreachable` | test-oracle | test-coverage | refused paragraph names zero concrete forms; ADV-1 fails by construction if `-s` joins `shellInterpreters["sh"]` |
| F5 | 6737 | `internal/table/inline_shell_authoring_0027_test.go:36` | REQ-0b guard accepts any `0027:C1` occurrence file-wide | KATA-BUG | test-oracle | test-coverage | filed (collapsed with F6, F7) |
| F6 | 6737 | `internal/table/inline_shell_authoring_0027_test.go:71` | REQ-35 guard only checks strings occur somewhere; REQ-74 loop can be bypassed | KATA-BUG | test-oracle | test-coverage | filed (collapsed into F5) |
| F7 | 6737 | `internal/table/inline_shell_authoring_0027_test.go:107` | REQ-40 guard satisfied by unrelated `q2q1`/`subsumed` tokens elsewhere in the file | KATA-BUG | test-oracle | test-coverage | filed (collapsed into F5) |

## Consolidation

F5/F6/F7 share one root cause — each guard asserts a file-wide
`strings.Contains` rather than scoping to its named locus — so they are
collapsed into a single kata. Jobs 6737 closes once, citing that kata.

## Evidence notes

**F1/F2 (DROP superseded-at-HEAD).** Both were flagged against `b51572b`,
the deliberately-RED Phase 1 commit, and both were fixed by the Phase 2
implementer, who recorded them as deviations D4 and D7. Verified at HEAD:
the `busybox` vector no longer carries `sh -c`, and the promise fixture
declares `[tags.recognized]`. This is the flow's own red-green rhythm,
not an escape.

**F3 (KATA-BUG).** `0027:MVV` step 3 says the first two admitted forms are
"shown, by running the model under `--allow-commands`, to behave as C1
states (the wrapper file runs; the `-S` string runs a shell …)". Only the
wrapper-file half executes. The comment at :196 claims the `-S` half
"skips on a BSD-`env` host", but no capability check and no `-S` subtest
exist, and `coverage.md` records it as skipped-on-this-host — untrue.
This host's `/usr/bin/env` DOES support `-S` (`env -S "echo hi"` → `hi`,
rc=0), so A4's darwin premise is stale here and the half is reachable and
unexercised. Contained and testable now, so KATA-BUG rather than RDR-SEED.

**F4 (DROP unreachable).** The refused paragraph
(`internal/table/category.go:126-133`) names zero concrete interpreter+flag
forms — it is structural ("a listed interpreter word followed, at ANY later
argv position, by one of that interpreter's own inline-code flags"), and
its only enumerations are wrappers. There is no concrete-form set on the
refused side to intersect. The scenario also requires `-s` to be added to
`shellInterpreters["sh"]`, which fails ADV-1 immediately, since
`{"sh -s", []string{"sh","-s"}}` sits in `admittedInDescription`. Phase 3a
independently found no over-claim. The suggested fix — rendering both
sections from shared structured claim data — would invert C1's deliberate
design: ADV-1 transcribes the left-hand side from RENDERED text precisely
because "transcribing from the record would test the record against
itself".

**F5-F7 (KATA-BUG).** Confirmed by three TARGETED mutations, all of which
left the guards green: rewriting `0027:C1`→`0025:C5` in both C5 comment
blocks (unrelated citations at :1035/:1182/:1200 still satisfy the check);
truncating the REQ-74 loop to `cases[:1]` and replacing `t.Fatalf` with
`return` (all four argv literals remain); and replacing D6's disposition
with `TBD` (other lines supply every token). The Phase 3c mutation
verification used coarse global `sed`, which is exactly why it passed.
`MVV/step_4` already independently executes all four REQ-74 vectors with
the category assertion, so F6's suggested "independently execute" fix is
redundant — the open question is whether that guard should cite step_4 and
be deleted rather than strengthened.
