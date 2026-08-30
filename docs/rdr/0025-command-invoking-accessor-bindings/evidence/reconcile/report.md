Model: claude-opus-5[1m]

# Stage 6 Reconcile — RDR 0025

Profile `foundational`. Preflight: `--outcome lens` = `/rdr-reconcile`
(cove -> 3amigo -> critique -> repeatability row complete); `--outcome critique`
= none (two base models, diffed); `--outcome repeatability` = none (variant
full, run-1/2/3 + diff). No caveats. Determinacy judgement not owed at
`foundational`.

## Open set (four sources)

1. **Pre-Lock needs-verification lists** — cove's carried block: A10 (new,
   Pending), A1/A4 (still Verified, precision-only edits). The 3amigo,
   critique and repeatability rounds resolved in-pass per the prelock contract.
2. **Still-Pending assumptions** — A11, A12, A13, A14, A15 (`ca_pending_ids`).
3. **Named-but-unrun spikes** — none. Both named spikes (`a1-deadline`,
   `a3-a7-a8-real-tools`) have `run.sh` + `output.txt` captured.
4. **Exactness-word delta** — `rdr lint 0025 | grep prose:exactness` returns
   ZERO findings. Baseline lint `blocking=0 resolution=0`.

Absorption audit (delegated, all four lens dirs): every finding absorbed or
dismissed-with-cite; **no residue**. critique D-1..D-15, repeatability
D1-D8/G1-G6, cove CV-1..CV-12, 3amigo iter-3's four FAILs all land at cited
elements. Dismissed-with-cite: critique D-11/D-13 (charted to the successor
config RDR), cove CV-7 (`verdictFor` disjoint by carrier), the
gate-sub-sections finding (Stage 7's obligation, correct state for a Draft).

## Dispositions

| item | source | disposition | evidence / plan |
|---|---|---|---|
| A11 | 2 | **VERIFIED** | `refusalOf` (16 sites) + `refusalWithKeys` (2, both in `Executor.Read`); every site can supply the error or `nil`; `selects`'s pre-selection literal bypasses both and structurally cannot carry one. No consumer branches on the discarded error (`disposition.go` reads `Class`; `flow_exec.go::accessorFailure` reads `Class`/`Keys`/`Expected`/`Observed`). Three net-new facts recorded: `Refusal.Detail`, `accessor.ExecError`, `readOutcome.err` |
| A12 | 2 | **VERIFIED** | Stronger than claimed: `buildRequest` obtains the path itself (`model, ref, ce := selectModel(cmd)`) and builds the registry three statements later; `flowRequest.modelRef` already carries it. No `table.Model` field, no caller change, no import cycle (`flowbind` already imports `os`/`path/filepath`; `internal/table` is pinned not to import `internal/cli`) |
| A13 | 2 | **VERIFIED** (one arm corrected) | One production `Registry` call site; four `buildRequest` callers; no other path constructs an `accessor.Registry` or reaches a binding. CORRECTED: pflag returns `(false, *NotExistError)` — "flag accessed but not defined" — not a bare false. C6 no longer rests on the discard idiom: ONE registration on the `flow` group's `PersistentFlags()` (subtree = exactly the four call sites; `lint` at root outside it) + lookup error CHECKED at `buildRequest` |
| A14 | 2 | **VERIFIED** | Zero `Path: ""` literals in `internal/`, production or test. `accessorTable` refuses empty declared paths pre-model; `validate.go::Validate` has no empty-path arm, so the loader is the sole existing refusal — which is what makes the constructor-side refusal load-bearing. First-write flows carry a non-empty DECLARED path; `load` operates on the caller's `art.Path`, a different field |
| A15 | 2 | **VERIFIED after in-gate repair** | Arm 3 was **REFUTED**: omission from both slices is UNREADABLE, not absent (`readOutcome.classify` `!ok` -> `unread` -> `ClassIncompleteRead`; pinned by `TestReq29_UnclassifiedKeyDefaultsToUnreadableNotAbsent`). Repaired to `KeyValue{Absent: true}` in `values`. `Invocations()` has zero production consumers; "count `Apply` entries, no internal respawn" is sound and is now normative in C7 |
| A10 | 1 | **VERIFIED** (unchanged) | Carried from cove as new-Pending; already dispositioned Verified at the cove pass with the absence established in source. No re-verification owed |
| A1, A4 | 1 | **VERIFIED** (unchanged) | cove edits were precision-only; no Verified assumption was invalidated |
| spikes | 3 | **none open** | Both captured |

`ca_total=15 ca_verified=15 ca_pending=0 ca_unverified=0`.

## The one refutation, and why it was repaired rather than routed back

A15's arm 3 met the hard rule's trigger. Handling:

- **Consult run before any drafting** (rdr-common §strong-consult). Verdict
  **PASS**: the repair is correct (`clearIsUnreadable` is guarded by
  `!v.Absent`, so an Absent record cannot be reclassified), keeps every
  contract total, touches no path-backed binding and no existing test.
- **A conflicting consult was adjudicated against source.** A second agent,
  reading 0004's RECORD TEXT, recommended the OPPOSITE repair — route absence
  through omission, citing `0004:C8` "what crosses the seam MUST leave the key
  absent from the owned snapshot". Source shows a **two-boundary** design and
  C8 governs only the outer one: `ReadBinding.Read`'s doc says it returns "one
  KeyValue per key it could read (present or **established-absent**)", and
  `KeyValue.Absent`'s doc says "This is the accessor layer's INTERNAL
  representation; what crosses the seam omits the key entirely (`0004:C8`)".
  `model.go::ReadResult.OwnedSnapshot` is the code performing the flag->omission
  conversion, citing C8 by name. Following that recommendation would have
  shipped the bug as the fix. Adjudication: `two-boundary-absence.md`.
- **Route-back not owed.** The corpus criterion (0008's A12 hard rule, and
  0002 M-13 / 0024 A8) licenses in-gate repair when the refutation hits
  evidence or enforcement reach and no normative rule, interface, or test
  changes. Here no interface changes, no path-backed binding changes, and no
  0004 test changes — `TestReq29` keeps pinning the defensive default; the
  repair stops RELYING on it. Repairs outnumber route-backs ~4:1 at this gate
  in-corpus. No §punt-ledger row is owed: no completed stage was reopened.
- **User accepted the repair** as the tiebreaker before this verdict.

## Edits landed in the RDR

- **C7** — `exit_absent` absence crosses as `KeyValue{Absent: true}` in
  `values`; omission named as UNREADABLE with the `classify` anchor; the
  two-boundary split stated. Plus the `Invocations()` no-respawn constraint,
  promoted from aside to normative.
- **A15** — header carries the corrected arm AND C3's whole-invocation
  qualifier; Evidence records the refutation, the repair, the two boundaries,
  the granularity restriction, and the prior art. Status -> Verified.
- **C3 prose** — the positive-declaration paragraph now names the return
  channel, not just the wire.
- **S3** — case 6's Expected asserts the channel, not just the verdict word.
- **disposition mini-check** — the `exit_absent` row names the channel.
- **A11** — header marks `ExecError`/`Detail` net-new; Evidence records the
  three net-new facts and the verification. Status -> Verified.
- **Capability Dependencies** — two rows added: `accessor.ExecError`
  (Not built), `readOutcome.err` (Extend).
- **C6** — `registration:` replaced: ONE registration on the `flow` group;
  new `lookup:` line requiring the error be CHECKED, with the fail-open
  hazard and the prior art named.
- **A13** — header and Evidence carry the corrected arm. Status -> Verified.
- **C6 rationale prose** — the stale "unregistered verb's lookup returns
  false" reasoning replaced.
- **S7** — oracle shifted from per-verb registration to group containment
  (in-group verbs resolve the flag, `lint` does not, an out-of-group verb
  fails), still deriving both sets rather than restating a literal four.
- **Phase 3** — names the group registration and the checked lookup.
- **A12, A14** — Evidence records the verification. Status -> Verified.

§amendment-sweep run for both amendments: pre-edit tokens swept
("neither slice", "established-absent", `exit_absent`, "returns false",
"registers the flag", `Invocations`, `Absent`, `allow_commands`); one stale
site found and fixed (the C6 rationale prose at the Decision Rationale), one
confirmed unaffected (the `authority` mini-check row, which is about where
permission is decided, not where the flag registers).

## Completeness check

No `_Draft placeholder._`, no `this is a seed skeleton` header. `## References`
fully authored — no template brackets. Closing lint: `blocking=0 resolution=0`,
placeholder=5 (all in the Finalization Gate, Stage 7's section). Advisory rose
6 -> 10: four `evidence:over-budget` on the Evidence fields expanded here; the
lint's own fix text says "do NOT truncate: the mass is usually verification
content the grounding sweep reads", and the budget call is Stage 7's.

## Verdict

**RECONCILED** — all items terminal, no BLOCKER. Ready for Finalize.
