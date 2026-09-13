Model: claude-opus-5

# Stage 6 Reconcile — RDR 0021 (lint's normalized-graph export)

Date: 2026-09-12 · Iteration 1 · Verdict: **NOT RECONCILED** (2 BLOCKERs, class `contract`)

## Stage 5 preflight

| Check | Answer | Note |
| --- | --- | --- |
| `--outcome lens` | `/rdr-reconcile` | `lens-large-row-complete`; Profile `large`, grounding→3amigo→critique row complete |
| `--outcome critique` | `none` | `critique-large-diffed`; `critique_models=differ` — no single-model fallback |
| `--outcome repeatability` | `none` | `repeatability-lite-complete` |
| `resolve:determinacy` chain | routed to `/rdr-prelock repeatability 1` | **Not taken as a stop** — see Caveats |

## Open set and dispositions

| # | Item | Source | Disposition | Evidence / plan |
| --- | --- | --- | --- | --- |
| 1 | A9 — `<opaque>` sentinel unforgeable; opaque-admits marking sound | 1, 2 | **BLOCKER** | Limb (a) REFUTED. Only `ClearSentinel` is reserved (`internal/table/model.go:11-14`, `normalize.go:64-66,617-620`, `load.go:1692-1695`); `conformKind`/`conformDomain` (`load.go:1853-1878`) have no `scalar` arm. Probe: a `scalar` tag, or an enum declaring `domain = ["<opaque>", …]`, loads clean with `[initial] free = "<opaque>"` and yields a node byte-identical to the synthesized one. Limb (b) HOLDS — `nodeMeetsAll` runs only over `splitNode` singleton-split nodes (`analysis.go:391-409,423-444`). |
| 2 | `0021:C2` restates a merged-node terminal quantifier JDR 0001 §JD-23 homes | 1 (round residue), 4 | **BLOCKER** | §JD-23 (`docs/jdr/0001-resolve-kernel-seam.md:941-960`): doctrine binds both liveness invariants; "neither record widens or narrows the split unilaterally. Both records cite this entry and neither restates the other's contract." Siblings 0015, 0022 — not 0021. `0015:C1` adopts it as 0015's answer. C2 writes the opposite (existential/opaque-admits) reading in its own block. |
| 3 | A8 — exported `graphlint` surface addable without altering `reach()` or lint's path | 1, 2 | **VERIFIED** (wording narrowed) | Operative limb holds: `reach()` returns `(nodes []Node, complete bool)` (`reach.go:73-76`); `Reach` discards it (`reach.go:90`); lint reaches `reach()` via `Run`→`newAnalysis` (`analysis.go:28-46`, `engine.go:59-65`), never via `Reach`; edges computed transiently in `successorsOf` and discarded, so recovery is additive. **Narrowing:** A8's rationale ("names a condition no caller can currently observe") is false — `CodeProductTooLarge` fires on exactly `complete == false` (`analysis.go:46,177-188`; `reach.go:128-133`) and is exposed on exported `Report.Findings` (`engine.go:13-32,59-80`), a faithful 1:1 proxy. Only the raw typed bool is unobservable. |
| 4 | F4 — wire `model` field identity unfixed | 1 (repeatability residue) | **Amendment recommended** (rides the route-back) | Carrier settled: `Model.ID` exists (`internal/table/model.go:497`) and A3 (Verified) already names `Model.ID`/`Class`/`Tags`/`Initial`/`Terminal`/`Rows` as C2's field-list carriers. C2 should pin `model` to `Model.ID`, not the `--model` path arg. Not MVV-load-bearing: MVV step 2 asserts identity/class are present, not their spelling. |
| 5 | Exactness sweep — "canonical" ×3 (420/423/454), "deterministic" (497) | 4 | **No action** | All `conformance` tier, all carry Evidence Records. 420/423 cite 0002's authored total row ordering and byte-lexicographic atom canonicalization (`0002:§normative-contracts`); A4 (Verified) carries the peer dependency. 454 rests on A7 (Verified, `(Node).key` as internal fingerprint) and A2 (Verified, byte-identical under `GODEBUG=randmapiter=1`). 497 is C3's own subject, oracled by S1/S2 and MVV step 2. |
| 6 | Spikes named but unrun | 3 | **Empty** | `spikes_unrun=[]`; absorption audit found no spike named in a findings file that the RDR does not name. |
| 7 | Hollow-body / References completeness | prompt gate | **PASS** | No `_Draft placeholder._`, no seed-skeleton header, no bracketed placeholders outside the Finalization Gate. `## References` fully authored (peer cites, source paths, prior art, kata). The five `placeholder:survived` + `gate:inline` findings are all in 1130-1250 (Stage 7's surface). |

## Absorption audit (sources 1 and 3)

Rounds audited: `grounding`, `3amigo` (+`iter-3`), `critique`, `repeatability`, `authors-round.md`.
Result: PASS with one residue. A7 Verified; C1/C2/C4/C5 and S1/S5/S6/S8 rewritten to match round findings; A8/A9 correctly left Pending pre-lock. A3's "FALSIFIED in part" (authors' round) was absorbed by narrowing — A3 is Verified with the narrowing recorded and notes "C2 is unchanged" — **not** a live refutation. Single residue is item 4 (F4). Detail: `absorption-audit.md`.

## Prior art consulted (per `{RDR_RESOURCES}`)

- **House precedent — decisive.** JDR 0001 §D5 posed this exact shape for `<clear>`: an in-band marker, unreserved, indistinguishable from an authored value, yielding "a success-shaped result over unverified state." Prior-art line: "a delete is a write of a marker (Cassandra tombstones; JSON Merge Patch's `null`; `kubectl label k-`) — sound exactly when the marker is **reserved**." It weighed reserve-at-load / separate typed field / typed absent value and **resolved (a)**, rejecting the typed-field arm as "a second re-lock for one field." Lands as "one load category refusing `<clear>` as an authored tag value", peers citing not restating.
- **Mechanism already exists.** `0002:C11` refuses `<clear>` "at load wherever a tag value is authored — a write-block value, an `[initial]` assignment, or a predicate literal — as a `reserved tag value`" — exactly the three sites left open for `<opaque>`.
- **Cost window.** `0029` (Final; cross-cutting owner of `0021:C2`) holds the schema explicitly unstable through `0.x` — "MAY change incompatibly in any release." Cheap now, expensive at 1.0.0. `0002` is `Implemented`.
- **Go practice** (`prior-art-go-sentinels.md`): ingress-time refusal dominates — Prometheus drops `__` labels in `PopulateLabels`; beads retrofitted reserved-prefix rejection after a silent-overwrite bug (GH#2474). Stronger variant: beads derives ids with a delimiter outside the legal input alphabet (collision-free by construction — cf. `0002:C11`'s own `#` rule). Kubernetes annotations are the counter-example (convention only, no ingress refusal).
- **Literature** (`prior-art-sentinel.md`, INCOMPLETE): SQL Antipatterns (p163-167) on in-band sentinels colliding with legitimate domain values; DDIA (p120-124) on protobuf/Avro representing absence structurally. Favors out-of-band. Reconciled: the literature's hazard is the *unreserved* sentinel, which §D5's reservation cures.

## Author decisions (2026-09-12)

1. **Sentinel** → reserve `<opaque>` in `0002:C11` beside `<clear>`, refused at the three authored-value sites. `0021:C2` stands as written and **cites** 0002; marking semantics unchanged. Accepted cost: amends an `Implemented` record (loader + tests) and rejects models that load today.
2. **Quantifier** → hoist to JDR 0001 §JD-23 as a new entry, widening siblings to include 0021, homing the export-side (published merged relation) quantifier beside the lint-side split rule. `0021:C2` cites rather than restates.

## Obligations leaving this stage

| Owner | Obligation |
| --- | --- |
| RDR 0002 (`C11`) | Add `<opaque>` as a second reserved tag value, refused at write-block value, `[initial]` assignment, and predicate literal. 0002 is `Implemented` — amendment + loader tests. |
| JDR 0001 (`§JD-23`) | New entry homing the export-side merged-node terminal quantifier; widen `Siblings:` to include 0021 (README: widen the cluster frontmatter "and say so"). |
| RDR 0021 (`C2`) | Replace the restated quantifier with a citation of §JD-23; cite `0002:C11` for the sentinel reservation; pin `model` to `Model.ID` (item 4). Per rdr-common, cite the owner — never declare in both. |
| RDR 0021 (`A8`) | Evidence narrowed (item 3) — applied this pass. |

## Caveats

- **Determinacy chain disagreement.** `--outcome repeatability` answers `none` (`repeatability-lite-complete`); the chained `resolve:determinacy` routes to `/rdr-prelock repeatability 1` (`determinacy-fired-large`). Resolved on evidence rather than by picking an answer: `run-1.md` line 2 carries `variant: lite (profile: large)` — the exact stamp `emit.surface` requires — and line 1 `model: claude-sonnet-5`, with `diff.md` present and `lens_repeatability_run1=true`. The designated outcome and the disk agree the pass ran, is stamped, and matches the owed variant; the determinacy group appears to route off the fired `Determinacy:` line without crediting the run on disk. Recorded, not silently dropped.
- **A8 packet shape.** The A8 clarification returned `verdict: NEEDS_DECISION` with an imperative in `next_action` where §return-packet requires a question. Substance verified and applied; the malformed field was not treated as licence to act.
