Model: claude-opus-5

# RDR 0010 — A1 spike: does the ∅ owned-state root make every group reachable?

**Verdict: A1 VERIFIED — by execution, not only by derivation.**

A1 claims: seeding reachability at the empty owned-state node makes every
group of a decision-table model reachable, so overlap and coverage produce
the same findings the scratch-tag PoC produced.

The spike reproduces the PoC baseline, shows the current zero-owned failure,
and then EXECUTES C5's proposed ∅ root without editing `internal/`.

## Reproducing

The harness is gated behind a build tag so it is evidence rather than a
member of `make check` (it reads `internal/` internals that implementation
will change):

    go test -tags rdr_spike ./docs/rdr/0010-stateless-decision-tables/evidence/spikes/emptyroot/ -v

## Model shape

`a1-baseline.toml` — 9 `bool` observed dimensions (`d0`..`d8`, each
`single_valued = true`, `required = true`) → **2^9 = 512 cells**. One scratch
owned tag `phase` (enum `start`/`done`) with `[initial] phase = "start"`,
`[read.scratch]` / `[write.scratch]` accessors, a `[context.done]` terminal,
and a `[rule.write] phase = "start"` on every rule — the PoC convention.

The rule set is deliberately PARTIAL. Three broad rules cover every
assignment except `d0=true AND d1=true AND d2=true`:

| rule | guard | cells |
|---|---|---|
| `r001` | `d0=false` | 256 |
| `r002` | `d0=true, d1=false` | 128 |
| `r003` | `d0=true, d1=true, d2=false` | 64 |
| `r004-dims` | all nine `= false` | (inside r001; pins d3..d8 as participating dimensions) |

Covered = 448. **Hole = 64 of 512.** (The PoC's 208/512 was a different
partition; A1's own Evidence clause asks for the same finding SET
cell-for-cell across arms, "not the count", which is what this spike tests.)

## Part 1 — baseline (scratch owned tag)

    $ bin/intrastate lint --model .../a1-baseline.toml --as=json     → a1-baseline.out

    graph-coverage-gap  blocking
      group a1-spike-baseline/decide over rows [r001 r002 r003 r004-dims]
      leaves 64 of 512 assignments in its scoped product uncovered for the
      no_match arm
      dimension: d0,d1,d2,d3,d4,d5,d6,d7,d8   rule: r001
    EXIT=2

POSITIVE coverage finding. **64/512 uncovered, 448/512 covered.**

Note the group context is `a1-spike-baseline/decide` — it carries NO owned
match atom, because the scratch tag never appears in `[rule.match]`. That is
exactly the shape whose reachability A1 is about.

## Part 2 — zero-owned-tag arm

### 2a. Deleting the owned tag, `[initial]`, accessors AND the write blocks

    $ bin/intrastate lint --model .../a1-zero-owned.toml --as=json   → a1-zero-owned.out

    model-invalid / malformed_rule_shape:
      "ordinary rule r001 carries no write block"
    EXIT=2

**Lint is never reached.** The LOADER refuses a zero-write ordinary rule
first. This is the gap 0010:S2 ("rule-shape normalization over an ordinary
rule with no write block") must close before C5 matters at all — the
zero-owned model does not currently load.

### 2b. Setting `class = "decision-table"` (A8's refusal, confirmed)

    $ bin/intrastate lint --model .../a1-zero-owned-class.toml --as=json
                                                        → a1-zero-owned-class.out

    model-invalid / unknown_schema_field:
      "strict mode: fields in the document are missing in the target struct"
    EXIT=2

Strict decoding refuses the unknown `[model].class` key exactly as **A8**
predicts — `unknown_schema_field`, never a silent load as a machine.

### 2c. No root, writes retained (isolating the reachability effect)

    $ bin/intrastate lint --model .../a1-no-initial.toml --as=json   → a1-no-initial.out

    graph-dangling-edge  blocking
      "the model declares no initial owned state"
    EXIT=2

**The 64/512 coverage finding DISAPPEARS.** This is the RDR's own analysis
confirmed by execution: `reach` returns no nodes, `contextReachable` is false
for every group, `checkGroups` `continue`s before `checkCoverage`, and
coverage is VACUOUS. The only survivor is the missing-root arm. This is the
false green that motivates C5.

## Part 3 — executing the ∅ root

C5 cannot be executed through the CLI (the class key does not decode), and
`internal/` may not be edited. But `graphlint.Run` / `NewRequest` / `Reach`
are EXPORTED, and `docs/` is inside the same Go module, so a spike test under
this directory can drive the engine over a hand-mutated `*table.Model`.

`emptyroot/emptyroot_test.go` loads `a1-baseline.toml`, then `strip()`s it —
deletes every owned tag, `m.Initial`, `m.Terminal`, and every row's `Writes`,
`NextTags`, `RequiresOwned` — producing the true zero-owned decision table
in memory. Three arms:

| arm | root | nodes | findings |
|---|---|---|---|
| A — baseline as loaded | `phase=start` | 1 | `graph-coverage-gap` (64/512) + `graph-redundant-row` |
| B — stripped, no root (CURRENT) | none | **0** | `graph-dangling-edge` only — **coverage vacuous** |
| C — stripped, ∅ root (C5) | `Node{Values:{}}` | **1** | `graph-coverage-gap` (64/512) + `graph-redundant-row` |

Output: `a1-emptyroot.out`.

The ∅ root is seeded WITHOUT editing `internal/`: `reach` bails only on
`len(m.Initial) == 0`, and then SKIPS every root assignment whose tag is not
owned-provenance (`reach.go:104-110`). So `m.Initial = [{Key:"d0"}]` — an
OBSERVED key — passes the guard and builds `root = Node{Values: {}}`, which
is literally the empty owned-state node C5 names. The test asserts
`len(nodes)==1 && len(nodes[0].Values)==0` before running lint.

**Arm A and arm C emit the identical coverage finding**, byte-for-byte on
every field except the model id — same code, same 64-of-512 message, same
dimension list, same rule, same element, same class, same severity. Verified
programmatically (see the comparison at the end of the spike transcript).

The one difference between A and C is the `graph-redundant-row` advisory's
FINGERPRINT: arm A's ends `...;#phase=start,;` and arm C's ends `...;#`.
The write block is part of `Fingerprint`. This does not affect A1 (the
coverage finding carries no fingerprint), but it is a data point for **A7**:
`Fingerprint` DOES read the write block, so a row's fingerprint changes when
its writes are removed. A7's actual claim is narrower — that `emit` is not
read — and this spike neither confirms nor refutes that.

## The source trace behind it

### (i) Does `nodeSatisfiesMatch` consult only owned atoms?

`internal/graphlint/analysis.go:78-91`:

```go
func (a *analysis) nodeSatisfiesMatch(n Node, match []table.Atom) bool {
	for _, atom := range match {
		if a.model.Tags[atom.Key].Provenance != table.ProvenanceOwned {
			continue
		}
		if !ownedAtomSatisfiable(n, atom) {
			return false
		}
	}
	return true
}
```

**Yes.** Every non-owned atom is skipped outright. Over a model with ZERO
owned tags, no atom can have owned provenance, so the loop body never runs
its test and the function returns `true` for EVERY match pattern, over ANY
node — the ∅ node included. This is unconditional, not incidental.

### (ii) Does `checkGroups` gate coverage on group reachability, and is that decided by "any reachable node satisfies the context"?

`internal/graphlint/groups.go:24-44`: `checkOwnedBeforeMatch` and
`checkVacuousAtoms` run for every group; then

```go
if !a.reachable[g.Context.String()] { continue }
a.checkOverlap(g); a.checkRedundantRows(g); a.checkCoverage(g)
```

So overlap, redundant-row, and coverage are ALL gated on
`a.reachable[ctx]`. That map is built once in `newAnalysis`
(`analysis.go:53-57`) from `contextReachable` (`analysis.go:66-73`):

```go
for _, n := range a.nodes {
	if a.nodeSatisfiesMatch(n, g.Context.Match) { return true }
}
return false
```

**Yes to both** — it is an existential over `a.nodes`. Combining with (i):
with one node in `a.nodes` and no owned tags in the model,
`nodeSatisfiesMatch` returns true immediately, so EVERY group is reachable.
With ZERO nodes the loop body never executes and every group is unreachable
— which is arm B.

### (iii) Reachability-dependent gates a one-node graph could fail

There are exactly two size gates, and neither is reachability-dependent —
both are decided over the group's DECLARED product, not over nodes:

**Product bound — `internal/graphlint/coverage.go:50`.**
`ProductBound()` = `guard.Bound()` = **2048** (`internal/guard/product.go:420`,
surfaced at `internal/graphlint/taxonomy.go:139`). The test is
`card > ProductBound()`.

> **512 ≤ 2048 → NOT withheld.** Confirmed empirically: `a1-thresholds.out`
> shows a 12-dimension (4096-cell) variant of the same table taking
> `graph-product-too-large` ("4096 assignments, above the published product
> bound of 2048"), while the 512-cell table does not. The single group over
> ~512 cells is REPORTED, not withheld.

**"Unprojectable group size" suppression.** There is no size-keyed
suppression distinct from the bound. What exists is the UNPROVABLE-DIMENSION
suppression, and it is keyed on DECLARATIONS, never on size:

- `emitUnprovableDimensions` (`coverage.go:106-129`) → `unprovableReason`
  (`coverage.go:178+`): fires when a dimension has no finite declared domain
  (`guard.AssignmentCount` returns `!ok`) or — per `0003::A21` — when a
  value atom (`eq`/`in`/`lt`/`lte`/`gt`/`gte`) sits over a tag lacking
  `single_valued`.
- `emitStructurallyUnprovable` (`coverage.go:147-174`): the same
  single-valued test, run on the over-bound path only.
- `guard.Product` (`internal/guard/product.go:182-186`) returns a
  non-projectable empty set when `!ok || card > Bound()`, and
  `emitCoverageArms` (`coverage.go:334-337`) returns silently on
  `!product.Projectable()`.

> **Threshold: not a size at all — a declaration predicate.** Every `d0`..`d8`
> is `kind = "bool"`, `single_valued = true`, `required = true`, so every
> dimension is finite and projectable and NOTHING is withheld. Confirmed
> empirically: `a1-nonsinglevalued.out` (same table, `single_valued` dropped
> from `d5`) takes `graph-unprovable-coverage` / `tag-not-single-valued`
> instead of the coverage gap. **A 512-cell single group of properly declared
> dimensions is NOT withheld.**

Crucially, in BOTH suppression cases the outcome is a BLOCKING finding
(`graph-product-too-large`, `graph-unprovable-coverage`), never a silent
green. The failure mode A1's "If wrong" clause fears — "exit 0 with an
unproven table" — is only producible by the arm-B reachability hole, which
is precisely what C5 closes.

### (iv) Node ceiling

`nodeCeiling = 4096` (`internal/graphlint/taxonomy.go:128`, exposed by
`NodeCeiling()` at `:144`). `reach` breaks and sets `complete = false` when
`len(nodes) > nodeCeiling` (`reach.go:116-119`); `checkNodeCeiling`
(`analysis.go:169-181`) emits `graph-product-too-large` against the traversal
only when `!a.complete`.

**A one-node graph trivially passes**: 1 ≤ 4096, the loop never breaks,
`complete` stays `true`, and `checkNodeCeiling` returns immediately.
Confirmed in arm C — `nodes=1`, no traversal finding.

## The exact answer

**Does the ∅ root make every group reachable? YES — unconditionally, for a
model with zero owned tags.**

The chain is closed and has no escape hatch:

1. `reach` with a non-empty `m.Initial` always produces at least the root
   node, so `a.nodes` is non-empty (`reach.go:112`).
2. A model with zero owned tags has no owned-provenance atom, so
   `nodeSatisfiesMatch` skips every atom in every match pattern and returns
   `true` (`analysis.go:78-91`).
3. `contextReachable` is therefore true for every group, because it is an
   existential over a non-empty `a.nodes` with an always-true predicate
   (`analysis.go:66-73`).
4. `checkGroups`'s `if !a.reachable[...] { continue }` never fires, so
   overlap, redundant-row and coverage run for every group
   (`groups.go:36-43`).
5. `checkCoverage` reads NO node — its doc comment is explicit that
   "coverage is a UNIVERSAL claim … computed from the group's authored rows
   and their declared guard domains alone and never from a reachability
   node. Reachability decides only whether the group is proven at all"
   (`coverage.go:28-32`). So once the gate opens, the finding set is a pure
   function of the rows and declarations — identical to the baseline's.

Step 5 is why arms A and C agree cell-for-cell: the owned tag was never an
input to coverage, only to the gate in front of it.

Note the stronger-than-stated result: the ∅ root makes every group reachable
because there are no owned atoms to fail, NOT because the ∅ node happens to
satisfy them. So the claim holds for any zero-owned model whatever its match
patterns — it does not depend on the table's shape.

## Caveats / what this spike does NOT show

- The 208/512 figure is not reproduced; this spike's partition gives 64/512.
  A1's Evidence clause asks for the same finding SET across arms, which is
  what was tested and what held.
- `class = "decision-table"` cannot be exercised; arm C simulates C5's
  ROOTING only. Whether the implementation reaches the ∅ root via the class
  key, and whether C5's "MUST NOT emit" list (dead end, always-present,
  owned-before-match, single-valued, terminal-escape) really stays silent,
  is A10's job — though arm C emitted none of them.
- Zero-write rules do not load today (`malformed_rule_shape`), so 0010:S2 is
  a hard prerequisite for any of this reaching lint through the CLI.

## Files

| file | what |
|---|---|
| `gen.py` | generates the two arms |
| `a1-baseline.toml` / `.out` | Part 1 — 64/512 coverage gap, EXIT=2 |
| `a1-zero-owned.toml` / `.out` | Part 2a — `malformed_rule_shape`, lint not reached |
| `a1-zero-owned-class.toml` / `.out` | Part 2b — `unknown_schema_field` (A8) |
| `a1-no-initial.toml` / `.out` | Part 2c — coverage vacuous, only `graph-dangling-edge` |
| `emptyroot/emptyroot_test.go` | Part 3 — arms A/B/C over the exported engine |
| `a1-emptyroot.out` | Part 3 transcript |
| `a1-thresholds.toml`s / `a1-thresholds.out` | bound 2048 and single-valued threshold, empirical |
