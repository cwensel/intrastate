Model: claude-sonnet-5

# Diff — RDR 0007 critique passes, iteration 2

Pass A: `critique.md`, model claude-opus-5, ledger C-1..C-14.
Pass B: `critique-modelB.md`, model claude-fable-5, ledger B-1..B-13.
Reconciled by RDR passage anchor first, then by prose claim (ledger IDs
do not correspond across files).

## Convergent

Both models independently drew the same defect from the same passage.
This is the strongest signal in the set — two independent hostile
passes over the same document landed on the same anchor without
coordination.

| A row | B row | Shared passage anchor | Shared claim |
|---|---|---|---|
| C-1 | B-2 | Normative existence clause `presence == literal` vs. A24's Evidence / A3 spike's literal-free existence rule | The only executed prototype decides `exists` from presence alone (ignores `Literal`), the opposite rule from the normative clause; Phase 1 is handed that spike as "the shape to start from." |
| C-6 | B-2 | A3 Status: "Pending — the SHAPE the spike proved is not the shape this RDR specifies" | The RDR's sole executable evidence (154 PASS, `escapeOrRefuse` zero-diff, byte-identical `gate`) is evidence for a rejected superset design (retained `Refusal.Guard` + `UndecidedAtoms` beside it), not for the specified replace-the-field shape; the re-spike that would prove the real shape has not been run, and Phase 1's exit condition is a number that doesn't exist yet. |
| C-5 | B-3 | A19 Status: "**This is a REVERSAL of a closed joint decision**"; Prerequisites "Does not block lock" | The payload's transport to the CLI is undecided at lock between a structured `omitempty` field and the currently-ruling `Detail`-flattened-text fallback; A19's own evidence shows the fallback "defeats the total-order determinism whose only purpose is machine consumption," so "reversible in both directions" is not shown for consumers who already parsed one shape. |
| C-8 | B-1 | A22 ("UNSTATED INFERENCE against RDR 0002's current text") and the shared producer-half gap with A16/A24 | The kernel's presence/token/canonicalization guarantees are "by structure" only if RDR 0002 (Draft) accepts unlanded duties it does not currently carry; this relocates, rather than closes, the trusted-but-unbuilt-component risk the RDR rejects Alternative 1 for. |
| C-9 | B-5 | §SURVIVOR MEMBERSHIP resolution-scope veto; Failure Modes "`unless` over a rarely-set tag … footgun in practice" | The non-escapable, table-wide veto vetoes every artifact lacking an optional tag the moment one row references it, with no writable sanctioned mitigation and no downgrade valve; both models independently name this the RDR's highest-risk shipped property. |
| C-10 | B-4 | §GATE, THEN COUNT precedence + disposition mini-check "no code row yet" for `owned_state_unavailable`; `flow-guard-unevaluable` scoped to "supplied facts" | The refusal that actually wins precedence on the flagship case cannot be rendered correctly: no stable `Code` for `owned_state_unavailable`, and the existing code for `guard_unevaluable` is scoped/exit-grouped for the wrong remedy class. |
| C-14 | B-11 | Prerequisites box: five unresolved external document changes (§D4 reopen, two 0002 duties, 0003 placement ruling, 0009 re-lock) | The RDR locks against a promissory note over peer documents, several without fallbacks or schedules, one (0009's re-lock) with a live code conflict and no sequencing. |
| C's §3 pick | B-6 / B's §3 pick | A5 ("Verified (runtime); the load-time half rides on A12") + A12 ("Deferred … open") + Phase 4 grammar question | The sanctioned "absent OR equals v" two-row / conjoined-atom idiom is not confirmed writable (open grammar question) or loadable (Deferred overlap-check question) — both models pick this as the assumption least likely to survive first contact, and both predict the fallback is sentinel-stamping, the named anti-pattern. |

## A-only

| A row | Passage | Claim |
|---|---|---|
| C-3 | A21 Evidence: "the kernel enforces no part of it" | Factually contradicted by RDR 0009's (Final) third normative block, which gives the kernel an entry precondition (`Table.CheckValid`, `ErrEscapeShapeBreach`) — A21 quoted 0009's producer clause and stopped reading. |
| C-4 | Testing Strategy row 15: `escapeRow` "non-conforming on BOTH fields and must not be reused as-is" | The RDR names the fixture defect for one test but leaves five other frozen escape tests on the same non-conforming constructor, so the suite stays green while the composition A21 asserts is violated everywhere else. |
| C-7 | §SEAM normative block: Fixup-1d "keeps its verdict but changes its REASON"; Testing Strategy row 19 nil-seam-now-yields-a-plan | The RDR diagnoses that its own change destroys the only nil-seam test and prescribes a fixture fix in prose, but adds no frozen/new test for the newly-introduced "plans without a seam" direction. |
| C-11 | Normative empty-`unless`-block clause ("the term DROPS OUT") vs. RDR 0003's subtractive algebra | Two different algebras for the same row: 0003's subtraction over an empty `unless` disables the row (subtracts the whole product), 0007's rule keeps it; A9's claim that 0003 "requires the same reading" is unsupported. |
| C-12 | §Load-Bearing Decisions Naming + A23 "Pending" | Nine-plus new exported kernel names are stated as normative surface while A23 (their non-conflict with frozen boundary tests `REQ-25/36/37`) is Pending with evidence literally "Pending" — the sweep was never run. |
| C-13 | §Approach / Decision Rationale SQL:2003 strict-routine citation | The precedent is inverted: SQL propagates NULL and continues (a `WHERE` clause evaluating UNKNOWN just doesn't select the row); it supports the per-atom unevaluable rule but is silent on — and actually cuts against — the resolution-scope veto the RDR built on top of it. |

## B-only

| B row | Passage | Claim |
|---|---|---|
| B-7 | A13 "accepted exposure"; Testing Strategy row 16 (`--tag` decides a previously-unevaluable guard, plan issues) | The `--tag` bypass is sanctioned and test-pinned but untraceable: `Plan` carries no provenance field distinguishing a guard decided from caller-supplied observed state vs. artifact state, so the accepted exposure is unauditable after the fact. |
| B-8 | Background "an overlap key absent entirely yields `no_match` … closed-world by design" vs. Phase 4 "predicate-PLACEMENT guidance" | The exact masking path the RDR's Problem Statement opens with is not closed by the domain rule at all — it just moves: the same predicate expressed as a `Match` tag instead of a guard atom still turns absence into escapable `no_match`, and the only control is a docs handoff to Phase 4, not a kernel or lint rule. |
| B-9 | Payload reason set `absent` / `uncomparable` | The closed two-reason set collapses three distinct producer-side causes (unparseable value, evaluator/grammar token skew, canonicalization drift) into one label each, so the payload systematically misattributes who can fix the problem even though it correctly names what blocked. |
| B-12 | Payload clause: "a row pruned as GuardFalse under `F ∧ U = F` … is not a survivor, so its absent keys are not reported" | The prune-for-soundness path (necessary for D8) leaves the single most common authoring question — "why didn't my row fire?" — with no surface anywhere: no payload, no refusal, no trace names the absent key on a pruned row. |
| B-13 | Phase 3: `contains` present-key contract "cannot be written from any current document" | The one operator with a named empty-set masking hazard ("absent set ≠ empty set") ships with no executable present-key contract test at lock, because it's blocked on an RDR 0003 declaration that doesn't exist yet. |

## Disagreements

None found. Every place the two passes address the same passage, they assert the same defect in the same direction (see Convergent table) — no case where one model asserts X about a passage and the other asserts not-X. The passes differ in *coverage* (A-only and B-only sections above) and occasionally in *framing/emphasis* — e.g. Pass A treats the existence-polarity inversion (C-1) and the spike/spec divergence (C-6) as two separate findings, where Pass B folds both into one (B-2) — but this is a granularity difference, not a substantive contradiction: both describe the identical underlying fact (spike code disagrees with normative clause on `exists` polarity) and draw the identical conclusion (Phase 1 will implement the wrong rule from the spike).

## Merged resolve queue

Union of distinct defects, ordered for the fix pass by severity (masking-path and lock-blocking issues first, then correctness/consistency gaps, then diagnosability/hygiene gaps). Not dispositioned here.

1. **Existence-operator polarity inversion** — spike/probe decide `exists` from presence alone (ignoring `Literal`), opposite of the normative `presence == literal` rule; Phase 1 is pointed at the spike as its starting shape. (C-1, B-2)
2. **A3 spike proves a rejected superset design, not the specified shape** — retained `Refusal.Guard`+`UndecidedAtoms`, `slices.MinFunc` representative selection, and the literal-free existence rule all diverge from the Normative Contracts; re-spike not yet run; Phase 1 exit condition is an unknown number. (C-6, B-2)
3. **Resolution-scope aggregation veto has no writable mitigation** — non-escapable, table-wide refusal the first time any row references a newly-optional tag; sanctioned two-row/conjoined-atom idiom rides on Deferred A12 and an open Phase 4 grammar question. (C-9, B-5, C's §3 pick, B-6/B's §3 pick)
4. **Payload transport to the CLI is undecided at lock** — A19 reverses a Closed JDR 0001 §D4 decision, marked "does not block lock," but the stated fallback (flatten into `Detail`) is shown by the RDR's own evidence to defeat the payload's only purpose. (C-5, B-3)
5. **A21's "kernel enforces no part of" the escape-shape composition is contradicted by RDR 0009 (Final)'s kernel-side precondition (`Table.CheckValid`/`ErrEscapeShapeBreach`)**, and the fix is unsequenced against the frozen escape fixtures (five tests share the non-conforming `escapeRow` constructor). (C-3, C-4)
6. **Winning refusal on the flagship case (owned-state-unavailable-before-guard-unevaluable) has no renderable Code / correct exit-group** in RDR 0005, and neither gap is scheduled for fix. (C-10, B-4)
7. **Producer-half guarantees for presence/token/canonicalization (A16, A22, A24) rest on unlanded RDR 0002 duties** — same trusted-but-unbuilt-component risk the RDR rejects Alternative 1 for, just relocated to the normalizer seam. (C-8, B-1)
8. **Five external-document prerequisites at lock, several without fallback or schedule**, including a live, unsequenced conflict with RDR 0009's Final kernel precondition. (C-14, B-11)
9. **Nil-seam narrowing (Testing Strategy row 19) has no discriminating test** for the newly-introduced "guard plans without a seam" behavior, only a prose fix for the frozen test it breaks. (C-7)
10. **Empty `unless` block algebra conflicts with RDR 0003's subtractive algebra** — 0007 drops the term, 0003's subtraction (read literally) disables the row; A9's claim of agreement is unsupported. (C-11)
11. **New exported kernel surface (A23) not yet swept against frozen boundary-symbol tests** (`REQ-25/36/37`) before lock. (C-12)
12. **SQL:2003 strict-routine citation supports the per-atom unevaluable rule but not the resolution-scope veto** — the veto ships with no cited prior art once the citation is read at the level it's used. (C-13)
13. **`--tag` bypass (accepted exposure, A13) is unauditable** — `Plan` carries no field distinguishing a guard decided from caller-supplied state vs. artifact state. (B-7)
14. **Masking path reopens one field over, via `Match` tags** — the same absent-key-masked-as-`no_match` pattern the RDR's Background opens with still exists when the predicate is authored as a `Match` tag instead of a guard atom; only mitigation is a Phase 4 docs handoff. (B-8)
15. **Closed reason set (`absent`/`uncomparable`) can't distinguish who-can-fix-it** — collapses distinct producer-side causes (bad value, token skew, canonicalization drift) into the same two labels. (B-9)
16. **Pruned rows (`F ∧ U = F`) report nothing** — no surface answers "why didn't my row fire" for the most common debugging question, even though the prune itself is sound. (B-12)
17. **`contains` operator's present-key contract is unwritable at lock**, blocked on an undelivered RDR 0003 declaration — the one operator with a named empty-set masking hazard ships untested on that hazard's other half. (B-13)
