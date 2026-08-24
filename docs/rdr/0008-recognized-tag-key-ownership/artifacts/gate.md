Model: claude-opus-5[1m]

# Finalization Gate — RDR 0008 `recognized-tag-key-ownership`

- **Date**: 2026-08-23 (re-lock; supersedes the 2026-08-21 record)
- **Verdict**: **READY — Gate PASS**
- **Mechanical sweep**: `evidence/tooling-pass/iter-3/tooling-pass.md` — PASS,
  no findings C1–C6; C9 advisory only.
- **Re-entry**: cluster `0002-0009` iteration-2 demotion, STAGE-SCOPED to
  {A9, A4}; closed by Stage 4 re-verify + Stage 6 iteration-3 reconcile
  (`evidence/reconcile/iter-3/report.md`, RECONCILED).

## 1. Contradiction Check

No conflict between Research Findings and Proposed Solution, and none between
planned features and stated principles.

The one contradiction this re-entry existed to close is closed in live text, not
merely dispositioned. Both halves were re-read at this gate:

- **Block 2's predicate-position clause** no longer claims predicate positions
  "need no separate reserved-key check … those references resolve to it." It now
  defers every predicate position to Final 0002's outcome-binding contract:
  `[rule.match.recognized]` is the mandatory outcome binding *lifted into*
  `Row.Outcome` and never left in the predicate set, and a `recognized` atom
  under `guard.all`/`guard.unless` is *refused at load* as a malformed outcome
  binding. That is 0002's fenced rule restated, not contradicted.
- **Scenario 5** no longer expects "two failures, not one." It expects **exactly
  one** refusal in the `reserved_tag_key` category, explicitly marks which of the
  two declarations gets reported as unspecified and non-assertable, and cites
  0002's fail-fast fence as the authority. The cardinality consequence is
  preserved but restated as "fix one, load again" rather than two reports from
  one load.

Block 2's reservation/binding split (reserved unconditionally; binding scoped to
resolves carrying an outcome) is consistent with `assemble` injecting only for a
non-empty `Input.Recognized` — the RDR states the empty-outcome path as RDR
0001's residual rather than asserting a breach it does not own.

The Direction sweep items are landed and re-checked: "0009 still Draft" corrected
(0009 is Final), stale 0007 quotes re-anchored, "§JD-4 open there" corrected to
answered-by-§D4.

One stale fact is knowingly carried and is **not** a contradiction in this RDR:
Final 0002's dependency table records RDR 0008 as `Final` while 0008 stood
`Draft [revised from Final]`. It is unamendable from this RDR's stage and
self-corrects at this lock.

## 2. Assumption Verification

Thirteen Critical Assumptions; every record is internally consistent —
Status/Method/Evidence agree and every "If wrong" branch is non-empty.

- **Statuses are all terminal.** Twelve `Verified` (four carrying a scoping
  qualifier that narrows a *consequence*, not the assumption); A12 carries
  `Refuted as stated`. Zero `Pending`/`Unverified`.
- **A12's refutation is terminal, not open.** The refuted claim was a
  *discoverability mechanism* (that a 0002 implementer would traverse to this
  RDR's constraint), never a design premise — no normative contract moves. Its
  "If wrong" branch is adopted as the repair and written as a hard Prerequisite;
  the cross-RDR-edit form has since happened (Final 0002 names
  `reserved_tag_key (RDR 0008)` in its validation-category block, its ownership
  table, and its dependency table — verified at this gate, 18 occurrences of
  `0008` in Final 0002). The unabsorbed residue (block 3's payload fields + the
  near-miss advisory carrier) rides JDR 0001 §JD-8 as a correctly-unchecked
  Prerequisite.
- **A9 and A4 — the re-entry's scope — are re-verified against Final 0002, not
  against its superseded spike.** A9's every quoted fence is verbatim in Final
  0002; spike `evidence/spikes/a9-lifted-outcome.md` reproduces 0002's output
  byte-for-byte. A4's census was re-derived from disk at Stage 6: 0002's two
  fixtures already carry `[tags.recognized]` (the JD-10 rename is landed), and
  the one remaining recognized-provenance guard atom is RDR 0003's fixture,
  which is a load *failure* under 0002 rather than a rename this RDR owes. The
  migration inventory is therefore empty, and Phase 2's rename scope was dropped
  on the author's round (3/3 approved).
- **Method vocabulary**: all thirteen in the sanctioned set (compound labels are
  pairs of sanctioned halves) — no paraphrase, no off-vocabulary label.
- **No load-bearing `Docs Only`**: zero `Docs Only` records exist.
- **No Source-Search self-reference**: every Source Search Evidence path resolves
  to `internal/resolve/*.go`, a Final peer RDR's section anchor, the
  `../state-machines` prior-art checkout, or captured spike output — none to this
  RDR or its artifact dir.

## 3. Scope Verification

The Minimum Viable Validation is **in scope, not deferred** — with a split that
is written and gated rather than silent.

- **Kernel half — executes during this RDR's implementation.** A row matching on
  the tag `recognized` fires against a recognized outcome through the real
  kernel's assembled view, and an `Input` supplying an owned/observed tag keyed
  `recognized` (or a row naming it in `RequiresOwned`) is rejected by the
  exported predicate and by `Resolve` at entry. Both run against
  `internal/resolve` as it stands today — no new capability required. Named
  proofs: Testing Strategy scenarios 1, 2, 6, 9, anchored on
  `internal/resolve/resolve.go::recognizedTagKey`, `::assemble`, `::missingOwned`
  and `resolve_test.go::TestReq17_OwnedObservedAndRecognizedTagsAllReachSelection`
  — all symbols confirmed resolving at this gate.
- **Normalizer half — executes inside RDR 0002's implementation**, because 0002
  owns the load/lint code. This is *gated by a written Prerequisite*, not
  deferred by choice, and its normative fixture is captured:
  `evidence/spikes/a9-lifted-outcome.md` § mutant 3 refuses with
  `reserved_tag_key: recognized declaration named "outcome"`, exit 1.

The kernel half alone does not deliver the Problem Statement's outcome, and the
RDR says so explicitly rather than letting the split pass unremarked.

## 4. Cross-Cutting Concerns

- **Canonical-form / determinism** — applies, and is owned here. The reserved key
  is a single literal (`internal/resolve/resolve.go::recognizedTagKey`) compared
  as an exact key string. The RDR rules identity questions explicitly: TOML
  quoting is a surface artifact (`[tags."recognized"]` ≡ `[tags.recognized]`),
  near-spellings (case variants, whitespace-bearing keys) are ordinary unreserved
  names, and duplicate tag keys are admissible input the predicate must *scan*
  for rather than look up (A13, since `Input.Owned`/`Observed` are ordered
  sequences, not maps). No hash, no content-addressed identity, no replay-stable
  byte layout is claimed, so the extended hashing checklist does not apply.
  Determinacy evidence exists regardless: `evidence/repeatability/run-1..3.md`
  (full variant, three distinct models) plus `diff.md` and `resolve.md`.
- **Incremental adoption** — applies. Narrowing `[tags.<tag>]` naming freedom for
  one provenance value invalidated 0002's two canonical fixtures; JDR 0001 §JD-10
  ratified "rename them; there are no users to migrate," and the rename landed at
  0002's re-lock. Migration inventory is empty.
- **Versioning / ambiguous field ownership** — applies, and is the RDR's subject.
  Ownership is settled by the Overrides field: this RDR ratifies RDR 0001's
  deviation D2 from latitude to normative contract, and extends RDR 0002's
  validation-category set with `reserved_tag_key` inside its "including at
  minimum" extensible list (A1). Enforcement reach is carried by an explicit
  Prerequisite, not by an assumed engine traversal (A12).

Omitted as inapplicable: build tool compatibility, licensing, deployment model,
IDE compatibility, secret/credential lifecycle, memory management, concurrency
model, character encoding.

## 5. Proportionality

Right-sized for `Profile: foundational`. The profile is earned rather than
inflated: this is a cross-RDR producer binding RDR 0001's kernel carrier to RDR
0002's declared model and extending 0002's validation surface, and the four
required lenses (cove, 3amigo, critique, repeatability) all ran with complete
evidence dirs.

Nothing to trim before locking. Two things worth naming, neither a blocker:

- **Evidence-field mass** (C9 advisory): six Evidence fields exceed the 30-line
  budget, A10 longest at 70 lines. Answering C9's one question per hit — the mass
  is genuine verification content (the accessor half-sweep, Final-0002 fence
  quotations, the refuted-mechanism record with its Prerequisite repair, the
  disk-re-derived census), and each field's load-bearing anchor is findable in its
  opening lines. The prose has not outgrown the record; no truncation, no
  relocation to `artifacts/`.
- **Single-model critique** (carried caveat from iterations 1–2): `critique` ran
  Pass A + Pass B both on `claude-opus-5[1m]` in fresh contexts, the fallback
  `resources.md` records. Repeatability supplied the genuine multi-model draw
  (three distinct models), so the profile's independence requirement is met in
  aggregate. Recorded, not repaired at lock.
