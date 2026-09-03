Model: claude-opus-5

# A5 — in-group backgrounded grandchild clears the 500 ms drain bound

Verification of RDR 0026 Critical Assumption A5 by MVV Test: run the existing
grandchild-leak oracle repeatedly under `-race` and CPU load, and separately
measure the quantity A5 actually bounds — the wall time from the group
`SIGKILL` to the parent's drain reaching EOF.

Host: darwin/arm64, 12 cores, go1.26.6. Date: 2026-09-03.

## The bound

`WaitDelay = 500` (milliseconds) — `internal/cli/cmdbind/cmdbind.go:66`,
documented at `cmdbind.go:63-65`. The RDR's claim of 500 ms is CONFIRMED; it is
applied as `cmd.WaitDelay = WaitDelay * time.Millisecond` at `cmdbind.go:224`.

## Oracles identified

| Role | Test | Location |
| --- | --- | --- |
| Grandchild-leak (in-group backgrounded helper, intrastate#x5vq) | `TestAdvBackgroundGrandchildDoesNotSurviveASuccessfulRead` | `internal/cli/cmdbind/adversarial_0025_test.go:313` |
| Grandchild-leak, drain-stall half (same shape, same tracker) | `TestAdvBackgroundGrandchildMustNotStallASuccessfulRead` | `internal/cli/cmdbind/adversarial_0025_test.go:258` |
| Cancel (intrastate#878e) | `TestAdvKilledChildOnParentCancelIsNotAnAnswer` | `internal/cli/cmdbind/adversarial_0025_test.go:131` |

Both grandchild oracles drive the `( sleep 20 ) &` wrapper shape — the helper
stays IN the child's process group and inherits the stdout/stderr pipes. No
build tag, env var, or flag gates them; `AllowCommands: true` is set in the
test's own `cmdbind.Config`. Nothing was skipped.

## Why the oracle alone does not measure A5

`spawn` joins the two drain goroutines with an UNBOUNDED `drains.Wait()`
(`cmdbind.go:298`), placed after `reapGroup(cmd.Process)` (`cmdbind.go:294`).
The oracle therefore proves the drain terminates, not how much slack sits under
500 ms. The spike measures the interval directly.

## Commands run

```
# oracles, 25 iterations each, race detector, under the 16-proc load
go test -race -count=25 -v -timeout 900s \
  -run 'TestAdvBackgroundGrandchildMustNotStallASuccessfulRead|TestAdvBackgroundGrandchildDoesNotSurviveASuccessfulRead|TestAdvKilledChildOnParentCancelIsNotAnAnswer' \
  ./internal/cli/cmdbind/

# release-latency spike (source: a5_ingroup_margin_spike.go.txt)
cd /tmp/a5spike && go build -o a5spike . 
./a5spike 40   # idle baseline
./a5spike 60   # under 16-proc load
./a5spike 60   # under 32-proc load

# load generators (started/stopped around the measured runs)
for i in $(seq 1 16); do nohup yes > /dev/null 2>&1 & echo $! >> load.pids; done
kill $(cat load.pids); pkill -x yes
```

## Load conditions

| Round | Load generators | 1-min load average (12 cores) |
| --- | --- | --- |
| Baseline | none | idle |
| Loaded | 16 x `yes > /dev/null` | 16.24 |
| Heavier | 32 x `yes > /dev/null` | 47.45 |

The oracle `-race -count=25` round ran concurrently with the 16-proc load.

## Release latency: kill(-pgid, SIGKILL) -> drain EOF

All figures in milliseconds. Bound = 500 ms.

| Round | N | min | median | p90 | p99 | max | mean | over bound | margin (bound/max) |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Idle baseline | 40 | 0.047 | 0.367 | 0.428 | 0.463 | 1.791 | 0.403 | 0/40 | 279x |
| 16-proc load | 60 | 0.186 | 1.187 | 8.127 | 19.874 | 27.685 | 3.414 | 0/60 | 18.1x |
| 32-proc load | 60 | 0.112 | 0.893 | 6.802 | 15.309 | 24.527 | 2.808 | 0/60 | 20.4x |

Aggregate over 160 measured trials: zero at or above the bound. The worst
single observation anywhere was 27.685 ms (16-proc load, trial 36) — 5.5% of
the bound. Every trial returned the correct envelope
(`{"state.phase":"review"}`); no trial errored.

## Verdict

**A5 HOLDS, with comfortable margin.** For the non-escaping in-group
backgrounded grandchild, release after the group SIGKILL is a sub-millisecond
to low-tens-of-milliseconds event dominated by kernel teardown, not by
scheduling. Saturating the machine at ~4x core count moved the median from
0.37 ms to ~1 ms and the max from 1.8 ms to 27.7 ms — an order of magnitude,
but still 18x below the bound. The ordinary "wrapper script backgrounds a
helper" success case pays nothing under a 500 ms bounded join and is never
refused.

The scaling under load is the honest caveat: the tail grew ~15x when the
machine went from idle to saturated, while the median grew only ~3x. The tail
is therefore load-sensitive, and a CI runner far more oversubscribed than
4x-core, or one with a much slower process teardown path, would erode margin.
At the observed scaling rate the bound still absorbs roughly another 18x of tail
growth, so the assumption is not marginal — but the margin is a function of
runner load, not a constant.

Nothing under `internal/` was modified. All load generators were stopped and no
grandchild (`sleep 20`) processes remained after the runs.
