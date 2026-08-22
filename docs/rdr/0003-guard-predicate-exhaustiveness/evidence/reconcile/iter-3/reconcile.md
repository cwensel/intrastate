Model: claude-opus-5[1m]

# Reconcile — RDR 0003 Guard Predicate Exhaustiveness (iteration 3)

## Verdict

**NOT RECONCILED — lock is blocked on A11 (with subfields A7, A9).**

The route is **not** a stage return inside this RDR. This RDR's approach,
wording, and evidence discipline are not what is defective — the missing field
belongs to `Draft` RDR 0002. The correct next action is the consolidated
producer request, then re-entry at reconcile.

## Preflight — Stage 5 completeness (the gate iteration 2 stopped at)

`Profile: large` → lens row `grounding → 3amigo → critique`, **plus
`repeatability` (lite) because the Determinacy trigger fires**.

| Lens | Required | On disk | State |
| --- | --- | --- | --- |
| grounding | yes (`large`) | `evidence/grounding/iter-2/{findings,dispositions}.md` | complete |
| 3amigo | yes (`large`) | `evidence/3amigo/iter-2/{persona-1..3,consolidation,dispositions}.md` | complete |
| critique | yes (`large`) | `evidence/critique/iter-2/{critique,critique-modelB,diff,dispositions}.md` — two distinct stamps | complete (dual-model) |
| repeatability (lite) | yes — Determinacy fired | `evidence/repeatability/iter-2/{run-1,diff,dispositions}.md` | **complete** |

Iteration 2 stopped here because `evidence/repeatability/` did not exist. It now
does. **Variant check passes**: `run-1.md` reads `variant: lite (profile: large)`
and no `run-2`/`run-3` sits beside it — lite was owed and lite ran. The
dispositions record deliberately declined the `diff.md` escalation
recommendation with a stated rationale (the findings are textual contradictions
inside the draft, not a sampling question), and correctly left the `variant:`
header unrewritten. That is a conforming record, not a variant mismatch.

**Preflight: PASS.** The open set was built.

## Absorption audit (delegated)

One sub-agent over the four iteration-2 lens output dirs, reporting per round
whether each finding was absorbed into the current RDR or survives as residue.

| Round | Findings | Absorbed | Residue |
| --- | --- | --- | --- |
| grounding | 1 | 1 | 0 |
| 3amigo | 22 raw → 10 ledger | 10 | 0 |
| critique (dual-model) | 20 ledger (R-01..R-20) | 20 | 0 |
| repeatability-lite | 7 (D-1..D-7) | 7 | 0 |
| **total** | **50 raw / 38 ledger** | **38** | **0** |

**Zero residue.** Both dismissals (R-19, D-7) carry cites that hold in the
current text. Every deferred finding is booked to a Critical Assumption.

## The open set

| # | Item | Source(s) | Disposition | Evidence pointer or plan |
| --- | --- | --- | --- | --- |
| A1 | Target flows fit a closed typed predicate vocabulary | 1 (critique demotion), 2, 3 | **VERIFIED** | `evidence/spikes/iter-2/a1-eval-harness/` — `go run .` → `output.txt`; 28 assertions pass / 0 fail |
| A2 | Exhaustiveness reduces to scoped finite declared domains | 1 (critique demotion), 2 | **BLOCKER** (behind A11) | Sound derivation, no producer. Closes when A11 lands |
| A7 | Existence atom projects as a presence dimension | 1, 2 | **BLOCKER** (A11 subfield) | RDR 0002 per-tag optionality field |
| A8 | RDR 0006 carries the same `guard_unevaluable` narrowing | 1, 2 | **DOWNGRADED** | Cluster reconcile 0003·0006·0007 + JDR 0001 §JD-4 assignment |
| A9 | Declared set-element universe exists for `contains` | 1, 2 | **BLOCKER** (A11 subfield) | RDR 0002 element universe |
| A10 | RDR 0006 accepts the row-group division of labour | 1, 2 | **DOWNGRADED** | Cluster reconcile, with A8/A12 |
| A11 | RDR 0002 carries a finite-domain field | 1, 2 | **BLOCKER — gates lock** | Consolidated four-field producer request to `Draft` RDR 0002 |
| A12 | Predecessor reachability is decidable | 1, 2 | **DOWNGRADED** | Cluster reconcile, with A8/A10 |
| A13 | Canonical set-literal spelling | 1 (repeatability closed it) | **ACCEPTED** (Design Decision) | `Normative Contracts` set-literal clause; rejected alternative recorded |
| A14 | Normalization preserves each atom's `block` | 1, 2 | **DOWNGRADED** | Rides the A11 four-field request |
| A15 | A single cardinality bound is model-independent | 1 (repeatability opened it), 2 | **DOWNGRADED** | MVV Scenario 3 equal-cardinality/differing-shape comparison |
| S-1 | `evidence/spikes/check.sh` (named spike, recorded negative) | 3 | **SUPERSEDED** | Replaced by the A1 evaluation harness; retained as the historical fixture |
| E-1 | Exactness-word delta over the round edits | 4 | **COVERED** | See below |

### Source 4 — exactness-word delta

Swept the round-introduced text (`8b76d62..a52055a`, the four iteration-2 lens
edits). Every new or edited exactness claim resolves to a normative clause with
a covering MVV scenario or an Evidence Record:

| Claim introduced by the rounds | Cover |
| --- | --- |
| participation — "every key any row carries, whether or not it discriminates" | `Normative Contracts` participation clause; MVV Scenarios 2, 3; A11 sizes the producer request to it |
| report-every-defect — "every defect it can decide in one pass", "a finding for each" | report-every-defect clause; MVV Scenarios 2, 3 |
| set-literal canonicalization — "one canonical spelling", "MUST parse to the same literal" | set-literal clause; **A13** (`Verified`); MVV Scenario 5 |
| bound determinism — "the same model MUST receive the same verdict on every conforming implementation" | too-large clause; **A15** (`Pending`, MVV Test); MVV Scenario 3 |
| gap attribution — "every source rule id in that group and at least one concrete uncovered assignment" | diagnostics clause; MVV Scenario 2 |
| "can refuse" totality — "decide this syntactically over declarations so the test is total" | narrowing clause; **A12** books the reachability input |
| escape-row participation — "participates in the coverage identity like any other row" | escape-row clause; A10 books the peer agreement on grouping |

No round-introduced exactness claim is uncovered.

## A1 — the one item this gate closed

A1's plan gave two routes: re-run the spike as an evaluation harness, or carry
it into implementation as the MVV's first acceptance gate. The first was taken.

The prior `check.sh` is an awk scan over operator *spellings* in TOML section
headers: its `coverage=` line is a hardcoded `printf` constant, and its report
loop is hardcoded to `split("eq in lt gte exists")`, so `lte`, `gt`, and
`contains` could not appear in its transcript even when used. It evaluates no
predicate.

The harness (`evidence/spikes/iter-2/a1-eval-harness/main.go`) encodes the
fixture's four target-flow rules as normalized atoms, adds a `contains` row over
a set-valued tag and a `gt`/`lte` band row, and **evaluates** them against
eleven representative tag views with asserted per-row qualify/prune verdicts.

- 28 assertions pass, 0 fail.
- All eight operators of the closed set carry an evaluated atom.
- Five negative controls confirm the vocabulary is closed *and typed*: unknown
  operator, `lt` on an enum, `contains` on an enum, literal/kind mismatch, and
  `in` with a scalar literal are each rejected at well-formedness.
- The seam split is reproduced, not assumed: `exists` decides on presence alone
  (absent `cluster_eligible` → decided FALSE), while a value atom over an absent
  key yields UNEVALUABLE under RDR 0007's veto.

**Scope discipline.** The harness consumes no declared domain — it reads only
`[[rule]]` blocks and supplied tag values, deliberately ignoring the fixture's
invented `domain`/`min`/`max` keys, so it does not re-create the circularity A2's
stamp was corrected for. It establishes evaluation, not exhaustiveness. A9 is
therefore **not** closed by the `contains` cell: evaluation needs the tag's
kind, an exhaustiveness claim needs the declared element universe.

## Why lock is blocked — the second hard rule

> No spike or assumption deferred past lock if the MVV depends on it.

The RDR's own **authorability gate** states that several MVV cases cannot be
written until their producers land: the `contains` case needs A9's element
universe, the possibly-absent-key case and its always-present negative control
need A7's optionality field, and every finite-domain case needs A11's domain
field.

Verified against source rather than taken from the draft's own prose: RDR 0002
`Technical Design` part 2 declares **tag name, provenance, value kind, and an
optional accessor reference** — and its tag normative contract adds only
provenance. No domain, bound, optionality, or element-universe field exists
anywhere in that document, and RDR 0002 books no reciprocal obligation to
supply one.

So A11 (with A7 and A9) pins the very cases the Minimum Viable Validation
consumes. DOWNGRADED is unavailable: locking here would freeze a contract whose
only falsifier is unrunnable, which is exactly what this gate exists to prevent.

## Strong consult (rdr-common §strong-consult)

One fresh-context consult at the strongest available tier, briefed on the two
claims in tension and the new evidence, without prior verdicts.

```
verdict: BLOCK
blocking: yes
next_action: Route RDR 0002 (Draft, already in refine re-entry) to record the
  consolidated A11 finite-domain tag-declaration field with A7 optionality and
  A9 element-universe subfields; then re-run RDR 0003 reconcile. Take
  A8/A10/A12 to cluster reconcile (0003·0006·0007) with a JD-4 disposition.
  No stage return inside RDR 0003.
```

It independently confirmed the field census against RDR 0002's source, and
added two points this report adopts:

1. RDR 0002 is not merely `Draft` — it is in an active refine re-entry, so the
   producer request has a live landing site. Claim B's "holding 0003 unlocked
   freezes both documents" does not hold.
2. The defect is not in this RDR's approach (Stage 2), wording (Stage 3), or
   evidence discipline (Stage 4). Returning 0003 to any of those stages would
   ask it to fix something it does not own. **The route is a peer action.**

## Refutation check

**No spike or source search refuted an assumption this RDR relies on.** The A1
harness *confirmed* the operator vocabulary rather than refuting it, and the
RDR 0002 field census confirmed a gap the draft already books rather than
contradicting a claim the draft makes. Nothing reopens a completed stage, so
**no §punt-ledger row is owed by this run** — the escape iteration 1 made
(stamping A1 `Verified` on the spelling-scan spike) was already caught and
recorded by the iteration-2 critique pass.

## Completeness check

- No `_Draft placeholder._` in any body section — grep clean.
- No `this is a seed skeleton` header — grep clean.
- `## References` carries real citations (peer RDRs, JDR 0001 §D1/§D4/§JD-4,
  the output contract, the prior-art corpus), no template brackets.

## Dispositions landed in the RDR

Every row above is written into the RDR's Critical Assumptions section, not only
here. A1's record is rewritten to `Verified` with the harness command and output
path; A2/A7/A9/A11 each carry a **Stage 6 disposition — BLOCKER** paragraph
naming why deferral is unavailable; A8/A10/A12/A14/A15 each carry a **Stage 6
disposition — DOWNGRADED** paragraph naming the venue, the survivability
argument, and the plan that will run.

Amendment sweep for the A1 flip: Metadata `Status`, Prerequisites checklist,
Phase-gating table (Phase 3 no longer names A1's harness as a blocker),
Assumption Verification prose, and the Testing Strategy spike citation were all
re-pointed. `grep` for stale "A1 is Pending" phrasing returns nothing.

## Next

```
/rdr-refine 0002        # land the four-field producer request (A11 + A7 + A9 + A14)
/rdr-cluster-reconcile  # 0003 · 0006 · 0007 — A8, A10, A12 and the §JD-4 assignment
/rdr-reconcile 0003     # re-entry once RDR 0002 records the domain field
```

The two peer actions are independent and can run in either order. This RDR is
otherwise ready: its grammar, identity tuple, `unless` semantics, coverage
derivation, and diagnostics contract are stated, internally consistent, and
carry no unabsorbed review finding.
