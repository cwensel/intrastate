Model: claude-sonnet-5

# Critique — RDR 0011 (iter-2, pass B)

## Findings Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `0011:A15` | The reader-demand-set extension is the mechanism that makes the whole predicate change SAFE (without it, every model with a match-only owned key silently degrades to `--all` with no error), yet A15's own Status is `Pending`, not `Verified`, and its Method is "Source Search + MVV Test" — deferred verification bolted to the same implementation pass that ships the predicate change it protects. | A model author adds a match-only owned key; `flow next` silently reports every row as a candidate with `{key, absent}` in every entry, and nothing in the build fails until someone notices the list never narrows. | §1, premortem |
| C-2 | `0011:C1` (lines 112–210), `0011:C2` (214–270), `0011:C3` (274–375) | Three normative contracts, each running 60–100 lines of prose with nested clauses, all serve ONE behavior change ("stop stripping match before the kernel probe") plus one incidental, unrelated payload rename (`unresolved`→`unknown`). The complexity is self-inflicted: two independent decisions (predicate change; payload shape change) are locked into one RDR and one implementation window. | Implementers cut corners under review pressure — the RDR's own Testing Strategy S6 admits "a green 0005 suite is NOT evidence C1 shipped," meaning the safety net for the *predicate* change is not the test suite people will actually run day to day; a partial, technically-passing implementation ships. | §1, §2 |
| C-3 | `0011:§approach` / `0011:C1` lines 160–178, "A match atom whose key is ABSENT … MUST be omitted from the probe" | The correctness of the whole design leans on one CLI-side presence test (`assembledView`) staying in lockstep with the kernel's own presence semantics (A11) and never diverging as the codebase evolves. There is no shared code path enforcing that agreement going forward — it is proven once, in the RDR, over the CURRENT source, and never re-verified by any test that would fail if a future change (e.g., a new tag provenance, a caching layer between `runReaders` and `assembledView`) broke the agreement. | Six months later a refactor to `flow_exec.go::runReaders` (e.g., adding a cache or a lazy-read optimization) silently reintroduces a presence-disagreement bug; `flow next` starts reporting or omitting `unknown` entries incorrectly, and no test in the suite is positioned to catch it because A11's invariant was checked by prose reasoning, not by an executable oracle. | premortem |
| C-4 | `0011:§load-bearing-decisions`, "Undecided reporting shape" — "One asymmetry against the cited precedents is accepted: … the widened rows carry strictly LESS [information]" | `--all` is supposed to be the "escape hatch" / diagnostic mode a caller reaches for when the default's narrowed list confuses them, but C2 explicitly strips match facts out of `unknown` under `--all`, so `--all` cannot tell the caller WHY a row was excluded by match — only THAT it is now included. The RDR's own Failure Modes section admits the remedy is manual eye-comparison against the model TOML file. | A skill author debugging "why isn't row X a candidate" runs `--all`, sees the row, and still cannot programmatically determine which match atom excluded it — must open the model source file and read `[rule.match]` by hand, which is not something an automated skill/agent caller can do at all. | §3, premortem |
| C-5 | `0011:§metadata` Overrides clause; `0011:C3` lines 292–338 | The `unresolved`→`unknown` rename is described as "mechanical" but is a breaking wire-format change for every out-of-repo consumer, bundled into an RDR whose headline feature is the match predicate. A consumer who only cares about "give me the match-conditioned list" is forced to also handle the field rename and the new `{key,reason}` object shape, in the same release, with no independent migration path or flag to decouple the two changes. | An external skill/tool parsing `candidates[].unresolved` as `[]string` breaks in one upgrade for a reason (the payload shape) that has nothing to do with the reason it opted into (narrower candidates), with no way to get one without the other. | §2, premortem |
| C-6 | `0011:C1` lines 127–139, "never by rebuilding match tags from `row.Atoms` in `internal/cli`" | The contract forbids re-deriving match tags from `row.Atoms` and mandates filtering the CONVERTED probe (`row.KernelRow()` then drop by `Key`) specifically to avoid re-implementing `seamValue`'s set-canonicalization. This is exactly the kind of subtle, easy-to-violate implementation constraint that a future maintainer "fixing" or refactoring `excluded()` will not know about — the reasoning lives only in RDR prose, not in a guard comment or a test that fails if someone reverts to filtering `row.Atoms` directly. | A future PR "simplifies" the probe-building code by filtering atoms before `KernelRow()` for readability; a set-valued match key silently mis-compares against the kernel (wrong JSON encoding), and `flow next` starts excluding or including rows incorrectly for set-kinded match atoms with no test catching it (the RDR records no test that would distinguish the two implementations for a set-kinded match key). | §1, premortem |
| C-7 | `0011:§technical-design`, "This is the one signature change in the file" (`excluded` returning `Result` instead of `bool`) | `excluded()` changes from returning `bool` to returning the kernel's `Result`, and `runFlowNext`'s loop (`if excluded(row, owned, req.observed) { continue }`) must be rewritten to consult the `Result`'s refusal kind, while `summarize` must additionally read `Refusal.Undecided` off it. This is a bigger, more error-prone refactor than the RDR's framing ("the one signature change") suggests — two call sites, a changed control-flow shape (bool test → refusal-kind switch → guard-payload extraction), and the ordering constraint (owned_state_unavailable precedes the undecided-guard collection, per A13) all have to be gotten right simultaneously in one function rewrite. | An off-by-one in the refactor (e.g., checking `result.Refused()` and excluding on ANY refusal instead of only `KindNoMatch`) silently converts `guard_unevaluable` or `owned_state_unavailable` rows from "reported candidate with facts" back into "silently dropped," reintroducing exactly the defect class 0005/0011 both exist to fix — with no distinguishing test failure until someone notices a row vanished. | premortem |
| C-8 | `0011:§implementation-plan`, Minimum Viable Validation step 6 and Testing Strategy S3 | The "three reachable match cases" fixture (present-equal, present-unequal, absent) is specified as ONE scenario across a hand-built model, but the RDR nowhere requires a fixture that exercises MULTIPLE match keys on the SAME row where one is present-equal, one present-unequal, and one absent simultaneously — the compound case that most resembles a real decision table (`models/rdr.toml` rows carry only a single match key each). The per-atom independence claim (A12: "no cross-atom interaction") is asserted, not tested against a multi-key row. | A row with two match keys — one satisfied, one not — is excluded correctly in isolation, but a row with two match keys where one is absent and one is unequal produces an ambiguous or wrong `unknown` payload (e.g., only one of the two facts surfaces, or exclusion happens for the wrong reason) because no fixture ever exercised the interaction, and the bug ships. | AT (acceptance tests) |
| C-9 | `0011:§decision-rationale`, "O3 wins on the two deciding rows" table, "peer-contract compatibility" row | The decision to leave `0007:REQ-78` (three-valued match at the kernel) permanently deferred, and to implement the absent/present-unequal distinction as a CLI-side presence test instead, creates a structural split: `flow next` and `flow resolve` will PERMANENTLY disagree about an absent-match-key row (next: candidate; resolve: `no_match`, escapable). The RDR calls this "the same report-vs-select relationship the verbs already have for `guard_unevaluable`" but that analogy is inexact — `guard_unevaluable` is a REFUSAL kernel-side in both verbs; here `flow next`'s "candidate" and `flow resolve`'s "excluded, escaped by a default row" are materially different verdicts for the identical input, sourced from two different presence tests (CLI's map-based one vs. kernel's `TagSet`) that the RDR asserts (A11) but does not mechanically enforce stay in sync forever. | A skill runs `flow next`, sees a row listed as a candidate with `unknown: [{key: "size", reason: "absent"}]`, picks it as "the" next action, calls `flow resolve` with the same state, and gets refused (`no_match`, rescued by an unrelated escape/default row) — a report/select mismatch that looks like a bug to any caller who has not internalized this RDR's Decision Rationale. | §3, premortem |

## 1. The three most likely ways implementation goes wrong

**(a) A15's demand-set extension ships unverified, or verified only against the fixtures the RDR itself hand-picked, and a match-only owned key model silently degrades to `--all` in production.**

Root cause: `0011:A15` is stamped `Status: Pending`. Its own text says "To verify" and lays out three sub-claims that "the verification must establish" — meaning the RDR is proposing to lock (or has locked, if pass reads a later revision) a "large" profile change while one of its two central mechanisms (the demand-set widening in `internal/cli/flow_exec.go::invokedReaders`) is not yet checked against real source, only argued about. The passage that enables it is `0011:C1` lines 141–158: "The assembled view MUST actually carry the keys the predicate reads, so `internal/cli/flow_exec.go::invokedReaders` MUST add each row's MATCH-block owned keys to its demand set... Without that, an owned key a model MATCHES on but never writes, clears, or guards demands no reader, is absent from the view, and has its atom omitted from every probe — so every row becomes a candidate carrying `{key, absent}` and the default silently degrades to --all for that model, with no caller remedy." This is the RDR admitting, in its own normative text, that getting A15 wrong doesn't produce a build failure or a test failure visible in the shipped suite — it produces a silent regression to the OLD (defective) behavior for a whole model class, with the caller unable to fix it (`--tag` is refused `flow-tag-owned` for an owned key). I confirmed in `internal/cli/flow_exec.go::guardOwnedKeys` (lines 131–145) that today's code explicitly skips any atom whose `Block` is not `all`/`unless` — so the gap A15 must close is real and currently unpatched in the pre-implementation source. Symptom: a model author declares a `match`-only owned tag (a legitimate, RDR-blessed pattern per S8), ships it, and `flow next` never narrows for that model — the exact defect this RDR exists to fix, reappearing through the one code path the RDR itself flags as unverified.

**(b) The `excluded()` signature change (bool → `Result`) is a bigger refactor than the RDR frames it as, and the refusal-kind dispatch gets the precedence wrong.**

Root cause: `0011:§technical-design` calls this "the one signature change in the file, and the reason the probe result is kept rather than discarded" — language that undersells the actual surface area. The function currently (`internal/cli/flow_next.go::excluded`, confirmed by direct read) returns a bare `bool` computed from `result.Refusal.Kind == resolve.KindNoMatch`. Under C1, the SAME function must now also feed `summarize()` a `Refusal.Undecided` payload for `guard_unevaluable`, while `owned_state_unavailable` must NOT try to read that payload (A13's "one fidelity limit," lines 200–205: "when the probe refuses `owned_state_unavailable`, `gate` returns before it collects the guard payload, so an `uncomparable` guard atom on that row is not on the wire and is not reported"). That means the CLI-side dispatch after calling `excluded`/`resolve.Resolve` must correctly branch on THREE refusal kinds (`no_match` excludes; `guard_unevaluable` reports with payload; `owned_state_unavailable` reports without payload) plus the success path, where today it branches on one boolean. Every one of the RDR's own mini-check tables (`disposition`, `trace`) enumerates this precedence as a design fact but the RDR provides no unit-level oracle isolating "does the CLI correctly distinguish owned_state_unavailable-without-payload from guard_unevaluable-with-payload" as its own assertion — S7 comes closest but is folded into a single scenario alongside the `uncomparable` fixture. Symptom: a row that should report `{key, absent}` for a missing owned key (and nothing else) instead also emits (or omits) an unrelated `uncomparable` guard fact, or vice versa — a payload-fidelity bug that manifests as a caller seeing wrong or missing diagnostic detail, silently, because nothing distinguishes "wrong facts" from "some facts" in casual testing.

**(c) The bundled `unresolved`→`unknown` payload rename breaks external consumers in the same release as the predicate change, and the "mechanical" five-site census undercounts the real blast radius.**

Root cause: `0011:Overrides` (metadata) and `0011:C3` bundle a wire-format-breaking rename into the same RDR as the behavioral change. The RDR itself distinguishes these as independent axes ("The two migrations are independent and a caller may owe both," `0011:§trade-offs` Consequences) but does not decouple them operationally — there is no flag, no versioned payload, no transition period; both land in the same build. The RDR's own text in `0011:C3` (lines 292–338) polices a five-site "mechanical" census of *shipped test reads* very precisely (verified correct against source at lines 103, 163, 166, 404 in `flow_next_0005_test.go`, plus `flow_adversarial_0005_test.go:446` and `flow_mvv_0005_test.go:66`), but that census is explicitly scoped to IN-REPO test reads — "The count is of test reads only." Nothing in the RDR inventories or warns about OUT-OF-REPO consumers of the JSON payload (skills, downstream tools) that parse `candidates[].unresolved` as a flat string array. `intrastate` is explicitly a generic library meant to be consumed by skills/tools outside this repo (stated directly in Context: "intrastate stays generic — no consumer (RDR-process) knowledge in code, docs, or fixtures") — meaning by design there ARE out-of-repo consumers, and the RDR's own "mechanical, low risk" framing of the rename only accounts for the risk it can see (in-repo tests), not the risk it structurally cannot see (external parsers). Symptom: every downstream skill/tool that reads `flow next`'s JSON output breaks on the same day it would have started getting a narrower, more useful candidate list — maximizing the pain of adopting the fix.

## 2. The one section that will be rewritten within 6 weeks of shipping

**`0011:§load-bearing-decisions` — "Undecided reporting shape" and its `--all` asymmetry, specifically the sentence: "One asymmetry against the cited precedents is accepted: `docker ps -a` and `git branch -a` widen to rows carrying the SAME information as the default's, whereas here the widened rows carry strictly LESS."**

This is the section the RDR itself flags as an accepted wart rather than a solved problem — it is candid that `--all` breaks the precedent it cites to justify its own name. In practice, `--all` is the verb's designated escape hatch for "I don't trust/understand the narrowed list," and the very first thing any caller will do with it is ask "okay, so WHY was row X excluded" — which `--all` cannot answer, by C2's explicit design (BlockMatch atoms are filtered OUT of `unknown` under `--all`). The RDR's own Briefly Rejected section names the fix (`excluded_by` field) and explicitly defers it: "If that proves to be the common case in use, the successor is `excluded_by` on a `--why`-style opt-in, not a default-mode payload widening — recorded here so the option is reopened deliberately rather than rediscovered." That sentence is a prediction of its own future amendment. Given that `flow next` exists specifically to be driven by skills/agents that need machine-actionable diagnosis (the whole Problem Statement is "the caller ... cannot constrain a skill's next step from it"), an automated caller hitting an excluded row with no `excluded_by` field has no recourse but string-matching model source files — which is not automatable at all for most skill callers. This gets escalated back into the RDR corpus almost immediately after the first real "why is my candidate list empty and confusing" support incident.

## 3. The one assumption that will not survive first contact with a real user

**A11: "CLI presence agrees with kernel presence for every key a row's match atom can name... so a match atom the CLI keeps because its key is present is never over a key the kernel's view lacks, and vice versa."**

This assumption is *correctly verified against the current source* — I confirmed both `assembledView` (CLI) and `assemble` (kernel) are built from the same `owned` slice with `OwnedSnapshot` uniformly omitting absent keys (`internal/accessor/model.go`, confirmed by the RDR's own citation, not independently re-verified by me but the mechanism described matches the accessor pattern seen elsewhere). The problem isn't that A11 is false today — it's that A11 is a STATIC snapshot of an invariant across TWO independently-evolving views (`flow_next.go::assembledView`, a plain `map[string]string`; `resolve.go::assemble`, a `TagSet` with provenance and conflict tracking) that happen to agree now because they're both built from one `owned []resolve.Tag` slice passed through unmodified. Nothing structurally *enforces* that agreement going forward — it's a coincidence of the current call graph (`runFlowNext` builds `owned` once and hands it to both `assembledView` and, inside `excluded`, to `resolve.Resolve`'s `Owned` field). A real user (in practice: the next engineer who touches `flow_exec.go` for an unrelated reason — caching, batching accessor calls, adding a new tag provenance) will refactor one side without realizing the other side's silent dependency on "same slice, same semantics." This assumption is exactly the kind of "true by construction today, unenforced tomorrow" invariant that erodes in every codebase under normal maintenance pressure, and the RDR provides no test that would fail specifically because A11 broke — only tests that would fail because SOME downstream symptom appeared, which is a much weaker signal to a maintainer debugging why `flow next` disagrees with itself.

## 4. The premortem

*Six weeks after RDR 0011 ships.*

`flow next` now reports a match-conditioned candidate list, and the `1mv1` defect that motivated this RDR is closed. But three tickets have landed against `internal/cli/flow_next.go`, and a fourth against the RDR corpus itself.

Ticket 1: a skill author using `models/consumer.toml` (structurally similar to `models/rdr.toml` but authored independently, outside this repo, per the "intrastate stays generic" design constraint) reports that `flow next` returns EVERY row as a candidate again, each with the same `{key, absent}` entry, for a model that used to narrow correctly in a prototype build. Root cause: the model declares an owned tag matched by several rules but never written, cleared, or guarded — precisely the shape `0011:S8` exists to test — and the shipped implementation of `internal/cli/flow_exec.go::invokedReaders` has a bug in the match-block-owned-key union (an off-by-one in `guardOwnedKeys`'s sibling function, or a missed edge case where the key also appears in an escape row and is filtered out by the `outcome != row.Outcome` skip in `invokedReaders`, which the RDR's own A15 evidence flags as a live edge: "A row binding a different outcome cannot serve this request... Retaining one used to be harmless because escape rows carry no `RequiresOwned`; once guard-owned keys joined the demand set (DEV-8) it stopped being harmless"). Because A15 shipped as `Pending` and its "MVV Test" verification ran only against the RDR's own hand-picked fixture (S8), the interaction with escape-row outcome filtering in a THIRD-PARTY model was never exercised. This is exactly failure mode (a) from Section 1, materialized.

Ticket 2: a different skill author reports that `flow next --evaluate-gates` is now invoking gate accessors on rows that should be excluded — a regression from the invariant "REQ-42: gates must not run on a guard-excluded row." Investigation traces it to the `excluded()` refactor: the engineer who implemented C1 changed the boolean check `if excluded(row, owned, req.observed) { continue }` to consult the new `Result`-returning function, but in doing so introduced a path where a row refused `owned_state_unavailable` (which should still be a REPORTED candidate per C1) was miscategorized as excluded in one code branch (a copy-paste of the `no_match` check that didn't get the kind comparison right), so some `owned_state_unavailable` rows silently vanish from the list instead of being reported with their missing-owned-key fact — the OPPOSITE of the design intent, and undetectable by the shipped test suite because S7's fixture for this precedence limit was scoped narrowly to the guard-payload case, not the row-disposition case. This is failure mode (b) from Section 1.

Ticket 3 (RDR corpus): three weeks after the previous two tickets, a follow-up RDR is opened to add the `excluded_by` field the Briefly Rejected section deferred — because, predictably, the top support request against `flow next` since launch has been "how do I find out WHY this row isn't a candidate without `--all`, and how do I find out why it's STILL not a candidate under `--all` either." This is Section 2's prediction landing on schedule.

Ticket 4: an external consumer (a CI pipeline driving `intrastate` as a library from another repo) files a compatibility break report: their parser for `candidates[].unresolved` (a `[]string`) now gets a type error because the field is `unknown` (`[]object`). They had NOT opted into the new predicate behavior deliberately — they simply upgraded the binary for an unrelated patch and their integration silently started failing JSON deserialization, because `intrastate`'s release process does not gate the rename behind any flag or major-version signal beyond ordinary semver, and the RDR's own text acknowledges "an out-of-repo consumer parsing the old field must change even if it adopts `--all`" without providing any mechanism to soften that for consumers who did not ask for the new predicate at all. This is failure mode (c) from Section 1.

Underlying all four: this RDR is enormous relative to what it changes. Three normative contracts, fifteen assumptions (one left `Pending` at what should be lock time), five mini-check tables, and ~800 lines of prose to change "the CLI stops stripping the match block before asking the kernel a yes/no question, and reports absent keys instead of dropping rows." The volume of grounding prose is not, on its own, a defect — much of it is genuinely well-verified against real source, which I confirmed independently for A2, A3, A6, and the mechanical structure of `KernelRow`, `checkAccessorBindings`, and `evaluateAtoms`. But volume this large is also exactly the condition under which an implementer's actual code diverges from the RDR's careful reasoning in exactly the corners the reasoning was most careful about — because there is too much to hold in working memory during a single implementation pass, and the RDR's own S6 admits the shipped test suite is *not* sufficient evidence the core contract landed correctly ("a green 0005 suite is NOT evidence C1 shipped — S2/S3 are"). When the safety net is explicitly not the ambient test suite but a small, hand-enumerated set of NEW fixtures (S1–S8), the actual defense against implementation drift rests entirely on those new fixtures being complete — and Section 1(a)'s finding shows at least one dimension (cross-model, cross-outcome interaction with escape rows) that S8 does not cover.

## 5. Acceptance tests that would have caught each failure at RDR-review time

```gherkin
Feature: Match-only owned key demand-set extension survives realistic model shapes (catches C-1)

  Scenario: A match-only owned key is also matched by an escape row for a DIFFERENT outcome
    Given a model where rule R1 (ordinary, outcome "advance") matches on owned key "tier"
      And rule R2 (escape, outcome "hold", rescuing no_match) also matches on owned key "tier"
      And "tier" is served by exactly one reader and written, cleared, and guarded by no row
    When "flow next" is run for outcome "advance" with no --tag
    Then the reader serving "tier" appears in "readers"
      And "tier" is present in the assembled view
      And no candidate reports "tier" as unknown/absent solely because R2's differing outcome
        caused invokedReaders to skip a key R1 still needs

  Scenario: The demand-set extension is verified with Status: Verified before lock, not Pending
    Given RDR 0011's Finalization Gate: Assumption Verification
    Then no assumption backing a Normative Contract may remain "Pending" at lock
      And A15 specifically must carry a Verified stamp with Method "Source Search" or "Spike",
        not "MVV Test" alone, since an MVV test proves the RDR's OWN fixture, not the
        general claim "invokes strictly more readers and never fewer" over arbitrary models


Feature: excluded() refusal-kind dispatch preserves the three-way disposition (catches C-2, C-7)

  Scenario: owned_state_unavailable is reported as a candidate, not excluded, and carries no guard payload
    Given a row R whose match holds, whose guard is decided, but whose RequiresOwned key
      is missing from the assembled view
      And R's guard ALSO contains an atom the seam could not decide (uncomparable), on a
      DIFFERENT present key
    When "flow next" runs without --tag for that key
    Then R is reported as a candidate (not excluded)
      And R's unknown list contains the missing owned key with reason "absent"
      And R's unknown list does NOT contain the uncomparable guard fact
        (gate returns before the undecided-guard collection runs; A13's stated limit)

  Scenario: guard_unevaluable is reported as a candidate with its undecided payload intact
    Given a row R whose match holds and whose guard contains a present-key atom the
      seam cannot decide, and no missing owned key
    When "flow next" runs
    Then R is reported as a candidate
      And R's unknown list contains {key, "uncomparable"} read off Refusal.Undecided

  Scenario: no_match is the ONLY disposition that excludes
    Given a row R whose match key is present and unequal
    When "flow next" runs, with and without --evaluate-gates
    Then R does not appear in candidates
      And no gate accessor for R is invoked


Feature: The unresolved-to-unknown payload rename is decoupled from the predicate change (catches C-5)

  Scenario: An external consumer can adopt the narrower candidate predicate without a payload
    schema break, or is given an explicit migration signal distinct from ordinary patch releases
    Given a consumer parses candidates[].unresolved as []string today
    When they upgrade to the release containing RDR 0011
    Then either:
      (a) the old field name/shape is still parseable for one deprecation window, or
      (b) the CLI's own versioning signals (e.g. a payload version marker, a major-version bump,
          or a documented breaking-change note distinct from routine patch notes) makes the
          break impossible to receive silently


Feature: --all provides an actionable diagnosis, not just a wider list (catches C-4)

  Scenario: A caller uses --all to find out why a row was match-excluded by default
    Given a row R excluded by default because match key "size" is present and unequal
    When "flow next --all" is run
    Then R appears in candidates
      And R's payload names WHICH match atom(s) would have excluded it by default
        (this fails against the current RDR design, which strips BlockMatch atoms from
        unknown under --all by C2's explicit mandate — the acceptance test that would have
        caught this is exactly the one the RDR's own Briefly Rejected section predicts
        will be needed and defers)


Feature: CLI presence and kernel presence cannot silently diverge (catches C-3)

  Scenario: A regression test pins A11's invariant directly, not just its symptoms
    Given any owned tag snapshot and any observed tag set flow next assembles
    When assembledView(owned, observed) and resolve.assemble's key set (owned ∪ observed ∪ recognized)
      are compared
    Then the two key sets are equal apart from the reserved "recognized" key
      (a property-based or table-driven test asserting this directly, independent of any
      specific model fixture, so a future refactor of runReaders or assembledView that
      breaks the coincidence fails a NAMED test rather than a downstream symptom)


Feature: Multi-key rows exercise match-atom independence under real-shaped models (catches C-8)

  Scenario: A row with two match keys in different presence/equality states reports both facts independently
    Given a row R matching on key A (present, unequal) and key B (absent)
    When "flow next" is run
    Then R does not appear in candidates (no_match on key A alone is sufficient to exclude)
    Given a row R' matching on key A (present, equal) and key B (absent)
    When "flow next" is run
    Then R' appears in candidates with unknown containing exactly {B, absent}, not A
```

## Verdict

O3 (the chosen design) is the least-bad of the five options weighed, but the RDR ships a fifteen-assumption, three-contract, near-800-line specification for a change whose actual code delta is one function's control flow, one demand-set union, and a field rename — and it leaves the single riskiest mechanism (A15, the reader-demand-set extension) unverified at what should be lock time. The design's own choices — deferring `--all`'s diagnostic value, permanently splitting `flow next`/`flow resolve` disposition on absent-match rows, and bundling a wire-breaking rename into the predicate fix — are each individually defensible but together guarantee at least one fast follow-up RDR and at least one production incident from an out-of-repo consumer who did not ask for any of this in the same release.
