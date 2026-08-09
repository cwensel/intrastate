# Triage — RDR 0001 Resolution Kernel

One-pass, unattended roborev triage of the Stage-8 implementation branch.

- **Window**: `BASE 991bd43` .. `HEAD efbeee5` (frozen at Phase 0)
- **Batch label**: `batch:rdr-0001`
- **Spine**: 9 in-window per-commit auto-reviews (jobs 4715, 4717, 4719,
  4721, 4723, 4724, 4725, 4726, 4727) → **30 findings**
- **Out of window**: 8 open jobs at-or-before BASE on `main`/`via-codex`
  (pre-existing; not this run's work)
- **Collapsed**: 30 findings → **6 root-cause clusters**, each grounded by
  one sub-agent against the RDR set + evidence roborev never sees

## Why so much collapsed

roborev reviews each commit in repo-sandbox isolation: it never reads the
RDRs. Every cluster below was raised repeatedly across commits because the
reviewer re-derived the same objection each time, and several rest on
premises the RDR text refutes. Three clusters, however, surfaced **real
unowned cross-RDR seams or under-specifications** that the Phase 3
verifiers did not catch — those are filed as seeds.

## Verdicts

| # | Cluster | Jobs | Verdict | Outcome |
| --- | --- | --- | --- | --- |
| 1 | Escape rows carry `Writes`, conflicting with RDR 0002 | 4715, 4719, 4721, 4723, 4725, 4726 | **RDR-SEED** | seam unowned |
| 2 | Recognized-outcome tag key hard-coded vs RDR 0002 declarations | 4715, 4717, 4719, 4723, 4726 | **RDR-SEED** | binding unbound |
| 3 | Guard-vs-owned precedence (D8 reversal) | 4717, 4719, 4721, 4723, 4724, 4726, 4727 | **RDR-SEED** (+1 DROP) | `RequiresOwned` conflates two roles |
| 4 | Diagnosis payloads sorted at the wrong layer | 4721, 4723, 4725, 4726 | **DROP**: over-engineering | sort is a correct fix |
| 5 | `GuardEvaluator` callback vs RDR 0003 typed atoms | 4715, 4717 | **DROP**: rdr-adjudicated | RDR 0001 delegates |
| 6a | "Commit leaves the suite failing" (High) | 4721 | **DROP**: rdr-adjudicated | red-before-green protocol |
| 6b | REQ-11 import denylist not exhaustive | 4715 | **DROP**: over-engineering | claim matches coverage |
| 6c | Uniformity test miscredited with FAIL-2 | 4726 | **FIX-NOW** | `coverage.md` corrected |
| 6d | `Table` omits tag declarations for write validation | 4715 | **DROP**: project-scoped-out | RDR 0002 owns validation |

## Grounding that overturned the reviewer

**Cluster 2 — premise factually false.** RDR 0002 L215/L333 use `recognized`
as one of three **provenance values** (`owned`, `observed`, `recognized`) —
never as a tag key. `outcome` appears only as an author-chosen tag *name* in
the spike fixtures (`[tags.outcome]` / `provenance = "recognized"`), nowhere
in the RDR body. The reviewer's "RDR 0002 declares the recognized tag as
`outcome`" is wrong. The *underlying* question — who names the key across
0001/0002 — is real and was already pre-registered in deviation **D2**.

**Cluster 1 — premise right, consequence wrong.** RDR 0002 L287-288 does say
"An escape rule MUST contain an `escape` list and MUST NOT contain a write
block or clear list." But RDR 0001 L296 delegates table shape outright:
"Kernel assumes a parsed, reviewable table shape." Malformed-shape rejection
is RDR 0002's validation duty (its L346 lists "malformed escape
declaration"). The kernel is not defective; the *seam* is unowned.
Verified separately: **no Phase 3 assertion depends on escape fixtures
carrying `Writes`** — strip them and every test still discriminates. The
Phase 3 evidence stands.

**Cluster 4 — refuted by RDR 0001's own identity rule.** REQ-2 defines the
input tuple as keyed on the transition table **revision**, an opaque value
the kernel "carries and compares but does not parse" — not on the row
sequence. Two permutations at `rev-1` therefore *are* the same input tuple,
so differing payloads **are** a determinism breach and the sort is a correct
fix. RDR 0002's ordering guarantee (L318-323) is scoped to the **dump**, not
to the slice handed to the resolver; and no normalizer exists at HEAD, so the
kernel is the only boundary that exists.

**Cluster 5 — RDR 0001 explicitly defers.** Capability Dependencies:
"Guard predicate evaluation | RDR 0003 | Deferred to peer". REQ-23 names
"guard evaluation **delegation**" as the Phase 2 deliverable, and Briefly
Rejected records host-language callbacks as "Rejected *here* because RDR 0003
must first define a predicate shape that lint can reason about." Demanding
typed atoms now would be implementing an unimplemented RDR. RDR 0003 is
**Final but not Implemented**.

**Cluster 3 — probed, not argued.** The grounding leaf ran live probes
against HEAD rather than reasoning from review text:

- Overlap key **present as observed** → `owned_state_unavailable` (correct;
  `gate()` fixed the jobs-4717/4719 shape) → `DROP: superseded-at-HEAD`.
- Overlap key **absent entirely** → `no_match`, which is *correct* per the
  req-list ASSUMPTION's own exclusion ("not when a candidate merely fails to
  match on a present tag").
- Guard FALSE + absent `RequiresOwned` key + modeled escape → plan via
  escape. **Reachable and demonstrated** — but wrong only under a reading the
  RDR does not state.

No input produces an unambiguously wrong disposition, so this is
under-specification, not a shipped defect. `RequiresOwned` currently means
both "state the guard reads" and "state the transition needs"; only the RDR
can separate them. **Do not revert D8** — that re-opens ADV-1/ADV-1b, which
two independent verifiers proved.

Job 4726's "verification.md says closed but D8 says open" inconsistency does
not hold: `verification.md` closes FAIL/ADV **defects**; D8 is a SPEC-UNDER
awaiting ratification, and `status.md` carries it explicitly as
unratified-but-non-blocking. Different registers, both true.

## FIX-NOW applied

`coverage.md:107` claimed the uniformity test "would have caught FAIL-1 and
FAIL-2 together". Its three cases (missing-owned/guard-TRUE, undecidable
guard, nil seam) never combine guard-FALSE with missing owned state — FAIL-2's
exact shape. Corrected to credit **ADV-1/ADV-1b** with FAIL-2 and narrow the
uniformity claim to FAIL-1. Deliberately documentation-only: adding the
collision case would pin the semantic D8 leaves unratified.

## Carried open questions

Both pre-existing SPEC-UNDER deviations remain open for author ratification
and are **not** blockers:

- **D1** — escape scope: RDR 0001 REQ-15 lists four escapable conditions;
  RDR 0002 L292 restricts escape lists to `no_match`/`ambiguous_match`. The
  narrow reading is implemented. Cluster 1's grounding independently
  corroborates the narrow reading.
- **D8** — guard-FALSE vs missing-owned precedence. Recommended resolution
  (from cluster 3): settle in RDR 0003 by making `GuardUnevaluable` the
  mandatory verdict when a predicate references absent state. That closes the
  masking path without touching `resolve.go` or weakening a frozen test.

Ratify D1 and D8 **before RDR 0003 is implemented** — 0003 is already
Final, and its guard contract's implementation depends on how D8 resolves
(the recommended fix lands *in* 0003's guard-evaluator contract).

## Noted, not filed

- `escapeRow` (`fixtures_test.go:249`) could drop `Writes`/`NextTags` for
  fixture realism. Changes no assertion outcome; tidy-up, not evidence repair.
- `GuardEvaluator` has no documented purity obligation, though RDR 0001's
  replay-determinism claim rests on it. Unreachable today (only test
  fixtures implement it); absorbed by RDR 0003's implementation.
- When RDR 0002 is implemented, its verification should assert the normalized
  candidate-row **slice** is identity-ordered, not only the rendered dump.
