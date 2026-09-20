# rdr-draft-to-lock run plan — cli/0030

What this run scheduled and decided. Not where the RDR stands — that is
derived from the record (`/rdr-status 0030`).

```
rdr: 0030-computed-write-value-grammar   profile: foundational   (latched at Stage 4; was large/provisional)
posture: upfront=true each=false finalize=true
lenses: cove,3amigo,critique,repeatability
                                         (critique/repeatability fan out)
stages: refine -> resolve -> [lenses] -> reconcile -> finalize   stop-after: finalize
```

## Ledger

| Stage | Verdict | Blocking | Note |
| --- | --- | --- | --- |
| refine | PASS | no | One load-bearing contradiction resolved (audit row reused 0013 lint ceiling at load); 4 round-narration passages cut. All 9 CAs Pending. Commit 586ecba. |
| resolve | NEEDS_DECISION | no | 6 CAs Verified, A8 Refuted, A5/A6 Pending-at-MVV. A7 spike ran. Profile latched large->foundational; lens row rewritten. 4 author items parked -> evidence/author-round.md. Commits 61f4cea, 2340098. |
| resolve (re-run) | PASS | no | Q1-Q4 rulings absorbed; rulings_open 1+ -> 0. Profile latched foundational. C3 narrowed, 4 audit rows, 1 phase step, 2 test scenarios. Commits 408722a, cb51995. |
| prelock cove | PASS | no | Dual-model fan-out, converged iter-2. 26 findings, 25 fixed. Two blocking: stepped literal hit one of two write carriers; composition rule contradicted the atom sort. A10-A12 booked Pending. Commits 85e4ff5, 71f9d1b. |
| prelock cove (delta) | PASS | no | Author ruling: C13 narrowing upheld, no JDR (0002 is Implemented). Unrecorded override was the defect -- Overrides now names 0002:C13 beside C4; typed edge resolves. Commits 40daf9e, 7c74b71. |
| prelock 3amigo | PASS | no | 3 isolated personas + loop-ordered delta re-run; converged iter-2. Retracted 2 self-contradictory/invented claims in C1/C2; D-identity ordering corrected; S8-S10 added. A13/A14 booked Pending. Commits b037316, 75cbdc4. |
| prelock critique | PASS | no | Dual-model (Opus 5 / Sonnet 5) + 3rd-model barrier diff; converged iter-2. New C3 obligation: graph-lint findings publish the cell. Delta re-run caught that the iter-1 fix would silently false-green groupHasOverlap. 3 items charted. A15 booked. Commits 84af232, 67ff621. |
| prelock repeatability | PASS | no | FULL variant: 3 cross-model runs (opus/fable/sonnet) + barrier diff; converged iter-2, 7/7 pins clean. C1 site-split pinned (renderWrites owns cell walk; expand stays total). C3 join key single-sourced. Determinacy written: fired. A16 booked. Commits 4d9b7e0, 1d8f72f. |

Lens row complete (`emit.row: none`). Nine CAs Pending with declared Methods — Stage 6 closes them.
