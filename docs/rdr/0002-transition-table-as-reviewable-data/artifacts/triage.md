# Triage — RDR 0002 Transition Table As Reviewable Data

Single-pass unattended roborev triage after the Stage-8 launch.

- **Window:** `BASE 51e0c87` → `HEAD fca46da` (frozen once, Phase 0)
- **Batch label:** `batch:rdr-0002`
- **Spine:** 12 in-window per-commit auto-review jobs → 39 findings
- **Range net:** one `--since BASE` review (job 6045) → 6 findings, 5 duplicating
  the spine, 1 new (R4)
- **Launch status:** `status.md` COMPLETE, worktree clean

Twelve further open roborev jobs sit on `main`/`via-*` outside this window and
were left untouched — they are the repo's pre-existing backlog, not this run's.

## Verdict table

| Finding | Location | Verdict | odc-type | odc-trigger | Outcome |
|---|---|---|---|---|---|
| 6023.1 | `testdata/neg/neg-no-model-table.toml` | KATA-BUG | test-oracle | boundary | kata `hjxk` |
| 6023.2 | `testdata/dup/dup-model-b.toml` | KATA-BUG | test-oracle | test-coverage | kata `hjxk` |
| 6025.1 | `format_test.go:189` | KATA-BUG | test-oracle | boundary | kata `fmhr` |
| 6025.2 | `format_test.go:329` | KATA-BUG | test-oracle | test-coverage | kata `g67s` |
| 6025.3 | `format_test.go:520` | DROP `unreachable` | test-oracle | boundary | probed: absent-version already refuses `malformed_model_declaration` |
| 6025.4 | `format_test.go:575` | KATA-BUG | test-oracle | boundary | kata `fmhr` |
| 6025.5 | `format_test.go:667` | KATA-BUG | test-oracle | test-coverage | kata `rb25` |
| 6030.1 | `rules_test.go:34` | KATA-BUG | test-oracle | test-coverage | kata `g67s` |
| 6030.2 | `rules_test.go:488` | KATA-BUG | test-oracle | test-coverage | kata `sje2` |
| 6030.3 | `atoms_test.go:455` | KATA-BUG | test-oracle | boundary | kata `fmhr` |
| 6030.4 | `atoms_test.go:501` | DROP `rdr-adjudicated` | test-oracle | design-conformance | REQ-59 binds the producer's rendering; `renderValue` quotes members per `0002:1594-1600` |
| 6033.1 | `dump_test.go:100` | KATA-BUG | test-oracle | test-coverage | katas `rb25`/`pwff` |
| 6033.2 | `dump_test.go:199` | KATA-BUG | test-oracle | test-coverage | kata `ekpq` |
| 6033.3 | `dump_test.go:359` | KATA-BUG | test-oracle | boundary | kata `sp9e` |
| 6033.4 | `dump_test.go:454` | DROP `rdr-adjudicated` | n/a-not-a-defect | design-conformance | `SourceLocator` carries no line/column; `TestReq92` covers the alleged hole |
| 6033.5 | `dump_test.go:482` | KATA-BUG | test-oracle | test-coverage | kata `hjxk` |
| 6033.6 | `dump_test.go:556` | DROP `rdr-adjudicated` | n/a-not-a-defect | design-conformance | including the locator IS REQ-114's strongest conforming oracle |
| 6033.7 | `normalize_test.go:663` | KATA-BUG | test-oracle | test-coverage | kata `pwff` |
| 6033.8 | `normalize_test.go:516` | KATA-BUG | test-oracle | test-coverage | kata `meab` |
| 6036.1 | `roundtrip_test.go:471` | DROP `superseded-at-HEAD` | test-oracle | test-coverage | deviation D6 added the self-exemption |
| 6036.2 | `mvv_test.go:140` | DROP `superseded-at-HEAD` | test-oracle | test-coverage | deviation D7 fixture repair |
| 6036.3 | `roundtrip_test.go:84,117` | DROP `rdr-adjudicated` | n/a-not-a-defect | design-conformance | as 6033.4/6033.6 |
| 6036.4 | `roundtrip_test.go:659` | KATA-BUG | test-oracle | test-coverage | kata `pwff` |
| 6036.5 | `roundtrip_test.go:245` | KATA-BUG | test-oracle | test-coverage | kata `rjdy` |
| 6037.1 | `atoms_test.go:81` | DROP `rdr-adjudicated` | test-oracle | design-conformance | a distinct `table.Block` asserted equal to the kernel's is the shipped design (D5, §D12) |
| 6037.2 | `format_test.go:680` | KATA-BUG | test-oracle | test-coverage | kata `rb25` |
| 6037.3 | `format_test.go:700` | KATA-BUG | test-oracle | test-coverage | kata `rb25` |
| 6037.4 | `accessors_test.go:539` | **FIX-NOW** | test-oracle | test-coverage | `7f8eba8` |
| 6038.1 | `model.go:213` | superseded → **FIX-NOW** | assignment | boundary | Phase 3c FAIL-1; residual family fixed in `2eb3b01` |
| 6038.2 | `normalize.go:75` | **RDR-SEED** | checking | design-conformance | seed `06pe` |
| 6038.3 | `load.go:193` | KATA-BUG | checking | design-conformance | kata `xgt0` |
| 6038.4 | `dump.go:43` | DROP `rdr-adjudicated` | algorithm | design-conformance | `0002:1332` scopes ordering "over one model's rows"; `DumpAll` has no production callers |
| 6038.5 | `load.go:50,190` | DROP `project-scoped-out` | documentation | doc-code-drift | `Failure.Error` is pre-CLI (REQ-106); A5 defers the user-facing mapping |
| 6039.1 | `roundtrip_test.go:478` | KATA-BUG | test-oracle | test-coverage | kata `v3yh` |
| 6040.1 | `adversarial_0002_test.go:126` | DROP `superseded-at-HEAD` | n/a-not-a-defect | — | Phase 3b's replacement oracle confirmed in the file |
| 6041.1 | `normalize.go:90` | DROP `superseded-at-HEAD` | algorithm | logic-flow | ADV-3 fix sorts `in` members |
| 6042.1 | `normalize.go:623` | **FIX-NOW** | assignment | boundary | `2eb3b01` |
| 6044.1 | `model.go:216` | **FIX-NOW** | assignment | boundary | `2eb3b01` |
| 6044.2 | `fixup_0002_test.go:105` | **RDR-SEED** | interface | design-conformance | seed `06pe` |
| R1 (6045) | `normalize.go:75` | **RDR-SEED** | checking | design-conformance | seed `06pe` |
| R2 (6045) | `model.go:216` | **FIX-NOW** | assignment | boundary | `2eb3b01` |
| R3 (6045) | `load.go:193` | KATA-BUG | checking | design-conformance | kata `xgt0` |
| R4 (6045) | `normalize.go:298` | KATA-BUG | checking | boundary | kata `q7jc` |
| R5 (6045) | `dump.go:43` | DROP `rdr-adjudicated` | algorithm | design-conformance | as 6038.4 |
| R6 (6045) | `load.go:50,190` | DROP `project-scoped-out` | documentation | doc-code-drift | as 6038.5 |

## FIX-NOW commits

- **`2eb3b01`** — guard literal carriage keys on the **operator alone**.
  Deviation D8 (Phase 3c) added the operator trigger with `||` but never removed
  the declared-kind one, so a legal `exists` guard on a `set`-kind tag
  array-encoded as `["true"]`. The kernel compares verbatim (`0007:C3`), matched
  neither `LiteralTrue` nor `LiteralFalse`, and returned `GuardUnevaluable` —
  which under RDR 0007's unevaluable-candidate veto poisons decidable sibling
  rows. The red run reproduced it exactly: `Literal:["true"] Reason:uncomparable`.
  Also fixes the `eq`/comparison family on set-kind tags. Tag **value** carriage
  stays kind-keyed (D8's "two rules answer different questions"). Recorded as
  deviation **D11**.
- **`7f8eba8`** — REQ-33's terminal-atom oracle now `reflect.DeepEqual`s the
  whole `Atom` instead of only `.Key`.

Both fix commits were auto-reviewed (jobs 6046, 6047): **no issues found**. The
fix-then-review sweep is bounded to that one round.

## Filed kata

Thirteen `type:bug` (`lifecycle:queued`) and one `kind:rdr-seed`, all carrying
`src:roborev` + `batch:rdr-0002` + `area:internal-table` + `release:v1.0`.

| short_id | Title | Sev / Pri |
|---|---|---|
| `rb25` | TestReq18 asserts difference and loading, never member-sequence survival | high / 1 |
| `g67s` | Oracles accept any refusal but a named few | medium / 2 |
| `fmhr` | Three controls author their mutation in the wrong scope | medium / 2 |
| `sje2` | TestReq44 shallow-copies rows and searches the whole dump | medium / 2 |
| `ekpq` | REQ-94 witnesses searched over the whole dump, not per column | medium / 2 |
| `pwff` | Handoff checks trace Match entries but never guard entries | medium / 2 |
| `rjdy` | Promoted-fixture check compares basenames only | medium / 2 |
| `xgt0` | Tag declaration accepts `single_valued` on `set` and `scalar` | medium / 2 |
| `q7jc` | A rule may omit its required local match block and still normalize | medium / 2 |
| `hjxk` | `dup-model` pair and `neg-no-model-table` cannot discriminate | low / 3 |
| `sp9e` | REQ-98 empty-suffix arm is decided by rule id, never by suffix | low / 3 |
| `meab` | No-alias control cannot observe a shared `TagValue.Value` array | low / 2 |
| `v3yh` | REQ-124 hash ban exempts the whole of `roundtrip_test.go` | medium / 2 |
| `06pe` | **seed** — operator/kind matrix fenced in RDR 0003, enforced by nobody | medium / 1 |

`rb25` carries the highest priority of the bugs because its blind spot is
precisely the `members[0]` truncation (FAIL-1 / D8) that shipped in this build
and had to be fixed in Phase 3c — the oracle was green throughout.

Seed `06pe` is the run's most consequential result. Probing found nine distinct
operator/kind violations that load clean at HEAD, including an empty match-block
`in` that silently collapses a rule to **zero rows**. RDR 0002's `C22` forbids
restating RDR 0003's model, RDR 0003 does not load documents, and RDR 0006's
lint does not exist — so the matrix has no enforcement site at all. It carries an
`## Open question` on whether the empty-literal arm severs from the matrix
question, and a coupling note that `fixup_0002_test.go:105` must be migrated if
the matrix is ever enforced at load.

## Notes

- Every one of the 45 findings reached exactly one terminal state; all 14
  in-window roborev jobs (12 spine + range net + 2 fix-commit) were commented
  with the citing reason and closed.
- Six drops overturned the finding on record evidence rather than on staleness.
  The `SourceLocator` cluster (6033.4/6033.6/6036.3) is the clearest: all three
  rest on the premise that the locator carries optional line/column detail, but
  it is built as `model.ID + ":" + id`, so including it in `DeepEqual` is the
  *strongest* conforming oracle — excluding it would have been the defect.
- The intentionally-RED pre-filter did **not** apply to the weak-oracle batch:
  `git log` shows no in-window commit touched `format_test.go`, `rules_test.go`,
  `atoms_test.go` or `accessors_test.go` after the red gate. Phase 3b/3c added
  new tests without strengthening old ones, so those oracles are live at HEAD.
