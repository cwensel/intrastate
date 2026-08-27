# REQ List — RDR 0011 `flow next` selects by match; `--all` enumerates the alphabet

Phase 0 audit artifact. Every clause below is a testable obligation drawn from
`docs/rdr/0011-flow-next-match-conditioned-candidates.md`. Quotes for fenced
elements are verbatim — copied from the projector (`rdr inspect --select <id>`),
never transcribed by hand; quotes for testable prose outside the fences are
copied from the record at the cited line.

Element ids (`0011:C1`, `0011:C2`, `0011:C3`, `0011:MVV`, `0011:S8`,
`0011:D-naming`) are carried wherever a REQ derives from a labelled element, so
a later stage can trace the REQ back to its contract.

**Coverage shape.** `counts.elements` for this record: `C = 3`, `MVV = 1`,
`S = 8`, `D = 6`, `F = 7`, `A = 12`, `ALT = 4`, `BR = 7`, `G = 1`. The three
`normative` fences carry the bulk of the obligation (C1 is 159 lines and
paragraph-dense). Testable prose outside the fences was read separately and
mined: `Technical Design` (the data-flow pipeline and the `excluded` signature
change), `Load-Bearing Decisions` (the flag spelling, the payload shape and its
text rendering, the vocabulary), the `disposition` mini-check table (the
input-class × outcome oracle grid), the `authority` mini-check table, `MVV`
steps 1–9, `Testing Strategy` S1–S8, `Failure Modes` F1–F7, `Consequences`, and
`Cross-Cutting Concerns` (`0011:G-cross-cutting`). `Alternatives Considered`
and `Briefly Rejected` are rationale, not obligation — with two exceptions
mined as negative REQs (BR4/BR7's ban on a CLI literal comparison, already
carried by C1; BR6's ban on scoping the demand term to `next`, already carried
by C1) — so no separate REQ is minted for ALT1/ALT2/ALT4/ALT5.

Source sections are abbreviated:

- `NC` = Proposed Solution / Technical Design / Normative Contracts (fenced)
- `AP` = Proposed Solution / Approach
- `TD` = Proposed Solution / Technical Design (unfenced prose)
- `MC` = Technical Design / Mini-checks (`authority`, `oracle`, `disposition`, `trace`)
- `LBD` = Technical Design / Load-Bearing Decisions
- `EIA` = Proposed Solution / Existing Infrastructure Audit
- `OVR` = Metadata / Overrides
- `CONS` = Trade-offs / Consequences
- `FM` = Trade-offs / Failure Modes
- `MVV` = Implementation Plan / Minimum Viable Validation
- `PH` = Implementation Plan / Phases 1–3
- `TS` = Validation / Testing Strategy (S1–S8)
- `CC` = Finalization Gate / Cross-Cutting Concerns

**Standing note on the override.** This RDR overrides one clause of the locked
predecessor `0005:C1` — the `flow next` candidate clause — and nothing else of
it. Every other `0005:C1` obligation on `flow next` (the data minimum, gate
opt-in, exit semantics, the narrowed reader set, non-mutation, the full
`outcomes` alphabet) is unchanged and rides in both modes; REQs below that
restate a 0005 obligation are marked **[0005-carried]** and are obligations to
NOT break, not licence to re-implement 0005.

---

## Candidate predicate (default mode)

- [REQ-1] "flow next MUST report as a candidate exactly each non-escape row of the requested model whose match atoms over keys PRESENT in the assembled view all hold and whose guard the kernel does not decide false." — (NC, `0011:C1`)

- [REQ-2] "Both verdicts MUST be the kernel's, asked over a one-row probe table that carries the row's match atoms RESTRICTED to keys present in the CLI's assembled view (owned tags from the invoked readers plus observed --tag tags), no escape list, and `Recognized` bound to the row's own outcome — the binding that keeps the probe modelling its own outcome, so `unmodeled_outcome` is unreachable." — (NC, `0011:C1`)

- [REQ-3] "The CLI decides PRESENCE only, by the same view-presence test `summarize` already applies; it MUST NOT compare a match value." — (NC, `0011:C1`) — the negative is the assertable half: no value comparison in `internal/cli`.

- [REQ-4] "Present means the view HOLDS A VALUE for the key: a reader that reports an owned key absent leaves it out of the view (`internal/accessor/model.go::OwnedSnapshot` omits `v.Absent`), so that key is `{key, absent}`, never a silent exclusion." — (NC, `0011:C1`)

- [REQ-5] "A build that decides equality in the CLI has forked the seam." — (NC, `0011:C1`) — the review-time form of REQ-3.

- [REQ-6] "The restriction MUST be applied as a filter over the CONVERTED probe — `probe := row.KernelRow()`, then drop from `probe.Match` every `resolve.Tag` whose `Key` the view lacks — never by rebuilding match tags from `row.Atoms` in `internal/cli`." — (NC, `0011:C1`) — discriminated by S3(a)'s set-kinded fixture (REQ-102).

- [REQ-7] "Filtering after conversion is sound because `resolve.Tag` carries `Key`, and presence is a property of the KEY: every match tag on one key shares one presence verdict, so the filter hands a key's tags to the kernel together or omits them together and can never split a key." — (NC, `0011:C1`) — no partial-keep on any key.

- [REQ-8] "A probe row MAY carry several match tags on one key" — (NC, `0011:C1`) — the filter MUST tolerate it; `atomsFromBlock`'s `eq`+`in` pairing is the shipped producer.

- [REQ-9] "Whether a key's tags hold TOGETHER is the kernel's conjunction (`TagSet.matches` is all-must-hold), never the CLI's: a dead row is a candidate carrying `{k, absent}` while `k` is absent, exactly as any undecided row is, and is excluded by the kernel once `k` is present at any value, exactly as `flow resolve` never selects it." — (NC, `0011:C1`)

- [REQ-10] "The CLI MUST NOT compare a row's literals against each other to detect a dead row — that is a literal-to-literal comparison this contract does not own, and `0002:C13` assigns dead rows to RDR 0006's lint, not to this verb (A12)." — (NC, `0011:C1`) — also `0011:BR4`.

- [REQ-11] "A match atom whose key is ABSENT from the assembled view MUST be omitted from the probe and MUST NOT exclude the row; it MUST appear in the candidate's `unknown` list as `{key, absent}`." — (NC, `0011:C1`)

- [REQ-12] "A match atom whose key is present and whose value the kernel finds unequal yields `no_match`, and `no_match` is the ONLY probe disposition that excludes: such a row MUST NOT be reported, and under --evaluate-gates its gates MUST NOT run." — (NC, `0011:C1`)

- [REQ-13] "`guard_unevaluable` and `owned_state_unavailable` MUST leave the row a candidate with their facts reported." — (NC, `0011:C1`)

- [REQ-14] "A row ALL of whose match atoms are omitted is probed with an empty match pattern, which matches unconditionally; it is a candidate and every omitted key is listed `absent`." — (NC, `0011:C1`)

- [REQ-15] "The invoked reader set and the assembled view are fixed ONCE per invocation, before the row loop, and are the same under --all; `unknown` is never a function of row order or of the flag." — (NC, `0011:C1`)

## Demand set (`invokedReaders`) — the term that also changes `flow resolve`

- [REQ-16] "The assembled view MUST actually carry the keys the predicate reads, so `internal/cli/flow_exec.go::invokedReaders` MUST add each row's MATCH-block owned keys to its demand set, which today is `Row.RequiresOwned` unioned with the row's GUARD-owned keys only" — (NC, `0011:C1`)

- [REQ-17] "It therefore binds BOTH callers of `invokedReaders` — `flow_next.go` (`invokedReaders(req.model, "")`) and `flow_resolve.go` (`invokedReaders(req.model, outcome)`) — and MUST NOT be scoped to `next`" — (NC, `0011:C1`) — also `0011:BR6`.

- [REQ-18] "The demand set is a property of the MODEL (and, for `resolve`, the requested outcome), not of the mode: it is computed once per invocation, before the row loop, and is identical under --all" — (NC, `0011:C1`)

- [REQ-19] "after this term it invokes the reader, so a run that used to refuse `flow-no-match` MAY now yield a plan, `flow-artifact-missing` (exit 2, the reader's role unbound), or the reader's own refusal (exit 3)." — (NC, `0011:C1`) — the three `flow resolve` arms S8 asserts.

- [REQ-20] "The demand set is a union over the outcome's rows, so the reader is invoked even when the row `resolve` would select does not itself match on the key: over an UNBOUND or REFUSING reader such a run turns from a plan into exit 2/3, before the kernel and above the escape phase — a class this contract NAMES AND ACCEPTS here" — (NC, `0011:C1`) — the accepted behaviour-change class; asserted by S8's BREAKING arm (REQ-116).

- [REQ-21] "No `flow resolve` run over a model WITHOUT such a key changes, and the term is empty over a model with zero owned tags (`0010:C4`)." — (NC, `0011:C1`)

## `unknown` payload — sources, vocabulary, dedup, order

- [REQ-22] "Every match atom over an absent key, every guard atom the kernel could not decide, every owned key no invoked reader established, and (absent --evaluate-gates) every gate id MUST appear in that candidate's `unknown` list as a `{key, reason}` pair." — (NC, `0011:C1`) — the four sources; three emission sites (`0011:A16`).

- [REQ-23] "Match atoms over absent keys and unestablished owned keys carry `absent`, derived from the CLI's own view walk." — (NC, `0011:C1`)

- [REQ-24] "Undecided guard atoms carry the reason the kernel reports — `absent` or `uncomparable`, `0007:C8`'s closed set — read off the probe's `guard_unevaluable` refusal payload (`Refusal.Undecided`); a guard atom over an absent key therefore appears once, the walk and the payload agreeing on `{key, absent}`." — (NC, `0011:C1`)

- [REQ-25] "That payload read MUST be guarded on the probe's disposition being `guard_unevaluable`, the only refusal that carries `Refusal.Undecided` (A13): `summarize` MUST NOT read the field off whatever `Result` it is handed and filter afterwards" — (NC, `0011:C1`)

- [REQ-26] "The gate ids are the ROW's own declared gates — `internal/table/model.go::Row.Gate`, the same field `internal/cli/flow_next.go` already passes to `runGates` when gates do run — not a model-level set and not `0005:C1`'s gate-accessor list" — (NC, `0011:C1`)

- [REQ-27] "A gate id un-run because --evaluate-gates was not passed is neither `absent` … nor `uncomparable` …, so it carries `not-evaluated`, which this RDR mints for the gate class ONLY and which never appears on a match, guard, or owned-key entry." — (NC, `0011:C1`)

- [REQ-28] "The gate-id entries are emitted on EVERY reported disposition, including `owned_state_unavailable`" — (NC, `0011:C1`)

- [REQ-29] "Entries are deduplicated by `{key, reason}` pair, not by key, and sorted by `(key, reason)`." — (NC, `0011:C1`) — sort pinned by S7 (REQ-113).

- [REQ-30] "Deduplication is the LAST step: each source contributes its entries, C2's --all filter is applied to the atom walk's contribution at emission, and only then are the merged entries deduped and sorted." — (NC, `0011:C1`)

- [REQ-31] "when the probe refuses `owned_state_unavailable`, `gate` returns before it collects the guard payload, so an `uncomparable` guard atom on that row is not on the wire and is not reported; every `absent` fact on the row still is" — (NC, `0011:C1`) — the stated fidelity limit; S7 pins it (REQ-112).

## Alphabet and non-selection (both modes)

- [REQ-32] "`outcomes` MUST remain the model's full declared alphabet in every mode." — (NC, `0011:C1`) — **[0005-carried]** (0005 DEV-4's reading stands).

- [REQ-33] "flow next MUST NOT choose among the candidates it reports and MUST NOT turn a gate deny into a refusal." — (NC, `0011:C1`) — **[0005-carried]**.

## Kernel untouched

- [REQ-34] "This contract changes nothing in the `internal/resolve` PACKAGE: the kernel's match seam stays two-valued, its selection and escape phases are untouched, `0007:REQ-78` stays deferred, `0007:C8`'s payload is not widened, and no sixth `RefusalKind` is minted (`0007:REQ-79`)." — (NC, `0011:C1`) — mechanically checkable form is `git diff --stat internal/resolve` empty (`0011:S5`, MVV 9).

- [REQ-35] "The `flow resolve` VERB changes in exactly the one class the demand-set paragraph names, and nowhere else." — (NC, `0011:C1`)

## `--all` (C2)

- [REQ-36] "flow next MUST accept a boolean flag --all, default false." — (NC, `0011:C2`)

- [REQ-37] "Under --all the candidate predicate MUST be exactly the one 0005:C1 specified: every non-escape row whose guard the kernel does not decide false, with the row's match pattern taking no part in the verdict." — (NC, `0011:C2`)

- [REQ-38] "A match atom is then neither an exclusion nor an entry in `unknown`, whatever its key's presence, achieved by the probe omitting the match pattern so the kernel never sees those atoms." — (NC, `0011:C2`)

- [REQ-39] "Under --all that entry MUST disappear — reporting a fact the mode ignores would be reporting on a predicate it does not apply — so the --all branch MUST filter BlockMatch atoms out of `unknown`, and an oracle MUST assert their absence." — (NC, `0011:C2`) — the one respect in which `--all` is NOT output-identical to 0005; outside A4's classification.

- [REQ-40] "The filter is scoped to the ATOM WALK's contribution ONLY. It MUST NOT suppress a `{key, absent}` pair the OWNED-key walk independently produces for the same key" — (NC, `0011:C2`)

- [REQ-41] "That scope is only achievable at EMISSION: the filter MUST be applied to the atom walk as it contributes, before the sources merge and before C1's dedup, never as a predicate over the merged list." — (NC, `0011:C2`)

- [REQ-42] "The oracle asserting the absence MUST therefore be written over a match key that is NOT in the row's `RequiresOwned`" — (NC, `0011:C2`)

- [REQ-43] "Everything else 0005:C1 requires of flow next — the data minimum, gate accessors only under --evaluate-gates and only for reported candidates, a deny reported on its candidate with exit 0, no invented guard facts, the narrowed invoked reader set — MUST hold identically in both modes." — (NC, `0011:C2`) — **[0005-carried]**.

- [REQ-44] "The success payload shape MUST be identical in both modes: `unknown` is present in both, as `[]` rather than omitted when it carries nothing, so one consumer struct parses either mode." — (NC, `0011:C2`)

- [REQ-45] "--all is part of request identity." — (NC, `0011:C2`) — with `0011:D-identity`.

- [REQ-46] "--all MUST NOT be accepted by flow resolve, flow read-state, or flow set-state — satisfied by NON-REGISTRATION (registered on next only, not on the shared registerSelectionFlags)." — (NC, `0011:C2`)

- [REQ-47] "The build MUST pin that negative STRUCTURALLY, asserting --all is absent from those three verbs' flag sets and from the flow group's own, guarded against vacuity by requiring the four verbs to exist first" — (NC, `0011:C2`) — the idiom is `TestReq69_NoPlanFlagShipsOnAnyVerb` / `TestReq12_FlowGroupDoesNotRedefineTheRootAsFlag`.

- [REQ-48] "`command-error` + exit 2 is the SHARED usage bucket for every unknown flag, unknown subcommand, and Args violation, so it MUST NOT be asserted as if it named this flag" — (NC, `0011:C2`) — no oracle may match on pflag's text.

- [REQ-49] "This contract mints no new refusal code for it." — (NC, `0011:C2`)

- [REQ-50] "the flag is `--all`, boolean." — (LBD, `0011:D-naming`) — exact spelling normative; `--enumerate`, `--ignore-match`, `--no-match`, a sense flip, and an enum (`--state=all|applicable`) are named rejected.

## Help, docs, and the test re-homing (C3)

- [REQ-51] "The command's Short and Long help MUST state that flow next reports the candidates the supplied state can take (match and guard both holding or undecided), that a match or guard key the state does not carry leaves a row a candidate with the key listed under `unknown` and its reason named, and that --all reports every row the guards do not exclude regardless of match." — (NC, `0011:C3`)

- [REQ-52] "The help MUST also state that a candidate is a row the supplied state does NOT EXCLUDE, not a row flow resolve will select: a candidate carrying an `unknown` entry may still be refused by flow resolve over the same state, and the entry names what to supply." — (NC, `0011:C3`)

- [REQ-53] "The 0005 flow next tests in internal/cli/flow_next_0005_test.go MUST keep passing unchanged where they assert the alphabet, gate handling, reader narrowing, determinism, and non-mutation; an assertion that depended on a match-excluded row being reported MUST be re-homed under --all, not deleted, and the file's header comment MUST name this RDR as the source of the default." — (NC, `0011:C3`)

- [REQ-54] "The build MUST mechanically re-home every shipped read of the `unresolved` key to `unknown` and its element type from string to `{key, reason}` — FIVE TEST reads on this build: `internal/cli/flow_next_0005_test.go` (3 — the presence check `:103`, the gate-id assertion `:163`, and the adversarial guard-fact assertion `:404`), `internal/cli/flow_adversarial_0005_test.go` (1, `:446`), and `internal/cli/flow_mvv_0005_test.go` (1, `:66`)." — (NC, `0011:C3`)

- [REQ-55] "the production rename additionally touches the `candidate.Unresolved` field and its `json:"unresolved"` tag, the field's doc comment, the file header comment naming `unresolved`, `summarize`'s `Unresolved: []string{}` initializer, its three `c.Unresolved = append(…)` assignments and the `slices.Contains(c.Unresolved, …)` guards beside them, and — the shipped USER-FACING string in this file, not a comment — the cobra `Long` help at `internal/cli/flow_next.go:70-71`" — (NC, `0011:C3`)

- [REQ-56] "Two further shipped user-facing descriptions state the OVERRIDDEN default and MUST be corrected with the help text, though neither reads the payload field: `docs/cli-output-contract.md:125` ("Enumerate the legal outcomes and their candidate rules") and `README.md:43` ("list legal next outcomes")." — (NC, `0011:C3`)

- [REQ-57] "The gate-id reads re-home to the `not-evaluated` reason C1 names; a re-homing that has to invent a reason token is not mechanical and is a defect." — (NC, `0011:C3`)

- [REQ-58] "THREE of the five discard the comma-ok of a helper (`flow_harness_0005_test.go::stringsAt` …) — `flow_next_0005_test.go:404`, `flow_mvv_0005_test.go:66`, and `flow_adversarial_0005_test.go:446` — and the last of those is a NEGATIVE assertion (`if containsString(unresolved, "flag")`) that passes vacuously on the empty slice, so a re-homed read MUST assert the ok bool, and a green suite MUST NOT be cited as evidence that this contract was implemented." — (NC, `0011:C3`)

- [REQ-59] "`stringsAt` itself MUST NOT be repointed at the new element type … The re-homing adds a SIBLING helper beside it (a `[]{key, reason}` reader over the same payload) and switches the five reads to it." — (NC, `0011:C3`)

- [REQ-60] "The same rename MUST also sweep the COMMENT and failure-message prose that names `unresolved` … On this build that is the doc and header comments and the `t.Errorf`/`t.Fatalf` strings at `internal/cli/flow_next_0005_test.go:23,66,104,135,137,141,169-171,404-408`, `internal/cli/flow_mvv_0005_test.go:46,73`, and `internal/cli/flow_adversarial_0005_test.go:448,450`." — (NC, `0011:C3`)

## Fixtures the build owes (C3)

- [REQ-61] "The build MUST therefore add fixtures that discriminate: a row whose match key is present and UNEQUAL (excluded by default, reported under --all; under --evaluate-gates its gates do not run), and a row whose match key is ABSENT (candidate in both modes, `{key, absent}` by default and no match entry under --all)." — (NC, `0011:C3`)

- [REQ-62] "That absent-key row MUST match on a key it does NOT write or clear, writing some OTHER key instead" — (NC, `0011:C3`) — else the owned-key walk produces the same pair in both modes and the fixture discriminates nothing.

- [REQ-63] "It MUST also add the three pins that license C1's presence test: a model declaring one owned key served by two readers is refused at load (`internal/table/load.go::checkAccessorBindings`); a repeated --tag key is refused `flow-tag-duplicate`; and a --tag on a key the model declares owned is refused `flow-tag-owned` — all before any probe is built." — (NC, `0011:C3`)

- [REQ-64] "There is no conflicted-key fixture for flow next: the case is unreachable (A3), and a fixture that reaches the kernel's conflicted fold by constructing a `TagSet` in-package tests `internal/resolve`, which this RDR does not change." — (NC, `0011:C3`) — a negative obligation.

- [REQ-65] "It MUST also add C2's structural negative … and register --all on next in the per-verb flag map the surface suite already carries (`flow_surface_0005_test.go::TestReq3_EachVerbRegistersItsNormativeFlagSpellings`). These are TWO oracles, not one" — (NC, `0011:C3`)

- [REQ-66] "The `uncomparable` fixture S7 needs is a present value the operator cannot parse (A13's one reachable producer): the other three producers are unreachable and MUST NOT be fixtured" — (NC, `0011:C3`)

## Data flow and signature (Technical Design prose)

- [REQ-67] "`excluded` must return the kernel's `Result` (not a `bool`) so `summarize` can read the `guard_unevaluable` payload; this is the one signature change in the file" — (TD, line 113) — with `0011:EIA`'s *Per-row kernel probe* row.

- [REQ-68] "a row whose retained match pattern is empty matches unconditionally as it does today" — (TD, line 113).

- [REQ-69] "The alphabet `outcomes` is copied from the model in both modes (0005's DEV-4 reading stands: the alphabet is what may be requested; the candidates are what varies with state)." — (TD, line 111).

- [REQ-70] "gates under `--evaluate-gates` for reported candidates only" — (TD, line 111) — **[0005-carried]**, and the match-excluded row's gates MUST NOT run (REQ-12).

## Payload rendering and text mode (Load-Bearing Decisions)

- [REQ-71] "the per-candidate field is `unknown`, a list of `{key, reason}`" — (LBD, `0011:D-undecided-reporting-shape`).

- [REQ-72] "Text mode (`--as=text`) carries the same pairs — both members of each, so the text caller gets the diagnosis and not just the symptom" — (LBD, `0011:D-undecided-reporting-shape`).

- [REQ-73] "It does NOT get a bespoke `key (reason)` form: `flow next` declares no `FindingCarrier` and `respond.OK` therefore renders its payload through the shared generic flattener … which emits path-qualified leaf lines — `candidates[0].unknown[0].key: stage`, `candidates[0].unknown[0].reason: absent`." — (LBD, `0011:D-undecided-reporting-shape`) — no per-verb text template may be minted; "no scenario asserts a text-mode string".

- [REQ-74] "In both modes `next` reports all survivors and chooses none; exactly-one selection and refusal remain `flow resolve`'s (`0005:D-selection-predicate`)." — (LBD, `0011:D-selection-predicate`).

- [REQ-75] "reuses `0007:C8`'s closed reason set (`absent`, `uncomparable`) for every fact that comes off the kernel's guard payload, rather than minting one there; the CLI's own walk emits only `absent`." — (LBD, `0011:D-undecided-vocabulary`).

- [REQ-76] "the same request in the same mode over the same model revision and artifact contents reports the same candidate set." — (LBD, `0011:D-identity`).

## Disposition grid (mini-check — one oracle row per reachable input class)

- [REQ-77] "match key present, equal | yes | none | yes | 0" — (MC, `disposition`).

- [REQ-78] "match key present, unequal | **no** (excluded) | n/a — not reported | **no** | 0" — (MC, `disposition`).

- [REQ-79] "match key absent from view | yes | `{key, absent}` | yes | 0" — (MC, `disposition`).

- [REQ-80] "dead row … yes while the key is absent (undecided); **no** once it is present at any value" — (MC, `disposition`).

- [REQ-81] "gate id, `--evaluate-gates` not passed | yes | `{gate-id, not-evaluated}` | no | 0" — (MC, `disposition`).

- [REQ-82] "only match atom is `recognized` | yes | none — lifted into `Row.Outcome` at normalize | yes | 0" — (MC, `disposition`).

- [REQ-83] "guard decided false | **no** (excluded) | n/a | no | 0" — (MC, `disposition`) — **[0005-carried]**.

- [REQ-84] "guard atom over present key, seam unevaluable | yes | `{key, uncomparable}` (payload only) | yes | 0" — (MC, `disposition`).

- [REQ-85] "owned key no reader established | yes | `{key, absent}`; an `uncomparable` guard atom on the same row is not reported (precedence limit) | yes | 0" — (MC, `disposition`).

- [REQ-86] "any of the above, under `--all` | match takes no part; guard rules alone decide | no match entry, whatever the key's presence | per guard | 0" — (MC, `disposition`).

## Minimum Viable Validation (`0011:MVV`)

- [REQ-87] "Assert: `candidates[]` is exactly `prelock`, `resolve-route-back`, `resolve-abandon` (3 rows, one per outcome), and no rule whose `match.stage` names another value; `outcomes[]` is still the model's full alphabet." — (MVV 3, `0011:MVV`).

- [REQ-88] "Re-run with `--all`; assert `candidates[]` has the 21 rows the baseline reported (the 22nd, `finalize-pass`, is guard-excluded on `gate_passed`, not match-excluded) and the same `outcomes[]`." — (MVV 4, `0011:MVV`).

- [REQ-89] "Run over the 0005 gated fixture (`flowGatedNextModel` …) with and without `--all`; assert the `gated-excluded` row is absent in both (guard-excluded), paired with its negative control `gated-reported` present in both. Then over `flowMVVModel` (`status=draft`), assert the row C3 adds … is absent by default and present under `--all`." — (MVV 5, `0011:MVV`).

- [REQ-90] "Run over a model exercising the three reachable match cases … and a dead row … Then add `--tag key=<value>` and assert the list narrows to the rows whose match holds. Assert the three pins" — (MVV 6, `0011:MVV`).

- [REQ-91] "Assert `unknown` is present as `[]` rather than omitted on a candidate with nothing undecided, in both modes — the fixture must SUPPLY one … and which is run WITH `--evaluate-gates`" — (MVV 7, `0011:MVV`).

- [REQ-92] "Run over a model declaring an owned key some row MATCHES on but no row writes, clears, or guards (A15): assert its reader appears in `readers`, the key is in the view, and the row is match-decided … assert the same `readers` set under `--all`." — (MVV 8, `0011:MVV`).

- [REQ-93] "Run `flow resolve --outcome <o>` over the same model with the reader bound, unbound, and refusing: assert a plan, `flow-artifact-missing` (exit 2), and the reader's refusal (exit 3) respectively, and `flow-no-match` for that key in none; assert `flow resolve` output over `models/rdr.toml` and every shipped fixture is unchanged." — (MVV 8, `0011:MVV`).

- [REQ-94] "`make check` passes; `go test ./internal/resolve` passes with no test changed and `git diff --stat internal/resolve` empty; the 0005 `next` suite passes with only C3's five mechanical `unresolved`→`unknown` re-homings (ok-bool asserted), and the file header names this RDR." — (MVV 9, `0011:MVV`).

## Testing Strategy S1–S8

- [REQ-95] "S1 … `candidates[]` is exactly `prelock`, `resolve-route-back`, `resolve-abandon`; `outcomes[]` is the full alphabet." — (TS, `0011:S1`).

- [REQ-96] "`gated-excluded` absent in both **over `flowGatedNextModel`, the fixture that authors it** … with `gated-reported` present in both as its negative control" — (TS, `0011:S2`) — asserting it over `flowMVVModel` "passes vacuously in every possible build and pins nothing".

- [REQ-97] "The filter/dedup ORDER is discriminated here: over a matched-AND-written key whose reader established nothing … `--all` MUST still report `{stage, absent}` from the owned-key walk, since only the atom walk's contribution is filtered (C2) and dedup runs last (C1)." — (TS, `0011:S2`).

- [REQ-98] "C2's negative is asserted STRUCTURALLY … not on pflag's error text; the behavioural run (exit 2, `command-error`, envelope on stdout under `--as=json`) corroborates it" — (TS, `0011:S2`).

- [REQ-99] "present-equal ⇒ candidate, **the match key** absent from `unknown`; present-unequal ⇒ excluded and its gates do not run under `--evaluate-gates`; absent ⇒ candidate with `{key, absent}`; a dead row … ⇒ candidate with `{key, absent}` while absent, excluded once the key is present at any value, with no CLI literal comparison (A12)." — (TS, `0011:S3`).

- [REQ-100] "the oracle asserts two match tags on the DEAD row only — the pairing's other expanded row collapses to one tag and is live, so asserting two there would assert something false." — (TS, `0011:S3`).

- [REQ-101] "The assertion is therefore on the MATCH key's absence from `unknown`, never on `unknown` being empty" — (TS, `0011:S3`).

- [REQ-102] "(a) A SET-KINDED match key … assert the row matches when the supplied set equals the authored one under the kernel's canonical form." — (TS, `0011:S3`) — the one fixture that separates a filter-before-conversion build (REQ-6).

- [REQ-103] "(b) A row carrying SEVERAL match keys in MIXED states at once — one present-and-equal, one present-and-unequal, one absent — asserting the row is excluded (the unequal key decides, by the kernel's conjunction) and that under `--all` it is a candidate whose `unknown` carries no match entry." — (TS, `0011:S3`).

- [REQ-104] "the `recognized`-only row is a candidate with `unknown` empty of match facts; the all-absent row is a candidate — its probe carries an empty match pattern, which matches unconditionally — and every omitted key is listed `{key, absent}`, so its `unknown` is NOT empty." — (TS, `0011:S4`).

- [REQ-105] "An AUTHORED empty `[rule.match]` is not a fixture here and MUST NOT be attempted: `internal/table/normalize.go` refuses it at load … the empty match pattern C1 names is minted by the probe builder at runtime, never authored." — (TS, `0011:S4`).

- [REQ-106] "a model with one owned key served by two readers fails to load on `checkAccessorBindings`' message …; `flow next --tag k=a --tag k=b` is refused `flow-tag-duplicate` and `--tag <owned-key>=v` is refused `flow-tag-owned`, with no probe built" — (TS, `0011:S5`).

- [REQ-107] "the kernel is untouched — the ONE mechanically checkable form of that claim is MVV 9's `git diff --stat internal/resolve` empty, and it is the form the build asserts" — (TS, `0011:S5`).

- [REQ-108] "`TestReq78_MatchPatternStillFoldsAbsenceIntoNonMatch`, the `match_conflicted_test.go` oracles, and the reason-set pins `TestReq49`/`TestReq47` … stay green as shipped." — (TS, `0011:S5`).

- [REQ-109] "A11's key-set agreement is PINNED here …: an oracle MUST compare `assembledView`'s key set against the kernel's `internal/resolve/resolve.go::assemble` over the same inputs, modulo the kernel's `recognized` key" — (TS, `0011:S5`).

- [REQ-110] "the `flow next` oracles pass, none re-homed under `--all`; the 5 test reads of `unresolved` are re-homed to `unknown` with the ok-bool asserted; file header names this RDR. A green 0005 suite is NOT evidence C1 shipped — S2/S3 are." — (TS, `0011:S6`).

- [REQ-111] ""None re-homed under `--all`" is a claim about the DIFF, not a runtime assertion — no test can observe it — so it is discharged at review by reading the changed test files" — (TS, `0011:S6`) — a review-time obligation, not an oracle.

- [REQ-112] "a guard atom over a present key the seam cannot decide appears as `{key, uncomparable}` … the same atom over an absent key appears once as `{key, absent}` (walk and payload agree, dedup on the pair); on a row that also lacks an owned key the probe refuses `owned_state_unavailable`, the row is still a candidate carrying the missing owned key and its `absent` facts, and the `uncomparable` guard atom is not reported — asserted so the limit is pinned rather than discovered." — (TS, `0011:S7`).

- [REQ-113] "`unknown` is SORTED by `(key, reason)` — asserted on a fixture whose facts are minted out of that order, so an unsorted build fails." — (TS, `0011:S7`).

- [REQ-114] "the reader IS invoked (it appears in the payload's `readers`), the key IS in the view, and the row is match-decided rather than reported `{key, absent}`." — (TS, `0011:S8`).

- [REQ-115] "Negative controls: the same model under `--all` reports the SAME reader set …; and a model with zero owned tags demands nothing and invokes no reader, unchanged (0010's `decision-table` class, `0010:C4`)." — (TS, `0011:S8`).

- [REQ-116] "The BREAKING arm: add a second ordinary row for the same outcome that does not match on the key; with the reader unbound, assert the pre-extension build returns that row's plan and the post-extension build `flow-artifact-missing` — pinned so the DEV-8 class is a known cost, not a surprise." — (TS, `0011:S8`).

- [REQ-117] "Both verbs' `readers` over the same model and outcome are EQUAL (one `invokedReaders`, one view), asserted on a fixture whose reader serves several keys." — (TS, `0011:S8`).

- [REQ-118] "`flow resolve` over `models/rdr.toml` and every shipped fixture is byte-identical before and after (A15 iii)." — (TS, `0011:S8`) — "the build owes the payload-diff form the spike could not reach".

- [REQ-119] "Done: S1–S8 green and `make check` passes." — (TS, line 793).

## Phase intents (Implementation Plan)

- [REQ-120] "in `flow_exec.go::invokedReaders`, add each row's match-block owned keys to the demand set — a term both callers take … — the one edit outside `flow_next.go`" — (PH 1, line 759).

- [REQ-121] "register `--all` on `next` only, route it to a probe with the match pattern omitted and the `BlockMatch` atoms filtered out of `unknown` — the ATOM WALK's contribution only, leaving the owned-key walk's entries intact (C2) — and rewrite Short/Long help per C3." — (PH 2, line 763).

- [REQ-122] "document `flow next`'s payload — including `unknown` and the `--all` invocation — in `docs/cli-output-contract.md`, which currently shows only the invocation grammar" — (PH 3, line 767).

## Failure modes and cross-cutting (assertable residue)

- [REQ-123] "`next` is effect-free in both modes" — (FM, `0011:F7`) — non-mutation, **[0005-carried]**, holds under `--all` too.

- [REQ-124] "Only guard atoms can carry this reason on `next`." — (FM, `0011:F3`) — `uncomparable` never appears on a match, owned-key, or gate entry.

- [REQ-125] "Version marker: none minted — the payload field rename `unresolved` → `unknown` IS the break" — (CC, `0011:G-cross-cutting`) — no schema/version field may be added.

- [REQ-126] "The term adds keys to a demand set computed once per invocation, before any probe, and invokes readers through the existing `internal/cli/flow_exec.go::runReaders` path with its shipped refusal arms (`codeArtifactMissing` exit 2, `accessorFailure` exit 3). No new execution path, no new concurrency surface" — (CC, `0011:G-cross-cutting`).

- [REQ-127] "the CLI compares no value, so no case folding, normalization, or collation decision is taken by this RDR. Set literals go through `seamValue`'s existing canonicalization." — (CC, `0011:G-cross-cutting`).

## Override scope (Metadata / Overrides)

- [REQ-128] "Every other `0005:C1` obligation on `flow next` is unchanged." — (OVR, line 20) — the boundary of the override; only the candidate clause moves.

- [REQ-129] "Renames the candidate's `unresolved` list to `unknown` and types its entries `{key, reason}`; adds the reason token `not-evaluated` for un-run gate ids on that CLI list only — `0007:C8`'s seam vocabulary stays at `absent`/`uncomparable` and its payload is not widened." — (OVR, line 20).

- [REQ-130] "Touches no clause of RDR 0007 or JDR 0001: `0007:REQ-78` (the two-valued match seam) is left deferred" — (OVR, line 20).

---

## ASSUMPTIONS

Implicit choices made where the wording is imprecise but one reading is
defensible. Each names the evidence it rests on.

- **A-1 — REQ-4's "the view HOLDS A VALUE" includes the empty string.** C1 says
  "Present means the view HOLDS A VALUE for the key" without saying whether
  `""` counts. `0011:A11`'s evidence settles it: "an empty-but-present value
  yields a `Tag` with `Value: ""`, which `assembledView` stores as a present map
  entry (`_, known := view[atom.Key]` is true) and the kernel stores as a
  present `taggedValue` — agreement, not divergence." Read as: presence is map
  membership (`_, ok := view[k]`), never a non-empty test. A non-empty test
  would also break REQ-109's key-set agreement oracle.

- **A-2 — the `unknown` entry's JSON member names are `key` and `reason`.** C1
  writes the pair as `{key, reason}` throughout and never gives a struct tag,
  but `0011:D-undecided-reporting-shape` fixes the text rendering as
  `candidates[0].unknown[0].key: stage` / `candidates[0].unknown[0].reason:
  absent`, and that rendering is the shared flattener over the JSON payload's
  own field names. Read as: the wire members are exactly `key` and `reason`,
  both always emitted (no `omitempty`), on a list named `unknown`.

- **A-3 — the reason values are the literal lowercase-hyphen tokens `absent`,
  `uncomparable`, `not-evaluated`.** The record spells them that way in every
  occurrence (C1, `0011:D-undecided-vocabulary`, the `disposition` table, F2/F3)
  and `0007:C8`'s shipped set is `internal/resolve/guard.go::Reason` with those
  two spellings. Read as: string values, not an enum type minted in
  `internal/cli`, and `not-evaluated` is spelled with a hyphen.

- **A-4 — the `--all` `BlockMatch` filter keys on the ATOM's `Block` field, not
  on a key-set difference.** C2 says "the --all branch MUST filter BlockMatch
  atoms out of `unknown`" and REQ-40/REQ-41 forbid a post-merge key-keyed
  predicate. `0011:A16` establishes the atom walk is a distinct emission site
  (`for _, atom := range row.Atoms`, `flow_next.go:186`), where `atom.Block` is
  in hand. Read as: inside the walk, skip an atom whose `Block == BlockMatch`
  when `--all` is set; the owned-key walk and the gate-id loop are untouched.

- **A-5 — the gate-id entry's `key` member carries the gate id itself.** The
  `disposition` table writes the pair as `{gate-id, not-evaluated}` and C1 calls
  the source "the row's gate ids", `internal/table/model.go::Row.Gate` — a
  `[]string`. Read as: the gate id string goes in the `key` member, sharing the
  namespace with view keys. `0011:S7` confirms the namespaces cannot collide in
  a way that matters: "`not-evaluated` is scoped to gate ids, a separate
  namespace from view keys", so no disambiguating field is owed.

- **A-6 — the `--all` probe omits the match pattern entirely rather than
  filtering it by presence.** C2: "achieved by the probe omitting the match
  pattern so the kernel never sees those atoms." Read as: under `--all`,
  `probe.Match = nil` (today's shipped shape), not a presence-filtered set —
  which is also what makes REQ-37's "exactly the one 0005:C1 specified"
  literally true.

- **A-7 — "no escape list" on the probe means `probe.Escape = nil`, as
  shipped.** C1 requires a probe with "no escape list" and does not name the
  mechanism; `0011:A1`'s evidence describes the shipped `excluded` as already
  doing this ("the stripped escape list leaves `escapeOrRefuse` no rescue row").
  Read as: the existing strip is retained unchanged in both modes; only
  `probe.Match` changes.

- **A-8 — `Outcomes: []string{row.Outcome}` is retained alongside `Recognized:
  row.Outcome`.** C1 names only the `Recognized` binding, but `0011:A1` states
  the pairing outright: "`internal/cli/flow_next.go::excluded` binds `Outcomes:
  []string{row.Outcome}` alongside `Recognized: row.Outcome`, so `Table.models`
  holds trivially — C1's probe builder MUST preserve that pairing." Read as: a
  MUST on the builder, inherited from A1's evidence rather than restated in the
  fence.

- **A-9 — the `--all` flag has no shorthand.** `0011:D-naming` fixes the
  spelling as `--all` and `0011:A14` records that `-a` produces different pflag
  text ("the `-a` shorthand yields the same code and exit but different message
  text"), observed as an unregistered flag. Read as: register the long form
  only; a `-a` shorthand is not authorized and would widen the surface REQ-47
  pins.

- **A-10 — "no rule whose `match.stage` names another value" (MVV 3) is an
  assertion about the reported set, not about the model file.** Read as: the
  three named rules are the whole of `candidates[]` and every other rule of the
  22 is absent, which is the discriminating form (a stripped-match build
  reports 21).

- **A-11 — REQ-118's "byte-identical before and after" is discharged against a
  pre-change build of the same tree, not against a committed golden.** `0011:S8`
  says the A15 spike compared "the two binaries (`cmp`)" and that "the build
  owes the payload-diff form the spike could not reach". Read as: the oracle may
  be a golden payload checked into the fixture set, provided it is captured
  from the pre-change behaviour and asserts the FULL payload rather than a
  subset — which is the gap S8 names ("a fixture asserting only a subset of its
  payload being exactly what suite-green cannot catch").

- **A-12 — Q1's resolution (below): the `--all` filter removes an entry a
  match atom would have contributed even when a GUARD atom on the same key
  would independently contribute `{key, absent}`.** See QUESTIONS Q1; the
  chosen reading is per-source, not per-key, grounded in REQ-40's own wording
  ("the OWNED-key walk independently produces") and `0011:A16`'s three-site
  separation.

- **A-13 — Q2's resolution (below): `--all` does not re-widen the demand set
  back to 0005's.** See QUESTIONS Q2; grounded in REQ-18 and CONS line 711
  ("`--all` on `next` restores 0005's PREDICATE, not 0005's reader set: on this
  class it invokes one more reader than 0005 did").

---

## QUESTIONS

Two clauses admit readings that would produce materially different behaviour.
This run is unattended; each is resolved below on the most defensible reading,
recorded as an ASSUMPTION, and the alternative left visible so a later phase can
overrule.

- **Q1 — Under `--all`, is the atom-walk filter scoped to MATCH-block atoms
  only, or to every atom of the walk over a key some match atom names?** C2
  says "the --all branch MUST filter BlockMatch atoms out of `unknown`", which
  is per-ATOM; but the walk is per-atom over `row.Atoms` and a row may carry
  BOTH a match atom and a guard atom on one key. Reading (a): filter per atom on
  `atom.Block == BlockMatch`, so a GUARD atom on the same absent key still
  contributes `{key, absent}` under `--all`. Reading (b): filter per key, so any
  key a match atom names is suppressed wholesale. **Proceeding on (a).** It is
  the literal wording ("BlockMatch atoms"); REQ-40 already establishes the
  contract reasons per SOURCE and forbids suppressing a pair another source
  independently produces ("It MUST NOT suppress a `{key, absent}` pair the
  OWNED-key walk independently produces for the same key"), and the guard walk
  is the same shape of independent contribution; C2's own ground is "reporting a
  fact the mode ignores would be reporting on a predicate it does not apply" —
  and `--all` does NOT ignore guards, it decides on them alone (REQ-37). Reading
  (b) would delete a guard fact `--all` is required to report and would violate
  REQ-43's "everything else 0005:C1 requires … MUST hold identically in both
  modes". Recorded as ASSUMPTION A-12.

- **Q2 — Does `--all` restore 0005's invoked reader SET as well as its
  predicate?** REQ-37 says the `--all` predicate is "exactly the one 0005:C1
  specified", which could be read to include 0005's demand set (no match-owned
  term), and REQ-43 requires "the narrowed invoked reader set" to hold
  identically in both modes — ambiguous as to WHICH narrowed set. Reading (a):
  the demand set is mode-independent and carries the match-owned term in both
  modes. Reading (b): `--all` reverts `invokedReaders` to `RequiresOwned` ∪
  guard-owned. **Proceeding on (a).** C1 states it directly — the demand set "is
  identical under --all (where the match pattern takes no part in the verdict
  but its keys are still read)" (REQ-18) — S8's negative control asserts "the
  same model under `--all` reports the SAME reader set" (REQ-115), and
  Consequences says it in as many words: "`--all` on `next` restores 0005's
  PREDICATE, not 0005's reader set: on this class it invokes one more reader
  than 0005 did" (line 711). Reading (b) would make the view mode-dependent and
  break REQ-15's "the invoked reader set and the assembled view are fixed ONCE
  per invocation … and are the same under --all". Recorded as ASSUMPTION A-13.
