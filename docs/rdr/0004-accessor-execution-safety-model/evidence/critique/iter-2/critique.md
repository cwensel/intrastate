Model: claude-sonnet-5

# Critique — RDR 0004 Accessor Execution Safety Model (re-entry iter-2)

## 6. Findings Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | A9 Evidence: "the existing Resolve spike does not witness any of the four [rules]"; Disposition Table rows tagged "not witnessed — see A9" (`read_back_incomplete`, timeout-outranks-`incomplete_read`, missing/empty requested-key-set) | Four normative rules ship to lock with zero executable evidence — only prose and an MVV promise — while A9 is graded "Pending" and the RDR proceeds to lock anyway | Implementation discovers `read_back_incomplete` can't be cleanly distinguished from `incomplete_read` at a real binding (see C-3); a rule that reads as settled in the normative-contracts section turns out to be the least-tested clause in the document | §Critical Assumptions A9 |
| C-2 | Load-Bearing Decision "Absence crosses the seam as omission": "how the accessor layer represents an absent value internally is free, but what reaches the resolver must leave the key absent from the owned snapshot" | The rule constrains the accessor→resolver *seam* but never specifies the accessor *definition* struct or the read function's *return type* that would make "omit the key" the only representable option; nothing stops a Phase 1 implementer from choosing `map[string]string` with a sentinel, which is exactly the encoding the clause exists to forbid | An owned key the artifact genuinely lacks gets read back as present with a placeholder value; `owned_state_unavailable` silently stops firing for that key; a transition plan gets computed and applied against state that was never actually there | §Load-Bearing Decisions; grounded against `internal/resolve/resolve.go::missingOwned`, `TagSet.has` |
| C-3 | Normative: "A read that cannot resolve every requested key MUST be reported as `incomplete_read`... When a read exceeds its timeout before resolving every requested key, `timeout` takes precedence over `incomplete_read`." | The precedence rule presumes an executor that can observe "some keys resolved, deadline hit" as one atomic state. Real bindings (a TOML parse, an HTTP GET, a single `os.ReadFile`) fail file-wide or request-wide, not key-wide — the spike's `unreadable map[string]bool` is a hand-authored oracle, not a derived one. There is no specified mechanism by which an implementation *detects* "N of M keys resolved, then timeout fired" as opposed to just "timeout fired, 0 keys examined" | Two operationally distinct failures (call that timed out having read nothing vs. call that read most keys and stalled on the last one) both surface as bare `timeout`, giving the operator no way to tell "retry it" from "this binding cannot do partial reads, redesign the accessor" | Oracle Discriminability row 7; MVV Scenario 7; 3amigo T-8 |
| C-4 | Normative: "One unreadable requested key MUST refuse the whole read; the executor MUST NOT return the keys that did resolve." + Load-Bearing "Read refusal granularity" | Whole-read refusal on any single unreadable key means one accessor definition that requests N keys is only as reliable as its least-reliable key. For an accessor role backed by a live external system (RDR mentions "External API accessors" in A5/Cross-Cutting), a single flaky field poisons every other field in the same read, even fields with no relationship to the flaky one | A transition that only needs `status` refuses `owned_state_unavailable`-adjacent (`incomplete_read`) because an unrelated key like `last_reviewed_by` in the same requested set happened to be unreadable this one time | §Load-Bearing Decisions "Read refusal granularity"; 3amigo T-9 |
| C-5 | Load-Bearing "Absent vs unreadable is the binding's call, defaulting to unreadable": "a key is *absent* only when the binding read the artifact successfully and the key was not there. Every other outcome... is *unreadable*" | This pushes the entire safety-critical absent/unreadable classification onto each binding author's judgment call, with no interface contract, no shared helper, and no validation-time check that a binding actually implements the distinction correctly. The RDR's own A8 evidence line concedes the spike "declares unreadability as a fixture field rather than deriving it" | Two different accessor bindings (say, a TOML-file reader and an HTTP-JSON reader) implement the absent/unreadable line differently; one is conservative (reports unreadable on any doubt) and one is optimistic (reports absent on any doubt); the same underlying artifact condition produces different resolver dispositions depending only on which binding happens to read it | §Load-Bearing Decisions; 3amigo T-5 |
| C-6 | Normative: "Every read accessor definition MUST declare the requested key set as validated metadata... A missing or empty requested key set MUST fail validation" | The spike's own default/fallback path (`expectedTagKeys`, cited by name in the RDR text itself as "exactly that circular shape — fixture convenience, not the contract") is the *only* implemented read path exercised by 10 of 13 transcript lines (all non-explicit-key calls, plus both replay runs). The explicit, contract-shaped path is exercised by exactly 3 lines. The disposition-table witness for A4 replay stability rides the circular path, not the normative one | Replay stability (A4, "Verified") is asserted, but the actual evidence covers a code path the RDR itself disowns as non-normative; nothing has tested whether replay is stable when the true requested-key-set path is used, which is a materially different code path (validated metadata lookup, not artifact introspection) | 3amigo T-13; §A4 Evidence vs. §A9 Evidence |
| C-7 | Failure Modes: "There are three silent-failure shapes, each with a mandatory guard" | The enumeration is reactive, not structural — it lists exactly the three failure shapes the last two review rounds happened to catch (write mismatch, truncated read, absent-as-present). There is no argument that this is the *complete* set of silent-failure shapes for a system whose executor is still unbuilt; the read-back-re-read-hits-unreadable-key case (`read_back_incomplete`) is a fourth near-silent shape (it could easily degrade to reporting success, since the pre-A9 code had no branch for it at all) that had to be discovered by a review pass, not derived from a stated method | A future accessor interaction the current three review passes didn't think to construct (e.g., a gate accessor whose indeterminate reason string is empty, or a write whose read-back reads a *different* artifact revision than the one written) ships as a fourth silent-failure shape nobody enumerated | §Failure Modes; premortem |
| C-8 | Proportionality: "Read completeness is part of that same contract... so carrying it here does not widen the RDR... it constrains this RDR's own output rather than reopening RDR 0001's `Input` shape, which is unchanged" | This claim is doing real defensive work — it's arguing the RDR doesn't need to touch `resolve.Input` — but C-2 shows the seam clause has no teeth without pinning something about the accessor's return type, which is the only artifact this RDR actually owns. The claim that "`Input` is unchanged" is true and irrelevant; the open question is whether the accessor package's return type can express "absent" without a sentinel, and that question is unanswered, not scoped-out | Implementation team treats "Input is unchanged" as license to design the accessor return type however is convenient, discovers late that a sentinel was the convenient choice, and now must choose between reopening this RDR or reopening RDR 0001 | §Proportionality; 3amigo T-1, T-11 |
| C-9 | MVV Scenario 6: "the resolver refuses `owned_state_unavailable` naming that key. This is the absence half of read completeness" | Scenario 6 tests only the *positive* case (an absent key does trigger the refusal). It has no negative control proving an implementation that used a sentinel value would *fail* this scenario — if the fixture happens to route the sentinel string through a code path where `TagSet.has` still returns false (e.g., a naive stub resolver in the MVV harness rather than the real accessor→`Input.Owned` wiring), the scenario passes vacuously while the production wiring still uses a sentinel | The MVV scenario passes at Resolve time; the actual CLI-integration wiring (RDR 0005, a different RDR/phase) is where the sentinel-vs-omission choice really gets made, outside this scenario's blast radius, and nothing forces that later wiring to be checked against this scenario's assumption | §Validation Scenario 6; Oracle Discriminability row 6 |

## 1. The three most likely ways implementation goes wrong

### 1a. The absence-as-omission rule gets implemented as a sentinel, because the RDR never pins the return type

**Root cause.** The Load-Bearing Decision "Absence crosses the seam as omission" states the *outcome* required at the seam ("what reaches the resolver must leave the key absent from the owned snapshot") but explicitly refuses to constrain the *mechanism*: "how the accessor layer represents an absent value internally is free." That freedom is fine internally — the danger is that the RDR also never specifies the shape of the accessor package's *external* return value (the thing that becomes `Input.Owned`). The spike returns `map[string]string` with a string sentinel `"<absent>"` for exactly this case, and the RDR spends a full paragraph disclaiming that the sentinel is "fixture shorthand, not that seam value" — but never says what the real value is instead.

**Specific passage.** "The accessor layer MAY represent absence however it chooses internally, but what crosses the seam MUST leave the key absent from the owned snapshot" (Normative Contracts) combined with "The spike's `<absent>` string is fixture shorthand, not that seam value" (Load-Bearing Decisions) — both true, neither constructive. Nowhere in Phase 1 ("Define the accessor definition structs...") or Phase 2 ("Implement context-bound invocation for typed accessor bindings. The executor returns structured success/refusal values") is the read accessor's success-value type specified as anything other than "typed tag values," which is exactly ambiguous enough to admit `map[string]string{"profile": "<absent>"}` as a compliant implementation, since that map still "carries typed tag values."

**Symptom.** A key the artifact genuinely lacks reads back into `Input.Owned` as present with a placeholder. `TagSet.has` returns true. `missingOwned` never reports it. `owned_state_unavailable` — the refusal this whole RDR exists to keep meaningful — stops firing for exactly the case A8/A9 were written to protect. A transition executes against state that was never read. The user sees a plan applied (a `status: Final` write, say) that should have been refused, and only notices when the artifact's real state and the tool's model of it diverge downstream — which, per this RDR's own Fidelity Table, is a divergence this system explicitly does not check for at the artifact-byte level.

### 1b. `read_back_incomplete` collapses into either `incomplete_read` or `read_back_mismatch` because it has no witnessed shape at all

**Root cause.** Every other refusal class in this RDR has a spike line: `timeout` at `output.txt:5`, `capability_mismatch` at `output.txt:7`, `read_back_mismatch` at `output.txt:9-10`, `incomplete_read` at `output.txt:12`. `read_back_incomplete` has none — the RDR says so itself, three times ("not witnessed — see A9"). It is pure prose sitting in a document whose entire methodology (the Method vocabulary table, the self-reference rule, the exactness-claim rule) exists to prevent exactly this: a normative claim with no executable evidence. The write function in the spike (`main.go::write`) re-reads by cloning the in-memory Go map (`observed := clone(art.tags)`), which can never fail to read a key — there is structurally no way for the current spike to produce this disposition even by accident.

**Specific passage.** "The read-back re-read is subject to read completeness. If it cannot read a key it must compare, the write MUST be reported as `read_back_incomplete`... and MUST NOT be reported as `read_back_mismatch`... nor as success" (Normative Contracts) plus A9's own admission: "the existing Resolve spike does not witness any of the four [rules]... it re-reads by cloning the tag map without going through `read`."

**Symptom.** When implementation actually wires the write path's re-read through the same `read` accessor function (as the normative clause requires — "subject to read completeness"), the engineer hits a three-way ambiguity the RDR never resolved with evidence: is an unreadable compared key during read-back reported through the *read* accessor's own `incomplete_read` refusal (reusing the class), or does it need a distinct `read_back_incomplete` code path threaded through the write accessor's result type? The RDR asserts the latter but the write accessor's `result` type in the spike has no field for it, and no MVV test exists yet to force the distinction. The most likely outcome: an engineer under time pressure treats an unreadable read-back key as a `read_back_mismatch` (it's already coded, it's adjacent, "close enough") — precisely the miscoding the normative clause forbids by name ("MUST NOT be reported as `read_back_mismatch`, which asserts the artifact is wrong"). The user sees a write reported as data corruption ("read_back_mismatch") when the actual problem was the verification step itself failing to read, e.g., a transient permission or network blip on the re-read — a false "your state is corrupted" alarm on an otherwise-successful write.

### 1c. Whole-read refusal on one bad key turns flaky external accessors into permanently-refusing accessors

**Root cause.** "One unreadable requested key MUST refuse the whole read; the executor MUST NOT return the keys that did resolve" is stated as a "conservative choice," but the RDR's own Load-Bearing Decisions concede it is "indistinguishable from an artifact of the spike's early return" — i.e., nobody actually decided this against an alternative with named trade-offs; the spike's control flow (`if art.unreadable[key] { return &result{...refusal: refusalIncompleteRead} }` inside a loop, an early return) produced the rule, and the rule was then written up as if it were a decision. Combined with A5's explicit scope-in of "External API accessors," this rule applies to accessors backed by real network calls, not just fixture maps.

**Specific passage.** "Read refusal granularity — one unreadable requested key refuses the entire read rather than returning the keys that did resolve. This is the conservative choice... a decision, not an artifact of the spike's early return" (Load-Bearing Decisions) — the text protests too much; the very next sentence in 3amigo T-9 (which fed this revision) flagged exactly this concern and the RDR's response was to assert the opposite of the finding rather than add evidence for it.

**Symptom.** For any accessor whose requested key set spans multiple underlying calls or fields (the RDR requires "declared metadata," not "one key per accessor," so this is the expected shape, not an edge case), a single consistently-flaky field turns the *entire* accessor into `incomplete_read` on every invocation, even though the other N-1 keys are perfectly readable every time. Downstream, that's `owned_state_unavailable` cascading from the resolver, on every run, for a transition that doesn't even depend on the flaky field. The user's experience: a flow that "used to work" now refuses unconditionally, and the diagnosis trail (per Failure Modes: "accessor identity, capability, artifact role, timeout, and expected versus observed tag values") gives no visibility into *which* key was the poison one, because the RDR only requires the refusal to "name the requested keys it could not read" — which is at least present, but the fix path (split the accessor into more granular bindings) is a re-authoring of the transition model, not a runtime remedy, so the user is stuck until an author intervenes.

## 2. The one section that will be rewritten within 6 weeks of shipping

**Load-Bearing Decisions → "Absent vs unreadable is the binding's call, defaulting to unreadable."**

This is the section that will be rewritten, and soon, because it is the only place in the RDR that hands a safety-critical binary classification to every future binding author individually, with no shared type, no validation-time check, and no test harness that can catch a binding that gets it wrong. Every other safety property in this RDR (capability mismatch, timeout, read-back mismatch, whole-read refusal) is enforced *centrally* by the executor — the executor owns the timeout, the executor owns the read-back comparison, the executor owns capability checking. Absent-vs-unreadable is the one property this RDR delegates to N different, independently-written binding implementations, each of which has to correctly answer "did I establish this key is missing, or did something merely fail?" for its own transport (TOML parse errors, HTTP partial responses, type coercion failures — the RDR names all three as unclear cases and does not resolve any of them). The first real binding (a TOML file reader, almost certainly, given `intrastate.toml` is already in the codebase per the sibling-path check) will hit a case the RDR's binary rule doesn't cover cleanly — a TOML table that's present but the specific key inside it is absent, versus a TOML *parse error* on the whole file, versus a key present with the wrong type — and whoever implements accessor binding #1 will either invent a convention that RDR 0004 never sanctioned, or come back to amend this section with the missing taxonomy. Given the number of times this document already says "this is the binding's call" without giving the binding a place to record that call in a way the executor can validate, this is the load-bearing gap most likely to force a rewrite.

## 3. The one assumption that will not survive first contact with a real user

**A9's implicit assumption that "the seam/boundary rules... are implementable as stated" as a single verified-or-not unit.**

A9 bundles four genuinely independent claims — seam omission, validated requested-key set, timeout-outranks-incomplete_read, and `read_back_incomplete` — under one Pending/Verified status. The RDR's own Evidence line for A9 already shows these four have wildly different implementation difficulty: seam omission is a type-design problem (solvable at Phase 1, no runtime uncertainty); validated requested-key set is a validation-rule problem (solvable, mechanical); timeout-outranks-incomplete_read requires an executor that can observe partial progress at the moment a deadline fires, which most real I/O primitives (a single blocking file read, a single HTTP round-trip) cannot report — the whole call either finished or it didn't, there is no "3 of 5 keys resolved so far" internal state to consult in the common case; and `read_back_incomplete` requires threading read-completeness semantics through a completely different code path (write, not read) that the spike doesn't even attempt.

The first real user-facing accessor — almost certainly a TOML file read, since that's the only I/O primitive already in the codebase (`internal/cli/config::Load`, cited in the Existing Infrastructure Audit) — will read the whole file in one `os.ReadFile` + one parse call. That operation is atomic: it either returns a fully-parsed document or an error. There is no way for such a binding to time out "partway through resolving keys," because there is no per-key resolution step at all — the whole artifact is parsed at once. The timeout-outranks-incomplete_read rule was written for an accessor architecture (key-by-key resolution with a deadline checked between keys) that the RDR's own primary example — a TOML-backed `state` artifact reading `status` and `profile` — does not actually have. When the first real binding is built, this precedence rule will turn out to be either vacuous (the case it's precedent for structurally cannot occur for file-backed accessors) or will force an artificial per-key polling loop into a binding that would otherwise be a single atomic read, purely to give the timeout clause something to preempt.

## 4. The Premortem

*Six weeks after RDR 0004's Phase 1-4 ship, engineering opens a retro titled "Why did the `flow advance` command silently promote a spec with no reviewer."*

The accessor package (`internal/accessor`, following the naming this RDR settles on: "read accessor," "gate accessor," "write accessor") was built against this RDR's normative contracts. The read accessor's success type is `map[string]string` — the most natural Go shape for "typed tag values," and the one the spike itself models — because Phase 1 ("Define the accessor definition structs, capability enum, refusal classes") never specified anything more constrained, and the RDR's own Load-Bearing Decision explicitly declined to pin the internal representation ("free... but what crosses the seam..."). The engineer who wired the accessor's output into `resolve.Input.Owned` (RDR 0001's `Input` struct, unchanged per this RDR's own Proportionality argument) needed *some* way to carry "this key was requested but the artifact doesn't have it" through a `map[string]string`. Two choices were live: omit the key from the map (matches the RDR's seam clause) or carry a sentinel (matches the spike's fixture code, which is the only executable reference implementation anyone had to copy from). Under deadline pressure, and because the spike — the only artifact anyone actually ran — uses `"<absent>"` as a real return value, the engineer copied the spike's shape into the "real" code, reasoning (correctly, on the text) that "the spike proves the branch machinery."

For weeks this was invisible: every fixture and integration test used artifacts where the required owned keys were always *present*. The `reviewer` field on a `spec.toml` artifact was optional in early adopters' files, and one team's spec template didn't set it. `flow advance` read that artifact through `spec.read`, got back `{"status": "Draft", "reviewer": "<absent>"}`, and passed it straight to `resolve.Input.Owned`. `TagSet.has("reviewer", ProvenanceOwned)` returned true — the map has the key. `missingOwned` never fired. The transition row requiring `reviewer` to gate promotion evaluated its guard against the literal string `"<absent>"`, which didn't match any configured guard predicate, so the row was pruned as GuardFalse rather than refused as `owned_state_unavailable` — an entirely different, non-alarming code path (per `resolve.go`'s own `gate` function: a pruned row "contributes nothing" and is invisible in the refusal). With that row pruned, a *less specific* row with no reviewer-gate matched instead, and the spec advanced from Draft to Final with no reviewer ever having been assigned.

Nobody saw a refusal. Nobody saw an error. `respond`/`clierr` never fired because no refusal was ever minted — the resolver's `Result` was a clean `Plan`, exactly as designed for the happy path. The first human to notice was a downstream consumer of the "Final" spec, three weeks later, asking who reviewed it. The incident review traced it to `internal/accessor/read.go`'s literal string `"<absent>"`, sitting in the same position the spike put it, and to RDR 0004's Load-Bearing Decision that called this representation choice "free."

The fix took one line (omit the key instead of writing the sentinel) and one very long conversation about why the RDR that was supposed to prevent exactly this — A8's whole argument is "the kernel cannot recover the distinction downstream" — didn't also pin the one type declaration that would have made the wrong choice impossible to write.

## 5. Acceptance tests that would have caught each failure at RDR-review time

```gherkin
Feature: Absence representation is structurally forced, not merely instructed

  Scenario: The read accessor's success return type cannot represent a sentinel
    Given the accessor definition struct and read-accessor return type from Phase 1
    When a requested key is genuinely absent from the artifact
    Then the return type MUST have no code path that inserts a string, enum, or
      other in-band value for that key into the returned tag set
    And the only representable way to signal absence MUST be omitting the key
      from the returned collection
    # Kills C-2 / C-8 / premortem: forces the type design question the RDR
    # explicitly declined to answer, before any binding is written against it.

  Scenario: A sentinel-based implementation fails the MVV
    Given a deliberately-wrong read accessor implementation that returns
      "<absent>" as a map value instead of omitting the key
    When its output is passed through the real accessor-to-resolver adapter
      (not a test stub) into resolve.Input.Owned
    Then TagSet.has for that key MUST report false
    And the sentinel-based implementation MUST fail this assertion
    # Kills C-9: MVV Scenario 6 as written has no negative control; this one
    # does, and it must run through the REAL adapter, not a harness stand-in.
```

```gherkin
Feature: read_back_incomplete has an executable witness before lock

  Scenario: Write read-back verification hits an unreadable compared key
    Given a write accessor that reports command-level success
    And the read-back re-read is executed through the same "read" function
      used by ordinary read accessors (not a raw map clone)
    And one of the keys the read-back must compare is unreadable on the re-read
    When the write accessor returns its disposition
    Then the refusal class MUST be "read_back_incomplete"
    And it MUST NOT be "read_back_mismatch"
    And it MUST NOT be reported as success
    # Kills C-1 / C-3 (partially) / failure 1b: forces the write path's
    # re-read through the SAME code as the read accessor, closing the gap
    # where write's re-read (main.go::write) clones a map and structurally
    # cannot exercise this refusal.
```

```gherkin
Feature: Timeout-vs-incomplete_read precedence is tested against a realistic binding shape

  Scenario: An atomic (non-key-granular) read binding times out
    Given a read binding that resolves its artifact in one atomic operation
      (e.g. one file read + one parse, not a per-key loop)
    And that operation exceeds its declared timeout
    When the read accessor returns its disposition
    Then the refusal class MUST be reported
    And the RDR's precedence rule MUST be demonstrated meaningful for this
      binding shape, not only for a hand-rolled per-key-loop fixture
    # Kills C-3 / A3 assumption failure: forces the precedence clause to be
    # checked against the binding shape the codebase actually has today
    # (TOML file config, atomic reads), not only a fixture built to make
    # the rule demonstrable.

  Scenario: A multi-key accessor with one permanently-flaky key
    Given a read accessor definition requesting keys A, B, and C
    And key B is unreadable on every invocation, while A and C always succeed
    When the read accessor is invoked repeatedly
    Then every invocation MUST refuse with incomplete_read naming B
    And the refusal MUST make clear that A and C were not the cause
    # Kills C-4 / failure 1c: surfaces the "one bad key poisons the whole
    # accessor" cost at review time instead of at first real flaky-integration
    # incident, and forces the diagnosis payload to actually name the culprit
    # distinctly enough to unblock a fix.
```

```gherkin
Feature: Absent-vs-unreadable classification is validated, not merely documented

  Scenario: A binding cannot distinguish "artifact lacks key" from "read failed"
    Given a binding whose underlying transport reports only success/failure at
      whole-artifact granularity (e.g. a single parse call with no per-key
      error information)
    When that binding is asked to classify one requested key as absent
    Then the binding MUST default to "unreadable" for every key it cannot
      independently verify was read successfully
    And this default MUST be enforced by a shared helper or validation check,
      not left to each binding's author to remember
    # Kills C-5: forces a mechanism (shared helper / validator check) instead
    # of a paragraph of guidance with no enforcement point.
```

```gherkin
Feature: Replay stability is proven on the normative requested-key-set path

  Scenario: Replay determinism using declared metadata, not artifact introspection
    Given a read accessor whose requested key set comes from validated
      definition metadata (not derived from the artifact's own keys)
    When the same model and fixture artifacts are resolved twice
    Then both runs MUST produce identical dispositions
    And this MUST be demonstrated on the metadata-driven path specifically,
      not only on a fallback path the RDR itself calls "fixture convenience,
      not the contract"
    # Kills C-6: closes the gap where A4's "Verified" evidence rides a code
    # path the RDR disowns as non-normative.
```

```gherkin
Feature: The silent-failure enumeration is derived, not just listed

  Scenario: Every refusal-adjacent branch is checked against the silent-failure inventory
    Given the full disposition table of accessor input classes and outcomes
    When each row is classified as "loud" (returns a value or names a refusal)
      or "silent" (could be mistaken for a different, better outcome)
    Then the Failure Modes section's enumerated silent-failure shapes MUST be
      derived by walking this table, not asserted from memory of prior review
      rounds
    And any newly-identified near-silent shape (e.g. read-back hitting an
      unreadable key with no prior branch at all) MUST be added with the same
      "mandatory guard" treatment as the other three
    # Kills C-7: turns the enumeration into a checkable process instead of a
    # count that happens to equal "however many review rounds have run so far."
```
