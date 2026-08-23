Model: claude-opus-5[1m]

# Second-Pass Critique — RDR 0002, Transition Table As Reviewable Data

## 6. Findings Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | §Technical Design, "The normalized candidate row is the kernel row: it carries … the predicate atoms (key, operator token, literal, block)" + §Normative Contracts guard-predicate clause | The RDR normalizes `match`, `all`, and `unless` into **one** predicate set, but the kernel row it targets has **two** fields — `Row.Match []Tag` and `Row.Guard string`. No clause anywhere says how the unified atom set splits back into those two fields. The single most load-bearing step of the normalizer is unspecified. | Two implementers produce two different tables from one TOML file; one puts `iter.lt=3` into `Match` (where `TagSet.matches` does raw equality on the string `"3"` and silently never matches), the other into `Guard`. Flow authors see edges that "should" fire silently refuse `no_match`. | §1, §2, premortem, AT-1 |
| C-2 | §Normative Contracts, "Every rule … MUST bind exactly one outcome: its combined predicate set (local match block plus inherited contexts) MUST contain exactly one atom on `recognized`" | The clause scopes outcome-binding to *match + contexts*, but the normalizer merges `guard.all` and `guard.unless` into the same set and lifts from it. A `recognized` atom authored under `[rule.guard.unless.recognized]` is either (a) illegally lifted, minting a negative outcome binding, or (b) left in the set as a guard atom the kernel can never see, because the lifted-outcome row already filtered on `Row.Outcome`. The spike implements (a): `liftOutcome` scans the post-merge map with no block filter. | An author writes `unless.recognized.eq = "finalized"` intending "any outcome but finalized". The rule binds outcome `finalized` — the exact inverse of the authored intent — and fires on precisely the outcome it was written to exclude. | §1, §3, premortem, AT-2 |
| C-3 | §Normative Contracts, dump clause: "atoms sort by (key, block, operator token, literal)" + spike `mergeAtoms` keying `out[block+"\x00"+key+"\x00"+op]` | The atom set is a **map keyed by (block, key, operator)**, so two atoms with the same key and operator in the same block silently overwrite each other, and inheritance overwrite is order-dependent on `use = [...]`. The C-002 determinism fix sorts the *output* of a map whose *contents* already lost data. Determinism was proven; correctness was not. | A child context narrowing `status.eq` from a parent silently wins or loses depending on the `use` list order; an author who writes two `in` constraints on one key sees one silently vanish from the dump. No error, no lint finding — the row is simply wrong. | §1, §2, premortem, AT-3 |
| C-4 | §Trace table, `expand in` row: "the single-member case is **unwitnessed** (spike suffixes on operator, not expansion count)" | The RDR's own trace admits the C-002 expansion-suffix rule is contradicted by the only running code that implements it, and books this as a "needs-verification item for Stage 6" rather than a defect. The spike suffixes on `len(outcomes) > 1`, which is the *right* rule — but the RDR asserts the fixtures are "the canonical examples implementation tests must promote," and no fixture exercises single-member `in`. | An author refactors `eq = "x"` to `in = ["x"]` for consistency. If the implementer reads the trace note instead of the contract, row identity changes from `r` to `r#x`, every golden test churns, and every RDR 0006 lint finding that named the old id is orphaned. | §1, §2, AT-4 |
| C-5 | §Normative Contracts, "`[model]` MUST contain `id` and `version`. Version `1` is the only version this RDR accepts" | The C-001 version gate is a **refusal-only** gate with no forward story. There is no clause on what a version-2 file looks like, no rule that version is checked before *any* other decoding, and — critically — TOML decoding of the whole document happens before the version field is read, so a v2 file with a renamed section fails as `unknown schema field` (or malformed TOML) rather than `unsupported version`. The two categories collide. | A user on a newer tool version opens an older repo and gets `unknown schema field: contexts` instead of "this file needs tool ≥ X." They edit the file to satisfy the parser and corrupt it. | §1, §2, AT-5 |
| C-6 | §Critical Assumptions A8, "with `Input.Recognized == ""` *and* the empty string in the alphabet, a row with an empty `Outcome` matches … The recognized-totality clause below is the sole barrier" | A8 correctly identifies a kernel hole and then declares this RDR's load-time clause the "sole barrier." But the kernel takes a `Table` value from *any* caller — the normalizer is not the only producer, and `internal/resolve` validates nothing. A load-time check in a package the kernel does not import is not a barrier; it is a convention. | A test helper, a future fixture builder, or RDR 0006's lint constructing a `Table` literal reproduces the exact hole A8 documents, and the kernel emits a plan against a view with no `recognized` key. | §1, §3, premortem, AT-6 |
| C-7 | §Normative Contracts, `RequiresOwned` clause: "the sorted, duplicate-free set of tag keys named by the rule's write block and clear list" | `RequiresOwned` is derived from **writes only**, excluding guard-read keys ("does not add guard-read keys to it"). But the kernel gates on `RequiresOwned` *before* evaluating guards, and `evaluateGuard` returns `GuardUnevaluable` when the seam cannot decide. A guard reading an owned tag absent from the snapshot therefore yields `guard_unevaluable`, not `owned_state_unavailable` — the precise diagnosis inversion the kernel's own `gate` comment says it exists to prevent. | An author whose `iter` tag is missing from the artifact gets `guard_unevaluable: iter.lt=3` — a bug-shaped message pointing at the guard seam — instead of `owned_state_unavailable: iter`, which names the actual fix. | §1, §3, premortem, AT-7 |
| C-8 | §Round-Trip / Inverse Invariants, "with the source locator's optional line/column detail excluded from the comparison" | The invariant compares over a field list that **includes** the source locator, then excludes the only part of the locator that varies. What remains is `(model id, rule id)` — already in the identity tuple. The locator contributes nothing to the invariant, so a normalizer that drops locators entirely passes the round-trip test, and every downstream diagnostic loses its back-reference. | Lint reports `ambiguous overlap between rule-a and rule-b` with no file or line. On a 200-rule model the author cannot find either rule. | §2, premortem, AT-8 |
| C-9 | §Normative Contracts, next-state/writes clause: "Normalization MUST populate the next-state tags … and the writes … — the same rendered set" | Two distinct kernel fields (`NextTags`, `Writes`) are mandated to always carry **identical** content. The spike confirms it: every row dumps `next=[…]` byte-identical to `write=[…]`. A contract that forces two fields to be equal has not modeled a distinction — it has duplicated a field and deferred the real question (which writes are persisted vs. which are merely next-state) to whoever hits it first. | The accessor layer writes `<clear>` sentinels and next-state-only tags to the artifact because writes == next tags, clearing fields the author never intended to persist. | §1, §2, premortem, AT-9 |
| C-10 | §Normative Contracts, tag-key identity clause: "no case folding, trimming, or namespace rewriting" | Byte-exact key identity is specified for *keys* but nothing is said about **literal** normalization, while the spike's `renderLiteral` does exactly what the key rule forbids: it renders set literals as space-joined sorted members. A literal containing a space is therefore indistinguishable from a two-member set. | `in = ["needs work"]` and `in = ["needs", "work"]` normalize to the identical atom `recognized.in=needs work`. One of the author's two edges silently disappears. | §1, premortem, AT-10 |
| C-11 | §Minimum Viable Validation, "two sibling candidate rows binding the same outcome" + §Trace table "no witness possible in the current fixture" | The RDR requires a two-sibling fixture to make the gate testable, then ships fixtures where "every outcome binds exactly one transition row," and states the property is "closed by the MVV's two-sibling requirement." The requirement is an IOU written against fixtures the RDR simultaneously declares canonical. Nothing forces the fixture to change. | The implementer promotes the canonical fixtures verbatim, scenario 4's three tag-sets have no sibling rows to draw from, and the gate-ordering property — the reason this row shape exists — ships untested. | §1, §2, premortem, AT-11 |
| C-12 | §Normative Contracts, escape clause: "Row kind … is a **derived view property, never a row field**" + §Illustrative Code `draft-no-match-escape` | An escape row carries a full predicate set including its `recognized` binding, and the kernel's `escapeOrRefuse` filters escapes on `row.Outcome != in.Recognized`. So an escape row rescues `no_match` **only for the one outcome it binds** — a per-outcome escape, not a per-flow one. The RDR never says this. Its own fixture models `no_match` for `round-clean` alone, leaving three of four alphabet outcomes unrescued. | An author adds one `no_match` escape row believing the flow now handles unmatched input. Three quarters of outcomes still hard-refuse. The dump shows the escape row present, so the model "looks" complete. | §1, §3, premortem, AT-12 |
| C-13 | §Existing Infrastructure Audit, "Add stable parse/lint refusal codes later — these are the first of their kind" + A5 | A5 honestly establishes that no parse-validation code has ever travelled the CLI envelope, that `config-invalid` is a TODO, and that code names, `Group`, and exit mapping are "new design work at implementation." The RDR then defines ~17 stable data-level categories (the C-003 fix) with **no mapping to CLI codes and no owner** — RDR 0005 is "Pending." Stable categories with no stable surface are internal constants. | Every load failure exits with `command-error` and a prose message, because that is the only general-purpose code that exists. Scripting against table validation is impossible; the C-003 category taxonomy is invisible to the user it was written for. | §1, §2, AT-13 |
| C-14 | §Normative Contracts, dump clause: "`[dump]` settings MAY reorder the rendered columns; they MUST NOT omit a field" | `[dump]` is authored **inside the transition model file** and the spike's fixture carries an `order` list. Render settings living in the semantic source mean the same model file produces different bytes depending on a field that has nothing to do with transition semantics — and the round-trip invariant excludes dump text, so nothing catches it. | A reviewer diffing two branches sees the whole table churn because someone reordered `[dump].order`. Golden tests fail with zero semantic change. | §2, AT-14 |
| C-15 | §Normative Contracts, "an absent expansion suffix sorts before any present one" | The total-ordering argument (the C-002 fix) rests on rule ids being unique and suffixes being alphabet members. But the RDR's own Round-Trip section admits "the expansion suffix has no reserved separator, so a rule id containing it makes `(rule id, suffix)` unrecoverable." A rule id containing `#` collides at the *identity* level, not just the render level — and the duplicate-rule-id check compares rule ids, not identity tuples. | Rules `a#x` (no expansion) and `a` (expanding on `x`) both produce identity `a#x`. The dump shows two rows with one identity; lint's "name both rows" diagnostic names the same row twice. | §1, §2, premortem, AT-15 |

---

## 1. The Three Most Likely Ways Implementation Goes Wrong

### 1.1 The normalizer cannot produce a kernel row, because the RDR never says how its atom set becomes `Row.Match` and `Row.Guard`

**Root cause in the RDR.** This RDR's entire deliverable is a function from sparse TOML to `internal/resolve::Row`. It specifies, in exhaustive detail, an intermediate object — "the predicate atoms (key, operator token, literal, block)" — and then asserts that object *is* the kernel row: "The normalized candidate row is the kernel row." It is not. The kernel row on `main` has `Match []Tag` (a flat equality list consumed by `TagSet.matches`, which does `tv.value != w.Value` and nothing else) and `Guard string` (an opaque blob handed to a seam). The RDR's unified atom set has to be **split** across those two fields, and no clause in the document performs that split.

**The specific passage that enabled it.** §Technical Design: *"The normalized candidate row is the kernel row: it carries the source identity …, the single outcome the row responds to, the predicate atoms (key, operator token, literal, block), the next-state tags and the writes …"* — followed by the Normative Contract: *"Normalization MUST combine both into one candidate-row predicate set before ambiguity checks."* Combine, yes. Split back out, never mentioned. The Prerequisites section notices the gap ("`Row.Guard` is a `string`, `Row.Match` is a flat `[]Tag` equality list … there is no `Atom` type") and then *reclassifies it as someone else's sequencing problem*: "Phases 2 and 3 therefore sequence behind the reshape." That is not a resolution. RDR 0007's reshape defines an atom-shaped row; it does not define how *this* RDR's `all`/`unless`/`match` blocks map onto whatever fields survive. The RDR has deferred its own core function to a document it cites but does not quote on this point.

**Symptom the user will see.** The atom `iter.lt = 3` has to go somewhere. If the implementer routes non-`eq` atoms to `Match`, `TagSet.matches` compares the view's `iter` value against the literal string `"3"` with `!=` and the row never matches — every prelock continuation silently refuses `no_match`, and the dump shows the row present and correct. If the implementer routes *all* atoms to `Guard`, then `Row.Match` is empty, every row for an outcome becomes a candidate, and the flow author gets `ambiguous_match` naming four rules that share nothing but an outcome. Both implementations pass every test the RDR names, because every MVV assertion is stated over the *normalized value* — the intermediate object — and never over the kernel row that the kernel actually consumes. The RDR's Testing Strategy scenario 4 does run `Resolve`, but only over "one matching ordinary tag-set," a single positive case that either wiring can be made to pass by choosing the fixture.

### 1.2 The predicate set is a map, so constraints silently overwrite each other, and the determinism fix certifies the loss

**Root cause in the RDR.** The C-002 fix established deterministic *ordering* of the expanded table. It did not establish deterministic *content*. The RDR describes the predicate set as a set — "Normalization MUST combine both into one candidate-row predicate set", "each atom in that set MUST retain the key, operator token, literal, and the block" — and specifies a sort over it. It never states the **cardinality rule**: how many atoms may share a key, or a (key, block), or a (key, block, operator). The only running implementation answers with a map keyed on `block + key + operator`, which means last-write-wins with no diagnostic. Contexts merge into that map in `use` list order, then `match`, then `all`, then `unless`.

**The specific passage that enabled it.** §Normative Contracts: *"Shared contexts MAY inherit from other contexts, but inheritance MUST normalize to an explicit predicate set before lint or resolution."* "Normalize to an explicit predicate set" is the entire specification of inheritance semantics. It does not say whether a child's `status.eq = "Final"` *overrides*, *conjoins with*, or *conflicts with* a parent's `status.eq = "Draft"`. Conjunction would make the row dead (an RDR 0006 finding). Override is a real, defensible statechart semantic. Conflict-is-an-error is the third option. The RDR picks none, and the fixture never exercises the case — `draft` → `prelock` → `large-prelock` each narrow on a *fresh* key, so the overwrite path is never reached in evidence. §A6 stamps this Verified on exactly that non-witnessing fixture: *"`rdr-fixture.toml` encodes `draft -> prelock -> large-prelock` inherited contexts … preserving ambiguity visibility in the expanded row."* Ambiguity that the fixture cannot produce is not visible; it is absent.

**Symptom the user will see.** An author builds a `final` context inheriting `draft` and overriding `status.eq`. Depending on which side the implementer chose, the rule either matches only `Draft` (parent wins), only `Final` (child wins), or nothing at all. The expanded table dump shows **one** `status` atom, deterministically sorted, byte-identical across runs — the C-002 fix working perfectly on corrupted data. The author reads the dump, sees a single plausible atom, and concludes the model is right. The determinism guarantee has been converted into a credibility guarantee for a wrong answer.

### 1.3 Outcome lifting reads from the wrong scope, so `unless` on `recognized` inverts the author's intent

**Root cause in the RDR.** The C-001-era outcome-binding contract defines the binding scope as *"its combined predicate set (local match block plus inherited contexts)"*. The guard-predicate contract, three clauses earlier, defines the combined predicate set as `match` ∪ `all` ∪ `unless`. These two definitions of "combined predicate set" are different sets, and the RDR uses the same phrase for both. The outcome-binding clause then says *"Normalization lifts that atom out of the predicate set"* — out of which one?

**The specific passage that enabled it.** §Normative Contracts, outcome-binding clause: *"MUST contain exactly one atom on `recognized`, using `eq` or `in`, whose literal(s) are members of the `outcomes` alphabet."* There is no block restriction on this clause. It does not say the atom must be authored in `match`. The only running implementation, `liftOutcome`, iterates the post-merge map — which includes `unless` atoms — and lifts whatever it finds on key `recognized`, discarding the block. The `Block` field it was so careful to retain is dropped precisely at the one decision where it is load-bearing.

**Symptom the user will see.** An author writes what reads as an obvious negative guard:

```toml
[rule.guard.unless.recognized]
eq = "finalized"
```

meaning "this rule applies to anything but a finalized outcome." Normalization lifts `finalized` into `Row.Outcome`. The rule now fires **only** on `finalized` — the single outcome it was written to exclude. There is no load failure (the atom count is exactly one, the literal is in the alphabet), no lint failure (the row is reachable and unambiguous), and the dump reads `outcome=finalized`, which an author scanning for the *absence* of finalized will not flag. This is the worst class of defect the RDR can ship: a semantically inverted edge that every layer of the specified validation stack certifies as correct.

---

## 2. The One Section That Will Be Rewritten Within Six Weeks

**§Normative Contracts, the expanded-table dump clause** — the C-002 fix — will be rewritten, and it will be rewritten because it over-specified a *rendering* while under-specifying the *value* it renders.

The clause is 30 lines long. It mandates that the dump carry every field of the normalized row, fixes a three-field sort key, fixes intra-row sort keys for five sub-collections, forbids map iteration, and cites Go's map-randomization and `encoding/json/v2` key-sorting behavior. It is the most carefully engineered passage in the document. It is also the one that has to change first, for four independent reasons that will arrive in roughly this order:

**Week 1–2: the `[dump]` conflict.** "`[dump]` settings MAY reorder the rendered columns; they MUST NOT omit a field" puts a *render preference* in the *semantic source file*. The first time two people set `order` differently on the same model, every golden test in the repo diffs. The fix is to move `[dump]` out of the model file, which changes the schema layout — and the schema layout is itself a Normative Contract ("The source schema MUST use the Resolve spike field layout: root `outcomes`, `[model]`, … and `[dump]`"). Two contracts have to move together.

**Week 2–3: the field list is wrong.** The clause enumerates the fields as "row identity, source locator, outcome, predicate atoms (with block), next-state tags, writes including `<clear>` entries, required-owned keys, and escape failure classes." That is the *intermediate* object. It has no `Row.Guard`. The moment §1.1's split is resolved, the dump either has to render the guard string (a field the clause does not list, so rendering it violates "MUST NOT omit a field"'s sibling expectation of a fixed list) or has to render the atoms it no longer holds. The field list was frozen against an object that does not survive contact with the kernel.

**Week 3–4: `next` and `writes` collapse.** The clause requires both to be rendered. C-9 establishes they are contractually identical. The first reviewer to read a dump will ask why every row prints the same list twice, and the answer — that the RDR mandated it — will not survive the asking. Either the fields diverge (and the contract's "the same rendered set" clause dies) or one column goes (and "MUST NOT omit a field" dies).

**Week 4–6: the total-ordering proof fails.** The clause's totality argument is: *"rule ids are unique within a model, and an expansion suffix is a member of the `outcomes` alphabet … so no two distinct rows compare equal and no positional tiebreak is needed."* C-15 breaks it: the RDR's own Round-Trip section concedes the suffix has no reserved separator. A rule id containing the separator makes two distinct rows share an identity tuple, which means the comparator returns 0 for distinct rows, which means the sort is not total, which means the ordering is implementation-defined — the exact failure the clause was written to prevent. Fixing it requires either reserving a character in rule ids (a new load-time category, changing the C-003 taxonomy) or making identity a structured tuple that the dump cannot render as a flat string (changing the render contract again).

Every one of these four is a change to the same clause, and three of them cascade into other Normative Contracts. The clause is not wrong in its intent — deterministic ordering is genuinely load-bearing. It is wrong in having been specified to the byte before the object it orders was pinned down.

---

## 3. The One Assumption That Will Not Survive First Contact With a Real User

**A6: "Shared contexts and positive/negative guard lists are sufficient to keep the RDR model sparse without hiding ambiguity."** Status: Verified. Method: Spike.

It will not survive, and the reason is visible in its own evidence line.

The evidence is: *"`rdr-fixture.toml` encodes `draft -> prelock -> large-prelock` inherited contexts plus `all.iter.lt = 3` and `unless.profile.eq = "small"` guards; `output.txt` shows the normalized row with inherited status/stage/profile predicates and combined `all`/`unless` predicates, preserving ambiguity visibility in the expanded row."*

Read what that fixture actually contains. `draft` sets `status`. `prelock` inherits `draft` and adds `stage`. `large-prelock` inherits `prelock` and adds `profile`. Three levels of inheritance, and **not one key is ever touched twice**. The chain is a pure union over disjoint keys. It is the single easiest inheritance case that can be constructed, and it is the only one the spike ran.

The assumption's claim is about *sufficiency without hiding ambiguity*. The fixture demonstrates sufficiency on a case where ambiguity is structurally impossible. The moment a real author writes a second flow — and the RDR's whole premise is that the RDR flow *and* the kata flow *and* future flows all live in this format — they will do the obvious thing: define a base context and specialize it. `[context.final] inherits = "draft"` with `status.eq = "Final"`. That is the first thing anyone does with an inheritance feature. It is the reason inheritance exists.

And at that moment the RDR has no answer. Not a wrong answer — *no* answer. §1.2 walks through the three plausible semantics and shows the RDR picks none. The spike picks last-write-wins by accident of using a Go map. The dump renders one atom and looks clean. The determinism proof holds. Every validation category in the C-003 taxonomy passes. The author gets a rule that matches `Draft`, or `Final`, or neither, and has no artifact anywhere in the system that tells them which.

The assumption compounds with the sparse-authoring premise that justifies the entire RDR. §Decision Rationale rejects Alternative 2 (fully expanded rows) because *"changing one shared condition requires finding every copied row."* The whole document is built on contexts being the mechanism that makes sparse authoring safe. If context override semantics are undefined, then the mechanism that makes sparseness safe is the mechanism that makes it dangerous, and Alternative 2 — repetitive but unambiguous — is retroactively the better choice for the exact dimension the RDR used to reject it.

A6 should be `Unverified`. Its spike proves A3 (sparse encoding is possible), not A6 (sparse encoding does not hide ambiguity). The RDR's own Method vocabulary forbids this: *"The cited proof must also support the specific claim, not an adjacent one: confirming a neighboring fact and stamping the assumption `Verified` is not verification."* A6 does exactly that, and the Finalization Gate's Assumption Verification item — which asks the reviewer to confirm no `Verified` stamp "proves only an adjacent claim" — is the item that will be answered wrongly.

---

## 4. Premortem

*Written from six months after ship.*

We shipped the normalizer in three weeks. The spike fixtures promoted cleanly, every MVV item had a failing control, and the C-003 category taxonomy gave us seventeen distinct load-time refusals with one mutated fixture apiece. The gate.md responses were the most thorough in the project. We were confident.

**The first failure was `normalize()` producing rows the kernel could not evaluate.** The RDR said "the normalized candidate row is the kernel row," so `normalize()` returned our `Row` type with an `Atoms []Atom` field. Then Phase 4 — Resolver Handshake — had to hand those to `resolve.Resolve`, which wanted `Row.Match []Tag` and `Row.Guard string`. RDR 0007's reshape had not landed and would not land for another two months. We wrote `toKernelRow()` in an afternoon. It put every `@match` atom into `Match` as a `Tag{Key, Value: literal}` and joined every `@all` and `@unless` atom into a `Guard` string with semicolons. Nobody reviewed the semicolon choice; it was three lines inside a function the RDR did not name.

That put `profile.in = ["large","foundational"]` — a `@match` atom — into `Match` as `Tag{Key: "profile", Value: "foundational large"}`. `TagSet.matches` compared it against the view's actual `profile` value of `"large"` with `!=`. Every RDR in the foundational or large profile stopped transitioning. `continue-prelock` refused `no_match` for six weeks before anyone connected it. The expanded table dump showed `profile.in=foundational large@match`, exactly as the contract specified, byte-identical across runs. We had built a perfectly deterministic renderer for a row the kernel silently discarded. **C-1.**

**The second failure was Priya's `final` context.** Priya was writing the kata-close flow — the second real model in the format, the whole reason the format is sparse. She wrote:

```toml
[context.final]
inherits = "draft"
[context.final.match.status]
eq = "Final"
```

`mergeContext` walked `draft` first, wrote `match\x00status\x00eq → Draft` into the atom map, then walked `final`'s own match block and wrote the same map key with `Final`. Child won. That was fine — it was what Priya wanted. Three weeks later Dan added `use = ["final", "draft"]` to a rule to pick up a shared predicate from `draft`, and the merge order flipped: `final` first, then `draft` overwrote `status` back to `Draft`. Dan's rule now matched Draft RDRs. There was no error, no lint finding, no diff in the dump beyond one atom's literal changing from `Final` to `Draft` inside a 400-line generated file. Nobody read that line. Dan's rule fired on drafts for a month, and the `terminal-archive` edge it shadowed stopped firing on finals — which read as "archiving is broken," and sent two people into `internal/resolve` looking for a kernel bug that was not there. **C-3.**

**The third failure was `liftOutcome` and the word "combined."** Sam was writing a catch-all rule and wanted "any outcome except `finalized`." The natural spelling was `[rule.guard.unless.recognized] eq = "finalized"`. `liftOutcome` scanned the merged map, found exactly one atom on key `recognized`, checked the literal against the alphabet — it was in the alphabet — lifted it into `Row.Outcome`, and deleted it from the predicate set. The block was never consulted; `Atom.Block` was retained everywhere except at the one decision where it decided meaning. Sam's catch-all fired on `finalized` and only on `finalized`. It ran a rewind. Two RDRs got rewound out of Final state before we noticed, and the audit trail said `rule=sam-catchall outcome=finalized`, which is exactly what the dump said, which is exactly what the contract said should happen. We had to explain to Sam that the document said "local match block plus inherited contexts" and the code said "everything," and that both were readings of the same clause. **C-2.**

**Then the diagnostics stopped naming anything.** Our round-trip test compared normalized values "with the source locator's optional line/column detail excluded." We had made the locator `rule.Source` — a hand-authored string like `"rdr:prelock"` — because the contract only required "at least the model id and rule id" and line/column was explicitly optional. Then somebody made it optional in practice: rules without a `source` key got `""`. The round-trip test passed (locator content was excluded from comparison in all the ways that varied). When RDR 0006's lint reported `ambiguous overlap between "continue-prelock" and "continue-prelock-large"` on Priya's 200-rule kata model, it printed no file and no line. Priya grepped. **C-8.**

**The escape rows were the quiet one.** We had `draft-no-match-escape` in the fixture, so the flow "had" a no-match escape. It bound `recognized.eq = "round-clean"`. `escapeOrRefuse` filters escape candidates on `row.Outcome != in.Recognized` before counting — so the escape rescued `no_match` for `round-clean` and for nothing else. `reconcile-block`, `verdict-flapping`, and `finalized` all hard-refused. The dump showed `kind=escape escape=[no_match]`, and every reviewer who scanned it read "this flow handles no-match." It took an incident to learn that escape coverage is per-outcome and that our four-outcome flow needed four escape rows. The RDR never said so; the kernel's filter said so, and we had read the RDR. **C-12.**

**And the sibling gate was never tested at all.** The MVV required "two sibling candidate rows binding the same outcome" and said the property was "closed by the MVV's two-sibling requirement." But the RDR also said the spike fixtures were "the canonical examples implementation tests must promote," and those fixtures have exactly one transition row per outcome. We promoted them. Scenario 4 needed three tag-sets "all three drawn from the fixture's two same-outcome sibling rows," and there were no sibling rows, so we wrote the test against single-candidate cases and it passed. The property that an unevaluable survivor refuses even when a decidable sibling and an escape row both exist — the property this RDR's entire row shape exists to support, the one the trace table flagged as having "no witness possible in the current fixture" — shipped with zero coverage. We found out when a missing accessor made one sibling unevaluable in production and the resolver took the decidable sibling's plan, because our `toKernelRow` had put the guard into `Match` where an undecidable predicate is just a non-match. **C-11, and C-1 again.**

The document was not careless. It was precise about everything downstream of a decision it never made.

---

## 5. Acceptance Tests That Would Have Caught Each Failure at RDR-Review Time

These are review-time tests: each is answerable by reading the RDR alone, and each fails on the current draft.

**AT-1 — the normalizer's output is a kernel row (C-1)**
```gherkin
Given the RDR claims "the normalized candidate row is the kernel row"
When a reviewer takes each field of internal/resolve::Row on main
  (RuleID, SourceLocator, Outcome, Match, RequiresOwned, Guard, NextTags, Writes, Escape)
Then each field is named by a Normative Contract that states what normalization puts in it
And in particular there is a clause stating which authored blocks (match / all / unless)
  produce Row.Match and which produce Row.Guard
```
FAILS: `Row.Guard` is named nowhere outside a Prerequisites parenthetical; no clause splits the atom set.

**AT-2 — outcome binding names its block (C-2)**
```gherkin
Given the outcome-binding contract says the recognized atom is in "the combined predicate set"
When a reviewer asks whether an atom authored under [rule.guard.unless.recognized] can be lifted
Then the RDR gives one unambiguous answer
```
FAILS: "combined predicate set" is defined once as match+contexts and once as match+all+unless; both definitions are Normative Contracts. Fix: the binding atom MUST be authored in a `match` block; a `recognized` atom in `all` or `unless` is a load failure in a named category.

**AT-3 — atom set cardinality and context override (C-3)**
```gherkin
Given two contexts, parent with status.eq="Draft" and child inheriting it with status.eq="Final"
When a rule uses the child
Then the RDR states whether the normalized row carries one atom or two
And if one, which one, and whether use-list order can change the answer
And if two, whether the row is dead (RDR 0006) or a load failure here
```
FAILS: no clause addresses same-key collision in any block, in inheritance, or across `use` entries. The A6 fixture cannot produce the case.

**AT-4 — single-member `in` (C-4)**
```gherkin
Given a rule authored [rule.match.recognized] in = ["x"]
When it is normalized
Then the row identity is exactly the identity of the same rule authored eq = "x"
And a fixture in the canonical evidence set witnesses this
```
FAILS: the contract asserts it; the trace table records the spike does the opposite discrimination and books it as "unwitnessed." An assertion the only implementation contradicts is not a contract.

**AT-5 — version gate ordering (C-5)**
```gherkin
Given a file with [model].version = 2 and a section name unknown to version 1
When it is loaded
Then it is refused with "unsupported version" and not "unknown schema field"
And the RDR states that the version field is read before any other schema validation
```
FAILS: the version clause says "MUST be refused before normalization," which is satisfied by refusing after full strict decoding. The C-001 gate does not order itself against the strict-decode rule the MVV mandates.

**AT-6 — the empty-outcome hole has an enforceable barrier (C-6)**
```gherkin
Given A8 documents that a kernel Table with "" in Outcomes and a row with Outcome ""
  emits a plan against a view carrying no recognized key
When a reviewer asks what prevents a non-normalizer producer of resolve.Table from doing this
Then the RDR names an enforcement point inside internal/resolve, or explicitly accepts the residual
```
FAILS: A8 calls this RDR's load clause "the sole barrier" for a value type the kernel accepts from any caller. A load-time check in a package the kernel does not import cannot be a barrier.

**AT-7 — diagnosis precedence for guard-read owned tags (C-7)**
```gherkin
Given a rule whose guard reads owned tag "iter" and whose write block does not name "iter"
And a resolve where "iter" is absent from the owned snapshot
When the kernel gates
Then the RDR states which refusal the author sees: owned_state_unavailable or guard_unevaluable
```
FAILS: the `RequiresOwned` contract excludes guard-read keys and defers meaning to RDR 0007, but the kernel's `gate` orders `missingOwned` before undecidability specifically to give the precise diagnosis. Excluding guard reads defeats that ordering, and the RDR does not say so.

**AT-8 — the round-trip invariant constrains the locator (C-8)**
```gherkin
Given the round-trip invariant compares over a field list including the source locator
  but excluding its line/column detail
When a candidate normalizer emits an empty locator for every row
Then the round-trip test fails
```
FAILS: it passes. The compared residue is `(model id, rule id)`, already in the identity tuple. Either drop the locator from the invariant and assert its presence separately, or make non-empty locators a load-time obligation.

**AT-9 — NextTags and Writes are distinguishable (C-9)**
```gherkin
Given the kernel carries NextTags and Writes as two fields
When a reviewer looks for one normalized row in the RDR or its evidence where they differ
Then such a row exists
```
FAILS: the contract mandates "the same rendered set," and all seven spike rows print `next` == `write`. A contract forcing two fields equal has not modeled a distinction.

**AT-10 — literal identity is as strict as key identity (C-10)**
```gherkin
Given the tag-key contract mandates exact byte equality with no rewriting
When a reviewer asks the same question of literals
Then the RDR states the literal normalization rule, including set literals
And a member containing the set separator is either escaped or a load failure
```
FAILS: literals are unaddressed; the spike space-joins sorted set members, so `["needs work"]` and `["needs","work"]` are one literal.

**AT-11 — the sibling fixture obligation is binding (C-11)**
```gherkin
Given the MVV requires two sibling candidate rows binding the same outcome
And the RDR names the spike fixtures as "the canonical examples implementation tests must promote"
When a reviewer checks whether the canonical fixtures satisfy the requirement
Then they do
```
FAILS: the trace table states no outcome in the fixture binds more than one transition row, and records "no witness possible in the current fixture." Two clauses of the same document require contradictory things of the same file.

**AT-12 — escape coverage scope (C-12)**
```gherkin
Given the kernel's escapeOrRefuse filters escape candidates on row.Outcome != in.Recognized
When a reviewer asks how many escape rows a flow with N outcomes needs to rescue no_match for all of them
Then the RDR answers N, and says so where escape rows are specified
```
FAILS: the escape contract never mentions that escape rescue is per-outcome. The illustrative fixture models one escape for a four-outcome alphabet and presents it as complete.

**AT-13 — the C-003 category taxonomy reaches the user (C-13)**
```gherkin
Given ~17 stable data-level categories are mandated "before CLI mapping"
When a reviewer asks which RDR maps them to CLI codes, Group, and exit status
Then a Final RDR owns that mapping, or this RDR defines it
```
FAILS: A5 establishes no parse-validation code has ever travelled the envelope and that names/Group/exit are "new design work at implementation"; RDR 0005 is Pending. Stability was specified for a surface with no owner.

**AT-14 — render settings are outside the semantic source (C-14)**
```gherkin
Given [dump] is authored inside the transition model file
When two authors set [dump].order differently for the same model
Then the RDR states whether that is a semantic change, and what protects golden tests from it
```
FAILS: the round-trip invariant is stated over the normalized value and explicitly excludes dump text, so nothing in the document constrains the churn `[dump]` can cause.

**AT-15 — the identity tuple is injective (C-15)**
```gherkin
Given rows are identified by (model id, rule id, expansion suffix)
And the Round-Trip section states the suffix has no reserved separator
When rule "a#x" (unexpanded) and rule "a" (expanding on outcome "x") are in one model
Then the RDR states that their identities are distinct, or names the load-time category that rejects one
```
FAILS: the totality argument assumes rule-id uniqueness suffices; the duplicate-rule-id check compares rule ids, not identity tuples, so the comparator returns 0 for two distinct rows and the "total ordering" is not total.
