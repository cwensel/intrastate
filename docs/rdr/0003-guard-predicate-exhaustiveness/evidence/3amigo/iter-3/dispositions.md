Model: claude-opus-5[1m]

# 3amigo Dispositions — iteration 3 (re-entry)

Origin ledger = this iteration's `consolidation.md` (14 entries over 21 persona
findings). Every finding traces to a ledger entry; none was net-new scope.

## Grounding gate

All 21 findings passed. Verified verbatim before any edit:

- `internal/resolve/resolve.go` — `Row` carries `Match []Tag` (:178) and
  `Guard string` (:185) as **separate** fields, and `Resolve`'s documented
  evaluation order states "candidate rows are the **non-escape** rows for that
  outcome whose match pattern holds" (step 2), with the guard gate at step 3.
  **P2-3 and P2-4 ground on shipped code.** `RuleID`/`SourceLocator` at
  :173-174 confirm A5's re-verification.
- `docs/rdr/0002-transition-table-as-reviewable-data.md` — `single-valued`
  occurs **zero** times; the type-model clause (`:342-350`) and schema list
  (`:217-220`) each enumerate four fields. **P1-2, P2-7, P3-7 ground.** Its
  escape clause (`:328-334`) admits a modeled escape disposition only when
  ordinary resolution has already failed and exactly one escape row matches —
  **P2-4 grounds on the peer's normative text too.** Its carriage requirement is
  already general ("every declared field through without loss"), so only the
  authoring enumeration is short.
- `docs/rdr/0006-graph-lint-authority-and-guarantees.md` — invariant 5
  (`:295-296`) constrains **writes** over a **tag class**; the lint input
  contract (`:260-267`) reads "single-valued grouping when applicable" from this
  RDR's model. **P2-5 and P3-2 ground.** Its Status line already carries
  "§JD-13's single-valued citation" and the §JD-14 invariant 3/4 repair as owed
  refine duties, which is the venue for both consumer-side halves.
- `docs/jdr/0001-resolve-kernel-seam.md:321-337` — §JD-14's decision text
  reasons from "RDR 0001's runtime refusal of ambiguity". **The premise is the
  defect** (see A17 below).
- `docs/rdr/0003-.../evidence/spikes/guard-fixture.toml` — declares seven tags,
  **none** carrying an optionality marker, and authors `status eq "Draft"` under
  `[rule.match.status]`. **P2-2 and P2-3 ground.**
- Cluster-wide grep: `conforming evaluation view` occurs exactly twice, both in
  the two new re-entry clauses, with no definition and no violation outcome.
  **P3-1 grounds.**

No finding was dismissed for failing the gate. No finding re-litigated an
already-adjudicated option — the one that touches a decided call (§JD-14) does
so on evidence the decision did not have, and is booked as an assumption rather
than silently reversed.

## Compute, don't argue

- **Assumption census** — the roster was rewritten from an executed extraction,
  not a claim. Parsing `- **A<N>` blocks for their `Status` line returned
  11 `Verified` (A1–A9, A11, A13) and 6 `Pending` (A10, A12, A14, A15, A16,
  A17). The roster now reads that computed set.
- **Four-vs-five-field split** — `grep -n 'single-valued'` returned five-field
  sites at `:17`, `:543`, `:687`, `:715`, `:729`, `:937`, `:940` against
  four-field sites at `:1055`, `:1740-1742`, and the Operator-semantics bullet.
  The subtraction is exactly ledger entry 2; all three short sites were fixed.
- **Single-valued product effect** — the clause previously said only what the
  marker "licenses". The two readings are now written as computed cardinalities:
  with the marker, one dimension of `|domain|` assignments; without, one boolean
  dimension per value (`2^|domain|`). For enum `{a,b,c}` that is 3 vs. 8 — which
  is what makes MVV Scenario 5's two verdicts differ by construction rather
  than by assertion.
- **Cross-lens check** — the `authority`, `disposition`, and `trace` tables were
  written by earlier lenses in this row and all three were re-read against this
  pass's clause edits. Three disagreements were found and fixed, not argued:
  `trace` step 3 counted a **match** key as a product dimension (contradicting
  the participation clause it was cited to justify), step 7 assumed one overlap
  population, and the `authority` census had no row for either fact.

## Dispositions

| Finding(s) | Ledger entry | Disposition | Section touched |
| --- | --- | --- | --- |
| P1-2, P2-7, P3-7 | 1 | **fixed** — new assumption **A16** | Critical Assumptions (A16, `Pending`, Method: Peer RDR); Capability Dependencies (own row); Prerequisites (§JD-13 checkbox → "Done for meaning", plus an A16 entry); MVV Authorability; Existing Infrastructure Audit |
| P1-1, P2-7 | 2 | **fixed** | Capability Dependencies decl-model row; Scope Verification; Load-Bearing Decisions (new `Declaration model` bullet fixing all five fields and requiring every enumeration to state them) |
| P3-1, P2-1, P2-2 | 3 | **fixed** | Normative Contracts — optionality clause gains the **default** (no marker = optional, the conservative arm) and a new paragraph **defining "conforming evaluation view"** with its violation routed to the kernel/RDR 0006 rather than left silent |
| P2-1 | 4 | **fixed** | Normative Contracts — "can refuse" restated over the **declared optional** field (one property, not three names); owned-tag graph condition explicitly separated as A12's subject; A6 Note and A12 Evidence swept to match |
| P2-3 | 5 | **fixed** | Normative Contracts — participation clause pinned to **guard** atoms, with a new paragraph excluding match keys on RDR 0002's authored shape and `resolve.go::Row`; `trace` step 3 witness corrected to `prelock_iterations` |
| P2-4 | 6 | **fixed** — new assumption **A17** | Normative Contracts — escape-row clause rewritten to **two overlap populations**; Critical Assumptions (A17); Prerequisites (§JD-14 checkbox re-opened, A17 entry); Contradiction Check; Metadata Status line; `disposition` (2 new rows); `trace` step 7; `authority` (new row) |
| P2-5, P3-2 | 7 | **fixed** | Normative Contracts — single-valued clause gains the **computed product effect** and a shape note naming RDR 0006's per-class/per-write mismatch as its refine's reconciliation; Testing Strategy gains **Scenario 5** (the acceptance test) |
| P2-6, P3-4 | 8 | **fixed** | Normative Contracts — five kind **tokens pinned** (`enum`/`bool`/`int`/`set`/`scalar`), the fifth kind named and given its domain rule; operator/kind matrix rewritten to those tokens; `exists` row restated so an always-present key is vacuous rather than rejected |
| P3-3, P2-9 | 9 | **fixed** | Normative Contracts — domain/kind clause names the kind **`literal-outside-declared-domain`** and splits the **two phase boundaries** with an owner each; A4's kind enumeration swept; Testing Strategy Scenario 4 gains the case; `disposition` gains 2 rows |
| P2-8 | 10 | **fixed** | Normative Contracts — `{min..max}` marked notation-not-wire-spelling, authored as `min`/`max`, both endpoints **inclusive** (so `{0..3}` has cardinality 4) |
| P3-5 | 11 | **fixed** | Testing Strategy — new **Scenario 7**, the post-normalization `block`-retention assertion A14's survivability argument names |
| P3-6 | 12 | **fixed** | Testing Strategy Scenario 3 — the shape pair is now constructed **relative to the published bound B** (equal cardinality C > B, then just under B), which is what makes the A15 comparison runnable |
| P2-10, P2-11 | 13 | **fixed** | Phase-gating table Phase 2 row — startability split: arithmetic/withholding/partition startable, group construction and escape overlap wait on A10/A17 |
| P1-3 | 14 | **fixed** | Illustrative Code — a `[tags.<tag>]` declaration example carrying all five fields, beside the existing guard block, with the `single_valued` spelling flagged as A16's to fix |

All 21 findings **fixed**. None dismissed, none charted to a successor.

## Tiebreaker reduction

Three forks collapsed on evidence rather than escalating:

1. **Do match keys enter the scoped product?** Apparent fork: honor the
   participation clause (guard-only) or the `trace`'s worked example (match key
   included). Collapsed — the group is *defined* by its selection context, so a
   match key constrained identically by every row is the grouping key, not a
   dimension within the group; it would contribute exactly one assignment,
   multiply every product by 1, and misreport the cardinality the too-large
   bound is measured against. RDR 0002's authored shape and `resolve.go::Row`
   both keep `Match` and `Guard` separate, and `Resolve` filters on match before
   the guard gate. The clause was right and its witness was wrong.
2. **Is `exists` over an always-present key a rejection or a legal vacuity?**
   Apparent fork: the matrix said `exists` accepts only optional tags; the
   presence-dimension clause said an always-present key contributes no
   `{absent}` assignment. Collapsed — the latter is the load-bearing one (it is
   what makes the narrowing's negative control provable rather than vacuous),
   and rejecting the atom outright would make the negative control unwritable.
   The matrix row now states well-formed-but-vacuous.
3. **Does the single-valued marker need a new field type?** Apparent fork:
   per-tag boolean (this RDR) vs. per-tag-class id (RDR 0006's invariant 5).
   Collapsed for *this* document — the marker is a property of a tag's declared
   domain, and every sibling field in the model is per-tag; a per-class grouping
   is a different predicate over writes, which is RDR 0006's invariant to
   reconcile at the refine it already owes. Recorded as a shape note, not a
   redesign.

One fork was **not** collapsed unilaterally and is booked instead: §JD-14's
overlap half. The evidence (kernel + RDR 0002) is strong enough to change this
RDR's clause, but a gate decision is a cross-document artifact, so **A17**
carries the confirmation to the venue that already owes the related repair
rather than this pass silently reversing a recorded decision.

## Needs (re)verification (carried to Stage 6)

- **A16** (new, `Pending`, Method: Peer RDR) — RDR 0002's authoring surface must
  carry the single-valued marker. Travels with A14 as one authoring/carriage
  request. Gates MVV Scenario 5's authoring, not this RDR's contract.
- **A17** (new, `Pending`, Method: Peer RDR) — escape-row overlap population;
  re-opens §JD-14's overlap half against the shipped kernel. Venue: RDR 0006's
  refine (which owes the invariant 3/4 repair) plus a §JD-14 correction at the
  next JDR 0001 touch.
- **A12** (carried, `Pending`) — **narrowed**: it now gates the owned-tag clause
  only. The "can refuse" clause reads the declared optionality field, so the
  withholding decision no longer waits on a predecessor relation. Recorded in
  A12's Evidence and A6's Note.
- **New load-bearing claims added this pass**, each grounded at its authority or
  carried by an assumption: the optionality **default** (conservative arm,
  grounded on the symmetric domain rule already in the draft); the
  **conformance** definition (states the premise of every lint claim, with
  enforcement routed to the kernel and RDR 0006 rather than asserted here); the
  **guard-only participation** rule (grounded on `resolve.go::Row` and RDR
  0002's authored shape); the **two overlap populations** (A17); the
  **single-valued product cardinality** (`|domain|` vs. `2^|domain|`, an
  arithmetic consequence, tested by Scenario 5); the **five pinned kind
  tokens** and the `scalar` domain rule; **`{min..max}` inclusivity**; and the
  **`literal-outside-declared-domain`** kind (swept into A4's enumeration).
- No previously `Verified` assumption was invalidated. A5, A7, A8, A11 and A13
  are untouched by this pass's edits.

## Amendment sweep (rdr-common §amendment-sweep)

1. **Pre-edit token grep** — `single-valued` (17 sites) reviewed: 3 short
   enumerations fixed, the rest already correct. `can refuse` (7 sites): the
   clause plus two A6/A12 sites swept to the narrowed reading; the `disposition`
   and `trace` rows re-read and still accurate. `escape` (18 sites): the Metadata
   Status line's "no edit owed" was stale and is fixed; §JD-14's checkbox
   re-opened; the two prose sites at `:91`/`:123` are unrelated (research
   findings about escape hatches as a design idea).
2. **Subject-token grep** — `participat` re-read across the participation
   clause, `trace` step 3, and the `authority` census; producers and consumers
   now agree that the product is built from guard keys. `conforming` re-read: the
   two clause uses are now backed by the definition paragraph.
3. **Added clause amends its contract block** — the new/edited normative clauses
   were re-read against their siblings. The conformance paragraph reuses the
   existing evaluator-scope clause rather than minting a second authority over
   the tag view; the two-population overlap clause was checked against the
   report-every-defect clause (both populations report every defect, so no
   conflict) and against the coverage-union clause (unchanged); the pinned kind
   tokens were checked against the finite-domain clause, which now covers
   `scalar` explicitly.
4. **New order/reachability/emission claim** — the guard-only participation rule
   is an emission claim about which dimensions enter the product; grounded at its
   authority (`resolve.go::Row`'s `Match`/`Guard` split and RDR 0002's authored
   shape) rather than asserted. The two-population overlap claim is registered as
   A17 rather than asserted as already agreed.

## Mini-checks

Cue read performed against the revised draft. The three tables the draft carries
— `authority`, `disposition`, `trace` — stay fired and were **all three updated
by this pass** (2 new `authority` rows, 4 new `disposition` rows, corrected
`trace` steps 3 and 7 plus an unexercised-clauses note). No new cue fired: the
fixes add no fallback path, derived output, oracle, or encode/decode pair.
`round-trip / fidelity` remains unfired by design (`Round-Trip / Inverse
Invariants`: "This RDR introduces no encode/decode pair"), and
`test-discriminability` stays discharged by the MVV's named negative controls —
Scenario 5 is itself a discriminability test, asserting the two verdicts must
*differ*.

**Mini-checks fired: `authority`, `disposition`, `trace` (all pre-existing;
updated, not newly fired).**

## Charted to successor

None. Every finding was in delta-scope and fixed in place; the two that reach
other documents (A16, A17) are booked as assumptions with named venues rather
than charted, because both are producer requests this RDR depends on rather than
net-new scope it declined.

## Tiebreakers escalated

None — three collapsed on evidence (above), and the one genuine cross-document
call (§JD-14's overlap half) is booked as A17 for its existing venue rather than
escalated as a blocking fork.
