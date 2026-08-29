# Finalization Gate — cli/0024 declared-emit-vocabulary

Date: 2026-08-29
Verdict: **READY — Gate PASS**
Mechanical pre-sweep: PASS (`evidence/tooling-pass/tooling-pass.md`)

## Contradiction Check

No contradictions found between Research Findings, the design principles this
RDR inherits, and the Proposed Solution.

Three places where a contradiction would be expected were checked against the
record's own text and each is resolved *in* the record, not left standing:

- **The DMN alignment claim.** Research Findings marks the DMN output-values
  class "SEARCHED AT RESOLVE AND NOT FOUND" across four corpora plus the
  sibling analysis repo, and A6 is Pending/downgraded to match. The Proposed
  Solution's alignment argument correspondingly runs on the in-repo mirror
  (`0002:C22`) and the two opened peer citations (ms-conductor
  `type: terminate`, scxmlcc `<final>`) — never on DMN. Finding and solution
  agree.
- **Enforcement tier.** Research Findings derives from `0002:C24` that
  single-rule checks belong to load and that the advisory tier is closed
  (`0006:C17`); C2 places both checks in the load pipeline ahead of
  `normalizeRules` and makes all three categories refusals, never advisory. No
  check in this RDR is advisory-shaped, as the discovery requires.
- **Fail-fast vs report-everything.** C2 states the tension explicitly rather
  than eliding it: the project accumulates findings only where a total order
  exists (naming `graphlint/engine.go::sortFindings` and
  `resolve.go::compareRefs` as those comparators), load has neither and
  declares its order unspecified, so this RDR takes fail-fast and states what
  it therefore does NOT claim — an emit refusal may mask a coexisting
  structural one. The Consequences section carries the matching adoption cost
  (an N-round fix-and-rerun loop on first adoption) rather than claiming a
  cost-free rollout. The stated principle and the stated consequence match.

The one planned behavior that touches a locked peer — regenerating 0011's 28
byte-identity goldens — is not a contradiction of `0011`'s guarantee but a
named amendment to it: A2 states the guarantee is scoped to 0011's own change
class (A15 iii), that `dispositions` is an append 0011 never contemplated, and
that the regeneration is licensed EXPLICITLY by name with the enumerated edit
set in Phase 3. Priced, not hidden.

## Assumption Verification

All nine Critical Assumption records are internally consistent: Status, Method,
and Evidence agree, and every "If wrong" is non-empty (A8's reads `n/a — the
claim is refuted, not carried`, followed by the consequence that actually
landed, which is a real disposition rather than an empty field).

- **Verified (6)**: A1, A2, A4, A8-as-refuted, A5, A9. A3 and A9 rest on
  captured spike runs on disk (`evidence/spikes/a3-corpus-sweep.md`,
  `a9-rule-side-locator.md`); A1/A2/A4 on source anchors; A5 on `0010:A6`, a
  Verified peer element, plus an empirical sweep of the consumer's two models.
- **Method `Docs Only`**: none. No record uses the label, so nothing is blocked
  on that ground.
- **Refuted (1) — A8**, at pre-lock by the critique lens on both models:
  `TestReq146_EveryEmittedSequenceIsASortedSlice` reflects over
  `table.Row` in all three sub-tests and is structurally unreachable from
  `Model.EmitDecls`. Absorbed with no contract change — C3 now states the sweep
  does not reach the carrier and names Testing Strategy scenario 4 as the sole
  determinism oracle; Phase 2 owes no `TestReq146` edit. The refutation costs
  the contract nothing because scenario 4 was already written as a direct
  assertion.
- **Pending (2), both downgraded at reconcile by decision, not omission**:
  - **A6** (DMN allowed-values on output clauses) — not load-bearing; a
    corroborating external citation only. Searched exhaustively and not found;
    the verification plan (read OMG DMN 1.3 §8 from the spec PDF, a manual
    corpus acquisition) is deliberately UNSCHEDULED because nothing in the
    design waits on it. No prose in the RDR treats it as settled.
  - **A7** (C1's two-shaped `domain` under strict decoding) — the DECODE legs
    are Verified by spike; only the CLOSURE leg is carried, i.e. that C1's
    hand-written arms are the complete set of shapes the decoder stops
    catching. Its plan **will run before implementation completes**: Testing
    Strategy scenario 1, a table-driven fixture per arm plus a negative sweep,
    inside Phase 1, which does not close until it is green. Not MVV-critical —
    the MVV drives the declared-then-green path, not the exhaustive malformed
    sweep — and the Prerequisites checklist carries it as the one unchecked box,
    correctly.
- **Self-reference**: none. Every `Source Search` Evidence path names a repo
  file (`internal/table/source.go`, `internal/cli/flow_resolve.go`,
  `internal/table/roundtrip_test.go`, `internal/table/load.go`); none resolves
  to this record or to its artifact directory.
- **Anchor resolution on `main`**: all 39 `source-anchor` edges report
  `resolved: true`. The two edges with `resolved` absent are a kata id
  (`issue/srz2`) and this gate artifact — neither is a source symbol.
- **Status consistency**: no `Pending` or `Refuted` property is relied on as a
  settled fact anywhere in the RDR. A6's claim is explicitly not leaned on in
  Research Findings, Approach, or any contract; A7's open leg is named in Risks
  as the residual cost with scenario 1 as its gate; A8's refutation is stated
  inside C3 itself and in scenario 4, so the two sources cannot disagree. The
  Prerequisites checklist boxes and this gate agree item for item.

## Scope Verification

The Minimum Viable Validation is **in scope and executed during
implementation**, not deferred.

The specific proof is the four-step MVV, and its steps land inside the phased
plan rather than after it:

1. Re-author the seed's two-rule adversarial table (misspelled command value;
   key typo'd `nxet`) and confirm it still lints exit 0 with zero findings —
   the defect reproduced as a negative control.
2. Add an `[emit.next]` enum declaration with a route/stop-partitioned domain
   and run `intrastate lint`. Expected: nonzero exit with ONE blocking finding
   carrying its category slug in `code` and a real source line in `locator`;
   fix and re-run to surface the second, so `emit_value_out_of_domain` and
   `unknown_emit_key` are both observed across the two runs. The one-finding
   cardinality is asserted, not assumed —
   `internal/table/reserved_key_0008_test.go` forbids a multi-error load and
   `evidence/spikes/c2-finding-multiplicity.md` measured it.
3. Fix both rules, lint exit 0, then `flow resolve` a rule whose value is
   listed under `stop`: the payload carries the authored `emit` unchanged and
   `dispositions` mapping the key to `stop`, positioned immediately after
   `emit`.
4. Delete the `[emit]` table and re-run lint and resolve over the original
   defective table: exit 0, no findings, `dispositions: {}` — today's behavior
   apart from the appended empty field.

The locator promise in step 2 is the leg that was open and is now closed: A9 is
Verified with both the declaration-side and rule-side techniques established
(`evidence/spikes/a9-rule-side-locator.md`), and the corrected id-anchored
rule-side technique is carried into C2 and Phase 1. Nothing the MVV proves rests
on a Pending assumption. Testing Strategy scenarios S1–S9 carry the balance,
with scenario 1 discharging A7's closure leg and scenario 4 the sole determinism
oracle for C3's sorted union.

## Proportionality

Right-sized. No section flagged for trimming before lock.

- **Contract count — one seam.** C1–C4 are four facets of the single contract
  this RDR owns, stated by surface: grammar (`[emit.<key>]`), proof (the load
  pipeline), carry (`Model.EmitDecls`), envelope (`dispositions` on the resolve
  payload). The record says this in its own Normative Contracts preamble, and
  it holds under inspection — none of the four is independently authorable
  without the others: C2 proves what C1 declares, C3 carries what C2 built, C4
  reads only C3's carry. No second seam to split off.
- **Profile re-validated: `foundational` stands.** One contract, but it locks a
  TOML authoring grammar, spans `internal/table` and `internal/cli`, appends
  three load categories to a locked peer's "at minimum" set, appends a field to
  a locked payload contract, and produces a declaration carry peers 0021/0023
  consume. The lens battery actually run matches the profile row: `cove`,
  `3amigo`, `critique` (two models, with `diff.md`), and `repeatability` in the
  **full** variant (`run-1.md` `variant: full`, plus run-2, run-3 and
  `diff.md`). No lens is missing and none was run against a stale draft — the
  Stage 6 reconcile postdates them and its absorption audit is complete
  (cove 9/9, 3amigo 30/30, critique 21/21, repeatability 13 disagreements + 7
  GUESS clusters). The field's form is correct: value plus one clause naming
  the contract; no matrix or provenance prose survives from the template.
- **No `Transient` contract**, so nothing is discounted in that count.
- **Length is earned, not padding.** The long passages are the ones carrying
  decisions an implementer would otherwise re-litigate or get wrong: C1's
  no-usable-domain single-arm ruling (a fixture per arm is unwritable, and
  saying so prevents a test that pins a discrimination the decoder cannot
  make), C2's fail-fast rationale with its named comparators, C3's
  safe-by-omission dependency on C1's arms (which fails silently if relaxed),
  C4's `omitempty`-deliberately-not-taken clause, and A2's enumerated
  peer-spec amendment. The four over-budget Evidence fields (C9 advisory) each
  keep their load-bearing anchor at the head and point at a spike file; the
  balance is verification content the grounding sweep reads, so it stays where
  it is rather than moving to artifacts.
