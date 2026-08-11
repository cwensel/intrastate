# Critique — resolution (RDR 0009, iteration 1)

Model: claude-fable-5 (resolve session)

Origin ledger: merged M-1…M-15 in `diff.md` (anchor-matched from
`critique.md` A-C-* and `critique-modelB.md` B-C-*). Grounding reads:
`internal/resolve/resolve.go` (whole file), `internal/cli/clierr/clierr.go`
(whole file), `internal/cli/root.go` (whole file),
`internal/resolve/resolve_test.go::TestReq1_DispositionIsExactlyOneOfPlanOrRefusal`.

## Dispositions

- **M-8 (A-C-2) — fixed.** Confirmed against `internal/cli/root.go::ExecuteAndEmit`:
  every non-`CLIError` returned by a verb is converted via `cobraErrorToCLIError`
  (`GroupUserEnv`, `Code: "command-error"`, generic `--help` hint) → exit **2**,
  not 1. The RDR's exit-1 tripwire was false; corrected in A2's evidence, the
  Normative CLI clause, the Prerequisites flow-verb bullet, and scenario 11 —
  the missing-wrap tripwire is now the stable `Code`, never the exit code.
- **M-1 (A-C-3 ≡ B-C-1) — fixed.** `resolve.go::Result` doc ("never both and
  never neither") and `Resolve`'s "returns exactly one disposition" contradict
  the zero-`Result`-plus-error return. Extended the *Doc-contract amendment*
  decision to authorize scoping both doc contracts to nil-error returns.
  Frozen `TestReq1_DispositionIsExactlyOneOfPlanOrRefusal` unaffected: inputs
  are conforming tables; `mustResolve` checks the error first (scenario 7b
  already pins breach-side behavior).
- **M-7 (A-C-1) — fixed.** Phase 1 landing before Phase 2 turns the frozen
  suite red (entry check errors on the sixteen `escapeRow`-fed fixtures;
  `mustResolve` fatals). Added the one-change landing requirement to Phase 1.
- **M-9 (A-C-6) — fixed.** Added to the predicate-export clause: `CheckValid`'s
  doc comment must name the one property checked and state that nothing else is
  checked (in particular not the Escape-class restriction Row's doc records),
  so nil is never read as general validity.
- **M-6 (A-C-7 ≡ B-C-6) — fixed.** CLI clause now binds the carrier *class*
  only, delegating field name / clierr-local type / `Count` serialization to
  A9; A9's plan gains item (d): the field must keep `clierr` a leaf (no
  `internal/resolve` import — clierr-local representation, not
  `resolve.RowRef`), and settles whether `Count` serializes.
- **M-15 (B-C-12) — fixed.** Scenario 5's baseline made regenerable: re-run the
  frozen suite at the implementation's base commit before Phase 1; the spike
  artifact records an instance of the oracle, not its definition.
- **M-3 (A-C-5 ≡ B-C-10) — dismissed-with-cite + one-line sharpening.** The
  peer dependency is adjudicated (A4 *Carried forward*; Prerequisites binding
  obligation on RDR 0004's implement stage). The narrow "checkbox has no test
  form" sub-point fixed: the Prerequisite now names the fixture 0004's
  implement copies (escaped plan, empty `Writes`, non-empty `NextTags` → zero
  owned-tag mutations).
- **M-2 (A-C-4 ≡ B-C-5) — dismissed-with-cite + clarity edit.** Collapse+Count
  is adjudicated (A8; multi-breach clause records the rejected per-row-entries
  branch and the degenerate-producer case). A's "REQ-2/REQ-10 bind only modeled
  payloads" is correct — made explicit in the clause that the extension to the
  breach report is a deliberate adoption, not an inherited obligation.
- **M-5 (A-C-9 ≡ B-C-8) — dismissed-with-cite.** The surface pin is the
  recorded decision ("The surface is fixed here, not left to the implementer,
  because it lands permanently on RDR 0001's locked package surface"); the
  type-assertion hazard is exactly what the errors.Is/As requirement and the
  Failure Modes traversal guidance exist to prevent; a future second
  entry-time error source renegotiates the contract by definition (a successor
  decision, recorded then).
- **M-10 (B-C-2) — dismissed-with-cite.** Symptom is false: scenarios 1–4, 7b,
  9, 10, 10b become permanent suite members at implementation, so a refactor
  dropping `CheckValid` fails them. They are excluded only from scenario 6's
  mutation *oracle*, to keep the non-vacuity proof non-tautological — which is
  what the Mutant A paragraph already records.
- **M-11 (B-C-4) — dismissed-with-cite.** Keeping the ADV-2b/Fixup-1e override
  values is unimplementable under the chosen design: a write-bearing escape
  fixture errors at `Resolve` entry and `mustResolve` fatals, so the suite
  cannot run at all. The discriminating assertions are on `Refusal.Kind`; the
  `Writes` reads sit only in `t.Fatalf` format args (A3 evidence).
- **M-12 (B-C-7) — dismissed-with-cite.** No conflation in the RDR: the
  "cheapness preference" bullet governs placement relative to `assemble`
  (pure, disposition-free), not relative to the alphabet check; precedence
  over `unmodeled_outcome` is separately normative, deliberate, and covered by
  scenario 4 plus the Trade-offs dormant-row bullet.
- **M-4 (A-C-8 ≡ B-C-3) — dismissed-with-cite.** The staging is adjudicated:
  Trade-offs "Staged, not immediate," Risk #1 mitigation, and the Prerequisite
  binding RDR 0002's implement stage to the Phase 3 fixtures. Re-raise of a
  recorded decision.
- **M-14 (B-C-11) — dismissed-with-cite.** The dormant-row behavior change is
  recorded as deliberate (Trade-offs negative bullet; MVV pins it), A5's
  expiry timing is weighed in the author-friction bullet, and the exported
  predicate is the migration path for future producers (A5's "If wrong" names
  the migration step).
- **M-13 (B-C-9) — dismissed-with-cite.** An unfilled Finalization Gate is
  Draft-normal (Stage 7 fills it); the Proportionality gate item re-counts
  contracts and re-validates Profile at finalize — that is where the split
  decision is owned. Not absorbed here.

## Needs (re)verification (carried to Stage 6)

- **A2** — evidence corrected at this pass (unwrapped verb error exits 2 via
  `ExecuteAndEmit`→`cobraErrorToCLIError`, not 1; grounded against
  `internal/cli/root.go` in this session). Stage 6 re-confirms the corrected
  wiring claim; A2's verdict (exit-2 routing, no RDR 0005 change) stands.
- **A9** — still Pending (owed to Stage 6); plan widened with item (d):
  clierr-leaf coupling constraint and `Count` serialization.

## Tiebreakers

None — every fork collapsed with grounding evidence.
