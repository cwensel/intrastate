Model: claude-opus-5[1m]

# Stage 6 Reconcile — RDR 0021 (lint's normalized-graph export)

Date: 2026-09-12 · Iteration 3 · Verdict: **RECONCILED**

Iteration 2 closed RECONCILED. The record then went to Finalize, locked
(`01412d2`), and was demoted again by Stage 7.1 cluster reconcile `0021-0029`
(`663a6f5`) over untiered machine-readable surfaces. Two repair commits followed
— `e4ed3ab` (refine: tier the three surfaces, re-run the joint-check against
0029) and `c65b494` (resolve: re-verify A4, clear the qualifier) — plus an
earlier scope trim, `d37a921`, that removed A9/S10 with the clause that
depended on them.

This pass therefore reconciles a **delta**, not a fresh draft: what the cluster
re-entry disturbed, plus the two obligations iteration 2 booked out. A1–A8 are
Verified and untouched by the re-entry except A4, which the re-entry named.

## Stage 5 preflight

| Check | Answer | Note |
| --- | --- | --- |
| `--outcome lens` | `/rdr-reconcile` | `lens-large-row-complete`; Profile `large`, grounding→3amigo→critique row complete |
| `--outcome critique` | `none` | `critique-large-diffed`; `critique_models=differ` — no single-model fallback |
| `--outcome repeatability` | `none` | `repeatability-lite-complete`; `variant: lite (profile: large)` stamped on `run-1.md` |
| Determinacy add-on | `determinacy=fired`, line written | No `stopped:determinacy-trigger-unjudged`. The iteration-1/2 chain disagreement did not recur at this call. |

Status is plain `Draft` (`status_form=none`): the `re-verify A4 @refine`
qualifier was cleared at `c65b494`, so nothing is owed to this stage by
qualifier.

## The open set

Sources 2, 3 and 4 are **empty on disk**, and no Pre-Lock list accompanied the
invocation (`lens_findings_open=0`, `lens_stale=none`, `rulings_open=0`).

| Source | State |
| --- | --- |
| 1 — Pre-Lock needs-verification list | None pasted; no lens left an open finding |
| 2 — Pending/Unverified CAs | Empty: `ca_total=8`, `ca_verified=8`, `ca_pending=0`, `ca_unverified=0`, `ca_off_vocabulary_ids=[]` |
| 3 — spikes named but unrun | `spikes_unrun=[]`. Six spike files on disk (a1, a2, a5, a6, a8, a9); the absorption audit found no findings file naming a spike the RDR does not — see item 5 for `a9`'s status |
| 4 — exactness-word delta | Empty: `rdr lint 0021 | grep prose:exactness` returns nothing. The five findings the pre-refine tooling pass captured (canonical ×3, byte-identical, deterministic) are gone with the prose the refines rewrote |

So the open set is the **carried set**: two obligations from iteration 2, the
six cluster-reconcile obligations, and the one structural question the A9
removal raises.

## Dispositions

| # | Item | Source | Disposition | Evidence pointer / plan |
| --- | --- | --- | --- | --- |
| 1 | A9 — `<opaque>` unforgeable (a, REFUTED) / opaque-admits marking sound (b) | carried from iter-2 | **ACCEPTED — subject removed, not papered over** | `d37a921` removed A9 **together with the clause that relied on it**. Old C2 asserted per-node terminal marking, the opaque-admits evaluator, and the discriminator; current C2 declares "NO per-node terminal marking … this record declares no evaluator of its own", defers the quantifier to RDR 0015 (JDR 0001 §JD-23), and says of the sentinel "this record attaches no meaning to that spelling". A6 was narrowed in the same pass (terminal-satisfaction dropped from the derivability claim). No surviving clause depends on either limb. Verified by grep over the projected record: the only `<opaque>` mention is C2's pass-through sentence. |
| 2 | `0002:C11` — reserve `<opaque>` beside `<clear>` | carried from iter-2 | **MOOT — obligation discharged by removal** | Iteration 2 booked this because C2 conditioned its discriminator on the reservation. C2 no longer does, and no citation of `0002:C11` survives anywhere in the record. Nothing in 0021 now waits on 0002. (The `<opaque>` forgeability fact remains true on `main` and is preserved in the ledger; it is simply no longer this record's dependency.) |
| 3 | JDR 0001 §JD-23 — widen `Siblings:` to include 0021 | carried from iter-2 | **MOOT — obligation discharged by removal** | Same cause. C2 declares no export-side quantifier to hoist; Decision Rationale states it explicitly: RDR 0015 "owns terminal satisfaction over merged nodes (JDR 0001 §JD-23), on which C2 declares nothing — the export … marks no node terminal, and takes no side of the split quantifier." §JD-23's `Siblings:` line correctly still reads `0015, 0022`; no widening is owed. The record cites §JD-23 as the home, which is the sanctioned cite-don't-restate shape. |
| 4 | Defect 1 / PW-1, PW-2, C-1, C-2, C-3 — three surfaces ship untiered, no enumeration seam | cluster re-entry | **VERIFIED — discharged in C1 and C2** | C1: the `--emit` format set is declared `append-only` (`0029:C2`) with a MUST NOT on cardinality/position, and its seam named — "a new exported `internal/cli` accessor over the two members … builds with no new import and no cycle". `graph-export-too-large` mints no vocabulary: it is a member of the CLIError `code` set `0029:C4` already tiers, inheriting that row's `seam: none (prose-only)`. C2: a STABILITY paragraph tiers the field spellings `append-only`, names the seam as the Go struct's json tags ("the compiler is the seam"), and separates the two markers. |
| 5 | Defect 2 / C-14 — C2 asserts a cardinality `0029:C2` forbids | cluster re-entry | **VERIFIED — restatement removed** | "a closed 11-member list with `emit` appended last" is gone; C2 now "carries the dump's field vocabulary rather than restating a cardinality", and the STABILITY clause adds the MUST NOT on cardinality, ordinal position and tail position. The Determinacy field-names bullet was corrected to match. |
| 6 | Defect 3 / PW-3, C-9 — joint-check never saw 0029 | cluster re-entry | **VERIFIED — re-run and recorded** | `JC1` now reads "fired → 0029 (home: `cli/0029 §Normative Contracts` C4), re-run 2026-09-12 on the post-0029 peer set", covering 0012–0020, 0022–0024 **and** 0029 on both arms. `joint_checks=1`, `joint_check_home=homed`. See item 9 for the residual. |
| 7 | C-11 — A4 verified against the pre-0029 envelope | cluster re-entry (the named `re-verify` set) | **VERIFIED** | `c65b494` appended the re-verification to A4's Evidence: `0029:C1`'s non-`omitempty` `schema_version` reconciles itself with `0005:C1`'s omitempty budget in its own text, and A4's premise is `0005:C1`'s **scope sentence**, not envelope immutability — so a field added to the envelope leaves the conclusion standing. C-11's second ask is discharged with a negative: the record asserts no exact envelope key set. |
| 8 | C-8 — S2's goldens straddle the envelope boundary | cluster re-entry | **VERIFIED — goldens scoped to the document** | S2 now captures the DOCUMENT (`--as=text`, or `jq .data`), so `0029:C1`'s `schema_version` "never enters the pinned bytes"; expectations add "an envelope version bump does not" fail the pin. The Risks mitigation was updated to match. No golden exists on disk, so this was a wording fix as the cluster report predicted. |
| 9 | C-7 — two-marker read order (open joint decision, home `0029:C4`) | cluster re-entry | **ACCEPTED — cited, not answered; obligation stands at the home** | This record answers it nowhere and claims nothing about it, which is the correct posture for a non-owner: `JC1` names the question, states it is "hoisted to C4, already the home for the envelope-versus-document boundary", and commits to "cite the home once C4 states it". Not a blocker here. **But the home does not yet state it** — see Obligations below. |
| 10 | Round residue across all five pre-lock rounds | 1, 3 | **ABSORBED** | Delegated absorption audit (PASS, non-blocking): every finding across grounding, 3amigo, critique, repeatability and the authors' round is absorbed or terminally dismissed with a named successor; iteration 1's sole open residue (F4, the wire `model` identity) is confirmed closed verbatim in live C2 text (`table.Model.ID`, "never the `--model <path>` argument or any path-derived string"). Detail: `absorption-audit.md` beside this file. |
| 11 | Hollow-body / References completeness | prompt gate | **PASS** | No `_Draft placeholder._`, no seed-skeleton header. `## References` fully authored (peer cites, source paths, prior art, kata) with no bracketed placeholders. Investigation, Prerequisites and the four Phases are authored prose. The one bracketed hit outside the gate is `[Gate key: cross-cutting]`, a gate-section label. |

## The two hard rules

**Refutation.** No live refutation. The one historical refutation — A9 limb (a)
— was routed back as a BLOCKER at iteration 1, dispositioned at iteration 2, and
has since been resolved the strongest available way: the clause that depended on
it was withdrawn, so there is no assumption left to refute. This is the opposite
of papering over; the refuted fact survives verbatim in the escaped-defect
ledger (row 1), which is its durable home now that refine has collapsed the
change history. No spike or source search run for this pass refutes anything the
current draft relies on.

**MVV dependency.** Nothing deferred is MVV-critical. The MVV's five steps
assert fixture load; double-emit byte identity plus presence of every C2 field
including the `reach` block and its abstraction marker; DOT node/edge id-set
equality; `jq .data` value equality; and lint neutrality. All five rest on
A1/A2/A3/A5/A6/A8 — every one Verified. The removed A9 pinned marked-terminal
membership, which the MVV never asserted and which C2 no longer emits. Item 9 is
a consumer-side read-order question that pins nothing the MVV proves.

## Obligations leaving this stage

One, and it is **another record's write**. It does not block this lock.

| Owner | Obligation | Consequence while unlanded |
| --- | --- | --- |
| RDR 0029 (`C4`) | Answer C-7: which version marker a consumer reads first — the envelope's `schema_version` (`0029:C1`) or the document's `schema` (`0021:C2`) — and what it does when one is supported and the other is not. | `0021:JC1` promises to "cite the home once C4 states it", and C4 does not yet state it. C4 covers the adjacent point (the envelope version says nothing about `data`'s shape; a consumer pinning the document marker refuses `data` and keeps the envelope) but never sequences the two reads. 0029 is `Final` with `status_form=none`, and `rdr index --open-joint` reports `open: []`, so **no qualifier anywhere records the tolerance** the cluster report said would "apply at re-lock". 0021 is correct as written — it answers nothing and claims nothing — but a reader following the citation finds no answer. |

## Caveats

- **The open-joint tracker reads clean, and that is itself the finding.** The
  cluster report (`cluster-reconcile/0021-0029/reconcile-report.md`, `open.txt`)
  carries `C-7` as the one finding neither repaired nor closed by a demotion,
  and recorded the tolerance in both records' re-entry notes rather than a
  Status qualifier, because both Status lines held demotion qualifiers at the
  time. Both qualifiers have since been cleared by the repair passes, and the
  tolerance was not migrated to a qualifier. So `open_joint_decisions=[]` on
  both records reflects where the note was written, not whether the decision is
  answered. Flagged for Finalize rather than silently inherited.
- **`section:unknown-to-template` advisory (1226–1327).** The `## Refinement
  Context (cluster re-entry — delete on re-lock)` block is the only lint
  finding (`blocking=0 resolution=0 placeholder=0 advisory=1`). It is
  self-marked for deletion at re-lock — Stage 7's write, not a reconcile
  regression.
- **A9's spike file remains on disk** (`evidence/spikes/a9-opaque-sentinel.md`)
  for an assumption the record no longer carries. Kept deliberately: it is the
  evidence behind ledger row 1, and deleting it would orphan that row.

## Verdict

**RECONCILED** — every item terminal, no BLOCKER. Ready for Finalize.
