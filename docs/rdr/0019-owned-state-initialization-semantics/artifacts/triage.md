# Triage — RDR 0019 owned-state initialization semantics

Single-pass roborev triage of the findings left after the Stage-8
implementation launch. Window frozen once: `BASE fcb85d6` →
`HEAD adb5bfb` (13 commits). Batch key `batch:rdr-0019`.

Repo-wide there were 79 open roborev jobs; 2 are in this window (the rest
are older `main` history). One cross-commit range review was run over the
window as an additive net.

## Verdicts

| # | Finding | Location | Verdict | odc-type | odc-trigger | Outcome |
|---|---|---|---|---|---|---|
| 1 | Escape predicate rejects the unescaped bytes the clause requires | `flow_initstate_encoder_0019_test.go:396` | DROP superseded-at-HEAD | n/a-not-a-defect | test-coverage | `drop:superseded-at-HEAD` |
| 2 | `command-error` conflates unregistered command with unknown flag | `flow_initstate_fixtures_0019_test.go:562` | DROP superseded-at-HEAD | n/a-not-a-defect | test-coverage | `drop:superseded-at-HEAD` |
| 3 | Seal-only predicate setup leaves the applied value beside the seal | `flow_initstate_predicate_0019_test.go:167` | DROP superseded-at-HEAD | n/a-not-a-defect | test-coverage | `drop:superseded-at-HEAD` |
| 4 | Second refusal asserted over the artifact the first run sealed | `flow_initstate_surface_0019_test.go:563` | DROP superseded-at-HEAD | n/a-not-a-defect | logic-flow | `drop:superseded-at-HEAD` |
| 5 | Mismatch fixture splits writer and reader across two roles | `flow_initstate_noop_0019_test.go:640` | DROP rdr-adjudicated | n/a-not-a-defect | design-conformance | `drop:rdr-adjudicated (DEV-9)` |
| 6 | Flag inheritance probed before cobra merges persistent flags | `flow_initstate_surface_0019_test.go:214` | DROP superseded-at-HEAD | n/a-not-a-defect | test-coverage | `drop:superseded-at-HEAD` |
| 7 | `initAbsentKeys` answers from the union of all bound stores | `flow_initstate.go:623` | DROP superseded-at-HEAD | algorithm | design-conformance | fixed in-window @`70b7424` |
| 8 | Missing role reader preempts carrier validation | `flow_initstate.go:472` / `:480` | KATA-BUG | algorithm | design-conformance | `kata es2r` |

Findings 1–6 are job #7017 (`b38b577`), 7–8 job #7019 (`0bcbc35`).
The range review (#7043) re-raised finding 8 at `HEAD` and is charged to
#7019; it independently confirmed finding 7 is fixed at `HEAD`.

## Why findings 1–6 are not escapes

`b38b577` is the flow's Phase-1 spec-test commit, committed **deliberately
red**: tests are written from the locked spec and must fail before any
implementation exists. Each finding named a real fixture defect in that
red commit, and each was corrected by a later in-window commit. Grounding
checked specifically for tests made green by *weakening* an assertion and
found none — every repair kept or strengthened the positive oracle. They
classify `n/a-not-a-defect` and are excluded from the defect distribution.

Two produced durable value beyond the fix: #2 exposed that
`cobraErrorToCLIError` collapses every cobra/pflag error into
`command-error` (which is what let 19 Phase-1 tests pass tautologically
against a missing verb until a `requireVerbRegistered` oracle was added),
and #5 produced the DEV-9 adjudication below.

## Finding 5 — the S8 recipe (DEV-9)

The reviewer's causal analysis matched the launch adjudication exactly:
`Registry.readerFor` selects by the writer's role, so a two-role fixture
yields `hasReader == false` — the exit-3 INCOMPLETE class the scenario
asks to be distinguished *from*. `0019:S8`'s stated construction is
unconstructible.

The clauses it defends (REQ-63, REQ-67) are correct and implemented; only
the scenario's binding topology is wrong, so this is TEST-FIXTURE, not
SPEC-DEFECT. The test was rebuilt on the route that does work — a
command-backed writer with a file-backed reader on ONE role, tripping
`verifyReadBack`'s protected-non-owned-key baseline loop — and now asserts
exit 2, distinctness from the exit-3 code, the clobbered protected key,
and the non-repairing re-run. Residual (scenario text in a Final record,
never amended in place) tracked as kata `ptn7`.

## Finding 8 — the one open defect

`initCarrierGate` emits an artifact-binding refusal from inside its own
loop, so a role-binding defect on an early writer preempts the carrier
check on a later one. `0019:C1` REQ-33 requires the opposite order; REQ-44
independently fixes the missing-reader *code*, and the implementation
satisfies REQ-44 while violating REQ-33.

Reproduced at `HEAD`: a two-writer model (file-backed writer on a
reader-less role sorting before a needed command-backed writer) lints
clean and returns `flow-artifact-missing`; giving every role its own
reader returns `flow-init-carrier-unsupported`. The only variable is
which writer the loop reaches first.

Diagnostic-quality only — both paths refuse at exit 2 with zero writes
and no accessor invoked, so no seeding, emptiness, or read-back guarantee
is affected. Filed as `es2r` rather than fixed inline because the gate's
ordering is pinned by four REQ-level tests (REQ-33/36/37/44) and warrants
its own red/green cycle. `TestReq33_0019` uses a single-writer model where
both defects sit on one accessor, so the two-writer interleaving was never
covered.

## Filed

| Kata | Kind | Labels | Subject |
|---|---|---|---|
| `es2r` | `type:bug` | `src:roborev`, `batch:rdr-0019`, `area:cli`, `severity:low`, `lifecycle:queued`, `release:post-1.0`, priority 2 | carrier-vs-role gate ordering |
| `ptn7` | `kind:rdr-seed` | `src:roborev`, `batch:rdr-0019` | S8's unconstructible fixture recipe |

No open questions: every verdict was settled on positive evidence.
