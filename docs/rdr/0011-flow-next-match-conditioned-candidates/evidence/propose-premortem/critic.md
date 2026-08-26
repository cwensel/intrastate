Model: claude-fable-5

# Premortem — RDR 0011 `flow next` match-conditioned candidates

Worked from the brief only. No repository file, RDR, or evidence file was read.

## Findings ledger

| ID | passage/claim | Failure mode | Symptom user sees | Origin |
|----|---------------|--------------|-------------------|--------|
| P-1 | C1 "conflicted key unreachable" (qualified: *within one provenance*) | Cross-provenance collision: a key both owned by an invoked reader and supplied via `--tag` with a differing value. Load checks reader-vs-reader; parse checks tag-vs-tag; nothing named checks reader-vs-tag. The flat `owned ∪ observed` map picks one value silently while the kernel's assembled view may hold both and fold conflicted-into-non-match. | Row vanishes from `next` with no `unknown` entry and no refusal; or `next` lists it while `resolve` refuses `no_match`. | negation |
| P-2 | C3 "omitting absent-key atoms cannot change the probe" | Probe of a row whose *every* match key is absent submits a row with zero match atoms. The brief never says what the kernel makes of a match-less ordinary row: matches-everything, normalization error, or reclassification as an escape/catch-all row. If the last, the probe's refusal class and gate handling change. | Candidate reported with the wrong disposition, or a probe error surfaced as a CLI crash/exit code on a decision-table with no `--tag`. | refutation-target (S1 class) |
| P-3 | C2 "CLI presence agrees with kernel presence" | Presence is not binary at the kernel: empty-string/zero values from a reader, key normalization (case, namespacing `artifact.stage`), or a reader that ran but returned no value for a declared key. CLI says "present", kernel says absent (or vice versa). | Row excluded via `no_match` for a key the CLI reports as present; or listed as `absent` when `resolve` would match it. | negation |
| P-4 | C2 view membership "owned tags from invoked readers" | Which readers are invoked is itself a choice. If 0005 invokes readers lazily (only those a probed row's atoms name) or skips readers under guard pruning, the view differs per row and per flag (`--all` omits every match atom, so may invoke no reader). Presence becomes probe-order dependent. | Same model, same state, different `unknown` lists between `next` and `next --all`, or between two rows naming the same key. | negation |
| P-5 | C6 "Refusal.Undecided is sufficient" | The payload may be short-circuited (first undecidable atom only), keyed by expression rather than by state key, or lack a reason distinguishing absent from uncomparable. A guard atom over an absent key appears both in the CLI's absence scan and in `Undecided`, producing a duplicate `unknown` entry with two reasons. | `unknown` shows one atom when three were undecidable; or `{stage, absent}` and `{stage, uncomparable}` for the same key. | refutation-target (S1 class: field named without opening the contract) |
| P-6 | C5 "all 0005 tests pass unchanged" | A 0005 fixture that proves guard-pruning by supplying only the *guard's* key while its rows match on another key will now list those rows as candidates-with-`absent` (passes) but a fixture asserting candidate *count* (e.g., 21) or a golden JSON file with `unresolved` breaks beyond the "six reads". | Test failures larger than scoped; implementer widens the diff under time pressure and edits assertions to fit. | negation |
| P-7 | C8 "21 → exactly 3" | Number asserted from prose. Depends on which readers are invoked (P-4) and whether `stage` is owned or observed on the project model. | Ships with 21→5 or 21→2 and nobody notices because no test pins the rule ids. | refutation-target (S2 class: cross-artifact claim from prose) |
| P-8 | C4 sibling decision-table need | S2 class exactly: "every ordinary row stays a candidate with no `--tag`" is asserted from prose about the sibling. If the sibling defines its catch-all as an *ordinary* row matching on `recognized`, or defines "otherwise" as an escape row that `next` never lists, a no-tag `next` either lists the catch-all as a peer of real candidates or lists nothing rescuable. | Skill author sees the "otherwise" branch as a legal next step, or sees an empty list on an empty-state decision table. | refutation-target (S2 class) |
| P-9 | C9 / R4 `--all` restores set, not shape | Existing 0005 skill scripts read `unresolved`; the rename ships with no shim. `--all` is advertised as "restores 0005" but 0005 consumers still break. Also unspecified: under `--all`, are omitted-but-present match keys listed as `absent`? | Every 0005-era skill that parses `next` breaks on upgrade; `--all` output either lies (`absent` for a present key) or omits match keys from `unknown` entirely, and the spec says neither. | negation |
| P-10 | R3 rejection ("strict kernel match empties partial-state lists") | Inverted: `flow resolve` *is* strict today (absent key → `no_match`, kernel unchanged). So `next` and `resolve` disagree by design on absent keys. "Candidate" is being redefined as "not yet excluded" without saying so. | Author runs `next`, sees rule X, runs `resolve` with the identical state, gets `no_match` or the escape row. Files a bug: "next lied". | negation |
| P-11 | R2 rejection ("buys nothing the CLI cannot already see") | Inverted by P-1: the one case the CLI cannot see is the cross-provenance conflict that R2's payload would have reported. | Same as P-1. | negation |
| P-12 | R1 rejection cites the locked peer clause as a *resolution-level* veto | Correct for `resolve`, but the brief's rejection also relies on "undecided payload is only for GUARD rows" — which is the same class-of-payload assumption P-5 attacks. If the payload does carry match facts, R1's third leg falls and R1's first leg (veto) does not apply to a one-row probe at all. | None directly; a documented escape hatch for a later "why not R1" reopen. | refutation-target (S1 class) |
| P-13 | C7 "cobra unknown-flag error is an adequate negative" | Adequate only if `flow` has no persistent `--all`/`-a`, no `FParseErrWhitelist`, no `DisableFlagParsing`, and `TraverseChildren` is not set. None of these is named. | `flow resolve --all` silently accepted and ignored. | negation |
| P-14 | Approach: "kernel `no_match` for a present unequal key excludes the row" | On a one-row probe `no_match` is a *table* outcome; a guard decided false on a one-row table may also surface as `no_match`. Fine for exclusion, but the CLI cannot then tell "excluded by match" from "excluded by guard". If a later RDR wants `next --why`, this collapses. Also match atoms with non-equality semantics (negation, absence tests, set membership) — omitting an absent-key `not:`/`absent:` atom changes meaning from "must be absent" to "unconstrained". | A row whose rule is "only when `stage` is unset" is listed as candidate with `unknown: {stage, absent}` even when it is in fact the *only* selectable row; conversely a row `stage: not resolved` is dropped when it should be kept. | negation |
| P-15 | Hindsight narrative | Combined P-1 + P-10: a skill author with an owned `stage` reader also passed `--tag stage=review` to "preview" the next state. | See narrative. | hindsight |

## 1. Prospective hindsight — the failure as accomplished fact

RDR 0011 shipped in the release after 0005. The flipped default did what it promised on the project's own model: `next` at `stage=resolved` went from 21 rows to a handful and the skill authors stopped complaining.

Three weeks later the decision-table sibling landed. Its authors wrote a skill that calls `flow next --model dt.yaml` with no tags to render a menu, then `flow next --model dt.yaml --tag kind=<choice>` to narrow. The first call listed the "otherwise" row alongside the real branches, because in the sibling model "otherwise" was an ordinary row matching on `recognized` and 0011's normalization lifted `match.recognized` into the outcome, leaving the row with no match atoms — a match-less ordinary row that the probe treated as matching everything. The skill offered "otherwise" as a legal choice. Nobody at review had asked what the kernel does with a row that has zero match atoms after omission; the RDR said "the kernel filters on match before the guard gate" and everyone nodded.

The second incident was the one that got a bug filed as "next lies". An author previewing a transition had a reader owning `stage` and also passed `--tag stage=review`. Tag parsing accepted it (only *repeated* tag keys are refused); load accepted it (only reader-vs-reader ownership is checked). The CLI's flat `owned ∪ observed` map kept whichever was written last, said `stage` was present, kept the match atom, and probed. The kernel's assembled view held both provenances, folded the conflict into non-match, and answered `no_match`. The row disappeared with no `unknown` entry and no refusal. `flow resolve` with the same arguments refused. The RDR's claim 1 had a qualifier — "within one provenance" — that quietly excluded the only reachable conflict.

The third was mundane: every 0005-era skill that parsed `unresolved` broke on upgrade, and `--all`, advertised as "restores 0005", did not restore the field. The RDR had said so (claim 9) and nobody read it as a migration cost.

## 2. Obstacle negation

### Claim 1 — conflicted key unreachable
Negated by P-1. The three guards named (reader-vs-reader at load, tag-vs-tag at parse, reader-declares-keys) leave reader-vs-tag unguarded. The brief's own phrasing "within one provenance" is the tell: cross-provenance collision is the case that must be either refused at the CLI or reported. The approach must answer: what does `flow next` do when an observed `--tag` key coincides with an owned reader key? Options that keep the approach: refuse at parse-time after load ("tag `stage` collides with owned key from reader `X`"), or define precedence and *report* it on the candidate (`unknown: {key: stage, reason: "shadowed"}`). Either is a spec sentence and a test; neither forces a switch.

### Claim 2 — CLI presence equals kernel presence
Negated by P-3 and P-4. "Present" must be defined identically on both sides: is an empty value present? Are keys normalized? And the view must be fixed *before* the per-row loop (one reader invocation set for the whole `next` call, identical under `--all`), otherwise `unknown` is order-dependent. The `recognized` lift is asserted, not shown — P-2 asks what a row is after the lift removes its only atom.

### Claim 3 — omission cannot change the probe
Negated by P-2 and P-14. Two sub-cases: (a) all atoms omitted → match-less row; the kernel's treatment of that row is unstated and may reclassify it. (b) Non-equality atoms (negation, absence, membership) — omission does not preserve meaning; "stage must be absent" becomes "unconstrained". If the match grammar is equality-only, say so in the RDR as a load-bearing assumption; if not, absence/negation atoms must be *kept* (they are decidable on an absent key) rather than omitted.

### Claim 4 — decision-table no-tag behavior
Negated by P-8 (S2 class). This is a control-flow claim about the sibling's row classification asserted from prose. It needs the sibling's actual row shapes: which rows are ordinary, which are escape, whether the catch-all matches on `recognized` only.

### Claim 5 — 0005 tests pass unchanged
Negated by P-6. "Each fixture supplies the key its rows match on" is a survey result that should be an enumerated table (fixture → match keys → supplied keys) in evidence, not a sentence. Count assertions and golden files are the usual breakage beyond field renames.

### Claim 6 — `Refusal.Undecided` suffices
Negated by P-5. Needs: the payload's cardinality (all undecidable atoms or first), its key (state key vs expression), and its reason vocabulary. Also the merge rule when a key is both view-absent and in `Undecided` — one entry, reason `absent`.

### Claim 7 — non-registration is an adequate negative
Negated by P-13. A one-line test per verb (`flow resolve --all` → non-zero exit, message contains `unknown flag`) turns the claim into a fact and guards against a future persistent flag.

### Claim 8 — 21 → 3
Negated by P-7. Pin the three rule ids in a test over the project model; the number is otherwise a prose artifact.

### Claim 9 — rename and flip ship together
Negated by P-9. Two gaps: (a) no consumer migration is named; (b) `--all`'s `unknown` contents for present-but-omitted match keys are unspecified. Decide: under `--all`, match keys are not listed in `unknown` at all (they were not probed), or they are listed with a distinct reason (`skipped`). Do not reuse `absent`.

### R1 — kernel three-valued match
The rejection's first leg (resolution-level veto) is sound for `flow resolve` and irrelevant to a one-row `next` probe, where a veto of a one-row table is just "not a candidate" — so R1 would *not* have been catastrophic for `next`. The rejection stands only on the second and third legs (escape rows cannot rescue; payload cannot carry match facts), and the third is the same assumption P-5 attacks. Keep R1 rejected, but rewrite the rejection so it does not rest on the veto for `next`.

### R2 — report undecided match atoms on `no_match`
Inverted by P-1/P-11: there is one thing the CLI cannot see (cross-provenance conflict). The approach can still reject R2 by closing that gap at the CLI (refuse or report the collision), which restores "nothing the CLI cannot see". Without that mitigation the rejection is false.

### R3 — strict kernel match
The rejection is correct for `next`, but P-10 shows the chosen approach makes `next` permissive while `resolve` stays strict. That is a defensible definition of "candidate" — "not yet excluded by supplied state" — but it must be stated in the payload contract and in the command help, or authors will read `next` as "would resolve".

### R4 — opt-in `--match`
Rejection stands. Note that `--all` is the mirror of R4 and inherits R4's discoverability problem in reverse: skills that *wanted* enumeration must now know the flag. Acceptable; the defect was the default.

## 3. Consumer artifacts — what would have caught each at review

| Finding | Artifact |
|---------|----------|
| P-1 | Test `TestNext_ObservedTagCollidesWithOwnedKey`: model with reader owning `stage`; run `next --tag stage=review`; assert either a refusal naming the collision *or* the row listed with an `unknown` entry naming `stage` — never a silent drop. Companion journey: "preview a transition by overriding an owned key". |
| P-2 | Test `TestNext_AllMatchKeysAbsent_RowStaysCandidate`: row matching only on `kind`; no tags; assert row listed, `unknown = [{kind, absent}]`, and the probe outcome class is unchanged (not escape, not error). Add `TestNext_RecognizedOnlyRow` for a row whose sole match is `recognized`. |
| P-3 | Test `TestNext_ReaderReturnsEmptyValue`: reader declares `stage`, returns `""`; assert the CLI and kernel agree (document which way). |
| P-4 | Test `TestNext_ReaderInvocationSetIsFlagIndependent`: assert the set of invoked readers is identical for `next` and `next --all` (fake reader records calls). |
| P-5 | Test `TestNext_GuardTwoUndecidableAtoms`: guard with two undecidable atoms over present keys; assert both appear in `unknown` with `uncomparable`. Test `TestNext_GuardOverAbsentKey_SingleEntry`: assert exactly one `unknown` entry, reason `absent`. |
| P-6 | Evidence table: every 0005 `next` fixture → rows' match keys → keys the fixture supplies → expected delta. Not a sentence. |
| P-7 | Test `TestNext_ProjectModel_StageResolved`: assert the exact rule-id set (three ids), not the count. |
| P-8 | Joint test in the sibling RDR's fixture: `next` with no tags over the decision-table model asserting the listed ids exclude the "otherwise" row, plus `next --tag kind=<mismatch>` asserting the drop. Written against the sibling's actual row shapes, per S2. |
| P-9 | Payload contract diff in the RDR listing `unknown` reasons per flag state, with a row for `--all`. Golden JSON for `next --all`. |
| P-10 | User journey in the RDR: "next lists X; resolve with the same state → ?" with the answer written down; command help text for `next` containing "candidate = not excluded by supplied state; may still fail to resolve". |
| P-13 | Test `TestAllFlag_RejectedOnSiblingVerbs` over `resolve`, `read-state`, `set-state`. |
| P-14 | Test over each match-atom operator the grammar supports, with the key absent, asserting the row's disposition; or an RDR assumption "match atoms are equality-only" verified at resolve stage. |

## 4. Refutation targets (seed classes applied here)

**S1 class — designing against a contract the draft never opened; naming a mechanism by description.** Present in: "the kernel filters on match before the guard gate, and the two are evaluated independently" (C3); "`Refusal.Undecided`" as a named field with a reason vocabulary (C6); "a row's `match.recognized` is lifted … and never appears as an atom" (C2); "the kernel's fold of conflicted-into-non-match" (C1). Each is a mechanism described, not cited. Before accepting, the RDR should quote the kernel's row-evaluation order and the refusal payload type, and state what a row with zero match atoms is. P-2, P-5, P-12 are this class.

**S2 class — cross-RDR safety claim asserted from prose instead of control flow.** Present in C4 (sibling decision table) and C8 (project-model count). P-7 and P-8 are this class; both convert to id-pinned tests.

## Verdict reasoning

Nothing here forces a switch away from match-conditioned `next` with CLI-side presence: every failure is closed by a spec sentence plus a named test. The two that must be folded before lock are the cross-provenance collision (P-1, which also repairs the R2 rejection) and the match-less-row/non-equality-atom semantics (P-2, P-14), because both produce silent wrong answers rather than visible errors.
