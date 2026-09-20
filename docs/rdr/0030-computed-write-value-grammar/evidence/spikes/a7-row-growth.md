Model: claude-opus-5

# A7 spike — step-expansion row growth needs no load-time ceiling

## Assumption under test

A7: "The step expansion's row growth is bounded by the declaration alone —
at most the admitted-cell count, itself at most the domain width — and
needs no load-time ceiling at the ladder sizes reported; the lint
enumeration limit is RDR 0013's caller-supplied analysis scope (`0013:C1`,
`0013:C3`) and gates nothing at load."

If wrong: `0030:C1` gains a load-time ceiling of its own.

## Method and its limit

The proposed `{ step = n }` write form is NOT implemented, so no step model
can be loaded. The spike measures the PROXY that decides the question:
what the existing loader costs on a table of N **literal** rows — exactly
what C1 says the expansion mints ("The loader expands a rule carrying a
step write into literal rows, one per ADMITTED CELL of the stepped tag").
Post-normalization C1 requires the shapes be indistinguishable ("After
normalization no surface distinguishes an expanded row from an authored
literal row"), so literal-row load cost is the right stand-in for
expanded-row load cost. What the proxy does NOT measure is the expansion
arithmetic itself, which is O(admitted cells) by construction.

Verb selection matters and cost a false start:

- `flow next` — load + normalize + candidate selection. No invariants, no
  product enumeration, no reachability traversal. This is the load probe.
- `graph` — load + normalize + reachability traversal. NOT a pure load
  probe: its traversal is superlinear (15.32s at 10k rows, emitting a
  1.4 GB document), so it is reported separately and not used for (a).
- `lint` — the full graph invariants, including the 0013 enumeration.

Binary built at commit 586ecba via `make build`.

## Generator

Saved at `/tmp/a7-spike/gen.py`; fixtures at `/tmp/a7-spike/m<M>.toml`
(`m99`, `m999`, `m9999`, `m100000`, `m500000`) and `/tmp/a7-spike/s<M>.toml`
for the lint curve. State artifact `/tmp/a7-spike/state.json` holds
`{"counter":"0","phase":"running"}`.

```python
#!/usr/bin/env python3
"""Generate a synthetic intrastate model whose rule table is exactly the
literal-row shape RDR 0030's `{ step = n }` write form would expand to:
one rule per admitted cell of an int domain `min = 0, max = M`, each
guarding `counter eq = i` and writing `counter = i+1`."""
import sys

M = int(sys.argv[1])          # max of the int domain; rows = M+1
out = sys.argv[2]

L = []
w = L.append
w('# A7 spike fixture: step expansion rendered as literal rows.')
w('outcomes = ["advance"]')
w('terminal = ["done"]')
w('')
w('[model]')
w(f'id = "a7-step-{M}"')
w('version = 1')
w(f'description = "Step expansion over an int domain 0..{M} as literal rows."')
w('')
w('[initial]')
w('counter = 0')
w('phase = "running"')
w('')
w('[tags.counter]')
w('provenance = "owned"')
w('kind = "int"')
w('min = 0')
w(f'max = {M}')
w('single_valued = true')
w('required = true')
w('')
w('[tags.phase]')
w('provenance = "owned"')
w('kind = "enum"')
w('domain = ["running", "done"]')
w('single_valued = true')
w('required = true')
w('')
w('[tags.recognized]')
w('provenance = "recognized"')
w('kind = "enum"')
w('single_valued = true')
w('required = true')
w('')
w('[read.state]')
w('role = "state"')
w('path = "state.counter"')
w('keys = ["counter", "phase"]')
w('timeout = "2s"')
w('')
w('[write.state]')
w('role = "state"')
w('path = "state.counter"')
w('keys = ["counter", "phase"]')
w('timeout = "2s"')
w('read_back = true')
w('')
w('[context.running.match.phase]')
w('eq = "running"')
w('')
w('[context.done.match.phase]')
w('eq = "done"')
w('')
for i in range(M + 1):
    w('[[rule]]')
    w(f'id = "step-{i}"')
    w('use = ["running"]')
    w(f'source = "a7:step-{i}"')
    w('[rule.match.recognized]')
    w('eq = "advance"')
    w('[rule.guard.all.counter]')
    w(f'eq = {i}')
    w('[rule.write]')
    w(f'counter = {min(i + 1, M)}')
    w('')

open(out, 'w').write('\n'.join(L) + '\n')
print(f'{out}: rows={M+1} bytes={sum(len(x)+1 for x in L)}')
```

Fixture validity was confirmed at M=9 before scaling: `graph --emit json`
returns a well-formed `intrastate.graph/1` document with all 10 rows, and
`lint` returns `{"type":"ok","schema_version":"0.1","data":{"findings":[]}}`
at exit 0. The generated shape is accepted, not merely parsed.

## Commands

```sh
cd /Users/cwensel/sandbox/newcoinc/intrastate && make build

for M in 99 999 9999 100000; do python3 /tmp/a7-spike/gen.py $M /tmp/a7-spike/m$M.toml; done
python3 /tmp/a7-spike/gen.py 499999 /tmp/a7-spike/m500000.toml

# load + normalize
/usr/bin/time -l ./bin/intrastate flow next \
    --model /tmp/a7-spike/m$M.toml --artifact state=/tmp/a7-spike/state.json

# full invariants (0013 enumeration)
/usr/bin/time -p ./bin/intrastate lint --model /tmp/a7-spike/s$M.toml

# load + normalize + reachability traversal (not a load probe)
/usr/bin/time -p ./bin/intrastate graph --model /tmp/a7-spike/m$M.toml --emit json
```

## Raw output — load/normalize (`flow next`, `/usr/bin/time -l`)

```
-- rows=100 --
        0.00 real         0.00 user         0.00 sys
             8257536  maximum resident set size
-- rows=1000 --
        0.01 real         0.01 user         0.00 sys
            15384576  maximum resident set size
-- rows=10000 --
        0.09 real         0.12 user         0.01 sys
            69468160  maximum resident set size
-- rows=100001 --
        0.96 real         1.27 user         0.09 sys
           584646656  maximum resident set size
-- rows=500000 --
        5.08 real         6.77 user         0.53 sys
          2836267008  maximum resident set size
```

Every run exited 0 and emitted the same 14-line candidate report. No
refusal, no truncation, no limit at any scale.

## Raw output — lint (`/usr/bin/time -p`)

```
rows=10 exit=0 real=0.00s out_lines=0
rows=20 exit=0 real=0.03s out_lines=0
rows=40 exit=0 real=0.29s out_lines=0
rows=60 exit=0 real=1.17s out_lines=0
rows=80 exit=0 real=3.33s out_lines=0
rows=100 exit=0 real=7.53s out_lines=0
rows=150 exit=0 real=36.02s out_lines=0
rows=200 exit=0 real=116.28s out_lines=0
```

A separate `/usr/bin/time -l` run at rows=100 recorded `6.30 real` and
14729216 maximum resident set size — lint's memory stays small while its
TIME explodes. lint at rows=1000 was still running when cancelled past
600s.

## Raw output — graph export (load + normalize + traversal)

```
rows=100   exit=0 graph_real=0.00s  bytes=145764
rows=1000  exit=0 graph_real=0.13s  bytes=12261768
rows=10000 exit=0 graph_real=15.32s bytes=1472729772
```

## Table

| rows | load+normalize (`flow next`) | RSS | lint | graph export |
|------|------|------|------|------|
| 10 | — | — | 0.00s | — |
| 20 | — | — | 0.03s | — |
| 40 | — | — | 0.29s | — |
| 60 | — | — | 1.17s | — |
| 80 | — | — | 3.33s | — |
| 100 | 0.00s | 8 MB | 7.53s | 0.00s |
| 150 | — | — | 36.02s | — |
| 200 | — | — | 116.28s | — |
| 1,000 | 0.01s | 15 MB | >600s (cancelled) | 0.13s |
| 10,000 | 0.09s | 69 MB | not attempted | 15.32s |
| 100,001 | **0.96s** | **584 MB** | not attempted | not attempted |
| 500,000 | 5.08s | 2.8 GB | not attempted | not attempted |

Load is linear in rows (~10 us/row, ~5.8 KB RSS/row). Lint is roughly
quartic: the 100→200 row doubling costs 15.4x.

## Findings

### (a) Does load/normalize of ~100k literal rows succeed?

Yes, and comfortably. 100,001 rows load and normalize in **0.96s wall /
584 MB RSS**, exit 0, correct candidate output. That is squarely in the
"order of seconds" acceptable band, not pathological. Pushed 5x past the
A7 ladder, 500,000 rows still load in 5.08s — linear, no cliff. The
17 MB TOML source at the 100k scale is itself the larger practical
concern, and it is an authoring artifact the expansion never writes to
disk: the step form mints these rows in memory from a two-line
declaration.

### (b) Does any load-time ceiling or row-count limit exist on the load path?

No. Confirmed two ways.

**Structurally.** `AssignmentCount` lives in `internal/guard/declaration.go:80`.
The loader package does not depend on `internal/guard` at all:

```
$ go list -deps ./internal/table | grep -E 'intrastate/internal/(guard|graphlint)'
  (no guard/graphlint dep)
```

Same for `./internal/resolve` and `./internal/accessor`. `AssignmentCount`
is not merely unreached from load — it is **unreachable**, enforced by the
package dependency graph.

**Call path.** Every non-test caller of `AssignmentCount` is in lint or
documentation-export code:

- `internal/graphlint/reach.go:437` — reachability traversal
- `internal/graphlint/coverage.go:241` — coverage invariant
- `internal/guard/lint.go:373` — lint scope check
- `internal/guard/product.go:270,307,374` — `unprovableDimension`,
  `productOver`, `Cardinality`; the product enumeration
- `internal/cli/graph_document.go:188` — decides whether the exported
  `domain` array is present; documentation only, emits no finding

Callers of those: `graphlint.Run` ← `internal/cli/lint.go`, and the graph
document builder ← `internal/cli/graph.go`. Neither is on the
load/normalize path.

The one CLI file on the runtime path that imports guard at all,
`internal/cli/flow_resolve.go`, uses exactly one symbol —
`guard.Evaluator` (line 561). `Evaluator.Evaluate`
(`internal/guard/grammar.go:104`) is a stateless per-atom switch over
`eq`/`in`/`lt`/`lte`/`gt`/`gte`/`contains`; it holds no model, enumerates
no domain, and never calls `AssignmentCount`. The type is documented as
carrying no state precisely so "never reads the tag view" is a property of
the type rather than a discipline.

A grep of `internal/table/` for `bound|ceiling|limit|max_rows|MaxRows`
returns only unrelated uses: comparison-bound parsing in
`normalize.go:168`, edit-anchor value binding in `edit.go`, and int
min/max declaration parsing in `load.go:865`. There is no row-count gate.

The two published bounds are both lint-side implementation constants, as
`lint --help` states verbatim ("Bounds enforced by this build
(model-independent implementation constants, not per-model inputs): product
bound 2048, node ceiling 4096") — `bound = 2048` at
`internal/guard/product.go:420` and `nodeCeiling = 4096` at
`internal/graphlint/taxonomy.go:158`.

### (c) Does lint hit the 0013 enumeration limit, and what does it emit?

Not on this fixture shape — and the way it fails to is itself the evidence.
Every lint run that completed returned **exit 0 with an empty findings
list** (`{"type":"ok","schema_version":"0.1","data":{"findings":[]}}`).
Lint does not refuse these models; it simply gets slower, quartically,
until it is unusable past a few hundred rows.

That is the cleanest possible demonstration of A7's claim. The 0013 limit
is an **analysis scope**: it bounds how much of the product the checker
enumerates before it gives up proving, and when exceeded it degrades the
verdict (an unprovable-coverage finding, `graphlint.CodeUnprovableCoverage`)
rather than rejecting the model. The node ceiling behaves the same way —
`reach.go:131` breaks the traversal and sets `complete = false`, and
`checkNodeCeiling` (`analysis.go:177`) reports it as a finding. Neither is
a gate; neither can stop a load. A model that lint cannot finish analysing
still loads in under a second and still runs.

The practical consequence for 0030 is real but is a LINT consequence: a
step write over a wide domain produces a model that loads instantly and
that `lint` cannot check in reasonable time. That belongs to 0013's
caller-supplied scope, exactly where A7 puts it.

### (d) Judgment

**A7 holds.** The headline: 100,001 rows load and normalize in 0.96s /
584 MB, with no ceiling anywhere on the load path — `AssignmentCount` is
unreachable from load by package dependency, not by convention. C1 does
NOT need a load-time ceiling of its own at the ladder sizes reported.

The expansion's own bound is what C1 already states — one row per admitted
cell, at most the domain width — and the loader absorbs that count
linearly. Load stays acceptable well past the ladder: 500,000 rows in
5.08s. If a ceiling were ever wanted it would be justified by memory, not
time, and the first scale where that argument bites is around 1,000,000
rows (extrapolating 5.8 KB/row to ~6 GB RSS) — an order of magnitude
beyond anything A7 contemplates.

One thing the spike surfaces that A7 does not claim and 0030 should not
mistake for a load problem: `lint` on these shapes is quartic and becomes
unusable around 200 rows, long before load notices anything. Any concern
about wide step domains is a concern about analysability under 0013's
scope, not about load. Answering it with a load-time ceiling would put the
gate in the wrong place — it would refuse models the loader handles fine,
for a cost the loader does not pay.
