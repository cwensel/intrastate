Model: claude-opus-5[1m]

# Critique Dispositions — iteration 2 (re-entry, dual-model)

Origin ledger = `diff.md`'s reconciled ledger (R-01..R-20), built by a barrier
sub-agent that authored neither pass. `C-N` IDs do not correspond across
`critique.md` (Model: claude-opus-5, 14 rows) and `critique-modelB.md`
(Model: claude-sonnet-5, 10 rows); reconciliation was by passage anchor.

## Dual-model draw

Two distinct stamps, fanned out in parallel per rdr-common §auto-fanout, diffed
behind the barrier. The disagreement was the signal and it was substantive:
pass B **accepted three `Verified` assumption stamps that pass A attacked**
(A1, A2, A6). Independent verification of pass A's six harshest claims returned
**zero REFUTED, five VERIFIED, one PARTIAL** — so the demotions are grounded,
not hostile framing. `large` profile: dual-model is strongly recommended and was
achieved; no single-model fallback recorded.

## Grounding gate

Every actioned finding was verified against source before any edit.

- `evidence/spikes/check.sh` — read verbatim: three awk passes matching operator
  *spellings* in TOML section headers; the `coverage=` line is a hardcoded
  `printf` constant and the report loop is hardcoded to
  `split("eq in lt gte exists")`. It evaluates no predicate. **R-06 grounds.**
- `docs/rdr/0002-*.md` Technical Design §2 — tag declaration is "tag name,
  provenance, value kind, and optional accessor reference". No enum, range, or
  domain field. **R-01, R-11, R-15 ground.**
- `docs/rdr/0002-*.md` Normative Contracts — "Normalization MUST combine both
  into one candidate-row predicate set"; retention promise covers rule id and
  locator only. Constrains the *container*, not per-atom fields. **R-07 grounds
  as PARTIAL** — a producer-contract gap, not the contradiction pass A claimed.
- `docs/rdr/0007-*.md:1249-1266` — atom fields are `Key`, `Operator`, `Literal`,
  `Block`; no source-identity field. **R-13 grounds.**
- `docs/rdr/0007-*.md:1585-1595` — assigns canonical set-literal spelling to
  this RDR, naming `in`. **R-14 grounds** (distinct from A9's element universe).
- `docs/rdr/0006-*.md:277-279` — coverage counts "one modeled row **or a
  declared escape row**"; `grep escape` over RDR 0003 returns 4 hits, none about
  rows. **R-08 grounds.**
- `docs/rdr/0006-*.md:352-356` — finding contract carries code, model identity,
  severity, message, source rule/context id or span. **No atom field. R-09
  grounds.**
- `row group` — 0 occurrences in RDR 0006, whose A2 delegates the coverage
  derivation back here. **R-02 grounds: the delegation cycle is real.**
- RDR 0003 `:592` — "Semantic equality is `(tag, operator, literal)`" beside the
  identity tuple, with no clause selecting which operation uses which. **R-17
  grounds** (pass B's strongest solo finding; pass A read the passage twice and
  never asked).

Nothing was dismissed for failing the gate. No finding re-litigated a decided
option.

## Dispositions

| Ledger | Disposition | Sections touched |
| --- | --- | --- |
| R-01 (A2 has no producer) | **fixed** — A2 `Verified` → `Pending`; new **A11** books the parent domain field | Critical Assumptions; Prerequisites; Contradiction Check; Assumption Verification |
| R-02 (row-group delegation cycle) | **fixed** — row group defined here; RDR 0006 supplies reachability only; new **A10** books peer confirmation | Technical Design; Critical Assumptions; Prerequisites |
| R-03 / H-1 (default-on vs. Loud) | **fixed** — "undeclared = opt-out" removed; undeclared dimension is the Loud blocking outcome, matching `disposition` row 8 | Technical Design |
| R-04 (reachable predecessor undefined) | **fixed** — split out as new **A12** | Critical Assumptions; Prerequisites |
| R-05 (A6 over-stamped) | **fixed** — A6 narrowed to labels-only; reachability half → A12 | Critical Assumptions; Assumption Verification |
| R-06 (A1 spike proves nothing) | **fixed** — A1 `Verified` → `Pending` with evaluation-harness plan | Critical Assumptions; Prerequisites; Assumption Verification |
| R-07 (block retention) | **fixed as producer request** (PARTIAL verdict honored — booked as **A14**, not litigated as a contradiction) | Capability Dependencies; Critical Assumptions |
| R-08 (escape rows) | **fixed** — new normative clause placing escape rows inside the coverage union and overlap checks | Normative Contracts |
| R-09 (atom payload on `Final` peer) | **fixed** — registered as an explicit A8 peer obligation rather than assumed | Normative Contracts (note block) |
| R-10 (gate audits label hygiene) | **fixed** — `[x] A1-A6 verified` removed; the section now says why conforming ≠ verified | Prerequisites; Assumption Verification |
| R-11 (MVV unauthorable) | **fixed** — authorability gate added; MVV stays in scope, schedule bound to A11/A7/A9 | MVV; Scope Verification |
| R-12 (four external blockers) | **fixed** — phase-gating table names what blocks each phase | Implementation Plan |
| R-13 (atom fields restated wrong) | **fixed** — cites RDR 0007's four fields; states source identity is row-level, not an atom field | Technical Design |
| R-14 (set-literal spelling for `in`) | **fixed** — new **A13**, owned here, distinct from A9 | Critical Assumptions; Prerequisites |
| R-15 / R-16 (A7/A8 gate no phase) | **fixed** — folded into the phase-gating table | Implementation Plan |
| R-17 (which equality) | **fixed** — identity for diagnostics/dedupe, semantic for domain computation; overlap is set intersection, not an equality test | Load-Bearing Decisions › Identity |
| R-18 ("too large" undefined) | **fixed** — new normative clause requiring a declared, model-independent, published bound | Normative Contracts |
| R-19 (cross-RDR error rendering) | **dismissed-with-cite** — RDR 0005 owns CLI envelope mapping (A4); this RDR explicitly disclaims it. Diff ranked NOISE. | — |
| R-20 (CON-1 ordering) | **fixed as two records** — A's producer gap (A11) precedes B's threshold gap (R-18); both booked | Critical Assumptions; Normative Contracts |

18 fixed, 1 dismissed-with-cite, 1 resolved as two records. None charted.

## Tiebreaker reduction

Three apparent forks collapsed on evidence:

1. **CON-1 — is A2's defect the missing producer or the unbounded product?**
   Collapsed on *ordering*: if no domain can be declared, no author reaches a
   product of any size. Pass A's gap precedes pass B's in time, and pass B's
   28,800-cell scenario silently assumed the `domain`/`min`/`max` schema R-01
   proves absent — it was read off this RDR's own spike fixture, the exact
   circularity A names. Both booked; neither dropped.
2. **CON-2 — one set-literal gap or two?** Collapsed on authority: RDR 0007
   files the canonical-spelling duty on *this* RDR naming `in`, while A9's
   element universe is a request to RDR 0002 for `contains`. Different owner,
   different operator, different failure. Two records (A13, A9).
3. **H-1 — does the default-on reading fail silent or loud?** Not a contest: the
   two passes predicted opposite symptoms because the draft genuinely said both
   things — Technical Design made an undeclared dimension a silent opt-out,
   `disposition` row 8 made it Loud. Two lenses had pinned one field two ways.
   Resolved to Loud, the reading the contract block already enforced.

None escalated.

## Needs (re)verification (carried to Stage 6)

Demoted from `Verified`:

- **A1** (`Pending`) — the recorded spike cannot support the claim. Closable
  here by re-running it as an evaluation harness.
- **A2** (`Pending`) — sound derivation, no producer. Blocked behind A11.
- **A6** — narrowed to provenance labels; still `Verified` for that scope only.

New records, all `Pending`:

- **A10** — RDR 0006 accepts the row-group division (cluster reconcile).
- **A11** — RDR 0002 carries a finite-domain field. **Parent of A7 and A9**;
  blocks lock. Live request — RDR 0002 is `Draft`.
- **A12** — predecessor reachability is decidable (cluster reconcile).
- **A13** — canonical set-literal spelling. **Owned here**, closable on this
  RDR's own initiative, and it affects `in`, inside the claimed-proven subset.
- **A14** — normalization preserves atom `block` (RDR 0002 producer request).

A11, A7, A9, A14 travel to RDR 0002 as **one** four-field producer request.
A8, A10, A12 travel to cluster reconcile as RDR 0006 agreements. A5 stays
`Verified` — R-13 sharpened how the RDR cites the seam without disturbing its
evidence.

## Amendment sweep (rdr-common §amendment-sweep)

1. **Pre-edit token grep** — `opts out` → 0 surviving sites (1 remaining hit is
   the retained "opt-in flag would let the guarantee be skipped" rationale,
   still accurate). `A1 through A6`/`A1-A6` → 0 sites after the Prerequisites
   and Assumption Verification rewrites. `travels with A7` → re-pointed at A11.
2. **Subject-token grep** — `source identity` (3 sites) re-read against R-13:
   A5's two uses are row-level identity, which is exactly what the corrected
   Technical Design states; the Consequences entry is unaffected. No stale site.
3. **Added clause amends its contract block** — the three new normative clauses
   (escape rows, declared bound, and the R-09 note) re-read against their
   siblings. The escape-row clause is compatible with the scoped-row-group
   clause it follows; the bound clause tightens the existing too-large clause
   rather than contradicting it; the R-09 note qualifies the withheld-claim
   clause without weakening the runtime veto.
4. **New order/reachability/emission claim** — the row-group definition is a new
   reachability-adjacent claim. Grounded at its authority (RDR 0001's exact-one
   selection set) **and** registered as A10 rather than asserted as agreed. The
   `reachable predecessor` quantifier was *removed* from assertion and booked as
   A12.

## Mini-checks

The cue read was performed at iteration 1 / the grounding pass; the `authority`,
`disposition`, and `trace` tables persist in the draft. This pass updated the
`disposition` table's consistency with Technical Design (H-1) but added no new
cue, so no new table is owed.

## Charted to successor

None.

## Tiebreakers escalated

None — all three collapsed on evidence.
