Model: claude-opus-5

# Critique — cli/0030, computed write value grammar

## Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `0030:C1` | The "placeholder in `assignments`" clause is incompatible with `renderWrites`' actual shape: `conform` runs on the authored value inside the write loop, and `assignments[key] []string` has nowhere to carry a step spec. The clause admits it "widens one signature pair" while in fact demanding the write pipeline be re-ordered so conformance moves from `renderWrites` to `expand` — two functions, two packages of test surface. | Phase 2 over-runs; the first implementation either conforms the placeholder (and a placeholder that conforms is a real literal a fixture can author) or skips conformance entirely on stepped keys, so an out-of-domain step loads clean. | §1, premortem, AT-1 |
| C-2 | `0030:C1` | `expand` has no tag declarations. Its signature is `expand(base Row, predicates []Atom, outcome Atom, writes []TagValue) []Row` and it is a package-level function with no `*loader` receiver, no `model`, no `TagDecl`. Every operation C1 puts in `expand` — domain enumeration, `intWidth`, per-cell `conform`, the `#`/duplicate-member guards — requires the declaration. The record names one signature change; the real one threads the model through `expand`. | none directly; it is the cost that makes the record's "one choice-point kind" estimate wrong by a factor. | §1, premortem, AT-1 |
| C-3 | `0030:D-identity` | "Step points order by KEY alone" is asserted total, but the sorted candidate list mixes step points (no `Atom`) with `Atom`-bearing choice points, including the outcome point whose key is the literal `"recognized"`. Two different comparators over one `slices.SortFunc` list is not a well-defined order; Go's sort is not stable and `compareAtoms` is a 4-tuple. A step on a tag named e.g. `region` and an `in` on `recognized` have no defined relative order under "key alone" vs "(key, block, operator, literal)". | Suffix element order — hence `Identity()`, hence `dump` rows and graph row ids — varies between builds or between Go versions. A model that dumped `retry#0#a` dumps `retry#a#0`. | §1, §2, premortem, AT-2 |
| C-4 | `0030:C3` | C3 asserts "every surface that names a ROW publishes the suffixed identity" and cites `dump` and the graph document. Lint is listed in C3's reader set but is not covered by the claim, and in fact **every** `graphlint` finding populates `Rule` from `row.RuleID`, never `Row.Identity()` — ten sites across `analysis.go`, `coverage.go`, `groups.go`. C3 is silent on this; the RDR never states that lint findings cannot name an expanded row. | A stepped ladder with a defect on one cell emits N findings all reading `rule: "retry"`. The author cannot tell which cell is broken. The unrolled ladder they replaced told them exactly (`retry-3`). The record's own selling point — review legibility — inverts at the exact moment review matters. | §1, §2, §3, premortem, AT-3 |
| C-5 | `0030:S1` | The MVV's central oracle projects findings to `(code, key, dimension, class, reason)` and explicitly drops `rule`, `span`, `element`, `fingerprint`. That projection is precisely the set of fields C-4 breaks. The one test the record leans on is constructed so it cannot observe the defect it most needs to observe. | The MVV goes green on an implementation that has made lint findings unattributable. | §1, §2, premortem, AT-3 |
| C-6 | `0030:A5` | A5 claims an `unless`-over-admitted interior cell yields a *dead row* lintable as advisory `graph-redundant-row`. `checkRedundantRows` fires only on a **proper subset** (`left.Subset(right) && left.Len() < right.Len()`), only against a sibling of the same `Kind()` in the same group, and only when **both** rows are `Projectable()`. An `unless`-carrying row is the canonical unprojectable row — `decidableAccepted`'s own comment says so. So the advisory is not merely fragile: on the common shape it does not fire at all. | The dead row is silent. It claims a cell at load, contributes nothing at runtime, and lint says nothing. This is the "silent no-op" the record cites `0002:C3` to forbid — reintroduced by the record's own `unless` exclusion. | §1, §3, premortem, AT-4 |
| C-7 | `0030:S5b` | S5b asserts "EXACTLY ONE advisory `graph-redundant-row`" and then concedes in the same scenario that the projectability precondition may not hold, asserting it "explicitly" as a precondition. A test that asserts its own precondition and only then asserts the outcome does not prove the outcome holds in the authored models users write — it proves it holds in the one fixture engineered to project. | The fixture passes; the user's real ladder does not fire the advisory. | §1, premortem, AT-4 |
| C-8 | `0030:§consequences` | The `enum` arm's residual is disclosed and then waved through: to exclude an enum's terminal member the author must write `in = [<every member but the last>]`, edited on every domain change. The Problem Statement's stated grievance is "every change to the tier domain means re-typing the unrolled block." The enum arm converts re-typing N rows into re-typing one N-member list. That is a smaller edit, not a different failure. | The tier author adds a fourth tier, forgets the `in` list, and the model refuses at load with a bound error pointing at a rule they did not touch. They conclude the feature does not work for enums. | §1, §3, §2, premortem, AT-5 |
| C-9 | `0030:C2` | The bound refusal is a **load** refusal that fires on a declaration-level interaction: widening `domain` or tightening `max` breaks a rule elsewhere in the file. The record calls this "existing refusals; no advisory is minted" — but no existing refusal is *action-at-a-distance across two declarations and a rule*. The remedy C2 prescribes (exclude with a positive atom) then hands the author a `graph-coverage-gap`, so one edit produces two errors in sequence. | Author edits `max = 9` to `max = 12`; model still loads (widening is safe). Author edits `domain` to add a tier; model refuses at load naming a rule they never opened. Two-step remedy, each step surfacing a different tool's error. | §1, §3, premortem, AT-6 |
| C-10 | `0030:A12` | The Evaluator **package move** — `internal/guard/grammar.go` to `internal/resolve`, with ~25 test references across seven `internal/guard/*_test.go` files, plus `intWidth`, plus repointing `guard/product.go::valueSatisfies`, `graphlint/reach.go::atomAdmitsValue` and `cli/flow_resolve.go::guardSeam` — is scoped as "sound and small" and "a relocation 0012 absorbs". It is a cross-package refactor of the guard seam, landed in the same phase as a new grammar feature, and A12's status is **Pending**. | A merge conflict class with RDR 0012 that the record declares away rather than sequences. If 0012 lands first the shim's signature assumption is void. | §1, premortem, AT-7 |
| C-11 | `0030:A10`, `0030:A11`, `0030:A13`, `0030:A14`, `0030:A5`, `0030:A6` | Six of fourteen Critical Assumptions are **Pending**, and four of them (A10, A11, A13, A14) are marked "new at this pre-lock pass." A10's own "If wrong" says both branches "reopen C1's expansion half" — i.e. C1 is contingent on an unverified assumption. Meanwhile C1 states the subsumption rule as settled normative prose. | The record locks a contract whose central mechanism is admitted-unverified. Implementation discovers the refutation, and C1 is rewritten post-lock. | §1, §2, premortem, AT-8 |
| C-12 | `0030:§finalization-gate` | The Finalization Gate is entirely unfilled template text — Contradiction Check, Assumption Verification, Scope Verification, Cross-Cutting Concerns and Proportionality are all bracketed placeholders. The gate's own Assumption Verification instruction says "no assumption marked `Pending` or `Unverified` may have settled-fact prose elsewhere in the RDR depending on it," which C-11 violates on six counts. | The record cannot lock as written, and the one mechanism that would have caught C-11 has not been run. | §1, §2 |
| C-13 | `0030:§proportionality` | The gate's Proportionality instruction requires "sole author of at most one independent load-bearing contract." This record authors a new write-value grammar (C1 grammar half), a new loader expansion mechanism (C1 expansion half), a new load-time domain-conformance check (C2), **and** a cross-package refactor of the guard evaluator seam (A12/Phase 2). That is at least two seams, plausibly three. | Scope creep is chartered rather than caught; the Phase 2 estimate is the one that slips. | §1, §2 |
| C-14 | `0030:§performance-expectations` | The record measures load (linear, cheap) and then reports lint is roughly quartic — 7.5 s at 100 rows, 116 s at 200, >600 s at 1,000 — and concludes no load-time ceiling is owed because the cost "lands on lint," which is 0013's scope. The user does not experience two budgets. They experience `intrastate lint` on their model. | Author writes `attempt = { step = 1 }` on `min=0, max=200` — an entirely reasonable retry counter. Model loads in milliseconds. `lint` never returns, or returns `complete=false`. The record predicted this and shipped it anyway. | §1, §3, premortem, AT-9 |
| C-15 | `0030:A7` | A7 concludes "no load-time ceiling is owed at the ladder sizes reported" from a spike at 100,001 rows and then pushes the real ceiling onto 0013's lint scope, which the same assumption proves is unreachable from load by the package graph. The record has proven that nothing in the system can stop the model between authoring and the quartic blowup. That is the argument *for* a load ceiling, presented as the argument against one. | See C-14. | §1, premortem, AT-9 |
| C-16 | `0030:§consequences` | "A `step` ladder and a hand-unrolled ladder are two authorings of one behaviour but not of one dump." The record then declines to extend `0002:§round-trip-inverse-invariants` to the pair. The RDR's whole motivating journey is an author *migrating* an existing unrolled ladder to step form. There is no migration story, no diff aid, no `intrastate` verb that shows the pair are equivalent. | The author converts a 24-row ladder to 8 rows, runs `lint`, gets a different-looking finding set (different `rule` names — C-4), and has no supported way to convince themselves the conversion was faithful. The MVV has that oracle. The user does not. | §1, §3, premortem, AT-10 |
| C-17 | `0030:C1` | "This NARROWS `0002:C13`'s suffix iff to the `in` expansion it was written about." A locked, Implemented record's normative clause is being narrowed by a Draft record's prose, in a paragraph that also admits "a reader of 0002 alone would otherwise be told something this record makes false." That is an amendment to 0002 performed inside 0030. | Future readers of 0002 hold a false invariant. Nothing in 0002 points at 0030. | §1, §2 |
| C-18 | `0030:C2` | C2 requires the refusal detail to name the rule, tag, cell and result, and to name any `unless` atom the rule authors on the stepped tag — all as free text in `Failure.Detail`, because `Failure` has no per-part fields and the record adds none. Tests assert detail text with a position argument ("the detail names the cell FIRST"). | The refusal message is a string contract asserted by substring position. Any rewording breaks tests; any test-driven wording freeze makes the message unimprovable. Machine consumers get nothing structured for the record's most important new error. | §1, premortem, AT-6 |
| C-19 | `0030:§decision-rationale` | "Per-row identity IN a flow payload is a possible successor record — disclosed here as a non-obligation, not an owed fix." Combined with C-4, the record has disclosed away the two surfaces (flow payloads, lint findings) on which a user diagnosing a stepped ladder actually works, and kept the claim only for `dump` and `graph` — the two surfaces a human reads least. | The author debugging a stepped ladder is told to stop using `flow resolve` and `lint` output and start reading `graph` JSON. | §1, §3, premortem, AT-3 |
| C-20 | `0030:A2` | A2 verifies that authored `domain` order survives today, and §consequences notes "reordering a domain that a rule steps changes the model" with **no lint to catch it** (C3 puts it in the authoring guide instead). The record creates a silent semantic dependency on source-file ordering and dispositions it to documentation. | A cosmetic alphabetization of a `domain` array — the kind of thing a formatter or a tidy-minded reviewer does — silently changes flow behaviour. No error, no finding, no diff signal. | §1, §3, §2, premortem, AT-11 |

---

## 1. The three most likely ways implementation goes wrong

### 1.1 The loader cannot do what C1 tells it to do, and Phase 2 discovers this

**Root cause.** C1 assigns the expansion to `internal/table/normalize.go::expand`. `expand` is a package-level function whose entire input is `(base Row, predicates []Atom, outcome Atom, writes []TagValue)`. It has no `*loader` receiver, no `l.model`, and therefore no `TagDecl`. Every single operation C1 places inside it needs a declaration: enumerating `{min..max}`, walking `domain` in authored order, applying `intWidth`, running the per-cell `conform`, checking for a `#` member, checking for a duplicate member. The record's own audit row calls this "one more choice-point kind."

**Enabling passage.** `0030:C1`: *"The stepped key is placed in `assignments` like any other written key, carrying a placeholder the expansion replaces per cell … The spec itself rides alongside as one more per-key record, since `::renderWrites` returns `([]TagValue, []string, error)` and `::expand` takes `writes []TagValue`, neither of which can carry a non-literal today — widening that pair is the one signature change this clause implies."*

That sentence is wrong twice. First, the placeholder must survive `conform(decl, "eq", members)` at `normalize.go:623`, which runs inside `renderWrites` on the authored member list — a placeholder that conforms is indistinguishable from a literal an author could have written; a placeholder that does not conform means the step path must branch around the one check C2 claims to reuse. Second, "widening that pair" is not one signature change: it is threading declarations into a function deliberately kept declaration-free, plus re-homing kind/domain conformance from the write loop to the expansion loop.

**Symptom.** The first working implementation takes the path of least resistance: it skips `conform` on stepped keys (because the placeholder cannot pass it) and defers the bound check to the per-cell loop — which is correct — but nothing in the record says the *authored* half of the write block is still checked. A model authoring `attempt = { step = 1 }` on a tag whose `max` was never declared, or on a tag whose declaration disagrees with its kind, takes a different code path from every other write. The refusals C1 promises are the ones the implementer remembers to re-add.

### 1.2 Suffix element order is not defined, so row identity is not stable

**Root cause.** `expand` builds one `[]choicePoint` and sorts it with `slices.SortFunc(candidates, func(x, y choicePoint) int { return compareAtoms(x.atom, y.atom) })`. Every candidate carries an `Atom`. D-identity says a step point **is not an `Atom`** and orders "by KEY alone." A single sort over a list where some elements are compared on a 4-tuple and others on a 1-tuple is not a total order and `slices.SortFunc` is documented as not stable.

**Enabling passage.** `0030:D-identity`: *"**Step points order by KEY alone.** A step point is not an `Atom`, so `compareAtoms`' `(Key, Block, Operator, Literal)` tuple does not apply to it … Key alone is total here because C1's subsumption leaves at most one suffix-emitting step point per key."*

The reasoning proves the *keys* are distinct. It does not prove the *comparator* is consistent. Two step points on distinct keys order fine. A step point on key `k` and an `in` point on key `k` are handled by subsumption. But a step point on key `region` against the outcome point whose key is the literal `"recognized"` (`internal/table/model.go:19`) orders by string comparison of `"region"` vs `"recognized"` — which happens to work — while a step point on `attempt` and an `in` on `attempt`'s *sibling context* atom orders by key then by a tuple the step point does not have. The record never writes down the comparator. Neither will the implementer, the same way.

**Symptom.** `retry#0#small` on one machine, `retry#small#0` on another. Row identity is what `dump`'s `identity` column and the graph document's `identity` field publish — the exact two surfaces C3 promises are the stable ones. The user's diff of a graph export churns for no source change.

### 1.3 The lint surface loses the ability to name the row it is complaining about

**Root cause.** C3 divides the world into surfaces that name a ROW (publish `identity`) and surfaces that name a RULE (publish `rule`). Lint is in neither list, and is in fact the surface that most needs the distinction. Grounded: every finding in `internal/graphlint` sets `Rule` from `row.RuleID` — `analysis.go:213`, `analysis.go:567`, `coverage.go:330`, `groups.go:80`, `groups.go:207`, `groups.go:318`, `groups.go:363`, `groups.go:399`, plus `coverage.go:91`/`:140` which use `firstRuleID(g)`. None calls `Identity()`.

**Enabling passage.** `0030:C3`: *"Every surface that names a ROW publishes the suffixed identity `rule#cell` — the shape `in` expansion already emits there (`dump`'s `identity` column, the graph document's `identity` row field) — and every surface that names a RULE publishes the authored rule id."* Lint is named in C3's reader list three lines earlier and then omitted from the split.

**Symptom.** Author steps `attempt` over cells 0–4. Cell 3's row has a defect — say `graph-owned-before-write` on a co-written tag. Lint emits one finding reading `"rule": "retry"`. The author opens `retry`, which is four lines of TOML and obviously correct, because the defect is in the *generated* row for cell 3. There is no cell in the finding. Before this record, the author had `retry-3` in the message and went straight there. The record's headline claim — better review legibility — is precisely inverted on the tool the author uses to review.

---

## 2. The one section rewritten within 6 weeks of shipping

**`0030:C3` (Invariance).**

C3 is the section that claims nothing changes. It is written as a promise about surfaces, and it is the promise that first contact breaks, for three compounding reasons already grounded above.

First, C3's surface split is incomplete. It enumerates `flow next`, `flow resolve`, `flow set-state`, `lint`, `dump`, `graph` as readers, then draws a ROW/RULE split covering only `dump`, `graph`, `flow next`, `flow resolve`. `lint` is dropped between the two sentences. Every `graphlint` finding names `row.RuleID`. Within the first week that an author debugs a real stepped ladder, someone files the bug: "lint says rule `retry` four times and I cannot tell which cell." The fix is to add the row identity to findings — which is a new field on `clierr.Finding` or a repurposing of `Element`, either of which is an **envelope member**, which C3 explicitly forbids ("No emitted vocabulary … gains a member, so no `0029:C4` stability tier is owed by this record"). The fix contradicts the clause.

Second, C3 pushes two genuine hazards into `docs/model-authoring.md`: that `domain` order is now semantic, and that `unless` never excludes a step cell. Documentation-as-mitigation for a silent semantic change survives until the first person reorders a domain. Then the request is for a lint finding, which is a new finding code, which C3 also forbids.

Third, `0030:§decision-rationale` pre-concedes the retreat: *"Per-row identity IN a flow payload is a possible successor record — disclosed here as a non-obligation, not an owed fix."* A record that names its own successor in the rationale has already told you which clause is provisional.

The Finalization Gate is unfilled (C-12), and its Proportionality instruction is the one that would have caught this: the record authors a grammar, an expansion mechanism, a conformance check, and a cross-package evaluator move. C3's job is to hold all four still. It cannot.

---

## 3. The one assumption that will not survive first contact with a real user

**That the author will author the cap exclusion.**

This is stated across `0030:§consequences`, `0030:C2` and `0030:D-selection-predicate`, most sharply in C2: *"The remedy is the author's, and it is ONE fix with a consequence, not two alternatives: exclude the cell with a positive atom on the stepped tag."*

The record's own worked example concedes the shape of the problem. From `0030:§illustrative-code`: `tier = { step = 1 }` with no atom admits `large` and refuses at load. The author's mental model — the one the Problem Statement says they arrived with — is *"step the tier."* Not *"step the tier, and separately, in a different block, restate which tiers are steppable."* The cap is the thing the loop-free formulation was supposed to make implicit.

It gets worse on `enum`, and `0030:§consequences` says so out loud: ordered operators are `int`-only, `unless` is not consulted, so the **only** atom that can exclude the terminal member is `in = [<every member but the last>]`. The record's stated grievance was "every change to the tier domain means re-typing." The remedy is a list that must be re-typed on every change to the tier domain. The record calls this "reduced but not eliminated" and points at a successor record to close it.

And the natural thing the author reaches for is the thing that silently does not work. `unless tier = "large"` is the intuitive spelling, it parses, it loads on an unstepped rule, and against a step it is simply not consulted when cells are admitted (`0030:C1`). C2 patches the symptom by making the refusal detail name the unconsulted `unless` atom — which is an admission that the record expects this exact mistake to be the common one, and chooses to explain it in an error string rather than to make the intuitive form work.

First contact: an author writes a three-tier escalation, uses `unless`, gets a load refusal quoting a cell they did not think about and an atom they thought was doing the work, rewrites it as an `in` list of every tier but the last, ships it, adds a fourth tier six weeks later, and the model refuses to load in a rule they did not touch.

---

## 4. Premortem

*Written at draft time, as if the failure has already happened.*

We shipped `{ step = n }` in three phases across five weeks. Phase 2 took three of them; the Evaluator relocation (`internal/guard/grammar.go` to `internal/resolve`, plus `intWidth`, plus repointing `guard/product.go::valueSatisfies`, `graphlint/reach.go::atomAdmitsValue` and `cli/flow_resolve.go::guardSeam`, plus ~25 test references across seven `internal/guard/*_test.go` files including the reflection assertion in `guard_evaluator_0003_test.go`) landed first and clean, exactly as `0030:A12` predicted. That was the part we were worried about. It was not the part that failed.

The first real failure was in `internal/table/normalize.go`. `expand` is declaration-free by construction — `expand(base Row, predicates []Atom, outcome Atom, writes []TagValue)` — and C1's admitted-cell filter needs `TagDecl`. We threaded `*loader` into `expand`, which meant `expand` is now a method, which meant the fixture-level tests that call `expand` directly all needed a model. Then the placeholder problem: `renderWrites` calls `conform(decl, "eq", members)` on the authored value before `expand` ever runs, and there is no placeholder string that both conforms and is unmistakable. We special-cased it — `if isStep { skip conform }` — and in doing so dropped the kind check on the *declaration* for stepped keys. That shipped. A model stepping an `enum` whose `domain` was empty loaded, expanded to zero cells, and hit the zero-cell refusal, which reported the wrong defect.

The second failure was `Identity()`. `0030:D-identity` said step points order by key alone; `expand` sorts one candidate list with `compareAtoms`. We gave the step point a synthetic `Atom{Key: k, Block: BlockAll, Operator: "eq"}` so the existing comparator would take it, which meant its literal was empty, which meant two step points on keys ordered fine but a step point sorted against an `in` on a *different* key ordered by the full tuple rather than by key. It was stable in our test matrix. It was not stable against the customer's model, where a step on `attempt` and a match `in` on `mode` produced `retry#0#fast` in CI and `retry#fast#0` on the author's machine after a Go toolchain bump. Their graph export diff churned. `0030:C3` promised the graph document's `identity` row field was the stable one.

The third failure is the one that ended the feature's adoption. An author converted the motivating 24-row retry/escalation ladder to 8 rules. `intrastate lint` reported `graph-coverage-gap` — correctly, on a cell they had excluded per C2 — and then reported `graph-owned-before-write` on a co-written tag. Both findings named `"rule": "retry"`. There are five `retry` rows. `internal/graphlint/groups.go:80`, `:207`, `:318`, `:363`, `:399`, `analysis.go:213`, `:567` and `coverage.go:330` all populate `Rule` from `row.RuleID`, and `coverage.go:91`/`:140` use `firstRuleID(g)`; not one calls `Row.Identity()`. The author asked which attempt value was broken. We told them to run `intrastate graph --as json` and match on the row `identity` field. They asked why the unrolled table, which told them `retry-3` directly, was worse. We did not have an answer.

Our MVV had passed. `0030:S1` compares finding sets under the projection `(code, key, dimension, class, reason)` and drops `rule`, `span`, `element` and `fingerprint` — the exact four fields that carry row attribution. The oracle was built blind to the defect by design, and `0030:§consequences` says so in plain text: *"the finding comparison is over S1's projection, which drops exactly those fields."*

Then the `unless` case. `0030:A5` and `0030:S5b` promised a dead row lints advisory `graph-redundant-row`. `checkRedundantRows` at `internal/graphlint/groups.go:299` requires `mine.Projectable()`, a same-`Kind()` sibling in the same group, and `properSubset` — `left.Subset(right) && left.Len() < right.Len()`. An `unless`-carrying row is the canonical unprojectable row; `decidableAccepted`'s comment names it. S5b's fixture was engineered to project, and asserted its own precondition to get there. The customer's ladder did not project. The dead row claimed a cell, matched nothing, and lint said nothing. `0030:BR1` had rejected saturation on the grounds that `0002:C3` forbids "a silent no-op." We shipped one by a different route.

Finally, `intrastate lint` stopped returning on a model that loaded in twelve milliseconds. The author had written `attempt = { step = 1 }` over `min = 0, max = 200` — a reasonable retry budget for a batch job. Load expanded to 201 rows in under a second, exactly as `0030:A7`'s spike predicted. `0030:§performance-expectations` had already published the other curve: 116 s at 200 rows, beyond 600 s at 1,000. We had measured the blowup, written it down, argued that it was `0013`'s analysis scope rather than ours, and concluded no load-time ceiling was owed. `0030:A7` proved the load path cannot reach the lint limit — `go list -deps ./internal/table` names neither `guard` nor `graphlint` — which we cited as the reason not to gate at load. It was the reason we should have.

Journeys that broke, in order of how often they were hit: *author converts an existing unrolled ladder and cannot verify the conversion* (no diff aid exists outside the MVV — `0030:§consequences` declines to extend `0002:§round-trip-inverse-invariants` to the pair); *author steps an enum tier and writes `unless tier = "large"`*; *author adds a fourth tier and a rule they did not touch refuses at load*; *author alphabetizes a `domain` array and the flow silently changes behaviour, because `0030:A2` made authored order semantic and `0030:C3` put the warning in a markdown file*.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

**AT-1 — the loader can actually carry a step spec (C-1, C-2)**

```gherkin
Given the RDR names the functions that carry the step spec
When a reviewer opens internal/table/normalize.go
Then expand's signature must already admit a tag declaration, or the RDR
     must name the declaration-threading change as a phase item
And renderWrites' conform(decl, "eq", members) call must have a stated
     disposition for a stepped key: either the placeholder conforms and the
     RDR names the literal used, or the check is relocated and the RDR names
     where the authored-value kind check survives
```
Failing this at review forces C1 to stop calling the change "one signature widening."

**AT-2 — suffix element order is a written comparator (C-3)**

```gherkin
Given a rule with a step point on "attempt" and a match `in` on "mode"
When the RDR is read for the ordering rule
Then it must state a single total comparator over the mixed candidate list,
     not two rules keyed on whether a candidate carries an Atom
And a fixture must assert Identity() is byte-identical across two runs with
     a shuffled input atom order
```

**AT-3 — a lint finding names the cell (C-4, C-5, C-19)**

```gherkin
Given a stepped rule "retry" expanding to cells 0..4
And one expanded row carries a defect lint reports
When intrastate lint --as json runs
Then exactly one finding is emitted for that cell
And the finding identifies WHICH cell — by rule#cell, or by a stated field
When the MVV's finding-set projection is read
Then it must not drop every field capable of carrying that identification
```
This is the test the RDR most needs and most nearly forbids: `0030:S1`'s projection excludes `rule`, `span`, `element` and `fingerprint` by construction.

**AT-4 — the dead row is actually loud (C-6, C-7)**

```gherkin
Given a stepped rule with guard.unless on the stepped tag excluding an
      interior cell
When the fixture is authored the way a real author would write it —
     not engineered to satisfy AcceptedAssignments().Projectable()
Then lint emits graph-redundant-row for the dead row
And if it does not, the RDR must say so, because a cell claimed by a row
     that matches nothing is the silent no-op 0002:C3 forbids
```
Run against `internal/graphlint/groups.go::checkRedundantRows` at review time this fails: `Projectable()` gate, same-`Kind()` gate, and strict `properSubset`.

**AT-5 — the enum arm solves the stated problem (C-8)**

```gherkin
Given the Problem Statement's grievance "every change to the tier domain
      means re-typing"
When a fourth tier is added to a stepped enum's domain
Then the number of author edits required must be zero
```
This fails on the record as written — the `in` exclusion list must be re-typed — and `0030:§consequences` already concedes it. The correct review outcome is to split the enum arm out or to widen the guard grammar first.

**AT-6 — the bound refusal is diagnosable and machine-readable (C-9, C-18)**

```gherkin
Given a model that loads today
When `max` is tightened or an enum `domain` gains a member
Then any resulting load refusal must name both the declaration that changed
     and the rule that broke
And the refusal's parts (rule, tag, cell, result) must be assertable without
     substring position, or the RDR must accept that the message is frozen
```

**AT-7 — the evaluator move is sequenced against RDR 0012 (C-10)**

```gherkin
Given RDR 0012 is Draft and lands on internal/guard/grammar.go::Evaluator
When 0030 relocates that file to internal/resolve in Phase 2
Then the RDR must state which record lands first and what the other absorbs
And "0012 absorbs it" must be an agreement recorded in 0012, not a claim
     made in 0030
```

**AT-8 — no Pending assumption carries settled normative prose (C-11, C-12)**

```gherkin
Given A5, A6, A10, A11, A13 and A14 are Pending
And A10's "If wrong" states that both branches reopen C1's expansion half
When the Finalization Gate's Assumption Verification is run
Then lock is refused: the gate forbids Pending assumptions with settled-fact
     prose depending on them, and C1 states subsumption as settled
```
The gate is unfilled template text. Running it is the test.

**AT-9 — a plausible model is lintable (C-14, C-15)**

```gherkin
Given attempt declared int min = 0, max = 200
And a rule writing attempt = { step = 1 }
When intrastate lint runs
Then it completes within an interactive budget
And if it does not, the RDR must own a load-time ceiling rather than
     attributing the cost to another record's analysis scope
```
The RDR's own numbers — 116 s at 200 rows, >600 s at 1,000 — fail this before any code is written.

**AT-10 — the migration journey is supported (C-16)**

```gherkin
Given an author converting an unrolled ladder to a stepped one
When they want to verify the conversion preserved behaviour
Then a supported command must show the equivalence
And "the MVV compares them" is not an answer, because the MVV is not a
     product surface
```

**AT-11 — domain reordering is not silent (C-20)**

```gherkin
Given an enum domain a rule steps
When the domain's members are reordered with no other change
Then the change must be observable: a lint finding, a graph diff, or a
     refusal
And documenting the hazard in docs/model-authoring.md does not satisfy this
```
