# Recommendation 0011: flow next selects by match; --all enumerates the alphabet

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-08-26
- **Status**: Draft [route-back from Stage 5 critique 2026-08-26; re-verify A3, A7, A8, A9, A10 — the kernel-placement justification, the `Refusal.Undecided` carrier, and the `0007:C10` aggregation question must be re-decided before the lens row resumes]
- **Type**: Feature
- **Profile**: large — three independent contracts: the three-valued match verdict at the kernel seam (C1, a semantics change consumed by `flow resolve` as well as `flow next`), the `--all` enumeration surface (C2), and the `unknown` reporting shape plus its discriminating fixtures (C3)
- **Priority**: High
- **Related Issues**: kata `1mv1` (defect tracker); umbrella rdr#thsc; consumer model rdr#tmxk
- **Predecessors**: 0005-skill-integration-cli-contract
- **Overrides**: `0005:C1`, the `flow next` candidate clause — "flow next MUST return the legal recognized-outcome alphabet for the supplied state, plus candidate summaries …" read as guard-conditioned, and its stated ground "because next enumerates rather than selects" — narrowed to the match-conditioned predicate of C1 below; 0005's enumeration becomes `--all` (C2). Every other `0005:C1` obligation on `flow next` is unchanged. `internal/cli/flow_next_0005_test.go` assertions that relied on a match-excluded row being reported move under `--all` — verified in Stage 4 (A4) to be an empty set over the shipped fixtures, so the clause is a forward guard and C3 owes the discriminating fixtures instead. Also closes `0007:REQ-78`, which scoped the match seam out of RDR 0007's build rather than settling it, and corrects the harmlessness claim JDR 0001's Block-cardinality decision rests on (that a match atom is never unevaluable), since a conflicted key makes one unevaluable in exactly `0007:C8`'s `uncomparable` sense. Widens `0007:C8`'s `Refusal.Undecided` payload to admit match-block atoms (A8) and extends `guard_unevaluable` to cover an undecidable MATCH predicate at the kernel's two selection sites — a `flow resolve` behaviour change that retires `guard_atoms_test.go::TestReq78` and re-expects three `internal/resolve/match_conflicted_test.go` oracles (S5). Adds the reason token `not-evaluated` for un-run gate ids on the CLI's `unknown` list only; `0007:C8`'s seam vocabulary stays at `absent`/`uncomparable`.
- **Seam Lineage**: `internal/resolve` match-seam semantics — one prior deferral, `0007:REQ-78` ("the match pattern still folds absence into non-match … NOT closed by this RDR"), which this RDR closes. Below the ≥2 accretion floor; the profile above is set by the contract count, not by the floor.

## Problem Statement

A skill author (or the human driving one) runs `flow next --model <model> --artifact <state>` to ask "what is legal from here?" so the next action can be chosen without re-deriving the transition table by hand. Today the answer is wrong for that question: RDR 0005 made `next` enumerate rather than select, stripping each row's `match` (and `escape`) before probing, so every row not pruned by a guard is reported as a candidate — 21 of the 22 rules on `models/rdr.toml` at `stage=resolved`. The caller discovers this when the candidate list is the whole alphabet regardless of the supplied state, and cannot constrain a skill's next step from it (0005's own premortem, "next omits enough condition detail to constrain a skill", materialized).

The system-internal requirement is a single decision on the candidate predicate of `flow next`: whether "the legal recognized-outcome alphabet for the supplied state" means guard-conditioned (0005's predicate: candidate = not guard-excluded, match is the kernel's selection pattern and `next` must not select) or match-conditioned (candidate = match holds AND not guard-excluded, evaluated by handing the kernel the row with its match intact so selection semantics stay in the kernel); which of the two is the default and which rides `--all`; and whether guard verdicts on match-excluded rows are still reported. A sub-question inside the same contract, sharpened by Stage 4 into the load-bearing one: what a match atom means when its key cannot be compared — absent from the supplied state, or carrying two differing values within one provenance (two invoked readers writing the same owned key; a repeated `--tag` is refused upstream). Today `internal/resolve/resolve.go::TagSet.matches` folds absent, conflicted, and unequal into a single `false`, so any of them excludes a row identically and silently; `0007:REQ-78` recorded that as deferred rather than settled. `--all` is the same contract's second surface, not a second contract. Help text follows the decision.

**What "constrained" means here, and where the outcome lands.** The outcome is a candidate list narrowed to what the supplied state can take — NOT a single next action. `flow next` still never selects (`0005:D-selection-predicate`, unchanged); choosing among survivors is `flow resolve`'s, and a caller facing several candidates is being told the state genuinely admits several. Success is therefore "no row is listed that the supplied state excludes", not a target cardinality. The narrowing is real only to the extent the match keys are DECIDABLE from the invoked readers plus supplied tags: where a match key is owned and served by a reader (`models/rdr.toml`'s `stage`, A6) the list narrows 21→3; where it is observed and unsupplied (0010's decision-table class with no `--tag`, A5) every row remains a candidate and each names the undecided key under `unknown`. That second case is the old symptom's shape with a diagnosis attached, and it is the intended behaviour, not a residual defect — the alternative, dropping those rows, is Alternative 1, which this RDR rejects for emptying the list on partial state. The caller's remedy is named in the payload (bind the reader, supply the tag); the RDR does not promise narrowing where the state was never supplied.

## Critical Assumptions

- **A1 A one-row probe table handed to the kernel with the row's match atoms retained can only yield a plan, `no_match`, `owned_state_unavailable`, or `guard_unevaluable`; `ambiguous_match` cannot arise, so the sub-question "does `ambiguous_match` leave the row a candidate" is vacuous on a probe.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/resolve/resolve.go::RefusalKinds` closes the set at five. `ambiguous_match` is produced only under the `default:` arm of `switch len(selected)` in `internal/resolve/resolve.go::Resolve` and of `switch len(viable)` in `internal/resolve/resolve.go::escapeOrRefuse`, both requiring ≥2; a one-row table makes `len(candidates) ≤ 1`, and the stripped escape list leaves `escapeOrRefuse` no rescue row, so it passes the original refusal through. `internal/resolve/resolve.go::gate` returns `owned_state_unavailable` then `guard_unevaluable` before that count. The fifth kind, `unmodeled_outcome`, is unreachable **by probe construction, not by one-row-ness**: `internal/cli/flow_next.go::excluded` binds `Outcomes: []string{row.Outcome}` alongside `Recognized: row.Outcome`, so `Table.models` holds trivially — C1's probe builder MUST preserve that pairing. The two non-refusal error returns (`internal/resolve/resolve.go::Table.CheckValid`, `internal/resolve/precondition.go::CheckInput`) are vacuous on the probe (escape stripped; no reserved `recognized` key in owned/observed/`RequiresOwned`) and the CLI already treats `err != nil` as non-excluding.
  - **If wrong**: a match-retained row could be dropped as "ambiguous" and vanish from the candidate list with nothing in `unknown` explaining it.
  - **Scope note (3amigo)**: this enumeration is unchanged by C1's selection-site rule — an indeterminate match now lands in `guard_unevaluable`, already a member of the list, which is precisely why the probe needs no new disposition and `Plan` needs no new field. The vacuity of `ambiguous_match` is a property of the ONE-ROW probe only; on a real multi-row table `flow resolve` can reach it, which is why C1 states explicitly that an indeterminate row does not count toward `len(selected)`.
- **A2 The CLI already reports an absent match key: `internal/cli/flow_next.go::summarize` walks `Row.Atoms` — which carries match-block atoms alongside guard atoms (`internal/table/model.go::Row.KernelRow` routes them by `Block`) — against `internal/cli/flow_next.go::assembledView`, so a match key the view lacks lands in the candidate's undecided list today. It does NOT hold the test C1 needs: `assembledView` cannot distinguish a conflicted key from a settled one, which is why the verdict is the kernel's (Alternative 3).**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/table/model.go::Row.Atoms` is documented as the one unified predicate set with each atom's authored block retained ("the Match/guard split is applied at the kernel handoff, never stored here", `0002:C16`) — there is no separate match field or pattern map on `table.Row`. `internal/table/model.go::Row.KernelRow` routes that same field by `a.Block == BlockMatch`, and `internal/cli/flow_next.go::summarize`'s loop `for _, atom := range row.Atoms { if _, known := view[atom.Key]; known { continue } … }` carries no `Block` test, so a match-block atom over a view-absent key lands in `unresolved` today (the loop's own comment says "guard atom", narrower than the code). Note `internal/cli/flow_next.go::assembledView` does **not** inject `recognized`, unlike the kernel's `assemble`; this is harmless because `internal/table/normalize.go` lifts the `recognized` atom out of the predicate set into `Row.Outcome`, so no `recognized`-keyed atom reaches `Row.Atoms` at all.
  - **If wrong**: an absent match key silently drops out of the report and C1's `unknown` clause needs new code in the CLI rather than the existing atom walk.
- **A3 The kernel's match seam folds THREE distinct cases into one `false` — `internal/resolve/resolve.go::TagSet.matches` returns false when `!ok || tv.conflicted || tv.value != w.Value` — so the distinction C1 needs cannot be recovered by any caller and must be made inside that function. A CLI-side presence pre-strip cannot substitute: a conflicted key is present-yet-uncomparable, so a bare presence test retains its atom and the row is excluded as `no_match` with nothing reported.**
  - **Status**: Pending
  - **Method**: Source Search
  - **REFUTED at Stage 5 (critique), 2026-08-26 — the reachability half only.** `internal/table/load.go::checkAccessorBindings` refuses at MODEL LOAD any model where an owned tag is served by `!= 1` readers (`owned tag %s is served by %d readers; want exactly one`), and because it counts key OCCURRENCES a single reader declaring `keys = ["k","k"]` is refused by the same test (probed live). `internal/accessor` classifies a reader's output against its declared keys, so a reader cannot emit an undeclared repeat. `TagSet.merge` conflicts only WITHIN one provenance (`prior.provenance == prov`), and the observed route is refused by `parseTags` as `flow-tag-duplicate` — which this RDR already states. Therefore PRESENT-but-CONFLICTED is **unreachable from the CLI over any model that loads**. The seam-shape half of this assumption (that `matches` folds three cases into one `bool`) stands and is unaffected. What falls is the reachability claim the Approach's fact (i), `D-placement`, and the Alternative 3 rejection ("the failure is real, not theoretical") all rest on. Re-verify what remains: whether SINGLE-SOURCING the match rule — not CLI incapacity — still justifies kernel placement, given that `flow_next.go::summarize` already distinguishes an ABSENT key by comma-ok today.
  - **Prior evidence (superseded)**: `internal/resolve/resolve.go::TagSet.matches` is `if !ok || tv.conflicted || tv.value != w.Value { return false }` — a three-legged disjunction collapsing to one `bool`, whose only callers (candidate selection and escape selection in `internal/resolve/resolve.go::Resolve` / `::escapeOrRefuse`) `continue` on false with no record of why; the guard seam by contrast does carry the distinction (`internal/resolve/guard.go::evaluateAtom` returns `ReasonAbsent` vs `ReasonUncomparable`). Atom independence holds: `matches` is a plain all-must-hold loop with no accumulator; match-block `in` atoms are expanded into per-member `eq` rows at normalization (`internal/table/normalize.go`, `0002:C13`) and set-kinded literals are canonicalized to one opaque string by `internal/table/model.go::seamValue`, so every match atom reaching the kernel is an independent single-key equality and `matches(nil)` is `true`. The conflicted exception is real and user-reachable: `internal/resolve/resolve.go::TagSet.merge` marks a key repeated within one provenance with differing values as `conflicted` while KEEPING it present in `s.tags`. The reachable producer is the OWNED provenance — `internal/cli/flow_exec.go::flowRequest.runReaders` appends each invoked reader's `OwnedSnapshot()` into one slice with no dedup, so two readers writing the same key with differing values conflicts it. It is NOT a repeated `--tag`: `internal/cli/flow_input.go::parseTags` refuses that with `flow-tag-duplicate` before `kernelTags` is reached (corrected at 3amigo; the earlier wording named `kernelTags`, which passes a repeat through only because the refusal already happened upstream). `internal/cli/flow_next.go::assembledView` has no conflicted concept at all (zero occurrences of `conflicted` in `internal/cli`), so the CLI cannot currently detect the condition. `matches`'s own doc names the posture it mirrors: "`ReasonUncomparable` … the key is present, and its value was not compared to a verdict" (`0007:C8`).
  - **If wrong**: the kernel already carries an undecided match verdict, and C1's Phase 1 is a reporting change rather than a semantics change.
- **A4 The 0005 `flow next` tests move rather than rewrite: every oracle that asserts on a named candidate seeds the owned key its rows match on (`status=draft` in `flowMVVModel` / `flowGatedNextModel`), and the alphabet, gate-handling, reader-narrowing, and non-mutation oracles do not depend on a match-excluded row being reported.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: all 28 `flow next` oracles across `internal/cli/flow_next_0005_test.go`, `flow_mvv_0005_test.go`, `flow_adversarial_0005_test.go`, `flow_input_0005_test.go`, `flow_encoder_0005_test.go`, `flow_resolve_0005_test.go`, `flow_harness_0005_test.go`, and `reserved_key_0008_test.go` were enumerated and classified against C1's predicate. Every match atom in the `next` fixtures is `status eq "draft"` (`flowMVVModel`'s `advance-draft` / `hold-draft`; `flowGatedNextModel`; `flowGateDenyModel`; `advUnlessExcludedModel`; `advInExcludedModel`), and every seed supplies `status=draft`, so `status` is present in the view and holds in every case. Every exclusion the suite exercises is a **guard** exclusion (`all.flag eq false`, `unless.flag`, `all.tier in [mid,high]`), which C1(b) leaves untouched. Baseline green: `go test ./internal/cli -run 'Next|Req4[0-7]|Req3[5-9]|MVV|Adv1'` → ok. Counts, **against C1's candidate predicate**: PASSES-UNCHANGED 28, MOVES-UNDER-`--all` **0**, BREAKS **0**. This classification is predicate-scoped and does NOT clear the suite for the build as a whole: the `unresolved` → `unknown` payload change is an independent edit that mechanically breaks 6 assertions in the same suite (C3 names them).
  - **If wrong**: the 0005 suite splits into default / `--all` variants beyond the moves C3 names, and Testing Strategy must enumerate them.
- **A5 Over 0010's decision-table class, `flow next` with no `--tag` reports every ordinary rule and exits 0, each rule's observed match keys appearing in the candidate's undecided list — with the exception of a row whose only match atom is `recognized`, whose list is empty by construction — so `0010:A11` and 0010's MVV step 5 hold under C1; and a partial `--tag` set narrows the list to the rows whose supplied keys hold, a supplied-but-mismatched key dropping its row.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `docs/rdr/0011-flow-next-match-conditioned-candidates/evidence/spikes/a5-decision-table-no-tag.md`. Live over a four-row two-dimension decision table, no `--tag`: exit 0, all 4 ordinary rows reported, each carrying its observed match keys `status` and `size` in `unresolved`. Supplying `--tag status=Draft` moves `status` out of every `unresolved` while dropping no row *today* (match is stripped) — under C1 `status` is present in the view, so `final-small`/`final-large` are supplied-but-mismatched and the kernel refuses their probes `no_match`, dropping them; `draft-small`/`draft-large` survive with `unresolved = ["size", "answer"]`. Two limits recorded rather than papered over: (i) a genuine zero-owned 0010 decision table is **unloadable on this build** — an ordinary rule with no write block is refused `malformed_rule_shape`, and a write to an observed tag is refused `write_to_non_owned_tag` — so the fixture carries one dummy owned tag and `required` is `["answer"]`, never empty; the "empty `required`" half of `0010:A11` / `0010:MVV` step 5 is downstream of 0010 shipping and is **not attestable here**. (ii) The universal quantifier in this assumption's earlier wording is falsified by fixture B: a row whose only match atom is `recognized` reports `"unresolved": []` on the current build, because `internal/table/normalize.go` lifts `recognized` into `Row.Outcome` and C1 binds it to the row's own outcome, so it can never be an unresolved fact. That refutes the quantifier, not `0010:A11` or `0010:MVV` step 5, which speak only of exit 0 and empty `required`.
  - **If wrong**: 0010 owes a class-specific clause on this seam, or Alternative 1's strict predicate is what a decision table needs and the two classes diverge.
- **A6 `models/rdr.toml` at `stage=resolved` narrows from 21 reported rules to exactly the 3 rows whose `match.stage` holds over the artifact's owned `stage` — `prelock`, `resolve-route-back`, `resolve-abandon`, one per declared outcome. The 22nd row's absence from the baseline 21 is guard-driven, not stage-driven.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `docs/rdr/0011-flow-next-match-conditioned-candidates/evidence/spikes/a6-rdr-toml-narrowing.md`. Live against the current build, artifact `{"stage":"resolved","status":"draft","gate_passed":"false"}`: 22 `[[rule]]` blocks, 22 non-escape normalized rows, **21 candidates** — the one excluded row is `finalize-pass`, pruned on its `[rule.guard.all.gate_passed] eq = "true"`, so the 21 is a property of `gate_passed=false`, not of `stage=resolved`. A prototype of C1's probe over the same model returns **3** candidates; cross-checked by `grep -n 'eq = "resolved"' models/rdr.toml` → exactly three `[rule.match.stage]` hits (lines 212/222/232). Every rule carries a `[rule.match.stage]` block and `stage` is `required = true` and served by the invoked `rdr-status` reader, so C1's absent-key provision never fires for this model. The `recognized` sub-claim is dropped as vacuous: `internal/table/normalize.go` lifts `match.recognized` into `Row.Outcome`, leaving **0** such atoms in `Row.Atoms`; the probe's `Recognized` binding instead serves the kernel's `row.Outcome != in.Recognized` filter and keeps the probe from refusing `unmodeled_outcome` (A1).
  - **If wrong**: the observed defect is not closed by a match predicate and the problem statement's diagnosis (match stripped) is incomplete.

- **A7 Applying the three-valued verdict at BOTH kernel selection sites — `internal/resolve/resolve.go::Resolve`'s ordinary candidate loop and `::escapeOrRefuse`'s escape-edge filter — changes `flow resolve`'s disposition on exactly the five shipped `internal/resolve` oracles S5 names, and no others.**
  - **Status**: Pending
  - **Method**: Source Search + MVV Test
  - **Evidence**: To verify at Stage 6 — enumerate the `internal/resolve` oracles across BOTH selection sites and classify each against the selection predicate, as A4 did for the 0005 `next` suite. The census MUST cover every file in the package that builds a row with a `Match` block, not only the escape-shaped ones: `match_conflicted_test.go`, `guard_atoms_test.go`, `reserved_key_0008_test.go`, `escape_shape_0009_test.go`, `escape_shape_mvv_0009_test.go`, `guard_fixtures_test.go`. The five identified by inspection (3amigo iter-1 + iter-2) are the claim's floor, not its verified extent; the assumption is the exhaustiveness of "and no others", and iter-2 already moved that floor from four to five by finding a match-block oracle in a file the first census scoped out — which is the specific way this assumption fails. Both sites are confirmed: `Resolve` filters on `!view.matches(row.Match)` before `gate`, `escapeOrRefuse` on the same call, and `guard_fixtures_test.go::conformingEscapeRow` carries `Match: status=Draft`, so escape rows do carry match blocks. Note the absent-key arm reaches further than conflicted keys: `assemble` omits `recognized` when `Input.Recognized` is empty, so any oracle matching on a key the view may lack is in scope.
  - **Census under-scoped at Stage 5 (critique), 2026-08-26.** The six files named above are not the population. Nine files in `internal/resolve` build rows carrying a `Match` block; the census omits `adversarial_test.go`, `fixup_test.go`, `reserved_key_sameview_0008_test.go`, `escape_shape_adv_0009_test.go`, and `resolve_test.go`. Two further oracles go red through A8's refuted mechanism rather than through the selection predicate — `guard_fixup_0007_test.go` and `guard_adversarial_0007_test.go` — and neither is in scope of this assumption as worded. The "and no others" clause is therefore unverified over the true population, while S5, MVV step 8, and Consequences all state its conclusion as settled fact.
  - **If wrong**: the blast radius on `flow resolve` is wider than five oracles and the Consequences understate it; if it is wider in KIND (a case where refusing `guard_unevaluable` is wrong), the site-class split in C1 needs revisiting.

- **A9 Routing an `indeterminate`-match row into the `guard_unevaluable` class is compatible with `0007:C10`'s resolution-level aggregation — i.e. `flow resolve` does not refuse a whole table because one unrelated row's match key was unbound.**
  - **Status**: Pending
  - **Method**: Source Search + MVV Test
  - **REFUTED as worded at Stage 5 (critique), 2026-08-26; booked so the redesign must answer it.** `0007:C10` is normative and this RDR cites it **zero** times (`0007:C11`, the pruning clause, is cited five times): "Aggregation is resolution-level: if any surviving candidate row's guard is GuardUnevaluable, the resolution MUST refuse `guard_unevaluable`; a decided-GuardTrue sibling MUST NOT be selected while an unevaluable candidate exists. This veto is the cost the RDR accepts: one unreadable row refuses a table whose other rows decide cleanly." Shipped `gate` implements exactly that — `if len(undecidable) > 0 { return nil, &Refusal{...} }` discards the accumulated `selected`. Its partition is `{GuardFalse → pruned, GuardTrue → selected, GuardUnevaluable → veto}`, with no channel for C1's "survives into `gate` as a non-selectable row". Measured blast radius: all 22 rules in `models/rdr.toml` carry a `[rule.match.stage]` block, so a single unbound `stage` would make every row indeterminate and refuse the entire table — where today the rows drop out and an escape row rescues. Verify by deciding, explicitly, one of: a row-scoped disposition distinct from `GuardUnevaluable`; or an accepted, stated decision that match-indeterminacy is table-fatal.
  - **If wrong**: `flow resolve` gains a whole-table refusal the Consequences never name, and the "same relationship `next` already has with `guard_unevaluable`" framing is wrong — that relationship is a table-wide veto, not a per-row one.

- **A10 `JC1`'s one-way-coupling disposition holds: C1 cannot turn a 0010 acceptance into a refusal.**
  - **Status**: Pending
  - **Method**: Source Search + MVV Test
  - **REFUTED as worded at Stage 5 (critique), 2026-08-26.** `internal/resolve/resolve.go::Resolve` returns on `if blocked != nil { return refuse(in, *blocked), nil }` — **above** the `switch len(selected)` that is the only path to `escapeOrRefuse`. A `guard_unevaluable` refusal is therefore unreachable by any escape row. `0010:A9` makes the decision table's "otherwise" row an escape row rescuing `no_match`, so under C1 a partially-supplied tag set on a 0010 table stops firing its catch-all and refuses instead. JC1 states the opposite ("strictly widens what a zero-owned decision table can report and cannot turn a 0010 acceptance into a refusal"). Re-run the joint check against the redesigned selection rule.
  - **If wrong**: 0010's headline default-row feature regresses, and its shipped MVV was written against the old behaviour.

- **A8 The `Refusal.Undecided` payload can carry match-block atoms without breaking the guard-path invariants `0007:C8` fixed: widening `guard.go::isGuardBlock`'s filter and the non-`GuardUnevaluable` payload discard admits `BlockMatch` entries while leaving every existing guard payload byte-identical.**
  - **Status**: Pending
  - **Method**: Source Search + MVV Test
  - **REFUTED at Stage 5 (critique), 2026-08-26 — the stated mechanism does not exist.** Three independent legs: (i) `internal/resolve/resolve.go::gate` calls `evaluateAtoms(row.Guard, ...)` — the GUARD field — so a match atom never reaches `isGuardBlock` at all; widening that filter admits nothing, because the exclusion happens one level up, by field selection. (ii) The payload cannot hold a match atom: `Row.Match` is `[]Tag{Key,Value}` while `UndecidedAtom` is `{Key,Block,Operator,Literal,Reason}`, so `Operator`/`Literal` must be SYNTHESIZED — and `compareUndecidedAtoms` sorts on `(Key,Block,Operator,Literal)`, making those synthesized values load-bearing on `0007:C8`'s pinned ordering. That contradicts "leaves every existing guard payload byte-identical" precisely because C1 requires match atoms to join a row's existing `Undecided` entry. (iii) `isGuardBlock`'s own comment forbids the edit, naming `BlockMatch` as the vector: sweeping such an atom in "would let it be DECIDED, and a decided-FALSE one would prune its row along with the row's owned-state obligation (D8) — reopening the masking path `0007:C11` ratifies pruning against". Were it widened anyway, a match atom over a PRESENT-and-EQUAL key would reach `seam.Evaluate`, whose `default` arm returns `GuardUnevaluable`/`uncomparable`, so a correctly-matching row would report a spurious undecided atom. Two shipped oracles go red, neither named in S5: `guard_fixup_0007_test.go:115` ("a match atom is never handed to the seam") and `guard_adversarial_0007_test.go:180` (the `"§D12 match block"` arm).
  - **Also refuted — the carrier itself.** `Refusal.Undecided` is `[]UndecidedRow`, not `[]UndecidedAtom`, and `evaluateAtoms` returns a `nil` payload early for any row whose GUARD verdict is decided. So a guard-free row with an indeterminate match produces **no `UndecidedRow` at all**, and its match facts have no carrier — exactly the row class this RDR exists to diagnose. C1's "no new field on `Plan` is required" does not follow from it.
  - **Prior evidence (superseded)**: To verify at Stage 6. The two sites are confirmed by inspection: `internal/resolve/guard.go` `if !isGuardBlock(atom.Block)` skips match atoms when collecting payload entries, and `if verdict != GuardUnevaluable { return verdict, nil }` discards payloads for decided rows.
  - **If wrong**: match facts need a separate `Refusal` field rather than sharing `Undecided`, and C1's carrier clause plus `D-undecided-reporting-shape` change with it.

## Proposed Solution

### Approach

`flow next` keeps asking the kernel one question per row, but stops taking the row's match pattern out of the question — and the kernel stops answering that question in two values. The match seam becomes three-valued, the way every other seam in this system already is: a match atom is `match` (its key is present and its value equal), `no-match` (present and unequal), or `indeterminate` (its key is absent from the view, or present but conflicted so no single value can be compared). Only `no-match` excludes a row from REPORTING. An `indeterminate` atom leaves the row a candidate of `flow next` and its key rides the candidate's `unknown` list carrying the reason — `absent` or `uncomparable` — that `0007:C8` already spells for the guard seam. At the kernel's SELECTION sites the same verdict means the opposite thing: `flow resolve` MUST NOT commit a plan (and its `Writes`) on a precondition it could not decide, so an `indeterminate` match refuses `guard_unevaluable` — the non-escapable class `0007:C8` reserves for exactly this — rather than selecting. Reporting shows the caller an undecided row; committing on one would be the masking `0007:C8` forbids. The verb still never chooses among the rows that survive, never turns a gate deny into a refusal, and still returns the model's full declared `outcomes` alphabet: it narrows the **candidate list** to rows the supplied state can actually take, which is what "for the supplied state" meant to the caller.

`--all` restores 0005's enumeration — the match pattern ignored entirely — for the caller who wants the guard-conditioned alphabet rather than the state-conditioned one. The two modes differ only in whether match participates; reader narrowing, gate opt-in, the payload shape, and exit semantics are the same in both.

The change lands in the **kernel**, not as a CLI probe reshaping. Three facts force that placement. (i) The CLI cannot implement the rule: conflictedness is `internal/resolve`'s `TagSet.conflicted`, reachable only through the unexported `conflicting()`, and `internal/cli/flow_next.go::assembledView` is a last-write-wins `map[string]string` with no conflicted concept at all (zero occurrences of `conflicted` in `internal/cli`) — a CLI-side presence test would retain a conflicted atom and the kernel would exclude the row as `no_match` with nothing in `unknown`, which is the masking `P1` forbids. (ii) `resolve.go::TagSet.matches` folds absent, conflicted, and unequal into one `bool`, so the distinction cannot be reported from outside it; JDR 0001 §D12's precedent is to widen a kernel vocabulary rather than fork a second one at the CLI. (iii) `0007:REQ-78` recorded the match seam's two-valuedness as a scope boundary on that build — "NOT closed by this RDR" — not as a settled exception; this RDR closes it. ⇒ a three-valued match result in `internal/resolve`, the `unknown` payload field on the candidate, a flag, and help text; no new refusal kind (`0007:REQ-79` pins the set at five).

Sibling-path check (step 5): searched `internal/cli` and `internal/resolve` for an existing "absent key ⇒ undecided" site at the match seam — none exists. The nearest is the GUARD seam's `internal/resolve/guard.go::evaluateAtom`, which already returns exactly this triple (`GuardUnevaluable` with `ReasonAbsent` / `ReasonUncomparable`), and whose reason vocabulary this RDR reuses rather than mints.

### Technical Design

Data flow, per ordinary row: `table.Row` → `Row.KernelRow()` (match/guard split by block) → probe in `excluded` (escape stripped; `Recognized = row.Outcome`; match carried in full by default, omitted under `--all`) → `resolve.Resolve` over the one-row table → the kernel evaluates each match atom three-valued: `no-match` ⇒ the probe refuses `no_match`; `indeterminate` ⇒ the probe refuses `guard_unevaluable` carrying the undecided match atoms in `Refusal.Undecided`; otherwise ⇒ a `Plan` → `no_match` ⇒ excluded, `guard_unevaluable` / `owned_state_unavailable` / `Plan` ⇒ candidate → `summarize` reads the kernel's indeterminate match atoms off the refusal payload and folds them in beside missing owned keys and un-run gate ids → gates under `--evaluate-gates` for reported candidates only → payload. The alphabet `outcomes` is copied from the model in both modes (0005's DEV-4 reading stands: the alphabet is what may be requested; the candidates are what varies with state).

Selection semantics stay wholly in the kernel: the CLI hands over the row and reads the answer, deciding neither presence nor equality. Match-block `in` atoms are already expanded into per-member `eq` rows at normalization (`0002:C13`) and set-kinded literals canonicalized by `internal/table/model.go::seamValue`, so every match atom reaching the kernel is an independent single-key equality — the three-valued rule is per-atom with no cross-atom interaction, and a row whose match pattern is empty matches unconditionally as it does today.

The kernel change is confined to the match seam: `TagSet.matches`' `bool` becomes a three-valued verdict plus the atoms that could not be decided, computed from the data `TagSet` already carries (`!ok` ⇒ `absent`; `tv.conflicted` ⇒ `uncomparable`; `tv.value != w.Value` ⇒ `no-match`). Both of the kernel's match call sites — `Resolve`'s ordinary candidate loop and `escapeOrRefuse`'s escape-edge selection (`internal/resolve/resolve.go::escapeOrRefuse` filters on `!view.matches(row.Match)` too, and shipped escape rows do carry match blocks) — are SELECTION sites, so at both an `indeterminate` row is excluded from selection and refuses `guard_unevaluable` with its undecided atoms named. `flow resolve`'s behaviour is unchanged in every case where the CLI supplies a complete, unconflicted view; where it does not, an undecidable match atom now surfaces as a named non-escapable fact instead of an escapable `no_match`, which is the same correction `0007:C8` made on the guard path. The reporting relaxation ("`indeterminate` does not exclude") is `flow next`'s alone, and `flow next` reaches it by reading the refusal's payload, not by selecting the row.

#### Normative Contracts

**C1**

```normative
The kernel MUST evaluate a match atom three-valued against the assembled
view: `match` when the atom's key is present, unconflicted, and its value
equal; `no-match` when the key is present, unconflicted, and its value
unequal; `indeterminate` when the key is ABSENT from the view, or PRESENT
but conflicted so the view carries no single value to compare. A row's
match verdict is `no-match` if any atom is `no-match`, otherwise
`indeterminate` if any atom is `indeterminate`, otherwise `match`. Only
`no-match` excludes unconditionally. `indeterminate` MUST NOT be folded
into the escapable `no_match` refusal; what it DOES exclude from is
site-dependent, and the next paragraph is normative on that split. Each undecided atom MUST be
named by its key and a reason drawn from the closed set `0007:C8` already
owns — `absent` (the key was not in the view) or `uncomparable` (the key
was present and its value was not compared to a verdict). This mints no
sixth `RefusalKind` (`0007:REQ-79`); it closes the match-seam gap
`0007:REQ-78` deferred.

The VERDICT is the seam's and MUST be computed identically wherever a match
pattern is evaluated. What a site DOES with `indeterminate` is the site's,
and the two kinds of site differ: a REPORTING site (`flow next`) MUST NOT
exclude on it, because the caller is entitled to see the row and the reason;
a SELECTION site (`flow resolve`, which commits a transition and its
`Writes`) MUST NOT select on it, because an undecidable precondition is
exactly what `0007:C8` forbids masking — "missing artifact state MUST NOT be
maskable behind an escapable refusal class" — and what JDR 0001's P2 settles
for the kernel generally: the kernel refuses rather than guesses, so
"Indeterminate MUST be a refusal-class result, not a false allow and not a
false deny" (`0004:C9`). Selecting a plan on an undecidable match would be
the false allow. A build that computes the
verdict two ways has forked the seam; a build that treats reporting and
committing alike has masked an undecidable precondition behind a write.

At the kernel's two selection sites the rule is therefore: an ordinary row
whose match verdict is `indeterminate` MUST NOT be selected as a plan, and
MUST NOT count toward the `len(selected)` that yields `ambiguous_match`
against a row that genuinely matched. It instead refuses
`guard_unevaluable`, the non-escapable class `0007:C8` already reserves for
an undecidable predicate, carrying its undecided atoms in
`Refusal.Undecided` — the same payload, now admitting `BlockMatch` entries
(this widens `0007:C8`'s payload to the match block; it mints no sixth
`RefusalKind`, `0007:REQ-79`).

The two verdicts are NOT applied at the same place, and the split is forced
by `0007:C11`. A `no-match` row MUST be pruned where match is filtered
today, before `gate` — pruning is safe there exactly because, in
`0007:C11`'s words, exclusion "can only arise from decided atoms", and
`no-match` is decided (the key is present and unequal). An `indeterminate`
row MUST NOT be pruned there: it is by definition undecided, so dropping it
upstream would take its owned-state obligation with it (D8) and reopen the
masking path `0007:C11` ratifies pruning against. An `indeterminate` row
therefore SURVIVES into `gate` as a non-selectable row, is carried
alongside the guard-undecidable rows, and its match atoms join that row's
`Undecided` entry. One consequence is normative: because
`internal/resolve/resolve.go::gate` returns `owned_state_unavailable`
before it collects the `guard_unevaluable` payload, a row that is BOTH
indeterminate-on-match AND missing an owned key refuses
`owned_state_unavailable` — the existing precedence is preserved, not
reordered, and `MissingOwned` carries the actionable fact. The match facts
for such a row are then NOT on the refusal, which C1's reporting clause
accounts for below. An escape edge whose own match verdict is
`indeterminate` is likewise NOT a viable rescuer: rescuing through an
undecidable match is the "operator rescues the positional plan through the
kernel's most permissive path" defect
`internal/resolve/match_conflicted_test.go::TestAdv0007_3_ConflictedKeyIsNotEscapableByNamingIt`
was written to close, and this RDR does not reopen it. Only `no-match` and
`indeterminate` disqualify an escape edge; only `no-match` and
`indeterminate` disqualify an ordinary row from selection. The distinction
this RDR adds at the selection sites is not WHETHER an undecidable match
excludes — it is that the exclusion is now REPORTED (a named
`guard_unevaluable` fact) instead of silently folded into an escapable
`no_match`.

flow next MUST report as a candidate exactly each non-escape row of the
requested model whose match verdict is not `no-match` and whose guard the
kernel does not decide false. Both verdicts MUST be the kernel's, asked
over a one-row probe table carrying the row's match atoms, no escape list,
and `Recognized` bound to the row's own outcome — the binding that keeps
the probe modelling its own outcome, so `unmodeled_outcome` is unreachable.
Because an indeterminate match now refuses `guard_unevaluable` at the
selection site, the probe's disposition for such a row is that refusal, and
its `Refusal.Undecided` payload is the CARRIER: the CLI reads the match
facts off it exactly as it reads guard facts today, and no new field on
`Plan` is required. A row whose probe returns a `Plan` has a fully decided
match by construction, so it contributes no match facts.

One disposition carries only part of the row's facts. When the probe
refuses `owned_state_unavailable` — which `gate` returns BEFORE collecting
the `guard_unevaluable` payload — the refusal carries `MissingOwned` and an
empty `Undecided`, so a row that is both indeterminate-on-match and
missing an owned key yields no match facts on the wire. For that row the
`unknown` list MUST carry the missing owned keys from `MissingOwned` plus
the row's match-block atoms over keys ABSENT from the assembled view,
which the CLI derives from its own `Row.Atoms` walk (A2) without needing
the kernel. The one fact class genuinely unavailable in this case is a
CONFLICTED match key on a row that is also missing owned state, because
conflictedness is kernel-only (A3): such a key MUST be reported with
reason `absent` rather than omitted, since the view carries no value the
CLI can distinguish, and the missing owned key is the actionable remedy
the caller acts on first. This is a reporting-fidelity limit of the
compound case, not a silent drop, and it is the ONLY case in which
`unknown` does not distinguish `uncomparable` from `absent`.

Every `indeterminate` match atom, every guard fact the kernel could not
decide, every owned key no invoked reader established, and (absent
--evaluate-gates) every gate id MUST appear in that candidate's `unknown`
list as a `{key, reason}` pair. The reason token is drawn from `0007:C8`'s
closed set for the first three classes. A gate id un-run because
--evaluate-gates was not passed is NOT one of them — it is neither `absent`
(the gate is declared and its id is reported) nor `uncomparable` (a gate id
is not a view key with a value) — so it carries the reason `not-evaluated`,
which this RDR mints for the gate class ONLY and which never appears on a
match or owned-key entry. `0007:C8`'s two-member set is unchanged for the
seam it governs; `unknown` is a CLI reporting surface that spans three
sources and needs the third token to stay total. Entries are deduplicated
by `{key, reason}` pair, not by key, so a key that is both an undecided
match atom and an unestablished owned key reports once per distinct reason.
`guard_unevaluable` and `owned_state_unavailable` MUST leave the row a
candidate with those facts reported. A row excluded by `no-match` is not a
candidate: it MUST NOT be reported, and under --evaluate-gates its gates
MUST NOT run. `outcomes` MUST remain the model's full declared alphabet in
every mode. flow next MUST NOT choose among the candidates it reports and
MUST NOT turn a gate deny into a refusal.
```

⇒ this is the clause that overrides `0005:C1`'s candidate reading; the exclusion mechanism (kernel probe, decided-false-only) is 0005's, extended from guard-only to match-and-guard, and the match seam is brought to the three-valued posture the gate seam (`0004:C9` `indeterminate`) and the guard seam (`0007:C6` strong-Kleene, `0007:C8` `absent`/`uncomparable`) already hold. Prior art for the seam itself: PTaCL's atomic target evaluates `1T`/`0T`/`⊥T` and "can distinguish between a non-matching value for an attribute and a missing attribute"; SAML 2.0 §2.5.1 requires `Indeterminate` where a condition "cannot be evaluated"; SQL's `col = NULL` yields UNKNOWN, not false.

**C2**

```normative
flow next MUST accept a boolean flag --all, default false. Under --all the
candidate predicate MUST be exactly the one 0005:C1 specified: every
non-escape row whose guard the kernel does not decide false, with the row's
match pattern taking no part in the verdict. A match atom is then neither
an exclusion nor an entry in `unknown`, whatever its key's presence,
achieved by the probe omitting the match pattern so the kernel never sees
those atoms. This is a DEPARTURE from shipped behaviour, not a
preservation of it: `summarize` today walks `Row.Atoms`, which carries
match-block atoms beside guard ones (A2), so a match key absent from the
view currently lands in `unresolved` even though match plays no part in
the verdict. Under --all that entry MUST disappear — reporting a fact the
mode ignores would be reporting on a predicate it does not apply — so the
--all branch MUST filter BlockMatch atoms out of `unknown`, and an oracle
MUST assert their absence. This is the one respect in which --all is not
output-identical to 0005, and it lies outside A4's predicate-scoped
classification.
Everything else 0005:C1 requires of flow next — the data minimum, gate
accessors only under --evaluate-gates and only for reported candidates, a
deny reported on its candidate with exit 0, no invented guard facts, the
narrowed invoked reader set — MUST hold identically in both modes. The
success payload shape MUST be identical in both modes: `unknown` is present
in both, as `[]` rather than omitted when it carries nothing, so one
consumer struct parses either mode. --all is part of request identity.
--all MUST NOT be accepted by flow resolve, flow read-state, or flow
set-state — satisfied by NON-REGISTRATION (registered on next only, not on
the shared registerSelectionFlags), so cobra's shipped unknown-flag path
supplies the error and exit code. This contract mints no new refusal code
for it.
```

⇒ 0005's behaviour is preserved one flag away, so a caller that wanted the enumeration loses nothing and the 0005 oracles for it move rather than die.

**C3**

```normative
The command's Short and Long help MUST state that flow next reports the
candidates the supplied state can take (match and guard both undecided or
holding), that a match or guard key that is unsupplied — or supplied more
than once with differing values — leaves a row a candidate with the key
listed under `unknown` and its reason named, and that --all reports every
row the guards do not exclude regardless of match. The 0005 flow next tests
in internal/cli/flow_next_0005_test.go MUST keep passing unchanged where
they assert the alphabet, gate handling, reader narrowing, determinism, and
non-mutation; an assertion that depended on a match-excluded row being
reported MUST be re-homed under --all, not deleted, and the file's header
comment MUST name this RDR as the source of the default.

Renaming the candidate's flat `unresolved` list to the `{key, reason}`
`unknown` list is a payload change independent of the candidate predicate,
and it breaks assertions the predicate leaves untouched. The build MUST
mechanically re-home every shipped read of the `unresolved` key to `unknown`
and its element type from string to `{key, reason}` — six TEST reads on
this build: `internal/cli/flow_next_0005_test.go` (4, at the presence check
and the gate-id assertions), `internal/cli/flow_adversarial_0005_test.go`
(1), and `internal/cli/flow_mvv_0005_test.go` (1). The count is of test
reads only; the production rename additionally touches the
`candidate.Unresolved` field and its `json:"unresolved"` tag, the field's
doc comment, the file header comment naming `unresolved`, and `summarize`'s
`slices.Contains(c.Unresolved, …)` uses — none of which are assertions and
none of which the six counts. These are mechanical re-homings, NOT the
`--all` moves above, and MUST NOT be counted against or excused by A4's
predicate-scoped BREAKS-0 result. The gate-id reads re-home to the
`not-evaluated` reason C1 names; a re-homing that has to invent a reason
token is not mechanical and is a defect.

Because no shipped 0005 fixture contains a row whose match key is present
and unequal, the 0005 suite cannot by itself distinguish this predicate
from the stripped-match one it replaces. The build MUST therefore add
fixtures that discriminate: a row whose match key is present and UNEQUAL
(excluded by default, reported under --all), a row whose match key is
ABSENT (candidate in both, reason `absent`), and a row whose match key is
CONFLICTED (candidate in both, reason `uncomparable`). The conflicted
fixture MUST be built from two invoked readers writing the same owned key
with differing values; a repeated --tag cannot reach the kernel
(flow_input.go::parseTags refuses it as flow-tag-duplicate), so a fixture
written that way tests the CLI's duplicate refusal instead of the match
verdict and MUST NOT be accepted as this fixture. A green 0005 suite MUST
NOT be cited as evidence that this contract was implemented.
```

⇒ the override is visible where a reader meets it — help text and the test file — not only in this record.

#### Mini-checks

Fired by cue at Stage 5 (grounding pass). Cues absent: round-trip / fidelity —
this RDR defines no import/export, parse/deparse, or serialize inverse.

**`authority` — source-authority census.** Cue: the match verdict moves from a
CLI-derived decision to a kernel one, with two kernel arms and a CLI reporting arm.

| input/decision | writer (canonical) | readers | call sites | sibling arms |
| --- | --- | --- | --- | --- |
| match verdict (3-valued) | **kernel** `resolve.go::TagSet.matches` | `Resolve`, `escapeOrRefuse` | `resolve.go:472` (candidates), `resolve.go:614` (escape rescue) | the VERDICT has none — C1 forbids a per-caller fork. Its USE has exactly two arms: selection sites (both call sites above) exclude on `indeterminate` and refuse `guard_unevaluable`; the reporting site (`flow next`) does not exclude |
| conflictedness | **kernel** `TagSet.merge` → `tv.conflicted` | `matches` | `resolve.go` only | `assembledView` has no conflicted concept; NOT an arm, it is why the verdict is kernel-side |
| guard verdict | **kernel** `guard.go::evaluateAtom` | `gate` | `resolve.go:518` | unchanged by this RDR |
| assembled view | kernel `assemble` (injects `recognized`) | `matches`, `gate` | `resolve.go` | CLI `flow_next.go::assembledView` (does NOT inject `recognized`; harmless — normalize lifts it into `Row.Outcome`) |
| candidate `unknown` list | **CLI** `flow_next.go::summarize` | payload consumers | `flow_next.go:168` | folds three sources: kernel indeterminate atoms (read off `Refusal.Undecided`, reasons `absent`/`uncomparable`), absent owned keys (`absent`), un-run gate ids (`not-evaluated`). Dedup key is the `{key, reason}` pair, not the key |

**`oracle` — test-discriminability.** Cue: C3 states outright that a green 0005
suite is not evidence the contract shipped.

| MVV / scenario | fails if X is wrong because Y | negative control |
| --- | --- | --- |
| MVV 2–3 / S1 | 21→3 narrowing: a stripped-match build reports 21, so the count discriminates | `--all` (MVV 4) must still report 21 |
| MVV 5 / S2 | the added row's `match.status` ≠ `draft` must be absent by default | same row present under `--all` |
| MVV 6 / S3 present-unequal | row excluded AND its gates do not run under `--evaluate-gates` | present-equal row: candidate, gates run |
| MVV 6 / S3 absent | candidate with `{key, absent}` | supplying `--tag key=<v>` moves it out of `unknown` |
| MVV 6 / S3 conflicted | candidate with `{key, uncomparable}`, NOT excluded from REPORTING — the case the pre-Stage-4 design got wrong. Reached by two readers writing the same owned key with differing values, NOT by a repeated `--tag` (`flow_input.go::parseTags` refuses that first) | a single-valued key: candidate, `unknown` empty |
| MVV 7 | `unknown` present as `[]`, not omitted | absent-key candidate: `unknown` non-empty |
| S5 escape (A7) | an `indeterminate` escape edge does NOT rescue (it is not a viable rescuer) and the refusal names its undecided atoms | a fully-matching escape edge still rescues |
| S5 census (A7) | five named oracles change disposition and no others — the enumeration is the oracle, since "suite green" is the wrong criterion | any UNLISTED `internal/resolve` oracle going red is a defect, not an expected inversion |
| S7 compound | indeterminate match + missing owned key ⇒ `owned_state_unavailable` (precedence preserved), `Undecided` empty, row still a candidate | either fault alone: `guard_unevaluable` / plain missing-owned |
| S6 | **anti-oracle**: the 0005 suite passing is NOT evidence for C1 — it cannot distinguish the predicates (A4) | S2/S3 are the discriminating oracles |

**`disposition` — input class × outcome.** Cue: C1 sorts match input classes to
reported-vs-excluded and to exit semantics.

| input class | candidate? | `unknown` entry | gates run under `--evaluate-gates` | exit |
| --- | --- | --- | --- | --- |
| match key present, equal | yes | none | yes | 0 |
| match key present, unequal | **no** (excluded) | n/a — not reported | **no** | 0 |
| match key absent from view | yes | `{key, absent}` | yes | 0 |
| match key present, conflicted | yes | `{key, uncomparable}` | yes | 0 |
| gate id, `--evaluate-gates` not passed | yes | `{gate-id, not-evaluated}` | no | 0 |
| empty match pattern | yes (matches unconditionally) | none | yes | 0 |
| only match atom is `recognized` | yes | none — lifted into `Row.Outcome` at normalize | yes | 0 |
| guard decided false | **no** (excluded) | n/a | no | 0 |
| guard undecidable | yes | guard fact named | yes | 0 |
| owned key no reader established | yes | `{key, absent}` | yes | 0 |
| **compound**: match indeterminate AND owned key missing | yes | missing owned key(s) from `MissingOwned`, plus absent match keys from the CLI's own atom walk; a CONFLICTED match key reports `{key, absent}` (the one case `uncomparable` is not recoverable — probe refuses `owned_state_unavailable`, whose `Undecided` is empty) | yes | 0 |
| any of the above, under `--all` | match takes no part; guard rules alone decide | no match entry, whatever the key's presence | per guard | 0 |

**`trace` — desk trace over the MVV.** Cue: C1, C2, C3 plus S1–S7 all bear on one
output surface (the candidate payload). Witnesses from the A6 spike.

| step | assertions in force | witness |
| --- | --- | --- |
| 1. build with C1 default + `--all` | C2 (flag exists, default false); C2 (`--all` rejected by `resolve`/`read-state`/`set-state` — by NON-REGISTRATION, so cobra's shipped unknown-flag converter at `internal/cli/root.go` supplies the error and exit code; no new typed refusal) | flag name free: no `"all"` in `internal/cli` |
| 2. seed artifact `stage=resolved` | A6 | `{"stage":"resolved","status":"draft","gate_passed":"false"}` |
| 3. `flow next … --as=json` | C1 (candidates = match-not-no-match ∧ guard-not-false); C1 (`outcomes` full alphabet) | 3 candidates `prelock`, `resolve-route-back`, `resolve-abandon`; `match.stage eq "resolved"` at `models/rdr.toml:212/222/232` |
| 4. re-run `--all` | C2 (predicate = 0005's); C1 (`outcomes` unchanged) | 21 of 22; `finalize-pass` guard-excluded on `gate_passed` — guard, not match, so it stays excluded in BOTH modes |
| 5. `flowMVVModel` ± `--all` | C3 (added discriminating row); C1(b) guard exclusions untouched | `gated-excluded` absent in both — guard-excluded |
| 6. four match classes | C1 (three-valued); `disposition` table above | conflicted row is a candidate, not an exclusion |
| 7. `unknown` as `[]` | C2 (identical payload shape both modes) | one consumer struct parses either |
| 8. `make check` + suites | C3 (6 mechanical re-homings); S5/A7 (both kernel selection sites); S6 (anti-oracle) | `TestReq78` retires at `guard_atoms_test.go:930`; the five oracles S5 enumerates (`TestReq78` retired, three in `match_conflicted_test.go`, one in `reserved_key_0008_test.go`) are re-expected, not silently red |

No CONTRADICTION row: step 4's `finalize-pass` is guard-excluded in both modes,
which is consistent — C1(b) leaves guard exclusions untouched and C2 changes only
match's participation.

#### Load-Bearing Decisions

- **Identity** — a `flow next` request is identified as `0005:D-identity` says, plus the `--all` bit; the same request in the same mode over the same model revision and artifact contents reports the same candidate set.
- **Naming** — the flag is `--all`, boolean. Precedent for exactly this meaning (widen past a conditioned default): `docker ps -a` "Show all containers (default shows just running)", `git branch -a`, `gh workflow list --all` "Include disabled workflows" (`langref/gh-cli/pkg/cmd/workflow/list/list.go`), `beads --all` "Include closed issues", `ps -A`. Rejected: `--enumerate` (describes 0005's verb, not the caller's intent, and has no precedent in the surveyed CLIs); `--ignore-match` / `--no-match` (names the mechanism rather than the result, and collides with the kernel refusal spelling); flipping the sense (`--select` / `--match` on a 0005 default) because the defect is the default and every skill would carry the flag (Alternative 2). Rejected: an enum (`--state=all|applicable`), because the axis has exactly two members — surveyed enums (`gh pr list --state open|closed|merged|all`) exist where a third value already does, and inventing one forces a value for the common case. If a third mode ever appears the migration is `gh`'s own (`--include-forks false|true|only`): promote to a `StringEnumFlag` and keep `--all` as the widest member.
- **Undecided reporting shape** — the per-candidate field is `unknown`, a list of `{key, reason}`. Rejected: keeping the flat `[]string` `unresolved`, which cannot carry the remedy. Precedent for the word and the shape: terraform/opentofu emits `after_unknown` in plan JSON as a structural sibling of `after` (`internal/command/jsonplan/plan.go`) and carries `"unknown"` as a first-class check status distinct from `"error"` (`internal/command/jsonchecks/status.go`), reserving the human-only rendering `(known after apply)` for text mode; kubectl's conditions use `True`/`False`/`Unknown` for the same "could not determine" sense. `unknown` is emitted as `[]` rather than omitted so one consumer struct parses both modes (C2). Ordering: entries are sorted by `(key, reason)`, so the merge of kernel-supplied and CLI-computed facts is stable across builds — `D-identity` requires only set equality, but an unstable order would churn any golden payload a consumer diffs. Text mode (`--as=text`) renders the same pairs as `key (reason)` on the candidate; the reason vocabulary carries the whole remedy story, so suppressing it in the human-facing mode would leave the text caller with the symptom and no diagnosis — terraform's precedent cited above splits the RENDERING, not the information.
- **Selection / predicate** — default: candidate ⇔ row match verdict ≠ `no-match` ∧ guard not decided false (C1). `--all`: candidate ⇔ guard not decided false (C2). In both modes `next` reports all survivors and chooses none; exactly-one selection and refusal remain `flow resolve`'s (`0005:D-selection-predicate`).
- **Placement** — the three-valued match verdict is the KERNEL's, not a CLI probe reshaping. Rejected: the CLI-side pre-strip (Alternative 3), because conflictedness is not visible to the CLI and a presence-only test silently excludes a conflicted row. Rejected: a separate `resolve.Candidates` entry point, because it would leave `TagSet.matches` two-valued for `flow resolve` and fork the semantics per caller — the verdict is a property of the seam, not of who asks.
- **Verdict use, by site class** — the shared verdict does NOT imply a shared disposition. A selection site (`Resolve`, `escapeOrRefuse`) commits a transition and its `Writes`, so it MUST NOT act on an undecidable precondition: `indeterminate` refuses `guard_unevaluable`, the non-escapable class `0007:C8` reserves for exactly that, on the same "MUST NOT be maskable behind an escapable refusal class" ground. A reporting site (`flow next`) commits nothing, so it reports the row and the reason instead. Rejected: `indeterminate` never excludes anywhere — that lets `flow resolve` emit an unescaped write-performing plan whose match precondition was never checked, and reopens the escapability defect `TestAdv0007_3_ConflictedKeyIsNotEscapableByNamingIt` closes. Rejected: `indeterminate` always excludes — that is Alternative 1, which empties the list for any partial-state caller. The carrier follows from this: the probe's refusal payload already reaches the CLI, so no new field on `Plan` is needed.
- **Undecided vocabulary** — reuses `0007:C8`'s closed reason set (`absent`, `uncomparable`) for every fact that comes off the SEAM, rather than minting one there. `absent` and `uncomparable` have different remedies (bind a reader or supply the tag; fix the overlapping readers), so the payload names which. The CLI's `unknown` list spans a third class the seam never sees — a gate id un-run because `--evaluate-gates` was not passed — which is neither `absent` (the id is declared and reported) nor `uncomparable` (a gate id is not a view key with a value). That class carries `not-evaluated`, minted here for the gate class only. `0007:C8`'s set stays at two members for the seam it governs; the split is that `0007:C8` types a REFUSAL payload while `unknown` is a CLI reporting surface with one more source. Rejected: overloading `absent` for un-run gates, which would tell the caller to supply something they cannot; rejected: a single opaque `unresolved` key list, which loses the remedy; rejected: a sixth `RefusalKind`, which `0007:REQ-79` forbids.

#### Illustrative Code

Illustrative — intent only.

```sh
# Default: rows the supplied state can take. At stage=resolved this lists the
# resolved→* rows, not 21 of the 22.
intrastate flow next --model models/rdr.toml --artifact rdr=./0011.md --as=json

# 0005's enumeration: every row the guards do not exclude, match ignored.
intrastate flow next --model models/rdr.toml --artifact rdr=./0011.md --all --as=json

# A decision table with no tags: every rule is a candidate; each observed
# match key is listed under `unknown`. Supplying --tag a=x narrows it.
intrastate flow next --model table.toml --as=json
```

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Per-row kernel probe | `internal/cli/flow_next.go::excluded` | strips every match atom (`probe.Match = nil`) | Extend | probe carries the row's match pattern by default; omits it under `--all` |
| Three-valued verdict + reason vocabulary | `internal/resolve/guard.go::evaluateAtom`, `::Reason`, `::UndecidedAtom` | scoped to guard blocks (`0007:REQ-78`) | Reuse | C1 reuses `absent` / `uncomparable` verbatim at the seam; mints `not-evaluated` for the CLI's un-run-gate class only |
| Undecided payload carrier | `internal/resolve/resolve.go::Refusal.Undecided`, `guard.go::isGuardBlock`, `::evaluateAtoms` | documented and populated on `guard_unevaluable` only; `isGuardBlock` filters `BlockMatch` out; non-`GuardUnevaluable` payloads discarded | **Extend** | admits `BlockMatch` entries so an indeterminate match refuses `guard_unevaluable` with its atoms named (C1, A8); this is the kernel→CLI carrier, so `Plan` gains no field |
| Match comparison | `internal/resolve/resolve.go::TagSet.matches` | folds absent, conflicted, and unequal into one `bool` | **Extend** | returns a three-valued verdict plus undecided atoms (C1) |
| Conflicted-key detection | `internal/resolve/resolve.go::TagSet.conflicting`, `::merge` | unexported; no CLI equivalent | Reuse | the reason the verdict is computed kernel-side, not in the CLI |
| Candidate reporting | `internal/cli/flow_next.go::summarize` over `assembledView` | flat `[]string`, no reason, no `Block` filter | Extend | `unknown` carries `{key, reason}` (C1), sourced from the probe's refusal payload for match facts; dedup on the pair; the `--all` branch filters `BlockMatch` atoms out (C2) |
| Match/guard split on the probe row | `internal/table/model.go::Row.KernelRow` | none | Reuse | no table change |
| Refusal kind set | `internal/resolve/resolve.go::RefusalKinds` | pinned at five (`0007:REQ-79`) | Reuse | no sixth kind; the facts ride a payload field |
| Flag registration | `internal/cli/flow.go::registerSelectionFlags` + per-verb `cmd.Flags()` | none | Reuse | `--all` registered on `next` only (C2) |
| Reader narrowing | `internal/cli/flow_exec.go::invokedReaders` | union over all rows for `next` | Reuse | unchanged in both modes |

### Decision Rationale

The caller's question is "what can I do from here", and every peer that answers it conditions the answer on the current state: stateless's `PermittedTriggers` returns only triggers whose guards are met in the current state; xstate v5 removed the unconditioned `nextEvents` and kept the conditioned `can()`; pytransitions keeps the declared form (`get_triggers`) but as a separate operator from the evaluated one (`may_*`); SCXML's `cond` gates enablement (Investigation). 0005's predicate was the declared-shape answer; in this table model the analogue of "current state" is the match block, so a state-conditioned `next` is a match-conditioned one. The mechanism is chosen to put the decision where the deciding information is. `excluded`'s own doc already states the posture for guards — "treating 'unknown' as 'false' would silently drop rows the caller is entitled to see" — and the fix generalizes it to match, which is the last two-valued seam in a system that is three-valued at the gate (`0004:C9`) and guard (`0007:C6`, `0007:C8`) seams. Alternative 1 (strict kernel match) was rejected because it empties the list for any caller that supplies partial state, including 0010's no-tag decision table (A5), and contradicts that posture at the same function. Alternative 2 (keep 0005's default, opt-in flag) was rejected because the defect is the default and the instance read says the conditioned form is the primary one. Alternative 3 (a CLI-side presence pre-strip, this RDR's own recommendation until Stage 4) was rejected on evidence: conflictedness lives behind the kernel's unexported `TagSet.conflicting` and has no representation in `assembledView`, so a presence-only test retains a conflicted atom and the row is excluded as `no_match` with nothing reported — the masking P1 forbids, reachable by a caller who repeats a `--tag` key. Placing the verdict in the kernel also fixes it once for `flow resolve` rather than per caller, closes the deferral `0007:REQ-78` recorded, and follows JDR 0001 §D12's own precedent of widening a kernel vocabulary rather than forking a second one at the CLI. The cost is a larger blast radius than the Stage-2 estimate — recorded in the `large` Profile — bought with the seam's containment and P7's pre-release licence: `TagSet` and `matches` are both unexported, so the signature change cannot leave `internal/resolve`, and `resolve.Resolve`, the entry point whose behaviour changes, has three non-test call sites — `internal/cli/flow_next.go` (the probe) and two in `internal/cli/flow_resolve.go`. (The package as a whole is imported by 16 non-test files across `internal/{accessor,cli,graphlint,guard,table}`; they consume its types, not the match verdict.)

Premortem (paragraph): this shipped, and a skill driving a model whose rows match on an owned key its reader does not serve (a typo'd `keys` list) saw every row listed as a candidate with the key under `unknown` — the same wall-of-candidates symptom as before, now with a hint. A second failure: a caller who scripted against 0005's list found candidates missing and no flag in the error, because there is no error — the list is just shorter. Both are answered by the design rather than forcing a switch: the first is the intended three-valued behaviour and the `unknown` entry, carrying reason `absent`, is the diagnostic (the alternative, silently dropping the row, is the worse failure); the second is the override's cost, mitigated by C3's help text, the test-file header, and `--all`. Neither shows a case the chosen predicate cannot answer.
Premortem: survived (paragraph)
Ground-sweep: clean (24 anchors). Ledger: `docs/rdr/0011-flow-next-match-conditioned-candidates/evidence/grounding-sweep/sweep.md`.
Joint-check: fired → 0010; disposition 2026-08-26: cite-don't-restate — C1 is the sole normative home of the `flow next` predicate; `0010:A11` cites it and 0010 restates nothing. Open peers at depth 1 with Status Draft/Final: cli/0010 only (0001–0009 are `Implemented`). Shared modify-anchor: `internal/cli/flow_exec.go::invokedReaders` (Reuse/unchanged in both). Shared literals: `no_match` (0010: its "otherwise" row's escape class, `0010:A9`; here: the excluding kernel refusal in C1) and `flow next` — the coupling is real (A5). Absence arm: this RDR narrows a reported set rather than converting a refusal into an acceptance; no `Final` peer exists, and the `Implemented` predecessor 0005's reliance on the enumeration is named under `Overrides` and rides to 7.1. Bridge sub-check: n/a — neither plan retires a surface the other introduces (0010 stays out of `flow_next.go`; this RDR stays out of `resolvePayload`, `normalizeRule`, `reach`, `checkGroups`). **Re-run 2026-08-26 after Stage 4 relocated the verdict into the kernel:** the shared surface widens from `flow next` to `internal/resolve/resolve.go::TagSet.matches`, which 0010's decision-table class reaches through `Resolve`. The disposition is unchanged and the coupling stays one-way — C1 makes an undecidable match atom a named fact instead of an escapable `no_match`, which strictly widens what a zero-owned decision table can report and cannot turn a 0010 acceptance into a refusal. `0010:A11` already pre-commits to re-verification "against cli/0011's candidate predicate once 0011 locks", and A5 records the one limit it cannot discharge here: a genuine zero-owned decision table is unloadable on this build, so `0010:A11`'s "empty `required`" half is downstream of 0010 shipping, not of this contract.

## Alternatives Considered

### Alternative 1: Strict kernel match (absent key excludes)

**Description**: Hand the kernel the row's full match block and let `TagSet.matches` decide: a match key the view lacks is a non-match, so the row is excluded.

**Pros**:

- Zero CLI logic beyond removing two `nil` assignments; `next`'s candidates equal exactly the rows `resolve` could select.

**Cons**:

- Any partially supplied state — a decision table queried with no `--tag`, a machine whose match key no reader serves — yields an empty list with nothing in `unknown` to act on.
- Contradicts the "undecided is not false" posture `excluded` already takes for guards and 0007 takes for atoms.

**Reason for rejection**: it breaks 0010's no-tag `flow next` (A5) and turns the verb's most useful degraded answer into silence.

### Alternative 2: Keep 0005's default; add an opt-in `--match` / `--select` flag

**Description**: Leave the guard-conditioned enumeration as the default and let a caller request the match-conditioned list with a flag.

**Pros**:

- No override of `0005:C1`; the 0005 tests do not move.

**Cons**:

- The defect is the default: every skill invocation must carry the flag or get the alphabet.
- The instance read points the other way — xstate retired the unconditioned form, stateless never offered it.

**Reason for rejection**: fixes the symptom for callers who know about the flag and leaves the observed defect in place for everyone else.

### Alternative 3: CLI-side presence pre-strip (no kernel change)

**Description**: Leave `TagSet.matches` two-valued and reshape the probe in `internal/cli/flow_next.go::excluded` instead — carry the row's match atoms over keys the assembled view holds, drop the atoms over keys it lacks, and read `no_match` as the sole exclusion. This was this RDR's own recommendation until Stage 4.

**Pros**:

- No kernel change; the diff is one function, a flag, and help text.
- The presence lookup already exists in the CLI for `unresolved` reporting and compares no value.

**Cons**:

- **The CLI cannot implement the rule correctly.** Conflictedness lives in `internal/resolve`'s unexported `TagSet.conflicting`, and `assembledView` is a last-write-wins `map[string]string` with no conflicted concept (zero occurrences of `conflicted` in `internal/cli`). A bare presence test retains a conflicted atom; `TagSet.matches` then returns false on `tv.conflicted` and the row is excluded as `no_match` with NOTHING in the candidate's undecided list — the silent drop `P1` forbids, reached through a different door than the one A1's "If wrong" watches.
- It forks a second presence vocabulary at the CLI beside the kernel's, which JDR 0001 §D12 widened `resolve.Block` specifically to avoid.
- It leaves `0007:REQ-78`'s deferred match-seam gap open, and leaves `flow resolve` still folding an undecidable match atom into an escapable `no_match`.

**Reason for rejection**: it buys a smaller diff by putting the decision where the deciding information is not. Verified in Stage 4 (A3): the conflicted-key case is user-reachable — `internal/cli/flow_exec.go::flowRequest.runReaders` appends each invoked reader's `OwnedSnapshot()` with no dedup, so two readers writing one key with differing values conflicts it — so the failure is real, not theoretical.

### Briefly Rejected

- **Filter `outcomes` to outcomes with a surviving candidate**: 0005 already refused this reading (DEV-4 in `runFlowNext`); the alphabet is what may be requested.
- **A separate verb (`flow candidates`)**: two verbs for one question; 0005's four-verb closure stands.
- **Report match-excluded rows with an `excluded_by` field instead of dropping them**: grows the payload for the same information `--all` already yields.

## Context

### Background

Observed 2026-08-26 while driving the consumer model (rdr#tmxk): `flow next --model models/rdr.toml --artifact rdr=<state with stage=resolved>` lists 21 of the 22 rules. The cause is by design in RDR 0005: `internal/cli/flow_next.go::excluded` nils `probe.Match` / `probe.Escape` before probing, under the comment "dropping a row whose match pattern the supplied facts do not satisfy would be a selection this verb was not asked to make" — the sentence this RDR revisits. 0005 rationale `0005:A3` holds that `next` exposes the alphabet without owning guard evaluation; selection (gate-then-count, exact-one survivor) belongs to the kernel. 0005's Briefly Rejected list also refused exit-code-only output because it "cannot carry legal outcome alphabets, conditional summaries" — the conditional summary is the surface this RDR sharpens. Constraints: RDRs are never amended in content (this is a new RDR that overrides 0005's clause); intrastate stays generic — no consumer (RDR-process) knowledge in code, docs, or fixtures; `make check` must pass. Not a facet of kata `zdat` / cli/0010 (owned-state optionality is a model-class decision; this is a verb-predicate decision) — cross-cite only; 0010 leaves `flow next` to this RDR (`0010:C4`, `0010:A11`).

### Technical Environment

Go CLI (`bin/intrastate`); `internal/cli/flow_next.go` (`excluded`, `summarize`, `assembledView`, candidate/guard reporting), kernel probe API `internal/resolve/resolve.go::Resolve` that today receives rows with match stripped, `internal/cli/flow_next_0005_test.go`. Conventions in `AGENTS.md` (respond gateway, CLIError codes, SilenceUsage). Design history: `docs/rdr/0001–0010`, `docs/jdr/0001`.

## Research Findings

### Investigation

Prior art was read first, class and instance. Class (StateMachineRes corpus, three queries, ledger in `docs/rdr/0011-flow-next-match-conditioned-candidates/evidence/research/prior-art.md`): the "what next" operator exists in every peer that has a current state, and is conditioned on it. Instance: qmuntal-stateless's README "Introspection" — "a list of the triggers that can be successfully fired within the current state via the `StateMachine.PermittedTriggers` property" — implemented in `states.go::stateRepresentation.PermittedTriggers` as triggers with `len(tb.UnmetGuardConditions(...)) == 0` ⇒ state- and guard-conditioned, the conditioned form is the only form. pytransitions `transitions/core.py::Machine.get_triggers` returns triggers declared from the given states without evaluating conditions, while `Machine._can_trigger` (the `may_*` family) evaluates them ⇒ the declared-shape enumeration exists but as a distinct operator, which is what `--all` is here. xstate's `packages/core/CHANGELOG.md` records "Removed `MachineSnapshot['nextEvents']`" and `packages/core/src/State.ts::machineSnapshotCan` keeps the guard-evaluated `can(event)` ⇒ when forced to keep one, the peer kept the conditioned query. scxmlcc `doc/user-manual.md` (`cond`): "The transition is only executed if the condition evaluates to true" ⇒ enablement is event- and condition-conditioned. The state-machine peers gave no coverage for distinguishing an *absent* state key from a mismatching one — every peer assumes a current state is always present — so Stage 4 searched the access-control and query literature, where partial input is the normal case, and found the seam formalized. **PTaCL** (policy target algebra) models a *target*: a set of key/value equalities selecting which rules apply, evaluated before the rule body — structurally this codebase's match block. Its atomic target is three-valued, `⟦(n, v)⟧(q) = 1T if (n, v′) ∈ q and v = v′; ⊥T if (n, v′) ∉ q; 0T otherwise` — absence tested *before* inequality — and the paper states the point outright: "PTaCL can distinguish between a non-matching value for an attribute and a missing attribute." **SAML 2.0 Core** §2.5.1 covers the uncomparable case: where a condition "is encountered that is not understood, the status of the condition cannot be evaluated and the validity status of the assertion MUST be considered to be Indeterminate." **XACML** is the ancestor of both, contributing the `Indeterminate` vocabulary. **SQL** is the mass-adoption instance: a comparison against NULL yields UNKNOWN, "even (NULL = NULL) is UNKNOWN", under an explicit three-valued logic. ⇒ no surveyed system folds present-but-uncomparable into a silent non-match; that is the one disposition the literature does not take, and it is the one the pre-Stage-4 draft of C1 took. Ledger: `docs/rdr/0011-flow-next-match-conditioned-candidates/evidence/research/three-valued-match.md`.

### Key Discoveries

- **Documented** — `internal/cli/flow_next.go::excluded` strips `Match` and `Escape` from the probe and excludes only on `resolve.KindNoMatch`; its doc gives the "unknown is not false" rule this RDR extends to match. ⇒ the change is a probe-shape change in one function.
- **Documented** — `internal/resolve/resolve.go::TagSet.matches` collapses three distinct cases into one `false`: absent (`!ok`), conflicted (`tv.conflicted`), and unequal. ⇒ the distinction cannot be recovered outside the function, so the three-valued verdict is computed there (C1, A3).
- **Documented** — the codebase is three-valued at every other seam: gates return `allow` / `deny` / `indeterminate`, where "Indeterminate MUST be a refusal-class result, not a false allow and not a false deny" (`0004:C9`); guards combine under strong-Kleene with normative truth tables (`0007:C6`) and name each undecided atom by key and a reason from a closed set (`0007:C8`). ⇒ the match seam is the sole two-valued holdout, and C1 reuses `absent` / `uncomparable` rather than minting a vocabulary.
- **Documented** — `0007:REQ-78` records the match seam's two-valuedness as a scope boundary on that build ("This is deliberate … and NOT closed by this RDR"), and JDR 0001 §D12 waved the case through as "harmless, since match atoms are never unevaluable". ⇒ the premise is false — a conflicted key makes a match atom unevaluable in exactly `0007:C8`'s `uncomparable` sense — so this is an open deferral, not a settled exception (Alternative 3).
- **Documented** — JDR 0001 rejected folding an undecidable or distinct condition into the escapable `no_match` three times (§D2 gate-then-count, §D6 block-keyed routing, §D9 gates-after-selection), under principles P1 "missing artifact state is never masked", P2 "the kernel refuses rather than guesses", P4 "a refusal names what was missing", and P7 "pre-release, prefer the clean shape". ⇒ the same principles decide this seam the same way.
- **Documented** — `internal/cli/flow_next.go::summarize` walks every `Row.Atoms` entry, match-block atoms included, against the view. ⇒ an absent match key is already an `unresolved` entry (A2).
- **Documented** — 0005 Decision Rationale premortem: "`next` omits enough condition detail to constrain a skill". ⇒ 0005 foresaw this failure and answered it with candidate *summaries*; the summaries were not enough because the *set* was unconditioned.
- **Documented** — `0010:C4`: "`flow next`, `flow read-state`, and `flow set-state` are unchanged by this RDR"; `0010:A11` defers the decision-table `next` behaviour to this RDR's predicate. ⇒ the absent-key rule must serve a no-tag decision table (A5).
- **Documented** — prior art (Investigation). ⇒ conditioned is the default, enumeration is the secondary surface.
- **Assumed** — every 0005 `next` oracle seeds the owned key its named rows match on (A4).

## Trade-offs

### Consequences

- Positive: a skill can read the next step from the candidate list at `stage=resolved` (A6) instead of the whole table.
- Positive: no table change and no new refusal code; the kernel diff is confined to the match seam, and the reason vocabulary is `0007:C8`'s, already shipped and tested.
- Positive: `flow resolve` stops folding an undecidable match atom into an escapable `no_match`, closing `0007:REQ-78` and correcting JDR 0001 §D12's false premise — a defect neither RDR had a caller for until now.
- Positive: `absent` and `uncomparable` have different remedies (bind a reader or supply the tag; deduplicate the input), and the payload names which, so the caller is told what to do rather than only that something is unknown.
- Negative: an override of a locked 0005 clause — callers scripted against the enumeration must add `--all`. `--all` restores the candidate SET, not the payload SHAPE: the per-candidate `unresolved` (`[]string`) becomes `unknown` (`[{key, reason}]`) in BOTH modes, so an out-of-repo consumer parsing the old field must change even if it adopts `--all`. The two migrations are independent and a caller may owe both. `docs/cli-output-contract.md` pins the envelope, not per-command payload fields, so this is a documented-surface change without a contract-doc edit; C3's help text and the test-file header are where a reader meets it.
- Negative: a candidate can still be listed that `resolve` would refuse (an undecidable match key) — the same relationship `next` already has with `guard_unevaluable`.
- Negative: `flow resolve` gains a refusal where it previously produced a plan. A row whose match is undecidable (an absent or conflicted match key) previously dropped out silently and let an escape row rescue; it now refuses `guard_unevaluable` naming the undecided atoms, and `guard_unevaluable` is not escapable. Five shipped `internal/resolve` oracles change disposition as a result and are re-expected by name in S5 — this is the largest behaviour change in the RDR, and it lands on `flow resolve`, not `flow next`.
- Negative: a row with an indeterminate match no longer participates in `len(selected)`, so it can neither be selected alone nor create an `ambiguous_match` against a row that genuinely matched. That is the intended reading — an undecidable row is not evidence of ambiguity — but it means a table whose rows overlap only under undecidability resolves where a naive reading might expect a refusal.
- Neutral: `flow next` is unaffected by the selection-site rule. It strips the escape list from its probe, and it reads the refusal payload rather than selecting, so an indeterminate row is still reported as a candidate with its reason.
- Negative: the RDR grew from one contract to three and from a CLI-local change to a kernel-seam one; `TestReq78` retires and the `internal/resolve` suite gains the three-valued oracles.

### Risks and Mitigations

- **Risk**: an absent match key makes the list look like the old wall of candidates.
  **Mitigation**: the key is named under `unknown` with reason `absent` on every such candidate (C1); the fix is on the caller's side (bind the reader / supply the tag) and is visible.
- **Risk**: a 0005 oracle silently passes for the wrong reason after the default flips.
  **Mitigation**: A4's run plus C3's explicit re-homing rule; Testing Strategy S2 pairs default and `--all` on the MVV fixture.

### Failure Modes

- **Visible**: a candidate the caller expected is absent. Diagnose with `--all`: if it appears there, a match atom over a supplied key failed (compare the row's `match` block in `dump` against `owned`/`observed` in the payload).
- **Visible**: every row is a candidate with the same key under `unknown`, reason `absent` — the match key is neither owned by an invoked reader nor supplied; bind or supply it.
- **Visible**: a key carrying two differing values within one provenance leaves every affected row a candidate carrying `{key, uncomparable}` under `flow next`, and refuses `guard_unevaluable` under `flow resolve`. The reachable producer is two invoked READERS writing the same owned key with differing values — `flowRequest.runReaders` appends each reader's `OwnedSnapshot()` with no dedup. It is NOT a repeated `--tag`: `internal/cli/flow_input.go::parseTags` refuses that first with `flow-tag-duplicate`. The remedy is to fix the overlapping readers, and the reason names it. (Before this RDR's Stage-4 revision the design excluded such a row silently; that failure mode is what C1's three-valued verdict removes.)
- **Recovery**: none needed — `next` is effect-free in both modes; re-run with the corrected inputs.

## Implementation Plan

### Prerequisites

- [x] All Critical Assumptions verified at Stage 4 (A1–A6). A3 relocated the verdict into the kernel; A4 found the 0005 suite cannot discriminate the new predicate, which C3 now answers with named fixtures.

### Minimum Viable Validation

1. Build with the C1 default and the `--all` flag.
2. Seed an artifact for `models/rdr.toml` whose owned `stage` is `resolved`; run `flow next --model models/rdr.toml --artifact rdr=<it> --as=json`.
3. Assert: `candidates[]` is exactly `prelock`, `resolve-route-back`, `resolve-abandon` (3 rows, one per outcome), and no rule whose `match.stage` names another value; `outcomes[]` is still the model's full alphabet.
4. Re-run with `--all`; assert `candidates[]` has the 21 rows the baseline reported (the 22nd, `finalize-pass`, is guard-excluded on `gate_passed`, not match-excluded) and the same `outcomes[]`.
5. Run over the 0005 MVV fixture (`flowMVVModel`, `status=draft`) with and without `--all`; assert the `gated-excluded` row is absent in both (guard-excluded), and that the row C3 adds — whose `match.status` names a value other than `draft` — is absent by default and present under `--all`.
6. Run over a model exercising all four match cases: present-and-equal (candidate, nothing in `unknown`), present-and-unequal (excluded; gates do not run under `--evaluate-gates`), absent (candidate, `{key, absent}`), and a key carrying two differing values from two invoked readers (candidate, `{key, uncomparable}` — NOT excluded; a repeated `--tag` cannot produce this, `parseTags` refuses it first). Then add `--tag key=<value>` and assert the list narrows to the rows whose match holds.
6b. Assert the selection-site half in `internal/resolve`: the same undecidable match atom refuses `guard_unevaluable` (not `no_match`, not a plan) and names the atom in `Refusal.Undecided`; an escape edge whose own match is indeterminate does not rescue.
6c. Assert the compound row (S7): indeterminate match AND a missing owned key refuses `owned_state_unavailable` with an empty `Undecided`; the CLI still reports the row as a candidate carrying the missing owned key plus its view-absent match keys, and a conflicted match key on that row reports `{key, absent}`.
7. Assert `unknown` is present as `[]` rather than omitted on a candidate with nothing undecided, in both modes, so one consumer struct parses either.
8. `make check` passes; the `internal/resolve` suite passes with `TestReq78` retired and the five oracles re-expected per S5 (`TestReq78` retired, three in `match_conflicted_test.go`, one in `reserved_key_0008_test.go`) — no OTHER oracle in that package changes disposition; the 0005 `next` suite passes with only C3's six mechanical `unresolved`→`unknown` re-homings, and the file header names this RDR.

End-state: `flow next` answers "what is legal from here" with the rows the supplied state can take, `--all` yields 0005's enumeration, and every 0005 obligation other than the candidate predicate is unchanged.

### Phase 1: Three-Valued Match Seam (kernel)

Intent: give `internal/resolve`'s match comparison a three-valued verdict plus the atoms it could not decide, computed from what `TagSet` already carries (`!ok` ⇒ `absent`, `tv.conflicted` ⇒ `uncomparable`, unequal ⇒ `no-match`), reusing `0007:C8`'s `Reason` vocabulary. `Resolve`'s candidate loop and `escapeOrRefuse`'s escape-edge filter are both SELECTION sites: neither may SELECT an `indeterminate` row. The two verdicts land at different places, per `0007:C11`: `no-match` prunes in the match filter as today (safe — a decided atom), while an `indeterminate` row survives that filter into `gate` as a non-selectable row so its owned-state obligation is not dropped with it (D8), where it joins the guard-undecidable set. `owned_state_unavailable` keeps its existing precedence over the undecidable payload. An `indeterminate` exclusion refuses `guard_unevaluable` carrying the undecided match atoms in `Refusal.Undecided` (widening that payload to admit `BlockMatch` entries — `guard.go::isGuardBlock` currently filters them out, and `evaluateAtoms` discards payloads for non-`GuardUnevaluable` verdicts, so both need the match-block path). Mint no sixth `RefusalKind`. Retire `TestReq78` and re-expect the four remaining oracles S5 names (three in `match_conflicted_test.go`, one in `reserved_key_0008_test.go`) — five in total with the retirement.

### Phase 2: Probe and Reporting (CLI)

Intent: stop stripping `probe.Match` in `excluded` so the kernel sees the row's match pattern; carry the kernel's undecided match atoms into the candidate's `unknown` as `{key, reason}` pairs beside the existing absent-owned-key and un-run-gate facts; emit `unknown` as `[]` rather than omitting it.

### Phase 3: `--all` and Help

Intent: register `--all` on `next` only, route it to a probe with the match pattern omitted, and rewrite Short/Long help per C3.

### Phase 4: Fixtures and Docs

Intent: add C3's discriminating fixtures (present-unequal, absent, conflicted), the default/`--all` pair, and the `recognized`-only and empty-match rows; re-home the 6 shipped reads of the `unresolved` key to `unknown` (C3); update the 0005 file header; document `flow next`'s payload — including `unknown` and the `--all` invocation — in `docs/cli-output-contract.md`, which currently shows only the invocation grammar.

## Validation

### Testing Strategy

1. **Scenario**: S1 — `models/rdr.toml` at owned `stage=resolved`, default mode (MVV 2–3). Backed by the A6 spike, which ran the baseline live and prototyped the predicate.
   **Expected**: `candidates[]` is exactly `prelock`, `resolve-route-back`, `resolve-abandon` — the three rules whose `match.stage` admits `resolved`, one per declared outcome — down from the 21 the baseline reports; `outcomes[]` is the full three-member alphabet. Normative fixture: `evidence/spikes/a6-rdr-toml-narrowing.md`.
2. **Scenario**: S2 — `flowMVVModel` (`status=draft`) with and without `--all` (MVV 4–5), plus the discriminating row C3 requires.
   **Expected**: `gated-excluded` absent in both (guard-excluded, untouched by C1); the added row whose `match.status` names a value other than `draft` is absent by default and present under `--all`; `unknown` present as `[]` in both modes on a fully-resolved candidate; under `--all`, a row with an absent match key reports NO match entry in `unknown` (the C2 departure — it does today), while the same row by default reports `{key, absent}`; `--all` rejected by `resolve` / `read-state` / `set-state` via cobra's unknown-flag error, asserted on that error class and exit code rather than a `flow-*` code.
3. **Scenario**: S3 — the three-valued match verdict, one row per value (MVV 6). Backed by the A5 spike's fixture shape.
   **Expected**, at the REPORTING site (`flow next`): a match key PRESENT and equal ⇒ candidate, key absent from `unknown`; PRESENT and unequal ⇒ row excluded and, under `--evaluate-gates`, its gates do not run; ABSENT ⇒ candidate with `{key, absent}` in `unknown`; carrying two differing values within one provenance ⇒ candidate with `{key, uncomparable}` in `unknown`, NOT excluded. `--tag key=<v>` narrows the list to the rows whose match holds. The fourth row is the one the pre-Stage-4 design got wrong and MUST be asserted.

   The conflicted fixture MUST be built from two invoked readers writing the same owned key with differing values (`flow_exec.go::runReaders` appends each `OwnedSnapshot()` with no dedup). A repeated `--tag` CANNOT produce it — `flow_input.go::parseTags` refuses with `flow-tag-duplicate` before the kernel is reached — so a fixture written that way asserts the wrong thing and passes for the wrong reason. The kernel-level arm of the same case is asserted directly in `internal/resolve` over a `TagSet` built in-package (S5), which is where the `guard_unevaluable` selection disposition is checked.
4. **Scenario**: S4 — a row whose only match atom is `recognized`, and a row with an empty match pattern (A5's fixture B).
   **Expected**: both are candidates with `unknown` empty of match facts — `recognized` is lifted into `Row.Outcome` at normalization and is never an undecided fact, and an empty match pattern matches unconditionally.
5. **Scenario**: S5 — `flow resolve` over the same seam change (C1's kernel half).
   **Expected**: every case with a complete, unconflicted view resolves exactly as before; an undecidable match atom refuses `guard_unevaluable` naming the atom, rather than dropping the row into an escapable `no_match`.

   The `internal/resolve` suite does NOT pass unchanged, and "green" is the wrong pass criterion. Five shipped oracles change disposition under C1, and the build MUST re-expect each by name — any OTHER oracle going red is a defect, which is the discrimination this scenario provides:
   - `guard_atoms_test.go::TestReq78_MatchPatternStillFoldsAbsenceIntoNonMatch` — **retired**. It asserts an absent match key yields the escapable `no_match`; C1 closes `0007:REQ-78`.
   - `match_conflicted_test.go::TestAdv0007_3_EscapeNotNamingTheConflictedKeyStillRescues` — **re-expected**. Its ordinary row `A` is guard-free with `Match: [status=Draft, gate=closed]` over a conflicted `gate`, so under C1 row `A` refuses `guard_unevaluable` instead of dropping out, and the clean escape is never reached (`guard_unevaluable` is not escapable). The new expectation is that refusal, naming `gate` `uncomparable`.
   - `match_conflicted_test.go::TestAdv0007_3_DuplicateKeyMakesRowSelectionAFunctionOfSlicePosition` — **re-expected**: refusal kind moves from `no_match` to `guard_unevaluable`. Its point (order-independence) is preserved and MUST still be asserted.
   - `match_conflicted_test.go::TestAdv0007_3_ConflictedKeyIsNotEscapableByNamingIt` — **kept, strengthened**. Its stated contract — an operator must not rescue a positional plan through the kernel's most permissive path — is what C1's selection rule preserves; the refusal kind moves to `guard_unevaluable`, which is non-escapable, so the property holds a fortiori.

   - `reserved_key_0008_test.go::TestReq3_AbsentOutcomeYieldsAViewWithNoRecognizedKey` — **re-expected**: refusal kind moves from `no_match` to `guard_unevaluable`. With `Input.Recognized == ""`, `resolve.go::assemble` skips the `recognized` binding entirely, so the row's `recognized` match atom names a key ABSENT from the view — `indeterminate` under C1, not `no-match`. Its actual subject (the kernel must not inject an empty binding, so the row is not selected) is PRESERVED and MUST still be asserted: `got.Plan` stays nil. Only the refusal kind changes.

   `TestAdv0007_3_MatchNarrowingsHold` and every other `internal/resolve` oracle MUST keep passing unchanged.
6. **Scenario**: S6 — the 0005 `next` suite after the default flips (A4, C3).
   **Expected**: the `flow next` oracles pass, none re-homed under `--all` (A4's MOVES and BREAKS classes are both empty against the predicate); the 6 test reads of the `unresolved` key are mechanically re-homed to `unknown` per C3 and pass on the new shape; file header names this RDR. A green 0005 suite is NOT evidence C1 shipped — S2/S3 are.

   The "28 oracles" figure A4 reports is a classification count, not a checkable pass criterion: the enumeration behind it was never written down, so the build cannot verify "28" or "none re-homed" against a list. The build MUST either (a) emit the enumeration as an artifact when it re-runs the classification, or (b) assert the weaker checkable form — `go test ./internal/cli` green except C3's six named re-homings, with no test re-homed under `--all`. Treat the number as provenance for A4's verdict, not as an oracle.

7. **Scenario**: S7 — the compound row: match verdict `indeterminate` AND an owned key no invoked reader established (C1's precedence clause).
   **Expected**: at the kernel, the probe refuses `owned_state_unavailable`, not `guard_unevaluable` — `gate` returns the missing-owned refusal before it collects the undecidable payload, and this RDR preserves that precedence rather than reordering it; `Refusal.Undecided` is empty on that disposition. At the CLI, the row is still a candidate and its `unknown` carries the missing owned key(s) from `MissingOwned` plus its view-absent match keys from the CLI's own atom walk. A match key that is present-but-CONFLICTED on such a row reports `{key, absent}` — the single case where `uncomparable` is not recoverable, asserted explicitly so the limit is pinned rather than discovered. Also assert the pruning split: a `no-match` row is excluded before `gate` and contributes no owned-state obligation, while an `indeterminate` row survives into `gate` and keeps its `RequiresOwned` obligation (`0007:C11`'s condition — pruning is safe only on decided atoms).

Done: S1–S7 green and `make check` passes.

## Finalization Gate

> Complete each item with a written response in
> `{ARTIFACT_DIR}/gate.md` before marking this RDR as
> **Final**. Written responses prevent rubber-stamping
> and produce a review record.
>
> First run the mechanical pre-sweep
> (`prompts/gate/tooling-pass.md`): TEMPLATE section
> coverage, Method-label vocabulary, `Source Search`
> self-reference, `Docs Only` on load-bearing claims. It
> catches what the review rounds disturbed; resolve any
> BLOCK before the written responses.
>
> At lock, replace this section's body with the
> one-line pointer to gate.md — responses are never
> inlined. The sub-sections below spec gate.md's
> content.

### Contradiction Check

[State any conflicts between Research Findings and
the Proposed Solution. If none exist, state
"No contradictions found between research findings,
design principles, and proposed solution."]

### Assumption Verification

[Confirm every Critical Assumption Evidence Record
is internally consistent: Status, Method, and
Evidence agree, and "If wrong" is non-empty. List
any record whose Method is `Docs Only` (these block
lock unless paired with a Spike or Source Search
plan) and any that remain `Pending` or `Unverified`
with a plan to verify before implementation begins.
Confirm no `Verified` stamp is self-referential or
proves only an adjacent claim, and that each cited
`path::Symbol` resolves on `main`. **Status
consistency:** no assumption marked `Pending` or
`Unverified` may have settled-fact prose elsewhere in
the RDR depending on it.]

### Scope Verification

[Confirm the Minimum Viable Validation is in scope
and will be executed during implementation, not
deferred. State the specific test or proof.]

### Cross-Cutting Concerns

[List only concerns that apply to this RDR. For each,
state either how this RDR addresses it, or which peer
RDR owns the project-wide policy this RDR conforms
to. Omit (rather than N/A-bullet) anything that does
not apply.]

Candidate concerns (include only those that apply):
versioning · build tool compatibility · licensing ·
deployment model · IDE compatibility · incremental
adoption · secret/credential lifecycle · memory
management · concurrency model · character encoding ·
canonical-form / determinism (see note below).

If this RDR claims byte-identical output,
content-addressed identity, or replay-stable hashes,
also confirm: hash function + library, pre-image
byte layout, primitive encodings, map iteration order,
whitespace policy, case folding, empty/null/absent
distinguishability, and a version marker for future
evolution.

### Proportionality

[Is the document right-sized for the change? Flag
any sections that should be trimmed before locking.
The split test is **contract count, not word count**:
confirm this RDR is the sole author of at most one
independent load-bearing contract (per the Normative
Contracts split signal). If it owns more than one
seam, flag it for splitting rather than locking the
seams together.

Re-validate the **Profile** Metadata field against the
contracts you just counted: confirm the value Resolve
wrote still matches (one contract + no user-facing
surface → `small`; etc. per the applicability matrix).
If the lenses that actually ran disagree with the
Profile (e.g. Profile says `small` but the change locks
a contract that warranted `mid`+ lenses, or the lenses
were skipped on a wrong `small`), correct the field and
do not lock until the missing lenses have run. This is
the latch's backstop — a wrong Profile cannot route
past the lens battery undetected. A `Transient`-marked
contract with a named deleting sibling and schedule is a
recorded lifespan disposition, not an under-sized
Profile — do not count it when re-deriving. Also confirm form:
value + one clause naming the contract(s); strip any
matrix/provenance prose left from the template or Seed
(it belongs in the template comment, not the instance).]

## References

- RDR 0005 `0005:C1` (the overridden clause), `0005:A3`, `0005:A7`, `0005:D-selection-predicate`, Decision Rationale premortem; RDR 0010 `0010:C4`, `0010:A11`; RDR 0007 `0007:C8`; JDR 0001 §D2.
- Source reviewed: `internal/cli/flow_next.go`, `internal/cli/flow_exec.go`, `internal/cli/flow.go`, `internal/table/model.go`, `internal/resolve/resolve.go`, `internal/cli/flow_next_0005_test.go`, `internal/cli/flow_fixtures_0005_test.go`, `docs/cli-output-contract.md`.
- Prior art, state machines (sibling checkouts under `../state-machines`): `repos/qmuntal-stateless` (`README.md` Introspection; `statemachine.go`, `states.go`), `study/transitions` (`transitions/core.py`), `repos/xstate` (`packages/core/CHANGELOG.md`, `packages/core/src/State.ts`), `repos/scxmlcc` (`doc/user-manual.md`). Ledger: `docs/rdr/0011-flow-next-match-conditioned-candidates/evidence/research/prior-art.md`.
- Prior art, three-valued selection seams: PTaCL via Griesmayer & Morisset "ATRAP" §2.1–2.2 and Table 1; SAML 2.0 Core §2.5.1; Cedar (Cutler et al. 2024) §Related Work; Elmasri & Navathe *Fundamentals of Database Systems* §4.1, §5.1 Table 5.1; Celko *SQL Database Programmers Handbook* §"The Null of It All". Ledger: `docs/rdr/0011-flow-next-match-conditioned-candidates/evidence/research/three-valued-match.md`.
- Prior art, CLI shape: `langref/gh-cli` (`pkg/cmd/workflow/list/list.go`, `pkg/cmd/pr/list/list.go`, `pkg/cmd/search/repos/repos.go`), `langref/helm` (`pkg/cmd/list.go`), `langref/beads` (`cmd/bd/query.go`), `langref/opentofu` (`internal/command/jsonplan/plan.go`, `internal/command/jsonchecks/status.go`); `docker ps`, `git branch`, `ps` man pages.
- Internal precedent: JDR 0001 principles P1/P2/P4/P7 and §D2, §D6, §D9, §D12; `0004:C9` (`indeterminate`); `0007:C6` (strong-Kleene), `0007:C8` (`absent`/`uncomparable`), `0007:REQ-78` (the deferred match seam), `0007:REQ-79` (five refusal kinds).
- Spike evidence: `evidence/spikes/a6-rdr-toml-narrowing.md`, `evidence/spikes/a5-decision-table-no-tag.md`.
- kata `1mv1`; rdr#thsc; rdr#tmxk.
