Model: claude-sonnet-5

## Findings Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | §JD-14 / A17 / `authority` table row "Escape-row overlap population" / Finalization Gate lines 1897-1907 | This RDR unilaterally reverses a JDR-recorded joint decision (JDR 0001 §JD-14, "Decided 2026-08-22: 0003's reading governs" — single-population overlap, escape rows never exempt) without amending the JDR or RDR 0006, both of which still state the opposite of what this RDR now asserts | Two implementers reading different documents in the cluster build mutually contradictory overlap lint — one blocks escape/guarded overlaps, one doesn't; a fixture flagged blocking under one reading passes clean under the other | §1, premortem, AT-1 |
| C-2 | A10 ("Pending"), A12 ("Pending") — row-group definition confirmation and predecessor-reachability contract, both homed at "RDR 0006's refine" | The RDR locks (or is on a path to lock) two of its own load-bearing derivations — the very unit coverage/overlap is computed over, and the decision procedure for the owned-tag clause — on a Draft peer document's *future* refine pass that has not started, with no forcing mechanism if RDR 0006's refine disagrees | Implementation of Phase 2 (finite-domain lint semantics) either stalls waiting on RDR 0006, or proceeds on this RDR's unconfirmed row-group definition and later has to be reworked when RDR 0006's refine defines it differently | §1, §2, premortem, AT-2 |
| C-3 | A14 ("Pending") — per-atom `Block` retention through normalization, requested from RDR 0002's refine; Normative Contracts identity-tuple clause (`(RuleID, SourceLocator, key, block, operator, literal)`) and `unless`-is-not-per-atom-negation clause both depend on `Block` surviving normalization | The entire `unless` semantics — and therefore every overlap/coverage proof involving a mixed `all`/`unless` row — depends on a normalization guarantee RDR 0002 has not yet made normative (RDR 0002:897-899 lists it as a to-do, not a contract); if the implemented normalizer collapses or drops `block`, guard exclusion silently becomes guard inclusion or vice versa | An author writes `unless profile eq "small"` expecting exclusion; if `Block` doesn't survive normalization as designed, the row is misclassified and either wrongly excludes valid transitions or wrongly admits excluded ones — with no compiler/lint error, only a wrong resolution or a wrong exhaustiveness verdict | §1, §2, premortem, AT-3 |
| C-4 | `disposition` table row "Value atom over an absent key" and shipped `internal/resolve/resolve.go::gate` lines 408-417 | The RDR's diagnostic promise ("names the row and the refusing atom", Normative Contracts "withheld claim MUST be observable... MUST name the participating row and the atom that can refuse") is written against an atom-slice `Row.Guard` shape that does not exist yet; the shipped kernel's `Refusal.Guard` is a single opaque `string` (`resolve.go:270-273`) and `gate()` picks the row with the lexicographically lowest `RowRef` among all undecidable rows (`slices.MinFunc`, line 408-411) as "the" refusing guard, discarding the others' guard strings entirely | When two rows are simultaneously `guard_unevaluable`, the CLI error names only one row's guard string (whichever sorts first by `RuleID`/`SourceLocator`), not "every refusing row" the RDR's own normative clause (`disposition` table, "one finding per refusing row, never just the first") requires — production behavior contradicts the RDR's own MVV Scenario 8 assertion until the RDR 0007 kernel reshape ships, which this RDR admits (Phase 1 gating table) has no start date | §1, premortem, AT-4 |
| C-5 | RDR 0006's `disposition` for "claims closed coverage" — "default-on for every scoped row group whose participating dimensions are all finitely declared" (Normative Contracts, `row group` block) | Default-on exhaustiveness checking means every row group with only finite dimensions is silently opted into a blocking lint gate the moment it becomes fully finite — an author who declares one more tag's domain (for an unrelated reason, e.g. documentation) can turn an existing, working, previously-unchecked transition table into a blocking lint failure without touching any guard or transition rule | A routine "let's document this tag's allowed values" PR breaks CI on an unrelated part of the transition table, because the newly-finite dimension makes a previously-unprovable row group provable — and not-exhaustive | §2, premortem, AT-5 |
| C-6 | Normative Contracts, "participation clause": "A dimension participates in a row group when any row in that group carries a guard atom over that key... A key every row constrains identically still bounds the product and still requires a finite declared domain" | Authors get a blocking lint failure for adding a *harmless, universally-true* guard atom (e.g., a defensive `all` clause every row in the group happens to share) purely because that key now "participates" and must have a finite domain declared — even though the atom contributes zero discriminating power and the group was previously exhaustive without it | An author adds `all: env eq "prod"` to every row in a group as a defensive/documentation gesture (all rows share it, so it changes nothing about which row is selected) and lint newly demands `env`'s domain be finite-declared, breaking a build that had nothing wrong with it before | §3, premortem, AT-6 |
| C-7 | A15 ("Pending", MVV Test) — the published cardinality bound is a single scalar, but the RDR's own evidence line admits "the proof cost of a symbolic or bitset representation may depend on dimension count or per-dimension width, not on cardinality alone" | The bound clause the RDR ships with is provisionally accepted on a test (MVV Scenario 3) that has not run — if the comparison of equal-cardinality/different-shape products diverges, "the bound needs a second term and this clause is restated before lock" (A15's own text) — i.e., the RDR's own author expects this clause may not survive its own validation | An author with a wide-but-low-cardinality guard set gets an inconsistent "too large to prove" verdict compared to a narrow-but-high-cardinality set of the same nominal size — the same size and shape and a different bound gate result depending on which implementation choice (bitset width vs. symbolic term count) actually gated the refusal | §2, §3, premortem, AT-7 |
| C-8 | Normative Contracts, optionality clause: "A declaration carrying no optionality marker declares the key optional — the conservative default" | Every tag an author forgets to mark `optional = false` silently becomes optional, which (a) adds a presence dimension to every row group referencing it, doubling/complicating the scoped product, and (b) licenses every guard row over that key to legally refuse `guard_unevaluable` at runtime — a correctness-relevant default that is invisible unless the author reads this RDR's prose, since RDR 0002's authoring surface has no normative requirement that `optional` be explicit | An author writes a tag declaration, forgets `optional = false` (an easy omission — the illustrative TOML block, lines 1219-1231, shows it explicitly set, but nothing requires it), and later gets confused why a row group that "should be" exhaustive is silently blocked pending A10/A12 lint support, or why resolution occasionally refuses `guard_unevaluable` for a key that is, in practice, always present | §3, premortem, AT-8 |
| C-9 | Metadata block, lines 9-25: five re-verification passes, four assumption status changes, and a rehomed model — all inside one iteration-3 re-entry of a Draft RDR whose Predecessors (RDR 0001) has already shipped and whose siblings (0002, 0006) remain Draft | The RDR's own churn rate — declaration model moved in (2026-08-21), narrowing moved in via cluster gate (2026-08-22), A16/A17 opened same day — signals the document is still actively finding its true shape at lock time, which is exactly the profile of a "large" RDR that gets rewritten again within weeks of shipping once real fixture pressure (not the hand-picked `guard-fixture.toml`) hits it | Six weeks after implementation starts, engineers file a follow-up RDR or a same-document Draft-reopening for whichever of A10/A12/A16/A17 landed wrong, because the "standing tolerance" mechanism let this RDR lock without those questions actually being settled | §2, premortem |
| C-10 | Illustrative Code block (lines 1219-1231) and MVV "Authorability" paragraph — the illustrative TOML uses `single_valued = true` as a tag field, explicitly flagged "illustrative only... the field's authoring location is RDR 0002's to fix (A16)" | The RDR's own worked example is not authorable today — a new engineer reading this RDR to understand the syntax will write exactly the TOML shown and it will not parse against RDR 0002's shipped/current schema (which enumerates only four fields), producing a confusing first-contact failure that has nothing to do with guard semantics | A developer copies the illustrative declaration block to author their first guarded transition, gets an "unknown field: single_valued" (or silently-ignored-field) error from the RDR 0002 loader, and has no way to know from RDR 0003 alone that the fix is pending in a different document | §3, premortem |

## 1. The three most likely ways implementation goes wrong

### 1.1 The escape-row overlap contradiction ships inconsistently across the cluster (C-1)

**Root cause.** This RDR's A17 and its two-population overlap clause assert that escape rows are checked for overlap only against other escape rows, never against guarded rows — directly contradicting JDR 0001 §JD-14, which the JDR itself still states as **"Decided 2026-08-22: 0003's reading governs — escape rows are ordinary participants in **both** the union and the overlap check"** (`docs/jdr/0001-resolve-kernel-seam.md:333-336`), and directly contradicting RDR 0006's cluster-reconcile disposition, which likewise records **"Decided: RDR 0003's reading governs; repair invariants 3 and 4 to match the Load-Bearing Decision"** (`docs/rdr/0006-graph-lint-authority-and-guarantees.md:911-912`) — meaning RDR 0006, too, currently commits to the single-population reading this draft now reverses.

**The specific passage.** The RDR's own `authority` table calls this out explicitly as **"Contested — A17"** (line 1100) and its Finalization Gate section admits: "One gate decision is partially re-opened, deliberately and on evidence" (line 1897). But "re-opened" is doing a lot of unearned work: nothing in JDR 0001 or RDR 0006 has actually been amended. The RDR is not describing a re-opened joint decision; it is describing itself *disagreeing* with a joint decision that, as recorded in the canonical joint-decision document, is still closed the other way. The user's own global instruction is unambiguous here: RDRs are not amended, and JDRs are presumably the analogous record-of-truth for joint decisions — yet this draft treats a closed JDR entry as overridable by unilateral assertion inside a peer RDR, with the reversal recorded only in RDR 0003's own prose (A17), not in the JDR.

**The symptom.** Two implementers — one building RDR 0006's graph lint from RDR 0006's text (which still says "repair invariants 3 and 4 to match... escape rows are ordinary participants... in the overlap check"), one building from this RDR's text (which says the opposite) — ship contradictory blocking-lint behavior. A fixture with an escape row overlapping a guarded row is a blocking `graph-overlap` finding under one reading and clean under the other. This is precisely the "single-source failure" A8/§JD-4 exists to prevent, reproduced by this RDR one clause later, and the irony is sharp: the RDR spends a full normative clause plus an `authority` table row explaining why cross-document agreement matters (A8), then produces exactly the disagreement it warned against, three sections later, in the same document.

### 1.2 The row-group and predecessor-reachability contracts are locked on a peer's unstarted work (C-2)

**Root cause.** Two structurally central concepts this RDR relies on — the row-group grouping predicate that scopes every coverage/overlap product (A10), and the reachability relation that makes the owned-tag clause ("every reachable predecessor sets or preserves that tag") decidable (A12) — are stamped `Pending`, with the resolution plan being "RDR 0006's refine confirms it reads the same [thing]." Both are marked "Not lock-blocking" via the "standing tolerance" mechanism.

**The specific passage.** A10's evidence line admits the underlying problem plainly: "This RDR's Technical Design previously took the grouping context *from* RDR 0006, while RDR 0006's A2 takes the coverage/overlap derivation *from* this RDR — a cycle with no floor." The "fix" is that this RDR now defines the row group unilaterally (Normative Contracts, "The row group is defined here"), but the done-condition for A10 is still "RDR 0006 confirming it reads the same division of labour" — i.e., the cycle is broken by fiat on one side, with the other side's agreement deferred to a document that is itself `Draft` and has not yet run its refine pass. Similarly A12's "Plan" line requires RDR 0006 to "publish a citable predecessor relation over selection contexts and state which decision procedure governs" — and flags that RDR 0006 currently frames the check as *graph reachability* while this RDR requires a *syntactic decision over declarations*, "which are different procedures over the same predicate." That is not a minor terminology gap; a graph-reachability decision procedure and a syntactic-declaration decision procedure can produce different verdicts on the same table, and nothing in either document commits to reconciling which one wins if they diverge.

**The symptom.** Implementation of Phase 2 (Finite-Domain Lint Semantics) either (a) stalls, blocked on RDR 0006's refine landing first, delaying this RDR's own Phase 3/4 despite the phase-gating table's claim that Phase 3 is "Startable today," since Phase 3's target-flow fixture needs group construction to test exhaustiveness at all; or (b) proceeds against this RDR's unilateral row-group definition, and when RDR 0006's refine lands with a different (e.g., graph-reachability-based) grouping, the lint code built in Phase 2 has to be reworked — silently reproducing the same "circularity" the RDR spent two review iterations closing for A2/A7/A9/A11.

### 1.3 `unless`-block semantics silently break because per-atom `Block` retention is not yet a normalization contract (C-3)

**Root cause.** The RDR's entire negative-guard model — the identity tuple `(RuleID, SourceLocator, key, block, operator, literal)`, the clause that `unless` denotes "a conjunctive exclusion set" rather than per-atom negation, and MVV Scenario 7 ("block retention") — depends on `Block` surviving RDR 0002's sparse-to-normalized transform as a per-atom field. But A14 is `Pending`, and the evidence line is explicit: "RDR 0002 requires normalization to 'combine both into one candidate-row predicate set' and guarantees only that a row retains its rule id and source locator" — nothing about per-atom block survival. Grounded directly: RDR 0002's own refine backlog (line 897-899) lists this as an open item ("RDR 0003 A14 — per-atom `block` retention. Normalization currently promises only that a candidate row retains rule id and locator. State that each atom also retains the block it was authored in") — i.e., RDR 0002 agrees this is missing and has not yet fixed it.

**The specific passage.** The RDR's own "If wrong" line for A14 states the failure mode precisely: "`unless` collapses into `all` after normalization, the identity tuple loses a component that makes it total, and a row's excluded intersection can no longer be subtracted from its accepted assignments." That is not a cosmetic bug — it is a semantic inversion. A guard exclusion that collapses into inclusion (or vice versa) changes which transitions are legal, silently.

**The symptom.** MVV Scenario 7 is supposed to catch this ("Normalize a row carrying atoms over the same key in both blocks... assert each atom still reports the block it was authored in"). But that test only exists in the MVV — it has not run, because RDR 0002's normalizer (and its `Block`-retention contract) do not exist yet. If the implementer builds RDR 0002's normalizer first without reading this RDR's A14 closely (a real risk since A14 is filed as "Not lock-blocking" and buried in a Pending-assumption list), the first place the defect surfaces is a wrong resolution or wrong exhaustiveness verdict in production data, not a compile-time or lint-time error — because nothing forces the `Block`-retention test to run before the normalizer ships.

## 2. The one section that will be rewritten within 6 weeks of shipping

**RDR 0006's "row group" definition and predecessor-reachability contract — and by extension, this RDR's Normative Contracts clauses that assume A10 and A12 close cleanly (the "row group" clause, the owned-tag clause, and the participation clause).**

This is the highest-confidence prediction because the RDR's own text tells you where the unresolved tension is, twice: once in A10's evidence ("a cycle with no floor"), and once in A12's evidence (two documents proposing genuinely different decision procedures — "graph reachability" vs. "syntactic decision over declarations" — over the same predicate, with no arbiter named). Both are deferred to "RDR 0006's refine" as a "standing tolerance," which sounds procedural but is actually a bet: the bet is that when RDR 0006's Draft author sits down to write its refine, they will (a) adopt exactly this RDR's row-group definition without modification, and (b) either adopt the "syntactic decision over declarations" reading for the owned-tag clause or produce a reachability contract this RDR can cite without contradiction.

Given that this exact cluster has *already* produced one closed-then-silently-reversed decision within the current review pass (the escape-row overlap contradiction, C-1) — meaning the peer RDRs in this cluster do not reliably converge even when a JDR explicitly decides the question — betting that RDR 0006's refine will converge cleanly on two more open questions (A10, A12) without a repeat of the same failure mode is optimistic. When RDR 0006's refine lands (likely within weeks of any implementation attempt, since Phase 2 of this RDR cannot meaningfully proceed without it), at least one of A10 or A12 will surface a mismatch that requires a follow-up correction — most likely the reachability decision-procedure question, since the RDR's own text flags it as "different procedures over the same predicate," which is a substantive disagreement, not a wording gap.

## 3. The one assumption that will not survive first contact with a real user

**A15 — that a single scalar cardinality bound predicts provability regardless of dimension shape.**

The RDR's own evidence line for A15 concedes the problem before a single real fixture is built: "the proof cost of a symbolic or bitset representation may depend on dimension count or per-dimension width, not on cardinality alone, in which case two products of equal cardinality could differ in provability and the bound would not predict the verdict it gates." This is not a hypothetical corner case an author is unlikely to hit — it is the *generic* case for any nontrivial guard model. Consider: a flow with two enum tags of 1000 values each (cardinality 1,000,000, low dimension count, "few wide dimensions") versus a flow with twenty boolean tags (cardinality 1,048,576, high dimension count, "many narrow ones"). These are exactly the shapes MVV Scenario 3 is designed to compare — and the RDR's own author has already written the contingency plan for when they diverge: "the bound needs a second term and this clause is restated before lock."

A real author will not construct these fixtures deliberately to probe the model's edges. They will simply write a transition table with a handful of enum tags and a handful of boolean flags, mixed — which is the overwhelmingly common shape for a real state machine's guard surface. The moment the implementation's chosen proof representation (a symbolic BDD-like structure, most likely, given "deterministic symbolic or bitset-equivalent representation") has a cost that depends on dimension count (which almost all practical symbolic representations do — this is precisely why BDD variable ordering is a classical hard problem), a user will hit a guard set that gets refused as "too large to prove" purely because it has many narrow dimensions, right next to a semantically simpler but nominally larger-cardinality guard set that proves fine. The user's mental model — "cardinality is the budget" — will be actively wrong, and the diagnostic (which reports cardinality and the bound, per the Normative Contracts clause) will mislead them into thinking they need to shrink their *domain*, when the actual fix is to shrink their *dimension count*, which the diagnostic does not tell them to do.

## 4. The premortem

It is six weeks after RDR 0003 shipped. Guard predicates are live in `internal/resolve` and `internal/lint` (or wherever RDR 0006's package lands). Three separate incidents have already been filed against this RDR's design, and a fourth is in progress.

**Incident 1 — the escape-overlap flip-flop.** The engineer who implemented RDR 0006's graph lint read RDR 0006's own text — which, per the `0003-0006-0007` cluster-reconcile disposition, still instructs "repair invariants 3 and 4 to match the Load-Bearing Decision" (single population, escape rows are ordinary overlap participants) — and shipped exactly that: `lint.checkOverlap` treats every row, escape or ordinary, as one population. A second engineer, working from RDR 0003's text (two populations, escape rows excluded from the ordinary check, A17), files a bug: "flow `cap-3-handling.toml` has an escape row that legitimately overlaps a guarded row by design (it's the rescue path), and lint now blocks it as `graph-overlap` — but resolve.go's runtime accepts this exact fixture and resolves it fine." The bug sits for a week because the two engineers each cite a different "final" document as ground truth, and nobody owns reconciling JDR 0001 (which says one thing), RDR 0006 (which says the same thing as the JDR), and RDR 0003 (which says the opposite) — because the mechanism that was supposed to prevent this (§JD-4's whole point, and A8's whole point) was itself bypassed by A17 without a JDR amendment. The eventual fix requires someone to re-open JDR 0001 §JD-14 formally — the very re-entry this project's workflow says RDRs don't get, but which a JDR apparently can, creating a precedent nobody had agreed to.

**Incident 2 — a `unless` regression from a "harmless" normalizer refactor.** RDR 0002's normalizer ships with `combine both into one candidate-row predicate set` implemented literally — i.e., `all` and `unless` atoms are merged into a single `[]Atom` slice with a `Block` field, because that's the shape RDR 0007 fixed and this RDR's identity tuple needs. Three weeks later, someone optimizing normalization performance flattens duplicate-key atoms across blocks as a memory optimization (an atom over `profile` in `all` and another over `profile` in `unless` get de-duplicated into one entry because the optimizer's dedup key was `(key, operator, literal)` — the *semantic* equality this RDR defines for domain computation, not the *identity* equality that requires `block` — and the two concepts got confused in the implementation because both are called "equality" in this RDR's own Load-Bearing Decisions section without enough separation in the code comments that got copied into the implementation). A `prelock-lens-routing` flow that used to correctly exclude rows once `prelock_iterations >= 3` (via `unless`) starts accepting them, because the `unless` atom silently vanished into the `all` set during the dedup pass. Nobody notices for eleven days because MVV Scenario 7 (block retention) tests the normalizer's *initial* output, not what a later performance pass does to it, and nothing in CI re-runs Scenario 7 against every code path that touches the atom slice.

**Incident 3 — a routine documentation change breaks CI.** An engineer adds `domain = ["prod", "staging", "dev"]` to the `environment` tag's declaration purely to document the allowed values for a wiki page generator RDR 0006 doesn't even use — no guard logic changes. Because RDR 0006's disposition table makes exhaustiveness "default-on for every scoped row group whose participating dimensions are all finitely declared," and `environment` participates (some row somewhere has `all: environment eq "prod"`), that row group — which was previously exempt from exhaustiveness checking because `environment` had no declared domain — is now checked, and it isn't exhaustive (nobody wrote a `staging`/`dev` branch because nobody thought they needed to). CI goes red on a docs-only PR. The engineer's first instinct is to revert the domain declaration, which "fixes" CI by making the model less precise, exactly the outcome A11's whole rehoming was meant to prevent.

**Incident 4 (in progress) — the cardinality-bound surprise.** A kata-authoring flow with twelve boolean-ish enum tags (each two or three values, `2-3^12 ≈ half a million to million range`) gets refused by the exhaustiveness prover as "too large to prove" even though a peer flow with two 800-value enum tags (cardinality 640,000, comparable size) proves fine. The author reads the diagnostic — cardinality and bound, per spec — shrinks their tag domains to make the number smaller, and it still fails, because the actual cost driver was dimension count (the symbolic prover's BDD-style representation scales with variable count, not raw cardinality), which the diagnostic never told them. They eventually give up on exhaustive lint for that flow and route around it with an escape row that swallows everything — the exact "sentinel-stamping anti-pattern" A7's own "If wrong" clause warned this design was supposed to prevent.

## 5. Acceptance tests that would have caught each failure at RDR-review time

**AT-1 (catches C-1, escape-overlap contradiction).**
```
Given JDR 0001 §JD-14 records "Decided: 0003's reading governs" for escape-row
  overlap, quoting a single-population reading
And RDR 0006's cluster-reconcile disposition independently records
  "Decided: RDR 0003's reading governs" citing the same single-population reading
When RDR 0003's own normative text is checked against both citations
Then RDR 0003 MUST NOT assert a two-population overlap reading
  without JDR 0001 §JD-14 itself being amended to record the reversal
And this RDR MUST NOT lock while any cited "Decided" JDR entry contradicts
  this RDR's own normative clause on the same question
```
This is a mechanical, greppable check: read every JDR entry this RDR cites as "Decided" or "Closed," and diff its stated substance against this RDR's own normative clauses. It would have caught C-1 in minutes — the contradiction is not subtle once the JDR and this RDR are read side by side, which is exactly what a pre-lock review pass should force.

**AT-2 (catches C-2, row-group/reachability locked on unstarted peer work).**
```
Given A10 and A12 are both status "Pending" with plan "confirmed at RDR 0006's refine"
And RDR 0006 is status "Draft" with no refine pass yet run
When this RDR's Phase 2 (Finite-Domain Lint Semantics) is scheduled
Then the Implementation Plan MUST NOT mark Phase 2 or Phase 3 as "startable"
  unless the row-group definition and predecessor-reachability contract used
  in Phase 2's actual test fixtures are the SAME definition RDR 0006 later adopts
And the RDR MUST name a fallback: what happens to Phase 2 work already built
  if RDR 0006's refine adopts a different row-group definition or a different
  reachability decision procedure than this RDR assumed
```
The phase-gating table already flags A10/A17 as blocking group construction ("should wait on A10/A17 rather than be built twice") — this is close to being caught by the RDR's own text, but no acceptance test forces the "built twice" risk to be resolved before lock rather than accepted as a known cost.

**AT-3 (catches C-3, `Block` retention).**
```
Given A14 is "Pending" and RDR 0002's own refine backlog lists per-atom block
  retention as an open normative gap, not yet a contract
When a row carries atoms over the same key in both `all` and `unless`
And that row is normalized by RDR 0002's shipped normalizer
Then MVV Scenario 7 (block retention) MUST run against RDR 0002's actual
  shipped normalizer output, not a hand-encoded fixture, before this RDR's
  guard-evaluation code is allowed to consume normalized rows
And this MUST be a build-time or CI gate, not merely an MVV scenario that
  exists in a Validation section without an enforcement mechanism tying it
  to every code path that constructs or transforms an atom slice
```

**AT-4 (catches C-4, `guard_unevaluable` diagnostic understates refusing rows).**
```
Given the shipped kernel's gate() function (internal/resolve/resolve.go:408-417)
  reports only the lexicographically-lowest undecidable row's Guard string
  on a KindGuardUnevaluable refusal
And this RDR's Normative Contracts require "one finding per refusing row,
  never just the first"
When two or more rows in one selection context are simultaneously
  guard_unevaluable
Then the runtime refusal payload MUST name every refusing row and its atom,
  not the first by sort order
And this RDR's Finalization Gate MUST NOT claim consistency with the shipped
  kernel while the shipped kernel's actual refusal payload (Refusal.Guard,
  a single string) structurally cannot carry more than one row's guard
```
This is directly testable against the existing shipped code today — no new fixture is needed, only reading `gate()`'s `slices.MinFunc` call against the RDR's own "never just the first" clause.

**AT-5 (catches C-5, default-on exhaustiveness breaks unrelated changes).**
```
Given a row group is not currently exhaustiveness-checked because one
  participating dimension lacks a finite domain
When an author declares a finite domain for that dimension for a reason
  unrelated to guard logic (e.g. documentation, tooling)
Then lint MUST NOT silently begin blocking that row group without a
  migration path — e.g. a grace-period warning-first rollout, or a lint
  rule that flags "this domain declaration newly enables exhaustiveness
  checking on N existing row groups, M of which do not currently pass"
  BEFORE promoting the check to blocking
And the RDR MUST state what happens to CI on the commit that crosses this
  threshold
```

**AT-6 (catches C-6, harmless shared guard atom triggers a false blocking finding).**
```
Given a guard atom over key K is authored identically in every row of a
  group (K does not discriminate between rows)
And K has no currently-declared finite domain
When an author adds that atom purely as documentation/defense-in-depth,
  changing no row's actual selection behavior
Then lint MUST NOT newly require K's domain to be finite-declared as a
  side effect of an atom that contributes zero discriminating power
Or, if the participation clause is kept as specified, the RDR MUST warn
  explicitly in Risks/Consequences that adding *any* guard atom — even a
  no-op one — can convert a previously-clean row group into a blocking
  lint failure, so authors are warned before they hit it in production
```

**AT-7 (catches C-7, cardinality bound doesn't predict provability).**
```
Given A15's own evidence line predicts the bound "may not predict the
  verdict it gates" for products of equal cardinality but different shape
When MVV Scenario 3's shape-comparison test is run
Then it MUST run BEFORE lock, not be deferred as a Pending assumption
  discharged "at implementation" — because if it fails, the fix is not a
  documentation tweak but a second bound term that changes the diagnostic
  contract's shape (what fields the refusal reports), which changes
  Normative Contracts text that downstream consumers (RDR 0005's envelope,
  RDR 0006's finding contract) already depend on
And the RDR MUST NOT lock with a load-bearing normative clause (the bound
  clause) whose correctness is explicitly conditional on a test that has
  not run
```

**AT-8 (catches C-8, silent-optional default is a correctness footgun).**
```
Given "no optionality marker" silently means "optional" (the conservative
  default)
When an author omits the optionality marker on a tag that is, in every
  real usage, always present
Then the declaration loader or a lint pass MUST warn (even if not block)
  that an undeclared-optionality tag is being treated as optional, so the
  author has a chance to notice before the presence dimension silently
  doubles their scoped product or before a guard row silently gains a
  legal `guard_unevaluable` escape hatch they didn't intend
And RDR 0002's authoring surface MUST make `optional` a required field
  (forcing an explicit choice) rather than an optional field with a
  silent conservative default, since "conservative" here is conservative
  for the exhaustiveness *proof*, not for the *author's* mental model of
  their own tag
```
