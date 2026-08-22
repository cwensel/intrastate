Model: claude-opus-5[1m]

# Grounding Findings — iteration 2 (re-entry)

Scope: delta against the propose/refine/resolve commits that produced this
re-entry (`3fea780`, `8b76d62`), i.e. every codebase/peer claim added or edited
since `Status: Final` at `6ccd164`. Iteration-1 grounding
(`../findings.md`) covered the pre-revision text; confirmed claims there are
not re-swept.

## REFUTED

### G2-1 — §JD-4 names RDR 0006 as the narrowing document, not this RDR

The RDR now claims JDR 0001 §JD-4 as the authority for a narrowing rule it
states about **its own** lint promise, in four places:

- Technical Design: "Where this exhaustiveness proof and RDR 0007's aggregation
  veto disagree, the **promise narrows** (JDR 0001 §JD-4)".
- A new `normative` clause: "An exhaustiveness claim MUST NOT be stronger than
  the runtime it describes … Where the two disagree the lint promise narrows".
- Decision Rationale: "narrows its own lint promise under §JD-4".
- Finalization Gate / Contradiction Check: "the strength of the exhaustiveness
  claim: where this RDR's proof and RDR 0007's aggregation veto disagree, the
  lint promise narrows and the runtime veto stands (§JD-4)".

What §JD-4 actually says (`docs/jdr/0001-resolve-kernel-seam.md:238-241`,
verbatim):

> **JD-4 Lint's promise is what gives.** Where **0006's** exhaustiveness proof
> and 0007's veto disagree, the *promise* narrows — P5 decides the substance.
> **Open only as to which document records the narrowing** and whether lint
> gains a warning category. *(0007×0006 F2, 0007×0003 F3)*

Two mismatches:

1. **Party.** §JD-4's head clause names *0006's* proof, not 0003's. RDR 0007
   confirms this reading in its own Predecessors block
   (`0007-guard-predicate-totality.md:159-161`): "**RDR 0006** (`Final`,
   tolerance §JD-4) narrows lint's promise where this RDR's veto and its proof
   disagree." RDR 0003 is listed there separately and is *not* the §JD-4
   tolerance-holder. RDR 0006 carries the qualifier on its own Status line:
   `Final [joint decision → JDR 0001 §JD-4]`
   (`0006-graph-lint-authority-and-guarantees.md:9`).
2. **Openness.** §JD-4 is an **open** ledger entry whose only remaining
   question is *which document records the narrowing*. RDR 0007 A12 restates
   this explicitly: "§JD-4 is not among the Closed entries; its SUBSTANCE is
   already decided … leaving open only which document records the narrowing and
   whether lint gains a warning category."

The RDR cites §JD-4 as if it were a closed decision assigning the narrowing to
0003. It is neither closed nor assigned here.

**Not a contradiction of substance.** §JD-4's substance ("the promise narrows;
the runtime veto stands") is exactly what the RDR wrote, and it is a defensible
claimant — the narrowing constrains a proof this RDR defines. What is refuted is
only the *citation*: the RDR presents a settled attribution that the cited
source records as open and currently points at a different document. The
citation, not the clause, is what needs to change.

**Live consequence — the sibling now disagrees.** RDR 0006 is `Final` and owns
graph-lint authority, but its exhaustiveness clause
(`0006-graph-lint-authority-and-guarantees.md:344-350`) narrows only for the
*non-finite-dimension* case:

> Graph lint MAY claim exhaustiveness only over finite declared domains supplied
> by the predicate/tag model. If a required dimension is not finite, lint MUST
> emit a blocking inability-to-prove finding for any contract that depends on
> closed coverage.

It carries **no** `guard_unevaluable`-veto narrowing — the case where every
dimension *is* finite but a participating row can still refuse at runtime. So
after this re-entry, 0003 states a narrowing rule that the `Final` document
§JD-4 actually tolerance-qualifies does not carry. Whichever document records
it, the two must agree; today they do not.

## NOT-FOUND

None. Every cited `path::Symbol` resolves on `main`.

## New rule with an existing sibling

None. The inverse search found no sibling path already deciding either new rule
(the §JD-4 narrowing, or A7's presence projection) — see Inverse search below.

## CONFIRMED

### Codebase claims (all resolve on `main`)

- `internal/resolve/resolve.go::GuardEvaluator` — declared at
  `resolve.go:91`; called via `resolve.go::evaluateGuard` (`:425`) from
  `resolve.go::gate` (`:376`). Investigation's claim that the kernel "calls the
  guard seam this RDR owns … but supplies no implementation of it" is
  confirmed: the only implementation in the repo is `fixtureGuards` in
  `internal/resolve/fixtures_test.go`; no non-test type satisfies the
  interface.
- `internal/resolve/resolve.go::Row` — carries `RuleID` and `SourceLocator`
  (`:173-174`), with the doc comment naming them as "the source identity RDR
  0002 requires every normalized row to retain". Confirms A5's shipped-identity
  claim.
- `Row.Guard string` (`resolve.go:185`) — confirms A5's honest caveat that the
  shipped kernel still carries an opaque predicate string and that 0007's
  reshape is specified, not implemented.
- `internal/cli/clierr/clierr.go::CLIError` — carries `Code`, `Message`,
  `Param`, `Detail`, `Hint`, `Group` (`:48-61`). A4 confirmed.
- `internal/cli/respond/respond.go::Fail` (`:133`). A4 confirmed.
- `internal/cli/config/config.go::Load` (`:74`). A4 confirmed.

Noted, not a finding: the shipped seam signature is
`Evaluate(guard string, view TagSet)` — it *does* receive the view, contrary to
the RDR's "never sees the tag view". The RDR attributes that clause to 0007's
specified reshape and states plainly that the reshape is "specified, not yet
implemented" and that its own implementation "sequences after it". The claim is
about the target seam, not the shipped one, and is correctly qualified.

### Peer-RDR / JDR claims

- **JDR 0001 §D1** fixes the parsed-atom row shape — confirmed
  (`0001-resolve-kernel-seam.md:86-90`, option (d) recommended): "Row carries a
  slice of atoms (key, operator token, literal, block) instead of a string. No
  reconstruction step". Matches A5's wording, including "no reconstruction step
  that could lose identity".
- **JDR 0001 §D4** resolves kernel-enforced guard domain — confirmed
  (`:189-208`, "**Resolved: (b)**"): the kernel "decides presence
  (provenance-blind), decides `exists` from presence alone, marks an absent-key
  value atom unevaluable **without calling the evaluator**", and "the evaluator
  decides value comparisons over present keys and never sees the view". Matches
  the RDR's split-enforcement paragraph and its new evaluator-scope normative
  clause.
- **RDR 0007 is `Final`** — confirmed (`0007-guard-predicate-totality.md:9`).
- **RDR 0007 fixes `exists` as the sole TOTAL operator, verdict
  `presence == literal`** — confirmed verbatim (`:1336-1344`).
- **RDR 0007 A12 routes the projection question here, `DOWNGRADED`, with the
  named plan "it lands when 0003 states its existence-atom projection"** —
  confirmed (`:586-604`).
- **RDR 0007 Phase 4 hands over six sibling items** — confirmed (`:2169-2183`):
  the existence token and `presence == literal` rule; the placement-MUST
  reading for `exists = false`; the two-row absence pattern and its conjoined
  value row; the A12/§JD-4 projection; the predicate-PLACEMENT guidance;
  sentinel-stamping as the anti-pattern; plus the set-valued element-encoding
  REQUEST. A7's "six sibling items" is a fair count of the hand-over list.
- **RDR 0002's tag declaration has no optionality field** — confirmed
  (`0002-transition-table-as-reviewable-data.md:217-218`): "Tag declarations:
  tag name, provenance (`owned`, `observed`, `recognized`), value kind, and
  optional accessor reference for observed or owned read-back." No optionality
  field. A7's blocker and the new Capability Dependencies row are accurate.
- **RDR 0002 is itself `Draft`** — confirmed (`:9`, re-entry re-verifying A2,
  A7).
- **RDR 0006 consumes source rule ids/spans and finite-domain semantics** —
  confirmed (`0006-…:105-121`, A1/A2 Peer RDR records citing 0003 A5 and A2).

### Spike / fixture claims

- `evidence/spikes/guard-fixture.toml` rule `foundational-to-cove` carries
  `[rule.guard.all.cluster_eligible] exists = true` over a tag declared
  `kind = "bool"`, `domain = [true, false]` — confirmed. A7's "If wrong" claim
  that refusing every `exists`-bearing group "disqualifies this RDR's own
  representative fixture row `foundational-to-cove`" is accurate, and the
  derivation's complete-domain/unprojectable-atom analysis holds against the
  file.
- The fixture's other six tags (`status`, `profile`, `stage`, `lens`,
  `prelock_iterations`, `rewind_target`) are matched or written on every path —
  confirmed, supporting the derivation's per-key (not global) presence-axis
  argument.

## Inverse search

The re-entry adds two new decision rules. Sibling-path check for each:

```sh
rg -n 'JD-4|narrow|guard_unevaluable|veto' docs/rdr/0006-graph-lint-authority-and-guarantees.md
rg -n 'optional|presence|absent|exists' docs/rdr/0002-transition-table-as-reviewable-data.md
rg -n 'Guard|presence|absent' internal/resolve/resolve.go
```

- **The §JD-4 narrowing rule** — searched RDR 0006 (the lint-authority owner and
  the §JD-4 tolerance-holder): it decides the *non-finite-dimension* refusal but
  **not** the `guard_unevaluable`-veto case. No sibling already makes this
  decision. The rule is genuinely unclaimed — which is precisely what §JD-4
  leaves open. This is the substance behind finding G2-1: not a duplicate rule,
  but a rule written here under a citation that points elsewhere.
- **A7's presence-dimension projection** — searched RDR 0002 (declaration
  producer), RDR 0006 (product/grouping consumer), RDR 0007 (totality owner),
  and `internal/resolve`. Searched, none exists: 0002 carries no optionality
  field, 0006 defers the product semantics to 0003, 0007 A12 explicitly routes
  the projection here, and the kernel has no presence-axis concept. A7 is
  correctly scoped as this RDR's to state.
