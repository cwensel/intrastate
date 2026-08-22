Model: claude-opus-5

# Critique — Guard Predicate Exhaustiveness (iteration 2, re-entry)

Scope: the CURRENT draft, including the grounding iter-2 and 3amigo iter-2 edits.
Iteration-1 findings (scoped product, set-valued proof limit, error-ownership
handoff) are dispositioned and are not re-raised. Everything below attacks text
that survived — or was introduced by — this iteration's fixes.

The thesis of this critique: **the previous two lenses tightened the vocabulary
of a proof whose inputs do not exist.** A7, A8, and A9 were booked as the three
open items. They are not the three open items. They are the three *noticed* open
items, and the largest one — that RDR 0002 declares no tag domains at all — sits
underneath A2, which this RDR stamps `Verified`.

---

## Findings Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | A2 Evidence (§Critical Assumptions, line 135): "lint receives … finite domains for the guard dimensions … enum/boolean values as declared sets … bounded integers as `{min..max}`"; `Status: Verified` | The producer does not exist. RDR 0002's tag declaration is "tag name, provenance …, value kind, and optional accessor reference" (`0002:217-218`) — no enum value list, no `min`/`max`, no domain of any kind. `value kind` occurs exactly once in all of RDR 0002. A2 is the load-bearing derivation and it is `Verified` against a field set that carries zero domain data. A7 and A9 book *two specific missing subfields* (optionality, element universe) while the base domain declaration they are subfields of is itself absent and unbooked. | `intrastate lint` can never certify any row group exhaustive, for any operator, because no authored model can declare a finite domain. The guarantee this RDR exists to provide is unreachable on day one, and the author is told nothing — the feature simply never fires. | §1, §3, premortem, AT-1 |
| C-2 | §Technical Design (lines 374-377): "RDR 0002 supplies the normalized candidate rows and RDR 0006 supplies the graph-lint grouping context, such as one source-state/recognized-outcome selection group" | Circular delegation with no owner. RDR 0006 contains no definition of a row group, selection context, or grouping rule — `grep` for "row group\|selection context\|grouping" over RDR 0006 returns only `:115` (quoting 0003's own A2), `:252`, and `:501`, none of which define grouping. Meanwhile RDR 0006 A2 (`0006:111-119`) delegates the whole coverage/overlap derivation *to RDR 0003*. Each document names the other as the supplier of the one input the proof cannot be computed without. | Two implementers group rows differently — one by `(source-state, outcome)`, one by outcome alone, one by shared match block. The same model lints green on one build and reports `graph-coverage-gap` on another. The author cannot tell which answer is correct because no document states the rule. | §1, §2, premortem, AT-2 |
| C-3 | §Technical Design (lines 378-383): "this RDR reads it as **default-on for every scoped row group whose participating dimensions are all finitely declared** — an author opts out by leaving a dimension undeclared" | The opt-out mechanism is the same signal as C-1's missing producer, so opt-out is the universal default and the guarantee is silently off everywhere. Worse, "opt out by leaving a dimension undeclared" makes *forgetting* indistinguishable from *deciding*. The 3amigo pass introduced this clause to close P2-6; it converted an undefined trigger into a default that fails open. | An author who intended a proven-exhaustive routing table gets a green lint that proves nothing, because one dimension was never given a domain. There is no diagnostic — the RDR explicitly makes silence the opt-out signal. The first time they learn coverage was never checked is a runtime `no_match` in production. | §1, §2, premortem, AT-3 |
| C-4 | §Normative Contracts (lines 495-501): "A row 'can refuse' when the row carries a value atom whose key is not declared present-for-every-reachable-predecessor … Lint MUST decide this syntactically over declarations so the test is total" | "Reachable predecessor" is never defined anywhere. Used 4× in 0003 (`:192`, `:417`, `:497`, `:517`) and 5× in RDR 0006 (`:52`, `:94`, `:130`, `:283`, `:725`) — always in prose, never as a definition. No document defines the predecessor relation, its carrier (states or rows), or how it is computed; RDR 0002 has **zero** uses of the word and defines no state-graph relation, so the notion has no producer. The citation chain is circular: RDR 0006 A3 (`0006:130`) cites RDR 0003's normative contracts as its evidence, and RDR 0003 A6 (`:165`) cites RDR 0006's Technical Design as its evidence. The 3amigo disposition's warrant is doubly broken: it cites "line 512, RDR 0006 via A6" — 0003 line 512 is the *overlap-diagnostics* clause tail (the owned-tag clause is `:515-519`), and **RDR 0006 has no A6** (it carries only A1-A5, at `0006:98,110,122,137,150`). The new clause is strictly *stronger* than the undefined notion it inherits: it demands syntactic decidability and totality. | Lint's withholding decision is implementation-defined. The same model withholds the exhaustiveness claim on one build and certifies it on another. Scenario 6's negative control (a group "over always-present keys certifies green") cannot be authored, because nothing says how a key is proved always-present. | §1, §2, premortem, AT-4 |
| C-5 | A6 (lines 162-167), `Status: Verified`, Method: Peer RDR, Evidence: "RDR 0002 `Normative Contracts` require the model to declare every matched or written tag, including … provenance" | Evidence proves an adjacent fact, not the claim. A6's own "If wrong" names the load: "it cannot prove that owned tags are set before they are matched." Provenance *labels* being declared does not make write-reachability *decidable*; that needs the predecessor graph C-4 shows is undefined. The RDR's own Method vocabulary block (lines 286-291) forbids exactly this: "confirming a neighboring fact and stamping the assumption `Verified` is not verification." A6 also carries the normative clause at `:515-519` (owned-tag read-before-write), so a `Verified` stamp here licenses a MUST that no one can implement. | `graph-owned-before-write` either never fires (vacuously passing every model) or fires on legal models. An author sets an owned tag in a predecessor row and lint still rejects the match, with no way to see what reachability lint computed. | §1, §3, premortem, AT-5 |
| C-6 | A1 Evidence (line 128): "`sh check.sh guard-fixture.toml` validated four representative guard rows … transcript captured in … `output.txt`"; A1 `Status: Verified`, Method: Spike | The spike validates nothing about predicates. `evidence/spikes/check.sh` is an `awk` script that regex-matches operator *tokens* against `^(eq\|in\|lt\|lte\|gt\|gte\|exists\|contains)$` and prints a count. It never evaluates a predicate, never builds a product, never computes coverage or overlap, never type-checks a literal against a tag kind. `output.txt` is four lines: `operators=eq,in,lt,gte,exists`, `rules=4`. The RDR's Method vocabulary defines Spike as "verified by running code against a live service or fixture" — this runs a spelling checker against a text file. A1 claims the flows "fit a closed typed predicate vocabulary"; the spike cannot observe typing, fit, or closure. | The operator set ships unvalidated against the actual target flows. The first real flow needs an operator (a range, a prefix, a set-disjointness test) that no one discovered, because the "verification" only confirmed the fixture author typed operator names from the list they were copying from. | §1, §3, premortem, AT-6 |
| C-7 | §Technical Design (lines 369-373) and §Load-Bearing Decisions › Identity (lines 584-592): `unless` as a distinct conjunctive block; atom identity tuple includes `block` | RDR 0002's normative contract erases the block: "Guard predicates MUST be represented as positive `all` predicates and negative `unless` predicates. Normalization MUST combine both into one candidate-row predicate set before ambiguity checks" (`0002:305-309`). One predicate set. This RDR's identity tuple `(RuleID, SourceLocator, key, block, operator, literal)` and its `unless`-block-as-single-excluded-intersection semantics both require the block to survive normalization. It is contracted away by the producer. | Two atoms that differ only by block collide in diagnostics, so an overlap finding names the wrong atom. Worse: if normalization flattens `unless` into the predicate set, `unless` becomes per-atom negation — precisely the misreading §Risks names at lines 913-915 — and rows the author intended as disabled become qualifying rows. Runtime returns `ambiguous_match` for a table the author reads as unambiguous. | §1, §2, premortem, AT-7 |
| C-8 | §Technical Design, §Normative Contracts, §disposition table — the entire RDR: `grep escape` over 0003 returns only lines 80, 112, 114, 689, none about escape rows | Escape rows carry guards and are invisible to this RDR, **while the consuming RDR requires them in the coverage check**. RDR 0006 invariant 4 (`0006:277-279`): "for each state/outcome pair that claims closed coverage, finite-domain input assignments must either match exactly one modeled row **or a declared escape row**" — so 0006's coverage check counts escape rows as covering. RDR 0003's coverage identity (A2, line 135) is `union(row_i accepted assignments) == scoped product` with no escape-row term at all. Shipped `resolve.go::escapeOrRefuse` (`:473-498`) runs escape candidates through the *identical* `gate()` (`:485`), and RDR 0002 normalizes them as candidate rows with "its normal predicate set" (`0002:294-298`). The `disposition` table (lines 545-557) enumerates ten input classes; none mentions escape. | Two documents define coverage over different row sets. Lint either reports a spurious `graph-overlap` on a legal escape edge, or (following 0006) counts escape rows as coverage and certifies a group exhaustive whose "coverage" is supplied by an edge the resolver only consults *after* a refusal has already occurred. The author's escape hatch silently satisfies the exhaustiveness proof, and the runtime still refuses. | §1, premortem, AT-8 |
| C-9 | §Normative Contracts (lines 475-482): withheld claim "MUST name the participating row and the atom that can refuse"; §authority row "Withheld-claim lint artifact … it reuses 0006's existing blocking code" (line 538) | This RDR assigns a payload requirement to `graph-unprovable-coverage`, a code owned by `Final` RDR 0006 (`0006:302`), whose own finding contract (`0006:352-356`) requires "stable code, model identity, severity, human-readable message, and the source rule/context id" — no atom-level field. The 3amigo pass removed a peer-binding MUST from the narrowing clause (entry 3) but left this one, which binds a `Final` peer's finding payload just as hard, only less visibly. A8 tracks the *narrowing's* recording document; it does not track this payload extension. | JSON consumers get a `graph-unprovable-coverage` finding with no atom field, because RDR 0006's implementer built to RDR 0006's payload contract. The author sees "cannot prove coverage" and no indication of which atom caused it — exactly the actionability failure the clause was written to prevent. | §1, §2, AT-9 |
| C-10 | §Prerequisites (lines 932-933): "- [x] A1-A6 verified"; §Assumption Verification (line 1104): "A1 through A6 are `Verified`, each uses an allowed Method label, each has concrete Evidence" | The self-audit is the failure. Three of the six checked-off assumptions do not survive their own stated standard: A1 on a token-scanner spike (C-6), A2 on a nonexistent producer (C-1), A6 on an adjacent fact (C-5). The Assumption Verification section re-asserts each as sound in prose without re-testing it, and the Finalization Gate treats the `[x]` as discharged. This is the mechanism by which C-1/C-5/C-6 reached iteration 2 unexamined: the gate audits *label hygiene* (allowed Method vocabulary, no self-reference, no `Docs Only`) and reports that as verification. | Nothing directly. This is the meta-defect: it is why the RDR will lock with three false `Verified` stamps, and why the implementer will trust A2's product derivation as settled when they open Phase 2. | §2, §3, premortem, AT-10 |
| C-11 | §Minimum Viable Validation (lines 960-969) and §Testing Strategy Scenario 6 (lines 1049-1060) | The MVV cannot be authored in the fixture format. It requires "one row group that is domain-exhaustive yet contains a possibly-absent guard key" plus a negative control "whose keys are all declared always-present" — both need the optionality field A7 says does not exist. It also requires `contains` over a declared set-valued tag (Phase 3, line 989), which A9 says cannot be authored. Meanwhile §Scope Verification (line 1148) asserts the MVV "is in scope for implementation, not deferred." Three of the MVV's required cases are unauthorable and the gate declares the MVV in scope. | The implementer reaches Phase 3, discovers the acceptance fixture cannot be written, and either invents the missing declaration fields unilaterally (forking the format from RDR 0002) or ships with the acceptance gate skipped. Either way the operator vocabulary is accepted without ever being exercised. | §2, premortem, AT-11 |
| C-12 | §Metadata Status (lines 9-12) and §Prerequisites (lines 946-953) | Status says `Draft [revised from Final 2026-08-12; re-verify A5 …]` and Prerequisites carries an unchecked item stating the RDR 0007 kernel reshape has not landed — "Phase 1 cannot begin against the atom-slice shape until that reshape lands." So this RDR's Phase 1 is blocked on an unimplemented `Final` peer, its A7/A9 are blocked on a `Draft` peer (RDR 0002) that has not accepted the request, and its A8 is blocked on an open JDR entry needing a `Final` peer's route-back. Every one of the four blockers is outside this document's authority, and the RDR contains no statement of what happens if RDR 0002 declines. | Not user-visible directly; the symptom is schedule. This RDR cannot lock, cannot implement, and has no fallback path. It will sit at `Draft` while three other documents move, and the eventual implementer will build against whatever the peers happened to become. | §2, premortem |
| C-13 | §Technical Design (lines 354-360): "A predicate atom has four conceptual fields: tag name, operator, expected value, and **source identity**. … The atom shape itself is fixed at the kernel seam by JDR 0001 §D1 and stated normatively in RDR 0007; this RDR cites it rather than restating it." | The RDR restates the shape and gets it wrong, then claims to be citing. RDR 0007's SEAM clause (`0007:1250-1253`, `:1260-1265`) fixes the four fields as `Key`, `Operator`, `Literal`, **`Block`** — where `Block` is an exported string type with constants `BlockAll`/`BlockUnless`. RDR 0003 lists four fields and substitutes **`source identity` for `Block`**. Both documents say "four fields"; the sets differ. This also breaks C-7's identity tuple from the other end — the tuple names `block` as a member, so the Technical Design field list contradicts the Load-Bearing Decision two hundred lines later in the same document. | An implementer building the atom type from §Technical Design omits `Block` and adds a source-identity field the kernel does not carry. Guard atoms fail to round-trip through the kernel seam; `unless` atoms become indistinguishable from `all` atoms, and every `unless`-block exclusion is evaluated as a positive requirement — inverting the meaning of every negative guard the author wrote. | §1, premortem, AT-12 |
| C-14 | §Approach operator matrix (lines 334-338): `in` takes a "non-empty typed scalar set"; `contains` takes a "non-empty typed element set" — and §A9, which books only the *declared element universe* as missing | Set-valued **literal encoding** is missing, and it is a distinct gap from A9's element universe. RDR 0007 states the consequence against its own contract (`0007:1585-1595`): "for the SET-valued literals RDR 0003 gives `in` … and `contains` … it does not yet [have a byte-comparable spelling] — `resolve.Tag.Value` is a bare `string` and no element encoding is declared … **two `in` atoms on one key differing only by literal-set ORDER are not distinguishable by this tuple**", and it assigns the duty back: "0003's declaration MUST make that spelling canonical." So a `Final` peer has filed a normative obligation on this RDR that the RDR does not carry, does not book as an assumption, and does not list in Capability Dependencies. Note `in` is affected, not just `contains` — and `in` is in the target-flow subset the spike "verified" and the MVV exercises. | Two authored rows whose `in` literal sets differ only in order are the same predicate but different atoms. Diagnostics report the wrong atom; the kernel's `guard_unevaluable` payload sort is non-total, so the same input tuple yields different payload orderings across runs — breaking the replay determinism RDR 0001 REQ-1 requires. | §1, §2, premortem, AT-13 |

---

## 1. The three most likely ways implementation goes wrong

### 1.1 The finite-domain proof has no input, and lint silently never runs

**Root cause.** A2 — the derivation the whole document rests on — is stamped
`Verified` against a producer contract that carries no domains.

**Enabling passage.** A2's Evidence (line 135) states as established fact:

> lint receives the candidate rows that share a selection context from the
> normalized model **plus finite domains for the guard dimensions that vary
> inside that group: enum/boolean values as declared sets, set-valued tags as a
> symbolic set over the declared element universe, and bounded integers as
> `{min..max}`**.

RDR 0002 §Technical Design item 2 (`0002:217-218`) is the complete tag
declaration:

> Tag declarations: tag name, provenance (`owned`, `observed`, `recognized`),
> value kind, and optional accessor reference for observed or owned read-back.

That is the entire schema. `grep -n "value kind\|finite\|declared domain"` over
RDR 0002 returns exactly one line: `:218`. There is no enum value list, no
integer `min`/`max`, no boolean domain, no element universe. RDR 0002's
normative contracts (`0002:265-360`) require declaring *provenance*; they never
require declaring a *domain*.

This is not the A7 gap and not the A9 gap. A7 asks RDR 0002 for a per-tag
*optionality* field. A9 asks for a set *element universe*. Both are subfields
of a domain declaration that does not exist. The draft's own Capability
Dependencies table (lines 644-645) books those two subfields as `Requested` and
books the base capability — "Sparse transition-table container … RDR 0002 …
Pending" (line 638) — as though the container is the only thing needed. The
spike fixture (`evidence/spikes/guard-fixture.toml`) *does* carry
`domain = ["Draft", "Final", "Implemented"]` and `min = 0 / max = 3` — the
fixture invented a schema RDR 0002 does not specify, and the RDR then read its
own invented fixture back as evidence that domains are available.

Note the third confirmation: RDR 0006's lint input contract (`0006:251-252`)
lists "declared tags with provenance, value kind, **finite-domain metadata when
exhaustiveness is claimed**, and single-valued grouping when applicable". Three
documents all expect the finite-domain metadata to arrive from somewhere. None
of the three is the producer. The same line carries a second orphan of the
identical shape — `single-valued` appears **zero** times in RDR 0002, so
RDR 0006's invariant 5 (`graph-single-valued-state`) has the same missing
producer. This is not a one-off oversight in this RDR; it is a systematic
pattern across the cluster, where consumers list inputs and no one audits the
producer's schema.

One more consequence worth stating plainly: RDR 0006 A2 is marked `Verified`
(`0006:110-119`) on evidence citing "declared set-universe" domains. Those do
not exist (A9 concedes it). So a `Final` RDR carries a `Verified` assumption
resting on a `Draft` RDR's unbuilt capability. Locking 0003 does not fix that;
it propagates it.

**Symptom.** Every authored model has at least one guard dimension with no
declared domain, because no dimension can have one. Under the default-on rule
(C-3), a group with an undeclared dimension is not exhaustiveness-eligible, so
lint proves nothing about any group in any model. `intrastate lint` returns
green. The author believes the routing table is proven complete. The resolver
returns `no_match` on the first uncovered combination in production.

### 1.2 Row grouping is delegated in a circle, so two implementations disagree

**Root cause.** The one input that determines *what* is being proved — which
rows are compared against each other — is claimed by neither document.

**Enabling passage.** §Technical Design (lines 374-377):

> For lint, RDR 0002 supplies the normalized candidate rows and **RDR 0006
> supplies the graph-lint grouping context**, such as one
> source-state/recognized-outcome selection group.

Now RDR 0006 A2 (`0006:111-119`), `Status: Verified`, Method: Peer RDR:

> **A2 Predicate lint can decide overlap and coverage for the finite domains
> this graph claims exhaustive.** … **Evidence**: RDR 0003 A2 derives coverage
> as `union(row_i accepted assignments) == scoped product` …

RDR 0003 points at RDR 0006 for grouping; RDR 0006 points at RDR 0003 for the
proof — and the proof is only defined *relative to a group*. The strings
`row group`, `selection context`, and `scoped` appear **zero** times in RDR 0006;
every `group` hit is `GroupUserEnv`/exit-code-group or the `flow` command group,
except `0006:252` ("single-valued grouping"), which is invariant 5's tag-class
grouping, not a row grouping. RDR 0002 is worse: `row group` and
`selection context` appear zero times there too, and RDR 0002 defines matching
as exactly-one over the **whole** table (`0002:312-315`, `:324-331`) with no
scoping notion at all.

The closest thing in either document is RDR 0006 invariant 4 (`0006:277-279`),
"for each state/outcome pair that claims closed coverage" — but that is a
predicate inside a numbered invariant list, not a definition: it never says how
a row is assigned to a pair, and RDR 0006's own Normative Contracts
(`0006:321-369`) never mention state/outcome pairs. RDR 0006 also never defines
"state" — its input contract (`0006:248-256`) lists rows, tags, the outcome
alphabet, terminal states, and escape rows, but never says how a row induces a
source or successor state. So the carrier for "state/outcome pair" is itself
undefined.

The draft half-notices this. Lines 376-383 observe that RDR 0006 "gates its
blocking finding on a contract 'that claims closed coverage', and no document
defines how a group makes that claim" — then resolves *how a group claims*
while never noticing that *what a group is* is equally undefined. Its own
normative clause at lines 462-466 then requires the undefined thing: "Coverage
and overlap checks MUST be scoped to a normalized row group **supplied by the
transition/lint model**." Neither named supplier supplies it.

The draft's own `trace` step 2 (line 568) walks this step and marks it **OK**:

> | 2. Group rows by selection context | … | `profile-to-grounding` and
> `foundational-to-cove` share `match status eq Draft` → one group | OK |

The witness *is the defect*. Those two rows are grouped because they share a
match predicate. But `continue-prelock-lenses` and `reconcile-rewind-legality`
have different match predicates and are placed in different groups. Nothing
states that a shared match block defines a group — the prose says
"source-state/recognized-outcome," and the fixture's rows carry no outcome at
all. The trace validated the grouping rule by applying an unstated rule to a
fixture and observing it produced a plausible answer. That is the definition of
a step that will diverge across implementations.

**Symptom.** Implementer A groups by `(match block, outcome)`. Implementer B
groups by outcome. On the same model, A reports `graph-coverage-gap` and B
reports green. The author files a bug against whichever they didn't expect, and
neither document adjudicates.

### 1.3 The default-on rule turns a missing declaration into a silent green

**Root cause.** The 3amigo pass closed P2-6 (an undefined trigger) by choosing
a default, and the default it chose fails open on exactly the input C-1 makes
universal.

**Enabling passage.** §Technical Design (lines 378-383):

> RDR 0006 gates its blocking finding on a contract "that claims closed
> coverage", and no document defines how a group makes that claim: this RDR
> reads it as **default-on for every scoped row group whose participating
> dimensions are all finitely declared** — **an author opts out by leaving a
> dimension undeclared**, not by omitting an annotation, since an opt-in flag
> would let the exhaustiveness guarantee be silently skipped exactly where it
> matters.

The reasoning is that opt-in "would let the guarantee be silently skipped." The
chosen alternative lets it be silently skipped *by the same mechanism*, with one
difference: an opt-in flag is a visible authored token a reviewer can grep for,
whereas an absent declaration is invisible. The clause optimizes for the wrong
failure. And note the clause's precondition — "whose participating dimensions
are all finitely declared" — is never satisfiable under C-1.

There is a second-order problem. The rule makes *forgetting* and *deciding*
produce identical artifacts. An author who deliberately leaves `notes` free-form
and an author who forgot to bound `iteration_count` both get the same silent
non-proof. §Risks (lines 903-905) anticipates the ergonomic half of this —
"Finite-domain declarations feel like boilerplate. **Mitigation**: … make
diagnostics explain when a missing domain blocks proof" — but no normative
clause requires that diagnostic. The `disposition` table row for "Guard
dimension lacks a finite declared domain" (line 554) marks it **Loud** with an
"Inability-to-prove finding," which contradicts the default-on paragraph: if an
undeclared dimension is the *opt-out mechanism*, it cannot also mint a blocking
finding, or opting out is impossible. The draft states both.

**Symptom.** Either every model with any undeclared dimension fails to lint
(making the feature unusable, since C-1 means that is every model), or every
such group silently skips its proof (making green meaningless). The
implementer picks one; the author gets the other.

---

## 2. The one section rewritten within six weeks of shipping

**§Technical Design, lines 345-418 — specifically the three paragraphs at
374-418 that define grouping, the default-on claim rule, and the narrowing.**

Iteration 1 already predicted this section would be rewritten, and it was
partially rewritten — that is the point. It was rewritten *without* fixing what
makes it unstable, and it acquired two new unstable clauses in the process
(default-on at 378-383; the narrowing's recording rationale at 408-418).

It will be rewritten again because it is the only place where four undefined
inputs converge and it has no authority over any of them:

1. **Grouping** (C-2) is delegated to a document that does not define it. The
   first implementer must invent the rule, and the invented rule lands here.
2. **The claim trigger** (C-3) resolves a peer's ambiguity by unilateral
   reading, with the alternative explicitly deferred: "If RDR 0006 intends an
   explicit per-group annotation instead, that is a divergence to settle at
   cluster reconcile alongside A8" (lines 381-383). A clause that names its own
   pending reversal is a clause scheduled for rewriting.
3. **The narrowing's recording document** (A8) is open by the RDR's own
   admission, and A8's Plan (lines 226-232) says both arms rewrite this text —
   either RDR 0006 cites this clause, or "this RDR's clause above collapses to
   a citation."
4. **"Can refuse"** (C-4) rests on `reachable predecessor`, undefined
   everywhere.

Second candidate, close behind: **§Critical Assumptions A1/A2/A6**, which will
be rewritten when the first implementer opens `check.sh` (C-6) or greps RDR 0002
for a domain field (C-1). But Technical Design is where the rewrite *lands*,
because that is where the semantics live.

---

## 3. The one assumption that will not survive first contact with a real user

**A2 — "Every exhaustiveness claim can be reduced to scoped finite declared
domains." Status: `Verified`. Method: Derivation.**

Iteration 1 also named A2, and the fix addressed the *mathematics* — it added
scoping, the symbolic/bitset representation requirement, and the
refuse-or-downgrade valve. The math is now correct. That is precisely why it
will fail: A2 is a valid derivation over inputs that do not exist, and a
`Verified` derivation is the strongest possible signal to an implementer that
the inputs are settled.

The word doing the damage is **`declared`**. A2 says domains are "declared" four
times. Nothing declares them. The Method is `Derivation` — "pure math or proof,
Evidence: the derivation, shown inline" (lines 269-271) — and a derivation is
the one Method that cannot detect a missing producer, because it never leaves
the document. A2 would be identically `Verified` if RDR 0002 did not exist.

First contact goes like this. A flow author writes their first real transition
model — the RDR prelock flow, the one the fixture gestures at. They declare
tags per RDR 0002: name, provenance, value kind. They author `all`/`unless`
guards. They run `intrastate lint`. It passes. They ask the obvious question —
"did it actually prove my profile routing is complete?" — and the honest answer
is no, it could not have, because they never declared what values `profile` can
take, and RDR 0002 gave them no place to say so.

The failure branch A2 states ("Lint may falsely claim guard coverage or miss
legal gaps in cap/profile/lens routing," lines 136-137) describes exactly this
and is booked as a hypothetical. It is the current state.

Runners-up, both of which fail for the same structural reason: **A6** (C-5) is
`Verified` on evidence for a neighboring claim, and **A1** (C-6) is `Verified` on
a spike that is a regex over operator spellings. All three share one pathology —
the Evidence line is *true*, and does not support the assumption.

---

## 4. Premortem

*Written from eight weeks after ship.*

We shipped guard predicate exhaustiveness behind `intrastate lint`, and the
lint has never once proved anything. It took six weeks to notice, because it
never failed.

The first flow author was Dana, encoding the RDR prelock path — the same journey
the fixture models. Dana wrote tag declarations exactly as RDR 0002 specifies:
`[tags.profile] provenance = "owned"`, `kind = "enum"`. She wrote
`[rule.guard.all.profile] in = ["mid", "large"]` on one row and
`eq = "foundational"` on another, and expected lint to tell her that `small` was
uncovered. She ran `intrastate lint`. Exit 0, no findings.

The implementer, Marco, had built `linter.groupRows` and `linter.provable` from
§Technical Design. `provable` asks each participating dimension for its declared
domain. `normalizer.TagDecl` — which Marco built from RDR 0002's schema, the
only normative source for it — has `Name`, `Provenance`, `Kind`, `Accessor`. No
domain field. So `provable` returned false for every dimension of every model,
and under the default-on rule at lines 378-383 ("an author opts out by leaving a
dimension undeclared") every group was correctly treated as opted out. Marco had
implemented the RDR exactly. `graph-coverage-gap` had never fired in production;
its only invocations were in `TestCoverageGap`, where the test fixture carried a
`domain` key that `TagDecl` ignored and the test asserted against a hand-built
domain map injected past the normalizer. The test proved the prover; nothing
proved the pipe.

Dana shipped. Three weeks later, `resolve.Resolve` returned `no_match` for a
`small`-profile RDR in the prelock stage. `escapeOrRefuse` found no modeled
escape for `no_match`, so the refusal reached the CLI as `respond.Fail` with a
`no_match` envelope, exit 2, and Dana's skill halted mid-flow with a message
naming nothing she could act on. She reopened her table, counted rows, and could
not see the gap — it is only visible in the `profile × prelock_iterations`
product, which is the one thing lint was built to see and had never looked at.

We fixed it by adding a `domain` field to `TagDecl`. That was the second
failure. Marco added it to the normalizer; Priya, implementing RDR 0006, had
already built `lint.Group` grouping rows by recognized outcome, reading
§Technical Design's "one source-state/recognized-outcome selection group."
Marco's `linter.groupRows` grouped by normalized match block, following the
draft's own `trace` step 2, which grouped `profile-to-grounding` and
`foundational-to-cove` because "both share `match status eq Draft`". The two
groupings disagree on every model with shared match contexts — which RDR 0002
item 4 encourages. On Dana's table, Marco's grouping reported an overlap that
Priya's did not. Neither could point at a document. We spent a week discovering
that RDR 0003 says RDR 0006 owns grouping and RDR 0006 says RDR 0003 owns the
proof, and that "row group" appears in RDR 0006 exactly three times, never as a
definition.

The third failure was the escape rows. Dana's table had an escape row for
`no_match` — added after the incident — carrying `unless status eq "Implemented"`.
`normalizer` emitted it as a candidate row with its full predicate set, per RDR
0002's escape clause. `linter.groupRows` had no rule excluding it, because RDR
0003 does not mention escape rows anywhere: four occurrences of the word, none
about rows. So the escape row joined the product and *supplied the coverage* for
the very `small` combination Dana had been missing. Lint went green. At runtime,
`resolve.Resolve` never consults escape rows until after the ordinary count is
zero (`escapeOrRefuse`, `resolve.go:473`), so the "covered" combination still
produced a refusal first. Lint proved coverage over a row set the resolver does
not use for coverage. That is the exact class of failure the §JD-4 narrowing was
written to prevent, and the narrowing did not catch it because the narrowing only
considers `guard_unevaluable`, not escape-edge ordering.

The fourth failure was quieter. Marco implemented the withheld-claim clause
(lines 475-482) and needed to decide "can refuse" — whether a key is "declared
present-for-every-reachable-predecessor." He searched for the definition of
reachable predecessor in RDR 0003 and RDR 0006 and found four uses and zero
definitions. He implemented it as "the key is written by some row whose
`NextTags` match this row's `Match`." Priya, implementing `graph-owned-before-write`
from the same undefined term, implemented a transitive closure over the whole
graph. On Dana's table with a rewind edge, Marco's version said a key was always
present and Priya's said it was not. Same model, two reachability answers, two
verdicts, no adjudicating text.

The fifth failure was found by a reviewer, not a user, and is the one that
would have cost the most. Marco had built the atom type from §Technical Design
line 354: "four conceptual fields: tag name, operator, expected value, and
source identity." RDR 0007's SEAM clause fixes four fields too — `Key`,
`Operator`, `Literal`, `Block`. Marco's type had no `Block`. Every `unless`
atom arrived at `evaluateGuard` indistinguishable from an `all` atom, so
`unless prelock_iterations gte 3` — the cap-3 guard, on the fixture's own
representative row — was evaluated as a positive requirement. The row fired
only when iterations *had* reached 3, the exact inverse of what Dana wrote.
It was caught in review because the RDR's own Load-Bearing Decisions › Identity
names `block` in the identity tuple, and the reviewer noticed the two sections
disagreed. Nothing in the pipeline would have caught it: the paragraph says
"this RDR cites it rather than restating it" immediately after restating it
wrong, and a claim to be citing is the strongest signal not to re-check.

The postmortem question was: how did this pass an RDR review that ran three
pre-lock lenses across two iterations? The answer is in §Prerequisites line 932:
`- [x] A1-A6 verified`. A2 was `Verified` by derivation, and the derivation was
correct. Nobody in either iteration opened RDR 0002 and read what a tag
declaration actually contains. Two lenses opened RDR 0002 — grounding iter-2 and
3amigo iter-2 both quote `0002:217-218` verbatim — and both used it to confirm
the *absence of an optionality field*, which is A7. Neither noticed that the same
line shows the absence of every domain field. They were looking for the subfield
and did not see that the field was missing.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

These are RDR-review-time tests: each is answerable by reading documents, before
any code exists. Each names the finding it catches.

### AT-1 — Every input A2 consumes has a named producer clause (C-1)

```gherkin
Given RDR 0003's A2 Evidence enumerates the inputs its proof consumes
  And those inputs include "finite domains for the guard dimensions" with
      enum/boolean declared sets, set element universes, and bounded integer
      {min..max}
When each input is traced to the normative clause in the producing RDR that
     obliges an author to supply it
Then every input resolves to a quoted normative clause in RDR 0002
  And no input resolves only to a spike fixture file or to RDR 0003's own text
```

**Fails today.** The finite-domain input resolves to
`evidence/spikes/guard-fixture.toml`, which invents `domain = [...]`, `min`, and
`max`. RDR 0002's tag declaration (`0002:217-218`) obliges name, provenance,
value kind, and accessor reference. A2 would have been `Pending`, not `Verified`,
and the missing producer request would have been booked alongside A7 and A9.

### AT-2 — No input is delegated in a cycle (C-2)

```gherkin
Given RDR 0003 delegates row grouping to RDR 0006
When RDR 0006 is searched for a definition of row group or selection context
Then RDR 0006 contains a normative clause defining how candidate rows are
     partitioned into groups
  And RDR 0006 does not, for the same computation, delegate back to RDR 0003
```

**Fails today.** RDR 0006 has no grouping definition, and its A2 (`0006:111-119`)
delegates the coverage/overlap derivation to RDR 0003 A2. Each names the other.

### AT-3 — The claim trigger cannot be satisfied by silence (C-3)

```gherkin
Given a row group with one guard dimension whose domain is not declared
When lint evaluates exhaustiveness for that group
Then the outcome is stated by exactly one clause in RDR 0003
  And that clause is not contradicted by another clause in the same document
```

**Fails today.** §Technical Design lines 378-383 make an undeclared dimension the
*opt-out* (no finding). The `disposition` table line 554 makes the same input a
**Loud** blocking "Inability-to-prove finding." Two clauses, opposite outcomes,
same input.

### AT-4 — Every term a normative clause decides on is defined (C-4)

```gherkin
Given the normative clause "A row 'can refuse' when the row carries a value
      atom whose key is not declared present-for-every-reachable-predecessor"
When "reachable predecessor" is traced to its definition
Then exactly one document defines the predecessor relation, its domain (states
     or rows), and how reachability is computed
```

**Fails today.** Four uses in RDR 0003 (`:192`, `:416`, `:497`, `:517`), two in
RDR 0006 (`:282-283`, `:94`), zero definitions. The 3amigo pass cited
"line 512, RDR 0006 via A6" as the existing notion; RDR 0003 line 512 is itself
an undefined use.

### AT-5 — Evidence supports the assumption's own "If wrong" (C-5)

```gherkin
Given A6 states "If wrong: it cannot prove that owned tags are set before they
      are matched"
When A6's Evidence is checked against that specific consequence
Then the Evidence establishes that write-reachability is decidable
  And not merely that provenance labels are declared
```

**Fails today.** A6's Evidence establishes that RDR 0002 requires provenance
declaration. Provenance labels do not make reachability decidable — that needs
the graph relation AT-4 shows is undefined. This test is the RDR's own rule at
lines 288-291 applied mechanically to each `Verified` record.

### AT-6 — A Spike Method executes the semantics it claims to verify (C-6)

```gherkin
Given A1 uses Method "Spike" to verify that target flows "fit a closed typed
      predicate vocabulary"
When evidence/spikes/check.sh is read
Then the script evaluates at least one predicate against at least one tag
     assignment
  And type-checks at least one literal against its tag's declared kind
  And its captured output shows a semantic verdict, not an operator inventory
```

**Fails today.** `check.sh` is `awk` regex-matching operator tokens against a
fixed list; `output.txt` reports `operators=eq,in,lt,gte,exists` and `rules=4`.
A1 would have been `Pending`, and the "closed vocabulary is sufficient" claim
would have carried its actual (zero) evidential weight.

### AT-7 — Producer normalization preserves every field the identity tuple names (C-7)

```gherkin
Given RDR 0003's atom identity is the tuple
      (RuleID, SourceLocator, key, block, operator, literal)
  And RDR 0003's unless semantics require the block to remain a single
      conjunctive exclusion set
When RDR 0002's normalization contract is read
Then normalized candidate rows retain the block discriminator per atom
```

**Fails today.** `0002:305-309`: "Normalization MUST combine both into one
candidate-row predicate set before ambiguity checks." One set. The block is not
required to survive, so both the identity tuple and `unless`-as-block-exclusion
lose their producer guarantee.

### AT-8 — Every normalized row kind has a stated disposition (C-8)

```gherkin
Given RDR 0002 normalizes escape rules into candidate rows carrying their
      normal predicate set
  And resolve.go::escapeOrRefuse runs escape rows through the same guard gate
When RDR 0003's disposition table and normative contracts are searched for
     escape rows
Then a clause states whether escape rows join their group's product, form a
     separate group, or are excluded from exhaustiveness proofs
```

**Fails today.** `grep escape` over RDR 0003 returns four hits, none about rows.
The ten-row `disposition` table has no escape row. This test also catches the
premortem's third failure — an escape edge supplying coverage the resolver will
not honor.

### AT-9 — A clause binding a peer's artifact is registered as a peer obligation (C-9)

```gherkin
Given RDR 0003 requires the withheld-claim finding to "name the participating
      row and the atom that can refuse" using RDR 0006's
      graph-unprovable-coverage code
When RDR 0006's finding payload contract is read
Then that contract carries an atom-level identity field
  Or RDR 0003 registers the payload extension as a Pending peer obligation
```

**Fails today.** `0006:352-356` requires "stable code, model identity, severity,
human-readable message, and the source rule/context id or source span" — no
atom field. A8 covers the narrowing's recording document, not this payload
extension, so the obligation is unregistered.

### AT-10 — Each `[x]` prerequisite is re-tested, not re-asserted (C-10)

```gherkin
Given §Prerequisites carries "- [x] A1-A6 verified"
When each of A1 through A6 is independently re-checked against its cited source
Then each Evidence line is confirmed to support that assumption's specific claim
  And the §Assumption Verification response cites the re-check, not the prior
      Verified stamp
```

**Fails today.** §Assumption Verification re-states A1/A2/A6 as sound in prose
(label hygiene: allowed Method, no self-reference, no `Docs Only`) without
re-testing any Evidence line. This is the gate that let C-1, C-5, and C-6 pass
two iterations.

### AT-11 — Every MVV case is authorable in the specified source format (C-11)

```gherkin
Given the MVV requires a row group with a possibly-absent guard key, a negative
      control over always-present keys, and a contains predicate over a declared
      set-valued tag
When each is checked against RDR 0002's source schema
Then each required fixture element can be expressed in the specified TOML layout
```

**Fails today.** Three of the MVV's required cases need declarations RDR 0002
does not carry (optionality — A7; element universe — A9; and per C-1, domains at
all). §Scope Verification nonetheless asserts the MVV "is in scope for
implementation, not deferred."

### AT-12 — A "cite, don't restate" claim is checked against the cited text (C-13)

```gherkin
Given §Technical Design states the atom has four fields — tag name, operator,
      expected value, and source identity
  And the same paragraph claims the shape is fixed by JDR 0001 §D1 and cited,
      not restated
When RDR 0007's SEAM clause is read
Then the four field names in RDR 0003 match the four field names in RDR 0007
```

**Fails today.** `0007:1260-1265` fixes `Key`, `Operator`, `Literal`, `Block`
(with `Block` an exported type carrying `BlockAll`/`BlockUnless`). RDR 0003
substitutes `source identity` for `Block`. A sentence claiming to cite is the
strongest possible signal *not* to re-check — which is why the substitution
survived two lenses. This test also catches the internal contradiction with the
identity tuple at lines 584-592, which does name `block`.

### AT-13 — Every obligation a peer files on this RDR is booked here (C-14)

```gherkin
Given RDR 0007 is Final and names obligations it assigns to RDR 0003
When each such obligation is traced into RDR 0003
Then each appears as a Critical Assumption, a Capability Dependency row, or a
     normative clause in RDR 0003
```

**Fails today.** `0007:1585-1595` requires RDR 0003 to declare a canonical
byte-spelling for set-valued literals (`in` and `contains`), because without it
0007's own payload sort is non-total. RDR 0003 books only the *element universe*
(A9) and only for `contains`. The literal-encoding duty for `in` — an operator
in the target-flow subset the spike covered and the MVV exercises — is
unbooked. RDR 0007 Phase 4 (`0007:2169-2183`) hands over seven items; RDR 0003
A7 acknowledges "six sibling items" (line 177).

---

## Closing

The two lenses that ran this iteration did careful work on the text in front of
them. Both opened RDR 0002 and quoted `0002:217-218` verbatim. Both read that
line to establish a missing *subfield* — optionality for A7, element universe for
A9 — and neither registered that the same line establishes there is no domain
declaration for the proof to consume at all. The 3amigo pass likewise opened
RDR 0007's SEAM clause and quoted its four atom fields (`0007:1249-1266`) to
prove there is no *position* field, without noticing that the four fields it had
just quoted are not the four fields §Technical Design lists (C-13).

That is the shape of the miss in both cases: the lens went to the right passage,
extracted the one fact it was looking for, and did not read the passage against
the draft's other claims. Precision about the question asked; no sweep for the
question not asked.

The result is a draft that now argues with precision about the recording document
for a narrowing rule (A8), the projection semantics of one operator (A7), and the
acceptance gate for another (A9) — three genuinely open questions, all correctly
booked — while its central `Verified` derivation reduces coverage to domains that
nothing in the system can declare, its atom field list contradicts the seam it
claims to cite, and its coverage identity omits a row class the consuming RDR
counts.

A7 and A8 are named as the two blockers to lock. They are not the binding
constraints. C-1 is: until RDR 0002 carries a finite-domain declaration, this
RDR specifies a proof with no inputs, and A2's `Verified` stamp is the reason
nobody has looked. C-13 is the one that ships a working system computing the
inverse of what authors wrote.

Recommended disposition: A1, A2, and A6 revert to `Pending`; C-1's
finite-domain declaration joins A7/A9 as a single RDR 0002 producer request;
C-2's grouping definition and C-13's field-list correction are both closable
inside this document today and should not wait for cluster reconcile.
