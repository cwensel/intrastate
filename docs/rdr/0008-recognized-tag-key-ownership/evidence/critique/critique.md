Model: claude-opus-5[1m]
Pass: A (single-model fallback — no alt model reachable; resources.md records "Alt-model roster: omitted (single model)")

# Hostile critique — RDR 0008 (recognized-tag-key ownership)

## Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | A4 Evidence: "The three committed TOML fixtures use `recognized` only as a provenance value under an author-named key … The reserved-key rule invalidates nothing at HEAD, so the mechanical-rename fallback below is not triggered." | A4 asks the wrong question. It sweeps for `recognized` in *owned/observed* position and concludes nothing breaks. The rule that actually breaks things is block 2's *first* clause — a recognized-provenance declaration MUST be named `recognized` — and every committed fixture in the repo violates it: `0002/evidence/spikes/rdr-fixture.toml:25` `[tags.outcome]`, `kata-fixture.toml:22` `[tags.outcome]`, `0003/evidence/spikes/guard-fixture.toml:36-37` `[tags.rewind_target]`. RDR 0002 Load-Bearing Decisions (`0002…md:361`) names these "the canonical examples implementation tests must promote." | The 0002 implementer promotes the canonical fixture, and the very first lint run rejects it in the `reserved_tag_key` category. The reserved-key rule's debut is breaking the spec's own worked example. | §1, §3, premortem, AT-1 |
| C-2 | A4 "If wrong": "implementation must rename the colliding fixture/table tags before the lint lands; scope grows by a mechanical rename, not a design change." | The rename is not mechanical. `[tags.outcome]` is referenced from `[rule.match.outcome]` at `rdr-fixture.toml:59,74,86` and `kata-fixture.toml:47,57`. Renaming the declaration forces renaming every match-position reference across two RDRs' canonical fixtures and 0003's guard fixture — and the RDR never states which side of `[rule.match.<tag>]` even normalizes into `Row.Outcome` versus `Row.Match`, so the implementer cannot tell whether the rename changes gating semantics. | Implementer stalls: the rename touches peer-RDR evidence artifacts that 0002's Wire/byte-format decision declares canonical, with no ruling on whether `[rule.match.outcome]` was ever the recognized-outcome gate. | §1, premortem, AT-2 |
| C-3 | Normative block 2: "The checked positions are the `[tags.<tag>]` declaration keys only. Reserved-name uses in predicate or write positions … need no separate reserved-key check: a conforming model declares `recognized` and those references resolve to it, while a model that does not declare it already fails RDR 0002's `unknown tag` rule." | The fall-through is wrong for the case that matters. An author who declares `[tags.recognized]` (legal) *and* an owned `[tags.outcome]` (legal), then writes `[rule.match.outcome]`, has a fully conforming model that never matches the recognized outcome — every name check passes. The clause proves only that *undeclared* references fail; it silently assumes reference position carries no ownership question. | The author's row still never fires. This is the exact Problem Statement symptom the RDR opens with, reproduced inside a model that passes every rule the RDR adds. | §1, §2, premortem, AT-3 |
| C-4 | Approach 2: "The author declares *that* they use the affordance, not what it is called." + Normative block 2's naming rule, with no lower-bound obligation. | The declaration is required to be *named* `recognized` but is never required to *exist*. A model whose rules gate on outcomes and declares no recognized-provenance tag at all passes every check in this RDR. Cove Q10 flagged the missing lower bound; the fix pass added an upper-bound cardinality paragraph and left the lower bound open. | A model that forgot the recognized declaration lints clean and refuses `no_match` at resolve time — the silent discovery the RDR was written to abolish, one declaration up. | §1, §2, premortem, AT-4 |
| C-5 | Normative block 4: "Producers of kernel `Input` MUST NOT supply an owned or observed tag keyed `recognized` … `Resolve` MUST apply that same predicate at entry, returning a non-nil error and no `Result` disposition on breach." | The obligation is stated on the wrong side of the accessor boundary. The owned snapshot is produced by RDR 0004's accessor layer from artifact contents at runtime; a tag key can be *data-derived*, not code-literal. The RDR's own evidence for "no breach exists" (A8) is a grep of `*_test.go` for the literal `Key: "recognized"` — a literal-source sweep that cannot see a key that arrives from a parsed artifact. | An artifact whose parsed content happens to yield an owned key `recognized` turns a working resolve into a hard Go error — no `Plan`, no `Refusal`, no diagnosis path, on the channel RDR 0001 reserves for *programmer* mistakes. The user's data becomes a programmer mistake. | §1, premortem, AT-5 |
| C-6 | Enforcement locus: "(b) is the only locus that sees *every* producer including non-TOML callers, and it lands on reserved-but-unexercised surface … so (b) adds a check without widening the signature." | "Adds no surface because the signature already returns `error`" conflates type surface with behavioral contract. Cove Q1 confirmed every `return` in `Resolve` ends `, nil` and the file imports neither `errors` nor `fmt`. Every existing caller and every existing test was written against a total function. Introducing the package's first non-nil error path is a behavioral break for callers that ignore the error — which is every caller, because none could ever be non-nil. | Callers written against a never-failing `Resolve` silently proceed on a zero-valued `Result` — `Refused()` returns false, `Plan` is nil — and dereference nil, or treat the breach as a successful non-refusal. | §1, §2, premortem, AT-6 |
| C-7 | Normative block 4: "The check is unconditional on `Input.Recognized`: a reserved-keyed owned or observed tag is a breach whether or not the resolve carries an outcome." | This makes the predicate strictly more aggressive than the harm it prevents. When `Input.Recognized == ""`, `assemble` (`resolve.go:151`) injects nothing, so there is no shadowing and no collision — the key is an ordinary owned tag. The RDR errors anyway, purely to match block 1's rhetorical "unconditional reservation." Scenario 6 pins this as *expected*. | A resolve that would have worked correctly — no recognized outcome in flight, an owned tag innocently keyed `recognized` — hard-errors. There is no possible mis-match to protect against in that input. | §1, premortem, AT-7 |
| C-8 | "Contract count (Resolve, evidence-grounded): **one** independent load-bearing contract … They are not separable seams" + Profile `foundational`. | Six normative blocks landing in four codebases (kernel constant, 0002 normalizer lint, kernel `Input` predicate + `Resolve` entry, normalizer producer obligation on `RequiresOwned`) are asserted to be one contract by declaring the enumeration complete. That enumeration was already extended once (blocks 1–3 → +4 → +5 after cove F-8), each time re-asserting completeness. The count is defended by rhetoric, not by a closure argument over the channels the key can arrive through. | The RDR ships, a sixth channel surfaces, and the "single contract" framing has already routed it past the split gate. Downstream, no one can revise the lint rule without touching the kernel predicate, because they were locked as one contract. | §2, premortem, AT-8 |
| C-9 | Testing Strategy: "**Carried into RDR 0002's implementation** as a written obligation this RDR authors and 0002 executes: scenarios 2, 4, 5, 7, and the lint half of the MVV." | The RDR's entire user-facing value — the typed load/lint failure that replaces the silent no-match — is in the half that does not execute. "Done for this RDR's own implementation" is scenarios 1(kernel half), 3, 6, 8, 9: a pointer comment, a view-plumbing test, and three tests of the new error path. Zero of them deliver the Problem Statement's outcome. | The RDR ships as Final-and-implemented having delivered the kernel's first hard-error path and nothing else. The author whose row never fires is exactly as stuck as before, and the RDR's status says the problem is solved. | §2, premortem, AT-9 |
| C-10 | Problem Statement: "**The outcome this RDR commits to is 'my row fires, or I am told why'**" + Consequences: "a first-time author's comprehension is therefore asserted at the data level only." | The committed outcome is user comprehension; the tested contract is three machine-readable fields. The RDR explicitly declines to specify or test any rendered message, and charts the surfacing question (C-2 in `charted.md`) to a peer that has no entry for any 0002 data-level category. Nothing in this RDR guarantees the author ever *sees* the fields. | The failure data is byte-perfect and unreachable: no verb the author runs surfaces a `reserved_tag_key` category. The author is told why, in a struct no CLI renders. | §2, §3, premortem, AT-10 |
| C-11 | Approach 1: "the shipped constant and frozen fixture are already conforming, so no kernel code changes" and Consequences: "Positive: the shipped constant, fixture, and D3 precedence all remain valid — the cheapest branch to implement." | "No kernel code changes" survives only three paragraphs. Block 4 requires a new exported predicate and a `Resolve`-entry call; Consequences' last bullet concedes it. But the QOC matrix's *deciding* "blast radius" row and the Trade-offs "cheapest branch" bullet were never rewritten against the settled locus — the fork was decided on a cost claim the design then invalidated. | The RDR was chosen as the zero-change option and lands as the option that puts the first error path into a kernel whose RDR is Final and implemented. Reviewers who read the matrix and not the fine print approve a different change than the one that ships. | §1, §2, premortem, AT-11 |
| C-12 | A5 Evidence: "by P1+P2 it cannot match or write a tag it did not declare; so no owned/observed tag keyed `recognized` reaches `assemble`" | The derivation confuses declaration space with data space. P1/P2 constrain what a *table* may reference. `Input.Owned` and `Input.Observed` are runtime data from the accessor layer and the caller — neither is a table reference, and RDR 0002's declare-every-tag rule does not reach them. So even for a fully conforming table, an owned snapshot carrying `recognized` reaches `assemble` and shadows under D3. A5 marks itself `Verified` by `Method: Derivation` on a premise it does not establish. | The exact silent-shadowing case A5 declares unreachable is reachable for conforming tables — and is now the case block 4 hard-errors on, so the RDR simultaneously claims it cannot happen and errors when it does. | §1, premortem, AT-12 |
| C-13 | Scenario 4 Expected: "the byte-exact rule on the post-parse key string treats `Recognized`, `RECOGNIZED`, and `" recognized"` as ordinary unreserved names (no folding, no trimming)" | The identity rule is authored for the validator's convenience, not the author's. `Recognized` is a legal owned tag name that a human reader cannot distinguish from the reserved word, and the RDR ratifies that as correct behavior with a test that pins it. The advisory-warning question is deferred to "a Resolve question" — but Resolve already ran. | An author writes `[tags.Recognized]` with `provenance = "recognized"`, gets a `reserved_tag_key` failure saying the required name is `recognized`, changes nothing visible, and cannot see the difference. | §1, premortem, AT-13 |
| C-14 | Normative block 5: "the obligation therefore binds any *future* derivation path that does not route through that rule, and is discharged by construction today." | A normative MUST that is discharged by construction, has no authored source form, no lint half, and no test on the producing side is unenforceable text. Scenario 9 tests the *kernel* refusal shape (which this RDR does not change) and explicitly has "**no lint half**." The block binds a future derivation path that no artifact will check. | A later normalizer change derives `RequiresOwned` from a new source, emits `recognized`, and produces `owned_state_unavailable` naming the reserved key — the exact confusing refusal, with a normative clause on the books that nothing executes. | §1, §2, premortem, AT-14 |
| C-15 | A7 **Status: Pending**, with Approach 5, block 5, Capability Dependencies, Phase 2, and scenario 9 all written as settled fact ("is already discharged", "needs no check here"). | The Finalization Gate's own Status-consistency rule forbids exactly this: "no assumption marked `Pending` or `Unverified` may have settled-fact prose elsewhere in the RDR depending on it." Five sections depend on A7's unverified ownership half. The Gate section is an unfilled template, so nothing in the document catches it. | The RDR locks with a dependency on RDR 0007's concurrence that was never obtained; cluster-reconcile later rules the name constraint reaches 0007's contract and block 5 must be withdrawn after implementation started. | §1, §2, premortem, AT-15 |

---

## 1. The three most likely ways implementation goes wrong

### (i) The lint rule's first act is to reject the spec's own canonical fixtures

**Root cause.** A4 verified the wrong proposition and then declared victory over migration.

**Enabling passage.** A4's Evidence: *"Repo-wide whole-token sweep for `recognized` returns zero uses as a tag *key* in owned/observed position… The three committed TOML fixtures use `recognized` only as a provenance value under an author-named key (`0002-…/evidence/spikes/rdr-fixture.toml` and `kata-fixture.toml`: `[tags.outcome]`; `0003-…/evidence/spikes/guard-fixture.toml`: `[tags.rewind_target]`). Every remaining hit is prose. The reserved-key rule invalidates nothing at HEAD, so the mechanical-rename fallback below is not triggered."*

A4 sweeps for the *second* clause of the reserved-key rule — is anything named `recognized` in owned/observed position? — and correctly finds nothing. It never sweeps for the *first* clause, which is the one that constrains real artifacts: a recognized-provenance declaration MUST be named `recognized`. A4 quotes the three fixtures that violate it, verbatim, as evidence that nothing is violated. The RDR literally prints the counterexample inside the assumption that denies it.

The fixtures are not incidental. RDR 0002's Load-Bearing Decisions (`0002…md:361`) reads: *"The RDR and kata spike fixtures are the canonical examples implementation tests must promote."* The `[tags.outcome]` / `provenance = "recognized"` shape at `rdr-fixture.toml:25` and `kata-fixture.toml:22` is the normative worked example of the schema this RDR constrains. RDR 0003's guard fixture (`guard-fixture.toml:36-37`) is a third. Every worked example of a recognized-provenance declaration in this repository is a `reserved_tag_key` failure under block 2.

The Failure Modes section knows this and files it under the *other* problem: *"both canonical RDR 0002 fixtures name the recognized-provenance tag `outcome` … so every existing example of an author naming this datum chose an innocent name."* It cites the fixtures as evidence that the *intent channel* is broad, and never notices that the same sentence is the migration inventory for the naming channel. The observation was made and misfiled.

**Symptom.** RDR 0002's implementer follows RDR 0002's own instruction to promote the canonical fixture into an implementation test, adds this RDR's block-2 check in the same work (Phase 2 puts them in the same path), and the fixture fails to load. The debut of the reserved-key rule is a red test on the spec's reference example. The implementer must then decide, with no ruling from either RDR, whether to change the canonical fixture (touching a Final peer's evidence artifact) or suppress the check.

### (ii) The producer obligation converts working data into a kernel error, and the kernel error is unhandled

**Root cause.** Block 4 asserts a compile-time-shaped obligation over a runtime-shaped channel, and verifies its safety with a literal-source grep.

**Enabling passages.** Block 4: *"Producers of kernel `Input` MUST NOT supply an owned or observed tag keyed `recognized` … A breach is a producer programmer mistake … and travels the Go error path RDR 0001 reserves for programmer mistakes."* And A8's blast-radius half: *"a sweep of `internal/resolve/*_test.go` for the literal `Key: "recognized"` returns exactly one hit … No existing fixture constructs an `Input` that would trip the precondition, so every current call site keeps its nil error."*

The owned snapshot is not written by a programmer. Per RDR 0001's own package doc it is *"an accessor-produced owned tag snapshot"* — RDR 0004's accessor layer reads artifact contents and produces tags whose keys derive from the artifact's own structure. The `rdr-fixture.toml` accessor block (`[accessors.rdr-status] path = "rdr.status"`) shows keys traced from document paths. A key named `recognized` can arrive from data, not from a literal. A grep for `Key: "recognized"` in `*_test.go` proves nothing about that channel; it proves only that no Go literal spells it today. A8 is `Verified` by `Method: Source Search` on a search that structurally cannot see the failure it is clearing.

Block 4 then compounds it with the unconditional clause: *"The check is unconditional on `Input.Recognized`: a reserved-keyed owned or observed tag is a breach whether or not the resolve carries an outcome."* When `Input.Recognized` is empty, `resolve.go:151` injects nothing — there is no recognized entry to shadow, no collision, no harm. The RDR errors anyway, and scenario 6 pins the empty-outcome variant as an expected breach. The predicate is strictly broader than its own justification.

And the error is unhandled by construction. Cove Q1 established that every `return` in `Resolve` ends `, nil`, and the file imports neither `errors` nor `fmt` — the `error` in the signature has never been non-nil in the package's life. Every caller and every test was written against a total function. The RDR treats "the signature already returns `error`" as proof of zero cost (*"(b) adds a check without widening the signature"*), which is true of the type and false of the contract.

**Symptom.** A user's artifact yields an owned tag keyed `recognized`. `Resolve` returns a zero-valued `Result` and a non-nil error. A caller that never had a non-nil error to check sees `Refused() == false` and `Plan == nil`, and either panics on the nil dereference or reports success. The user's data has been reclassified as a programmer mistake, and the diagnosis path RDR 0001 built — typed refusals with `Flow`, `Revision`, `Recognized`, `MissingOwned` — is bypassed entirely.

### (iii) The row still does not fire, inside a model that passes every new check

**Root cause.** Block 2 checks declaration keys only and argues reference positions away with a proof that covers a different case.

**Enabling passage.** Block 2: *"The checked positions are the `[tags.<tag>]` declaration keys only. Reserved-name uses in predicate or write positions (`[rule.match.<tag>]`, `[rule.guard.all.<tag>]`, `[rule.guard.unless.<tag>]`, `[rule.write]`) need no separate reserved-key check: a conforming model declares `recognized` and those references resolve to it, while a model that does not declare it already fails RDR 0002's `unknown tag` rule."*

Consider the model an author actually writes after this rule lands. They have been told the recognized-provenance declaration must be named `recognized`. Their existing model — modeled on the canonical fixture — matched `[rule.match.outcome]`. They rename the declaration to `[tags.recognized]` and, because `outcome` is a name they use elsewhere, keep or add an owned `[tags.outcome]`. Their rules still say `[rule.match.outcome]`. Every check in this RDR passes: the recognized declaration is named `recognized`; no owned/observed declaration is named `recognized`; `outcome` is declared, so `unknown tag` does not fire.

Their row does not fire.

The fall-through clause's disjunction is not exhaustive. It covers (i) a model that declares `recognized` and references it — fine; and (ii) a model that references an undeclared name — caught by 0002. It does not cover (iii) a model that declares `recognized` and references a *different, declared* tag where it meant the recognized one. That is the Problem Statement's opening scenario, and it is precisely the case the reference-position check would have caught. The RDR removed the only check that reaches it and justified the removal on the two cases it does not need.

The same hole opens one level up (C-4): block 2 constrains what a recognized-provenance declaration must be *named* and never requires one to *exist*. Cove Q10 flagged the missing lower bound explicitly. The fix pass responded by adding an upper-bound cardinality paragraph — "at most one" — and left the lower bound unwritten.

**Symptom.** The author reads the failure message, renames the declaration as instructed, re-lints clean, and the row still never fires. The RDR's fix has produced one round of visible remediation and returned them to the original silent failure — now with more confidence that the model is correct, because it linted green.

---

## 2. The one section that will be rewritten within six weeks of shipping

**§ Normative Contracts, block 4 (the `Input` producer obligation) — and with it the Enforcement locus decision that produced it.**

This is the section that will be rewritten, and the rewrite will be forced by an incident rather than by a review.

It is the only block in the RDR that changes running behavior. Blocks 1 and 6 are ratifications ("no kernel code changes"). Blocks 2 and 3 land in a normalizer that does not exist and are explicitly scoped out of this RDR's Done. Block 5 is discharged by construction with no lint half and no producing-side test. Block 4 alone puts new executable code into a Final, implemented kernel — an exported predicate and a `Resolve`-entry call that introduces the package's first non-nil error path.

It will be rewritten for four converging reasons, all visible in the RDR now:

1. **It was added late and settled under review pressure, not designed.** The 3amigo dispositions record it as IMP-2/QA-1, a fork the RDR itself had deferred ("Pick the final form at Pre-Lock") and the lens then *collapsed* rather than escalated. A6 carries a "Residual closed at Pre-Lock" paragraph; A8 was created in the same pass specifically to backfill the load-bearing claim the new decision introduced, and stamped `Verified` immediately — the disposition file says so: *"Verified immediately because both halves were cheap."* A decision this structurally consequential arriving in the final review pass with its supporting assumption authored alongside it is a rewrite waiting for a trigger.

2. **The decision's stated cost basis is already contradicted inside the document.** The QOC matrix's deciding "blast radius" row and Trade-offs' *"the cheapest branch to implement"* were written for the no-kernel-change design. Consequences' final bullet then concedes *"the kernel gains a narrow public surface it did not have."* The Enforcement locus section concedes it again: *"The 'no kernel change' claim in the Decision Rationale is therefore scoped to disposition for conforming input, not to zero surface."* The RDR chose R over M partly on M's *"new public kernel surface on an implemented, Final RDR"* and then took new public kernel surface on an implemented, Final RDR. The moment anyone re-reads the matrix against what shipped, the row is wrong and the block that made it wrong is the thing that moves.

3. **Its safety evidence cannot survive first production data.** A8's blast-radius sweep is a grep of test files for a Go literal (C-5, C-12). The owned channel is accessor-produced from artifact content. The first data-derived `recognized` key turns a resolve into a hard error, and the fix will be to scope the predicate — most likely to `Input.Recognized != ""`, deleting the unconditional clause block 4 currently mandates.

4. **The unconditional clause has no defender.** Block 4's *"The check is unconditional on `Input.Recognized`"* is justified only by symmetry with block 1's rhetorical reservation. When it fires on a resolve where no recognized outcome exists and nothing could have been shadowed, no one will defend it. It goes first, and the block is reopened.

The rewrite will look like: the exported predicate demoted to advisory (candidate (c) alone, which the RDR already names as the fallback in A8's "If wrong"), the `Resolve`-entry call removed, and D3 restored as the sole backstop — which is candidate (a), the option the RDR rejected. It will get there by incident, not by design, which is the expensive way to arrive at the option that was on the table.

---

## 3. The one assumption that will not survive first contact with a real user

**A4 — "No fixture or authored table at HEAD declares or writes an owned/observed tag named `recognized`, so the reserved-key rule invalidates nothing that exists."** `Status: Verified`.

It does not survive because it is not the assumption the design needs. The design needs: *no existing artifact violates the reserved-key rule.* A4 verifies: *no existing artifact names an owned/observed tag `recognized`.* Those are different propositions, and the gap between them is the entire naming rule's first clause.

The evidence paragraph contains its own refutation. A4 writes: *"The three committed TOML fixtures use `recognized` only as a provenance value under an author-named key (`0002-…/evidence/spikes/rdr-fixture.toml` and `kata-fixture.toml`: `[tags.outcome]`; `0003-…/evidence/spikes/guard-fixture.toml`: `[tags.rewind_target]`)."* Read against block 2's first clause — *"A tag declaration with provenance `recognized` MUST be named `recognized`"* — that sentence enumerates three violations. A4 reads it as three passes because it is asking about owned/observed occupancy.

**Why the real user is the one who breaks it.** The first real user of this rule is not an end author — it is RDR 0002's implementer, and their first act is prescribed by RDR 0002's own Load-Bearing Decisions: *"The RDR and kata spike fixtures are the canonical examples implementation tests must promote."* They promote `rdr-fixture.toml`, add the block-2 check from Phase 2 in the same normalizer work, and the reference example fails to load. There is no scenario where this is discovered late; it is discovered on the first run, by the person the RDR delegates its own testing to.

**Why the "If wrong" clause does not save it.** A4's fallback reads: *"implementation must rename the colliding fixture/table tags before the lint lands; scope grows by a mechanical rename, not a design change."* The rename is not mechanical (C-2). `[tags.outcome]` is referenced from `[rule.match.outcome]` at five sites across two fixtures, plus 0003's guard fixture. Renaming the declaration to `recognized` forces every match-position reference to follow, and the RDR nowhere states whether `[rule.match.outcome]` normalizes into `Row.Match` (a tag pattern) or `Row.Outcome` (the first-class non-tag gate that A2 goes to some length to show is unaffected). The implementer cannot tell whether the rename is cosmetic or changes what the rule gates on. That is a design question landing in the middle of what was promised as a mechanical edit — and it lands in a *Final* peer's canonical evidence, which this RDR's Overrides field claims to narrow nothing of.

The second-order damage is worse than the fixture edit. The RDR's Failure Modes reasons from the same fixtures to argue the *intent channel* is broad: *"every existing example of an author naming this datum chose an innocent name."* That is correct and important — and it means the population this RDR aims at has an established convention (`outcome`) that the rule outlaws with no migration story. The RDR's answer, in Risks and Mitigations, is *"none needed as migration — the project carries no back-compat obligation and the reservation lands before any normalizer or authored table exists (A4)."* That mitigation is load-bearing on A4, and A4 is the assumption that fails.

---

## 4. Premortem

*Written from twelve weeks after RDR 0008 was marked Final and implemented.*

**Week 1 — the fixture.** RDR 0002's implementation began, and the implementer did what `0002…md:361` instructs: promoted `docs/rdr/0002-.../evidence/spikes/rdr-fixture.toml` into `internal/model/testdata/` as the schema conformance fixture. They added RDR 0008's block-2 check to the load path in the same commit, per Phase 2. The fixture failed to load: `[tags.outcome]` at line 25 carries `provenance = "recognized"` and is not named `recognized`. The `reserved_tag_key` payload was correct — offending name `outcome`, required name `recognized`, rule `reserved-tag-key/kernel-owned` — and it was pointed at the specification's own reference model. `kata-fixture.toml` failed identically. `0003/evidence/spikes/guard-fixture.toml`, promoted for the guard-evaluator work, failed on `[tags.rewind_target]`.

The implementer went to RDR 0008's A4 to find the migration plan, and found the three failing files quoted as proof that nothing fails. They edited the fixtures — renaming `[tags.outcome]` to `[tags.recognized]` and chasing `[rule.match.outcome]` at `:59`, `:74`, `:86` and at `kata-fixture.toml:47,57` — then stopped, because nothing told them whether `[rule.match.outcome]` had been feeding `Row.Outcome` or `Row.Match`, and the two normalize into different kernel fields with different gating semantics (`Resolve` compares `row.Outcome != in.Recognized` at `resolve.go:332`; `view.matches(row.Match)` at `:335`). They picked `Row.Match`, edited a Final peer's canonical evidence artifact, and moved on. The RDR 0002 conformance suite now pins a fixture that no longer matches the RDR 0002 document.

**Week 3 — the model that lints clean and never fires.** The first internal author migrated a real flow. They renamed the declaration to `[tags.recognized]` as instructed, kept `outcome` as an owned tag for their own reporting, and left `[rule.match.outcome]` in place because nothing told them to change it. Block 2's fall-through clause is explicit that reference positions *"need no separate reserved-key check"* — `outcome` is declared, so `unknown tag` does not fire; the recognized declaration is correctly named, so `reserved_tag_key` does not fire. The model linted green.

At resolve time, `assemble` (`resolve.go:148-164`) injected the outcome under `recognized`; the owned `outcome` tag carried a stale value; `view.matches(row.Match)` returned false for every candidate; `escapeOrRefuse` found no rescue; `refuse` returned `no_match`. The author had followed the remediation text exactly, re-linted clean, and landed back in the Problem Statement's opening paragraph — with more confidence, because the lint was green. They filed it as a resolver bug.

**Week 5 — the accessor.** A flow whose owned snapshot is produced by RDR 0004's accessor layer from a document with a `recognized:` field started returning a non-nil error from `Resolve`. The exported `Input` predicate did exactly what block 4 mandates: an owned tag keyed `recognized`, breach, zero-valued `Result`, non-nil error. A8 had cleared this by grepping `internal/resolve/*_test.go` for the literal `Key: "recognized"` — a sweep that cannot see a key produced from parsed artifact content.

The calling code did not check the error. It had never needed to: before this RDR, every `return` in `Resolve` ended `, nil` and the package imported neither `errors` nor `fmt` (cove Q1). The caller read `res.Refused()` — false, because `Refusal` was nil — took the success branch, and dereferenced `res.Plan`. Nil pointer panic, in a CLI verb, on a user's document. The user saw a stack trace where RDR 0001 had promised *"exactly one disposition — a transition plan or a typed refusal."*

**Week 6 — the empty-outcome variant.** The same predicate fired on a resolve where `Input.Recognized` was empty. `assemble` had injected nothing at `resolve.go:151`; there was no recognized entry, no shadowing, no possible mis-match. Block 4's *"The check is unconditional on `Input.Recognized`"* errored anyway, and scenario 6 pinned that behavior as correct, so the test suite defended it. The fix — scoping the predicate to non-empty outcomes — required contradicting a normative MUST and a green test.

**Week 8 — the surfacing question comes due.** Support asked which command shows a `reserved_tag_key` failure. The answer was: none. RDR 0008's block 3 guarantees the failure *data* carries three fields; `charted.md` C-2 records that RDR 0005's stable code-string table has zero entries for any RDR 0002 data-level category, and Phase 2 holds RDR 0006's `intrastate lint` at arm's length as out of scope. Scenario 7 asserts the payload against a *test stub* consumer, byte-for-byte, and passes. The RDR's committed outcome — *"my row fires, or I am told why"* — was satisfied against a stub and never against a terminal.

**Week 10 — block 5.** A normalizer refactor changed how `Row.RequiresOwned` is derived, taking a path that does not route through `[rule.write]`. A row emitted `RequiresOwned: ["recognized"]`. `missingOwned` (`resolve.go:444-458`) tested `view.has("recognized", ProvenanceOwned)`, found the key present under `ProvenanceRecognized`, and refused `owned_state_unavailable` naming the reserved key — the confusing refusal cove's F-8 identified and block 5 was written to prevent. Block 5's MUST was on the books. It had no lint half, no producing-side test, and A7 was still `Pending`: it was *"discharged by construction"* against a derivation path that no longer existed. Scenario 9 was green throughout — it tests the kernel refusal shape, which this RDR never changed.

**Week 12 — the postmortem.** Read against the RDR, the pattern is uniform: every load-bearing safety claim was verified against a proxy. A4 verified the wrong clause. A5 derived data-channel unreachability from table-declaration rules that do not govern `Input`. A8 verified the absence of a data-derived key by grepping Go literals in test files. Scenario 7 verified user comprehension against a stub. Block 5 verified a producer obligation by asserting it was structurally impossible to breach. The failures were not in what the RDR did not consider — the Failure Modes section names the fixture convention, cove named the `RequiresOwned` channel, and the Enforcement locus section names the surface it concedes. The failures were in the evidence standard: `Verified` was stamped on the nearest available proposition, and the six-block normative surface was defended as "one contract" by an enumeration that had already been extended twice.

The RDR's Done for its own implementation — scenarios 1(kernel half), 3, 6, 8, 9 plus a pointer comment — all stayed green through every one of these incidents. They were green in week 12. None of them touched the outcome the RDR committed to.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

These are review-time tests over the RDR and the repo as it stands — every one is runnable before a line of implementation code exists.

**AT-1 (C-1) — Reserved-key rule vs. committed artifacts**
```gherkin
Given the RDR states "A tag declaration with provenance `recognized` MUST be named `recognized`"
When every committed *.toml under docs/rdr/**/evidence/ is scanned for
     a [tags.<name>] block whose provenance = "recognized"
Then the set of <name> values that are not exactly "recognized" MUST be empty
     OR the RDR MUST carry a named migration item enumerating each one
```
Run today this fails with three entries: `rdr-fixture.toml:25` (`outcome`), `kata-fixture.toml:22` (`outcome`), `guard-fixture.toml:36` (`rewind_target`). A4 would have been forced from `Verified` to `Refuted`, and the "no migration needed" mitigation in Risks would have been struck.

**AT-2 (C-2) — Rename is not mechanical**
```gherkin
Given A4's "If wrong" claims scope grows "by a mechanical rename, not a design change"
When each [tags.<name>] declaration identified by AT-1 is traced to its references
Then every reference site MUST be enumerable, and for each the RDR MUST state
     which normalized kernel field it produces (Row.Outcome or Row.Match)
```
Fails: five `[rule.match.outcome]` sites across two fixtures, and the RDR nowhere maps `[rule.match.<tag>]` to `Row.Outcome` vs `Row.Match`. The rename cannot be executed without a ruling the RDR does not contain.

**AT-3 (C-3) — Reference-position fall-through is not exhaustive**
```gherkin
Given block 2's claim that reference positions need no reserved-key check
When a model declares [tags.recognized] provenance="recognized"
     AND declares [tags.outcome] provenance="owned"
     AND a rule carries [rule.match.outcome]
Then the RDR MUST state whether this model is conforming
     And if conforming, MUST state whether the row fires
```
The RDR has no answer. The model passes every check it adds, and the row does not fire — the Problem Statement's symptom, reproduced inside the fix.

**AT-4 (C-4) — Missing lower bound on the declaration**
```gherkin
Given block 2 constrains the NAME of a recognized-provenance declaration
When a model's rules gate on outcomes but declare no recognized-provenance tag at all
Then the RDR MUST state a validation disposition for that model
```
No clause exists. Cove Q10 raised it; the fix added only an upper bound ("at most one"). Review-time failure: the naming rule has no lower bound, so the affordance can be silently absent.

**AT-5 (C-5) — Producer obligation vs. the accessor channel**
```gherkin
Given block 4 binds "Producers of kernel Input"
When the producers of Input.Owned are enumerated from RDR 0001 and RDR 0004
Then each MUST be identified as constructing keys from Go literals or from parsed data
     And A8's evidence method MUST be able to see the data-derived case
```
Fails: `Input.Owned` is "an accessor-produced owned tag snapshot" (RDR 0001 package doc, `resolve.go:1-11`); A8's method is `grep 'Key: "recognized"' internal/resolve/*_test.go`. A literal grep cannot clear a data-derived channel. A8's `Verified` stamp does not survive this question.

**AT-6 (C-6) — First non-nil error path is a caller-contract break**
```gherkin
Given `Resolve` has returned nil error at every return site since RDR 0001 shipped
When every caller of resolve.Resolve in the repo is enumerated
Then each MUST be shown to handle a non-nil error without treating a zero-valued
     Result as a non-refusal success
```
Fails at review time: cove Q1 established `resolve.go:322,344,348,350,352` all end `, nil` and the package imports neither `errors` nor `fmt`. No caller was written against a fallible `Resolve`. The Enforcement locus' "adds a check without widening the signature" is refuted as a contract claim.

**AT-7 (C-7) — Unconditional check has no harm to prevent**
```gherkin
Given block 4's check is "unconditional on Input.Recognized"
When Input.Recognized == "" and Input.Owned contains a tag keyed "recognized"
Then the RDR MUST name the harm the error prevents in that input
```
Fails: `resolve.go:151` gates injection on non-empty, so nothing is shadowed and no collision exists. Scenario 6 pins the error as expected behavior with no harm behind it. Review-time verdict: the predicate is broader than its justification and the test defends the excess.

**AT-8 (C-8) — Contract count closure**
```gherkin
Given "one independent load-bearing contract" and "They are not separable seams"
When the channels the key can arrive through are enumerated
Then the enumeration MUST be closed by an argument, not by assertion
     And the count MUST NOT have changed during the review rounds
```
Fails on the second clause: the enumeration was two channels (blocks 1–3, 4), then three (block 5 after cove F-8), each version asserting completeness. An enumeration that has been extended under review is not a closure argument. AT-3, AT-4, and AT-14 each name a further gap.

**AT-9 (C-9) — Done does not deliver the committed outcome**
```gherkin
Given the RDR commits to "my row fires, or I am told why"
When "Done for this RDR's own implementation" is enumerated
Then at least one item MUST demonstrate that outcome end to end
```
Fails: Done is scenarios 1(kernel half), 3, 6, 8, 9 plus a pointer comment. Scenario 1's kernel half asserts the view binding (which already works); 3 pins view plumbing; 6 and 8 test the new error path; 9 tests an unchanged refusal shape. Every scenario that delivers the user outcome (2, 4, 5, 7, MVV lint half) is carried into a peer's implementation. The RDR can be marked implemented having delivered nothing the Problem Statement asked for.

**AT-10 (C-10) — Guidance reaches a terminal**
```gherkin
Given block 3 requires three machine-readable fields carrying the guidance
When the path from a `reserved_tag_key` failure to a command an author runs is traced
Then some named verb MUST surface it
```
Fails: `charted.md` C-2 records RDR 0005's code-string table has zero entries for any 0002 data-level category; Phase 2 holds RDR 0006's `intrastate lint` out of scope; scenario 7's consumer is a test stub. The outcome is asserted against a stub and never against a terminal.

**AT-11 (C-11) — QOC matrix consistency with the settled locus**
```gherkin
Given the Enforcement locus was settled to (b)+(c) at Pre-Lock
When the QOC "blast radius" row and the Trade-offs "cheapest branch" bullet are re-read
Then neither MUST claim zero kernel surface
```
Fails: the blast-radius row and Trade-offs still carry the pre-settlement cost basis while Consequences and Load-Bearing Decisions concede the surface. R was chosen over M partly on M's "new public kernel surface on an implemented, Final RDR," and R now takes exactly that. The fork's cost basis is stale.

**AT-12 (C-12) — A5's derivation covers the data channel**
```gherkin
Given A5 derives that no owned/observed tag keyed `recognized` reaches `assemble`
When A5's premises P1 and P2 are checked against Input.Owned / Input.Observed
Then the premises MUST govern runtime input, not only table declarations
```
Fails: P1 is RDR 0002's "The model MUST declare every tag it matches or writes" and P2 is 0002's load-time rejection of undeclared references. Neither reaches `Input.Owned` or `Input.Observed`, which are runtime data, not table references. A5's `Verified`/`Derivation` stamp rests on a scope error — and the RDR then contradicts itself by writing block 4 to hard-error on the case A5 declares unreachable.

**AT-13 (C-13) — Identity rule is human-distinguishable**
```gherkin
Given `Recognized` and `RECOGNIZED` are ruled ordinary unreserved names
When an author declares [tags.Recognized] with provenance = "recognized"
Then the failure text MUST be distinguishable by a human reader from the correct form
```
Fails: the required-name field reads `recognized`, the offending-name field reads `Recognized`, and the RDR defers any case-variant warning to "a Resolve question" — after Resolve has run. Scenario 4 pins the confusing behavior as correct.

**AT-14 (C-14) — Block 5 is enforceable**
```gherkin
Given block 5 is a normative MUST on the normalizer
When the artifacts that would fail if block 5 were violated are enumerated
Then at least one MUST exist on the producing side
```
Fails: block 5 has no authored source form, no lint rule (Phase 2 explicitly declines one), and scenario 9 states "There is **no lint half**." Scenario 9 tests the *kernel* refusal shape, which this RDR does not change. Block 5 is normative text with zero executing artifacts.

**AT-15 (C-15) — Pending assumption has no settled-fact dependents**
```gherkin
Given A7 Status is Pending
When every passage depending on A7 is enumerated
Then none MUST state A7's conclusion as settled fact
```
Fails on five: Approach 5 ("is already discharged"), block 5 ("discharged by construction today"), Capability Dependencies ("discharged today by 0002's `write to non-owned tag`"), Phase 2 ("needs no check here"), scenario 9 ("the normalizer-side obligation is discharged by construction"). This is the Finalization Gate's own Status-consistency rule, verbatim — and the Gate section is an unfilled template, so nothing in the document runs it.
