Model: claude-opus-5[1m]

# Stage 6 Reconcile — RDR 0021 (lint's normalized-graph export)

Date: 2026-09-13 · Iteration 4 · Verdict: **RECONCILED**

Iteration 3 closed RECONCILED. The record went to Finalize, locked, was
implemented at Stage 8, and was routed back by the author on 2026-09-13 over
`C2`/REQ-20's `domain` presence rule — unsatisfiable as written, because
"finite" and "has authored members" are different predicates. Refine narrowed
REQ-20 (`5154721`) and resolve re-verified A3 against `domainSize` and cleared
the qualifier (`5fe32cf`).

This pass reconciles that **delta**: whether the narrowing landed in the
contract and whether it matches the code that was already shipped against the
old wording. A1–A8 are Verified; only A3 was named by the route-back.

## Stage 5 preflight

| Check | Answer | Note |
| --- | --- | --- |
| `--outcome lens` | `/rdr-reconcile` | `lens-large-row-complete`; Profile `large`, grounding→3amigo→critique row complete |
| `--outcome critique` | `none` | `critique-large-diffed`; `critique_models=differ` — no single-model fallback |
| `--outcome repeatability` | `none` | `repeatability-lite-complete`; `variant: lite (profile: large)` stamped on `run-1.md` |
| Determinacy add-on | `determinacy=fired`, line written | `Determinacy: fired — C3 states byte-for-byte identity…` (§normative-contracts:213). No `stopped:determinacy-trigger-unjudged`; the iteration-1/2 chain disagreement did not recur. |

Status is plain `Draft` (`status_form=none`): the `re-verify A3 @refine`
qualifier was cleared at `5fe32cf`, so nothing is owed to this stage by
qualifier.

## The open set

All four sources are empty on disk; no Pre-Lock list accompanied the
invocation.

| Source | State |
| --- | --- |
| 1 — Pre-Lock needs-verification list | None pasted; `lens_findings_open=0`, `lens_stale=none`, `rulings_open=0` (all 10 rulings carry `absorbed @`) |
| 2 — Pending/Unverified CAs | Empty: `ca_total=8`, `ca_verified=8`, `ca_pending=0`, `ca_unverified=0`, `ca_off_vocabulary_ids=[]` |
| 3 — spikes named but unrun | `spikes_unrun=[]`. Six spike files on disk (a1, a2, a5, a6, a8, a9); the absorption audit found no findings file naming a spike the record does not. `a9-opaque-sentinel.md` survives for an assumption the record no longer carries — retained deliberately as the evidence behind ledger row 1 (iteration 3's disposition) |
| 4 — exactness-word delta | Empty: `rdr lint 0021 \| grep prose:exactness` returns nothing |

So the open set is the **route-back delta** plus the one obligation iteration 3
booked out.

## Dispositions

| # | Item | Source | Disposition | Evidence pointer / plan |
| --- | --- | --- | --- | --- |
| 1 | A3 / REQ-20 — `domain` presence rule, the Stage 8 route-back's subject | route-back (`re-verify A3`) | **VERIFIED — record and code agree** | C2 now binds `domain` to the AUTHORED members: "present exactly when there are any — an `enum` with a non-empty `domain` — and ABSENT, its key omitted, otherwise", with `bool`, `int`, `set`, `scalar` and a member-less `enum` all omitting it, and one sentence naming why finiteness is the wrong predicate ("Presence is NOT `guard.AssignmentCount` finiteness"). Delegated source search (PASS): the exporter gates on `finite && len(decl.Domain) > 0` (`internal/cli/graph_document.go:185`) — the authored-members reading; `domainSize` confirms `enum` finite only with a non-empty `Domain` (`declaration.go:116-119`), `bool` `spread(2, …), true` (`:123`), `int` `spread(width, …), true` (`:128`, ceiling arm `:136`), `set` `spread(len(d.Elements), false), true` (`:145`), `scalar` `0, false` (`:149`) — three kinds finite with no `Domain`. No live refutation: the clause the route-back refuted has been withdrawn and replaced by the reading the code implements. Detail: `a3-domain-presence.md`. A3's Evidence records the re-verification this pass. |
| 2 | No derived value vocabulary minted for the omitting kinds | route-back (resolution direction) | **VERIFIED — discharged in C2** | The re-entry note required that no new wire vocabulary be minted, none being backed by a REQ-N. C2 states it: "neither `bool`'s two literals, nor a `set`'s `Elements` or its powerset, nor an `int`'s `{min..max}` bound is projected into this member", and records the rejected alternative — "A compact bound form for `int` was considered and rejected — adding it later is additive under this clause's own `/1` rule, removing it is not." |
| 3 | Coverage gap — `bool`, bounded `int`, member-less `enum` unasserted | route-back (named gap) | **VERIFIED — oracle stated in S2** | Testing Strategy scenario 2 now asserts both sides: "present for an `enum` with members, and ABSENT for `bool`, for a bounded `int`, and for a member-less `enum` — the three kinds `guard.AssignmentCount` reports finite while authoring no `decl.Domain`, so a test asserting only the enum arm would pass against the pre-narrowing finiteness rule too." The record owes the oracle; the shipped test does not yet meet it — item 5. |
| 4 | C-7 — two-marker read order (open joint decision, home `0029:C4`) | carried from iter-3 | **VERIFIED — the home now answers it** | Iteration 3 left this as the one obligation, with `0029:C4` covering the adjacent point but never sequencing the two reads. C4 now does, verbatim: "A consumer reads the envelope's major to decide whether it can parse at all, and the document's `schema` to decide what the payload means. That ordering is normative", and settles both mixed cases — unsupported envelope major is rejected under C1 before `data` is read; supported envelope with an unrecognized document marker leaves the envelope trustworthy and `data` opaque. `0021:JC1` promised to "cite the home once C4 states it"; the citation now resolves. 0029 is `Implemented`. Detail: `c7-joint-decision.md`. |
| 5 | Implementation-side defects on the unmerged branch | item 1's search | **DOWNGRADED — survivable, named plan, implementation's to fix** | Two defects in `worktree-rdr-0021` @ `143d919` (no graph files on `main`): (a) the doc comment at `internal/cli/graph_document.go:51-54` still asserts the refuted finite-only claim, "the finite arm always carries at least one member", contradicting the `:185` gate it documents, with the same stale framing repeated at `:180-181`; (b) `TestReq20And28` asserts only the enum-with-members (`:229-233`) and `scalar` (`:239-243`) arms — `recognized` (member-less `enum`) is decoded but never asserted, and `bool`/bounded `int` are absent from the fixture. Neither is a record defect: the record's text is correct and S2 already states the owed oracle. Survivable because the shipped behavior is right and the golden already pins it (no `domain` for `bool flag` or member-less `enum recognized`, `["a","b"]` for `enum status`); the exposure is a misleading comment and a test that would pass against the pre-narrowing rule. Plan: both fixed in the implementation pass before the branch merges, asserted against S2 scenario 2. |
| 6 | Round residue across all pre-lock rounds | 1, 3 | **ABSORBED** | Delegated absorption audit (PASS, non-blocking) over grounding, 3amigo (base + iter-3), critique, repeatability: every finding absorbed into current text. F3's base fix was withdrawn with A9 at cluster reconcile, but C2 now declares no terminal marking and cites RDR 0015 / JDR 0001 §JD-23 as owner — a superseding absorption, not residue. Detail: `absorption-audit.md`. |
| 7 | Hollow-body / References completeness | prompt gate | **PASS** | Grep, not judgment: no `_Draft placeholder._`, no `this is a seed skeleton` header. `## References` fully authored (peer cites, source paths, prior art, kata) with no bracketed placeholders. The one bracket hit in the body is `[Gate key: cross-cutting]` at :1191, a gate-section label. `rdr lint` reads `blocking=0 resolution=0 placeholder=0 advisory=1`. |

## The two hard rules

**Refutation.** No live refutation, and this pass tested the rule directly
rather than assuming it. The Stage 8 route-back WAS a refutation — REQ-20's
presence predicate was unsatisfiable against `domainSize` — and it was handled
the way the rule demands: the clause was narrowed to the predicate the
authority actually supports, not reworded to hide the gap. The delegated search
then re-ran the falsifier against the code and found agreement
(`graph_document.go:185` vs C2's authored-members rule). The refuted fact
survives verbatim in the escaped-defect ledger (row 5), which is its durable
home now that refine has collapsed the change history. No new route-back is
issued, so no new §punt-ledger row is owed this pass.

**MVV dependency.** Nothing deferred is MVV-critical. The MVV's five steps
assert fixture load; double-emit byte identity plus PRESENCE of every C2 field;
DOT node/edge id-set equality; `jq .data` value equality; and lint neutrality.
Item 5's two defects pin neither: MVV step 2 asserts that the `tags` member is
present and byte-stable, never which declaration kinds carry a `domain`
sub-member. The presence rule's own oracle is S2, an in-scope scenario with
assertions on both sides of the rule and a stated negative control ("a test
asserting only the enum arm would pass against the pre-narrowing finiteness
rule too"). So the downgraded item has a named in-scope test, and nothing the
MVV proves depends on it.

## Obligations leaving this stage

| Owner | Obligation | Consequence while unlanded |
| --- | --- | --- |
| Implementation (`worktree-rdr-0021`) | Correct the doc comment at `internal/cli/graph_document.go:51-54` (and the `:180-181` framing) to the authored-members reading, and extend `TestReq20And28` to assert `domain` ABSENCE for `bool`, bounded `int`, and member-less `enum` per S2 scenario 2. | The shipped behavior is already correct and the golden pins it; the exposure is a comment that contradicts its own gate and a test that cannot distinguish the narrowed rule from the one it replaced. Branch is UNMERGED at `143d919`. |

Iteration 3's obligation on `0029:C4` is **discharged** (item 4). No obligation
leaves this stage for any peer record.

## Caveats

- **The implementation is unmerged and predates the narrowing.** All
  implementation paths (`internal/cli/graph_document.go`, its test, the
  fixture) exist only on `worktree-rdr-0021` @ `143d919`; `main` carries only
  `internal/guard/declaration.go`, byte-identical on both refs. The code was
  written against the OLD REQ-20 wording and happens to implement the NEW one —
  which is why the route-back called the shipped gate "the correct behavior".
  The record is now the thing that changed; Stage 8's re-entry reconciles the
  comment and the test to it.
- **`evidence:over-budget` advisory on A3 (177–219) — created by this pass.**
  Appending the re-verification took A3's Evidence field to 43 lines against
  TEMPLATE.md's soft cap of 30. `conformance` tier, non-blocking, and the
  lint's own fix text routes the call to the Gate rather than to a rewrite:
  answer there whether the load-bearing anchor is still findable in the field,
  and whether the balance belongs in `{ARTIFACT_DIR}` with the field keeping
  the anchor and a pointer. Deliberately NOT truncated here — the mass is the
  `domainSize` per-kind verification the grounding sweep reads, and A3 is the
  assumption the route-back named. Flagged as this iteration's addition so
  Finalize does not read it as inherited.
- **`section:unknown-to-template` advisory (1283–1341).** The `## Refinement
  Context (re-entry — delete on re-lock)` block is the only lint finding. It is
  self-marked for deletion at re-lock — Stage 7's write, not a reconcile
  regression. It carries the Stage 8 defect narrative and the branch pointer, so
  it must not be dropped before Finalize reads it.
- **A9's spike file remains on disk** (`evidence/spikes/a9-opaque-sentinel.md`)
  for an assumption the record no longer carries — retained across iterations 3
  and 4 as the evidence behind ledger row 1.

## Verdict

**RECONCILED** — every item terminal, no BLOCKER. One DOWNGRADED item (5), on
the implementation branch rather than the record, with a named plan and a
stated in-scope oracle. Ready for Finalize.
