# RDR 0015 — Stage 2 prior-art search record

Date: 2026-08-28. Budget: ≤3 corpus queries per claim, ≤5 opened hits per
claim. Sibling 0022's committed search record
(`docs/rdr/0022-terminal-reachability-liveness-invariant/evidence/research/prior-art.md`)
already covers the shared problem class (merged-node liveness quantifier);
its negative results are reused, not re-run.

## Claim 1 — scoped relational tracking ("packing") is the standard remedy
for correlations lost by attribute-independent abstraction
(frames Alternative 2's costing)

Queries:

1. `arc search semantic --corpus PapersFast --limit 5 "relational abstract
   domain variable packing recover correlations lost by cartesian
   non-relational abstraction"` — top 5 hits are citation-fragment chunks
   from unrelated database/data-mining papers (scores ~0.80 on reference
   lists, no content). Rejected.
2. `arc search semantic --corpus StateMachineLit --limit 5 "abstract
   interpretation precision relational domain octagon packs Astree
   correlated variables"` — scores ≤0.49, all one Armada (Donaldson 2020)
   refinement-proof paper; no packing/relational-domain statement. Rejected.

**Result: ⚠ no prior-art coverage** for the relational-packing framing in
the local corpora. The Cousot/Astrée-style "restrict relational precision
to small packs of correlated variables" framing stays a Resolve assumption
(Method: Prior Art) and is NOT load-bearing at Stage 2 — Alternative 2's
rejection rests on the in-repo measured evidence (empty residual
population, D12's costed widening), not on this citation.

## Claim 2 — peer state-machine tools' disposition for dead/trap-state
detection over an abstracted state space (instance question)

Queries:

1. `arc search semantic --corpus StateMachineRes --limit 5 "dead state trap
   state detection state machine verification tool concrete state
   enumeration"` — hits are statewright/javascript-state-machine/jfsm doc
   headers; none specify a dead-state check's quantifier. Rejected.

Reused from 0022's record (same problem class, 4 queries, all rejected):
**⚠ no prior-art coverage** in the indexed peer-tool corpus for a
liveness-shaped lint at all, let alone its quantifier over an abstraction.

## Accepted (in-repo, quotable) prior art the choice rests on

- `docs/rdr/0006-graph-lint-authority-and-guarantees/artifacts/req-list.md`
  — REQ-37: "The split is bounded by the declared domains of
  terminal-participating keys only, never the whole lattice." (MUST);
  REQ-111: "**Universal checks must not.**" … dead end's remedy is to
  "**split** the node on terminal-participating keys first, recovering
  exactness"; REQ-117: "the cure is a clearer model, never a weaker lint";
  REQ-122: SC-23 census, "Zero is the pass condition".
- `docs/rdr/0006-graph-lint-authority-and-guarantees/artifacts/deviations.md`
  D12 — the measured widening experiment: seven correlated
  `stage = "dropped"` / `status = "abandoned"` rows in `models/rdr.toml`;
  widening "mints NINE false `graph-dead-end` findings on the model REQ-122
  requires lint clean"; "Both readings have record support and the record
  does not rank them"; "This entry is the trigger."
- `internal/graphlint/analysis.go::checkDeadEnd` /
  `::hasOutgoingOrdinaryRow` — the shipped quantifier and its documented
  residual ("the accepted false-NEGATIVE the record books against
  invariant 2 … See D12").
- `internal/graphlint/adversarial_0006_test.go::TestAdvDeadEndExistentialOnMergedNode`
  — the pin: asserts the bounded split's guarantee, pins the miss, carries
  a positive arm on a terminal-participating key.
- `docs/jdr/0001-resolve-kernel-seam.md` §JD-23 — the doctrine's home;
  charters this RDR's ranking; "neither record widens or narrows the split
  unilaterally."
- `cli/0022:C2` (PROPOSED peer, via projector) — invariant 8 adopts the
  same split bound and accepts the path-scoped D12-class residual, citing
  §JD-23.
