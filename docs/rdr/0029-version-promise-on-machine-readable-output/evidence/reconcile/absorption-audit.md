# Absorption audit — cli/0029 pre-lock rounds

Status: Draft (pre-lock). All rounds show iter-1 findings fixed, and the
iter-2 delta re-checks (which specifically target defects introduced by
the iter-1 fixes) confirm resolution with no further residue.

## Grounding round
- GF1-GF4 (iter-1): tier census gaps (findings[].operator/block untiered),
  A3 strictness understatement, A2 "only" overclaim, decodeStrict count —
  all fixed; iter-2 delta-scoped re-check covered only the two claims
  iter-1 added (block/operator tiers).
- GF5 (iter-2): `findings[].block` wrongly assigned `frozen`, contradicting
  Final `JDR 0001 §D12`. Fixed — current C4 assigns `append-only` and cites
  §D12 (grew 2→3 members).
- GF6 (iter-2): `findings[].operator` verdict right, reasoning wrong (cited
  wrong mechanism). Fixed — current C4 cites `0003:C7` and
  `TestReq2_OperatorVocabularyIsClosedAndTyped`.
- iter-2's own "Still open" section: None.
Verdict: fully absorbed, no residue.

## 3amigo round
- Consolidation: 27 findings (hotspots on C4 seam census, C2 `closed`
  retirement scope, C1 schema_version home, planEnvelope misdescription,
  S1/S2 key-set baselines, S8 cardinality-on-growing).
- iter-2 delta re-check (F1-F4, all on the iter-1 rewrite): F1 (C1/MVV
  mischaracterized `planEnvelope` as an input decoder) — current C1 text
  correctly states it consumes an *emitted* envelope, output rule governs,
  not `0002:C3`. F2 (C4 seam table omitted 6 of 16 tiered vocabularies,
  arithmetic disagreed with §capability-dependencies) — current C4 census
  table lists all 16 rows including envelope `type` and `findings[].block`.
  F3 (S7/S8 retiring only one of three cardinality-on-append-only tests)
  — current S7 names all three (TestReq73, TestReq74, TestReq80) as the
  anti-pattern, backed by A7's CI snapshot. F4 (S9/§step-2 grep scope
  ambiguous, ~91 files) — not independently re-verified against current
  §step-2 text in this pass; flagged below as the one item not directly
  re-confirmed post-fix (no iter-3 3amigo file exists).
- Charted (not residue, dispositioned): tier table exposed via binary
  (successor RDR), release-note mechanism for C3 (successor kata/RDR),
  conformance fixture for C1 consumer rules (successor RDR if 3rd party
  appears), docs-vs-contract drift check (successor kata).
Verdict: absorbed, with one caveat — F4 (§step-2 grep scope) has no
explicit post-fix confirmation on file; current §step-2 text should be
checked at Stage 6/7 sweep if not already narrowed to C4-tiered
vocabularies only.

## Critique round
- Consolidated ledger D-1..D-20: 5 blocking-lock (D-1 envelope type={ok}
  seam, D-2 schema_version-is-not-discriminator self-contradiction, D-3
  unanswered Proportionality gate, D-4 C3 disclosure unbacked, D-5 A5
  Pending on largest vocabulary), 13 should-fix, D-19 charted (A1 peer-RDR
  amendment — routed to Stage 7.1 cluster reconcile over {0006,0029}, per
  existing JDR precedent), D-20 dismissed (A2 forcing-function claim,
  refuted by A2's own recorded "If wrong" clause and §phase-1 text —
  dismissed.md cites this exactly).
- iter-2 critique (post-fix): verdict "No blocking gaps." All D-1..D-18
  passages independently re-verified SOUND against current tree/text
  (envelope {ok} seam correct, A7 added and consistently Pending
  everywhere cited, A6 correctly Pending-on-mechanics, Proportionality
  gate now has a full written response, S7/A7 distinction real). One nit
  (R-1): C4's "failed is reserved, never emitted" framing reads sharper
  than stale doc comments in respond.go/flow_input.go that still say
  `{"type":"failed"}` — not a rewrite-introduced defect; Step 2 already
  schedules the doc-comment fix.
Verdict: fully absorbed. D-19 correctly terminally routed to 7.1 (not
residue of this lock); D-20 correctly terminally dismissed.

## Repeatability round
run-1 reconstruction vs. current C1-C4/MVV/S1-S9: zero findings admitted.
The two live GUESS-shaped items (A5's seam-mechanics outcome, A7's CI
snapshot generator mechanism) are not residue — both are RDR's own
tracked Pending assumptions with spikes and stated "If wrong" fallbacks
(seam: none (prose-only) / review-only disclosure), not silent gaps.

## Open Critical Assumptions implied by findings (not residue — tracked)
- A5 (Pending, Spike): whether CLIError `code` and `flow next` reason
  seams are buildable without inverting clierr's leaf-package layering.
  Fallback stated if wrong.
- A6 (Pending, Source Search): schema_version constant home in clierr —
  layering already grounded, pending only on edit landing.
- A7 (Pending, Spike): CI snapshot diff mechanism for seam/severity
  disclosure — generator shape (hidden verb / go:generate / fixture) and
  import-boundary readability unresolved. Fallback (review-only
  disclosure, recorded as a limit) stated if wrong.
None of these are silently open — each has a Method, Evidence and an
"If wrong" clause naming the fallback, so they are Stage 4/6 work, not
lock blockers left unaddressed by the record.

## New claims introduced during fix passes with no independent verification
- A7 itself is a new claim (critique round introduced it to back C3).
  iter-2 critique verified A7's Evidence section cites real, checked
  facts (Capability Dependencies' prior "caught by a reviewer or not at
  all" concession, `graphlint::AdvisoryCodes`/`BlockingCodes` existing)
  and gives a falsifiable verification method. Not verified: that the
  spike will actually succeed — correctly left Pending, not asserted.
