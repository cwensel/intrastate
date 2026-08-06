Model: gpt-5

# COVE findings — RDR 0008

## Step 0 — grounding sweep

| Claim | Verdict | Ground |
| --- | --- | --- |
| The resolver seam is `internal/resolver/resolver.go::{Edge,TransitionPlan,Resolve,plan}`. | NOT-FOUND on `main`; CONFIRMED on the current integration branch. | `git show main:internal/resolver/resolver.go` fails; the current branch defines all four symbols. |
| `Edge`, `TransitionPlan`, and `plan` discard selected-rule identity. | CONFIRMED on the current integration branch. | `internal/resolver/resolver.go::Edge`, `::TransitionPlan`, and `::plan` carry/copy only match and action fields. |
| Ordinary and modeled-escape success share one construction boundary. | CONFIRMED on the current integration branch. | `internal/resolver/resolver.go::Resolve` sends both exact-one branches through `plan`. |
| `plan` already owns the defensive action copy. | CONFIRMED on the current integration branch. | `internal/resolver/resolver.go::plan` uses `maps.Clone` for `NextTags` and copies `Writes`. |
| No sibling selected-rule carrier already exists. | CONFIRMED on the current integration branch. | Repository-wide search of Go source found no `SelectedRule`, `RuleID`, `SourceLocator`, `ExpansionSuffix`, `ModelID`, or `TransitionAction` symbol. |
| RDR 0002 supplies model/rule identity, expansion suffix, and source locator on normalized rows. | CONFIRMED. | RDR 0002 `Technical Design`, `Normative Contracts`, and `Load-Bearing Decisions / Identity`. |
| RDR 0005 is a translating consumer of matched-rule identity. | CONFIRMED as a contract; no production consumer exists yet. | RDR 0005 `Technical Design` and `Normative Contracts`; code search finds no implemented resolver CLI surface. |
| RDR 0007 owns table-revision binding. | CONFIRMED in its current Draft. | RDR 0007 `Approach`, `Technical Design`, and `Normative Contracts`. |

## Step 1–2 — independent verification questions

1. **Does every cited resolver symbol resolve on `main`?** No. `internal/resolver/resolver.go` is absent from `main`; it is introduced by predecessor implementation commits on the current integration branch.
2. **Does the current `Edge` retain any selected-rule identity or source locator?** No. `internal/resolver/resolver.go::Edge` has `Outcome`, `EscapeFor`, `Guard`, `NextTags`, and `Writes` only.
3. **Does the current `TransitionPlan` expose only action data?** Yes. `internal/resolver/resolver.go::TransitionPlan` exposes only `NextTags` and `Writes`.
4. **Do ordinary and modeled-escape success use the same constructor?** Yes. Both exact-one branches in `internal/resolver/resolver.go::Resolve` call `plan`.
5. **Does that constructor isolate returned action storage from the supplied table?** Yes for the current action shape. `internal/resolver/resolver.go::plan` clones the tag map and copies the write slice; `TestResolveAdversarialReturnedPlanCannotMutateReplay` covers the reverse mutation direction, and the RDR 0008 spike covers input mutation/reuse.
6. **Is there another Go symbol that already carries the proposed identity?** No. The Go-source sweep found none, so the proposal extends the existing result boundary rather than duplicating a sibling discriminator.
7. **Does RDR 0002 require ordinary and escape rows to retain the source information RDR 0008 needs?** Yes. Its `Normative Contracts` require source rule id and source locator on normalized rows, and its identity decision adds model id plus deterministic expansion suffix; modeled escapes are normalized candidate rows.
8. **Can the current production CLI prove the proposed projection without a second lookup?** No. RDR 0005 specifies the projection, but no resolver CLI implementation exists; the RDR 0008 implementation cannot own a production `flow resolve` fixture without crossing RDR 0005's boundary.
9. **Does the draft state who guarantees that `SourceLocator` points to the same model/rule as `ModelID` and `RuleID`?** No. It requires non-empty values and calls the locator opaque, but RDR 0002 says the locator identifies at least model id and rule id. A non-empty but conflicting locator would pass the described `plan` check.
10. **Will missing required identity take the Go error path on both ordinary and escape success?** Yes in the proposed shape: both successful branches call `plan(edge) (Disposition, error)`, and the normative contract assigns missing fields to the invariant-error path.
11. **Can a modeled refusal claim selected-rule identity under the proposed carrier?** No. Identity is nested in `TransitionPlan`; `internal/resolver/resolver.go::Disposition` keeps plan and refusal as separate branches, and the contract forbids identity when no row is selected.
12. **Does the draft accidentally make table revision part of selected-rule identity?** No. RDR 0008's `Load-Bearing Decisions / Identity` excludes revision, and RDR 0007 separately binds revision to the complete normalized table before selection.

## Step 3 — origin ledger

- **COVE-1 — REFUTED codebase claim:** the cited resolver symbols do not resolve on `main`, despite the Research Findings describing them as current. The predecessor implementation exists only on the integration branch, so landing it must be an explicit implementation prerequisite.
- **COVE-2 — internal ownership contradiction:** Phase 3 and Validation scenario 6 assign a production `flow resolve` payload fixture to this RDR while the draft repeatedly assigns CLI rendering and the consumer implementation to RDR 0005, whose production surface does not yet exist.
- **COVE-3 — RDR silence:** the draft requires non-empty `ModelID`, `RuleID`, and opaque `SourceLocator`, but does not require the locator to designate the same authored rule. RDR 0002 makes the locator identify at least model id and rule id, so disagreement would produce misleading diagnostics unless normalization rejects it.
