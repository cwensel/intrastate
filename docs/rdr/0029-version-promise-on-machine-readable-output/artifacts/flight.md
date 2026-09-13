# Flight drain — batch:rdr-0029

`/kata-flight --label batch:rdr-0029 --drain`, run after RDR 0029 landed
(implementation `3d86027`, docs flip `895db0d`).

```
flight: 1 shipped, 0 stopped, 0 skipped over 1 wave
  shipped: 5h4a=a626736
  stopped: —
  skipped: —
```

## Pre-flight

Repo anchor `intrastate`, non-bare, primary tree clean, `TARGET_BRANCH=main`
at `895db0d`, only the primary worktree listed (no foreign debris),
`roborev status` and `kata health` both exit 0. `--drain` carries a standing
selector (`--label`), so the re-sweep gate is satisfied.

## Wave 1 — resolution

One eligible kata: `5h4a` (open, unowned, priority 3, no `blocks` edges, none
of the `kind:rdr-seed` / `inbox:hold` / `umbrella` exclusions). No drops, no
cycles.

## Wave 1 — scope review

The review gate ran by default; `5h4a` entered `lifecycle:queued`, not
`lifecycle:reviewed`.

| kata | verdict | evidence |
|---|---|---|
| 5h4a | IN-SCOPE | Defect reproduced at `895db0d`. `internal/cli/schema_register_0029_test.go::TestReq21_TheSnapshotIsComparedAgainstTheCurrentTree` walked `internal/graphlint` and `internal/cli` with no self-exclusion; a detached-worktree probe deleting both genuine readers left it green, and deleting a third matcher — leaving the scanner as the sole match — still left it green. |

No overlaps (sole member; no other open kata cites these files), seam-accretion
count 0 at this locus, no prior demote, no `blocks` edges so no dependent
guard action. `batch:rdr-0029` fits the kept verdict, so §batch-correction
left it in place.

Shape gate: not RDR-SHAPED. REQ-21's disclosure obligation is already decided
by the record and discharged by
`internal/graphlint/vocabulary_snapshot_0029_test.go::TestReq21_TheVocabularySnapshotMatchesTheTree`
(req-list.md:46, `0029:C3`); nothing in `deviations.md` adjudicates the
meta-check. Tightening a test's own scan predicate is a bug against a decided
clause.

Stale fact recorded on the kata: the body cites
`TestReq11And12And13And14...` at lines 211-214. Both are wrong — the
meta-check is `TestReq21_TheSnapshotIsComparedAgainstTheCurrentTree`. The ship
anchored on the symbol.

Disposition amended the filed "redundant scaffolding" reading: the probe was
**kept**, not deleted. It guards a hazard the graphlint oracle cannot — if the
oracle file itself is deleted, the graphlint test vanishes with it and nothing
goes red.

## Wave 1 — ship

| phase | verdict |
|---|---|
| resolve | `ok` — tip `a626736`, one file, 91+/3- |
| rebase + refine | `passed` — rebase clean (already on main tip), 0 iterations, 0 findings, 0 dismissals, lint clean |
| ship | `shipped` — full suite and lint green on the tip, ff-merge `895db0d` → `a626736` |

Scope held to `internal/cli/schema_register_0029_test.go`, exactly the
reported `in_scope_paths`. Squash `n/a` (one commit, no fixups). No
`KATA_PUSH:` spin-offs, no tiebreakers raised.

The shipped fix makes the probe unable to satisfy itself: a self-file
exclusion (`snapshotProbeSelfFile`) plus a read-evidencing predicate
(`evidencesSnapshotRead`) requiring `vocabularies.txt` **and** `os.ReadFile`,
replacing the loose `testdata` + (`vocabular` | `seams`) substring test. Both
design points are load-bearing: the exclusion is required because the scanning
file must literally name the literals it searches for, and the helper plus its
table-driven test live inside that same excluded file so they cannot become a
fresh self-satisfier.

Kata closed with typed evidence at `a626736`; worktree and branch torn down.

## Wave 2 — re-sweep

Standing selector re-resolved to 0 eligible katas, and no spin-off was tagged
onto it during wave 1. Normal termination, not a refusal.

## Post-flight

`main` at `a626736`, primary tree clean, only the primary worktree listed.
The 22 open roborev jobs on `main` are pre-existing debris at unrelated
historical refs — zero for this branch or its tip — and were left untouched.
