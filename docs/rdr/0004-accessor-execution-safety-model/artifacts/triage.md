# Triage — RDR 0004 Accessor Execution Safety Model

Single-pass unattended roborev triage after the Stage 8 launch.

- **Window (frozen once)**: `BASE f3fad20` .. `HEAD b9b4916`
- **Batch label**: `batch:rdr-0004`
- **Spine**: 7 in-window per-commit auto-reviews (jobs 6210, 6212, 6213, 6215,
  6216, 6217, 6218). Job 6214 (`5798214`) and 6219 (`f7c7a1e`) passed.
- **Cross-commit net**: job 6220 (`--since f3fad20`), run once.
- **Raw findings**: 22 → **8 clusters** after dedup.
- **Fix-commit sweep**: bounded to one round (jobs 6221, 6222 clean; 6223
  filed, never re-fixed).

## Verdicts

| # | Finding (location) | Jobs | Verdict | odc-type | odc-trigger | Outcome |
|---|---|---|---|---|---|---|
| F1 | `readerFor` selects the first same-role reader — `model.go:231` | 6215, 6220 | RDR-SEED | interface | design-conformance | kata `p63c` ★ |
| F2 | Slice passed to `Apply` reused as the read-back expectation — `executor.go:251/316/352` | 6215, 6220 | FIX-NOW | algorithm | design-conformance | `7dd30bc` |
| F3 | Non-cooperative binding blocks forever — `executor.go:168/200/316` | 6220 | DROP `rdr-adjudicated` | interface | design-conformance | closed |
| F4 | `Lookup`/`bound` ignore `Identity.Flow` — `model.go:206`, `validate.go:95` (+ test arm `capability:136`) | 6210, 6220 | DROP `unreachable` | interface | design-conformance | closed |
| F5a | `Write` not restricted to `def.OwnedKeys()` — `executor.go:251` | 6215 | DROP `superseded-at-HEAD` | n/a-not-a-defect | — | fixed by `4b6d01d` |
| F5b | Pre-write snapshot failure ignored — `executor.go:271` | 6215 | DROP `superseded-at-HEAD` | n/a-not-a-defect | — | fixed by `a914edb` |
| F6 | Planned-owned read-back skipped on incomplete baseline — `executor.go:346` | 6218 | KATA-BUG | checking | rare-situation | kata `bkf3` |
| F7a | ADV-2/ADV-4 accept `read_back_mismatch` where REQ-62 forbids it — `adversarial:138,289` | 6216, 6218 | FIX-NOW | test-oracle | design-conformance | `12289bd` |
| F7b | ADV-2 requires post-mutation refusal — `adversarial:148` | 6216 | DROP `rdr-adjudicated` | test-oracle | design-conformance | closed (D17) |
| F7c | ADV-1 asserts refusal existence, not identity — `adversarial:68` | 6216 | FIX-NOW | test-oracle | mutation-sensitivity | `2c4b5b3` |
| F7d | `nonOwnedPlanKeys`' two conjuncts individually untested — `executor.go:429` | 6217 | FIX-NOW | test-oracle | mutation-sensitivity | `2c4b5b3` |
| F8a | Non-owned read-back unsatisfiable — `fixtures:281`, `validation:205` | 6210 | DROP `superseded-at-HEAD` | n/a-not-a-defect | — | closed (D14) |
| F8b | `NextTags` == `Writes` in every fixture — `fixtures:305` | 6210 | KATA-BUG | test-oracle | logic-flow | kata `742a` |
| F8c | Timeout oracle compares against `blockFor`/1h, not 50ms — `boundary:41,58,72,92`, `desk_trace:121` | 6210, 6213 | KATA-BUG | test-oracle | timing-serialization | kata `742a` |
| F8d | MVV non-owned mismatch control unreachable — `mvv:359` | 6212 | DROP `superseded-at-HEAD` | n/a-not-a-defect | — | closed (D14) |
| F8e | REQ-101 simultaneity never asserts `sawUnreadable` — `desk_trace:53,111` | 6213 | KATA-BUG | test-oracle | logic-flow | kata `742a` |
| F8f | Context-bound test accepts any context/deadline — `desk_trace:232` | 6213 | KATA-BUG | test-oracle | interface | kata `742a` |
| F8g | "missing or empty" tests only a nil slice — `desk_trace:282` | 6213 | DROP `unreachable` | n/a-not-a-defect | — | closed (`TestReq17`) |
| S1 | `Refusal.Keys` asserted by membership, not exact set — `adversarial:78` | 6223 | KATA-BUG (attached) | test-oracle | mutation-sensitivity | kata `742a` (5th site) |

★ = carries an `## Open question`.

## Counts

- **dropped**: 9 — `superseded-at-HEAD`=4, `rdr-adjudicated`=2, `unreachable`=3
- **fixed-now**: 4 findings in 3 commits — `7dd30bc`, `12289bd`, `2c4b5b3`
- **kata-bug**: 2 katas (`bkf3`, `742a` — the latter collapsing 5 sites)
- **rdr-seed**: 1 kata (`p63c`)

## The one real safety defect

**F2** was the finding that mattered. A single `slices.Clone` backed the
`Apply` argument, the `verifyReadBack` oracle, seven `Expected` diagnostics,
and `verifiedWritten` — so a binding that mutated its argument rewrote the
expectation it would be judged against. The red-before run showed the defect
was worse than "verification passes": `Written` actively reported the binding's
substituted value as verified. On a write path mutating authoritative artifacts
with no ledger and no undo, that is a self-referential oracle. Fixed by giving
`Apply` its own clone; regression test ADV-5 confirms red-before/green-after.

## Where the RDR's own adjudications won

Three findings were dropped because the RDR had already decided the question
and roborev — sandboxed to the repo, without the RDR — reversed it:

- **F3**: `0004:C15` requires a bounded timeout *reported as its own refusal
  class*; REQ-71 names the mechanism (the ratified spike's
  `context.WithTimeout` with a cooperative binding). The proposed goroutine
  abandonment would create an in-flight mutation the executor cannot observe —
  precisely what `0004:C14` forbids.
- **F7b**: REQ-63/REQ-64 fix `read_back_incomplete` as post-mutation by
  contract; a pre-write abort would mint it for a write that never occurred.
  Deviation D17 got this right.
- **F5a/F5b**: already fixed in-window by Phase 3c.

## Follow-ups

`/kata-flight --label batch:rdr-0004 --drain` ships the two `type:bug`
children (`bkf3`, `742a`). The `kind:rdr-seed` child (`p63c`) routes to RDR
authoring — `kata-ship` gate 4 refuses it by design.
