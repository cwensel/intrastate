Model: claude-opus-5[1m]

# 3amigo iter-3 Persona 3 — QA / Tester (RDR cli/0011)

Read via projector: all of `0011:S1`–`S7`, `0011:MVV`, `0011:C1`/`C2`/`C3` (owned set), and — WIDENED, because several scenarios name a fixture, a payload shape, or a rendering whose pass/fail criterion is not stated inside the scenario — `§approach`, `§mini-checks` (`authority`/`oracle`/`disposition`/`trace`), `§load-bearing-decisions` (`0011:D-undecided-reporting-shape`, `0011:D-identity`), `§implementation-plan` Phase 1–3, `§trade-offs`, and assumptions `A1`–`A6`, `A11`–`A14`. What sent me out: S2 and MVV 5 name a fixture row by literal id; S4 names a fixture shape; C1's `unknown` merge is a payload whose text rendering only `D-undecided-reporting-shape` describes. Source checked: `internal/cli/flow_next.go`, `internal/cli/flow_fixtures_0005_test.go`, `internal/cli/flow_next_0005_test.go`, `internal/cli/flow_mvv_0005_test.go`, `internal/cli/flow_adversarial_0005_test.go`, `internal/cli/flow_surface_0005_test.go`, `internal/cli/flow_harness_0005_test.go`, `internal/cli/respond/text.go`, `internal/resolve/guard.go`, `internal/resolve/resolve.go`, `internal/table/normalize.go`, `internal/table/model.go`.

---

## BLOCKER 1 — `0011:S2` / `0011:MVV` step 5 / `0011:§mini-checks` trace row 5: the named fixture does not contain the row the scenario asserts on

`0011:MVV` step 5 says "Run over the 0005 MVV fixture (`flowMVVModel`, `status=draft`) with and without `--all`; assert the `gated-excluded` row is absent in both". `0011:S2`'s Expected repeats it ("`gated-excluded` absent in both"), and the `trace` mini-check's step-5 row repeats it a third time ("`flowMVVModel` ± `--all` … `gated-excluded` absent in both — guard-excluded").

`gated-excluded` is not a row of `flowMVVModel`. `flowMVVModel` is `internal/cli/flow_fixtures_0005_test.go:77-165` and declares exactly two ordinary rules, `advance-draft` and `hold-draft`; neither carries a `guard` block at all. `gated-excluded` is a row of a DIFFERENT fixture, `flowGatedNextModel` (`internal/cli/flow_fixtures_0005_test.go:317`), which is the model the two shipped oracles that mention the row actually load (`flow_next_0005_test.go:198`, `flow_mvv_0005_test.go:455`).

Test this prevents: the guard-exclusion-unchanged half of S2. Written literally, the assertion passes VACUOUSLY — `gated-excluded` is absent from a `flowMVVModel` run in every possible build, including one that broke guard exclusion entirely, because the row does not exist in that model. This is precisely the anti-oracle failure `0011:S6` and the `oracle` mini-check warn about elsewhere in this same record, reintroduced in the scenario that is supposed to carry the discriminating evidence. As written the implementer must either (a) load `flowGatedNextModel` — but then the `status=draft` seed and the "row C3 adds" clause in the same sentence do not apply, since C3's added row is specified against `flowMVVModel`'s `match.status` domain; or (b) add a guard-excluded row to `flowMVVModel` — a fixture edit neither `0011:C3` nor Phase 3 authorises, and one that ripples into the 64 uses of `flowMVVModel` across 14 files plus the two derived constants `flowEscapeModel` and `flowEscapeOtherOutcomeModel` (`flow_fixtures_0005_test.go:171`, `:770`). Neither branch has a stated pass/fail criterion.

Decision blocked: which fixture S2 runs over, and therefore whether the guard-exclusion negative control the `oracle` mini-check pairs against MVV 5 exists at all.

---

## BLOCKER 2 — `0011:S4`: the "empty match pattern" fixture is refused at load; the scenario asks for a row that cannot exist

`0011:S4` reads: "a row whose only match atom is `recognized`, and a row with an empty match pattern (A5's fixture B)". Its Expected: "both are candidates with `unknown` empty of match facts".

The grammatical reading is two rows. An AUTHORED empty match pattern is unbuildable: `internal/table/normalize.go:379` refuses any transition rule whose `[rule.match]` key is absent, with the comment stating explicitly that "both an absent `[rule.match]` and a present-but-empty one decode to nil, and both are equally malformed" (`CatMalformedRuleShape`, "rule %s carries no match block"). So no loadable model can contribute a second row of that shape, and the second half of S4 has no achievable input.

The charitable reading is ONE row — a `recognized`-only rule, whose match atom normalize lifts into `Row.Outcome`, leaving a normalized (kernel-level) match pattern that is empty. A5's fixture B is exactly and only that: `row-recognized-only`, one rule, `[rule.match.recognized] eq = "fallback"` (spike file, lines 528-533). Under that reading S4's two clauses name the same row twice and the scenario asserts nothing about the class C1's sentence "A row ALL of whose match atoms are omitted is probed with an empty match pattern, which matches unconditionally" actually mints — a row with real match atoms, all of them over ABSENT keys, which is a distinct input class from a `recognized`-only row (its `unknown` is non-empty, carrying `{key, absent}` per key, whereas S4's Expected demands "`unknown` empty of match facts").

Test this prevents: the C1 empty-probe clause has no scenario. On the two-row reading the fixture cannot be built; on the one-row reading S4's Expected ("`unknown` empty of match facts") is correct for the `recognized`-only row and WRONG for the all-atoms-omitted row, so an implementer cannot write the assertion for the latter without contradicting S4.

Decision blocked: whether "empty match pattern" in `0011:C1` means authored-empty (unbuildable) or all-omitted (buildable, but with the opposite `unknown` expectation from the one S4 states).

---

## MAJOR 3 — `0011:D-undecided-reporting-shape` vs `0011:S1`–`S7` and `0011:MVV`: a text-mode rendering is mandated with no scenario, and it contradicts the shipped renderer's ownership

`0011:D-undecided-reporting-shape` states: "Text mode (`--as=text`) renders the same pairs as `key (reason)` on the candidate; the reason vocabulary carries the whole remedy story, so suppressing it in the human-facing mode would leave the text caller with the symptom and no diagnosis".

Two problems, both testability problems.

(a) No scenario covers it. Every one of `0011:S1`–`S7` and all eight `0011:MVV` steps run `--as=json` or assert on flag sets and suite greenness. `--as=text` appears in no scenario, no MVV step, and no row of the `oracle` or `disposition` mini-checks. A rendering the record calls load-bearing enough to argue about has no pass/fail criterion anywhere.

(b) The stated rendering is not the one the code can produce without a change the record does not scope. `internal/cli/respond/text.go::writeTextPayload` is respond-owned and generic: it marshals the payload and emits path-qualified leaf lines via `flatten`. A `[]{key, reason}` list renders as `candidates[0].unknown[0].key: status` and `candidates[0].unknown[0].reason: absent` — never `key (reason)`. Producing `key (reason)` requires either a per-verb text template or a custom marshaller, and `text.go`'s own header names that as the drift it exists to prevent (REQ-11/REQ-120: "a verb never prints its payload itself"; `internal/cli/flow_surface_0005_test.go:322-324` carries the REQ-11/REQ-120 oracles). Phase 1–3 scopes no `respond` change.

Test this prevents: any text-mode oracle for `unknown`. An implementer following `D-undecided-reporting-shape` literally writes a test expecting `status (absent)` and it fails against a build that satisfies every other clause of C1/C2; an implementer following the shipped renderer writes nothing, and the decision goes unasserted either way.

---

## MAJOR 4 — `0011:C2` / `0011:S2`: the structural `--all` negative names a positive-only oracle as its home, so the negative has no stated site

`0011:C2` requires the `--all` negative be pinned "STRUCTURALLY, asserting `--all` is absent from those three verbs' flag sets and from the flow group's own, guarded against vacuity by requiring the four verbs to exist first — the idiom already shipped for an absent flag (`flow_input_0005_test.go::TestReq69_NoPlanFlagShipsOnAnyVerb`, `flow_surface_0005_test.go::TestReq12_FlowGroupDoesNotRedefineTheRootAsFlag`)". `0011:C3` then requires the build "register `--all` on next in the per-verb flag map the surface suite already carries (`flow_surface_0005_test.go::TestReq3_EachVerbRegistersItsNormativeFlagSpellings`)".

`TestReq3` (`flow_surface_0005_test.go:128-163`) asserts PRESENCE only. Its `want` map lists, per verb, flags that MUST be registered, and the body only ever reports `sub.Flags().Lookup(f) == nil`. Adding `"all"` to `next`'s entry satisfies C3's registration clause and asserts nothing whatsoever about the other three verbs — the map has no absent-flag column, and nothing in it is negative. So C2's structural negative and C3's registration clause are two different oracles, and the record names a home for the second while leaving the first homeless: it cites the two absent-flag idioms as PRECEDENT ("the idiom already shipped") without saying whether the new oracle is added to `flow_surface_0005_test.go`, to `flow_input_0005_test.go`, or to a new `0011` file, and without fixing its name.

Test this prevents: nothing outright — the oracle is writable — but the record leaves it ambiguous whether one test or two is owed, and a build that only edits `TestReq3`'s map can claim C2's negative is satisfied. Since `0011:S2` says only "C2's negative is asserted STRUCTURALLY", it does not close the gap either. Given C3's own insistence that a green suite must not be cited as evidence, the negative deserves the same explicit siting the five `unresolved` re-homings get.

---

## MAJOR 5 — `0011:S3` / `0011:MVV` step 6 / `0011:C3`: the discriminating fixtures are specified by row property, not by model, and one of them collides with a load refusal

`0011:C3` mandates "fixtures that discriminate: a row whose match key is present and UNEQUAL … and a row whose match key is ABSENT". `0011:S3` runs "the three reachable match cases, one row each (MVV 6). Backed by the A5 spike's fixture shape", and `0011:MVV` step 6 adds "Then add `--tag key=<value>` and assert the list narrows to the rows whose match holds."

The MISSING criterion is the model. "Backed by the A5 spike's fixture shape" points at `dt-write.toml`, whose four rows match on two OBSERVED dimensions — and the spike itself records (its own §"Constraint discovered first") that this shape only loads because it carries a dummy owned tag `answer` plus the reader/writer pair it drags in, because `normalize.go:366` refuses an ordinary rule with no write block and a write to an observed tag is refused `write_to_non_owned_tag`. So the S3 fixture is not a free choice: whatever model the implementer writes must carry a dummy owned key, and that key will appear in every candidate's `required` and — until the `nav` reader establishes it — in every candidate's `unknown` as `{answer, absent}`.

S3's Expected states "present-equal ⇒ candidate, key absent from `unknown`". Under the spike's own fixture that is FALSE at the payload level: the candidate's `unknown` is non-empty, carrying the dummy owned key. The Expected is only true if read as "the MATCH key is absent from `unknown`", which is not what it says, and the difference decides whether an oracle asserting `len(unknown) == 0` passes or fails.

Test this prevents: the present-equal oracle, and by extension MVV 6's "present-and-equal (candidate, nothing in `unknown`)" — same wording, same ambiguity. As written, "nothing in `unknown`" is unachievable on any model of the shape S3 cites unless the artifact seeds the dummy owned key too, which no step says.

---

## MINOR 6 — `0011:MVV` step 7 / `0011:S2`: "`unknown` present as `[]` … on a fully-resolved candidate" names no such candidate in any fixture

`0011:MVV` step 7 and `0011:S2`'s Expected both require asserting `unknown` is `[]` rather than omitted "on a candidate with nothing undecided, in both modes". `0011:C2` makes this the shape guarantee ("`unknown` is present in both, as `[]` rather than omitted when it carries nothing, so one consumer struct parses either mode"), and the `oracle` mini-check's MVV-7 row makes it a named oracle.

No fixture in the record produces such a candidate. In `flowMVVModel` at `status=draft`, `advance-draft` carries `gate = ["approval"]`, so without `--evaluate-gates` its `unknown` carries `{approval, not-evaluated}` (C1's new reason); `hold-draft` carries no gate and no guard, so it is the only shipped candidate that could have an empty list — but only if the seed establishes every owned key the row requires, and the shipped seeds are `status=draft` alone (`flow_next_0005_test.go:74`, `flow_mvv_0005_test.go`), leaving `labels`/`stale`/`note` unestablished. `models/rdr.toml` (S1) has `gate_passed` and gates throughout. So the implementer must construct a candidate the record never specifies, and cannot tell from the record whether an unestablished non-required owned key contributes an `unknown` entry — C1 says "every owned key no invoked reader established", and today's `summarize` walks `row.RequiresOwned`, not all declared owned keys, which are different sets.

Test this prevents: the `[]`-not-omitted oracle. Also leaves open whether C1's "every owned key no invoked reader established" preserves `summarize`'s current `RequiresOwned` scope or widens it — a silent widening would change `unknown` on every candidate and break the empty-list assertion.

---

## MINOR 7 — `0011:C1` / `0011:S7`: `unknown` sort and dedup are stated as behaviour but no scenario asserts them

`0011:C1` fixes the merge: "Entries are deduplicated by `{key, reason}` pair, not by key, and sorted by `(key, reason)`." `0011:D-undecided-reporting-shape` explains why ("an unstable order would churn any golden payload a consumer diffs"), and `0011:D-identity` only requires set equality.

`0011:S7` asserts the dedup case in one direction only — "the same atom over an absent key appears once as `{key, absent}` (walk and payload agree, dedup on the pair)" — and the `oracle` mini-check's S7 row echoes it. Nothing asserts SORT, and nothing asserts the not-deduplicated-by-key direction: a key that legitimately carries two different reasons. Per C1 that pair-not-key rule is load-bearing, but per `A13` the only key that can carry two reasons at once would need both an `absent` and an `uncomparable` entry, which is contradictory for one key, and the gate class carries `not-evaluated` on a gate ID that is "not a view key" — so the record does not establish that any input can produce two entries for one key. That makes "deduplicated by `{key, reason}` pair, not by key" either untestable or vacuous, and the record does not say which.

Test this prevents: (a) any sort-stability oracle — `0011:S1`–`S7` name none, so an implementation returning `unknown` in row-atom order passes every stated scenario; (b) the pair-vs-key dedup discriminator — a build that deduplicates by KEY alone passes every scenario in the record, because no scenario supplies an input where the two rules differ.

---

## MINOR 8 — `0011:S6` / `0011:C3`: "none re-homed under `--all`" is asserted, but the re-homing rule that would produce a move is stated as conditional

`0011:S6`'s Expected: "the `flow next` oracles pass, none re-homed under `--all`", and `0011:MVV` step 8 requires "the 0005 `next` suite passes with only C3's five mechanical `unresolved`→`unknown` re-homings". `0011:A4` supplies the count (MOVES-UNDER-`--all` 0, BREAKS 0) and is Verified.

`0011:C3` nonetheless states a conditional obligation: "an assertion that depended on a match-excluded row being reported MUST be re-homed under `--all`, not deleted". Per A4 no such assertion exists, so the clause is currently vacuous — which is fine — but S6 turns its vacuity into a positive assertion ("none re-homed") without stating how a build DEMONSTRATES it. Counting diff hunks is not an oracle; `go test ./internal/cli` green is exactly the anti-oracle S6 itself disclaims two sentences later.

Test this prevents: nothing critical, but "none re-homed under `--all`" as written is a claim about the diff, not about the program, and has no automatable pass/fail criterion. Restating it as "no `flow next` test in `internal/cli` passes `--all`" would make it grep-checkable; as it stands it is a review obligation dressed as a scenario expectation.

---

## OBSERVATION 9 — `0011:S5`: the kernel-untouched oracle is stated twice with two different criteria

`0011:S5`'s Expected includes "`go test ./internal/resolve` passes with no test file changed"; `0011:MVV` step 8 states "`go test ./internal/resolve` passes with no test changed and `git diff --stat internal/resolve` empty"; the `oracle` mini-check's S5-kernel row says "with NO test changed". Three statements, and only the MVV one carries the `git diff --stat` check that actually makes the claim mechanical — "no test file changed" is not observable from a test run. Not a blocker (the MVV form is checkable and is the strongest of the three), but S5 and the mini-check row should carry the same criterion the MVV does, or an implementer reading S5 alone has no way to check it. Verified against source that the three named kernel oracles exist and are reachable: `TestReq49`/`TestReq47` pin `Reasons()` at exactly `{absent, uncomparable}` (`internal/resolve/guard.go:81-83`), and the `owned_state_unavailable`-before-`guard_unevaluable` precedence S7 pins is real (`internal/resolve/resolve.go:532-538` returns before the `undecided` collection loop at `:540-556`).

---

## Verified-as-testable (no finding)

For the record, these I checked against source and found buildable with stated criteria: `0011:C3`'s five `unresolved` reads resolve exactly at the cited lines (`flow_next_0005_test.go:103`, `:163`, `:404`; `flow_adversarial_0005_test.go:446`; `flow_mvv_0005_test.go:66`), and `:166` is indeed the `%#v` argument of `:163`'s `t.Fatalf`, not a sixth read. `stringsAt` (`flow_harness_0005_test.go:426-444`) does return `nil,false` on a missing key, a non-array, or a non-string element, so C3's ok-bool requirement is well-founded and the three comma-ok discards are real. `0011:C2`'s `BlockMatch` filter is buildable — `table.Atom` retains `Block` (`internal/table/model.go:142-147`). `0011:S7`'s `uncomparable` producer is reachable as A13 says. `0011:S1`'s 21→3 count has a live spike behind it.
