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
  - §strong-consult run once at the bound ceiling (`claude-fable-5-1`; Profile
    foundational, `status_form=none` so no prior `consult: strong`). ADVISORY only —
    it did not close the fork. Its verdict: Q1 → (a) canonicalize int literals at
    load, retire A4's proxy claim (it argues lint and runtime move together, so the
    measured flip is the correct verdict for the new runtime, not a hidden defect);
    Q2 → both edits, owned-reader ingress is the primary door, build the arm;
    Q3 → confirm JD-1/JD-2, confirm JD-3's substance but hold its form (dependency
    note, not a present site in fenced C4). Notes F1's expectation and F3's role
    both change under (a). Bottom line: approach stands, lockable after three edits.
    Relayed to the user 2026-09-20; the six items remain unruled.
  - Adversarial verification of the consult's load-bearing empirical claim
    (LINT-EQ-RUNTIME), run once at the ceiling. New work, not a second opinion:
    it tested a checkable claim rather than re-asking the design question.
    **VERDICT: CLAIM REFUTED**, with a measured counterexample.
    - Shared-seam premise holds for GUARD atoms, is FALSE for MATCH atoms: lint
      decides match atoms through the seam (`reach.go:502-513`), the kernel
      byte-compares them (`resolve.go:178`, via `model.go:367`). C1 requires that
      divergence while C4 lists `atomAdmitsValue` as a typed construction site.
    - Fixture `matchlit` (`[initial] n = 0`; `[rule.match.n] eq = "00"`): typed lint
      exit 0 / no findings, while the runtime is a real dead end at the root
      (`flow-no-match`, exit 2) under BOTH binaries. The dropped `graph-dead-end`
      was correct. Same from the held side (`matchheld`).
    - `cover07` itself survives — its clean verdict does re-describe the typed
      runtime. The narrow claim held; the generalization did not.
    - Consequences for Q1: option (b) is dead; option (a) is a SAFETY fix and is
      insufficient as scoped (must extend to int `[initial]` and `[rule.write]`
      values, not just guard literals); and a THIRD option the author's round never
      listed is on the table — exclude `atomAdmitsValue` from typed comparison so
      lint's match admission stays byte-equal to the kernel's, which contradicts
      C4 as written and so also bears on Q3.
    - Blast radius partially refuted: 220 .toml files (not 225); zero non-canonical
      spellings CONFIRMED; but "only authorable as a quoted string" is FALSE —
      bare float `eq = -0.0` renders literal `"-0"`, which `Atoi` accepts.
    - `--write n=07` misroute CONFIRMED by measurement (base prunes to `fallback`,
      typed selects `guarded`).
    - Corpus diff byte-identical across all 123 lintable models: the divergence is
      LATENT, not currently firing — an urgency argument, not a correctness one.
    Relayed to the user 2026-09-20. Still advisory; the six items remain unruled.
