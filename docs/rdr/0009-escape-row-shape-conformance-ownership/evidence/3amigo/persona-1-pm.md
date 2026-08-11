Model: claude-opus-5[1m]
Lens: 3amigo (persona 1 — product manager)

# Persona 1 — Product Manager: does RDR 0009 deliver the stated user outcome?

**Stated user outcome (Problem Statement):** a table author writes an escape
rule expecting a pure exit route, and is harmed when "something downstream
persists a tag they never authorized." The harm is *unauthorized persistence
of owned state*, discovered late, with nothing in the disposition disclosing
it.

Judged against that outcome, the RDR is internally rigorous but the population
it protects and the population it names are not the same one, and the one
protection the *table author* actually receives is the half the RDR does not
build.

---

## PM-1 — The RDR's only shipped enforcement is on the path the named user
cannot reach; the author-facing half is deferred to a peer that does not exist

**Severity:** high

**Passage:** *Capability Dependencies* table —

> | Load-time rejection of malformed escape declarations | Predecessor
> (RDR 0002, Final, unimplemented) | **Deferred** | Conformance fixtures from
> this RDR bind the normalizer build |

and *Approach*:

> On the authored path, RDR 0002's normalizer rejects the rule at load — its
> existing "malformed escape declaration" validation class, diagnosable to the
> source rule.

**The gap:** The Problem Statement's user is a *table author* — someone who
"writes an escape rule" in an authored table. The RDR's own Context concedes
that path is not reachable at all today: "the RDR 0002 normalizer is not yet
implemented, so the combination is reachable only from hand-constructed
`resolve.Row` values — tests and any future non-TOML table producer (A5)."
Verified at HEAD: `internal/resolve` has no production importer
(`internal/cli/root.go::NewRootCmd` registers only `newVersionCmd()`;
`cmd/intrastate/main.go` is `func main() { cli.Execute() }`), and `Row` is
constructed only in `internal/resolve/*_test.go`. So the two things this RDR
actually *ships* — Phase 1 (kernel entry precondition) and Phase 2 (fixture
conformance) — protect Go programmers and this repo's test suite. The table
author gets Phase 3, which is a *fixture set binding a future build*, not an
enforcement.

That is not a scoping flaw by itself; the flaw is that the document never says
so plainly in user terms. It says "no write-bearing escape plan can be emitted
by any producer path" (Consequences) as if the outcome were delivered, while
the enforcement for the *named* user is `Deferred` in the RDR's own table. A
reader takes away "the author is protected now." The truthful statement is
"the author is protected when RDR 0002 is implemented; until then the
guarantee is a fixture set nobody has run."

**Decision it blocks:** Whether RDR 0009 can be considered *done* on merge of
Phases 1–2, or whether "done" for the stated user outcome is gated on RDR
0002's implement stage. The Prerequisites checklist has the RDR 0002 binding
as an unchecked box (`[ ] RDR 0002's implement stage binds to this RDR's
conformance fixtures (Phase 3) once both are Final`) — so an implementer
cannot tell whether shipping Phases 1–2 closes this RDR or leaves it open, and
an operator cannot tell what protection is live today.

---

## PM-2 — The user-visible symptom in the Problem Statement is never
traced to a user-visible remedy; the whole solution stops at the kernel
boundary and hands the actual harm to a peer contract

**Severity:** high

**Passage:** *Problem Statement* —

> They discover the problem only if something downstream persists a tag they
> never authorized ... afterwards an owned tag has a value no rule in their
> table ever set.

vs. **A4 → Carried forward (peer-contract dependency)** —

> `resolve.go::planOf` copies `NextTags` and `Writes` symmetrically with no
> escape-awareness, so this safety rests **entirely on RDR 0004's
> write-accessor scoping, not on any kernel-local property**. A future RDR
> 0004 implementer widening "apply the plan" beyond planned writes would
> breach this RDR's invariant.

**The gap:** The stated harm is *persistence*. Persistence happens in RDR
0004's accessor layer, not in the kernel. This RDR's entire mechanism is
upstream of persistence: it guarantees `Plan.Writes` is empty on escape plans.
Whether an empty `Writes` actually means "nothing persisted" is, by A4's own
admission, a property of RDR 0004 — which is `Final` but unimplemented (per
the RDR's own Technical Environment listing). Verified at HEAD:
`resolve.go::planOf` at line ~506 does copy both fields symmetrically
(`NextTags: copyTags(row.NextTags), Writes: copyTags(row.Writes)`), and the
`Plan` type doc distinguishes them only in prose ("NextTags is the next
state." / "Writes are the owned-tag writes for the accessor layer.") — there
is no structural barrier.

So the RDR delivers "the escape row's `Writes` field is empty," and *asserts*
that this equals "no unauthorized persistence" via a chain of three normative
citations to an unimplemented peer. That is a defensible design position, but
the RDR states the user outcome as achieved ("the ADV-2b 'kernel guessing at
persistence' class is closed before RDR 0004 can apply an unsanctioned write")
without ever writing down that the closure is conditional on an RDR 0004
implementation that does not exist and is not bound by this RDR's
Prerequisites. Compare: RDR 0002's obligation *is* bound in Prerequisites;
RDR 0004's is not, despite carrying the actual user harm.

**Decision it blocks:** Whether a Prerequisite or a normative clause must bind
RDR 0004's implement stage the way Phase 3 binds RDR 0002's. Right now
nothing stops an RDR 0004 implementer from widening "apply the plan" and
silently reintroducing the exact user symptom this RDR was opened to prevent
— and A4 says so in as many words while filing it as "not a predicate gap."

---

## PM-3 — Success is defined in kernel-internal terms; there is no
statement of what the user observes when the fix works

**Severity:** medium

**Passage:** *Validation → Testing Strategy* opening —

> Done means: every producer path is closed, no conforming table changes
> disposition, and the check itself is proven non-vacuous.

**The gap:** All three success clauses are properties of the codebase, none is
a property of the user's experience. "Every producer path is closed" is
falsified in the same document (PM-1: the authored path's enforcer is
`Deferred`). "No conforming table changes disposition" is a
no-regression statement. "The check is non-vacuous" is a test-quality
statement. Nowhere does the RDR state the outcome in the Problem Statement's
own vocabulary — something like *"a table author can never end a flow with an
owned tag their table did not set"* — nor does it state how anyone would
observe that this is now true.

This matters because the twelve validation scenarios are all kernel/CLI
mechanics (error types, `errors.Is`, sort order, exit codes) plus one scenario
(8) that binds a peer's unwritten normalizer. Scenario 8 is the only one
touching the authored path, and its Expected is a *load-time rejection*, i.e.
a diagnostic — not a demonstration that the user's owned state is safe. The
document never closes the loop from "escape row has no `Writes`" back to "the
author's owned tags are unmolested," which is the sentence the Problem
Statement promised.

**Decision it blocks:** An author or reviewer cannot judge whether the RDR's
scope is right-sized to the problem it opened, because there is no user-terms
success criterion to measure the scope against. Concretely: the Finalization
Gate's Scope Verification section is unanswered and cannot be answered
meaningfully without one.

---

## PM-4 — Two different users are conflated under one "diagnosability"
claim, and the one who is actually served is not the one in the Problem
Statement

**Severity:** medium

**Passage:** *Trade-offs → Consequences* —

> Positive: authored-path failures stay author-diagnosable (load-time, named
> source rule); programmatic-path failures name the offending row.

and *Failure Modes → Diagnosis (revised at Resolve)* —

> The runbook entry is therefore "branch on the code; the named rows identify
> the broken producer," not "read the error text."

**The gap:** These describe three distinct audiences — the table author, the
Go programmer who built a bad producer, and the operator reading CLI output —
and the RDR treats "diagnosable" as one property satisfied for all of them.
It is not. The named user (table author) gets diagnosis only through the
`Deferred` RDR 0002 half. The Go programmer gets the typed error (real,
shipped in Phase 1). The operator gets exit 2 plus a `Code`, `Detail`, and
`Hint` — but only if a `flow` verb wraps it, and no `flow` verb exists at
HEAD (verified: `internal/cli/root.go`).

Worse, the operator-facing story is actively unsettled: A2's "Second
obligation (found at Pre-Lock)" concedes that `clierr.CLIError.Cause` is
`json:"-"` (verified at `internal/cli/clierr/clierr.go:66`, and `Group` is
`json:"-"` too at line 61), so the `RowRef` identities do not reach the wire
without an *additive envelope change the RDR leaves as a choice*: "either the
existing `Detail` ... or a new `omitempty` field." The RDR never decides
which. For a document whose Problem Statement is about a user not being told
what happened, leaving the wire carrier undecided is squarely on-topic, not a
detail.

**Decision it blocks:** The CLI/verb implementer cannot decide the serialized
carrier for row identities (reuse `Detail` vs. add a field), and therefore
cannot write scenario 11's assertion ("an assertion that reads identities off
the JSON output is what proves the wire carrier exists") — the scenario names
a proof for a carrier the RDR declines to specify.

---

## PM-5 — The deliberate behavior change (dormant-row breach) is
justified only by "no production constructor exists," with no statement of
what a table author would experience if one did

**Severity:** low

**Passage:** *Trade-offs → Consequences*, final bullet —

> Negative (behavior change, deliberate): a hand-built table carrying a
> *dormant* malformed row — one no resolution path reaches — previously
> resolved fine and now errors on every `Resolve` (whole-table precondition).
> ... Whether any deployed table can regress turns on A5 (no production
> `resolve.Row` constructor at HEAD); if A5 fails, a migration step becomes a
> prerequisite.

**The gap:** The whole-table scope is the RDR's most user-affecting choice: it
converts a working table into a hard error because of a row that never fires.
The justification is entirely "nobody is running this yet" (A5 — verified true
at HEAD). But the whole point of Phase 3 is that authored tables *will* exist,
and an authored table's dormant escape rule with a stray write block will,
once RDR 0002 ships, fail the *author's whole table at load*. The RDR argues
the whole-table scope on a correctness ground ("a malformed table is malformed
as a value") and never weighs it as a user-experience trade: is failing an
entire flow table for an unreachable row the behavior an author wants, versus
a lint warning? RDR 0006 (graph lint) is cited as a neighbour but only for
boundary-drawing, never for "could this be a warning instead of a hard
failure."

The RDR does pre-empt the *escalation* question (Risks: "RDR 0002's
implementer treats the kernel precondition as *the* enforcer") but not the
*severity* question (should the authored-path breach be blocking?). The
authored-path clause simply says "MUST be rejected at table load," and Phase 3
extends that to `writes = []` — strictly *more* aggressive than the kernel —
without a user-outcome argument for why an empty block is worth failing an
author's table over. The normative text calls it "a statement of intent worth
refusing," which is an assertion of taste, not a user-outcome argument.

**Decision it blocks:** The RDR 0002 implementer cannot decide whether a
malformed escape declaration is a blocking load failure or a lint-level
diagnostic for the whole class, and cannot calibrate how aggressive the
empty-write-block rule should be against author friction.

---

## Not raised (out of persona)

- Error type/sentinel spelling, `errors.Join` vs. bare error, `CheckValid` vs.
  `Validate` naming, `compareRefs` collapse mechanics — Implementer.
- Whether scenarios 9/10/10b are sufficient, mutation-check adequacy,
  non-vacuity method — QA.

## Verification notes

Claims checked against HEAD (`internal/resolve/resolve.go`,
`internal/cli/clierr/clierr.go`, `internal/cli/root.go`,
`cmd/intrastate/main.go`):

- `Resolve(in Input) (Result, error)` — confirmed, line 318.
- Five refusal kinds, closed set — confirmed, lines 49–71.
- `planOf` copies `NextTags` and `Writes` symmetrically with an `escaped`
  flag — confirmed, lines ~506–509.
- `CLIError.Cause` and `.Group` are both `json:"-"` — confirmed, lines 61, 66;
  `Detail` and `Hint` are serialized (`json:"detail,omitempty"`,
  `json:"hint,omitempty"`), lines 56, 58.
- No `flow` verb; `internal/resolve` has no production importer — confirmed.
- **Unverified by me:** the contents of RDR 0002, 0004, 0005, 0006 normative
  text (quoted only as this RDR reports them); the A3 spike outputs.
