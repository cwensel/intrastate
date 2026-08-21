Model: claude-opus-5[1m]

# 3amigo resolve — disposition ledger (iteration 2)

Origin ledger = the 35 persona findings (PM-1..10, IMP-1..11, QA-1..14),
consolidated into hotspots H-1..H-8 in `consolidation.md`. Every entry exits
exactly one way. Grounding gate run against code on `main` (working tree clean
at `6ec198c`), `{RDR_RESOURCES}` design docs, JDR 0001, and the RDR's own
decided text.

Delta scope: the re-entry `Status:` line (re-verify A3, A16, A17, A18, A19,
A21, A22) plus the passages the §D1/§D4 reshape touched. Two findings their
authors marked `[out-of-delta]` (IMP-11, QA-14) are dispositioned anyway —
both landed inside delta-scoped passages on inspection.

| # | Disposition | Origin | Section touched |
|---|---|---|---|
| IMP-2 / QA-1 / QA-7 | **fixed** | H-1 | Normative Contracts refusal clause — payload named: `Refusal.Undecided []UndecidedRow`, `UndecidedAtom`, closed `Reason` set; A23 added (Pending) |
| IMP-4 / QA-2 | **fixed** | H-2 | Normative Contracts refusal clause — sort key made TOTAL: `(RuleID, SourceLocator, key, block, operator, literal)`; propagated to 4 sites (§amendment-sweep) |
| PM-7 / IMP-5 / QA-14 | **fixed** | H-3 | Testing Strategy preamble — the 154-PASS baseline qualified: the spike proved a SUPERSET; Fixup-1d is re-decided, not re-encoded; Phase 1's exit count is the re-spike's |
| PM-8 / IMP-10 / QA-3 | **fixed** | H-4 | Testing Strategy row 15 — tautology closed: the kernel enforces no part of A21, so the test must build a conforming escape row explicitly; shipped `escapeRow` named as non-conforming on both fields |
| IMP-6 / QA-5 | **fixed** | H-6 | Normative Contracts existence clause — constants named (`OpExists`, `LiteralTrue`, `LiteralFalse`); malformed/missing literal ⇒ kernel-unevaluable `uncomparable`; load-rejection leg placed upstream; A24 added (Pending) |
| QA-8 | **fixed** | H-8 | `disposition` mini-check — `owned_state_unavailable` row corrected to "no code row yet (§JD-8)"; two exit-group notes added |
| PM-2 | **fixed** | H-8 | `disposition` mini-check note — `flow-guard-unevaluable` scoped to "supplied facts" mis-signals missing ARTIFACT state; recorded as RDR 0005's to fix |
| IMP-1 | **fixed** | — | Normative Contracts empty-block clause + K3 clause — the `¬(unless_conj)` TERM DROPS OUT; verdict reduces to `all_result` |
| QA-4 | **fixed** | — | Testing Strategy row 19 — discriminating leg named (`exists` over a PRESENT key under a nil seam ⇒ plan); absent-key leg marked non-discriminating |
| QA-10 | **fixed** | — | Testing Strategy row 16 — assert the exact disposition, not "a decided verdict"; "not `guard_unevaluable`" named as insufficient |
| QA-12 | **fixed** | — | `trace` step 5 — witness marked non-discriminating; ordering's real witness pointed at row 13 |
| IMP-9 / QA-9 / PM-10 | **fixed** | H-5 | Testing Strategy row 11 legs split (kernel vs load); row 8's blocked present-key half already stated in Phase 3 — cross-reference now accurate |
| PM-1 | **dismissed-with-cite** | H-8 | Re-raise: Prerequisites already adjudicate "Does not block lock" with a stated reversibility argument and a `Detail`-flattening fallback. The determinism tension PM-1 names is real but is A19's own recorded cost, not a new defect. |
| PM-6 | **dismissed-with-cite** | — | The closed reason set (`absent`/`uncomparable`) is deliberate; a third "producer defect" value would require the kernel to distinguish a typo from genuine absence, which it cannot (A7: the typo fails at load; A22: canonical keys are an input precondition). Failure Modes already states the misdiagnosis honestly. |
| PM-9 | **dismissed-with-cite** | — | Re-raise: GATE-THEN-COUNT adjudicates the ordering explicitly and states the two-round cost as "the honest cost," pinned to shipped behavior rather than to D8's superseded rationale. |
| IMP-3 | **dismissed-with-cite** | — | `GuardAtom`'s field types are JDR 0001 §D1's, not this RDR's — §D1 fixes the atom as "(key, operator token, literal, block)" and this RDR "cites it and does not restate the grammar" (SEAM clause). Set-shaped literals for `in`/`contains` are RDR 0003's typed semantics; Phase 4 already carries the element-encoding REQUEST. |
| IMP-7 | **dismissed-with-cite** | — | A22 already states both arms: the kernel "performs none" (documentation of an input precondition), and the fallback narrows the assumption, not the code. That the two produce identical kernel code is the point — the difference is where the obligation is filed, not what is built. No `Resolve`-entry check is implied. |
| IMP-11 | **dismissed-with-cite** | H-7 | The MVV's escape row is load-bearing as a NEGATIVE control: it is precisely what would rescue the row under the pre-RDR kernel, and `oracle` row 1 names "the pre-RDR kernel returns a plan here." Inert-under-the-new-ordering IS the assertion. |
| QA-11 | **dismissed-with-cite** | H-7 | Grounded and correct that `fixtures_test.go::escapeRow` is non-conforming — but that is now stated at row 15 (H-4 fix), which is the row that owns escape-row conformance. Fixing it twice would pin one fact in two places (the drift the resolve prompt's COMPUTE-DON'T-ARGUE clause warns about). |
| PM-3 | **charted-to-successor** | — | `Charted.md` — end-to-end "told plainly" acceptance belongs to RDR 0005; also the reason `xg7p` cannot fully close on this lock alone |
| PM-4 | **charted-to-successor** | — | `Charted.md` — veto dry-run/diagnostic mode is a CLI/lint surface (RDR 0005 / 0006); warrants `/rdr-seed` |
| PM-5 | **charted-to-successor** | — | `Charted.md` — lint gate for `unless`-over-rarely-set-tag is RDR 0003's grammar obligation |
| QA-6 | **charted-to-successor** | — | `Charted.md` — mutation-testing tooling is repo infrastructure, not this contract |
| IMP-8 / QA-13 | **charted-to-successor** | — | `Charted.md` — "callable regardless of sibling verdicts" belongs in Phase 3's exported seam contract test |

## Grounding gate notes

Verified firsthand against `main` before any edit:

- `fixtures_test.go::escapeRow` builds an escape row with
  `RequiresOwned: []string{"status"}` AND non-empty `Writes` — non-conforming
  on both fields A21/RDR 0009 require empty. QA-3's tautology claim is a fact,
  not a reading.
- The A3 spike's existence atom is
  `{Key: "guard.missing", Op: resolve.OpExists}` with no `Literal` field —
  it decides on a zero value, contradicting `presence == literal`. QA-5
  confirmed; this is what A24 now covers.
- RDR 0005's stable-code table (lines 603–614) has rows for
  `flow-guard-unevaluable` and `flow-fact-missing` but **none** for
  `owned_state_unavailable`; JDR 0001 §JD-8 names that exact gap. QA-8
  confirmed — the mini-check was asserting a mapping that does not exist.
- `clierr.go` GroupUserEnv and GroupInternal both exit 2; `GroupEnvUnavailable`
  exits 3. PM-2's "wrong remedy" reading is grounded, and §JD-8 already flags
  the exit-3 question as live.
- `resolve.go::gate` — the `slices.MinFunc` representative-row pick, the
  prune → owned → undecidable partition, and `evaluateGuard`'s whole-guard
  `seam == nil` branch are all exactly as the draft describes them.
- JDR 0001 §D4 does route the payload "to the CLI through §JD-8's `Detail`",
  and §JD-12 records §D4 as Closed — A19's "this is a reversal" framing is
  correct as written.
- RDR 0002's Refinement Context Direction list names only the `RequiresOwned`
  producer (JD-3) and the fixture rename (JD-10) — neither the existence-token
  duty (A16) nor canonicalization (A22). The Prerequisite is accurate.

Re-raise check: PM-1, PM-9, and IMP-7 were each tested against the RDR's own
decided text and dismissed as re-litigation of adjudicated calls. In every case
the decision was already explicit in the draft, so no "make the implicit
explicit" edit was owed — unlike iteration 1's PM-4, where it was.

COMPUTE-DON'T-ARGUE: the H-2 fix states the executed comparator tuple rather
than arguing that the old key was "determinate enough"; the H-3 fix states what
the spike actually ran (a superset retaining `Refusal.Guard`) rather than
restating the 154-PASS claim.

## Needs (re)verification — carried to Stage 6

| ID | State | Why |
| --- | --- | --- |
| A23 | **new, Pending** | The named payload surface (`Refusal.Undecided []UndecidedRow`, `UndecidedAtom`, closed `Reason` set) is new normative signature. Folds into A3's re-spike. Check no frozen test constrains `Refusal`'s field set. |
| A24 | **new, Pending** | Verbatim existence-literal comparison against `LiteralTrue`/`LiteralFalse`, and malformed ⇒ `uncomparable`. Kernel half is this RDR's to decide; producer half rides the same unlanded 0002 duty as A16/A22. |
| A3 | **Pending** (unchanged) | Unchanged by this lens, but its re-spike now closes three assumptions (A3, A23, A24) — the reshape with `Refusal.Guard` actually removed is where all three land. |
| A16 | Verified (kernel half) | Kernel half now names the three constants explicitly; producer half still acknowledged-not-clause. No status flip. |
| A19 | **Pending** (unchanged) | Untouched. PM-1's challenge dismissed as a re-raise; the `Detail`-vs-structured-field reversal still needs the JDR reopening. |
| A21 | Verified (composition) | Unchanged, but row 15 now states that the kernel enforces no part of it — the composition is a producer property, and the test says so. |
| A22 | **Pending** (unchanged) | IMP-7 dismissed; both arms were already stated. |
| — | **new claim** | The total sort tuple is now a normative ordering contract. Row 17 must exercise the tie case (two atoms, one key, permuted) — verify at Phase 2 that the comparator is total in code, not just on paper. |
| — | **new claim** | Empty-`unless` term-drop-out. Verify at Phase 1 that the implementation reduces the verdict to `all_result` rather than substituting a truth value for `unless_conj`. |

## Mini-checks

No new cue fired. The four cues already fired on this draft (authority, oracle,
disposition, trace) have their tables in the RDR from the cove pass; this pass's
fixes edited three of those tables in place (`disposition` row + exit-group
notes, `trace` step 5) rather than adding a cue. Per the skill's inheritance
rule, the cue read is not re-run.

## Tiebreakers

None escalated. The one fork that could have gone either way — whether to NAME
the payload surface here or defer it to Phase 2 as an implementation choice —
collapsed on evidence without a §strong-consult: six assertion sites (matrix
rows 1, 12, 17, 20; MVV scenario 1; `oracle` row 1) plus RDR 0005's renderer and
Phase 3's contract-test export all type against a field that had no name, and
the RDR's own Load-Bearing Decisions already committed to "a field on `Refusal`,
not a new kind." Naming it is executing a decision the draft had made, not
making a new one. The residual risk — that a frozen test constrains `Refusal`'s
shape — is booked as A23 rather than absorbed.
