model: gpt-5.6-sol
variant: full (profile: foundational)

# Public API reconstruction

```go
package resolver

// SelectedRule is copied from the exact normalized row selected by Resolve.
type SelectedRule struct {
	ModelID        string
	RuleID         string
	ExpansionSuffix string
	SourceLocator  SourceLocator
}

// TransitionAction is independently usable by a state-application caller.
type TransitionAction struct {
	NextTags TagSet
	Writes   []OwnedTagWrite
}

type TransitionPlan struct {
	SelectedRule SelectedRule
	Action       TransitionAction
}

// GUESS: RDR 0008 requires an opaque, snapshot-owned table but does not name
// its concrete type.
type ValidatedTable struct {
	rows []normalizedRow
}

// GUESS: The construction boundary's name and input shape remain open under
// A5. The input form shown here permits validation but must be defensively
// copied and must not be accepted directly by Resolve.
func NewValidatedTable(rows []NormalizedRowInput) (ValidatedTable, error)

// GUESS: The RDR does not reproduce RDR 0001's exact Resolve signature or
// input type. This rendering preserves its required table-first boundary.
func Resolve(table ValidatedTable, input ResolveInput) (Disposition, error)

// GUESS: Inspection is optional and its exact shape is unspecified. If
// exposed, it returns values with no mutable backing storage shared with the
// table.
func (t ValidatedTable) Rows() []NormalizedRowView
```

`SourceLocator`, `TagSet`, `OwnedTagWrite`, `ResolveInput`, `Disposition`, and
the modeled refusal types are referenced contracts owned by peer RDRs; RDR 0008
does not determine their exact Go declarations.

Error modes:

- `NewValidatedTable` rejects empty model or rule identity, an invalid or zero
  locator, a locator whose model/rule differs from `SelectedRule`, and any row
  whose match inputs, selected identity, locator, and action cannot be bound as
  one owned snapshot. GUESS: these failures are returned as ordinary Go errors
  that the load/config layer later maps to a stable `CLIError`; the concrete
  error types and codes are not specified here.
- `Resolve` returns a Go error when an impossible invalid value appears inside
  an already validated table. This is a programmer/invariant failure, not a
  modeled refusal.
- No match, ambiguity, unavailable owned state, unevaluable guard, and an
  unmodeled outcome are modeled dispositions owned by RDR 0001, not Go errors;
  none claims a selected rule when no row was selected.
- GUESS: modeled escape returns an identity-bearing non-success disposition
  with no action. A6 explicitly leaves this fork unresolved, so a conforming
  implementation could instead make it a success only if normalization
  supplies a complete, source-defined action.

# Three most important internal helpers

```go
// validateAndOwnRow validates identity/locator agreement and clones every
// caller-owned map, slice, locator component, and action component into one
// stable normalized row. GUESS: exact name and parameters.
func validateAndOwnRow(in NormalizedRowInput) (normalizedRow, error)

// plan checks required selected-rule fields and returns defensive value copies
// of SelectedRule and TransitionAction from the exact ordinary row selected by
// RDR 0001 matching. The RDR names this helper but not its exact signature.
func plan(row normalizedRow) (TransitionPlan, error)

// cloneSelection clones the nested values without retaining mutable aliases.
// GUESS: this may be inlined into construction and plan if SourceLocator is
// deeply immutable and value-copyable; A2/A4 decide that contract.
func cloneSelection(row normalizedRow) (SelectedRule, TransitionAction)
```

# Boundary data model

The resolver boundary passes, but does not independently persist, these values:

- `ValidatedTable`: an opaque owned snapshot containing ordinary rows whose
  match inputs, `SelectedRule`, `SourceLocator`, and `TransitionAction` are
  inseparable after validation. GUESS: escape rows occupy a distinct internal
  row variant because their action/disposition meaning remains unresolved.
- `SelectedRule`: `ModelID`, `RuleID`, `ExpansionSuffix`, and typed
  `SourceLocator`. Logical selected-rule equality within one revision-bound
  input uses `(ModelID, RuleID, ExpansionSuffix)`; the locator is diagnostic
  provenance, not another identity component.
- `TransitionAction`: `NextTags` plus `Writes`. It contains no guards or match
  predicates and has no mutable backing storage shared with constructor input,
  normalized-table inspection, or another returned plan.
- `TransitionPlan`: the successful ordinary disposition payload, containing
  exactly one `SelectedRule` and its separately usable `TransitionAction`.
- `Disposition`: GUESS: a tagged result that contains exactly one success,
  refusal, or selected-escape branch. Its concrete representation is not fixed
  by this RDR. Kernel refusals without a selected row expose no identity.
- No standalone historical record is introduced. Table revision identity is
  external to these values and owned by RDR 0007 and its successor. RDR 0005
  projects the returned values into text and JSON but does not reconstruct
  them.

# Top-level operation

```text
function Resolve(validatedTable, input) -> (Disposition, error):
    assert validatedTable is the opaque production snapshot
    candidates = evaluate ordinary row match inputs against input

    if evaluating a guard fails:
        return modeled UnevaluableGuard disposition with no selected rule
    if required owned state is unavailable:
        return modeled UnavailableOwnedState disposition with no selected rule

    if count(candidates) > 1:
        return modeled Ambiguous disposition with no selected rule

    if count(candidates) == 1:
        selectedRow = candidates[0]
        if selectedRow identity or locator violates validated invariants:
            return programmer/invariant Go error
        selectedRule, action = clone exact selectedRow values
        transitionPlan = TransitionPlan(selectedRule, action)
        return successful disposition containing transitionPlan

    escapeRows = evaluate modeled escape rows against the same input
    if count(escapeRows) > 1:
        return modeled Ambiguous disposition with no selected rule

    if count(escapeRows) == 1:
        selectedEscape = escapeRows[0]
        if selectedEscape identity or locator violates validated invariants:
            return programmer/invariant Go error
        copiedRule = clone selectedEscape.SelectedRule
        GUESS: return identity-bearing selected-escape non-success
               disposition containing copiedRule and no action
        GUESS: alternatively, if A6 defines a complete normalized escape
               action, clone it and return a successful TransitionPlan

    if the outcome is unmodeled:
        return modeled UnmodeledOutcome disposition with no selected rule
    return modeled NoMatch disposition with no selected rule
```
