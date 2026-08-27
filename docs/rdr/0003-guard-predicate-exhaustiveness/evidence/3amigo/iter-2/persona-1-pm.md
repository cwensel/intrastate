Model: claude-opus-5[1m]

# 3 Amigos iter-2 — Persona 1: Product Manager

Delta-scope: A5, A7, A8, and the passages stating or depending on the narrowed
lint promise (Technical Design narrowing paragraph, Normative Contracts
exhaustiveness-narrowing clause, Decision Rationale, Contradiction Check,
Prerequisites, Minimum Viable Validation).

Lens: does this RDR deliver the user outcome? The stated outcome (Problem
Statement) is that "a flow author needs conditional edges ... to be expressed in
a way the lint can prove exhaustive and mutually exclusive." The re-entry's new
material is about what lint may NOT claim. That is the right substance, but the
author-facing half of it is largely unwritten.

---

## HIGH

### P1-1 — The narrowed promise has no author-visible outcome: "withheld" is defined only as the absence of a green result

**Anchor**: Normative Contracts, the exhaustiveness-narrowing clause — "lint
MUST NOT certify a row group exhaustive when a participating row can refuse
`guard_unevaluable` under RDR 0007's aggregation veto"; and the `disposition`
table's last row, "Row group is domain-exhaustive but a participating row can
refuse → Lint outcome: `Claim withheld — the narrowing`; Diagnostic minted:
`RDR 0007 payload; no green certification`; Silent or loud: `Loud`."

**Objection**: The disposition table marks this case **Loud**, but the only
diagnostic it names is "RDR 0007 payload" — which is a *runtime refusal*
payload, minted when resolution actually runs, not a lint finding. At lint time
nothing has refused yet; the author has run `lint` on a table, not resolved a
flow. So the row asserts a loud lint outcome while naming no lint diagnostic,
and the column literally reads "no green certification" — the absence of a
result, not the presence of one.

The neighbouring rows in the same table each name a real lint code:
`graph-coverage-gap` (RDR 0006), the overlap code (RDR 0006), the
inability-to-prove finding for the non-finite and too-large cases. This new row
is the only Loud row with no code. And it cannot borrow the existing one: RDR
0006's `graph-unprovable-coverage` is bound in that document's finding table to
"Required finite-domain proof unavailable"
(`./docs/rdr/0006-graph-lint-authority-and-guarantees.md:302`)
and its Testing Strategy scenario 3 scopes it to "a model claims closed coverage
over an input dimension that is **not finite** under the predicate contract"
(same file, lines 696–699). The narrowing case is by construction the opposite:
the product IS fully finite. There is no code for it in either document.

Consequence for the outcome: the author whose table hits this case gets — per
the text as written — a lint run that does not say the group is exhaustive and
does not say why. That is indistinguishable from lint simply not having checked.
The user outcome in the Problem Statement is that lint can *prove* things about
the author's edges; a silent non-proof delivers neither the proof nor the reason
it is unavailable.

**Decision blocked**: Whether RDR 0006 needs a *new* blocking finding code for
the withheld-claim case (versus widening `graph-unprovable-coverage`'s
definition), and whether the withheld claim is blocking or advisory. Both are
part of the A8 route-back that Prerequisites gates lock on, and A8 as written
asks RDR 0006 only to "adopt the same wording" — it does not ask for the code.
The route-back will be scoped too narrowly to fix this.

---

### P1-2 — The MVV asserts what lint must NOT do, and never asserts the author can act on it

**Anchor**: Minimum Viable Validation — "The fixture must also include one row
group that is domain-exhaustive yet contains a possibly-absent guard key, and
assert that lint withholds the exhaustiveness claim there rather than
contradicting the runtime refusal." And Testing Strategy Scenario 6 —
"**Expected**: ... lint does not certify that row group exhaustive — the
green-lint-implies-resolution-succeeds promise holds."

**Objection**: Both are negative assertions. Neither asserts that anything is
emitted, that a code is returned, or that the offending key/row is named. This
matters against this RDR's own standard elsewhere: the MVV's earlier sentence
requires "one intentional gap and one intentional overlap ... **with source
rule/context ids in the diagnostic**," and Normative Contracts require "Overlap
and coverage diagnostics MUST name the source rule id or context id that
contributed each predicate involved in the finding." The withheld case is
exempted from that standard without saying so — it is a coverage-class outcome
whose diagnostic is never required to name anything.

A test that passes when lint emits nothing at all is satisfied by an
implementation that simply never runs the check. The MVV is the validation the
Scope Verification section calls "the specific proof"; here it proves only that
a false claim is not made, not that a true and useful one is.

**Decision blocked**: What implementation Phase 2 must actually build for this
case — a suppression of the green result, or a finding with a payload. As
written, Phase 2 ("Define how ... domains are converted into scoped row-group
coverage and overlap checks, including refusal/downgrade behavior for unbounded
dimensions and finite products too large to prove") does not mention the
narrowing case at all, so an implementer reading Phase 2 plus the MVV builds the
suppression and stops.

---

## MEDIUM

### P1-3 — A7's `Pending` status is scoped to a projection rule, but the same gap silently removes the RDR's headline capability for the pattern RDR 0007 sanctioned

**Anchor**: A7 "If wrong" — lint "refuses every group containing an `exists`
atom, which disqualifies this RDR's own representative fixture row
`foundational-to-cove` and leaves RDR 0007's sanctioned two-row absence pattern
permanently unprovable." And the `trace` table Step 5 — "**GAP — booked as A7,
not a contradiction.** Without the presence dimension this step has no defined
result."

**Objection**: A7 is framed and severity-scoped as a *representation* question
(how does an atom project onto a product). Read against the user outcome, it is
a scope question: whichever way it lands, one of the two arms deletes a
capability the rest of the document sells. The `exists` operator is in the
closed operator matrix, is one of the five operators the A1 spike ran, and is
listed in the MVV. If A7 resolves to its P2 arm, `exists` remains a legal
runtime operator whose presence in any row makes that row's whole group
unprovable — i.e. authors of the sanctioned pattern get no exhaustiveness
guarantee at all, which is the single thing this RDR exists to give them.

The document never states this consequence at the level where a reader picks up
scope. Proportionality says "The document is right-sized for lock. It owns one
independent load-bearing contract"; Trade-offs / Consequences says "Static lint
can produce gap and overlap diagnostics from the same predicate model runtime
resolution uses" with no carve-out; the operator matrix lists `exists` beside
`eq`/`in` with no note that its provability is unsettled. Only A7's own "If
wrong" carries it.

**Decision blocked**: Whether the RDR can lock with `exists` in the closed
operator vocabulary at all, or whether `exists`-bearing groups need a declared,
documented exhaustiveness carve-out in Trade-offs. That is a scope call for the
user-facing capability, and it is currently only reachable by reading an
assumption's failure branch.

---

### P1-4 — A7's plan depends on a producer request that no document has accepted, and the dependency runs to a `Draft` peer

**Anchor**: A7 **Plan** — "state the projection rule as a normative clause
before lock, and carry the RDR 0002 producer request it depends on (below).
Blocked on that request: RDR 0002's tag declaration carries name, provenance,
value kind, and optional accessor reference, with no optionality field ... RDR
0002 is itself `Draft`." And Capability Dependencies, last row —
"Per-tag optionality declaration (which keys may be absent) | RDR 0002 |
**Requested**".

**Objection**: Verified against
`./docs/rdr/0002-transition-table-as-reviewable-data.md:218`
— the tag declaration is "value kind, and optional accessor reference for
observed or owned read-back," with no optionality field, exactly as A7 says. The
claim is accurate. The problem is the plan's shape: `Requested` is a status this
RDR assigns to itself. Nothing in RDR 0002 records the request, and RDR 0002 is
`Draft` with its own re-verify list (`Status`: "re-verify A2, A7"). So A7's plan
has no counterparty who has agreed to it, and Prerequisites states the block as
"**A7 is Pending**, blocked on RDR 0002's optionality declaration" — a
dependency on a document that does not know it owes anything.

Compare A8, which the same Prerequisites line handles better: it names the
mechanism ("blocked on a route-back to `Final` RDR 0006"). A7 names no
mechanism.

**Decision blocked**: Sequencing — whether RDR 0003 can lock before RDR 0002
re-locks, or whether the two must lock together. Prerequisites currently implies
0003 waits on 0002, but no artifact makes 0002 aware, so the wait has no
termination condition.

---

### P1-5 — Decision Rationale claims the outcome is delivered, in a paragraph that sits directly above two open blockers to it

**Anchor**: Decision Rationale — "The fixed symbolic atom model best matches the
user's outcome: flow authors can write conditional edges, and lint can still
prove whether those edges are complete and mutually exclusive."

**Objection**: This is the document's one explicit statement of the user
outcome, and it is stated unconditionally ("can still prove"). At this revision
that is not what the rest of the document says. Per A7, groups containing an
`exists` atom have no defined proof result (`trace` Step 5: "this step has no
defined result"). Per the narrowing this re-entry adopts, domain-exhaustive
groups containing a possibly-absent value key are withheld from proof. Those are
two named classes of conditional edge for which lint *cannot* prove completeness
as the document currently stands.

The Contradiction Check is more honest ("The strength of the exhaustiveness
claim is decided in substance but not yet agreed across documents"), and I am
not filing this as a contradiction — the paragraph is comparative rationale for
choosing Alternative 1 over the alternatives, and against callbacks it holds.
But it is the sentence a reader takes as the promise, and it now overstates the
delivered scope by exactly the two Pending assumptions.

**Decision blocked**: How the capability is communicated at lock — whether the
selling sentence needs the qualifier the operator matrix and Trade-offs also
lack (P1-3). Without it, downstream readers (RDR 0006 consumes this RDR's
exhaustiveness semantics; RDR 0005 maps its findings to the envelope) inherit
the unqualified promise.

---

## LOW

### P1-6 — A5 is re-verified as a citation change, but its user-facing consequence ("If wrong") is no longer the risk the evidence addresses

**Anchor**: A5 Evidence — "The normalized atom shape is fixed at the kernel seam
by JDR 0001 §D1; RDR 0007 is the landing document and states it normatively — a
row carries parsed atoms (key, operator token, literal, block), not an opaque
predicate string, so there is no reconstruction step that could lose identity."
And A5 **If wrong** — "Lint may detect an error but fail to point reviewers at
the guard to fix."

**Objection**: The re-verification did its job: the shape is cited, not
restated, and the source-identity claim is anchored to shipping code
(`internal/resolve/resolve.go::Row` carrying `RuleID` and `SourceLocator`) —
which I take as sound. But the Evidence also records "RDR 0007 is `Final`, not
yet implemented: the shipped kernel still carries `Row.Guard string`, and
0007's A3 scopes the reshape." So the identity-preserving shape is *specified*
and the shipped shape is the opaque string the assumption is about. The
assumption reads `Verified` on the strength of a peer document's unimplemented
contract.

That is defensible under `Method: Peer RDR`, and Prerequisites already sequences
this RDR after 0007's reshape. It is Low because the sequencing is stated. But
"If wrong" describes a user-visible failure (reviewers cannot find the guard to
fix) whose actual trigger is now "0007's reshape does not ship as specified" —
not anything this RDR or RDR 0002 controls. The failure branch names a
consequence without naming the condition that would produce it.

**Decision blocked**: Nothing at lock. It affects whether A5 needs re-checking
if RDR 0007's implementation deviates — worth a line in the Prerequisites
sequencing note rather than a change to A5.

---

## Not filed

- The `contains` operator's set-element encoding gap (RDR 0007 Phase 4 routes
  the request to this RDR; Phase 3 and Testing Strategy require a `contains`
  test before the vocabulary is accepted) is real and user-affecting, but it is
  outside this delta-scope and is not part of A5/A7/A8 or the narrowing. Not
  filed as OUT-OF-SCOPE either, because the RDR already carries it explicitly in
  Capability Dependencies ("Requested alongside the set-valued element encoding
  RDR 0007 Phase 4 routes here") and in Phase 3 — it is charted, not missed.
- A8's substance checks out against source: JDR 0001 §JD-4
  (`./docs/jdr/0001-resolve-kernel-seam.md:238-241`)
  is open only as to the recording document, and RDR 0006's exhaustiveness
  clause (lines 344–350) covers only the non-finite case, exactly as A8 states.
  No finding on A8's accuracy — see P1-1 for its scope.
