Model: claude-opus-5[1m]

# Persona 1 — Product Manager (iter-3 re-entry, delta-scoped)

Question: does this RDR actually deliver the user outcome?
Scope: the re-entry delta only — the tag declaration model rehomed from RDR 0002
(2026-08-21), the 0002/0003 split, A8 `Verified` via §JD-4, §JD-13's
single-valued marker, §JD-14, A5 re-verification, and the open A10/A12/A14/A15.

The user outcome this RDR names in its own Problem Statement: **a flow author
needs conditional edges the lint can prove exhaustive and mutually exclusive.**
Findings are judged against that outcome, not against internal document hygiene.

---

## P1-1 — MEDIUM — `Capability Dependencies` row "**Tag declaration model** (value kind, finite domain, optionality, set-element universe)" (line 1055), `Scope Verification` (lines 1740-1742), and `Load-Bearing Decisions` → *Operator semantics* (lines 1012-1016)

**Concern.** §JD-13 homed a *fifth* field — the single-valued marker — into this
RDR's declaration model, and the normative clause (lines 714-725) and the
`authority` table (lines 937, 940) carry it. But three enumerations of the same
model inside this document still list only **four** fields:

- `Capability Dependencies`, the row that tells an implementer what this RDR
  introduces: "(value kind, finite domain, optionality, set-element universe)".
- `Scope Verification`, the passage that states what scope grew to: "(value
  kind, finite domain, optionality, element universe)".
- `Load-Bearing Decisions` → *Operator semantics*, which enumerates what a
  declaration supplies per operator and never reaches the marker.

Metadata (line 17), the Approach (line 543), and the normative clause all say
five. The document disagrees with itself about the size of the contract it owns.

From the product angle this is not a typo: `Capability Dependencies` and
`Scope Verification` are the two passages a planner reads to size the work and
to decide what "done" covers. An implementer scoping Phase 1 from the
`Capability Dependencies` row builds a four-field declaration model, and RDR
0006's **mandatory blocking** `graph-single-valued-state` — which JDR 0001
§JD-13 confirms is gated on this field (`0006:315`, MVV at `0006:640`) — has no
producer at runtime. The exact unhomed-producer defect §JD-13 was opened to
close reappears one layer down, in this document's own scoping prose.

**Decision blocked.** Phase 1's scope definition — whether the declaration model
an implementer builds has four fields or five — and therefore whether RDR 0006's
mandatory blocking single-valued invariant is deliverable at integration.

---

## P1-2 — MEDIUM — `Normative Contracts` single-valued clause (lines 714-725), against RDR 0002's authoring surface

**Concern.** The split this re-entry establishes is: RDR 0003 owns *meaning*,
RDR 0002 owns *authoring location and normalization carriage*. For the four
original fields that split closes cleanly — RDR 0002's normative type-model
clause (`0002:342-350`) names "value kind, and optionally a finite domain, an
optionality marker, and a set-element universe" and requires normalization to
carry "every declared field through to the normalized model without loss", and
its schema list (`0002:217-220`) repeats the same four.

For the single-valued marker it does not. `single-valued` occurs **zero** times
in RDR 0002 (verified by count). So the field this RDR now declares the meaning
of has **no stated authoring location** and is not named in the carriage clause's
enumeration. The generic "every declared field ... without loss" arguably covers
it, but the enumerated four are what an author reads and what a normalizer
implementer codes against.

This RDR's MVV `Authorability` paragraph (lines 1454-1457) asserts the opposite
as settled — "a single-valued enum tag declares the marker §JD-13 homed here" —
and the whole point of that paragraph is that the MVV is authorable against
owned declarations rather than peer requests. For four fields that is true. For
the fifth, authorability still rests on a peer edit that no Direction list names:
RDR 0002's duty list (`0002:895-908`) carries A14's `block` retention and
§JD-10(b), and its ownership-correction paragraph enumerates the same four
fields — it does not carry a single-valued authoring slot.

**Decision blocked.** Whether the MVV's single-valued case is authorable today
(as `Authorability` claims) or is a fifth peer request that belongs beside A14 —
which decides whether Phase 3's fixture can be written before RDR 0002's refine.

---

## P1-3 — LOW-MEDIUM — `Normative Contracts` declaration clauses (lines 684-735) and `Illustrative Code` (lines 1038-1046)

**Concern.** The re-entry moved a five-field type system into this document, and
the document states every field in MUST-language — but never once shows an
author what a tag declaration *looks like*. `Illustrative Code` is unchanged from
before the rehoming: it shows only `[rule.guard.all]` / `[rule.guard.unless]`,
the surface this RDR does **not** own the authoring of. The one surface this
re-entry newly made this document's own is the one with no illustration.

The user outcome depends on this. The RDR's own `Risks and Mitigations` names
"Finite-domain declarations feel like boilerplate" as a live risk, and its
mitigation is "keep declarations close to tag definitions in RDR 0002's model" —
a mitigation that now points at a document that, per the split, defers meaning
back here. Neither document shows the shape. The proportionality between "author
writes a conditional edge" (the outcome) and "author first declares kind, domain,
optionality, single-valuedness, and element universe for every participating
tag" (the cost) is nowhere visible in a form a reviewer can judge.

Note this is a *legible-cost* gap, not a correctness gap: the clauses are
individually unambiguous, and RDR 0002 owns the placement (`[tags.<tag>]`), so an
illustration here would be the meaning-side sketch, not a competing schema.

**Decision blocked.** Whether the declaration burden is acceptable to flow
authors — i.e. whether the boilerplate risk is mitigated or merely named — which
cannot be judged without seeing the authored shape.

---

## Passages reviewed and found sound (no finding)

Recorded so the delta scope is visibly covered rather than silently skipped.

- **A8 `Verified` / `authority` row "Exhaustiveness verdict for a row group" /
  `Contradiction Check`.** Grounded: JDR 0001 §JD-4 is CLOSED 2026-08-22 and
  reads exactly as this RDR reports — this RDR is the recording document, RDR
  0006 cites and mints no code, no non-blocking tier, `graph-unprovable-coverage`
  reused. RDR 0006's own Status line (`0006:10-13`) records the same assignment.
  The consequent atom-level-field duty is correctly booked on RDR 0006 and
  correctly stated as non-gating here. The user-visible promise — a green
  exhaustiveness result means resolution succeeds — is stated once and cited
  once, which is the outcome §JD-4 exists to protect.
- **§JD-14 / the escape-row clause (lines 799-807).** §JD-14 decides this RDR's
  reading governs and that the repair is RDR 0006's; this RDR correctly makes no
  edit and claims none. No PM concern: the author-visible behavior (an escape row
  that overlaps a guarded row is a blocking finding, not an exemption) is the
  one that matches RDR 0001's runtime refusal, so lint and runtime tell the
  author the same story.
- **A5 re-verification.** All four anchors resolve; `RuleID`/`SourceLocator`
  ship at `internal/resolve/resolve.go:173-174`. The outcome that depends on it —
  a diagnostic that points the author at the specific guard atom to fix — has a
  shipped producer for the row half and a booked producer for the atom half.
- **A11 / the 0002-0003 split.** The rehoming rationale is verifiable from both
  sides: RDR 0002's ownership-correction paragraph (`0002:909-919`) states the
  same split independently, and RDR 0007 (`Final`) assigns value kinds and set
  universes here. The declaration model genuinely was homeless. Scope growth is
  justified for the outcome: a grammar that cannot say what a tag may hold cannot
  prove anything about it.
- **A10 / A12 open tolerances.** Both are homed at RDR 0006's refine, and RDR
  0006 is `Draft` with both items on its Status line — a scheduled producer, not
  an open venue. Neither blocks the author outcome at this document's lock.
- **A14 open.** Correctly scoped as carriage rather than semantics, correctly
  left with RDR 0002, and present on RDR 0002's duty list at `0002:897-902`.
- **A15 open.** Discharged by MVV Scenario 3, which is written to make the
  comparison the record needs. Bounded "If wrong".
