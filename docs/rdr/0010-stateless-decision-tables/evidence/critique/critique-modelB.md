Model: claude-sonnet-5

## Findings Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | 0010:§technical-design (guard-vs-match discriminator) | Authors will discriminate decision tables with `[rule.match.<key>]` instead of `[rule.guard.all.<key>]` because match reads exactly like a decision-table condition in every other rule engine, and only a lint finding (not the grammar) stops it | Author writes what looks like a correct, complete decision table; `intrastate lint` passes clean; the table is silently unprovable and rows can be missing without any signal until C5's fence lands correctly | §1, premortem, AT-1 |
| C-2 | 0010:C1, 0010:§authority table | The "one accessor, four readers" invariant is a discipline convention, not a compiler-enforced one — a future PR touching `normalizeRule`, `reach`, `checkDanglingEdge`, or `checkCoverage` can read `len(owned)==0` directly instead of the class accessor, and nothing stops it | Two of the four sites disagree after an unrelated future change: loader accepts a decision table, lint reports it as a rootless machine (or vice versa) — a resurrection of exactly the Alternative-2 misclassification bug this RDR was written to prevent | §2, premortem |
| C-3 | 0010:A6, 0010:C3 ("values MUST be strings") | The RDR's own evidence for A6 quotes a consumer example (`Next: /rdr-prelock 0046 critique`) that is a command *with arguments* jammed into one string — A6 explicitly defers structuring it, but ships C3 as a hard string-only contract with no escape hatch | First real consumer beyond the motivating one needs a second field (arguments, a list, a nested value) and cannot express it; `emit` becomes a string-concatenation micro-language the moment two consumers disagree on shape, and the "widening, not a break" promised in A6 requires a second RDR before it's usable | §3, premortem |
| C-4 | 0010:C2, 0010:A9 | The "otherwise" idiom (an escape row rescuing `no_match` with `[rule.emit]`) is authoring convention documented only in Phase 4 docs, not enforced by the grammar; an author who writes an overlapping catch-all ordinary row gets `flow-ambiguous-match` with no guidance toward the correct idiom | First-time author hits a refusal that names a symptom (`ambiguous_match`) but not the fix (use an escape row); support/Slack churn until the doc is read | §1, AT-4 |
| C-5 | 0010:C4 ("`--outcome` remains required in both classes") | A decision table by definition has no state to select an outcome from either — the single-outcome-per-call ergonomic gap (BR5) is explicitly parked as "a plain kata," not solved here, yet decision tables are exactly the class most likely to have one outcome and most annoyed by having to spell `--outcome` every time | Every single decision-table invocation in the wild types a redundant `--outcome <the-only-one>`; the CLI ergonomics this RDR was supposed to fix (no scratch artifact, no dummy tag) are only half-fixed — the other half is deferred indefinitely with no tracking commitment beyond "parked" | premortem |
| C-6 | 0010:A5, 0010:§decision-rationale | Hit policy is "verified" against exactly one consumer (rdr#tmxk) with an explicit "⚠ no corpus coverage for decision-table hit policies" flag; exact-one is locked as *the* policy, not *a* policy, with no versioning or extension point in C3/C5 for first-hit/priority | Second consumer wants overlapping rows with priority (a completely standard DMN "first-hit" or "priority" policy) and discovers the kernel has no mechanism — this becomes RDR 0001's problem, not 0010's, but 0010 shipped a grammar (`[rule.emit]`, `class`) that implies decision-table generality it does not have | §3, premortem |
| C-7 | 0010:§technical-design Phase 4, A9 consequence note | "A table over an N-outcome alphabet needs N escape rows" is called out as a consequence the Phase 4 docs "must say," but nothing in C2/C5 requires N escape rows to exist — a table with 3 outcomes and 1 escape row silently leaves 2 outcomes' coverage open with no distinguishing lint code from "correctly covered" | An author adds a second outcome to an existing decision table months later, forgets the second escape row, and coverage silently degrades for the new outcome only — caught only by an alert reader of `intrastate lint --as=json`'s coverage-gap output, not by any structural signal | §1, premortem, AT-4 |
| C-8 | 0010:A12 (Pending), 0010:A13 (Pending) | Two of thirteen Critical Assumptions are still `Pending` at what reads like a near-final draft (Premortem PASS, grounding sweep clean, 3amigo run) — both gate the correctness of C5, the RDR's most novel contract | If A12 or A13 fails verification during Stage 4/6, C5 needs rework post-lock in a document the team has committed never to amend — forcing either a hasty reinterpretation or a full new RDR to patch a contract that shipped with known-unverified load-bearing assumptions | §2, premortem |

## 1. The three most likely ways implementation goes wrong

### Failure 1: The guard-vs-match discriminator is a trap the grammar does not prevent

**Root cause in the RDR:** The entire soundness of decision-table coverage rests on authors discriminating with `[rule.guard.all.<key>]` atoms, never `[rule.match.<key>]` atoms — because `0006:C7`/`0003:C13` define match atoms as scoping the *group*, contributing zero assignment to the coverage product. This is stated repeatedly (Approach, C5, Illustrative Code, Risks, MVV step 1) as if repetition were the mitigation. It is not: the mitigation is C5's fence (`graph-unprovable-coverage` on a zero-dimension group), which is still `Pending` (A13) at the time of this review.

**The specific passage that enabled it:** `0010:§technical-design` Illustrative Code block: "A decision table's discriminating dimensions are authored as **guard** atoms, not match atoms. This is 0006/0003's existing split, not a rule this RDR introduces." That sentence is true and completely non-obvious to any newcomer. Every other rule/decision-table system a working engineer has touched (DMN, Drools, a switch statement) treats the "if this field equals this value" clause uniformly — there is no engineer intuition that "match" and "guard" mean structurally different things with respect to a coverage *proof*. The RDR itself concedes in Key Discoveries: "only this RDR's illustrative example had drifted from it" — i.e., even the RDR's own authors got this wrong once while writing the example.

**The symptom the user will see:** An author builds a decision table using `[rule.match.status]` for every discriminating field (the natural, DMN-like way to write it), gets a green `intrastate lint`, ships it, and later discovers uncovered cells resolve to `no_match` in production with no warning at authoring time — *unless* A13 lands correctly and unless the author encounters the `graph-unprovable-coverage` finding and correctly diagnoses "unprovable" (a word implying "we can't tell," not "you did it wrong") as "rewrite your match atoms as guard atoms." The fence, even when working, produces an advisory-sounding blocking finding whose fix is non-obvious from the message alone.

### Failure 2: The four-site class invariant drifts under normal maintenance pressure

**Root cause in the RDR:** C1 declares "one accessor, four readers" (`normalizeRule`, `reach`, `checkDanglingEdge`, `checkCoverage`) as the mechanism that prevents the loader and lint from disagreeing about class. The RDR's own Authority table names this "the drift risk." But the *mitigation* offered — "the agreement check lives in the loader only (C1); lint and normalize read the loaded class, never the owned set" — is a coding convention, not a compiler-enforced invariant, and not a lint rule against the intrastate codebase itself. Nothing stops a future PR (fixing an unrelated bug, or written by someone who hasn't read this RDR) from writing `if len(owned) == 0` directly in a fifth call site instead of `model.Class()`.

**The specific passage that enabled it:** `0010:§authority`: "Cue: the class is one decision read at four sites that 'must agree' (the drift risk below), and C1 fixes exactly one writer for it." The RDR names the risk explicitly and then declares it mitigated by discipline ("reads the loaded class, never the owned set") with no test, linter rule, or code-level enforcement cited anywhere in the Testing Strategy that would catch a fifth site reading `len(owned)` instead of the accessor. Scenario 1 in Testing Strategy tests *behavior* (load refuses/accepts) not *implementation discipline* (that all call sites use the one accessor) — so a future regression that reads owned-set directly and happens to agree with class today would pass every test in this RDR and only diverge later when someone adds a state machine with a temporarily-empty owned set mid-refactor, or vice versa.

**The symptom the user will see:** Months post-ship, an unrelated refactor to `internal/graphlint` or `internal/table` re-derives "is this a decision table" from `len(owned) == 0` in a new call site instead of `model.Class()`. Silent misclassification returns — exactly the failure Alternative 2 was rejected for, reintroduced through the back door by a maintainer who never read this RDR (which, per policy, is a prompt not documentation, and will not be re-read by every future contributor).

### Failure 3: `emit` is string-only and the RDR's own evidence shows it won't stay that way

**Root cause in the RDR:** A6 verifies "sufficiency" by checking the motivating consumer's actual answer shape and finding it's a command *with arguments* (`/rdr-prelock 0046 critique`) packed into one string. The assumption's own "If wrong" clause states plainly: "consumers pack structure into strings... and the wire shape needs an `any`-typed value later (a widening, not a break)." This is not a hedge — it is a documented near-miss. The RDR ships C3 with `sourceRule.Emit` hard-typed `map[string]string`, and BR3 explicitly rejects `any`-typed values for this cycle ("a dump column needs a deterministic rendering").

**The specific passage that enabled it:** `0010:A6`: "sufficiency rests on the consumer treating the command as one opaque string, and a consumer needing the verb and its arguments as separate fields is the widening A6 defers." The RDR is aware and honest about this, which makes it worse, not better, as a design risk: it is a known, named, accepted limitation shipped anyway because only one consumer was checked.

**The symptom the user will see:** The second real consumer of decision tables (not the motivating rdr#tmxk model) wants to emit a structured answer — say, a verb plus a list of arguments, or a numeric priority alongside a label — and cannot, because `[rule.emit]` values must be strings. They either (a) invent their own string-encoding convention inside the string value (exactly the "consumer-invented format" this RDR's Problem Statement condemns for the old scratch-artifact workaround), or (b) file a new RDR to widen `emit` to `any`-typed, which per this RDR's own Testing Strategy Done-clause discipline ("no golden-file tree... changed line differs only by...") means another full grammar/dump/payload change cycle through the same three seams.

## 2. The one section that will be rewritten within 6 weeks of shipping

**C5, specifically the zero-participating-dimension fence (A13's territory).**

This is the one section where the RDR's own confidence outruns its verification. A13 is `Pending` — "to verify" — at what is otherwise a heavily pre-verified document (11 of 13 assumptions Verified, premortem PASS, grounding sweep clean, 3amigo run). The mechanism A13 needs is grounded (`checkCoverage`'s `len(dims) == 0` branch, `guard.Dimensions` collecting only guard atoms) but the actual emission — "whether `emitStructurallyUnprovable`'s existing per-dimension shape admits a group-level emission with no dimension to name (its `Dimension` field would be empty)" — is admittedly unverified. The RDR's own Testing Strategy Trace table flags this as "the one Trace row with no witness."

Why this is the one that gets rewritten: it is the single mechanism standing between "decision tables prove exhaustiveness" (the RDR's entire stated reason to exist — "the property that makes stateless tables worth supporting at all," per the Problem Statement) and "decision tables silently lint clean while incomplete" (the exact bug the scratch-tag PoC was built to expose). If the `Dimension`-field-empty shape doesn't fit cleanly into the existing finding renderer, or if it turns out some checked-in state-machine fixture has a legitimate zero-dimension group that would false-positive under a wrongly-keyed arm (a risk A13's own "If wrong" names), this clause gets patched at the *first* real decision-table author's angry bug report — because it is precisely the "an author discriminates with match atoms" failure mode from Finding C-1 above, which is the most intuitive way to author a decision table and the one the RDR's own writers drifted into once already.

## 3. The one assumption that will not survive first contact with a real user

**A5/A6 combined: "the kernel's exact-one selection is sufficient" and "a flat string-valued emit table is sufficient" — both verified against exactly one consumer, rdr#tmxk.**

A5's own evidence section admits it outright: "⚠ no corpus coverage for decision-table hit policies... this is verified as *sufficient for the motivating consumer*, not as a general survey of hit policies." A6 similarly checks only one consumer's actual answer shape.

Real decision tables in the wild — the artifact DMN, Drools, and every production rules engine model — very commonly want a *priority* or *first-match* hit policy specifically because authors write overlapping catch-all rows on purpose (a "default" row that's broader than the specific rows above it, evaluated in order). This RDR's kernel forces authors into the "otherwise = escape row rescuing no_match" idiom (A9) instead, which works for exactly one failure class (`no_match`) and refuses `ambiguous_match` the instant two ordinary rows legitimately overlap by design. The very first external team that tries to model a table with graceful degradation ("if nothing specific matches, fall back to X, and if two specific things match, prefer the more specific one") will hit `flow-ambiguous-match` on a table they consider correctly authored, and will discover the kernel — RDR 0001's territory, explicitly out of scope here — has no priority mechanism. This is the load-bearing assumption most likely to be falsified by the second consumer rather than the first, because the RDR verified it against a single internal, cooperative, already-known-shape use case (an RDR-status-navigator table) rather than against any external decision-table corpus — a gap the RDR itself flags with a warning symbol and then locks anyway.

## 4. Premortem

It is six weeks after RDR 0010 shipped. `class = "decision-table"` is live, `internal/table`, `internal/graphlint`, and `internal/cli` all carry the changes, `models/rdr.toml` is untouched as promised, and the motivating consumer (rdr#tmxk's navigator model) has been rewritten to use it. Three tickets have landed against the feature.

**Ticket 1 (filed by a new team building a "which reviewer owns this PR" decision table):** They authored their table with `[rule.match.file-extension]`, `[rule.match.team]`, `[rule.match.risk-level]` — three match atoms, because that's how every rule they'd previously seen worked, and because "match" reads as "the thing that decides." `intrastate lint` passed with zero findings the first three times they ran it while iterating. They shipped it. Two weeks later a fourth `risk-level` value was added to the domain, and no new rule was authored for it — the table was now incomplete for that value. Lint still passed clean, because the scoped product over zero guard dimensions closes vacuously on any single row (exactly the mechanism `internal/graphlint/coverage.go::checkCoverage`'s `len(dims) == 0` branch produces). A resolve call for the new `risk-level` value returned `no_match` in a CI job that assumed exhaustive coverage was lint-enforced — which the Problem Statement of 0010 promised it would be. The bug was traced back to A13's fence (`graph-unprovable-coverage` for zero-participating-dimension groups) either not having landed correctly, or having landed but not being obvious to the author as "you used the wrong discriminator." Either way, `checkCoverage` and `emitStructurallyUnprovable` are now getting a second pass to make the finding's `Message` field say, in plain language, "author your discriminating fields as `[rule.guard.all.<key>]`, not `[rule.match.<key>]`" — a UX fix that should have been in the original finding text from day one, not left to a doc reader.

**Ticket 2 (filed against `internal/graphlint/analysis.go::checkDanglingEdge`):** A contributor refactoring the reachability check for an unrelated performance fix (batch-processing multiple models) introduced a helper that checked `len(model.OwnedTags()) == 0` directly, instead of calling `model.Class()`, because it was faster to write and the two conditions "happened to be equivalent" in every fixture the contributor tested against. Three months later, a state machine under active development temporarily had zero owned tags mid-refactor (someone was migrating an accessor and hadn't re-added the write block yet) — the class was still declared `state-machine`, but the helper's owned-set check made it get treated as a decision table by the new code path, silencing `graph-dangling-edge`'s missing-root arm precisely when it should have fired. The "one writer, four readers" invariant C1's Authority table promised had quietly become "one writer, five readers, one of which cheats." This is fixed by adding a lint rule against the intrastate codebase itself (banning direct `len(owned)` class inference outside `model.Class()`) that this RDR never proposed.

**Ticket 3 (filed by rdr#tmxk's own maintainers, the motivating consumer):** Having successfully migrated off the scratch-artifact workaround, they immediately hit the wall A6 named as a known limitation: they wanted the `emit` block to carry not just `next = "/rdr-prelock 0046 critique"` as a single string, but the verb and the RDR id as separate fields so their own CLI tooling could compose the command programmatically instead of shell-splitting a string. They filed a "widen `[rule.emit]` to `any`-typed values" request within the six-week window, which is functionally a request to reopen RDR 0010's C3 — an RDR whose content, per project policy, is never amended, meaning this becomes a new RDR (0010's stated escape hatch: "a widening, not a break") going through the full Stage 1–8 cycle to unlock what A6 predicted on day one it would need.

Across all three tickets, the common thread is not that any individual contract (C1–C5) was wrong as stated — they were rigorously verified against source, and the spikes back them. The common thread is that **the contracts protect the mechanism, not the author.** The RDR proves that the lint *can* catch an incomplete table when authored correctly, and that the four class-reading sites *can* stay in agreement if every future contributor reads and honors the Authority table. Neither guarantee holds against an author who doesn't know the guard/match distinction exists, or a future maintainer who never reads this RDR (which, by the team's own stated policy, is a prompt for implementation, not living documentation, so nothing forces a re-read).

## 5. Acceptance tests that would have caught each failure at RDR-review time

```gherkin
Feature: Decision-table authoring does not silently under-prove coverage

  Scenario: A newcomer discriminates with match atoms instead of guard atoms
    Given a decision-table model with three ordinary rules
    And each rule discriminates on "status" using "[rule.match.status]" (not guard)
    And the table is deliberately incomplete (one status value has no rule)
    When "intrastate lint --model m.toml --as=json" is run
    Then the exit code MUST NOT be 0
    And the finding message MUST name the specific fix ("author guard.all/guard.unless
      atoms, not match atoms, to make this field a coverage dimension")
    # Catches Finding C-1 / Premortem Ticket 1: today's RDR text only requires the
    # finding CODE (graph-unprovable-coverage) to fire; it does not require the
    # MESSAGE to be actionable without reading the RDR's authoring docs.

  Scenario: A fifth call site infers class from the owned set instead of the accessor
    Given the intrastate codebase itself
    When a static check scans internal/table, internal/graphlint, and internal/cli
      for any comparison of the form "len(owned) == 0" or equivalent outside
      model.Class()'s own implementation
    Then zero such comparisons MUST exist
    # Catches Finding C-2 / Premortem Ticket 2: the RDR's Authority table declares
    # "one writer, four readers" as the invariant but ships no enforcement — this
    # AT would have forced either a linter rule or an explicit statement in the
    # RDR that no such enforcement exists and drift is an accepted risk.

  Scenario: emit must carry structured values for a second consumer
    Given a decision-table model whose answer is a verb plus an argument list
    When the author writes "[rule.emit]" with a nested table or a list value
    Then the load MUST either succeed with the structure preserved
      OR the RDR MUST document, before lock, a concrete widening path with an
      owner and no open-ended "later RDR" deferral
    # Catches Finding C-3 / Premortem Ticket 3: A6 already knows the string-only
    # shape is insufficient for a command-with-arguments case (its own cited
    # example). This AT would have forced the choice to be made now, not deferred
    # on a documented near-miss.

  Scenario: An "otherwise" row must exist per outcome, and its absence is caught
    Given a decision-table model with two declared outcomes
    And an escape "otherwise" row rescuing no_match for only one outcome
    When "intrastate lint --model m.toml --as=json" is run over an incomplete
      second-outcome coverage
    Then the finding MUST distinguish "this outcome has no otherwise row" from
      "this outcome's table is complete" -- not merely report a coverage gap
      indistinguishable from an ordinary missing-cell finding
    # Catches Finding C-4/C-7: the RDR states the N-escape-rows-per-outcome
    # consequence must be documented (Phase 4) but ships no lint signal
    # specific to "you forgot the otherwise row for outcome N", only the
    # generic coverage-gap code that also fires for ordinary missing cells.

  Scenario: A single-outcome decision table does not require --outcome
    Given a decision-table model declaring exactly one outcome
    When "intrastate flow resolve --model m.toml --tag k=v" is run with no
      "--outcome" flag
    Then the call MUST succeed using the sole declared outcome
    # Catches Finding C-5: BR5 parks this as "a plain kata," but for the class
    # this RDR exists to serve (decision-only models), a mandatory --outcome
    # flag on every single-outcome table is friction this RDR's own Problem
    # Statement complains about in spirit (unwanted authoring/invocation
    # surface) and never actually eliminates.

  Scenario: A second consumer needs an overlapping-priority hit policy
    Given a decision-table model with two ordinary rules whose guard atoms
      legitimately overlap by design (a specific rule and a broader fallback)
    When "intrastate flow resolve" is invoked with tags matching both rows
    Then the RDR MUST state, before lock, whether this is out of scope for
      the kernel (RDR 0001) with a named follow-up, or resolved via a
      documented idiom equivalent in ergonomics to hit-policy priority
    # Catches Finding C-6: A5 already flags "no corpus coverage for
    # decision-table hit policies" with a warning; this AT would have forced
    # an explicit scope statement instead of silent reliance on one
    # consumer's requirements matching the kernel's one available policy.

  Scenario: A13's fence is verified before lock, not left Pending
    Given the zero-participating-dimension group fence C5 requires
    When the MVV's Trace table is reviewed at the Finalization Gate
    Then no Trace row bearing on a locked Normative Contract may be
      "unwitnessed" — A12 and A13 MUST resolve to Verified or the affected
      contract (C5) MUST be narrowed to exclude the unverified clause
    # Catches Finding C-8: the RDR's own Trace table already flags step 2″ as
    # unwitnessed. This AT is simply enforcing the RDR's own Finalization
    # Gate "Assumption Verification" clause, which requires every Pending
    # assumption to carry "a plan to verify before implementation begins" —
    # A12/A13 have plans, but the gate should refuse lock while two
    # load-bearing assumptions behind the RDR's most novel contract remain
    # Pending this close to lock.
```
