Model: claude-opus-5

# Hostile critique — RDR 0003, iteration 3

## Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `trace` step 6, `:1141` ("Compute coverage … OK") against `evidence/spikes/guard-fixture.toml` | The document's only end-to-end worked example asserts coverage OK for a group that does **not** cover `profile = small`, and whose `cluster_eligible` carries no `optional` marker so the default-optional rule adds a `{present,absent}` presence dimension the trace never counts. The one desk-check that would have caught a coverage-arithmetic bug was run without doing the arithmetic. | The MVV's "one complete partition" case is built from this fixture, fails on first run, and the team spends a sprint deciding whether lint is wrong or the fixture is — with no worked example in the RDR to arbitrate. | §1, premortem, AT-1 |
| C-2 | Cardinality clause `:1053-1055` — "the product of every participating dimension's **declared domain size**"; kind table `:841` (`set` → "its element universe") | A `set` dimension's assignment count is `2^N` over an element universe of size N, not N. The bound clause specifies the *wrong quantity* for exactly the kind the too-large clause exists to guard. `contains` is also the one operator with no evaluated-plus-declared evidence anywhere. | A model with three 20-element set tags reports cardinality 8000 (under the bound, "proved") while the real assignment space is 2^60. Lint hangs or OOMs, or silently proves nothing — the "bound discovered by exhausting memory is not a conforming bound" clause is violated by the clause's own arithmetic. | §1, §3, premortem, AT-2 |
| C-3 | Single-valued clause `:810-818` ("Absent the marker … one independent boolean dimension per value (`2^\|domain\|`)") vs. optionality default `:775-777` ("no optionality marker declares the key optional") | Two independently-defaulted-to-worst-case rules multiply. An undeclared 5-value enum is `2^5 = 32` assignments, times a presence dimension = 64, where the author meant 5. The RDR never states the two defaults compose, and never gives the composed cardinality for any real tag. | Every author's first model refuses with `graph-unprovable-coverage` (or too-large) until they discover two markers nobody told them are required. The RDR's own §"Risks" mitigation ("declarations feel like boilerplate → keep them close to tag definitions") does not touch this. | §1, §2, §3, premortem, AT-3 |
| C-4 | Disposition table `:1114` ("Full `unless` block decides true \| Row disabled \| Excluded intersection subtracted") vs. RDR 0007 `:1435-1440` (`all_result ∧ ¬(unless_conj)`, `¬U = U`) | The RDR states `unless` semantics as a two-valued set subtraction while the kernel it delegates to computes it in three-valued Kleene where `¬U = U`. An unevaluable atom **inside `unless`** makes the whole row unevaluable at runtime, but the lint model subtracts a set as if it were decided. The disposition table has no row for "`unless` atom unevaluable". | A row whose `unless` guards an optional key is certified green by lint and refuses `guard_unevaluable` at runtime — the exact false-green the narrowing clause exists to forbid, arriving through the one block the narrowing's worked reasoning never examines. | §1, premortem, AT-4 |
| C-5 | Participation clause `:910-919` ("Match keys are not product dimensions") vs. `guard-fixture.toml` rule `reconcile-rewind-legality` (`[rule.guard.unless.status]`) and rule `profile-to-grounding` (`[rule.match.status]`) | The clause partitions keys by *authoring block*, not by key. The RDR's own fixture uses `status` as a match key in one rule and a guard key in another. The clause gives no rule for a key that is both, and `resolve.go::Row` carries `Match` and `Guard` as separate fields, so the ambiguity survives into the shipped shape. | Two implementers produce different scoped products for the same model — one includes `status` (3 values), one does not — so cardinality, coverage verdict, and the too-large bound all diverge across builds. A17's whole worry (cross-implementation determinism) reappears in the dimension count. | §1, §2, premortem, AT-5 |
| C-6 | A17 `:526-546` + two-population overlap clause `:943-958` + disposition `:1125` | A17 reads `resolve.go::Resolve` step 2 correctly but stops one function short. `gate()` (`resolve.go::gate`) returns `blocked` for `guard_unevaluable` and `owned_state_unavailable` **before** the exact-one count, so `Resolve` returns `refuse(in, *blocked)` and **never calls `escapeOrRefuse`**. An escape row therefore cannot rescue those two classes at all — yet the RDR treats escape rows as a coverage-closing population and cites this kernel as its authority. | An author writes a bare escape row expecting it to close coverage; lint certifies the group green on the union identity; at runtime a `guard_unevaluable` on any ordinary row refuses without ever consulting the escape row. Green lint, refused resolution — P5 breached by the clause that cites the kernel. | §1, premortem, AT-6 |
| C-7 | Coverage clause `:936-941` ("An escape row carrying no guard atoms denotes the whole scoped product and therefore closes coverage by itself") | Under the two-population overlap rule, a bare escape row closes coverage while generating **no** overlap finding against any peer. This is a one-line, always-available, blocking-free exemption from the RDR's central guarantee — precisely the "opt-out" the RDR spends two paragraphs (`:685-698`, "Leaving a dimension undeclared is not an opt-out") forbidding for declarations. | The first team member who hits a `graph-coverage-gap` adds a bare escape row. Lint goes green forever. The exhaustiveness guarantee is dead in the repo within a month and nothing reports it. | §1, §2, premortem, AT-7 |
| C-8 | A14 `:434-464` "Plan: single-field request to RDR 0002" vs. RDR 0002 `:309-311` ("Normalization MUST combine both into one candidate-row predicate set before ambiguity checks") | A14 calls this "a gap to close, not a contradiction to resolve." RDR 0002's clause is normative and says *combine*; the identity tuple at `:1162-1170` requires `block` to survive as a per-atom field. That is a contradiction between two normative clauses, downgraded to a request by asserting compatibility without citing a reading under which both hold. | Implementer builds normalization to RDR 0002's letter, then MVV Scenario 7 fails, then the identity tuple is not total, then the same-key-both-blocks row (RDR 0007 A5's conjoined case) reports one atom where two exist. Reopened as an RDR-level dispute mid-implementation. | §1, §2, premortem, AT-8 |
| C-9 | Metadata `:22-24` + Prerequisites `:1567-1659`: A10, A12, A14, A15, A16, A17 all Pending, four homed at *other documents' refine passes* | Six of seventeen assumptions are open; four are homed at RDR 0006's and RDR 0002's refine, both `Draft`. Every one is individually argued "not lock-blocking." Nobody argues the aggregate. The phase table `:1699-1704` then declares Phases 1, 2, 3 "partially" startable and Phase 4 "No" — i.e. the RDR locks a contract whose integration phase is admittedly not startable. | Implementation starts, produces Phase 1/2/3 work against unconfirmed group construction (A10), unconfirmed overlap population (A17), unwritable declaration fields (A16), and uncarried atom blocks (A14) — then redoes it when the peers refine. The "abandon code and iterate RDR" header line fires. | §1, §2, premortem, AT-9 |
| C-10 | A15 `:465-499` + MVV Scenario 3 `:1784-1800` — "Reading B at fixture setup is what makes the case runnable" | The test reads the implementation's published bound `B` at setup and builds products relative to it. A test that derives its own oracle from the system under test cannot falsify the claim it is named to discharge: whatever `B` the implementation publishes, the fixture is constructed to straddle it, so the test passes by construction unless the implementation is internally inconsistent. A15 asks whether cardinality *predicts provability*; this test never measures provability. | A15 is stamped closed at implementation on a test that could not have failed. The first model whose shape (not cardinality) exceeds the proof representation refuses unpredictably, and the "same model, same verdict on every conforming implementation" promise is discovered false in production. | §1, §2, premortem, AT-10 |
| C-11 | Withheld-claim clause `:968-975` ("MUST name the participating row and the atom that can refuse") + the `Consequent duty` note `:977-984` | The RDR states a MUST whose producer field it simultaneously documents as nonexistent in RDR 0006's finding contract (`0006:363-367`), then declares this "does not gate this RDR's lock." A normative MUST that no consumer can satisfy is not a contract; it is a scheduled defect. | MVV Scenario 8 asserts the finding "names that row and atom." RDR 0006's finding record has no atom field. The test cannot be written, or is written against a field that does not exist, and Scenario 8 — the RDR's flagship false-green defense — is the first test deleted. | §1, §2, premortem, AT-11 |
| C-12 | A6 `:234-245` + owned-tag clause `:1075-1078` ("every reachable predecessor sets or preserves that tag") | The RDR keeps a normative MUST quantifying over "every reachable predecessor" while A12 concedes no citable predecessor relation exists, and the cluster gate flagged that RDR 0006 frames it as *graph reachability* while this RDR requires a *syntactic decision*. Two different decision procedures over one predicate, both normative, neither chosen. | Lint either rejects rows it should accept (syntactic reading: any owned tag not written in the same rule) or accepts rows it should reject (graph reading with an incomplete traversal). Authors see `graph-owned-before-write` fire on legal flows and stop trusting the lint. | §1, premortem, AT-12 |
| C-13 | Illustrative Code `:1219-1245` — the only declaration example, marked "illustrative"; `single_valued` spelling explicitly disclaimed at `:1243-1245` | The RDR owns a five-field type system and gives exactly one example, two of whose fields it disowns the spelling of. The one artifact an implementer would copy is stamped non-normative. No example anywhere shows a `set` declaration with `contains`, a computed product cardinality, or a withheld claim. | The implementer invents the wire spelling; RDR 0002's author invents a different one; the fixture in `evidence/spikes/` (which declares none of the five fields — no `optional`, no `elements`, no `single_valued`) matches neither. Three spellings in the repo at once. | §1, §2, §3, premortem, AT-13 |
| C-14 | Scope Verification `:2010-2025` + Proportionality `:2051-2064` — "Scope grew deliberately at Stage 6 … That is one contract, not two" | The RDR absorbed an entire tag type system at Stage 6 — after grounding, 3amigo, and critique had already run against a document that did not contain it — and then re-ran only "that row." The five-field declaration model has no Alternatives section, no prior-art pass, and no assumption record about *its own* design (only about whether it exists). Six alternatives are analyzed for the operator grammar; zero for the type system. | The type system is wrong in a way nobody looked for: no nullable-vs-optional distinction, no `int` beyond `{min..max}`, no string/scalar domain, no unit or ordering for enum comparison, no versioning story for adding an enum value to a shipped domain. Each is discovered by a user, not a reviewer. | §2, §3, premortem, AT-14 |
| C-15 | Metadata Status block `:9-25` (17 lines of cross-document bookkeeping) and the assumption records generally — A8 `:276-306`, A17 `:522-546` | The document has become a ledger of its own governance. A8's record is 30 lines about which document *records* a clause both documents agree on. The reader who needs to implement guard evaluation must first parse a four-document dependency graph, two JDR decisions, one re-opened JDR decision, and a cluster-gate report. | Onboarding an implementer takes days. The implementer reads the Normative Contracts block and skips the assumptions — where every unresolved arithmetic question actually lives. | §2, premortem |
| C-16 | `disposition` table `:1110-1126` — no row for a value atom whose key is *declared always-present* but absent in a non-conforming view | The conformance clause `:787-798` declares non-conforming views out of scope ("this RDR MUST NOT be read as promising anything about a view that violates the declarations") and routes the check to "the kernel's view assembly" — but no clause requires anyone to *perform* that check, and `resolve.go::assemble` performs none. Conformance is a premise with no enforcer. | An observed tag's accessor returns nothing for a key declared always-present. Lint proved coverage over a product with no `{absent}` assignment. Runtime refuses `guard_unevaluable`. Green lint, refused resolution, and the RDR's answer is that the view was out of scope. | §1, §3, premortem, AT-15 |
| C-17 | A1 `:139-173` + Testing Strategy `:1758-1767` | A1 is `Verified` on a harness that "carries its own TOML-to-atom encoding because the kernel reshape is unshipped, and **consumes no declared domain**." It proves that eight operator *spellings* evaluate under an encoding the implementation will not use, against a fixture declaring none of the five fields this RDR now owns. It is evidence for a vocabulary claim about a document that has since changed its central contract. | The closed-vocabulary claim survives lock on a spike that predates the declaration model. Phase 3 discovers `contains` over a *declared* element universe needs semantics the harness never exercised (subset vs. superset vs. non-empty-intersection), and the operator set reopens after lock. | §1, §3, premortem, AT-16 |

---

## 1. The three most likely ways implementation goes wrong

### 1.1 The product cardinality is computed wrong, and the RDR's own worked example proves nobody checked

**Root cause in the RDR.** Three clauses each define a piece of the scoped product's size, and no passage anywhere multiplies them out against a real model.

- The cardinality clause (`:1053-1055`) fixes the counted quantity as "the number of assignments in it, the product of every participating dimension's **declared domain size**."
- The kind table (`:841`) declares a `set` tag's finite domain to be "its element universe."
- The single-valued clause (`:810-818`) says an unmarked finite-domain tag contributes "one **independent boolean dimension per value** (`2^|domain|`)."
- The optionality clause (`:775-777`) says "A declaration carrying **no optionality marker declares the key optional**," and the presence-dimension clause (`:922-934`) says an optional key contributes a `{present, absent}` dimension.

Compose them on a `set` tag with a 20-element universe: the cardinality clause reports 20. The actual assignment count is `2^20`. The clause specifies the wrong quantity by an exponential factor for exactly the kind whose blowup the too-large bound was written to catch. Compose them on an unmarked 5-value enum: `2^5` value assignments times 2 presence assignments = 64, where every author on earth means 5. The RDR states each rule in isolation and never once shows the composed number.

**The specific passage that enabled it.** The `trace` mini-check (`:1128-1158`) is the document's only end-to-end walk and the place this would have surfaced. Step 3 (`:1138`) builds the product as "`profile` = 4 enum values; `prelock_iterations` = `{0..3}`" — 16 assignments — and stamps OK. Step 6 (`:1141`) computes coverage and stamps OK. Now open `evidence/spikes/guard-fixture.toml`, which the trace names as its witness source:

```toml
[tags.profile]
domain = ["small", "mid", "large", "foundational"]

[[rule]]
id = "profile-to-grounding"
[rule.guard.all.profile]
in = ["mid", "large"]

[[rule]]
id = "foundational-to-cove"
[rule.guard.all.profile]
eq = "foundational"
```

`profile = small` is accepted by **no row in the group**. The union is `{mid, large} ∪ {foundational}` = three of four values. The group has a coverage gap, and the trace says OK. Worse: `cluster_eligible` in that same fixture declares no `optional` marker, so under `:775-777` it is optional, so under `:922-934` it contributes a presence dimension — which step 3 does not include in the product and step 5 (`:1140`) discusses without ever adding to the count. The trace's step 5 even walks the `exists` atom in detail and still does not put the dimension into step 3's product.

The trace concludes (`:1147`) "No CONTRADICTION row and no GAP row." The fixture it walks contains a coverage gap. The check that exists to catch arithmetic errors was performed without arithmetic.

**Symptom the user sees.** The MVV requires "one complete partition" (`:1665-1666`) and Phase 3 builds it from "representative RDR and kata guards" — i.e. this fixture. On first run, lint reports `graph-coverage-gap` on the group the RDR certified. The implementer has no worked example in the RDR to check against, because the only one is wrong. Half a sprint goes to deciding whether the lint, the fixture, or the RDR is at fault, and the likeliest resolution is the wrong one: relax the coverage check until the fixture passes.

### 1.2 `unless` is specified as set subtraction in a system that computes it in three-valued logic

**Root cause in the RDR.** This RDR delegates atom-verdict combination to the kernel (`:875-880`, "the combination of per-atom verdicts belong to the kernel (JDR 0001 §D4, normative in RDR 0007)") and then models `unless` as a Boolean set operation in its own proof.

RDR 0007's normative combination clause (`docs/rdr/0007-guard-predicate-totality.md:1435-1440`) is:

> The row verdict is `all_result ∧ ¬(unless_conj)`, where `unless_conj` is the conjunction of the `unless` block's atoms (`unless` is block-level negation, NOT per-atom negation).

with `¬U = U` and `T ∧ U = U` (`:1412-1421`). So if any `unless` atom is UNEVALUABLE and no `unless` atom is FALSE, `unless_conj ∈ {T, U}`; if it is `U`, `¬(unless_conj) = U`, and a row whose `all` block is fully TRUE has verdict `U` — `guard_unevaluable`, a blocking refusal.

RDR 0003's model of the same object (`:698-701`):

> Within that group, a row's accepted assignments are the intersection of all positive `all` atom domains minus the single conjunctive assignment set matched by the row's full `unless` block.

and the disposition table (`:1114`):

| Full `unless` block decides true | Row disabled | Excluded intersection subtracted | none | Silent by design |

That is two-valued. The disposition table has rows for "value atom over an absent key" and "existence atom over an absent key" but **no row for an unevaluable atom inside `unless`**. The narrowing clause (`:961-965`) and the "can refuse" clause (`:1011-1023`) do formally cover it — "a value atom over a key **declared optional**" does not say which block — but nothing in the coverage derivation, the trace, the disposition table, or the MVV walks the case, and the subtraction language actively contradicts it: you cannot both subtract a set and withhold the claim.

**The specific passage that enabled it.** A3 (`:191-197`) stamps `Verified` by **Design Decision** — "This RDR chooses separate positive and negative guard lists" — and never checks its own choice against the kernel's aggregation table. A3's "If wrong" is about author ergonomics ("authors will duplicate rows"), not about three-valued semantics. The one assumption record that owns `unless` never looks at how `unless` is computed.

**Symptom the user sees.** An author writes an ordinary exclusion — `unless prelock_iterations gte 3` — on a key that is optional because they never wrote `optional = false` (see 1.3). Lint certifies the group green: the subtraction is performed, the union closes, the claim is issued. At runtime, the key is absent, `unless_conj = U`, the row verdict is `U`, and `resolve.go::gate` returns a `guard_unevaluable` refusal. Green lint, refused resolution — exactly the failure the narrowing clause (`:961-965`) is the RDR's centerpiece defense against, arriving through the one block the defense's worked reasoning never touches.

### 1.3 Two conservative defaults compose into a product nobody can prove, on a document that never composes them

**Root cause in the RDR.** The RDR sets two independent defaults to their worst case:

- No optionality marker ⇒ **optional** (`:775-777`), justified as "the conservative default, since assuming always-present would let lint certify a group green on an undeclared property."
- No single-valued marker ⇒ **values co-occur**, `2^|domain|` (`:810-818`), justified as "the model does not assume at least one holds."

Both justifications are individually sound. Together, an author who declares `kind = "enum"` and `domain = [a,b,c,d,e]` — the complete declaration the RDR's own illustrative block shows minus two lines — gets a dimension of `2^5 × 2 = 64` assignments, of which 59 are combinations the tag can never hold. Four such tags in a group: `64^4 ≈ 1.7 × 10^7`. Add the `set` miscount from 1.1 and you exceed any plausible published bound `B` on a five-tag model.

**The specific passage that enabled it.** The single-valued clause spells out its arithmetic (`:810-818`) and the optionality clause spells out its default (`:775-777`), and the two live 30 lines apart in the same `Normative Contracts` block, and no passage anywhere multiplies them. The `Risks and Mitigations` entry that should have caught it (`:1519-1521`) says:

> **Risk**: Finite-domain declarations feel like boilerplate. **Mitigation**: Keep declarations close to tag definitions in RDR 0002's model and make diagnostics explain when a missing domain blocks proof.

That mitigation addresses a *missing domain*. It does not address a *present domain with two missing markers*, which is the actual failure — the author declared the domain and still cannot prove anything. The `Illustrative Code` block (`:1219-1231`) shows `profile` with all five fields, which is exactly the declaration that would *not* fail, and then disclaims its own spelling (`:1243-1245`).

**Symptom the user sees.** The first real model an author writes — declaring `kind` and `domain` on every tag, as the RDR's stated MUST (`:742-746`) requires and its MAY (same) makes the rest optional — refuses. Every group takes the blocking inability-to-prove outcome, which is A11's own stated "If wrong" (`:381-383`): "no row group is ever exhaustiveness-eligible, every group takes the blocking inability-to-prove outcome, and this RDR's central guarantee is unreachable in practice while remaining true in theory." A11 is stamped `Verified` and that consequence arrives anyway, because A11 verifies that the fields *exist*, not that their defaults are usable.

---

## 2. The one section that will be rewritten within 6 weeks of shipping

**`Normative Contracts` — the declaration model clauses (`:741-861`), specifically the four defaults and the domain/kind agreement table.**

This is the newest text in the document, added at Stage 6 on 2026-08-21 and extended on 2026-08-22 by §JD-13, and it is the only load-bearing contract in the RDR that has never been through a full lens cycle in its current form. Scope Verification concedes the timing (`:2013-2016`, "Scope grew deliberately at Stage 6"), and Proportionality re-argues the *placement* (`:2051-2064`) — but neither re-argues the *design*. Six alternatives are analyzed for the operator grammar. Zero for the type system. The five fields have no prior-art pass, no rejected-alternative record about the fields themselves, and one illustrative example whose spelling the document disowns.

Three forces guarantee the rewrite:

1. **The defaults are unusable as composed** (§1.3). The first week of real authoring produces a demand to invert one or both — probably "no single-valued marker means single-valued," since that is what every author means and what the fixture in `evidence/spikes/` already assumes by declaring none. Inverting a default in a normative clause is a rewrite of the clause and of every consumer citing it (RDR 0002's authoring schema, RDR 0006's `graph-single-valued-state`, the cardinality bound).

2. **The kind table is missing the cases real models need.** There is no distinction between "may be absent" and "may hold null." There is no `int` form except `{min..max}` — no enumerated int set, no open-above bound. `scalar` carries no domain and therefore poisons any group it appears in with a blocking finding (`:750-755`), which means the escape hatch for "I have a string tag" is *the lint always fails*. The first author with a `string` tag they only ever compare with `eq` will demand a finite string domain, and the kind table gains a sixth row or `scalar` gains a domain.

3. **A16 guarantees the clause gets reopened anyway.** The single-valued marker's authoring location is a Pending request to RDR 0002 (`:500-521`). When RDR 0002's refine takes it up, that document's authors will read this clause for the first time as consumers, and the `set`-kind rejection (`:841`) plus the "meaningful only for a kind carrying a finite domain" restriction (`:804`) will get argued. A field whose meaning lives in one document and whose spelling lives in another is renegotiated the first time both are edited.

The rewrite is not a refinement. Changing a default changes every previously-green model's verdict, which is a breaking change to the acceptance gate the whole cluster is built on.

---

## 3. The one assumption that will not survive first contact with a real user

**A11 (`:357-383`) — "The tag declaration model carries a finite domain, so any exhaustiveness claim has a producer."**

A11 is stamped `Verified` by **Design Decision**, and the decision it records is *that the fields exist in this document*. Its evidence (`:361-377`) is entirely an ownership argument: RDR 0002 has no normative tag-type vocabulary, RDR 0007 says value kinds are RDR 0003's, therefore the model is homed here. Every word is about *where* the model lives. Not one word is about whether the model, once written, produces a finite domain for a tag a real author declares.

The claim that will not survive is the causal one in A11's own title: *carries a finite domain* ⇒ *any exhaustiveness claim has a producer*. Between those two, the RDR inserted:

- a default that makes every unmarked key optional (`:775-777`) — adding a presence dimension the author did not ask for,
- a default that makes every unmarked finite-domain tag a power set (`:810-818`) — squaring the dimension the author did ask for,
- a `scalar` kind that can never bear a claim (`:750-755`) — so any string tag is fatal,
- a cardinality rule that counts a `set`'s element universe rather than its assignment space (`:1053-1055`) — so the bound does not measure the thing it gates,
- and a `MAY` on four of the five fields (`:742-746`) — so declaring the domain is optional, and the conservative defaults fire whenever it is omitted.

A real user's first model declares `provenance` and `kind` (both MUST) and `domain` on the tags they thought about, and omits `optional`, `single_valued`, and `elements` — the same shape as `evidence/spikes/guard-fixture.toml`, which is the only real model in this repo and which declares **none** of the three. That model produces zero provable row groups. A11's stated "If wrong" (`:381-383`) describes this outcome precisely, and A11 is `Verified` anyway, because the verification asked "do the fields exist?" and the failure mode is "do the defaults compose?"

The tell is in the record itself. A11's rejected alternative (`:378-380`) is "filing four fields into RDR 0002" — a placement alternative. A Design Decision whose only rejected alternative is *where to put the decision* has not evaluated the decision.

---

## 4. The premortem

*Written from twelve weeks after lock.*

Phase 1 landed clean. RDR 0007's kernel reshape shipped first, `Row.Guard` became `[]Atom`, and `guard.Parse` produced the four-field atoms with the operator/kind matrix enforced. `guard.Evaluate` decided value semantics over a present value and never read the tag view, exactly as `:875-880` specified. Nobody argued with any of it. The operator vocabulary was right. The atoms-not-callbacks decision was right. Everything this RDR was *about* was right.

It failed at the arithmetic, and it failed in the first hour of Phase 3.

**Week 4.** `lint.ScopedProduct` was written from `:1053-1055` — "the product of every participating dimension's declared domain size." `lint.dimensionSize` returned `len(decl.Domain)` for `enum`, `2` for `bool`, `max-min+1` for `int`, and `len(decl.Elements)` for `set`, because the kind table (`:841`) names the element universe as the `set` kind's finite domain and the cardinality clause says to multiply declared domain sizes. Nobody caught that a set of 20 elements has 2^20 assignments and not 20; the RDR never showed a computed cardinality for any model, so there was nothing to check against. The `lint.tooLargeBound` constant was published at 10^6 as the clause required, and reported beside the computed cardinality in every refusal, and the whole thing was internally consistent and wrong.

**Week 5, day 1.** Phase 3 encoded `evidence/spikes/guard-fixture.toml` as the MVV's "one complete partition" case, because the RDR's `trace` (`:1136-1145`) walks that exact fixture and stamps every step OK. `lint.CheckGroup` reported `graph-coverage-gap` on the `status eq Draft` group: `profile = small` accepted by no row. Two days went to deciding whether `lint.coverageUnion` had a bug. It did not. The RDR's own desk trace had certified a group with a hole in it, and the fixture had been carrying that hole since iteration 1 because no pass ever expanded the four-value enum by hand.

The team's fix was the one available under deadline: they added `small` to `profile-to-grounding`'s `in` list. The fixture went green. Nobody asked why the RDR's trace said OK, and the trace's method — walk the clauses, stamp OK — stayed the team's model of what "checked" meant.

**Week 5, day 3.** The same fixture's `cluster_eligible` declares no `optional` marker. `decl.Optional` defaulted true per `:775-777`, so `lint.presenceDimension` adjoined `{present, absent}`, so the group's product doubled, so `foundational-to-cove`'s value atoms left the `absent` half uncovered, so the group failed again. This one was harder: the RDR's step 5 (`:1140`) discusses the `exists` atom at length and never adds the dimension to step 3's product, so the document simultaneously requires the dimension and demonstrates a product without it. The team read `:930` — "A key declared always-present contributes no presence dimension" — and added `optional = false` to `cluster_eligible`. Which made the `exists` atom vacuous, which the operator matrix (`:633`) says lint "reports as such rather than rejecting it," which meant the fixture's only `exists` atom now produced an advisory instead of exercising A7's presence projection. The RDR's flagship derivation went untested in its own MVV and nobody noticed, because the fixture was green.

**Week 6.** First real flow. Eight tags, five enums, one `set` for labels, two ints. The author declared `provenance`, `kind`, and `domain` on all eight and stopped — the same shape as every fixture in the repo. `lint.dimensionSize` applied `:810-818`: no `single_valued` marker, so each 4-to-6-value enum became `2^n` independent boolean dimensions; no `optional` marker, so each got a presence dimension too. `lint.ScopedProduct` returned 4.3 × 10^11. `lint.tooLargeBound` was 10^6. Every group in the flow refused with `graph-unprovable-coverage`, and the diagnostic dutifully reported both figures exactly as `:1049-1051` requires, which told the author their model was 400,000× too big and nothing about why.

The author's question — "why is a five-value enum thirty-two dimensions?" — has an answer in `:810-818` and the answer is correct and the author did not care. They wrote `single_valued = true` on eight tags, which RDR 0002's schema did not yet accept (A16, `:500-521`, still Pending at RDR 0002's refine), so the field was dropped by `table.Normalize` before `lint` ever saw it, so nothing changed, so they filed a bug against lint.

**Week 7 — the failure that mattered.** While the cardinality fight ran, a second flow shipped through with green lint. It had one row:

```toml
[rule.guard.all.stage]
eq = "prelock"
[rule.guard.unless.prelock_iterations]
gte = 3
```

`prelock_iterations` had no `optional` marker. `lint.rowAccepted` implemented `:698-701` — intersection of `all` domains minus the `unless` subtraction — and the presence dimension was in the product, and the subtraction was performed over the `present` half, and the union closed, and the group certified green.

Then `resolve.Resolve` ran on a real transition where the accessor had not yet produced `prelock_iterations`. `guard.Evaluate` was never called for that atom — the kernel decided presence first, per `:875-880`, and marked the atom UNEVALUABLE. `unless_conj` was UNEVALUABLE. Per RDR 0007's normative table, `¬U = U`, `T ∧ U = U`, and the row verdict was `GuardUnevaluable`. `resolve.gate` collected it into `undecidable`, returned `blocked = &Refusal{Kind: KindGuardUnevaluable}`, and `Resolve` refused.

Green lint. Refused resolution. The single promise this RDR exists to make — `:718-719`, "a green exhaustiveness result must mean resolution succeeds" — broken by a two-line guard.

The narrowing clause (`:961-965`) should have caught it. The "can refuse" clause (`:1011-1023`) formally covers it: a value atom over a key declared optional, and `prelock_iterations` was optional by default. But `lint.canRefuse` had been implemented from the disposition table (`:1110-1126`), which is where implementers go for input-class-to-outcome mapping, and the disposition table has a row for "Value atom over an absent key" and a row for "Full `unless` block decides true → Excluded intersection subtracted, Silent by design" and **no row for an unevaluable atom inside `unless`**. The implementer wired `canRefuse` over `all` atoms, because that is where the table's absent-key row lives and because `:698-701` says `unless` is *subtracted*, and subtraction is a thing you do to a decided set.

**Week 8.** The escape-row discovery. `graph-coverage-gap` was firing across the repo and the fixes were tedious, and someone read `:936-941`:

> An escape row carrying no guard atoms denotes the whole scoped product and therefore closes coverage by itself.

together with `:952-954`: escape rows "are therefore excluded from the ordinary-row overlap check." One bare escape row per flow closed coverage, generated no overlap finding against anything, and turned every red group green. It spread through four flows in a week. It was fully conforming with the RDR.

Then a flow with a bare escape row hit a `guard_unevaluable` on an ordinary row. `resolve.gate` returned `blocked` for the `KindGuardUnevaluable` refusal, and `Resolve` returned `refuse(in, *blocked)` on the spot — it **never called `escapeOrRefuse`**. The escape row that closed lint's coverage cannot rescue a `guard_unevaluable` at all; `escapeOrRefuse` is only reachable from the `len(selected) == 0` and `default` arms, after `blocked == nil`. A17 (`:522-546`) had read `Resolve` step 2 carefully enough to see that ordinary candidates exclude escape rows, and had cited that reading as its whole evidence base, and had not read `gate` — the function immediately after — where two of the four refusal classes are decided before escape rows are ever consulted.

**What the postmortem said.** Not that the guard model was wrong. The symbolic-atom decision, the closed operator set, the `all`/`unless` split, the value-only evaluator seam: all correct, all shipped, all still in the codebase.

What was wrong was that a 2,092-line document specifying a finite-domain proof contained no computed finite domain. Its one desk trace stamped OK on a fixture with a coverage hole. Its one illustrative declaration disclaimed its own field spelling. Its `set` cardinality was off by an exponent. Its two conservative defaults were each defended in isolation and never multiplied. Its `unless` was specified as subtraction in a system that computes it in Kleene logic. Six of seventeen assumptions were open at lock, four homed at other documents' unwritten refine passes, each individually declared "not lock-blocking" and never assessed together.

Every one of those was checkable at review time with a calculator and the fixture already in the repo.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

These are RDR-review gates, not implementation tests. Each is executable against the document and the repo as they stood at iteration 3.

**AT-1 — the desk trace must compute, not stamp.** (catches C-1)
```gherkin
Given the `trace` mini-check names `evidence/spikes/guard-fixture.toml` as its witness
When a reviewer enumerates the scoped product for the `status eq "Draft"` group
  by listing every assignment of every participating dimension
And computes union(row_accepted) for that group
Then the trace row for step 6 must record the computed cardinality,
  the computed union size, and the uncovered assignments by name
And the step MUST NOT be stampable "OK" while any assignment is uncovered
```
Run today: `profile` domain is `{small, mid, large, foundational}`; row accepted sets are `{mid,large}` and `{foundational}`; `small` is uncovered. Step 6 fails. `:1141` says OK.

**AT-2 — every kind's assignment count must be stated, not just its domain.** (catches C-2)
```gherkin
Given the kind table at `:836-842` names a finite domain per kind
And the cardinality clause at `:1053-1055` counts "declared domain size"
When a reviewer writes, for each of the five kinds, the number of distinct
  values the dimension can take in one conforming view
Then the table MUST carry that number as its own column
And for `set` with element universe of size N the column MUST read 2^N
And the cardinality clause MUST reference that column, not "domain size"
```
Run today: the column does not exist; `set` reads "its element universe," and the cardinality clause multiplies that. Fails.

**AT-3 — the defaults must be multiplied out on a stated model.** (catches C-3)
```gherkin
Given the optionality default at `:775-777` and the single-valued default at `:810-818`
When a reviewer takes the minimal conforming declaration
  ("provenance", "kind", "domain" only — the shape of every fixture in the repo)
And computes the scoped product for a group over four such 5-value enum tags
Then the RDR MUST state that number
And MUST state whether it falls under any plausible published bound
And if it does not, the RDR MUST justify the defaults against that number
  or change them
```
Run today: `(2^5 × 2)^4 = 64^4 ≈ 1.7 × 10^7`. No passage computes it. Fails.

**AT-4 — `unless` must be specified in the logic the kernel computes it in.** (catches C-4)
```gherkin
Given RDR 0007 `:1435-1440` fixes the row verdict as all_result ∧ ¬(unless_conj)
And RDR 0007 `:1412-1421` fixes ¬U = U and T ∧ U = U
When a reviewer enumerates the nine (all_result, unless_conj) verdict pairs
Then the disposition table at `:1110-1126` MUST carry a row for every pair
  whose row verdict is U
And in particular MUST carry the pair (all_result = T, unless_conj = U)
And the coverage derivation at `:698-701` MUST state what a row's accepted
  assignment set is when its `unless` block is undecidable
```
Run today: no such row exists; the derivation says "subtracted." Fails.

**AT-5 — a key used as both match key and guard key must have a stated rule.** (catches C-5)
```gherkin
Given the participation clause at `:910-919` partitions by authoring block
When a reviewer greps the RDR's own fixture for a key appearing under both
  [rule.match.*] and [rule.guard.*]
Then the RDR MUST state whether that key is a product dimension,
  and MUST state it for the case where the two rules are in the same group
```
Run today: `status` appears as `[rule.match.status]` in `profile-to-grounding` and as `[rule.guard.unless.status]` in `reconcile-rewind-legality`. No clause covers it. Fails.

**AT-6 — every claim citing `resolve.go` must be traced to the refusal, not the candidate build.** (catches C-6)
```gherkin
Given A17 and the two-population clause cite `resolve.go::Resolve` step 2
When a reviewer follows the control flow from the cited line to every `return`
Then the RDR MUST state, for each of the four RefusalKinds,
  whether an escape row can rescue it
And the coverage clause at `:936-941` MUST NOT let a bare escape row close
  coverage for a refusal class the kernel cannot route to escapeOrRefuse
```
Run today: `gate` returns `blocked` for `KindGuardUnevaluable` and `KindOwnedStateUnavailable` before the exact-one count, so `escapeOrRefuse` is unreachable for both. The RDR states neither. Fails.

**AT-7 — no authoring gesture may close coverage without a matching check.** (catches C-7)
```gherkin
Given the RDR forbids silent opt-out at `:685-698`
When a reviewer searches for any single authored construct that turns a group
  from "coverage gap" to "green" with no accompanying finding
Then no such construct may exist
Or the RDR MUST state the finding that accompanies it
```
Run today: a bare escape row closes coverage by itself (`:938-940`) and is excluded from the ordinary overlap check (`:952-954`). One line, always available, no finding. Fails.

**AT-8 — a "gap" against a peer's normative MUST must be shown compatible, not asserted.** (catches C-8)
```gherkin
Given A14 calls the block-retention question "a gap to close, not a contradiction"
And RDR 0002 `:309-311` normatively requires normalization to
  "combine both into one candidate-row predicate set"
When a reviewer is asked for one reading under which both clauses hold
Then A14 MUST cite that reading in its Evidence line
```
Run today: A14 asserts "The container-level combination is compatible with per-atom `block` retention" (`:446-447`) and cites its own spike's encoding, which A14 itself says "corroborates rather than discharges" (`:454-455`). No reading is given. Fails.

**AT-9 — the aggregate of open assumptions must be assessed, not only each one.** (catches C-9)
```gherkin
Given six assumptions are Pending, four homed at peer documents' refine passes
When a reviewer lists, for each implementation phase, the Pending records that
  gate it (from the phase table at `:1699-1704`)
Then the RDR MUST state which phases produce work that survives every
  possible resolution of its gating records
And Phase 4 being "startable: No" MUST be reconciled with locking the contract
```
Run today: Phases 1, 2, 3 are "Partially"; Phase 4 is "No"; no passage assesses the aggregate. Fails.

**AT-10 — a test that reads its oracle from the system under test does not discharge an assumption.** (catches C-10)
```gherkin
Given MVV Scenario 3 constructs products relative to the implementation's
  published bound B, read at fixture setup
When a reviewer asks which implementation behavior makes the scenario fail
Then the answer MUST be a behavior the assumption predicts is possible
```
Run today: A15 asks whether cardinality predicts *provability*; Scenario 3 measures only whether both products land on the same side of `B`, which any implementation that gates on cardinality alone passes trivially, and which no implementation that gates on shape can be made to reveal, since the fixture is sized to `B` and never to the proof representation. Fails.

**AT-11 — no normative MUST may name a field no consumer carries.** (catches C-11)
```gherkin
Given the withheld-claim clause requires the finding to name the refusing atom
When a reviewer resolves "the finding" to its producer's contract
Then that contract MUST carry the named field
Or the clause MUST be marked non-effective until it does
```
Run today: RDR 0006 `:363-367` has no atom-level field, the RDR says so in a block quote at `:977-984`, and the MUST stands unqualified. Fails.

**AT-12 — one predicate, one decision procedure.** (catches C-12)
```gherkin
Given the owned-tag clause at `:1075-1078` quantifies over "every reachable predecessor"
When a reviewer names the decision procedure
Then exactly one procedure MUST be named
And it MUST be the same procedure the consuming document uses
```
Run today: A12 concedes none is citable; the cluster gate flagged that RDR 0006 frames it as graph reachability while this RDR requires a syntactic decision. Two procedures, neither chosen, one normative MUST. Fails.

**AT-13 — the normative example must be normative.** (catches C-13)
```gherkin
Given the RDR owns a five-field declaration model
When a reviewer looks for one complete, non-disclaimed example
  showing all five fields with their wire spelling
Then exactly one MUST exist
And the repo's fixtures MUST match it
```
Run today: the sole example (`:1219-1231`) disclaims the `single_valued` spelling (`:1243-1245`), shows no `contains` guard against its `set` declaration, and `evidence/spikes/guard-fixture.toml` declares none of `optional`, `single_valued`, or `elements`. Fails.

**AT-14 — absorbed scope re-runs the lenses that absorbed scope invalidates.** (catches C-14)
```gherkin
Given the tag declaration model was absorbed at Stage 6, after grounding,
  3amigo, and critique had run on a document that did not contain it
When a reviewer checks the Alternatives section for the absorbed contract
Then the absorbed contract MUST carry its own alternatives analysis
  and its own prior-art pass, at the depth the operator grammar received
```
Run today: six alternatives for the operator grammar, zero for the type system; A11's only rejected alternative is a placement (`:378-380`). Fails.

**AT-15 — every stated premise must name an enforcer that exists.** (catches C-16)
```gherkin
Given the conformance clause at `:787-798` makes every lint claim conditional
  on view conformance
When a reviewer asks which code checks conformance
Then the RDR MUST name a clause in a peer document that requires the check
And that clause MUST exist
```
Run today: the RDR routes the check to "the kernel's view assembly" and to "RDR 0007's and RDR 0006's respectively" without citing a clause in either; `resolve.go::assemble` performs no conformance check. Fails.

**AT-16 — a spike predating a contract change is not evidence for the changed contract.** (catches C-17)
```gherkin
Given A1's harness "consumes no declared domain" and predates the declaration model
And the RDR's central contract changed at Stage 6
When a reviewer checks A1's Verified stamp
Then the harness MUST exercise at least one operator against a *declared* domain
Or A1's scope note MUST say the vocabulary claim is untested against
  the model the RDR now owns
```
Run today: A1's scope note (`:166-172`) does concede independence from the declaration model — and then A1 stays `Verified` while `contains` is the one operator with no declared-universe evidence and Phase 3 is required to "add at least one `contains` predicate over a declared set-valued tag before the full operator vocabulary is accepted" (`:1726-1728`). The vocabulary is accepted at lock on a condition the RDR schedules for after lock. Fails.
