Model: claude-opus-5

# Charted — 3amigo findings not resolved into this RDR

- **PM-5 — `§decision-rationale`'s Premortem/Ground-sweep/Joint-check ledger
  crowds out the rationale.** Out of scope for this lens: the block is durable
  provenance the Finalization Gate reads (`G-proportionality` explicitly asks
  whether provenance prose should be stripped), and the same content lives in
  `evidence/joint-check/arms.md`. Trimming it is a proportionality judgement at
  Stage 7, not a 3amigo defect fix — and cutting provenance the gate has not yet
  read would be the wrong order. Successor: answer it in the `G-proportionality`
  gate response at `/rdr-finalize`, where the ledger's per-arm dispositions are
  either summarized to their verdicts or left whole.

- **iter-2 delta, inherited — the `${2}` group reference in the `fidelity` and
  `trace` mini-check tables vs. an anchor defining one group.** Not introduced by
  this lens's fixes: it rides in from A3 (Status: Verified) and the tables the
  grounding pass wrote from it. Out of scope for 3amigo, which reviews the draft
  rather than re-opening a Verified assumption's spike. Successor: check at
  `/rdr-reconcile` or the Finalization Gate, alongside A10/A11 — if A3's anchor
  defines one capture group, the tables' `${2}` should read `${1}`, or the
  anchor in the illustrative/normative form should define two.
