# Triage — RDR 0010 Stateless Decision Tables

Single-pass unattended triage of the roborev findings left by the Stage-8
launch. Window frozen once: `BASE=f990bee`, `HEAD=7d81d59` (18 commits).
Batch label: `batch:rdr-0010`.

- **Spine** — 7 in-window per-commit auto-reviews (jobs 6313–6317, 6320,
  6323), 14 distinct findings. 27 further open jobs were repo-wide backlog
  whose refs are not ancestors of `HEAD`; out of window, untouched.
- **Range net** — one `--since BASE` review (job 6325), 6 findings, **all**
  duplicates of the spine. Deduped by `(source-anchor, problem-gist)` and
  charged to the originating per-commit job.
- **Fix-commit sweep** — bounded to one round: jobs 6326 (1 finding, filed)
  and 6327 (clean).

Every finding is grounded against the RDR + `{RDR_RESOURCES}`, which roborev's
repo-sandboxed prompt never sees. Two findings were correct on the facts but
wrong on the spec, and one carried a factually wrong remedy — see DROP notes.

## Dispositions

| # | Job | Finding | Verdict | ODC type / trigger | Outcome |
|---|---|---|---|---|---|
| 1 | 6313 | `class_0010_test.go:160` — `dtClassOmitted` declares owned tag `status` | DROP | n/a-not-a-defect / design-conformance | `superseded-at-HEAD` — D11 repointed at class-stripped `smZeroOwned` |
| 2 | 6313 | `class_0010_test.go:259` — requires `malformed model declaration` | DROP | n/a-not-a-defect / design-conformance | `superseded-at-HEAD` — d8e5b9a moved the check (ADV-3/D14) |
| 3 | 6313 | `emit_0010_test.go:69` — REQ-14 clear-list prohibition unexercised | KATA-BUG | test-oracle / test-coverage | `kata stn1` |
| 4 | 6313 | `emit_0010_test.go:161` — sm-escape-with-emit quadrant missing | KATA-BUG | test-oracle / test-coverage | `kata stn1` (collapsed) |
| 5 | 6314 | `class_0010_test.go:127` — escape-row group naming / single-finding counts | DROP | n/a-not-a-defect / design-conformance | `superseded-at-HEAD` — D12 adjudicated; exact-set assertions now |
| 6 | 6314 | `class_0010_test.go:748` — two-outcome + stray-terminal controls unbounded | FIX-NOW | test-oracle / test-coverage | `a1a3156` |
| 7 | 6314 | `class_0010_test.go:315` — REQ-54 arms unexercised over the class | RDR-SEED | test-oracle / test-coverage | `kata bycb` ★ |
| 8 | 6315 | `decision_table_0010_test.go:419` — vacuous gate invariant | KATA-BUG | test-oracle / test-coverage | `kata stn1` ★ |
| 9 | 6315 | `decision_table_0010_test.go:446` — `bare-otherwise` unreachable | KATA-BUG | test-oracle / test-coverage | `kata stn1` (collapsed) |
| 10 | 6316 | `failure_modes_0010_test.go:163` — persistence counts entries only | KATA-BUG | test-oracle / boundary | `kata stn1` (collapsed) |
| 11 | 6316 | `emit_shape_0010_test.go:292` — `TestReq110` tautological | DROP | test-oracle / design-conformance | `rdr-adjudicated` — REQ-110 is deliberately NOT a review-time orphan |
| 12 | 6317 | `flow_next_0011_test.go:751` — `roundtrip_test.go` missing from guard | DROP | n/a-not-a-defect / design-conformance | `superseded-at-HEAD` — a32cf53 extended the allow-list |
| 13 | 6317 | `normalize.go:363` — empty `[rule.write]` / `clear = []` accepted | DROP | checking / design-conformance | `rdr-adjudicated` — C2 keys on CONTENT; REQ-15 forbids the proposed arm |
| 14 | 6320 | `adversarial_0010_test.go:108` vs REQ-9 — mutually exclusive contracts | DROP | n/a-not-a-defect / design-conformance | `superseded-at-HEAD` — both now assert `CatUnknownTag` (D14) |
| 15 | 6323 | `model-authoring.md:34` — example lints blocking | FIX-NOW | documentation / doc-code-drift | `1810f22` |
| 16 | 6323 | `model-authoring.md:79` — `guard.all` claimed required | FIX-NOW | documentation / doc-code-drift | `1810f22` |
| 17 | 6323 | `model-authoring.md:102` — "names each uncovered cell" | FIX-NOW | documentation / doc-code-drift | `1810f22` |
| 18 | 6326 | `model-authoring.md:78` — rows vs cells (fix-commit review) | KATA-BUG | documentation / doc-code-drift | `kata e8dt` |

★ = carries an `## Open question` for the shipping run.

## The two spec-level DROPs

Both were factually correct and would have been wrong to act on.

**#13 `normalize.go:363`.** Empty `[rule.write]` and `clear = []` do load and
normalize to absent — verified empirically. But C2's prohibition clause is
self-qualifying: *"None of these is a new refusal: each is already
unauthorable once the owned set is empty"*, and every arm C2 names keys on
**content** (a write/clear *key* that is not an owned tag). An empty block
names no key, so C2's own carrying mechanism cannot reach it, and **REQ-15
explicitly forbids adding the refusal arm the fix proposed**. C3 draws the
presence-vs-content line deliberately for this same function and assigns
"even an empty one" to the **escape** arm — exactly where `normalize.go:343-347`
implements it. The code conforms to the record in all four arms.

**#11 `TestReq110`.** `coverage.md` §Orphans lists exactly five review-time
REQs (86, 87, 88, 100, 106). REQ-110 is deliberately **not** among them — it
is tabled in §L as a runtime characterization pin on fixture-identifier
provenance. Deleting it, as proposed, would contradict the RDR's own orphan
accounting.

## Corrections to the findings themselves

- **#17** asked the doc to promise "one witness assignment". The shipped
  finding carries no witness — `clierr.Finding` has no such field. The
  correction omits it rather than documenting a promise the code does not
  keep. Adding a witness would be a code change and an RDR seed; not filed.
- **#6**'s stated rationale (that bounding these two fixtures guards the
  ADV-1/ADV-2 class-suppression arms) is **false**, and was verified so:
  disabling both class guards leaves both reports byte-identical. Those arms
  do not fire over these fixtures — `TestAdv0010_*` owns them, on fixtures
  built to trip them. The bound was landed for the guarantee it does provide
  ("nothing else fired"), with the comments corrected to say only that.

## FIX-NOW detail

- **`1810f22`** — three `docs/model-authoring.md` drifts, each verified
  against the built binary: the primary example lints blocking
  (`graph-coverage-gap`, 3 of 4 cells) and is now labelled a fragment;
  `guard.unless` participates as a dimension alongside `guard.all`
  (`guard.Dimensions` collects both via `guardAtoms`), so the section no
  longer claims otherwise and its heading was narrowed to match;
  `graph-coverage-gap` reports a per-arm uncovered count and dimension list,
  not a per-cell enumeration.
- **`a1a3156`** — bounded SC-4's two unbounded oracles. Probed the real
  finding sets first (`dtTwoOutcomeEscape0010` → exactly
  `[graph-coverage-closed-by-escape@dt/decide, graph-coverage-gap@dt/review]`;
  `dtStrayTerminal0010` → exactly `[graph-dangling-edge@terminal[0]]`). Added
  `codesOf`, a duplicate-**preserving** companion to `codesIn`, since a
  deduping bound hides a regression emitting one code twice — falsified
  against an injected spurious finding, which `TestReq96` caught *doubled*.
  Extended 0011's diff guard with `fixtures_0006_test.go` per REQ-87/88 and
  the D7/D9 convention; the guard fired on this edit as designed.

## Filed

| Kata | Kind | Title | Labels |
|---|---|---|---|
| `stn1` | `type:bug` | Tighten 0010's vacuous and quadrant-incomplete test oracles | `src:roborev` `batch:rdr-0010` `area:test` `severity:low` `release:post-1.0` `lifecycle:queued`, p3 |
| `bycb` | `kind:rdr-seed` | REQ-54's non-coverage invariant arms are unexercised over the decision-table class | `src:roborev` `batch:rdr-0010` `area:graphlint` `severity:medium` `release:post-1.0`, p3 |
| `e8dt` | `type:bug` | model-authoring.md conflates rows with cells in the completeness rule | `src:roborev` `batch:rdr-0010` `area:docs` `severity:low` `release:post-1.0` `lifecycle:queued`, p4 |

All three are `--related zdat` (the `kind:rdr-tracked` tracker). `bycb` takes
no `lifecycle:*` — it exits to RDR authoring, and kata-ship gate 4 refuses a
seed.

## Open questions carried into the katas

- **`stn1`** — `TestReq42`'s **deny** half may not be dischargeable:
  `verification.md` §Undecidable records that no runtime gate-deny witness
  could be authored on a decision-table fixture. The `verdict`/`result`
  field-name correction is required regardless of how the deny half resolves.
- **`bycb`** — whether the **node-ceiling** arm is reachable at all for this
  class (a ∅ root means one node). Determine reachability per arm first and
  document the vacuous arms as vacuous; a test asserting an unreachable arm is
  the same defect in the other direction.

## Not a defect

Findings 1, 2, 5, 12, 14 are the flow's own red-then-green rhythm — raised on
deliberately-red Phase 1 commits and resolved by later in-window commits.
Recorded `odc-type: n/a-not-a-defect` and excluded from the defect
distribution. Genuine defects surviving at `HEAD`: **9** of 14 spine findings
(4 fixed now, 5 filed), plus 1 from the fix-commit sweep.
