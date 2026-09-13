Model: claude-sonnet-5

# Pre-lock absorption audit — RDR 0021 (lint normalized-graph export), iteration 3

Scope per task: verify findings raised in the grounding, 3amigo (base +
iter-3), critique, repeatability rounds, and authors-round.md are reflected
in the CURRENT RDR text (read only via the `rdr` projector). This is a
re-run of the prior audit at `evidence/reconcile/absorption-audit.md`
(iteration 1), which found one residue (F4) and one soft/accepted residue
(`--flow` dead flag). Both were dispositioned at Stage 6 Reconcile iteration
2 (`evidence/reconcile/iter-2/reconcile.md`). This pass re-verifies those
dispositions directly against the live text and checks nothing has drifted
since, including through the intervening Stage 7.1 cluster-reconcile
(0021-0029) and Stage 4 re-verify of A4.

Current record state: `Status: Draft` (demoted from Final 2026-09-12 by the
0021-0029 cluster-reconcile gate; not yet re-locked). A
`## Refinement Context (cluster re-entry — delete on re-lock)` block is still
present — this is Stage 7.1 scaffolding, outside the five rounds this audit
is scoped to, and is reported separately below rather than folded into the
round ledgers.

---

## Round: grounding

`evidence/grounding/findings.md` is a confirmation sweep, not a findings
ledger — explicit: "Refuted / Not-Found: None. Every claim swept ...
confirmed against `main`." It ran an inverse-rule sweep (no sibling `graph`
verb, no `--emit`-like flag, no existing DOT writer, no existing
schema-versioned document — all confirmed genuinely new) plus a coverage
confirmation over A1-A5, C1, C4, C5, 0006:C19/C20, 0002's round-trip section,
JDR 0002 §D1, and the external state-machine-cat prior art.

**Classification**: N/A — no findings raised, nothing to absorb or leave as
residue.

---

## Round: 3amigo (base pass, iteration 1 + iteration 3 delta pass)

`evidence/3amigo/consolidation.md` (12 findings) + `evidence/3amigo/iter-3/consolidation.md` (2 more, delta-scoped after the base rewrite).

| # | id | Finding | Disposition | Verified in current RDR? |
| - | --- | --- | --- | --- |
| 1/4 | C2/A7 | `values` shape contradicted (object-keyed-by-tag vs array-of-strings) | fixed | ABSORBED. Current C2: "`values` an OBJECT keyed by tag name whose every value is that tag's sorted, deduplicated value ARRAY (the shape `reach.go::Node`'s `Values map[string][]string` projects without invention)". A7 status Verified. |
| 2 | C4 | `graph-export-too-large` depends on a `complete` bool `Reach` discards; no exported completeness surface | fixed (new A8, Pending) | ABSORBED as stated requirement. A8 remains Pending by design (pre-implementation verification task, not a text gap). |
| 3 | S6 | DOT oracle is exit-0 + id-set only; misses mis-escaped-but-well-formed labels | fixed | ABSORBED (confirmed again at iter-3, see below). |
| 5 | C5 | Forward-compat clause cites JDR 0002 §D1, read as unresolved | dismissed-with-cite | TERMINAL — citation confirmed correct (separate namespace from RDR 0002). |
| 6 | S5 | Byte-equality claim doesn't pin trailing-newline accounting | fixed | ABSORBED per consolidation.md basis (F1 pins text-mode stdout as document+"\n"; S5 now names the newline explicitly). |
| 7 | C1 | Unclear if "mirrors lint's arm set verbatim" covers help text | fixed | ABSORBED. Current C1: "mirroring is scoped to the ARM SET and its CODE SPELLINGS, not to lint's help/usage text". |
| 8 | S1 | No stated map-order provocation mechanism | fixed | ABSORBED per consolidation.md (S1 now names `GODEBUG=randmapiter=1`). |
| 9 | S8 | "No document on stdout on refusal" lived only in mini-check table | fixed | ABSORBED per consolidation.md. |
| 10 | §problem-statement | Motivates RDR partly on "external formal checker," out of §background's scope | dismissed-with-cite | TERMINAL — §background/§consequences frame formal models as a derivable downstream consumer; no contradiction. |
| 11 | MVV/S6 | Diagram-review outcome validated only by well-formedness + id-set, never legibility | charted-to-successor | TERMINAL, not residue — DOT styling is explicitly non-normative (Load-Bearing Decisions); reversing would undo a decided call. Recorded in `evidence/3amigo/charted.md`. Successor: a styling/layout RDR or kata. |
| 12 | §decision-rationale | CI-diffability claim has no real cross-commit/cross-build diff test | charted-to-successor | TERMINAL, not residue — a true test needs a two-build harness, outside C3's (model, build) narrowing. Successor: a CI-integration kata. |

Iteration 3 delta pass (2 more findings, both fixed):

| # | id | Finding | Verified in current RDR? |
| - | --- | --- | --- |
| 1 | S6 | Rewritten oracle said "decode back to source value" but named no decoder; contradicted stale mini-check rows | ABSORBED. Current S6-backing text (via C2/A6) inverts A6's recorded escaping order explicitly; mini-check rows swept to match per the consolidation's `§amendment-sweep` note. |
| 2 | C4 | Closing sentence called `checkNodeCeiling` "a distinct guard-product bound" — factually backwards | ABSORBED per iter-3 consolidation.md basis (same bound, same constant, two call sites; `ProductBound()` named as the bound NOT involved). |

**Residue**: none. Both charted items are terminal dismissals with a named successor, not unabsorbed findings.

---

## Round: critique (dual-model, iteration 1)

`evidence/critique/diff.md` + `resolve.md`. 20 rows across two independent
model passes, converged on 4 anchors (A8, C2, C4, C5). Outcome: "4 fixed (A8,
C2, C5, RT1) · 13 dismissed-with-cite · 1 charted-to-successor. 1 new Pending
assumption (A9)."

| Row | id | Finding | Disposition | Verified? |
| --- | --- | --- | --- | --- |
| A C-2/B C-1 | A8 | Pending assumption carries a normative MUST, ungated by Prerequisites | fixed | ABSORBED — resolve.md: "Added an A8 Prerequisites gate on Phase 2." A8 itself stays Pending by design. |
| A C-3 | C2 | `<opaque>` sentinel reaches `values[]` unnamed | fixed (new A9, Pending) | ABSORBED then SUPERSEDED — see A9 note below: C2 now names the sentinel and explicitly disclaims any meaning attached to it ("this record attaches no meaning to that spelling"), which closes the question A9 was tracking more completely than the original fix; A9 no longer exists as an id in the current record (8 assumptions, A1-A8). |
| A C-7/B C-5 | C5 | "One `jq -r` step" promised, no key named | fixed | ABSORBED. C5/C2 now name `data.dot`, unwrap `jq -r .data.dot`. |
| A C-10 | RT1 | "value identity on the exported projection" self-scoped/circular | fixed | ABSORBED per resolve.md (requantified over C2's field list). |
| A C-1/C-9 | C4, §metadata | `graph-export-too-large` vs `graph-product-too-large` code confusion | dismissed-with-cite | TERMINAL — traced to a stale doc comment in `taxonomy.go:150-152`, a source bug not an RDR defect; charted separately. |
| A C-4 | A6 | Re-raise of provisional spike spellings | dismissed-with-cite | TERMINAL — A6's Evidence already states spellings are provisional, C2 governs. |
| A C-5 | S6 | Escaping order exists only in a deleted `/tmp` spike file | dismissed-with-cite | TERMINAL — A6's Evidence records the order durably. |
| A C-6 | C2 | `[]`/`{}` vs `null` re-raise | dismissed-with-cite | TERMINAL — confirmed in current C2: "never `null`... ABSENT, never `null`". |
| A C-8 | A1 | TextLiner spike scale doubted | dismissed-with-cite | TERMINAL — refuted by A1 spike evidence (multi-KB payload tested). |
| B C-2 | C1 | `--flow` on `graph` resolves no ids — dead flag | dismissed-with-cite | TERMINAL (soft) — accepted by ruling, not silently dropped; current C1 states the behavior verbatim ("this build resolves no ids"); S8 covers the refusal arm. Re-confirmed not residue (see below). |
| B C-3 | C2 | Authored-TOML fidelity out of frame | dismissed-with-cite | TERMINAL — correctly routed to A3, separately handled by the authors-round refutation (absorbed, see below). |
| B C-4 | C4 | Oracle untested against future observer arm | dismissed-with-cite | TERMINAL — C4 text confirms mechanism-independence. |
| B C-6 | C2 | No machine-checkable schema artifact | dismissed-with-cite | TERMINAL — accepted risk, documented in §risks-and-mitigations. |
| B C-7 | §decision-rationale | Re-raise of settled ALT1/ALT2 choice | dismissed-with-cite | TERMINAL — settled in Decision Rationale/Alternatives Considered. |
| B C-8 | C4 | Refusal-over-partial-document re-raise | dismissed-with-cite | TERMINAL — confirmed in F1/C4. |
| B C-9 | §failure-modes | Killed-process case re-raise | dismissed-with-cite | TERMINAL — confirmed in F1. |
| B C-10 | §problem-statement | Abstraction-marker soundness re-raise | dismissed-with-cite | TERMINAL — confirmed; C2's soundness sentence present; §references names RDR 0022. |
| — | `taxonomy.go:150-152` | Stale doc comment (source bug) | charted-to-successor | TERMINAL — out of RDR scope, kata tracked separately. |

**Residue**: none as of this iteration. The `--flow` dead-flag item (B C-2)
was flagged "soft residue" in iteration 1 of this audit; Stage 6 Reconcile
iteration 2 re-examined it explicitly (item 5 in its open-set table) and
confirmed "ABSORBED — deferred by design": the behavior is stated verbatim in
C1, not silently dropped, and S8 covers the arm. Re-confirmed directly
against current C1 text this pass — holds.

---

## Round: repeatability (lite, run-1, profile: large)

`evidence/repeatability/diff.md` — independent reconstruction diffed against
the record. Verdict: "Healthy." 12 GUESS markers (8 discarded as taste, 1
reading error), 4 admitted findings (F1-F4), all genuine gaps.

| id | Finding | Disposition | Verified? |
| --- | --- | --- | --- |
| F1 | `--emit` validation has no fixed position in the arm order | fixed | ABSORBED. Current C1: "`--emit` validity is checked with the argument-shaped arms, AFTER the `--model`/`--flow` selection arms and BEFORE any file I/O." |
| F2 | C2 declares field spellings normative by citing a partial, self-disclaiming exhibit | fixed | ABSORBED. Current C2 states spellings "AS SPELLED HERE" and enumerates the full list inline, explicitly noting the Illustrative Code "disclaims literal assertion... binds nothing." |
| F3 | DOT terminal-node marking has no satisfaction predicate; `<opaque>` evaluation undefined | fixed, then superseded by a stronger fix | ABSORBED — and the question is now dissolved rather than merely answered: current C2 states "The document carries NO verdict or finding field ... and NO per-node terminal marking: which merged nodes satisfy a `terminal` predicate set is the dead-end quantifier RDR 0015 owns (JDR 0001 §JD-23), and this record declares no evaluator of its own." This is a stronger, later fix than the one the original iteration-1 audit verified (which had the document doing its own existential/opaque-admits evaluation) — the ownership was hoisted to RDR 0015 per a Stage 6 JDR-hoist decision, closing F3's predicate-ambiguity question by removing the document's claim to evaluate it at all. |
| F4 | The `model` member's wire value is unfixed (`Model.ID` vs `--model` path argument) | fixed | ABSORBED — CONFIRMED directly against current text. C2: "`model` (`table.Model.ID`, the AUTHORED `[model] id`, never the `--model <path>` argument or any path-derived string)". `D-identity` restates it identically per Stage 6 Reconcile iteration 2's delegated check. This closes the residue the prior iteration-1 audit flagged and left open. |

Additionally `evidence/repeatability/Charted.md` charts one out-of-scope
finding (`.rdr/resources.md` naming a nonexistent `CLAUDE.md`) — a seam/index
defect, correctly routed away from this RDR, not residue here.

**Residue**: none. F4 (the sole open item from the prior audit) is now
confirmed absorbed verbatim in the live C2 text.

---

## authors-round.md (Stage 4 Resolve) — context, not a scoped round, reported per instruction to flag refutations loudly

**Q3 / A3 falsification**: A3 ("the normalized model value carries everything
the document needs") was reported "FALSIFIED in part" — the model value does
not carry the reachability edges; `reach.go::Reach` returns only `[]Node`.

**CONFIRMED ABSORBED**, re-verified directly this pass (`0021:A3`): Status is
Verified, Evidence states explicitly "The edge relation is NOT carried by the
model value... recovered in the export path per A2... C2 predicates *carries*
of the document, which does carry edges." The assumption was narrowed rather
than the document contract weakened; no open loop remains.

**Q2 (`[]`/`null` commitment)** and **Q1 (contract-count Profile latch)** and
**Q3(b)/(c)** (A2 conditional on `subsumes`, edges-carrier as new code) were
all author rulings recorded in `rulings.md` and absorbed into the current
C2/C3/A2 text (confirmed by the repeated "never `null`" and "ABSENT, never
`null`" language and A2's now-Verified, source-grounded status).

---

## Context: Stage 7.1 cluster-reconcile (0021-0029) — outside the five scoped
## rounds, reported because it is the most recent disturbance to the text

`docs/rdr/cluster-reconcile/0021-0029/reconcile-report.md` demoted 0021 from
Final to Draft (STAGE-SCOPED re-entry at Stage 3/refine) over four defects.
All four are now absorbed in the live text, verified directly this pass:

- **Defect 1 (PW-1, C-1/C-2/C-3)** — tiered-surface obligation undischarged.
  ABSORBED: current C1 assigns `append-only` to the `--emit` format set with
  a named enumeration seam, and states `graph-export-too-large` inherits the
  CLIError `code` vocabulary's existing `seam: none (prose-only)` tier.
- **Defect 2 (C-14)** — C2's "closed 11-member list" cardinality claim
  contradicted 0029's append-only no-cardinality-assertion rule. ABSORBED:
  current C2's STABILITY paragraph states the field list is append-only and
  a consumer "MUST NOT assert on the field set's cardinality, a member's
  ordinal position, or a tail position"; "closed" survives only describing
  `DumpColumns` by name, which `0029:S9` exempts explicitly.
- **Defect 3 (PW-3, C-9)** — joint-check never saw 0029. ABSORBED: current
  Decision Rationale's Joint-check row states it was "re-run 2026-09-12 on
  the post-0029 peer set" and that "C4's tier obligation is discharged in C1
  and C2 above."
- **Re-verify A4 (C-11)** — envelope gained a field after A4 verified.
  ABSORBED: current A4 Evidence carries a dated re-verification paragraph
  ("Re-verified 2026-09-12 against the post-0029 envelope (C-11)") concluding
  the surface survives because A4's premise is `0005:C1`'s scope sentence,
  not envelope immutability.
- **Carried forward (C-8)** — S2 goldens must scope to the document, not the
  envelope. ABSORBED: current S2 states goldens "capture the DOCUMENT, not
  the envelope."
- **Standing joint decision (C-7)** — two-marker read order, hoisted to
  `cli/0029 §Normative Contracts` C4. Current C2 states the boundary and
  that "this `schema` field versions the DOCUMENT... `0029:C1`'s
  `schema_version` versions the ENVELOPE... and neither substitutes for the
  other" — 0021 cites the home rather than restating it, consistent with the
  hoist. Not yet independently verified whether 0029:C4 itself states the
  answer (0029 is Final per the git log's "finalize cli/0029" commit); out of
  this audit's scope (0029's text, not 0021's).

The `## Refinement Context (cluster re-entry — delete on re-lock)` block
itself is mechanical Stage 7.1 scaffolding still present in the record,
marked for deletion at re-lock — not a finding requiring further absorption,
since its four defects are all independently confirmed fixed in the
normative text above.

---

## SPIKE / assumption apparatus cross-check

Cross-checked every spike-shaped identifier named across the evidence tree
against the RDR's own S1-S10 test items and A1-A8:

- A2's differential test suite (36 real models + synthetic fixtures A-I) is
  named and cited in A2's own Evidence (Verified) — in the apparatus, not
  residue.
- A6's spike artifact is cited by A6's Evidence — present, not residue.
- A9 (the `<opaque>`-sentinel Pending assumption critique/repeatability
  raised) no longer exists as an id in the current record (`rdr inspect
  0021` lists only A1-A8). This is consistent with Stage 6 Reconcile
  iteration 2's DOWNGRADED disposition, which stopped C2 from depending on
  the refuted unforgeability limb rather than merely conditioning it — the
  current C2 text goes further, stating the document "attaches no meaning"
  to the spelling at all, which dissolves the need for a tracking assumption
  rather than leaving one open. Confirmed not residue: no normative clause
  in the current text depends on the sentinel's unforgeability or on any
  evaluation semantics over it (the terminal-evaluation question itself was
  separately hoisted to RDR 0015, per F3 above).

No other evidence-tree spike names a mechanism the RDR's own S1-S10/A1-A8
apparatus does not already cover.

---

## Summary table — all rounds

| Round | Findings raised | Absorbed | Terminal (dismissed/charted) | Residue |
| --- | --- | --- | --- | --- |
| grounding | 0 (confirmation only) | n/a | n/a | none |
| 3amigo (base + iter-3) | 14 | 12 | 2 (charted, named successors) | none |
| critique (dual-model) | 20 (~17 unique dispositions) | ~16 | ~13 dismissed-with-cite (1 soft, re-confirmed terminal) + 1 charted | none |
| repeatability | 4 admitted + 1 charted | 4 (F1-F4, F4 now confirmed) | 1 (charted, out of scope) | none |
| authors-round (context) | A3 falsified-in-part | confirmed absorbed (narrowed, Verified) | — | none |

**Overall verdict basis**: every finding across all five scoped rounds is
either absorbed into the current normative text (verified directly against
the live record this pass, not merely trusted from a round file's
self-report) or terminally dismissed/charted with a stated reason and, where
applicable, a named successor. The one residue the prior iteration of this
audit (iteration 1, pre-dating Stage 6 Reconcile) left open — F4, the
`model` wire-field identity — is now confirmed closed verbatim in C2. No
round finding refutes a claim the current RDR text relies on without that
refutation already being folded into a narrowing or hoisting fix (A3's
narrowing, A9's dissolution via RDR 0015 hoist). The only open material is
Stage 7.1 cluster-reconcile scaffolding (the re-entry block awaiting
mechanical deletion at re-lock) — outside this audit's five-round scope, and
itself already absorbed at the text level.
