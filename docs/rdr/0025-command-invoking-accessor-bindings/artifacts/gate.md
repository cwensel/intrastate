# Finalization Gate — cli/0025 Command-Invoking Accessor Bindings

Date: 2026-08-30
Verdict: **READY** — Gate PASS, locked to Final.

Mechanical pre-sweep (tooling-pass, iter 1): PASS. `lint --locking` exit 0,
`blocking=0 resolution=0 placeholder=5 advisory=10`. The five
`placeholder:survived` findings and `gate:inline` are this section's own
unwritten template text, discharged by writing this file and applying the lock
edit. All 81 `source-anchor` edges resolve `true` (C5); no `Docs Only` record
(C4); no off-vocabulary Method member (C2); no Source-Search self-reference
(C3); no `Pending`/`Unverified` assumption (C6). Report:
`evidence/tooling-pass/tooling-pass.md`.

## 1. Contradiction Check

No contradictions found between research findings, design principles, and
proposed solution.

The Key Discoveries and the Approach agree on the deciding axis. The reversal
ledger (helm, consul, hugo, Go `ErrDot`/`GOVCS`; 32 peer state-machine tools)
records movement in exactly one direction — string→argv, inherit→allowlist,
open→gated — and C1/C2 (fixed argv, closed placeholder vocabulary), C4's env
allowlist and C6's execution gate each sit on the settled side. No finding
pushes the other way.

The one place research pressed *against* the design is recorded rather than
smoothed: the A3 spike establishes that established-tool **writes** cannot bind
directly (a fixed argv cannot carry the planned value; the tools ignore stdin),
which partly concedes the "no wrappers" ambition. The Approach states this
plainly ("integration honesty", premortem P-1/P-8/P-9), the Decision Rationale
scores the row on the read/gate half only and says so, and Consequences carries
it as a named negative. That is a disclosed limit, not a conflict.

Likewise the per-invocation `--allow-commands` gate is the one shape none of
the ledger's peers *kept* — Risks names this against the RDR's own evidence,
states no v1 mitigation exists (A10: no config surface), and argues only the
ordering (ship gated, add a second grant path later). Contradiction between
principle and plan would be claiming the flag is sufficient; the RDR claims it
is the cheap direction and books the cost.

`0025:C1`'s relaxation of the load-time path-required rule was checked against
the locked peers (0001–0011) in the joint-check Arm 3 sweep: no locked peer
relies on that refusal, and 0004's Alt-2/Alt-3 rejections (shell, global
allowlist) are both preserved by C2/C5 rather than re-litigated.

## 2. Assumption Verification

All 15 Critical Assumptions terminal and `Verified`; `ca=all-terminal`,
`ca_pending=0`, `ca_unverified=0`, `ca_placeholder=0`, `ca_off_vocabulary=0`.

Method distribution, all within the sanctioned vocabulary: Spike ×4
(A1, A3, A7, A8), Source Search ×9 (A2, A5, A9, A10, A11, A12, A13, A14, A15),
Peer RDR ×2 (A4, A6). No `Docs Only` record exists, so the load-bearing
`Docs Only` blocker cannot apply. Each record carries a non-empty "If wrong",
and Status/Method/Evidence agree on every row.

No self-referential `Verified` stamp: the 45 distinct anchor targets are all
`internal/…` source symbols or spike outputs; none resolves to this RDR or its
artifact directory (CHECK 3 clean — and this is the check that signals a
disturbed evidence record, so its silence here is meaningful).

Every cited `path::Symbol` resolves on `main`: 81 of 81 `source-anchor` edges
`resolved: true`, none `false`, none absent. The two `Peer RDR` records cite
elements, not bare records (A4 → `0016:C4`; A6 → 0004's element ids), so
`peer-evidence:no-element` does not fire.

Status consistency holds trivially in the direction that matters — nothing is
`Pending`, so no settled-fact prose depends on an unverified property. The two
assumptions whose subject matter is genuinely *not yet built* are honest about
it in the other direction: A11 (`accessor.ExecError` does not exist on `main`)
and A10 (no config subsystem) are `Verified` statements **of absence**, both
confirmed by source search, and the Capability Dependencies table marks the
corresponding capabilities `Not built` with C4/C6 naming what this RDR
specifies. That is a verified absence carried into a spec obligation, not a
fact assumed into existence.

A15 was refuted and repaired in-gate during reconcile (the `exit_absent`
channel), not routed back; A11–A15 were re-verified after that repair.

## 3. Scope Verification

The Minimum Viable Validation is in scope and executed during implementation,
not deferred. Phase 4 ("Surface and proof") carries it, and the Testing
Strategy pins it to a named spec: **S6 — the MVV end to end against an
established tool present in CI**, supported by S1 (load-time C1/C5 validation
plus six rejected mutants), S3 (envelope-before-exit ordering), S4 (deadline,
≈1.0s with 0 surviving processes), S5 (write read-back) and S7 (the C6
execution gate).

The specific proof: a linted model declaring one command-backed read and one
command-backed write over a single artifact role, delegating to `git config`;
`intrastate lint` accepts it and rejects six mutated copies each with its own
C5 defect; with `--allow-commands` the declared write applies a real state
change through `internal/cli/flow_state.go`'s existing set-state path (and
without the flag the same invocation refuses `execution_failure` before spawn);
read-back through the declared command reader verifies the written value; and a
command sleeping past its declared timeout refuses `timeout` leaving no orphan
process.

The MVV is honest about its own boundary: automatic resolve→apply wiring is
explicitly *not* in scope (stated in the Problem Statement and in the MVV
itself, and no Phase builds it). What the MVV proves is that a declared command
can carry the apply — the half that is missing today. Scope is bounded, not
quietly shrunk: the end state still requires a real state change applied and
verified through a declared command.

## 5. Proportionality

Right-sized. Nothing to trim before locking.

**Contract count — the split test.** Seven labelled contracts (C1–C7), but one
independent load-bearing seam: *what carries a declared accessor's real
external behaviour, and what bounds the authority it executes with*. C1
(carrier, exactly-one) and C5 (the load-time defects that make C1 checkable)
are the carrier; C2 (closed placeholder vocabulary) and C3 (stdin/stdout
envelope and exit maps) are the invocation surface; C4 (env allowlist, deadline
triple, envelope-before-exit ordering, `ExecError`) is the authority bound; C6
is the gate that governs whether any of it executes; C7 quotes the implemented
seam. Removing any one leaves the others incoherent — C1 without C5 is
unenforceable, C2 without C4 is an unbounded hole, C6 without C1 has nothing to
gate. This is one seam expressed at the granularity lint can check, not several
seams locked together. The Profile field states this and the Resolve recount
agrees: "one seam, the declared-command carrier + authority bound (C1/C2/C5/C6)
with its invocation envelope (C3/C4), inseparable — no split."

**Profile re-validation.** `foundational` is correct and matches the lenses
that actually ran. Seam Lineage ≥2 floors it there independently
(`internal/accessor/binding.go::WriteBinding` — 3rd point-fix; trail dkcf,
742a). The full foundational row is complete and not stale: cove → 3amigo →
critique → repeatability, with `critique_models=differ` (the dual-model
obligation genuinely discharged, not a same-model fallback) and
`repeatability_variant=full` (run-1/2/3 plus diff). `lens_stale=none`. The
determinacy obligation is therefore satisfied by a full repeatability lens, not
an `n/a` disposition. No lens is owed; the Profile did not route past anything.

Form is correct: value plus one clause naming the contracts, no template or
Seed matrix prose left in the field.

**Document size.** 1917 lines is large, and the mass is where it should be:
Normative Contracts 34.8K and Critical Assumptions 30.7K — the two sections a
foundational RDR exists to carry — against 1.8K of Problem Statement and 2.4K
of Approach. The Testing Strategy (11.3K, S1–S7) and Pre-Lock Mini-Checks
(4.8K) are verification content the implementation prompt consumes directly.

**C9 evidence-field budget (advisory, 4 hits: A11 38 lines, A12 31, A13 38,
A15 62).** Answering the question C9 puts to the author: the load-bearing
anchor is still findable in each — A11 anchors `executor.go::refusalOf` /
`refusalWithKeys` / `readOutcome.classify`, A12 `flow_input.go::parseArtifacts`
and `selectModel`, A13 `flow.go::newFlowCmd` and `flow_exec.go`, A15
`binding.go::ReadBinding.Read` and `model.go::ReadResult.OwnedSnapshot` (7
anchors). The balance is not padding: these four are precisely the assumptions
that establish **absence** on `main` (A11's missing `ExecError`, A13's
single-construction-site claim) or map a specified envelope onto an existing
signature (A15), where the verification is the enumeration and a pointer would
discard what the grounding sweep reads. A15's 62 lines carry the `exit_absent`
repair this record's reconcile produced. Keeping them in the field, per C9's
own instruction not to truncate.
