model: claude-opus-5

# Tooling Pass — RDR 0021 (lint's normalized-graph export), iter-4

Run: Stage 7 mechanical pre-sweep, immediately before the Finalization Gate's
written responses. Source: `rdr lint --locking 0021` at `./lint.txt` (baseline)
and `./lint-after.txt` (post-edit receipt), plus `inspect --json --filter
assumptions,edges,metadata`. Findings only; the two RDR edits this pass made
(the 0012 joint-check citation, the re-entry note deletion) were made by the
Gate, not by this sweep.

Re-entry context: this is the re-lock pass after a Stage 8 (Implement)
route-back — `TARGET RE-ENTRY STAGE: 3 (STAGE-SCOPED)`, re-verify A3, over
REQ-20's unsatisfiable `domain` presence rule. Refine (`5154721`), resolve
(`5fe32cf`) and reconcile (`ce2065a`) all ran after the note was written; the
status qualifier is already cleared (`status_form=none`,
`status_reentry=false`, `lens_stale=none`).

## CHECK 1 — Template section coverage

- Baseline carried one `parse:section:unknown-to-template` (1283-1341): the
  `## Refinement Context (re-entry — delete on re-lock)` block. Stage 7's rule
  is that it MUST be gone by re-lock with its defect folded into live text.
  Verified closed in live text (table below), so the note is discharged and
  deleted in this pass; `lint-after.txt` no longer reports it. Not a spine hole.
- No `template:missing-section`. No `placeholder:survived`, no `scaffold:row`,
  no `contract:template-example` — the four gate responses this lock moves to
  `gate.md` are authored, not bracketed.
- No surviving `this is a seed skeleton` header. `## References` is authored
  (peer clauses, source reviewed, prior art with a search record, related kata).

Re-entry note discharge, item by item:

| Note item | Closed in live text | Where |
| --- | --- | --- |
| The defect — REQ-20 binds `domain` presence to `AssignmentCount` finiteness, unsatisfiable because `bool`/bounded `int`/`set` report finite while populating no `decl.Domain` | yes | C2 now reads "`domain` carries the declaration's AUTHORED members and is present exactly when there are any — an `enum` with a non-empty `domain` — and ABSENT, its key omitted, otherwise", with the explicit carve-out "Presence is NOT `guard.AssignmentCount` finiteness: that predicate reports finite for `bool`, for a bounded `int`, and for `set` — three kinds that author no `decl.Domain`" |
| Resolution direction — narrow to the authored-members arm; `int` omits `domain` too; mint no derived value vocabulary | yes | C2 states all three: the enum-with-members arm, `int` omitted, and "No derived value vocabulary is minted here: neither `bool`'s two literals, nor a `set`'s `Elements` or its powerset, nor an `int`'s `{min..max}` bound" — with the rejected compact-bound form recorded as additive-later |
| Re-verify A3 against the shipped exporter | yes | A3 carries "Re-verified 2026-09-13 at reconcile against the narrowed C2", citing the exporter's `finite && len(decl.Domain) > 0` gate at `graph_document.go:185` and `domainSize`'s per-kind arms; `evidence/reconcile/iter-4/a3-domain-presence.md` |
| Scope fence — C1/C3/C4/C5 and the edge work not reopened | yes | those four clauses are byte-unchanged across the refine/resolve/reconcile commits; only C2's `domain` clause and A3's Evidence moved |
| Residual: doc comment `graph_document.go:51-54` still asserts the refuted finite-only claim | implementation-side, correctly scoped out | verified STILL OPEN on the unmerged branch (`worktree-rdr-0021`): the comment reads "the finite arm always carries at least one member". A3's Evidence names it as "the implementation's to fix, not this record's" — the record's rule is now correct and the code already implements it; only the comment documenting the gate is stale |
| Residual: `TestReq20And28` asserts only the enum-with-members and `scalar` arms | implementation-side, correctly scoped out | same disposition; A3 names the `bool`/bounded-`int`/member-less-`enum` coverage gap as "the coverage S2 requires on both sides of the rule", owed by the implementation pass |

Both residuals are code/test defects on an unmerged branch, not record
defects: the RDR's normative text is what this gate locks, and it now states
the behavior the shipped exporter already has. Neither is a spine hole and
neither blocks Final.

## CHECK 2 — Method label vocabulary

PASS. Eight rows, every `method.off_vocabulary[]` empty
(`ca_off_vocabulary=0`). Members: `Spike` (A1, A2, A5, A6), `Source Search`
(A3, A7, A8), `Peer RDR` (A4) — all sanctioned. No row lacks a Method field.
No record was added or relabelled by the refine/resolve/reconcile passes; A3
kept `Source Search` while its Evidence grew.

## CHECK 3 — Source Search self-reference

PASS. The `Source Search` rows are A3, A7, A8; every Evidence anchor resolves
into the product tree (`graphlint/reach.go::Reach`, `::Node`, `::reach`,
`guard/declaration.go::AssignmentCount`, `guard/product.go::Groups`,
`graphlint/analysis.go::newAnalysis`, `graphlint/engine.go::Run`,
`graphlint/taxonomy.go::CodeProductTooLarge`). None resolves to this record or
its artifact dir. A3's re-verification cites
`evidence/reconcile/iter-4/a3-domain-presence.md` alongside a source anchor —
an evidence-dir path beside a resolving symbol, which is where a reconcile
artifact belongs, not self-reference. No regression from the rounds.

## CHECK 4 — Docs Only on load-bearing claims

PASS, vacuously. No assumption carries `method.members == ["Docs Only"]`.

## CHECK 5 — Symbol resolution of Source Search / Spike anchors

PASS. 31 `source-anchor` edges, every one `resolved: true` — none false, none
ABSENT, so `--repo` was supplied (`RDR_SOURCE_REPO` is bound) and the lookups
genuinely ran. Read from `--filter edges`, which is the authority here: the
`assumptions` filter's `evidence.anchors[]` objects carry only `to` and omit
the `resolved` key entirely, so that projection must not be read as ABSENT. No blocking
`edge:unresolved`. The three non-`true` edges in the projection are two
`issue/` references (`issue/4hps`, `issue/yybx`) and the
`{ARTIFACT_DIR}/gate.md` artifact pointer — none is a source anchor, and the
gate.md pointer resolves once this lock writes the file. A5's
`clierr.go:174::WriteJSONLine` carries a symbol that resolves; the stale line
component is a documented NON-finding.

## CHECK 6 — Status consistency

PASS. Status `Draft`, no qualifier (the Stage 8 route-back qualifier was
cleared at resolve, `5fe32cf`). All eight assumptions `Verified`
(`ca=all-terminal`), so no Pending/Unverified property is relied on as settled
fact. The gate responses move to `gate.md` in this same lock, so there is no
second copy to disagree with the record.

## CHECK 9 — Evidence-field budget

ADVISORY, reported not blocking. One `evidence:over-budget`: A3's Evidence
field is 43 lines against a soft cap of 30 (177-219). Profile is `large`, not
`foundational`, so `lint --locking` does not mark it blocking. The Gate's
item-2 response answers the author's question: the load-bearing anchors stay
findable (four `path::Symbol` anchors — `reach.go::Reach`,
`declaration.go::AssignmentCount`, `product.go::Groups`, `reach.go::heldValues`
— each resolving `true` under `--repo`), and the balance is the
per-kind `domainSize` reasoning that IS the re-verification this route-back
demanded — the content the next grounding sweep reads. Not truncated; not
relocated.

## CHECK 10 — Linking

PASS. C1-C5 labelled (no `label:contracts`); A4's `Method: Peer RDR` cites an
element, not a bare record (no `peer-evidence:no-element`); no
`edge:unresolved`, no `edge:unresolved-terminal`.

## Joint-decision fence — one blocker, fired and cleared in this pass

The fence first resolved `stopped:overlap-uncited` (rule `fence-uncited`):
`index --literal-intersect --record 0021` reported `0012 0021 UNCITED 5
shared: bool enum int scalar set`. This is a NEW pair, created by this
re-entry's own REQ-20 narrowing, which introduced the five kind tokens into
C2's presence rule; the 0017 pair iter-3 fired on now reads `cited`. Fired the
check rather than syncing copies, grounded per §ground-before-ask
(`--outcome ground` → `code`, settled in source and records):

- `0003:C2` (Status **Implemented** — terminal) OWNS the vocabulary outright:
  "The value kinds are exactly five, spelled with these tokens wherever a kind
  is named". `0003:C6` fixes which kind admits a finite domain; `0003:C11`
  what lint may claim over them.
- `0012:C1` inherits it by explicit citation — "the five-kind vocabulary
  `enum | bool | int | set | scalar`, RDR 0003's spelling" — and decides a
  CARRIER: how the declared kind reaches the value-comparison seam
  (`grammar.go::Evaluator`'s `Evaluate(atom, value) GuardResult`, constructed
  at `flow_resolve.go::guardSeam`, `product.go::valueSatisfies`,
  `reach.go::atomAdmitsValue`), and what `eq`/`in` answer per kind.
- `0021:C2` decides when a DOCUMENT field is emitted, reading
  `guard.AssignmentCount` finiteness and `decl.Domain` authored members.
- The seams are disjoint in source: the exporter constructs no evaluator
  (`grep NewEvaluator|Evaluate(` over `graph_document.go` → no match), and
  0012's source anchors touch neither `AssignmentCount` nor `decl.Domain`.
  `AssignmentCount`'s non-test callers are cardinality consumers
  (`coverage.go`, `product.go`, `lint.go`, `reach.go`), never the comparison
  seam.

So both records are downstream readers of one terminal owner's vocabulary,
over seams that do not touch; neither constrains the other. Recorded as a
citation on the Joint-check line naming 0012 and the four shared symbols.
Post-edit the pair reads `cited` on both arms, and the fence resolves
`op = none`, `rule = fence-clear`. JC1's own fire (→ 0029, home
`cli/0029 §Normative Contracts` C4) remains `homed`. `rulings_open=0` — all
nine author rulings marked absorbed.

## LOOP-BREAKER

`rdr anchors --record 0021` returns `{C1, C4}` over iter-3 and
`{A3, C1, C2, C4, JC1}` over this pass, so the `comm -12` intersection is
NON-EMPTY: `0021:C1`, `0021:C4`.

Judged, not waved through. The rule fires on a FINDING re-reported after its
named stage ran; `rdr anchors` extracts every record id a report mentions,
finding or not, so the intersection is a superset of the re-report set. In
iter-3 both ids appear only as CLOSED items — the discharge table's "Defect 1
— tiers never assigned" row recording that C1 and C2 now carry their tier
declarations, plus `0021:C4` in the 0017 fence prose and `C1-C5 labelled` in
CHECK 10. In this pass they appear only in the scope-fence row (C1/C3/C4/C5
byte-unchanged across the three commits), the same CHECK 10 PASS line, a
reference to `0012:C1`, and this section's own narration. Neither pass reports
C1 or C4 as an open finding, and no stage was asked to clear anything about
them.

So: anchor-extraction overlap, not a routing loop. No
`stopped:finalize-routing-loop` — the two ids name clauses both reports cite
while PASSING on them.

## Verdict

PASS — no blocking finding; proceed to the Gate's written responses.
`lint --locking` exit 0, `blocking=0 resolution=0 placeholder=0 advisory=1`
(C9, answered at the Gate).
