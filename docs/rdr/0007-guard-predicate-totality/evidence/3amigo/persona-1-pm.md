Model: claude-opus-5[1m]

# 3amigo Persona 1 — Product Manager (RDR 0007)

Lens: does this RDR deliver the user outcome it names? The stated outcome is
in the Problem Statement: "They want a refusal they can act on: told plainly
that the artifact state needed to decide was missing, not handed a plan that
quietly routed around it."

Findings are severity-ranked. Severity reflects distance between the promised
user outcome and what the RDR actually commits to ship.

---

## PM-1 — The promised refusal never names the missing state; only the "not masked" half of the outcome ships

**Severity: High**

**Passage** — Problem Statement: "They want a refusal they can act on: told
plainly that the artifact state needed to decide was missing" — read against
Failure Modes, *Diagnostic gap*: "`GuardResult` carries no missing-tag
payload, so 'which tag was absent' is inferable only from the guard text; if
that proves insufficient in practice, the enrichment belongs to RDR 0003's
semantic-kind diagnostics (Phase 3 hands it off)".

**Problem.** The Problem Statement makes a two-part promise: (a) don't hand me
a plan that routed around missing state, and (b) tell me plainly *what state
was missing*. The chosen approach delivers (a) by construction — `guard_
unevaluable` is non-escapable per RDR 0002. It does not deliver (b), and the
RDR downgrades (b) to a conditional hand-off ("if that proves insufficient in
practice"). Grounding the shipped surface: `internal/resolve/resolve.go`'s
`Refusal` carries `MissingOwned []string` on `owned_state_unavailable` but
only `Guard string` on `guard_unevaluable` — a raw predicate string, not the
absent key. So the user who today gets an escaped plan will tomorrow get a
refusal reading roughly `guard_unevaluable, guard="..."`, and must diff the
guard's referenced tags against the view by hand to learn which one was
missing.

That is a materially different user outcome than the one the Problem
Statement sells, and the Consequences section reasserts the strong version
anyway: "the user always learns 'I could not tell' as a distinct,
non-escapable refusal" — true, but "I could not tell" without "about what" is
exactly the actionability the Problem Statement said was the point. The
asymmetry with the sibling refusal (`owned_state_unavailable` *does* name its
keys) is the tell: the kernel already knows how to name missing keys for the
other absence refusal, and this RDR declines to for its own.

**Decision blocked:** whether naming the absent tag is in scope for *this*
RDR's shipped outcome or is a genuinely separate RDR 0003 concern — and
therefore whether RDR 0007 can be called done when the user gets a refusal
they still cannot act on without manual inspection. Also blocks: should the
`Refusal` struct grow a missing-tag payload here (the RDR forecloses this as
"a kernel enum change" without pricing it against the outcome it promised).

---

## PM-2 — Zero migration story for the tables and authors the change breaks

**Severity: High**

**Passage** — Trade-offs / Consequences: "Negative: tables whose authors
expected closed-world reads ('x == 1 just fails when x is missing') now refuse
where they previously escaped or pruned; those guards must be rewritten with
explicit existence atoms."

**Problem.** This is a one-line acknowledgment of a behavior change that flips
working tables into hard, *non-escapable* refusals, and the RDR carries no
answer to any of the questions an affected user asks. Not stated anywhere in
the document: how many such tables exist today (the Background probes
established the masking path exists but never counted real tables that ride
it); whether the break is detectable before it fires in production (a lint
that flags value operators over tags that may be absent is explicitly rejected
under Briefly Rejected — "Lint-only totality" — as a *design* alternative,
with no consideration of it as a *migration* aid, which is a different job);
or whether affected authors get any warning phase. The Implementation Plan's
three phases are contract doc alignment, conformance vectors, and a
diagnostics/authoring handoff — no migration phase, no inventory step, no
"find affected tables" tooling.

Compounding it: the break is *silent until triggered*. A table with a
closed-world guard keeps working on every input where the tag happens to be
present; it fails only on the input where it previously escaped. So the user
discovers the break at the worst moment, in production, on the exact case they
had modeled an escape for.

**Decision blocked:** whether this RDR ships as a breaking semantics change
with no migration path, or whether an inventory/lint/warn step is a
prerequisite. Also blocks scheduling: without a count of affected tables,
nobody can decide if this is a quiet doc change or a coordinated migration.

---

## PM-3 — The "degraded artifact has no hatch" consequence is stated as a footnote, not weighed as a product decision

**Severity: Medium-High**

**Passage** — Trade-offs / Consequences: "Negative: `guard_unevaluable` is
non-escapable by design, so an artifact in genuinely degraded state has no
modeled-escape hatch through a guard that cannot read it; recovery is fixing
the artifact state or the accessor, not the table."

**Problem.** The whole value proposition of RDR 0002's modeled escape is that
a table author can say, in reviewable data, "here is what to do when the
resolver cannot pick an edge." This RDR removes that authority for an entire
class of situation — and the class it removes it for is precisely the
degraded-state case where an author most plausibly *wants* a modeled fallback
("if I can't read the gate, route to manual review"). The RDR treats this as
a cost line item rather than a design question, and never asks whether the
legitimate use of that hatch exists.

The Decision Rationale scoring matrix reinforces the gap: its criteria are
correctness fit, prior-art alignment, reversibility, blast radius, and cost.
There is no row for *author expressiveness* or *operator recovery ergonomics*
— the two dimensions on which B is strictly worse than A. So the matrix that
"decides" the approach cannot see the cost this bullet names. The
counter-argument the RDR could make (a modeled escape from unevaluable is
indistinguishable from the masking the RDR exists to prevent — which is
probably the right answer) is never actually written down; the reader is left
to construct it.

Related unpriced consequence: A6b is carried OPEN, and Failure Modes /
*Conflated recovery signal* concludes "operators should treat the refusal as
'inspect the artifact and accessor,' not 'retry'". So the operator-facing
outcome is: a hard, non-escapable, non-retryable refusal whose cause could be
a transient accessor failure. That combination — no table hatch, no retry
guidance, no missing-key name (PM-1) — is a strictly worse operator
experience than today's escape, and no section adds those three up.

**Decision blocked:** whether a distinct escapable-on-purpose class (or an
operator-level retry/override affordance outside the table) is needed before
this ships, versus accepting a hard stop. Also blocks whether A6b's open state
is really "not a blocker for this RDR's lock" (Prerequisites) once the
operator-facing consequence is priced.

---

## PM-4 — "unevaluable blocks a decided-true sibling" is a real user-visible behavior change presented only as internal aggregation

**Severity: Medium**

**Passage** — Normative Contracts (aggregation block): "if any surviving
candidate row's guard is GuardUnevaluable, the resolution MUST refuse
guard_unevaluable — a decided-GuardTrue sibling row MUST NOT be selected while
an unevaluable candidate exists."

**Problem.** From the table author's seat this is a large and surprising
property: an *unrelated* row — one that does not match, that the author may
consider dead weight, that merely happens to be a candidate and reads a tag
the view lacks — vetoes a row that unambiguously applies. That is a
coupling-between-rows story with direct authoring consequences (a table gets
more fragile as it grows; adding a row can break resolution of an existing
one), and it is exactly the kind of property that belongs in the Problem
Statement's user narrative and in the authoring guidance.

Instead it appears only as a normative clause and as MVV scenario 3
(*unevaluable-blocks-true-sibling*), justified purely by what the kernel
already does (A3: "`gate`'s `len(undecidable) > 0` return overwrites the named
return `selected` with `nil`"). Note the direction of that argument — the
shipped implementation is being ratified into contract because it exists, not
because the user outcome was reasoned to it. The RDR itself flags that no
shipped test covers it ("**No shipped test covers this**"), i.e. it was
incidental behavior, now being made a promise. Whether that is the *right*
promise for a table author is never asked. Note the contrast with D5, which
the RDR cites approvingly for the opposite instinct — "so that an unrelated
row's key cannot poison a legal resolution" (Context/Background). D5 exists to
stop unrelated rows poisoning resolution; this clause lets them, one verdict
over, and the RDR does not reconcile the two from the author's point of view.

Grounding the diagnostic consequence: `gate` reports `Guard: lowest.Guard`
from `slices.MinFunc` over the undecidable set — the user is shown one
arbitrary-by-row-order guard among possibly several, compounding PM-1.

**Decision blocked:** whether veto-by-unrelated-candidate is the intended
authoring model (and therefore belongs in Phase 3's authoring guidance as a
first-class rule authors must design around), or an implementation artifact
that should be revisited before it is frozen as normative.

---

## PM-5 — The scope is spec-plus-fixtures, but the outcome only lands when RDR 0003 is built — and nothing in scope forces that

**Severity: Medium**

**Passage** — Risks and Mitigations, first risk: "Mitigation: partial, and the
residue is real. … `go test ./internal/resolve/` stays green against the stub
no matter what 0003 builds — the same blind spot as premortem P-10. … RDR
0003's implement stage MUST instantiate it against its own evaluator. …
**until RDR 0003 accepts it, drift is caught by review only.**"

**Problem.** Stated as a delivery question rather than an engineering one:
this RDR ships a rule, a doc-comment narrowing (Phase 1), and a conformance
harness (Phase 2). None of those change what any user experiences. The masking
path recorded under Background — "A FALSE guard plus an absent
`RequiresOwned` key plus a modeled `no_match` escape yields a plan via the
escape, `Escaped:true`. This is the path where missing artifact state is
masked" — remains open in the shipped product on the day this RDR is marked
done, and stays open until RDR 0003's evaluator is built and voluntarily opts
into the harness. The RDR is candid that the binding is not enforceable from
here ("until RDR 0003 accepts it, drift is caught by review only") and that
the obligation is only "handed to RDR 0003 in Phase 3."

So the honest scope line is: this RDR does not close the masking path; it
writes down how someone else will close it later, and asks them to run a test
suite. Yet Consequences claims the outcome in the present tense — "Positive:
missing artifact state is never maskable behind an escape" — and the
Approach's closing clause says the rule means missing state "can never be
masked behind an escapable refusal class, which is the user outcome this RDR
exists to secure." Both read as shipped guarantees. Nothing in the document
states when a user actually stops being able to get a masked plan, or who owns
that date. Prerequisites bind sequencing in only one direction ("This RDR
Final **before** RDR 0003's implementation begins") — which prevents 0003 from
shipping wrong, but places no obligation on 0003 shipping at all.

**Decision blocked:** whether this RDR can be accepted as "delivering the user
outcome" or must be tracked as an unshipped promise until RDR 0003 lands —
and, downstream, whether an interim mitigation (a kernel-side guard on the
recorded masking probe, or a lint on tables using the vulnerable pattern) is
warranted so users are not exposed for the whole gap between the two RDRs.
