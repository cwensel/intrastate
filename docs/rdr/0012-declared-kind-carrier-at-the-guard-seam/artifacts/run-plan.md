# Run plan — /rdr-draft-to-lock 0012

What this run scheduled and decided. Not where the RDR stands — that is
derived (`/rdr-status 0012`).

```
rdr: 0012-declared-kind-carrier-at-the-guard-seam   profile: foundational   (as read; Draft = provisional)
posture: upfront=true each=false finalize=true
lenses: cove,3amigo,critique,repeatability
                                             (critique/repeatability fan out)
stages: refine -> resolve -> [lenses] -> reconcile -> finalize   stop-after: finalize
```

Lens row from `rdr-status.toml` `--outcome lens`, rule `lens-foundational-cove`:
"Profile foundational — cove leads the row (it subsumes grounding's sweep as
Step 0) and has not run." Accretion floor: `none` (rule `floor-below-two`).
Posture from `rdr-cascade.toml`, rule `posture-draft-foundational`.

Not a re-entry: `status_form=none`, `reentry_target=none`.

## Ledger

| Stage | Verdict | Blocking | Note |
| --- | --- | --- | --- |
| (up-front confirm) | — | — | plan put to user; go-ahead received |
| refine | PASS | no | 4 contradictions + 3 redundancies resolved; no change-history; commit 6c22a12. Router re-answered /rdr-refine (locate-draft-refine); advanced to Stage 4 on this run's Ledger PASS per the refine carve-out. |
| resolve | NEEDS_DECISION | yes | A1/A2/A3/A5 Verified; A4 refuted on two legs. Profile latched foundational. Commits 60a3954, ca709d9. PARKED on Stage 4's author's round (6 items: Q1-Q3 + fixtures F1-F3) — put to the user; rulings go to evidence/rulings.md. |
| resolve | NEEDS_DECISION | yes | A1/A2/A3/A5 Verified against source; A4 REFUTED on its literal side (authored `n eq "00"` flips lint blocking→clean) and on its malformed leg (`--tag` conforms upstream, so the seam is never reached). Profile latched `foundational` (1 durable C1; C2/C3/C4 marked `Surface — of C1`). Evidence-body authored; Performance Expectations omitted per template. 3 questions + 3 fixtures put to the author in `evidence/author-round.md`; delegated run, so the round is not self-approved. Commits 60a3954, ca709d9. |

## Parked forks

- **Author's round, Stage 4** (`evidence/author-round.md`, 2026-09-20) — 3 questions
  + 3 fixtures awaiting the author. Q1 (how A4 narrows) is load-bearing: it rewrites
  C2's parsed-comparison leg. Q2 may re-scope the Problem Statement. Q3 asks whether
  JDR 0004's JD-1/JD-2/JD-3 land. Hard stop — a delegated run may not approve fixtures.
  Resume: `/rdr-resolve 0012` once `evidence/rulings.md` carries the rulings.
