Model: claude-sonnet-5

# Premortem critique: RDR 0024 — Declared emit vocabulary

## Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `0024:A5`, `0024:C1` (`kind = "enum"` closed domain) | The enum mechanism has zero headroom for any interpolated or per-invocation emit value; `0010:A6`'s deferred "table interpolates a value" widening is not merely future work, it is a value a real author hits the day they add ONE parameterized route (`Next: /rdr-prelock {id} critique`) to an otherwise-enumerable key. That key cannot be `enum`, only `scalar` (unvalidated) — so the whole declaration for that key silently loses the proof the RDR sells, on the exact key most worth proving. | A model author declares `next` as `enum` for 21 of 22 values, hits one route needing an id, and either (a) can't declare `next` at all and gives up on the whole key, or (b) declares it `scalar` and the one column the kata filed against (`next`) reverts to unvalidated strings — the RDR's own MVV consumer model (`rdr-status.toml`, cited in A5) is one interpolated value away from this. | §1, premortem |
| C-2 | `0024:C2` fail-fast, "one refusal per run" + Trade-offs "N-round fix-and-rerun loop" | First adoption on any non-trivial model (the RDR's own example: 54 `[rule.emit]` blocks) is an N-round manual loop with NO tooling to shortcut it — no `--report-all`, no batch mode, nothing. The RDR names this cost explicitly and dismisses it as "a cost of the tier this RDR lands in" rather than owning a mitigation. | An author who declares `next` on `rdr-status.toml` runs lint, fixes one typo, reruns, fixes the next, 20+ times, with the pipeline order unspecified so they cannot even predict which defect surfaces next. They give up partway and either abandon adoption or blanket-declare `scalar` to make it stop — which is exactly the "scalar hollowing" failure mode the RDR names as a risk (premortem P-2/P-11) and declines to mitigate beyond "reviewers will notice." | §1, premortem, Trade-offs |
| C-3 | `0024:C1` "domain is `any`-typed" + hand-written refusal arms; `0024:A7` (Pending — closure unverified) | The two-shaped `domain` field is the one place strict decoding is deliberately weakened, and C1's list of hand-written refusal arms is derived "from the type structure rather than enumerated exhaustively" (A7's own words). A7 is Pending at Draft with the closure question explicitly unresolved — whether five hand-written arms are the COMPLETE set of malformed shapes. If they are not, a malformed domain shape loads silently, and the declaration then *overclaims* what it proves — worse than no declaration, because an author reasonably trusts a present `[emit]` table. | A model author writes a subtly malformed disposition sub-table (a shape not in the five enumerated arms — e.g. a domain value that is an array of arrays, or a disposition key holding a table instead of an array), lint passes, and the author believes the key is proven when it silently isn't. | §1, §3, premortem |
| C-4 | `0024:§Trade-offs` "Negative: opt-in means an undeclared model keeps today's silent gap" + A3 | The RDR ships zero coverage on merge day by its own admission (A3: zero declarations exist repo-wide) and schedules the motivating consumer's actual adoption as "out of scope here and unscheduled by this record." The seed kata that justified the whole RDR (`vt9n`) is not closed by this RDR shipping — only by a SEPARATE, unscheduled adoption effort in a different repo. | Six weeks after this RDR ships, `intrastate lint` on the real consumer models still exits 0 on the exact adversarial two-rule table the Problem Statement describes, because nobody has gone back and declared `[emit.next]` on `rdr-status.toml`. The defect this RDR was written to close is still open in production. | premortem, §2 |
| C-5 | `0024:A8` (Pending), `TestReq146` scope | The determinism sweep A8 leans on (`TestReq146_EveryEmittedSequenceIsASortedSlice`) is confirmed (by reading the test) to be a HAND-MAINTAINED field list over `table.Row` (`Atoms, NextTags, Writes, RequiresOwned, Gate, Escape, Suffix`) — it does not reflectively cover new fields, and `Model.EmitDecls` is a `Model`-level field, not even a `Row` field, so the sweep cannot reach it by construction regardless of what Phase 2 does. A8 is honest about this being unverified, but the RDR still cites the sweep in C3's normative text ("Sorting is what makes the carry deterministic... the property `TestReq146`... already pins") as if it were load-bearing corroboration, when grounding shows it structurally cannot be. | An implementer reads C3, sees `TestReq146` cited as existing proof of the determinism property, and either skips writing scenario 4's dedicated assertion (treating it as redundant) or discovers at Phase 2 that the "proof" was rhetorical, costing a review round. | §1, §5 |
| C-6 | `0024:C4` "the dispositions field is registered under JDR 0002 §D1" + `0024:S7` | C4 and scenario 7 both depend on a reflective oracle from JDR §D1 that "does not exist at HEAD" because 0023 — the RDR that would ship it — is itself Final-but-unimplemented. The Done criterion for scenario 7 is explicitly "unsatisfiable while it is unmet," and landing order between 0023 and 0024 is undetermined at RDR-lock time. This is two Draft/Final-but-not-built RDRs each assuming the other's scaffolding will exist when they land. | If 0024 lands first (plausible — it is Draft, 0023 status is unclear from this text but described as unimplemented), Phase 3's `dispositions` field ships with NO reflective registration check at all, silently violating the very obligation C4 calls "not optional," discovered only later when 0023 finally implements and finds an unregistered field it must now retrofit against a shipped payload. | §1, premortem |
| C-7 | `0024:D-naming` (payload field is `dispositions`) + `0024:C4` | `dispositions` is plural but keyed per-emit-key, one token each — a consumer must learn that `dispositions["next"]` is a single string, not a list, despite the field's own name suggesting a set. This is a naming footgun the RDR resolves by fiat ("rejected: `emit_dispositions`... `kinds`") without weighing the singular/plural mismatch at all — the rejected-alternatives list only compares sibling names, never interrogates the chosen one's shape-signaling. | A new consumer author, skimming `docs/cli-output-contract.md`, writes `dispositions["next"][0]` or iterates it as an array before discovering (from a runtime type error or from reading the source) that each value is a bare string. | §3 |
| C-8 | `0024:§Trade-offs` Risk "domain drift" and "disposition-token drift", both Mitigation: "out of intrastate's scope by design... the consumer's cross-file seam test... closes the loop" | Both of the two riskiest failure classes the premortem itself names (P-1 stale-but-valid domain, P-12 disposition-token rename) are mitigated by a consumer-side test that does NOT EXIST, is not scheduled by this RDR, is not tracked by any kata cited here, and lives in a "sibling engine repo" this RDR admits it cannot reach. The mitigation is aspirational, not built — the RDR cites its own absence as its safety net. | Six months later, `rdr-status.toml` renames `stop` to `halt` in one key's disposition table (a one-line, lint-clean edit per C1's "author-owned and unversioned" clause), every downstream skill matching `dispositions["next"] == "stop"` goes quietly false, and nothing in the toolchain — because nothing was built — catches it. This is a stop being read as a route again, the EXACT original kata defect, now laundered through a declaration that made everyone confident it couldn't happen. | premortem, §3, Risks |
| C-9 | `0024:A2` "28 inline whole-payload assertion sites... every one of which pins the `emit`/`next` adjacency" + Phase 3 licensed diff | The 28-site mechanical golden regeneration is real and enumerated, which is good — but it is also the RDR's own admission that this repo's test suite has near-zero payload-shape abstraction: any future append to `resolvePayload` (there will be more — `0023`'s `--plan-only` echo group is coming next per JC1) repeats this exact 28-site brittle-golden tax. The RDR treats this as someone else's problem ("0011's own adversarial oracle... deliberately brittle by its authors' comment") rather than flagging that the test architecture itself is what will need rewriting, not just this one field's fixtures. | The team ships 0024's 28-site diff cleanly, then hits an even larger version of the same diff when 0023 lands its own payload field, and a third when the next payload field after that ships — the fixture tax compounds per RDR rather than the suite ever amortizing it. | §2 |

## 1. The three most likely ways implementation goes wrong

### Failure 1 — The enum mechanism cannot survive contact with one interpolated value (C-1)

**Root cause in the RDR.** C1 defines exactly four kinds — `enum | bool | int | scalar` — and gives `enum` a closed, hand-authored domain. A5 (Verified, Peer RDR + empirical sweep) establishes that the motivating consumer's routing keys are "enumerable at authoring time" and backs this with real numbers: `rdr-status.toml`'s `next` key has 22 distinct fixed-literal values, `rdr-write.toml`'s `op` key has 10. Both A5 and A6 (`0010:A6`, cited directly) are careful to say the *current* motivating consumer's answers are closed — but `0010:A6` also states, in the same breath, what it explicitly defers: "A consumer needing the table itself to interpolate — a template language in an emit value — is the widening A6 defers." That deferred widening is not a hypothetical; it is one route away from every real routing table the RDR cites as its motivating evidence. `rdr-status.toml` itself already renders record numbers into commands (`Next: /rdr-prelock 0046 critique`) — per A6's own analysis, that particular composition currently stays outside the table (the number is caller-supplied, not table-authored), but that boundary is a discipline the model author maintains by hand, not something C1 enforces or even detects. Nothing in C1 stops an author from later writing one templated value into an otherwise-clean enum key, and nothing in the RDR helps them notice they've crossed the line until the declaration for that whole key becomes either impossible or must retreat to `scalar`.

**Specific passage.** `0024:A5`: "The motivating consumer's emit answers are enumerable at authoring time... so an enum domain with a disposition partition can actually be authored for each of its routing keys." This is true today, stated as a durable property, when it is actually a snapshot. `0024:C1`'s kind list has no accommodation for "closed except for one caller-supplied fragment" — the RDR's only escape is the whole-key `scalar` fallback, which throws away validation for the entire key, not just the parameterized member.

**Symptom the user sees.** A model author who wants to close the exact gap this RDR was built for (`next` on `rdr-status.toml`) discovers mid-authoring that one route needs a parameter, and is forced to choose between (a) declaring nothing for that key — reverting the whole key to today's unproven state, defeating the RDR's purpose for the highest-value key in the corpus — or (b) declaring it `scalar`, which the RDR's own Risks section names as "scalar hollowing" and admits is undetectable except by manual review. Either way, the flagship consumer scenario this RDR was designed to fix does not get fixed for its highest-traffic key.

### Failure 2 — Fail-fast, one-refusal-per-run makes first adoption on a real model an unbounded manual grind (C-2)

**Root cause in the RDR.** C2 inherits fail-fast, one-error-per-`Load()`-call cardinality from the load tier wholesale ("Reporting is FAIL-FAST, one refusal per run — inherited from the load tier, not decided here"), and the Trade-offs section is unusually candid about the cost: "first adoption on a large model is an N-round fix-and-rerun loop... where N is the number of defects the declaration exposes, not the number of runs an author expects." On the RDR's own example model (54 `[rule.emit]` blocks, `rdr-status.toml`), N could plausibly be dozens on a single first pass — every misdeclared kind, every out-of-domain value, every unknown key surfaces one at a time, in an order the RDR explicitly declines to make deterministic ("the order in which independent defects are checked is deliberately unspecified"). This is a design decision this RDR is entitled to inherit, but it does not follow that the RDR should ship with NO tooling story for the one adoption workflow it exists to enable. The RDR's own MVV walks exactly two defects across two runs and calls that sufficient evidence; it never demonstrates or estimates the workflow at N=20+.

**Specific passage.** `0024:C2`: "The choice here rests on the load tier's own contract and the missing comparator above... Multi-defect load reporting would require inventing the total order `load.go` declines to fix; it is a pipeline-wide change out of scope here, left to a later RDR that owns it." This defers the actual adoption-usability problem to an unscheduled future RDR while shipping the feature whose value proposition depends on exactly that adoption happening.

**Symptom the user sees.** An author trying to declare `rdr-status.toml`'s `next` key for the first time runs `intrastate lint`, gets one finding, fixes it, reruns, gets a different finding in an unpredictable order, and repeats this 15-30 times before the model is clean. This is precisely the kind of tedium that produces the RDR's own predicted failure mode: blanket `scalar` declarations to make the loop stop, which the RDR calls out by name (premortem P-2/P-11, "scalar hollowing") without shipping a mitigation stronger than "reviewers will notice a `scalar` kind in the declaration table."

### Failure 3 — Two risk mitigations that don't exist are load-bearing for the two worst failure modes the RDR itself names (C-8)

**Root cause in the RDR.** The Risks and Mitigations section names domain drift (a declared value goes stale relative to the real command surface) and disposition-token drift (a consumer's string-comparison test silently breaks when a model author renames a disposition token) as risks, and both are mitigated identically: "out of intrastate's scope by design (generic core); the consumer's cross-file seam test... closes the loop." That seam test is not built, not scheduled, not tracked by a kata in this record, and explicitly said to live in "a sibling engine repo... not reachable from a phase here." The RDR is aware this is aspirational — Phase 4 even declines to add a worked disposition-partition example inside this repo for the same reason (A5's motivating models "live in the sibling engine repo and are not reachable from a phase here") — but it still treats the unbuilt seam test as the closing move for two premortem-flagged risks (P-1, P-12) rather than as a genuine open gap.

**Specific passage.** `0024:§Risks and Mitigations`, disposition-token drift: "a model renames `stop` to `halt` and every consumer's `dispositions[k] == "stop"` test goes quietly false (premortem P-12). **Mitigation**: same seam as domain drift... one authoritative list to pin is what the declaration provides." Providing a list to pin against is not a mitigation for silent drift unless something actually pins against it — and nothing in this RDR, this repo, or any cited kata does.

**Symptom the user sees.** A model author renames a disposition token in a one-line, fully lint-clean edit (C1: "author-owned and unversioned... a domain may be widened or narrowed by editing the model, and intrastate holds no history to check the edit against"). Every downstream skill comparing against the old token string goes silently false. This is functionally identical to the ORIGINAL seed defect this RDR exists to prevent — a stop read as a route — except now it happens through the new mechanism, one layer removed, with the team's guard down because "we declared it, so it's proven."

## 2. The one section that will be rewritten within 6 weeks of shipping

**Trade-offs → Consequences, the "Negative: opt-in means an undeclared model keeps today's silent gap" bullet, together with the Phase 4 scope boundary that keeps the motivating consumer's adoption "out of scope here and unscheduled by this record."**

This RDR is entirely honest that it ships *capability*, not *coverage*: A3 proves zero declarations exist anywhere in the corpus at merge time, and the only Phase 4 artifact that exercises the disposition-partition grammar is a brand-new, purpose-built toy example the RDR invents specifically because the real motivating models (`rdr-status.toml`, `rdr-write.toml`) live in a sibling repo this RDR cannot touch. That means the actual kata this RDR was filed against — `vt9n`, the two-rule adversarial table lints clean — remains literally unfixed in the real consumer's models the day this RDR ships. Within six weeks, someone will notice the seed kata is still open in the engine repo, ask why a "Final" fix for it shipped with zero real-world coverage, and the Trade-offs section's blithe "adoption... is real work in a different repo... out of scope here" will get rewritten into either (a) a follow-up RDR that actually schedules the sibling-repo adoption, or (b) a walked-back claim about what this RDR achieves. The gap between "the mechanism exists" and "the kata is fixed" is the single largest credibility risk in the document, and the RDR resolves it by declaring the gap out of scope rather than closing it or scheduling its closure.

## 3. The one assumption that will not survive first contact with a real user

**A5 — "the motivating consumer's emit answers are enumerable at authoring time... so an enum domain... can actually be authored for each of its routing keys."**

A5 is marked Verified, backed by real numbers pulled from the live consumer models, and it is honestly the strongest-evidenced assumption in the document — that is exactly why its failure mode is dangerous rather than obviously flagged. It is a snapshot claim dressed as a durable one. The moment a real author needs one parameterized emit value on an otherwise-closed key (a near-certainty as the model grows — `0010:A6`, which A5 cites approvingly, explicitly names table-side interpolation as deferred future work, not ruled out), A5's "can actually be authored for each of its routing keys" stops being true for that key, and C1 offers no partial escape: the choice is binary, `enum` (all-or-nothing closed) or `scalar` (all-or-nothing open), per key. A5's evidence is real; A5's implicit promise — that this shape scales to the consumer's evolving needs — is not tested anywhere in the record, and the record itself (via `0010:A6`) names the exact scenario that breaks it.

## 4. Premortem

It is six weeks after RDR 0024 shipped. The load pipeline change landed cleanly, all 28 licensed golden-fixture sites in `internal/cli/flow_demand_0011_test.go` were regenerated exactly as A2 enumerated, `TestReq39_EmitSitsImmediatelyAfterGatesOnTheWire` was updated to assert 14 wire keys and `NumField() == 15`, and `models/examples/pricing-decision-table.toml` lints clean with `plan` and `dpa` declared. Every scenario in the Testing Strategy is green. By its own acceptance bar, the RDR delivered.

And the kata it was filed against — `intrastate#vt9n`, the two-rule table where a misspelled command value and a typo'd `nxet` key both lint exit 0 — is still open, because nobody has gone into the sibling engine repo and actually declared `[emit.next]` on `rdr-status.toml`. A3 predicted this precisely: at merge, the only declared models anywhere are the two Phase 4 toy examples this RDR invented for itself. The declared-emit-vocabulary feature exists; it validates nothing that matters yet.

Three weeks later, an engineer picks up the adoption work as an unplanned side task. They start declaring `next` on `rdr-status.toml` — 54 `[rule.emit]` blocks, 22 distinct `next` values across 14 route commands, six `stopped:` tokens, plus `none` and `resolve:lens`. Twelve `lint`-fix-`lint` cycles in, at defect #9, they hit a route that needs a record id folded into the command text — a case A5 didn't anticipate because at Draft-time no such route existed. They cannot declare `next` as `enum` without either omitting that one route from the domain (which makes it refuse under `unknown_emit_key` from every rule that legitimately emits it) or falling back to `scalar` for the whole key, which throws away the proof for the 21 other, perfectly enumerable values. Under deadline pressure they choose `scalar`. The declaration table now says `next` is validated. It is not. A `dispositions["next"]` reader downstream trusts the declared-key list as its contract; nothing tells them one of the 22 branches is unchecked.

Two months after that, a model author renames the `stop` disposition token to `halt` in a routine cleanup pass — a one-line, fully lint-clean edit, exactly as C1's "author-owned and unversioned" clause permits, with "intrastate holds no history to check the edit against." Every skill downstream comparing `dispositions["next"] == "stop"` goes quietly false. Nothing catches it: the consumer-side seam test the Risks section names as the mitigation for exactly this ("every declared emit domain member is a real command or stop token") was never built — it was always described as living in a sibling repo, out of this RDR's reach, unscheduled by any kata this record cites. A stop gets read as a route again. It is the original kata defect, recurring through the mechanism built to prevent it, one abstraction layer removed, at a moment when the team's guard is down because "we declared it, it's proven."

Meanwhile 0023 — the sibling RDR whose JDR §D1 reflective oracle C4 and Testing Strategy scenario 7 both depend on for registering the new `dispositions` field — never landed first. 0024's `dispositions` field shipped with the registration obligation unmet, "not discovered as a red reflective test" because there was no reflective test to fail against. When 0023 finally implements months later, its author finds an unregistered payload field already in production and must retrofit the assignment against a shipped, versioned wire contract instead of gating it at 0023's own build time as the RDR assumed would happen.

The retrospective conclusion: the grammar, load proof, and envelope join all work exactly as specified — the RDR's internal correctness claims hold up completely. What failed is that the RDR treated "the mechanism is sound" and "the kata is fixed" as the same milestone, when the second one required unscheduled, unfunded, cross-repo work the RDR explicitly declined to own.

## 5. Acceptance tests that would have caught each failure at RDR-review time

```gherkin
Feature: RDR 0024 review-time acceptance checks

  # Catches C-1 / Failure 1
  Scenario: A closed enum key needs one interpolated value
    Given a model author has declared an emit key as `kind = "enum"`
      with N-1 fixed-literal members already covering every current route
    When the author needs to add one route whose value embeds a
      caller-supplied fragment (a record id, a path)
    Then the RDR's grammar provides a way to keep the other N-1 members
      validated while marking that one member as open
    And the answer is not "redeclare the whole key as scalar"

  # Catches C-2 / Failure 2
  Scenario: First declaration on a 50+-block model
    Given a model with 54 `[rule.emit]` blocks and no existing declarations
    When an author declares one key for the first time
    Then the RDR names a bounded number of lint-fix-rerun cycles,
      or a batch/report-all path, required to reach green
    And the number is not "N, where N is the number of defects, discovered
      one per run, in an order the tier declines to fix"

  # Catches C-3 / A7 closure gap
  Scenario: A7's closure claim is proven before lock, not carried Pending
    Given C1 lists five hand-written refusal arms for the two-shaped
      `domain` field
    When the Finalization Gate's Assumption Verification step runs
    Then A7 is Verified, not Pending, with a negative sweep proving no
      sixth malformed shape reaches the pipeline unrefused
    And "verify at implementation" is rejected as a lock condition for
      a claim the grammar's soundness depends on

  # Catches C-4 / premortem's central failure
  Scenario: The seed kata is actually closed, not just closable
    Given kata `intrastate#vt9n` is the reason this RDR exists
    When this RDR is marked Final
    Then either the motivating consumer's model in the sibling engine
      repo has `[emit.next]` declared and its lint is green,
      or a scheduled, tracked follow-up (kata, not prose) closes that
      gap within a stated window
    And "adoption... is out of scope here and unscheduled by this
      record" is rejected as a closing disposition for the RDR's own
      motivating defect

  # Catches C-5 / A8's misattributed corroboration
  Scenario: TestReq146 is checked for actual scope before being cited
    Given C3 cites `TestReq146_EveryEmittedSequenceIsASortedSlice` as
      already pinning the new carrier's determinism
    When a reviewer reads the test body
    Then the review confirms the sweep enumerates `Model.EmitDecls`
      (or `EmitDecl.Domain`) by name or reflectively
    And if it does not, C3 does not cite it as corroboration — it
      states plainly that no existing test reaches the new carrier and
      Phase 2 owes a fresh assertion

  # Catches C-6 / cross-RDR sequencing risk
  Scenario: The JDR §D1 registration obligation is enforceable regardless
    of landing order
    Given C4 depends on an oracle that "does not exist at HEAD" and
      whose existence depends on 0023 landing first
    When 0024 is implemented before 0023
    Then Phase 3 includes its own local assertion that `dispositions`
      is a known, accounted-for payload field — not merely a note that
      the obligation "transfers" to 0023's future implementation
    And Done is not satisfiable by an obligation with no owner until a
      second, independently-timed RDR lands

  # Catches C-8 / unbuilt mitigations
  Scenario: A cited mitigation must exist or be scheduled, not just named
    Given the Risks section mitigates domain drift and disposition-token
      drift with "the consumer's cross-file seam test"
    When the Finalization Gate's Scope Verification step runs
    Then that seam test is either already built, or tracked by a kata
      with an owner and a repo
    And "out of intrastate's scope by design" is rejected as a
      mitigation for a risk the same RDR's premortem rated severe
      enough to name explicitly (P-1, P-12)
```
