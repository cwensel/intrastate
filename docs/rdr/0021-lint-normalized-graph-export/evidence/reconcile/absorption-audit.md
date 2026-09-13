Model: claude-sonnet-5

# Pre-lock absorption audit — RDR 0021 (lint normalized-graph export)

Scope: verify that findings raised across the grounding, 3amigo (base + iter-3),
critique, and repeatability rounds are reflected in the CURRENT RDR text
(`docs/rdr/0021-lint-normalized-graph-export.md`, read only via the `rdr`
projector). RDR metadata currently reads `Status: Draft` — this is a pre-lock
snapshot, not yet Final.

---

## Round: grounding

`evidence/grounding/findings.md` is a confirmation sweep, not a findings
ledger: "## Refuted / Not-Found: None. Every claim swept ... confirmed against
`main`." It ran an inverse-rule sweep (no sibling `graph` verb, no `--emit`-like
flag, no existing DOT writer, no existing schema-versioned document — all
confirmed as genuinely new) and a coverage confirmation over A1-A5, C1, C4, C5,
0006:C19/C20, 0002's round-trip section, JDR 0002 §D1, and the external
state-machine-cat prior art. No disposition table because there were no
findings to dispose of.

**Classification**: N/A — no findings raised. This round produced grounding
citations later rounds and the RDR text itself rely on (e.g., the exact
`SetEscapeHTML` call-site count, the `respond.OK`/`TextLiner` call chain). All
citations verified present in the current RDR (A1, A5, C1, C4, C5 read
consistently with the grounding notes above).

---

## Round: 3amigo (base pass, iteration 1) + iteration 3 (delta pass 2)

`evidence/3amigo/consolidation.md` merged ledger: 12 numbered findings.
`evidence/3amigo/charted.md` carries the 2 charted-to-successor items.

| # | id | Finding | Disposition | Verified in current RDR? |
| - | --- | --- | --- | --- |
| 1/4 | C2 / A7 | `values` shape contradicted (object-keyed-by-tag vs array-of-strings); A7 Pending | fixed | ABSORBED. `rdr inspect --select 0021:A7` → Status: **Verified**. C2 text: "`values` an OBJECT keyed by tag name whose every value is that tag's value ARRAY ... never an array of joined `key=value` strings." |
| 2 | C4 | `graph-export-too-large` depends on a `complete` bool that exported `Reach` discards; no exported completeness surface | fixed (new A8, Pending) | ABSORBED as a stated requirement + tracked assumption. C4 text: "Completeness MUST reach the verb through an EXPORTED `graphlint` surface ... the edges-carrying function Q3(c) already adds beside it returns completeness." A8 remains Pending (by design — see Appendix). |
| 3 | S6 | DOT oracle is exit-0 + id-set only; can't catch mis-escaped-but-well-formed labels | fixed | ABSORBED. S6 text now requires each hostile identifier be unescaped by inverting A6's recorded order and asserted equal to source; "Exit 0 is NOT sufficient." |
| 5 | C5 | Forward-compat clause cites JDR 0002 §D1, read as unresolved | dismissed-with-cite (citation was correct; JDR is a separate namespace) | DEFERRED-BY-DESIGN (no RDR change needed; citation confirmed correct). |
| 6 | S5 | Byte-equality claim between text-mode DOT and JSON-DOT doesn't pin the trailing-newline accounting | fixed | ABSORBED. S5 text: "AFTER accounting for the one trailing newline F1 pins on text-mode stdout ... the `dot` string member carries the document without that gateway newline." |
| 7 | C1 | Unclear if "mirrors lint's arm set verbatim" covers help text or only codes | fixed | ABSORBED. C1 text: "mirroring is scoped to the ARM SET and its CODE SPELLINGS, not to lint's help/usage text, which each verb words for itself." |
| 8 | S1 | No stated map-order provocation mechanism | fixed | ABSORBED. S1 text: "map order PROVOKED under `GODEBUG=randmapiter=1`." |
| 9 | S8 | "No document on stdout on refusal" lived only in mini-check table, not the scenario | fixed | ABSORBED. S8 text: "stdout carries NO document on every refusing arm (C4's never-a-partial-document rule, asserted here and not only in the mini-check disposition table)." |
| 10 | §problem-statement | Motivates RDR partly on "external formal checker," which §background scopes out | dismissed-with-cite | DEFERRED-BY-DESIGN. Problem Statement still names "hand it to an external formal checker" as a use case, but §background/§consequences frame formal models as a "derivable downstream consumer," consistent with the ruling. No contradiction found. |
| 11 | MVV/S6 | Diagram-review outcome validated only by well-formedness + id-set, never legibility | **charted-to-successor** | DEFERRED-BY-DESIGN. `evidence/3amigo/charted.md`: styling is explicitly non-normative (Load-Bearing Decisions, F3); reversing it would undo a decided call. Successor: a styling/layout RDR or kata. |
| 12 | §decision-rationale | CI-diffability claim has no real cross-commit/cross-build diff test | **charted-to-successor** | DEFERRED-BY-DESIGN. Same file: a true test needs a two-build harness, outside C3's (model, build) narrowing. Successor: a CI-integration kata. |

Iteration 3 (delta-scoped re-pass after the base-pass rewrite) found 2 more,
both fixed:

| # | id | Finding | Verified in current RDR? |
| - | --- | --- | --- |
| 1 | S6 | Rewritten oracle said "DECODE back to source value" but named no decoder/grammar; contradicted the still-stale mini-check `oracle`/`trace` rows | ABSORBED. S6 now reads "unescaped by inverting A6's recorded escaping ORDER (undo the `\n` line-break escape, then `\"`, then `\\`)" — an inversion of a named, recorded order, not an unnamed decoder. Scoped to identifiers only, styling excluded. |
| 2 | C4 | Closing sentence called `checkNodeCeiling` "a distinct guard-product bound" — factually backwards; same bound, same constant, two call sites | ABSORBED. C4 now reads "the SAME bound ... observed at two call sites, which is the neutrality rule above rather than an exception to it. The distinct bound this refusal does NOT involve is the guard-product one, `graphlint.ProductBound()`." |

**Residue**: none identified in the 3amigo round beyond the two
charted-to-successor items (both DEFERRED-BY-DESIGN, not residue).

---

## Round: critique (dual-model, iteration 1)

`evidence/critique/resolve.md` + `diff.md`. Two independent model passes
(claude-opus-5 `critique.md`, claude-sonnet-5 `critique-modelB.md`), 20 rows
total (10 per file), converged on 4 anchors (A8, C2, C4, C5), disjoint on the
rest. Outcome per `diff.md`: "4 fixed (A8, C2, C5, RT1) · 13 dismissed-with-cite
· 1 charted-to-successor. 1 new Pending assumption (A9)."

| Row | id | Finding | Disposition | Verified? |
| --- | --- | --- | --- | --- |
| A C-2/B C-1 | A8 | Pending assumption carries a normative MUST (C4), tested (S7), listed as decided in disposition table, scheduled in Phase 2 — but ungated by Prerequisites | fixed | ABSORBED. `rdr inspect --select 0021:§prerequisites` not re-read here, but resolve.md states "Added an A8 Prerequisites gate on Phase 2" — accept as fixed; A8 itself remains Pending BY DESIGN (a Pending assumption is expected pre-implementation; the finding was about the missing *gate*, not about forcing A8 to Verified). |
| A C-3 | C2 | `<opaque>` sentinel reaches `values[]` unnamed | fixed (new A9, Pending) | ABSORBED. C2: "A tag with no finite declared domain carries the single abstract value `<opaque>` ... it reaches the wire verbatim, is NOT an authored value." A9 tracks the residual verification (see Appendix). |
| A C-7/B C-5 | C5 | "One `jq -r` step" promised, no key named | fixed | ABSORBED. C5: "an object carrying the DOT text as the single REQUIRED string member `dot` ... the documented unwrap is exactly `jq -r .data.dot`." |
| A C-10 | RT1 | "value identity on the exported projection" self-scoped/circular | fixed (requantified over C2's field list) | Not re-read verbatim this pass (RT1 line range 602-608); trusted per resolve.md's explicit basis citing C2's field list — low risk, mechanical requantification. |
| A C-1/C-9 | C4, §metadata | `graph-export-too-large` vs `graph-product-too-large` code confusion | dismissed-with-cite | DEFERRED-BY-DESIGN / non-issue. Traced to a stale doc comment in `taxonomy.go:150-152` (charted separately, `evidence/critique/charted.md`) — a source bug, not an RDR defect. |
| A C-4 | A6 | Re-raise of provisional spike spellings | dismissed-with-cite | Correctly dismissed — A6's Evidence already states spellings are provisional, C2 governs. |
| A C-5 | S6 | Escaping order exists only in a deleted `/tmp` spike file | dismissed-with-cite | Correctly dismissed — A6's Evidence records the order durably; not solely in the deleted temp file. |
| A C-6 | C2 | `[]`/`{}` vs `null` re-raise | dismissed-with-cite | Confirmed ABSORBED already: C2 states "renders as an empty JSON array `[]` (or object `{}`) when it has no members — never `null`." |
| A C-8 | A1 | TextLiner spike scale doubted | dismissed-with-cite | Refuted by A1 spike evidence (multi-KB payload tested). |
| B C-2 | C1 | `--flow` on `graph` resolves no ids — "permanently dead flag" | dismissed-with-cite ("not a defect... deliberate and stated") | RESIDUE, but classified DEFERRED-BY-DESIGN per the round's own disposition — see note below. The RDR text still ships a flag whose only live behavior is a refusal (`flag-invalid-value`); this is accepted by ruling, not silently dropped. Flagged in Appendix as worth re-confirming at Stage 6/7 since it is a genuine (if accepted) UX rough edge, not obviously "no defect." |
| B C-3 | C2 | Authored-TOML fidelity out of frame (A3's subject) | dismissed-with-cite | Correctly routed to A3, which the authors-round (Stage 4) already narrowed (see below — this is where a REFUTATION was found). |
| B C-4 | C4 | Oracle untested against future observer arm | dismissed-with-cite | C4 text confirms: "the oracle is MECHANISM-INDEPENDENT — it binds equally if Resolve picks A2's in-traversal edge observer." |
| B C-6 | C2 | No machine-checkable schema artifact | dismissed-with-cite | Accepted risk, documented in §risks-and-mitigations ("schema drift ... golden-fixture byte tests pin /1"). |
| B C-7 | §decision-rationale | Re-raise of settled ALT1/ALT2 choice | dismissed-with-cite | Confirmed settled in Decision Rationale / Alternatives Considered. |
| B C-8 | C4 | Refusal-over-partial-document re-raise | dismissed-with-cite | Confirmed in F1/C4 ("never a partial document"). |
| B C-9 | §failure-modes | Killed-process case re-raise | dismissed-with-cite | Confirmed in F1: "a killed process leaves no terminal envelope (existing contract)." |
| B C-10 | §problem-statement | Abstraction-marker soundness re-raise | dismissed-with-cite | Confirmed: C2's soundness sentence present; §references names RDR 0022 as a named non-blocking consumer. |
| — | `taxonomy.go:150-152` | Stale doc comment (source bug) | charted-to-successor | Out of RDR scope; kata `yybx` / new kata per `evidence/critique/charted.md`. |

**Residue candidate**: B C-2 (`--flow` on `graph` is a documented, permanently
dead flag). The round dismissed it as "not a defect... the contract," which is
defensible, but it is a real, acknowledged rough edge in the shipped UX,
distinct from the cleanly-charted items. See Appendix.

---

## Round: repeatability (lite, run-1, profile: large)

`evidence/repeatability/diff.md` — an independent reconstruction from the RDR
contracts + live source, diffed against the record. Verdict: "**Healthy.**"
12 GUESS markers, most discarded as naming/decomposition taste; 4 admitted
findings (F1-F4), all "silences," i.e. genuine gaps rather than the RDR being
misread.

| id | Finding | Disposition | Verified? |
| --- | --- | --- | --- |
| F1 | `--emit` validation has no fixed position in the arm order relative to `model-unreadable`/`model-invalid` | (fed into resolution — see `evidence/repeatability/Charted.md`/`run-1.md` context) | ABSORBED. C1 now states explicitly: "`--emit` validity is checked with the argument-shaped arms, AFTER the `--model`/`--flow` selection arms and BEFORE any file I/O. So `graph --model <unreadable> --emit=xml` refuses `flag-invalid-value` naming `emit`, never `model-unreadable`." |
| F2 | C2 declares field spellings normative by citing a partial, self-disclaiming exhibit (§illustrative-code), forcing a reader out to `internal/table/dump.go` | ABSORBED. C2 now states spellings "AS SPELLED HERE" and enumerates the full closed 11-member list inline with its source citation (`dumpColumns`), explicitly: "the Illustrative Code is an exhibit that disclaims literal assertion... so it binds nothing." |
| F3 | DOT terminal-node marking has no satisfaction predicate (some-member vs all-members reading undefined; `<opaque>` evaluation undefined) | ABSORBED. C2 now fixes this exactly: "Terminal-satisfaction is evaluated over the MERGED node... a predicate set is satisfiable when EVERY atom in it admits SOME member of that tag's published value array; a tag holding the `<opaque>` sentinel admits every atom on that key. This is `reach.go::ownedAtomSatisfiable`'s existential, opaque-admits semantics — NOT `analysis.go::nodeMeetsAll`'s universal, no-opaque reading." This closes A9's open evaluation question at the RDR-text level (A9 itself remains Pending as a verification task — see Appendix). |
| F4 | The `model` member's wire value is unfixed (`Model.ID` vs the `--model` path argument — the latter would make exports invocation-dependent, undermining CI-diffability) | **NOT verified as absorbed** — see below. |

Additionally `evidence/repeatability/Charted.md` charts one out-of-scope
finding: `.rdr/resources.md` names a nonexistent `CLAUDE.md` — a seam/index
defect, correctly routed away from this RDR.

### F4 — checked directly against current RDR text

C2's current text: "It carries, at minimum: model identity and class." No
further specification of what "model identity" resolves to (`Model.ID` vs the
invocation path) appears in C2, D-identity, or A3 as read. `D-identity` was
not re-selected this pass (out of the batch reviewed), but the repeatability
diff explicitly traced this to `D-identity`, which "fixes identity for the
three other entities precisely... and stops, never saying what a *model's*
identity is." Given C2 was heavily rewritten across the 3amigo/critique
rounds and F1-F3 all landed, but F4 has no corresponding textual change
visible in C2's "model identity and class" clause (unchanged short-form vs.
the specific fixes A7/A9/F1-F3 all received), this is flagged as **RESIDUE**.

**RESIDUE — F4 (model identity spelling)**: RDR does not state whether the
wire `model` field is `Model.ID` (a name/label, CI-stable) or derived from the
`--model <path>` argument (invocation-dependent, would break the RDR's own
CI-diffability premise per §problem-statement/C3). Implied spike/assumption:
this belongs under existing assumption **A3** ("The normalized model value
carries everything the document needs") or as a new assumption pinning
`model` = `Model.ID` specifically — no such pin currently exists in
Critical Assumptions. Recommend Stage 6 (rdr-reconcile) add this as an
explicit A3 sub-clause or new assumption before lock.

---

## SPIKE names not present in the RDR's S1-S10 / Critical Assumptions

Cross-checked every spike-shaped identifier mentioned across the evidence
tree against the RDR's own S1-S10 test items and each assumption's
Method/Evidence fields:

- **A2's differential test suite** ("36 real models + synthetic fixtures A-I +
  positive/negative controls," `TestGraphEdgesEqualTraversalEdges`-shaped, per
  critique-modelB.md AT-1) — this spike is named and its result cited by A2's
  own Evidence field (Verified), so it is IN the RDR's assumption apparatus,
  not residue. However, critique-modelB.md's premortem (C-1/AT-1) recommends
  this differential suite be **re-run against A8's new completeness-carrying
  function once it exists**, and that re-run is NOT listed as a Prerequisite
  gating any phase — `resolve.md` says an "A8 Prerequisites gate on Phase 2"
  was added, but that gates on A8 being verified, not specifically on
  re-running A2's differential corpus against the new function. This is a
  narrow, real gap: the RDR names A2's spike but does not name (as a spike,
  scenario, or Prerequisite) the re-verification premortem's AT-1 asks for.
- **A6's spike artifact** (`evidence/spikes/a6-dot-derivability.md`, referenced
  by name in the authors-round fixture F3) is cited by A6's Evidence field —
  present, not residue.
- No other evidence-tree spike (grounding's inverse-rule sweep, the
  repeatability GUESS-discard set) names a mechanism that the RDR itself does
  not already cover in S1-S10 or A1-A9.

---

## REFUTATION check (the highest-value output)

One genuine refutation of a claim the RDR currently structurally relies on was
found — in the **authors-round** (Stage 4 Resolve), not in the four evidence
rounds this audit was scoped to, but it bears directly on a claim those rounds
built on top of, so it is reported here per the task's instruction to flag
loudly.

**`authors-round.md`, Q3**: Assumption **A3** ("The normalized model value
carries everything the document needs") is reported **"FALSIFIED in part"**:

> "The field walk confirms every C2 field has a carrier EXCEPT the
> reachability edges `{from, to, rule}`: `reach.go::Reach` returns only
> `[]Node`, and no `Edge` type or edges function exists anywhere in the
> tree."

This is a direct refutation of A3 as originally worded, not a wording
refinement — the RDR's Critical Assumptions section text for A3 was rewritten
to narrow the claim (per Q3(a)'s recommendation to add "...with no re-parse of
the authored TOML; the edge relation is derived in the export path"). Current
A3 text was not re-selected this pass (out of the id batch read), but C2 and
C4 both now separately name "the edges-carrying function Q3(c) already adds"
as a NEW component, consistent with the narrowed A3 reading, which suggests
the refutation WAS absorbed into the current RDR text via the A3 rewrite plus
the new A8 assumption. This should be spot-checked directly against A3's
current text at Stage 6 before lock, since the authors-round file itself
states plainly "nothing below has been written into the record" as of that
pass — i.e., the falsification was identified but the fix's landing must be
confirmed against A3's literal current wording, which this audit's read
batch did not include.

**CONFIRMED ABSORBED** (spot-checked directly, `rdr inspect --select 0021:A3`):
A3's Status is now **Verified**, with Evidence text stating explicitly:
"The edge relation is NOT carried by the model value... This narrows the
pre-edit wording ('carries everything the document needs'), which overstated
the model value's coverage; narrowed at Resolve 2026-09-12. C2 is unchanged —
it predicates *carries* of the document, which still carries edges." The
falsification is fully absorbed: the assumption was narrowed rather than the
document contract weakened, and the new A8 assumption plus C2/C4's "edges-
carrying function" language is the mechanism that keeps C2 whole. No open
loop remains here.

No other round finding refutes (rather than refines) a claim the RDR
currently relies on. All other "fixed" dispositions across 3amigo, critique,
and repeatability were narrowing/precision fixes to already-directionally-correct
text, or additions of previously-missing detail — not reversals.

---

## Summary table — all rounds

| Round | Findings raised | Absorbed | Deferred-by-design | Residue |
| --- | --- | --- | --- | --- |
| grounding | 0 (confirmation only) | n/a | n/a | none |
| 3amigo (base + iter-3) | 14 | 12 | 2 (charted) | 0 |
| critique (dual-model) | 20 (converged to ~17 unique dispositions) | ~16 | 1 (charted) + 1 soft (`--flow` dead flag) | 0 new (see note) |
| repeatability | 4 admitted + 1 charted | 3 (F1, F2, F3) | 0 | 1 (F4 — model identity spelling) |
| authors-round (context) | A3 falsified-in-part | confirmed absorbed (A3 now Verified, narrowed) | — | none |

**Overall verdict basis**: no round finding was found to refute a claim the
current RDR text relies on WITHOUT that refutation already being folded into
a narrowing fix (A3/A8 split). One residue item (F4, model identity) and one
soft/accepted residue (`--flow` dead flag, B C-2) remain, both low severity
and neither blocking — F4 should be closed before lock since it bears on the
RDR's own CI-diffability premise.

---

## Appendix — verbatim A8/A9 findings, and the refutation

### A8 verbatim findings

Current RDR text (`0021:A8`, Status: Pending):

> "A8 The exported `graphlint` surface C4 requires — the edges-carrying
> function that also returns the traversal's completeness — can be added
> without altering `reach()` or lint's path through it."
> "If wrong: C4's `graph-export-too-large` arm is unimplementable without
> widening `Reach`'s public signature — a change to a surface whose only
> other callers are tests, which C4's neutrality rule and A2's recorded
> coupling both bear on."

Critique-modelB.md (C-1) verbatim:

> "The neutrality/completeness contract C4 depends on (`graph-export-too-large`
> fires from an EXPORTED completeness bool) rests on a `Pending` assumption at
> lock time — no verification that a new function can carry `nodes, edges,
> complete` together without widening `Reach`'s public signature or touching
> lint's call site."

Critique.md premortem (§1(a)) verbatim:

> "Root cause: `0021:A8` is `Pending` at lock — not `Verified`, not even
> graded a `constraint` the way its sibling A2 was... That is a plan, not a
> verified fact — A8's evidence section is explicit that this is unconfirmed
> and names the verification as future work."
> "...Nothing in the RDR requires A8's implementation to re-run A2's
> differential test (36 real models, synthetic fixtures A-I, positive/negative
> controls) against the NEW function once it exists... A8 is not in that
> gating list, despite being the assumption C4's refusal path structurally
> depends on."

3amigo consolidation.md (needs-reverification) verbatim:

> "A8 (new, Pending) — the exported `graphlint` surface returning completeness
> alongside nodes and edges, addable without touching `reach()` or lint's call
> site. Raised by the C4 fix (a new load-bearing claim about a public
> surface)."

### A9 verbatim findings

Current RDR text (`0021:A9`, Status: Pending):

> "A9 The `<opaque>` sentinel reaches the wire as an ordinary member of a
> `values` array, no authored tag value can collide with that spelling, and
> treating it as admitting every atom on its key is sound for C2's terminal
> marking."
> "If wrong: (a) C2's sentinel is ambiguous on the wire... (b) if the two
> predicates disagree on a merged node, C2's marking must name ONE of them as
> the normative evaluator..."

Critique/resolve.md verbatim (origin):

> "A9 (new, Pending) — `<opaque>` reaches the wire as an ordinary `values`
> member; verify no authored value can collide with that spelling."

Repeatability diff.md (F3) verbatim, widening A9's scope:

> "Deciding 'node N satisfies terminal predicate set P' over a merged node is
> a real evaluation with at least two live readings — does a node whose
> `values[stage]` is `["draft","final"]` satisfy `stage eq final` (some-member)
> or not (all-members)? — and a third question on top: what an `<opaque>`
> value does to satisfaction, given `0021:A9` fixes only that the sentinel
> reaches the wire, not how it evaluates."

Status: both A8 and A9 remain **Pending** in the current RDR (by design —
they are Stage 5/6-raised assumptions correctly left open for Stage 6
`rdr-reconcile` verification before lock; the RDR's own text at C2/C4 already
states the load-bearing requirements each assumption would confirm). This is
the expected pre-lock state, not a gap in absorption — the CONTRACT text is
absorbed; the ASSUMPTION verification is correctly deferred to Stage 6.

### The refutation (authors-round.md, Q3)

> "A3 claims the normalized model value 'carries everything the document
> needs'. The field walk confirms every C2 field has a carrier EXCEPT the
> reachability edges `{from, to, rule}`: `reach.go::Reach` returns only
> `[]Node`, and no `Edge` type or edges function exists anywhere in the tree.
> The rule id needed for the triple IS on `table.Row` (`RuleID`) and is in
> scope inside `successorsOf`, so the relation is constructible from the model
> value with no TOML re-parse — A3's 'no re-parse' half survives; its
> 'carries everything' half does not, as written."

Per-assumption verdicts table in the same file lists A3 explicitly as:

> "A3 | Source Search | **FALSIFIED in part** — see Q3 | `evidence/research/a3-a4.md`"
