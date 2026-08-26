Model: claude-fable-5

# Premortem critic — RDR 0011, iter-2 (draft-free, brief-only)

Independence note: this critique was written from the brief alone. No RDR, no source, no search. Every claim below about kernel or CLI behaviour is therefore a *hypothesis the draft must refute with a cited line or a named test*, not an assertion about the code.

## Findings ledger

| ID | passage/claim | Failure mode | Symptom user sees | Origin |
|---|---|---|---|---|
| P-1 | N4 "changes `flow resolve` only on the match-only-owned-key class — where every existing run refuses `flow-no-match`" | The demand set is a UNION over all rows, not over the selected row. On a model with a match-only owned key K, runs that select a row NOT matching on K work today (K's reader is never invoked). After the change K's reader is demanded on every invocation. | `flow resolve --artifact rdr=x` that produced a plan yesterday now exits 2 `flow-artifact-missing` (K's role unbound) or exits 3 (K's reader refuses) — with no rule involving K anywhere near the selected row. | refutation-target (S2 class: "changes only on class Y" asserted from prose, not from the union in `invokedReaders`) |
| P-2 | N4 "strictly additive ... nothing that worked stops working" | Additive on the reader SET is not additive on the OUTCOME: a newly invoked reader can refuse (exit 3), an unbound role can refuse (exit 2), and a newly read value can fail owned-state validation in the kernel. More readers = more veto points, before the escape phase runs. | Escape rows that rescued `no_match` yesterday never get the chance: refusal happens above the escape switch. | refutation-target (S2 class exactly: escape reachability) |
| P-3 | (b)1 / problem statement: ABSENT = "neither owned by an invoked reader nor supplied as `--tag`" | Presence is defined by *reader invoked*, not *reader returned a value*. A freshly seeded artifact whose front matter has no `stage:` yet has the key "present" (reader ran) but valueless; kernel `TagSet.matches` folds it into false → row excluded, silently. | The most common journey — `flow next` on a new artifact — returns an EMPTY candidate list (the exact R1 failure the approach claims to avoid) instead of every seed-stage row with `{stage, absent}`. | negation (N1/N6) |
| P-4 | N1 "a `--tag` on an owned key is refused" | Reachability/order: if the `--tag`-vs-owned refusal is checked against the *invoked* reader set rather than the model's *owned* set, a `--tag` on a match-only owned key was accepted yesterday (reader not invoked) and today collides with the newly invoked reader — either refused (behaviour change) or, worse, both land in the view → present-and-conflicted IS reachable. | Skill that passed `--tag approved=yes` as a workaround for the unread key now gets exit 2, or gets a candidate list computed over a conflicted key. | refutation-target (S1/S2 class: peer-check semantics asserted without stating which set the check runs against) |
| P-5 | N7 "omitting absent-key tags changes only whether the row reaches the kernel's `gate`" | A row whose ONLY match key is absent yields a probe with an EMPTY match list. If the kernel classifies empty-match rows as escape/default rows, or validates "row must match on ≥1 key", the filtered probe takes a different phase (escape) or a table-validation refusal — not `gate`. | `flow next` reports the row under an `escape` disposition, or errors on table validation, on a query with no `--tag`. | refutation-target (S2 class: phase/reachability from prose) |
| P-6 | N2 "the filter is over the CONVERTED probe (`row.KernelRow()`)" | If `KernelRow()` materialises match tags into a map keyed by `Key` (dedupe) or into a canonical value set, the `eq`+`in` two-tag case collapses before the filter — dead rows become live rows. The proposal fixed S3 in the *CLI* filter but relies on the *converter* preserving cardinality, which is the same S3 assumption one hop upstream. | A dead row (`stage eq a` + `stage in [b]`) is selected by `flow resolve` when stage=a, silently. | refutation-target (S3 class, moved one producer up) |
| P-7 | N6 "the CLI's view and the kernel's view hold the same key SET" | Two assembly paths. A reader that serves several owned keys (front-matter reader serving `stage` AND `status`) is invoked for `stage`; whether `status` enters the CLI view, the kernel view, both, or neither depends on whether each path filters reader output to demanded keys. | A match on `status` is reported `{status, absent}` by `next` but excluded/selected by `resolve`, or vice versa — the exact `next`/`resolve` divergence N5 was written to prevent. | negation (N6), refutation-target (cardinality claim about two data structures asserted from prose) |
| P-8 | N8 "`--all` yields exactly 0005's predicate; the two modes share reader set" | The demand set grows for `--all` too. A 0005 `flow next` over a model with an unbound match-only key role listed 21 rows; `--all` now exits 2 on the same inputs. | "`--all` restores 0005" is false on precisely the models the feature targets. | negation (N8) |
| P-9 | N3 "dead row needs no CLI handling; lint (RDR 0006) is the remedy" | RDR 0006 is unimplemented. Until it ships, every no-`--tag` query lists every dead row — and `in`-expansion multiplies them (one dead row per expanded literal). | User sees N phantom candidates, each carrying only `{key, absent}`, indistinguishable from live rows. | negation (N3) |
| P-10 | R1 rejection: "partially supplied state yields an empty list with nothing to act on" | Under the chosen approach the *no-tag* query over `models/rdr.toml` still yields ~21 candidates each with `{stage, absent}` — the enumeration defect, relabelled. The approach only outperforms R1 when the state IS supplied, where R1 gives the identical answer. | The motivating symptom ("21 of 22") survives the fix for the no-tag journey. | negation (R1) |
| P-11 | (b)4 rename `unresolved` → `unknown [{key, reason}]` | Payload shape change with no version signal; skills reading `unresolved` get nil and treat every row as fully resolved. | Skill silently treats un-run gates as passed. | negation (N8 payload-shape) |
| P-12 | (b)1 "the CLI decides PRESENCE only, the kernel decides EQUALITY" | Kernel `no_match` is also the fold for *conflicted*; with P-4 open, `next` maps a conflict to "present-and-unequal → excluded" with no diagnostic. | Row vanishes from the list; `unknown` is empty; nothing explains why. | negation (N1) |
| P-13 | R4 rejection cites "the payload carrier for match facts does not exist" | Correct today, but the proposal's own `unknown` field IS a carrier for match facts built CLI-side; R5's rejection ("a fact the CLI already has") is contradicted by P-3/P-12: the CLI does NOT have value-absent or conflict facts. | — (design-consistency finding) | negation (R4/R5) |

## 1. Prospective hindsight — the failure narrative

RDR 0011 shipped in week 35. Match-conditioned `flow next` worked on every fixture and on `models/rdr.toml`. The first external model to hit it was a team's release model with an `approvals` role: one owned key, `approved`, served by a sign-off reader and matched on by exactly two of fourteen rows. Nobody wrote, cleared, or guarded it — a match-only owned key, the very class the RDR named.

Their `flow resolve --artifact rdr=…` runs had been green for weeks; the rows they selected never touched `approved`, so the reader was never demanded. After the upgrade, `invokedReaders` unioned `approved` into the demand set on every invocation. Every run without `--artifact approvals=…` exited 2 `flow-artifact-missing`. The incident ticket quoted the RDR: "nothing that worked stops working." It had, because the RDR's "class" was defined by the *selected row* while the demand set is a union over *all candidate rows*.

They bound the role. Now the sign-off reader ran on artifacts that had no sign-off section yet and refused (exit 3). Two escape rows that had rescued `no_match` for the pre-approval stages never executed: the reader refusal happened above the escape switch. That was the S2 failure, verbatim, one layer up.

Meanwhile a skill author ran `flow next --artifact rdr=fresh.md` on a just-seeded RDR. The front-matter reader ran, found no `stage:`, and the key was "present" by the RDR's definition (reader invoked). The kernel folded the valueless key to false; every row matching on `stage` was excluded; the list was empty. The RDR had rejected R1 for producing "an empty list with nothing to act on." Its own presence definition produced one on the most common journey.

Post-incident, the fixes were three lines of policy the premortem could have demanded: presence means *value returned*, unbound roles for match-only keys degrade to absent rather than refuse, and reader refusals on non-selected-row keys are reported under `unknown`, not raised.

## 2. Obstacle negation

### N1 — "present-but-conflicted is unreachable; presence test is exact"
Negated by P-3 and P-4. Two distinct failures:
(a) *Valueless present.* "Owned by an invoked reader" is a fact about the demand set, not about the artifact. A reader that ran and returned no value (missing field) produces a key the CLI counts present and the kernel counts absent. The exact presence test is exact over the wrong predicate.
(b) *Check ordering.* "A `--tag` on an owned key is refused" — against which set? If the refusal compares `--tag` keys with the *invoked* readers' keys (natural, since that is the set the CLI has assembled), a `--tag` on an un-demanded owned key passed yesterday. The demand-set extension moves that key into the invoked set; the outcome is either a new refusal or a conflict, depending on the order of tag parsing vs reader invocation. This is a control-flow claim that must be shown from the check, not from prose.

### N2 — "presence-per-key licenses the per-key filter"
Negated by P-6. The filter is correct *given* that `KernelRow()` emits one `resolve.Tag` per authored atom. That is a cardinality property of the converter. If the converter (or `resolve.Row`'s constructor, or the kernel's `TagSet`) keys by `Key`, the second tag on a key is dropped before the filter runs — the S3 assumption, relocated. Also unaddressed: negative operators. If the model language has `ne`/`not-in` and the normaliser emits a Tag with a negate flag, "absent → omitted → candidate" is fine; if it expands `ne` into the enum complement as multiple rows, an absent key produces one `{key, absent}` candidate per complement literal.

### N3 — "dead rows need no CLI handling"
Negated by P-9. The remedy is a future RDR. Until then the no-tag query reports dead rows as candidates, and `in`-expansion multiplies them. The approach must at minimum say what the user sees in the interim and that it is acceptable.

### N4 — "strictly additive, no-op on zero-owned, changes resolve only on the match-only class where every run refuses"
Negated by P-1 and P-2. "Every existing run refuses `flow-no-match`" is true only for runs whose selected row matches on the key. Runs that select other rows produce plans today and fail tomorrow. The demand set is a union; the class is "every model with a match-only owned key," and within it the broken runs are the *working* ones. "Strictly additive" on the reader set is "strictly more refusal points" on the outcome, and the new refusals fire above the escape phase. The no-op claim for zero-owned models is safe; the rest is not.

### N5 — "scoping to `next` only is worse"
Partly negated. The asymmetry N5 describes is real, but the symmetric fix imports P-1/P-2 into the verb that has locked semantics. The alternative N5 never weighs: symmetric demand set with *soft* binding for match-only keys — unbound role → key absent (reported under `unknown` by `next`; folded to `no_match` by `resolve`, escape rows still reachable). That keeps one view per model without turning acceptance into refusal.

### N6 — "same key set in both views"
Negated by P-7. Two assembly paths, and readers that serve multiple keys. The proposal needs a single function that both verbs call for view assembly, or a test that asserts key-set equality per fixture. Also: is a `--tag` key over a non-owned key in the kernel view? If the kernel view is built from owned state only and `--tag` observed tags are merged separately, the key set differs by construction.

### N7 — "omission changes only gate reachability"
Negated by P-5. The empty-match probe is a new input the kernel has never received from `flow next` (0005 stripped match, so the kernel already saw empty-match probes — *unless* 0005's strip left a sentinel or the kernel distinguishes "stripped" from "empty"). The draft must show the kernel's row classification (escape vs transition) does not depend on `len(row.Match)`, and that table validation does not reject an empty match list.

### N8 — "`--all` yields exactly 0005; shared reader set"
Negated by P-8. Sharing the reader set with the new demand term means `--all` is 0005's predicate over a *larger* demand set; on models with an unbound match-only key role it exits 2 where 0005 listed. "Exactly 0005" is true of the predicate, false of the observable behaviour.

### R1 — strict kernel match rejected for empty lists
P-10: the chosen approach reproduces the enumeration defect on the no-tag query and (P-3) reproduces R1's empty list on valueless-present keys. R1's rejection stands only if presence means value-returned.

### R2 — opt-in `--match` rejected
Stands. But note the mirror: the proposal ships an opt-in `--all`; the same "every skill must carry the flag" argument applies to anyone who wanted 0005 semantics. Acceptable, state it.

### R4 — kernel three-valued match rejected
The rejection now cites the peer clause (S1 fixed). But R4's second reason — "no payload carrier for match facts" — is undermined by the proposal building exactly such a carrier CLI-side (`unknown`). The honest rejection is "the kernel must not change under the locked contract," not "no carrier exists."

### R5 — kernel reports undecided atoms rejected: "a fact the CLI already has"
P-12/P-13: the CLI has *presence*, not *value-absent* nor *conflict*. Under P-3 the CLI does not have the fact for valueless-present keys; the kernel does (it folded it). Either redefine presence (then R5's reason holds) or R5's reason is false.

### R6 — CLI dead-row detection rejected
Stands, provided P-6 is closed at the converter.

### R7 — `next`-only scoping rejected
See N5: the rejection is correct but the space is not two-valued; soft binding is the unweighed third option.

## 3. Consumer artifacts (what would have caught each at review)

- **P-1** — test `TestFlowResolve_MatchOnlyOwnedKey_UnrelatedRowStillResolves`: fixture with rows A (match `stage=x`, writes `stage`) and B (match `approved=yes`), `approved` served by role `approvals`; run `flow resolve --artifact rdr=…` with `approvals` UNBOUND and state `stage=x`. Expected today: plan for A. Expected after: the draft must state which — and if it is exit 2, the RDR's "nothing that worked stops working" line is struck.
- **P-2** — test `TestFlowResolve_NewlyDemandedReaderRefuses_EscapeRowStillRescues`: same fixture plus an escape row; bind `approvals` to an artifact the reader refuses. Assert whether the escape row runs. If it cannot (refusal above the switch), the RDR must say so.
- **P-3** — user journey `flow next --artifact rdr=<just-seeded, no stage:>`: expected N candidates carrying `{stage, absent}`; if the list is empty, presence is defined on the wrong predicate. Companion unit test `TestPresence_ReaderInvokedButValueless_IsAbsent`.
- **P-4** — test `TestTagRefusal_OwnedKeyNotDemanded`: model with match-only owned key K; run `--tag K=v` before and after the demand-set change; assert the refusal is identical and comes from the model's owned set, not the invoked set. Cite the check's line in the draft.
- **P-5** — test `TestFlowNext_AllMatchKeysAbsent_ProbeHasEmptyMatch_StillGate`: row with one match key, no `--tag`, no reader; assert the candidate disposition is gate-side and no table-validation error surfaces.
- **P-6** — test `TestKernelRow_EqPlusInOnOneKey_EmitsTwoTags` at the converter; and `TestFlowResolve_DeadRow_ExcludedWhenKeyPresent`. Cite the converter loop in the draft (the S3 seed says "cite the producer").
- **P-7** — test `TestViewKeySet_CLIEqualsKernel` over every fixture, including a reader serving two owned keys where only one is demanded; and a `--tag` on a non-owned key. Or: one shared assembly function, cited.
- **P-8** — test `TestFlowNextAll_UnboundMatchOnlyRole`: 0005 behaviour (list) vs after (exit 2); the RDR states which and drops "exactly 0005" if it changed.
- **P-9** — golden output for `flow next` on a fixture with one dead row and no `--tag`: the RDR shows the phantom candidate and accepts it in writing.
- **P-10** — golden output for `flow next --model models/rdr.toml` with NO `--artifact`/`--tag`: if it still lists ~21 rows, the problem statement's motivating number must be re-scoped to "with state supplied."
- **P-11** — `unknown` field presence asserted in the payload golden files; skills in-repo grepped for `unresolved`.
- **P-12** — test `TestFlowNext_ConflictedKey_IsReportedNotSilentlyExcluded` (only meaningful if P-4 opens a path; otherwise the draft cites the three refusals with lines).

## 4. Refutation-target sweep (seed classes found in this proposal)

- **S1 class (peer semantics asserted without opening the peer):** N1's three refusals are cited as facts about the loader and the tag parser with no line references; P-4 shows the ordering matters. N7's "independent of `row.Match`" is a claim about kernel row classification with no citation; P-5.
- **S2 class (order/reachability from prose):** N4's "strictly additive" and "nothing that worked stops working" — P-1/P-2 are the same shape as the seed: a widening at one layer that becomes a veto above the escape switch. N8's "exactly 0005" is a reachability claim over a larger demand set; P-8.
- **S3 class (cardinality from spec not producer):** N2 relocates the assumption from the CLI filter to `KernelRow()`; P-6. N6 asserts key-set equality of two structures built by two paths; P-7.
- **"X changes only on class Y" (N4):** the class is mis-scoped — defined by selected row, enforced by union. P-1.

## Disposition

The approach is recoverable without a switch, but not as written. Three mitigations are load-bearing: (1) presence = reader returned a value, valueless → `{key, absent}`; (2) match-only demanded keys bind *softly* — unbound role or reader refusal on a key not on the selected row degrades to absent (reported under `unknown` by `next`, folded to `no_match` with escape still reachable by `resolve`), or else N4's "nothing stops working" is deleted and the breaking class is named honestly; (3) cite the producer for every cardinality and reachability claim (converter loop, tag-refusal check set, kernel empty-match classification, view assembly function) and land the tests named above. If (2) cannot be made consistent with 0005:C1's locked wording, that is the genuine fork and NEEDS_DECISION.
