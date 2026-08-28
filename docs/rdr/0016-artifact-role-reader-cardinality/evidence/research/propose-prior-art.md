# 0016 Propose — prior-art search record

Stage 2 prior-art pass (2026-08-28). Problem class: a resolution key
(artifact role) can match multiple registered handlers (read bindings);
peers either refuse the ambiguity at registration/validation, define a
deterministic total order, or aggregate over the set.

## Queries run (budget: ≤3/claim, ≤5 opened hits/claim)

1. `arc search semantic --corpus StateMachineRes --limit 5 --json
   "conflicting transitions multiple match resolution document order determinism"`
   — top hits were heading-only chunks (`resolver-vs-roundtrip.md` §3 header,
   `ms-conductor/docs/design/registry.md` "Design decisions" header) and
   off-topic merge-conflict docs. No quotable passage on handler
   cardinality. REJECTED branch: did not open the ms-conductor registry doc
   further — headings gave no cardinality signal and the budget stops at
   first non-confirmation.
2. `arc search semantic --corpus DevRef --limit 5 --json
   "duplicate registration ambiguous handler binding refuse at registration unique constraint"`
   — hits were database uniqueness-constraint pages (DDIA p543, db books);
   analogous (uniqueness enforced by the system, not by caller ordering) but
   not about handler registries. Not load-bearing; not cited.
3. `arc search semantic --corpus StateMachineLit --limit 5 --json
   "SCXML statecharts conflicting enabled transitions selection document order priority"`
   — Harel statecharts paper chunks, none on conflict-selection ordering.
   Not cited.

Verdict: ⚠ no prior-art coverage in the indexed corpora for
handler-cardinality-at-registration; the class read below comes from an
openable local source and in-repo precedent instead.

## Accepted citations

- **Go stdlib `net/http` ServeMux** (go1.26.6,
  `$GOROOT/src/net/http/server.go::ServeMux.registerErr`, ~line 2912):
  registration refuses a conflicting pattern with
  `pattern %q (registered at %s) conflicts with pattern %q (registered at
  %s)` — ambiguity that no specificity rule resolves is refused **at
  registration time**, never resolved by registration order at serve time.
  ⇒ supports refusing same-role duplicate read bindings at model load
  rather than defining a runtime tiebreak.
- **In-repo precedent** (the sibling-path exhibit, step 5):
  - `internal/accessor/validate.go` multiply-bound-identity arm
    (`CodeMultiplyBoundAccessor`): a second binding of the same
    `(flow, name, capability)` triple is a validation finding — uniqueness
    of a resolution key is already refused, not tiebroken.
  - `internal/table/load.go::checkAccessorBindings`: per-key reader arity
    ("owned tag … served by %d readers; want exactly one") — load-time
    arity walks over bindings are the established mechanism and category
    (`CatMalformedAccessorBinding`).
  - `internal/cli/flow_state.go::writerFor`: exactly-one-writer arity
    enforced at the CLI seam (kata 8dg3 test) — the codebase's uniform
    answer to N-candidates-one-key is refusal, not ordering.

## Freshness/grounding reads made during this pass

- `internal/accessor/model.go::readerFor` (231–238): first-match over
  `reg.Definitions` order on `Capability==CapRead && Accessor.Role==role`.
- `internal/accessor/validate.go::Validate`: eight arms, none over role.
- `internal/accessor/executor.go::nonOwnedPlanKeys` (469–): a planned key
  outside the writer definition's own `keys` is refused before the write
  command — planned keys ⊆ writer's declared keys at execution.
- `internal/accessor/executor.go` (~275, ~363): read-back resolves
  `readerFor(def.Accessor.Role)`; compared set = planned keys ∪ protected
  (reader-declared minus planned).
- `internal/table/model.go::Accessor` (144–150): `Role`, `Keys` available
  to load-time walks for both `Readers` and `Writers` maps.
