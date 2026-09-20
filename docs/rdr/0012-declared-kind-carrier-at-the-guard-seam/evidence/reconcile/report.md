Model: claude-opus-5[1m]

# Stage 6 Reconcile — RDR 0012

Preflight: `--outcome lens` → `/rdr-reconcile` (row complete, Profile
foundational, `cove → 3amigo → critique → repeatability`).
`--outcome critique` → `none` (two base models, diffed).
`--outcome repeatability` → `none` (variant `full`; run-1/2/3 + diff).
No caveat on either `surface`. `spikes_unrun=[]`.

A new author ruling (`evidence/rulings.md`, `## 2026-09-20 — prelock`)
post-dated the lens row and re-ruled A4 as option (ii). It was read
before the record and applied in this pass; it named two verifications
owed HERE rather than assumed, both discharged below.

## Open set and dispositions

| Item | Source | Disposition | Evidence / plan |
| --- | --- | --- | --- |
| A4 — lint/runtime agreement on `eq`/`in` over `int` | 1, 2 | VERIFIED | Scoped by the ruling to the GUARD path; the held-side leak is match-path only and no longer typed. `assignment.go::valueAssignments` (`strconv.Itoa`), `spikes/a4-render-path.md`, `spikes/a4-lint-diff.md`; match-path derivation §A4 below |
| A6 — no silent verdict flip in the installed base | 1, 2 | VERIFIED | Census run at reconcile; population is empty. `models/rdr.toml::[tags.gate_passed]`, `models/examples/release-grammar.toml::[tags.risk]`/`[tags.urgent]`; §A6 below |
| A7 — three-valued widening of graphlint's atom check | 1, 2 | VERIFIED (as RETIRED) | Subject absent from the design after C4's narrowing; `atomVerdictForValue` appears nowhere in the record. §A7 below |
| A3 — construction sites have the model in scope | 1 (cove carry) | VERIFIED (unchanged) | Headline and measurement stand (three `Evaluator{}` sites on `main`). Only the cost sentence was corrected: C4 migrates two, allow-lists the third |
| C3 — second implementer + re-keyed suite cases | 1 (cove carry) | DOWNGRADED | Claims about test code that does not exist yet; Phase 2 implementation obligations, pinned by S1/S3 and MVV step 4. Survivable: a miss fails `TestGuardEvaluatorContract` loudly, not silently |
| C5 — CLI reach of the canonicalization | 1 (cove carry) | DOWNGRADED | The diagnostic being identical at both venues is an implementation obligation, verified when C5 lands (S7, MVV step 5). Survivable: a divergent diagnostic is a wording defect, not a verdict defect |
| Exactness sweep — "canonical" ×9 in C5 | 4 | ACCEPTED | All nine sit in C5, the clause that DEFINES the term (`strconv.Itoa(n) == authored`), and are covered by fixture S7 + MVV step 5. Not a post-mutation delta: the rounds did not touch C5's definition |
| Absorption audit — 4 lens rounds | 3 | PASS (no residue) | Delegated sweep over `3amigo/`, `critique/`, `repeatability/`, `cove/` (incl. `iter-2/`): every finding fixed or dispositioned; one CLI-canonicality item charted to a successor and recorded in F1 |

No BLOCKER. No refutation.

## §A4 — the ruling's first owed verification: MVV coherence under (ii)

(ii) narrows what C2 covers, so the MVV had to be re-checked for
coherence and for whether it still witnesses the parsed-comparison leg.

Verified coherent. Every MVV step is a guard-path or load-path
assertion; none touches a match atom or `atomAdmitsValue`:

- Steps 1-3 author a row GUARDED `iter eq 7` plus an owned accessor.
  The value travels `accessor::OwnedSnapshot` → `resolve::assemble` →
  `evaluateAtom` → the typed seam. Guard atoms, so C2 applies.
- Step 3 IS the parsed-comparison leg (`07` against `eq 7`, matching
  where raw-string comparison pruned). It remains witnessed, on the
  guard path, through the owned door — the one ingress load cannot
  reach.
- Step 4 is the conformance suite (C3); step 5 is C5's predicate-
  literal refusal. Neither is affected by the narrowing.

The narrowing therefore removes no MVV step. A note stating this was
added to MVV step 3.

## §A4 — the match-path derivation (moved here from A4's Evidence)

Both callers of `reach.go::atomAdmitsValue` are reached ONLY through
`BlockMatch` atoms, which is why excluding it from typing keeps lint
byte-equal to the kernel rather than merely differently-typed:

- `reach.go::matchSatisfiable` skips any atom where
  `a.Block != table.BlockMatch`, then calls `ownedAtomSatisfiable` →
  `atomAdmitsValue`.
- `analysis.go::nodeMeetsAll` (via `satisfiesSomeTerminal`) walks
  `a.model.Terminal`. Those atoms are dereferenced context predicates:
  `load.go::loadTerminal` → `contextAtoms` → `loadContexts`, which
  builds them with `atomsFromBlock(ctx.Match, BlockMatch, …)`.

So no guard atom reaches `atomAdmitsValue`, and the kernel byte-
compares the same population at `resolve.go::TagSet.matches`. The
guard path is disjoint and typed: `product.go::valueSatisfies` under
`Denotation`, with `product.go::selectionOf` routing `BlockMatch`
atoms away from the guard product.

## §A6 — census (narrowed to the guard path by the ruling)

An owned value can only reach a typed comparison through a tag that is
BOTH `provenance = "owned"` AND declared `int` or `bool`. Across the
committed models exactly one such dimension exists:

- `models/rdr.toml::[tags.gate_passed]` — `owned`, `bool`. Every
  authored value is the canonical token: `[initial] = "false"`, guards
  `eq = "true"` and `eq = "false"`, one `[rule.write] = "false"`.
- `models/examples/release-grammar.toml::[tags.risk]` (`int`) and
  `::[tags.urgent]` (`bool`) — both `provenance = "observed"`, which
  carries no persisted owned value.

Canonical `bool` tokens compare identically under raw-string and token
equality. The population of non-canonical-yet-parseable held values on
a typed dimension is EMPTY on `main`, which is exactly A6's claim.
Residual risk (an out-of-corpus deployment artifact) is carried by
S9's owned-reader fixture.

## §A7 — the ruling's second owed verification: retired, not re-scoped

The ruling asserts A7 ceases to exist; that was verified rather than
taken as given.

A7 quantified over a CHANGE: widening `atomAdmitsValue`'s return from
`bool` to `resolve.GuardResult`. The test for "retired vs re-scoped"
is whether any residual claim survives at a narrower home. None does:

- The narrowed C4 leaves `atomAdmitsValue` at its `main` shape
  (`bool`, `guard.Evaluator{}`, byte comparison).
- No other clause proposes touching it. `atomVerdictForValue` now
  appears nowhere in the record outside A7's own retirement note and
  the decision recording the reversal.
- The population A7 quantified over — lint verdicts changed by the
  widening — is empty BY CONSTRUCTION, not by measurement, because
  the widening is not performed.

That is retirement. A re-scoped assumption would still have a verdict
population at a narrower home; this one has none. The opposite-polarity
hazard A7 named is preserved in C4 as the stated reason the widening
is declined, so the knowledge outlives the assumption.

Status is `Verified` with `Method: Design Decision` — `Retired` is off
the CA vocabulary (`recs` flagged `ca_off_vocabulary_ids=["0012:A7"]`
when tried), and the prompt's ACCEPTED disposition (a design decision
with the rejected alternative named) is the correct shape.

## RDR edits made this pass

- C2 — opens by stating the guard/match split explicitly, as the
  ruling requires, pointing at §Load-Bearing Decisions for the
  reasoning and residual.
- C4 — construction sites narrowed from three to two; new
  "THE MATCH PATH IS EXCLUDED, DELIBERATELY" paragraph with the
  `BlockMatch` derivation; consumption obligation scoped to the guard
  path; `atomVerdictForValue` signature block removed and replaced by
  the stated rejection of the widening; dependency line "three-site"
  → "two-site".
- C1 — "three C4 sites span three packages" → two/two (×3 sites).
- A3 — Verified evidence untouched per the ruling; only the closing
  cost sentence corrected (two migrated, graphlint allow-listed).
- A4, A6, A7 — dispositioned as above.
- §Load-Bearing Decisions — the "open choice" passage replaced by the
  settled decision "The guard/match split, and why typing stops at the
  guard path", carrying the (i)-vs-(ii) reasoning, the effort argument,
  the Non-Redundancy answer, and the Residual + kata pointer. The
  earlier rejection (c) is marked REVERSED with a pointer.
- S4 (Validation scenario 4) — the blanket "no zero-value construction
  survives" assertion gains ONE named allow-listed exception
  (`reach.go::atomAdmitsValue`), asserted by name so the exception is
  as auditable as the rule; closing rationale re-aimed at guard-path
  sites.
- MVV step 3 — coherence note (guard path, unaffected by the exclusion).
- Mini-check tables (`fidelity`, `oracle`, `authority`) and the
  Existing Infrastructure Audit row — counts and the match-atom row
  corrected.
- Phase 3 — graphlint removed from the migration list, with the reason.
- §References — filled from citations the record already carried
  (peer records, source paths, spike artifacts, prior art, katas);
  the template's bracketed placeholders are gone.

## Kata filed

`intrastate#ch99` — the held-ingress canonicalization residual
(`[initial]`/`[rule.write]` → `heldValues` → `canonicalValues`, no int
round-trip), carrying the pre-lock critique evidence and the RDR 0024
`REQ-7`/`REQ-9` conflict as the reason it is out of scope for 0012.
Recorded in A4, §Load-Bearing Decisions and §References.

## Gate

- Open set complete (all four sources; source 3 empty, source 4 swept
  and dispositioned).
- No refutation → no return stage, no §punt-ledger row owed.
- DOWNGRADED items (C3, C5 implementation obligations) are survivable
  and named, each with the test that will run.
- Spikes: no new spike run this pass; the census was a read-only
  source search over committed models, so no worktree was owed.
  `recs clean` exits 0.
- Dispositions landed IN the RDR: `ca_total=7 ca_verified=7
  ca_pending=0 ca_off_vocabulary=0 ca=all-terminal`.
- `rulings_open=0`.
- Lint: `blocking=0 resolution=0`. One advisory carried to the Gate:
  A4's Evidence field is 34 lines against a soft cap of 30 — it was
  already over budget before this pass, the load-bearing anchors are
  findable, and the match-path derivation was moved here with a
  pointer rather than truncated.

Verdict: RECONCILED.
