# RDR 0001 Phase 1 Coverage

Test file: `internal/resolver/resolver_test.go`

## Requirement-to-test map

| Requirement | Behavior test |
| --- | --- |
| REQ-1 | `TestResolve_ReplayReturnsValueIdenticalPlans` |
| REQ-2 | `TestResolve_RefusesUnmodeledOutcome` |
| REQ-3 | `TestResolve_AmbiguityRequiresOneExplicitEscapeEdge` |
| REQ-4 | `TestResolve_GuardsDistinguishTagProvenance` |
| REQ-5 | `TestResolve_ReturnsInertOwnedTagWrites` |
| REQ-6 | `TestResolve_GuardsDistinguishTagProvenance` |
| REQ-7 | `TestResolve_ReturnsInertOwnedTagWrites` |
| REQ-8 | `TestResolve_RefusesUnmodeledOutcome` |
| REQ-9 | `TestResolve_GuardsDistinguishTagProvenance` |
| REQ-10 | `TestResolve_GuardsDistinguishTagProvenance` |
| REQ-11 | `TestResolve_ReturnsInertOwnedTagWrites` |
| REQ-12 | `TestResolve_ReturnsExactlyOneStructuredDisposition` |
| REQ-13 | `TestResolve_ReplayReturnsValueIdenticalPlans`, `TestResolve_ReturnsExactlyOneStructuredDisposition` |
| REQ-14 | `TestResolve_RefusesUnmodeledOutcome`, `TestResolve_RefusesNoMatch`, `TestResolve_AmbiguityRequiresOneExplicitEscapeEdge`, `TestResolve_RefusesUnavailableOwnedStateAsValue`, `TestResolve_RefusesUnevaluableGuardAsValue` |
| REQ-15 | `TestResolve_RefusesUnavailableOwnedStateAsValue`, `TestResolve_RefusesUnevaluableGuardAsValue`, `TestResolve_ReturnsExactlyOneStructuredDisposition` |
| REQ-16 | `TestResolve_RefusalKindsAreStableValues` |
| REQ-17 | `TestResolve_RefusalKindsAreStableValues` |
| REQ-18 | `TestResolve_IgnoresAmbientProcessStateAndProducesNoSideEffects` |
| REQ-19 | `TestResolve_IgnoresAmbientProcessStateAndProducesNoSideEffects` |
| REQ-20 | `TestResolve_IgnoresAmbientProcessStateAndProducesNoSideEffects` |
| REQ-21 | `TestResolve_IgnoresAmbientProcessStateAndProducesNoSideEffects` |
| REQ-22 | `TestResolve_ReturnsInertOwnedTagWrites`, `TestResolve_IgnoresAmbientProcessStateAndProducesNoSideEffects` |
| REQ-23 | `TestResolve_ReturnsInertOwnedTagWrites`, `TestResolve_IgnoresAmbientProcessStateAndProducesNoSideEffects` |
| REQ-24 | `TestResolve_ReplayReturnsValueIdenticalPlans` |
| REQ-25 | `TestResolve_ReplayReturnsValueIdenticalPlans` |
| REQ-26 | `TestResolve_GuardsDistinguishTagProvenance` |
| REQ-27 | `TestResolve_AmbiguityRequiresOneExplicitEscapeEdge` |
| REQ-28 | `TestResolve_ReturnsInertOwnedTagWrites` |
| REQ-29 | `TestResolve_GuardsDistinguishTagProvenance` |
| REQ-30 | `TestResolve_ReplayReturnsValueIdenticalPlans` |
| REQ-31 | `TestResolve_ReturnsInertOwnedTagWrites` |
| REQ-32 | `TestResolver_PublicContractUsesOnlyStandardGoValues` |
| REQ-33 | `TestResolve_ReplayReturnsValueIdenticalPlans` |
| REQ-34 | `TestResolve_RefusesNoMatch` |
| REQ-35 | `TestResolve_AmbiguityRequiresOneExplicitEscapeEdge` |
| REQ-36 | `TestResolve_RefusesUnmodeledOutcome` |
| REQ-37 | `TestResolve_RefusesUnavailableOwnedStateAsValue`, `TestResolve_RefusesUnevaluableGuardAsValue` |
| REQ-38 | `TestResolve_GuardsDistinguishTagProvenance` |
| REQ-39 | `TestResolve_ReplayReturnsValueIdenticalPlans` |
| REQ-MVV | `TestResolve_MinimumViableValidation` |

## Orphan audit

- Requirements without tests: none.
- Tests without `REQ-*` headers: none.
- Non-test helpers are fixtures/assertion helpers only: `legalInput`,
  `noMatchInput`, `unmodeledOutcomeInput`, `ambiguousInput`,
  `unevaluableInput`, and `assertRefusal`.

## Red gate

- Command: `go test ./internal/resolver`
- Result: expected failure before implementation: `no non-test Go files`.
- New behavior tests: 13.
- Numbered requirements covered: 39 of 39, plus REQ-MVV.
- REQ-MVV runner after implementation:
  `go test ./internal/resolver -run '^TestResolve_MinimumViableValidation$'`.

## Phase 2 implementation validation

- Package command: `go test ./internal/resolver`
- Package result: `ok github.com/newcoinc/intrastate/internal/resolver`
- Full-suite command: `go test ./...`
- Full-suite result: PASS for all packages.
- Project gate command: `make check`
- Project gate result: PASS (`go vet`, `golangci-lint` with 0 issues, and
  race-enabled atomic coverage tests; resolver coverage 92.6%).
- REQ-MVV command:
  `go test -v ./internal/resolver -run '^TestResolve_MinimumViableValidation$'`
- REQ-MVV actual result: PASS. The value-level replay assertion
  `reflect.DeepEqual(first, second)` passed for the complete disposition, and
  the `no_match`, `ambiguous_match`, `owned_state_unavailable`,
  `guard_unevaluable`, and `unmodeled_outcome` value-refusal subtests each
  passed without using the Go error path.
- Orphan audit after implementation: unchanged; every REQ has a green behavior
  test and every behavior test has at least one REQ header.
