Model: claude-opus-5[1m]

# Charted to successor — 3amigo iteration 2, RDR 0007

Real findings that are net-new scope for this RDR. Each is recorded here and
dismissed from this loop; none was absorbed into the draft (the
scope-expansion wormhole).

- **PM-4 — the aggregation veto is unsized and has no dry-run affordance.**
  One unevaluable row refuses a whole table; the RDR accepts this explicitly
  ("This veto is the cost the RDR accepts") but offers operators no way to ask
  "which rows would refuse against today's artifact state?" before it bites.
  Why out of scope: a diagnostic/dry-run mode is a CLI surface, and this RDR is
  the kernel's domain rule — RDR 0005 owns verb surfaces and RDR 0006 owns the
  lint/proof story. Suggested successor: an RDR seed against RDR 0005 for a
  guard-decidability dry-run, or an extension of RDR 0006's lint promise
  (already narrowed at §JD-4). Substantial enough to warrant `/rdr-seed`.

- **PM-5 — the `unless`-over-a-rarely-set-tag footgun is routed to authoring
  prose with no gate.** Failure Modes concedes "skip when `legal_hold` is true"
  refuses on every artifact that never had `legal_hold`, and hands the
  conjoined-existence idiom to RDR 0003's authoring guidance in Phase 4. PM's
  objection is sharp: prose-only mitigation is the same "held by discipline"
  weakness this RDR rejects Alternative 1 for. Why out of scope: making it a
  *lint* obligation is RDR 0003's (operator vocabulary) and RDR 0006's (what
  lint may promise); this RDR cannot impose a load-time rule on a Draft peer's
  grammar. Suggested successor: fold into RDR 0003's re-entry as a lint rule
  request — "a value atom in `unless` over a tag with no existence sibling is a
  lint warning." Phase 4 already carries the handoff; the *gate* is the new ask.

- **PM-3 — no scenario validates the Problem Statement's "told plainly".**
  All 20 Testing Strategy rows and all 4 MVV rows assert on Go struct fields;
  the user-visible refusal text is never exercised. Why out of scope: the CLI
  rendering is RDR 0005's and is blocked on the A19 envelope question; this RDR
  correctly declines to specify the renderer. Suggested successor: an
  acceptance scenario on RDR 0005 that asserts the rendered refusal names the
  absent key — the end-to-end closure of kata `xg7p`. Note this is the finding
  that says `xg7p` cannot fully close on this RDR's lock alone.

- **QA-6 — the mutation MUST names no tool, mutant set, or threshold.** The
  Testing Strategy says any future edit to the combination tables "MUST be
  re-mutation-tested," which is an unexecutable instruction without naming how.
  Why out of scope: choosing a Go mutation-testing tool and wiring it to CI is
  repo tooling, not this RDR's contract. Suggested successor: a kata for
  mutation-testing infrastructure under `internal/resolve`, with the
  strong-Kleene combinator as its first target (the spike's 6-case probe at
  `evidence/spikes/a3-reshape/spike_kleene_test.go.txt` is the seed).

- **IMP-8 / QA-13 — "no short-circuit" makes seam call-count an observable
  contract on RDR 0003's code.** If the kernel MUST evaluate every atom of
  every survivor, then a conforming evaluator must tolerate being called for
  atoms whose block is already decided — an obligation on 0003's evaluator that
  this RDR states only implicitly, and one with no payload-visible consequence
  beyond row 20. Why out of scope: it is a clause on the *seam contract* that
  Phase 3's exported contract-test function should carry, and Phase 3 is
  already blocked on 0003's element-encoding declaration. Suggested successor:
  add "the evaluator MUST be callable for any present-key value atom regardless
  of sibling verdicts" to Phase 3's contract-test export when that unblocks.
