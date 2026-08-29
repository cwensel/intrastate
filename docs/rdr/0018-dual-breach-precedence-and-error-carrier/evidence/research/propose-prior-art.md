# RDR 0018 Propose — prior-art pass (Stage 2, bounded)

Scope: selection-depth reads only; spikes belong to Resolve.

## Claim 1 — dual-precondition precedence ordering (class + instance)

Queries (budget ≤3/claim):

1. `arc search semantic --corpus DevRef --limit 5 --json "report first
   validation failure fail fast versus aggregate all errors into one
   report"` — top 5 hits are CI/test-pipeline fail-fast prose
   (Refactoring 2nd ed. p120, Continuous Delivery p205, DevOps Handbook
   p193, Effective Debugging p181, Writing Effective Use Cases p73);
   none address API entry-precondition ordering. REJECTED (off-class).
2. `arc search semantic --corpus StateMachineRes --limit 5 --json
   "which error is reported when a state machine input violates
   multiple entry preconditions ordering"` — top hits: stateless
   CHANGELOG transition-precedence fix, awf-cli error-code taxonomy,
   javascript-state-machine API docs. None decide multi-precondition
   report order. REJECTED (off-class).

Verdict: ⚠ no prior-art coverage for dual-precondition precedence
ordering in the corpora; the ordering choice rests on repo-internal
anchors (HEAD order, `0009:C3`'s value-vs-tuple rationale, the CLI's
per-side rendering richness). The "0009 reports first at HEAD" claim
is an emission/order claim — demoted to a Resolve assumption
(Method: MVV Test), never quote-confirmed.

## Claim 2 — error-carrier convention (sentinel vs typed vs none)

Accepted citation (opened directly; langref is the sanctioned Go/CLI
fallback per `.rdr/resources.md`):

- `langref/uber-go-guide/style.md` §Error Types (lines 755-781):
  "Does the caller need to match the error so that they can handle
  it? If yes, we must support the errors.Is or errors.As functions by
  declaring a top-level error variable or a custom type." Guidance
  table row: Error matching **Yes** + **static** message → "top-level
  `var` with errors.New". ⇒ 0008's breach needs matching (the pinning
  test and producers branch on it) and its messages are static per
  channel → exported sentinel var, no custom type.
- `0009:C4` (peer RDR, projector): the in-repo precedent for the full
  typed-error surface — justified there by a structural payload
  (RowRef identity for findings). 0008's breach carries no payload a
  caller needs (the key is a compile-time constant; `0008:C4` fixes
  first-breach single error), so the typed half is not mirrored.

## Claim 3 — aggregate-both alternative's prior art

- `langref/kubebuilder/hack/docs/internal/cronjob-tutorial/webhook_implementation.go:113-126`:
  Kubernetes-style validation appends every violation to
  `field.ErrorList` and returns one invalid error — the real-world
  precedent for "join both so order stops mattering". Context check:
  that pattern serves user-facing declarative input validation.
  Intrastate's dual breach travels the programmer-mistake Go-error
  channel (`0008:C4`: a producer "is not owed an exhaustive list"),
  and `0009:C4` fences the returned aggregate — VERBATIM
  `CheckValid` error, `Unwrap() []error` exactly one level deep,
  every element `*EscapeShapeBreachError`. A cross-precondition join
  cannot satisfy that fence. Accepted as the alternative's honest
  prior art; rejected as the choice on fence contradiction.

## Repo-internal anchors read (freshness + instance question)

- `internal/resolve/resolve.go::Resolve` — `in.Table.CheckValid()`
  (line 444) precedes `CheckInput(in)` (line 452); comment block
  states placement "a cheapness preference, not an observable
  contract. JDR 0001 §JD-5 leaves the relative order ... open".
- `internal/resolve/precondition.go` — three unexported plain
  sentinels (`errReservedOwnedTag`, `errReservedObservedTag`,
  `errReservedRequiresOwned`); first-breach-wins.
- `internal/cli/flow_resolve.go::kernelResolveFailure` — the only
  consumer discrimination today: `errors.Is(err,
  resolve.ErrEscapeShapeBreach)`; everything else falls to generic
  `codeAccessorFailed`, comment: "A blanket recode would mislabel RDR
  0008's reserved-key breach."
- `0008:C4` (projector) — "if the table also breaches RDR 0009's
  shape rule, either error is conforming ... JD-5 may narrow this to
  a fixed order; nothing here forecloses that."
- `0009:C3`/`0009:C4` (projector) — precedence fenced only against
  *modeled dispositions*; aggregate shape and verbatim-return fenced.
- `docs/rdr/0009-.../artifacts/req-list.md` REQ-19/21/22 + the
  REQ-22 ASSUMPTION ("no test asserts that assemble was not called").
- `internal/table/category.go::CatReservedTagKey` = literal
  `reserved_tag_key` — the existing name for this breach category on
  the authored path.
- Dual-breach construction EXISTS in test at HEAD —
  `internal/resolve/reserved_key_0008_test.go::TestReq44And45And90_ReservedKeyBreachIsNotSkippedByACoincidentShapeBreach`
  (found by the Stage-2 grounding sweep; the initial
  `rg "dual.?breach|doubly.?breach"` sweep missed it because the
  test names the not-skipped clause, not the phrase). It asserts
  only that some breach reports, never which — so the ordering gap
  stands; only the "no test constructs a dual breach" phrasing was
  wrong and has been corrected in the record.
