Model: claude-opus-5[1m] (dispatcher; persona passes each stamped in their own file)

# 3amigo — consolidation — cli/0012

Three isolated persona passes (PM, Implementer, QA), no cross-persona
visibility. Hotspots below are a mechanical count over `recs anchors`, not a
re-judging of the three files.

## Hotspots (an id raised by >=2 independent personas)

| id | personas | what converged |
| --- | --- | --- |
| `0012:C5` | **3/3** (PM-1, IMPL-5, IMPL-6, QA-1, QA-7) | C5's ingress list is the record's least-settled clause: scope (which ingresses), correctness (a named symbol that does not exist), and a live contradiction with RDR 0024's BOUNDARY tests |
| `0012:C1` | 2 (IMPL-2, IMPL-3, QA-2) | the nil-mapping/zero-value arm — no named producer for the kinds map, and no scenario pinning the LOUD verdict F6 rests on |
| `0012:C2` | 2 (IMPL-4, QA-4, QA-5) | arm-level under-specification: `in`-member poisoning, `enum`/`scalar` indiscriminability, and the re-keyed `gte`/`contains` cases |
| `0012:C3` | 2 (IMPL-1, QA-3) | the published fixture has no delivery mechanism and the second implementer no stated shape |
| `0012:C4` | 2 (IMPL-7, QA-6) | the "same loaded model" half — plumbing fork unresolved (A3), and no observable for a mismatch |
| `0012:MVV` | 2 (PM-2, QA-*) | the MVV's population is a fresh synthetic model; the installed base is unsized |
| `0012:F1` | 2 (PM-1, PM-2, QA-7) | the user-visible break inventory is where C5's second/third break is discoverable, not in the problem statement |

Overlap marks a hotspot PASSAGE, not a validated finding; a single-persona
finding is not thereby weaker (the personas never saw each other, so their
agreement is evidence, not conformity).

## Origin ledger (all 19 findings, by severity)

### High
- **QA-1** `0012:C5` — the `[emit]` ingress contradicts two live BOUNDARY-marked
  RDR 0024 tests on `main` (`TestReq9_0024_NonCanonicalIntLiteralsAreAdmitted`
  admits `03`/`+5`/`-0`; `TestReq7_0024_EveryKindCheckIsLexical...` carries
  `"007"`). No deviation recorded. Prevents: any `[emit]` C5 regression test —
  two contradictory oracles.
- **QA-2** `0012:C1`,`0012:S4` — nil-mapping LOUD disposition has NO scenario;
  S4 explicitly disclaims it. Prevents: the one-line unit test that is F6's
  actual proof. The "leave the raw-string arm alone" implementation passes
  S1-S8 in full.
- **IMPL-1** `0012:C3` — the published fixture has no delivery mechanism and
  `TestGuardEvaluatorContract`'s fixed 2-arg signature has no seat for it.
  Blocks: the exported API shape, with two in-tree callers.
- **IMPL-2** `0012:C1` — no producer named for `NewEvaluator`'s
  `map[string]string`; only `guard.DeclaredKinds(m)` in illustrative code
  stamped not-load-bearing. `internal/graphlint` needs it too, forcing export.
  Blocks: new exported symbol vs three copies of a loop; total over `m.Tags`
  vs filtered.

### Medium
- **PM-1** `0012:C5` vs `0012:§problem-statement` — the stated problem clears
  the CLI path that C5 then breaks. Blocks: whether C5 ships here, and whether
  Phase 4 is atomic with Phase 1.
- **PM-2** `0012:MVV`,`0012:F1` — no account of existing persisted artifacts;
  the silent False-to-True leg (non-canonical `"07"` reader value) is unsized.
  Blocks: the rollout decision.
- **IMPL-3** `0012:C1` — `Evaluator{}`'s LOUD disposition holds for `eq`/`in`
  only; `gte`/`contains` still answer correctly, so a missed site with only
  those operators migrates SILENTLY. F6 overclaims.
- **IMPL-4** `0012:C2` — the `in` arm does not say whether one unparseable
  member poisons the list or only itself; C3's case is authorable either way.
- **IMPL-5** `0012:C5` — S5's zero-occurrence measurement is about GUARD ATOMS;
  `[emit]`, `[initial]`, `[rule.write]` are three other populations the same
  tightening covers, unmeasured. (Independently reached QA-1's territory.)
- **IMPL-7** `0012:C4`,`0012:A3` — A3 offers "widen `probeRow` or pass the
  evaluator down" as an unchosen either/or; the two differ observably
  (per-request vs per-row construction).
- **QA-3** `0012:C3` — `conformingContractSeam` is bound as a second
  implementer but has no stated fixture-carrying shape, and is wrapped in
  `recordingSeam` for 0007's tests (a second consumer, unscoped).
- **QA-4** `0012:C2` — `enum` and `scalar` are behaviorally identical, so S2
  cannot discriminate them and C3's per-kind-token meta-check is vacuously
  satisfiable. A `default: string-equal` implementation breaks the two
  defensive arms.
- **QA-5** `0012:C2` — the "unchanged" clauses have no regression pin, and C3
  MOVES the 5 `gte` / 4 `contains` cases onto `int`/`set` keys at the same
  moment the seam becomes kind-aware. C2 does not say whether the kind lookup
  is consulted for `contains`.
- **QA-6** `0012:C4`,`0012:F4` — the "SAME loaded model" half is uncovered;
  no stated observable for a mixed-model evaluator. Defensible as untested,
  but belongs in Testing Strategy as a known gap, not only inside F4.

### Low
- **PM-3** `0012:§problem-statement` — the named beneficiary (the rdr
  navigator) never carried to an outcome; blocks acceptance sign-off.
- **PM-4** `0012:G-cross-cutting`,`0012:G-proportionality` — still TEMPLATE
  guidance blocks (Draft, so expected sequencing); flagged because those two
  gates are what would force written answers to PM-1 and PM-2.
- **IMPL-6** `0012:C5`,`0012:§existing-infrastructure-audit` — `normalize.go::atom`
  and `normalize.go::atomsFromBlock` do not exist on `main`; the described
  behavior is real (`normalize.go:172` `conform(...)`, `badAtom(...)`), the
  symbol names are wrong.
- **IMPL-8** `0012:C1`,`0012:§phase-3-producer-alignment` — eleven in-tree
  `Evaluator{}` test literals: Phase 1 or Phase 3? Between them those tests
  evaluate `eq` against a nil mapping and go red.
- **QA-7** `0012:C5`,`0012:S7` — the NEW CLI break (`--tag n=07`,
  `--write n=+1`) has no scenario, and its diagnostic envelope is undecided
  between C5's load-site wording and `flow_input.go:753`'s existing prefix.

## Source-anchor check (from the implementer pass)

Every symbol C1-C5 rests on EXISTS on `main` as described — verified
individually, not from the `resolved: true` edge flags. Two anchors name
symbols that do not exist (`normalize.go::atom`, `normalize.go::atomsFromBlock`);
neither is a load-bearing site for any contract. That is IMPL-6.

## Isolation

No persona file references another's output. Convergence on C5 (3/3) and
C1-C4 (2/3) is independent agreement.
