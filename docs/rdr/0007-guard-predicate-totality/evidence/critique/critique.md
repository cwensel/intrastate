Model: claude-opus-5[1m]

# Critique — RDR 0007 Guard predicate totality

## Findings Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | A5 Evidence (c) "the two rows are provably disjoint: no assignment has X both absent and present"; A5 **Carried constraint** "RDR 0003 does not say whether *absence* is itself an element of the finite domain product" | The blessed two-row absence pattern is unlintable and probably lint-rejected. RDR 0003's overlap/coverage check runs over the **declared finite domain product** (`a row's accepted assignments are the intersection of all positive `all` atom domains minus … the row's full `unless` block`). Absence is not an element of any declared domain in RDR 0002/0003 — every fixture tag declares `domain = [...]` with no absent/null member. So the absence-row and the value-row project onto the *same* assignment set, and RDR 0002's `ambiguous overlap` validation category (a stable data-level MUST) fires at table load. A5 admits the domain question is open and then declares the assumption **Verified** anyway. | Author writes the RDR's own blessed pattern, table fails to load with an `ambiguous overlap` lint error naming two rows they were told are disjoint. The RDR's only sanctioned escape from the domain rule is unusable. | §1, §3, premortem, AT-1 |
| C-2 | A10 **Status: Pending**; Prerequisites "A10 verified … Must close before lock"; Normative Contracts *SCOPE OF THE ATOM VOCABULARY* "the implementation of `Evaluate` MUST recover the atom structure of the guard it is given" | Every atom-level normative clause binds a seam that cannot see atoms. `Row.Guard` is `string`; `Evaluate(guard string, view TagSet)` is the whole seam; no RDR states the mapping. The RDR's stated resolution is to *hand the obligation to RDR 0003* — i.e. lock a contract whose implementability is delegated to a Final, already-locked peer that must then be reopened. The RDR's own house rule ("we never amend rdrs") makes that handoff a dead letter. | Implementer reaches Phase 2, cannot encode vectors 4/6/7/8 because a vector's input has no representation, and either invents a guard grammar inside `internal/resolve` (silently owning RDR 0003's contract) or ships Phase 1 doc-comments and calls the RDR done. | §1, §2, premortem, AT-2 |
| C-3 | Normative Contracts, aggregation block: "if any surviving candidate row's guard is GuardUnevaluable, the resolution MUST refuse guard_unevaluable — a decided-GuardTrue sibling row MUST NOT be selected"; "This veto is intended, and it is the cost the RDR accepts" | The resolution-level veto is a table-wide denial-of-service that no user consented to, and it is *strictly wider* than the masking path the RDR exists to close. One unreadable observed tag in one unrelated-but-matching row halts every transition for that outcome, including the row that decides cleanly and whose writes touch nothing the absent tag governs. It is also non-escapable by construction, so there is no table-level remedy. | Flow that worked yesterday returns `guard_unevaluable` today for every artifact in a given state, naming a rule the user is not trying to take. No retry, no escape row, no workaround short of editing the table or the artifact. | §1, §3, premortem, AT-3 |
| C-4 | Trade-offs / Consequences: "Positive: missing artifact state becomes unmaskable behind an escape … The masking path recorded under *Background* stays open until that evaluator is built"; Phase 1 "doc comments in `internal/resolve/resolve.go` only; no behavior change" | The RDR ships zero behavior. Phase 1 edits comments; Phase 2's vectors run against a stub the RDR itself must build; Phase 3 hands obligations to an unimplemented RDR. The masking probe recorded in *Background* — the concrete defect kata `xg7p` raised — is still reproducible on `main` after this RDR closes. The tracker closes, the bug does not. | User files the same finding again against the shipped kernel and is told it is "already decided by RDR 0007." | §1, §2, premortem, AT-4 |
| C-5 | Risks and Mitigations, evaluator-drift risk: "`go test ./internal/resolve/` stays green against the stub no matter what 0003 builds"; Phase 2 item 2 "An exported conformance harness"; Phase 3 "RDR 0003's implement stage must instantiate Phase 2's exported harness" | The single mitigation for the RDR's headline risk is an obligation on a *different, already-Final* RDR that has not accepted it. The RDR states the gap in its own prose ("until RDR 0003 accepts it, drift is caught by review only") and then keeps the risk marked mitigated. An exported harness in a non-`_test` package also puts `testing`-shaped test code into production build surface under `internal/`, which nothing in the tree currently does. | RDR 0003's evaluator ships folding absence into false for `contains`; CI is green; the masking path the RDR exists to close is reintroduced with no test failure anywhere. | §1, §2, premortem, AT-5 |
| C-6 | A11 Evidence "`TagSet.Lookup` … is the only *exported* accessor"; Normative Contracts *PRESENCE IS PROVENANCE-BLIND* | Provenance-blind presence hands the caller control over guard decidability. `Input.Observed` is caller-supplied and unauthenticated; `assemble` writes observed keys into the same map. A caller who supplies `X=<anything>` as an observed tag converts a `guard_unevaluable` into a decided verdict on a value the artifact never had. ADV-4/ADV-5 exist precisely because the kernel elsewhere refuses to let observed tags stand in for owned state; this clause reverses that for guards without addressing the conflict. | Operator "fixes" a `guard_unevaluable` refusal by passing an observed tag on the command line and gets a plan whose writes are computed off caller input rather than artifact state. Nothing in the refusal or plan flags it. | §1, §3, premortem, AT-6 |
| C-7 | Failure Modes, *Diagnostic gap*: "`guard_unevaluable` carries `Refusal.Guard` (one row's guard text) and `Refusal.Rows` … it does not name the absent tag directly"; Problem Statement "told plainly that the artifact state needed to decide was missing" | The RDR does not deliver its own Problem Statement. The user outcome is "told plainly the state was missing"; the shipped refusal names a rule id and an opaque guard string. Since `Row.Guard` may not even carry atom structure (C-2), the guard text may be a bare key like `"is-fast-lane"` — naming nothing. The RDR classifies this as "diagnostic enrichment" and defers it to an unimplemented peer. | Refusal reads `guard_unevaluable: rdr.draft.successful.fast (is-fast-lane)`. User has no idea which tag was missing, tries the obvious ones, gets the same refusal. Non-escapable, so they cannot route around it either. | §1, §3, premortem, AT-7 |
| C-8 | Overrides: "This RDR *fixes* the meaning of `Row.RequiresOwned` … a field no RDR ever defined"; Normative Contracts "Row.RequiresOwned names the owned tag keys the row's post-guard transition depends on — the keys its `Writes` require" | This is not a narrowing, it is a new normative contract on RDR 0001's frozen `Row` shape, asserted from outside RDR 0001, with no corresponding validation. Nothing checks that `RequiresOwned` matches `Writes`; nothing rejects a `RequiresOwned` key the row does not write; the frozen fixtures (`escapeRow`, `singleMatchTable`) all declare `RequiresOwned: ["status"]` for rows whose match reads `status` — i.e. the *old* meaning. The new contract is violated by the RDR's own cited fixtures on day one. | Nothing visible. The doc says one thing, every table in the tree means another, and the next author reads the fixtures. | §1, §2, premortem, AT-8 |
| C-9 | Normative Contracts, empty-block clause "an empty `unless` block MUST be treated as absent"; A9 **Status: Pending**; Prerequisites "A9 verified … Must close before lock" | A normative MUST is written conditional on an unverified assumption, with an inline self-cancelling clause ("if RDR 0002's normalization already fixes this identity, the clause is 0002's and is struck from here"). A normative block that may be struck is not normative. Testing Strategy scenario 9(b) then admits it is "*not* assertable at the kernel seam" — so the clause is unverifiable where it is written and unowned where it matters. | If the wrong identity is chosen, every row carrying an omitted `unless` block is disabled: whole-table silent failure, every transition returns `no_match`. | §1, §2, premortem, AT-9 |
| C-10 | Investigation: "SCXML §5.9.1 … MUST place the error 'error.execution' in the internal event queue … Intrastate's kernel has no second channel — the verdict is the only observable — so an honest transplant surfaces the error in the verdict itself" | The prior-art argument is inverted. SCXML's rule is *continue execution, signal the error on a side channel*. The RDR's rule is *halt resolution, signal nothing about which tag*. The transplant preserves the error signal and discards SCXML's actual disposition (the machine keeps running). Calling that "honest" while dropping the payload SCXML mandates (the error naming what failed) uses the citation to authorize the opposite of what it says. | Indirect: the RDR's strongest external justification does not support the veto (C-3) or the missing diagnostic payload (C-7); both were argued through on this citation. | §1, premortem, AT-10 |
| C-11 | Decision Rationale: "Operator recovery is the criterion the matrix does not score, and B is genuinely the worst of the three on it"; scored matrix rows Correctness/Reversibility/Blast radius | The scored matrix omits the criterion the chosen approach loses on, then re-admits it in prose after the decision is stated. Every scored row is written in terms B wins (masking closure, kernel-change avoidance); operability — the thing the user experiences — is scored nowhere. A matrix whose criteria are selected after the winner is known is decoration. | Indirect: the operability cost (C-3, C-7) was never weighed against the correctness benefit, so no one asked whether a narrower veto scope would have bought most of the correctness at a fraction of the cost. | §1, §2, premortem |
| C-12 | Testing Strategy preamble: "Rows 1–3 are the MVV and run against `internal/resolve/fixtures_test.go`'s `fixtureGuards`"; Existing Infrastructure Audit row "Verdict-only stub evaluator for the MVV … Reuse (rows 1–3 only)" | The MVV validates nothing this RDR decides. Rows 1–3 hand the kernel a verdict and assert the kernel's already-frozen mapping — behavior A3 says is unchanged and already tested. The domain rule (rows 4–8) is the RDR's entire content and is excluded from the MVV, blocked on A10 (C-2). An MVV that passes on day one against shipped behavior is not validation. | Implementation is declared "MVV green" having exercised zero clauses of the Normative Contracts. | §1, §2, premortem, AT-11 |
| C-13 | Normative Contracts: "Among surviving rows, absent owned state is reported BEFORE an undecidable guard … pinned here because no test currently contends the two within one survivor set"; Testing Strategy *Owned-state ordering gap* "flagged for Phase 2" | The RDR pins a precedence it acknowledges is untested, in the same document that defers the test to a phase blocked on A10. The precedence is also load-bearing for C-8's narrowing: under the new `RequiresOwned` meaning (post-guard write deps), reporting `owned_state_unavailable` ahead of `guard_unevaluable` reports a *write* problem in preference to a *decidability* problem — the less relevant diagnosis, inverting the ordering's own stated rationale ("absent owned state … is frequently the reason the seam could not decide"). That rationale is only true under the *old* meaning the RDR just removed. | Row cannot decide its guard *and* is missing a write dependency. User is told the write dep is missing, supplies it, and gets `guard_unevaluable` on the next run — two round trips for one defect. | §1, §2, premortem, AT-12 |
| C-14 | Metadata **Profile**: "foundational — one contract"; Proportionality gate "confirm this RDR is the sole author of at most one independent load-bearing contract" | The RDR authors at least five independent contracts: the atom domain rule, the strong-Kleene combination rule, provenance-blind presence, the resolution-level aggregation veto, and the `RequiresOwned` redefinition. Each has distinct consumers (0003 evaluator, 0003 evaluator, 0003 evaluator + kernel, kernel, kernel + 0002 authors) and distinct failure modes. The Finalization Gate's own split signal is triggered and unaddressed — the gate section is an unfilled template. | Indirect: the seams lock together, so a later fix to the veto scope (C-3) or the `RequiresOwned` meaning (C-8) cannot be made without reopening the domain rule, which is the one part that is right. | §2, premortem |
| C-15 | Consequences: "The exposure today is bounded — no `GuardEvaluator` implementation exists yet, so no shipped table depends on closed-world reads and there is no installed base to migrate" | "No installed base" is measured against implementations, not against authored intent. RDR 0002's shipped fixtures already author guards (`rdr-fixture.toml`, `guard-fixture.toml`) and RDR 0003's fixture declares `[tags.cluster_eligible] provenance = "observed"` with `exists = true` — an observed tag whose presence is a runtime fact. Every value-comparing guard in those fixtures over an observed or bounded-int tag becomes a `guard_unevaluable` risk the moment an evaluator exists. The migration is real and pre-dated; it is only invisible because nothing runs yet. | First evaluator build turns the project's own reference fixtures into refusal generators, and there is no inventory tooling (the RDR says so) to find them. | §3, premortem, AT-13 |

Nothing in the prose below is unindexed: every defect raised in §1–§5 resolves to a row above.

---

## 1. The three most likely ways implementation goes wrong

### 1.1 The implementer cannot encode the vectors, so the RDR ships as comments

**Root cause.** The RDR writes eight atom-level normative clauses that quantify over *operator, referenced tag key, literal, and `all`/`unless` placement*, and binds them to a seam that receives an opaque `string`. It knows this. A10 is `Pending`. The Prerequisites say A10 "must close before lock." The Normative Contracts open with a clause that concedes the structural problem and then delegates it: "How guard structure survives the trip from the authored table … into `Row.Guard` … is RDR 0003's to state, and it MUST be stated before that evaluator is built."

**The enabling passage.** *SCOPE OF THE ATOM VOCABULARY* (Normative Contracts, first block), read with Prerequisites bullet 2: "If the mapping is unspecified anywhere, the resolution is to hand it to RDR 0003 as a named obligation, not to specify a grammar here."

That sentence is the failure. It converts an unresolved blocker into a resolved one by relabelling it. RDR 0003 is `Final`. This repo's stated house rule is that RDRs are never amended. So the "named obligation" has no destination: it cannot be written into 0003, and 0007 refuses to write it. Phase 3 lists it as a handoff item alongside authoring guidance, as if it were documentation rather than the precondition for the document's entire normative content.

**What actually happens.** The implementer runs Phase 1 (doc comments in `resolve.go` — fifteen minutes), reaches Phase 2, discovers that "operators × present/absent × strong-Kleene verdict pairs × `unless` atoms × row aggregation" has no Go representation anywhere in the tree, and faces three bad options: (a) invent a guard AST inside `internal/resolve` — silently taking ownership of RDR 0003's grammar contract, which is the exact "implementer decides a cross-RDR contract silently" failure the RDR rejects Alternative *Evaluator-defined* for; (b) encode vectors against a hand-rolled evaluator keyed on guard *text* — i.e. rebuild `fixtureGuards`, which the RDR explicitly says "cannot express presence/absence at all"; (c) ship Phase 1 and mark Phases 2–3 deferred. (c) is the path of least resistance and is what will happen, because Phase 1 alone satisfies every checkbox a landing skill can mechanically verify.

**Symptom.** The user sees nothing change. `internal/resolve/resolve.go` gets two revised doc comments. The masking probe recorded in *Background* — FALSE guard + absent `RequiresOwned` key + modeled `no_match` escape → `Escaped:true` — still reproduces on `main`. Kata `xg7p` closes. Six months later someone re-derives the same finding. (C-2, C-4, C-12)

### 1.2 The two-row absence pattern — the RDR's only sanctioned outlet — is lint-rejected

**Root cause.** The RDR makes value operators partial and then owes authors a way to say "this row applies when X is absent." A5 supplies it: the two-row pattern, one row guarded `unless … exists` on X, one guarded `X eq v`. A5 is marked **Verified**, Method **Design Decision**.

It is not verified. It is contradicted by the peers it cites.

RDR 0003's overlap check is set-theoretic over the **declared finite domain product**: "a row's accepted assignments are the intersection of all positive `all` atom domains minus the single conjunctive assignment set matched by the row's full `unless` block," scoped to a row group and evaluated "as one product rather than as independent one-dimensional checks." Every tag declaration in the tree — RDR 0002's schema (`tag name, provenance, value kind`), RDR 0003's fixture (`domain = ["Draft", "Final", "Implemented"]`, `min = 0 / max = 3`) — enumerates *values*. None declares an absent element. There is no null in the domain product.

So both rows' accepted assignments project onto the same set: row 2 accepts `{X = v}`; row 1, whose only guard is an existence atom over a dimension the product has no absent member for, accepts either everything or nothing depending on how the prover handles an atom outside the product. Neither reading yields disjointness. Under the "accepts everything" reading, RDR 0002's `ambiguous overlap` validation category — a normative MUST in the stable data-level list — fires at table load.

**The enabling passage.** A5 Evidence (c): "the two rows are provably disjoint: no assignment has X both absent and present." That is a claim about *runtime views*, not about the *declared domain product* the lint reasons over. The RDR then states the real objection itself, in A5's **Carried constraint**: "RDR 0003 does not say whether *absence* is itself an element of the finite domain product used for coverage proofs. If a lint must ever *prove* an absence/value row pair exhaustive (rather than merely evaluate it), that domain question belongs to RDR 0003 and is handed off with Phase 3."

Two things are wrong with that paragraph. First, it scopes the problem to *exhaustiveness proofs* ("prove … rather than merely evaluate") when the binding failure is **overlap**, which RDR 0003 checks in the same product and RDR 0002 lists as a mandatory rejection category — not an optional claim that can be downgraded. Exhaustiveness may be downgraded; overlap is a hard reject. Second, having identified an unresolved domain question that the pattern depends on, it stamps the assumption **Verified** and hands the question to Phase 3.

**Symptom.** An author hits the domain rule (their `status eq "Draft"` guard now returns unevaluable because `status` was not in the view), reads the authoring guidance, writes exactly the blessed two-row pattern, and the table fails to load with an overlap lint error naming both rows. The RDR told them this is provably disjoint. They now believe the lint is broken, or the RDR is, and they have no third option — the RDR explicitly forbids the implicit-existence reading and flags sentinel-stamping as the anti-pattern. (C-1, C-15)

### 1.3 The aggregation veto ships and turns one unreadable tag into a table-wide outage

**Root cause.** The RDR promotes an implementation artifact of `gate` — `if len(undecidable) > 0 { return nil, &Refusal{...} }`, which overwrites the accumulated `selected` slice — into a normative resolution-level rule, and then defends it as a deliberate cost.

The defense conflates two different things. The masking path the RDR exists to close is: *a row's own guard cannot be decided, and that undecidability is laundered into an escapable `no_match`*. The veto goes much further: *any* matching row's undecidability blocks *every other* row, including a row whose guard is decided TRUE from present state and whose writes have nothing to do with the absent tag.

**The enabling passage.** The aggregation normative block: "This veto is intended, and it is the cost the RDR accepts: one unreadable row refuses a table whose other rows decide cleanly. It does not contradict D5's anti-poisoning rationale, which scopes evaluation to *matching* candidates so an UNRELATED row's obligation cannot poison a legal resolution. An unevaluable row that matched the outcome and the tag pattern is not unrelated — it is a live edge whose applicability is unknown."

Read that against D5's actual text in `deviations.md`: "Checking every row would make an unrelated row's `RequiresOwned` key poison an otherwise legal resolution." D5 and D8 exist *because* the project already decided, twice, with two independent verifiers, that a row which does not apply must not decide another row's fate. The RDR's rebuttal — a matching-but-undecidable row is "not unrelated" — is exactly the argument the old owned-before-guard ordering made and lost. ADV-1 is the frozen test that killed it. The RDR cites ADV-1 approvingly four times and then re-derives the losing position at a different scope.

Worse: the veto is non-escapable. RDR 0002 closes escape lists to `no_match` and `ambiguous_match`, which the RDR celebrates as closing the masking path "by construction" — but it is the same construction that removes every table-level remedy from the veto. `owned_state_unavailable` has the same property and D8 called that out as the reason to prune: "the old order therefore converted an escapable condition into an inescapable one."

**Symptom.** A flow that worked yesterday stops working for every artifact in a given state. The refusal names `rdr.draft.successful.slow` — a rule the user is not trying to take and has never seen fire — because that row's guard reads an observed tag the accessor stopped supplying. The row the user wants (`rdr.draft.successful.fast`) decides TRUE and is discarded. There is no escape row that can rescue it, no retry that helps (A6b: the refusal cannot distinguish transient read failure from real absence), and the only fixes are editing the table or the artifact. (C-3, C-6, C-10)

---

## 2. The one section that will be rewritten within 6 weeks of shipping

**`#### Normative Contracts`, the aggregation block** (Normative Contracts, "Aggregation is resolution-level, not merely per-guard…" through "…which is absence-as-false at resolution scope").

It goes first because it is the only block that is simultaneously (a) already implemented, (b) load-bearing on every real transition, and (c) defended on grounds the codebase has already rejected.

The rewrite arrives by this route. Phases 2 and 3 stall on A10 (§1.1), so the domain rule never becomes executable and the atom-level blocks stay inert prose. Meanwhile the veto is live *today* — `gate` already discards `selected` when any survivor is undecidable — and the first real table with two matching rows behind different guards hits it. Because the veto is non-escapable, the first report is an outage report, not a feature request. Someone opens `resolve.go`, finds the four-line branch, finds the RDR block that froze it, and discovers the block's own justification cites D5 to argue against D5's conclusion.

The rewrite will narrow the veto: an unevaluable row blocks selection only when no sibling decides TRUE, or only when the unevaluable row's referenced tags intersect the selected row's writes, or (most likely) an unevaluable row degrades to `ambiguous_match` when a decided sibling exists, restoring escapability. Any of those reopens the block wholesale, and — because RDR 0014's author will discover that C-14's five contracts are welded together — it will drag the `RequiresOwned` narrowing and the ordering clause with it.

The second-place candidate is **Testing Strategy**, which will be rewritten because rows 4–8 cannot be encoded and rows 1–3 pass against shipped behavior — an MVV that is green before implementation begins is a section that gets rewritten as soon as anyone reads it critically. (C-3, C-12, C-13, C-14)

---

## 3. The one assumption that will not survive first contact with a real user

**A5** — that authors can express absence-conditional rows with the two-row existence pattern.

Not A10, which is a *drafting* defect the RDR itself flags as blocking and which a careful implementer will at least trip over. Not A6b, which is honestly carried open. A5 is the one marked **Verified**, argued at length across three lettered sub-clauses, blessed as the "authoring handoff" pattern, and load-bearing on the entire user-facing story — and it dies the first time a real author uses it.

Here is why the user contact is what breaks it, not review. A5 is verified against *evaluation semantics*: at runtime, a view either has X or does not, so the two rows never both fire. That is true and it is the wrong question. The author does not experience evaluation semantics; the author experiences **table load**. Between authoring and evaluation sits RDR 0002's normalization and lint and RDR 0003's overlap proof — both of which reason over the declared finite domain product, in which absence is not an element, and both of which the RDR consults only for the grammar-legality of `unless … exists` (A5(a)) and for the absence of intra-guard disjunction (A5(b)). Sub-clause (c), the disjointness claim, cites row-group scoping and `unless`-priority prose — neither of which speaks to whether the prover can represent an absent assignment at all.

The RDR knows the gap. It writes it down in A5's **Carried constraint** in almost the right words. Then it scopes it to exhaustiveness (downgradeable) instead of overlap (a hard MUST-reject in RDR 0002's validation category list), stamps **Verified**, and defers to Phase 3 — a phase that, per §1.1, will not run.

The first author who needs "route this way when the artifact has no `profile` tag yet" writes both rows, runs the lint, and gets an ambiguous-overlap rejection on the exact pattern the docs blessed. They will then do one of two things, both bad: stamp a `profile = "unset"` sentinel upstream — the anti-pattern the RDR's own Risks section names and whose only mitigation *is* A5 — or file a bug against the lint. Either way the domain rule's user-facing story collapses, because the rule's whole claim to being livable was that absence remained expressible.

Runners-up, for the record: A11's provenance-blind presence (C-6) fails on first contact with an operator who discovers observed tags are caller-supplied and uses them to unstick a refusal; and A3's "no kernel behavior change" (C-8) fails on first contact with the existing fixtures, every one of which uses `RequiresOwned` in the sense the RDR just deleted. (C-1, C-6, C-8, C-15)

---

## 4. Premortem

*Written at draft time, in the past tense, as if the failure had already occurred.*

RDR 0007 locked on 2026-08-14 with A9 and A10 still `Pending` and the Prerequisites checkboxes unticked; the finalize pass read them as "handed to Phase 3" and let the lock through.

**Phase 1 landed in an afternoon.** The implementer edited two doc comments in `internal/resolve/resolve.go` — the `Row.RequiresOwned` block ("names owned tag keys the row's evaluation needs" → "names the owned tag keys the row's post-guard transition depends on") and the `GuardEvaluator` block, which gained a paragraph restating the domain rule. No behavior changed. `go test ./internal/resolve/` was green before and after, which is exactly what A3 predicted and exactly why nobody noticed anything was missing. Kata `xg7p` was closed against that commit.

**Phase 2 stalled inside a day.** The implementer opened `fixtures_test.go`, confirmed that `fixtureGuards.Evaluate(guard string, _ resolve.TagSet)` discards the view, and set out to write the view-reading evaluator the RDR requires. The first vector — scenario 7, `contains` over an absent set-valued tag — needed an input. There was nothing to write. `Row.Guard` is a `string`; there is no atom type; `evaluateGuard` special-cases only `""`. The implementer went looking for the mapping A10 says RDR 0003 owes, found RDR 0003 `Final` at 44,914 bytes with no mention of `Row.Guard` at all, and correctly declined to invent a guard grammar inside `internal/resolve`. Phase 2 was marked *blocked on RDR 0003*; Phase 3 followed it. The exported conformance harness — the only mitigation the RDR offers for evaluator drift — was never built. The status file said `Implemented`.

**Then the veto shipped, because it was already shipped.** The aggregation block froze `gate`'s existing `len(undecidable) > 0 → return nil, &Refusal{KindGuardUnevaluable}` branch as normative. Nobody wrote new code for it; the RDR simply declared the four lines that overwrite the named return `selected` to be the contract.

**The first outage came through the RDR flow itself.** A table modeled two rows on `status:Draft` + `recognized:successful`: `rdr.draft.successful.fast`, guarded on `profile in ["mid","large"]`, and `rdr.draft.successful.slow`, guarded on `cluster_eligible exists` — the observed tag RDR 0003's own fixture declares. The cluster accessor started timing out. Under RDR 0004's contract a read accessor returns "typed tag values or a typed refusal," with no completeness requirement on the values branch — exactly the hole A6b was carried open on — so the accessor returned a partial snapshot with `cluster_eligible` simply absent. `assemble` faithfully built a view without the key, per A6a.

`Resolve` then did the following, in order. `matches` admitted both rows as candidates (neither row's `Match` block reads `cluster_eligible`; the guard does). `gate` called `evaluateGuard` on each: `.fast` decided TRUE from present tags; `.slow` returned `GuardUnevaluable`, per the domain rule, correctly. `missingOwned(survivors, view)` found nothing missing — `cluster_eligible` is observed, and per the narrowed contract it is not a write dependency, so no row declared it in `RequiresOwned`; C-8's redefinition guaranteed the more precise diagnosis was unavailable. The loop over `survivors` appended `.fast` to `selected` and `.slow` to `undecidable`. Then `len(undecidable) > 0` fired, `selected` was overwritten with `nil`, and the kernel returned `guard_unevaluable` naming `rdr.draft.successful.slow`.

Every `rdr` flow transition in the org stopped. The users were RDR authors running `intrastate` from the CLI to advance drafts through prelock; the journey was the ordinary one — recognize `successful`, get a plan, advance `Draft → Final`. What they saw was:

```
guard_unevaluable: rdr.draft.successful.slow (flows/rdr.toml:70)
  guard: cluster-eligible
```

The refusal named a rule none of them had authored intent toward, and a guard string that named no tag — because `Row.Guard` was still opaque (C-2), so `Refusal.Guard` carried whatever text normalization happened to stuff in it. The absent tag, `cluster_eligible`, appeared nowhere in the output. `MissingOwned` was empty, because the missing tag was observed, not owned. The Problem Statement's promise — "told plainly that the artifact state needed to decide was missing" — was not delivered by the refusal that the RDR built to deliver it (C-7).

**The four attempted recoveries all failed, in the order people tried them.**

*Retry.* Natural, since a timing-out accessor sounds transient. The Failure Modes section had anticipated this and said not to — "operators should treat the refusal as 'inspect the artifact and accessor,' not 'retry'" — but that text lived in the RDR, not in the CLI output. A6b's open seam meant the refusal genuinely could not tell them which it was.

*Model an escape.* The obvious table-level fix. RDR 0002 L292 closes escape lists to `no_match` and `ambiguous_match`, so `escape = ["guard_unevaluable"]` was rejected at load as a malformed escape declaration. This was by design and was the RDR's proudest structural claim; it was also the reason the outage had no table-level remedy.

*Supply the tag on the command line.* This one worked, and that was the worse outcome. Per A11 and the *PRESENCE IS PROVENANCE-BLIND* block, an observed tag satisfies presence exactly as an owned one does. An operator passed `--observed cluster_eligible=true`, `assemble` wrote it into `view.tags` with `ProvenanceObserved`, the guard decided TRUE, and both rows became viable — producing `ambiguous_match` for some artifacts and, where only `.slow` matched, a plan whose `Writes` were computed off a value the operator had typed rather than one the cluster reported. `ADV-4` and `ADV-5` exist precisely to stop caller input standing in for artifact state; neither test covers the guard path, because before RDR 0007 nothing said guards read the view provenance-blind (C-6).

*Split the row on absence.* The last resort, and the one the RDR had blessed. An author took A5's two-row pattern: row 1 guarded `unless.cluster_eligible.exists = true`, row 2 guarded `all.cluster_eligible.eq = true`. RDR 0002's normalizer expanded both into the same row group; RDR 0003's overlap proof intersected their accepted assignments over the declared product — in which `cluster_eligible` declares `domain = [true, false]` and no absent member — and reported a non-empty intersection. The table was rejected at load with `ambiguous overlap`, one of RDR 0002's mandatory validation categories (C-1). The author had followed the documented pattern exactly and been told, by the same project, both that the rows are "provably disjoint" and that they overlap.

**The postscript.** The fix that shipped was a one-line change to `gate`: when `selected` is non-empty, do not discard it. That reopened the masking path the RDR existed to close — under the emergency patch, an unevaluable row beside a decided sibling silently yields the sibling's plan — and it was merged anyway, because the alternative was an outage. The domain rule was still inert prose; the conformance harness was still unbuilt; when RDR 0003's evaluator finally landed, it folded absent-key `contains` to false-by-empty-set, and `go test ./internal/resolve/` stayed green, exactly as the Risks section had predicted it would (C-5).

The RDR's ledger recorded it as `Implemented`.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

These are review-time gates, executable against the RDR text and the peer RDRs before lock — not implementation tests. Each names the finding it catches.

**AT-1 — the blessed absence pattern survives the lint the project already mandates.** (C-1)
```gherkin
Given RDR 0007 A5 blesses a two-row absence pattern
  And RDR 0003 computes overlap as the intersection of rows' accepted
      assignments over a declared finite domain product
  And RDR 0002 lists "ambiguous overlap" as a mandatory data-level
      validation rejection category
When a reviewer projects both blessed rows onto the declared domain
     product for a tag whose declaration enumerates only values
Then absence must be an element of that product,
      or RDR 0003 must state an overlap rule for atoms outside it
 And if neither holds, A5 MUST NOT be marked Verified
```
Executable at review time as a hand proof against `guard-fixture.toml`'s `[tags.cluster_eligible] domain = [true, false]`. The RDR reaches the right question in A5's *Carried constraint* and scopes it to exhaustiveness; this gate forces the overlap question, which is not downgradeable.

**AT-2 — no normative clause may quantify over structure its bound seam cannot see.** (C-2)
```gherkin
Given a normative block quantifying over operator, tag key, literal,
      or all/unless placement
 And the seam it binds is GuardEvaluator.Evaluate(guard string, view TagSet)
When the reviewer asks "what Go value is a conforming vector's input?"
Then a concrete answer must exist in this RDR or in a Final peer
 And "handed to RDR 0003 as a named obligation" is NOT an answer,
      because RDR 0003 is Final and this project does not amend RDRs
```
Fails on the *SCOPE OF THE ATOM VOCABULARY* block. Forces the resolution before lock: either specify the mapping here, or cut every atom-level clause and ship only the resolution-level ones.

**AT-3 — a veto wider than the defect must be scored against a narrower one.** (C-3, C-11)
```gherkin
Given the aggregation block blocks selection of a decided-TRUE row
      whenever any matching sibling is unevaluable
 And the masking path in Background requires only that an unevaluable
      row's own undecidability not become an escapable no_match
When the reviewer constructs a table with row A decided TRUE and row B
      unevaluable, where B's referenced tags do not intersect A's Writes
Then the RDR must show that selecting A masks missing artifact state
 And must score at least one narrower alternative (e.g. degrade to an
      escapable class when a decided sibling exists)
 And "the cost the RDR accepts" without such a comparison is not a
      rationale
```
Fails immediately. The block asserts the cost is accepted without ever constructing the case where the veto is strictly wider than the defect. Note this is the same case D5 and D8 already litigated at row scope.

**AT-4 — the RDR must state which shipped behavior changes and when.** (C-4)
```gherkin
Given Background records a reproducible masking probe on main
When the reviewer asks "after this RDR is Implemented, does that probe
     still reproduce?"
Then the answer must be No,
      or the Implementation Plan must name the peer RDR whose landing
      makes it No, and this RDR's status must not reach Implemented first
```
Fails: the RDR answers "yes, the masking path stays open," in its own Consequences section, and plans an `Implemented` status regardless.

**AT-5 — the mitigation for the headline risk must be enforceable by this RDR.** (C-5)
```gherkin
Given the headline risk is evaluator drift from the domain rule
 And the mitigation is a conformance harness RDR 0003 must instantiate
When the reviewer asks what fails if RDR 0003 never instantiates it
Then something in this repo must fail
 And "drift is caught by review only" means the risk is UNMITIGATED
      and must be recorded as such
```
The RDR writes the failure of its own mitigation in plain text and leaves the risk presented as mitigated.

**AT-6 — a presence rule must be checked against caller-controlled inputs.** (C-6)
```gherkin
Given PRESENCE IS PROVENANCE-BLIND makes any provenance satisfy presence
 And Input.Observed is caller-supplied
 And ADV-4/ADV-5 freeze that observed tags cannot stand in for owned state
When a caller supplies an observed tag for a key a guard reads
Then the RDR must state whether the resulting decided verdict is legitimate
 And must reconcile that answer with ADV-4/ADV-5's rationale
```
Unaddressed anywhere in the RDR. A11 verifies only that provenance-blindness matches shipped `Lookup` behavior — never that it is *safe*.

**AT-7 — the refusal must deliver the Problem Statement's promise.** (C-7)
```gherkin
Given the Problem Statement promises the user is "told plainly that the
      artifact state needed to decide was missing"
When the RDR's own Failure Modes state the refusal does not name the
     absent tag
Then either the Problem Statement is overstated and must be narrowed,
      or the diagnostic is in scope and must ship in this RDR
 And deferring it to an unimplemented peer satisfies neither
```
The RDR states both halves of this contradiction, one section apart, and resolves it by relabelling the gap "diagnostic enrichment."

**AT-8 — a redefined field must be checked against its existing uses.** (C-8)
```gherkin
Given Row.RequiresOwned is redefined as post-guard write-dependency keys
When the reviewer greps every RequiresOwned usage in the tree
Then each must be consistent with the new meaning,
      or the RDR must name them as a migration with an owner
```
Fails on `fixtures_test.go`: `singleMatchTable`, `missingOwnedTable`, `escapeRow`, `twoRowsOneGuardFalseTable` and `observedSensitiveTable` all declare `RequiresOwned: ["status"]` for rows that *read* `status` in `Match`. Under the new contract those declarations are wrong, and `escapeRow`'s is wrong even for rows that write `status`, since nothing ties the two.

**AT-9 — no normative block may be conditional on a Pending assumption.** (C-9)
```gherkin
Given a ```normative``` block containing MUST language
When any assumption it depends on has Status Pending
Then the RDR MUST NOT lock
 And a block containing its own strike-out condition
      ("if 0002 already fixes this, the clause is struck")
      is not normative and must be moved to an open question
```
Fails on the empty-block clause (A9 Pending) and, transitively, on every atom-level clause (A10 Pending). This gate alone blocks the lock and would have forced §1.1 and §1.2 to the surface before drafting went further.

**AT-10 — a prior-art citation must support the disposition, not only the vocabulary.** (C-10)
```gherkin
Given SCXML §5.9.1 is cited as the decisive external anchor
When the reviewer compares dispositions rather than error signals
Then SCXML continues execution and signals on a side channel
 And this RDR halts resolution and signals nothing about which tag
 And the citation therefore supports the three-valued verdict but NOT
     the resolution-level veto or the omitted diagnostic payload
 And the RDR must say so
```
Catches the citation being used to authorize two conclusions it does not reach.

**AT-11 — the MVV must fail against the pre-RDR tree.** (C-12)
```gherkin
Given the MVV is rows 1-3 against fixtureGuards
When the reviewer runs the MVV against main before any change
Then at least one row must FAIL
 And if every row passes, the MVV validates shipped behavior,
     not this RDR's contract
```
Rows 1–3 hand the kernel a verdict and assert `gate`'s frozen mapping. They pass today. Scenario 3 (*unevaluable-blocks-true-sibling*) is the only row that pins anything untested, and it pins the veto — which is the clause most likely to be reverted (§2).

**AT-12 — a pinned precedence must have a stated rationale that survives the RDR's own redefinitions.** (C-13)
```gherkin
Given owned-state-before-unevaluable is pinned normatively
 And its rationale is "absent owned state is frequently the reason the
     seam could not decide the predicate"
When Row.RequiresOwned is simultaneously narrowed to post-guard write
     dependencies
Then the rationale no longer holds: a write dependency is not a guard input
 And the RDR must either re-derive the precedence or unpin it
```
Both clauses sit in the same Normative Contracts section and contradict each other's premises. The RDR also concedes no test contends the two conditions.

**AT-13 — "no installed base" must be measured against authored artifacts.** (C-15)
```gherkin
Given the RDR claims no installed base because no evaluator exists
When the reviewer greps the tree for authored guard blocks
Then rdr-fixture.toml and guard-fixture.toml already author guards
      over observed and bounded-integer tags
 And each must be classified as safe or as migration work under the
     domain rule
 And an unclassified authored guard means the migration is real
```
`guard-fixture.toml`'s `cluster_eligible` (observed) and `prelock_iterations` (int, `min=0 max=3`) are exactly the guards that become unevaluable when the artifact does not carry them. The RDR's reference fixtures are its own first migration and are not counted.
