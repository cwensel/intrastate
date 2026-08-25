# Triage — RDR 0007 Guard predicate totality

Single-pass unattended triage of the roborev findings left by the Stage-8
implementation launch. Window frozen once at Phase 0 and never re-read.

- **Window**: `d74998a..e1db52a` (11 commits)
- **Batch label**: `batch:rdr-0007`
- **Spine**: per-commit auto-reviews, jobs 6015–6020
- **Range net**: not run separately — the five in-window per-commit jobs
  cover every commit that touches `internal/`; no cross-commit interaction
  defect was outstanding.

Of 17 open roborev jobs repo-wide, **5 were in-window**. The other 12 are
pre-existing findings on `main` or other branches — the class a moving
`merge-base..HEAD` refine loop would have re-flagged as this branch's work.

## Findings × verdicts

| # | Job | Location | Verdict | odc-type | odc-trigger | Outcome |
| --- | --- | --- | --- | --- | --- | --- |
| F1 | 6015 | `guard_mvv_test.go:243` | KATA-BUG | test-oracle | test-coverage | kata `5wrc` |
| F2 | 6016 | `guard.go:236` | DROP over-engineering | algorithm | boundary | drop |
| F3 | 6017 | `guard_adversarial_0007_test.go:75` | DROP superseded-at-HEAD | n/a-not-a-defect | — | drop |
| F4 | 6017 | `guard_adversarial_0007_test.go:86` | DROP rdr-adjudicated | n/a-not-a-defect | design-conformance | drop |
| F5 | 6017 | `guard_adversarial_0007_test.go:271` | DROP superseded-at-HEAD | n/a-not-a-defect | — | drop |
| F6 | 6018 | `guard_fixup_0007_test.go:164` | DROP superseded-at-HEAD | n/a-not-a-defect | — | drop |
| F7 | 6020 | `resolve.go:156` | KATA-BUG | algorithm | rare-situation | kata `qdnj` |

**Defect distribution** (excludes the four `n/a-not-a-defect` pre-filtered
rows): `test-oracle` 1, `algorithm` 2.

## The two filed

### `qdnj` — P1, `severity:high` — duplicate-key positionality escapes into `Row.Match`

The one genuine escape. `assemble` marks a within-provenance duplicate
`conflicted` but retains the last value; `TagSet.matches` (`resolve.go:156`)
compares that value and never consults `conflicted`, while `conflicting()`
IS consulted on the guard path. Reproduced over real `Resolve`:
`Observed=[gate=open, gate=closed]` → `no_match`; reversed → `PLAN A`.
Same multiset, opposite disposition, silently.

Phase 3b found this defect class (ADV-3) and `a87fc9e` closed the
**guard-verdict** leg; the orchestrator scoped that fix to the kernel guard
boundary (D9). `Row.Match` selection is a **second consumer of the same root
cause** the ruling never had in view — surfaced only by reviewing the fix
commit itself. Out of RDR 0007's contract scope (`0007:C8` is scoped by
`req-list.md:146` to refusal mapping and the payload; `0002:279` states the
candidate filter in RDR 0002's voice), so filed rather than fixed.

Carries an `## Open question`: fail-closed `no_match` is escapable under
0002, so a modeled escape row could still route around a conflicted key.
Recommendation baked into the kata — fail-closed now, route escapability to
RDR 0002/0009.

### `5wrc` — P2, `severity:medium` — REQ-71 meta-test is a weak oracle

`guard_mvv_test.go:243` is byte-identical from the RED commit to HEAD, so
the red-green pre-filter does not apply. It records which questions
`TestGuardEvaluatorContract` asked, never the verdicts required; floors are
`>=2` operators + one varied value axis, and the line-301 unparseable
heuristic misclassifies valid non-numeric `in` values. `guardcontract.go` is
strong today (4 operators, 12 cases, asserted verdicts) and the artifacts
record the right mutants killed — but those were one-time manual acts.

Cross-RDR: RDR 0003 instantiates this exported harness at BUILD-ORDER step 3.
Kata targets landing before `/rdr-implement 3`. Note `deviations.md` D2 item 4
forecloses the suggested subprocess fix (Go propagates subtest failure); the
kata proposes assertable alternatives.

## The drops worth recording

**F4 — `rdr-adjudicated`, the branch's central interpretation.** roborev
argued the fail-closed block boundary tests malformed input outside the
normalized contract. Refuted on three independent grounds:

1. *"Defines exactly three valid block values"* inverts the change under
   review. The RDR's own fence (`0007:1269-1271`) says **exactly two**;
   JDR 0001 §D12 widened it to three, landing `BlockMatch` in this RDR's
   Phase 1 — which is what makes ADV-1/ADV-2 reachable at all.
2. The *"unreachable normalized rows"* premise rests on a contract not yet
   written: §D6 assigns block-keyed match routing to **RDR 0002**,
   unimplemented at BUILD-ORDER step 2.
3. The ADV-1 fixture is misread — `guardedRow` sets a **satisfied** `Match`
   pattern, so the row is a genuine candidate; the absent key is on a stray
   `BlockMatch` atom in `Row.Guard`.

Two of ADV-2's three legs (`Block("")`, a future token) are excludable by no
normalization contract. Q1 reading (a) stands; zero new public surface.

**F2 — `over-engineering`.** The cited text is real (`0007:2361`,
"Evaluation stays linear in the atom count per candidate row") but sits under
Validation, opens with "No performance dimension", and carries no `0007:CN`
element id — a characterization, not a fenced bound. `appendAtom` runs only
on the undecided branch, so the quadratic term is over *distinct unevaluable*
atoms within one hand-authored row. The dedupe is load-bearing and documented;
replacing it would disturb a comparator whose totality D8 just repaired.

*Judgment call flagged*: if that sentence is later read as normative, F2
flips from a drop to a low-priority kata.

## Fix-now

None. No finding met the fix-now bar (contained, unambiguous,
data-safety/security, clearly in this RDR's scope, cheap).

## Closure

All 5 in-window jobs closed via `roborev close` after a `respond` comment
citing the verdict, evidence, and kata short_id or drop reason. Zero
in-window findings left open; nothing dropped silently.
