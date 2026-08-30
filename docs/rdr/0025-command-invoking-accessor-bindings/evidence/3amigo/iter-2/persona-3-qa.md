Model: claude-opus-5[1m]

# Persona 3 — QA / Tester (iter-2, delta)

Delta scope: 0025:§testing-strategy scenarios 1, 3, 4, 4b, 5, 5b, 7; contracts
0025:C3, 0025:C4, 0025:C5, 0025:C6; Failure Modes 0025:F7 (Windows).

Method: read each element once via the projector; verified every fixture
citation against the fixture file directly and every code anchor against
source in /Users/cwensel/sandbox/newcoinc/intrastate.

## Closed by the rewrite

- **Scenario 3 fixture citations — VERIFIED, closed.** This was the defect
  class this pass exists to catch, so I checked each mapping line by line
  against `evidence/spikes/a3-a7-a8-real-tools/output.txt`. Every asserted
  value is present in the cited case:
  - "E1b/E1e supply the mapped-verdict codes" — E1b `exit=1 stdout=""`
    (line 118–119), E1e `exit=1 stdout=""` (line 140–141). Both exit 1,
    both empty stdout, matching the `exit_verdicts = { "1" = "deny" }` arm.
  - "E1c/E1f the error codes" — E1c `exit=2 stdout=""` (line 125–126),
    E1f `exit=128 stdout=""` (line 147–148). These back cases 5 (gate exit 2
    under the map) and 7 (read exit 128) respectively.
  - "R6/E1g the established-absent read" — R6 `exit=1 stdout=""` (line
    102–103), E1g `exit=1 stdout=""` (line 155–156). Both back the
    `exit_absent = [1]` arm.
  - "every case of which is empty-stdout" — confirmed: E1a, E1b, E1c, E1d,
    E1e, E1f, E1g, E1h and R6 all carry `stdout=""`. The claim holds with no
    exception.
  The empty-stdout/non-empty-stdout split is stated and correct: cases 3–7
  are the five empty-stdout arms, cases 1/2/8 the three non-empty ones, and
  the record explicitly declines to cite FX-exit-codes for the latter with a
  stated reason (contract-derived from C3/C4 ordering). Eight cases, eight
  ordered Expected values, one-to-one. A tester can write all eight.

- **Scenario 3 stderr-tail oracle — closed.** "asserted non-empty on case 3,
  where no applied-sense text competes for `CLIError.Detail`" is grounded:
  `internal/cli/flow_exec.go::accessorFailure` sets
  `ce.Detail = detailMayHaveApplied` only on `ClassTimeout` at `phaseWrite`
  (line 378) and `ClassReadBackIncomplete` (line 397). The
  `ClassExecutionFailure` arm (line 385–387) sets no `Detail`. Case 3 is an
  `execution_failure`, so the field is genuinely uncontested. C4's "NEW
  `Detail string` field on `accessor.Refusal`" is also accurate — `Refusal`
  (`internal/accessor/model.go:280`) has Class/Accessor/Capability/Role/
  Timeout/Keys and no `Detail`; the existing `Detail` is on `Finding`
  (line 139), a different type.

- **Scenario 1 oracle strength — closed.** Pass/fail is now assertions, not
  adjectives: per-mutant named `table.Category`, "never by 'validation
  returned non-empty'", membership in `table.Categories()` at the tail in
  clause order, and `ValidationCodes` at eight. Grounded: `validationCodes`
  in `internal/accessor/model.go:117-126` has exactly eight members and
  `validation_0004_test.go:182-190` already asserts the set by equality, so
  the "stays at eight" arm is a live regression oracle. The open-deny-list
  negative (`perl -e` loads) gives C5's `interpreter set: OPEN` its own
  falsifier — previously unfalsifiable.

- **Scenario 4 ablations — closed.** The prior concern was two arms with no
  runnable oracle. The rewrite demotes them from test arms with a stated
  reason (no injection seam; adding one would ship a way to weaken the
  deadline) and cites S3/S1 as evidence instead. FX-deadline checks out:
  S2 1004ms/0 survivors, S4 1002ms/0, S5 1002ms/0 — "S2/S4/S5 ≈1.0s with
  0 survivors" is exact. S3 1501ms with `survivors=1 [79351]` (the orphan)
  and S1 blocked past the 5s watchdog are both as described. Three test arms
  with a numeric bound and a zero-survivor count: writable.

- **Scenario 4b env composition — closed.** C4 now fixes precedence
  explicitly (`overlay > entry env > env_pass > parent allowlist`) and the
  scenario asserts on the child's observed environment "not on the read
  result, since an env defect can leave the value correct" — which is the
  right oracle, and the fixture proves the hazard (V4c reads `draft` under
  `GIT_CONFIG_COUNT` injection while the no-`--file` control V4c'' reads
  `FROM_ENV_INJECTION`, lines 250–272). The `LC_ALL` passes / `LCFOO` does
  not arm gives C4's literal-prefix rule a falsifier. FX-raw-read is exact:
  R1 line 15 `stdout="draft\n"`, line 18 `RAW-mode bind: key=state.phase
  value="draft" (TrimRight \n)` — the asserted `draft` is in the fixture.

- **Scenario 5 anchors — closed.** `internal/accessor/executor.go:96` lands
  on the `<clear>`-is-UNREADABLE comment block (lines 96–100) with
  `raw.classify(requested, true)` immediately below — the anchor is correct
  to the line. Spike R3b (line 42–49: `exit=0`, `artifact_after` still
  `phase = draft`) and R4b (line 76–83: `exit=0`,
  `artifact_after="{\"state.phase\":\"review\"}"`) both support their
  asserted outcomes, and R4c (line 85–91, `exit=128`, `bad config line 1`)
  supports `read_back_incomplete`. "the write's exit status decides none of
  them" is the assertion that makes this falsifiable.

- **Scenario 5b two-halves oracle — closed.** Both halves are grounded and
  separately assertable: `internal/cli/flowbind/registry.go:26` uses
  `slices.Sorted(keys(m.Readers))`, and `internal/accessor/model.go:231`
  `readerFor` returns the first match. Asserting "by the selected reader's
  identity, never by 'read-back succeeded'" is the correct oracle — under
  the weak form either reader passes.

- **Scenario 7 caller enumeration — closed and verified.** Exactly three
  non-test `NewExecutor` callers exist: `internal/cli/flow_state.go:306`,
  `internal/cli/flow_exec.go:256`, `internal/cli/flow_exec.go:644`. The
  record's "(`flow_state.go`, `flow_exec.go` ×2)" is exact. The
  absence-of-spawn oracle (sentinel argv0 leaving an observable trace) is a
  real assertion, not "non-zero exit". The apparent C6-vs-scenario wording
  difference ("command binding's constructor" vs "registry constructor") is
  not a defect: `internal/cli/flowbind/registry.go:20`
  `func Registry(m *table.Model) accessor.Registry` is the site where
  command bindings are constructed from the model, so both phrasings name
  the same seam and the scenario's assertions hold under either reading.

## Findings

### 1. MEDIUM — 0025:F7 / 0025:C4 `platform`: the Windows refusal has no scenario and no reachable oracle on CI

`0025:F7` states a command entry invoked on Windows refuses
`execution_failure` naming the unsupported platform before spawn, and C4
fixes the mechanism as "a runtime `runtime.GOOS` check in the command
binding". C4 explicitly rejects a build constraint on the grounds that it
"would make the refusal unbuildable-on-Windows rather than observable" —
so observability is the stated reason for choosing this mechanism.

But no scenario in §testing-strategy exercises it. The projector summary
shows `0025:F7` is the only platform-bearing element in the record; scenario
7 covers the C6 gate, scenario 3 the exit/envelope ordering, neither touches
`GOOS`. And as written the branch is unreachable from the Unix CI the rest
of the strategy assumes: a bare `runtime.GOOS` read admits no injection, so
a darwin/linux runner can never take the windows arm. The mechanism was
chosen for observability and then left with nothing that observes it.

The record already knows how to handle this shape — scenario 4 demotes the
two deadline ablations to non-arms *with a stated reason* rather than
leaving them as untestable arms. F7 gets neither treatment: it is neither a
scenario nor explicitly declared untestable.

Two ways to close, both cheap: either state that the check reads a
package-level `goos` variable (or equivalent seam) that the test sets, which
makes the refusal assertable on any host and is the only form under which
C4's "observable" claim is true; or say plainly, as scenario 4 does for the
ablations, that F7 is asserted by inspection and carries no test arm.

**Test prevented**: "a command entry on an unsupported platform refuses
`execution_failure` naming the platform, before any spawn" — cannot be
written to run in CI, so C4's `platform` clause is the one normative line in
the delta with no failing test when its rule is dropped. That is a direct
miss against the stated coverage goal ("every C1–C6 clause has a test that
fails when its rule is dropped").

### 2. LOW — 0025:S3: E1h is cited as a fixture arm but no case exercises a spawn error

Scenario 3's fixture paragraph enumerates what FX-exit-codes supplies and
ends with "E1h the spawn error". E1h is real (fixture line 160–165:
`spawn_error=exec: "definitely-not-a-binary-xyz": executable file not found
in $PATH  is_ErrNotFound=true  exit=(none)`), but the scenario's case list
has eight cases and none of them is a missing executable — they are the
deny-envelope, malformed-envelope, unmapped non-zero, `exit_verdicts` hit,
`exit_verdicts` miss, `exit_absent` hit, exit 128, and over-cap stdout.

This is the inverse of the defect class this pass targets: not an asserted
value absent from the fixture, but a cited fixture arm with no assertion
consuming it. A tester following the citation list will look for a ninth
case and not find one.

Note the substantive point underneath: a spawn error is genuinely different
from a non-zero exit — there is no exit code at all (`exit=(none)`), so
neither `exit_absent` nor `exit_verdicts` can match, and the classification
has to come from somewhere else. C3's `exit_absent`/`exit_verdicts` lines
are both phrased as "a listed exit with empty stdout", which says nothing
about the no-exit case.

Either add a case 9 (unresolvable argv0 refuses `execution_failure` before
any envelope parse, per E1h) or drop E1h from the citation list.

**Test prevented**: "a command whose argv0 does not resolve refuses
`execution_failure`, and no exit map can claim it" — currently unwritten,
and the citation implies it is covered.
