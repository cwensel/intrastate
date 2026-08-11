Model: claude-opus-5[1m]
Persona: Product Manager

# 3amigo — Persona 1: Product Manager (RDR 0008)

Lens: does this RDR deliver the outcome its own Problem Statement promises — an
author who names a recognized-provenance tag in their own vocabulary, whose row
silently never fires, and whom nothing tells why? Findings are ordered by
severity, then by how directly they hit that user.

---

## PM-1 — The chosen branch removes the capability the Problem Statement says the user wanted, and never says so

**Passage**: Problem Statement — "They name their tags themselves — that is the
point of the declared model — so they name the recognized-provenance one
whatever reads best in their flow." Against Proposed Solution / Approach item 2 —
"the declaration with provenance `recognized` MUST be named `recognized`. The
author declares *that* they use the affordance, not what it is called." And
Alternatives Considered / Alternative 2 Cons — "Deletes the exact capability the
Problem Statement's user wanted (matching the recognized outcome in tag
vocabulary)."

**Concern**: The Problem Statement frames the user's want as *using their own tag
vocabulary for this tag*. The chosen approach denies exactly that want — it fixes
the name — and the RDR nowhere acknowledges that it is refusing the user's stated
desire rather than satisfying it. Worse, the RDR uses the "deletes the capability
the user wanted" argument as a *rejection reason against Alternative 2* while the
chosen branch removes the same naming freedom; the phrase "the capability the
user wanted" silently shifts meaning between the two places. In the Alternative 2
Cons it means "matching the recognized outcome in tag vocabulary at all," but in
the Problem Statement the want was "using *their own* name for it." The
Trade-offs / Consequences bullet ("Negative: authors lose naming freedom for
exactly one tag") is the only place this lands, and it is filed as an accepted
cost rather than as a redefinition of the goal. The RDR should say plainly: the
user's real outcome is *"my row fires, or I am told why"* — not *"I get to name
this tag."* That is a defensible and probably correct reframing, but it is
performed silently, which is what makes it a PM problem rather than a design one.

**Decision it blocks**: An author (and the 0002 implementer writing the lint
message) cannot tell whether the reserved name is a deliberate product stance
they should teach — "this one tag is engine-owned, like `event` in XState" — or
an incidental implementation constraint they should apologize for and eventually
lift. That determines the whole tone and content of the `reserved tag key`
message, which per contract block 3 is the *only* thing the user ever sees.

**Severity**: High

---

## PM-2 — "Authoring-time diagnosis" is asserted but never traced to a moment the author is actually present for

**Passage**: Proposed Solution / Approach item 3 — "reported by RDR 0002's
normalizer before any resolution, which converts the seed's silent no-match
discovery into an authoring-time diagnosis." Also Trade-offs / Consequences —
"the naming collision surfaces at authoring time as a typed load/lint failure,
replacing the silent no-match discovery in the Problem Statement."

**Concern**: "Authoring time" is doing load-bearing work in the outcome claim and
is never defined. What the contracts actually specify is *load/lint time* — the
moment the normalizer runs. Nothing in this RDR establishes that the author runs
the normalizer while authoring. RDR 0005's CLI code-string table
(`flow-tag-invalid`, `flow-model-not-found`, `flow-zero-match`, …) contains no
entry for a data-level validation category from this RDR, and RDR 0006 owns the
lint command surface but is explicitly held at arm's length in Phase 2 ("Not RDR
0006's graph lint"). So the RDR reserves a category, requires it carry good data
(block 3), and then routes it through no named user-facing command. The failure
may well arrive at the same moment as the old silent no-match — at the first
resolve the author runs — with the only improvement being that the message is now
correct. That is still a real improvement, but it is a *much smaller* claim than
"converts silent discovery into authoring-time diagnosis," and the RDR banks the
larger claim.

**Decision it blocks**: Whether this RDR is done when the category exists in the
normalizer, or is only done when some command the author already runs surfaces
it. That is the difference between Phase 2 being the last phase and this RDR
owing a dependency on 0005 or 0006 for delivery. It also blocks the Capability
Dependencies table: there is a row for "Load/lint name validation" but none for
"a way the author sees the failure."

**Severity**: High

---

## PM-3 — The residual silent-no-match case is the *more likely* version of the Problem Statement's user, and it is filed as an accepted residual

**Passage**: Trade-offs / Failure Modes — "Silent, residual: an author who
declares an *observed* tag under an innocent name (`result`, `outcome`)
intending recognizer semantics passes every name check and the row never fires —
the original silent no-match survives this narrow path, because intent is
invisible to a name-and-provenance validator (premortem P-3)."

**Concern**: This is described as a "narrow path," but read against the Problem
Statement it is arguably the *main* path. The Problem Statement's user is someone
who does not know a recognized-provenance tag key is a distinct, engine-owned
thing — that is precisely why they named it themselves. Such a user is at least
as likely to reach for `[tags.outcome]` with `provenance = "observed"` (or to
omit the declaration and match on a name they invented) as to declare
`provenance = "recognized"` under a wrong name. The RDR's own A1 evidence
confirms the pull of this shape: *both* canonical 0002 fixtures name the
recognized-provenance tag `outcome`, not `recognized` — i.e. every existing
example of an author naming this thing chose the innocent name. The fix only
catches the user who already got the *provenance* right and only the *name*
wrong. The RDR never states what fraction of the Problem Statement's population
that is, nor does it revise the Problem Statement to scope itself to that subset.

**Decision it blocks**: Whether the advisory heuristic mentioned in the same
paragraph ("an advisory lint heuristic for that shape is a Resolve question") is
optional polish or is required for this RDR to close its own Problem Statement.
Until that is settled, an implementer cannot size Phase 2, and a reviewer cannot
judge the Minimum Viable Validation as sufficient — the MVV pair exercises only
the wrong-name case, never the innocent-name case.

**Severity**: High

---

## PM-4 — Success is stated in contract terms; no criterion is written from the author's seat

**Passage**: Validation / Testing Strategy — "Done means: every scenario below
has a green test, and the guard/matcher same-view property (scenario 3) is pinned
by a test that would fail if the two views diverged." And Implementation Plan /
Minimum Viable Validation — "(b) the same table with the declaration renamed …
fails load/lint with the `reserved tag key` category whose failure data names the
required name — not a silent no-match at resolve time."

**Concern**: Every "done" statement is expressed as a property of the system
(category emitted, data fields present, views identical). None is expressed as a
property of the *author's experience*: that a person who hit the original problem
can, from the message alone and without reading this RDR or RDR 0002, get their
row firing. The MVV's user-facing half stops at "failure data names the required
name" — a data-shape assertion, checked in scenario 7 as a golden-text check *at
the data level, not the renderer*. That is a deliberate and reasonable division
of labor, but it means no scenario anywhere asserts that a human-readable
remediation reaches a human. The RDR asserts in the Consequences that the loss of
naming freedom is "mitigated by the lint message naming the rule," yet the lint
message's user-facing text is explicitly outside every test.

**Decision it blocks**: The 0002 implementer cannot tell whether writing the
three data fields discharges this RDR's obligation, or whether they also owe
message wording that a first-time author can act on. Since block 3 says "the
guidance travels in the failure data, not the renderer," they will reasonably
choose the former — and the Problem Statement's "nothing in the refusal tells
them" is then only half fixed.

**Severity**: Medium

---

## PM-5 — The Problem Statement's second paragraph swaps the user's problem for an internal ownership problem, and the rest of the RDR follows the swap

**Passage**: Problem Statement — "Internally, this RDR must fix **who owns the
name** of the recognized-provenance tag key in the assembled evaluation view."
And the Metadata / Profile field — "one contract: who owns the *name* of the
recognized-provenance tag key in the assembled view." And Decision Rationale —
"QOC matrix (question: *who owns the name of the recognized-provenance tag key in
the assembled view?*)."

**Concern**: The user's problem is "my row silently never fires and nothing tells
me why." The ownership question is a *means* to that end, not the end. But from
the second paragraph onward, ownership becomes the stated question everywhere
that matters — the Profile, the QOC matrix question, the contract count, the
deciding rows. The QOC matrix's criteria are correctness fit, prior-art
alignment, blast radius, reversibility, peer consistency, and cost; only the
first mentions the user's outcome, and it is not one of the two deciding rows
("Deciding rows: **blast radius** and **prior-art alignment**"). So the branch is
chosen on internal-cost and precedent grounds, with the user outcome relegated to
a non-deciding row. This may still yield the right answer — R does look best for
the user too — but the RDR does not demonstrate that, and a reader cannot check
whether a branch that served the user better but cost more would have won.

**Decision it blocks**: Whether the open Enforcement-locus decision (candidates
a/b/c, defaulting to (b)) should be settled on user-outcome grounds or on
surface-cost grounds. Candidate (a) is explicitly "zero new surface, no
detection" — which is the *cheapest* on the deciding criteria and the *worst* for
the user, and the RDR's own A6 Residual concedes that under (a) the silent
shadowing "detected by nothing" survives. With the user outcome absent from the
deciding rows, nothing in the RDR's decision framework rules (a) out; only an
informally stated "lean" does.

**Severity**: Medium
