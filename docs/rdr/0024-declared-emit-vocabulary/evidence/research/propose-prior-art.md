# RDR 0024 — Stage 2 prior-art search record

Stage: Propose (selection reads, not spikes). Budget: ≤3 queries and ≤5
opened hits per claim.

## Claim 1 (problem class): peer decision-table systems validate output
entries against a declared allowed-value set

Queries run (all `arc search semantic --limit 5 --json`):

1. `--corpus StateMachineRes` — "decision table output values allowed
   values validation DMN output entries" → top hits generic validation
   docs (BehaviorTree.CPP name rules, inngest parameter validation);
   none about output-vocabulary declaration. Rejected.
2. `--corpus PapersFast` — "decision table completeness output actions
   validation" → TOC pages and unrelated ML papers. Rejected.
3. `--corpus StateMachineLit` — "decision table verification anomaly
   detection completeness DMN" → scores ≤0.42, agent-protocol papers.
   Rejected.

⚠ no prior-art coverage for the DMN "output values / allowed values on an
output clause" class claim in the available corpora. The class claim is
therefore NOT load-bearing at Propose; the DMN alignment note is demoted
to a Resolve assumption (spec read), and the choice rests on the in-repo
mirror instead: `0002:C22` tag type-model declaration + RDR 0003 kind
vocabulary (`internal/table/model.go::declaredKinds`), which is openable
and quoted in the RDR body.

## Claim 2 (instance): peers mark a stop/terminal answer structurally
(a declared marker), not by string naming convention

Query run:

1. `--corpus StateMachineRes` — "terminal state final state marker
   distinguish continue versus halt result token" → 5 hits, 2 opened:

Accepted citations:

- ms-conductor `examples/README.md` §Explicit Termination
  (`/Users/cwensel/sandbox/newcoinc/state-machines/repos/ms-conductor/examples/README.md`):
  a workflow with "multiple legitimate end states" uses `type: terminate`
  steps with a declared `status: success|failed` field — the stop is a
  structural declaration carried on the step, surfaced distinctly
  (exit code, `is_explicit: true`), never a prefix on the output string.
- scxmlcc `doc/user-manual.md` §Final State (`<final>`)
  (`/Users/cwensel/sandbox/newcoinc/state-machines/repos/scxmlcc/doc/user-manual.md`):
  SCXML marks termination with the `<final>` element — a declared node
  kind that machinery reacts to (`done.state.ID` events, machine
  termination) — not a naming discipline over state ids.

Rejected branches: xgrammar structural-tag API (token-format grammar,
not a stop marker), Archon CLI quick reference and conductor CLI docs
(command lists, no disposition semantics).

Verdict: the instance read confirms the declared-disposition shape over
the reserved-prefix convention; search stopped at confirmation.
