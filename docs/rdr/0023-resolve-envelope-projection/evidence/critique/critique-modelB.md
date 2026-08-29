Model: claude-sonnet-5

## Findings Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `0023:A5` | The whole-tree registration oracle's implementability is Pending at Draft, and the "verification" is scheduled as Phase 2 implementation work rather than pre-lock proof. The claimed mechanism (a total walker over a root with both `help` and `completion` force-materialized) has never been written or run against the real `NewRootCmd()`. | The build stalls mid-Phase-2, or ships an oracle that is quietly weaker than C1's text (reuses a name-skip pattern under time pressure) because the "total walker" turns out to interact badly with cobra's lazy command materialization in ways not anticipated by a desk description. | §1 (Root Cause 1), premortem, AT-4 |
| C-2 | `0023:A2`, `0023:D-wire-byte-format`, Phase 1 | The chosen mechanism — converting five struct fields from concrete types (`string`, `map[string]string`, `[]string`) to pointer types — is asserted as a settled implementation detail ("deferred to Resolve (A2)") but is actually a breaking internal-API change to every reader of `resolvePayload`'s fields, and the RDR's own source-authority census shows all fifteen fields are written at one literal, which the pointer conversion must now populate with `&`-taken locals throughout `runFlowResolve`. | Phase 1 takes longer than a "one flag" change suggests; nil-pointer panics in text-mode `flatten` or in any future reader that dereferences a field expecting a value, not a pointer-to-value. | §1 (Root Cause 2), premortem, AT-2 |
| C-3 | `0023:A8`, `0023:MVV step 3` | The strict-width clause in TEXT mode is Pending — unmeasured — at Draft, while C1 already asserts it as a MUST binding both modes. The RDR admits "the claim is very likely true" but books it as an assumption rather than proving it before lock. If measurement at Phase 2 shows the projected text is NOT strictly shorter on some fixture (e.g., a model where all six echo-only lines are short and the plan group is verbose), C1 either needs a late amendment or ships with a self-contradicting MUST. | CI red on the text-width oracle mid-implementation, forcing a scope renegotiation of a "locked" contract after work has started. | §1 (Root Cause 3), §2, premortem, AT-3 |
| C-4 | `0023:§approach`, `0023:Consequences` | The RDR ships a capability with zero scheduled adopters. No consumer code changes, no migration plan, no owner named for realizing the saving. "Adoption is not scheduled by this RDR and is not one of its phases." | Six weeks post-ship: the flag exists, all oracles are green, and the transcript-token cost that motivated the RDR (kata `intrastate#srz2`) is completely unchanged in production, because nobody adopted `--plan-only`. The org re-discovers the same P1 problem and either has to chase down which consumers should adopt it, or files a new kata to do the migration nobody owns. | §2, §3, premortem |
| C-5 | `0023:C2` (always-keep core), `0023:Trade-offs/Risks` | `model` is explicitly acknowledged as the field a consumer needs to identify which model produced a plan, yet it is placed in the ECHO group (projected away) rather than the always-keep PLAN core, with the stated justification that "revision" will someday cover this need — except `revision` is contractually vacant (renders `""` on every payload the CLI can emit today, per C2 itself). | A chained caller building a multi-model audit trail loses the model identity on every `--plan-only` call, with no substitute — `revision` is silently empty, offering no information despite being "always kept." A consumer debugging "which model decided this" from a `--plan-only` log has nothing to go on. | §3, premortem, AT-5 |
| C-6 | `0023:§problem-statement`, `0023:Consequences (Scope of the guarantee)` | The RDR openly states the token-saving benefit is unmeasured for the actual traffic mix: "what share of the motivating consumer's traffic refuses" is "unsized here." The entire economic justification (the reason this RDR exists) rests on an assumption that success calls dominate the chained-call pattern, which is never verified against real traffic. | The flag ships, gets adopted, and the measured real-world token savings turn out to be a fraction of A1's synthetic 79.4%, because the actual consumer's chain refuses (retries, ambiguous matches) far more than the synthetic fixture models — undermining the "High priority" framing in front of stakeholders who were sold the 79.4% number. | §2, §3, premortem |
| C-7 | `0023:C1` (whole-tree walker text) | The oracle design requires reusing none of the existing recursive-walk infrastructure and hand-writing a *new* total walker that must correctly force-materialize `help` and `completion` via cobra's lazy-init functions (`InitDefaultHelpCmd`, `InitDefaultCompletionCmd` — the latter only reachable "from `ExecuteC`" per A5's own evidence). This is exactly the kind of cobra-internals-dependent test that breaks silently on a cobra version bump. | A future `go get -u` on cobra silently breaks the S4 oracle's precondition (the walked set no longer contains `completion`, or contains it via a different init path), and the oracle passes vacuously — the exact failure mode A5 was written to prevent, now reintroduced one dependency bump away. | §1, premortem, AT-4 |
| C-8 | `0023:Joint-decision check` | The joint-decision sweep found only 0024 as a fired collision on `resolvePayload`, but did not enumerate 0018, which also modifies `internal/cli/flow_resolve.go` (`kernelResolveFailure`, `escapeShapeBreaches`) in the same file, same RDR-open window. The RDR's own text treats "same file, disjoint symbol" as merely "context," but that standard was applied inconsistently — 0012 got a "context beside the fire" callout for a disjoint-symbol touch on the same file, while 0018 (equally disjoint, also touching `flow_resolve.go`, also currently Draft) is not mentioned at all. | A future reader auditing this RDR's joint-decision completeness finds an unlisted concurrent Draft RDR editing the same file, and has to manually re-derive whether the omission was deliberate (refusal-path vs success-path, truly disjoint) or an oversight — the RDR gives no evidence it considered 0018 explicitly. | §1, premortem |
| C-9 | `0023:A1`, `0023:Problem Statement` | A1's "verification" is a hand-approved spike table on synthetic fixtures the RDR's own authors built to match the target shape, not a byte-width measurement on the actual motivating consumer's real chained calls. The 1529→430 number that justified filing the kata is explicitly retired as unrepresentative, replaced by a new synthetic number (708→146) built for the RDR. This is evidence self-manufactured to fit the shipped design, not independent confirmation the real-world saving holds. | The org lands the flag, migrates the actual motivating consumer, and finds the real saving is meaningfully different from 79.4% because the real model's fact/gate ratio does not match either synthetic fixture — the number quoted in every retro and status update turns out to have been measured against a fixture built after the design was chosen, not before. | §3, premortem |

Ungrounded items: none — every ledger row cites a passage read directly from the projector output or the RDR file, and every code-side claim (struct field types, `NumField()==14` pin, `readersOf` fail-hard behavior, `walkCommandTree`'s name-skip, `NewRootCmd()`'s `AddCommand` calls, `help_all.go`'s force-init of `help`, 0018's touch on `flow_resolve.go`) was verified directly against `internal/cli/*.go` and the docs/rdr tree, not asserted from RDR prose alone.

---

## 1. The three most likely ways implementation goes wrong

### Failure 1: The whole-tree registration oracle (S4/A5) becomes the schedule's long pole, and it ships weaker than C1 requires

**Root cause in the RDR.** `0023:A5` is marked `Status: Pending` at Draft — not Verified, not even a completed spike — with the "Verification plan" reading: "write the total walker... against the real root and confirm it enumerates `help` and `completion`... Method: Spike, at implementation Phase 2." In other words, the single most structurally novel piece of this RDR (a new kind of oracle that has *no shipped exemplar anywhere in the repo*, per A5's own evidence) is deferred to be discovered during implementation, not proven before lock. The RDR is explicit that this is deliberate: "Phase 1 does not depend on it," and only lock/Phase-2 completion is gated on it.

**The specific passage that enabled it.** `0023:A5`'s evidence section is unusually candid about how much is *not yet known*: "the IMPLEMENTABILITY half was re-opened by the cove sweep... and its original evidence retracted: no shipped test walks `Commands()` recursively... the repo's only recursive walker... SKIPS children named `help` or `completion`... the whole-tree idiom C1 mandates has no shipped exemplar." That is a candid admission that nobody has built or run the thing C1 mandates. The Prerequisites section (`0023:§prerequisites`) then converts this Pending status into an accepted risk rather than a blocker: "A5's implementability half is discharged by writing the total walker itself, which is Phase 2 work... What it does gate is lock."

**The symptom the user (implementer/reviewer) will see.** Two sub-failures are plausible, both grounded in the actual cobra mechanics this RDR itself documents:

- The walker has to force-materialize `completion` via `InitDefaultCompletionCmd`, which per A5's own evidence is "creat[ed]... from `ExecuteC`" — i.e., normally only wired during actual command execution, not during a bare tree construction. Writing a test that calls cobra's private/internal init path correctly, in a way that survives a future cobra upgrade, is exactly the kind of test-infrastructure work that eats days, not hours, and that a "Phase 2, oracles" line item drastically undersells.
- If the walker proves awkward to write "in the house idiom" (A5's own phrase), the temptation under schedule pressure is to fall back to something closer to the existing `walkCommandTree` pattern with a *name-based allowlist* extended to cover `plan-only` specifically — which A5 itself says is a "strictly weaker negative" than what C1 demands. The result: an oracle that looks green, has a docstring citing C1, but does not actually cover the `help`/`completion` class the whole point of A5's cove correction was to force coverage of.

I verified directly: `internal/cli/help_all.go::walkCommandTree` does skip `help`/`completion` by name (confirmed at lines 343-351), and `NewRootCmd()` (`internal/cli/root.go`) constructs the tree via plain `AddCommand` calls with no `completion` subcommand registered — confirming A5's claim that a bare root lacks `completion` until `ExecuteC` runs. This is real, not a hypothetical: the RDR is asking Phase 2 to write and validate, for the first time in this codebase, a test that reaches into cobra's lazy command initialization internals. That is genuinely uncertain engineering work being scheduled as if it were a checkbox.

### Failure 2: The A2 "encoder mechanism" is undersold as a detail — it is a type-signature change to a heavily-read struct, and Phase 1 will run long

**Root cause in the RDR.** `0023:D-wire-byte-format` states: "The in-code mechanism (pointer-`omitempty` nilling vs a projected struct) is deferred to Resolve (A2); the wire result is fixed here." This frames the mechanism choice as a late, low-stakes implementation detail. But A2's own evidence describes the winning mechanism as "pointer-valued echo fields with `omitempty`, always populated in default mode, nilled at projection" — which means five fields on `resolvePayload` (`Model`, `Observed`, `Owned`, `Readers`, `Outcome`) change from concrete Go types to pointer types.

**The specific passage that enabled it.** The Source-authority census (`0023:§source-authority-census`) states plainly: "All fifteen are written at one assembly site (`internal/cli/flow_resolve.go` `resolvePayload{…}` literal) — there is no fallback arm and no second writer, which is what makes the partition a single-site edit." This "single-site edit" framing is doing a lot of load-bearing work to make the change sound small. I verified the struct directly (`internal/cli/flow_resolve.go:31-55`): today every one of these fields is a concrete type — `Model string`, `Observed map[string]string`, `Owned map[string]string`, `Readers []string`, `Outcome string`. None are pointers. Converting them to pointers means:
- The literal construction site (`payload := resolvePayload{...}`) must now take addresses of locals (`&req.modelRef` where `req.modelRef` is presumably a value, or a fresh local variable) for scalar fields, since Go does not let you take `&"literal"` or `&someFunctionCall()` directly in a composite literal — this typically forces the introduction of intermediate named variables purely to get addressable values.
- `decision_table_0010_test.go` pins `resolvePayload{}` `reflect.TypeOf(...).NumField() == 14` (verified at line 435) — the RDR's own A3 assumption says the *field count* survives, but every type assertion, JSON-shape assertion, or reflection-based test that inspects field *kinds* (not just count) is now touching a different `reflect.Kind` for five of fourteen fields. The RDR asserts NumField is stable; it never asserts or checks that no test inspects field Kind.
- Map- and slice-typed fields (`Observed`, `Owned` are maps; `Readers` is a slice) already have a natural nil-vs-empty distinction in Go's JSON encoder without needing a pointer wrapper for `omitempty` on non-empty-vs-absent — but the RDR's own A2 evidence flags "the bare non-pointer `omitempty` hazard reproduced (drops empty `{}`/`[]`)" as the very problem pointer-wrapping is meant to solve. That means the mechanism has to special-case container fields (pointer-to-map, pointer-to-slice) distinctly from scalar fields (pointer-to-string), which is more surface area than "one mechanism" suggests.

**The symptom the user will see.** Phase 1, sized in the plan as "Register `--plan-only` on `resolve` only and hand `respond.OK` the projected verb result when set," runs long because the actual work is a struct-wide type migration with knock-on edits at every read site, plus new intermediate locals purely for addressability — the kind of mechanical, easy-to-get-subtly-wrong-once (e.g., forgetting to nil a pointer field on one code path, or nil-dereferencing it in a future reader) refactor that eats a day of debugging for what the plan sizes as "one flag." Compilation and the existing test suite will likely catch outright breakage, but a **shape** bug — a pointer field that renders differently from the non-pointer original in some non-obvious code path (e.g., something computing `len(payload.Observed)` that now needs a nil-check first) — is exactly the class of bug this RDR's own A2 evidence admits happened once already during the spike ("the bare non-pointer `omitempty` hazard reproduced").

### Failure 3: The text-mode strict-width MUST (A8, S5) is unverified at Draft and may not hold, forcing a late contract change mid-implementation

**Root cause in the RDR.** `0023:A8` is `Status: Pending` with "Evidence: not yet measured," yet `0023:C1` already states as a MUST: "the projected encoding STRICTLY SHORTER than the default... it is asserted in BOTH output modes." A normative MUST is locked into the contract before the assumption that makes it satisfiable has been checked.

**The specific passage that enabled it.** A8 itself flags the risk candidly: "The claim is very likely true — the six dropped lines... carry non-empty content — but C1 now asserts it as a MUST in both modes, so it is booked rather than assumed... **If wrong**: C1's both-modes width clause is unsatisfiable as written and the clause narrows to JSON, leaving S5's subset assertion as text mode's only width guard." This is the RDR's own author acknowledging the contract may need to shrink after lock — which is precisely the kind of finding a pre-lock gate exists to prevent from reaching implementation.

**The symptom the user will see.** Phase 2 (Oracles) implements S5, measures actual text-mode byte widths per MVV step 6, and either (a) confirms the width savings hold on every tested fixture — in which case this was correctly Pending, just unlucky in review timing — or (b) finds a fixture shape (plausible: a model with very few facts/gates but a large `writes`/`next` block, i.e., the "release begin" shape the RDR's own A1 table shows saving only 42.2%, tested here in *text* rather than JSON where line-overhead-per-key differs) where the text rendering does not strictly shrink. In case (b), the team discovers mid-Phase-2 that a MUST clause in a "locked" contract cannot be satisfied as written, and has to either narrow C1 (an amendment to Normative Contracts after other work has started against it) or invent a special case. Either path is schedule risk introduced by locking a MUST before its enabling assumption was checked.

---

## 2. The one section that will be rewritten within 6 weeks of shipping

**`Consequences` / `Scope of the deliverable`** (`0023:§consequences`) will be rewritten first, and specifically the paragraph that reads: "this RDR ships the CAPABILITY, not its adoption... a green build delivers zero measured saving on day one... Adoption is not scheduled by this RDR and is not one of its phases."

This is the section that will not survive contact with reality, for a simple reason: the entire Problem Statement is an economic argument ("An agent-driven caller of `flow resolve` pays for every output byte as transcript tokens... raised to P1"), the entire justification for High priority and foundational-profile review overhead is that this problem is costing something *now*, and yet the RDR explicitly declines to schedule the one thing that would actually reduce that cost — a consumer migration. Six weeks after this ships, someone will look at the P1 kata that motivated it (`intrastate#srz2`), find it still open or still costing tokens in production, and ask "didn't we fix this?" The answer will be "we shipped the mechanism, but nobody adopted it," and the natural response is to either (a) retroactively write an adoption/migration plan into this RDR's Implementation Plan — which the RDR's own rule against amending locked RDR content ("we never amend rdr content... code is the source of truth") would forbid, forcing a *new* RDR or kata just to schedule the obvious next step — or (b) discover that the "natural first adopters" named in Consequences (kata `rg0e`, "the consumer that dropped its status call") were never actually committed to adopt, and the sentence naming them was aspirational, not a plan. Either way, the "Scope of the deliverable" framing — clever as a scoping device to keep this RDR's blast radius small — is the section whose optimism about adoption happening on its own will be falsified first, and loudest.

---

## 3. The one assumption that will not survive first contact with a real user

**A7's framing of `model` as a pure echo field, married to C2's decision to keep `model` out of the always-keep core** is the assumption that breaks first.

A7 argues, correctly on its own narrow terms, that `model`, `observed`, `owned`, `readers`, and `outcome` are all "the caller's own input echoed back" and therefore safe to drop under `--plan-only` because "the caller that opts in is the caller that already holds" that data. The RDR is aware this could hurt an adopter — Risks and Mitigations even calls it out directly: "a consumer treats plan-only output as the full record and loses the model reference — `model` is the projected-away field that actually names which model produced the plan, since `revision` is constant-empty today." The RDR's mitigation is: "opt-in with an unchanged default; the help line states the omitted group by name; the caller that opts in is the caller that already holds the request, `model` included."

That mitigation assumes every caller of `--plan-only` retains, correlates, and logs its own `--model` argument alongside the projected response at the point where it later needs to explain "which model produced this plan." In practice, the first real user of this flag is exactly the chained-call pattern the RDR was built for: an agent that calls `resolve` once, gets a plan, and later (in a *different* turn, possibly a different process, possibly after the model reference has scrolled out of context or been summarized away) needs to answer "was this the pricing model or the eu-region model?" The RDR's own economic argument is that **transcript tokens are re-paid on every turn** — meaning the caller's own copy of `--model` is exactly the kind of thing an agent under token pressure will *not* keep re-quoting either. The assumption that "the caller already holds it" quietly assumes the caller behaves like a well-disciplined batch script with the original invocation in scope, not like the token-starved chained-call agent this RDR exists to serve. The RDR's own always-keep design leaves `revision` as the nominal identity slot for exactly this need, then documents in the same breath that `revision` is "presently constant-empty" and "no model can declare one" — so the actual safety net for "which model produced this" is empty for the foreseeable future. The first real user who logs a projected plan without separately retaining the model argument will hit precisely the failure mode Risks and Mitigations predicted and shrugged off.

---

## 4. Premortem

*Six weeks post-ship. Written as if it already happened.*

`--plan-only` shipped, all five oracles (S1–S5) are green, `make check` is clean, and the CHANGELOG says "flow resolve: add --plan-only for chained calls." Here is what actually happened.

**Week 1–2 (Phase 1 slips).** The implementer starts Phase 1 expecting "one flag, one projection function." They hit A2's mechanism head-on: converting `Model`, `Observed`, `Owned`, `Readers`, `Outcome` on `resolvePayload` to pointer types breaks addressability at the `payload := resolvePayload{...}` composite literal in `runFlowResolve` (`internal/cli/flow_resolve.go`), because `req.modelRef` and `outcome` are not addressable expressions inline. Four new named locals get introduced purely to take their addresses. `decision_table_0010_test.go`'s `NumField() == 14` pin stays green (A3 holds on count), but a second, previously-unnoticed test elsewhere that does `reflect.TypeOf(payload.Model).Kind() == reflect.String` (or the moral equivalent — a JSON-shape fixture comparison that assumed a concrete type) breaks, and nobody had inventoried "every test that touches `resolvePayload`'s field *kinds*, not just its field *count*" before starting, because A3's Verified claim was scoped to "no test asserts an echo field's presence" — not to type reflection. Two days lost.

**Week 2–3 (A5's oracle turns out to be the real project).** Phase 2 begins the whole-tree registration walker. The implementer discovers, as A5's own evidence already flagged, that `completion` only exists on a tree that has gone through `ExecuteC` — but `ExecuteC` also *executes* the command, which is not what a registration-audit test wants to do (it wants to enumerate, not run, `os.Exit`-triggering subcommands). Getting `InitDefaultCompletionCmd` invoked in isolation, safely, in a unit test, without triggering full command execution side effects, turns into a half-day spelunking exercise through cobra internals that the RDR's Implementation Plan sized as one bullet in a five-oracle Phase 2 list. The walker ships, but under time pressure the "no name-skip, no Hidden gate" requirement gets satisfied literally while the surrounding test harness quietly special-cases `completion` with a manual `cmd.InitDefaultCompletionCmd()` call bolted on right before the walk — technically compliant with C1's letter, but nobody re-derives whether a *future* auto-generated cobra command (say, if cobra ships a new default subcommand in a later version) would be caught by this walker or silently missed the way `completion` originally was. The oracle guards against yesterday's blind spot, not tomorrow's.

**Week 3 (A8 breaks as predicted).** MVV step 6 / S5 measures actual text-mode byte widths. On the "release begin" fixture (gate+writes shape, the one A1's own JSON table already shows saving only 42.2% — the smallest margin measured), the text-mode rendering saves bytes but by a much smaller margin than the JSON case, and on a fixture nobody tested — a model with a single gate, no writes, minimal `next`, but a long `observed` tag set — the projected text version is *not* strictly shorter, because the six dropped echo lines are short (`model:`, single-tag `observed.*` lines) while the retained plan group (`gates[]` empty-line placeholder, `writes{}`/`clear[]` placeholders) renders the same "(none)" boilerplate lines `flatten` already emits for empty containers regardless of projection. C1's both-modes strict-width MUST goes red. The team narrows C1 exactly as A8's own "If wrong" clause predicted — but this happens as an urgent Slack thread mid-Phase-2, not as a calm pre-lock decision, and the fix is to add a special-case exclusion for "all-empty-plan" shapes that the original contract text never anticipated.

**Week 4–6 (the actual point of the RDR is not realized).** The flag ships. Every oracle is green. `intrastate#srz2`, the P1 kata that started this whole RDR, is still open, because the consumer whose skill prompt was measured in the Problem Statement never got migrated — the RDR explicitly declined to schedule that ("Adoption is not scheduled by this RDR and is not one of its phases"). A garden-pass review six weeks later asks "did we fix the token cost problem," and the honest answer is: we built the mechanism, we did not use it. The 79.4% saving quoted everywhere in this RDR's evidence and Performance Expectations was measured on a *synthetic fixture the RDR's own authors built after the design was chosen* (A1's `evidence/spikes/a1-byte-width.md`, S1b), not on the real consumer's real chained call. Nobody has yet measured what the real saving is, because nobody adopted the flag. The org either files a follow-up kata to actually wire the flag into the motivating consumer's prompt (the work this RDR's economic argument was entirely about, and the work it did not schedule), or quietly treats "capability shipped" as "done" and moves on, leaving the P1 problem effectively unaddressed while the CHANGELOG says otherwise.

Separately, a consumer who *did* opt into `--plan-only` early (the seam-mate `rg0e` ask) hits the `model`-field loss predicted in §3: their multi-turn audit log shows a plan with no model reference and an empty `revision`, and they cannot answer "which model decided this" from the projected record alone — exactly the gap Risks and Mitigations named and waved off with "the caller already holds it."

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

**AT-1 (targets C-4, C-6, C-9 — unmeasured/unscheduled adoption and self-manufactured evidence).**
```
Scenario: The RDR names a specific consumer and commits to measuring
  real-world savings before claiming the economic goal is met
  Given the Problem Statement cites a specific consumer's chained-call
    pattern as the motivating cost
  When the Implementation Plan is reviewed
  Then it MUST name an owner and a phase (even if a follow-up kata)
    for migrating at least one real consumer call site
  And the RDR MUST NOT claim "done" or close the tracking kata until
    a real (not synthetic) before/after byte measurement is recorded
    against that consumer's actual traffic
  And A1's evidence MUST include at least one measurement against
    real historical call logs, not solely fixtures authored for the
    RDR after the design was already chosen
```
This would have forced either an adoption phase into the Implementation Plan, or an explicit, reviewed decision to defer adoption with a named follow-up tracked at the same priority as the original kata — instead of the current silent "capability, not adoption" scoping that reads as complete but leaves the P1 problem unsolved.

**AT-2 (targets C-2 — undersized Phase 1 due to the pointer-type migration).**
```
Scenario: The A2 mechanism's blast radius on resolvePayload's readers
  is enumerated before Phase 1 is sized
  Given A2 proposes converting five resolvePayload fields from
    concrete types to pointer types
  When the Implementation Plan sizes Phase 1
  Then it MUST list every existing call site and test file that reads
    resolvePayload fields by reflection, type assertion, or direct
    field access (not just those asserting "presence")
  And it MUST confirm no test inspects reflect.Kind() of the affected
    fields, not merely NumField()
  And Phase 1's estimate MUST account for introducing addressable
    locals at the construction site, not just "hand respond.OK the
    projected form"
```
This would have surfaced the addressability problem and the reflect.Kind() blind spot in A3's "no predecessor oracle moves" claim before Phase 1 began, rather than discovering it mid-implementation.

**AT-3 (targets C-3 — A8 Pending with a MUST already locked on it).**
```
Scenario: No MUST clause depends on a Pending assumption at lock
  Given C1 asserts a strict-width MUST in both JSON and text modes
  And A8 (text-mode width) is Status: Pending with no evidence
  When the Finalization Gate's Assumption Verification runs
  Then it MUST flag: "Status consistency: no assumption marked
    Pending... may have settled-fact prose elsewhere in the RDR
    depending on it" (the gate's own stated rule)
  And C1's text-mode width clause MUST be measured (or narrowed to
    "JSON mode only, text mode best-effort") before Status flips to
    Final, not left to resolve during Phase 2
```
The Finalization Gate template *already contains this exact rule* ("no assumption marked `Pending`... may have settled-fact prose elsewhere... depending on it") — this test would have caught the gate's own rule being violated by C1 depending on Pending A8, and forced either an A8 spike before lock or an explicit narrowing of C1's MUST to JSON-only until proven.

**AT-4 (targets C-1, C-7 — the whole-tree walker's real implementability and future-fragility).**
```
Scenario: The whole-tree registration oracle is proven against the
  real command tree before lock, not scheduled as Phase 2 discovery
  Given A5 proposes a total walker that must materialize help and
    completion via cobra's lazy init functions
  When the RDR is reviewed for lock
  Then a spike MUST demonstrate the walker running against the real
    NewRootCmd(), asserting the walked set contains both help and
    completion, without invoking full command execution (ExecuteC)
    side effects
  And the spike MUST record what cobra API guarantees (documented,
    not incidental) the walker relies on for completion's
    materialization, so a future cobra version bump has a named
    contract to break against rather than a silent test regression
```
This converts A5 from "Pending, verify during Phase 2" into a pre-lock spike — exactly the discipline the RDR applies to A1 and A2 (both Verified via spike before lock) but withholds from A5 for no stated reason other than scheduling convenience.

**AT-5 (targets C-5 — `model` field loss on real chained-call audit trails).**
```
Scenario: A --plan-only caller can still identify which model produced
  a plan without separately retaining the --model argument
  Given a caller issues `flow resolve --plan-only` and logs only the
    response body (the realistic behavior of a token-constrained
    agent, per this RDR's own Problem Statement)
  When the caller later needs to determine which model produced the
    plan
  Then the projected payload MUST supply enough information to
    answer that — either a live (non-empty) `revision` value, or
    `model` moved into the always-keep core
  And if neither is true, the RDR MUST document this as an accepted
    gap with a named consumer obligation ("log --model yourself"),
    not merely assert "the caller already holds it" as a mitigation
```
This would have forced the RDR to confront that its own mitigation for the `model`-loss risk assumes caller behavior that contradicts the RDR's own stated motivation (agents drop context to save tokens) — and either move `model` into always-keep, prioritize `revision`'s activation, or explicitly accept and document the gap rather than wave it off in a Risks bullet.
