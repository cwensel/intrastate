# RDR 0001 Self-Verification

## Phase 3b — Adversarial Review

### ADV-1 — Unevaluable competitor must prevent a guessed transition

- **Failure mode:** A resolver may select one visibly matching edge while
  silently ignoring another candidate whose guard cannot be evaluated.
- **Failure Modes anchor:** failures remain visible typed refusals; the kernel
  must not guess a transition.
- **Catching test:**
  `TestResolveAdversarialRefusesInsteadOfGuessing/one_apparent_match_cannot_hide_an_unevaluable_competitor`.
- **Observed:** PASS — the disposition is `guard_unevaluable`, with no plan and
  no Go error.

### ADV-2 — Ambiguous escape rows must not become an implicit tie-breaker

- **Failure mode:** An implementation may treat the presence of escape data as
  permission to pick the first of multiple matching escape rows.
- **Failure Modes anchor:** ambiguous resolution is visible as the typed
  `ambiguous_match` refusal; guessing an edge is prohibited.
- **Catching test:**
  `TestResolveAdversarialRefusesInsteadOfGuessing/multiple_ambiguity_escapes_remain_ambiguous`.
- **Observed:** PASS — multiple matching ambiguity escapes return
  `ambiguous_match`, with no plan and no Go error.

### ADV-3 — A returned plan must not mutate the table used for replay

- **Failure mode:** A transition plan may alias table-owned next tags or write
  descriptions, allowing caller mutation to change later resolution and act as
  an implicit persistence side effect.
- **Failure Modes anchor:** direct persistence is prohibited, and diagnosis is
  grounded in a replayable input tuple and table revision.
- **Catching test:** `TestResolveAdversarialReturnedPlanCannotMutateReplay`.
- **Observed:** PASS — mutating the first plan leaves the supplied table intact;
  replay returns the original next tags and write descriptions.

## Phase 3b Verdict

PASS — all three adversarial probes pass against the current implementation;
no Phase 3c defect was found by this review.

## Phase 3a — Chain-of-Verification

- **Requirements challenged:** REQ-1 through REQ-39 and REQ-MVV (40 total).
- **Independent violating inputs:** repeated identical tuples; provenance-name
  collisions; zero, multiple, unavailable, unevaluable, and unmodeled
  candidates; explicit and competing escapes; returned-plan mutation; changed
  process arguments and environment; refusal-taxonomy and dependency-boundary
  probes.
- **FAIL-N entries:** None.
- **Probe commands:**
  - `go test -list . ./internal/resolver`
  - `go test -count=1 ./internal/resolver`
  - `go list -f '{{join .Imports "\n"}}' ./internal/resolver`
- **Observed test result:** `go test -count=1` reported
  `ok github.com/newcoinc/intrastate/internal/resolver`.
- **Category-to-test mapping:** each name below is emitted by
  `go test -list . ./internal/resolver`:
  - legal: `TestResolve_ReturnsExactlyOneStructuredDisposition`
  - provenance: `TestResolve_GuardsDistinguishTagProvenance`
  - replay: `TestResolve_ReplayReturnsValueIdenticalPlans` and
    `TestResolve_MinimumViableValidation`
  - copy: `TestResolve_ReturnsInertOwnedTagWrites` and
    `TestResolveAdversarialReturnedPlanCannotMutateReplay`
  - refusals: `TestResolve_RefusesUnmodeledOutcome`,
    `TestResolve_RefusesDeclaredOutcomeWithoutCandidateRowAsNoMatch`,
    `TestResolve_RefusesUndeclaredOutcomeDespiteCandidateRow`,
    `TestResolve_RefusesNoMatch`,
    `TestResolve_RefusesUnavailableOwnedStateAsValue`,
    `TestResolve_RefusesUnevaluableGuardAsValue`, and
    `TestResolveAdversarialRefusesInsteadOfGuessing`
  - escapes: `TestResolve_AmbiguityRequiresOneExplicitEscapeEdge` and
    `TestResolveAdversarialNoMatchEscapes`
  - ambient: `TestResolve_IgnoresAmbientProcessStateAndProducesNoSideEffects`
  - taxonomy: `TestResolve_RefusalKindsAreStableValues`
- **Observed direct imports:** `go list` emitted the standard-library packages
  `maps` and `slices`, one per line.

## Phase 3a Verdict

PASS — all 40 requirements were challenged independently; zero implementation
violations and zero FAIL-N entries were found.
