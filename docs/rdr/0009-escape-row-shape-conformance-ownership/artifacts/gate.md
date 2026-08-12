Model: claude-opus-5[1m]

# Finalization Gate — RDR 0009: Ownership of escape-row shape conformance

- **RDR**: `0009-escape-row-shape-conformance-ownership`
- **Date**: 2026-08-11
- **Profile**: foundational (re-validated below)
- **Mechanical pre-sweep**: PASS — `evidence/tooling-pass/tooling-pass.md`
- **Verdict**: **READY — Gate PASS**

## Contradiction Check

No contradictions found between research findings, design principles, and
proposed solution.

Three places where the draft could have contradicted itself were checked
directly, because each is where a review round rewrote the text:

**Research vs. solution.** The decisive research finding — RDR 0001 reserves
`Resolve`'s Go error return for programmer mistakes, and `Resolve` returns a
literal `nil` error on all five existing paths — is exactly what the chosen
approach consumes. Nothing in Research Findings argues for an enforcement point
the solution then abandons.

**The one finding that reversed.** A6's Propose-stage claim (load-time
rejection is the industry standard) was researched at Stage 4 and **falsified**:
SCXML leaves the enforcement point to the implementation. The RDR carries the
reversal honestly — Key Discoveries labels it "Researched and falsified," A6's
statement is restated to the falsified form, and its Consequence records that
the two-point enforcement now rests on in-repo authority (RDR 0002's Final
prohibition, RDR 0001's reserved channel) rather than external precedent. No
surviving prose still asserts the withdrawn claim, so the falsification is a
citation loss, not a contradiction.

**Principles vs. planned features.** The RDR's stated principles are that the
five-kind refusal taxonomy stays closed, that no peer's locked surface is
reopened, and that no consumer parses error prose. Each planned feature is
checked against them: the breach travels the error path, not a refusal kind
(taxonomy closed); the `clierr` envelope gains an `omitempty` field under the
type's own documented extension allowance, verified additive by A9 (0005's
refusal-to-exit-code mapping untouched); and every diagnostic is asserted
structurally — `errors.Is` on the sentinel, `Ref`/`Count` struct fields, the
stable `Code` — never on message text. The one apparent tension, that the
authored layer rejects an *empty* write block while the kernel predicate is
length-based and accepts an empty `Writes` slice, is stated as deliberate
asymmetry with its reason (an authored empty block is a statement of intent; an
empty Go slice is a representation artifact) rather than left as an
inconsistency.

## Assumption Verification

All ten Critical Assumption Evidence Records are internally consistent and
terminal.

**Status / Method / Evidence agreement.** Every record is `Verified` with a
Method in vocabulary (Source Search ×7, Spike ×1, Peer RDR ×1, Prior Art ×1)
and Evidence of the form that Method requires — `path::Symbol` for Source
Search, captured command output for the Spike, peer normative quotations for
Peer RDR, an external standards quotation plus query ledger for Prior Art. Five
records carry a qualifying clause on the Verified verdict (A2 wiring
obligation, A4 settled-closed, A6 verified-contrary-to-original, A7
structured-form-adopted, A9 amendment obligation); each qualifier narrows what
was verified rather than softening the status, and each names its residue as a
Prerequisite or a carried-forward dependency. Every "If wrong" is non-empty and
states a consequence specific to that assumption.

**No `Docs Only`.** Zero records use it, so nothing blocks on an unverified
documentation claim.

**No `Pending` / `Unverified`.** None remain, so the status-consistency rule is
satisfied vacuously — no settled-fact prose can be leaning on an unsettled
assumption.

**No self-reference, no adjacent-claim proof.** Every Source Search resolves to
product source under `internal/`, never to this RDR or its artifacts (C3
clean). The records prove their own statements rather than neighbours: A1
proves the error slot is unused by enumerating all five return paths *and*
their callers; A5 proves no production constructor by showing no production
file imports the package at all; A9 proves the envelope addition is additive by
checking the three shapes that would actually break (whole-JSON comparison,
`DeepEqual` over an unmarshalled map, `DisallowUnknownFields`) and finding none
in the repo, rather than merely asserting the allowance.

**Cited symbols resolve.** C5 re-grepped every anchor; all resolve in their
cited file (`resolve.go` ×17, `clierr.go` ×7, `root.go` ×3, `config.go`,
`fixtures_test.go`). `Table.CheckValid`, `ErrEscapeShapeBreach`, and
`EscapeShapeBreachError` deliberately resolve nowhere — they are the surfaces
this RDR introduces, pinned by its Normative Contracts, not claims about
shipped code. Two source facts underpinning load-bearing verdicts were
re-confirmed verbatim at this gate: `ExitCodeFor`'s
`case GroupUserEnv, GroupInternal: return 2`, and the `CLIError` doc comment's
"Extend with new optional fields as needed — keep them `omitempty`" allowance
alongside `Cause`'s now-stale "the wire-visible cause surface is Detail."

**Assumptions the implementation must not lose.** A4's verdict rests on a peer
contract (RDR 0004 scoping the write accessor to planned writes) rather than on
any kernel-local property, and A9's field addition requires amending the stale
`Cause` doc comment in the same change. Both are bound as Prerequisites with
test forms, not left as notes.

## Scope Verification

The Minimum Viable Validation is **in scope and executed by this
implementation** — not deferred.

The MVV is: a hand-constructed table containing an escape row with populated
`Writes` causes `Resolve` to return a non-nil error naming that row's `RuleID`
and `SourceLocator`, with no `Plan` and no `Refusal`, while the same table with
the writes removed resolves, escapes, and refuses exactly as the frozen suite
proves today.

It is stated stepwise as three normative fixtures, each with an expected
end-state, and each traceable to a test scenario this RDR runs:

1. **breach-yields-error-not-plan** → scenario 1. The specific proof: an escape
   row `RowRef{"rdr.escape.needsowned", "flows/rdr.toml:90"}` carrying
   `Writes: []Tag{{Key: "status", Value: "Escaped"}}` yields the zero `Result`
   plus an `errors.Join` aggregate whose `Unwrap() []error` has length 1 and
   whose single `*EscapeShapeBreachError` carries `Count == 1` and that `Ref`,
   recovered via `errors.As`/`AsType` with `errors.Is` classifying the
   sentinel. No assertion parses message text.
2. **dormant-row-still-errors** → scenario 2, pinning whole-table scope.
3. **empty-not-nil-conforms** → scenario 3, pinning the length-based predicate,
   with the discriminating assertion on the *input row* rather than the plan
   (both nil and empty reach the plan as `len == 0`, so a plan-side assertion
   would not separate the two states).

The fixture values are read from the frozen suite and the A3 spike runs, not
invented — the `{status Escaped}` override and the 154-PASS conforming baseline
were both re-confirmed in the working tree at this gate.

Scope is stated explicitly rather than left ambiguous: scenarios 1–7b, 9, 10,
and 10b execute here; scenario 8 binds RDR 0002's build and scenario 11 binds
the future `flow` verb, both marked DEFERRED with the reason (no `flow` verb
exists at HEAD, per A5). Nothing the MVV depends on sits behind a deferral —
all three MVV fixtures are kernel-side, and the kernel is implemented.

Non-vacuity is proven rather than assumed, which is what makes this MVV a proof
and not a smoke test: scenario 6 runs two named mutants and records the honest
result that Mutant A (`CheckValid` returns nil) leaves the frozen suite green —
because after Phase 2 no pre-existing fixture breaches — which is precisely why
scenarios 1–4 must exist. Mutant B proves the check is live.

## Cross-Cutting Concerns

**Concurrency model** — the kernel is stateless by contract, and this RDR keeps
it so. The precondition is a pure read-only scan over the caller's own table
value; no caching or memoization is introduced, explicitly because a cache
keyed on a caller-supplied value would reintroduce the ambient state RDR 0001
forbids. Addressed here.

**Canonical form / determinism (value-level)** — this RDR makes one determinism
claim: a multi-breach report is a function of the table *value*, not of row
order. It is addressed structurally: ordering by `RowRef` identity via the
existing `compareRefs`, equal identities collapsed to one entry (because
`compareRefs` does not totally order distinct rows), per-identity pre-collapse
`Count` carried as a struct field, and a positional tiebreak explicitly
forbidden — the rule REQ-2/REQ-10 bind and that ADV-3b and Fixup-3c freeze.
Scenarios 10 and 10b assert it across both permutations and both identity
classes. The byte-level determinism checklist does **not** apply: no hash, no
canonical serialization, no replay-stable byte output is introduced.

**Incremental adoption** — the guarantee lands staged, and the RDR says so
rather than claiming immediate coverage. Phases 1–2 close the whole reachable
producer population today (A5: no production `resolve.Row` constructor at
HEAD); Phase 3 ships fixtures binding RDR 0002's future normalizer build. The
invariant holds everywhere a row can be constructed at this RDR's completion,
and becomes an authoring-time guarantee when RDR 0002 is built.

**Error envelope / exit-code policy** — owned by **RDR 0005**, which this RDR
conforms to rather than reopening. The breach routes through the existing
`GroupInternal` → exit 2 mapping with a stable `CLIError.Code`; the new
identity-carrying field is additive under `CLIError`'s own `omitempty`
extension allowance (A9). The one documentation defect the addition creates —
`Cause`'s stale "the wire-visible cause surface is Detail" — is bound as a
same-change Prerequisite alongside `docs/cli-output-contract.md`'s field list.

**Refusal taxonomy** — owned by **RDR 0001**, deliberately closed at five kinds
and gaining no member here. This is the RDR's central alignment claim and its
reason for rejecting Alternative 2; RDR 0007's precedent (rejecting
`guard_input_missing` at the adjacent seam) is weighed in Decision Rationale.

**Row-shape vs. graph-lint authority** — the boundary with **RDR 0006** is
stated explicitly: 0006's blocking authority covers graph invariants over the
normalized model and names parse/render fidelity as RDR 0002's. Row shape sits
on the RDR 0002 side of that line, so this RDR assigns it rather than
inheriting an assignment, with no overlap resulting.

Concerns that do not apply are omitted rather than N/A-bulleted: versioning,
build tool compatibility, licensing, deployment model, IDE compatibility,
secret/credential lifecycle, memory management, character encoding.

## Proportionality

**Right-sized. No sections flagged for trimming, and no split warranted.**

**Contract count (the split test).** This RDR is the sole author of **one**
independent load-bearing contract: escape-row shape conformance — an escape row
bears no owned-state mutation. The eight `normative` blocks are not eight
contracts; they are one contract plus the surface an implementer needs to hold
it: the predicate itself, its two enforcement points (which the RDR requires be
literally the same function, precisely so they cannot become two contracts),
the typed error carrying it, the multi-breach report shape, the vacuous-empty
case, the CLI carrier, and the explicit statement that the `Row` type does not
split. An implementer holds one invariant in working memory. The CLI envelope
field is additive under RDR 0005's ownership, not sole-authored here; the
authored-path half is RDR 0002's to implement. No second seam is being locked
alongside the first.

**Profile re-validation.** `foundational` is correct and unchanged. The
contract axis puts it there directly — this is a cross-RDR producer obligation
spanning modules: it binds every constructor of `resolve.Row` values, is
enforced in RDR 0001's kernel *and* RDR 0002's normalizer, and is consumed at
RDR 0005's CLI surface. The accretion axis does not apply (`Seam Lineage`
records no prior accretion), so no hard floor is in play; the contract axis
alone carries it. The form is correct: value plus one clause naming the
contract, with no matrix or provenance prose left over from the template.

**Lenses agree with the Profile.** The foundational row requires cove, 3amigo,
critique, and repeatability, and all four ran with evidence on disk — including
the full repeatability variant (`run-1`, `run-2`, `run-3`, `diff.md`) that the
contract's identity/ordering surface demands. No lens is missing, and no lens
disagrees with the sizing. The latch's backstop finds nothing to correct.

**Length is earned, not bloat.** The document is long, but the length sits
where the decisions are contested and expensive to reverse: the multi-breach
ordering clause is long because identity does not totally order distinct rows
and the naive fixes (positional tiebreak, dropping `Count`) each violate a
frozen test or force the prose re-parsing the RDR forbids; scenario 6 is long
because it records a mutation result that *refutes* the convenient claim rather
than hiding behind it. Trimming either would remove the reasoning that keeps an
implementer from re-deriving the wrong answer. The prose is proportionate to a
cross-module invariant with three enforcement points and a permanently exported
Go surface.

---

**Gate verdict: READY.** No blockers. Mechanical sweep PASS, ten terminal
assumptions, MVV in scope with non-vacuity proven, Profile re-validated against
the lenses that ran. Locking to Final.
