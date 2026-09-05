Model: claude-opus-5[1m]

# Persona 1 — Product Manager

Owned starting set read: `0019:§problem-statement`, `0019:§approach`,
`0019:§decision-rationale`, `0019:MVV`.

**Widening declared.** The PM question — "does the outcome actually reach a
user?" — could not be answered inside the owned set, because the owned set
asserts the outcome ("the first-run wall is gone") without naming who hits
the wall or how many of them the verb reaches. Three silences sent me out:

1. `§approach` and `§decision-rationale` both claim the wall is removed, but
   neither states the population it is removed for. That sent me to
   `0019:C1` (carrier scope) and `§consequences`, where a population carve-out
   turns out to exist.
2. `§problem-statement` names an operator but no priority or blast radius.
   That sent me to `§background`, which carries the Priority: Low
   justification the problem statement omits.
3. `0019:MVV` validates against a *fixture* model only. That sent me to
   `models/` in the source repo to ask whether any shipped model is in the
   affected population.

Findings are severity-ranked.

---

## HIGH — `0019:§background` — the Priority: Low justification is
## contradicted by a shipped model in this repo

**Claim.** `§background` states: "Priority is Low on the blocking axis: the
shipped `models/rdr.toml` is a decision table with zero owned tags, so no
current first run is blocked." Every clause of that sentence is false of the
file it names.

**Evidence** (`/Users/cwensel/sandbox/newcoinc/intrastate/models/rdr.toml`):

- It declares **three owned tags** — `[tags.stage]`, `[tags.status]`,
  `[tags.gate_passed]`, each `provenance = "owned"` (lines 28, 47, 57) — not
  zero.
- It declares an **`[initial]` block** (line 91) seeding all three
  (`stage = "seeded"`, `status = "draft"`, `gate_passed = "false"`).
- It declares **no `class` key**. Per
  `/Users/cwensel/sandbox/newcoinc/intrastate/internal/table/model.go:501`
  ("empty string — reads everywhere as ClassStateMachine") and
  `::IsDecisionTable` (line 68, `m.Class == ClassDecisionTable`), it is a
  **state-machine** model. It is therefore in the class the verb serves, not
  the class `0019:C1` refuses.
- Its write accessor `[write.rdr-status]` (line 79) has neither `edit` nor
  `command`, so it is **file-backed** — inside `0019:C1`'s carrier scope, not
  in the disclosed carrier gap.

So the single shipped, non-example model in this repo is precisely a model
whose first run this RDR unblocks. `§background`'s own framing ("this is the
acceptance gate's SUBJECT", per the file's header comment) makes it the most
load-bearing model in the repo, not an incidental one.

**Why this is a PM finding and not a nit.** Priority is the field that says
how much user pain exists today. It is set to `Low` on a factual premise that
does not hold, and the premise is the *only* evidence offered for it. The RDR
simultaneously argues in `§decision-rationale` that option A (do nothing)
"leaves the user outcome unmet" — an argument whose force depends entirely on
somebody actually hitting the wall. `§background` says nobody does;
`models/rdr.toml` says the repo's own gate model does.

**Decision it blocks.** The Priority field, and with it the scheduling
decision (`Low` vs. the priority a live blocked first-run warrants) — and,
downstream, whether Phase 2's doc activation is a nice-to-have or a
prerequisite. It also weakens `§decision-rationale`'s rejection of
Alternative 1: A was rejected partly for leaving a wall nobody is
currently standing at, which is a materially different trade if the wall is
live in-repo.

---

## MEDIUM — `0019:§problem-statement` + `0019:§approach` — the stated user
## outcome is not traceable to the population the approach reaches

**Claim.** `§problem-statement` opens with an unqualified user outcome: "An
operator standing up a state-machine flow hits a first-run wall … The first
run of every state-machine flow needs a manual `set-state`." `§approach`
answers with an equally unqualified "the first-run wall is gone" framing
(echoed in `0019:MVV` step 4). Neither passage carries the carve-out that
`0019:C1` and `§consequences` establish.

**Evidence.** `0019:C1` (CARRIER SCOPE of the predicate) restricts the verb to
the **file-backed** write carrier and requires refusal against any model whose
bound write accessor is edit-carried (0028) or command-backed (0025). The
final `§consequences` bullet states this plainly and calls it "a live scope
gap, not a hypothetical one: those flows keep the manual first-run wall this
RDR removes elsewhere." The registry branch it rests on is real —
`internal/cli/flowbind/registry.go` selects among three writers.

So the honest outcome is "the first-run wall is removed for the file-backed
subset of state-machine flows." The problem statement promises it for
*every* state-machine flow, and the approach does not narrow it. The
qualification exists only two sections downstream and in a Negative bullet.

**Why it matters to a PM.** The problem statement is the passage that
survives into release notes, the kata, and the `--help-all` text Phase 2
writes. As written it will be read as a complete fix. A user on an
edit-carried model reads "first run of every state-machine flow" and then
gets a refusal.

**Decision it blocks.** The scope boundary of the shipped promise — i.e.
what Phase 2's `docs/cli-reference.md` and `--help-all` text may claim, and
whether the successor (gated on `0019:A6`) is a follow-up or part of this
outcome. Also blocks a clean answer to "is this RDR done?": on the stated
problem it is not, on the scoped problem it is.

---

## MEDIUM — `0019:MVV` — the validation never runs against a model that has
## the problem

**Claim.** Every step of `0019:MVV` runs against "a fixture state-machine
model declaring `[initial]` with one always-present owned key and one plain
owned key, both writer-served." No step runs against a model that ships.

**Evidence.** Steps 1–9 are all fixture-scoped; the only non-fixture arm
(step 8) is a decision-table refusal, also fixture-shaped. Meanwhile
`models/rdr.toml` is a shipped, file-backed, `[initial]`-declaring
state-machine model (see the High finding) whose owned keys include a `bool`
(`gate_passed`) and two `enum`s — a value-kind mix the fixture ("one
always-present owned key and one plain owned key") does not pin.

**Why it matters to a PM.** MVV is the passage that answers "will the user
actually get the outcome?" A fixture-only MVV proves the mechanism, not the
outcome. `§decision-rationale` claims the shipped doc
`docs/model-authoring.md` "already promise[s] start-state semantics the
runtime does not deliver" — I confirmed that promise at
`docs/model-authoring.md:572` ("`[initial]` declares the owned state a model
starts from"). The MVV never demonstrates that the promise is now kept for
any model a user would actually encounter.

**Decision it blocks.** Whether the MVV as written is sufficient for the
Scope Verification gate (`0019:G-scope`, currently unanswered). A one-line
addition — run `init-state` against `models/rdr.toml` bound to a fresh
artifact and assert the three `[initial]` keys land — would make the outcome
claim demonstrable rather than inferred.

---

## LOW — `0019:§consequences` — the two-command first-run cost is recorded
## but never sized against the outcome

**Claim.** The bullet "first-run UX is two commands (init, then work) where
option B/C would have been zero/one — the explicitness is bought with a step"
records the cost honestly but never says whether the user is expected to type
`init-state` by hand, or whether some wrapper/skill invokes it.

**Evidence.** `§approach` describes the verb as "an explicit lifecycle step
at session start." `0019:BR4` rejects "Seed at artifact creation by an
external tool/wrapper" and `0019:BR5` rejects "Auto-init on first
`set-state`" — so both automatic-invocation routes are closed, which means
the answer is "by hand," but no passage states it. Phase 2
(`0019:§activation-step-1-contract-and-reference-docs`) lists the doc
surfaces that must change but does not name a discovery path: nothing tells
an operator hitting `unknown[].reason: absent` in `flow next` that
`init-state` is the cure.

**Why it matters to a PM.** The problem statement's user is an operator who
did not know they needed `set-state`. Replacing an undiscoverable manual
`set-state` with an undiscoverable manual `init-state` moves the wall rather
than removing it, unless something points at the new verb. `§failure-modes`
F3 does name "`flow next`'s own unknown report" as *a* diagnosis surface,
but for the different case of a non-`[initial]` owned key.

**Decision it blocks.** Whether Phase 2 owes a discovery affordance (e.g. the
absent-key path naming the verb) or whether reference-doc text is the whole
activation. Low because it is an activation-scope question, not a contract
one.

---

## Not findings (checked, clean)

- The problem statement names a **user-visible** outcome (a first run that
  needs a manual `set-state`), not an engine-internal irritation. Grounded:
  `internal/graphlint/analysis.go:128` emits "the model declares no initial
  owned state", and no non-test `internal/cli` code reads `Model.Initial`
  (only `mvv_0010_test.go` and `lint_gate_0006_test.go`) — the stated
  lint-demands-it / runtime-ignores-it gap is real.
- `§approach` traces to what the design does: the
  `internal/cli/flow_state.go::groupByWriter` (line 645) → executor →
  read-back path it names exists and is the `set-state` path.
- Scope did **not** expand past the problem: the `[initial]`-key-set
  boundary in `0019:C1` ("the verb claims nothing about owned keys
  `[initial]` does not assign") holds the line, and `0019:BR1`–`BR5` record
  five adjacent temptations as rejected rather than absorbed.
- `0019:G-scope` and `0019:G-proportionality` are unanswered template
  placeholders. Not reported as findings — this is a pre-lock pass and the
  Finalization Gate is Stage 7's to answer.
