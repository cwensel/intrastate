Model: gpt-5

# 3-amigo — QA / Tester

1. **QA-1 — Critical: locator/model-rule disagreement has no writable
   pass/fail test.** Trigger: `Critical Assumptions / A4` and the locator
   agreement `Normative Contract`. The relation "designates" and callable
   validation boundary are undefined, and the matrix omits the mismatch case.
   **Prevents:** proving a conflicting locator is rejected before selection.
   This re-raises COVE-3/A4.
2. **QA-2 — High: the prescribed mutation cannot prove `Writes` isolation.**
   Trigger: Validation scenario 3. Replacing a slice does not mutate its backing
   storage. **Prevents:** an adversarial test that fails when
   `TransitionPlan.Action.Writes` aliases the selected edge.
3. **QA-3 — High: expansion-suffix preservation lacks positive coverage.**
   Trigger: `Technical Design`, `Normative Contracts`, and Validation scenarios
   1–2. "Complete identity" does not require either success fixture to carry a
   non-empty suffix. **Prevents:** detecting accidental omission of the suffix.
4. **QA-4 — Medium: "programmer/invariant error" has no testable
   classification.** Trigger: `Normative Contracts` and Validation scenario 4.
   QA can assert only a non-nil error and empty disposition. **Prevents:** a
   stable error-class test distinguishing malformed identity from unrelated Go
   errors.
5. **QA-5 — Medium: replay-consumer acceptance is absent.** Trigger: `Problem
   Statement` versus `Validation / Testing Strategy`. Resolver copy tests are
   specified and CLI acceptance is delegated, but no replay-facing assertion
   is named. **Prevents:** proving replay consumes selected identity unchanged
   rather than reconstructing it.

Lower-severity overflow: 1 — whether a write-free escape exposes nil or non-nil
empty `Writes`, which may affect later JSON projection.
