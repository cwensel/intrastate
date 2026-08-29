# RDR 0013 — Stage 2 prior-art search record

Date: 2026-08-28. Budget: ≤3 corpus queries per claim, ≤5 opened hits per
claim. Tool: `arc search semantic --corpus <C> --limit N --json "<q>"`.

## Claim P1 (class): peer bounded-analysis tools treat the analysis bound
as a declared scope, and proof cost may not track state-space size

- Query 1 (StateMachineLit, limit 4): "bounded model checking bound depth
  completeness threshold analysis scope" — HIT.
  - Accepted citation A: TraceFix (Proceedings of the ACM Conference on AI
    and Agentic Systems, 2026), page 7: "The critical observation is that
    *verification time does not track state-space size*." — supports
    0003:A15's original doubt that a single scalar cardinality predicts
    proof cost for non-enumerating representations.
  - Accepted citation B: same paper, page 9: "bounded counters, bounded
    queues); outside those bounds the guarantee is incomplete." — the
    definitional-scope stance: a bounded verifier's guarantee is *scoped by*
    its bound; the bound is not an empirical prediction to be refuted.
  - Rejected hits: page 12 (metrics table, off-topic), page 9 chunk 23
    (related-work list, off-topic).

## Claim P2 (class): a test oracle must know the answer independently of
the system under test

- Query 1 (PapersFast, limit 4): "test oracle problem pseudo-oracle
  differential testing independent oracle" — MISS (all four hits were
  index/bibliography noise; rejected).
- Query 2 (DevRef, limit 4): "test oracle must be independent of the
  implementation under test self-verifying tautology" — HIT.
  - Accepted citation: developer-testing.pdf (DevRefGit/Books2/tdd + bdd),
    page 171: "an oracle is a black box that knows the answer to a problem
    ... How does it know? You program it to!" — the oracle's answer comes
    from outside the computation under test; an oracle defined *as* the
    system's own gate (`ProofCompletes` = `card <= Bound()`) observes
    nothing.
  - Rejected hits: xUnitTestPatterns pages 917/531/69 (index pages and
    behavior-verification pattern text, not oracle independence).

## Instance question: what do peer state-machine tools do when a state
space is too large to enumerate (refuse with published bound vs cap vs
enumerate)?

- Query 1 (StateMachineRes, limit 5): "state space too large to enumerate
  bound refusal explicit limit exhaustiveness checking" — MISS (heading
  fragments from unrelated repos; all rejected).
- Query 2 (StateMachineRes, limit 5): "Alloy analyzer scope bound small
  scope hypothesis bounded verification within scope" — MISS (all
  irrelevant; rejected).
- ⚠ no prior-art coverage for the instance question (peer state-machine
  linters' too-large disposition) in the resolving corpora. The
  model-prior claims about Alloy scopes and BMC completeness thresholds
  are therefore NOT load-bearing in the record; the definitional-scope
  stance rests on the quoted TraceFix page-9 passage and on 0003's own
  D6/REQ-86 text instead, and the broader peer-tool comparison is demoted
  to a Resolve assumption (Method: Prior Art).

## Query ledger

| # | Corpus | Query | Outcome |
| - | ------ | ----- | ------- |
| 1 | StateMachineLit | bounded model checking bound depth completeness threshold analysis scope | 2 accepted (TraceFix p7, p9) |
| 2 | PapersFast | test oracle problem pseudo-oracle differential testing independent oracle | miss |
| 3 | DevRef | test oracle must be independent of the implementation under test self-verifying tautology | 1 accepted (developer-testing p171) |
| 4 | StateMachineRes | state space too large to enumerate bound refusal explicit limit exhaustiveness checking | miss |
| 5 | StateMachineRes | Alloy analyzer scope bound small scope hypothesis bounded verification within scope | miss |
