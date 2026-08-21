Model: claude-opus-5[1m]

# 3amigo Dispositions — iter-2 (re-entry, delta-scoped to A8)

Origin ledger = `consolidation.md` (T-1 … T-14). Every entry exits exactly one
way. Grounding gate run against code on `main`, `{RDR_RESOURCES}`, and the RDR's
own decided text before any edit.

## Grounding gate — what was confirmed before editing

- `internal/resolve/resolve.go::Tag` is `{Key, Value string}` — no absence
  channel. CONFIRMED.
- `internal/resolve/resolve.go::assemble` writes every `in.Owned` tag into the
  view keyed by `t.Key` with `ProvenanceOwned`. CONFIRMED.
- `internal/resolve/resolve.go::TagSet.has` returns `ok && provenance match` —
  map presence. CONFIRMED.
- `internal/resolve/resolve.go::missingOwned` accumulates a key only when `has`
  is false, so a sentinel-valued tag retires `owned_state_unavailable` for that
  key. CONFIRMED — this is the mechanism behind T-1.
- `internal/resolve/resolve.go::Refusal.MissingOwned` is documented as "owned tag
  keys **absent from the snapshot**". CONFIRMED — the kernel's defined channel
  for absence is omission, which is what collapsed T-1's apparent fork.
- Spike `main.go::read` falls back to `main.go::expectedTagKeys(art)` when no
  keys are passed; transcript line 1 and both replay runs use that path.
  CONFIRMED — the circularity behind T-2 and T-13.
- Spike `main.go::write` clones `art.tags` directly and never calls `read`.
  CONFIRMED — the gap behind T-4.
- `main.go::newPartialArtifacts` hand-populates `unreadable`. CONFIRMED — the
  gap behind T-5.
- JDR 0001 §D3 resolves (b) and parks (c) as "available if implementation shows
  the kernel needs to distinguish read-failed from absent". CONFIRMED. §D3's
  fork is read-failed vs absent; the personas' *absent vs present* case is a
  third one §D3 did not settle, so no re-raise applies.

## Dispositions

- **fixed** — T-1 (hotspot, 3 personas: absence representation unpinned) —
  sections touched: Normative Contracts (new seam clause), Load-Bearing
  Decisions (new "Absence crosses the seam as omission"), Fidelity Table `read`
  row, Disposition Table absent-key row + closing sentence, Desk Trace step 2
  and the no-CONTRADICTION paragraph, Failure Modes, Risks, Proportionality,
  Approach. The fork (sentinel / omission / typed marker) was collapsed by
  evidence, not by tiebreaker: RDR 0001 already defines `Refusal.MissingOwned`
  as keys *absent from the snapshot*, so omission is the encoding the locked
  kernel contract expects and a sentinel is the one that breaks it. The RDR now
  pins what crosses the seam while leaving the accessor layer's internal
  representation free.
- **fixed** — T-2 (hotspot, 2: requested key set has no origin) — sections
  touched: Normative Contracts (new validated-key-set clause), Load-Bearing
  Decisions (new "Requested key set"), Technical Design phase one and two,
  Validation Scenario 1 (now eight arms), Disposition Table validation row.
  Derivation from what was read is now explicitly non-conformant.
- **fixed** — T-3 (hotspot, 2: `incomplete_read` never minted normatively) —
  sections touched: Normative Contracts (new clause naming the class, requiring
  distinctness from execution failure and timeout, and requiring the refusal to
  name the unreadable keys). Parallels the existing timeout clause.
- **fixed** — T-4 (hotspot, 2: read completeness at the read-back re-read) —
  sections touched: Normative Contracts (new `read_back_incomplete` clause),
  Desk Trace step 5, Disposition Table (new row), Validation Scenario 3, Failure
  Modes refusal list. Resolved as its own class rather than folding into
  `read_back_mismatch`: "I could not verify" and "the artifact is wrong" are
  different operator facts, and the existing table already separates
  verification failures from mismatches.
- **fixed** — T-5 (hotspot, 2: no observable criterion for "could not read") —
  sections touched: Load-Bearing Decisions (new "Absent vs unreadable is the
  binding's call, defaulting to unreadable"). Pinned the property rather than
  per-binding mechanics: absence requires a successful read that found no key;
  every other outcome is unreadable, and a binding that cannot tell them apart
  must report `incomplete_read`. Guessing absence is the failure the contract
  exists to prevent.
- **fixed** — T-6 (hotspot, 2: MVV asserts the branch and stops) — sections
  touched: Minimum Viable Validation, Validation Scenario 2, Oracle
  Discriminability row 2 (second negative control for the derived-key-set bug),
  new Scenario 6 carrying an absent required key through to the resolver.
- **fixed** — T-7 (PM: silent-failure inventory names one shape where the delta
  created two) — sections touched: Failure Modes (now three shapes, the third
  being absence-read-as-present with the seam clause as its guard), Risks (new
  matching risk/mitigation pair).
- **fixed** — T-8 (QA: no timeout / incomplete-read precedence) — sections
  touched: Normative Contracts (precedence sentence), Technical Design phase
  two, Disposition Table (new row), new Validation Scenario 7, Oracle row 7.
- **fixed** — T-9 (Implementer: whole-read refusal asserted, not argued) —
  sections touched: Normative Contracts (explicit whole-read sentence),
  Load-Bearing Decisions (new "Read refusal granularity" recording it as a
  decision aligned with JDR 0001's refuse-don't-guess rule).
- **fixed** — T-10 (QA: totality is a lower bound; over-return unforbidden) —
  sections touched: Normative Contracts ("exactly the keys it was asked for — no
  requested key missing, no unrequested key added"), Fidelity Table `read` row
  (now key-set equality). Grounded: an unrequested key entering `assemble` with
  owned provenance would shadow caller-supplied observed context.
- **fixed** — T-11 (PM: Proportionality covers the branch rule, reads as
  covering representation) — sections touched: Proportionality. Records that
  naming the seam encoding constrains this RDR's own output and leaves RDR
  0001's `Input` shape unchanged, so the scoping claim now covers what it is
  doing work for.
- **fixed** — T-12 (PM: no-CONTRADICTION argued against the wrong pair) —
  sections touched: Desk Trace closing paragraph, which now examines the
  cross-layer pair as well and shows it is non-colliding *because* of the seam
  clause.
- **fixed** — T-13 (QA: witnesses come from the explicit-key path) — sections
  touched: Load-Bearing Decisions "Requested key set" names
  `main.go::expectedTagKeys` as fixture convenience rather than contract; A9's
  Evidence records that the spike's default path, sentinel absence, non-
  overlapping timeout, and cloned read-back are exactly what it does not
  witness. The spike is not edited — it is Resolve evidence, and the RDR now
  states where it stops rather than overclaiming it.
- **dismissed-with-cite** — T-14 (Implementer: Scenario 2 bundles five
  dispositions under one number) — the Oracle Discriminability table keys off
  MVV scenario numbers, and the Implementer's own finding names that as the
  reason a silent split is harmful. Renumbering an existing scenario would break
  every Oracle row and Witness cite that references it. Scenario 2's assertions
  were strengthened in place and the genuinely new cases were added as scenarios
  6 and 7 instead, which preserves the mapping. Test decomposition below the
  scenario number is an implementation choice this RDR does not constrain.

## Net-new scope

None charted. Every fix binds to the read-completeness contract this re-entry
reopened — the branch rule's success predicate and what that branch hands across
the accessor→resolver seam. No fix widens the RDR into RDR 0001's `Input` shape,
RDR 0002's carrier, or RDR 0005's CLI mapping.

## Needs verification (Stage 6 closes these — not verified here)

- **A9 added, Status: Pending, Method: MVV Test.** Books the four load-bearing
  claims these fixes introduced that the Resolve spike does not witness:
  (1) absence crosses to the resolver as omission from the owned snapshot;
  (2) a missing/empty requested key set is rejectable before execution (the
  eighth validation arm); (3) `timeout` outranks `incomplete_read` on a partial
  read; (4) a read-back re-read that cannot read a compared key is reportable as
  `read_back_incomplete`. Verified by MVV Scenarios 1, 2, 3, 6, and 7.
- **Prerequisites checkbox flipped** from `[x] All Critical Assumptions verified
  (A1-A8)` to unchecked, naming A9 Pending. The re-entry can no longer claim a
  fully-verified assumption set at lock.
- **A8 is untouched and stays Verified.** Its three read dispositions still hold;
  these fixes constrain what the success branch emits downstream, which A8 never
  claimed either way.

## §amendment-sweep

Pre-edit tokens grepped and every surviving site updated in the same pass:
`total over the requested keys` (Fidelity, Load-Bearing) → key-set equality;
`is not pinned` (Fidelity) → scoped to internal representation only;
`absent value` / `<absent>` (8 sites: A8 Evidence, Load-Bearing, Disposition
Table, Oracle row 2, Fidelity, Desk Trace step 2, Validation Scenario 2 ×2) →
re-read against the seam clause, each either amended or confirmed to be
describing the spike fixture rather than the contract; `(7 shapes)` → `(8
shapes)`; `two silent-failure shapes` → three; `typed absence marker` removed
(it was the escape hatch T-1 turns on). The refusal-class list in Failure Modes
gained `read_back_incomplete` so the enumeration matches the new clause.

Producer/consumer re-read: the new seam clause was checked against RDR 0001's
`Input`/`Refusal` shape (unchanged by it), against A8's own Evidence sentence
(which cites `missingOwned` and now agrees with the clause instead of cutting
against it), and against the Disposition Table's no-silent-exit guarantee (which
needed the amendment recorded under T-1).

Cross-lens check: the iter-2 grounding pass edited A8's Evidence citation form
and MVV Scenario 5's capture assertion. Neither field is touched by these fixes —
no two-lenses-pinned-one-field conflict.

## Cross-RDR note (not a finding — recorded for reconcile)

RDR 0007 carries **A6b as Pending**, downgraded at its Stage 6 with the
obligation ROUTED to RDR 0004: "must state whether the 'typed tag values' branch
carries a completeness guarantee". This RDR now discharges it — the completeness
clause plus the seam-omission clause together answer both halves — and RDR 0007
is named in References. Closing A6b is RDR 0007's own edit, not this pass's;
flagged so reconcile can route it.

## Tiebreakers

None escalated. T-1 looked like the one genuine either/or (sentinel vs omission
vs typed marker) but collapsed against `Refusal.MissingOwned`'s documented
meaning — the kernel already defines absence as omission from the snapshot, so
the evidence chose, not the author. No §strong-consult was needed.

## Mini-check cue read

Not owed by this pass. The iter-2 grounding pass ran the cue read and its four
fired tables (Disposition Table, Oracle Discriminability, Fidelity Table, Desk
Trace) are in the draft; tables persist across passes. These fixes added no new
cue: the source-authority census remains `no` — the executor is still the sole
writer of the read result, and the seam clause names one encoding rather than
introducing a competing source of truth.
