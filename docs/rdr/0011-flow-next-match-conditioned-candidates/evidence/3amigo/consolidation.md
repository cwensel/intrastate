Model: claude-opus-5[1m] (dispatcher; persona passes stamped in their own files)

# 3amigo Consolidation — RDR cli/0011

Three isolated persona passes (`persona-1-pm.md`, `persona-2-implementer.md`,
`persona-3-qa.md`), each run as a separate sub-agent seeing only the RDR and its
own persona block. No file references another's output — isolation held.

Consolidated **mechanically**: hotspots are the set intersection over the element
ids each persona anchored to, not a re-judgement by a model that read all three.
Overlap marks a hotspot passage, not a validated finding; a single-persona
finding is not thereby weaker.

## Hotspots (≥2 personas anchored the same id)

| id | personas | what they independently landed on |
| --- | --- | --- |
| `0011:C1` | PM, Impl, QA | the contract every finding orbits: outcome claim, missing return channel, ordinary-site blast radius, reason vocabulary |
| `0011:MVV` | PM, Impl, QA | step 6 is unreachable (Impl-1), unwritable (QA-2), and asserts cardinality not outcome (PM-F1) |
| `0011:A4` | PM, Impl, QA | the predicate-scoped BREAKS-0 result and what it does/doesn't excuse |
| `0011:A7` | Impl, QA | still `Pending`; its census scopes to escape oracles only, and QA shows ordinary-site oracles break too |
| `0011:S5` | Impl, QA | asserts a shipped oracle "keeps passing" that C1 inverts |
| `0011:S3` | Impl, QA | the `uncomparable` fixture cannot be built as written |
| `0011:D-undecided-vocabulary` | Impl, QA | closed reason set has no token for an un-run gate id |
| `0011:§cross-cutting-concerns` | PM, Impl | unauthored template text; versioning + determinism both apply |

**The strongest signal**: Impl-2 and QA-2 are the *same defect* reached from two
directions with no shared context — the probe returns a `Plan`, and `Plan`
carries no undecided-facts field, so C1's `unknown` obligation has no carrier.

## Origin ledger

Each entry is a finding this round raised. Dispositions are written by the
resolve half; net-new scope is charted, never folded in.

| # | id | sev | origin | finding (one line) | disposition |
| --- | --- | --- | --- | --- | --- |
| L1 | `0011:C1` `0011:MVV` `0011:S3` `0011:C3` | BLOCK | Impl-1 | C3's mandated `uncomparable` fixture is unreachable via the CLI: `internal/cli/flow_input.go::parseTags` hard-refuses a repeated `--tag` (`codeTagDuplicate`) before `kernelTags`. A3's evidence sentence is true of `kernelTags` in isolation, false of the path. | open |
| L2 | `0011:C1` `0011:MVV` `0011:S3` | BLOCK | Impl-2 + QA-2 (hotspot) | C1 obliges every indeterminate match atom into `unknown`, but a match-indeterminate probe returns a `Plan`, and `resolve.Plan` has six fields and no undecided payload; `Refusal.Undecided` is documented `guard_unevaluable`-only. No named kernel→CLI channel. | open |
| L3 | `0011:C1` `0011:A7` `0011:§consequences` | BLOCK | Impl-3 | C1's seam-uniformity clause changes ORDINARY `flow resolve` selection at `resolve.go::Resolve` (`!view.matches(row.Match)`), not just the escape site: (a) an unescaped write-performing plan where `TestReq78` expects a rescue; (b) reachable `ambiguous_match` on real multi-row tables. A7 and Consequences account only for the escape half. | open |
| L4 | `0011:S5` `0011:A7` | BLOCK | QA-1 | S5 asserts `TestAdv0007_3_EscapeNotNamingTheConflictedKeyStillRescues` "keeps passing"; under C1 it fails (guard-free row `A` survives as an unescaped plan, rescue never entered). Two siblings break unnamed, one — `..._ConflictedKeyIsNotEscapableByNamingIt` — is a direct inversion of C1. S5's "suite passes with TestReq78 retired" is unsatisfiable as written. | open |
| L5 | `0011:C1` `0011:D-undecided-vocabulary` | HIGH | QA-3 | C1 puts un-run gate ids into `unknown` as `{key, reason}`, but the closed set (`absent`, `uncomparable`) has no token for "never consulted because `--evaluate-gates` was not passed". Blocks C3's "mechanical" re-homing of the gate-id assertions. | open |
| L6 | `0011:C1` `0011:D-undecided-vocabulary` | HIGH | Impl-4 | C1 makes match atoms ride `Refusal.Undecided`, whose doc scopes it to `guard_unevaluable` and whose producer discards non-`GuardUnevaluable` payloads; `guard.go::isGuardBlock` deliberately excludes `BlockMatch`. Does `0007:C8`'s payload widen, or does C1 mint a separate field? | open |
| L7 | `0011:C1` `0011:§mini-checks` | HIGH | Impl-5 | `unknown` merges a CLI-derived `absent` (from `assembledView`, which reports a conflicted key as present) with a kernel-derived `absent`/`uncomparable` — the two views disagree by construction. Dedup key with `{key, reason}` pairs is unstated. | open |
| L8 | `0011:§problem-statement` `0011:MVV` `0011:§consequences` | HIGH | PM-F1 | The stated outcome ("a skill can constrain its next step") has no acceptance criterion; every MVV/S step asserts cardinality. 3 candidates across 3 outcomes is still not a selection, and `next` by design never selects. Record never reconciles narrowing-vs-selecting. | open |
| L9 | `0011:§problem-statement` `0011:F2` `0011:A5` | HIGH | PM-F2 | The wall-of-candidates symptom survives for observed-key models: A5 shows 0010's decision-table class with no `--tag` reports every row. Outcome lands where match keys are reader-owned; record holds both facts and never joins them. | open |
| L10 | `0011:C2` `0011:§consequences` `0011:C3` | MED | PM-F3 | `--all` restores 0005's candidate SET but not the field NAME/TYPE (`unresolved` `[]string` → `unknown` `[{key,reason}]`); `flow_next.go` emits `json:"unresolved"` on a documented `--as=json` scripting surface. Escape hatch does not restore the old contract; record nowhere says so. | open |
| L11 | `0011:C2` `0011:A4` | MED | Impl-6 | Under `--all`, `summarize` walks `row.Atoms`, which already carries match-block atoms (A2) — so match keys already land in `unresolved` today. C2's "neither an exclusion nor an entry in `unknown`" is a REMOVAL of shipped 0005 behaviour, framed as preservation. | open |
| L12 | `0011:C2` `0011:D-naming` | MED | Impl-7 | C2's "`--all` MUST NOT be accepted by `resolve`/`read-state`/`set-state`" has no stated mechanism: non-registration (cobra `unknown flag`) vs an explicit typed refusal produce different exit codes and different oracles. | open |
| L13 | `0011:S6` `0011:A4` | MED | QA-4 | S6's "all 28 oracles pass, none re-homed" is uncheckable: the 28-name enumeration A4 claims exists is not in the record or `evidence/`. Partly self-mitigating (S6 is an anti-oracle). | open |
| L14 | `0011:C3` | LOW | Impl-8 | "six shipped reads" counts test assertions only; a field rename also touches `candidate.Unresolved`, the `json:` tag, the header comment, and three `slices.Contains` uses. The count is load-bearing for the gate. | open |
| L15 | `0011:D-undecided-reporting-shape` `0011:C3` | LOW | PM-F4 | Every MVV step is `--as=json`; the text-mode rendering of `unknown` is undecided, though `D-undecided-reporting-shape` cites terraform's text-mode precedent and README advertises `--as text`. | open |
| L16 | `0011:C1` `0011:D-identity` | LOW | QA-5 | `unknown` element ordering across the two-source merge (kernel-sorted atoms + CLI insertion order) is unspecified. `D-identity` says "set", so satisfiable — a spec gap, not an unwritable test. | open |
| L17 | `0011:§cross-cutting-concerns` `0011:§contradiction-check` `0011:§scope-verification` `0011:§proportionality` | LOW | Impl-9 + PM-F3 tail | Finalization Gate sections are unauthored template text. `§scope-verification` is what would say whether MVV 6 is a scope deferral or a defect (upstream of L1). | open |

## Grounded at consolidation (dispatcher, against `main` @ 36e411e)

Confirmed before any edit — these are code facts, not persona claims:

- `internal/cli/flow_input.go:256` — `case seen[key]: return nil, userErr(codeTagDuplicate, …)`. L1 holds.
- `internal/resolve/resolve.go` `Plan` — six fields (`RuleID`, `SourceLocator`, `NextTags`, `Writes`, `Revision`, `Escaped`), no undecided payload. L2 holds.
- `internal/resolve/resolve.go:370-375` — `Undecided` doc reads "on `guard_unevaluable`". L2/L6 hold.
- `internal/resolve/guard.go:162` — `if verdict != GuardUnevaluable { return verdict, nil }` discards non-unevaluable payloads; `:121` `if !isGuardBlock(atom.Block)`. L6 holds.
- `internal/resolve/resolve.go::Resolve` ordinary loop — `if !view.matches(row.Match) { continue }`, then `switch len(selected)`. L3 holds (both arms).
- `internal/resolve/guard_atoms_test.go:930` `TestReq78…` asserts `plan.Escaped`. L3(a) holds.
- `internal/resolve/match_conflicted_test.go` — `dupKey = "gate"`; `gateMatchRow` is guard-free with `Match: [status=Draft, gate=closed]`; `…_ConflictedKeyIsNotEscapableByNamingIt:151` asserts `got.Refused()` on the stated ground that the escape path "runs the same match filter". Direct inversion of C1. L4 holds.
- `0011:§cross-cutting-concerns` projects as verbatim template brackets. L17 holds.
- `internal/cli/flow_next.go:49` — `Unresolved []string \`json:"unresolved"\``; `docs/cli-output-contract.md:19` documents `--as=json` as a scripting surface (envelope-level; it does not pin per-command fields). L10 holds on the emit site and the scripting posture.

## Review gate

**Healthy.** Every finding across all three files anchors to a named passage and
states the decision or test it blocks. No persona file cites another's output.
The PM and QA passes each report a "not findings — checked and clear" section,
which is the shape of a real read rather than a quota fill. Two personas
independently reaching L2 from different directions is agreement, not conformity.

---

# Iteration 2 — delta-scoped re-run

The iter-1 resolve rewrote C1's second and third paragraphs (a substantial fix),
so the lens re-ran **delta-scoped to the open ledger entries** — Implementer and
QA only, the two personas whose findings drove the rewrite. PM's entries (L8–L10,
L15) were closed by scope statements that restructured no contract, so they were
not re-run. Files: `iter-2/persona-2-implementer.md`, `iter-2/persona-3-qa.md`.

## Iter-1 ledger, closing dispositions

| # | disposition | note |
| --- | --- | --- |
| L1 | **fixed** | conflicted producer corrected to two-reader owned collision; `parseTags` refusal named at every site. iter-2: CLOSED |
| L2 | **fixed** (2 passes) | carrier named (`Refusal.Undecided`); iter-2 found the compound case still uncarried → L18 |
| L3 | **fixed** | C1 split into selection vs reporting sites; both ordinary-site cases now in Consequences. iter-2: CLOSED |
| L4 | **fixed** (2 passes) | S5 rewritten to enumerate changed oracles; iter-2 found a fifth → L19 |
| L5 | **fixed** | `not-evaluated` minted for the un-run-gate class only. iter-2: CLOSED |
| L6 | **fixed** | payload widening booked as A8 (Pending). iter-2: CLOSED, scoping correct |
| L7 | **fixed** | dedup on the `{key, reason}` pair. iter-2: CLOSED |
| L8 | **fixed** | Problem Statement now states what "constrained" means and that `next` narrows rather than selects |
| L9 | **fixed** | same clause states the outcome lands where match keys are decidable; A5's observed-key case named as intended, not residual |
| L10 | **fixed** + charted | Consequences states `--all` restores the set not the field; policy question charted |
| L11 | **fixed** | C2 now says the `--all` match-atom suppression is a departure from shipped behaviour and owes an oracle. iter-2: CLOSED |
| L12 | **fixed** | non-registration named as the mechanism. iter-2: CLOSED |
| L13 | **fixed** | S6's "28" demoted to provenance; checkable form stated. iter-2: CLOSED |
| L14 | **fixed** | C3 clarifies six = TEST reads; production rename sites listed |
| L15 | **fixed** | text-mode rendering and `unknown` ordering decided in `D-undecided-reporting-shape`. iter-2: CLOSED |
| L16 | **fixed** | sorted by `(key, reason)`. iter-2: CLOSED |
| L17 | **dismissed-with-cite** | Finalization Gate sections are Stage 7's to author into `gate.md`; template text is correct at Draft |

## Net-new from iter-2

| # | id | sev | origin | finding | disposition |
| --- | --- | --- | --- | --- | --- |
| L18 | `0011:C1` `0011:A1` | BLOCK | iter-2 Impl-2 | The carrier hole moved one disposition over: `gate` returns `owned_state_unavailable` BEFORE collecting the `guard_unevaluable` payload, so a row that is both indeterminate-on-match and missing an owned key carries no match facts. Also unstated: where the match check sits relative to `missingOwned` — pruning an indeterminate row upstream would drop its owned-state obligation. | **fixed** — C1 now fixes the precedence (preserved, not reordered), states the CLI's fallback for the compound row and the single `uncomparable`-not-recoverable limit, and settles the placement via `0007:C11`: `no-match` prunes (decided), `indeterminate` survives into `gate` (undecided). New S7 + MVV 6c + two mini-check rows. |
| L19 | `0011:S5` `0011:A7` | BLOCK | iter-2 QA-1 | S5's enumeration was incomplete: `reserved_key_0008_test.go::TestReq3_AbsentOutcomeYieldsAViewWithNoRecognizedKey` also inverts — `assemble` skips the `recognized` binding when `Input.Recognized` is empty, so the row's `recognized` atom names a key ABSENT from the view (`indeterminate`, not `no-match`). | **fixed** — added as the fifth re-expected oracle with its subject preserved (`got.Plan` stays nil, only the kind changes); A7 rescoped to five and its census widened to every file building a `Match` block; MVV 8 updated. |

## Grounded at iter-2 (dispatcher, against `main` @ 36e411e)

- `internal/resolve/resolve.go::gate` — `missingOwned` → `KindOwnedStateUnavailable` returns before the `GuardUnevaluable` collection loop. L18's precedence claim holds.
- `internal/resolve/guard.go` block-boundary comment — sweeping a non-guard atom into `all_result` "would prune its row along with the row's owned-state obligation (D8) — reopening the masking path `0007:C11` ratifies pruning against". Names the exact hazard, and `BlockMatch` explicitly.
- `0007:C11` — ratifies D8 pruning "conditional on the domain rule: pruning is safe exactly because GuardFalse can only arise from decided atoms". This is what decides the placement fork: `no-match` decided ⇒ prune; `indeterminate` undecided ⇒ must not.
- `internal/resolve/resolve.go::assemble` — `if in.Recognized != ""` guards the binding, so an empty `Recognized` leaves the key absent. L19 holds.
- `internal/resolve/reserved_key_0008_test.go::TestReq3_AbsentOutcomeYieldsAViewWithNoRecognizedKey` — asserts `KindNoMatch` on exactly that input. L19 holds.
- `docs/jdr/0001-resolve-kernel-seam.md` §Principles — P1 "Missing artifact state is never masked", P2 "The kernel refuses rather than guesses", P7 "Pre-release, prefer the clean shape" all resolve; the RDR's leans are grounded (an earlier dispatcher note claiming otherwise was wrong and is corrected in `Charted.md`).

## Convergence

**Converged.** No open ledger entries: L1–L17 fixed / dismissed-with-cite /
charted at iter-1; L18–L19 fixed at iter-2. Both iter-2 personas reported every
other entry CLOSED against the revised text, and QA reported no new untestable
claim from the rewrite (entry 6). Two iterations, well inside the cap of 3.
