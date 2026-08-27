Model: claude-opus-5[1m]

# Finalization Gate — cli/0010 (0010-stateless-decision-tables)

**RDR**: 0010 — Owned state is optional: stateless decision tables are
first-class
**Date**: 2026-08-26
**Verdict**: **READY — Gate PASS.** Mechanical sweep PASS
(`evidence/tooling-pass/tooling-pass.md`); no blocker; lock.

Responses 1, 2, 3 and 5 are recorded here. `### Cross-Cutting Concerns`
(item 4) stays in the record as `cli/0010:G-cross-cutting`, because peers
cite it.

## 1. Contradiction Check

**No contradictions found** between research findings, design principles, and
the proposed solution.

Three places where a conflict was plausible, each checked and clean:

- **Prior art vs. the declared class.** Research Findings records that peers
  *require* an initial state and that pytransitions silently injects a default
  one — prior art for the dummy-state workaround, i.e. for the very thing the
  Problem Statement rejects. The solution does not follow that prior art; it
  reads it as the negative case, and the Decision Rationale's scored matrix
  says so explicitly (A's prior-art cell: "matches pytransitions' auto-injected
  `'initial'` state — the workaround peers had to add an opt-out for"). Citing
  prior art as a cautionary case, with the reasoning written down, is not a
  contradiction.
- **"Declared, not inferred" vs. the empty owned set being the real
  discriminator.** The Approach argues the sibling in-repo discriminator
  (`Row.Kind`) *is* presence-inferred, then declines to follow it — on the
  stated ground that nothing in the tree infers a **class** from an absence and
  that 0006 rests on "an absent root is a defect, never an empty graph"
  (`0006:C18`). The one-directional check (C1) is what keeps both readings
  true: a declared table carrying owned state is refused, a zero-owned machine
  stays a lint finding. ALT2 is the inferred design, weighed and rejected on
  the misclassification and overrides rows. Consistent.
- **"Coverage runs unchanged" vs. the new `graph-unprovable-coverage` arm.**
  C5 both asserts the 0006 invariants run unchanged and adds a class-keyed
  emission. These are not in tension: the arm fires over **guard** dimensions'
  absence, and the unchanged claim is scoped to "over the authored guard
  dimensions". Research Findings carries the same distinction (match atoms
  scope the group and contribute no assignment, `0006:C7`/`0003:C13`), and the
  Risks section names the match-discriminated silent-green as the defect the
  fence closes. The RDR is internally consistent about which atoms participate.

Planned features vs. stated principles: the RDR's principle is "intrastate
stays generic — no consumer knowledge, no artifact discovery" (Problem
Statement). Nothing in C1–C5 encodes consumer knowledge; the MVV explicitly
routes the consumer rewrite to rdr#tmxk in the consumer repo and leaves
`models/rdr.toml` untouched. `emit` values are literal strings with no
interpolation (C3, A6), which is the same principle held at the value level.

## 2. Assumption Verification

**All 15 Critical Assumption Evidence Records are internally consistent and
Verified.** Zero `Pending`, zero `Unverified` — so the status-consistency
sub-check has no candidate: no settled-fact prose anywhere in the RDR can be
resting on an unsettled assumption.

- **Status / Method / Evidence agree** on every record; each `If wrong` is
  non-empty and states a real consequence (e.g. A3's "the class needs a kernel
  change RDR 0001 owns", A2's "the problem statement's refusal survives").
- **No `Docs Only` record exists**, so nothing blocks on the Docs-Only rule.
  Methods used: `Spike` (A1, A8), `Source Search` (A2, A3, A4, A15),
  `Prior Art` (A5), `Peer RDR` (A6), `MVV Test` (A7, A9–A14) — all sanctioned
  (tooling-pass CHECK 2, `off_vocabulary` empty throughout).
- **No `Verified` stamp is self-referential.** The four `Source Search`
  records resolve to `internal/` symbols, never to this RDR or its artifact
  dir (CHECK 3). Nor does any prove only an adjacent claim: the five records
  re-opened at pre-lock were each re-verified against the thing actually in
  question — A6 against the consumer's real emit keys (the critique lens had
  correctly refuted the old supporting reading), A12 against twelve enumerated
  consumers, A13 by executing a spike package against the real loader and
  engine, A14 by a green full-suite run on a scratch copy, A15 by tracing all
  four readers to the model where each runs.
- **Every cited `path::Symbol` resolves.** All `source-anchor` edges report
  `resolved: true`; none unresolved, none skipped (CHECK 5).

Two records carry a load-bearing nuance worth naming at lock, both already
written into the record rather than left to the implementer:

- **A14** is the one clause not additive on its owner's grammar — appending
  `no-participating-dimension` to 0006's `reason` set. Its assent is
  **structural, not solicited**: 0006 declares the set "closed, append-only",
  and reserves "closed at …" for sets that may not grow (`0006:C17`). That
  reading is recorded in the Prerequisites checkbox, so lock does not depend
  on reopening 0006.
- **A12** surfaced an implementation trap and C5 already mandates the fix: the
  root-arm predicate must be the augmenting OR
  (`class == "decision-table" || len(Initial) > 0`), never class alone, which
  would silence the root arm for rootless *state-machine* models.

## 3. Scope Verification

**The Minimum Viable Validation is in scope and executed during
implementation, not deferred.**

The specific proof is the six-step MVV (`0010:MVV`), whose load-bearing step
is **step 2 → step 3**: a deliberately partial decision-table model, with two
`single_valued`/`required` guard dimensions of two values each, must lint to
**exit 2 with one coverage finding naming the uncovered cell** — a *positive*
finding, which is what proves coverage actually ran over a model with no root
declared — and then, with the fourth rule added, lint to **exit 0 with
`findings` exactly `[]`**. A lint that never ran and a lint that ran clean are
otherwise indistinguishable; requiring the positive finding first is what
makes the pair discriminating. Step 4 pins the answer
(`data.rule` = the fourth rule's id, `data.emit` = its authored block,
`data.next`/`writes`/`owned` empty, `data.readers` `[]`, no `--artifact`
supplied), and step 6's negative controls pin the non-changes (`class`
omitted → `graph-dangling-edge`; an owned tag added → load refusal;
`models/rdr.toml` unchanged).

This is carried in the Testing Strategy as scenarios S1–S6 with a Trace, so it
is scheduled implementation work, not an aspiration. Nothing MVV-critical was
deferred past lock: the Stage 6 reconcile applied that rule explicitly and ran
A13 and A14 **now** — both are MVV prerequisites, because step 3 asserts the
empty finding set and the match-discriminated control asserts C5's fence —
rather than downgrading either.

One scope boundary is stated and correct: the **consumer** rewrite
(rdr#tmxk's navigator model) is *not* in this RDR's MVV. intrastate#zdat
closes on the generic fixture passing; the consumer rewrite is tracked
separately. The RDR names these as distinct completions where neither
substitutes for the other — that is a scoping decision, not a deferral of the
validation.

## 5. Proportionality

**Right-sized. Nothing to trim before locking. No split.**

- **Contract count — one seam, not five.** The split test is contract count,
  and this RDR labels five (C1–C5). They are **facets of one independent
  load-bearing contract**, not five seams: the question is "is a zero-owned-tag
  model a legal class, and what does resolve answer over it". C1 declares the
  class, C2 is what that class means for rule shape, C5 is what it means for
  lint, and C3/C4 are the answer such a model returns and how it reaches the
  caller. Remove any one and the others are unimplementable or pointless — C5's
  ∅ root is meaningless without C1's class; C4 has nothing to carry without C3.
  Splitting would put a grammar key in one record and its only consumer in
  another.
- **Profile re-validated: `foundational` is correct, and the lenses that ran
  agree.** The field's own clause says five contracts are facets of one seam
  spanning `internal/table`, `internal/graphlint`, and `internal/cli`, produced
  for 0002/0005/0006 — which matches the contracts just counted and the
  three-seam Technical Design. The full `foundational` lens row ran and
  completed: `cove` (findings + step0-grounding), `3amigo` (3 personas +
  consolidation + resolve), `critique` (dual-model, `claude-opus-5` /
  `claude-sonnet-5`, converged), `repeatability` (**full** variant — run-1/2/3
  across three distinct model stamps, plus diff and resolve). No lens was
  skipped on a wrong Profile, and no lens outcome argues the Profile was too
  large. No `Transient` contract is involved. The field's form is correct:
  value plus one clause naming the contracts, no template/Seed matrix prose.
- **Length is proportionate to the seam, and the mass is where the risk is.**
  1951 lines is large, but it is concentrated in Critical Assumptions (525) and
  the Normative Contracts (334) — the two places a wrong call is expensive.
  C5 alone is ~100 lines because it fixes an emission *site* and a precedence
  rule that were both proven necessary by execution (A13 established that the
  obvious insertion double-reports). The five over-budget Evidence fields
  (tooling-pass CHECK 9: A14 57, A10 49, A13 43, A9 31, A12 31) are verification
  content, each with its anchor at the head and a spike file pointer; cutting
  them would blind the grounding sweep that reads those anchors. **No
  truncation proposed.** The Alternatives, Briefly Rejected (6 entries), and
  Decision Rationale matrix are each terse relative to what they decide.
