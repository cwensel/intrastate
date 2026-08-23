Model: claude-opus-5[1m]

# Charted to successor — critique iteration 3

Net-new scope, recorded durably and dismissed from this loop. Neither was folded
into RDR 0003 (the scope-expansion wormhole).

- **R-10 — RDR 0006's default-on rider makes a docs-only declaration a CI
  break.** Declaring a tag's finite domain for documentation reasons silently
  opts an existing row group into a blocking exhaustiveness gate, because
  "claims closed coverage" is default-on for every group whose participating
  dimensions are all finitely declared (`docs/rdr/0006-graph-lint-authority-and-guarantees.md:913-915`).
  The cheapest author fix is reverting the declaration — the opposite of the
  incentive the model wants. **Out of scope here**: the rider is RDR 0006's
  decided text, and this RDR reads it rather than owning it; reversing it from
  this document would re-decide a peer's rider unilaterally, the exact failure
  §JD-4 exists to prevent. **Suggested successor**: raise at RDR 0006's refine
  alongside the §JD-14 invariant repair — either a grace path for newly-declared
  domains, or an explicit acknowledgement that declaring a domain is opting in.

- **R-18 — the tag declaration model shipped without an alternatives pass.**
  RDR 0003 analyses six alternatives for the operator grammar and zero for the
  five-field type system it absorbed at Stage 6 (value kind, finite domain,
  optionality, single-valuedness, element universe). A11's only rejected
  alternative is a *placement* (RDR 0002 vs. here), not a design. Open questions
  nobody has examined: no nullable-vs-optional distinction, no `int` domain form
  beyond `{min..max}`, no string/scalar domain, and no versioning story for
  adding a value to a shipped enum domain. **Out of scope here**: authoring an
  alternatives pass for the declaration model is a Stage-2-shaped body of work on
  a document at lock, and the model as written is internally consistent and
  computable — the gap is review coverage, not a known defect. **Suggested
  successor**: a follow-up RDR seed for enum-domain evolution/versioning, which
  is the one of the four that will bind first in practice. The aggregate
  open-assumption assessment in `Assumption Verification` names this as the thing
  to watch.
