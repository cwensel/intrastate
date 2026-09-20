Model: claude-sonnet-5

# Critique: RDR 0030 (Computed write-value grammar — `step`)

## Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | 0030:A7 | No load-time ceiling on row expansion; the record explicitly declines one, defending on a spike that stops at 500,000 rows / 5s / 584MB–6GB extrapolated. A modest cross-product of two stepped tags (`attempt` × `tier`, each declared wide) blows past "tens of rows" silently. | `intrastate lint` hangs (quartic, 600s+ at 1,000 rows per the RDR's own Performance Expectations) or the process is OOM-killed, with no refusal, no warning, no way to know at `load` time that the model they just authored is unanalyzable. | premortem, §2 |
| C-2 | 0030:A8 (Refuted), C3 | The identity split (dump/graph show `rule#cell`, `flow next`/`flow resolve` show bare `rule`) is asymmetric by design, discovered as a *refutation* of the original assumption mid-design. Authors and downstream tools that join `flow resolve`'s `rule` back to a `dump`/`graph` identity get one-of-N ambiguity, resolved only by "first match wins" (`0010:C3`), silently. | A user watching `flow next`/`flow resolve` output cannot tell WHICH of the N expanded cells fired; they see `retry`, not `retry#3`, on a machine that just consumed a counter. | §1, premortem |
| C-3 | 0030:§consequences ("The `enum` arm delivers the stepped VALUE once, not the cap once") | The RDR's own headline promise — "say it once" — is admitted, in its own Trade-offs section, to be false for the primary motivating case (tier escalation, an enum). The cap exclusion for an enum's terminal member requires `in = [<every member but the last>]`, a full domain re-listing that must be hand-edited on every domain change — the exact defect (re-typing on domain change) the Problem Statement exists to kill. | Author adds a fourth tier to the domain, ships it, and the model silently admits stepping past the old terminal member into the new one (or off the end, if they forgot to update the `in` list) — the same "recompute by hand or the model breaks" workflow reintroduced in the one place the record calls out as the residual. | §2, §3 |
| C-4 | 0030:A10, A11, A12, A13, A14 | Five of C1's normative sub-clauses (subsumption, both-carriers, the evaluator relocation, bound-check ordering, step+clear conflict) are `Pending`, i.e. unverified until the MVV runs at Phase 2 implementation time — not before lock. The record locks (is about to) a set of behavioral guarantees that have never been executed. | Any one of these five going the other way at implementation time (e.g. A10's subsumption interacting badly with `compareAtoms`' fixed sort tuple, or A12's package move hitting an unnamed cycle) forces a mid-implementation redesign of C1, not a bug fix — the contract itself was wrong. | §1, §4 |
| C-5 | 0030:§technical-design (Phase 2) | The implementation plan casually schedules a **package relocation** of `internal/guard/grammar.go::Evaluator` plus `declaration.go::intWidth` into `internal/resolve`, repointing three non-test call sites and rewriting ~25 references across seven `_test.go` files — inside a record whose Approach section frames the change as "one choice-point kind in `expand`." The scale mismatch between the Approach's plain-English framing and the actual blast radius is not an oversight; it's flagged ("Accepted scope … named here rather than discovered at implementation") but never re-scored against Proportionality's contract-count test. | The 3-week "add a grammar form" ships as a multi-file cross-package refactor; review and QA are sized off the Approach section, not the actual diff. | §1, premortem |
| C-6 | 0030:§consequences, C1 | Authored `enum` `domain` order silently acquires load-bearing semantic weight (step order) the moment any rule steps it — with **no lint, no warning, at declaration time**, that reordering the domain is now a breaking change. The record documents this as an authoring-guide hazard, not a load-time or lint-time check. | An author reorders an enum's `domain` for readability (alphabetizing, grouping) months after a `step` rule was added, with zero tooling feedback; the model silently changes runtime behavior — steps now point somewhere else. | §2, §3, premortem |
| C-7 | 0030:A5, C1 (`unless` not consulted) | `unless` atoms on a stepped tag are defined to be inert at admission time — a design choice that produces a *dead row* (advisory-only `graph-redundant-row`, never a blocking finding) rather than a refusal, when an author's intuitive `unless tier = "large"` does nothing. This is the single most natural thing an author reaches for (excluding a case with `unless`) and the record documents it will silently not work. | Author writes `unless tier = "large"` expecting it to exclude the terminal cell from stepping; the model loads clean, lints only an *advisory* (easy to ignore/suppress), and the row that "looks like" it protects the boundary is dead weight while the real boundary bug (an out-of-domain step) fires elsewhere or is masked by another row. | §2, §3, premortem |
| C-8 | 0030:§problem-statement, §decision-rationale (row 6) | The chosen design explicitly trades away the stated review benefit for two of six shipped surfaces: dump and graph now show `retry#0 … retry#4` (N rows), not one row — the review cost the Problem Statement says "lands twice" is only halved, not eliminated, for anyone auditing via `dump`/`graph` rather than raw source. | A reviewer using `intrastate dump` or `graph` to sanity-check a table — the tool this project explicitly built for reviewability — sees the exact unrolled-row explosion the RDR promises to eliminate; the "readability win" is source-only. | §2, decision-rationale |
| C-9 | 0030:A4, C1 (third-arm exclusion) | The admits filter treats the runtime evaluator's UNDECIDED verdict as "exclude, never admit," documented as unreachable for the currently-admitted kinds/operators — a soundness argument that depends on the operator/kind matrix never growing. The record itself flags (§consequences) that extending the guard grammar with an ordered/`neq` atom on enums is "the natural successor" once the enum residual (C-3) hurts. | The moment someone ships that natural extension, the "unreachable" third arm becomes reachable, and C1's admits filter's soundness argument — verified today via A4 — silently stops holding until someone re-derives it. A future RDR author, unaware the exclusion logic assumed unreachability, ships a regression that only shows up as models suddenly admitting fewer cells than expected. | premortem, §3 |
| C-10 | 0030:§finalization-gate (Proportionality) | The Proportionality gate's own split test ("sole author of at most one independent load-bearing contract") is in tension with C1 alone carrying: a new write-value grammar, a load-time admits filter reimplementing runtime evaluation semantics, a package relocation of a shared evaluator, AND a redefinition of suffix/identity ordering (D-identity) that narrows an existing RDR 0002 contract (C13). That is at least three independently-shippable seams bundled under one contract id. | Nothing user-visible immediately — this is a process/maintainability risk: the next RDR that needs to touch identity ordering, or the evaluator location, or the admits filter, has to re-litigate all of 0030 to change one part of it. | §2 |

## 1. The three most likely ways implementation goes wrong

**(a) The Phase 2 package move breaks something the RDR's own audit didn't census.**

Root cause: C1's admits filter requires reusing `internal/guard/grammar.go::Evaluator` from the loader, which today lives in a package (`internal/guard`) that imports `internal/table` — the very package doing the reusing. The record's fix is not "write a new evaluator" but "relocate the existing one, plus `intWidth`, to `internal/resolve`, and repoint three non-test call sites plus roughly 25 test references across seven files" (§technical-design, Phase 2). This is scored in the Decision Rationale matrix as merely a "behaviour-preserving move," and Alternative 1 is rejected partly *because* it avoids exactly this kind of package churn elsewhere — but the chosen approach pays its own churn cost here, just relabeled a "relocation" instead of a "widening."

The specific passage that enables this: A12 is graded `Pending`, not `Verified` — "the move is new at this pre-lock pass ... it resolves when Phase 2 compiles." A compile-time-only verification for a structural refactor of a shared evaluator, gated behind `go build ./...` and a `go list -deps` check, is a weak assurance for a move that touches `guard/product.go::valueSatisfies`, `graphlint/reach.go::atomAdmitsValue`, and `cli/flow_resolve.go::guardSeam` simultaneously — three call sites, each required (per the audit table) to preserve its OWN third-arm disposition after the shared core moves out from under it. The record's own audit table flags this as the delicate part: "the extraction is therefore of the two-valued core only... which is the condition under which repointing is behaviour-preserving." A single missed call site, or a test whose reflection assertion (`guard_evaluator_0003_test.go`, checking the evaluator type holds no state) doesn't follow the type cleanly to its new package, produces a regression in guard evaluation that has nothing to do with `step` and everything to do with an incidental refactor bundled into this RDR.

Symptom: guard evaluation subtly diverges between the loader's admits filter and the runtime kernel/lint paths after the move — a class of bug the RDR explicitly says it exists to prevent ("a divergence from the runtime is a load refusal rather than a lint finding") but is introduced by the very mechanism meant to prevent it, during a mechanical refactor nobody is watching as closely as the new grammar.

**(b) The `enum` cap-exclusion footgun (C-3 above) ships and is discovered by a real user, not by review.**

Root cause: C1's step mechanism has no analog to an `int`'s `lt`/`lte` bound for enums — ordered comparison operators are `int`-only in the frozen operator/kind matrix. The record itself states, in §Trade-offs → Consequences: "for `enum` the Problem Statement's 'every change to the tier domain means re-typing' is reduced ... but not eliminated: the write is said once, the cap is not." The specific passage: "the only positive atom that can exclude an `enum`'s terminal member is `in = [<every member but the last>]` — a re-listing of the domain, edited on every domain change."

This means the flagship example in the RDR's own Background/Illustrative Code (`tier` stepping `small → medium → large`) ships with exactly the footgun the RDR was written to eliminate, for the harder of its two motivating cases. The record acknowledges this and defers the fix ("the natural successor... out of scope here").

Symptom: an author adds a new tier to a domain (a routine, low-risk-feeling edit) months after the `step` rule shipped, and does not update the `in` re-listing that excludes the *old* terminal member — because nothing tells them to. The model now either refuses to load (if the old terminal is caught by C2's bound check against the new domain — but the new member has just been ADDED, so the risk is actually the inverse: the step now silently walks into the new tier where the `in` exclusion previously stopped it at the old terminal) or, worse, quietly changes which state the escalation ladder lands in. This is a runtime behavior change with zero load-time signal, in a system whose entire design ethos (per RDR 0002's charter, cited throughout 0030) is "the table means what it says, reviewably."

**(c) The lint-side quartic blowup (C-1) is hit by an ordinary author, not a stress-test.**

Root cause: A7 and Performance Expectations are candid that load is cheap and linear but lint is "roughly quartic" and unbounded from the load path by design (package graph enforced — `internal/table` cannot see `internal/guard`/`internal/graphlint`, so no ceiling is reachable at load). The record's mitigation is "degrade the verdict" (`complete=false`) rather than refuse the model — but that degrade only happens *if lint finishes at all*; the RDR's own numbers show lint taking beyond 600 seconds at 1,000 rows. A `step` write is, by construction, the easiest possible way to *accidentally* generate hundreds of rows from one rule — that's the entire point of the feature. Two stepped tags on one rule (the record's own MVV fixture pattern, e.g. `attempt` × `tier` both stepped) multiplies admitted cells across dimensions.

The specific passage: "a step write over a wide domain yields a model that loads instantly and that `lint` cannot finish analysing" — stated as an accepted risk, mitigated only by a mitigation whose enforcement lives in a *different, Draft* RDR (0013) that 0030 explicitly declines to depend on or gate against ("Not reusable (A7)... none: the A7 spike measured load linear and cheap, so no load ceiling is owed").

Symptom: a model author with a `step` rule over an `int` tag with `min=0, max=500` (not exotic — a percentage-like counter, a page count, a retry budget in a system with generous limits) runs `intrastate lint` and it never returns, or the CI job times out, with no actionable signal about *why* — because the RDR's answer to "why is lint slow" is a separate, not-yet-final RDR's problem to solve, if it ever gets picked up.

## 2. The section that will be rewritten within 6 weeks of shipping

**§Trade-offs → Consequences, specifically the enum residual paragraph, and C1's enum arm.**

This is the section the record itself flags as unsatisfying — "the enum arm delivers the stepped VALUE once, not the cap once" — and immediately proposes its own successor ("Extending the guard grammar with an ordered or `neq` atom over an authored `enum` domain would close it ... it is the natural successor if the residual proves to hurt"). This is a design that ships *knowing* its own half-measure and *naming* the fix it isn't doing. That is close to a self-fulfilling prophecy: the record predicts its own follow-up RDR. Combined with the fact that tier-escalation (the enum case) is literally the second half of the motivating example in the Background section — not an edge case, but one of the two examples the whole RDR exists to serve — the pressure to close this gap will land almost immediately once any real ladder with more than two tiers is authored. The rewrite will either widen the operator/kind matrix (reopening RDR 0003's "frozen vocabulary," which 0030 explicitly declines to touch) or add an enum-specific `step`-adjacent bound syntax, and in either case it revisits C1 and the operator model both.

## 3. The assumption that will not survive first contact with a real user

**A2/§consequences: "Authored enum `domain` order acquires meaning (step order) where today it has none; reordering a domain that a rule steps changes the model."**

This is stated as a documented hazard, addressed purely by prose in the authoring guide, with no load-time or lint-time detection. It fails a basic test any experienced engineer applies to "silent semantic coupling introduced by refactoring-looking edits": *if the tool cannot tell the difference between a safe edit and a breaking edit, users will make the breaking edit thinking it's safe.* Reordering a domain list is exactly the kind of edit a linter, formatter, or "let me alphabetize this for readability" pass would make without a second thought — and once any `step` rule references that tag, the record admits the reorder silently changes runtime successor states. No lint finding, no load refusal, nothing. This is the single assumption in the record most likely to produce a production incident that traces back, hours later, to "someone reordered a TOML array."

## 4. Premortem

*(Written as if the failure has already happened — six months post-ship.)*

Three incidents landed within the first quarter after `step` shipped, and all three trace to gaps this record named and accepted rather than closed.

**Incident 1 — the "silent tier skip."** A model author on the payments-retry team added a fourth escalation tier (`"critical"`) to an existing `tier` enum stepped by an `escalate` rule. The existing rule excluded the prior terminal member (`"large"`) via `match.tier in = ["small", "medium", "large"]` minus... no, via the documented pattern: `in` listing every member but the last. Adding `"critical"` to the domain did not touch that `in` list — nobody remembered it needed to, because nothing in `lint`, `load`, or the authoring guide's tooling (only its prose) flagged the coupling. `"large"` was no longer the terminal member; the step rule silently began admitting a cell into `"critical"` that the author never intended to be reachable from that transition. `graph-coverage-gap` did not fire because the cell was, technically, claimed — just claimed by the wrong intent. The bug was caught three weeks later by a support ticket: production escalation flows were reaching a `"critical"` alert tier from paths that were supposed to terminate at `"large"`. Root cause, once found: `internal/table/normalize.go::expand`'s step choice-point, exactly as C1 specifies, indexed positions into an enum domain whose order had quietly become load-bearing (C-6, C-3).

**Incident 2 — the CI timeout nobody could explain.** A platform team modeling a multi-stage rollout process authored a `step` rule over an `int` progress counter (`min=0, max=400`, representing rollout percentage in quarter-points) crossed with an enum `stage` tier also stepped in the same rule family. Both were "reasonable" declarations — nothing like the RDR's stress-test fixtures. `intrastate lint` in CI began timing out. The team's first three debugging sessions assumed a CI infrastructure problem, because `intrastate load`/`intrastate lint --model x` on a laptop "eventually" finished (in about nine minutes) and nobody thought to check row counts, because the source file was maybe sixty lines and eight rules — exactly the "one row per intent" promise the RDR delivered on load. The actual row count post-expansion, feeding lint's quartic analysis, was in the low thousands. This was foreseeable — the RDR's own Performance Expectations section states lint crosses ten minutes near 1,000 rows — but no tooling connects "your eight authored rules expand to 2,400 rows" to "lint will time out," because that ceiling belongs to a different, still-Draft RDR (0013) that 0030 declined to gate against (C-1).

**Incident 3 — the ambiguous `flow resolve` debugging session.** An on-call engineer investigating why a state machine seemed "stuck" repeating the same retry outcome used `flow resolve --as json` to inspect what rule fired. The payload named `rule: "retry"` — the authored id — giving no indication of which of the five expanded cells (`retry#0` through `retry#4`) actually matched, because `flow next`/`flow resolve` intentionally do not surface the suffix (C3, by design, per the Decision Rationale's explicit rejection of "publishing the suffix on `rule`"). The engineer had to separately run `dump` and cross-reference guard state by hand to figure out which literal cell had fired — a debugging step the RDR's decision rationale acknowledges ("a consumer wanting per-row identity from a flow payload reads `graph` or `dump`") but that isn't documented anywhere an on-call engineer would find it at 2am, and that doesn't exist at all for the `in`-expansion precedent the record leans on to justify the asymmetry (C-2).

All three incidents share a pattern: each is a gap the RDR itself documents, in its own prose, as an accepted trade-off or a "natural successor" — not a defect nobody saw coming, but a defect everybody saw coming and shipped anyway, betting that "tens of rows" and "an author who reads the authoring guide" would hold in production. They didn't.

## 5. Acceptance tests that would have caught each failure at RDR-review time

```gherkin
Feature: step grammar safety nets the RDR declines to build

  Scenario: reordering a stepped enum's domain is flagged
    Given a model with an enum tag "tier" whose domain is stepped by a rule
    When the model's authored domain order is changed (a reorder, not add/remove)
    Then "intrastate lint" MUST emit at least an advisory finding
      naming the tag and the rule that steps it
    # Currently: no finding exists for this at all (C-6)

  Scenario: adding a domain member to a stepped enum does not silently
            widen an existing step rule's reach
    Given an enum tag "tier" = ["small", "medium", "large"], stepped by
      a rule whose only cap protection is `in = ["small", "medium"]`
    When a fourth member "critical" is appended to the domain
    Then "intrastate lint" MUST flag that the rule's cap exclusion no
      longer covers the new terminal member
    # Currently: silently admits the new cell (C-3, C-6, Incident 1)

  Scenario: lint has a load-time-visible cost estimate before it runs
    Given a model whose step rules expand to N literal rows at load
    When N exceeds a documented threshold
    Then "intrastate load" or "intrastate lint --dry-run" reports an
      estimated lint cost or row count BEFORE a full lint run is attempted
    # Currently: no such signal exists; users discover cost via timeout
    #   (C-1, Incident 2)

  Scenario: an author-written `unless` on a stepped tag that does nothing
             is surfaced as more than an easily-ignored advisory
    Given a step rule with `unless tier = "large"` intended to exclude
      the terminal cell
    When the model is linted
    Then the finding severity communicates that the `unless` had NO
      effect on cell admission, distinctly from an ordinary redundant row
    # Currently: `graph-redundant-row`, an undifferentiated advisory (C-7)

  Scenario: flow resolve exposes which expanded cell fired
    Given a step rule expanded into N rows, one of which is selected by
      "flow resolve"
    When the JSON payload is inspected
    Then the payload names which cell/suffix fired, not just the bare
      authored rule id
    # Currently: `rule` is bare; suffix identity is dump/graph-only (C-2)

  Scenario: the guard evaluator relocation preserves identical verdicts
            across all three call sites, asserted as a single cross-package
            contract test, not by per-file inspection
    Given the moved `Evaluator` in `internal/resolve` is called from
      `guard::valueSatisfies`, `graphlint::atomAdmitsValue`, and the new
      loader admits filter
    When the same (atom, candidate value) is evaluated through all three
      paths
    Then all three return identical three-valued verdicts, verified by
      an automated contract test that fails the build if any path diverges
    # Currently: verification is "go build ./..." plus a dep-graph check,
    #   not a behavioral contract test across all three call sites (C-4, §1a)

  Scenario: RDR locks with zero Pending assumptions
    Given the Finalization Gate's Assumption Verification checklist
    When any Critical Assumption's Status is "Pending"
    Then the RDR does not lock; the MVV runs and settles every assumption
      status to Verified/Refuted before Final
    # Currently: A8 (Refuted, narrowed), A10-A14 (five, Pending) are
    #   carried into lock territory with "resolves at the MVV" (C-4)
```

## Closing note

This is an unusually well-instrumented RDR — the assumption ledger, mini-check tables, and MVV design are more rigorous than most records of this size get. That rigor is precisely what makes its self-documented gaps damning: this is not a record that failed to notice its weak points. It is a record that noticed all nine of them (C-1 through C-9 above are each named, in the record's own words, as accepted risk, residual cost, or "natural successor" work) and chose to ship anyway. The critique is not that the authors were sloppy; it's that "we said so in the Trade-offs section" is not a substitute for closing the gap, and several of these gaps (enum cap exclusion, domain-reorder hazard, lint cost cliff) sit exactly on the path of the RDR's own headline example.
