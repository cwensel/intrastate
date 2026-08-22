Model: claude-opus-5[1m]

# 3amigo Consolidation — iteration 2 (re-entry)

Three isolated persona passes (PM 6, Implementer 8, QA 8 = 22 findings), each
delta-scoped to the re-entry (A5, A7, A8, and the passages stating or depending
on the narrowed lint promise). No persona file references another persona's
output; isolation held. Consolidation is **mechanical**: hotspots are computed
by which RDR passage two or more personas independently named, not by re-judging
their findings. Overlap marks a hotspot passage, not a validated finding; a
single-persona finding is not thereby weaker.

## Origin ledger (this iteration)

| # | Passage anchor | Personas | Hotspot | Concern |
| --- | --- | --- | --- | --- |
| **1** | Normative Contracts — narrowing clause: "withheld"/"MUST NOT certify"; `disposition` table's last two rows; MVV withholding sentence; Testing Strategy Scenario 6 | P1-1, P1-2, P3-2, P3-6 | **HOTSPOT (3)** | The narrowed promise has no observable lint form. Marked **Loud** but the only named artifact is RDR 0007's *runtime* payload; no lint finding code exists for a claim withheld on a product that *does* prove finite. Asserted only negatively ("withholds", "does not certify"), which an implementation emitting nothing satisfies. "can refuse" is a modal with no decision procedure (syntactic vs semantic readings give opposite outcomes). |
| **2** | "lint MUST refuse **or downgrade**" — Normative Contracts (×2), A2, Technical Design, `disposition` rows, Testing Strategy Scenario 3 | P2-4, P3-1 | **HOTSPOT (2)** | Free disjunction with no discriminator; "downgrade" has no finding code, severity, or exit behavior in any document. §JD-4 leaves "whether lint gains a warning category" explicitly open; RDR 0006 (`0006:345-349`) states only the blocking arm. Scenario 3 passes under either reading. |
| **3** | Normative Contracts — narrowing clause's closing MUST: "RDR 0006 … MUST carry this same narrowing; the two documents MUST NOT state it differently"; A8 | P2-3, P2-5, P3-5 | **HOTSPOT (3)** | A normative MUST binding a `Final` sibling this RDR has no authority over, constraining document *prose* rather than code — no test surface, no conformance check. Verified: RDR 0006 has zero occurrences of `guard_unevaluable`/"narrow". A8's own Plan agrees the fix is a route-back. A8's two convergence arms ("same wording" vs "assign to one document") imply different implementations, and the "same wording" arm recreates the duplication A8's "If wrong" forbids. The `authority` table says "Contested"; the normative block reads settled. |
| **4** | A7 `Pending`; `trace` Step 5 ("this step has no defined result"); MVV operator list; Phase 2; fixture rule `foundational-to-cove` | P1-3, P2-2, P3-4 | **HOTSPOT (3)** | `exists` is inside the closed operator vocabulary and is the sole TOTAL operator, but its projection is undefined, so Phase 2 has no defined behavior for it and the fixture's only multi-row group is exactly the undefined case. The scope consequence (no exhaustiveness for any `exists`-bearing group, including RDR 0007's sanctioned two-row absence pattern) is reachable only from an assumption's failure branch — the operator matrix, Consequences, and Proportionality carry no carve-out. No interim contract is stated. |
| **5** | A7 Plan + Capability Dependencies — per-tag optionality request to RDR 0002 | P1-4, P2-8, P3-3 | **HOTSPOT (3)** | The blocking producer request is recorded only here; RDR 0002 (`0002:217-218`) has no optionality field and is itself `Draft`, so no document has accepted the request and the wait has no termination condition. Consequence for test: MVV Scenario 6 requires a "possibly-absent guard key" the fixture format cannot express, contradicting Scope Verification's "in scope for implementation, not deferred". Capability Dependencies names a *second* producer request (set-valued element encoding) with no assumption ID or status while Phase 3 gates on `contains`. |
| **6** | A5 Evidence + Load-Bearing Decisions › Identity | P2-1, P3-7 | **HOTSPOT (2)** | Identity is defined as source rule/context id plus "**position within** `all` or `unless`", but the atom shape A5 cites (`0007:1249-1266`: `Key`/`Operator`/`Literal`/`Block`) carries no position field. A slice index collides with Scenario 5's reorder invariance; `(Key,Block)` cannot distinguish RDR 0007 A5's two-atoms-one-key row; a fifth field is a change request to `Final` RDR 0007. A5 asks only a row-level source-identity assertion while Identity is defined per-atom — no atom-granularity regression criterion. |
| **7** | Technical Design — row-group grouping delegated to RDR 0006 | P2-6 | single | RDR 0006 invariant 4 gates on a group "that claims closed coverage", and no document defines how a group makes that claim. Blocks opt-in vs default-on lint. |
| **8** | Prerequisites — the single all-assumptions checkbox | P2-7, P3-8 | **HOTSPOT (2)** | One checkbox gates on A7 and A8, whose closures differ in kind (a peer's field + an undrafted normative clause vs. amending a `Final` peer / a JDR disposition); neither names a verifier or artifact, so no mechanical pre-lock check is possible. Separately the RDR 0007 sequencing item is `[x]` while the shipped kernel still carries `Row.Guard string` (`internal/resolve/resolve.go:185`) and `Evaluate(guard string, view TagSet)` (`:91`). |
| **9** | Decision Rationale — "lint can still prove whether those edges are complete and mutually exclusive" | P1-5 | single | States the outcome unconditionally directly above two Pending assumptions that each remove a class of edge from provability; RDR 0006 and RDR 0005 inherit the unqualified promise. |
| **10** | A5 "If wrong" | P1-6 | single | No longer names the condition that produces it: the identity-preserving shape is specified but unshipped, so the real trigger is RDR 0007's reshape deviating, not a normalization step losing identity. |

## Hotspot summary

Six of ten entries are multi-persona hotspots. Entries **1, 2, 3** converge on one
root: *the narrowing's behavioral outcome is never given an observable lint form* —
three personas reached that passage from three different questions (does the author
get an outcome / what do I code / what do I assert). Entries **4 and 5** are the two
faces of A7: the undefined projection and the unaccepted producer request that
blocks it.

## Not filed

- PM checked A8's factual claims against `docs/jdr/0001-resolve-kernel-seam.md:238-241`
  and `docs/rdr/0006-graph-lint-authority-and-guarantees.md:344-350` — clean.
- QA recorded four passages checked and found testable as written: Testing Strategy
  Scenarios 4 and 5, the Phase 3 `contains` gate, and the Decision Rationale /
  Contradiction Check §JD-4 narrative.
- The `contains` element-encoding gap is real but already charted in Capability
  Dependencies and Phase 3; PM judged it outside this delta-scope (it surfaces here
  only as part of entry 5).

No OUT-OF-SCOPE findings were reported by any persona.
