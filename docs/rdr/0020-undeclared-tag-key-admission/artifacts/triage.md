# Triage — RDR 0020 Undeclared `--tag` key admission

Single-pass, unattended. Window frozen once:
`BASE d97f2c5 .. HEAD 9bca897` on `worktree-rdr-0020`.
Batch label: `batch:rdr-0020`.

## Collection

| source | result |
|---|---|
| spine (per-commit auto-reviews) | 4 reviews on the branch: `886e600` P, `1b17dcb` P, `455b792` P, `299f27b` **F**. The five `chore(rdr)` artifact-only commits drew no review. |
| range net (`review --since BASE`, job #6959) | 1 finding — the SAME defect as #6953. Deduped against the spine; charged to #6953. No cross-commit-interaction defect emerged. |

**Findings: 1** (spine 1, range-net 1, deduped 1).

Notably the range review's own summary is otherwise clean: "admits
undeclared `--tag` values verbatim while preserving admission guards and
declared-tag validation, with documentation and extensive regression
coverage."

## Findings

### F1 — the carrier-persistence oracle is vacuous

- **Jobs**: #6953 (spine, `299f27b`) + #6959 (range net) — same defect.
- **Severity**: Low (roborev). Kept Low: test-only, no shipped behaviour
  affected, no false green about the PRODUCT — but see grounding.
- **Location**: `internal/cli/flow_carrier_adv_0020_test.go:348`
  (`TestAdv2b_0020_TheAdmittedArrayRunsAccessorsWithoutReachingOwnedState`)
- **Problem**: the test admits a carrier, then asserts the carried key did
  not reach owned state via
  `if strings.Contains(state, "stray")`. `state` comes from
  `flow read-state`, and `flowbind.Reader.Read`
  (`internal/cli/flowbind/flowbind.go:189`) loops over `requested` ONLY:

  ```go
  for _, key := range requested {
      v, held := s[key]
      values = append(values, accessor.KeyValue{Key: key, Value: v, Absent: !held})
  }
  ```

  Neither fixture reader requests `stray`, so `state` can never contain
  `"stray"` whatever is actually persisted. The assertion cannot fail.
  The test also never checks the write set despite claiming to.
- **Verdict**: **FIX-NOW**. Contained, unambiguous, in this RDR's scope
  (the test is this branch's own code), cheap.
- **Grounding**: the assertion is the guard on the record's **Decision
  Rationale (b)** — a carrier is unreadable and unwritable by every
  model-side path. That is a load-bearing claim of 0020, and an oracle
  that cannot fail does not protect it. Verified against source before
  accepting the finding rather than taking the reviewer's word:
  `Reader.Read`'s `requested` filter is exactly as described.
- **odc-type**: `test-oracle`
- **odc-trigger**: `test-coverage`
- **Outcome**: `fix-now:dc4e221` — replaced the vacuous assertion with a
  read-only OBSERVER model that declares and requests `stray` over the same
  artifact role (so the read can actually see the key), asserted absent;
  added a guard that the observer really requested it, a `status=final`
  control read, and exact write-set verification from the `set-state`
  payload.
- **Non-vacuity proof** (required before accepting the fix): the fixer
  temporarily patched `runFlowSetState` to inject every observed carrier
  into the artifact and the reported write set — the world Rationale (b)
  forbids. Under that perturbation BOTH new assertions FAILED, while a
  scratch probe of the OLD `strings.Contains(state, "stray")` assertion
  PASSED — a differential proof that the old oracle was vacuous and the new
  one is not. Perturbation and probe reverted; only the test file changed;
  full suite re-run green.

## Not re-litigated here

The three record-level findings from the launch's own Phase 3 were already
routed before triage and are NOT roborev findings; they are listed so this
file is a complete picture of the run's spin-offs:

| kata | route | source |
|---|---|---|
| `intrastate#r41t` | `type:bug` `lifecycle:queued` | ADV-1: ill-formed UTF-8 collapses to U+FFFD in the `observed` echo |
| `intrastate#fk4n` | `kind:rdr-seed` | D3/ADV-3: silent-typo acceptance rests on an echo the refusal envelope lacks |
| `intrastate#g3ks` | `kind:rdr-seed` | D1 residue: REQ-40/41 names a non-existent `docs/model-schema.md` |


## Fix-commit sweep (bounded, one round)

`dc4e221`'s own auto-review was collected once, per the skill's
bounded-sweep rule. Result recorded below; a further FIX-NOW in a
fix-commit review would be FILED as a kata, never re-fixed.

## Report

```
triage: docs/rdr/0020-undeclared-tag-key-admission.md @ 9bca897  (base d97f2c5)  batch: batch:rdr-0020
  findings: 1  (spine 1  range-net 1  deduped 1)
  dropped:  0
  odc:      test-oracle=1  (trigger: test-coverage)
  fixed-now: 6953=dc4e221   (6959 deduped into 6953)
  kata-bug:  none from roborev
  rdr-seed:  none from roborev
  next: /kata-flight --label batch:rdr-0020 --drain   # ships r41t (the two rdr-seeds are refused by design)
```

Zero DROPs: the single finding was real, confirmed against source, and
fixed. No false positives to charge to the reviewer this run.
