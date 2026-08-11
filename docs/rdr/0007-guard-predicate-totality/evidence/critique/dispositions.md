Model: claude-opus-5[1m]

# Critique resolve — disposition ledger (iteration 1)

Origin ledger = `critique.md`'s C-1..C-15 rows (single-model pass A; the
dual-model obligation is recorded at the bottom). Grounding gate run against code
on the working tree (branch `via-claude` — `internal/resolve/` does not exist on
`main`, which is docs-only at HEAD `991bd43`; doc claims hold on both),
`{RDR_RESOURCES}` design docs, and the RDR's own decided text.

| # | Disposition | Origin | Section touched |
|---|---|---|---|
| C-1 | **fixed** | §1, §3, premortem, AT-1 | A5 Status → Pending; *Carried constraint* rewritten from exhaustiveness to overlap; A12 added (Pending); Scenario 8 load-time precondition |
| C-2 | **dismissed-with-cite** (sub-claim **fixed**) | §1, §2, premortem, AT-2 | Prerequisites A10 — handoff destination named |
| C-3 | **dismissed-with-cite** | §1, §3, premortem, AT-3 | — (3amigo PM-4 already adjudicated; C-11 carries the live half) |
| C-4 | **fixed** | §1, §2, premortem, AT-4 | Consequences — status-vs-defect distinction stated |
| C-5 | **fixed** | §1, §2, premortem, AT-5 | Risks — residual status corrected to UNMITIGATED |
| C-6 | **fixed** | §1, §3, premortem, AT-6 | A13 added (Pending); Failure Modes — new *Refusal defeated by caller-supplied state* |
| C-7 | **dismissed-with-cite** | §1, §3, premortem, AT-7 | — (3amigo PM-1 already forced the asymmetry statement) |
| C-8 | **dismissed-with-cite** | §1, §2, premortem, AT-8 | — (refuted on source) |
| C-9 | **fixed** | §1, §2, premortem, AT-9 | Normative Contracts — strike-out condition moved out of the normative block into prose |
| C-10 | **fixed** | §1, premortem, AT-10 | Investigation — *What this citation does and does not authorize* |
| C-11 | **fixed** | §1, §2, premortem | Decision Rationale — Operability + authoring-blast-radius scored as matrix rows; narrower veto added to *Briefly Rejected* |
| C-12 | **fixed** | §1, §2, premortem, AT-11 | MVV — *What the MVV does and does not prove* |
| C-13 | **fixed** | §1, §2, premortem, AT-12 | Normative Contracts — ordering pinned to behavior, superseded rationale retired |
| C-14 | **charted-to-successor** | §2, premortem | `Charted.md` → `/rdr-finalize` Proportionality |
| C-15 | **fixed** | §3, premortem, AT-13 | Consequences — authored guards named as the first migration inventory |

## Grounding gate notes

Verified firsthand:

- **C-8 REFUTED.** All 11 `RequiresOwned` sites in `fixtures_test.go` declare only
  keys the row also `Writes` — consistent with the narrowed meaning, not the old
  one. The single divergent row cuts the *other* way: `missingOwnedTable` :140
  declares `gate`, which it neither reads in `Match` nor writes, so it is
  compatible with either reading. The fixture corpus does not disambiguate the
  two meanings at all, so it cannot contradict the narrowing. What the critique
  *should* have cited is the live doc comment (`resolve.go` :180, "names owned
  tag keys the row's **evaluation** needs") — which is the read meaning, and is
  exactly what Phase 1 exists to rewrite. The narrowing is real work; the
  fixture-contradiction claim is not.
- **C-13 CONFIRMED** with attribution corrected. The rationale is verbatim at
  `0001-…/artifacts/deviations.md` :224-226 and copied into `resolve.go::gate`
  :372-374 — but it is *not* in RDR 0001's locked text; its home is a
  `SPEC-UNDER` deviation still typed "needs author decision". So the ordering is
  pinned to shipped behavior, and the rationale that justified it is retired
  here rather than restated.
- **C-1 CONFIRMED, and stronger than argued.** RDR 0002 :342-348 puts `ambiguous
  overlap` in a ```normative``` MUST list ("at minimum"). RDR 0003 :322-326 binds
  coverage *and* overlap to one declared-domain product — but every "refuse or
  downgrade" escape in 0003 (:279-282, :316-320, :328-333, A2's closing sentence)
  attaches to the **exhaustiveness** claim and none mentions overlap. No
  `[tags.*]` block in the tree admits an absent member, and `guard-fixture.toml`
  guards `cluster_eligible` (`domain = [true, false]`) with `exists = true` — an
  operator ranging over a state the product cannot represent.
- **C-6 CONFIRMED.** `Input.Observed` :228 is caller-supplied and unvalidated (no
  allowlist, no schema check); `assemble` :157-159 writes it provenance-tagged
  into the same map, and the full view reaches `Evaluate` :380. ADV-4 :400 /
  ADV-5 :436 freeze the owned-state path only — both use `allGuardsTrue()` with
  zero named predicates, so no row in either test reaches the seam.
- **C-3 / C-12 mechanics CONFIRMED.** `gate` :396-417 returns `nil` for
  `selected` on any undecidable survivor; ordering is prune → `missingOwned(
  survivors)` → undecidable. D5's "unrelated row … poison" and D8's "converted an
  escapable condition into an inescapable one" are both verbatim in
  `deviations.md`.

Re-raise check (grounding-gate source 3):

- **C-3** re-litigates the aggregation veto that 3amigo PM-4 adjudicated, where
  the D5 tension was answered *in the draft* rather than dismissed. The critique
  adds no new evidence — it re-argues D5's text, which PM-4 already weighed. What
  it does add is the observation that no *narrower* alternative was scored; that
  is C-11 and is fixed there, and the narrower variant is now named and rejected
  on its merits in *Briefly Rejected*. Note D8 is a `SPEC-UNDER` author decision,
  not settled spec — which is precisely why this RDR ratifies it rather than
  assuming it.
- **C-7** re-raises the diagnostic asymmetry that 3amigo PM-1 already forced into
  Failure Modes in the RDR's own words. Nothing new; the gap is stated, owned,
  and routed to Phase 3.
- **C-2** re-raises the A10 hinge (3amigo IMP-1/QA-1), already converted into a
  Pending assumption blocking lock. Its one net-new observation — that "hand it
  to RDR 0003" has no destination under the never-amend doctrine — is a real
  defect in the *resolution*, not the finding, and is fixed in Prerequisites.

## Needs (re)verification — carried to Stage 6

- **A5 (was Verified → now Pending)** — the runtime half stands; the load-time
  half is refuted pending A12. Flipped because the fix touches a load-bearing
  authoring claim.
- **A12 (new, Pending)** — two-row pattern vs RDR 0003's load-time overlap check.
  Method: Source Search against 0003's overlap proof and 0002's validation
  categories. Blocks lock (added to Prerequisites) and gates Scenario 8.
- **A13 (new, Pending)** — provenance-blind presence vs caller-supplied observed
  tags. Method: Source Search against RDR 0004/0005's input boundary. Does not
  block the domain rule; blocks the claim that the refusal is unbypassable.
- **A9, A10 (pre-existing, Pending)** — unchanged by this lens. A10's handoff
  destination is now a named obligation rather than an address on a locked doc.
- **A6b (pre-existing, open by decision)** — unchanged.

No previously-Verified assumption other than A5 was invalidated.

## Tiebreakers

None escalated. The one genuinely load-bearing fork — whether C-3's veto
objection should reopen the aggregation clause — collapsed on evidence: the
narrower alternatives it implies are either absence-as-false at resolution scope
(rejected on the same grounds as Alternative 1) or uncomputable without the atom
structure A10 shows the evaluator lacks. Both are now recorded in *Briefly
Rejected* so the next lens does not re-derive them.

## Dual-model obligation — NOT yet satisfied

This RDR is `foundational`, for which `2-critique.md` requires dual-model
convergence: "A lone single-model `critique.md` with no diff does **not** complete
the lens." This file records pass A only (`claude-opus-5[1m]`). `{RDR_RESOURCES}`
notes the alt-model roster as omitted, but two evidence files for this same RDR
(`research/propose-prior-art.md`, `propose-premortem/critic.md`) are stamped
`claude-fable-5`, so a second model has been reachable here before.

Pass B must run in a fresh session on a different base model, write
`critique-modelB.md` (own stamp), and diff the ledgers **by passage anchor** (C-N
IDs do not correspond across files). Only if no alt model is reachable does the
recorded single-model fallback apply — two fresh-context runs, diffed, and noted
in `critique.md` as a fallback rather than passed off as a dual-model pass.
