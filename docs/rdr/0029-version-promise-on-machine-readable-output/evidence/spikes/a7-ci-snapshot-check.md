Model: claude-sonnet-5

# A7 spike: committed snapshot of tiered-vocabulary seams + severity, diffed in CI

## Ground truth read

`internal/graphlint/taxonomy.go::AdvisoryCodes` / `::BlockingCodes` / `::IsBlocking`
/ `::severityFor` confirm the two-tier severity partition claimed by A7's
premise: `blockingCodes` (10 members) and `advisoryCodes` (4 members, closed
per `0006:C17`), with `IsBlocking(code)` doing `slices.Contains(blockingCodes,
code)` and `severityFor` mapping onto `SeverityBlocking`/`SeverityInfo`.

`.github/workflows/ci.yml` already hosts a structurally identical
regenerate-then-diff gate: the `docs` job runs `make docs-check`, which
(per `Makefile::docs-check`) builds the binary, renders `docs
--dir "$tmp"`, and diffs the committed `docs/cli-reference.md` / `llms.txt`
against the fresh render — "do the committed files match what this binary
emits?" This is strong precedent for the golden-file shape, though it is
CLI-binary-driven, not a `go test` golden file.

No `go:generate` directives exist anywhere in the repo today
(`grep -rl go:generate **/*.go` → empty). A hidden-verb or go:generate shape
would be the first of its kind; a `go test` golden file has no such
green-field cost — testdata/golden-file patterns already exist
(`internal/table/testdata/`).

## Seam inventory (today-existing only, per task scope)

- `internal/resolve::RefusalKinds` — internal/resolve/resolve.go:67
- `internal/table::Categories` — internal/table/category.go:90
- `internal/accessor::Verdicts` — internal/accessor/model.go:309
- `internal/graphlint::Reasons` — internal/graphlint/taxonomy.go:112
- `internal/graphlint::AdvisoryCodes` / `::BlockingCodes` / `::IsBlocking` —
  internal/graphlint/taxonomy.go:106,109,116-121

Future seams this check would extend to (NOT built yet, per A4/C4's
per-tiered-vocabulary census obligation): any new enumeration seam C4
requires when a new machine-readable surface is added. The spike's
snapshot renderer would gain one more `## package.Func` block per seam;
no structural change.

## Import-graph check (the second unresolved question)

`go list -deps ./internal/graphlint/...` today already includes
`internal/resolve`, `internal/table`, `internal/guard` in PRODUCTION code
(non-test): `internal/graphlint/reach.go` imports `internal/resolve` and
`internal/table` directly; `analysis.go`, `coverage.go`, `engine.go`,
`groups.go` import `internal/table`. `internal/accessor` depends on
`internal/resolve`/`internal/table` but NOT on `internal/graphlint`
(`grep -rn internal/graphlint internal/accessor/*.go` → no matches), so
`internal/graphlint` (or an external `graphlint_test` package beside it)
importing `internal/accessor` from a test creates no cycle.

The originally-intended host, `internal/cli` (already importing all of
resolve/table/accessor/graphlint in production code per `docs.go`,
`lint.go`, `flow_exec.go` etc.), was ALSO checked and is architecturally
fine for this purpose — `internal/cli/lint_0006_test.go`'s import-ban check
(`TestReq95And96_FindingLivesInClierrWithStringTypedAtomFields`) only
forbids `internal/graphlint`/`internal/table`/`internal/resolve` inside
`internal/cli/clierr/clierr.go` specifically (the wire-format package, kept
subsystem-agnostic), not in `internal/cli` generally. However, at spike
time `internal/cli` had a pre-existing, unrelated build break from
concurrent work in the working tree (`internal/cli/respond/respond.go:141:
undefined: clierr.SchemaVersion`) — NOT caused by this spike. The spike was
rehosted in `internal/graphlint` (external test package `graphlint_test`)
to avoid depending on that unrelated broken state, and this confirms the
severity partition is readable, from a test, WITHOUT an awkward import: no
seam package needs to import `graphlint`, and `graphlint_test` importing
`resolve`/`table`/`accessor` is a normal downward-only test import with no
cycle.

## Spike commands, in order

```
cd /Users/cwensel/sandbox/newcoinc/intrastate
git status --porcelain            # captured pre-existing dirty state (NOT mine, see below)
```

Pre-existing (baseline, NOT touched by this spike):
```
 M internal/cli/clierr/clierr.go
 M internal/cli/flow_next.go
 M internal/cli/respond/respond.go
?? docs/rdr/0029-version-promise-on-machine-readable-output/evidence/reconcile/
```
(This baseline visibly shifted mid-spike — e.g. `docs/rdr/0029-...-output.md`
and additional `evidence/spikes/*.md` files appeared/changed between spike
steps — confirming a concurrent session editing the same working tree in
parallel. None of that drift originated from this spike; only
`internal/graphlint/taxonomy.go` and the throwaway test/golden files
described below were touched here, and all three are reverted below.)

Spike file written: `internal/graphlint/a7_snapshot_spike_test.go`
(package `graphlint_test`) — renders a stable text snapshot from
`resolve.RefusalKinds()`, `table.Categories()`, `accessor.Verdicts()`,
`graphlint.Reasons()`, and the `graphlint.AdvisoryCodes()`/
`BlockingCodes()`/`IsBlocking()` severity partition (sorted, tab-separated
code→severity), then compares against a committed golden at
`internal/graphlint/testdata/a7_snapshot_spike.golden`.

```
go test ./internal/graphlint/ -run TestA7SnapshotSpike -v
```
GREEN baseline:
```
=== RUN   TestA7SnapshotSpike
--- PASS: TestA7SnapshotSpike (0.00s)
PASS
ok  	github.com/cwensel/intrastate/internal/graphlint	0.295s
```

### Failing case: undisclosed promotion (advisory → blocking)

Edited `internal/graphlint/taxonomy.go`: moved `CodeVacuousAtom` from
`advisoryCodes` to `blockingCodes` (an undisclosed promotion, exactly the
event C3 requires release-note disclosure for and A7 claims CI would
catch without a reviewer).

```
go test ./internal/graphlint/ -run TestA7SnapshotSpike -v
```

VERBATIM RED OUTPUT (elided to the load-bearing lines; full diff was
symmetric top-to-bottom with one differing line):

```
=== RUN   TestA7SnapshotSpike
    a7_snapshot_spike_test.go:74: A7 snapshot drifted from committed golden.
        --- got ---
        ...
        ## internal/graphlint severity partition
        graph-always-present-owned	blocking
        graph-coverage-closed-by-escape	info
        graph-coverage-gap	blocking
        graph-dangling-edge	blocking
        graph-dead-end	blocking
        graph-overlap	blocking
        graph-owned-before-write	blocking
        graph-product-too-large	blocking
        graph-redundant-row	info
        graph-single-valued-state	blocking
        graph-terminal-escape	blocking
        graph-unprovable-coverage	blocking
        graph-unreachable-rule	info
        graph-vacuous-atom	blocking
        --- want ---
        ...
        graph-vacuous-atom	info
--- FAIL: TestA7SnapshotSpike (0.00s)
FAIL
FAIL	github.com/cwensel/intrastate/internal/graphlint	0.295s
FAIL
```

The diff isolates exactly the moved code (`graph-vacuous-atom`:
`info`→`blocking`), nothing else — a reviewable, readable failure a CI log
would show directly. This demonstrates A7's claim: the undisclosed
promotion fails the build.

### Restore and confirm green

```
git checkout -- internal/graphlint/taxonomy.go
go test ./internal/graphlint/ -run TestA7SnapshotSpike -v
```
```
=== RUN   TestA7SnapshotSpike
--- PASS: TestA7SnapshotSpike (0.00s)
PASS
ok  	github.com/cwensel/intrastate/internal/graphlint	0.368s
```

### Cleanup

```
rm -f internal/graphlint/a7_snapshot_spike_test.go
rm -rf internal/graphlint/testdata
git status --porcelain
```
Final status (my changes fully reverted; only the pre-existing concurrent
session's own files remain, none touched by this spike):
```
 M docs/rdr/0029-version-promise-on-machine-readable-output.md
?? docs/rdr/0029-version-promise-on-machine-readable-output/evidence/reconcile/
?? docs/rdr/0029-version-promise-on-machine-readable-output/evidence/spikes/a5-enumeration-seams.md
?? docs/rdr/0029-version-promise-on-machine-readable-output/evidence/spikes/a6-schema-version-home.md
```

## Answers to the two unresolved questions

**(a) Generation shape.** Recommend a **`go test` golden-file test**, not a
hidden verb and not `go:generate`. Reasons:
- The repo already has zero `go:generate` directives — introducing one
  here would be a new mechanism for one check, where `docs-check`'s
  binary-driven regenerate-diff pattern exists but is heavier (requires
  `make build` first) and is scoped to CLI-emitted docs, not internal Go
  package seams.
- A `go test` golden file needs no binary build, runs under the existing
  `test` CI job (`.github/workflows/ci.yml::test` → `make test-ci`) with
  zero new CI job, and a failing diff is immediately visible in `go test -v`
  output exactly as captured above — no separate CI step or script to
  maintain.
- "Hidden verb" (an internal CLI subcommand no one runs) adds a second
  code path with its own flag/help surface for a check that is purely a
  build-time invariant, not an operator action — the docs example
  (`intrastate docs --dir`) is precedent for hidden verbs but that verb
  emits real deliverables (docs consumers read); this check has no such
  deliverable.
- Regenerating the golden on legitimate change is a one-line addition
  (a `-update` flag or a tiny generator test, as spiked ad hoc above),
  matching the existing `internal/table/testdata` golden convention.

**(b) Import readability.** YES, readable without an awkward import.
Hosted in `internal/graphlint` as an external test package
(`package graphlint_test`), importing `internal/resolve`,
`internal/table`, `internal/accessor`, and `internal/graphlint` itself.
This is architecturally identical to imports `internal/graphlint`'s own
production code already makes (`resolve`, `table`) plus one new
test-only import (`accessor`), which introduces no cycle because
`accessor` does not depend on `graphlint`. The alternative host,
`internal/cli`, is also import-clean for this purpose (its own
production code already imports all four seam packages; the only import
ban in that package is scoped narrowly to `internal/cli/clierr/clierr.go`)
but was not used for the demonstration because it had an unrelated,
pre-existing build break from concurrent work in the shared working tree
at spike time.

## Verdict on A7

A7's claim holds: a committed snapshot, diffed in CI, catches an
undisclosed severity promotion and fails the build with a readable diff,
using only seams that exist today. No fallback to review-only disclosure
is needed. The check should ship as a `go test` golden-file test hosted
in `internal/graphlint` (or `internal/cli` once its current build is
green), wired into the existing `test` CI job — no new CI job required.
