# Recommendation 0011: flow next selects by match; --all enumerates the alphabet

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-08-26
- **Status**: Draft [re-verify A12, A15]
- **Type**: Feature
- **Profile**: large — one independent contract, the `flow next` candidate predicate (C1), locking the `unknown` payload shape and its reason vocabulary; C2 (`--all`) is that predicate's second surface and C3 pins its wording and fixtures, so neither is separately held
- **Priority**: High
- **Related Issues**: kata `1mv1` (defect tracker); umbrella rdr#thsc; consumer model rdr#tmxk
- **Predecessors**: 0005-skill-integration-cli-contract
- **Overrides**: `0005:C1`, the `flow next` candidate clause — "flow next MUST return the legal recognized-outcome alphabet for the supplied state, plus candidate summaries …" read as guard-conditioned, and its stated ground "because next enumerates rather than selects" — narrowed to the match-conditioned predicate of C1 below; 0005's enumeration becomes `--all` (C2). Every other `0005:C1` obligation on `flow next` is unchanged. `internal/cli/flow_next_0005_test.go` assertions that relied on a match-excluded row being reported move under `--all` — an empty set over the shipped fixtures (A4), so the clause is a forward guard and C3 owes the discriminating fixtures instead. Renames the candidate's `unresolved` list to `unknown` and types its entries `{key, reason}`; adds the reason token `not-evaluated` for un-run gate ids on that CLI list only — `0007:C8`'s seam vocabulary stays at `absent`/`uncomparable` and its payload is not widened. Touches no clause of RDR 0007 or JDR 0001: `0007:REQ-78` (the two-valued match seam) is left deferred, with the reason recorded in Decision Rationale. `0005:C1`'s reader-narrowing clause ("exactly those readers serving an owned key some candidate row of the requested model requires") is READ, not overridden: the demand set gains each row's match-block owned keys for BOTH callers of `internal/cli/flow_exec.go::invokedReaders` — the widening 0005's own DEV-8 named as a seed and declined as unverified. That reading changes `flow resolve` on one model class (a match-only owned key), in both directions: a run that refused `flow-no-match` over an artifact holding the fact may now yield a plan; and — because the demand set is a union over the outcome's rows — a run whose selected row does not itself match on that key, over an unbound or refusing reader, turns from a plan into `flow-artifact-missing` (exit 2) or the reader's own refusal (exit 3), the class DEV-8 created for guard-only keys and accepted. Recorded in Consequences and pinned by S8; no shipped flow-reachable fixture or `models/rdr.toml` rule is in that class (A15).
- **Seam Lineage**: `flow next`'s probe shape, `internal/cli/flow_next.go::excluded` — no prior point-fix at this locus. The adjacent `internal/resolve` match seam carries one deferral, `0007:REQ-78`, which this RDR leaves in place. Below the ≥2 accretion floor; the profile above is set by the contract count, not by the floor.

## Problem Statement

A skill author (or the human driving one) runs `flow next --model <model> --artifact <state>` to ask "what is legal from here?" so the next action can be chosen without re-deriving the transition table by hand. Today the answer is wrong for that question: RDR 0005 made `next` enumerate rather than select, stripping each row's `match` (and `escape`) before probing, so every row not pruned by a guard is reported as a candidate — 21 of the 22 rules on `models/rdr.toml` at `stage=resolved`. The caller discovers this when the candidate list is the whole alphabet regardless of the supplied state, and cannot constrain a skill's next step from it (0005's own premortem, "next omits enough condition detail to constrain a skill", materialized).

The system-internal requirement is a single decision on the candidate predicate of `flow next`: whether "the legal recognized-outcome alphabet for the supplied state" means guard-conditioned (0005's predicate: candidate = not guard-excluded, match is the kernel's selection pattern and `next` must not select) or match-conditioned (candidate = match holds AND not guard-excluded); which of the two is the default and which rides `--all`; and whether guard verdicts on match-excluded rows are still reported. A sub-question inside the same contract: what a match atom means when its key is ABSENT from the supplied state — neither owned by an invoked reader nor supplied as a `--tag`. The kernel's `internal/resolve/resolve.go::TagSet.matches` folds an absent key into non-match, and `0007:REQ-78` recorded that fold as deferred; whether `flow next` inherits the fold (an absent key excludes) or reports it (the row stays a candidate and names the key) is the decision. The third case the fold covers, a present-but-conflicted key, is unreachable from the CLI over any model that loads (A3) and is therefore not a case this RDR has to sort. `--all` is the same contract's second surface, not a second contract. Help text follows the decision.

**What "constrained" means here, and where the outcome lands.** The outcome is a candidate list narrowed to what the supplied state can take — NOT a single next action. `flow next` still never selects (`0005:D-selection-predicate`, unchanged); choosing among survivors is `flow resolve`'s, and a caller facing several candidates is being told the state genuinely admits several. Success is therefore "no row is listed that the supplied state excludes, and every row listed that the supplied state could not decide names the undecided key", not a target cardinality. That criterion is universally quantified while the validation below is a set of existential fixture checks, MVV 3's `21→3` among them; the gap is closed by enumeration rather than by a property test — the `disposition` mini-check table lists every reachable input class and each row carries an oracle — so the `21→3` figure is a witness for one model, never proof of the general claim. The narrowing is real only to the extent the match keys are DECIDABLE from the invoked readers plus supplied tags: where a match key is owned and served by a reader (`models/rdr.toml`'s `stage`, A6) the list narrows 21→3; where it is observed and unsupplied (0010's decision-table class with no `--tag`, A5) every row remains a candidate and each names the undecided key under `unknown`. That second case is the old symptom's shape with a diagnosis attached, and it is the intended behaviour, not a residual defect — the alternative, dropping those rows, is Alternative 1, which this RDR rejects for emptying the list on partial state. The caller's remedy is named in the payload (bind the reader, supply the tag); the RDR does not promise narrowing where the state was never supplied. One more limit on the witness: `models/rdr.toml` authors one ordinary rule per (stage, outcome) — the guard-partitioned `finalize-pass`/`finalize-blocked` pair at `reconciled`/`advance` is the one exception — so its three survivors at `stage=resolved` are one per declared outcome — the narrowed candidate set names the same three outcomes `outcomes[]` already carried. What the caller gains there is the per-outcome ROW (rule id, gate ids, next tags) instead of 21 rows to sift, not a smaller outcome set; the outcome set itself narrows only on a model where some declared outcome has no row the state can take, which that model never authors.

## Critical Assumptions

- **A1 A one-row probe table handed to the kernel with the row's present-key match atoms retained can only yield a plan, `no_match`, `owned_state_unavailable`, or `guard_unevaluable`; `ambiguous_match` cannot arise, so the sub-question "does `ambiguous_match` leave the row a candidate" is vacuous on a probe.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/resolve/resolve.go::RefusalKinds` closes the set at five. `ambiguous_match` is produced only under the `default:` arm of `switch len(selected)` in `internal/resolve/resolve.go::Resolve` and of `switch len(viable)` in `internal/resolve/resolve.go::escapeOrRefuse`, both requiring ≥2; a one-row table makes `len(candidates) ≤ 1`, and the stripped escape list leaves `escapeOrRefuse` no rescue row, so it passes the original refusal through. `internal/resolve/resolve.go::gate` returns `owned_state_unavailable` then `guard_unevaluable` before that count. The fifth kind, `unmodeled_outcome`, is unreachable **by probe construction, not by one-row-ness**: `internal/cli/flow_next.go::excluded` binds `Outcomes: []string{row.Outcome}` alongside `Recognized: row.Outcome`, so `Table.models` holds trivially — C1's probe builder MUST preserve that pairing. The two non-refusal error returns (`internal/resolve/resolve.go::Table.CheckValid`, `internal/resolve/precondition.go::CheckInput`) are vacuous on the probe (escape stripped; no reserved `recognized` key in owned/observed/`RequiresOwned`) and the CLI already treats `err != nil` as non-excluding.
  - **If wrong**: a match-retained row could be dropped as "ambiguous" and vanish from the candidate list with nothing in `unknown` explaining it.
- **A2 The CLI already reports an absent match key: `internal/cli/flow_next.go::summarize` walks `Row.Atoms` — which carries match-block atoms alongside guard atoms (`internal/table/model.go::Row.KernelRow` routes them by `Block`) — against `internal/cli/flow_next.go::assembledView`, so a match key the view lacks lands in the candidate's undecided list today. That same presence test is the one C1's probe builder uses to decide which match atoms to hand the kernel.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/table/model.go::Row.Atoms` is documented as the one unified predicate set with each atom's authored block retained ("the Match/guard split is applied at the kernel handoff, never stored here", `0002:C16`) — there is no separate match field or pattern map on `table.Row`. `internal/table/model.go::Row.KernelRow` routes that same field by `a.Block == BlockMatch`, and `internal/cli/flow_next.go::summarize`'s loop `for _, atom := range row.Atoms { if _, known := view[atom.Key]; known { continue } … }` carries no `Block` test, so a match-block atom over a view-absent key lands in `unresolved` today (the loop's own comment says "guard atom", narrower than the code). Note `internal/cli/flow_next.go::assembledView` does **not** inject `recognized`, unlike the kernel's `assemble`; this is harmless because `internal/table/normalize.go` lifts the `recognized` atom out of the predicate set into `Row.Outcome`, so no `recognized`-keyed atom reaches `Row.Atoms` at all.
  - **If wrong**: an absent match key silently drops out of the report and C1's `unknown` clause needs new code in the CLI rather than the existing atom walk.
- **A3 A PRESENT-but-CONFLICTED match key — one key carrying two differing values within one provenance, the case `internal/resolve/resolve.go::TagSet.matches` folds into `false` beside absence — is unreachable from the CLI over any model that loads, so the CLI's presence test is a sound stand-in for the kernel's undecidable-match case and the kernel's fold can never fire on a CLI-built probe.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/table/load.go::checkAccessorBindings` refuses at MODEL LOAD any model where an owned tag is served by `!= 1` readers (`owned tag %s is served by %d readers; want exactly one`), and because it counts key OCCURRENCES — `for _, key := range r.Keys { readerCount[key]++ }` over every reader's declared key slice — a single reader declaring `keys = ["k","k"]` counts 2 and is refused by that same test. `internal/accessor/executor.go::readOutcome.classify` reduces a reader's output to exactly its declared `requested` keys through a map keyed on `v.Key`, so an undeclared key is dropped and a repeat collapses; a reader cannot emit an undeclared repeat. `internal/resolve/resolve.go::TagSet.merge` conflicts only WITHIN one provenance (`prior.provenance == prov`), and the observed route is refused by `internal/cli/flow_input.go::parseTags` as `flow-tag-duplicate`. The CROSS-provenance case — an observed `--tag` on a key an invoked reader owns — is not a conflict in the kernel (`merge` resolves it by provenance precedence, owned over observed) and is refused before that anyway: `parseTags` returns `flow-tag-owned` for any key in `flowbind.OwnedTags(m)`. C3's three pins turn the invariant the presence test rests on into asserted fact rather than assumption. The seam-shape half (that `matches` folds three cases into one `bool`) stands and is why the kernel, not the CLI, keeps deciding equality. A DEAD ROW (`0002:C13`) is a different thing from a conflicted key and does not reopen this: the conflict is between two literals in the row's PATTERN, which the kernel's conjunction decides once the key is present; a conflicted key is two values for one key in the VIEW, which is what the three refusals keep out.
  - **If wrong**: a conflicted key would reach the kernel through the CLI, `matches` would exclude its row as `no_match`, and C1 would report nothing for it — the silent drop P1 forbids. The recorded upgrade path is the exported per-atom verdict in Briefly Rejected, and `checkAccessorBindings` is the invariant a future relaxation must re-examine.
- **A4 The 0005 `flow next` tests move rather than rewrite: every oracle that asserts on a named candidate seeds the owned key its rows match on (`status=draft` in `flowMVVModel` / `flowGatedNextModel`), and the alphabet, gate-handling, reader-narrowing, and non-mutation oracles do not depend on a match-excluded row being reported.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: all 28 `flow next` oracles across `internal/cli/flow_next_0005_test.go`, `flow_mvv_0005_test.go`, `flow_adversarial_0005_test.go`, `flow_input_0005_test.go`, `flow_encoder_0005_test.go`, `flow_resolve_0005_test.go`, `flow_harness_0005_test.go`, and `reserved_key_0008_test.go` were enumerated and classified against C1's predicate. Every match atom in the `next` fixtures is `status eq "draft"` (`flowMVVModel`'s `advance-draft` / `hold-draft`; `flowGatedNextModel`; `flowGateDenyModel`; `advUnlessExcludedModel`; `advInExcludedModel`), and every seed supplies `status=draft`, so `status` is present in the view and holds in every case. Every exclusion the suite exercises is a **guard** exclusion (`all.flag eq false`, `unless.flag`, `all.tier in [mid,high]`), which C1 leaves untouched. Baseline green: `go test ./internal/cli -run 'Next|Req4[0-7]|Req3[5-9]|MVV|Adv1'` → ok. Counts, **against C1's candidate predicate**: PASSES-UNCHANGED 28, MOVES-UNDER-`--all` **0**, BREAKS **0**. This classification is predicate-scoped and does NOT clear the suite for the build as a whole: the `unresolved` → `unknown` payload change is an independent edit that mechanically breaks 5 assertions in the same suite (C3 names them, with line numbers).
  - **If wrong**: the 0005 suite splits into default / `--all` variants beyond the moves C3 names, and Testing Strategy must enumerate them.
- **A5 Over 0010's decision-table class, `flow next` with no `--tag` reports every ordinary rule and exits 0, each rule's observed match keys appearing in the candidate's undecided list — with the exception of a row whose only match atom is `recognized`, whose list is empty by construction — so `0010:A11` and 0010's MVV step 5 hold under C1; and a partial `--tag` set narrows the list to the rows whose supplied keys hold, a supplied-but-mismatched key dropping its row.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `docs/rdr/0011-flow-next-match-conditioned-candidates/evidence/spikes/a5-decision-table-no-tag.md`. Live over a four-row two-dimension decision table, no `--tag`: exit 0, all 4 ordinary rows reported, each carrying its observed match keys `status` and `size` in `unresolved`. Supplying `--tag status=Draft` moves `status` out of every `unresolved` while dropping no row *today* (match is stripped) — under C1 `status` is present in the view, so `final-small`/`final-large` are supplied-but-mismatched and the kernel refuses their probes `no_match`, dropping them; `draft-small`/`draft-large` survive with `unresolved = ["size", "answer"]`. Two limits recorded rather than papered over: (i) a genuine zero-owned 0010 decision table is **unloadable on this build** — an ordinary rule with no write block is refused `malformed_rule_shape`, and a write to an observed tag is refused `write_to_non_owned_tag` — so the fixture carries one dummy owned tag and `required` is `["answer"]`, never empty; the "empty `required`" half of `0010:A11` / `0010:MVV` step 5 is downstream of 0010 shipping and is **not attestable here**. (ii) Fixture B establishes this assumption's stated exception: a row whose only match atom is `recognized` reports `"unresolved": []` on the current build, because `internal/table/normalize.go` lifts `recognized` into `Row.Outcome` and C1 binds it to the row's own outcome, so it can never be an unresolved fact. It bounds the "every rule names its keys" claim without touching `0010:A11` or `0010:MVV` step 5, which speak only of exit 0 and empty `required`.
  - **If wrong**: 0010 owes a class-specific clause on this seam, or Alternative 1's strict predicate is what a decision table needs and the two classes diverge.
- **A6 `models/rdr.toml` at `stage=resolved` narrows from 21 reported rules to exactly the 3 rows whose `match.stage` holds over the artifact's owned `stage` — `prelock`, `resolve-route-back`, `resolve-abandon`, one per declared outcome. The 22nd row's absence from the baseline 21 is guard-driven, not stage-driven.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `docs/rdr/0011-flow-next-match-conditioned-candidates/evidence/spikes/a6-rdr-toml-narrowing.md`. Live against the current build, artifact `{"stage":"resolved","status":"draft","gate_passed":"false"}`: 22 `[[rule]]` blocks, 22 non-escape normalized rows, **21 candidates** — the one excluded row is `finalize-pass`, pruned on its `[rule.guard.all.gate_passed] eq = "true"`, so the 21 is a property of `gate_passed=false`, not of `stage=resolved`. A prototype of C1's probe over the same model — "probe carries the row's match atoms over keys PRESENT in the assembled view, drops match atoms over absent keys, drops the escape list, and binds `Recognized: row.Outcome`", i.e. exactly the shape C1 now fixes — returns **3** candidates; cross-checked by `grep -n 'eq = "resolved"' models/rdr.toml` → exactly three `[rule.match.stage]` hits (lines 212/222/232). Every rule carries a `[rule.match.stage]` block and `stage` is `required = true` and served by the invoked `rdr-status` reader, so C1's absent-key provision never fires for this model. The `recognized` sub-claim is dropped as vacuous: `internal/table/normalize.go` lifts `match.recognized` into `Row.Outcome`, leaving **0** such atoms in `Row.Atoms`; the probe's `Recognized` binding instead serves the kernel's `row.Outcome != in.Recognized` filter and keeps the probe from refusing `unmodeled_outcome` (A1).
  - **If wrong**: the observed defect is not closed by a match predicate and the problem statement's diagnosis (match stripped) is incomplete.
- **A11 CLI presence agrees with kernel presence for every key a row's match atom can name: `internal/cli/flow_next.go::assembledView` (owned ∪ observed, last write wins) and `internal/resolve/resolve.go::assemble` (owned ∪ observed ∪ `recognized`, provenance-precedenced) hold the same key SET over the inputs `flow next` passes, so a match atom the CLI keeps because its key is present is never over a key the kernel's view lacks, and vice versa.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: Both views are built from ONE `owned` slice: `internal/cli/flow_exec.go::flowRequest.runReaders` returns it once and `runFlowNext` passes that same slice to `assembledView(owned, req.observed)` and to `excluded(row, owned, req.observed)`. `internal/cli/flow_exec.go::kernelTags` is a pure type conversion over the observed tags — an unconditional append, no filter. Neither builder drops what the other keeps: `internal/accessor/model.go::OwnedSnapshot` omits absent keys from the slice both consume (`if v.Absent { continue }`), and an empty-but-present value yields a `Tag` with `Value: ""`, which `assembledView` stores as a present map entry (`_, known := view[atom.Key]` is true) and the kernel stores as a present `taggedValue` — agreement, not divergence. A reader that refuses aborts the whole verb (`if result.Refused() { return nil, nil, accessorFailure(...) }`), so no partial-view skew arises. The kernel's extra `recognized` key is unnameable by any atom: `internal/table/normalize.go` lifts the match-block `recognized` atom into `Row.Outcome`, and a `recognized` atom under `guard.all`/`guard.unless` is refused `CatMalformedOutcomeBinding` at load. Value disagreement is irrelevant — the CLI compares no value.
  - **If wrong**: a match atom over a key the CLI sees but the kernel does not would reach the kernel and be folded into `no_match`, dropping the row silently; the fix is to build the probe's key set from the kernel's view, which would need an exported presence query.
- **A12 Omitting the absent-key match tags from a one-row probe changes only whether the row reaches `gate`, and does so per KEY: the kernel filters on `matches` before `gate`, `gate` evaluates `row.Guard` and `RequiresOwned` independently of `row.Match`, and because every tag on one key shares that key's presence the filter hands a key's tags over together or omits them together — so a probe's guard verdict and owned-state disposition equal the full-table `flow resolve`'s for that row, and a key carrying several match tags (the `eq`+`in` pairing `atomsFromBlock` admits; `0002:C13`'s dead row when the literals differ) is decided by the kernel's conjunction exactly as `resolve` decides it, never split by the CLI.**
  - **Status**: Pending
  - **Method**: Source Search + Spike
  - **Evidence**: In `internal/resolve/resolve.go::Resolve` the match filter's `continue` (`if !view.matches(row.Match) { continue }`) precedes `selected, blocked := gate(candidates, in.Guards, view)`. `::gate(rows []Row, seam GuardEvaluator, view TagSet)` references `row.Guard` (`evaluateAtoms(row.Guard, seam, view)`) and, through `missingOwned`, `row.RequiresOwned`; neither `gate` nor `missingOwned` reads `row.Match`. A row ALL of whose match tags are omitted is an ordinary row with an empty match pattern, which `TagSet.matches` accepts unconditionally (an all-must-hold loop over zero tags returns `true`) and which nothing reclassifies — `Table.CheckValid` inspects only `len(row.Escape)` / `len(row.Writes)`. Every match tag is an equality: `internal/table/normalize.go::atomsFromBlock` refuses any operator but `eq`/`in` under a match block, the `in` expansion emits per-member `eq` rows (`0002:C13`), and `internal/table/model.go::Row.KernelRow` renders each `BlockMatch` atom as a bare `resolve.Tag{Key, Value}` — so omitting a tag relaxes the row rather than inverting a negation. A key may carry SEVERAL match tags: `atomsFromBlock` loops over OPERATORS inside its loop over KEYS, so `[rule.match.k]` carrying both `eq` and `in` emits two tags on `k` in each expanded row (reproduced live: `table.Load` accepts, `KernelRow().Match` carries two tags on one key); `0002:C13` names the differing-literal product a dead row and not a load failure. The per-key filter is licensed by presence being per-key, not by one-tag-per-key. To verify: (i) that no loader tier refuses the pairing (or which does), over a model authoring `eq` and `in` on one match key; (ii) a live dead-row fixture showing candidate-with-`{k, absent}` while `k` is absent and excluded once `k` is present at any value; (iii) that a live pairing (`eq = "x"` with `in = ["x", "y"]`) yields one live expanded row and one dead one, each decided by the kernel as `flow resolve` would decide it.
  - **If wrong**: the CLI would have to reason about a key's tags jointly — a literal comparison C1 forbids — and the recorded fallback is to demote such rows to RDR 0006's lint (a dead-row finding) rather than sort them in `flow next`.

- **A13 Reading `Refusal.Undecided` off the probe's `guard_unevaluable` refusal names every GUARD atom the kernel could not decide, with the kernel's own reason (`absent` / `uncomparable`), and adds exactly one fact class the CLI's view walk cannot see today: a guard atom over a PRESENT key the seam answered unevaluable. Every `absent` entry it carries duplicates one the walk already produces, so the merge is a set union on `{key, reason}`.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/resolve/guard.go::evaluateAtom` returns `ReasonAbsent` on `!present` and `ReasonUncomparable` on each of the ways a PRESENT key fails to decide that `0007:C8` enumerates — the seam answered unevaluable, a NIL seam, and a foreign or missing literal on `OpExists` (which returns BEFORE the presence check, so it fires whatever the key's presence). Of those four producers exactly ONE is reachable from `flow next`: the seam-unevaluable arm, `internal/guard/grammar.go::Evaluator.Evaluate` returning `GuardUnevaluable` for a present value the operator cannot parse (`lt`/`lte`/`gt`/`gte` on a non-integer, `in`/`contains` on a non-array literal or held value, and any unknown operator). The other three cannot arise: a conflicted key is refused before a probe exists (A3); `internal/cli/flow_resolve.go::guardSeam` returns `guard.Evaluator{}` and is never nil; and a non-boolean existence literal is refused at load by `internal/table/normalize.go` (`an existence literal is a boolean`, `0002:C8`), the case that function's own comment calls "the kernel's fail-closed backstop for a normalizer that did not reject it at load". The reason VALUE is identical across all four, so C1's `{key, reason}` set and its dedup rule are unaffected, and the two-member closure is independently pinned upstream by `TestReq49_ReasonIsAClosedNamedStringSetWithAnEnumerator` / `TestReq47_ReasonSetStaysAtExactlyTwoMembers` in `internal/resolve`, which S5 requires to stay green unchanged. `internal/resolve/guard.go::evaluateAtoms` walks every atom of the row with no `break` and appends each unevaluable one before combining the verdict, returning the payload only when the row verdict is `GuardUnevaluable`, sorted; `internal/resolve/resolve.go::gate` attaches it to the `guard_unevaluable` refusal, rows sorted. Today `internal/cli/flow_next.go::excluded` reduces the `Result` to a `bool`, so that payload never reaches `summarize`. Payload shape as `0007:C8` fixes it: `Refusal.Undecided []UndecidedRow`, each `{RuleID, SourceLocator, Atoms []UndecidedAtom{Key, Block, Operator, Literal, Reason}}`, every unevaluable atom of the row (no short-circuit), sorted; on a one-row probe at most one `UndecidedRow`. `summarize` projects each atom to `{Key, Reason}` and unions with the walk's entries on that pair — the walk (`for _, atom := range row.Atoms { if _, known := view[atom.Key]; known { continue } … }`) produces exactly the `absent` class, so a guard atom over an absent key yields one `{key, absent}`, never a second reason. The stated precedence limit is confirmed: `gate`'s `missingOwned` check returns the `owned_state_unavailable` refusal BEFORE the `undecided` collection loop begins, so an `uncomparable` guard atom on such a row never reaches the wire.
  - **If wrong**: `uncomparable` never appears on the CLI's list and the `{key, reason}` shape carries one live reason fewer; the shape still stands on `absent` vs `not-evaluated`.
- **A14 Non-registration of `--all` on `flow resolve` / `flow read-state` / `flow set-state` yields a deterministic disposition — pflag fails the parse before any `RunE`, and the CLI's OWN error gateway, not cobra's printer, maps it — so no `flow-*` code is reachable on that path. But the class it maps to, `command-error` + exit 2, is the SHARED usage bucket, so the class does not by itself name this flag: C2's negative is carried by a STRUCTURAL absence oracle, with the behavioural run as corroboration.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `docs/rdr/0011-flow-next-match-conditioned-candidates/evidence/spikes/a14-unknown-flag-class.md`. Live against the current build: all four verbs (the other three plus `next`, today's baseline) exit **2** with stderr `error: command-error: unknown flag: --all` in text mode, and under `--as=json` exactly one NDJSON line on stdout, byte-identical across the four: `{"code":"command-error","message":"unknown flag: --all","hint":"run with `--help` to list supported commands and flags"}`. The emission path is the STRUCTURED gateway, not cobra's printer — `SilenceErrors: true` on root and every verb suppresses cobra's own output; `internal/cli/root.go::cobraErrorToCLIError` stamps `Code: "command-error"`, `Group: clierr.GroupUserEnv`, and `respond.Fail` emits it, with `clierr.ExitCodeFor` mapping `GroupUserEnv` to exit 2 (`docs/cli-output-contract.md`: "Every other failure is exit 2"). `primeAsFlag` is why `--as=json` is still honoured despite the failed parse. All four structural conditions confirmed by source read: the sole persistent flag is `--as` (`internal/cli/root.go`), `flow` registers none of its own and both shared registrars write to `cmd.Flags()`; and repo-wide greps for `FParseErrWhitelist`/`UnknownFlags`, `DisableFlagParsing`, and `TraverseChildren` return zero matches. Three limits recorded rather than papered over: (i) `command-error` + exit 2 is the shared bucket for every unknown flag, unknown subcommand, and `Args` violation — the bare-`flow` RunE hand-authors the same pair — so the code alone cannot attribute a refusal to `--all`; (ii) the flag name survives only in the free-text `message` (`param` is empty), text owned by pflag rather than this repo; (iii) the `-a` shorthand yields the same code and exit but different message text (`unknown shorthand flag: 'a' in -a`). Hence the oracle C2 mandates is the structural one — the house idiom already shipped as `internal/cli/flow_input_0005_test.go::TestReq69_NoPlanFlagShipsOnAnyVerb` and `internal/cli/flow_surface_0005_test.go::TestReq12_FlowGroupDoesNotRedefineTheRootAsFlag`, both asserting a flag's ABSENCE by `Flags().Lookup(...) != nil` with a vacuity guard — and not a match on pflag's text. This repo is on record against the alternative: `internal/cli/flow_harness_0005_test.go`'s header calls asserting on cobra's generic `command-error` "exactly the tautology the red gate must exclude", and `flow_resolve_0005_test.go::TestReq56And119_EveryRefusalFamilyKeepsItsOwnIdentity` ships an adversarial oracle that fails if a refusal collapses into it.
  - **If wrong**: C2 owes a typed refusal for the flag on the other verbs, which is the flag-surface negative `evidence/critique/Charted.md` item 1 routes to a successor.

- **A15 Extending `internal/cli/flow_exec.go::invokedReaders`' demand set with each row's MATCH-block owned keys invokes strictly MORE readers and never fewer, for BOTH of its callers (`flow_next.go:97`, `flow_resolve.go:98`); it changes no reader's contract, is a no-op over a model with zero owned tags (0010's `decision-table` class, `0010:C4`), and changes `flow resolve` in exactly one class — a model with a match-only owned key — in two directions: a run that refused `flow-no-match` for want of a reader it never invoked may now yield a plan; and a run whose selected row does not need the key, over an unbound or refusing reader, turns from a plan into `flow-artifact-missing` (exit 2) or the reader's own refusal (exit 3), because the demand set is a union over the outcome's rows.**
  - **Status**: Pending
  - **Method**: Source Search + MVV Test
  - **Evidence**: To verify. The demand set today is `Row.RequiresOwned` ∪ `guardOwnedKeys(m, row)` (`flow_exec.go::invokedReaders`), and `guardOwnedKeys` skips any atom whose `Block` is not `all`/`unless`; `RequiresOwned` is derived from the write block and clear list alone (`normalize.go`, `renderWrites`). Adding match-block owned keys is a union with a third term, so the set only grows. The reading is `0005:C1`'s own — "exactly those readers serving an owned key some candidate row … requires" — applied the way 0005's DEV-8 (`docs/rdr/0005-skill-integration-cli-contract/artifacts/deviations.md`) applied it to guard keys; DEV-8 recorded the match-key case explicitly ("a `match` atom on an owned key a row does not write would demand a reader by the same argument, and would produce `flow-no-match` rather than `flow-guard-unevaluable` … a genuine RDR-SEED candidate") and declined it only as unverified. `flow_exec.go::runReaders` pre-checks every invoked reader's artifact role and aborts on an unbound role (`codeArtifactMissing`, exit 2) or a reader refusal (`accessorFailure`, exit 3), so the `flow resolve` arms above are the exact set. What the verification must establish: (i) the added term is EMPTY over a zero-owned-tag model; (ii) over a match-only-owned-key fixture, `flow next` invokes the reader and match-decides the row (S8), and `flow resolve --outcome <o>` over the same fixture yields a plan when the reader answers, `flow-artifact-missing` when its role is unbound, and the reader's refusal when it refuses — and `flow-no-match` for that key in none of them; (ii-b) the BREAKING arm, pinned rather than discovered: a second ordinary row for the same outcome that does not match on the key, the reader unbound — the pre-extension run yields that row's plan, the post-extension run `flow-artifact-missing` (exit 2), the DEV-8 class; (iii) `flow resolve` over every shipped fixture and `models/rdr.toml` is byte-identical before and after (every owned match key there is also written, so the term adds nothing).
  - **If wrong** (the term cannot be uniform): C1 owes a stated limit instead of the extension — the match-only owned key stays absent, the row reports `{key, absent}` with no caller remedy (`--tag` is refused `flow-tag-owned`), and the help C3 mandates must say so rather than implying the key can be supplied. A `next`-only term is NOT the fallback: it makes the two verbs disagree on the view.

## Proposed Solution

### Approach

`flow next` keeps asking the kernel one question per row, but stops taking the whole match pattern out of the question. It hands the kernel the row's match atoms over the keys the supplied state actually carries — owned tags from the invoked readers plus `--tag` observed tags — and lets the kernel decide equality on them. A match atom over a key the state does NOT carry is not handed over: the CLI cannot ask a question the kernel would answer with a fold (`TagSet.matches` returns `false` for an absent key exactly as for an unequal one, and `flow next`'s caller is entitled to the difference), so that atom is instead reported on the candidate under `unknown` as `{key, absent}`. The kernel's `no_match` — now reachable only for a present-and-unequal key — excludes the row; `guard_unevaluable` and `owned_state_unavailable` leave it a candidate with their facts reported, as today. The CLI decides PRESENCE only, which it already decides for reporting (A2); EQUALITY stays the kernel's. The restriction is per TAG, keyed on presence: every match tag on one key shares its key's presence, so a key's tags are handed over together or omitted together, and whether the tags handed over hold TOGETHER — including `0002:C13`'s dead row, two `eq` literals on one key that no value satisfies — is the kernel's conjunction, exactly as in `flow resolve`. The verb still never chooses among the rows that survive, never turns a gate deny into a refusal, and still returns the model's full declared `outcomes` alphabet: it narrows the **candidate list** to rows the supplied state can actually take, which is what "for the supplied state" meant to the caller.

`--all` restores 0005's enumeration — the match pattern ignored entirely — for the caller who wants the guard-conditioned alphabet rather than the state-conditioned one. The two modes differ only in whether match participates; reader narrowing, gate opt-in, the payload shape, and exit semantics are the same in both.

The change lands in the **CLI**, and `internal/resolve` is not touched. Two facts settle that placement. (i) The only match case the CLI cannot distinguish from the kernel's fold is a PRESENT-but-CONFLICTED key, and that case cannot reach the kernel from the CLI: `internal/table/load.go::checkAccessorBindings` refuses a model at load when an owned tag is served by other than exactly one reader (`owned tag %s is served by %d readers; want exactly one`), and `internal/cli/flow_input.go::parseTags` refuses a repeated `--tag` key as `flow-tag-duplicate` (A3) ⇒ over any model that loads, "present" and "comparable" coincide, so a presence test is exact, not approximate. (ii) Making the kernel's match seam three-valued cannot stop at reporting: the kernel's undecidable class, `guard_unevaluable`, is a RESOLUTION-level veto under `0007:C10` — "if any surviving candidate row's guard is GuardUnevaluable, the resolution MUST refuse `guard_unevaluable`; a decided-GuardTrue sibling MUST NOT be selected while an unevaluable candidate exists" — and `internal/resolve/resolve.go::Resolve` returns on `if blocked != nil { return refuse(in, *blocked), nil }` ABOVE the `switch len(selected)` that reaches `escapeOrRefuse`, so routing an undecidable match there would refuse a whole table on one unbound key (every rule in `models/rdr.toml` matches on `stage`) and would make the refusal unreachable by the escape row 0010's decision tables rescue through ⇒ the kernel's disposition for an undecidable match is a `flow resolve` design decision with a locked-peer collision, not a by-product of fixing `flow next`, and this RDR leaves it where `0007:REQ-78` left it. ⇒ a probe-shape change in `internal/cli/flow_next.go::excluded`, one added term in `internal/cli/flow_exec.go::invokedReaders`' demand set so the view carries the match keys the predicate reads (A15) — a term both of that function's callers take, so `flow resolve` over a match-only owned key stops refusing `flow-no-match` for want of a reader it never invoked (Overrides) — the `unknown` payload field on the candidate, a flag, and help text; no change under `internal/resolve`, no new refusal kind, no widening of `0007:C8`'s payload.

Sibling-path check (step 5): the presence decision C1 adds already exists at this locus — `internal/cli/flow_next.go::summarize`'s `if _, known := view[atom.Key]; known` over `assembledView` is the CLI's one presence test, and C1's probe builder reuses it rather than adding a second; searched `internal/cli` for any other view-presence test, none. The nearest kernel-side signal is the GUARD seam's `internal/resolve/guard.go::evaluateAtom` (`ReasonAbsent` / `ReasonUncomparable`), whose reason vocabulary the `unknown` list reuses for the guard facts it reads off the probe's payload rather than minting.

### Technical Design

Data flow, per ordinary row: `table.Row` → `Row.KernelRow()` (match/guard split by block) → probe in `excluded` (escape stripped; `Recognized = row.Outcome`; match atoms restricted to keys present in `assembledView` by default, omitted entirely under `--all`) → `resolve.Resolve` over the one-row table → `no_match` ⇒ excluded; `guard_unevaluable` / `owned_state_unavailable` / `Plan` ⇒ candidate → `summarize` builds `unknown` from three sources: the row's match-block atoms over view-absent keys and its owned keys the view lacks (the existing atom walk, reason `absent`); the guard atoms the kernel could not decide, read off a `guard_unevaluable` refusal's `Refusal.Undecided` (reasons as the kernel reports them, `absent` / `uncomparable`); and, absent `--evaluate-gates`, the row's gate ids (reason `not-evaluated`) → gates under `--evaluate-gates` for reported candidates only → payload. The alphabet `outcomes` is copied from the model in both modes (0005's DEV-4 reading stands: the alphabet is what may be requested; the candidates are what varies with state).

Equality stays wholly in the kernel: the CLI hands over the atoms whose keys the state carries and reads the answer. Match-block `in` atoms are already expanded into per-member `eq` rows at normalization (`0002:C13`) and set-kinded literals canonicalized by `internal/table/model.go::seamValue`, so every match atom reaching the kernel is an independent single-key equality — restricting the set by key presence is per-atom with no cross-atom interaction, and a row whose retained match pattern is empty matches unconditionally as it does today. `excluded` must return the kernel's `Result` (not a `bool`) so `summarize` can read the `guard_unevaluable` payload; this is the one signature change in the file, and the reason the probe result is kept rather than discarded.

The kernel is unchanged. `flow resolve` therefore still folds an absent match key into `no_match` at both selection sites — escapable, so a decision table's "otherwise" row (`0010:A9`) keeps rescuing — and `flow next` reports the same row as a candidate naming the key. That is the intended relationship between the two verbs (`next` reports; `resolve` selects) and is the same relationship they already have for `guard_unevaluable`, in the escapable direction rather than the vetoing one. Whether the kernel should ever distinguish the absent case at selection is `0007:REQ-78`'s open question; this RDR records why the answer is not free (Decision Rationale) and does not take it.

#### Normative Contracts

**C1**

```normative
flow next MUST report as a candidate exactly each non-escape row of the
requested model whose match atoms over keys PRESENT in the assembled view
all hold and whose guard the kernel does not decide false. Both verdicts
MUST be the kernel's, asked over a one-row probe table that carries the
row's match atoms RESTRICTED to keys present in the CLI's assembled view
(owned tags from the invoked readers plus observed --tag tags), no escape
list, and `Recognized` bound to the row's own outcome — the binding that
keeps the probe modelling its own outcome, so `unmodeled_outcome` is
unreachable. The CLI decides PRESENCE only, by the same view-presence test `summarize` already applies; it MUST NOT compare a match value. Present means the view HOLDS A VALUE for the key: a reader that reports an owned key absent leaves it out of the view (`internal/accessor/model.go::OwnedSnapshot` omits `v.Absent`), so that key is `{key, absent}`, never a silent exclusion. A build
that decides equality in the CLI has forked the seam.

The restriction MUST be applied as a filter over the CONVERTED probe —
`probe := row.KernelRow()`, then drop from `probe.Match` every
`resolve.Tag` whose `Key` the view lacks — never by rebuilding match tags
from `row.Atoms` in `internal/cli`. `Row.KernelRow` renders each match
atom's literal through `internal/table/model.go::seamValue`, which
canonicalizes a set-kinded value (sorted, compacted JSON array,
`SetEscapeHTML(false)`) against the row's own set-key list; re-deriving
those tags CLI-side would reimplement that encoding in a second place and
silently mis-compare a set-kinded match key against the kernel. Filtering after conversion is sound because `resolve.Tag` carries `Key`, and
presence is a property of the KEY: every match tag on one key shares one
presence verdict, so the filter hands a key's tags to the kernel together or
omits them together and can never split a key. A probe row MAY carry several
match tags on one key — `internal/table/normalize.go::atomsFromBlock` loops
over operators inside its loop over keys and a match block admits both `eq`
and `in`, so `[rule.match.k]` authoring both emits two tags on `k` in each
expanded row, and `0002:C13` names the unsatisfiable product ("two `eq` on
one key, different literals") a DEAD ROW, "not a load failure". Whether a
key's tags hold TOGETHER is the kernel's conjunction (`TagSet.matches` is
all-must-hold), never the CLI's: a dead row is a candidate carrying
`{k, absent}` while `k` is absent, exactly as any undecided row is, and is
excluded by the kernel once `k` is present at any value, exactly as
`flow resolve` never selects it. The CLI MUST NOT compare a row's literals
against each other to detect a dead row — that is a literal-to-literal
comparison this contract does not own, and `0002:C13` assigns dead rows to
RDR 0006's lint, not to this verb (A12).

The assembled view MUST actually carry the keys the predicate reads, so
`internal/cli/flow_exec.go::invokedReaders` MUST add each row's MATCH-block
owned keys to its demand set, which today is `Row.RequiresOwned` unioned
with the row's GUARD-owned keys only (`::guardOwnedKeys` skips any atom
whose `Block` is not `all`/`unless`, and `RequiresOwned` is derived from the
write block and clear list alone). Without that, an owned key a model
MATCHES on but never writes, clears, or guards demands no reader, is absent
from the view, and has its atom omitted from every probe — so every row
becomes a candidate carrying `{key, absent}` and the default silently
degrades to --all for that model, with no caller remedy: the key is
`provenance = "owned"`, so `--tag` is refused `flow-tag-owned` (A3).
`checkAccessorBindings` does not prevent this — it requires every DECLARED
owned tag to have exactly one reader, so the reader exists and is simply
never invoked. The added term is a READING of `0005:C1`'s narrowing clause
— "exactly those readers serving an owned key some candidate row of the
requested model requires" — by the argument 0005's own DEV-8 used to admit
guard-owned keys and, in the same entry, named for match keys and declined
as unverified: a row cannot be match-decided without the key. It therefore
binds BOTH callers of `invokedReaders` — `flow_next.go`
(`invokedReaders(req.model, "")`) and `flow_resolve.go`
(`invokedReaders(req.model, outcome)`) — and MUST NOT be scoped to `next`:
a `next`-only term would have the two verbs assemble different views over
one model, so `next` would report a row a candidate that `resolve` then
refuses `flow-no-match` for the key it never read. The demand set is a
property of the MODEL (and, for `resolve`, the requested outcome), not of
the mode: it is computed once per invocation, before the row loop, and is
identical under --all (where the match pattern takes no part in the verdict
but its keys are still read), which is what keeps the invoked reader set
and the view mode-independent as the paragraph below requires. The
`flow resolve` consequence is stated, not denied: over a model with a
match-only owned key, `resolve` today never invokes the serving reader and
refuses `flow-no-match` for every row matching on that key, over an
artifact that HOLDS the fact — DEV-8's ADV-2 shape with the other refusal;
after this term it invokes the reader, so a run that used to refuse
`flow-no-match` MAY now yield a plan, `flow-artifact-missing` (exit 2, the
reader's role unbound), or the reader's own refusal (exit 3). Each of those names the remedy the old refusal hid. The demand set is a union over the outcome's rows, so the reader is invoked even when the row `resolve` would select does not itself match on the key: over an UNBOUND or REFUSING reader such a run turns from a plan into exit 2/3, before the kernel and above the escape phase — the same class DEV-8 created for guard-only keys and accepted, because the model declares that reader (`checkAccessorBindings`) and binding it restores the plan. That class is named in Consequences and pinned by S8; it is not denied. No `flow resolve` run over a model WITHOUT such a key changes, and the term is empty over a model with zero owned tags (`0010:C4`).

A match atom whose key is ABSENT from the assembled view MUST be omitted
from the probe and MUST NOT exclude the row; it MUST appear in the
candidate's `unknown` list as `{key, absent}`. A match atom whose key is
present and whose value the kernel finds unequal yields `no_match`, and
`no_match` is the ONLY probe disposition that excludes: such a row MUST
NOT be reported, and under --evaluate-gates its gates MUST NOT run.
`guard_unevaluable` and `owned_state_unavailable` MUST leave the row a
candidate with their facts reported. A row ALL of whose match atoms are
omitted is probed with an empty match pattern, which matches
unconditionally; it is a candidate and every omitted key is listed
`absent`. The presence test is exact, not approximate, because a
PRESENT-but-CONFLICTED key cannot reach the kernel from the CLI over any
model that loads (A3; C3 pins the three refusals that make it so — two
readers on one owned key at load, a repeated --tag, and a --tag on an
owned key) — the only case in which the kernel's own `matches` would
return false for a present key without comparing it. The invoked reader
set and the assembled view are fixed ONCE per invocation, before the row
loop, and are the same under --all; `unknown` is never a function of row
order or of the flag.

This contract changes nothing in the `internal/resolve` PACKAGE: the
kernel's match seam stays two-valued, its selection and escape phases are
untouched, `0007:REQ-78` stays deferred, `0007:C8`'s payload is not
widened, and no sixth `RefusalKind` is minted (`0007:REQ-79`). The
`flow resolve` VERB changes in exactly the one class the demand-set
paragraph names, and nowhere else.

Every match atom over an absent key, every guard atom the kernel could not
decide, every owned key no invoked reader established, and (absent
--evaluate-gates) every gate id MUST appear in that candidate's `unknown`
list as a `{key, reason}` pair. Match atoms over absent keys and
unestablished owned keys carry `absent`, derived from the CLI's own view
walk. Undecided guard atoms carry the reason the kernel reports — `absent`
or `uncomparable`, `0007:C8`'s closed set — read off the probe's
`guard_unevaluable` refusal payload (`Refusal.Undecided`); a guard atom
over an absent key therefore appears once, the walk and the payload
agreeing on `{key, absent}`. A gate id un-run because --evaluate-gates was
not passed is neither `absent` (the gate is declared and its id is
reported) nor `uncomparable` (a gate id is not a view key with a value),
so it carries `not-evaluated`, which this RDR mints for the gate class
ONLY and which never appears on a match, guard, or owned-key entry.
Entries are deduplicated by `{key, reason}` pair, not by key, and sorted by
`(key, reason)`. One fidelity limit is stated rather than hidden: when the
probe refuses `owned_state_unavailable`, `gate` returns before it collects
the guard payload, so an `uncomparable` guard atom on that row is not on
the wire and is not reported; every `absent` fact on the row still is,
from the walk, and the missing owned key is the remedy the caller acts on
first.

`outcomes` MUST remain the model's full declared alphabet in every mode.
flow next MUST NOT choose among the candidates it reports and MUST NOT
turn a gate deny into a refusal.
```

⇒ this is the clause that overrides `0005:C1`'s candidate reading; the exclusion mechanism (kernel probe, decided-false-only) is 0005's, extended from guard-only to match-and-guard, and the absent-key disposition follows the posture `excluded`'s own doc already states for guards ("treating 'unknown' as 'false' would silently drop rows the caller is entitled to see"). Prior art for distinguishing absent from unequal at a selection seam: PTaCL's atomic target evaluates `1T`/`0T`/`⊥T` with absence tested before inequality and "can distinguish between a non-matching value for an attribute and a missing attribute"; SQL's `col = NULL` yields UNKNOWN, not false. Prior art for leaving the kernel's SELECTION fold alone: SQL at statement scope — a `WHERE` clause evaluating to UNKNOWN does not qualify the row and the statement proceeds — which `0007:C10` itself names as the counter-example to its veto.

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
The filter is scoped to the ATOM WALK's contribution ONLY. It MUST NOT
suppress a `{key, absent}` pair the OWNED-key walk independently produces
for the same key: C1 names `every owned key no invoked reader established`
as its own source, 0005 reports that fact today, and the sentence below
("everything else 0005:C1 requires … MUST hold identically in both modes")
forbids deleting it. The distinction is load-bearing wherever a key is BOTH
matched and written — which is every rule in `models/rdr.toml`, where
`stage` is a `[rule.match]` key AND a `[rule.write]` key and so enters
`RequiresOwned`. The oracle asserting the absence MUST therefore be written
over a match key that is NOT in the row's `RequiresOwned`; asserted over a
matched-and-written key it fails on the very model MVV 4 runs, and forcing
it to pass would delete a shipped 0005 fact.
Everything else 0005:C1 requires of flow next — the data minimum, gate
accessors only under --evaluate-gates and only for reported candidates, a
deny reported on its candidate with exit 0, no invented guard facts, the
narrowed invoked reader set — MUST hold identically in both modes. The
success payload shape MUST be identical in both modes: `unknown` is present
in both, as `[]` rather than omitted when it carries nothing, so one
consumer struct parses either mode. --all is part of request identity.
--all MUST NOT be accepted by flow resolve, flow read-state, or flow
set-state — satisfied by NON-REGISTRATION (registered on next only, not on
the shared registerSelectionFlags). The build MUST pin that negative
STRUCTURALLY, asserting --all is absent from those three verbs' flag sets
and from the flow group's own, guarded against vacuity by requiring the
four verbs to exist first — the idiom already shipped for an absent flag
(`flow_input_0005_test.go::TestReq69_NoPlanFlagShipsOnAnyVerb`,
`flow_surface_0005_test.go::TestReq12_FlowGroupDoesNotRedefineTheRootAsFlag`).
That oracle is the contract, because it holds whatever pflag's wording is
and it fails if a future persistent --all on flow or root silently widens
the surface. The behaviour is corroboration, not the assertion: pflag
fails the parse before any RunE, so no `flow-*` code is reachable and the
CLI's own gateway — not cobra's printer, which SilenceErrors suppresses —
maps it to `command-error` / GroupUserEnv / exit 2 (A14). `command-error`
+ exit 2 is the SHARED usage bucket for every unknown flag, unknown
subcommand, and Args violation, so it MUST NOT be asserted as if it named
this flag; the flag name lives only in pflag's free-text message. This
contract mints no new refusal code for it.
```

⇒ 0005's behaviour is preserved one flag away, so a caller that wanted the enumeration loses nothing and the 0005 oracles for it move rather than die.

**C3**

```normative
The command's Short and Long help MUST state that flow next reports the
candidates the supplied state can take (match and guard both holding or
undecided), that a match or guard key the state does not carry leaves a
row a candidate with the key listed under `unknown` and its reason named,
and that --all reports every row the guards do not exclude regardless of
match. The help MUST also state that a candidate is a row the supplied
state does NOT EXCLUDE, not a row flow resolve will select: a candidate
carrying an `unknown` entry may still be refused by flow resolve over the
same state, and the entry names what to supply. The 0005 flow next tests in internal/cli/flow_next_0005_test.go
MUST keep passing unchanged where they assert the alphabet, gate handling,
reader narrowing, determinism, and non-mutation; an assertion that
depended on a match-excluded row being reported MUST be re-homed under
--all, not deleted, and the file's header comment MUST name this RDR as
the source of the default.

Renaming the candidate's flat `unresolved` list to the `{key, reason}`
`unknown` list is a payload change independent of the candidate predicate,
and it breaks assertions the predicate leaves untouched. The build MUST
mechanically re-home every shipped read of the `unresolved` key to `unknown`
and its element type from string to `{key, reason}` — FIVE TEST reads on
this build: `internal/cli/flow_next_0005_test.go` (3 — the presence check
`:103`, the gate-id assertion `:163`, and the adversarial guard-fact
assertion `:404`), `internal/cli/flow_adversarial_0005_test.go` (1, `:446`),
and `internal/cli/flow_mvv_0005_test.go` (1, `:66`). A sixth occurrence of
the literal, `flow_next_0005_test.go:166`, is the `%#v` argument of the
`t.Fatalf` reporting `:163`'s failure — it re-homes with that read and is
not a read of its own. The count is of test reads only; the production
rename additionally touches the
`candidate.Unresolved` field and its `json:"unresolved"` tag, the field's
doc comment, the file header comment naming `unresolved`, `summarize`'s
`Unresolved: []string{}` initializer, its three `c.Unresolved = append(…)`
assignments and the `slices.Contains(c.Unresolved, …)` guards
beside them, and — the shipped USER-FACING string in this file, not a
comment — the cobra `Long` help at `internal/cli/flow_next.go:70-71` ("next
ENUMERATES rather than selects … their ids are reported as unresolved
facts"), which the help rewrite this contract's first paragraph mandates
already reaches. Two further shipped user-facing descriptions state the
OVERRIDDEN default and MUST be corrected with the help text, though neither
reads the payload field: `docs/cli-output-contract.md:125` ("Enumerate the
legal outcomes and their candidate rules") and `README.md:43` ("list legal
next outcomes").
None of these is an assertion and none of them the five counts. These are
mechanical re-homings, NOT the `--all` moves above, and MUST NOT be counted
against or excused by A4's
predicate-scoped BREAKS-0 result. The gate-id reads re-home to the
`not-evaluated` reason C1 names; a re-homing that has to invent a reason
token is not mechanical and is a defect. THREE of the five discard the
comma-ok of a helper (`flow_harness_0005_test.go::stringsAt`, which returns
`nil,false` on a missing key, a non-array value, or a non-string element) —
`flow_next_0005_test.go:404`, `flow_mvv_0005_test.go:66`, and
`flow_adversarial_0005_test.go:446` — and the last of those is a NEGATIVE
assertion (`if containsString(unresolved, "flag")`) that passes vacuously
on the empty slice, so a re-homed read MUST assert the ok bool, and a green
suite MUST NOT be cited as evidence that this contract was implemented.
`stringsAt` itself MUST NOT be repointed at the new element type: it is
shared harness used well beyond these five sites, and changing its return
shape would touch readers this contract has not counted. The re-homing adds
a SIBLING helper beside it (a `[]{key, reason}` reader over the same
payload) and switches the five reads to it. That helper is harness, not a
sixth re-homing: the census counts shipped READS of the `unresolved` key,
and a green suite still proves nothing without the ok-bool assertion above.

Because no shipped 0005 fixture contains a row whose match key is present
and unequal, the 0005 suite cannot by itself distinguish this predicate
from the stripped-match one it replaces. The build MUST therefore add
fixtures that discriminate: a row whose match key is present and UNEQUAL
(excluded by default, reported under --all; under --evaluate-gates its
gates do not run), and a row whose match key is ABSENT (candidate in both
modes, `{key, absent}` by default and no match entry under --all). That
absent-key row MUST match on a key it does NOT write or clear, writing some
OTHER key instead — an ordinary rule with no write block is refused
`malformed_rule_shape`, and a matched-AND-written key enters `RequiresOwned`,
where C1's owned-key walk produces the same `{key, absent}` pair that C2's
--all filter does not touch, so the row would report the pair in both modes
and discriminate nothing. It MUST
also add the three pins that license C1's presence test: a model declaring
one owned key served by two readers is refused at load
(`internal/table/load.go::checkAccessorBindings`); a repeated --tag key is
refused `flow-tag-duplicate`; and a --tag on a key the model declares
owned is refused `flow-tag-owned` — all before any probe is built. There is no
conflicted-key fixture for flow next: the case is unreachable (A3), and a
fixture that reaches the kernel's conflicted fold by constructing a
`TagSet` in-package tests `internal/resolve`, which this RDR does not
change. It MUST also add C2's structural negative — --all absent from the
flag sets of flow resolve, flow read-state, flow set-state, and the flow
group itself, with the vacuity guard — and register --all on next in the
per-verb flag map the surface suite already carries
(`flow_surface_0005_test.go::TestReq3_EachVerbRegistersItsNormativeFlagSpellings`).
These are TWO oracles, not one: `TestReq3` is presence-only — it asserts
each verb registers its normative spellings and says nothing about what a
verb does NOT register — so it homes the positive (`--all` on `next`) and
cannot carry C2's structural negative. The negative's home is the absent-flag
idiom C2 names (`TestReq69`/`TestReq12`), a separate assertion over the other
three verbs' flag sets and the flow group's own.
The `uncomparable` fixture S7 needs is a present value the operator cannot
parse (A13's one reachable producer): the other three producers are
unreachable and MUST NOT be fixtured, since a fixture that forces one
tests `internal/resolve`, which this RDR does not change.
```

⇒ the override is visible where a reader meets it — help text and the test file — not only in this record.

#### Mini-checks

Fired by cue. Cues absent: round-trip / fidelity — this RDR defines no import/export, parse/deparse, or serialize inverse.

**`authority` — source-authority census.** Cue: the candidate verdict is split across a CLI presence decision and a kernel equality decision.

| input/decision | writer (canonical) | readers | call sites | sibling arms |
| --- | --- | --- | --- | --- |
| match-key PRESENCE (which atoms the probe carries) | **CLI** `flow_next.go::excluded` over `assembledView` | the probe builder; `summarize` (same test, for `absent` entries) | `flow_next.go::excluded`, `::summarize` | one test, two uses — C1 forbids a second presence vocabulary; the kernel's `assemble` presence agrees on key set (A11) |
| which owned keys the view CARRIES (the demand set) | **CLI** `flow_exec.go::invokedReaders` | `runReaders` → `assembledView` | `flow_next.go:97` and `flow_resolve.go:98` (once per invocation, before any probe or kernel call) | three terms after C1: `RequiresOwned` ∪ guard-owned ∪ **match-owned** (A15), for BOTH callers — a `next`-only term would give the verbs different views; mode-independent, so `--all` reads the same view |
| match EQUALITY | **kernel** `resolve.go::TagSet.matches` | `Resolve`, `escapeOrRefuse` | unchanged | the CLI compares no value (C1) |
| guard verdict + undecided payload | **kernel** `guard.go::evaluateAtoms` | `gate`; now also `summarize` via `Refusal.Undecided` | unchanged in the kernel | `excluded` returns the `Result` instead of a `bool` (A13) |
| candidate `unknown` list | **CLI** `flow_next.go::summarize` | payload consumers | `flow_next.go` | three sources — walk (`absent`), payload (`absent`/`uncomparable`), gate ids (`not-evaluated`); dedup on the `{key, reason}` pair |

**`oracle` — test-discriminability.** Cue: C3 states outright that a green 0005 suite is not evidence the contract shipped.

| MVV / scenario | fails if X is wrong because Y | negative control |
| --- | --- | --- |
| MVV 2–3 / S1 | 21→3 narrowing: a stripped-match build reports 21, so the count discriminates | `--all` (MVV 4) must still report 21 |
| MVV 5 / S2 | the added row's `match.status` ≠ `draft` must be absent by default | same row present under `--all` |
| MVV 6 / S3 present-unequal | row excluded AND its gates do not run under `--evaluate-gates` | present-equal row: candidate, gates run |
| MVV 6 / S3 absent | candidate with `{key, absent}` | supplying `--tag key=<v>` moves it out of `unknown` |
| S5 pins | the two-reader model is refused at load; the repeated `--tag` is refused `flow-tag-duplicate`; a `--tag` on an owned key is refused `flow-tag-owned` — the invariants C1's presence test rests on | a one-reader model loads; a single observed `--tag` parses |
| S8 demand set | a match-ONLY owned key invokes its reader and is match-decided — on the un-extended demand set that reader is never invoked and every row returns `{key, absent}`, the default degraded to `--all` (A15); `flow resolve --outcome <o>` over the same model yields a plan (reader bound and answering), `flow-artifact-missing` (role unbound) or the reader's refusal (refusing) — `flow-no-match` for that key in none | same reader set under `--all` (mode-independent); zero-owned-tag model invokes nothing (`0010:C4`); `flow resolve` over `models/rdr.toml` and every shipped fixture byte-identical |
| S5 kernel | `go test ./internal/resolve` passes with NO test changed — the kernel is untouched | a kernel diff in the build is a defect, not an expected inversion |
| S7 guard `uncomparable` | a guard atom over a present key the seam cannot decide appears as `{key, uncomparable}` — only reachable through the payload read (A13) | the same atom over an absent key: `{key, absent}` once, not twice |
| MVV 7 | `unknown` present as `[]`, not omitted | absent-key candidate: `unknown` non-empty |
| S6 | **anti-oracle**: the 0005 suite passing is NOT evidence for C1 — it cannot distinguish the predicates (A4) | S2/S3 are the discriminating oracles |

**`disposition` — input class × outcome.** Cue: C1 sorts match input classes to reported-vs-excluded and to exit semantics.

| input class | candidate? | `unknown` entry | gates run under `--evaluate-gates` | exit |
| --- | --- | --- | --- | --- |
| match key present, equal | yes | none | yes | 0 |
| match key present, unequal | **no** (excluded) | n/a — not reported | **no** | 0 |
| match key absent from view | yes | `{key, absent}` | yes | 0 |
| match key present, conflicted | unreachable from the CLI (A3) — refused at load or at `--tag` parse before any probe | n/a | n/a | load / usage error |
| dead row — two match tags on one key with differing literals (`eq`+`in` authored on one key, `0002:C13`) | yes while the key is absent (undecided); **no** once it is present at any value — the kernel's conjunction fails, as it does in `flow resolve` | `{key, absent}` while absent; n/a once present | while absent yes; once present no | 0 |
| gate id, `--evaluate-gates` not passed | yes | `{gate-id, not-evaluated}` | no | 0 |
| all match atoms omitted ⇒ empty match PATTERN on the probe (never authored — an authored empty `[rule.match]` is refused at load) | yes (matches unconditionally) | `{key, absent}` per omitted key | yes | 0 |
| only match atom is `recognized` | yes | none — lifted into `Row.Outcome` at normalize | yes | 0 |
| guard decided false | **no** (excluded) | n/a | no | 0 |
| guard atom over absent key | yes | `{key, absent}` (walk and payload agree) | yes | 0 |
| guard atom over present key, seam unevaluable | yes | `{key, uncomparable}` (payload only) | yes | 0 |
| owned key no reader established | yes | `{key, absent}`; an `uncomparable` guard atom on the same row is not reported (precedence limit) | yes | 0 |
| any of the above, under `--all` | match takes no part; guard rules alone decide | no match entry, whatever the key's presence | per guard | 0 |

**`trace` — desk trace over the MVV.** Cue: C1, C2, C3 plus S1–S7 all bear on one output surface (the candidate payload). Witnesses from the A6 spike.

| step | assertions in force | witness |
| --- | --- | --- |
| 1. build with C1 default + `--all` | C2 (flag exists, default false); C2 (`--all` rejected by `resolve`/`read-state`/`set-state` by NON-REGISTRATION, A14) | flag name free: no `"all"` in `internal/cli` |
| 2. seed artifact `stage=resolved` | A6 | `{"stage":"resolved","status":"draft","gate_passed":"false"}` |
| 3. `flow next … --as=json` | C1 (candidates = present-key match holds ∧ guard-not-false); C1 (`outcomes` full alphabet) | 3 candidates `prelock`, `resolve-route-back`, `resolve-abandon`; `match.stage eq "resolved"` at `models/rdr.toml:212/222/232` |
| 4. re-run `--all` | C2 (predicate = 0005's); C1 (`outcomes` unchanged) | 21 of 22; `finalize-pass` guard-excluded on `gate_passed` — guard, not match, so it stays excluded in BOTH modes |
| 5. `flowGatedNextModel` ± `--all`, then `flowMVVModel` ± `--all` | C3 (added discriminating row); C1 guard exclusions untouched | `gated-excluded` absent in both — guard-excluded — over the fixture that AUTHORS it (`:317`), with `gated-reported` present as the control |
| 6. three match classes + the three pins | C1; `disposition` table above; C3 pins | the conflicted row is refused before a probe exists |
| 7. `unknown` as `[]` | C2 (identical payload shape both modes) | one consumer struct parses either |
| 8. match-only owned key | C1 (demand set = `RequiresOwned` ∪ guard-owned ∪ match-owned, both callers); A15 | the reader is invoked and the row is match-decided; same `readers` under `--all`; `flow resolve` over the same model no longer refuses `flow-no-match` for that key |
| 9. `make check` + suites | C3 (5 mechanical re-homings, ok-bool asserted); S5 (kernel suite unchanged); S6 (anti-oracle) | `git diff --stat internal/resolve` empty |

No CONTRADICTION row: step 4's `finalize-pass` is guard-excluded in both modes, which is consistent — C1 leaves guard exclusions untouched and C2 changes only match's participation.

#### Load-Bearing Decisions

- **Identity** — a `flow next` request is identified as `0005:D-identity` says, plus the `--all` bit; the same request in the same mode over the same model revision and artifact contents reports the same candidate set.
- **Naming** — the flag is `--all`, boolean. Precedent for exactly this meaning (widen past a conditioned default): `docker ps -a` "Show all containers (default shows just running)", `git branch -a`, `gh workflow list --all` "Include disabled workflows" (`langref/gh-cli/pkg/cmd/workflow/list/list.go`), `beads --all` "Include closed issues", `ps -A`. Rejected: `--enumerate` (describes 0005's verb, not the caller's intent, and has no precedent in the surveyed CLIs); `--ignore-match` / `--no-match` (names the mechanism rather than the result, and collides with the kernel refusal spelling); flipping the sense (`--select` / `--match` on a 0005 default) because the defect is the default and every skill would carry the flag (Alternative 2). Rejected: an enum (`--state=all|applicable`), because the axis has exactly two members — surveyed enums (`gh pr list --state open|closed|merged|all`) exist where a third value already does, and inventing one forces a value for the common case. If a third mode ever appears the migration is `gh`'s own (`--include-forks false|true|only`): promote to a `StringEnumFlag` and keep `--all` as the widest member. One asymmetry against the cited precedents is accepted: `docker ps -a` and `git branch -a` widen to rows carrying the SAME information as the default's, whereas here the widened rows carry strictly LESS — C2 strips match facts from `unknown` under `--all`, so a row visible only under `--all` does not say why the default excluded it. The name still fits (the axis is "how many rows", which is what every precedent's `--all` selects); the information difference is a consequence of `--all` meaning "do not apply the match predicate", and it is recorded in Failure Modes and Briefly Rejected rather than papered over.
- **Undecided reporting shape** — the per-candidate field is `unknown`, a list of `{key, reason}`. Rejected: keeping the flat `[]string` `unresolved`, which cannot carry the remedy. Precedent for the word and the shape: terraform/opentofu emits `after_unknown` in plan JSON as a structural sibling of `after` (`internal/command/jsonplan/plan.go`) and carries `"unknown"` as a first-class check status distinct from `"error"` (`internal/command/jsonchecks/status.go`), reserving the human-only rendering `(known after apply)` for text mode; kubectl's conditions use `True`/`False`/`Unknown` for the same "could not determine" sense. `unknown` is emitted as `[]` rather than omitted so one consumer struct parses both modes (C2). Ordering: entries are sorted by `(key, reason)`, so the merge of walk-derived, payload-derived, and gate-id facts is stable across builds — `D-identity` requires only set equality, but an unstable order would churn any golden payload a consumer diffs. Text mode (`--as=text`) carries the same pairs — both members of each, so the text caller gets the diagnosis and not just the symptom, which is the property that matters (terraform's precedent cited above splits the RENDERING, not the information). It does NOT get a bespoke `key (reason)` form: `flow next` declares no `FindingCarrier` and `respond.OK` therefore renders its payload through the shared generic flattener (`internal/cli/respond/text.go::writeTextPayload` → `::flatten`), which emits path-qualified leaf lines — `candidates[0].unknown[0].key: stage`, `candidates[0].unknown[0].reason: absent`. That flattener exists precisely so "the two modes cannot disagree about what the run reported", an invariant its own doc pins to `0005`'s REQ-7/REQ-9/REQ-11/REQ-120 — clauses this RDR does not override — so minting a per-verb text template to get `stage (absent)` would edit a predecessor's surface for cosmetics. The rendering is the gateway's to change, uniformly, if it is ever worth changing; this contract owns the payload's INFORMATION, not its text layout, and no scenario asserts a text-mode string.
- **Selection / predicate** — default: candidate ⇔ every match atom over a present key holds (kernel) ∧ guard not decided false (kernel); an absent match key neither excludes nor decides, it is reported (C1). `--all`: candidate ⇔ guard not decided false (C2). In both modes `next` reports all survivors and chooses none; exactly-one selection and refusal remain `flow resolve`'s (`0005:D-selection-predicate`).
- **Placement** — PRESENCE is decided in the CLI, EQUALITY in the kernel, and `internal/resolve` is not modified. The split is exact because the only present-key case the kernel would fold without comparing — a conflicted key — is refused before a probe exists (A3), and the CLI already owns the presence test (A2). Rejected: a kernel three-valued match verdict with `flow resolve` refusing on `indeterminate` (Alternative 4), because `guard_unevaluable` is a resolution-level veto under `0007:C10`, is unreachable by escape rows, and has no carrier for match facts. Rejected: a kernel-side report of undecided match atoms on `no_match` with selection unchanged (Alternative 5), because it changes a kernel payload to carry a fact the CLI already holds. The recorded upgrade path if A3's invariant ever relaxes is the exported per-atom verdict in Briefly Rejected — not a silent CLI comparison.
- **Undecided vocabulary** — reuses `0007:C8`'s closed reason set (`absent`, `uncomparable`) for every fact that comes off the kernel's guard payload, rather than minting one there; the CLI's own walk emits only `absent`. `absent` and `uncomparable` have different remedies (bind a reader or supply the tag; fix the value or the atom's literal), so the payload names which. The CLI's `unknown` list spans a third class the seam never sees — a gate id un-run because `--evaluate-gates` was not passed — which is neither `absent` (the id is declared and reported) nor `uncomparable` (a gate id is not a view key with a value). That class carries `not-evaluated`, minted here for the gate class only. `0007:C8`'s set stays at two members for the seam it governs; the split is that `0007:C8` types a REFUSAL payload while `unknown` is a CLI reporting surface with one more source. Rejected: overloading `absent` for un-run gates, which would tell the caller to supply something they cannot; rejected: a single opaque `unresolved` key list, which loses the remedy; rejected: a sixth `RefusalKind`, which `0007:REQ-79` forbids.

#### Illustrative Code

Illustrative — intent only.

```sh
# Default: rows the supplied state can take. At stage=resolved this lists the
# resolved→* rows, not 21 of the 22.
intrastate flow next --model models/rdr.toml --artifact rdr=./0011.md --as=json

# 0005's enumeration: every row the guards do not exclude, match ignored.
intrastate flow next --model models/rdr.toml --artifact rdr=./0011.md --all --as=json

# A decision table with no tags: every rule is a candidate; each observed
# match key is listed under `unknown` as absent. Supplying --tag a=x narrows it.
intrastate flow next --model table.toml --as=json
```

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Per-row kernel probe | `internal/cli/flow_next.go::excluded` | strips every match atom (`probe.Match = nil`); returns `bool`, discarding the `Result` | Extend | probe carries the row's match atoms over view-present keys by default, none under `--all`; returns the `Result` so `summarize` can read the guard payload |
| View presence test | `internal/cli/flow_next.go::summarize` over `::assembledView` | applied for reporting only; no `Block` filter | Reuse | the same test selects which match atoms the probe carries (C1) and which land in `unknown` as `absent` |
| Undecided reason vocabulary | `internal/resolve/guard.go::Reason` (`absent`, `uncomparable`), `::UndecidedAtom` | guard blocks only (`0007:REQ-78`) | Reuse | reused verbatim for the guard facts read off the payload; `not-evaluated` minted for the CLI's un-run-gate class only |
| Undecided payload carrier | `internal/resolve/resolve.go::Refusal.Undecided` on `guard_unevaluable` | guard atoms only; empty on `owned_state_unavailable` | Reuse | read by `summarize` for guard facts (A13); not widened |
| Match comparison | `internal/resolve/resolve.go::TagSet.matches` | folds absent and conflicted into `false` | Reuse (unchanged) | the CLI never hands it an absent key; conflicted cannot arise (A3) |
| Conflicted-key exclusion | `internal/table/load.go::checkAccessorBindings`; `internal/cli/flow_input.go::parseTags` (`flow-tag-duplicate`, `flow-tag-owned`) | none | Reuse | the three invariants C1's presence test rests on; C3 pins all three |
| Candidate reporting | `internal/cli/flow_next.go::summarize` | flat `[]string`, no reason | Extend | `unknown` carries `{key, reason}` from three sources (C1); dedup on the pair; the `--all` branch filters `BlockMatch` atoms out (C2) |
| Match/guard split on the probe row | `internal/table/model.go::Row.KernelRow` | none | Reuse | no table change |
| Refusal kind set | `internal/resolve/resolve.go::RefusalKinds` | pinned at five (`0007:REQ-79`) | Reuse | no sixth kind |
| Flag registration | `internal/cli/flow.go::registerSelectionFlags` + per-verb `cmd.Flags()` | no deny-list; a negative held only by non-registration is unpinned today | Reuse | `--all` registered on `next` only, and its absence on the other three pinned structurally the way `TestReq69` already pins `--plan` (C2) |
| Reader narrowing | `internal/cli/flow_exec.go::invokedReaders` | demand set is `RequiresOwned` ∪ guard-owned keys; `::guardOwnedKeys` skips `BlockMatch`, so a match-only owned key demands no reader | **Extend** | match-block owned keys join the demand set (C1) for BOTH callers — `flow next` and `flow resolve` — so each verb's view carries the keys its predicate reads; the set stays mode-independent; `flow resolve` changes on the match-only-owned-key class only (A15) |

### Decision Rationale

The caller's question is "what can I do from here", and every peer that answers it conditions the answer on the current state: stateless's `PermittedTriggers` returns only triggers whose guards are met in the current state; xstate v5 removed the unconditioned `nextEvents` and kept the conditioned `can()`; pytransitions keeps the declared form (`get_triggers`) but as a separate operator from the evaluated one (`may_*`); SCXML's `cond` gates enablement (Investigation). 0005's predicate was the declared-shape answer; in this table model the analogue of "current state" is the match block, so a state-conditioned `next` is a match-conditioned one. That settles the predicate. The second decision is WHERE the absent-key case is sorted — the CLI's presence test or the kernel's match seam.

Questions-Options-Criteria. Options: **O1** strict kernel match, absent excludes (Alternative 1); **O2** keep 0005's default, opt-in flag (Alternative 2); **O3** CLI presence-split, kernel untouched (chosen); **O4** kernel three-valued match, selection refuses `guard_unevaluable` (Alternative 4); **O5** kernel reports undecided match atoms on `no_match`, selection unchanged (Alternative 5). Cells are one clause; the deciding rows are marked.

| criterion | O1 strict | O2 opt-in | **O3 CLI split** | O4 kernel refuse | O5 kernel report |
| --- | --- | --- | --- | --- | --- |
| answers the user's question by default | no — partial state yields an empty list | no — the alphabet stays the default | **yes** — narrows where state is supplied, names the key where it is not | yes at `next` | yes at `next` |
| prior-art alignment (absent ≠ unequal; conditioned default) | contradicts PTaCL/SQL | contradicts the instance read | **aligned** — distinguishes at the atom, folds only at selection, as SQL does at statement scope | aligned at the atom; no precedent for the veto (0007:C10 says so of itself) | aligned |
| **peer-contract compatibility (0007:C10, 0010:A9) — deciding** | compatible | compatible | **compatible — kernel unchanged** | **collides**: whole-table refusal on one unbound key; escape rows cannot rescue `guard_unevaluable` | compatible — selection unchanged |
| **mechanism exists in the code — deciding** | yes | yes | **yes** — presence test (A2), load refusal (A3), payload carrier for guard facts (A13) | **no** — no "survives but non-selectable" partition; `Refusal.Undecided` is guard-only | needs a new payload field on `no_match` |
| blast radius | one function | one flag | **`internal/cli/flow_next.go` + fixtures**; `internal/resolve` diff empty | kernel seam, both selection sites, ≥7 kernel oracles, two callers of `Resolve` | kernel payload + CLI |
| reversibility | trivial | trivial | **trivial** — `--all` restores the set; the upgrade path if A3 relaxes is named | poor — a `flow resolve` semantics change with a locked-peer dependency | moderate |
| cost | lowest | lowest | **low** | highest | medium |

O3 wins on the two deciding rows and ties or leads on the rest. O4 loses on both deciding rows: it needs a kernel partition that does not exist, and where it would land, `0007:C10` makes it a table-wide veto. O5 is compatible but buys nothing: the only match fact the CLI cannot see for itself is the conflicted case, and that case is refused before a probe exists. Alternative 1 (strict kernel match) was rejected because it empties the list for any caller that supplies partial state, including 0010's no-tag decision table (A5), and contradicts the "undecided is not false" posture `excluded` already takes for guards. Alternative 2 (keep 0005's default, opt-in flag) was rejected because the defect is the default and the instance read says the conditioned form is the primary one.

Two clauses of O3 carry their own rationale, and neither moves the choice. **The filter's licence.** The per-key probe filter is NOT licensed on "at most one match tag per key" — `internal/table/normalize.go::atomsFromBlock` admits `eq`+`in` on one key, emitting two tags. It is licensed on presence-per-key: a key's tags are handed over or omitted together, and the kernel's all-must-hold conjunction decides them — which also disposes of `0002:C13`'s dead row without a CLI literal comparison (Briefly Rejected). O1/O4/O5 never filtered, so the matrix is unaffected. **The demand-set term's scope.** Three options: scope the match-owned term to `next`; extend it for both callers of `invokedReaders`; drop it and state the limit. `next`-only was rejected because the two verbs would assemble different views over one model — `next` reports a candidate, `resolve` refuses `flow-no-match` for a key it never read — a worse disagreement than the escapable one this RDR already accepts. Dropping it was rejected because the default silently degrades to `--all` on that model class with no caller remedy. Both-callers was chosen because it is `0005:C1`'s narrowing clause read the way 0005's DEV-8 already read it for guard keys — DEV-8 named this exact widening and declined it only as unverified — and because the `flow resolve` change it causes is confined to one model class: from a refusal the caller cannot repair (`flow-no-match` over an artifact holding the fact) to a plan or a refusal that names the reader — and, on a run whose selected row does not need the key, from a plan to that same named refusal when the reader is unbound, the class DEV-8 accepted for guard keys (the hardened critic's P-1/P-2; folded, not denied). It is recorded in Overrides and Consequences as a `flow resolve` behaviour change.

Why `0007:REQ-78` stays deferred rather than closed: closing it means deciding what `flow resolve` does with an undecidable match, and the only non-escapable class the kernel has for "undecidable" is a resolution-level veto by locked contract. Taking that on is a `flow resolve` design decision — with a counter-example in SQL's statement-scope fold and a direct dependency on 0010's default-row rescue — that no caller has asked for and this RDR's user does not need. The cost is recorded honestly in Consequences: `next` and `resolve` disagree on an absent-key row, in the escapable direction.

Premortem (hardened, iter-2): one draft-free critic, briefed on the re-decided approach, its load-bearing claims N1–N8, the rejection reasons, and three ledger seeds; ledger at `docs/rdr/0011-flow-next-match-conditioned-candidates/evidence/propose-premortem/iter-2/critic.md`, verdict PASS with mitigations. Folded: P-1/P-2 (the demand set is a union over the outcome's rows, so a `resolve` whose selected row does not need the match-only key still breaks over an unbound/refusing reader — named in C1, A15 ii-b, Consequences, Failure Modes, S8 as DEV-8's accepted class; "nothing that worked stops working" withdrawn); P-3 (presence = the view holds a value; a reader-reported-absent key is `{key, absent}`, C1); P-7 (both verbs' `readers` equal over one model and outcome, S8); P-8 (`--all` restores 0005's predicate, not its reader set, Consequences); P-9 (the shipped lint reports no dead rows; successor scope, Failure Modes). Not folded, with cite: P-4 (`flowbind.OwnedTags` is already A3's cite), P-5/P-6 (A12's evidence already carries `CheckValid` and the per-atom converter), P-10 (the Problem Statement already scopes the win to state-supplied queries). The iter-1 ledger (`propose-premortem/critic.md`) stands for the parts of the approach this pass did not re-decide.
Ground-sweep: clean (14 anchors) — the re-propose's added anchors, checked by a draft-free checker on `main` (2026-08-26, iter-3): 13 CONFIRMED (`atomsFromBlock`'s nested operator loop, the `eq`/`in` admission, `0002:C13`'s dead-row sentence, `TagSet.matches`, both `invokedReaders` call sites, `runReaders`' two refusal arms and their exit classes — `codeArtifactMissing` exit 2, `accessorFailure` exit 3 in its execution/incomplete-read classes, kernel `no_match` exit 2 — `0005:C1`'s narrowing clause, DEV-8's match-key sentences, `0010:C4`/`0010:A2`, `Resolve`'s filter-before-`gate`, `KernelRow`/`CheckValid`, and no match-only owned key in any flow-reachable fixture); 1 REFUTED, cosmetic and corrected in place (the Problem Statement's "one rule per (stage, outcome)" — `reconciled`/`advance` carries a guard-partitioned pair). The iter-2 sweep (58 anchors, all CONFIRMED) stands for the unchanged text; its note on `0007:REQ-78`/`REQ-79` living in RDR 0007's req-list artifact carries forward.
Joint-check: fired → 0010 (home: cli/0011:C1), re-run 2026-08-26 (iter-3) on the re-decided proposal. Open peers at depth 1 with Status Draft/Final: cli/0010 only (0001–0009 are `Implemented`). Shared modify-anchor: `internal/cli/flow_exec.go::invokedReaders` (9 hits in 0010) — this RDR EXTENDS its demand set with match-block owned keys for BOTH callers (C1), where 0010 records it Reuse/unchanged. Not a collision: `0010:A2` declares the `decision-table` class has zero owned tags ("With zero owned tags the reader demand set is empty, so `flow resolve` and `flow next` invoke no reader"), so the added term is empty over that class and `0010:C4`'s "unchanged" holds verbatim — including for `flow resolve`, which this pass now admits changes on the match-only-owned-key class of the `state-machine` class only, a class 0010 does not touch. Shared literals: `no_match` (0010: its "otherwise" escape row's class, `0010:A9`; here: the excluding probe disposition) and `flow next`/`flow resolve` as verb names; `flow-no-match`, the new literal this pass adds, has 0 hits in 0010; 0010's `unknown` hits are the phrase "unknown tag". Disposition unchanged — cite-don't-restate: C1 is the sole normative home of the `flow next` predicate and of the demand-set term, `0010:A11` cites it, and 0010's own `Joint-check:` line records the same fire and home, so symmetry holds without a peer edit. The coupling still runs one way: `internal/resolve` is untouched, so `0010:A9`'s escape row keeps rescuing `no_match`; `runFlowNext` never lists escape rows. Absence arm: this RDR narrows a reported set and converts no refusal into an acceptance; the `flow resolve` change is refusal→plan or plan→named refusal on a class no `Final` peer relies on (none exists); 0005's reliance on the enumeration and on DEV-8's demand set is named under `Overrides` and rides to 7.1. Bridge sub-check: n/a — neither plan retires a surface the other introduces. Not paused: the fork is answered and its home is unchanged.

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

### Alternative 4: Kernel three-valued match; selection refuses `guard_unevaluable` on an undecidable match

**Description**: Make `internal/resolve/resolve.go::TagSet.matches` return `match` / `no-match` / `indeterminate`; at both selection sites (`Resolve`, `escapeOrRefuse`) an `indeterminate` row is not selected and the resolution refuses `guard_unevaluable` carrying the undecided match atoms in `Refusal.Undecided`; `flow next` reads them off the probe's refusal.

**Pros**:

- One verdict for every caller of the seam; closes `0007:REQ-78`; `flow resolve` stops folding an undecidable match into an escapable `no_match`.

**Cons**:

- `guard_unevaluable` is a RESOLUTION-level veto: `0007:C10` — "if any surviving candidate row's guard is GuardUnevaluable, the resolution MUST refuse `guard_unevaluable`; a decided-GuardTrue sibling MUST NOT be selected while an unevaluable candidate exists" — and shipped `gate` discards `selected` on `len(undecidable) > 0`. Every rule in `models/rdr.toml` matches on `stage`, so one unbound key refuses the whole table where today the rows drop out and an escape row rescues.
- `Resolve` returns on `blocked != nil` above the `switch len(selected)` that reaches `escapeOrRefuse`, so the refusal is unreachable by any escape row — 0010's "otherwise" row (`0010:A9`, an escape rescuing `no_match`) stops rescuing on a partially-supplied decision table.
- The carrier does not exist: `Refusal.Undecided` is emitted only for rows whose GUARD verdict is `GuardUnevaluable` (`evaluateAtoms` returns `nil` for a decided guard), `gate` passes `row.Guard` — not `row.Atoms` — to it, and `Row.Match` is `[]Tag{Key,Value}` while `UndecidedAtom` needs `{Key,Block,Operator,Literal,Reason}` on `0007:C8`'s pinned sort key.
- Its justification rested on the conflicted case being user-reachable, which A3 refutes.

**Reason for rejection**: it needs a kernel partition ("survives into `gate` as a non-selectable row") the code does not have, and where it would land it collides with a locked peer contract and regresses 0010 — a `flow resolve` design decision smuggled in as a `flow next` fix. Recorded in `docs/rdr/0011-flow-next-match-conditioned-candidates-postmortem.md`.

### Alternative 5: Kernel reports undecided match atoms on `no_match`; selection unchanged

**Description**: Leave selection two-valued, but have `matches` also return the atoms it folded (absent / conflicted), carried on the `no_match` refusal in a new payload field, so `flow next` reads the distinction off the kernel without the CLI testing presence.

**Pros**:

- Single-sources the absent-vs-unequal distinction in the kernel; no locked-peer collision; escape rescue unchanged.

**Cons**:

- A kernel payload change (a new field on `Refusal`, or `Undecided` populated on a second kind against `0007:C8`'s wording) to carry a fact the CLI already has: absence is CLI-visible (A2) and the conflicted case cannot reach the kernel from the CLI (A3).
- Widens `flow resolve`'s `no_match` payload for every caller to serve one.

**Reason for rejection**: correct but unearned — it pays a kernel surface change for a distinction the CLI can make exactly over any model that loads.

### Briefly Rejected

- **Filter `outcomes` to outcomes with a surviving candidate**: 0005 already refused this reading (DEV-4 in `runFlowNext`); the alphabet is what may be requested.
- **A separate verb (`flow candidates`)**: two verbs for one question; 0005's four-verb closure stands.
- **Report match-excluded rows with an `excluded_by` field instead of dropping them**: rejected, but NOT because `--all` substitutes for it — C2 strips match facts from `unknown` under `--all`, so the two are not equivalent and `--all` recovers only the row's identity, never the atom that excluded it. Rejected because it inverts the contract's shape: C1's whole point is that an excluded row is one the supplied state CANNOT take, and a verb that reports non-candidates alongside candidates hands the caller back the filtering problem 0005 already failed at. The cost is real and accepted: a caller who wants to know WHY a row was excluded reads the model's match block against the payload's `owned`/`observed` by eye (Failure Modes). If that proves to be the common case in use, the successor is `excluded_by` on a `--why`-style opt-in, not a default-mode payload widening — recorded here so the option is reopened deliberately rather than rediscovered.
- **CLI-side dead-row detection** (compare a key's match literals against each other and drop or flag the row before probing): a literal-to-literal comparison the CLI does not own; `0002:C13` names the dead row RDR 0006's (lint), and the kernel's conjunction already excludes it the moment its key is present. `flow next` treats it as it treats any undecided row (Failure Modes).
- **Scoping the match-owned demand term to `next`** (leave `flow resolve`'s reader set as shipped): the two verbs would assemble different views over one model, so `next` reports a candidate `resolve` refuses `flow-no-match` — a disagreement in the unrepairable direction. Rejected for the uniform term (Decision Rationale).
- **An exported per-atom match verdict (`resolve.MatchVerdict(view, atom)`) the CLI calls instead of testing presence**: exports a seam helper for a case (conflicted) no CLI caller can reach. Named here as the UPGRADE PATH: if `checkAccessorBindings` or `parseTags` ever relax and a conflicted key becomes reachable, this is the change — not a CLI-side value comparison.

## Context

### Background

Observed 2026-08-26 while driving the consumer model (rdr#tmxk): `flow next --model models/rdr.toml --artifact rdr=<state with stage=resolved>` lists 21 of the 22 rules. The cause is by design in RDR 0005: `internal/cli/flow_next.go::excluded` nils `probe.Match` / `probe.Escape` before probing, under the comment "dropping a row whose match pattern the supplied facts do not satisfy would be a selection this verb was not asked to make" — the sentence this RDR revisits. 0005 rationale `0005:A3` holds that `next` exposes the alphabet without owning guard evaluation; selection (gate-then-count, exact-one survivor) belongs to the kernel. 0005's Briefly Rejected list also refused exit-code-only output because it "cannot carry legal outcome alphabets, conditional summaries" — the conditional summary is the surface this RDR sharpens. Constraints: RDRs are never amended in content (this is a new RDR that overrides 0005's clause); intrastate stays generic — no consumer (RDR-process) knowledge in code, docs, or fixtures; `make check` must pass. Not a facet of kata `zdat` / cli/0010 (owned-state optionality is a model-class decision; this is a verb-predicate decision) — cross-cite only; 0010 leaves `flow next` to this RDR (`0010:C4`, `0010:A11`).

### Technical Environment

Go CLI (`bin/intrastate`); `internal/cli/flow_next.go` (`excluded`, `summarize`, `assembledView`, candidate/guard reporting), kernel probe API `internal/resolve/resolve.go::Resolve` that today receives rows with match stripped, `internal/table/load.go` load-time accessor-binding checks, `internal/cli/flow_next_0005_test.go`. Conventions in `AGENTS.md` (respond gateway, CLIError codes, SilenceUsage). Design history: `docs/rdr/0001–0010`, `docs/jdr/0001`.

## Research Findings

### Investigation

Prior art was read class and instance. Class (StateMachineRes corpus, three queries, ledger `docs/rdr/0011-flow-next-match-conditioned-candidates/evidence/research/prior-art.md`): the "what next" operator exists in every peer that has a current state, and is conditioned on it. Instance: qmuntal-stateless's README "Introspection" — "a list of the triggers that can be successfully fired within the current state via the `StateMachine.PermittedTriggers` property" — implemented in `states.go::stateRepresentation.PermittedTriggers` as triggers with `len(tb.UnmetGuardConditions(...)) == 0` ⇒ state- and guard-conditioned, the conditioned form is the only form. pytransitions `transitions/core.py::Machine.get_triggers` returns triggers declared from the given states without evaluating conditions, while `Machine._can_trigger` (the `may_*` family) evaluates them ⇒ the declared-shape enumeration exists but as a distinct operator, which is what `--all` is here. xstate's `packages/core/CHANGELOG.md` records "Removed `MachineSnapshot['nextEvents']`" and `packages/core/src/State.ts::machineSnapshotCan` keeps the guard-evaluated `can(event)` ⇒ when forced to keep one, the peer kept the conditioned query. scxmlcc `doc/user-manual.md` (`cond`): "The transition is only executed if the condition evaluates to true" ⇒ enablement is event- and condition-conditioned. The state-machine peers gave no coverage for distinguishing an *absent* state key from a mismatching one — every peer assumes a current state is always present — so the search widened to the access-control and query literature, where partial input is the normal case (ledger `docs/rdr/0011-flow-next-match-conditioned-candidates/evidence/research/three-valued-match.md`). **PTaCL** (via Griesmayer & Morisset, ATRAP §2.1–2.2) models a *target* — key/value equalities selecting which rules apply — with a three-valued atomic evaluation, `⟦(n, v)⟧(q) = 1T if (n, v′) ∈ q and v = v′; ⊥T if (n, v′) ∉ q; 0T otherwise`, absence tested *before* inequality, and states outright: "PTaCL can distinguish between a non-matching value for an attribute and a missing attribute" ⇒ absent ≠ unequal at a selection seam is an established shape, which C1's reporting adopts. **SQL** is the mass-adoption instance and cuts both ways, which is what decides placement: at the atom, "When a NULL is involved in a comparison operation, the result is considered to be UNKNOWN" (Elmasri & Navathe §5.1) ⇒ do not fold absent into false where the caller can see it; at statement scope a `WHERE` clause evaluating to UNKNOWN does not qualify the row and the statement proceeds — the reading `0007:C10` itself gives as the counter-example to its own veto ⇒ leaving the kernel's selection fold in place has precedent, and replacing it with a veto (Alternative 4) has none. **SAML 2.0 Core** §2.5.1 ("the validity status … MUST be considered to be Indeterminate") and **XACML**'s `Indeterminate` cover the uncomparable case, which A3 shows the CLI cannot produce. ⇒ the literature supports distinguishing at the reporting surface and does not support a selection-scope veto; the chosen placement follows both halves.

### Key Discoveries

- **Documented** — `internal/cli/flow_next.go::excluded` strips `Match` and `Escape` from the probe and excludes only on `resolve.KindNoMatch`; its doc gives the "unknown is not false" rule this RDR extends to match. ⇒ the change is a probe-shape change in one function.
- **Documented** — `internal/resolve/resolve.go::TagSet.matches` collapses three cases into one `false`: absent (`!ok`), conflicted (`tv.conflicted`), and unequal. ⇒ the CLI must not hand it an absent key if the caller is to see the difference; it never hands it a conflicted one (A3).
- **Documented** — `internal/table/load.go::checkAccessorBindings` refuses a model whose owned tag is served by other than exactly one reader; `internal/cli/flow_input.go::parseTags` refuses a repeated `--tag` as `flow-tag-duplicate`. ⇒ over any model that loads, a present key is a comparable key; the CLI's presence test is exact (A3).
- **Documented** — `0007:C10` makes `guard_unevaluable` a resolution-level veto ("a decided-GuardTrue sibling MUST NOT be selected while an unevaluable candidate exists"), and `internal/resolve/resolve.go::Resolve` returns on `blocked != nil` above the escape switch. ⇒ routing an undecidable match into that class refuses whole tables and defeats escape rescue; the kernel is left untouched (Alternative 4).
- **Documented** — `internal/resolve/resolve.go::gate` calls `evaluateAtoms(row.Guard, seam, view)`; `Refusal.Undecided` is populated only for guard-undecidable rows. ⇒ the payload carries guard facts, which `summarize` can read (A13); it is not a carrier for match facts.
- **Documented** — the codebase is three-valued at the gate (`0004:C9`) and guard (`0007:C6`, `0007:C8`) seams and two-valued at the match seam by `0007:REQ-78`'s explicit deferral. ⇒ the reporting surface reuses `absent` / `uncomparable`; the seam's own posture is left as deferred.
- **Documented** — `internal/cli/flow_next.go::summarize` walks every `Row.Atoms` entry, match-block atoms included, against the view. ⇒ an absent match key is already an `unresolved` entry (A2), and the same test picks the probe's atoms.
- **Documented** — 0005 Decision Rationale premortem: "`next` omits enough condition detail to constrain a skill". ⇒ 0005 foresaw this failure and answered it with candidate *summaries*; the summaries were not enough because the *set* was unconditioned.
- **Documented** — `0010:C4`: "`flow next`, `flow read-state`, and `flow set-state` are unchanged by this RDR"; `0010:A11` defers the decision-table `next` behaviour to this RDR's predicate; `0010:A9` makes the "otherwise" row an escape rescuing `no_match`. ⇒ the absent-key rule must serve a no-tag decision table (A5), and the kernel's `no_match` fold must stay escapable.
- **Documented** — `internal/table/normalize.go::atomsFromBlock` loops over operators inside keys, so `[rule.match.k]` may author `eq` and `in` together and emit two match tags on `k`; `0002:C13`: "A product row whose expanded atoms cannot be satisfied together (two `eq` on one key, different literals) is a dead row for RDR 0006, not a load failure." ⇒ the probe filter's licence is presence-per-key, never one-tag-per-key; a dead row is the kernel's to exclude and lint's to report.
- **Documented** — 0005 DEV-8 (`docs/rdr/0005-skill-integration-cli-contract/artifacts/deviations.md`): the demand set was already widened once, uniformly for both verbs, by reading `0005:C1`'s "an owned key some candidate row requires"; the same entry names the match-key case ("would demand a reader by the same argument … a genuine RDR-SEED candidate"). ⇒ the match-owned term is that seed, taken for both callers; `flow resolve`'s change on the match-only-owned-key class is the DEV-8 shape with `flow-no-match` in place of `flow-guard-unevaluable`.
- **Documented** — prior art (Investigation). ⇒ conditioned is the default, enumeration is the secondary surface; distinguish at the atom, do not veto at selection.
- **Assumed** — every 0005 `next` oracle seeds the owned key its named rows match on (A4).

## Trade-offs

### Consequences

- Positive: a skill can read the next step from the candidate list at `stage=resolved` (A6) instead of the whole table.
- Positive: no table change, no kernel change, no new refusal code; the diff is `internal/cli/flow_next.go` plus one demand-set term in `internal/cli/flow_exec.go::invokedReaders` (A15) that both verbs take, a flag, help text, and fixtures, and the reason vocabulary is `0007:C8`'s, already shipped and tested.
- Positive: `absent` and `uncomparable` have different remedies (bind a reader or supply the tag; fix the value or the atom), and the payload names which, so the caller is told what to do rather than only that something is unknown. Guard atoms the seam could not decide over a present key become visible on `next` for the first time (A13). The limit, stated so the bullet is not read as more than it is: `absent` names the KEY, not which of its own two remedies applies — an owned key needs a reader bound and REFUSES `--tag` (`flow-tag-owned`), an observed key needs the `--tag`. The caller reads that off the model's `[tags.<key>] provenance`, which the payload does not carry; putting provenance on the entry is a payload widening this RDR does not take.
- Negative: an override of a locked 0005 clause — callers scripted against the enumeration must add `--all`. `--all` restores the candidate SET, not the payload SHAPE: the per-candidate `unresolved` (`[]string`) becomes `unknown` (`[{key, reason}]`) in BOTH modes, so an out-of-repo consumer parsing the old field must change even if it adopts `--all`. The two migrations are independent and a caller may owe both. `docs/cli-output-contract.md` pins the envelope, not per-command payload fields, so no CONTRACT clause of that doc is invalidated — but its `flow next` example still describes the verb as "Enumerate the legal outcomes and their candidate rules" (`:125`), which is the predicate this RDR overrides, and `README.md:43` describes it as "list legal next outcomes". Both are shipped user-facing prose stating the old default, so Phase 3 updates them alongside C3's help text and the test-file header. The point stands that no per-command payload contract had to be written to make this change legal; the edits are descriptive-prose corrections, not a contract amendment.
- Negative: `flow next` and `flow resolve` disagree on a row whose match key is absent: `next` lists it as a candidate naming the key; `resolve` folds it into `no_match` (escapable, so a default row may rescue). This is the same report-vs-select relationship the verbs already have for `guard_unevaluable`, in the escapable direction, and `0007:REQ-78` stays open — closing it is a `flow resolve` decision with a locked-peer collision (Decision Rationale).
- Negative: `flow resolve` changes on one model class. A model with a match-only owned key today never invokes that key's reader and refuses `flow-no-match` for every row matching on it — over an artifact that holds the fact; with the demand-set term it invokes the reader, so the run yields a plan, or `flow-artifact-missing` (exit 2) / the reader's own refusal (exit 3), each naming what the old refusal hid. The same term BREAKS one shape: a `resolve` whose selected row does not itself match on that key still invokes the reader (the demand set is a union over the outcome's rows), so over an unbound or refusing reader a run that returned a plan now exits 2/3 — DEV-8's class for guard-only keys, accepted here for the same reason (the model declares the reader; binding it restores the plan). `--all` on `next` restores 0005's PREDICATE, not 0005's reader set: on this class it invokes one more reader than 0005 did. No shipped flow-reachable fixture and no `models/rdr.toml` rule is in the class (every owned match key there is also written), so no existing output changes (A15 iii). Stated under Overrides as a reading of `0005:C1`'s narrowing clause — the one DEV-8 named and deferred.
- Negative: the presence test's exactness rests on three load/parse invariants outside `flow_next.go` (A3); C3 pins them so a relaxation fails a named oracle rather than silently reopening the fold. The upgrade path is recorded (Briefly Rejected).
- Negative: a candidate can still be listed that `resolve` would refuse (an undecidable guard, a missing owned key) — unchanged from 0005.

### Risks and Mitigations

- **Risk**: an absent match key makes the list look like the old wall of candidates.
  **Mitigation**: the key is named under `unknown` with reason `absent` on every such candidate (C1); the fix is on the caller's side (bind the reader / supply the tag) and is visible.
- **Risk**: a 0005 oracle silently passes for the wrong reason after the default flips.
  **Mitigation**: A4's run plus C3's explicit re-homing rule with the ok-bool asserted; Testing Strategy S2 pairs default and `--all` on the MVV fixture.
- **Risk**: a future relaxation of `checkAccessorBindings` or `parseTags` lets a conflicted key reach the probe, and `matches` folds it into a silent `no_match`.
  **Mitigation**: C3's three pins fail first; the recorded upgrade is the exported per-atom verdict (Briefly Rejected), never a CLI value comparison.

### Failure Modes

- **Visible**: a candidate the caller expected is absent. Diagnose with `--all`: if it appears there, the row was match-excluded — a match atom over a PRESENT key compared unequal (had the key been absent, C1 would have kept the row and named the key). `--all` localizes the cause to match, and no further, because C2 strips match facts from `unknown` in that mode and the payload carries no `match` block in either mode; the row's authored pattern is read from the model file, compared by eye against `owned`/`observed`, which the payload does carry. That eye-comparison is the whole remedy today: `intrastate` ships no model-inspection verb (`version`, `lint`, `flow next|resolve|read-state|set-state`; `dump` is asserted FOREIGN to the flow group at `flow_surface_0005_test.go:105`), and this RDR adds none — the `excluded_by` field that would have named the failing atom on the row itself is Briefly Rejected, and its reason is recorded there honestly.
- **Visible**: every row is a candidate with the same key under `unknown`, reason `absent` — the match key is neither owned by an invoked reader nor supplied; bind or supply it.
- **Visible**: a candidate carries `{key, uncomparable}` — a guard atom over a present key the seam could not decide (a literal the operator cannot compare against the value); fix the value or the atom. Only guard atoms can carry this reason on `next`.
- **Visible**: a model whose owned key is served by two readers does not load; a repeated `--tag` is refused `flow-tag-duplicate`. Both precede any probe, which is why `next` never has to sort a conflicted key.
- **Visible**: a candidate listed while a key is absent disappears once that key is supplied at ANY value — the row is `0002:C13`'s dead row (two differing literals on one match key, `eq`+`in` authored together). `next` cannot see it without comparing literals, which C1 forbids, and `flow resolve` would never select it either. Remedy: fix the rule. `0002:C13` assigns dead rows to RDR 0006's lint, but the shipped `lint` reports none today (searched `internal/`, no dead-row check) — a lint finding for them is successor scope, not this RDR's.
- **Visible**: a `flow resolve` that used to return a plan exits 2 (`flow-artifact-missing`) or 3 — the model matches on an owned key no rule writes, the requested outcome has a row on that key, and its declared reader is unbound or refusing; bind the reader (C1's demand-set clause, DEV-8's class).
- **Recovery**: none needed — `next` is effect-free in both modes; re-run with the corrected inputs.

## Implementation Plan

### Prerequisites

- [x] A1, A2, A4, A5, A6 verified. A6's prototype is the probe shape C1 fixes.
- [x] A3, A11, A13, A14 verified — the reachability invariant, the presence agreement, the guard-payload read, and the `--all` negative.
- [ ] A12 verified — the per-key filter under a key carrying several match tags (the `eq`+`in` pairing) and the dead-row disposition. Stage 4/6 closes it.
- [ ] A15 verified — the demand-set extension (match-block owned keys) invokes strictly more readers for both callers, is a no-op over the zero-owned-tag class, and changes `flow resolve` only on the match-only-owned-key class, in the plan / exit-2 / exit-3 arms named. Stage 6 closes it.

### Minimum Viable Validation

1. Build with the C1 default and the `--all` flag.
2. Seed an artifact for `models/rdr.toml` whose owned `stage` is `resolved`; run `flow next --model models/rdr.toml --artifact rdr=<it> --as=json`.
3. Assert: `candidates[]` is exactly `prelock`, `resolve-route-back`, `resolve-abandon` (3 rows, one per outcome), and no rule whose `match.stage` names another value; `outcomes[]` is still the model's full alphabet.
4. Re-run with `--all`; assert `candidates[]` has the 21 rows the baseline reported (the 22nd, `finalize-pass`, is guard-excluded on `gate_passed`, not match-excluded) and the same `outcomes[]`.
5. Run over the 0005 gated fixture (`flowGatedNextModel`, which is where `gated-excluded` is authored — `flow_fixtures_0005_test.go:317`; `flowMVVModel` has no such row, so asserting its absence there would pass vacuously in every build) with and without `--all`; assert the `gated-excluded` row is absent in both (guard-excluded), paired with its negative control `gated-reported` present in both. Then over `flowMVVModel` (`status=draft`), assert the row C3 adds — whose `match.status` names a value other than `draft` — is absent by default and present under `--all`.
6. Run over a model exercising the three reachable match cases: present-and-equal (candidate, the MATCH key absent from `unknown` — not necessarily an empty list, since a written owned key enters `RequiresOwned` and may report its own `absent` entry), present-and-unequal (excluded; gates do not run under `--evaluate-gates`), absent (candidate, `{key, absent}`), and a dead row (`eq`+`in` on one key with the `eq` literal outside the set): candidate carrying `{key, absent}` while the key is absent, excluded once it is supplied at any value. Then add `--tag key=<value>` and assert the list narrows to the rows whose match holds. Assert the three pins: a model declaring one owned key served by two readers is refused at load; a repeated `--tag` is refused `flow-tag-duplicate`; a `--tag` on an owned key is refused `flow-tag-owned`.
7. Assert `unknown` is present as `[]` rather than omitted on a candidate with nothing undecided, in both modes — the fixture must SUPPLY one, since no shipped 0005 fixture produces a fully-resolved candidate: a row whose match and guard keys are all established by the invoked reader, whose owned writes that reader also serves, and which is run WITH `--evaluate-gates` (without it, any declared gate id lands in `unknown` as `not-evaluated` and the list is non-empty by construction). Assert a guard atom the seam cannot decide over a present key reports `{key, uncomparable}` (A13).
8. Run over a model declaring an owned key some row MATCHES on but no row writes, clears, or guards (A15): assert its reader appears in `readers`, the key is in the view, and the row is match-decided rather than every row coming back a candidate carrying `{key, absent}`; assert the same `readers` set under `--all`. Run `flow resolve --outcome <o>` over the same model with the reader bound, unbound, and refusing: assert a plan, `flow-artifact-missing` (exit 2), and the reader's refusal (exit 3) respectively, and `flow-no-match` for that key in none; assert `flow resolve` output over `models/rdr.toml` and every shipped fixture is unchanged.
9. `make check` passes; `go test ./internal/resolve` passes with no test changed and `git diff --stat internal/resolve` empty; the 0005 `next` suite passes with only C3's five mechanical `unresolved`→`unknown` re-homings (ok-bool asserted), and the file header names this RDR.

End-state: `flow next` answers "what is legal from here" with the rows the supplied state can take, `--all` yields 0005's enumeration, every 0005 obligation other than the candidate predicate is unchanged, and the kernel is untouched.

### Phase 1: Probe and Reporting (CLI)

Intent: in `flow_exec.go::invokedReaders`, add each row's match-block owned keys to the demand set — a term both callers take — so each verb's view carries the keys its predicate reads (C1, A15) — the one edit outside `flow_next.go`; in `excluded`, carry the row's match atoms over view-present keys instead of stripping them (filtering `probe.Match` by `Tag.Key` AFTER `KernelRow()`, never rebuilding tags from `row.Atoms`), and return the kernel's `Result`; in `summarize`, build `unknown` as `{key, reason}` from the view walk (`absent`), the `guard_unevaluable` payload (`absent`/`uncomparable`), and un-run gate ids (`not-evaluated`), sorted by `(key, reason)`; emit `unknown` as `[]` rather than omitting it.

### Phase 2: `--all` and Help

Intent: register `--all` on `next` only, route it to a probe with the match pattern omitted and the `BlockMatch` atoms filtered out of `unknown` — the ATOM WALK's contribution only, leaving the owned-key walk's entries intact (C2) — and rewrite Short/Long help per C3.

### Phase 3: Fixtures and Docs

Intent: add C3's discriminating fixtures (present-unequal, absent), the three invariant pins, the default/`--all` pair, the `recognized`-only and all-absent rows, and the guard-`uncomparable` row; re-home the 5 shipped reads of the `unresolved` key to `unknown` (C3) via the sibling harness helper, the production sites C3 enumerates, and the cobra `Long` help; update the 0005 file header; document `flow next`'s payload — including `unknown` and the `--all` invocation — in `docs/cli-output-contract.md`, which currently shows only the invocation grammar and whose example comment still reads "Enumerate the legal outcomes and their candidate rules" (`:125`); and correct `README.md:43` ("list legal next outcomes"), the other shipped description of the old default.

## Validation

### Testing Strategy

Tightened at Resolve against the verified assumptions: every scenario names the code paths reviewed or the spike output that backs it.

1. **Scenario**: S1 — `models/rdr.toml` at owned `stage=resolved`, default mode (MVV 2–3). Backed by the A6 spike, which ran the baseline live and prototyped exactly this probe shape.
   **Expected**: `candidates[]` is exactly `prelock`, `resolve-route-back`, `resolve-abandon`; `outcomes[]` is the full alphabet. Normative fixture: `evidence/spikes/a6-rdr-toml-narrowing.md`.
2. **Scenario**: S2 — `flowGatedNextModel` and `flowMVVModel` (`status=draft`) with and without `--all` (MVV 4–5), plus the discriminating row C3 requires.
   **Expected**: `gated-excluded` absent in both **over `flowGatedNextModel`, the fixture that authors it** (`flow_fixtures_0005_test.go:317`), with `gated-reported` present in both as its negative control — asserted over `flowMVVModel`, which carries no such row, the oracle passes vacuously in every possible build and pins nothing; the added row absent by default and present under `--all`; `unknown` present as `[]` in both modes on a fully-resolved candidate; under `--all` a row with an absent match key reports NO match entry (the C2 departure) while by default it reports `{key, absent}`. C2's negative is asserted STRUCTURALLY — `--all` absent from the flag sets of the other three verbs and of the `flow` group, vacuity-guarded on the four verbs existing — not on pflag's error text; the behavioural run (exit 2, `command-error`, envelope on stdout under `--as=json`) corroborates it and is backed by the A14 spike.
3. **Scenario**: S3 — the three reachable match cases, one row each (MVV 6). Backed by the A5 spike's fixture shape.
   **Expected**: present-equal ⇒ candidate, **the match key** absent from `unknown`; present-unequal ⇒ excluded and its gates do not run under `--evaluate-gates`; absent ⇒ candidate with `{key, absent}`; a dead row (`eq`+`in` on one key, `eq` literal outside the set — `0002:C13`) ⇒ candidate with `{key, absent}` while absent, excluded once the key is present at any value, with no CLI literal comparison (A12). `--tag key=<v>` narrows the list to the rows whose match holds. The scenario is specified by row property, and the fixture must satisfy the model's own shape rules — an ordinary rule needs a write block, and a written owned key enters `RequiresOwned`, so a present-equal row may still carry an unrelated `{owned-key, absent}` entry from C1's owned walk. The assertion is therefore on the MATCH key's absence from `unknown`, never on `unknown` being empty; the fixture SHOULD write a key the invoked reader establishes so the list is empty in fact, but the oracle must not depend on it.
4. **Scenario**: S4 — a row whose only match atom is `recognized` (A5's fixture B), and a row ALL of whose non-`recognized` match atoms are over ABSENT keys, which is the reachable form of C1's empty-probe clause.
   **Expected**: the `recognized`-only row is a candidate with `unknown` empty of match facts; the all-absent row is a candidate — its probe carries an empty match pattern, which matches unconditionally — and every omitted key is listed `{key, absent}`, so its `unknown` is NOT empty. An AUTHORED empty `[rule.match]` is not a fixture here and MUST NOT be attempted: `internal/table/normalize.go` refuses it at load (`rule <id> carries no match block`, `malformed_rule_shape`), keying on key PRESENCE so an absent and a present-but-empty block are equally malformed — the empty match pattern C1 names is minted by the probe builder at runtime, never authored.
5. **Scenario**: S5 — the invariants C1's presence test rests on, and the kernel left untouched.
   **Expected**: a model with one owned key served by two readers fails to load on `checkAccessorBindings`' message (which counts key OCCURRENCES, so one reader declaring a key twice fails the same test); `flow next --tag k=a --tag k=b` is refused `flow-tag-duplicate` and `--tag <owned-key>=v` is refused `flow-tag-owned`, with no probe built; the kernel is untouched — the ONE mechanically checkable form of that claim is MVV 9's `git diff --stat internal/resolve` empty, and it is the form the build asserts; `go test ./internal/resolve` passing with no test file changed is its corroboration, not a second criterion, since a green suite cannot by itself prove nothing changed. `TestReq78_MatchPatternStillFoldsAbsenceIntoNonMatch`, the `match_conflicted_test.go` oracles, and the reason-set pins `TestReq49`/`TestReq47` (which hold `Reasons()` at exactly `absent`/`uncomparable`, the closure A13's merge relies on) stay green as shipped.
6. **Scenario**: S6 — the 0005 `next` suite after the default flips (A4, C3).
   **Expected**: the `flow next` oracles pass, none re-homed under `--all`; the 5 test reads of `unresolved` are re-homed to `unknown` with the ok-bool asserted; file header names this RDR. A green 0005 suite is NOT evidence C1 shipped — S2/S3 are. The "28 oracles" figure A4 reports is provenance for A4's verdict, not an oracle: the build asserts the checkable form (`go test ./internal/cli` green except C3's five named re-homings). "None re-homed under `--all`" is a claim about the DIFF, not a runtime assertion — no test can observe it — so it is discharged at review by reading the changed test files: A4 predicts zero such moves, and any move found refutes A4 rather than satisfying this scenario.
7. **Scenario**: S7 — the guard-payload read (A13) and its one precedence limit.
   **Expected**: a guard atom over a present key the seam cannot decide appears as `{key, uncomparable}` — the fixture builds that from A13's one reachable producer, a present value the operator cannot parse (`gt` on a non-integer, `contains` on a non-array); the same atom over an absent key appears once as `{key, absent}` (walk and payload agree, dedup on the pair); on a row that also lacks an owned key the probe refuses `owned_state_unavailable`, the row is still a candidate carrying the missing owned key and its `absent` facts, and the `uncomparable` guard atom is not reported — asserted so the limit is pinned rather than discovered. One payload invariant C1 states is pinned here rather than left to inspection: `unknown` is SORTED by `(key, reason)` — asserted on a fixture whose facts are minted out of that order, so an unsorted build fails. The dedup rule is pinned only in its reachable half: the walk and the payload both producing `{k, absent}` for one guard atom over an absent key MUST collapse to ONE entry (S7's fixture asserts exactly this). Whether dedup is on the pair or on the key is NOT discriminated by any reachable input, because no single key can carry two reasons — `uncomparable` requires a present key and `absent` an absent one, and `not-evaluated` is scoped to gate ids, a separate namespace from view keys. C1 specifies the pair because that is the honest description of the merge; the build is not asked to prove a distinction the input space cannot express.

8. **Scenario**: S8 — the match-only owned key (A15, C1's demand-set clause). A model declaring an owned key that some row MATCHES on but no row writes, clears, or guards, served by exactly one reader.
   **Expected**: the reader IS invoked (it appears in the payload's `readers`), the key IS in the view, and the row is match-decided rather than reported `{key, absent}`. This is the discriminating oracle for the demand-set extension: on the un-extended set the reader is never invoked and every row comes back a candidate carrying that key `absent` — the default silently degraded to `--all`. Negative controls: the same model under `--all` reports the SAME reader set (the demand set is mode-independent, C1); and a model with zero owned tags demands nothing and invokes no reader, unchanged (0010's `decision-table` class, `0010:C4`). The `flow resolve` arms over the same model: with the reader bound and answering, a plan; with its role unbound, `flow-artifact-missing` (exit 2); with the reader refusing, the reader's refusal (exit 3) — and `flow-no-match` for that key in none, which is the un-extended set's answer in all three. The BREAKING arm: add a second ordinary row for the same outcome that does not match on the key; with the reader unbound, assert the pre-extension build returns that row's plan and the post-extension build `flow-artifact-missing` — pinned so the DEV-8 class is a known cost, not a surprise. Both verbs' `readers` over the same model and outcome are EQUAL (one `invokedReaders`, one view), asserted on a fixture whose reader serves several keys. `flow resolve` over `models/rdr.toml` and every shipped fixture is byte-identical before and after (A15 iii).

Done: S1–S8 green and `make check` passes.

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

- RDR 0005 `0005:C1` (the overridden clause), `0005:A3`, `0005:A7`, `0005:D-selection-predicate`, Decision Rationale premortem; RDR 0010 `0010:C4`, `0010:A9`, `0010:A11`; RDR 0007 `0007:C8`, `0007:C10`, `0007:C11`, `0007:REQ-78`, `0007:REQ-79`; JDR 0001 §D2, §D12.
- Source reviewed: `internal/cli/flow_next.go`, `internal/cli/flow_exec.go`, `internal/cli/flow_input.go`, `internal/cli/flow.go`, `internal/table/model.go`, `internal/table/load.go`, `internal/resolve/resolve.go`, `internal/resolve/guard.go`, `internal/cli/flow_next_0005_test.go`, `internal/cli/flow_fixtures_0005_test.go`, `docs/cli-output-contract.md`.
- Prior art, state machines (sibling checkouts under `../state-machines`): `repos/qmuntal-stateless` (`README.md` Introspection; `statemachine.go`, `states.go`), `study/transitions` (`transitions/core.py`), `repos/xstate` (`packages/core/CHANGELOG.md`, `packages/core/src/State.ts`), `repos/scxmlcc` (`doc/user-manual.md`). Ledger: `docs/rdr/0011-flow-next-match-conditioned-candidates/evidence/research/prior-art.md`.
- Prior art, three-valued selection seams: PTaCL via Griesmayer & Morisset "ATRAP" §2.1–2.2 and Table 1; SAML 2.0 Core §2.5.1; Cedar (Cutler et al. 2024) §Related Work; Elmasri & Navathe *Fundamentals of Database Systems* §4.1, §5.1 Table 5.1; Celko *SQL Database Programmers Handbook* §"The Null of It All". Ledger: `docs/rdr/0011-flow-next-match-conditioned-candidates/evidence/research/three-valued-match.md`.
- Prior art, CLI shape: `langref/gh-cli` (`pkg/cmd/workflow/list/list.go`, `pkg/cmd/pr/list/list.go`, `pkg/cmd/search/repos/repos.go`), `langref/helm` (`pkg/cmd/list.go`), `langref/beads` (`cmd/bd/query.go`), `langref/opentofu` (`internal/command/jsonplan/plan.go`, `internal/command/jsonchecks/status.go`); `docker ps`, `git branch`, `ps` man pages.
- Internal precedent: JDR 0001 principles P1/P2/P4/P7 and §D2, §D6, §D9, §D12; `0004:C9` (`indeterminate`); `0007:C6` (strong-Kleene), `0007:C8` (`absent`/`uncomparable`), `0007:C10` (resolution-level veto), `0007:REQ-78` (the deferred match seam), `0007:REQ-79` (five refusal kinds).
- Spike and lens evidence: `evidence/spikes/a6-rdr-toml-narrowing.md`, `evidence/spikes/a5-decision-table-no-tag.md`, `evidence/critique/diff.md`, `docs/rdr/0011-flow-next-match-conditioned-candidates-postmortem.md` (escaped-defect ledger).
- kata `1mv1`; rdr#thsc; rdr#tmxk.
