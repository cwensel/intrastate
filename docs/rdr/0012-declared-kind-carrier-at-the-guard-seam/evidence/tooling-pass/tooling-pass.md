Model: claude-opus-5

# Tooling Pass — cli/0012 (Stage 7 mechanical sweep)

Date: 2026-09-20 · Iteration 1 · Verdict: **PASS** (after one in-pass fix)

Lint: `recs lint --locking 0012` → exit 0, `blocking=0 resolution=0
placeholder=5 advisory=15` (captured at `lint.txt`).

## CHECK 1 — Template section coverage

No `template:missing-section` finding: the Required spine is complete
(Metadata, Problem Statement, Critical Assumptions, Proposed Solution +
Normative Contracts, Alternatives, Context, Research Findings, Trade-offs,
Implementation Plan + MVV, Validation, Finalization Gate, References).

Five `placeholder:survived` blocks (1882-1885, 1891-1903, 1909-1911,
1922-1926, 1947-1972) and one `gate:inline` (1852-1972) — all inside
§Finalization Gate, which is unfilled by construction before this stage.
Not a hollow section: Stage 7 authors the four responses into gate.md and
the `lock` writer moves the sub-sections out behind the pointer, deleting
those blocks. Not a finding against the record.

No seed-skeleton header. No `scaffold:row`, no `contract:template-example`.

## CHECK 2 — Method-label vocabulary

All 7 assumptions: `method.off_vocabulary` empty. Labels in use — Peer RDR,
Source Search, Spike, Design Decision — all sanctioned. Every assumption has
a Method field. PASS.

## CHECK 3 — Source Search self-reference

Source-Search rows: A1, A2, A3, A4, A5. Evidence anchors resolve to
`internal/…` source symbols and to `evidence/spikes/`, `evidence/reconcile/`
paths (where spike evidence belongs). None resolves to the record itself or
to `artifacts/`. PASS.

## CHECK 4 — Docs Only on load-bearing claims

No assumption carries `method.members == ["Docs Only"]`. PASS.

## CHECK 5 — Symbol resolution of anchors

`anchors_total=78 anchors_unresolved=0 anchors_unlooked=0
peer_evidence_unresolved=0` against `$RDR_SOURCE_REPO`. No `edge:unresolved`,
no bare `file:line` anchor. PASS.

## CHECK 6 — Status consistency

Status Draft, no qualifier, no re-entry. `ca=all-terminal`: 7/7 Verified,
0 Pending, 0 Unverified — so no settled-fact prose can rest on an unverified
assumption. No checklist/gate disagreement (the gate was unwritten until this
pass). PASS.

## CHECK 9 — Evidence-field budget (BLOCKS on `foundational`)

**BLOCK, fixed in-pass.** `evidence:over-budget` at 225-258: A4's Evidence
field ran 34 lines against a 30-line soft cap, and this record's Profile is
`foundational`, where `lint --locking` marks the hit blocking.

Fix applied per the check's own remedy — relocate, never truncate. The full
derivation (held-side product path, the `reach.go::heldValues` exclusion, the
literal-side `conformKind`/`valueMembers` measurements, the 123-model corpus
coverage statement) moved verbatim to `artifacts/a4-evidence.md`. The field
retains every load-bearing anchor (`valueAssignments`, `product.go::valueSatisfies`,
`heldValues`, `atomAdmitsValue`, `conformKind`, `valueMembers`), the spike and
reconcile pointers, and a pointer to the artifact. Re-lint: exit 0, blocking=0.

## CHECK 10 — Linking

No `label:contracts` / `label:contracts-required`: C1-C5 are labelled. No
`peer-evidence:no-element` — Peer RDR evidence cites elements
(`0001:§d1`, `0004:§jd-2`, …). No `edge:unresolved`,
no `edge:unresolved-terminal`. PASS.

## Advisory carried, not fixed

Nine `prose:exactness` hits on "canonical" (815, 833-834, 862-899). The term
is not loose here: C5 defines it operationally as
`strconv.Itoa(n) == authored`, and it is pinned by named fixtures — S6 (the
123-model corpus diff) and S7 (`n eq "00"` regression), both in the MVV and
the conformance suite. Advisory tier; carried.

## Also fixed in this pass (mechanical, per Stage 7's split)

Two further in-pass fixes, both fillable from material the record already
carried, so neither is routed to Refine:

- **CHECK 6 — internal status contradiction.** §Key Discoveries' match-path
  bullet still described the pre-ruling state and closed with "A4 carries
  the open question of whether the held leg must be closed too … or the
  seam's match arm must byte-compare." A4's reconcile ruling (option (ii))
  settled that question, and C4 implements it. Rewritten to the settled
  state: the match path is closed by EXCLUSION from typing, not by
  canonicalization, and the surviving held-ingress gap is identical on both
  sides and carried by kata `intrastate#ch99`.
- **Fence — `overlap_uncited=1+`.** `index --anchor-intersect` reported two
  uncited open-peer pairs, 0013 and 0022. Both were adjudicated, not
  synced: 0013's own `JC1` had already fired this joint check and ruled it
  disjoint ("0012 modifies `grammar.go::Evaluator` / `valueSatisfies`,
  disjoint from this seam's entries") — recorded here symmetrically; 0022
  reads the two functions C4 excludes and relies on their byte-based
  behaviour, which is exactly what C4's exclusion guarantees, with the
  dependency running one way. Citations landed in §References and
  §Load-Bearing Decisions (NOT inside the ```normative fence, which emits
  no reference edge). `overlap_uncited` now 0; the fence emits
  `op = none`.

---

VERDICT: **PASS** — no blocking finding; proceed to the Gate's written responses.
