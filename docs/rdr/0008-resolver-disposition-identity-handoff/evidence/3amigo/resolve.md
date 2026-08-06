Model: gpt-5

# 3-amigo resolve — RDR 0008

- **fixed** — origin: `AMIGO-1` / `PM-1` / `IMP-1` / `QA-1`; sections
  touched: `Critical Assumptions / A4`, `Approach`, Validation scenario 6. The
  draft now requires Stage 6 to name the normalized-table constructor or
  validator and its mismatch test before lock, keeps `plan` as a copying
  consumer, and makes the normalized-table rejection test explicit. A4 remains
  Pending rather than being silently treated as settled.
- **fixed** — origin: `IMP-2`; sections touched: `Approach`, `Implementation
  Plan / Phase 1`. RDR 0008 owns one resolver-facing Go carrier while RDR 0002
  owns identity semantics and population; normalized edges and plans share the
  same `SelectedRule` value without a parallel identity struct.
- **fixed** — origin: `AMIGO-2` / `IMP-4` / `QA-2`; sections touched: `Minimum
  Viable Validation`, Validation scenario 3. The implementation test must
  mutate a write element in the caller-owned backing array, matching the A2
  spike, before replacing or reusing the edge.
- **fixed** — origin: `AMIGO-3` / `IMP-5` / `QA-3`; sections touched: `Minimum
  Viable Validation`, Validation scenarios 1 and 7. The ordinary and replay
  fixtures must carry a non-empty expansion suffix.
- **dismissed-with-cite** — origin: `AMIGO-4` / `PM-3` / `QA-5`; sections made
  explicit: `Technical Design`, Validation scenario 7. RDR 0001 `Normative
  Contracts` and `Minimum Viable Validation` already assign replay equality to
  the resolver, and `internal/resolver/resolver_test.go::TestResolve_ReplayReturnsValueIdenticalPlans`
  is the existing owner; this RDR extends that test with selected identity.
- **fixed** — origin: `PM-2`; section touched: `Implementation Plan / Phase 3`.
  Carrier plus resolver/replay tests complete RDR 0008 implementation;
  operator-visible diagnostics require the separately owned RDR 0005 consumer
  acceptance to land.
- **dismissed-with-cite** — origin: `IMP-3` / `QA-4`; section touched: none.
  RDR 0001 `Normative Contracts` requires stability for modeled refusal kinds
  while reserving the Go error path for programmer mistakes; it does not make a
  stable programmer-error taxonomy part of the kernel contract. A non-nil
  error and empty disposition are sufficient here.
- **dismissed-with-cite** — origin: implementer overflow on exported `Edge`
  migration; section touched: none. `Consequences` already requires `Edge`,
  `TransitionPlan`, construction code, and fixtures to change in lockstep; the
  repository-internal literal migration is implementation work, not a second
  contract.
- **charted-to-successor** — origin: implementer/QA overflow on nil-versus-empty
  `Writes`; section touched: none. **Charted:** RDR 0005 consumer acceptance
  should pin JSON's empty planned-write representation if it is wire-visible;
  RDR 0008 owns semantic no-writes and the internal defensive copy, not JSON
  rendering.

Needs verification:

- A4 remains Pending. Stage 6 must verify locator/model-rule agreement against
  RDR 0002 and name the callable normalized-table validation boundary plus its
  mismatch test before lock.
- Confirm RDR 0001's resolver implementation is present on the implementation
  baseline before Stage 8 work begins.

Tiebreakers: none.
