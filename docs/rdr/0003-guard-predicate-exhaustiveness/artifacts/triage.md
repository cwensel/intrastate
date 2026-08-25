# Triage — RDR 0003 Guard Predicate Exhaustiveness

Single-pass unattended triage of the roborev findings left by the Stage 8
launch. Window frozen once: `BASE=f4e6b19`, `HEAD=fda8e70` (18 commits).
Batch label `batch:rdr-0003`.

Spine: 9 in-window jobs (6083, 6084, 6085, 6087, 6092, 6093, 6094, 6096,
6097) carrying 20 findings. Cross-commit net: job 6099 over
`f4e6b19..fda8e70`, deduped against the spine (it independently re-raised
the held-`null` and typed-held-value defects, charged to 6097 and 6085).
Job 6098 reviewed the D15 retirement commit and returned **no issues
found**.

Every finding was grounded in a sub-agent against the RDR set, RDR 0002,
0006 and 0007, the artifacts, and `.rdr/resources.md` — the context
roborev's repo-sandboxed prompt never sees.

## Verdicts

| # | Job | Finding | Verdict | ODC type | ODC trigger | Outcome |
|---|---|---|---|---|---|---|
| 1 | 6097 | `parseHeldSet` reads a held `null` as the empty set | FIX-NOW | checking | rare-situation | `072c7a0` |
| 2 | 6087 | dimensionless group returns before the overlap check | FIX-NOW | algorithm | logic-flow | `8e8ca08` |
| 3 | 6084 | operator/kind gate preempts the `<clear>` check | FIX-NOW | checking | boundary | `5aef90e` |
| 4 | 6096 | int domain width overflows before saturation | KATA-BUG | algorithm | boundary | kata `x0fp` |
| 5 | 6083+6085 | `graph-vacuous-exists` ≠ RDR 0006's `graph-vacuous-atom`; blocking flags unset | KATA-BUG | interface | design-conformance | kata `fyf4` |
| 6 | 6094 | value-dimension unprojectability discards the presence dimension | KATA-BUG | algorithm | rare-situation | kata `x2bp` |
| 7 | 6085 | `eq`/`in` compare held values as raw strings | KATA-BUG | checking | boundary | kata `cq5p` |
| 8 | 6093+6092 | vacuous REQ-43 union guard; wall-clock ADV-2 oracle | KATA-BUG | test-oracle | test-coverage | kata `9yeq` |
| 9 | 6085 | `ambiguous_match` arm not vacuously closed when overlap-free | RDR-SEED ★ | interface | design-conformance | kata `r5ja` |
| 10 | 6096 | REQ-93 "known finite" vs REQ-94 "provable" | RDR-SEED ★ | interface | rare-situation | kata `ge67` |
| 11 | 6092 | adversarial tests committed red (CI-red) | DROP `superseded-at-HEAD` | n/a-not-a-defect | — | suite green at HEAD |
| 12 | 6092 | ADV-1 skips a green report excluding the empty view | DROP `superseded-at-HEAD` | n/a-not-a-defect | — | `Covers([])` = true |
| 13 | 6083 | exact 2^80 cardinality representation | DROP `superseded-at-HEAD` | n/a-not-a-defect | — | `skeleton.go` deleted |
| 14 | 6085 | saturated 2^40 instead of exact computed size | DROP `rdr-adjudicated` | n/a-not-a-defect | — | D6 fixes B=2048 |
| 15 | 6083 | REQ-66 assertion bans `info` severity | DROP `rdr-adjudicated` | n/a-not-a-defect | — | two-tier vocabulary is deliberate |
| 16 | 6083 | `shapePairSource` unequal cardinalities | DROP `rdr-adjudicated` | n/a-not-a-defect | — | D6 chose B=2048 for this |

★ = carries an `## Open question` resolved when the kata ships.

## The three FIX-NOWs

Each was mutation-tested — the fix reverted, the new test confirmed red,
the fix restored — so none is a test written to match the code.

1. **Held `null` (`072c7a0`).** `json.Unmarshal("null", &members)`
   returns NO error with a nil slice, while `[]` yields a non-nil empty
   one, so a malformed held value was DECIDED false. That is the
   direction `0007:C1` forbids ("never to false and never to true") and
   the mirror image of the defect D14 fixed. Pinned in the cross-RDR
   contract (`contains/null value`) so every evaluator inherits it.
2. **Dimensionless overlap (`8e8ca08`).** The only true false-green in
   the run: two guard-atom-free ordinary rows in one selection context
   produced NO verdict and NO finding — lint exited 0 on a model the
   kernel always refuses `ambiguous_match`. Withholding the exhaustiveness
   CLAIM is not withholding the overlap invariant, which carries no
   finitely-declared precondition.
3. **`<clear>` precedence (`5aef90e`).** RDR 0003's new operator/kind gate
   was inserted ABOVE the sentinel check, inverting `0002:C17` ("binds
   every atom regardless of operator"). The code comment asserted the
   opposite of the behaviour. Every pre-existing REQ-61 site pairs the
   sentinel with a matrix-LEGAL operator/kind, so none could catch it.

## A fix implemented and deliberately reverted

RDR 0006's rule that the `ambiguous_match` arm is "vacuously closed for a
group whose ordinary population is overlap-free" was implemented, and it
turned `TestReq77_CoverageUnionIsPerGroupTimesDeclaredRescuableClass` red.
REQ-77's own fixture has a single ordinary row — necessarily overlap-free —
so the two LOCKED records prescribe opposite outcomes on the same model.

The change was reverted and the conflict recorded as seed `r5ja` with the
tension stated. Picking a side would have silently overruled a locked
record; editing the inconvenient test would have been worse. The current
behaviour is a false RED (a spurious blocking finding), never a false
green, so deferring is safe.

## Note on the DROP population

10 of 20 findings dropped, 6 of them `n/a-not-a-defect` — mostly the
flow's own red-green rhythm (Phase 3b commits adversarial tests
deliberately red; Phase 3c fixes them) and cosmetic large-cardinality
complaints that D6's published bound of 2048 already adjudicates. Two
DROPs (12, 13) were superseded by work that landed later in the same
window.
