Model: claude-sonnet-5

# Findings Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | 0019:A6 | The safety property the whole design rests on ("store carries no key") has no implementation carrier — no exported cardinality probe exists anywhere on the accessor seam, and the RDR ships to Final review with this Pending | Implementation stalls, or ships a bespoke unreviewed API on `flowbind` that nobody 3-amigo'd, because the RDR authorized the semantics but not the plumbing | §technical-design, premortem |
| C-2 | 0019:A7 | The carrier-detection type-switch (`*flowbind.Writer` / `*flowbind.EditWriter` / `*cmdbind.Writer`) is asserted "mutually exclusive and exhaustive" but not yet verified against `Registry`'s actual branches; if a fourth arm or shared type exists, the verb silently seeds on an emptiness answer C1 declares UNDEFINED | Data corruption: an edit-carried or command-backed artifact gets written to when it should refuse, or a legitimate file-backed model refuses when it shouldn't | §technical-design, premortem |
| C-3 | 0019:C1 (CARRIER SCOPE) | The verb is scoped to file-backed writers only; edit-carried (0028) and command-backed (0025) flows keep the exact wall this RDR exists to remove | Half the flows in the fleet still require hand-typed `set-state` on first run; users of those flows get a confusing "carrier unsupported" refusal instead of the promised fix | §trade-offs (Negative), premortem |
| C-4 | 0019:C1 (Seeding is ALL-OR-NOTHING…) | The ALL-quantifier-over-STORE-keys predicate makes "operator clears keys one at a time to reset" silently reseed on the last clear, indistinguishable from intentional full reset | Operator who believes they escalated toward a full reset instead gets a resurrected state machine mid-workflow, potentially inside an automated pipeline | §trade-offs (Risk), Decision Rationale, premortem |
| C-5 | 0019:A4 | Extending a closed, Final peer's verb enumeration is justified by an analogy to `0002:C19`, but the actual blast radius (8 separate consumer sites: two hardcoded "four verbs" strings, a command-error message, docs generation, `llms.txt`, cli-output-contract taxonomy, `0005:D-naming`, README, model-authoring.md) is large and undischarged busywork disguised as "recordable override" | A shipped release where one of the 8 sites still says "four verbs" or omits `init-state` from `--help-all`, discovered by a user via a support ticket | §critical-assumptions (A4), premortem |
| C-6 | 0019:S9 / 0019:C1 (sealedKey) | The sealed-artifact test scenario requires constructing state via an exit-non-zero SETUP step whose non-zero exit is declared "expected... not a failure of the scenario" — a test-authoring foot-gun that will be miscoded as a real failure by whoever implements it, or silently skipped by a test runner that treats non-zero as fatal | CI flakiness or a silently-skipped regression test for the sealed-store boundary, discovered only when a real sealed artifact gets seeded into | §S9, premortem |
| C-7 | 0019:§failure-modes (Empty scalar) | The empty-scalar asymmetry (loader admits `note = ""` in `[initial]`, but `canonicalValue` refuses it unconditionally on the argv route) is explicitly routed out of scope to "RDR 0002 as a seed," yet C1's own SCALAR encoder arm (`members[0]`) will silently emit an empty string for it if the asymmetry is ever resolved toward admission — meaning this RDR's own normative text pre-commits the encoder's behavior without deciding it | A model author declares `note = ""` in `[initial]`, model loads clean, `init-state` either mysteriously refuses (today) or silently writes an empty string with no test coverage (post RDR-0002 resolution) — a spec gap masquerading as "not decided here" | §failure-modes, §Testing Strategy (S4 exclusion) |
| C-8 | 0019:§consequences (Negative, "two commands") | The two-command UX ("init, then work") for every fresh artifact is asserted as the accepted cost, but the RDR never asks how an operator DISCOVERS they need to run `init-state` on a brand-new artifact before their first real command — the only discovery path named (Activation Step 1) is that `flow next`'s `unknown[].reason: absent` report will mention it, which requires the operator to already know to interpret that report as "go run init-state" | New users hit `flow next`, get a wall of `unknown[].reason: absent`, and have no direct pointer telling them "run flow init-state" — same UX confusion the RDR claims to fix, just moved one command downstream | §Implementation Plan (Activation Step 1), premortem |
| C-9 | 0019:C1 (Torn/multi-writer no-repair) | A torn seed (role A committed, role B failed) is permanently unrepairable by re-running init-state — the only recovery is "explicit set-state of the listed keys, or discarding the artifact and re-running init" — but discarding a multi-role artifact means discarding role A's already-good, possibly externally-consumed state too | An operator who has a partially-seeded, multi-artifact flow either hand-repairs individual keys (re-introducing the exact transcription-drift failure mode this RDR exists to eliminate) or destroys good state to get a clean re-seed | §Failure Modes (Torn), §trade-offs, premortem |
| C-10 | 0019:§metadata (Status: Draft) + 0019:A6/A7 Pending | The RDR is being critiqued at what should be a pre-lock/Final gate, yet two Critical Assumptions remain Pending with "If wrong" consequences the RDR itself calls "approach-level, not editorial" — this is a structural readiness mismatch, not a nitpick | Finalization Gate's Assumption Verification is asked to "confirm no assumption... remains Pending... without a plan to verify before implementation begins" — the plan exists in name (booked as A6/A7) but the actual mechanism is undesigned | §critical-assumptions, §Finalization Gate |

# 1. The three most likely ways implementation goes wrong

## 1a. The emptiness probe (A6) gets bolted on as an unreviewed side-door, not a designed capability

**Root cause.** The entire safety argument of this RDR — the empty-store predicate that makes "cleared" survive "unconditional automated re-invocation" (A5) — depends on being able to answer one question: does the bound artifact carry any key at all? The RDR's own Technical Design section admits, in the clause labeled "EMPTINESS IS A NEW DATA FLOW," that no such capability exists on the shipped accessor seam: `internal/accessor/binding.go::ReadBinding.Read` is key-scoped and reports no cardinality, and `internal/cli/flowbind/flowbind.go::store` is package-private with no exported surface returning a key count. I confirmed this directly against source: `flowbind.go` exports exactly `Reader`, `Writer`, and `Gate` — no cardinality probe, no exported store accessor. A6 is marked **Pending**, and its own "If wrong" text says plainly: "the empty-store predicate cannot be implemented as specified, and the verb must fall back to a different gate... or the approach reopens... This is the predicate the whole safety argument rests on, so a refutation is approach-level, not editorial."

**Enabling passage.** `0019:A6` — Pending, Method "Source Search + Spike," with the fallback candidates listed as "a new `accessor.ReadBinding` capability, an exported `flowbind` cardinality probe, or a `read-state`-family surface" — three genuinely different API shapes, none chosen.

**Symptom.** Whoever implements Stage 8 hits this exact gap on day one, discovers the RDR authorized *what* to compute but not *how*, and one of two things happens: (a) implementation stalls waiting for a follow-up RDR/spike that should have run before lock, burning the "Draft" cycle time this record already spent; or (b) under schedule pressure, someone adds a quick unexported peek into `flowbind`'s internals (e.g., a private helper reused across package boundaries via an internal-only import, or literally re-opening the file to `len(store)`), which is exactly the kind of undesigned surface expansion that 0004:C3 and the accessor abstraction exist to prevent. Either way, the RDR's stated intent ("no implementation may substitute... a per-key variant") is enforceable only by review vigilance, not by anything the RDR actually built.

## 1b. The carrier type-switch is wrong on day one because A7 was never actually checked

**Root cause.** C1's carrier-scope clause is entirely dependent on being able to reliably distinguish file-backed writers from edit-carried and command-backed ones, from package `internal/cli`, without calling the unexported `flowbind::commandBacked` discriminator. The RDR's answer is a type-switch on `Definition.Binding`'s dynamic type, admitting only `*flowbind.Writer`. A7 says this needs verification that "the three constructed binding types are mutually exclusive and exhaustive over the registry's branches" and is marked **Pending**. I confirmed against `internal/cli/flowbind/registry.go` that `commandBacked` is indeed unexported (line 121) and that `Registry` branches into exactly `NewEditWriter`, `cmdbind.Writer`, and the fallthrough `flowbind.Writer` — so on today's code the exhaustiveness claim likely holds. But "likely holds today" is not what A7 asks for: it asks whether a future *fourth* binding type, or two carriers sharing one Go type, can silently break the type-switch. The RDR itself says the failure is silent: "the type-switch admits or refuses the wrong carrier silently — seeding on an emptiness answer C1 declares UNDEFINED."

**Enabling passage.** `0019:A7`, Pending, "If wrong: ...seeding on an emptiness answer C1 declares UNDEFINED, which is the failure S10 exists to catch." S10 catches the *known* carriers today; it structurally cannot catch a carrier that doesn't exist yet at RDR-lock time.

**Symptom.** Months later, someone adds a fourth accessor carrier (a plausible future extension given the project already has three) without revisiting `init-state`'s type-switch, because nothing forces that connection — there is no compile-time exhaustiveness check named anywhere in C1 (Go has no sealed-interface enforcement without an explicit unexported marker method, which the RDR doesn't ask for). The new carrier silently falls into the `default` refusal-or-not branch depending on how the switch is written, and either wrongly refuses a working carrier or — worse — is admitted by accident and gets seeded into with an undefined emptiness answer.

## 1c. The "no runtime consumer" analysis undercounts what breaks when `init-state` ships, because Activation Step 1's discovery story is circular

**Root cause.** The Consequences section names the true UX cost honestly: first run becomes two commands, not one. But the *discovery* path for the second command is `flow next`'s existing `unknown[].reason: absent` report — the same report operators already see today and already find opaque enough that this RDR exists to fix the underlying wall. Nothing in C1, the payload contract, or Activation Step 1 requires `flow next` (or `resolve`) to change its message to say "run `flow init-state`" when the absent keys are exactly the artifact's `[initial]` set on an otherwise-empty store. The RDR explicitly rules out automatic invocation (BR4, BR5) on defensible contract grounds, but then treats "name the verb in reference docs" as sufficient replacement for an actionable, in-band pointer.

**Enabling passage.** `0019:§implementation-plan` (Activation Step 1): "it names the verb from the surface where an operator meets the wall — `flow next`'s `unknown[].reason: absent` report" — but this is a doc-writing instruction, not a normative contract on `flow next`'s payload. C1 contains no requirement that the `unknown[].reason` value itself change.

**Symptom.** A user runs `flow next` against a fresh artifact, gets a wall of `unknown[].reason: absent` entries exactly as they do today, and has no stronger signal than before that the fix is "go run `init-state`" unless they already read the docs. The RDR removes the *mechanical* wall (manual transcription) but leaves the *discoverability* wall (users don't know the mechanical wall exists) essentially untouched, and ships this as if it were resolved because reference docs were updated.

# 2. The section that will be rewritten within 6 weeks of shipping

**Load-Bearing Decisions / Selection-Predicate, and its restatement in C1 ("Seeding is ALL-OR-NOTHING").**

This is the section most likely to get rewritten, because it is the one place where the RDR's chosen safety property (empty-store, ALL-quantifier, STORE-key-count) collides hardest with an operator mental model that the RDR itself predicts and then dismisses as an accepted risk: "operators who clear keys individually to reset must discard the artifact or expect the reseed." That is not a hypothetical edge case — it is the single most natural way anyone manually resets state without knowing the artifact's internal representation. The first real incident report — someone cleared three tags to "start fresh," and a nightly CI job's `init-state` reseeded the fourth-cleared-to-zero artifact and silently un-did their reset in a way that surprised them — will force a redesign of the predicate itself (almost certainly toward the previously-rejected tombstone mechanism, or toward a `--force`-style explicit-intent flag `set-state`-side), because the RDR's own "Briefly Rejected" analysis already identifies tombstones as "the only mechanism that distinguishes cleared from never-set" and rejects them purely on wire-format blast radius, not on merit. When the blast radius argument loses to a real production incident, this section is the one that gets reopened, and C1's ALL-quantifier and store-emptiness language will need a second normative pass.

# 3. The assumption that will not survive first contact with a real user

**A5 / C1's premise that "unconditional automated re-invocation" is the threat model that matters, and that a human operator's *manual*, *intentional*, *incremental* reset is out of scope because it maps onto the same store-emptiness signal.**

The RDR spends enormous care proving that a scheduled job calling `init-state` every deploy will never resurrect a cleared key — that is the automation axis, and it is genuinely well-defended (A5's spike, RT4, MVV steps 6–9, S3). But the actual first user story that will break this design is not automation at all: it's a human who wants to reset one flow instance back to its initial state and reaches for the tool that's *named* `init-state`. They will clear the owned keys (the only exposed mechanism, since there's no `reset` verb), watch the store empty out, and get an automatic, correct-per-spec, but completely undocumented-in-their-mental-model reseed on the next `init-state` call — except the RDR's own risk-and-mitigation entry admits this is *expected behavior*, not a bug. The gap is that "expected behavior, formally correct" and "matches what a human means by reset" are different claims, and the RDR only defends the first. The disclosed residual is honest, but it is disclosed to reviewers of a 1600-line record, not to the operator typing `flow set-state --clear`. First contact will be: an operator "resets" a flow by clearing keys, something else automated (or a colleague) invokes `init-state`, and the flow silently jumps back to `[initial]` instead of staying at whatever partial-reset state the operator expected. This is not a corner case; it's the most intuitive interaction with the exact verb this RDR ships.

# 4. The premortem, written as if the failure already happened

Six weeks after `flow init-state` shipped, `intrastate` filed its own incident against `models/rdr.toml` — the very gate model the RDR cites as proof the wall was real and the fix mattered.

A maintainer working the RDR pipeline itself wanted to back out a bad `gate_passed` write during a botched Stage-7 run. Following the only documented undo mechanism (`flow set-state --clear gate_passed`), they cleared `gate_passed`, then, mid-investigation, cleared `status` too, intending to leave `stage` alone as an anchor while they figured out what had actually gone wrong. A CI job — the project's own nightly RDR-index rebuild, which runs `flow init-state` unconditionally as a defensive "make sure the artifact is seedable" step exactly per the automation pattern this RDR was built to make safe — ran between the two clears and the third. It saw a non-empty store (one key, `stage`, remained) and correctly, per C1, did nothing. Reassured that automation "couldn't touch it," the maintainer then cleared `stage` as well, emptying the store completely, intending to hand-reconstruct all three values from the incident log afterward.

The next scheduled run of the same CI job — unchanged, still unconditional, still doing exactly what A5 and RT4 certify as safe — saw a zero-key store and reseeded all three tags to their `[initial]` values: `stage = "draft"`, `status = "proposed"`, `gate_passed = false`. This silently overwrote the maintainer's in-progress incident reconstruction with the model's static bootstrap values, with a `flow init-state` payload that reported "seeded-all" — a correct, no-op-free, fully-compliant success exit — with nothing in the CLI output flagging that this reseed had just clobbered manual recovery-in-progress. The maintainer discovered the loss only when `flow next` produced garbage transitions inconsistent with their notes, twenty minutes into re-deriving what they'd already reconstructed.

Root cause, traced to the record: `0019:C1`'s "Seeding is ALL-OR-NOTHING... the count is over STORE keys, not owned keys" combined with the disclosed-but-unenforced residual "clearing the LAST remaining key empties the store... a subsequent init-state therefore RESEEDS such a store" (also `0019:C1`). The RDR's own Risk-and-Mitigations entry for this exact scenario says only: "the no-op payload and the docs must state the boundary in store terms... MVV step 9 pins it so it cannot move silently." Documentation and a regression test are not a mitigation against an operator's incremental clear sequence interacting with an unrelated automated job's polling cadence — they are a mitigation against the design *changing*, not against the design's *documented behavior being exactly the trap*.

The specific functions implicated: `internal/cli/flowbind/flowbind.go::store.Apply` (clear-as-delete, no tombstone — cited directly in C1 and A5), the not-yet-written `init-state` verb's predicate check (over STORE keys per C1's Load-Bearing Decisions), and `internal/cli/flow_state.go::runFlowSetState`'s clear path, which the RDR never asks to warn "this is the last key; the store will now read as unseeded to `init-state`."

# 5. Acceptance tests that would have caught each failure at RDR-review time

**AT-1 (catches 1a — emptiness probe undesigned).**
```
Given the RDR is proposed for Final status
When the reviewer checks Critical Assumptions
Then no assumption whose "If wrong" text is marked "approach-level, not editorial"
  may remain in status Pending
And the RDR must name the exact chosen carrier for the emptiness read
  (not a list of three candidates) before lock
```

**AT-2 (catches 1b — carrier exhaustiveness unverified).**
```
Given internal/cli/flowbind/registry.go::Registry's write-construction branches
When a reviewer enumerates every branch that can be reached
Then the RDR's normative text must cite the exhaustiveness check by test name,
  not by narrative assertion
And a fourth accessor carrier, if added later, must be required (by an explicit
  cross-reference in either RDR or code comment) to update this verb's type-switch
  before it can be considered complete
```

**AT-3 (catches 1c and C-8 — discovery gap).**
```
Given a fresh, empty, file-backed artifact bound to a state-machine model
When an operator runs `flow next` without first running `flow init-state`
Then the `unknown[].reason: absent` payload for `[initial]`-declared keys
  must name `flow init-state` directly, not merely be documented externally
```

**AT-4 (catches the premortem and C-4 — incremental-clear resurrect surprise).**
```
Given a seeded artifact with three owned keys
When an operator clears the keys one at a time in three separate invocations,
  with an unrelated automated `init-state` call interleaved after the second clear
Then the third clear (which empties the store) must not result in a silent,
  unflagged reseed on the next automated `init-state` call
Or, if the current design is retained, the RDR must require `flow set-state --clear`
  on the LAST remaining key to emit an explicit warning that the artifact has
  returned to the initializable class
```
This is the test the RDR does not have. MVV step 9 tests the boundary in isolation ("clear the REMAINING key... init-state... seeds again") but never tests it as the *last op in a human reset sequence interleaved with automation* — which is the actual failure shape.

**AT-5 (catches C-9 — torn-seed recovery destroys good state).**
```
Given a two-writer artifact where role A seeded successfully and role B failed (torn)
When the operator wants to recover
Then the RDR must specify a repair path that does not require discarding role A's
  already-committed, possibly externally-consumed state
And "discarding the artifact and re-running init" must not be the only documented
  recovery option for a multi-role torn seed
```

**AT-6 (catches C-6 — sealed-artifact test fragility).**
```
Given the S9 test fixture construction, which requires a SETUP step to exit non-zero
When a reviewer reads the test as written
Then the test framework must explicitly assert on the persisted artifact bytes
  after the non-zero-exit setup step, not rely on a comment explaining the exit
  code is not a failure
And CI must not be configured to treat this step's non-zero exit as a test failure
  by default
```
