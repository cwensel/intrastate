# RDR 0022 — Stage 2 prior-art search record

Date: 2026-08-28. Budget: ≤3 corpus queries per claim, ≤5 opened hits per
claim.

## Claim 1 — peer state-machine tools' disposition on terminal reachability
(instance question: what does each peer validator do for "state cannot reach
a final state"?)

Queries (StateMachineRes, semantic + text):

1. `arc search semantic --corpus StateMachineRes "validator check that every
   state can reach a final state livelock unreachable terminal"` — top 5 hits
   all Archon workflow-process docs (validate-phase checklists), none about
   graph validation. Rejected.
2. `arc search semantic --corpus StateMachineRes "liveness property terminal
   reachable from every state TLC model checker recommendation lint"` — hits
   are awf-cli code-quality linter docs and statewright model-compatibility
   spec; none address terminal reachability. Rejected.
3. `arc search text "livelock" --corpus StateMachineRes` — 0 results.
4. `arc search semantic --corpus StateMachineRes "detect state that can never
   reach a final or terminal state; statechart validation unreachable final"`
   — hits are Archon validate phases and scxmlcc's `target` attribute doc;
   none implement or specify the check. Rejected.

**Result: ⚠ no prior-art coverage** in the indexed peer-tool corpus for a
terminal-reachability lint; no peer checkout indexed here ships one. The
instance disposition therefore rests on the in-repo prior art (RDR 0006's
invariant 2/7 machinery) and on the review-sourced framing already fixed in
the Problem Statement.

## Claim 2 — literature: AG EF terminal as the standard livelock property,
computed by backward reachability

Queries:

1. `arc search semantic --corpus StateMachineLit "backward reachability from
   goal states liveness AG EF model checking"` — top hits all
   TraceFix (TLA+ agent-protocol repair, CAIS 2026); confirms TLC-class
   checking exists as prior art but contains no quotable statement of the
   EF-fixpoint algorithm. Rejected as load-bearing.
2. `arc search text "AG EF" --corpus StateMachineLit --corpus PapersFast` —
   noise (tokenized matches); no quotable hit.
3. `arc search semantic --corpus StateMachineLit "CTL operator EF least
   fixpoint backward image computation model checking algorithm"` — scores
   ≤0.49, no algorithmic source. Rejected.

**Result: ⚠ no quotable prior-art coverage** for the CTL algorithm claim in
the local corpora. Per the propose prompt, order/reachability claims are
never quote-confirmed at this stage anyway: the Clarke/Emerson/Sistla 1986 /
Newcombe 2015 framing stays a Resolve assumption (Method: Prior Art), not a
load-bearing Stage-2 citation. The choice does not rest on it — it rests on
in-repo quotable prior art.

## Accepted (in-repo, quotable) prior art the choice rests on

- `internal/graphlint/analysis.go::checkDeadEnd` — "The test runs on SPLIT
  nodes, not merged ones. Universal satisfaction is anti-monotone in
  merging… The split is exact… bounded by the declared domains of the
  terminal-participating keys ONLY, never the whole lattice."
- `internal/graphlint/analysis.go::hasOutgoingOrdinaryRow` — the pinned
  residual false-negative doctrine ("the accepted false-NEGATIVE the record
  books against invariant 2 … See D12").
- `0006:D-soundness-direction-is-per-invariant-not-global` — universal
  checks must not read merged nodes; dead end's remedy is the split.
- `0006:ALT4` — external model checker rejected: "the authoritative gate
  should be native lint over the model intrastate actually consumes."
- `0006:C3` — blocking tier is a floor ("at least these blocking invariant
  classes"); `0006:C17` — advisory tier closed at four.
- `internal/graphlint/taxonomy.go::ReasonNoParticipatingDimension` — RDR
  0010's append precedent on a closed-but-append-only vocabulary.
- `0017:C1` (PROPOSED peer) — findings[].code drawn from a producer-declared
  per-finding vocabulary; a later-added finding class conforms by extending
  its producer's declared vocabulary (0017:A5).
- `0024:§approach` (PROPOSED peer) — emit declarations land in the load
  pipeline; emit values "stay uninterpreted", never enter graphlint — no
  interaction with the owned-state relation this RDR traverses.
