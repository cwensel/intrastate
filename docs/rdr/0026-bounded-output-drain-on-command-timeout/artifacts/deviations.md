# Deviations — RDR 0026 bounded-output-drain-on-command-timeout

Entries pre-seeded by the 7.1 cluster reconcile (0026-0027-0028, 2026-09-03)
are left OPEN for Stage 8: running the named check is the disposition; an
entry escalates only if the check contradicts a contract.

## D1 — Negative witness for the `Applied()`-keyed rendering

**Type**: TEST-FIXTURE. **Status**: OPEN (pre-seeded, 7.1 pairwise 0026×0028 PW3; critique C-9).

**Situation.** `0026:S1` asserts only the positive: `Applied()` true on both
write legs renders `detailMayHaveApplied`. No scenario asserts the negative —
a WRITE-phase `execution_failure` with `Applied()` false (every refusal 0028's
`edit` binding mints, and today's non-zero-exit command writer) must render
WITHOUT `detailMayHaveApplied`; and a read or gate refusal wrapping the typed
held-pipe error stays `Applied()` false (`0026:C1` `refusal:` "as today"). The
neighbouring `ClassTimeout` arm in `flow_exec.go::accessorFailureOf` keys on
phase; copying that shape would label 0028's pre-mutation refusals applied.

**Check.** A unit test on `accessorFailureOf`: (write phase,
`ClassExecutionFailure`, `Applied()` false) → no `detailMayHaveApplied`;
(read phase, refusal wrapping the typed held-pipe error) → `Applied()` false.
The `errors.As` for the held-pipe error lives in `Executor.Write`'s
`err != nil` arm, never in the shared `refusalOf` — `0026:A8`'s "either
placement compiles" is the trap.
