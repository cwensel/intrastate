Model: claude-opus-5[1m]

# Critique — RDR 0002, Transition Table As Reviewable Data (iteration 3)

## 6. Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | §Normative Contracts, dump clause: "The identifiers are the **field names enumerated in the sentence above**, snake_cased … (`row_identity`, `source_locator`, … `write` …) — that closed list is the column vocabulary" vs. `evidence/spikes/iter-2/rdr-fixture.toml` `[dump] order = ["identity", "source", "kind", "outcome", "atoms", "next", "writes", …]` | Both normative fixtures author three column identifiers (`identity`, `source`, `writes`) outside the closed vocabulary, so by this RDR's own `malformed dump declaration` rule both fixtures MUST refuse at load. A conforming implementation refuses the exact files Testing Strategy 1/2/3/5 name as its inputs. | First `go test` after the dump validator lands: every fixture-backed test fails `malformed dump declaration: unknown column "identity"`. Holding a §D7-blessed fixture and a Stage-4-approved SHA, the implementer deletes the validator — retiring the category. | §1, premortem, AT-1 |
| C-2 | §Normative Contracts, dump clause: "An `order` naming an unknown identifier, repeating one, or omitting one is a `malformed dump declaration` refused at load" vs. Testing Strategy 3: "the spike decodes `[dump]` and never reads it" | `[dump]` is the only sub-schema with a full normative validation rule and zero executed evidence — no negative fixture, no positive load, no implementation. Verified: mutating `order` to `["identity","NOT_A_COLUMN"]` loads clean and dumps nine rows. | A reviewer edits `order` to drop `escape`; the dump renders nine columns instead of ten; "a dump missing any field is not an expanded table dump" is violated with no diagnostic. Golden diffs shrink and nobody knows why. | §1, §2, premortem, AT-1 |
| C-3 | §Normative Contracts, operator clause: "A token under `[rule.guard.all.<tag>]` or `[rule.guard.unless.<tag>]` outside the set is a `malformed predicate atom` (`unknown operator`) refused at load — the guard-side mirror of the match-block restriction above, **which without this was absent**" | The guard-block validator in the sole spike checks only *key declared* and *not `recognized`* (`main.go:529-537`). No operator-membership check. Verified: `[rule.guard.unless.profile] frobnicate = "small"` loads clean and emits `profile.frobnicate=small@unless`. The RDR notices the rule "was absent," adds the prose, and adds no fixture and no implementation. | The author types `greater_than` for `gt`. Table loads, lint passes, and RDR 0007's seam answers `uncomparable` naming the *key*, not the typo. RDR 0007 predicts the response: "it pushes authors toward sentinel-stamping to make the kernel talk, which is the named anti-pattern." | §1, §3, premortem, AT-2 |
| C-4 | §Normative Contracts, reserved-value clause: "The string `<clear>` MUST be refused at load wherever a tag value is authored — a write-block value, an `[initial]` assignment, **or a predicate literal** — as a `reserved tag value`" | Implemented only inside `checkMatchBlock` (`main.go:563-565`). Verified: `[rule.guard.unless.profile] eq = "<clear>"` loads clean and emits `profile.eq=<clear>@unless`. The clause's own stated reason — "a silent dead rule is the more dangerous failure" — is exactly what the evidence produces. | An unsatisfiable guard atom triggers RDR 0007's resolution-scope veto ("a decided-GuardTrue sibling MUST NOT be selected while an unevaluable candidate exists"). Every resolve touching that outcome refuses non-escapably. The flow deadlocks. | §1, premortem, AT-2 |
| C-5 | §Normative Contracts, tag-type clause: "`malformed predicate atom` (a literal outside the tag's declared domain, decided after declarations load and before rows are yielded)" | Domain membership is enforced in match blocks only. Verified: `[rule.guard.unless.profile] eq = "NOT_A_PROFILE"` — against `domain = ["small","mid","large","foundational"]` — loads clean. The clause is unqualified as to block; the evidence enforces it in one of three. | An `unless` guard misspells an enum member. The negation never fires, so a rule the author explicitly carved out now applies. The artifact answers "this edge is legal" for an edge the author forbade. No refusal, no lint. | §1, §3, premortem, AT-2 |
| C-6 | §Normative Contracts, next/writes clause: "**RDR 0009 A4 fixes that these are distinct fields** and that only the writes reach the accessor layer" and "They stay distinct fields rather than collapsing to one because **RDR 0009 A4 fixes them as distinct**" | Fabricated peer citation. RDR 0009 A4 (Final) settles the **opposite**: "Status: Verified — settled **closed**: the predicate does **not** widen; it stays `Writes`-only," and its normative fence reads "The predicate is Writes-only and does NOT extend to NextTags." A4's holding is that **no rule about `NextTags` on escape rows is needed**. RDR 0002 cites it as authority for a MUST requiring `NextTags` empty. | The implementer follows the citation to RDR 0009, finds it says the predicate does not widen, and cannot tell whether the empty-`NextTags` MUST is real. Whichever way they resolve it, one of two locked documents is being disobeyed. | §1, §2, premortem, AT-3 |
| C-7 | §Normative Contracts, reserved-value clause: "**RDR 0004 removes the key on that write and asserts absence on read-back**" | Fabricated peer citation. The word "clear" appears **zero times in RDR 0004's body** — only in its Status line and Refinement Context. RDR 0004's Fidelity Table asserts the opposite with no exemption: "Every planned owned-tag key/value **present** in the re-read equals the transition plan's expected value … Lossy exemptions: **none**." RDR 0004's own Refinement Context names this defect as unrepaired and schedules it for a future stage. | Every clearing rule read-backs as a write failure under RDR 0004 as written. `reconcile-rewind` clears `prelock_lens`; the accessor re-reads, finds the key absent, and reports the write failed. Rewind is unusable. | §1, premortem, AT-3 |
| C-8 | §Normative Contracts, guard clause: "the block it was authored in — a three-valued domain, `match`, `all`, or `unless` — the atom shape JDR 0001 §D1 fixes and **RDR 0007 spells**" | Fabricated peer citation. RDR 0007 is Final and spells a **two-valued** domain: "a slice of parsed atoms — key, operator token, literal, block ∈ {all, unless} — the shape JDR 0001 §D1 fixes." RDR 0003 quotes the same two-valued domain twice. The widening is legitimate but its authority is JDR 0001 §D6, not RDR 0007. | The implementer reads RDR 0007 to build the atom type, gets a two-valued block, and either mints a two-valued enum that cannot carry `match` (breaking block retention) or ignores the Final peer. | §1, premortem, AT-3 |
| C-9 | §Normative Contracts, escape clause: "Row kind (`transition` / `escape`) is a **derived view property, never a row field** … normalization introduces no field to carry it" | Contradicts RDR 0006 (Final), whose **minimum input contract** lists row kind as a required input distinct from the class list: "writes, clears, and — for escape rows — **row kind and** the declared list of failure classes the row rescues (RDR 0002)." RDR 0009 (Final) separately quotes RDR 0002 verbatim as rendering "a candidate row **with row kind `escape`**" — a sentence RDR 0002 has since deleted and inverted, and on which RDR 0009 A4's Verified consistency note rests. | Lint cannot be built to its own input contract. The implementer either adds the field RDR 0002 forbids, or ships a lint that must re-derive kind — and RDR 0009's Verified evidence is grounded on text that no longer exists. | §1, §2, premortem, AT-3, AT-4 |
| C-10 | §Normative Contracts, root/stop-set clause: "root `terminal` is a list of context ids, each of which MUST resolve (`unknown context`)" | Contradicts RDR 0006 (Final) A6's explicit shape prohibition: "both declarations are **tag predicates, not state names** … what this RDR requires is that **neither is a bare identifier a row references by name**, since invariant 1 checks them as tag keys and values and invariant 2 evaluates a terminal as a predicate over a node." RDR 0002 ships exactly a bare identifier referenced by name and never states the dereference to a predicate set. | RDR 0006's invariant 2 splits nodes on "terminal-participating keys" — a notion that requires a terminal to *be* key/value atoms. Handed a context id, lint cannot compute reachability or dead-ends, and the stop set is unenforceable. | §1, §2, premortem, AT-4 |
| C-11 | §Conditional Mini-Checks `disposition` table: "\| expansion counts per rule \| lint \| RDR 0006 \| reported, non-blocking; no threshold here \|" and "ambiguous expansion" routed to RDR 0006 | Contradicts RDR 0006 (Final), whose advisory tier is normatively **closed at four members**: "The advisory tier is closed at `graph-coverage-closed-by-escape`, `graph-redundant-row`, `graph-unreachable-rule`, and `graph-vacuous-atom`." RDR 0002 invents a fifth advisory output and a `graph-ambiguous-expansion` code RDR 0006 does not mint. Additionally, RDR 0006 contains **zero** occurrences of "expansion" in any relevant sense — it does not know rows expand. | The Risks section's mitigation for "sparse contexts hide an accidental Cartesian product" is "lint must report expansion counts per rule." That report is forbidden by the lint's own closed tier. The stated mitigation for the format's headline risk cannot be built. | §1, §2, premortem, AT-4 |
| C-12 | §Normative Contracts, literals clause: "Wherever atoms are compared, keyed, deduplicated, or sorted … the literal field MUST be compared as the member sequence, element by element. An implementation MUST NOT derive that key, or any other identity, by joining members into a string." | Contradicts RDR 0006 (Final), whose mandatory finding payload carries `Literal` as a flat Go **string**: "`clierr.Finding` carries those as declared string-typed fields, not enums: `Key`, `Operator`, `Literal`, `Block`, and `Class` are `string`," bound by "A finding attributed to one guard atom MUST carry that atom's `Key`, `Operator`, `Literal`, and `Block`." The banned joined rendering is reintroduced at the locked diagnostic boundary. | The lint finding for a set-literal atom prints `["a,b","c"]` and `["a","b,c"]` identically. The reviewer reading the diagnostic cannot tell which atom is at fault — the exact ambiguity the literals clause spends two fences banning. | §1, §2, premortem, AT-4 |
| C-13 | §Conditional Mini-Checks `trace`, rows "merge atoms" and "write values": two **CONTRADICTION** verdicts; A13 `Status: Pending` | The RDR locks with two self-declared unresolved contradictions between its normative clauses and its only executable evidence, on the arm it describes as "changes match semantics … inverts the user outcome this RDR exists to deliver" and that "trips **zero** load categories." | The author writes `in = ["a,b", "c"]` in one context and `in = ["a", "b,c"]` in another. One constraint vanishes in the merge. The table answers "this edge is legal" for an edge the author did not write. No load error, no lint finding, no dump difference. | §1, §3, premortem, AT-5 |
| C-14 | §Normative Contracts, `RequiresOwned` clause: "the sorted, duplicate-free set of tag keys named by the rule's write block and clear list" + RDR 0007 (Final): "predicate PLACEMENT decides whether an absent key refuses non-escapably or escapes. The only control is authoring guidance … which Phase 4 hands to RDR 0003" | The most consequential authoring decision in the format — match block vs. guard block — is delegated by a Final peer to RDR 0003, and RDR 0002, which owns the authoring surface, offers no guidance, no lint, no diagnostic. Its routing clause calls the match block "select on," which is precisely how an author reads `stage`. | The accessor fails to read `stage`. `TagSet.matches` returns false, every candidate drops out, `escapeOrRefuse` runs, and `draft-no-match-escape` returns a **plan**. Missing owned state is laundered into a successful transition — what JDR 0001 §D2 exists to forbid. | §1, §3, premortem, AT-6 |
| C-15 | §Technical Design: "(The iter-2 spike populates `Locator` from `rule.Source` directly … that is a spike defect against this clause, not a permitted reading.)" vs. §Normative Contracts schema clause: "each `[[rule]]` an optional `source` — the authored provenance string the locator derives from" | The RDR names a defect in the artifact it simultaneously declares normative, and gives the locator two incompatible origins one clause apart ("composed from `(model id, rule id)`" vs. "derives from" the authored `source`). `output.txt` — SHA-pinned as "the expected value" — carries two kata rows whose `source=kata:review` collide. | Two rows are implicated in an `ambiguous_match`. `compareRefs` orders on `(RuleID, SourceLocator)` and RDR 0005 prints the locator. Both conflicting rows report `kata:review`. The user is told two rows conflict and handed one identical pointer for both. | §1, premortem, AT-7 |
| C-16 | §Testing Strategy 2: "The normative fixture for this scenario is `evidence/spikes/iter-2/output.txt` (approved Stage 4, A1; SHA-256 `6ccfe901…`), whose nine RDR rows and two kata rows are the expected value" | A SHA-pinned rendered-text artifact is made the oracle for a value-level contract, contradicting this RDR's own dump clause: "Golden tests SHOULD therefore assert over the normalized value … and reserve rendered-text goldens for tests of rendering itself." The pinned text is produced by a renderer that joins set literals (`main.go:758`) — the identity the literals clause bans. | Any correct implementation of the literals clause changes the rendered bytes and fails the pinned SHA. The implementer sees a Stage-4-approved hash fail and reverts the correct behavior to match the hash. | §1, §2, premortem, AT-1, AT-5 |
| C-17 | §Prerequisites: "RDR 0007 is `Final`, so the *specification* is settled — but its kernel reshape is unimplemented … Phases 2 and 3 therefore sequence behind the reshape **in their entirety**" + "This item gates implementation sequencing, not lock" | The RDR's own plan makes Phases 2 and 3 — normalizer, dump, and all validation — unstartable, and no RDR, kata, or phase owns the reshape. Four of fourteen assumptions (A9, A10, A12, A13) are Pending on it or on owed fixtures. The RDR locks a plan it has just shown cannot start. | Implementation opens, finds `resolve.Row.Guard` is still `string` and `OpExists` does not exist, and builds a local mirror of the constants — the exact drift the Risks section names, whose only mitigation "compiles only once RDR 0007's reshape exports them." | §1, §2, premortem, AT-8 |
| C-18 | §Minimum Viable Validation: "The overlap item is the one MVV assertion this RDR cannot discharge alone … **deferred to the Phase 5 lint handshake**" + §Technical Design: "Ambiguous overlap between candidate rows is therefore a lint finding, not a load failure" | Every safety net for the format's dangerous outputs — dead rules, unreachable `[initial]`, ambiguous overlap, read-before-write, and C-13's live-row-where-a-dead-rule-was-written — routes to RDR 0006. The RDR discharges its own hazards by naming a lint that has no implementation, eight times. | The RDR fixture's `[initial] stage = "propose"` is matched by no rule. The declared root is unreachable. This loads clean, dumps clean, is the RDR's canonical example, and is caught only by the unimplemented reachability check. It ships as the reference the team copies. | §1, §2, §3, premortem, AT-9 |
| C-19 | §Normative Contracts, contexts clause: "Merging is idempotent on identical atoms … **Inheritance therefore never overrides — it only accumulates.** A narrowing 'override' of an inherited constraint is not expressible" | The entire case for Alternative 1 over Alternative 2 rests on a factoring mechanism with no override, verified on a fixture with three contexts and five rules against a Decision Rationale naming seven interacting dimensions. Any rule needing a different value on an inherited key must fork the whole chain. | Six weeks in the model has `large-prelock`, `mid-prelock`, `small-prelock`, `large-prelock-cluster`, each re-stating `status.eq = "Draft"`. Changing "Draft" means finding every copy — verbatim Alternative 2's rejection reason, relocated from rows to contexts where the dump cannot show it. | §2, §3, premortem, AT-10 |
| C-20 | §Normative Contracts, load clause: "**Load is fail-fast: the first category a document trips is the refusal.** … the order in which independent defects are checked is deliberately **unspecified**" | Fail-fast plus unspecified order chosen for a format whose primary journey is hand-authoring. The RDR concedes the fit ("a poor fit for the hand-authoring loop this RDR exists to serve — one edit-reload cycle per defect") and defers the remedy to RDR 0006 or "a distinct RDR 0005 verb," neither scheduled. | A new author writes a 300-line model. Lint reports one error. Fix, reload, one more. Fifty cycles later they abandon the sparse format and hand-write expanded rows — Alternative 2, the rejected option. | §2, §3, premortem, AT-11 |
| C-21 | §Critical Assumptions A1 Evidence: "35 negative fixtures refuse one category each … covering **18 of the 25 categories**"; Testing Strategy 3: "The seven still owed" → six names → "Two of those **five**" | Three different cardinalities in one paragraph for the completeness claim of the RDR's central API surface ("The Testing Strategy asserts on the category, so these identifiers are an API surface"). Nobody can determine what is covered. | The implementer promotes 35 fixtures, believes the floor is met, and ships with `malformed TOML`, `missing recognized outcome alphabet`, `malformed model declaration`, and `malformed dump declaration` unimplemented — the four categories a real author's first file is most likely to hit. | §1, §2, premortem, AT-12 |
| C-22 | §Normative Contracts, accessor clause: "every **owned** tag MUST be served by exactly one reader" | An unconditional load refusal on an operational property. A write-only owned tag (an audit stamp, a one-way flag) is unauthorable. The RDR carefully provenance-scopes the *observed* arm — "zero is legal … refusing it would make every `--tag`-supplied key unauthorable" — and performs no equivalent check on the owned arm. | Author adds `[tags.audit_stamp] provenance = "owned"` and writes it once. Refused: `malformed accessor binding: owned tag "audit_stamp" served by 0 readers`. To ship it they invent a fake reader with a `path` and `timeout` RDR 0004 will execute on every resolve. | §3, premortem, AT-13 |
| C-23 | §Normative Contracts, next/writes clause: "MUST NOT populate one by aliasing the other" + Testing Strategy 2: "**The no-alias obligation … is not assertable by value comparison** … Enforcement is by review of the normalizer" | A MUST NOT whose oracle the RDR itself declares unassertable, falling back to human review — the one control the RDR's `oracle` mini-check elsewhere rejects ("a new assertion is not accepted here until a wrong implementation that fails it is named"). | Latent. A later RDR admits a next-state tag that is not an owned write; both fields move together through a shared backing array; the accessor writes a tag the model never authorized. Found as artifact corruption, not a compile error. | §1, premortem, AT-14 |
| C-24 | §Round-Trip / Inverse Invariants: "set-valued and multi-entry fields are rendered with unescaped separators. This makes the rendered form of a **set-valued atom literal** non-recoverable" + Problem Statement: "every legal edge to live in one reviewable artifact, so 'is this transition legal' has one answer" | The stated user outcome is one reviewable answer; the RDR then specifies the review surface to be provably ambiguous and accepts it — "A reviewer needing the exact members reads the source rule the locator names." The artifact defers to the source it exists to replace. | A reviewer diffs two dumps, sees `labels=needs work` in both, approves. One was `["needs work"]`, the other `["needs", "work"]`. Set-membership behaves differently in production; the diff meant to catch it showed no change. | §2, §3, premortem, AT-5 |

---

## 1. The three most likely ways implementation goes wrong

### 1.1 The RDR cites Final peers as authority for three contracts those peers do not contain, and contradicts four contracts they do

This is the deepest failure, and it is structural rather than incidental. RDR 0002 is the cluster's wire-format producer and locks **last** — RDR 0006, 0007, 0008, 0009 are all `Final`. A producer locking after its consumers has exactly one obligation: reconcile against what they actually say. RDR 0002 instead paraphrases them from memory, and the paraphrases are wrong in both directions.

**Three fabricated citations.**

*One.* The next/writes clause states, twice:

> RDR 0009 A4 fixes that these are distinct fields and that only the writes reach the accessor layer, so an escape row, which carries neither a write block nor a clear list, normalizes to a row with both empty.

> They stay distinct fields rather than collapsing to one because RDR 0009 A4 fixes them as distinct…

RDR 0009 A4 says the opposite. Its subject is whether the escape conformance predicate should widen to `NextTags`, and its verdict is that it should not:

> - **Status**: Verified — settled **closed**: the predicate does **not** widen; it stays `Writes`-only.

and normatively:

> The predicate is Writes-only and does NOT extend to NextTags: A4 settled that owned state is reachable only through a write accessor, which RDR 0004 scopes to "planned owned-tag writes," so an escape row's NextTags cannot mutate owned state.

A4's holding is that **no rule constraining an escape row's `NextTags` is needed at all**. RDR 0002 cites it as the authority for a MUST requiring `NextTags` empty on escape rows, and for the distinctness of two fields A4 merely observes in passing. The design choice may be right on its own merits; the citation is invented.

*Two.* The reserved-value clause states as settled fact:

> RDR 0004 removes the key on that write and asserts absence on read-back, and RDR 0006 reads it as removal.

The word "clear" appears **zero times in RDR 0004's body** — every occurrence is in its Status line or its Refinement Context. RDR 0004's Fidelity Table asserts the contrary with no exemption:

> \| `write -> read`, owned tags \| Every planned owned-tag key/value **present** in the re-read equals the transition plan's expected value \| value equality over the planned key set \| **none** \|

RDR 0004's own Refinement Context names the defect and schedules the repair for a future stage — it has not landed. RDR 0002's Capability Dependencies row is honest about this ("RDR 0004 owns … what a `<clear>` write and its read-back mean" — **Pending**) while a normative fence 260 lines earlier asserts the contract exists.

*Three.* The guard clause states:

> the block it was authored in — a three-valued domain, `match`, `all`, or `unless` — the atom shape JDR 0001 §D1 fixes and **RDR 0007 spells**.

RDR 0007 is Final and spells two values:

> A candidate row carries its guard as a slice of parsed atoms — key, operator token, literal, block ∈ {all, unless} — the shape JDR 0001 §D1 fixes; this RDR cites it and does not restate the grammar.

RDR 0003 quotes the same two-valued domain twice. The widening to three values is legitimate and authorized — by **JDR 0001 §D6**, which explicitly books it to land in RDR 0002. RDR 0002 has the right decision and the wrong citation, and it attributes the widening to a locked document that says something else.

**Four contradicted contracts.**

RDR 0006 (Final) publishes a **minimum input contract**. RDR 0002 violates four of its clauses.

*Row kind.* RDR 0006 requires it as an input, listed separately from the class list:

> normalized candidate rows with deterministic row identity, source rule id, optional source span, match predicates, guard predicates …, writes, clears, and — for escape rows — **row kind and** the declared list of failure classes the row rescues (RDR 0002);

RDR 0002 forbids it: "Row kind … is a **derived view property, never a row field** … normalization introduces no field to carry it." And RDR 0009 (Final) *quotes RDR 0002 verbatim* — twice — as rendering "a candidate row **with row kind `escape`**," a sentence RDR 0002 has since deleted and inverted, and on which RDR 0009 A4's Verified consistency note is grounded.

*Terminal shape.* RDR 0006 A6 states a shape requirement **and an explicit prohibition**:

> both declarations are **tag predicates, not state names** … an `initial` table of owned `tag = value` assignments fixing the root node, and a `terminal` list of predicates over owned tags in the same atom shape rules already use. RDR 0002 owns the final spelling; what this RDR requires is that **neither is a bare identifier a row references by name**, since invariant 1 checks them as tag keys and values and invariant 2 evaluates a terminal as a predicate over a node.

RDR 0002 ships exactly a bare identifier referenced by name — "root `terminal` is a list of context ids, each of which MUST resolve" — and never states the dereference from context id to predicate set. A context id *can* dereference to a predicate set, which is why JDR 0001 §D7(i) blessed the spelling; but the dereference is the whole content of the obligation, and RDR 0002 does not write it. RDR 0006's invariant 2 splits nodes on "terminal-participating keys," which requires the terminal to *be* atoms.

*Advisory tier.* RDR 0006's advisory tier is normatively closed:

> The advisory tier is closed at `graph-coverage-closed-by-escape`, `graph-redundant-row`, `graph-unreachable-rule`, and `graph-vacuous-atom`.

RDR 0002's disposition table routes a fifth: "\| expansion counts per rule \| lint \| RDR 0006 \| reported, non-blocking \|" — and this is not decorative. It is the stated mitigation for the format's headline risk: "**Risk**: Sparse contexts hide an accidental Cartesian product. **Mitigation**: … lint must report expansion counts per rule." The mitigation for the single risk that most threatens the sparse format is an output the lint's own closed tier forbids. Compounding it, RDR 0006 contains zero occurrences of "expansion," "expansion suffix," or "per-member" — it does not know rows expand at all, and its finding identity is keyed on a rule id that expansion makes non-unique.

*Literal type.* RDR 0002 spends two normative fences banning joined literal renderings — "An implementation MUST NOT derive that key, or any other identity, by joining members into a string." RDR 0006's mandatory finding payload declares `Literal` a flat Go `string`, bound by a MUST. The banned rendering is reintroduced at the locked diagnostic boundary, and RDR 0002 never notices — even though its own desk trace flags the joined-rendering defect as a live CONTRADICTION twice.

**Symptom.** The implementer resolving RDR 0002 against RDR 0006 has two documents, one Final, giving contradictory input contracts, and no rule for which wins. RDRs are not amendable. What actually happens is a Slack thread, a coin flip, and a comment in the source saying the other RDR is stale — which is how a cluster stops being a specification.

### 1.2 The normative fixtures do not satisfy the normative contracts, so the implementer deletes contracts to make the fixtures pass

**Root cause.** This RDR promotes `evidence/spikes/iter-2/` from "a spike showing the shape is possible" to "the normative fixture set implementation must promote and must not narrow" — in Testing Strategy 1, 2, 3 and 5, in Phase 1 ("**Discharged at Stage 4**"), and in the MVV ("a fixture may not be narrowed on promotion, only extended"). Meanwhile the Normative Contracts were tightened past what the spike implements, repeatedly, across three pre-lock lenses, without re-running the spike to conformance. The contracts now refuse their own evidence.

**The enabling passage.** The dump clause:

> `[dump]` carries exactly one key, `order`: a list of column identifiers. The identifiers are the **field names enumerated in the sentence above**, snake_cased … (`row_identity`, `source_locator`, `outcome`, `atoms`, `gate`, `next`, `write`, `requires_owned`, `escape`, `kind`) — that closed list is the column vocabulary … An `order` naming an unknown identifier, repeating one, or omitting one is a `malformed dump declaration` refused at load — not a silently truncated dump.

Both normative fixtures author:

```toml
[dump]
order = [
  "identity", "source", "kind", "outcome",
  "atoms", "next", "writes", "requires_owned", "gate", "escape",
]
```

`identity` ≠ `row_identity`. `source` ≠ `source_locator`. `writes` ≠ `write`. Three of ten outside the closed set. Per this RDR's own clause both fixtures MUST refuse — so Testing Strategy scenarios 1, 2, 3 and 5 all begin by loading a document the contract forbids.

I ran the spike. Both fixtures load clean; so does one whose `order` I mutated to `["identity", "NOT_A_COLUMN"]`, which dumped all nine rows. `[dump]` is decoded into `main.go:122 type Dump struct { Order []string }` and never read again. The RDR states this itself: "the spike decodes `[dump]` and never reads it."

So `malformed dump declaration` — one of the two categories this pre-lock pass *added* — has a fully specified rule, zero implementations, zero fixtures, and two normative fixtures that violate it.

**Symptom.** The implementer writes the validator to spec; every fixture-backed test fails at load. They hold a §D7-blessed fixture set, a SHA-pinned `output.txt`, and a Prerequisites checkbox saying Phase 1 is discharged. The cheap reconciliation is that the validator is wrong. The category is retired the week it is written — exactly the failure mode the RDR warns about for `unknown schema field`: "a parser swap that silently lost it would retire the … category without any contract appearing to change."

The trap is set twice. Testing Strategy 2 pins `output.txt` by SHA-256 as "the expected value" for a **value-level** contract, contradicting the dump clause's own advice ("Golden tests SHOULD therefore assert over the normalized value … reserve rendered-text goldens for tests of rendering itself"). The pinned text comes from a renderer that joins set literals (`main.go:758`). Any correct implementation of the literals clause changes those bytes and fails the hash.

### 1.3 Guard blocks are validated by a weaker rulebook than match blocks, so three normative categories are enforced in one block out of three

**Root cause.** The predicate clauses are stated in block-neutral language; the only executable evidence enforces them in one block. Nobody noticed because every negative fixture that exercises a predicate rule mutates a *match* block.

**The enabling passages.** Three clauses, none qualified as to block. The operator clause — which explicitly notices the gap and then leaves it:

> A token under `[rule.guard.all.<tag>]` or `[rule.guard.unless.<tag>]` outside the set is a `malformed predicate atom` (`unknown operator`) refused at load — the guard-side mirror of the match-block restriction above, **which without this was absent**.

The reserved-value clause: "`<clear>` MUST be refused at load wherever a tag value is authored — a write-block value, an `[initial]` assignment, or a predicate literal." The type-model clause: "`malformed predicate atom` (a literal outside the tag's declared domain)."

`main.go:529-537` is the entire guard-block validation:

```go
for _, block := range []map[string]map[string]any{rule.Guard.All, rule.Guard.Unless} {
    for key := range block {
        if key == recognizedKey { return ... malformed outcome binding ... }
        if _, ok := m.Tags[key]; !ok { return ... unknown tag ... }
    }
}
```

By contrast `checkMatchBlock` (`main.go:552-577`) checks operator membership, `<clear>`, `#`-in-`in`-members, and domain conformance. I ran all three mutations:

| Mutation on `[rule.guard.unless.profile]` | RDR says | Spike does |
|---|---|---|
| `frobnicate = "small"` | `malformed predicate atom` (unknown operator) | loads; emits `profile.frobnicate=small@unless` |
| `eq = "NOT_A_PROFILE"` (domain `["small","mid","large","foundational"]`) | `malformed predicate atom` (literal outside domain) | loads; emits `profile.eq=NOT_A_PROFILE@unless` |
| `eq = "<clear>"` | `reserved tag value` | loads; emits `profile.eq=<clear>@unless` |

The `<clear>` case is the one the RDR argues hardest for: "A predicate literal of `<clear>` is refused rather than admitted as a dead atom because no view can ever hold it and a silent dead rule is the more dangerous failure." The evidence produces the more dangerous failure.

**Symptom.** Three different failures, all bad.

*Unknown operator* reaches RDR 0007's seam, which "MAY answer unevaluable for a present value it cannot compare (A18)." The refusal is `guard_unevaluable`, reason `uncomparable`, naming the **key** — not the typo. RDR 0007 predicts the response: "the risk is that it pushes authors toward sentinel-stamping to make the kernel talk, which is the named anti-pattern."

*Out-of-domain literal* is silent and inverts the rule. An `unless` guard that can never be true never excludes anything, so a rule the author explicitly carved out now fires. The reviewable artifact answers "this edge is legal" for an edge the author wrote a guard to forbid — the same user-outcome inversion the RDR identifies as C-13's severe arm, reached by a route the RDR never considers.

*`<clear>` literal* triggers RDR 0007's resolution-scope veto: "if any surviving candidate row's guard is GuardUnevaluable, the resolution MUST refuse `guard_unevaluable`; a decided-GuardTrue sibling MUST NOT be selected while an unevaluable candidate exists." One unsatisfiable atom, and every resolve touching that outcome refuses non-escapably. The flow deadlocks.

---

## 2. The one section that will be rewritten within 6 weeks of shipping

**§Normative Contracts — the `[dump]` clause, and with it the whole Round-Trip / Inverse Invariants section.**

Three commitments in this section cannot coexist with a working review loop.

*One.* It fixes a closed column vocabulary its own normative fixtures violate in three places (C-1). The first implementation session forces a choice: change the fixtures, change the vocabulary, or delete the validator. Whichever wins, this clause is edited.

*Two.* It makes `[dump]` a load-refusable schema — "refused at load — not a silently truncated dump" — for a feature that is pure presentation, has no consumer, no implementation, and no fixture, while advising implementers in the next paragraph not to golden-test the thing it validates: "the cost is that a `[dump]` edit churns golden files asserted over rendered text. Golden tests SHOULD therefore assert over the normalized value." A refusable schema for a field nobody should assert on is the first thing cut when the category floor is triaged against the seven-or-six-or-five owed categories (C-21).

*Three, and fatally.* The Round-Trip section concedes the review artifact is ambiguous and accepts it:

> One site makes the rendered form lossy today … set-valued and multi-entry fields are rendered with unescaped separators. This makes the rendered form of a **set-valued atom literal** non-recoverable — `["a,b", "c"]` and `["a", "b,c"]` render identically…

> **Phase 2's dump is not required to carry an escaping grammar, and that is a decision, not an omission.** … What the dump owes the reviewer instead is that the ambiguity be **visible rather than silent** … A reviewer needing the exact members reads the source rule the locator names.

Read against the Problem Statement:

> A flow author needs every legal edge to live in one reviewable artifact, so "is this transition legal" has one answer instead of being reconstructed from scattered prose.

The dump *is* the reviewable artifact. "One answer" has become "an answer, plus a bracket meaning *a difference may be hiding*, plus a trip back to the source." The first time a reviewer approves a diff that showed no change while semantics changed (C-24), an escaping grammar gets added — and the RDR names the cost: "Defining one would make the dump a re-readable format — an inverse this RDR explicitly does not claim." The inverse gets claimed.

The Round-Trip invariant will not survive either, for a reason internal to it. It compares "the locator's rule-identifying part" — but the RDR gives the locator two origins one clause apart (C-15), and the SHA-pinned witness shows two kata rows sharing `kata:review`. The invariant is stated over a projection no witness computes, of a field whose derivation the document specifies two incompatible ways.

Two close runners-up, both of which will also be rewritten, and one of which may be rewritten first:

**§Prerequisites** (C-17) — a lock document whose second checkbox is unchecked, whose body concedes "Phases 2 and 3 therefore sequence behind the reshape in their entirety," and which then asserts "This item gates implementation sequencing, not lock." It also never mentions RDR 0006, a `Final` peer that has booked two unchecked prerequisites naming RDR 0002 by name.

**§Conditional Mini-Checks `disposition` table** (C-11) — the moment someone tries to build the expansion-count report the Risks section promises and discovers RDR 0006's advisory tier is normatively closed at four members that do not include it.

---

## 3. The one assumption that will not survive first contact with a real user

**That shared contexts with accumulate-only inheritance are sufficient factoring for the RDR flow — A6, and the entire case for Alternative 1 over Alternative 2.**

A6: "**Shared contexts and positive/negative guard lists are sufficient to keep the RDR model sparse without hiding ambiguity.** — Status: Verified — Method: Spike."

The mechanism, fixed normatively:

> Merging is idempotent on identical atoms: the same atom contributed by a rule and by one or more inherited contexts collapses to one. **Inheritance therefore never overrides — it only accumulates.** A narrowing "override" of an inherited constraint is not expressible, and is not silently approximated: authoring two atoms on one key that no view can satisfy together yields a dead rule, which is RDR 0006's unreachable-rule finding, not a load failure here.

The evidence is a fixture with **three** contexts in one linear chain and **five** rules. The workload the Decision Rationale describes:

> A fully expanded table would make the RDR model's dimensions multiply: status, profile, stage, prelock iteration, cluster eligibility, rewind scope, and guards would force authors to copy the same predicate fragments across rows. That is exactly the DX failure this RDR must avoid.

Seven dimensions. Apply accumulate-only inheritance to seven dimensions.

`large-prelock` inherits `prelock` inherits `draft`, contributing `status.eq=Draft`, `stage.eq=prelock`, `profile.in=["large","foundational"]`. A rule needing the same status and stage but `profile = "mid"` **cannot use this chain** — not "should not," *cannot*. Adding `profile.eq=mid` to a rule inheriting `large-prelock` produces two `eq` atoms on `profile` with different literals, which the RDR itself classifies as a dead rule routed to an unimplemented lint. The only option is a parallel chain re-stating the shared atoms.

That works while the dimensions form a clean tree. They do not. Cluster eligibility crosses profile. Rewind scope crosses stage. Prelock iteration crosses profile — the cap-3 counter the Ragel POC contrast names as one of "the hard parts." Every crossing that is not a strict refinement of an existing chain forces a new chain that re-states its shared atoms, because the only composition operator in the language is *add more constraints* and there is no operator for *same, except*.

This is Alternative 2's rejection reason, relocated: "Repetition makes edits risky: changing one shared condition requires finding every copied row." At six weeks the model has `large-prelock`, `mid-prelock`, `small-prelock`, `large-prelock-cluster`, `mid-prelock-rewind`, each re-stating `status.eq = "Draft"`. The Cartesian product was not eliminated; it moved from rows to contexts, where it is **less** visible, because a context definition is not a row and the dump — the review artifact — shows only rows. The reviewer sees nine expanded rows and cannot see that four came from three near-identical contexts.

The spike cannot detect this. Three contexts and five rules is below the threshold where accumulate-only inheritance binds, and A6's "Verified" rests entirely on it. `use = ["a","b"]` helps with orthogonal dimensions and not at all with crossing ones, because two contexts constraining one key with different literals is the dead-rule case again.

Four secondary assumptions fail in the same shape — verified on a five-rule fixture, deployed against a seven-dimension model:

- **Fail-fast load (C-20).** The RDR concedes the fit ("a poor fit for the hand-authoring loop this RDR exists to serve — one edit-reload cycle per defect, with no guarantee about which defect surfaces first") and prices the cost against a 200-line fixture. Against a 300-line hand-authored model, K defects means K cycles with "deliberately **unspecified**" ordering. The batch remedy is deferred to two unimplemented surfaces.

- **Every owned tag needs a reader (C-22).** A hard load refusal that makes a write-only owned tag unauthorable. The rule was written to catch a *typo*; it is enforced as a *modeling* constraint, and the first flow wanting an audit stamp must invent a fake reader RDR 0004 will execute on every resolve.

- **The overlap safety net (C-18).** "The overlap item is the one MVV assertion this RDR cannot discharge alone … deferred to the Phase 5 lint handshake." Every hazard this RDR creates and declines to refuse routes to RDR 0006, which has no implementation. Until it exists the format has no safety net at all — and the RDR's own canonical fixture already trips one of the checks: `[initial] stage = "propose"` is matched by no rule, so the declared root is unreachable.

- **The predicate-placement decision (C-14).** RDR 0007 (Final) states that placement decides escapability and hands the mitigation to RDR 0003, not to RDR 0002. RDR 0002 owns the authoring surface, calls the match block "select on," and never tells the author that a match-block predicate over a possibly-absent key produces an escapable `no_match` an escape row will convert into a plan.

---

## 4. Premortem

*Written from twelve weeks after ship.*

We shipped `internal/table` in week five. The flow team wrote the first real RDR transition model in weeks six through nine. In week eleven a rewind silently skipped two stages on a live RDR, and the postmortem took four days because the reviewable artifact said the transition was legal.

**Week 1 — the fixture contradiction, and the first contract deleted.**

We wrote `Load`, `normalize`, and the dump validator straight from Normative Contracts. `validateDump` enforced the closed vocabulary.

Every fixture-backed test failed at load:

```
malformed dump declaration: unknown column "identity"
```

Both promoted fixtures author `["identity", "source", "kind", ...]`. Three identifiers outside the closed set. The RDR marked Phase 1 **Discharged**, pinned `output.txt` with SHA-256 `6ccfe901…`, and said the fixtures "may not be narrowed on promotion, only extended." Nobody believed they had standing to edit a Stage-4-approved fixture set to satisfy a validator for a field with no consumer.

`validateDump` became a `//nolint` warning on day three. `malformed dump declaration` never shipped. Neither did `malformed model declaration` — same triage, same session, because the paragraph naming them could not agree whether seven, six, or five categories were outstanding. We shipped 18 of 25 and recorded it complete.

**Weeks 2–4 — the guard-block asymmetry ships.**

`checkMatchBlock` was written from the match-block clause and is thorough: operator restriction, `<clear>`, `#`-in-`in`-members, domain conformance. The guard path was written from `main.go:529-537`, the only implementation in evidence, and checks two things: key declared, key not `recognized`.

Nobody caught it, because no fixture would. All 35 promoted negatives mutate match blocks, `[initial]`, accessor tables, or the alphabet. Not one mutates a guard-block operator, a guard-block literal's domain, or a guard-block `<clear>`. The RDR added all three obligations in prose during pre-lock — the operator clause even says the guard-side rule "without this was absent" — and added no fixture for any. `gen-cases.py` derives cases from the RDR fixture by single mutation, and the mutations it knows are the ones that already existed.

**Week 5 — the lint input contract does not typecheck.**

Starting Phase 5, we opened RDR 0006's minimum input contract and could not satisfy it. It requires, for escape rows, "**row kind and** the declared list of failure classes" — two items. RDR 0002 forbids the first: "normalization introduces no field to carry it." It requires `clears` as a field distinct from `writes`; RDR 0002 folds clears into writes as `<clear>` entries and lists no `clears` in the dump field list. It requires `terminal` to be tag predicates and states the prohibition in so many words — "neither is a bare identifier a row references by name" — and RDR 0002 ships a list of context ids.

Then we found RDR 0009, also `Final`, quoting RDR 0002 as producing "a candidate row **with row kind `escape`**" — a sentence RDR 0002 no longer contains and has inverted. RDR 0009's A4 Verified consistency note is built on that quotation.

Two days on a Slack thread. Resolution: follow 0002, add a comment in `row.go` saying RDR 0009's citation is stale, and hand-derive `clears` by string-matching `<clear>` in the writes list. The cluster now has two `Final` documents whose load-bearing citations of their producer are false, and no mechanism that will notice.

**Week 6 — the expansion-count report cannot be built.**

The Risks section's mitigation for "sparse contexts hide an accidental Cartesian product" is "lint must report expansion counts per rule." RDR 0006's advisory tier is normatively closed at four members and expansion counts is not one; RDR 0006 contains no occurrence of "expansion" at all. Worse, RDR 0006's finding identity is keyed on the source rule id, and after expansion one rule yields N rows sharing that id. We shipped no expansion diagnostic. The mitigation for the format's headline risk does not exist.

**Week 8 — `stage` goes missing and the escape row rescues it.**

`resolveRDR` on `rdr-flow.toml`, against a repo where `docs/rdr/0047-*.md` had truncated front matter. `readAccessor("rdr-status")` returned `status=Draft` and no `stage`.

`assemble` built the view. `Resolve` filtered on outcome — `round-clean` — and reached `view.matches(row.Match)` for `continue-prelock#large`. `Match` carries `stage.eq=prelock` because `stage` is a **selection** and the author put it in `[context.prelock.match.stage]`, exactly as the RDR's own fixture and Illustrative Code do. `TagSet.matches` hit `tv, ok := s.tags["stage"]`, got `!ok`, returned false. Same for `#foundational`, same for both `continue-prelock-cluster` rows.

Zero candidates. `gate` returned `nil, nil` — `missingOwned` scans **survivors**, and a match-failed row is never a survivor, so `stage` was never reported missing despite being in `RequiresOwned`. `switch len(selected)` hit `case 0`. `escapeOrRefuse(in, view, Refusal{Kind: KindNoMatch})`.

`draft-no-match-escape#round-clean` rescues `no_match`, binds `round-clean`, and its `Match` is `status.eq=Draft` — which held. Empty guard, empty `RequiresOwned`. `gate` passed it. One viable escape. `planOf(in, viable[0], true)`.

We returned a **plan**. `Escaped: true`. The accessor applied it.

The RDR's own Load-Bearing Decision forbids this: "Count-first pruning is rejected because it launders missing state into an escapable `no_match` (JDR 0001 §D2)." We did not count-first prune. We match-first pruned — which the authoring surface makes the default — and got identical laundering. RDR 0007 (Final) documented it months earlier: "predicate PLACEMENT decides whether an absent key refuses non-escapably or escapes … The only control is authoring guidance … which Phase 4 hands to RDR 0003." RDR 0002, which owns the surface, never picked it up, and its routing clause calls the match block "select on."

**Week 9 — the unless-guard that never fired.**

Investigating week 8 we found `[rule.guard.unless.profile] eq = "smal"`. A typo, six weeks old. The declared domain is `["small","mid","large","foundational"]`. Per Normative Contracts this is `malformed predicate atom` — "a literal outside the tag's declared domain." Our loader enforces domain in match blocks only.

The negation never matched, so it never excluded anything, so `continue-prelock` had been firing for `small` profiles since week 6. Every small-profile RDR ran the large-profile prelock path. The dump showed `profile.eq=smal@unless` on two rows. Nobody reading nine rows in review noticed one missing `l`.

**Week 11 — the rewind that skipped two stages.**

`reconcile-rewind` writes `stage`, `status`, `rewind_scope` and clears `prelock_lens`. In week 7 the flow team refactored `[context.draft]` and added a second context carrying a `profile` constraint. Both were reached through `use`; one contributed `in = ["small,mid", "large"]` — a paste artifact from a spreadsheet column — while the rule contributed `in = ["small", "mid,large"]`.

Our `mergeAtoms` keys on `Atom.identity()`, which calls `literalString()`, which is `strings.Join(a.Literal, ",")`. We inherited that from the spike. Both sort to `large,mid,small`. One key. One surviving atom. The other constraint vanished.

The RDR predicted it exactly. Its `trace` table carries a row marked **CONTRADICTION**:

> **CONTRADICTION** — the spike keys the merge on a comma-joined rendering (`main.go::Atom.identity` → `literalString`), so `in = ["a,b", "c"]` and `in = ["a", "b,c"]` on one key and block collapse to one atom and the rule normalizes to a live expanding row instead of a dead rule. Scalar witnesses cannot see it.

And A13's "If wrong" describes what we got:

> the surviving rule normalizes to a **live expanding row where the author wrote a self-contradictory (dead) rule** … the artifact answers "this edge is legal" for an edge the author did not write, which inverts the user outcome this RDR exists to deliver. It is also the arm that trips **zero** load categories.

The RDR locked with that contradiction open, A13 `Pending`, and the delimiter-bearing control "owed." It was owed to us. We inherited `identity()` verbatim because it was the only implementation in evidence and it passed every promoted merge fixture — `merge-idempotent.toml` and `merge-distinct.toml` both use single-member scalar literals, under which a joined rendering and a member sequence are the same string.

**Week 11, day 2 — the postmortem could not name the row.**

Two kata rows were implicated in the ambiguity that followed. `compareRefs` orders on `(RuleID, SourceLocator)`; RDR 0005 printed the locator. Both printed `kata:review`.

Our `Locator` is `rule.Source`, from the spike. RDR 0002 says the locator is "**derived, not authored**: the loader composes it from `(model id, rule id)`" — and one clause earlier says each rule carries "an optional `source` — the authored provenance string the locator derives from." Two origins, one document. The RDR flags the consequence in a parenthetical — "which is why both `kata-fixture.toml` rules carry one identical locator; that is a spike defect against this clause" — in the same paragraph that promotes the spike as normative.

Four days. The reviewable artifact, the one that was supposed to give "is this transition legal" one answer, gave us `kata:review` twice and a nine-row dump in which `smal` and `small` are one character apart and `["a,b","c"]` and `["a","b,c"]` are the same bytes.

**Week 12 — the clear that read-back as a failure.**

We wired RDR 0004's write accessor. `reconcile-rewind` clears `prelock_lens`; the executor applied the plan, re-read the role, and reported a write failure. RDR 0004's Fidelity Table: "Every planned owned-tag key/value **present** in the re-read equals the transition plan's expected value … Lossy exemptions: **none**." A cleared key is absent.

RDR 0002 told us this was settled: "RDR 0004 removes the key on that write and asserts absence on read-back." The word *clear* appears zero times in RDR 0004's body. RDR 0004's own Refinement Context names the contradiction and schedules the repair for a stage that has not run. We special-cased `<clear>` in the executor against a Final RDR that says otherwise.

**What we never got.** RDR 0006 never shipped. Every check that would have caught weeks 8, 9 and 11 — unreachable-rule, ambiguous overlap, read-before-write, root reachability — is RDR 0006's by the arity split. RDR 0002 routes hazards there eight times. We had none of them. We also never got the reshape: `resolve.Row.Guard` is still `string`, `OpExists` still does not exist, and `internal/table/atom.go` carries a local mirror of constants that RDR 0002's Risks section named as the drift it feared and whose only mitigation was an assertion that "compiles only once RDR 0007's reshape exports them."

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

Review-time gates, executable against the RDR document and its evidence directory before lock.

**AT-1 — Normative fixtures must satisfy every normative clause of the RDR that promotes them.** *(C-1, C-2, C-16)*

```gherkin
Given an RDR designating fixtures as "normative", "promoted", or "the expected value"
When every normative MUST clause in that RDR is applied to those fixtures
Then no fixture violates any clause
And for each clause a fixture exercises, a witness exists in the evidence transcript
```
Steps: extract the closed `[dump]` vocabulary; extract `order` from each fixture; assert set membership both ways. **Fails today**: `identity`, `source`, `writes` are outside the set. The RDR could not have been marked Phase-1-Discharged with this gate in place.

**AT-2 — Every block-neutral predicate clause must have one negative fixture per block.** *(C-3, C-4, C-5)*

```gherkin
Given a normative clause constraining predicate atoms
And the clause does not name a specific authored block
When the fixture set is enumerated
Then a negative fixture exists exercising it under `match`
And one under `guard.all`
And one under `guard.unless`
```
Applied to the three block-neutral clauses this demands nine fixtures. Zero exist; all 35 promoted negatives mutate match blocks or non-predicate surfaces. This turns "the guard-side mirror … which without this was absent" from a prose repair into a coverage obligation.

**AT-3 — Every citation of a peer RDR must resolve to text in that peer that supports the claim.** *(C-6, C-7, C-8)*

```gherkin
Given a sentence in this RDR attributing a contract to a named peer
When that peer is searched for the attributed content
Then text is found that states it
And the peer's own verdict on that question is not the opposite
```
Steps: for each `RDR 000N <assumption> fixes/states/requires …`, locate the assumption and read its Status and normative landing. **Fails three times today**: RDR 0009 A4 settles that the predicate does **not** widen to `NextTags`, the inverse of what RDR 0002 credits it with; the word "clear" appears zero times in RDR 0004's body while RDR 0002 asserts RDR 0004 "asserts absence on read-back"; RDR 0007 spells `block ∈ {all, unless}` while RDR 0002 attributes a three-valued spelling to it. This gate is mechanical, cheap, and is the single highest-yield check available to a producer that locks after its consumers.

**AT-4 — Every Final peer's stated input contract must be satisfiable from what this RDR produces.** *(C-9, C-10, C-11, C-12)*

```gherkin
Given a peer RDR with Status Final that publishes an input contract
When each item is matched against this RDR's normalized value and dump field list
Then every item is produced
And no item is normatively forbidden by this RDR
```
Steps: enumerate RDR 0006's minimum input contract. **Fails four times**: row kind is required and forbidden; `clears` is required as a distinct field and folded into writes; `terminal` is required as tag predicates with "not a bare identifier a row references by name" stated as the prohibition, and shipped as context ids; `Literal` is required as a flat `string` while this RDR bans joined literal renderings in every comparison. A fifth failure appears in the reverse direction: this RDR routes "expansion counts per rule" to an advisory tier RDR 0006 closed at four members.

**AT-5 — No rendered-text artifact may be the oracle for a value-level contract, and no clause may ship with a self-declared CONTRADICTION.** *(C-13, C-16, C-24)*

```gherkin
Given a normative contract stated over a normalized value
When the Testing Strategy names its oracle
Then the oracle is a value comparison, not a SHA over rendered text
And no correct implementation of any clause would change those pinned bytes
And no desk-trace row carries an unresolved CONTRADICTION verdict at lock
```
**Fails today**: Testing Strategy 2 pins `output.txt` by SHA as "the expected value" for a value-level contract, contradicting this RDR's own dump clause; and the `trace` table carries two open CONTRADICTIONs whose fixtures are "owed."

**AT-6 — The authoring surface must diagnose any placement decision a Final peer says decides escapability.** *(C-14)*

```gherkin
Given a Final peer stating that an authoring choice this format exposes
      decides whether a refusal is escapable
When this RDR's authoring guidance is examined
Then it names the choice, states the consequence of each side, and specifies
     a load refusal, a lint finding, or explicit guidance
```
**Fails today**: RDR 0007 says "predicate PLACEMENT decides whether an absent key refuses non-escapably or escapes … the only control is authoring guidance … which Phase 4 hands to RDR 0003." RDR 0002 owns the surface, describes the match block as "select on," and never states that a match-block predicate over a possibly-absent key produces an escapable `no_match` an escape row converts into a plan. This gate forces the week-8 laundering path into the document before lock.

**AT-7 — Every derived field must have one origin and one witness.** *(C-15)*

```gherkin
Given a field the RDR describes as derived
Then exactly one clause states its derivation
And no other clause states a different origin
And a witness in the evidence shows the derived value, not an authored proxy
```
**Fails today**: "the loader composes it from `(model id, rule id)`" and "an optional `source` — the authored provenance string the locator derives from" are two origins; the witness shows `source=kata:review` twice, so the uniqueness the Round-Trip's "rule-identifying part" comparison depends on is unwitnessed. An RDR may not both promote an artifact as normative and describe it as defective in the same paragraph.

**AT-8 — Every phase must be startable on the repository as it exists at lock.** *(C-17)*

```gherkin
Given an Implementation Plan with numbered phases
When each phase's first deliverable is checked against current source
Then either the required surface exists
Or a named, scheduled, owned work item creates it
```
RDR 0002's Prerequisites already contain the failing evidence — "no `Atom` type, no operator-token type, no per-atom block field, and no `OpExists` / `LiteralTrue` / `LiteralFalse` anywhere in the repo" — and conclude "Phases 2 and 3 therefore sequence behind the reshape in their entirety." No item owns the reshape. This forces that owner to exist before lock instead of allowing "This item gates implementation sequencing, not lock."

**AT-9 — No hazard may be discharged onto an unimplemented consumer without a residual statement.** *(C-18)*

```gherkin
Given a hazard this RDR admits at load and routes to another RDR
When that RDR's implementation status is checked
Then it is implemented
Or this RDR states the residual and names what protects users until it lands
```
Steps: count routings to RDR 0006 (dead rule, unreachable rule, ambiguous overlap, read-before-write, missing root, missing stop set, terminal-on-non-owned, expansion counts); check RDR 0006 for an implementation. **Fails on eight routings with zero implementation.** The RDR does this correctly once — the MVV overlap item is explicitly deferred rather than counted — and then routes seven more without the same honesty.

**AT-10 — The factoring mechanism must be exercised at the dimensionality the RDR claims to serve.** *(C-19, and A6)*

```gherkin
Given an RDR whose selection rationale rests on a factoring mechanism
And whose Decision Rationale names N interacting dimensions
When the fixture is examined
Then it exercises at least N dimensions
And it includes at least one rule needing a narrowing of an inherited constraint,
    showing what the author must write
```
**Fails today**: the Decision Rationale names seven dimensions; the fixture has three contexts in one linear chain and five rules, and no rule narrows an inherited constraint. A6 is `Verified` by `Spike` on evidence an order of magnitude below the claimed workload. Running this gate surfaces the accumulate-only wall while Alternatives 2 and 5 are still live.

**AT-11 — The primary user journey must be walked end to end in the document.** *(C-20)*

```gherkin
Given a Problem Statement naming a human user and a task
When the RDR is read for that journey
Then a walkthrough exists covering author -> load -> diagnose -> fix -> re-load
And the number of iterations to fix K independent defects is stated
```
**Fails today**: the user is "a flow author"; the only user-facing loop is fail-fast with "deliberately **unspecified**" ordering, K defects requiring K cycles, and the batch mode deferred to two unimplemented surfaces. The RDR concedes "That cost is accepted here" without writing down what the cost is. Writing K down at review time is what makes it visible that the sparse format's authoring loop is worse than the expanded format it rejected.

**AT-12 — The category floor must have one arithmetic.** *(C-21)*

```gherkin
Given a category set the Testing Strategy asserts on
When every count in the RDR is extracted
Then all counts agree
And covered + owed == total
And each owed category names its fixture and its implementation owner
```
**Fails today**: "The seven still owed" → six names → "Two of those five." Three cardinalities in one paragraph, for the completeness claim of an API surface the RDR itself calls an API surface.

**AT-13 — Every load refusal must be checked against a legitimate authoring it forbids.** *(C-22)*

```gherkin
Given a normative load refusal
When one legitimate model shape is constructed that trips it
Then either the shape is genuinely illegal
Or the refusal is narrowed
```
Applied to "every **owned** tag MUST be served by exactly one reader" with a write-only owned tag: **the shape is legitimate and refuses.** The RDR provenance-scopes the *observed* arm with exactly this reasoning — "zero is legal … refusing it would make every `--tag`-supplied key unauthorable" — and never applies it to the owned arm.

**AT-14 — No MUST NOT may ship with an oracle the RDR declares unassertable.** *(C-23)*

```gherkin
Given a normative MUST or MUST NOT
When its verification is located
Then it is a mechanical assertion or a compile-time property
And it is not "by review", "by inspection", or "not assertable"
```
**Fails today**: "MUST NOT populate one by aliasing the other" is verified by "review of the normalizer, or by a test that mutates one field and asserts the other is unchanged." The mutation test *is* mechanical, so the remedy is to require it rather than offer review as an alternative. The RDR's own `oracle` mini-check already demands this — "a new assertion is not accepted here until a wrong implementation that fails it is named" — and this clause exempted itself from the rule it states.
