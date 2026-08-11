# Finalization Gate — RDR 0008: Ownership of the recognized-outcome tag key name

- **RDR**: `0008-recognized-tag-key-ownership.md`
- **Date**: 2026-08-11
- **Profile**: foundational
- **Mechanical pre-sweep**: PASS (`evidence/tooling-pass/tooling-pass.md`, C1–C6 no findings)
- **Verdict**: **PASS — READY to lock**

## Contradiction Check

No contradictions found between research findings, design principles, and
proposed solution. Three places where a conflict could plausibly have hidden
were checked directly, because each is a spot where the rounds moved text:

1. **Prior art vs. the chosen mechanism.** Key Discoveries reports that
   surveyed engines fix the injected datum's name in a namespace *disjoint*
   from author vocabulary (XState's `GuardArgs.event`, SCXML's `_`-sigil
   reservation), while this RDR carves a bare word out of the author's single
   flat `[tags.*]` namespace. That is a genuine gap, and the Problem Statement
   concedes it in text rather than papering over it — the prior art supports
   the *posture* (carrier owns the injected datum's name), not the
   *mechanism*, which is exactly why this RDR needs a validation category
   where XState needs none. Concession stated, not contradicted.

2. **"No kernel change" vs. the settled Enforcement locus.** The Approach's
   original framing (ratify D2, no code change) would contradict block 4's
   exported `Input` predicate and `Resolve`-entry call. It does not, because
   the claim is scoped in three places to *disposition for conforming input*:
   Load-Bearing Decisions ("Surface conceded, stated plainly"), Trade-offs
   (third Negative), and the Decision Rationale's re-scored blast-radius row,
   which explicitly records that R "is no longer the zero-kernel-change branch
   it was chosen as" and that the margin narrowed. The QOC ranking is
   re-derived under the settled cost rather than inherited from the
   pre-settlement scoring.

3. **A4's finding vs. the Approach.** Research findings establish that all
   three committed recognized-provenance declarations violate the rule's first
   clause — every observed author exercised the freedom this RDR withdraws.
   The Proposed Solution does not contradict this; Phase 2 carries the rename
   of three declarations plus six reference sites, and the Problem Statement
   names the cost against observed convention. A finding that cuts against the
   decision is carried as cost, which is the conforming shape.

The Failure Modes section likewise states the intent-channel gap (a
conforming model whose row still never fires) rather than letting block 2 be
read as a guarantee it does not make.

## Assumption Verification

Thirteen Critical Assumption Evidence Records, all internally consistent —
Status, Method, and Evidence agree, and every record carries a non-empty "If
wrong" or an explicit statement of what its falsification changes.

- **Status distribution**: A1–A11, A13 `Verified`; **A12 `Refuted as
  stated`**. Zero `Pending`, zero `Unverified`.
- **No `Docs Only` records** — the category that blocks lock is empty (C4).
- **No `Source Search` self-reference** — every Source Search record (A2, A4,
  A7, A8, A10, A12, A13) resolves to `internal/resolve/*.go`,
  `internal/cli/respond/respond.go`, peer RDR bodies, peer spike fixtures, or
  the RDR engine's prompt tree; none cites this RDR or its artifacts (C3).
- **Every cited `path::Symbol` resolves at HEAD `1caf456`** (C5, delegated
  verification): `recognizedTagKey`, `assemble`, `Resolve`, `missingOwned`,
  `gate`, `evaluateGuard`, `escapeOrRefuse`, `refuse`, `Input`, `Tag`,
  `TagSet.has`/`Lookup`/`Len`, `Row`, `ProvenanceOwned`,
  `ProvenanceRecognized`, `Table.models`, `recognizedTagSensitiveTable`,
  `fixtureGuards.Evaluate`, `mustResolve`,
  `TestReq17_OwnedObservedAndRecognizedTagsAllReachSelection`. Two
  independently re-checked load-bearing properties held:
  `fixtureGuards.Evaluate` does discard the `TagSet` (so A2's test gap is
  real and scenario 3 is net-new coverage), and `Resolve`'s body returns
  `nil` on all four paths with no `errors.`/`fmt.Errorf` in the file (so A8's
  "first non-nil error path" is accurate).
- **Peer-RDR anchors resolve**: RDR 0002's "including at minimum" fence and
  its canonical-examples sentence; RDR 0009 carries zero `Input` occurrences,
  confirming A6's negative conclusion and A11's disjointness leg; RDR 0007
  (`Final`) cites this RDR's producer sentence as a peer's *name* constraint,
  which is what closes A7's ownership half without a concurrence round.
- **A4's inventory verifies exactly**: three declarations
  (`rdr-fixture.toml:25`, `kata-fixture.toml:22`, `guard-fixture.toml:36`)
  and six reference sites (`rdr-fixture.toml:59,74,86`,
  `kata-fixture.toml:47,57`, `guard-fixture.toml:72`) — including the guard
  site the critique lens found missing at iter-2, now present and propagated
  into Phase 2.

**A12's refutation is not a blocker, and the reason is recorded rather than
asserted.** A12 claimed a *discoverability mechanism* (that RDR 0002's
implementer would encounter this constraint via the implementation prompt),
not a design premise. Its falsification changes no normative contract — the
rule, the category, the payload, and the predicate all stand — but it does
change enforcement *reach*, and the RDR repairs that in-document by promoting
the handoff to an explicit **Prerequisite** ("this RDR's Normative Contracts
are a named input to RDR 0002's Stage 8 run, or a pointer lands in 0002's own
surface"). The obligation is settled here; only its *form* routes to Stage
7.1. That is a refuted assumption converted into a written gate, which is the
disposition the flow wants.

**Status consistency**: no settled-fact prose depends on an unsettled
assumption (C6). A11's split status (`Verified` for independence; report
order a Stage 7.1 item by design) is matched by block 4, which licenses an
interim rule — "either error is conforming; what is *not* conforming is
skipping this check because another fired" — rather than asserting an order
neither RDR may set unilaterally. The Prerequisites checklist's three
unchecked boxes gate *implementation*, not lock, and do not contradict this
gate.

## Scope Verification

The Minimum Viable Validation is **in scope and split across two owners**,
and the split is a real property of the seam rather than a deferral.

**Kernel half — executed during this RDR's implementation, runnable at HEAD:**

- A row matching on the tag `recognized` fires against a recognized outcome,
  asserted through the **real kernel's** assembled view — never by comparing
  two hardcoded literals (scenario 1's kernel half; the behavioral form
  premortem P-5 requires).
- An `Input` supplying an owned or observed tag keyed `recognized` — or a row
  naming it in `RequiresOwned` — is rejected by the exported predicate *and*
  by `Resolve` at entry, in three variants including the duplicate-key case
  A13 surfaced (scenario 6), with the kernel-side `owned_state_unavailable`
  residual and its predicate-side breach both asserted (scenario 9).

Done for this RDR's own implementation is scenarios 1's kernel half, 3, 6, 8,
and 9 green plus Phase 1's pointer comment.

**Normalizer half — executed inside RDR 0002's implementation** (scenarios 2,
4, 5, 7 and the MVV's lint half): a table declaring a recognized-provenance
tag named `recognized` lints clean; the same table renamed, and separately
with an owned tag named `recognized`, fails load/lint in the
`reserved_tag_key` category with the direction-appropriate remedy name and
rule identifier.

This half is **not deferred in the Gate's sense** — it is scoped to the peer
that owns the code (no normalizer exists at HEAD; the reuse audit confirms
it), and it is gated by a written Prerequisite rather than by intent. The RDR
states the binding consequence itself: every scenario delivering the Problem
Statement's user outcome sits in the carried half, so "this RDR's Close MUST
NOT claim the Problem Statement outcome until the carried half is green."
Declaring that in the document is what makes the split honest rather than an
accounting artifact, and it is the right disposition for a cross-RDR producer
whose peer is Final-but-unimplemented.

## Cross-Cutting Concerns

- **Character encoding** — load-bearing here, and addressed. Identity is
  byte-exact on the **post-parse** key string: no case folding, no trimming,
  no aliasing. TOML quoting is resolved explicitly as a surface artifact
  (`[tags."recognized"]` and `[tags.recognized]` parse to the identical key
  and are both reserved), while `Recognized`, `RECOGNIZED`, and
  `" recognized"` are ordinary unreserved names. Because that rule *creates* a
  human-invisible near-miss hazard, the non-blocking advisory is settled as
  normative alongside it, with a deliberately **disjunctive** trigger (fold
  alone or trim alone), pinned by scenario 4 on all three near-misses plus a
  negative control (`[tags.result]`).
- **Incremental adoption** — addressed, with the cost named. The rule's first
  clause breaks all three committed recognized-provenance declarations on day
  one, including the two RDR 0002 designates canonical. A4 carries the
  inventory, Phase 2 carries the edits, and A9 (Spike, against RDR 0002's own
  normalizer) establishes the rename is mechanical — every tag predicate lands
  in the candidate row's predicate set, the normalized row has no outcome
  field, so no gating semantics move. The project carries no back-compat
  obligation, so no migration shim is owed.
- **Versioning / API compatibility** — addressed. The kernel gains one
  exported `Input` predicate and its `Resolve`-entry call; `Resolve`'s
  signature does not widen (already `(Result, error)`), and the error channel
  is one RDR 0001 REQ-6 affirmatively reserves for programmer mistakes.
  Conforming callers see no behavior change; the predicate's exported name is
  implementation latitude and MUST NOT bind to any symbol from the still-Draft
  RDR 0009.
- **Determinism** — applies in the weak form only, and is satisfied. This RDR
  claims no hash, no content-addressed identity, and no byte-identical output,
  so the hash/pre-image checklist does not apply. What it does claim is that
  RDR 0001 D3's precedence remains the deterministic backstop for input
  bypassing the precondition; scenario 8 pins that the same non-conforming
  input yields the same disposition on repeat runs, with the honest scope note
  that A5's *unreachability* half is a derivation, not a test, because the
  overwrite in `assemble` has no observable at the package boundary.
- **Concurrency model** — does not apply; omitted from the response rather
  than N/A-bulleted, along with build tooling, licensing, deployment model,
  IDE compatibility, secret lifecycle, and memory management.

Cross-RDR policy ownership, stated where it is not this RDR's: the
user-facing surface and exit-code mapping for the `reserved_tag_key` category
belong to **RDR 0005** (charted, not settled here); the *meaning* of
`RequiresOwned` belongs to **RDR 0007** (cited, not restated, per 0007's own
convention); the normalized-graph lint layer belongs to **RDR 0006**, and
Phase 2 argues explicitly that declaration-name validation is upstream of
0006's input contract rather than a competing acceptance rule.

## Proportionality

**Right-sized; nothing flagged for trimming.**

**Contract count: one.** This RDR is the sole author of exactly one
independent load-bearing contract — the identity of the recognized-provenance
tag key. The six normative blocks are that one contract expressed across the
three channels the key can arrive through (declaration: blocks 1–3;
resolve-time data: block 4; row-level owned-key references: block 5), plus
block 6 which is a negative contract stating what does not change. They are
not separable seams: reserving only the declaration channel leaves the
silent-shadowing hole this RDR exists to close, and reserving only the data
channel leaves the load-time diagnosis unfixed.

The RDR does not merely assert "three channels" — it owes and pays a
**closure argument**, derived from the code rather than from inspection of its
own text: `assemble` is the sole constructor of the view (A2, whole-package
sweep) and writes exactly three sources; rows reference view keys through
`Match`, `RequiresOwned`, and the guard seam, of which `RequiresOwned` is the
one field read *by provenance*. A fourth channel would require a second view
constructor or a new provenance-sensitive row field — both kernel changes
this RDR would have to reopen anyway. The ≥2 split signal is therefore not
tripped.

**Profile re-validated: `foundational` is correct and stands.** It is earned
on the cross-RDR axis, not on contract count — this RDR binds RDR 0001's
kernel carrier to RDR 0002's declared model, extends 0002's validation
surface, cites 0007 for `RequiresOwned` semantics, borrows 0009's pattern, and
charts to 0005. The accretion axis does not force it independently (`Seam
Lineage` records no prior accretion at this locus), so the contract axis'
cross-RDR-producer trigger is what sets it. The lenses the matrix requires for
`foundational` all ran and are resolved: `cove`, `3amigo`, `critique` (passes
A and B plus an iter-2 delta), and `repeatability` in the **full** variant —
three runs across three distinct base models (`claude-opus-5[1m]`,
`claude-fable-5`, `glm-5.2:cloud`) plus `diff.md` and `resolve.md`. The
determinacy requirement for a contract naming identity is satisfied by that
full repeatability set rather than by an `n/a` disposition.

The Profile field's form is correct: value plus one clause naming the
contract, with no matrix or provenance prose left from the template.

Length is proportionate to a foundational cross-RDR producer whose
enumeration was extended twice under review — the closure argument, the
direction-specific payload split, and the A12 repair are each load-bearing
text that a shorter document would have had to omit at the cost of the
contract. No section is redundant with another, and the change-history
narration Refine removes has not crept back.
