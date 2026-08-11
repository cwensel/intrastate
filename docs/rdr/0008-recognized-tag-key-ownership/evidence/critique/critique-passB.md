Model: claude-opus-5[1m]
Pass: B (single-model fallback — fresh context, independent of pass A)

# Hostile critique — RDR 0008, Recognized-tag-key ownership

## Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | Normative Contracts block 4: "the kernel MUST export a construction-time predicate over `Input` … and `Resolve` MUST apply that same predicate at entry"; Load-Bearing Decisions / Enforcement locus "(b) and (c) together — one predicate, two call sites" | Collides with RDR 0009's identically-shaped, identically-justified `Resolve`-entry precondition ("The conformance predicate MUST be exported by the kernel package … and Resolve's entry precondition MUST be that same function — one predicate, two call sites"). Two RDRs each write a singular normative "one predicate, two call sites" clause over the same entry point with no ordering, no composition rule, and an explicit *prohibition* on sharing a name. The implementer must invent the composition contract the RDRs refuse to write. | Two `Resolve`-entry validators run in unspecified order; a table that breaches both escape-row shape and the reserved key reports whichever error the implementer happened to sequence first, and the message differs between builds/branches. Downstream 0005 exit-code mapping is nondeterministic for doubly-breaching input. | §1, §4, AT-1 |
| C-2 | A6 "Dependency qualification (RDR 0009 is `Draft`, not Final)" and Joint-check "One inbound reference, not a dependency: RDR 0007 (`Final`) quotes this RDR's `Input` producer sentence in its A13 evidence" | RDR 0007 is **Final** and its A13 "Verified" stamp cites this Draft RDR's producer sentence as evidence that "The single input-side producer obligation in the whole RDR set constrains one reserved key *name*". A Final RDR's verified assumption rests on unlocked Draft text. The RDR calls this "corroborating a negative-existential," which is exactly backwards: 0007 asserts a *uniqueness* claim ("the single … obligation"), and 0008 is the sole witness. If 0008's locus moves, 0007's A13 evidence is stale in a locked document. | Nothing at first. Then a maintainer reading 0007's A13 follows the citation to a sentence that no longer exists in the form quoted, and cannot tell whether the accepted-exposure verdict still holds. The `--tag` CLI channel ships with an evidence chain that dead-ends. | §1, §3, premortem, AT-2 |
| C-3 | Failure Modes, "Silent, residual — **and not a narrow path**"; A4 evidence "both canonical fixtures name the recognized-provenance tag `outcome`"; A1 corroborating "(A1d) … the name `recognized` is unclaimed vacancy, not a collision" | The RDR's own evidence proves the rule misses the actual population. Every committed fixture — `rdr-fixture.toml`, `kata-fixture.toml` (`[tags.outcome]`), `guard-fixture.toml` (`[tags.rewind_target]`) — is a *correct-provenance, wrong-name* declaration, i.e. exactly what the new rule now **rejects**. A4 reports "invalidates nothing at HEAD" by scoping to owned/observed position only, which is the wrong denominator: the rule that bites is block 2's "A tag declaration with provenance `recognized` MUST be named `recognized`", and 3/3 fixtures violate it. | Every committed example table fails load with `reserved_tag_key` the day the normalizer lands. The first thing a new author copies is a fixture, and the fixture does not lint. | §1, §2, §3, premortem, AT-3 |
| C-4 | A4 status "Verified"; "If wrong: implementation must rename the colliding fixture/table tags before the lint lands; scope grows by a mechanical rename, not a design change" | A4 is verified against the wrong predicate. It sweeps for `recognized` "as a tag *key* in owned/observed position" and concludes the rule invalidates nothing. It never sweeps for recognized-*provenance* declarations under other names — the strictly larger violating set, which A1d and Failure Modes both enumerate elsewhere in the same document. The assumption is internally contradicted by two other sections of its own RDR. | The "mechanical rename" fallback is triggered on day one, not never — and it is not mechanical: renaming `[tags.outcome]` to `[tags.recognized]` requires rewriting every `[rule.match.outcome]` reference in both fixtures and both spike `main.go` expectations. | §1, §3, AT-3, AT-4 |
| C-5 | Normative Contracts block 2: "The checked positions are the `[tags.<tag>]` declaration keys only. Reserved-name uses in predicate or write positions … need no separate reserved-key check: a conforming model declares `recognized` and those references resolve to it, while a model that does not declare it already fails RDR 0002's `unknown tag` rule." | The reasoning silently assumes `[rule.write]` is covered by `unknown tag`, but a model that legitimately declares `[tags.recognized]` with `provenance = "recognized"` makes `recognized` a *known* tag — so `[rule.write] recognized = "x"` passes `unknown tag` and must be caught by 0002's `write to non-owned tag`. That category exists, but this RDR never states the dependency in block 2; it states it only in block 5, about `RequiresOwned`. The write-position case is asserted safe by an argument that does not cover it. | An author writes `[rule.write] recognized = "done"`. Whether it is caught depends on whether 0002's implementer wired `write to non-owned tag` to check declared provenance or only accessor binding. If not: a plan carrying a write to a recognized tag reaches RDR 0004's accessor, which normatively "MUST NOT write observed or recognized tags" — an executor-layer failure for a table that linted clean. | §1, §4, AT-5 |
| C-6 | Normative Contracts block 3: "a stable rule identifier whose value is the literal `reserved-tag-key/kernel-owned`" vs. same block "The category's stable data-level value is the token `reserved_tag_key`, matching the snake_case discriminator grammar RDR 0001 uses" | Two machine-comparable identifiers for one failure, in two different grammars (kebab-with-slash and snake_case), both normative, both golden-tested byte-for-byte, with no statement of which one a consumer keys on. RDR 0005's mapping layer is handed two discriminators and no rule for precedence. | A consumer maps on `reserved_tag_key`; a second consumer (or a later renderer) maps on `reserved-tag-key/kernel-owned`. The two disagree after any rename, and the golden tests pin both, so neither can be changed without breaking a test that claims to be normative. | §1, §2, AT-6 |
| C-7 | Normative Contracts block 1: "an absent outcome yields a view with no `recognized` key (unreachable past the outcome-alphabet gate for any alphabet that excludes the empty string)" | Reachable at HEAD and asserted otherwise. `Table.models` is `slices.Contains(t.Outcomes, outcome)`; nothing in RDR 0002 or the kernel forbids `""` in the `outcomes` alphabet, and `assemble` runs *before* the gate (`resolve.go:319-321`). The parenthetical launders a reachable state into an unreachable one by assuming a constraint no peer writes. | A table with an empty-string entry in `outcomes` resolves with a view that has no `recognized` key; a row matching on `recognized` never fires, and the refusal is `no_match` with no indication the key was absent — the exact silent no-match the Problem Statement promises to eliminate. | §1, premortem, AT-7 |
| C-8 | Normative Contracts block 4: "The check is unconditional on `Input.Recognized`: a reserved-keyed owned or observed tag is a breach whether or not the resolve carries an outcome" | Makes `Resolve` fail for inputs that are today legal, well-formed, and harmless — a resolve with no recognized outcome and an owned tag incidentally keyed `recognized` cannot shadow anything, because `assemble` injects nothing. The RDR converts a non-collision into a hard programmer-mistake error, and justifies it only by symmetry with block 1's "unconditional reservation," never by user harm. | A caller with an artifact whose accessor legitimately produces a tag named `recognized` (a plausible domain word: "recognized: true") gets a Go error from `Resolve` with no path to disable it. Accessor-produced owned tags are named by the *artifact*, not by the table author, so this is not a naming choice the reserved-key rule can reach. | §1, §3, premortem, AT-8 |
| C-9 | Approach step 4 / Capability Dependencies "Input-boundary producer enforcement … Available"; A8 "Blast radius: a sweep of `internal/resolve/*_test.go` for the literal `Key: \"recognized\"` returns exactly one hit" | The blast-radius sweep is scoped to a literal that cannot appear in `Input` construction sites, then reports "no existing fixture constructs an `Input` that would trip the precondition" as if the sweep proved it. The sweep searched `Key: "recognized"` — which finds `Row.Match` entries — not `Input.Owned`/`Input.Observed` builders. The conclusion is correct by luck, not by the evidence offered. A8 is stamped "Verified" on a search that does not test the claim. | None immediately; the defect is in the verification, so it recurs the next time someone adds an `Input` fixture and trusts A8's method. | §1, AT-9 |
| C-10 | Load-Bearing Decisions / Identity: "`\" recognized\"` — a whitespace-bearing key — is likewise ordinary and unreserved"; Testing Strategy scenario 4 | Deliberately blesses a near-miss that is invisible in a diff and in most editors. `[tags." recognized"]` declares an owned tag whose key differs from the reserved key by one leading space; it lints clean, and at resolve time it is a *different* key, so the row matching `recognized` never fires while the author sees what looks like the right declaration. The RDR names this and then routes the mitigation to "a Resolve question" — i.e. defers it out of the document that owns the identity rule. | An author declares `[tags." recognized"]`, everything lints clean, and the row silently never fires. This is the Problem Statement's failure, reproduced *by* the fix, and the RDR pins it as correct behavior in scenario 4. | §1, §2, premortem, AT-10 |
| C-11 | A5 "Scope of the guarantee": "the derivation is conditional on the table having passed RDR 0002 validation. Hand-constructed and non-TOML-produced input bypasses that path entirely" + block 6 "D3 remains the deterministic backstop behind the producer obligation for input that bypasses the precondition" | Under the *chosen* enforcement locus, "input that bypasses the precondition" is the empty set — `Resolve` applies the predicate at entry, so nothing reaches `assemble` unchecked. The D3-as-backstop rationale is inherited from the rejected candidate (a) and never re-derived after (b)+(c) was chosen. Scenario 8 then requires a test of that backstop "with the producer precondition bypassed (calling `assemble`'s behavior through the package-internal test path)" — a test of a path the contract makes unreachable. | The implementer writes an internal-only test asserting behavior no caller can observe, and D3's justification in block 6 cites a channel the same document closed. Wasted test, and a stale rationale a future reader will act on. | §1, §2, AT-11 |
| C-12 | Testing Strategy preamble: "Scenarios 2, 4, 5, and 7 land inside RDR 0002's normalizer work (no normalizer exists at HEAD)" and "**Carried into RDR 0002's implementation** as a written obligation this RDR authors and 0002 executes" | The MVV — the thing that proves the RDR solved its stated problem — is split so that the *entire user-facing half* (typed load/lint failure replacing the silent no-match) is unexecutable at lock and delegated to a peer with no mechanism to force it. Prerequisites list "RDR 0002 implementation underway" as an unchecked box. The RDR can be marked done with zero evidence that its Problem Statement outcome was achieved. | The RDR ships "complete" having changed nothing an author can perceive: the kernel gains a predicate, and the diagnosis the author was promised does not exist until an unscheduled peer implementation lands. | §1, §2, §4, AT-12 |
| C-13 | Consequences: "a first-time author's comprehension is therefore asserted at the data level only. Whether that suffices is a question for whichever CLI surface renders it (charted with the user-facing-surface item)" | The RDR's committed outcome is "my row fires, **or I am told why**." It then normatively specifies only three data fields and explicitly disclaims the telling. No peer RDR owns rendering this category: 0005's Predecessors are 0001/0002/0003 and it never enumerates 0002's data-level categories (A1's own evidence: "zero occurrences"). The obligation is charted to a surface that has not agreed to it. | The failure renders through `CLIError`'s generic path as an unmapped category with a `Detail` blob. The author sees a code they cannot look up and three fields with no sentence connecting them. | §2, §3, premortem, AT-13 |
| C-14 | Approach step 5 / block 5: "`RequiresOwned` has no authored source form … so this reservation is a **producer obligation on the normalizer** … and it is already discharged by RDR 0002's `write to non-owned tag` rule" | A normative MUST that the RDR simultaneously declares vacuous ("discharged by construction today", "A normalizer that emits a conforming model cannot produce a breaching row"). It binds a future derivation path that does not exist, tests only the kernel-side residual (scenario 9), and its Capability Dependencies row is stamped **Pending (A7)** — an unverified assumption carrying a normative block into lock. Block 5 is contract-shaped decoration over a no-op. | Nothing, which is the problem: a normative MUST with no enforcement point and no failing test teaches implementers that this RDR's MUSTs are aspirational. | §2, §4, AT-14 |
| C-15 | Problem Statement: "**The outcome this RDR commits to is \"my row fires, or I am told why\" — not \"I get to name this tag.\"** … the naming freedom was a *means* they assumed, not the end" | The RDR re-specifies the user's want on the user's behalf and then declares the substitution non-negotiable, with no user evidence. The only observed authors in the repo — the fixture authors — exercised naming freedom three times out of three (`outcome`, `outcome`, `rewind_target`), which is behavioral evidence *against* the claim that the name was incidental. The XState analogy is misapplied: XState's `event` is a *slot in a struct*, never in the author's tag namespace; `recognized` is a bare word in the same flat `[tags.*]` namespace the author owns. | The author is told their existing, natural spelling is illegal and must be replaced with an engine word, in a system whose entire premise (RDR 0002) is that the model is *reviewable data authors own*. The first reaction is not "a rule I learn once"; it is "why is the engine in my namespace." | §3, premortem, AT-15 |
| C-16 | Metadata Status `Draft` + Overrides "extends RDR 0002's tag-declaration surface with a name constraint" | This RDR imposes a new MUST on the authored TOML surface of a **Final** RDR (0002) and lands the enforcement inside 0002's unwritten implementation, while itself remaining Draft. A1 argues only that the *validation-category list* is extensible; it never argues that the `[tags.<tag>]` **naming grammar** — a distinct Normative Contract in 0002 ("The source schema MUST use the Resolve spike field layout … `[tags.<tag>]`") — is open to new constraints. The extensibility verification covers the wrong contract. | 0002's implementer, reading 0002 alone, builds a normalizer that accepts any `[tags.<name>]`. The constraint arrives later as a bug report against a Final spec that does not mention it. | §1, §4, AT-16 |

---

## 1. The three most likely ways implementation goes wrong

### 1.1 Two RDRs each claim the sole `Resolve`-entry precondition, and neither owns composition

**Root cause in the RDR.** The Enforcement locus decision was made by reasoning about *this* RDR's channel in isolation, then validated against RDR 0009 for the wrong property. A6 asked "does 0009's obligation *reach* `Input`?", answered no, and concluded the locus was "genuinely this RDR's to make." That is the correct answer to a question about scope overlap and the wrong question entirely. The live risk was never that 0009 already covers `Input`; it is that 0009 installs a structurally identical mechanism at the identical function entry, and two mechanisms at one entry point need an ordering.

**The passage that enabled it.** Block 4: "Enforcement is one predicate at two call sites: the kernel MUST export a construction-time predicate over `Input` that producers may call, and `Resolve` MUST apply that same predicate at entry, returning a non-nil error and no `Result` disposition on breach." Compare RDR 0009's normative block verbatim: "The conformance predicate MUST be exported by the kernel package as a construction-time check callable by any table producer … and Resolve's entry precondition MUST be that same function — one predicate, two call sites, so a non-TOML producer can fail at build or construction time rather than at first production Resolve, and the two enforcement points cannot drift."

The phrase "one predicate, two call sites" appears in both documents as a normative singular. 0008 even acknowledges the borrowing — "Pairing them is RDR 0009's own shape" — and then, instead of writing the composition rule, writes a *prohibition*: "the predicate's exported name MUST NOT be bound to any symbol name from RDR 0009." That clause forbids the one thing that would have made the collision visible and manageable (a shared validator surface) while leaving the collision itself unaddressed. It also makes the union hostile to refactoring: an implementer who later wants a single `ValidateInput`-and-`ValidateTable` entry gate is normatively barred from unifying the names.

Note the second-order problem: 0009's precondition is over `Table.Rows` and is evaluated "over the WHOLE table before any evaluation step"; 0008's is over `Input.Owned`/`Input.Observed`. `Input` *contains* `Table`. So 0008's predicate is over a superstructure of 0009's subject, and neither RDR says whether the `Input` predicate subsumes, precedes, follows, or is independent of the `Table` predicate.

**Symptom the user sees.** A table that both carries an escape row with writes *and* is resolved against an `Input` with an owned tag keyed `recognized` produces one error, not two, and which one depends on the order the implementer chose. Because RDR 0005 maps errors to exit codes and the two breaches are different classes, the exit code for a doubly-breaching input is nondeterministic across implementations. Worse, a user who fixes the reported breach gets a *second* error on the next run, with no indication the first fix was incomplete — the classic serial-validator experience the whole "one predicate" framing was supposed to avoid.

### 1.2 The reserved key is enforced against a population the RDR proved it would break

**Root cause in the RDR.** A4 verified the wrong predicate and stamped it Verified, and no later section reconciled the contradiction that A1 and Failure Modes independently surface.

**The passages that enabled it.** A4: "No fixture or authored table at HEAD declares or writes an owned/observed tag named `recognized`, so the reserved-key rule invalidates nothing that exists," with evidence scoped to "zero uses as a tag *key* in owned/observed position." That is true. It is also irrelevant to the rule that will actually reject things. The binding rule is block 2's *first* sentence: "A tag declaration with provenance `recognized` MUST be named `recognized`." The rule bites recognized-provenance declarations under other names — and the RDR itself catalogues them twice:

- A1's corroborating clause (A1d): "both canonical fixtures name the recognized-provenance tag `outcome` (`0002-…/evidence/spikes/rdr-fixture.toml`, `kata-fixture.toml`: `[tags.outcome]` / `provenance = \"recognized\"`)".
- Failure Modes: "both canonical RDR 0002 fixtures name the recognized-provenance tag `outcome` … so every existing example of an author naming this datum chose an innocent name."

I confirmed a third: `docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes/guard-fixture.toml` declares `[tags.rewind_target]` with `provenance = "recognized"`, and `[[rule]] id = "reconcile-rewind-legality"` guards on `[rule.guard.all.rewind_target]`. Three committed fixtures, three violations, 100% of the observed population.

A4's "If wrong" is therefore already true at authoring time — and it is worse than A4's own estimate of it ("scope grows by a mechanical rename"). Renaming `[tags.outcome]` → `[tags.recognized]` in `rdr-fixture.toml` also requires editing `[rule.match.outcome]` in three rules; `kata-fixture.toml` in two; and the 0002 spike's `main.go` and `output.txt`, which encode the normalized dump. The 0003 fixture's `rewind_target` carries a `domain = ["resolve", "refine", "propose"]` and a guard reference; renaming it to `recognized` makes the guard read `[rule.guard.all.recognized]`, which is legible but no longer says what it means — the fixture's whole point was that the recognized datum in that flow is a *rewind target*.

**Symptom the user sees.** The day the 0002 normalizer lands with this rule, every committed example fails `reserved_tag_key`. A first-time author copies `rdr-fixture.toml`, runs lint, and is told the shipped example is invalid. Their reasonable conclusion is that the tool is broken, not that the example is stale.

### 1.3 The write- and reference-position exemption rests on an argument that does not cover the conforming case

**Root cause in the RDR.** Block 2 disposes of four syntactic positions with one sentence of reasoning that splits on "conforming vs. non-conforming model" and never considers the conforming-model-with-illegal-write case.

**The passage that enabled it.** "The checked positions are the `[tags.<tag>]` declaration keys only. Reserved-name uses in predicate or write positions (`[rule.match.<tag>]`, `[rule.guard.all.<tag>]`, `[rule.guard.unless.<tag>]`, `[rule.write]`) need no separate reserved-key check: a conforming model declares `recognized` and those references resolve to it, while a model that does not declare it already fails RDR 0002's `unknown tag` rule."

Trace `[rule.write] recognized = "x"` against that dichotomy. The model *is* conforming — it declares `[tags.recognized]` with `provenance = "recognized"`, exactly as block 2 requires. So the reference "resolves to it" and `unknown tag` does not fire. The only thing standing between that write and a `Plan` is RDR 0002's `write to non-owned tag` category — which block 2 never names. Block 5 names it, but block 5 is about `RequiresOwned`, a derived field, and its argument is explicitly about *derivation*, not about authored write blocks.

The dependency is load-bearing and unstated. And it is not obviously safe: `write to non-owned tag` must be implemented as a check against *declared provenance* for this to hold. If 0002's implementer instead keys it on accessor binding (a plausible reading — 0002's tag declarations carry `accessor = ...` precisely for owned/observed read-back), a `recognized`-provenance tag with no accessor is "non-owned" in the sense of *unbound*, and the check may or may not reject it.

**Symptom the user sees.** A table that lints clean produces a `Plan` whose `Writes` contains `{Key: "recognized"}`. That plan reaches the RDR 0004 accessor layer, where the normative rule is unambiguous: "A write accessor MUST apply only planned owned-tag writes … It MUST NOT write observed or recognized tags." So the failure lands at *execution* time, three layers past the lint that was supposed to catch it, as an accessor-layer error the author has no vocabulary to connect back to a tag name. This is a strictly worse version of the Problem Statement's silent no-match: later, less diagnosable, and after side effects may have begun.

---

## 2. The one section that will be rewritten within 6 weeks of shipping

**`#### Normative Contracts`, block 3 — the failure-payload contract — together with the Consequences bullet that disclaims its rendering.**

Block 3 will be rewritten because it is the only part of this RDR that a human being will actually encounter, and it is specified in a form that cannot survive contact with the renderer.

What block 3 commits to: three machine-readable fields (offending name as authored, required name, rule identifier), a rule identifier pinned byte-for-byte as `reserved-tag-key/kernel-owned`, a category token `reserved_tag_key`, and an explicit refusal to specify wording — "human-readable wording is the renderer's to choose."

Four forces converge on it inside six weeks:

1. **Two identifiers, one failure (C-6).** The same block normatively pins `reserved-tag-key/kernel-owned` *and* `reserved_tag_key`, in two different grammars, with no statement of which one is the discriminator. RDR 0005's mapping layer needs exactly one. The first time someone writes the map, one of the two is demoted — and both are pinned by golden tests the RDR calls normative, so the demotion is a spec change, not a refactor.

2. **The renderer disclaimer is unowned (C-13).** Consequences says comprehension "is a question for whichever CLI surface renders it (charted with the user-facing-surface item)." No such surface has agreed. RDR 0005's Predecessors are 0001/0002/0003 — not 0008 — and this RDR's own A1 evidence establishes that 0005 "never enumerates RDR 0002's data-level categories at all (zero occurrences)." So the category renders through the generic `CLIError` path. The first bug report is "the error doesn't tell me what to do," and the fix is to move wording into block 3 — reversing its central design choice.

3. **The payload is under-specified for the case it exists to serve.** Block 3 requires "the offending name as authored" and "the required name." For the *most common* violation established in §1.2 — `[tags.outcome]` with `provenance = "recognized"` — those two fields render as `outcome` and `recognized`, which does not say the offense was the *provenance-name pair*. An owned tag named `recognized` produces the same two fields with the same shape and the opposite required action (rename away from `recognized`, not toward it). The payload cannot distinguish the two directions, so any renderer built on it produces a message that is wrong half the time. A fourth field — the violated direction — gets added, and block 3's "three distinct machine-readable fields" is rewritten.

4. **Block 3 extends a peer's pre-existing category by fiat.** "Where the offending key is the reserved name, this payload requirement extends RDR 0002's pre-existing `unknown tag` failure; that extension is authored by this RDR and lands in 0002's implementation." This bolts a conditional payload onto a category 0002 owns, triggered by inspecting the offending key. 0002's implementer will either miss it or push back, and the reconciliation rewrites the clause.

Runners-up, for the record: Testing Strategy (rewritten when the 0002-side half of the MVV stays unexecutable — C-12) and Load-Bearing Decisions / Identity (rewritten the first time `[tags." recognized"]` reaches a user — C-10).

---

## 3. The one assumption that will not survive first contact with a real user

Not A1 through A8. The assumption that dies is the one stated as settled fact in the Problem Statement and never registered as an assumption at all:

> **"The outcome this RDR commits to is 'my row fires, or I am told why' — not 'I get to name this tag.'** Those diverge, and the divergence is a deliberate product stance rather than an incidental constraint. The author's underlying want is a working match with a diagnosable failure; the naming freedom was a *means* they assumed, not the end."

This is a claim about user psychology, load-bearing for the entire QOC matrix's "correctness fit" row — the row the Decision Rationale names first and calls decisive — and it carries no evidence whatsoever. It is not an A-record. It has no Method, no Status, and no "If wrong."

It will not survive, for four grounded reasons.

**The repo's only observed authors contradict it.** Three committed TOML fixtures declare a recognized-provenance tag. All three chose a domain name: `outcome`, `outcome`, `rewind_target`. That is 3/3 exercising the "means" the RDR says was incidental. And `rewind_target` is the telling one: it is not a lazy synonym for `recognized`, it is a *semantically specific* name in a flow where the recognized datum is a rewind destination, guarded by `[rule.guard.all.rewind_target] in = ["resolve", "refine", "propose"]`. Renaming it `recognized` destroys information the author deliberately encoded. The RDR's own Failure Modes concedes the pattern is "common rather than exotic" — and then does not revisit the stance the concession undermines.

**The XState analogy is category-wrong.** The RDR leans on it repeatedly: "on the same footing as XState's fixed `event` slot: the engine-injected datum has an engine-owned name." But XState's `event` is a *field of a struct handed to the guard* (`GuardArgs` — `{ context, event }`). It occupies no position in the author's namespace; authors name event *types*, and the RDR says so. `recognized` is the opposite: it is a bare key in the flat `[tags.*]` namespace, adjacent to `status`, `profile`, `stage`, `owner`. The engine is not reserving a slot in its own struct; it is claiming a word inside the author's dictionary. The nearest true analogue in this repo is 0002's `<clear>` sentinel — which the RDR itself notes is "syntactically un-declarable" and therefore *cannot* collide, and whose mechanism it then explicitly rejects (Load-Bearing Decisions / Naming). The one in-repo precedent that actually supports the stance is the one the RDR declines to follow.

**The system's stated premise is author ownership.** RDR 0002's entire framing is a *reviewable* model authors maintain by hand. "The rendered table is review support, not the source authors maintain." In that context, a reserved bare word is not a rule learned once; it is the engine appearing inside the artifact whose reviewability is the product.

**Accessor-produced names are not the author's to change (see C-8).** Owned tags come from artifacts via accessors. If an artifact has a field named `recognized`, the author cannot rename it — and block 4 makes supplying it a hard Go error, unconditionally, even when no recognized outcome is present and no collision is possible.

**What first contact looks like.** An author writes `[tags.outcome] provenance = "recognized"` — because every example they have seen does. Lint rejects it. They fix it to `[tags.recognized]`, and now their rules read `[rule.match.recognized] eq = "round-clean"`, which reads as a tautology. They file: "why can't I call it `outcome`?" The honest answer — "so a rename mistake becomes a lint error instead of a silent no-match" — is a good answer to a question they never asked, about a failure they had not hit. The RDR has no user-facing story for this exchange, because it decided in advance that the question was illegitimate.

---

## 4. Premortem: written as if it already failed

**Date: eleven weeks after RDR 0008 locked. Subject: why `flow lint` was reverted and RDR 0002's implementation slipped a cycle.**

RDR 0008 locked cleanly. Its own implementation was small and green: a pointer comment on `internal/resolve/resolve.go::recognizedTagKey`, an exported `ValidateInput`-shaped predicate over `Input`, a call at the top of `Resolve`, and four kernel-side tests (scenarios 3, 6, 8, 9). Everything the RDR could execute at HEAD, it executed. That is precisely how the failure hid.

**Week 2 — the two preconditions meet.** RDR 0009 landed its own `Resolve`-entry precondition. Both RDRs had written "one predicate, two call sites"; neither had written what happens when there are two predicates and one entry. 0008 additionally forbade sharing 0009's symbol name, so unifying them was normatively barred. The implementer stacked them — `ValidateInput(in)` then `ValidateTable(in.Table)` — because 0008's predicate takes the wider argument. Nobody wrote down that the order was a choice.

Three weeks later a contributor reversed them while extracting a helper. `TestReq6_ResolveErrorPathIsProgrammerMistakesOnly` still passed; every test asserted "some non-nil error," never *which*. The exit code for a doubly-breaching table flipped from `flow-bad-input` to `flow-bad-table` in a patch release, and RDR 0005's exit-code contract — the thing 0005 exists to hold stable — broke without a single failing test.

**Week 4 — `assemble` runs before the gate.** A flow shipped with `""` in its `outcomes` alphabet (a legitimate encoding for "recognizer produced nothing"). `Resolve` called `assemble(in)` at line 319, before `in.Table.models(in.Recognized)` at 321. With `in.Recognized == ""`, `assemble` injected nothing. `Table.models("")` returned true, because `""` was in the alphabet. Every row matching on the tag `recognized` failed `view.matches`, and the kernel returned `no_match`. Block 1 had declared this "unreachable past the outcome-alphabet gate for any alphabet that excludes the empty string" — a constraint no RDR imposes and `Table.models` does not check. The user got the Problem Statement's exact silent no-match, from the document that promised to abolish it.

**Week 5 — the fixtures do not lint.** The 0002 normalizer reached the point of running `reserved_tag_key` against the committed spikes. `rdr-fixture.toml` failed. `kata-fixture.toml` failed. `guard-fixture.toml` failed. Three for three. A4 said "the reserved-key rule invalidates nothing at HEAD, so the mechanical-rename fallback below is not triggered" — verified against owned/observed position only, while the rule that fired was the recognized-provenance naming rule that A1d and Failure Modes had *both* documented as universally violated in the same document.

The rename was not mechanical. `rdr-fixture.toml` needed `[tags.outcome]` → `[tags.recognized]` plus three `[rule.match.outcome]` edits; `kata-fixture.toml` two; both spikes' `main.go` and the frozen `output.txt` dumps had to be regenerated, which meant re-deciding what the normalized dump ordering was supposed to be. `guard-fixture.toml`'s `rewind_target` — a recognized tag with a real domain, read by `[rule.guard.all.rewind_target]` — could not be renamed without destroying what the fixture was demonstrating, so RDR 0003's spike evidence was marked stale.

**Week 6 — the write position.** A flow author, now correctly declaring `[tags.recognized] provenance = "recognized"`, wrote `[rule.write] recognized = "done"` to record the outcome. Block 2 had exempted write positions on the grounds that "a conforming model declares `recognized` and those references resolve to it" — true, and the reason nothing rejected it. `unknown tag` did not fire; the tag was known. `write to non-owned tag` had been implemented against accessor binding rather than declared provenance, so it did not fire either. Normalization emitted a `Row` whose `Writes` carried `{Key: "recognized"}`. `resolve.go::planOf` copied it into the `Plan` unconditionally. The RDR 0004 accessor layer refused at execution — "MUST NOT write observed or recognized tags" — after the transition had been reported as successful. The author's diagnosis started at the accessor and never reached the tag declaration.

**Week 7 — the leading space.** A different author hand-typed `[tags." recognized"]`. Byte-exact, no trimming, by normative decision: unreserved. It linted clean as an owned tag. The row matching `recognized` never fired. `no_match`. Testing Strategy scenario 4 asserted this outcome as *correct*, and it passed.

**Week 8 — the payload says nothing.** The `reserved_tag_key` failures from week 5 rendered through `internal/cli/clierr::CLIError`'s generic branch, because nothing in RDR 0005 mapped the new category — exactly as A1's own evidence had established ("RDR 0005 never enumerates RDR 0002's data-level categories at all"). Users saw a code they could not look up and a `Detail` blob with `outcome`, `recognized`, and `reserved-tag-key/kernel-owned`. Nobody could tell whether the fix was to rename *toward* `recognized` or *away* from it — the payload's three fields are identical in shape for both violation directions. The golden test asserting all three fields byte-for-byte passed throughout.

**Week 9 — a Final RDR's evidence goes stale.** Fixing the enforcement locus meant editing block 4's producer sentence. RDR 0007 — `Final` — cites that sentence verbatim in A13's Verified evidence as proof that "The single input-side producer obligation in the whole RDR set constrains one reserved key *name*, not provenance origin." 0008's Joint-check had called this "one inbound reference, not a dependency … it corroborates a negative-existential rather than supplying a premise." But A13's claim is a *uniqueness* claim, and 0008 was its only witness. With the sentence edited, a locked document's verified assumption pointed at text that no longer existed. Per house rule, RDRs are never amended. The evidence chain for the `--tag` accepted-exposure simply dead-ended.

**Week 11 — the accessor case.** A flow bound an artifact whose frontmatter carried `recognized: true`. The accessor produced it as an owned tag. Block 4's check is "unconditional on `Input.Recognized`," so `Resolve` returned a Go error even on resolves carrying no recognized outcome, where `assemble` injects nothing and no collision is possible. The tag name came from the artifact, not from a table author, so no reserved-key rule could reach it and no rename was available. The precondition was made conditional in a hotfix, and block 4's "unconditional" MUST was breached by the shipping code within eleven weeks of lock.

**Post-mortem line.** The RDR was validated as a *document* — every A-record stamped, every peer cross-referenced, contract count argued down to one — and never as a *system*. Its own text contained the disproof of A4 in two separate sections. Its enforcement mechanism was a verbatim copy of a peer's mechanism at a shared entry point, and the peer comparison it ran asked only whether the *scope* overlapped. And the half of its Minimum Viable Validation that would have exposed all of this was, by its own Testing Strategy, "carried into RDR 0002's implementation" — which is to say, not run.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

These are review-time gates. Each is executable against the RDR corpus and the repo as they stood at review.

**AT-1 — one `Resolve`-entry precondition, or a written composition rule (C-1)**
```gherkin
Given every RDR whose Normative Contracts place a precondition at resolve.Resolve's entry
When more than one such RDR exists
Then exactly one of them MUST own the entry-gate composition (ordering, short-circuit, and
     whether one predicate's argument subsumes the other's)
And each participating RDR MUST cite that owner
And no participating RDR may forbid sharing a symbol name with another participant
     without naming the mechanism that keeps the two from drifting
```
Result at review: FAIL. 0008 block 4 and 0009's normative block both say "one predicate, two call sites" over `Resolve` entry. No composition owner exists. 0008 forbids name-sharing and supplies no alternative drift control.

**AT-2 — no Final RDR's Verified evidence may rest on Draft text (C-2)**
```gherkin
Given RDR 0007 has Status Final and assumption A13 stamped Verified
When A13's evidence quotes a sentence from an RDR whose Status is Draft
And that quotation supports a uniqueness claim ("the single input-side producer obligation")
Then the Draft RDR MUST NOT lock until either
     (a) the quoted sentence is marked frozen against edit, or
     (b) RDR 0007's A13 records the dependency and its re-verification trigger
```
Result at review: FAIL. 0007 A13 cites 0008's producer sentence as sole witness to a uniqueness claim; 0008's Joint-check classifies this as "not a dependency" and locks.

**AT-3 — the rule is run against every committed fixture (C-3, C-4)**
```gherkin
Given the RDR introduces a naming rule over [tags.<tag>] declarations
When the rule is applied to every committed *.toml fixture under docs/rdr/*/evidence/spikes/
Then the RDR MUST state the pass/fail count
And any failing fixture MUST be listed by path with its required edit
```
Result at review: FAIL, 0/3 passing.
- `0002-…/evidence/spikes/rdr-fixture.toml` — `[tags.outcome]` / `provenance = "recognized"` → violates block 2
- `0002-…/evidence/spikes/kata-fixture.toml` — `[tags.outcome]` / `provenance = "recognized"` → violates block 2
- `0003-…/evidence/spikes/guard-fixture.toml` — `[tags.rewind_target]` / `provenance = "recognized"` → violates block 2

**AT-4 — an assumption's evidence must test the assumption's own rule (C-4, C-9)**
```gherkin
Given an assumption claims "the rule invalidates nothing that exists"
When its Method is Source Search
Then the search predicate MUST be the rule's own predicate, not a narrower one
```
Result at review: FAIL twice.
- A4 searches "`recognized` as a tag key in owned/observed position"; the rule also forbids recognized-provenance declarations under other names. The narrower predicate hides 3 violations.
- A8 searches `Key: "recognized"` in `*_test.go` (a `Row.Match` shape) and concludes "No existing fixture constructs an `Input` that would trip the precondition." The search cannot find `Input.Owned`/`Input.Observed` construction. Correct conclusion, invalid evidence.

**AT-5 — every exempted syntactic position names its covering rule (C-5)**
```gherkin
Given block 2 exempts [rule.match.<tag>], [rule.guard.all.<tag>],
      [rule.guard.unless.<tag>], and [rule.write] from the reserved-key check
When the model IS conforming (it declares [tags.recognized] with provenance recognized)
Then for each exempted position the RDR MUST name the peer rule that rejects a
     reserved-key use, and state the implementation property that rule depends on
```
Result at review: FAIL. Block 2's dichotomy ("conforming → resolves to it" / "non-conforming → `unknown tag`") leaves `[rule.write] recognized = ...` in the conforming branch with nothing rejecting it. The covering rule (`write to non-owned tag`) is named only in block 5, about a different field, and its dependency on provenance-based checking is never stated.

**AT-6 — one failure carries one discriminator (C-6)**
```gherkin
Given block 3 pins machine-comparable identifiers for one failure category
When more than one identifier is pinned byte-for-byte
Then the RDR MUST name which identifier consumers key on, and what the other is for
```
Result at review: FAIL. `reserved_tag_key` and `reserved-tag-key/kernel-owned` are both normative, both golden-tested, in two grammars, with no precedence.

**AT-7 — unreachability claims are checked against the shipped control flow (C-7)**
```gherkin
Given block 1 claims a state is "unreachable past the outcome-alphabet gate"
When resolve.go is read
Then the gate MUST actually precede the state's construction
And the constraint the claim assumes MUST be normative in some RDR
```
Result at review: FAIL. `assemble(in)` is line 319; `in.Table.models(in.Recognized)` is line 321 — the view is built *before* the gate. And `Table.models` is `slices.Contains(t.Outcomes, outcome)`; no RDR forbids `""` in `outcomes`. Both halves of the parenthetical are false.

**AT-8 — a precondition rejects only inputs that can cause the harm (C-8)**
```gherkin
Given block 4's check is "unconditional on Input.Recognized"
When Input.Recognized is empty
Then assemble injects no recognized key, so no owned/observed tag can shadow it
And the RDR MUST justify rejecting that input by harm, not by symmetry with block 1
And the RDR MUST state the remedy for an owned tag whose key comes from an artifact
     (accessor-produced) rather than from a table author
```
Result at review: FAIL. Justification is symmetry only ("matching block 1's unconditional reservation"). Accessor-produced owned keys are never considered anywhere in the RDR.

**AT-9 — kept as part of AT-4 (C-9).** See above.

**AT-10 — a near-miss the rule blesses must have a detection story (C-10)**
```gherkin
Given the Identity rule declares " recognized" (leading space) an ordinary unreserved name
When an author declares [tags." recognized"] as an owned tag
Then the row matching `recognized` does not fire, with no diagnostic
And the RDR MUST NOT route the mitigation to a later stage while pinning the
     silent outcome as expected in its own Testing Strategy
```
Result at review: FAIL. Scenario 4 asserts the silent-pass as **Expected**; the mitigation is deferred ("a Resolve question, not identity"). The Problem Statement's failure is reproduced and blessed.

**AT-11 — a backstop's justification must survive the chosen locus (C-11)**
```gherkin
Given the enforcement locus settled on (b)+(c): Resolve applies the predicate at entry
When block 6 justifies D3 as "the deterministic backstop … for input that bypasses the precondition"
Then the RDR MUST identify a caller that can reach assemble without passing the entry check
Or the backstop rationale MUST be re-derived for the chosen locus
```
Result at review: FAIL. Under (b) no such caller exists at the package boundary. Scenario 8 concedes it by requiring "the package-internal test path" — a test of a path the contract makes unreachable. The rationale is a leftover from rejected candidate (a).

**AT-12 — the MVV must be executable by the RDR that owns it (C-12)**
```gherkin
Given the Minimum Viable Validation is split across this RDR and a peer
When the peer's half carries the user-facing outcome named in the Problem Statement
Then that half MUST be executable at lock, or the RDR MUST NOT claim the Problem
     Statement outcome as delivered
```
Result at review: FAIL. Scenarios 2, 4, 5, 7 and "the lint half of the MVV" are "carried into RDR 0002's implementation." The typed load/lint failure — the entire user-visible fix — is unexecutable at lock, and the Prerequisite ("RDR 0002 implementation underway") is an unchecked box with no scheduling authority behind it.

**AT-13 — a committed user outcome needs an owner for the surface that delivers it (C-13)**
```gherkin
Given the RDR commits to "my row fires, or I am told why"
When the telling is disclaimed to "whichever CLI surface renders it"
Then some peer RDR MUST list this RDR as a predecessor, or name the category, or
     accept the obligation in its own text
```
Result at review: FAIL. RDR 0005's Predecessors are 0001/0002/0003. This RDR's own A1 evidence establishes 0005 "never enumerates RDR 0002's data-level categories at all (zero occurrences)." No surface has accepted the obligation.

**AT-14 — a normative MUST needs a reachable breach (C-14)**
```gherkin
Given block 5 states "Row.RequiresOwned MUST NOT name recognized"
When the RDR also states it is "discharged by construction today" and
     "A normalizer that emits a conforming model cannot produce a breaching row"
Then the RDR MUST name a test that fails if the MUST is removed
And its supporting assumption MUST NOT be Pending at lock
```
Result at review: FAIL. Scenario 9 tests the kernel-side *residual* (`owned_state_unavailable` naming the key), which is present with or without block 5. A7 is **Pending**, and the Capability Dependencies row is stamped Pending — an unverified assumption carrying a normative block into lock.

**AT-15 — a re-specified user want needs evidence (C-15)**
```gherkin
Given the Problem Statement replaces the user's stated want ("I get to name this tag")
      with a different one ("my row fires, or I am told why")
When that substitution is the deciding row of the QOC matrix (correctness fit)
Then it MUST be registered as a Critical Assumption with Method, Status, and "If wrong"
And the observed behavior of every author in the repo MUST be reported for or against it
```
Result at review: FAIL. The claim carries no A-record. Observed behavior: 3/3 fixture authors chose a domain name (`outcome`, `outcome`, `rewind_target`), and one of those (`rewind_target`) encodes flow-specific meaning that renaming destroys. Evidence runs against the stance; the stance is asserted anyway.

**AT-16 — extending a Final peer's grammar requires verifying the right contract (C-16)**
```gherkin
Given this RDR adds a MUST over RDR 0002's authored [tags.<tag>] surface
When A1 verifies extensibility
Then A1 MUST address the contract being extended — the source-schema field layout —
     not only the validation-category list
```
Result at review: FAIL. A1 verifies only that "Validation failures MUST retain stable data-level categories … including at minimum …" is a floor. The naming constraint lands on a *different* 0002 normative block ("The source schema MUST use the Resolve spike field layout: root `outcomes`, `[model]`, `[tags.<tag>]`, …"), which carries no extensibility fence and was never examined.
