# Reconcile - RDR 0003 Guard Predicate Exhaustiveness

## Inputs

- Pre-lock lists:
  - Grounding: none.
  - 3amigo: `contains` MVV coverage, `unless` disabled-row MVV case.
  - Critique: multi-dimensional gap/overlap MVV cases, set-valued deterministic proof, too-large product refusal/downgrade, RDR 0005/RDR 0006 ownership chain.
- Pending or unverified Critical Assumptions: none.
- Named spikes: `docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes/check.sh guard-fixture.toml`, captured in `docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes/output.txt`.
- Exactness-word delta: pre-lock fixes introduced or touched exactness claims for multi-dimensional product proof, deterministic set-valued proof, source-order independence, and full conjunctive `unless` semantics.

## Dispositions

| item | source | disposition | evidence pointer or plan |
| --- | --- | --- | --- |
| `contains` predicate over declared set-valued tag must be covered before accepting full vocabulary | 1, 4 | DOWNGRADED | Not a pre-lock fact about the existing spike. Recorded in A1 Reconciliation, Phase 3, MVV Scenario 1, and Testing Strategy as an implementation MVV requirement. |
| `unless` disabled-row behavior must be tested when all positives match and the full exclusion block matches | 1, 4 | DOWNGRADED | Recorded in A3 Reconciliation and MVV Scenario 1. |
| Multi-dimensional row-group gap and overlap cases must be visible only in the scoped product | 1, 4 | DOWNGRADED | Recorded in A2 Reconciliation, Risks and Mitigations, MVV, and MVV Scenario 2. |
| Set-valued domains must use deterministic symbolic or bitset-equivalent proof | 1, 4 | DOWNGRADED | Recorded in A2 Evidence, Normative Contracts, Technical Design, Performance Expectations, and MVV Scenario 2. |
| Finite scoped product too large to prove must refuse or downgrade, not silently cap | 1, 4 | DOWNGRADED | Recorded in A2 Reconciliation, Normative Contracts, Failure Modes, Phase 2, and MVV Scenario 3. |
| RDR 0005 and RDR 0006 must preserve predicate semantic-kind ownership chain | 1, 4 | DOWNGRADED | Recorded in A4 Reconciliation, Load-Bearing Decisions, Phase 4, and MVV Scenario 4. |
| Resolve spike named in A1 | 3 | VERIFIED | `cd docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes && sh check.sh guard-fixture.toml`; captured output at `docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes/output.txt`. |

## Completeness Check

- No Critical Assumption remains Pending or Unverified.
- No `Docs Only` load-bearing assumption remains.
- No named spike lacks captured output.
- No draft-placeholder or seed-skeleton marker remains in the body sections.
- Finalization Gate written-response placeholders remain intentionally unfilled for Stage 7.

## Verdict

RECONCILED. All open items are terminal, no blocker refutes the approach, and the remaining obligations are named implementation MVV requirements rather than pre-lock spikes.
