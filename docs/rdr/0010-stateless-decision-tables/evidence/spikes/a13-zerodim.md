Model: claude-opus-5[1m]

# A13 — zero-participating-dimension group: emission site and fixture sweep

Assumption under test (RDR 0010, A13):

> A decision-table group with zero participating guard dimensions can take
> `graph-unprovable-coverage` from inside `checkCoverage`'s `len(dims) == 0`
> branch, ahead of `emitCoverageArms`, without disturbing the state-machine
> class, where a zero-dimension group is legitimate and stays silent.

Read through the projector only:

```
/Users/cwensel/sandbox/newcoinc/rdr/bin/rdr inspect --select 0010:A13 0010
/Users/cwensel/sandbox/newcoinc/rdr/bin/rdr inspect --select 0010:C5 0010
```

Method: a Go spike package under this evidence dir, driving the REAL
`internal/table` loader and the REAL `internal/graphlint` engine through
`graphlint.Run(graphlint.NewRequest(m))`. Nothing under `internal/` was
edited. Spike source and captured output:

- `evidence/spikes/a13-zerodim/a13_zerodim_spike_test.go`
- `evidence/spikes/a13-zerodim/a13-spike.out`

## Baseline

```
$ go test ./internal/graphlint/... ./internal/table/...
ok  	github.com/newcoinc/intrastate/internal/graphlint	3.370s
ok  	github.com/newcoinc/intrastate/internal/table	0.502s
```

GREEN, both before and after the spike (`internal/` is clean per
`git status --porcelain internal/`, which returns nothing).

## Q1 — does an emission ahead of `emitCoverageArms` double-report?

### What `emitCoverageArms` does in the zero-dimension case

Traced at `internal/graphlint/coverage.go:332` (`emitCoverageArms`) with
`len(dims) == 0`:

1. `guard.Product(a.model, g)` is the EMPTY product — one empty assignment —
   and it IS projectable, so the early `!product.Projectable()` return does
   not fire.
2. For each rescuable class, `coverageUnionFor` takes its own
   `len(guard.Dimensions(...)) == 0` arm (`coverage.go:415`): membership
   alone decides, so any row in the closing population makes the union
   equal the product. **Every zero-dimension group closes every arm.**
3. Because the arm closed, `bareEscapeFor(g, class)` runs. An escape row
   carrying no guard atoms sets `closedBy`, and the function then emits
   `graph-coverage-closed-by-escape`.
4. No `graph-coverage-gap` can ever accompany it: the union equals the
   product by construction, so the gap branch is unreachable in this case.

### Demonstrated

`TestQ1_ZeroDimGroupWithBareEscapeTakesClosedByEscapeToday` — a
zero-dimension group (`dims=[]`) with an ordinary row plus a bare escape row:

```
groups:
  group "t/go status.eq=a" rows=[ordinary rescue] dims=[] zero=true
report:
  graph-coverage-closed-by-escape sev=info rule="rescue" element="t/go status.eq=a"
    msg=the coverage of group t/go status.eq=a is closed by the bare escape row
        "rescue" rather than proved over its declared domains
Q1: closed-by-escape count over zero-dim group with bare escape = 1
```

`TestQ1_ZeroDimGroupWithoutEscapeIsSilent` — the same group with no escape
row: `report: (no findings)`.

`TestQ1b_ZeroDimBranchEmitsOnlyTheClosureArm` — over a zero-dimension group,
`emitCoverageArms` produces the closure advisory and NOTHING else from
invariant 4 (no gap, no unprovable-coverage, no product-too-large), in both
the bare-escape and no-escape shapes.

### Verdict on Q1: PASS, with the precedence being load-bearing

A naive insertion — emit, then fall through to `emitCoverageArms` — **WOULD
double-report**: over a zero-dimension group carrying a bare escape row the
run would carry both the new `graph-unprovable-coverage` and the existing
`graph-coverage-closed-by-escape`, two findings for one group. That is
precisely the outcome C5 forbids with "Where both would apply, this finding
is reported and `graph-coverage-closed-by-escape` MUST NOT be".

The precedence C5 states IS achievable at that exact site, and cheaply. The
branch at `internal/graphlint/coverage.go:34-42` is a self-contained
`a.emitCoverageArms(g); return` pair. The class-keyed arm therefore takes
the shape:

```go
if len(dims) == 0 {
    if <declared class is "decision-table"> {
        a.emit(clierr.Finding{
            Code:   CodeUnprovableCoverage,
            Reason: ReasonNoParticipatingDimension,   // A14's append
            ...
        })
        return                                        // suppresses the closure arm
    }
    a.emitCoverageArms(g)
    return
}
```

The `return` is what delivers C5's precedence. Because Q1b establishes that
the closure advisory is the ONLY invariant-4 output over a zero-dimension
group, returning early suppresses exactly that one finding and forfeits
nothing else — no gap, no withholding, no bound refusal is lost, since none
of them can arise on this path. The state-machine path is byte-identical to
today's: the `else` limb still calls `emitCoverageArms(g); return`.

One dependency stands unresolved and is already booked: the `reason`
discriminator `no-participating-dimension` does not exist in
`internal/graphlint/taxonomy.go` — the closed set there is exactly
`ReasonDimensionNotFinite`, `ReasonTagNotSingleValued`, `ReasonRowCanRefuse`.
That append is A14 and belongs to 0006, as C5 already says. A13's site
question is settled independently of it.

## Q2 — fixture sweep for zero-participating-dimension groups

Method: `TestQ2_FixtureSweepForZeroDimensionGroups` walks every `*.toml`
under the checked-in fixture roots, loads each with the real
`table.Load`, partitions with `guard.Groups`, and reports every group where
`guard.Dimensions` is empty. Roots swept: `internal/table/testdata`,
`models`, plus `internal/graphlint/testdata`, `internal/guard/testdata`,
`internal/resolve/testdata` (all three absent — graphlint's fixtures are
authored inline in `fixtures_0006_test.go`, not on disk).

```
$ go test ./docs/rdr/0010-stateless-decision-tables/evidence/spikes/a13-zerodim/ -v
sweep: loaded=37 refused=67 models-with-no-groups=0 zero-dimension-groups=151
```

The 67 refusals are `internal/table/testdata/neg/`, the deliberate
loader-negative fixtures; they never reach lint and are excluded.

**Result: 151 zero-dimension groups across ALL 37 loadable fixtures.** Every
single loadable checked-in model carries at least one. The distinct files:

```
internal/table/testdata/delim/merge-delim-{comma,empty,pipe,semi,space}.toml
internal/table/testdata/delim/write-delim-{comma,empty,pipe,semi,space}-{a,b}.toml
internal/table/testdata/dup/dup-model-{a,b,distinct}.toml
internal/table/testdata/kata-fixture.toml
internal/table/testdata/merge-delim-atom.toml
internal/table/testdata/merge-distinct{,-rev}.toml
internal/table/testdata/merge-idempotent.toml
internal/table/testdata/perm/rdr-{eq-as-in,keyorder,ruleorder}.toml
internal/table/testdata/pos-near-miss-{both,folded,space,upper}.toml
internal/table/testdata/pos-no-{initial,terminal}.toml
internal/table/testdata/pos-not-near-miss.toml
internal/table/testdata/rdr-fixture.toml
internal/table/testdata/write-delim-{a,b}.toml
models/rdr.toml
```

Class check: NONE of these declares `[model].class`, so all 151 groups sit
in the state-machine default class. Typical shapes are terminal rows and
bare escape rows that match on state but carry no `guard.all`/`guard.unless`
atom at all — e.g. `models/rdr.toml`'s `draft-no-match-escape` and
`terminal-archive` groups.

### Verdict on Q2: PASS on the assumption, but the blast radius is total

A13's two halves both hold, and the second is the sharp one:

- **Zero-dimension groups over a state machine are not a corner case — they
  are the norm.** 151 of them, in 37 of 37 loadable fixtures, including the
  production model `models/rdr.toml`. A13's clause "where a zero-dimension
  group is legitimate and stays silent" is confirmed empirically: they are
  silent today (or take only the closure advisory), and they must stay that
  way.
- **Therefore the arm MUST be class-keyed, and an unconditionally-keyed arm
  is catastrophic, not marginally noisy.** Keying it wrongly would newly
  fire `graph-unprovable-coverage` — a BLOCKING code — on 151 groups and
  turn every checked-in fixture and the production model red at once. The
  0006 suite would fail loudly and immediately, so the mis-keying is not a
  silent-drift risk; but the fixture sweep is what converts "the class fence
  matters" from a design preference into a hard requirement with a measured
  cost.

## Overall verdict

A13 holds. Both unverified halves are now grounded:

| Question | Verdict |
|---|---|
| Q1 — emission ahead of `emitCoverageArms` without double-reporting | PASS; requires the early `return`, which the branch shape admits |
| Q2 — fixture sweep for newly-firing zero-dimension groups | PASS; 151 groups in 37/37 fixtures, all state-machine class, all must stay silent |

Residual dependency, unchanged and already booked: A14's
`no-participating-dimension` reason append to
`internal/graphlint/taxonomy.go`, which is 0006's to make. C5 already states
this and A13 does not depend on it for its site claim.

## Commands run

```
/Users/cwensel/sandbox/newcoinc/rdr/bin/rdr inspect --select 0010:A13 0010
/Users/cwensel/sandbox/newcoinc/rdr/bin/rdr inspect --select 0010:C5 0010
go test ./internal/graphlint/... ./internal/table/...
go test ./docs/rdr/0010-stateless-decision-tables/evidence/spikes/a13-zerodim/ -v \
  > docs/rdr/0010-stateless-decision-tables/evidence/spikes/a13-zerodim/a13-spike.out 2>&1
git status --porcelain internal/     # empty: internal/ untouched
```
