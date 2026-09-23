# Triage — RDR 0012 Declared-Kind Carrier at the Guard Seam

Single-pass unattended triage of the roborev findings left by the Stage-8
launch. Window frozen once: `BASE=901ae3d`, `HEAD=6c113a5` (5 commits).
Batch label: `batch:rdr-0012`.

- **Spine** — 2 in-window per-commit auto-reviews carrying findings (jobs
  8229 on `38a7e4b`, 8247 on `6c113a5`), one finding each. Jobs 8236, 8237
  and 8245 reviewed passing. Ancestry confirmed mechanically per job ref.
- **Range net** — one `--since BASE` review (job 8250), 2 findings; both
  duplicate the spine and are charged to jobs 8229 and 8247.
- **Pre-filter** — no finding dropped without a leaf. Neither cited hunk
  changed between its ref and `HEAD`.

Each surviving finding was grounded in a read-only leaf against the RDR +
`{RDR_RESOURCES}`. The leaf confirmed each one with a falsification probe in a
detached scratch worktree, never in the launch worktree.

## Dispositions

| # | Source | Location | Verdict | odc-type | odc-trigger | Outcome |
|---|---|---|---|---|---|---|
| 1 | 8229 + 8250 | `internal/cli/guard_typed_0012_test.go:897`: the invalid-spelling cases accept any lint failure without "canonical spelling", so a loader that trims `" 1"` still passes on `graph-coverage-gap` | IN-SCOPE (`0012:C5` / REQ-45) | test-oracle | test-coverage | fixed:974b546. Fix: in the invalid-spelling loop, require `env.Code == "model-invalid"` plus a `malformed_predicate_atom` finding whose message contains `strconv.Quote(spelled)+" is not an int"`. The mutant trim probe passes at HEAD and fails with the fix. Grounding: `/tmp/batch:rdr-0012/8229.md` |
| 2 | 8247 + 8250 | `internal/guard/seam_carrier_0012_test.go:603`: the go/types zero-value construction check misses implicitly zero-initialized array elements (`[1]guard.Evaluator{}[0]`, `[2]guard.Evaluator{guard.NewEvaluator(nil)}[1]`) | IN-SCOPE (`0012:S4` / REQ-48, backstop of `0012:C4` / REQ-6; D4 "an array of the type counts as holding it") | checking | boundary | fixed:31ba7bb. Fix: in the CompositeLit arm of `r12FileZeroValueSites`, also flag an array literal whose element type holds an Evaluator when `int64(len(x.Elts)) < a.Len()`. Add the `[1]guard.Evaluator{}[0]` plant to `TestAdv1_0012` (the partial-fill `[2]guard.Evaluator{guard.NewEvaluator(nil)}[1]` plant too). Grounding: `/tmp/batch:rdr-0012/8247.md` |

No finding routed OUT-OF-SCOPE, so no kata was filed. Every verdict rests on
positive evidence.

## Close-ledger

The fixup commits' per-commit reviews (jobs 8255 on `974b546`, 8256 on
`31ba7bb`) came back clean ("No issues found"), so the round-2 sweep filed
nothing. No row is `held:`, so no kata was filed. Jobs 8229, 8247 and 8250
were responded to with their `fixed:<sha>` Outcomes and closed. The
completion gate reads `COMPLETE` at `31ba7bb`.
