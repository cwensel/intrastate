Model: claude-opus-5[1m]

# Stage 6 Reconcile — RDR cli/0011

Verdict: **RECONCILED**.

## Stage 5 preflight

`Profile: large` → lens row `grounding → 3amigo → critique`, all present with
iterations (`grounding/` + `iter-2`, `3amigo/` + `iter-2` + `iter-3`,
`critique/` + `iter-2`). The Determinacy trigger fired (C1 is an algorithmic
contract: a per-key probe filter, a three-source merge with a pinned dedup
order, and a sorted payload), and is discharged by
`evidence/repeatability/run-1.md` — header `variant: lite (profile: large)`,
the correct variant for `large` — plus `diff.md` and `resolve.md`. No variant
mismatch: lite ran where lite was owed, and no stray `run-2`/`run-3` sits
beside it. Preflight PASS.

## Open set

Built from the four sources. Sources 2 and 3 came back empty; source 1 (the
pre-lock needs-verification lists) resolved to items already closed by the
rounds that raised them; source 4 (the post-mutation exactness delta) was
clean. The live residue came from the absorption audit — four findings that
`critique/iter-2/diff.md` listed as "absorb in-draft on the way back through"
which the re-propose did not land.

| item | source | disposition | evidence pointer / plan |
| --- | --- | --- | --- |
| A1–A6, A11–A16 (all 12 CAs) | 2 | VERIFIED (no change) | All `Status: Verified` before this stage; spot-checked against the projector. Source-2 open set is empty. |
| Spikes a5, a6, a12, a14, a15 | 3 | VERIFIED (no change) | Every spike named in the RDR has a captured run under `evidence/spikes/`. No named-but-unrun spike. |
| A7 (minted by `grounding/findings.md` iter-1) | 1 | ACCEPTED — retired, not dangling | The escape-site census died with Alternative 3 at the critique iter-2 route-back; `internal/resolve` is now untouched (C1, S5, MVV 9). `grounding/iter-2/findings.md:8` records the drop. |
| A12 (flipped Pending by `critique/iter-2/resolve.md`) | 1 | VERIFIED — closed before this stage | `evidence/spikes/a12-per-key-match-filter.md`; zero-splits detector over 28 row×view observations, `eq`+`in` pairing loaded live. |
| A15 (Pending, scope corrected, `critique/iter-2` + `3amigo/iter-3`) | 1 | VERIFIED — closed before this stage | `evidence/spikes/a15-demand-set-match-owned.md`; breaking arm (plan → `flow-artifact-missing`) reproduced live. |
| A3 ("re-read when C1 is re-decided") | 1 | VERIFIED — re-read, no flip | Evidence now carries the dead-row-vs-conflicted-key distinction explicitly; A12's spike confirms the conflicted path stays CLI-unreachable. |
| A4, A11 ("re-read, not flipped", twice) | 1 | VERIFIED — no flip | A4's BREAKS-0 is predicate-scoped and untouched. A11's verdict stands; its Evidence prose corrected this pass (see below) and its property now pinned by an oracle (R1). |
| R1 — A11's key-set agreement holds only by the current call graph; no scenario fails when it breaks | audit residue (`critique/iter-2`, B C-3) | VERIFIED — oracle added | S5 now mandates comparing `assembledView`'s key set against `internal/resolve/resolve.go::assemble` modulo `recognized`. New mini-check row `S5 key-set`. A11's Evidence + If-wrong updated. |
| R2 — C1's "filter after `KernelRow()`" is prose-only; no fixture uses a set-kinded match key, so filter-before is observationally identical | audit residue (`critique/iter-2`, B C-6) | VERIFIED — oracle added | S3 now mandates a set-kinded match key, which is the shape `seamValue`'s canonicalization (sorted, compacted JSON, `SetEscapeHTML(false)`) separates. New mini-check row `MVV 6 / S3 set-kinded`. |
| R3 — no fixture exercises several match keys on ONE row in mixed states; A12's spike records this as limit (c) | audit residue (`critique/iter-2`, B C-8) | VERIFIED — oracle added | S3 now mandates a row carrying present-equal + present-unequal + absent keys at once. New mini-check row `MVV 6 / S3 multi-key mixed`. A12's limit (c) amended to record the discharge. |
| R4 — C3's census counts payload reads and production symbols but not the `unresolved` comment / `t.Errorf` prose | audit residue (`critique/iter-2`, A C-5 second half) | VERIFIED — census clause added | Occurrences confirmed live on `main` at `flow_next_0005_test.go:23,66,104,135,137,141,169-171,404-408`, `flow_mvv_0005_test.go:46,73`, `flow_adversarial_0005_test.go:448,450`. C3 now requires the prose sweep, explicitly outside the five-read count. |
| R5 — the rename census is blind to out-of-repo consumers | audit residue (`critique/iter-2`, B C-5) | ACCEPTED — already recorded | Consequences already states an out-of-repo consumer parsing `unresolved` must change even if it adopts `--all`. `intrastate` is a generic library by design; an in-repo census cannot see external parsers, and that is a stated scope limit, not a verifiable fact. |
| `0007:REQ-78` (two-valued match seam) left deferred | 4 | ACCEPTED — Design Decision | Decision Rationale names the rejected alternative (a resolution-level veto on undecidable match, `0011:ALT4`) and records the cost in Consequences: `next` and `resolve` disagree on an absent-key row, in the escapable direction. |
| A5 limit (i) — a genuine zero-owned 0010 decision table is unloadable on this build | 4 | ACCEPTED — scoped to 0010 | The "empty `required`" half of `0010:A11` is downstream of 0010 shipping and not attestable here; recorded in A5's Evidence rather than papered over. Belongs to 7.1 cluster-reconcile, not Stage 6. |
| A15 limit (b) — shipped fixtures attested by suite-green, not payload diff | 4 | ACCEPTED — carried as a build obligation | S8 already carries it forward: the build owes the payload-diff form the spike could not reach. Stated limit, not an open spike. |

## Exactness-word delta (source 4)

Only the post-mutation delta is in scope; Stage 4 swept the rest. The
repeatability round was the last to mutate the draft, landing four pins
(F1 gate-id source `Row.Gate`; F2 `not-evaluated` emitted on every reported
disposition; F3 the payload read guarded on `guard_unevaluable`; F5 dedup-last).
All four are in C1/C2 as written, and each new exactness claim they introduced
carries a record: S2's filter/dedup-order oracle grounds on the
matched-and-written `stage` key in `models/rdr.toml`; S7's corrected dedup
scope is reasoned in place; S8's "byte-identical" claim cites the A15 spike,
which ran the comparison live and names its own limit. No unrecorded delta.

## Completeness

No `_Draft placeholder._`, no seed-skeleton header, no surviving template
bracket. `## References` is fully authored (peer RDRs, source paths, three
prior-art ledgers, spike and lens evidence, trackers). `rdr lint 0011` = PASS;
its one `conformance gate:inline` note is the Stage 7 lock-time move of the
gate body to `artifacts/gate.md`, not a Stage 6 item.

## No blockers

Nothing refuted. Every residue item is additive — a scenario or census clause
pinning a fact the RDR already establishes — so no route-back to Propose,
Refine, or Resolve is owed, and no §punt-ledger row is due. No MVV-critical
spike or assumption is deferred past lock: the two limits that touch MVV
oracles (A12 limit (c), A15 limit (b)) are now either discharged by a mandated
fixture or carried as a named build obligation.
