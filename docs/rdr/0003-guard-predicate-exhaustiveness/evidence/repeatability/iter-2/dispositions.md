Model: claude-opus-5[1m]

# Repeatability-lite resolve — RDR 0003 (iteration 2)

Origin ledger: `diff.md` D-1..D-7 (seven findings from one alternate-model
reconstruction, `run-1.md` @ claude-sonnet-5). Every finding exits once.

## Dispositions

| ID | Disposition | Kind | RDR section touched |
| --- | --- | --- | --- |
| D-1 | **fixed** | pin + single-source | Technical Design → `Normative Contracts` (new participation clause); A2 Evidence |
| D-2 | **fixed** | pin | `Normative Contracts` (new report-every-defect clause); `disposition` table |
| D-3 | **fixed** (partial — remainder booked as A15) | pin | `Normative Contracts` (too-large bound clause); `disposition` table; Testing Strategy Scenario 3 |
| D-4 | **fixed** | pin | `Normative Contracts` (gap attribution); `disposition` table; Scenario 2 |
| D-5 | **fixed** | pin | `Normative Contracts` (new set-literal clause); A13 → `Verified`; Scenario 5 |
| D-6 | **fixed** | pin | `Normative Contracts` (report-every-defect clause, same edit as D-2); `disposition` table |
| D-7 | **dismissed-with-cite** | leave non-normative | none |

### Detail

- **D-1 (participating dimension).** Grounded against the draft: `:163` said
  dimensions "that vary inside that group" while `:442`, `:507`, and `:599` said
  "participating" unqualified — two readings, different scoped products, and a
  differently-sized A11 producer request. Pinned the *referenced* reading: a
  dimension participates when any row carries an atom over that key, whether or
  not it discriminates, because a non-varying key still bounds the product and
  can still lack a declared domain. Amendment sweep updated the `:163` site; no
  other "vary" phrasing survives.

- **D-2 + D-6 (finding-set scope and completeness).** One edit, since both are
  the same silence seen from two sides: the draft fixed the outcome *class*
  (blocking) but never its scope or cardinality, so `run-1` invented three
  short-circuit returns. Pinned: report every decidable defect in one pass; a
  withheld claim does not suppress overlap/coverage findings for the same group;
  one finding per unprovable dimension and per refusing row. The consequence the
  draft already forbids elsewhere is now closed here too — a two-refusing-row
  group must not emit a single finding chosen by iteration order, which would
  make output depend on source order (the property Scenario 5 tests).

- **D-3 (too-large bound).** Split. The unit and declaration shape are this
  RDR's to fix and are now pinned: cardinality of the scoped product, one
  published integer constant, reported beside the computed cardinality, not
  per-model/per-group/configurable. What this RDR *cannot* settle by drafting is
  whether a scalar cardinality actually predicts provability — proof cost may
  track dimension count or per-dimension width instead. That is a new
  load-bearing claim, so it is booked as **A15** (`Pending`, Method: MVV Test)
  with Scenario 3 extended to compare two equal-cardinality products of
  differing shape. Not closed by assertion (COMPUTE-DON'T-ARGUE: the comparison
  is not runnable at draft time).

- **D-4 (coverage-gap attribution).** Grounded: Scenario 2 requires gaps to
  "fail with source rule/context ids", but the diagnostics clause quantifies
  over *contributing predicates* and a gap has none — a real asymmetry, not a
  reading error. Pinned a distinct attribution for gaps: selection context,
  every rule id in the group, and one concrete uncovered assignment as witness.

- **D-5 (set-literal spelling).** Grounded against A13, whose own Plan said
  "state the set-literal encoding as a normative clause before lock — owned here
  and needs no peer". `run-1` independently produced the ordered, non-canonical
  slice A13's "If wrong" predicts, confirming the gap is reachable by an
  ordinary reader. Pinned as unordered, duplicate-free, canonicalized before
  entering the identity tuple, repeats rejected at parse. A13 → `Verified`
  (Design Decision, with the rejected alternative recorded). Phase 1's gating
  row and the Finalization Gate's A13 line updated.

- **D-7 (interim atom-naming representation).** Dismissed with cite: the draft
  already carries the A8 blockquote stating the atom-naming half is "a stated
  requirement with no producer", and A8 books the route-back with a venue.
  Prescribing an interim internal representation would specify implementation
  freedom the RDR deliberately leaves open (over-spec trap), and the finding
  itself rates it lower-weight. No edit.

## Inadmissible GUESS clusters (no action)

Package layout, all Go identifiers, `EvaluateAtom`'s signature, helper
decomposition, `SourceLocator` shape, and `Subset`'s opaque representation —
the RDR leaves these open by design and `run-1` correctly stopped at the
contract in each case. Recorded here only to show they were considered.

## Needs (re)verification (Stage 6)

- **A15** (new, `Pending`, Method: MVV Test) — a single published cardinality
  bound predicts provability. Plan: Scenario 3's equal-cardinality/differing-shape
  comparison; a divergence forces a second bound term and a restated clause.
- **A13** flipped `Pending` → `Verified` on a Design Decision made in this pass.
  Stage 6 should confirm the clause as written is what A13 claims, not re-derive it.
- No previously-`Verified` record was invalidated by these edits.

## Charted to successor

None — every finding landed inside this RDR's existing scope.

## Tiebreakers escalated

None.

## Escalation verdict (lens criteria)

**Not escalated to the full ×3 lens** — deliberately, against `diff.md`'s own
recommendation. The lens escalates when a lite diff surfaces a load-bearing or
cross-RDR silence *because one run under-powers the sample*. That rationale is
statistical, and it does not apply here: D-1 is a textual contradiction between
`:163` and `:507`/`:599` **inside the draft**, confirmed by reading the RDR, not
by counting reconstructions. D-3 and D-5 are likewise self-evident from the
draft — the RDR demands a bound it never declares, and A13 books an unstated
clause it already owns. Two further runs would poll model opinion on which
reading to prefer; they cannot make a contradiction more or less present. The
author's job is to pin the clause, which this pass did. `run-1` also showed no
overfit signature (34 GUESS markers, seven substantive silences), so the other
escalation criterion is unmet.

Escalating would also have required rewriting `run-1.md`'s `variant:` line to
`full (escalated: <reason>)` — deliberately not done; the header still reads
`variant: lite (profile: large)`, which is the correct durable record.
