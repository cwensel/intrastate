# Flight — batch:rdr-0008

`/kata-flight --label batch:rdr-0008 --drain`, run as the `§land-flight` tail of
`/rdr-implement-triage 8 --close-and-flight`.

```
flight: 3 shipped, 0 stopped, 1 skipped (over 1 wave)
  shipped: wagv=1959384  wtne=ce3e158  sf4e=01c9ee7
  stopped: —
  skipped: jb9v=closed-at-scope-review (over-eng)
```

## Wave 1 — resolution

4 katas resolved from the standing selector, all eligible (unowned, no
`kind:*`, no `inbox:*`, no `umbrella`, no open `blocks` edges). No intra-batch
dependency edges, so topo-sort reduced to priority then short_id:
`wagv` (P3) → `wtne` (P3) → `jb9v` (P4) → `sf4e` (P4).

## Scope-review gate (default-on, ran before the ship loop)

```
scope-review: batch:rdr-0008  reviewed: 4
  in-scope:     wagv wtne sf4e   (lifecycle:reviewed, plans embedded)
  out-of-scope: jb9v=over-eng
```

No overlap clusters: all pairwise overlaps were judged compatible-but-distinct
— three disjoint REQ clusters at three code sites (`internal/cli` carrier,
`internal/table` declaration guard, `internal/table` advisory trigger).
Consolidating would have produced one unit whose DoD spans unrelated oracles.
No seam-accretion tripwire (no repeated point-fixes at these loci).

**`jb9v` closed at review (over-eng).** No mutation exists where a looped
assertion fails and the current single-load version passes; `load.go` scans
`slices.Sorted(maps.Keys(decls))`, so the payload is fixed by construction, and
commit `0ead6e5` had deliberately *replaced* a 200-sample histogram with the
exact assertion. Re-adding repetition would reinstate the rejected construct.
Its `batch:*` labels were stripped per §batch-correction.

A factual correction was recorded on close: the kata body (inherited from the
triage note) claimed the pre-fix map range "returned the same wrong entry every
run". The record measured **161/39 over 200 loads** — ~20%-per-load detectable,
not 100%. The exactness of the assertion, not the loop count, is what makes the
oracle sound.

## Shipped

### `wagv` → `1959384` (2 commits) — `area:internal-cli`, severity medium

`lint.go` and `flow_input.go` now carry load categories and the RDR 0008
reserved-key payload as typed `clierr.Findings` through a shared `loadFindings`
helper; near-miss advisories ride the success payload.

The kata's own premise dissolved under review: it asked for a 0005/0006 joint
decision on the carrier field, but **JDR 0001 §D10 item 3 had already made it**,
and both carriers (`Finding.Rule`, `lintPayload.Findings`) already existed.
Zero new public surface — a wiring defect, not a design fork.

Two findings caught during the run, both real:
- Resolve's own first draft echoed the offending name (`recognized`) in the
  rename hint; the REQ-30 test it wrote caught it. The offending name travels
  on `param` instead.
- Refine found `nearMissFindings` setting `Code: table.CatReservedTagKey`,
  which **REQ-63 forbids** (the advisory must carry no category discriminator
  and must not participate in category dispatch). Fixed to `Code: a.Rule`. The
  agent corrected its own resolve-phase design note rather than defending it.

One roborev job (#6265) was stale-open at the pre-merge gate — it reviewed
`eaf063b`, before the refine fix landed in `1959384`. Verified independently
that `lint.go:149` reads `Code: a.Rule`, exactly the remedy requested, then
closed with that evidence.

### `wtne` → `ce3e158` (1 commit) — `area:internal-table`, severity medium

Added `testdata/neg/neg-recognized-observed.toml` plus category and payload
assertions pinning REQ-11/REQ-30's **observed-provenance** half of the reserved
key. Refine surfaced zero findings.

Deliberately a *third* fixture rather than a mutation of
`neg-recognized-owned.toml`, which nine call sites depend on — flipping it would
have silently retargeted them and lost the owned-side pin. Fixture uses the
minimal legal observed shape (`kind = "scalar"`); a probe confirmed the guard
fires in `loadTags` before accessor binding.

Mutation-verified: narrowing the guard to `ProvenanceOwned` failed all three new
assertions while the entire pre-existing suite stayed green; setting `Remedy` to
`RecognizedTagKey` failed the new payload row. Both reverted.

### `sf4e` → `01c9ee7` (1 commit) — `area:internal-table`, severity low

Added `testdata/pos-near-miss-both.toml` with a quoted leading-space key
(`[tags." Recognized"]`) plus table cases pinning REQ-59's **combined**
fold-and-trim clause. Refine surfaced zero findings.

Mutation-verified: under a disjunctive `isNearMiss`
(`EqualFold(key, r) || TrimSpace(key) == r`) only the new case goes red while
all four pre-existing cases stay green — the unique discriminating witness the
scope review predicted by direct predicate evaluation.

## Drain

Wave 1 tagged no `KATA_PUSH` spin-offs onto the standing selector (every ship
reported `pushed_to_kata: empty`), and the wave-2 re-sweep resolved empty
(0 open `batch:rdr-0008` members). Normal termination after one wave.

## Invariants held

Every ship ran in its own worktree on `worktree-<short_id>`, ff-merged, and tore
down clean; the primary checkout stayed clean and on `main` at every
between-katas gate. All three shipped katas verified closed + unowned +
phase-label-free by reading kata state, not by trusting the ship narration.
The two test-only katas kept a strict zero-production-`.go`-change diff.
