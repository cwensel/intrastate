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
| resolve (re-run) | PASS | no | All 6 rulings absorbed; A4 restated and Verified (5/5 CAs, rulings_open=0). New contract C5 (canonical int spelling at load). Problem Statement + MVV re-scoped onto the owned ingress; JDR 0004 JD-1/JD-2/JD-3 applied. Profile re-resolved foundational (idempotent). Commits ef5312d, 37bee1a. |
| prelock cove | PASS | no | Converged in 2 iterations, 16 findings all fixed. REFUTED a Key Discovery: graphlint's reachability already runs MATCH atoms through the guard seam, so "shared for guard atoms only" was false — design survives (C5 makes the divergence unauthorable) but the stated reason was wrong. C5's venue re-pinned to "the ingress that admits them". Commits d652204, e77f86f. |
| prelock 3amigo | PASS | no | Converged in 2 iterations, 27 findings. **C5 NARROWED** from the shared `conformKind` int arm to the predicate ingress (`normalize.go::(*loader).atom`) — the wider rule was refuted by two BOUNDARY-marked tests passing on main (RDR 0024 REQ-7/REQ-9). CLI + `[emit]`/`[initial]`/`[rule.write]` now OUT of scope; CLI int canonicality charted to a successor. A4 flipped Verified->Pending (basis changed); new Pending A6 (installed base of persisted owned values unmeasured). Commits 08919f4, ccd2c99. |
| prelock critique | PASS | no | Dual-model (opus-5 + fable-5, third-model barrier diff), 2 iterations, 23 findings. **REVERSED a Verified Key Discovery**: canonicalization does NOT close the divergence — graphlint reaches held values via `heldValues`->`canonicalValues` with no int round-trip, fed by `[initial]`/`[rule.write]` cells `conformKind` admits as `"007"`. Reopens the `matchlit` shape. C4 gained a consumption obligation (two-valued collapse at consumers would have made lint green on this RDR's own defects). Commits ab4aefa, 6198deb. |
| resolve | NEEDS_DECISION | yes | A1/A2/A3/A5 Verified against source; A4 REFUTED on its literal side (authored `n eq "00"` flips lint blocking→clean) and on its malformed leg (`--tag` conforms upstream, so the seam is never reached). Profile latched `foundational` (1 durable C1; C2/C3/C4 marked `Surface — of C1`). Evidence-body authored; Performance Expectations omitted per template. 3 questions + 3 fixtures put to the author in `evidence/author-round.md`; delegated run, so the round is not self-approved. Commits 60a3954, ca709d9. |

Re-ask after resolve (2026-09-20): profile `foundational` unchanged; lens row
`cove,3amigo,critique,repeatability` unchanged (rule `lens-foundational-cove`);
posture unchanged. Plan line above still current — no rewrite owed.

## OPEN DESIGN CHOICE — A4, for the user (critique declined to collapse it)

Critique established the divergence is real but did NOT select a remedy, stating
that is Stage 6's call. Two live options, both of which change a fenced contract:
  (i) extend canonicalization to the HELD ingress (`[initial]`/`[rule.write]`
      cells) — changes C5, and is close to what ruling Q1 originally said before
      3amigo narrowed it;
  (ii) make the seam's match arm byte-compare — changes C2, and is option (d)
      from the original Q1 fork, which the ruling rejected.
Note the shape: 3amigo narrowed C5 away from the held ingress on RDR 0024
grounding; critique then found the held ingress is exactly where the leak is.
The ruling's rejection of (d) rested on the premise that canonicalization could
close the divergence everywhere — critique reversed that premise.
This is a ruling-level question, not a stage-level one. Cascade stops after
repeatability rather than carrying it into reconcile.

## Ruling divergence — for the user, not resolved here

3amigo narrowed C5 to the predicate ingress. This is NARROWER than ruling Q1,
which said "every authoring site: guard literals, match literals, `[initial]`,
`[rule.write]`". The narrowing is grounded (two live BOUNDARY tests from locked
RDR 0024 refute the wider rule; `--outcome ground` returned `apply`), and the
stage charted the CLI/`[initial]`/`[rule.write]` scope to a successor rather than
retiring 0024's locked REQs — out of this RDR's mandate.

Consequence the stage states plainly: `--tag iter=07` now matches SILENTLY
(accepted residual, recorded in F1). That is the misroute the ruling intended to
close. The ruling's own rationale is also undercut: "zero violators" rested on a
wrong-unit measurement (S5 counted guard atoms; the rule was homed in an arm
serving five populations).

The user should decide whether the charted successor is an acceptable home for
the CLI scope, or whether 0012 must carry it. Not this skill's call.

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
  - Literature/prior-art pass at the ceiling on the (a)-vs-(d) fork — the evidence
    neither prior pass consulted. **Verdict: (a), canonicalize at load, widened.**
    - Grounding is principle + prior art, not agent judgment: the measured defect is
      a Non-Redundancy violation (Meyer, OOSC §11.7 pp.354-355) / SPOT violation
      (Raymond, AoUP p.124) — one correctness condition decided in two places with
      two answers. Google BSRS p.105 (parse-don't-validate), DDIA p.39
      (schema-on-write), Refactoring Databases p.186.
    - DX / least-burden answer is (a): AoUP Rule of Repair p.54 — under (d) the
      author gets lint exit 0 and a runtime dead end, the "much later" case; (d) has
      no natural discovery point. Raymond/Spencer p.53 on "almost the same" fits (d)
      exactly (match and guard atoms are visually identical, subtly different).
    - Peer CLIs agree: helm `typedVal` refuses leading-zero int coercion (commit
      609e72b35, issue #2693); OpenTofu converts at ingress with a sited diagnostic.
    - Timing: zero violators argues FOR acting now (SWE-at-Google pp.161-162, gofmt
      vs buildifier's 6 engineer-weeks over 200k files; Refactoring Databases p.186).
    - REFINEMENT: C4 stays as written — `atomAdmitsValue` remains a construction site
      — with the RDR stating lint's typed admission is safe ONLY because load made
      typed-vs-bytes unobservable. This DECOUPLES Q3 from Q1.
    - Refusal must be diagnostic (Postgres style guide p.2642) and cover the
      bare-float `-0.0` -> `"-0"` path, not just quoted strings.
    - Rejected third framings: canonicalize-on-read (K8s Quantity's own "or don't
      diff" caveat; models are reviewable data), accept-with-deprecation (zero
      violators), `--fix` (a later affordance layered on (a), not a substitute).
    - Corpora empty where expected: SchemaEvo* is relational-DDL-oriented,
      CodeMaintenance named no hazard, StateMachine*/SpecDrivenDev nothing on point.
      Grounding rests on DevRef + peer CLI source.
    Relayed to the user 2026-09-20. Advisory; the six items remain unruled.
