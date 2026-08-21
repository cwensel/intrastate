Model: claude-opus-5[1m]

# Cove iteration 2 — dispositions

One line per finding: disposition · origin ledger entry · RDR section touched.

| ID | Disposition | Section touched |
| --- | --- | --- |
| F-1 | **fixed** (direction upheld via §strong-consult; framing corrected) | A19 Status; Prerequisites |
| F-2 | **fixed** — A3 flipped Verified → Pending with a re-spike plan | A3 Status |
| F-3 | **fixed** — 0009 re-lock prerequisite added | Prerequisites |
| F-4 | **fixed** — A16's producer half separated from its kernel half | A16 Evidence; Prerequisites |
| F-5 | **fixed** — nil-seam narrowed to a per-atom rule | Normative Contracts (SEAM); Testing Strategy row 19 |
| F-6 | **fixed** — 0009 quote restored to its producer-obligation scope | A21 Evidence |
| F-7 | **fixed** — gap named, encoding left to RDR 0003 | Phase 3; Phase 4 |
| F-8 | **fixed** — A19 given a fallback; "Blocks lock" removed | A19 Status; Prerequisites |
| F-9 | **fixed** — proposed symbols marked | A3 Status; Testing Strategy row 17 |
| F-10 | **fixed** — audit row states the first-nested-payload cost | Existing Infrastructure Audit |
| F-11 | **fixed** — both quotes corrected | A16, A17 Evidence |
| F-12 | **fixed** — `MinFunc` retirement stated | Normative Contracts (payload); Phase 2; Testing Strategy row 20 |
| F-13 | **fixed** — doc-comment rewrite added to Phase 2 | Phase 2 |

No finding was dismissed and none was charted — every one
grounded against source and every one landed inside the
delta scope (the seven re-verify assumptions and the
contracts depending on them).

## Mini-checks fired

Four of five cues fired; tables written into the RDR under
**Validation → Mini-checks**: `authority`, `oracle`,
`disposition`, `trace`. Round-trip / fidelity did NOT fire
(the draft's "migrate" is fixture migration; no inverse
invariant is claimed).

**The desk trace found one CONTRADICTION, now fixed.** MVV
scenario 1 gave its row an absent `RequiresOwned` key AND an
unevaluable atom while asserting `guard_unevaluable`. Under
the draft's own SURVIVOR MEMBERSHIP and GATE-THEN-COUNT
clauses — and under shipped `resolve.go::gate`, which returns
`KindOwnedStateUnavailable` at L388 before reaching the
undecidable loop at L396 — that input yields
`owned_state_unavailable`. The RDR's primary validation
scenario contradicted three of its own clauses and would
have failed on first implementation. Fixed by requiring the
row's `RequiresOwned` to be satisfied; the both-ways pairing
remains Testing Strategy row 13's subject.

This is the pass's second-most valuable outcome after F-1,
and neither cove half found it — it took the stepwise trace.

## §strong-consult record (F-1)

The A19 fork — reverse JDR 0001 §D4's `Detail` routing, or
comply with it — was taken to one fresh-context consult
before any escalation, per the tiebreaker-reduction gate.

Returned **PASS** (not NEEDS_DECISION), collapsing the fork
toward A19's structured field on engineering grounds:

- `internal/cli/clierr/clierr.go:46-47` — the type's own doc
  invites it: "Extend with new optional fields as needed —
  keep them `omitempty` so the envelope stays append-only and
  stable for tools."
- §JD-8's "not new envelope fields" is a SUFFICIENCY claim
  about the `Code` values RDR 0005 enumerated, inherited
  verbatim from 0005's A5 evidence
  (`0005-skill-integration-cli-contract.md:188`) and written
  before this payload existed — not a prohibition.
- `Detail` is documented human prose; nesting a sorted
  two-level array inside it defeats the determinism whose
  only consumer is a machine.
- Fully reversible in both directions, which is why this
  does not warrant blocking lock.

Both claims were re-verified in the main context before
being acted on. No escalation to the author was needed —
but the REVERSAL still must be taken back to JDR 0001,
since §D4 is Closed. That is now a Prerequisite, not a
silent in-draft correction.

## Needs (re)verification — carried to Stage 6

| ID | State | Why |
| --- | --- | --- |
| A3 | Verified → **Pending** | The spike proved a superset (kept `Refusal.Guard`). Re-spike with the field actually removed and record the frozen suite's disposition for `TestFixup1d…`. Bound + importer halves stand (re-run this pass: CLEAN, 154 PASS / 0 FAIL). |
| A19 | **Pending** (unchanged) | Now correctly framed as a reversal of a CLOSED §D4. Needs the JDR reopening; fallback recorded. |
| A22 | **Pending** (unchanged) | Duty still absent from 0002's direction list; now paired with A16's producer duty in one Prerequisite. |
| A16 | Verified (kernel half) | Producer half explicitly downgraded to "acknowledged, not a clause" — no status flip, but the Prerequisite now carries it. |
| A17 | Verified | New load-bearing gap named (set-valued element encoding is RDR 0003's). Not a flip: the absent-key rule A17 supports needs no encoding. |
| — | **new claim** | Nil-seam per-atom narrowing (Normative Contracts) is a behavior change against shipped `evaluateGuard`; Testing Strategy row 19 is its test. Verify at Phase 1 that Fixup-1d's fixture is re-read, not re-encoded. |
| — | **new claim** | `MinFunc` retirement (row 20). Verify no other reader depends on a single representative row. |

## Tiebreakers

None outstanding. F-1 was the only genuine fork and the
consult collapsed it.
