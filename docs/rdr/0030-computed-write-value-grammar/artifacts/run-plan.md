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
