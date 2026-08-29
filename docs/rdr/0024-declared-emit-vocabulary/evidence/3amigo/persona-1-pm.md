Model: claude-opus-5[1m]

# Persona 1 — Product Manager: does 0024 deliver the user outcome?

Starting set: `§problem-statement`, `§approach`, `§decision-rationale`, `MVV`.
**Widened** to `§implementation-plan` (Phase 4), `§consequences`, `§risks-and-mitigations`,
`§testing-strategy`, `§capability-dependencies`, `A5`, `A3`, `C4`, `§finalization-gate`
sub-sections. What sent me: the Problem Statement names a harm that lands on a
*consumer* ("a verbatim executor executes it", "a stop printed as a command to run
is a skipped check reading as a passed one"), but every contract and every phase
lands inside intrastate. A product question about whether the harm is actually
removed cannot be answered inside the four owned entries — the answer is a silence
in the plan, and a silence has no line range.

---

## PM-1 — The stop/route harm is not closed by anything this RDR ships

**Anchor**: `0024:§problem-statement` (harm statement) vs `0024:C4` + `0024:§consequences`
**Severity**: high

The Problem Statement states the folded facet's harm in consumer terms: today
"the distinction lives in a string-prefix naming discipline every consumer
re-implements and no lint checks." `C4` replaces that with `dispositions` on the
resolve payload — but the token is "surfaced verbatim; intrastate never interprets
it," and `§risks-and-mitigations` (disposition-token drift) puts the *meaning* of
`stop` back on "the consumer pins its own token vocabulary with a contract test
against its model." So after this RDR the consumer still re-implements a
discipline; the discipline has moved from a string *prefix* to a string *token in
a different payload field*. `§consequences` claims "consumers stop re-implementing
a string-prefix discipline" — that is a substitution, not a removal, and the record
does not say what the consumer's code does differently on the day `dispositions`
arrives. Nothing in `§implementation-plan` produces a consumer-visible artifact
that makes the difference, and `§testing-strategy` scenario 5 asserts only the wire
shape.

**Blocks**: the decision of whether `C4` earns its blast radius (28 golden
rewrites in Phase 3 + an amendment to a peer RDR's spec test `TestReq39`). If the
delivered outcome is "one more model-authored string in a different field," the
cheaper `ALT1` scoring of 2 on correctness is contestable on the *disposition*
half — `ALT1` loses on the vocabulary proof, which is a different argument. Also
blocks sizing: `C4` could be deferred to a follow-on with no loss to the
Problem Statement's *primary* (kata `vt9n`) harm.

---

## PM-2 — The motivating consumer's model is not in this repo, and no phase reaches it

**Anchor**: `0024:§phase-4-example-and-docs` (widened; contrast `0024:A5`)
**Severity**: high

`A5` grounds the entire disposition design on `rdr-status.toml` (54 emit blocks,
22 distinct `next` values, the three-valued partition) and `rdr-write.toml` (29
blocks, `op` key, `edit` needing `scalar`). Neither file exists in this repo
(`find` over the source tree returns nothing; `models/` holds `rdr.toml` and three
examples). Phase 4 declares emit keys for `models/examples/pricing-decision-table.toml`
instead — whose emit keys are `plan` (`basic|pro`) and `dpa` (`required|none`), a
model with **no stop, no terminal, and no route**. So the shipped worked example
exercises `C1`'s flat-domain leg and demonstrates the disposition partition
*nowhere*. `§implementation-plan` has no phase, and `§testing-strategy` no
scenario, that authors the motivating model's declaration; `A5`'s three-valued
finding — the empirical case that killed `ALT1` — is never exercised end to end
by anything this RDR delivers.

**Blocks**: the go/no-go on the disposition sub-table grammar. The one live
justification for `[emit.<key>.domain]` over the flat `domain = [...]` array is a
model in another repo that this plan never touches. Also blocks the "adoption
step" question in `§consequences` ("whole-model strictness makes adopting the
first declaration a step") — the cost of that step is measured only against
`rdr-status.toml`/`rdr-write.toml`, and the plan does not pay it, so the estimate
is never validated.

---

## PM-3 — Opt-in with no adopter named means the shipped outcome is zero

**Anchor**: `0024:§approach` (opt-in framing) + `0024:§consequences` (first Negative)
**Severity**: high

`§approach` and `C2` make the whole guarantee opt-in; `A3` verifies that **zero**
models, fixtures, or test sources in the repo declare an emit vocabulary; `§consequences`
concedes "the guarantee exists only where an author spends the declaration effort."
Phase 4 adds one declaration to the pricing example. Net user-facing outcome at
merge: the pricing example gains a declaration nobody depends on, and every real
consumer (the RDR-navigator models in the sibling repo) still ships the exact
defect the Problem Statement reproduces at `bc9f2a0`. The record never states an
adoption commitment, a target adopter, or a "done when" for the *outcome* as
distinct from the *mechanism*. `§testing-strategy`'s "Done =" is entirely
mechanism-side ("every scenario green and the MVV run recorded").

**Blocks**: the Priority `High` claim in `§metadata` and the phase ordering. If no
adopter is scheduled, this is infrastructure ahead of demand and Phase 3/Phase 4
sequencing is arbitrary. Also blocks the reviewer's ability to falsify the
recommendation — there is no stated post-merge observation that would show the
Problem Statement's harm was actually removed.

---

## PM-4 — Fail-fast one-finding-per-run is a real authoring-experience regression against the stated outcome, and the RDR only argues it is *permitted*

**Anchor**: `0024:C2` (fail-fast clause) + `0024:MVV` step 2
**Severity**: medium

The user-side promise in `§problem-statement` is that lint proves the table — the
author's bargain. `C2` spends substantial space establishing that one refusal per
run is *inherited* from the load tier and *defensible*, and concedes "aggregation
would be defensible on this channel and is declined because the tier it lands in
has already fixed the opposite." `MVV` step 2 then bakes the consequence into the
acceptance walk: "Fix that rule and re-run: the second defect now surfaces." For
the first adopter — `A5`'s model with 54 emit blocks and 22 values across two
routing keys, adopting whole-model strictness in one step per `§consequences` —
first-declaration adoption is an N-round fix-and-rerun loop where N is the number
of typos, with the pipeline free to interleave unrelated structural refusals in
unspecified order. The record never states this cost from the author's seat, never
sizes N for the named adopter, and `§consequences` does not list it as a Negative.

**Blocks**: the enforcement-point decision (load vs. a lint-side tier). The
`0002:C24` arity argument in `§approach` establishes that load is the *correct*
tier, but the product question — is a one-at-a-time authoring loop acceptable for
the bulk-adoption event this design makes mandatory — is never asked. If the
answer is no, either the whole-model strictness in `C2` or the tier choice is the
thing to revisit, and both are decided here.

---

## PM-5 — `scalar` gives the adopter a one-line way to take the credit and none of the proof, and the only stated defense is "a reviewer will notice"

**Anchor**: `0024:C1` (`scalar` arm) + `0024:§risks-and-mitigations` (scalar hollowing)
**Severity**: medium

`C1` makes `scalar` a per-key escape; `C2` never refuses a scalar value. The
Risks entry names the exact failure ("whole-model strictness pushes a bulk adopter
to declare everything `scalar`, and 'declared' reads as 'checked'") and mitigates
with (a) the declaration table being reviewable and (b) an explicit refusal to add
the advisory finding, because `0006:C17` is closed. `§decision-rationale` restates
this as a structural refusal of Meyer's own attached condition ("identify [its
loopholes] clearly, if possible providing tools to flag any software using them").
The honesty is good; the product consequence is that the RDR ships a guarantee
whose most likely adoption path is the one that provides no guarantee, with human
review as the only control — and `A5` confirms the named adopter *already* has a
key (`edit`, a multi-line shell script) that must be `scalar`, so the escape is
exercised on day one. Nothing in `§testing-strategy` or `MVV` makes scalar
coverage observable; `§consequences` does not list hollowing as a Negative.

**Blocks**: whether whole-model strictness (the `C2` clause forcing every key
declared at once) is the right coupling. Per-key opt-in would remove the pressure
toward `scalar`, at the cost of not catching the seed's `nxet` key typo — which
`C2` states plainly. That is a genuine product tradeoff the record resolves in one
sentence ("per-key checking cannot catch the seed's `nxet` defect") without pricing
the hollowing it creates.

---

## PM-6 — Author discoverability of the new grammar is unplanned

**Anchor**: `0024:§phase-4-example-and-docs`
**Severity**: medium

Phase 4 extends `docs/cli-output-contract.md` — a *payload* document — with "the
declaration grammar." `docs/model-authoring.md` is the author-facing home and is
quoted approvingly twice in `§decision-rationale` ("a misspelled facet is a stable
refusal rather than a silent no-op" and "declaring the class rather than deriving
it is what makes the refusal possible" — called "this RDR's argument in the
project's own words"). It is not in any phase's edit list. `§capability-dependencies`
has no row for author documentation. For an opt-in feature whose entire value
depends on an author choosing to write a declaration, the grammar is documented in
the wrong book.

**Blocks**: nothing structural, but it blocks the adoption path PM-3 identifies
as missing — an author cannot opt in to a grammar documented only in the output
contract.

---

## PM-7 — Success is defined only mechanically; there is no outcome-side "done"

**Anchor**: `0024:MVV` + `0024:§testing-strategy` ("Done =")
**Severity**: medium

Every `MVV` step and every Testing Strategy scenario asserts an intrastate-internal
value: exit codes, finding categories, a 14-key wire order, `NumField() == 15`.
The Problem Statement's harm is stated in outcome terms ("a model can route
perfectly to a command that does not exist, and a verbatim executor executes it";
"a rare row's typo ships and fires at a close-out"). Neither the MVV nor any
scenario asserts that a real routing model's real typo is now caught — the MVV's
table is a synthetic two-rule adversarial fixture authored for the test. The
`§risks-and-mitigations` domain-drift entry twice defers the outcome-closing check
to a consumer-side seam test ("every declared emit domain member is a real command
or stop token") that is explicitly out of scope and unscheduled.

**Blocks**: the Scope Verification gate (`0024:G-scope`), which asks whether the
MVV "is in scope and will be executed" and to "state the specific test or proof."
The MVV proves the mechanism works. It does not prove the outcome arrives, and the
one check that would is assigned to a repo this plan does not touch.

---

## PM-8 — Finalization Gate is unfilled template; the outcome claims are unadjudicated

**Anchor**: `0024:G-contradiction`, `0024:G-assumptions`, `0024:G-scope`, `0024:G-proportionality`
**Severity**: low (status is `Draft`; flagged as a lock blocker, not a design defect)

All four gate sub-sections still carry template instruction text, not written
responses. `G-proportionality` in particular is the gate that would have to
re-derive the `foundational` Profile against a contract count of four (`C1`–`C4`)
spanning `internal/table` and `internal/cli` plus a user-facing payload change —
and PM-1 argues `C4` is a separable seam. `G-assumptions` would have to record
that `A6` is carried Pending (`§prerequisites` does state this, so the material
exists but is not in the gate).

**Blocks**: lock. Also blocks the split question PM-1 raises: `G-proportionality`
is the designated place to decide whether the envelope contract (`C4`) and the
vocabulary proof (`C1`–`C3`) are one seam or two, and it is currently silent.

---

## PM-9 — The `dispositions: {}` field on every undeclared model is a user-visible change with no user benefit

**Anchor**: `0024:C4` (never-omitted clause) + `0024:§approach` ("its one observable delta")
**Severity**: low

`C4` requires `dispositions` "present as `{}` — never `null`, never omitted" even
for a model with zero declarations. `§approach` frames this as "the same class of
additive, fixture-updated change `0010:C4` shipped." For the 100% of current
models that will not declare (`A3`), the entire shipped user-visible outcome of
this RDR is a permanently empty JSON object in every resolve payload, plus a
28-golden rewrite. `omitempty` is available (`C4`'s own text notes `escape_class`
uses it and is absent from the wire) and is not weighed; the never-omitted choice
is stated as contract without a rationale in `§decision-rationale` or
`§load-bearing-decisions` — only the *position* (after `emit`) and the *name* are
given rejected alternatives.

**Blocks**: nothing, if the always-present shape is a deliberate consumer-parsing
guarantee — but that reason is not written down, so a reviewer cannot tell
whether the 28-golden cost buys a stable-shape promise or is an unexamined default.
