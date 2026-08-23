Model: claude-opus-5[1m]

# Grounding Dispositions — iteration 2

Origin ledger = `findings.md` F1–F4 (this sweep's own findings; grounding's
first-run findings *are* the originating concerns).

- **F1 view-level conformance misattributed to RDR 0007** — **fixed**;
  origin: F1; section touched: `Technical Design` invariant 5, plus new
  Critical Assumption **A7**. Invariant 5 no longer assigns the view-level check
  to RDR 0007. It now states that the model-level check discharges RDR 0003's
  A18 premise on its own (A18 names two sufficient producers and requires one),
  and records the view-level check as unowned, citing both that RDR 0007 states
  no such obligation and that `internal/resolve/resolve.go::assemble` performs
  none. A7 books the residual: it flips to `Verified` when RDR 0003 A18 records
  invariant 5 as its producer.

- **F2 `single-step` cited as RDR 0001 text** — **fixed**; origin: F2; section
  touched: A6 Evidence. Reworded to what RDR 0001 normatively says (stateless;
  one legal transition plan per supplied snapshot) and marks single-step
  explicitly as this RDR's gloss. The downstream conclusion — multi-step
  reachability exists nowhere else in the cluster — is unchanged and still
  holds. The `Capability Dependencies` cell repeating it is a summary pointing
  at A6, which now carries the qualification; no second edit owed.

- **F3 two partitions of the five declaration fields** — **fixed**; origin: F3;
  section touched: `Technical Design` input contract. Aligned to RDR 0003's own
  partition (value kind; finite domain as enum set or `{min..max}`; optionality
  marker; single-valued marker; set-element universe as its own field), matching
  A2's list. The normative exhaustiveness clause was deliberately **not**
  widened: it is scoped to what lint reads for proofs (finite domain,
  optionality, single-valuedness), not a restatement of the model.

- **F4 four peer records discharged but unnamed** — **fixed**; origin: F4;
  section touched: A2 Evidence. Each done-condition was checked against the
  current draft before claiming closure: A17 ← invariant 3's two-population
  overlap; A18 ← invariant 5 as model-level producer (with A7); A19 ←
  `graph-coverage-closed-by-escape`; A20's lint half ← the atom field on the
  finding contract. Named so RDR 0003's ledger can close by citation. This
  records existing content, not new scope — no draft behavior changed.

## Dismissed with cite

- **Escape-row rescue alphabet not enforced by the `Row` type** —
  **dismissed-with-cite**; the closed alphabet is fixed upstream by
  `0002::Normative Contracts` ("An `escape` list MUST contain only ... `no_match`
  and `ambiguous_match`"), and this RDR's invariant 1 requires every named
  element to resolve to a declared model element. The per-class partition is
  correctly general as written. Recorded in `findings.md` for implementers.
- **`{"type":"failed"}` envelope** — **dismissed-with-cite**; the stale claim is
  in `respond.go`'s package comment, not in this RDR. The draft says findings
  ride "both the error envelope and the success payload" and never asserts the
  discriminator, so no edit is grounded.
- **Aggregation veto's literal subject** — **dismissed-with-cite**; "a row can
  refuse" is RDR 0003's defined row-level predicate (a syntactic test over the
  declared optionality field), which this RDR cites by name. Rewriting to RDR
  0007's resolution-level phrasing would break the citation, not fix it.
- **§JD-14 reading differently from invariant 3 on its face** —
  **dismissed-with-cite**; RDR 0003 A17 records that §JD-14's overlap half
  "rested on a premise the shipped kernel refutes" and that this RDR's reading
  "was closer to right than the gate credited". RDR 0003's current clause checks
  two populations, which invariant 3 matches. The stale document is the JDR; its
  correction is already scheduled at the next JDR 0001 touch (A17's plan).
- **CI push trigger / `CLIError.Cause` / default-on provenance** —
  **dismissed-with-cite**; each is a recorded imprecision that does not change a
  claim the RDR makes. Detail in `findings.md`.

## Needs verification

- **A7** (new, `Status: Pending`, Method: Peer RDR) — flips to `Verified` when
  RDR 0003 A18 records invariant 5 as its model-level producer. Not
  lock-blocking by A18's own Stage 6 disposition (DOWNGRADED, survivable).

Pre-existing and unchanged by this pass: **A5** (Pending, MVV Test — production
gate wiring via Validation scenario 7) and **A6** (Pending, Peer RDR — RDR 0002
declaring initial owned state and terminal states).

## Charted to successor

- None. F4 recorded existing content rather than absorbing new scope.

## Tiebreakers

- None. Every fork collapsed on evidence: F1 resolved by A18 naming two
  sufficient producers and requiring one; the §JD-14 apparent contradiction
  resolved by A17's recorded correction of the gate's premise.

## Mini-checks (first lens pass on this RDR owes the cue read)

Cues read from Normative Contracts, the MVV, and Testing Strategy — not a file
grep. Two fired; both tables written into the RDR.

- **disposition — FIRED.** Cue: the draft sets exit/op outcomes across input
  classes (nine finding codes split blocking/info, one aggregate `CLIError`, a
  `GroupUserEnv` exit, a success payload), and silent-vs-loud is load-bearing
  here — invariant 3 deliberately reports *nothing* for escape/ordinary overlap
  while RDR 0003 forbids a bare escape row being a silent opt-out. Table added
  after the finding-codes table: input class × exit · envelope · finding minted
  · silent-or-loud. It makes the one intentional silence explicit and cites why.
- **desk trace — FIRED.** Cue: several normative clauses bear on one output
  surface (the findings envelope) and on the MVV's end state. Table added at the
  end of Testing Strategy, walking one invocation in eight steps.
  **No CONTRADICTION row.** Step 7 is the only joint-satisfiability point — a
  group both closable by a bare escape row and carrying a row that can refuse
  `guard_unevaluable`. The narrowing clause settles it (withholding dominates:
  blocking `graph-unprovable-coverage`, not a success with the info code),
  corroborated by RDR 0003: "a bare escape row MUST NOT be read as discharging
  the narrowing ... at runtime that refusal is returned before the escape row is
  ever consulted." Recorded so the info code is not read as an escape hatch.
- **source-authority census — did not fire.** The RDR concerns authority, but
  every arm is already enumerated in `Capability Dependencies`, the
  `Existing Infrastructure Audit`, and the command-authority clause. No fallback
  path, derived/propagated output, or ambiguous source-of-truth pair.
- **test-discriminability — did not fire.** The MVV asserts exact stable codes
  per invariant class and carries explicit negative controls (scenario 11's
  always-present control, 9(a)'s "no `graph-overlap`", 8's "no
  `graph-coverage-gap`"). No absence-of-error, exit-0, or fixture-name oracle.
- **round-trip / fidelity — did not fire.** The draft declares the class absent:
  "This RDR introduces no encode/decode, import/export, or inverse operation."
