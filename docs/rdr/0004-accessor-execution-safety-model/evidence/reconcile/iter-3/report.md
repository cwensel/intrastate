Model: claude-opus-5[1m]

# Reconcile Report — iter-3 (re-entry, A8)

RDR: `0004-accessor-execution-safety-model`
Status line: `Draft [revised from Final 2026-08-12; re-verified A8 — …]`

## Stage 5 preflight — PASS

`Profile: large` → lens row (§lens-row) is **grounding → 3amigo → critique**.
The iter-2 reconcile returned NOT RECONCILED because 3amigo and critique had not
run for this iteration. Both have since run.

| Lens | iter-2 evidence | Complete |
| --- | --- | --- |
| grounding | `grounding/iter-2/{findings,dispositions}.md` | yes |
| 3amigo | `3amigo/iter-2/{consolidation,dispositions,persona-1-pm,persona-2-implementer,persona-3-qa}.md` | yes |
| critique | `critique/iter-2/{critique,critique-modelB,diff,grounding,dispositions}.md` | yes (dual-model: Pass A `claude-sonnet-5`, Pass B `claude-fable-5`) |

Determinacy trigger: discharged — `grounding/iter-2/dispositions.md` carries the
written `Determinacy trigger: n/a` with a reason. No `repeatability` row entry is
appended, so no `run-1.md`/`diff.md` is owed and no variant mismatch exists.

All four mini-check tables are in the draft: Disposition Table, Oracle
Discriminability, Fidelity Table, Desk Trace.

## Open set

Pre-Lock needs-verification lists pasted from `grounding/iter-2`,
`3amigo/iter-2`, and `critique/iter-2` dispositions.

| item | source | disposition | evidence pointer or plan |
| --- | --- | --- | --- |
| A9 — read-completeness seam/boundary rules implementable as stated (4 independent rules) | 1, 2, 4 | **DOWNGRADED** | `Status: Pending — DOWNGRADED at Stage 6`, `Method: MVV Test`. Named plan: MVV Scenario 1 arm 8 (missing/empty requested key set), Scenario 6 (absence → `owned_state_unavailable` at the resolver), Scenario 7 (timeout precedence), Scenario 3 (`read_back_incomplete`). Survivable: "If wrong" names the fallback — move the absence encoding into RDR 0001's `Input` shape (JDR 0001 §D3 option (c)). Retires per-rule. |
| A10 — post-mutation refusal reports applied-but-unverified without compensation | 1, 2, 4 | **DOWNGRADED** | `Status: Pending — DOWNGRADED at Stage 6`, `Method: MVV Test`. Named plan: MVV Scenario 8 — refusal carries the applied-but-unverified sense, exactly one write invocation, no undo/retry/re-derivation. Survivable: "If wrong" names the two fallbacks (caller treats every write refusal as possibly-applied, or the RDR claims transactional semantics Alternative 4 declined). |
| A1–A8 | 2 | no change — **VERIFIED** | Carried forward. A8's cited symbols re-checked against `main` this pass: `internal/resolve/resolve.go::Input` (`:224`, `Owned []Tag`), `::Tag` (`:20`, two fields), `::TagSet.has` (`:125`), `::missingOwned` (`:444`), `::RefusalKinds` (`:64`, closed five-kind set), `::escapeOrRefuse` (`:473`), `::TagSet.matches` (`:132`), `internal/resolve/adversarial_test.go::TestAdv1b_AllGuardsFalseIsNoMatchNotOwnedStateUnavailable` (`:114`). All resolve. |
| Named-but-unrun spikes | 3 | **none open** | One named spike; code and output captured (`evidence/spikes/main.go`, `output.txt` 23 lines). `{SPIKE_DIR}` did not need to grow this pass — no item's disposition required a new run. |
| Exactness-word delta (post-mutation) | 4 | **covered** | The rounds' new exactness claims — "exactly the keys it was asked for", "no unrequested key added", "every requested key", timeout precedence — all trace to A8 (Verified, `output.txt:11-13`) or to A9/A10's named MVV scenarios. No uncovered exactness word introduced. |
| Scope Verification narrower than the MVV | 4 (sweep miss) | **fixed in place** | §amendment-sweep residue: both iter-2 passes widened the MVV (gate denied, unsafe definition validation, `read_back_incomplete`, A9/A10 boundary cases) without propagating to the Finalization Gate's Scope Verification restatement. Propagated this pass. Not a finding either round claimed and missed — a site neither sweep listed. |
| A9 Evidence cited stale scenario numbers | 4 (sweep miss) | **fixed in place** | A9 pointed at "Scenario 2's absent-required-key case" and "its timeout-with-unresolved-keys case"; 3amigo T-6/T-8 moved both to Scenarios 6 and 7. Re-anchored. |
| Cross-RDR: RDR 0007 A6b routed here | 3amigo iter-2 cross-RDR note | **routed, not closed here** | This RDR discharges the obligation (read-completeness clause + seam-omission clause) and names 0007 in References. Closing A6b is RDR 0007's own edit; 0007 is itself `Draft [re-verify A3, A10, A15]`, so it will be re-run. Not a 0004 blocker. |
| Cross-RDR: `Row.RequiresOwned` has no obliged producer | critique iter-2 charted-to-successor | **routed to cluster-reconcile** | Recorded as an **Open** row in Capability Dependencies and cross-referenced under the Disposition Table. JDR 0001 §JD-3 already carries it as an open joint decision over 0002/0007/0009. Out of scope for 0004 — 0002 owns the normalized row. |

## Why A9 and A10 may be downgraded rather than verified

The Stage 6 hard rule bars deferring an assumption **the MVV depends on** — a
byte-parity reference, a normative fixture the MVV consumes, load-bearing
external behavior. A9 and A10 are the converse: they are properties the MVV
*proves*. Each of the five rules is asserted by a named scenario (1 arm 8, 3, 6,
7, 8), and none is an input any scenario consumes. `Method: MVV Test` is defined
as "pending implementation at lock time", so this is the state that Method
exists to express, not a deferral of the MVV itself.

Precedent in this repo: RDR 0006 locked to **Final** carrying A5 as
`Status: Pending` / `Method: MVV Test`, DOWNGRADED at its Stage 6 on the same
reasoning (the production-gate proof cannot run until the command exists).

Neither assumption is refuted. Nothing in the RDR treats A9's or A10's
properties as settled empirical fact — the clauses that state them are normative
requirements this RDR is legislating, which is what the MVV will check.

## Absorption audit

Delegated over `evidence/{grounding,3amigo,critique}/` for both iterations
(six rounds). Verdict: **ABSORBED — no residue.** Every finding dispositioned
`fixed` is carried in the current RDR text. Both iter-2 §amendment-sweep token
claims verified present: `(8 shapes)`, `read_back_incomplete` in the Failure
Modes list (11 sites), "artifact unavailable" removed and routed to
`execution_failure`, `A9`→`A9 and A10` at Prerequisites and Assumption
Verification, the "kernel cannot recover" phrasing corrected at both sites, the
bare silent-shape count replaced by a derivation rule, `typed absence marker`
gone. The audit surfaced the Scope Verification narrowing, dispositioned above.

## Completeness check

- No `_Draft placeholder._`, seed-skeleton header, or verbatim template bracket.
- `## References` is filled from citations the RDR carries (no bracketed
  placeholders).
- Every Critical Assumption is terminal: A1–A8 Verified, A9/A10 Pending with an
  explicit DOWNGRADED disposition, a named MVV plan, and a stated survivable
  "If wrong".
- Dispositions are written **in the RDR**, not only here.

## Punt ledger

No row owed. This pass routes nothing back. The Final→Draft demotion that opened
this re-entry came from applying JDR 0001's decisions across the 0002–0009
cluster (`76a9e67`), a deliberate upstream decision change, not an escaped
defect caught by a route-back.

## Verdict

**RECONCILED.** All items terminal, no BLOCKER, no refutation.
