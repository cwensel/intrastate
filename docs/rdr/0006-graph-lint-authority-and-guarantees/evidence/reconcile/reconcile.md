Model: claude-opus-5[1m]

# Reconcile — RDR 0006 (post-demotion cycle)

Supersedes the reconcile report of the pre-demotion lifecycle (RDR 0006 was
`Final`, demoted to `Draft` on 2026-08-21 by the JDR 0001 rehoming). That report
predates all six rounds below and was regenerated, not read.

## Stage 5 preflight

`Profile: large` → lens row `grounding → 3amigo → critique` (§lens-row). All
three have completed evidence: `grounding/` (iter-1, iter-2), `3amigo/` (iter-1,
iter-2, iter-3), `critique/` (iter-1, iter-2 dual-model with `diff.md`,
`strong-consult.md`, and a delta pass). Determinacy trigger did **not** fire —
`determinacy` appears nowhere in the RDR or the round output, so no
`repeatability` row entry is owed. Stage 5 complete.

## Open set

| item | source(1-4) | disposition | evidence pointer or plan |
| --- | --- | --- | --- |
| **A5** CI runs `intrastate lint` as blocking graph-acceptance authority | 1, 2 | DOWNGRADED | Plan: Validation scenario 7 — the `graph-lint` job in `.github/workflows/ci.yml` asserting the JSON `code`. Gate surfaces re-verified: `.github/workflows/ci.yml::jobs` = `test`/`lint`/`vuln`, `lint` is `name: Lint` (golangci-lint + `make fmt-check`), none invokes `make check`; `Makefile::check` = `fmt-check vet lint test`, no `build` edge; `Makefile::build` → `./bin/intrastate` from `./cmd/intrastate`. |
| **A6** model declares initial owned state + terminals | 1, 2 | DOWNGRADED | Plan: additive clause on RDR 0002's authoring schema (`Draft` — scheduled peer edit). Absence re-verified: `initial`/`terminal` absent from RDR 0001 and RDR 0003; in RDR 0002 only the rule id `terminal-archive`; both fixture `[model]` blocks carry only `id`, `version`, `description`. Window is loud — a rootless model takes blocking `graph-dangling-edge`. |
| **A7** invariant 5 + `graph-always-present-owned` discharge the owned half of `0003::A18` | 1, 2 | DOWNGRADED | Plan: route-back at RDR 0003's next touch. `0003::A18` is `Status: Pending`, itself **DOWNGRADED at RDR 0003's own Stage 6** as not lock-blocking, venue "RDR 0007's next touch and/or RDR 0006's refine — either discharges it". `0003` Status line lists eight open records, "none lock-blocking". |
| **A8** a transition model is checked into this repo | 1, 2 | DOWNGRADED | Plan: authored + homed at implementation (now booked in `Prerequisites`). Re-verified: complete `.toml` set is `.roborev.toml`, `.kata.toml`, `0002/evidence/spikes/{rdr,kata}-fixture.toml`, `0003/evidence/spikes/guard-fixture.toml`. Gates scenarios 7 and 23 only. |
| **A9** `--flow <id>` resolves via RDR 0005 config discovery | 1, 2 | DOWNGRADED | Plan: flips when RDR 0005 lands discovery. Re-verified: `0005` is `Final` and its audit row reads "Parser placeholder; no transition-model config yet. \| Extend later". `--model <path>` is the instantiable form — now recorded in `Prerequisites`. |
| **A10** write-replaces for single-valued tags | 1, 2 | DOWNGRADED | Plan: RDR 0002 states the rule (`Draft` — scheduled peer edit, now in `Prerequisites`). Re-verified: RDR 0002 states only "Absence from both the write block and the clear list MUST NOT imply deletion"; no replace/accumulate rule anywhere; RDR 0003 owns the marker only. |
| **Stale A5 disjunction in `Prerequisites`, `Phase 3`, `Cross-Cutting Concerns`** | 4 | VERIFIED (fixed) | §amendment-sweep miss. The 3amigo H5 + critique iter-2 fix pinned A5's gate to one target and A5's Evidence says "The disjunctions ... are resolved", but three sites still carried "`make check` **or** the GitHub workflow" / "model **or** fixture corpus". All three rewritten to the pinned target. |
| **Misquoted locked peer in invariant 6** | 4 | VERIFIED (fixed) | The RDR quoted RDR 0007 as "guard-input decidability is not `RequiresOwned`'s job"; `0007` says "Guard decidability is not `RequiresOwned`'s job: guard-input coverage is enforced by the domain rule above". Substance held, quote did not. Corrected to verbatim. Sole body site (§amendment-sweep). |
| **A9/A10 absent from `Prerequisites`** | 4 | VERIFIED (fixed) | Both were booked as assumptions at critique iter-2 but never surfaced as implementation prerequisites. Added, with A8. |
| Named-but-unrun spikes | 3 | n/a — none | No assumption in this RDR uses `Method: Spike`. The six `spike` hits are the method-vocabulary boilerplate and citations of *other* RDRs' fixtures. The absent `evidence/spikes/` dir is correct, not a gap. |
| Exactness-word delta (rounds-introduced claims) | 4 | VERIFIED | Each round-introduced exactness claim carries MVV coverage or an inline derivation: per-invariant soundness split → scenario 20 + derivation (`Load-Bearing Decisions`); node-splitting exactness → scenario 21; published node ceiling → scenario 22; escape-row self-loops → grounded in `0002::Normative Contracts`; merged-node fixpoint no-false-green → inline derivation; deterministic finding order + canonical-not-hash fingerprint → scenario 4 and `Cross-Cutting Concerns`. |

## Absorption audit

Delegated over all six round dirs; every claimed-`fixed` disposition was
spot-checked against current RDR text.

- **grounding iter-1** — residue: none (zero refuted/missing source claims).
- **grounding iter-2** — residue: none. F1-F4 all verified in body; five
  dismissed-with-cite items each carry a peer or source citation.
- **3amigo iter-1** — residue: none (all four themes verified).
- **3amigo iter-2** — residue: none.
- **3amigo iter-3** — residue: none. H1-H8 and S1-S16 all verified in body.
- **critique iter-1** — residue: none. C1-C3 verified.
- **critique iter-2 + delta pass** — residue: none. All 23 reconciled defects
  verified, including the quantifier split of the no-false-green claim and the
  delta pass's three self-caught defects (`node-not-exact` correctly absent,
  dead-end reclassified universal, scenario 22 added with 23 renumbered).

What the round-level audit could not catch — and this pass did — is the four
source-4 items above: three sites the amendment sweep missed and one
misquotation. Absorption was complete *per finding*; propagation was not.

## Completeness

- No `_Draft placeholder._` body section, no seed-skeleton header, no surviving
  template bracket.
- `## References` is filled from citations the RDR already carries.
- A1-A10 sequential; A1-A4 `Verified`, A5-A10 `Pending` with a Stage 6
  DOWNGRADED disposition written into each record.

## Hard rules

- **Refutation** — none. Every claim re-verified against the working tree was
  found exactly as recorded (delegated verification, seven claim groups, all
  VERIFIED). No route-back; no §punt-ledger row owed.
- **MVV-critical deferral** — none fires. The Minimum Viable Validation is
  **fixture-backed**: it proves one legal model and one illegal model per
  blocking invariant class against authored fixtures through the production
  command path. It consumes no shipped model (A8), no CI wiring (A5), no
  `--flow` discovery (A9), and neither peer clause (A6, A10). A10 governs
  whether invariant 5's check is *complete*, not whether it runs.

## Verdict

**RECONCILED.** All ten items terminal, no BLOCKER. Six assumptions DOWNGRADED
with named, survivable plans; four source-4 defects found and fixed in place.

Next: `/rdr-finalize 0006`
