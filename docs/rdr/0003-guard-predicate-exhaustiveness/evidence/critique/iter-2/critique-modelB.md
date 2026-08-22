Model: claude-sonnet-5

# Critique — RDR 0003 Guard Predicate Exhaustiveness (iteration 2 re-entry)

## Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | §Technical Design, "an author opts out by leaving a dimension undeclared" (default-on exhaustiveness reading) | An author who forgets to declare a bound gets a *silent* strengthening of the exhaustiveness claim on unrelated dimensions in the same row group, not a warning at the point of omission — the failure surfaces later as a blocking lint finding on a group the author never intended to submit for proof | Flow author adds one new row to an existing, previously-green row group; a passing table now fails lint with `graph-unprovable-coverage` on a dimension they never touched, with no diff-local explanation of why this row's addition triggered scrutiny of an unrelated tag | §2/premortem |
| C-2 | A7, "Pending... this RDR is NOT lockable until it resolves" + Prerequisites: RDR 0002 must accept the producer request first | The presence-dimension projection for `exists` is not just unresolved research, it is unresolved *design* — the derivation lives in a side file (`evidence/research/iter-2-projection-derivation.md`) that Implementation Plan never requires re-reading before Phase 1 begins, and Phase 1 ("Define the tag-kind/operator compatibility matrix") does not gate on A7 explicitly | `exists`-bearing row groups (RDR 0007's own sanctioned two-row absence pattern) implement first, get exhaustiveness-certified incorrectly during the gap between Phase 1 landing and A7's normative clause landing, or block indefinitely because RDR 0002 (also Draft) has not accepted the optionality-field request | §1/premortem |
| C-3 | A8, "Pending and also blocks lock" — cross-document agreement with `Final` RDR 0006 that has no route-back mechanism owned by this RDR | This RDR cannot lock without a document it cannot edit (RDR 0006, `Final`) agreeing to carry language RDR 0006 does not currently carry, and the venue named ("cluster reconcile") is not a step in *this* RDR's own Implementation Plan — it is an external process this document merely gestures at | Implementation begins against this RDR's normative text; RDR 0006's shipped lint never emits the narrowed refusal because RDR 0006's Final text still reads "non-finite dimension only"; a row group with a possibly-absent guard key gets certified green in violation of this RDR's own §"exhaustiveness claim MUST NOT be stronger than the runtime" clause | §1/premortem |
| C-4 | §Technical Design / `disposition` table row "Row-group claim semantics... default-on for every scoped row group whose participating dimensions are all finitely declared" | The RDR itself flags this as a possible divergence from RDR 0006 ("If RDR 0006 intends an explicit per-group annotation instead, that is a divergence to settle at cluster reconcile") but proceeds to *implement against* the unsettled reading rather than blocking on it the way A7/A8 block | Two independently-correct implementations (one following 0003's default-on reading, one following a hypothetical 0006 opt-in reading) diverge in which groups get proven at all; a flow that passed lint under one build silently stops being checked — or starts failing — under a rebuild against the other reading | §2/premortem |
| C-5 | §Technical Design, Load-Bearing Decisions "Identity" — atom identity is `(RuleID, SourceLocator, key, block, operator, literal)` | "Semantic equality is `(tag, operator, literal)`" is stated as a separate, narrower equality from atom identity, but no normative clause or MVV scenario says which operations use which equality (e.g., does overlap detection dedupe on full identity or semantic equality? Does the "same atom under reordering" claim in MVV Scenario 5 test identity or semantic equality?) | Two authors write byte-identical guards in two different rule ids; lint either treats them as separate (correct, since RuleID differs) or an implementer conflates identity with semantic equality during overlap computation and silently drops one — no test in the MVV distinguishes these | §5/AT-3 |
| C-6 | §Technical Design, "if the finite product is too large for the implementation to prove deterministically, lint must also refuse or downgrade" | No normative clause, MVV scenario, or Performance Expectations paragraph names *how large is too large* — there is no stated bound, algorithm, or even an order-of-magnitude heuristic, so "too large" is undecidable by an implementer without inventing a threshold the RDR never sanctions | Two flows of similar size get inconsistent lint behavior depending on which implementer's undocumented threshold shipped; a CI run is green on one machine/build and blocking on another because the "too large" boundary is not part of the spec | §2/premortem |
| C-7 | §Metadata Status line + repeated "RDR 0006 is Final... does not carry this narrowing" | The RDR asks implementers to build against a normative promise ("lint must not certify exhaustive when a participating row can refuse") that the document responsible for actually emitting that lint finding (RDR 0006) does not yet implement or agree to; this RDR's own Prerequisites checkbox for A8 is unchecked, yet Phase 1-4 of the Implementation Plan do not reference A8 as a blocking gate on any specific phase | A team starts Phase 1/2 work under time pressure while A7/A8 are still open (both routed to slow external processes — RDR 0002 Draft acceptance, cluster reconcile), producing a working evaluator that cannot yet be lint-certified against its own stated promise, so "guard evaluation" ships ahead of "guard evaluation lint can trust" | §1/§2/premortem |
| C-8 | Illustrative Code example + operator matrix: no example shows a `contains` predicate | A9 (Pending) means the entire `contains` operator — one of five in the closed vocabulary — has zero worked example anywhere in a 1203-line document; the only fixture evidence (Resolve spike) explicitly excludes it ("the spike fixture declares no set-valued tag") | First author who needs set-containment (the RDR's own prior-art discussion assumes set-membership guards exist for lens sets) writes a `contains` guard, and either the operator is rejected at parse time (blocking basic flow authoring) or accepted but never exhaustiveness-provable, with no documented example of what "declared element universe" syntax looks like | §3/premortem |
| C-9 | §Failure Modes, "A guard that cannot be decided at runtime surfaces as RDR 0007's `guard_unevaluable` refusal, not as a predicate parse kind owned here" | The user-facing distinction between "this predicate is malformed" (this RDR's problem, rejected before resolution) and "this predicate is fine but the key might be absent" (RDR 0007's problem, refused at runtime) is a distinction only a reader of three RDRs can reconstruct; the CLI output contract (docs/cli-output-contract.md) says nothing about how these two failure classes render differently to a flow author debugging a broken transition | A flow author sees `error: guard_unevaluable` in text mode with no indication of *which* declared-tag optionality gap caused it versus a genuine authoring mistake, because RDR 0005 (not yet built) is the only place these get unified into one envelope, and this RDR explicitly disclaims owning that mapping | §3/premortem |
| C-10 | Implementation Plan, Prerequisites: three of five items are unchecked/blocked on other Draft or process-external documents, yet MVV and Phase 1-4 are written as if implementation proceeds linearly | The plan reads as sequential (Phase 1 → 2 → 3 → 4) but the actual dependency graph is: Phase 1 can start now, Phase 3's full operator acceptance is blocked on A9, the exhaustiveness promise Phase 2 implements is blocked on A7 and A8, and RDR 0007's kernel reshape (unchecked item, "gates implementation sequencing, not lock") is a hidden precondition for Phase 1 itself since the shipped kernel still has `Row.Guard string` | Implementers pick up "Phase 1: Predicate Model" per the plan's document order, build against the current `Row.Guard string`/`Evaluate(guard string, ...)` shape because that's what compiles today, and produce work that must be redone once RDR 0007's atom-slice reshape actually lands — a rewrite the RDR itself predicts ("Phase 1 cannot begin against the atom-slice shape until that reshape lands") but does not sequence as a blocking precondition in the Phase list itself | §2/premortem |

## 1. The three most likely ways implementation goes wrong

### Failure 1: The exhaustiveness promise ships without the document that has to honor it

**Root cause.** This RDR's central selling point — "a green exhaustiveness result must mean resolution succeeds" — is a promise this RDR *states* but cannot *fulfill alone*. The mechanism that would make the promise real is a blocking lint finding that RDR 0006 must emit whenever a domain-exhaustive row group contains a row that can still refuse `guard_unevaluable` at runtime. RDR 0006 is `Final`. Its exhaustiveness clause (verified at `docs/rdr/0006-graph-lint-authority-and-guarantees.md:344-350`, read verbatim) is scoped to exactly one case: "If a required dimension is not finite, lint MUST emit a blocking inability-to-prove finding." It says nothing about a fully-finite product whose participating row can still refuse at runtime. This RDR's own Assumption Verification concedes this directly: "`Final` RDR 0006 carries no `guard_unevaluable` narrowing."

**Specific passage.** A8: "**A8 is Pending and also blocks lock.**... Closing it needs a route-back to RDR 0006 or a §JD-4 disposition, not an edit here." And the Normative Contracts block: "An exhaustiveness claim MUST NOT be stronger than the runtime it describes: lint MUST NOT certify a row group exhaustive when a participating row can refuse `guard_unevaluable`..." — a MUST this RDR issues to a lint authority it does not own and that has already shipped `Final` without it.

**Symptom the user will see.** A flow author writes a row group that is domain-exhaustive on paper (every enum/bool/int dimension is fully declared and partitioned) but contains a row whose guard reads a tag that is not actually present on every reachable predecessor path. Lint — built against the *shipped* RDR 0006 contract, because that is the `Final` document an implementer actually codes against — reports success. The author trusts the green result. At runtime, that exact row refuses `guard_unevaluable` on the input the "exhaustive" proof claimed was covered. This is precisely the failure this RDR exists to prevent, and its own design leaves the door open because the enforcing document was locked before the rule existed. The dependency is not abstract: A8 is explicitly a **lock blocker** for this RDR, but nothing in this RDR's own Implementation Plan phases treats it as an implementation blocker with a concrete phase gate — it is a checkbox in Prerequisites, disconnected from Phase 1-4.

### Failure 2: `exists` — the operator RDR 0007 calls "total" — ships unprovable, and nobody notices until the sanctioned pattern breaks

**Root cause.** A7 is Pending. The presence-dimension projection needed to make an `exists`-bearing row group provable at all does not exist as a normative clause in this document — it exists as a *derivation* in a side evidence file (`evidence/research/iter-2-projection-derivation.md`) that the Implementation Plan does not require anyone to promote into the normative section before Phase 1 starts. Worse, this RDR's own trace table (`trace` step 5) states the gap is not hypothetical: it disqualifies the RDR's *own representative fixture row* (`foundational-to-cove`) and "leaves RDR 0007's sanctioned two-row absence pattern permanently unprovable" if A7 never closes. A7 is itself blocked on a producer request to RDR 0002 (a per-tag optionality declaration) that RDR 0002 — itself `Draft`, not `Final` — does not currently carry (verified: `docs/rdr/0002-transition-table-as-reviewable-data.md:217-218` lists tag declaration fields as name, provenance, value kind, optional accessor reference — no optionality field).

**Specific passage.** A7's own "If wrong" clause: lint "either drops `exists` atoms from the product — certifying a row group exhaustive that refuses `guard_unevaluable` at runtime... or refuses every group containing an `exists` atom, which disqualifies this RDR's own representative fixture row `foundational-to-cove`... pushing authors back to the sentinel-stamping anti-pattern." Both arms of that disjunction are bad; the RDR names them but resolves neither, because A7 is not closed.

**Symptom the user will see.** Since RDR 0007 fixes `exists` as "the sole TOTAL operator" and explicitly sanctions the two-row absence idiom as the correct way to encode "this tag may or may not be set," any implementer who follows RDR 0007's guidance today will write guard rows using `exists`. Under this RDR's *unclosed* A7, that author's row group cannot be exhaustiveness-proven — not because their guards are wrong, but because the projection rule that would make them provable hasn't been written yet. The practical result: the RDR's own sanctioned idiom for expressing optional-tag guards produces flows that cannot pass the lint this RDR promises, and the fallback the RDR itself names is regression to a "sentinel-stamping anti-pattern" — i.e., authors work around the feature by not using it as designed.

### Failure 3: The closed operator vocabulary ships with one operator (`contains`) that is unauthored, untested, and unbound

**Root cause.** A9 is Pending: "a declared set-element universe for `contains`" does not exist in RDR 0002's tag declaration today, and "the Resolve spike fixture declares no set-valued tag." The Phase 3 acceptance gate explicitly requires "at least one `contains` predicate over a declared set-valued tag" before "the full operator vocabulary is accepted," which means the vocabulary marketed throughout the document ("equality, membership, bounded integer comparison, existence, and set containment" — five operators, stated as closed and final in the Normative Contracts) is, as of lock, four operators plus one aspirational one with no worked syntax anywhere in 1203 lines.

**Specific passage.** The operator/kind matrix table lists `contains` with "Literal shape: non-empty typed element set" and "Lint proof role: Narrows the set-valued domain to assignments containing every listed element" as if this is settled behavior, immediately followed by A9's own admission: "the operator vocabulary is closed *as specified* but exercised only over the target-flow subset the spike proved" — i.e., `contains` is specified but never exercised.

**Symptom the user will see.** A flow author whose guard genuinely needs set-containment (the RDR's own domain examples repeatedly invoke "lens sets" and similar collection-shaped tags) reaches for `contains`, finds it in the operator table, writes a guard using it — and either (a) the implementation rejects it because RDR 0002 has no element-universe declaration syntax to parse, producing a confusing "operator not supported" error for a documented operator, or (b) an implementer ad-hoc invents an element-universe declaration syntax to unblock themselves, and that invented syntax becomes de facto load-bearing before RDR 0002 (Draft) ever ratifies it — exactly the kind of undocumented, unreviewed extension this RDR's entire premise (symbolic, reviewable, lint-provable guards) exists to prevent.

## 2. The section that will be rewritten within 6 weeks of shipping

**The "default-on for every scoped row group whose participating dimensions are all finitely declared" reading in Technical Design**, paired with the entire A7/A8/A9 apparatus that surrounds it, will be rewritten first and fastest.

The RDR itself names the risk and does not resolve it: "RDR 0006 gates its blocking finding on a contract 'that claims closed coverage', and no document defines how a group makes that claim: this RDR reads it as **default-on**... If RDR 0006 intends an explicit per-group annotation instead, that is a divergence to settle at cluster reconcile." This is not a hedge about a minor detail — it is a hedge about whether exhaustiveness checking is opt-out (this RDR's reading, silently strong) or opt-in (a plausible alternative reading of a `Final` sibling this RDR does not control). An opt-out design means every incomplete or exploratory row group an author writes during flow development is *by default* subject to a blocking lint check the author did not ask for and may not be ready to satisfy — the RDR's own Risk register even acknowledges "Finite-domain declarations feel like boilerplate" as a known risk, but does not connect that risk to the churn a wrong default-on/opt-in choice will cause once real authors hit it in week one of use.

Combine that with A7 (the presence-dimension projection) and A8 (cross-document narrowing agreement) both being *lock blockers this RDR cannot close by itself* — routed to a `Draft` peer (RDR 0002) and a `Final` peer via an unowned process (cluster reconcile) respectively — and the picture is a document whose central "the promise narrows correctly" guarantee is provisional on decisions made by other documents after this one locks. The moment RDR 0002 adds the optionality field, or §JD-4 assigns the narrowing to RDR 0006 instead of here, or cluster reconcile picks the opt-in reading over default-on, the Technical Design section's prose (not just a cross-reference) has to change, because the *shape* of the normative clauses ("default-on... an author opts out by leaving a dimension undeclared") is itself the contested reading, not a detail downstream of it.

## 3. The assumption that will not survive first contact with a real user

**A2's "or the finite product cannot be represented by the implementation's symbolic/bitset-equivalent proof, lint must refuse or downgrade" — and more precisely, the unstated threshold for "too large."**

Every other Critical Assumption in this document has a stated verification method with concrete evidence: a spike transcript, a source-code symbol, a peer-RDR citation. A2 is `Verified` by `Method: Derivation` — a closed-form mathematical description of coverage and overlap as set operations over a scoped product. That derivation is sound as math. But the moment a real flow author writes a row group with, say, three declared enum tags of 6, 8, and 12 values each plus a bounded integer range of 50, the "scoped product" is 6×8×12×50 = 28,800 elements, and the RDR has no stated bound, no stated algorithm complexity class, and no stated fallback behavior other than "refuse or downgrade... rather than silently capping the proof." The Performance Expectations section confirms this is genuinely unmeasured: "Resolve measured representative fixture size rather than setting a throughput target: the spike uses four rows."

A first real user is not going to hit the four-row fixture that was spiked. They are going to hit whatever combination of enum/bool/set/bounded-int tags their actual flow needs, and the first time the scoped product crosses whatever undocumented threshold an implementer chose, that user's lint call either (a) hangs or is slow enough to feel broken, or (b) silently downgrades to a blocking "cannot prove" finding for a group the author reasonably expected to be provable, with no way to know in advance whether their flow design is within budget. Nothing in the MVV or Testing Strategy scenario 3 pins a concrete size — it says "lint an... declared finite product too large for the deterministic proof representation" as a scenario to test, but the RDR never says what "too large" means in a way an author could design against. This is the assumption that looks airtight on paper (pure math, `Derivation`, verified) and is actually a UX and performance cliff the document has not measured.

## 4. Premortem

*Six weeks after RDR 0003 ships and RDR 0007's kernel reshape lands, guard predicates are live in production flows. Here is what happened.*

The team implemented Phase 1 and Phase 2 against the shipped `internal/resolve/resolve.go::GuardEvaluator.Evaluate(guard string, view TagSet) GuardResult` interface first, because that's what compiled — RDR 0007's atom-slice reshape hadn't landed yet, and the Implementation Plan's unchecked Prerequisite item ("Phase 1 cannot begin against the atom-slice shape until that reshape lands... gates implementation sequencing, not lock") was read as a soft caveat, not a hard blocker, because nothing enforced it as a phase gate. Two engineers built a string-based interim guard grammar to unblock Phase 2 work, planning to migrate once the atom slice landed. That interim grammar shipped to three internal flows before the migration happened, because "it works" beat "it's the spec'd shape" under a deadline.

Meanwhile, flow authors writing the two target flows — RDR-flow status/profile routing and kata cluster eligibility — used `exists` guards exactly as RDR 0007 sanctions, encoding "this tag may be present or absent" with the two-row absence pattern. Lint, built against A7's still-open projection (the team shipped Phase 2 with a stopgap that simply excluded any row group containing an `exists` atom from the exhaustiveness product, the "refuse every group" arm A7's own "If wrong" clause predicted), silently downgraded every such group to "cannot prove" — including `foundational-to-cove`, the RDR's own representative fixture row. Authors saw the downgrade finding so often on ordinary flows that they started treating "cannot prove — inability to prove coverage" as background noise rather than a signal, exactly the alert-fatigue failure mode the RDR's Failure Modes section warns against ("Silent failure would be a false exhaustiveness claim") without naming its cousin: a *ubiquitous* non-blocking-feeling "cannot prove" is functionally the same as silence, because nobody reads noise.

The real incident: a kata cluster-eligibility row group was, in fact, domain-exhaustive and contained one row with a value atom over a key that was absent on one reachable predecessor path — precisely the scenario this RDR's own MVV Scenario 6 was written to catch ("lint emits the blocking inability-to-prove finding naming that row and the refusing atom"). But the team's shipped RDR 0006 lint — the `Final` document, unmodified, because A8's route-back never happened; cluster reconcile got deprioritized behind feature work — only implements the non-finite-dimension refusal clause. The row group in question had *all finite* dimensions. RDR 0006 certified it exhaustive. The flow shipped. A kata run hit the uncovered predecessor path, `Resolve` refused with `guard_unevaluable` at runtime, and the on-call engineer spent two hours confused because the CI lint gate — the entire reason this RDR exists — had said green.

The postmortem traced it to exactly the gap A8 named at RDR-review time: "an implementer reading only RDR 0006 certifies a row group exhaustive that this RDR forbids." The fix was not a code bug. It was that RDR 0003's promise was never actually binding on the tool that emits the lint verdict, because RDR 0003 locked (or was implemented against) before RDR 0006 agreed to carry the narrowing, and nothing in the Implementation Plan made that agreement a hard gate on writing the `graph-unprovable-coverage` finding emitter.

## 5. Acceptance tests that would have caught each failure at RDR-review time

**AT-1 (catches Failure 1 / C-3, C-7).**
```
Given RDR 0003's Normative Contracts state a MUST binding RDR 0006's lint output
And RDR 0006 is status Final and does not contain that MUST's language
When the Finalization Gate's Contradiction Check runs
Then the gate MUST fail (not merely note an A8 assumption as Pending)
And RDR 0003 MUST NOT be eligible to move to Final
 until either RDR 0006 is amended to carry the narrowing
 or the Implementation Plan names a specific phase-gate task
 that blocks "graph-unprovable-coverage" emission on the narrowing
 being verified present in RDR 0006's shipped code.
```
This would have forced A8 out of "Pending, blocks lock" limbo and into either an actual route-back completed before merge, or an explicit, code-verifiable phase gate — not a Prerequisites checkbox disconnected from the Phase 1-4 sequence.

**AT-2 (catches Failure 2 / C-2).**
```
Given A7's presence-dimension projection is Pending
And the RDR's own trace table shows the RDR's representative fixture row
 (foundational-to-cove) is disqualified without it
When Phase 1 ("Define the tag-kind/operator compatibility matrix") is scheduled
Then Phase 1 MUST explicitly declare whether `exists` atoms
 are in scope for Phase 1's matrix
And if excluded, the Implementation Plan MUST state what happens to
 `exists`-bearing row groups in the interim
 (blocked entirely? downgraded silently? downgraded loudly with a
 distinct diagnostic from "not yet implemented"?)
 rather than leaving the interim behavior to whichever implementer
 hits it first.
```

**AT-3 (catches Failure 3 / C-8).**
```
Given the operator vocabulary table lists `contains` as specified
And A9 states no worked example or fixture exists for `contains`
When Phase 3's acceptance gate is defined
Then the RDR MUST include at least an illustrative TOML fragment
 showing the element-universe declaration syntax `contains` depends on,
 even if RDR 0002 has not yet ratified that field,
 so implementers are not left to invent the syntax under deadline
And the closed-vocabulary claim in Normative Contracts MUST be scoped
 to "four operators verified, one operator specified-pending-fixture"
 rather than stated as one undifferentiated closed set.
```

**AT-4 (catches C-1, C-4 — the default-on ambiguity).**
```
Given this RDR's own Technical Design flags a possible divergence
 with RDR 0006 over opt-in vs default-on exhaustiveness scope
When the Finalization Gate's Contradiction Check runs
Then the gate MUST treat this as an open contradiction with a Final
 sibling document, not a resolved reading
And lock MUST be blocked until RDR 0006 confirms (in writing, cited by
 symbol/section) which reading it implements,
 the same bar A8 is already held to for the narrowing question.
```

**AT-5 (catches C-6 — undocumented "too large" threshold).**
```
Given Performance Expectations states no throughput or size target
 beyond "the spike uses four rows"
When the MVV's Scenario 3 ("finite product too large to prove
 deterministically") is written
Then the RDR MUST state a concrete scoped-product-size threshold
 (or an explicit formula/complexity bound) above which refusal is
 mandatory, sourced from a spike that exercises a product at least
 two orders of magnitude larger than the four-row fixture
And the MVV test MUST assert on that stated threshold, not merely on
 the existence of "a" too-large case chosen by the implementer.
```

**AT-6 (catches C-9 — undifferentiated user-facing errors).**
```
Given a flow author can hit either "predicate malformed" (this RDR)
 or "predicate fine, key possibly absent" (RDR 0007's runtime veto)
 or "cannot prove, dimension not finite" (RDR 0006's lint)
When a flow author runs `--as=text` against a broken guard
Then the RDR (or a cited peer RDR) MUST specify a rendered example
 showing how these three distinct failure classes are visually and
 textually distinguishable to the author,
 not merely internally distinguishable by RefusalKind/finding code,
 because docs/cli-output-contract.md's text-mode format
 ("error: <code>: <message>") gives the author one line and this
 RDR must confirm <code> alone is enough to self-diagnose which of
 three RDRs' failure surface they are looking at.
```
