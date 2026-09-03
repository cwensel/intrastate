Model: claude-opus-5

# Stage 6 — Reconcile: 0026 (bounded output drain on command timeout)

Date: 2026-09-03 · Profile: foundational · Verdict: **RECONCILED**

## Stage 5 preflight

All four routing calls clear, so Stage 5 is genuinely complete:

| call | emit.next | why |
| --- | --- | --- |
| `--outcome lens` | `/rdr-reconcile` | foundational row cove → 3amigo → critique → repeatability complete |
| `--outcome critique` | `none` | two base models, diffed — dual-model obligation discharged |
| `--outcome repeatability` | `none` | variant full — run-1/2/3 and the diff written |
| `--outcome determinacy` | `none` | `Determinacy: fired`; the full lens is already on the row |

No caveat to carry: no single-model fallback, no unstamped pass, no variant mismatch.

## The open set — four sources

1. **Pre-Lock needs-verification list** — A9, A10, A8 (the Stage-5 gate's own three).
2. **Still-Pending Critical Assumptions** — A9, A10, A11, A12, A8. Five, not three: the
   gate text predates A11 (critique D-8) and A12 (repeatability G-6).
3. **`spikes_unrun`** — empty. Corroborated by an absorption audit over all four lens
   dirs (3amigo, critique, repeatability, cove), both iterations each: every finding
   traces to an absorbed record change or an explicitly charted dismissal with a named
   successor. No unabsorbed residue implying a spike.
4. **Exactness-word delta** — `rdr lint 0026 | grep prose:exactness` returns zero
   findings over what the rounds touched.

## Dispositions

| item | src | disposition | evidence / plan |
| --- | --- | --- | --- |
| A9 — per-read re-arm bounds the idle gap, not the tail | 1,2 | **VERIFIED**, witness corrected | `evidence/spikes/a9-a10-regrace/`, both OSes. Paced 1 MiB: 12/12 whole on EOF, max per-read gap 26.8 ms (darwin) / 31.6 ms (linux) under the 50 ms grace while total elapsed ran 12.2x / 13.3x the grace |
| A10 — pollability probe is safe and detecting | 1,2 | **VERIFIED** | same dir. `os.ErrNoDeadline` on both non-pollable shapes incl. a genuine dup'd pipe fd; parked-read 12/12 at 300 ms against a 10 ms probe window; round-trip identical in every cell |
| A8 — executor derives `Applied()` from the typed held-pipe error | 1,2 | **VERIFIED** | all 11 sub-claims hold, zero line drift. Decisive: `flow_exec.go:410-412` `ClassExecutionFailure` arm has no phase check and no `Detail`; `Applied()` has zero non-test callers (23 hits, all `_test.go`) |
| A12 — the stall seam is new; a nil-in-production hook suffices | 2 | **VERIFIED** | only injection var is `goos` (`cmdbind.go:97`); timing tests drive real `sleep`; hook inserts at `cmdbind.go:273`/`:277`; only build tags are the unix/other split |
| A11 — the bounded join preserves the happens-before edge | 2 | **DOWNGRADED** | Unverifiable pre-implementation by construction: its check is S3+S7 under `-race`, and S7 needs the seam A12 confirms this record *adds*. Plan: Phase 1 closes only on S3+S7 passing `-race -count=25`; a race report there is a Phase-1 blocker |

No BLOCKER. No refutation of any assumption the design relies on.

## The one finding that needed adjudication

The A9 spike confirmed the **mechanism** and refuted the **witness**.

MVV row 6 (marked Normative) specified a child writing 1 MiB in one `Write` and exiting,
claiming to be "what distinguishes the two candidate mechanisms." Measured, it does not: a
buffered 1 MiB tail drains in 1.6 ms (darwin) / 9.5 ms (linux) across 32 non-waiting reads,
so the correct re-armed mechanism AND the rejected single-absolute deadline both recover it
whole. A wrongly-implemented drain would have PASSED row 6.

The discriminating shape is a **paced** child (32 chunks, 20 ms idle gaps = 0.4x grace each,
summing to ~12x grace): re-armed 12/12 complete on EOF, single-absolute **0/12**, truncating
at exactly 163840 of 1048576 bytes (15.6%) after 5 reads at ~51 ms with terminal
`os.ErrDeadlineExceeded` — a false held report on a live writer, identical on both OSes.
Even a 1 KiB paced payload loses 62%. **Duration, not size, is the operative variable.**

This mattered because three sites lean on row 6 as the catcher — C1 `precedence:`, and two
Illustrative Code notes recording that an earlier reconstruction *did* rebuild the one-shot
parent deadline. The record's own `oracle` mini-check table had flagged the question
("its 1 MiB tail under a re-armed grace is A9, unmeasured until that spike runs").

**Strong-consult (fresh context, factored brief): PASS** — contract correct and vindicated
by measurement; a test-specification defect, not a design defect; the under-powered witness
is repaired by strengthening the witness, not by reopening the approach. No reading found
in which the paced shape indicts the mechanism itself.

**User tiebreaker: amend in place.** Applied:

- MVV row 6 → two shapes: 6(a) burst (guards against a deadline of literally `now`) and
  6(b) paced (the mechanism discriminator), with 6(a) explicitly marked insufficient alone.
- C1 `precedence:` → "never the size **or total duration** of the tail"; the rationale's
  "a large buffered payload" → "a slow-arriving tail", since a large *buffered* payload
  demonstrably does not trip the single-deadline variant.

### §amendment-sweep

Pre-edit tokens grepped and every surviving site updated in the same pass:

| site | change |
| --- | --- |
| C1 `precedence:` (2 sentences) | the two precision edits above |
| Illustrative Code comment (`~:354`) | "size or total duration", with A9's gap-vs-total numbers |
| Illustrative Code prose (2 sites) | both now name row **6(b)**; the reconstruction note records that 6(a) would have let it pass |
| `oracle` mini-check table, row 6 | now states both shapes and which one discriminates; the "unmeasured" note retired |
| Gate §Assumption Verification | Stage 6 disposition appended; the contingency it declared is discharged |

Consumers re-read for semantic agreement, not just staleness: C1 `whole:` and the
`fidelity` table's "grace bounds the idle gap between reads, not the tail" row already
state the rule in gap terms and needed no change.

## Completeness check

- No `_Draft placeholder._` and no seed-skeleton header survives.
- `## References` carries no template placeholders; the new spike dir was added.
- `rdr lint 0026`: **PASS**, blocking=0 resolution=0. The 5 `placeholder:survived`
  findings are all Finalization Gate guidance blocks — Stage 7's section, unchanged in
  count from this stage's baseline lint.
- Post-edit tags: `ca_total=12 ca_verified=11 ca_pending=1 ca_unverified=0
  ca_off_vocabulary=0`, `ca_pending_ids=["0026:A11"]`, `spikes_unrun=[]`.

## Known sensitivity recorded (not blocking)

The measured per-read gap (27-32 ms) sits under 2x `DrainGrace` (50 ms), so a heavily
loaded host or slower writer cadence could push a legitimate gap past the grace. Recorded
on A9: a tuning question about the constant's VALUE, which argues for a comfortable grace
and never for the single-absolute mechanism. Revisit if S7/S9 flake under CI load.

## Verdict

**RECONCILED.** Every item terminal — four VERIFIED, one DOWNGRADED with a binding
pre-phase plan, no BLOCKER. No MVV-critical item deferred: A9 and A10, the two the MVV and
C1 leaned on, are both verified by spike on both OSes; A11 is not MVV-reachable and is
bounded by a Phase-1 gate. Ready for Finalize.
