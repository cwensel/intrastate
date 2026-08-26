# Finalization Gate — cli/0011 flow-next-match-conditioned-candidates

- **Date**: 2026-08-26
- **Verdict**: **PASS** — READY, locked to Final this pass.
- **Mechanical pre-sweep**: `evidence/tooling-pass/tooling-pass.md` — PASS
  (`rdr lint --locking 0011` exit 0; its one `conformance gate:inline` finding is
  resolved by this lock's move of the gate body to this file).

## 1. Contradiction Check

No contradictions found between research findings, design principles, and
proposed solution. Three places where the record could have contradicted itself
were checked directly, because each is a point where a review round rewrote the
draft:

- **Research vs. solution on the absent-key fold.** Investigation reads SQL as
  cutting both ways — UNKNOWN at the atom, fold-and-proceed at statement scope —
  and the solution takes exactly that split: report the distinction at the CLI
  reporting surface (C1's `{key, absent}`), leave the kernel's selection fold
  alone (`0007:REQ-78` stays deferred). The Consequences bullet naming the
  `next`/`resolve` disagreement states the cost rather than denying it, so the
  research and the design agree on both halves.
- **Prior art vs. the chosen default.** Every state-machine peer read
  (stateless, xstate, pytransitions, SCXML) conditions "what next" on state, and
  C1 makes the conditioned form the default with `--all` as the secondary
  surface — matching pytransitions' two-operator shape rather than contradicting
  it. The one asymmetry against the CLI-shape precedents (`docker ps -a` widens
  to rows carrying the SAME information, `--all` here widens to rows carrying
  strictly LESS) is stated in Load-Bearing Decisions → Naming and again in
  Failure Modes, not papered over.
- **Contract vs. contract on filter/dedup order.** C1 requires dedup last; C2
  requires the `--all` filter at emission. These are complementary, not
  conflicting — the repeatability round's F5 fix single-sourced them, and S2
  carries the oracle that fails on the wrong order. Checked the `authority`
  mini-check row, the Existing Infrastructure Audit `Candidate reporting` row,
  and the Illustrative Code intent line: all four sites state the same order.

The `Overrides` field asserts a locked-peer narrowing (`0005:C1`) and the body
must not contradict it. It does not: Approach, C1's demand-set paragraph,
Decision Rationale, and Consequences all state the same `flow resolve` change on
the same one model class, in the same three arms, including the breaking arm.

## 2. Assumption Verification

Twelve Critical Assumption records (A1–A6, A11–A16). Every one is internally
consistent: `Status`, `Method`, and `Evidence` agree, and `If wrong` is
non-empty.

- **Status**: all twelve `Verified`. None `Pending`, none `Unverified`, so the
  status-consistency rule has nothing to catch — no settled-fact prose anywhere
  in the RDR rests on an unverified record. The Prerequisites checklist's four
  `[x]` boxes name exactly the twelve and agree with the records.
- **Method**: `Source Search` (A1, A2, A3, A4, A11, A13, A16), `Spike` (A5, A6,
  A14), `Source Search + Spike` (A12, A15). **Zero `Docs Only`**, so no
  load-bearing claim rests on documentation without a verification plan.
- **Self-reference**: none. Every `Source Search` record's Evidence cites
  `internal/**` symbols or peer-RDR clauses. The five
  `evidence/spikes/*.md` paths belong to `Spike` / `Source Search + Spike`
  records, where a spike output is the sanctioned Evidence, not self-reference.
- **Anchors resolve on `main`**: 95 `source-anchor` edges, all `resolved: true`
  (`rdr inspect --json --filter edges`, `--repo` supplied). No bare `file:line`
  in a normative or Evidence position without a symbol.
- **Adjacent-claim check** — the failure mode where a `Verified` stamp proves
  something narrower than the assumption states. Two records were re-read
  closely because their scope was corrected late:
  - **A12**'s spike records three limits, and limit (c) (two keys each carrying
    multiple tags un-measured) is discharged forward by S3's mandated multi-key
    mixed-state and set-kinded fixtures rather than left as a gap.
  - **A15**'s spike attests the shipped fixtures by suite-green rather than
    payload diff (its limit (b)); S8 carries that as a named build obligation.
    Its breaking arm (plan → `flow-artifact-missing`) was reproduced live, not
    inferred — which is what licenses Consequences and C1 to name and accept the
    class.
  Both limits are recorded in the records themselves and carried by a scenario;
  neither is a stamp proving an adjacent claim.
- **Note, not a finding**: A15 and A16 spell the field
  `**If wrong** (<qualifier>): …`. The projector indexes only the bare spelling,
  so those fields do not appear in `elements[].fields[]` — the text is present
  and non-empty in both (lines 89, 95). Established corpus idiom (cli/0008 uses
  it seven times); `lint` raises nothing. No amendment made.

## 3. Scope Verification

The Minimum Viable Validation is **in scope and executed during
implementation**, not deferred. Its nine steps map onto the three
Implementation Plan phases with no step landing outside them: MVV 1–4 on Phase 1
(probe + reporting) and Phase 2 (`--all`), MVV 5–8 on Phase 3 (fixtures), MVV 9
on `make check` across all three.

**The specific proof** is MVV 2–3, backed by the A6 spike, which already ran it
live: build with the C1 default, seed an artifact for `models/rdr.toml` whose
owned `stage` is `resolved`, run
`flow next --model models/rdr.toml --artifact rdr=<it> --as=json`, and assert
`candidates[]` is exactly `prelock`, `resolve-route-back`, `resolve-abandon`
while `outcomes[]` stays the model's full alphabet. That is the 21→3 narrowing —
the observed defect closed — and it discriminates: a stripped-match build
reports 21, and MVV 4's `--all` re-run must still report those 21.

Two scope facts the RDR states about its own validation, confirmed rather than
taken on faith:

- The MVV is a set of existential fixture checks while C1's success criterion is
  universally quantified. The record closes that gap by enumeration — the
  `disposition` mini-check table lists every reachable input class with an oracle
  per row — and says outright that `21→3` is a witness for one model, never
  proof of the general claim. Honest, and the enumeration is present.
- S6 is an explicit **anti-oracle**: a green 0005 suite is NOT evidence C1
  shipped, because no shipped 0005 fixture has a present-and-unequal match key
  (A4). The discriminating oracles are S2/S3, and C3 mandates the fixtures that
  make them discriminate. Scope therefore includes writing those fixtures, not
  only running the existing suite.

Testing Strategy S1–S8 are all in scope; the closing line is "S1–S8 green and
`make check` passes".

## 5. Proportionality

**Right-sized. Locking as one RDR, not split.**

**Contract count — the split test.** Three labelled contracts, one independent
load-bearing seam. C1 is the `flow next` candidate predicate and owns the
`unknown` payload shape and its reason vocabulary. C2 (`--all`) is that same
predicate's second surface — it is defined *as* "exactly the one 0005:C1
specified", so it cannot be held separately from C1 without one of the two
becoming unstatable. C3 is help text, fixtures, and the mechanical
`unresolved`→`unknown` re-homing census — a build obligation pinning C1/C2's
wording and discriminability, authoring no independent seam. So this RDR is the
sole author of **one** independent contract. No split flag.

The demand-set term in C1 reaches `flow resolve` through
`internal/cli/flow_exec.go::invokedReaders`, which could look like a second seam.
It is not: one function, one module, `internal/resolve` untouched
(`git diff --stat internal/resolve` empty is an MVV assertion), and the term is
a *reading* of `0005:C1`'s existing narrowing clause rather than a new contract.
One seam, not a module span.

**Profile re-validation.** Metadata says `large`. Re-derived against the
contracts just counted: one independent contract plus a user-facing CLI surface
(a new flag, a renamed payload field, rewritten help) and an override of a
locked peer clause — above `mid`, below `foundational` (no new seam, no
cross-module span, kernel untouched). `large` is correct and the field's form is
right: value plus one clause naming the contracts, no matrix or provenance prose
left from the template.

**The lens row matches the Profile.** `large` → grounding → 3amigo → critique,
plus `repeatability` appended by the Stage-5 Determinacy trigger (C1 names
payload identity, entry dedup, and sort order). All four ran with evidence on
disk, and repeatability carries distinct model stamps across `run-1.md`
(`claude-sonnet-5`) and `diff.md` (`claude-opus-5[1m]`). No lens routed past.

**Nothing to trim before locking.** The record is long, and the length was
tested against the standard rather than against a word count: the mass sits in
C1's normative block, the four assumption records whose spikes ran live
(A12, A15, A5, A6), and the Decision Rationale matrix — all load-bearing. Two
specific checks:

- The Evidence fields carrying the most prose (A12, A15, A6) are
  source-verified content with their anchors findable inline. Cutting them would
  blind the grounding sweep that reads those anchors, and the tooling pass's C9
  budget check reports **0 fields over budget**.
- The Decision Rationale's prior-art paragraph (Cedar, OpenTofu, Meyer) earns
  its length by recording a **withdrawn claim** — no source was found for "a
  static over-approximation is preferable because behaviour is a function of the
  query rather than the data", so the determinism argument is re-grounded on
  `0005:REQ-110`/DEV-1 instead. That is exactly the kind of content that must
  not be trimmed.

The four `Alternative N` sections each carry a distinct rejection ground
(O1 empties the list on partial state; O2 leaves the defect as the default;
O4 collides with `0007:C10`; O5 pays a kernel surface change for a distinction
the CLI already has). None is filler.
