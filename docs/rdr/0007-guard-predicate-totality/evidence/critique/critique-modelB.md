Model: claude-fable-5

# Critique — RDR 0007 Guard Predicate Totality (model B)

## Findings Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| B-1 | §Risks "Residual status: UNMITIGATED"; Phase 2 "exported conformance harness" | The conformance harness binds no one: RDR 0003 is Final and unamendable, its text names no harness, and nothing in this repo fails if 0003's implement stage never calls it. Drift ships green. | RDR 0003's evaluator folds absence into false for one operator; every test in the tree is green; masked plans return in production. | §1, premortem, AT-1 |
| B-2 | §Consequences "the tracking kata (`xg7p`) is resolved by the contract, not by the behavior … the honest status of the *defect* is open even when the status of the *RDR* is Implemented" | The defect is closed in the tracker while the RDR's own masking probe still reproduces on the shipped kernel after Phases 1–3. | The exact pre-RDR behavior — `Escaped:true` plan masking missing artifact state — continues; the finding is re-derived later against a kernel that never changed. | §1, premortem, AT-1 |
| B-3 | §Normative "SCOPE OF THE ATOM VOCABULARY"; A10 (Pending); Testing "Prerequisite for rows 4–8 (blocking)" | Every atom-level MUST quantifies over structure (`operator, referenced tag key, literal, placement`) that no shipped or specified seam carries; `Row.Guard` is an opaque string. Phase 2's "vector-local and non-normative" encoding becomes the de facto normative mapping — the exact drift-by-implementation the RDR forbids. | Two evaluators legally disagree about what the same guard string means; the conformance vectors test an encoding nothing else ships. | §1, AT-3 |
| B-4 | §Normative "TOTALITY OF THE MAPPING … An evaluator MUST surface it … outside the verdict channel" vs the pinned seam `Evaluate(guard string, view TagSet) GuardResult` | The mandated out-of-band report has no channel: `Evaluate` returns a bare `GuardResult`, the same block pins that as "deliberate," and the evaluator cannot reach `Resolve`'s error return without panicking. | Evaluator authors either panic inside the pure kernel path or quietly return `GuardUnevaluable` for parse failures — the exact conflation the clause exists to forbid. | §1, AT-4 |
| B-5 | Testing Scenario 8 ("Expected: exactly one row selected in each case") vs §Normative "SURVIVOR MEMBERSHIP … a decided-GuardTrue sibling row MUST NOT be selected while an unevaluable candidate exists" | The blessed two-row absence pattern self-vetoes: with X absent, the value row's `X eq v` atom is unevaluable, the row survives as undecidable, and the aggregation veto refuses — the exists-guarded row is never selected. | Author writes the documented pattern; every absent-X resolution refuses `guard_unevaluable` instead of selecting the absence row. The pattern works only when the tag is present. | §2, AT-2 |
| B-6 | A5 (Pending, "load-time half is refuted pending A12"), A12 (Pending), A14 (Pending); §Risks "this mitigation stands or fails with A12 … it is a prohibition" | The entire authoring story rests on a stack of Pending assumptions, while the normative existence clause states A14's answer (`all … exists = false` is legal) as settled fact — violating the RDR's own Finalization Gate status-consistency rule. | Authors have no working absence idiom; the sanctioned fallback is the sentinel-stamping anti-pattern, so guards decide on placeholder values and partiality is defeated wholesale. | §2, §3, AT-9 |
| B-7 | §Background "An overlap key absent entirely yields `no_match` — the row genuinely does not match"; kernel `TagSet.matches` (`!ok → false`) | Closed-world absence-as-false is blessed in the `match` block while forbidden in guards. Predicate placement — `[rule.match.<tag>]` vs `[rule.guard.all.<tag>]` — silently decides whether absence is escapable `no_match` or a table-wide non-escapable halt, and no RDR states a placement rule. | The masking probe reproduces with zero guards involved: a row whose match pattern references the absent tag drops out of candidacy, `no_match` fires, the modeled escape rescues, plan emitted. | §3, premortem, AT-5 |
| B-8 | Problem Statement "told plainly that the artifact state needed to decide was missing" vs §Failure Modes "Diagnostic gap … it does not name the absent tag directly" | The refusal the RDR ships carries one row's opaque guard text and row refs, never the missing tag key; the promised actionability is deferred to RDR 0003 Phase 3 diagnostics — an obligation B-1 shows cannot be bound. | `guard_unevaluable guard="…"` with no tag name; the operator cannot tell which state was missing, whether to retry, or what to fix. | §1, premortem, AT-6 |
| B-9 | A13 (Pending); §Failure Modes "Refusal defeated by caller-supplied state … unroutable *around within the table*, not unbypassable by the caller" | The "non-escapable" guarantee is bypassable from the command line: supplying the missing key as an observed tag converts the refusal into a decided verdict. ADV-4/ADV-5's defect reintroduced one seam over, acknowledged, left Pending, untested. | Operator "fixes" the refusal with `--observed`; the plan's `Writes` derive from CLI input rather than artifact state, with nothing in the output marking the difference. | premortem, AT-7 |
| B-10 | A6b (open); §Failure Modes "Conflated recovery signal … 'inspect the artifact and accessor,' not 'retry'" | Transient read failure and genuine absence are byte-identical at the kernel boundary; the read-completeness obligation is routed to RDR 0004, which is also Final and unamendable, so it has no owner. | Refusal storms against a healthy artifact after a flaky accessor read; the documented guidance is to never retry, so operators hand-inspect artifacts for state that was there all along. | §1, premortem |
| B-11 | A3 "**Correction found at Resolve**: … the Normative Contracts aggregation block was corrected to the shipped behavior rather than the kernel changed" | Adding a guarded escape row converts an escapable `no_match` into a non-escapable `guard_unevaluable` (`escapeOrRefuse` returns the escape set's blocking refusal). The RDR's design position was refuted by frozen code and the spec rewrote itself to match — implementation-captured spec. | An author who models an escape edge gets a strictly harder refusal than an author who modeled nothing; the diagnosis names the escape row's guard, not the candidate condition being escaped. | premortem, AT-8 |
| B-12 | §Normative "The ordering is pinned to shipped behavior, NOT to its original rationale, which this RDR's narrowing invalidates" | The owned-before-unevaluable precedence is retained after its own justification ("more precise diagnosis") is admitted dead; a row failing both ways reports two independent defects serially. | Operator fixes the owned-state problem, re-runs, and only then learns the guard was also undecidable — two production round-trips per row. | §2 |
| B-13 | §Decision Rationale "how *often* the veto fires … an empirical question with no installed base to answer it"; §Investigation "on both of those points the citation cuts *against* the current design" | The resolution-wide veto — one unreadable row refuses a table whose other rows decide cleanly — is adopted at unknown frequency, and the only external prior art (SCXML) contradicts both the halt and the payload omission; the two most user-visible choices stand on in-house argument alone. | A 30-row table with one guard over a rarely-present observed tag refuses every resolution; the "accepted cost" is paid by every user of the table, on every run. | §2, premortem |
| B-14 | A9 (Pending); "Open ownership question … A normative MUST cannot carry its own strike-out condition" | The empty-`unless` identity is published as a normative MUST while its ownership is explicitly unresolved; if RDR 0002's normalization answers the other way, the clause is either duplicated or wrong. | An omitted `unless` block disables every row carrying one — the silent whole-table failure A9 itself names — or the identity ends up stated in two places that can diverge. | AT-9 |
| B-15 | Metadata "Profile: foundational — one contract" vs the nine ```normative``` blocks | The RDR locks the atom domain rule, `exists` semantics, provenance-blind presence, empty-block identities, strong-Kleene combination, refusal ordering, `Refusal.Rows` payload shape, `RequiresOwned` narrowing, and the aggregation veto as one seam, against its own Proportionality split test ("contract count, not word count"). | Any single clause failing (B-5 alone suffices) reopens the entire lock; there is no independent revision path for the authoring story vs the evaluation domain. | §2 |

---

## 1. The three most likely ways implementation goes wrong

### 1a. The contract enforces nothing, so the drift it exists to prevent ships anyway

**Root cause.** The RDR's entire delivery is Phase 1 (doc comments), Phase 2 (a test harness), and Phase 3 (a "handoff" of obligations to RDR 0003). But RDR 0003 is Final, this project never amends RDRs, and the RDR admits — in its own Risks section, in bold — that the residual status is "**UNMITIGATED, not partially mitigated.** Nothing in this repo fails if RDR 0003 never instantiates the harness." The enforcement point is a future implement stage's voluntary acceptance of an obligation recorded in a *different* document's Phase 3 notes. There is no mechanism: no failing test, no lint, no build break, no prerequisite checkable by tooling. `go test ./internal/resolve/` stays green against `fixtures_test.go::fixtureGuards` regardless of what 0003 builds — the RDR says so itself.

**Enabling passage.** §Risks: "the acceptance must happen at 0003's *implement* stage (which reads its Prerequisites), not by editing 0003's text — and if that stage declines, the drift risk is carried openly rather than recorded as closed." This is a design document pre-authorizing its own core guarantee to fail, and calling the failure "carried openly." Compounding it, §Consequences instructs that kata `xg7p` be closed "against the contract, not the behavior" — while conceding "this RDR's own probe still reproduces on the shipped kernel after Phases 1–3 land."

**Symptom the user sees.** Nothing changes, which is the symptom. The masking probe from the Problem Statement — FALSE-folded absence, `no_match`, modeled escape, `Escaped:true` plan — reproduces byte-for-byte after the RDR is marked Implemented. Months later, a 0003-built evaluator that folds absence into false for `contains` (the one operator whose set-theoretic reading invites it, per the RDR's own Scenario 7) passes every test in the tree. The user gets a plan where the honest answer was "I could not tell," which is the opening sentence of the Problem Statement, verbatim, after two RDRs and a kata closure.

### 1b. The normative clauses bind a seam that cannot see what the rules are about — and forbid the only escape

**Root cause.** Every atom-level clause quantifies over operator, referenced tag key, literal, and `all`/`unless` placement. At the only seam the RDR binds — `internal/resolve/resolve.go::GuardEvaluator.Evaluate(guard string, view TagSet)` — a guard is an opaque string the kernel attaches no meaning to, and no RDR anywhere specifies the mapping from authored TOML structure into that string. The RDR knows this: A10 is Pending, flagged as "the hinge finding," and Testing rows 4–8 — the rows carrying "this RDR's actual content" — are marked unencodable until A10 closes. The proposed resolution is to hand the mapping to RDR 0003, which (see 1a) cannot receive it. Meanwhile Phase 2 must invent a "vector-local and non-normative" guard encoding to test against. There will then be exactly one executable definition of guard structure in the repository, living in the conformance vectors, and it will be treated as normative by every future implementer because it is the only thing that runs — "deciding A10 by implementation, which is the drift this RDR exists to prevent," in the RDR's own words, as a prediction it fulfills.

The same block contains a sharper defect: "a well-loaded guard the evaluator cannot parse … MUST NOT be reported as GuardUnevaluable … An evaluator MUST surface it … outside the verdict channel." There is no outside. `Evaluate` returns a bare `GuardResult` — the RDR pins this as "deliberate and not an oversight" in the adjacent paragraph — and the evaluator has no path to `Resolve`'s error return. The clause mandates behavior the pinned signature cannot express. The only implementations are `panic` (inside a kernel whose doc comment promises a pure decision boundary) or returning `GuardUnevaluable` anyway (the forbidden conflation).

**Enabling passage.** §Normative Contracts, "TOTALITY OF THE MAPPING IS LOAD-BEARING": "`Evaluate` has no error return (below), which is deliberate and not an oversight: the seam reports verdicts, and a mapping failure is not a verdict." A verdict-only seam plus a MUST-NOT on the only in-band verdict plus no out-of-band channel is not a contract; it is a trap laid for the first implementer.

**Symptom the user sees.** Either a panic surfacing as a Go error from `Resolve` on a table that passed lint (a crash where the taxonomy promised a typed refusal), or — the likelier, quieter path — evaluator bugs reported as `guard_unevaluable`, sending operators to inspect artifacts for missing state that was never missing. The refusal cannot distinguish the cases (B-8), so the misdirection is undetectable from the output.

### 1c. Authors have no legal way to express absence, so they route around the rule

**Root cause.** The rule makes value operators partial, so every legitimate "row applies when X is absent (or defaulted)" needs the existence operator. The RDR's own record shows every route is broken or unproven: the single-atom route (`all … exists = false`) rests on A14, Pending, with the repeatability lens already reconstructing `exists` as unary; the two-row route rests on A12, Pending, with the RDR itself noting that RDR 0003's overlap proof runs over a domain product "that cannot represent absence" and that `ambiguous overlap` is a mandatory RDR 0002 rejection with no downgrade valve; and — worse than either open question — the two-row pattern as written in Scenario 8 is refuted by the RDR's own aggregation veto (see §2 below): the value row goes unevaluable precisely when the tag is absent, vetoing the absence row. The Risks section then concedes the endgame: "an anti-pattern warning with no working alternative is not a mitigation, it is a prohibition, and it would make the sentinel path the only route to an absence-conditional row."

And the sentinel path has a cheaper sibling the RDR blesses in its own Background: put the predicate in the `match` block. `TagSet.matches` returns false on an absent key — closed-world absence-as-false, escapable `no_match`, exactly the masking semantics the RDR forbids for guards — and the RDR endorses it ("the row genuinely does not match"). RDR 0002 authors match predicates and guard predicates in adjacent TOML tables with no rule about which predicate belongs where. Authors who hit the guard-side refusal will move predicates into `match` and get the old behavior back.

**Enabling passage.** §Background: "An overlap key absent entirely yields `no_match` — the row genuinely does not match, per the `req-list.md` ASSUMPTION's own exclusion clause." One paragraph after describing absence-as-false as "the path where missing artifact state is masked," the same semantics one field over is ratified as correct, and the RDR never returns to the asymmetry.

**Symptom the user sees.** Tables accrete `status = "none"` sentinels and match-block predicates that silently drop rows on absent state. The domain rule holds perfectly over guards nobody writes anymore. Missing artifact state is masked behind `no_match` escapes — through the front door this time.

## 2. The one section that will be rewritten within 6 weeks of shipping

**The existence clause + Scenario 8 + A5 cluster — the authoring story — and it will drag the SURVIVOR MEMBERSHIP block's veto text with it.**

It will be rewritten because it is internally contradictory today, mechanically, with no empirical input needed. Walk the RDR's own clauses: row 1 guarded `unless … exists` on X, row 2 guarded `X eq v`, X absent. Row 1: `exists` is total, decides FALSE on absence, `unless_conj = F`, `¬F = T`, row 1 is GuardTrue. Row 2: `X eq v` is a value-comparing atom over an absent key, unevaluable by the first normative block; per SURVIVOR MEMBERSHIP, row 2 "is a SURVIVOR"; per the aggregation clause, "if any surviving candidate row's guard is GuardUnevaluable, the resolution MUST refuse guard_unevaluable — a decided-GuardTrue sibling row MUST NOT be selected while an unevaluable candidate exists." Scenario 8's expected value — "exactly one row selected in each case" — is unsatisfiable in the absent case. The pattern the RDR blesses as the sole legitimate outlet for absence-conditional rows is refused by the RDR's own veto, in exactly the case the pattern exists to serve. MVV Scenario 3 (*unevaluable-blocks-true-sibling*) is the identical shape, celebrated two pages earlier as the RDR's key new behavioral pin; the document never connects the two.

The only repair is to require every value-comparing atom over an optional tag to carry a companion `exists = true` atom in the same `all` block (so `F ∧ U = F` prunes via the exists atom — strong-Kleene makes it work). That requirement appears nowhere in the RDR, changes the blessed pattern, changes Scenario 8's fixture, changes A5's disjointness argument (which analyzes load-time assignment sets, not runtime three-valued verdicts — a category error the Carried Constraint half-notices for A12 and misses for the veto), and hands RDR 0003's lint a new obligation ("warn on a value atom over an optional tag with no exists companion") with, per B-1, no way to deliver it. First author to encode the pattern — plausibly Phase 2 itself, encoding Scenario 8 — hits the contradiction, and the rewrite follows. A12 and A14 resolving in any direction forces further rewrites of the same text; three Pending assumptions plus one live contradiction concentrated in one clause cluster is not a section that survives six weeks of contact.

## 3. The one assumption that will not survive first contact with a real user

**A5 — "Authors who need 'row applies when tag X is absent' … can express it with existence atoms … Neither route is ambiguous under exact-one selection."**

It is already three-quarters dead inside the document: its Status line says the load-time half is "refuted pending A12," its single-atom leg hangs on A14 (Pending, with a named dissenting reconstruction), and §2 above shows its two-row leg is refuted at *runtime* by the RDR's own aggregation veto — a refutation the RDR did not find because A5's disjointness proof works in RDR 0003's two-valued set semantics while the runtime it governs is three-valued. The first real user is a table author with an optional observed tag — the RDR's own reference fixture, `cluster_eligible`, is exactly this — who needs "absent or false ⇒ this row." They will try the documented pattern, get `guard_unevaluable` on every absent-case run, read the Failure Modes guidance ("rewrite the guard with an explicit existence atom if absence was intended" — which they already did), and conclude the feature does not work. Then they will do one of the three things the RDR cannot stop: stamp a sentinel upstream, move the predicate into the `match` block (B-7), or supply the tag with `--observed` (B-9). Each defeats the domain rule; the third produces plans computed from caller input. A5 is the assumption that authors will stay inside the sanctioned idiom; there is no working sanctioned idiom to stay inside.

## 4. Premortem

*Written 2027-05, as if looking back.*

RDR 0007 locked in September 2026 after A9, A12, and A14 were argued closed on textual readings of two Final peer documents. Phase 1 landed as doc comments on `Row.RequiresOwned` and `GuardEvaluator`. Phase 2 landed `resolveconform.RunGuardVectors` plus an unexported view-reading evaluator with a vector-local guard encoding — colon-delimited `op:key:literal` strings, chosen in an afternoon because A10 had been "handed to RDR 0003." Kata `xg7p` was closed against the contract. The RDR's masking probe still reproduced on `main` the day the status flipped to Implemented; everyone knew, because the Consequences section said it would, and that was treated as the section working.

RDR 0003's implementation began in November under a different session. Its launch prompt read 0003's own Prerequisites — which name no harness, because 0003 locked four months before the harness existed. The evaluator parsed the TOML guard structure directly and serialized rows into its own guard keys; nothing resembling `op:key:literal` ever reached it, so when an engineer eventually tried `RunGuardVectors` in week three, every vector failed at parse, the failure was read as "the harness encodes a stub format, not our format" — which was true, the RDR itself had stamped the encoding non-normative — and the harness was dropped with a TODO. The evaluator's `contains` treated an absent set-valued tag as the empty set, because that is what the implementer's set library returned and no red test said otherwise. `go test ./...` was green in both packages. B-1 and B-3, exactly as written.

First production contact came from the kata-flow table. A row guarded `cluster_eligible eq true` — copied from 0003's own reference fixture — hit a view where the observed tag wasn't supplied. Under the drifted evaluator this folded to false, pruned in `gate`, and the table's `no_match` escape emitted an `Escaped:true` plan: the Problem Statement's opening scenario, verbatim, now with a closed kata asserting it was fixed. Meanwhile the RDR-flow table, evaluated through a path that *did* route absence to `GuardUnevaluable`, produced the opposite failure: one row guarded on `prelock_iterations` refused every resolution of a thirty-row table whenever the accessor's read of the state file was truncated — A6b's conflation, `owned_state_unavailable`'s sibling storm. The refusal printed `guard_unevaluable guard="lt:prelock_iterations:3"` and two row refs. Nothing named the missing tag (B-8). The operator did what the refusal text made inevitable: re-ran with `--observed prelock_iterations=0`. It resolved. The workaround went into a runbook. Plans were now being computed from a CLI flag (B-9), and the runbook called it the fix.

The author who tried to do it right — model "absent or 0" with the blessed two-row pattern — filed the bug that unwound the section: her exists-guarded row was GuardTrue, her value row was unevaluable, and `gate`'s `len(undecidable) > 0` branch discarded the selection, exactly as the SURVIVOR MEMBERSHIP clause required and exactly as Scenario 8 promised would not happen (B-5). The fix — companion `exists = true` atoms on every value atom over an optional tag — was folded into a successor RDR, 0011, which re-stated the domain rule, the veto, and the authoring story in one place, superseding three of 0007's nine normative blocks. 0007's remaining novel content was, in the end, the two doc comments Phase 1 shipped. The post-mortem's first line: the RDR chose to be the "single normative home" of a rule whose implementer, enforcer, and diagnostic channel all lived in documents it was forbidden to touch, and every failure above was recorded inside the RDR itself — as a Pending assumption, an UNMITIGATED risk, or an accepted cost — before lock.

## 5. Acceptance tests that would have caught each failure at RDR-review time

**AT-1 (catches B-1, B-2): Enforcement is checkable, not narrated.**
```gherkin
Given RDR 0007 claims the conformance harness binds RDR 0003's evaluator
When I search RDR 0003's text and its implement-stage prerequisite inputs
Then some artifact that 0003's implementation MUST consume names the harness
And a mechanical check (test, lint, or gate) fails when it is not instantiated
And the kata xg7p closure text states the probe still reproduces, with a
  named successor tracking the behavioral fix
```
Fails at review time: no such artifact exists, and the RDR says so. Verdict should have been "do not lock until the enforcement point exists," not "carried openly."

**AT-2 (catches B-5): Desk-execute the blessed pattern against the RDR's own clauses.**
```gherkin
Given row R1 guarded "unless X exists = true" and row R2 guarded "X eq v"
And an assembled view where X is absent
When I evaluate both rows per the Normative Contracts (exists total,
  value operators partial, strong-Kleene, survivor membership, veto)
Then R1 is GuardTrue and R2 is GuardUnevaluable
And the aggregation clause requires refusal guard_unevaluable
And Scenario 8's expected value "exactly one row selected" is contradicted
```
A pencil-and-paper walkthrough; five minutes at review time. It was never done because A5's disjointness argument ran in two-valued load-time semantics.

**AT-3 (catches B-3): Every normative clause must be encodable as a concrete vector before lock.**
```gherkin
Given Testing Strategy row 4 (operator × presence × strong-Kleene)
When I attempt to write one vector as a literal Go value against the
  shipped seam Evaluate(guard string, view TagSet)
Then the guard input is expressible without inventing an encoding
```
Fails: no encoding exists. The RDR knows (rows 4–8 "cannot be encoded") and locks anyway with the blocker as a Pending assumption. The test converts "Pending" into "not lockable."

**AT-4 (catches B-4): The mapping-failure clause must name a real channel.**
```gherkin
Given an evaluator handed a well-loaded guard it cannot parse
When it obeys "MUST NOT return GuardUnevaluable" and "MUST surface it
  outside the verdict channel"
Then there exists a compilable implementation of Evaluate that does both
  without panicking in the kernel's pure path
```
Fails: the signature admits no such implementation. The clause needed either an error return on the seam (a kernel change the RDR refuses) or an explicit "panic is the sanctioned surface" statement; it has neither.

**AT-5 (catches B-7): Placement parity probe.**
```gherkin
Given predicate "phase = prelock" authored once in [rule.match.phase]
  and once in [rule.guard.all.phase], with tag phase absent from the view
When both tables resolve against a table modeling a no_match escape
Then the two dispositions are stated by some RDR to be intentionally
  different, with a rule telling authors which block to use
```
Fails: match-placement yields an escapable `no_match` (via `TagSet.matches`), guard-placement a non-escapable `guard_unevaluable`, and no document states the placement rule. The masking path survives in `match`.

**AT-6 (catches B-8): Actionability from the payload alone.**
```gherkin
Given a guard_unevaluable Refusal for a guard referencing tags a and b,
  where only b is absent
When a user reads only the Refusal value (Kind, Guard, Rows, Revision)
Then the user can identify b as the missing state without consulting the
  table source or the evaluator
```
Fails: `Refusal` carries opaque guard text and row refs. The Problem Statement's promise ("told plainly that the artifact state needed to decide was missing") is the RDR's stated reason to exist; this test holds the refusal payload to it and the RDR's own Failure Modes section concedes the miss.

**AT-7 (catches B-9): The non-escapability claim is scoped or defended.**
```gherkin
Given a resolution refusing guard_unevaluable on absent key k
When the caller re-runs with k supplied via Input.Observed
Then either the kernel still refuses (test pinning the guard path like
  ADV-4/5 pin the owned path), or the RDR's masking-closed claim is
  restated as "within the table only" everywhere it appears
```
At review time this forces A13 from Pending to either a frozen test or honest language in the Approach section, which currently says "missing artifact state can never be masked" unqualified.

**AT-8 (catches B-11): Escape-row surprise probe.**
```gherkin
Given a table whose candidate set yields no_match and whose no_match
  escape row carries a guard over an absent tag
When the table resolves
Then the refusal kind, and the fact that ADDING the escape row changed
  an escapable refusal into a non-escapable one, is documented in the
  authoring-facing Failure Modes, not only inside assumption A3's Evidence
```
The behavior is frozen and arguably right; the defect is that its only statement is buried in a correction note inside A3, where no table author will find it.

**AT-9 (catches B-6, B-14): The RDR's own status-consistency gate, run at draft.**
```gherkin
Given the Finalization Gate rule "no assumption marked Pending may have
  settled-fact prose elsewhere in the RDR depending on it"
When I list normative-block sentences depending on A9, A12, or A14
Then the list is empty
```
Fails three times today: the existence clause states A14's conclusion as fact; the empty-`unless` MUST stands while A9's ownership is open; Scenario 8 and the sentinel mitigation state A12's conclusion conditionally at best. The gate exists in the document; it was not run against the document.

---

*15 findings. Independence note: authored without reading any file under `evidence/critique/` other than writing this one.*
