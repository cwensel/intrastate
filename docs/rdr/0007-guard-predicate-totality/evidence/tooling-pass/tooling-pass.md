Model: claude-opus-5[1m]

# Tooling Pass — RDR 0007 guard predicate totality

Mechanical adherence sweep, run as the Stage 7 pre-step to the Finalization
Gate. Post-mutation regression check: four pre-lock lenses and the Stage 6
reconcile rewrote this draft.

Run 1 (this file). No prior reports in `evidence/tooling-pass/`, so no
loop-breaker comparison applies.

## CHECK 1 — Template section coverage

Every **Required** (spine) section is Present-substantive: Metadata, Problem
Statement, Critical Assumptions, Proposed Solution / Approach / Technical
Design / Normative Contracts, Decision Rationale, Alternatives Considered,
Context, Research Findings, Trade-offs (Consequences / Risks / Failure Modes),
Implementation Plan, Minimum Viable Validation, Validation / Testing Strategy,
References (12 real citations), Finalization Gate.

**Conditional sections legitimately deleted** (PASS, not Missing, per the
false-positive guard): Round-Trip / Inverse Invariants (no
encode/decode pair — the A10 guard-structure mapping is a one-way recovery
obligation, not an inverse pair), Illustrative Code, Day 2 Operations (creates
no persistent resource), New Dependencies, Performance Expectations, Metadata
`Cluster` (stands alone).

**Conditional sections present and filled**: Load-Bearing Decisions
(Selection/predicate + Naming), Capability Dependencies (4 rows), Existing
Infrastructure Audit (6 rows, two marked Build).

No `_Draft placeholder._`, no seed-skeleton header, no surviving bracketed
template instruction outside the Finalization Gate section, and no block from
an older TEMPLATE.md. The Finalization Gate section still carries its template
body — correct pre-lock state; it is replaced by the gate.md pointer at lock.

Verdict: **PASS**.

## CHECK 2 — Method label vocabulary

16 Evidence Records; every Method is in the sanctioned eight.

- Source Search ×13 (A1, A3, A6a, A6b, A7, A8, A9, A10, A11, A12, A13, A14, A15)
- Derivation ×1 (A2)
- Prior Art ×1 (A4)
- Design Decision + Source Search ×1 (A5 — compound, but both halves are
  sanctioned labels and the record splits them explicitly by clause: "Design
  Decision (runtime) + Source Search (load-time)". Not paraphrase, not
  off-vocabulary.)

No record is missing a Method. No `Docs Only` anywhere.

Verdict: **PASS**.

## CHECK 3 — Source Search self-reference

No Source Search Evidence path resolves to `{RDR_PATH}` or to this RDR's
artifact directory. The five in-document self-path references are all
evidence-artifact citations (`evidence/research/*`,
`evidence/propose-premortem/critic.md`, `evidence/spikes/aggregation-probe.md`),
which is the sanctioned form — cached quotes and spike output, not the RDR
proving itself.

Verdict: **PASS**.

## CHECK 4 — Docs Only on load-bearing claims

Zero `Docs Only` records in the RDR. The only occurrences of the string are in
the Finalization Gate's own template text describing the check.

Verdict: **PASS**.

## CHECK 5 — Symbol resolution of Source Search / Spike anchors

Delegated; all 40+ anchors resolve **in their cited files** on branch
`via-claude`. Structural claims the normative blocks depend on were confirmed
rather than assumed:

- `GuardEvaluator` is an **interface** with exactly
  `Evaluate(guard string, view TagSet) GuardResult` — the Normative Contracts'
  "INTERFACE with the single method … not a function type" clause holds.
- `Row.Guard` is `string`; `Row.RequiresOwned` is `[]string`.
- `Refusal.Rows` is `[]RowRef`; `RowRef` is `(RuleID, SourceLocator)` — the
  element type the aggregation clause pins as a normative assertion target.
- `fixtureGuards.Evaluate(guard string, _ resolve.TagSet)` genuinely discards
  the view, and `fixtures_test.go` is `package resolve_test` (unimportable) —
  both load-bearing for Phase 2's "not free reuse" claim.
- `TagSet.Lookup` is the only exported accessor (A11's forcing argument).
- Test anchors resolve: ADV-1, ADV-1b, ADV-2 (with subtest "guard UNEVALUABLE
  must not rescue"), ADV-3, ADV-4, ADV-5, Fixup1d,
  TestFixupGateIsUniformAcrossOrdinaryAndEscapeCandidates, Req7.

Negative claims verified: no `testdata/` repo-wide; `internal/cli/root.go`
registers only `version`; no `flow` package; no production importer of
`internal/resolve`. These bound A13's exposure and Phase 2's harness gap.

All cited evidence files exist (0007 research/spikes/premortem, 0001
deviations.md, 0002 + 0003 spike fixtures).

No bare `file:line` anchors in the RDR body.

Verdict: **PASS**.

## CHECK 6 — Status consistency

One assumption is not `Verified`: **A6b** (`Pending` — downgraded by explicit
decision at Stage 6, obligation routed to RDR 0004's implement stage). Every
one of its nine references treats it as **open** (Risks "A6b remains OPEN",
Failure Modes "confirmed open", Prerequisites "carried open by decision",
Testing Strategy "Out of scope here", the matrix's operability row "while A6b
is open"). No settled-fact prose depends on it.

The premortem ledger's "A6 widened" (Decision Rationale) is historical
narration of the A6a/A6b split, not a live reference to a nonexistent A6.

No checklist-box vs gate disagreement: the Finalization Gate carried no
responses before this pass, so there was no second copy to contradict.

Peer-RDR citations independently verified (39 quotations, delegated): 37
verbatim hits; all peer statuses match the RDR's claims (0001 Implemented;
0002–0006 Final; 0003 Final-and-unimplemented). The RDR's demonstrated
negatives hold — RDR 0005 contains zero occurrences of "observed", and RDR 0006
genuinely leaves row-group membership undefined beyond "the same
state/outcome", which is exactly what A12 routes to 0006's implement stage.

Verdict: **PASS** (after the mechanical fixes below).

## Mechanical findings fixed in-pass

Per the Stage 7 contract, mechanical findings are repaired here and the sweep
re-run rather than routed to Refine (conformance is outside Refine's contract).

1. **C6 — MVV scope contradiction (load-bearing).** The Minimum Viable
   Validation closed by placing "mixed-verdict combination cases … assert the
   strong-Kleene selection" *in MVV scope*, while the qualifier paragraph
   immediately below and Testing Strategy row 4 both place strong-Kleene
   combination in Phase 2, blocked on A10. Verified structurally unreachable:
   `fixtureGuards` maps a whole guard *string* to one verdict and discards the
   `TagSet`, so atom-level `F ∧ U = F` cannot be expressed against it at all.
   Fixed by scoping the MVV sentence to what it can prove and pointing at row 4.
2. **C6 — Risks mitigation resting on (1).** "Mitigation: A2's truth-table
   derivation at Resolve; MVV scenario exercises the mixed-verdict
   combinations" asserted an executable check that is not in MVV scope. Fixed
   to "partial", naming derivation as the whole of what holds today and row 4
   as the blocked executable check. Both fixes are drawn from the RDR's own
   existing text, not new evidence.
3. **C5 — quotation drift.** RDR 0002 reads "normalization **renders as** a
   `<clear>` write"; the RDR dropped "as". Corrected.
4. **C5 — citation over-attribution.** RDR 0003's A4 is titled "Guard predicate
   errors can use the existing structured CLI failure gateway" and *enumerates*
   "unevaluable guard" among the semantic kinds it owns; it does not "claim" the
   kind as its subject. Softened at both sites (Technical Environment, Approach).

Sweep re-run after the fixes: all four resolved, no regression, no new finding.

## Verdict

**PASS** — no findings survive; proceed to the Gate's written responses.
