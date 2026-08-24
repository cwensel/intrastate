model: claude-opus-5
stage: 06-reconcile (re-entry pass — rewritten after the repeatability lens)

# Reconcile Report — RDR 0005

The earlier reconcile pass (recorded here previously) predates the
repeatability lens. Stage 5's row grew after it: the Determinacy trigger fired
on this `mid` RDR, `repeatability` (lite) ran, and its diff/resolve mutated the
draft and opened A7. This pass rebuilds the open set over the current draft.

## Stage 5 completeness preflight

- `Profile`: `mid` → required row `grounding → 3amigo`. Both complete:
  `evidence/grounding/findings.md`, `evidence/3amigo/` (consolidation +
  3 personas + resolve, iter-2, iter-3).
- Determinacy trigger fired → `repeatability` (lite) appended to the row.
  Discharged by evidence, not waiver: `evidence/repeatability/run-1.md`
  (`variant: lite (profile: mid)`, `model: claude-sonnet-5`) plus
  `diff.md` and `resolve.md` (`model: claude-opus-5`).
- Variant check: `run-1.md` reads `lite`, and no `run-2`/`run-3` sit beside
  it — lite ran where lite was owed. No variant mismatch.
- Stage 5 complete. Reconcile proceeds.

## Pre-Lock needs-verification inputs

- `evidence/3amigo/resolve.md`, `iter-2/{delta-review,resolve}.md`,
  `iter-3/{delta-review,resolve}.md`
- `evidence/repeatability/resolve.md` (§Needs verification — the new input
  this pass adds over the prior report)

## Absorption audit

- O1–O5 (3amigo iters 1–3): absorbed, as recorded in the prior pass. Iteration 2
  reported no open findings; iteration 3 left A4, which the prior pass verified.
- Iter-3 residue re-checked against the current draft: iter-3's wording
  ("gates … before calling the pure resolver kernel") is *superseded*, not
  contradicted — A4 records the June fencing as replaced by JDR 0001 §D8/§D9,
  and the body is consistent on post-selection gating (§Technical Design line
  ~445, A4, and the §D9 citation). Absorbed, no residue.
- Repeatability R-1..R-4: all four `fixed` in the draft and confirmed present
  in this pass — R-1 `revision` pin (§Technical Design *Model selection*),
  R-2 `--as` cite-don't-restate (§Technical Design *Envelope* tail),
  R-3 narrowed reader set + R-4 `next` gate scope (§Normative Contracts),
  plus the §Failure Modes disposition table and MVV/Testing scenario 8.
  No unabsorbed finding survives.

## Open set

| Item | Source | Disposition | Evidence pointer or plan |
| --- | --- | --- | --- |
| A7 normalized model exposes reader→owned-key mapping and per-row gate lists, pre-evaluation | 1 (repeatability needs-verification), 2 (Status: Pending), 4 (`only` / `every` / `no others` in the two new narrowing clauses) | **VERIFIED** | Peer RDR, delegated read of RDR 0002/0004. Leg (a): `0002::Normative Contracts` *Accessor tables* — entries carry `keys` ("a non-empty list of declared tag keys"), and "every **owned** tag MUST be served by exactly one reader"; `0002::Technical Design` rejects at load "accessor `keys` that bind a tag to zero or several readers"; `0004::Normative Contracts` — "MUST declare the requested key set as validated metadata. The set MUST NOT be derived from the keys a read actually resolved." Leg (b): `0002::Normative Contracts` *Gate references* — "the list is carried on the normalized row and is part of its value"; `Row.RequiresOwned` supplies the per-row owned demand. Both are dumpable normalized-row columns (`… requires_owned, gate, escape`), produced by a "Load" pipeline that refuses before yielding candidate rows. Written into A7 and the RDR's Reconciliation Report. |
| R-1 `revision` pin — new exactness claim ("`revision` on every payload", carried "verbatim") | 4 | **VERIFIED** (covered in-body) | Grounded on a live symbol, re-checked this pass: `internal/resolve::Table.Revision` at `internal/resolve/resolve.go:205-207` carries the quoted doc comment verbatim and is echoed at `:508,:516`. Ownership sits with RDR 0002's loader; the CLI never derives it. |
| R-2 `--as` default | 4 | **ACCEPTED** (Design Decision, in-body) | Cite-don't-restate: the default stays the root contract's. Confirmed live — `internal/cli/root.go:51` wires `respond.FlagName` with default `"text"`. The rejected alternative (restating "default text" in this RDR) is the single-source drift the disposition names. |
| Named spikes with no captured run | 3 | **NONE OPEN** | All five spikes cited in the body have output: `d13-canonical-set.out`, `escaping-surfaces.out`, `finding-producer-matrix.out`, `finding-shape-options.out`, `finding-shape-tiebreak.out`. |
| A1–A6 | 2 | **UNCHANGED — VERIFIED** | Repeatability invalidated none: R-1/R-2 resolved *toward* already-verified surfaces; R-3/R-4 added scope clauses inside A3/A4's verified verb/accessor boundary rather than moving it. Re-confirmed against the current draft. |

## MVV deferral check (hard rule)

A7 pins behavior the **MVV itself proves** — the MVV requires the unneeded-reader
pair ("the only assertion that distinguishes the narrowed invoked set from
'every declared reader'") and the `next --evaluate-gates` scope assertion. A7 was
therefore **not deferrable**: DOWNGRADED was unavailable and it had to reach
VERIFIED before lock. It did.

## Refutation check (hard rule)

No spike or source-search refuted an assumption the RDR relies on. A7's fallback
branch ("run every declared reader; widen `flow-artifact-missing`") is **not**
triggered — the peer text supports the narrow reading the draft pins. No BLOCKER,
no route-back, no §punt-ledger row owed.

## Completeness check

- No `_Draft placeholder._`; no `this is a seed skeleton` header; no surviving
  template bracket (`[Required` / `[Conditional` / `[Resource]` / `[Capability]`).
- No `Status: Pending` or `Unverified` assumption remains.
- `## References`: no bracketed template placeholder. Collected this pass —
  the stale "the escaping edge is open" annotation on `d13-canonical-set.out`
  (settled by A6) corrected, `escaping-surfaces.out` added (cited normatively in
  three body sites but previously unlisted), and the two anchors the
  repeatability round introduced added: `internal/cli::NewRootCmd`,
  `internal/resolve::Table.Revision`. Both resolve on `main`.

## Other edits landed this pass

- §Finalization Gate *Scope Verification* enumerated the MVV without the two
  assertions the repeatability round added — stale against the MVV body, and it
  would have read as a scope reduction at Stage 7. Enumeration completed.
- §Finalization Gate *Assumption Verification* stub said "after A3–A6 re-verify";
  now also names A7's Stage 6 verification. (Stage 7 still rewrites both.)

Verdict: **RECONCILED**

Next: `/rdr-finalize 0005`
