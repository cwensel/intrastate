Model: claude-opus-5[1m]

# Reconcile — RDR 0003 Guard Predicate Exhaustiveness (iteration 2)

## Verdict

**NOT RECONCILED — return to Stage 5.**

Stopped at the Stage 5 completeness preflight. The open set was **not** built:
the skill's step 1 makes the preflight a gate, and a reconcile run over an
incomplete lens row would disposition assumptions the missing lens is meant to
disturb — the same ordering failure that let iteration 1 pass.

## Preflight — Stage 5 completeness

`Profile: large` → lens row `grounding → 3amigo → critique`, **plus
`repeatability` (lite) when the Determinacy trigger fires** (rdr-common
§lens-row: the trigger appends a row entry, not a footnote).

| Lens | Required | On disk | State |
| --- | --- | --- | --- |
| grounding | yes (`large`) | `evidence/grounding/iter-2/{findings,dispositions}.md` | complete |
| 3amigo | yes (`large`) | `evidence/3amigo/iter-2/{persona-1..3,consolidation,dispositions}.md` | complete |
| critique | yes (`large`) | `evidence/critique/iter-2/{critique,critique-modelB,diff,dispositions}.md` — two distinct stamps (opus-5 / sonnet-5), diffed behind a barrier | complete (dual-model) |
| **repeatability (lite)** | **yes — Determinacy fired** | **absent** | **MISSING** |

`evidence/repeatability/` does not exist, in iteration 1 or 2, and no written
`determinacy: n/a — <reason>` disposition exists anywhere under the RDR's
evidence tree (`grep -ri determinacy` → 0 hits). The obligation is therefore
undischarged by either of its two allowed routes.

### Determinacy trigger read

Read from the fenced ` ```normative ` blocks only (`0003-…md:554-700`), per the
Stage 5 narrow-read rule — not a whole-RDR grep. Confirmed by an independent
sub-agent pass. **FIRES**, on four disjunctive cues:

1. **Parse/deparse.** "Unknown operators MUST be rejected during parse or lint
   before resolution"; "A predicate whose literal cannot be parsed as the
   declared tag kind MUST be rejected before resolution." A grammar contract.
2. **Identity / compose-decompose.** The escape-row clause legislates a
   *coverage identity* — "its accepted assignments are computed from its guard
   atoms and unioned with its peers'" — over a product the adjacent clause
   requires be evaluated "as one product rather than as independent
   one-dimensional checks."
3. **Data-model field whose ownership/semantics could be inferred more than one
   way** — the loudest cue, and self-evidenced. The RDR's own `authority` census
   marks the row-group exhaustiveness verdict **"Contested — A8"** (two
   candidate writers: this RDR vs RDR 0006) and marks per-tag optionality and
   the `contains` element universe **"not yet declared."**
4. **Model-independent determinism.** "'Too large to prove' MUST be a declared,
   model-independent bound… the same model MUST receive the same verdict on
   every conforming implementation"; "MUST be proved with a deterministic
   symbolic or bitset-equivalent representation." A clause legislating
   cross-implementation verdict identity is precisely what an alternate-model
   reconstruction diff tests.

Counterweight noted and rejected as dispositive: "Successful row matching MUST
NOT depend on source order or first-match priority" negates the *step-ordering*
cue by design. The trigger is disjunctive; three other cues fire independently.

No clause carries a `Transient` marker, so nothing in the block is excluded from
the trigger read.

## Why this is not waivable here

Two independent grounds, either sufficient on its own:

1. **The lens is owed and undischarged.** §lens-row: "Whole row complete →
   `/rdr-reconcile`." The row is not complete. Reconcile cannot mint the missing
   lens's disposition, and the `determinacy: n/a` escape is a Stage 5 artifact,
   not a Stage 6 one.
2. **The trigger's subject matter is exactly what is contested.** The three
   assumptions the RDR itself names as lock-blocking (A7, A8, A11) are all
   ownership/producer questions over the guard data model — cue 3. Running a
   reconcile that dispositions A8 as "route-back at cluster reconcile" without
   first taking an alternate-model reconstruction of the contract block would
   disposition the very ambiguity the lens exists to surface.

## Notes carried forward (not dispositions)

Recorded so the Stage 5 re-entry and the subsequent reconcile do not re-derive
them. **These are observations, not terminal dispositions** — no item below has
been dispositioned, and the open set remains unbuilt.

- **Open-set size, preliminary.** Ten records are `Pending`: A1, A2, A7, A8, A9,
  A10, A11, A12, A13, A14. Sources 1 and 2 overlap almost completely — the three
  iter-2 lens dispositions carry A1, A2 (demoted from `Verified`), A6 (narrowed),
  and A7–A14 (new or carried).
- **A1's spike is a recorded negative.** `evidence/spikes/check.sh` is an awk
  scan over operator *spellings* in TOML section headers; its `coverage=` line
  is a hardcoded `printf` constant and its report loop is hardcoded to
  `split("eq in lt gte exists")`. It evaluates no predicate. Source 3 is
  therefore "named spike exists but does not support the claim," not "unrun."
- **Iteration 1's reconcile is the escape this gate exists to catch.**
  `evidence/reconcile/reconcile.md` stamped that same spike **VERIFIED** and
  closed with "No Critical Assumption remains Pending or Unverified." The iter-2
  critique pass refuted it. A punt-ledger row is owed for that escape at the
  route-back that reopens it (rdr-common §punt-ledger) — filed under the Stage 5
  re-entry, not here, since this run reopens no completed stage's decision.
- **MVV-critical deferral, pre-flagged.** The MVV's own authorability gate
  states the `contains` case, the possibly-absent-key case, and its negative
  control are **unauthorable** against RDR 0002's current tag declaration
  (A9, A7, A11). The prompt's second hard rule — no assumption deferred past
  lock if the MVV depends on it — is live and will bind the reconcile that
  follows Stage 5. It is flagged, not dispositioned.
- **Body completeness (grep, clean).** No `_Draft placeholder._`, no seed-skeleton
  header, no surviving template bracket. `## References` is filled with real
  citations, not bracketed placeholders. These would not have blocked.

## Next

```
/rdr-prelock 0003 repeatability
```

Discharge by one of exactly two routes:

- run the **repeatability-lite** lens (one alternate-model reconstruction of the
  Normative Contracts block + a focused diff), escalating to the full lens on the
  criteria the lens names; **or**
- write the disposition `determinacy: n/a — <reason>` into
  `evidence/repeatability/`, which given the four cues above would need to rebut
  all four.

Then re-run `/rdr-reconcile 0003`.
