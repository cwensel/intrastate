# Finalization Gate — RDR 0008: Ownership of the recognized-outcome tag key name

- **RDR**: `0008-recognized-tag-key-ownership.md`
- **Date**: 2026-08-21
- **Profile**: foundational
- **Lock**: re-lock following the 08.1 cluster demotion of 2026-08-11
  (supersedes the 2026-08-11 gate record)
- **Mechanical pre-sweep**: PASS
  (`evidence/tooling-pass/iter-2/tooling-pass.md`, C1–C6 no findings;
  C9 advisory, 1 field 1 line over budget)
- **Verdict**: **PASS — READY to lock**

## Contradiction Check

No contradictions between Research Findings and the Proposed Solution, and none
between planned features and stated principles. The 2026-08-11 gate cleared
three candidate conflicts (prior art vs. mechanism; "no kernel change" vs. the
settled enforcement locus; A4's finding vs. the Approach); all three passages
survive this pass unchanged and their resolutions still hold. This pass
re-checked only what the demotion and the two re-entry stages moved:

1. **A6/A11's re-derived conclusion vs. the Approach.** The demotion falsified
   the *evidence* for A6's conclusion, not the conclusion. Final RDR 0009 does
   place a whole-table precondition at `Resolve` entry, so the boundary A6 had
   called virgin is shared. The design is unchanged because the predicates read
   disjoint fields — 0009's is a method on `Table` and cannot read
   `Input.Owned`/`Input.Observed` at all, which is exactly what this RDR's
   check must read. What changed is the rationale, now weaker in the direction
   A6's own "If wrong" line anticipated: less novel surface, not more. Research
   Findings and Technical Design tell the same story.

2. **`Overrides` vs. A4's three-violation finding.** This was the cluster's
   `contradiction / blocks-impl` item, and it is closed. The body always
   disclosed the collision (A4: three committed declarations violate the rule);
   only the Metadata summary denied it with "narrows nothing in either peer."
   The field now states the narrowing and cites JDR 0001 §JD-10's ruling
   ("rename them; there are no users to migrate"). Summary and body agree.

3. **Failure Modes' reclassification claim vs. A10's actual reach.** The
   reconcile refuted "no user data is reclassified as a programmer mistake" as
   written — A10 sweeps only the accessor half, while RDR 0005's planned `flow`
   verbs feed `Input.Observed` from `--tag name=value`, making
   `--tag recognized=x` a reachable user invocation. The claim is now scoped to
   what A10 proves, A10's "Consequence for A8" is corrected from "closes it" to
   "closes the accessor half," and the open remainder is latched at JD-9. The
   over-reach is gone rather than papered over.

4. **Block 4's interim ordering rule vs. JD-5's openness.** Block 4 licenses
   "either error is conforming; what is *not* conforming is skipping this check
   because another fired." JD-5 says "either order is defensible." These agree;
   block 4 is strictly the weaker claim and forecloses nothing JD-5 might pick.

## Assumption Verification

Thirteen Evidence Records, every one internally consistent — Status, Method, and
Evidence agree, and no "If wrong" line is empty.

- **Status terminality**: zero `Pending`, zero `Unverified`. Twelve `Verified`
  (four carrying a qualifier that narrows a *consequence*, not the assumption);
  A12 `Refuted as stated`, a terminal disposition whose repair is written as a
  hard Prerequisite gate item rather than left open.
- **The re-verify set {A6, A11}**: both re-derived against the *Final* text of
  RDR 0009 and against source, not against the retired Draft. A6's six quoted
  claims were located verbatim in Final 0009; `::Resolve`, `::Input`, `::Row`,
  `::recognizedTagKey` all confirm the shape at HEAD, including that all five
  `return` sites pair `Result` with a literal `nil`. A11's disjointness now
  rests on the two field sets, not on the withdrawn `Input`-count negative.
  The false negative-existential survives only inside `Corrects:` bullets and
  the delete-on-re-lock defect record, correctly labeled withdrawn.
- **Method vocabulary**: all thirteen in the sanctioned eight (compound labels
  are `+`-joined members, not paraphrases). No load-bearing `Docs Only` exists —
  the label appears zero times.
- **Source Search self-reference**: none. No Evidence path resolves to this RDR
  or its artifact directory.
- **Anchor durability**: the mechanism that produced this demotion — bare
  peer-RDR line citations — is repaired at the root. All 15 live ones are now
  section-heading or assumption-ID anchors (0002 ×6, 0004 ×4, 0006 ×3,
  0007 ×2). Source anchors are `path::Symbol` and every one resolves at HEAD.

The postmortem (`0008-recognized-tag-key-ownership-postmortem.md`) records both
defects this cycle surfaced and their shared root cause — an assumption verified
against a *moving* peer or a *future* caller, stamped with a point-in-time fact
and no re-verification trigger. That ledger is opened, not owed.

## Scope Verification

The Minimum Viable Validation is in scope and split in two, matching the Done
split in Testing Strategy. The split is a genuine ownership boundary, not a
deferral of convenience:

- **Kernel half — runs during this RDR's implementation.** A row matching on the
  tag `recognized` fires against a recognized outcome, asserted through the real
  kernel's assembled view (behavioral conformance, per premortem P-5); and an
  `Input` supplying an owned or observed tag keyed `recognized`, or a row naming
  it in `RequiresOwned`, is rejected by the exported predicate and by `Resolve`
  at entry. Both run against `internal/resolve` as it stands.
- **Normalizer half — runs inside RDR 0002's implementation**, because the
  normalizer is 0002's code. A table declaring a recognized-provenance tag named
  `recognized` loads and lints clean; renamed (and separately, with an owned tag
  so keyed) it fails load/lint in the `reserved_tag_key` category carrying a
  direction-appropriate remedy name and rule identifier — not a silent no-match
  at resolve time.

The RDR states plainly that the kernel half alone does not deliver the Problem
Statement's outcome, and gates the normalizer half behind a **written hard
Prerequisite** rather than behind an intention. The gate blocks claiming the
Problem Statement outcome, so the split cannot ship silently half-done. Nothing
MVV-critical is deferred: no JD this RDR latches (JD-4/5/8/9/10) is load-bearing
for what the MVV proves — each bears on peer surfaces, report channels, or lint
semantics downstream of the naming rule the MVV exercises.

Determinacy: `foundational` requires the full repeatability variant. Present —
three runs with three **distinct** model stamps (`claude-opus-5[1m]`,
`claude-fable-5`, `glm-5.2:cloud`) plus `diff.md` and `resolve.md`. No variant
mismatch, so no `determinacy: n/a` disposition is owed.

## Cross-Cutting Concerns

- **Incremental adoption.** The constraint invalidates all three of RDR 0002's
  canonical fixtures — this RDR's largest cross-cutting cost, and it is
  disclosed in `Overrides` rather than buried. The policy is not this RDR's to
  set: JDR 0001 §JD-10 owns it and has ruled ("rename them; there are no users
  to migrate"). Phase 2 carries the three renames plus six reference sites.
- **Canonical-form / determinism.** Applies only in the weak sense: this RDR
  fixes an *identifier* (`recognized`) in the assembled view. It claims no
  byte-identical output, no content-addressed identity, and no replay-stable
  hash, so the hash/pre-image/iteration-order sub-checklist does not engage. The
  one determinism property it does assert — that a reserved-key breach is
  detected whenever present, never skipped because another precondition fired —
  is stated normatively in block 4 and tested by scenario 6.
- **Versioning of the validation surface.** The `reserved_tag_key` category is
  additive within RDR 0002's "including at minimum" extensible list, and RDR
  0005's A-block pre-authorizes resolver-specific `Code` values as needing "not
  new envelope fields or exit groups." Where the code lands in the CLI envelope
  is JD-8's, and the RDR disclaims it to RDR 0005 rather than deciding it.

Omitted as inapplicable: licensing, deployment model, IDE compatibility,
secret/credential lifecycle, memory management, concurrency model, character
encoding, build-tool compatibility.

## Proportionality

Right-sized; nothing to trim before locking.

**Contract count — the split test.** Six ```normative``` fences, but they are
clauses of a **single** independent contract: the reserved-key identity rule and
its enforcement. No second independent seam is present — no distinct hash, wire
format, taxonomy, or destructive-op policy — so the ≥2 split signal does not
fire and there is nothing to split out.

**Profile re-validation.** `foundational` confirmed against the contracts just
counted, and it comes from the cross-RDR-producer trigger rather than the
contract axis. The re-verification *strengthens* it: 0008 binds RDR 0001's
kernel carrier to RDR 0002's declared model, extends 0002's validation surface,
and now additionally co-resides with RDR 0009 at `Resolve` entry. That is more
cross-RDR coupling than at the original lock, never less. The accretion floor
does not apply — `Seam Lineage` records zero prior closed point-fixes, below the
≥2 threshold — and could only have raised a profile already at the ceiling. All
four required lenses (`cove`, `3amigo`, `critique`, `repeatability`) have
evidence dirs; the documented single-model `critique` fallback is recorded in
`resources.md` and was compensated by repeatability's genuine multi-model draw.

**Length.** 2,567 lines is at the top of the corpus, and the mass is
load-bearing: thirteen Evidence Records carrying the verification anchors the
grounding sweep reads, plus the JD latches this RDR must cite rather than
restate. C9 flags exactly one Evidence field one line over an advisory budget.
Cutting there would blind the sweep for no gain.

## Note on the JD-5 lock question

The Stage 4 re-entry's author round recorded "cite JD-5, stay Draft" on the
reasoning that block 4's interim rule depends on JD-5, and that a Status latch
"would compose awkwardly with the `re-verify` qualifier Stage 8 expects to
self-clear." Both halves are resolved, so this is not a blocker at the Gate:

- **The corpus convention is settled and unanimous.** An open joint decision is
  carried as a Status latch on a **Final** RDR — RDR 0005 `Final [joint decision
  → JDR 0001 §JD-8, §JD-9]`, RDR 0006 `Final [… §JD-4]`, and RDR 0009 `Final […
  §JD-5]`, which is this RDR's own peer on this very decision. Holding 0008 at
  Draft for JD-5 while 0009 is Final for JD-5 would make the two sides of one
  joint decision disagree about what an open JD means.
- **The Stage 6 reconcile answered the composition worry.** The re-lock
  overwrites the whole Status value, so the latch *replaces* the `re-verify`
  qualifier cleanly — there is nothing to compose.
- **The gate is on implementing, not on locking.** The cluster verdict is NOT
  RECONCILED, "Do not implement over these," and the RDR's own
  implementation-ordering passage already holds itself at Final-unimplemented
  until 0002 starts and the depended-on JDs close. The reconcile states it
  directly: "Locking remains correct; the gate is on implementing."
- **Block 4 is implementable as written today.** Its interim rule is the weaker
  claim JD-5 may later narrow, and the RDR says so: "nothing here forecloses
  that, and no implementer needs to invent one to proceed." Scenario 6's fourth
  variant asserts only the non-silent-skip property and is marked as the test
  that pins JD-5 once it closes.

Locked as `Final [joint decision → JDR 0001 §JD-5]`.
